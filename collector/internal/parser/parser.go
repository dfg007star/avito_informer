package parser

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"

	"github.com/dfg007star/avito_informer/collector/internal/config"
	"github.com/dfg007star/avito_informer/collector/internal/model"
)

// ErrBlocked means Avito served a bot challenge / the results grid never
// appeared. The caller should back off, not retry immediately.
var ErrBlocked = errors.New("avito: results grid did not load (blocked or challenged)")

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 Edg/132.0.0.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 YaBrowser/132.0.0.0 Yowser/2.5 Safari/537.36",
}

// antiBotScript is injected on every new document. Kept minimal — the real
// defence against PerimeterX is the persistent logged-in profile plus a slow,
// jittered request cadence, not navigator spoofing.
const antiBotScript = `
	Object.defineProperty(navigator, 'webdriver', { get: () => undefined });
	Object.defineProperty(navigator, 'platform', { get: () => 'Win32' });
	Object.defineProperty(navigator, 'vendor', { get: () => 'Google Inc.' });
	window.chrome = { runtime: {} };
	Object.defineProperty(navigator, 'plugins', { get: () => [1, 2, 3] });
	Object.defineProperty(navigator, 'languages', { get: () => ['ru-RU', 'ru', 'en-US', 'en'] });
`

// Parser owns a single long-lived Chromium process and one browser context
// (one tab), reused for every link. The context is rebuilt by Recycle to bound
// Chromium's memory growth; the underlying process and the on-disk profile
// survive a recycle.
type Parser struct {
	allocatorCtx    context.Context
	cancelAllocator context.CancelFunc

	mu         sync.Mutex
	browserCtx context.Context
	cancelBrowser context.CancelFunc
}

// New builds the ExecAllocator (the Chromium process options) but does not
// start the browser. Call Start for that.
func New() *Parser {
	cfg := config.AppConfig().Parser

	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.Flag("lang", "ru-RU"),
		chromedp.WindowSize(1366, 768),
		chromedp.UserDataDir(cfg.ProfileDir()),
		chromedp.UserAgent(userAgents[rand.Intn(len(userAgents))]),
	}
	if cfg.Headless() {
		opts = append(opts, chromedp.Flag("headless", true))
	} else {
		// Run headful on the Xvfb display. chromedp defaults to
		// ozone-platform=headless; force X11 so it uses the virtual display.
		opts = append(opts,
			chromedp.Flag("headless", false),
			chromedp.Flag("ozone-platform", "x11"),
		)
	}
	if p := cfg.Proxy(); p != "" {
		opts = append(opts, chromedp.ProxyServer(p))
	}

	allocatorCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)

	return &Parser{
		allocatorCtx:    allocatorCtx,
		cancelAllocator: cancel,
	}
}

// NewParser keeps the old constructor name working for existing call sites.
func NewParser() *Parser { return New() }

// Start launches Chromium and creates the reusable browser context, enabling
// network/fetch and installing the anti-bot script once. Safe to call again
// after Shutdown of the context (Recycle uses it).
func (p *Parser) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.startLocked(ctx)
}

func (p *Parser) startLocked(ctx context.Context) error {
	if p.browserCtx != nil {
		return nil
	}

	bctx, cancel := chromedp.NewContext(p.allocatorCtx, chromedp.WithLogf(log.Printf))

	// The browser is bound to the context of the first Run. That must be bctx
	// itself, not a WithTimeout child — cancelling the child would kill the
	// browser. Guard the cold start with a watchdog goroutine instead.
	done := make(chan error, 1)
	go func() {
		done <- chromedp.Run(bctx,
			network.Enable(),
			fetch.Enable(),
			chromedp.ActionFunc(func(c context.Context) error {
				blockImages(c)
				_, err := page.AddScriptToEvaluateOnNewDocument(antiBotScript).Do(c)
				return err
			}),
		)
	}()

	var err error
	select {
	case err = <-done:
	case <-time.After(90 * time.Second):
		err = errors.New("timed out")
	case <-ctx.Done():
		err = ctx.Err()
	}
	if err != nil {
		cancel()
		return fmt.Errorf("failed to start browser: %w", err)
	}

	p.browserCtx = bctx
	p.cancelBrowser = cancel
	return nil
}

// Recycle tears down the browser context and builds a fresh one, keeping the
// Chromium process and the profile. Called once per cycle.
func (p *Parser) Recycle(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cancelBrowser != nil {
		p.cancelBrowser()
		p.browserCtx = nil
		p.cancelBrowser = nil
	}
	return p.startLocked(ctx)
}

// Shutdown stops the browser context and the Chromium process.
func (p *Parser) Shutdown() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cancelBrowser != nil {
		p.cancelBrowser()
		p.browserCtx = nil
		p.cancelBrowser = nil
	}
	if p.cancelAllocator != nil {
		p.cancelAllocator()
	}
}

// Scrape navigates the existing tab to the link's URL and extracts the items
// shown. Returns ErrBlocked if the results grid never appears.
func (p *Parser) Scrape(ctx context.Context, link *model.Link) ([]*model.Item, error) {
	p.mu.Lock()
	bctx := p.browserCtx
	p.mu.Unlock()
	if bctx == nil {
		return nil, errors.New("parser: Start not called")
	}

	cfg := config.AppConfig().Parser
	maxRetries := cfg.GetItemsMaxRetry()

	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		navCtx, cancel := context.WithTimeout(bctx, cfg.NavTimeout())
		html, err := loadResultsPage(navCtx, link.Url)
		cancel()

		if err == nil {
			return extractItems(html, link)
		}

		lastErr = err
		log.Printf("scrape %s attempt %d/%d: %s", link.Name, i+1, maxRetries, err)
		if i < maxRetries-1 {
			if !sleepCtx(ctx, cfg.GetItemsRetryDelay()) {
				return nil, ctx.Err()
			}
		}
	}

	// Every attempt failed to render the grid within the timeout — treat as a
	// block/challenge so the caller backs off instead of hammering.
	return nil, fmt.Errorf("%w: %v", ErrBlocked, lastErr)
}

// LoggedIn reports whether the current session shows a logged-in Avito UI.
// ponytail: checks for the profile/avatar marker on the main page; if Avito
// changes that marker this returns a false negative — adjust the selector.
func (p *Parser) LoggedIn(ctx context.Context) (bool, error) {
	p.mu.Lock()
	bctx := p.browserCtx
	p.mu.Unlock()
	if bctx == nil {
		return false, errors.New("parser: Start not called")
	}

	navCtx, cancel := context.WithTimeout(bctx, config.AppConfig().Parser.NavTimeout())
	defer cancel()

	var loggedIn bool
	err := chromedp.Run(navCtx,
		chromedp.Navigate("https://www.avito.ru/profile"),
		chromedp.WaitReady("body"),
		chromedp.Sleep(2*time.Second),
		chromedp.Evaluate(`!document.location.pathname.startsWith('/login') &&
			!!document.querySelector('[data-marker="header/username"], [data-marker="account-menu"]')`, &loggedIn),
	)
	return loggedIn, err
}

// Screenshot captures the current tab as PNG bytes.
func (p *Parser) Screenshot(ctx context.Context) ([]byte, error) {
	p.mu.Lock()
	bctx := p.browserCtx
	p.mu.Unlock()
	if bctx == nil {
		return nil, errors.New("parser: Start not called")
	}

	shotCtx, cancel := context.WithTimeout(bctx, 15*time.Second)
	defer cancel()

	var buf []byte
	if err := chromedp.Run(shotCtx, chromedp.CaptureScreenshot(&buf)); err != nil {
		return nil, err
	}
	return buf, nil
}

func blockImages(ctx context.Context) {
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		reqEvent, ok := ev.(*fetch.EventRequestPaused)
		if !ok {
			return
		}
		go func() {
			if reqEvent.ResourceType == network.ResourceTypeImage {
				_ = fetch.FailRequest(reqEvent.RequestID, network.ErrorReasonBlockedByClient).Do(ctx)
				return
			}
			_ = fetch.ContinueRequest(reqEvent.RequestID).Do(ctx)
		}()
	})
}

func sleepCtx(ctx context.Context, d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-t.C:
		return true
	}
}

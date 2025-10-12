package parser

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/dfg007star/avito_informer/collector/internal/config"
	"github.com/dfg007star/avito_informer/collector/internal/model"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 Edg/132.0.0.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 YaBrowser/132.0.0.0 Yowser/2.5 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Vivaldi/4.3.2439.65 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 OPR/80.0.4170.72 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Brave/1.32.106 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Iron/132.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Sleipnir/6.4.10 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Comodo_Dragon/132.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Maxthon/6.1.0.2000 Safari/537.36",
}

type Parser struct {
	AllocatorCtx    context.Context
	CancelAllocator context.CancelFunc
	proxyURL        string
}

func NewParser() *Parser {
	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("start-maximized", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.WindowSize(1920, 1080),
		chromedp.UserAgent(userAgents[rand.Intn(len(userAgents))]),
	}
	allocatorCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)

	return &Parser{
		AllocatorCtx:    allocatorCtx,
		CancelAllocator: cancel,
	}
}

func (p *Parser) SetProxies(proxyURLs []string) error {
	if len(proxyURLs) > 0 {
		p.proxyURL = proxyURLs[0] // chromedp only supports one proxy at a time
	} else {
		p.proxyURL = ""
	}

	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("start-maximized", true),
		chromedp.WindowSize(1920, 1080),
		chromedp.UserAgent(userAgents[rand.Intn(len(userAgents))]),
	}
	if p.proxyURL != "" {
		opts = append(opts, chromedp.ProxyServer(p.proxyURL))
	}

	if p.CancelAllocator != nil {
		p.CancelAllocator()
	}

	p.AllocatorCtx, p.CancelAllocator = chromedp.NewExecAllocator(context.Background(), opts...)

	return nil
}

func (p *Parser) Parse(link *model.Link, cookies map[string]string) ([]*model.Item, error) {
	var items []*model.Item
	priceRegex := regexp.MustCompile(`[^0-9] `)

	taskCtx, cancelTask := chromedp.NewContext(p.AllocatorCtx)
	defer cancelTask()

	taskCtx, cancelTimeout := context.WithTimeout(taskCtx, 60*time.Second)
	defer cancelTimeout()

	var err error
	maxRetries := config.AppConfig().Parser.GetItemsMaxRetry()
	for i := 0; i < maxRetries; i++ {
		var res string
		err = chromedp.Run(taskCtx,
			network.Enable(),
			fetch.Enable(),
			chromedp.ActionFunc(func(ctx context.Context) error {
				chromedp.ListenTarget(ctx, func(ev interface{}) {
					if reqEvent, ok := ev.(*fetch.EventRequestPaused); ok {
						go func() {
							if reqEvent.ResourceType == network.ResourceTypeImage {
								err := fetch.FailRequest(reqEvent.RequestID, network.ErrorReasonBlockedByClient).Do(ctx)
								if err != nil {
									log.Printf("failed to block request %s: %v", reqEvent.Request.URL, err)
								}
							} else {
								err := fetch.ContinueRequest(reqEvent.RequestID).Do(ctx)
								if err != nil {
									log.Printf("failed to continue request %s: %v", reqEvent.Request.URL, err)
								}
							}
						}()
					}
				})

				script := `
					Object.defineProperty(navigator, 'webdriver', { get: () => undefined });
					Object.defineProperty(navigator, 'platform', { get: () => 'Win32' });
					Object.defineProperty(navigator, 'vendor', { get: () => 'Google Inc.' });
					window.chrome = { runtime: {} };
					Object.defineProperty(navigator, 'plugins', { get: () => [1, 2, 3] });
					Object.defineProperty(navigator, 'languages', { get: () => ['en-US', 'en'] });
				`
				_, err := page.AddScriptToEvaluateOnNewDocument(script).Do(ctx)
				if err != nil {
					return err
				}

				for name, value := range cookies {
					expr := cdp.TimeSinceEpoch(time.Now().Add(7 * 24 * time.Hour))
					err := network.SetCookie(name, value).
						WithExpires(&expr).
						WithDomain("www.avito.ru").
						WithHTTPOnly(false).
						Do(ctx)
					if err != nil {
						return err
					}
				}
				return nil
			}),
			chromedp.Navigate(link.Url),
			chromedp.WaitVisible("div[data-marker='item']", chromedp.ByQuery),
			chromedp.OuterHTML("html", &res),
		)

		if err == nil {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
			if err != nil {
				return nil, fmt.Errorf("failed to parse HTML: %w", err)
			}

			doc.Find("div[data-marker='item']").Each(func(i int, s *goquery.Selection) {
				uid := s.AttrOr("data-item-id", "")
				priceStr := s.Find("meta[itemprop='price']").AttrOr("content", "")
				priceStr = priceRegex.ReplaceAllString(priceStr, "")
				price, _ := strconv.Atoi(priceStr)
				//var itemPreviewUrl string
				//if imgSelection := s.Find("img[itemprop='image']"); imgSelection.Length() > 0 {
				//	itemPreviewUrl = imgSelection.AttrOr("src", "")
				//}

				var imageUrls []string
				dataMarkerRegex := regexp.MustCompile(`slider-image/image-(https?://.*)`)
				s.Find("ul.photo-slider-list-R0jle li").Each(func(_ int, li *goquery.Selection) {
					dataMarker, exists := li.Attr("data-marker")
					if exists {
						matches := dataMarkerRegex.FindStringSubmatch(dataMarker)
						if len(matches) > 1 {
							imageUrls = append(imageUrls, matches[1])
						}
					}
				})

				item := &model.Item{
					LinkId:      link.ID,
					Uid:         uid,
					Title:       s.Find("h2[itemprop='name']").Text(),
					Price:       price,
					Description: s.Find("meta[itemprop='description']").AttrOr("content", ""),
					Url:         s.Find("a[itemprop='url']").AttrOr("href", ""),
					PreviewUrl:  imageUrls[0],
					IsNotify:    link.ItemsCount == 0,
				}
				items = append(items, item)
			})
			break
		}

		log.Printf("failed to visit %s (attempt %d/%d): %s", link.Name, i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(config.AppConfig().Parser.GetItemsRetryDelay())
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to visit url after %d attempts: %w", maxRetries, err)
	}

	return items, nil
}

func (p *Parser) GetCookies(ctx context.Context, urlStr string) (map[string]string, error) {
	taskCtx, cancelTask := chromedp.NewContext(p.AllocatorCtx)
	defer cancelTask()

	taskCtx, cancelTimeout := context.WithTimeout(taskCtx, config.AppConfig().Parser.GetCookieTimeout())
	defer cancelTimeout()

	log.Printf("attempting to get cookies from %s", urlStr)

	var cookies []*network.Cookie
	err := chromedp.Run(taskCtx,
		network.Enable(),
		chromedp.Navigate(urlStr),
		chromedp.WaitReady("body"),
		chromedp.Sleep(2*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			cookies, err = network.GetCookies().Do(ctx)
			return err
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get cookies: %w", err)
	}

	cookieMap := make(map[string]string)
	for _, cookie := range cookies {
		cookieMap[cookie.Name] = cookie.Value
	}

	return cookieMap, nil
}

func (p *Parser) Shutdown() {
	p.CancelAllocator()
}

func parseProxyURL(proxyURL string) (string, string, bool) {
	u, err := url.Parse(proxyURL)
	if err != nil {
		return "", "", false
	}
	if u.User != nil {
		password, _ := u.User.Password()
		return u.User.Username(), password, true
	}
	return "", "", false
}

package parser

import (
	"context"
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
	"github.com/chromedp/chromedp"
	"github.com/dfg007star/avito_informer/collector/internal/config"
	"github.com/dfg007star/avito_informer/collector/internal/model"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36",
}

type Parser struct {
	AllocatorCtx    context.Context
	CancelAllocator context.CancelFunc
	ProxyUser       string
	ProxyPass       string
}

func NewParser() *Parser {
	proxy := config.AppConfig().Parser.Proxy()
	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.WindowSize(1366, 768),
		chromedp.UserAgent(userAgents[rand.Intn(len(userAgents))]),
	}

	var proxyUser, proxyPass string
	if proxy != "" {
		if u, err := url.Parse(proxy); err == nil && u.User != nil {
			proxyUser = u.User.Username()
			proxyPass, _ = u.User.Password()
			u.User = nil
			proxy = u.String()
		}
		opts = append(opts, chromedp.ProxyServer(proxy))
		log.Printf("Parser initialized with proxy: %s", proxy)
	}

	allocatorCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)

	return &Parser{
		AllocatorCtx:    allocatorCtx,
		CancelAllocator: cancel,
		ProxyUser:       proxyUser,
		ProxyPass:       proxyPass,
	}
}

// https://github.com/chromedp/examples/blob/master/proxy/main.go
func (p *Parser) handleProxyAuth(ctx context.Context) {
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *fetch.EventRequestPaused:
			go func() {
				_ = chromedp.Run(ctx, fetch.ContinueRequest(ev.RequestID))
			}()
		case *fetch.EventAuthRequired:
			if p.ProxyUser != "" && ev.AuthChallenge.Source == fetch.AuthChallengeSourceProxy {
				log.Printf("Proxy auth required for %s. Providing credentials for user: %s", ev.AuthChallenge.Origin, p.ProxyUser)
				go func() {
					err := chromedp.Run(ctx, fetch.ContinueWithAuth(ev.RequestID, &fetch.AuthChallengeResponse{
						Response: fetch.AuthChallengeResponseResponseProvideCredentials,
						Username: p.ProxyUser,
						Password: p.ProxyPass,
					}))
					if err != nil {
						log.Printf("Failed to provide proxy auth: %v", err)
					}
				}()
			} else {
				log.Printf("Auth required (not proxy): %s", ev.AuthChallenge.Source)
				go func() {
					_ = chromedp.Run(ctx, fetch.ContinueWithAuth(ev.RequestID, &fetch.AuthChallengeResponse{
						Response: fetch.AuthChallengeResponseResponseDefault,
					}))
				}()
			}
		}
	})
}

func (p *Parser) Parse(link *model.Link, cookies map[string]string) ([]*model.Item, error) {
	taskCtx, cancelTask := chromedp.NewContext(p.AllocatorCtx)
	defer cancelTask()

	taskCtx, cancelTimeout := context.WithTimeout(taskCtx, 60*time.Second)
	defer cancelTimeout()

	p.handleProxyAuth(taskCtx)

	var res string
	err := chromedp.Run(taskCtx,
		fetch.Enable().WithHandleAuthRequests(true),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for name, value := range cookies {
				expr := cdp.TimeSinceEpoch(time.Now().Add(7 * 24 * time.Hour))
				_ = network.SetCookie(name, value).WithExpires(&expr).WithDomain("www.avito.ru").Do(ctx)
			}
			return nil
		}),
		chromedp.Navigate(link.Url),
		chromedp.WaitVisible("div[data-marker='item']", chromedp.ByQuery),
		chromedp.OuterHTML("html", &res),
	)

	if err != nil {
		return nil, err
	}

	return p.parseHTML(res, link)
}

func (p *Parser) GetCookies(urlStr string) (map[string]string, error) {
	taskCtx, cancelTask := chromedp.NewContext(p.AllocatorCtx)
	defer cancelTask()

	taskCtx, cancelTimeout := context.WithTimeout(taskCtx, 60*time.Second)
	defer cancelTimeout()

	p.handleProxyAuth(taskCtx)

	var cookies []*network.Cookie
	err := chromedp.Run(taskCtx,
		fetch.Enable().WithHandleAuthRequests(true),
		chromedp.Navigate(urlStr),
		chromedp.Sleep(5*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			cookies, err = network.GetCookies().Do(ctx)
			return err
		}),
	)
	if err != nil {
		return nil, err
	}

	cookieMap := make(map[string]string)
	for _, cookie := range cookies {
		cookieMap[cookie.Name] = cookie.Value
	}
	return cookieMap, nil
}

func (p *Parser) parseHTML(res string, link *model.Link) ([]*model.Item, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		return nil, err
	}

	var items []*model.Item
	priceRegex := regexp.MustCompile(`[^0-9]`)
	doc.Find("div[data-marker='item']").Each(func(i int, s *goquery.Selection) {
		uid := s.AttrOr("data-item-id", "")
		priceStr := priceRegex.ReplaceAllString(s.Find("meta[itemprop='price']").AttrOr("content", ""), "")
		price, _ := strconv.Atoi(priceStr)

		item := &model.Item{
			LinkId:      link.ID,
			Uid:         uid,
			Title:       s.Find("h2[itemprop='name']").Text(),
			Price:       price,
			Description: s.Find("meta[itemprop='description']").AttrOr("content", ""),
			Url:         s.Find("a[itemprop='url']").AttrOr("href", ""),
			IsNotify:    link.ItemsCount == 0,
		}
		items = append(items, item)
	})
	return items, nil
}

func (p *Parser) Shutdown() {
	p.CancelAllocator()
}

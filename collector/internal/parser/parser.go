package parser

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/dfg007star/avito_informer/collector/internal/model"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/110.0",
}

type Parser struct {
	collector *colly.Collector
}

func NewParser() *Parser {
	c := colly.NewCollector()
	c.SetRequestTimeout(30 * time.Second)
	return &Parser{
		collector: c,
	}
}

func (p *Parser) SetProxies(proxyURLs []string) error {
	rp, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		return err
	}
	p.collector.SetProxyFunc(rp)
	return nil
}

func (p *Parser) Parse(link *model.Link) ([]*model.Item, error) {
	var items []*model.Item
	priceRegex := regexp.MustCompile(`[^0-9] `)

	p.collector.OnRequest(func(r *colly.Request) {
		fmt.Printf("visiting %s\n", r.URL)
		r.Headers.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])
	})

	p.collector.OnHTML("div[data-marker='item']", func(e *colly.HTMLElement) {
		priceStr := e.ChildText("meta[itemprop='price']")
		priceStr = priceRegex.ReplaceAllString(priceStr, "")
		price, _ := strconv.Atoi(priceStr)

		item := &model.Item{
			LinkId:      link.ID,
			Title:       e.ChildText("h3[itemprop='name']"),
			Price:       price,
			Description: e.ChildText("meta[itemprop='description']"),
			Url:         e.ChildAttr("a[itemprop='url']", "href"),
		}
		items = append(items, item)
	})

	delay := time.Duration(1+rand.Intn(1)) * time.Second
	fmt.Printf("waiting for %s before visiting\n", delay)
	time.Sleep(delay)

	err := p.collector.Visit(link.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to visit url: %w", err)
	}

	p.collector.Wait()

	fmt.Printf("found %d items\n", len(items))

	return items, nil
}

package parser

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"

	"github.com/dfg007star/avito_informer/collector/internal/model"
)

var (
	priceRegex      = regexp.MustCompile(`[^0-9]`)
	dataMarkerRegex = regexp.MustCompile(`slider-image/image-(https?://.*)`)
)

// loadResultsPage navigates the tab to url and returns the rendered HTML once
// the results grid is visible. The caller supplies the timeout via ctx.
func loadResultsPage(ctx context.Context, url string) (string, error) {
	var html string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.WaitVisible("div[data-marker='item']", chromedp.ByQuery),
		chromedp.OuterHTML("html", &html),
	)
	if err != nil {
		return "", err
	}
	return html, nil
}

// extractItems parses the search-results HTML into items.
func extractItems(html string, link *model.Link) ([]*model.Item, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var items []*model.Item
	doc.Find("div[data-marker='item']").Each(func(_ int, s *goquery.Selection) {
		uid := s.AttrOr("data-item-id", "")

		priceStr := priceRegex.ReplaceAllString(s.Find("meta[itemprop='price']").AttrOr("content", ""), "")
		price, _ := strconv.Atoi(priceStr)

		var imageUrls []string
		s.Find("ul.photo-slider-list-R0jle li").Each(func(_ int, li *goquery.Selection) {
			if dm, ok := li.Attr("data-marker"); ok {
				if m := dataMarkerRegex.FindStringSubmatch(dm); len(m) > 1 {
					imageUrls = append(imageUrls, m[1])
				}
			}
		})
		var previewURL string
		if len(imageUrls) > 0 {
			previewURL = imageUrls[0]
		}

		items = append(items, &model.Item{
			LinkId:      link.ID,
			Uid:         uid,
			Title:       s.Find("h2[itemprop='name']").Text(),
			Price:       price,
			Description: s.Find("meta[itemprop='description']").AttrOr("content", ""),
			Url:         s.Find("a[itemprop='url']").AttrOr("href", ""),
			PreviewUrl:  previewURL,
			IsNotify:    link.ItemsCount == 0,
		})
	})

	return items, nil
}

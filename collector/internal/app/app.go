package app

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type App struct {
	diContainer *diContainer
}

func (a *App) Run(ctx context.Context) error {
	return a.collect(ctx)
}

func (a *App) collect(ctx context.Context) error {
	parser := a.diContainer.Parser()
	//proxyURL := "https://"
	//if err := parser.SetProxies([]string{proxyURL}); err != nil {
	//	return fmt.Errorf("failed to set proxies: %w", err)
	//}

	for {
		fmt.Println("starting new collection cycle")
		links, err := a.diContainer.Service(ctx).GetAllLinks(ctx)
		if err != nil {
			return fmt.Errorf("failed to get all links: %w", err)
		}

		fmt.Printf("found %d links to process\n", len(links))

		for _, link := range links {
			fmt.Printf("collecting items for link: %s\n", link.Url)
			items, err := parser.Parse(link)
			if err != nil {
				fmt.Printf("failed to parse link %s: %s\n", link.Url, err)
				continue
			}

			err = a.diContainer.Service(ctx).CreateItems(ctx, items)
			if err != nil {
				fmt.Printf("failed to create items for link %s: %s\n", link.Url, err)
				continue
			}

			delay := time.Duration(15+rand.Intn(15)) * time.Second
			fmt.Printf("waiting for %s before next link\n", delay)
			time.Sleep(delay)
		}

		fmt.Println("collection cycle finished, waiting for 30 seconds")
		time.Sleep(30 * time.Second)
	}
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

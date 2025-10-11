package app

import (
	"context"
	"fmt"
	"time"

	"github.com/dfg007star/avito_informer/collector/internal/config"
)

type App struct {
	diContainer *diContainer
}

func (a *App) Run(ctx context.Context) error {
	defer a.diContainer.Parser().Shutdown()
	return a.collect(ctx)
}

func (a *App) collect(ctx context.Context) error {
	parser := a.diContainer.Parser()
	//proxyURL := "http://LUXzR9:QV93K7@185.202.2.10:8000"
	//if err := parser.SetProxies([]string{proxyURL}); err != nil {
	//	return fmt.Errorf("failed to set proxies: %w", err)
	//}

	dummyAvitoURL := "https://www.avito.ru/moskva/telefony?q=iphone"
	initialCookies, err := parser.GetCookies(ctx, dummyAvitoURL)
	if err != nil {
		return fmt.Errorf("failed to get initial cookies: %w", err)
	}
	fmt.Printf("initial cookies obtained: %v", initialCookies)

	for {
		fmt.Println("starting new collection cycle")
		links, err := a.diContainer.Service(ctx).GetAllLinks(ctx)
		if err != nil {
			return fmt.Errorf("failed to get all links: %w", err)
		}

		for _, link := range links {
			fmt.Printf("collecting items for link name: %s", link.Name)
			items, err := parser.Parse(link, initialCookies)
			if err != nil {
				fmt.Printf("failed to parse link %s: %s", link.Name, err)
				continue
			}

			err = a.diContainer.Service(ctx).CreateItems(ctx, items)
			if err != nil {
				fmt.Printf("failed to create items for link %s: %s", link.Name, err)
				continue
			}

			delay := config.AppConfig().Parser.DelayBetweenLinks()
			time.Sleep(delay)
		}
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

package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dfg007star/avito_informer/collector/internal/config"
	"github.com/dfg007star/avito_informer/collector/internal/parser"
)

type App struct {
	diContainer *diContainer
}

func (a *App) Run(ctx context.Context) error {
	chromeInstance := a.diContainer.Parser()
	return a.collect(ctx, chromeInstance)
}

func (a *App) collect(ctx context.Context, chromeInstance *parser.Parser) error {
	dummyAvitoURL := "https://www.avito.ru/moskva/telefony?q=iphone"

	initialCookies, err := chromeInstance.GetCookies(dummyAvitoURL)
	if err != nil {
		return fmt.Errorf("failed to get initial cookies: %w", err)
	}
	log.Printf("initial cookies obtained: %v", initialCookies)

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		log.Println("starting new collection cycle")
		links, err := a.diContainer.Service(ctx).GetAllLinks(ctx)
		if err != nil {
			// Transient DB errors must not kill the process; log and retry next cycle.
			log.Printf("failed to get all links: %s", err)
			if !sleepCtx(ctx, config.AppConfig().Parser.DelayBetweenLinks()) {
				return ctx.Err()
			}
			continue
		}

		for _, link := range links {
			if err := ctx.Err(); err != nil {
				return err
			}

			log.Printf("collecting items for link name: %s", link.Name)

			items, err := chromeInstance.Parse(link, initialCookies)
			if err != nil {
				log.Printf("failed to parse link %s: %s", link.Name, err)
				continue
			}

			err = a.diContainer.Service(ctx).CreateItems(ctx, items)
			if err != nil {
				log.Printf("failed to create items for link %s: %s", link.Name, err)
				continue
			}

			if !sleepCtx(ctx, config.AppConfig().Parser.DelayBetweenLinks()) {
				return ctx.Err()
			}
		}

		// Always pause between cycles, even when there are no links, so an empty
		// links table does not busy-spin.
		if !sleepCtx(ctx, config.AppConfig().Parser.DelayBetweenLinks()) {
			return ctx.Err()
		}
	}
}

// sleepCtx sleeps for d or until ctx is cancelled. Returns false if cancelled.
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

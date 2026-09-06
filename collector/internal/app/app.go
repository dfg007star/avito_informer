package app

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/dfg007star/avito_informer/collector/internal/config"
	"github.com/dfg007star/avito_informer/collector/internal/parser"
)

type App struct {
	diContainer *diContainer
}

func (a *App) Run(ctx context.Context) error {
	p := a.diContainer.Parser()
	defer p.Shutdown()

	if err := p.Start(ctx); err != nil {
		return err
	}

	if loggedIn, err := p.LoggedIn(ctx); err != nil {
		log.Printf("login check failed: %s", err)
	} else if !loggedIn {
		log.Printf("WARNING: Avito session is not logged in. Bootstrap the profile over VNC.")
	} else {
		log.Printf("Avito session is logged in")
	}

	return a.collect(ctx, p)
}

func (a *App) collect(ctx context.Context, p *parser.Parser) error {
	// backoff steps applied after a block is detected; the final step repeats.
	backoff := []time.Duration{10 * time.Minute, 30 * time.Minute, time.Hour}
	blockedStreak := 0

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		links, err := a.diContainer.Service(ctx).GetAllLinks(ctx)
		if err != nil {
			log.Printf("failed to get all links: %s", err)
			if !sleepCtx(ctx, a.betweenLinks()) {
				return ctx.Err()
			}
			continue
		}

		log.Printf("cycle start: %d links", len(links))
		rand.Shuffle(len(links), func(i, j int) { links[i], links[j] = links[j], links[i] })

		cycleBlocked := false
		for _, link := range links {
			if err := ctx.Err(); err != nil {
				return err
			}

			items, err := p.Scrape(ctx, link)
			if err != nil {
				if errors.Is(err, parser.ErrBlocked) {
					log.Printf("BLOCKED on link %q: %s", link.Name, err)
					cycleBlocked = true
					break
				}
				log.Printf("failed to scrape link %q: %s", link.Name, err)
				continue
			}

			if err := a.diContainer.Service(ctx).CreateItems(ctx, items); err != nil {
				log.Printf("failed to create items for link %q: %s", link.Name, err)
			} else {
				log.Printf("link %q: %d items scraped", link.Name, len(items))
			}

			if !sleepCtx(ctx, a.betweenLinks()) {
				return ctx.Err()
			}
		}

		if cycleBlocked {
			d := backoff[min(blockedStreak, len(backoff)-1)]
			blockedStreak++
			log.Printf("backing off for %s (blocked streak %d)", d, blockedStreak)
			if !sleepCtx(ctx, d) {
				return ctx.Err()
			}
			continue
		}
		blockedStreak = 0

		if err := p.Recycle(ctx); err != nil {
			log.Printf("failed to recycle browser: %s", err)
		}

		if !sleepCtx(ctx, a.betweenLinks()) {
			return ctx.Err()
		}
	}
}

// betweenLinks returns a randomised pause. Falls back to DelayBetweenLinks when
// the jitter range is not configured.
func (a *App) betweenLinks() time.Duration {
	lo, hi := config.AppConfig().Parser.LinkDelayRange()
	if hi <= lo {
		return config.AppConfig().Parser.DelayBetweenLinks()
	}
	return lo + time.Duration(rand.Int63n(int64(hi-lo)))
}

func New(ctx context.Context) (*App, error) {
	a := &App{}
	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
	}
	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
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

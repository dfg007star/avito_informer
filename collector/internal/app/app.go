package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/stdlib"

	"github.com/dfg007star/avito_informer/collector/internal/config"
)

type App struct {
	diContainer *diContainer
}

func (a *App) Run(ctx context.Context) error {
	// here need run collctor

	return nil
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

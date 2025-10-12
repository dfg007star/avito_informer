package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/stdlib"

	"github.com/dfg007star/avito_informer/http/internal/config"
	pgMigrator "github.com/dfg007star/avito_informer/platform/pkg/migrator/pg"
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
}

func (a *App) Run(ctx context.Context) error {
	err := a.runHTTPServer(ctx)
	if err != nil {
		fmt.Errorf("http server crashed: %v", err)
	}

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
		a.initMigrator,
		a.initHTTPServer,
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

func (a *App) initMigrator(ctx context.Context) error {
	migrator := pgMigrator.New(
		stdlib.OpenDB(*a.diContainer.PostgresClient(ctx).Config().Copy()),
		config.AppConfig().Postgres.MigrationDirectory(),
	)
	err := migrator.Up()
	if err != nil {
		panic(fmt.Errorf("failed to migrate db: %w", err))
	}

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	linkHandler := a.diContainer.LinkHandler(ctx)
	authHandler := a.diContainer.AuthHandler(ctx)
	router := a.diContainer.LinkRouter(ctx, linkHandler, authHandler)

	a.httpServer = &http.Server{
		Addr:              config.AppConfig().HTTP.Address(),
		Handler:           router,
		ReadHeaderTimeout: config.AppConfig().HTTP.ReadTimeout(),
	}

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	err := a.httpServer.ListenAndServe()
	if err != nil {
		fmt.Errorf("failed to start http server: %w", err)
		return err
	}

	return nil
}

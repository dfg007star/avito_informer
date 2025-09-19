package main

import "fmt"

import (
	"context"
	"os/signal"
	"syscall"

	app "github.com/dfg007star/avito_informer/http/internal/app"
	"github.com/dfg007star/avito_informer/http/internal/config"
)

const configPath = "../deploy/compose/core/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()

	a, err := app.New(appCtx)
	if err != nil {
		fmt.Errorf("failed to create http app: %w", err)
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		fmt.Errorf("failed to run http app: %w", err)
		return
	}
}

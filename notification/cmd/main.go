package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/dfg007star/avito_informer/notification/internal/app"
	"github.com/dfg007star/avito_informer/notification/internal/config"
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
		fmt.Errorf("failed to create notification service: %w", err)
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		fmt.Errorf("failed to run notification service: %w", err)
		return
	}
}

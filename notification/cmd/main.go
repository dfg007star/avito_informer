package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dfg007star/avito_informer/notification/internal/app"
	"github.com/dfg007star/avito_informer/notification/internal/config"
)

func main() {
	var configPath string
	if os.Getenv("RUNNING_IN_DOCKER") == "true" {
		configPath = ""
	} else {
		configPath = "../deploy/compose/core/.env.local"
	}

	if err := config.Load(configPath); err != nil {
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

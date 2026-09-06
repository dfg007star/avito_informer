package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	app "github.com/dfg007star/avito_informer/collector/internal/app"
	"github.com/dfg007star/avito_informer/collector/internal/config"
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
		log.Fatalf("failed to create collector service: %s", err)
	}

	err = a.Run(appCtx)
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("collector service stopped with error: %s", err)
	}
}

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dfg007star/avito_informer/notification/internal/app"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	application, err := app.New(ctx)
	if err != nil {
		fmt.Printf("failed to create app: %v\n", err)
		os.Exit(1)
	}

	if err := application.Run(ctx); err != nil {
		fmt.Printf("app stopped with error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("app stopped gracefully")
}

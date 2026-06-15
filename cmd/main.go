package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/vikikurnia87/service-order/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a := app.New(ctx)
	defer a.Cleanup()

	a.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case <-quit:
		a.Logger.WarnContext(ctx, "🛑 Received shutdown signal...")
		cancel()
	case <-ctx.Done():
	}

	a.Shutdown()
	a.Logger.WarnContext(ctx, "👋 Process exiting")
}

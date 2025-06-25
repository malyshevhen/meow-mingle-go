package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/malyshEvhen/meow_mingle/internal/app"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	mingle, err := app.New(ctx)
	if err != nil {
		os.Exit(1)
	}

	go func() {
		if err := mingle.Start(ctx); err != nil {
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelShutdown()

	if err := mingle.Stop(shutdownCtx); err != nil {
		os.Exit(1)
	}
}

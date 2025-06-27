package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/malyshEvhen/meow_mingle/cmd/mingle"
	"github.com/malyshEvhen/meow_mingle/pkg/logger"
)

func main() {
	// Initialize logger first
	appLogger := logger.InitLogger()
	appLogger.LogStartup("meow-mingle", map[string]any{
		"version": "1.0.0",
		"env":     os.Getenv("ENV"),
	})

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	mingle, err := mingle.New(ctx)
	if err != nil {
		appLogger.Error("Failed to initialize application", "error", err.Error())
		os.Exit(1)
	}

	go func() {
		appLogger.Info("Starting HTTP server")
		if err := mingle.Start(ctx); err != nil {
			appLogger.Error("HTTP server failed", "error", err.Error())
			os.Exit(1)
		}
	}()

	appLogger.Info("Application started successfully")

	<-ctx.Done()
	appLogger.Info("Received shutdown signal")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelShutdown()

	appLogger.Info("Shutting down application")
	if err := mingle.Stop(shutdownCtx); err != nil {
		appLogger.Error("Failed to shutdown gracefully", "error", err.Error())
		os.Exit(1)
	}

	appLogger.LogShutdown("meow-mingle", "graceful shutdown completed")
}

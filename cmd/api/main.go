package main

import (
	"context"
	"log"
	"os"

	"github.com/malyshEvhen/meow_mingle/internal/app"
)

func main() {
	ctx := context.Background()

	shutDown, err := app.Start(ctx)
	if err != nil {
		log.Fatalf("Application failed to start: %s", err.Error())
	}
	defer func() {
		if err := shutDown(ctx); err != nil {
			log.Fatalf("Application failed to shut down: %s", err.Error())
		}
	}()

	os.Exit(0)
}

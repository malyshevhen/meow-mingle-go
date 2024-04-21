package app

import (
	"context"
	"fmt"
	"log"

	"github.com/malyshEvhen/meow_mingle/internal/config"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/internal/router"
	"github.com/malyshEvhen/meow_mingle/internal/server"
)

func Start(ctx context.Context) (func() error, error) {
	cfg := config.InitConfig()

	DB, err := db.NewDB(cfg)
	if err != nil {
		log.Fatalf("DB connection refused: %s", err.Error())
	}

	store := db.NewSQLStore(DB)
	mux := router.RegisterRoutes(store, cfg)
	if err := server.Serve(mux, cfg); err != nil {
		return nil, fmt.Errorf("an error occured while server starts: %s", err.Error())
	}

	return func() error {
		return DB.Close()
	}, nil
}

package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/malyshEvhen/meow_mingle/internal/config"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/internal/router"
)

const MIGRATION_SOURCE_URL string = "file://./db/migration"

func Start(ctx context.Context) (closerFunc func() error, appError error) {
	var (
		cfg       = config.InitConfig()
		DB        *sql.DB
		migration *migrate.Migrate
		store     *db.SQLStore
		mux       *mux.Router
	)

	DB, err := db.NewDB(cfg)
	if err != nil {
		log.Printf("%-15s ==> Database connection refused: %s\n", "Application", err.Error())
		appError = fmt.Errorf("database connection refused: %s", err.Error())
		return
	}
	log.Printf("%-15s ==> Database connection createt successfully", "Application")

	if err := DB.Ping(); err != nil {
		log.Printf("%-15s ==> Database is not reachable: %s\n", "Application", err.Error())
		appError = fmt.Errorf("database is not reachable: %s", err.Error())
		return
	}
	log.Printf("%-15s ==> Database connection is reachable", "Application")

	migration, err = migrate.New(MIGRATION_SOURCE_URL, cfg.DBSource)
	if err != nil {
		log.Printf("%-15s ==> Migration failed to prepare: %s\n", "Application", err.Error())
		appError = fmt.Errorf("migration configuration failed: %s", err.Error())
		return
	}
	log.Printf("%-15s ==> Migration configured successfully", "Application")

	if err := migration.Up(); err != nil {
		if err.Error() != "no change" {
			log.Printf("%-15s ==> Migration failed to apply: %s\n", "Application", err.Error())
			appError = fmt.Errorf("migration failed: %s", err.Error())
			return
		}
		log.Printf("%-15s ==> Migration not applied: %s", "Application", err.Error())
	} else {
		log.Printf("%-15s ==> Migration applied successfully", "Application")
	}

	store = db.NewSQLStore(DB)
	mux = router.RegisterRoutes(store, cfg)
	if err := http.ListenAndServe(cfg.ServerPort, mux); err != nil {
		log.Printf("%-15s ==> Server failed to start: %s\n", "Application", err.Error())
		appError = fmt.Errorf("an error occured while server starts: %s", err.Error())
		return
	}

	closerFunc = func() error {
		return DB.Close()
	}
	return
}

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
	_ "github.com/lib/pq"
	"github.com/malyshEvhen/meow_mingle/internal/config"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/internal/router"
	"github.com/malyshEvhen/meow_mingle/internal/server"
)

const MIGRATION_SOURCE_URL string = "file://./db/migration"

func Start(ctx context.Context) (closer func() error, err error) {
	var (
		cfg       = config.InitConfig()
		DB        *sql.DB
		migration *migrate.Migrate
		store     *db.SQLStore
		mux       *http.ServeMux
	)

	DB, err = db.NewDB(cfg)
	if err != nil {
		log.Printf("%-15s ==> Database connection refused: %s\n", "Application", err.Error())
		err = fmt.Errorf("database connection refused: %s", err.Error())
		return
	}
	log.Printf("%-15s ==> Database connection createt successfully", "Application")

	err = DB.Ping()
	if err != nil {
		log.Printf("%-15s ==> Database is not reachable: %s\n", "Application", err.Error())
		return
	}
	log.Printf("%-15s ==> Database connection is alive", "Application")

	migration, err = migrate.New(MIGRATION_SOURCE_URL, cfg.DBSource)
	if err != nil {
		log.Printf("%-15s ==> Migration failed to prepare: %s\n", "Application", err.Error())
		return
	}
	log.Printf("%-15s ==> Migration configured successfully", "Application")

	err = migration.Up()
	if err != nil {
		log.Printf("%-15s ==> Migration failed to apply: %s\n", "Application", err.Error())
		return
	}
	log.Printf("%-15s ==> Migration applied successfully", "Application")

	store = db.NewSQLStore(DB)
	mux = router.RegisterRoutes(store, cfg)
	err = server.Serve(mux, cfg)
	if err != nil {
		err = fmt.Errorf("an error occured while server starts: %s", err.Error())
		return
	}

	closer = func() error {
		return DB.Close()
	}
	return
}

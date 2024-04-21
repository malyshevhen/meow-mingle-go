package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/malyshEvhen/meow_mingle/internal/config"
)

func NewDB(cfg config.Config) (*sql.DB, error) {
	dbURL := cfg.DBSource
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return db, err
	}

	if err := db.Ping(); err != nil {
		return db, err
	} else {
		log.Println("Successfully Connected")
	}
	return db, err
}

package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/malyshEvhen/meow_mingle/internal/config"
)

type ConnectionPool struct {
	*sql.DB
	Err error
}

func NewDB() *ConnectionPool {
	dbURL := config.Envs.DBSource
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return &ConnectionPool{
			DB:  db,
			Err: err,
		}
	}

	if err := db.Ping(); err != nil {
		return &ConnectionPool{
			DB:  db,
			Err: err,
		}
	} else {
		log.Println("Successfully Connected")
	}
	return &ConnectionPool{
		DB:  db,
		Err: nil,
	}
}

package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

type MySQLDB struct {
	db *sql.DB
}

func NewMySQLStorage(config mysql.Config) *MySQLDB {
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatal("😱 Failed to open MySQL connection: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("😨 Failed to ping MySQL: ", err)
	}

	log.Println("🎉 Connected to the MySQL DB")

	return &MySQLDB{db: db}
}

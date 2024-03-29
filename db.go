package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

type MySQLStorage struct {
	db *sql.DB
}

func NewMySQLStorage(config mysql.Config) *MySQLStorage {
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatal("😱 Failed to open MySQL connection: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("😨 Failed to ping MySQL: ", err)
	}

	log.Println("🎉 Connected to the MySQL DB")

	return &MySQLStorage{db: db}
}

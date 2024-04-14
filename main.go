package main

import (
	"log"
	"os"

	"github.com/malyshEvhen/meow_mingle/application"
	db "github.com/malyshEvhen/meow_mingle/db/sqlc"
)

func main() {
	DB := db.NewDB()
	if err := DB.Err; err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	app, err := application.New(DB)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if err := app.Start(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	os.Exit(0)
}

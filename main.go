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
		log.Fatalf("Failed to initialized DB, due to: %s", err.Error())
	}

	app, err := application.New(DB)
	if err != nil {
		// TODO: errorf
		log.Fatal(err)
	}

	if err := app.Start(); err != nil {
		// TODO: errorf
		log.Fatal(err)
	}
	os.Exit(0)
}

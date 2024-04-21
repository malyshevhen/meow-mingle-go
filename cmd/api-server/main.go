package main

import (
	"log"
	"os"

	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/internal/db"
)

func main() {
	DB := db.NewDB()
	if err := DB.Err; err != nil {
		log.Fatalf("Failed to initialized DB, due to: %s", err.Error())
	}

	application, err := app.New(DB)
	if err != nil {
		// TODO: errorf
		log.Fatal(err)
	}

	if err := application.Start(); err != nil {
		// TODO: errorf
		log.Fatal(err)
	}
	os.Exit(0)
}

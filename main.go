package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	db "github.com/malyshEvhen/meow_mingle/db/sqlc"
)

func main() {
	DB, err := sql.Open("postgres", Envs.DBSource)
	if err != nil {
		log.Fatalln(err)
	}

	defer DB.Close()

	if err := DB.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully Connected")
	}

	store := db.NewStore(DB)

	server := NewApiServer(":8080", store)
	server.Serve()
}

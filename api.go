package main

import (
	"log"
	"net/http"

	db "github.com/malyshEvhen/meow_mingle/db/sqlc"
)

type ApiServer struct {
	store *db.Store
	addr  string
}

func NewApiServer(addr string, store *db.Store) *ApiServer {
	return &ApiServer{
		addr:  addr,
		store: store,
	}
}

func (s *ApiServer) Serve() {
	submuxer := http.NewServeMux()

	router := NewRouter(s.store)
	router.RegisterRoutes(submuxer)

	muxer := http.NewServeMux()
	muxer.Handle("/api/v1/", http.StripPrefix("/api/v1", submuxer))

	log.Printf("Server starting at port: %s\n", s.addr)

	log.Fatal(http.ListenAndServe(s.addr, muxer))
}

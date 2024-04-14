package api

import (
	"log"
	"net/http"

	db "github.com/malyshEvhen/meow_mingle/db/sqlc"
)

type Server struct {
	store *db.Store
	addr  string
}

func NewServer(addr string, store *db.Store) *Server {
	return &Server{
		addr:  addr,
		store: store,
	}
}

func (s *Server) Serve() error {
	submuxer := http.NewServeMux()

	router := NewRouter(s.store)
	router.RegisterRoutes(submuxer)

	muxer := http.NewServeMux()
	muxer.Handle("/api/v1/", http.StripPrefix("/api/v1", submuxer))

	log.Printf("Server starting at port: %s\n", s.addr)

	return http.ListenAndServe(s.addr, muxer)
}

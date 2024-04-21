package api

import (
	"log"
	"net/http"

	"github.com/malyshEvhen/meow_mingle/internal/db"
)

type Server struct {
	addr string
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (s *Server) Serve(store db.IStore) error {
	submuxer := http.NewServeMux()

	router := NewRouter(store)
	router.RegisterRoutes(submuxer)

	muxer := http.NewServeMux()
	muxer.Handle("/api/v1/", http.StripPrefix("/api/v1", submuxer))

	log.Printf("Server starting at port: %s\n", s.addr)

	return http.ListenAndServe(s.addr, muxer)
}

package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	addr  string
	store Store
}

func NewApiServer(addr string, store Store) *ApiServer {
	return &ApiServer{
		addr:  addr,
		store: store,
	}
}

func (s *ApiServer) Serve() {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userService := NewUserService(s.store)
	userService.RegisterRoutes(subrouter)

	postsService := NewPostService(s.store)
	postsService.RegisterRoutes(subrouter)

	log.Printf("Server starting at port: %s\n", s.addr)

	log.Fatal(http.ListenAndServe(s.addr, subrouter))
}

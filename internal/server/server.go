package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/malyshEvhen/meow_mingle/internal/config"
)

func Serve(muxer *mux.Router, cfg config.Config) error {
	log.Printf("%-15s ==> Starting at port: %s\n", "Server", cfg.ServerPort)

	return http.ListenAndServe(cfg.ServerPort, muxer)
}

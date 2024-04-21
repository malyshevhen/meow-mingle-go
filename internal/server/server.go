package server

import (
	"log"
	"net/http"

	"github.com/malyshEvhen/meow_mingle/internal/config"
)

func Serve(muxer *http.ServeMux, cfg config.Config) error {
	log.Printf("Server starting at port: %s\n", cfg.ServerPort)

	return http.ListenAndServe(cfg.ServerPort, muxer)
}

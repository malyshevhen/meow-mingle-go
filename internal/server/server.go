package server

import (
	"log"
	"net/http"

	"github.com/malyshEvhen/meow_mingle/internal/config"
)

func Serve(muxer *http.ServeMux, cfg config.Config) error {
	log.Printf("%-15s ==> Starting at port: %s\n", "Server", cfg.ServerPort)

	return http.ListenAndServe(cfg.ServerPort, muxer)
}

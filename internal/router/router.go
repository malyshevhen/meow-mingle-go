package router

import (
	"net/http"

	"github.com/malyshEvhen/meow_mingle/internal/config"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/internal/middleware"
)

func Authenticated(
	store db.IUserReposytory,
	cfg config.Config,
	handler func(w http.ResponseWriter, r *http.Request) error,
) http.HandlerFunc {
	return middleware.MiddlewareChain(
		handler,
		middleware.LoggerMW,
		middleware.ErrorHandler,
		middleware.WithJWTAuth(store, cfg),
	)
}

func Public(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return middleware.MiddlewareChain(
		handler,
		middleware.LoggerMW,
		middleware.ErrorHandler,
	)
}

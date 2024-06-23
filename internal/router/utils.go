package router

import (
	"net/http"

	"github.com/malyshEvhen/meow_mingle/internal/middleware"
	"github.com/malyshEvhen/meow_mingle/internal/types"
)

func Authenticated(
	handler types.Handler,
	authMW func(h types.Handler) types.Handler,
) http.HandlerFunc {
	return middleware.MiddlewareChain(
		handler,
		middleware.LoggerMW,
		middleware.ErrorHandler,
		authMW,
	)
}

func Public(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return middleware.MiddlewareChain(
		handler,
		middleware.LoggerMW,
		middleware.ErrorHandler,
	)
}

package router

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/malyshEvhen/meow_mingle/internal/config"
	"github.com/malyshEvhen/meow_mingle/internal/handlers"
	"github.com/malyshEvhen/meow_mingle/internal/middleware"
	"github.com/malyshEvhen/meow_mingle/internal/types"
)

type LoginRouter struct {
	loginHandler *handlers.LoginHandler
	authMW       *middleware.AuthProvider
}

func NewLoginHandler(
	loginHandler *handlers.LoginHandler,
	authMW *middleware.AuthProvider,
) *LoginRouter {
	return &LoginRouter{
		loginHandler: loginHandler,
		authMW:       authMW,
	}
}

func (lr *LoginRouter) RegisterRouts(ctx context.Context, mux *mux.Router, cfg config.Config) *mux.Router {
	loginMux := mux.PathPrefix("/comments").Subrouter()

	auth := func(handler types.Handler) http.HandlerFunc {
		return Authenticated(handler, lr.authMW.WithBasicAuth)
	}

	loginMux.HandleFunc("/login", auth(lr.loginHandler.HahdleLogin())).Methods("POST")

	return mux
}

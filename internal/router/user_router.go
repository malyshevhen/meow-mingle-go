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

type UserRouter struct {
	authMW      *middleware.AuthProvider
	userHandler *handlers.UserHandler
	postHandler *handlers.PostHandler
}

func NewUserRouter(
	authMW *middleware.AuthProvider,
	userHandler *handlers.UserHandler,
	postHandler *handlers.PostHandler,
) *UserRouter {
	return &UserRouter{
		authMW:      authMW,
		userHandler: userHandler,
		postHandler: postHandler,
	}
}

func (ur *UserRouter) RegisterRouts(ctx context.Context, mux *mux.Router, cfg config.Config) *mux.Router {
	usersMux := mux.PathPrefix("/users").Subrouter()

	auth := func(handler types.Handler) http.HandlerFunc {
		return Authenticated(handler, ur.authMW.WithJWTAuth)
	}

	usersMux.HandleFunc("/register", Public(ur.userHandler.HandleCreateUser(cfg))).Methods("POST")
	usersMux.HandleFunc("/{id}/feed", Public(ur.postHandler.HandleUsersFeed())).Methods("GET")
	usersMux.HandleFunc("/{id}/posts", Public(ur.postHandler.HandleGetUserPosts())).Methods("GET")
	usersMux.HandleFunc("/{id}", auth(ur.userHandler.HandleGetUser())).Methods("GET")
	usersMux.HandleFunc("/feed", auth(ur.postHandler.HandleOwnersFeed())).Methods("GET")
	usersMux.HandleFunc("/{id}/subscriptions", auth(ur.userHandler.HandleSubscribe())).Methods("POST")
	usersMux.HandleFunc("/{id}/subscriptions", auth(ur.userHandler.HandleUnsubscribe())).Methods("DELETE")

	return mux
}

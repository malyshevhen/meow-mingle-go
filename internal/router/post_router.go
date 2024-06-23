package router

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/malyshEvhen/meow_mingle/internal/config"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/internal/handlers"
	"github.com/malyshEvhen/meow_mingle/internal/middleware"
	"github.com/malyshEvhen/meow_mingle/internal/types"
)

type PostRouter struct {
	authMW         *middleware.AuthProvider
	postHandler    *handlers.PostHandler
	commentHandler *handlers.CommentHandler
	userRepo       db.IUserReposytory
}

func NewPostRouter(
	authMW *middleware.AuthProvider,
	postHandler *handlers.PostHandler,
	commentHandler *handlers.CommentHandler,
	userRepo db.IUserReposytory,
) *PostRouter {
	return &PostRouter{
		authMW:         authMW,
		postHandler:    postHandler,
		commentHandler: commentHandler,
		userRepo:       userRepo,
	}
}

func (pr *PostRouter) RegisterRouts(ctx context.Context, mux *mux.Router, cfg config.Config) *mux.Router {
	postsMux := mux.PathPrefix("/posts").Subrouter()

	auth := func(handler types.Handler) http.HandlerFunc {
		return Authenticated(handler, pr.authMW.WithJWTAuth)
	}

	postsMux.HandleFunc("/{id}", Public(pr.postHandler.HandleGetPostsById())).Methods("GET")
	postsMux.HandleFunc("/{id}/comments", Public(pr.commentHandler.HandleGetComments())).Methods("GET")
	postsMux.HandleFunc("", auth(pr.postHandler.HandleCreatePost())).Methods("POST")
	postsMux.HandleFunc("/{id}", auth(pr.postHandler.HandleUpdatePostsById())).Methods("PUT")
	postsMux.HandleFunc("/{id}", auth(pr.postHandler.HandleDeletePostsById())).Methods("DELETE")
	postsMux.HandleFunc("/{id}/likes", auth(pr.postHandler.HandleLikePost())).Methods("POST")
	postsMux.HandleFunc("/{id}/likes", auth(pr.postHandler.HandleRemoveLikeFromPost())).Methods("DELETE")
	postsMux.HandleFunc("/{id}/comments", auth(pr.commentHandler.HandleCreateComment())).Methods("POST")

	return mux
}

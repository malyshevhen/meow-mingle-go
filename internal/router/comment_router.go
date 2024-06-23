package router

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/malyshEvhen/meow_mingle/internal/config"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/internal/handlers"
	"github.com/malyshEvhen/meow_mingle/internal/types"
)

type CommentRouter struct {
	commentHandler *handlers.CommentHandler
	userRepo       db.IUserReposytory
}

func NewCommentRouter(
	userRepo db.IUserReposytory,
	commentHandler *handlers.CommentHandler,
) *CommentRouter {
	return &CommentRouter{
		commentHandler: commentHandler,
		userRepo:       userRepo,
	}
}

func (cr *CommentRouter) RegisterRouts(ctx context.Context, mux *mux.Router, cfg config.Config) *mux.Router {
	commentsMux := mux.PathPrefix("/comments").Subrouter()

	auth := func(handler types.Handler) http.HandlerFunc {
		return Authenticated(cr.userRepo, cfg, handler)
	}

	commentsMux.HandleFunc("/{id}", auth(cr.commentHandler.HandleUpdateComments())).Methods("PUT")
	commentsMux.HandleFunc("/{id}", auth(cr.commentHandler.HandleDeleteComments())).Methods("DELETE")
	commentsMux.HandleFunc("/{id}/likes", auth(cr.commentHandler.HandleLikeComment())).Methods("POST")
	commentsMux.HandleFunc("/{id}/likes", auth(cr.commentHandler.HandleRemoveLikeFromComment())).Methods("DELETE")

	return mux
}

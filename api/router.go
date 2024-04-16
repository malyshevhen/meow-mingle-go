package api

import (
	"net/http"

	db "github.com/malyshEvhen/meow_mingle/db/sqlc"
)

type Handler func(w http.ResponseWriter, r *http.Request) error

type Router struct {
	store *db.Store
}

func NewRouter(store *db.Store) *Router {
	return &Router{
		store: store,
	}
}

func (rr *Router) RegisterRoutes(mux *http.ServeMux) {
	authenticated := func(handler Handler) http.HandlerFunc {
		return MiddlewareChain(
			handler,
			LoggerMiddleware,
			ErrorHandler,
			WithJWTAuth(rr.store, handler),
		)
	}

	noAuth := func(handler Handler) http.HandlerFunc {
		return MiddlewareChain(
			handler,
			LoggerMiddleware,
			ErrorHandler,
		)
	}

	mux.HandleFunc("GET /users/{id}", authenticated(rr.handleGetUser))
	mux.HandleFunc("POST /users/register", noAuth(rr.handleCreateUser))
	mux.HandleFunc("GET /users/{id}/posts", noAuth(rr.handleGetUserPosts))
	mux.HandleFunc("POST /users/{id}/subscriptions", authenticated(rr.handleSubscribe))
	mux.HandleFunc("DELETE /users/{id}/subscriptions", authenticated(rr.handleUnsubscribe))
	mux.HandleFunc("GET /users/feed", authenticated(rr.handleOwnersFeed))
	mux.HandleFunc("GET /users/{id}/feed", noAuth(rr.handleUsersFeed))

	mux.HandleFunc("POST /posts", authenticated(rr.handleCreatePost))
	mux.HandleFunc("POST /posts/{id}/likes", authenticated(rr.handleLikePost))
	mux.HandleFunc("POST /posts/{id}/comments", authenticated(rr.handleCreateComment))
	mux.HandleFunc("PUT /posts/{id}", authenticated(rr.handleUpdatePostsById))
	mux.HandleFunc("DELETE /posts/{id}", authenticated(rr.handleDeletePostsById))
	mux.HandleFunc("DELETE /posts/{id}/likes", authenticated(rr.handleRemoveLikeFromPost))
	mux.HandleFunc("GET /posts/{id}", noAuth(rr.handleGetPostsById))
	mux.HandleFunc("GET /posts/{id}/comments", noAuth(rr.handleGetComments))

	mux.HandleFunc("PUT /comments/{id}", authenticated(rr.handleUpdateComments))
	mux.HandleFunc("POST /comments/{id}/likes", authenticated(rr.handleLikeComment))
	mux.HandleFunc("DELETE /comments/{id}", authenticated(rr.handleDeleteComments))
	mux.HandleFunc("DELETE /comments/{id}/likes", authenticated(rr.handleRemoveLikeFromComment))
}

package api

import (
	"net/http"

	db "github.com/malyshEvhen/meow_mingle/db/sqlc"
)

type Handler func(w http.ResponseWriter, r *http.Request) error

type Router struct {
	store db.IStore
}

func NewRouter(store db.IStore) *Router {
	return &Router{
		store: store,
	}
}

func (rr *Router) RegisterRoutes(mux *http.ServeMux) {
	authenticated := func(handler Handler) http.HandlerFunc {
		return MiddlewareChain(
			handler,
			LoggerMW,
			ErrorHandler,
			WithJWTAuth(rr.store),
		)
	}

	noAuth := func(handler Handler) http.HandlerFunc {
		return MiddlewareChain(
			handler,
			LoggerMW,
			ErrorHandler,
		)
	}

	mux.HandleFunc("GET /users/{id}", authenticated(handleGetUser(rr.store)))
	mux.HandleFunc("POST /users/register", noAuth(handleCreateUser(rr.store)))
	mux.HandleFunc("POST /users/{id}/subscriptions", authenticated(handleSubscribe(rr.store)))
	mux.HandleFunc("DELETE /users/{id}/subscriptions", authenticated(handleUnsubscribe(rr.store)))
	mux.HandleFunc("GET /users/feed", authenticated(handleOwnersFeed(rr.store)))
	mux.HandleFunc("GET /users/{id}/feed", noAuth(handleUsersFeed(rr.store)))

	mux.HandleFunc("GET /users/{id}/posts", noAuth(handleGetUserPosts(rr.store)))
	mux.HandleFunc("POST /posts", authenticated(handleCreatePost(rr.store)))
	mux.HandleFunc("POST /posts/{id}/likes", authenticated(handleLikePost(rr.store)))
	mux.HandleFunc("PUT /posts/{id}", authenticated(handleUpdatePostsById(rr.store)))
	mux.HandleFunc("DELETE /posts/{id}", authenticated(handleDeletePostsById(rr.store)))
	mux.HandleFunc("DELETE /posts/{id}/likes", authenticated(handleRemoveLikeFromPost(rr.store)))
	mux.HandleFunc("GET /posts/{id}", noAuth(handleGetPostsById(rr.store)))

	mux.HandleFunc("POST /posts/{id}/comments", authenticated(handleCreateComment(rr.store)))
	mux.HandleFunc("GET /posts/{id}/comments", noAuth(handleGetComments(rr.store)))
	mux.HandleFunc("PUT /comments/{id}", authenticated(handleUpdateComments(rr.store)))
	mux.HandleFunc("POST /comments/{id}/likes", authenticated(handleLikeComment(rr.store)))
	mux.HandleFunc("DELETE /comments/{id}", authenticated(handleDeleteComments(rr.store)))
	mux.HandleFunc(
		"DELETE /comments/{id}/likes",
		authenticated(handleRemoveLikeFromComment(rr.store)),
	)
}

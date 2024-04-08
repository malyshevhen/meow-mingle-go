package main

import (
	"net/http"

	db "github.com/malyshEvhen/meow_mingle/db/sqlc"
)

type Router struct {
	store *db.Store
}

func NewRouter(store *db.Store) *Router {
	return &Router{
		store: store,
	}
}

func (r *Router) RegisterRoutes(mux *http.ServeMux) {
	middlewareChain := func(handler apiHandler) http.HandlerFunc {
		return MiddlewareChain(
			handler,
			LoggerMiddleware,
			ErrorHandler,
			WithJWTAuth(r.store, handler),
		)
	}

	mux.HandleFunc("POST /users/register", r.handleCreateUser)
	mux.HandleFunc("GET /users/{id}", middlewareChain(r.handleGetUser))
	mux.HandleFunc("GET /users/{id}/posts", middlewareChain(r.handleGetUserPosts))
	mux.HandleFunc("POST /posts", middlewareChain(r.handleCreatePost))
	mux.HandleFunc("GET /posts/{id}", middlewareChain(r.handleGetPostsById))
	mux.HandleFunc("PUT /posts/{id}", middlewareChain(r.handleUpdatePostsById))
	mux.HandleFunc("DELETE /posts/{id}", middlewareChain(r.handleDeletePostsById))
	mux.HandleFunc("POST /posts/{id}/comments", middlewareChain(r.handleCreateComment))
	mux.HandleFunc("GET /posts/{id}/comments", middlewareChain(r.handleGetComments))
	mux.HandleFunc("PUT /comments/{id}", middlewareChain(r.handleUpdateComments))
	mux.HandleFunc("DELETE /comments/{id}", middlewareChain(r.handleDeleteComments))
	mux.HandleFunc("POST /posts/{id}/likes", middlewareChain(r.handleLikePost))
	mux.HandleFunc("DELETE /posts/{id}/likes", middlewareChain(r.handleRemoveLikeFromPost))
	mux.HandleFunc("POST /comments/{id}/likes", middlewareChain(r.handleLikeComment))
	mux.HandleFunc("DELETE /comments/{id}/likes", middlewareChain(r.handleRemoveLikeFromComment))
}

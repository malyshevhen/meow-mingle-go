package main

import (
	"net/http"

	db "github.com/malyshEvhen/meow_mingle/db/sqlc"
)

type ApiHandler func(w http.ResponseWriter, r *http.Request) error

type Router struct {
	store *db.Store
}

func NewRouter(store *db.Store) *Router {
	return &Router{
		store: store,
	}
}

func (r *Router) RegisterRoutes(mux *http.ServeMux) {
	var (
		userMux    = http.NewServeMux()
		postMux    = http.NewServeMux()
		commentMux = http.NewServeMux()
	)

	authenticated := func(handler ApiHandler) http.HandlerFunc {
		return MiddlewareChain(
			handler,
			LoggerMiddleware,
			ErrorHandler,
			WithJWTAuth(r.store, handler),
		)
	}

	noAuth := func(handler ApiHandler) http.HandlerFunc {
		return MiddlewareChain(
			handler,
			LoggerMiddleware,
			ErrorHandler,
		)
	}

	userMux.HandleFunc("GET /{id}", authenticated(r.handleGetUser))
	userMux.HandleFunc("POST /register", noAuth(r.handleCreateUser))
	userMux.HandleFunc("GET /{id}/posts", noAuth(r.handleGetUserPosts))
	mux.Handle("/users/", http.StripPrefix("/users", userMux))

	postMux.HandleFunc("POST /", authenticated(r.handleCreatePost))
	postMux.HandleFunc("POST /{id}/likes", authenticated(r.handleLikePost))
	postMux.HandleFunc("POST /{id}/comments", authenticated(r.handleCreateComment))
	postMux.HandleFunc("PUT /{id}", authenticated(r.handleUpdatePostsById))
	postMux.HandleFunc("DELETE /{id}", authenticated(r.handleDeletePostsById))
	postMux.HandleFunc("DELETE /{id}/likes", authenticated(r.handleRemoveLikeFromPost))
	postMux.HandleFunc("GET /{id}", noAuth(r.handleGetPostsById))
	postMux.HandleFunc("GET /{id}/comments", noAuth(r.handleGetComments))
	mux.Handle("/posts/", http.StripPrefix("/posts", postMux))

	commentMux.HandleFunc("PUT /{id}", authenticated(r.handleUpdateComments))
	commentMux.HandleFunc("POST /{id}/likes", authenticated(r.handleLikeComment))
	commentMux.HandleFunc("DELETE /{id}", authenticated(r.handleDeleteComments))
	commentMux.HandleFunc("DELETE /{id}/likes", authenticated(r.handleRemoveLikeFromComment))
	mux.Handle("/comments/", http.StripPrefix("/comments", commentMux))
}

package router

import (
	"net/http"

	"github.com/malyshEvhen/meow_mingle/internal/api"
	"github.com/malyshEvhen/meow_mingle/internal/config"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/internal/middleware"
)

func RegisterRoutes(store db.IStore, cfg config.Config) *http.ServeMux {
	mux := http.NewServeMux()

	authenticated := func(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
		return middleware.MiddlewareChain(
			handler,
			middleware.LoggerMW,
			middleware.ErrorHandler,
			middleware.WithJWTAuth(store, cfg),
		)
	}

	noAuth := func(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
		return middleware.MiddlewareChain(
			handler,
			middleware.LoggerMW,
			middleware.ErrorHandler,
		)
	}

	mux.HandleFunc("GET /users/{id}", authenticated(api.HandleGetUser(store)))
	mux.HandleFunc("POST /users/register", noAuth(api.HandleCreateUser(store, cfg)))
	mux.HandleFunc("POST /users/{id}/subscriptions", authenticated(api.HandleSubscribe(store)))
	mux.HandleFunc("DELETE /users/{id}/subscriptions", authenticated(api.HandleUnsubscribe(store)))
	mux.HandleFunc("GET /users/feed", authenticated(api.HandleOwnersFeed(store)))
	mux.HandleFunc("GET /users/{id}/feed", noAuth(api.HandleUsersFeed(store)))

	mux.HandleFunc("GET /users/{id}/posts", noAuth(api.HandleGetUserPosts(store)))
	mux.HandleFunc("POST /posts", authenticated(api.HandleCreatePost(store)))
	mux.HandleFunc("POST /posts/{id}/likes", authenticated(api.HandleLikePost(store)))
	mux.HandleFunc("PUT /posts/{id}", authenticated(api.HandleUpdatePostsById(store)))
	mux.HandleFunc("DELETE /posts/{id}", authenticated(api.HandleDeletePostsById(store)))
	mux.HandleFunc("DELETE /posts/{id}/likes", authenticated(api.HandleRemoveLikeFromPost(store)))
	mux.HandleFunc("GET /posts/{id}", noAuth(api.HandleGetPostsById(store)))

	mux.HandleFunc("POST /posts/{id}/comments", authenticated(api.HandleCreateComment(store)))
	mux.HandleFunc("GET /posts/{id}/comments", noAuth(api.HandleGetComments(store)))
	mux.HandleFunc("PUT /comments/{id}", authenticated(api.HandleUpdateComments(store)))
	mux.HandleFunc("POST /comments/{id}/likes", authenticated(api.HandleLikeComment(store)))
	mux.HandleFunc("DELETE /comments/{id}", authenticated(api.HandleDeleteComments(store)))
	mux.HandleFunc(
		"DELETE /comments/{id}/likes",
		authenticated(api.HandleRemoveLikeFromComment(store)),
	)

	muxer := http.NewServeMux()
	muxer.Handle("/api/v1/", http.StripPrefix("/api/v1", mux))

	return muxer
}

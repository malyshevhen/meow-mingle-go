package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/malyshEvhen/meow_mingle/internal/config"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/internal/handlers"
	"github.com/malyshEvhen/meow_mingle/internal/middleware"
)

func RegisterRoutes(store db.IStore, cfg config.Config) *mux.Router {
	mux := mux.NewRouter()
	apiMux := mux.PathPrefix("/api/v1").Subrouter()

	authenticated := func(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
		return middleware.MiddlewareChain(
			handler,
			middleware.LoggerMW,
			middleware.ErrorHandler,
			middleware.WithJWTAuth(store, cfg),
		)
	}

	public := func(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
		return middleware.MiddlewareChain(
			handler,
			middleware.LoggerMW,
			middleware.ErrorHandler,
		)
	}

	usersMux := apiMux.PathPrefix("/users").Subrouter()
	postsMux := apiMux.PathPrefix("/posts").Subrouter()
	commentsMux := apiMux.PathPrefix("/comments").Subrouter()

	usersMux.HandleFunc("/register", public(handlers.HandleCreateUser(store, cfg))).Methods("POST")
	usersMux.HandleFunc("/{id}/feed", public(handlers.HandleUsersFeed(store))).Methods("GET")
	usersMux.HandleFunc("/{id}/posts", public(handlers.HandleGetUserPosts(store))).Methods("GET")
	usersMux.HandleFunc("/{id}", authenticated(handlers.HandleGetUser(store))).Methods("GET")
	usersMux.HandleFunc("/feed", authenticated(handlers.HandleOwnersFeed(store))).Methods("GET")
	usersMux.HandleFunc("/{id}/subscriptions", authenticated(handlers.HandleSubscribe(store))).Methods("POST")
	usersMux.HandleFunc("/{id}/subscriptions", authenticated(handlers.HandleUnsubscribe(store))).Methods("DELETE")

	postsMux.HandleFunc("/{id}", public(handlers.HandleGetPostsById(store))).Methods("GET")
	postsMux.HandleFunc("/{id}/comments", public(handlers.HandleGetComments(store))).Methods("GET")
	postsMux.HandleFunc("", authenticated(handlers.HandleCreatePost(store))).Methods("POST")
	postsMux.HandleFunc("/{id}", authenticated(handlers.HandleUpdatePostsById(store))).Methods("PUT")
	postsMux.HandleFunc("/{id}", authenticated(handlers.HandleDeletePostsById(store))).Methods("DELETE")
	postsMux.HandleFunc("/{id}/likes", authenticated(handlers.HandleLikePost(store))).Methods("POST")
	postsMux.HandleFunc("/{id}/likes", authenticated(handlers.HandleRemoveLikeFromPost(store))).Methods("DELETE")
	postsMux.HandleFunc("/{id}/comments", authenticated(handlers.HandleCreateComment(store))).Methods("POST")

	commentsMux.HandleFunc("/{id}", authenticated(handlers.HandleUpdateComments(store))).Methods("PUT")
	commentsMux.HandleFunc("/{id}", authenticated(handlers.HandleDeleteComments(store))).Methods("DELETE")
	commentsMux.HandleFunc("/{id}/likes", authenticated(handlers.HandleLikeComment(store))).Methods("POST")
	commentsMux.HandleFunc("/{id}/likes", authenticated(handlers.HandleRemoveLikeFromComment(store))).Methods("DELETE")

	return mux
}

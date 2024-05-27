package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/malyshEvhen/meow_mingle/internal/api"
	"github.com/malyshEvhen/meow_mingle/internal/config"
	"github.com/malyshEvhen/meow_mingle/internal/db"
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

	puplic := func(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
		return middleware.MiddlewareChain(
			handler,
			middleware.LoggerMW,
			middleware.ErrorHandler,
		)
	}

	usersMux := apiMux.PathPrefix("/users").Subrouter()
	postsMux := apiMux.PathPrefix("/posts").Subrouter()
	commentsMux := apiMux.PathPrefix("/comments").Subrouter()

	usersMux.HandleFunc("/register", puplic(api.HandleCreateUser(store, cfg))).Methods("POST")
	usersMux.HandleFunc("/{id}/feed", puplic(api.HandleUsersFeed(store))).Methods("GET")
	usersMux.HandleFunc("/{id}/posts", puplic(api.HandleGetUserPosts(store))).Methods("GET")
	usersMux.HandleFunc("/{id}", authenticated(api.HandleGetUser(store))).Methods("GET")
	usersMux.HandleFunc("/feed", authenticated(api.HandleOwnersFeed(store))).Methods("GET")
	usersMux.HandleFunc("/{id}/subscriptions", authenticated(api.HandleSubscribe(store))).Methods("POST")
	usersMux.HandleFunc("/{id}/subscriptions", authenticated(api.HandleUnsubscribe(store))).Methods("POST")

	postsMux.HandleFunc("/{id}", puplic(api.HandleGetPostsById(store))).Methods("GET")
	postsMux.HandleFunc("/{id}/comments", puplic(api.HandleGetComments(store))).Methods("GET")
	postsMux.HandleFunc("/", authenticated(api.HandleCreatePost(store))).Methods("POST")
	postsMux.HandleFunc("/{id}", authenticated(api.HandleUpdatePostsById(store))).Methods("PUT")
	postsMux.HandleFunc("/{id}", authenticated(api.HandleDeletePostsById(store))).Methods("DELETE")
	postsMux.HandleFunc("/{id}/likes", authenticated(api.HandleLikePost(store))).Methods("POST")
	postsMux.HandleFunc("/{id}/likes", authenticated(api.HandleRemoveLikeFromPost(store))).Methods("DELETE")
	postsMux.HandleFunc("/{id}/comments", authenticated(api.HandleCreateComment(store))).Methods("POST")

	commentsMux.HandleFunc("/{id}", authenticated(api.HandleUpdateComments(store))).Methods("PUT")
	commentsMux.HandleFunc("/{id}", authenticated(api.HandleDeleteComments(store))).Methods("DELETE")
	commentsMux.HandleFunc("/{id}/likes", authenticated(api.HandleLikeComment(store))).Methods("POST")
	commentsMux.HandleFunc("/{id}/likes", authenticated(api.HandleRemoveLikeFromComment(store))).Methods("DELETE")

	return mux
}

package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/pkg/api"
	"github.com/malyshEvhen/meow_mingle/pkg/auth"
)

func RegisterRouts(
	authMW *auth.Provider,
	profileRepo db.IProfileRepository,
	commentRepo db.ICommentRepository,
	postRepo db.IPostRepository,
	secret string,
) *mux.Router {
	auth := func(handler api.Handler) http.HandlerFunc {
		return authenticated(handler, authMW.WithJWTAuth)
	}

	r := mux.NewRouter().PathPrefix("/api/v1").Subrouter()

	// Feed API
	r.HandleFunc("/feed", auth(handleGetFeed(postRepo))).Methods("GET")

	// Post API
	r.HandleFunc("/posts", auth(handleCreatePost(postRepo))).Methods("POST")
	r.HandleFunc("/posts", auth(handleGetPosts(postRepo))).Methods("GET")
	r.HandleFunc("/posts/{id}", auth(handleGetPostById(postRepo))).Methods("GET")
	r.HandleFunc("/posts/{id}", auth(handleUpdatePostById(postRepo))).Methods("PATCH")
	r.HandleFunc("/posts/{id}", auth(handleDeletePostById(postRepo))).Methods("DELETE")

	// Comment API
	r.HandleFunc("/comments", auth(handleCreateComment(commentRepo))).Methods("POST")
	r.HandleFunc("/comments", auth(handleGetComments(commentRepo))).Methods("GET")
	r.HandleFunc("/comments/{id}", auth(handleUpdateComment(commentRepo))).Methods("PUT")
	r.HandleFunc("/comments/{id}", auth(handleDeleteComment(commentRepo))).Methods("DELETE")

	// Profile API
	r.HandleFunc("/profiles", public(handleCreateProfile(profileRepo))).Methods("POST")
	r.HandleFunc("/profiles/{id}", auth(handleGetProfile(profileRepo))).Methods("GET")
	r.HandleFunc("/profiles/{id}/subscriptions", auth(handleSubscribe(profileRepo))).Methods("POST")
	r.HandleFunc("/profiles/{id}/subscriptions", auth(handleUnsubscribe(profileRepo))).Methods("DELETE")

	// Reaction API
	r.HandleFunc("/reactions", auth(nil)).Methods("PUT")         // TODO: implement
	r.HandleFunc("/reactions", auth(nil)).Methods("GET")         // TODO: implement
	r.HandleFunc("/reactions/{id}", auth(nil)).Methods("DELETE") // TODO: implement

	return r
}

func authenticated(handler api.Handler, authMW api.Middleware) http.HandlerFunc {
	return middlewareChain(handler, loggerMW, ErrorHandler, authMW)
}

func public(handler api.Handler) http.HandlerFunc {
	return middlewareChain(handler, loggerMW, ErrorHandler)
}

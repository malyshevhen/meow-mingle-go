package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/pkg/api"
	"github.com/malyshEvhen/meow_mingle/pkg/auth"
)

func RegisterRouts(
	authMW *auth.Provider,
	profileService app.ProfileService,
	commentService app.CommentService,
	postService app.PostService,
	subscriptionService app.SubscriptionService,
) *mux.Router {
	auth := func(handler api.Handler) http.HandlerFunc {
		return authenticated(handler, authMW.WithJWTAuth)
	}

	r := mux.NewRouter().PathPrefix("/api/v1").Subrouter()

	// Feed API
	r.HandleFunc("/feed", auth(handleGetFeed(postService))).Methods("GET")

	// Post API
	r.HandleFunc("/posts", auth(handleCreatePost(postService))).Methods("POST")
	r.HandleFunc("/posts", auth(handleGetPosts(postService))).Methods("GET")
	r.HandleFunc("/posts/{id}", auth(handleGetPostById(postService))).Methods("GET")
	r.HandleFunc("/posts/{id}", auth(handleUpdatePostById(postService))).Methods("PATCH")
	r.HandleFunc("/posts/{id}", auth(handleDeletePostById(postService))).Methods("DELETE")

	// Comment API
	r.HandleFunc("/comments", auth(handleCreateComment(commentService))).Methods("POST")
	r.HandleFunc("/comments", auth(handleGetComments(commentService))).Methods("GET")
	r.HandleFunc("/comments/{id}", auth(handleUpdateComment(commentService))).Methods("PUT")
	r.HandleFunc("/comments/{id}", auth(handleDeleteComment(commentService))).Methods("DELETE")

	// Profile API
	r.HandleFunc("/profiles", public(handleCreateProfile(profileService))).Methods("POST")
	r.HandleFunc("/profiles/{id}", auth(handleGetProfile(profileService))).Methods("GET")

	// Subscription API
	r.HandleFunc("/subscriptions{id}", auth(handleSubscribe(subscriptionService))).Methods("POST")
	r.HandleFunc("/subscriptions{id}", auth(handleUnsubscribe(subscriptionService))).Methods("DELETE")

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

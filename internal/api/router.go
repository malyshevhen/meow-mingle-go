package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/internal/auth"
	"github.com/malyshEvhen/meow_mingle/pkg/api"
)

func RegisterRouts(
	authMW *auth.Provider,
	profileService app.ProfileService,
	commentService app.CommentService,
	postService app.PostService,
	subscriptionService app.SubscriptionService,
	reactionService app.ReactionService,
) *mux.Router {
	auth := func(handler api.Handler) http.Handler {
		return authenticated(handler, authMW.Basic)
	}

	r := mux.NewRouter().PathPrefix("/api/v1").Subrouter()

	// Feed API
	r.Handle("/feed", auth(handleGetFeed(postService))).Methods("GET")

	// Post API
	r.Handle("/posts", auth(handleCreatePost(postService))).Methods("POST")
	r.Handle("/posts", auth(handleGetPosts(postService))).Methods("GET")
	r.Handle("/posts/{id}", auth(handleGetPostById(postService))).Methods("GET")
	r.Handle("/posts/{id}", auth(handleUpdatePostById(postService))).Methods("PATCH")
	r.Handle("/posts/{id}", auth(handleDeletePostById(postService))).Methods("DELETE")

	// Comment API
	r.Handle("/comments", auth(handleCreateComment(commentService))).Methods("POST")
	r.Handle("/comments", auth(handleGetComments(commentService))).Methods("GET")
	r.Handle("/comments/{id}", auth(handleUpdateComment(commentService))).Methods("PUT")
	r.Handle("/comments/{id}", auth(handleDeleteComment(commentService))).Methods("DELETE")

	// Profile API
	r.Handle("/profiles", public(handleCreateProfile(profileService))).Methods("POST")
	r.Handle("/profiles/{id}", auth(handleGetProfile(profileService))).Methods("GET")

	// Subscription API
	r.Handle("/subscriptions{id}", auth(handleSubscribe(subscriptionService))).Methods("POST")
	r.Handle("/subscriptions{id}", auth(handleUnsubscribe(subscriptionService))).Methods("DELETE")

	// Reaction API
	r.Handle("/reactions", auth(handleCreateReaction(reactionService))).Methods("PUT")
	r.Handle("/reactions/{id}", auth(handleDeleteReaction(reactionService))).Methods("DELETE")

	return r
}

func authenticated(handler api.Handler, authMW api.Middleware) http.Handler {
	return middlewareChain(handler, loggerMW, ErrorHandler, authMW)
}

func public(handler api.Handler) http.Handler {
	return middlewareChain(handler, loggerMW, ErrorHandler)
}

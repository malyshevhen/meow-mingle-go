package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/malyshEvhen/meow_mingle/internal/db"
)

func RegisterRouts(
	ctx context.Context,
	authMW *AuthProvider,
	userRepo db.IUserRepository,
	commentRepo db.ICommentRepository,
	postRepo db.IPostRepository,
	secret string,
) *mux.Router {
	auth := func(handler Handler) http.HandlerFunc {
		return authenticated(handler, authMW.WithJWTAuth)
	}

	mux := mux.NewRouter().PathPrefix("/api/v1").Subrouter()

	mux.HandleFunc("/login", auth(handleLogin(secret))).Methods("POST")
	// TODO: logout
	mux.HandleFunc("/register", public(handleRegistration(userRepo, secret))).Methods("POST")

	mux.HandleFunc("/feed", auth(handleGetFeed(postRepo))).Methods("GET")
	mux.HandleFunc("/posts/{id}", auth(handleGetPostsById(postRepo))).Methods("GET")
	mux.HandleFunc("/posts/{id}/comments", auth(handleGetComments(commentRepo))).Methods("GET")
	mux.HandleFunc("/posts", auth(handleCreatePost(postRepo))).Methods("POST")
	mux.HandleFunc("/posts/{id}", auth(handleUpdatePostsById(postRepo))).Methods("PUT")
	mux.HandleFunc("/posts/{id}", auth(handleDeletePostsById(postRepo))).Methods("DELETE")
	mux.HandleFunc("/posts/{id}/likes", auth(handleLikePost(postRepo))).Methods("POST")
	mux.HandleFunc("/posts/{id}/likes", auth(handleRemoveLikeFromPost(postRepo))).Methods("DELETE")
	mux.HandleFunc("/posts/{id}/comments", auth(handleCreateComment(commentRepo))).Methods("POST")

	mux.HandleFunc("/comments/{id}", auth(handleUpdateComments(commentRepo))).Methods("PUT")
	mux.HandleFunc("/comments/{id}", auth(handleDeleteComments(commentRepo))).Methods("DELETE")
	mux.HandleFunc("/comments/{id}/likes", auth(handleLikeComment(commentRepo))).Methods("POST")
	mux.HandleFunc("/comments/{id}/likes", auth(handleRemoveLikeFromComment(commentRepo))).Methods("DELETE")

	mux.HandleFunc("/users/{id}/feed", auth(handleUsersFeed(postRepo))).Methods("GET")
	mux.HandleFunc("/users/{id}/posts", auth(handleGetUserPosts(postRepo))).Methods("GET")
	mux.HandleFunc("/users/{id}", auth(handleGetUser(userRepo))).Methods("GET")
	mux.HandleFunc("/users/{id}/subscriptions", auth(handleSubscribe(userRepo))).Methods("POST")
	mux.HandleFunc("/users/{id}/subscriptions", auth(handleUnsubscribe(userRepo))).Methods("DELETE")

	return mux
}

func authenticated(handler Handler, authMW Middleware) http.HandlerFunc {
	return MiddlewareChain(handler, LoggerMW, ErrorHandler, authMW)
}

func public(handler Handler) http.HandlerFunc {
	return MiddlewareChain(handler, LoggerMW, ErrorHandler)
}

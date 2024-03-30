package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	addr           string
	userService    *UserService
	postService    *PostService
	commentService *CommentService
}

func NewApiServer(addr string, userService *UserService, postService *PostService,
	commentService *CommentService) *ApiServer {
	return &ApiServer{
		addr:           addr,
		userService:    userService,
		postService:    postService,
		commentService: commentService,
	}
}

func (s *ApiServer) Serve() {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userService := NewUserController(s.userService)
	userService.RegisterRoutes(subrouter)

	postsService := NewPostController(s.userService, s.postService)
	postsService.RegisterRoutes(subrouter)

	commentController := NewCommentController(s.userService, s.commentService)
	commentController.RegisterRoutes(subrouter)

	log.Printf("Server starting at port: %s\n", s.addr)

	log.Fatal(http.ListenAndServe(s.addr, subrouter))
}

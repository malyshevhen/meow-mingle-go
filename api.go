package main

import (
	"log"
	"net/http"
)

type apiHandler func(w http.ResponseWriter, r *http.Request) error

type ApiServer struct {
	addr           string
	userService    *UserService
	postService    *PostService
	commentService *CommentService
	sCtx           *SecurityContextHolder
}

func NewApiServer(addr string, us *UserService, ps *PostService,
	cs *CommentService, sCtx *SecurityContextHolder) *ApiServer {
	return &ApiServer{
		addr:           addr,
		userService:    us,
		postService:    ps,
		commentService: cs,
		sCtx:           sCtx,
	}
}

func (s *ApiServer) Serve() {
	subrouter := http.NewServeMux()

	userController := NewUserController(s.sCtx, s.userService)
	userController.RegisterRoutes(subrouter)

	postsController := NewPostController(s.postService, s.sCtx)
	postsController.RegisterRoutes(subrouter)

	commentController := NewCommentController(s.commentService, s.sCtx)
	commentController.RegisterRoutes(subrouter)

	router := http.NewServeMux()
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", subrouter))

	log.Printf("Server starting at port: %s\n", s.addr)

	log.Fatal(http.ListenAndServe(s.addr, router))
}

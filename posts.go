package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
)

type PostController struct {
	postService *PostService
	sCtx        *SecurityContextHolder
}

func NewPostController(postService *PostService, sCtx *SecurityContextHolder) *PostController {
	return &PostController{
		postService: postService,
		sCtx:        sCtx,
	}
}

func (ps *PostController) RegisterRoutes(r *http.ServeMux) {
	middlewareChain := func(handler apiHandler) http.HandlerFunc {
		return MiddlewareChain(
			handler,
			LoggerMiddleware,
			ErrorHandler,
			ps.sCtx.WithJWTAuth,
		)
	}
	r.HandleFunc("GET /users/{id}/posts", middlewareChain(ps.handleGetUserPosts))
	r.HandleFunc("POST /posts", middlewareChain(ps.handleCreatePost))
	r.HandleFunc("GET /posts/{id}", middlewareChain(ps.handleGetPostsById))
	r.HandleFunc("PUT /posts/{id}", middlewareChain(ps.handleUpdatePostsById))
	r.HandleFunc("DELETE /posts/{id}", middlewareChain(ps.handleDeletePostsById))
}

func (ps *PostController) handleCreatePost(w http.ResponseWriter, r *http.Request) error {
	postRequest, err := readPostReqType(r)
	if err != nil {
		log.Printf("%-15s ==> 😞 Error reading post request: %v\n", "PostController", err)
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Error reading post request",
		}
	}

	userId, err := ps.sCtx.GetAuthUserId(r)
	if err != nil {
		log.Printf("%-15s ==> 😱 Error getting user Id from token %v\n", "PostController ", err)
		return err
	}

	postResponse, err := ps.postService.CreatePost(userId, postRequest)
	if err != nil {
		log.Printf("%-15s ==> 🤯 Error creating post in store %v\n", "PostController", err)
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error creating post"))
		return err
	}

	log.Printf("%-15s ==> 🎉 Successfully created new post\n", "PostController")

	return WriteJson(w, http.StatusCreated, postResponse)
}

func (ps *PostController) handleGetUserPosts(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id param %v\n", "PostController", err)
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Error reading post request",
		}
	}

	postResponse, err := ps.postService.GetUserPosts(id)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error getting user posts from store %v\n", "PostController", err)
		return err
	}

	log.Printf("%-15s ==> 🤩 Successfully retrieved user posts\n", "PostController")

	return WriteJson(w, http.StatusOK, postResponse)
}

func (ps *PostController) handleGetPostsById(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id para:%v\n", "PostController", err)
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Error reading post request",
		}
	}

	postResponse, err := ps.postService.GetPostById(id)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error getting post by Id from stor:%v\n", "PostController", err)
		return err
	}

	log.Printf("%-15s ==> 🤩 Successfully retrieved post by Id\n", "PostController")

	return WriteJson(w, http.StatusOK, postResponse)
}

func (ps *PostController) handleUpdatePostsById(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id para %v\n", "PostController", err)
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Invalid param",
		}
	}

	postRequest, err := readPostReqType(r)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error reading post request %v\n", "PostController", err)
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Invalid content",
		}
	}

	postResponse, err := ps.postService.UpdatePostById(id, postRequest)
	if err != nil {
		log.Printf("%-15s ==> 🤯 Error updating post by Id in store %v\n", "PostController", err)
		return err
	}

	log.Printf("%-15s ==> 🎉 Successfully updated post by Id\n", "PostController")

	return WriteJson(w, http.StatusOK, postResponse)
}

func (ps *PostController) handleDeletePostsById(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id param %v\n", "PostController", err)
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Invalid param",
		}
	}

	if err := ps.postService.DeletePostById(id); err != nil {
		log.Printf("%-15s ==> 😫 Error deleting post by Id from store %v\n", "PostController", err)
		return err
	}

	log.Printf("%-15s ==> 🗑️ Successfully deleted post by Id\n", "PostController")

	return WriteJson(w, http.StatusNoContent, nil)
}

func readPostReqType(r *http.Request) (*PostRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	p, err := Unmarshal[PostRequest](body)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func parseIdParam(r *http.Request) (int64, error) {
	id := r.PathValue("id")

	numId, err := strconv.Atoi(id)
	if err != nil {
		return 0, nil
	}

	return int64(numId), nil
}

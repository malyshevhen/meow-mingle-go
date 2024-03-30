package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type PostController struct {
	userService *UserService
	postService *PostService
}

func NewPostController(userService *UserService, postService *PostService) *PostController {
	return &PostController{
		userService: userService,
		postService: postService,
	}
}

func (ps *PostController) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users/posts", WithJWTAuth(ps.handleCreatePost, ps.userService)).Methods("POST")
	r.HandleFunc("/users/{id}/posts", WithJWTAuth(ps.handleGetUserPosts, ps.userService)).Methods("GET")
	r.HandleFunc("/users/posts/{id}", WithJWTAuth(ps.handleGetPostsById, ps.userService)).Methods("GET")
	r.HandleFunc("/users/posts/{id}", WithJWTAuth(ps.handleUpdatePostsById, ps.userService)).Methods("PUT")
	r.HandleFunc("/users/posts/{id}", WithJWTAuth(ps.handleDeletePostsById, ps.userService)).Methods("DELETE")
}

func (ps *PostController) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	p, err := readPostReqType(r)
	if err != nil {
		log.Printf("%-15s ==> 😞 Error reading post request: %v\n", "PostService", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Error reading post request"))
		return
	}

	token := GetTokenFromRequest(r)
	userId, err := GetAuthUserId(token)
	if err != nil {
		log.Printf("%-15s ==> 😱 Error getting user Id from token %v\n", "PostService ", err)
		WriteJson(w, http.StatusUnauthorized, NewErrorResponse("Error getting user Id from token"))
		return
	}

	pResp, err := ps.postService.CreatePost(userId, p)
	if err != nil {
		log.Printf("%-15s ==> 🤯 Error creating post in store %v\n", "PostService ", err)
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error creating post"))
		return
	}

	log.Printf("%-15s ==> 🎉 Successfully created new post\n", "PostService!")
	WriteJson(w, http.StatusCreated, pResp)
}

func (ps *PostController) handleGetUserPosts(w http.ResponseWriter, r *http.Request) {
	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id param %v\n", "PostService ", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Error parsing Id param"))
		return
	}

	p, err := ps.postService.GetUserPosts(id)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error getting user posts from store %v\n", "PostService ", err)
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error getting user posts from store"))
		return
	}

	log.Printf("%-15s ==> 🤩 Successfully retrieved user posts\n", "PostService")
	WriteJson(w, http.StatusOK, p)
}

func (ps *PostController) handleGetPostsById(w http.ResponseWriter, r *http.Request) {
	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id para:%v\n", "PostService ", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Error parsing Id param"))
		return
	}

	p, err := ps.postService.GetPostById(id)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error getting post by Id from stor:%v\n", "PostService ", err)
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error getting post by Id from store"))
		return
	}

	log.Printf("%-15s ==> 🤩 Successfully retrieved post by Id\n", "PostService")
	WriteJson(w, http.StatusOK, p)
}

func (ps *PostController) handleUpdatePostsById(w http.ResponseWriter, r *http.Request) {
	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id para %v\n", "PostService: ", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Error parsing Id param"))
		return
	}

	p, err := readPostReqType(r)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error reading post request %v\n", "PostService: ", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Error reading post request"))
		return
	}

	pr, err := ps.postService.UpdatePostById(id, p)
	if err != nil {
		log.Printf("%-15s ==> 🤯 Error updating post by Id in store %v\n", "PostService: ", err)
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error updating post by Id in store"))
		return
	}

	log.Printf("%-15s ==> 🎉 Successfully updated post by Id\n", "PostService ")
	WriteJson(w, http.StatusOK, pr)
}

func (ps *PostController) handleDeletePostsById(w http.ResponseWriter, r *http.Request) {
	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id param %v\n", "PostService ", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Error parsing Id param"))
		return
	}

	err = ps.postService.DeletePostById(id)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error deleting post by Id from store %v\n", "PostService ", err)
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error deleting post by Id from store"))
		return
	}

	log.Printf("%-15s ==> 🗑️ Successfully deleted post by Id\n", "PostService")
	WriteJson(w, http.StatusNoContent, nil)
}

func readPostReqType(r *http.Request) (*PostRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var p *PostRequest
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, err
	}
	return p, nil
}

func parseIdParam(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	id := vars["id"]

	numId, err := strconv.Atoi(id)
	if err != nil {
		return 0, nil
	}

	return int64(numId), nil
}

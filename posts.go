package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type PostService struct {
	store Store
}

func NewPostService(s Store) *PostService {
	return &PostService{store: s}
}

func (ps *PostService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users/posts", WithJWTAuth(ps.handleCreatePost, ps.store)).Methods("POST")
	r.HandleFunc("/users/{id}/posts", WithJWTAuth(ps.handleGetUserPosts, ps.store)).Methods("GET")
	r.HandleFunc("/users/posts/{id}", WithJWTAuth(ps.handleGetPostsById, ps.store)).Methods("GET")
	r.HandleFunc("/users/posts/{id}", WithJWTAuth(ps.handleUpdatePostsById, ps.store)).Methods("PUT")
	r.HandleFunc("/users/posts/{id}", WithJWTAuth(ps.handleDeletePostsById, ps.store)).Methods("DELETE")
	r.HandleFunc("/users/posts/{id}/comments", WithJWTAuth(ps.handleCreateComment, ps.store)).Methods("POST")
	r.HandleFunc("/users/posts/{id}/comments", WithJWTAuth(ps.handleGetComments, ps.store)).Methods("GET")
	r.HandleFunc("/users/posts/comments/{id}", WithJWTAuth(ps.handleUpdateComments, ps.store)).Methods("PUT")
	r.HandleFunc("/users/posts/comments/{id}", WithJWTAuth(ps.handleDeleteComments, ps.store)).Methods("DELETE")
}

func (ps *PostService) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	p, err := readPostReqType(r)
	if err != nil {
		log.Printf("%-15s ==> 😞 Error reading post request: %v\n", "PostService", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Error reading post request"))
		return
	}

	token := GetTokenFromRequest(r)
	userId, err := GetAuthUserId(token)
	if err != nil {
		log.Printf("%-15s ==> 😱 Error getting user ID from token %v\n", "PostService ", err)
		WriteJson(w, http.StatusUnauthorized, NewErrorResponse("Error getting user ID from token"))
		return
	}

	pResp, err := ps.store.CreatePost(userId, p)
	if err != nil {
		log.Printf("%-15s ==> 🤯 Error creating post in store %v\n", "PostService ", err)
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error creating post"))
		return
	}

	log.Printf("%-15s ==> 🎉 Successfully created new post\n", "PostService!")
	WriteJson(w, http.StatusCreated, pResp)
}

func (ps *PostService) handleGetUserPosts(w http.ResponseWriter, r *http.Request) {
	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing ID param %v\n", "PostService ", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Error parsing ID param"))
		return
	}

	p, err := ps.store.GetUserPosts(id)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error getting user posts from store %v\n", "PostService ", err)
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error getting user posts from store"))
		return
	}

	log.Printf("%-15s ==> 🤩 Successfully retrieved user posts\n", "PostService")
	WriteJson(w, http.StatusOK, p)
}

func (ps *PostService) handleGetPostsById(w http.ResponseWriter, r *http.Request) {
	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing ID para:%v\n", "PostService ", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Error parsing ID param"))
		return
	}

	p, err := ps.store.GetPostById(id)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error getting post by ID from stor:%v\n", "PostService ", err)
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error getting post by ID from store"))
		return
	}

	log.Printf("%-15s ==> 🤩 Successfully retrieved post by ID\n", "PostService")
	WriteJson(w, http.StatusOK, p)
}

func (ps *PostService) handleUpdatePostsById(w http.ResponseWriter, r *http.Request) {
	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing ID para %v\n", "PostService: ", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Error parsing ID param"))
		return
	}

	p, err := readPostReqType(r)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error reading post request %v\n", "PostService: ", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Error reading post request"))
		return
	}

	pr, err := ps.store.UpdatePostById(id, p)
	if err != nil {
		log.Printf("%-15s ==> 🤯 Error updating post by ID in store %v\n", "PostService: ", err)
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error updating post by ID in store"))
		return
	}

	log.Printf("%-15s ==> 🎉 Successfully updated post by ID\n", "PostService ")
	WriteJson(w, http.StatusOK, pr)
}

func (ps *PostService) handleDeletePostsById(w http.ResponseWriter, r *http.Request) {
	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing ID param %v\n", "PostService ", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Error parsing ID param"))
		return
	}

	err = ps.store.DeletePostById(id)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error deleting post by ID from store %v\n", "PostService ", err)
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error deleting post by ID from store"))
		return
	}

	log.Printf("%-15s ==> 🗑️ Successfully deleted post by ID\n", "PostService")
	WriteJson(w, http.StatusNoContent, nil)
}

func (ps *PostService) handleCreateComment(w http.ResponseWriter, r *http.Request) {
	postId, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing post ID param %v\n", "PostService ", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Error parsing post ID param"))
		return
	}

	cReq, err := readCommentReqType(r)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error reading comment request %v\n", "PostService ", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Error reading comment request"))
		return
	}

	token := GetTokenFromRequest(r)
	userId, err := GetAuthUserId(token)
	if err != nil {
		log.Printf("%-15s ==> 😱 Error getting authenticated user ID %v\n", "PostService ", err)
		WriteJson(w, http.StatusUnauthorized, NewErrorResponse("Error getting authenticated user ID"))
		return
	}

	cResp, err := ps.store.CreateComment(postId, userId, cReq)
	if err != nil {
		log.Printf("%-15s ==> 🤯 Error creating comment in store %v\n", "PostService ", err)
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error creating comment in store"))
		return
	}

	log.Printf("%-15s ==> 🎉 Successfully created comment\n", "PostService")
	WriteJson(w, http.StatusCreated, cResp)
}

func (ps *PostService) handleGetComments(w http.ResponseWriter, r *http.Request) {
	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing ID para %v\n", "PostService ", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Error parsing ID param"))
		return
	}

	c, err := ps.store.GetCommentsByPostId(id)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error getting comment by ID from stor %v\n", "PostService ", err)
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error getting comment by ID from store"))
		return
	}

	log.Printf("%-15s ==> 🎉 Successfully got comment by ID\n", "PostService!")
	WriteJson(w, http.StatusOK, c)
}

func (ps *PostService) handleUpdateComments(w http.ResponseWriter, r *http.Request) {
	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing ID para %v\n", "PostService ", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Error parsing ID param"))
		return
	}

	c, err := readCommentReqType(r)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error reading comment request %v\n", "PostService ", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Error reading comment request"))
		return
	}

	cr, err := ps.store.UpdateCommentById(id, c)
	if err != nil {
		log.Printf("%-15s ==> 😱 Error updating comment by ID in stor %v\n", "PostService ", err)
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error updating comment by ID in store"))
		return
	}

	log.Printf("%-15s ==> 🎉 Successfully updated comment by ID\n", "PostService")
	WriteJson(w, http.StatusOK, cr)
}

func (ps *PostService) handleDeleteComments(w http.ResponseWriter, r *http.Request) {
	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing ID para\n ", "PostService")
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Error parsing ID param"))
		return
	}

	err = ps.store.DeleteCommentById(id)
	if err != nil {
		log.Printf("%-15s ==> 😱 Error deleting comment by ID from stor\n ", "PostService")
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error deleting comment by ID from store"))
		return
	}

	log.Printf("%-15s ==> 🎉 Successfully deleted comment by ID\n", "PostService")
	WriteJson(w, http.StatusNoContent, nil)
}

func readCommentReqType(r *http.Request) (*CommentRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var c *CommentRequest
	if err := json.Unmarshal(body, &c); err != nil {
		return nil, err
	}
	return c, nil
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

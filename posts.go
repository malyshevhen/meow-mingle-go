package main

import (
	"encoding/json"
	"io"
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
}

func (ps *PostService) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	p, err := readPostReqType(r)
	if err != nil {
		return
	}

	// TODO: validate payload

	token := GetTokenFromRequest(r)
	userId, err := GetAuthUserId(token)

	pResp, err := ps.store.CreatePost(userId, p)
	if err != nil {
		return
	}

	WriteJson(w, http.StatusCreated, pResp)
}

func (ps *PostService) handleGetUserPosts(w http.ResponseWriter, r *http.Request) {
	id, err := parseIdParam(r)
	if err != nil {
		return
	}

	p, err := ps.store.GetUserPosts(id)
	if err != nil {
		return
	}

	WriteJson(w, http.StatusOK, p)
}

func (ps *PostService) handleGetPostsById(w http.ResponseWriter, r *http.Request) {
	id, err := parseIdParam(r)
	if err != nil {
		return
	}

	p, err := ps.store.GetPostById(id)
	if err != nil {
		return
	}

	WriteJson(w, http.StatusOK, p)
}

func (ps *PostService) handleUpdatePostsById(w http.ResponseWriter, r *http.Request) {
	id, err := parseIdParam(r)
	if err != nil {
		return
	}

	p, err := readPostReqType(r)
	if err != nil {
		return
	}

	pr, err := ps.store.UpdatePostById(id, p)
	if err != nil {
		return
	}

	WriteJson(w, http.StatusOK, pr)
}

func (ps *PostService) handleDeletePostsById(w http.ResponseWriter, r *http.Request) {
	id, err := parseIdParam(r)
	if err != nil {
		return
	}

	err = ps.store.DeletePostById(id)
	if err != nil {
		return
	}

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

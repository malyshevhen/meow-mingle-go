package main

import (
	"encoding/json"
	"io"
	"net/http"

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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	defer r.Body.Close()

	var p *PostRequest
	if err := json.Unmarshal(body, &p); err != nil {
		return
	}

	// TODO: validation

	pResp, err := ps.store.CreatePost(p)
	if err != nil {
		return
	}

	WriteJson(w, http.StatusCreated, pResp)
}

func (ps *PostService) handleGetUserPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	p, err := ps.store.GetUserPosts(id)
	if err != nil {
		return
	}

	WriteJson(w, http.StatusOK, p)
}

func (ps *PostService) handleGetPostsById(w http.ResponseWriter, r *http.Request)    {
	vars := mux.Vars(r)
	id := vars["id"]

	p, err := ps.store.GetPostsById(id)
	if err != nil {
		return
	}

	WriteJson(w, http.StatusOK, p)
}

func (ps *PostService) handleUpdatePostsById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	defer r.Body.Close()

	var p *PostRequest
	if err := json.Unmarshal(body, &p); err != nil {
		return
	}

	pr, err := ps.store.UpdatePostsById(id, p)
	if err != nil {
		return
	}

	WriteJson(w, http.StatusOK, pr)
}

func (ps *PostService) handleDeletePostsById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := ps.store.DeletePostsById(id)
	if err != nil {
		return
	}

	WriteJson(w, http.StatusNoContent, nil)
}

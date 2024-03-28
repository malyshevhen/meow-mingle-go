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
	r.HandleFunc("/users/posts/{id}/comments", WithJWTAuth(ps.handleCreateComment, ps.store)).Methods("POST")
	r.HandleFunc("/users/posts/{id}/comments", WithJWTAuth(ps.handleGetComments, ps.store)).Methods("GET")
	r.HandleFunc("/users/posts/comments/{id}", WithJWTAuth(ps.handleUpdateComments, ps.store)).Methods("PUT")
	r.HandleFunc("/users/posts/comments/{id}", WithJWTAuth(ps.handleDeleteComments, ps.store)).Methods("DELETE")
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

func (ps *PostService) handleCreateComment(w http.ResponseWriter, r *http.Request) {
	postId, err := parseIdParam(r)
	if err != nil {
		return
	}

	cReq, err := readCommentReqType(r)
	if err != nil {
		return
	}

	token := GetTokenFromRequest(r)
	userId, err := GetAuthUserId(token)

	cResp, err := ps.store.CreateComment(postId, userId, cReq)
	if err != nil {
		return
	}

	WriteJson(w, http.StatusCreated, cResp)
}

func (ps *PostService) handleGetComments(w http.ResponseWriter, r *http.Request) {
	id, err := parseIdParam(r)
	if err != nil {
		return
	}

	c, err := ps.store.GetCommentById(id)
	if err != nil {
		return
	}

	WriteJson(w, http.StatusOK, c)

}

func (ps *PostService) handleUpdateComments(w http.ResponseWriter, r *http.Request) {
	id, err := parseIdParam(r)
	if err != nil {
		return
	}

	c, err := readCommentReqType(r)
	if err != nil {
		return
	}

	cr, err := ps.store.UpdateCommentById(id, c)
	if err != nil {
		return
	}

	WriteJson(w, http.StatusOK, cr)
}

func (ps *PostService) handleDeleteComments(w http.ResponseWriter, r *http.Request) {
	id, err := parseIdParam(r)
	if err != nil {
		return
	}

	err = ps.store.DeleteCommentById(id)
	if err != nil {
		return
	}

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

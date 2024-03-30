package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type CommentController struct {
	userService    *UserService
	commentService *CommentService
}

func NewCommentController(userService *UserService, commService *CommentService) *CommentController {
	return &CommentController{
		userService:    userService,
		commentService: commService,
	}
}

func (ps *CommentController) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users/posts/{id}/comments", WithJWTAuth(ps.handleCreateComment, ps.userService)).Methods("POST")
	r.HandleFunc("/users/posts/{id}/comments", WithJWTAuth(ps.handleGetComments, ps.userService)).Methods("GET")
	r.HandleFunc("/users/posts/comments/{id}", WithJWTAuth(ps.handleUpdateComments, ps.userService)).Methods("PUT")
	r.HandleFunc("/users/posts/comments/{id}", WithJWTAuth(ps.handleDeleteComments, ps.userService)).Methods("DELETE")
}

func (ps *CommentController) handleCreateComment(w http.ResponseWriter, r *http.Request) {
	postId, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing post Id param %v\n", "PostService ", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Error parsing post Id param"))
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
		log.Printf("%-15s ==> 😱 Error getting authenticated user Id %v\n", "PostService ", err)
		WriteJson(w, http.StatusUnauthorized, NewErrorResponse("Error getting authenticated user Id"))
		return
	}

	cResp, err := ps.commentService.CreateComment(postId, userId, cReq)
	if err != nil {
		log.Printf("%-15s ==> 🤯 Error creating comment in store %v\n", "PostService ", err)
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error creating comment in store"))
		return
	}

	log.Printf("%-15s ==> 🎉 Successfully created comment\n", "PostService")
	WriteJson(w, http.StatusCreated, cResp)
}

func (ps *CommentController) handleGetComments(w http.ResponseWriter, r *http.Request) {
	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id para %v\n", "PostService ", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Error parsing Id param"))
		return
	}

	c, err := ps.commentService.GetCommentsByPostId(id)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error getting comment by Id from stor %v\n", "PostService ", err)
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error getting comment by Id from store"))
		return
	}

	log.Printf("%-15s ==> 🎉 Successfully got comment by Id\n", "PostService!")
	WriteJson(w, http.StatusOK, c)
}

func (ps *CommentController) handleUpdateComments(w http.ResponseWriter, r *http.Request) {
	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id para %v\n", "PostService ", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Error parsing Id param"))
		return
	}

	c, err := readCommentReqType(r)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error reading comment request %v\n", "PostService ", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Error reading comment request"))
		return
	}

	cr, err := ps.commentService.UpdateCommentById(id, c)
	if err != nil {
		log.Printf("%-15s ==> 😱 Error updating comment by Id in stor %v\n", "PostService ", err)
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error updating comment by Id in store"))
		return
	}

	log.Printf("%-15s ==> 🎉 Successfully updated comment by Id\n", "PostService")
	WriteJson(w, http.StatusOK, cr)
}

func (ps *CommentController) handleDeleteComments(w http.ResponseWriter, r *http.Request) {
	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id para\n ", "PostService")
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Error parsing Id param"))
		return
	}

	err = ps.commentService.DeleteCommentById(id)
	if err != nil {
		log.Printf("%-15s ==> 😱 Error deleting comment by Id from stor\n ", "PostService")
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error deleting comment by Id from store"))
		return
	}

	log.Printf("%-15s ==> 🎉 Successfully deleted comment by Id\n", "PostService")
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

package main

import (
	"io"
	"log"
	"net/http"
)

type CommentController struct {
	commentService *CommentService
	sCtx           *SecurityContextHolder
}

func NewCommentController(commService *CommentService, sCtx *SecurityContextHolder) *CommentController {
	return &CommentController{
		commentService: commService,
		sCtx:           sCtx,
	}
}

func (ps *CommentController) RegisterRoutes(r *http.ServeMux) {
	middlewareStack := func(handler apiHandler) http.HandlerFunc {
		return MiddlewareChain(
			handler,
			LoggerMiddleware,
			ErrorHandler,
			ps.sCtx.WithJWTAuth,
		)
	}
	r.HandleFunc("POST /posts/{id}/comments", middlewareStack(ps.handleCreateComment))
	r.HandleFunc("GET /posts/{id}/comments", middlewareStack(ps.handleGetComments))
	r.HandleFunc("PUT /comments/{id}", middlewareStack(ps.handleUpdateComments))
	r.HandleFunc("DELETE /comments/{id}", middlewareStack(ps.handleDeleteComments))
}

func (ps *CommentController) handleCreateComment(w http.ResponseWriter, r *http.Request) error {
	postId, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing post Id param %v\n", "PostService ", err)
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Not valid ID param",
		}
	}

	cReq, err := readCommentReqType(r)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error reading comment request %v\n", "PostService ", err)
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Error reading comment request",
		}
	}

	userId, err := ps.sCtx.GetAuthUserId(r)
	if err != nil {
		log.Printf("%-15s ==> 😱 Error getting authenticated user Id %v\n", "PostService ", err)
		return err
	}

	cResp, err := ps.commentService.CreateComment(postId, userId, cReq)
	if err != nil {
		log.Printf("%-15s ==> 🤯 Error creating comment in store %v\n", "PostService ", err)
		return err
	}

	log.Printf("%-15s ==> 🎉 Successfully created comment\n", "PostService")

	return WriteJson(w, http.StatusCreated, cResp)
}

func (ps *CommentController) handleGetComments(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id para %v\n", "PostService ", err)
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Error parsing Id param",
		}
	}

	c, err := ps.commentService.GetCommentsByPostId(id)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error getting comment by Id from stor %v\n", "PostService ", err)
		return err
	}

	log.Printf("%-15s ==> 🎉 Successfully got comment by Id\n", "PostService!")

	return WriteJson(w, http.StatusOK, c)
}

func (ps *CommentController) handleUpdateComments(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id para %v\n", "PostService ", err)
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Not valid ID param",
		}

	}

	c, err := readCommentReqType(r)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error reading comment request %v\n", "PostService ", err)
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Not valid ID param",
		}

	}

	cr, err := ps.commentService.UpdateCommentById(id, c)
	if err != nil {
		log.Printf("%-15s ==> 😱 Error updating comment by Id in stor %v\n", "PostService ", err)
		return err
	}

	log.Printf("%-15s ==> 🎉 Successfully updated comment by Id\n", "PostService")

	return WriteJson(w, http.StatusOK, cr)
}

func (ps *CommentController) handleDeleteComments(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id para\n ", "PostService")
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Not valid ID param",
		}

	}

	err = ps.commentService.DeleteCommentById(id)
	if err != nil {
		log.Printf("%-15s ==> 😱 Error deleting comment by Id from stor\n ", "PostService")
		return err
	}

	log.Printf("%-15s ==> 🎉 Successfully deleted comment by Id\n", "PostService")

	return WriteJson(w, http.StatusNoContent, nil)
}

func readCommentReqType(r *http.Request) (*CommentRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	c, err := Unmarshal[CommentRequest](body)
	if err != nil {
		return nil, err
	}

	return c, nil
}

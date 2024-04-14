package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	db "github.com/malyshEvhen/meow_mingle/db/sqlc"
)


func (rr *Router) handleCreateComment(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	postId, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing post Id param %v\n", "PostService ", err)
		return err
	}

	params, err := readCommentCreateParams(r)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error reading comment request %v\n", "PostService ", err)
		return err
	}

	userId, err := getAuthUserId(r)
	if err != nil {
		log.Printf("%-15s ==> 😱 Error getting authenticated user Id %v\n", "PostService ", err)
		return err
	}

	params.AuthorID = userId
	params.PostID = postId

	cResp, err := rr.store.CreateComment(ctx, *params)
	if err != nil {
		log.Printf("%-15s ==> 🤯 Error creating comment in store %v\n", "PostService ", err)
		return err
	}

	log.Printf("%-15s ==> 🎉 Successfully created comment\n", "PostService")

	return WriteJson(w, http.StatusCreated, cResp)
}

func (rr *Router) handleGetComments(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id para %v\n", "PostService ", err)
		return err
	}

	comments, err := rr.store.ListPostComments(ctx, id)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error getting comment by Id from stor %v\n", "PostService ", err)
		return err
	}

	log.Printf("%-15s ==> 🎉 Successfully got comment by Id\n", "PostService!")

	return WriteJson(w, http.StatusOK, comments)
}

func (rr *Router) handleUpdateComments(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id para %v\n", "PostService ", err)
		return err

	}

	c, err := readCommentCreateParams(r)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error reading comment request %v\n", "PostService ", err)
		return err

	}

	params := &db.UpdateCommentParams{
		ID:      id,
		Content: c.Content,
	}

	cr, err := rr.store.UpdateComment(ctx, *params)
	if err != nil {
		log.Printf("%-15s ==> 😱 Error updating comment by Id in stor %v\n", "PostService ", err)
		return err
	}

	log.Printf("%-15s ==> 🎉 Successfully updated comment by Id\n", "PostService")

	return WriteJson(w, http.StatusOK, cr)
}

func (rr *Router) handleDeleteComments(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id para\n ", "PostService")
		return err

	}

	err = rr.store.DeleteComment(ctx, id)
	if err != nil {
		log.Printf("%-15s ==> 😱 Error deleting comment by Id from stor\n ", "PostService")
		return err
	}

	log.Printf("%-15s ==> 🎉 Successfully deleted comment by Id\n", "PostService")

	return WriteJson(w, http.StatusNoContent, nil)
}

func readCommentCreateParams(r *http.Request) (*db.CreateCommentParams, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	c, err := Unmarshal[db.CreateCommentParams](body)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (rr *Router) handleLikeComment(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIdParam(r)
	if err != nil {
		return err
	}

	userId, err := getAuthUserId(r)
	if err != nil {
		return err
	}

	params := db.CreateCommentLikeParams{
		UserID:    userId,
		CommentID: id,
	}

	if err := rr.store.CreateCommentLike(context.Background(), params); err != nil {
		return err
	}

	return WriteJson(w, http.StatusNoContent, nil)
}

func (rr *Router) handleRemoveLikeFromComment(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIdParam(r)
	if err != nil {
		return err
	}

	userId, err := getAuthUserId(r)
	if err != nil {
		return err
	}

	comment, err := rr.store.GetComment(context.Background(), id)
	if err != nil {
		return err
	}

	if comment.AuthorID != userId {
		return fmt.Errorf("user with ID: %d can not modify post of author with ID: %d", userId, comment.ID)
	}

	params := db.DeleteCommentLikeParams{
		CommentID: id,
		UserID:    userId,
	}

	if err := rr.store.DeleteCommentLike(context.Background(), params); err != nil {
		return err
	}

	return WriteJson(w, http.StatusNoContent, nil)
}

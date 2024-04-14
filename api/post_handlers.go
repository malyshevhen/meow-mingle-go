package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	db "github.com/malyshEvhen/meow_mingle/db/sqlc"
	"github.com/malyshEvhen/meow_mingle/errors"
)

func (rr *Router) handleCreatePost(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	postRequest, err := readPostReqType(r)
	if err != nil {
		log.Printf("%-15s ==> 😞 Error reading post request: %v\n", "Post Handler", err)
		return err
	}

	userId, err := getAuthUserId(r)
	if err != nil {
		log.Printf("%-15s ==> 😱 Error getting user Id from token %v\n", "Post Handler ", err)
		return err
	}

	savedPost, err := rr.store.CreatePostTx(ctx, userId, postRequest.Content)
	if err != nil {
		log.Printf("%-15s ==> 🤯 Error creating post in store %v\n", "Post Handler", err)
		return err
	}

	log.Printf("%-15s ==> 🎉 Successfully created new post\n", "Post Handler")

	return WriteJson(w, http.StatusCreated, savedPost)
}

func (rr *Router) handleGetUserPosts(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id param %v\n", "Post Handler", err)
		return err
	}

	postResponses, err := rr.store.ListUserPosts(ctx, id)
	if err != nil {
		return err
	}

	log.Printf("%-15s ==> 🤩 Successfully retrieved user posts\n", "Post Handler")

	return WriteJson(w, http.StatusOK, postResponses)
}

func (rr *Router) handleGetPostsById(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id para:%v\n", "Post Handler", err)
		return err
	}

	post, err := rr.store.GetPost(ctx, id)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error getting post by Id from stor:%v\n", "Post Handler", err)
		return err
	}

	log.Printf("%-15s ==> 🤩 Successfully retrieved post by Id\n", "Post Handler")

	return WriteJson(w, http.StatusOK, post)
}

func (rr *Router) handleUpdatePostsById(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id para %v\n", "Post Handler", err)
		return err
	}

	postRequest, err := readPostReqType(r)
	if err != nil {
		log.Printf("%-15s ==> 😫 Error reading post request %v\n", "Post Handler", err)
		return err
	}

	params := &db.UpdatePostParams{
		ID:      id,
		Content: postRequest.Content,
	}

	postResponse, err := rr.store.UpdatePost(ctx, *params)
	if err != nil {
		log.Printf("%-15s ==> 🤯 Error updating post by Id in store %v\n", "Post Handler", err)
		return err
	}

	log.Printf("%-15s ==> 🎉 Successfully updated post by Id\n", "Post Handler")

	return WriteJson(w, http.StatusOK, postResponse)
}

func (rr *Router) handleDeletePostsById(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	id, err := parseIdParam(r)
	if err != nil {
		log.Printf("%-15s ==> 😿 Error parsing Id param %v\n", "Post Handler", err)
		return err
	}

	if err := rr.store.DeletePost(ctx, id); err != nil {
		log.Printf("%-15s ==> 😫 Error deleting post by Id from store %v\n", "Post Handler", err)
		return err
	}

	log.Printf("%-15s ==> 🗑️ Successfully deleted post by Id\n", "Post Handler")

	return WriteJson(w, http.StatusNoContent, nil)
}

func readPostReqType(r *http.Request) (*db.CreatePostParams, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.NewValidationError("parameter ID is not valid")
	}
	defer r.Body.Close()

	p, err := Unmarshal[db.CreatePostParams](body)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (rr *Router) handleLikePost(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIdParam(r)
	if err != nil {
		return err
	}

	userId, err := getAuthUserId(r)
	if err != nil {
		return err
	}

	params := db.CreatePostLikeParams{
		UserID: userId,
		PostID: id,
	}

	if err := rr.store.CreatePostLike(context.Background(), params); err != nil {
		return err
	}

	return WriteJson(w, http.StatusNoContent, nil)
}

func (rr *Router) handleRemoveLikeFromPost(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIdParam(r)
	if err != nil {
		return err
	}

	userId, err := getAuthUserId(r)
	if err != nil {
		return err
	}

	post, err := rr.store.GetPost(context.Background(), id)
	if err != nil {
		return err
	}

	if post.AuthorID != userId {
		return fmt.Errorf("user with ID: %d can not modify post of author with ID: %d", userId, post.AuthorID)
	}

	params := db.DeletePostLikeParams{
		PostID: id,
		UserID: userId,
	}

	if err := rr.store.DeletePostLike(context.Background(), params); err != nil {
		return err
	}

	return WriteJson(w, http.StatusNoContent, nil)
}

package api

import (
	"context"
	"log"
	"net/http"

	db "github.com/malyshEvhen/meow_mingle/db/sqlc"
)

func handleCreatePost(store db.IStore) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		params, err := readCreatePostParams(r)
		if err != nil {
			log.Printf("%-15s ==> Error reading post request: %v\n", "Post Handler", err)
			return err
		}

		userId, err := getAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> Error getting user Id from token %v\n", "Post Handler ", err)
			return err
		}

		params.AuthorID = userId

		if err := Validate(params); err != nil {
			return err
		}

		savedPost, err := store.CreatePostTx(ctx, *params)
		if err != nil {
			log.Printf("%-15s ==> Error creating post in store %v\n", "Post Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully created new post\n", "Post Handler")

		return WriteJson(w, http.StatusCreated, savedPost)
	}
}

func handleGetUserPosts(store db.IStore) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := ParseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id param %v\n", "Post Handler", err)
			return err
		}

		postResponses, err := store.ListUserPostsTx(ctx, id)
		if err != nil {
			return err
		}

		log.Printf("%-15s ==> Successfully retrieved user posts\n", "Post Handler")

		return WriteJson(w, http.StatusOK, postResponses)
	}
}

func handleGetPostsById(store db.IStore) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := ParseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id para:%v\n", "Post Handler", err)
			return err
		}

		post, err := store.GetPostTx(ctx, id)
		if err != nil {
			log.Printf("%-15s ==> Error getting post by Id from stor:%v\n", "Post Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully retrieved post by Id\n", "Post Handler")

		return WriteJson(w, http.StatusOK, post)
	}
}

func handleUpdatePostsById(store db.IStore) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := ParseIdParam(r)
		if err != nil {
			return err
		}

		params, err := readUpdatePostParams(r)
		if err != nil {
			return err
		}
		params.ID = id

		if err := Validate(params); err != nil {
			return err
		}

		userId, err := getAuthUserId(r)
		if err != nil {
			return err
		}

		postResponse, err := store.UpdatePostTx(ctx, userId, *params)
		if err != nil {
			log.Printf("%-15s ==> Error updating post by Id in store %v\n", "Post Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully updated post by Id\n", "Post Handler")

		return WriteJson(w, http.StatusOK, postResponse)
	}
}

func handleDeletePostsById(store db.IStore) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := ParseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id param %v\n", "Post Handler", err)
			return err
		}

		userId, err := getAuthUserId(r)
		if err != nil {
			return err
		}

		if err := store.DeletePostTx(ctx, userId, id); err != nil {
			log.Printf("%-15s ==> Error deleting post by Id from store %v\n", "Post Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully deleted post by Id\n", "Post Handler")

		return WriteJson(w, http.StatusNoContent, nil)
	}
}

func handleLikePost(store db.IStore) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := ParseIdParam(r)
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

		if err := Validate(params); err != nil {
			return err
		}

		if err := store.CreatePostLikeTx(context.Background(), params); err != nil {
			return err
		}

		return WriteJson(w, http.StatusNoContent, nil)
	}
}

func handleRemoveLikeFromPost(store db.IStore) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := ParseIdParam(r)
		if err != nil {
			return err
		}

		userId, err := getAuthUserId(r)
		if err != nil {
			return err
		}

		params := db.DeletePostLikeParams{
			PostID: id,
			UserID: userId,
		}

		if err := Validate(params); err != nil {
			return err
		}

		if err := store.DeletePostLikeTx(ctx, params); err != nil {
			return err
		}

		return WriteJson(w, http.StatusNoContent, nil)
	}
}

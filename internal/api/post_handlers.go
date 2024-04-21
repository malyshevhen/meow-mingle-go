package api

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/internal/errors"
	"github.com/malyshEvhen/meow_mingle/internal/types"
	"github.com/malyshEvhen/meow_mingle/internal/utils"
)

func HandleCreatePost(store db.IStore) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		params, err := readCreatePostParams(r)
		if err != nil {
			log.Printf("%-15s ==> Error reading post request: %v\n", "Post Handler", err)
			return err
		}

		userId, err := utils.GetAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> Error getting user Id from token %v\n", "Post Handler ", err)
			return err
		}

		params.AuthorID = userId

		if err := utils.Validate(params); err != nil {
			return err
		}

		savedPost, err := store.CreatePostTx(ctx, *params)
		if err != nil {
			log.Printf("%-15s ==> Error creating post in store %v\n", "Post Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully created new post\n", "Post Handler")

		return utils.WriteJson(w, http.StatusCreated, savedPost)
	}
}

func HandleGetUserPosts(store db.IStore) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := utils.ParseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id param %v\n", "Post Handler", err)
			return err
		}

		postResponses, err := store.ListUserPostsTx(ctx, id)
		if err != nil {
			return err
		}

		log.Printf("%-15s ==> Successfully retrieved user posts\n", "Post Handler")

		return utils.WriteJson(w, http.StatusOK, postResponses)
	}
}

func HandleGetPostsById(store db.IStore) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := utils.ParseIdParam(r)
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

		return utils.WriteJson(w, http.StatusOK, post)
	}
}

func HandleUpdatePostsById(store db.IStore) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := utils.ParseIdParam(r)
		if err != nil {
			return err
		}

		params, err := readUpdatePostParams(r)
		if err != nil {
			return err
		}
		params.ID = id

		if err := utils.Validate(params); err != nil {
			return err
		}

		userId, err := utils.GetAuthUserId(r)
		if err != nil {
			return err
		}

		postResponse, err := store.UpdatePostTx(ctx, userId, *params)
		if err != nil {
			log.Printf("%-15s ==> Error updating post by Id in store %v\n", "Post Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully updated post by Id\n", "Post Handler")

		return utils.WriteJson(w, http.StatusOK, postResponse)
	}
}

func HandleDeletePostsById(store db.IStore) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := utils.ParseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id param %v\n", "Post Handler", err)
			return err
		}

		userId, err := utils.GetAuthUserId(r)
		if err != nil {
			return err
		}

		if err := store.DeletePostTx(ctx, userId, id); err != nil {
			log.Printf("%-15s ==> Error deleting post by Id from store %v\n", "Post Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully deleted post by Id\n", "Post Handler")

		return utils.WriteJson(w, http.StatusNoContent, nil)
	}
}

func HandleLikePost(store db.IStore) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := utils.ParseIdParam(r)
		if err != nil {
			return err
		}

		userId, err := utils.GetAuthUserId(r)
		if err != nil {
			return err
		}

		params := db.CreatePostLikeParams{
			UserID: userId,
			PostID: id,
		}

		if err := utils.Validate(params); err != nil {
			return err
		}

		if err := store.CreatePostLikeTx(context.Background(), params); err != nil {
			return err
		}

		return utils.WriteJson(w, http.StatusNoContent, nil)
	}
}

func HandleRemoveLikeFromPost(store db.IStore) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := utils.ParseIdParam(r)
		if err != nil {
			return err
		}

		userId, err := utils.GetAuthUserId(r)
		if err != nil {
			return err
		}

		params := db.DeletePostLikeParams{
			PostID: id,
			UserID: userId,
		}

		if err := utils.Validate(params); err != nil {
			return err
		}

		if err := store.DeletePostLikeTx(ctx, params); err != nil {
			return err
		}

		return utils.WriteJson(w, http.StatusNoContent, nil)
	}
}

func readCreatePostParams(r *http.Request) (*db.CreatePostParams, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.NewValidationError("parameter ID is not valid")
	}
	defer r.Body.Close()

	p, err := utils.Unmarshal[db.CreatePostParams](body)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func readUpdatePostParams(r *http.Request) (*db.UpdatePostParams, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("%-15s ==> Error reading post request %v\n", "Post Handler", err)
		return nil, errors.NewValidationError("parameter ID is not valid")
	}
	defer r.Body.Close()

	p, err := utils.Unmarshal[db.UpdatePostParams](body)
	if err != nil {
		log.Printf("%-15s ==> Error reading post request %v\n", "Post Handler", err)
		return nil, err
	}

	return &p, nil
}

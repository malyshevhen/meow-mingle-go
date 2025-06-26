package api

import (
	"log"
	"net/http"

	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/pkg/api"
	"github.com/malyshEvhen/meow_mingle/pkg/auth"
	"github.com/malyshEvhen/meow_mingle/pkg/errors"
)

func handleCreatePost(postRepo db.IPostRepository) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		postContent, err := readReqBody[ContentForm](r)
		if err != nil {
			log.Printf("%-15s ==> Error reading post request: %v\n", "Post Handler", err)
			return err
		}

		params := Map(postContent, func(s ContentForm) db.CreatePostParams {
			return db.CreatePostParams{Content: postContent.Content}
		})

		userId, err := auth.GetAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> Error getting user Id from token %v\n", "Post Handler ", err)
			return err
		}

		params.AuthorID = userId

		savedPost, err := postRepo.CreatePost(ctx, params)
		if err != nil {
			log.Printf("%-15s ==> Error creating post in store %v\n", "Post Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully created new post\n", "Post Handler")

		return writeJson(w, http.StatusCreated, savedPost)
	}
}

func handleGetPosts(postRepo db.IPostRepository) api.Handler {
	// TODO: implement
	return nil
}

func handleGetPostById(postRepo db.IPostRepository) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		id, err := parseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id para:%v\n", "Post Handler", err)
			return err
		}

		post, err := postRepo.GetPost(ctx, id)
		if err != nil {
			log.Printf("%-15s ==> Error getting post by Id from stor:%v\n", "Post Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully retrieved post by Id\n", "Post Handler")

		return writeJson(w, http.StatusOK, post)
	}
}

func handleUpdatePostById(postRepo db.IPostRepository) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		id, err := parseIdParam(r)
		if err != nil {
			return err
		}

		postContent, err := readReqBody[ContentForm](r)
		if err != nil {
			log.Printf("%-15s ==> Error reading update request: %v\n", "Post Handler", err)
			return err
		}

		userId, err := auth.GetAuthUserId(r)
		if err != nil {
			return err
		}

		params := Map(postContent, func(content ContentForm) db.UpdatePostParams {
			return db.UpdatePostParams{
				ID:       id,
				Content:  content.Content,
				AuthorId: userId,
			}
		})

		postResponse, err := postRepo.UpdatePost(ctx, params)
		if err != nil {
			log.Printf("%-15s ==> Error updating post by Id in store %v\n", "Post Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully updated post by Id\n", "Post Handler")

		return writeJson(w, http.StatusOK, postResponse)
	}
}

func handleDeletePostById(postRepo db.IPostRepository) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		id, err := parseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id param %v\n", "Post Handler", err)
			return err
		}

		userId, err := auth.GetAuthUserId(r)
		if err != nil {
			return err
		}

		if err := postRepo.DeletePost(ctx, userId, id); err != nil {
			log.Printf("%-15s ==> Error deleting post by Id from store %v\n", "Post Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully deleted post by Id\n", "Post Handler")

		return writeJson(w, http.StatusNoContent, nil)
	}
}

func handleGetFeed(postRepo db.IPostRepository) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		authUserID, err := auth.GetAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> No authenticated user found", "User Handler")
			return err
		}

		feed, err := postRepo.GetFeed(ctx, authUserID)
		if err != nil {
			return err
		}

		return writeJson(w, http.StatusOK, feed)
	}
}

func handleUsersFeed(postRepo db.IPostRepository) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		id, err := parseIdParam(r)
		if err != nil {
			return errors.NewValidationError("ID parameter is invalid")
		}

		feed, err := postRepo.GetFeed(ctx, id)
		if err != nil {
			return err
		}
		return writeJson(w, http.StatusOK, feed)
	}
}

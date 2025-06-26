package api

import (
	"log"
	"net/http"

	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/pkg/api"
	"github.com/malyshEvhen/meow_mingle/pkg/auth"
	"github.com/malyshEvhen/meow_mingle/pkg/errors"
)

func handleCreatePost(postService app.PostService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		req, err := readReqBody[ContentForm](r)
		if err != nil {
			log.Printf("%-15s ==> Error reading post request: %v\n", "Post Handler", err)
			return err
		}

		userId, err := auth.GetAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> Error getting user Id from token %v\n", "Post Handler ", err)
			return err
		}

		// TODO: add New function to Post struct with validation
		post := app.Post{
			Content:  req.Content,
			AuthorID: userId,
		}

		if err := postService.CreatePost(ctx, &post); err != nil {
			log.Printf("%-15s ==> Error creating post in store %v\n", "Post Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully created new post\n", "Post Handler")

		return writeJson(w, http.StatusCreated, post)
	}
}

func handleGetPosts(postService app.PostService) api.Handler {
	// TODO: implement
	return nil
}

func handleGetPostById(postService app.PostService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		id, err := parseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id para:%v\n", "Post Handler", err)
			return err
		}

		post, err := postService.GetPost(ctx, id)
		if err != nil {
			log.Printf("%-15s ==> Error getting post by Id from stor:%v\n", "Post Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully retrieved post by Id\n", "Post Handler")

		return writeJson(w, http.StatusOK, post)
	}
}

func handleUpdatePostById(postService app.PostService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		id, err := parseIdParam(r)
		if err != nil {
			return err
		}

		req, err := readReqBody[ContentForm](r)
		if err != nil {
			log.Printf("%-15s ==> Error reading update request: %v\n", "Post Handler", err)
			return err
		}

		if err := postService.UpdatePost(ctx, id, req.Content); err != nil {
			log.Printf("%-15s ==> Error updating post by Id in store %v\n", "Post Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully updated post by Id\n", "Post Handler")

		return writeJson(w, http.StatusNoContent, nil)
	}
}

func handleDeletePostById(postService app.PostService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		id, err := parseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id param %v\n", "Post Handler", err)
			return err
		}

		if err := postService.DeletePost(ctx, id); err != nil {
			log.Printf("%-15s ==> Error deleting post by Id from store %v\n", "Post Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully deleted post by Id\n", "Post Handler")

		return writeJson(w, http.StatusNoContent, nil)
	}
}

func handleGetFeed(postService app.PostService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		authUserID, err := auth.GetAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> No authenticated user found", "User Handler")
			return err
		}

		feed, err := postService.GetFeed(ctx, authUserID)
		if err != nil {
			return err
		}

		return writeJson(w, http.StatusOK, feed)
	}
}

func handleUsersFeed(postService app.PostRepository) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		id, err := parseIdParam(r)
		if err != nil {
			return errors.NewValidationError("ID parameter is invalid")
		}

		feed, err := postService.GetFeed(ctx, id)
		if err != nil {
			return err
		}
		return writeJson(w, http.StatusOK, feed)
	}
}

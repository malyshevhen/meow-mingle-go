package api

import (
	"log"
	"net/http"

	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/pkg/api"
	"github.com/malyshEvhen/meow_mingle/pkg/auth"
)

func handleCreatePost(postService app.PostService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		req, err := readBody[ContentForm](r)
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

		if err := postService.Create(ctx, &post); err != nil {
			log.Printf("%-15s ==> Error creating post in store %v\n", "Post Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully created new post\n", "Post Handler")

		return writeJSON(w, http.StatusCreated, post)
	}
}

func handleGetPosts(postService app.PostService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		profileId := r.URL.Query().Get("profileId")

		posts, err := postService.List(ctx, profileId)
		if err != nil {
			log.Printf("%-15s ==> Error getting posts from store %v\n", "Post Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully retrieved posts\n", "Post Handler")

		return writeJSON(w, http.StatusOK, posts)
	}
}

func handleGetPostById(postService app.PostService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		id, err := iaPathParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id para:%v\n", "Post Handler", err)
			return err
		}

		post, err := postService.Get(ctx, id)
		if err != nil {
			log.Printf("%-15s ==> Error getting post by Id from stor:%v\n", "Post Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully retrieved post by Id\n", "Post Handler")

		return writeJSON(w, http.StatusOK, post)
	}
}

func handleUpdatePostById(postService app.PostService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		id, err := iaPathParam(r)
		if err != nil {
			return err
		}

		req, err := readBody[ContentForm](r)
		if err != nil {
			log.Printf("%-15s ==> Error reading update request: %v\n", "Post Handler", err)
			return err
		}

		if err := postService.Edit(ctx, id, req.Content); err != nil {
			log.Printf("%-15s ==> Error updating post by Id in store %v\n", "Post Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully updated post by Id\n", "Post Handler")

		return writeJSON(w, http.StatusNoContent, nil)
	}
}

func handleDeletePostById(postService app.PostService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		id, err := iaPathParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id param %v\n", "Post Handler", err)
			return err
		}

		if err := postService.Delete(ctx, id); err != nil {
			log.Printf("%-15s ==> Error deleting post by Id from store %v\n", "Post Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully deleted post by Id\n", "Post Handler")

		return writeJSON(w, http.StatusNoContent, nil)
	}
}

func handleGetFeed(postService app.PostService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		feed, err := postService.Feed(ctx)
		if err != nil {
			return err
		}

		return writeJSON(w, http.StatusOK, feed)
	}
}

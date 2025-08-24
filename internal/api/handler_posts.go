package api

import (
	"net/http"

	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/pkg/api"
	"github.com/malyshEvhen/meow_mingle/pkg/logger"
)

func handleCreatePost(postService app.PostService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		logger := logger.GetLogger().WithComponent("post_handler")
		ctx := r.Context()

		req, err := readValidBody[ContentForm](r)
		if err != nil {
			logger.WithError(err).Error("Error reading post request")
			return err
		}

		post, err := app.NewPost(ctx, req.Content)
		if err != nil {
			logger.WithError(err).Error("Error creating post")
			return err
		}

		if err := postService.Create(ctx, post); err != nil {
			logger.WithError(err).Error("Error creating post")
			return err
		}

		logger.Info("Successfully created new post")

		return writeJSON(w, http.StatusCreated, post)
	}
}

func handleGetPosts(postService app.PostService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		logger := logger.GetLogger().WithComponent("post_handler")
		ctx := r.Context()

		profileId := r.URL.Query().Get("profileId")

		posts, err := postService.List(ctx, profileId)
		if err != nil {
			logger.WithError(err).Error("Error getting posts from store")
			return err
		}

		logger.Info("Successfully retrieved posts")

		return writeJSON(w, http.StatusOK, posts)
	}
}

func handleGetPostByID(postService app.PostService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		logger := logger.GetLogger().WithComponent("post_handler")
		ctx := r.Context()

		id, err := idPathParam(r)
		if err != nil {
			logger.WithError(err).Error("Error parsing Id parameter")
			return err
		}

		post, err := postService.Get(ctx, id)
		if err != nil {
			logger.WithError(err).Error("Error getting post by Id from store")
			return err
		}

		logger.Info("Successfully retrieved post by Id")

		return writeJSON(w, http.StatusOK, post)
	}
}

func handleUpdatePostByID(postService app.PostService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		logger := logger.GetLogger().WithComponent("post_handler")
		ctx := r.Context()

		id, err := idPathParam(r)
		if err != nil {
			logger.WithError(err).Error("Error parsing Id parameter")
			return err
		}

		req, err := readValidBody[ContentForm](r)
		if err != nil {
			logger.WithError(err).Error("Error reading update request")
			return err
		}

		if err := postService.Edit(ctx, id, req.Content); err != nil {
			logger.WithError(err).Error("Error updating post by Id in store")
			return err
		}

		logger.Info("Successfully updated post by Id")

		return writeJSON(w, http.StatusNoContent, nil)
	}
}

func handleDeletePostByID(postService app.PostService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		logger := logger.GetLogger().WithComponent("post_handler")
		ctx := r.Context()

		id, err := idPathParam(r)
		if err != nil {
			logger.WithError(err).Error("Error parsing Id parameter")
			return err
		}

		if err := postService.Delete(ctx, id); err != nil {
			logger.WithError(err).Error("Error deleting post by Id from store")
			return err
		}

		logger.Info("Successfully deleted post by Id")

		return writeJSON(w, http.StatusNoContent, nil)
	}
}

func handleGetFeed(postService app.PostService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		logger := logger.GetLogger().WithComponent("post_handler")
		ctx := r.Context()

		feed, err := postService.Feed(ctx)
		if err != nil {
			logger.WithError(err).Error("Error getting feed")
			return err
		}

		logger.Info("Successfully retrieved feed")

		return writeJSON(w, http.StatusOK, feed)
	}
}

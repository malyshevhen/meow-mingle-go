package api

import (
	"net/http"

	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/pkg/api"
	"github.com/malyshEvhen/meow_mingle/pkg/errors"
	"github.com/malyshEvhen/meow_mingle/pkg/logger"
)

func handleCreateComment(commentService app.CommentService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		logger := logger.GetLogger().WithComponent("comment_handler")
		ctx := r.Context()

		content, err := readBody[CreateCommentRequest](r)
		if err != nil {
			logger.WithError(err).Error("Error reading comment request")
			return err
		}

		comment, err := app.NewComment(ctx, content.PostID, content.Content)
		if err != nil {
			logger.WithError(err).Error("Error creating comment")
			return err
		}

		if err := commentService.Add(ctx, comment); err != nil {
			logger.WithError(err).Error("Error creating comment")
			return err
		}

		logger.Info("Successfully created comment")

		return writeJSON(w, http.StatusCreated, comment)
	}
}

// TODO: add pagination
func handleGetComments(commentService app.CommentService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		logger := logger.GetLogger().WithComponent("comment_handler")
		ctx := r.Context()

		postID := r.URL.Query().Get("postId")
		if len(postID) == 0 {
			err := errors.NewValidationError("Post ID is required")
			logger.WithError(err).Error("Error getting comment by Id")
			return err
		}

		comments, err := commentService.List(ctx, postID)
		if err != nil {
			logger.WithError(err).Error("Error getting comment by Id")
			return err
		}

		logger.Info("Successfully got comment by Id")

		return writeJSON(w, http.StatusOK, comments)
	}
}

func handleUpdateComment(commentService app.CommentService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		logger := logger.GetLogger().WithComponent("comment_handler")
		ctx := r.Context()

		id, err := iaPathParam(r)
		if err != nil {
			logger.WithError(err).Error("Error parsing Id parameter")
			return err
		}

		content, err := readBody[ContentForm](r)
		if err != nil {
			logger.WithError(err).Error("Error reading comment request")
			return err
		}

		if err := commentService.Update(ctx, id, content.Content); err != nil {
			logger.WithError(err).Error("Error updating comment by Id")
			return err
		}

		logger.Info("Successfully updated comment by Id")

		return writeJSON(w, http.StatusNoContent, nil)
	}
}

func handleDeleteComment(commentService app.CommentService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		logger := logger.GetLogger().WithComponent("comment_handler")
		ctx := r.Context()

		id, err := iaPathParam(r)
		if err != nil {
			logger.WithError(err).Error("Error parsing Id parameter")
			return err
		}

		err = commentService.Remove(ctx, id)
		if err != nil {
			logger.WithError(err).Error("Error deleting comment by Id")
			return err
		}

		logger.Info("Successfully deleted comment by Id")

		return writeJSON(w, http.StatusNoContent, nil)
	}
}

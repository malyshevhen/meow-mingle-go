package api

import (
	"net/http"

	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/pkg/api"
	"github.com/malyshEvhen/meow_mingle/pkg/logger"
)

func handleCreateReaction(reactionService app.ReactionService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		logger := logger.GetLogger().WithComponent("reaction_handler")
		ctx := r.Context()

		content, err := readValidBody[CreateCommentRequest](r)
		if err != nil {
			logger.WithError(err).Error("Error reading reaction request")
			return err
		}

		reaction := app.Reaction{
			Content: content.Content,
		}

		if err := reactionService.Add(ctx, &reaction); err != nil {
			logger.WithError(err).Error("Error creating reaction")
			return err
		}

		logger.Info("Successfully created reaction")

		return writeJSON(w, http.StatusNoContent, nil)
	}
}

func handleDeleteReaction(reactionService app.ReactionService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		logger := logger.GetLogger().WithComponent("reaction_handler")
		ctx := r.Context()

		id, err := idPathParam(r)
		if err != nil {
			logger.WithError(err).Error("Error getting reaction id")
			return err
		}

		if err := reactionService.Remove(ctx, id); err != nil {
			logger.WithError(err).Error("Error deleting reaction")
			return err
		}

		logger.Info("Successfully deleted reaction")

		return writeJSON(w, http.StatusNoContent, nil)
	}
}

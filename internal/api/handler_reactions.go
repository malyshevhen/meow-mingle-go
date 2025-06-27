package api

import (
	"log"
	"net/http"

	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/pkg/api"
)

func handleCreateReaction(reactionService app.ReactionService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		content, err := readBody[CreateCommentRequest](r)
		if err != nil {
			log.Printf("%-15s ==> Error reading comment request %v\n", "Comment Handler", err)
			return err
		}

		reaction := app.Reaction{
			Content: content.Content,
		}

		if err := reactionService.Add(ctx, &reaction); err != nil {
			log.Printf("%-15s ==> Error creating reaction %v\n", "Comment Handler", err)
			return err
		}

		return writeJSON(w, http.StatusNoContent, nil)
	}
}

func handleDeleteREaction(reactionService app.ReactionService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		id, err := iaPathParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error getting reaction id %v\n", "Comment Handler", err)
			return err
		}

		if err := reactionService.Remove(ctx, id); err != nil {
			log.Printf("%-15s ==> Error deleting reaction %v\n", "Comment Handler", err)
			return err
		}

		return writeJSON(w, http.StatusNoContent, nil)
	}
}

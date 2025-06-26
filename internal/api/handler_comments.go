package api

import (
	"log"
	"net/http"

	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/pkg/api"
	"github.com/malyshEvhen/meow_mingle/pkg/auth"
	"github.com/malyshEvhen/meow_mingle/pkg/errors"
)

func handleCreateComment(commentService app.CommentService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		content, err := readBody[CreateCommentRequest](r)
		if err != nil {
			log.Printf("%-15s ==> Error reading comment request %v\n", "Comment Handler", err)
			return err
		}

		userId, err := auth.GetAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> Error getting authenticated user Id %v\n", "Comment Handler", err)
			return err
		}

		// TODO: add New function to Comment struct with validation
		comment := app.Comment{
			Content:  content.Content,
			AuthorID: userId,
			PostID:   content.PostID,
		}

		if err := commentService.Create(ctx, &comment); err != nil {
			log.Printf("%-15s ==> Error creating comment in store %v\n", "Comment Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully created comment\n", "Comment Handler")

		return writeJSON(w, http.StatusCreated, comment)
	}
}

// TODO: add pagination
func handleGetComments(commentService app.CommentService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		postID := r.URL.Query().Get("postId")
		if len(postID) == 0 {
			return errors.NewValidationError("Post ID is required")
		}

		comments, err := commentService.List(ctx, postID)
		if err != nil {
			log.Printf(
				"%-15s ==> Error getting comment by Id from stor %v\n",
				"Comment Handler",
				err,
			)
			return err
		}

		log.Printf("%-15s ==> Successfully got comment by Id\n", "Comment Handler")

		return writeJSON(w, http.StatusOK, comments)
	}
}

func handleUpdateComment(commentService app.CommentService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		id, err := iaPathParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id para %v\n", "Comment Handler", err)
			return err
		}

		content, err := readBody[ContentForm](r)
		if err != nil {
			log.Printf("%-15s ==> Error reading comment request %v\n", "Comment Handler", err)
			return err
		}

		if err := commentService.Update(ctx, id, content.Content); err != nil {
			log.Printf(
				"%-15s ==> Error updating comment by Id in stor %v\n",
				"Comment Handler",
				err,
			)
			return err
		}

		log.Printf("%-15s ==> Successfully updated comment by Id\n", "Comment Handler")

		return writeJSON(w, http.StatusNoContent, nil)
	}
}

func handleDeleteComment(commentService app.CommentService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		id, err := iaPathParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id para\n ", "Comment Handler")
			return err
		}

		err = commentService.Delete(ctx, id)
		if err != nil {
			log.Printf("%-15s ==> Error deleting comment by Id from stor\n ", "Comment Handler")
			return err
		}

		log.Printf("%-15s ==> Successfully deleted comment by Id\n", "Comment Handler")

		return writeJSON(w, http.StatusNoContent, nil)
	}
}

package api

import (
	"log"
	"net/http"

	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/pkg/api"
	"github.com/malyshEvhen/meow_mingle/pkg/auth"
	"github.com/malyshEvhen/meow_mingle/pkg/errors"
)

func handleCreateComment(commentRepo db.ICommentRepository) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		content, err := readReqBody[CreateCommentRequest](r)
		if err != nil {
			log.Printf("%-15s ==> Error reading comment request %v\n", "Comment Handler", err)
			return err
		}

		userId, err := auth.GetAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> Error getting authenticated user Id %v\n", "Comment Handler", err)
			return err
		}

		params := Map(content, func(c CreateCommentRequest) db.CreateCommentParams {
			return db.CreateCommentParams{
				Content:  c.Content,
				PostID:   c.PostID,
				AuthorID: userId,
			}
		})

		comment, err := commentRepo.CreateComment(ctx, params)
		if err != nil {
			log.Printf("%-15s ==> Error creating comment in store %v\n", "Comment Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully created comment\n", "Comment Handler")

		return writeJson(w, http.StatusCreated, comment)
	}
}

func handleGetComments(commentRepo db.ICommentRepository) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		postID := r.URL.Query().Get("postId")
		if len(postID) == 0 {
			return errors.NewValidationError("Post ID is required")
		}

		comments, err := commentRepo.ListPostComments(ctx, postID)
		if err != nil {
			log.Printf(
				"%-15s ==> Error getting comment by Id from stor %v\n",
				"Comment Handler",
				err,
			)
			return err
		}

		log.Printf("%-15s ==> Successfully got comment by Id\n", "Comment Handler")

		return writeJson(w, http.StatusOK, comments)
	}
}

func handleUpdateComment(commentRepo db.ICommentRepository) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		id, err := parseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id para %v\n", "Comment Handler", err)
			return err
		}

		content, err := readReqBody[ContentForm](r)
		if err != nil {
			log.Printf("%-15s ==> Error reading comment request %v\n", "Comment Handler", err)
			return err
		}

		userId, err := auth.GetAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> Error getting authenticated user Id %v\n", "Comment Handler", err)
			return err
		}

		params := Map(content, func(c ContentForm) db.UpdateCommentParams {
			return db.UpdateCommentParams{
				ID:       id,
				Content:  c.Content,
				AuthorId: userId,
			}
		})

		comment, err := commentRepo.UpdateComment(ctx, params)
		if err != nil {
			log.Printf(
				"%-15s ==> Error updating comment by Id in stor %v\n",
				"Comment Handler",
				err,
			)
			return err
		}

		log.Printf("%-15s ==> Successfully updated comment by Id\n", "Comment Handler")

		return writeJson(w, http.StatusOK, comment)
	}
}

func handleDeleteComment(commentRepo db.ICommentRepository) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		id, err := parseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id para\n ", "Comment Handler")
			return err

		}

		userId, err := auth.GetAuthUserId(r)
		if err != nil {
			return err
		}

		err = commentRepo.DeleteComment(ctx, userId, id)
		if err != nil {
			log.Printf("%-15s ==> Error deleting comment by Id from stor\n ", "Comment Handler")
			return err
		}

		log.Printf("%-15s ==> Successfully deleted comment by Id\n", "Comment Handler")

		return writeJson(w, http.StatusNoContent, nil)
	}
}

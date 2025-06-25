package api

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/malyshEvhen/meow_mingle/internal/db"
)

func handleCreateComment(commentRepo db.ICommentRepository) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		postId, err := parseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing post Id param %v\n", "Comment Handler", err)
			return err
		}

		content, err := readReqBody[ContentForm](r)
		if err != nil {
			log.Printf("%-15s ==> Error reading comment request %v\n", "Comment Handler", err)
			return err
		}

		userId, err := GetAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> Error getting authenticated user Id %v\n", "Comment Handler", err)
			return err
		}

		params := Map(content, func(c ContentForm) db.CreateCommentParams {
			return db.CreateCommentParams{
				Content:  c.Content,
				AuthorID: userId,
				PostID:   postId,
			}
		})

		comment, err := commentRepo.CreateComment(ctx, params)
		if err != nil {
			log.Printf("%-15s ==> Error creating comment in store %v\n", "Comment Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully created comment\n", "Comment Handler")

		return WriteJson(w, http.StatusCreated, comment)
	}
}

func handleGetComments(commentRepo db.ICommentRepository) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id := mux.Vars(r)["id"]

		comments, err := commentRepo.ListPostComments(ctx, id)
		if err != nil {
			log.Printf(
				"%-15s ==> Error getting comment by Id from stor %v\n",
				"Comment Handler",
				err,
			)
			return err
		}

		log.Printf("%-15s ==> Successfully got comment by Id\n", "Comment Handler")

		return WriteJson(w, http.StatusOK, comments)
	}
}

func handleUpdateComments(commentRepo db.ICommentRepository) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

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

		userId, err := GetAuthUserId(r)
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

		return WriteJson(w, http.StatusOK, comment)
	}
}

func handleDeleteComments(commentRepo db.ICommentRepository) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := parseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id para\n ", "Comment Handler")
			return err

		}

		userId, err := GetAuthUserId(r)
		if err != nil {
			return err
		}

		err = commentRepo.DeleteComment(ctx, userId, id)
		if err != nil {
			log.Printf("%-15s ==> Error deleting comment by Id from stor\n ", "Comment Handler")
			return err
		}

		log.Printf("%-15s ==> Successfully deleted comment by Id\n", "Comment Handler")

		return WriteJson(w, http.StatusNoContent, nil)
	}
}

func handleLikeComment(commentRepo db.ICommentRepository) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := parseIdParam(r)
		if err != nil {
			return err
		}

		userId, err := GetAuthUserId(r)
		if err != nil {
			return err
		}

		params := db.CreateCommentLikeParams{
			UserID:    userId,
			CommentID: id,
		}

		if err := validate(params); err != nil {
			return err
		}

		if err := commentRepo.CreateCommentLike(ctx, params); err != nil {
			return err
		}

		return WriteJson(w, http.StatusNoContent, nil)
	}
}

func handleRemoveLikeFromComment(commentRepo db.ICommentRepository) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		commentId, err := parseIdParam(r)
		if err != nil {
			return err
		}

		userId, err := GetAuthUserId(r)
		if err != nil {
			return err
		}

		params := db.DeleteCommentLikeParams{
			CommentID: commentId,
			UserID:    userId,
		}

		if err := validate(params); err != nil {
			return err
		}

		if err := commentRepo.DeleteCommentLike(ctx, params); err != nil {
			return err
		}

		return WriteJson(w, http.StatusNoContent, nil)
	}
}

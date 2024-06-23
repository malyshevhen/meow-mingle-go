package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/internal/types"
	"github.com/malyshEvhen/meow_mingle/internal/utils"
)

type CommentHandler struct {
	commentRepo db.ICommentRepository
}

func NewCommentHandler(commentRepo db.ICommentRepository) *CommentHandler {
	return &CommentHandler{
		commentRepo: commentRepo,
	}
}

func (ch *CommentHandler) HandleCreateComment() types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		postId, err := utils.ParseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing post Id param %v\n", "Comment Handler", err)
			return err
		}

		content, err := ReadReqBody[ContentForm](r)
		if err != nil {
			log.Printf("%-15s ==> Error reading comment request %v\n", "Comment Handler", err)
			return err
		}

		userId, err := utils.GetAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> Error getting authenticated user Id %v\n", "Comment Handler", err)
			return err
		}

		params := utils.Map(content, func(c ContentForm) db.CreateCommentParams {
			return db.CreateCommentParams{
				Content:  c.Content,
				AuthorID: userId,
				PostID:   postId,
			}
		})

		comment, err := ch.commentRepo.CreateComment(ctx, params)
		if err != nil {
			log.Printf("%-15s ==> Error creating comment in store %v\n", "Comment Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully created comment\n", "Comment Handler")

		return utils.WriteJson(w, http.StatusCreated, comment)
	}
}

func (ch *CommentHandler) HandleGetComments() types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id := mux.Vars(r)["id"]

		comments, err := ch.commentRepo.ListPostComments(ctx, id)
		if err != nil {
			log.Printf(
				"%-15s ==> Error getting comment by Id from stor %v\n",
				"Comment Handler",
				err,
			)
			return err
		}

		log.Printf("%-15s ==> Successfully got comment by Id\n", "Comment Handler")

		return utils.WriteJson(w, http.StatusOK, comments)
	}
}

func (ch *CommentHandler) HandleUpdateComments() types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := utils.ParseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id para %v\n", "Comment Handler", err)
			return err
		}

		content, err := ReadReqBody[ContentForm](r)
		if err != nil {
			log.Printf("%-15s ==> Error reading comment request %v\n", "Comment Handler", err)
			return err
		}

		userId, err := utils.GetAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> Error getting authenticated user Id %v\n", "Comment Handler", err)
			return err
		}

		params := utils.Map(content, func(c ContentForm) db.UpdateCommentParams {
			return db.UpdateCommentParams{
				ID:       id,
				Content:  c.Content,
				AuthorId: userId,
			}
		})

		comment, err := ch.commentRepo.UpdateComment(ctx, params)
		if err != nil {
			log.Printf(
				"%-15s ==> Error updating comment by Id in stor %v\n",
				"Comment Handler",
				err,
			)
			return err
		}

		log.Printf("%-15s ==> Successfully updated comment by Id\n", "Comment Handler")

		return utils.WriteJson(w, http.StatusOK, comment)
	}
}

func (ch *CommentHandler) HandleDeleteComments() types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := utils.ParseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id para\n ", "Comment Handler")
			return err

		}

		userId, err := utils.GetAuthUserId(r)
		if err != nil {
			return err
		}

		err = ch.commentRepo.DeleteComment(ctx, userId, id)
		if err != nil {
			log.Printf("%-15s ==> Error deleting comment by Id from stor\n ", "Comment Handler")
			return err
		}

		log.Printf("%-15s ==> Successfully deleted comment by Id\n", "Comment Handler")

		return utils.WriteJson(w, http.StatusNoContent, nil)
	}
}

func (ch *CommentHandler) HandleLikeComment() types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := utils.ParseIdParam(r)
		if err != nil {
			return err
		}

		userId, err := utils.GetAuthUserId(r)
		if err != nil {
			return err
		}

		params := db.CreateCommentLikeParams{
			UserID:    userId,
			CommentID: id,
		}

		if err := utils.Validate(params); err != nil {
			return err
		}

		if err := ch.commentRepo.CreateCommentLike(ctx, params); err != nil {
			return err
		}

		return utils.WriteJson(w, http.StatusNoContent, nil)
	}
}

func (ch *CommentHandler) HandleRemoveLikeFromComment() types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		commentId, err := utils.ParseIdParam(r)
		if err != nil {
			return err
		}

		userId, err := utils.GetAuthUserId(r)
		if err != nil {
			return err
		}

		params := db.DeleteCommentLikeParams{
			CommentID: commentId,
			UserID:    userId,
		}

		if err := utils.Validate(params); err != nil {
			return err
		}

		if err := ch.commentRepo.DeleteCommentLike(ctx, params); err != nil {
			return err
		}

		return utils.WriteJson(w, http.StatusNoContent, nil)
	}
}

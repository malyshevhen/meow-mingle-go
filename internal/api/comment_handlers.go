package api

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/internal/utils"
)

func handleCreateComment(store db.IStore) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		postId, err := utils.ParseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing post Id param %v\n", "Comment Handler", err)
			return err
		}

		params, err := readCreateCommentParams(r)
		if err != nil {
			log.Printf("%-15s ==> Error reading comment request %v\n", "Comment Handler", err)
			return err
		}

		userId, err := utils.GetAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> Error getting authenticated user Id %v\n", "Comment Handler", err)
			return err
		}

		params.AuthorID = userId
		params.PostID = postId

		if err := utils.Validate(params); err != nil {
			return err
		}

		comment, err := store.CreateCommentTx(ctx, *params)
		if err != nil {
			log.Printf("%-15s ==> Error creating comment in store %v\n", "Comment Handler", err)
			return err
		}

		log.Printf("%-15s ==> Successfully created comment\n", "Comment Handler")

		return utils.WriteJson(w, http.StatusCreated, comment)
	}
}

func handleGetComments(store db.IStore) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := utils.ParseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id para %v\n", "Comment Handler", err)
			return err
		}

		comments, err := store.ListPostCommentsTx(ctx, id)
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

func handleUpdateComments(store db.IStore) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := utils.ParseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id para %v\n", "Comment Handler", err)
			return err

		}

		params, err := readUpdateCommentParams(r)
		if err != nil {
			log.Printf("%-15s ==> Error reading comment request %v\n", "Comment Handler", err)
			return err

		}

		params.ID = id

		if err := utils.Validate(params); err != nil {
			return err
		}

		userId, err := utils.GetAuthUserId(r)
		if err != nil {
			return err
		}

		comment, err := store.UpdateCommentTx(ctx, userId, *params)
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

func handleDeleteComments(store db.IStore) Handler {
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

		err = store.DeleteCommentTx(ctx, userId, id)
		if err != nil {
			log.Printf("%-15s ==> Error deleting comment by Id from stor\n ", "Comment Handler")
			return err
		}

		log.Printf("%-15s ==> Successfully deleted comment by Id\n", "Comment Handler")

		return utils.WriteJson(w, http.StatusNoContent, nil)
	}
}

func handleLikeComment(store db.IStore) Handler {
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

		if err := store.CreateCommentLikeTx(ctx, params); err != nil {
			return err
		}

		return utils.WriteJson(w, http.StatusNoContent, nil)
	}
}

func handleRemoveLikeFromComment(store db.IStore) Handler {
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

		if err := store.DeleteCommentLikeTx(ctx, params); err != nil {
			return err
		}

		return utils.WriteJson(w, http.StatusNoContent, nil)
	}
}

func readCreateCommentParams(r *http.Request) (*db.CreateCommentParams, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	c, err := utils.Unmarshal[db.CreateCommentParams](body)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func readUpdateCommentParams(r *http.Request) (*db.UpdateCommentParams, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	c, err := utils.Unmarshal[db.UpdateCommentParams](body)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

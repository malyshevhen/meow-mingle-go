package api

import (
	"context"
	"io"
	"log"
	"net/http"

	db "github.com/malyshEvhen/meow_mingle/db/sqlc"
)

func handleCreateComment(store db.IStore) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		postId, err := ParseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing post Id param %v\n", "PostService ", err)
			return err
		}

		params, err := readCreateCommentParams(r)
		if err != nil {
			log.Printf("%-15s ==> Error reading comment request %v\n", "PostService ", err)
			return err
		}

		userId, err := getAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> Error getting authenticated user Id %v\n", "PostService ", err)
			return err
		}

		params.AuthorID = userId
		params.PostID = postId

		comment, err := store.CreateCommentTx(ctx, *params)
		if err != nil {
			log.Printf("%-15s ==> Error creating comment in store %v\n", "PostService ", err)
			return err
		}

		log.Printf("%-15s ==> Successfully created comment\n", "PostService")

		return WriteJson(w, http.StatusCreated, comment)
	}
}

func handleGetComments(store db.IStore) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := ParseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id para %v\n", "PostService ", err)
			return err
		}

		comments, err := store.ListPostCommentsTx(ctx, id)
		if err != nil {
			log.Printf("%-15s ==> Error getting comment by Id from stor %v\n", "PostService ", err)
			return err
		}

		log.Printf("%-15s ==> Successfully got comment by Id\n", "PostService!")

		return WriteJson(w, http.StatusOK, comments)
	}
}

func handleUpdateComments(store db.IStore) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := ParseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id para %v\n", "PostService ", err)
			return err

		}

		params, err := readUpdateCommentParams(r)
		if err != nil {
			log.Printf("%-15s ==> Error reading comment request %v\n", "PostService ", err)
			return err

		}

		params.ID = id

		userId, err := getAuthUserId(r)
		if err != nil {
			return err
		}

		comment, err := store.UpdateCommentTx(ctx, userId, *params)
		if err != nil {
			log.Printf("%-15s ==> Error updating comment by Id in stor %v\n", "PostService ", err)
			return err
		}

		log.Printf("%-15s ==> Successfully updated comment by Id\n", "PostService")

		return WriteJson(w, http.StatusOK, comment)
	}
}

func handleDeleteComments(store db.IStore) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := ParseIdParam(r)
		if err != nil {
			log.Printf("%-15s ==> Error parsing Id para\n ", "PostService")
			return err

		}

		userId, err := getAuthUserId(r)
		if err != nil {
			return err
		}

		err = store.DeleteCommentTx(ctx, userId, id)
		if err != nil {
			log.Printf("%-15s ==> Error deleting comment by Id from stor\n ", "PostService")
			return err
		}

		log.Printf("%-15s ==> Successfully deleted comment by Id\n", "PostService")

		return WriteJson(w, http.StatusNoContent, nil)
	}
}

func handleLikeComment(store db.IStore) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := ParseIdParam(r)
		if err != nil {
			return err
		}

		userId, err := getAuthUserId(r)
		if err != nil {
			return err
		}

		params := db.CreateCommentLikeParams{
			UserID:    userId,
			CommentID: id,
		}

		if err := store.CreateCommentLikeTx(ctx, params); err != nil {
			return err
		}

		return WriteJson(w, http.StatusNoContent, nil)
	}
}

func handleRemoveLikeFromComment(store db.IStore) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		commentId, err := ParseIdParam(r)
		if err != nil {
			return err
		}

		userId, err := getAuthUserId(r)
		if err != nil {
			return err
		}

		if err := store.DeleteCommentLikeTx(ctx, db.DeleteCommentLikeParams{
			CommentID: commentId,
			UserID:    userId,
		}); err != nil {
			return err
		}

		return WriteJson(w, http.StatusNoContent, nil)
	}
}

func readCreateCommentParams(r *http.Request) (*db.CreateCommentParams, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	c, err := Unmarshal[db.CreateCommentParams](body)
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

	c, err := Unmarshal[db.UpdateCommentParams](body)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

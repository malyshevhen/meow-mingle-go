// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: post_likes.sql

package db

import (
	"context"
)

const createPostLike = `-- name: CreatePostLike :exec
INSERT INTO post_likes (
    user_id, post_id
) VALUES ($1, $2)
`

type CreatePostLikeParams struct {
	UserID int64 `json:"user_id"`
	PostID int64 `json:"post_id"`
}

func (q *Queries) CreatePostLike(ctx context.Context, arg CreatePostLikeParams) error {
	_, err := q.db.ExecContext(ctx, createPostLike, arg.UserID, arg.PostID)
	return err
}

const deletePostLike = `-- name: DeletePostLike :exec
DELETE FROM post_likes
WHERE post_id = $1 AND user_id = $2
`

type DeletePostLikeParams struct {
	PostID int64 `json:"post_id"`
	UserID int64 `json:"user_id"`
}

func (q *Queries) DeletePostLike(ctx context.Context, arg DeletePostLikeParams) error {
	_, err := q.db.ExecContext(ctx, deletePostLike, arg.PostID, arg.UserID)
	return err
}

const listPostLikes = `-- name: ListPostLikes :many
SELECT user_id, post_id FROM post_likes
WHERE post_id = $1
`

func (q *Queries) ListPostLikes(ctx context.Context, postID int64) ([]PostLike, error) {
	rows, err := q.db.QueryContext(ctx, listPostLikes, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PostLike{}
	for rows.Next() {
		var i PostLike
		if err := rows.Scan(&i.UserID, &i.PostID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
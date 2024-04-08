-- name: CreateCommentLike :exec
INSERT INTO comment_likes (
    user_id, comment_id
) VALUES ($1, $2);

-- name: ListCommentLikes :many
SELECT * FROM comment_likes
WHERE comment_id = $1;

-- name: DeleteCommentLike :exec
DELETE FROM comment_likes
WHERE comment_id = $1 AND user_id = $2;
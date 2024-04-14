
-- name: CreatePostLike :exec
INSERT INTO post_likes (
    user_id, post_id
) VALUES ($1, $2);

-- name: ListPostLikes :many
SELECT * FROM post_likes
WHERE post_id = $1;

-- name: CountPostLikes :one
SELECT COUNT(id) FROM post_likes
WHERE post_id = $1;

-- name: DeletePostLike :exec
DELETE FROM post_likes
WHERE post_id = $1 AND user_id = $2;
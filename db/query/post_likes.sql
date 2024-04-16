
-- name: CreatePostLike :exec
INSERT INTO post_likes (
    user_id, post_id
) VALUES ($1, $2);

-- name: GetPostLike :one
SELECT * FROM post_likes
WHERE post_id = $1 AND user_id = $2
LIMIT 1;

-- name: ListPostLikes :many
SELECT * FROM post_likes
WHERE post_id = $1;

-- name: CountPostLikes :one
SELECT COUNT(id) FROM post_likes
WHERE post_id = $1;

-- name: DeletePostLike :exec
DELETE FROM post_likes
WHERE post_id = $1 AND user_id = $2;
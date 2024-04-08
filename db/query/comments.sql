-- name: CreateComment :one
INSERT INTO comments (
    content, author_id, post_id
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetComment :one
SELECT * FROM comments
WHERE id = $1 LIMIT 1;

-- name: ListPostComments :many
SELECT * FROM comments 
WHERE post_id = $1;

-- name: UpdateComment :one
UPDATE comments
SET content = $2
WHERE id = $1
RETURNING *;

-- name: DeleteComment :exec
DELETE FROM comments
WHERE id = $1;

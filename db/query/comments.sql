-- name: CreateComment :one
INSERT INTO comments (
    content, author_id, post_id
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetComment :one
SELECT * FROM comment_info
WHERE id = $1
LIMIT 1;

-- name: ListPostComments :many
SELECT * FROM comment_info
WHERE post_id = $1
ORDER BY id;

-- name: GetCommentsAuthorID :one
SELECT c.author_id
FROM comments c
WHERE c.id = $1 LIMIT 1;

-- name: UpdateComment :one
UPDATE comments
SET content = $2
WHERE id = $1
RETURNING *;

-- name: DeleteComment :exec
DELETE FROM comments
WHERE id = $1;

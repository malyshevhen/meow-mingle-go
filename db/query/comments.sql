-- name: CreateComment :one
INSERT INTO comments (
    content, author_id, post_id
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetComment :one
SELECT
    c.id,
    c.author_id,
    c.content,
    c.created_at,
    c.updated_at,
    COALESCE(lc.count_likes, 0) as likes
FROM comments c
LEFT JOIN (
    SELECT comment_id, COUNT(*) as count_likes
    FROM comment_likes
    GROUP BY comment_id
) lc ON c.id = lc.comment_id
WHERE c.id = $1 LIMIT 1;

-- name: ListPostComments :many
SELECT
    c.id,
    c.author_id,
    c.content,
    c.created_at,
    c.updated_at,
    COALESCE(lc.count_likes, 0) as likes
FROM comments c
LEFT JOIN (
    SELECT comment_id, COUNT(*) as count_likes
    FROM comment_likes
    GROUP BY comment_id
) lc ON c.id = lc.comment_id
WHERE c.post_id = $1
ORDER BY c.id;

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

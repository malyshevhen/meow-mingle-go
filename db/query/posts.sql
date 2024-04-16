-- name: CreatePost :one
INSERT INTO posts (
    content, author_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetPost :one
SELECT
    p.id,
    p.author_id,
    p.content,
    p.created_at,
    p.updated_at,
    COALESCE(lc.count_likes, 0) as likes
FROM posts p
LEFT JOIN (
    SELECT post_id, COUNT(*) as count_likes
    FROM post_likes
    GROUP BY post_id
) lc ON p.id = lc.post_id
WHERE p.id = $1
LIMIT 1;

-- name: GetPostsAuthorID :one
SELECT p.author_id
FROM posts p
WHERE p.id = $1 LIMIT 1;

-- name: ListUserPosts :many
SELECT
    p.id,
    p.author_id,
    p.content,
    p.created_at,
    p.updated_at,
    COALESCE(lc.count_likes, 0) as likes
FROM posts p
LEFT JOIN (
    SELECT post_id, COUNT(*) as count_likes
    FROM post_likes
    GROUP BY post_id
) lc ON p.id = lc.post_id
WHERE p.author_id = $1
ORDER BY p.id;

-- name: UpdatePost :one
UPDATE posts
SET content = $2
WHERE id = $1
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1;

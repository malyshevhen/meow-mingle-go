-- name: CreatePost :one
INSERT INTO posts (
    content, author_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetPost :one
SELECT * FROM post_info
WHERE id = $1
LIMIT 1;

-- name: GetPostsAuthorID :one
SELECT p.author_id
FROM posts p
WHERE p.id = $1 LIMIT 1;

-- name: ListUserPosts :many
SELECT * FROM post_info
WHERE author_id = $1
ORDER BY id;

-- name: UpdatePost :one
UPDATE posts
SET content = $2
WHERE id = $1
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1;

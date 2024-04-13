-- name: CreateUser :one
INSERT INTO users (
	email, first_name, last_name, password
) VALUES (
	$1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: IsUserExists :one
SELECT COUNT(*) FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id;

-- name: UpdateUser :one
UPDATE users
SET first_name = $2, last_name = $3
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

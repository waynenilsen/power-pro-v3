-- name: CreateUser :exec
INSERT INTO users (id, created_at, updated_at)
VALUES (?, ?, ?);

-- name: GetUser :one
SELECT id, created_at, updated_at
FROM users
WHERE id = ?;

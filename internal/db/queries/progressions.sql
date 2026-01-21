-- name: CreateProgression :exec
INSERT INTO progressions (id, name, type, parameters, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?);

-- name: GetProgression :one
SELECT id, name, type, parameters, created_at, updated_at
FROM progressions
WHERE id = ?;

-- name: ListProgressions :many
SELECT id, name, type, parameters, created_at, updated_at
FROM progressions
ORDER BY name ASC
LIMIT ? OFFSET ?;

-- name: ListProgressionsByType :many
SELECT id, name, type, parameters, created_at, updated_at
FROM progressions
WHERE type = ?
ORDER BY name ASC
LIMIT ? OFFSET ?;

-- name: UpdateProgression :exec
UPDATE progressions
SET name = ?, type = ?, parameters = ?, updated_at = ?
WHERE id = ?;

-- name: DeleteProgression :exec
DELETE FROM progressions WHERE id = ?;

-- name: CountProgressions :one
SELECT COUNT(*) FROM progressions;

-- name: CountProgressionsByType :one
SELECT COUNT(*) FROM progressions WHERE type = ?;

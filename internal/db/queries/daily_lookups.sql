-- name: GetDailyLookup :one
SELECT id, name, entries, program_id, created_at, updated_at
FROM daily_lookups
WHERE id = ?;

-- name: ListDailyLookupsByNameAsc :many
SELECT id, name, entries, program_id, created_at, updated_at
FROM daily_lookups
ORDER BY name ASC
LIMIT ? OFFSET ?;

-- name: ListDailyLookupsByNameDesc :many
SELECT id, name, entries, program_id, created_at, updated_at
FROM daily_lookups
ORDER BY name DESC
LIMIT ? OFFSET ?;

-- name: ListDailyLookupsByCreatedAtAsc :many
SELECT id, name, entries, program_id, created_at, updated_at
FROM daily_lookups
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: ListDailyLookupsByCreatedAtDesc :many
SELECT id, name, entries, program_id, created_at, updated_at
FROM daily_lookups
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountDailyLookups :one
SELECT COUNT(*) FROM daily_lookups;

-- name: CreateDailyLookup :exec
INSERT INTO daily_lookups (id, name, entries, program_id, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdateDailyLookup :exec
UPDATE daily_lookups
SET name = ?, entries = ?, program_id = ?, updated_at = ?
WHERE id = ?;

-- name: DeleteDailyLookup :exec
DELETE FROM daily_lookups WHERE id = ?;

-- name: DailyLookupIsUsedByPrograms :one
SELECT EXISTS(
    SELECT 1 FROM programs p
    WHERE p.daily_lookup_id = ?
) AS is_used;

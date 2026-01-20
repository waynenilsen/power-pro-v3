-- name: GetWeeklyLookup :one
SELECT id, name, entries, program_id, created_at, updated_at
FROM weekly_lookups
WHERE id = ?;

-- name: ListWeeklyLookupsByNameAsc :many
SELECT id, name, entries, program_id, created_at, updated_at
FROM weekly_lookups
ORDER BY name ASC
LIMIT ? OFFSET ?;

-- name: ListWeeklyLookupsByNameDesc :many
SELECT id, name, entries, program_id, created_at, updated_at
FROM weekly_lookups
ORDER BY name DESC
LIMIT ? OFFSET ?;

-- name: ListWeeklyLookupsByCreatedAtAsc :many
SELECT id, name, entries, program_id, created_at, updated_at
FROM weekly_lookups
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: ListWeeklyLookupsByCreatedAtDesc :many
SELECT id, name, entries, program_id, created_at, updated_at
FROM weekly_lookups
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountWeeklyLookups :one
SELECT COUNT(*) FROM weekly_lookups;

-- name: CreateWeeklyLookup :exec
INSERT INTO weekly_lookups (id, name, entries, program_id, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdateWeeklyLookup :exec
UPDATE weekly_lookups
SET name = ?, entries = ?, program_id = ?, updated_at = ?
WHERE id = ?;

-- name: DeleteWeeklyLookup :exec
DELETE FROM weekly_lookups WHERE id = ?;

-- name: WeeklyLookupIsUsedByPrograms :one
SELECT EXISTS(
    SELECT 1 FROM programs p
    WHERE p.weekly_lookup_id = ?
) AS is_used;

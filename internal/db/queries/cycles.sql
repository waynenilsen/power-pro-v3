-- name: GetCycle :one
SELECT id, name, length_weeks, created_at, updated_at
FROM cycles
WHERE id = ?;

-- name: ListCyclesByNameAsc :many
SELECT id, name, length_weeks, created_at, updated_at
FROM cycles
ORDER BY name ASC
LIMIT ? OFFSET ?;

-- name: ListCyclesByNameDesc :many
SELECT id, name, length_weeks, created_at, updated_at
FROM cycles
ORDER BY name DESC
LIMIT ? OFFSET ?;

-- name: ListCyclesByCreatedAtAsc :many
SELECT id, name, length_weeks, created_at, updated_at
FROM cycles
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: ListCyclesByCreatedAtDesc :many
SELECT id, name, length_weeks, created_at, updated_at
FROM cycles
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListCyclesByLengthWeeksAsc :many
SELECT id, name, length_weeks, created_at, updated_at
FROM cycles
ORDER BY length_weeks ASC
LIMIT ? OFFSET ?;

-- name: ListCyclesByLengthWeeksDesc :many
SELECT id, name, length_weeks, created_at, updated_at
FROM cycles
ORDER BY length_weeks DESC
LIMIT ? OFFSET ?;

-- name: CountCycles :one
SELECT COUNT(*) FROM cycles;

-- name: CreateCycle :exec
INSERT INTO cycles (id, name, length_weeks, created_at, updated_at)
VALUES (?, ?, ?, ?, ?);

-- name: UpdateCycle :exec
UPDATE cycles
SET name = ?, length_weeks = ?, updated_at = ?
WHERE id = ?;

-- name: DeleteCycle :exec
DELETE FROM cycles WHERE id = ?;

-- name: CycleIsUsedByPrograms :one
SELECT EXISTS(
    SELECT 1 FROM programs p
    WHERE p.cycle_id = ?
) AS is_used;

-- name: CountWeeksByCycleID :one
SELECT COUNT(*) FROM weeks WHERE cycle_id = ?;

-- name: ListWeeksByCycleID :many
SELECT id, week_number, variant, cycle_id, created_at, updated_at
FROM weeks
WHERE cycle_id = ?
ORDER BY week_number ASC;

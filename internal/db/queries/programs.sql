-- name: GetProgram :one
SELECT id, name, slug, description, cycle_id, weekly_lookup_id, daily_lookup_id, default_rounding, created_at, updated_at
FROM programs
WHERE id = ?;

-- name: GetProgramBySlug :one
SELECT id, name, slug, description, cycle_id, weekly_lookup_id, daily_lookup_id, default_rounding, created_at, updated_at
FROM programs
WHERE slug = ?;

-- name: ListProgramsByNameAsc :many
SELECT id, name, slug, description, cycle_id, weekly_lookup_id, daily_lookup_id, default_rounding, created_at, updated_at
FROM programs
ORDER BY name ASC
LIMIT ? OFFSET ?;

-- name: ListProgramsByNameDesc :many
SELECT id, name, slug, description, cycle_id, weekly_lookup_id, daily_lookup_id, default_rounding, created_at, updated_at
FROM programs
ORDER BY name DESC
LIMIT ? OFFSET ?;

-- name: ListProgramsByCreatedAtAsc :many
SELECT id, name, slug, description, cycle_id, weekly_lookup_id, daily_lookup_id, default_rounding, created_at, updated_at
FROM programs
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: ListProgramsByCreatedAtDesc :many
SELECT id, name, slug, description, cycle_id, weekly_lookup_id, daily_lookup_id, default_rounding, created_at, updated_at
FROM programs
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountPrograms :one
SELECT COUNT(*) FROM programs;

-- name: CreateProgram :exec
INSERT INTO programs (id, name, slug, description, cycle_id, weekly_lookup_id, daily_lookup_id, default_rounding, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateProgram :exec
UPDATE programs
SET name = ?, slug = ?, description = ?, cycle_id = ?, weekly_lookup_id = ?, daily_lookup_id = ?, default_rounding = ?, updated_at = ?
WHERE id = ?;

-- name: DeleteProgram :exec
DELETE FROM programs WHERE id = ?;

-- name: ProgramSlugExists :one
SELECT EXISTS(
    SELECT 1 FROM programs WHERE slug = ?
) AS slug_exists;

-- name: ProgramSlugExistsExcluding :one
SELECT EXISTS(
    SELECT 1 FROM programs WHERE slug = ? AND id != ?
) AS slug_exists;

-- name: ProgramHasEnrolledUsers :one
SELECT EXISTS(
    SELECT 1 FROM user_program_states ups
    WHERE ups.program_id = ?
) AS has_enrolled;

-- name: CountEnrolledUsers :one
SELECT COUNT(*) FROM user_program_states WHERE program_id = ?;

-- name: GetCycleForProgram :one
SELECT c.id, c.name, c.length_weeks, c.created_at, c.updated_at
FROM cycles c
WHERE c.id = ?;

-- name: GetWeeklyLookupForProgram :one
SELECT id, name, entries, program_id, created_at, updated_at
FROM weekly_lookups
WHERE id = ?;

-- name: GetDailyLookupForProgram :one
SELECT id, name, entries, program_id, created_at, updated_at
FROM daily_lookups
WHERE id = ?;

-- name: GetDay :one
SELECT id, name, slug, metadata, program_id, created_at, updated_at
FROM days
WHERE id = ?;

-- name: GetDayBySlug :one
SELECT id, name, slug, metadata, program_id, created_at, updated_at
FROM days
WHERE slug = ? AND (program_id = ? OR (program_id IS NULL AND ? IS NULL));

-- name: ListDaysByNameAsc :many
SELECT id, name, slug, metadata, program_id, created_at, updated_at
FROM days
ORDER BY name ASC
LIMIT ? OFFSET ?;

-- name: ListDaysByNameDesc :many
SELECT id, name, slug, metadata, program_id, created_at, updated_at
FROM days
ORDER BY name DESC
LIMIT ? OFFSET ?;

-- name: ListDaysByCreatedAtAsc :many
SELECT id, name, slug, metadata, program_id, created_at, updated_at
FROM days
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: ListDaysByCreatedAtDesc :many
SELECT id, name, slug, metadata, program_id, created_at, updated_at
FROM days
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListDaysFilteredByProgramByNameAsc :many
SELECT id, name, slug, metadata, program_id, created_at, updated_at
FROM days
WHERE program_id = ?
ORDER BY name ASC
LIMIT ? OFFSET ?;

-- name: ListDaysFilteredByProgramByNameDesc :many
SELECT id, name, slug, metadata, program_id, created_at, updated_at
FROM days
WHERE program_id = ?
ORDER BY name DESC
LIMIT ? OFFSET ?;

-- name: ListDaysFilteredByProgramByCreatedAtAsc :many
SELECT id, name, slug, metadata, program_id, created_at, updated_at
FROM days
WHERE program_id = ?
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: ListDaysFilteredByProgramByCreatedAtDesc :many
SELECT id, name, slug, metadata, program_id, created_at, updated_at
FROM days
WHERE program_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountDays :one
SELECT COUNT(*) FROM days;

-- name: CountDaysFilteredByProgram :one
SELECT COUNT(*) FROM days WHERE program_id = ?;

-- name: CreateDay :exec
INSERT INTO days (id, name, slug, metadata, program_id, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: UpdateDay :exec
UPDATE days
SET name = ?, slug = ?, metadata = ?, program_id = ?, updated_at = ?
WHERE id = ?;

-- name: DeleteDay :exec
DELETE FROM days WHERE id = ?;

-- name: DaySlugExists :one
SELECT EXISTS(SELECT 1 FROM days WHERE slug = ? AND program_id = ? AND id != ?) AS slug_exists;

-- name: DaySlugExistsForNew :one
SELECT EXISTS(SELECT 1 FROM days WHERE slug = ? AND program_id = ?) AS slug_exists;

-- name: DaySlugExistsNullProgram :one
SELECT EXISTS(SELECT 1 FROM days WHERE slug = ? AND program_id IS NULL AND id != ?) AS slug_exists;

-- name: DaySlugExistsNullProgramForNew :one
SELECT EXISTS(SELECT 1 FROM days WHERE slug = ? AND program_id IS NULL) AS slug_exists;

-- name: DayIsUsedInWeeks :one
SELECT EXISTS(SELECT 1 FROM week_days WHERE day_id = ?) AS is_used;

-- Day Prescriptions queries

-- name: GetDayPrescription :one
SELECT id, day_id, prescription_id, "order", created_at
FROM day_prescriptions
WHERE id = ?;

-- name: ListDayPrescriptions :many
SELECT dp.id, dp.day_id, dp.prescription_id, dp."order", dp.created_at
FROM day_prescriptions dp
WHERE dp.day_id = ?
ORDER BY dp."order" ASC;

-- name: CreateDayPrescription :exec
INSERT INTO day_prescriptions (id, day_id, prescription_id, "order", created_at)
VALUES (?, ?, ?, ?, ?);

-- name: DeleteDayPrescription :exec
DELETE FROM day_prescriptions WHERE id = ?;

-- name: DeleteDayPrescriptionByDayAndPrescription :exec
DELETE FROM day_prescriptions WHERE day_id = ? AND prescription_id = ?;

-- name: GetDayPrescriptionByDayAndPrescription :one
SELECT id, day_id, prescription_id, "order", created_at
FROM day_prescriptions
WHERE day_id = ? AND prescription_id = ?;

-- name: UpdateDayPrescriptionOrder :exec
UPDATE day_prescriptions
SET "order" = ?
WHERE id = ?;

-- name: GetMaxDayPrescriptionOrder :one
SELECT COALESCE(MAX("order"), -1) as max_order
FROM day_prescriptions
WHERE day_id = ?;

-- name: CountDayPrescriptions :one
SELECT COUNT(*) FROM day_prescriptions WHERE day_id = ?;

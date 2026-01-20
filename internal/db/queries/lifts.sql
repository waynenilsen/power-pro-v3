-- name: GetLift :one
SELECT id, name, slug, is_competition_lift, parent_lift_id, created_at, updated_at
FROM lifts
WHERE id = ?;

-- name: GetLiftBySlug :one
SELECT id, name, slug, is_competition_lift, parent_lift_id, created_at, updated_at
FROM lifts
WHERE slug = ?;

-- name: ListLiftsByNameAsc :many
SELECT id, name, slug, is_competition_lift, parent_lift_id, created_at, updated_at
FROM lifts
ORDER BY name ASC
LIMIT ? OFFSET ?;

-- name: ListLiftsByNameDesc :many
SELECT id, name, slug, is_competition_lift, parent_lift_id, created_at, updated_at
FROM lifts
ORDER BY name DESC
LIMIT ? OFFSET ?;

-- name: ListLiftsByCreatedAtAsc :many
SELECT id, name, slug, is_competition_lift, parent_lift_id, created_at, updated_at
FROM lifts
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: ListLiftsByCreatedAtDesc :many
SELECT id, name, slug, is_competition_lift, parent_lift_id, created_at, updated_at
FROM lifts
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListLiftsFilteredByCompetitionByNameAsc :many
SELECT id, name, slug, is_competition_lift, parent_lift_id, created_at, updated_at
FROM lifts
WHERE is_competition_lift = ?
ORDER BY name ASC
LIMIT ? OFFSET ?;

-- name: ListLiftsFilteredByCompetitionByNameDesc :many
SELECT id, name, slug, is_competition_lift, parent_lift_id, created_at, updated_at
FROM lifts
WHERE is_competition_lift = ?
ORDER BY name DESC
LIMIT ? OFFSET ?;

-- name: ListLiftsFilteredByCompetitionByCreatedAtAsc :many
SELECT id, name, slug, is_competition_lift, parent_lift_id, created_at, updated_at
FROM lifts
WHERE is_competition_lift = ?
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: ListLiftsFilteredByCompetitionByCreatedAtDesc :many
SELECT id, name, slug, is_competition_lift, parent_lift_id, created_at, updated_at
FROM lifts
WHERE is_competition_lift = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountLifts :one
SELECT COUNT(*) FROM lifts;

-- name: CountLiftsFilteredByCompetition :one
SELECT COUNT(*) FROM lifts WHERE is_competition_lift = ?;

-- name: CreateLift :exec
INSERT INTO lifts (id, name, slug, is_competition_lift, parent_lift_id, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: UpdateLift :exec
UPDATE lifts
SET name = ?, slug = ?, is_competition_lift = ?, parent_lift_id = ?, updated_at = ?
WHERE id = ?;

-- name: DeleteLift :exec
DELETE FROM lifts WHERE id = ?;

-- name: SlugExists :one
SELECT EXISTS(SELECT 1 FROM lifts WHERE slug = ? AND id != ?) AS slug_exists;

-- name: SlugExistsForNew :one
SELECT EXISTS(SELECT 1 FROM lifts WHERE slug = ?) AS slug_exists;

-- name: LiftHasChildReferences :one
SELECT EXISTS(SELECT 1 FROM lifts WHERE parent_lift_id = ?) AS has_references;

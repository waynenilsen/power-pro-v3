-- name: GetLiftMax :one
SELECT id, user_id, lift_id, type, value, effective_date, created_at, updated_at
FROM lift_maxes
WHERE id = ?;

-- name: ListLiftMaxesByUserByEffectiveDateDesc :many
SELECT id, user_id, lift_id, type, value, effective_date, created_at, updated_at
FROM lift_maxes
WHERE user_id = ?
ORDER BY effective_date DESC
LIMIT ? OFFSET ?;

-- name: ListLiftMaxesByUserByEffectiveDateAsc :many
SELECT id, user_id, lift_id, type, value, effective_date, created_at, updated_at
FROM lift_maxes
WHERE user_id = ?
ORDER BY effective_date ASC
LIMIT ? OFFSET ?;

-- name: ListLiftMaxesByUserFilterLiftByEffectiveDateDesc :many
SELECT id, user_id, lift_id, type, value, effective_date, created_at, updated_at
FROM lift_maxes
WHERE user_id = ? AND lift_id = ?
ORDER BY effective_date DESC
LIMIT ? OFFSET ?;

-- name: ListLiftMaxesByUserFilterLiftByEffectiveDateAsc :many
SELECT id, user_id, lift_id, type, value, effective_date, created_at, updated_at
FROM lift_maxes
WHERE user_id = ? AND lift_id = ?
ORDER BY effective_date ASC
LIMIT ? OFFSET ?;

-- name: ListLiftMaxesByUserFilterTypeByEffectiveDateDesc :many
SELECT id, user_id, lift_id, type, value, effective_date, created_at, updated_at
FROM lift_maxes
WHERE user_id = ? AND type = ?
ORDER BY effective_date DESC
LIMIT ? OFFSET ?;

-- name: ListLiftMaxesByUserFilterTypeByEffectiveDateAsc :many
SELECT id, user_id, lift_id, type, value, effective_date, created_at, updated_at
FROM lift_maxes
WHERE user_id = ? AND type = ?
ORDER BY effective_date ASC
LIMIT ? OFFSET ?;

-- name: ListLiftMaxesByUserFilterLiftAndTypeByEffectiveDateDesc :many
SELECT id, user_id, lift_id, type, value, effective_date, created_at, updated_at
FROM lift_maxes
WHERE user_id = ? AND lift_id = ? AND type = ?
ORDER BY effective_date DESC
LIMIT ? OFFSET ?;

-- name: ListLiftMaxesByUserFilterLiftAndTypeByEffectiveDateAsc :many
SELECT id, user_id, lift_id, type, value, effective_date, created_at, updated_at
FROM lift_maxes
WHERE user_id = ? AND lift_id = ? AND type = ?
ORDER BY effective_date ASC
LIMIT ? OFFSET ?;

-- name: CountLiftMaxesByUser :one
SELECT COUNT(*) FROM lift_maxes WHERE user_id = ?;

-- name: CountLiftMaxesByUserFilterLift :one
SELECT COUNT(*) FROM lift_maxes WHERE user_id = ? AND lift_id = ?;

-- name: CountLiftMaxesByUserFilterType :one
SELECT COUNT(*) FROM lift_maxes WHERE user_id = ? AND type = ?;

-- name: CountLiftMaxesByUserFilterLiftAndType :one
SELECT COUNT(*) FROM lift_maxes WHERE user_id = ? AND lift_id = ? AND type = ?;

-- name: CreateLiftMax :exec
INSERT INTO lift_maxes (id, user_id, lift_id, type, value, effective_date, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateLiftMax :exec
UPDATE lift_maxes
SET value = ?, effective_date = ?, updated_at = ?
WHERE id = ?;

-- name: DeleteLiftMax :exec
DELETE FROM lift_maxes WHERE id = ?;

-- name: GetCurrentOneRM :one
SELECT id, user_id, lift_id, type, value, effective_date, created_at, updated_at
FROM lift_maxes
WHERE user_id = ? AND lift_id = ? AND type = 'ONE_RM'
ORDER BY effective_date DESC
LIMIT 1;

-- name: GetCurrentMax :one
SELECT id, user_id, lift_id, type, value, effective_date, created_at, updated_at
FROM lift_maxes
WHERE user_id = ? AND lift_id = ? AND type = ?
ORDER BY effective_date DESC
LIMIT 1;

-- name: UniqueConstraintExists :one
SELECT EXISTS(
    SELECT 1 FROM lift_maxes
    WHERE user_id = ? AND lift_id = ? AND type = ? AND effective_date = ?
) AS constraint_exists;

-- name: UniqueConstraintExistsExcluding :one
SELECT EXISTS(
    SELECT 1 FROM lift_maxes
    WHERE user_id = ? AND lift_id = ? AND type = ? AND effective_date = ? AND id != ?
) AS constraint_exists;

-- name: LiftHasMaxReferences :one
SELECT EXISTS(SELECT 1 FROM lift_maxes WHERE lift_id = ?) AS has_references;

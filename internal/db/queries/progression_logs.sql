-- name: CreateProgressionLog :exec
INSERT INTO progression_logs (id, user_id, progression_id, lift_id, previous_value, new_value, delta, trigger_type, trigger_context, applied_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetProgressionLog :one
SELECT id, user_id, progression_id, lift_id, previous_value, new_value, delta, trigger_type, trigger_context, applied_at
FROM progression_logs
WHERE id = ?;

-- name: CheckIdempotency :one
SELECT EXISTS(
    SELECT 1 FROM progression_logs
    WHERE user_id = ? AND progression_id = ? AND lift_id = ? AND trigger_type = ? AND applied_at = ?
) AS already_applied;

-- name: ListProgressionLogsByUser :many
SELECT id, user_id, progression_id, lift_id, previous_value, new_value, delta, trigger_type, trigger_context, applied_at
FROM progression_logs
WHERE user_id = ?
ORDER BY applied_at DESC
LIMIT ? OFFSET ?;

-- name: ListProgressionLogsByUserAndLift :many
SELECT id, user_id, progression_id, lift_id, previous_value, new_value, delta, trigger_type, trigger_context, applied_at
FROM progression_logs
WHERE user_id = ? AND lift_id = ?
ORDER BY applied_at DESC
LIMIT ? OFFSET ?;

-- name: CountProgressionLogsByUser :one
SELECT COUNT(*) FROM progression_logs WHERE user_id = ?;

-- name: CountProgressionLogsByUserAndLift :one
SELECT COUNT(*) FROM progression_logs WHERE user_id = ? AND lift_id = ?;

-- name: DeleteProgressionLog :exec
DELETE FROM progression_logs WHERE id = ?;

-- name: ListProgressionHistoryBase :many
SELECT
    pl.id,
    pl.user_id,
    pl.progression_id,
    pl.lift_id,
    pl.previous_value,
    pl.new_value,
    pl.delta,
    pl.trigger_type,
    pl.trigger_context,
    pl.applied_at,
    p.name AS progression_name,
    p.type AS progression_type,
    l.name AS lift_name
FROM progression_logs pl
JOIN progressions p ON pl.progression_id = p.id
JOIN lifts l ON pl.lift_id = l.id
WHERE pl.user_id = ?
ORDER BY pl.applied_at DESC
LIMIT ? OFFSET ?;

-- name: CountProgressionHistoryBase :one
SELECT COUNT(*) FROM progression_logs pl
WHERE pl.user_id = ?;

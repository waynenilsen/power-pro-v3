-- User Progression States Queries
-- Used to track per-user, per-lift, per-progression state (e.g., current stage)

-- name: CreateUserProgressionState :exec
INSERT INTO user_progression_states (id, user_id, lift_id, progression_id, current_stage, state_data, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'));

-- name: GetUserProgressionState :one
SELECT id, user_id, lift_id, progression_id, current_stage, state_data, created_at, updated_at
FROM user_progression_states
WHERE user_id = ? AND lift_id = ? AND progression_id = ?;

-- name: UpdateUserProgressionStateStage :exec
UPDATE user_progression_states
SET current_stage = ?, updated_at = datetime('now')
WHERE user_id = ? AND lift_id = ? AND progression_id = ?;

-- name: UpsertUserProgressionState :exec
INSERT INTO user_progression_states (id, user_id, lift_id, progression_id, current_stage, state_data, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))
ON CONFLICT(user_id, lift_id, progression_id) DO UPDATE SET
    current_stage = excluded.current_stage,
    state_data = excluded.state_data,
    updated_at = datetime('now');

-- name: DeleteUserProgressionState :exec
DELETE FROM user_progression_states
WHERE user_id = ? AND lift_id = ? AND progression_id = ?;

-- name: ListUserProgressionStatesByUser :many
SELECT id, user_id, lift_id, progression_id, current_stage, state_data, created_at, updated_at
FROM user_progression_states
WHERE user_id = ?
ORDER BY updated_at DESC;

-- name: ListUserProgressionStatesByProgression :many
SELECT id, user_id, lift_id, progression_id, current_stage, state_data, created_at, updated_at
FROM user_progression_states
WHERE progression_id = ?
ORDER BY updated_at DESC;

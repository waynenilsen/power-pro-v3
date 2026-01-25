-- name: CreateWorkoutSession :exec
INSERT INTO workout_sessions (id, user_program_state_id, week_number, day_index, status, started_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetWorkoutSessionByID :one
SELECT id, user_program_state_id, week_number, day_index, status, started_at, finished_at, created_at, updated_at
FROM workout_sessions
WHERE id = ?;

-- name: GetActiveWorkoutSession :one
SELECT id, user_program_state_id, week_number, day_index, status, started_at, finished_at, created_at, updated_at
FROM workout_sessions
WHERE user_program_state_id = ? AND status = 'IN_PROGRESS'
LIMIT 1;

-- name: GetWorkoutSessionsByState :many
SELECT id, user_program_state_id, week_number, day_index, status, started_at, finished_at, created_at, updated_at
FROM workout_sessions
WHERE user_program_state_id = ?
ORDER BY created_at DESC;

-- name: UpdateWorkoutSessionStatus :exec
UPDATE workout_sessions
SET status = ?, updated_at = ?
WHERE id = ?;

-- name: CompleteWorkoutSession :exec
UPDATE workout_sessions
SET status = 'COMPLETED', finished_at = ?, updated_at = ?
WHERE id = ?;

-- name: AbandonWorkoutSession :exec
UPDATE workout_sessions
SET status = 'ABANDONED', finished_at = ?, updated_at = ?
WHERE id = ?;

-- name: DeleteWorkoutSession :exec
DELETE FROM workout_sessions WHERE id = ?;

-- name: CountWorkoutSessionsByUserID :one
SELECT COUNT(*) as count
FROM workout_sessions ws
JOIN user_program_states ups ON ws.user_program_state_id = ups.id
WHERE ups.user_id = ?;

-- name: CountWorkoutSessionsByUserIDWithStatus :one
SELECT COUNT(*) as count
FROM workout_sessions ws
JOIN user_program_states ups ON ws.user_program_state_id = ups.id
WHERE ups.user_id = ? AND ws.status = ?;

-- name: GetWorkoutSessionsByUserID :many
SELECT ws.id, ws.user_program_state_id, ws.week_number, ws.day_index, ws.status, ws.started_at, ws.finished_at, ws.created_at, ws.updated_at
FROM workout_sessions ws
JOIN user_program_states ups ON ws.user_program_state_id = ups.id
WHERE ups.user_id = ?
ORDER BY ws.created_at DESC
LIMIT ? OFFSET ?;

-- name: GetWorkoutSessionsByUserIDWithStatus :many
SELECT ws.id, ws.user_program_state_id, ws.week_number, ws.day_index, ws.status, ws.started_at, ws.finished_at, ws.created_at, ws.updated_at
FROM workout_sessions ws
JOIN user_program_states ups ON ws.user_program_state_id = ups.id
WHERE ups.user_id = ? AND ws.status = ?
ORDER BY ws.created_at DESC
LIMIT ? OFFSET ?;

-- name: GetActiveWorkoutSessionByUserID :one
SELECT ws.id, ws.user_program_state_id, ws.week_number, ws.day_index, ws.status, ws.started_at, ws.finished_at, ws.created_at, ws.updated_at
FROM workout_sessions ws
JOIN user_program_states ups ON ws.user_program_state_id = ups.id
WHERE ups.user_id = ? AND ws.status = 'IN_PROGRESS'
LIMIT 1;

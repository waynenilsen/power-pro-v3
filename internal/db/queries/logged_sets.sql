-- name: CreateLoggedSet :exec
INSERT INTO logged_sets (id, user_id, session_id, prescription_id, lift_id, set_number, weight, target_reps, reps_performed, is_amrap, rpe, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetLoggedSet :one
SELECT id, user_id, session_id, prescription_id, lift_id, set_number, weight, target_reps, reps_performed, is_amrap, rpe, created_at
FROM logged_sets
WHERE id = ?;

-- name: ListLoggedSetsBySession :many
SELECT id, user_id, session_id, prescription_id, lift_id, set_number, weight, target_reps, reps_performed, is_amrap, rpe, created_at
FROM logged_sets
WHERE session_id = ?
ORDER BY created_at ASC, set_number ASC;

-- name: ListLoggedSetsByUser :many
SELECT id, user_id, session_id, prescription_id, lift_id, set_number, weight, target_reps, reps_performed, is_amrap, rpe, created_at
FROM logged_sets
WHERE user_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountLoggedSetsByUser :one
SELECT COUNT(*) FROM logged_sets WHERE user_id = ?;

-- name: GetLatestAMRAPForLift :one
SELECT id, user_id, session_id, prescription_id, lift_id, set_number, weight, target_reps, reps_performed, is_amrap, rpe, created_at
FROM logged_sets
WHERE user_id = ? AND lift_id = ? AND is_amrap = TRUE
ORDER BY created_at DESC
LIMIT 1;

-- name: DeleteLoggedSet :exec
DELETE FROM logged_sets WHERE id = ?;

-- name: DeleteLoggedSetsBySession :exec
DELETE FROM logged_sets WHERE session_id = ?;

-- name: ListLoggedSetsBySessionAndPrescription :many
SELECT id, user_id, session_id, prescription_id, lift_id, set_number, weight, target_reps, reps_performed, is_amrap, rpe, created_at
FROM logged_sets
WHERE session_id = ? AND prescription_id = ?
ORDER BY set_number ASC;

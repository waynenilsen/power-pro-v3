-- name: CreateFailureCounter :exec
INSERT INTO failure_counters (id, user_id, lift_id, progression_id, consecutive_failures, last_failure_at, last_success_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetFailureCounter :one
SELECT id, user_id, lift_id, progression_id, consecutive_failures, last_failure_at, last_success_at, created_at, updated_at
FROM failure_counters
WHERE id = ?;

-- name: GetFailureCounterByKey :one
SELECT id, user_id, lift_id, progression_id, consecutive_failures, last_failure_at, last_success_at, created_at, updated_at
FROM failure_counters
WHERE user_id = ? AND lift_id = ? AND progression_id = ?;

-- name: ListFailureCountersByUser :many
SELECT id, user_id, lift_id, progression_id, consecutive_failures, last_failure_at, last_success_at, created_at, updated_at
FROM failure_counters
WHERE user_id = ?
ORDER BY updated_at DESC;

-- name: ListFailureCountersByUserAndLift :many
SELECT id, user_id, lift_id, progression_id, consecutive_failures, last_failure_at, last_success_at, created_at, updated_at
FROM failure_counters
WHERE user_id = ? AND lift_id = ?
ORDER BY updated_at DESC;

-- name: ListFailureCountersByProgression :many
SELECT id, user_id, lift_id, progression_id, consecutive_failures, last_failure_at, last_success_at, created_at, updated_at
FROM failure_counters
WHERE progression_id = ?
ORDER BY updated_at DESC;

-- name: UpdateFailureCounter :exec
UPDATE failure_counters
SET consecutive_failures = ?,
    last_failure_at = ?,
    last_success_at = ?,
    updated_at = ?
WHERE id = ?;

-- name: IncrementFailureCounter :exec
UPDATE failure_counters
SET consecutive_failures = consecutive_failures + 1,
    last_failure_at = ?,
    updated_at = ?
WHERE user_id = ? AND lift_id = ? AND progression_id = ?;

-- name: ResetFailureCounter :exec
UPDATE failure_counters
SET consecutive_failures = 0,
    last_success_at = ?,
    updated_at = ?
WHERE user_id = ? AND lift_id = ? AND progression_id = ?;

-- name: DeleteFailureCounter :exec
DELETE FROM failure_counters WHERE id = ?;

-- name: DeleteFailureCounterByKey :exec
DELETE FROM failure_counters WHERE user_id = ? AND lift_id = ? AND progression_id = ?;

-- name: CountFailureCountersByUser :one
SELECT COUNT(*) FROM failure_counters WHERE user_id = ?;

-- name: UpsertFailureCounterOnFailure :exec
INSERT INTO failure_counters (id, user_id, lift_id, progression_id, consecutive_failures, last_failure_at, last_success_at, created_at, updated_at)
VALUES (?, ?, ?, ?, 1, ?, NULL, ?, ?)
ON CONFLICT(user_id, lift_id, progression_id) DO UPDATE SET
    consecutive_failures = consecutive_failures + 1,
    last_failure_at = excluded.last_failure_at,
    updated_at = excluded.updated_at;

-- name: UpsertFailureCounterOnSuccess :exec
INSERT INTO failure_counters (id, user_id, lift_id, progression_id, consecutive_failures, last_failure_at, last_success_at, created_at, updated_at)
VALUES (?, ?, ?, ?, 0, NULL, ?, ?, ?)
ON CONFLICT(user_id, lift_id, progression_id) DO UPDATE SET
    consecutive_failures = 0,
    last_success_at = excluded.last_success_at,
    updated_at = excluded.updated_at;

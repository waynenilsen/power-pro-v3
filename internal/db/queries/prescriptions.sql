-- name: GetPrescription :one
SELECT id, lift_id, load_strategy, set_scheme, "order", notes, rest_seconds, created_at, updated_at
FROM prescriptions
WHERE id = ?;

-- name: ListPrescriptionsByOrderAsc :many
SELECT id, lift_id, load_strategy, set_scheme, "order", notes, rest_seconds, created_at, updated_at
FROM prescriptions
ORDER BY "order" ASC
LIMIT ? OFFSET ?;

-- name: ListPrescriptionsByOrderDesc :many
SELECT id, lift_id, load_strategy, set_scheme, "order", notes, rest_seconds, created_at, updated_at
FROM prescriptions
ORDER BY "order" DESC
LIMIT ? OFFSET ?;

-- name: ListPrescriptionsByCreatedAtAsc :many
SELECT id, lift_id, load_strategy, set_scheme, "order", notes, rest_seconds, created_at, updated_at
FROM prescriptions
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: ListPrescriptionsByCreatedAtDesc :many
SELECT id, lift_id, load_strategy, set_scheme, "order", notes, rest_seconds, created_at, updated_at
FROM prescriptions
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListPrescriptionsFilterLiftByOrderAsc :many
SELECT id, lift_id, load_strategy, set_scheme, "order", notes, rest_seconds, created_at, updated_at
FROM prescriptions
WHERE lift_id = ?
ORDER BY "order" ASC
LIMIT ? OFFSET ?;

-- name: ListPrescriptionsFilterLiftByOrderDesc :many
SELECT id, lift_id, load_strategy, set_scheme, "order", notes, rest_seconds, created_at, updated_at
FROM prescriptions
WHERE lift_id = ?
ORDER BY "order" DESC
LIMIT ? OFFSET ?;

-- name: ListPrescriptionsFilterLiftByCreatedAtAsc :many
SELECT id, lift_id, load_strategy, set_scheme, "order", notes, rest_seconds, created_at, updated_at
FROM prescriptions
WHERE lift_id = ?
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: ListPrescriptionsFilterLiftByCreatedAtDesc :many
SELECT id, lift_id, load_strategy, set_scheme, "order", notes, rest_seconds, created_at, updated_at
FROM prescriptions
WHERE lift_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountPrescriptions :one
SELECT COUNT(*) FROM prescriptions;

-- name: CountPrescriptionsFilterLift :one
SELECT COUNT(*) FROM prescriptions WHERE lift_id = ?;

-- name: CreatePrescription :exec
INSERT INTO prescriptions (id, lift_id, load_strategy, set_scheme, "order", notes, rest_seconds, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdatePrescription :exec
UPDATE prescriptions
SET lift_id = ?, load_strategy = ?, set_scheme = ?, "order" = ?, notes = ?, rest_seconds = ?, updated_at = ?
WHERE id = ?;

-- name: DeletePrescription :exec
DELETE FROM prescriptions WHERE id = ?;

-- name: LiftHasPrescriptionReferences :one
SELECT EXISTS(SELECT 1 FROM prescriptions WHERE lift_id = ?) AS has_references;

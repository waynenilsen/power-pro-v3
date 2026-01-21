-- name: CreateProgramProgression :exec
INSERT INTO program_progressions (id, program_id, progression_id, lift_id, priority, enabled, override_increment, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetProgramProgression :one
SELECT id, program_id, progression_id, lift_id, priority, enabled, override_increment, created_at, updated_at
FROM program_progressions
WHERE id = ?;

-- name: ListProgramProgressionsByProgram :many
SELECT id, program_id, progression_id, lift_id, priority, enabled, override_increment, created_at, updated_at
FROM program_progressions
WHERE program_id = ?
ORDER BY priority ASC;

-- name: ListEnabledProgramProgressionsByProgram :many
SELECT id, program_id, progression_id, lift_id, priority, enabled, override_increment, created_at, updated_at
FROM program_progressions
WHERE program_id = ? AND enabled = 1
ORDER BY priority ASC;

-- name: ListProgramProgressionsByProgramAndLift :many
SELECT id, program_id, progression_id, lift_id, priority, enabled, override_increment, created_at, updated_at
FROM program_progressions
WHERE program_id = ? AND (lift_id = ? OR lift_id IS NULL) AND enabled = 1
ORDER BY priority ASC;

-- name: UpdateProgramProgression :exec
UPDATE program_progressions
SET priority = ?, enabled = ?, override_increment = ?, updated_at = ?
WHERE id = ?;

-- name: DeleteProgramProgression :exec
DELETE FROM program_progressions WHERE id = ?;

-- name: DeleteProgramProgressionsByProgram :exec
DELETE FROM program_progressions WHERE program_id = ?;

-- name: CountProgramProgressionsByProgram :one
SELECT COUNT(*) FROM program_progressions WHERE program_id = ?;

-- name: CountProgramProgressionsByProgression :one
SELECT COUNT(*) FROM program_progressions WHERE progression_id = ?;

-- name: GetProgramProgressionByProgramProgressionLift :one
SELECT id, program_id, progression_id, lift_id, priority, enabled, override_increment, created_at, updated_at
FROM program_progressions
WHERE program_id = ? AND progression_id = ? AND COALESCE(lift_id, '00000000-0000-0000-0000-000000000000') = COALESCE(?, '00000000-0000-0000-0000-000000000000');

-- name: ListProgramProgressionsWithDetailsByProgram :many
SELECT
    pp.id,
    pp.program_id,
    pp.progression_id,
    pp.lift_id,
    pp.priority,
    pp.enabled,
    pp.override_increment,
    pp.created_at,
    pp.updated_at,
    p.name as progression_name,
    p.type as progression_type,
    p.parameters as progression_parameters
FROM program_progressions pp
JOIN progressions p ON pp.progression_id = p.id
WHERE pp.program_id = ?
ORDER BY pp.priority ASC;

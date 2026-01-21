-- name: GetUserProgramStateByUserID :one
SELECT id, user_id, program_id, current_week, current_cycle_iteration, current_day_index, enrolled_at, updated_at
FROM user_program_states
WHERE user_id = ?;

-- name: GetUserProgramStateByID :one
SELECT id, user_id, program_id, current_week, current_cycle_iteration, current_day_index, enrolled_at, updated_at
FROM user_program_states
WHERE id = ?;

-- name: CreateUserProgramState :exec
INSERT INTO user_program_states (id, user_id, program_id, current_week, current_cycle_iteration, current_day_index, enrolled_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateUserProgramState :exec
UPDATE user_program_states
SET program_id = ?, current_week = ?, current_cycle_iteration = ?, current_day_index = ?, updated_at = ?
WHERE user_id = ?;

-- name: DeleteUserProgramStateByUserID :exec
DELETE FROM user_program_states WHERE user_id = ?;

-- name: UserIsEnrolled :one
SELECT EXISTS(
    SELECT 1 FROM user_program_states WHERE user_id = ?
) AS is_enrolled;

-- name: GetEnrollmentWithProgram :one
SELECT
    ups.id,
    ups.user_id,
    ups.program_id,
    ups.current_week,
    ups.current_cycle_iteration,
    ups.current_day_index,
    ups.enrolled_at,
    ups.updated_at,
    p.name AS program_name,
    p.slug AS program_slug,
    p.description AS program_description,
    c.length_weeks AS cycle_length_weeks
FROM user_program_states ups
JOIN programs p ON ups.program_id = p.id
JOIN cycles c ON p.cycle_id = c.id
WHERE ups.user_id = ?;

-- name: GetStateAdvancementContext :one
SELECT
    ups.id,
    ups.user_id,
    ups.program_id,
    ups.current_week,
    ups.current_cycle_iteration,
    ups.current_day_index,
    ups.enrolled_at,
    ups.updated_at,
    c.id AS cycle_id,
    c.length_weeks AS cycle_length_weeks,
    (
        SELECT COUNT(*)
        FROM week_days wd
        JOIN weeks w ON wd.week_id = w.id
        WHERE w.cycle_id = c.id AND w.week_number = ups.current_week
    ) AS days_in_current_week
FROM user_program_states ups
JOIN programs p ON ups.program_id = p.id
JOIN cycles c ON p.cycle_id = c.id
WHERE ups.user_id = ?;

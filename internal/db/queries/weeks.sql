-- name: GetWeek :one
SELECT id, week_number, variant, cycle_id, created_at, updated_at
FROM weeks
WHERE id = ?;

-- name: ListWeeksByWeekNumberAsc :many
SELECT id, week_number, variant, cycle_id, created_at, updated_at
FROM weeks
ORDER BY week_number ASC
LIMIT ? OFFSET ?;

-- name: ListWeeksByWeekNumberDesc :many
SELECT id, week_number, variant, cycle_id, created_at, updated_at
FROM weeks
ORDER BY week_number DESC
LIMIT ? OFFSET ?;

-- name: ListWeeksByCreatedAtAsc :many
SELECT id, week_number, variant, cycle_id, created_at, updated_at
FROM weeks
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: ListWeeksByCreatedAtDesc :many
SELECT id, week_number, variant, cycle_id, created_at, updated_at
FROM weeks
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListWeeksFilteredByCycleByWeekNumberAsc :many
SELECT id, week_number, variant, cycle_id, created_at, updated_at
FROM weeks
WHERE cycle_id = ?
ORDER BY week_number ASC
LIMIT ? OFFSET ?;

-- name: ListWeeksFilteredByCycleByWeekNumberDesc :many
SELECT id, week_number, variant, cycle_id, created_at, updated_at
FROM weeks
WHERE cycle_id = ?
ORDER BY week_number DESC
LIMIT ? OFFSET ?;

-- name: ListWeeksFilteredByCycleByCreatedAtAsc :many
SELECT id, week_number, variant, cycle_id, created_at, updated_at
FROM weeks
WHERE cycle_id = ?
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: ListWeeksFilteredByCycleByCreatedAtDesc :many
SELECT id, week_number, variant, cycle_id, created_at, updated_at
FROM weeks
WHERE cycle_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountWeeks :one
SELECT COUNT(*) FROM weeks;

-- name: CountWeeksFilteredByCycle :one
SELECT COUNT(*) FROM weeks WHERE cycle_id = ?;

-- name: CreateWeek :exec
INSERT INTO weeks (id, week_number, variant, cycle_id, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdateWeek :exec
UPDATE weeks
SET week_number = ?, variant = ?, cycle_id = ?, updated_at = ?
WHERE id = ?;

-- name: DeleteWeek :exec
DELETE FROM weeks WHERE id = ?;

-- name: WeekNumberExistsInCycle :one
SELECT EXISTS(SELECT 1 FROM weeks WHERE cycle_id = ? AND week_number = ? AND id != ?) AS week_number_exists;

-- name: WeekNumberExistsInCycleForNew :one
SELECT EXISTS(SELECT 1 FROM weeks WHERE cycle_id = ? AND week_number = ?) AS week_number_exists;

-- name: WeekIsUsedInActiveCycle :one
SELECT EXISTS(
    SELECT 1 FROM user_program_states ups
    JOIN programs p ON ups.program_id = p.id
    JOIN cycles c ON p.cycle_id = c.id
    JOIN weeks w ON w.cycle_id = c.id
    WHERE w.id = ?
) AS is_used;

-- Week Days queries

-- name: GetWeekDay :one
SELECT id, week_id, day_id, day_of_week, created_at
FROM week_days
WHERE id = ?;

-- name: ListWeekDays :many
SELECT wd.id, wd.week_id, wd.day_id, wd.day_of_week, wd.created_at
FROM week_days wd
WHERE wd.week_id = ?
ORDER BY
    CASE wd.day_of_week
        WHEN 'MONDAY' THEN 1
        WHEN 'TUESDAY' THEN 2
        WHEN 'WEDNESDAY' THEN 3
        WHEN 'THURSDAY' THEN 4
        WHEN 'FRIDAY' THEN 5
        WHEN 'SATURDAY' THEN 6
        WHEN 'SUNDAY' THEN 7
    END ASC;

-- name: CreateWeekDay :exec
INSERT INTO week_days (id, week_id, day_id, day_of_week, created_at)
VALUES (?, ?, ?, ?, ?);

-- name: DeleteWeekDay :exec
DELETE FROM week_days WHERE id = ?;

-- name: DeleteWeekDayByWeekAndDay :exec
DELETE FROM week_days WHERE week_id = ? AND day_id = ? AND day_of_week = ?;

-- name: GetWeekDayByWeekAndDayAndDayOfWeek :one
SELECT id, week_id, day_id, day_of_week, created_at
FROM week_days
WHERE week_id = ? AND day_id = ? AND day_of_week = ?;

-- name: CountWeekDays :one
SELECT COUNT(*) FROM week_days WHERE week_id = ?;

-- name: GetCycleByID :one
SELECT id, name, length_weeks, created_at, updated_at
FROM cycles
WHERE id = ?;

-- Workout Generation Queries
-- These queries support the workout generation API endpoint.

-- name: GetWeekByNumberAndCycle :one
SELECT id, week_number, variant, cycle_id, created_at, updated_at
FROM weeks
WHERE cycle_id = ? AND week_number = ?;

-- name: GetDaysForWeek :many
SELECT d.id, d.name, d.slug, d.metadata, d.program_id, d.created_at, d.updated_at,
       wd.day_of_week
FROM days d
JOIN week_days wd ON d.id = wd.day_id
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

-- name: GetDayBySlugAndWeek :one
SELECT d.id, d.name, d.slug, d.metadata, d.program_id, d.created_at, d.updated_at
FROM days d
JOIN week_days wd ON d.id = wd.day_id
WHERE d.slug = ? AND wd.week_id = ?;

-- name: GetPrescriptionsForDay :many
SELECT p.id, p.lift_id, p.load_strategy, p.set_scheme, p."order", p.notes, p.rest_seconds, p.created_at, p.updated_at
FROM prescriptions p
JOIN day_prescriptions dp ON p.id = dp.prescription_id
WHERE dp.day_id = ?
ORDER BY dp."order" ASC;

-- name: GetProgramWithCycle :one
SELECT
    p.id AS program_id,
    p.name AS program_name,
    p.slug AS program_slug,
    p.description AS program_description,
    p.cycle_id,
    p.weekly_lookup_id,
    p.daily_lookup_id,
    p.default_rounding,
    c.length_weeks AS cycle_length_weeks
FROM programs p
JOIN cycles c ON p.cycle_id = c.id
WHERE p.id = ?;

-- name: GetEnrollmentForWorkout :one
SELECT
    ups.id,
    ups.user_id,
    ups.program_id,
    ups.current_week,
    ups.current_cycle_iteration,
    ups.current_day_index,
    p.name AS program_name,
    p.slug AS program_slug,
    p.cycle_id,
    p.weekly_lookup_id,
    p.daily_lookup_id,
    p.default_rounding,
    c.length_weeks AS cycle_length_weeks
FROM user_program_states ups
JOIN programs p ON ups.program_id = p.id
JOIN cycles c ON p.cycle_id = c.id
WHERE ups.user_id = ?;

-- name: GetDayByIndexInWeek :one
SELECT d.id, d.name, d.slug, d.metadata, d.program_id, d.created_at, d.updated_at
FROM days d
JOIN week_days wd ON d.id = wd.day_id
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
    END ASC
LIMIT 1 OFFSET ?;

-- name: CountDaysInWeek :one
SELECT COUNT(*) FROM week_days WHERE week_id = ?;

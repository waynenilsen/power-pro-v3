-- Dashboard aggregation queries

-- name: GetRecentCompletedWorkouts :many
-- Get recent completed workouts for a user with day name and sets completed
-- day_index is used as an offset into the ordered days for the week
SELECT
    ws.id,
    ws.finished_at,
    (
        SELECT d.name
        FROM days d
        JOIN week_days wd ON d.id = wd.day_id
        JOIN weeks w ON wd.week_id = w.id
        JOIN cycles c ON w.cycle_id = c.id
        JOIN programs p ON p.cycle_id = c.id
        WHERE p.id = ups.program_id AND w.week_number = ws.week_number
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
        LIMIT 1 OFFSET ws.day_index
    ) AS day_name,
    (SELECT COUNT(*) FROM logged_sets ls WHERE ls.session_id = ws.id) AS sets_completed
FROM workout_sessions ws
JOIN user_program_states ups ON ws.user_program_state_id = ups.id
WHERE ups.user_id = ? AND ws.status = 'COMPLETED' AND ws.finished_at IS NOT NULL
ORDER BY ws.finished_at DESC
LIMIT ?;

-- name: GetCurrentMaxesByUser :many
-- Get the most recent max for each lift a user has recorded
-- For each lift, find the record with the maximum effective_date
SELECT
    lm.id,
    lm.lift_id,
    l.name AS lift_name,
    lm.type,
    lm.value,
    lm.effective_date
FROM lift_maxes lm
JOIN lifts l ON lm.lift_id = l.id
WHERE lm.user_id = ?
    AND NOT EXISTS (
        SELECT 1
        FROM lift_maxes lm2
        WHERE lm2.user_id = lm.user_id
        AND lm2.lift_id = lm.lift_id
        AND lm2.effective_date > lm.effective_date
    )
ORDER BY l.name ASC;

-- name: GetDayForWeekPosition :one
-- Get the day at a specific position in a week for a program
-- Uses day_index as an offset into the ordered days by day_of_week
SELECT d.id, d.name, d.slug
FROM days d
JOIN week_days wd ON d.id = wd.day_id
JOIN weeks w ON wd.week_id = w.id
JOIN cycles c ON w.cycle_id = c.id
JOIN programs p ON p.cycle_id = c.id
WHERE p.id = ? AND w.week_number = ?
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

-- name: GetDayExerciseAndSetCounts :one
-- Count distinct exercises and estimate total sets for a day
SELECT
    COUNT(DISTINCT pr.lift_id) AS exercise_count,
    COALESCE(SUM(
        CASE
            WHEN pr.scheme_type IN ('fixed', 'greyskull') THEN 3
            WHEN pr.scheme_type = 'ramp' THEN 5
            WHEN pr.scheme_type = 'fatigue_drop' THEN 3
            WHEN pr.scheme_type = 'mrs' THEN 4
            WHEN pr.scheme_type = 'total_reps' THEN 5
            WHEN pr.scheme_type = 'amrap' THEN 1
            ELSE 3
        END
    ), 0) AS total_sets
FROM day_prescriptions dp
JOIN prescriptions pr ON dp.prescription_id = pr.id
WHERE dp.day_id = ?;

-- name: CountLoggedSetsBySession :one
-- Count logged sets for a session
SELECT COUNT(*) FROM logged_sets WHERE session_id = ?;

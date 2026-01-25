-- name: GetProgram :one
SELECT id, name, slug, description, cycle_id, weekly_lookup_id, daily_lookup_id, default_rounding, difficulty, days_per_week, focus, has_amrap, created_at, updated_at
FROM programs
WHERE id = ?;

-- name: GetProgramBySlug :one
SELECT id, name, slug, description, cycle_id, weekly_lookup_id, daily_lookup_id, default_rounding, difficulty, days_per_week, focus, has_amrap, created_at, updated_at
FROM programs
WHERE slug = ?;

-- name: ListProgramsFilteredByNameAsc :many
SELECT id, name, slug, description, cycle_id, weekly_lookup_id, daily_lookup_id, default_rounding, difficulty, days_per_week, focus, has_amrap, created_at, updated_at
FROM programs
WHERE (sqlc.narg('difficulty') IS NULL OR difficulty = sqlc.narg('difficulty'))
  AND (sqlc.narg('days_per_week') IS NULL OR days_per_week = sqlc.narg('days_per_week'))
  AND (sqlc.narg('focus') IS NULL OR focus = sqlc.narg('focus'))
  AND (sqlc.narg('has_amrap') IS NULL OR has_amrap = sqlc.narg('has_amrap'))
  AND (sqlc.narg('search') IS NULL OR name LIKE '%' || sqlc.narg('search') || '%' COLLATE NOCASE)
ORDER BY name ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ListProgramsFilteredByNameDesc :many
SELECT id, name, slug, description, cycle_id, weekly_lookup_id, daily_lookup_id, default_rounding, difficulty, days_per_week, focus, has_amrap, created_at, updated_at
FROM programs
WHERE (sqlc.narg('difficulty') IS NULL OR difficulty = sqlc.narg('difficulty'))
  AND (sqlc.narg('days_per_week') IS NULL OR days_per_week = sqlc.narg('days_per_week'))
  AND (sqlc.narg('focus') IS NULL OR focus = sqlc.narg('focus'))
  AND (sqlc.narg('has_amrap') IS NULL OR has_amrap = sqlc.narg('has_amrap'))
  AND (sqlc.narg('search') IS NULL OR name LIKE '%' || sqlc.narg('search') || '%' COLLATE NOCASE)
ORDER BY name DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ListProgramsFilteredByCreatedAtAsc :many
SELECT id, name, slug, description, cycle_id, weekly_lookup_id, daily_lookup_id, default_rounding, difficulty, days_per_week, focus, has_amrap, created_at, updated_at
FROM programs
WHERE (sqlc.narg('difficulty') IS NULL OR difficulty = sqlc.narg('difficulty'))
  AND (sqlc.narg('days_per_week') IS NULL OR days_per_week = sqlc.narg('days_per_week'))
  AND (sqlc.narg('focus') IS NULL OR focus = sqlc.narg('focus'))
  AND (sqlc.narg('has_amrap') IS NULL OR has_amrap = sqlc.narg('has_amrap'))
  AND (sqlc.narg('search') IS NULL OR name LIKE '%' || sqlc.narg('search') || '%' COLLATE NOCASE)
ORDER BY created_at ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ListProgramsFilteredByCreatedAtDesc :many
SELECT id, name, slug, description, cycle_id, weekly_lookup_id, daily_lookup_id, default_rounding, difficulty, days_per_week, focus, has_amrap, created_at, updated_at
FROM programs
WHERE (sqlc.narg('difficulty') IS NULL OR difficulty = sqlc.narg('difficulty'))
  AND (sqlc.narg('days_per_week') IS NULL OR days_per_week = sqlc.narg('days_per_week'))
  AND (sqlc.narg('focus') IS NULL OR focus = sqlc.narg('focus'))
  AND (sqlc.narg('has_amrap') IS NULL OR has_amrap = sqlc.narg('has_amrap'))
  AND (sqlc.narg('search') IS NULL OR name LIKE '%' || sqlc.narg('search') || '%' COLLATE NOCASE)
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountProgramsFiltered :one
SELECT COUNT(*) FROM programs
WHERE (sqlc.narg('difficulty') IS NULL OR difficulty = sqlc.narg('difficulty'))
  AND (sqlc.narg('days_per_week') IS NULL OR days_per_week = sqlc.narg('days_per_week'))
  AND (sqlc.narg('focus') IS NULL OR focus = sqlc.narg('focus'))
  AND (sqlc.narg('has_amrap') IS NULL OR has_amrap = sqlc.narg('has_amrap'))
  AND (sqlc.narg('search') IS NULL OR name LIKE '%' || sqlc.narg('search') || '%' COLLATE NOCASE);

-- name: ListProgramsByNameAsc :many
SELECT id, name, slug, description, cycle_id, weekly_lookup_id, daily_lookup_id, default_rounding, difficulty, days_per_week, focus, has_amrap, created_at, updated_at
FROM programs
ORDER BY name ASC
LIMIT ? OFFSET ?;

-- name: ListProgramsByNameDesc :many
SELECT id, name, slug, description, cycle_id, weekly_lookup_id, daily_lookup_id, default_rounding, difficulty, days_per_week, focus, has_amrap, created_at, updated_at
FROM programs
ORDER BY name DESC
LIMIT ? OFFSET ?;

-- name: ListProgramsByCreatedAtAsc :many
SELECT id, name, slug, description, cycle_id, weekly_lookup_id, daily_lookup_id, default_rounding, difficulty, days_per_week, focus, has_amrap, created_at, updated_at
FROM programs
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: ListProgramsByCreatedAtDesc :many
SELECT id, name, slug, description, cycle_id, weekly_lookup_id, daily_lookup_id, default_rounding, difficulty, days_per_week, focus, has_amrap, created_at, updated_at
FROM programs
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountPrograms :one
SELECT COUNT(*) FROM programs;

-- name: CreateProgram :exec
INSERT INTO programs (id, name, slug, description, cycle_id, weekly_lookup_id, daily_lookup_id, default_rounding, difficulty, days_per_week, focus, has_amrap, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateProgram :exec
UPDATE programs
SET name = ?, slug = ?, description = ?, cycle_id = ?, weekly_lookup_id = ?, daily_lookup_id = ?, default_rounding = ?, difficulty = ?, days_per_week = ?, focus = ?, has_amrap = ?, updated_at = ?
WHERE id = ?;

-- name: DeleteProgram :exec
DELETE FROM programs WHERE id = ?;

-- name: ProgramSlugExists :one
SELECT EXISTS(
    SELECT 1 FROM programs WHERE slug = ?
) AS slug_exists;

-- name: ProgramSlugExistsExcluding :one
SELECT EXISTS(
    SELECT 1 FROM programs WHERE slug = ? AND id != ?
) AS slug_exists;

-- name: ProgramHasEnrolledUsers :one
SELECT EXISTS(
    SELECT 1 FROM user_program_states ups
    WHERE ups.program_id = ?
) AS has_enrolled;

-- name: CountEnrolledUsers :one
SELECT COUNT(*) FROM user_program_states WHERE program_id = ?;

-- name: GetCycleForProgram :one
SELECT c.id, c.name, c.length_weeks, c.created_at, c.updated_at
FROM cycles c
WHERE c.id = ?;

-- name: GetWeeklyLookupForProgram :one
SELECT id, name, entries, program_id, created_at, updated_at
FROM weekly_lookups
WHERE id = ?;

-- name: GetDailyLookupForProgram :one
SELECT id, name, entries, program_id, created_at, updated_at
FROM daily_lookups
WHERE id = ?;

-- name: GetProgramSampleWeek :many
-- Returns days for the first week of a program with prescription counts
-- For programs with week_days, uses week 1; otherwise falls back to days.program_id
SELECT
    d.id,
    d.name,
    COALESCE(wd.day_of_week, 'MONDAY') as day_of_week,
    (SELECT COUNT(*) FROM day_prescriptions dp WHERE dp.day_id = d.id) as exercise_count
FROM days d
LEFT JOIN week_days wd ON wd.day_id = d.id
LEFT JOIN weeks w ON w.id = wd.week_id
LEFT JOIN cycles c ON c.id = w.cycle_id
LEFT JOIN programs p ON p.cycle_id = c.id
WHERE (p.id = ? OR d.program_id = ?)
  AND (w.week_number = 1 OR w.id IS NULL)
GROUP BY d.id
ORDER BY
    CASE wd.day_of_week
        WHEN 'MONDAY' THEN 1
        WHEN 'TUESDAY' THEN 2
        WHEN 'WEDNESDAY' THEN 3
        WHEN 'THURSDAY' THEN 4
        WHEN 'FRIDAY' THEN 5
        WHEN 'SATURDAY' THEN 6
        WHEN 'SUNDAY' THEN 7
        ELSE 8
    END,
    d.name ASC;

-- name: GetProgramLiftRequirements :many
-- Returns unique lift names used in a program, sorted alphabetically
SELECT DISTINCT l.name
FROM lifts l
INNER JOIN prescriptions pr ON pr.lift_id = l.id
INNER JOIN day_prescriptions dp ON dp.prescription_id = pr.id
INNER JOIN days d ON d.id = dp.day_id
WHERE d.program_id = ?
ORDER BY l.name ASC;

-- name: GetProgramSessionStats :one
-- Returns total sets and exercises per average day for session duration estimation
SELECT
    COALESCE(SUM(total_sets), 0) as total_sets,
    COALESCE(COUNT(DISTINCT d.id), 0) as total_days,
    COALESCE(SUM(exercise_count), 0) as total_exercises
FROM days d
LEFT JOIN (
    SELECT
        dp.day_id,
        COUNT(*) as exercise_count,
        SUM(
            CASE
                WHEN pr.set_scheme LIKE '%x%' THEN
                    CAST(SUBSTR(pr.set_scheme, 1, INSTR(pr.set_scheme, 'x') - 1) AS INTEGER)
                ELSE 3
            END
        ) as total_sets
    FROM day_prescriptions dp
    INNER JOIN prescriptions pr ON pr.id = dp.prescription_id
    GROUP BY dp.day_id
) stats ON stats.day_id = d.id
WHERE d.program_id = ?;

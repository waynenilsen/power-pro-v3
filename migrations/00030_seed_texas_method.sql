-- +goose Up
-- Seed Texas Method intermediate program
-- This migration creates the complete Texas Method program with Volume, Recovery,
-- and Intensity days using weekly periodization and percentage-based prescriptions.

-- =============================================================================
-- DETERMINISTIC UUIDs FOR SEEDED DATA
-- =============================================================================
-- Lifts (existing from Starting Strength seed):
--   squat:          00000000-0000-0000-0000-000000000001
--   bench-press:    00000000-0000-0000-0000-000000000002
--   deadlift:       00000000-0000-0000-0000-000000000003
--   overhead-press: 00000000-0000-0000-0000-000000000004
--   power-clean:    00000000-0000-0000-0000-000000000005
--
-- Program entities:
--   program:        texas-method--0000-0000-000000000001
--   cycle:          texas-method--0000-0000-000000000002
--   week:           texas-method--0000-0000-000000000003
--   day-volume:     texas-method--0000-0000-000000000004
--   day-recovery:   texas-method--0000-0000-000000000005
--   day-intensity:  texas-method--0000-0000-000000000006
--
-- Prescriptions (Volume Day - Monday):
--   squat-5x5-vol:     texas-method--0000-0000-000000000010
--   bench-5x5-vol:     texas-method--0000-0000-000000000011
--   press-5x5-vol:     texas-method--0000-0000-000000000012
--   deadlift-1x5:      texas-method--0000-0000-000000000013
--
-- Prescriptions (Recovery Day - Wednesday):
--   squat-2x5-rec:     texas-method--0000-0000-000000000020
--   bench-3x5-rec:     texas-method--0000-0000-000000000021
--   press-3x5-rec:     texas-method--0000-0000-000000000022
--
-- Prescriptions (Intensity Day - Friday):
--   squat-1x5-int:     texas-method--0000-0000-000000000030
--   bench-1x5-int:     texas-method--0000-0000-000000000031
--   press-1x5-int:     texas-method--0000-0000-000000000032
--   clean-5x3:         texas-method--0000-0000-000000000033
--
-- Day prescriptions (join table):
--   vol-squat:         texas-method--0000-0000-000000000040
--   vol-bench:         texas-method--0000-0000-000000000041
--   vol-press:         texas-method--0000-0000-000000000042
--   vol-deadlift:      texas-method--0000-0000-000000000043
--   rec-squat:         texas-method--0000-0000-000000000044
--   rec-bench:         texas-method--0000-0000-000000000045
--   rec-press:         texas-method--0000-0000-000000000046
--   int-squat:         texas-method--0000-0000-000000000047
--   int-bench:         texas-method--0000-0000-000000000048
--   int-press:         texas-method--0000-0000-000000000049
--   int-clean:         texas-method--0000-0000-000000000050
--
-- Week days (join table):
--   week-monday:       texas-method--0000-0000-000000000060
--   week-wednesday:    texas-method--0000-0000-000000000061
--   week-friday:       texas-method--0000-0000-000000000062
--
-- Progressions:
--   weekly-5lb:        texas-method--0000-0000-000000000070
--
-- Program progressions (join table):
--   squat-prog:        texas-method--0000-0000-000000000080
--   bench-prog:        texas-method--0000-0000-000000000081
--   press-prog:        texas-method--0000-0000-000000000082
--   deadlift-prog:     texas-method--0000-0000-000000000083
--   clean-prog:        texas-method--0000-0000-000000000084

-- =============================================================================
-- STEP 1: Lifts already exist from Starting Strength seed (OHP, Power Clean)
-- =============================================================================
-- No new lifts needed - all lifts are already seeded

-- =============================================================================
-- STEP 2: Create cycle (1 week repeating cycle)
-- =============================================================================
-- +goose StatementBegin
INSERT OR IGNORE INTO cycles (id, name, length_weeks, created_at, updated_at) VALUES
    ('texas-method--0000-0000-000000000002', 'Texas Method Weekly Cycle', 1, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 3: Create week 1 in the cycle
-- =============================================================================
-- +goose StatementBegin
INSERT OR IGNORE INTO weeks (id, week_number, variant, cycle_id, created_at, updated_at) VALUES
    ('texas-method--0000-0000-000000000003', 1, NULL, 'texas-method--0000-0000-000000000002', datetime('now'), datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 4: Create workout days (Volume, Recovery, Intensity)
-- =============================================================================
-- +goose StatementBegin
INSERT OR IGNORE INTO days (id, name, slug, metadata, program_id, created_at, updated_at) VALUES
    ('texas-method--0000-0000-000000000004', 'Volume Day', 'volume-day', '{"dayType": "VOLUME"}', NULL, datetime('now'), datetime('now')),
    ('texas-method--0000-0000-000000000005', 'Recovery Day', 'recovery-day', '{"dayType": "RECOVERY"}', NULL, datetime('now'), datetime('now')),
    ('texas-method--0000-0000-000000000006', 'Intensity Day', 'intensity-day', '{"dayType": "INTENSITY"}', NULL, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 5: Create prescriptions for each exercise
-- =============================================================================
-- Texas Method percentage relationships:
--   Friday (Intensity) = 100% of TRAINING_MAX (the anchor)
--   Monday (Volume) = 90% of Friday
--   Wednesday (Recovery) = 80% of Monday (~72% of Friday)

-- Volume Day prescriptions (Monday)
-- +goose StatementBegin
INSERT OR IGNORE INTO prescriptions (id, lift_id, load_strategy, set_scheme, "order", notes, rest_seconds, created_at, updated_at) VALUES
    -- Squat: 5x5 @ 90% of Friday weight
    ('texas-method--0000-0000-000000000010', '00000000-0000-0000-0000-000000000001',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 90.0}',
     '{"type": "FIXED", "sets": 5, "reps": 5}',
     0, 'High volume squat work', 180, datetime('now'), datetime('now')),
    -- Bench Press: 5x5 @ 90% (both press lifts available, user alternates focus)
    ('texas-method--0000-0000-000000000011', '00000000-0000-0000-0000-000000000002',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 90.0}',
     '{"type": "FIXED", "sets": 5, "reps": 5}',
     1, NULL, 180, datetime('now'), datetime('now')),
    -- Overhead Press: 5x5 @ 90%
    ('texas-method--0000-0000-000000000012', '00000000-0000-0000-0000-000000000004',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 90.0}',
     '{"type": "FIXED", "sets": 5, "reps": 5}',
     2, NULL, 180, datetime('now'), datetime('now')),
    -- Deadlift: 1x5 @ 90%
    ('texas-method--0000-0000-000000000013', '00000000-0000-0000-0000-000000000003',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 90.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     3, 'Single heavy set', 180, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- Recovery Day prescriptions (Wednesday)
-- +goose StatementBegin
INSERT OR IGNORE INTO prescriptions (id, lift_id, load_strategy, set_scheme, "order", notes, rest_seconds, created_at, updated_at) VALUES
    -- Squat: 2x5 @ 72% (~80% of Monday which is 90% of Friday)
    ('texas-method--0000-0000-000000000020', '00000000-0000-0000-0000-000000000001',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 72.0}',
     '{"type": "FIXED", "sets": 2, "reps": 5}',
     0, 'Light recovery work', 120, datetime('now'), datetime('now')),
    -- Bench Press: 3x5 @ 90%
    ('texas-method--0000-0000-000000000021', '00000000-0000-0000-0000-000000000002',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 90.0}',
     '{"type": "FIXED", "sets": 3, "reps": 5}',
     1, NULL, 120, datetime('now'), datetime('now')),
    -- Overhead Press: 3x5 @ 90%
    ('texas-method--0000-0000-000000000022', '00000000-0000-0000-0000-000000000004',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 90.0}',
     '{"type": "FIXED", "sets": 3, "reps": 5}',
     2, NULL, 120, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- Intensity Day prescriptions (Friday)
-- +goose StatementBegin
INSERT OR IGNORE INTO prescriptions (id, lift_id, load_strategy, set_scheme, "order", notes, rest_seconds, created_at, updated_at) VALUES
    -- Squat: 1x5 @ 100% (PR attempt)
    ('texas-method--0000-0000-000000000030', '00000000-0000-0000-0000-000000000001',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 100.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     0, 'PR attempt - new 5RM', 300, datetime('now'), datetime('now')),
    -- Bench Press: 1x5 @ 100%
    ('texas-method--0000-0000-000000000031', '00000000-0000-0000-0000-000000000002',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 100.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     1, 'PR attempt', 300, datetime('now'), datetime('now')),
    -- Overhead Press: 1x5 @ 100%
    ('texas-method--0000-0000-000000000032', '00000000-0000-0000-0000-000000000004',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 100.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     2, 'PR attempt', 300, datetime('now'), datetime('now')),
    -- Power Clean: 5x3 @ 100% (optional lift)
    ('texas-method--0000-0000-000000000033', '00000000-0000-0000-0000-000000000005',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 100.0}',
     '{"type": "FIXED", "sets": 5, "reps": 3}',
     3, 'Technical precision emphasis', 180, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 6: Link prescriptions to days
-- =============================================================================
-- +goose StatementBegin
INSERT OR IGNORE INTO day_prescriptions (id, day_id, prescription_id, "order", created_at) VALUES
    -- Volume Day (Monday): Squat 5x5, Bench 5x5, Press 5x5, Deadlift 1x5
    ('texas-method--0000-0000-000000000040', 'texas-method--0000-0000-000000000004', 'texas-method--0000-0000-000000000010', 0, datetime('now')),
    ('texas-method--0000-0000-000000000041', 'texas-method--0000-0000-000000000004', 'texas-method--0000-0000-000000000011', 1, datetime('now')),
    ('texas-method--0000-0000-000000000042', 'texas-method--0000-0000-000000000004', 'texas-method--0000-0000-000000000012', 2, datetime('now')),
    ('texas-method--0000-0000-000000000043', 'texas-method--0000-0000-000000000004', 'texas-method--0000-0000-000000000013', 3, datetime('now')),
    -- Recovery Day (Wednesday): Squat 2x5, Bench 3x5, Press 3x5
    ('texas-method--0000-0000-000000000044', 'texas-method--0000-0000-000000000005', 'texas-method--0000-0000-000000000020', 0, datetime('now')),
    ('texas-method--0000-0000-000000000045', 'texas-method--0000-0000-000000000005', 'texas-method--0000-0000-000000000021', 1, datetime('now')),
    ('texas-method--0000-0000-000000000046', 'texas-method--0000-0000-000000000005', 'texas-method--0000-0000-000000000022', 2, datetime('now')),
    -- Intensity Day (Friday): Squat 1x5, Bench 1x5, Press 1x5, Power Clean 5x3
    ('texas-method--0000-0000-000000000047', 'texas-method--0000-0000-000000000006', 'texas-method--0000-0000-000000000030', 0, datetime('now')),
    ('texas-method--0000-0000-000000000048', 'texas-method--0000-0000-000000000006', 'texas-method--0000-0000-000000000031', 1, datetime('now')),
    ('texas-method--0000-0000-000000000049', 'texas-method--0000-0000-000000000006', 'texas-method--0000-0000-000000000032', 2, datetime('now')),
    ('texas-method--0000-0000-000000000050', 'texas-method--0000-0000-000000000006', 'texas-method--0000-0000-000000000033', 3, datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 7: Link days to week (Volume/Recovery/Intensity: Mon/Wed/Fri)
-- =============================================================================
-- +goose StatementBegin
INSERT OR IGNORE INTO week_days (id, week_id, day_id, day_of_week, created_at) VALUES
    ('texas-method--0000-0000-000000000060', 'texas-method--0000-0000-000000000003', 'texas-method--0000-0000-000000000004', 'MONDAY', datetime('now')),
    ('texas-method--0000-0000-000000000061', 'texas-method--0000-0000-000000000003', 'texas-method--0000-0000-000000000005', 'WEDNESDAY', datetime('now')),
    ('texas-method--0000-0000-000000000062', 'texas-method--0000-0000-000000000003', 'texas-method--0000-0000-000000000006', 'FRIDAY', datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 8: Create program
-- =============================================================================
-- +goose StatementBegin
INSERT OR IGNORE INTO programs (id, name, slug, description, cycle_id, weekly_lookup_id, daily_lookup_id, default_rounding, created_at, updated_at) VALUES
    ('texas-method--0000-0000-000000000001', 'Texas Method', 'texas-method',
     'Intermediate strength program using weekly periodization with Volume, Recovery, and Intensity days. Designed for lifters who have exhausted linear progression. Friday intensity day drives new personal records, with Monday and Wednesday supporting adaptation.',
     'texas-method--0000-0000-000000000002', NULL, NULL, 5.0, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 9: Update days to reference the program
-- =============================================================================
-- +goose StatementBegin
UPDATE days SET program_id = 'texas-method--0000-0000-000000000001'
WHERE id IN ('texas-method--0000-0000-000000000004', 'texas-method--0000-0000-000000000005', 'texas-method--0000-0000-000000000006');
-- +goose StatementEnd

-- =============================================================================
-- STEP 10: Create progression rules
-- =============================================================================
-- Texas Method uses weekly progression: +5lb per week for all lifts
-- Progression is based on successful completion of Intensity Day (Friday)
-- +goose StatementBegin
INSERT OR IGNORE INTO progressions (id, name, type, parameters, created_at, updated_at) VALUES
    -- +5lb weekly progression (all lifts)
    ('texas-method--0000-0000-000000000070', 'Texas Method Weekly +5lb', 'LINEAR_PROGRESSION',
     '{"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_WEEK"}',
     datetime('now'), datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 11: Link progressions to program for each lift
-- =============================================================================
-- All lifts use +5lb/week progression in Texas Method
-- +goose StatementBegin
INSERT OR IGNORE INTO program_progressions (id, program_id, progression_id, lift_id, priority, enabled, override_increment, created_at, updated_at) VALUES
    -- Squat: +5lb per week
    ('texas-method--0000-0000-000000000080', 'texas-method--0000-0000-000000000001', 'texas-method--0000-0000-000000000070', '00000000-0000-0000-0000-000000000001', 1, 1, NULL, datetime('now'), datetime('now')),
    -- Bench Press: +5lb per week (can override to +2.5lb if needed)
    ('texas-method--0000-0000-000000000081', 'texas-method--0000-0000-000000000001', 'texas-method--0000-0000-000000000070', '00000000-0000-0000-0000-000000000002', 2, 1, 2.5, datetime('now'), datetime('now')),
    -- Overhead Press: +5lb per week (can override to +2.5lb if needed)
    ('texas-method--0000-0000-000000000082', 'texas-method--0000-0000-000000000001', 'texas-method--0000-0000-000000000070', '00000000-0000-0000-0000-000000000004', 3, 1, 2.5, datetime('now'), datetime('now')),
    -- Deadlift: +5lb per week
    ('texas-method--0000-0000-000000000083', 'texas-method--0000-0000-000000000001', 'texas-method--0000-0000-000000000070', '00000000-0000-0000-0000-000000000003', 4, 1, NULL, datetime('now'), datetime('now')),
    -- Power Clean: +5lb per week (can override to +2.5lb if needed)
    ('texas-method--0000-0000-000000000084', 'texas-method--0000-0000-000000000001', 'texas-method--0000-0000-000000000070', '00000000-0000-0000-0000-000000000005', 5, 1, 2.5, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- +goose Down
-- =============================================================================
-- DOWN MIGRATION: Remove all Texas Method seeded data
-- =============================================================================
-- Order matters due to foreign key constraints - delete in reverse order

-- Remove program progressions
-- +goose StatementBegin
DELETE FROM program_progressions WHERE id IN (
    'texas-method--0000-0000-000000000080',
    'texas-method--0000-0000-000000000081',
    'texas-method--0000-0000-000000000082',
    'texas-method--0000-0000-000000000083',
    'texas-method--0000-0000-000000000084'
);
-- +goose StatementEnd

-- Remove progressions
-- +goose StatementBegin
DELETE FROM progressions WHERE id IN (
    'texas-method--0000-0000-000000000070'
);
-- +goose StatementEnd

-- Remove week_days
-- +goose StatementBegin
DELETE FROM week_days WHERE id IN (
    'texas-method--0000-0000-000000000060',
    'texas-method--0000-0000-000000000061',
    'texas-method--0000-0000-000000000062'
);
-- +goose StatementEnd

-- Remove day_prescriptions
-- +goose StatementBegin
DELETE FROM day_prescriptions WHERE id IN (
    'texas-method--0000-0000-000000000040',
    'texas-method--0000-0000-000000000041',
    'texas-method--0000-0000-000000000042',
    'texas-method--0000-0000-000000000043',
    'texas-method--0000-0000-000000000044',
    'texas-method--0000-0000-000000000045',
    'texas-method--0000-0000-000000000046',
    'texas-method--0000-0000-000000000047',
    'texas-method--0000-0000-000000000048',
    'texas-method--0000-0000-000000000049',
    'texas-method--0000-0000-000000000050'
);
-- +goose StatementEnd

-- Remove prescriptions
-- +goose StatementBegin
DELETE FROM prescriptions WHERE id IN (
    'texas-method--0000-0000-000000000010',
    'texas-method--0000-0000-000000000011',
    'texas-method--0000-0000-000000000012',
    'texas-method--0000-0000-000000000013',
    'texas-method--0000-0000-000000000020',
    'texas-method--0000-0000-000000000021',
    'texas-method--0000-0000-000000000022',
    'texas-method--0000-0000-000000000030',
    'texas-method--0000-0000-000000000031',
    'texas-method--0000-0000-000000000032',
    'texas-method--0000-0000-000000000033'
);
-- +goose StatementEnd

-- Remove program (must come before days due to FK)
-- +goose StatementBegin
DELETE FROM programs WHERE id = 'texas-method--0000-0000-000000000001';
-- +goose StatementEnd

-- Remove days
-- +goose StatementBegin
DELETE FROM days WHERE id IN (
    'texas-method--0000-0000-000000000004',
    'texas-method--0000-0000-000000000005',
    'texas-method--0000-0000-000000000006'
);
-- +goose StatementEnd

-- Remove week
-- +goose StatementBegin
DELETE FROM weeks WHERE id = 'texas-method--0000-0000-000000000003';
-- +goose StatementEnd

-- Remove cycle
-- +goose StatementBegin
DELETE FROM cycles WHERE id = 'texas-method--0000-0000-000000000002';
-- +goose StatementEnd

-- Note: We do NOT remove lifts in the down migration as they are shared
-- canonical reference data that may be used by other programs or user data.

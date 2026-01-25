-- +goose Up
-- Seed Starting Strength novice program
-- This migration creates the complete Starting Strength program with A/B rotation,
-- all prescriptions, and linear progression rules.

-- =============================================================================
-- DETERMINISTIC UUIDs FOR SEEDED DATA
-- =============================================================================
-- Lifts (existing):
--   squat:       00000000-0000-0000-0000-000000000001
--   bench-press: 00000000-0000-0000-0000-000000000002
--   deadlift:    00000000-0000-0000-0000-000000000003
--
-- New lifts for Starting Strength:
--   overhead-press: 00000000-0000-0000-0000-000000000004
--   power-clean:    00000000-0000-0000-0000-000000000005
--
-- Program entities:
--   program:         starting-strength-0000-0000-000000000001
--   cycle:           starting-strength-0000-0000-000000000002
--   week:            starting-strength-0000-0000-000000000003
--   day-a:           starting-strength-0000-0000-000000000004
--   day-b:           starting-strength-0000-0000-000000000005
--
-- Prescriptions:
--   squat-3x5-a:     starting-strength-0000-0000-000000000010
--   bench-3x5:       starting-strength-0000-0000-000000000011
--   deadlift-1x5:    starting-strength-0000-0000-000000000012
--   squat-3x5-b:     starting-strength-0000-0000-000000000013
--   press-3x5:       starting-strength-0000-0000-000000000014
--   clean-5x3:       starting-strength-0000-0000-000000000015
--
-- Day prescriptions (join table):
--   day-a-squat:     starting-strength-0000-0000-000000000020
--   day-a-bench:     starting-strength-0000-0000-000000000021
--   day-a-deadlift:  starting-strength-0000-0000-000000000022
--   day-b-squat:     starting-strength-0000-0000-000000000023
--   day-b-press:     starting-strength-0000-0000-000000000024
--   day-b-clean:     starting-strength-0000-0000-000000000025
--
-- Week days (join table):
--   week-monday:     starting-strength-0000-0000-000000000030
--   week-wednesday:  starting-strength-0000-0000-000000000031
--   week-friday:     starting-strength-0000-0000-000000000032
--
-- Progressions:
--   linear-5lb:      starting-strength-0000-0000-000000000040
--   linear-10lb:     starting-strength-0000-0000-000000000041
--
-- Program progressions (join table):
--   squat-prog:      starting-strength-0000-0000-000000000050
--   bench-prog:      starting-strength-0000-0000-000000000051
--   press-prog:      starting-strength-0000-0000-000000000052
--   deadlift-prog:   starting-strength-0000-0000-000000000053
--   clean-prog:      starting-strength-0000-0000-000000000054

-- =============================================================================
-- STEP 1: Add OHP and Power Clean lifts (if not exists)
-- =============================================================================
-- +goose StatementBegin
INSERT OR IGNORE INTO lifts (id, name, slug, is_competition_lift, parent_lift_id, created_at, updated_at) VALUES
    ('00000000-0000-0000-0000-000000000004', 'Overhead Press', 'overhead-press', 0, NULL, datetime('now'), datetime('now')),
    ('00000000-0000-0000-0000-000000000005', 'Power Clean', 'power-clean', 0, NULL, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 2: Create cycle (1 week repeating cycle)
-- =============================================================================
-- +goose StatementBegin
INSERT OR IGNORE INTO cycles (id, name, length_weeks, created_at, updated_at) VALUES
    ('starting-strength-0000-0000-000000000002', 'Starting Strength A/B Cycle', 1, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 3: Create week 1 in the cycle
-- =============================================================================
-- +goose StatementBegin
INSERT OR IGNORE INTO weeks (id, week_number, variant, cycle_id, created_at, updated_at) VALUES
    ('starting-strength-0000-0000-000000000003', 1, NULL, 'starting-strength-0000-0000-000000000002', datetime('now'), datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 4: Create workout days (Day A and Day B)
-- =============================================================================
-- +goose StatementBegin
INSERT OR IGNORE INTO days (id, name, slug, metadata, program_id, created_at, updated_at) VALUES
    ('starting-strength-0000-0000-000000000004', 'Workout A', 'workout-a', NULL, NULL, datetime('now'), datetime('now')),
    ('starting-strength-0000-0000-000000000005', 'Workout B', 'workout-b', NULL, NULL, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 5: Create prescriptions for each exercise
-- =============================================================================
-- Load strategy: PERCENT_OF at 100% of TRAINING_MAX (user works at their training max)
-- Set schemes: FIXED with specified sets x reps

-- Day A prescriptions
-- +goose StatementBegin
INSERT OR IGNORE INTO prescriptions (id, lift_id, load_strategy, set_scheme, "order", notes, rest_seconds, created_at, updated_at) VALUES
    -- Squat: 3x5 (Day A)
    ('starting-strength-0000-0000-000000000010', '00000000-0000-0000-0000-000000000001',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 100.0}',
     '{"type": "FIXED", "sets": 3, "reps": 5}',
     0, NULL, 180, datetime('now'), datetime('now')),
    -- Bench Press: 3x5
    ('starting-strength-0000-0000-000000000011', '00000000-0000-0000-0000-000000000002',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 100.0}',
     '{"type": "FIXED", "sets": 3, "reps": 5}',
     1, NULL, 180, datetime('now'), datetime('now')),
    -- Deadlift: 1x5
    ('starting-strength-0000-0000-000000000012', '00000000-0000-0000-0000-000000000003',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 100.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     2, NULL, 180, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- Day B prescriptions
-- +goose StatementBegin
INSERT OR IGNORE INTO prescriptions (id, lift_id, load_strategy, set_scheme, "order", notes, rest_seconds, created_at, updated_at) VALUES
    -- Squat: 3x5 (Day B - separate prescription for independent tracking)
    ('starting-strength-0000-0000-000000000013', '00000000-0000-0000-0000-000000000001',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 100.0}',
     '{"type": "FIXED", "sets": 3, "reps": 5}',
     0, NULL, 180, datetime('now'), datetime('now')),
    -- Overhead Press: 3x5
    ('starting-strength-0000-0000-000000000014', '00000000-0000-0000-0000-000000000004',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 100.0}',
     '{"type": "FIXED", "sets": 3, "reps": 5}',
     1, NULL, 180, datetime('now'), datetime('now')),
    -- Power Clean: 5x3
    ('starting-strength-0000-0000-000000000015', '00000000-0000-0000-0000-000000000005',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 100.0}',
     '{"type": "FIXED", "sets": 5, "reps": 3}',
     2, NULL, 120, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 6: Link prescriptions to days
-- =============================================================================
-- +goose StatementBegin
INSERT OR IGNORE INTO day_prescriptions (id, day_id, prescription_id, "order", created_at) VALUES
    -- Day A: Squat, Bench, Deadlift
    ('starting-strength-0000-0000-000000000020', 'starting-strength-0000-0000-000000000004', 'starting-strength-0000-0000-000000000010', 0, datetime('now')),
    ('starting-strength-0000-0000-000000000021', 'starting-strength-0000-0000-000000000004', 'starting-strength-0000-0000-000000000011', 1, datetime('now')),
    ('starting-strength-0000-0000-000000000022', 'starting-strength-0000-0000-000000000004', 'starting-strength-0000-0000-000000000012', 2, datetime('now')),
    -- Day B: Squat, Press, Power Clean
    ('starting-strength-0000-0000-000000000023', 'starting-strength-0000-0000-000000000005', 'starting-strength-0000-0000-000000000013', 0, datetime('now')),
    ('starting-strength-0000-0000-000000000024', 'starting-strength-0000-0000-000000000005', 'starting-strength-0000-0000-000000000014', 1, datetime('now')),
    ('starting-strength-0000-0000-000000000025', 'starting-strength-0000-0000-000000000005', 'starting-strength-0000-0000-000000000015', 2, datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 7: Link days to week (A/B/A pattern: Mon/Wed/Fri)
-- =============================================================================
-- +goose StatementBegin
INSERT OR IGNORE INTO week_days (id, week_id, day_id, day_of_week, created_at) VALUES
    ('starting-strength-0000-0000-000000000030', 'starting-strength-0000-0000-000000000003', 'starting-strength-0000-0000-000000000004', 'MONDAY', datetime('now')),
    ('starting-strength-0000-0000-000000000031', 'starting-strength-0000-0000-000000000003', 'starting-strength-0000-0000-000000000005', 'WEDNESDAY', datetime('now')),
    ('starting-strength-0000-0000-000000000032', 'starting-strength-0000-0000-000000000003', 'starting-strength-0000-0000-000000000004', 'FRIDAY', datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 8: Create program
-- =============================================================================
-- +goose StatementBegin
INSERT OR IGNORE INTO programs (id, name, slug, description, cycle_id, weekly_lookup_id, daily_lookup_id, default_rounding, created_at, updated_at) VALUES
    ('starting-strength-0000-0000-000000000001', 'Starting Strength', 'starting-strength',
     'Mark Rippetoe''s classic novice linear progression program. Features A/B rotation with compound barbell movements, designed for beginners to build strength through progressive overload.',
     'starting-strength-0000-0000-000000000002', NULL, NULL, 5.0, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 9: Update days to reference the program
-- =============================================================================
-- +goose StatementBegin
UPDATE days SET program_id = 'starting-strength-0000-0000-000000000001'
WHERE id IN ('starting-strength-0000-0000-000000000004', 'starting-strength-0000-0000-000000000005');
-- +goose StatementEnd

-- =============================================================================
-- STEP 10: Create progression rules
-- =============================================================================
-- Linear progression: +5lb for upper body lifts, +10lb for lower body lifts
-- +goose StatementBegin
INSERT OR IGNORE INTO progressions (id, name, type, parameters, created_at, updated_at) VALUES
    -- +5lb progression (upper body: bench, press, power clean)
    ('starting-strength-0000-0000-000000000040', 'Starting Strength +5lb', 'LINEAR_PROGRESSION',
     '{"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"}',
     datetime('now'), datetime('now')),
    -- +10lb progression (lower body: squat, deadlift)
    ('starting-strength-0000-0000-000000000041', 'Starting Strength +10lb', 'LINEAR_PROGRESSION',
     '{"increment": 10.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"}',
     datetime('now'), datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 11: Link progressions to program for each lift
-- =============================================================================
-- +goose StatementBegin
INSERT OR IGNORE INTO program_progressions (id, program_id, progression_id, lift_id, priority, enabled, override_increment, created_at, updated_at) VALUES
    -- Squat: +10lb per session (lower body)
    ('starting-strength-0000-0000-000000000050', 'starting-strength-0000-0000-000000000001', 'starting-strength-0000-0000-000000000041', '00000000-0000-0000-0000-000000000001', 1, 1, NULL, datetime('now'), datetime('now')),
    -- Bench Press: +5lb per session (upper body)
    ('starting-strength-0000-0000-000000000051', 'starting-strength-0000-0000-000000000001', 'starting-strength-0000-0000-000000000040', '00000000-0000-0000-0000-000000000002', 2, 1, NULL, datetime('now'), datetime('now')),
    -- Overhead Press: +5lb per session (upper body)
    ('starting-strength-0000-0000-000000000052', 'starting-strength-0000-0000-000000000001', 'starting-strength-0000-0000-000000000040', '00000000-0000-0000-0000-000000000004', 3, 1, NULL, datetime('now'), datetime('now')),
    -- Deadlift: +10lb per session (lower body)
    ('starting-strength-0000-0000-000000000053', 'starting-strength-0000-0000-000000000001', 'starting-strength-0000-0000-000000000041', '00000000-0000-0000-0000-000000000003', 4, 1, NULL, datetime('now'), datetime('now')),
    -- Power Clean: +5lb per session (upper body classification for progression)
    ('starting-strength-0000-0000-000000000054', 'starting-strength-0000-0000-000000000001', 'starting-strength-0000-0000-000000000040', '00000000-0000-0000-0000-000000000005', 5, 1, NULL, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- +goose Down
-- =============================================================================
-- DOWN MIGRATION: Remove all Starting Strength seeded data
-- =============================================================================
-- Order matters due to foreign key constraints - delete in reverse order

-- Remove program progressions
-- +goose StatementBegin
DELETE FROM program_progressions WHERE id IN (
    'starting-strength-0000-0000-000000000050',
    'starting-strength-0000-0000-000000000051',
    'starting-strength-0000-0000-000000000052',
    'starting-strength-0000-0000-000000000053',
    'starting-strength-0000-0000-000000000054'
);
-- +goose StatementEnd

-- Remove progressions
-- +goose StatementBegin
DELETE FROM progressions WHERE id IN (
    'starting-strength-0000-0000-000000000040',
    'starting-strength-0000-0000-000000000041'
);
-- +goose StatementEnd

-- Remove week_days
-- +goose StatementBegin
DELETE FROM week_days WHERE id IN (
    'starting-strength-0000-0000-000000000030',
    'starting-strength-0000-0000-000000000031',
    'starting-strength-0000-0000-000000000032'
);
-- +goose StatementEnd

-- Remove day_prescriptions
-- +goose StatementBegin
DELETE FROM day_prescriptions WHERE id IN (
    'starting-strength-0000-0000-000000000020',
    'starting-strength-0000-0000-000000000021',
    'starting-strength-0000-0000-000000000022',
    'starting-strength-0000-0000-000000000023',
    'starting-strength-0000-0000-000000000024',
    'starting-strength-0000-0000-000000000025'
);
-- +goose StatementEnd

-- Remove prescriptions
-- +goose StatementBegin
DELETE FROM prescriptions WHERE id IN (
    'starting-strength-0000-0000-000000000010',
    'starting-strength-0000-0000-000000000011',
    'starting-strength-0000-0000-000000000012',
    'starting-strength-0000-0000-000000000013',
    'starting-strength-0000-0000-000000000014',
    'starting-strength-0000-0000-000000000015'
);
-- +goose StatementEnd

-- Remove program (must come before days due to FK)
-- +goose StatementBegin
DELETE FROM programs WHERE id = 'starting-strength-0000-0000-000000000001';
-- +goose StatementEnd

-- Remove days
-- +goose StatementBegin
DELETE FROM days WHERE id IN (
    'starting-strength-0000-0000-000000000004',
    'starting-strength-0000-0000-000000000005'
);
-- +goose StatementEnd

-- Remove week
-- +goose StatementBegin
DELETE FROM weeks WHERE id = 'starting-strength-0000-0000-000000000003';
-- +goose StatementEnd

-- Remove cycle
-- +goose StatementBegin
DELETE FROM cycles WHERE id = 'starting-strength-0000-0000-000000000002';
-- +goose StatementEnd

-- Note: We do NOT remove the OHP and Power Clean lifts in the down migration
-- because they may be used by other programs or user data. Lift types are
-- considered canonical reference data that persists.

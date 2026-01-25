-- +goose Up
-- Seed GZCLP (GZCL Linear Progression) novice/early-intermediate program
-- This migration creates the 4-day GZCLP program with its tiered system (T1/T2),
-- AMRAP sets, and unique stage-based progression on failure.

-- =============================================================================
-- DETERMINISTIC UUIDs FOR SEEDED DATA
-- =============================================================================
-- Lifts (existing from Starting Strength seed):
--   squat:          00000000-0000-0000-0000-000000000001
--   bench-press:    00000000-0000-0000-0000-000000000002
--   deadlift:       00000000-0000-0000-0000-000000000003
--   overhead-press: 00000000-0000-0000-0000-000000000004
--
-- Program entities:
--   program:        gzclp----------0000-0000-000000000001
--   cycle:          gzclp----------0000-0000-000000000002
--   week:           gzclp----------0000-0000-000000000003
--
-- Days:
--   day-1:          gzclp----------0000-0000-000000000010
--   day-2:          gzclp----------0000-0000-000000000011
--   day-3:          gzclp----------0000-0000-000000000012
--   day-4:          gzclp----------0000-0000-000000000013
--
-- Prescriptions (T1 - 5x3+ default):
--   d1-t1-squat:    gzclp----------0000-0000-000000000100
--   d2-t1-ohp:      gzclp----------0000-0000-000000000101
--   d3-t1-bench:    gzclp----------0000-0000-000000000102
--   d4-t1-deadlift: gzclp----------0000-0000-000000000103
--
-- Prescriptions (T2 - 3x10 default):
--   d1-t2-bench:    gzclp----------0000-0000-000000000110
--   d2-t2-deadlift: gzclp----------0000-0000-000000000111
--   d3-t2-squat:    gzclp----------0000-0000-000000000112
--   d4-t2-ohp:      gzclp----------0000-0000-000000000113
--
-- Day prescriptions (join table):
--   d1-t1:          gzclp----------0000-0000-000000000200
--   d1-t2:          gzclp----------0000-0000-000000000201
--   d2-t1:          gzclp----------0000-0000-000000000202
--   d2-t2:          gzclp----------0000-0000-000000000203
--   d3-t1:          gzclp----------0000-0000-000000000204
--   d3-t2:          gzclp----------0000-0000-000000000205
--   d4-t1:          gzclp----------0000-0000-000000000206
--   d4-t2:          gzclp----------0000-0000-000000000207
--
-- Week days (join table):
--   week-d1:        gzclp----------0000-0000-000000000300
--   week-d2:        gzclp----------0000-0000-000000000301
--   week-d3:        gzclp----------0000-0000-000000000302
--   week-d4:        gzclp----------0000-0000-000000000303
--
-- Progressions:
--   t1-lower-5lb:   gzclp----------0000-0000-000000000400
--   t1-upper-2.5lb: gzclp----------0000-0000-000000000401
--   t2-2.5lb:       gzclp----------0000-0000-000000000402
--   t1-stage-prog:  gzclp----------0000-0000-000000000403
--   t2-stage-prog:  gzclp----------0000-0000-000000000404
--
-- Program progressions:
--   squat-t1:       gzclp----------0000-0000-000000000500
--   bench-t1:       gzclp----------0000-0000-000000000501
--   deadlift-t1:    gzclp----------0000-0000-000000000502
--   ohp-t1:         gzclp----------0000-0000-000000000503
--   squat-t2:       gzclp----------0000-0000-000000000504
--   bench-t2:       gzclp----------0000-0000-000000000505
--   deadlift-t2:    gzclp----------0000-0000-000000000506
--   ohp-t2:         gzclp----------0000-0000-000000000507

-- =============================================================================
-- STEP 1: Lifts already exist from Starting Strength seed
-- =============================================================================
-- No new lifts needed - all lifts (squat, bench, deadlift, OHP) are already seeded

-- =============================================================================
-- STEP 2: Create cycle (1 week repeating cycle)
-- =============================================================================
-- GZCLP uses a simple weekly cycle with 4 training days
-- +goose StatementBegin
INSERT OR IGNORE INTO cycles (id, name, length_weeks, created_at, updated_at) VALUES
    ('gzclp----------0000-0000-000000000002', 'GZCLP Weekly Cycle', 1, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 3: Create week 1 in the cycle
-- =============================================================================
-- +goose StatementBegin
INSERT OR IGNORE INTO weeks (id, week_number, variant, cycle_id, created_at, updated_at) VALUES
    ('gzclp----------0000-0000-000000000003', 1, NULL, 'gzclp----------0000-0000-000000000002', datetime('now'), datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 4: Create workout days
-- =============================================================================
-- GZCLP 4-day program:
--   Day 1: T1 Squat, T2 Bench
--   Day 2: T1 OHP, T2 Deadlift
--   Day 3: T1 Bench, T2 Squat
--   Day 4: T1 Deadlift, T2 OHP
-- +goose StatementBegin
INSERT OR IGNORE INTO days (id, name, slug, metadata, program_id, created_at, updated_at) VALUES
    ('gzclp----------0000-0000-000000000010', 'Day 1 - Squat/Bench', 'gzclp-day-1', '{"t1Lift": "squat", "t2Lift": "bench-press"}', NULL, datetime('now'), datetime('now')),
    ('gzclp----------0000-0000-000000000011', 'Day 2 - OHP/Deadlift', 'gzclp-day-2', '{"t1Lift": "overhead-press", "t2Lift": "deadlift"}', NULL, datetime('now'), datetime('now')),
    ('gzclp----------0000-0000-000000000012', 'Day 3 - Bench/Squat', 'gzclp-day-3', '{"t1Lift": "bench-press", "t2Lift": "squat"}', NULL, datetime('now'), datetime('now')),
    ('gzclp----------0000-0000-000000000013', 'Day 4 - Deadlift/OHP', 'gzclp-day-4', '{"t1Lift": "deadlift", "t2Lift": "overhead-press"}', NULL, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 5: Create prescriptions
-- =============================================================================
-- T1 Prescriptions (5x3+ default, last set AMRAP)
-- T1 uses working weight (TRAINING_MAX) and has staged progression on failure
-- Load strategy uses TRAINING_MAX which represents the lifter's current working weight
-- +goose StatementBegin
INSERT OR IGNORE INTO prescriptions (id, lift_id, load_strategy, set_scheme, "order", notes, rest_seconds, created_at, updated_at) VALUES
    -- Day 1 T1: Squat 5x3+ (last set AMRAP)
    ('gzclp----------0000-0000-000000000100', '00000000-0000-0000-0000-000000000001',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 100.0}',
     '{"type": "FIXED", "sets": 5, "reps": 3, "isAmrap": true, "tier": "T1", "stage": 1}',
     0, 'T1: 5x3+ (Stage 1)', 180, datetime('now'), datetime('now')),
    -- Day 2 T1: OHP 5x3+ (last set AMRAP)
    ('gzclp----------0000-0000-000000000101', '00000000-0000-0000-0000-000000000004',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 100.0}',
     '{"type": "FIXED", "sets": 5, "reps": 3, "isAmrap": true, "tier": "T1", "stage": 1}',
     0, 'T1: 5x3+ (Stage 1)', 180, datetime('now'), datetime('now')),
    -- Day 3 T1: Bench 5x3+ (last set AMRAP)
    ('gzclp----------0000-0000-000000000102', '00000000-0000-0000-0000-000000000002',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 100.0}',
     '{"type": "FIXED", "sets": 5, "reps": 3, "isAmrap": true, "tier": "T1", "stage": 1}',
     0, 'T1: 5x3+ (Stage 1)', 180, datetime('now'), datetime('now')),
    -- Day 4 T1: Deadlift 5x3+ (last set AMRAP)
    ('gzclp----------0000-0000-000000000103', '00000000-0000-0000-0000-000000000003',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 100.0}',
     '{"type": "FIXED", "sets": 5, "reps": 3, "isAmrap": true, "tier": "T1", "stage": 1}',
     0, 'T1: 5x3+ (Stage 1)', 180, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- T2 Prescriptions (3x10 default)
-- T2 uses a percentage of T1 weight (typically 50-65% recommended)
-- +goose StatementBegin
INSERT OR IGNORE INTO prescriptions (id, lift_id, load_strategy, set_scheme, "order", notes, rest_seconds, created_at, updated_at) VALUES
    -- Day 1 T2: Bench 3x10
    ('gzclp----------0000-0000-000000000110', '00000000-0000-0000-0000-000000000002',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 100.0}',
     '{"type": "FIXED", "sets": 3, "reps": 10, "tier": "T2", "stage": 1}',
     1, 'T2: 3x10 (Stage 1)', 120, datetime('now'), datetime('now')),
    -- Day 2 T2: Deadlift 3x10
    ('gzclp----------0000-0000-000000000111', '00000000-0000-0000-0000-000000000003',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 100.0}',
     '{"type": "FIXED", "sets": 3, "reps": 10, "tier": "T2", "stage": 1}',
     1, 'T2: 3x10 (Stage 1)', 120, datetime('now'), datetime('now')),
    -- Day 3 T2: Squat 3x10
    ('gzclp----------0000-0000-000000000112', '00000000-0000-0000-0000-000000000001',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 100.0}',
     '{"type": "FIXED", "sets": 3, "reps": 10, "tier": "T2", "stage": 1}',
     1, 'T2: 3x10 (Stage 1)', 120, datetime('now'), datetime('now')),
    -- Day 4 T2: OHP 3x10
    ('gzclp----------0000-0000-000000000113', '00000000-0000-0000-0000-000000000004',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 100.0}',
     '{"type": "FIXED", "sets": 3, "reps": 10, "tier": "T2", "stage": 1}',
     1, 'T2: 3x10 (Stage 1)', 120, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 6: Link prescriptions to days
-- =============================================================================
-- +goose StatementBegin
INSERT OR IGNORE INTO day_prescriptions (id, day_id, prescription_id, "order", created_at) VALUES
    -- Day 1: T1 Squat, T2 Bench
    ('gzclp----------0000-0000-000000000200', 'gzclp----------0000-0000-000000000010', 'gzclp----------0000-0000-000000000100', 0, datetime('now')),
    ('gzclp----------0000-0000-000000000201', 'gzclp----------0000-0000-000000000010', 'gzclp----------0000-0000-000000000110', 1, datetime('now')),
    -- Day 2: T1 OHP, T2 Deadlift
    ('gzclp----------0000-0000-000000000202', 'gzclp----------0000-0000-000000000011', 'gzclp----------0000-0000-000000000101', 0, datetime('now')),
    ('gzclp----------0000-0000-000000000203', 'gzclp----------0000-0000-000000000011', 'gzclp----------0000-0000-000000000111', 1, datetime('now')),
    -- Day 3: T1 Bench, T2 Squat
    ('gzclp----------0000-0000-000000000204', 'gzclp----------0000-0000-000000000012', 'gzclp----------0000-0000-000000000102', 0, datetime('now')),
    ('gzclp----------0000-0000-000000000205', 'gzclp----------0000-0000-000000000012', 'gzclp----------0000-0000-000000000112', 1, datetime('now')),
    -- Day 4: T1 Deadlift, T2 OHP
    ('gzclp----------0000-0000-000000000206', 'gzclp----------0000-0000-000000000013', 'gzclp----------0000-0000-000000000103', 0, datetime('now')),
    ('gzclp----------0000-0000-000000000207', 'gzclp----------0000-0000-000000000013', 'gzclp----------0000-0000-000000000113', 1, datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 7: Link days to week
-- =============================================================================
-- Typical 4-day schedule: Mon, Tue, Thu, Fri (or any 4 days with rest between)
-- +goose StatementBegin
INSERT OR IGNORE INTO week_days (id, week_id, day_id, day_of_week, created_at) VALUES
    ('gzclp----------0000-0000-000000000300', 'gzclp----------0000-0000-000000000003', 'gzclp----------0000-0000-000000000010', 'MONDAY', datetime('now')),
    ('gzclp----------0000-0000-000000000301', 'gzclp----------0000-0000-000000000003', 'gzclp----------0000-0000-000000000011', 'TUESDAY', datetime('now')),
    ('gzclp----------0000-0000-000000000302', 'gzclp----------0000-0000-000000000003', 'gzclp----------0000-0000-000000000012', 'THURSDAY', datetime('now')),
    ('gzclp----------0000-0000-000000000303', 'gzclp----------0000-0000-000000000003', 'gzclp----------0000-0000-000000000013', 'FRIDAY', datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 8: Create program
-- =============================================================================
-- +goose StatementBegin
INSERT OR IGNORE INTO programs (id, name, slug, description, cycle_id, weekly_lookup_id, daily_lookup_id, default_rounding, created_at, updated_at) VALUES
    ('gzclp----------0000-0000-000000000001', 'GZCLP', 'gzclp',
     'GZCL Linear Progression - A tiered training system for novice to early-intermediate lifters. Features T1 (high intensity, 5x3+), T2 (moderate volume, 3x10), and optional T3 accessory work. Uses unique staged progression: on failure, rep scheme changes before weight increases. T1 stages: 5x3+ → 6x2+ → 10x1+ → retest. T2 stages: 3x10 → 3x8 → 3x6 → reset.',
     'gzclp----------0000-0000-000000000002', NULL, NULL, 2.5, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 9: Update days to reference the program
-- =============================================================================
-- +goose StatementBegin
UPDATE days SET program_id = 'gzclp----------0000-0000-000000000001'
WHERE id IN (
    'gzclp----------0000-0000-000000000010',
    'gzclp----------0000-0000-000000000011',
    'gzclp----------0000-0000-000000000012',
    'gzclp----------0000-0000-000000000013'
);
-- +goose StatementEnd

-- =============================================================================
-- STEP 10: Create progression rules
-- =============================================================================
-- GZCLP has unique progression:
-- - T1 Weight: +5lb lower body (squat, deadlift), +2.5lb upper body (bench, OHP)
-- - T2 Weight: +2.5lb for all lifts
-- - On failure: progress through stages (rep scheme changes) before weight increases
--
-- T1 Progression (on success):
--   Lower body: +5lb per successful workout
--   Upper body: +2.5lb per successful workout
--
-- T1 Stage progression (on failure):
--   Stage 1: 5x3+ (15 rep target)
--   Stage 2: 6x2+ (12 rep target) - same weight
--   Stage 3: 10x1+ (10 rep target) - same weight
--   After Stage 3 failure: Retest 5RM, reset to Stage 1
--
-- T2 Progression (on success):
--   All lifts: +2.5lb per successful workout
--
-- T2 Stage progression (on failure):
--   Stage 1: 3x10 (30 rep target)
--   Stage 2: 3x8 (24 rep target) - same weight
--   Stage 3: 3x6 (18 rep target) - same weight
--   After Stage 3 failure: Reset weight, return to Stage 1

-- +goose StatementBegin
INSERT OR IGNORE INTO progressions (id, name, type, parameters, created_at, updated_at) VALUES
    -- T1 Lower body linear progression (+5lb)
    ('gzclp----------0000-0000-000000000400', 'GZCLP T1 Lower +5lb', 'LINEAR_PROGRESSION',
     '{"increment": 5.0, "tier": "T1", "liftCategory": "LOWER"}',
     datetime('now'), datetime('now')),
    -- T1 Upper body linear progression (+2.5lb)
    ('gzclp----------0000-0000-000000000401', 'GZCLP T1 Upper +2.5lb', 'LINEAR_PROGRESSION',
     '{"increment": 2.5, "tier": "T1", "liftCategory": "UPPER"}',
     datetime('now'), datetime('now')),
    -- T2 linear progression (+2.5lb all lifts)
    ('gzclp----------0000-0000-000000000402', 'GZCLP T2 +2.5lb', 'LINEAR_PROGRESSION',
     '{"increment": 2.5, "tier": "T2"}',
     datetime('now'), datetime('now')),
    -- T1 Stage progression (rep scheme changes on failure)
    ('gzclp----------0000-0000-000000000403', 'GZCLP T1 Stage Progression', 'STAGE_PROGRESSION',
     '{"tier": "T1", "stages": [{"stage": 1, "sets": 5, "reps": 3, "isAmrap": true, "targetVolume": 15}, {"stage": 2, "sets": 6, "reps": 2, "isAmrap": true, "targetVolume": 12}, {"stage": 3, "sets": 10, "reps": 1, "isAmrap": true, "targetVolume": 10}], "onFinalFailure": "RETEST_5RM"}',
     datetime('now'), datetime('now')),
    -- T2 Stage progression (rep scheme changes on failure)
    ('gzclp----------0000-0000-000000000404', 'GZCLP T2 Stage Progression', 'STAGE_PROGRESSION',
     '{"tier": "T2", "stages": [{"stage": 1, "sets": 3, "reps": 10, "targetVolume": 30}, {"stage": 2, "sets": 3, "reps": 8, "targetVolume": 24}, {"stage": 3, "sets": 3, "reps": 6, "targetVolume": 18}], "onFinalFailure": "RESET_WEIGHT"}',
     datetime('now'), datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 11: Link progressions to program for each lift
-- =============================================================================
-- Each lift needs both weight progression and stage progression rules
-- +goose StatementBegin
INSERT OR IGNORE INTO program_progressions (id, program_id, progression_id, lift_id, priority, enabled, override_increment, created_at, updated_at) VALUES
    -- Squat T1: +5lb (lower body)
    ('gzclp----------0000-0000-000000000500', 'gzclp----------0000-0000-000000000001', 'gzclp----------0000-0000-000000000400', '00000000-0000-0000-0000-000000000001', 1, 1, NULL, datetime('now'), datetime('now')),
    -- Bench T1: +2.5lb (upper body)
    ('gzclp----------0000-0000-000000000501', 'gzclp----------0000-0000-000000000001', 'gzclp----------0000-0000-000000000401', '00000000-0000-0000-0000-000000000002', 2, 1, NULL, datetime('now'), datetime('now')),
    -- Deadlift T1: +5lb (lower body)
    ('gzclp----------0000-0000-000000000502', 'gzclp----------0000-0000-000000000001', 'gzclp----------0000-0000-000000000400', '00000000-0000-0000-0000-000000000003', 3, 1, NULL, datetime('now'), datetime('now')),
    -- OHP T1: +2.5lb (upper body)
    ('gzclp----------0000-0000-000000000503', 'gzclp----------0000-0000-000000000001', 'gzclp----------0000-0000-000000000401', '00000000-0000-0000-0000-000000000004', 4, 1, NULL, datetime('now'), datetime('now')),
    -- Squat T2: +2.5lb
    ('gzclp----------0000-0000-000000000504', 'gzclp----------0000-0000-000000000001', 'gzclp----------0000-0000-000000000402', '00000000-0000-0000-0000-000000000001', 5, 1, NULL, datetime('now'), datetime('now')),
    -- Bench T2: +2.5lb
    ('gzclp----------0000-0000-000000000505', 'gzclp----------0000-0000-000000000001', 'gzclp----------0000-0000-000000000402', '00000000-0000-0000-0000-000000000002', 6, 1, NULL, datetime('now'), datetime('now')),
    -- Deadlift T2: +2.5lb
    ('gzclp----------0000-0000-000000000506', 'gzclp----------0000-0000-000000000001', 'gzclp----------0000-0000-000000000402', '00000000-0000-0000-0000-000000000003', 7, 1, NULL, datetime('now'), datetime('now')),
    -- OHP T2: +2.5lb
    ('gzclp----------0000-0000-000000000507', 'gzclp----------0000-0000-000000000001', 'gzclp----------0000-0000-000000000402', '00000000-0000-0000-0000-000000000004', 8, 1, NULL, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- +goose Down
-- =============================================================================
-- DOWN MIGRATION: Remove all GZCLP seeded data
-- =============================================================================
-- Order matters due to foreign key constraints - delete in reverse order

-- Remove program progressions
-- +goose StatementBegin
DELETE FROM program_progressions WHERE id IN (
    'gzclp----------0000-0000-000000000500',
    'gzclp----------0000-0000-000000000501',
    'gzclp----------0000-0000-000000000502',
    'gzclp----------0000-0000-000000000503',
    'gzclp----------0000-0000-000000000504',
    'gzclp----------0000-0000-000000000505',
    'gzclp----------0000-0000-000000000506',
    'gzclp----------0000-0000-000000000507'
);
-- +goose StatementEnd

-- Remove progressions
-- +goose StatementBegin
DELETE FROM progressions WHERE id IN (
    'gzclp----------0000-0000-000000000400',
    'gzclp----------0000-0000-000000000401',
    'gzclp----------0000-0000-000000000402',
    'gzclp----------0000-0000-000000000403',
    'gzclp----------0000-0000-000000000404'
);
-- +goose StatementEnd

-- Remove week_days
-- +goose StatementBegin
DELETE FROM week_days WHERE id IN (
    'gzclp----------0000-0000-000000000300',
    'gzclp----------0000-0000-000000000301',
    'gzclp----------0000-0000-000000000302',
    'gzclp----------0000-0000-000000000303'
);
-- +goose StatementEnd

-- Remove day_prescriptions
-- +goose StatementBegin
DELETE FROM day_prescriptions WHERE id IN (
    'gzclp----------0000-0000-000000000200',
    'gzclp----------0000-0000-000000000201',
    'gzclp----------0000-0000-000000000202',
    'gzclp----------0000-0000-000000000203',
    'gzclp----------0000-0000-000000000204',
    'gzclp----------0000-0000-000000000205',
    'gzclp----------0000-0000-000000000206',
    'gzclp----------0000-0000-000000000207'
);
-- +goose StatementEnd

-- Remove prescriptions
-- +goose StatementBegin
DELETE FROM prescriptions WHERE id IN (
    'gzclp----------0000-0000-000000000100',
    'gzclp----------0000-0000-000000000101',
    'gzclp----------0000-0000-000000000102',
    'gzclp----------0000-0000-000000000103',
    'gzclp----------0000-0000-000000000110',
    'gzclp----------0000-0000-000000000111',
    'gzclp----------0000-0000-000000000112',
    'gzclp----------0000-0000-000000000113'
);
-- +goose StatementEnd

-- Remove program (must come before days due to FK)
-- +goose StatementBegin
DELETE FROM programs WHERE id = 'gzclp----------0000-0000-000000000001';
-- +goose StatementEnd

-- Remove days
-- +goose StatementBegin
DELETE FROM days WHERE id IN (
    'gzclp----------0000-0000-000000000010',
    'gzclp----------0000-0000-000000000011',
    'gzclp----------0000-0000-000000000012',
    'gzclp----------0000-0000-000000000013'
);
-- +goose StatementEnd

-- Remove week
-- +goose StatementBegin
DELETE FROM weeks WHERE id = 'gzclp----------0000-0000-000000000003';
-- +goose StatementEnd

-- Remove cycle
-- +goose StatementBegin
DELETE FROM cycles WHERE id = 'gzclp----------0000-0000-000000000002';
-- +goose StatementEnd

-- Note: We do NOT remove lifts in the down migration as they are shared
-- canonical reference data that may be used by other programs or user data.

-- +goose Up
-- Seed Wendler 5/3/1 intermediate program
-- This migration creates the complete 5/3/1 program with its 4-week cycle,
-- percentage-based prescriptions, and AMRAP sets.

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
--   program:        531------------0000-0000-000000000001
--   cycle:          531------------0000-0000-000000000002
--   week-1 (5s):    531------------0000-0000-000000000003
--   week-2 (3s):    531------------0000-0000-000000000004
--   week-3 (531):   531------------0000-0000-000000000005
--   week-4 (deload):531------------0000-0000-000000000006
--
-- Days (shared across weeks):
--   day-press:      531------------0000-0000-000000000010
--   day-deadlift:   531------------0000-0000-000000000011
--   day-bench:      531------------0000-0000-000000000012
--   day-squat:      531------------0000-0000-000000000013
--
-- Week 1 (5s Week) Prescriptions:
--   press-w1-s1:    531------------0000-0000-000000000100
--   press-w1-s2:    531------------0000-0000-000000000101
--   press-w1-s3:    531------------0000-0000-000000000102
--   deadlift-w1-s1: 531------------0000-0000-000000000103
--   deadlift-w1-s2: 531------------0000-0000-000000000104
--   deadlift-w1-s3: 531------------0000-0000-000000000105
--   bench-w1-s1:    531------------0000-0000-000000000106
--   bench-w1-s2:    531------------0000-0000-000000000107
--   bench-w1-s3:    531------------0000-0000-000000000108
--   squat-w1-s1:    531------------0000-0000-000000000109
--   squat-w1-s2:    531------------0000-0000-000000000110
--   squat-w1-s3:    531------------0000-0000-000000000111
--
-- Week 2 (3s Week) Prescriptions:
--   press-w2-s1:    531------------0000-0000-000000000200
--   press-w2-s2:    531------------0000-0000-000000000201
--   press-w2-s3:    531------------0000-0000-000000000202
--   deadlift-w2-s1: 531------------0000-0000-000000000203
--   deadlift-w2-s2: 531------------0000-0000-000000000204
--   deadlift-w2-s3: 531------------0000-0000-000000000205
--   bench-w2-s1:    531------------0000-0000-000000000206
--   bench-w2-s2:    531------------0000-0000-000000000207
--   bench-w2-s3:    531------------0000-0000-000000000208
--   squat-w2-s1:    531------------0000-0000-000000000209
--   squat-w2-s2:    531------------0000-0000-000000000210
--   squat-w2-s3:    531------------0000-0000-000000000211
--
-- Week 3 (5/3/1 Week) Prescriptions:
--   press-w3-s1:    531------------0000-0000-000000000300
--   press-w3-s2:    531------------0000-0000-000000000301
--   press-w3-s3:    531------------0000-0000-000000000302
--   deadlift-w3-s1: 531------------0000-0000-000000000303
--   deadlift-w3-s2: 531------------0000-0000-000000000304
--   deadlift-w3-s3: 531------------0000-0000-000000000305
--   bench-w3-s1:    531------------0000-0000-000000000306
--   bench-w3-s2:    531------------0000-0000-000000000307
--   bench-w3-s3:    531------------0000-0000-000000000308
--   squat-w3-s1:    531------------0000-0000-000000000309
--   squat-w3-s2:    531------------0000-0000-000000000310
--   squat-w3-s3:    531------------0000-0000-000000000311
--
-- Week 4 (Deload) Prescriptions:
--   press-w4-s1:    531------------0000-0000-000000000400
--   press-w4-s2:    531------------0000-0000-000000000401
--   press-w4-s3:    531------------0000-0000-000000000402
--   deadlift-w4-s1: 531------------0000-0000-000000000403
--   deadlift-w4-s2: 531------------0000-0000-000000000404
--   deadlift-w4-s3: 531------------0000-0000-000000000405
--   bench-w4-s1:    531------------0000-0000-000000000406
--   bench-w4-s2:    531------------0000-0000-000000000407
--   bench-w4-s3:    531------------0000-0000-000000000408
--   squat-w4-s1:    531------------0000-0000-000000000409
--   squat-w4-s2:    531------------0000-0000-000000000410
--   squat-w4-s3:    531------------0000-0000-000000000411
--
-- Day prescriptions (join table): 531------------0000-0000-000000000500+
-- Week days (join table): 531------------0000-0000-000000000600+
--
-- Progressions:
--   cycle-5lb:      531------------0000-0000-000000000700
--   cycle-10lb:     531------------0000-0000-000000000701
--
-- Program progressions: 531------------0000-0000-000000000800+

-- =============================================================================
-- STEP 1: Lifts already exist from Starting Strength seed
-- =============================================================================
-- No new lifts needed - all lifts (squat, bench, deadlift, OHP) are already seeded

-- =============================================================================
-- STEP 2: Create cycle (4-week repeating cycle)
-- =============================================================================
-- +goose StatementBegin
INSERT OR IGNORE INTO cycles (id, name, length_weeks, created_at, updated_at) VALUES
    ('531------------0000-0000-000000000002', 'Wendler 5/3/1 Cycle', 4, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 3: Create 4 weeks in the cycle
-- =============================================================================
-- Note: Week variants are stored in day metadata; the variant column is for A/B rotation only.
-- Week 1 = 5s Week, Week 2 = 3s Week, Week 3 = 5/3/1 Week, Week 4 = Deload
-- +goose StatementBegin
INSERT OR IGNORE INTO weeks (id, week_number, variant, cycle_id, created_at, updated_at) VALUES
    ('531------------0000-0000-000000000003', 1, NULL, '531------------0000-0000-000000000002', datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000004', 2, NULL, '531------------0000-0000-000000000002', datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000005', 3, NULL, '531------------0000-0000-000000000002', datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000006', 4, NULL, '531------------0000-0000-000000000002', datetime('now'), datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 4: Create workout days (Press, Deadlift, Bench, Squat)
-- =============================================================================
-- +goose StatementBegin
INSERT OR IGNORE INTO days (id, name, slug, metadata, program_id, created_at, updated_at) VALUES
    ('531------------0000-0000-000000000010', 'Press Day', 'press-day', '{"mainLift": "overhead-press"}', NULL, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000011', 'Deadlift Day', 'deadlift-day', '{"mainLift": "deadlift"}', NULL, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000012', 'Bench Day', 'bench-day', '{"mainLift": "bench-press"}', NULL, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000013', 'Squat Day', 'squat-day', '{"mainLift": "squat"}', NULL, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 5: Create prescriptions for each week
-- =============================================================================
-- Note: 5/3/1 uses Training Max (90% of 1RM) as the reference.
-- All percentages below are of the Training Max.
-- AMRAP sets use is_amrap: true in set_scheme.

-- Week 1 (5s Week): 65% x5, 75% x5, 85% x5+
-- +goose StatementBegin
INSERT OR IGNORE INTO prescriptions (id, lift_id, load_strategy, set_scheme, "order", notes, rest_seconds, created_at, updated_at) VALUES
    -- Press Day Week 1
    ('531------------0000-0000-000000000100', '00000000-0000-0000-0000-000000000004',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 65.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     0, NULL, 180, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000101', '00000000-0000-0000-0000-000000000004',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 75.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     1, NULL, 180, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000102', '00000000-0000-0000-0000-000000000004',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5, "isAmrap": true}',
     2, '5+ AMRAP', 180, datetime('now'), datetime('now')),
    -- Deadlift Day Week 1
    ('531------------0000-0000-000000000103', '00000000-0000-0000-0000-000000000003',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 65.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     0, NULL, 180, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000104', '00000000-0000-0000-0000-000000000003',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 75.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     1, NULL, 180, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000105', '00000000-0000-0000-0000-000000000003',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5, "isAmrap": true}',
     2, '5+ AMRAP', 180, datetime('now'), datetime('now')),
    -- Bench Day Week 1
    ('531------------0000-0000-000000000106', '00000000-0000-0000-0000-000000000002',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 65.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     0, NULL, 180, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000107', '00000000-0000-0000-0000-000000000002',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 75.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     1, NULL, 180, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000108', '00000000-0000-0000-0000-000000000002',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5, "isAmrap": true}',
     2, '5+ AMRAP', 180, datetime('now'), datetime('now')),
    -- Squat Day Week 1
    ('531------------0000-0000-000000000109', '00000000-0000-0000-0000-000000000001',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 65.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     0, NULL, 180, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000110', '00000000-0000-0000-0000-000000000001',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 75.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     1, NULL, 180, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000111', '00000000-0000-0000-0000-000000000001',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5, "isAmrap": true}',
     2, '5+ AMRAP', 180, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- Week 2 (3s Week): 70% x3, 80% x3, 90% x3+
-- +goose StatementBegin
INSERT OR IGNORE INTO prescriptions (id, lift_id, load_strategy, set_scheme, "order", notes, rest_seconds, created_at, updated_at) VALUES
    -- Press Day Week 2
    ('531------------0000-0000-000000000200', '00000000-0000-0000-0000-000000000004',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 70.0}',
     '{"type": "FIXED", "sets": 1, "reps": 3}',
     0, NULL, 180, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000201', '00000000-0000-0000-0000-000000000004',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 80.0}',
     '{"type": "FIXED", "sets": 1, "reps": 3}',
     1, NULL, 180, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000202', '00000000-0000-0000-0000-000000000004',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 90.0}',
     '{"type": "FIXED", "sets": 1, "reps": 3, "isAmrap": true}',
     2, '3+ AMRAP', 180, datetime('now'), datetime('now')),
    -- Deadlift Day Week 2
    ('531------------0000-0000-000000000203', '00000000-0000-0000-0000-000000000003',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 70.0}',
     '{"type": "FIXED", "sets": 1, "reps": 3}',
     0, NULL, 180, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000204', '00000000-0000-0000-0000-000000000003',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 80.0}',
     '{"type": "FIXED", "sets": 1, "reps": 3}',
     1, NULL, 180, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000205', '00000000-0000-0000-0000-000000000003',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 90.0}',
     '{"type": "FIXED", "sets": 1, "reps": 3, "isAmrap": true}',
     2, '3+ AMRAP', 180, datetime('now'), datetime('now')),
    -- Bench Day Week 2
    ('531------------0000-0000-000000000206', '00000000-0000-0000-0000-000000000002',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 70.0}',
     '{"type": "FIXED", "sets": 1, "reps": 3}',
     0, NULL, 180, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000207', '00000000-0000-0000-0000-000000000002',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 80.0}',
     '{"type": "FIXED", "sets": 1, "reps": 3}',
     1, NULL, 180, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000208', '00000000-0000-0000-0000-000000000002',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 90.0}',
     '{"type": "FIXED", "sets": 1, "reps": 3, "isAmrap": true}',
     2, '3+ AMRAP', 180, datetime('now'), datetime('now')),
    -- Squat Day Week 2
    ('531------------0000-0000-000000000209', '00000000-0000-0000-0000-000000000001',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 70.0}',
     '{"type": "FIXED", "sets": 1, "reps": 3}',
     0, NULL, 180, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000210', '00000000-0000-0000-0000-000000000001',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 80.0}',
     '{"type": "FIXED", "sets": 1, "reps": 3}',
     1, NULL, 180, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000211', '00000000-0000-0000-0000-000000000001',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 90.0}',
     '{"type": "FIXED", "sets": 1, "reps": 3, "isAmrap": true}',
     2, '3+ AMRAP', 180, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- Week 3 (5/3/1 Week): 75% x5, 85% x3, 95% x1+
-- +goose StatementBegin
INSERT OR IGNORE INTO prescriptions (id, lift_id, load_strategy, set_scheme, "order", notes, rest_seconds, created_at, updated_at) VALUES
    -- Press Day Week 3
    ('531------------0000-0000-000000000300', '00000000-0000-0000-0000-000000000004',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 75.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     0, NULL, 180, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000301', '00000000-0000-0000-0000-000000000004',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85.0}',
     '{"type": "FIXED", "sets": 1, "reps": 3}',
     1, NULL, 180, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000302', '00000000-0000-0000-0000-000000000004',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 95.0}',
     '{"type": "FIXED", "sets": 1, "reps": 1, "isAmrap": true}',
     2, '1+ AMRAP', 180, datetime('now'), datetime('now')),
    -- Deadlift Day Week 3
    ('531------------0000-0000-000000000303', '00000000-0000-0000-0000-000000000003',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 75.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     0, NULL, 180, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000304', '00000000-0000-0000-0000-000000000003',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85.0}',
     '{"type": "FIXED", "sets": 1, "reps": 3}',
     1, NULL, 180, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000305', '00000000-0000-0000-0000-000000000003',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 95.0}',
     '{"type": "FIXED", "sets": 1, "reps": 1, "isAmrap": true}',
     2, '1+ AMRAP', 180, datetime('now'), datetime('now')),
    -- Bench Day Week 3
    ('531------------0000-0000-000000000306', '00000000-0000-0000-0000-000000000002',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 75.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     0, NULL, 180, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000307', '00000000-0000-0000-0000-000000000002',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85.0}',
     '{"type": "FIXED", "sets": 1, "reps": 3}',
     1, NULL, 180, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000308', '00000000-0000-0000-0000-000000000002',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 95.0}',
     '{"type": "FIXED", "sets": 1, "reps": 1, "isAmrap": true}',
     2, '1+ AMRAP', 180, datetime('now'), datetime('now')),
    -- Squat Day Week 3
    ('531------------0000-0000-000000000309', '00000000-0000-0000-0000-000000000001',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 75.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     0, NULL, 180, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000310', '00000000-0000-0000-0000-000000000001',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 85.0}',
     '{"type": "FIXED", "sets": 1, "reps": 3}',
     1, NULL, 180, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000311', '00000000-0000-0000-0000-000000000001',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 95.0}',
     '{"type": "FIXED", "sets": 1, "reps": 1, "isAmrap": true}',
     2, '1+ AMRAP', 180, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- Week 4 (Deload): 40% x5, 50% x5, 60% x5
-- +goose StatementBegin
INSERT OR IGNORE INTO prescriptions (id, lift_id, load_strategy, set_scheme, "order", notes, rest_seconds, created_at, updated_at) VALUES
    -- Press Day Week 4
    ('531------------0000-0000-000000000400', '00000000-0000-0000-0000-000000000004',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 40.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     0, 'Deload', 120, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000401', '00000000-0000-0000-0000-000000000004',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 50.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     1, 'Deload', 120, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000402', '00000000-0000-0000-0000-000000000004',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 60.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     2, 'Deload', 120, datetime('now'), datetime('now')),
    -- Deadlift Day Week 4
    ('531------------0000-0000-000000000403', '00000000-0000-0000-0000-000000000003',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 40.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     0, 'Deload', 120, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000404', '00000000-0000-0000-0000-000000000003',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 50.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     1, 'Deload', 120, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000405', '00000000-0000-0000-0000-000000000003',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 60.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     2, 'Deload', 120, datetime('now'), datetime('now')),
    -- Bench Day Week 4
    ('531------------0000-0000-000000000406', '00000000-0000-0000-0000-000000000002',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 40.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     0, 'Deload', 120, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000407', '00000000-0000-0000-0000-000000000002',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 50.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     1, 'Deload', 120, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000408', '00000000-0000-0000-0000-000000000002',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 60.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     2, 'Deload', 120, datetime('now'), datetime('now')),
    -- Squat Day Week 4
    ('531------------0000-0000-000000000409', '00000000-0000-0000-0000-000000000001',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 40.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     0, 'Deload', 120, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000410', '00000000-0000-0000-0000-000000000001',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 50.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     1, 'Deload', 120, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000411', '00000000-0000-0000-0000-000000000001',
     '{"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 60.0}',
     '{"type": "FIXED", "sets": 1, "reps": 5}',
     2, 'Deload', 120, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 6: Link prescriptions to days via day_prescriptions
-- =============================================================================
-- Since prescriptions vary by week, we need to link them week-specifically.
-- We'll create separate day instances per week to handle different prescriptions.

-- Actually, looking at the schema more carefully, we need week-specific days.
-- Let me create week-specific day prescription links using a different approach.
-- The prescriptions table has order field that works within a day context.
-- We need to create week_day_prescriptions that tie week -> day -> prescriptions.

-- For 5/3/1, each week has different prescriptions for the same conceptual day.
-- The cleanest approach is to have the days be generic templates and use
-- a week_day_prescriptions join table, but looking at the existing schema,
-- we're using day_prescriptions without week context.

-- Let me re-examine: in 5/3/1, the SAME day (e.g., "Press Day") has DIFFERENT
-- prescriptions in week 1 vs week 2 vs week 3 vs week 4.

-- The existing schema ties prescriptions to days, and days to weeks.
-- For week-varying prescriptions, we need to create week-specific day instances.

-- Create week-specific day instances
-- +goose StatementBegin
INSERT OR IGNORE INTO days (id, name, slug, metadata, program_id, created_at, updated_at) VALUES
    -- Week 1 days
    ('531------------0000-0000-000000000020', 'Press Day - 5s Week', 'press-day-w1', '{"mainLift": "overhead-press", "week": 1}', NULL, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000021', 'Deadlift Day - 5s Week', 'deadlift-day-w1', '{"mainLift": "deadlift", "week": 1}', NULL, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000022', 'Bench Day - 5s Week', 'bench-day-w1', '{"mainLift": "bench-press", "week": 1}', NULL, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000023', 'Squat Day - 5s Week', 'squat-day-w1', '{"mainLift": "squat", "week": 1}', NULL, datetime('now'), datetime('now')),
    -- Week 2 days
    ('531------------0000-0000-000000000024', 'Press Day - 3s Week', 'press-day-w2', '{"mainLift": "overhead-press", "week": 2}', NULL, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000025', 'Deadlift Day - 3s Week', 'deadlift-day-w2', '{"mainLift": "deadlift", "week": 2}', NULL, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000026', 'Bench Day - 3s Week', 'bench-day-w2', '{"mainLift": "bench-press", "week": 2}', NULL, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000027', 'Squat Day - 3s Week', 'squat-day-w2', '{"mainLift": "squat", "week": 2}', NULL, datetime('now'), datetime('now')),
    -- Week 3 days
    ('531------------0000-0000-000000000028', 'Press Day - 5/3/1 Week', 'press-day-w3', '{"mainLift": "overhead-press", "week": 3}', NULL, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000029', 'Deadlift Day - 5/3/1 Week', 'deadlift-day-w3', '{"mainLift": "deadlift", "week": 3}', NULL, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000030', 'Bench Day - 5/3/1 Week', 'bench-day-w3', '{"mainLift": "bench-press", "week": 3}', NULL, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000031', 'Squat Day - 5/3/1 Week', 'squat-day-w3', '{"mainLift": "squat", "week": 3}', NULL, datetime('now'), datetime('now')),
    -- Week 4 days (Deload)
    ('531------------0000-0000-000000000032', 'Press Day - Deload', 'press-day-w4', '{"mainLift": "overhead-press", "week": 4}', NULL, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000033', 'Deadlift Day - Deload', 'deadlift-day-w4', '{"mainLift": "deadlift", "week": 4}', NULL, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000034', 'Bench Day - Deload', 'bench-day-w4', '{"mainLift": "bench-press", "week": 4}', NULL, datetime('now'), datetime('now')),
    ('531------------0000-0000-000000000035', 'Squat Day - Deload', 'squat-day-w4', '{"mainLift": "squat", "week": 4}', NULL, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- Delete the original generic days we don't need
-- +goose StatementBegin
DELETE FROM days WHERE id IN (
    '531------------0000-0000-000000000010',
    '531------------0000-0000-000000000011',
    '531------------0000-0000-000000000012',
    '531------------0000-0000-000000000013'
);
-- +goose StatementEnd

-- Link Week 1 prescriptions to Week 1 days
-- +goose StatementBegin
INSERT OR IGNORE INTO day_prescriptions (id, day_id, prescription_id, "order", created_at) VALUES
    -- Press Day Week 1
    ('531------------0000-0000-000000000500', '531------------0000-0000-000000000020', '531------------0000-0000-000000000100', 0, datetime('now')),
    ('531------------0000-0000-000000000501', '531------------0000-0000-000000000020', '531------------0000-0000-000000000101', 1, datetime('now')),
    ('531------------0000-0000-000000000502', '531------------0000-0000-000000000020', '531------------0000-0000-000000000102', 2, datetime('now')),
    -- Deadlift Day Week 1
    ('531------------0000-0000-000000000503', '531------------0000-0000-000000000021', '531------------0000-0000-000000000103', 0, datetime('now')),
    ('531------------0000-0000-000000000504', '531------------0000-0000-000000000021', '531------------0000-0000-000000000104', 1, datetime('now')),
    ('531------------0000-0000-000000000505', '531------------0000-0000-000000000021', '531------------0000-0000-000000000105', 2, datetime('now')),
    -- Bench Day Week 1
    ('531------------0000-0000-000000000506', '531------------0000-0000-000000000022', '531------------0000-0000-000000000106', 0, datetime('now')),
    ('531------------0000-0000-000000000507', '531------------0000-0000-000000000022', '531------------0000-0000-000000000107', 1, datetime('now')),
    ('531------------0000-0000-000000000508', '531------------0000-0000-000000000022', '531------------0000-0000-000000000108', 2, datetime('now')),
    -- Squat Day Week 1
    ('531------------0000-0000-000000000509', '531------------0000-0000-000000000023', '531------------0000-0000-000000000109', 0, datetime('now')),
    ('531------------0000-0000-000000000510', '531------------0000-0000-000000000023', '531------------0000-0000-000000000110', 1, datetime('now')),
    ('531------------0000-0000-000000000511', '531------------0000-0000-000000000023', '531------------0000-0000-000000000111', 2, datetime('now'));
-- +goose StatementEnd

-- Link Week 2 prescriptions to Week 2 days
-- +goose StatementBegin
INSERT OR IGNORE INTO day_prescriptions (id, day_id, prescription_id, "order", created_at) VALUES
    -- Press Day Week 2
    ('531------------0000-0000-000000000512', '531------------0000-0000-000000000024', '531------------0000-0000-000000000200', 0, datetime('now')),
    ('531------------0000-0000-000000000513', '531------------0000-0000-000000000024', '531------------0000-0000-000000000201', 1, datetime('now')),
    ('531------------0000-0000-000000000514', '531------------0000-0000-000000000024', '531------------0000-0000-000000000202', 2, datetime('now')),
    -- Deadlift Day Week 2
    ('531------------0000-0000-000000000515', '531------------0000-0000-000000000025', '531------------0000-0000-000000000203', 0, datetime('now')),
    ('531------------0000-0000-000000000516', '531------------0000-0000-000000000025', '531------------0000-0000-000000000204', 1, datetime('now')),
    ('531------------0000-0000-000000000517', '531------------0000-0000-000000000025', '531------------0000-0000-000000000205', 2, datetime('now')),
    -- Bench Day Week 2
    ('531------------0000-0000-000000000518', '531------------0000-0000-000000000026', '531------------0000-0000-000000000206', 0, datetime('now')),
    ('531------------0000-0000-000000000519', '531------------0000-0000-000000000026', '531------------0000-0000-000000000207', 1, datetime('now')),
    ('531------------0000-0000-000000000520', '531------------0000-0000-000000000026', '531------------0000-0000-000000000208', 2, datetime('now')),
    -- Squat Day Week 2
    ('531------------0000-0000-000000000521', '531------------0000-0000-000000000027', '531------------0000-0000-000000000209', 0, datetime('now')),
    ('531------------0000-0000-000000000522', '531------------0000-0000-000000000027', '531------------0000-0000-000000000210', 1, datetime('now')),
    ('531------------0000-0000-000000000523', '531------------0000-0000-000000000027', '531------------0000-0000-000000000211', 2, datetime('now'));
-- +goose StatementEnd

-- Link Week 3 prescriptions to Week 3 days
-- +goose StatementBegin
INSERT OR IGNORE INTO day_prescriptions (id, day_id, prescription_id, "order", created_at) VALUES
    -- Press Day Week 3
    ('531------------0000-0000-000000000524', '531------------0000-0000-000000000028', '531------------0000-0000-000000000300', 0, datetime('now')),
    ('531------------0000-0000-000000000525', '531------------0000-0000-000000000028', '531------------0000-0000-000000000301', 1, datetime('now')),
    ('531------------0000-0000-000000000526', '531------------0000-0000-000000000028', '531------------0000-0000-000000000302', 2, datetime('now')),
    -- Deadlift Day Week 3
    ('531------------0000-0000-000000000527', '531------------0000-0000-000000000029', '531------------0000-0000-000000000303', 0, datetime('now')),
    ('531------------0000-0000-000000000528', '531------------0000-0000-000000000029', '531------------0000-0000-000000000304', 1, datetime('now')),
    ('531------------0000-0000-000000000529', '531------------0000-0000-000000000029', '531------------0000-0000-000000000305', 2, datetime('now')),
    -- Bench Day Week 3
    ('531------------0000-0000-000000000530', '531------------0000-0000-000000000030', '531------------0000-0000-000000000306', 0, datetime('now')),
    ('531------------0000-0000-000000000531', '531------------0000-0000-000000000030', '531------------0000-0000-000000000307', 1, datetime('now')),
    ('531------------0000-0000-000000000532', '531------------0000-0000-000000000030', '531------------0000-0000-000000000308', 2, datetime('now')),
    -- Squat Day Week 3
    ('531------------0000-0000-000000000533', '531------------0000-0000-000000000031', '531------------0000-0000-000000000309', 0, datetime('now')),
    ('531------------0000-0000-000000000534', '531------------0000-0000-000000000031', '531------------0000-0000-000000000310', 1, datetime('now')),
    ('531------------0000-0000-000000000535', '531------------0000-0000-000000000031', '531------------0000-0000-000000000311', 2, datetime('now'));
-- +goose StatementEnd

-- Link Week 4 prescriptions to Week 4 days
-- +goose StatementBegin
INSERT OR IGNORE INTO day_prescriptions (id, day_id, prescription_id, "order", created_at) VALUES
    -- Press Day Week 4
    ('531------------0000-0000-000000000536', '531------------0000-0000-000000000032', '531------------0000-0000-000000000400', 0, datetime('now')),
    ('531------------0000-0000-000000000537', '531------------0000-0000-000000000032', '531------------0000-0000-000000000401', 1, datetime('now')),
    ('531------------0000-0000-000000000538', '531------------0000-0000-000000000032', '531------------0000-0000-000000000402', 2, datetime('now')),
    -- Deadlift Day Week 4
    ('531------------0000-0000-000000000539', '531------------0000-0000-000000000033', '531------------0000-0000-000000000403', 0, datetime('now')),
    ('531------------0000-0000-000000000540', '531------------0000-0000-000000000033', '531------------0000-0000-000000000404', 1, datetime('now')),
    ('531------------0000-0000-000000000541', '531------------0000-0000-000000000033', '531------------0000-0000-000000000405', 2, datetime('now')),
    -- Bench Day Week 4
    ('531------------0000-0000-000000000542', '531------------0000-0000-000000000034', '531------------0000-0000-000000000406', 0, datetime('now')),
    ('531------------0000-0000-000000000543', '531------------0000-0000-000000000034', '531------------0000-0000-000000000407', 1, datetime('now')),
    ('531------------0000-0000-000000000544', '531------------0000-0000-000000000034', '531------------0000-0000-000000000408', 2, datetime('now')),
    -- Squat Day Week 4
    ('531------------0000-0000-000000000545', '531------------0000-0000-000000000035', '531------------0000-0000-000000000409', 0, datetime('now')),
    ('531------------0000-0000-000000000546', '531------------0000-0000-000000000035', '531------------0000-0000-000000000410', 1, datetime('now')),
    ('531------------0000-0000-000000000547', '531------------0000-0000-000000000035', '531------------0000-0000-000000000411', 2, datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 7: Link days to weeks
-- =============================================================================
-- 5/3/1 has 4 training days per week on a 4-day schedule
-- Typically: Day 1 (Mon), Day 2 (Tue), Day 3 (Thu), Day 4 (Fri)
-- +goose StatementBegin
INSERT OR IGNORE INTO week_days (id, week_id, day_id, day_of_week, created_at) VALUES
    -- Week 1 (5s Week)
    ('531------------0000-0000-000000000600', '531------------0000-0000-000000000003', '531------------0000-0000-000000000020', 'MONDAY', datetime('now')),
    ('531------------0000-0000-000000000601', '531------------0000-0000-000000000003', '531------------0000-0000-000000000021', 'TUESDAY', datetime('now')),
    ('531------------0000-0000-000000000602', '531------------0000-0000-000000000003', '531------------0000-0000-000000000022', 'THURSDAY', datetime('now')),
    ('531------------0000-0000-000000000603', '531------------0000-0000-000000000003', '531------------0000-0000-000000000023', 'FRIDAY', datetime('now')),
    -- Week 2 (3s Week)
    ('531------------0000-0000-000000000604', '531------------0000-0000-000000000004', '531------------0000-0000-000000000024', 'MONDAY', datetime('now')),
    ('531------------0000-0000-000000000605', '531------------0000-0000-000000000004', '531------------0000-0000-000000000025', 'TUESDAY', datetime('now')),
    ('531------------0000-0000-000000000606', '531------------0000-0000-000000000004', '531------------0000-0000-000000000026', 'THURSDAY', datetime('now')),
    ('531------------0000-0000-000000000607', '531------------0000-0000-000000000004', '531------------0000-0000-000000000027', 'FRIDAY', datetime('now')),
    -- Week 3 (5/3/1 Week)
    ('531------------0000-0000-000000000608', '531------------0000-0000-000000000005', '531------------0000-0000-000000000028', 'MONDAY', datetime('now')),
    ('531------------0000-0000-000000000609', '531------------0000-0000-000000000005', '531------------0000-0000-000000000029', 'TUESDAY', datetime('now')),
    ('531------------0000-0000-000000000610', '531------------0000-0000-000000000005', '531------------0000-0000-000000000030', 'THURSDAY', datetime('now')),
    ('531------------0000-0000-000000000611', '531------------0000-0000-000000000005', '531------------0000-0000-000000000031', 'FRIDAY', datetime('now')),
    -- Week 4 (Deload)
    ('531------------0000-0000-000000000612', '531------------0000-0000-000000000006', '531------------0000-0000-000000000032', 'MONDAY', datetime('now')),
    ('531------------0000-0000-000000000613', '531------------0000-0000-000000000006', '531------------0000-0000-000000000033', 'TUESDAY', datetime('now')),
    ('531------------0000-0000-000000000614', '531------------0000-0000-000000000006', '531------------0000-0000-000000000034', 'THURSDAY', datetime('now')),
    ('531------------0000-0000-000000000615', '531------------0000-0000-000000000006', '531------------0000-0000-000000000035', 'FRIDAY', datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 8: Create program
-- =============================================================================
-- +goose StatementBegin
INSERT OR IGNORE INTO programs (id, name, slug, description, cycle_id, weekly_lookup_id, daily_lookup_id, default_rounding, created_at, updated_at) VALUES
    ('531------------0000-0000-000000000001', 'Wendler 5/3/1', '531',
     'Jim Wendler''s percentage-based intermediate strength program. Features a 4-week cycle with progressively heavier loads culminating in AMRAP sets. Training max is 90% of 1RM. Includes built-in deload week for recovery.',
     '531------------0000-0000-000000000002', NULL, NULL, 5.0, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 9: Update days to reference the program
-- =============================================================================
-- +goose StatementBegin
UPDATE days SET program_id = '531------------0000-0000-000000000001'
WHERE id LIKE '531------------0000-0000-00000000002%'
   OR id LIKE '531------------0000-0000-00000000003%';
-- +goose StatementEnd

-- =============================================================================
-- STEP 10: Create progression rules
-- =============================================================================
-- 5/3/1 uses cycle-based progression: +5lb for upper body, +10lb for lower body per cycle
-- +goose StatementBegin
INSERT OR IGNORE INTO progressions (id, name, type, parameters, created_at, updated_at) VALUES
    -- +5lb per cycle (upper body: bench, press)
    ('531------------0000-0000-000000000700', '5/3/1 Cycle +5lb', 'LINEAR_PROGRESSION',
     '{"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_CYCLE"}',
     datetime('now'), datetime('now')),
    -- +10lb per cycle (lower body: squat, deadlift)
    ('531------------0000-0000-000000000701', '5/3/1 Cycle +10lb', 'LINEAR_PROGRESSION',
     '{"increment": 10.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_CYCLE"}',
     datetime('now'), datetime('now'));
-- +goose StatementEnd

-- =============================================================================
-- STEP 11: Link progressions to program for each lift
-- =============================================================================
-- +goose StatementBegin
INSERT OR IGNORE INTO program_progressions (id, program_id, progression_id, lift_id, priority, enabled, override_increment, created_at, updated_at) VALUES
    -- Squat: +10lb per cycle (lower body)
    ('531------------0000-0000-000000000800', '531------------0000-0000-000000000001', '531------------0000-0000-000000000701', '00000000-0000-0000-0000-000000000001', 1, 1, NULL, datetime('now'), datetime('now')),
    -- Bench Press: +5lb per cycle (upper body)
    ('531------------0000-0000-000000000801', '531------------0000-0000-000000000001', '531------------0000-0000-000000000700', '00000000-0000-0000-0000-000000000002', 2, 1, NULL, datetime('now'), datetime('now')),
    -- Deadlift: +10lb per cycle (lower body)
    ('531------------0000-0000-000000000802', '531------------0000-0000-000000000001', '531------------0000-0000-000000000701', '00000000-0000-0000-0000-000000000003', 3, 1, NULL, datetime('now'), datetime('now')),
    -- Overhead Press: +5lb per cycle (upper body)
    ('531------------0000-0000-000000000803', '531------------0000-0000-000000000001', '531------------0000-0000-000000000700', '00000000-0000-0000-0000-000000000004', 4, 1, NULL, datetime('now'), datetime('now'));
-- +goose StatementEnd

-- +goose Down
-- =============================================================================
-- DOWN MIGRATION: Remove all 5/3/1 seeded data
-- =============================================================================
-- Order matters due to foreign key constraints - delete in reverse order

-- Remove program progressions
-- +goose StatementBegin
DELETE FROM program_progressions WHERE id LIKE '531------------0000-0000-00000000080%';
-- +goose StatementEnd

-- Remove progressions
-- +goose StatementBegin
DELETE FROM progressions WHERE id IN (
    '531------------0000-0000-000000000700',
    '531------------0000-0000-000000000701'
);
-- +goose StatementEnd

-- Remove week_days
-- +goose StatementBegin
DELETE FROM week_days WHERE id LIKE '531------------0000-0000-00000000060%';
-- +goose StatementEnd

-- Remove day_prescriptions
-- +goose StatementBegin
DELETE FROM day_prescriptions WHERE id LIKE '531------------0000-0000-00000000050%';
-- +goose StatementEnd

-- Remove prescriptions
-- +goose StatementBegin
DELETE FROM prescriptions WHERE id LIKE '531------------0000-0000-0000000001%'
    OR id LIKE '531------------0000-0000-0000000002%'
    OR id LIKE '531------------0000-0000-0000000003%'
    OR id LIKE '531------------0000-0000-0000000004%';
-- +goose StatementEnd

-- Remove program (must come before days due to FK)
-- +goose StatementBegin
DELETE FROM programs WHERE id = '531------------0000-0000-000000000001';
-- +goose StatementEnd

-- Remove days
-- +goose StatementBegin
DELETE FROM days WHERE id LIKE '531------------0000-0000-00000000002%'
    OR id LIKE '531------------0000-0000-00000000003%';
-- +goose StatementEnd

-- Remove weeks
-- +goose StatementBegin
DELETE FROM weeks WHERE id IN (
    '531------------0000-0000-000000000003',
    '531------------0000-0000-000000000004',
    '531------------0000-0000-000000000005',
    '531------------0000-0000-000000000006'
);
-- +goose StatementEnd

-- Remove cycle
-- +goose StatementBegin
DELETE FROM cycles WHERE id = '531------------0000-0000-000000000002';
-- +goose StatementEnd

-- Note: We do NOT remove lifts in the down migration as they are shared
-- canonical reference data that may be used by other programs or user data.

-- +goose Up
-- Backfill discovery metadata for the four canonical programs
-- This migration updates existing programs with accurate metadata values
-- for filtering and discovery purposes.

-- Starting Strength: Beginner, 3 days/week, strength focus, no AMRAP
-- +goose StatementBegin
UPDATE programs SET
    difficulty = 'beginner',
    days_per_week = 3,
    focus = 'strength',
    has_amrap = 0
WHERE slug = 'starting-strength';
-- +goose StatementEnd

-- Texas Method: Intermediate, 3 days/week, strength focus, no AMRAP
-- +goose StatementBegin
UPDATE programs SET
    difficulty = 'intermediate',
    days_per_week = 3,
    focus = 'strength',
    has_amrap = 0
WHERE slug = 'texas-method';
-- +goose StatementEnd

-- Wendler 5/3/1: Intermediate, 4 days/week, strength focus, has AMRAP
-- +goose StatementBegin
UPDATE programs SET
    difficulty = 'intermediate',
    days_per_week = 4,
    focus = 'strength',
    has_amrap = 1
WHERE slug = '531';
-- +goose StatementEnd

-- GZCLP: Beginner, 4 days/week, strength focus, has AMRAP
-- +goose StatementBegin
UPDATE programs SET
    difficulty = 'beginner',
    days_per_week = 4,
    focus = 'strength',
    has_amrap = 1
WHERE slug = 'gzclp';
-- +goose StatementEnd

-- +goose Down
-- Reset metadata to column defaults
-- Note: This is a no-op if programs don't exist

-- +goose StatementBegin
UPDATE programs SET
    difficulty = 'beginner',
    days_per_week = 3,
    focus = 'strength',
    has_amrap = 0
WHERE slug IN ('starting-strength', 'texas-method', '531', 'gzclp');
-- +goose StatementEnd

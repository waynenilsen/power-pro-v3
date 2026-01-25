-- +goose Up
-- Add program discovery metadata columns to the programs table
-- SQLite requires separate ALTER TABLE statements for each column

-- Add difficulty column
-- +goose StatementBegin
ALTER TABLE programs ADD COLUMN difficulty TEXT NOT NULL DEFAULT 'beginner' CHECK(difficulty IN ('beginner', 'intermediate', 'advanced'));
-- +goose StatementEnd

-- Add days_per_week column
-- +goose StatementBegin
ALTER TABLE programs ADD COLUMN days_per_week INTEGER NOT NULL DEFAULT 3 CHECK(days_per_week BETWEEN 1 AND 7);
-- +goose StatementEnd

-- Add focus column
-- +goose StatementBegin
ALTER TABLE programs ADD COLUMN focus TEXT NOT NULL DEFAULT 'strength' CHECK(focus IN ('strength', 'hypertrophy', 'peaking'));
-- +goose StatementEnd

-- Add has_amrap column
-- +goose StatementBegin
ALTER TABLE programs ADD COLUMN has_amrap INTEGER NOT NULL DEFAULT 0 CHECK(has_amrap IN (0, 1));
-- +goose StatementEnd

-- Create indexes for filtering
-- +goose StatementBegin
CREATE INDEX idx_programs_difficulty ON programs(difficulty);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_programs_days_per_week ON programs(days_per_week);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_programs_focus ON programs(focus);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_programs_has_amrap ON programs(has_amrap);
-- +goose StatementEnd

-- +goose Down
-- Drop indexes first
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_programs_has_amrap;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_programs_focus;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_programs_days_per_week;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_programs_difficulty;
-- +goose StatementEnd

-- SQLite does not support DROP COLUMN directly in older versions
-- We need to recreate the table without the columns
-- +goose StatementBegin
CREATE TABLE programs_backup (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK(length(name) > 0 AND length(name) <= 100),
    slug TEXT NOT NULL UNIQUE CHECK(length(slug) > 0 AND length(slug) <= 100 AND slug GLOB '[a-z0-9-]*'),
    description TEXT,
    cycle_id TEXT NOT NULL,
    weekly_lookup_id TEXT,
    daily_lookup_id TEXT,
    default_rounding REAL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (cycle_id) REFERENCES cycles(id) ON DELETE RESTRICT,
    FOREIGN KEY (weekly_lookup_id) REFERENCES weekly_lookups(id) ON DELETE SET NULL,
    FOREIGN KEY (daily_lookup_id) REFERENCES daily_lookups(id) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO programs_backup (id, name, slug, description, cycle_id, weekly_lookup_id, daily_lookup_id, default_rounding, created_at, updated_at)
SELECT id, name, slug, description, cycle_id, weekly_lookup_id, daily_lookup_id, default_rounding, created_at, updated_at
FROM programs;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE programs;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE programs_backup RENAME TO programs;
-- +goose StatementEnd

-- Recreate existing indexes
-- +goose StatementBegin
CREATE INDEX idx_programs_name ON programs(name);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_programs_cycle_id ON programs(cycle_id);
-- +goose StatementEnd

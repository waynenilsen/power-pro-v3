-- +goose Up
-- +goose StatementBegin
CREATE TABLE cycles (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK(length(name) > 0 AND length(name) <= 100),
    length_weeks INTEGER NOT NULL CHECK(length_weeks >= 1),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
-- +goose StatementEnd

-- Index for name lookups
-- +goose StatementBegin
CREATE INDEX idx_cycles_name ON cycles(name);
-- +goose StatementEnd

-- Add foreign key constraint to weeks table
-- +goose StatementBegin
CREATE TABLE weeks_new (
    id TEXT PRIMARY KEY,
    week_number INTEGER NOT NULL CHECK(week_number >= 1),
    variant TEXT CHECK(variant IS NULL OR variant IN ('A', 'B')),
    cycle_id TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (cycle_id) REFERENCES cycles(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- Copy data from old weeks table to new table
-- +goose StatementBegin
INSERT INTO weeks_new SELECT * FROM weeks;
-- +goose StatementEnd

-- Drop old weeks table
-- +goose StatementBegin
DROP TABLE weeks;
-- +goose StatementEnd

-- Rename new table to weeks
-- +goose StatementBegin
ALTER TABLE weeks_new RENAME TO weeks;
-- +goose StatementEnd

-- Recreate indexes on weeks table
-- +goose StatementBegin
CREATE UNIQUE INDEX idx_weeks_cycle_id_week_number ON weeks(cycle_id, week_number);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_weeks_cycle_id ON weeks(cycle_id);
-- +goose StatementEnd

-- +goose Down
-- Remove foreign key constraint from weeks by recreating without it
-- +goose StatementBegin
CREATE TABLE weeks_old (
    id TEXT PRIMARY KEY,
    week_number INTEGER NOT NULL CHECK(week_number >= 1),
    variant TEXT CHECK(variant IS NULL OR variant IN ('A', 'B')),
    cycle_id TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO weeks_old SELECT * FROM weeks;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE weeks;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE weeks_old RENAME TO weeks;
-- +goose StatementEnd

-- Recreate indexes on weeks table
-- +goose StatementBegin
CREATE UNIQUE INDEX idx_weeks_cycle_id_week_number ON weeks(cycle_id, week_number);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_weeks_cycle_id ON weeks(cycle_id);
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_cycles_name;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS cycles;
-- +goose StatementEnd

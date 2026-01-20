-- +goose Up
-- +goose StatementBegin
CREATE TABLE programs (
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

-- Index for name lookups
-- +goose StatementBegin
CREATE INDEX idx_programs_name ON programs(name);
-- +goose StatementEnd

-- Index for cycle_id lookups
-- +goose StatementBegin
CREATE INDEX idx_programs_cycle_id ON programs(cycle_id);
-- +goose StatementEnd

-- Add foreign key constraints to lookup tables for program_id
-- SQLite requires table recreation to add foreign keys

-- Recreate weekly_lookups with foreign key
-- +goose StatementBegin
CREATE TABLE weekly_lookups_new (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK(length(name) > 0 AND length(name) <= 100),
    entries TEXT NOT NULL CHECK(json_valid(entries)),
    program_id TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (program_id) REFERENCES programs(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO weekly_lookups_new SELECT * FROM weekly_lookups;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE weekly_lookups;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE weekly_lookups_new RENAME TO weekly_lookups;
-- +goose StatementEnd

-- Recreate indexes on weekly_lookups
-- +goose StatementBegin
CREATE INDEX idx_weekly_lookups_name ON weekly_lookups(name);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_weekly_lookups_program_id ON weekly_lookups(program_id);
-- +goose StatementEnd

-- Recreate daily_lookups with foreign key
-- +goose StatementBegin
CREATE TABLE daily_lookups_new (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK(length(name) > 0 AND length(name) <= 100),
    entries TEXT NOT NULL CHECK(json_valid(entries)),
    program_id TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (program_id) REFERENCES programs(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO daily_lookups_new SELECT * FROM daily_lookups;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE daily_lookups;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE daily_lookups_new RENAME TO daily_lookups;
-- +goose StatementEnd

-- Recreate indexes on daily_lookups
-- +goose StatementBegin
CREATE INDEX idx_daily_lookups_name ON daily_lookups(name);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_daily_lookups_program_id ON daily_lookups(program_id);
-- +goose StatementEnd

-- +goose Down
-- Recreate lookup tables without foreign key constraints

-- Recreate weekly_lookups without foreign key
-- +goose StatementBegin
CREATE TABLE weekly_lookups_old (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK(length(name) > 0 AND length(name) <= 100),
    entries TEXT NOT NULL CHECK(json_valid(entries)),
    program_id TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO weekly_lookups_old SELECT * FROM weekly_lookups;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE weekly_lookups;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE weekly_lookups_old RENAME TO weekly_lookups;
-- +goose StatementEnd

-- Recreate indexes on weekly_lookups
-- +goose StatementBegin
CREATE INDEX idx_weekly_lookups_name ON weekly_lookups(name);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_weekly_lookups_program_id ON weekly_lookups(program_id);
-- +goose StatementEnd

-- Recreate daily_lookups without foreign key
-- +goose StatementBegin
CREATE TABLE daily_lookups_old (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK(length(name) > 0 AND length(name) <= 100),
    entries TEXT NOT NULL CHECK(json_valid(entries)),
    program_id TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO daily_lookups_old SELECT * FROM daily_lookups;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE daily_lookups;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE daily_lookups_old RENAME TO daily_lookups;
-- +goose StatementEnd

-- Recreate indexes on daily_lookups
-- +goose StatementBegin
CREATE INDEX idx_daily_lookups_name ON daily_lookups(name);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_daily_lookups_program_id ON daily_lookups(program_id);
-- +goose StatementEnd

-- Drop programs table and indexes
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_programs_cycle_id;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_programs_name;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS programs;
-- +goose StatementEnd

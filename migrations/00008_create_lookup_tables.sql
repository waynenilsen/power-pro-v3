-- +goose Up
-- +goose StatementBegin
CREATE TABLE weekly_lookups (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK(length(name) > 0 AND length(name) <= 100),
    entries TEXT NOT NULL CHECK(json_valid(entries)),
    program_id TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
-- +goose StatementEnd

-- Index for name lookups
-- +goose StatementBegin
CREATE INDEX idx_weekly_lookups_name ON weekly_lookups(name);
-- +goose StatementEnd

-- Index for program_id lookups
-- +goose StatementBegin
CREATE INDEX idx_weekly_lookups_program_id ON weekly_lookups(program_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE daily_lookups (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK(length(name) > 0 AND length(name) <= 100),
    entries TEXT NOT NULL CHECK(json_valid(entries)),
    program_id TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
-- +goose StatementEnd

-- Index for name lookups
-- +goose StatementBegin
CREATE INDEX idx_daily_lookups_name ON daily_lookups(name);
-- +goose StatementEnd

-- Index for program_id lookups
-- +goose StatementBegin
CREATE INDEX idx_daily_lookups_program_id ON daily_lookups(program_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_daily_lookups_program_id;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_daily_lookups_name;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS daily_lookups;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_weekly_lookups_program_id;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_weekly_lookups_name;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS weekly_lookups;
-- +goose StatementEnd

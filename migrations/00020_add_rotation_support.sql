-- +goose Up
-- Create rotation_lookups table for programs that rotate through different lifts
-- +goose StatementBegin
CREATE TABLE rotation_lookups (
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
CREATE INDEX idx_rotation_lookups_name ON rotation_lookups(name);
-- +goose StatementEnd

-- Index for program_id lookups
-- +goose StatementBegin
CREATE INDEX idx_rotation_lookups_program_id ON rotation_lookups(program_id);
-- +goose StatementEnd

-- Add rotation tracking fields to user_program_states
-- +goose StatementBegin
ALTER TABLE user_program_states ADD COLUMN rotation_position INTEGER NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE user_program_states ADD COLUMN cycles_since_start INTEGER NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_program_states DROP COLUMN cycles_since_start;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE user_program_states DROP COLUMN rotation_position;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_rotation_lookups_program_id;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_rotation_lookups_name;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS rotation_lookups;
-- +goose StatementEnd

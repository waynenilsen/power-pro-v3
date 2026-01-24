-- +goose Up
-- +goose StatementBegin
-- SQLite doesn't support ALTER TABLE to modify CHECK constraints directly.
-- We need to recreate the table with the updated constraint.

-- Step 1: Create new table with updated constraint
CREATE TABLE progressions_new (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK(length(name) > 0 AND length(name) <= 100),
    type TEXT NOT NULL CHECK(type IN (
        'LINEAR_PROGRESSION',
        'CYCLE_PROGRESSION',
        'AMRAP_PROGRESSION',
        'DELOAD_ON_FAILURE',
        'STAGE_PROGRESSION'
    )),
    parameters TEXT NOT NULL CHECK(json_valid(parameters)),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
-- Step 2: Copy data from old table
INSERT INTO progressions_new SELECT * FROM progressions;
-- +goose StatementEnd

-- +goose StatementBegin
-- Step 3: Drop old table
DROP TABLE progressions;
-- +goose StatementEnd

-- +goose StatementBegin
-- Step 4: Rename new table
ALTER TABLE progressions_new RENAME TO progressions;
-- +goose StatementEnd

-- +goose StatementBegin
-- Step 5: Recreate index
CREATE INDEX idx_progressions_type ON progressions(type);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert to original constraint (removes DELOAD_ON_FAILURE and STAGE_PROGRESSION)
CREATE TABLE progressions_old (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK(length(name) > 0 AND length(name) <= 100),
    type TEXT NOT NULL CHECK(type IN ('LINEAR_PROGRESSION', 'CYCLE_PROGRESSION', 'AMRAP_PROGRESSION')),
    parameters TEXT NOT NULL CHECK(json_valid(parameters)),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
-- This will fail if there are any DELOAD_ON_FAILURE or STAGE_PROGRESSION rows
INSERT INTO progressions_old SELECT * FROM progressions;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE progressions;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE progressions_old RENAME TO progressions;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_progressions_type ON progressions(type);
-- +goose StatementEnd

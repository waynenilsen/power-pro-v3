-- +goose Up
-- +goose StatementBegin
-- SQLite doesn't support ALTER TABLE to modify CHECK constraints directly.
-- We need to recreate the table with the updated constraint that includes
-- ON_FAILURE and AFTER_SET trigger types.

-- Step 1: Create new table with updated constraint
CREATE TABLE progression_logs_new (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    progression_id TEXT NOT NULL REFERENCES progressions(id) ON DELETE CASCADE,
    lift_id TEXT NOT NULL REFERENCES lifts(id),
    previous_value REAL NOT NULL CHECK(previous_value >= 0),
    new_value REAL NOT NULL CHECK(new_value >= 0),
    delta REAL NOT NULL,
    trigger_type TEXT NOT NULL CHECK(trigger_type IN (
        'AFTER_SESSION',
        'AFTER_WEEK',
        'AFTER_CYCLE',
        'AFTER_SET',
        'ON_FAILURE'
    )),
    trigger_context TEXT,
    applied_at TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
-- Step 2: Copy data from old table
INSERT INTO progression_logs_new SELECT * FROM progression_logs;
-- +goose StatementEnd

-- +goose StatementBegin
-- Step 3: Drop old table
DROP TABLE progression_logs;
-- +goose StatementEnd

-- +goose StatementBegin
-- Step 4: Rename new table
ALTER TABLE progression_logs_new RENAME TO progression_logs;
-- +goose StatementEnd

-- +goose StatementBegin
-- Step 5: Recreate indexes
CREATE UNIQUE INDEX idx_progression_logs_idempotency
    ON progression_logs(user_id, progression_id, lift_id, trigger_type, applied_at);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_progression_logs_user_progression
    ON progression_logs(user_id, progression_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert to original constraint (removes ON_FAILURE and AFTER_SET)
CREATE TABLE progression_logs_old (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    progression_id TEXT NOT NULL REFERENCES progressions(id) ON DELETE CASCADE,
    lift_id TEXT NOT NULL REFERENCES lifts(id),
    previous_value REAL NOT NULL CHECK(previous_value >= 0),
    new_value REAL NOT NULL CHECK(new_value >= 0),
    delta REAL NOT NULL,
    trigger_type TEXT NOT NULL CHECK(trigger_type IN ('AFTER_SESSION', 'AFTER_WEEK', 'AFTER_CYCLE')),
    trigger_context TEXT,
    applied_at TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
-- This will fail if there are any ON_FAILURE or AFTER_SET rows
INSERT INTO progression_logs_old SELECT * FROM progression_logs;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE progression_logs;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE progression_logs_old RENAME TO progression_logs;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE UNIQUE INDEX idx_progression_logs_idempotency
    ON progression_logs(user_id, progression_id, lift_id, trigger_type, applied_at);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_progression_logs_user_progression
    ON progression_logs(user_id, progression_id);
-- +goose StatementEnd

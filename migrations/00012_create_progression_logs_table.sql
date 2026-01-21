-- +goose Up
-- +goose StatementBegin
CREATE TABLE progression_logs (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    progression_id TEXT NOT NULL,
    lift_id TEXT NOT NULL,
    previous_value REAL NOT NULL,
    new_value REAL NOT NULL,
    delta REAL NOT NULL,
    trigger_type TEXT NOT NULL CHECK(trigger_type IN ('AFTER_SESSION', 'AFTER_WEEK', 'AFTER_CYCLE')),
    trigger_context TEXT NOT NULL CHECK(json_valid(trigger_context)),
    applied_at TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (progression_id) REFERENCES progressions(id) ON DELETE CASCADE,
    FOREIGN KEY (lift_id) REFERENCES lifts(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- Unique constraint for idempotency: prevent double-application of progressions
-- +goose StatementBegin
CREATE UNIQUE INDEX idx_progression_logs_idempotency
    ON progression_logs(user_id, progression_id, trigger_type, applied_at);
-- +goose StatementEnd

-- Index for history queries (looking up progression history for a user's lift)
-- +goose StatementBegin
CREATE INDEX idx_progression_logs_user_lift
    ON progression_logs(user_id, lift_id);
-- +goose StatementEnd

-- Index for date range queries
-- +goose StatementBegin
CREATE INDEX idx_progression_logs_applied_at
    ON progression_logs(applied_at);
-- +goose StatementEnd

-- Index for foreign key lookups
-- +goose StatementBegin
CREATE INDEX idx_progression_logs_progression_id
    ON progression_logs(progression_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_progression_logs_progression_id;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_progression_logs_applied_at;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_progression_logs_user_lift;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_progression_logs_idempotency;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS progression_logs;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE failure_counters (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    lift_id TEXT NOT NULL,
    progression_id TEXT NOT NULL,
    consecutive_failures INT NOT NULL DEFAULT 0,
    last_failure_at TEXT,
    last_success_at TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (lift_id) REFERENCES lifts(id) ON DELETE CASCADE,
    FOREIGN KEY (progression_id) REFERENCES progressions(id) ON DELETE CASCADE,
    UNIQUE(user_id, lift_id, progression_id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_failure_counters_user ON failure_counters(user_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_failure_counters_lift ON failure_counters(lift_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_failure_counters_progression ON failure_counters(progression_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_failure_counters_lookup ON failure_counters(user_id, lift_id, progression_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_failure_counters_lookup;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_failure_counters_progression;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_failure_counters_lift;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_failure_counters_user;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS failure_counters;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE lift_maxes (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    lift_id TEXT NOT NULL,
    type TEXT NOT NULL CHECK(type IN ('ONE_RM', 'TRAINING_MAX')),
    value REAL NOT NULL CHECK(value > 0),
    effective_date TEXT NOT NULL DEFAULT (datetime('now')),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (lift_id) REFERENCES lifts(id) ON DELETE RESTRICT,
    UNIQUE (user_id, lift_id, type, effective_date)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_lift_maxes_user_lift_type ON lift_maxes(user_id, lift_id, type);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_lift_maxes_effective_date ON lift_maxes(effective_date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_lift_maxes_effective_date;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_lift_maxes_user_lift_type;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS lift_maxes;
-- +goose StatementEnd

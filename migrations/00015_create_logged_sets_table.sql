-- +goose Up
-- +goose StatementBegin
CREATE TABLE logged_sets (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    session_id TEXT NOT NULL,
    prescription_id TEXT NOT NULL,
    lift_id TEXT NOT NULL,
    set_number INT NOT NULL,
    weight REAL NOT NULL,
    target_reps INT NOT NULL,
    reps_performed INT NOT NULL,
    is_amrap BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (lift_id) REFERENCES lifts(id) ON DELETE CASCADE,
    UNIQUE(session_id, prescription_id, set_number)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_logged_sets_user ON logged_sets(user_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_logged_sets_session ON logged_sets(session_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_logged_sets_lift ON logged_sets(lift_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_logged_sets_lift_amrap ON logged_sets(user_id, lift_id, is_amrap, created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_logged_sets_lift_amrap;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_logged_sets_lift;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_logged_sets_session;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_logged_sets_user;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS logged_sets;
-- +goose StatementEnd

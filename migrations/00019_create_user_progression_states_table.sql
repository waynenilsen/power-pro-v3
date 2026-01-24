-- +goose Up
-- +goose StatementBegin
-- User progression states track per-user, per-lift, per-progression state.
-- This is used by StageProgression to track which stage a user is at.
CREATE TABLE user_progression_states (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    lift_id TEXT NOT NULL,
    progression_id TEXT NOT NULL,
    current_stage INT NOT NULL DEFAULT 0,
    state_data TEXT, -- JSON blob for additional progression-specific state
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (lift_id) REFERENCES lifts(id) ON DELETE CASCADE,
    FOREIGN KEY (progression_id) REFERENCES progressions(id) ON DELETE CASCADE,
    UNIQUE(user_id, lift_id, progression_id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_user_progression_states_user ON user_progression_states(user_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_user_progression_states_lift ON user_progression_states(lift_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_user_progression_states_progression ON user_progression_states(progression_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_user_progression_states_lookup ON user_progression_states(user_id, lift_id, progression_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_progression_states_lookup;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_progression_states_progression;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_progression_states_lift;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_progression_states_user;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS user_progression_states;
-- +goose StatementEnd

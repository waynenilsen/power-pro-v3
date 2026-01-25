-- +goose Up
-- +goose StatementBegin
CREATE TABLE workout_sessions (
    id TEXT PRIMARY KEY,
    user_program_state_id TEXT NOT NULL,
    week_number INTEGER NOT NULL,
    day_index INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'IN_PROGRESS'
        CHECK (status IN ('IN_PROGRESS', 'COMPLETED', 'ABANDONED')),
    started_at TEXT NOT NULL DEFAULT (datetime('now')),
    finished_at TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (user_program_state_id) REFERENCES user_program_states(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_workout_sessions_state_week_day
    ON workout_sessions(user_program_state_id, week_number, day_index);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_workout_sessions_status
    ON workout_sessions(user_program_state_id, status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_workout_sessions_status;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_workout_sessions_state_week_day;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS workout_sessions;
-- +goose StatementEnd

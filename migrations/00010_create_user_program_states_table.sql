-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_program_states (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL UNIQUE,
    program_id TEXT NOT NULL,
    current_week INTEGER NOT NULL CHECK(current_week >= 1),
    current_cycle_iteration INTEGER NOT NULL CHECK(current_cycle_iteration >= 1),
    current_day_index INTEGER,
    enrolled_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (program_id) REFERENCES programs(id) ON DELETE RESTRICT
);
-- +goose StatementEnd

-- Index for program_id lookups (user_id already indexed via UNIQUE constraint)
-- +goose StatementBegin
CREATE INDEX idx_user_program_states_program_id ON user_program_states(program_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_program_states_program_id;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS user_program_states;
-- +goose StatementEnd

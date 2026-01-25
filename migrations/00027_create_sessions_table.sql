-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    token TEXT NOT NULL,
    expires_at TEXT NOT NULL,
    created_at TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- Unique index on token for fast authentication lookup
-- +goose StatementBegin
CREATE UNIQUE INDEX idx_sessions_token ON sessions(token);
-- +goose StatementEnd

-- Index on user_id for looking up user's sessions
-- +goose StatementBegin
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
-- +goose StatementEnd

-- Index on expires_at for cleanup queries (deleting expired sessions)
-- +goose StatementBegin
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_sessions_expires_at;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_sessions_user_id;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_sessions_token;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS sessions;
-- +goose StatementEnd

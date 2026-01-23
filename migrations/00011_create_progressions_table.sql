-- +goose Up
-- +goose StatementBegin
CREATE TABLE progressions (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK(length(name) > 0 AND length(name) <= 100),
    type TEXT NOT NULL CHECK(type IN ('LINEAR_PROGRESSION', 'CYCLE_PROGRESSION', 'AMRAP_PROGRESSION')),
    parameters TEXT NOT NULL CHECK(json_valid(parameters)),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
-- +goose StatementEnd

-- Index for type lookups (filtering by progression type)
-- +goose StatementBegin
CREATE INDEX idx_progressions_type ON progressions(type);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_progressions_type;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS progressions;
-- +goose StatementEnd

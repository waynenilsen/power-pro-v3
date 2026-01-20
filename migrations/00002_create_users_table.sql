-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
-- +goose StatementEnd

-- Seed test users with deterministic IDs
-- These are used for E2E testing
-- +goose StatementBegin
INSERT INTO users (id, created_at, updated_at) VALUES
    ('test-user-001', datetime('now'), datetime('now')),
    ('test-admin-001', datetime('now'), datetime('now')),
    ('create-test-user', datetime('now'), datetime('now')),
    ('date-test-user', datetime('now'), datetime('now')),
    ('missing-lift-user', datetime('now'), datetime('now')),
    ('missing-type-user', datetime('now'), datetime('now')),
    ('invalid-type-user', datetime('now'), datetime('now')),
    ('invalid-precision-user', datetime('now'), datetime('now')),
    ('zero-value-user', datetime('now'), datetime('now')),
    ('bad-lift-user', datetime('now'), datetime('now')),
    ('duplicate-test-user', datetime('now'), datetime('now')),
    ('tm-warning-user', datetime('now'), datetime('now')),
    ('update-user', datetime('now'), datetime('now')),
    ('delete-user', datetime('now'), datetime('now')),
    ('format-user', datetime('now'), datetime('now')),
    ('sort-test-user', datetime('now'), datetime('now')),
    ('convert-test-user', datetime('now'), datetime('now')),
    ('convert-tm-user', datetime('now'), datetime('now')),
    ('round-test-user', datetime('now'), datetime('now')),
    ('admin-user', datetime('now'), datetime('now')),
    ('current-max-user', datetime('now'), datetime('now')),
    ('single-max-user', datetime('now'), datetime('now')),
    -- Authorization test users
    ('user-a', datetime('now'), datetime('now')),
    ('user-b', datetime('now'), datetime('now')),
    ('current-user-a', datetime('now'), datetime('now')),
    ('current-user-b', datetime('now'), datetime('now')),
    ('convert-user-a', datetime('now'), datetime('now')),
    ('convert-user-b', datetime('now'), datetime('now'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

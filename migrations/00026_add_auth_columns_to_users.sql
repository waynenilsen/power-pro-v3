-- +goose Up
-- Add authentication columns to users table
-- Existing test users will have NULL email/password_hash and continue to work via X-User-ID header

-- +goose StatementBegin
ALTER TABLE users ADD COLUMN email TEXT CHECK(email IS NULL OR length(email) <= 255);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE users ADD COLUMN password_hash TEXT;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE users ADD COLUMN name TEXT CHECK(name IS NULL OR length(name) <= 100);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE users ADD COLUMN is_admin INTEGER NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- Partial unique index on email - only enforces uniqueness for non-null emails
-- +goose StatementBegin
CREATE UNIQUE INDEX idx_users_email ON users(email) WHERE email IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_email;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE users DROP COLUMN is_admin;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE users DROP COLUMN name;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE users DROP COLUMN password_hash;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE users DROP COLUMN email;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_program_states
    ADD COLUMN enrollment_status TEXT NOT NULL DEFAULT 'ACTIVE'
    CHECK (enrollment_status IN ('ACTIVE', 'BETWEEN_CYCLES', 'QUIT'));
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE user_program_states
    ADD COLUMN cycle_status TEXT NOT NULL DEFAULT 'PENDING'
    CHECK (cycle_status IN ('PENDING', 'IN_PROGRESS', 'COMPLETED'));
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE user_program_states
    ADD COLUMN week_status TEXT NOT NULL DEFAULT 'PENDING'
    CHECK (week_status IN ('PENDING', 'IN_PROGRESS', 'COMPLETED'));
-- +goose StatementEnd

-- +goose Down
-- SQLite doesn't support DROP COLUMN in older versions, so we need to recreate the table
-- For development, we can just leave this as a placeholder since we're moving forward
-- +goose StatementBegin
SELECT 1; -- Placeholder for down migration
-- +goose StatementEnd

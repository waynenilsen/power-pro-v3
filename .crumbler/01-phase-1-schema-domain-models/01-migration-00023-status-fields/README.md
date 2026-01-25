# Migration 00023 - Status Fields

## Task

Create migration `migrations/00023_add_state_machine_status_fields.sql` to add status fields to the user_program_states table.

## SQL to Create

```sql
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
```

## Verification

- Run `make migrate` or equivalent to apply migration
- Verify the columns exist using `sqlite3`

## Done When

- Migration file exists at correct path
- Migration applies successfully (or tests still pass if migration runs automatically)

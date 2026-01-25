# Sessions Schema Migration

## Ticket Reference
`phases/in-progress/002-frontend-readiness/sprints/in-progress/001-authentication-system/tickets/todo/002-sessions-schema-migration.md`

## Task
Create a goose migration to add the sessions table for server-side session storage.

## Implementation

1. Create migration file in `internal/db/migrations/`
2. Create sessions table with:
   - `id` (TEXT, primary key, UUID)
   - `user_id` (TEXT, required, FK to users.id with ON DELETE CASCADE)
   - `token` (TEXT, required, unique)
   - `expires_at` (TEXT, required, ISO8601)
   - `created_at` (TEXT, required, ISO8601)
3. Create indexes:
   - Unique index on token
   - Index on user_id
   - Index on expires_at
4. Include proper down migration that drops table

## Acceptance Criteria
- [ ] Sessions table created with all columns
- [ ] Foreign key to users with CASCADE delete
- [ ] Unique index on token
- [ ] Index on user_id and expires_at
- [ ] Down migration drops table

## When Done
Move ticket from `tickets/todo/` to `tickets/done/` then run `crumbler delete`

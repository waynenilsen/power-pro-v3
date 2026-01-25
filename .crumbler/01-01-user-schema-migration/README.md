# User Schema Migration

## Ticket Reference
`phases/in-progress/002-frontend-readiness/sprints/in-progress/001-authentication-system/tickets/todo/001-user-schema-migration.md`

## Task
Create a goose migration to add authentication columns to the users table.

## Implementation

1. Create migration file in `internal/db/migrations/` (follow existing naming convention)
2. Add columns to users table:
   - `email` (TEXT, nullable, unique when not null)
   - `password_hash` (TEXT, nullable)
   - `name` (TEXT, nullable, max 100 chars)
   - `is_admin` (INTEGER, default 0)
3. Create partial unique index on email: `CREATE UNIQUE INDEX idx_users_email ON users(email) WHERE email IS NOT NULL`
4. Add CHECK constraints for max lengths (email: 255, name: 100)
5. Include proper down migration that removes columns
6. Run and verify migration works

## Acceptance Criteria
- [ ] Migration adds email, password_hash, name, is_admin columns
- [ ] Partial unique index on email where not null
- [ ] CHECK constraints for email (255) and name (100) max lengths
- [ ] Down migration removes columns correctly
- [ ] Existing test users preserved with NULL email/password_hash

## When Done
Move ticket from `tickets/todo/` to `tickets/done/` then run `crumbler delete`

# 002: Backfill Migration for Canonical Programs

## Task
Create a goose migration that backfills discovery metadata for the four canonical programs.

## Program Metadata Values

| Program | Slug | Difficulty | Days/Week | Focus | Has AMRAP |
|---------|------|------------|-----------|-------|-----------|
| Starting Strength | starting-strength | beginner | 3 | strength | 0 |
| Texas Method | texas-method | intermediate | 3 | strength | 0 |
| Wendler 5/3/1 | 531 | intermediate | 4 | strength | 1 |
| GZCLP | gzclp | beginner | 4 | strength | 1 |

## Acceptance Criteria
- Create migration file in `migrations/` directory
- Use UPDATE statements matching on slug
- Migration must be idempotent (safe to run multiple times)
- Migration should not fail if programs don't exist
- Include down migration (reset to defaults or no-op)
- Verify metadata values after running

## Reference
- Ticket: `phases/in-progress/002-frontend-readiness/sprints/todo/006-program-discovery/tickets/todo/002-backfill-migration.md`
- Depends on: 001-schema-migration (columns must exist)

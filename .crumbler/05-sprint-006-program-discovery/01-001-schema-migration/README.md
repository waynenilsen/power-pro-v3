# 001: Program Metadata Schema Migration

## Task
Create a goose migration that adds program discovery metadata columns to the programs table.

## Columns to Add

1. **difficulty** - TEXT NOT NULL DEFAULT 'beginner'
   - CHECK: value IN ('beginner', 'intermediate', 'advanced')
   - Index: `idx_programs_difficulty`

2. **days_per_week** - INTEGER NOT NULL DEFAULT 3
   - CHECK: value BETWEEN 1 AND 7
   - Index: `idx_programs_days_per_week`

3. **focus** - TEXT NOT NULL DEFAULT 'strength'
   - CHECK: value IN ('strength', 'hypertrophy', 'peaking')
   - Index: `idx_programs_focus`

4. **has_amrap** - INTEGER NOT NULL DEFAULT 0
   - CHECK: value IN (0, 1)
   - Index: `idx_programs_has_amrap`

## Acceptance Criteria
- Create migration file in `migrations/` directory
- SQLite requires separate ALTER TABLE for each column
- Include proper down migration (DROP indices, DROP columns)
- Run migration and verify columns exist
- Verify existing programs get default values

## Reference
- Ticket: `phases/in-progress/002-frontend-readiness/sprints/todo/006-program-discovery/tickets/todo/001-schema-migration.md`
- Existing schema: `migrations/` directory

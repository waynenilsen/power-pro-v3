# 001: Program Metadata Schema Migration

## ERD Reference
Implements: REQ-SCHEMA-001, REQ-SCHEMA-002, REQ-SCHEMA-003, REQ-SCHEMA-004

## Description
Create a goose migration that adds program discovery metadata columns to the programs table. This includes difficulty, days_per_week, focus, and has_amrap columns with appropriate constraints and indices.

## Context / Background
Program discovery requires metadata fields that enable filtering. These columns must be added before any filtering logic can be implemented. The migration must be backward-compatible with existing programs by providing sensible defaults.

## Acceptance Criteria
- [ ] Create goose migration file: `NNNN_add_program_discovery_columns.sql`
- [ ] Add `difficulty` column:
  - Type: TEXT NOT NULL DEFAULT 'beginner'
  - CHECK constraint: value IN ('beginner', 'intermediate', 'advanced')
  - Create index: `idx_programs_difficulty`
- [ ] Add `days_per_week` column:
  - Type: INTEGER NOT NULL DEFAULT 3
  - CHECK constraint: value BETWEEN 1 AND 7
  - Create index: `idx_programs_days_per_week`
- [ ] Add `focus` column:
  - Type: TEXT NOT NULL DEFAULT 'strength'
  - CHECK constraint: value IN ('strength', 'hypertrophy', 'peaking')
  - Create index: `idx_programs_focus`
- [ ] Add `has_amrap` column:
  - Type: INTEGER NOT NULL DEFAULT 0
  - CHECK constraint: value IN (0, 1)
  - Create index: `idx_programs_has_amrap`
- [ ] Migration includes proper down migration (DROP columns, DROP indices)
- [ ] Run migration locally and verify columns exist with correct constraints
- [ ] Verify existing programs have default values applied

## Technical Notes
- Use ALTER TABLE ADD COLUMN for each new column
- SQLite requires separate ALTER TABLE statements for each column
- Create indices after columns are added
- Down migration must drop indices before dropping columns
- Default values ensure existing programs remain valid after migration
- Consider adding composite index for common filter combinations in future optimization

## Dependencies
- Blocks: 002 (backfill needs columns), 003 (filtering needs columns), 004 (search benefits from schema), 005 (detail needs columns), 006 (tests need columns)
- Blocked by: None (first ticket in sequence)

## Resources / Links
- Existing schema: `migrations/00009_create_programs_table.sql`
- ERD requirements: `phases/in-progress/002-frontend-readiness/sprints/todo/006-program-discovery/erd.md`

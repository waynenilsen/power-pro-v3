# 005: Program Entity Schema and Migration

## ERD Reference
Implements: REQ-PROG-001

## Description
Create the database schema for the Program entity, including the goose migration file. The Program entity bundles cycles, lookups, and defaults into a named configuration that users can enroll in.

## Context / Background
A Program is a named configuration that bundles cycles, lookups, and defaults. Programs reference one cycle and can reference lookup tables. Examples: "5/3/1 BBB", "Starting Strength", "Bill Starr 5x5".

## Acceptance Criteria
- [ ] Program table created with: id (UUID), name (VARCHAR, NOT NULL), slug (VARCHAR, UNIQUE, NOT NULL), description (TEXT, nullable), cycle_id (UUID FK), default_rounding (DECIMAL, nullable)
- [ ] Optional foreign keys to weekly_lookup_id and daily_lookup_id
- [ ] Slug uniqueness constraint
- [ ] Goose migration file created
- [ ] Migration tested (up and down)

## Technical Notes
- Program table: `programs` with columns: id, name, slug, description, cycle_id, weekly_lookup_id (nullable FK), daily_lookup_id (nullable FK), default_rounding, created_at, updated_at
- slug: lowercase alphanumeric with hyphens, unique
- default_rounding: weight rounding increment (e.g., 2.5 for 2.5lb plates)
- One program references one cycle
- Program can optionally reference WeeklyLookup and DailyLookup tables

## Dependencies
- Blocks: 006, 012, 013
- Blocked by: 003 (Cycle schema), 004 (Lookup tables schema)
- Related: All Day/Week/Cycle tickets

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/003-schedule-periodization/erd.md

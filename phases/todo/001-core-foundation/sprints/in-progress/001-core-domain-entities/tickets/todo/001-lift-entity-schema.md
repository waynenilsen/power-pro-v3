# 001: Lift Entity Schema and Migration

## ERD Reference
Implements: REQ-LIFT-001, REQ-LIFT-002, REQ-LIFT-003, REQ-LIFT-004, REQ-LIFT-005

## Description
Create the database schema and migration for the Lift entity. This establishes the foundational table structure for storing lift definitions (exercises) in the system.

## Context / Background
The Lift entity is the atomic building block of PowerPro. Every prescription, progression, and schedule references a lift. This ticket creates the database foundation that all Lift-related features depend on.

## Acceptance Criteria
- [ ] Create `lifts` table with the following columns:
  - `id` (UUID, primary key)
  - `name` (VARCHAR(100), required, non-empty)
  - `slug` (VARCHAR(100), unique, lowercase alphanumeric with hyphens only)
  - `is_competition_lift` (BOOLEAN, default false)
  - `parent_lift_id` (UUID, nullable, foreign key to lifts.id)
  - `created_at` (TIMESTAMP, required)
  - `updated_at` (TIMESTAMP, required)
- [ ] Create unique constraint on `slug` column
- [ ] Create index on `parent_lift_id` for efficient parent-child queries
- [ ] Create goose migration file with proper up/down migrations
- [ ] Migration handles circular reference prevention at database level (CHECK constraint or trigger)
- [ ] Seed competition lifts (Squat, Bench Press, Deadlift) with `is_competition_lift = true`

## Technical Notes
- Use PostgreSQL with TypeORM entities as specified in ERD
- Slug should be enforced as lowercase alphanumeric with hyphens via CHECK constraint
- Consider adding a CHECK constraint to prevent `parent_lift_id = id` (self-reference)
- Seed data should use deterministic UUIDs for competition lifts to enable cross-environment consistency

## Dependencies
- Blocks: 002, 003 (Lift domain logic and CRUD API depend on this schema)
- Blocked by: None
- Related: 004 (LiftMax schema references this table)

## Resources / Links
- ERD: phases/todo/001-core-foundation/sprints/in-progress/001-core-domain-entities/erd.md
- Tech Stack: prompts/tech-stack.md

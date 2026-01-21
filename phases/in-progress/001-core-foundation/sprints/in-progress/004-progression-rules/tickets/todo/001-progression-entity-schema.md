# 001: Progression Entity Schema and Migration

## ERD Reference
Implements: REQ-PROG-001
Related to: NFR-006, NFR-007

## Description
Create the database schema and migration for the Progression entity. This establishes the foundational table structure for storing progression rule definitions with type discrimination and JSONB parameters.

## Context / Background
Progressions are rules that mutate LiftMax values over time. Different programs use different progression strategies (linear, cycle-based, etc.). This table stores all progression types using a discriminated union pattern where the `type` field determines how to interpret the `parameters` JSON.

## Acceptance Criteria
- [ ] Create `progressions` table with the following columns:
  - `id` (UUID, primary key)
  - `name` (VARCHAR(100), required, non-empty)
  - `type` (VARCHAR(50), required) - discriminator: LINEAR_PROGRESSION, CYCLE_PROGRESSION
  - `parameters` (JSONB, required) - type-specific configuration
  - `created_at` (TIMESTAMP, required)
  - `updated_at` (TIMESTAMP, required)
- [ ] Create index on `type` column for filtering by progression type
- [ ] Create goose migration file with proper up/down migrations
- [ ] Add CHECK constraint ensuring `type` is a valid progression type
- [ ] Validate that migration handles empty parameters gracefully (empty JSON object `{}`)

## Technical Notes
- Use SQLite JSONB for parameters storage (or TEXT with JSON validation)
- Type discriminator enables polymorphic handling in Go code
- Parameters structure varies by type:
  - LINEAR_PROGRESSION: `{"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"}`
  - CYCLE_PROGRESSION: `{"increment": 5.0, "maxType": "TRAINING_MAX"}`
- Consider CHECK constraint on type enum values

## Dependencies
- Blocks: 003, 004, 005, 006, 009 (ProgramProgression, interface, implementations, and CRUD depend on this)
- Blocked by: None
- Related: 002 (ProgressionLog references this table)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/004-progression-rules/erd.md
- Tech Stack: prompts/tech-stack.md

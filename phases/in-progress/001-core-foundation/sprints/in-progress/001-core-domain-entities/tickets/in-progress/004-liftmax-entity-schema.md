# 004: LiftMax Entity Schema and Migration

## ERD Reference
Implements: REQ-MAX-001, REQ-MAX-002, REQ-MAX-003, REQ-MAX-004, REQ-MAX-005, REQ-MAX-006, REQ-MAX-007

## Description
Create the database schema and migration for the LiftMax entity. This establishes the table structure for storing user-specific reference values (1RM, Training Max) used in load calculations.

## Context / Background
LiftMax stores the reference numbers used for load calculation. Different programs use different reference types (1RM, Training Max). A user has multiple LiftMax records - one for each (lift, type, effective_date) combination they track.

## Acceptance Criteria
- [ ] Create `lift_maxes` table with the following columns:
  - `id` (UUID, primary key)
  - `user_id` (UUID, required, foreign key to users.id)
  - `lift_id` (UUID, required, foreign key to lifts.id)
  - `type` (ENUM: 'ONE_RM', 'TRAINING_MAX', required)
  - `value` (DECIMAL, required, positive, precision to 0.25)
  - `effective_date` (TIMESTAMP, required, defaults to now)
  - `created_at` (TIMESTAMP, required)
  - `updated_at` (TIMESTAMP, required)
- [ ] Create unique constraint on (`user_id`, `lift_id`, `type`, `effective_date`)
- [ ] Create foreign key constraint to `lifts` table with ON DELETE RESTRICT
- [ ] Create foreign key constraint to `users` table with ON DELETE CASCADE
- [ ] Create index on (`user_id`, `lift_id`, `type`) for efficient current max lookups
- [ ] Create index on `effective_date` for temporal queries
- [ ] Create CHECK constraint ensuring `value > 0`
- [ ] Create goose migration file with proper up/down migrations
- [ ] Create max_type ENUM type in database

## Technical Notes
- Use PostgreSQL DECIMAL type for value to maintain precision (not FLOAT)
- The value is unit-agnostic per ERD (stored as raw number, unit handling at API layer)
- Consider using DECIMAL(8,2) to support up to 999,999.99 with 2 decimal precision
- The 0.25 precision requirement means values like 315.25, 315.5, 315.75 are valid
- User entity is assumed to exist per ERD Section 6 Assumptions

## Dependencies
- Blocks: 005, 006, 007, 008 (All LiftMax tickets depend on this schema)
- Blocked by: 003 (Lift CRUD API must exist for referential integrity testing)
- Related: 001 (References lifts table)

## Resources / Links
- ERD: phases/todo/001-core-foundation/sprints/in-progress/001-core-domain-entities/erd.md
- Tech Stack: prompts/tech-stack.md

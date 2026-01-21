# 002: ProgressionLog Entity Schema and Migration

## ERD Reference
Implements: REQ-HIST-001, REQ-TRIG-005
Related to: NFR-004

## Description
Create the database schema and migration for the ProgressionLog entity. This table tracks all progression applications for audit trail, debugging, and idempotency enforcement.

## Context / Background
Every progression application must be logged. This log serves two purposes: (1) providing an audit trail of all LiftMax changes, and (2) enabling idempotency checks to prevent double-application on retry/error scenarios. The unique constraint on trigger context ensures progressions are never applied multiple times for the same event.

## Acceptance Criteria
- [ ] Create `progression_logs` table with the following columns:
  - `id` (UUID, primary key)
  - `user_id` (UUID, required, foreign key to users.id)
  - `progression_id` (UUID, required, foreign key to progressions.id)
  - `lift_id` (UUID, required, foreign key to lifts.id)
  - `previous_value` (DECIMAL, required) - LiftMax value before progression
  - `new_value` (DECIMAL, required) - LiftMax value after progression
  - `delta` (DECIMAL, required) - the increment applied
  - `trigger_type` (VARCHAR(50), required) - AFTER_SESSION, AFTER_WEEK, AFTER_CYCLE
  - `trigger_context` (JSONB, required) - context about when trigger fired
  - `applied_at` (TIMESTAMP, required) - when progression was applied
- [ ] Create unique constraint on (`user_id`, `progression_id`, `trigger_type`, `applied_at`) for idempotency
- [ ] Create index on (`user_id`, `lift_id`) for history queries
- [ ] Create index on `applied_at` for date range queries
- [ ] Create goose migration file with proper up/down migrations
- [ ] Foreign key constraints with appropriate ON DELETE behavior

## Technical Notes
- trigger_context JSONB structure varies by trigger type:
  - AFTER_SESSION: `{"sessionId": "uuid", "daySlug": "day-a", "weekNumber": 2}`
  - AFTER_WEEK: `{"weekNumber": 2, "cycleIteration": 1}`
  - AFTER_CYCLE: `{"cycleIteration": 1, "totalWeeks": 4}`
- The unique constraint enables idempotency: before applying, query for existing log entry
- Consider composite index for the idempotency lookup

## Dependencies
- Blocks: 007, 008, 011 (Trigger integration and history query depend on this)
- Blocked by: 001 (References progressions table)
- Related: ERD-001 (References lifts and users tables)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/004-progression-rules/erd.md
- Tech Stack: prompts/tech-stack.md

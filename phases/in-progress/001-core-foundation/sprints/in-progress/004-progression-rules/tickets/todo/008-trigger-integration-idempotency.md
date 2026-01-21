# 008: Trigger Integration and Idempotency Enforcement

## ERD Reference
Implements: REQ-TRIG-002, REQ-TRIG-003, REQ-TRIG-004, REQ-TRIG-005
Related to: NFR-003, NFR-004, NFR-005

## Description
Implement the trigger integration layer that receives events from the Schedule system and applies progressions atomically with idempotency guarantees. This is the orchestration layer connecting state advancement to progression application.

## Context / Background
When users complete sessions or advance through their program, the Schedule system fires trigger events. This ticket implements the handler that receives those events, looks up applicable progressions, applies them, and logs the results. Critical requirements include atomic transactions (no partial updates) and idempotency (prevent double-application on retry).

## Acceptance Criteria
- [ ] Implement `ProgressionService` with trigger handlers:
  - `HandleSessionComplete(ctx, SessionTriggerContext) error`
  - `HandleWeekAdvance(ctx, WeekTriggerContext) error`
  - `HandleCycleComplete(ctx, CycleTriggerContext) error`
- [ ] Trigger handlers:
  - Look up user's enrolled program
  - Fetch applicable ProgramProgressions (enabled, matching trigger type)
  - For each progression, call Apply() with appropriate context
  - Log all applications to ProgressionLog
  - Return aggregate result
- [ ] AFTER_SESSION trigger implementation:
  - Fired after completing a training day
  - Passes session context with lifts performed
  - Only applies progressions to lifts in the completed session
- [ ] AFTER_WEEK trigger implementation:
  - Fired when user advances from week N to week N+1
  - Passes week context
  - Applies to all lifts configured in program progressions
- [ ] AFTER_CYCLE trigger implementation:
  - Fired when cycle completes (week wraps to 1)
  - Passes cycle context
  - Applies to all lifts configured in program progressions
- [ ] Idempotency enforcement:
  - Before applying, check ProgressionLog for existing entry
  - Use unique constraint (user_id, progression_id, trigger_type, applied_at)
  - Skip application if already exists (return applied=false)
  - 100% reliability in preventing double-application (NFR-004)
- [ ] Atomic transactions:
  - Single transaction for: LiftMax update + ProgressionLog creation
  - Rollback on any failure (NFR-003)
  - No partial updates
- [ ] Batch progression performance:
  - Multiple lifts in single transaction where possible
  - Complete in < 500ms for batch (NFR-002)
- [ ] Integration tests covering:
  - Full flow from trigger to LiftMax update
  - Idempotency on retry
  - Transaction rollback on failure
  - Priority ordering of progressions

## Technical Notes
- Hook into State Advancement API from ERD-003
- Consider event-driven approach (internal event bus) or direct call
- Transaction isolation level should prevent concurrent progression
- Use context.Context for cancellation and timeout
- Consider batch insert for ProgressionLog entries
- Priority field in ProgramProgression determines evaluation order

## Dependencies
- Blocks: 012 (Manual trigger API uses this service)
- Blocked by: 002, 004, 005, 006, 007 (Log schema, interface, trigger types, implementations)
- Related: ERD-003 (State advancement fires triggers), ERD-001 (LiftMax updates)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/004-progression-rules/erd.md
- State Advancement: ERD-003 ticket 015-state-advancement-api.md

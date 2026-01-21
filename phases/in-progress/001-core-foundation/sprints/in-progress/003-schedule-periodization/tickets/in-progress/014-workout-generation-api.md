# 014: Workout Generation API

## ERD Reference
Implements: REQ-GEN-001, REQ-GEN-002, REQ-GEN-003

## Description
Implement the workout generation endpoint that returns resolved prescriptions for a user based on their program state and schedule position.

## Context / Background
This is the core API function - returning "today's training" with fully resolved weights based on program state, lookups, and user maxes. This ties together all schedule entities with the prescription system.

## Acceptance Criteria
- [ ] GET /users/{userId}/workout - generate workout
  - Uses current program state if no overrides
  - Optional query params: date, weekNumber, daySlug
  - Returns resolved prescriptions for specified day
- [ ] GET /users/{userId}/workout/preview - preview future workouts
  - Query params: week, day (required)
  - Does not require state advancement
  - Returns resolved prescriptions for specified position
- [ ] Workout includes: userId, programId, cycleIteration, weekNumber, daySlug, date, exercises[]
- [ ] Each exercise includes: prescriptionId, lift, sets[], notes, restSeconds
- [ ] Each set includes: setNumber, weight (resolved), targetReps, isWorkSet
- [ ] Week number passed to WeeklyLookup during resolution
- [ ] Day identifier passed to DailyLookup during resolution
- [ ] Resolved weights reflect all lookup modifications
- [ ] NFR-001: Generation completes in <500ms (p95)
- [ ] Unit tests with >80% coverage
- [ ] Integration tests with various lookup scenarios

## Technical Notes
- Resolution chain:
  1. Get user's program state
  2. Determine day from state (or override)
  3. Get day's prescriptions in order
  4. For each prescription, resolve with lookup context
  5. Return assembled workout
- Lookup context: {weekNumber, daySlug}
- If user not enrolled, return 400 or 404
- Preview does not modify state

## Dependencies
- Blocks: None
- Blocked by: 007-013 (All entity APIs), 011 (Lookup integration)
- Related: ERD-002 Prescription resolution

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/003-schedule-periodization/erd.md
- API Response Format: See ERD Section 5 "Generated Workout Response"

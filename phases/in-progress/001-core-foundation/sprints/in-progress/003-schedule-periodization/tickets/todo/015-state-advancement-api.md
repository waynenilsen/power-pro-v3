# 015: State Advancement API

## ERD Reference
Implements: REQ-STATE-002, REQ-STATE-003

## Description
Implement the state advancement endpoint that progresses a user through their program, including cycle completion detection.

## Context / Background
After completing a workout, users advance through their program. This endpoint advances to the next day, wrapping to the next week, wrapping to the next cycle. Cycle completion is detected and can trigger progressions.

## Acceptance Criteria
- [ ] POST /users/{userId}/program-state/advance - advance user state
  - Advances to next day within week
  - If end of week, advances to next week
  - If end of cycle, wraps to week 1 and increments cycle iteration
  - Returns new state
- [ ] Cycle completion detection
  - Flag/event when advancing from last week to week 1
  - Returned in response: { cycleCompleted: boolean }
  - Can be used as progression trigger
- [ ] State advancement is atomic (no partial updates)
- [ ] NFR-002: Advancement completes in <100ms (p95)
- [ ] NFR-004: Atomic state advancement (transaction)
- [ ] NFR-005: Reliable cycle completion detection
- [ ] Unit tests with >80% coverage
- [ ] Integration tests for state transitions
- [ ] Edge case tests: end of week, end of cycle

## Technical Notes
- Advancement logic:
  1. Get current state (week, day index, cycle iteration)
  2. Increment day index
  3. If day index >= days in week, day index = 0, increment week
  4. If week > cycle length, week = 1, increment cycle iteration
  5. Save and return
- Cycle completion: when week goes from N to 1
- Atomic: use database transaction
- Consider: should advancement be idempotent? (probably not - each call advances)

## Dependencies
- Blocks: None
- Blocked by: 013 (User enrollment), 006 (UserProgramState schema)
- Related: 014 (Workout generation uses state)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/003-schedule-periodization/erd.md

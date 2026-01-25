# 002: Enrollment Status Aggregation

## ERD Reference
Implements: REQ-DASH-003

## Description
Implement the enrollment section aggregation that returns the user's current program enrollment status, including program name, cycle position, and week position.

## Context / Background
The enrollment section tells users where they are in their training program. It combines data from enrollment, program, cycle, and week tables to provide a complete picture of program status.

## Acceptance Criteria
- [ ] Implement `AggregateEnrollment(ctx context.Context, userID uuid.UUID) (*EnrollmentSummary, error)`
- [ ] Return nil if user has no active enrollment
- [ ] EnrollmentSummary contains: Status, ProgramName, CycleIteration, CycleStatus, WeekNumber, WeekStatus
- [ ] Status is "ACTIVE" or "PAUSED"
- [ ] CycleStatus and WeekStatus are "IN_PROGRESS" or "COMPLETED"
- [ ] CycleIteration is 1-indexed (first cycle = 1)
- [ ] WeekNumber is 1-indexed within the current cycle
- [ ] Query is efficient (single query with joins preferred)
- [ ] Unit tests covering: active enrollment, paused enrollment, no enrollment, completed cycle

## Technical Notes
- Query should join enrollment → program → current_cycle → current_week
- Consider whether this belongs in enrollment service or as standalone function
- Use existing Enrollment, Program, Cycle, Week domain models
- Handle edge cases: enrollment exists but cycle/week not started

## Dependencies
- Blocks: 001 (service needs this aggregation)
- Blocked by: Existing enrollment domain model

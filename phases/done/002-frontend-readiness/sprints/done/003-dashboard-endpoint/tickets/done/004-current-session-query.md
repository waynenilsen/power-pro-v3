# 004: Current Session Query

## ERD Reference
Implements: REQ-DASH-005

## Description
Implement the current session detection that identifies if the user has an in-progress workout session and returns its details.

## Context / Background
When a user starts a workout but hasn't finished, they have an "active" session. The dashboard should surface this immediately so users can resume their workout rather than starting a new one.

## Acceptance Criteria
- [ ] Implement `GetCurrentSession(ctx context.Context, userID uuid.UUID) (*SessionSummary, error)`
- [ ] Return nil if no active session
- [ ] SessionSummary contains: SessionID, DayName, StartedAt, SetsCompleted, TotalSets
- [ ] SessionID is the UUID of the active session
- [ ] StartedAt in time.Time (will be serialized to ISO 8601)
- [ ] SetsCompleted is count of completed sets in this session
- [ ] TotalSets is total sets planned for this workout
- [ ] Only one active session per user enforced at query level
- [ ] Unit tests covering: active session, no session, session with partial completion

## Technical Notes
- Active session = session where completed_at is NULL
- May need to join with session_sets or similar to count completed sets
- TotalSets comes from program template for the workout day
- Consider adding index on user_id + completed_at for performance

## Dependencies
- Blocks: 001 (service needs this query)
- Blocked by: Existing WorkoutSession domain model

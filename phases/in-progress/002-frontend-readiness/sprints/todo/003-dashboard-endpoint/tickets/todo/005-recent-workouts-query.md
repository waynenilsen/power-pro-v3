# 005: Recent Workouts Query

## ERD Reference
Implements: REQ-DASH-006

## Description
Implement the recent workouts query that returns the user's last N completed workouts with summary information.

## Context / Background
Users want to see their recent workout history on the home screen. This provides context about their training consistency and allows them to verify that workouts were logged correctly.

## Acceptance Criteria
- [ ] Implement `GetRecentWorkouts(ctx context.Context, userID uuid.UUID, limit int) ([]WorkoutSummary, error)`
- [ ] Return empty slice if no completed workouts
- [ ] WorkoutSummary contains: Date, DayName, SetsCompleted
- [ ] Date is date only (YYYY-MM-DD), not full timestamp
- [ ] DayName is human-readable (e.g., "Intensity Day")
- [ ] SetsCompleted is count of sets logged in that session
- [ ] Limit parameter controls max results (default 5)
- [ ] Ordered by date descending (most recent first)
- [ ] Only includes completed sessions (completed_at is NOT NULL)
- [ ] Unit tests covering: multiple workouts, no workouts, respects limit

## Technical Notes
- Query should be: WHERE user_id = ? AND completed_at IS NOT NULL ORDER BY completed_at DESC LIMIT ?
- May need to extract date from completed_at timestamp
- Join with program day to get DayName
- Count sets from session_sets or similar table

## Dependencies
- Blocks: 001 (service needs this query)
- Blocked by: Existing WorkoutSession domain model

# 008: Dashboard E2E Tests

## ERD Reference
Implements: All REQs (REQ-DASH-001 through REQ-DASH-008)

## Description
Create comprehensive end-to-end tests for the dashboard endpoint covering all scenarios: empty state, active enrollment, mid-workout, and authorization.

## Context / Background
E2E tests validate the entire stack from HTTP request to database queries and back. These tests ensure the dashboard endpoint works correctly in real-world scenarios and catches integration issues between layers.

## Acceptance Criteria
- [ ] Test: Unauthenticated request returns 401
- [ ] Test: Authenticated user can access their own dashboard
- [ ] Test: Authenticated user cannot access another user's dashboard (403)
- [ ] Test: Non-existent user returns 404
- [ ] Test: User with no enrollment returns null enrollment, empty arrays
- [ ] Test: User with active enrollment returns full enrollment section
- [ ] Test: User with active session returns currentSession populated
- [ ] Test: User with active session has null nextWorkout
- [ ] Test: User with completed workouts returns recentWorkouts
- [ ] Test: User with training maxes returns currentMaxes
- [ ] Test: Weight values respect user's weight_unit preference
- [ ] Test: Response structure matches expected JSON format
- [ ] Test: Response time is within acceptable bounds (< 200ms)
- [ ] All tests pass in CI

## Technical Notes
- Use existing E2E test infrastructure
- Create test fixtures for various states (enrolled user, mid-workout user, etc.)
- May need to seed test data in setup phase
- Consider table-driven tests for authorization scenarios
- Measure and assert on response time for performance requirement

## Test Scenarios

1. **Empty State**
   - New user with no activity
   - Expected: all sections null or empty

2. **Active Enrollment**
   - User enrolled in Texas Method, Week 2 of Cycle 1
   - Expected: enrollment populated, next workout available

3. **Mid-Workout**
   - User started Volume Day, completed 3 of 10 sets
   - Expected: currentSession populated, nextWorkout null

4. **Historical Data**
   - User completed 7 workouts, has 3 maxes set
   - Expected: recentWorkouts has 5 entries, currentMaxes has 3 entries

5. **Authorization**
   - User A tries to access User B's dashboard
   - Expected: 403 Forbidden

## Dependencies
- Blocks: None (final ticket)
- Blocked by: 007 (needs endpoint to test)

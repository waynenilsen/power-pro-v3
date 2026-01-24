# E2E Tests for TotalReps and Building the Monolith

Create end-to-end integration tests demonstrating the TotalReps scheme in action, particularly for 5/3/1 Building the Monolith style programming.

## Building the Monolith Context

From 5/3/1 Building the Monolith:
- Main lifts follow standard 5/3/1 (already supported)
- Accessory work has total rep targets: "100 chin-ups", "100-200 dips"
- User distributes reps however they want across sets

## E2E Test Scenarios

### 1. Basic TotalReps Flow
- Create prescription with TotalReps scheme (100 reps target)
- Start session, get first set
- Log sets with varying rep counts (15, 12, 10, 8, etc.)
- Verify session tracks cumulative progress
- Verify termination when target reached

### 2. Exact Target Achievement
- Target: 50 reps
- Log: 15, 15, 15, 5 = 50 exactly
- Verify clean termination

### 3. Overshoot Scenario
- Target: 50 reps
- Log: 20, 20, 15 = 55 (over target)
- Verify terminates immediately on exceeding target

### 4. Max Sets Safety
- Target: 1000 reps (unrealistic)
- Max sets: 5
- Log 5 sets
- Verify terminates due to max sets, not reps

### 5. Building the Monolith Day Simulation
- Create a session with:
  - Squat 5/3/1 sets (existing scheme)
  - Press 5/3/1 sets (existing scheme)
  - Chin-ups: 100 total reps
  - Dips: 100-200 total reps (use 100 minimum)
  - Face pulls or band pull-aparts: 100 total reps
- Execute through the session
- Verify all prescriptions can complete

## Test Location

`internal/integration/totalreps_test.go` or similar E2E test directory

## Acceptance Criteria

- [ ] E2E tests cover full session workflow
- [ ] Tests demonstrate Building the Monolith accessory patterns
- [ ] All integration tests pass
- [ ] Tests are well-documented with program context

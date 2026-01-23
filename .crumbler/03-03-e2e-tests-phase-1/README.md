# E2E Tests for Phase 1 Programs

## Objective
Create end-to-end tests that validate each Phase 1 program can be fully configured and executed through the API.

## Phase 1 Programs Requiring E2E Tests

### 1. Starting Strength
- A/B rotation
- Fixed 3x5 sets
- LinearProgression per session

### 2. Bill Starr 5x5
- Heavy/Light/Medium days (DailyLookup)
- Ramp sets
- LinearProgression per week

### 3. Wendler 5/3/1 BBB
- WeeklyLookup for 4-week wave
- CycleProgression

### 4. Sheiko Beginner
- Many Fixed sets at various percentages
- No autoregulation needed

### 5. Greg Nuckols High Frequency
- DailyLookup + WeeklyLookup
- CycleProgression

## Test Structure
Each E2E test should:
1. Create a user
2. Create the program with all required entities (lifts, prescriptions, days, weeks, cycles, lookups)
3. Set up user's lift maxes
4. Simulate a full training cycle
5. Verify prescription resolution works correctly
6. Verify progression rules apply correctly

## Location
Tests should go in `internal/api/e2e/` or similar appropriate location.

## Acceptance Criteria
- [ ] Starting Strength e2e test passing
- [ ] Bill Starr 5x5 e2e test passing
- [ ] Wendler 5/3/1 BBB e2e test passing
- [ ] Sheiko Beginner e2e test passing
- [ ] Greg Nuckols High Frequency e2e test passing
- [ ] All tests runnable via `go test`

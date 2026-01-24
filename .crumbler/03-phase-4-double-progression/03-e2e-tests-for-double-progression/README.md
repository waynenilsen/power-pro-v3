# E2E Tests for Double Progression

Create end-to-end tests demonstrating double progression in action, specifically for Reddit PPL 6-Day program accessory work.

## Implementation

Create `/home/wayne/git/power-pro-v3/internal/integration/double_progression_test.go`:

### Test Scenarios

1. **Basic Double Progression Cycle**:
   - User performs 3x8-12 curls at 50lb
   - Session 1: 8, 8, 8 reps - no progression (not at ceiling)
   - Session 2: 10, 10, 10 reps - no progression (not at ceiling)
   - Session 3: 12, 12, 12 reps - progression triggered, weight increases to 55lb

2. **Partial Ceiling Hit**:
   - User hits ceiling on some but not all sets
   - Test behavior (should NOT progress until all target sets hit ceiling)

3. **Integration with RepRangeSetScheme**:
   - Create prescription with RepRangeSetScheme (3 sets, 8-12 reps)
   - Associate DoubleProgression
   - Verify GenerateSets produces correct output
   - Verify progression triggers correctly

4. **Reddit PPL Accessory Pattern**:
   - Simulate accessory exercise (e.g., bicep curls)
   - 3 sets x 8-12 reps
   - Progress weight when all sets hit 12 reps

## Acceptance Criteria

- E2E tests demonstrate complete double progression cycles
- Tests verify rep range set generation
- Tests verify progression triggers at ceiling
- Tests verify weight increases correctly
- All tests pass with `go test ./internal/integration/... -v`

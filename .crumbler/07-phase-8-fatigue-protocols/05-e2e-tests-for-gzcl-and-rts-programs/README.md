# E2E Tests for GZCL and RTS Programs

Create end-to-end integration tests demonstrating the fatigue protocols work correctly for real programs.

## Programs to Test

### GZCL Compendium (VDIP Pattern)

Test the T1 tier using MRS:
1. Create program with T1 using MRS (3 Max Rep Sets)
2. Simulate workout:
   - Set 1: 10 reps
   - Set 2: 8 reps
   - Set 3: 6 reps (total 24 reps - done after 3 MRS)
3. Verify termination after 3 MRS blocks
4. Verify total rep tracking

Test the T3 tier using MRS:
1. Create program with T3 using MRS (4 Max Rep Sets)
2. Simulate reaching target reps
3. Verify different MRS count works

### RTS Intermediate (Load Drop Pattern)

Test the Load Drop method using FatigueDrop:
1. Create program with FatigueDrop scheme
2. Simulate workout:
   - Set 1: 315 lbs x 3 @ RPE 8.0
   - Set 2: 299 lbs x 3 @ RPE 8.5
   - Set 3: 284 lbs x 3 @ RPE 9.0
   - Set 4: 270 lbs x 3 @ RPE 10.0 (stop)
3. Verify weight drops correctly
4. Verify RPE threshold triggers termination

Test the Repeat Sets method:
1. Create program repeating at same weight
2. Simulate RPE increasing over sets
3. Verify termination when RPE exceeds target

## Test Structure

Create `internal/integration/fatigue_protocols_test.go`:

```go
func TestGZCL_VDIP_T1_ThreeMRS(t *testing.T)
func TestGZCL_VDIP_T3_FourMRS(t *testing.T)
func TestGZCL_MRS_FailureTermination(t *testing.T)
func TestRTS_LoadDrop_RPETermination(t *testing.T)
func TestRTS_LoadDrop_MaxSetsLimit(t *testing.T)
func TestRTS_RepeatSets_RPEIncrease(t *testing.T)
```

## What to Verify

For each test:
1. **Setup**: Program created with correct scheme
2. **Execution**: Simulated sets logged correctly
3. **Termination**: Stopped at right condition
4. **Data**: All logged sets recorded with correct weights/reps/RPE
5. **Progression**: (If applicable) Next session adjusts correctly

## Dependencies

- Tasks 1-4 must be complete
- Requires working FatigueDrop and MRS schemes
- Requires working session handling

## Acceptance Criteria

- [ ] GZCL VDIP T1 test (3 MRS) passing
- [ ] GZCL VDIP T3 test (4 MRS) passing
- [ ] MRS failure termination test passing
- [ ] RTS Load Drop RPE termination test passing
- [ ] RTS Load Drop max sets safety test passing
- [ ] At least 6 E2E tests total
- [ ] Tests document expected behavior clearly
- [ ] Tests can serve as examples for users

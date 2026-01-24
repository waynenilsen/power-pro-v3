# RTS Intermediate E2E Tests

Implement end-to-end tests demonstrating partial RTS Intermediate program support.

## Background

RTS Intermediate uses RPE-based load calculation. With RPETarget LoadStrategy and RPEChart Lookup implemented, we can demonstrate basic RTS workout generation.

**Note:** This is a "partial" implementation because RTS Intermediate also requires:
- Fatigue management (Phase 8)
- Volume targets (Phase 9)

This phase focuses on: **RPE-based weight calculation for prescribed sets**

## Test Scenarios

### 1. Basic RPE-Based Workout Generation

```go
func TestRTSIntermediateBasicWorkout(t *testing.T) {
    // Setup:
    // - Create user with squat 1RM = 365 lbs
    // - Create prescription with RPETarget strategy (5 reps @ RPE 9)

    // Expected:
    // - RPE chart: 5 reps @ RPE 9 = 80% (0.80)
    // - Calculated weight: 365 × 0.80 = 292 → rounds to 290 lbs

    // Verify:
    // - Workout generation produces correct weight
    // - Generated sets have correct target reps (5)
}
```

### 2. Multiple RPE Prescriptions in Day

```go
func TestRTSIntermediateMultiplePrescriptions(t *testing.T) {
    // Setup:
    // - Squat: 4 reps @ RPE 9 (82% = 300 lbs from 365 1RM)
    // - Bench: 5 reps @ RPE 8 (77% = 230 lbs from 300 1RM)

    // Verify:
    // - Both prescriptions resolve correctly
    // - Different lifts use their own 1RMs
}
```

### 3. RPE with Different Rounding

```go
func TestRTSIntermediateRounding(t *testing.T) {
    // Test NEAREST, UP, DOWN rounding with RPE calculations
    // Verify edge cases where rounding affects final weight
}
```

### 4. Log Set with RPE

```go
func TestRTSIntermediateLogSetWithRPE(t *testing.T) {
    // Generate workout with target RPE 8
    // Log set with actual RPE 9 (felt harder)
    // Verify RPE is stored correctly

    // This data will be used in Phase 7 for e1RM calculations
}
```

### 5. Custom RPE Chart

```go
func TestRTSIntermediateCustomRPEChart(t *testing.T) {
    // Some users may have custom RPE charts based on their training history
    // Test that non-default RPE charts work correctly
}
```

## File Location

Create: `internal/integration/rts_intermediate_test.go`

## Test Structure

Follow existing integration test patterns:
- `internal/integration/greyskull_lp_test.go`
- `internal/integration/juggernaut_test.go`

```go
func TestRTSIntermediate(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    // Setup test database
    // Register RPETarget strategy
    // Create test fixtures
    // Run scenarios
}
```

## Acceptance Criteria

- [ ] E2E test file created at `internal/integration/rts_intermediate_test.go`
- [ ] Test: RPE-based weight calculation works correctly
- [ ] Test: Multiple RPE prescriptions in a day resolve correctly
- [ ] Test: Logging sets with actual RPE works
- [ ] Test: Different rounding directions work with RPE
- [ ] All tests pass with `go test ./internal/integration/... -v`
- [ ] Tests demonstrate RTS Intermediate "partial" capability

## Not In Scope (Future Phases)

- Fatigue percentage drops (Phase 8)
- Volume load tracking (Phase 9)
- Rep drop / Load drop methods (Phase 8)
- Auto-regulation based on logged RPE (future)

# Calgary Barbell E2E Tests

Demonstrate E1RM, FindRM, and RelativeTo working together in real program scenarios.

## Programs to Test

### GZCL Jacked & Tan 2.0
- Week 1: Find 10RM, then 3×10 @ 85% of found weight
- Week 2: Find 8RM, then 3×8 @ 85%
- Uses: FindRM + RelativeTo

### Calgary Barbell 8-Week
- RPE-based top set → calculate E1RM
- Back-off sets at % of top set
- Uses: RPETarget + E1RM calculation + RelativeTo

## Test Structure

Create `internal/integration/e1rm_test.go` (or add to existing integration tests)

## Test Scenarios

### Test 1: GZCL FindRM Flow
```go
func TestGZCL_FindRM_WithBackoffs(t *testing.T) {
    // 1. Create prescription with FindRM(10) + RelativeTo(0, 85%)
    // 2. Log the FindRM set: 315 × 10
    // 3. Calculate RelativeTo: 315 × 0.85 = 267.75 → 270
    // 4. Verify back-off prescription shows 270 lbs
}
```

### Test 2: Calgary Barbell RPE→E1RM→Backoff
```go
func TestCalgaryBarbell_RPE_E1RM_Backoffs(t *testing.T) {
    // 1. Create RPETarget prescription for top single @ RPE 8
    // 2. Log the set: 365 × 1 @ RPE 8
    // 3. Calculate E1RM: 365 / 0.92 = 396.7 → 397.5
    // 4. Back-offs at 80% of top set: 365 × 0.80 = 292
    // 5. Verify calculations match expected
}
```

### Test 3: E1RM Storage and Retrieval
```go
func TestE1RM_StoredAsLiftMax(t *testing.T) {
    // 1. Calculate E1RM from logged set
    // 2. Store as LiftMax with type E1RM
    // 3. Use in subsequent PercentOf prescription
    // 4. Verify weight calculation uses stored E1RM
}
```

## Acceptance Criteria

- [ ] GZCL FindRM→RelativeTo flow works end-to-end
- [ ] Calgary Barbell RPE→E1RM→RelativeTo flow works
- [ ] E1RM can be stored and used in future prescriptions
- [ ] All edge cases handled (missing data, invalid state)
- [ ] Tests are integration-level (touch multiple domains)

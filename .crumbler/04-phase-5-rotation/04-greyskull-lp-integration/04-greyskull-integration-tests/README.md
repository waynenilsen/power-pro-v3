# GreySkull Integration Tests

Write comprehensive integration tests for GreySkull LP covering the full workout generation flow.

## Overview

Integration tests should verify that all GreySkull LP components work together correctly:
- Week variant rotation (A/B pattern)
- Set scheme generation (2x5 + 1x5+)
- Progression application (deload, standard, double)
- Load calculation

## Test File Location
`internal/integration/greyskull_lp_test.go`

## Test Scenarios

### 1. Week Variant Rotation Tests
```go
func TestGreySkullWeekVariantRotation(t *testing.T) {
    testCases := []struct {
        weekNumber  int
        dayPosition int
        expected    string
    }{
        {1, 1, "A"}, {1, 2, "B"}, {1, 3, "A"},
        {2, 1, "B"}, {2, 2, "A"}, {2, 3, "B"},
        {3, 1, "A"}, {3, 2, "B"}, {3, 3, "A"},
        {4, 1, "B"}, {4, 2, "A"}, {4, 3, "B"},
    }
    // ...
}
```

### 2. Set Scheme Generation Tests
```go
func TestGreySkullSetSchemeGeneration(t *testing.T) {
    testCases := []struct {
        name       string
        weight     float64
        fixedSets  int
        fixedReps  int
        amrapSets  int
        expected   []Set
    }{
        {
            name: "main lift 2x5+1x5+",
            weight: 135.0,
            fixedSets: 2, fixedReps: 5, amrapSets: 1,
            expected: []Set{
                {Weight: 135, Reps: 5, SetType: "WORK"},
                {Weight: 135, Reps: 5, SetType: "WORK"},
                {Weight: 135, Reps: 5, SetType: "AMRAP"},
            },
        },
    }
    // ...
}
```

### 3. Progression Tests
```go
func TestGreySkullProgression(t *testing.T) {
    testCases := []struct {
        name           string
        currentWeight  float64
        repsPerformed  int
        increment      float64
        minReps        int
        doubleThreshold int
        deloadPercent  float64
        expectedWeight float64
        expectedReason string
    }{
        {"deload on failure", 135, 3, 2.5, 5, 10, 0.10, 121.5, "deload"},
        {"standard increment", 135, 7, 2.5, 5, 10, 0.10, 137.5, "standard"},
        {"double increment", 135, 12, 2.5, 5, 10, 0.10, 140, "double"},
    }
    // ...
}
```

### 4. Full Workout Flow Tests
```go
func TestGreySkullFullWorkoutFlow(t *testing.T) {
    // Test complete flow:
    // 1. Create program with GreySkull configuration
    // 2. Create cycle with weeks
    // 3. Resolve prescriptions for each day
    // 4. Verify correct exercises appear based on variant
    // 5. Verify correct sets generated
    // 6. Apply progression after AMRAP result
    // 7. Verify next session has correct weight
}
```

### 5. Multi-Week Cycle Tests
```go
func TestGreySkullMultiWeekCycle(t *testing.T) {
    // Test 6-week cycle:
    // - Week 1-6 with alternating A/B pattern
    // - Progressive overload with standard increments
    // - One deload scenario
    // - One double progression scenario
}
```

### 6. Day Template Selection Tests
```go
func TestGreySkullDayTemplateSelection(t *testing.T) {
    testCases := []struct {
        variant        string
        expectedLifts  []string
    }{
        {"A", []string{"bench-press", "barbell-row", "squat", "tricep-extension", "ab-rollout"}},
        {"B", []string{"overhead-press", "chin-up", "deadlift", "bicep-curl", "shrug"}},
    }
    // ...
}
```

## Test Fixtures

### Program Fixture
```go
func createGreySkullProgram() *program.Program {
    return &program.Program{
        ID:   "greyskull-lp",
        Name: "GreySkull LP",
        Slug: "greyskull-lp",
        // ...
    }
}
```

### Cycle Fixture
```go
func createGreySkullCycle() *cycle.Cycle {
    return &cycle.Cycle{
        ID:          "greyskull-cycle",
        Name:        "GreySkull LP Cycle",
        LengthWeeks: 0, // Ongoing, no fixed length
        // ...
    }
}
```

### Week Fixtures
```go
func createGreySkullWeeks(count int) []*week.Week {
    weeks := make([]*week.Week, count)
    for i := 0; i < count; i++ {
        weekNum := i + 1
        // Variant is determined at resolution time, not stored
        weeks[i] = &week.Week{
            ID:         fmt.Sprintf("week-%d", weekNum),
            WeekNumber: weekNum,
            CycleID:    "greyskull-cycle",
        }
    }
    return weeks
}
```

## Tasks

1. Create `internal/integration/greyskull_lp_test.go`
2. Write week variant rotation tests
3. Write set scheme generation tests
4. Write progression tests (deload, standard, double)
5. Write full workout flow integration test
6. Write multi-week cycle test
7. Write day template selection tests
8. Create test fixtures and helpers

## Dependencies

This task depends on:
- Task 1: GreySkull Progression Strategy
- Task 2: Week Variant Rotation Logic
- Task 3: Composite Set Scheme

## Acceptance Criteria

- [ ] All rotation tests pass
- [ ] All set scheme tests pass
- [ ] All progression tests pass
- [ ] Full workflow integration test passes
- [ ] Multi-week cycle test passes
- [ ] Day template selection tests pass
- [ ] Tests are table-driven and comprehensive
- [ ] Test coverage > 80% for greyskull package

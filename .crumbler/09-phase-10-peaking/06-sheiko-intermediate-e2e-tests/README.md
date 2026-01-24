# Sheiko Intermediate E2E Tests

Comprehensive end-to-end tests for Sheiko Intermediate program with peaking.

## Test Scenarios

### Test 1: Full 13-Week Program Execution
1. Seed Sheiko Intermediate program with 3 phases
2. Enroll user with meet date 13 weeks out
3. Verify schedule type is "days_out"
4. Generate sessions for each week
5. Verify correct phase at each point:
   - Weeks 1-4: Prep 1 (higher volume, moderate intensity)
   - Weeks 5-8: Prep 2 (max testing, 100-105% loads)
   - Weeks 9-13: Competition (taper, volume reduction)
6. Verify taper multiplier reduces in final weeks

### Test 2: Phase Transitions
1. Set meet date 5 weeks out (start in Competition phase)
2. Verify immediate Competition phase detection
3. Generate sessions, verify taper applied

### Test 3: Meet Date Changes
1. Enroll with meet date 10 weeks out
2. Change meet date to 6 weeks out
3. Verify phase recalculates correctly

### Test 4: Opener Practice
1. In final week before meet
2. Verify opener-specific prescriptions (if modeled)
3. Light technical work only

## Sheiko Intermediate Structure

| Phase | Weeks | Primary Intensity | Volume | Focus |
|-------|-------|-------------------|--------|-------|
| Prep 1 | 1-4 | 70-75% | High | Base building |
| Prep 2 | 5-8 | 80-90% + 100-105% max tests | Medium-High | Strength peaks |
| Comp | 9-13 | 85-95% → 50-70% | Decreasing | Taper to meet |

## Acceptance Criteria

- [ ] E2E test demonstrates complete 13-week cycle
- [ ] All phase transitions verified
- [ ] Taper behavior verified
- [ ] Load calculations correct for each phase
- [ ] All tests pass

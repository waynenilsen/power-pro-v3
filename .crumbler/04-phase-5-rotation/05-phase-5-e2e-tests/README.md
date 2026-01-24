# Phase 5 E2E Tests

Comprehensive end-to-end tests demonstrating all three unlocked programs work correctly.

## Test Files

- `tests/integration/phase5_cap3_rotation_test.go`
- `tests/integration/phase5_juggernaut_rotation_test.go`
- `tests/integration/phase5_greyskull_rotation_test.go`
- `tests/acceptance/phase5_rotation_acceptance_test.go`

## nSuns CAP3 Tests

### Test: Full 3-Week Rotation Cycle
1. Create CAP3 program with rotation lookup
2. Enroll user at rotation position 0
3. Generate Week 1 workout → verify Deadlift gets AMRAP percentages
4. Verify Squat/Bench get volume percentages
5. Advance to Week 2 → verify Squat gets AMRAP
6. Advance to Week 3 → verify Bench gets AMRAP
7. Complete cycle → verify rotation position resets

### Test: Multiple Cycles
1. Run through 2 complete 3-week cycles
2. Verify rotation advances correctly (position cycles 0→1→2→0→1→2)
3. Verify TM adjustments based on AMRAP performance

## Inverted Juggernaut Tests

### Test: Full 16-Week Cycle
1. Create Juggernaut program
2. Test Week 1 (10s wave, accumulation): 9 sets @ 60%
3. Test Week 5 (8s wave, accumulation): 7 sets @ 65%
4. Test Week 9 (5s wave, accumulation): 5 sets @ 70%
5. Test Week 13 (3s wave, accumulation): 6 sets @ 75%
6. Test deload weeks (4, 8, 12, 16)
7. Verify TM progression after each wave's realization week

### Test: Wave Transitions
1. Complete 10s wave (weeks 1-4)
2. Verify transition to 8s wave
3. Verify volume set count changes (9→7)
4. Verify base percentage changes (60%→65%)

### Test: 5/3/1 Overlay
1. Week 1 of any wave: 65/75/85%
2. Week 2: 70/80/90%
3. Week 3: 75/85/95% with AMRAP
4. Week 4: 40/50/60% deload

## GreySkull LP Tests

### Test: A/B Rotation
1. Create GreySkull program with alternating weeks
2. Week 1 (Variant A): Day 1 = Bench, Day 2 = OHP, Day 3 = Bench
3. Week 2 (Variant B): Day 1 = OHP, Day 2 = Bench, Day 3 = OHP
4. Verify correct exercises appear on correct days

### Test: AMRAP Progression
1. User performs AMRAP set
2. If 5-9 reps: standard progression (+2.5)
3. If 10+ reps: double progression (+5)
4. If <5 reps: 10% deload

### Test: Complete Training Block
1. Simulate 6 weeks of training
2. Track weight progression
3. Verify deloads occur appropriately
4. Verify double progression triggers when applicable

## Acceptance Tests

### Business Requirements Validation

#### nSuns CAP3
- [ ] Each lift gets AMRAP focus once per 3-week cycle
- [ ] Non-focus lifts use medium/volume percentages
- [ ] Rotation cycles correctly (0→1→2→0)
- [ ] TM adjusts based on AMRAP performance

#### Inverted Juggernaut
- [ ] 4 waves cycle correctly (10s→8s→5s→3s)
- [ ] Volume sets match wave target (9/7/5/6)
- [ ] Base percentages match wave (60/65/70/75%)
- [ ] 5/3/1 percentages apply correctly
- [ ] 16-week cycle completes and resets

#### GreySkull LP
- [ ] Days alternate A/B/A, B/A/B correctly
- [ ] AMRAP final sets work on all main lifts
- [ ] Double progression triggers on 10+ reps
- [ ] 10% deload triggers on failure
- [ ] 3-week cycle repeats correctly

## Acceptance Criteria

- [ ] All CAP3 E2E tests pass
- [ ] All Juggernaut E2E tests pass
- [ ] All GreySkull E2E tests pass
- [ ] All acceptance criteria validated
- [ ] Test coverage adequate for rotation logic
- [ ] Tests document expected behavior clearly

# nSuns CAP3 Implementation

Implement the 3-week cyclical AMRAP rotation for nSuns CAP3 program.

## Program Mechanics

CAP3 uses a 3-week rotation where each main lift (Deadlift, Squat, Bench) gets AMRAP focus once per cycle:

| Week | Deadlift | Squat | Bench |
|------|----------|-------|-------|
| Week 1 | HIGH INTENSITY AMRAP | Medium Volume | Volume |
| Week 2 | Medium Volume | HIGH INTENSITY AMRAP | Medium Volume |
| Week 3 | Volume | Volume | HIGH INTENSITY AMRAP |

### High Intensity AMRAP Day Percentages
- Set 1: 79.5% TM x 6 reps
- Set 2: 83.5% TM x 4 reps
- Set 3 (AMRAP): 88.5% TM x 2+ reps

### Volume Day Percentages
- Sets 1-4: 73.5% TM x 4 reps
- Sets 5-6: 73.5% TM + 5lb x 4 reps
- Set 7 (AMRAP): 73.5% TM + 10lb x 4+ reps

### Medium Day Percentages
- Sets 1-4: 77% TM x 3 reps
- Sets 5-6: 77% TM + 5lb x 3 reps
- Sets 7-8 (AMRAP): 77% TM + 10lb x 3+ reps

## Tasks

### 1. Create CAP3 RotationLookup Fixture
Define the 3-position rotation:
- Position 0: Deadlift Focus
- Position 1: Squat Focus
- Position 2: Bench Focus

### 2. Create CAP3 WeeklyLookup Entries
Configure percentages for each lift based on whether it's the focus lift that week:
- Focus week: High intensity AMRAP percentages (79.5/83.5/88.5%)
- Medium week: Medium percentages (77%)
- Volume week: Volume percentages (73.5%)

### 3. Integrate Rotation in Prescription Resolution
Modify prescription resolution to:
1. Get current rotation position from UserProgramState
2. Check if this lift matches the focus lift for current position
3. Apply appropriate percentages based on focus state

### 4. Create CAP3 Test Fixtures
Build test helpers for:
- Program setup with rotation
- User enrollment
- Cycle/Week/Day structure
- Prescription configurations

### 5. Write Integration Tests
Test scenarios:
- Week 1: Deadlift gets AMRAP, others get volume/medium
- Week 2: Squat gets AMRAP, others get medium
- Week 3: Bench gets AMRAP, others get volume
- Rotation advances after 3-week cycle

## Acceptance Criteria

- [ ] CAP3 rotation lookup correctly maps positions to lifts
- [ ] Prescriptions resolve with correct percentages based on rotation
- [ ] Rotation advances after 3-week cycle completes
- [ ] Integration tests pass for complete rotation cycle
- [ ] AMRAP sets correctly configured for focus lift

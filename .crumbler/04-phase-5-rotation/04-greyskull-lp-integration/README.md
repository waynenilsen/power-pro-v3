# GreySkull LP Integration

Integrate the existing A/B rotation with AMRAP and DoubleProgression for GreySkull LP.

## Program Mechanics

GreySkull uses a 3-day A/B alternating pattern:
- **Week 1**: A, B, A
- **Week 2**: B, A, B
- **Week 3**: A, B, A (repeats)

### Day Templates
**Day A**: Bench Press, Squat, Tricep Extension, Ab Rollout
**Day B**: Overhead Press, Deadlift, Bicep Curl, Shrug

### Rep Scheme
All main lifts: 2x5 + 1x5+ (AMRAP final set)
Accessories: 3x12 (or AMRAP on final)

### Progression Rules
- AMRAP < 5 reps: 10% deload
- AMRAP 5-9 reps: +2.5 lbs (standard)
- AMRAP 10+ reps: +5 lbs (double progression)

## Current State

Week.Variant field already supports A/B alternation. What needs integration:

1. AMRAP final set on all exercises
2. DoubleProgression when hitting rep ceiling (10+)
3. Proper day template selection based on variant

## Tasks

### 1. Verify Week.Variant Integration
Confirm that Week.Variant correctly filters days:
- Variant "A" → Day A exercises
- Variant "B" → Day B exercises

### 2. Configure AMRAP SetScheme
Ensure AMRAP set scheme works for main lifts:
- 2 fixed sets of 5 reps
- 1 AMRAP set with min 5 reps

### 3. Integrate DoubleProgression
Use existing DoubleProgression from Phase 4:
- MinReps: 5, MaxReps: 10
- On hitting 10+ reps, double the weight increment

### 4. Configure Accessory Progression
Accessories use tolerance band progression:
- < 10 reps: 10% deload
- 10-14 reps: maintain
- 15+ reps: add increment

### 5. Create GreySkull Test Fixtures
Build fixtures for:
- Program with 3-week cycle
- Weeks with alternating variants (A, B, A, B, A, B, ...)
- Days with appropriate prescriptions
- Prescriptions with AMRAP and DoubleProgression

### 6. Write Integration Tests
Test scenarios:
- Week 1, Day 1 (Variant A): Bench day exercises
- Week 2, Day 1 (Variant B): OHP day exercises
- AMRAP triggers weight increase
- Double progression triggers on 10+ reps
- 10% deload on failure (<5 reps)

## Acceptance Criteria

- [ ] Week.Variant correctly selects day templates
- [ ] AMRAP final set configured for all main lifts
- [ ] DoubleProgression works (5-9 reps = +2.5, 10+ = +5)
- [ ] 10% deload triggers on failure
- [ ] Accessory tolerance band progression works
- [ ] Integration tests pass for full A/B rotation cycle

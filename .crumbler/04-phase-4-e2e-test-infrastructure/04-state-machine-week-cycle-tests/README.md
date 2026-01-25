# State Machine Week/Cycle Tests

## Overview

Create `state_machine_week_cycle_test.go` to verify week and cycle state transitions.

## Test Cases

### Test: First workout of week transitions week PENDING → IN_PROGRESS
- Enroll user (weekStatus starts PENDING)
- Start workout
- Verify weekStatus = "IN_PROGRESS"

### Test: Last workout of week transitions to COMPLETED with auto-advance
- Set up program with known days per week
- Complete all workouts in week
- Call advanceWeek
- Verify previous week marked COMPLETED
- Verify new week starts PENDING

### Test: Last workout of cycle triggers CYCLE_COMPLETED transition
- Set up program with known cycle length (e.g., 1 week for simplicity)
- Complete all workouts in cycle
- Advance past final week
- Verify enrollmentStatus = "BETWEEN_CYCLES"
- Verify cycleStatus = "COMPLETED"

### Test: Week status resets to PENDING on new week
- Complete a week
- Advance to next week
- Verify weekStatus = "PENDING"
- Verify currentWeek incremented

### Test: Cycle status resets to PENDING on new cycle
- Complete full cycle (reach BETWEEN_CYCLES)
- Start next cycle
- Verify cycleStatus = "PENDING"
- Verify weekStatus = "PENDING"
- Verify cycleIteration incremented
- Verify currentWeek = 1

## File Location

`internal/api/e2e/state_machine_week_cycle_test.go`

## Acceptance Criteria

- [ ] All 5 test cases implemented
- [ ] Tests create minimal programs with known week/cycle lengths
- [ ] Week transitions verified
- [ ] Cycle transitions verified
- [ ] Auto-advance behavior verified

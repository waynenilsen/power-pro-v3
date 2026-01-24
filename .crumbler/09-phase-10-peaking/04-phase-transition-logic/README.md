# Phase Transition Logic

Implement automatic phase transitions based on calendar/meet date.

## Tasks

1. Extend state advancement logic to check meet date
2. Implement `ShouldTransitionPhase(state, now)` function
3. Auto-transition between Prep1 → Prep2 → Competition phases
4. Update UserProgramState with `CurrentPhase` tracking (optional - may be derived)
5. Handle edge cases:
   - Meet date changed mid-program
   - Meet date passed
   - No meet date set (stay in rotation mode)
6. Add integration tests for phase transitions

## Transition Rules

- If `ScheduleType` is "rotation": use existing week-based progression
- If `ScheduleType` is "days_out": derive phase from meet date
- Phase is derived each time, not stored (simpler, no sync issues)

## Integration with Session Generation

When generating a session:
1. Check schedule type
2. If days_out, calculate current phase from meet date
3. Look up correct week within that phase
4. Apply taper multiplier if in Competition phase

## Acceptance Criteria

- [ ] Phase transitions happen at correct times
- [ ] Both rotation and days_out schedules work
- [ ] Edge cases handled gracefully
- [ ] Integration tests verify transitions

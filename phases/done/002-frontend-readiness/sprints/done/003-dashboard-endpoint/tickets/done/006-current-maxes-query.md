# 006: Current Maxes Query

## ERD Reference
Implements: REQ-DASH-007

## Description
Implement the current maxes query that returns the user's current training maxes for all tracked lifts, converted to user's preferred weight unit.

## Context / Background
Training maxes are fundamental to powerlifting programs - they determine the weights used in each workout. Users constantly reference their maxes and need them readily available on the home screen.

## Acceptance Criteria
- [ ] Implement `GetCurrentMaxes(ctx context.Context, userID uuid.UUID) ([]MaxSummary, error)`
- [ ] Return empty slice if no maxes set
- [ ] MaxSummary contains: Lift, Value, Type
- [ ] Lift is exercise name (e.g., "Squat", "Bench Press", "Deadlift")
- [ ] Value is weight in user's preferred unit (from profile weight_unit)
- [ ] Type is max type: "TRAINING_MAX", "ONE_REP_MAX", "ESTIMATED"
- [ ] Convert weights based on user's weight_unit preference
- [ ] Order by lift name alphabetically
- [ ] Only return current/active maxes (not historical)
- [ ] Unit tests covering: multiple maxes, no maxes, lb preference, kg preference

## Technical Notes
- Need to fetch user's weight_unit from profile (Sprint 002)
- Maxes may be stored in one unit and need conversion
- Conversion: 1 kg = 2.20462 lb (round to nearest 0.5 or configurable precision)
- Join with exercise table to get lift name
- Consider: should maxes be stored normalized (always kg) and converted on output?

## Dependencies
- Blocks: 001 (service needs this query)
- Blocked by: Existing TrainingMax domain model, Sprint 002 weight_unit

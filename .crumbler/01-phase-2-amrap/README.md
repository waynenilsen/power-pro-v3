# Phase 2: AMRAP

Implement AMRAP (As Many Reps As Possible) functionality and AMRAP-driven progression.

## Features to Implement

- **SetScheme: `AMRAP`** - Sets where user performs maximum reps
- **LoggedSet** - Record actual reps performed in a set
- **Progression: `AMRAPProgression`** - Weight adjustments based on AMRAP performance thresholds
- **Trigger: `AfterSet(amrap)`** - Fire progression rules after AMRAP sets

## Programs Unlocked

| Program | Why |
|---------|-----|
| Wendler 5/3/1 BBB | AMRAP final sets (5+, 3+, 1+), WeeklyLookup for 4-week wave, CycleProgression |
| Greg Nuckols High Frequency | AMAP sets in Week 3, DailyLookup + WeeklyLookup, CycleProgression |
| nSuns 5/3/1 LP 5-Day | Multiple AMRAP sets, complex AMRAPProgression thresholds |

## Acceptance Criteria

- Users can define AMRAP set schemes
- Actual reps are logged and persisted
- Progression rules fire based on AMRAP performance
- E2E tests demonstrate all three unlocked programs

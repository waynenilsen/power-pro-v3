# DaysOut Schedule Type

Implement the DaysOut schedule type that counts down to a meet date.

## Tasks

1. Create `internal/domain/schedule/daysout.go` package
2. Implement `DaysOutCalculator` that computes:
   - Days remaining until meet
   - Current phase based on meet proximity (Prep1, Prep2, Competition)
   - Week within current phase
3. Define phase durations (for Sheiko: Prep1=4wk, Prep2=4wk, Comp=5wk = 13 total)
4. Implement `GetCurrentPhase(meetDate, now)` function
5. Implement `GetDaysOut(meetDate, now)` function
6. Add unit tests for all calculations

## Phase Logic

For a 13-week program:
- **Prep 1**: Weeks 1-4 (91-64 days out)
- **Prep 2**: Weeks 5-8 (63-36 days out)
- **Competition**: Weeks 9-13 (35-0 days out)

## Acceptance Criteria

- [ ] DaysOut calculation is accurate
- [ ] Phase determination works correctly based on days remaining
- [ ] Week-within-phase calculation works
- [ ] Unit tests cover edge cases (meet day, day after meet, etc.)

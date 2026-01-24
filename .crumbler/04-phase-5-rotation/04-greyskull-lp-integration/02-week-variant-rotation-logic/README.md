# Week Variant Rotation Logic

Implement the A/B rotation pattern for GreySkull LP week variants.

## Overview

GreySkull LP uses a 3-day per week A/B alternating pattern. The variant determines which day template is used.

## Rotation Pattern

```
Week 1: Day 1 (A), Day 2 (B), Day 3 (A)
Week 2: Day 1 (B), Day 2 (A), Day 3 (B)
Week 3: Day 1 (A), Day 2 (B), Day 3 (A)
Week 4: Day 1 (B), Day 2 (A), Day 3 (B)
...
```

The pattern is: Odd weeks start with A, even weeks start with B.

## Day Templates

### Variant A (Bench Day)
- Bench Press: 2x5 + 1x5+ (AMRAP)
- Barbell Row: 2x5 + 1x5+ (AMRAP)
- Squat: 2x5 + 1x5+ (AMRAP)
- Tricep Extension: 3x12 (or AMRAP final)
- Ab Rollout: 3x10+

### Variant B (OHP Day)
- Overhead Press: 2x5 + 1x5+ (AMRAP)
- Chin-ups/Pull-ups: 2x5 + 1x5+ (AMRAP)
- Deadlift/Squat: 2x5 + 1x5+ (AMRAP)
- Bicep Curl: 3x12 (or AMRAP final)
- Shrug: 3x12 (or AMRAP final)

## Implementation

### File Location
`internal/domain/greyskull/rotation.go`

### Core Function
```go
// GetVariantForDay returns "A" or "B" based on week number and day-of-week position
func GetVariantForDay(weekNumber int, dayPosition int) string {
    // dayPosition is 1, 2, or 3 (for Mon, Wed, Fri typically)
    // Week 1: A, B, A (positions 1, 2, 3)
    // Week 2: B, A, B (positions 1, 2, 3)

    isOddWeek := weekNumber % 2 == 1
    isOddPosition := dayPosition % 2 == 1

    if isOddWeek == isOddPosition {
        return "A"
    }
    return "B"
}
```

### Alternative: Use Existing Week.Variant
The `Week.Variant` field already exists. This task may be about:
1. Creating helper functions to determine variant
2. Ensuring rotation lookup integrates with existing infrastructure

### Integration with Existing RotationLookup
Check if we can leverage `internal/domain/rotationlookup/` for this pattern.

## Tasks

1. Review existing `Week.Variant` implementation in `internal/domain/week/`
2. Review `RotationLookup` implementation for reuse opportunity
3. Create `internal/domain/greyskull/rotation.go` with rotation calculation
4. Create `GetVariantForDay(weekNumber, dayPosition) string` function
5. Create `GetDayTemplate(variant string) []string` function (returns lift slugs)
6. Write unit tests in `internal/domain/greyskull/rotation_test.go`

## Test Cases

- Week 1, Day 1 → "A"
- Week 1, Day 2 → "B"
- Week 1, Day 3 → "A"
- Week 2, Day 1 → "B"
- Week 2, Day 2 → "A"
- Week 2, Day 3 → "B"
- Week 10, Day 2 → "A" (even week, even position)
- Week 11, Day 1 → "A" (odd week, odd position)

## Acceptance Criteria

- [ ] Rotation logic correctly alternates A/B pattern
- [ ] Week.Variant field integration verified
- [ ] GetVariantForDay function implemented
- [ ] Day template mapping works
- [ ] All unit tests pass

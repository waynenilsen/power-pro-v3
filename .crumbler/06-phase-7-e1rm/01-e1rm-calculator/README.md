# E1RM Calculator

Implement a calculator that estimates 1RM from a performed set using the RPE chart.

## Formula

```
E1RM = Weight / RPEChart.GetPercentage(RepsPerformed, RPE)
```

Example: 315 lbs × 5 reps @ RPE 8 → 315 / 0.81 = 388.9 → rounded to 390 lbs

## Location

Create `internal/domain/e1rm/calculator.go`

## Interface

```go
type Calculator struct {
    rpeChart *rpechart.RPEChart
}

func NewCalculator(chart *rpechart.RPEChart) *Calculator

func (c *Calculator) Calculate(weight float64, reps int, rpe float64) (float64, error)
```

## Validation

- Weight must be > 0
- Reps must be 1-12 (RPE chart range)
- RPE must be 6.5-10.0 (valid RPE range)
- RPE chart lookup must succeed

## Rounding

- Use RoundWeight with 2.5 increment (standard for maxes)
- Round to nearest

## Tests

- Basic calculation: 315 × 5 @ RPE 8 → expected ~389
- Edge case: 1 rep @ RPE 10 → E1RM equals weight
- Invalid inputs: negative weight, out-of-range reps/RPE
- RPE chart lookup failure propagation

## Acceptance Criteria

- [ ] Calculator takes RPE chart as dependency
- [ ] Returns error for invalid inputs
- [ ] Returns error if RPE chart lookup fails
- [ ] Rounds to 2.5 lb increments
- [ ] Comprehensive unit tests pass

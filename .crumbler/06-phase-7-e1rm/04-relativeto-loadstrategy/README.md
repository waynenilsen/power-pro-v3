# RelativeTo LoadStrategy

Implement a LoadStrategy that calculates weight as a percentage of today's top set.

## Use Case

After finding a rep max or hitting an RPE-based top set, back-off sets are often prescribed as "X% of what you just did". This enables prescription chaining within a session.

## Location

Update existing stub: `internal/domain/loadstrategy/relativeto.go` (or create if missing)

## Interface

```go
type RelativeToLoadStrategy struct {
    ReferenceSetIndex int              // Which set to reference (0 = first/top set)
    Percentage        float64          // Percentage of reference weight (e.g., 85.0)
    RoundingIncrement float64          // Optional, defaults to 5.0
    RoundingDirection RoundingDirection // Optional, defaults to NEAREST

    sessionLookup SessionLookup `json:"-"`  // Injected dependency
}

type SessionLookup interface {
    GetLoggedSetByIndex(ctx context.Context, sessionID string, liftID string, setIndex int) (*LoggedSetResult, error)
}

type LoggedSetResult struct {
    Weight float64
    Reps   int
    RPE    *float64
}
```

## Calculation

```go
func (s *RelativeToLoadStrategy) CalculateLoad(ctx context.Context, params LoadCalculationParams) (float64, error) {
    // 1. Get sessionID from params.Context
    // 2. Look up the referenced set via sessionLookup
    // 3. Calculate: weight = referenceWeight * (percentage / 100)
    // 4. Round to increment
    // 5. Return
}
```

## Params Context

```go
params.Context = map[string]interface{}{
    "sessionID": "session-uuid",
    "liftID":    "squat",
}
```

## JSON Format

```json
{
  "type": "RELATIVE_TO",
  "referenceSetIndex": 0,
  "percentage": 85.0,
  "roundingIncrement": 5.0,
  "roundingDirection": "NEAREST"
}
```

## Edge Cases

- Reference set not yet logged → return error (must execute in order)
- Reference set index out of bounds → return error
- Session ID missing from context → return error

## Tests

- Basic: 400 lbs top set × 85% = 340 lbs
- Rounding: 405 lbs × 87% = 352.35 → 355 (round nearest 5)
- Error: Reference set not found
- Error: Missing session context

## Acceptance Criteria

- [ ] RelativeTo calculates percentage of reference set
- [ ] Proper error handling for missing/invalid references
- [ ] Rounding options work correctly
- [ ] Factory registration complete
- [ ] SessionLookup interface defined and injectable
- [ ] Comprehensive unit tests pass

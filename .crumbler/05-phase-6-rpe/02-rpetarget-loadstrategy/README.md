# RPETarget LoadStrategy

Implement the RPETarget load strategy that calculates weight based on target RPE and rep count.

## Background

The `RPE_TARGET` strategy type is already declared in `loadstrategy.go`. This task implements the concrete strategy.

RPETarget calculates load using: `Weight = 1RM × RPEChart(reps, RPE)`

## Implementation Tasks

### 1. Create `internal/domain/loadstrategy/rpetarget.go`

Follow the PercentOf implementation pattern in `percentof.go`:

```go
type RPETargetLoadStrategy struct {
    TargetReps        int               `json:"target_reps"`        // e.g., 5
    TargetRPE         float64           `json:"target_rpe"`         // e.g., 8.0
    RoundingIncrement float64           `json:"rounding_increment"` // e.g., 5.0
    RoundingDirection RoundingDirection `json:"rounding_direction"` // NEAREST, UP, DOWN

    // Injected dependencies
    maxLookup MaxLookup
    rpeChart  *rpechart.RPEChart
}
```

### 2. Implement LoadStrategy Interface

```go
func (s *RPETargetLoadStrategy) Type() LoadStrategyType { return TypeRPETarget }

func (s *RPETargetLoadStrategy) CalculateLoad(ctx context.Context, params LoadCalculationParams) (float64, error) {
    // 1. Validate params
    // 2. Get user's 1RM for the lift (always uses ONE_RM)
    // 3. Get RPE chart from params.LookupContext OR use injected rpeChart
    // 4. Look up percentage: rpeChart.GetPercentage(TargetReps, TargetRPE)
    // 5. Calculate: weight = 1RM * percentage
    // 6. Round using RoundWeight(weight, RoundingIncrement, RoundingDirection)
    // 7. Return rounded weight
}

func (s *RPETargetLoadStrategy) Validate() error {
    // TargetReps must be 1-12
    // TargetRPE must be 7.0-10.0
    // RoundingIncrement must be > 0
    // RoundingDirection must be valid
}
```

### 3. Dependency Injection

```go
func (s *RPETargetLoadStrategy) SetMaxLookup(ml MaxLookup) {
    s.maxLookup = ml
}

func (s *RPETargetLoadStrategy) SetRPEChart(chart *rpechart.RPEChart) {
    s.rpeChart = chart
}
```

### 4. JSON Serialization

```go
func (s *RPETargetLoadStrategy) MarshalJSON() ([]byte, error) {
    type Alias RPETargetLoadStrategy
    return json.Marshal(&struct {
        Type LoadStrategyType `json:"type"`
        *Alias
    }{
        Type:  TypeRPETarget,
        Alias: (*Alias)(s),
    })
}

func UnmarshalRPETarget(data json.RawMessage) (LoadStrategy, error) {
    var s RPETargetLoadStrategy
    if err := json.Unmarshal(data, &s); err != nil {
        return nil, err
    }
    if err := s.Validate(); err != nil {
        return nil, err
    }
    return &s, nil
}

func RegisterRPETarget(factory *StrategyFactory) {
    factory.Register(TypeRPETarget, UnmarshalRPETarget)
}
```

### 5. Register in Server

Update `internal/server/server.go`:
```go
loadstrategy.RegisterRPETarget(strategyFactory)
```

### 6. Unit Tests (`internal/domain/loadstrategy/rpetarget_test.go`)

- Test `CalculateLoad` with known RPE chart values
- Test rounding behavior (NEAREST, UP, DOWN)
- Test validation of invalid configs
- Test JSON serialization/deserialization
- Test factory registration

## Example Usage

```json
{
  "type": "RPE_TARGET",
  "target_reps": 5,
  "target_rpe": 8.0,
  "rounding_increment": 5.0,
  "rounding_direction": "NEAREST"
}
```

For a user with 1RM of 400 lbs:
- RPE chart lookup: 5 reps @ RPE 8 = 77% (0.77)
- Calculated weight: 400 × 0.77 = 308 lbs
- Rounded (NEAREST, 5): 310 lbs

## Acceptance Criteria

- [ ] `RPETargetLoadStrategy` implements `LoadStrategy` interface
- [ ] Uses RPEChart lookup for percentage calculation
- [ ] Weight rounding works correctly
- [ ] JSON serialization includes type discriminator
- [ ] Factory registration allows deserialization
- [ ] Registered in server initialization
- [ ] Comprehensive unit tests pass

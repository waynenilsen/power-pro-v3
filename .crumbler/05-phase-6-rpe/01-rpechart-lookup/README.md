# RPEChart Lookup

Implement the RPEChart lookup table following the existing WeeklyLookup/DailyLookup patterns.

## Background

The RPE Chart maps (reps, RPE) combinations to a percentage of 1RM. This is the core lookup table used by RTS (Reactive Training Systems) programs.

### Standard RPE Chart (from RTS Intermediate docs)

| RPE / Reps | 1     | 2     | 3     | 4     | 5     | 6     | 7     | 8     | 9     | 10    | 11    | 12    |
|------------|-------|-------|-------|-------|-------|-------|-------|-------|-------|-------|-------|-------|
| **7**      | 0.88  | 0.82  | 0.80  | 0.74  | 0.74  | 0.68  | 0.66  | 0.64  | 0.62  | 0.60  | 0.58  | 0.56  |
| **8**      | 0.91  | 0.88  | 0.82  | 0.80  | 0.77  | 0.71  | 0.68  | 0.66  | 0.64  | 0.62  | 0.60  | 0.58  |
| **9**      | 0.95  | 0.91  | 0.89  | 0.82  | 0.80  | 0.74  | 0.71  | 0.68  | 0.66  | 0.64  | 0.62  | 0.60  |
| **10**     | 1.00  | 0.95  | 0.92  | 0.88  | 0.82  | 0.80  | 0.74  | 0.71  | 0.68  | 0.66  | 0.64  | 0.62  |

## Implementation Tasks

### 1. Domain Entity (`internal/domain/rpechart/`)

Create `rpechart.go` following the patterns in `weeklylookup/` and `dailylookup/`:

```go
type RPEChartEntry struct {
    TargetReps int     // 1-12
    TargetRPE  float64 // 7.0, 7.5, 8.0, 8.5, 9.0, 9.5, 10.0
    Percentage float64 // 0.0-1.0 (e.g., 0.82 for 82%)
}

type RPEChart struct {
    Entries []RPEChartEntry
}
```

**Methods:**
- `GetPercentage(reps int, rpe float64) (float64, error)` - Look up percentage for reps+RPE
- `Validate() error` - Ensure entries are valid (RPE 7-10, reps 1-12, percentage 0-1)
- `NewRPEChart(entries []RPEChartEntry) (*RPEChart, error)` - Constructor with validation
- `NewDefaultRPEChart() *RPEChart` - Creates the standard RTS RPE chart

**Validation:**
- RPE must be between 7.0 and 10.0 (in 0.5 increments)
- Reps must be between 1 and 12
- Percentage must be between 0.0 and 1.0

### 2. Integrate with LookupContext

Update `internal/domain/loadstrategy/lookup_context.go`:

```go
type LookupContext struct {
    // ... existing fields ...
    RPEChart *rpechart.RPEChart
}

func (lc *LookupContext) HasRPEChart() bool
func (lc *LookupContext) GetRPEPercentage(reps int, rpe float64) (float64, error)
```

### 3. Unit Tests (`internal/domain/rpechart/rpechart_test.go`)

- Test `GetPercentage` for known values
- Test interpolation edge cases (if supported)
- Test validation for invalid entries
- Test `NewDefaultRPEChart()` returns correct values

## Acceptance Criteria

- [ ] RPEChart domain entity created at `internal/domain/rpechart/rpechart.go`
- [ ] Standard RTS chart embedded as `NewDefaultRPEChart()`
- [ ] `GetPercentage(reps, rpe)` returns correct percentage
- [ ] LookupContext extended with RPEChart support
- [ ] Comprehensive unit tests pass
- [ ] Half-RPE values (7.5, 8.5, 9.5) are supported via interpolation or explicit entries

# Extend SetScheme for Variable Counts

Add support for variable set counts in the SetScheme system. Currently all sets are determined at generation time - this task adds the foundation for schemes where set count depends on session performance.

## What to Implement

### 1. TerminationCondition Interface

Create `internal/domain/setscheme/termination.go`:

```go
// TerminationCondition determines when to stop generating more sets
type TerminationCondition interface {
    Type() string
    ShouldTerminate(ctx TerminationContext) bool
    Validate() error
}

type TerminationContext struct {
    SetNumber       int
    LastRPE         *float64  // nil if not tracked
    LastReps        int
    TotalReps       int
    TotalSets       int
    TargetReps      int       // what we wanted
}
```

### 2. Concrete Termination Conditions

Implement these conditions:
- `RPEThreshold` - Stop when RPE >= target (e.g., RPE 10)
- `RepFailure` - Stop when reps < target
- `MaxSets` - Safety limit (never exceed N sets)

### 3. Extend SetScheme Interface

Add optional method via interface embedding:

```go
// VariableSetScheme is optionally implemented by schemes with dynamic set counts
type VariableSetScheme interface {
    SetScheme
    IsVariableCount() bool
    GetTerminationCondition() TerminationCondition
    GenerateNextSet(ctx SetGenerationContext, history []GeneratedSet, termCtx TerminationContext) (*GeneratedSet, bool)
}
```

### 4. Update GeneratedSet

Add field to track provisional status:
```go
type GeneratedSet struct {
    // ... existing fields
    IsProvisional bool // true for variable schemes until logged
}
```

## Key Files to Modify

- `internal/domain/setscheme/setscheme.go` - Add VariableSetScheme interface
- `internal/domain/setscheme/termination.go` - New file for termination conditions
- `internal/domain/setscheme/factory.go` - Register termination condition types

## Acceptance Criteria

- [ ] TerminationCondition interface defined
- [ ] RPEThreshold termination implemented with tests
- [ ] RepFailure termination implemented with tests
- [ ] MaxSets termination implemented with tests
- [ ] VariableSetScheme interface defined (existing schemes unchanged)
- [ ] GeneratedSet has IsProvisional field
- [ ] All existing tests still pass
- [ ] Unit tests for all new code

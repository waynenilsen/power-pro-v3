# AMRAP SetScheme Implementation

Implement the AMRAP (As Many Reps As Possible) set scheme following the existing patterns in `internal/domain/setscheme/`.

## What to Implement

### 1. Create `internal/domain/setscheme/amrap.go`

Follow the pattern from `fixed.go` and `ramp.go`:

```go
type AMRAPSetScheme struct {
    // Sets is the number of AMRAP sets to generate (usually 1)
    Sets int `json:"sets"`
    // MinReps is the minimum expected reps (for display/logging purposes)
    MinReps int `json:"minReps"`
}
```

**Key behaviors:**
- `Type()` returns `TypeAMRAP` (already defined in `setscheme.go` line 22)
- `GenerateSets()` creates sets with `TargetReps` = `MinReps` and `IsWorkSet` = `true`
- `Validate()` ensures `Sets >= 1` and `MinReps >= 1`
- Include `MarshalJSON`, `UnmarshalAMRAPSetScheme`, and `RegisterAMRAPScheme`

### 2. Create `internal/domain/setscheme/amrap_test.go`

Unit tests covering:
- Validation (positive sets/reps, zero/negative rejection)
- `GenerateSets()` output correctness
- JSON serialization round-trip with type discriminator

### 3. Register in `internal/server/server.go`

Add `setscheme.RegisterAMRAPScheme(schemeFactory)` alongside existing scheme registrations.

## Design Notes

- AMRAP sets work exactly like Fixed sets for *generation* purposes - the difference is semantic (user performs max reps, not target reps)
- The actual "AMRAP" behavior (logging reps performed) happens at workout logging time, handled by LoggedSet (next crumb)
- `MinReps` is the minimum to "succeed" (e.g., Wendler's "5+" means MinReps=5)

## Files to Create/Modify
- `internal/domain/setscheme/amrap.go` (new)
- `internal/domain/setscheme/amrap_test.go` (new)
- `internal/server/server.go` (register scheme)

## Verification
- `go test ./internal/domain/setscheme/...` passes
- `go build ./...` succeeds

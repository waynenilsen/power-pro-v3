# RepRange SetScheme

Implement a flexible rep range set scheme that allows defining sets with minimum and maximum rep targets.

## Implementation

Create `/home/wayne/git/power-pro-v3/internal/domain/setscheme/rep_range.go`:

1. **Add Type Constant** - Add `TypeRepRange SetSchemeType = "REP_RANGE"` to `ValidSchemeTypes` map in `setscheme.go`

2. **Struct Definition**:
   ```go
   type RepRangeSetScheme struct {
       Sets    int `json:"sets"`
       MinReps int `json:"minReps"`
       MaxReps int `json:"maxReps"`
   }
   ```

3. **Methods**:
   - `Type()` - Returns TypeRepRange
   - `GenerateSets(baseWeight float64, ctx SetGenerationContext)` - Generates sets with MinReps as target
   - `Validate()` - Checks Sets >= 1, MinReps >= 1, MaxReps >= MinReps
   - `MarshalJSON()` - Includes type discriminator
   - `UnmarshalRepRangeSetScheme()` - Factory deserialization
   - `RegisterRepRangeScheme()` - Factory registration

4. **Write Unit Tests** in `rep_range_test.go`:
   - Type returns correct constant
   - Validate accepts valid schemes
   - Validate rejects invalid schemes (sets < 1, minReps < 1, maxReps < minReps)
   - GenerateSets produces correct output
   - JSON marshaling/unmarshaling roundtrip

## Acceptance Criteria

- RepRangeSetScheme exists with Sets, MinReps, MaxReps fields
- GenerateSets returns Sets number of sets with MinReps as target
- All tests pass

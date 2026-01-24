# Core Rotation Infrastructure

Build the foundational domain entities and infrastructure needed for rotation schedules.

## Tasks

### 1. Create RotationLookup Domain Entity
**File**: `internal/domain/rotationlookup/rotationlookup.go`

Create a new domain entity similar to WeeklyLookup that maps rotation positions to lift identifiers:

```go
type RotationLookup struct {
    ID        string
    Name      string
    Entries   []RotationLookupEntry
    ProgramID *string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type RotationLookupEntry struct {
    Position       int     // 0-based position in rotation
    LiftIdentifier string  // e.g., "deadlift", "squat", "bench"
    Description    string  // e.g., "Deadlift Focus - High Intensity AMRAP"
}
```

Include:
- `Create()` and `Update()` functions with validation
- `GetByPosition(position int) *RotationLookupEntry` method
- `Validate() *validation.Result` method
- Unit tests in `rotationlookup_test.go`

### 2. Extend LookupContext
**File**: `internal/domain/loadstrategy/lookup_context.go`

Add rotation fields to LookupContext:
- `RotationPosition int`
- `RotationLookup *rotationlookup.RotationLookup`

Add methods:
- `HasRotationLookup() bool`
- `GetRotationEntry() *RotationLookupEntry`
- `IsLiftInRotationFocus(liftIdentifier string) bool`

### 3. Extend UserProgramState
**File**: `internal/domain/userprogramstate/userprogramstate.go`

Add rotation state tracking:
- `RotationPosition int` - Current position in rotation (0-based)
- `CyclesSinceStart int` - Number of cycles completed

Add methods:
- `AdvanceRotation(rotationLength int)` - Advances position and wraps around

### 4. Database Migration
Create migration to add rotation fields to user_program_state table:
- `rotation_position INTEGER DEFAULT 0`
- `cycles_since_start INTEGER DEFAULT 0`

## Acceptance Criteria

- [ ] RotationLookup entity with full validation
- [ ] LookupContext extended with rotation fields and methods
- [ ] UserProgramState extended with rotation tracking
- [ ] Database migration applied successfully
- [ ] All unit tests pass
- [ ] Follows existing domain patterns (validation, factory, etc.)

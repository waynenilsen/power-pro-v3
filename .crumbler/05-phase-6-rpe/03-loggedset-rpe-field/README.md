# LoggedSet RPE Field

Add RPE field to LoggedSet for recording actual perceived exertion on performed sets.

## Background

LoggedSet records actual workout performance. Adding an RPE field allows users to log how hard each set felt, enabling:
- Fatigue tracking over time
- Auto-regulation based on actual vs target RPE
- Data for future e1RM calculations (Phase 7)

## Implementation Tasks

### 1. Database Migration

Create `migrations/NNNN_add_rpe_to_logged_sets.sql`:

```sql
-- Add RPE field to logged_sets table
ALTER TABLE logged_sets ADD COLUMN rpe REAL;
```

**Notes:**
- RPE is nullable (sets logged before this feature won't have RPE)
- Using REAL for decimal RPE values (7.5, 8.5, etc.)
- No NOT NULL constraint for backwards compatibility

### 2. Update Domain Entity

Update `internal/domain/loggedset/loggedset.go`:

```go
type LoggedSet struct {
    // ... existing fields ...

    // RPE is the rate of perceived exertion (7.0-10.0).
    // Optional - nil means RPE was not recorded.
    RPE *float64
}
```

Add validation:
```go
func (ls *LoggedSet) Validate() error {
    // ... existing validation ...

    if ls.RPE != nil {
        if *ls.RPE < 5.0 || *ls.RPE > 10.0 {
            return fmt.Errorf("RPE must be between 5.0 and 10.0")
        }
    }
    return nil
}
```

### 3. Update SQL Queries

Update `internal/db/queries/logged_sets.sql` to include RPE in:

- `CreateLoggedSet` - Add `rpe` parameter
- `GetLoggedSet` - Include `rpe` in SELECT
- `ListLoggedSetsBySession` - Include `rpe` in SELECT
- `UpdateLoggedSet` - Add `rpe` to UPDATE (if exists)

### 4. Regenerate sqlc

Run `sqlc generate` to regenerate:
- `internal/db/logged_sets.sql.go`
- `internal/db/models.go`

### 5. Update Repository

Update `internal/repository/loggedset_repository.go`:

- Map new `RPE` field in `toLoggedSet()` conversion
- Include `RPE` in `CreateLoggedSetParams`
- Handle nullable properly (sql.NullFloat64)

### 6. Update API Handlers

Update logged set API handlers to accept and return RPE:

- `POST /sessions/{id}/sets` - Accept optional `rpe` in request body
- `GET /sessions/{id}/sets` - Return `rpe` in response (null if not recorded)

### 7. Update Tests

- Add tests for creating logged sets with RPE
- Add tests for RPE validation
- Update existing tests to handle new field

## Acceptance Criteria

- [ ] Migration adds `rpe` column to `logged_sets` table
- [ ] LoggedSet domain entity has `RPE *float64` field
- [ ] RPE validation: 5.0-10.0 range (or nil)
- [ ] CRUD operations support RPE field
- [ ] API endpoints accept/return RPE
- [ ] Existing tests pass (backwards compatible)
- [ ] New tests for RPE functionality pass

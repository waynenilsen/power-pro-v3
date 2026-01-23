# LoggedSet Entity and API

Implement the LoggedSet entity to record actual performance (reps performed) for sets, enabling AMRAP-driven progressions.

## What to Implement

### 1. Database Migration

Create `migrations/00015_create_logged_sets_table.sql`:

```sql
-- +goose Up
CREATE TABLE logged_sets (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id),
    session_id TEXT NOT NULL,              -- Groups sets from same workout
    prescription_id TEXT NOT NULL,         -- Links back to the prescription
    lift_id TEXT NOT NULL REFERENCES lifts(id),
    set_number INT NOT NULL,
    weight REAL NOT NULL,
    target_reps INT NOT NULL,              -- What was prescribed
    reps_performed INT NOT NULL,           -- What was actually done
    is_amrap BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE(session_id, prescription_id, set_number)
);

CREATE INDEX idx_logged_sets_user ON logged_sets(user_id);
CREATE INDEX idx_logged_sets_session ON logged_sets(session_id);
CREATE INDEX idx_logged_sets_lift ON logged_sets(lift_id);

-- +goose Down
DROP TABLE IF EXISTS logged_sets;
```

### 2. SQLC Queries

Add to `internal/db/queries/logged_sets.sql`:
- `CreateLoggedSet` - Insert a logged set
- `GetLoggedSet` - Get by ID
- `ListLoggedSetsBySession` - Get all sets for a session
- `ListLoggedSetsByUser` - Get all sets for a user (with pagination)
- `GetLatestAMRAPForLift` - Get most recent AMRAP set for a lift/user

### 3. Domain Model

Create `internal/domain/loggedset/loggedset.go`:
- `LoggedSet` struct with validation
- Constructor `NewLoggedSet(...)`

### 4. Repository

Create `internal/repository/logged_set_repository.go`:
- Implement persistence using sqlc-generated code

### 5. API Handlers

Create `internal/api/logged_set_handler.go`:

**Endpoints:**
- `POST /sessions/{sessionId}/sets` - Log a set (array of sets)
- `GET /sessions/{sessionId}/sets` - Get all sets for a session
- `GET /users/{userId}/logged-sets` - Get user's logged sets (paginated)

### 6. Tests

- Unit tests for domain validation
- Integration tests for API endpoints

## Design Notes

- `session_id` is a UUID generated client-side (or by the API on first set logged)
- A session represents a single workout day
- `is_amrap` marks whether this was an AMRAP set (important for progression logic)
- The `prescription_id` links to what was prescribed, enabling comparison

## Files to Create/Modify

- `migrations/00015_create_logged_sets_table.sql` (new)
- `internal/db/queries/logged_sets.sql` (new)
- `internal/domain/loggedset/loggedset.go` (new)
- `internal/domain/loggedset/loggedset_test.go` (new)
- `internal/repository/logged_set_repository.go` (new)
- `internal/api/logged_set_handler.go` (new)
- `internal/api/logged_set_handler_test.go` (new)
- `internal/server/server.go` (wire up routes)

## Verification

- Run `goose up` to apply migration
- `sqlc generate` succeeds
- `go test ./...` passes
- API endpoints work via manual curl testing

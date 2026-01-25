# Phase 1: Schema & Domain Models

## Overview

Add state machine status fields to the database schema and create domain packages for state management.

## Tasks

### 1. Migration: Add status fields to user_program_states

Create migration `00023_add_state_machine_status_fields.sql`:

```sql
ALTER TABLE user_program_states
    ADD COLUMN enrollment_status TEXT NOT NULL DEFAULT 'ACTIVE'
    CHECK (enrollment_status IN ('ACTIVE', 'BETWEEN_CYCLES', 'QUIT'));

ALTER TABLE user_program_states
    ADD COLUMN cycle_status TEXT NOT NULL DEFAULT 'PENDING'
    CHECK (cycle_status IN ('PENDING', 'IN_PROGRESS', 'COMPLETED'));

ALTER TABLE user_program_states
    ADD COLUMN week_status TEXT NOT NULL DEFAULT 'PENDING'
    CHECK (week_status IN ('PENDING', 'IN_PROGRESS', 'COMPLETED'));
```

### 2. Migration: Create workout_sessions table

Create migration `00024_create_workout_sessions_table.sql`:

```sql
CREATE TABLE workout_sessions (
    id TEXT PRIMARY KEY,
    user_program_state_id TEXT NOT NULL,
    week_number INTEGER NOT NULL,
    day_index INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'IN_PROGRESS'
        CHECK (status IN ('IN_PROGRESS', 'COMPLETED', 'ABANDONED')),
    started_at TEXT NOT NULL DEFAULT (datetime('now')),
    finished_at TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (user_program_state_id) REFERENCES user_program_states(id) ON DELETE CASCADE
);

CREATE INDEX idx_workout_sessions_state_week_day
    ON workout_sessions(user_program_state_id, week_number, day_index);

CREATE INDEX idx_workout_sessions_status
    ON workout_sessions(user_program_state_id, status);
```

### 3. Update sqlc queries

- Add new columns to user_program_states queries
- Add CRUD queries for workout_sessions table
- Regenerate sqlc code

### 4. Domain: Update UserProgramState entity

Update `internal/domain/userprogramstate/userprogramstate.go`:
- Add `EnrollmentStatus`, `CycleStatus`, `WeekStatus` fields
- Add status type constants
- Add validation for status fields
- Update `EnrollUser` to set initial statuses

### 5. Domain: Create WorkoutSession entity

Create `internal/domain/workoutsession/` package:
- `workoutsession.go` - Entity with fields matching schema
- `validation.go` - Validation rules for state transitions

### 6. Domain: Create State Machine interfaces

Create `internal/domain/statemachine/` package:
- `statemachine.go` - Generic state machine interface
- `enrollment.go` - Enrollment state machine (ACTIVE, BETWEEN_CYCLES, QUIT)
- `cycle.go` - Cycle state machine (PENDING, IN_PROGRESS, COMPLETED)
- `week.go` - Week state machine (PENDING, IN_PROGRESS, COMPLETED)
- `workout.go` - Workout state machine (IN_PROGRESS, COMPLETED, ABANDONED)

## Acceptance Criteria

- [ ] New migrations run successfully
- [ ] sqlc generates updated code
- [ ] UserProgramState has status fields
- [ ] WorkoutSession domain package exists with validation
- [ ] State machine package exists with all state machines
- [ ] All existing tests pass
- [ ] New unit tests for status validation

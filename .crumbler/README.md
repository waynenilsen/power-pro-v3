# State Machine Implementation Plan

## Overview

PowerPro currently has implicit state transitions scattered throughout the codebase. User position is tracked via `user_program_states` fields (`current_week`, `current_day_index`, `current_cycle_iteration`) but there are no explicit state machines governing transitions. This leads to:

1. **Computed state on every read** - Performance overhead
2. **Implicit transitions** - Hard to reason about, test, and extend
3. **No workout session entity** - `logged_sets.session_id` is an orphan FK
4. **Poor UX** - Can't show user "you're between cycles" or "workout in progress"

This plan introduces explicit, persisted state machines at each level of the hierarchy.

---

## Scope & Responsibilities

### What This Plan Covers (Backend API Team)

- State machine design and implementation
- Database schema changes
- API endpoints with proper state in responses
- Event system for progression triggers
- **E2E tests that document the API contract** (critical for frontend team)

### What This Plan Does NOT Cover (Frontend Team Responsibility)

- UI/UX for state transitions (showing "between cycles" screen, etc.)
- Client-side state management
- Optimistic updates or offline handling
- Push notifications or real-time updates

The frontend team consumes our API. Our E2E tests serve as **living documentation** showing exactly how the API is intended to be used, what state transitions look like, and what responses to expect at each step.

---

## Frontend Team Considerations

While we don't implement the frontend, we must design the API with their needs in mind:

### Clear State at Every Level

Every enrollment response includes all state levels so frontend can render appropriate UI:

```json
{
  "enrollment_status": "ACTIVE",      // Top-level: show program or "pick next" screen
  "cycle_status": "IN_PROGRESS",      // Show cycle progress indicator
  "week_status": "IN_PROGRESS",       // Show week progress
  "current_workout_session": {...}    // Show "resume workout" or "start workout"
}
```

### Actionable Next Steps

API responses should make it obvious what actions are available:

```json
{
  "enrollment_status": "BETWEEN_CYCLES",
  "available_actions": [
    {"action": "start_next_cycle", "endpoint": "POST /enrollment/next-cycle"},
    {"action": "change_program", "endpoint": "DELETE /enrollment"},
    {"action": "quit", "endpoint": "DELETE /enrollment"}
  ]
}
```

### Predictable State Transitions

Frontend needs to know:
- What state comes next after an action
- Which transitions are automatic (no user action needed)
- Which transitions require user decision

This is documented in E2E tests - frontend devs should **read the E2E tests** as their primary API documentation.

### Error States Are Explicit

When an action isn't allowed (e.g., starting workout when one is IN_PROGRESS):

```json
{
  "error": "workout_already_in_progress",
  "message": "Complete or abandon current workout before starting a new one",
  "current_workout_session_id": "..."
}
```

---

## State Machine Hierarchy (Outermost to Innermost)

```
┌─────────────────────────────────────────────────────────────────┐
│ Level 1: ENROLLMENT STATE (user ↔ program relationship)        │
│   States: ACTIVE | BETWEEN_CYCLES | QUIT                        │
│   Persisted: user_program_states.enrollment_status              │
├─────────────────────────────────────────────────────────────────┤
│ Level 2: CYCLE STATE (current cycle instance)                   │
│   States: PENDING | IN_PROGRESS | COMPLETED                     │
│   Persisted: user_program_states.cycle_status                   │
├─────────────────────────────────────────────────────────────────┤
│ Level 3: WEEK STATE (current week within cycle)                 │
│   States: PENDING | IN_PROGRESS | COMPLETED                     │
│   Persisted: user_program_states.week_status                    │
├─────────────────────────────────────────────────────────────────┤
│ Level 4: WORKOUT SESSION STATE (individual workout)             │
│   States: IN_PROGRESS | COMPLETED | ABANDONED                   │
│   Persisted: workout_sessions table (NEW)                       │
├─────────────────────────────────────────────────────────────────┤
│ Level 5: SET STATE (individual set within workout)              │
│   States: (implicit) PENDING (not logged) | LOGGED (exists)     │
│   Persisted: logged_sets table (existing)                       │
└─────────────────────────────────────────────────────────────────┘
```

---

## Level 1: Enrollment State Machine

### States

| State | Description |
|-------|-------------|
| `ACTIVE` | User is actively working through the program |
| `BETWEEN_CYCLES` | User completed a cycle, awaiting decision (continue same program or switch) |
| `QUIT` | User explicitly quit the program |

### Transitions

```
                         enroll()
            ┌───────────────────────────────────┐
            │                                   ▼
        ┌───────┐                         ┌──────────┐
        │ (none)│                         │  ACTIVE  │
        └───────┘                         └────┬─────┘
                                               │
                              ┌────────────────┼────────────────┐
                              │                │                │
                         quit()         cycle_completes()   quit()
                              │                │                │
                              ▼                ▼                │
                        ┌──────────┐    ┌─────────────┐         │
                        │   QUIT   │    │  BETWEEN_   │         │
                        │          │    │   CYCLES    │─────────┘
                        └──────────┘    └──────┬──────┘
                                               │
                              ┌────────────────┼────────────────┐
                              │                │                │
                      start_new_cycle()  change_program()   quit()
                              │                │                │
                              ▼                ▼                ▼
                        ┌──────────┐      ┌───────┐       ┌──────────┐
                        │  ACTIVE  │      │ (none)│       │   QUIT   │
                        │(iter++)  │      │       │       │          │
                        └──────────┘      └───────┘       └──────────┘
```

### API Triggers

| Trigger | Endpoint | Auto? |
|---------|----------|-------|
| `enroll()` | `POST /users/{id}/enrollment` | No |
| `quit()` | `DELETE /users/{id}/enrollment` | No |
| `cycle_completes()` | (internal) when last week completes | Yes |
| `start_new_cycle()` | `POST /users/{id}/enrollment/next-cycle` | No |
| `change_program()` | `DELETE` then `POST /enrollment` | No |

### Events Emitted

- `ENROLLED` - User enrolled in program
- `CYCLE_BOUNDARY_REACHED` - User completed cycle, awaiting decision
- `QUIT` - User quit program

---

## Level 2: Cycle State Machine

### States

| State | Description |
|-------|-------------|
| `PENDING` | Cycle exists but no workouts started yet |
| `IN_PROGRESS` | At least one workout has been started |
| `COMPLETED` | All weeks in cycle completed |

### Transitions

```
    first_workout_starts()           last_week_completes()
┌───────────┐              ┌─────────────┐              ┌───────────┐
│           │─────────────►│             │─────────────►│           │
│  PENDING  │              │ IN_PROGRESS │              │ COMPLETED │
│           │              │             │              │           │
└───────────┘              └─────────────┘              └───────────┘
```

### Behavior

- When user starts new cycle: `cycle_status = PENDING`, `current_cycle_iteration++`
- When first workout of cycle starts: `cycle_status = IN_PROGRESS`
- When last week completes: `cycle_status = COMPLETED`, triggers `enrollment_status = BETWEEN_CYCLES`

### Events Emitted

- `CYCLE_STARTED` - First workout of cycle began
- `CYCLE_COMPLETED` - All weeks done (triggers progression evaluation for AFTER_CYCLE)

---

## Level 3: Week State Machine

### States

| State | Description |
|-------|-------------|
| `PENDING` | Week is current but no workouts started |
| `IN_PROGRESS` | At least one workout started this week |
| `COMPLETED` | All days in week completed or week manually advanced |

### Transitions

```
    first_workout_starts()           all_days_complete() OR advance_week()
┌───────────┐              ┌─────────────┐              ┌───────────┐
│           │─────────────►│             │─────────────►│           │
│  PENDING  │              │ IN_PROGRESS │              │ COMPLETED │
│           │              │             │              │           │
└───────────┘              └─────────────┘              └───────────┘
```

### Behavior

- When advancing to new week: `week_status = PENDING`, `current_week++`
- When first workout of week starts: `week_status = IN_PROGRESS`
- When all days done OR manual advance: `week_status = COMPLETED`
- If `current_week > cycle.length_weeks`: triggers cycle completion

### Events Emitted

- `WEEK_STARTED` - First workout of week began
- `WEEK_COMPLETED` - Week finished (triggers progression evaluation for AFTER_WEEK)

---

## Level 4: Workout Session State Machine

### States

| State | Description |
|-------|-------------|
| `IN_PROGRESS` | Workout started, user logging sets |
| `COMPLETED` | User finished workout normally |
| `ABANDONED` | Workout started but never finished (timeout or explicit) |

### Transitions

```
         start_workout()
              │
              ▼
        ┌─────────────┐
        │             │
        │ IN_PROGRESS │
        │             │
        └──────┬──────┘
               │
    ┌──────────┼──────────┐
    │          │          │
finish()    abandon()   timeout()
    │          │          │
    ▼          ▼          ▼
┌─────────┐ ┌─────────┐
│COMPLETED│ │ABANDONED│
└─────────┘ └─────────┘
```

### New Table: `workout_sessions`

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

-- Allow multiple sessions per day (retries, partial workouts)
CREATE INDEX idx_workout_sessions_state_week_day
    ON workout_sessions(user_program_state_id, week_number, day_index);

-- Find in-progress sessions
CREATE INDEX idx_workout_sessions_status
    ON workout_sessions(user_program_state_id, status);
```

### Relationship to logged_sets

`logged_sets.session_id` will reference `workout_sessions.id`. We may add an FK constraint or keep it soft (for backwards compat with existing data).

### Events Emitted

- `WORKOUT_STARTED` - Workout began
- `WORKOUT_COMPLETED` - Workout finished (triggers progression evaluation for AFTER_SESSION)
- `WORKOUT_ABANDONED` - Workout abandoned

---

## Level 5: Set State (Existing)

Sets don't have an explicit state machine - their state is implicit:
- **PENDING**: No `logged_sets` row exists for this prescription/set_number in current session
- **LOGGED**: Row exists in `logged_sets`

### Events Emitted

- `SET_LOGGED` - Set was recorded (triggers progression evaluation for AFTER_SET, ON_FAILURE)

---

## Schema Changes Summary

### Migration: Add status fields to user_program_states

```sql
-- Add enrollment status
ALTER TABLE user_program_states
    ADD COLUMN enrollment_status TEXT NOT NULL DEFAULT 'ACTIVE'
    CHECK (enrollment_status IN ('ACTIVE', 'BETWEEN_CYCLES', 'QUIT'));

-- Add cycle status
ALTER TABLE user_program_states
    ADD COLUMN cycle_status TEXT NOT NULL DEFAULT 'PENDING'
    CHECK (cycle_status IN ('PENDING', 'IN_PROGRESS', 'COMPLETED'));

-- Add week status
ALTER TABLE user_program_states
    ADD COLUMN week_status TEXT NOT NULL DEFAULT 'PENDING'
    CHECK (week_status IN ('PENDING', 'IN_PROGRESS', 'COMPLETED'));
```

### Migration: Create workout_sessions table

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

---

## Event System Design

### Event Types

```go
type EventType string

const (
    EventEnrolled             EventType = "ENROLLED"
    EventCycleBoundaryReached EventType = "CYCLE_BOUNDARY_REACHED"
    EventQuit                 EventType = "QUIT"
    EventCycleStarted         EventType = "CYCLE_STARTED"
    EventCycleCompleted       EventType = "CYCLE_COMPLETED"
    EventWeekStarted          EventType = "WEEK_STARTED"
    EventWeekCompleted        EventType = "WEEK_COMPLETED"
    EventWorkoutStarted       EventType = "WORKOUT_STARTED"
    EventWorkoutCompleted     EventType = "WORKOUT_COMPLETED"
    EventWorkoutAbandoned     EventType = "WORKOUT_ABANDONED"
    EventSetLogged            EventType = "SET_LOGGED"
)
```

### Event Structure

```go
type StateEvent struct {
    Type       EventType
    UserID     string
    ProgramID  string
    Timestamp  time.Time
    Payload    map[string]interface{}
}
```

### Progression Trigger Mapping

| Progression TriggerType | Subscribes to Event |
|------------------------|---------------------|
| `AFTER_SET` | `SET_LOGGED` |
| `AFTER_SESSION` | `WORKOUT_COMPLETED` |
| `AFTER_WEEK` | `WEEK_COMPLETED` |
| `AFTER_CYCLE` | `CYCLE_COMPLETED` |
| `ON_FAILURE` | `SET_LOGGED` (filtered by failure) |

---

## State Transition Rules (Detailed)

### Auto-Transitions

These happen automatically as side effects:

1. **First workout of cycle starts**
   - If `cycle_status == PENDING`: set `cycle_status = IN_PROGRESS`
   - Emit `CYCLE_STARTED`

2. **First workout of week starts**
   - If `week_status == PENDING`: set `week_status = IN_PROGRESS`
   - Emit `WEEK_STARTED`

3. **All days in week completed**
   - Set `week_status = COMPLETED`
   - Emit `WEEK_COMPLETED`
   - If `current_week == cycle.length_weeks`:
     - Set `cycle_status = COMPLETED`
     - Emit `CYCLE_COMPLETED`
     - Set `enrollment_status = BETWEEN_CYCLES`
     - Emit `CYCLE_BOUNDARY_REACHED`
   - Else:
     - Set `current_week++`
     - Set `week_status = PENDING`

4. **Workout completed**
   - Set `workout_session.status = COMPLETED`
   - Set `workout_session.finished_at = now()`
   - Emit `WORKOUT_COMPLETED`
   - Check if all days in week done (trigger week completion if so)

### Manual Transitions

These require explicit API calls:

1. **Start workout** - `POST /workouts/start`
2. **Finish workout** - `POST /workouts/{id}/finish`
3. **Abandon workout** - `POST /workouts/{id}/abandon`
4. **Advance week manually** - `POST /users/{id}/enrollment/advance-week`
5. **Start new cycle** - `POST /users/{id}/enrollment/next-cycle`
6. **Quit program** - `DELETE /users/{id}/enrollment`

---

## API Endpoints (New/Modified)

### Workout Session Endpoints

```
POST   /workouts/start                    # Start new workout, returns session_id
GET    /workouts/{id}                     # Get workout session details
POST   /workouts/{id}/finish              # Mark workout complete
POST   /workouts/{id}/abandon             # Mark workout abandoned
GET    /users/{id}/workouts               # List user's workout history
GET    /users/{id}/workouts/current       # Get current in-progress workout (if any)
```

### Enrollment State Endpoints

```
GET    /users/{id}/enrollment             # Get enrollment with all state info
POST   /users/{id}/enrollment/next-cycle  # Start new cycle (when BETWEEN_CYCLES)
POST   /users/{id}/enrollment/advance-week # Manual week advance
```

### Response Shape Changes

```json
// GET /users/{id}/enrollment
{
  "id": "...",
  "user_id": "...",
  "program_id": "...",
  "enrollment_status": "ACTIVE",
  "cycle_status": "IN_PROGRESS",
  "week_status": "IN_PROGRESS",
  "current_week": 2,
  "current_cycle_iteration": 1,
  "current_day_index": 1,
  "current_workout_session": {
    "id": "...",
    "status": "IN_PROGRESS",
    "started_at": "2024-01-15T10:00:00Z"
  }
}
```

---

## Domain Package Structure

```
internal/domain/
├── statemachine/
│   ├── statemachine.go      # Generic state machine interface
│   ├── enrollment.go        # Enrollment state machine
│   ├── cycle.go             # Cycle state machine
│   ├── week.go              # Week state machine
│   └── workout.go           # Workout state machine
├── event/
│   ├── event.go             # Event types and structures
│   ├── bus.go               # In-memory event bus
│   └── handlers.go          # Event handler registration
└── workoutsession/
    ├── workoutsession.go    # WorkoutSession domain entity
    └── validation.go        # Validation rules
```

---

## E2E Testing Strategy

### Philosophy: E2E Tests as Frontend Documentation

The program-specific E2E tests in `internal/api/e2e/` are the **primary documentation** for how the frontend team should integrate with our API. They must:

1. **Show the complete happy path** - Every state transition the user experiences
2. **Verify state at every step** - Frontend knows exactly what to expect
3. **Demonstrate the intended API usage** - Not just "it works" but "this is how you use it"
4. **Be readable as a narrative** - A frontend dev should be able to read the test and understand the flow

### Current E2E Test Gaps

Looking at existing tests like `wendler_531_test.go`:

```go
// CURRENT: Implicit state, manual progression triggers
advanceUserState(t, ts, userID)  // What states changed? Unknown.
advanceUserState(t, ts, userID)  // Did week complete? Did cycle complete? Unknown.

// Manually force progression (not realistic)
triggerBody := ManualTriggerRequest{Force: true}  // Frontend can't do this
```

**Problems:**
- No `enrollment_status`, `cycle_status`, `week_status` in responses
- `advanceUserState()` hides all state transitions
- Progressions manually triggered, not event-driven
- No workout session lifecycle (start → log sets → finish)
- No "between cycles" state handling shown

### Required E2E Test Updates

Every program E2E test must be updated to show the **full state machine flow**:

```go
// AFTER: Explicit state verification at every step

// 1. ENROLL - verify initial state
enrollResp := userPost("/users/{id}/enrollment", enrollBody)
assertState(t, enrollResp, State{
    EnrollmentStatus: "ACTIVE",
    CycleStatus:      "PENDING",    // Not started yet
    WeekStatus:       "PENDING",
    CurrentWeek:      1,
    CurrentCycleIteration: 1,
})

// 2. START WORKOUT - state transitions happen
startResp := userPost("/workouts/start", startBody)
assertState(t, startResp, State{
    CycleStatus: "IN_PROGRESS",     // Auto-transitioned!
    WeekStatus:  "IN_PROGRESS",     // Auto-transitioned!
})
sessionID := startResp.Data.SessionID

// 3. LOG SETS - within session
logResp := userPost("/sessions/{sessionID}/sets", setsBody)
// verify sets logged

// 4. FINISH WORKOUT - state may transition
finishResp := userPost("/workouts/{sessionID}/finish", nil)
assertState(t, finishResp, State{
    // If more days in week:
    WeekStatus: "IN_PROGRESS",
    // OR if last day of week:
    WeekStatus: "COMPLETED",  // Auto-transitioned!
})

// ... continue through full cycle ...

// 5. CYCLE COMPLETES - progression fires automatically
finishLastWorkoutResp := userPost("/workouts/{sessionID}/finish", nil)
assertState(t, finishLastWorkoutResp, State{
    EnrollmentStatus: "BETWEEN_CYCLES",  // Auto-transitioned!
    CycleStatus:      "COMPLETED",
    WeekStatus:       "COMPLETED",
})

// Verify progression was applied automatically (not forced!)
maxResp := userGet("/users/{id}/lift-maxes/{liftId}/current")
assert(t, maxResp.Data.Value == originalValue + expectedDelta)

// 6. START NEW CYCLE - user decision required
nextCycleResp := userPost("/users/{id}/enrollment/next-cycle", nil)
assertState(t, nextCycleResp, State{
    EnrollmentStatus: "ACTIVE",
    CycleStatus:      "PENDING",
    WeekStatus:       "PENDING",
    CurrentCycleIteration: 2,  // Incremented!
})
```

### Program-Specific E2E Tests to Update

All 22 program E2E tests need updating:

| Test File | Key State Transitions to Show |
|-----------|-------------------------------|
| `wendler_531_test.go` | 4-week cycle, CYCLE_COMPLETED, AFTER_CYCLE progression |
| `starting_strength_test.go` | SESSION progression (every workout), failure handling |
| `greyskull_lp_test.go` | AMRAP-driven progression, AFTER_SET trigger |
| `texas_method_test.go` | AFTER_WEEK progression on intensity day |
| `gzclp_t1_test.go` | Stage progression, failure → stage advance |
| `gzclp_t2_test.go` | T2 stage progression |
| `nsuns_531_lp_test.go` | High frequency, multiple sessions per week |
| `nuckols_frequency_test.go` | 3x/week per lift, AMRAP-based |
| `sheiko_beginner_test.go` | Percentage-based, no AMRAP |
| `building_the_monolith_test.go` | 6-week cycle, specific deload |
| `inverted_juggernaut_test.go` | Wave progression with AMRAP |
| `jacked_and_tan_test.go` | Block periodization |
| `reddit_ppl_test.go` | 6-day rotation |
| `calgary_barbell_8_test.go` | Peaking program, meet_date handling |
| `calgary_barbell_16_test.go` | Longer peaking, phase transitions |
| `rts_intermediate_test.go` | RPE-based, fatigue management |
| `nsuns_cap3_test.go` | CAP3 periodization |
| `gzcl_compendium_test.go` | VDIP approach |
| `bill_starr_test.go` | Classic 5x5 progression |
| `nuckols_beginner_test.go` | Beginner frequency |
| `sheiko_intermediate_test.go` | Intermediate periodization |
| `phase5_rotation_e2e_test.go` | Rotation-based programs |

### New E2E Tests to Add

Beyond updating existing tests, add new tests specifically for state machine behavior:

#### 1. `state_machine_enrollment_test.go`

```go
func TestEnrollmentStateTransitions(t *testing.T) {
    // Test: NONE → ACTIVE (enroll)
    // Test: ACTIVE → BETWEEN_CYCLES (cycle completes)
    // Test: BETWEEN_CYCLES → ACTIVE (start new cycle)
    // Test: ACTIVE → QUIT (quit)
    // Test: BETWEEN_CYCLES → QUIT (quit while deciding)
    // Test: Can't start workout when BETWEEN_CYCLES
    // Test: Can't start new cycle when ACTIVE
}
```

#### 2. `state_machine_workout_session_test.go`

```go
func TestWorkoutSessionLifecycle(t *testing.T) {
    // Test: Start workout creates IN_PROGRESS session
    // Test: Can't start second workout while one IN_PROGRESS
    // Test: Finish workout transitions to COMPLETED
    // Test: Abandon workout transitions to ABANDONED
    // Test: Can start new workout after COMPLETED
    // Test: Can start new workout after ABANDONED
    // Test: Logging sets requires active session
    // Test: Can't log sets to COMPLETED session
}
```

#### 3. `state_machine_week_cycle_test.go`

```go
func TestWeekAndCycleTransitions(t *testing.T) {
    // Test: First workout of week: PENDING → IN_PROGRESS
    // Test: Last workout of week: → COMPLETED, auto-advance
    // Test: Last workout of cycle: triggers CYCLE_COMPLETED
    // Test: Week status resets to PENDING on new week
    // Test: Cycle status resets to PENDING on new cycle
}
```

#### 4. `state_machine_progression_events_test.go`

```go
func TestProgressionEventTriggers(t *testing.T) {
    // Test: AFTER_SET progression fires on set log
    // Test: AFTER_SESSION progression fires on workout finish
    // Test: AFTER_WEEK progression fires on week complete
    // Test: AFTER_CYCLE progression fires on cycle complete
    // Test: ON_FAILURE progression fires on failed set
    // Test: Progression does NOT fire when Force=false and conditions not met
}
```

### E2E Test Helpers to Add/Update

```go
// assertEnrollmentState verifies all state fields in enrollment response
func assertEnrollmentState(t *testing.T, resp *http.Response, expected EnrollmentState) {
    t.Helper()
    var envelope EnrollmentResponse
    json.NewDecoder(resp.Body).Decode(&envelope)

    if envelope.Data.EnrollmentStatus != expected.EnrollmentStatus {
        t.Errorf("enrollment_status: expected %s, got %s",
            expected.EnrollmentStatus, envelope.Data.EnrollmentStatus)
    }
    if envelope.Data.CycleStatus != expected.CycleStatus {
        t.Errorf("cycle_status: expected %s, got %s",
            expected.CycleStatus, envelope.Data.CycleStatus)
    }
    if envelope.Data.WeekStatus != expected.WeekStatus {
        t.Errorf("week_status: expected %s, got %s",
            expected.WeekStatus, envelope.Data.WeekStatus)
    }
    // ... etc
}

// startWorkoutAndVerify starts a workout and verifies state transitions
func startWorkoutAndVerify(t *testing.T, ts *TestServer, userID string,
    expectedCycleStatus, expectedWeekStatus string) string {
    t.Helper()
    // POST /workouts/start
    // Verify response state
    // Return session ID
}

// finishWorkoutAndVerify finishes a workout and verifies state
func finishWorkoutAndVerify(t *testing.T, ts *TestServer, sessionID string,
    expectedEnrollmentStatus, expectedCycleStatus, expectedWeekStatus string) {
    t.Helper()
    // POST /workouts/{id}/finish
    // Verify response state
    // Check if progressions fired (if applicable)
}

// completeFullCycle runs through an entire cycle verifying states
func completeFullCycle(t *testing.T, ts *TestServer, userID string,
    program ProgramConfig) CycleCompletionResult {
    t.Helper()
    // For each week in cycle:
    //   For each day in week:
    //     startWorkout, logSets, finishWorkout
    //     Verify state transitions
    // Verify BETWEEN_CYCLES at end
    // Return progression results
}
```

### Deprecate/Remove

- **Remove `advanceUserState()` usage** - Replace with explicit workout start/finish
- **Remove `Force: true` in progression triggers** - Test event-driven behavior
- **Remove any direct DB manipulation** - Everything through API

---

## Implementation Order (Tickets)

### Phase 1: Schema & Domain Models

1. **Migration: Add status fields to user_program_states**
   - Add `enrollment_status`, `cycle_status`, `week_status` columns
   - Default existing rows to appropriate states

2. **Migration: Create workout_sessions table**
   - Create table with indexes
   - Add sqlc queries

3. **Domain: WorkoutSession entity**
   - Create `internal/domain/workoutsession/` package
   - Validation, state transitions

4. **Domain: State machine interfaces**
   - Create `internal/domain/statemachine/` package
   - Generic interface + concrete implementations

### Phase 2: Event System

5. **Domain: Event types and bus**
   - Create `internal/domain/event/` package
   - In-memory event bus with subscriber pattern

6. **Integration: Wire events to progression system**
   - Progression service subscribes to relevant events
   - Replace manual trigger calls with event-driven

### Phase 3: API Layer

7. **Handler: Workout session endpoints**
   - Start, finish, abandon, list, get current

8. **Handler: Enrollment state endpoints**
   - Get enrollment (with state), next-cycle, advance-week

9. **Handler: Modify existing endpoints**
   - Logged sets endpoint triggers events
   - Enrollment endpoint uses state machine

### Phase 4: E2E Test Infrastructure

10. **E2E: Add state assertion helpers**
    - `assertEnrollmentState()`, `startWorkoutAndVerify()`, etc.
    - Response types with state fields

11. **E2E: Add state machine specific tests**
    - `state_machine_enrollment_test.go`
    - `state_machine_workout_session_test.go`
    - `state_machine_week_cycle_test.go`
    - `state_machine_progression_events_test.go`

### Phase 5: Update All Program E2E Tests

12. **E2E: Update program tests batch 1 (LP programs)**
    - `starting_strength_test.go`
    - `greyskull_lp_test.go`
    - `bill_starr_test.go`
    - `gzclp_t1_test.go`
    - `gzclp_t2_test.go`

13. **E2E: Update program tests batch 2 (531 variants)**
    - `wendler_531_test.go`
    - `nsuns_531_lp_test.go`
    - `building_the_monolith_test.go`

14. **E2E: Update program tests batch 3 (periodized)**
    - `texas_method_test.go`
    - `inverted_juggernaut_test.go`
    - `jacked_and_tan_test.go`
    - `gzcl_compendium_test.go`

15. **E2E: Update program tests batch 4 (frequency/RPE)**
    - `nuckols_frequency_test.go`
    - `nuckols_beginner_test.go`
    - `rts_intermediate_test.go`
    - `nsuns_cap3_test.go`

16. **E2E: Update program tests batch 5 (peaking/other)**
    - `calgary_barbell_8_test.go`
    - `calgary_barbell_16_test.go`
    - `sheiko_beginner_test.go`
    - `sheiko_intermediate_test.go`
    - `reddit_ppl_test.go`
    - `phase5_rotation_e2e_test.go`

### Phase 6: Polish & Backfill

17. **Migration: Backfill existing data**
    - Set appropriate states for existing enrollments based on position
    - Create placeholder workout_sessions for existing logged_sets if needed

---

## Open Questions

1. **Workout session timeout**: How long before IN_PROGRESS → ABANDONED automatically?
   - Option A: Background job checks for stale sessions
   - Option B: Check on next workout start, mark old one abandoned
   - Option C: No auto-timeout, only explicit abandon

2. **Multiple workouts per day**: Allow or enforce one completed workout per day?
   - Current: `UNIQUE(session_id, prescription_id, set_number)` allows multiple sessions
   - Proposed: Allow multiple, track which is "canonical" via COMPLETED status

3. **Backwards compatibility**: Soft or hard FK from logged_sets to workout_sessions?
   - Existing logged_sets have session_ids that don't exist in workout_sessions
   - Option A: Soft reference (no FK constraint)
   - Option B: Create placeholder sessions for existing data
   - Option C: New column `workout_session_id` with FK, keep old `session_id`

---

## Performance Considerations

1. **State is persisted on write, not computed on read**
   - No joins or calculations needed to determine current state
   - Single column reads for status checks

2. **Event bus is in-memory**
   - No persistence overhead for events
   - Events are fire-and-forget (progression results are persisted separately)

3. **Indexes support common queries**
   - Find in-progress workout: `idx_workout_sessions_status`
   - List workouts for week: `idx_workout_sessions_state_week_day`

---

## Success Criteria

- [ ] All state transitions are explicit and logged
- [ ] No computed state on read paths
- [ ] Progression system is fully event-driven
- [ ] API responses include current state at each level
- [ ] E2E tests cover full enrollment lifecycle
- [ ] E2E tests demonstrate happy path for all 22 programs
- [ ] E2E tests serve as usable documentation for frontend team
- [ ] Existing functionality unchanged (backwards compatible)
- [ ] Frontend team can implement full workout flow using only E2E tests as reference

---

## Appendix: Frontend Integration Quick Reference

> **Note:** This section is a summary for the frontend team. Implementation is their responsibility.
> The E2E tests are the authoritative source - read them for exact request/response shapes.

### Typical User Flow (Happy Path)

```
1. USER ENROLLS IN PROGRAM
   POST /users/{id}/enrollment {programId: "..."}
   Response: enrollment_status=ACTIVE, cycle_status=PENDING, week_status=PENDING

   Frontend: Show "Ready to start" screen, "Begin Workout" button

2. USER STARTS FIRST WORKOUT
   POST /workouts/start {dayIndex: 0}  (or auto-detect current day)
   Response: session_id, cycle_status=IN_PROGRESS, week_status=IN_PROGRESS

   Frontend: Show workout screen with exercises/sets

3. USER LOGS SETS (repeat for each set)
   POST /sessions/{sessionId}/sets {sets: [...]}
   Response: logged set data

   Frontend: Update UI, show next set

4. USER FINISHES WORKOUT
   POST /workouts/{sessionId}/finish
   Response: updated enrollment state

   Frontend: Check response state:
   - If week_status=IN_PROGRESS: Show "See you next workout"
   - If week_status=COMPLETED: Show "Week complete!" celebration
   - If enrollment_status=BETWEEN_CYCLES: Show "Cycle complete!" + decision screen

5. USER CONTINUES (repeat steps 2-4 for each workout)
   ...

6. CYCLE COMPLETES (automatic when last workout finished)
   Response from finish: enrollment_status=BETWEEN_CYCLES

   Frontend: Show decision screen:
   - "Start another cycle of [Program]" → POST /enrollment/next-cycle
   - "Choose different program" → DELETE /enrollment, then new POST /enrollment
   - "Take a break" → (do nothing, stay in BETWEEN_CYCLES)

7. USER STARTS NEW CYCLE
   POST /users/{id}/enrollment/next-cycle
   Response: enrollment_status=ACTIVE, cycle_status=PENDING, cycle_iteration=2

   Frontend: Back to "Ready to start" screen
```

### State-Based UI Rendering

```
enrollment_status | What to show
------------------|------------------------------------------
(no enrollment)   | Program selection screen
ACTIVE            | Current workout or "start workout" button
BETWEEN_CYCLES    | "Cycle complete" decision screen
QUIT              | (enrollment deleted, show program selection)
```

```
cycle_status | week_status | What to show
-------------|-------------|------------------------------------------
PENDING      | PENDING     | "Ready to begin cycle" + start button
IN_PROGRESS  | PENDING     | "New week" indicator + start button
IN_PROGRESS  | IN_PROGRESS | Normal workout view
IN_PROGRESS  | COMPLETED   | "Week done" + next week preview
COMPLETED    | COMPLETED   | (enrollment_status will be BETWEEN_CYCLES)
```

```
current_workout_session | What to show
------------------------|------------------------------------------
null                    | "Start Workout" button
{status: IN_PROGRESS}   | "Resume Workout" button, show session
{status: COMPLETED}     | (null after completion)
{status: ABANDONED}     | "Workout abandoned" message, start new
```

### Error Handling

```
Action attempted              | Error response                    | Frontend action
------------------------------|-----------------------------------|------------------
Start workout while one open  | workout_already_in_progress       | Prompt to finish/abandon
Start new cycle while ACTIVE  | cannot_start_cycle_while_active   | Show current workout
Log sets without session      | no_active_session                 | Prompt to start workout
Finish non-existent session   | session_not_found                 | Refresh state
```

### Progression Visibility

Progressions happen automatically. Frontend can show results by:

1. **Compare lift maxes before/after cycle:**
   ```
   GET /users/{id}/lift-maxes/{liftId}/current
   ```

2. **Check progression logs:**
   ```
   GET /users/{id}/progression-logs?since={cycleStartDate}
   ```

3. **Response from workout finish may include:**
   ```json
   {
     "progressions_applied": [
       {"lift": "Squat", "previous": 315, "new": 325, "delta": 10}
     ]
   }
   ```

### Key Timestamps for Display

- `enrolled_at` - When user started this program
- `workout_session.started_at` - Workout duration tracking
- `workout_session.finished_at` - Completion time
- `progression_logs.applied_at` - When max increased

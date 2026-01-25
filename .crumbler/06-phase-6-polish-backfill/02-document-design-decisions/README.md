# Document Design Decisions

## Overview

Document the architectural decisions made during the state machine implementation for future reference.

## Tasks

### 1. Create design decisions document

Create `docs/design-decisions/state-machine.md` with the following sections:

#### Workout Session Timeout Decision
**Decision**: Option B - Check on next workout start, mark old one abandoned

**Rationale**:
- No background jobs needed - simpler infrastructure
- Abandoned sessions are rare edge cases
- Users naturally self-correct by starting new workouts
- Abandoned state properly recorded for analytics

**Implementation**:
- When starting a new workout via `POST /workouts/start`:
  - Check for existing IN_PROGRESS session
  - If found, return error with session ID (user must explicitly finish/abandon)
- Explicit abandon endpoint: `POST /workouts/{id}/abandon`

#### Multiple Workouts Per Day Decision
**Decision**: Allow multiple, COMPLETED status marks canonical

**Rationale**:
- Users may abandon workouts partway through (phone dies, gym emergency, etc.)
- Re-attempts are valid use cases
- COMPLETED status indicates the "real" workout
- Historical data preserved for analytics

**Implementation**:
- Multiple workout_sessions can exist for same week_number + day_index
- Only one IN_PROGRESS session allowed at a time (enforced at API level)
- COMPLETED workouts counted for progression triggers
- ABANDONED workouts tracked but don't trigger progressions

#### Backwards Compatibility Decision
**Decision**: Option A - Soft reference for logged_sets.session_id

**Rationale**:
- Simplest approach - no data migration needed
- Pre-migration logged_sets have client-generated session_ids
- No FK constraint means existing data remains valid
- New logged_sets can optionally reference workout_sessions.id

**Implementation**:
- logged_sets.session_id has no FK to workout_sessions
- Historical session_ids are preserved
- Frontend determines session_id format (legacy or new)
- API accepts both formats

### 2. Add state transition diagram

Include Mermaid diagram showing:
- Enrollment states: ACTIVE ↔ BETWEEN_CYCLES → QUIT
- Cycle states: PENDING → IN_PROGRESS → COMPLETED
- Week states: PENDING → IN_PROGRESS → COMPLETED
- Workout session states: IN_PROGRESS → COMPLETED | ABANDONED

### 3. Document event types

List all events emitted by the system:
- ENROLLED
- QUIT
- CYCLE_STARTED
- CYCLE_BOUNDARY_REACHED
- WEEK_STARTED (implicit via workout start)
- WEEK_COMPLETED
- WORKOUT_STARTED
- WORKOUT_COMPLETED
- WORKOUT_ABANDONED
- PROGRESSION_APPLIED

## Acceptance Criteria

- [ ] Design decisions document created
- [ ] Rationale clearly explained for each decision
- [ ] State transition diagram included
- [ ] Event types documented
- [ ] Document follows existing docs conventions

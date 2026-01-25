# Phase 6: Polish & Backfill

## Overview

Handle backwards compatibility for existing data and add any finishing touches.

## Tasks

### 1. Migration: Backfill existing data

Create migration `00025_backfill_state_machine_data.sql`:

For existing user_program_states:
- Set `enrollment_status = 'ACTIVE'` (already defaulted)
- Set `cycle_status = 'IN_PROGRESS'` for rows with `current_day_index IS NOT NULL`
- Set `week_status = 'IN_PROGRESS'` for rows with `current_day_index IS NOT NULL`

### 2. Create placeholder workout_sessions for existing logged_sets

Decide approach (from open questions):
- **Option A**: Soft reference (no FK constraint) - existing session_ids remain orphan
- **Option B**: Create placeholder sessions for existing data
- **Option C**: New column `workout_session_id` with FK, keep old `session_id`

Recommended: **Option A** - keep soft reference, document that pre-migration logged_sets have session_ids that don't map to workout_sessions table.

### 3. Address open questions

Document decisions for:

1. **Workout session timeout**:
   - Recommend: Option B - Check on next workout start, mark old one abandoned
   - Rationale: No background jobs needed, simple implementation

2. **Multiple workouts per day**:
   - Recommend: Allow multiple, COMPLETED status marks canonical
   - Rationale: Users may abandon workouts, re-attempts are valid

3. **Backwards compatibility**:
   - Recommend: Option A - Soft reference
   - Rationale: Simplest, no data migration needed

### 4. Performance verification

Verify:
- State reads are single column lookups (no joins)
- Event bus doesn't add noticeable latency
- Indexes support common query patterns

### 5. Documentation updates

Update relevant documentation:
- API documentation with new endpoints
- State transition diagrams
- Frontend integration guide (derived from E2E tests)

## Acceptance Criteria

- [ ] Existing data migrated appropriately
- [ ] Open questions resolved and documented
- [ ] Performance acceptable
- [ ] All tests pass (unit, integration, E2E)
- [ ] Documentation updated
- [ ] Success criteria from parent README all met:
  - [ ] All state transitions are explicit and logged
  - [ ] No computed state on read paths
  - [ ] Progression system is fully event-driven
  - [ ] API responses include current state at each level
  - [ ] E2E tests cover full enrollment lifecycle
  - [ ] E2E tests demonstrate happy path for all 22 programs
  - [ ] E2E tests serve as usable documentation for frontend team
  - [ ] Existing functionality unchanged (backwards compatible)

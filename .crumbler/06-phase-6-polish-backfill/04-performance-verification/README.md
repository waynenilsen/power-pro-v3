# Performance Verification

## Overview

Verify that the state machine implementation meets performance requirements.

## Tasks

### 1. Verify state reads are efficient

Check that state reads are single column lookups with no joins:

```sql
-- User program state read should be:
SELECT * FROM user_program_states WHERE user_id = ?
-- This is an indexed lookup on user_id (UNIQUE constraint)
```

Verify:
- `GetEnrollmentWithProgram` uses efficient JOIN (programs table only)
- Status fields are stored directly (no computed state)
- No N+1 queries in list operations

### 2. Verify event bus performance

Check that event bus doesn't add noticeable latency:
- Events are published asynchronously (`PublishAsync`)
- Handlers run in background goroutines
- No blocking on event delivery
- Handler errors don't affect API responses

### 3. Verify indexes support common queries

Check these query patterns have indexes:

```sql
-- Workout sessions by user_program_state and status
-- Covered by: idx_workout_sessions_status (user_program_state_id, status)

-- Workout sessions by user_program_state, week, day
-- Covered by: idx_workout_sessions_state_week_day

-- User program state by user_id
-- Covered by: UNIQUE constraint on user_id

-- logged_sets by session_id
-- Covered by: idx_logged_sets_session
```

### 4. Run basic performance test

Run the test suite with race detection and timing:

```bash
go test -race -v ./... 2>&1 | tee test-output.txt
# Check for race conditions
# Note any tests that take unusually long
```

### 5. Document findings

Create a brief performance notes section in the design decisions doc noting:
- Which operations are O(1) lookups
- Which operations require joins
- Any potential bottlenecks identified

## Acceptance Criteria

- [ ] State reads verified as efficient (no computed state on read)
- [ ] Event bus verified as async (no blocking)
- [ ] Index coverage verified for common queries
- [ ] Test suite passes with -race flag
- [ ] No race conditions detected
- [ ] Performance acceptable (tests complete in reasonable time)

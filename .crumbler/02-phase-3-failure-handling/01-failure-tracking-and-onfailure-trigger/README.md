# Failure Tracking and OnFailure Trigger

Implement the foundation for failure detection and the `OnFailure` trigger type.

## Context

The LoggedSet entity already tracks `TargetReps` vs `RepsPerformed` with helper methods:
- `ExceededTarget()` - returns true if reps_performed > target_reps
- `RepsDifference()` - returns reps_performed - target_reps (negative = failure)

This crumb adds:
1. **Failure counter tracking** - Track consecutive failures per (user, lift, progression)
2. **OnFailure trigger type** - A new trigger that fires when a set fails

## Implementation Requirements

### 1. Failure Counter Entity/Storage

Create a `FailureCounter` to track consecutive failures:
- `UserID` - The user
- `LiftID` - The lift being tracked
- `ProgressionID` - The progression rule this counter is for
- `ConsecutiveFailures` - Number of consecutive failures (resets on success)
- `LastFailureAt` - Timestamp of last failure

Storage options (choose one):
- New SQLite table `failure_counters` with migration
- Or extend existing `progression_logs` to track failure state

### 2. OnFailure Trigger Type

Add to `internal/domain/progression/trigger.go`:
- New `TriggerOnFailure` constant in `TriggerType`
- `FailureTriggerContext` struct containing:
  - `LoggedSet` - The set that failed
  - `ConsecutiveFailures` - Current failure count
  - `TargetReps` - What was prescribed
  - `RepsPerformed` - What was achieved

### 3. Failure Detection Service

Add to progression service or create new service:
- `CheckForFailure(loggedSet)` - Determines if a logged set is a failure
- `IncrementFailureCounter(userID, liftID, progressionID)` - Increments counter
- `ResetFailureCounter(userID, liftID, progressionID)` - Resets on success
- `GetFailureCount(userID, liftID, progressionID)` - Gets current count

### 4. Integration with LoggedSet Creation

When a LoggedSet is created:
1. Check if `RepsPerformed < TargetReps` (failure)
2. If failure:
   - Increment failure counter
   - Fire `OnFailure` trigger for applicable progressions
3. If success:
   - Reset failure counter

## Files to Create/Modify

- `internal/domain/progression/trigger.go` - Add TriggerOnFailure type and context
- `internal/domain/progression/failure_counter.go` - New failure counter entity (if new file)
- `internal/repository/failure_counter_repository.go` - Repository for counters
- `internal/service/failure_service.go` - Failure detection and counter management
- `migrations/00016_create_failure_counters_table.sql` - Database migration
- `internal/db/queries/failure_counter.sql` - SQL queries for sqlc

## Acceptance Criteria

- [ ] `TriggerOnFailure` type added to trigger system
- [ ] `FailureTriggerContext` properly carries failure information
- [ ] Failure counters persist in database
- [ ] Consecutive failures increment on each failure
- [ ] Counter resets to 0 on successful set
- [ ] OnFailure trigger fires with correct context
- [ ] Unit tests cover failure detection logic
- [ ] Integration tests verify counter persistence

## Dependencies

- Phase 2 (AMRAP) must be complete - it is
- LoggedSet entity with TargetReps/RepsPerformed - exists

## Notes

- Keep the trigger system polymorphic like existing triggers
- The failure counter is per (user, lift, progression) to allow different progressions to have different failure thresholds
- Consider idempotency - logging the same set twice shouldn't double-count failures

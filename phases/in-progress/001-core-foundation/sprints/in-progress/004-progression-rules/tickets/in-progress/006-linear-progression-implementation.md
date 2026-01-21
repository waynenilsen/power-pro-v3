# 006: LinearProgression Strategy Implementation

## ERD Reference
Implements: REQ-LINEAR-001, REQ-LINEAR-002, REQ-LINEAR-003, REQ-LINEAR-004
Related to: REQ-PROG-002

## Description
Implement the LinearProgression strategy that adds a fixed increment at regular intervals. This is the most common progression type for beginner and intermediate programs.

## Context / Background
Linear progression is the workhorse of strength training. Starting Strength adds 5lb per session, Bill Starr adds 5lb per week. This strategy supports all interval types (per-session, per-week) with configurable increments. It's the foundation for programs targeting novice and intermediate lifters.

## Acceptance Criteria
- [ ] Implement `LinearProgression` struct with:
  - `ID` (uuid) - from base Progression entity
  - `Name` (string) - human-readable name
  - `Increment` (decimal) - weight to add on each application
  - `MaxType` (string) - which max to update (1RM or TRAINING_MAX)
  - `TriggerType` (TriggerType) - when to fire (AFTER_SESSION or AFTER_WEEK)
- [ ] Implement `Type()` method returning `LINEAR_PROGRESSION`
- [ ] Implement `Apply()` method:
  - Verify trigger event type matches configured TriggerType
  - For AFTER_SESSION: only apply if lift was in performed lifts
  - Fetch current LiftMax for (userId, liftId, maxType)
  - Create new LiftMax with value = current + increment
  - Return ProgressionResult with applied=true, delta=increment
  - Return applied=false if trigger type mismatch or lift not in session
- [ ] Support per-session linear progression (Starting Strength pattern):
  - TriggerType = AFTER_SESSION
  - Only fires for lifts performed in that session
  - Configurable increment per lift
- [ ] Support per-week linear progression (Bill Starr pattern):
  - TriggerType = AFTER_WEEK
  - Fires once per week for all configured lifts
  - Configurable increment per lift
- [ ] Unit tests covering:
  - Correct increment applied
  - Trigger type matching/rejection
  - Session-specific lift filtering
  - Integration with LiftMax service
- [ ] Performance: Apply completes in < 100ms (NFR-001)

## Technical Notes
- Implements Progression interface from ticket 004
- Uses LiftMax service from ERD-001 for current value lookup and update
- Trigger filtering happens in Apply() - return applied=false for mismatches
- Consider using increment override from ProgramProgression if present
- Transaction required for LiftMax update + ProgressionLog creation (ticket 008)

## Dependencies
- Blocks: 008 (Trigger integration uses this implementation)
- Blocked by: 001, 004, 005 (Schema, interface, trigger types)
- Related: ERD-001 (Uses LiftMax service), 003 (ProgramProgression configuration)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/004-progression-rules/erd.md
- Starting Strength: programs/starting-strength.md
- Bill Starr 5x5: programs/bill-starr-5x5.md

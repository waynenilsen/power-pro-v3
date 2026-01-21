# 007: CycleProgression Strategy Implementation

## ERD Reference
Implements: REQ-CYCLE-001, REQ-CYCLE-002, REQ-CYCLE-003
Related to: REQ-PROG-002

## Description
Implement the CycleProgression strategy that adds a fixed increment at cycle completion. This supports periodized programs like 5/3/1 and Greg Nuckols High Frequency.

## Context / Background
Cycle-based progression is used by periodized programs that need longer adaptation windows. 5/3/1 uses a 4-week cycle with +5lb upper/+10lb lower at cycle end. Greg Nuckols HF uses a 3-week cycle. This strategy fires only when the cycle completes (week wraps back to 1).

## Acceptance Criteria
- [ ] Implement `CycleProgression` struct with:
  - `ID` (uuid) - from base Progression entity
  - `Name` (string) - human-readable name
  - `Increment` (decimal) - weight to add on cycle completion
  - `MaxType` (string) - which max to update (1RM or TRAINING_MAX)
  - (TriggerType is implicitly AFTER_CYCLE)
- [ ] Implement `Type()` method returning `CYCLE_PROGRESSION`
- [ ] Implement `Apply()` method:
  - Verify trigger event type is AFTER_CYCLE
  - Fetch current LiftMax for (userId, liftId, maxType)
  - Create new LiftMax with value = current + increment
  - Return ProgressionResult with applied=true, delta=increment
  - Return applied=false if trigger type is not AFTER_CYCLE
- [ ] Support lift-specific increment overrides:
  - Check ProgramProgression.override_increment first
  - Fall back to CycleProgression.Increment if no override
  - Enables 5/3/1 pattern: +5lb upper, +10lb lower
- [ ] Works with any cycle length:
  - 3-week cycles (Greg Nuckols HF)
  - 4-week cycles (5/3/1)
  - N-week cycles (future programs)
- [ ] Unit tests covering:
  - Correct increment applied at cycle end
  - Lift-specific override increments
  - Trigger type rejection for non-AFTER_CYCLE
  - Integration with LiftMax service
- [ ] Performance: Apply completes in < 100ms (NFR-001)

## Technical Notes
- Implements Progression interface from ticket 004
- Implicit AFTER_CYCLE trigger - no need to store triggerType
- Uses LiftMax service from ERD-001 for current value lookup and update
- Check ProgramProgression for lift-specific override_increment
- Transaction required for LiftMax update + ProgressionLog creation (ticket 008)
- Cycle length detection comes from trigger context, not progression config

## Dependencies
- Blocks: 008 (Trigger integration uses this implementation)
- Blocked by: 001, 004, 005 (Schema, interface, trigger types)
- Related: ERD-001 (Uses LiftMax service), 003 (ProgramProgression for overrides)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/004-progression-rules/erd.md
- 5/3/1: programs/531.md
- Greg Nuckols HF: programs/greg-nuckols-hf.md

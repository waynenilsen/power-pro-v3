# 005: Trigger Type Definitions and Event Structure

## ERD Reference
Implements: REQ-TRIG-001
Related to: REQ-TRIG-002, REQ-TRIG-003, REQ-TRIG-004

## Description
Define the trigger type enumeration and event structures that control when progressions fire. This establishes the vocabulary for progression timing.

## Context / Background
Progressions fire at different times for different programs: Starting Strength adds weight after each session, Bill Starr adds weight after each week, and 5/3/1 adds weight after each cycle. The trigger system provides clean separation between "when" (triggers) and "what" (progression logic).

## Acceptance Criteria
- [ ] Define `TriggerType` enum/constants:
  - `AFTER_SESSION` - fires after completing a training day
  - `AFTER_WEEK` - fires when advancing from week N to week N+1
  - `AFTER_CYCLE` - fires when advancing from week N to week 1 (cycle wrap)
- [ ] Define `TriggerEvent` struct with:
  - `Type` (TriggerType) - which trigger type fired
  - `UserID` (uuid) - which user triggered the event
  - `Timestamp` (time.Time) - when the trigger occurred
  - `Context` (TriggerContext) - type-specific context
- [ ] Define `TriggerContext` interface and concrete types:
  - `SessionTriggerContext`: sessionId, daySlug, weekNumber, liftsPerformed []uuid
  - `WeekTriggerContext`: previousWeek, newWeek, cycleIteration
  - `CycleTriggerContext`: completedCycle, newCycle, totalWeeks
- [ ] Define JSON serialization for TriggerContext (for ProgressionLog storage)
- [ ] Test coverage for trigger type validation and context serialization

## Technical Notes
- TriggerType maps to LinearProgression.triggerType parameter
- CycleProgression implicitly uses AFTER_CYCLE
- Context structures enable progressions to make decisions based on what happened
- liftsPerformed in SessionTriggerContext allows per-lift progression (Starting Strength)
- Consider Go interface with concrete implementations for type safety

## Dependencies
- Blocks: 004, 006, 007, 008 (Interface, implementations, and trigger integration use these)
- Blocked by: None
- Related: ERD-003 (State advancement API generates these events)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/004-progression-rules/erd.md
- Schedule & Periodization ERD-003: triggers originate from state advancement

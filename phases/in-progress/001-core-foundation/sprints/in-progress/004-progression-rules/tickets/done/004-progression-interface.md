# 004: Progression Interface and Strategy Pattern

## ERD Reference
Implements: REQ-PROG-001, REQ-PROG-002
Related to: NFR-006, NFR-007

## Description
Define the common interface for all progression strategies and the core types used for progression application. This establishes the polymorphic contract that all progression implementations must follow.

## Context / Background
Per the ERD, Progression enables polymorphic progression handling. Different programs progress differently (linear per-session, linear per-week, cycle-based, etc.). This interface ensures all strategies can be used interchangeably while remaining extensible for future implementations (AMRAP, failure-based, RPE-based).

## Acceptance Criteria
- [ ] Define `Progression` interface with:
  - `Type() ProgressionType` method - returns the type discriminator
  - `Apply(ctx context.Context, params ProgressionContext) (ProgressionResult, error)` method signature
- [ ] Define `ProgressionType` enum/constants:
  - `LINEAR_PROGRESSION` - fixed increment at intervals
  - `CYCLE_PROGRESSION` - fixed increment at cycle completion
  - (Document future types: AMRAP_PROGRESSION, DELOAD_ON_FAILURE, etc.)
- [ ] Define `ProgressionContext` struct with:
  - `UserID` (uuid) - for LiftMax lookup
  - `LiftID` (uuid) - for LiftMax lookup
  - `MaxType` (string) - 1RM or TRAINING_MAX
  - `CurrentValue` (decimal) - current LiftMax value
  - `TriggerEvent` (TriggerEvent) - what caused this progression to evaluate
- [ ] Define `ProgressionResult` struct with:
  - `Applied` (bool) - whether progression was actually applied
  - `PreviousValue` (decimal) - value before progression
  - `NewValue` (decimal) - value after progression
  - `Delta` (decimal) - the increment applied
  - `LiftID` (uuid) - which lift was modified
  - `MaxType` (string) - which max type was modified
  - `AppliedAt` (timestamp) - when applied
- [ ] Define JSON serialization format for Progression polymorphic storage/retrieval
- [ ] Implement factory function: `NewProgression(progressionType string, params json.RawMessage) (Progression, error)`
- [ ] Test coverage for interface contracts and type discrimination

## Technical Notes
- Use discriminated union pattern: `{"type": "LINEAR_PROGRESSION", ...type-specific-fields}`
- Go implementation: interface + concrete types that implement it
- Factory function handles deserialization from database JSONB
- Consider using Go generics or type assertions for parameter handling
- Apply method should be idempotent when combined with ProgressionLog checks

## Dependencies
- Blocks: 005, 006, 008 (LinearProgression, CycleProgression, trigger integration implement this)
- Blocked by: 001 (Needs Progression schema for persistence)
- Related: 007 (Trigger types used in ProgressionContext)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/004-progression-rules/erd.md
- Strategy Pattern: NFR-006, NFR-007 in ERD

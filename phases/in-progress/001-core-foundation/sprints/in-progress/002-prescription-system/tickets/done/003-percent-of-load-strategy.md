# 003: PercentOf LoadStrategy Implementation

## ERD Reference
Implements: REQ-LOAD-002, REQ-LOAD-003, REQ-LOAD-006
Related to: REQ-LOAD-001

## Description
Implement the PercentOf load strategy, which calculates target weight as a percentage of a reference max (1RM or Training Max). This is the foundational load calculation method used by most powerlifting programs.

## Context / Background
The vast majority of programs prescribe load as a percentage of a reference max. For example, "5x5 at 85% of Training Max" or "3x3 at 90% of 1RM". This strategy fetches the user's current max for the specified lift and calculates the target weight.

## Acceptance Criteria
- [ ] Implement `PercentOfLoadStrategy` struct with fields:
  - `ReferenceType` (enum: ONE_RM, TRAINING_MAX) - required
  - `Percentage` (float64) - required, percentage value (e.g., 85 for 85%)
  - `RoundingIncrement` (float64, optional) - weight increment (default 5.0)
  - `RoundingDirection` (enum: NEAREST, DOWN, UP, default NEAREST)
- [ ] Implement `CalculateLoad` method:
  - Fetch current max for user/lift/reference type from LiftMax repository
  - Calculate: `max * (percentage / 100)`
  - Apply rounding (uses rounding system from ticket 004)
  - Return calculated weight
- [ ] Implement validation (REQ-LOAD-006):
  - Percentage must be > 0
  - Percentage > 100 allowed (for overload work)
  - ReferenceType must be valid enum value
  - RoundingIncrement must be > 0 if specified
- [ ] Return clear error if max not found for user/lift/type
- [ ] JSON serialization/deserialization for strategy storage
- [ ] Test coverage > 90% including:
  - Normal percentage calculations
  - Percentages > 100% (overload)
  - Missing max handling
  - Invalid parameter validation
  - Both 1RM and TM reference types

## Technical Notes
- Depends on LiftMax repository from sprint 001
- ReferenceType enum values: `ONE_RM = "1RM"`, `TRAINING_MAX = "TM"`
- Rounding logic delegated to rounding system (ticket 004) for reusability
- Consider caching max lookups during batch operations (handled in ticket 010)

## Dependencies
- Blocks: 008, 010 (Domain logic and resolution use this strategy)
- Blocked by: 002 (Interface must be defined first), 004 (Rounding system)
- Related: Sprint 001 (LiftMax entity for max lookup)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/002-prescription-system/erd.md
- ERD-001: LiftMax entity from sprint 001

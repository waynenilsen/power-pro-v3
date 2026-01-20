# 006: Fixed SetScheme Implementation

## ERD Reference
Implements: REQ-SET-002, REQ-SET-005
Related to: REQ-SET-001

## Description
Implement the Fixed set scheme, which generates a specified number of sets with the same weight and reps. This is the most common scheme in powerlifting (e.g., 5x5, 3x8, 4x6).

## Context / Background
Fixed sets are the workhorse of strength programming. Starting Strength uses 3x5, Bill Starr uses 5x5, BBB uses 5x10. The Fixed scheme is simple: N sets of M reps at the prescribed weight. All sets are work sets.

## Acceptance Criteria
- [ ] Implement `FixedSetScheme` struct with fields:
  - `Sets` (int) - number of sets to generate (required)
  - `Reps` (int) - repetitions per set (required)
- [ ] Implement `GenerateSets` method:
  - Generate `Sets` number of `GeneratedSet` structs
  - Each set has:
    - `SetNumber`: 1 through Sets
    - `Weight`: baseWeight (unchanged)
    - `TargetReps`: Reps value
    - `IsWorkSet`: true (all Fixed sets are work sets)
- [ ] Implement validation (REQ-SET-005):
  - `Sets` must be >= 1
  - `Reps` must be >= 1
  - Return clear error message on validation failure
- [ ] JSON serialization/deserialization for scheme storage
- [ ] Test coverage > 90% including:
  - Various set/rep combinations (5x5, 3x8, 1x5)
  - Validation error cases
  - Edge cases (1 set, 1 rep)

## Technical Notes
- All sets are identical, just incrementing SetNumber
- No rounding needed at this level (weight already rounded by LoadStrategy)
- Simple implementation but critical for correct output

## Dependencies
- Blocks: 008, 010 (Domain logic and resolution use this scheme)
- Blocked by: 005 (Interface must be defined first)
- Related: 007 (Ramp scheme is the other scheme type)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/002-prescription-system/erd.md
- Example programs using Fixed: Starting Strength (3x5), Bill Starr (5x5), BBB (5x10)

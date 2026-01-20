# 005: SetScheme Interface Definition

## ERD Reference
Implements: REQ-SET-001
Related to: NFR-006, NFR-007

## Description
Define the common interface for all set/rep schemes. This establishes the polymorphic contract that all SetScheme implementations must follow for generating concrete sets from a base weight.

## Context / Background
SetScheme defines the structure of sets and reps for an exercise. Different programs use different schemes: fixed sets (5x5), ramping warmups, AMRAP sets, etc. This interface ensures all schemes can be used interchangeably while remaining extensible.

## Acceptance Criteria
- [ ] Define `SetScheme` interface with:
  - `Type` field/discriminator (string enum: "FIXED", "RAMP" for now)
  - `GenerateSets(baseWeight float64, context SetGenerationContext) ([]GeneratedSet, error)` method signature
- [ ] Define `SetGenerationContext` struct with:
  - Additional context fields as needed (may be empty initially)
- [ ] Define `GeneratedSet` struct with:
  - `SetNumber` (int) - 1-indexed set number
  - `Weight` (float64) - target weight for this set
  - `TargetReps` (int) - target number of repetitions
  - `IsWorkSet` (bool) - true if this is a work set, false if warmup
- [ ] Define `SetSchemeType` enum/constants:
  - `FIXED` - same weight and reps for all sets
  - `RAMP` - progressive percentages across sets
  - (Others documented as future: AMRAP, TOP_BACKOFF, etc.)
- [ ] Define JSON serialization format for SetScheme polymorphic storage
- [ ] Interface supports strategy pattern for extensibility (NFR-006)
- [ ] New schemes can be added without modifying existing code (NFR-007)
- [ ] Test coverage for interface contracts and type discrimination

## Technical Notes
- Use discriminated union pattern: `{"type": "FIXED", ...scheme-specific-fields}`
- Go implementation: interface + concrete types that implement it
- JSON unmarshaling should use type field to determine concrete type
- Consider factory function: `NewSetScheme(schemeType string, params json.RawMessage) (SetScheme, error)`
- SetNumber is 1-indexed for user-facing display

## Dependencies
- Blocks: 006, 007 (Fixed and Ramp implementations implement this interface)
- Blocked by: None
- Related: 001 (Schema stores SetScheme as JSON)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/002-prescription-system/erd.md
- Strategy Pattern: NFR-006 in ERD

# 002: LoadStrategy Interface Definition

## ERD Reference
Implements: REQ-LOAD-001
Related to: NFR-006, NFR-007

## Description
Define the common interface for all load calculation strategies. This establishes the polymorphic contract that all LoadStrategy implementations must follow.

## Context / Background
Per the ERD, LoadStrategy enables polymorphic load calculation. Different programs calculate target weights differently (percentage of max, RPE-based, fixed weight, etc.). This interface ensures all strategies can be used interchangeably while remaining extensible for future implementations.

## Acceptance Criteria
- [ ] Define `LoadStrategy` interface with:
  - `Type` field/discriminator (string enum: "PERCENT_OF" for now, extensible)
  - `CalculateLoad(ctx context.Context, params LoadCalculationParams) (float64, error)` method signature
- [ ] Define `LoadCalculationParams` struct with:
  - `UserID` (uuid) - for max lookup
  - `LiftID` (uuid) - for max lookup
  - Additional context fields as needed
- [ ] Define `LoadStrategyType` enum/constants:
  - `PERCENT_OF` - percentage of max
  - (Others documented as future: RPE_TARGET, FIXED_WEIGHT, RELATIVE_TO)
- [ ] Define JSON serialization format for LoadStrategy polymorphic storage
- [ ] Interface supports strategy pattern for extensibility (NFR-006)
- [ ] New strategies can be added without modifying existing code (NFR-007)
- [ ] Test coverage for interface contracts and type discrimination

## Technical Notes
- Use discriminated union pattern: `{"type": "PERCENT_OF", ...strategy-specific-fields}`
- Go implementation: interface + concrete types that implement it
- JSON unmarshaling should use type field to determine concrete type
- Consider factory function: `NewLoadStrategy(strategyType string, params json.RawMessage) (LoadStrategy, error)`

## Dependencies
- Blocks: 003 (PercentOf implementation implements this interface)
- Blocked by: None
- Related: 001 (Schema stores LoadStrategy as JSON)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/002-prescription-system/erd.md
- Strategy Pattern: NFR-006 in ERD

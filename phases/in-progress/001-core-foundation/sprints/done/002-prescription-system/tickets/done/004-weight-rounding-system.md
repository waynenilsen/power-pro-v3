# 004: Weight Rounding System

## ERD Reference
Implements: REQ-LOAD-004, REQ-LOAD-005

## Description
Implement the weight rounding system that adjusts calculated weights to the nearest available plate increment. This is a reusable utility used by LoadStrategy implementations and potentially other weight-related calculations.

## Context / Background
Gyms have specific plate increments (2.5lb, 5lb, 2.5kg, 5kg, etc.). Calculated weights like 142.5 need to be rounded to practical values. The rounding system supports configurable increments and rounding directions to accommodate different user preferences and equipment availability.

## Acceptance Criteria
- [ ] Implement `RoundWeight` function with parameters:
  - `weight` (float64) - the weight to round
  - `increment` (float64) - the rounding increment (e.g., 2.5, 5.0)
  - `direction` (enum) - rounding direction
- [ ] Implement `RoundingDirection` enum:
  - `NEAREST` - round to nearest increment (default)
  - `DOWN` - always round down (conservative)
  - `UP` - always round up
- [ ] Rounding logic:
  - `NEAREST`: standard rounding (142.5 with 5.0 increment -> 145 if >= 142.5, else 140)
  - `DOWN`: floor to increment (147 with 5.0 -> 145)
  - `UP`: ceil to increment (143 with 5.0 -> 145)
- [ ] Handle edge cases:
  - Weight already at increment boundary
  - Very small increments (0.5, 1.0)
  - Zero increment (return original weight or error)
  - Negative weights (return error)
- [ ] Test coverage > 95% including:
  - All rounding directions
  - Various increment sizes (2.5, 5.0, 10.0, 1.0)
  - Boundary conditions
  - Edge cases (exact multiples, very small differences)

## Technical Notes
- Implement as pure functions for testability and reusability
- Use precise decimal arithmetic to avoid floating-point errors
- Consider creating `RoundedWeight` value object with increment/direction metadata
- Formula for NEAREST: `math.Round(weight / increment) * increment`
- Formula for DOWN: `math.Floor(weight / increment) * increment`
- Formula for UP: `math.Ceil(weight / increment) * increment`

## Dependencies
- Blocks: 003 (PercentOf uses rounding)
- Blocked by: None
- Related: None

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/002-prescription-system/erd.md
- REQ-LOAD-004, REQ-LOAD-005 for rounding requirements

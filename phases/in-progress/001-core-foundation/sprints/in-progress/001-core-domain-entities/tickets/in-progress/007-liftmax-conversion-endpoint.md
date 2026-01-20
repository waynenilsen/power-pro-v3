# 007: LiftMax Conversion Endpoint

## ERD Reference
Implements: REQ-MAX-009
Priority: Should

## Description
Implement API endpoint for converting between max types (1RM ↔ Training Max). This allows users to calculate what their Training Max should be based on their 1RM, or vice versa.

## Context / Background
Different programs use different reference types for percentages. 5/3/1 uses Training Max (typically 90% of 1RM), while others use 1RM directly. Users need an easy way to convert between types without manually calculating.

## Acceptance Criteria
- [ ] `GET /lift-maxes/{id}/convert` - Convert max to different type
  - Query params:
    - `to_type` (required): Target type ('ONE_RM' or 'TRAINING_MAX')
    - `percentage` (optional): Conversion percentage, default 90
  - Returns calculated value without persisting
  - Response includes:
    - `originalValue`: The source max value
    - `originalType`: The source max type
    - `convertedValue`: The calculated converted value
    - `convertedType`: The target type
    - `percentage`: The percentage used for conversion
  - Returns 400 if:
    - `to_type` is same as current type
    - `to_type` is invalid
    - `percentage` is not between 1 and 100
  - Returns 404 if max not found
  - Returns 403 if requesting user is not the owner or admin
- [ ] Conversion formulas:
  - 1RM to TM: `tm = oneRm * (percentage / 100)`
  - TM to 1RM: `oneRm = tm / (percentage / 100)`
- [ ] Converted value rounded to nearest 0.25 for consistency with value precision
- [ ] Response time < 50ms (p95) - calculation only, no database writes
- [ ] Integration tests cover:
  - Both conversion directions
  - Default percentage (90%)
  - Custom percentage values
  - Edge cases (percentage boundaries)
  - Authorization

## Technical Notes
- Use domain layer conversion logic from Ticket 005
- This is a read-only operation - no database writes
- Consider caching since the conversion is purely computational
- Round to nearest 0.25: `Math.round(value * 4) / 4`

## Dependencies
- Blocks: None
- Blocked by: 005 (Domain logic with conversion), 006 (Base CRUD API)
- Related: 008 (Current max query - similar authorization pattern)

## Resources / Links
- ERD: phases/todo/001-core-foundation/sprints/in-progress/001-core-domain-entities/erd.md
- REQ-MAX-009: Conversion requirement specification

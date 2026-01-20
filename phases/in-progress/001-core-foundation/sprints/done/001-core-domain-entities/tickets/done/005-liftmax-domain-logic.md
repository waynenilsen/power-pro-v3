# 005: LiftMax Domain Logic and Validation

## ERD Reference
Implements: REQ-MAX-001, REQ-MAX-004, REQ-MAX-005, REQ-MAX-008, REQ-MAX-009
Related to: NFR-008, NFR-009

## Description
Implement the domain logic layer for the LiftMax entity, including value validation, max type handling, TM validation warnings, and conversion logic between max types.

## Context / Background
Per NFR-008 and NFR-009, domain logic must be isolated from the API layer. LiftMax has critical business rules: the value must be positive, TM values should be validated against 1RM ranges, and conversion between max types must use consistent formulas.

## Acceptance Criteria
- [ ] Create LiftMax domain entity/model with all required fields
- [ ] Implement value validation:
  - Required, positive decimal
  - Precision to 0.25 (values like 315.25, 315.5, 315.75 are valid)
  - Returns clear error message on violation
- [ ] Implement max type validation:
  - Must be ONE_RM or TRAINING_MAX
  - Required field, no default
  - Returns clear error message on invalid type
- [ ] Implement TM validation warning (REQ-MAX-008):
  - When creating/updating a TRAINING_MAX, check if user has a ONE_RM for the same lift
  - If TM < 80% or > 95% of existing 1RM, generate warning (not error)
  - Warning should be logged and optionally returned in response
  - Operation proceeds despite warning
- [ ] Implement max conversion logic (REQ-MAX-009):
  - Convert between ONE_RM and TRAINING_MAX
  - Default TM percentage: 90%
  - Support custom percentage parameter
  - 1RM to TM: `tm = oneRm * (percentage / 100)`
  - TM to 1RM: `oneRm = tm / (percentage / 100)`
  - Returns calculated value without persisting
- [ ] Implement effective date validation:
  - Required timestamp
  - Defaults to current time if not provided
- [ ] Domain logic is testable in isolation (no database dependencies)
- [ ] Test coverage > 90% for all validation and conversion logic

## Technical Notes
- Implement as pure TypeScript classes/functions for testability
- Use value objects where appropriate (e.g., MaxValue, MaxType)
- Consider a MaxCalculator service for conversion logic
- Round conversion results to nearest 0.25 for consistency
- TM validation warning should use a Result type that can carry warnings

## Dependencies
- Blocks: 006, 007, 008 (CRUD API and endpoints use domain logic)
- Blocked by: 004 (Schema must exist first for TypeORM entity)
- Related: 002 (Similar pattern to Lift domain logic)

## Resources / Links
- ERD: phases/todo/001-core-foundation/sprints/in-progress/001-core-domain-entities/erd.md
- REQ-MAX-008: TM validation warning requirement
- REQ-MAX-009: Max conversion requirement

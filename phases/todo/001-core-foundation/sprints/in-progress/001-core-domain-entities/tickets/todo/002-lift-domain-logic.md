# 002: Lift Domain Logic and Validation

## ERD Reference
Implements: REQ-LIFT-001, REQ-LIFT-002, REQ-LIFT-003, REQ-LIFT-004, REQ-LIFT-005
Related to: NFR-008, NFR-009

## Description
Implement the domain logic layer for the Lift entity, including validation rules, slug generation, and circular reference detection. This layer enforces business rules independently from the API layer.

## Context / Background
Per NFR-008 and NFR-009, domain logic must be isolated from the API layer, and entity changes must be validated at the domain layer. This ensures consistent validation whether data comes from API requests, database seeds, or internal operations.

## Acceptance Criteria
- [ ] Create Lift domain entity/model with all required fields
- [ ] Implement name validation:
  - Required, non-empty
  - Maximum 100 characters
  - Returns clear error message on violation
- [ ] Implement slug validation:
  - Unique (validated at domain level, enforced at DB level)
  - Lowercase alphanumeric with hyphens only
  - Auto-generate from name if not provided
  - Returns clear error message on invalid format
- [ ] Implement slug generation from name:
  - Convert to lowercase
  - Replace spaces and special characters with hyphens
  - Remove consecutive hyphens
  - Trim leading/trailing hyphens
- [ ] Implement competition lift flag handling:
  - Default to false
  - Boolean validation
- [ ] Implement parent lift reference validation:
  - Optional field
  - Circular reference detection (lift cannot be its own ancestor)
  - Returns clear error message if circular reference detected
- [ ] Domain logic is testable in isolation (no database dependencies)
- [ ] Test coverage > 90% for all validation logic

## Technical Notes
- Implement as pure TypeScript classes/functions for testability
- Use value objects where appropriate (e.g., LiftSlug value object)
- Circular reference detection should use recursive ancestor checking
- Consider factory pattern for Lift creation with validation

## Dependencies
- Blocks: 003 (CRUD API uses domain logic for validation)
- Blocked by: 001 (Schema must exist first for TypeORM entity)
- Related: None

## Resources / Links
- ERD: phases/todo/001-core-foundation/sprints/in-progress/001-core-domain-entities/erd.md
- NFR-008: Domain logic isolation requirement
- NFR-009: Domain layer validation requirement

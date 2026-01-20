# 008: Prescription Domain Logic and Validation

## ERD Reference
Implements: REQ-PRSC-001, REQ-PRSC-002, REQ-PRSC-003, REQ-PRSC-004, REQ-PRSC-005, REQ-PRSC-006, REQ-PRSC-007
Related to: NFR-004, NFR-006, NFR-007

## Description
Implement the domain logic layer for the Prescription entity, including validation rules, entity construction, and the orchestration of LoadStrategy and SetScheme for prescription resolution.

## Context / Background
The Prescription entity is the fundamental unit of programming, linking a lift to load and set specifications. The domain layer enforces business rules, validates inputs, and coordinates the LoadStrategy and SetScheme components to resolve prescriptions into concrete workout instructions.

## Acceptance Criteria
- [ ] Create Prescription domain entity with all required fields:
  - `ID` (UUID)
  - `LiftID` (UUID, required)
  - `LoadStrategy` (LoadStrategy interface, required)
  - `SetScheme` (SetScheme interface, required)
  - `Order` (int, default 0)
  - `Notes` (string, optional, max 500 chars)
  - `RestSeconds` (int pointer, optional)
  - `CreatedAt`, `UpdatedAt` (timestamps)
- [ ] Implement validation for all fields:
  - LiftID required and must be valid UUID format
  - LoadStrategy required and valid (delegates to strategy validation)
  - SetScheme required and valid (delegates to scheme validation)
  - Order >= 0
  - Notes max 500 characters
  - RestSeconds >= 0 when provided
- [ ] Implement Prescription factory/constructor with validation
- [ ] Implement `Resolve(ctx, userID)` method:
  - Call LoadStrategy.CalculateLoad to get base weight
  - Call SetScheme.GenerateSets with base weight
  - Return ResolvedPrescription with sets, notes, rest
- [ ] Define `ResolvedPrescription` struct:
  - `PrescriptionID` (UUID)
  - `Lift` (Lift info: ID, Name, Slug)
  - `Sets` ([]GeneratedSet)
  - `Notes` (string, optional)
  - `RestSeconds` (int pointer, optional)
- [ ] Graceful error handling (NFR-004):
  - Clear error if max not found
  - Clear error if lift not found
  - Partial success possible in batch (handled in ticket 010)
- [ ] Domain logic testable in isolation
- [ ] Test coverage > 90%

## Technical Notes
- Domain entity should be independent of database representation
- Use value objects where appropriate
- Validation errors should be descriptive and actionable
- Resolution requires access to LiftMax and Lift repositories

## Dependencies
- Blocks: 009, 010 (CRUD and Resolution APIs use domain logic)
- Blocked by: 001 (Schema), 002, 003 (LoadStrategy), 005, 006, 007 (SetScheme)
- Related: Sprint 001 (Lift and LiftMax entities)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/002-prescription-system/erd.md
- NFR-004: Graceful failure handling

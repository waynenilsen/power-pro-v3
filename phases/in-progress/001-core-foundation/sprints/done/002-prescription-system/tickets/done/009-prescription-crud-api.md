# 009: Prescription CRUD API Endpoints

## ERD Reference
Implements: REQ-PRSC-011
Related to: NFR-003

## Description
Implement the REST API endpoints for prescription management: create, read, update, delete, and list prescriptions.

## Context / Background
Prescriptions need to be managed through the API for program configuration. This ticket implements standard CRUD operations following the patterns established in sprint 001 for Lift and LiftMax entities.

## Acceptance Criteria
- [ ] Implement `GET /prescriptions` endpoint:
  - List prescriptions with filtering support
  - Filter by lift_id (optional)
  - Pagination support (limit, offset)
  - Returns array of prescription objects
- [ ] Implement `GET /prescriptions/{id}` endpoint:
  - Get single prescription by ID
  - Returns 404 if not found
  - Returns full prescription with LoadStrategy and SetScheme details
- [ ] Implement `POST /prescriptions` endpoint:
  - Create new prescription
  - Request body: liftId, loadStrategy, setScheme, order?, notes?, restSeconds?
  - Validates all fields via domain layer
  - Returns created prescription with 201 status
- [ ] Implement `PUT /prescriptions/{id}` endpoint:
  - Update existing prescription
  - Returns 404 if not found
  - Validates all fields via domain layer
  - Returns updated prescription
- [ ] Implement `DELETE /prescriptions/{id}` endpoint:
  - Delete prescription
  - Returns 404 if not found
  - Returns 204 on success
- [ ] All endpoints meet NFR-003 performance requirement (< 100ms p95)
- [ ] Follow API patterns from sprint 001 (error format, response structure)
- [ ] Test coverage > 90% including:
  - Happy path for all endpoints
  - Validation errors (400)
  - Not found errors (404)
  - Invalid JSON (400)

## Technical Notes
- Use existing router/handler patterns from sprint 001
- LoadStrategy and SetScheme transmitted as JSON objects with type discriminator
- Validation delegated to domain layer (ticket 008)
- Consider using DTOs for request/response transformation

## API Response Format
```json
{
  "id": "uuid",
  "liftId": "uuid",
  "loadStrategy": {
    "type": "PERCENT_OF",
    "referenceType": "TRAINING_MAX",
    "percentage": 85,
    "roundingIncrement": 5.0,
    "roundingDirection": "NEAREST"
  },
  "setScheme": {
    "type": "FIXED",
    "sets": 5,
    "reps": 5
  },
  "order": 1,
  "notes": "Focus on depth",
  "restSeconds": 180,
  "createdAt": "2024-01-15T10:00:00Z",
  "updatedAt": "2024-01-15T10:00:00Z"
}
```

## Dependencies
- Blocks: None
- Blocked by: 001 (Schema), 008 (Domain logic)
- Related: 010 (Resolution API)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/002-prescription-system/erd.md
- Sprint 001 patterns for CRUD API implementation

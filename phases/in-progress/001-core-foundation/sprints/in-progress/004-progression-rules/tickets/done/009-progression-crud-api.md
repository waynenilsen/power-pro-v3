# 009: Progression CRUD API Endpoints

## ERD Reference
Implements: REQ-API-001
Related to: REQ-PROG-001

## Description
Implement the REST API endpoints for managing Progression entities. This enables creating, reading, updating, and deleting progression rules through the API.

## Context / Background
Progression rules need to be managed through the API for program configuration. Administrators and program designers need to create progression configurations that define how users advance their training loads. This API provides full CRUD operations for the Progression entity.

## Acceptance Criteria
- [ ] Implement `GET /progressions` endpoint:
  - Returns list of all progressions
  - Supports filtering by `type` query parameter
  - Returns pagination metadata
  - Response includes id, name, type, parameters, createdAt, updatedAt
- [ ] Implement `GET /progressions/{id}` endpoint:
  - Returns single progression by ID
  - Returns 404 if not found
  - Response includes full progression details with parameters
- [ ] Implement `POST /progressions` endpoint:
  - Creates new progression
  - Validates request body: name required, type must be valid, parameters validated per type
  - Returns 201 with created progression
  - Returns 400 for validation errors
- [ ] Implement `PUT /progressions/{id}` endpoint:
  - Updates existing progression
  - Validates request body similar to POST
  - Returns 200 with updated progression
  - Returns 404 if not found
  - Returns 400 for validation errors
- [ ] Implement `DELETE /progressions/{id}` endpoint:
  - Deletes progression by ID
  - Returns 204 on success
  - Returns 404 if not found
  - Returns 409 if progression is referenced by ProgramProgressions
- [ ] Request/response validation:
  - Type discriminator must be valid enum value
  - Parameters must match type-specific schema
  - LinearProgression params: increment (required, positive), maxType (required), triggerType (required)
  - CycleProgression params: increment (required, positive), maxType (required)
- [ ] Error responses follow standard API error format
- [ ] Test coverage for all endpoints and error cases

## Technical Notes
- Use discriminated union pattern for request/response bodies
- Parameters field is JSONB - validate structure based on type
- Consider OpenAPI/Swagger documentation
- Follow TDD approach: endpoint structure first, hardcode responses, write failing tests, implement

## Dependencies
- Blocks: 010 (Program Progression Configuration uses progressions)
- Blocked by: 001 (Schema must exist)
- Related: 004 (Uses Progression interface/types)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/004-progression-rules/erd.md
- API Response Format: Section 5 of ERD

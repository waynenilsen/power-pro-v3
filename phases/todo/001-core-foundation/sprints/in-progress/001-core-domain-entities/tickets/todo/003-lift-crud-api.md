# 003: Lift CRUD API Endpoints

## ERD Reference
Implements: REQ-LIFT-006
Related to: NFR-001, NFR-003, NFR-007

## Description
Implement RESTful API endpoints for Lift entity CRUD operations. This provides the interface for creating, reading, updating, and deleting lift definitions.

## Context / Background
The Lift API enables program configuration and administration. All authenticated users can read lift data (NFR-007), while write operations require appropriate authorization. These endpoints must meet performance requirements (NFR-001) and support pagination (NFR-003).

## Acceptance Criteria
- [ ] `GET /lifts` - List all lifts
  - Supports pagination (default page size: 20)
  - Returns array of lift objects
  - Filterable by `is_competition_lift` query param
  - Sortable by `name` or `created_at`
- [ ] `GET /lifts/{id}` - Get single lift
  - Returns lift object by UUID
  - Returns 404 if lift not found
  - Supports lookup by slug as alternative: `GET /lifts/by-slug/{slug}`
- [ ] `POST /lifts` - Create lift
  - Accepts: name (required), slug (optional, auto-generated if omitted), isCompetitionLift (optional), parentLiftId (optional)
  - Returns 201 with created lift object
  - Returns 400 for validation errors
  - Returns 409 if slug already exists
- [ ] `PUT /lifts/{id}` - Update lift
  - Accepts partial updates
  - Returns 200 with updated lift object
  - Returns 404 if lift not found
  - Returns 400 for validation errors
  - Returns 409 if slug conflict
- [ ] `DELETE /lifts/{id}` - Delete lift
  - Returns 204 on success
  - Returns 404 if lift not found
  - Returns 409 if lift is referenced by LiftMax records (referential integrity per NFR-005)
- [ ] All endpoints respond in < 100ms (p95) per NFR-001
- [ ] Response format matches ERD specification (JSON with id, name, slug, isCompetitionLift, parentLiftId, createdAt, updatedAt)
- [ ] Authentication required for all endpoints
- [ ] Integration tests cover all endpoints and error cases

## Technical Notes
- Follow TDD approach: create endpoint structure first, write failing tests, then implement
- Use domain layer (Ticket 002) for validation
- Pagination should use cursor-based or offset pagination
- Consider ETag headers for caching
- Use proper HTTP status codes and error response format

## Dependencies
- Blocks: 004 (LiftMax schema references lifts table)
- Blocked by: 001 (Schema), 002 (Domain logic)
- Related: 009 (Authorization rules)

## Resources / Links
- ERD: phases/todo/001-core-foundation/sprints/in-progress/001-core-domain-entities/erd.md
- ERD Section 5: External Interfaces - Lift API Response format

# 006: LiftMax CRUD API Endpoints

## ERD Reference
Implements: REQ-MAX-010
Related to: NFR-001, NFR-003, NFR-004, NFR-005, NFR-006

## Description
Implement RESTful API endpoints for LiftMax entity CRUD operations. This provides the interface for users to manage their lift max values (1RM, Training Max).

## Context / Background
LiftMax API enables user onboarding (entering initial maxes) and progression tracking (updating maxes over time). Access is restricted to the owning user or admin (NFR-006). These endpoints must meet performance requirements and maintain referential integrity.

## Acceptance Criteria
- [ ] `GET /users/{userId}/lift-maxes` - List user's maxes
  - Supports pagination (default page size: 20)
  - Returns array of LiftMax objects for the specified user
  - Filterable by `lift_id` and `type` query params
  - Sortable by `effective_date` (default: descending)
  - Returns 403 if requesting user is not the owner or admin
- [ ] `GET /lift-maxes/{id}` - Get single max
  - Returns LiftMax object by UUID
  - Returns 404 if max not found
  - Returns 403 if requesting user is not the owner or admin
- [ ] `POST /users/{userId}/lift-maxes` - Create max
  - Accepts: liftId (required), type (required), value (required), effectiveDate (optional, defaults to now)
  - Returns 201 with created LiftMax object
  - Returns 400 for validation errors
  - Returns 409 if unique constraint violated (user, lift, type, effective_date)
  - Returns 403 if requesting user is not the owner or admin
  - Includes TM validation warning in response if applicable (REQ-MAX-008)
- [ ] `PUT /lift-maxes/{id}` - Update max
  - Accepts partial updates (value, effectiveDate)
  - Type and lift cannot be changed after creation
  - Returns 200 with updated LiftMax object
  - Returns 404 if max not found
  - Returns 400 for validation errors
  - Returns 403 if requesting user is not the owner or admin
  - Includes TM validation warning in response if applicable
- [ ] `DELETE /lift-maxes/{id}` - Delete max
  - Returns 204 on success
  - Returns 404 if max not found
  - Returns 403 if requesting user is not the owner or admin
- [ ] All endpoints respond in < 100ms (p95) per NFR-001
- [ ] Database transactions ensure consistency per NFR-004
- [ ] Referential integrity prevents orphaned maxes per NFR-005
- [ ] Response format matches ERD specification (JSON with id, userId, liftId, type, value, effectiveDate, createdAt, updatedAt)
- [ ] Authentication required for all endpoints
- [ ] Integration tests cover all endpoints, authorization, and error cases

## Technical Notes
- Follow TDD approach: create endpoint structure first, write failing tests, then implement
- Use domain layer (Ticket 005) for validation
- Include warnings in response body, not just logs
- Consider a response wrapper that can include warnings: `{ data: {...}, warnings: [...] }`
- Pagination should match Lift API pattern for consistency

## Dependencies
- Blocks: 007, 008 (Conversion and current max endpoints extend this)
- Blocked by: 004 (Schema), 005 (Domain logic)
- Related: 003 (Lift CRUD API - similar pattern), 009 (Authorization rules)

## Resources / Links
- ERD: phases/todo/001-core-foundation/sprints/in-progress/001-core-domain-entities/erd.md
- ERD Section 5: External Interfaces - LiftMax API Response format

# Ticket 005: API Documentation Update

## ERD Reference
Implements: REQ-TD2-009

## Description
Update API documentation to include all Phase 2 endpoints with accurate request/response schemas.

## Acceptance Criteria

### Auth Endpoints
- [ ] POST /auth/register documented
- [ ] POST /auth/login documented
- [ ] POST /auth/logout documented
- [ ] Authentication header format documented (Bearer token)

### Profile Endpoints
- [ ] GET /profile documented
- [ ] PUT /profile documented
- [ ] Authentication requirements documented

### Dashboard Endpoint
- [ ] GET /dashboard documented with response schema
- [ ] All aggregated data fields documented
- [ ] Authentication requirements documented

### General
- [ ] Error response format documented
- [ ] Common error codes documented (401, 403, 404, 422)
- [ ] Example requests/responses provided

## Technical Notes
- Update existing API documentation in docs/api/
- Use consistent documentation format
- Include curl examples

## Dependencies
- Blocked by: 004-phase2-test-coverage-review

# 005: API Documentation Update

## ERD Reference
Implements: REQ-TD2-009

## Description
Update API documentation to include all Phase 2 endpoints. Documentation should accurately reflect request/response schemas, authentication requirements, and error responses.

## Context / Background
Phase 2 added authentication, profile, and dashboard endpoints. API documentation must be updated to help frontend developers and API consumers integrate with these new endpoints.

## Acceptance Criteria

### Auth Endpoints
- [ ] POST /auth/register documented with request/response schema
- [ ] POST /auth/login documented with request/response schema
- [ ] POST /auth/logout documented with requirements
- [ ] Authentication header format documented (Bearer token)

### Profile Endpoints
- [ ] GET /profile documented with response schema
- [ ] PUT /profile documented with request/response schema
- [ ] Authentication requirements documented

### Dashboard Endpoint
- [ ] GET /dashboard documented with response schema
- [ ] Response includes all aggregated data fields
- [ ] Authentication requirements documented

### General
- [ ] Error response format documented for all endpoints
- [ ] Common error codes documented (401, 403, 404, 422)
- [ ] Rate limiting documented (if applicable)
- [ ] Example requests/responses provided

## Technical Notes
- Update existing API documentation files
- Use consistent documentation format across endpoints
- Include curl examples for common operations
- Document the X-User-ID header backward compatibility (if still supported)

## Dependencies
- Blocks: None
- Blocked by: All Phase 2 endpoints complete
- Related: 004-phase2-test-coverage-review

## Resources / Links
- ERD: phases/in-progress/002-frontend-readiness/sprints/todo/005-technical-debt-phase2/erd.md
- Existing API docs: docs/api/ (if exists)

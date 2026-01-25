# 005: Profile E2E Tests

## ERD Reference
Implements: All REQs (REQ-PROFILE-001 through REQ-PROFILE-006, REQ-AUTH-006, REQ-AUTH-007)

## Description
Implement comprehensive end-to-end tests for the profile endpoints. These tests verify the full request/response cycle including authentication, authorization, validation, and data persistence.

## Context / Background
E2E tests provide confidence that the profile feature works correctly from the API consumer's perspective. They test the integration of all components: routes, middleware, handlers, services, and database.

## Acceptance Criteria
- [ ] Test GET /users/{id}/profile:
  - Returns 200 with correct profile data for owner
  - Returns 200 with correct profile data for admin viewing other user
  - Returns 401 without authentication
  - Returns 403 for non-owner, non-admin
  - Returns 404 for non-existent user
  - Response has correct JSON structure (camelCase fields)
- [ ] Test PUT /users/{id}/profile:
  - Returns 200 and updates name
  - Returns 200 and updates weightUnit to "kg"
  - Returns 200 and updates weightUnit to "lb"
  - Returns 200 with partial update (only name)
  - Returns 200 with partial update (only weightUnit)
  - Returns 200 with empty body (no-op)
  - Clears name when empty string provided
  - Returns 400 for invalid weightUnit
  - Returns 400 for name > 100 chars
  - Returns 401 without authentication
  - Returns 403 for non-owner
  - Returns 403 for admin trying to update other user
  - Returns 404 for non-existent user
  - Verify updated_at changes on update
- [ ] Test weight_unit default:
  - New users have weight_unit = "lb"
- [ ] All tests use proper test isolation (clean database state)
- [ ] Tests work with both session auth and X-User-ID test header

## Technical Notes
- Follow existing E2E test patterns in the codebase
- Use test helpers for authentication (create session, get token)
- Create test users with known data for assertions
- Consider table-driven tests for validation cases
- Test both success and error response bodies

### Test Data Setup
- Create regular user A
- Create regular user B
- Create admin user
- Use these for owner/non-owner/admin test scenarios

### Key Scenarios
1. Owner views own profile - success
2. Owner updates own profile - success
3. Admin views other's profile - success
4. Admin updates other's profile - forbidden
5. User views other's profile - forbidden
6. Unauthenticated request - unauthorized

## Dependencies
- Blocks: None (final ticket in sprint)
- Blocked by: 001, 002, 003, 004 (all implementation must be complete)

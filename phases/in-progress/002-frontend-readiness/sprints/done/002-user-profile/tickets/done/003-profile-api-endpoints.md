# 003: Profile API Endpoints

## ERD Reference
Implements: REQ-PROFILE-002, REQ-PROFILE-003, REQ-PROFILE-004, REQ-PROFILE-005, REQ-PROFILE-006

## Description
Implement the HTTP handlers for GET and PUT /users/{id}/profile endpoints. These endpoints use the profile service for business logic and handle HTTP-specific concerns like request parsing and response formatting.

## Context / Background
The profile endpoints allow users to view and update their profile data. The GET endpoint returns the full profile. The PUT endpoint accepts partial updates. Both endpoints require authentication and will have authorization applied (separate ticket).

## Acceptance Criteria
- [ ] Implement GET /users/{id}/profile handler:
  - Extract user ID from path parameter
  - Call profile service GetProfile
  - Return 200 with profile JSON
  - Return 404 if user not found
  - Return 401 if not authenticated (middleware handles this)
- [ ] Implement PUT /users/{id}/profile handler:
  - Extract user ID from path parameter
  - Parse request body (name, weightUnit - both optional)
  - Call profile service UpdateProfile
  - Return 200 with updated profile JSON
  - Return 400 for validation errors
  - Return 404 if user not found
  - Return 401 if not authenticated (middleware handles this)
- [ ] Register routes with router
- [ ] Apply authentication middleware to both endpoints
- [ ] Integration tests for endpoint behavior

## Technical Notes
- Use the existing router pattern from the codebase
- Profile endpoints should be under /users/{id}/profile path
- Request body should use camelCase (weightUnit, not weight_unit)
- Response should use camelCase consistently
- Handler should not contain business logic - delegate to service
- Consider using a request struct for PUT body parsing

### Request/Response Examples

GET /users/{id}/profile Response (200):
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "name": "John Doe",
  "weightUnit": "lb",
  "createdAt": "2024-01-15T10:00:00Z",
  "updatedAt": "2024-01-15T10:00:00Z"
}
```

PUT /users/{id}/profile Request:
```json
{
  "name": "Jane Doe",
  "weightUnit": "kg"
}
```

## Dependencies
- Blocks: 005 (tests use these endpoints)
- Blocked by: 001 (schema), 002 (domain logic)

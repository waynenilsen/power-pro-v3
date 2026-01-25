# 002: Profile Domain Logic

## ERD Reference
Implements: REQ-PROFILE-002, REQ-PROFILE-003, REQ-PROFILE-004, REQ-PROFILE-005, REQ-PROFILE-006

## Description
Implement the profile service/domain logic that handles profile retrieval, updates, and validation. This provides the business logic layer between the API handlers and the database.

## Context / Background
The profile service encapsulates all profile-related business logic. This includes fetching user profiles with the correct fields, validating update requests, and applying partial updates. The service should be independent of HTTP concerns.

## Acceptance Criteria
- [ ] Create profile service/package with clear interface
- [ ] Implement GetProfile(userID) method:
  - Returns profile with id, email, name, weight_unit, created_at, updated_at
  - Returns error if user not found
- [ ] Implement UpdateProfile(userID, updates) method:
  - Accepts partial updates (only provided fields)
  - Validates name length (<= 100 chars)
  - Validates weight_unit ("lb" or "kg")
  - Empty string for name sets it to null
  - Returns updated profile
  - Returns validation error with descriptive message
- [ ] Profile struct/type with JSON tags matching API contract:
  - `id`, `email`, `name`, `weightUnit`, `createdAt`, `updatedAt`
- [ ] Unit tests for validation logic
- [ ] Unit tests for partial update behavior

## Technical Notes
- Profile service should depend on a repository/store interface for database access
- Keep validation logic in the service, not in the handler
- Consider using a separate UpdateProfileRequest struct for updates
- Empty string vs nil distinction: empty string clears the field, nil/omitted leaves it unchanged
- Use camelCase for JSON field names (Go convention for JSON APIs)

## Dependencies
- Blocks: 003 (endpoints use this service), 005 (tests verify this logic)
- Blocked by: 001 (schema must exist)

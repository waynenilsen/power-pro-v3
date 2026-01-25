# 004: Profile Authorization

## ERD Reference
Implements: REQ-AUTH-006, REQ-AUTH-007

## Description
Implement authorization rules for profile endpoints. Users can only access their own profile, with admins having read-only access to any profile. Update operations are owner-only (even admins cannot update other users' profiles).

## Context / Background
Authorization ensures users cannot view or modify other users' profiles. This is distinct from authentication (which verifies identity). The auth middleware from Sprint 001 provides the current user context; this ticket adds the authorization layer.

## Acceptance Criteria
- [ ] GET /users/{id}/profile authorization:
  - Allow if requester ID matches profile user ID (owner)
  - Allow if requester is admin (is_admin = true)
  - Return 403 Forbidden otherwise
- [ ] PUT /users/{id}/profile authorization:
  - Allow only if requester ID matches profile user ID (owner)
  - Return 403 Forbidden for non-owners (including admins)
- [ ] 403 error response format:
  - `{"error": "Forbidden: you can only access your own profile"}`
  - For admin PUT attempts: `{"error": "Forbidden: profile updates are owner-only"}`
- [ ] Authorization check happens before any database access
- [ ] Unit tests for authorization logic
- [ ] Integration tests verifying 403 responses

## Technical Notes
- Authorization can be implemented as middleware or within handlers
- Current user info (ID, is_admin) should come from request context (set by auth middleware)
- Consider extracting authorization logic into a reusable function/middleware
- Error messages should not leak information about whether user exists
- Pattern: check auth → check authz → perform operation

### Authorization Matrix

| Requester | GET own profile | GET other profile | PUT own profile | PUT other profile |
|-----------|-----------------|-------------------|-----------------|-------------------|
| Owner     | 200             | 403               | 200             | 403               |
| Admin     | 200             | 200               | 200             | 403               |
| Other     | 403             | 403               | 403             | 403               |

## Dependencies
- Blocks: 005 (tests verify authorization)
- Blocked by: 001 (schema), 002 (domain logic), 003 (endpoints exist first)

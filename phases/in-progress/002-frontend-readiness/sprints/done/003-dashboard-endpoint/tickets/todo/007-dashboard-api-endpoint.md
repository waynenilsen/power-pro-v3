# 007: Dashboard API Endpoint

## ERD Reference
Implements: REQ-DASH-001, REQ-DASH-008

## Description
Create the HTTP handler for the dashboard endpoint that wires up the dashboard service and handles authentication, authorization, and response formatting.

## Context / Background
This ticket creates the API layer that exposes the dashboard service. It follows existing API patterns in PowerPro, using the session middleware for authentication and implementing owner-only authorization.

## Acceptance Criteria
- [ ] Create `GET /users/{id}/dashboard` endpoint
- [ ] Require valid session token in Authorization header
- [ ] Return 401 if no valid session
- [ ] Return 403 if authenticated user is not the owner (no admin override)
- [ ] Return 404 if user does not exist
- [ ] Return 200 with dashboard JSON on success
- [ ] Response uses camelCase for JSON fields (Go struct tags)
- [ ] Error responses follow existing API error format
- [ ] Register route in router configuration
- [ ] Integration test for the endpoint

## Technical Notes
- Follow existing handler patterns in `internal/api/` or similar
- Use existing auth middleware to extract user from session
- Compare path parameter {id} with session user ID for authorization
- Dashboard service handles all business logic - handler is thin
- Consider response struct with json tags for camelCase output

## Response Examples

Success (200):
```json
{
  "enrollment": { ... },
  "nextWorkout": { ... },
  "currentSession": null,
  "recentWorkouts": [ ... ],
  "currentMaxes": [ ... ]
}
```

Unauthorized (401):
```json
{
  "error": "Unauthorized"
}
```

Forbidden (403):
```json
{
  "error": "Forbidden"
}
```

## Dependencies
- Blocks: 008 (E2E tests)
- Blocked by: 001 (dashboard service)

# Sprint 002: User Profile

## Overview

Implement user profile viewing and updating functionality. Users can view and update their profile data including name and weight unit preference (lb/kg).

## Reference Documents

- PRD: `phases/in-progress/002-frontend-readiness/sprints/todo/002-user-profile/prd.md`
- ERD: `phases/in-progress/002-frontend-readiness/sprints/todo/002-user-profile/erd.md`
- Tickets: `phases/in-progress/002-frontend-readiness/sprints/todo/002-user-profile/tickets/todo/`

## Requirements Summary

### Schema Change
- Add `weight_unit` column to users table (TEXT, values: "lb" or "kg", default: "lb", NOT NULL)

### Endpoints
- `GET /users/{id}/profile` - Get user profile (id, email, name, weightUnit, createdAt, updatedAt)
- `PUT /users/{id}/profile` - Update profile (name, weightUnit)

### Authorization
- Users can only access/update their own profile
- Admins can view (but not update) any profile
- Return 403 for unauthorized access

### Validation
- Name max 100 characters
- weightUnit must be "lb" or "kg"
- Empty string for name clears it (sets to null)

## Tickets (in order)

1. `001-profile-schema-migration.md` - Add weight_unit column
2. `002-profile-domain-logic.md` - Profile validation and business logic
3. `003-profile-api-endpoints.md` - GET/PUT handlers
4. `004-profile-authorization.md` - Owner/admin access rules
5. `005-profile-e2e-tests.md` - End-to-end tests

## Dependencies

- Sprint 001 (Authentication) must be complete - it is!

## Success Criteria

- User can view their profile via GET endpoint
- User can update name and weight_unit via PUT endpoint
- Authorization properly enforced
- All E2E tests pass

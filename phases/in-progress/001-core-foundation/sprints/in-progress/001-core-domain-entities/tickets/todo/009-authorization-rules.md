# 009: Authorization Rules for Lift and LiftMax

## ERD Reference
Implements: NFR-006, NFR-007

## Description
Implement authorization rules for Lift and LiftMax entities as specified in the non-functional requirements. Lift data should be readable by all authenticated users, while LiftMax data should only be accessible to the owning user or admins.

## Context / Background
Security is critical for user data protection. While lifts are shared reference data (all users see the same Squat definition), LiftMax records contain personal training data that users expect to be private.

## Acceptance Criteria
- [ ] Lift authorization (NFR-007):
  - All authenticated users can read lift data (GET /lifts, GET /lifts/{id})
  - Only admins can create/update/delete lifts (POST, PUT, DELETE /lifts)
  - Unauthenticated requests receive 401
  - Non-admin write requests receive 403
- [ ] LiftMax authorization (NFR-006):
  - Users can only access their own LiftMax data
  - Admins can access any user's LiftMax data
  - Ownership verified by comparing authenticated user ID with resource's user_id
  - Unauthenticated requests receive 401
  - Non-owner, non-admin requests receive 403
- [ ] Authorization middleware/guards implemented as reusable components
- [ ] Authorization checks happen before any database operations
- [ ] Clear, consistent error messages for authorization failures
- [ ] Security tests cover:
  - Authenticated user accessing own data (allowed)
  - Authenticated user accessing other user's data (denied)
  - Admin accessing other user's data (allowed)
  - Unauthenticated access attempts (denied)
  - Non-admin attempting lift write operations (denied)

## Technical Notes
- Implement as middleware/guards that can be applied to routes
- Consider using decorators or route-level middleware depending on framework
- Authorization should be checked after authentication but before route handler
- Log authorization failures for security monitoring (without exposing sensitive data)
- Assumes authentication system already in place per ERD Section 6 Assumptions

## Dependencies
- Blocks: None (authorization is applied across other tickets)
- Blocked by: None (can be developed in parallel, applied during API integration)
- Related: 003 (Lift CRUD API), 006 (LiftMax CRUD API), 007, 008 (LiftMax endpoints)

## Resources / Links
- ERD: phases/todo/001-core-foundation/sprints/in-progress/001-core-domain-entities/erd.md
- NFR-006: LiftMax access control requirement
- NFR-007: Lift access control requirement
- Auth System: prompts/auth.md

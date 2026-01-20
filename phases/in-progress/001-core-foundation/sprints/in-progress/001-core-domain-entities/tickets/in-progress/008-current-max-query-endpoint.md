# 008: Current Max Query Endpoint

## ERD Reference
Implements: REQ-MAX-011
Related to: NFR-002

## Description
Implement efficient API endpoint for looking up a user's current (most recent) max value for a specific lift and type. This is the primary endpoint used by load calculation features.

## Context / Background
Load calculation requires finding the most recent max for a given (user, lift, type) combination. This is a high-frequency operation that must be optimized for performance (NFR-002 requires < 50ms p95). The current max is determined by the most recent `effective_date`.

## Acceptance Criteria
- [ ] `GET /users/{userId}/lift-maxes/current` - Get current max
  - Query params:
    - `lift` (required): Lift ID (UUID)
    - `type` (required): Max type ('ONE_RM' or 'TRAINING_MAX')
  - Returns the LiftMax record with the most recent `effective_date` for the specified (user, lift, type)
  - Response includes full LiftMax object
  - Returns 404 if no max exists for the specified combination
  - Returns 400 if:
    - `lift` param is missing or invalid UUID
    - `type` param is missing or invalid
  - Returns 403 if requesting user is not the owner or admin
- [ ] Response time < 50ms (p95) per NFR-002
- [ ] Query uses appropriate index on (`user_id`, `lift_id`, `type`, `effective_date`)
- [ ] Integration tests cover:
  - Successful lookup with single max
  - Successful lookup with multiple maxes (returns most recent)
  - No max exists (404)
  - Invalid parameters (400)
  - Authorization (403)
  - Performance under load

## Technical Notes
- Query should use: `ORDER BY effective_date DESC LIMIT 1`
- Ensure index from Ticket 004 is utilized
- Consider query plan analysis to verify index usage
- This endpoint is critical path for load calculation - optimize aggressively
- Could benefit from caching with invalidation on max create/update

## Dependencies
- Blocks: None (this is a leaf endpoint)
- Blocked by: 004 (Schema with index), 006 (Base CRUD API for auth pattern)
- Related: 007 (Similar authorization pattern), 009 (Authorization rules)

## Resources / Links
- ERD: phases/todo/001-core-foundation/sprints/in-progress/001-core-domain-entities/erd.md
- REQ-MAX-011: Current max query specification
- NFR-002: Response time requirement (< 50ms p95)

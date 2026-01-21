# 011: Progression History Query API

## ERD Reference
Implements: REQ-HIST-002
Related to: REQ-HIST-001

## Description
Implement the REST API endpoint for querying progression history. This enables users to see their progression applications over time, supporting progress tracking and debugging.

## Context / Background
Users want to see how their training loads have progressed over time. This API exposes the ProgressionLog data in a user-friendly format, supporting filtering by lift, date range, and progression type. It's also valuable for debugging progression issues and verifying correct application.

## Acceptance Criteria
- [ ] Implement `GET /users/{userId}/progression-history` endpoint:
  - Returns chronological list of progression applications
  - Default order: most recent first (descending by appliedAt)
  - Includes: id, progressionId, liftId, previousValue, newValue, delta, triggerType, triggerContext, appliedAt
  - Join progression name for readability
  - Join lift name for readability
- [ ] Query parameter filters:
  - `liftId` - filter to specific lift
  - `progressionType` - filter by progression type (LINEAR_PROGRESSION, CYCLE_PROGRESSION)
  - `triggerType` - filter by trigger type (AFTER_SESSION, AFTER_WEEK, AFTER_CYCLE)
  - `startDate` - filter applications after this date (ISO 8601)
  - `endDate` - filter applications before this date (ISO 8601)
- [ ] Pagination support:
  - `limit` - number of results (default 20, max 100)
  - `offset` - skip N results for pagination
  - Response includes total count for pagination UI
- [ ] Authorization:
  - Users can only query their own progression history
  - Return 403 for unauthorized access attempts
  - Admin users can query any user's history
- [ ] Response format matches ERD ProgressionLog Entry specification
- [ ] Performance: Query returns in < 200ms for typical history (< 1000 entries)
- [ ] Test coverage for filtering, pagination, and authorization

## Technical Notes
- Use indexed queries on (user_id, lift_id) and applied_at
- Consider cursor-based pagination for large histories
- TriggerContext is JSONB - include as-is in response
- Join progressions and lifts tables for names
- Date filtering uses applied_at column

## Dependencies
- Blocks: None
- Blocked by: 002 (ProgressionLog schema)
- Related: 008 (Trigger integration creates log entries)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/004-progression-rules/erd.md
- ProgressionLog Entry format: Section 5 of ERD

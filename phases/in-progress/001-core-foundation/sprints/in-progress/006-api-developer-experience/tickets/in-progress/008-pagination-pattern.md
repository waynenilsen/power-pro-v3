# 008: Pagination Pattern

## ERD Reference
Implements: REQ-API-003

## Description
Implement and document consistent pagination for all list endpoints. Different pagination patterns across endpoints confuse developers.

## Context / Background
List endpoints can return large result sets. Consistent pagination enables clients to implement generic paging logic that works across all list endpoints.

## Acceptance Criteria
- [ ] All list endpoints support pagination
- [ ] Pagination parameters are consistent (limit, offset or cursor)
- [ ] Response includes total count where practical
- [ ] Response includes pagination metadata
- [ ] Pagination behavior documented

## Technical Notes
- **Pagination Strategy Options**:
  - **Offset-based**: `?limit=10&offset=20` - Simple, but can have issues with large offsets
  - **Cursor-based**: `?limit=10&cursor=abc123` - Better for real-time data, more complex

- **Recommendation**: Offset-based for simplicity (SQLite handles well)

- **Pagination Parameters** (suggested):
  - `limit` - Number of items to return (default: 20, max: 100)
  - `offset` - Number of items to skip (default: 0)

- **Pagination Response** (suggested):
  ```json
  {
    "data": [...],
    "meta": {
      "total": 150,
      "limit": 20,
      "offset": 40,
      "hasMore": true
    }
  }
  ```

- **List Endpoints to Support Pagination**:
  - GET /lifts
  - GET /lift-maxes
  - GET /prescriptions
  - GET /days
  - GET /weeks
  - GET /cycles
  - GET /programs
  - GET /progressions
  - GET /progression-logs

- Audit current pagination implementations
- Standardize parameters and response format
- Document pagination pattern

## Dependencies
- Blocks: None
- Blocked by: 007-response-envelope-standardization
- Related: 009-filtering-pattern

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/006-api-developer-experience/erd.md

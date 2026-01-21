# 009: Filtering Pattern

## ERD Reference
Implements: REQ-API-004

## Description
Implement and document consistent filtering for list endpoints. Consistent filtering makes queries predictable.

## Context / Background
Developers need to filter list results to find specific resources. Inconsistent filter parameter naming or behavior across endpoints creates confusion and errors.

## Acceptance Criteria
- [ ] Filter parameters follow consistent naming
- [ ] Filter behavior is documented
- [ ] Complex filters use consistent syntax

## Technical Notes
- **Filter Parameter Naming Convention** (suggested):
  - Use field names directly: `?userId=123&liftId=456`
  - For date ranges: `?createdAfter=2024-01-01&createdBefore=2024-12-31`
  - For status filters: `?status=active`

- **Filter Operators** (if needed):
  - For simple equality, use direct value: `?status=active`
  - For ranges, use suffixes: `?weightGte=100&weightLte=200` (gte = greater than or equal, lte = less than or equal)

- **Endpoints with Filtering**:
  - GET /lifts - filter by userId
  - GET /lift-maxes - filter by liftId, userId
  - GET /prescriptions - filter by liftId
  - GET /progression-logs - filter by userId, liftId, date range
  - Others as applicable

- **Documentation Requirements**:
  - List available filters for each endpoint
  - Document filter syntax
  - Explain behavior when multiple filters are combined (AND vs OR)
  - Document default behavior when no filters provided

- Audit current filtering implementations
- Standardize filter parameter naming
- Document filtering pattern

## Dependencies
- Blocks: None
- Blocked by: 007-response-envelope-standardization
- Related: 008-pagination-pattern

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/006-api-developer-experience/erd.md

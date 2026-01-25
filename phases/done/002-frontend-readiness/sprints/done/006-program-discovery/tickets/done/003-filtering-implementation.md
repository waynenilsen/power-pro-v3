# 003: Program Filtering Query Implementation

## ERD Reference
Implements: REQ-FILTER-001, REQ-FILTER-002, REQ-FILTER-003, REQ-FILTER-004, REQ-FILTER-005

## Description
Implement filtering capabilities for the GET /programs endpoint. Users should be able to filter programs by difficulty, days_per_week, focus, and has_amrap. Multiple filters should combine with AND logic.

## Context / Background
With discovery metadata columns in place, users need API support to filter programs based on their constraints and goals. A beginner with 3 days per week to train should be able to find programs matching both criteria. Filtering is essential for discovery at scale.

## Acceptance Criteria
- [ ] Update domain/program types to include filter options struct:
  - Difficulty filter (optional string)
  - DaysPerWeek filter (optional int)
  - Focus filter (optional string)
  - HasAmrap filter (optional bool)
- [ ] Update ProgramRepository.List to accept filter options
- [ ] Implement filter validation:
  - Difficulty must be: beginner, intermediate, advanced
  - DaysPerWeek must be: 1-7
  - Focus must be: strength, hypertrophy, peaking
  - HasAmrap must be: true, false
- [ ] Update program_handler to parse query parameters:
  - `?difficulty=beginner`
  - `?days_per_week=3`
  - `?focus=strength`
  - `?has_amrap=true`
- [ ] Return 400 Bad Request for invalid filter values with clear error message
- [ ] Combined filters use AND logic:
  - `?difficulty=beginner&days_per_week=3` returns programs matching BOTH
- [ ] Empty result returns 200 with empty array (not 404)
- [ ] Update Program domain model to include new fields
- [ ] Update ProgramResponse DTO to include new fields (camelCase)
- [ ] Add unit tests for filter validation
- [ ] Add integration tests for filter queries

## Technical Notes
- Build WHERE clause dynamically based on which filters are provided
- Use parameterized queries to prevent SQL injection
- Filter validation should use shared validation functions (similar to existing validation patterns)
- Query builder pattern may help manage optional WHERE clauses
- Consider: `WHERE 1=1 AND (difficulty = ? OR ? IS NULL) ...` pattern OR dynamic query building
- Update both list and count queries if pagination exists
- Ensure indices from ticket 001 are utilized

## Dependencies
- Blocks: 006 (E2E tests for filtering)
- Blocked by: 001 (schema must have columns)

## Resources / Links
- Existing handler patterns: `internal/api/program_handler.go`
- Domain model: `internal/domain/program/program.go`
- Repository: `internal/repository/program_repository.go`
- ERD requirements: section 3, "Program Filtering"

# 004: Program Search Implementation

## ERD Reference
Implements: REQ-SEARCH-001, REQ-SEARCH-002

## Description
Implement search capability for the GET /programs endpoint. Users should be able to search programs by name using a substring match. Search should be case-insensitive and combinable with filters.

## Context / Background
Users who know a program name (or partial name) should be able to find it quickly. "strength" should find "Starting Strength". "531" should find "Wendler 5/3/1". Search combined with filters enables powerful discovery (e.g., search "5" with difficulty=intermediate).

## Acceptance Criteria
- [ ] Update filter options struct to include Search field (optional string)
- [ ] Update ProgramRepository.List to handle search parameter
- [ ] Implement search query:
  - Case-insensitive matching (use LOWER() or COLLATE NOCASE)
  - Substring matching (not just prefix)
  - `?search=strength` returns Starting Strength
  - `?search=531` returns Wendler 5/3/1
  - `?search=gz` returns GZCLP
- [ ] Empty search parameter is ignored (returns all programs)
- [ ] Search combines with filters using AND logic:
  - `?search=strength&difficulty=beginner` returns Starting Strength
  - `?search=method&days_per_week=4` returns nothing (Texas Method is 3 days)
- [ ] Update program_handler to parse `search` query parameter
- [ ] Add unit tests for search functionality
- [ ] Add integration tests for search + filter combinations

## Technical Notes
- SQLite LIKE pattern: `WHERE LOWER(name) LIKE '%' || LOWER(?) || '%'`
- Alternative: `WHERE name LIKE '%' || ? || '%' COLLATE NOCASE`
- Sanitize search input (escape % and _ characters if used literally)
- Index on name exists from original schema; LIKE with prefix can use index but substring cannot
- For small datasets (< 1000 programs), full table scan is acceptable
- Consider: if search is empty string, skip adding search condition entirely
- Integrate with existing filter query building from ticket 003

## Dependencies
- Blocks: 006 (E2E tests for search)
- Blocked by: 003 (search builds on filter implementation)

## Resources / Links
- Existing handler patterns: `internal/api/program_handler.go`
- Repository: `internal/repository/program_repository.go`
- ERD requirements: section 3, "Program Search"

# 012: Program Entity CRUD API

## ERD Reference
Implements: REQ-PROG-001, REQ-PROG-002

## Description
Implement the Program entity repository, service, and CRUD API endpoints. A Program bundles cycles, lookups, and defaults into a named configuration.

## Context / Background
A Program is the top-level configuration users enroll in. It references a cycle, optional lookup tables, and default settings. Examples: "5/3/1 BBB", "Starting Strength", "Bill Starr 5x5".

## Acceptance Criteria
- [ ] Program repository implemented with CRUD operations
- [ ] Program service implemented with business logic
- [ ] GET /programs - list all programs
- [ ] GET /programs/{id} - get program with full structure (cycle, lookups)
- [ ] POST /programs - create program
- [ ] PUT /programs/{id} - update program
- [ ] DELETE /programs/{id} - delete program (fails if users enrolled)
- [ ] Program includes: name, slug, description, cycle_id, lookup references, default_rounding
- [ ] Slug uniqueness validation
- [ ] Unit tests with >80% coverage
- [ ] Integration tests for all endpoints
- [ ] NFR-003: Operations complete in <100ms (p95)

## Technical Notes
- Program response includes embedded cycle structure
- Program response includes lookup table references (not full entries)
- DELETE fails with 409 Conflict if any users are enrolled
- default_rounding: weight rounding increment (e.g., 2.5, 5.0)

## Dependencies
- Blocks: 013, 014
- Blocked by: 005 (Program schema), 009 (Cycle API), 010 (Lookup API)
- Related: All schedule entities

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/003-schedule-periodization/erd.md

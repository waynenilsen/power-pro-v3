# 009: Cycle Entity CRUD API

## ERD Reference
Implements: REQ-CYCLE-001, REQ-CYCLE-002, REQ-CYCLE-003, REQ-CYCLE-004

## Description
Implement the Cycle entity repository, service, and CRUD API endpoints. A Cycle is the repeating unit of a program (1-week, 3-week, 4-week, etc.).

## Context / Background
Programs repeat in cycles. 5/3/1 uses 4-week cycles, Greg Nuckols uses 3-week cycles, Starting Strength uses 1-week cycles. The Cycle entity defines the structure and length of these repeating units.

## Acceptance Criteria
- [ ] Cycle repository implemented with CRUD operations
- [ ] Cycle service implemented with business logic
- [ ] GET /cycles - list all cycles (with pagination)
- [ ] GET /cycles/{id} - get cycle with weeks
- [ ] POST /cycles - create cycle with name, length_weeks
- [ ] PUT /cycles/{id} - update cycle
- [ ] DELETE /cycles/{id} - delete cycle (fails if used by programs)
- [ ] Validation: count of associated weeks must equal length_weeks
- [ ] Cycle length validation (>= 1)
- [ ] Unit tests with >80% coverage
- [ ] Integration tests for all endpoints

## Technical Notes
- Cycle response includes embedded weeks with week_number
- Validation at application layer: weeks.count() == cycle.length_weeks
- Cycles can have 1, 3, 4, or any positive integer weeks
- A cycle with no weeks is valid initially (weeks added separately)

## Dependencies
- Blocks: 012, 014
- Blocked by: 002, 003 (Week and Cycle schemas)
- Related: 008 (Week API)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/003-schedule-periodization/erd.md
- API Response Format: See ERD Section 5 "Cycle API Response"

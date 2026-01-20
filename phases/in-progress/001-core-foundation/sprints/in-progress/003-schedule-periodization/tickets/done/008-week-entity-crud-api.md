# 008: Week Entity CRUD API

## ERD Reference
Implements: REQ-WEEK-001, REQ-WEEK-002, REQ-WEEK-003, REQ-WEEK-004

## Description
Implement the Week entity repository, service, and CRUD API endpoints. A Week represents a collection of training days within a program cycle.

## Context / Background
A Week contains training days mapped to specific days of the week. Weeks are numbered within cycles and support A/B week variants for alternating programs.

## Acceptance Criteria
- [ ] Week repository implemented with CRUD operations
- [ ] Week service implemented with business logic
- [ ] GET /weeks - list all weeks (with pagination)
- [ ] GET /weeks/{id} - get week with day mappings
- [ ] POST /weeks - create week with week_number, variant, cycle_id
- [ ] PUT /weeks/{id} - update week
- [ ] DELETE /weeks/{id} - delete week (fails if part of active cycle)
- [ ] POST /weeks/{id}/days - add day to week with day_of_week
- [ ] DELETE /weeks/{id}/days/{dayId} - remove day from week
- [ ] Week number unique within cycle validation
- [ ] Support A/B variant field (null, "A", or "B")
- [ ] Unit tests with >80% coverage
- [ ] Integration tests for all endpoints

## Technical Notes
- Week response includes embedded day mappings
- day_of_week: MONDAY, TUESDAY, WEDNESDAY, THURSDAY, FRIDAY, SATURDAY, SUNDAY
- Multiple days can map to same day_of_week (for two-a-days)
- Variant field enables A/B week rotation patterns

## Dependencies
- Blocks: 009, 014
- Blocked by: 001, 002 (Day and Week schemas)
- Related: 007 (Day API)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/003-schedule-periodization/erd.md
- API Response Format: See ERD Section 5 "Week API Response"

# 010: Lookup Tables CRUD API

## ERD Reference
Implements: REQ-LOOKUP-001, REQ-LOOKUP-002, REQ-LOOKUP-003, REQ-LOOKUP-004

## Description
Implement the WeeklyLookup and DailyLookup entity repositories, services, and CRUD API endpoints. Lookup tables enable programs to vary percentages based on week number or day identifier.

## Context / Background
Programs like 5/3/1 vary percentages by week (Week 1 = 65/75/85%). Bill Starr varies intensity by day (Heavy/Light/Medium). Lookup tables enable this variation without duplicating prescriptions.

## Acceptance Criteria
- [ ] WeeklyLookup repository and service implemented
- [ ] DailyLookup repository and service implemented
- [ ] GET /weekly-lookups - list all weekly lookups
- [ ] GET /weekly-lookups/{id} - get weekly lookup with entries
- [ ] POST /weekly-lookups - create weekly lookup
- [ ] PUT /weekly-lookups/{id} - update weekly lookup
- [ ] DELETE /weekly-lookups/{id} - delete (fails if used by programs)
- [ ] GET /daily-lookups - list all daily lookups
- [ ] GET /daily-lookups/{id} - get daily lookup with entries
- [ ] POST /daily-lookups - create daily lookup
- [ ] PUT /daily-lookups/{id} - update daily lookup
- [ ] DELETE /daily-lookups/{id} - delete (fails if used by programs)
- [ ] WeeklyLookup.GetByWeekNumber(weekNum) method
- [ ] DailyLookup.GetByDayIdentifier(daySlug) method
- [ ] Unit tests with >80% coverage
- [ ] Integration tests for all endpoints

## Technical Notes
- WeeklyLookup entries: array of {weekNumber, percentages[], reps[], percentageModifier}
- DailyLookup entries: map of {dayIdentifier -> percentageModifier, intensityLevel}
- Entries stored as JSONB
- Lookup queries: provide lookup ID + position key, get parameters back

## Dependencies
- Blocks: 011, 012
- Blocked by: 004 (Lookup tables schema)
- Related: ERD-002 LoadStrategy

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/003-schedule-periodization/erd.md
- API Response Format: See ERD Section 5 "WeeklyLookup Example"

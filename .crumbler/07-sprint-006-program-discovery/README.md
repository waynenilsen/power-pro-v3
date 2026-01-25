# Sprint 006 - Program Discovery

## Overview

Enhance program listing and detail endpoints to help users find appropriate programs.

## Enhancements to GET /programs

Add filter parameters:
- `difficulty` (beginner, intermediate, advanced)
- `days_per_week` (1-7)
- `focus` (strength, hypertrophy, peaking)
- `has_amrap` (true/false)
- `search` (name substring)

## Schema Changes

Add metadata columns to `programs` table:
- `difficulty`
- `days_per_week`
- `focus`
- `has_amrap`

Backfill existing programs with appropriate values.

## Program Detail Enhancement

`GET /programs/{id}` should include:
- Sample week preview (day names, exercise counts)
- Lift requirements (what lifts the program uses)
- Estimated session duration

## Tasks

1. Create sprint directory: `phases/in-progress/002-frontend-readiness/sprints/todo/006-program-discovery/`
2. Write `prd.md` - Product requirements for program discovery
3. Write `erd.md` - Engineering requirements with detailed specs
4. Create ticket directory structure
5. Create tickets (at least 5):
   - Schema migration for program metadata columns
   - Backfill migration for existing programs
   - Program filtering query implementation
   - Program detail enhancement (sample week, lift requirements)
   - Program search implementation
   - E2E tests for program discovery

## Acceptance Criteria

- [ ] Sprint directory structure exists
- [ ] PRD document completed
- [ ] ERD document completed with REQ-XXX requirements
- [ ] At least 5 tickets created in `tickets/todo/`
- [ ] Tickets reference ERD requirements

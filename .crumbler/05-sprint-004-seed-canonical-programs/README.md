# Sprint 004 - Seed Canonical Programs

## Overview

Seed 3-5 popular powerlifting programs so users have programs to enroll in.

## Programs to Seed

- **Starting Strength** - beginner, 3 days/week, linear progression
- **Texas Method** - intermediate, 3 days/week, weekly progression
- **5/3/1** - intermediate, 4 days/week, monthly progression
- **GZCLP** - beginner/intermediate, 3-4 days/week, tiered progression

## Implementation Approach

- Add migration that creates full program structure:
  - Lookups, prescriptions, days, weeks, cycles, programs, progressions
- Use canonical slugs: `starting-strength`, `texas-method`, `531`, `gzclp`
- E2E tests continue using unique slugs (`texas-method-{testID}`) - no conflict

## Design Notes

- This is a large migration but mostly mechanical
- Translate E2E test setup code patterns into SQL INSERTs
- Reference program documentation in `programs/*.md` for accuracy

## Tasks

1. Create sprint directory: `phases/in-progress/002-frontend-readiness/sprints/todo/004-seed-canonical-programs/`
2. Write `prd.md` - Product requirements for program seeding
3. Write `erd.md` - Engineering requirements with detailed specs
4. Create ticket directory structure
5. Create tickets (at least 5):
   - Starting Strength program seed migration
   - Texas Method program seed migration
   - 5/3/1 program seed migration
   - GZCLP program seed migration
   - Program seed verification tests
   - Documentation for canonical programs

## Acceptance Criteria

- [ ] Sprint directory structure exists
- [ ] PRD document completed
- [ ] ERD document completed with REQ-XXX requirements
- [ ] At least 5 tickets created in `tickets/todo/`
- [ ] Tickets reference ERD requirements

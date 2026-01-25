# Sprint 006: Program Discovery

## Overview

Implement program discovery features: filtering by metadata, search by name, and enhanced program detail views. Enables users to find programs matching their experience level, schedule, and goals.

## Reference Documents

- PRD: `phases/in-progress/002-frontend-readiness/sprints/todo/006-program-discovery/prd.md`
- ERD: `phases/in-progress/002-frontend-readiness/sprints/todo/006-program-discovery/erd.md`
- Tickets: `phases/in-progress/002-frontend-readiness/sprints/todo/006-program-discovery/tickets/todo/`

## Requirements Summary

### Schema Changes (programs table)
- `difficulty` TEXT: 'beginner', 'intermediate', 'advanced' (default: 'beginner')
- `days_per_week` INTEGER: 1-7 (default: 3)
- `focus` TEXT: 'strength', 'hypertrophy', 'peaking' (default: 'strength')
- `has_amrap` INTEGER: 0 or 1 (default: 0)

### Backfill Values
- Starting Strength: beginner, 3, strength, 0
- Texas Method: intermediate, 3, strength, 0
- Wendler 5/3/1: intermediate, 4, strength, 1
- GZCLP: beginner, 4, strength, 1

### Filter Endpoints
```
GET /programs?difficulty=beginner
GET /programs?days_per_week=3
GET /programs?focus=strength
GET /programs?has_amrap=true
GET /programs?difficulty=beginner&days_per_week=3  # Combined filters
```

### Search
```
GET /programs?search=strength  # Case-insensitive name search
GET /programs?search=strength&difficulty=beginner  # Combined with filters
```

### Enhanced Program Detail
`GET /programs/{id}` now includes:
- `sampleWeek` - Array of days with exercise counts
- `liftRequirements` - Unique lifts used in program
- `estimatedSessionMinutes` - Approximate workout duration

## Tickets (in order)

1. `001-schema-migration.md` - Add metadata columns with indices
2. `002-backfill-migration.md` - Backfill canonical programs
3. `003-filtering-implementation.md` - Filter query support
4. `004-search-implementation.md` - Name search support
5. `005-detail-enhancement.md` - Sample week, lifts, duration
6. `006-e2e-tests.md` - End-to-end tests

## Dependencies

- Sprint 004 (Seed Programs) must be complete for backfill

## Success Criteria

- All filters work individually and combined
- Search returns relevant results
- Program detail includes new fields
- All queries < 100ms
- E2E tests pass

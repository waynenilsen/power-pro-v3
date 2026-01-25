# Sprint 004: Seed Canonical Programs

## Overview

Seed four canonical powerlifting programs into the database via migrations. These programs provide users with immediate access to proven training methodologies without requiring manual configuration.

## Reference Documents

- PRD: `phases/in-progress/002-frontend-readiness/sprints/todo/004-seed-canonical-programs/prd.md`
- ERD: `phases/in-progress/002-frontend-readiness/sprints/todo/004-seed-canonical-programs/erd.md`
- Tickets: `phases/in-progress/002-frontend-readiness/sprints/todo/004-seed-canonical-programs/tickets/todo/`
- Program details: `programs/*.md`

## Programs to Seed

### 1. Starting Strength (`starting-strength`)
- Novice program, 3 days/week, linear progression
- Workout A: Squat 3x5, Bench 3x5, Deadlift 1x5
- Workout B: Squat 3x5, OHP 3x5, Power Clean 5x3
- +5 lbs/session (upper), +10 lbs/session (deadlift)

### 2. Texas Method (`texas-method`)
- Intermediate program, 3 days/week, weekly periodization
- Volume Day: Squat 5x5@90%, Bench/Press 5x5@90%, Deadlift 1x5
- Recovery Day: Squat 2x5@80%, Bench/Press 3x5@90%
- Intensity Day: Squat 1x5@100%, Bench/Press 1x5@100%
- +5 lbs/week

### 3. Wendler 5/3/1 (`531`)
- Intermediate program, 4 days/week, monthly cycles
- 4-week cycle: 5s week, 3s week, 5/3/1 week, deload
- Week 1: 65%x5, 75%x5, 85%x5+
- Week 2: 70%x3, 80%x3, 90%x3+
- Week 3: 75%x5, 85%x3, 95%x1+
- Week 4: 40%x5, 50%x5, 60%x5
- +5 lbs upper, +10 lbs lower per cycle

### 4. GZCLP (`gzclp`)
- Beginner/intermediate program, 4 days/week, tiered progression
- T1 lifts: 5x3+ -> 6x2+ -> 10x1+ (on failure)
- T2 lifts: 3x10 -> 3x8 -> 3x6 (on failure)
- Linear progression per successful session

## Tickets (in order)

1. `001-starting-strength-seed.md` - Starting Strength migration
2. `002-texas-method-seed.md` - Texas Method migration
3. `003-531-seed.md` - Wendler 5/3/1 migration
4. `004-gzclp-seed.md` - GZCLP migration
5. `005-program-verification-tests.md` - Structure and prescription tests
6. `006-canonical-programs-documentation.md` - Document canonical slugs

## Dependencies

- Existing program infrastructure from Phase 001

## Success Criteria

- All 4 programs seeded with correct structure
- Prescriptions accurate per program documentation
- Verification tests pass
- Programs discoverable via /programs endpoint

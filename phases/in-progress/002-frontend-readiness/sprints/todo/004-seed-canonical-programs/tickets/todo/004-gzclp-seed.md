# 004: GZCLP Program Seed

## ERD Reference
Implements: REQ-SEED-004

## Description
Create a goose migration that seeds the GZCLP (GZCL Linear Progression) program into the database. This program uses a tiered approach with T1 (high intensity), T2 (moderate intensity), and progression through rep scheme changes on failure.

## Context / Background
GZCLP is built on the tiered system: T1 exercises are heavy compound movements, T2 exercises are moderate-intensity compounds that support T1 lifts. The unique feature is progression through failure - when you can't complete the prescribed reps, you move to a different rep scheme at the same weight before eventually adding weight.

Reference: `programs/002-gzclp-linear-progression.md`

## Acceptance Criteria
- [ ] Create goose migration file: `NNNN_seed_gzclp.sql`
- [ ] Create program record:
  - slug: `gzclp`
  - name: `GZCLP`
  - description: Brief description of the tiered progression philosophy
  - author: system user
- [ ] Create 4 workout days with T1/T2 pairings:
  - Day 1: T1 Squat, T2 Bench Press
  - Day 2: T1 Overhead Press, T2 Deadlift
  - Day 3: T1 Bench Press, T2 Squat
  - Day 4: T1 Deadlift, T2 Overhead Press
- [ ] Create T1 prescriptions (default scheme):
  - Stage 1: 5 sets x 3 reps, last set AMRAP (5x3+)
  - Stage 2: 6 sets x 2 reps, last set AMRAP (6x2+) - on failure
  - Stage 3: 10 sets x 1 rep, last set AMRAP (10x1+) - on failure
- [ ] Create T2 prescriptions (default scheme):
  - Stage 1: 3 sets x 10 reps (3x10)
  - Stage 2: 3 sets x 8 reps (3x8) - on failure
  - Stage 3: 3 sets x 6 reps (3x6) - on failure
- [ ] Create progression model:
  - T1/T2 Lower body (Squat, Deadlift): +5 lbs per successful session
  - T1/T2 Upper body (Bench, OHP): +2.5 lbs per successful session
- [ ] Model failure progression (stage changes):
  - T1: 5x3+ -> 6x2+ -> 10x1+ -> retest 5RM
  - T2: 3x10 -> 3x8 -> 3x6 -> reset weight
- [ ] Migration includes proper down migration
- [ ] Migration is idempotent

## Technical Notes
- GZCLP is unique in that failure doesn't reduce weight - it changes rep scheme
- The "+" notation indicates AMRAP (last set)
- Success criteria:
  - T1: total reps >= base volume (e.g., 15 reps for 5x3)
  - T2: complete all prescribed reps
- May need to model "stages" or "phases" for rep scheme progression
- T3 (accessory) work is user-defined and not seeded
- Consider how to store the stage progression rules:
  - Could be separate progression_rules table
  - Could be metadata on prescriptions
  - Could be documented and handled by application logic

## Dependencies
- Blocks: 005 (verification tests need all programs seeded)
- Blocked by: 003 (run in sequence)

## Resources / Links
- Program specification: `programs/002-gzclp-linear-progression.md`
- GZCL method: https://swoleateveryheight.blogspot.com/2016/02/gzcl-applications-adaptations.html

# Ticket 004: GZCLP Program Seed

## Reference
- Ticket: `phases/in-progress/002-frontend-readiness/sprints/todo/004-seed-canonical-programs/tickets/todo/004-gzclp-seed.md`
- Program Spec: `programs/002-gzclp-linear-progression.md`
- ERD: REQ-SEED-004

## Task
Create a goose migration that seeds the GZCLP program.

## Program Structure (4 days/week)

### Day 1
- T1: Squat (5x3+)
- T2: Bench (3x10)

### Day 2
- T1: OHP (5x3+)
- T2: Deadlift (3x10)

### Day 3
- T1: Bench (5x3+)
- T2: Squat (3x10)

### Day 4
- T1: Deadlift (5x3+)
- T2: OHP (3x10)

### T1 Progression (on failure)
- 5x3+ -> 6x2+ -> 10x1+ -> retest

### T2 Progression (on failure)
- 3x10 -> 3x8 -> 3x6 -> reset

### Weight Progression
- Lower body: +5 lbs per success
- Upper body: +2.5 lbs per success

## Acceptance Criteria
- [ ] Create migration `NNNN_seed_gzclp.sql`
- [ ] Create program with slug `gzclp`
- [ ] Create 4 workout days with T1/T2 pairings
- [ ] Create T1 prescriptions (5x3+ default)
- [ ] Create T2 prescriptions (3x10 default)
- [ ] Create progression rules (including stage changes)
- [ ] Include down migration
- [ ] Ensure idempotency

## Technical Notes
- Unique failure progression (rep scheme changes, not weight)
- "+" = AMRAP on last set
- T3 is user-defined, not seeded

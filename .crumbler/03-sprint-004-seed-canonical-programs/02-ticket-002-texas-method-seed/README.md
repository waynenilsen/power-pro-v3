# Ticket 002: Texas Method Program Seed

## Reference
- Ticket: `phases/in-progress/002-frontend-readiness/sprints/todo/004-seed-canonical-programs/tickets/todo/002-texas-method-seed.md`
- Program Spec: `programs/006-texas-method.md`
- ERD: REQ-SEED-002

## Task
Create a goose migration that seeds the Texas Method intermediate program.

## Program Structure (3 days/week)

### Volume Day (Monday)
- Squat: 5x5 @ 90% of Friday weight
- Bench/Press: 5x5 @ 90%
- Deadlift: 1x5

### Recovery Day (Wednesday)
- Squat: 2x5 @ 80% of Monday (~72% of Friday)
- Bench/Press: 3x5 @ 90%

### Intensity Day (Friday)
- Squat: 1x5 @ 100% (PR attempt)
- Bench/Press: 1x5 @ 100%

### Progression
- All lifts: +5 lbs per week
- Friday weight is the anchor

## Acceptance Criteria
- [ ] Create migration `NNNN_seed_texas_method.sql`
- [ ] Create program with slug `texas-method`
- [ ] Create 3 workout days
- [ ] Create percentage-based prescriptions
- [ ] Create weekly progression rules
- [ ] Include down migration
- [ ] Ensure idempotency

## Technical Notes
- Percentage relationships are key
- Weekly (not session) progression
- May need to model Bench/Press alternation

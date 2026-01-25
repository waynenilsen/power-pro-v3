# Ticket 003: Wendler 5/3/1 Program Seed

## Reference
- Ticket: `phases/in-progress/002-frontend-readiness/sprints/todo/004-seed-canonical-programs/tickets/todo/003-531-seed.md`
- Program Spec: `programs/003-wendler-531-bbb.md`
- ERD: REQ-SEED-003

## Task
Create a goose migration that seeds the Wendler 5/3/1 program.

## Program Structure (4 days/week, 4-week cycle)

### Days
- Day 1: Overhead Press
- Day 2: Deadlift
- Day 3: Bench Press
- Day 4: Squat

### Week 1 ("5s Week")
- Set 1: 65% x 5
- Set 2: 75% x 5
- Set 3: 85% x 5+ (AMRAP)

### Week 2 ("3s Week")
- Set 1: 70% x 3
- Set 2: 80% x 3
- Set 3: 90% x 3+ (AMRAP)

### Week 3 ("5/3/1 Week")
- Set 1: 75% x 5
- Set 2: 85% x 3
- Set 3: 95% x 1+ (AMRAP)

### Week 4 ("Deload")
- Set 1: 40% x 5
- Set 2: 50% x 5
- Set 3: 60% x 5

### Progression
- Upper body: +5 lbs per cycle
- Lower body: +10 lbs per cycle

## Acceptance Criteria
- [ ] Create migration `NNNN_seed_531.sql`
- [ ] Create program with slug `531`
- [ ] Create 4 workout days
- [ ] Create 4 weeks per cycle with different prescriptions
- [ ] Mark AMRAP sets appropriately
- [ ] Create cycle progression rules
- [ ] Include down migration
- [ ] Ensure idempotency

## Technical Notes
- Most complex program (4 weeks x 4 days)
- Percentages based on Training Max (90% of 1RM)
- "+" notation = AMRAP

# Ticket 001: Starting Strength Program Seed

## Reference
- Ticket: `phases/in-progress/002-frontend-readiness/sprints/todo/004-seed-canonical-programs/tickets/todo/001-starting-strength-seed.md`
- Program Spec: `programs/004-starting-strength.md`
- ERD: REQ-SEED-001, REQ-LOOKUP-001, REQ-LOOKUP-002

## Task
Create a goose migration that seeds the Starting Strength novice program.

## Program Structure

### Workout A
- Squat: 3x5
- Bench Press: 3x5
- Deadlift: 1x5

### Workout B
- Squat: 3x5
- Overhead Press: 3x5
- Power Clean: 5x3

### Progression
- Squat: +5 lbs per session
- Bench Press: +5 lbs per session
- Overhead Press: +5 lbs per session
- Deadlift: +10 lbs per session
- Power Clean: +5 lbs per session

## Acceptance Criteria
- [ ] Create migration `NNNN_seed_starting_strength.sql`
- [ ] Ensure OHP and Power Clean lift types exist (squat/bench/deadlift already seeded)
- [ ] Create program with slug `starting-strength`
- [ ] Create 2 workout days (A and B)
- [ ] Create prescriptions for each exercise
- [ ] Create progression rules
- [ ] Include down migration
- [ ] Ensure idempotency (INSERT OR IGNORE pattern)

## Technical Notes
- Use deterministic UUIDs for seeded data
- Check existing lifts table - squat/bench/deadlift already exist
- Need to add OHP and Power Clean to lifts if missing
- Use cycle/week/day structure per existing schema

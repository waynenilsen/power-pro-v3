# 002: Texas Method Program Seed

## ERD Reference
Implements: REQ-SEED-002

## Description
Create a goose migration that seeds the Texas Method intermediate program into the database. This program uses weekly periodization with Volume, Recovery, and Intensity days.

## Context / Background
The Texas Method is designed for intermediate lifters who have exhausted linear progression. It operates on a weekly cycle: Monday (Volume Day) accumulates stress, Wednesday (Recovery Day) maintains technique with reduced load, and Friday (Intensity Day) tests strength with new personal records.

Reference: `programs/006-texas-method.md`

## Acceptance Criteria
- [ ] Create goose migration file: `NNNN_seed_texas_method.sql`
- [ ] Create program record:
  - slug: `texas-method`
  - name: `Texas Method`
  - description: Brief description of weekly periodization philosophy
  - author: system user
- [ ] Create Volume Day (Day 1/Monday) prescriptions:
  - Squat: 5 sets x 5 reps @ 90% of Friday weight
  - Bench Press OR Overhead Press: 5 sets x 5 reps @ 90% (alternating weeks)
  - Deadlift: 1 set x 5 reps
- [ ] Create Recovery Day (Day 2/Wednesday) prescriptions:
  - Squat: 2 sets x 5 reps @ 80% of Monday weight (~72% of Friday)
  - Bench Press OR Overhead Press: 3 sets x 5 reps @ 90% of 5RM (alternating)
- [ ] Create Intensity Day (Day 3/Friday) prescriptions:
  - Squat: 1 set x 5 reps @ 100% (PR attempt)
  - Bench Press OR Overhead Press: 1 set x 5 reps @ 100% (alternating)
  - Power Clean: 5 sets x 3 reps (optional)
- [ ] Create progression model:
  - Weekly progression: +5 lbs per week for all lifts
  - Friday weight is the anchor; Monday/Wednesday derive from it
- [ ] Model the Bench/Press alternation:
  - Week A: Bench on Volume/Recovery/Intensity, Press on Recovery only
  - Week B: Press on Volume/Recovery/Intensity, Bench on Recovery only
  - Or simplify: both lifts trained, alternating which gets intensity work
- [ ] Migration includes proper down migration
- [ ] Migration is idempotent

## Technical Notes
- The percentage relationships are key:
  - Monday = 90% of Friday
  - Wednesday = 80% of Monday (or ~72% of Friday)
- For initial seed, can simplify alternation by having both lifts available each day
- Progression is weekly, not session-based (differs from Starting Strength)
- May need to model 2-week cycle to capture Bench/Press alternation properly
- Reference existing schema for how to store percentage-based prescriptions

## Dependencies
- Blocks: 005 (verification tests need all programs seeded)
- Blocked by: 001 (run in sequence to ensure lookup data exists)

## Resources / Links
- Program specification: `programs/006-texas-method.md`

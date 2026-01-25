# 001: Starting Strength Program Seed

## ERD Reference
Implements: REQ-SEED-001, REQ-LOOKUP-001, REQ-LOOKUP-002

## Description
Create a goose migration that seeds the Starting Strength novice program into the database. This includes ensuring required lookup data exists and creating the complete program structure with all days, prescriptions, and progression rules.

## Context / Background
Starting Strength is Mark Rippetoe's classic novice program. It uses an A/B rotation with 3 training days per week. The program focuses on linear progression where weight increases every session. This is typically the first program new lifters should use.

Reference: `programs/004-starting-strength.md`

## Acceptance Criteria
- [ ] Create goose migration file: `NNNN_seed_starting_strength.sql`
- [ ] Ensure required lift type lookups exist (or create if missing):
  - Squat, Bench Press, Deadlift, Overhead Press, Power Clean
- [ ] Ensure required rep scheme lookups exist (or create if missing):
  - 3x5, 1x5, 5x3
- [ ] Create program record:
  - slug: `starting-strength`
  - name: `Starting Strength`
  - description: Brief description of the program philosophy
  - author: system user (admin/seed user)
- [ ] Create Workout A prescriptions:
  - Squat: 3 sets x 5 reps
  - Bench Press: 3 sets x 5 reps
  - Deadlift: 1 set x 5 reps
- [ ] Create Workout B prescriptions:
  - Squat: 3 sets x 5 reps
  - Overhead Press: 3 sets x 5 reps
  - Power Clean: 5 sets x 3 reps
- [ ] Create progression model:
  - Squat: +5 lbs per session
  - Bench Press: +5 lbs per session
  - Overhead Press: +5 lbs per session
  - Deadlift: +10 lbs per session
  - Power Clean: +5 lbs per session
- [ ] Migration includes proper down migration (DELETE statements)
- [ ] Migration is idempotent (uses INSERT OR IGNORE or similar pattern)
- [ ] Run migration locally and verify program appears in database

## Technical Notes
- Use UUID generation compatible with SQLite (may need to use hex() and randomblob())
- Reference existing E2E test setup code for INSERT patterns
- Program author should be a system user - check if one exists or create as part of migration
- Consider using a transaction to ensure atomicity
- A/B rotation can be modeled as 2 program_days with day_number 1 and 2
- The A/B pattern repeats indefinitely; no fixed week/cycle structure needed for basic version
- Progression is per-session (linear), stored in progression table linked to prescription

## Dependencies
- Blocks: 005 (verification tests need all programs seeded)
- Blocked by: None (first migration in sequence)

## Resources / Links
- Program specification: `programs/004-starting-strength.md`
- Existing schema: check `internal/db/migrations/` for table structures
- E2E test patterns: check `e2e/` for INSERT examples

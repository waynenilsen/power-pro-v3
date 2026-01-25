# 002: Backfill Migration for Canonical Programs

## ERD Reference
Implements: REQ-BACKFILL-001

## Description
Create a goose migration that backfills discovery metadata for the four canonical programs seeded in Sprint 004. Each program needs accurate difficulty, days_per_week, focus, and has_amrap values based on program documentation.

## Context / Background
Sprint 004 seeded four canonical programs (Starting Strength, Texas Method, 5/3/1, GZCLP) but without discovery metadata columns. Now that those columns exist from ticket 001, we need to update the canonical programs with accurate metadata values that reflect each program's characteristics.

## Acceptance Criteria
- [ ] Create goose migration file: `NNNN_backfill_canonical_program_metadata.sql`
- [ ] Update Starting Strength:
  - slug: `starting-strength`
  - difficulty: 'beginner'
  - days_per_week: 3
  - focus: 'strength'
  - has_amrap: 0
- [ ] Update Texas Method:
  - slug: `texas-method`
  - difficulty: 'intermediate'
  - days_per_week: 3
  - focus: 'strength'
  - has_amrap: 0
- [ ] Update Wendler 5/3/1:
  - slug: `531`
  - difficulty: 'intermediate'
  - days_per_week: 4
  - focus: 'strength'
  - has_amrap: 1
- [ ] Update GZCLP:
  - slug: `gzclp`
  - difficulty: 'beginner'
  - days_per_week: 4
  - focus: 'strength'
  - has_amrap: 1
- [ ] Migration is idempotent (safe to run multiple times)
- [ ] Down migration resets values to defaults (or is no-op)
- [ ] Run migration locally and verify metadata values

## Technical Notes
- Use UPDATE statements matching on slug (canonical identifier)
- Migration should not fail if programs don't exist (for fresh databases)
- Consider using UPDATE ... WHERE EXISTS pattern for safety
- Reference programs/*.md documentation for accuracy
- All canonical programs are strength-focused (no hypertrophy or peaking programs seeded yet)
- Starting Strength and Texas Method use straight sets (no AMRAP)
- 5/3/1 and GZCLP use AMRAP sets (+) on main lifts

## Dependencies
- Blocks: 006 (tests verify backfill correctness)
- Blocked by: 001 (columns must exist before backfill)

## Resources / Links
- Program documentation: `programs/004-starting-strength.md`, `programs/002-texas-method.md`, `programs/001-531.md`, `programs/003-gzclp.md`
- Schema migration: ticket 001
- ERD requirements: `phases/in-progress/002-frontend-readiness/sprints/todo/006-program-discovery/erd.md`

# 003: ProgramProgression Join Table Schema and Migration

## ERD Reference
Implements: REQ-PROG-003, REQ-CYCLE-003, REQ-MANUAL-002
Related to: ERD-003 (Program entity)

## Description
Create the database schema and migration for the ProgramProgression join table. This entity links progressions to programs with optional lift-specific configuration and enables fine-grained control over progression behavior.

## Context / Background
Different programs use different progressions, and within a program, different lifts may have different progression configurations. For example, 5/3/1 uses +5lb for upper body lifts and +10lb for lower body lifts. This join table enables that flexibility while supporting priority ordering and enable/disable functionality.

## Acceptance Criteria
- [ ] Create `program_progressions` table with the following columns:
  - `id` (UUID, primary key)
  - `program_id` (UUID, required, foreign key to programs.id)
  - `progression_id` (UUID, required, foreign key to progressions.id)
  - `lift_id` (UUID, nullable, foreign key to lifts.id) - null means applies to all lifts
  - `priority` (INTEGER, required, default 0) - order of evaluation
  - `enabled` (BOOLEAN, required, default true) - allows disabling without deletion
  - `override_increment` (DECIMAL, nullable) - lift-specific increment override
  - `created_at` (TIMESTAMP, required)
  - `updated_at` (TIMESTAMP, required)
- [ ] Create unique constraint on (`program_id`, `progression_id`, `lift_id`) - can't have same progression twice for same lift
- [ ] Create index on `program_id` for program lookup
- [ ] Create index on (`program_id`, `lift_id`) for lift-specific progression lookup
- [ ] Create goose migration file with proper up/down migrations
- [ ] Foreign key constraints with appropriate ON DELETE CASCADE behavior

## Technical Notes
- When `lift_id` is NULL, the progression applies as a program-level default
- Lift-specific entries (non-null `lift_id`) take precedence over program-level defaults
- Priority field allows ordering when multiple progressions could apply
- `enabled = false` supports Sheiko-style programs with no auto-progression
- `override_increment` allows 5/3/1's lift-category-specific increments

## Dependencies
- Blocks: 010 (Program Progression Configuration API depends on this)
- Blocked by: 001 (References progressions table), ERD-003 (References programs table)
- Related: 005, 006 (LinearProgression and CycleProgression use this for configuration)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/004-progression-rules/erd.md
- Tech Stack: prompts/tech-stack.md

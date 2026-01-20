# 003: Cycle Entity Schema and Migration

## ERD Reference
Implements: REQ-CYCLE-001, REQ-CYCLE-002, REQ-CYCLE-004

## Description
Create the database schema for the Cycle entity, including the goose migration file. The Cycle entity represents the repeating unit of a program (e.g., 4 weeks for 5/3/1, 3 weeks for Greg Nuckols).

## Context / Background
A Cycle is the repeating unit of a program with configurable length. Different programs have different cycle lengths: 5/3/1 uses 4-week cycles, Greg Nuckols uses 3-week cycles, Starting Strength uses 1-week cycles.

## Acceptance Criteria
- [ ] Cycle table created with: id (UUID), name (VARCHAR 100, NOT NULL), length_weeks (INTEGER >= 1, NOT NULL)
- [ ] Proper constraints on length_weeks
- [ ] Goose migration file created
- [ ] Migration tested (up and down)

## Technical Notes
- Cycle table: `cycles` with columns: id, name, length_weeks, created_at, updated_at
- The weeks relationship (REQ-CYCLE-003) is handled by the weeks table's cycle_id foreign key
- Validation: count of weeks in cycle must equal length_weeks (enforced at application layer)

## Dependencies
- Blocks: 004, 009, 010
- Blocked by: 002 (Week schema required for proper FK setup, though Week references Cycle)
- Related: REQ-CYCLE-003 (Cycle Weeks relationship)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/003-schedule-periodization/erd.md

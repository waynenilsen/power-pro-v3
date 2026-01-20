# 006: User Program State Schema and Migration

## ERD Reference
Implements: REQ-STATE-001

## Description
Create the database schema for the UserProgramState entity, including the goose migration file. This entity tracks the user's current position within a program (current week, cycle iteration).

## Context / Background
Users need to track their position in a program to know what workout to do "today" and to trigger progressions at cycle completion. The state tracks current week within cycle and which cycle iteration the user is on.

## Acceptance Criteria
- [ ] UserProgramState table created with: id (UUID), user_id (UUID FK), program_id (UUID FK), current_week (INTEGER >= 1), current_cycle_iteration (INTEGER >= 1), current_day_index (INTEGER, nullable)
- [ ] Unique constraint on user_id (one program per user at a time)
- [ ] Proper foreign key constraints to users and programs tables
- [ ] Goose migration file created
- [ ] Migration tested (up and down)

## Technical Notes
- UserProgramState table: `user_program_states` with columns: id, user_id, program_id, current_week, current_cycle_iteration, current_day_index, enrolled_at, updated_at
- current_week: 1 to cycle.length_weeks
- current_cycle_iteration: which time through the cycle (1, 2, 3...)
- current_day_index: optional, for tracking day position within week
- enrolled_at: timestamp when user enrolled in program
- User can only be enrolled in one program at a time (unique constraint on user_id)

## Dependencies
- Blocks: 013, 014, 015
- Blocked by: 005 (Program schema required)
- Related: REQ-STATE-002, REQ-STATE-003, REQ-PROG-003

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/003-schedule-periodization/erd.md

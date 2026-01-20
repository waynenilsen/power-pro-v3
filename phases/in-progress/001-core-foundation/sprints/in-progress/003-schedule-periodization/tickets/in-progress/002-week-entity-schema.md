# 002: Week Entity Schema and Migration

## ERD Reference
Implements: REQ-WEEK-001, REQ-WEEK-002, REQ-WEEK-004

## Description
Create the database schema for the Week entity and WeekDay join table, including the goose migration file. The Week entity represents a collection of training days within a program cycle.

## Context / Background
A Week is a collection of training days. Most programs operate on weekly cycles. Weeks are numbered within cycles and can support A/B week rotation for alternating programs.

## Acceptance Criteria
- [ ] Week table created with: id (UUID), week_number (INTEGER >= 1), variant (VARCHAR, nullable - for A/B support), cycle_id (UUID FK)
- [ ] WeekDay join table created with: week_id, day_id, day_of_week (ENUM or INTEGER 1-7)
- [ ] Week number unique within cycle constraint
- [ ] Proper foreign key constraints
- [ ] Goose migration file created
- [ ] Migration tested (up and down)

## Technical Notes
- Week table: `weeks` with columns: id, week_number, variant (nullable), cycle_id, created_at, updated_at
- WeekDay join table: `week_days` with columns: id, week_id, day_id, day_of_week, created_at
- day_of_week: Use TEXT with values MONDAY-SUNDAY or INTEGER 1-7
- Multiple days can map to same day-of-week (for multiple sessions per day)
- Variant field (A, B, null) supports A/B week rotation patterns

## Dependencies
- Blocks: 003, 008, 009
- Blocked by: 001 (Day schema required)
- Related: REQ-WEEK-003 (Week Days relationship)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/003-schedule-periodization/erd.md

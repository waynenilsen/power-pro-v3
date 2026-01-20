# 004: Lookup Tables Schema and Migration

## ERD Reference
Implements: REQ-LOOKUP-001, REQ-LOOKUP-002, REQ-LOOKUP-003, REQ-LOOKUP-004

## Description
Create the database schema for WeeklyLookup and DailyLookup entities, including the goose migration file. These lookup tables enable programs to vary percentages and parameters based on week number or day identifier.

## Context / Background
Programs like 5/3/1 vary percentages by week (Week 1 = 65/75/85%, Week 2 = 70/80/90%). Bill Starr varies intensity by day (Heavy/Light/Medium). Lookup tables enable this variation without duplicating prescriptions.

## Acceptance Criteria
- [ ] WeeklyLookup table created with: id (UUID), name (VARCHAR, NOT NULL), entries (JSONB)
- [ ] DailyLookup table created with: id (UUID), name (VARCHAR, NOT NULL), entries (JSONB)
- [ ] WeeklyLookup entries structure: array of {weekNumber, percentages[], reps[], percentageModifier}
- [ ] DailyLookup entries structure: map of {dayIdentifier -> percentageModifier, intensityLevel}
- [ ] Goose migration file created
- [ ] Migration tested (up and down)

## Technical Notes
- WeeklyLookup table: `weekly_lookups` with columns: id, name, entries (JSONB), program_id (optional FK), created_at, updated_at
- DailyLookup table: `daily_lookups` with columns: id, name, entries (JSONB), program_id (optional FK), created_at, updated_at
- JSONB used for flexibility in parameter structure
- Example WeeklyLookup entries:
  ```json
  [
    { "weekNumber": 1, "percentages": [65, 75, 85], "reps": [5, 5, 5] },
    { "weekNumber": 2, "percentages": [70, 80, 90], "reps": [3, 3, 3] }
  ]
  ```
- Example DailyLookup entries:
  ```json
  {
    "heavy": { "percentageModifier": 1.0, "intensityLevel": "HEAVY" },
    "light": { "percentageModifier": 0.7, "intensityLevel": "LIGHT" },
    "medium": { "percentageModifier": 0.8, "intensityLevel": "MEDIUM" }
  }
  ```

## Dependencies
- Blocks: 010, 011, 014
- Blocked by: None (can be created independently)
- Related: ERD-002 LoadStrategy (for integration)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/003-schedule-periodization/erd.md

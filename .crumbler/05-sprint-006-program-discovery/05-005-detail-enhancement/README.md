# 005: Program Detail Enhancement

## Task
Enhance GET /programs/{id} to include sampleWeek, liftRequirements, and estimatedSessionMinutes.

## New Response Fields

### sampleWeek
Array of day objects showing first week structure:
```json
[
  {"day": 1, "name": "Workout A", "exerciseCount": 3},
  {"day": 2, "name": "Workout B", "exerciseCount": 3}
]
```
- Query program_days and count prescriptions per day
- For A/B rotation programs, show both A and B days

### liftRequirements
Array of unique lift names, sorted alphabetically:
```json
["Bench Press", "Deadlift", "Overhead Press", "Squat"]
```
- Query: SELECT DISTINCT lift_name FROM prescriptions WHERE program_id = ?

### estimatedSessionMinutes
Integer calculated as:
- (total sets per average day * 3 minutes) + (exercises per day * 2 minutes warmup)

## Changes Required

1. **DTO** - Add fields to ProgramDetailResponse
2. **Service/Repository** - Query program_days, prescriptions, lifts tables
3. **Handler** - Assemble enhanced response

## Reference
- Ticket: `phases/in-progress/002-frontend-readiness/sprints/todo/006-program-discovery/tickets/todo/005-detail-enhancement.md`
- Depends on: 003-filtering-implementation (domain model updates)

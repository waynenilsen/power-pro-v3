# 006: E2E Tests for Program Discovery

## Task
Create comprehensive end-to-end tests for all program discovery features.

## Test File
`e2e/program_discovery_test.go`

## Test Cases

### Filtering Tests
- Filter by difficulty (beginner, intermediate, advanced, invalid)
- Filter by days_per_week (3, 4, 5, invalid values)
- Filter by focus (strength, hypertrophy, invalid)
- Filter by has_amrap (true, false, invalid)
- Combined filters (AND logic)
- Empty result returns 200 with empty array

### Search Tests
- search=strength returns Starting Strength
- search=531 returns 5/3/1
- search=gz returns GZCLP
- Case-insensitive search
- Search + filter combinations

### Detail Enhancement Tests
- Response includes sampleWeek array
- Response includes liftRequirements array (sorted)
- Response includes estimatedSessionMinutes
- Verify structure for each canonical program

### List Response Tests
- List includes difficulty, daysPerWeek, focus, hasAmrap
- Values match backfilled data

## Reference
- Ticket: `phases/in-progress/002-frontend-readiness/sprints/todo/006-program-discovery/tickets/todo/006-e2e-tests.md`
- Depends on: All previous tickets (001-005)

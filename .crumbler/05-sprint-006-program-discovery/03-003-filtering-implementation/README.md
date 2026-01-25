# 003: Program Filtering Query Implementation

## Task
Implement filtering for GET /programs endpoint by difficulty, days_per_week, focus, and has_amrap.

## Query Parameters
- `?difficulty=beginner` - beginner, intermediate, advanced
- `?days_per_week=3` - 1-7
- `?focus=strength` - strength, hypertrophy, peaking
- `?has_amrap=true` - true, false

## Changes Required

1. **Domain Model** (`internal/domain/program/`)
   - Add filter options struct
   - Add new fields to Program model

2. **Repository** (`internal/repository/program_repository.go`)
   - Update List() to accept filter options
   - Build WHERE clause dynamically
   - Use parameterized queries

3. **Handler** (`internal/api/program_handler.go`)
   - Parse query parameters
   - Validate filter values (400 for invalid)
   - Pass filters to service/repository

4. **Response DTO**
   - Add new fields: difficulty, daysPerWeek, focus, hasAmrap

## Behavior
- Combined filters use AND logic
- Empty result returns 200 with empty array (not 404)
- Invalid filter values return 400 Bad Request

## Reference
- Ticket: `phases/in-progress/002-frontend-readiness/sprints/todo/006-program-discovery/tickets/todo/003-filtering-implementation.md`
- Depends on: 001-schema-migration

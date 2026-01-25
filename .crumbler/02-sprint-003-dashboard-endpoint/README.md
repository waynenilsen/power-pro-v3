# Sprint 003: Dashboard Endpoint

## Overview

Implement an aggregated dashboard endpoint that combines multiple data sources into a single API response. This reduces frontend complexity by providing enrollment status, next workout, current session, recent workouts, and current maxes in one call.

## Reference Documents

- PRD: `phases/in-progress/002-frontend-readiness/sprints/todo/003-dashboard-endpoint/prd.md`
- ERD: `phases/in-progress/002-frontend-readiness/sprints/todo/003-dashboard-endpoint/erd.md`
- Tickets: `phases/in-progress/002-frontend-readiness/sprints/todo/003-dashboard-endpoint/tickets/todo/`

## Requirements Summary

### Endpoint
- `GET /users/{id}/dashboard` - Returns aggregated dashboard data

### Response Sections
1. **enrollment** - Current program status (status, programName, cycleIteration, cycleStatus, weekNumber, weekStatus)
2. **nextWorkout** - Upcoming workout preview (dayName, daySlug, exerciseCount, estimatedSets)
3. **currentSession** - In-progress workout if any (sessionId, dayName, startedAt, setsCompleted, totalSets)
4. **recentWorkouts** - Last 5 completed workouts (date, dayName, setsCompleted)
5. **currentMaxes** - Training maxes for all lifts (lift, value, type)

### Authorization
- Owner-only access (not even admins can view other users' dashboards)

### Design Principles
- No new database tables - aggregation only
- Empty data returns null/empty array, not errors
- All sections always present in response

## Tickets (in order)

1. `001-dashboard-service.md` - Dashboard aggregation service structure
2. `002-enrollment-aggregation.md` - Enrollment status section
3. `003-next-workout-calculation.md` - Next workout calculation
4. `004-current-session-query.md` - Active session detection
5. `005-recent-workouts-query.md` - Recent workout history
6. `006-current-maxes-query.md` - Current training maxes
7. `007-dashboard-api-endpoint.md` - API handler and routing
8. `008-dashboard-e2e-tests.md` - End-to-end tests

## Dependencies

- Sprint 001 (Authentication) - complete
- Sprint 002 (User Profile) - for weight_unit preference

## Success Criteria

- Single endpoint returns all dashboard sections
- Response < 200ms p95
- Authorization enforced
- E2E tests pass

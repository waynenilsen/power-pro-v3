# Sprint 003 - Dashboard Endpoint

## Overview

Create a single aggregated endpoint for the home screen that combines multiple queries into one response.

## Endpoint

```
GET /users/{id}/dashboard
```

## Response Shape

```json
{
  "enrollment": {
    "status": "ACTIVE",
    "program_name": "Texas Method",
    "cycle_iteration": 2,
    "cycle_status": "IN_PROGRESS",
    "week_number": 1,
    "week_status": "IN_PROGRESS"
  },
  "next_workout": {
    "day_name": "Volume Day",
    "day_slug": "volume",
    "exercise_count": 2,
    "estimated_sets": 10
  },
  "current_session": null,
  "recent_workouts": [
    {"date": "2024-01-15", "day_name": "Intensity Day", "sets_completed": 4}
  ],
  "current_maxes": [
    {"lift": "Squat", "value": 320, "type": "TRAINING_MAX"},
    {"lift": "Bench", "value": 227.5, "type": "TRAINING_MAX"}
  ]
}
```

## Design Notes

- This is an aggregation endpoint - no new data, just combines existing queries
- Should handle users with no enrollment gracefully
- Should handle users mid-workout (current_session populated)

## Tasks

1. Create sprint directory: `phases/in-progress/002-frontend-readiness/sprints/todo/003-dashboard-endpoint/`
2. Write `prd.md` - Product requirements for dashboard
3. Write `erd.md` - Engineering requirements with detailed specs
4. Create ticket directory structure
5. Create tickets (at least 5):
   - Dashboard service aggregation logic
   - Enrollment status aggregation
   - Next workout calculation
   - Recent workouts query
   - Current maxes query
   - Dashboard API endpoint
   - E2E tests for dashboard

## Acceptance Criteria

- [ ] Sprint directory structure exists
- [ ] PRD document completed
- [ ] ERD document completed with REQ-XXX requirements
- [ ] At least 5 tickets created in `tickets/todo/`
- [ ] Tickets reference ERD requirements

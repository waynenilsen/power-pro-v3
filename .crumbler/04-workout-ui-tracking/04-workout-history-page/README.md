# Workout History Page

Implement the workout history list and detail view.

## Location

`frontend/src/routes/pages/History.tsx` (replace existing stub)

## Features

### 1. History List View

Display paginated list of past workout sessions:
- Use useWorkoutSessions hook
- Show: date, week/day info, duration (if available), status
- Visual distinction for COMPLETED vs ABANDONED sessions
- Click to navigate to detail view

### 2. Workout Session Card

For each session in the list:
```
[Status Icon] Week X, Day Y
Date: Jan 26, 2026
Duration: 45 mins (if finishedAt exists)
Status: Completed / Abandoned
```

### 3. Detail View (Optional Enhancement)

Could add `/history/:sessionId` route for detailed view:
- Full workout breakdown
- Would need to store actual workout data with session
- For now, just show the summary info

## Hooks Used

- useWorkoutSessions(userId, { page, pageSize, status })
- useWorkoutSession(userId, sessionId) - for detail view if implemented

## Pagination

- Default page size: 10
- Load more or pagination controls
- Filter by status (All / Completed / Abandoned)

## Empty States

- "No workouts yet" message when list is empty
- Link to start a workout if enrolled

## API Response

```typescript
interface WorkoutSession {
  id: string;
  userProgramStateId: string;
  weekNumber: number;
  dayIndex: number;
  status: 'IN_PROGRESS' | 'COMPLETED' | 'ABANDONED';
  startedAt: string;
  finishedAt?: string;
}
```

## UI Notes

- Use consistent card styling with rest of app
- Status badges with appropriate colors
- Date formatting (relative or absolute)
- Use /frontend-design skill

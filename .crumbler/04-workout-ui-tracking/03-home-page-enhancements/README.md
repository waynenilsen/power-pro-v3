# Home Page Enhancements

Enhance the Home page to handle additional workout states.

## Location

`frontend/src/routes/pages/Home.tsx` (modify existing)

## Current Implementation

Home.tsx already shows:
- Enrolled program card with progress
- "Start Workout" button linking to /workout
- Browse programs CTA when not enrolled

## Enhancements Needed

### 1. In-Progress Workout Detection

Check `enrollment.data.currentWorkoutSession` for active sessions:

```typescript
const hasActiveSession = enrollment?.data?.currentWorkoutSession?.status === 'IN_PROGRESS';
```

If active session exists:
- Change button from "Start Workout" to "Continue Workout"
- Add visual indicator that workout is in progress
- Show session start time or elapsed time

### 2. BETWEEN_CYCLES State Handling

Check `enrollment.data.enrollmentStatus`:

```typescript
const isBetweenCycles = enrollment?.data?.enrollmentStatus === 'BETWEEN_CYCLES';
```

If BETWEEN_CYCLES:
- Hide "Start Workout" button
- Show "Cycle Complete!" message
- Show "Start Next Cycle" button
- Call advanceState with advanceType: 'week' to start next cycle

### 3. Current Day Display

Add display of current day in the cycle:
- "Week X, Day Y" indicator
- Or day name from the program if available

## API Types Reference

```typescript
interface EnrollmentResponse {
  data: Enrollment & {
    enrollmentStatus: 'ACTIVE' | 'BETWEEN_CYCLES' | 'QUIT';
    cycleStatus: 'PENDING' | 'IN_PROGRESS' | 'COMPLETED';
    weekStatus: 'PENDING' | 'IN_PROGRESS' | 'COMPLETED';
    currentWorkoutSession?: WorkoutSessionSummary;
  };
}
```

## Implementation Notes

- Keep existing component structure
- Add new conditional rendering in EnrolledProgramCard
- May need useAdvanceState mutation hook (add if not exists)
- Use /frontend-design skill for new UI elements

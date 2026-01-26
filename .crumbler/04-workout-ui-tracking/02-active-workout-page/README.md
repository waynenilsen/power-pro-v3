# Active Workout Page

Implement the full workout execution page where users track their workout in real-time.

## Location

`frontend/src/routes/pages/Workout.tsx` (replace existing stub)

## Dependencies

- ExerciseCard component (must be completed first)
- Hooks: useCurrentWorkout, useCurrentWorkoutSession, useEnrollment
- API: startWorkoutSession, finishWorkoutSession, abandonWorkoutSession, advanceState
- ConfirmDialog component (existing)

## Page States

1. **Loading** - Fetching enrollment/workout data
2. **No Enrollment** - User not enrolled, redirect or show message
3. **No Active Session** - Show workout preview with "Start Workout" button
4. **Active Session** - Show exercises with tracking
5. **Session Complete** - Show completion summary and advancement options

## Features

### Workout Preview (Pre-Start)
- Display all exercises from getCurrentWorkout
- "Start Workout" button to begin session
- Cancel/back navigation

### Active Workout View
- List all exercises using ExerciseCard components
- Local state for completed sets tracking
- "Finish Workout" button (enabled when workout started)
- "Abandon Workout" button with confirmation dialog

### Session Completion Flow
1. User clicks "Finish Workout"
2. Call finishWorkoutSession API
3. Show completion summary
4. If not last day of week: Option to "Go Home"
5. If last day of week: Prompt to advance week
6. If last week of cycle: Prompt to start next cycle

## Mutations Needed

Add these to hooks/useWorkouts.ts:
- useStartWorkoutSession
- useFinishWorkoutSession
- useAbandonWorkoutSession
- useAdvanceState

## Local State

```typescript
// Track which sets are completed locally (doesn't persist to server currently)
const [completedSets, setCompletedSets] = useState<Set<string>>(new Set());
// Format: "exerciseIndex-setNumber" e.g., "0-1", "0-2", "1-1"
```

## UI/UX Notes

- Large touch targets for mobile
- Minimal UI during active workout
- Clear visual feedback for completed sets
- Confirmation before abandoning
- Use /frontend-design skill

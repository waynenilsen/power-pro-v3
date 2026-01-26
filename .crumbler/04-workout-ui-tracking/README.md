# Workout UI & Tracking

The core workout experience - viewing and tracking workouts.

## API Endpoints Used
- `GET /users/{userId}/workout` - Generate current workout
- `GET /users/{userId}/workout/preview` - Preview a specific week/day
- `POST /workouts/start` - Start a workout session
- `POST /workouts/{id}/finish` - Complete workout session
- `POST /workouts/{id}/abandon` - Abandon workout session
- `GET /users/{userId}/workouts/current` - Get in-progress workout
- `GET /users/{userId}/workouts` - List workout history
- `GET /users/{userId}/enrollment/advance-week` - Advance to next week
- `GET /users/{userId}/enrollment/next-cycle` - Start next cycle

## Tasks

1. Workout Dashboard (Home):
   - Show current program & enrollment state
   - Show current week/day
   - "Start Workout" button
   - If in-progress workout exists, show "Continue Workout"
   - If BETWEEN_CYCLES state, show "Start Next Cycle" option

2. Active Workout View:
   - List all exercises with sets/reps/weight
   - Checkable sets to mark completion
   - Rest timer between sets (optional enhancement)
   - "Finish Workout" button
   - "Abandon Workout" button (with confirmation)

3. Exercise Card Component:
   - Exercise name
   - Weight (in lbs or kg based on user preference if available)
   - Sets and reps breakdown
   - Notes if any
   - Checkboxes for set completion

4. Workout History List (`/history`):
   - Paginated list of past workouts
   - Show: date, program day, duration, status (completed/abandoned)
   - Click to view details

5. Workout Detail View:
   - Full workout breakdown
   - What was prescribed vs what actually happened (if tracking)

6. State Advancement:
   - After finishing workout, prompt about advancing week
   - At end of cycle, prompt about starting next cycle
   - Handle BETWEEN_CYCLES state

## UI Design Notes
- Workout view is the most used screen - make it excellent
- Large touch targets for mobile
- Easy to mark sets done during workout
- Minimal distractions during workout
- Use /frontend-design skill

# 003: Next Workout Calculation

## ERD Reference
Implements: REQ-DASH-004

## Description
Implement the next workout calculation that determines which workout day the user should do next and returns a preview of that workout.

## Context / Background
Powerlifting programs have structured workout days (e.g., Volume Day, Intensity Day, Recovery Day). Given a user's current position in their program, we need to calculate which day comes next and provide summary information about that workout.

## Acceptance Criteria
- [ ] Implement `CalculateNextWorkout(ctx context.Context, userID uuid.UUID) (*NextWorkoutPreview, error)`
- [ ] Return nil if user has no active enrollment
- [ ] Return nil if user has completed all workouts in program
- [ ] NextWorkoutPreview contains: DayName, DaySlug, ExerciseCount, EstimatedSets
- [ ] DayName is human-readable (e.g., "Volume Day")
- [ ] DaySlug is URL-friendly (e.g., "volume-day")
- [ ] ExerciseCount is number of distinct exercises in the workout
- [ ] EstimatedSets is total sets across all exercises
- [ ] If user has active session, next workout is null (they're already working out)
- [ ] Unit tests covering: next day available, cycle complete, program complete, mid-workout

## Technical Notes
- Need to query program template to get day structure
- Consider: what defines "next"? Assume days are sequential within a week
- Program templates define exercises and set counts per day
- May need to look at completed sessions to determine last completed day

## Dependencies
- Blocks: 001 (service needs this calculation)
- Blocked by: Existing program template structure, enrollment domain

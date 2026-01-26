# Current Program Indicator

Show which program the user is enrolled in throughout the app.

## Requirements

### 1. Update Home Page
Update `frontend/src/routes/pages/Home.tsx`:
- Fetch enrollment using `useEnrollment(userId)`
- If enrolled: Show current program card with:
  - Program name
  - Current week/cycle progress
  - Link to "/workout" to start training
  - Link to program details
- If NOT enrolled: Show prompt to browse programs with link to "/programs"

### 2. Add Programs to Navigation
Update `SideNav.tsx` and `BottomNav.tsx`:
- Add "Programs" nav item linking to `/programs`
- Icon suggestion: `BookOpen` or `Library` from lucide-react
- Position: After Home, before Workout

### 3. Optional: Current Program Badge in Nav
If user is enrolled, optionally show small indicator in nav (e.g., dot or small text showing program name).

## Data Available (EnrollmentResponse)
```typescript
{
  data: {
    id: string;
    userId: string;
    program: {
      id: string;
      name: string;
      slug: string;
      description?: string;
      cycleLengthWeeks: number;
    };
    state: {
      currentWeek: number;
      currentCycleIteration: number;
      currentDayIndex?: number;
    };
    enrollmentStatus: 'ACTIVE' | 'BETWEEN_CYCLES' | 'QUIT';
    cycleStatus: 'PENDING' | 'IN_PROGRESS' | 'COMPLETED';
    weekStatus: 'PENDING' | 'IN_PROGRESS' | 'COMPLETED';
    currentWorkoutSession?: WorkoutSessionSummary;
    enrolledAt: string;
    updatedAt: string;
  }
}
```

## Design Notes
- Home page should feel informative but not cluttered
- Progress display should be encouraging
- Clear CTA to start workout or browse programs
- Use existing design patterns from Login/Profile pages

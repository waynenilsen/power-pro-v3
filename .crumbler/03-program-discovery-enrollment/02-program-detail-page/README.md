# Program Detail Page

Build the program detail page at `/programs/:id` route.

## Current State
- Skeleton page exists at `frontend/src/routes/pages/ProgramDetails.tsx`
- Uses `useParams()` to get program id
- API endpoint `getProgram(id)` available returning `ProgramDetail`
- No hook exists yet for fetching single program

## Requirements

1. **Create `useProgram(id)` hook** in `frontend/src/hooks/usePrograms.ts`
   - Query key: `queryKeys.programs.detail(id)`
   - Fetch using `programs.getProgram(id)`
   - Enable only when id is defined

2. **Fetch and display program details**:
   - Program name (large heading)
   - Description
   - Cycle info: name, length in weeks
   - Week breakdown (show week numbers from cycle.weeks)

3. **Handle states**:
   - Loading state
   - Not found / error state
   - Back button to programs list

4. **Enroll button**:
   - Show "Enroll in Program" button
   - If already enrolled in THIS program, show "Currently Enrolled" indicator
   - If enrolled in DIFFERENT program, show "Switch to this Program" button
   - Actual enrollment logic in separate crumb

## Data Available (ProgramDetail)
```typescript
{
  id: string;
  name: string;
  slug: string;
  description?: string;
  cycle: {
    id: string;
    name: string;
    lengthWeeks: number;
    weeks: Array<{ id: string; weekNumber: number }>;
  };
  weeklyLookup?: { id: string; name: string };
  dailyLookup?: { id: string; name: string };
  defaultRounding?: number;
  createdAt: string;
  updatedAt: string;
}
```

## Design Notes
- Use `/frontend-design` skill for polished UI
- Clean layout with clear hierarchy
- Make the enroll CTA prominent
- Show week structure in a visual way

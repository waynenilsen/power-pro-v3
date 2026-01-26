# Programs List Page

Build the programs list page at `/programs` route.

## Current State
- Skeleton page exists at `frontend/src/routes/pages/Programs.tsx`
- `usePrograms` hook exists in `frontend/src/hooks/usePrograms.ts`
- API endpoint `listPrograms` available with pagination support
- Navigation to /programs NOT in nav (need to add route or link from home)

## Requirements

1. **Fetch and display programs** using `usePrograms()` hook
2. **Show program cards/list items** with:
   - Program name
   - Description (if available)
   - Cycle length in weeks (from `ProgramListItem` - need to fetch detail or add to list)
   - Simple clean presentation

3. **Handle states**:
   - Loading state (skeleton or spinner)
   - Empty state (no programs available)
   - Error state

4. **Navigation**:
   - Click on program navigates to `/programs/:id`
   - Add "Browse Programs" link to home page or navigation

## Data Available (ProgramListItem)
```typescript
{
  id: string;
  name: string;
  slug: string;
  description?: string;
  cycleId: string;
  weeklyLookupId?: string;
  dailyLookupId?: string;
  defaultRounding?: number;
  createdAt: string;
  updatedAt: string;
}
```

## Design Notes
- Use existing Tailwind patterns from other pages (dark theme, accent color)
- Mobile-first with clean desktop view
- Use `/frontend-design` skill for polished UI
- Keep it simple - list items can work well, don't overdo cards

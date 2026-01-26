# Set Up TanStack Query

Configure TanStack Query for data fetching and caching.

## Tasks

1. Install TanStack Query:
   ```bash
   bun add @tanstack/react-query
   bun add -D @tanstack/react-query-devtools
   ```

2. Configure QueryClient in `src/lib/query.ts`:
   - Default stale time (5 minutes suggested)
   - Default cache time
   - Retry configuration
   - Error handling defaults

3. Wrap app with QueryClientProvider in main.tsx

4. Add React Query Devtools (dev only)

5. Create example hooks in `src/hooks/`:
   ```
   src/hooks/
   ├── usePrograms.ts      # List programs
   ├── useProgram.ts       # Single program
   ├── useWorkouts.ts      # Workout history
   └── useCurrentUser.ts   # User profile
   ```

6. Integrate hooks with API client from previous crumb

## Query Key Convention

Use consistent query key structure:
```typescript
['programs']           // list
['programs', id]       // single
['workouts', { userId }] // filtered list
```

## Acceptance Criteria

- QueryClient configured with sensible defaults
- DevTools available in development
- Example hooks work with API client
- Proper loading/error states handled

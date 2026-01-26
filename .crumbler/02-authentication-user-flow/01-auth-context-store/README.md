# Auth Context Store

Create a React context for authentication state management with localStorage persistence.

## Requirements

1. Create `/frontend/src/contexts/AuthContext.tsx`:
   - `AuthContext` with provider and hook
   - State: `userId: string | null`, `isAuthenticated: boolean`, `isLoading: boolean`
   - Actions: `login(userId: string)`, `logout()`, `createUser(): string`
   - Initialize from localStorage on mount

2. localStorage integration:
   - Key: `powerpro_user_id`
   - Load on app startup
   - Persist on login
   - Clear on logout

3. Auto-configure API client:
   - When userId changes, call `configureClient({ userId })`
   - This integrates with existing `/frontend/src/api/client.ts`

4. Export from `/frontend/src/contexts/index.ts`

## Usage Pattern
```tsx
const { userId, isAuthenticated, login, logout, createUser } = useAuth();
```

## Files to Create/Modify
- Create: `/frontend/src/contexts/AuthContext.tsx`
- Create: `/frontend/src/contexts/index.ts`
- Modify: `/frontend/src/main.tsx` - wrap app with AuthProvider

## Acceptance Criteria
- Auth state persists across page refreshes
- API client automatically configured with userId
- Loading state while checking localStorage
- TypeScript types exported

# Protected Routes

Implement route protection for authenticated-only pages.

## Requirements

1. Create `/frontend/src/components/auth/ProtectedRoute.tsx`:
   - Wrapper component that checks `isAuthenticated` from AuthContext
   - Shows loading spinner while auth state initializing
   - Redirects to `/login` if not authenticated
   - Renders children (Outlet) if authenticated

2. Update routing in `/frontend/src/routes/index.tsx`:
   - Wrap authenticated routes with ProtectedRoute
   - Keep `/login` public
   - Protected routes: `/`, `/programs/*`, `/workout`, `/history`, `/profile`

3. Handle redirect return:
   - Store intended destination before redirect
   - After login, return to intended destination (optional enhancement)

## Component Pattern
```tsx
<Route element={<ProtectedRoute />}>
  <Route element={<Layout />}>
    <Route index element={<Home />} />
    ...
  </Route>
</Route>
```

## Files to Create/Modify
- Create: `/frontend/src/components/auth/ProtectedRoute.tsx`
- Create: `/frontend/src/components/auth/index.ts`
- Modify: `/frontend/src/routes/index.tsx`

## Acceptance Criteria
- Unauthenticated users redirected to /login
- Loading state shown during auth check
- Authenticated users can access protected routes
- Login page accessible without auth

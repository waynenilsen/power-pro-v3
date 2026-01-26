# Logout Flow

Implement complete logout functionality.

## Requirements

1. Logout action in AuthContext (should exist from auth-context-store):
   - Clear localStorage
   - Reset userId to null
   - Call `configureClient({ userId: undefined })`

2. Add logout button to Profile page:
   - Update `/frontend/src/routes/pages/Profile.tsx`
   - Add "Sign Out" button
   - Confirm dialog (optional): "Are you sure you want to sign out?"
   - On confirm: call `logout()`, redirect to `/login`

3. Add logout to navigation:
   - Consider adding to SideNav or Profile dropdown
   - Mobile: could be in BottomNav profile section or Profile page only

4. Clear TanStack Query cache on logout:
   - Call `queryClient.clear()` to remove cached user data
   - Prevents data leakage between users

## Files to Modify
- `/frontend/src/routes/pages/Profile.tsx` - add logout button
- `/frontend/src/contexts/AuthContext.tsx` - ensure query cache cleared
- Optionally: `/frontend/src/components/layout/SideNav.tsx`

## Acceptance Criteria
- Logout clears all user data from localStorage
- Logout clears TanStack Query cache
- User redirected to login page
- API client reset (no stale userId in headers)
- No cached data visible after re-login as different user

# Auth Header Interceptor

Ensure API client properly handles auth headers and error responses.

## Requirements

1. Verify `/frontend/src/api/client.ts` integration:
   - Confirm `X-User-ID` header sent on all requests when configured
   - Ensure `configureClient()` is called by AuthContext on userId change

2. Add 401/403 response handling:
   - On 401 Unauthorized: Clear auth state, redirect to login
   - On 403 Forbidden: Show appropriate error (user may need different permissions)

3. Create error handling hook or utility:
   - `useApiErrorHandler()` hook or integrate into TanStack Query setup
   - Automatically handle auth errors globally

4. Update TanStack Query config in `/frontend/src/lib/query.ts`:
   - Add global `onError` handler for auth errors
   - Trigger logout on 401 responses

## Integration Points
- AuthContext logout function
- API client error handling
- TanStack Query global error handling

## Files to Modify
- `/frontend/src/api/client.ts` - verify/enhance error handling
- `/frontend/src/lib/query.ts` - add global auth error handling
- May need to create utility for handling auth errors

## Acceptance Criteria
- All API requests include X-User-ID when authenticated
- 401 responses trigger logout and redirect
- 403 responses handled gracefully
- Error handling doesn't break existing functionality

# Authentication & User Flow

Implement authentication and user management UI.

## Current Auth Model
The API uses simple header-based auth:
- `X-User-ID: {uuid}` - User ID in header
- `Authorization: Bearer {userId}` - Alternative bearer format

For MVP, we'll use local storage to persist user ID. Full auth (OAuth, etc.) is out of scope per project constraints.

## Tasks

1. Create auth context/store:
   - Store userId in localStorage
   - Provide auth context to app
   - Handle "logged out" state

2. Create Login/Register page:
   - Simple UI to enter or generate a user ID
   - For MVP: "Create New User" button generates UUID
   - Store in localStorage
   - Redirect to enrollment or home

3. Protected routes:
   - Wrap routes that require auth
   - Redirect to login if no user ID

4. Auth header interceptor:
   - Add X-User-ID header to all API requests
   - Handle 401/403 responses

5. Logout flow:
   - Clear localStorage
   - Reset app state
   - Redirect to login

## User States
1. **Not logged in**: No userId - show login page
2. **Logged in, no enrollment**: Has userId but 404 on enrollment - prompt to select program
3. **Logged in, enrolled**: Has userId and active enrollment - show workout dashboard

## UI Design Notes
- Keep login simple and clean
- Use /frontend-design skill for implementation
- Mobile-first approach

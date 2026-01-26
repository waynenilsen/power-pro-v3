# Configure API Client and Types

Create a typed API client for the Go backend.

## Tasks

1. Create `frontend/src/api/` directory structure:
   ```
   src/api/
   ├── client.ts      # Base fetch wrapper
   ├── types.ts       # API response types
   └── endpoints/     # Organized by resource
       ├── programs.ts
       ├── workouts.ts
       └── users.ts
   ```

2. Create base API client (`client.ts`):
   - Wrapper around fetch
   - Automatic JSON parsing
   - Error handling
   - Base URL configuration (uses Vite proxy)
   - X-User-ID header injection for dev auth

3. Generate/create TypeScript types:
   - Check if `/docs/openapi.yaml` exists and generate types
   - Or manually create types based on API responses
   - Types for: User, Program, Workout, Exercise, Set, etc.

4. Create endpoint functions:
   - Programs: list, get, enroll
   - Workouts: list, get, create, update
   - Users: get profile, update profile

## API Authentication

For development, use X-User-ID header:
```typescript
headers: {
  'X-User-ID': userId,
  'Content-Type': 'application/json'
}
```

## Acceptance Criteria

- Typed API client with proper error handling
- Types match backend API responses
- Easy to use endpoint functions

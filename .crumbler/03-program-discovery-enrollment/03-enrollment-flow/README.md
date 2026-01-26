# Enrollment Flow

Implement the full enrollment, switch, and unenroll functionality.

## API Endpoints
- `POST /users/{userId}/enrollment` - Enroll in program (via `programs.enrollInProgram`)
- `GET /users/{userId}/enrollment` - Get current enrollment (via `users.getEnrollment`)
- `DELETE /users/{userId}/enrollment` - Unenroll (via `users.unenroll`)

## Requirements

### 1. Enroll Mutation Hook
Create `useEnrollInProgram` mutation in `frontend/src/hooks/useCurrentUser.ts`:
- Call `programs.enrollInProgram(userId, { programId })`
- On success, invalidate enrollment query
- Return mutation with loading/error states

### 2. Unenroll Mutation Hook
Create `useUnenroll` mutation:
- Call `users.unenroll(userId)`
- On success, invalidate enrollment query
- Include confirmation before executing

### 3. Enrollment UI on Program Detail Page
Update `ProgramDetails.tsx`:
- If NOT enrolled: Show "Enroll" button
- If enrolled in THIS program: Show "You're enrolled" badge + "Unenroll" option
- If enrolled in DIFFERENT program: Show "Switch to [Program Name]" button
  - Switching = unenroll current + enroll new (or backend handles atomically)

### 4. Confirmation Dialogs
- **Enroll confirmation**: "Start training with [Program Name]?"
- **Switch confirmation**: "Switch from [Current] to [New]? Your progress will reset."
- **Unenroll confirmation**: "Are you sure? Your enrollment progress will be lost."

### 5. Error Handling
- Handle enrollment failures gracefully
- Show toast/alert on success/error
- Handle case where user doesn't have required training maxes

## State Management
- Use TanStack Query mutations
- Invalidate `queryKeys.users.enrollment(userId)` on changes
- Consider invalidating program-related queries if needed

## Design Notes
- Use modal/dialog for confirmations (can use simple custom dialog)
- Make destructive actions (unenroll, switch) visually distinct
- Loading states on buttons during mutation

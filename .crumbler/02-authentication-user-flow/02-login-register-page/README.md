# Login/Register Page

Implement the login page UI using /frontend-design skill.

## Requirements

1. Replace placeholder `/frontend/src/routes/pages/Login.tsx`:
   - Dark industrial theme matching existing design system
   - Mobile-first responsive design

2. Two modes:
   - **New User**: "Create Account" button generates UUID, calls `createUser()`
   - **Existing User**: Input field to enter existing user ID (for development/testing)

3. UI Components:
   - App branding/logo area
   - Welcome message
   - "Create New Account" primary CTA button
   - Collapsible "Use Existing ID" section for development
   - Loading state during user creation

4. Flow:
   - On successful login/create → redirect to home or program enrollment
   - Use `useNavigate()` from react-router-dom

## Design Notes
- Use existing color palette from index.css
- Accent color (#f97316) for primary actions
- Match the dark industrial aesthetic
- Use lucide-react for icons
- Use /frontend-design skill for implementation

## Files to Modify
- `/frontend/src/routes/pages/Login.tsx`

## Acceptance Criteria
- Clean, professional login UI
- Creates new user with UUID
- Optional existing ID input for dev use
- Redirects after successful auth
- Mobile and desktop responsive

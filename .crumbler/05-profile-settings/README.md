# Profile & Settings

User profile and lift maxes management.

## API Endpoints Used
- `GET /users/{userId}/lift-maxes` - List user's lift maxes
- `POST /users/{userId}/lift-maxes` - Create lift max
- `PUT /lift-maxes/{id}` - Update lift max
- `DELETE /lift-maxes/{id}` - Delete lift max
- `GET /lift-maxes/{id}/convert` - Convert between 1RM and training max
- `GET /lifts` - List all lifts (for selecting which lift to set max for)

## Tasks

1. Profile Page (`/profile`):
   - Current user ID display (for debugging/sharing)
   - Current program enrollment info
   - Quick link to lift maxes
   - Logout button

2. Lift Maxes List:
   - Show all user's recorded maxes
   - Group by lift
   - Show both 1RM and Training Max if both exist
   - Add new max button

3. Add/Edit Lift Max Form:
   - Select lift (dropdown from available lifts)
   - Enter weight value
   - Select type (1RM or Training Max)
   - Date (default today)
   - Submit/Cancel

4. Training Max Calculator:
   - Input 1RM
   - Calculate training max (typically 90%)
   - Quick "Save as Training Max" button

5. Lift Max History:
   - Show progression over time for a specific lift
   - Simple list or chart

## UI Design Notes
- Lift maxes are important for program to work
- Make it easy to add/update training maxes
- Highlight if missing required training maxes for current program
- Use /frontend-design skill

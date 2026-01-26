# Program Discovery & Enrollment

Build UI for browsing programs and enrolling.

## API Endpoints Used
- `GET /programs` - List all programs (paginated)
- `GET /programs/{id}` - Get program details with cycles, weeks, days
- `POST /users/{userId}/program` - Enroll in program
- `GET /users/{userId}/program` - Get current enrollment
- `DELETE /users/{userId}/program` - Unenroll

## Tasks

1. Programs list page (`/programs`):
   - Grid/list of available programs
   - Show: name, description, metadata (weeks, days/week, lifts required, estimated session time)
   - Pagination or infinite scroll

2. Program detail page (`/programs/:id`):
   - Full program details
   - Sample week breakdown showing what a typical training week looks like
   - Lifts required by the program
   - Estimated session duration
   - "Enroll" button (or "Switch to this program" if already enrolled)

3. Enrollment flow:
   - When enrolling, user needs training maxes set up
   - If missing training maxes, prompt to set them first
   - Confirmation dialog before enrolling
   - Handle switch from existing program gracefully

4. Current program indicator:
   - Show which program user is enrolled in
   - Quick access from home/navigation

5. Unenroll flow:
   - Confirmation dialog
   - Clear enrollment state

## UI Design Notes
- Programs should feel discoverable and inspiring
- Don't overuse cards - simple list items can work well
- Mobile-first with clean desktop view
- Use /frontend-design skill

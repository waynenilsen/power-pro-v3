# Set Up React Router

Configure React Router with the application route structure.

## Tasks

1. Install react-router-dom:
   ```bash
   bun add react-router-dom
   ```

2. Create route structure in `src/routes/`:
   ```
   src/routes/
   ├── index.tsx       # Route definitions
   └── pages/
       ├── Home.tsx
       ├── Login.tsx
       ├── Programs.tsx
       ├── ProgramDetails.tsx
       ├── Workout.tsx
       ├── History.tsx
       └── Profile.tsx
   ```

3. Define routes:
   - `/` - Home/Landing
   - `/login` - Login page
   - `/programs` - Program discovery
   - `/programs/:id` - Program details
   - `/workout` - Current workout view
   - `/history` - Workout history
   - `/profile` - Profile & settings

4. Create placeholder page components (just titles for now)

5. Set up BrowserRouter in main.tsx

## Acceptance Criteria

- All routes accessible and render correct placeholder
- URL navigation works
- 404 handling for unknown routes

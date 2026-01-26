# Frontend Project Setup

Set up the React + TypeScript + Vite frontend project with Bun.

## Stack
- Bun (package manager & runtime)
- Vite (build tool)
- React 18+ (UI framework)
- TypeScript (type safety)
- TailwindCSS (styling)
- React Router (routing)
- TanStack Query (data fetching/caching)
- Zustand (state management, if needed)

## Tasks

1. Initialize frontend project:
   - Create `frontend/` directory in project root
   - Initialize with `bun create vite frontend --template react-ts`
   - Install dependencies

2. Configure TailwindCSS:
   - Install tailwindcss, postcss, autoprefixer
   - Create tailwind.config.js with mobile-first responsive breakpoints (2 breakpoints only: mobile and desktop)
   - Set up base styles

3. Configure API client:
   - Create typed API client for Go backend (use fetch or axios)
   - Generate types from OpenAPI spec or manually create them
   - Set up authentication headers (X-User-ID for now)

4. Set up routing:
   - Install react-router-dom
   - Create route structure:
     - `/` - Home/Landing
     - `/login` - Login page
     - `/programs` - Program discovery
     - `/programs/:id` - Program details
     - `/workout` - Current workout view
     - `/history` - Workout history
     - `/profile` - Profile & settings

5. Set up TanStack Query:
   - Configure QueryClient with defaults
   - Create custom hooks for API calls
   - Set up devtools for development

6. Create base layout:
   - Bottom navigation for mobile (4 tabs: Home, Workout, History, Profile)
   - Side navigation for desktop
   - Responsive container

7. Add proxy configuration for dev server to backend at localhost:8080

## Mobile-first Design
- Base styles for mobile
- Single desktop breakpoint at `md:` (768px)
- No tablet-specific breakpoint

## API Backend
- Go backend running at `http://localhost:8080`
- Uses X-User-ID header for auth (development mode)
- See `/docs/openapi.yaml` for full API spec

## Use /frontend-design skill
Use the `/frontend-design` skill when implementing UI components to ensure high design quality.

# Initialize Vite React TypeScript Project

Set up the base frontend project with Bun, Vite, React, and TypeScript.

## Tasks

1. Create `frontend/` directory in project root
2. Initialize project with `bun create vite frontend --template react-ts`
3. Install core dependencies:
   - react, react-dom (should come with template)
   - typescript (should come with template)
4. Verify the project builds and runs with `bun run dev`
5. Clean up default Vite template files (remove demo content from App.tsx)
6. Update .gitignore if needed for frontend artifacts
7. Add proxy configuration in vite.config.ts for backend at localhost:8080

## Acceptance Criteria

- `frontend/` directory exists with working Vite + React + TypeScript setup
- `bun run dev` starts development server successfully
- `bun run build` creates production build without errors
- Proxy configured to forward `/api` requests to `http://localhost:8080`

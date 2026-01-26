# PowerPro Frontend

React-based frontend for the PowerPro powerlifting tracking application.

## Tech Stack

- **React 19** - UI framework
- **TypeScript** - Type safety
- **Vite** - Build tool and dev server
- **TailwindCSS v4** - Utility-first CSS
- **React Router v7** - Client-side routing
- **TanStack Query v5** - Server state management
- **Vitest** - Testing framework
- **React Testing Library** - Component testing utilities

## Getting Started

### Prerequisites

- Node.js 20+ or Bun 1.x
- PowerPro API backend running on `http://localhost:8080`

### Installation

```bash
# Install dependencies
bun install
# or
npm install
```

### Development

```bash
# Start dev server with hot reload
bun run dev
# or
npm run dev
```

The app will be available at `http://localhost:5173`. API requests to `/api/*` are proxied to the backend.

### Building

```bash
# Type-check and build for production
bun run build
# or
npm run build
```

### Testing

```bash
# Run tests in watch mode
bun run test

# Run tests once
bun run test:run

# Run tests with UI
bun run test:ui

# Run tests with coverage
bun run test:coverage
```

### Linting and Formatting

```bash
# Run ESLint
bun run lint

# Format code with Prettier
bun run format

# Check formatting without changes
bun run format:check
```

## Project Structure

```
src/
├── api/                 # API client and type definitions
│   ├── client.ts        # HTTP client with fetch wrapper
│   ├── types.ts         # TypeScript types for API responses
│   ├── endpoints/       # Endpoint-specific API functions
│   └── index.ts         # Barrel export
├── components/          # Reusable UI components
│   └── layout/          # Layout components (SideNav, BottomNav, etc.)
├── hooks/               # Custom React hooks
│   ├── usePrograms.ts   # Program-related hooks
│   ├── useWorkouts.ts   # Workout-related hooks
│   └── useCurrentUser.ts
├── lib/                 # Utilities and configuration
│   └── query.ts         # TanStack Query setup
├── routes/              # React Router configuration
│   ├── index.tsx        # Route definitions
│   └── pages/           # Page components
├── test/                # Test utilities and setup
│   └── setup.ts         # Vitest setup
├── main.tsx             # Application entry point
└── index.css            # Global styles and design system
```

## Design System

The frontend uses a dark industrial aesthetic with an electric orange accent color. Key design tokens are defined in `src/index.css`:

- **Background**: `#0a0a0b`
- **Surface**: `#18181b`
- **Accent**: `#f97316` (orange)
- **Foreground**: `#fafafa`

## API Integration

The API client in `src/api/client.ts` provides typed functions for all backend endpoints. Configure the client with user credentials:

```typescript
import { configureClient } from './api/client';

configureClient({
  userId: 'user-uuid',
  isAdmin: false,
});
```

All subsequent requests will include the appropriate headers.

## Available Scripts

| Script | Description |
|--------|-------------|
| `dev` | Start development server |
| `build` | Type-check and build for production |
| `preview` | Preview production build |
| `lint` | Run ESLint |
| `format` | Format code with Prettier |
| `format:check` | Check formatting |
| `test` | Run tests in watch mode |
| `test:run` | Run tests once |
| `test:ui` | Run tests with Vitest UI |
| `test:coverage` | Run tests with coverage |

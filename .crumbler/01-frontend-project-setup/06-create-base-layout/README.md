# Create Base Layout

Create the responsive base layout with navigation.

## Tasks

1. Create layout components in `src/components/layout/`:
   ```
   src/components/layout/
   ├── Layout.tsx          # Main layout wrapper
   ├── BottomNav.tsx       # Mobile bottom navigation
   ├── SideNav.tsx         # Desktop side navigation
   └── Container.tsx       # Responsive content container
   ```

2. Bottom Navigation (mobile):
   - Fixed to bottom of screen
   - 4 tabs: Home, Workout, History, Profile
   - Icon + label for each
   - Active state styling
   - Hidden on desktop (`md:hidden`)

3. Side Navigation (desktop):
   - Fixed to left side
   - Same 4 navigation items
   - Hidden on mobile (`hidden md:flex`)
   - Collapsible or always visible

4. Layout wrapper:
   - Handles responsive switching
   - Proper content padding for nav
   - Mobile: padding-bottom for bottom nav
   - Desktop: padding-left for side nav

5. Use icons (lucide-react or heroicons):
   ```bash
   bun add lucide-react
   ```

## Use /frontend-design Skill

Use the `/frontend-design` skill when implementing the layout to ensure high design quality.

## Mobile-First Approach

- Default styles for mobile
- `md:` prefix for desktop overrides
- Bottom nav: `fixed bottom-0 md:hidden`
- Side nav: `hidden md:fixed md:left-0`

## Acceptance Criteria

- Seamless responsive navigation
- No layout shift on resize
- Active route highlighting
- Clean, modern design

# Configure TailwindCSS

Set up TailwindCSS with mobile-first responsive design.

## Tasks

1. Install TailwindCSS and dependencies:
   ```bash
   bun add -D tailwindcss postcss autoprefixer
   bunx tailwindcss init -p
   ```

2. Configure tailwind.config.js:
   - Set content paths for React files
   - Configure theme with only 2 breakpoints:
     - Default: mobile (< 768px)
     - `md`: desktop (>= 768px)
   - No tablet-specific breakpoint

3. Update CSS:
   - Add Tailwind directives to index.css
   - Remove default Vite/React styles

4. Test that Tailwind classes work in a component

## Mobile-First Design

- All base styles are for mobile
- Use `md:` prefix for desktop-specific styles
- Example: `text-sm md:text-base`

## Acceptance Criteria

- TailwindCSS compiles and works in dev/build
- Only mobile and desktop breakpoints exist
- Default Vite styles removed

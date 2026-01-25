# 004: Program Search Implementation

## Task
Implement search capability for GET /programs endpoint by name substring.

## Query Parameter
- `?search=strength` - case-insensitive substring match

## Expected Behavior
- `?search=strength` returns Starting Strength
- `?search=531` returns Wendler 5/3/1
- `?search=gz` returns GZCLP
- Empty search is ignored (returns all)
- Search combines with filters using AND logic

## Changes Required

1. **Filter Options** - Add Search field (optional string)

2. **Repository Query**
   - Use: `WHERE LOWER(name) LIKE '%' || LOWER(?) || '%'`
   - Or: `WHERE name LIKE '%' || ? || '%' COLLATE NOCASE`
   - Integrate with existing filter query building

3. **Handler** - Parse `search` query parameter

## Reference
- Ticket: `phases/in-progress/002-frontend-readiness/sprints/todo/006-program-discovery/tickets/todo/004-search-implementation.md`
- Depends on: 003-filtering-implementation

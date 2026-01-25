# 006: Canonical Programs Documentation

## ERD Reference
Implements: REQ-DOC-001

## Description
Document the canonical programs for API consumers. This includes the slugs, program descriptions, and any special handling needed for canonical vs. user-created programs.

## Context / Background
Frontend teams need to know which programs are canonical (system-provided) so they can present them appropriately in the UI. Canonical programs should be discoverable via API, and the documentation should explain how to identify and use them.

## Acceptance Criteria
- [ ] Create or update API documentation with canonical program section
- [ ] Document all canonical slugs:
  - `starting-strength` - Starting Strength novice program
  - `texas-method` - Texas Method intermediate program
  - `531` - Wendler 5/3/1 program
  - `gzclp` - GZCLP linear progression program
- [ ] Document program characteristics:
  - Training days per week
  - Progression model type (linear/weekly/monthly)
  - Target experience level (beginner/intermediate)
- [ ] Document how to identify canonical programs:
  - Field or flag that indicates canonical status
  - Or: author is system user
- [ ] Document canonical program restrictions:
  - Users cannot modify canonical programs
  - Users can enroll in canonical programs
  - Users can copy canonical programs to create their own version (future)
- [ ] Document enrollment process:
  - How to enroll in a canonical program via API
  - What happens when user enrolls (copy prescriptions? reference?)
- [ ] Add to README or dedicated API docs file

## Technical Notes
- Canonical programs may be identified by:
  - `is_canonical` boolean field on programs table
  - `author_id` matching a known system user
  - Slug prefix or namespace (e.g., `canonical:starting-strength`)
- Check existing program model to understand how canonical status is represented
- Consider adding a dedicated endpoint: `GET /programs/canonical`
- Or filter parameter: `GET /programs?canonical=true`

## Dependencies
- Blocks: None
- Blocked by: 001, 002, 003, 004 (programs must be seeded to document accurately)

## Resources / Links
- Existing API documentation (if any)
- Program specifications: `programs/*.md`
- README.md for project overview

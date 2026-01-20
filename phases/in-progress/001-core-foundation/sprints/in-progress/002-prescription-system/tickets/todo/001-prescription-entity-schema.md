# 001: Prescription Entity Schema and Migration

## ERD Reference
Implements: REQ-PRSC-001, REQ-PRSC-002, REQ-PRSC-003, REQ-PRSC-004, REQ-PRSC-005, REQ-PRSC-006, REQ-PRSC-007

## Description
Create the database schema and migration for the Prescription entity. A prescription links a lift to a load strategy and set scheme, forming the fundamental unit of programming.

## Context / Background
The Prescription entity is the core abstraction that defines what a lifter should do for a single exercise slot. It combines a lift reference with polymorphic LoadStrategy and SetScheme stored as JSON. This ticket creates the database foundation for the prescription system.

## Acceptance Criteria
- [ ] Create `prescriptions` table with the following columns:
  - `id` (UUID, primary key)
  - `lift_id` (UUID, required, foreign key to lifts.id)
  - `load_strategy` (JSON, required, stores LoadStrategy polymorphic data)
  - `set_scheme` (JSON, required, stores SetScheme polymorphic data)
  - `order` (INTEGER, default 0, for ordering within context)
  - `notes` (TEXT, nullable, max 500 characters)
  - `rest_seconds` (INTEGER, nullable, for rest period specification)
  - `created_at` (TIMESTAMP, required)
  - `updated_at` (TIMESTAMP, required)
- [ ] Create foreign key constraint on `lift_id` referencing `lifts.id`
- [ ] Create index on `lift_id` for efficient lift-based queries
- [ ] Create CHECK constraint to ensure `notes` length <= 500 characters
- [ ] Create CHECK constraint to ensure `rest_seconds` >= 0 when not null
- [ ] Create goose migration file with proper up/down migrations
- [ ] JSON columns validated at application layer (schema ticket just stores JSON)

## Technical Notes
- LoadStrategy and SetScheme stored as JSON for flexibility and extensibility
- The `order` field supports ordering prescriptions within any parent context (workout day, template, etc.)
- JSON structure validation happens at domain layer (ticket 008), not schema level
- Use SQLite JSON1 extension for JSON storage
- Consider partial index on `order` for unique ordering within context (may need parent_context_id in future)

## Dependencies
- Blocks: 008, 009, 010 (Domain logic, CRUD API, and Resolution API depend on this schema)
- Blocked by: None (lifts table from sprint 001 already exists)
- Related: 002, 003, 005, 006, 007 (LoadStrategy and SetScheme implementations)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/002-prescription-system/erd.md
- Tech Stack: prompts/tech-stack.md

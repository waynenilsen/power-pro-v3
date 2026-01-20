# 001: Day Entity Schema and Migration

## ERD Reference
Implements: REQ-DAY-001, REQ-DAY-002, REQ-DAY-004, REQ-DAY-005

## Description
Create the database schema for the Day entity and DayPrescription join table, including the goose migration file. The Day entity represents a single training session with ordered exercises.

## Context / Background
A Day is a named training session with ordered exercise slots. All programs organize work by training days. Days can have metadata (like intensity level), a unique slug for API URLs, and contain ordered prescriptions.

## Acceptance Criteria
- [x] Day table created with: id (UUID), name (VARCHAR 50, NOT NULL), slug (VARCHAR 50, NOT NULL), metadata (JSONB, nullable)
- [x] DayPrescription join table created with: day_id, prescription_id, order (INTEGER)
- [x] Proper foreign key constraints to prescriptions table
- [x] Slug uniqueness constraint (within program context)
- [x] Goose migration file created
- [x] Migration tested (up and down)

## Technical Notes
- Day table: `days` with columns: id, name, slug, metadata, program_id, created_at, updated_at
- DayPrescription join table: `day_prescriptions` with columns: id, day_id, prescription_id, `order`, created_at
- Slug should be unique within a program (composite unique constraint on program_id + slug)
- Metadata is JSONB for flexibility (common keys: intensityLevel, focus)
- Order field enables same prescription in multiple days at different positions

## Dependencies
- Blocks: 002, 007, 008
- Blocked by: None (ERD-002 Prescription schema assumed complete)
- Related: ERD-002 (Prescription system)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/003-schedule-periodization/erd.md

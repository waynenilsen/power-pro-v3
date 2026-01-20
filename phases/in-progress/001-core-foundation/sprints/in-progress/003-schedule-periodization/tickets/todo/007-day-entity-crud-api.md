# 007: Day Entity CRUD API

## ERD Reference
Implements: REQ-DAY-001, REQ-DAY-002, REQ-DAY-003, REQ-DAY-004, REQ-DAY-005

## Description
Implement the Day entity repository, service, and CRUD API endpoints. A Day represents a training session with ordered prescriptions.

## Context / Background
The Day entity is the building block for weekly schedules. Days can be named (e.g., "Day A", "Heavy Day", "Monday"), have a URL-safe slug, optional metadata, and contain ordered prescriptions.

## Acceptance Criteria
- [ ] Day repository implemented with CRUD operations
- [ ] Day service implemented with business logic
- [ ] GET /days - list all days (with pagination)
- [ ] GET /days/{id} - get day with prescriptions
- [ ] POST /days - create day with name, slug, metadata
- [ ] PUT /days/{id} - update day
- [ ] DELETE /days/{id} - delete day (fails if used in weeks)
- [ ] POST /days/{id}/prescriptions - add prescription to day with order
- [ ] DELETE /days/{id}/prescriptions/{prescriptionId} - remove prescription from day
- [ ] PUT /days/{id}/prescriptions/reorder - reorder prescriptions
- [ ] Slug uniqueness validation within program
- [ ] Unit tests with >80% coverage
- [ ] Integration tests for all endpoints

## Technical Notes
- Follow existing patterns from ERD-002 implementation
- Day response includes embedded prescriptions with order
- Metadata stored as JSONB, validated common keys: intensityLevel (HEAVY/LIGHT/MEDIUM), focus (string)
- Same prescription can appear in multiple days (many-to-many)
- Order field determines prescription sequence within day

## Dependencies
- Blocks: 008, 014
- Blocked by: 001 (Day schema)
- Related: ERD-002 Prescription API

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/003-schedule-periodization/erd.md
- API Response Format: See ERD Section 5 "Day API Response"

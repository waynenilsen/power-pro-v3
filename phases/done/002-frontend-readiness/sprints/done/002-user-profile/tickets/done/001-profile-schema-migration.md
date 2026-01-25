# 001: Profile Schema Migration

## ERD Reference
Implements: REQ-PROFILE-001

## Description
Add the weight_unit column to the users table to store user preference for weight display (lb or kg). This migration extends the users table created/modified in Sprint 001.

## Context / Background
Powerlifters use different weight systems. The weight_unit preference will determine how weights are displayed in program calculations and throughout the API. This is the first user preference field and establishes the pattern for future preferences.

The column should have a sensible default (lb) so existing users don't need to set it immediately.

## Acceptance Criteria
- [ ] Create goose migration file adding column to users table:
  - `weight_unit` (TEXT, NOT NULL, default 'lb')
- [ ] Add CHECK constraint: weight_unit IN ('lb', 'kg')
- [ ] Migration has proper down migration that removes column
- [ ] Existing users get default value 'lb'
- [ ] Run migration in test environment and verify no errors
- [ ] Verify existing user data is preserved

## Technical Notes
- Use SQLite TEXT type for weight_unit
- CHECK constraint enforces valid values at database level
- Default value ensures all existing rows get 'lb'
- Migration should be idempotent where possible
- Consider migration file naming: `YYYYMMDDHHMMSS_add_weight_unit_to_users.sql`

## Dependencies
- Blocks: 002 (domain logic needs schema), 003 (endpoints need schema), 004 (authorization needs schema), 005 (tests need schema)
- Blocked by: Sprint 001 user schema migration (users table must exist with current structure)

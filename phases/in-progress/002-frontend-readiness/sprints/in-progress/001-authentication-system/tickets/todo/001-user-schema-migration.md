# 001: User Schema Migration

## ERD Reference
Implements: REQ-USER-001, REQ-USER-002, REQ-USER-003, REQ-MIDDLEWARE-003

## Description
Extend the existing users table with authentication fields: email, password_hash, name, and is_admin. This migration adds the columns needed for real authentication while preserving existing test users.

## Context / Background
The current users table only has id, created_at, and updated_at. For real authentication, we need email for identification, password_hash for credential verification, name for personalization, and is_admin for authorization.

Existing test users in the database should continue to work - they'll have NULL email/password_hash and can only be accessed via X-User-ID test header.

## Acceptance Criteria
- [ ] Create goose migration file adding columns to users table:
  - `email` (TEXT, nullable, unique when not null)
  - `password_hash` (TEXT, nullable)
  - `name` (TEXT, nullable, max 100 chars)
  - `is_admin` (INTEGER/BOOLEAN, default 0/false)
- [ ] Create unique index on email for efficient lookup (partial index, WHERE email IS NOT NULL)
- [ ] Add CHECK constraint for email max length (255 chars)
- [ ] Add CHECK constraint for name max length (100 chars)
- [ ] Migration has proper down migration that removes columns
- [ ] Existing test users are preserved with NULL email/password_hash
- [ ] Run migration in test environment and verify no errors

## Technical Notes
- Use SQLite TEXT type for all string columns
- Email uniqueness should use partial unique index: `CREATE UNIQUE INDEX idx_users_email ON users(email) WHERE email IS NOT NULL`
- is_admin uses INTEGER (0/1) since SQLite doesn't have native BOOLEAN
- Consider email normalization (lowercase) at application layer, not database
- password_hash should accommodate bcrypt output (~60 chars) or argon2 output (~97 chars)

## Dependencies
- Blocks: 002 (sessions table references users), 003 (auth service needs schema), 004 (middleware needs schema), 005 (endpoints need schema)
- Blocked by: None

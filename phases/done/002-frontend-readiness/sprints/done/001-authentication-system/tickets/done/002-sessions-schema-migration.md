# 002: Sessions Schema Migration

## ERD Reference
Implements: REQ-SESSION-001, REQ-SESSION-002, REQ-SESSION-003, REQ-SESSION-004, REQ-SESSION-005

## Description
Create the sessions table for server-side session storage. Sessions link authentication tokens to users with expiration tracking.

## Context / Background
PowerPro uses server-side sessions (not JWT) for authentication. Each session has a unique token that the client sends in the Authorization header. The server looks up the session to identify the user. Sessions expire after 7 days.

## Acceptance Criteria
- [ ] Create goose migration file creating sessions table:
  - `id` (TEXT, primary key, UUID)
  - `user_id` (TEXT, required, foreign key to users.id with ON DELETE CASCADE)
  - `token` (TEXT, required, unique)
  - `expires_at` (TEXT, required, ISO8601 timestamp)
  - `created_at` (TEXT, required, ISO8601 timestamp)
- [ ] Create unique index on token for efficient lookup
- [ ] Create index on user_id for user's sessions lookup
- [ ] Create index on expires_at for cleanup queries
- [ ] Migration has proper down migration that drops table
- [ ] Foreign key constraint to users table with CASCADE delete

## Technical Notes
- Token should be indexed for fast lookup during authentication
- Consider composite index on (user_id, expires_at) if querying active sessions per user
- expires_at stored as TEXT in ISO8601 format (SQLite convention)
- Session cleanup (deleting expired sessions) can be a background job, not part of this ticket
- Token format: 32 bytes from crypto/rand, base64 encoded (~44 chars)

## Dependencies
- Blocks: 003 (auth service creates sessions), 004 (middleware validates sessions), 005 (endpoints use sessions)
- Blocked by: 001 (users table must have required columns first - though FK only needs id which exists)

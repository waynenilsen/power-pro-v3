# 012: Manual Progression Trigger API

## ERD Reference
Implements: REQ-API-003
Related to: REQ-TRIG-005, REQ-MANUAL-001

## Description
Implement the REST API endpoint for manually triggering progressions. This enables testing, administrative override, and correction scenarios.

## Context / Background
During development, testing, and administrative operations, it's useful to manually trigger progressions outside the normal state advancement flow. This endpoint allows triggering a specific progression for a user with optional force mode to bypass idempotency checks.

## Acceptance Criteria
- [ ] Implement `POST /users/{userId}/progressions/trigger` endpoint:
  - Request body: progressionId (required), liftId (optional), force (default false)
  - Looks up the progression configuration
  - Creates a synthetic trigger event
  - Applies the progression using ProgressionService
  - Returns ProgressionResult
- [ ] Normal mode (force=false):
  - Respects idempotency checks
  - Returns applied=false if progression already applied for this context
  - Creates standard ProgressionLog entry
  - Safe for repeated calls
- [ ] Force mode (force=true):
  - Bypasses idempotency check
  - Always applies progression regardless of previous applications
  - Creates ProgressionLog entry with special marker (e.g., trigger_context.manual=true)
  - Use with caution - enables corrections but can cause double-increments
- [ ] Lift-specific triggering:
  - If liftId provided, only apply to that lift
  - If liftId not provided, apply to all lifts in progression configuration
- [ ] Authorization:
  - Users can trigger their own progressions
  - Admin users can trigger for any user
  - Return 403 for unauthorized access
- [ ] Response format:
  - Returns ProgressionResult with applied status
  - If multiple lifts affected, return array of ProgressionResults
- [ ] Error handling:
  - 404 if progressionId not found
  - 404 if liftId not found (when specified)
  - 400 for invalid request body
- [ ] Audit logging for force=true operations
- [ ] Test coverage including force mode scenarios

## Technical Notes
- Uses ProgressionService.Apply() internally
- Synthetic trigger event has triggerType = MANUAL (or similar marker)
- force=true should be logged/audited for security
- Consider rate limiting to prevent abuse
- Manual triggers don't affect state advancement (just LiftMax values)

## Dependencies
- Blocks: None
- Blocked by: 008 (Uses ProgressionService)
- Related: REQ-MANUAL-001 (Manual max updates via LiftMax CRUD)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/004-progression-rules/erd.md
- ProgressionResult format: Section 5 of ERD

# 010: Program Progression Configuration API

## ERD Reference
Implements: REQ-API-002
Related to: REQ-PROG-003, REQ-CYCLE-003

## Description
Implement the REST API endpoints for configuring which progressions apply to which programs and lifts. This enables program designers to specify progression rules with optional lift-specific overrides.

## Context / Background
Programs need to specify which progressions apply to their lifts. 5/3/1 needs different increments for upper vs lower body. Some programs (Sheiko) may disable progressions entirely. This API enables full configuration of the program-progression relationship with priority ordering and lift-specific overrides.

## Acceptance Criteria
- [ ] Implement `GET /programs/{programId}/progressions` endpoint:
  - Returns list of ProgramProgressions for the program
  - Includes progression details (joined from progressions table)
  - Ordered by priority
  - Response includes: id, programId, progressionId, liftId, priority, enabled, overrideIncrement
- [ ] Implement `POST /programs/{programId}/progressions` endpoint:
  - Creates new ProgramProgression link
  - Request body: progressionId (required), liftId (optional), priority (default 0), enabled (default true), overrideIncrement (optional)
  - Validates progressionId exists
  - Validates liftId exists if provided
  - Returns 201 with created configuration
  - Returns 400 for validation errors
  - Returns 409 if duplicate (same program, progression, lift combination)
- [ ] Implement `PUT /programs/{programId}/progressions/{configId}` endpoint:
  - Updates existing configuration
  - Can modify: priority, enabled, overrideIncrement
  - Cannot modify: programId, progressionId, liftId (delete and recreate instead)
  - Returns 200 with updated configuration
  - Returns 404 if not found
- [ ] Implement `DELETE /programs/{programId}/progressions/{configId}` endpoint:
  - Removes progression from program
  - Returns 204 on success
  - Returns 404 if not found
- [ ] Priority ordering:
  - Lower priority numbers evaluated first
  - Enables defining fallback progressions
- [ ] Lift-specific vs program-wide:
  - liftId = null applies to all lifts (program default)
  - liftId specified applies only to that lift
  - Lift-specific takes precedence over program-wide
- [ ] Enable/disable support:
  - enabled = false prevents progression from firing
  - Supports Sheiko-style "no auto-progression" programs
- [ ] Test coverage for all endpoints and configuration scenarios

## Technical Notes
- Consider nested resource or flat API design
- Priority 0 is highest priority (first to evaluate)
- overrideIncrement allows 5/3/1 pattern: one CycleProgression with +5lb default, override to +10lb for squat/deadlift
- Deleting a progression should fail if ProgramProgressions reference it

## Dependencies
- Blocks: 008 (Trigger integration reads these configurations)
- Blocked by: 003, 009 (Schema and Progression CRUD)
- Related: ERD-003 (Programs entity), 006, 007 (Progression implementations use these configs)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/004-progression-rules/erd.md
- ProgramProgression schema: ticket 003

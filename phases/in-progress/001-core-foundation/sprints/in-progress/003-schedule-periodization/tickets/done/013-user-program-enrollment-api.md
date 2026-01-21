# 013: User Program Enrollment API

## ERD Reference
Implements: REQ-STATE-001, REQ-PROG-003

## Description
Implement the user program enrollment functionality including the enrollment endpoint and initial state creation.

## Context / Background
Users need to enroll in programs to follow them. Enrollment creates a UserProgramState record with initial position (week 1, cycle iteration 1). Users can only be enrolled in one program at a time.

## Acceptance Criteria
- [ ] POST /users/{userId}/program - enroll user in program
  - Request body: { programId }
  - Creates UserProgramState with initial position
  - Returns created state
- [ ] GET /users/{userId}/program - get user's current program enrollment
  - Returns program info and current state
  - Returns 404 if not enrolled
- [ ] DELETE /users/{userId}/program - unenroll user from program
  - Deletes UserProgramState
  - Returns 204 No Content
- [ ] User can only be enrolled in one program at a time
- [ ] Re-enrollment replaces existing enrollment (or requires explicit unenroll)
- [ ] Unit tests with >80% coverage
- [ ] Integration tests for all endpoints

## Technical Notes
- Initial state: currentWeek = 1, currentCycleIteration = 1
- Enrollment records enrolled_at timestamp
- Consider: should re-enrollment be allowed? Or require DELETE first?
- Enrollment validates program exists

## Dependencies
- Blocks: 014, 015
- Blocked by: 006 (UserProgramState schema), 012 (Program API)
- Related: 015 (State advancement)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/003-schedule-periodization/erd.md

# 005: Workflow Documentation

## ERD Reference
Implements: REQ-DOC-005
Related to: NFR-001, NFR-002

## Description
Document common API workflows that involve multiple endpoints. Multi-step workflows need guidance beyond individual endpoint documentation.

## Context / Background
Real-world API usage involves sequences of API calls. Workflow documentation shows developers how to accomplish common tasks end-to-end, reducing integration time and errors.

## Acceptance Criteria
- [ ] User onboarding workflow documented
- [ ] Program enrollment workflow documented
- [ ] Workout generation workflow documented
- [ ] Progression trigger workflow documented
- [ ] Documentation is clear and accessible (NFR-001)
- [ ] Examples are copy-paste ready (NFR-002)

## Technical Notes
- **User Onboarding Workflow**:
  1. Register/authenticate user
  2. Create lifts for the user
  3. Set initial maxes for lifts
  4. Enroll user in a program

- **Program Enrollment Workflow**:
  1. Get available programs
  2. View program details (cycle, days, prescriptions)
  3. Enroll user in program
  4. Set initial program state

- **Workout Generation Workflow**:
  1. Get current user program state
  2. Generate workout for current day
  3. Advance state after workout completion

- **Progression Trigger Workflow**:
  1. Understand program's progression rules
  2. Trigger progression (manual or automatic)
  3. Verify max was updated
  4. View progression history

Each workflow should include:
- Prerequisites
- Step-by-step API calls
- Example requests/responses
- Common pitfalls

## Dependencies
- Blocks: None
- Blocked by: 001-endpoint-documentation, 002-example-requests, 003-example-responses, 004-error-documentation
- Related: Program configuration tickets (010-014)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/006-api-developer-experience/erd.md

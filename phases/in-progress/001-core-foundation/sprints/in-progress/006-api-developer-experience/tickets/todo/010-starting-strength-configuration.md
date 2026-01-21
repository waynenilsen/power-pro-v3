# 010: Starting Strength Configuration

## ERD Reference
Implements: REQ-PROG-001
Related to: NFR-004, NFR-006

## Description
Document the complete Starting Strength program configuration, validating linear per-session progression. This serves as both documentation and validation that the system can represent this program.

## Context / Background
Starting Strength is a beginner program with linear progression (add weight every session). Documenting this configuration validates the core architecture and provides an example for developers.

## Acceptance Criteria
- [ ] Complete configuration documented (lifts, prescriptions, days, cycle, progression)
- [ ] Example workout output shown
- [ ] Progression behavior demonstrated
- [ ] Program example uses real API calls that can be verified (NFR-004)
- [ ] Examples produce expected outputs when executed (NFR-006)

## Technical Notes
- **Starting Strength Program Characteristics**:
  - Linear progression: add weight every session
  - 3x5 for main lifts (squat, bench/press, deadlift)
  - Alternating A/B days
  - 1-week cycle (3 training days)

- **Configuration Components**:
  - **Lifts**: Squat, Bench Press, Overhead Press, Deadlift, Power Clean
  - **Prescriptions**: 3x5 at working weight
  - **Days**:
    - Day A: Squat 3x5, Bench 3x5, Deadlift 1x5
    - Day B: Squat 3x5, Press 3x5, Power Clean 5x3
  - **Cycle**: 1 week, alternating A-B-A, B-A-B
  - **Progression**: Linear, add 5lbs per session (upper), 10lbs (lower/deadlift)

- **Documentation Format**:
  - JSON configuration for each entity
  - API calls to create the program
  - Example workout generation output
  - Example progression trigger and result

- **Validation**:
  - Configuration should be executable against the API
  - Workout output should match expected values
  - Progression should update maxes correctly

## Dependencies
- Blocks: None
- Blocked by: None (can be done in parallel with other program configs)
- Related: 011-014 (other program configurations), 005-workflow-documentation

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/006-api-developer-experience/erd.md
- Program Reference: programs/starting-strength.md (if exists)

# 012: Wendler 5/3/1 BBB Configuration

## ERD Reference
Implements: REQ-PROG-003
Related to: NFR-004, NFR-006

## Description
Document the complete Wendler 5/3/1 Boring But Big (BBB) program configuration, validating weekly percentage variation and cycle-based progression.

## Context / Background
Wendler 5/3/1 uses a 4-week cycle with different rep schemes each week and percentage-based loading. This validates the system's periodization and cycle progression capabilities.

## Acceptance Criteria
- [ ] Complete configuration documented
- [ ] 4-week cycle structure shown
- [ ] Weekly percentage lookup demonstrated
- [ ] Cycle-end progression shown
- [ ] Program example uses real API calls that can be verified (NFR-004)
- [ ] Examples produce expected outputs when executed (NFR-006)

## Technical Notes
- **Wendler 5/3/1 BBB Program Characteristics**:
  - 4-week mesocycle
  - Weekly rep scheme changes (5/5/5+, 3/3/3+, 5/3/1+, deload)
  - Percentage-based loading via lookup tables
  - Progression at end of cycle (add to training max)

- **Configuration Components**:
  - **Lifts**: Squat, Bench, Deadlift, Overhead Press
  - **Prescriptions**: Main lifts with AMRAP final set, BBB assistance (5x10)
  - **Weeks** (4-week cycle):
    - Week 1 (5s week): 65%x5, 75%x5, 85%x5+
    - Week 2 (3s week): 70%x3, 80%x3, 90%x3+
    - Week 3 (5/3/1 week): 75%x5, 85%x3, 95%x1+
    - Week 4 (deload): 40%x5, 50%x5, 60%x5
  - **Cycle**: 4 weeks
  - **Progression**: Cycle-end, add 5lbs (upper) / 10lbs (lower) to training max

- **Lookup Table Usage**:
  - Week number determines percentages
  - Lookup key: week_in_cycle, set_number
  - Lookup value: percentage of training max

- **BBB Assistance**:
  - 5x10 at 50% of training max
  - Same lift as main or supplementary

- **Documentation Format**:
  - JSON configuration showing lookup tables for weekly percentages
  - Full 4-week cycle example
  - Workout output for each week
  - Cycle progression trigger and result

## Dependencies
- Blocks: None
- Blocked by: None (can be done in parallel with other program configs)
- Related: 010-011, 013-014 (other program configurations)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/006-api-developer-experience/erd.md
- Program Reference: programs/wendler-531.md (if exists)

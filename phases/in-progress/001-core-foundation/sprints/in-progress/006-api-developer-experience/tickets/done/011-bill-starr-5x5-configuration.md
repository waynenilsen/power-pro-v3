# 011: Bill Starr 5x5 Configuration

## ERD Reference
Implements: REQ-PROG-002
Related to: NFR-004, NFR-006

## Description
Document the complete Bill Starr 5x5 program configuration, validating ramping sets and Heavy/Light/Medium day structure. This validates daily intensity variation.

## Context / Background
Bill Starr 5x5 uses ramping sets (increasing weight across sets) and a Heavy/Light/Medium weekly structure. This tests the system's ability to handle daily intensity variation.

## Acceptance Criteria
- [ ] Complete configuration documented
- [ ] Daily intensity variation demonstrated (Heavy/Light/Medium)
- [ ] Ramping set scheme shown
- [ ] Program example uses real API calls that can be verified (NFR-004)
- [ ] Examples produce expected outputs when executed (NFR-006)

## Technical Notes
- **Bill Starr 5x5 Program Characteristics**:
  - 5x5 with ramping weights (not all sets at same weight)
  - Heavy/Light/Medium day structure
  - 1-week cycle
  - Linear progression week over week

- **Configuration Components**:
  - **Lifts**: Squat, Bench Press, Bent Over Row (or Power Clean)
  - **Prescriptions**: 5x5 ramping (sets increase in weight)
  - **Days**:
    - Monday (Heavy): Full intensity
    - Wednesday (Light): ~80% of heavy day
    - Friday (Medium): ~90% of heavy day
  - **Cycle**: 1 week
  - **Progression**: Linear, weekly increase

- **Ramping Set Scheme**:
  - Set 1: ~60% of top set
  - Set 2: ~70% of top set
  - Set 3: ~80% of top set
  - Set 4: ~90% of top set
  - Set 5: 100% (top set)

- **Daily Intensity (via lookup tables)**:
  - Heavy day: 100%
  - Light day: 80%
  - Medium day: 90%

- **Documentation Format**:
  - JSON configuration showing lookup table usage
  - API calls to create the program
  - Example workout output for each day type
  - Show how intensity modifier affects workout generation

## Dependencies
- Blocks: None
- Blocked by: None (can be done in parallel with other program configs)
- Related: 010, 012-014 (other program configurations)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/006-api-developer-experience/erd.md
- Program Reference: programs/bill-starr-5x5.md (if exists)

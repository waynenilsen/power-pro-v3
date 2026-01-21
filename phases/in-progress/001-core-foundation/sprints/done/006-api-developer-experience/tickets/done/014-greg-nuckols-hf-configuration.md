# 014: Greg Nuckols High Frequency Configuration

## ERD Reference
Implements: REQ-PROG-005
Related to: NFR-004, NFR-006

## Description
Document the complete Greg Nuckols High Frequency program configuration, validating daily undulation and 3-week cycle structure.

## Context / Background
Greg Nuckols' High Frequency programs feature daily undulating periodization within a 3-week cycle structure. This validates the system's ability to handle complex daily variation.

## Acceptance Criteria
- [ ] Complete configuration documented
- [ ] Daily intensity variation demonstrated
- [ ] 3-week cycle structure shown
- [ ] Program example uses real API calls that can be verified (NFR-004)
- [ ] Examples produce expected outputs when executed (NFR-006)

## Technical Notes
- **Greg Nuckols High Frequency Program Characteristics**:
  - High frequency per lift (3-4x per week)
  - Daily undulating periodization (different rep/intensity each day)
  - 3-week progression cycle
  - Can be customized per lift (squat, bench, deadlift programs mix and match)

- **Configuration Components**:
  - **Lifts**: Squat, Bench Press, Deadlift (each with own template)
  - **Prescriptions**: Vary by day within week
  - **Days**: 5-6 days, each day has different rep/intensity scheme
  - **Weeks** (3-week cycle):
    - Week 1: Base intensities
    - Week 2: Slightly higher
    - Week 3: Peak/test week
  - **Cycle**: 3 weeks
  - **Progression**: Cycle-based (adjust after each 3-week block)

- **Daily Undulation Example** (for one lift):
  - Day 1: Heavy singles (90%+)
  - Day 2: Volume (4x8 @ 70%)
  - Day 3: Moderate (3x5 @ 80%)
  - Day 4: Light technique (3x3 @ 65%)

- **3-Week Cycle Structure**:
  - Week 1: Baseline intensities (from lookup table)
  - Week 2: +2.5% on all work
  - Week 3: +5% on all work, may include PR attempts
  - Then reset and increase training max

- **Lookup Table Complexity**:
  - Keys: week_in_cycle, day_in_week
  - Values: intensity percentages
  - More complex lookup than simpler programs

- **Documentation Format**:
  - JSON configuration with multi-dimensional lookup tables
  - Full 3-week cycle workout examples
  - Show daily variation within a week
  - Show weekly progression within cycle
  - Cycle reset and progression

## Dependencies
- Blocks: None
- Blocked by: None (can be done in parallel with other program configs)
- Related: 010-013 (other program configurations)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/006-api-developer-experience/erd.md
- Program Reference: programs/greg-nuckols.md (if exists)

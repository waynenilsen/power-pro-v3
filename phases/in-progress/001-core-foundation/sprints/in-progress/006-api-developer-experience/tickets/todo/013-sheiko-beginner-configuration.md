# 013: Sheiko Beginner Configuration

## ERD Reference
Implements: REQ-PROG-004
Related to: NFR-004, NFR-006

## Description
Document the complete Sheiko Beginner program configuration, validating high-volume programming with manual (no auto) progression.

## Context / Background
Sheiko programs are known for high volume and specific percentage-based work. The beginner variant doesn't auto-progress - the coach/athlete adjusts maxes manually. This validates the system's flexibility.

## Acceptance Criteria
- [ ] Complete configuration documented
- [ ] Multiple daily sessions shown if applicable
- [ ] Manual progression approach documented
- [ ] Program example uses real API calls that can be verified (NFR-004)
- [ ] Examples produce expected outputs when executed (NFR-006)

## Technical Notes
- **Sheiko Beginner Program Characteristics**:
  - High volume, moderate intensity
  - Percentage-based loading (typically 70-80% range)
  - No automatic progression (manual max adjustments)
  - Multiple training phases/blocks

- **Configuration Components**:
  - **Lifts**: Squat, Bench Press, Deadlift
  - **Prescriptions**: Multiple sets at various percentages
    - Example: 5x3@70%, 4x4@75%, 3x5@80%
  - **Days**: 3-4 training days per week
  - **Cycle**: Varies (often 4+ weeks per block)
  - **Progression**: None (manual) - user updates maxes when ready

- **High Volume Structure**:
  - Typical day might have:
    - Squat: 5x3@70%
    - Bench: 6x4@75%
    - Deadlift: 4x3@80%
  - Competition lift frequency: 2-3x per week each

- **Manual Progression Documentation**:
  - Explain that no progression rule is attached
  - Document how to manually update maxes via API
  - Document when/why a user would increase maxes

- **Documentation Format**:
  - JSON configuration without progression rules
  - Sample week of workouts
  - API calls to manually update maxes
  - Contrast with auto-progression programs

## Dependencies
- Blocks: None
- Blocked by: None (can be done in parallel with other program configs)
- Related: 010-012, 014 (other program configurations)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/006-api-developer-experience/erd.md
- Program Reference: programs/sheiko.md (if exists)

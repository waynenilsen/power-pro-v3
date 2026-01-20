# 007: Ramp SetScheme Implementation

## ERD Reference
Implements: REQ-SET-003, REQ-SET-004, REQ-SET-005
Related to: REQ-SET-001

## Description
Implement the Ramp set scheme, which generates sets with progressive percentages across a series of sets. Used for warmup progressions and Bill Starr style ramping sets.

## Context / Background
Ramp sets progress from lighter to heavier weights, typically used for warmups or programs like Bill Starr that ramp to a top set. Each set has a percentage of the target weight and specified reps. Sets below a threshold are classified as warmups.

## Acceptance Criteria
- [ ] Implement `RampSetScheme` struct with fields:
  - `Steps` ([]RampStep) - array of percentage/rep pairs (required)
  - `WorkSetThreshold` (float64) - percentage above which sets are work sets (default 80)
- [ ] Implement `RampStep` struct with:
  - `Percentage` (float64) - percentage of baseWeight for this set
  - `Reps` (int) - repetitions for this set
- [ ] Implement `GenerateSets` method:
  - Generate one `GeneratedSet` per step
  - Each set has:
    - `SetNumber`: 1 through len(Steps)
    - `Weight`: baseWeight * (step.Percentage / 100)
    - `TargetReps`: step.Reps
    - `IsWorkSet`: step.Percentage >= WorkSetThreshold
- [ ] Implement validation (REQ-SET-005):
  - At least one step required
  - All percentages must be > 0
  - All reps must be >= 1
  - WorkSetThreshold must be > 0 and <= 100
  - Return clear error message on validation failure
- [ ] Implement work set classification (REQ-SET-004):
  - Sets at or above WorkSetThreshold are work sets
  - Sets below threshold are warmup sets
  - Default threshold: 80%
- [ ] JSON serialization/deserialization for scheme storage
- [ ] Test coverage > 90% including:
  - Typical warmup ramp (50%, 63%, 75%, 88%, 100%)
  - Custom WorkSetThreshold values
  - Validation error cases
  - Edge cases (single step, all warmups, all work sets)

## Technical Notes
- Weight calculation at this level (unlike Fixed) since each set is different
- Should apply rounding to each set's weight (use rounding system from ticket 004)
- Steps should be in order (first step is first set)
- Consider: should steps be validated for ascending order? (probably not, allows flexibility)

## Dependencies
- Blocks: 008, 010 (Domain logic and resolution use this scheme)
- Blocked by: 004 (Needs rounding), 005 (Interface must be defined first)
- Related: 006 (Fixed scheme is the other scheme type)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/002-prescription-system/erd.md
- ERD example: `[{"percentage": 50, "reps": 5}, {"percentage": 63, "reps": 5}, ...]`

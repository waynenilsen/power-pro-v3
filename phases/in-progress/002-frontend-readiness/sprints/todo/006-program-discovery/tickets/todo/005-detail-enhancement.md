# 005: Program Detail Enhancement

## ERD Reference
Implements: REQ-DETAIL-001, REQ-DETAIL-002, REQ-DETAIL-003

## Description
Enhance the GET /programs/{id} endpoint to include sample week preview, lift requirements, and estimated session duration. These additions help users evaluate a program before enrolling.

## Context / Background
Program metadata (difficulty, days, focus) helps users filter, but users need more detail to make enrollment decisions. A sample week shows the training structure. Lift requirements tell users what equipment is needed. Estimated duration helps users plan their gym time.

## Acceptance Criteria
- [ ] Add `sampleWeek` to program detail response:
  - Array of day objects: `{day: number, name: string, exerciseCount: number}`
  - Shows first week (or first cycle for cyclical programs)
  - For A/B rotation programs like Starting Strength, show both A and B days
  - Example: `[{day: 1, name: "Workout A", exerciseCount: 3}, {day: 2, name: "Workout B", exerciseCount: 3}]`
- [ ] Add `liftRequirements` to program detail response:
  - Array of unique lift names used in the program
  - Sorted alphabetically
  - Example: `["Bench Press", "Deadlift", "Overhead Press", "Power Clean", "Squat"]`
- [ ] Add `estimatedSessionMinutes` to program detail response:
  - Calculation: (total sets per average day * 3 minutes) + (exercises per day * 2 minutes warmup)
  - Return as integer (minutes)
  - Use average across all days if days have different exercise counts
- [ ] Update ProgramDetailResponse DTO with new fields
- [ ] Update domain types if needed for intermediate data
- [ ] Query program_days and prescriptions tables for sample week data
- [ ] Query prescriptions and lifts tables for lift requirements
- [ ] Add integration tests verifying response structure
- [ ] Document that estimatedSessionMinutes is approximate

## Technical Notes
- Sample week query joins program_days and counts prescriptions per day
- Lift requirements query: SELECT DISTINCT lift_name FROM prescriptions WHERE program_id = ?
- For programs with weeks/cycles, use week 1 or first cycle
- For A/B rotation (no explicit weeks), show all days as the "sample week"
- Estimated duration formula is intentionally conservative
- Consider caching these derived values if queries become expensive
- New response fields should use camelCase per existing conventions

## Dependencies
- Blocks: 006 (E2E tests for detail enhancement)
- Blocked by: 003 (domain model updates from filtering)

## Resources / Links
- Existing detail handler: `internal/api/program_handler.go` (GetProgram)
- Program days schema: `migrations/` (find program_days table)
- Prescriptions schema: `migrations/` (find prescriptions table)
- ERD requirements: section 3, "Program Detail Enhancement"

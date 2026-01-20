# 011: Lookup Integration with LoadStrategy

## ERD Reference
Implements: REQ-LOOKUP-005

## Description
Integrate WeeklyLookup and DailyLookup tables with the LoadStrategy resolution from ERD-002. This enables lookups to modify base percentages during prescription resolution.

## Context / Background
During workout generation, the resolved weight must reflect lookup modifications. A prescription using 85% of training max with a week 2 lookup might resolve to 90% instead (if week 2 = 90%).

## Acceptance Criteria
- [ ] PercentOf LoadStrategy accepts optional lookup reference (weekly_lookup_id, daily_lookup_id)
- [ ] ResolvePrescription accepts context: {weekNumber, daySlug}
- [ ] Weekly lookup values modify base percentage during resolution
- [ ] Daily lookup values modify base percentage during resolution
- [ ] Lookup context passed through full resolution chain
- [ ] If lookup has percentages array, use set-specific percentages
- [ ] If lookup has percentageModifier, apply as multiplier
- [ ] Unit tests for all lookup integration scenarios
- [ ] Integration tests verifying resolved weights

## Technical Notes
- Resolution flow:
  1. Get base percentage from LoadStrategy
  2. Look up week context in WeeklyLookup (if present)
  3. Look up day context in DailyLookup (if present)
  4. Apply modifiers to base percentage
  5. Calculate weight from training max
- Modifiers can be additive or multiplicative (configured in lookup)
- If both weekly and daily lookups apply, determine precedence (typically weekly first, then daily)
- Context struct: `type LookupContext struct { WeekNumber int; DaySlug string }`

## Dependencies
- Blocks: 014
- Blocked by: 010 (Lookup CRUD API), ERD-002 LoadStrategy
- Related: 014 (Workout Generation)

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/003-schedule-periodization/erd.md
- ERD-002: Prescription System (LoadStrategy)

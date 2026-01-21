# 007: Code Comment Review for Complex Logic

## ERD Reference
Implements: REQ-DEBT-007

## Description
Review codebase for complex logic that lacks adequate comments and add explanatory comments where needed. Comments aid understanding for future developers and AI assistants.

## Context / Background
This is a COULD-priority requirement. While code should be self-documenting where possible, complex algorithms and business rules benefit from explanatory comments.

## Acceptance Criteria
- [ ] Identify complex algorithms and business rules in domain code
- [ ] Add explanatory comments for prescription resolution logic
- [ ] Add explanatory comments for progression evaluation logic
- [ ] Add explanatory comments for workout generation logic
- [ ] Comments explain "why" not just "what"
- [ ] No excessive commenting of obvious code

## Technical Notes
- Focus on non-obvious business logic
- Key areas to document:
  - Prescription resolution algorithm
  - Progression rule evaluation
  - Periodization calculations
  - Exercise ordering logic
- Comment style:
  - Use Go doc comments for exported functions
  - Use inline comments for complex logic blocks
  - Explain business rules, not implementation details
- Avoid:
  - Commenting obvious code
  - Redundant comments that repeat the code
  - Outdated comments

## Dependencies
- Blocks: None
- Blocked by: None
- Related: 008-api-documentation-sync

## Resources / Links
- ERD: phases/in-progress/001-core-foundation/sprints/in-progress/005-technical-debt-phase1/erd.md

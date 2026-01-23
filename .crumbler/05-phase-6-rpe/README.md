# Phase 6: RPE

Implement RPE-based load calculation.

## Features to Implement

- **LoadStrategy: `RPETarget`** - Calculate weight based on target RPE
- **Lookup: `RPEChart`** - Table mapping (reps, RPE) to percentage of 1RM
- **LoggedSet RPE field** - Record RPE for performed sets

## Programs Unlocked

| Program | Why |
|---------|-----|
| RTS Intermediate (partial) | RPETarget for load calculation, RPEChart lookup |

## Acceptance Criteria

- Users can prescribe sets at target RPE (e.g., "3 reps @RPE 8")
- System calculates appropriate weight using RPEChart lookup
- Actual RPE can be logged for performed sets
- E2E tests demonstrate RTS Intermediate partial implementation

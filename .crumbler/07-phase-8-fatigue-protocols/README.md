# Phase 8: Fatigue Protocols

Implement fatigue-based set schemes with variable set counts.

## Features to Implement

- **SetScheme: `FatigueDrop`** - Repeat sets until RPE hits target, dropping weight
- **SetScheme: `MRS`** - Max Rep Sets until failure
- **Variable set count** - Number of sets unknown until session complete
- **Trigger: `AfterSet` with RPE condition** - Fire when RPE threshold reached
- **"Repeat until" logic** - Continue sets until condition met

## Programs Unlocked

| Program | Why |
|---------|-----|
| GZCL Compendium | VDIP, MRS protocols, fatigue-based volume |
| RTS Intermediate (complete) | Load drop method, repeat sets method |

## Acceptance Criteria

- FatigueDrop continues until target RPE reached
- MRS continues until technical failure
- Session tracks variable number of sets dynamically
- E2E tests demonstrate both unlocked programs

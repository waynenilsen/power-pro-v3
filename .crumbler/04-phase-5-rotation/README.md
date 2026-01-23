# Phase 5: Rotation

Implement rotation schedules for lift focus changes.

## Features to Implement

- **Schedule: `Rotation`** - Rotating patterns for exercise selection
- **Exercise slot rotation** - Which lift is "primary" changes by week/cycle
- **Cycle-based lift focus** - Different lifts emphasized in different cycles

## Programs Unlocked

| Program | Why |
|---------|-----|
| nSuns CAP3 | 3-week rotation of which lift gets AMRAP focus |
| Inverted Juggernaut 5/3/1 | 16-week wave rotation (10s/8s/5s/3s phases) |
| Greyskull LP | AMRAP final set, AMRAPProgression with rotation |

## Acceptance Criteria

- Users can define rotation patterns for exercises
- System tracks rotation state across cycles
- Correct lift focus is applied based on rotation position
- E2E tests demonstrate all three unlocked programs

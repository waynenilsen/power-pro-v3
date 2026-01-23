# Phase 7: E1RM

Implement estimated 1RM calculations and relative load strategies.

## Features to Implement

- **LiftMax: `E1RM`** - Estimated 1RM calculated from performed sets
- **LoadStrategy: `FindRM`** - User discovers rep max (no prescribed weight)
- **LoadStrategy: `RelativeTo`** - Weight as percentage of today's top set
- **E1RM calculation** - Calculate estimated max from LoggedSet data
- **Prescription chaining** - Today's top set feeds back-off prescription

## Programs Unlocked

| Program | Why |
|---------|-----|
| GZCL Jacked & Tan 2.0 | FindRM (10RM→8RM→6RM etc), back-offs at % of found RM |
| Calgary Barbell 8-Week | RPE top set → E1RM → RelativeTo back-offs |
| Calgary Barbell 16-Week | Same as 8-week, extended phases |

## Acceptance Criteria

- System calculates E1RM from weight/reps/RPE
- FindRM sets have no prescribed weight, user works up to find it
- RelativeTo sets calculate weight from today's completed sets
- E2E tests demonstrate all three unlocked programs

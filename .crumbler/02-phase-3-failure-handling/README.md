# Phase 3: Failure Handling

Implement failure tracking and failure-based progression rules.

## Features to Implement

- **Failure tracking** - Detect when reps < target (failed set)
- **Progression: `DeloadOnFailure`** - Reduce weight after consecutive failures
- **Progression: `StageProgression`** - Change set/rep scheme on failure (e.g., 5x3 → 6x2 → 10x1)
- **Trigger: `OnFailure`** - Fire progression rules when failure is detected

## Programs Unlocked

| Program | Why |
|---------|-----|
| GZCLP | StageProgression (5x3 → 6x2 → 10x1), DeloadOnFailure |
| Texas Method | Implicit failure handling, reset protocols |
| Greg Nuckols Beginner | AMRAP drives weekly TM adjustment with failure handling |

## Acceptance Criteria

- System detects when user fails to hit prescribed reps
- DeloadOnFailure reduces weight after configurable failure count
- StageProgression transitions between set/rep schemes on failure
- E2E tests demonstrate all three unlocked programs

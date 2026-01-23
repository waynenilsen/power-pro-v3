# E2E Tests for Failure-Based Programs

Create end-to-end tests demonstrating the failure handling system with real program configurations.

## Context

Phase 3 unlocks these programs:
1. **GZCLP** - StageProgression (5x3 -> 6x2 -> 10x1), DeloadOnFailure
2. **Texas Method** - Implicit failure handling, reset protocols
3. **Greg Nuckols Beginner** - AMRAP drives weekly TM adjustment with failure handling

## Test Requirements

### 1. GZCLP T1 Progression Test

Demonstrate the full GZCLP T1 cycle:
```
Week 1: 5x3+ at 200 lbs - SUCCESS (hit 15+ reps) -> increase weight
Week 2: 5x3+ at 205 lbs - FAIL (only 13 reps) -> move to 6x2+
Week 3: 6x2+ at 205 lbs - SUCCESS (hit 12+ reps) -> increase weight
Week 4: 6x2+ at 210 lbs - FAIL (only 10 reps) -> move to 10x1+
Week 5: 10x1+ at 210 lbs - SUCCESS (hit 10+ reps) -> increase weight
Week 6: 10x1+ at 215 lbs - FAIL (only 8 reps) -> RESET to 5x3+ with deload
Week 7: 5x3+ at 185 lbs (15% deload) - continue fresh
```

### 2. GZCLP T2 Progression Test

Demonstrate T2 stage changes (no AMRAP):
```
Session 1: 3x10 at 100 lbs - SUCCESS -> increase weight
Session 2: 3x10 at 102.5 lbs - FAIL (only 27 reps) -> move to 3x8
Session 3: 3x8 at 102.5 lbs - SUCCESS -> increase weight
Session 4: 3x8 at 105 lbs - FAIL -> move to 3x6
Session 5: 3x6 at 105 lbs - FAIL -> RESET (manual intervention)
```

### 3. Texas Method Failure Test

Demonstrate stall handling:
```
Week 1: Friday 1x5 at 315 lbs - SUCCESS (5+ reps) -> increase weight
Week 2: Friday 1x5 at 320 lbs - FAIL (only 4 reps) -> keep weight
Week 3: Friday 1x5 at 320 lbs - FAIL (only 3 reps) -> 2 consecutive -> deload
Week 4: Friday 1x5 at 305 lbs (deload) -> continue
```

### 4. Greg Nuckols Beginner Bench Test

Demonstrate AMRAP-based progression with failure:
```
Day 1 (8-rep day):
  Week 1: 2x8 + AMRAP at 70% - 11 reps -> increase 10 lbs
  Week 2: 2x8 + AMRAP at 70% + 10 - 8 reps -> keep weight (hit target, no bonus)
  Week 3: 2x8 + AMRAP at 70% + 10 - 6 reps -> FAIL (under target) -> keep weight

Day 3 (4-rep day):
  Week 1: 2x4 + AMRAP at 80% - 6 reps -> increase 5 lbs
  Week 2: 2x4 + AMRAP at 80% + 5 - 4 reps -> keep weight (hit target exactly)
```

## File Structure

Create in `examples/` directory (matching existing pattern):
- `examples/gzclp-t1-progression.sh` - Full T1 cycle demonstration
- `examples/gzclp-t2-progression.sh` - T2 stage progression
- `examples/texas-method-failure.sh` - Stall and deload handling
- `examples/greg-nuckols-beginner-failure.sh` - AMRAP with failure

Each script should:
1. Set up a user with test data
2. Create program with appropriate progressions
3. Simulate logged sets (success and failure)
4. Verify weight/stage changes via API
5. Print clear output showing the progression

## Integration Tests

Create in `internal/integration/`:
- `failure_tracking_test.go` - Test failure counter increment/reset
- `deload_on_failure_test.go` - Test deload triggers correctly
- `stage_progression_test.go` - Test stage transitions
- `gzclp_program_test.go` - Full GZCLP program integration

## Acceptance Criteria

- [ ] GZCLP T1 e2e test passes all stage transitions
- [ ] GZCLP T2 e2e test passes stage transitions (no AMRAP)
- [ ] Texas Method e2e test shows stall detection and deload
- [ ] Greg Nuckols test shows AMRAP failure handling
- [ ] All example scripts are executable and produce clear output
- [ ] Integration tests cover edge cases:
  - [ ] First failure (counter = 1)
  - [ ] Threshold exactly met (e.g., 3 failures -> deload)
  - [ ] Success resets counter
  - [ ] Stage exhaustion triggers reset
- [ ] Tests can run independently (proper setup/teardown)

## Dependencies

- Crumb 01 (Failure Tracking) - complete
- Crumb 02 (DeloadOnFailure) - complete
- Crumb 03 (StageProgression) - complete

## Example Script Template

```bash
#!/usr/bin/env bash
# examples/gzclp-t1-progression.sh
# Demonstrates GZCLP T1 stage-based progression through failure

set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080/api/v1}"
USER_ID=""
PROGRAM_ID=""
LIFT_ID=""

# Helper functions
create_user() { ... }
create_lift() { ... }
create_program_with_stage_progression() { ... }
log_set() { ... }
verify_current_stage() { ... }
verify_weight() { ... }

echo "=== GZCLP T1 Stage Progression Demo ==="
echo ""

# Setup
create_user
create_lift "Squat"
create_program_with_stage_progression

# Week 1: 5x3+ SUCCESS
echo "Week 1: 5x3+ at 200 lbs"
log_set 200 5 3  # Set 1: 3 reps
log_set 200 5 3  # Set 2: 3 reps
log_set 200 5 3  # Set 3: 3 reps
log_set 200 5 3  # Set 4: 3 reps
log_set 200 5 5 --amrap  # Set 5: 5 reps (AMRAP) - Total: 17 reps
verify_current_stage 0
echo "SUCCESS: Total 17 reps >= 15 minimum"
echo ""

# ... continue for full demo

echo "=== Demo Complete ==="
```

## Notes

- Follow the pattern established in `examples/` directory
- Use curl for API calls
- Include verbose output explaining what's happening
- Tests should be deterministic (no randomness)
- Clean up test data after each run if possible

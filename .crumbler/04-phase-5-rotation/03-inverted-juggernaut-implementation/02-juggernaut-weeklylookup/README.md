# Juggernaut WeeklyLookup

Create factory functions to build WeeklyLookup configurations for Inverted Juggernaut 5/3/1.

## Context

The program needs WeeklyLookup entries for the 5/3/1 main work. The 5/3/1 percentages are CONSISTENT across all 4 waves - only the volume work changes per wave (handled separately).

## 5/3/1 Set Structure (All Waves)

### Week 1/5/9/13 (Accumulation)
| Set | Percentage | Reps |
|-----|------------|------|
| 1 | 65% | 5 |
| 2 | 75% | 5 |
| 3 | 85% | 5+ (AMRAP) |
| 4 | 75% | 5 |
| 5 | 65% | 5+ (AMRAP) |

### Week 2/6/10/14 (Intensification)
| Set | Percentage | Reps |
|-----|------------|------|
| 1 | 70% | 3 |
| 2 | 80% | 3 |
| 3 | 90% | 3+ (AMRAP) |
| 4 | 80% | 3 |
| 5 | 70% | 3+ (AMRAP) |

### Week 3/7/11/15 (Realization)
| Set | Percentage | Reps |
|-----|------------|------|
| 1 | 75% | 5 |
| 2 | 85% | 3 |
| 3 | 95% | 1+ (AMRAP) |
| 4 | 85% | 3 |
| 5 | 75% | 5+ (AMRAP) |

### Week 4/8/12/16 (Deload)
| Set | Percentage | Reps |
|-----|------------|------|
| 1 | 40% | 5 |
| 2 | 50% | 5 |
| 3 | 60% | 5 |

## Implementation

Create `internal/domain/juggernaut/weeklylookup.go`:

```go
package juggernaut

import "github.com/waynenilsen/power-pro-v3/internal/domain/weeklylookup"

// Create531WeeklyLookup creates the standard 5/3/1 weekly lookup.
// This defines the 5/3/1 main work sets that are consistent across all waves.
// The 4-week pattern repeats for all 4 waves (weeks 1-4 = same as 5-8 = 9-12 = 13-16).
func Create531WeeklyLookup(id string, programID *string) *weeklylookup.WeeklyLookup
```

The function should create entries for weeks 1-4 only, since the pattern repeats.

When resolving, we use `WeekInWave` (1-4) rather than `CurrentWeek` (1-16) to look up the correct entry.

## Approach

Since the 5/3/1 percentages repeat every 4 weeks, we only need 4 entries:
- Entry 1 (Accumulation): 65/75/85/75/65 @ 5/5/5+/5/5+
- Entry 2 (Intensification): 70/80/90/80/70 @ 3/3/3+/3/3+
- Entry 3 (Realization): 75/85/95/85/75 @ 5/3/1+/3/5+
- Entry 4 (Deload): 40/50/60 @ 5/5/5

The LookupContext will be set with `WeekNumber = WeekInWave` (from WaveInfo).

## Tests

Create `internal/domain/juggernaut/weeklylookup_test.go`:

1. Week 1 entry has correct 5/3/1 percentages (65/75/85/75/65)
2. Week 2 entry has correct percentages (70/80/90/80/70)
3. Week 3 entry has correct percentages (75/85/95/85/75)
4. Week 4 deload has correct percentages (40/50/60)
5. Reps arrays match expected values
6. GetByWeekNumber returns correct entry for each week

## Acceptance Criteria

- [ ] Create531WeeklyLookup creates valid WeeklyLookup
- [ ] All 4 weekly entries have correct percentages
- [ ] All 4 weekly entries have correct reps
- [ ] Deload week has only 3 sets (not 5)
- [ ] Unit tests pass

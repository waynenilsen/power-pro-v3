# Taper Load Strategy

Implement taper protocol for competition preparation (volume/intensity reduction).

## Tasks

1. Create `TaperLoadStrategy` that modifies base load as meet approaches
2. Implement taper multiplier curve based on days out:
   - Week 5 of Comp (final week): ~30-50% volume reduction
   - Week 4: ~40% volume reduction
   - Week 3: ~30% volume reduction
   - Week 2: ~20% volume reduction
   - Week 1: ~10% volume reduction
3. Integrate with existing LoadStrategy interface
4. Support configurable taper parameters per program
5. Add unit tests for taper calculations

## Taper Formula

```go
// Example taper curve
func GetTaperMultiplier(daysOut int) float64 {
    if daysOut > 35 { return 1.0 } // No taper outside Comp phase
    if daysOut <= 7 { return 0.5 }  // Final week
    if daysOut <= 14 { return 0.6 }
    if daysOut <= 21 { return 0.7 }
    if daysOut <= 28 { return 0.8 }
    return 0.9
}
```

## Acceptance Criteria

- [ ] TaperLoadStrategy implements LoadStrategy interface
- [ ] Volume reduces progressively as meet approaches
- [ ] Intensity maintains or increases slightly
- [ ] Unit tests verify taper curve

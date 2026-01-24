# Inverted Juggernaut 5/3/1 Implementation

Implement the 16-week wave rotation for Inverted Juggernaut program.

## Program Mechanics

The program runs 16 weeks with 4 rep waves, each lasting 4 weeks:

| Wave | Weeks | Focus | Volume Sets | Percentage |
|------|-------|-------|-------------|------------|
| 10s Wave | 1-4 | High volume | 9 sets | 60% TM |
| 8s Wave | 5-8 | Moderate volume | 7 sets | 65% TM |
| 5s Wave | 9-12 | Strength | 5 sets | 70% TM |
| 3s Wave | 13-16 | Peak intensity | 6 sets | 75% TM |

### Each Wave Contains 4 Phases
1. **Accumulation** (Week 1): High volume, moderate intensity
2. **Intensification** (Week 2): Moderate volume, higher intensity
3. **Realization** (Week 3): Lower volume, peak AMRAP
4. **Deload** (Week 4): Recovery

### 5/3/1 Percentage Structure (Each Week)
- Week 1: 65/75/85% (5 reps)
- Week 2: 70/80/90% (3 reps)
- Week 3: 75/85/95% (1+ AMRAP)
- Week 4: 40/50/60% (deload)

## Tasks

### 1. Extend Rotation Model for Waves
The Juggernaut needs a different rotation model - wave-based rather than lift-based.

Add to UserProgramState or create WaveState:
- `CurrentWaveIndex int` (0-3 for 10s/8s/5s/3s)
- `WeekWithinWave int` (1-4)

### 2. Create Juggernaut WeeklyLookup
Configure 16 weeks of percentages:
- Weeks 1-4: 10s wave percentages
- Weeks 5-8: 8s wave percentages
- Weeks 9-12: 5s wave percentages
- Weeks 13-16: 3s wave percentages

Each week also needs the 5/3/1 overlay (Accumulation/Intensification/Realization/Deload).

### 3. Implement Volume Set Generation
Different waves require different numbers of volume sets before the main 5/3/1 work:
- 10s wave: 9 sets x 5 reps @ 60%
- 8s wave: 7 sets x 5 reps @ 65%
- 5s wave: 5 sets x 5 reps @ 70%
- 3s wave: 6 sets x 3 reps @ 75%

### 4. Wave Advancement Logic
After each 4-week wave:
- Advance to next wave (0→1→2→3→0)
- Calculate new TM based on AMRAP performance
- Reset week counter

### 5. Create Integration Tests
Test scenarios:
- Week 1 (10s wave, accumulation): Correct volume sets and percentages
- Week 5 (8s wave, accumulation): Wave transition works
- Week 12 (5s wave, realization): AMRAP at correct percentage
- Week 16: Full cycle completes, resets properly

## Acceptance Criteria

- [ ] Wave state tracking (current wave, week within wave)
- [ ] Correct volume set count per wave (9/7/5/6)
- [ ] Correct base percentages per wave (60/65/70/75%)
- [ ] 5/3/1 overlay applies correctly each week
- [ ] Wave advancement after 4 weeks
- [ ] Full 16-week cycle completes and resets
- [ ] Integration tests pass for complete cycle

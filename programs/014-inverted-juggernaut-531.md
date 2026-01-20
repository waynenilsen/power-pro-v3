# Inverted Juggernaut 5/3/1 Method - Algorithm Specification

## Overview

The Inverted Juggernaut 5/3/1 is a hybrid powerlifting program that combines the Juggernaut Method's wave periodization with Wendler's 5/3/1 rep scheme. The program runs for 16 weeks (one complete cycle) and is designed to build both strength and work capacity through intelligent volume manipulation.

**Created by:** CoffeePWO

## Core Concepts

### Training Max Foundation
The entire program is built upon the Training Max (TM), which is **90% of the lifter's true or estimated 1 Rep Max (1RM)**:

```
Training Max = 1RM x 0.90
```

The program supports four main lifts:
- Bench Press
- Squat
- Overhead Press
- Deadlift

### 1RM Estimation Formula
When estimating 1RM from a weight/rep performance, the program uses:

```
Estimated 1RM = (Weight x Reps x 0.0333) + Weight
```

This is the Epley formula for 1RM estimation.

---

## Program Structure: 16-Week Cycle

The program consists of **four Rep Waves**, each lasting 4 weeks. Each wave contains four phases:

| Rep Wave | Duration | Focus |
|----------|----------|-------|
| 10s Wave | Weeks 1-4 | High volume foundation |
| 8s Wave | Weeks 5-8 | Moderate volume |
| 5s Wave | Weeks 9-12 | Strength emphasis |
| 3s Wave | Weeks 13-16 | Peak intensity |

Each 4-week wave follows the same phase structure:
1. **Week 1: Accumulation Phase** - High volume, moderate intensity
2. **Week 2: Intensification Phase** - Moderate volume, higher intensity
3. **Week 3: Realization Phase** - Lower volume, peak intensity with AMRAP
4. **Week 4: Deload Phase** - Recovery

---

## The "Inverted" Rep Scheme Explained

Traditional Juggernaut uses a descending set scheme within each session (e.g., 10, 8, 6 reps). The **Inverted** version reverses this, using an ascending pyramid that builds to a top AMRAP set, similar to 5/3/1.

### Key Innovation: 5/3/1 Set Structure Within Juggernaut Waves

Each main lift session uses the classic 5/3/1 percentage progression but with wave-specific top set targets:

**First Set (FS)** - Warmup working set
**Second Set (SS)** - Moderate intensity set
**Top Set (+)** - AMRAP set at peak intensity
**Back-Off Sets** - Volume work after the top set

---

## 5/3/1 Percentage Structure (Used Across All Waves)

The main lift 5/3/1 sets follow this exact structure each week:

### Week 1 (Accumulation)
| Set | Percentage | Reps |
|-----|------------|------|
| FS | 65% | 5 |
| SS | 75% | 5 |
| + | 85% | 5+ (AMRAP) |
| SS | 75% | 5 |
| FS+ | 65% | 5+ (AMRAP) |

### Week 2 (Intensification)
| Set | Percentage | Reps |
|-----|------------|------|
| FS | 70% | 3 |
| SS | 80% | 3 |
| + | 90% | 3+ (AMRAP) |
| SS | 80% | 3 |
| FS+ | 70% | 3+ (AMRAP) |

### Week 3 (Realization)
| Set | Percentage | Reps |
|-----|------------|------|
| FS | 75% | 5 |
| SS | 85% | 3 |
| + | 95% | 1+ (AMRAP) |
| SS | 85% | 3 |
| FS+ | 75% | 5+ (AMRAP) |

### Week 4 (Deload)
| Set | Percentage | Reps |
|-----|------------|------|
| 1 | 40% | 5 |
| 2 | 50% | 3 |
| 3 | 60% | 1 |

**Note:** All weights are calculated from the Training Max and rounded up to the nearest 2.5 units (using CEILING function).

---

## Juggernaut Wave-Specific Volume Work

Before the 5/3/1 sets, each wave includes wave-specific volume work that decreases as the program progresses. This is the "Juggernaut" component that builds work capacity.

### 10 Rep Wave - Volume Sets

**Accumulation Phase (Week 1):**
- 9 sets x 5 reps @ 60% TM
- RPE: 2-3 reps shy of failure
- Final set: 5+ (AMRAP)

**Intensification Phase (Week 2):**
- First 2 sets @ 55% and 62.5% TM
- Remaining sets @ 67.5% TM
- All sets: 3 reps minimum
- RPE: 1-2 reps shy of failure
- Final set: 3+ (AMRAP)

**Realization Phase (Week 3):**
- 4 sets with ascending percentages: 50%, 60%, 70%, 75% TM
- Reps: 5, 3, 1, AMRAP (to failure)

### 8 Rep Wave - Volume Sets

**Accumulation Phase (Week 1):**
- 7 sets x 5 reps @ 65% TM
- RPE: 2-3 reps shy of failure
- Final set: 5+ (AMRAP)

**Intensification Phase (Week 2):**
- First 2 sets @ 60% and 67.5% TM
- Remaining sets @ 72.5% TM
- All sets: 3 reps minimum
- RPE: 1-2 reps shy of failure
- Final set: 3+ (AMRAP)

**Realization Phase (Week 3):**
- 5 sets with ascending percentages: 50%, 60%, 70%, 75%, 80% TM
- Reps: 5, 3, 2, 1, AMRAP (to failure)

### 5 Rep Wave - Volume Sets

**Accumulation Phase (Week 1):**
- 5 sets x 5 reps @ 70% TM
- RPE: 2-3 reps shy of failure
- Final set: 5+ (AMRAP)

**Intensification Phase (Week 2):**
- First 2 sets @ 65% and 72.5% TM (2 reps each)
- Remaining 3 sets @ 77.5% TM (5 reps each)
- RPE: 1-2 reps shy of failure
- Final set: 5+ (AMRAP)

**Realization Phase (Week 3):**
- 5 sets with ascending percentages: 50%, 60%, 70%, 75%, 80% TM
- Reps: 5, 3, 2, 1, 1
- Final AMRAP @ 85% TM

### 3 Rep Wave - Volume Sets

**Accumulation Phase (Week 1):**
- 6 sets x 3 reps @ 75% TM
- RPE: 2-3 reps shy of failure
- Final set: 3+ (AMRAP)

**Intensification Phase (Week 2):**
- First 2 sets @ 70% and 77.5% TM (1 rep each)
- Remaining 4 sets @ 82.5% TM (3 reps each)
- RPE: 1-2 reps shy of failure
- Final set: 3+ (AMRAP)

**Realization Phase (Week 3):**
- 6 sets with ascending percentages: 50%, 60%, 70%, 75%, 80%, 85% TM
- Reps: 5, 3, 2, 1, 1, 1
- Final AMRAP @ 90% TM

---

## Realization Week Protocol

The Realization week (Week 3 of each wave) is designed to peak performance with:

1. **Ascending pyramid** - Weights increase with each set
2. **Decreasing reps** - Volume drops as intensity rises
3. **AMRAP finale** - The final set is performed to technical failure for as many reps as possible
4. **Progress marker** - AMRAP reps are recorded to calculate new Training Max for the next wave

### AMRAP Target Percentages by Wave

| Wave | AMRAP Percentage | Rep Standard |
|------|-----------------|--------------|
| 10s | 75% TM | 10 reps |
| 8s | 80% TM | 8 reps |
| 5s | 85% TM | 5 reps |
| 3s | 90% TM | 3 reps |

The "Rep Standard" is the minimum expected reps at that percentage. Exceeding this indicates the Training Max can be increased.

---

## Training Max Progression Algorithm

After each Rep Wave's Realization phase, a new Training Max is calculated based on AMRAP performance:

```
New TM = Previous TM + ((AMRAP Reps - Rep Standard) x Weight Increment)
```

Where:
- **Weight Increment for Bench/Press:** 2.5 units per rep over standard
- **Weight Increment for Squat/Deadlift:** 5 units per rep over standard

### Example Calculation
- Current Bench TM: 200
- 10s Wave Realization AMRAP: 13 reps (vs. 10 standard)
- Excess reps: 13 - 10 = 3
- Increase: 3 x 2.5 = 7.5
- New TM: 200 + 7.5 = 207.5

### Cycle-to-Cycle Progression

For subsequent 16-week cycles, Training Maxes increase automatically:
- **Squat/Deadlift:** +10 units per cycle
- **Bench/Press:** +5 units per cycle

---

## Weekly Training Split

The program uses a 4-day split with main lifts paired with accessory work:

### Day 1: Overhead Press + Upper Accessories
- Overhead Press (main work)
- Upper pec/front delt: 3 sets
- Side delts: 5 sets
- Triceps + Biceps: 3 sets

### Day 2: Squat + Lower Accessories
- Squat (main work)
- Horizontal pulling: 5 sets
- Quad + hamstrings: 3 sets
- Triceps + Biceps: 5 sets

### Day 3: Bench Press + Upper Accessories
- Bench Press (main work)
- Delts: 3 sets
- Chest: 3 sets
- Triceps + Biceps: 3 sets

### Day 4: Deadlift + Lower Accessories
- Deadlift (main work)
- Vertical rowing: 5 sets
- Quad/hamstrings: 3 sets
- Triceps + Biceps: 3 sets

### Accessory Rep Progression
Accessory work follows its own wave pattern:

| Phase | Reps per Set |
|-------|-------------|
| Accumulation | 15 |
| Intensification | 12 |
| Realization | 10 |
| Deload | (No accessories) |

---

## Deload Protocol

Week 4 of each wave is a mandatory deload:

### Main Lift Deload
| Set | Percentage | Reps |
|-----|------------|------|
| 1 | 40% | 5 |
| 2 | 50% | 5 |
| 3 | 60% | 5 |

- No AMRAP sets
- No back-off sets
- Minimal accessory work (or none)

---

## RPE Guidelines

The program uses Rate of Perceived Exertion to autoregulate volume work:

| Phase | RPE Target |
|-------|-----------|
| Accumulation | 2-3 reps shy of failure |
| Intensification | 1-2 reps shy of failure |
| Realization | To technical failure (AMRAP set only) |
| Deload | Very easy, recovery focus |

---

## Weight Calculation Formula

All working weights are calculated using:

```
Working Weight = CEILING(Training Max x Percentage, 2.5)
```

The CEILING function rounds up to the nearest 2.5 increment, ensuring consistent loading.

---

## Complete Wave Summary Table

| Wave | Accum % | Intens % | Real AMRAP % | Volume (Accum Sets) |
|------|---------|----------|--------------|---------------------|
| 10s | 60% | 67.5% | 75% | 9 sets |
| 8s | 65% | 72.5% | 80% | 7 sets |
| 5s | 70% | 77.5% | 85% | 5 sets |
| 3s | 75% | 82.5% | 90% | 6 sets |

---

## Program Execution Order

1. **Start:** Enter 1RM or calculate from recent performance
2. **Calculate TM:** Multiply 1RM by 0.90
3. **Run 10s Wave:** Weeks 1-4 (Accum, Intens, Real, Deload)
4. **Calculate new TM** from AMRAP performance
5. **Run 8s Wave:** Weeks 5-8
6. **Calculate new TM**
7. **Run 5s Wave:** Weeks 9-12
8. **Calculate new TM**
9. **Run 3s Wave:** Weeks 13-16
10. **Calculate new TM** for next cycle
11. **Repeat** with accumulated progression

---

## Key Algorithm Principles

1. **Wave Loading:** Volume decreases (10 -> 8 -> 5 -> 3) while intensity increases across the 16-week cycle
2. **Inverted Pyramid:** Each session builds to a peak set rather than starting heavy
3. **Dual AMRAP:** Both the top set and the final back-off set are taken to near-failure
4. **Autoregulation:** AMRAP performance drives progression
5. **Built-in Deload:** Every 4th week prevents accumulated fatigue
6. **5/3/1 Integration:** The proven 65/75/85 -> 70/80/90 -> 75/85/95 weekly progression

This hybrid approach combines Juggernaut's intelligent volume management with 5/3/1's proven intensity progression, creating a comprehensive strength program suitable for intermediate to advanced lifters.

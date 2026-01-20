# GZCLP Linear Progression Program Specification

## Program Overview

**GZCLP** (GZCL Linear Progression) is a beginner-to-intermediate strength training program created by Cody Lefever (/u/GZCL). This implementation (v4.1 by /u/Blacknoir) provides a 12-week spreadsheet with both 3-day and 4-day variants.

### Core Philosophy

GZCLP is built on the **tiered system** approach to training:

- **Tier 1 (T1)**: High intensity, low volume compound movements (main lifts)
- **Tier 2 (T2)**: Moderate intensity, moderate volume compound movements
- **Tier 3 (T3)**: Low intensity, high volume isolation/accessory movements

This pyramid structure ensures adequate stimulus across different rep ranges while managing fatigue and allowing for consistent progression.

### Program Variants

- **3-Day Program**: 3 training days per week, rotating through lift patterns
- **4-Day Program**: 4 training days per week with dedicated days for each main lift

---

## Configuration Options

### User Settings

The program requires the following initial configuration:

| Setting | Options | Purpose |
|---------|---------|---------|
| Gender | XY / XX | Used for Wilks score calculations |
| Training Days | 3 or 4 | Determines program structure |
| Unit System | Kg / Lbs | Weight display and rounding |
| Microplates | Yes / No | Enables finer weight increments |
| T1 Rep Scheme | Default / Modified | Sets T1 progression pattern |
| T2 Rep Scheme | Default / Modified | Sets T2 progression pattern |
| T3 Rep Scheme | Default / Modified | Sets T3 progression pattern |
| AMRAP Weight Increases | Default / Modified | Determines increment logic |

### Rep Scheme Options

| Tier | Default Scheme | Modified Scheme |
|------|----------------|-----------------|
| T1 | 5x3+ > 6x2+ > 10x1+ | 3x5+ > 4x3+ > 5x2+ |
| T2 | 3x10 > 3x8 > 3x6 | 4x8 > 4x6 > 4x4 |
| T3 | 3x15+ (25 rep progression) | 4x12+ (18 rep progression) |

---

## Starting Weights

### Main Lifts (T1)

Users enter starting weights for four main lifts:
- **Bench Press**
- **Squat**
- **Overhead Press (OHP)**
- **Deadlift**

**Recommendation**: Use 85% of current 5-rep max (5RM) as starting weight.

### T2 Main Lift Weights

T2 main lift weights are calculated as a percentage of the corresponding T1 starting weight:

```
T2 Weight = T1 Starting Weight x T2 Percentage
```

The T2 percentage is user-configurable (default suggestion: 35-80%, with 35-50% recommended for beginners).

### Weight Increments

Users configure increment values for each lift category:

| Lift Type | Typical Increment (Lbs) | Typical Increment (Kg) |
|-----------|-------------------------|------------------------|
| Bench | 2.5 | 2.5 |
| Squat | 5.0 | 5.0 |
| OHP | 2.5 | 2.5 |
| Deadlift | 5.0 | 5.0 |
| T2 Lifts | 2.5 | 2.5 |
| T3 Lifts | 5.0 | 5.0 |

---

## Tier 1 (T1) Progression Algorithm

### Structure

T1 exercises are the primary strength movements, performed with high intensity (heavy weight) and lower volume.

### Default Rep Scheme Progression

**Stage 1: 5 sets of 3 reps, last set AMRAP (5x3+)**
- Base volume target: 15 reps
- Weight increases each successful session

**Stage 2: 6 sets of 2 reps, last set AMRAP (6x2+)**
- Base volume target: 12 reps
- Triggered when Stage 1 fails
- Same weight as failed Stage 1 attempt

**Stage 3: 10 sets of 1 rep, last set AMRAP (10x1+)**
- Base volume target: 10 reps
- Triggered when Stage 2 fails
- Same weight as failed Stage 2 attempt

### Modified Rep Scheme Progression

**Stage 1: 3 sets of 5 reps, last set AMRAP (3x5+)**
- Base volume target: 15 reps

**Stage 2: 4 sets of 3 reps, last set AMRAP (4x3+)**
- Base volume target: 12 reps

**Stage 3: 5 sets of 2 reps, last set AMRAP (5x2+)**
- Base volume target: 10 reps

### Weight Progression Logic

#### Success Condition
Total reps completed across all sets must meet or exceed the base volume target.

#### Default AMRAP Increment Mode
When base volume is achieved, weight increases by the configured increment for that lift.

```
New Weight = Previous Weight + Configured Increment
```

#### Modified AMRAP Increment Mode
Weight increase varies based on AMRAP performance:

**For 5x3+ (Default scheme):**
| AMRAP Reps | Increment Tier |
|------------|----------------|
| 3-4 reps | Tier 1 (lowest) |
| 5-7 reps | Tier 2 (middle) |
| 8+ reps | Tier 3 (highest) |

**For 3x5+ (Modified scheme):**
| AMRAP Reps | Increment Tier |
|------------|----------------|
| 5-6 reps | Tier 1 (lowest) |
| 7-8 reps | Tier 2 (middle) |
| 9+ reps | Tier 3 (highest) |

### Failure Protocol

**Definition of Failure**: Unable to achieve base volume (not technical failure on a rep, but anticipated inability to complete 1-2 more quality reps).

**Progression through failure:**
1. Fail at 5x3+ (or 3x5+) -> Move to 6x2+ (or 4x3+) at SAME weight
2. Fail at 6x2+ (or 4x3+) -> Move to 10x1+ (or 5x2+) at SAME weight
3. Fail at 10x1+ (or 5x2+) -> Rest 2-3 days, then RETEST new 5RM

After retest, the program resets to Stage 1 with the new tested weight.

---

## Tier 2 (T2) Progression Algorithm

### Structure

T2 exercises support T1 lifts with moderate intensity and volume. Primary T2 lifts mirror T1 lifts (e.g., if T1 is Squat, T2 could be Bench).

### Default Rep Scheme Progression

**Stage 1: 3 sets of 10 reps (3x10)**
- Base volume target: 30 reps

**Stage 2: 3 sets of 8 reps (3x8)**
- Base volume target: 24 reps

**Stage 3: 3 sets of 6 reps (3x6)**
- Base volume target: 18 reps

### Modified Rep Scheme Progression

**Stage 1: 4 sets of 8 reps (4x8)**
- Base volume target: 32 reps

**Stage 2: 4 sets of 6 reps (4x6)**
- Base volume target: 24 reps

**Stage 3: 4 sets of 4 reps (4x4)**
- Base volume target: 16 reps

### Weight Progression Logic

#### Success Condition
Total reps across all sets meets or exceeds base volume target.

```
If Total Reps >= Base Volume:
    New Weight = Previous Weight + T2 Increment
    Rep Scheme = Same (no change)
```

#### Failure Handling
```
If Total Reps < Base Volume:
    Weight = Same
    Rep Scheme = Next Stage (lower reps per set)
```

### Reset Protocol

When failing at the final stage (3x6 or 4x4):
1. Mark the lift for "RESET"
2. User enters a new starting weight
3. Program returns to Stage 1 (3x10 or 4x8) at the new weight

**Recommendation**: New reset weight should be slightly higher than original starting weight.

---

## Tier 3 (T3) Progression Algorithm

### Structure

T3 exercises are accessory/isolation movements performed with low weight and high volume.

### Default Rep Scheme: 3x15+

- 3 sets of 15 reps minimum
- Final set is AMRAP
- **Progression threshold**: 25 total reps on final (AMRAP) set

### Modified Rep Scheme: 4x12+

- 4 sets of 12 reps minimum
- Final set is AMRAP
- **Progression threshold**: 18 total reps on final (AMRAP) set

### Weight Progression Logic

```
If Final AMRAP Set >= Progression Threshold:
    New Weight = Previous Weight + T3 Increment
Else:
    Weight = Same
```

The rep scheme never changes for T3 - only weight increases.

---

## Warmup Calculations

The program calculates three warmup sets based on working weight:

| Warmup Set | Percentage of Working Weight |
|------------|------------------------------|
| Set 1 | 60% |
| Set 2 | 70% |
| Set 3 | 85% |

### Rounding Logic

Warmup weights are rounded based on equipment settings:

**With Microplates (Yes):**
- Lbs: Round to nearest 0.5 lb
- Kg: Round to nearest 2.5 kg

**Without Microplates (No):**
- Lbs: Round to nearest 5 lbs
- Kg: Round to nearest 2.5 kg

---

## Estimated 1RM Calculation

The program calculates estimated 1-rep max using the Epley formula:

```
Estimated 1RM = Weight x (1 + (Reps / 30))
```

Or equivalently:
```
Estimated 1RM = Weight x Reps x 0.0333 + Weight
```

This uses the maximum reps achieved in any set of the T1 exercise.

---

## Wilks Score Calculation

The program calculates estimated Wilks score for powerlifting comparison purposes.

### Formula

The Wilks coefficient uses polynomial equations based on body weight:

**For XY (Male):**
```
Coefficient = 500 / (-216.0475144 + 16.2606339x - 0.002388645x^2 - 0.00113732x^3 + 0.00000701863x^4 - 0.00000001291x^5)
```

**For XX (Female):**
```
Coefficient = 500 / (594.31747775582 - 27.23842536447x + 0.82112226871x^2 - 0.00930733913x^3 + 0.00004731582x^4 - 0.00000009054x^5)
```

Where x = body weight in kg

```
Wilks Score = (Sum of Estimated 1RMs) x Coefficient
```

---

## 3-Day Program Structure

### Week Layout

Each week contains 3 training days with the following lift rotation:

**Day 1:**
- T1: Squat
- T2: Bench (at T2 percentage)
- T3: User selected

**Day 2:**
- T1: Bench
- T2: Deadlift (at T2 percentage)
- T3: User selected

**Day 3:**
- T1: Deadlift
- T2: OHP (at T2 percentage)
- T3: User selected

The fourth main lift (OHP as T1, Squat as T2) cycles in on subsequent training blocks.

### Lift Rotation Pattern

The program alternates which lift appears as T1, ensuring each main lift gets trained as both T1 (high intensity) and T2 (moderate intensity) throughout the training cycle.

---

## 4-Day Program Structure

### Week Layout

**Day 1:**
- T1: Squat
- T2: Bench (at T2 percentage)
- T3: User selected

**Day 2:**
- T1: OHP
- T2: Deadlift (at T2 percentage)
- T3: User selected

**Day 3:**
- T1: Bench
- T2: Squat (at T2 percentage)
- T3: User selected

**Day 4:**
- T1: Deadlift
- T2: OHP (at T2 percentage)
- T3: User selected

Each main lift appears once as T1 and once as T2 per week.

---

## Optional Additional Exercises

### Optional T2 Exercises

Users can add a second T2 exercise per day. Options include:
- Legs Up Bench
- Close Grip Bench Press
- Paused Bench
- Incline Bench
- Front Squat
- Barbell Row
- Sumo Deadlift
- Custom exercises (3 slots)

### T3 Exercise Options

Users can add up to 3 T3 exercises per day. Pre-configured options include:
- Bicep variations (Curls, DB Curls, Hammer Curls, EZ Bar Curls)
- Back exercises (Cable Row, Lat Pulldowns, DB Row, Chest-Supported Row)
- Tricep exercises (Skull Crushers, Tricep Pushdowns, Dips)
- Leg accessories (Leg Press, Lunges, Leg Curls, Leg Extensions, Calf Raises)
- Shoulder work (Lateral Raises, Face Pulls, Rear Delt Flyes)
- Core work (Cable Crunch)
- Custom exercises (3 slots)

---

## Rest Periods

| Exercise Tier | Recommended Rest |
|--------------|------------------|
| T1 | 3-5 minutes (heavy sets) |
| T2 | 120 seconds |
| T3 | 90 seconds |

---

## Bar Loading Calculator

The program includes a plate calculator feature that:

1. Takes the target weight
2. Subtracts the bar weight (45 lbs or 20 kg by default)
3. Calculates optimal plate combinations per side

### Available Plates (Per Side)

**Lbs System:**
55, 45, 35, 25, 15, 10, 5, 2.5, 1.5, 1, 0.75, 0.5, 0.25 lbs

**Kg System:**
25, 20, 15, 10, 5, 2.5, 2, 1.5, 1.25, 1, 0.5, 0.25 kg

Users can configure which plates they have available for accurate loading suggestions.

---

## Week-to-Week Progression Summary

### T1 Progression Flow
```
Success at Stage -> Add weight, stay at stage
Failure at 5x3+/3x5+ -> Keep weight, move to 6x2+/4x3+
Failure at 6x2+/4x3+ -> Keep weight, move to 10x1+/5x2+
Failure at 10x1+/5x2+ -> Retest 5RM, reset to Stage 1
```

### T2 Progression Flow
```
Success at Stage -> Add weight, stay at stage
Failure at 3x10/4x8 -> Keep weight, move to 3x8/4x6
Failure at 3x8/4x6 -> Keep weight, move to 3x6/4x4
Failure at 3x6/4x4 -> Reset with new weight, back to 3x10/4x8
```

### T3 Progression Flow
```
AMRAP meets threshold -> Add weight
AMRAP below threshold -> Keep weight
```

---

## Total Volume Tracking

The program tracks total volume lifted per exercise category:

```
Volume = Weight x Total Reps
```

Tracked metrics include:
- Individual lift totals (Bench, Squat, OHP, Deadlift)
- T2-1 Total (primary T2 lifts)
- T2-2 Total (optional T2 lifts)
- T3-1, T3-2, T3-3 Totals
- Combined total weight lifted per training block

---

## Key Program Principles

1. **Linear Progression**: Add weight every successful session
2. **Autoregulation**: AMRAP sets provide feedback for progression decisions
3. **Staged Deloading**: Rather than reducing weight, reduce reps per set when failing
4. **Skill Development**: High frequency on main lifts (each trained 2x per week minimum)
5. **Recovery Management**: Lower tiers provide active recovery while building volume
6. **Flexible Accessories**: T3 work customized to individual weaknesses

---

## References

- Original GZCL Method: https://swoleateveryheight.blogspot.com/2016/02/gzcl-applications-adaptations.html
- GZCLP Infographic: https://saynotobroscience.com/gzclp-infographic/
- Reddit Community: /r/gzcl

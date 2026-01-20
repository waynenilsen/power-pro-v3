# Greg Nuckols High Frequency Program Specification

## Overview

This program implements a high-frequency training approach for the main compound lifts (Bench Press, Squat, Deadlift, Overhead Press) along with accessory back work (Rows, Pulldowns). The program operates on a 3-week mesocycle with progressive volume increases while maintaining consistent intensity percentages.

## User Inputs

The program requires the following user-defined maximums:

### 1RM (One Rep Max) Inputs
- **Bench Press 1RM**
- **Overhead Press (OHP) 1RM**
- **Squat 1RM**
- **Deadlift 1RM**

### 6RM (Six Rep Max) Inputs
- **Row 6RM** (used for T-Bar Rows)
- **Pulldown 6RM** (used for Pulldowns)

### Rounding Factor
- A configurable rounding factor for weight calculations
- Use **2.5** for kilograms
- Use **5** for pounds
- All calculated weights are rounded to the nearest multiple of this factor

---

## Weekly Structure

The program runs 6 training days per week with the following daily focus:

| Day | Primary Focus | Secondary Focus |
|-----|--------------|-----------------|
| Monday | Bench + Squat | Triceps, Quadriceps, Conditioning |
| Tuesday | Bench + Squat | Pecs, Posterior Chain, Rows, Biceps |
| Wednesday | Bench + Squat | Triceps, Quadriceps, Conditioning |
| Thursday | Bench + Squat + Pulldowns | Triceps, Quadriceps, Biceps |
| Friday | OHP + Deadlift | Side Delts, Posterior Chain, Conditioning |
| Saturday | Bench + Squat + Rows | Pecs, Posterior Chain, Biceps |

**Training Frequency:**
- Bench Press: 5x per week
- Squat: 5x per week
- Deadlift: 1x per week
- OHP: 1x per week
- Rows: 2x per week (Tuesday + Saturday)
- Pulldowns: 1x per week (Thursday)

---

## Intensity Distribution by Day

### Main Lifts (Bench & Squat)

Each training day uses a specific intensity percentage applied to the 1RM:

| Day | Intensity | Character |
|-----|-----------|-----------|
| Monday | **75% of 1RM** | Moderate - Technical practice |
| Tuesday | **80% of 1RM** | Moderate-Heavy - Strength building |
| Wednesday | **70% of 1RM** | Light - Recovery/Volume |
| Thursday | **85% of 1RM** | Heavy - Intensity day |
| Saturday | **65% of 1RM** | Light - High volume hypertrophy |

### OHP and Deadlift
- Both train at **85% of 1RM** on Friday

### Back Exercises (Rows & Pulldowns)

Back exercises use a 3-week intensity progression based on the 6RM:

| Week | Intensity | Sets | Reps |
|------|-----------|------|------|
| Week 1 | **77.5% of 6RM** | 4 | 12 |
| Week 2 | **85% of 6RM** | 4 | 8 |
| Week 3 | **100% of 6RM** | 1+ sets | 6 |

---

## Volume Progression (3-Week Cycle)

### Sets Progression for Main Lifts

The program progressively increases sets across the 3-week cycle:

#### Monday (75% intensity)
| Week | Sets | Reps |
|------|------|------|
| Week 1 | 4 | 3 |
| Week 2 | 5 | 3 |
| Week 3 | 5 | 4 |

#### Tuesday (80% intensity)
| Week | Sets | Reps |
|------|------|------|
| Week 1 | 3 | 2 |
| Week 2 | 4 | 2 |
| Week 3 | 4 | 3 |

#### Wednesday (70% intensity)
| Week | Sets | Reps |
|------|------|------|
| Week 1 | 4 | 4 |
| Week 2 | 5 | 4 |
| Week 3 | 5 | 5 |

#### Thursday (85% intensity)
| Week | Sets | Reps |
|------|------|------|
| Week 1 | 3 | 1 |
| Week 2 | 4 | 1 |
| Week 3 | 1 | AMAP (As Many As Possible) |

#### Friday - OHP & Deadlift (85% intensity)
| Week | Sets | Reps |
|------|------|------|
| Week 1 | 3 | 1 |
| Week 2 | 4 | 1 |
| Week 3 | 1 | AMAP |

#### Saturday (65% intensity)
| Week | Sets | Reps |
|------|------|------|
| Week 1 | 5 | 5 |
| Week 2 | 6 | 5 |
| Week 3 | 6 | 6 |

---

## Progression Model

### AMAP Sets (Week 3 Test Days)

Week 3 includes AMAP (As Many As Possible) sets for the main lifts at 85% intensity. Based on performance, adjust the 1RM as follows:

| Reps Achieved | 1RM Adjustment |
|---------------|----------------|
| 6-7 reps | Increase by 5 lbs / 2.5 kg |
| 8-9 reps | Increase by 10 lbs / 5 kg |
| 10+ reps | Increase by 15 lbs / 7.5 kg |

### 1+ Set Protocol (Back Exercises)

Week 3 back exercises use a "1+" set notation. This means:
1. Perform the first set with the prescribed weight and reps
2. For every additional set completed, increase your 6RM by 2%
3. Continue performing sets until you can no longer complete the prescribed reps

This creates an autoregulated volume approach where stronger performance leads to greater training max increases.

---

## Exercise Categories and Rotation

### Primary Compound Lifts
These are fixed and form the core of the program:
- Bench Press
- Squat
- Deadlift
- Overhead Press (OHP)

### Back Work
Two exercise slots with prescribed movements:
- **Pulldowns** (Thursday)
- **T-Bar Rows** (Saturday)
- **Kroc Rows** (Tuesday) - RPE-based, not percentage-based

### Accessory Work (RPE-Based)

Accessory exercises are performed at prescribed RPE (Rate of Perceived Exertion) with unspecified reps ("?"). The lifter selects appropriate weight to hit the target RPE within a reasonable rep range.

| Category | RPE | Sets | Frequency |
|----------|-----|------|-----------|
| Tricep Work | RPE 8 | 2 | Mon, Wed, Thu |
| Pec Work | RPE 8 | 2 | Tue, Sat |
| Quadricep Work | RPE 8 | 2 | Mon, Wed, Thu |
| Posterior Chain Work | RPE 8 | 2 | Tue, Thu, Fri, Sat |
| Side Delt Work | RPE 8 | 2 | Fri |
| Biceps | RPE 9 | 3 | Tue, Thu, Sat |
| Kroc Rows | RPE 9-9.5 | 4 | Tue |

### Exercise Substitution Guidelines

Accessory exercises can be substituted based on individual weak points:
- "Tricep Work" on Bench day can be substituted with "Pec Work"
- The specific exercises performed are up to the lifter as long as they isolate the target muscle group

### Additional Recommendations
- After every workout, perform 2 sets of:
  - Ab work (medium intensity)
  - Rear delt work (medium intensity)

---

## Conditioning

Conditioning is included on:
- Monday
- Wednesday
- Friday

The program does not prescribe specific conditioning protocols, allowing the lifter to choose appropriate modalities.

---

## Weight Calculation Formula

All main lift weights are calculated using:

```
Working Weight = ROUND(1RM * Intensity%, Rounding Factor)
```

Where:
- **1RM** is the user's one-rep maximum for the lift
- **Intensity%** is the day-specific percentage (65%, 70%, 75%, 80%, or 85%)
- **Rounding Factor** is 2.5 (kg) or 5 (lbs)
- **ROUND** rounds to the nearest multiple of the rounding factor

For back exercises using 6RM:
```
Working Weight = ROUND(6RM * Intensity%, Rounding Factor)
```

---

## Summary of Key Program Principles

1. **High Frequency**: Train bench and squat 5x per week to maximize skill acquisition and muscle protein synthesis
2. **Daily Undulating Periodization**: Vary intensity daily (65-85%) to manage fatigue while providing varied stimuli
3. **Progressive Overload**: Increase volume weekly within the 3-week cycle
4. **Autoregulation**: Week 3 AMAP sets and 1+ protocols allow performance-based progression
5. **Fixed Intensity/Variable Volume**: Percentages stay constant; sets and reps increase
6. **Balanced Upper/Lower**: Each session includes both upper (bench) and lower (squat) work
7. **Accessory Autoregulation**: RPE-based accessory work allows for daily readiness adjustments
8. **Specialization Flexibility**: Accessory exercises can be modified based on individual needs

---

## Cycle Reset

After completing Week 3:
1. Update all 1RM values based on AMAP performance
2. Update 6RM values for back exercises based on 1+ set performance (2% increase per additional set completed)
3. Return to Week 1 with the new maxes
4. Repeat the 3-week cycle

This creates a continuous progression system where each mesocycle builds upon the previous one's performance gains.

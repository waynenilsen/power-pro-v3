# nSuns 5/3/1 Linear Progression - 5-Day Variant

## Program Overview

nSuns 5/3/1 LP is a high-volume, linear progression strength training program based on Jim Wendler's 5/3/1 methodology. It combines the periodization principles of 5/3/1 with aggressive weekly progression suitable for intermediate lifters.

### Philosophy

The program is built around several core principles:

1. **Training Max (TM) Based Loading**: All working weights are calculated as percentages of your Training Max, not your true 1 Rep Max. This builds in fatigue management.

2. **High Volume Primary Work**: Each primary lift is performed for 9 sets with varying intensities and rep ranges, creating a pyramid-style loading pattern.

3. **Complementary T2 Movements**: Secondary lifts (Tier 2) support the primary movements with lighter weights and higher cumulative volume.

4. **Weekly Linear Progression**: Training maxes increase every week based on AMRAP performance, allowing rapid strength gains for those who can recover.

5. **Autoregulation via AMRAP Sets**: The 1+ (one-or-more) set on primary lifts determines progression rate, providing built-in autoregulation.

---

## Input Parameters

### User Inputs

The program requires the following inputs from the user:

1. **1 Rep Maxes (1RMs)** for each of the five main lifts:
   - Squat
   - Bench Press
   - Deadlift
   - Overhead Press (OHP)
   - Barbell Row

2. **Unit Selection**: Choose between pounds (Lb) or kilograms (Kg)

### Derived Values

**Training Maxes (TMs)**: Initially calculated as 90% of the user's true 1RM. After the first week, users manually update their TMs based on AMRAP performance (the original formula is intentionally overwritten).

```
Initial TM = 1RM * 0.90
```

**Rounding Increment**: Weights are rounded to the nearest practical plate increment based on unit selection:
- Pounds: Round to nearest 5 lb
- Kilograms: Round to nearest 2.5 kg

---

## Weight Calculation Formula

All working weights are calculated using the following formula:

```
Working Weight = ROUND(Training Max * Percentage, Rounding Increment)
```

Where:
- `Training Max` is the TM for the relevant lift
- `Percentage` is the prescribed intensity for that set
- `Rounding Increment` depends on unit selection (5 lb or 2.5 kg)

---

## Weekly Progression Logic

Progression is determined by performance on the 1+ AMRAP set (typically Set 3 of primary lifts):

| Reps Achieved on 1+ Set | TM Increase |
|------------------------|-------------|
| 0-1 reps | No increase |
| 2-3 reps | +5 lb (2.5 kg) |
| 4-5 reps | +5-10 lb (2.5-5 kg) |
| 5+ reps | +10-15 lb (5-7.5 kg) |

**Note**: Users manually update their Training Max values each week. The program does not auto-calculate progression.

---

## Program Structure

### Training Split

| Day | Primary Lift | Secondary/T2 Lifts | Assistance Focus |
|-----|-------------|-------------------|------------------|
| Monday | Bench Press | Row (mirror pattern), OHP | Chest, Arms, Back |
| Tuesday | Squat | Deficit Deadlift | Legs, Abs |
| Wednesday | Overhead Press | Row, Pause Bench | Shoulders, Chest |
| Thursday | Deadlift | Front Squat | Back, Abs |
| Friday | Bench Press | Row (primary pattern), Incline Close-Grip Bench | Arms, Other |

---

## Set and Rep Schemes

### Primary Lift Patterns

The program uses two distinct primary lift patterns:

#### Pattern A: Volume Day (Monday Bench)
Used for the first bench session of the week with emphasis on accumulated volume at moderate intensity.

| Set | % of TM | Reps |
|-----|---------|------|
| 1 | 65% | 8 |
| 2 | 75% | 6 |
| 3 | 85% | 4 |
| 4 | 85% | 4 |
| 5 | 85% | 4 |
| 6 | 80% | 5 |
| 7 | 75% | 6 |
| 8 | 70% | 7 |
| 9 | 65% | 8 |

This pattern pyramids up to 85%, holds for three sets, then pyramids back down with increasing reps.

#### Pattern B: Intensity Day (Standard 5/3/1+)
Used for Squat, OHP, Deadlift, and Friday Bench - features the classic 5/3/1 top set structure.

| Set | % of TM | Reps |
|-----|---------|------|
| 1 | 75% | 5 |
| 2 | 85% | 3 |
| 3 | 95% | 1+ (AMRAP) |
| 4 | 90% | 3 |
| 5 | 85% | 3 |
| 6 | 80% | 3 |
| 7 | 75% | 5 |
| 8 | 70% | 5 |
| 9 | 65% | 5 |

**Note on Set 3 (1+ Set)**: This is an AMRAP (As Many Reps As Possible) set. Perform as many reps as possible with good form. This performance drives weekly progression decisions.

#### Pattern B Variant: Deadlift
Deadlifts use the same percentage pattern but with 3 reps on all back-off sets:

| Set | % of TM | Reps |
|-----|---------|------|
| 1 | 75% | 5 |
| 2 | 85% | 3 |
| 3 | 95% | 1+ (AMRAP) |
| 4 | 90% | 3 |
| 5 | 85% | 3 |
| 6 | 80% | 3 |
| 7 | 75% | 3 |
| 8 | 70% | 3 |
| 9 | 65% | 3 |

#### Pattern B Variant: Friday Bench
Friday Bench uses alternating 5/3 rep patterns on back-off sets:

| Set | % of TM | Reps |
|-----|---------|------|
| 1 | 75% | 5 |
| 2 | 85% | 3 |
| 3 | 95% | 1+ (AMRAP) |
| 4 | 90% | 3 |
| 5 | 85% | 5 |
| 6 | 80% | 3 |
| 7 | 75% | 5 |
| 8 | 70% | 3 |
| 9 | 65% | 5 |

---

### Row Patterns

Rows follow the corresponding bench pattern for that day, using the Row Training Max:

#### Monday Row (Matches Monday Bench - Volume Pattern)
Uses the same Volume Day percentages: 65/75/85/85/85/80/75/70/65% with reps 8/6/4/4/4/5/6/7/8.

#### Wednesday Row (Light T2 Pattern)
Uses reduced percentages as a lighter secondary movement.

#### Friday Row (Matches Friday Bench - Intensity Pattern)
Uses the same percentages as Friday Bench (75/85/95/90/85/80/75/70/65%) with 5/3/1+/3/5/3/5/3/5+ rep scheme, including an AMRAP on the final set.

---

### T2 (Tier 2) Movement Patterns

T2 exercises use a "Pyramid and Hold" structure - weights increase through the first three sets, then the heaviest weight is maintained for the remaining sets with higher reps.

#### T2 Pattern Structure

| Set | Percentage Pattern | Reps |
|-----|-------------------|------|
| 1 | Low % | 5-6 |
| 2 | Medium % | 5 |
| 3 | High % | 3 |
| 4 | = Set 3 weight | 5 |
| 5 | = Set 3 weight | 7 |
| 6 | = Set 3 weight | 4 |
| 7 | = Set 3 weight | 6 |
| 8 | = Set 3 weight | 8 |

**Total**: 8 sets for T2 movements (not 9 like primary lifts)

#### Specific T2 Percentages by Exercise

**Monday OHP (T2)** - Uses OHP Training Max:
- Sets 1-3: 45% / 55% / 65%
- Sets 4-8: Hold at 65% weight

**Tuesday Deficit Deadlift** - Uses Deadlift Training Max:
- Sets 1-3: 50% / 60% / 70%
- Sets 4-8: Hold at 70% weight

**Wednesday Row (T2)** - Uses Row Training Max:
- Sets 1-3: 40% / 50% / 60%
- Sets 4-8: Hold at 60% weight

**Wednesday Pause Bench** - Uses Bench Training Max:
- Sets 1-3: 35% / 45% / 55%
- Sets 4-8: Hold at 55% weight

**Thursday Front Squat** - Uses Squat Training Max:
- Sets 1-3: 35% / 45% / 55%
- Sets 4-8: Hold at 55% weight

**Friday Incline Close-Grip Bench** - Uses Bench Training Max:
- Sets 1-3: 40% / 50% / 60%
- Sets 4-8: Hold at 60% weight

---

## Complete Daily Breakdown

### Monday

**Bench Press** (Primary - Volume Pattern)
| Set | % of Bench TM | Reps |
|-----|---------------|------|
| 1 | 65% | 8 |
| 2 | 75% | 6 |
| 3 | 85% | 4 |
| 4 | 85% | 4 |
| 5 | 85% | 4 |
| 6 | 80% | 5 |
| 7 | 75% | 6 |
| 8 | 70% | 7 |
| 9 | 65% | 8 |

**Barbell Row** (Secondary - Volume Pattern)
| Set | % of Row TM | Reps |
|-----|-------------|------|
| 1 | 65% | 8 |
| 2 | 75% | 6 |
| 3 | 85% | 4 |
| 4 | 85% | 4 |
| 5 | 85% | 4 |
| 6 | 80% | 5 |
| 7 | 75% | 6 |
| 8 | 70% | 7 |
| 9 | 65% | 8 |

**Overhead Press** (T2)
| Set | % of OHP TM | Reps |
|-----|-------------|------|
| 1 | 45% | 6 |
| 2 | 55% | 5 |
| 3 | 65% | 3 |
| 4 | 65% | 5 |
| 5 | 65% | 7 |
| 6 | 65% | 4 |
| 7 | 65% | 6 |
| 8 | 65% | 8 |

**Assistance**: Chest, Arms, Back (bodybuilding rep ranges)

---

### Tuesday

**Squat** (Primary - Intensity Pattern)
| Set | % of Squat TM | Reps |
|-----|---------------|------|
| 1 | 75% | 5 |
| 2 | 85% | 3 |
| 3 | 95% | 1+ |
| 4 | 90% | 3 |
| 5 | 85% | 3 |
| 6 | 80% | 3 |
| 7 | 75% | 5 |
| 8 | 70% | 5 |
| 9 | 65% | 5 |

**Deficit Deadlift** (T2)
| Set | % of Deadlift TM | Reps |
|-----|------------------|------|
| 1 | 50% | 5 |
| 2 | 60% | 5 |
| 3 | 70% | 3 |
| 4 | 70% | 5 |
| 5 | 70% | 7 |
| 6 | 70% | 4 |
| 7 | 70% | 6 |
| 8 | 70% | 8 |

**Assistance**: Legs, Abs

---

### Wednesday

**Overhead Press** (Primary - Intensity Pattern)
| Set | % of OHP TM | Reps |
|-----|-------------|------|
| 1 | 75% | 5 |
| 2 | 85% | 3 |
| 3 | 95% | 1+ |
| 4 | 90% | 3 |
| 5 | 85% | 3 |
| 6 | 80% | 3 |
| 7 | 75% | 5 |
| 8 | 70% | 5 |
| 9 | 65% | 5 |

**Barbell Row** (T2)
| Set | % of Row TM | Reps |
|-----|-------------|------|
| 1 | 40% | 6 |
| 2 | 50% | 5 |
| 3 | 60% | 3 |
| 4 | 60% | 5 |
| 5 | 60% | 7 |
| 6 | 60% | 4 |
| 7 | 60% | 6 |
| 8 | 60% | 8 |

**Pause Bench** (T2)
| Set | % of Bench TM | Reps |
|-----|---------------|------|
| 1 | 35% | 6 |
| 2 | 45% | 5 |
| 3 | 55% | 3 |
| 4 | 55% | 5 |
| 5 | 55% | 7 |
| 6 | 55% | 4 |
| 7 | 55% | 6 |
| 8 | 55% | 8 |

**Assistance**: Shoulders, Chest

---

### Thursday

**Deadlift** (Primary - Intensity Pattern with 3-rep back-offs)
| Set | % of Deadlift TM | Reps |
|-----|------------------|------|
| 1 | 75% | 5 |
| 2 | 85% | 3 |
| 3 | 95% | 1+ |
| 4 | 90% | 3 |
| 5 | 85% | 3 |
| 6 | 80% | 3 |
| 7 | 75% | 3 |
| 8 | 70% | 3 |
| 9 | 65% | 3 |

**Front Squat** (T2)
| Set | % of Squat TM | Reps |
|-----|---------------|------|
| 1 | 35% | 5 |
| 2 | 45% | 5 |
| 3 | 55% | 3 |
| 4 | 55% | 5 |
| 5 | 55% | 7 |
| 6 | 55% | 4 |
| 7 | 55% | 6 |
| 8 | 55% | 8 |

**Assistance**: Back, Abs

---

### Friday

**Bench Press** (Primary - Intensity Pattern with alternating reps)
| Set | % of Bench TM | Reps |
|-----|---------------|------|
| 1 | 75% | 5 |
| 2 | 85% | 3 |
| 3 | 95% | 1+ |
| 4 | 90% | 3 |
| 5 | 85% | 5 |
| 6 | 80% | 3 |
| 7 | 75% | 5 |
| 8 | 70% | 3 |
| 9 | 65% | 5 |

**Barbell Row** (Primary Pattern - with AMRAP)
| Set | % of Row TM | Reps |
|-----|-------------|------|
| 1 | 75% | 5 |
| 2 | 85% | 3 |
| 3 | 95% | 1+ |
| 4 | 90% | 3 |
| 5 | 85% | 5 |
| 6 | 80% | 3 |
| 7 | 75% | 5 |
| 8 | 70% | 3 |
| 9 | 65% | 5+ |

**Note**: Friday Row includes an AMRAP on the final set (5+).

**Incline Close-Grip Bench** (T2)
| Set | % of Bench TM | Reps |
|-----|---------------|------|
| 1 | 40% | 6 |
| 2 | 50% | 5 |
| 3 | 60% | 3 |
| 4 | 60% | 5 |
| 5 | 60% | 7 |
| 6 | 60% | 4 |
| 7 | 60% | 6 |
| 8 | 60% | 8 |

**Assistance**: Arms, Other

---

## Assistance Work Guidelines

Assistance work is performed with bodybuilding-style sets and reps (typically 3-4 sets of 8-12 reps). The specific exercises can be customized based on individual needs.

### Suggested Assistance Categories

| Focus Area | Example Exercises |
|------------|------------------|
| Chest | Dumbbell flyes, cable crossovers, dips |
| Arms | Bicep curls, tricep pushdowns, skull crushers |
| Back | Lat pulldowns, cable rows, face pulls |
| Shoulders | Lateral raises, rear delt flyes, shrugs |
| Legs | Leg press, leg curls, leg extensions, lunges |
| Abs | Planks, hanging leg raises, cable crunches |

---

## Algorithm Summary

### To Calculate Any Working Weight:

1. Identify the lift's Training Max (TM)
2. Identify the prescribed percentage for that set
3. Multiply: `TM * Percentage`
4. Round to the nearest increment:
   - Pounds: Round to nearest 5
   - Kilograms: Round to nearest 2.5

### Weekly Workflow:

1. Complete all workouts for the week
2. Review AMRAP performance on each 1+ set
3. Determine TM increases based on the progression table
4. Update Training Max values for the following week
5. Repeat

### Key Implementation Notes:

- All weight calculations use the MROUND function (round to nearest multiple)
- T2 movements sets 4-8 reference the weight from set 3 (hold at top weight)
- The program originally calculates TMs as 90% of 1RM, but this formula is intentionally overwritten when users update their TMs weekly
- Row is treated as a primary lift with its own TM and appears with different patterns across the week

---

## Example Calculation

**Given:**
- Bench 1RM: 100 kg
- Bench TM: 85 kg (after user adjustment)
- Unit: Kilograms (2.5 kg rounding)

**Friday Bench Set 3 Calculation:**
```
Working Weight = MROUND(85 * 0.95, 2.5)
Working Weight = MROUND(80.75, 2.5)
Working Weight = 80 kg
```

**Monday Bench Set 1 Calculation:**
```
Working Weight = MROUND(85 * 0.65, 2.5)
Working Weight = MROUND(55.25, 2.5)
Working Weight = 55 kg
```

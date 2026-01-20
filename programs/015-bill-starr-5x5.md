# Bill Starr 5x5 Linear Program Specification

## Overview

The Bill Starr 5x5 is a classic strength training program designed for intermediate lifters. This implementation follows the "Madcow" linear version, featuring a Heavy/Light/Medium weekly structure with ramping sets that progressively increase in weight each set. The program uses calculated percentages to auto-generate training weights based on user-input maxes.

## User Input Parameters

### Required Inputs

The program requires the following user-configurable parameters for each lift:

1. **Current Max** - The lifter's current tested weight for a given rep range
2. **Reps** - The number of reps the Current Max was achieved at (must be less than 12)
3. **Weekly Increase** - The amount of weight to add to the top set each week
   - Specified in pounds for the "by Pound" version
   - Specified as percentage for the "by Percentage" version (typically 2.5%)
4. **Percent to Reset** - How far below the calculated 5RM to begin training (typically 7%, meaning old PRs are reached around Week 4)
5. **Smallest Weight Increment** - The minimum weight that can be added to the barbell (typically 5 lbs; used for rounding)

### Exercises Tracked

- **Squat** - Primary lower body movement
- **Bench Press** - Primary horizontal push
- **Row** (or Power Clean) - Primary pull movement
- **Incline Press** (or Military Press) - Secondary pressing movement
- **Deadlift** (or High Pulls) - Primary hip hinge movement

## Rep Max Conversion Formulas

### 1RM Calculation (Brzycki Formula)

To convert a tested rep max to an estimated one-rep max:

```
1RM = Current Max / (1.0278 - (0.0278 * Reps))
```

### 5RM Calculation

To calculate the estimated 5-rep max from the 1RM:

```
5RM = ROUND((1RM * (1.0278 - (0.0278 * 5))) / Smallest_Weight) * Smallest_Weight
```

This rounds to the nearest weight increment (e.g., nearest 5 lbs).

## Week 1 Starting Weight Calculation

The Week 1 top set weight is calculated by taking a percentage reduction from the calculated 5RM:

```
Week_1_Top_Set = ROUND((5RM - (5RM * Percent_to_Reset)) / Smallest_Weight) * Smallest_Weight
```

With a typical 7% reset, this places the lifter approximately 3-4 weeks below their current 5RM, allowing for PR attempts around Week 4.

## Weekly Progression Model

The program offers two progression models:

### Pound-Based Linear Progression (Primary Method)

Each week, the top set increases by a fixed pound amount:

```
Next_Week_Top_Set = Current_Week_Top_Set + Weekly_Pound_Increase
```

This is calculated by taking Friday's top set weight and adding the specified pound increase to determine the following Monday's top set.

### Percentage-Based Linear Progression (Alternative Method)

Each week, the top set increases by a fixed percentage:

```
Next_Week_Top_Set = ROUND(Current_Week_Top_Set * (1 + (Percentage_Increase / 100)))
```

## Weekly Training Structure

The program follows a three-day-per-week schedule with a Heavy/Light/Medium rotation:

### Monday - Heavy Day

**Exercise 1: Squat (5 ramping sets of 5)**
- Set 1: 50% of top set weight x 5
- Set 2: 63% of top set weight x 5
- Set 3: 75% of top set weight x 5
- Set 4: 88% of top set weight x 5
- Set 5: 100% (top set) x 5

**Exercise 2: Bench Press (5 ramping sets of 5)**
- Same ramping protocol as Squat

**Exercise 3: Row or Power Clean (5 ramping sets of 5)**
- Same ramping protocol as Squat

**Assistance Work:**
- Weighted Hyperextensions: 2 sets
- Weighted Sit-Ups: 4 sets

### Wednesday - Light Day

**Exercise 1: Squat (4 sets of 5, capped intensity)**
- Set 1: Same as Monday Set 1 (50% of Monday top set)
- Set 2: Same as Monday Set 2 (63% of Monday top set)
- Set 3: Same as Monday Set 3 (75% of Monday top set)
- Set 4: Same as Monday Set 3 (75% of Monday top set) - repeated, not ramped higher

Light Day squats stop at 75% of Monday's top set weight. The lifter performs 4 sets, with the final two sets using the same weight.

**Exercise 2: Incline Press or Military Press (4 ramping sets of 5)**
- Set 1: 63% of top set weight x 5
- Set 2: 75% of top set weight x 5
- Set 3: 88% of top set weight x 5
- Set 4: 100% (top set) x 5

Note: Incline/Military Press uses only 4 sets, starting at 63% rather than 50%.

**Exercise 3: Deadlift or High Pulls (4 ramping sets of 5)**
- Set 1: 63% of top set weight x 5
- Set 2: 75% of top set weight x 5
- Set 3: 88% of top set weight x 5
- Set 4: 100% (top set) x 5

Note: Deadlift/High Pulls uses the same 4-set protocol as Incline Press.

**Assistance Work:**
- Sit-Ups: 3 sets

### Friday - Medium Day

**Exercise 1: Squat (6 sets with mixed rep scheme)**
- Set 1: 50% of top set weight x 5
- Set 2: 63% of top set weight x 5
- Set 3: 75% of top set weight x 5
- Set 4: 88% of top set weight x 5
- Set 5: 100% (top set) x 3 (triple at new PR weight)
- Set 6: 75% of top set weight x 8 (back-off set)

The Friday top set is calculated as:
```
Friday_Top_Set = Monday_Top_Set + Weekly_Increase
```

This means Friday's triple is performed at next week's target 5RM weight, serving as an overreach/PR attempt.

**Exercise 2: Bench Press (6 sets with mixed rep scheme)**
- Same protocol as Friday Squat

**Exercise 3: Row or Power Clean (6 sets with mixed rep scheme)**
- Same protocol as Friday Squat

**Assistance Work:**
- Weighted Dips: 3 sets of 5-8
- Barbell Curls: 3 sets of 8
- Triceps Extensions: 3 sets of 8

## Ramping Set Percentages Summary

### Standard 5x5 Ramp (Monday Heavy, Main Lifts)
| Set | Percentage | Reps |
|-----|------------|------|
| 1   | 50%        | 5    |
| 2   | 63%        | 5    |
| 3   | 75%        | 5    |
| 4   | 88%        | 5    |
| 5   | 100%       | 5    |

### 4x5 Ramp (Wednesday Pressing/Pulling)
| Set | Percentage | Reps |
|-----|------------|------|
| 1   | 63%        | 5    |
| 2   | 75%        | 5    |
| 3   | 88%        | 5    |
| 4   | 100%       | 5    |

### Medium Day Protocol (Friday)
| Set | Percentage | Reps |
|-----|------------|------|
| 1   | 50%        | 5    |
| 2   | 63%        | 5    |
| 3   | 75%        | 5    |
| 4   | 88%        | 5    |
| 5   | 100%       | 3    |
| 6   | 75%        | 8    |

## Weight Rounding

All calculated weights are rounded down to the nearest weight increment using a floor function:

```
Ramping_Set_Weight = FLOOR(Percentage * Top_Set_Weight, Smallest_Weight)
```

This ensures all weights are achievable with available plates (e.g., multiples of 5 lbs).

## Periodization Logic

### Intra-Week Periodization

The Heavy/Light/Medium structure creates a stress-recovery-adaptation cycle within each week:

1. **Monday (Heavy)**: Maximum intensity for the week on all lifts
2. **Wednesday (Light)**: Reduced volume and intensity for squats; introduces accessory pressing and pulling work
3. **Friday (Medium)**: Sets new PRs with a triple, followed by back-off volume

### Inter-Week Periodization

The program uses simple linear periodization with weekly weight increases. The starting point is deliberately set below current capacity (via Percent to Reset) to allow for:

1. **Weeks 1-3**: Sub-maximal work, technique refinement, building momentum
2. **Week 4**: Typically matches previous PRs
3. **Weeks 5+**: New personal records each week

## Program Duration

The spreadsheet calculates 10-12 weeks of programming. Lifters typically run the program until:
- Progress stalls on multiple lifts
- A deload or reset is needed
- Transitioning to an intermediate periodization model

## Key Implementation Notes

### Exercise Substitutions

The program allows for the following substitutions:
- Row can be replaced with Power Clean
- Incline Press can be replaced with Military Press
- Deadlift can be replaced with High Pulls

### Progression Linking

- Friday's top set weight feeds into the following Monday's progression
- Monday's top set becomes the baseline for Wednesday's light work
- The cycle creates natural auto-regulation of intensity across the week

### Assistance Work Guidelines

Assistance exercises are prescribed with set ranges rather than strict progressions:
- Focus is on the main lifts
- Assistance work supports but does not drive the program
- Volume is moderate to allow recovery for main lift progression

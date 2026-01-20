# Wendler 5/3/1 Program Specification

## Program Overview

The Wendler 5/3/1 is a percentage-based strength training program designed by Jim Wendler. This specification documents the algorithm as implemented in the spreadsheet calculator, which provides automatic weight calculations, progress tracking, and stall management.

**Note**: While the file is named "BBB" (Boring But Big), this spreadsheet implements the core 5/3/1 main lift programming without the BBB supplemental template. BBB supplemental work (5 sets of 10 reps at 50-60% of training max) would be performed after the main lifts but is not calculated in this spreadsheet.

## Core Exercises

The program is built around four primary compound lifts:
1. Military Press (Overhead Press)
2. Deadlift
3. Bench Press
4. Squat

Each lift is trained once per week on its own day, resulting in a 4-day training week.

---

## Training Max Calculation

### Initial Setup

The program uses a **Training Max (TM)** rather than a true 1-rep max (1RM) for all weight calculations. The training max is calculated as:

```
Training Max = Estimated 1RM x 0.90 (90%)
```

### Estimating 1RM from Submaximal Lifts

The spreadsheet uses the **Brzycki Formula** to estimate 1RM from a weight lifted for multiple reps:

```
Estimated 1RM = (36 x Weight) / (37 - Reps)
```

This formula is used throughout the program for:
- Calculating initial training max from a rep test
- Estimating 1RM from AMRAP (As Many Reps As Possible) sets
- Generating rep goals for progressive overload

### Rounding

All weights are rounded to the nearest practical increment:
- **Pounds**: Round to nearest 5 lbs
- **Kilograms**: Round to nearest 2.5 kg

The rounding function ensures weights are achievable with standard gym equipment.

---

## Weekly Percentage Progression (5/3/1 Scheme)

Each training cycle consists of 4 weeks: three working weeks followed by a deload week.

### Week 1: "5s Week" (3x5)

| Set | Percentage of TM | Target Reps |
|-----|------------------|-------------|
| 1   | 65%              | 5           |
| 2   | 75%              | 5           |
| 3   | 85%              | 5+ (AMRAP)  |

### Week 2: "3s Week" (3x3)

| Set | Percentage of TM | Target Reps |
|-----|------------------|-------------|
| 1   | 70%              | 3           |
| 2   | 80%              | 3           |
| 3   | 90%              | 3+ (AMRAP)  |

### Week 3: "5/3/1 Week"

| Set | Percentage of TM | Target Reps |
|-----|------------------|-------------|
| 1   | 75%              | 5           |
| 2   | 85%              | 3           |
| 3   | 95%              | 1+ (AMRAP)  |

### Week 4: Deload Week

| Set | Percentage of TM | Target Reps |
|-----|------------------|-------------|
| 1   | 40%              | 5           |
| 2   | 50%              | 5           |
| 3   | 60%              | 5           |

The "+" notation indicates an AMRAP set where the lifter performs as many reps as possible while maintaining good form.

---

## Cycle Progression

### Training Max Increases

After each 4-week cycle, the training max increases automatically:

| Lift          | Increment per Cycle |
|---------------|---------------------|
| Military Press | +5 lbs (2.5 kg)    |
| Bench Press    | +5 lbs (2.5 kg)    |
| Deadlift       | +10 lbs (5 kg)     |
| Squat          | +10 lbs (5 kg)     |

Upper body lifts increase by smaller increments due to their slower strength progression compared to lower body lifts.

### Training Max Formula for Subsequent Cycles

```
New Training Max = Previous Training Max + Cycle Increment
```

For example, if Cycle 1 Squat TM = 300 lbs:
- Cycle 2 Squat TM = 300 + 10 = 310 lbs
- Cycle 3 Squat TM = 310 + 10 = 320 lbs
- And so on...

---

## Rep Goal System

The spreadsheet includes a sophisticated rep goal calculator to help lifters beat their previous estimated 1RM each week.

### Rep Goal Calculation

The rep goal tells the lifter how many reps they need to perform on the top AMRAP set to achieve a target 1RM. Using the inverse of the Brzycki formula:

```
Required Reps = 37 - (36 x Weight / Target 1RM)
```

### Week 1 Rep Goal Logic

The Week 1 rep goal is calculated to beat the current training max by a configurable increment (default: 5 lbs):

```
Target 1RM = Current Training Max + Rep Goal Increment
Required Reps = 37 - (36 x Week1_Top_Set_Weight / Target_1RM)
```

### Week 2 and Week 3 Rep Goal Logic

For subsequent weeks, the rep goal adapts based on performance:

1. **If previous week's estimated 1RM >= Training Max**:
   ```
   Target 1RM = Previous Week's Estimated 1RM + Rep Goal Increment
   ```

2. **If previous week's estimated 1RM < Training Max**:
   ```
   Target 1RM = Training Max + Rep Goal Increment
   ```

This ensures the lifter always has a meaningful target, even if they had a bad week.

### Rep Goal Increment Settings

The default rep goal increment is 5 lbs for all lifts, but this is configurable per exercise. As lifters become more advanced and weights get heavier, increasing this value may provide more accurate rep targets.

---

## Stall Management System

### Purpose

The stall adjustment feature allows lifters to reduce their training max when progress stalls, without losing historical workout data or requiring a complete spreadsheet reset.

### How It Works

1. **Trigger**: When a lifter determines they've stalled (consistently missing rep goals, form breakdown, etc.), they mark the "Stall" indicator for the specific lift.

2. **Backoff Amount**: A configurable weight reduction is applied (default: 10 lbs for all lifts).

3. **Application**: The reduction is applied to the **next cycle's** training max, not the current cycle. This allows completion of the current training block.

### Stall Adjustment Formula

```
If Stall Triggered:
    Next Cycle TM = Current Cycle TM - Backoff Amount
Else:
    Next Cycle TM = Current Cycle TM + Normal Increment
```

### Benefits

- Independent adjustment per lift (can stall on one lift while progressing on others)
- Preserves all historical training data
- Automatic recalculation of all future cycles
- Can be triggered multiple times throughout training history

---

## 1RM Estimation Tools

The spreadsheet includes three utility calculators:

### 1RM Calculator

Given weight and reps performed, estimates the 1RM:
```
1RM = (36 x Weight) / (37 - Reps)
```

### Rep Goal Calculator

Given weight and target 1RM, calculates required reps:
```
Reps Needed = 37 - (36 x Weight / Target_1RM)
```

### Weight Calculator

Given target 1RM and reps to be performed, calculates required weight:
```
Weight = Target_1RM / (36 / (37 - Reps))
```

All calculations round to the configured increment (5 lbs or 2.5 kg).

---

## Progress Tracking

### Automatic 1RM Logging

The spreadsheet automatically tracks estimated 1RMs for each lift across all training weeks. When a lifter enters their actual reps performed on AMRAP sets, the estimated 1RM is calculated and logged.

### Highest 1RM Display

For each exercise, the highest estimated 1RM achieved across all recorded weeks is displayed, providing a clear view of personal records.

### Week-by-Week Tracking

The progress tracker logs data for each working week (deload weeks are excluded from progress tracking):
- Cycle number
- Week number within the program
- Estimated 1RM achieved for each lift

---

## Program Structure Summary

### Single Cycle (4 Weeks)

| Week | Name       | Set 1  | Set 2  | Set 3   |
|------|------------|--------|--------|---------|
| 1    | 5s Week    | 65%x5  | 75%x5  | 85%x5+  |
| 2    | 3s Week    | 70%x3  | 80%x3  | 90%x3+  |
| 3    | 5/3/1 Week | 75%x5  | 85%x3  | 95%x1+  |
| 4    | Deload     | 40%x5  | 50%x5  | 60%x5   |

### Long-Term Progression

The spreadsheet supports tracking up to 13 cycles (52 weeks / 1 year) of continuous programming with automatic weight progression.

---

## Alternative: 3/5/1 Powerlifting Variant

The spreadsheet also includes a 3/5/1 variant designed for powerlifters, which reorders the weeks:

| Week | Name       | Set 1  | Set 2  | Set 3   | Additional    |
|------|------------|--------|--------|---------|---------------|
| 1    | 3s Week    | 70%x3  | 80%x3  | 90%x3+  | Singles x3    |
| 2    | 5s Week    | 65%x5  | 75%x5  | 85%x5   | -             |
| 3    | 5/3/1 Week | 75%x5  | 85%x3  | 95%x1+  | Singles x2    |
| 4    | Deload     | 40%x5  | 50%x5  | 60%x5   | -             |

The 3/5/1 order places the heaviest work at the beginning and end of the cycle, with the higher-volume 5s week in the middle for recovery.

---

## BBB Supplemental Work (Not Included in Spreadsheet)

Although named "BBB," the spreadsheet does not calculate Boring But Big supplemental sets. For reference, the standard BBB template is:

### After Main Lift Work

| Sets | Reps | Percentage of TM |
|------|------|------------------|
| 5    | 10   | 50-60%           |

The supplemental work uses the same exercise as the main lift (e.g., 5x10 squats after main squat work) or an opposing movement (e.g., 5x10 deadlifts after main squat work).

### BBB Percentage Progression Options

- **Beginner**: 50% of TM for all cycles
- **Intermediate**: 50% -> 55% -> 60% progression across cycles
- **Advanced**: 60% of TM, or various challenge templates (65-70%)

---

## Implementation Notes

### Weight Calculation Formula

For any percentage-based weight calculation:
```
Working Weight = ROUND(Training_Max x Percentage, Rounding_Increment)
```

### Formula Reference: Brzycki Equation

The core formula used throughout:
```
1RM = (36 x Weight) / (37 - Reps)
```

This formula is most accurate in the 1-10 rep range. Accuracy decreases for sets of 12+ reps.

### Input Requirements

To use the program, a lifter needs to provide:
1. A recent weight lifted for a known number of reps for each of the four lifts
2. Selected rounding preference (5 lbs or 2.5 kg)

All other calculations (training max, working weights, rep goals, progress tracking) are derived automatically.

---

## Credits and Attribution

This specification documents the algorithm from a spreadsheet created by Pazienxa (v2.1), designed as a companion calculator for Jim Wendler's 5/3/1 program. The 5/3/1 methodology and training philosophy are the intellectual property of Jim Wendler.

For complete program understanding, training philosophy, and assistance work recommendations, consult Jim Wendler's official 5/3/1 books.

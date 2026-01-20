# Starting Strength Program Specification

A comprehensive strength training program calculator based on Mark Rippetoe's Starting Strength methodology. This specification documents the complete algorithm for implementing a spreadsheet-based training log with automatic weight calculations.

## Table of Contents

1. [Overview](#overview)
2. [User Inputs](#user-inputs)
3. [Core Calculations](#core-calculations)
4. [Program Variants](#program-variants)
5. [Warmup Protocols](#warmup-protocols)
6. [Linear Progression Model](#linear-progression-model)
7. [Weight Rounding](#weight-rounding)
8. [Plate Loading Reference](#plate-loading-reference)

---

## Overview

The Starting Strength program is a novice strength training program designed around linear progression. The lifter performs compound barbell movements three times per week, adding weight to the bar every session. The spreadsheet automates weight calculations for both working sets and warmup sets across multiple weeks of training.

### Core Principles

- **Frequency**: 3 sessions per week (typically Monday/Wednesday/Friday or A/B/A rotation)
- **Progression**: Weight increases every workout (linear progression)
- **Compound Movements**: Squat, Deadlift, Bench Press, Overhead Press, Power Clean
- **Rep Schemes**: Primarily 3x5 for main lifts, 1x5 for deadlift, 5x3 for power cleans

---

## User Inputs

### Required Configuration Parameters

| Parameter | Description | Typical Default |
|-----------|-------------|-----------------|
| **Test Weight** | Starting weight or previous rep max for each exercise | User-defined |
| **Reps** | Number of reps performed at test weight (used for 1RM calculation) | 5 |
| **lb Increase** | Weight increment per workout for each exercise | Varies by lift |
| **% to Reset** | Percentage to subtract when starting below a previous max | 0% |
| **Smallest Weight Increment** | Minimum weight that can be added to the barbell | 5 lbs |

### Recommended Weight Increments by Exercise

| Exercise | Recommended Increment |
|----------|----------------------|
| Squat | 5 lbs per workout |
| Bench Press | 5 lbs per workout |
| Overhead Press | 5 lbs per workout |
| Deadlift | 10-15 lbs per workout |
| Power Clean | 5 lbs per workout |
| Front Squat | 5 lbs per workout |

---

## Core Calculations

### 1RM Estimation (Brzycki Formula)

The program uses the Brzycki formula to estimate a one-rep max from a rep max:

```
1RM = Weight / (1.0278 - (0.0278 * Reps))
```

Where:
- `Weight` = The test weight lifted
- `Reps` = Number of reps performed (must be less than 12)

### 5RM Calculation

The estimated 5-rep max is derived from the 1RM:

```
5RM = ROUND((1RM * (1.0278 - (0.0278 * 5))) / SmallestWeight) * SmallestWeight
```

This formula:
1. Converts the 1RM back to a 5RM using the Brzycki formula in reverse
2. Rounds the result to the nearest smallest weight increment

### Starting Working Set Weight

The initial working set weight accounts for any reset percentage:

```
Starting_Weight = ROUND((5RM - (5RM * Reset%)) / SmallestWeight) * SmallestWeight
```

This allows lifters returning from a layoff or deload to start below their previous max.

---

## Program Variants

The spreadsheet includes five distinct program variants, each with slightly different exercise selection and scheduling.

### 1. Original Novice Program

The classic A/B workout rotation.

**Workout A:**
- Squat: 3x5
- Bench Press: 3x5
- Deadlift: 1x5

**Workout B:**
- Squat: 3x5
- Overhead Press: 3x5
- Power Clean: 5x3

**Rotation**: A/B/A/B... (alternating each session)

### 2. Onus Wunsler Program

A variation that staggers deadlift and power clean frequency.

**Workout A:**
- Squat: 3x5
- Overhead Press: 3x5
- Deadlift: 1x5 (every other A workout)
- Power Clean: 5x3 (every other A workout, alternating with deadlift)

**Workout B:**
- Squat: 3x5
- Bench Press: 3x5
- Back Extensions: 3-5x10 (or Glute Ham Raises)
- Chin-Ups: 3 sets to failure

**Special Feature**: Deadlift and Power Clean appear on alternating Workout A sessions (odd sessions only for one, even sessions only for the other).

### 3. Practical Programming Novice Program

A three-day weekly structure with specific days assigned.

**Monday:**
- Squat: 3x5
- Bench Press: 3x5 (alternating weeks)
- Overhead Press: 5x3 (alternating weeks)
- Chin-Ups: 3 sets to failure

**Wednesday:**
- Squat: 3x5
- Bench Press: 3x5 (alternating weeks)
- Overhead Press: 5x3 (alternating weeks)
- Deadlift: 1x5

**Friday:**
- Squat: 3x5
- Bench Press: 3x5 (alternating weeks)
- Overhead Press: 5x3 (alternating weeks)
- Pull-Ups: 3 sets to failure

**Note**: Bench Press and Overhead Press alternate which sessions they appear on, creating a weekly rotation pattern.

### 4. Wichita Falls Novice Program

Similar to Practical Programming but includes Power Clean.

**Monday:**
- Squat: 3x5
- Bench Press: 3x5 (alternating)
- Overhead Press: 5x3 (alternating)
- Chin-Ups: 3 sets to failure

**Wednesday:**
- Squat: 3x5
- Bench Press: 3x5 (alternating)
- Overhead Press: 5x3 (alternating)
- Deadlift: 1x5 (alternating with Power Clean)
- Power Clean: 5x3 (alternating with Deadlift)

**Friday:**
- Squat: 3x5
- Bench Press: 3x5 (alternating)
- Overhead Press: 5x3 (alternating)
- Pull-Ups: 3 sets to failure

### 5. Advanced Novice Program

Introduces Front Squat as a light day variation and weighted chin-ups/pull-ups.

**Monday:**
- Squat: 3x5
- Bench Press: 3x5 (alternating)
- Overhead Press: 5x3 (alternating)
- Chin-Ups: 3 sets to failure (alternating weighted/unweighted)

**Wednesday:**
- Front Squat: 3x5 (replaces back squat as lighter variation)
- Bench Press: 3x5 (alternating)
- Overhead Press: 5x3 (alternating)
- Deadlift: 1x5 (alternating with Power Clean)
- Power Clean: 5x3 (alternating with Deadlift)

**Friday:**
- Squat: 3x5
- Bench Press: 3x5 (alternating)
- Overhead Press: 5x3 (alternating)
- Pull-Ups: 3 sets to failure (alternating weighted/unweighted)

---

## Warmup Protocols

Warmup sets are calculated as percentages of the working set weight, with all weights rounded down to the nearest plate-loadable weight (typically 5 lbs).

### Warmup Calculation Formula

```
Warmup_Weight = FLOOR(Working_Set_Weight * Percentage, SmallestWeight)
```

The FLOOR function rounds DOWN to the nearest increment, ensuring warmups are never heavier than intended.

### Protocol by Exercise Type

#### Squat (and Front Squat) Warmup Protocol

| Set | Sets x Reps | Percentage |
|-----|-------------|------------|
| 1 | 2x5 | Empty bar (45 lbs) |
| 2 | 1x5 | 40% of working weight |
| 3 | 1x3 | 60% of working weight |
| 4 | 1x2 | 80% of working weight |
| **Working** | **3x5** | **100%** |

#### Bench Press Warmup Protocol

| Set | Sets x Reps | Percentage |
|-----|-------------|------------|
| 1 | 2x5 | Empty bar (45 lbs) |
| 2 | 1x5 | 50% of working weight |
| 3 | 1x3 | 70% of working weight |
| 4 | 1x2 | 90% of working weight |
| **Working** | **3x5** | **100%** |

#### Deadlift Warmup Protocol

| Set | Sets x Reps | Percentage |
|-----|-------------|------------|
| 1 | 2x5 | 40% of working weight |
| 2 | 1x3 | 60% of working weight |
| 3 | 1x2 | 85% of working weight |
| **Working** | **1x5** | **100%** |

Note: Deadlift starts with weight on the bar (no empty bar warmup) due to proper pulling position requirements.

#### Overhead Press Warmup Protocol

| Set | Sets x Reps | Percentage |
|-----|-------------|------------|
| 1 | 2x5 | Empty bar (45 lbs) |
| 2 | 1x5 | 55% of working weight |
| 3 | 1x3 | 70% of working weight |
| 4 | 1x2 | 85% of working weight |
| **Working** | **3x5** (or **5x3**) | **100%** |

#### Power Clean Warmup Protocol

| Set | Sets x Reps | Percentage |
|-----|-------------|------------|
| 1 | 2x5 | Empty bar (45 lbs) |
| 2 | 1x5 | 55% of working weight |
| 3 | 1x3 | 70% of working weight |
| 4 | 1x2 | 85% of working weight |
| **Working** | **5x3** | **100%** |

---

## Linear Progression Model

### Session-to-Session Progression

The core principle of Starting Strength is linear progression: adding weight every workout.

```
Next_Session_Weight = Current_Session_Weight + lb_Increase
```

### Workout A/B Interleaving

For programs using A/B rotation, the squat progresses continuously across both workouts:

**Example with 5 lb increments:**
- Session 1 (Workout A): Squat 100 lbs
- Session 2 (Workout B): Squat 105 lbs
- Session 3 (Workout A): Squat 110 lbs
- Session 4 (Workout B): Squat 115 lbs

The working weight for Workout A, Session 3 equals:
```
Workout_B_Session_2_Weight + Squat_Increment
```

This creates a dependency chain where each workout references the previous workout's weight (regardless of A or B) plus the increment.

### Alternating Exercise Progression

For exercises that appear every other session (like Bench Press and Overhead Press alternating):

```
Next_Appearance_Weight = Previous_Appearance_Weight + lb_Increase
```

The increment is applied per appearance, not per total session. So if Bench Press appears on sessions 1, 3, 5, etc., each appearance adds the increment from the previous Bench Press session.

### Deadlift Double Progression

In some variants, deadlift uses a larger increment (10-15 lbs) but appears less frequently:

- **Original Novice**: Deadlift every Workout A (every other session), +10 lbs per appearance
- **Later variants**: Deadlift once per week, +15 lbs per week

---

## Weight Rounding

All calculated weights are rounded to loadable values using the smallest weight increment parameter.

### Rounding Rules

1. **Working Sets**: Round to nearest increment
   ```
   ROUND(Weight / SmallestWeight) * SmallestWeight
   ```

2. **Warmup Sets**: Round DOWN to nearest increment
   ```
   FLOOR(Weight / SmallestWeight) * SmallestWeight
   ```

### Empty Bar Minimum

The first warmup set is always the empty barbell (45 lbs for standard Olympic bar). If calculated warmup weights would be below 45 lbs, they remain at 45 lbs.

---

## Plate Loading Reference

The spreadsheet includes a plate weight calculator showing how to load the bar for each working weight.

### Standard Plate Set (per side)

| Plate Weight | Quantity Available |
|--------------|-------------------|
| 45 lbs | Multiple |
| 35 lbs | Multiple |
| 25 lbs | Multiple |
| 10 lbs | Multiple |
| 5 lbs | Multiple |
| 2.5 lbs | Multiple |

### Loading Calculation

For a target weight:
```
Weight_Per_Side = (Target_Weight - Bar_Weight) / 2
```

Where `Bar_Weight` = 45 lbs for standard Olympic bar.

### Example Loading Table

| Target Weight | 45 | 35 | 25 | 10 | 5 | 2.5 |
|---------------|----|----|----|----|---|-----|
| 215 | 1 | 1 | 0 | 0 | 1 | 0 |
| 220 | 1 | 1 | 0 | 0 | 1 | 1 |
| 225 | 2 | 0 | 0 | 0 | 0 | 0 |
| 230 | 2 | 0 | 0 | 0 | 0 | 1 |
| 235 | 2 | 0 | 0 | 0 | 1 | 0 |

---

## Accessory Exercises

### Chin-Ups / Pull-Ups

- **Sets x Reps**: 3 sets to failure
- **Progression**: No weight calculation; record reps achieved
- **Variation**: Advanced program alternates weighted and unweighted sessions

### Back Extensions / Glute Ham Raises

- **Sets x Reps**: 3-5 sets of 10 reps
- **Progression**: Bodyweight; add weight when 5x10 becomes easy

---

## Reset Protocol

When a lifter fails to complete all prescribed reps for multiple sessions:

1. Reduce working weight by 10-15% (configurable via "% to Reset" parameter)
2. Resume linear progression from the reduced weight
3. The reset percentage applies to the 5RM calculation:
   ```
   Reset_Weight = 5RM * (1 - Reset_Percentage)
   ```

---

## Implementation Notes

### Session Numbering

Sessions are numbered sequentially (1, 2, 3, ...) regardless of workout type (A or B) or day of week. This allows the spreadsheet to calculate 12+ weeks of training in advance.

### Formula Dependencies

Working set weights form a chain of dependencies:
- Session 1: Calculated from test weight and reset percentage
- Session N+1: Session N weight + increment

Warmup weights depend only on the current session's working weight.

### Minimum 45 lb Rule

Many warmup calculations use a minimum of 45 lbs (empty bar). When percentage calculations yield values below 45, the displayed warmup weight should be 45 lbs, but formula calculations may show lower values (40, etc.) in some cells when the working weight is very light.

---

## Summary

This Starting Strength calculator implements:

1. **Brzycki formula** for 1RM/5RM estimation
2. **Linear progression** with configurable increments per exercise
3. **Percentage-based warmups** with exercise-specific protocols
4. **Weight rounding** using FLOOR for warmups and ROUND for working sets
5. **Multiple program variants** accommodating different training schedules
6. **Reset functionality** for managing failed progressions

The algorithm enables a lifter to input their starting weights once and have all training weights calculated automatically for multiple weeks of progressive training.

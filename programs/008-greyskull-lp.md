# GreySkull LP Program Specification

## Overview

GreySkull LP (Linear Progression) is a beginner-to-intermediate strength training program that uses a 3-day-per-week rotating A/B schedule. The program is designed around the AMRAP (As Many Reps As Possible) final set protocol, which allows for autoregulated progression and built-in reset mechanisms.

## Program Structure

### Weekly Schedule

The program uses a 3-day rotating structure where workouts alternate between two templates (A and B):

- **Week 1**: A, B, A
- **Week 2**: B, A, B
- **Week 3**: A, B, A
- (Pattern repeats)

Each workout consists of:
1. Upper body pressing movement (with antagonist)
2. Lower body compound movement
3. Upper body accessory work
4. Lower body accessory work

### Day Templates

**Day A Structure:**
- Bench Press (with Row as antagonist)
- Squat
- Tricep Extension
- Ab Rollout

**Day B Structure:**
- Overhead Press (with Chin-up as antagonist)
- Deadlift OR Squat (Deadlift appears once per week on Day 2)
- Bicep Curl
- Shrug

### Deadlift Frequency

Deadlifts are programmed only once per week (on Day 2 of each week). The remaining lower body sessions are squats. This results in a weekly pattern of:
- Day 1: Squat
- Day 2: Deadlift
- Day 3: Squat

## Exercise Categories and Configuration

### Main Lifts

| Exercise | Default Increment | Default Rep Target | Notes |
|----------|-------------------|-------------------|-------|
| Bench Press | 2.5 units | 5 reps | Alternates with OHP |
| Overhead Press (OHP) | 2.5 units | 5 reps | Alternates with Bench |
| Squat | 2.5 units | 5 reps | 2x per week |
| Deadlift | 5 units | 5 reps | 1x per week, higher increment |

### Antagonist Movements (Bench/OHP Days)

| Exercise | Default Increment | Default Rep Target |
|----------|-------------------|-------------------|
| Row | 2.5 units | 5 reps |
| Chin-up | 2.5 units | 5 reps |

### Upper Body Accessories

| Exercise | Default Increment | Default Rep Target |
|----------|-------------------|-------------------|
| Tricep Extension | 2.5 units | 12 reps |
| Bicep Curl | 2.5 units | 12 reps |

### Lower Body Accessories

| Exercise | Default Increment | Default Rep Target |
|----------|-------------------|-------------------|
| Ab Rollout | N/A (bodyweight) | 12 reps |
| Shrug | 10 units | 12 reps |

## Set and Rep Scheme

### Main Lifts (Bench, OHP, Squat, Deadlift)

**Working Sets:** 3 sets total
- Set 1: Target reps (e.g., 5)
- Set 2: Target reps (e.g., 5)
- Set 3: AMRAP (target+) - perform as many reps as possible

**Note:** Deadlifts use only 1 working set (the AMRAP set), not 3 sets.

### Antagonist and Accessory Movements

**Working Sets:** 3 sets total
- Set 1: Target reps
- Set 2: Target reps
- Set 3: AMRAP (target+)

## Warm-up Protocol

### Upper Body Lifts (Bench, OHP)

Warm-up sets are calculated as percentages of the working weight, rounded down to the configured rounding increment (default 2.5). The minimum warm-up weight is the empty bar (45 lbs / 20 kg):

| Set | Reps | Percentage of Working Weight |
|-----|------|------------------------------|
| 1 | 5 | Empty bar (45 lbs or 20 kg) |
| 2 | 4 | 55% |
| 3 | 3 | 70% |
| 4 | 2 | 85% |

**Calculation:** `MAX(FLOOR(working_weight * percentage, rounding), empty_bar_weight)`

### Lower Body Lifts (Squat, Deadlift)

Warm-up uses 4 sets with different percentages:

| Set | Reps | Percentage of Working Weight |
|-----|------|------------------------------|
| 1 | 5 | 40% |
| 2 | 4 | 55% |
| 3 | 3 | 70% |
| 4 | 2 | 85% |

**Empty Bar:** 45 lbs for pounds, 20 kg for kilograms

## Linear Progression Rules

### Main Lifts Progression Algorithm

After each session, the working weight for the next session is calculated based on AMRAP set performance:

```
IF reps_achieved = 0 (no data entered):
    next_weight = current_weight + increment
ELSE IF reps_achieved < rep_target:
    next_weight = current_weight * 0.9  (10% deload)
ELSE IF reps_achieved < (rep_target * 2):
    next_weight = current_weight + increment  (standard progression)
ELSE IF reps_achieved >= (rep_target * 2):
    next_weight = current_weight + (increment * 2)  (double progression)
```

**Example with 5 rep target and 2.5 lb increment:**
- 0-4 reps: Reset to 90% of current weight
- 5-9 reps: Add 2.5 lbs
- 10+ reps: Add 5 lbs (double increment)

The result is rounded to the nearest configured rounding increment (default 2.5).

### Accessory Movements Progression Algorithm

Accessories use a modified progression with a tolerance band around the rep target:

```
IF reps_achieved = 0 (no data entered):
    next_weight = current_weight  (maintain)
ELSE IF reps_achieved < (rep_target - 2):
    next_weight = current_weight * 0.9  (10% deload)
ELSE IF reps_achieved > (rep_target + 2):
    next_weight = current_weight + increment  (progress)
ELSE:
    next_weight = current_weight  (maintain)
```

**Example with 12 rep target:**
- 0-9 reps: Reset to 90% of current weight
- 10-14 reps: Maintain current weight
- 15+ reps: Add increment

## Reset/Deload Protocol

### The 10% Reset Rule

When the lifter fails to achieve the minimum rep target on the AMRAP set:

1. **Trigger:** AMRAP reps < rep_target (e.g., less than 5 reps for main lifts)
2. **Action:** Reduce working weight by 10%
3. **Rounding:** Result is rounded to the nearest rounding increment

### Purpose of Resets

- Provides recovery opportunity without complete deload weeks
- Allows lifter to build momentum through previously stalled weights
- Creates natural periodization through the rep ranges

### Double Progression Bonus

When performance significantly exceeds expectations (hitting double the rep target or more), the program rewards this with double the normal increment. This encourages effort on AMRAP sets and accelerates progression when appropriate.

## Estimated 1RM Tracking

The Progress sheet calculates estimated one-rep max (e1RM) using the Brzycki formula:

```
e1RM = weight / (1.0278 - (0.0278 * reps))
```

This is calculated for each AMRAP set where data is entered, allowing tracking of strength progress over time independent of the rep/weight combinations used.

## Configuration Options

### User-Configurable Settings

| Setting | Purpose | Default (lbs) | Default (kg) |
|---------|---------|---------------|--------------|
| Units | Weight system | lb | kg |
| Warm-up Rounding | Precision for warm-up weights | 2.5 | 2.5 |
| Starting Weight | Initial working weight per lift | User-defined | User-defined |
| Increment | Weight added per successful session | Varies by lift | Varies by lift |
| Rep Target | Minimum reps for progression | Varies by lift | Varies by lift |
| Rounding | Precision for working weights | 2.5 | 2.5 |

### Exercise Customization

Exercise names can be customized while maintaining the program structure. The spreadsheet uses dynamic labels that pull from the configuration.

## Usage Workflow

1. **Initial Setup:**
   - Enter starting weights for all lifts (typically based on a percentage of 1RM or conservative estimates)
   - Configure unit system (lb or kg)
   - Adjust increments and rep targets if desired

2. **Training Session:**
   - Perform warm-up sets as prescribed
   - Complete working sets
   - On the final set, perform AMRAP
   - Record the reps achieved in the AMRAP set

3. **Automatic Calculations:**
   - The program calculates next session's working weight
   - Warm-up weights are recalculated automatically
   - Progress tracking updates estimated 1RM

## Implementation Notes

### Weight Rounding

All weights are rounded using `MROUND` (round to nearest multiple) to ensure practical plate loading:
- Default rounding: 2.5 units (allows for fractional plates)
- Can be set to 5 for gyms without small plates

### Empty Bar Handling

The minimum warm-up weight is always the empty bar:
- 45 lbs for pound-based programs
- 20 kg for kilogram-based programs

### Bodyweight Exercises

Exercises like Ab Rollout that don't use external weight:
- Weight field left blank
- Progression tracked by reps only
- Reset logic still applies based on rep performance

### Optional Exercises

Some exercises (like Row/Chin-up) can be disabled by setting their starting weight to 0, which hides them from the tracking.

## Summary

GreySkull LP is characterized by:

1. **AMRAP Final Sets:** Every exercise ends with a maximum effort set
2. **Autoregulated Progression:** Weight increases based on performance, not fixed schedule
3. **Built-in Deloads:** 10% resets occur automatically on failure
4. **Double Progression:** Exceptional performance earns bonus weight
5. **Flexible Structure:** Customizable exercises while maintaining proven template
6. **Deadlift Economy:** Once-weekly deadlifts allow recovery while maintaining skill
7. **Balanced Development:** Antagonist pairings and accessory work prevent imbalances

# RTS (Reactive Training Systems) Generalized Intermediate Program

## Overview

This spreadsheet implements the Reactive Training Systems (RTS) Generalized Intermediate Program, an autoregulated strength training system developed by Mike Tuchscherer. The program uses RPE (Rate of Perceived Exertion) as the primary autoregulation mechanism, combined with fatigue percentage calculations to manage training stress and volume.

**Original Source:** [RTS General Intermediate Program](http://forum.reactivetrainingsystems.com/content.php?126-Generalized-Intermediate-Program)

**Fatigue Methods Reference:** [Fatigue Percents Revisited](http://forum.reactivetrainingsystems.com/content.php?62-Fatigue-Percents-Revisited)

---

## Core Concepts

### RPE (Rate of Perceived Exertion) Scale

The RPE scale quantifies proximity to failure:
- **RPE 10:** Maximum effort, no additional reps possible
- **RPE 9:** Could do 1 more rep
- **RPE 8:** Could do 2 more reps
- **RPE 7:** Could do 3 more reps

### RPE-to-Percentage Chart

The system uses a lookup table that converts RPE and rep count to a percentage of 1RM:

| RPE / Reps | 1     | 2     | 3     | 4     | 5     | 6     | 7     | 8     | 9     | 10    | 11    | 12    |
|------------|-------|-------|-------|-------|-------|-------|-------|-------|-------|-------|-------|-------|
| **7**      | 0.88  | 0.82  | 0.80  | 0.74  | 0.74  | 0.68  | 0.66  | 0.64  | 0.62  | 0.60  | 0.58  | 0.56  |
| **8**      | 0.91  | 0.88  | 0.82  | 0.80  | 0.77  | 0.71  | 0.68  | 0.66  | 0.64  | 0.62  | 0.60  | 0.58  |
| **9**      | 0.95  | 0.91  | 0.89  | 0.82  | 0.80  | 0.74  | 0.71  | 0.68  | 0.66  | 0.64  | 0.62  | 0.60  |
| **10**     | 1.00  | 0.95  | 0.92  | 0.88  | 0.82  | 0.80  | 0.74  | 0.71  | 0.68  | 0.66  | 0.64  | 0.62  |

**Interpretation:** A set of 5 reps at RPE 9 should be performed at approximately 80% of 1RM.

---

## Input Parameters

### User-Configurable Inputs

1. **Lift Reference Maxes (1RM values)**
   - Each exercise variant has an associated 1RM value
   - Exercise names include variations (e.g., "squat", "squat 2ct pause", "squat tempo 303")
   - Each lift is classified into a movement category: Squat, Press, or Pull

2. **Goal Weekly Volume**
   - Bench: 11,000 (volume-load units)
   - Squat: 9,000 (volume-load units)
   - Deadlift: 10,000 (volume-load units)

3. **Round Increment**
   - Default: 5 (pounds/kg)
   - All calculated weights are rounded to this increment

### Lift Database Structure

Each lift entry contains:
- **Exercise Name:** Full name with modifiers (lowercase, concatenated)
- **1RM:** Current estimated one-rep max
- **Movement Class:** Category for volume tracking (Squat, Press, Pull)

Example entries:
```
squat           | 365 | Squat
squat 2ct pause | 315 | Squat
bench tng       | 300 | Press
bench competition | 285 | Press
deadlift        | 425 | Pull
```

---

## Workout Prescription Algorithm

### Step 1: Estimated Top Weight Calculation

For each prescribed set, the system calculates the estimated top weight using:

```
Estimated Top Weight = CEILING(
    RPE_Percentage(Target_RPE, Target_Reps) × Exercise_1RM / Round_Increment
) × Round_Increment
```

Where:
- `RPE_Percentage` is looked up from the RPE chart using target RPE and rep count
- `Exercise_1RM` is retrieved from the lift database by matching the exercise name (combining base exercise and modifiers)
- Result is rounded up to the nearest increment

**Example:**
- Exercise: Squat at 4 reps, RPE 9
- 1RM: 365 lbs
- RPE percentage for 4@9 = 0.82
- Calculated: 365 × 0.82 = 299.3 → rounds to 300 lbs

### Step 2: Warm-up/Working Weight Progression

The system generates three working weights as a comma-separated progression:

```
Working Weights = [
    Top_Weight - 2 × MROUND(Top_Weight × 0.05, Round_Increment),
    Top_Weight - MROUND(Top_Weight × 0.05, Round_Increment),
    Top_Weight
]
```

This creates a 3-set warm-up/build-up sequence with approximately 5% jumps:
- First weight: ~90% of top weight
- Second weight: ~95% of top weight
- Third weight: 100% (the top weight)

**Example:** Top weight of 300 lbs generates: "270, 285, 300"

---

## Fatigue Management System

### Fatigue Percentage

The fatigue percentage (ranging from 0 to 0.10+) determines how much weight to drop for subsequent work sets:

```
Fatigue Weight = Top_Weight - CEILING(Top_Weight × Fatigue_Percent / Round_Increment) × Round_Increment
```

**Example:** Top weight 340 with 5% fatigue = 340 - (340 × 0.05 → 17 → 20) = 320 lbs

### Fatigue Methods

The program uses three distinct fatigue methods:

#### 1. Load Drop Method
**Usage:** Primary method for main working sets

**Protocol:**
1. Work up to the initial set at prescribed reps and RPE
2. Reduce weight by the programmed fatigue percentage
3. Continue performing sets at the same rep count until RPE reaches the target

**When to use:** Listed as "Load Drop" in the Fatigue Method column

#### 2. Repeat Sets Method
**Usage:** For accumulating volume at consistent intensity

**Protocol:**
1. Perform initial set at prescribed reps and RPE
2. Repeat sets at the same weight until RPE increases (based on fatigue percentage and RPE chart relationship)

**When to use:** Listed as "Repeat" in the Fatigue Method column

#### 3. Rep Drop Method
**Usage:** For maintaining intensity while reducing per-set fatigue

**Protocol:**
1. Perform initial set at prescribed reps and RPE
2. Repeat sets with fewer reps (calculated from fatigue percentage and RPE chart) until reaching target RPE

**When to use:** Less commonly prescribed in this template

---

## Volume Calculation

### Expected Volume

Volume is calculated as volume-load (weight × reps × sets):

```
Expected Volume = SUM(Working_Weights) × Reps + (Fatigue_Reps × Fatigue_Sets × Fatigue_Weight)
```

The system tracks:
- Volume from the three warm-up/working sets
- Additional volume from fatigue sets (recorded as "RxS" format, e.g., "4x3" = 4 reps for 3 sets)

### Done Volume (Actual)

Actual volume is calculated from user-entered data:
- **Done Working Weights:** Comma-separated list of actual weights used
- **Done Working Reps:** Comma-separated list of actual reps performed
- **Fatigue (RxS):** Additional sets performed in "RepsxSets" format

The calculation includes a 2x multiplier for dumbbell exercises (detected by "DB" in the exercise name).

---

## Movement Type Classification

Exercises are classified into three categories for volume tracking:

| Category | Examples |
|----------|----------|
| **Squat** | Squat, Front Squat, Squat Tempo, Leg Press |
| **Press** | Bench Competition, Bench TnG, OHP, Push Press, JM Press |
| **Pull** | Deadlift, Deadlift Conventional, Pendlay Row, Lat Pulldown |

Classification is determined by looking up the exercise in the lift database.

---

## Program Structure

### Weekly Organization

The program spans 9 weeks organized into mesocycles:

**Week 1:** Introduction/Baseline
- No fatigue method specified
- Establishes working weights
- RPE targets typically 9-10

**Weeks 2-4:** Development Phase
- Load Drop method introduced
- Fatigue percentages: 0-7.5%
- Progressive intensity increases

**Weeks 5-8:** Intensification Phase
- Mix of Load Drop and Repeat methods
- Fatigue percentages: 5-7%
- Higher intensities, lower volumes

**Week 9:** Peaking/Testing
- Singles at RPE 7-10
- No fatigue percentage (0%)
- Test new maxes

### Training Day Structure

Each training day includes 3-4 exercises:
- Primary competition lift or variant
- Secondary pressing movement
- Accessory work

Days are separated by blank rows in the data structure.

---

## Volume Tracking System

### Daily Volume Aggregation

The Volume sheet tracks cumulative volume by movement category:

1. **Date Column:** Extracts unique dates from the Main sheet
2. **Expected Volume:** Sums expected volume for each movement type on each date
3. **Actual Volume:** Sums completed volume for each movement type on each date

### Weekly Volume Summation

Weekly totals are calculated at 4-day intervals (every 4th row):
- Rolling sum of daily volumes
- Compared against goal weekly volumes
- Allows monitoring of training load progression

---

## User Workflow

### Pre-Workout
1. Enter the date in the Date column
2. Review prescribed exercises, reps, RPE, and fatigue method
3. Note the estimated top weight and working weight progression

### During Workout
1. Perform warm-up sets building to the top weight
2. Execute the first working set at target reps
3. Apply the appropriate fatigue method for additional sets
4. Stop when RPE reaches the target

### Post-Workout
1. Enter actual weights used in "Done Working Weights" (comma-separated)
2. Enter actual reps performed in "Done Working Reps" (comma-separated)
3. Enter additional fatigue work in "Fatigue (RxS)" format (e.g., "4x3")

---

## Key Formulas Summary

### Estimated Top Weight
```
RPE_Percentage × Exercise_1RM, rounded up to increment
```

### Working Weight Progression
```
[90%, 95%, 100%] of top weight (approximately)
```

### Fatigue Weight
```
Top_Weight × (1 - Fatigue_Percent), rounded to increment
```

### Volume Load
```
Weight × Reps × Sets (summed across all work)
```

---

## Technical Notes

### Exercise Name Matching
- Exercise names are matched case-insensitively
- Base exercise and modifier are concatenated with a space
- Trailing/leading spaces are trimmed

### Rounding Behavior
- All weights are rounded to the configured increment (default: 5)
- RPE chart lookup uses exact matching for both RPE and rep values

### Error Handling
- Missing lift data returns blank estimated weights
- Division by zero is protected in fatigue calculations
- Invalid exercise lookups return empty strings

---

## Customization Guidelines

### Modifying the Program

1. **Changing Exercises:** Update the lift database with new exercises and their 1RM values
2. **Adjusting Volume Goals:** Modify the Goal Weekly Volume values
3. **Changing Increment:** Update the Round Increment for finer or coarser weight jumps
4. **Adding Weeks:** Copy existing week structure and adjust prescriptions

### Adding New Exercises

1. Add exercise name (lowercase) to the lift database
2. Enter the estimated 1RM
3. Assign appropriate movement class (Squat, Press, or Pull)
4. Use the exact exercise name (with modifiers) in the Main sheet

---

## References

- Tuchscherer, M. "Reactive Training Systems Manual"
- RTS Forum: Generalized Intermediate Program documentation
- RTS Forum: Fatigue Percents Revisited

---

*Last Updated: 2017-09-28*
*Spreadsheet by: [MagnesiumCarbonate](https://www.reddit.com/user/MagnesiumCarbonate)*

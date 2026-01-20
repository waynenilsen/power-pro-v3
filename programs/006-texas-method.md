**Texas Method Training Algorithm**

*Complete Technical Specification*

1. Program Overview
===================

The Texas Method is an intermediate strength training program designed
for lifters who have exhausted linear progression. It operates on a
weekly periodization cycle with three distinct training days that
manipulate volume, intensity, and recovery to drive continued strength
adaptation.

1.1 Core Philosophy
-------------------

The program alternates stress and recovery within each week. Monday
provides a high-volume stimulus, Wednesday allows active recovery, and
Friday tests intensity with new personal records. This weekly undulation
creates the stress-recovery-adaptation cycle necessary for intermediate
trainees.

1.2 Weekly Structure
--------------------

  **Day**     **Training Focus**               **Primary Purpose**
  ----------- -------------------------------- ------------------------------------------------
  Monday      Volume Day (Medium Intensity)    Accumulate training stress through high volume
  Wednesday   Recovery Day (Light Intensity)   Facilitate recovery with reduced load
  Friday      Intensity Day (Heavy)            Test strength with new personal records

2. Mathematical Foundations
===========================

2.1 Rep Max Estimation Formula
------------------------------

The program uses the Epley formula to convert any rep max to an
estimated one-rep maximum (1RM). This formula provides the mathematical
basis for all weight calculations throughout the program.

### Epley Formula

> **Estimated 1RM = Weight × (1 + Reps ÷ 30)**

Algebraically rearranged for spreadsheet implementation:

> **Estimated 1RM = Weight ÷ (1.0278 - (0.0278 × Reps))**

### Reverse Calculation (1RM to Rep Max)

To calculate the weight for a target number of reps from a known 1RM:

> **Weight for N Reps = 1RM × (1.0278 - (0.0278 × N))**

2.2 Rounding Algorithm
----------------------

All calculated weights must be rounded to practical increments that can
be loaded on a barbell. The rounding behavior differs by calculation
type:

-   Intensity Day weights: Round DOWN to nearest increment (conservative
    approach for max attempts)

-   Volume/Recovery Day weights: Round to NEAREST increment (standard
    rounding)

-   Warm-up weights: Round DOWN to nearest 5 units (simplifies plate
    loading)

2.3 Standard Increments by Lift
-------------------------------

  **Lift**         **Weight Increment (lbs)**   **Weight Increment (kg)**
  ---------------- ---------------------------- ---------------------------
  Squat            5                            2.5
  Deadlift         5                            2.5
  Bench Press      2.5                          1.25
  Overhead Press   2.5                          1.25
  Power Clean      2.5                          1.25

3. Weight Calculation Algorithm
===============================

3.1 Initial Setup
-----------------

The user provides their current working weights and the number of reps
achieved at that weight for each lift. From this input, the system
calculates the estimated 1RM which becomes the baseline for all
subsequent calculations.

### Step-by-Step Initialization

1.  User enters current weight and reps for each lift (typically a 5-rep
    max)

2.  System calculates estimated 1RM using the Epley formula

3.  Friday (Intensity Day) weight is set by rounding the 5RM equivalent
    DOWN to the nearest increment

4.  Monday and Wednesday weights are derived from Friday weight using
    percentage calculations

3.2 Day-Specific Weight Formulas
--------------------------------

### Friday (Intensity Day) - The Anchor

Friday weight serves as the reference point. It represents the target
for a new 5-rep (or lift-specific rep scheme) personal record.

> **Friday Weight = ROUND\_DOWN(Estimated 1RM ÷ Increment) × Increment**

### Monday (Volume Day)

Monday uses 90% of Friday's weight, providing sufficient stimulus
without pre-fatiguing the lifter before Friday's intensity work.

> **Monday Weight = ROUND((Friday Weight × 0.90) ÷ Increment) ×
> Increment**

### Wednesday (Recovery Day)

Wednesday uses 80% of Monday's weight (effectively \~72% of Friday),
allowing active recovery while maintaining movement patterns.

> **Wednesday Weight = ROUND((Monday Weight × 0.80) ÷ Increment) ×
> Increment**

4. Progression Algorithm
========================

4.1 Core Progression Logic
--------------------------

Progression occurs weekly based on performance on the intensity day. The
algorithm uses a simple conditional check:

> **IF Reps Achieved \> Target Reps THEN Add Increment ELSE Maintain
> Weight**

### Detailed Progression Rules

  **Lift**         **Rep Target**     **Progression Trigger**      **Increment Added**
  ---------------- ------------------ ---------------------------- ---------------------
  Squat            5 reps on Friday   More than 4 reps hit         5 lbs / 2.5 kg
  Bench Press      5 reps on Friday   More than 4 reps hit         2.5 lbs / 1.25 kg
  Overhead Press   5 reps on Friday   More than 4 reps hit         2.5 lbs / 1.25 kg
  Deadlift         5 reps on Monday   More than 4 reps hit         5 lbs / 2.5 kg
  Power Clean      3 sets on Friday   More than 3 sets completed   2.5 lbs / 1.25 kg

4.2 Progression Cascade
-----------------------

When Friday weight increases, Monday and Wednesday weights automatically
recalculate based on their percentage relationships. This creates a
cascading effect:

1.  User hits target reps on Friday

2.  Next week's Friday weight increases by the lift-specific increment

3.  Next week's Monday weight recalculates to 90% of new Friday weight

4.  Next week's Wednesday weight recalculates to 80% of new Monday
    weight

4.3 Stall Handling
------------------

If the target reps are not achieved, the weight remains constant for the
following week. The algorithm does not prescribe automatic deloads;
repeated stalls require manual intervention by the user to adjust
training variables.

5. Warm-up Protocol Algorithm
=============================

5.1 Standard Warm-up Structure
------------------------------

Each lift begins with a standardized warm-up progression that prepares
the lifter for the working weight while minimizing fatigue accumulation.

### Warm-up Set Percentages

  **Set**   **Percentage of Work Weight**   **Reps**   **Rounding**
  --------- ------------------------------- ---------- -------------------
  Set 1     Empty Bar (45 lbs / 20 kg)      2 × 5      Fixed
  Set 2     40%                             1 × 5      Down to nearest 5
  Set 3     60%                             1 × 3      Down to nearest 5
  Set 4     80%                             1 × 2      Down to nearest 5

5.2 Warm-up Calculation Formula
-------------------------------

> **Warm-up Weight = FLOOR(Work Weight × Percentage ÷ 5) × 5**

The FLOOR function ensures weights always round down to the nearest
5-unit increment, making plate loading practical and consistent.

5.3 Lift-Specific Variations
----------------------------

Some lifts use modified warm-up percentages based on their
characteristics:

-   Bench Press: Uses 50%, 70%, 90% progression (higher percentages due
    to shorter range of motion)

-   Overhead Press: Uses 55%, 70%, 85% progression

-   Power Clean: Uses 55%, 70%, 85% progression

6. Lift-Specific Programming
============================

6.1 Squat Programming
---------------------

  **Day**     **Sets × Reps**   **Intensity**       **Calculation**
  ----------- ----------------- ------------------- ---------------------------
  Monday      5 × 5             90% of Friday       Round(Friday × 0.9)
  Wednesday   2 × 5             80% of Monday       Round(Monday × 0.8)
  Friday      1 × 5             100% (PR attempt)   Base weight + progression

6.2 Bench Press Programming
---------------------------

Bench press alternates with overhead press, appearing twice per week on
a rotating schedule.

  **Day**              **Sets × Reps**   **Intensity**   **Notes**
  -------------------- ----------------- --------------- -----------------
  Monday (Week A)      5 × 5             90% of Friday   Volume emphasis
  Friday (Week A)      1 × 5             100%            PR attempt
  Wednesday (Week B)   3 × 5             90% of 5RM      Moderate volume

6.3 Overhead Press Programming
------------------------------

Press follows a similar alternating pattern with bench, with an optional
intensity reduction on the lighter day.

### Wednesday Decrement Factor

An optional multiplier (default: 0.95) can reduce Wednesday press volume
to prevent excessive shoulder fatigue:

> **Wednesday Weight = Round((Friday 5RM × 0.9 × Decrement Factor))**

6.4 Deadlift Programming
------------------------

Deadlift receives less volume due to its systemic fatigue. It typically
appears once per week.

  **Day**   **Sets × Reps**   **Intensity**   **Calculation**
  --------- ----------------- --------------- ------------------------
  Monday    1 × 5             90% of 1RM      Round\_Down(1RM × 0.9)

Progression occurs when Monday's reps exceed the target (typically more
than 4 reps achieved).

6.5 Power Clean Programming
---------------------------

  **Day**   **Sets × Reps**   **Intensity**   **Notes**
  --------- ----------------- --------------- ------------------------------
  Friday    5 × 3             100%            Technical precision emphasis

Power clean progression triggers when the lifter completes more than 3
sets successfully at the prescribed weight.

7. Volume Tracking Algorithm
============================

7.1 Total Volume Calculation
----------------------------

The system tracks total training volume for each session using the
standard volume formula:

> **Volume = Weight × Sets × Reps**

7.2 Session Volume Summation
----------------------------

Total session volume aggregates across all exercises, including warm-up
sets:

> **Session Volume = Σ(Weight\_i × Sets\_i × Reps\_i) for all exercises
> i**

### Example Calculation

For a squat session with warm-ups and work sets:

-   Bar (45 lbs) × 2 sets × 5 reps = 450 lbs

-   40% (80 lbs) × 1 set × 5 reps = 400 lbs

-   60% (120 lbs) × 1 set × 3 reps = 360 lbs

-   80% (160 lbs) × 1 set × 2 reps = 320 lbs

-   100% (200 lbs) × 5 sets × 5 reps = 5,000 lbs

**Total Squat Volume:** 6,530 lbs

8. Statistics and Performance Metrics
=====================================

8.1 Estimated 1RM Tracking
--------------------------

The system continuously tracks estimated 1RM for each lift based on the
most recent intensity day performance:

> **Current 1RM Estimate = Weight Lifted ÷ (1.0278 - (0.0278 × Reps
> Completed))**

8.2 Powerlifting Total
----------------------

For competition tracking, the system calculates a powerlifting total
from the three competition lifts:

> **PL Total = Squat 1RM + Bench 1RM + Deadlift 1RM**

8.3 Wilks Score Calculation
---------------------------

The Wilks coefficient normalizes strength performance across different
body weights, enabling fair comparison between lifters.

### Wilks Formula

> **Wilks Score = Total × (500 ÷ Coefficient)**

Where the coefficient is calculated using a 5th-order polynomial based
on body weight:

> **Coefficient = a + b(x) + c(x²) + d(x³) + e(x⁴) + f(x⁵)**

Where x = body weight in kilograms, and coefficients for male lifters
are:

-   a = -216.0475144

-   b = 16.2606339

-   c = -0.002388645

-   d = -0.00113732

-   e = 7.01863 × 10⁻⁶

-   f = -1.291 × 10⁻⁸

8.4 Strength-to-Weight Ratios
-----------------------------

The system calculates relative strength by dividing each lift's 1RM by
the lifter's body weight:

> **Strength Ratio = Estimated 1RM ÷ Body Weight**

9. Unit Conversion
==================

9.1 Metric/Imperial Handling
----------------------------

The system supports both pounds (lbs) and kilograms (kg) with automatic
adjustments:

### Empty Bar Weight

  **Unit System**   **Standard Bar Weight**
  ----------------- -------------------------
  Pounds (lbs)      45 lbs
  Kilograms (kg)    20 kg

### Rounding Increments

All rounding operations use the appropriate increment for the selected
unit system. The increment values specified in the configuration
(typically 2.5 or 5) are interpreted in the active unit system.

10. Cycle Management
====================

10.1 Cycle Structure
--------------------

The program organizes training into cycles, with each cycle representing
one week of training (Monday, Wednesday, Friday). Multiple cycles are
tracked for long-term progression analysis.

10.2 Inter-Cycle Progression
----------------------------

At the end of each cycle:

1.  User records actual reps/sets achieved on intensity days

2.  Progression algorithm evaluates performance against targets

3.  Next cycle weights are calculated based on progression rules

4.  Statistics are updated with new 1RM estimates and totals

10.3 Data Flow Summary
----------------------

The complete algorithm flow for each training cycle:

-   INPUT: User's current lifts (weight × reps), body weight, unit
    preference, progression increments

-   PROCESS: Calculate 1RMs, derive day-specific weights, generate
    warm-up progressions

-   OUTPUT: Complete training prescription for Mon/Wed/Fri including all
    weights, sets, and reps

-   FEEDBACK: User records performance, triggering progression
    calculations for next cycle

# nSuns CAP3: Cyclical AMRAP Progression - 3 Week Variant

## Program Overview

CAP3 (Cyclical AMRAP Progression, 3-Week Variant) is a strength training program designed by nSuns focused on increasing the Squat, Bench Press, and Deadlift through a structured 3-week rotating cycle. The program uses estimated maxes rather than tested 1RMs and features a cyclical design where each primary lift receives a high-intensity AMRAP (As Many Reps As Possible) test once every three weeks.

### Key Principles

- **Cyclical Rotation**: Each of the three main lifts (Bench, Squat, Deadlift) gets one maximum effort AMRAP day per 3-week cycle
- **Dual Progression System**: Major progression through high-intensity AMRAP tests, secondary progression through regular AMRAP sets
- **Training Max Based**: All weights calculated from a Training Max (TM) set at 90% of estimated 1RM
- **Volume Accumulation**: Higher volume phases build into intensity peaks
- **EMOM Protocol**: Secondary lifts use Every Minute On the Minute sets for volume regulation

---

## Training Max Calculation

### Initial Setup

The Training Max (TM) serves as the foundation for all weight calculations:

```
Training Max = Estimated Max * 0.9
```

Where Estimated Max can be either:
- A true tested 1RM
- A calculated estimate from a 2-3 rep max using the Epley formula

### Epley Formula (Estimated 1RM Calculator)

```
Estimated 1RM = Weight * (1 + 0.0333 * Reps)
```

Or rearranged to find target weight for a rep goal:
```
Target Weight = TM / (0.995 + 0.0333 * Target_Reps)
```

Note: The formula uses 0.995 instead of 1.0 as the base constant in this implementation.

---

## The 3-Week Cycle Structure

### Cyclical AMRAP Rotation

The program rotates which lift receives maximum effort testing each week:

| Week | Deadlift | Squat | Bench |
|------|----------|-------|-------|
| Week 1 | HIGH INTENSITY AMRAP | Medium Volume | Volume |
| Week 2 | Medium Volume | HIGH INTENSITY AMRAP | Medium Volume |
| Week 3 | Volume | Volume | HIGH INTENSITY AMRAP |

This ensures each primary lift is tested at maximum effort exactly once per cycle, while accumulating volume in the other weeks.

### Day Numbering by Week

| Training Day | Week 1 | Week 2 | Week 3 |
|--------------|--------|--------|--------|
| Day 1 | 1 | 8 | 16 |
| Day 2 | 2 | 9 | 17 |
| REST | - | Day 10 | - |
| Day 3 | 3 | 11 | 18 |
| Day 4 | 4 | 12 | 19 |
| REST | Day 5 | - | - |
| Day 5/6 | 6 | 13 | 20 |
| Day 6/7 | 7 | 14 | 21 |
| REST | - | - | Day 15 |

Rest days are strategically placed before maximum effort AMRAP days.

---

## Weekly Training Schedule

### Day 1: Chest and Biceps
- **Primary Lift**: Bench Press (main variation)
- **Secondary Lift**: Bench Press (close variation - e.g., Close Grip Bench)
- **Accessory Focus**: Chest, Biceps

### Day 2: Back, Abs, and Triceps
- **Primary Lift**: Deadlift (secondary variation - e.g., Paused Deadlift)
- **Secondary Lift**: Row Variation (e.g., Cable Row)
- **Accessory Focus**: Back, Abs, Triceps

### Day 3: Shoulders and Legs
- **Primary Lift**: Squat (main variation - e.g., Low Bar Squat)
- **Secondary Lift**: Press Variation (e.g., Overhead Press)
- **Accessory Focus**: Shoulders, Legs

### Day 4: Chest and Biceps
- **Primary Lift**: Bench Press (main variation) - different rep scheme from Day 1
- **Secondary Lift**: Bench Press (close variation)
- **Accessory Focus**: Chest, Biceps

### Day 5/6: Back, Abs, and Triceps
- **Primary Lift**: Deadlift (main variation)
- **Secondary Lift**: Row Variation (main variation)
- **Accessory Focus**: Back, Abs, Triceps

### Day 6/7: Legs and Shoulders
- **Primary Lift**: Squat (secondary variation - e.g., Front Squat)
- **Secondary Lift**: Press Variation
- **Accessory Focus**: Legs, Shoulders

---

## Primary Lift Set/Rep Schemes

### High Intensity AMRAP Day (Major Progression)

When a lift is scheduled for its high-intensity day, the user chooses ONE of three options:

| Option | Percentage | Minimum Reps |
|--------|------------|--------------|
| Option 1 | 79.5% TM | 6+ |
| Option 2 | 83.5% TM | 4+ |
| Option 3 | 88.5% TM | 2+ |

**Structure**:
1. First set at 79.5% TM x 6 reps (warm-up/work-up)
2. Optional set marker
3. Second set at 83.5% TM x 4 reps
4. Optional set marker
5. Final AMRAP set at 88.5% TM x 2+ reps

The "- o r -" markers indicate optional sets - perform them only if feeling strong.

### Volume Day (Accumulation Phase)

Higher rep ranges at lower percentages with multiple sets:

**Bench Example (Week 1 - Volume Day)**:
- Sets 1-4: 73.5% TM x 4 reps
- Sets 5-6: 73.5% TM + 5lb x 4 reps
- Set 7: 73.5% TM + 10lb x 4+ reps (AMRAP)

**Total**: 7 sets, all at 4 reps with final set as AMRAP

### Medium Volume Day

Moderate intensity with moderate rep ranges:

**Bench Example (Week 2 - Medium Day)**:
- Sets 1-4: 77% TM x 3 reps
- Sets 5-6: 77% TM + 5lb x 3 reps
- Set 7: 77% TM + 10lb x 3 reps
- Set 8: 77% TM + 10lb x 3+ reps (AMRAP)

**Total**: 8 sets at 3 reps with final set as AMRAP

---

## Primary Lift Percentage Tables by Week

### Bench Press (Day 1)

| Week | Phase | Base % | Rep Scheme | Sets | Final Set |
|------|-------|--------|------------|------|-----------|
| 1 | Volume | 73.5% | 4 reps | 7 | 4+ AMRAP |
| 2 | Medium | 77% | 3 reps | 8 | 3+ AMRAP |
| 3 | High Intensity | 79.5/83.5/88.5% | 6/4/2 | 3-5 | 2+ AMRAP |

### Bench Press (Day 4) - Secondary Bench Day

| Week | Phase | Base % | Rep Scheme | Sets | Final Set |
|------|-------|--------|------------|------|-----------|
| 1 | Volume | 68% | 6 reps | 7 | 5+ AMRAP |
| 2 | Medium | 72.5-74% | 5-4 reps | 8 | 4+ AMRAP |
| 3 | High Volume | 62-65% | 8-6 reps | 6 | 6+ AMRAP |

### Squat (Day 3)

| Week | Phase | Base % | Rep Scheme | Sets | Final Set |
|------|-------|--------|------------|------|-----------|
| 1 | Medium | 77% | 3 reps | 7 | 3+ AMRAP |
| 2 | High Intensity | 79.5/83.5/88.5% | 6/4/2 | 3-5 | 2+ AMRAP |
| 3 | Volume | 73.5-75.5% | 4 reps | 6 | 4+ AMRAP |

### Squat Variant (Day 7) - e.g., Front Squat

| Week | Phase | Base % | Rep Scheme | Sets | Final Set |
|------|-------|--------|------------|------|-----------|
| 1 | Medium | 73-75.5% | 5-4 reps | 8 | 4+ AMRAP |
| 2 | High Volume | 62-65% | 8-6 reps | 6 | 6+ AMRAP |
| 3 | Volume | 68-70% | 6-5 reps | 7 | 5+ AMRAP |

### Deadlift (Day 6)

| Week | Phase | Base % | Rep Scheme | Sets | Final Set |
|------|-------|--------|------------|------|-----------|
| 1 | High Intensity | 79.5/83.5/88.5% | 6/4/2 | 3-5 | 2+ AMRAP |
| 2 | Medium | 74-75% | 4 reps | 6 | 4+ AMRAP |
| 3 | Volume | 77% | 3 reps | 7 | 3+ AMRAP |

### Deadlift Variant (Day 2) - e.g., Paused Deadlift

| Week | Phase | Base % | Rep Scheme | Sets | Final Set |
|------|-------|--------|------------|------|-----------|
| 1 | Medium | 73-75.5% | 5-4 reps | 8 | 4+ AMRAP |
| 2 | High Volume | 62-65% | 8-6 reps | 6 | 6+ AMRAP |
| 3 | Volume | 68-70% | 6-5 reps | 7 | 5+ AMRAP |

---

## Secondary Lift Protocol

Secondary lifts use a fundamentally different approach based on rep PR targets and EMOM volume accumulation.

### Weight Calculation for Secondary Lifts

Weight is calculated using the inverse Epley formula to target specific rep ranges:

```
Target Weight = TM / (0.995 + 0.0333 * Target_Reps)
```

### Secondary Lift Structure

Each secondary lift session consists of three phases:

**Phase 1: AMRAP Set**
Perform one AMRAP set at the calculated weight for the target rep range.

| Week | Target Reps (varies by lift) |
|------|------------------------------|
| 1 | 6-10+ reps |
| 2 | 8-12+ reps |
| 3 | 10-15+ reps |

**Phase 2: Max Rep Set (MRS)**
Drop weight to a lower percentage and perform a max rep set:

| Week | Target Reps for MRS |
|------|---------------------|
| 1 | 4-6+ reps |
| 2 | 6-8+ reps |
| 3 | 8-10+ reps |

**Phase 3: EMOM Repeats**
Continue performing sets at the MRS weight every minute on the minute until you fail to complete the minimum rep threshold:

| Week | EMOM Minimum Reps |
|------|-------------------|
| 1 | 3-4 reps |
| 2 | 4-5 reps |
| 3 | 5-6 reps |

### Secondary Lift Example (Close Grip Bench, Week 1)

1. **AMRAP Set**: TM / (0.995 + 0.0333 * 12) = ~71% TM for 10+ reps
2. **MRS Set**: TM / (0.995 + 0.0333 * 15) = ~67% TM for 8+ reps
3. **EMOM Repeats**: Continue at MRS weight until failing to complete 5 reps

---

## Progression System

### Major Progression (High Intensity AMRAP Days)

Evaluated once every 3 weeks when a lift has its high-intensity day:

**If you complete minimum reps but don't set a new estimated max:**
- Add 2.5 kg (5 lb) to Training Max
- Continue to next cycle

**If you complete enough reps to set a new estimated max:**
- Calculate new estimated max using Epley formula
- Set Training Max = New Estimated Max (rounded up)
- Continue to next cycle

**If you fail to complete minimum reps:**
- Reduce Training Max by 5-10%
- Focus on pushing AMRAP and accessory work harder
- Continue to next cycle

### Secondary Progression (Yellow AMRAP Sets)

Evaluated weekly on regular training days:

**If you set a rep PR on a yellow AMRAP set:**
- Calculate new estimated max
- Adjust Training Max to new estimated max

**If you complete all required reps but no rep PR by cycle end:**
- Add 2.5 kg (5 lb) to Training Max
- Continue to next cycle

**If you fail to complete required reps:**
- Reduce Training Max by 5-10%
- Push AMRAP sets harder next cycle

---

## Weight Rounding

All calculated weights are rounded to the nearest increment:

| Unit System | Rounding Increment |
|-------------|-------------------|
| Pounds (Lb) | 5 lb |
| Kilograms (Kg) | 2.5 kg |

Formulas use FLOOR() for conservative rounding down or CEILING() for rounding up, depending on the context:
- Warm-up/early sets: FLOOR (round down)
- Working/AMRAP sets: CEILING (round up)

---

## Volume and Intensity Tracking

### Volume Calculation

Daily volume is calculated as:
```
Volume = Sum of (Weight * Reps) for all sets
```

### INOL (Intensity x Number of Lifts)

INOL provides a fatigue metric:
```
INOL = Reps / (100 - (100 * (Weight / TM)))
```

Or simplified:
```
INOL = Reps / (100 - Percentage)
```

Where Percentage = (Weight / TM) * 100

Higher INOL indicates greater training stress.

---

## Accessory Work Guidelines

After completing the two programmed lifts, incorporate accessory work based on:

- Time available
- Energy levels
- Identified weakpoints
- Muscle group focus for the day

### Recommended Approach
- 1-4 exercises per muscle group
- 6-8+ reps per set
- As many sets as can be recovered from

### Daily Prehab Work
Each day includes recommended band and mobility work:
- 150-200 Band Pull-Aparts
- 75-100 Banded Shoulder Dislocations
- Chin-ups (bodyweight or weighted, varying by day)
- Mobility work on leg days

---

## Lift Selection Options

The program includes dropdown selections for customizing exercise variants:

### Squat Variations
Low Bar, High Bar, Front Squat, SSB Squat, Box Squats (Low/High), Paused Squat, Belt Squat, Split Squat, Stance Variations, Duffalo Bar, Lunges, Hatfield Squat

### Bench Variations
Paused, Flat, Close Grip, Incline, DB Bench, Slingshot, Spoto, Feet Up, JM Press, Duffalo Bar, Chains

### Deadlift Variations
Conventional, Sumo, Deficit (conventional/sumo), Rack Pull (conventional/sumo), Block Pull, Paused, Romanian, SLDL, Axle, Chains

### Press Variations
Strict OHP, Push Press, DB OHP, Seated Press, Z Press, One-Arm variations

### Row Variations
Strict Barbell, Cheat Row, Reverse Grip, DB Row, Chest Supported, T-Bar, Cable Row

---

## Summary Algorithm

### Weekly Flow

1. **Select Week Number** (1, 2, or 3)
2. **For each training day**:
   - Complete warm-up and prehab work
   - Perform Primary Lift 1 according to week's protocol
   - Perform Secondary Lift 2 with AMRAP + MRS + EMOM structure
   - Complete accessory work as desired
3. **Track AMRAP performance** for progression decisions
4. **At cycle end**:
   - Evaluate all AMRAP performances
   - Adjust Training Maxes according to progression rules
   - Change week selector back to Week 1
   - Begin new cycle

### Cycle Progression Logic

```
FOR each lift:
    IF high_intensity_week:
        IF completed_minimum_reps AND new_estimated_max:
            TM = new_estimated_max (rounded up)
        ELSE IF completed_minimum_reps:
            TM = TM + 5lb/2.5kg
        ELSE:
            TM = TM * 0.90 to 0.95
    ELSE:
        IF set_rep_PR on AMRAP:
            TM = new_estimated_max
        ELSE IF completed_all_reps at cycle_end:
            TM = TM + 5lb/2.5kg
        ELSE:
            TM = TM * 0.90 to 0.95
```

---

## Deload Protocol

- Deloads are taken as needed after completing a full cycle
- If deloads are needed more frequently than every 2-3 cycles, address recovery factors (sleep, nutrition, stress)
- Training frequency should not drop below 4 sessions per week even during recovery periods

---

## Program Attribution

CAP3 was created by nSuns as a personal training program combining elements from various programs run over years of training. The cyclical structure allows for consistent incremental progress while managing fatigue through strategic intensity rotation.

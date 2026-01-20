# Greg Nuckols Beginner Program Specification

## Overview

This is a general beginner strength training program designed by Greg Nuckols. It is a 3-day per week program that trains the squat, bench press, and deadlift with varying frequencies optimized for beginners.

## Program Philosophy

The program follows several key principles:

1. **Autoregulation through AMAP sets**: The bench press routine uses "As Many As Possible" (AMAP) sets to gauge progress and determine weight increases
2. **Progressive overload**: Systematic weekly increases in either weight or volume
3. **Periodization**: The squat and deadlift follow 4-week cycles with planned intensity progressions
4. **Frequency optimization**: Each lift is trained at a frequency suited to its recovery demands:
   - Bench Press: 3x per week (highest frequency for upper body)
   - Squat: 2x per week
   - Deadlift: 2x per week

## Weekly Training Schedule

The program uses a 3-day training week with the following structure:

| Day | Session 1 | Session 2 | Session 3 |
|-----|-----------|-----------|-----------|
| **Day 1** | Squat Day 1 | Bench Day 1 | Deadlift Day 1 |
| **Day 2** | Rest | | |
| **Day 3** | Bench Day 2 | Squat Day 2 | |
| **Day 4** | Rest | | |
| **Day 5** | Bench Day 3 | Deadlift Day 2 | |
| **Day 6** | Rest | | |
| **Day 7** | Rest | | |

Then repeat from Day 1.

## Input Parameters

The program requires only two inputs:

1. **1 Rep Max (1RM)** for each lift:
   - Bench Press 1RM
   - Squat 1RM
   - Deadlift 1RM

2. **Rounding increment**: The smallest weight increment available (typically 1 lb or 2.5 kg)

All weights are calculated from these inputs. Units are user-defined (can be pounds, kilograms, etc.).

---

## Bench Press Routine (3x per Week)

### Structure

The bench press is trained 3 times per week with different rep schemes each day:
- **Day 1**: High reps (8 rep sets)
- **Day 2**: Medium reps (6 rep sets)
- **Day 3**: Low reps (4 rep sets)

### Week 1 Protocol

#### Day 1 (High Rep Day)
| Set Type | Intensity | Sets | Reps |
|----------|-----------|------|------|
| Working Sets | 70% of 1RM | 2 | 8 |
| AMAP Set | 70% of 1RM | 1 | As Many As Possible |

#### Day 2 (Medium Rep Day)
| Set Type | Intensity | Sets | Reps |
|----------|-----------|------|------|
| Working Sets | 75% of 1RM | 2 | 6 |
| AMAP Set | 75% of 1RM | 1 | As Many As Possible |

#### Day 3 (Low Rep Day)
| Set Type | Intensity | Sets | Reps |
|----------|-----------|------|------|
| Working Sets | 80% of 1RM | 2 | 4 |
| AMAP Set | 80% of 1RM | 1 | As Many As Possible |

### Weight Calculation Formula

```
Working Weight = FLOOR(1RM * Intensity Percentage, Rounding Increment)
```

The FLOOR function rounds DOWN to the nearest rounding increment, ensuring the weight is achievable with available plates.

**Example** (100 lb 1RM, 1 lb rounding):
- Day 1: FLOOR(100 * 0.70, 1) = 70 lbs
- Day 2: FLOOR(100 * 0.75, 1) = 75 lbs
- Day 3: FLOOR(100 * 0.80, 1) = 80 lbs

### Week 2+ Progression Logic

From Week 2 onward, weight adjustments are based on AMAP performance from the previous week. Each day has its own progression rules:

#### Day 1 Progression (8-Rep Day)
| AMAP Reps Achieved | Weight Adjustment |
|--------------------|-------------------|
| 8 or fewer | Keep same weight |
| 9-10 reps | Increase by 5 lbs |
| 11-12 reps | Increase by 10 lbs |
| More than 12 | Increase by 15 lbs |

#### Day 2 Progression (6-Rep Day)
| AMAP Reps Achieved | Weight Adjustment |
|--------------------|-------------------|
| 6 or fewer | Keep same weight |
| 7-8 reps | Increase by 5 lbs |
| 9-10 reps | Increase by 10 lbs |
| More than 10 | Increase by 15 lbs |

#### Day 3 Progression (4-Rep Day)
| AMAP Reps Achieved | Weight Adjustment |
|--------------------|-------------------|
| 4 or fewer | Keep same weight |
| 5-6 reps | Increase by 5 lbs |
| 7-8 reps | Increase by 10 lbs |
| More than 8 | Increase by 15 lbs |

### Bench Press Summary

- **Total weekly sets**: 9 (3 sets x 3 days)
- **Rep ranges**: 4-8 prescribed, plus AMAP
- **Intensity range**: 70-80% of 1RM (Week 1)
- **Progression method**: Autoregulated based on AMAP performance

---

## Squat Routine (2x per Week)

### Structure

The squat follows a 4-week periodized cycle with two distinct training days:
- **Day 1**: Volume work at prescribed percentage
- **Day 2**: Rep max (RM) work with back-off sets

### 4-Week Cycle Overview

| Week | Day 1 Intensity | Day 1 Volume | Day 2 RM Goal | Day 2 Back-off Reps |
|------|-----------------|--------------|---------------|---------------------|
| 1 | 75% | 6x6 | 8RM | 5-6 reps x 3 sets |
| 2 | 80% | 5x5 | 5RM | 3-4 reps x 3 sets |
| 3 | 85% | 3x1 | 3RM | 1-2 reps x 3 sets |
| 4 | 70% (deload) | 3x3 | 1RM test | 1 rep x 1 set |

### Day 1 Protocol (Volume Day)

#### Week 1
- **Intensity**: 75% of 1RM
- **Sets**: 6
- **Reps**: 6 per set
- **Weight formula**: FLOOR(1RM * 0.75, Rounding Increment)

#### Week 2
- **Intensity**: 80% of 1RM
- **Sets**: 5
- **Reps**: 5 per set
- **Weight formula**: FLOOR(1RM * 0.80, Rounding Increment)

#### Week 3
- **Intensity**: 85% of 1RM
- **Sets**: 3
- **Reps**: 1 per set (heavy singles)
- **Weight formula**: FLOOR(1RM * 0.85, Rounding Increment)

#### Week 4 (Deload/Test Week)
- **Intensity**: 70% of 1RM
- **Sets**: 3
- **Reps**: 3 per set
- **Weight formula**: FLOOR(1RM * 0.70, Rounding Increment)

### Day 2 Protocol (RM Day)

Day 2 uses rep maxes (RM) rather than percentage-based weights. The lifter works up to a rep max, then performs back-off sets.

#### Week 1
- **Top Set**: Work up to an 8RM (1 set of 8 reps)
- **Back-off Sets**: Same weight, 3 sets of 5-6 reps

#### Week 2
- **Top Set**: Work up to a 5RM (1 set of 5 reps)
- **Back-off Sets**: Same weight, 3 sets of 3-4 reps

#### Week 3
- **Top Set**: Work up to a 3RM (1 set of 3 reps)
- **Back-off Sets**: Same weight, 3 sets of 1-2 reps

#### Week 4
- **Top Set**: Work up to a new 1RM (1 set of 1 rep)
- No back-off sets

### Squat Cycle Summary

- **Total weekly sets**: 7-10 depending on week
- **Intensity progression**: 75% -> 80% -> 85% -> 70% (deload)
- **Volume progression**: Higher reps early, lower reps later
- **Peak**: New 1RM attempt in Week 4

---

## Deadlift Routine (2x per Week)

### Structure

The deadlift follows a 4-week periodized cycle with two training days:
- **Day 1**: Primary stance work (conventional or sumo, whichever is your main stance)
- **Day 2**: Opposite stance work (the other stance for variation and development)

### 4-Week Cycle Overview

| Week | Day 1 Top Set | Day 1 Back-offs | Day 2 Sets | Day 2 Reps |
|------|---------------|-----------------|------------|------------|
| 1 | Heavy set of 5 | 2x3 same weight | 2 | 6 |
| 2 | +5-10 lbs | 3x3 same weight | 3 | 6 |
| 3 | +5-10 lbs | 4x3 same weight | 5 | 5 |
| 4 | Deload (-30%) | 4x3 | 1RM Test | 1 |

### Day 1 Protocol (Primary Stance)

The deadlift uses autoregulated loading rather than strict percentages.

#### Week 1
- **Top Set**: Work up to a heavy (but not max) set of 5 reps
- **Back-off Sets**: 2 sets of 3 reps at the same weight

#### Week 2
- **Top Set**: 5-10 pounds heavier than Week 1's top set, 1 set of 5 reps
- **Back-off Sets**: 3 sets of 3 reps at the same weight

#### Week 3
- **Top Set**: 5-10 pounds heavier than Week 2's top set, 1 set of 5 reps
- **Back-off Sets**: 4 sets of 3 reps at the same weight

#### Week 4 (Deload)
- **Top Set**: 30% lighter than Week 3's top set of 5
- **Back-off Sets**: 4 sets of 3 reps at the deload weight

### Day 2 Protocol (Opposite Stance)

Day 2 uses the opposite deadlift stance for development and to address weaknesses.

#### Week 1
- **Sets**: 2
- **Reps**: 6 per set
- **Intensity guidance**: Should have 2-3 solid reps in the tank at the end of each set

#### Week 2
- **Sets**: 3
- **Reps**: 6 per set
- **Progression**: Same weight as Week 1 for the first set; add 5-10 pounds for each additional set if the preceding set left more than 2 solid reps in the tank

#### Week 3
- **Sets**: 5
- **Reps**: 5 per set
- **Progression**: Same weight as Week 1 for the first set; add 5-10 pounds for each additional set if the preceding set left more than 2 solid reps in the tank

#### Week 4 (Test Week)
- **Exercise**: Primary stance (not opposite stance)
- **Goal**: New 1RM with no form deviations
- **Sets**: 1
- **Reps**: 1

### Deadlift Cycle Summary

- **Total weekly sets**: 3-9 depending on week
- **Progression method**: Autoregulated (RPE-based)
- **Key principle**: Always keep 2-3 reps in reserve on working sets
- **Peak**: New 1RM attempt in Week 4

---

## Volume and Frequency Summary

### Weekly Training Volume (Sets)

| Lift | Week 1 | Week 2 | Week 3 | Week 4 |
|------|--------|--------|--------|--------|
| Bench Press | 9 | 9 | 9 | 9 |
| Squat | 10 | 9 | 7 | 5 |
| Deadlift | 5 | 7 | 10 | 6 |
| **Total** | **24** | **25** | **26** | **20** |

### Training Frequency

| Lift | Sessions/Week |
|------|---------------|
| Bench Press | 3 |
| Squat | 2 |
| Deadlift | 2 |

### Intensity Distribution

| Lift | Low Intensity | Medium Intensity | High Intensity |
|------|---------------|------------------|----------------|
| Bench | 70% (Day 1) | 75% (Day 2) | 80% (Day 3) |
| Squat | 70-75% | 80% | 85%+ / RM work |
| Deadlift | Deload week | Heavy 5s | Near-max work |

---

## Key Algorithm Components

### Weight Rounding Formula

All calculated weights use the FLOOR function to round down to the nearest available increment:

```
Calculated Weight = FLOOR(Base Weight * Percentage, Rounding Increment)
```

This ensures:
1. Weights are always achievable with available plates
2. Rounding is consistent (always down, never up)
3. Users can set their own rounding increment based on available equipment

### Progression Algorithms

#### Bench Press (Autoregulated)
```
IF AMAP_reps <= target_reps THEN
    next_weight = current_weight
ELSE IF AMAP_reps <= target_reps + 2 THEN
    next_weight = current_weight + 5
ELSE IF AMAP_reps <= target_reps + 4 THEN
    next_weight = current_weight + 10
ELSE
    next_weight = current_weight + 15
```

Where target_reps = 8 (Day 1), 6 (Day 2), or 4 (Day 3)

#### Squat (Percentage-Based)
```
Week 1: weight = 1RM * 0.75
Week 2: weight = 1RM * 0.80
Week 3: weight = 1RM * 0.85
Week 4: weight = 1RM * 0.70 (deload)
```

#### Deadlift (RPE-Based)
```
Week 1: Work up to heavy-but-not-max 5
Week 2: Week 1 weight + 5-10 lbs
Week 3: Week 2 weight + 5-10 lbs
Week 4: Week 3 weight * 0.70 (deload before 1RM test)
```

---

## Program Cycling

After completing a 4-week cycle:

1. **Update 1RMs**: If new maxes were achieved in Week 4, update the input values
2. **Restart**: Begin a new 4-week cycle with the updated maxes
3. **Bench adjustments**: The bench press continues its week-to-week autoregulation regardless of cycle position

---

## Implementation Notes

1. **Warm-up sets**: Each session begins with warm-up sets (not specified in detail; lifter's discretion)

2. **Units**: The program is unit-agnostic. Use pounds, kilograms, or any unit consistently

3. **Rest periods**: Not specified; use standard powerlifting rest periods (2-5 minutes for main lifts)

4. **Accessory work**: Not included in the base program; can be added based on individual needs

5. **Failed reps**:
   - Bench: AMAP performance will naturally slow progression
   - Squat/Deadlift: If failing sets, may need to reduce starting max or extend deload periods

---

## Attribution

Program designed by Greg Nuckols. Spreadsheet formatting assistance credited to Reddit user /u/yesmanwriter.

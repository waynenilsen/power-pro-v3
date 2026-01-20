# Calgary Barbell 16-Week Powerlifting Program

## Overview

The Calgary Barbell 16-Week Program is a comprehensive powerlifting meet preparation program designed by Bryce Lewis of Calgary Barbell. It follows a structured periodization model progressing from higher-volume hypertrophy work through strength development to competition peaking, culminating in a detailed taper protocol.

## Program Structure

### Training Schedule

- **Frequency**: 4 training days per week
- **Duration**: 16 weeks total (15 weeks of training + 1 taper week)
- **Split Organization**:
  - Day 1: Squat primary + Bench secondary
  - Day 2: Deadlift primary + Bench variation + Squat secondary
  - Day 3: Squat variation + Bench primary + Deadlift secondary
  - Day 4: Deadlift variation + Bench variation

### Phase Structure

The program divides into distinct phases with different training emphases:

| Phase | Weeks | Focus | Characteristics |
|-------|-------|-------|-----------------|
| Hypertrophy/Base Building | 1-4 | Volume accumulation | Higher reps (5-7), moderate intensity (63-75%), 4-5 sets |
| Strength Development | 5-8 | Intensity introduction | Lower reps (2-5), increasing intensity (65-86%), paired heavy/light sets |
| Peaking | 9-11 | Competition specificity | RPE-regulated top sets with fatigue management |
| Intensification | 12-15 | Max strength | RPE singles/doubles/triples + percentage-based volume |
| Taper | Week 16 | Recovery & priming | Opener practice + reduced volume |

---

## Input Parameters

The program requires three user inputs for individualization:

1. **Squat Training Max**: The user's current squat 1RM or conservative training max
2. **Bench Training Max**: The user's current bench press 1RM or conservative training max
3. **Deadlift Training Max**: The user's current deadlift 1RM or conservative training max
4. **Rounding Increment** (optional): Weight rounding precision (default: 5 lbs/kg)

---

## Load Calculation Algorithm

### Percentage-Based Loading

For exercises with prescribed intensity percentages:

```
Working Load = MROUND(Training Max * Intensity Percentage, Rounding Increment)
```

The MROUND function rounds the calculated load to the nearest specified increment (typically 5 lbs) for practical loading.

### RPE-Based Loading

Exercises marked with "RPE" intensity values (e.g., "8RPE", "9RPE") do not have auto-calculated loads. The lifter self-selects weight based on the target RPE (Rate of Perceived Exertion) scale:

| RPE | Description |
|-----|-------------|
| 7 | Could do 3 more reps |
| 8 | Could do 2 more reps |
| 9 | Could do 1 more rep |
| 10 | Maximum effort, no reps left |

### Fatigue Set Notation

In Weeks 9-15, a special notation indicates fatigue management sets:

- **"1+2F"**: 1 top set followed by 2 "fatigue" back-off sets (reduced weight, same RPE target)
- **"1+3R"**: 1 top set followed by 3 "repeat" sets (same weight as top set)
- **"1+1F"**: 1 top set followed by 1 fatigue set

The distinction:
- **Fatigue sets (F)**: Drop weight to maintain target RPE after accumulated fatigue
- **Repeat sets (R)**: Maintain same weight, allowing RPE to drift higher

---

## Phase-by-Phase Breakdown

### Phase 1: Hypertrophy/Base Building (Weeks 1-4)

#### Intensity Progression

| Week | Competition Lifts Intensity | Variation Intensity |
|------|---------------------------|---------------------|
| 1 | 67% | 60-65% |
| 2 | 70% | 60-70% |
| 3 | 73% | 65-68% |
| 4 | 75% | 70-73% |

#### Rep Scheme Progression

| Week | Main Lift Reps | Sets |
|------|---------------|------|
| 1 | 7 | 4 |
| 2 | 6 | 4 |
| 3 | 6 | 4 |
| 4 | 5 | 5 |

#### Day 1 Structure
1. **Competition Squat**: Primary volume work at prescribed percentage
2. **Paused Bench Press**: Primary bench volume with competition pause
3. **Overhead Press**: 3 sets x 7-8 reps (autoregulated, no percentage)
4. **Bent Over Row**: 3 sets x 8-12 reps (accessory)
5. **GHR/Back Extensions**: 4-5 sets x 10-12 reps (posterior chain)

#### Day 2 Structure
1. **Competition Deadlift**: Primary volume work at prescribed percentage
2. **3-Count Pause Bench**: Long pause variation (60-70% intensity)
3. **SSB/Front/High Bar Pause Squat**: 3 sets x 4-6 reps at 8RPE
4. **Wide Grip Seated Row**: 4-5 sets x 8 reps (bench-specific back work)

#### Day 3 Structure
1. **Pin Squat**: 3 sets x 5-6 reps (65-73% intensity)
2. **2-Board Press**: 3 sets x 4-6 reps at 8RPE
3. **1-Arm Dumbbell Rows**: 5 sets x 8-10 reps
4. **Birddogs**: 3-4 sets x 6-8 reps per side (core stability)

#### Day 4 Structure
1. **2-Count Pause Deadlift**: 3 sets x 5-6 reps (63-73% intensity)
2. **Rep Bench (Touch and Go)**: 4 sets x 8-10 reps (63-70% intensity)
3. **SLDL**: 3 sets x 6-8 reps (40-48% of deadlift max)
4. **Vertical Pull**: 4 sets x 8-10 reps (lat pulldowns or chin-ups)
5. **Tricep Movement**: 4 sets x 8-10 reps (accessory)

---

### Phase 2: Strength Development (Weeks 5-8)

#### Key Changes from Phase 1
- Introduction of **paired heavy/light sets** on main lifts
- Higher intensity on primary sets (80-86%)
- More competition-specific variations
- Reduced accessory volume

#### Intensity Progression

| Week | Heavy Set Intensity | Light Set Intensity |
|------|--------------------|--------------------|
| 5 | 80% | 65-68% |
| 6 | 82% | 68-70% |
| 7 | 86% | 71-72% |
| 8 | 85% | 74-75% |

#### Paired Set Structure (Day 1 & Day 2 Primary Lifts)

Each main lift now uses two distinct rep/intensity brackets:

**Heavy Work:**
- Week 5: 3 sets x 3 reps @ 80%
- Week 6: 4 sets x 3 reps @ 82%
- Week 7: 5 sets x 2 reps @ 86%
- Week 8: 4 sets x 3 reps @ 85%

**Light Work (immediately following):**
- Week 5: 2 sets x 5 reps @ 68%
- Week 6: 2 sets x 5 reps @ 70%
- Week 7: 2 sets x 4 reps @ 72%
- Week 8: 3 sets x 4 reps @ 75%

#### Day 3 Emphasis
- **2-Count Pause Squat**: RPE-regulated (8-9 RPE)
- **Competition Pause Bench**: High volume (6 sets x 3-5 reps @ 68-75%)
- **Feet Up Bench**: 3-4 sets at 8RPE
- **Competition Deadlift**: Light work (65-74%)

#### Day 4 Emphasis
- **2-Count Pause Deadlift**: RPE-regulated
- **Touch and Go Bench**: Higher RPE work (8-10 RPE)
- **Close Grip Incline Press**: 4-5 sets x 6-10 reps
- **1-Arm DB Rows**: 6 sets x 8-10 reps

#### Core Work Addition
- **Side Planks**: 3-4 sets x 30-45 seconds per side

---

### Phase 3: Peaking (Weeks 9-11)

#### Key Characteristics
- Introduction of **fatigue management sets** (1+2F, 1+3R notation)
- Higher intensity on main lifts (82-85%)
- Reduced volume with quality emphasis
- RPE-regulated variation work

#### Intensity Progression

| Week | Main Lift Heavy Sets | Main Lift Light Sets |
|------|---------------------|---------------------|
| 9 | 4-5 sets x 4 @ 82% | 2 sets x 4 @ 71% |
| 10 | 5-6 sets x 3 @ 85% | 2 sets x 4 @ 74% |
| 11 | 3 sets x 3 @ 83% | 2 sets x 4 @ 76% |

#### Day 1 Structure (Weeks 9-11)
1. **Competition Squat**: Primary heavy work + light back-off
2. **Competition Pause Bench**: Primary heavy work + light back-off
3. **Bent Over Row**: 4 sets x 5-8 reps at 8RPE
4. **Rolling Planks**: 3-4 sets x 20-24 reps total

#### Day 2 Structure (Weeks 9-11)
1. **Competition Deadlift**: Primary heavy work + light back-off
2. **Pin Press**: 4 sets x 3-5 reps at 8RPE
3. **Competition Squat**: Light work (68-74%)
4. **Wide Grip Seated Row**: 4 sets x 8 reps

#### Day 3 Structure (Weeks 9-11)
1. **2-Count Pause Squat**: Fatigue set protocol (1+2F/1+3F/1+1R)
2. **Competition Pause Bench**: 5-6 sets x 3-5 reps (72-78%)
3. **Close Grip Bench**: Fatigue set protocol
4. **Competition Deadlift**: Light work (68-74%)
5. **Vertical Pull**: 4 sets x 8 reps

#### Day 4 Structure (Weeks 9-11)
1. **2-Count Pause Deadlift**: Fatigue set protocol
2. **Banded Bench**: 1+1F notation at 9RPE
3. **Barbell Overhead Press**: Fatigue set protocol (1+3R/1+2R/1+1R)
4. **1-Arm DB Rows**: 5 sets x 8 reps

---

### Phase 4: Intensification (Weeks 12-15)

#### Key Characteristics
- **RPE-based top singles/doubles/triples** for competition lifts
- **Percentage-based volume work** calculated from Estimated 1RM (E1RM)
- Progressive reduction in volume
- Maintained high intensity

#### RPE Singles Progression

| Week | Squat | Bench | Deadlift |
|------|-------|-------|----------|
| 12 | 1x3 @ 8RPE | 1x3 @ 8RPE | 1x3 @ 8RPE |
| 13 | 1x2 @ 8RPE | 1x2 @ 8RPE | 1x2 @ 8RPE |
| 14 | 1x1 @ 8-9RPE | 1x1 @ 9RPE | 1x1 @ 8RPE |
| 15 | 1x1 @ 8-9RPE | 1x1 @ 9RPE | 1x1 @ 8RPE |

#### E1RM-Based Volume Work

Following the RPE top sets, percentage-based volume is calculated from the Estimated 1RM derived from the top set performance:

| Week | Sets x Reps | % of E1RM |
|------|-------------|-----------|
| 12 | 6-7 x 5 | 65% |
| 13 | 6-7 x 5 | 68% |
| 14 | 4-5 x 4 | 72% |
| 15 | 3-4 x 3 | 76% |

#### Variation Work Structure

All variation work in Weeks 12-15 uses fatigue set protocols with RPE regulation:

**Day 1:**
- Overhead Press: Fatigue sets at 8-9 RPE

**Day 2:**
- 2-Count Pause Bench: Fatigue sets at 8-9 RPE
- High Bar Squat: Fatigue sets at 8-9 RPE

**Day 3:**
- Pin Squat: RPE top set + fatigue sets
- Close Grip Bench: Fatigue sets
- Feet Up Bench: Fatigue sets

**Day 4:**
- 2-Count Pause Deadlift: RPE top set + fatigue sets
- Touch and Go Bench: Fatigue sets

---

### Phase 5: Taper Week (Week 16 - Competition Week)

The taper week is structured around days-out from competition:

#### 5 Days Out (Day 1)
| Exercise | Sets | Reps | Load |
|----------|------|------|------|
| Competition Squat | 1 | 1 | Opener |
| Competition Squat | 3 | 2 | 82% |
| Competition Pause Bench | 1 | 1 | Opener |
| Competition Pause Bench | 3 | 2 | 84% |

**Purpose**: Practice openers; moderate volume at sub-maximal intensity.

#### 4 Days Out (Day 2)
| Exercise | Sets | Reps | Load |
|----------|------|------|------|
| Competition Deadlift | 1 | 1 | Opener |
| Competition Deadlift | 2 | 2 | 82% |
| Competition Pause Bench | 4 | 1 | 85% |

**Purpose**: Deadlift opener practice; bench nervous system priming.

#### 3 Days Out (Day 3)
| Exercise | Sets | Reps | Load |
|----------|------|------|------|
| Competition Squat | 1 | 1 | 85% |
| Competition Squat | 2 | 2 | 78% |
| Competition Pause Bench | 1 | 1 | 85% |
| Competition Pause Bench | 3 | 2 | 78% |
| Competition Deadlift | 1 | 1 | 82% |
| Competition Deadlift | 2 | 2 | 75% |

**Purpose**: Final heavy touches on all three lifts; groove competition movement patterns.

#### 2 Days Out (Day 4)
| Exercise | Sets | Reps | Load |
|----------|------|------|------|
| Competition Squat | 2 | 3 | 75% |
| Competition Pause Bench | 3 | 3 | 78% |
| Competition Deadlift | 1 | 3 | 75% |

**Purpose**: Light movement practice; maintain neuromuscular activation without fatigue.

#### 1 Day Out
Complete rest.

#### Meet Day
Competition.

---

## Tempo Prescription

The program uses tempo notation in format: **Eccentric.Pause.Concentric**

| Notation | Meaning |
|----------|---------|
| 1.0.1 | 1 second down, no pause, 1 second up (competition squat/deadlift) |
| 1.1.1 | 1 second down, 1 second pause, 1 second up (competition bench) |
| 1.2.1 | 1 second down, 2 second pause, 1 second up (pause variations) |
| 1.3.1 | 1 second down, 3 second pause, 1 second up (long pause bench) |
| x | No specific tempo prescribed |

---

## Rest Period Guidelines

| Exercise Type | Rest Period (seconds) |
|---------------|----------------------|
| Competition lifts (heavy) | 180 |
| Competition lifts (light) | 180 |
| Variation lifts | 120-180 |
| Accessory work | 60-90 |
| Core work | 60 |

---

## Volume Progression Summary

### Weekly Set Volume Trends (Competition Lifts Only)

| Phase | Squat Sets/Week | Bench Sets/Week | Deadlift Sets/Week |
|-------|-----------------|-----------------|-------------------|
| Weeks 1-4 | 10-14 | 14-18 | 10-13 |
| Weeks 5-8 | 11-15 | 18-22 | 11-14 |
| Weeks 9-11 | 12-16 | 16-20 | 10-13 |
| Weeks 12-15 | 8-12 | 12-16 | 8-10 |
| Taper | 6-7 | 8-12 | 4-5 |

### Intensity Zones by Phase

| Phase | Primary Zone | Secondary Zone |
|-------|--------------|----------------|
| Hypertrophy | 67-75% | 60-73% |
| Strength | 80-86% | 65-75% |
| Peaking | 82-85% | 68-78% |
| Intensification | 8-9 RPE | 65-76% |
| Taper | 75-85% | Openers |

---

## Exercise Selection Rationale

### Competition Movement Patterns
- **Competition Squat**: Low bar back squat to competition depth
- **Competition Pause Bench**: Bench press with competition-length pause
- **Competition Deadlift**: Conventional or sumo (athlete's competition stance)

### Primary Variations
| Variation | Purpose |
|-----------|---------|
| Pin Squat | Concentric strength, bottom position power |
| Pause Squat (2ct) | Positional strength, rate of force development |
| Pause Bench (3ct) | Bottom position strength, tightness |
| 2-Board Press | Lockout strength, overload capacity |
| Pause Deadlift (2ct) | Breaking the floor, positional strength |
| SLDL | Posterior chain development, hamstring emphasis |

### Accessory Classifications
| Category | Exercises | Purpose |
|----------|-----------|---------|
| Horizontal Pull | Bent Over Row, Seated Row, DB Rows | Upper back development, bench stability |
| Vertical Pull | Lat Pulldowns, Chin-ups | Lat development, deadlift assistance |
| Pressing Assistance | Overhead Press, Close Grip, Feet Up | Tricep/shoulder development |
| Posterior Chain | GHR, Back Extensions, SLDL | Hamstring/glute development |
| Core | Planks, Birddogs, Rolling Planks | Trunk stability |

---

## Implementation Notes

### Auto-Regulation Integration

The program blends percentage-based and RPE-based loading:

1. **Weeks 1-8**: Primarily percentage-based with RPE for variations
2. **Weeks 9-11**: Mixed approach with fatigue management protocols
3. **Weeks 12-15**: RPE top sets dictate subsequent percentage work
4. **Taper**: Fixed percentages for consistency and predictability

### Weight Selection for RPE Work

For exercises without calculated loads:
1. Start conservatively in Week 1
2. Track weights used and RPE achieved
3. Progressively increase loads while maintaining target RPE
4. Use the same weight or slightly less when reps decrease

### Fatigue Set Implementation

When encountering notation like "1+2F":
1. Work up to a top set at target RPE
2. For fatigue sets: reduce weight 5-10% to maintain same RPE
3. For repeat sets: keep same weight, allow RPE to increase

### Opener Selection

Openers should be:
- A weight you could triple on your worst day
- Approximately 90-92% of projected competition max
- Practiced multiple times during the taper

---

## Summary

The Calgary Barbell 16-Week Program provides a systematic approach to powerlifting meet preparation through:

1. **Progressive overload** via increasing intensity across phases
2. **Specificity progression** from variations to competition movements
3. **Fatigue management** through strategic deload patterns within the volume work
4. **RPE integration** for auto-regulation of training stress
5. **Structured taper** to optimize performance on meet day

The program balances sufficient volume for adaptation with appropriate intensity exposure for strength expression, making it suitable for intermediate to advanced powerlifters preparing for competition.

# Calgary Barbell 8-Week Peaking Program

## Program Overview

The Calgary Barbell 8-Week Peaking Program is a condensed powerlifting meet preparation program designed to bring a lifter to peak performance for competition. This is a shortened version of the more common 16-week Calgary Barbell program, designed for athletes who have limited time before a meet or who respond well to shorter peaking cycles.

### Key Characteristics
- **Duration**: 8 weeks plus a competition taper week
- **Frequency**: 4 training days per week
- **Structure**: Two distinct 4-week phases followed by a taper
- **Target**: Powerlifting competition preparation (squat, bench, deadlift)

---

## Program Inputs

### Training Maxes
The program requires three input values:
1. **Squat Training Max** (in pounds or kilograms)
2. **Bench Training Max** (in pounds or kilograms)
3. **Deadlift Training Max** (in pounds or kilograms)

### Rounding Parameter
A configurable rounding value (default: 5 lbs/kg) ensures all calculated loads round to practical plate increments.

### Load Calculation Formula
All percentage-based loads are calculated as:
```
Load = ROUND(Training Max * Intensity Percentage, Rounding Value)
```

---

## Phase Structure

### Phase 1: Weeks 1-4 (Accumulation/Hypertrophy)

This phase emphasizes building work capacity through percentage-based work on competition lifts combined with RPE-regulated accessory variations.

#### Characteristics:
- Competition lifts performed with prescribed percentages
- Dual-tier system: Heavy work (lower reps, higher intensity) followed by volume work (higher reps, lower intensity)
- Accessory movements use RPE for autoregulation
- Higher total volume compared to Phase 2
- Tempo prescriptions for all movements

### Phase 2: Weeks 5-8 (Intensification/Peaking)

This phase transitions to heavier loads with RPE-based top sets followed by percentage-based back-off work.

#### Characteristics:
- Competition lifts begin with RPE-based heavy singles, doubles, or triples
- Back-off sets calculated from Estimated 1RM (E1RM) based on RPE top set
- Introduction of fatigue-based set notation (e.g., "1+2F")
- Reduced overall volume with maintained or increased intensity
- Progressive reduction of reps on heavy work (3 -> 2 -> 1 -> 1)

### Taper Week: Final Competition Preparation

A 4-day taper in the final week before competition, with sessions organized by days-out from the meet.

---

## Weekly Training Schedule

### Phase 1 (Weeks 1-4) - Day Structure

#### Day 1: Squat/Bench Focus
| Exercise | Purpose |
|----------|---------|
| Competition Squat (Heavy) | Primary squat work, higher intensity |
| Competition Squat (Volume) | Back-off squat work, lower intensity |
| Competition Pause Bench (Heavy) | Primary bench work, higher intensity |
| Competition Pause Bench (Volume) | Back-off bench work, lower intensity |
| SLDL | Posterior chain accessory (RPE-based) |
| Side Planks | Core stability (timed) |

#### Day 2: Deadlift/Bench Variation
| Exercise | Purpose |
|----------|---------|
| Competition Deadlift (Heavy) | Primary deadlift work |
| Competition Deadlift (Volume) | Back-off deadlift work |
| 2ct Pause Bench | Bench variation (RPE-based) |
| Competition Squat (Light) | Technical practice, lower intensity |
| Wide Grip Seated Row | Back accessory |

#### Day 3: Squat Variation/Bench Volume
| Exercise | Purpose |
|----------|---------|
| 2ct Pause Squat | Squat variation (RPE-based) |
| Competition Pause Bench (Volume) | High-volume bench work |
| Feet Up Bench | Bench variation (RPE-based) |
| Competition Deadlift (Light) | Technical practice |
| Vertical Pull | Back accessory |

#### Day 4: Deadlift Variation/Bench Accessories
| Exercise | Purpose |
|----------|---------|
| 2ct Pause Deadlifts | Deadlift variation (RPE-based) |
| Touch and Go Bench | Bench variation (RPE-based) |
| Close Grip Incline Press | Tricep/pressing accessory |
| 1-Arm DB Rows | Back accessory |

---

### Phase 2 (Weeks 5-8) - Day Structure

#### Day 1: Squat/Bench Peak
| Exercise | Purpose |
|----------|---------|
| Competition Squat (RPE Top Set) | Heavy single/double/triple at RPE |
| Competition Squat (% of E1RM) | Back-off volume from top set |
| Competition Pause Bench (RPE Top Set) | Heavy single/double/triple at RPE |
| Competition Pause Bench (% of E1RM) | Back-off volume from top set |
| Overhead Press | Pressing accessory (fatigue sets) |

#### Day 2: Deadlift Peak
| Exercise | Purpose |
|----------|---------|
| Competition Deadlift (RPE Top Set) | Heavy single/double/triple at RPE |
| Competition Deadlift (% of E1RM) | Back-off volume from top set |
| 2ct Pause Bench | Bench variation (fatigue sets) |
| High Bar Squat | Squat variation (fatigue sets) |

#### Day 3: Variation Day
| Exercise | Purpose |
|----------|---------|
| Pin Squat (RPE Top Set) | Squat variation, overload focus |
| Pin Squat (Fatigue Sets) | Volume work |
| Close Grip Bench | Bench variation (fatigue sets) |
| Feet Up Bench | Bench variation (fatigue sets) |

#### Day 4: Deadlift Variation
| Exercise | Purpose |
|----------|---------|
| 2ct Pause Deadlifts (RPE Top Set) | Deadlift variation |
| 2ct Pause Deadlifts (Fatigue Sets) | Volume work |
| Touch and Go Bench | Bench variation (fatigue sets) |

---

## Intensity Progression

### Phase 1 Percentage Progression

#### Competition Lifts - Heavy Work
| Week | Sets | Reps | Intensity |
|------|------|------|-----------|
| 1 | 3-4 | 3 | 80% |
| 2 | 4-5 | 3 | 82% |
| 3 | 5 | 2 | 86% |
| 4 | 4-5 | 3 | 85% |

#### Competition Lifts - Volume Work
| Week | Sets | Reps | Intensity |
|------|------|------|-----------|
| 1 | 2 | 5 | 68% |
| 2 | 2-3 | 5 | 70% |
| 3 | 2 | 4 | 72% |
| 4 | 2-3 | 4 | 75% |

#### Secondary Volume (Day 2/3 Light Work)
| Week | Sets | Reps | Intensity |
|------|------|------|-----------|
| 1 | 2 | 5 | 65% |
| 2 | 3 | 5 | 68% |
| 3 | 2 | 5 | 71% |
| 4 | 2 | 4 | 74% |

#### Day 3 Bench Volume
| Week | Sets | Reps | Intensity |
|------|------|------|-----------|
| 1 | 6 | 5 | 70% |
| 2 | 6 | 4 | 73% |
| 3 | 6 | 3 | 75% |
| 4 | 6 | 5 | 68% |

Note: Week 4 acts as a partial deload with reduced intensity on volume work before Phase 2.

### Phase 2 RPE/Percentage Progression

#### Competition Lift Top Sets
| Week | Sets | Reps | RPE |
|------|------|------|-----|
| 5 | 1 | 3 | 8 |
| 6 | 1 | 2 | 8 |
| 7 | 1 | 1 | 8-9 |
| 8 | 1 | 1 | 8-9 |

#### Competition Lift Back-Off Sets (% of E1RM)
| Week | Sets | Reps | Intensity |
|------|------|------|-----------|
| 5 | 6-7 | 5 | 65% |
| 6 | 6-7 | 5 | 68% |
| 7 | 4-5 | 4 | 72% |
| 8 | 3-4 | 3 | 76% |

---

## RPE Usage

### Phase 1 RPE Guidelines

RPE is used for accessory movements and variations:

| RPE | Description | Usage |
|-----|-------------|-------|
| 8 | Could do 2 more reps | Most accessory work |
| 9 | Could do 1 more rep | Heavy variation work |
| 10 | Maximum effort | Occasional use on TnG bench |

### Phase 2 RPE and Fatigue Set System

#### Top Set RPE
- **RPE 8**: Work up to a weight where 2 more reps could be completed
- **RPE 9**: Work up to a weight where 1 more rep could be completed

#### Fatigue Set Notation
Phase 2 uses a unique fatigue-based set notation:

| Notation | Meaning |
|----------|---------|
| 1+2F | 1 top set, then add weight for 2 more sets (Fatigue ascending) |
| 1+2R | 1 top set, then reduce weight for 2 more sets (Fatigue descending) |
| 1+1F | 1 top set, then add weight for 1 more set |
| 1+3R | 1 top set, then reduce weight for 3 more sets |

**F (Forward/Fatigue)**: After the top set, increase weight slightly for subsequent sets, accepting that fatigue will make these harder.

**R (Reduce)**: After the top set, reduce weight slightly to maintain quality reps as fatigue accumulates.

---

## Peaking Protocol

### Volume Tapering
The program uses a progressive volume reduction approach:

| Phase | Relative Volume |
|-------|-----------------|
| Weeks 1-3 | High (building) |
| Week 4 | Moderate (transition) |
| Weeks 5-6 | Moderate-High (intensification) |
| Weeks 7-8 | Reduced (peaking) |
| Taper Week | Minimal (competition prep) |

### Intensity Peaking
Peak intensity is achieved through:
1. Progressive reduction in competition lift reps (3 -> 2 -> 1)
2. Maintenance of RPE at 8-9 on heavy work
3. Increase in back-off percentages (65% -> 76%)

### Taper Week Protocol

The final week uses a countdown format based on days from competition:

#### 5 Days Out
- Competition Squat: 1x1 @ Opener, then 3x2 @ 82%
- Competition Pause Bench: 1x1 @ Opener, then 3x2 @ 84%

#### 4 Days Out
- Competition Deadlift: 1x1 @ Opener, then 2x2 @ 82%
- Competition Pause Bench: 4x1 @ 85%

#### 3 Days Out
- Competition Squat: 1x1 @ 85%, 2x2 @ 78%
- Competition Pause Bench: 1x1 @ 85%, 3x2 @ 78%
- Competition Deadlift: 1x1 @ 82%, 2x2 @ 75%

#### 2 Days Out
- Competition Squat: 2x3 @ 75%
- Competition Pause Bench: 3x3 @ 78%
- Competition Deadlift: 1x3 @ 75%

#### 1 Day Out
- Rest

---

## Tempo Notation

The program uses a 3-digit tempo system:

| Position | Meaning |
|----------|---------|
| First digit | Eccentric (lowering) time in seconds |
| Second digit | Pause at bottom in seconds |
| Third digit | Concentric (lifting) time in seconds |

### Standard Tempos Used
| Tempo | Application |
|-------|-------------|
| 1.0.1 | Standard squat/deadlift (controlled eccentric, no pause, explosive concentric) |
| 1.1.1 | Competition pause bench (controlled eccentric, 1-second pause, explosive concentric) |
| 1.2.1 | Pin squats, 2ct pause bench (longer pause) |
| 1.3.1 | 2ct pause bench variation (extended pause) |

---

## Rest Periods

| Exercise Type | Rest Period (seconds) |
|---------------|----------------------|
| Competition Lifts (Heavy) | 180 |
| Competition Lifts (Volume) | 180 |
| RPE-Based Variations | 180 |
| Accessories | 60-120 |
| Core Work | 60 |

---

## Volume/Intensity Progression Summary

### Competition Squat Weekly Volume (Estimated Working Sets)
| Week | Heavy Sets | Volume Sets | Total |
|------|------------|-------------|-------|
| 1 | 3 | 4 | 7 |
| 2 | 4 | 5 | 9 |
| 3 | 5 | 4 | 9 |
| 4 | 4 | 5 | 9 |
| 5 | 1 + 6 | - | 7 |
| 6 | 1 + 6 | - | 7 |
| 7 | 1 + 4 | - | 5 |
| 8 | 1 + 3 | - | 4 |

### Key Periodization Principles

1. **Dual-Quality Training**: Phase 1 separates heavy and volume work within the same session
2. **RPE-to-Percentage Transition**: Phase 2 uses RPE for top sets with percentage-based back-offs
3. **Fatigue Management**: Week 4 and Week 8 reduce volume before transitions
4. **Specificity Increase**: Variations decrease as competition approaches
5. **Opener Practice**: Taper week includes opener attempts for competition preparation

---

## Implementation Notes

### Training Max Selection
- Use approximately 95% of true 1RM for training max inputs
- Alternatively, use a recent competition PR or gym max

### Autoregulation Guidelines
- If RPE targets cannot be hit, reduce load by 2.5-5%
- If sets feel significantly easier than prescribed RPE, increase load by 2.5%
- Track estimated 1RMs from Phase 2 top sets to gauge peaking progress

### Exercise Substitutions
The following substitutions maintain program intent:

| Prescribed | Acceptable Alternatives |
|------------|------------------------|
| Pin Squat | Box Squat, Anderson Squat |
| High Bar Squat | Safety Bar Squat |
| 2ct Pause Deadlift | Deficit Deadlift |
| Feet Up Bench | Larsen Press |
| Close Grip Bench | 2-Board Press |

### Accessories
Accessory work without prescribed intensity should be performed at a moderate effort that allows completion of prescribed sets and reps while leaving 2-3 reps in reserve.

---

## Condensed vs. 16-Week Program

### Key Differences
| Aspect | 8-Week | 16-Week |
|--------|--------|---------|
| Accumulation Phase | 4 weeks | 8 weeks |
| Peaking Phase | 4 weeks | 8 weeks |
| Volume Build | Compressed | Gradual |
| Best For | Experienced lifters, quick peaking | Newer lifters, major meets |

### When to Use 8-Week Version
- Time constraints before competition
- Lifter responds well to shorter cycles
- Maintaining fitness between close competitions
- Experienced lifters with established strength base

---

## Program Attribution

This program is based on the Calgary Barbell methodology, a respected powerlifting coaching system. The spreadsheet implementation was distributed through LiftVault.com.

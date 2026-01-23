# Power Pro Roadmap

## Core Domain Objects

### Lift
The exercise itself (Squat, Bench, Deadlift, etc.)

### LiftMax
A user's reference number for a lift:
- `1RM` - true one-rep max
- `TM` - training max (usually 85-90% of 1RM)
- `xRM` - rep max (5RM, 6RM, 8RM, etc.)
- `E1RM` - estimated 1RM calculated from a performed set

### Prescription
What you're supposed to do for a single exercise slot:
- Which lift
- How to determine weight (LoadStrategy)
- How many sets/reps (SetScheme)

### LoadStrategy
How to calculate the weight:
- `PercentOf(reference, percentage)` - e.g., 85% of TM
- `RPETarget(reps, rpe)` - use RPE chart lookup
- `LinearAdd(last_weight, increment)` - last time + 5lb
- `FindRM(target_rm)` - no prescription, user discovers it
- `RelativeTo(set_id, percentage)` - % of today's top set

### SetScheme
How sets and reps are structured:
- `Fixed(sets, reps)` - 5×5
- `Ramp(percentages, reps)` - [50%, 63%, 75%, 88%, 100%] × 5
- `TopBackoff(top_sets, top_reps, backoff_sets, backoff_reps, backoff_percent)`
- `AMRAP(sets, min_reps)` - 1×5+
- `MRS(initial_reps, max_sets)` - max rep sets until failure
- `FatigueDrop(drop_percent, stop_rpe)` - repeat until RPE hits target
- `TotalReps(target)` - 100 reps however you want

### Progression
Rules that mutate LiftMax values:
- `LinearProgression(increment, frequency)` - +5lb per session/week/cycle
- `AMRAPProgression(thresholds)` - if reps > X then +Y
- `DeloadOnFailure(failure_count, multiplier)` - after 2 fails, ×0.9
- `StageProgression(stages)` - 5×3 → 6×2 → 10×1 on failure
- `DoubleProgression(rep_range, increment)` - add reps until ceiling, then add weight
- `CycleProgression(increment)` - +5lb at end of 4-week cycle

### Schedule
When and what:
- `Day` - a named training day with exercise slots
- `Week` - a collection of days
- `Cycle` - repeating unit (1 week, 3 weeks, 4 weeks, etc.)
- `Phase` - a block with different parameters
- `Rotation` - A/B alternation, or 3-week lift rotation

### Trigger
When a progression rule fires:
- `AfterSet(filter)` - after specific set type (AMRAP, top set)
- `AfterSession`
- `AfterWeek`
- `AfterCycle`
- `OnFailure`

### Lookup
Tables that vary prescriptions:
- `WeeklyLookup` - week 1 = 65/75/85%, week 2 = 70/80/90%
- `DailyLookup` - Monday = 75%, Tuesday = 80%
- `RPEChart` - (reps, RPE) → percentage

### LoggedSet
What actually happened:
- lift, weight, reps, rpe, timestamp
- Used to calculate E1RM, trigger progressions

---

## Implementation Phases

### Phase 1: Core
**Features:**
- Lift, LiftMax (1RM, TM)
- Prescription
- LoadStrategy: `PercentOf`
- SetScheme: `Fixed`, `Ramp`
- Progression: `LinearProgression`, `CycleProgression`
- Schedule: Day, Week, Cycle
- Lookup: `WeeklyLookup`, `DailyLookup`
- Rounding to configurable increment

**Unlocks:**
| Program | Why |
|---------|-----|
| Starting Strength | Fixed 3×5, LinearProgression per session. A/B alternation is basic Week/Cycle functionality (not Rotation feature). Reset protocols are manual user parameters, not automatic failure tracking. |
| Bill Starr 5x5 | Heavy/Light/Medium days, Ramp sets, LinearProgression per week. Reset/deload are manual user actions, not system features. |
| Sheiko Beginner | Many Fixed sets at various %, no autoregulation needed. Prep phase (Weeks 1-3) fully implementable; Comp phase taper would benefit from Phase 10 but not required. |

**Programs unlocked: 3**

**Note:** Phases are cumulative - each phase includes all features from previous phases. Programs listed here can be **fully implemented** with Phase 1 features.

---

### Phase 2: AMRAP
**Features:**
- SetScheme: `AMRAP`
- LoggedSet (record actual reps)
- Progression: `AMRAPProgression`
- Trigger: `AfterSet(amrap)`

**Unlocks:**
| Program | Why |
|---------|-----|
| Wendler 5/3/1 BBB | AMRAP final sets (5+, 3+, 1+), WeeklyLookup for 4-week wave, CycleProgression |
| Greg Nuckols High Frequency | AMAP sets in Week 3, DailyLookup + WeeklyLookup, CycleProgression |
| nSuns 5/3/1 LP 5-Day | Multiple AMRAP sets, complex AMRAPProgression thresholds |

**Programs unlocked: 3 (Total: 6)**

---

### Phase 3: Failure Handling
**Features:**
- Failure tracking (reps < target)
- Progression: `DeloadOnFailure`, `StageProgression`
- Trigger: `OnFailure`

**Unlocks:**
| Program | Why |
|---------|-----|
| GZCLP | StageProgression (5×3 → 6×2 → 10×1), DeloadOnFailure. Requires AMRAP (Phase 2) + Failure Handling (Phase 3). Phases are cumulative, so Phase 3 includes all Phase 2 features. |
| Texas Method | Implicit failure handling, reset protocols. No AMRAP required. |
| Greg Nuckols Beginner | AMRAP drives weekly TM adjustment, requires failure handling (Phase 3). Phases are cumulative, so Phase 3 includes all Phase 2 features. |

**Programs unlocked: 3 (Total: 9)**

---

### Phase 4: Double Progression
**Features:**
- Progression: `DoubleProgression`
- SetScheme: `RepRange(min, max)`

**Unlocks:**
| Program | Why |
|---------|-----|
| Reddit PPL 6-Day | Requires double progression for accessories (Phase 4) and failure handling (Phase 3, included in Phase 4). Phases are cumulative, so Phase 4 includes all Phase 2 and Phase 3 features. |

**Programs unlocked: 1 (Total: 10)**

---

### Phase 5: Rotation
**Features:**
- Schedule: `Rotation`
- Exercise slot rotation (which lift is "primary" this week)
- Cycle-based lift focus changes

**Unlocks:**
| Program | Why |
|---------|-----|
| nSuns CAP3 | 3-week rotation of which lift gets AMRAP focus |
| Inverted Juggernaut 5/3/1 | 16-week wave rotation (10s/8s/5s/3s phases) |
| Greyskull LP | AMRAP final set, AMRAPProgression, requires failure handling (Phase 3) and rotation (Phase 5). Phases are cumulative, so Phase 5 includes all Phase 2 and Phase 3 features. |

**Programs unlocked: 3 (Total: 13)**

---

### Phase 6: RPE
**Features:**
- LoadStrategy: `RPETarget`
- Lookup: `RPEChart`
- LoggedSet extended with RPE field

**Unlocks:**
| Program | Why |
|---------|-----|
| RTS Intermediate (partial) | RPETarget for load calculation, RPEChart lookup |

**Programs unlocked: 1 partial (Total: 13)**

---

### Phase 7: E1RM
**Features:**
- LiftMax: `E1RM`
- LoadStrategy: `FindRM`, `RelativeTo`
- E1RM calculation from LoggedSet
- Prescription chaining (today's top set feeds back-off prescription)

**Unlocks:**
| Program | Why |
|---------|-----|
| GZCL Jacked & Tan 2.0 | FindRM (10RM→8RM→6RM etc), back-offs at % of found RM |
| Calgary Barbell 8-Week | Phase 2: RPE top set → E1RM → RelativeTo back-offs |
| Calgary Barbell 16-Week | Same as 8-week, extended phases |

**Programs unlocked: 3 (Total: 16)**

---

### Phase 8: Fatigue Protocols
**Features:**
- SetScheme: `FatigueDrop`, `MRS`
- Variable set count (unknown until session complete)
- Trigger: `AfterSet` with RPE condition
- "Repeat until" logic

**Unlocks:**
| Program | Why |
|---------|-----|
| GZCL Compendium | VDIP, MRS protocols, fatigue-based volume |
| RTS Intermediate (complete) | Load drop method, repeat sets method |

**Programs unlocked: 2 (Total: 18)**

---

### Phase 9: Volume Targets
**Features:**
- SetScheme: `TotalReps`
- Session-level rep tracking
- Superset notation

**Unlocks:**
| Program | Why |
|---------|-----|
| 5/3/1 Building the Monolith | "100 chin-ups", "100-200 dips" however you get them |

**Programs unlocked: 1 (Total: 19)**

---

### Phase 10: Peaking
**Features:**
- Schedule: `DaysOut` (countdown to meet)
- Phase transitions based on calendar
- Opener practice sessions
- Taper protocols

**Unlocks:**
| Program | Why |
|---------|-----|
| Sheiko Intermediate | 13-week with prep/comp phases, max testing protocols |

**Programs unlocked: 1 (Total: 20)**

---

## Summary

| Phase | Key Features | Programs Unlocked | Running Total |
|-------|--------------|-------------------|---------------|
| 1 | Core (%, Fixed, Ramp, Linear) | Starting Strength, Bill Starr, Sheiko Beginner | 3 |
| 2 | AMRAP | Wendler 5/3/1 BBB, Greg Nuckols HF, nSuns 5/3/1 | 6 |
| 3 | Failure Handling | GZCLP, Texas Method, Greg Nuckols Beginner | 9 |
| 4 | Double Progression | Reddit PPL 6-Day | 10 |
| 5 | Rotation | nSuns CAP3, Inverted Juggernaut, Greyskull LP | 13 |
| 6 | RPE | RTS (partial) | 13 |
| 7 | E1RM | J&T 2.0, Calgary 8wk, Calgary 16wk | 16 |
| 8 | Fatigue Protocols | GZCL Compendium, RTS (complete) | 18 |
| 9 | Volume Targets | Building the Monolith | 19 |
| 10 | Peaking | Sheiko Intermediate | 20 |

---

## Quick Reference

| Concept | Examples |
|---------|----------|
| **LoadStrategy** | PercentOf, RPETarget, LinearAdd, FindRM, RelativeTo |
| **SetScheme** | Fixed, Ramp, AMRAP, MRS, FatigueDrop, TotalReps |
| **Progression** | Linear, AMRAP-driven, Deload, Stage, Double, Cycle |
| **Trigger** | AfterSet, AfterSession, AfterWeek, AfterCycle, OnFailure |
| **Lookup** | WeeklyLookup, DailyLookup, RPEChart |
| **Schedule** | Day, Week, Cycle, Phase, Rotation, DaysOut |

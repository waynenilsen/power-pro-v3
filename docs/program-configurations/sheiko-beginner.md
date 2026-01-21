# Sheiko Beginner Program Configuration

This document provides the complete API configuration for the Sheiko Beginner program, demonstrating high-volume percentage-based training with **manual (no auto) progression**.

## Table of Contents

- [Overview](#overview)
- [Program Characteristics](#program-characteristics)
- [Configuration Components](#configuration-components)
  - [Step 1: Create Lifts](#step-1-create-lifts)
  - [Step 2: Create Prescriptions (High-Volume Sets)](#step-2-create-prescriptions-high-volume-sets)
  - [Step 3: Create Days](#step-3-create-days)
  - [Step 4: Create Cycle](#step-4-create-cycle)
  - [Step 5: Create Weeks](#step-5-create-weeks)
  - [Step 6: Create Program (No Progressions)](#step-6-create-program-no-progressions)
- [User Setup](#user-setup)
- [Example Workout Output](#example-workout-output)
  - [Prep Phase Week 1 Day 1](#prep-phase-week-1-day-1)
  - [Prep Phase Week 1 Day 3](#prep-phase-week-1-day-3)
  - [Prep Phase Week 1 Day 5](#prep-phase-week-1-day-5)
- [Manual Progression Approach](#manual-progression-approach)
  - [When to Update Maxes](#when-to-update-maxes)
  - [How to Update Maxes via API](#how-to-update-maxes-via-api)
  - [Contrast with Auto-Progression Programs](#contrast-with-auto-progression-programs)
- [Complete Configuration Summary](#complete-configuration-summary)

---

## Overview

The Sheiko Beginner program is a structured powerlifting training system designed by Boris Sheiko, the legendary Russian powerlifting coach. This program emphasizes high frequency, moderate intensity, and technical mastery through submaximal training.

**Key Characteristics:**
- High volume: 20-30+ reps per lift per day
- Moderate intensity: Most work at 70-85% of 1RM
- 3 training days per week (M/W/F)
- Percentage-based loading with RAMP set schemes
- **No automatic progression** - coach/athlete adjusts maxes manually
- 4-week Prep phase (repeatable) + optional 4-week Comp phase

**What This Configuration Demonstrates:**
- **RAMP set schemes** with built-in warm-up sets
- **High-volume prescriptions** - multiple sets at various percentages
- **No progression rules** - manual max adjustments only
- **Manual max update API calls** - how users increase their maxes when ready

---

## Program Characteristics

| Property | Value |
|----------|-------|
| Cycle Length | 4 weeks (Prep Phase) |
| Days per Week | 3 (M/W/F) |
| Progression Type | **NONE (Manual)** |
| Intensity Range | 50-90% (Prep), 50-105% (Comp) |
| Primary Focus | Technical perfection, volume accumulation |

### Sheiko Philosophy

1. **Technical Perfection** - Weights should challenge technique but allow perfect form
2. **High Frequency** - Each competition lift trained 2-3 times per week
3. **Submaximal Training** - Most work at 70-85% of 1RM
4. **Volume Accumulation** - Progress through accumulated volume, not maximal efforts
5. **Manual Progression** - Athlete/coach decides when maxes should increase

### Program Structure (Prep Phase)

| Day | Focus | Volume |
|-----|-------|--------|
| Day 1 (Monday) | Squat + Bench | High (40+ total reps main lifts) |
| Day 3 (Wednesday) | Deadlift variations + Bench | Moderate (35+ total reps) |
| Day 5 (Friday) | Squat + Bench variations | Moderate (35+ total reps) |

### Rep and Set Schemes

Competition lifts use ramping percentages with built-in warm-ups:

| Percentage | Typical Reps |
|------------|--------------|
| 50% | 3-5 (warm-up) |
| 60% | 3-5 (warm-up) |
| 70% | 3-4 (work) |
| 75% | 2-3 (work) |
| 80% | 2-3 (work) |
| 85% | 2 (work) |
| 90% | 1-2 (heavy) |

---

## Configuration Components

All API calls below assume:
- Base URL: `http://localhost:8080`
- Admin authentication: `X-Admin: true`
- User ID header: `X-User-ID: {userId}`

### Step 1: Create Lifts

Create the main competition lifts and variations used in Sheiko.

```bash
# Squat (Competition)
curl -X POST "http://localhost:8080/lifts" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Squat",
    "slug": "squat",
    "isCompetitionLift": true
  }'
```

**Response:**
```json
{
  "data": {
    "id": "lift-squat-uuid",
    "name": "Squat",
    "slug": "squat",
    "isCompetitionLift": true,
    "parentLiftId": null,
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

```bash
# Bench Press (Competition)
curl -X POST "http://localhost:8080/lifts" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bench Press",
    "slug": "bench-press",
    "isCompetitionLift": true
  }'

# Deadlift (Competition)
curl -X POST "http://localhost:8080/lifts" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Deadlift",
    "slug": "deadlift",
    "isCompetitionLift": true
  }'

# Deadlift Up to Knees (Variation)
curl -X POST "http://localhost:8080/lifts" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Deadlift Up to Knees",
    "slug": "deadlift-up-to-knees",
    "isCompetitionLift": false,
    "parentLiftId": "33333333-3333-3333-3333-333333333333"
  }'

# Deadlift from Boxes (Variation)
curl -X POST "http://localhost:8080/lifts" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Deadlift from Boxes",
    "slug": "deadlift-from-boxes",
    "isCompetitionLift": false,
    "parentLiftId": "33333333-3333-3333-3333-333333333333"
  }'

# Incline Bench Press (Variation)
curl -X POST "http://localhost:8080/lifts" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Incline Bench Press",
    "slug": "incline-bench-press",
    "isCompetitionLift": false,
    "parentLiftId": "22222222-2222-2222-2222-222222222222"
  }'

# Good Morning (Accessory)
curl -X POST "http://localhost:8080/lifts" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Good Morning",
    "slug": "good-morning",
    "isCompetitionLift": false
  }'
```

**Lift IDs for subsequent steps** (example UUIDs):
| Lift | ID |
|------|-----|
| Squat | `11111111-1111-1111-1111-111111111111` |
| Bench Press | `22222222-2222-2222-2222-222222222222` |
| Deadlift | `33333333-3333-3333-3333-333333333333` |
| Deadlift Up to Knees | `44444444-4444-4444-4444-444444444444` |
| Deadlift from Boxes | `55555555-5555-5555-5555-555555555555` |
| Incline Bench Press | `66666666-6666-6666-6666-666666666666` |
| Good Morning | `77777777-7777-7777-7777-777777777777` |

---

### Step 2: Create Prescriptions (High-Volume Sets)

Sheiko uses RAMP prescriptions with warm-up sets built into the scheme. Each prescription includes both warm-up and working sets.

**Day 1 Squat Prescription (Prep Week 1):**

50%x5x1, 60%x5x1, 70%x4x4 = 26 total reps

```bash
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "11111111-1111-1111-1111-111111111111",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "TRAINING_MAX",
      "percentage": 100.0,
      "roundTo": 2.5
    },
    "setScheme": {
      "type": "RAMP",
      "steps": [
        {"percentage": 50, "reps": 5},
        {"percentage": 60, "reps": 5},
        {"percentage": 70, "reps": 4},
        {"percentage": 70, "reps": 4},
        {"percentage": 70, "reps": 4},
        {"percentage": 70, "reps": 4}
      ],
      "workSetThreshold": 70
    },
    "order": 1,
    "notes": "Focus on technique. Weight should feel moderate.",
    "restSeconds": 180
  }'
```

**Response:**
```json
{
  "data": {
    "id": "rx-squat-d1-w1-uuid",
    "liftId": "11111111-1111-1111-1111-111111111111",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "TRAINING_MAX",
      "percentage": 100.0,
      "roundTo": 2.5
    },
    "setScheme": {
      "type": "RAMP",
      "steps": [
        {"percentage": 50, "reps": 5},
        {"percentage": 60, "reps": 5},
        {"percentage": 70, "reps": 4},
        {"percentage": 70, "reps": 4},
        {"percentage": 70, "reps": 4},
        {"percentage": 70, "reps": 4}
      ],
      "workSetThreshold": 70
    },
    "order": 1,
    "notes": "Focus on technique. Weight should feel moderate.",
    "restSeconds": 180,
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

**Day 1 Bench Press Prescription (Prep Week 1):**

50%x5x1, 60%x4x1, 70%x3x1, 75%x3x4 = 24 total reps

```bash
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "22222222-2222-2222-2222-222222222222",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "TRAINING_MAX",
      "percentage": 100.0,
      "roundTo": 2.5
    },
    "setScheme": {
      "type": "RAMP",
      "steps": [
        {"percentage": 50, "reps": 5},
        {"percentage": 60, "reps": 4},
        {"percentage": 70, "reps": 3},
        {"percentage": 75, "reps": 3},
        {"percentage": 75, "reps": 3},
        {"percentage": 75, "reps": 3},
        {"percentage": 75, "reps": 3}
      ],
      "workSetThreshold": 70
    },
    "order": 2,
    "notes": "Controlled descent, pause on chest, drive up.",
    "restSeconds": 180
  }'
```

**Day 1 Good Morning Prescription:**

75%x4x5 = 20 total reps

```bash
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "77777777-7777-7777-7777-777777777777",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "TRAINING_MAX",
      "percentage": 75.0,
      "roundTo": 2.5
    },
    "setScheme": {
      "type": "FIXED",
      "sets": 5,
      "reps": 4,
      "isAmrap": false
    },
    "order": 3,
    "notes": "Maintain flat back, hinge at hips.",
    "restSeconds": 120
  }'
```

**Day 3 Deadlift Up to Knees Prescription (Prep Week 1):**

50%x3x1, 60%x3x1, 70%x3x1, 75%x2x4 = 17 total reps

```bash
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "44444444-4444-4444-4444-444444444444",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "TRAINING_MAX",
      "percentage": 100.0,
      "roundTo": 2.5
    },
    "setScheme": {
      "type": "RAMP",
      "steps": [
        {"percentage": 50, "reps": 3},
        {"percentage": 60, "reps": 3},
        {"percentage": 70, "reps": 3},
        {"percentage": 75, "reps": 2},
        {"percentage": 75, "reps": 2},
        {"percentage": 75, "reps": 2},
        {"percentage": 75, "reps": 2}
      ],
      "workSetThreshold": 70
    },
    "order": 1,
    "notes": "Stop at knee height, control descent.",
    "restSeconds": 180
  }'
```

**Day 3 Incline Bench Press Prescription:**

85%x3x5 = 15 total reps (of bench TM)

```bash
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "66666666-6666-6666-6666-666666666666",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "TRAINING_MAX",
      "percentage": 85.0,
      "roundTo": 2.5
    },
    "setScheme": {
      "type": "FIXED",
      "sets": 5,
      "reps": 3,
      "isAmrap": false
    },
    "order": 2,
    "notes": "Upper chest focus, full range of motion.",
    "restSeconds": 120
  }'
```

**Day 3 Deadlift from Boxes Prescription (Prep Week 1):**

55%x3x1, 65%x3x1, 75%x3x1, 85%x2x3 = 15 total reps

```bash
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "55555555-5555-5555-5555-555555555555",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "TRAINING_MAX",
      "percentage": 100.0,
      "roundTo": 2.5
    },
    "setScheme": {
      "type": "RAMP",
      "steps": [
        {"percentage": 55, "reps": 3},
        {"percentage": 65, "reps": 3},
        {"percentage": 75, "reps": 3},
        {"percentage": 85, "reps": 2},
        {"percentage": 85, "reps": 2},
        {"percentage": 85, "reps": 2}
      ],
      "workSetThreshold": 75
    },
    "order": 3,
    "notes": "Elevated start position, focus on lockout.",
    "restSeconds": 180
  }'
```

**Day 5 Squat Prescription (Prep Week 1):**

50%x5x1, 60%x4x1, 70%x3x1, 75%x3x4 = 24 total reps

```bash
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "11111111-1111-1111-1111-111111111111",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "TRAINING_MAX",
      "percentage": 100.0,
      "roundTo": 2.5
    },
    "setScheme": {
      "type": "RAMP",
      "steps": [
        {"percentage": 50, "reps": 5},
        {"percentage": 60, "reps": 4},
        {"percentage": 70, "reps": 3},
        {"percentage": 75, "reps": 3},
        {"percentage": 75, "reps": 3},
        {"percentage": 75, "reps": 3},
        {"percentage": 75, "reps": 3}
      ],
      "workSetThreshold": 70
    },
    "order": 1,
    "notes": "Slightly heavier than Day 1. Maintain technique.",
    "restSeconds": 180
  }'
```

**Day 5 Bench Press Prescription (Prep Week 1):**

50%x5x1, 60%x4x1, 70%x3x1, 80%x2x4 = 20 total reps

```bash
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "22222222-2222-2222-2222-222222222222",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "TRAINING_MAX",
      "percentage": 100.0,
      "roundTo": 2.5
    },
    "setScheme": {
      "type": "RAMP",
      "steps": [
        {"percentage": 50, "reps": 5},
        {"percentage": 60, "reps": 4},
        {"percentage": 70, "reps": 3},
        {"percentage": 80, "reps": 2},
        {"percentage": 80, "reps": 2},
        {"percentage": 80, "reps": 2},
        {"percentage": 80, "reps": 2}
      ],
      "workSetThreshold": 70
    },
    "order": 2,
    "notes": "Heavier doubles. Maintain bar path.",
    "restSeconds": 180
  }'
```

**Prescription IDs for subsequent steps** (example UUIDs):
| Prescription | ID |
|--------------|-----|
| Squat D1 W1 | `aaaa1111-1111-1111-1111-111111111111` |
| Bench D1 W1 | `aaaa2222-2222-2222-2222-222222222222` |
| Good Morning D1 | `aaaa3333-3333-3333-3333-333333333333` |
| DL Up to Knees D3 W1 | `aaaa4444-4444-4444-4444-444444444444` |
| Incline Bench D3 | `aaaa5555-5555-5555-5555-555555555555` |
| DL from Boxes D3 W1 | `aaaa6666-6666-6666-6666-666666666666` |
| Squat D5 W1 | `aaaa7777-7777-7777-7777-777777777777` |
| Bench D5 W1 | `aaaa8888-8888-8888-8888-888888888888` |

---

### Step 3: Create Days

Create the three training days for the Prep phase Week 1.

**Day 1 (Monday):**
```bash
curl -X POST "http://localhost:8080/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Day 1 - Squat/Bench",
    "slug": "sheiko-prep-w1-d1",
    "metadata": {
      "focus": "squat-bench-volume",
      "program": "sheiko-beginner",
      "phase": "prep",
      "week": 1
    }
  }'
```

**Response:**
```json
{
  "data": {
    "id": "day-w1d1-uuid",
    "name": "Day 1 - Squat/Bench",
    "slug": "sheiko-prep-w1-d1",
    "metadata": {
      "focus": "squat-bench-volume",
      "program": "sheiko-beginner",
      "phase": "prep",
      "week": 1
    },
    "programId": null,
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

```bash
# Day 3 (Wednesday)
curl -X POST "http://localhost:8080/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Day 3 - Deadlift Variations",
    "slug": "sheiko-prep-w1-d3",
    "metadata": {
      "focus": "deadlift-variations-bench",
      "program": "sheiko-beginner",
      "phase": "prep",
      "week": 1
    }
  }'

# Day 5 (Friday)
curl -X POST "http://localhost:8080/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Day 5 - Squat/Bench",
    "slug": "sheiko-prep-w1-d5",
    "metadata": {
      "focus": "squat-bench-heavier",
      "program": "sheiko-beginner",
      "phase": "prep",
      "week": 1
    }
  }'
```

**Day IDs:**
| Day | ID |
|-----|-----|
| Day 1 (W1) | `bbbb1111-1111-1111-1111-111111111111` |
| Day 3 (W1) | `bbbb2222-2222-2222-2222-222222222222` |
| Day 5 (W1) | `bbbb3333-3333-3333-3333-333333333333` |

**Add prescriptions to Day 1:**
```bash
curl -X POST "http://localhost:8080/days/bbbb1111-1111-1111-1111-111111111111/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"prescriptionId": "aaaa1111-1111-1111-1111-111111111111", "order": 1}'

curl -X POST "http://localhost:8080/days/bbbb1111-1111-1111-1111-111111111111/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"prescriptionId": "aaaa2222-2222-2222-2222-222222222222", "order": 2}'

curl -X POST "http://localhost:8080/days/bbbb1111-1111-1111-1111-111111111111/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"prescriptionId": "aaaa3333-3333-3333-3333-333333333333", "order": 3}'
```

**Add prescriptions to Day 3:**
```bash
curl -X POST "http://localhost:8080/days/bbbb2222-2222-2222-2222-222222222222/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"prescriptionId": "aaaa4444-4444-4444-4444-444444444444", "order": 1}'

curl -X POST "http://localhost:8080/days/bbbb2222-2222-2222-2222-222222222222/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"prescriptionId": "aaaa5555-5555-5555-5555-555555555555", "order": 2}'

curl -X POST "http://localhost:8080/days/bbbb2222-2222-2222-2222-222222222222/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"prescriptionId": "aaaa6666-6666-6666-6666-666666666666", "order": 3}'
```

**Add prescriptions to Day 5:**
```bash
curl -X POST "http://localhost:8080/days/bbbb3333-3333-3333-3333-333333333333/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"prescriptionId": "aaaa7777-7777-7777-7777-777777777777", "order": 1}'

curl -X POST "http://localhost:8080/days/bbbb3333-3333-3333-3333-333333333333/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"prescriptionId": "aaaa8888-8888-8888-8888-888888888888", "order": 2}'
```

---

### Step 4: Create Cycle

Sheiko Prep phase uses a 4-week cycle. For simplicity, we configure Week 1 here (full implementation would include all 4 weeks with varying prescriptions).

```bash
curl -X POST "http://localhost:8080/cycles" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Sheiko Beginner Prep Phase 4-Week Cycle",
    "lengthWeeks": 4
  }'
```

**Response:**
```json
{
  "data": {
    "id": "cccc1111-1111-1111-1111-111111111111",
    "name": "Sheiko Beginner Prep Phase 4-Week Cycle",
    "lengthWeeks": 4,
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

---

### Step 5: Create Weeks

Create Week 1 of the Prep phase. A complete implementation would include all 4 weeks.

```bash
curl -X POST "http://localhost:8080/weeks" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "cycleId": "cccc1111-1111-1111-1111-111111111111",
    "weekNumber": 1,
    "name": "Prep Week 1"
  }'
```

**Response:**
```json
{
  "data": {
    "id": "dddd1111-1111-1111-1111-111111111111",
    "cycleId": "cccc1111-1111-1111-1111-111111111111",
    "weekNumber": 1,
    "name": "Prep Week 1",
    "days": [],
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

**Add days to Week 1 (M/W/F pattern):**
```bash
# Position 0: Day 1 (Monday)
curl -X POST "http://localhost:8080/weeks/dddd1111-1111-1111-1111-111111111111/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "bbbb1111-1111-1111-1111-111111111111", "position": 0}'

# Position 1: Day 3 (Wednesday)
curl -X POST "http://localhost:8080/weeks/dddd1111-1111-1111-1111-111111111111/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "bbbb2222-2222-2222-2222-222222222222", "position": 1}'

# Position 2: Day 5 (Friday)
curl -X POST "http://localhost:8080/weeks/dddd1111-1111-1111-1111-111111111111/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "bbbb3333-3333-3333-3333-333333333333", "position": 2}'
```

---

### Step 6: Create Program (No Progressions)

Create the Sheiko Beginner program **without linking any progressions**. This is the key difference from auto-progression programs.

```bash
curl -X POST "http://localhost:8080/programs" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Sheiko Beginner",
    "slug": "sheiko-beginner",
    "description": "Boris Sheiko'\''s beginner powerlifting program. High volume, moderate intensity, manual progression. Run Prep cycles repeatedly, updating maxes manually when technique improves at current weights.",
    "cycleId": "cccc1111-1111-1111-1111-111111111111",
    "defaultRounding": 2.5
  }'
```

**Response:**
```json
{
  "data": {
    "id": "program-sheiko-uuid",
    "name": "Sheiko Beginner",
    "slug": "sheiko-beginner",
    "description": "Boris Sheiko's beginner powerlifting program. High volume, moderate intensity, manual progression. Run Prep cycles repeatedly, updating maxes manually when technique improves at current weights.",
    "cycleId": "cccc1111-1111-1111-1111-111111111111",
    "weeklyLookupId": null,
    "dailyLookupId": null,
    "defaultRounding": 2.5,
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

**Program ID:** `eeee1111-1111-1111-1111-111111111111`

**Note: No progressions are linked to this program.** The user will manually update their training maxes when ready.

---

## User Setup

Before a user can generate workouts, they need training maxes for all lifts.

### Create User Training Maxes

```bash
# User ID for examples
USER_ID="user-12345678-1234-1234-1234-123456789012"

# Squat Training Max: 200kg
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "11111111-1111-1111-1111-111111111111",
    "type": "TRAINING_MAX",
    "value": 200.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'

# Bench Press Training Max: 140kg
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "22222222-2222-2222-2222-222222222222",
    "type": "TRAINING_MAX",
    "value": 140.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'

# Deadlift Training Max: 220kg
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "33333333-3333-3333-3333-333333333333",
    "type": "TRAINING_MAX",
    "value": 220.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'

# Deadlift Up to Knees (uses deadlift max)
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "44444444-4444-4444-4444-444444444444",
    "type": "TRAINING_MAX",
    "value": 220.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'

# Deadlift from Boxes (uses deadlift max)
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "55555555-5555-5555-5555-555555555555",
    "type": "TRAINING_MAX",
    "value": 220.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'

# Incline Bench Press (uses bench max)
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "66666666-6666-6666-6666-666666666666",
    "type": "TRAINING_MAX",
    "value": 140.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'

# Good Morning Training Max: 100kg
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "77777777-7777-7777-7777-777777777777",
    "type": "TRAINING_MAX",
    "value": 100.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'
```

### Enroll User in Program

```bash
curl -X POST "http://localhost:8080/users/${USER_ID}/program" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "programId": "eeee1111-1111-1111-1111-111111111111"
  }'
```

**Response:**
```json
{
  "data": {
    "id": "enrollment-uuid",
    "userId": "user-12345678-1234-1234-1234-123456789012",
    "program": {
      "id": "eeee1111-1111-1111-1111-111111111111",
      "name": "Sheiko Beginner",
      "slug": "sheiko-beginner",
      "description": "Boris Sheiko's beginner powerlifting program. High volume, moderate intensity, manual progression.",
      "cycleLengthWeeks": 4
    },
    "state": {
      "currentWeek": 1,
      "currentCycleIteration": 1,
      "currentDayIndex": 0
    },
    "enrolledAt": "2024-01-15T10:00:00Z",
    "updatedAt": "2024-01-15T10:00:00Z"
  }
}
```

---

## Example Workout Output

### Prep Phase Week 1 Day 1

Training maxes: Squat 200kg, Bench 140kg, Good Morning 100kg

```bash
curl -X GET "http://localhost:8080/users/${USER_ID}/workout" \
  -H "X-User-ID: ${USER_ID}"
```

**Response:**
```json
{
  "data": {
    "userId": "user-12345678-1234-1234-1234-123456789012",
    "programId": "eeee1111-1111-1111-1111-111111111111",
    "cycleIteration": 1,
    "weekNumber": 1,
    "weekName": "Prep Week 1",
    "daySlug": "sheiko-prep-w1-d1",
    "date": "2024-01-15",
    "exercises": [
      {
        "prescriptionId": "aaaa1111-1111-1111-1111-111111111111",
        "lift": {
          "id": "11111111-1111-1111-1111-111111111111",
          "name": "Squat",
          "slug": "squat"
        },
        "sets": [
          {"setNumber": 1, "weight": 100.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 2, "weight": 120.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 3, "weight": 140.0, "targetReps": 4, "isWorkSet": true},
          {"setNumber": 4, "weight": 140.0, "targetReps": 4, "isWorkSet": true},
          {"setNumber": 5, "weight": 140.0, "targetReps": 4, "isWorkSet": true},
          {"setNumber": 6, "weight": 140.0, "targetReps": 4, "isWorkSet": true}
        ],
        "notes": "Focus on technique. Weight should feel moderate.",
        "restSeconds": 180
      },
      {
        "prescriptionId": "aaaa2222-2222-2222-2222-222222222222",
        "lift": {
          "id": "22222222-2222-2222-2222-222222222222",
          "name": "Bench Press",
          "slug": "bench-press"
        },
        "sets": [
          {"setNumber": 1, "weight": 70.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 2, "weight": 85.0, "targetReps": 4, "isWorkSet": false},
          {"setNumber": 3, "weight": 97.5, "targetReps": 3, "isWorkSet": true},
          {"setNumber": 4, "weight": 105.0, "targetReps": 3, "isWorkSet": true},
          {"setNumber": 5, "weight": 105.0, "targetReps": 3, "isWorkSet": true},
          {"setNumber": 6, "weight": 105.0, "targetReps": 3, "isWorkSet": true},
          {"setNumber": 7, "weight": 105.0, "targetReps": 3, "isWorkSet": true}
        ],
        "notes": "Controlled descent, pause on chest, drive up.",
        "restSeconds": 180
      },
      {
        "prescriptionId": "aaaa3333-3333-3333-3333-333333333333",
        "lift": {
          "id": "77777777-7777-7777-7777-777777777777",
          "name": "Good Morning",
          "slug": "good-morning"
        },
        "sets": [
          {"setNumber": 1, "weight": 75.0, "targetReps": 4, "isWorkSet": true},
          {"setNumber": 2, "weight": 75.0, "targetReps": 4, "isWorkSet": true},
          {"setNumber": 3, "weight": 75.0, "targetReps": 4, "isWorkSet": true},
          {"setNumber": 4, "weight": 75.0, "targetReps": 4, "isWorkSet": true},
          {"setNumber": 5, "weight": 75.0, "targetReps": 4, "isWorkSet": true}
        ],
        "notes": "Maintain flat back, hinge at hips.",
        "restSeconds": 120
      }
    ]
  }
}
```

**Day 1 Volume Summary:**
| Exercise | Work Sets | Total Reps | Top Weight |
|----------|-----------|------------|------------|
| Squat | 4 | 26 | 140kg (70%) |
| Bench Press | 5 | 24 | 105kg (75%) |
| Good Morning | 5 | 20 | 75kg (75%) |

---

### Prep Phase Week 1 Day 3

After advancing to day index 1:

```bash
curl -X POST "http://localhost:8080/users/${USER_ID}/program-state/advance" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{"advanceType": "day"}'
```

Then generate the workout:

```bash
curl -X GET "http://localhost:8080/users/${USER_ID}/workout" \
  -H "X-User-ID: ${USER_ID}"
```

**Response:**
```json
{
  "data": {
    "userId": "user-12345678-1234-1234-1234-123456789012",
    "programId": "eeee1111-1111-1111-1111-111111111111",
    "cycleIteration": 1,
    "weekNumber": 1,
    "weekName": "Prep Week 1",
    "daySlug": "sheiko-prep-w1-d3",
    "date": "2024-01-17",
    "exercises": [
      {
        "prescriptionId": "aaaa4444-4444-4444-4444-444444444444",
        "lift": {
          "id": "44444444-4444-4444-4444-444444444444",
          "name": "Deadlift Up to Knees",
          "slug": "deadlift-up-to-knees"
        },
        "sets": [
          {"setNumber": 1, "weight": 110.0, "targetReps": 3, "isWorkSet": false},
          {"setNumber": 2, "weight": 132.5, "targetReps": 3, "isWorkSet": false},
          {"setNumber": 3, "weight": 155.0, "targetReps": 3, "isWorkSet": true},
          {"setNumber": 4, "weight": 165.0, "targetReps": 2, "isWorkSet": true},
          {"setNumber": 5, "weight": 165.0, "targetReps": 2, "isWorkSet": true},
          {"setNumber": 6, "weight": 165.0, "targetReps": 2, "isWorkSet": true},
          {"setNumber": 7, "weight": 165.0, "targetReps": 2, "isWorkSet": true}
        ],
        "notes": "Stop at knee height, control descent.",
        "restSeconds": 180
      },
      {
        "prescriptionId": "aaaa5555-5555-5555-5555-555555555555",
        "lift": {
          "id": "66666666-6666-6666-6666-666666666666",
          "name": "Incline Bench Press",
          "slug": "incline-bench-press"
        },
        "sets": [
          {"setNumber": 1, "weight": 120.0, "targetReps": 3, "isWorkSet": true},
          {"setNumber": 2, "weight": 120.0, "targetReps": 3, "isWorkSet": true},
          {"setNumber": 3, "weight": 120.0, "targetReps": 3, "isWorkSet": true},
          {"setNumber": 4, "weight": 120.0, "targetReps": 3, "isWorkSet": true},
          {"setNumber": 5, "weight": 120.0, "targetReps": 3, "isWorkSet": true}
        ],
        "notes": "Upper chest focus, full range of motion.",
        "restSeconds": 120
      },
      {
        "prescriptionId": "aaaa6666-6666-6666-6666-666666666666",
        "lift": {
          "id": "55555555-5555-5555-5555-555555555555",
          "name": "Deadlift from Boxes",
          "slug": "deadlift-from-boxes"
        },
        "sets": [
          {"setNumber": 1, "weight": 122.5, "targetReps": 3, "isWorkSet": false},
          {"setNumber": 2, "weight": 142.5, "targetReps": 3, "isWorkSet": false},
          {"setNumber": 3, "weight": 165.0, "targetReps": 3, "isWorkSet": true},
          {"setNumber": 4, "weight": 187.5, "targetReps": 2, "isWorkSet": true},
          {"setNumber": 5, "weight": 187.5, "targetReps": 2, "isWorkSet": true},
          {"setNumber": 6, "weight": 187.5, "targetReps": 2, "isWorkSet": true}
        ],
        "notes": "Elevated start position, focus on lockout.",
        "restSeconds": 180
      }
    ]
  }
}
```

**Day 3 Volume Summary:**
| Exercise | Work Sets | Total Reps | Top Weight |
|----------|-----------|------------|------------|
| Deadlift Up to Knees | 5 | 17 | 165kg (75%) |
| Incline Bench Press | 5 | 15 | 120kg (85% of bench) |
| Deadlift from Boxes | 4 | 15 | 187.5kg (85%) |

---

### Prep Phase Week 1 Day 5

After advancing to day index 2:

```bash
curl -X POST "http://localhost:8080/users/${USER_ID}/program-state/advance" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{"advanceType": "day"}'
```

Then generate the workout:

```bash
curl -X GET "http://localhost:8080/users/${USER_ID}/workout" \
  -H "X-User-ID: ${USER_ID}"
```

**Response:**
```json
{
  "data": {
    "userId": "user-12345678-1234-1234-1234-123456789012",
    "programId": "eeee1111-1111-1111-1111-111111111111",
    "cycleIteration": 1,
    "weekNumber": 1,
    "weekName": "Prep Week 1",
    "daySlug": "sheiko-prep-w1-d5",
    "date": "2024-01-19",
    "exercises": [
      {
        "prescriptionId": "aaaa7777-7777-7777-7777-777777777777",
        "lift": {
          "id": "11111111-1111-1111-1111-111111111111",
          "name": "Squat",
          "slug": "squat"
        },
        "sets": [
          {"setNumber": 1, "weight": 100.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 2, "weight": 120.0, "targetReps": 4, "isWorkSet": false},
          {"setNumber": 3, "weight": 140.0, "targetReps": 3, "isWorkSet": true},
          {"setNumber": 4, "weight": 150.0, "targetReps": 3, "isWorkSet": true},
          {"setNumber": 5, "weight": 150.0, "targetReps": 3, "isWorkSet": true},
          {"setNumber": 6, "weight": 150.0, "targetReps": 3, "isWorkSet": true},
          {"setNumber": 7, "weight": 150.0, "targetReps": 3, "isWorkSet": true}
        ],
        "notes": "Slightly heavier than Day 1. Maintain technique.",
        "restSeconds": 180
      },
      {
        "prescriptionId": "aaaa8888-8888-8888-8888-888888888888",
        "lift": {
          "id": "22222222-2222-2222-2222-222222222222",
          "name": "Bench Press",
          "slug": "bench-press"
        },
        "sets": [
          {"setNumber": 1, "weight": 70.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 2, "weight": 85.0, "targetReps": 4, "isWorkSet": false},
          {"setNumber": 3, "weight": 97.5, "targetReps": 3, "isWorkSet": true},
          {"setNumber": 4, "weight": 112.5, "targetReps": 2, "isWorkSet": true},
          {"setNumber": 5, "weight": 112.5, "targetReps": 2, "isWorkSet": true},
          {"setNumber": 6, "weight": 112.5, "targetReps": 2, "isWorkSet": true},
          {"setNumber": 7, "weight": 112.5, "targetReps": 2, "isWorkSet": true}
        ],
        "notes": "Heavier doubles. Maintain bar path.",
        "restSeconds": 180
      }
    ]
  }
}
```

**Day 5 Volume Summary:**
| Exercise | Work Sets | Total Reps | Top Weight |
|----------|-----------|------------|------------|
| Squat | 5 | 24 | 150kg (75%) |
| Bench Press | 5 | 20 | 112.5kg (80%) |

---

## Manual Progression Approach

Sheiko programs do not use automatic progression. The coach/athlete decides when to increase training maxes based on technique quality and perceived difficulty.

### When to Update Maxes

Update training maxes when:

1. **Weights feel consistently easy** - If prescribed weights feel light and technique is solid across multiple weeks
2. **After completing a full Prep cycle** - Evaluate after finishing all 4 weeks
3. **Technical improvement** - When previously challenging weights now move with better form
4. **Competition results** - After a successful meet, update based on actual lifts

**Do NOT update maxes when:**
- Form is breaking down
- Fatigue is accumulating
- Within 4 weeks of competition
- Recovering from illness or injury

### How to Update Maxes via API

To manually increase a training max, create a new lift max entry:

**Increase Squat Max from 200kg to 210kg:**
```bash
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "11111111-1111-1111-1111-111111111111",
    "type": "TRAINING_MAX",
    "value": 210.0,
    "effectiveDate": "2024-02-12T00:00:00Z"
  }'
```

**Response:**
```json
{
  "data": {
    "id": "liftmax-new-uuid",
    "userId": "user-12345678-1234-1234-1234-123456789012",
    "liftId": "11111111-1111-1111-1111-111111111111",
    "type": "TRAINING_MAX",
    "value": 210.0,
    "effectiveDate": "2024-02-12T00:00:00Z",
    "createdAt": "2024-02-12T00:00:00Z",
    "updatedAt": "2024-02-12T00:00:00Z"
  }
}
```

**Increase Bench Max from 140kg to 145kg:**
```bash
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "22222222-2222-2222-2222-222222222222",
    "type": "TRAINING_MAX",
    "value": 145.0,
    "effectiveDate": "2024-02-12T00:00:00Z"
  }'
```

**Increase Deadlift Max from 220kg to 230kg:**
```bash
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "33333333-3333-3333-3333-333333333333",
    "type": "TRAINING_MAX",
    "value": 230.0,
    "effectiveDate": "2024-02-12T00:00:00Z"
  }'
```

**Remember to also update variation maxes:**
```bash
# Update Deadlift Up to Knees
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "44444444-4444-4444-4444-444444444444",
    "type": "TRAINING_MAX",
    "value": 230.0,
    "effectiveDate": "2024-02-12T00:00:00Z"
  }'

# Update Deadlift from Boxes
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "55555555-5555-5555-5555-555555555555",
    "type": "TRAINING_MAX",
    "value": 230.0,
    "effectiveDate": "2024-02-12T00:00:00Z"
  }'

# Update Incline Bench
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "66666666-6666-6666-6666-666666666666",
    "type": "TRAINING_MAX",
    "value": 145.0,
    "effectiveDate": "2024-02-12T00:00:00Z"
  }'
```

### Verify Updated Maxes

After updating, verify the new maxes are in effect:

```bash
curl -X GET "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}"
```

**Response:**
```json
{
  "data": [
    {
      "liftId": "11111111-1111-1111-1111-111111111111",
      "liftName": "Squat",
      "type": "TRAINING_MAX",
      "currentValue": 210.0,
      "effectiveDate": "2024-02-12T00:00:00Z"
    },
    {
      "liftId": "22222222-2222-2222-2222-222222222222",
      "liftName": "Bench Press",
      "type": "TRAINING_MAX",
      "currentValue": 145.0,
      "effectiveDate": "2024-02-12T00:00:00Z"
    },
    {
      "liftId": "33333333-3333-3333-3333-333333333333",
      "liftName": "Deadlift",
      "type": "TRAINING_MAX",
      "currentValue": 230.0,
      "effectiveDate": "2024-02-12T00:00:00Z"
    }
  ]
}
```

### Workout After Max Update

The next workout will automatically use the new maxes:

**Day 1 Squat with new 210kg max:**
| Set | Old Weight (200kg TM) | New Weight (210kg TM) |
|-----|----------------------|----------------------|
| 1 | 100kg (50%) | 105kg (50%) |
| 2 | 120kg (60%) | 127.5kg (60%) |
| 3-6 | 140kg (70%) | 147.5kg (70%) |

### Contrast with Auto-Progression Programs

| Program | Progression Type | How It Works |
|---------|-----------------|--------------|
| **Starting Strength** | AFTER_SESSION | +5/10 lbs every workout automatically |
| **Wendler 5/3/1** | AFTER_CYCLE | +5/10 lbs after 4 weeks automatically |
| **Bill Starr 5x5** | AFTER_WEEK | +5 lbs weekly automatically |
| **Sheiko Beginner** | **MANUAL** | User calls API to update maxes when ready |

**Why Sheiko uses manual progression:**
1. Emphasis on technique over load
2. Progress measured by volume tolerance, not weight
3. Coach/athlete can assess readiness
4. Prevents premature weight increases
5. Accommodates individual recovery rates

---

## Complete Configuration Summary

### Entity Relationship Diagram

```
Program: Sheiko Beginner
├── Cycle: 4-week Prep Phase
│   ├── Week 1
│   │   ├── Position 0: Day 1 (Squat/Bench/Good Morning)
│   │   ├── Position 1: Day 3 (DL Variations/Incline Bench)
│   │   └── Position 2: Day 5 (Squat/Bench - heavier)
│   ├── Week 2
│   │   └── [Similar structure with different prescriptions]
│   ├── Week 3
│   │   └── [Peak intensity within prep phase]
│   └── Week 4
│       └── [Transition week]
│
├── Days (Week 1 shown)
│   ├── Day 1 - Squat/Bench
│   │   ├── Prescription: Squat 50%x5, 60%x5, 70%x4x4
│   │   ├── Prescription: Bench 50%x5, 60%x4, 70%x3, 75%x3x4
│   │   └── Prescription: Good Morning 75%x4x5
│   │
│   ├── Day 3 - Deadlift Variations
│   │   ├── Prescription: DL Up to Knees 50%x3, 60%x3, 70%x3, 75%x2x4
│   │   ├── Prescription: Incline Bench 85%x3x5
│   │   └── Prescription: DL from Boxes 55%x3, 65%x3, 75%x3, 85%x2x3
│   │
│   └── Day 5 - Squat/Bench
│       ├── Prescription: Squat 50%x5, 60%x4, 70%x3, 75%x3x4
│       └── Prescription: Bench 50%x5, 60%x4, 70%x3, 80%x2x4
│
└── Program Progressions: NONE
    └── Manual max updates via POST /users/{id}/lift-maxes
```

### JSON Configuration Export

```json
{
  "program": {
    "name": "Sheiko Beginner",
    "slug": "sheiko-beginner",
    "description": "Boris Sheiko's beginner powerlifting program. High volume, moderate intensity, manual progression.",
    "defaultRounding": 2.5
  },
  "cycle": {
    "name": "Sheiko Beginner Prep Phase 4-Week Cycle",
    "lengthWeeks": 4
  },
  "weeks": [
    {
      "weekNumber": 1,
      "name": "Prep Week 1",
      "days": [
        {"position": 0, "daySlug": "sheiko-prep-w1-d1"},
        {"position": 1, "daySlug": "sheiko-prep-w1-d3"},
        {"position": 2, "daySlug": "sheiko-prep-w1-d5"}
      ]
    }
  ],
  "days": [
    {
      "name": "Day 1 - Squat/Bench",
      "slug": "sheiko-prep-w1-d1",
      "prescriptions": [
        {
          "liftSlug": "squat",
          "setScheme": {
            "type": "RAMP",
            "steps": [
              {"percentage": 50, "reps": 5},
              {"percentage": 60, "reps": 5},
              {"percentage": 70, "reps": 4},
              {"percentage": 70, "reps": 4},
              {"percentage": 70, "reps": 4},
              {"percentage": 70, "reps": 4}
            ]
          },
          "order": 1
        },
        {
          "liftSlug": "bench-press",
          "setScheme": {
            "type": "RAMP",
            "steps": [
              {"percentage": 50, "reps": 5},
              {"percentage": 60, "reps": 4},
              {"percentage": 70, "reps": 3},
              {"percentage": 75, "reps": 3},
              {"percentage": 75, "reps": 3},
              {"percentage": 75, "reps": 3},
              {"percentage": 75, "reps": 3}
            ]
          },
          "order": 2
        },
        {
          "liftSlug": "good-morning",
          "setScheme": {"type": "FIXED", "sets": 5, "reps": 4, "percentage": 75},
          "order": 3
        }
      ]
    },
    {
      "name": "Day 3 - Deadlift Variations",
      "slug": "sheiko-prep-w1-d3",
      "prescriptions": [
        {
          "liftSlug": "deadlift-up-to-knees",
          "setScheme": {
            "type": "RAMP",
            "steps": [
              {"percentage": 50, "reps": 3},
              {"percentage": 60, "reps": 3},
              {"percentage": 70, "reps": 3},
              {"percentage": 75, "reps": 2},
              {"percentage": 75, "reps": 2},
              {"percentage": 75, "reps": 2},
              {"percentage": 75, "reps": 2}
            ]
          },
          "order": 1
        },
        {
          "liftSlug": "incline-bench-press",
          "setScheme": {"type": "FIXED", "sets": 5, "reps": 3, "percentage": 85},
          "order": 2
        },
        {
          "liftSlug": "deadlift-from-boxes",
          "setScheme": {
            "type": "RAMP",
            "steps": [
              {"percentage": 55, "reps": 3},
              {"percentage": 65, "reps": 3},
              {"percentage": 75, "reps": 3},
              {"percentage": 85, "reps": 2},
              {"percentage": 85, "reps": 2},
              {"percentage": 85, "reps": 2}
            ]
          },
          "order": 3
        }
      ]
    },
    {
      "name": "Day 5 - Squat/Bench",
      "slug": "sheiko-prep-w1-d5",
      "prescriptions": [
        {
          "liftSlug": "squat",
          "setScheme": {
            "type": "RAMP",
            "steps": [
              {"percentage": 50, "reps": 5},
              {"percentage": 60, "reps": 4},
              {"percentage": 70, "reps": 3},
              {"percentage": 75, "reps": 3},
              {"percentage": 75, "reps": 3},
              {"percentage": 75, "reps": 3},
              {"percentage": 75, "reps": 3}
            ]
          },
          "order": 1
        },
        {
          "liftSlug": "bench-press",
          "setScheme": {
            "type": "RAMP",
            "steps": [
              {"percentage": 50, "reps": 5},
              {"percentage": 60, "reps": 4},
              {"percentage": 70, "reps": 3},
              {"percentage": 80, "reps": 2},
              {"percentage": 80, "reps": 2},
              {"percentage": 80, "reps": 2},
              {"percentage": 80, "reps": 2}
            ]
          },
          "order": 2
        }
      ]
    }
  ],
  "progressions": []
}
```

---

## See Also

- [API Reference](../api-reference.md) - Complete endpoint documentation
- [Workflows](../workflows.md) - Multi-step API workflows
- [Wendler 5/3/1 BBB Configuration](./wendler-531-bbb.md) - Cycle-end auto-progression example
- [Bill Starr 5x5 Configuration](./bill-starr-5x5.md) - Weekly auto-progression example
- [Sheiko Beginner Program Specification](../../programs/010-sheiko-beginner.md) - Full program algorithm

# Greg Nuckols High Frequency Program Configuration

This document provides the complete API configuration for the Greg Nuckols High Frequency program, demonstrating **daily undulating periodization** within a **3-week cycle structure** using multi-dimensional lookup tables.

## Table of Contents

- [Overview](#overview)
- [Program Characteristics](#program-characteristics)
- [Configuration Components](#configuration-components)
  - [Step 1: Create Lifts](#step-1-create-lifts)
  - [Step 2: Create Daily Lookup (Intensity by Day)](#step-2-create-daily-lookup-intensity-by-day)
  - [Step 3: Create Weekly Lookup (Volume Progression)](#step-3-create-weekly-lookup-volume-progression)
  - [Step 4: Create Prescriptions](#step-4-create-prescriptions)
  - [Step 5: Create Days](#step-5-create-days)
  - [Step 6: Create Cycle](#step-6-create-cycle)
  - [Step 7: Create Weeks](#step-7-create-weeks)
  - [Step 8: Create Progressions](#step-8-create-progressions)
  - [Step 9: Create Program](#step-9-create-program)
  - [Step 10: Link Progressions to Program](#step-10-link-progressions-to-program)
- [User Setup](#user-setup)
- [Example Workout Output](#example-workout-output)
  - [Week 1 Daily Variation](#week-1-daily-variation)
  - [Week 2 Volume Increase](#week-2-volume-increase)
  - [Week 3 Peak/Test Week](#week-3-peaktest-week)
- [Cycle Progression Behavior](#cycle-progression-behavior)
- [Complete Configuration Summary](#complete-configuration-summary)

---

## Overview

Greg Nuckols' High Frequency program is designed for intermediate to advanced lifters who benefit from practicing the competition lifts multiple times per week. The program combines **daily undulating periodization** (DUP) with a **3-week progression cycle**.

**Key Characteristics:**
- High frequency: Bench and Squat 5x per week, Deadlift 1x, OHP 1x
- Daily undulating periodization: Each day has different intensity (65-85%)
- 3-week cycle with progressive volume
- Week 3 includes AMAP (As Many As Possible) sets for autoregulated progression
- Back work uses 6RM-based loading with separate progression

**What This Configuration Demonstrates:**
- **Daily lookup tables** - Different intensity percentages for each day
- **Weekly lookup tables** - Volume progression across weeks
- **Multi-dimensional lookups** - Combining day AND week position
- **AMAP-based progression** - Training max updates based on Week 3 performance
- **Multiple max types** - 1RM for main lifts, 6RM for back work

---

## Program Characteristics

| Property | Value |
|----------|-------|
| Cycle Length | 3 weeks |
| Days per Week | 6 (Mon-Sat, Sun rest) |
| Progression Type | LINEAR (cycle-end, AMAP-guided) |
| Progression Trigger | AFTER_CYCLE |
| Intensity Range | 65-85% of 1RM |
| Volume Strategy | Progressive sets across 3 weeks |

### Daily Intensity Distribution

| Day | Intensity | Character | Focus |
|-----|-----------|-----------|-------|
| Monday | 75% | Moderate | Technical practice |
| Tuesday | 80% | Moderate-Heavy | Strength building |
| Wednesday | 70% | Light | Recovery/Volume |
| Thursday | 85% | Heavy | Intensity day |
| Friday | 85% (DL/OHP) | Heavy | Pull focus |
| Saturday | 65% | Light | High volume hypertrophy |

### 3-Week Volume Progression

The program uses fixed intensity with progressive volume (sets × reps):

**Monday (75% intensity):**
| Week | Sets | Reps |
|------|------|------|
| Week 1 | 4 | 3 |
| Week 2 | 5 | 3 |
| Week 3 | 5 | 4 |

**Tuesday (80% intensity):**
| Week | Sets | Reps |
|------|------|------|
| Week 1 | 3 | 2 |
| Week 2 | 4 | 2 |
| Week 3 | 4 | 3 |

**Wednesday (70% intensity):**
| Week | Sets | Reps |
|------|------|------|
| Week 1 | 4 | 4 |
| Week 2 | 5 | 4 |
| Week 3 | 5 | 5 |

**Thursday (85% intensity):**
| Week | Sets | Reps |
|------|------|------|
| Week 1 | 3 | 1 |
| Week 2 | 4 | 1 |
| Week 3 | 1 | AMAP |

**Friday - OHP/Deadlift (85% intensity):**
| Week | Sets | Reps |
|------|------|------|
| Week 1 | 3 | 1 |
| Week 2 | 4 | 1 |
| Week 3 | 1 | AMAP |

**Saturday (65% intensity):**
| Week | Sets | Reps |
|------|------|------|
| Week 1 | 5 | 5 |
| Week 2 | 6 | 5 |
| Week 3 | 6 | 6 |

---

## Configuration Components

All API calls below assume:
- Base URL: `http://localhost:8080`
- Admin authentication: `X-Admin: true`
- User ID header: `X-User-ID: {userId}`

### Step 1: Create Lifts

Create the main lifts used in the program.

```bash
# Bench Press
curl -X POST "http://localhost:8080/lifts" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bench Press",
    "slug": "bench-press",
    "isCompetitionLift": true
  }'
```

**Response:**
```json
{
  "data": {
    "id": "lift-bench-uuid",
    "name": "Bench Press",
    "slug": "bench-press",
    "isCompetitionLift": true,
    "parentLiftId": null,
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

```bash
# Squat
curl -X POST "http://localhost:8080/lifts" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Squat",
    "slug": "squat",
    "isCompetitionLift": true
  }'

# Deadlift
curl -X POST "http://localhost:8080/lifts" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Deadlift",
    "slug": "deadlift",
    "isCompetitionLift": true
  }'

# Overhead Press
curl -X POST "http://localhost:8080/lifts" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Overhead Press",
    "slug": "overhead-press",
    "isCompetitionLift": false
  }'

# T-Bar Row (uses 6RM)
curl -X POST "http://localhost:8080/lifts" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "T-Bar Row",
    "slug": "t-bar-row",
    "isCompetitionLift": false
  }'

# Pulldown (uses 6RM)
curl -X POST "http://localhost:8080/lifts" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Pulldown",
    "slug": "pulldown",
    "isCompetitionLift": false
  }'
```

**Lift IDs for subsequent steps** (example UUIDs):
| Lift | ID |
|------|-----|
| Bench Press | `11111111-1111-1111-1111-111111111111` |
| Squat | `22222222-2222-2222-2222-222222222222` |
| Deadlift | `33333333-3333-3333-3333-333333333333` |
| Overhead Press | `44444444-4444-4444-4444-444444444444` |
| T-Bar Row | `55555555-5555-5555-5555-555555555555` |
| Pulldown | `66666666-6666-6666-6666-666666666666` |

---

### Step 2: Create Daily Lookup (Intensity by Day)

The daily lookup maps day position to intensity percentage. This enables daily undulating periodization.

```bash
curl -X POST "http://localhost:8080/daily-lookups" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Greg Nuckols HF Daily Intensity",
    "entries": [
      {"dayPosition": 0, "percentage": 75.0},
      {"dayPosition": 1, "percentage": 80.0},
      {"dayPosition": 2, "percentage": 70.0},
      {"dayPosition": 3, "percentage": 85.0},
      {"dayPosition": 4, "percentage": 85.0},
      {"dayPosition": 5, "percentage": 65.0}
    ]
  }'
```

**Response:**
```json
{
  "data": {
    "id": "dlookup-hf-uuid",
    "name": "Greg Nuckols HF Daily Intensity",
    "entries": [
      {"dayPosition": 0, "percentage": 75.0},
      {"dayPosition": 1, "percentage": 80.0},
      {"dayPosition": 2, "percentage": 70.0},
      {"dayPosition": 3, "percentage": 85.0},
      {"dayPosition": 4, "percentage": 85.0},
      {"dayPosition": 5, "percentage": 65.0}
    ],
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

**Daily Lookup ID:** `dlookup-1111-1111-1111-111111111111`

**How Daily Lookups Work:**

When generating a workout, the system:
1. Determines the current day position within the week (0-5)
2. Looks up the matching entry for that day position
3. Uses the entry's percentage to calculate working weight

---

### Step 3: Create Weekly Lookup (Volume Progression)

The weekly lookup maps week number and day position to sets and reps. This creates the 3-week volume progression.

```bash
curl -X POST "http://localhost:8080/weekly-lookups" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Greg Nuckols HF Volume Progression",
    "entries": [
      {"weekNumber": 1, "dayPosition": 0, "sets": 4, "reps": 3, "isAmrap": false},
      {"weekNumber": 1, "dayPosition": 1, "sets": 3, "reps": 2, "isAmrap": false},
      {"weekNumber": 1, "dayPosition": 2, "sets": 4, "reps": 4, "isAmrap": false},
      {"weekNumber": 1, "dayPosition": 3, "sets": 3, "reps": 1, "isAmrap": false},
      {"weekNumber": 1, "dayPosition": 4, "sets": 3, "reps": 1, "isAmrap": false},
      {"weekNumber": 1, "dayPosition": 5, "sets": 5, "reps": 5, "isAmrap": false},

      {"weekNumber": 2, "dayPosition": 0, "sets": 5, "reps": 3, "isAmrap": false},
      {"weekNumber": 2, "dayPosition": 1, "sets": 4, "reps": 2, "isAmrap": false},
      {"weekNumber": 2, "dayPosition": 2, "sets": 5, "reps": 4, "isAmrap": false},
      {"weekNumber": 2, "dayPosition": 3, "sets": 4, "reps": 1, "isAmrap": false},
      {"weekNumber": 2, "dayPosition": 4, "sets": 4, "reps": 1, "isAmrap": false},
      {"weekNumber": 2, "dayPosition": 5, "sets": 6, "reps": 5, "isAmrap": false},

      {"weekNumber": 3, "dayPosition": 0, "sets": 5, "reps": 4, "isAmrap": false},
      {"weekNumber": 3, "dayPosition": 1, "sets": 4, "reps": 3, "isAmrap": false},
      {"weekNumber": 3, "dayPosition": 2, "sets": 5, "reps": 5, "isAmrap": false},
      {"weekNumber": 3, "dayPosition": 3, "sets": 1, "reps": 1, "isAmrap": true},
      {"weekNumber": 3, "dayPosition": 4, "sets": 1, "reps": 1, "isAmrap": true},
      {"weekNumber": 3, "dayPosition": 5, "sets": 6, "reps": 6, "isAmrap": false}
    ]
  }'
```

**Response:**
```json
{
  "data": {
    "id": "wlookup-hf-uuid",
    "name": "Greg Nuckols HF Volume Progression",
    "entries": [
      {"weekNumber": 1, "dayPosition": 0, "sets": 4, "reps": 3, "isAmrap": false},
      {"weekNumber": 1, "dayPosition": 1, "sets": 3, "reps": 2, "isAmrap": false},
      {"weekNumber": 1, "dayPosition": 2, "sets": 4, "reps": 4, "isAmrap": false},
      {"weekNumber": 1, "dayPosition": 3, "sets": 3, "reps": 1, "isAmrap": false},
      {"weekNumber": 1, "dayPosition": 4, "sets": 3, "reps": 1, "isAmrap": false},
      {"weekNumber": 1, "dayPosition": 5, "sets": 5, "reps": 5, "isAmrap": false},
      {"weekNumber": 2, "dayPosition": 0, "sets": 5, "reps": 3, "isAmrap": false},
      {"weekNumber": 2, "dayPosition": 1, "sets": 4, "reps": 2, "isAmrap": false},
      {"weekNumber": 2, "dayPosition": 2, "sets": 5, "reps": 4, "isAmrap": false},
      {"weekNumber": 2, "dayPosition": 3, "sets": 4, "reps": 1, "isAmrap": false},
      {"weekNumber": 2, "dayPosition": 4, "sets": 4, "reps": 1, "isAmrap": false},
      {"weekNumber": 2, "dayPosition": 5, "sets": 6, "reps": 5, "isAmrap": false},
      {"weekNumber": 3, "dayPosition": 0, "sets": 5, "reps": 4, "isAmrap": false},
      {"weekNumber": 3, "dayPosition": 1, "sets": 4, "reps": 3, "isAmrap": false},
      {"weekNumber": 3, "dayPosition": 2, "sets": 5, "reps": 5, "isAmrap": false},
      {"weekNumber": 3, "dayPosition": 3, "sets": 1, "reps": 1, "isAmrap": true},
      {"weekNumber": 3, "dayPosition": 4, "sets": 1, "reps": 1, "isAmrap": true},
      {"weekNumber": 3, "dayPosition": 5, "sets": 6, "reps": 6, "isAmrap": false}
    ],
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

**Weekly Lookup ID:** `wlookup-1111-1111-1111-111111111111`

**How Combined Lookups Work:**

When generating a workout, the system:
1. Determines current week number (1-3) and day position (0-5)
2. Looks up intensity from daily lookup using day position
3. Looks up sets/reps from weekly lookup using week number + day position
4. Calculates working weight: `weight = 1RM × intensity%`
5. Generates the specified number of sets at the calculated weight

---

### Step 4: Create Prescriptions

Create prescriptions that use the combined daily and weekly lookups.

**Bench Press Main Lift Prescription:**
```bash
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "11111111-1111-1111-1111-111111111111",
    "loadStrategy": {
      "type": "DAILY_LOOKUP",
      "maxType": "TRAINING_MAX",
      "roundTo": 2.5
    },
    "setScheme": {
      "type": "WEEKLY_LOOKUP"
    },
    "order": 1,
    "notes": "Focus on technique and bar speed. Week 3 heavy days are AMAP test sets.",
    "restSeconds": 180
  }'
```

**Response:**
```json
{
  "data": {
    "id": "rx-bench-main-uuid",
    "liftId": "11111111-1111-1111-1111-111111111111",
    "loadStrategy": {
      "type": "DAILY_LOOKUP",
      "maxType": "TRAINING_MAX",
      "roundTo": 2.5
    },
    "setScheme": {
      "type": "WEEKLY_LOOKUP"
    },
    "order": 1,
    "notes": "Focus on technique and bar speed. Week 3 heavy days are AMAP test sets.",
    "restSeconds": 180,
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

**Squat Main Lift Prescription:**
```bash
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "22222222-2222-2222-2222-222222222222",
    "loadStrategy": {
      "type": "DAILY_LOOKUP",
      "maxType": "TRAINING_MAX",
      "roundTo": 2.5
    },
    "setScheme": {
      "type": "WEEKLY_LOOKUP"
    },
    "order": 2,
    "notes": "Maintain depth and bracing. Week 3 heavy days are AMAP test sets.",
    "restSeconds": 180
  }'
```

**Deadlift Prescription (Friday only, 85%):**
```bash
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "33333333-3333-3333-3333-333333333333",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "TRAINING_MAX",
      "percentage": 85.0,
      "roundTo": 2.5
    },
    "setScheme": {
      "type": "WEEKLY_LOOKUP"
    },
    "order": 1,
    "notes": "Heavy singles - maintain tightness off the floor. Week 3 is AMAP.",
    "restSeconds": 240
  }'
```

**Overhead Press Prescription (Friday only, 85%):**
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
      "percentage": 85.0,
      "roundTo": 2.5
    },
    "setScheme": {
      "type": "WEEKLY_LOOKUP"
    },
    "order": 2,
    "notes": "Strict press - no leg drive. Week 3 is AMAP.",
    "restSeconds": 180
  }'
```

**Back Work Prescriptions (6RM-based):**

The back work uses a 3-week intensity progression based on 6RM:

```bash
# T-Bar Row - Week 1 (77.5% of 6RM, 4x12)
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "55555555-5555-5555-5555-555555555555",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "SIX_REP_MAX",
      "percentage": 77.5,
      "roundTo": 2.5
    },
    "setScheme": {
      "type": "FIXED",
      "sets": 4,
      "reps": 12,
      "isAmrap": false
    },
    "order": 3,
    "notes": "Controlled movement, full stretch at bottom.",
    "restSeconds": 90
  }'

# Pulldown - Week 1 (77.5% of 6RM, 4x12)
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "66666666-6666-6666-6666-666666666666",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "SIX_REP_MAX",
      "percentage": 77.5,
      "roundTo": 2.5
    },
    "setScheme": {
      "type": "FIXED",
      "sets": 4,
      "reps": 12,
      "isAmrap": false
    },
    "order": 3,
    "notes": "Full range of motion, squeeze at bottom.",
    "restSeconds": 90
  }'
```

**Back Work - Week 2 (85% of 6RM, 4x8):**
```bash
# T-Bar Row Week 2
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "55555555-5555-5555-5555-555555555555",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "SIX_REP_MAX",
      "percentage": 85.0,
      "roundTo": 2.5
    },
    "setScheme": {
      "type": "FIXED",
      "sets": 4,
      "reps": 8,
      "isAmrap": false
    },
    "order": 3,
    "notes": "Heavier weight, maintain form.",
    "restSeconds": 90
  }'

# Pulldown Week 2
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "66666666-6666-6666-6666-666666666666",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "SIX_REP_MAX",
      "percentage": 85.0,
      "roundTo": 2.5
    },
    "setScheme": {
      "type": "FIXED",
      "sets": 4,
      "reps": 8,
      "isAmrap": false
    },
    "order": 3,
    "notes": "Heavier weight, maintain form.",
    "restSeconds": 90
  }'
```

**Back Work - Week 3 (100% of 6RM, 1+ sets of 6):**
```bash
# T-Bar Row Week 3 (AMAP sets)
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "55555555-5555-5555-5555-555555555555",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "SIX_REP_MAX",
      "percentage": 100.0,
      "roundTo": 2.5
    },
    "setScheme": {
      "type": "FIXED",
      "sets": 1,
      "reps": 6,
      "isAmrap": true
    },
    "order": 3,
    "notes": "Test set - continue for additional sets of 6 to gauge progress. +2% 6RM per additional set completed.",
    "restSeconds": 120
  }'

# Pulldown Week 3 (AMAP sets)
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "66666666-6666-6666-6666-666666666666",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "SIX_REP_MAX",
      "percentage": 100.0,
      "roundTo": 2.5
    },
    "setScheme": {
      "type": "FIXED",
      "sets": 1,
      "reps": 6,
      "isAmrap": true
    },
    "order": 3,
    "notes": "Test set - continue for additional sets of 6 to gauge progress. +2% 6RM per additional set completed.",
    "restSeconds": 120
  }'
```

**Prescription IDs for subsequent steps** (example UUIDs):
| Prescription | ID |
|--------------|-----|
| Bench Main | `aaaa1111-1111-1111-1111-111111111111` |
| Squat Main | `aaaa2222-2222-2222-2222-222222222222` |
| Deadlift | `aaaa3333-3333-3333-3333-333333333333` |
| OHP | `aaaa4444-4444-4444-4444-444444444444` |
| T-Bar Row W1 | `aaaa5555-5555-5555-5555-555555555555` |
| Pulldown W1 | `aaaa6666-6666-6666-6666-666666666666` |
| T-Bar Row W2 | `aaaa7777-7777-7777-7777-777777777777` |
| Pulldown W2 | `aaaa8888-8888-8888-8888-888888888888` |
| T-Bar Row W3 | `aaaa9999-9999-9999-9999-999999999999` |
| Pulldown W3 | `aaaaAAAA-AAAA-AAAA-AAAA-AAAAAAAAAAAA` |

---

### Step 5: Create Days

Create the six training days.

**Monday (Bench + Squat, 75%):**
```bash
curl -X POST "http://localhost:8080/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Monday - Moderate",
    "slug": "hf-monday",
    "metadata": {
      "intensity": "75%",
      "character": "technical-practice",
      "program": "greg-nuckols-hf"
    }
  }'
```

**Response:**
```json
{
  "data": {
    "id": "day-monday-uuid",
    "name": "Monday - Moderate",
    "slug": "hf-monday",
    "metadata": {
      "intensity": "75%",
      "character": "technical-practice",
      "program": "greg-nuckols-hf"
    },
    "programId": null,
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

```bash
# Tuesday (Bench + Squat, 80%)
curl -X POST "http://localhost:8080/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Tuesday - Moderate-Heavy",
    "slug": "hf-tuesday",
    "metadata": {
      "intensity": "80%",
      "character": "strength-building",
      "program": "greg-nuckols-hf"
    }
  }'

# Wednesday (Bench + Squat, 70%)
curl -X POST "http://localhost:8080/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Wednesday - Light",
    "slug": "hf-wednesday",
    "metadata": {
      "intensity": "70%",
      "character": "recovery-volume",
      "program": "greg-nuckols-hf"
    }
  }'

# Thursday (Bench + Squat + Pulldown, 85%)
curl -X POST "http://localhost:8080/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Thursday - Heavy",
    "slug": "hf-thursday",
    "metadata": {
      "intensity": "85%",
      "character": "intensity-day",
      "program": "greg-nuckols-hf"
    }
  }'

# Friday (OHP + Deadlift, 85%)
curl -X POST "http://localhost:8080/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Friday - Pull Focus",
    "slug": "hf-friday",
    "metadata": {
      "intensity": "85%",
      "character": "pull-focus",
      "program": "greg-nuckols-hf"
    }
  }'

# Saturday (Bench + Squat + Row, 65%)
curl -X POST "http://localhost:8080/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Saturday - Volume",
    "slug": "hf-saturday",
    "metadata": {
      "intensity": "65%",
      "character": "high-volume-hypertrophy",
      "program": "greg-nuckols-hf"
    }
  }'
```

**Day IDs:**
| Day | ID |
|-----|-----|
| Monday | `bbbb1111-1111-1111-1111-111111111111` |
| Tuesday | `bbbb2222-2222-2222-2222-222222222222` |
| Wednesday | `bbbb3333-3333-3333-3333-333333333333` |
| Thursday | `bbbb4444-4444-4444-4444-444444444444` |
| Friday | `bbbb5555-5555-5555-5555-555555555555` |
| Saturday | `bbbb6666-6666-6666-6666-666666666666` |

**Add prescriptions to Monday:**
```bash
# Bench Press
curl -X POST "http://localhost:8080/days/bbbb1111-1111-1111-1111-111111111111/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"prescriptionId": "aaaa1111-1111-1111-1111-111111111111", "order": 1}'

# Squat
curl -X POST "http://localhost:8080/days/bbbb1111-1111-1111-1111-111111111111/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"prescriptionId": "aaaa2222-2222-2222-2222-222222222222", "order": 2}'
```

**Add prescriptions to Tuesday:**
```bash
curl -X POST "http://localhost:8080/days/bbbb2222-2222-2222-2222-222222222222/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"prescriptionId": "aaaa1111-1111-1111-1111-111111111111", "order": 1}'

curl -X POST "http://localhost:8080/days/bbbb2222-2222-2222-2222-222222222222/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"prescriptionId": "aaaa2222-2222-2222-2222-222222222222", "order": 2}'
```

**Add prescriptions to Wednesday:**
```bash
curl -X POST "http://localhost:8080/days/bbbb3333-3333-3333-3333-333333333333/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"prescriptionId": "aaaa1111-1111-1111-1111-111111111111", "order": 1}'

curl -X POST "http://localhost:8080/days/bbbb3333-3333-3333-3333-333333333333/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"prescriptionId": "aaaa2222-2222-2222-2222-222222222222", "order": 2}'
```

**Add prescriptions to Thursday (Bench + Squat + Pulldown):**
```bash
curl -X POST "http://localhost:8080/days/bbbb4444-4444-4444-4444-444444444444/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"prescriptionId": "aaaa1111-1111-1111-1111-111111111111", "order": 1}'

curl -X POST "http://localhost:8080/days/bbbb4444-4444-4444-4444-444444444444/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"prescriptionId": "aaaa2222-2222-2222-2222-222222222222", "order": 2}'

# Pulldown (week-specific prescription linked at week level)
```

**Add prescriptions to Friday (OHP + Deadlift):**
```bash
curl -X POST "http://localhost:8080/days/bbbb5555-5555-5555-5555-555555555555/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"prescriptionId": "aaaa4444-4444-4444-4444-444444444444", "order": 1}'

curl -X POST "http://localhost:8080/days/bbbb5555-5555-5555-5555-555555555555/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"prescriptionId": "aaaa3333-3333-3333-3333-333333333333", "order": 2}'
```

**Add prescriptions to Saturday (Bench + Squat + Row):**
```bash
curl -X POST "http://localhost:8080/days/bbbb6666-6666-6666-6666-666666666666/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"prescriptionId": "aaaa1111-1111-1111-1111-111111111111", "order": 1}'

curl -X POST "http://localhost:8080/days/bbbb6666-6666-6666-6666-666666666666/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"prescriptionId": "aaaa2222-2222-2222-2222-222222222222", "order": 2}'

# T-Bar Row (week-specific prescription linked at week level)
```

---

### Step 6: Create Cycle

Greg Nuckols HF uses a 3-week cycle.

```bash
curl -X POST "http://localhost:8080/cycles" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Greg Nuckols HF 3-Week Cycle",
    "lengthWeeks": 3
  }'
```

**Response:**
```json
{
  "data": {
    "id": "cccc1111-1111-1111-1111-111111111111",
    "name": "Greg Nuckols HF 3-Week Cycle",
    "lengthWeeks": 3,
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

---

### Step 7: Create Weeks

Create all three weeks with the same day structure. The lookups handle the volume/intensity differences.

```bash
# Week 1 (Base)
curl -X POST "http://localhost:8080/weeks" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "cycleId": "cccc1111-1111-1111-1111-111111111111",
    "weekNumber": 1,
    "name": "Base Volume"
  }'

# Week 2 (Increased Volume)
curl -X POST "http://localhost:8080/weeks" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "cycleId": "cccc1111-1111-1111-1111-111111111111",
    "weekNumber": 2,
    "name": "Increased Volume"
  }'

# Week 3 (Peak/Test)
curl -X POST "http://localhost:8080/weeks" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "cycleId": "cccc1111-1111-1111-1111-111111111111",
    "weekNumber": 3,
    "name": "Peak/Test Week"
  }'
```

**Week IDs:**
| Week | ID |
|------|-----|
| Week 1 (Base) | `dddd1111-1111-1111-1111-111111111111` |
| Week 2 (Increased) | `dddd2222-2222-2222-2222-222222222222` |
| Week 3 (Peak/Test) | `dddd3333-3333-3333-3333-333333333333` |

**Add days to each week:**

```bash
# Week 1 days
for pos in 0 1 2 3 4 5; do
  dayId=$(echo "bbbb$(printf '%04d' $((pos+1)))-$(printf '%04d' $((pos+1)))-$(printf '%04d' $((pos+1)))-$(printf '%04d' $((pos+1)))-$(printf '%012d' $((pos+1)))")
  # Note: Use the actual day IDs from above
done

# Week 1 - Monday through Saturday
curl -X POST "http://localhost:8080/weeks/dddd1111-1111-1111-1111-111111111111/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "bbbb1111-1111-1111-1111-111111111111", "position": 0}'

curl -X POST "http://localhost:8080/weeks/dddd1111-1111-1111-1111-111111111111/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "bbbb2222-2222-2222-2222-222222222222", "position": 1}'

curl -X POST "http://localhost:8080/weeks/dddd1111-1111-1111-1111-111111111111/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "bbbb3333-3333-3333-3333-333333333333", "position": 2}'

curl -X POST "http://localhost:8080/weeks/dddd1111-1111-1111-1111-111111111111/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "bbbb4444-4444-4444-4444-444444444444", "position": 3}'

curl -X POST "http://localhost:8080/weeks/dddd1111-1111-1111-1111-111111111111/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "bbbb5555-5555-5555-5555-555555555555", "position": 4}'

curl -X POST "http://localhost:8080/weeks/dddd1111-1111-1111-1111-111111111111/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "bbbb6666-6666-6666-6666-666666666666", "position": 5}'

# Week 2 - same day structure
curl -X POST "http://localhost:8080/weeks/dddd2222-2222-2222-2222-222222222222/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "bbbb1111-1111-1111-1111-111111111111", "position": 0}'

curl -X POST "http://localhost:8080/weeks/dddd2222-2222-2222-2222-222222222222/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "bbbb2222-2222-2222-2222-222222222222", "position": 1}'

curl -X POST "http://localhost:8080/weeks/dddd2222-2222-2222-2222-222222222222/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "bbbb3333-3333-3333-3333-333333333333", "position": 2}'

curl -X POST "http://localhost:8080/weeks/dddd2222-2222-2222-2222-222222222222/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "bbbb4444-4444-4444-4444-444444444444", "position": 3}'

curl -X POST "http://localhost:8080/weeks/dddd2222-2222-2222-2222-222222222222/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "bbbb5555-5555-5555-5555-555555555555", "position": 4}'

curl -X POST "http://localhost:8080/weeks/dddd2222-2222-2222-2222-222222222222/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "bbbb6666-6666-6666-6666-666666666666", "position": 5}'

# Week 3 - same day structure
curl -X POST "http://localhost:8080/weeks/dddd3333-3333-3333-3333-333333333333/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "bbbb1111-1111-1111-1111-111111111111", "position": 0}'

curl -X POST "http://localhost:8080/weeks/dddd3333-3333-3333-3333-333333333333/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "bbbb2222-2222-2222-2222-222222222222", "position": 1}'

curl -X POST "http://localhost:8080/weeks/dddd3333-3333-3333-3333-333333333333/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "bbbb3333-3333-3333-3333-333333333333", "position": 2}'

curl -X POST "http://localhost:8080/weeks/dddd3333-3333-3333-3333-333333333333/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "bbbb4444-4444-4444-4444-444444444444", "position": 3}'

curl -X POST "http://localhost:8080/weeks/dddd3333-3333-3333-3333-333333333333/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "bbbb5555-5555-5555-5555-555555555555", "position": 4}'

curl -X POST "http://localhost:8080/weeks/dddd3333-3333-3333-3333-333333333333/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "bbbb6666-6666-6666-6666-666666666666", "position": 5}'
```

---

### Step 8: Create Progressions

Greg Nuckols uses AMAP-guided cycle-end progression.

**Main Lift Progression (AMAP-based):**
```bash
curl -X POST "http://localhost:8080/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Nuckols AMAP Progression",
    "type": "AMAP_GUIDED",
    "parameters": {
      "maxType": "TRAINING_MAX",
      "triggerType": "AFTER_CYCLE",
      "repRanges": [
        {"minReps": 6, "maxReps": 7, "increment": 2.5},
        {"minReps": 8, "maxReps": 9, "increment": 5.0},
        {"minReps": 10, "maxReps": 999, "increment": 7.5}
      ]
    }
  }'
```

**Response:**
```json
{
  "data": {
    "id": "prog-nuckols-amap-uuid",
    "name": "Nuckols AMAP Progression",
    "type": "AMAP_GUIDED",
    "parameters": {
      "maxType": "TRAINING_MAX",
      "triggerType": "AFTER_CYCLE",
      "repRanges": [
        {"minReps": 6, "maxReps": 7, "increment": 2.5},
        {"minReps": 8, "maxReps": 9, "increment": 5.0},
        {"minReps": 10, "maxReps": 999, "increment": 7.5}
      ]
    },
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

**Back Work Progression (Set-based 6RM increase):**
```bash
curl -X POST "http://localhost:8080/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Back Work 6RM Progression (+2% per extra set)",
    "type": "SET_BASED",
    "parameters": {
      "maxType": "SIX_REP_MAX",
      "triggerType": "AFTER_CYCLE",
      "incrementPerSet": 2.0
    }
  }'
```

**Progression IDs:**
| Progression | ID |
|-------------|-----|
| AMAP Main Lifts | `eeee1111-1111-1111-1111-111111111111` |
| Back Work 6RM | `eeee2222-2222-2222-2222-222222222222` |

---

### Step 9: Create Program

Create the Greg Nuckols HF program, linking cycle and both lookup tables.

```bash
curl -X POST "http://localhost:8080/programs" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Greg Nuckols High Frequency",
    "slug": "greg-nuckols-hf",
    "description": "Greg Nuckols'\'' High Frequency program for intermediate to advanced lifters. Features daily undulating periodization with 3-week volume progression cycles. Bench and squat trained 5x per week with varying intensities each day.",
    "cycleId": "cccc1111-1111-1111-1111-111111111111",
    "dailyLookupId": "dlookup-1111-1111-1111-111111111111",
    "weeklyLookupId": "wlookup-1111-1111-1111-111111111111",
    "defaultRounding": 2.5
  }'
```

**Response:**
```json
{
  "data": {
    "id": "program-hf-uuid",
    "name": "Greg Nuckols High Frequency",
    "slug": "greg-nuckols-hf",
    "description": "Greg Nuckols' High Frequency program for intermediate to advanced lifters. Features daily undulating periodization with 3-week volume progression cycles. Bench and squat trained 5x per week with varying intensities each day.",
    "cycleId": "cccc1111-1111-1111-1111-111111111111",
    "dailyLookupId": "dlookup-1111-1111-1111-111111111111",
    "weeklyLookupId": "wlookup-1111-1111-1111-111111111111",
    "defaultRounding": 2.5,
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

**Program ID:** `ffff1111-1111-1111-1111-111111111111`

---

### Step 10: Link Progressions to Program

Link each progression to the appropriate lifts.

```bash
# Bench Press Progression
curl -X POST "http://localhost:8080/programs/ffff1111-1111-1111-1111-111111111111/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "eeee1111-1111-1111-1111-111111111111",
    "liftId": "11111111-1111-1111-1111-111111111111",
    "priority": 1
  }'

# Squat Progression
curl -X POST "http://localhost:8080/programs/ffff1111-1111-1111-1111-111111111111/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "eeee1111-1111-1111-1111-111111111111",
    "liftId": "22222222-2222-2222-2222-222222222222",
    "priority": 1
  }'

# Deadlift Progression
curl -X POST "http://localhost:8080/programs/ffff1111-1111-1111-1111-111111111111/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "eeee1111-1111-1111-1111-111111111111",
    "liftId": "33333333-3333-3333-3333-333333333333",
    "priority": 1
  }'

# OHP Progression
curl -X POST "http://localhost:8080/programs/ffff1111-1111-1111-1111-111111111111/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "eeee1111-1111-1111-1111-111111111111",
    "liftId": "44444444-4444-4444-4444-444444444444",
    "priority": 1
  }'

# T-Bar Row Progression (6RM)
curl -X POST "http://localhost:8080/programs/ffff1111-1111-1111-1111-111111111111/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "eeee2222-2222-2222-2222-222222222222",
    "liftId": "55555555-5555-5555-5555-555555555555",
    "priority": 1
  }'

# Pulldown Progression (6RM)
curl -X POST "http://localhost:8080/programs/ffff1111-1111-1111-1111-111111111111/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "eeee2222-2222-2222-2222-222222222222",
    "liftId": "66666666-6666-6666-6666-666666666666",
    "priority": 1
  }'
```

---

## User Setup

Before a user can generate workouts, they need training maxes for all lifts.

### Create User Training Maxes

```bash
# User ID for examples
USER_ID="user-12345678-1234-1234-1234-123456789012"

# Bench Press Training Max: 100kg
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "11111111-1111-1111-1111-111111111111",
    "type": "TRAINING_MAX",
    "value": 100.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'

# Squat Training Max: 140kg
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "22222222-2222-2222-2222-222222222222",
    "type": "TRAINING_MAX",
    "value": 140.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'

# Deadlift Training Max: 180kg
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "33333333-3333-3333-3333-333333333333",
    "type": "TRAINING_MAX",
    "value": 180.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'

# OHP Training Max: 60kg
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "44444444-4444-4444-4444-444444444444",
    "type": "TRAINING_MAX",
    "value": 60.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'

# T-Bar Row 6RM: 80kg
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "55555555-5555-5555-5555-555555555555",
    "type": "SIX_REP_MAX",
    "value": 80.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'

# Pulldown 6RM: 70kg
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "66666666-6666-6666-6666-666666666666",
    "type": "SIX_REP_MAX",
    "value": 70.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'
```

### Enroll User in Program

```bash
curl -X POST "http://localhost:8080/users/${USER_ID}/program" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "programId": "ffff1111-1111-1111-1111-111111111111"
  }'
```

**Response:**
```json
{
  "data": {
    "id": "enrollment-uuid",
    "userId": "user-12345678-1234-1234-1234-123456789012",
    "program": {
      "id": "ffff1111-1111-1111-1111-111111111111",
      "name": "Greg Nuckols High Frequency",
      "slug": "greg-nuckols-hf",
      "description": "Greg Nuckols' High Frequency program for intermediate to advanced lifters.",
      "cycleLengthWeeks": 3
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

### Week 1 Daily Variation

This demonstrates how **intensity changes daily** while volume stays consistent within the week.

Training maxes: Bench 100kg, Squat 140kg

**Week 1, Day 1 (Monday - 75%):**

```bash
curl -X GET "http://localhost:8080/users/${USER_ID}/workout" \
  -H "X-User-ID: ${USER_ID}"
```

**Response:**
```json
{
  "data": {
    "userId": "user-12345678-1234-1234-1234-123456789012",
    "programId": "ffff1111-1111-1111-1111-111111111111",
    "cycleIteration": 1,
    "weekNumber": 1,
    "weekName": "Base Volume",
    "daySlug": "hf-monday",
    "dayPosition": 0,
    "date": "2024-01-15",
    "exercises": [
      {
        "prescriptionId": "aaaa1111-1111-1111-1111-111111111111",
        "lift": {
          "id": "11111111-1111-1111-1111-111111111111",
          "name": "Bench Press",
          "slug": "bench-press"
        },
        "sets": [
          {"setNumber": 1, "weight": 75.0, "targetReps": 3, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 2, "weight": 75.0, "targetReps": 3, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 3, "weight": 75.0, "targetReps": 3, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 4, "weight": 75.0, "targetReps": 3, "isAmrap": false, "isWorkSet": true}
        ],
        "notes": "Focus on technique and bar speed. Week 3 heavy days are AMAP test sets.",
        "restSeconds": 180
      },
      {
        "prescriptionId": "aaaa2222-2222-2222-2222-222222222222",
        "lift": {
          "id": "22222222-2222-2222-2222-222222222222",
          "name": "Squat",
          "slug": "squat"
        },
        "sets": [
          {"setNumber": 1, "weight": 105.0, "targetReps": 3, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 2, "weight": 105.0, "targetReps": 3, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 3, "weight": 105.0, "targetReps": 3, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 4, "weight": 105.0, "targetReps": 3, "isAmrap": false, "isWorkSet": true}
        ],
        "notes": "Maintain depth and bracing. Week 3 heavy days are AMAP test sets.",
        "restSeconds": 180
      }
    ]
  }
}
```

**Monday Weight Calculation:**
| Lift | TM | Intensity | Calculation | Weight |
|------|-----|-----------|-------------|--------|
| Bench | 100kg | 75% | 100 × 0.75 | 75kg |
| Squat | 140kg | 75% | 140 × 0.75 | 105kg |

**Volume (from weekly lookup, week 1, day position 0):** 4 sets × 3 reps

---

**Week 1, Day 2 (Tuesday - 80%):**

After advancing:
```bash
curl -X POST "http://localhost:8080/users/${USER_ID}/program-state/advance" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{"advanceType": "day"}'
```

**Response:**
```json
{
  "data": {
    "weekNumber": 1,
    "weekName": "Base Volume",
    "daySlug": "hf-tuesday",
    "dayPosition": 1,
    "exercises": [
      {
        "lift": {"name": "Bench Press"},
        "sets": [
          {"setNumber": 1, "weight": 80.0, "targetReps": 2, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 2, "weight": 80.0, "targetReps": 2, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 3, "weight": 80.0, "targetReps": 2, "isAmrap": false, "isWorkSet": true}
        ]
      },
      {
        "lift": {"name": "Squat"},
        "sets": [
          {"setNumber": 1, "weight": 112.5, "targetReps": 2, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 2, "weight": 112.5, "targetReps": 2, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 3, "weight": 112.5, "targetReps": 2, "isAmrap": false, "isWorkSet": true}
        ]
      }
    ]
  }
}
```

**Tuesday Weight Calculation:**
| Lift | TM | Intensity | Calculation | Weight |
|------|-----|-----------|-------------|--------|
| Bench | 100kg | 80% | 100 × 0.80 | 80kg |
| Squat | 140kg | 80% | 140 × 0.80 | 112.5kg (rounded) |

**Volume (week 1, day position 1):** 3 sets × 2 reps

---

**Week 1 Complete Daily Variation Summary:**

| Day | Intensity | Bench Weight | Squat Weight | Sets × Reps |
|-----|-----------|--------------|--------------|-------------|
| Monday | 75% | 75kg | 105kg | 4×3 |
| Tuesday | 80% | 80kg | 112.5kg | 3×2 |
| Wednesday | 70% | 70kg | 97.5kg | 4×4 |
| Thursday | 85% | 85kg | 120kg | 3×1 |
| Saturday | 65% | 65kg | 92.5kg | 5×5 |

This demonstrates **daily undulating periodization** - the same lifts are performed with different intensities each day, providing varied stimuli while maintaining high frequency.

---

### Week 2 Volume Increase

After advancing to Week 2:

```bash
curl -X POST "http://localhost:8080/users/${USER_ID}/program-state/advance" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{"advanceType": "week"}'
```

**Week 2, Day 1 (Monday - 75%):**

```json
{
  "data": {
    "weekNumber": 2,
    "weekName": "Increased Volume",
    "daySlug": "hf-monday",
    "exercises": [
      {
        "lift": {"name": "Bench Press"},
        "sets": [
          {"setNumber": 1, "weight": 75.0, "targetReps": 3, "isAmrap": false},
          {"setNumber": 2, "weight": 75.0, "targetReps": 3, "isAmrap": false},
          {"setNumber": 3, "weight": 75.0, "targetReps": 3, "isAmrap": false},
          {"setNumber": 4, "weight": 75.0, "targetReps": 3, "isAmrap": false},
          {"setNumber": 5, "weight": 75.0, "targetReps": 3, "isAmrap": false}
        ]
      },
      {
        "lift": {"name": "Squat"},
        "sets": [
          {"setNumber": 1, "weight": 105.0, "targetReps": 3, "isAmrap": false},
          {"setNumber": 2, "weight": 105.0, "targetReps": 3, "isAmrap": false},
          {"setNumber": 3, "weight": 105.0, "targetReps": 3, "isAmrap": false},
          {"setNumber": 4, "weight": 105.0, "targetReps": 3, "isAmrap": false},
          {"setNumber": 5, "weight": 105.0, "targetReps": 3, "isAmrap": false}
        ]
      }
    ]
  }
}
```

**Week 2 Volume Comparison (Monday):**

| Week | Intensity | Sets | Reps | Total Volume |
|------|-----------|------|------|--------------|
| Week 1 | 75% | 4 | 3 | 12 reps |
| Week 2 | 75% | 5 | 3 | 15 reps (+25%) |

The intensity stays the same (75%), but volume increases from 4 to 5 sets.

---

### Week 3 Peak/Test Week

**Week 3, Day 4 (Thursday - 85% AMAP):**

After advancing to Week 3, Thursday:

```json
{
  "data": {
    "weekNumber": 3,
    "weekName": "Peak/Test Week",
    "daySlug": "hf-thursday",
    "dayPosition": 3,
    "exercises": [
      {
        "lift": {"name": "Bench Press"},
        "sets": [
          {"setNumber": 1, "weight": 85.0, "targetReps": 1, "isAmrap": true, "isWorkSet": true}
        ],
        "notes": "AMAP test - go for as many reps as possible with good form."
      },
      {
        "lift": {"name": "Squat"},
        "sets": [
          {"setNumber": 1, "weight": 120.0, "targetReps": 1, "isAmrap": true, "isWorkSet": true}
        ],
        "notes": "AMAP test - go for as many reps as possible with good form."
      }
    ]
  }
}
```

**Week 3 Thursday (Heavy Day) Comparison:**

| Week | Intensity | Sets | Reps | Notes |
|------|-----------|------|------|-------|
| Week 1 | 85% | 3 | 1 | Fixed singles |
| Week 2 | 85% | 4 | 1 | More singles |
| Week 3 | 85% | 1 | AMAP | Test set for progression |

The Week 3 Thursday heavy day converts from multiple singles to a single AMAP set, allowing performance-based progression decisions.

---

## Cycle Progression Behavior

After completing all 3 weeks, use the AMAP results from Week 3 to determine training max increases.

### AMAP Results and Progression

**AMAP Progression Table:**
| Reps Achieved at 85% | Training Max Increase |
|----------------------|----------------------|
| 6-7 reps | +2.5kg |
| 8-9 reps | +5.0kg |
| 10+ reps | +7.5kg |

### Triggering Progression

```bash
# Example: User hit 8 reps on bench at 85kg (85% of 100kg TM)
curl -X POST "http://localhost:8080/users/${USER_ID}/progressions/trigger" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "eeee1111-1111-1111-1111-111111111111",
    "liftId": "11111111-1111-1111-1111-111111111111",
    "amapReps": 8
  }'
```

**Response:**
```json
{
  "data": {
    "results": [
      {
        "progressionId": "eeee1111-1111-1111-1111-111111111111",
        "liftId": "11111111-1111-1111-1111-111111111111",
        "applied": true,
        "result": {
          "previousValue": 100.0,
          "newValue": 105.0,
          "delta": 5.0,
          "maxType": "TRAINING_MAX",
          "amapReps": 8,
          "appliedAt": "2024-02-05T18:00:00Z"
        }
      }
    ]
  }
}
```

### Training Max Progression Over Cycles

| Cycle | Bench TM | Squat TM | Deadlift TM | OHP TM |
|-------|----------|----------|-------------|--------|
| 1 | 100kg | 140kg | 180kg | 60kg |
| 2 | 105kg | 147.5kg | 187.5kg | 62.5kg |
| 3 | 110kg | 155kg | 195kg | 65kg |

---

## Complete Configuration Summary

### Entity Relationship Diagram

```
Program: Greg Nuckols High Frequency
├── Cycle: 3-week cycle
│   ├── Week 1 (Base Volume)
│   │   ├── Position 0: Monday (Bench+Squat @ 75%)
│   │   ├── Position 1: Tuesday (Bench+Squat @ 80%)
│   │   ├── Position 2: Wednesday (Bench+Squat @ 70%)
│   │   ├── Position 3: Thursday (Bench+Squat+Pulldown @ 85%)
│   │   ├── Position 4: Friday (OHP+Deadlift @ 85%)
│   │   └── Position 5: Saturday (Bench+Squat+Row @ 65%)
│   ├── Week 2 (Increased Volume)
│   │   └── [same day structure, +1 set per day]
│   └── Week 3 (Peak/Test)
│       └── [same day structure, heavy days become AMAP]
│
├── Daily Lookup: Intensity by Day Position
│   ├── Day 0 (Mon): 75%
│   ├── Day 1 (Tue): 80%
│   ├── Day 2 (Wed): 70%
│   ├── Day 3 (Thu): 85%
│   ├── Day 4 (Fri): 85%
│   └── Day 5 (Sat): 65%
│
├── Weekly Lookup: Volume by Week + Day
│   ├── Week 1: Base sets (4, 3, 4, 3, 3, 5)
│   ├── Week 2: +1 set (5, 4, 5, 4, 4, 6)
│   └── Week 3: Peak (5, 4, 5, 1 AMAP, 1 AMAP, 6)
│
└── Program Progressions
    ├── Bench → AMAP-Guided (AFTER_CYCLE)
    ├── Squat → AMAP-Guided (AFTER_CYCLE)
    ├── Deadlift → AMAP-Guided (AFTER_CYCLE)
    ├── OHP → AMAP-Guided (AFTER_CYCLE)
    ├── T-Bar Row → Set-Based 6RM (AFTER_CYCLE)
    └── Pulldown → Set-Based 6RM (AFTER_CYCLE)
```

### JSON Configuration Export

```json
{
  "program": {
    "name": "Greg Nuckols High Frequency",
    "slug": "greg-nuckols-hf",
    "description": "Greg Nuckols' High Frequency program. Daily undulating periodization with 3-week volume progression.",
    "defaultRounding": 2.5
  },
  "cycle": {
    "name": "Greg Nuckols HF 3-Week Cycle",
    "lengthWeeks": 3
  },
  "dailyLookup": {
    "name": "Greg Nuckols HF Daily Intensity",
    "entries": [
      {"dayPosition": 0, "percentage": 75.0},
      {"dayPosition": 1, "percentage": 80.0},
      {"dayPosition": 2, "percentage": 70.0},
      {"dayPosition": 3, "percentage": 85.0},
      {"dayPosition": 4, "percentage": 85.0},
      {"dayPosition": 5, "percentage": 65.0}
    ]
  },
  "weeklyLookup": {
    "name": "Greg Nuckols HF Volume Progression",
    "entries": [
      {"weekNumber": 1, "dayPosition": 0, "sets": 4, "reps": 3, "isAmrap": false},
      {"weekNumber": 1, "dayPosition": 1, "sets": 3, "reps": 2, "isAmrap": false},
      {"weekNumber": 1, "dayPosition": 2, "sets": 4, "reps": 4, "isAmrap": false},
      {"weekNumber": 1, "dayPosition": 3, "sets": 3, "reps": 1, "isAmrap": false},
      {"weekNumber": 1, "dayPosition": 4, "sets": 3, "reps": 1, "isAmrap": false},
      {"weekNumber": 1, "dayPosition": 5, "sets": 5, "reps": 5, "isAmrap": false},
      {"weekNumber": 2, "dayPosition": 0, "sets": 5, "reps": 3, "isAmrap": false},
      {"weekNumber": 2, "dayPosition": 1, "sets": 4, "reps": 2, "isAmrap": false},
      {"weekNumber": 2, "dayPosition": 2, "sets": 5, "reps": 4, "isAmrap": false},
      {"weekNumber": 2, "dayPosition": 3, "sets": 4, "reps": 1, "isAmrap": false},
      {"weekNumber": 2, "dayPosition": 4, "sets": 4, "reps": 1, "isAmrap": false},
      {"weekNumber": 2, "dayPosition": 5, "sets": 6, "reps": 5, "isAmrap": false},
      {"weekNumber": 3, "dayPosition": 0, "sets": 5, "reps": 4, "isAmrap": false},
      {"weekNumber": 3, "dayPosition": 1, "sets": 4, "reps": 3, "isAmrap": false},
      {"weekNumber": 3, "dayPosition": 2, "sets": 5, "reps": 5, "isAmrap": false},
      {"weekNumber": 3, "dayPosition": 3, "sets": 1, "reps": 1, "isAmrap": true},
      {"weekNumber": 3, "dayPosition": 4, "sets": 1, "reps": 1, "isAmrap": true},
      {"weekNumber": 3, "dayPosition": 5, "sets": 6, "reps": 6, "isAmrap": false}
    ]
  },
  "weeks": [
    {"weekNumber": 1, "name": "Base Volume", "days": ["hf-monday", "hf-tuesday", "hf-wednesday", "hf-thursday", "hf-friday", "hf-saturday"]},
    {"weekNumber": 2, "name": "Increased Volume", "days": ["hf-monday", "hf-tuesday", "hf-wednesday", "hf-thursday", "hf-friday", "hf-saturday"]},
    {"weekNumber": 3, "name": "Peak/Test Week", "days": ["hf-monday", "hf-tuesday", "hf-wednesday", "hf-thursday", "hf-friday", "hf-saturday"]}
  ],
  "days": [
    {"name": "Monday - Moderate", "slug": "hf-monday", "prescriptions": ["bench-main", "squat-main"]},
    {"name": "Tuesday - Moderate-Heavy", "slug": "hf-tuesday", "prescriptions": ["bench-main", "squat-main"]},
    {"name": "Wednesday - Light", "slug": "hf-wednesday", "prescriptions": ["bench-main", "squat-main"]},
    {"name": "Thursday - Heavy", "slug": "hf-thursday", "prescriptions": ["bench-main", "squat-main", "pulldown"]},
    {"name": "Friday - Pull Focus", "slug": "hf-friday", "prescriptions": ["ohp", "deadlift"]},
    {"name": "Saturday - Volume", "slug": "hf-saturday", "prescriptions": ["bench-main", "squat-main", "t-bar-row"]}
  ],
  "progressions": [
    {
      "name": "Nuckols AMAP Progression",
      "type": "AMAP_GUIDED",
      "parameters": {
        "maxType": "TRAINING_MAX",
        "triggerType": "AFTER_CYCLE",
        "repRanges": [
          {"minReps": 6, "maxReps": 7, "increment": 2.5},
          {"minReps": 8, "maxReps": 9, "increment": 5.0},
          {"minReps": 10, "maxReps": 999, "increment": 7.5}
        ]
      },
      "liftSlugs": ["bench-press", "squat", "deadlift", "overhead-press"]
    },
    {
      "name": "Back Work 6RM Progression",
      "type": "SET_BASED",
      "parameters": {
        "maxType": "SIX_REP_MAX",
        "triggerType": "AFTER_CYCLE",
        "incrementPerSet": 2.0
      },
      "liftSlugs": ["t-bar-row", "pulldown"]
    }
  ]
}
```

---

## See Also

- [API Reference](../api-reference.md) - Complete endpoint documentation
- [Workflows](../workflows.md) - Multi-step API workflows
- [Wendler 5/3/1 BBB Configuration](./wendler-531-bbb.md) - Weekly lookup example
- [Sheiko Beginner Configuration](./sheiko-beginner.md) - Manual progression example
- [Greg Nuckols HF Program Specification](../../programs/020-greg-nuckols-frequency.md) - Full program algorithm

# Bill Starr 5x5 Program Configuration

This document provides the complete API configuration for the Bill Starr 5x5 program (Madcow Linear variant), demonstrating ramping sets and Heavy/Light/Medium daily intensity variation.

## Table of Contents

- [Overview](#overview)
- [Program Characteristics](#program-characteristics)
- [Configuration Components](#configuration-components)
  - [Step 1: Create Lifts](#step-1-create-lifts)
  - [Step 2: Create Daily Lookup (H/L/M Intensity)](#step-2-create-daily-lookup-hlm-intensity)
  - [Step 3: Create Prescriptions (Ramping Sets)](#step-3-create-prescriptions-ramping-sets)
  - [Step 4: Create Days](#step-4-create-days)
  - [Step 5: Create Cycle](#step-5-create-cycle)
  - [Step 6: Create Week](#step-6-create-week)
  - [Step 7: Create Progressions](#step-7-create-progressions)
  - [Step 8: Create Program](#step-8-create-program)
  - [Step 9: Link Progressions to Program](#step-9-link-progressions-to-program)
- [User Setup](#user-setup)
- [Example Workout Output](#example-workout-output)
  - [Heavy Day (Monday)](#heavy-day-monday)
  - [Light Day (Wednesday)](#light-day-wednesday)
  - [Medium Day (Friday)](#medium-day-friday)
- [Progression Behavior](#progression-behavior)
- [Complete Configuration Summary](#complete-configuration-summary)

---

## Overview

Bill Starr 5x5 is a classic intermediate strength program featuring ramping sets and a Heavy/Light/Medium weekly structure. This implementation follows the "Madcow" linear variant.

**Key Characteristics:**
- Ramping sets: weight increases across sets (not all sets at same weight)
- Heavy/Light/Medium daily variation: Monday 100%, Wednesday 80%, Friday 90%
- 3 training days per week
- Linear progression: add weight weekly to top set
- 1-week cycle that repeats

**What This Configuration Demonstrates:**
- **RAMP set scheme** - Different weights per set
- **Daily lookup tables** - Intensity modifiers per day
- **Weekly linear progression** - Weight increases each week

---

## Program Characteristics

| Property | Value |
|----------|-------|
| Cycle Length | 1 week |
| Days per Week | 3 (M/W/F) |
| Progression Type | LINEAR (weekly) |
| Progression Trigger | AFTER_WEEK |
| Weekly Increment | 5 lbs (all lifts) |

### Program Structure

**Weekly Structure:**
- Monday (Heavy Day): 100% intensity, 5x5 ramping
- Wednesday (Light Day): 80% intensity, 4x5 ramping (capped)
- Friday (Medium Day): 90% intensity, 6 sets (PR attempt)

### Ramping Set Percentages

**Standard 5x5 Ramp (Heavy Day):**
| Set | % of Top Set | Reps |
|-----|--------------|------|
| 1 | 50% | 5 |
| 2 | 63% | 5 |
| 3 | 75% | 5 |
| 4 | 88% | 5 |
| 5 | 100% | 5 |

**Light Day 4x5 (Wednesday):**
| Set | % of Top Set | Reps |
|-----|--------------|------|
| 1 | 50% | 5 |
| 2 | 63% | 5 |
| 3 | 75% | 5 |
| 4 | 75% | 5 |

*Note: Light day caps at 75% and repeats the final weight.*

**Medium Day Protocol (Friday):**
| Set | % of Top Set | Reps |
|-----|--------------|------|
| 1 | 50% | 5 |
| 2 | 63% | 5 |
| 3 | 75% | 5 |
| 4 | 88% | 5 |
| 5 | 100% | 3 |
| 6 | 75% | 8 |

*Note: Set 5 is a triple at PR weight; Set 6 is a back-off set.*

---

## Configuration Components

All API calls below assume:
- Base URL: `http://localhost:8080`
- Admin authentication: `X-Admin: true`
- User ID header: `X-User-ID: {userId}`

### Step 1: Create Lifts

Create the main lifts used in Bill Starr 5x5.

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

# Bent Over Row
curl -X POST "http://localhost:8080/lifts" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bent Over Row",
    "slug": "bent-over-row",
    "isCompetitionLift": false
  }'

# Incline Press (for Light Day pressing)
curl -X POST "http://localhost:8080/lifts" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Incline Press",
    "slug": "incline-press",
    "isCompetitionLift": false
  }'

# Deadlift (for Light Day pulling)
curl -X POST "http://localhost:8080/lifts" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Deadlift",
    "slug": "deadlift",
    "isCompetitionLift": true
  }'
```

**Lift IDs for subsequent steps** (example UUIDs):
| Lift | ID |
|------|-----|
| Squat | `11111111-1111-1111-1111-111111111111` |
| Bench Press | `22222222-2222-2222-2222-222222222222` |
| Bent Over Row | `33333333-3333-3333-3333-333333333333` |
| Incline Press | `44444444-4444-4444-4444-444444444444` |
| Deadlift | `55555555-5555-5555-5555-555555555555` |

---

### Step 2: Create Daily Lookup (H/L/M Intensity)

The daily lookup table defines the Heavy/Light/Medium intensity structure. Each day identifier maps to a percentage modifier applied to all prescriptions on that day.

```bash
curl -X POST "http://localhost:8080/daily-lookups" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bill Starr H/L/M Daily Intensity",
    "entries": [
      {
        "dayIdentifier": "heavy",
        "percentageModifier": 100.0,
        "intensityLevel": "HEAVY"
      },
      {
        "dayIdentifier": "light",
        "percentageModifier": 80.0,
        "intensityLevel": "LIGHT"
      },
      {
        "dayIdentifier": "medium",
        "percentageModifier": 90.0,
        "intensityLevel": "MEDIUM"
      }
    ]
  }'
```

**Response:**
```json
{
  "data": {
    "id": "dlookup-hlm-uuid",
    "name": "Bill Starr H/L/M Daily Intensity",
    "entries": [
      {
        "dayIdentifier": "heavy",
        "percentageModifier": 100.0,
        "intensityLevel": "HEAVY"
      },
      {
        "dayIdentifier": "light",
        "percentageModifier": 80.0,
        "intensityLevel": "LIGHT"
      },
      {
        "dayIdentifier": "medium",
        "percentageModifier": 90.0,
        "intensityLevel": "MEDIUM"
      }
    ],
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

**Daily Lookup ID:** `dddd1111-1111-1111-1111-111111111111`

**How Daily Lookups Work:**

When generating a workout, the system:
1. Looks up the day's identifier (e.g., "heavy", "light", "medium")
2. Finds the matching entry in the daily lookup table
3. Multiplies all prescription weights by the `percentageModifier`

For example, if a prescription calculates a top set weight of 300 lbs:
- Heavy day (100%): 300 lbs
- Light day (80%): 240 lbs
- Medium day (90%): 270 lbs

---

### Step 3: Create Prescriptions (Ramping Sets)

Bill Starr uses ramping sets where weight increases across sets. We use the **RAMP** set scheme.

**Heavy Day 5x5 Ramping (Squat, Bench, Row):**
```bash
# Squat 5x5 Ramping
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
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "RAMP",
      "steps": [
        {"percentage": 50, "reps": 5},
        {"percentage": 63, "reps": 5},
        {"percentage": 75, "reps": 5},
        {"percentage": 88, "reps": 5},
        {"percentage": 100, "reps": 5}
      ],
      "workSetThreshold": 80
    },
    "order": 1,
    "notes": "Ramping to top set. Add weight each week.",
    "restSeconds": 180
  }'
```

**Response:**
```json
{
  "data": {
    "id": "rx-squat-5x5-ramp-uuid",
    "liftId": "11111111-1111-1111-1111-111111111111",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "TRAINING_MAX",
      "percentage": 100.0,
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "RAMP",
      "steps": [
        {"percentage": 50, "reps": 5},
        {"percentage": 63, "reps": 5},
        {"percentage": 75, "reps": 5},
        {"percentage": 88, "reps": 5},
        {"percentage": 100, "reps": 5}
      ],
      "workSetThreshold": 80
    },
    "order": 1,
    "notes": "Ramping to top set. Add weight each week.",
    "restSeconds": 180,
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

```bash
# Bench Press 5x5 Ramping
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
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "RAMP",
      "steps": [
        {"percentage": 50, "reps": 5},
        {"percentage": 63, "reps": 5},
        {"percentage": 75, "reps": 5},
        {"percentage": 88, "reps": 5},
        {"percentage": 100, "reps": 5}
      ],
      "workSetThreshold": 80
    },
    "order": 2,
    "notes": "Ramping to top set",
    "restSeconds": 180
  }'

# Bent Over Row 5x5 Ramping
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "33333333-3333-3333-3333-333333333333",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "TRAINING_MAX",
      "percentage": 100.0,
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "RAMP",
      "steps": [
        {"percentage": 50, "reps": 5},
        {"percentage": 63, "reps": 5},
        {"percentage": 75, "reps": 5},
        {"percentage": 88, "reps": 5},
        {"percentage": 100, "reps": 5}
      ],
      "workSetThreshold": 80
    },
    "order": 3,
    "notes": "Ramping to top set",
    "restSeconds": 180
  }'
```

**Light Day 4x5 Capped Ramp (Squat only):**
```bash
# Squat 4x5 Light Day (capped at 75%)
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
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "RAMP",
      "steps": [
        {"percentage": 50, "reps": 5},
        {"percentage": 63, "reps": 5},
        {"percentage": 75, "reps": 5},
        {"percentage": 75, "reps": 5}
      ],
      "workSetThreshold": 80
    },
    "order": 1,
    "notes": "Light day - capped at 75% of Monday top set",
    "restSeconds": 120
  }'
```

**Light Day 4x5 Ramp (Incline Press, Deadlift):**
```bash
# Incline Press 4x5 (starts at 63%)
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
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "RAMP",
      "steps": [
        {"percentage": 63, "reps": 5},
        {"percentage": 75, "reps": 5},
        {"percentage": 88, "reps": 5},
        {"percentage": 100, "reps": 5}
      ],
      "workSetThreshold": 80
    },
    "order": 2,
    "notes": "Wednesday pressing movement",
    "restSeconds": 120
  }'

# Deadlift 4x5 (starts at 63%)
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
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "RAMP",
      "steps": [
        {"percentage": 63, "reps": 5},
        {"percentage": 75, "reps": 5},
        {"percentage": 88, "reps": 5},
        {"percentage": 100, "reps": 5}
      ],
      "workSetThreshold": 80
    },
    "order": 3,
    "notes": "Wednesday pulling movement",
    "restSeconds": 180
  }'
```

**Medium Day 6-Set Protocol (Friday):**
```bash
# Squat Medium Day (PR attempt + back-off)
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
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "RAMP",
      "steps": [
        {"percentage": 50, "reps": 5},
        {"percentage": 63, "reps": 5},
        {"percentage": 75, "reps": 5},
        {"percentage": 88, "reps": 5},
        {"percentage": 100, "reps": 3},
        {"percentage": 75, "reps": 8}
      ],
      "workSetThreshold": 80
    },
    "order": 1,
    "notes": "PR triple at top set, then back-off set",
    "restSeconds": 180
  }'

# Bench Press Medium Day
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
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "RAMP",
      "steps": [
        {"percentage": 50, "reps": 5},
        {"percentage": 63, "reps": 5},
        {"percentage": 75, "reps": 5},
        {"percentage": 88, "reps": 5},
        {"percentage": 100, "reps": 3},
        {"percentage": 75, "reps": 8}
      ],
      "workSetThreshold": 80
    },
    "order": 2,
    "notes": "PR triple at top set, then back-off set",
    "restSeconds": 180
  }'

# Bent Over Row Medium Day
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "33333333-3333-3333-3333-333333333333",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "TRAINING_MAX",
      "percentage": 100.0,
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "RAMP",
      "steps": [
        {"percentage": 50, "reps": 5},
        {"percentage": 63, "reps": 5},
        {"percentage": 75, "reps": 5},
        {"percentage": 88, "reps": 5},
        {"percentage": 100, "reps": 3},
        {"percentage": 75, "reps": 8}
      ],
      "workSetThreshold": 80
    },
    "order": 3,
    "notes": "PR triple at top set, then back-off set",
    "restSeconds": 180
  }'
```

**Prescription IDs for subsequent steps** (example UUIDs):
| Prescription | ID |
|--------------|-----|
| Squat 5x5 Ramp (Heavy) | `aaaa1111-1111-1111-1111-111111111111` |
| Bench 5x5 Ramp (Heavy) | `aaaa2222-2222-2222-2222-222222222222` |
| Row 5x5 Ramp (Heavy) | `aaaa3333-3333-3333-3333-333333333333` |
| Squat 4x5 Light | `aaaa4444-4444-4444-4444-444444444444` |
| Incline 4x5 | `aaaa5555-5555-5555-5555-555555555555` |
| Deadlift 4x5 | `aaaa6666-6666-6666-6666-666666666666` |
| Squat 6-set Medium | `aaaa7777-7777-7777-7777-777777777777` |
| Bench 6-set Medium | `aaaa8888-8888-8888-8888-888888888888` |
| Row 6-set Medium | `aaaa9999-9999-9999-9999-999999999999` |

---

### Step 4: Create Days

Create Heavy, Light, and Medium days with their respective prescriptions.

**Monday - Heavy Day:**
```bash
curl -X POST "http://localhost:8080/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Monday - Heavy",
    "slug": "bs-monday-heavy",
    "dailyLookupIdentifier": "heavy",
    "metadata": {
      "dayType": "heavy",
      "program": "bill-starr-5x5"
    }
  }'
```

**Response:**
```json
{
  "data": {
    "id": "day-heavy-uuid",
    "name": "Monday - Heavy",
    "slug": "bs-monday-heavy",
    "dailyLookupIdentifier": "heavy",
    "metadata": {
      "dayType": "heavy",
      "program": "bill-starr-5x5"
    },
    "programId": null,
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

```bash
# Wednesday - Light Day
curl -X POST "http://localhost:8080/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Wednesday - Light",
    "slug": "bs-wednesday-light",
    "dailyLookupIdentifier": "light",
    "metadata": {
      "dayType": "light",
      "program": "bill-starr-5x5"
    }
  }'

# Friday - Medium Day
curl -X POST "http://localhost:8080/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Friday - Medium",
    "slug": "bs-friday-medium",
    "dailyLookupIdentifier": "medium",
    "metadata": {
      "dayType": "medium",
      "program": "bill-starr-5x5"
    }
  }'
```

**Day IDs:**
| Day | ID |
|-----|-----|
| Monday Heavy | `bbbb1111-1111-1111-1111-111111111111` |
| Wednesday Light | `bbbb2222-2222-2222-2222-222222222222` |
| Friday Medium | `bbbb3333-3333-3333-3333-333333333333` |

**Add prescriptions to Heavy Day (Monday):**
```bash
# Squat 5x5 Ramp
curl -X POST "http://localhost:8080/days/bbbb1111-1111-1111-1111-111111111111/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionId": "aaaa1111-1111-1111-1111-111111111111",
    "order": 1
  }'

# Bench 5x5 Ramp
curl -X POST "http://localhost:8080/days/bbbb1111-1111-1111-1111-111111111111/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionId": "aaaa2222-2222-2222-2222-222222222222",
    "order": 2
  }'

# Row 5x5 Ramp
curl -X POST "http://localhost:8080/days/bbbb1111-1111-1111-1111-111111111111/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionId": "aaaa3333-3333-3333-3333-333333333333",
    "order": 3
  }'
```

**Add prescriptions to Light Day (Wednesday):**
```bash
# Squat 4x5 Light
curl -X POST "http://localhost:8080/days/bbbb2222-2222-2222-2222-222222222222/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionId": "aaaa4444-4444-4444-4444-444444444444",
    "order": 1
  }'

# Incline Press 4x5
curl -X POST "http://localhost:8080/days/bbbb2222-2222-2222-2222-222222222222/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionId": "aaaa5555-5555-5555-5555-555555555555",
    "order": 2
  }'

# Deadlift 4x5
curl -X POST "http://localhost:8080/days/bbbb2222-2222-2222-2222-222222222222/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionId": "aaaa6666-6666-6666-6666-666666666666",
    "order": 3
  }'
```

**Add prescriptions to Medium Day (Friday):**
```bash
# Squat 6-set Medium
curl -X POST "http://localhost:8080/days/bbbb3333-3333-3333-3333-333333333333/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionId": "aaaa7777-7777-7777-7777-777777777777",
    "order": 1
  }'

# Bench 6-set Medium
curl -X POST "http://localhost:8080/days/bbbb3333-3333-3333-3333-333333333333/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionId": "aaaa8888-8888-8888-8888-888888888888",
    "order": 2
  }'

# Row 6-set Medium
curl -X POST "http://localhost:8080/days/bbbb3333-3333-3333-3333-333333333333/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionId": "aaaa9999-9999-9999-9999-999999999999",
    "order": 3
  }'
```

---

### Step 5: Create Cycle

Bill Starr 5x5 uses a 1-week cycle that repeats.

```bash
curl -X POST "http://localhost:8080/cycles" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bill Starr 5x5 1-Week Cycle",
    "lengthWeeks": 1
  }'
```

**Response:**
```json
{
  "data": {
    "id": "cccc1111-1111-1111-1111-111111111111",
    "name": "Bill Starr 5x5 1-Week Cycle",
    "lengthWeeks": 1,
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

---

### Step 6: Create Week

Create the week with Heavy/Light/Medium day structure.

```bash
curl -X POST "http://localhost:8080/weeks" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "cycleId": "cccc1111-1111-1111-1111-111111111111",
    "weekNumber": 1,
    "name": "Week 1"
  }'
```

**Response:**
```json
{
  "data": {
    "id": "week-1-uuid",
    "cycleId": "cccc1111-1111-1111-1111-111111111111",
    "weekNumber": 1,
    "name": "Week 1",
    "days": [],
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

**Week ID:** `eeee1111-1111-1111-1111-111111111111`

**Add days to week (M/W/F pattern):**
```bash
# Position 0: Monday Heavy
curl -X POST "http://localhost:8080/weeks/eeee1111-1111-1111-1111-111111111111/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "dayId": "bbbb1111-1111-1111-1111-111111111111",
    "position": 0
  }'

# Position 1: Wednesday Light
curl -X POST "http://localhost:8080/weeks/eeee1111-1111-1111-1111-111111111111/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "dayId": "bbbb2222-2222-2222-2222-222222222222",
    "position": 1
  }'

# Position 2: Friday Medium
curl -X POST "http://localhost:8080/weeks/eeee1111-1111-1111-1111-111111111111/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "dayId": "bbbb3333-3333-3333-3333-333333333333",
    "position": 2
  }'
```

---

### Step 7: Create Progressions

Bill Starr uses weekly linear progression - add weight to each lift's top set after each week.

```bash
# Linear +5lb (All Lifts)
curl -X POST "http://localhost:8080/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Linear +5lb Weekly",
    "type": "LINEAR",
    "parameters": {
      "increment": 5.0,
      "maxType": "TRAINING_MAX",
      "triggerType": "AFTER_WEEK"
    }
  }'
```

**Response:**
```json
{
  "data": {
    "id": "prog-linear-5lb-uuid",
    "name": "Linear +5lb Weekly",
    "type": "LINEAR",
    "parameters": {
      "increment": 5.0,
      "maxType": "TRAINING_MAX",
      "triggerType": "AFTER_WEEK"
    },
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

**Progression ID:** `ffff1111-1111-1111-1111-111111111111`

---

### Step 8: Create Program

Create the Bill Starr 5x5 program, linking to cycle and daily lookup.

```bash
curl -X POST "http://localhost:8080/programs" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bill Starr 5x5",
    "slug": "bill-starr-5x5",
    "description": "Classic intermediate strength program with ramping sets and Heavy/Light/Medium weekly structure. Three days per week with linear weekly progression.",
    "cycleId": "cccc1111-1111-1111-1111-111111111111",
    "dailyLookupId": "dddd1111-1111-1111-1111-111111111111",
    "defaultRounding": 5.0
  }'
```

**Response:**
```json
{
  "data": {
    "id": "program-bs5x5-uuid",
    "name": "Bill Starr 5x5",
    "slug": "bill-starr-5x5",
    "description": "Classic intermediate strength program with ramping sets and Heavy/Light/Medium weekly structure. Three days per week with linear weekly progression.",
    "cycleId": "cccc1111-1111-1111-1111-111111111111",
    "weeklyLookupId": null,
    "dailyLookupId": "dddd1111-1111-1111-1111-111111111111",
    "defaultRounding": 5.0,
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

**Program ID:** `gggg1111-1111-1111-1111-111111111111`

---

### Step 9: Link Progressions to Program

Link the linear progression to each lift in the program.

```bash
# Squat Progression
curl -X POST "http://localhost:8080/programs/gggg1111-1111-1111-1111-111111111111/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "ffff1111-1111-1111-1111-111111111111",
    "liftId": "11111111-1111-1111-1111-111111111111",
    "priority": 1
  }'

# Bench Press Progression
curl -X POST "http://localhost:8080/programs/gggg1111-1111-1111-1111-111111111111/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "ffff1111-1111-1111-1111-111111111111",
    "liftId": "22222222-2222-2222-2222-222222222222",
    "priority": 1
  }'

# Bent Over Row Progression
curl -X POST "http://localhost:8080/programs/gggg1111-1111-1111-1111-111111111111/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "ffff1111-1111-1111-1111-111111111111",
    "liftId": "33333333-3333-3333-3333-333333333333",
    "priority": 1
  }'

# Incline Press Progression
curl -X POST "http://localhost:8080/programs/gggg1111-1111-1111-1111-111111111111/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "ffff1111-1111-1111-1111-111111111111",
    "liftId": "44444444-4444-4444-4444-444444444444",
    "priority": 1
  }'

# Deadlift Progression
curl -X POST "http://localhost:8080/programs/gggg1111-1111-1111-1111-111111111111/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "ffff1111-1111-1111-1111-111111111111",
    "liftId": "55555555-5555-5555-5555-555555555555",
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

# Squat Training Max: 315 lbs (this becomes the top set for Heavy day)
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "11111111-1111-1111-1111-111111111111",
    "type": "TRAINING_MAX",
    "value": 315.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'

# Bench Press Training Max: 225 lbs
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "22222222-2222-2222-2222-222222222222",
    "type": "TRAINING_MAX",
    "value": 225.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'

# Bent Over Row Training Max: 185 lbs
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "33333333-3333-3333-3333-333333333333",
    "type": "TRAINING_MAX",
    "value": 185.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'

# Incline Press Training Max: 155 lbs
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "44444444-4444-4444-4444-444444444444",
    "type": "TRAINING_MAX",
    "value": 155.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'

# Deadlift Training Max: 365 lbs
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "55555555-5555-5555-5555-555555555555",
    "type": "TRAINING_MAX",
    "value": 365.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'
```

### Enroll User in Program

```bash
curl -X POST "http://localhost:8080/users/${USER_ID}/program" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "programId": "gggg1111-1111-1111-1111-111111111111"
  }'
```

**Response:**
```json
{
  "data": {
    "id": "enrollment-uuid",
    "userId": "user-12345678-1234-1234-1234-123456789012",
    "program": {
      "id": "gggg1111-1111-1111-1111-111111111111",
      "name": "Bill Starr 5x5",
      "slug": "bill-starr-5x5",
      "description": "Classic intermediate strength program with ramping sets and Heavy/Light/Medium weekly structure.",
      "cycleLengthWeeks": 1
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

### Heavy Day (Monday)

Training maxes: Squat 315, Bench 225, Row 185

```bash
curl -X GET "http://localhost:8080/users/${USER_ID}/workout" \
  -H "X-User-ID: ${USER_ID}"
```

**Response:**
```json
{
  "data": {
    "userId": "user-12345678-1234-1234-1234-123456789012",
    "programId": "gggg1111-1111-1111-1111-111111111111",
    "cycleIteration": 1,
    "weekNumber": 1,
    "daySlug": "bs-monday-heavy",
    "intensityLevel": "HEAVY",
    "dailyModifier": 100.0,
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
          {"setNumber": 1, "weight": 155.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 2, "weight": 200.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 3, "weight": 235.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 4, "weight": 275.0, "targetReps": 5, "isWorkSet": true},
          {"setNumber": 5, "weight": 315.0, "targetReps": 5, "isWorkSet": true}
        ],
        "notes": "Ramping to top set. Add weight each week.",
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
          {"setNumber": 1, "weight": 110.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 2, "weight": 140.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 3, "weight": 170.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 4, "weight": 200.0, "targetReps": 5, "isWorkSet": true},
          {"setNumber": 5, "weight": 225.0, "targetReps": 5, "isWorkSet": true}
        ],
        "notes": "Ramping to top set",
        "restSeconds": 180
      },
      {
        "prescriptionId": "aaaa3333-3333-3333-3333-333333333333",
        "lift": {
          "id": "33333333-3333-3333-3333-333333333333",
          "name": "Bent Over Row",
          "slug": "bent-over-row"
        },
        "sets": [
          {"setNumber": 1, "weight": 90.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 2, "weight": 115.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 3, "weight": 140.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 4, "weight": 160.0, "targetReps": 5, "isWorkSet": true},
          {"setNumber": 5, "weight": 185.0, "targetReps": 5, "isWorkSet": true}
        ],
        "notes": "Ramping to top set",
        "restSeconds": 180
      }
    ]
  }
}
```

**Heavy Day Weight Calculation (Squat example):**
| Set | Percentage | Calculation | Weight |
|-----|------------|-------------|--------|
| 1 | 50% | 315 × 0.50 × 1.00 (heavy) | 155 lbs |
| 2 | 63% | 315 × 0.63 × 1.00 | 200 lbs |
| 3 | 75% | 315 × 0.75 × 1.00 | 235 lbs |
| 4 | 88% | 315 × 0.88 × 1.00 | 275 lbs |
| 5 | 100% | 315 × 1.00 × 1.00 | 315 lbs |

---

### Light Day (Wednesday)

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
    "programId": "gggg1111-1111-1111-1111-111111111111",
    "cycleIteration": 1,
    "weekNumber": 1,
    "daySlug": "bs-wednesday-light",
    "intensityLevel": "LIGHT",
    "dailyModifier": 80.0,
    "date": "2024-01-17",
    "exercises": [
      {
        "prescriptionId": "aaaa4444-4444-4444-4444-444444444444",
        "lift": {
          "id": "11111111-1111-1111-1111-111111111111",
          "name": "Squat",
          "slug": "squat"
        },
        "sets": [
          {"setNumber": 1, "weight": 125.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 2, "weight": 160.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 3, "weight": 190.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 4, "weight": 190.0, "targetReps": 5, "isWorkSet": false}
        ],
        "notes": "Light day - capped at 75% of Monday top set",
        "restSeconds": 120
      },
      {
        "prescriptionId": "aaaa5555-5555-5555-5555-555555555555",
        "lift": {
          "id": "44444444-4444-4444-4444-444444444444",
          "name": "Incline Press",
          "slug": "incline-press"
        },
        "sets": [
          {"setNumber": 1, "weight": 80.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 2, "weight": 95.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 3, "weight": 110.0, "targetReps": 5, "isWorkSet": true},
          {"setNumber": 4, "weight": 125.0, "targetReps": 5, "isWorkSet": true}
        ],
        "notes": "Wednesday pressing movement",
        "restSeconds": 120
      },
      {
        "prescriptionId": "aaaa6666-6666-6666-6666-666666666666",
        "lift": {
          "id": "55555555-5555-5555-5555-555555555555",
          "name": "Deadlift",
          "slug": "deadlift"
        },
        "sets": [
          {"setNumber": 1, "weight": 185.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 2, "weight": 220.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 3, "weight": 255.0, "targetReps": 5, "isWorkSet": true},
          {"setNumber": 4, "weight": 290.0, "targetReps": 5, "isWorkSet": true}
        ],
        "notes": "Wednesday pulling movement",
        "restSeconds": 180
      }
    ]
  }
}
```

**Light Day Weight Calculation (Squat example):**

The daily modifier (80%) is applied to all weights:

| Set | Percentage | Calculation | Weight |
|-----|------------|-------------|--------|
| 1 | 50% | 315 × 0.50 × 0.80 (light) | 125 lbs |
| 2 | 63% | 315 × 0.63 × 0.80 | 160 lbs |
| 3 | 75% | 315 × 0.75 × 0.80 | 190 lbs |
| 4 | 75% | 315 × 0.75 × 0.80 | 190 lbs |

The light day squat tops out at 190 lbs (vs 315 on heavy day) - a 40% reduction in top set weight, providing recovery between heavy days.

---

### Medium Day (Friday)

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
    "programId": "gggg1111-1111-1111-1111-111111111111",
    "cycleIteration": 1,
    "weekNumber": 1,
    "daySlug": "bs-friday-medium",
    "intensityLevel": "MEDIUM",
    "dailyModifier": 90.0,
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
          {"setNumber": 1, "weight": 140.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 2, "weight": 180.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 3, "weight": 215.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 4, "weight": 250.0, "targetReps": 5, "isWorkSet": true},
          {"setNumber": 5, "weight": 285.0, "targetReps": 3, "isWorkSet": true},
          {"setNumber": 6, "weight": 215.0, "targetReps": 8, "isWorkSet": false}
        ],
        "notes": "PR triple at top set, then back-off set",
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
          {"setNumber": 1, "weight": 100.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 2, "weight": 125.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 3, "weight": 150.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 4, "weight": 180.0, "targetReps": 5, "isWorkSet": true},
          {"setNumber": 5, "weight": 205.0, "targetReps": 3, "isWorkSet": true},
          {"setNumber": 6, "weight": 150.0, "targetReps": 8, "isWorkSet": false}
        ],
        "notes": "PR triple at top set, then back-off set",
        "restSeconds": 180
      },
      {
        "prescriptionId": "aaaa9999-9999-9999-9999-999999999999",
        "lift": {
          "id": "33333333-3333-3333-3333-333333333333",
          "name": "Bent Over Row",
          "slug": "bent-over-row"
        },
        "sets": [
          {"setNumber": 1, "weight": 85.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 2, "weight": 105.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 3, "weight": 125.0, "targetReps": 5, "isWorkSet": false},
          {"setNumber": 4, "weight": 145.0, "targetReps": 5, "isWorkSet": true},
          {"setNumber": 5, "weight": 165.0, "targetReps": 3, "isWorkSet": true},
          {"setNumber": 6, "weight": 125.0, "targetReps": 8, "isWorkSet": false}
        ],
        "notes": "PR triple at top set, then back-off set",
        "restSeconds": 180
      }
    ]
  }
}
```

**Medium Day Weight Calculation (Squat example):**

The daily modifier (90%) is applied:

| Set | Percentage | Calculation | Weight | Reps |
|-----|------------|-------------|--------|------|
| 1 | 50% | 315 × 0.50 × 0.90 | 140 lbs | 5 |
| 2 | 63% | 315 × 0.63 × 0.90 | 180 lbs | 5 |
| 3 | 75% | 315 × 0.75 × 0.90 | 215 lbs | 5 |
| 4 | 88% | 315 × 0.88 × 0.90 | 250 lbs | 5 |
| 5 | 100% | 315 × 1.00 × 0.90 | 285 lbs | 3 (PR) |
| 6 | 75% | 315 × 0.75 × 0.90 | 215 lbs | 8 (back-off) |

---

## Progression Behavior

Bill Starr uses weekly linear progression - weights increase after each full week.

### Triggering Progression After Week 1

After completing Friday (end of week), trigger progressions:

```bash
# Trigger all progressions for the week
curl -X POST "http://localhost:8080/users/${USER_ID}/progressions/trigger" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "ffff1111-1111-1111-1111-111111111111",
    "liftId": "11111111-1111-1111-1111-111111111111"
  }'
```

**Response:**
```json
{
  "data": {
    "results": [
      {
        "progressionId": "ffff1111-1111-1111-1111-111111111111",
        "liftId": "11111111-1111-1111-1111-111111111111",
        "applied": true,
        "skipped": false,
        "result": {
          "previousValue": 315.0,
          "newValue": 320.0,
          "delta": 5.0,
          "maxType": "TRAINING_MAX",
          "appliedAt": "2024-01-19T18:00:00Z"
        }
      }
    ],
    "totalApplied": 1,
    "totalSkipped": 0,
    "totalErrors": 0
  }
}
```

Repeat for other lifts (Bench, Row, Incline, Deadlift).

### Progression Over Multiple Weeks

Here's how the training maxes progress:

| Week | Squat TM | Bench TM | Row TM | Incline TM | Deadlift TM |
|------|----------|----------|--------|------------|-------------|
| 1 | 315 | 225 | 185 | 155 | 365 |
| 2 | 320 | 230 | 190 | 160 | 370 |
| 3 | 325 | 235 | 195 | 165 | 375 |
| 4 | 330 | 240 | 200 | 170 | 380 |

**Weekly Summary:**
- All lifts: +5 lbs per week
- Over 4 weeks: +20 lbs per lift

### Week 2 Heavy Day Workout (After Progression)

After advancing to week 2:

```bash
curl -X POST "http://localhost:8080/users/${USER_ID}/program-state/advance" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{"advanceType": "week"}'
```

The Heavy Day squat now shows:

```json
{
  "sets": [
    {"setNumber": 1, "weight": 160.0, "targetReps": 5, "isWorkSet": false},
    {"setNumber": 2, "weight": 200.0, "targetReps": 5, "isWorkSet": false},
    {"setNumber": 3, "weight": 240.0, "targetReps": 5, "isWorkSet": false},
    {"setNumber": 4, "weight": 280.0, "targetReps": 5, "isWorkSet": true},
    {"setNumber": 5, "weight": 320.0, "targetReps": 5, "isWorkSet": true}
  ]
}
```

Top set increased from 315 to 320 lbs.

---

## Complete Configuration Summary

### Entity Relationship Diagram

```
Program: Bill Starr 5x5
├── Cycle: 1-week cycle
│   └── Week 1
│       ├── Position 0: Monday Heavy (Squat/Bench/Row 5x5 Ramp)
│       ├── Position 1: Wednesday Light (Squat 4x5/Incline/Deadlift)
│       └── Position 2: Friday Medium (Squat/Bench/Row 6-set)
│
├── Daily Lookup: H/L/M Intensity
│   ├── "heavy" → 100%
│   ├── "light" → 80%
│   └── "medium" → 90%
│
├── Days
│   ├── Monday Heavy (dailyLookupIdentifier: "heavy")
│   │   ├── Prescription: Squat 5x5 RAMP @ 100% TM
│   │   ├── Prescription: Bench 5x5 RAMP @ 100% TM
│   │   └── Prescription: Row 5x5 RAMP @ 100% TM
│   │
│   ├── Wednesday Light (dailyLookupIdentifier: "light")
│   │   ├── Prescription: Squat 4x5 RAMP (capped 75%)
│   │   ├── Prescription: Incline 4x5 RAMP @ 100% TM
│   │   └── Prescription: Deadlift 4x5 RAMP @ 100% TM
│   │
│   └── Friday Medium (dailyLookupIdentifier: "medium")
│       ├── Prescription: Squat 6-set RAMP (PR triple + back-off)
│       ├── Prescription: Bench 6-set RAMP (PR triple + back-off)
│       └── Prescription: Row 6-set RAMP (PR triple + back-off)
│
└── Program Progressions
    ├── Squat → Linear +5lb (AFTER_WEEK)
    ├── Bench → Linear +5lb (AFTER_WEEK)
    ├── Row → Linear +5lb (AFTER_WEEK)
    ├── Incline → Linear +5lb (AFTER_WEEK)
    └── Deadlift → Linear +5lb (AFTER_WEEK)
```

### JSON Configuration Export

```json
{
  "program": {
    "name": "Bill Starr 5x5",
    "slug": "bill-starr-5x5",
    "description": "Classic intermediate strength program with ramping sets and Heavy/Light/Medium weekly structure.",
    "defaultRounding": 5.0
  },
  "cycle": {
    "name": "Bill Starr 5x5 1-Week Cycle",
    "lengthWeeks": 1
  },
  "dailyLookup": {
    "name": "Bill Starr H/L/M Daily Intensity",
    "entries": [
      {"dayIdentifier": "heavy", "percentageModifier": 100.0, "intensityLevel": "HEAVY"},
      {"dayIdentifier": "light", "percentageModifier": 80.0, "intensityLevel": "LIGHT"},
      {"dayIdentifier": "medium", "percentageModifier": 90.0, "intensityLevel": "MEDIUM"}
    ]
  },
  "weeks": [
    {
      "weekNumber": 1,
      "name": "Week 1",
      "days": [
        {"position": 0, "daySlug": "bs-monday-heavy"},
        {"position": 1, "daySlug": "bs-wednesday-light"},
        {"position": 2, "daySlug": "bs-friday-medium"}
      ]
    }
  ],
  "days": [
    {
      "name": "Monday - Heavy",
      "slug": "bs-monday-heavy",
      "dailyLookupIdentifier": "heavy",
      "prescriptions": [
        {
          "liftSlug": "squat",
          "setScheme": {
            "type": "RAMP",
            "steps": [
              {"percentage": 50, "reps": 5},
              {"percentage": 63, "reps": 5},
              {"percentage": 75, "reps": 5},
              {"percentage": 88, "reps": 5},
              {"percentage": 100, "reps": 5}
            ]
          },
          "order": 1
        },
        {
          "liftSlug": "bench-press",
          "setScheme": {"type": "RAMP", "steps": "...same as squat"},
          "order": 2
        },
        {
          "liftSlug": "bent-over-row",
          "setScheme": {"type": "RAMP", "steps": "...same as squat"},
          "order": 3
        }
      ]
    },
    {
      "name": "Wednesday - Light",
      "slug": "bs-wednesday-light",
      "dailyLookupIdentifier": "light",
      "prescriptions": [
        {
          "liftSlug": "squat",
          "setScheme": {
            "type": "RAMP",
            "steps": [
              {"percentage": 50, "reps": 5},
              {"percentage": 63, "reps": 5},
              {"percentage": 75, "reps": 5},
              {"percentage": 75, "reps": 5}
            ]
          },
          "order": 1
        },
        {
          "liftSlug": "incline-press",
          "setScheme": {
            "type": "RAMP",
            "steps": [
              {"percentage": 63, "reps": 5},
              {"percentage": 75, "reps": 5},
              {"percentage": 88, "reps": 5},
              {"percentage": 100, "reps": 5}
            ]
          },
          "order": 2
        },
        {
          "liftSlug": "deadlift",
          "setScheme": {"type": "RAMP", "steps": "...same as incline"},
          "order": 3
        }
      ]
    },
    {
      "name": "Friday - Medium",
      "slug": "bs-friday-medium",
      "dailyLookupIdentifier": "medium",
      "prescriptions": [
        {
          "liftSlug": "squat",
          "setScheme": {
            "type": "RAMP",
            "steps": [
              {"percentage": 50, "reps": 5},
              {"percentage": 63, "reps": 5},
              {"percentage": 75, "reps": 5},
              {"percentage": 88, "reps": 5},
              {"percentage": 100, "reps": 3},
              {"percentage": 75, "reps": 8}
            ]
          },
          "order": 1
        },
        {
          "liftSlug": "bench-press",
          "setScheme": {"type": "RAMP", "steps": "...same as squat medium"},
          "order": 2
        },
        {
          "liftSlug": "bent-over-row",
          "setScheme": {"type": "RAMP", "steps": "...same as squat medium"},
          "order": 3
        }
      ]
    }
  ],
  "progressions": [
    {
      "name": "Linear +5lb Weekly",
      "type": "LINEAR",
      "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_WEEK"},
      "liftSlugs": ["squat", "bench-press", "bent-over-row", "incline-press", "deadlift"]
    }
  ]
}
```

---

## See Also

- [API Reference](../api-reference.md) - Complete endpoint documentation
- [Workflows](../workflows.md) - Multi-step API workflows
- [Starting Strength Configuration](./starting-strength.md) - Simpler linear progression example
- [Bill Starr 5x5 Program Specification](../../programs/015-bill-starr-5x5.md) - Full program algorithm

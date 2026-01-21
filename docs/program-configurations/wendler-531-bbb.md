# Wendler 5/3/1 BBB Program Configuration

This document provides the complete API configuration for the Wendler 5/3/1 Boring But Big (BBB) program, demonstrating weekly percentage variation via lookup tables and cycle-based progression.

## Table of Contents

- [Overview](#overview)
- [Program Characteristics](#program-characteristics)
- [Configuration Components](#configuration-components)
  - [Step 1: Create Lifts](#step-1-create-lifts)
  - [Step 2: Create Weekly Lookup (Percentage Tables)](#step-2-create-weekly-lookup-percentage-tables)
  - [Step 3: Create Prescriptions](#step-3-create-prescriptions)
  - [Step 4: Create Days](#step-4-create-days)
  - [Step 5: Create Cycle](#step-5-create-cycle)
  - [Step 6: Create Weeks](#step-6-create-weeks)
  - [Step 7: Create Progressions](#step-7-create-progressions)
  - [Step 8: Create Program](#step-8-create-program)
  - [Step 9: Link Progressions to Program](#step-9-link-progressions-to-program)
- [User Setup](#user-setup)
- [Example Workout Output](#example-workout-output)
  - [Week 1 (5s Week)](#week-1-5s-week)
  - [Week 2 (3s Week)](#week-2-3s-week)
  - [Week 3 (5/3/1 Week)](#week-3-531-week)
  - [Week 4 (Deload)](#week-4-deload)
- [Cycle Progression Behavior](#cycle-progression-behavior)
- [Complete Configuration Summary](#complete-configuration-summary)

---

## Overview

Wendler 5/3/1 is Jim Wendler's percentage-based strength training program. The BBB (Boring But Big) variant adds high-volume supplemental work to build muscle mass alongside strength gains.

**Key Characteristics:**
- 4-week mesocycle with distinct rep schemes each week
- Percentage-based loading via weekly lookup tables
- AMRAP (As Many Reps As Possible) final sets
- 4 training days per week (one main lift per day)
- Progression at end of each 4-week cycle (not per session)
- BBB supplemental work: 5x10 at 50% of training max

**What This Configuration Demonstrates:**
- **Weekly lookup tables** - Different percentages for each week of the cycle
- **AMRAP final sets** - Last work set is "as many reps as possible"
- **Cycle-end progression** - Training max increases after completing all 4 weeks
- **Different increments** - Upper body +5 lbs, lower body +10 lbs

---

## Program Characteristics

| Property | Value |
|----------|-------|
| Cycle Length | 4 weeks |
| Days per Week | 4 (Squat/Bench/Deadlift/Press) |
| Progression Type | LINEAR (cycle-end) |
| Progression Trigger | AFTER_CYCLE |
| Lower Body Increment | 10 lbs (squat, deadlift) |
| Upper Body Increment | 5 lbs (bench, press) |

### Weekly Rep Scheme Structure

| Week | Name | Set 1 | Set 2 | Set 3 | Notes |
|------|------|-------|-------|-------|-------|
| 1 | 5s Week | 65%x5 | 75%x5 | 85%x5+ | AMRAP final set |
| 2 | 3s Week | 70%x3 | 80%x3 | 90%x3+ | AMRAP final set |
| 3 | 5/3/1 Week | 75%x5 | 85%x3 | 95%x1+ | AMRAP final set |
| 4 | Deload | 40%x5 | 50%x5 | 60%x5 | Recovery week |

### BBB Supplemental Work

After main lift work, perform:
- 5 sets of 10 reps at 50% of training max
- Same exercise as main lift

---

## Configuration Components

All API calls below assume:
- Base URL: `http://localhost:8080`
- Admin authentication: `X-Admin: true`
- User ID header: `X-User-ID: {userId}`

### Step 1: Create Lifts

Create the four main lifts used in Wendler 5/3/1.

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
```

**Lift IDs for subsequent steps** (example UUIDs):
| Lift | ID |
|------|-----|
| Squat | `11111111-1111-1111-1111-111111111111` |
| Bench Press | `22222222-2222-2222-2222-222222222222` |
| Deadlift | `33333333-3333-3333-3333-333333333333` |
| Overhead Press | `44444444-4444-4444-4444-444444444444` |

---

### Step 2: Create Weekly Lookup (Percentage Tables)

The weekly lookup table is the key to Wendler 5/3/1. It maps the current week number and set number to the appropriate percentage and rep count.

```bash
curl -X POST "http://localhost:8080/weekly-lookups" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Wendler 5/3/1 Main Lift Percentages",
    "entries": [
      {
        "weekNumber": 1,
        "setNumber": 1,
        "percentage": 65.0,
        "reps": 5,
        "isAmrap": false
      },
      {
        "weekNumber": 1,
        "setNumber": 2,
        "percentage": 75.0,
        "reps": 5,
        "isAmrap": false
      },
      {
        "weekNumber": 1,
        "setNumber": 3,
        "percentage": 85.0,
        "reps": 5,
        "isAmrap": true
      },
      {
        "weekNumber": 2,
        "setNumber": 1,
        "percentage": 70.0,
        "reps": 3,
        "isAmrap": false
      },
      {
        "weekNumber": 2,
        "setNumber": 2,
        "percentage": 80.0,
        "reps": 3,
        "isAmrap": false
      },
      {
        "weekNumber": 2,
        "setNumber": 3,
        "percentage": 90.0,
        "reps": 3,
        "isAmrap": true
      },
      {
        "weekNumber": 3,
        "setNumber": 1,
        "percentage": 75.0,
        "reps": 5,
        "isAmrap": false
      },
      {
        "weekNumber": 3,
        "setNumber": 2,
        "percentage": 85.0,
        "reps": 3,
        "isAmrap": false
      },
      {
        "weekNumber": 3,
        "setNumber": 3,
        "percentage": 95.0,
        "reps": 1,
        "isAmrap": true
      },
      {
        "weekNumber": 4,
        "setNumber": 1,
        "percentage": 40.0,
        "reps": 5,
        "isAmrap": false
      },
      {
        "weekNumber": 4,
        "setNumber": 2,
        "percentage": 50.0,
        "reps": 5,
        "isAmrap": false
      },
      {
        "weekNumber": 4,
        "setNumber": 3,
        "percentage": 60.0,
        "reps": 5,
        "isAmrap": false
      }
    ]
  }'
```

**Response:**
```json
{
  "data": {
    "id": "wlookup-531-uuid",
    "name": "Wendler 5/3/1 Main Lift Percentages",
    "entries": [
      {"weekNumber": 1, "setNumber": 1, "percentage": 65.0, "reps": 5, "isAmrap": false},
      {"weekNumber": 1, "setNumber": 2, "percentage": 75.0, "reps": 5, "isAmrap": false},
      {"weekNumber": 1, "setNumber": 3, "percentage": 85.0, "reps": 5, "isAmrap": true},
      {"weekNumber": 2, "setNumber": 1, "percentage": 70.0, "reps": 3, "isAmrap": false},
      {"weekNumber": 2, "setNumber": 2, "percentage": 80.0, "reps": 3, "isAmrap": false},
      {"weekNumber": 2, "setNumber": 3, "percentage": 90.0, "reps": 3, "isAmrap": true},
      {"weekNumber": 3, "setNumber": 1, "percentage": 75.0, "reps": 5, "isAmrap": false},
      {"weekNumber": 3, "setNumber": 2, "percentage": 85.0, "reps": 3, "isAmrap": false},
      {"weekNumber": 3, "setNumber": 3, "percentage": 95.0, "reps": 1, "isAmrap": true},
      {"weekNumber": 4, "setNumber": 1, "percentage": 40.0, "reps": 5, "isAmrap": false},
      {"weekNumber": 4, "setNumber": 2, "percentage": 50.0, "reps": 5, "isAmrap": false},
      {"weekNumber": 4, "setNumber": 3, "percentage": 60.0, "reps": 5, "isAmrap": false}
    ],
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

**Weekly Lookup ID:** `wlookup-1111-1111-1111-111111111111`

**How Weekly Lookups Work:**

When generating a workout, the system:
1. Determines the current week number within the cycle (1-4)
2. For each set in the prescription, looks up the matching entry
3. Uses the entry's percentage, reps, and AMRAP flag to generate the set

For example, in Week 3 with a 300 lb training max:
- Set 1: 300 × 0.75 = 225 lbs for 5 reps
- Set 2: 300 × 0.85 = 255 lbs for 3 reps
- Set 3: 300 × 0.95 = 285 lbs for 1+ reps (AMRAP)

---

### Step 3: Create Prescriptions

Wendler uses two types of prescriptions:
1. **Main lift** - 3 working sets using weekly lookup percentages
2. **BBB supplemental** - 5x10 at fixed 50% of training max

**Main Lift Prescription (uses weekly lookup):**
```bash
# Squat Main Work
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "11111111-1111-1111-1111-111111111111",
    "loadStrategy": {
      "type": "WEEKLY_LOOKUP",
      "maxType": "TRAINING_MAX",
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "LOOKUP",
      "workSets": 3
    },
    "order": 1,
    "notes": "Main work - hit your AMRAP on final set",
    "restSeconds": 180
  }'
```

**Response:**
```json
{
  "data": {
    "id": "rx-squat-main-uuid",
    "liftId": "11111111-1111-1111-1111-111111111111",
    "loadStrategy": {
      "type": "WEEKLY_LOOKUP",
      "maxType": "TRAINING_MAX",
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "LOOKUP",
      "workSets": 3
    },
    "order": 1,
    "notes": "Main work - hit your AMRAP on final set",
    "restSeconds": 180,
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

**BBB Supplemental Prescription (fixed percentage):**
```bash
# Squat BBB 5x10
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "11111111-1111-1111-1111-111111111111",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "TRAINING_MAX",
      "percentage": 50.0,
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "FIXED",
      "sets": 5,
      "reps": 10,
      "isAmrap": false
    },
    "order": 2,
    "notes": "BBB assistance - focus on volume and form",
    "restSeconds": 90
  }'
```

**Response:**
```json
{
  "data": {
    "id": "rx-squat-bbb-uuid",
    "liftId": "11111111-1111-1111-1111-111111111111",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "TRAINING_MAX",
      "percentage": 50.0,
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "FIXED",
      "sets": 5,
      "reps": 10,
      "isAmrap": false
    },
    "order": 2,
    "notes": "BBB assistance - focus on volume and form",
    "restSeconds": 90,
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

**Create prescriptions for other lifts:**
```bash
# Bench Press Main Work
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "22222222-2222-2222-2222-222222222222",
    "loadStrategy": {
      "type": "WEEKLY_LOOKUP",
      "maxType": "TRAINING_MAX",
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "LOOKUP",
      "workSets": 3
    },
    "order": 1,
    "notes": "Main work - hit your AMRAP on final set",
    "restSeconds": 180
  }'

# Bench Press BBB 5x10
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "22222222-2222-2222-2222-222222222222",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "TRAINING_MAX",
      "percentage": 50.0,
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "FIXED",
      "sets": 5,
      "reps": 10,
      "isAmrap": false
    },
    "order": 2,
    "notes": "BBB assistance - focus on volume and form",
    "restSeconds": 90
  }'

# Deadlift Main Work
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "33333333-3333-3333-3333-333333333333",
    "loadStrategy": {
      "type": "WEEKLY_LOOKUP",
      "maxType": "TRAINING_MAX",
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "LOOKUP",
      "workSets": 3
    },
    "order": 1,
    "notes": "Main work - hit your AMRAP on final set",
    "restSeconds": 180
  }'

# Deadlift BBB 5x10
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "33333333-3333-3333-3333-333333333333",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "TRAINING_MAX",
      "percentage": 50.0,
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "FIXED",
      "sets": 5,
      "reps": 10,
      "isAmrap": false
    },
    "order": 2,
    "notes": "BBB assistance - focus on volume and form",
    "restSeconds": 90
  }'

# Overhead Press Main Work
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "44444444-4444-4444-4444-444444444444",
    "loadStrategy": {
      "type": "WEEKLY_LOOKUP",
      "maxType": "TRAINING_MAX",
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "LOOKUP",
      "workSets": 3
    },
    "order": 1,
    "notes": "Main work - hit your AMRAP on final set",
    "restSeconds": 180
  }'

# Overhead Press BBB 5x10
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "44444444-4444-4444-4444-444444444444",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "TRAINING_MAX",
      "percentage": 50.0,
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "FIXED",
      "sets": 5,
      "reps": 10,
      "isAmrap": false
    },
    "order": 2,
    "notes": "BBB assistance - focus on volume and form",
    "restSeconds": 90
  }'
```

**Prescription IDs for subsequent steps** (example UUIDs):
| Prescription | ID |
|--------------|-----|
| Squat Main | `aaaa1111-1111-1111-1111-111111111111` |
| Squat BBB | `aaaa2222-2222-2222-2222-222222222222` |
| Bench Main | `aaaa3333-3333-3333-3333-333333333333` |
| Bench BBB | `aaaa4444-4444-4444-4444-444444444444` |
| Deadlift Main | `aaaa5555-5555-5555-5555-555555555555` |
| Deadlift BBB | `aaaa6666-6666-6666-6666-666666666666` |
| Press Main | `aaaa7777-7777-7777-7777-777777777777` |
| Press BBB | `aaaa8888-8888-8888-8888-888888888888` |

---

### Step 4: Create Days

Create one day for each main lift. Each day has its main lift work plus BBB assistance.

**Squat Day:**
```bash
curl -X POST "http://localhost:8080/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Squat Day",
    "slug": "531-squat-day",
    "metadata": {
      "mainLift": "squat",
      "program": "wendler-531-bbb"
    }
  }'
```

**Response:**
```json
{
  "data": {
    "id": "day-squat-uuid",
    "name": "Squat Day",
    "slug": "531-squat-day",
    "metadata": {
      "mainLift": "squat",
      "program": "wendler-531-bbb"
    },
    "programId": null,
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

```bash
# Bench Day
curl -X POST "http://localhost:8080/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bench Day",
    "slug": "531-bench-day",
    "metadata": {
      "mainLift": "bench-press",
      "program": "wendler-531-bbb"
    }
  }'

# Deadlift Day
curl -X POST "http://localhost:8080/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Deadlift Day",
    "slug": "531-deadlift-day",
    "metadata": {
      "mainLift": "deadlift",
      "program": "wendler-531-bbb"
    }
  }'

# Press Day
curl -X POST "http://localhost:8080/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Press Day",
    "slug": "531-press-day",
    "metadata": {
      "mainLift": "overhead-press",
      "program": "wendler-531-bbb"
    }
  }'
```

**Day IDs:**
| Day | ID |
|-----|-----|
| Squat Day | `bbbb1111-1111-1111-1111-111111111111` |
| Bench Day | `bbbb2222-2222-2222-2222-222222222222` |
| Deadlift Day | `bbbb3333-3333-3333-3333-333333333333` |
| Press Day | `bbbb4444-4444-4444-4444-444444444444` |

**Add prescriptions to Squat Day:**
```bash
# Squat Main
curl -X POST "http://localhost:8080/days/bbbb1111-1111-1111-1111-111111111111/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionId": "aaaa1111-1111-1111-1111-111111111111",
    "order": 1
  }'

# Squat BBB
curl -X POST "http://localhost:8080/days/bbbb1111-1111-1111-1111-111111111111/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionId": "aaaa2222-2222-2222-2222-222222222222",
    "order": 2
  }'
```

**Add prescriptions to Bench Day:**
```bash
# Bench Main
curl -X POST "http://localhost:8080/days/bbbb2222-2222-2222-2222-222222222222/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionId": "aaaa3333-3333-3333-3333-333333333333",
    "order": 1
  }'

# Bench BBB
curl -X POST "http://localhost:8080/days/bbbb2222-2222-2222-2222-222222222222/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionId": "aaaa4444-4444-4444-4444-444444444444",
    "order": 2
  }'
```

**Add prescriptions to Deadlift Day:**
```bash
# Deadlift Main
curl -X POST "http://localhost:8080/days/bbbb3333-3333-3333-3333-333333333333/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionId": "aaaa5555-5555-5555-5555-555555555555",
    "order": 1
  }'

# Deadlift BBB
curl -X POST "http://localhost:8080/days/bbbb3333-3333-3333-3333-333333333333/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionId": "aaaa6666-6666-6666-6666-666666666666",
    "order": 2
  }'
```

**Add prescriptions to Press Day:**
```bash
# Press Main
curl -X POST "http://localhost:8080/days/bbbb4444-4444-4444-4444-444444444444/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionId": "aaaa7777-7777-7777-7777-777777777777",
    "order": 1
  }'

# Press BBB
curl -X POST "http://localhost:8080/days/bbbb4444-4444-4444-4444-444444444444/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionId": "aaaa8888-8888-8888-8888-888888888888",
    "order": 2
  }'
```

---

### Step 5: Create Cycle

Wendler 5/3/1 uses a 4-week cycle.

```bash
curl -X POST "http://localhost:8080/cycles" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Wendler 5/3/1 BBB 4-Week Cycle",
    "lengthWeeks": 4
  }'
```

**Response:**
```json
{
  "data": {
    "id": "cccc1111-1111-1111-1111-111111111111",
    "name": "Wendler 5/3/1 BBB 4-Week Cycle",
    "lengthWeeks": 4,
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

---

### Step 6: Create Weeks

Create all four weeks with the same day structure. The weekly lookup handles the percentage differences.

```bash
# Week 1 (5s Week)
curl -X POST "http://localhost:8080/weeks" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "cycleId": "cccc1111-1111-1111-1111-111111111111",
    "weekNumber": 1,
    "name": "5s Week"
  }'

# Week 2 (3s Week)
curl -X POST "http://localhost:8080/weeks" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "cycleId": "cccc1111-1111-1111-1111-111111111111",
    "weekNumber": 2,
    "name": "3s Week"
  }'

# Week 3 (5/3/1 Week)
curl -X POST "http://localhost:8080/weeks" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "cycleId": "cccc1111-1111-1111-1111-111111111111",
    "weekNumber": 3,
    "name": "5/3/1 Week"
  }'

# Week 4 (Deload)
curl -X POST "http://localhost:8080/weeks" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "cycleId": "cccc1111-1111-1111-1111-111111111111",
    "weekNumber": 4,
    "name": "Deload Week"
  }'
```

**Week IDs:**
| Week | ID |
|------|-----|
| Week 1 (5s) | `dddd1111-1111-1111-1111-111111111111` |
| Week 2 (3s) | `dddd2222-2222-2222-2222-222222222222` |
| Week 3 (5/3/1) | `dddd3333-3333-3333-3333-333333333333` |
| Week 4 (Deload) | `dddd4444-4444-4444-4444-444444444444` |

**Add days to each week (same structure for all weeks):**

```bash
# Week 1 days
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

# Week 2 days (same pattern)
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

# Week 3 days (same pattern)
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

# Week 4 days (same pattern)
curl -X POST "http://localhost:8080/weeks/dddd4444-4444-4444-4444-444444444444/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "bbbb1111-1111-1111-1111-111111111111", "position": 0}'

curl -X POST "http://localhost:8080/weeks/dddd4444-4444-4444-4444-444444444444/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "bbbb2222-2222-2222-2222-222222222222", "position": 1}'

curl -X POST "http://localhost:8080/weeks/dddd4444-4444-4444-4444-444444444444/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "bbbb3333-3333-3333-3333-333333333333", "position": 2}'

curl -X POST "http://localhost:8080/weeks/dddd4444-4444-4444-4444-444444444444/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "bbbb4444-4444-4444-4444-444444444444", "position": 3}'
```

---

### Step 7: Create Progressions

Wendler uses cycle-end progression with different increments for upper and lower body lifts.

**Lower Body Progression (+10 lbs):**
```bash
curl -X POST "http://localhost:8080/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Linear +10lb Lower Body (Cycle-End)",
    "type": "LINEAR",
    "parameters": {
      "increment": 10.0,
      "maxType": "TRAINING_MAX",
      "triggerType": "AFTER_CYCLE"
    }
  }'
```

**Response:**
```json
{
  "data": {
    "id": "prog-lower-uuid",
    "name": "Linear +10lb Lower Body (Cycle-End)",
    "type": "LINEAR",
    "parameters": {
      "increment": 10.0,
      "maxType": "TRAINING_MAX",
      "triggerType": "AFTER_CYCLE"
    },
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

**Upper Body Progression (+5 lbs):**
```bash
curl -X POST "http://localhost:8080/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Linear +5lb Upper Body (Cycle-End)",
    "type": "LINEAR",
    "parameters": {
      "increment": 5.0,
      "maxType": "TRAINING_MAX",
      "triggerType": "AFTER_CYCLE"
    }
  }'
```

**Progression IDs:**
| Progression | ID |
|-------------|-----|
| Lower Body +10lb | `eeee1111-1111-1111-1111-111111111111` |
| Upper Body +5lb | `eeee2222-2222-2222-2222-222222222222` |

---

### Step 8: Create Program

Create the Wendler 5/3/1 BBB program, linking to cycle and weekly lookup.

```bash
curl -X POST "http://localhost:8080/programs" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Wendler 5/3/1 BBB",
    "slug": "wendler-531-bbb",
    "description": "Jim Wendler'\''s 5/3/1 Boring But Big program. 4-week cycles with percentage-based main lifts and 5x10 supplemental volume. Training max increases at end of each cycle.",
    "cycleId": "cccc1111-1111-1111-1111-111111111111",
    "weeklyLookupId": "wlookup-1111-1111-1111-111111111111",
    "defaultRounding": 5.0
  }'
```

**Response:**
```json
{
  "data": {
    "id": "program-531-bbb-uuid",
    "name": "Wendler 5/3/1 BBB",
    "slug": "wendler-531-bbb",
    "description": "Jim Wendler's 5/3/1 Boring But Big program. 4-week cycles with percentage-based main lifts and 5x10 supplemental volume. Training max increases at end of each cycle.",
    "cycleId": "cccc1111-1111-1111-1111-111111111111",
    "weeklyLookupId": "wlookup-1111-1111-1111-111111111111",
    "dailyLookupId": null,
    "defaultRounding": 5.0,
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

**Program ID:** `ffff1111-1111-1111-1111-111111111111`

---

### Step 9: Link Progressions to Program

Link each progression to the appropriate lifts.

```bash
# Squat Progression (Lower Body +10lb)
curl -X POST "http://localhost:8080/programs/ffff1111-1111-1111-1111-111111111111/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "eeee1111-1111-1111-1111-111111111111",
    "liftId": "11111111-1111-1111-1111-111111111111",
    "priority": 1
  }'

# Bench Press Progression (Upper Body +5lb)
curl -X POST "http://localhost:8080/programs/ffff1111-1111-1111-1111-111111111111/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "eeee2222-2222-2222-2222-222222222222",
    "liftId": "22222222-2222-2222-2222-222222222222",
    "priority": 1
  }'

# Deadlift Progression (Lower Body +10lb)
curl -X POST "http://localhost:8080/programs/ffff1111-1111-1111-1111-111111111111/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "eeee1111-1111-1111-1111-111111111111",
    "liftId": "33333333-3333-3333-3333-333333333333",
    "priority": 1
  }'

# Overhead Press Progression (Upper Body +5lb)
curl -X POST "http://localhost:8080/programs/ffff1111-1111-1111-1111-111111111111/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "eeee2222-2222-2222-2222-222222222222",
    "liftId": "44444444-4444-4444-4444-444444444444",
    "priority": 1
  }'
```

---

## User Setup

Before a user can generate workouts, they need training maxes for all lifts.

### Create User Training Maxes

The training max in 5/3/1 should be 90% of your true 1RM.

```bash
# User ID for examples
USER_ID="user-12345678-1234-1234-1234-123456789012"

# Squat Training Max: 300 lbs (based on 335 lb 1RM)
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "11111111-1111-1111-1111-111111111111",
    "type": "TRAINING_MAX",
    "value": 300.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'

# Bench Press Training Max: 200 lbs (based on 225 lb 1RM)
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "22222222-2222-2222-2222-222222222222",
    "type": "TRAINING_MAX",
    "value": 200.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'

# Deadlift Training Max: 350 lbs (based on 390 lb 1RM)
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "33333333-3333-3333-3333-333333333333",
    "type": "TRAINING_MAX",
    "value": 350.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'

# Overhead Press Training Max: 120 lbs (based on 135 lb 1RM)
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "44444444-4444-4444-4444-444444444444",
    "type": "TRAINING_MAX",
    "value": 120.0,
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
      "name": "Wendler 5/3/1 BBB",
      "slug": "wendler-531-bbb",
      "description": "Jim Wendler's 5/3/1 Boring But Big program.",
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

### Week 1 (5s Week)

Training maxes: Squat 300, Bench 200, Deadlift 350, Press 120

**Week 1, Day 1: Squat Day**

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
    "weekName": "5s Week",
    "daySlug": "531-squat-day",
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
          {"setNumber": 1, "weight": 195.0, "targetReps": 5, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 2, "weight": 225.0, "targetReps": 5, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 3, "weight": 255.0, "targetReps": 5, "isAmrap": true, "isWorkSet": true}
        ],
        "notes": "Main work - hit your AMRAP on final set",
        "restSeconds": 180
      },
      {
        "prescriptionId": "aaaa2222-2222-2222-2222-222222222222",
        "lift": {
          "id": "11111111-1111-1111-1111-111111111111",
          "name": "Squat",
          "slug": "squat"
        },
        "sets": [
          {"setNumber": 1, "weight": 150.0, "targetReps": 10, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 2, "weight": 150.0, "targetReps": 10, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 3, "weight": 150.0, "targetReps": 10, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 4, "weight": 150.0, "targetReps": 10, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 5, "weight": 150.0, "targetReps": 10, "isAmrap": false, "isWorkSet": true}
        ],
        "notes": "BBB assistance - focus on volume and form",
        "restSeconds": 90
      }
    ]
  }
}
```

**Week 1 Weight Calculation (Squat):**

Main Work (from weekly lookup, week 1):
| Set | Percentage | Calculation | Weight | Reps |
|-----|------------|-------------|--------|------|
| 1 | 65% | 300 × 0.65 = 195 | 195 lbs | 5 |
| 2 | 75% | 300 × 0.75 = 225 | 225 lbs | 5 |
| 3 | 85% | 300 × 0.85 = 255 | 255 lbs | 5+ (AMRAP) |

BBB Supplemental (fixed 50%):
| Sets | Percentage | Calculation | Weight | Reps |
|------|------------|-------------|--------|------|
| 1-5 | 50% | 300 × 0.50 = 150 | 150 lbs | 10 |

---

### Week 2 (3s Week)

After advancing through week 1:

```bash
# Advance through all days in week 1
curl -X POST "http://localhost:8080/users/${USER_ID}/program-state/advance" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{"advanceType": "week"}'
```

**Week 2, Day 1: Squat Day**

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
    "weekNumber": 2,
    "weekName": "3s Week",
    "daySlug": "531-squat-day",
    "date": "2024-01-22",
    "exercises": [
      {
        "prescriptionId": "aaaa1111-1111-1111-1111-111111111111",
        "lift": {
          "id": "11111111-1111-1111-1111-111111111111",
          "name": "Squat",
          "slug": "squat"
        },
        "sets": [
          {"setNumber": 1, "weight": 210.0, "targetReps": 3, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 2, "weight": 240.0, "targetReps": 3, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 3, "weight": 270.0, "targetReps": 3, "isAmrap": true, "isWorkSet": true}
        ],
        "notes": "Main work - hit your AMRAP on final set",
        "restSeconds": 180
      },
      {
        "prescriptionId": "aaaa2222-2222-2222-2222-222222222222",
        "lift": {
          "id": "11111111-1111-1111-1111-111111111111",
          "name": "Squat",
          "slug": "squat"
        },
        "sets": [
          {"setNumber": 1, "weight": 150.0, "targetReps": 10, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 2, "weight": 150.0, "targetReps": 10, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 3, "weight": 150.0, "targetReps": 10, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 4, "weight": 150.0, "targetReps": 10, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 5, "weight": 150.0, "targetReps": 10, "isAmrap": false, "isWorkSet": true}
        ],
        "notes": "BBB assistance - focus on volume and form",
        "restSeconds": 90
      }
    ]
  }
}
```

**Week 2 Weight Calculation (Squat):**

Main Work (from weekly lookup, week 2):
| Set | Percentage | Calculation | Weight | Reps |
|-----|------------|-------------|--------|------|
| 1 | 70% | 300 × 0.70 = 210 | 210 lbs | 3 |
| 2 | 80% | 300 × 0.80 = 240 | 240 lbs | 3 |
| 3 | 90% | 300 × 0.90 = 270 | 270 lbs | 3+ (AMRAP) |

BBB Supplemental remains at 50% (150 lbs).

---

### Week 3 (5/3/1 Week)

After advancing to week 3:

**Week 3, Day 1: Squat Day**

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
    "weekNumber": 3,
    "weekName": "5/3/1 Week",
    "daySlug": "531-squat-day",
    "date": "2024-01-29",
    "exercises": [
      {
        "prescriptionId": "aaaa1111-1111-1111-1111-111111111111",
        "lift": {
          "id": "11111111-1111-1111-1111-111111111111",
          "name": "Squat",
          "slug": "squat"
        },
        "sets": [
          {"setNumber": 1, "weight": 225.0, "targetReps": 5, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 2, "weight": 255.0, "targetReps": 3, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 3, "weight": 285.0, "targetReps": 1, "isAmrap": true, "isWorkSet": true}
        ],
        "notes": "Main work - hit your AMRAP on final set",
        "restSeconds": 180
      },
      {
        "prescriptionId": "aaaa2222-2222-2222-2222-222222222222",
        "lift": {
          "id": "11111111-1111-1111-1111-111111111111",
          "name": "Squat",
          "slug": "squat"
        },
        "sets": [
          {"setNumber": 1, "weight": 150.0, "targetReps": 10, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 2, "weight": 150.0, "targetReps": 10, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 3, "weight": 150.0, "targetReps": 10, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 4, "weight": 150.0, "targetReps": 10, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 5, "weight": 150.0, "targetReps": 10, "isAmrap": false, "isWorkSet": true}
        ],
        "notes": "BBB assistance - focus on volume and form",
        "restSeconds": 90
      }
    ]
  }
}
```

**Week 3 Weight Calculation (Squat):**

Main Work (from weekly lookup, week 3):
| Set | Percentage | Calculation | Weight | Reps |
|-----|------------|-------------|--------|------|
| 1 | 75% | 300 × 0.75 = 225 | 225 lbs | 5 |
| 2 | 85% | 300 × 0.85 = 255 | 255 lbs | 3 |
| 3 | 95% | 300 × 0.95 = 285 | 285 lbs | 1+ (AMRAP) |

Note: Week 3 is the heaviest week. The final set at 95% is the intensity peak.

---

### Week 4 (Deload)

After advancing to week 4:

**Week 4, Day 1: Squat Day**

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
    "weekNumber": 4,
    "weekName": "Deload Week",
    "daySlug": "531-squat-day",
    "date": "2024-02-05",
    "exercises": [
      {
        "prescriptionId": "aaaa1111-1111-1111-1111-111111111111",
        "lift": {
          "id": "11111111-1111-1111-1111-111111111111",
          "name": "Squat",
          "slug": "squat"
        },
        "sets": [
          {"setNumber": 1, "weight": 120.0, "targetReps": 5, "isAmrap": false, "isWorkSet": false},
          {"setNumber": 2, "weight": 150.0, "targetReps": 5, "isAmrap": false, "isWorkSet": false},
          {"setNumber": 3, "weight": 180.0, "targetReps": 5, "isAmrap": false, "isWorkSet": false}
        ],
        "notes": "Main work - hit your AMRAP on final set",
        "restSeconds": 180
      },
      {
        "prescriptionId": "aaaa2222-2222-2222-2222-222222222222",
        "lift": {
          "id": "11111111-1111-1111-1111-111111111111",
          "name": "Squat",
          "slug": "squat"
        },
        "sets": [
          {"setNumber": 1, "weight": 150.0, "targetReps": 10, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 2, "weight": 150.0, "targetReps": 10, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 3, "weight": 150.0, "targetReps": 10, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 4, "weight": 150.0, "targetReps": 10, "isAmrap": false, "isWorkSet": true},
          {"setNumber": 5, "weight": 150.0, "targetReps": 10, "isAmrap": false, "isWorkSet": true}
        ],
        "notes": "BBB assistance - focus on volume and form",
        "restSeconds": 90
      }
    ]
  }
}
```

**Week 4 Weight Calculation (Squat):**

Main Work (from weekly lookup, week 4):
| Set | Percentage | Calculation | Weight | Reps |
|-----|------------|-------------|--------|------|
| 1 | 40% | 300 × 0.40 = 120 | 120 lbs | 5 |
| 2 | 50% | 300 × 0.50 = 150 | 150 lbs | 5 |
| 3 | 60% | 300 × 0.60 = 180 | 180 lbs | 5 |

Note: Deload week has no AMRAP sets and significantly lower intensity for recovery.

---

## Cycle Progression Behavior

After completing all 4 weeks, trigger progression to increase training maxes for the next cycle.

### Triggering Cycle-End Progression

```bash
# Trigger squat progression (+10 lbs)
curl -X POST "http://localhost:8080/users/${USER_ID}/progressions/trigger" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "eeee1111-1111-1111-1111-111111111111",
    "liftId": "11111111-1111-1111-1111-111111111111"
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
        "skipped": false,
        "result": {
          "previousValue": 300.0,
          "newValue": 310.0,
          "delta": 10.0,
          "maxType": "TRAINING_MAX",
          "appliedAt": "2024-02-12T18:00:00Z"
        }
      }
    ],
    "totalApplied": 1,
    "totalSkipped": 0,
    "totalErrors": 0
  }
}
```

**Trigger all progressions:**
```bash
# Bench Press (+5 lbs)
curl -X POST "http://localhost:8080/users/${USER_ID}/progressions/trigger" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "eeee2222-2222-2222-2222-222222222222",
    "liftId": "22222222-2222-2222-2222-222222222222"
  }'

# Deadlift (+10 lbs)
curl -X POST "http://localhost:8080/users/${USER_ID}/progressions/trigger" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "eeee1111-1111-1111-1111-111111111111",
    "liftId": "33333333-3333-3333-3333-333333333333"
  }'

# Overhead Press (+5 lbs)
curl -X POST "http://localhost:8080/users/${USER_ID}/progressions/trigger" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "eeee2222-2222-2222-2222-222222222222",
    "liftId": "44444444-4444-4444-4444-444444444444"
  }'
```

### Training Max Progression Over Cycles

| Cycle | Squat TM | Bench TM | Deadlift TM | Press TM |
|-------|----------|----------|-------------|----------|
| 1 | 300 | 200 | 350 | 120 |
| 2 | 310 | 205 | 360 | 125 |
| 3 | 320 | 210 | 370 | 130 |
| 4 | 330 | 215 | 380 | 135 |

**Yearly Progression (13 cycles):**
- Squat: 300 → 420 (+120 lbs)
- Bench: 200 → 260 (+60 lbs)
- Deadlift: 350 → 470 (+120 lbs)
- Press: 120 → 180 (+60 lbs)

### Cycle 2, Week 1 Workout (After Progression)

After advancing to cycle 2:

```bash
curl -X POST "http://localhost:8080/users/${USER_ID}/program-state/advance" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{"advanceType": "cycle"}'
```

**Response (Squat Day, Week 1, Cycle 2):**
```json
{
  "data": {
    "cycleIteration": 2,
    "weekNumber": 1,
    "weekName": "5s Week",
    "daySlug": "531-squat-day",
    "exercises": [
      {
        "lift": {"name": "Squat"},
        "sets": [
          {"setNumber": 1, "weight": 200.0, "targetReps": 5, "isAmrap": false},
          {"setNumber": 2, "weight": 230.0, "targetReps": 5, "isAmrap": false},
          {"setNumber": 3, "weight": 265.0, "targetReps": 5, "isAmrap": true}
        ]
      },
      {
        "lift": {"name": "Squat"},
        "sets": [
          {"setNumber": 1, "weight": 155.0, "targetReps": 10, "isAmrap": false},
          {"setNumber": 2, "weight": 155.0, "targetReps": 10, "isAmrap": false},
          {"setNumber": 3, "weight": 155.0, "targetReps": 10, "isAmrap": false},
          {"setNumber": 4, "weight": 155.0, "targetReps": 10, "isAmrap": false},
          {"setNumber": 5, "weight": 155.0, "targetReps": 10, "isAmrap": false}
        ]
      }
    ]
  }
}
```

The training max increased from 300 to 310, resulting in:
- Week 1 Set 3: 310 × 0.85 = 263.5 → 265 lbs (rounded)
- BBB: 310 × 0.50 = 155 lbs

---

## Complete Configuration Summary

### Entity Relationship Diagram

```
Program: Wendler 5/3/1 BBB
├── Cycle: 4-week cycle
│   ├── Week 1 (5s Week)
│   │   ├── Position 0: Squat Day
│   │   ├── Position 1: Bench Day
│   │   ├── Position 2: Deadlift Day
│   │   └── Position 3: Press Day
│   ├── Week 2 (3s Week)
│   │   └── [same day structure]
│   ├── Week 3 (5/3/1 Week)
│   │   └── [same day structure]
│   └── Week 4 (Deload)
│       └── [same day structure]
│
├── Weekly Lookup: 5/3/1 Percentages
│   ├── Week 1: 65%x5, 75%x5, 85%x5+
│   ├── Week 2: 70%x3, 80%x3, 90%x3+
│   ├── Week 3: 75%x5, 85%x3, 95%x1+
│   └── Week 4: 40%x5, 50%x5, 60%x5
│
├── Days
│   ├── Squat Day
│   │   ├── Prescription: Squat Main (WEEKLY_LOOKUP)
│   │   └── Prescription: Squat BBB 5x10 @ 50%
│   ├── Bench Day
│   │   ├── Prescription: Bench Main (WEEKLY_LOOKUP)
│   │   └── Prescription: Bench BBB 5x10 @ 50%
│   ├── Deadlift Day
│   │   ├── Prescription: Deadlift Main (WEEKLY_LOOKUP)
│   │   └── Prescription: Deadlift BBB 5x10 @ 50%
│   └── Press Day
│       ├── Prescription: Press Main (WEEKLY_LOOKUP)
│       └── Prescription: Press BBB 5x10 @ 50%
│
└── Program Progressions
    ├── Squat → Linear +10lb (AFTER_CYCLE)
    ├── Bench → Linear +5lb (AFTER_CYCLE)
    ├── Deadlift → Linear +10lb (AFTER_CYCLE)
    └── Press → Linear +5lb (AFTER_CYCLE)
```

### JSON Configuration Export

```json
{
  "program": {
    "name": "Wendler 5/3/1 BBB",
    "slug": "wendler-531-bbb",
    "description": "Jim Wendler's 5/3/1 Boring But Big program. 4-week cycles with percentage-based main lifts and 5x10 supplemental volume.",
    "defaultRounding": 5.0
  },
  "cycle": {
    "name": "Wendler 5/3/1 BBB 4-Week Cycle",
    "lengthWeeks": 4
  },
  "weeklyLookup": {
    "name": "Wendler 5/3/1 Main Lift Percentages",
    "entries": [
      {"weekNumber": 1, "setNumber": 1, "percentage": 65.0, "reps": 5, "isAmrap": false},
      {"weekNumber": 1, "setNumber": 2, "percentage": 75.0, "reps": 5, "isAmrap": false},
      {"weekNumber": 1, "setNumber": 3, "percentage": 85.0, "reps": 5, "isAmrap": true},
      {"weekNumber": 2, "setNumber": 1, "percentage": 70.0, "reps": 3, "isAmrap": false},
      {"weekNumber": 2, "setNumber": 2, "percentage": 80.0, "reps": 3, "isAmrap": false},
      {"weekNumber": 2, "setNumber": 3, "percentage": 90.0, "reps": 3, "isAmrap": true},
      {"weekNumber": 3, "setNumber": 1, "percentage": 75.0, "reps": 5, "isAmrap": false},
      {"weekNumber": 3, "setNumber": 2, "percentage": 85.0, "reps": 3, "isAmrap": false},
      {"weekNumber": 3, "setNumber": 3, "percentage": 95.0, "reps": 1, "isAmrap": true},
      {"weekNumber": 4, "setNumber": 1, "percentage": 40.0, "reps": 5, "isAmrap": false},
      {"weekNumber": 4, "setNumber": 2, "percentage": 50.0, "reps": 5, "isAmrap": false},
      {"weekNumber": 4, "setNumber": 3, "percentage": 60.0, "reps": 5, "isAmrap": false}
    ]
  },
  "weeks": [
    {"weekNumber": 1, "name": "5s Week", "days": ["531-squat-day", "531-bench-day", "531-deadlift-day", "531-press-day"]},
    {"weekNumber": 2, "name": "3s Week", "days": ["531-squat-day", "531-bench-day", "531-deadlift-day", "531-press-day"]},
    {"weekNumber": 3, "name": "5/3/1 Week", "days": ["531-squat-day", "531-bench-day", "531-deadlift-day", "531-press-day"]},
    {"weekNumber": 4, "name": "Deload Week", "days": ["531-squat-day", "531-bench-day", "531-deadlift-day", "531-press-day"]}
  ],
  "days": [
    {
      "name": "Squat Day",
      "slug": "531-squat-day",
      "prescriptions": [
        {
          "liftSlug": "squat",
          "loadStrategy": {"type": "WEEKLY_LOOKUP"},
          "setScheme": {"type": "LOOKUP", "workSets": 3},
          "order": 1
        },
        {
          "liftSlug": "squat",
          "loadStrategy": {"type": "PERCENT_OF", "percentage": 50.0},
          "setScheme": {"type": "FIXED", "sets": 5, "reps": 10},
          "order": 2
        }
      ]
    },
    {
      "name": "Bench Day",
      "slug": "531-bench-day",
      "prescriptions": [
        {"liftSlug": "bench-press", "loadStrategy": {"type": "WEEKLY_LOOKUP"}, "setScheme": {"type": "LOOKUP", "workSets": 3}, "order": 1},
        {"liftSlug": "bench-press", "loadStrategy": {"type": "PERCENT_OF", "percentage": 50.0}, "setScheme": {"type": "FIXED", "sets": 5, "reps": 10}, "order": 2}
      ]
    },
    {
      "name": "Deadlift Day",
      "slug": "531-deadlift-day",
      "prescriptions": [
        {"liftSlug": "deadlift", "loadStrategy": {"type": "WEEKLY_LOOKUP"}, "setScheme": {"type": "LOOKUP", "workSets": 3}, "order": 1},
        {"liftSlug": "deadlift", "loadStrategy": {"type": "PERCENT_OF", "percentage": 50.0}, "setScheme": {"type": "FIXED", "sets": 5, "reps": 10}, "order": 2}
      ]
    },
    {
      "name": "Press Day",
      "slug": "531-press-day",
      "prescriptions": [
        {"liftSlug": "overhead-press", "loadStrategy": {"type": "WEEKLY_LOOKUP"}, "setScheme": {"type": "LOOKUP", "workSets": 3}, "order": 1},
        {"liftSlug": "overhead-press", "loadStrategy": {"type": "PERCENT_OF", "percentage": 50.0}, "setScheme": {"type": "FIXED", "sets": 5, "reps": 10}, "order": 2}
      ]
    }
  ],
  "progressions": [
    {
      "name": "Linear +10lb Lower Body (Cycle-End)",
      "type": "LINEAR",
      "parameters": {"increment": 10.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_CYCLE"},
      "liftSlugs": ["squat", "deadlift"]
    },
    {
      "name": "Linear +5lb Upper Body (Cycle-End)",
      "type": "LINEAR",
      "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_CYCLE"},
      "liftSlugs": ["bench-press", "overhead-press"]
    }
  ]
}
```

---

## See Also

- [API Reference](../api-reference.md) - Complete endpoint documentation
- [Workflows](../workflows.md) - Multi-step API workflows
- [Starting Strength Configuration](./starting-strength.md) - Simple linear progression example
- [Bill Starr 5x5 Configuration](./bill-starr-5x5.md) - Daily lookup tables example
- [Wendler 5/3/1 Program Specification](../../programs/003-wendler-531-bbb.md) - Full program algorithm

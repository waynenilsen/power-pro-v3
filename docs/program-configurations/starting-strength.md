# Starting Strength Program Configuration

This document provides the complete API configuration for the Starting Strength program (Original Novice Program variant), demonstrating linear per-session progression.

## Table of Contents

- [Overview](#overview)
- [Program Characteristics](#program-characteristics)
- [Configuration Components](#configuration-components)
  - [Step 1: Create Lifts](#step-1-create-lifts)
  - [Step 2: Create Prescriptions](#step-2-create-prescriptions)
  - [Step 3: Create Days](#step-3-create-days)
  - [Step 4: Create Cycle](#step-4-create-cycle)
  - [Step 5: Create Weeks](#step-5-create-weeks)
  - [Step 6: Create Progressions](#step-6-create-progressions)
  - [Step 7: Create Program](#step-7-create-program)
  - [Step 8: Link Progressions to Program](#step-8-link-progressions-to-program)
- [User Setup](#user-setup)
- [Example Workout Output](#example-workout-output)
- [Progression Behavior](#progression-behavior)
- [Complete Configuration Summary](#complete-configuration-summary)

---

## Overview

Starting Strength is Mark Rippetoe's novice linear progression program. It is the simplest program to configure in PowerPro and serves as a reference implementation for linear progression.

**Key Characteristics:**
- Linear progression: add weight every session
- 3 training days per week (A/B/A, B/A/B rotation)
- 3x5 for main lifts, 1x5 for deadlift, 5x3 for power clean
- No weekly lookups (constant percentages)
- No daily lookups (constant intensity within each day)

---

## Program Characteristics

| Property | Value |
|----------|-------|
| Cycle Length | 1 week |
| Days per Week | 3 (alternating A/B) |
| Progression Type | LINEAR (per-session) |
| Progression Trigger | AFTER_SESSION |
| Lower Body Increment | 5 lbs (squat), 10 lbs (deadlift) |
| Upper Body Increment | 5 lbs (bench, press, power clean) |

### Program Structure

**Week Pattern (A/B rotation):**
- Session 1: Day A (Squat, Bench, Deadlift)
- Session 2: Day B (Squat, Press, Power Clean)
- Session 3: Day A (Squat, Bench, Deadlift)
- Next week: B/A/B, then A/B/A, etc.

### Working Sets

| Exercise | Sets x Reps | Notes |
|----------|-------------|-------|
| Squat | 3x5 | Every session |
| Bench Press | 3x5 | Day A only |
| Overhead Press | 3x5 | Day B only |
| Deadlift | 1x5 | Day A only |
| Power Clean | 5x3 | Day B only |

---

## Configuration Components

All API calls below assume:
- Base URL: `http://localhost:8080`
- Admin authentication: `X-Admin: true`
- User ID header: `X-User-ID: {userId}`

### Step 1: Create Lifts

Create the five main lifts used in Starting Strength.

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
```

```bash
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

```bash
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
```

```bash
# Power Clean
curl -X POST "http://localhost:8080/lifts" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Power Clean",
    "slug": "power-clean",
    "isCompetitionLift": false
  }'
```

**Lift IDs for subsequent steps** (example UUIDs):
| Lift | ID |
|------|-----|
| Squat | `11111111-1111-1111-1111-111111111111` |
| Bench Press | `22222222-2222-2222-2222-222222222222` |
| Overhead Press | `33333333-3333-3333-3333-333333333333` |
| Deadlift | `44444444-4444-4444-4444-444444444444` |
| Power Clean | `55555555-5555-5555-5555-555555555555` |

---

### Step 2: Create Prescriptions

Create prescriptions for each exercise slot. All working sets use 100% of training max since progression happens between sessions.

**Squat 3x5:**
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
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "FIXED",
      "sets": 3,
      "reps": 5,
      "isAmrap": false
    },
    "order": 1,
    "notes": "Focus on depth - crease of hip below top of knee",
    "restSeconds": 180
  }'
```

**Response:**
```json
{
  "data": {
    "id": "rx-squat-3x5-uuid",
    "liftId": "11111111-1111-1111-1111-111111111111",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "TRAINING_MAX",
      "percentage": 100.0,
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "FIXED",
      "sets": 3,
      "reps": 5,
      "isAmrap": false
    },
    "order": 1,
    "notes": "Focus on depth - crease of hip below top of knee",
    "restSeconds": 180,
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

**Bench Press 3x5:**
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
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "FIXED",
      "sets": 3,
      "reps": 5,
      "isAmrap": false
    },
    "order": 2,
    "notes": "Touch chest, pause briefly, press",
    "restSeconds": 180
  }'
```

**Overhead Press 3x5:**
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
      "percentage": 100.0,
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "FIXED",
      "sets": 3,
      "reps": 5,
      "isAmrap": false
    },
    "order": 2,
    "notes": "Strict press - no leg drive",
    "restSeconds": 180
  }'
```

**Deadlift 1x5:**
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
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "FIXED",
      "sets": 1,
      "reps": 5,
      "isAmrap": false
    },
    "order": 3,
    "notes": "Reset each rep - no touch and go",
    "restSeconds": 180
  }'
```

**Power Clean 5x3:**
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
      "roundTo": 5.0
    },
    "setScheme": {
      "type": "FIXED",
      "sets": 5,
      "reps": 3,
      "isAmrap": false
    },
    "order": 3,
    "notes": "Explosive triple extension",
    "restSeconds": 120
  }'
```

**Prescription IDs for subsequent steps** (example UUIDs):
| Prescription | ID |
|--------------|-----|
| Squat 3x5 | `aaaa1111-1111-1111-1111-111111111111` |
| Bench 3x5 | `aaaa2222-2222-2222-2222-222222222222` |
| Press 3x5 | `aaaa3333-3333-3333-3333-333333333333` |
| Deadlift 1x5 | `aaaa4444-4444-4444-4444-444444444444` |
| Power Clean 5x3 | `aaaa5555-5555-5555-5555-555555555555` |

---

### Step 3: Create Days

Create Day A (Squat/Bench/Deadlift) and Day B (Squat/Press/Power Clean).

**Day A:**
```bash
curl -X POST "http://localhost:8080/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Day A",
    "slug": "ss-day-a",
    "metadata": {
      "focus": "squat-bench-deadlift",
      "program": "starting-strength"
    }
  }'
```

**Response:**
```json
{
  "data": {
    "id": "day-a-uuid",
    "name": "Day A",
    "slug": "ss-day-a",
    "metadata": {
      "focus": "squat-bench-deadlift",
      "program": "starting-strength"
    },
    "programId": null,
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

**Day B:**
```bash
curl -X POST "http://localhost:8080/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Day B",
    "slug": "ss-day-b",
    "metadata": {
      "focus": "squat-press-clean",
      "program": "starting-strength"
    }
  }'
```

**Day IDs:**
| Day | ID |
|-----|-----|
| Day A | `bbbb1111-1111-1111-1111-111111111111` |
| Day B | `bbbb2222-2222-2222-2222-222222222222` |

**Add prescriptions to Day A:**
```bash
# Squat (order 1)
curl -X POST "http://localhost:8080/days/bbbb1111-1111-1111-1111-111111111111/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionId": "aaaa1111-1111-1111-1111-111111111111",
    "order": 1
  }'

# Bench Press (order 2)
curl -X POST "http://localhost:8080/days/bbbb1111-1111-1111-1111-111111111111/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionId": "aaaa2222-2222-2222-2222-222222222222",
    "order": 2
  }'

# Deadlift (order 3)
curl -X POST "http://localhost:8080/days/bbbb1111-1111-1111-1111-111111111111/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionId": "aaaa4444-4444-4444-4444-444444444444",
    "order": 3
  }'
```

**Add prescriptions to Day B:**
```bash
# Squat (order 1)
curl -X POST "http://localhost:8080/days/bbbb2222-2222-2222-2222-222222222222/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionId": "aaaa1111-1111-1111-1111-111111111111",
    "order": 1
  }'

# Overhead Press (order 2)
curl -X POST "http://localhost:8080/days/bbbb2222-2222-2222-2222-222222222222/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionId": "aaaa3333-3333-3333-3333-333333333333",
    "order": 2
  }'

# Power Clean (order 3)
curl -X POST "http://localhost:8080/days/bbbb2222-2222-2222-2222-222222222222/prescriptions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionId": "aaaa5555-5555-5555-5555-555555555555",
    "order": 3
  }'
```

---

### Step 4: Create Cycle

Starting Strength uses a 1-week cycle that repeats indefinitely.

```bash
curl -X POST "http://localhost:8080/cycles" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Starting Strength 1-Week Cycle",
    "lengthWeeks": 1
  }'
```

**Response:**
```json
{
  "data": {
    "id": "cccc1111-1111-1111-1111-111111111111",
    "name": "Starting Strength 1-Week Cycle",
    "lengthWeeks": 1,
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

---

### Step 5: Create Weeks

Create a single week with alternating A/B/A days.

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
    "id": "dddd1111-1111-1111-1111-111111111111",
    "cycleId": "cccc1111-1111-1111-1111-111111111111",
    "weekNumber": 1,
    "name": "Week 1",
    "days": [],
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

**Add days to week (A/B/A pattern):**
```bash
# Position 0: Day A (Monday)
curl -X POST "http://localhost:8080/weeks/dddd1111-1111-1111-1111-111111111111/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "dayId": "bbbb1111-1111-1111-1111-111111111111",
    "position": 0
  }'

# Position 1: Day B (Wednesday)
curl -X POST "http://localhost:8080/weeks/dddd1111-1111-1111-1111-111111111111/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "dayId": "bbbb2222-2222-2222-2222-222222222222",
    "position": 1
  }'

# Position 2: Day A (Friday)
curl -X POST "http://localhost:8080/weeks/dddd1111-1111-1111-1111-111111111111/days" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "dayId": "bbbb1111-1111-1111-1111-111111111111",
    "position": 2
  }'
```

**Note on A/B rotation:** The A/B/A pattern in week 1 becomes B/A/B in week 2 through cycle iteration logic. When `currentCycleIteration` is odd, the week starts with Day A. When even, it starts with Day B. This alternation happens automatically when the user advances through the program.

---

### Step 6: Create Progressions

Create LINEAR progressions for each lift. Different lifts have different increments.

**Linear +5lb (Upper Body - Bench, Press, Power Clean):**
```bash
curl -X POST "http://localhost:8080/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Linear +5lb Upper Body",
    "type": "LINEAR",
    "parameters": {
      "increment": 5.0,
      "maxType": "TRAINING_MAX",
      "triggerType": "AFTER_SESSION"
    }
  }'
```

**Response:**
```json
{
  "data": {
    "id": "eeee1111-1111-1111-1111-111111111111",
    "name": "Linear +5lb Upper Body",
    "type": "LINEAR",
    "parameters": {
      "increment": 5.0,
      "maxType": "TRAINING_MAX",
      "triggerType": "AFTER_SESSION"
    },
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

**Linear +5lb (Squat):**
```bash
curl -X POST "http://localhost:8080/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Linear +5lb Squat",
    "type": "LINEAR",
    "parameters": {
      "increment": 5.0,
      "maxType": "TRAINING_MAX",
      "triggerType": "AFTER_SESSION"
    }
  }'
```

**Linear +10lb (Deadlift):**
```bash
curl -X POST "http://localhost:8080/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Linear +10lb Deadlift",
    "type": "LINEAR",
    "parameters": {
      "increment": 10.0,
      "maxType": "TRAINING_MAX",
      "triggerType": "AFTER_SESSION"
    }
  }'
```

**Progression IDs:**
| Progression | ID |
|-------------|-----|
| Linear +5lb Upper | `eeee1111-1111-1111-1111-111111111111` |
| Linear +5lb Squat | `eeee2222-2222-2222-2222-222222222222` |
| Linear +10lb Deadlift | `eeee3333-3333-3333-3333-333333333333` |

---

### Step 7: Create Program

Create the Starting Strength program, linking it to the cycle.

```bash
curl -X POST "http://localhost:8080/programs" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Starting Strength",
    "slug": "starting-strength",
    "description": "Mark Rippetoe'\''s novice linear progression program. Three days per week, alternating A/B workouts with linear progression every session.",
    "cycleId": "cccc1111-1111-1111-1111-111111111111",
    "defaultRounding": 5.0
  }'
```

**Response:**
```json
{
  "data": {
    "id": "ffff1111-1111-1111-1111-111111111111",
    "name": "Starting Strength",
    "slug": "starting-strength",
    "description": "Mark Rippetoe's novice linear progression program. Three days per week, alternating A/B workouts with linear progression every session.",
    "cycleId": "cccc1111-1111-1111-1111-111111111111",
    "weeklyLookupId": null,
    "dailyLookupId": null,
    "defaultRounding": 5.0,
    "createdAt": "2024-01-15T00:00:00Z",
    "updatedAt": "2024-01-15T00:00:00Z"
  }
}
```

---

### Step 8: Link Progressions to Program

Link each progression to the program, specifying which lift each applies to.

**Squat Progression:**
```bash
curl -X POST "http://localhost:8080/programs/ffff1111-1111-1111-1111-111111111111/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "eeee2222-2222-2222-2222-222222222222",
    "liftId": "11111111-1111-1111-1111-111111111111",
    "priority": 1
  }'
```

**Bench Press Progression:**
```bash
curl -X POST "http://localhost:8080/programs/ffff1111-1111-1111-1111-111111111111/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "eeee1111-1111-1111-1111-111111111111",
    "liftId": "22222222-2222-2222-2222-222222222222",
    "priority": 1
  }'
```

**Overhead Press Progression:**
```bash
curl -X POST "http://localhost:8080/programs/ffff1111-1111-1111-1111-111111111111/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "eeee1111-1111-1111-1111-111111111111",
    "liftId": "33333333-3333-3333-3333-333333333333",
    "priority": 1
  }'
```

**Deadlift Progression:**
```bash
curl -X POST "http://localhost:8080/programs/ffff1111-1111-1111-1111-111111111111/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "eeee3333-3333-3333-3333-333333333333",
    "liftId": "44444444-4444-4444-4444-444444444444",
    "priority": 1
  }'
```

**Power Clean Progression:**
```bash
curl -X POST "http://localhost:8080/programs/ffff1111-1111-1111-1111-111111111111/progressions" \
  -H "X-User-ID: admin-user-id" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "eeee1111-1111-1111-1111-111111111111",
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

# Squat Training Max: 225 lbs
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "11111111-1111-1111-1111-111111111111",
    "type": "TRAINING_MAX",
    "value": 225.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'

# Bench Press Training Max: 155 lbs
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "22222222-2222-2222-2222-222222222222",
    "type": "TRAINING_MAX",
    "value": 155.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'

# Overhead Press Training Max: 95 lbs
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "33333333-3333-3333-3333-333333333333",
    "type": "TRAINING_MAX",
    "value": 95.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'

# Deadlift Training Max: 275 lbs
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "44444444-4444-4444-4444-444444444444",
    "type": "TRAINING_MAX",
    "value": 275.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'

# Power Clean Training Max: 135 lbs
curl -X POST "http://localhost:8080/users/${USER_ID}/lift-maxes" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "55555555-5555-5555-5555-555555555555",
    "type": "TRAINING_MAX",
    "value": 135.0,
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
      "name": "Starting Strength",
      "slug": "starting-strength",
      "description": "Mark Rippetoe's novice linear progression program.",
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

### Session 1: Day A (Week 1, Day Index 0)

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
    "daySlug": "ss-day-a",
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
          {"setNumber": 1, "weight": 225.0, "targetReps": 5, "isWorkSet": true},
          {"setNumber": 2, "weight": 225.0, "targetReps": 5, "isWorkSet": true},
          {"setNumber": 3, "weight": 225.0, "targetReps": 5, "isWorkSet": true}
        ],
        "notes": "Focus on depth - crease of hip below top of knee",
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
          {"setNumber": 1, "weight": 155.0, "targetReps": 5, "isWorkSet": true},
          {"setNumber": 2, "weight": 155.0, "targetReps": 5, "isWorkSet": true},
          {"setNumber": 3, "weight": 155.0, "targetReps": 5, "isWorkSet": true}
        ],
        "notes": "Touch chest, pause briefly, press",
        "restSeconds": 180
      },
      {
        "prescriptionId": "aaaa4444-4444-4444-4444-444444444444",
        "lift": {
          "id": "44444444-4444-4444-4444-444444444444",
          "name": "Deadlift",
          "slug": "deadlift"
        },
        "sets": [
          {"setNumber": 1, "weight": 275.0, "targetReps": 5, "isWorkSet": true}
        ],
        "notes": "Reset each rep - no touch and go",
        "restSeconds": 180
      }
    ]
  }
}
```

### Session 2: Day B (Week 1, Day Index 1)

After advancing to the next day:

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
    "programId": "ffff1111-1111-1111-1111-111111111111",
    "cycleIteration": 1,
    "weekNumber": 1,
    "daySlug": "ss-day-b",
    "date": "2024-01-17",
    "exercises": [
      {
        "prescriptionId": "aaaa1111-1111-1111-1111-111111111111",
        "lift": {
          "id": "11111111-1111-1111-1111-111111111111",
          "name": "Squat",
          "slug": "squat"
        },
        "sets": [
          {"setNumber": 1, "weight": 225.0, "targetReps": 5, "isWorkSet": true},
          {"setNumber": 2, "weight": 225.0, "targetReps": 5, "isWorkSet": true},
          {"setNumber": 3, "weight": 225.0, "targetReps": 5, "isWorkSet": true}
        ],
        "notes": "Focus on depth - crease of hip below top of knee",
        "restSeconds": 180
      },
      {
        "prescriptionId": "aaaa3333-3333-3333-3333-333333333333",
        "lift": {
          "id": "33333333-3333-3333-3333-333333333333",
          "name": "Overhead Press",
          "slug": "overhead-press"
        },
        "sets": [
          {"setNumber": 1, "weight": 95.0, "targetReps": 5, "isWorkSet": true},
          {"setNumber": 2, "weight": 95.0, "targetReps": 5, "isWorkSet": true},
          {"setNumber": 3, "weight": 95.0, "targetReps": 5, "isWorkSet": true}
        ],
        "notes": "Strict press - no leg drive",
        "restSeconds": 180
      },
      {
        "prescriptionId": "aaaa5555-5555-5555-5555-555555555555",
        "lift": {
          "id": "55555555-5555-5555-5555-555555555555",
          "name": "Power Clean",
          "slug": "power-clean"
        },
        "sets": [
          {"setNumber": 1, "weight": 135.0, "targetReps": 3, "isWorkSet": true},
          {"setNumber": 2, "weight": 135.0, "targetReps": 3, "isWorkSet": true},
          {"setNumber": 3, "weight": 135.0, "targetReps": 3, "isWorkSet": true},
          {"setNumber": 4, "weight": 135.0, "targetReps": 3, "isWorkSet": true},
          {"setNumber": 5, "weight": 135.0, "targetReps": 3, "isWorkSet": true}
        ],
        "notes": "Explosive triple extension",
        "restSeconds": 120
      }
    ]
  }
}
```

---

## Progression Behavior

Starting Strength uses linear progression - weights increase after each session.

### Triggering Progression After Session 1

After completing Day A (session 1), trigger progressions for the lifts performed:

**Trigger Squat Progression (+5 lbs):**
```bash
curl -X POST "http://localhost:8080/users/${USER_ID}/progressions/trigger" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "eeee2222-2222-2222-2222-222222222222",
    "liftId": "11111111-1111-1111-1111-111111111111"
  }'
```

**Response:**
```json
{
  "data": {
    "results": [
      {
        "progressionId": "eeee2222-2222-2222-2222-222222222222",
        "liftId": "11111111-1111-1111-1111-111111111111",
        "applied": true,
        "skipped": false,
        "result": {
          "previousValue": 225.0,
          "newValue": 230.0,
          "delta": 5.0,
          "maxType": "TRAINING_MAX",
          "appliedAt": "2024-01-15T18:00:00Z"
        }
      }
    ],
    "totalApplied": 1,
    "totalSkipped": 0,
    "totalErrors": 0
  }
}
```

**Trigger Bench Press Progression (+5 lbs):**
```bash
curl -X POST "http://localhost:8080/users/${USER_ID}/progressions/trigger" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "eeee1111-1111-1111-1111-111111111111",
    "liftId": "22222222-2222-2222-2222-222222222222"
  }'
```

**Response:**
```json
{
  "data": {
    "results": [
      {
        "progressionId": "eeee1111-1111-1111-1111-111111111111",
        "liftId": "22222222-2222-2222-2222-222222222222",
        "applied": true,
        "skipped": false,
        "result": {
          "previousValue": 155.0,
          "newValue": 160.0,
          "delta": 5.0,
          "maxType": "TRAINING_MAX",
          "appliedAt": "2024-01-15T18:00:00Z"
        }
      }
    ],
    "totalApplied": 1,
    "totalSkipped": 0,
    "totalErrors": 0
  }
}
```

**Trigger Deadlift Progression (+10 lbs):**
```bash
curl -X POST "http://localhost:8080/users/${USER_ID}/progressions/trigger" \
  -H "X-User-ID: ${USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "eeee3333-3333-3333-3333-333333333333",
    "liftId": "44444444-4444-4444-4444-444444444444"
  }'
```

**Response:**
```json
{
  "data": {
    "results": [
      {
        "progressionId": "eeee3333-3333-3333-3333-333333333333",
        "liftId": "44444444-4444-4444-4444-444444444444",
        "applied": true,
        "skipped": false,
        "result": {
          "previousValue": 275.0,
          "newValue": 285.0,
          "delta": 10.0,
          "maxType": "TRAINING_MAX",
          "appliedAt": "2024-01-15T18:00:00Z"
        }
      }
    ],
    "totalApplied": 1,
    "totalSkipped": 0,
    "totalErrors": 0
  }
}
```

### Progression Over Multiple Sessions

Here's how the training maxes progress over the first week:

| Session | Day | Squat TM | Bench TM | Press TM | Deadlift TM | Clean TM |
|---------|-----|----------|----------|----------|-------------|----------|
| Start | - | 225 | 155 | 95 | 275 | 135 |
| After 1 | A | 230 | 160 | 95 | 285 | 135 |
| After 2 | B | 235 | 160 | 100 | 285 | 140 |
| After 3 | A | 240 | 165 | 100 | 295 | 140 |

**Week 1 Summary:**
- Squat: 225 → 240 (+15 lbs, 3 sessions)
- Bench: 155 → 165 (+10 lbs, 2 sessions - Day A only)
- Press: 95 → 100 (+5 lbs, 1 session - Day B only)
- Deadlift: 275 → 295 (+20 lbs, 2 sessions - Day A only)
- Power Clean: 135 → 140 (+5 lbs, 1 session - Day B only)

### Session 4: Day A (Next Week, After Progressions)

After advancing to week 2 and generating the workout:

```json
{
  "data": {
    "cycleIteration": 2,
    "weekNumber": 1,
    "daySlug": "ss-day-a",
    "exercises": [
      {
        "lift": {"name": "Squat"},
        "sets": [
          {"setNumber": 1, "weight": 240.0, "targetReps": 5},
          {"setNumber": 2, "weight": 240.0, "targetReps": 5},
          {"setNumber": 3, "weight": 240.0, "targetReps": 5}
        ]
      },
      {
        "lift": {"name": "Bench Press"},
        "sets": [
          {"setNumber": 1, "weight": 165.0, "targetReps": 5},
          {"setNumber": 2, "weight": 165.0, "targetReps": 5},
          {"setNumber": 3, "weight": 165.0, "targetReps": 5}
        ]
      },
      {
        "lift": {"name": "Deadlift"},
        "sets": [
          {"setNumber": 1, "weight": 295.0, "targetReps": 5}
        ]
      }
    ]
  }
}
```

### Viewing Progression History

```bash
curl -X GET "http://localhost:8080/users/${USER_ID}/progression-history?limit=10" \
  -H "X-User-ID: ${USER_ID}"
```

**Response:**
```json
{
  "data": [
    {
      "id": "history-1",
      "progressionId": "eeee2222-2222-2222-2222-222222222222",
      "progressionName": "Linear +5lb Squat",
      "progressionType": "LINEAR",
      "liftId": "11111111-1111-1111-1111-111111111111",
      "liftName": "Squat",
      "previousValue": 235.0,
      "newValue": 240.0,
      "delta": 5.0,
      "triggerType": "AFTER_SESSION",
      "appliedAt": "2024-01-19T18:00:00Z"
    },
    {
      "id": "history-2",
      "progressionId": "eeee1111-1111-1111-1111-111111111111",
      "progressionName": "Linear +5lb Upper Body",
      "progressionType": "LINEAR",
      "liftId": "22222222-2222-2222-2222-222222222222",
      "liftName": "Bench Press",
      "previousValue": 160.0,
      "newValue": 165.0,
      "delta": 5.0,
      "triggerType": "AFTER_SESSION",
      "appliedAt": "2024-01-19T18:00:00Z"
    },
    {
      "id": "history-3",
      "progressionId": "eeee3333-3333-3333-3333-333333333333",
      "progressionName": "Linear +10lb Deadlift",
      "progressionType": "LINEAR",
      "liftId": "44444444-4444-4444-4444-444444444444",
      "liftName": "Deadlift",
      "previousValue": 285.0,
      "newValue": 295.0,
      "delta": 10.0,
      "triggerType": "AFTER_SESSION",
      "appliedAt": "2024-01-19T18:00:00Z"
    }
  ],
  "meta": {
    "total": 9,
    "limit": 10,
    "offset": 0,
    "hasMore": false
  }
}
```

---

## Complete Configuration Summary

### Entity Relationship Diagram

```
Program: Starting Strength
├── Cycle: 1-week cycle
│   └── Week 1
│       ├── Position 0: Day A (Squat/Bench/Deadlift)
│       ├── Position 1: Day B (Squat/Press/Clean)
│       └── Position 2: Day A (Squat/Bench/Deadlift)
│
├── Days
│   ├── Day A
│   │   ├── Prescription: Squat 3x5 @ 100% TM
│   │   ├── Prescription: Bench 3x5 @ 100% TM
│   │   └── Prescription: Deadlift 1x5 @ 100% TM
│   │
│   └── Day B
│       ├── Prescription: Squat 3x5 @ 100% TM
│       ├── Prescription: Press 3x5 @ 100% TM
│       └── Prescription: Power Clean 5x3 @ 100% TM
│
└── Program Progressions
    ├── Squat → Linear +5lb (AFTER_SESSION)
    ├── Bench → Linear +5lb (AFTER_SESSION)
    ├── Press → Linear +5lb (AFTER_SESSION)
    ├── Deadlift → Linear +10lb (AFTER_SESSION)
    └── Power Clean → Linear +5lb (AFTER_SESSION)
```

### JSON Configuration Export

```json
{
  "program": {
    "name": "Starting Strength",
    "slug": "starting-strength",
    "description": "Mark Rippetoe's novice linear progression program.",
    "defaultRounding": 5.0
  },
  "cycle": {
    "name": "Starting Strength 1-Week Cycle",
    "lengthWeeks": 1
  },
  "weeks": [
    {
      "weekNumber": 1,
      "name": "Week 1",
      "days": [
        {"position": 0, "daySlug": "ss-day-a"},
        {"position": 1, "daySlug": "ss-day-b"},
        {"position": 2, "daySlug": "ss-day-a"}
      ]
    }
  ],
  "days": [
    {
      "name": "Day A",
      "slug": "ss-day-a",
      "prescriptions": [
        {"liftSlug": "squat", "sets": 3, "reps": 5, "percentage": 100, "order": 1},
        {"liftSlug": "bench-press", "sets": 3, "reps": 5, "percentage": 100, "order": 2},
        {"liftSlug": "deadlift", "sets": 1, "reps": 5, "percentage": 100, "order": 3}
      ]
    },
    {
      "name": "Day B",
      "slug": "ss-day-b",
      "prescriptions": [
        {"liftSlug": "squat", "sets": 3, "reps": 5, "percentage": 100, "order": 1},
        {"liftSlug": "overhead-press", "sets": 3, "reps": 5, "percentage": 100, "order": 2},
        {"liftSlug": "power-clean", "sets": 5, "reps": 3, "percentage": 100, "order": 3}
      ]
    }
  ],
  "progressions": [
    {
      "name": "Linear +5lb Squat",
      "type": "LINEAR",
      "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"},
      "liftSlug": "squat"
    },
    {
      "name": "Linear +5lb Upper Body",
      "type": "LINEAR",
      "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"},
      "liftSlugs": ["bench-press", "overhead-press", "power-clean"]
    },
    {
      "name": "Linear +10lb Deadlift",
      "type": "LINEAR",
      "parameters": {"increment": 10.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"},
      "liftSlug": "deadlift"
    }
  ]
}
```

---

## See Also

- [API Reference](../api-reference.md) - Complete endpoint documentation
- [Workflows](../workflows.md) - Multi-step API workflows
- [Starting Strength Program Specification](../../programs/004-starting-strength.md) - Full program algorithm

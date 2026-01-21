# PowerPro API Workflows

This document describes common multi-step workflows for integrating with the PowerPro API. Each workflow shows the sequence of API calls needed to accomplish common tasks, with copy-paste ready examples.

## Table of Contents

- [User Onboarding Workflow](#user-onboarding-workflow)
- [Program Enrollment Workflow](#program-enrollment-workflow)
- [Workout Generation Workflow](#workout-generation-workflow)
- [Progression Trigger Workflow](#progression-trigger-workflow)

---

## User Onboarding Workflow

Set up a new user with lifts and training maxes so they can generate workouts.

### Prerequisites

- A valid user ID (authentication handled externally)
- Knowledge of the user's current maxes for competition lifts

### Step 1: Verify Available Lifts

First, check which lifts exist in the system. Competition lifts (squat, bench, deadlift) are typically pre-configured by admins.

```bash
curl -X GET "http://localhost:8080/lifts?is_competition_lift=true" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

**Expected Response** `200 OK`:
```json
{
  "data": [
    {
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "name": "Squat",
      "slug": "squat",
      "isCompetitionLift": true
    },
    {
      "id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
      "name": "Bench Press",
      "slug": "bench-press",
      "isCompetitionLift": true
    },
    {
      "id": "d4e5f6a7-b8c9-0123-defg-456789012345",
      "name": "Deadlift",
      "slug": "deadlift",
      "isCompetitionLift": true
    }
  ],
  "page": 1,
  "pageSize": 20,
  "totalItems": 3,
  "totalPages": 1
}
```

Note the lift IDs for the next step.

### Step 2: Create Training Maxes for Each Lift

Create a training max for each competition lift. Training maxes are typically 85-90% of your true 1RM.

**Squat Training Max:**
```bash
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/lift-maxes" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "type": "TRAINING_MAX",
    "value": 315.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'
```

**Bench Press Training Max:**
```bash
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/lift-maxes" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
    "type": "TRAINING_MAX",
    "value": 225.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'
```

**Deadlift Training Max:**
```bash
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/lift-maxes" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "d4e5f6a7-b8c9-0123-defg-456789012345",
    "type": "TRAINING_MAX",
    "value": 405.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'
```

**Expected Response** `201 Created` (for each):
```json
{
  "id": "e5f6a7b8-c9d0-1234-efgh-456789012345",
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "type": "TRAINING_MAX",
  "value": 315.0,
  "effectiveDate": "2024-01-15T00:00:00Z",
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

### Step 3: Enroll User in a Program

Now enroll the user in a training program.

```bash
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/program" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "programId": "g7h8i9j0-k1l2-3456-mnop-789012345678"
  }'
```

**Expected Response** `201 Created`:
```json
{
  "id": "q4r5s6t7-u8v9-0123-zabc-456789012345",
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "program": {
    "id": "g7h8i9j0-k1l2-3456-mnop-789012345678",
    "name": "Wendler 5/3/1 BBB",
    "slug": "wendler-531-bbb",
    "description": "The classic Boring But Big variant of 5/3/1.",
    "cycleLengthWeeks": 4
  },
  "state": {
    "currentWeek": 1,
    "currentCycleIteration": 1,
    "currentDayIndex": 0
  },
  "enrolledAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

### Step 4: Generate First Workout

The user is now ready to generate their first workout.

```bash
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/workout" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

**Expected Response** `200 OK`:
```json
{
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "programId": "g7h8i9j0-k1l2-3456-mnop-789012345678",
  "cycleIteration": 1,
  "weekNumber": 1,
  "daySlug": "squat-day",
  "date": "2024-01-15",
  "exercises": [
    {
      "prescriptionId": "h8i9j0k1-l2m3-4567-mnop-789012345678",
      "lift": {
        "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "name": "Squat",
        "slug": "squat"
      },
      "sets": [
        {"setNumber": 1, "weight": 265.0, "targetReps": 5, "isWorkSet": true}
      ],
      "notes": "Focus on depth",
      "restSeconds": 180
    }
  ]
}
```

### Common Pitfalls

1. **Missing training maxes**: If you try to generate a workout without setting up training maxes, you'll get a `422 Unprocessable Entity` error with details about which lift maxes are missing.

2. **Wrong max type**: Programs typically use `TRAINING_MAX`. If you create a `ONE_RM` but the program expects `TRAINING_MAX`, workout generation will fail.

3. **Enrolling before maxes**: Always create training maxes before enrolling in a program, so the user can immediately generate workouts.

4. **Using 1RM instead of Training Max**: Training maxes should be ~90% of your true 1RM. Using your actual 1RM will make weights too heavy.

---

## Program Enrollment Workflow

Explore available programs and enroll a user.

### Prerequisites

- User has training maxes set up for lifts used in the program
- User authentication configured

### Step 1: List Available Programs

Get a list of all programs to show the user.

```bash
curl -X GET "http://localhost:8080/programs" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

**Expected Response** `200 OK`:
```json
{
  "data": [
    {
      "id": "g7h8i9j0-k1l2-3456-mnop-789012345678",
      "name": "Wendler 5/3/1 BBB",
      "slug": "wendler-531-bbb",
      "description": "The classic Boring But Big variant of 5/3/1.",
      "cycleId": "t0u1v2w3-x4y5-6789-abcd-901234567890",
      "defaultRounding": 5.0
    },
    {
      "id": "w3x4y5z6-a7b8-9012-fghi-345678901234",
      "name": "Starting Strength",
      "slug": "starting-strength",
      "description": "Classic linear progression for novice lifters.",
      "cycleId": "a7b8c9d0-e1f2-3456-hijk-678901234567",
      "defaultRounding": 5.0
    }
  ],
  "page": 1,
  "pageSize": 20,
  "totalItems": 2,
  "totalPages": 1
}
```

### Step 2: Get Program Details

Get detailed information about a specific program, including its cycle structure.

```bash
curl -X GET "http://localhost:8080/programs/g7h8i9j0-k1l2-3456-mnop-789012345678" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

**Expected Response** `200 OK`:
```json
{
  "id": "g7h8i9j0-k1l2-3456-mnop-789012345678",
  "name": "Wendler 5/3/1 BBB",
  "slug": "wendler-531-bbb",
  "description": "The classic Boring But Big variant of 5/3/1.",
  "cycle": {
    "id": "t0u1v2w3-x4y5-6789-abcd-901234567890",
    "name": "4 Week 5/3/1 Cycle",
    "lengthWeeks": 4,
    "weeks": [
      {"id": "s9t0u1v2-w3x4-5678-zabc-890123456789", "weekNumber": 1},
      {"id": "u1v2w3x4-y5z6-7890-bcde-012345678901", "weekNumber": 2},
      {"id": "y5z6a7b8-c9d0-1234-fghi-456789012345", "weekNumber": 3},
      {"id": "h4i5j6k7-l8m9-0123-pqrs-456789012345", "weekNumber": 4}
    ]
  },
  "weeklyLookup": {
    "id": "c9d0e1f2-a3b4-5678-jklm-890123456789",
    "name": "5/3/1 Top Set Percentages"
  },
  "dailyLookup": null,
  "defaultRounding": 5.0
}
```

### Step 3: Check User's Current Enrollment (Optional)

Before enrolling, check if the user is already enrolled in a program.

```bash
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/program" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

**If enrolled** `200 OK`: Returns current enrollment details.

**If not enrolled** `404 Not Found`:
```json
{
  "error": "enrollment not found"
}
```

### Step 4: Enroll in Program

Enroll the user in the selected program. If already enrolled, this replaces the existing enrollment.

```bash
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/program" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "programId": "g7h8i9j0-k1l2-3456-mnop-789012345678"
  }'
```

**Expected Response** `201 Created`:
```json
{
  "id": "q4r5s6t7-u8v9-0123-zabc-456789012345",
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "program": {
    "id": "g7h8i9j0-k1l2-3456-mnop-789012345678",
    "name": "Wendler 5/3/1 BBB",
    "slug": "wendler-531-bbb",
    "description": "The classic Boring But Big variant of 5/3/1.",
    "cycleLengthWeeks": 4
  },
  "state": {
    "currentWeek": 1,
    "currentCycleIteration": 1,
    "currentDayIndex": 0
  },
  "enrolledAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

### Step 5: View Program Progression Rules (Optional)

Check what progression rules apply to this program.

```bash
curl -X GET "http://localhost:8080/programs/g7h8i9j0-k1l2-3456-mnop-789012345678/progressions" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

**Expected Response** `200 OK`:
```json
[
  {
    "id": "m0n1o2p3-q4r5-6789-vwxy-012345678901",
    "programId": "g7h8i9j0-k1l2-3456-mnop-789012345678",
    "progressionId": "j7k8l9m0-n1o2-3456-stuv-789012345678",
    "liftId": null,
    "priority": 1
  }
]
```

### Common Pitfalls

1. **Re-enrolling resets state**: Enrolling in a program (even the same one) resets the user's state to week 1, day 0. Warn users before switching programs.

2. **Missing training maxes**: If the new program uses lifts the user doesn't have training maxes for, workout generation will fail.

3. **Program ID not found**: Always validate the program ID exists before enrolling to give users a better error message.

---

## Workout Generation Workflow

Generate workouts for daily training and advance through the program.

### Prerequisites

- User is enrolled in a program
- User has training maxes set up for all lifts used in the program

### Step 1: Get Current Workout

Generate the workout for the user's current position in the program.

```bash
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/workout" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

**Expected Response** `200 OK`:
```json
{
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "programId": "g7h8i9j0-k1l2-3456-mnop-789012345678",
  "cycleIteration": 1,
  "weekNumber": 1,
  "daySlug": "squat-day",
  "date": "2024-01-15",
  "exercises": [
    {
      "prescriptionId": "h8i9j0k1-l2m3-4567-mnop-789012345678",
      "lift": {
        "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "name": "Squat",
        "slug": "squat"
      },
      "sets": [
        {"setNumber": 1, "weight": 205.0, "targetReps": 5, "isWorkSet": true},
        {"setNumber": 2, "weight": 235.0, "targetReps": 5, "isWorkSet": true},
        {"setNumber": 3, "weight": 265.0, "targetReps": 5, "isWorkSet": true}
      ],
      "notes": "Focus on depth",
      "restSeconds": 180
    },
    {
      "prescriptionId": "j0k1l2m3-n4o5-6789-opqr-901234567890",
      "lift": {
        "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "name": "Squat",
        "slug": "squat"
      },
      "sets": [
        {"setNumber": 1, "weight": 155.0, "targetReps": 10, "isWorkSet": true},
        {"setNumber": 2, "weight": 155.0, "targetReps": 10, "isWorkSet": true},
        {"setNumber": 3, "weight": 155.0, "targetReps": 10, "isWorkSet": true},
        {"setNumber": 4, "weight": 155.0, "targetReps": 10, "isWorkSet": true},
        {"setNumber": 5, "weight": 155.0, "targetReps": 10, "isWorkSet": true}
      ],
      "notes": "BBB volume work",
      "restSeconds": 90
    }
  ]
}
```

### Step 2: Complete Workout and Advance Day

After completing the workout, advance to the next training day.

```bash
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/program-state/advance" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "advanceType": "day"
  }'
```

**Expected Response** `200 OK`:
```json
{
  "id": "q4r5s6t7-u8v9-0123-zabc-456789012345",
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "program": {
    "id": "g7h8i9j0-k1l2-3456-mnop-789012345678",
    "name": "Wendler 5/3/1 BBB",
    "slug": "wendler-531-bbb",
    "cycleLengthWeeks": 4
  },
  "state": {
    "currentWeek": 1,
    "currentCycleIteration": 1,
    "currentDayIndex": 1
  },
  "enrolledAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-16T08:00:00Z"
}
```

### Step 3: Preview Future Workouts (Optional)

Preview what a workout will look like for a specific week/day without changing state.

```bash
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/workout/preview?week=3&day=squat-day" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

**Expected Response** `200 OK`:
```json
{
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "programId": "g7h8i9j0-k1l2-3456-mnop-789012345678",
  "cycleIteration": 1,
  "weekNumber": 3,
  "daySlug": "squat-day",
  "date": "2024-01-15",
  "exercises": [
    {
      "prescriptionId": "h8i9j0k1-l2m3-4567-mnop-789012345678",
      "lift": {
        "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "name": "Squat",
        "slug": "squat"
      },
      "sets": [
        {"setNumber": 1, "weight": 235.0, "targetReps": 5, "isWorkSet": true},
        {"setNumber": 2, "weight": 270.0, "targetReps": 3, "isWorkSet": true},
        {"setNumber": 3, "weight": 300.0, "targetReps": 1, "isWorkSet": true}
      ],
      "notes": "Week 3 - heavy singles",
      "restSeconds": 180
    }
  ]
}
```

### Step 4: Advance Week (End of Week)

When all days in a week are complete, advance to the next week.

```bash
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/program-state/advance" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "advanceType": "week"
  }'
```

**Expected Response** `200 OK`:
```json
{
  "state": {
    "currentWeek": 2,
    "currentCycleIteration": 1,
    "currentDayIndex": 0
  }
}
```

### End of Cycle Behavior

When advancing past the last week of a cycle:
- `currentWeek` resets to 1
- `currentCycleIteration` increments
- `currentDayIndex` resets to 0

This is when cycle-based progressions (like 5/3/1's +5lb per cycle) should be triggered.

### Common Pitfalls

1. **Forgetting to advance**: Call the advance endpoint after each workout to track progress correctly.

2. **Advancing without completing**: Consider adding client-side tracking of completed sets before allowing advancement.

3. **Missing lift max**: If a lift max is deleted or missing, workout generation will fail with a `400` or `422` error.

4. **Week vs Day advancement**: Use `"advanceType": "day"` for daily training, `"advanceType": "week"` to skip remaining days in a week.

---

## Progression Trigger Workflow

Understand and trigger progression rules to increase weights over time.

### Prerequisites

- User is enrolled in a program with configured progressions
- User has training maxes for lifts in the program

### Step 1: Understand Program's Progression Rules

Get the progression configurations for the user's program.

```bash
curl -X GET "http://localhost:8080/programs/g7h8i9j0-k1l2-3456-mnop-789012345678/progressions" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

**Expected Response** `200 OK`:
```json
[
  {
    "id": "m0n1o2p3-q4r5-6789-vwxy-012345678901",
    "programId": "g7h8i9j0-k1l2-3456-mnop-789012345678",
    "progressionId": "j7k8l9m0-n1o2-3456-stuv-789012345678",
    "liftId": null,
    "priority": 1
  },
  {
    "id": "n1o2p3q4-r5s6-7890-wxyz-123456789012",
    "programId": "g7h8i9j0-k1l2-3456-mnop-789012345678",
    "progressionId": "k8l9m0n1-o2p3-4567-tuvw-890123456789",
    "liftId": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
    "priority": 2
  }
]
```

### Step 2: View Progression Rule Details

Get details about a specific progression rule.

```bash
curl -X GET "http://localhost:8080/progressions/j7k8l9m0-n1o2-3456-stuv-789012345678" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

**Expected Response** `200 OK`:
```json
{
  "id": "j7k8l9m0-n1o2-3456-stuv-789012345678",
  "name": "Cycle +5lb (lower body)",
  "type": "CYCLE",
  "parameters": {
    "increment": 5.0,
    "maxType": "TRAINING_MAX"
  },
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

**Progression Types:**

| Type | Trigger | Description |
|------|---------|-------------|
| `LINEAR` | `AFTER_SESSION` or `AFTER_WEEK` | Add weight after each session or week |
| `CYCLE` | End of cycle | Add weight at the end of each training cycle |

### Step 3: Trigger Progression Manually

At the appropriate time (e.g., end of cycle), trigger the progression to increase maxes.

**Trigger for all configured lifts:**
```bash
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/progressions/trigger" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "j7k8l9m0-n1o2-3456-stuv-789012345678"
  }'
```

**Expected Response** `200 OK`:
```json
{
  "results": [
    {
      "progressionId": "j7k8l9m0-n1o2-3456-stuv-789012345678",
      "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "applied": true,
      "skipped": false,
      "result": {
        "previousValue": 315.0,
        "newValue": 320.0,
        "delta": 5.0,
        "maxType": "TRAINING_MAX",
        "appliedAt": "2024-01-28T10:30:00Z"
      }
    },
    {
      "progressionId": "j7k8l9m0-n1o2-3456-stuv-789012345678",
      "liftId": "d4e5f6a7-b8c9-0123-defg-456789012345",
      "applied": true,
      "skipped": false,
      "result": {
        "previousValue": 405.0,
        "newValue": 410.0,
        "delta": 5.0,
        "maxType": "TRAINING_MAX",
        "appliedAt": "2024-01-28T10:30:00Z"
      }
    }
  ],
  "totalApplied": 2,
  "totalSkipped": 0,
  "totalErrors": 0
}
```

**Trigger for a specific lift only:**
```bash
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/progressions/trigger" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "j7k8l9m0-n1o2-3456-stuv-789012345678",
    "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  }'
```

### Step 4: Verify Max Was Updated

Check the user's current training max to verify the progression was applied.

```bash
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/lift-maxes/current?lift=a1b2c3d4-e5f6-7890-abcd-ef1234567890&type=TRAINING_MAX" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

**Expected Response** `200 OK`:
```json
{
  "id": "new-liftmax-uuid",
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "type": "TRAINING_MAX",
  "value": 320.0,
  "effectiveDate": "2024-01-28T10:30:00Z",
  "createdAt": "2024-01-28T10:30:00Z",
  "updatedAt": "2024-01-28T10:30:00Z"
}
```

### Step 5: View Progression History

Review all past progressions for the user.

```bash
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/progression-history" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

**Expected Response** `200 OK`:
```json
{
  "data": [
    {
      "id": "r5s6t7u8-v9w0-1234-bcde-567890123456",
      "userId": "550e8400-e29b-41d4-a716-446655440000",
      "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "progressionId": "j7k8l9m0-n1o2-3456-stuv-789012345678",
      "previousValue": 315.0,
      "newValue": 320.0,
      "delta": 5.0,
      "appliedAt": "2024-01-28T10:30:00Z",
      "triggeredBy": "CYCLE_END",
      "cycleIteration": 1,
      "weekNumber": 4
    }
  ],
  "page": 1,
  "pageSize": 20,
  "totalItems": 1,
  "totalPages": 1
}
```

**Filter history by lift:**
```bash
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/progression-history?lift_id=a1b2c3d4-e5f6-7890-abcd-ef1234567890" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### Force Triggering (Override Skip Protection)

If a progression was already applied this period but you want to apply it again:

```bash
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/progressions/trigger" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "j7k8l9m0-n1o2-3456-stuv-789012345678",
    "force": true
  }'
```

### Common Pitfalls

1. **Triggering too early**: For `CYCLE` progressions, only trigger at the end of the cycle. Triggering mid-cycle will increase weights prematurely.

2. **Double progression**: Without `force: true`, the API prevents applying the same progression twice in one period. This is a safety feature.

3. **Wrong progression type**: `LINEAR` progressions with `AFTER_SESSION` trigger type should be applied after each workout, not at cycle end.

4. **Forgetting to trigger**: Unlike some systems, progressions don't happen automatically. Your client must call the trigger endpoint.

5. **Upper vs lower body increments**: Many programs use +5lb for lower body and +2.5lb for upper body. Make sure the correct progression is configured for each lift.

---

## See Also

- [API Reference](./api-reference.md) - Complete endpoint documentation
- [Example Requests](./example-requests.md) - Copy-paste ready examples
- [Example Responses](./example-responses.md) - Response format examples
- [Error Documentation](./errors.md) - Error handling guide

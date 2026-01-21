# PowerPro API Reference

PowerPro is a headless API for powerlifting program management. This document provides comprehensive documentation for all API endpoints.

## Base URL

```
http://localhost:{port}
```

## Authentication

All endpoints (except `/health`) require authentication via HTTP headers.

### Headers

| Header | Description |
|--------|-------------|
| `Authorization` | Bearer token format: `Bearer {userId}` |
| `X-User-ID` | Alternative: User ID directly (for development/testing) |
| `X-Admin` | Set to `"true"` for admin privileges |

### Authentication Levels

- **Public**: No authentication required
- **Authenticated**: Any authenticated user
- **Admin**: Requires `X-Admin: true`
- **Owner/Admin**: User must own the resource or be admin

---

## Common Response Formats

### Paginated Response

```json
{
  "data": [...],
  "page": 1,
  "pageSize": 20,
  "totalItems": 100,
  "totalPages": 5
}
```

### Error Response

```json
{
  "error": "error message",
  "details": ["detail 1", "detail 2"]
}
```

### Response with Warnings

```json
{
  "data": {...},
  "warnings": ["warning message"]
}
```

---

## Common Query Parameters (List Endpoints)

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | int | 1 | Page number (1-indexed) |
| `pageSize` | int | 20 | Items per page (max: 100) |
| `sortBy` | string | varies | Field to sort by |
| `sortOrder` | string | "asc" | "asc" or "desc" |

---

## HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | Success (GET, PUT) |
| 201 | Created (POST) |
| 204 | No Content (DELETE) |
| 400 | Bad Request (validation error, malformed request) |
| 401 | Unauthorized (missing authentication) |
| 403 | Forbidden (insufficient permissions) |
| 404 | Not Found |
| 409 | Conflict (duplicate slug, FK constraint) |
| 422 | Unprocessable Entity (missing lift max, etc.) |
| 500 | Internal Server Error |

---

## Endpoints

### Health Check

#### GET /health

Check API health status.

**Auth**: Public

**Response** `200 OK`:
```json
{"status": "ok"}
```

---

### Lifts

Manage exercises (lifts) in the system.

#### GET /lifts

List all lifts with pagination.

**Auth**: Authenticated

**Query Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `page` | int | Page number |
| `pageSize` | int | Items per page |
| `sortBy` | string | "name" or "created_at" |
| `sortOrder` | string | "asc" or "desc" |
| `is_competition_lift` | bool | Filter by competition lift status |

**Response** `200 OK`:
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "Squat",
      "slug": "squat",
      "isCompetitionLift": true,
      "parentLiftId": null,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "page": 1,
  "pageSize": 20,
  "totalItems": 3,
  "totalPages": 1
}
```

#### GET /lifts/{id}

Get a lift by ID.

**Auth**: Authenticated

**Response** `200 OK`:
```json
{
  "id": "uuid",
  "name": "Squat",
  "slug": "squat",
  "isCompetitionLift": true,
  "parentLiftId": null,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### GET /lifts/by-slug/{slug}

Get a lift by slug.

**Auth**: Authenticated

**Response**: Same as GET /lifts/{id}

#### POST /lifts

Create a new lift.

**Auth**: Admin

**Request Body**:
```json
{
  "name": "Squat",
  "slug": "squat",
  "isCompetitionLift": true,
  "parentLiftId": "uuid"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Lift name |
| `slug` | string | No | URL-friendly identifier (auto-generated from name if omitted) |
| `isCompetitionLift` | bool | No | Whether this is a competition lift |
| `parentLiftId` | string | No | Parent lift ID for variations |

**Response** `201 Created`: Lift object

#### PUT /lifts/{id}

Update a lift.

**Auth**: Admin

**Request Body**:
```json
{
  "name": "Back Squat",
  "slug": "back-squat",
  "isCompetitionLift": true,
  "parentLiftId": "uuid",
  "clearParentLift": false
}
```

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | New name |
| `slug` | string | New slug |
| `isCompetitionLift` | bool | Competition lift status |
| `parentLiftId` | string | New parent lift ID |
| `clearParentLift` | bool | Set to true to remove parent lift |

**Response** `200 OK`: Updated lift object

#### DELETE /lifts/{id}

Delete a lift.

**Auth**: Admin

**Response** `204 No Content`

**Errors**:
- `409 Conflict`: Lift is referenced by child lifts or other records

---

### Lift Maxes

Manage user's personal records (1RM, training max).

#### GET /users/{userId}/lift-maxes

List all lift maxes for a user.

**Auth**: Owner/Admin

**Query Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `page` | int | Page number |
| `pageSize` | int | Items per page |
| `sortOrder` | string | "asc" or "desc" (default: desc by effective_date) |
| `lift_id` | string | Filter by lift ID |
| `type` | string | Filter by type: "ONE_RM" or "TRAINING_MAX" |

**Response** `200 OK`:
```json
{
  "data": [
    {
      "id": "uuid",
      "userId": "user-uuid",
      "liftId": "lift-uuid",
      "type": "TRAINING_MAX",
      "value": 315.0,
      "effectiveDate": "2024-01-01T00:00:00Z",
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "page": 1,
  "pageSize": 20,
  "totalItems": 1,
  "totalPages": 1
}
```

#### GET /users/{userId}/lift-maxes/current

Get the most recent lift max for a user, lift, and type.

**Auth**: Owner/Admin

**Query Parameters** (required):
| Parameter | Type | Description |
|-----------|------|-------------|
| `lift` | string | Lift ID (UUID) |
| `type` | string | "ONE_RM" or "TRAINING_MAX" |

**Response** `200 OK`: Single LiftMax object

#### GET /lift-maxes/{id}

Get a lift max by ID.

**Auth**: Owner/Admin (ownership checked on resource)

**Response** `200 OK`: LiftMax object

#### GET /lift-maxes/{id}/convert

Convert a lift max between 1RM and Training Max.

**Auth**: Owner/Admin

**Query Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `to_type` | string | Yes | Target type: "ONE_RM" or "TRAINING_MAX" |
| `percentage` | float | No | Conversion percentage (default: 90) |

**Response** `200 OK`:
```json
{
  "originalValue": 350.0,
  "originalType": "ONE_RM",
  "convertedValue": 315.0,
  "convertedType": "TRAINING_MAX",
  "percentage": 90.0
}
```

#### POST /users/{userId}/lift-maxes

Create a new lift max.

**Auth**: Owner/Admin

**Request Body**:
```json
{
  "liftId": "uuid",
  "type": "TRAINING_MAX",
  "value": 315.0,
  "effectiveDate": "2024-01-01T00:00:00Z"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `liftId` | string | Yes | Lift ID |
| `type` | string | Yes | "ONE_RM" or "TRAINING_MAX" |
| `value` | float | Yes | Weight value (must be positive) |
| `effectiveDate` | datetime | No | Date when this max became effective |

**Response** `201 Created`: LiftMax object

#### PUT /lift-maxes/{id}

Update a lift max.

**Auth**: Owner/Admin

**Request Body**:
```json
{
  "value": 320.0,
  "effectiveDate": "2024-01-15T00:00:00Z"
}
```

Note: `type` and `liftId` cannot be changed after creation.

**Response** `200 OK`: Updated LiftMax object

#### DELETE /lift-maxes/{id}

Delete a lift max.

**Auth**: Owner/Admin

**Response** `204 No Content`

---

### Prescriptions

Manage exercise prescriptions (what to do for a single exercise slot).

#### GET /prescriptions

List all prescriptions.

**Auth**: Authenticated

**Query Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `page` | int | Page number |
| `pageSize` | int | Items per page |
| `sortBy` | string | "order" or "created_at" |
| `sortOrder` | string | "asc" or "desc" |
| `lift_id` | string | Filter by lift ID |

**Response** `200 OK`:
```json
{
  "data": [
    {
      "id": "uuid",
      "liftId": "lift-uuid",
      "loadStrategy": {
        "type": "PERCENT_OF",
        "maxType": "TRAINING_MAX",
        "percentage": 85.0
      },
      "setScheme": {
        "type": "FIXED",
        "sets": 5,
        "reps": 5
      },
      "order": 1,
      "notes": "Focus on depth",
      "restSeconds": 180,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "page": 1,
  "pageSize": 20,
  "totalItems": 1,
  "totalPages": 1
}
```

#### GET /prescriptions/{id}

Get a prescription by ID.

**Auth**: Authenticated

**Response** `200 OK`: Prescription object

#### POST /prescriptions

Create a new prescription.

**Auth**: Admin

**Request Body**:
```json
{
  "liftId": "uuid",
  "loadStrategy": {
    "type": "PERCENT_OF",
    "maxType": "TRAINING_MAX",
    "percentage": 85.0
  },
  "setScheme": {
    "type": "FIXED",
    "sets": 5,
    "reps": 5
  },
  "order": 1,
  "notes": "Focus on depth",
  "restSeconds": 180
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `liftId` | string | Yes | Lift ID |
| `loadStrategy` | object | Yes | How to calculate weight |
| `setScheme` | object | Yes | Sets and reps structure |
| `order` | int | No | Display order (default: 0) |
| `notes` | string | No | Notes for the exercise |
| `restSeconds` | int | No | Rest time between sets |

**LoadStrategy Types**:

```json
// PERCENT_OF - Calculate weight as percentage of a max
{
  "type": "PERCENT_OF",
  "maxType": "TRAINING_MAX",
  "percentage": 85.0,
  "lookupKey": "week",
  "roundTo": 5.0
}
```

**SetScheme Types**:

```json
// FIXED - Fixed sets and reps (e.g., 5x5)
{
  "type": "FIXED",
  "sets": 5,
  "reps": 5,
  "isAmrap": false
}

// RAMP - Ramping percentages (e.g., warmup sets)
{
  "type": "RAMP",
  "sets": [
    {"percentage": 50, "reps": 5},
    {"percentage": 60, "reps": 5},
    {"percentage": 70, "reps": 3}
  ]
}
```

**Response** `201 Created`: Prescription object

#### PUT /prescriptions/{id}

Update a prescription.

**Auth**: Admin

**Request Body**: Same fields as POST (all optional)

**Response** `200 OK`: Updated Prescription object

#### DELETE /prescriptions/{id}

Delete a prescription.

**Auth**: Admin

**Response** `204 No Content`

#### POST /prescriptions/{id}/resolve

Resolve a prescription to concrete sets/reps/weights for a user.

**Auth**: Authenticated

**Request Body**:
```json
{
  "userId": "user-uuid"
}
```

**Response** `200 OK`:
```json
{
  "prescriptionId": "uuid",
  "lift": {
    "id": "lift-uuid",
    "name": "Squat",
    "slug": "squat"
  },
  "sets": [
    {
      "setNumber": 1,
      "weight": 265.0,
      "targetReps": 5,
      "isWorkSet": true
    }
  ],
  "notes": "Focus on depth",
  "restSeconds": 180
}
```

**Errors**:
- `422 Unprocessable Entity`: Missing lift max for the user

#### POST /prescriptions/resolve-batch

Resolve multiple prescriptions at once.

**Auth**: Authenticated

**Request Body**:
```json
{
  "prescriptionIds": ["uuid1", "uuid2"],
  "userId": "user-uuid"
}
```

**Response** `200 OK`:
```json
{
  "results": [
    {
      "prescriptionId": "uuid1",
      "status": "success",
      "resolved": { ... }
    },
    {
      "prescriptionId": "uuid2",
      "status": "error",
      "error": "missing lift max for squat (TRAINING_MAX)"
    }
  ]
}
```

---

### Days

Manage training days.

#### GET /days

List all days.

**Auth**: Authenticated

**Query Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `page` | int | Page number |
| `pageSize` | int | Items per page |
| `sortBy` | string | "name" or "created_at" |
| `sortOrder` | string | "asc" or "desc" |
| `program_id` | string | Filter by program ID |

**Response** `200 OK`:
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "Squat Day",
      "slug": "squat-day",
      "metadata": {"focus": "legs"},
      "programId": "program-uuid",
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "page": 1,
  "pageSize": 20,
  "totalItems": 1,
  "totalPages": 1
}
```

#### GET /days/{id}

Get a day by ID with its prescriptions.

**Auth**: Authenticated

**Response** `200 OK`:
```json
{
  "id": "uuid",
  "name": "Squat Day",
  "slug": "squat-day",
  "metadata": {"focus": "legs"},
  "programId": "program-uuid",
  "prescriptions": [
    {
      "id": "day-prescription-uuid",
      "prescriptionId": "prescription-uuid",
      "order": 1,
      "createdAt": "2024-01-01T00:00:00Z"
    }
  ],
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### GET /days/by-slug/{slug}

Get a day by slug.

**Auth**: Authenticated

**Query Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `program_id` | string | Optional: scope to specific program |

**Response**: Same as GET /days/{id}

#### POST /days

Create a new day.

**Auth**: Admin

**Request Body**:
```json
{
  "name": "Squat Day",
  "slug": "squat-day",
  "metadata": {"focus": "legs"},
  "programId": "program-uuid"
}
```

**Response** `201 Created`: Day object (without prescriptions)

#### PUT /days/{id}

Update a day.

**Auth**: Admin

**Request Body**:
```json
{
  "name": "Heavy Squat Day",
  "slug": "heavy-squat-day",
  "metadata": {"focus": "strength"},
  "clearMetadata": false,
  "programId": "program-uuid",
  "clearProgramId": false
}
```

**Response** `200 OK`: Updated Day object (without prescriptions)

#### DELETE /days/{id}

Delete a day.

**Auth**: Admin

**Response** `204 No Content`

**Errors**:
- `409 Conflict`: Day is used in one or more weeks

#### POST /days/{id}/prescriptions

Add a prescription to a day.

**Auth**: Admin

**Request Body**:
```json
{
  "prescriptionId": "prescription-uuid",
  "order": 1
}
```

**Response** `201 Created`: Day with prescriptions

#### DELETE /days/{id}/prescriptions/{prescriptionId}

Remove a prescription from a day.

**Auth**: Admin

**Response** `204 No Content`

#### PUT /days/{id}/prescriptions/reorder

Reorder prescriptions within a day.

**Auth**: Admin

**Request Body**:
```json
{
  "prescriptionIds": ["uuid1", "uuid2", "uuid3"]
}
```

**Response** `200 OK`: Day with reordered prescriptions

---

### Weeks

Manage training weeks.

#### GET /weeks

List all weeks.

**Auth**: Authenticated

**Query Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `page` | int | Page number |
| `pageSize` | int | Items per page |
| `sortBy` | string | "week_number" or "created_at" |
| `sortOrder` | string | "asc" or "desc" |

**Response** `200 OK`: Paginated list of Week objects

#### GET /weeks/{id}

Get a week by ID.

**Auth**: Authenticated

**Response** `200 OK`:
```json
{
  "id": "uuid",
  "cycleId": "cycle-uuid",
  "weekNumber": 1,
  "name": "Week 1",
  "days": [
    {
      "id": "week-day-uuid",
      "dayId": "day-uuid",
      "position": 0
    }
  ],
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### POST /weeks

Create a new week.

**Auth**: Admin

**Request Body**:
```json
{
  "cycleId": "cycle-uuid",
  "weekNumber": 1,
  "name": "Week 1"
}
```

**Response** `201 Created`: Week object

#### PUT /weeks/{id}

Update a week.

**Auth**: Admin

**Response** `200 OK`: Updated Week object

#### DELETE /weeks/{id}

Delete a week.

**Auth**: Admin

**Response** `204 No Content`

#### POST /weeks/{id}/days

Add a day to a week.

**Auth**: Admin

**Request Body**:
```json
{
  "dayId": "day-uuid",
  "position": 0
}
```

**Response** `201 Created`: Updated Week with days

#### DELETE /weeks/{id}/days/{dayId}

Remove a day from a week.

**Auth**: Admin

**Response** `204 No Content`

---

### Cycles

Manage training cycles.

#### GET /cycles

List all cycles.

**Auth**: Authenticated

**Response** `200 OK`: Paginated list of Cycle objects

#### GET /cycles/{id}

Get a cycle by ID.

**Auth**: Authenticated

**Response** `200 OK`:
```json
{
  "id": "uuid",
  "name": "4 Week Cycle",
  "lengthWeeks": 4,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### POST /cycles

Create a new cycle.

**Auth**: Admin

**Request Body**:
```json
{
  "name": "4 Week Cycle",
  "lengthWeeks": 4
}
```

**Response** `201 Created`: Cycle object

#### PUT /cycles/{id}

Update a cycle.

**Auth**: Admin

**Response** `200 OK`: Updated Cycle object

#### DELETE /cycles/{id}

Delete a cycle.

**Auth**: Admin

**Response** `204 No Content`

---

### Weekly Lookups

Manage weekly lookup tables (varying parameters by week).

#### GET /weekly-lookups

List all weekly lookups.

**Auth**: Authenticated

**Response** `200 OK`: Paginated list of WeeklyLookup objects

#### GET /weekly-lookups/{id}

Get a weekly lookup by ID.

**Auth**: Authenticated

**Response** `200 OK`:
```json
{
  "id": "uuid",
  "name": "5/3/1 Percentages",
  "entries": {
    "1": 65,
    "2": 70,
    "3": 75,
    "4": 40
  },
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### POST /weekly-lookups

Create a new weekly lookup.

**Auth**: Admin

**Request Body**:
```json
{
  "name": "5/3/1 Percentages",
  "entries": {
    "1": 65,
    "2": 70,
    "3": 75,
    "4": 40
  }
}
```

**Response** `201 Created`: WeeklyLookup object

#### PUT /weekly-lookups/{id}

Update a weekly lookup.

**Auth**: Admin

**Response** `200 OK`: Updated WeeklyLookup object

#### DELETE /weekly-lookups/{id}

Delete a weekly lookup.

**Auth**: Admin

**Response** `204 No Content`

---

### Daily Lookups

Manage daily lookup tables (varying parameters by day).

#### GET /daily-lookups

List all daily lookups.

**Auth**: Authenticated

**Response** `200 OK`: Paginated list of DailyLookup objects

#### GET /daily-lookups/{id}

Get a daily lookup by ID.

**Auth**: Authenticated

**Response** `200 OK`:
```json
{
  "id": "uuid",
  "name": "Heavy/Light/Medium",
  "entries": {
    "monday": 100,
    "wednesday": 80,
    "friday": 90
  },
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### POST /daily-lookups

Create a new daily lookup.

**Auth**: Admin

**Request Body**:
```json
{
  "name": "Heavy/Light/Medium",
  "entries": {
    "monday": 100,
    "wednesday": 80,
    "friday": 90
  }
}
```

**Response** `201 Created`: DailyLookup object

#### PUT /daily-lookups/{id}

Update a daily lookup.

**Auth**: Admin

**Response** `200 OK`: Updated DailyLookup object

#### DELETE /daily-lookups/{id}

Delete a daily lookup.

**Auth**: Admin

**Response** `204 No Content`

---

### Programs

Manage training programs.

#### GET /programs

List all programs.

**Auth**: Authenticated

**Query Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `page` | int | Page number |
| `pageSize` | int | Items per page |
| `sortBy` | string | "name" or "created_at" |
| `sortOrder` | string | "asc" or "desc" |

**Response** `200 OK`:
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "Wendler 5/3/1 BBB",
      "slug": "wendler-531-bbb",
      "description": "Boring But Big variant",
      "cycleId": "cycle-uuid",
      "weeklyLookupId": "lookup-uuid",
      "dailyLookupId": null,
      "defaultRounding": 5.0,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "page": 1,
  "pageSize": 20,
  "totalItems": 1,
  "totalPages": 1
}
```

#### GET /programs/{id}

Get a program by ID with embedded cycle details.

**Auth**: Authenticated

**Response** `200 OK`:
```json
{
  "id": "uuid",
  "name": "Wendler 5/3/1 BBB",
  "slug": "wendler-531-bbb",
  "description": "Boring But Big variant",
  "cycle": {
    "id": "cycle-uuid",
    "name": "4 Week Cycle",
    "lengthWeeks": 4,
    "weeks": [
      {"id": "week-uuid", "weekNumber": 1}
    ]
  },
  "weeklyLookup": {
    "id": "lookup-uuid",
    "name": "5/3/1 Percentages"
  },
  "dailyLookup": null,
  "defaultRounding": 5.0,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### POST /programs

Create a new program.

**Auth**: Admin

**Request Body**:
```json
{
  "name": "Wendler 5/3/1 BBB",
  "slug": "wendler-531-bbb",
  "description": "Boring But Big variant",
  "cycleId": "cycle-uuid",
  "weeklyLookupId": "lookup-uuid",
  "dailyLookupId": null,
  "defaultRounding": 5.0
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Program name |
| `slug` | string | Yes | URL-friendly identifier (must be unique) |
| `description` | string | No | Program description |
| `cycleId` | string | Yes | Associated cycle ID |
| `weeklyLookupId` | string | No | Weekly lookup table ID |
| `dailyLookupId` | string | No | Daily lookup table ID |
| `defaultRounding` | float | No | Default weight rounding (e.g., 5.0 for 5lb plates) |

**Response** `201 Created`: Program object (list format)

#### PUT /programs/{id}

Update a program.

**Auth**: Admin

**Response** `200 OK`: Updated Program object (list format)

#### DELETE /programs/{id}

Delete a program.

**Auth**: Admin

**Response** `204 No Content`

**Errors**:
- `409 Conflict`: Users are enrolled in this program

---

### Progressions

Manage progression rules (how to increase weights over time).

#### GET /progressions

List all progressions.

**Auth**: Authenticated

**Query Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `page` | int | Page number |
| `pageSize` | int | Items per page |
| `type` | string | Filter by type: "LINEAR" or "CYCLE" |

**Response** `200 OK`:
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "Linear +5lb",
      "type": "LINEAR",
      "parameters": {
        "increment": 5.0,
        "maxType": "TRAINING_MAX",
        "triggerType": "AFTER_SESSION"
      },
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "page": 1,
  "pageSize": 20,
  "totalItems": 1,
  "totalPages": 1
}
```

#### GET /progressions/{id}

Get a progression by ID.

**Auth**: Authenticated

**Response** `200 OK`: Progression object

#### POST /progressions

Create a new progression.

**Auth**: Admin

**Request Body** (LINEAR type):
```json
{
  "name": "Linear +5lb",
  "type": "LINEAR",
  "parameters": {
    "increment": 5.0,
    "maxType": "TRAINING_MAX",
    "triggerType": "AFTER_SESSION"
  }
}
```

**Request Body** (CYCLE type):
```json
{
  "name": "Cycle +5lb",
  "type": "CYCLE",
  "parameters": {
    "increment": 5.0,
    "maxType": "TRAINING_MAX"
  }
}
```

**Progression Types**:

| Type | Description | Parameters |
|------|-------------|------------|
| `LINEAR` | Add weight after session/week | `increment`, `maxType`, `triggerType` |
| `CYCLE` | Add weight at end of cycle | `increment`, `maxType` |

**Parameter Values**:

- `maxType`: "ONE_RM" or "TRAINING_MAX"
- `triggerType` (LINEAR only): "AFTER_SESSION" or "AFTER_WEEK"

**Response** `201 Created`: Progression object

#### PUT /progressions/{id}

Update a progression.

**Auth**: Admin

**Response** `200 OK`: Updated Progression object

#### DELETE /progressions/{id}

Delete a progression.

**Auth**: Admin

**Response** `204 No Content`

**Errors**:
- `409 Conflict`: Progression is referenced by program progressions

---

### Program Progressions

Configure which progressions apply to which programs/lifts.

#### GET /programs/{programId}/progressions

List progression configurations for a program.

**Auth**: Authenticated

**Response** `200 OK`: List of ProgramProgression objects

#### GET /programs/{programId}/progressions/{configId}

Get a specific progression configuration.

**Auth**: Authenticated

**Response** `200 OK`:
```json
{
  "id": "config-uuid",
  "programId": "program-uuid",
  "progressionId": "progression-uuid",
  "liftId": "lift-uuid",
  "priority": 1,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### POST /programs/{programId}/progressions

Create a new progression configuration.

**Auth**: Admin

**Request Body**:
```json
{
  "progressionId": "progression-uuid",
  "liftId": "lift-uuid",
  "priority": 1
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `progressionId` | string | Yes | Progression rule ID |
| `liftId` | string | No | Specific lift (null = all lifts) |
| `priority` | int | No | Priority when multiple progressions apply |

**Response** `201 Created`: ProgramProgression object

#### PUT /programs/{programId}/progressions/{configId}

Update a progression configuration.

**Auth**: Admin

**Response** `200 OK`: Updated ProgramProgression object

#### DELETE /programs/{programId}/progressions/{configId}

Delete a progression configuration.

**Auth**: Admin

**Response** `204 No Content`

---

### User Program Enrollment

Manage user enrollment in programs.

#### GET /users/{userId}/program

Get a user's current program enrollment.

**Auth**: Owner/Admin

**Response** `200 OK`:
```json
{
  "id": "enrollment-uuid",
  "userId": "user-uuid",
  "program": {
    "id": "program-uuid",
    "name": "Wendler 5/3/1 BBB",
    "slug": "wendler-531-bbb",
    "description": "Boring But Big variant",
    "cycleLengthWeeks": 4
  },
  "state": {
    "currentWeek": 1,
    "currentCycleIteration": 1,
    "currentDayIndex": 0
  },
  "enrolledAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

**Errors**:
- `404 Not Found`: User is not enrolled in any program

#### POST /users/{userId}/program

Enroll a user in a program. If already enrolled, replaces the existing enrollment.

**Auth**: Owner/Admin

**Request Body**:
```json
{
  "programId": "program-uuid"
}
```

**Response** `201 Created`: Enrollment object

#### DELETE /users/{userId}/program

Unenroll a user from their current program.

**Auth**: Owner/Admin

**Response** `204 No Content`

---

### State Advancement

Advance a user's program state (move to next day/week).

#### POST /users/{userId}/program-state/advance

Advance the user's program state.

**Auth**: Owner/Admin

**Request Body**:
```json
{
  "advanceType": "day"
}
```

| Field | Type | Description |
|-------|------|-------------|
| `advanceType` | string | "day" or "week" |

**Response** `200 OK`: Updated enrollment with new state

---

### Workout Generation

Generate workouts based on user's program and state.

#### GET /users/{userId}/workout

Generate the current workout for a user.

**Auth**: Owner/Admin

**Query Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `date` | string | Override workout date (YYYY-MM-DD) |
| `weekNumber` | int | Override week number |
| `daySlug` | string | Override day slug |

**Response** `200 OK`:
```json
{
  "userId": "user-uuid",
  "programId": "program-uuid",
  "cycleIteration": 1,
  "weekNumber": 1,
  "daySlug": "squat-day",
  "date": "2024-01-15",
  "exercises": [
    {
      "prescriptionId": "prescription-uuid",
      "lift": {
        "id": "lift-uuid",
        "name": "Squat",
        "slug": "squat"
      },
      "sets": [
        {
          "setNumber": 1,
          "weight": 265.0,
          "targetReps": 5,
          "isWorkSet": true
        }
      ],
      "notes": "Focus on depth",
      "restSeconds": 180
    }
  ]
}
```

**Errors**:
- `404 Not Found`: User not enrolled in a program
- `400 Bad Request`: Missing lift max (set up training maxes first)

#### GET /users/{userId}/workout/preview

Preview a workout for a specific week/day without state advancement.

**Auth**: Owner/Admin

**Query Parameters** (required):
| Parameter | Type | Description |
|-----------|------|-------------|
| `week` | int | Week number to preview |
| `day` | string | Day slug to preview |

**Response** `200 OK`: Same as GET /users/{userId}/workout

---

### Progression History

Query a user's progression history.

#### GET /users/{userId}/progression-history

List progression history entries for a user.

**Auth**: Owner/Admin

**Query Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `page` | int | Page number |
| `pageSize` | int | Items per page |
| `lift_id` | string | Filter by lift ID |
| `progression_id` | string | Filter by progression ID |

**Response** `200 OK`:
```json
{
  "data": [
    {
      "id": "history-uuid",
      "userId": "user-uuid",
      "liftId": "lift-uuid",
      "progressionId": "progression-uuid",
      "previousValue": 315.0,
      "newValue": 320.0,
      "delta": 5.0,
      "appliedAt": "2024-01-15T10:30:00Z",
      "triggeredBy": "AFTER_SESSION",
      "cycleIteration": 1,
      "weekNumber": 1
    }
  ],
  "page": 1,
  "pageSize": 20,
  "totalItems": 1,
  "totalPages": 1
}
```

---

### Manual Progression Trigger

Manually trigger a progression for a user.

#### POST /users/{userId}/progressions/trigger

Manually apply a progression.

**Auth**: Owner/Admin

**Request Body**:
```json
{
  "progressionId": "progression-uuid",
  "liftId": "lift-uuid",
  "force": false
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `progressionId` | string | Yes | Progression to apply |
| `liftId` | string | No | Specific lift (null = all configured lifts) |
| `force` | bool | No | Force apply even if already applied this period |

**Response** `200 OK`:
```json
{
  "results": [
    {
      "progressionId": "progression-uuid",
      "liftId": "lift-uuid",
      "applied": true,
      "skipped": false,
      "result": {
        "previousValue": 315.0,
        "newValue": 320.0,
        "delta": 5.0,
        "maxType": "TRAINING_MAX",
        "appliedAt": "2024-01-15T10:30:00Z"
      }
    }
  ],
  "totalApplied": 1,
  "totalSkipped": 0,
  "totalErrors": 0
}
```

**Errors**:
- `404 Not Found`: Progression or lift not found
- `400 Bad Request`: User not enrolled / no applicable progressions

---

## Error Reference

| Error Message | Status | Cause |
|---------------|--------|-------|
| `Authentication required` | 401 | Missing or invalid auth headers |
| `Admin privileges required` | 403 | Admin-only endpoint accessed by non-admin |
| `Access denied: you do not have permission` | 403 | Owner-only resource accessed by non-owner |
| `not found` | 404 | Resource does not exist |
| `validation failed` | 400 | Request body validation errors |
| `slug already exists` | 409 | Duplicate slug |
| `cannot delete: it is referenced` | 409 | Foreign key constraint violation |
| `missing lift max` | 422 | Required lift max not set up for user |

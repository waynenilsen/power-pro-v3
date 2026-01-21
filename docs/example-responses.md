# PowerPro API Example Responses

This document provides complete example responses for all PowerPro API endpoints. These examples show the exact JSON structure returned by the API, including success responses and common error responses.

For example requests, see [`docs/example-requests.md`](./example-requests.md).

## Table of Contents

- [Response Formats](#response-formats)
  - [Paginated Response Format](#paginated-response-format)
  - [Error Response Format](#error-response-format)
  - [Response with Warnings Format](#response-with-warnings-format)
- [HTTP Status Codes](#http-status-codes)
- [Common Error Responses](#common-error-responses)
- [Endpoint Responses](#endpoint-responses)
  - [Health Check](#health-check)
  - [Lifts](#lifts)
  - [Lift Maxes](#lift-maxes)
  - [Prescriptions](#prescriptions)
  - [Days](#days)
  - [Weeks](#weeks)
  - [Cycles](#cycles)
  - [Weekly Lookups](#weekly-lookups)
  - [Daily Lookups](#daily-lookups)
  - [Programs](#programs)
  - [Progressions](#progressions)
  - [Program Progressions](#program-progressions)
  - [User Program Enrollment](#user-program-enrollment)
  - [State Advancement](#state-advancement)
  - [Workout Generation](#workout-generation)
  - [Progression History](#progression-history)
  - [Manual Progression Trigger](#manual-progression-trigger)

---

## Response Formats

### Paginated Response Format

All list endpoints return data in this paginated format:

```json
{
  "data": [...],
  "page": 1,
  "pageSize": 20,
  "totalItems": 100,
  "totalPages": 5
}
```

| Field | Type | Description |
|-------|------|-------------|
| `data` | array | Array of resource objects |
| `page` | int | Current page number (1-indexed) |
| `pageSize` | int | Number of items per page |
| `totalItems` | int | Total number of items across all pages |
| `totalPages` | int | Total number of pages |

### Error Response Format

All error responses follow this format:

```json
{
  "error": "error message",
  "details": ["detail 1", "detail 2"]
}
```

| Field | Type | Description |
|-------|------|-------------|
| `error` | string | Human-readable error message |
| `details` | array | Optional array of specific error details |

### Response with Warnings Format

Some operations succeed but include warnings:

```json
{
  "data": {...},
  "warnings": ["warning message 1", "warning message 2"]
}
```

---

## HTTP Status Codes

| Code | Description | When Used |
|------|-------------|-----------|
| `200 OK` | Success | GET (read), PUT (update) |
| `201 Created` | Resource created | POST (create) |
| `204 No Content` | Success with no body | DELETE |
| `400 Bad Request` | Invalid request | Validation errors, malformed JSON |
| `401 Unauthorized` | Authentication required | Missing or invalid auth headers |
| `403 Forbidden` | Access denied | Insufficient permissions |
| `404 Not Found` | Resource not found | ID/slug doesn't exist |
| `409 Conflict` | Conflict | Duplicate slug, FK constraint violation |
| `422 Unprocessable Entity` | Business rule violation | Missing lift max, etc. |
| `500 Internal Server Error` | Server error | Unexpected errors |

---

## Common Error Responses

### 400 Bad Request - Validation Error

**When**: Request body fails validation

```http
HTTP/1.1 400 Bad Request
Content-Type: application/json
```

```json
{
  "error": "validation failed",
  "details": [
    "name is required",
    "value must be positive"
  ]
}
```

### 400 Bad Request - Invalid Request Body

**When**: Request body is malformed JSON

```http
HTTP/1.1 400 Bad Request
Content-Type: application/json
```

```json
{
  "error": "invalid request body"
}
```

### 400 Bad Request - Missing Required Parameter

**When**: Required query parameter is missing

```http
HTTP/1.1 400 Bad Request
Content-Type: application/json
```

```json
{
  "error": "week: missing required parameter"
}
```

### 401 Unauthorized

**When**: No authentication headers provided

```http
HTTP/1.1 401 Unauthorized
Content-Type: application/json
```

```json
{
  "error": "Authentication required"
}
```

### 403 Forbidden - Admin Required

**When**: Non-admin tries to access admin-only endpoint

```http
HTTP/1.1 403 Forbidden
Content-Type: application/json
```

```json
{
  "error": "Admin privileges required"
}
```

### 403 Forbidden - Access Denied

**When**: User tries to access another user's resources

```http
HTTP/1.1 403 Forbidden
Content-Type: application/json
```

```json
{
  "error": "you can only view your own enrollment"
}
```

### 404 Not Found

**When**: Resource doesn't exist

```http
HTTP/1.1 404 Not Found
Content-Type: application/json
```

```json
{
  "error": "lift not found"
}
```

### 409 Conflict - Duplicate Slug

**When**: Creating or updating with a slug that already exists

```http
HTTP/1.1 409 Conflict
Content-Type: application/json
```

```json
{
  "error": "slug already exists"
}
```

### 409 Conflict - Foreign Key Constraint

**When**: Trying to delete a resource that is referenced by other records

```http
HTTP/1.1 409 Conflict
Content-Type: application/json
```

```json
{
  "error": "cannot delete lift: it is referenced by other lifts as a parent"
}
```

### 422 Unprocessable Entity - Missing Lift Max

**When**: Trying to generate a workout without required lift maxes

```http
HTTP/1.1 422 Unprocessable Entity
Content-Type: application/json
```

```json
{
  "error": "missing lift max: set up your training maxes to generate workouts",
  "details": [
    "missing lift max for squat (TRAINING_MAX)"
  ]
}
```

### 500 Internal Server Error

**When**: Unexpected server error (details logged server-side)

```http
HTTP/1.1 500 Internal Server Error
Content-Type: application/json
```

```json
{
  "error": "internal server error"
}
```

---

## Endpoint Responses

### Health Check

#### GET /health

**Success Response** `200 OK`:

```json
{
  "status": "ok"
}
```

---

### Lifts

#### GET /lifts - List Lifts

**Success Response** `200 OK`:

```json
{
  "data": [
    {
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "name": "Squat",
      "slug": "squat",
      "isCompetitionLift": true,
      "parentLiftId": null,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    },
    {
      "id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
      "name": "Bench Press",
      "slug": "bench-press",
      "isCompetitionLift": true,
      "parentLiftId": null,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    },
    {
      "id": "c3d4e5f6-a7b8-9012-cdef-234567890123",
      "name": "Pause Squat",
      "slug": "pause-squat",
      "isCompetitionLift": false,
      "parentLiftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "createdAt": "2024-01-02T00:00:00Z",
      "updatedAt": "2024-01-02T00:00:00Z"
    }
  ],
  "page": 1,
  "pageSize": 20,
  "totalItems": 3,
  "totalPages": 1
}
```

#### GET /lifts/{id} - Get Lift by ID

**Success Response** `200 OK`:

```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "name": "Squat",
  "slug": "squat",
  "isCompetitionLift": true,
  "parentLiftId": null,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

**Error Response** `404 Not Found`:

```json
{
  "error": "lift not found"
}
```

#### GET /lifts/by-slug/{slug} - Get Lift by Slug

**Success Response** `200 OK`:

```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "name": "Squat",
  "slug": "squat",
  "isCompetitionLift": true,
  "parentLiftId": null,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### POST /lifts - Create Lift

**Success Response** `201 Created`:

```json
{
  "id": "d4e5f6a7-b8c9-0123-defg-345678901234",
  "name": "Deadlift",
  "slug": "deadlift",
  "isCompetitionLift": true,
  "parentLiftId": null,
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

**Error Response** `400 Bad Request` (validation):

```json
{
  "error": "validation failed",
  "details": [
    "name is required"
  ]
}
```

**Error Response** `409 Conflict` (duplicate slug):

```json
{
  "error": "slug already exists"
}
```

#### PUT /lifts/{id} - Update Lift

**Success Response** `200 OK`:

```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "name": "Back Squat",
  "slug": "back-squat",
  "isCompetitionLift": true,
  "parentLiftId": null,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### DELETE /lifts/{id} - Delete Lift

**Success Response** `204 No Content`:

*(No response body)*

**Error Response** `409 Conflict`:

```json
{
  "error": "cannot delete lift: it is referenced by other lifts as a parent"
}
```

---

### Lift Maxes

#### GET /users/{userId}/lift-maxes - List User's Lift Maxes

**Success Response** `200 OK`:

```json
{
  "data": [
    {
      "id": "e5f6a7b8-c9d0-1234-efgh-456789012345",
      "userId": "550e8400-e29b-41d4-a716-446655440000",
      "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "type": "TRAINING_MAX",
      "value": 315.0,
      "effectiveDate": "2024-01-01T00:00:00Z",
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    },
    {
      "id": "f6a7b8c9-d0e1-2345-fghi-567890123456",
      "userId": "550e8400-e29b-41d4-a716-446655440000",
      "liftId": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
      "type": "TRAINING_MAX",
      "value": 225.0,
      "effectiveDate": "2024-01-01T00:00:00Z",
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "page": 1,
  "pageSize": 20,
  "totalItems": 2,
  "totalPages": 1
}
```

#### GET /users/{userId}/lift-maxes/current - Get Current Lift Max

**Success Response** `200 OK`:

```json
{
  "id": "e5f6a7b8-c9d0-1234-efgh-456789012345",
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "type": "TRAINING_MAX",
  "value": 315.0,
  "effectiveDate": "2024-01-01T00:00:00Z",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

**Error Response** `404 Not Found`:

```json
{
  "error": "lift max not found"
}
```

#### GET /lift-maxes/{id} - Get Lift Max by ID

**Success Response** `200 OK`:

```json
{
  "id": "e5f6a7b8-c9d0-1234-efgh-456789012345",
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "type": "TRAINING_MAX",
  "value": 315.0,
  "effectiveDate": "2024-01-01T00:00:00Z",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### GET /lift-maxes/{id}/convert - Convert Lift Max

**Success Response** `200 OK`:

```json
{
  "originalValue": 350.0,
  "originalType": "ONE_RM",
  "convertedValue": 315.0,
  "convertedType": "TRAINING_MAX",
  "percentage": 90.0
}
```

#### POST /users/{userId}/lift-maxes - Create Lift Max

**Success Response** `201 Created`:

```json
{
  "id": "g7h8i9j0-k1l2-3456-ijkl-678901234567",
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "liftId": "d4e5f6a7-b8c9-0123-defg-345678901234",
  "type": "TRAINING_MAX",
  "value": 405.0,
  "effectiveDate": "2024-01-15T00:00:00Z",
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

**Error Response** `400 Bad Request`:

```json
{
  "error": "validation failed",
  "details": [
    "value must be positive"
  ]
}
```

#### PUT /lift-maxes/{id} - Update Lift Max

**Success Response** `200 OK`:

```json
{
  "id": "e5f6a7b8-c9d0-1234-efgh-456789012345",
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "type": "TRAINING_MAX",
  "value": 320.0,
  "effectiveDate": "2024-01-15T00:00:00Z",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### DELETE /lift-maxes/{id} - Delete Lift Max

**Success Response** `204 No Content`:

*(No response body)*

---

### Prescriptions

#### GET /prescriptions - List Prescriptions

**Success Response** `200 OK`:

```json
{
  "data": [
    {
      "id": "h8i9j0k1-l2m3-4567-mnop-789012345678",
      "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "loadStrategy": {
        "type": "PERCENT_OF",
        "maxType": "TRAINING_MAX",
        "percentage": 85.0,
        "lookupKey": null,
        "roundTo": 5.0
      },
      "setScheme": {
        "type": "FIXED",
        "sets": 5,
        "reps": 5,
        "isAmrap": false
      },
      "order": 1,
      "notes": "Focus on depth and controlled descent",
      "restSeconds": 180,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    },
    {
      "id": "i9j0k1l2-m3n4-5678-nopq-890123456789",
      "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "loadStrategy": {
        "type": "PERCENT_OF",
        "maxType": "TRAINING_MAX",
        "percentage": 95.0,
        "lookupKey": "week",
        "roundTo": 5.0
      },
      "setScheme": {
        "type": "FIXED",
        "sets": 1,
        "reps": 1,
        "isAmrap": true
      },
      "order": 2,
      "notes": "Top set - as many reps as possible",
      "restSeconds": 300,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "page": 1,
  "pageSize": 20,
  "totalItems": 2,
  "totalPages": 1
}
```

#### GET /prescriptions/{id} - Get Prescription by ID

**Success Response** `200 OK`:

```json
{
  "id": "h8i9j0k1-l2m3-4567-mnop-789012345678",
  "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "loadStrategy": {
    "type": "PERCENT_OF",
    "maxType": "TRAINING_MAX",
    "percentage": 85.0,
    "lookupKey": null,
    "roundTo": 5.0
  },
  "setScheme": {
    "type": "FIXED",
    "sets": 5,
    "reps": 5,
    "isAmrap": false
  },
  "order": 1,
  "notes": "Focus on depth and controlled descent",
  "restSeconds": 180,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### POST /prescriptions - Create Prescription

**Success Response** `201 Created`:

```json
{
  "id": "j0k1l2m3-n4o5-6789-opqr-901234567890",
  "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "loadStrategy": {
    "type": "PERCENT_OF",
    "maxType": "TRAINING_MAX",
    "percentage": 50.0,
    "lookupKey": null,
    "roundTo": 5.0
  },
  "setScheme": {
    "type": "FIXED",
    "sets": 5,
    "reps": 10,
    "isAmrap": false
  },
  "order": 3,
  "notes": "BBB volume work - focus on quality reps",
  "restSeconds": 90,
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### Prescription with RAMP Set Scheme

**Example Response**:

```json
{
  "id": "k1l2m3n4-o5p6-7890-pqrs-012345678901",
  "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "loadStrategy": {
    "type": "PERCENT_OF",
    "maxType": "TRAINING_MAX",
    "percentage": 100.0,
    "lookupKey": null,
    "roundTo": null
  },
  "setScheme": {
    "type": "RAMP",
    "sets": [
      {"percentage": 40, "reps": 5},
      {"percentage": 50, "reps": 5},
      {"percentage": 60, "reps": 3},
      {"percentage": 70, "reps": 2}
    ]
  },
  "order": 0,
  "notes": "Warmup - build up to working weight",
  "restSeconds": null,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### PUT /prescriptions/{id} - Update Prescription

**Success Response** `200 OK`:

```json
{
  "id": "h8i9j0k1-l2m3-4567-mnop-789012345678",
  "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "loadStrategy": {
    "type": "PERCENT_OF",
    "maxType": "TRAINING_MAX",
    "percentage": 90.0,
    "lookupKey": null,
    "roundTo": 5.0
  },
  "setScheme": {
    "type": "FIXED",
    "sets": 5,
    "reps": 5,
    "isAmrap": false
  },
  "order": 1,
  "notes": "Pause at the bottom for 2 seconds",
  "restSeconds": 240,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### DELETE /prescriptions/{id} - Delete Prescription

**Success Response** `204 No Content`:

*(No response body)*

#### POST /prescriptions/{id}/resolve - Resolve Prescription

**Success Response** `200 OK`:

```json
{
  "prescriptionId": "h8i9j0k1-l2m3-4567-mnop-789012345678",
  "lift": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "name": "Squat",
    "slug": "squat"
  },
  "sets": [
    {
      "setNumber": 1,
      "weight": 265.0,
      "targetReps": 5,
      "isWorkSet": true
    },
    {
      "setNumber": 2,
      "weight": 265.0,
      "targetReps": 5,
      "isWorkSet": true
    },
    {
      "setNumber": 3,
      "weight": 265.0,
      "targetReps": 5,
      "isWorkSet": true
    },
    {
      "setNumber": 4,
      "weight": 265.0,
      "targetReps": 5,
      "isWorkSet": true
    },
    {
      "setNumber": 5,
      "weight": 265.0,
      "targetReps": 5,
      "isWorkSet": true
    }
  ],
  "notes": "Focus on depth and controlled descent",
  "restSeconds": 180
}
```

**Error Response** `422 Unprocessable Entity`:

```json
{
  "error": "missing lift max for squat (TRAINING_MAX)"
}
```

#### POST /prescriptions/resolve-batch - Resolve Multiple Prescriptions

**Success Response** `200 OK`:

```json
{
  "results": [
    {
      "prescriptionId": "h8i9j0k1-l2m3-4567-mnop-789012345678",
      "status": "success",
      "resolved": {
        "prescriptionId": "h8i9j0k1-l2m3-4567-mnop-789012345678",
        "lift": {
          "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
          "name": "Squat",
          "slug": "squat"
        },
        "sets": [
          {"setNumber": 1, "weight": 265.0, "targetReps": 5, "isWorkSet": true},
          {"setNumber": 2, "weight": 265.0, "targetReps": 5, "isWorkSet": true},
          {"setNumber": 3, "weight": 265.0, "targetReps": 5, "isWorkSet": true},
          {"setNumber": 4, "weight": 265.0, "targetReps": 5, "isWorkSet": true},
          {"setNumber": 5, "weight": 265.0, "targetReps": 5, "isWorkSet": true}
        ],
        "notes": "Focus on depth",
        "restSeconds": 180
      }
    },
    {
      "prescriptionId": "i9j0k1l2-m3n4-5678-nopq-890123456789",
      "status": "error",
      "error": "missing lift max for bench-press (TRAINING_MAX)"
    }
  ]
}
```

---

### Days

#### GET /days - List Days

**Success Response** `200 OK`:

```json
{
  "data": [
    {
      "id": "l2m3n4o5-p6q7-8901-qrst-123456789012",
      "name": "Squat Day",
      "slug": "squat-day",
      "metadata": {
        "focus": "lower body",
        "primaryLift": "squat"
      },
      "programId": "m3n4o5p6-q7r8-9012-rstu-234567890123",
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    },
    {
      "id": "n4o5p6q7-r8s9-0123-stuv-345678901234",
      "name": "Bench Day",
      "slug": "bench-day",
      "metadata": {
        "focus": "upper body",
        "primaryLift": "bench"
      },
      "programId": "m3n4o5p6-q7r8-9012-rstu-234567890123",
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "page": 1,
  "pageSize": 20,
  "totalItems": 2,
  "totalPages": 1
}
```

#### GET /days/{id} - Get Day by ID (with prescriptions)

**Success Response** `200 OK`:

```json
{
  "id": "l2m3n4o5-p6q7-8901-qrst-123456789012",
  "name": "Squat Day",
  "slug": "squat-day",
  "metadata": {
    "focus": "lower body",
    "primaryLift": "squat"
  },
  "programId": "m3n4o5p6-q7r8-9012-rstu-234567890123",
  "prescriptions": [
    {
      "id": "o5p6q7r8-s9t0-1234-uvwx-456789012345",
      "prescriptionId": "h8i9j0k1-l2m3-4567-mnop-789012345678",
      "order": 1,
      "createdAt": "2024-01-01T00:00:00Z"
    },
    {
      "id": "p6q7r8s9-t0u1-2345-vwxy-567890123456",
      "prescriptionId": "i9j0k1l2-m3n4-5678-nopq-890123456789",
      "order": 2,
      "createdAt": "2024-01-01T00:00:00Z"
    }
  ],
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### GET /days/by-slug/{slug} - Get Day by Slug

**Success Response** `200 OK`:

```json
{
  "id": "l2m3n4o5-p6q7-8901-qrst-123456789012",
  "name": "Squat Day",
  "slug": "squat-day",
  "metadata": {
    "focus": "lower body",
    "primaryLift": "squat"
  },
  "programId": "m3n4o5p6-q7r8-9012-rstu-234567890123",
  "prescriptions": [
    {
      "id": "o5p6q7r8-s9t0-1234-uvwx-456789012345",
      "prescriptionId": "h8i9j0k1-l2m3-4567-mnop-789012345678",
      "order": 1,
      "createdAt": "2024-01-01T00:00:00Z"
    }
  ],
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### POST /days - Create Day

**Success Response** `201 Created`:

```json
{
  "id": "q7r8s9t0-u1v2-3456-wxyz-678901234567",
  "name": "Deadlift Day",
  "slug": "deadlift-day",
  "metadata": {
    "focus": "posterior chain",
    "primaryLift": "deadlift"
  },
  "programId": "m3n4o5p6-q7r8-9012-rstu-234567890123",
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### PUT /days/{id} - Update Day

**Success Response** `200 OK`:

```json
{
  "id": "l2m3n4o5-p6q7-8901-qrst-123456789012",
  "name": "Heavy Squat Day",
  "slug": "heavy-squat-day",
  "metadata": {
    "focus": "max effort",
    "primaryLift": "squat",
    "intensity": "high"
  },
  "programId": "m3n4o5p6-q7r8-9012-rstu-234567890123",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### DELETE /days/{id} - Delete Day

**Success Response** `204 No Content`:

*(No response body)*

**Error Response** `409 Conflict`:

```json
{
  "error": "cannot delete day: it is used in one or more weeks"
}
```

#### POST /days/{id}/prescriptions - Add Prescription to Day

**Success Response** `201 Created`:

```json
{
  "id": "l2m3n4o5-p6q7-8901-qrst-123456789012",
  "name": "Squat Day",
  "slug": "squat-day",
  "metadata": {
    "focus": "lower body",
    "primaryLift": "squat"
  },
  "programId": "m3n4o5p6-q7r8-9012-rstu-234567890123",
  "prescriptions": [
    {
      "id": "o5p6q7r8-s9t0-1234-uvwx-456789012345",
      "prescriptionId": "h8i9j0k1-l2m3-4567-mnop-789012345678",
      "order": 1,
      "createdAt": "2024-01-01T00:00:00Z"
    },
    {
      "id": "r8s9t0u1-v2w3-4567-xyza-789012345678",
      "prescriptionId": "j0k1l2m3-n4o5-6789-opqr-901234567890",
      "order": 2,
      "createdAt": "2024-01-15T10:30:00Z"
    }
  ],
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### DELETE /days/{id}/prescriptions/{prescriptionId} - Remove Prescription from Day

**Success Response** `204 No Content`:

*(No response body)*

#### PUT /days/{id}/prescriptions/reorder - Reorder Prescriptions

**Success Response** `200 OK`:

```json
{
  "id": "l2m3n4o5-p6q7-8901-qrst-123456789012",
  "name": "Squat Day",
  "slug": "squat-day",
  "metadata": {
    "focus": "lower body",
    "primaryLift": "squat"
  },
  "programId": "m3n4o5p6-q7r8-9012-rstu-234567890123",
  "prescriptions": [
    {
      "id": "r8s9t0u1-v2w3-4567-xyza-789012345678",
      "prescriptionId": "j0k1l2m3-n4o5-6789-opqr-901234567890",
      "order": 1,
      "createdAt": "2024-01-15T10:30:00Z"
    },
    {
      "id": "o5p6q7r8-s9t0-1234-uvwx-456789012345",
      "prescriptionId": "h8i9j0k1-l2m3-4567-mnop-789012345678",
      "order": 2,
      "createdAt": "2024-01-01T00:00:00Z"
    }
  ],
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

---

### Weeks

#### GET /weeks - List Weeks

**Success Response** `200 OK`:

```json
{
  "data": [
    {
      "id": "s9t0u1v2-w3x4-5678-zabc-890123456789",
      "cycleId": "t0u1v2w3-x4y5-6789-abcd-901234567890",
      "weekNumber": 1,
      "name": "Week 1 - 5s",
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    },
    {
      "id": "u1v2w3x4-y5z6-7890-bcde-012345678901",
      "cycleId": "t0u1v2w3-x4y5-6789-abcd-901234567890",
      "weekNumber": 2,
      "name": "Week 2 - 3s",
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "page": 1,
  "pageSize": 20,
  "totalItems": 2,
  "totalPages": 1
}
```

#### GET /weeks/{id} - Get Week by ID (with days)

**Success Response** `200 OK`:

```json
{
  "id": "s9t0u1v2-w3x4-5678-zabc-890123456789",
  "cycleId": "t0u1v2w3-x4y5-6789-abcd-901234567890",
  "weekNumber": 1,
  "name": "Week 1 - 5s",
  "days": [
    {
      "id": "v2w3x4y5-z6a7-8901-cdef-123456789012",
      "dayId": "l2m3n4o5-p6q7-8901-qrst-123456789012",
      "position": 0
    },
    {
      "id": "w3x4y5z6-a7b8-9012-defg-234567890123",
      "dayId": "n4o5p6q7-r8s9-0123-stuv-345678901234",
      "position": 1
    },
    {
      "id": "x4y5z6a7-b8c9-0123-efgh-345678901234",
      "dayId": "q7r8s9t0-u1v2-3456-wxyz-678901234567",
      "position": 2
    }
  ],
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### POST /weeks - Create Week

**Success Response** `201 Created`:

```json
{
  "id": "y5z6a7b8-c9d0-1234-fghi-456789012345",
  "cycleId": "t0u1v2w3-x4y5-6789-abcd-901234567890",
  "weekNumber": 3,
  "name": "Week 3 - 5/3/1",
  "days": [],
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### PUT /weeks/{id} - Update Week

**Success Response** `200 OK`:

```json
{
  "id": "s9t0u1v2-w3x4-5678-zabc-890123456789",
  "cycleId": "t0u1v2w3-x4y5-6789-abcd-901234567890",
  "weekNumber": 1,
  "name": "Week 1 - Volume",
  "days": [
    {
      "id": "v2w3x4y5-z6a7-8901-cdef-123456789012",
      "dayId": "l2m3n4o5-p6q7-8901-qrst-123456789012",
      "position": 0
    }
  ],
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### DELETE /weeks/{id} - Delete Week

**Success Response** `204 No Content`:

*(No response body)*

#### POST /weeks/{id}/days - Add Day to Week

**Success Response** `201 Created`:

```json
{
  "id": "s9t0u1v2-w3x4-5678-zabc-890123456789",
  "cycleId": "t0u1v2w3-x4y5-6789-abcd-901234567890",
  "weekNumber": 1,
  "name": "Week 1 - 5s",
  "days": [
    {
      "id": "v2w3x4y5-z6a7-8901-cdef-123456789012",
      "dayId": "l2m3n4o5-p6q7-8901-qrst-123456789012",
      "position": 0
    },
    {
      "id": "z6a7b8c9-d0e1-2345-ghij-567890123456",
      "dayId": "n4o5p6q7-r8s9-0123-stuv-345678901234",
      "position": 1
    }
  ],
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### DELETE /weeks/{id}/days/{dayId} - Remove Day from Week

**Success Response** `204 No Content`:

*(No response body)*

---

### Cycles

#### GET /cycles - List Cycles

**Success Response** `200 OK`:

```json
{
  "data": [
    {
      "id": "t0u1v2w3-x4y5-6789-abcd-901234567890",
      "name": "4 Week 5/3/1 Cycle",
      "lengthWeeks": 4,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    },
    {
      "id": "a7b8c9d0-e1f2-3456-hijk-678901234567",
      "name": "Weekly Linear Cycle",
      "lengthWeeks": 1,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "page": 1,
  "pageSize": 20,
  "totalItems": 2,
  "totalPages": 1
}
```

#### GET /cycles/{id} - Get Cycle by ID

**Success Response** `200 OK`:

```json
{
  "id": "t0u1v2w3-x4y5-6789-abcd-901234567890",
  "name": "4 Week 5/3/1 Cycle",
  "lengthWeeks": 4,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### POST /cycles - Create Cycle

**Success Response** `201 Created`:

```json
{
  "id": "b8c9d0e1-f2a3-4567-ijkl-789012345678",
  "name": "3 Week Undulating Cycle",
  "lengthWeeks": 3,
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### PUT /cycles/{id} - Update Cycle

**Success Response** `200 OK`:

```json
{
  "id": "t0u1v2w3-x4y5-6789-abcd-901234567890",
  "name": "4 Week BBB Cycle",
  "lengthWeeks": 4,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### DELETE /cycles/{id} - Delete Cycle

**Success Response** `204 No Content`:

*(No response body)*

---

### Weekly Lookups

#### GET /weekly-lookups - List Weekly Lookups

**Success Response** `200 OK`:

```json
{
  "data": [
    {
      "id": "c9d0e1f2-a3b4-5678-jklm-890123456789",
      "name": "5/3/1 Top Set Percentages",
      "entries": {
        "1": 85,
        "2": 90,
        "3": 95,
        "4": 60
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

#### GET /weekly-lookups/{id} - Get Weekly Lookup by ID

**Success Response** `200 OK`:

```json
{
  "id": "c9d0e1f2-a3b4-5678-jklm-890123456789",
  "name": "5/3/1 Top Set Percentages",
  "entries": {
    "1": 85,
    "2": 90,
    "3": 95,
    "4": 60
  },
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### POST /weekly-lookups - Create Weekly Lookup

**Success Response** `201 Created`:

```json
{
  "id": "d0e1f2a3-b4c5-6789-klmn-901234567890",
  "name": "Nuckols 3-Week Cycle",
  "entries": {
    "1": 80,
    "2": 85,
    "3": 90
  },
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### PUT /weekly-lookups/{id} - Update Weekly Lookup

**Success Response** `200 OK`:

```json
{
  "id": "c9d0e1f2-a3b4-5678-jklm-890123456789",
  "name": "5/3/1 Top Set Percentages (Modified)",
  "entries": {
    "1": 80,
    "2": 85,
    "3": 90,
    "4": 50
  },
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### DELETE /weekly-lookups/{id} - Delete Weekly Lookup

**Success Response** `204 No Content`:

*(No response body)*

---

### Daily Lookups

#### GET /daily-lookups - List Daily Lookups

**Success Response** `200 OK`:

```json
{
  "data": [
    {
      "id": "e1f2a3b4-c5d6-7890-lmno-012345678901",
      "name": "Heavy/Light/Medium",
      "entries": {
        "monday": 100,
        "wednesday": 80,
        "friday": 90
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

#### GET /daily-lookups/{id} - Get Daily Lookup by ID

**Success Response** `200 OK`:

```json
{
  "id": "e1f2a3b4-c5d6-7890-lmno-012345678901",
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

#### POST /daily-lookups - Create Daily Lookup

**Success Response** `201 Created`:

```json
{
  "id": "f2a3b4c5-d6e7-8901-mnop-123456789012",
  "name": "Daily Undulating Periodization",
  "entries": {
    "day1": 85,
    "day2": 70,
    "day3": 90,
    "day4": 75,
    "day5": 95
  },
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### PUT /daily-lookups/{id} - Update Daily Lookup

**Success Response** `200 OK`:

```json
{
  "id": "e1f2a3b4-c5d6-7890-lmno-012345678901",
  "name": "Modified H/L/M",
  "entries": {
    "monday": 100,
    "wednesday": 70,
    "friday": 85
  },
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### DELETE /daily-lookups/{id} - Delete Daily Lookup

**Success Response** `204 No Content`:

*(No response body)*

---

### Programs

#### GET /programs - List Programs

**Success Response** `200 OK`:

```json
{
  "data": [
    {
      "id": "m3n4o5p6-q7r8-9012-rstu-234567890123",
      "name": "Wendler 5/3/1 BBB",
      "slug": "wendler-531-bbb",
      "description": "The classic Boring But Big variant of 5/3/1 featuring high-volume accessory work at 50% of training max.",
      "cycleId": "t0u1v2w3-x4y5-6789-abcd-901234567890",
      "weeklyLookupId": "c9d0e1f2-a3b4-5678-jklm-890123456789",
      "dailyLookupId": null,
      "defaultRounding": 5.0,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    },
    {
      "id": "g3a4b5c6-d7e8-9012-nopq-345678901234",
      "name": "Starting Strength",
      "slug": "starting-strength",
      "description": "Classic linear progression program for novice lifters.",
      "cycleId": "a7b8c9d0-e1f2-3456-hijk-678901234567",
      "weeklyLookupId": null,
      "dailyLookupId": null,
      "defaultRounding": 5.0,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "page": 1,
  "pageSize": 20,
  "totalItems": 2,
  "totalPages": 1
}
```

#### GET /programs/{id} - Get Program by ID (with embedded details)

**Success Response** `200 OK`:

```json
{
  "id": "m3n4o5p6-q7r8-9012-rstu-234567890123",
  "name": "Wendler 5/3/1 BBB",
  "slug": "wendler-531-bbb",
  "description": "The classic Boring But Big variant of 5/3/1 featuring high-volume accessory work at 50% of training max.",
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
  "defaultRounding": 5.0,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### POST /programs - Create Program

**Success Response** `201 Created`:

```json
{
  "id": "h5i6j7k8-l9m0-1234-qrst-567890123456",
  "name": "Bill Starr 5x5",
  "slug": "bill-starr-5x5",
  "description": "Heavy/Light/Medium 5x5 program with ramping sets and daily intensity variation.",
  "cycleId": "a7b8c9d0-e1f2-3456-hijk-678901234567",
  "weeklyLookupId": null,
  "dailyLookupId": "e1f2a3b4-c5d6-7890-lmno-012345678901",
  "defaultRounding": 5.0,
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

**Error Response** `409 Conflict`:

```json
{
  "error": "slug already exists"
}
```

#### PUT /programs/{id} - Update Program

**Success Response** `200 OK`:

```json
{
  "id": "m3n4o5p6-q7r8-9012-rstu-234567890123",
  "name": "Wendler 5/3/1 BBB",
  "slug": "wendler-531-bbb",
  "description": "Updated description with more detail about the program structure.",
  "cycleId": "t0u1v2w3-x4y5-6789-abcd-901234567890",
  "weeklyLookupId": "c9d0e1f2-a3b4-5678-jklm-890123456789",
  "dailyLookupId": null,
  "defaultRounding": 2.5,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### DELETE /programs/{id} - Delete Program

**Success Response** `204 No Content`:

*(No response body)*

**Error Response** `409 Conflict`:

```json
{
  "error": "cannot delete program: users are enrolled in this program"
}
```

---

### Progressions

#### GET /progressions - List Progressions

**Success Response** `200 OK`:

```json
{
  "data": [
    {
      "id": "i6j7k8l9-m0n1-2345-rstu-678901234567",
      "name": "Linear +5lb per session",
      "type": "LINEAR",
      "parameters": {
        "increment": 5.0,
        "maxType": "TRAINING_MAX",
        "triggerType": "AFTER_SESSION"
      },
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    },
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
    },
    {
      "id": "k8l9m0n1-o2p3-4567-tuvw-890123456789",
      "name": "Cycle +2.5lb (upper body)",
      "type": "CYCLE",
      "parameters": {
        "increment": 2.5,
        "maxType": "TRAINING_MAX"
      },
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

#### GET /progressions/{id} - Get Progression by ID

**Success Response** `200 OK` (LINEAR type):

```json
{
  "id": "i6j7k8l9-m0n1-2345-rstu-678901234567",
  "name": "Linear +5lb per session",
  "type": "LINEAR",
  "parameters": {
    "increment": 5.0,
    "maxType": "TRAINING_MAX",
    "triggerType": "AFTER_SESSION"
  },
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

**Success Response** `200 OK` (CYCLE type):

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

#### POST /progressions - Create Progression

**Success Response** `201 Created`:

```json
{
  "id": "l9m0n1o2-p3q4-5678-uvwx-901234567890",
  "name": "Weekly +5lb",
  "type": "LINEAR",
  "parameters": {
    "increment": 5.0,
    "maxType": "TRAINING_MAX",
    "triggerType": "AFTER_WEEK"
  },
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### PUT /progressions/{id} - Update Progression

**Success Response** `200 OK`:

```json
{
  "id": "i6j7k8l9-m0n1-2345-rstu-678901234567",
  "name": "Modified Linear +10lb",
  "type": "LINEAR",
  "parameters": {
    "increment": 10.0,
    "maxType": "TRAINING_MAX",
    "triggerType": "AFTER_SESSION"
  },
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### DELETE /progressions/{id} - Delete Progression

**Success Response** `204 No Content`:

*(No response body)*

**Error Response** `409 Conflict`:

```json
{
  "error": "cannot delete progression: it is referenced by program progressions"
}
```

---

### Program Progressions

#### GET /programs/{programId}/progressions - List Program Progressions

**Success Response** `200 OK`:

```json
[
  {
    "id": "m0n1o2p3-q4r5-6789-vwxy-012345678901",
    "programId": "m3n4o5p6-q7r8-9012-rstu-234567890123",
    "progressionId": "j7k8l9m0-n1o2-3456-stuv-789012345678",
    "liftId": null,
    "priority": 1,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  },
  {
    "id": "n1o2p3q4-r5s6-7890-wxyz-123456789012",
    "programId": "m3n4o5p6-q7r8-9012-rstu-234567890123",
    "progressionId": "k8l9m0n1-o2p3-4567-tuvw-890123456789",
    "liftId": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
    "priority": 2,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
]
```

#### GET /programs/{programId}/progressions/{configId} - Get Program Progression

**Success Response** `200 OK`:

```json
{
  "id": "m0n1o2p3-q4r5-6789-vwxy-012345678901",
  "programId": "m3n4o5p6-q7r8-9012-rstu-234567890123",
  "progressionId": "j7k8l9m0-n1o2-3456-stuv-789012345678",
  "liftId": null,
  "priority": 1,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### POST /programs/{programId}/progressions - Create Program Progression

**Success Response** `201 Created`:

```json
{
  "id": "o2p3q4r5-s6t7-8901-xyza-234567890123",
  "programId": "m3n4o5p6-q7r8-9012-rstu-234567890123",
  "progressionId": "i6j7k8l9-m0n1-2345-rstu-678901234567",
  "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "priority": 2,
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### PUT /programs/{programId}/progressions/{configId} - Update Program Progression

**Success Response** `200 OK`:

```json
{
  "id": "m0n1o2p3-q4r5-6789-vwxy-012345678901",
  "programId": "m3n4o5p6-q7r8-9012-rstu-234567890123",
  "progressionId": "j7k8l9m0-n1o2-3456-stuv-789012345678",
  "liftId": null,
  "priority": 3,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### DELETE /programs/{programId}/progressions/{configId} - Delete Program Progression

**Success Response** `204 No Content`:

*(No response body)*

---

### User Program Enrollment

#### GET /users/{userId}/program - Get Current Enrollment

**Success Response** `200 OK`:

```json
{
  "id": "p3q4r5s6-t7u8-9012-yzab-345678901234",
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "program": {
    "id": "m3n4o5p6-q7r8-9012-rstu-234567890123",
    "name": "Wendler 5/3/1 BBB",
    "slug": "wendler-531-bbb",
    "description": "The classic Boring But Big variant of 5/3/1.",
    "cycleLengthWeeks": 4
  },
  "state": {
    "currentWeek": 2,
    "currentCycleIteration": 1,
    "currentDayIndex": 1
  },
  "enrolledAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-10T08:00:00Z"
}
```

**Error Response** `404 Not Found`:

```json
{
  "error": "enrollment not found"
}
```

#### POST /users/{userId}/program - Enroll User in Program

**Success Response** `201 Created`:

```json
{
  "id": "q4r5s6t7-u8v9-0123-zabc-456789012345",
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "program": {
    "id": "m3n4o5p6-q7r8-9012-rstu-234567890123",
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

**Error Response** `400 Bad Request`:

```json
{
  "error": "programId: program not found"
}
```

#### DELETE /users/{userId}/program - Unenroll User

**Success Response** `204 No Content`:

*(No response body)*

---

### State Advancement

#### POST /users/{userId}/program-state/advance - Advance State

**Success Response** `200 OK` (advance day):

```json
{
  "id": "p3q4r5s6-t7u8-9012-yzab-345678901234",
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "program": {
    "id": "m3n4o5p6-q7r8-9012-rstu-234567890123",
    "name": "Wendler 5/3/1 BBB",
    "slug": "wendler-531-bbb",
    "description": "The classic Boring But Big variant of 5/3/1.",
    "cycleLengthWeeks": 4
  },
  "state": {
    "currentWeek": 1,
    "currentCycleIteration": 1,
    "currentDayIndex": 1
  },
  "enrolledAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

**Success Response** `200 OK` (advance week - new cycle):

```json
{
  "id": "p3q4r5s6-t7u8-9012-yzab-345678901234",
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "program": {
    "id": "m3n4o5p6-q7r8-9012-rstu-234567890123",
    "name": "Wendler 5/3/1 BBB",
    "slug": "wendler-531-bbb",
    "description": "The classic Boring But Big variant of 5/3/1.",
    "cycleLengthWeeks": 4
  },
  "state": {
    "currentWeek": 1,
    "currentCycleIteration": 2,
    "currentDayIndex": 0
  },
  "enrolledAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-28T10:30:00Z"
}
```

---

### Workout Generation

#### GET /users/{userId}/workout - Generate Current Workout

**Success Response** `200 OK`:

```json
{
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "programId": "m3n4o5p6-q7r8-9012-rstu-234567890123",
  "cycleIteration": 1,
  "weekNumber": 1,
  "daySlug": "squat-day",
  "date": "2024-01-15",
  "exercises": [
    {
      "prescriptionId": "k1l2m3n4-o5p6-7890-pqrs-012345678901",
      "lift": {
        "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "name": "Squat",
        "slug": "squat"
      },
      "sets": [
        {"setNumber": 1, "weight": 125.0, "targetReps": 5, "isWorkSet": false},
        {"setNumber": 2, "weight": 155.0, "targetReps": 5, "isWorkSet": false},
        {"setNumber": 3, "weight": 185.0, "targetReps": 3, "isWorkSet": false},
        {"setNumber": 4, "weight": 220.0, "targetReps": 2, "isWorkSet": false}
      ],
      "notes": "Warmup - build up to working weight",
      "restSeconds": null
    },
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
      "notes": "Main work sets",
      "restSeconds": 180
    },
    {
      "prescriptionId": "i9j0k1l2-m3n4-5678-nopq-890123456789",
      "lift": {
        "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "name": "Squat",
        "slug": "squat"
      },
      "sets": [
        {"setNumber": 1, "weight": 265.0, "targetReps": 5, "isWorkSet": true}
      ],
      "notes": "Top set - as many reps as possible",
      "restSeconds": 300
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
      "notes": "BBB volume work - focus on quality reps",
      "restSeconds": 90
    }
  ]
}
```

**Error Response** `404 Not Found` (not enrolled):

```json
{
  "error": "enrollment not found"
}
```

**Error Response** `400 Bad Request` (missing lift max):

```json
{
  "error": "missing lift max: set up your training maxes to generate workouts",
  "details": [
    "missing lift max for squat (TRAINING_MAX)"
  ]
}
```

#### GET /users/{userId}/workout/preview - Preview Workout

**Success Response** `200 OK`:

```json
{
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "programId": "m3n4o5p6-q7r8-9012-rstu-234567890123",
  "cycleIteration": 1,
  "weekNumber": 3,
  "daySlug": "squat-day",
  "date": "2024-01-15",
  "exercises": [
    {
      "prescriptionId": "i9j0k1l2-m3n4-5678-nopq-890123456789",
      "lift": {
        "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "name": "Squat",
        "slug": "squat"
      },
      "sets": [
        {"setNumber": 1, "weight": 300.0, "targetReps": 1, "isWorkSet": true}
      ],
      "notes": "Top set - as many reps as possible",
      "restSeconds": 300
    }
  ]
}
```

**Error Response** `400 Bad Request`:

```json
{
  "error": "week: missing required parameter"
}
```

---

### Progression History

#### GET /users/{userId}/progression-history - List Progression History

**Success Response** `200 OK`:

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
    },
    {
      "id": "s6t7u8v9-w0x1-2345-cdef-678901234567",
      "userId": "550e8400-e29b-41d4-a716-446655440000",
      "liftId": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
      "progressionId": "k8l9m0n1-o2p3-4567-tuvw-890123456789",
      "previousValue": 225.0,
      "newValue": 227.5,
      "delta": 2.5,
      "appliedAt": "2024-01-28T10:30:00Z",
      "triggeredBy": "CYCLE_END",
      "cycleIteration": 1,
      "weekNumber": 4
    }
  ],
  "page": 1,
  "pageSize": 20,
  "totalItems": 2,
  "totalPages": 1
}
```

---

### Manual Progression Trigger

#### POST /users/{userId}/progressions/trigger - Trigger Progression

**Success Response** `200 OK`:

```json
{
  "results": [
    {
      "progressionId": "j7k8l9m0-n1o2-3456-stuv-789012345678",
      "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "applied": true,
      "skipped": false,
      "result": {
        "previousValue": 320.0,
        "newValue": 325.0,
        "delta": 5.0,
        "maxType": "TRAINING_MAX",
        "appliedAt": "2024-01-15T10:30:00Z"
      }
    },
    {
      "progressionId": "j7k8l9m0-n1o2-3456-stuv-789012345678",
      "liftId": "d4e5f6a7-b8c9-0123-defg-345678901234",
      "applied": true,
      "skipped": false,
      "result": {
        "previousValue": 405.0,
        "newValue": 410.0,
        "delta": 5.0,
        "maxType": "TRAINING_MAX",
        "appliedAt": "2024-01-15T10:30:00Z"
      }
    }
  ],
  "totalApplied": 2,
  "totalSkipped": 0,
  "totalErrors": 0
}
```

**Success Response** `200 OK` (with skipped):

```json
{
  "results": [
    {
      "progressionId": "j7k8l9m0-n1o2-3456-stuv-789012345678",
      "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "applied": false,
      "skipped": true,
      "skipReason": "progression already applied this period"
    }
  ],
  "totalApplied": 0,
  "totalSkipped": 1,
  "totalErrors": 0
}
```

**Error Response** `404 Not Found`:

```json
{
  "error": "progression not found"
}
```

**Error Response** `400 Bad Request`:

```json
{
  "error": "user not enrolled in a program"
}
```

---

## Response Headers

All responses include the following headers:

```http
Content-Type: application/json
```

For successful responses:
```http
HTTP/1.1 200 OK
Content-Type: application/json
```

For created resources:
```http
HTTP/1.1 201 Created
Content-Type: application/json
```

For deletions:
```http
HTTP/1.1 204 No Content
```

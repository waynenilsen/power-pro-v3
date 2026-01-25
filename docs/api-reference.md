# PowerPro API Reference

PowerPro is a headless API for powerlifting program management. This document provides comprehensive documentation for all API endpoints.

## Machine-Readable Specification

An OpenAPI 3.0 specification is available at [`docs/openapi.yaml`](./openapi.yaml). This specification can be used with tools like Swagger UI, Postman, or code generators to explore and interact with the API.

## Example Requests

Copy-paste ready example requests for all endpoints are available at [`docs/example-requests.md`](./example-requests.md). These examples use realistic data and demonstrate common workflows.

## Example Responses

Complete example responses for all endpoints are available at [`docs/example-responses.md`](./example-responses.md). This includes success responses, common error responses, and pagination format details.

## Error Documentation

Comprehensive error handling documentation is available at [`docs/errors.md`](./errors.md). This includes all HTTP status codes, error response formats, common error scenarios, and error handling best practices.

## Workflow Documentation

Common multi-step API workflows are documented at [`docs/workflows.md`](./workflows.md). This includes:
- User onboarding workflow
- Program enrollment workflow
- Workout generation workflow
- Progression trigger workflow
- State machine workout flow

## Base URL

```
http://localhost:{port}
```

## Authentication

PowerPro uses session-based authentication with Bearer tokens. Users register and login to obtain a session token, which is then used to authenticate subsequent requests.

### Auth Endpoints

See the [Authentication Endpoints](#authentication-1) section below for register, login, and logout operations.

### Headers

| Header | Description |
|--------|-------------|
| `Authorization` | Bearer token format: `Bearer {session-token}` (obtained from login response) |
| `X-User-ID` | Alternative: User ID directly (for development/testing only) |
| `X-Admin` | Set to `"true"` for admin privileges |

### Authentication Levels

- **Public**: No authentication required (health check, register, login)
- **Authenticated**: Any authenticated user with valid session token
- **Admin**: Requires `X-Admin: true`
- **Owner/Admin**: User must own the resource or be admin
- **Owner-only**: Only the resource owner can access (not even admins)

---

## Standard Response Envelope

All API responses follow a consistent envelope format for predictable client handling.

### Success Response (Single Entity)

```json
{
  "data": {
    "id": "abc123",
    "name": "Squat",
    ...
  }
}
```

### Success Response with Warnings

```json
{
  "data": {
    "id": "abc123",
    ...
  },
  "warnings": ["informational message"]
}
```

### Paginated Response

```json
{
  "data": [
    {"id": "abc123", ...},
    {"id": "def456", ...}
  ],
  "meta": {
    "total": 100,
    "limit": 20,
    "offset": 0,
    "hasMore": true
  }
}
```

### Error Response

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "lift not found: abc123",
    "details": {
      "validationErrors": ["field is required"]
    }
  }
}
```

#### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `NOT_FOUND` | 404 | Resource not found |
| `VALIDATION_ERROR` | 400 | Request validation failed |
| `BAD_REQUEST` | 400 | Malformed request |
| `CONFLICT` | 409 | Resource conflict (e.g., duplicate) |
| `FORBIDDEN` | 403 | Permission denied |
| `UNAUTHORIZED` | 401 | Authentication required |
| `INTERNAL_ERROR` | 500 | Server error |

---

## Pagination

All list endpoints use consistent offset-based pagination with the following query parameters:

| Parameter | Type | Default | Max | Description |
|-----------|------|---------|-----|-------------|
| `limit` | int | 20 | 100 | Number of items to return |
| `offset` | int | 0 | - | Number of items to skip |
| `sortBy` | string | varies | - | Field to sort by |
| `sortOrder` | string | "asc" | - | "asc" or "desc" |

### Response Metadata

All paginated responses include a `meta` object with pagination information:

```json
{
  "data": [...],
  "meta": {
    "total": 150,
    "limit": 20,
    "offset": 40,
    "hasMore": true
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| `total` | int64 | Total number of items in the dataset |
| `limit` | int | Number of items per page (as requested) |
| `offset` | int | Current offset (as requested) |
| `hasMore` | bool | Whether more items exist beyond the current page |

### Pagination Example

```bash
# Get the first 10 lifts
GET /lifts?limit=10&offset=0

# Get the next 10 lifts
GET /lifts?limit=10&offset=10

# Get the third page (items 21-30)
GET /lifts?limit=10&offset=20
```

---

## Filtering

All list endpoints support consistent filtering through query parameters. This section documents the filtering conventions used across the API.

### Filter Parameter Naming

Filter parameters use **snake_case** naming convention to match query parameter conventions:

| Pattern | Description | Example |
|---------|-------------|---------|
| Simple filters | Use field names directly | `?lift_id=123&user_id=456` |
| Boolean filters | Use "true"/"false" or "1"/"0" | `?is_competition_lift=true` |
| Date ranges | Use `_after`/`_before` suffixes | `?start_date=2024-01-01&end_date=2024-12-31` |
| Numeric ranges | Use `_gte`/`_lte` suffixes | `?weight_gte=100&weight_lte=200` |
| Enum filters | Provide the enum value (case-insensitive) | `?type=TRAINING_MAX` |

### Filter Behavior

- **AND logic**: Multiple filters are combined with AND logic. For example, `?lift_id=123&type=TRAINING_MAX` returns only items matching both conditions.
- **Unknown parameters**: Unknown filter parameters are silently ignored.
- **Invalid values**: Invalid filter values return a `400 Bad Request` with a validation error.
- **Empty values**: Empty filter values (e.g., `?lift_id=`) are treated as if the parameter was not provided.

### Date Format

Date filters accept ISO 8601 formats:
- Full datetime: `2024-01-15T10:30:00Z` (RFC3339)
- Date only: `2024-01-15` (interpreted as start of day for "after" filters, end of day for "before" filters)

### Filter Examples

```bash
# Filter lifts by competition status
GET /lifts?is_competition_lift=true

# Filter lift maxes by lift and type
GET /users/{userId}/lift-maxes?lift_id=abc123&type=TRAINING_MAX

# Filter prescriptions by lift
GET /prescriptions?lift_id=abc123

# Filter days by program
GET /days?program_id=abc123

# Filter progression history by date range
GET /users/{userId}/progression-history?start_date=2024-01-01&end_date=2024-03-31

# Filter progression history by lift and type
GET /users/{userId}/progression-history?lift_id=abc123&progression_type=LINEAR_PROGRESSION
```

### Available Filters by Endpoint

| Endpoint | Available Filters |
|----------|-------------------|
| `GET /lifts` | `is_competition_lift` (bool) |
| `GET /users/{userId}/lift-maxes` | `lift_id` (string), `type` (enum: ONE_RM, TRAINING_MAX) |
| `GET /prescriptions` | `lift_id` (string) |
| `GET /days` | `program_id` (string) |
| `GET /users/{userId}/progression-history` | `lift_id` (string), `progression_type` (enum), `trigger_type` (enum), `start_date` (date), `end_date` (date) |
| `GET /progressions` | `type` (enum: LINEAR, CYCLE) |

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

### Authentication

User registration, login, and session management.

#### POST /auth/register

Register a new user account.

**Auth**: Public

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "securepassword123",
  "name": "John Doe"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `email` | string | Yes | User's email address (must be unique) |
| `password` | string | Yes | Password (min 8 characters recommended) |
| `name` | string | No | User's display name |

**Response** `201 Created`:
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "John Doe",
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T10:30:00Z"
  }
}
```

**Errors**:
- `400 Bad Request`: Invalid JSON, missing email or password
- `409 Conflict`: Email already exists

#### POST /auth/login

Authenticate and obtain a session token.

**Auth**: Public

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `email` | string | Yes | User's email address |
| `password` | string | Yes | User's password |

**Response** `200 OK`:
```json
{
  "data": {
    "token": "session-token-string",
    "expiresAt": "2024-01-22T10:30:00Z",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "name": "John Doe",
      "createdAt": "2024-01-15T10:30:00Z",
      "updatedAt": "2024-01-15T10:30:00Z"
    }
  }
}
```

**Notes**:
- Session tokens are valid for 7 days from creation
- Use the returned `token` in the `Authorization: Bearer {token}` header for authenticated requests

**Errors**:
- `400 Bad Request`: Invalid JSON, missing email or password
- `401 Unauthorized`: Invalid email or password

#### POST /auth/logout

End the current session and invalidate the token.

**Auth**: Authenticated

**Request Body**: None

**Response** `204 No Content`

**Errors**:
- `401 Unauthorized`: Missing or invalid authentication token

---

### User Profile

Manage user profile information.

#### GET /users/{userId}/profile

Get a user's profile information.

**Auth**: Owner/Admin

**Path Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `userId` | string | User ID (UUID) |

**Response** `200 OK`:
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "John Doe",
    "weightUnit": "lb",
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T10:30:00Z"
  }
}
```

**Errors**:
- `400 Bad Request`: Missing user ID in path
- `403 Forbidden`: Accessing another user's profile (without admin privileges)
- `404 Not Found`: User not found

#### PUT /users/{userId}/profile

Update a user's profile information.

**Auth**: Owner-only (admin cannot override)

**Path Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `userId` | string | User ID (UUID) |

**Request Body**:
```json
{
  "name": "Jane Doe",
  "weightUnit": "kg"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | No | User's display name |
| `weightUnit` | string | No | Preferred weight unit ("lb" or "kg") |

**Response** `200 OK`:
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "Jane Doe",
    "weightUnit": "kg",
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T12:00:00Z"
  }
}
```

**Notes**:
- Profile updates are strictly owner-only; even admins cannot modify another user's profile

**Errors**:
- `400 Bad Request`: Invalid JSON, missing user ID
- `403 Forbidden`: Not the profile owner (even admins are blocked)
- `404 Not Found`: User not found

---

### Dashboard

User dashboard with aggregated data.

#### GET /users/{id}/dashboard

Get aggregated dashboard data for a user including enrollment status, current/next workout info, and recent activity.

**Auth**: Owner-only (admin cannot access)

**Path Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string | User ID (UUID) |

**Response** `200 OK`:
```json
{
  "data": {
    "enrollment": {
      "status": "ACTIVE",
      "programName": "Wendler 5/3/1 BBB",
      "cycleIteration": 1,
      "cycleStatus": "IN_PROGRESS",
      "weekNumber": 1,
      "weekStatus": "IN_PROGRESS"
    },
    "nextWorkout": {
      "dayName": "Squat Day",
      "daySlug": "squat-day",
      "exerciseCount": 4,
      "estimatedSets": 12
    },
    "currentSession": {
      "sessionId": "session-uuid",
      "dayName": "Squat Day",
      "startedAt": "2024-01-15T08:00:00Z",
      "setsCompleted": 3,
      "totalSets": 5
    },
    "recentWorkouts": [
      {
        "date": "2024-01-15",
        "dayName": "Squat Day",
        "setsCompleted": 5
      }
    ],
    "currentMaxes": [
      {
        "lift": "Squat",
        "value": 315.0,
        "type": "TRAINING_MAX"
      }
    ]
  }
}
```

**Response Fields**:

| Field | Type | Description |
|-------|------|-------------|
| `enrollment` | object | Current program enrollment status (null if not enrolled) |
| `enrollment.status` | string | Enrollment status: "ACTIVE", "BETWEEN_CYCLES", "QUIT" |
| `enrollment.programName` | string | Name of enrolled program |
| `enrollment.cycleIteration` | int | Current cycle number (1, 2, 3...) |
| `enrollment.cycleStatus` | string | Cycle status: "PENDING", "IN_PROGRESS", "COMPLETED" |
| `enrollment.weekNumber` | int | Current week number within cycle |
| `enrollment.weekStatus` | string | Week status: "PENDING", "IN_PROGRESS", "COMPLETED" |
| `nextWorkout` | object | Info about the next workout (null if none scheduled) |
| `nextWorkout.dayName` | string | Name of the training day |
| `nextWorkout.daySlug` | string | Slug identifier for the day |
| `nextWorkout.exerciseCount` | int | Number of exercises in the workout |
| `nextWorkout.estimatedSets` | int | Total estimated sets |
| `currentSession` | object | Current in-progress workout session (null if none) |
| `recentWorkouts` | array | Recent completed workouts (up to 5) |
| `currentMaxes` | array | User's current training maxes for competition lifts |

**Notes**:
- This endpoint is owner-only; even admins cannot access another user's dashboard
- Fields may be null if the user hasn't enrolled in a program or has no workout history

**Errors**:
- `400 Bad Request`: Missing user ID
- `403 Forbidden`: Not the dashboard owner (no admin override)
- `404 Not Found`: User or enrollment not found

---

### Lifts

Manage exercises (lifts) in the system.

#### GET /lifts

List all lifts with pagination.

**Auth**: Authenticated

**Query Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `limit` | int | Number of items to return (default: 20, max: 100) |
| `offset` | int | Number of items to skip (default: 0) |
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
  "meta": {
    "total": 3,
    "limit": 20,
    "offset": 0,
    "hasMore": false
  }
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
| `limit` | int | Number of items to return (default: 20, max: 100) |
| `offset` | int | Number of items to skip (default: 0) |
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
  "meta": {
    "total": 1,
    "limit": 20,
    "offset": 0,
    "hasMore": false
  }
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
| `limit` | int | Number of items to return (default: 20, max: 100) |
| `offset` | int | Number of items to skip (default: 0) |
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
  "meta": {
    "total": 1,
    "limit": 20,
    "offset": 0,
    "hasMore": false
  }
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
| `limit` | int | Number of items to return (default: 20, max: 100) |
| `offset` | int | Number of items to skip (default: 0) |
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
  "meta": {
    "total": 1,
    "limit": 20,
    "offset": 0,
    "hasMore": false
  }
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
| `limit` | int | Number of items to return (default: 20, max: 100) |
| `offset` | int | Number of items to skip (default: 0) |
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
| `limit` | int | Number of items to return (default: 20, max: 100) |
| `offset` | int | Number of items to skip (default: 0) |
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
  "meta": {
    "total": 1,
    "limit": 20,
    "offset": 0,
    "hasMore": false
  }
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
| `limit` | int | Number of items to return (default: 20, max: 100) |
| `offset` | int | Number of items to skip (default: 0) |
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
  "meta": {
    "total": 1,
    "limit": 20,
    "offset": 0,
    "hasMore": false
  }
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
| `limit` | int | Number of items to return (default: 20, max: 100) |
| `offset` | int | Number of items to skip (default: 0) |
| `lift_id` | string | Filter by lift ID |
| `progression_type` | string | Filter by progression type: "LINEAR_PROGRESSION" or "CYCLE_PROGRESSION" |
| `trigger_type` | string | Filter by trigger type: "AFTER_SESSION", "AFTER_WEEK", or "AFTER_CYCLE" |
| `start_date` | date | Filter entries on or after this date (ISO 8601) |
| `end_date` | date | Filter entries on or before this date (ISO 8601) |

**Response** `200 OK`:
```json
{
  "data": [
    {
      "id": "history-uuid",
      "progressionId": "progression-uuid",
      "progressionName": "Linear +5lb",
      "progressionType": "LINEAR_PROGRESSION",
      "liftId": "lift-uuid",
      "liftName": "Squat",
      "previousValue": 315.0,
      "newValue": 320.0,
      "delta": 5.0,
      "triggerType": "AFTER_SESSION",
      "triggerContext": {},
      "appliedAt": "2024-01-15T10:30:00Z"
    }
  ],
  "meta": {
    "total": 1,
    "limit": 20,
    "offset": 0,
    "hasMore": false
  }
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

### Workout Sessions

Manage workout sessions for tracking workout lifecycle.

#### POST /workouts/start

Start a new workout session for the authenticated user.

**Auth**: Authenticated (uses current user's enrollment)

**Request Body**: None required (uses current enrollment state)

**Response** `201 Created`:
```json
{
  "data": {
    "id": "session-uuid",
    "userProgramStateId": "enrollment-uuid",
    "weekNumber": 1,
    "dayIndex": 0,
    "status": "IN_PROGRESS",
    "startedAt": "2024-01-15T08:00:00Z",
    "createdAt": "2024-01-15T08:00:00Z",
    "updatedAt": "2024-01-15T08:00:00Z"
  }
}
```

**Errors**:
- `404 Not Found`: User not enrolled in a program
- `400 Bad Request`: Enrollment not in ACTIVE state
- `409 Conflict`: User already has an in-progress workout session

#### GET /workouts/{id}

Get a workout session by ID.

**Auth**: Owner/Admin

**Response** `200 OK`:
```json
{
  "data": {
    "id": "session-uuid",
    "userProgramStateId": "enrollment-uuid",
    "weekNumber": 1,
    "dayIndex": 0,
    "status": "IN_PROGRESS",
    "startedAt": "2024-01-15T08:00:00Z",
    "finishedAt": null,
    "createdAt": "2024-01-15T08:00:00Z",
    "updatedAt": "2024-01-15T08:00:00Z"
  }
}
```

#### POST /workouts/{id}/finish

Complete a workout session.

**Auth**: Owner/Admin

**Request Body**: None required

**Response** `200 OK`:
```json
{
  "data": {
    "id": "session-uuid",
    "userProgramStateId": "enrollment-uuid",
    "weekNumber": 1,
    "dayIndex": 0,
    "status": "COMPLETED",
    "startedAt": "2024-01-15T08:00:00Z",
    "finishedAt": "2024-01-15T09:30:00Z",
    "createdAt": "2024-01-15T08:00:00Z",
    "updatedAt": "2024-01-15T09:30:00Z"
  }
}
```

**Errors**:
- `404 Not Found`: Session not found
- `409 Conflict`: Session already completed or abandoned
- `400 Bad Request`: Session not in IN_PROGRESS state

#### POST /workouts/{id}/abandon

Abandon a workout session.

**Auth**: Owner/Admin

**Request Body**: None required

**Response** `200 OK`:
```json
{
  "data": {
    "id": "session-uuid",
    "userProgramStateId": "enrollment-uuid",
    "weekNumber": 1,
    "dayIndex": 0,
    "status": "ABANDONED",
    "startedAt": "2024-01-15T08:00:00Z",
    "finishedAt": "2024-01-15T08:30:00Z",
    "createdAt": "2024-01-15T08:00:00Z",
    "updatedAt": "2024-01-15T08:30:00Z"
  }
}
```

**Errors**:
- `404 Not Found`: Session not found
- `409 Conflict`: Session already completed or abandoned
- `400 Bad Request`: Session not in IN_PROGRESS state

#### GET /users/{userId}/workouts

List a user's workout history with pagination.

**Auth**: Owner/Admin

**Query Parameters**:
| Parameter | Type | Description |
|-----------|------|-------------|
| `limit` | int | Number of items to return (default: 20, max: 100) |
| `offset` | int | Number of items to skip (default: 0) |
| `status` | string | Filter by status: "IN_PROGRESS", "COMPLETED", or "ABANDONED" |

**Response** `200 OK`:
```json
{
  "data": [
    {
      "id": "session-uuid",
      "userProgramStateId": "enrollment-uuid",
      "weekNumber": 1,
      "dayIndex": 0,
      "status": "COMPLETED",
      "startedAt": "2024-01-15T08:00:00Z",
      "finishedAt": "2024-01-15T09:30:00Z",
      "createdAt": "2024-01-15T08:00:00Z",
      "updatedAt": "2024-01-15T09:30:00Z"
    }
  ],
  "meta": {
    "total": 10,
    "limit": 20,
    "offset": 0,
    "hasMore": false
  }
}
```

#### GET /users/{userId}/workouts/current

Get the user's current in-progress workout session if any.

**Auth**: Owner/Admin

**Response** `200 OK`: WorkoutSession object (same format as GET /workouts/{id})

**Errors**:
- `404 Not Found`: No active workout session

---

### Enrollment State Management

Manage enrollment state transitions for cycles and weeks.

#### POST /users/{userId}/enrollment/next-cycle

Start a new cycle when the enrollment is in BETWEEN_CYCLES state.

**Auth**: Owner/Admin

**Request Body**: None required

**Response** `200 OK`:
```json
{
  "data": {
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
      "currentCycleIteration": 2,
      "currentDayIndex": null
    },
    "enrollmentStatus": "ACTIVE",
    "cycleStatus": "PENDING",
    "weekStatus": "PENDING",
    "currentWorkoutSession": null,
    "enrolledAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-29T10:00:00Z"
  }
}
```

**Errors**:
- `404 Not Found`: User not enrolled
- `400 Bad Request`: Enrollment not in BETWEEN_CYCLES state

#### POST /users/{userId}/enrollment/advance-week

Advance to the next week in the cycle. If at the final week, transitions enrollment to BETWEEN_CYCLES.

**Auth**: Owner/Admin

**Request Body**: None required

**Response** `200 OK`:
```json
{
  "data": {
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
      "currentWeek": 2,
      "currentCycleIteration": 1,
      "currentDayIndex": null
    },
    "enrollmentStatus": "ACTIVE",
    "cycleStatus": "IN_PROGRESS",
    "weekStatus": "PENDING",
    "currentWorkoutSession": null,
    "enrolledAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-22T10:00:00Z"
  }
}
```

**Note**: When advancing from the final week of a cycle:
- `enrollmentStatus` transitions to `BETWEEN_CYCLES`
- `cycleStatus` transitions to `COMPLETED`
- `weekStatus` transitions to `COMPLETED`
- User must call `POST /users/{userId}/enrollment/next-cycle` to start the next cycle

**Errors**:
- `404 Not Found`: User not enrolled
- `400 Bad Request`: Enrollment not in ACTIVE state

---

### Updated Enrollment Response

The enrollment response (`GET /users/{userId}/program`) now includes state machine status fields:

```json
{
  "data": {
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
    "enrollmentStatus": "ACTIVE",
    "cycleStatus": "IN_PROGRESS",
    "weekStatus": "IN_PROGRESS",
    "currentWorkoutSession": {
      "id": "session-uuid",
      "weekNumber": 1,
      "dayIndex": 0,
      "status": "IN_PROGRESS",
      "startedAt": "2024-01-15T08:00:00Z"
    },
    "enrolledAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-15T08:00:00Z"
  }
}
```

**Status Fields**:

| Field | Values | Description |
|-------|--------|-------------|
| `enrollmentStatus` | `ACTIVE`, `BETWEEN_CYCLES`, `QUIT` | Overall enrollment state |
| `cycleStatus` | `PENDING`, `IN_PROGRESS`, `COMPLETED` | Current cycle state |
| `weekStatus` | `PENDING`, `IN_PROGRESS`, `COMPLETED` | Current week state |

**Workflow Notes**:
- `enrollmentStatus: ACTIVE` - User can start workouts
- `enrollmentStatus: BETWEEN_CYCLES` - User must call next-cycle to continue
- `currentWorkoutSession` - Present only if user has an in-progress workout

---

## Canonical Programs

PowerPro ships with pre-configured canonical programs that users can enroll in immediately. These programs represent proven, real-world powerlifting methodologies.

### Available Canonical Programs

| Slug | Name | Days/Week | Cycle Length | Level | Progression Model |
|------|------|-----------|--------------|-------|-------------------|
| `starting-strength` | Starting Strength | 3 | 1 week | Novice | Linear (+5lb/+10lb per session) |
| `texas-method` | Texas Method | 3 | 1 week | Intermediate | Weekly periodization |
| `531` | Wendler 5/3/1 | 4 | 4 weeks | Intermediate | Monthly cycles with AMRAP |
| `gzclp` | GZCLP | 4 | 1 week | Beginner/Intermediate | Tiered progression (T1/T2) |

### Program Details

#### Starting Strength (`starting-strength`)

The classic novice linear progression program by Mark Rippetoe. Features an A/B day rotation with compound barbell movements.

- **Structure**: A/B alternating workouts (Mon/Wed/Fri pattern)
- **Workout A**: Squat 3x5, Bench Press 3x5, Deadlift 1x5
- **Workout B**: Squat 3x5, Overhead Press 3x5, Power Clean 5x3
- **Progression**: Lower body +10lb/session, Upper body +5lb/session
- **Best for**: Complete beginners, rapid strength gains

#### Texas Method (`texas-method`)

An intermediate weekly periodization program with volume, recovery, and intensity phases.

- **Structure**: 3 distinct training days per week
- **Volume Day** (Monday): 5x5 at 90% - accumulate volume
- **Recovery Day** (Wednesday): Light weights at 72% - active recovery
- **Intensity Day** (Friday): Heavy singles/triples at 100%+ - set new PRs
- **Progression**: Weekly weight increases (+5lb upper, +5lb lower per week)
- **Best for**: Lifters who have exhausted linear progression

#### Wendler 5/3/1 (`531`)

Jim Wendler's submaximal training program with monthly progression cycles.

- **Structure**: 4-week mesocycle with 4 training days per week
- **Week 1** (5s): 65% x5, 75% x5, 85% x5+ (AMRAP)
- **Week 2** (3s): 70% x3, 80% x3, 90% x3+ (AMRAP)
- **Week 3** (5/3/1): 75% x5, 85% x3, 95% x1+ (AMRAP)
- **Week 4** (Deload): 40% x5, 50% x5, 60% x5
- **Progression**: +10lb lower body, +5lb upper body per cycle (after week 4)
- **Training Max**: Uses 90% of actual 1RM for calculations
- **Best for**: Sustainable long-term strength building

#### GZCLP (`gzclp`)

Cody Lefever's linear progression using the GZCL tiered methodology.

- **Structure**: 4 training days per week with T1/T2 tier system
- **T1 (Main Lift)**: 5x3+ at 85% - heavy work with AMRAP final set
- **T2 (Secondary Lift)**: 3x10 at 65% - volume work
- **Day Pairings**:
  - Day 1: T1 Squat, T2 Bench
  - Day 2: T1 OHP, T2 Deadlift
  - Day 3: T1 Bench, T2 Squat
  - Day 4: T1 Deadlift, T2 OHP
- **Progression**: T1 lower +5lb, T1 upper +2.5lb, T2 all +2.5lb per session
- **Failure Protocol**: Progress through stages (5x3 → 6x2 → 10x1) before resetting
- **Best for**: Beginners wanting structured tier-based training

### Identifying Canonical Programs

Canonical programs are identified by their slugs. Use `GET /programs/by-slug/{slug}` with one of the canonical slugs to retrieve the program:

```bash
# Get Starting Strength program
GET /programs/by-slug/starting-strength

# Get 5/3/1 program
GET /programs/by-slug/531
```

All canonical programs appear in the regular `GET /programs` endpoint alongside any user-created programs.

### Canonical Program Restrictions

1. **Read-only**: Canonical programs cannot be modified via PUT/DELETE endpoints
2. **Enrollable**: Users can enroll in any canonical program via `POST /users/{userId}/program`
3. **System-owned**: Canonical programs are not associated with any user account

### Enrolling in a Canonical Program

To enroll a user in a canonical program:

1. Get the program ID by slug:
```bash
GET /programs/by-slug/starting-strength
```

2. Enroll the user:
```bash
POST /users/{userId}/program
Content-Type: application/json

{
  "programId": "{program-id-from-step-1}"
}
```

3. Set up required training maxes for the lifts used in the program:
```bash
POST /users/{userId}/lift-maxes
Content-Type: application/json

{
  "liftId": "{squat-lift-id}",
  "type": "TRAINING_MAX",
  "value": 315.0
}
```

4. Start generating workouts:
```bash
GET /users/{userId}/workout
```

### Required Lift Maxes by Program

Each program requires training maxes to be set for specific lifts:

| Program | Required Lift Maxes |
|---------|---------------------|
| `starting-strength` | Squat, Bench Press, Deadlift, Overhead Press, Power Clean |
| `texas-method` | Squat, Bench Press, Deadlift, Overhead Press, Power Clean |
| `531` | Squat, Bench Press, Deadlift, Overhead Press |
| `gzclp` | Squat, Bench Press, Deadlift, Overhead Press |

If a required lift max is missing, workout generation will return a `422 Unprocessable Entity` error.

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
| `user not enrolled in a program` | 404 | User must enroll before starting workouts |
| `cannot perform action in current enrollment state` | 400 | Invalid state transition (e.g., starting workout when BETWEEN_CYCLES) |
| `workout already in progress` | 409 | User already has an active workout session |
| `session already completed` | 409 | Cannot finish/abandon an already completed session |
| `session not in progress` | 400 | Session must be IN_PROGRESS to finish/abandon |

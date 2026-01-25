# PowerPro API Example Requests

This document provides copy-paste ready example requests for all PowerPro API endpoints. All examples use realistic data and can be executed directly against a running API instance.

## Base URL and Authentication

```bash
# Base URL (adjust port as needed)
BASE_URL="http://localhost:8080"

# For session-based auth (recommended), use the token from login response:
AUTH_TOKEN="your-session-token-from-login"

# Standard headers for authenticated requests (using session token)
AUTH_HEADERS="-H 'Authorization: Bearer $AUTH_TOKEN'"

# Alternative: User ID for development/testing (bypasses session auth)
USER_ID="550e8400-e29b-41d4-a716-446655440000"
DEV_AUTH_HEADERS="-H 'X-User-ID: $USER_ID'"

# Admin headers (for admin-only endpoints)
ADMIN_HEADERS="-H 'Authorization: Bearer $AUTH_TOKEN' -H 'X-Admin: true'"
```

---

## Health Check

### GET /health

Check API health status.

```bash
curl http://localhost:8080/health
```

---

## Authentication

User registration, login, and session management.

### POST /auth/register - Register a new user

```bash
# Register a new user
curl -X POST "http://localhost:8080/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "password": "securepassword123",
    "name": "John Doe"
  }'

# Register without optional name
curl -X POST "http://localhost:8080/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "minimal@example.com",
    "password": "securepassword123"
  }'
```

### POST /auth/login - Authenticate and get session token

```bash
# Login and get session token
curl -X POST "http://localhost:8080/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123"
  }'
# Response contains token to use in Authorization header for subsequent requests
```

### POST /auth/logout - End current session

```bash
# Logout (invalidates the session token)
curl -X POST "http://localhost:8080/auth/logout" \
  -H "Authorization: Bearer your-session-token"
```

---

## User Profile

Manage user profile information.

### GET /users/{userId}/profile - Get user profile

```bash
# Get your own profile
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/profile" \
  -H "Authorization: Bearer your-session-token"

# Admin viewing another user's profile
curl -X GET "http://localhost:8080/users/another-user-uuid/profile" \
  -H "Authorization: Bearer admin-session-token" \
  -H "X-Admin: true"
```

### PUT /users/{userId}/profile - Update user profile

```bash
# Update name only
curl -X PUT "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/profile" \
  -H "Authorization: Bearer your-session-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe"
  }'

# Update weight unit preference
curl -X PUT "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/profile" \
  -H "Authorization: Bearer your-session-token" \
  -H "Content-Type: application/json" \
  -d '{
    "weightUnit": "kg"
  }'

# Update multiple fields
curl -X PUT "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/profile" \
  -H "Authorization: Bearer your-session-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "weightUnit": "kg"
  }'
```

---

## Dashboard

User dashboard with aggregated data.

### GET /users/{id}/dashboard - Get user dashboard

```bash
# Get your dashboard (owner-only endpoint)
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/dashboard" \
  -H "Authorization: Bearer your-session-token"
```

---

## Lifts

Manage exercises (lifts) in the system.

### GET /lifts - List all lifts

```bash
# Basic request
curl -X GET "http://localhost:8080/lifts" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# With pagination
curl -X GET "http://localhost:8080/lifts?page=1&pageSize=10" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# Filter competition lifts only
curl -X GET "http://localhost:8080/lifts?is_competition_lift=true" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# Sorted by name descending
curl -X GET "http://localhost:8080/lifts?sortBy=name&sortOrder=desc" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### GET /lifts/{id} - Get lift by ID

```bash
curl -X GET "http://localhost:8080/lifts/a1b2c3d4-e5f6-7890-abcd-ef1234567890" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### GET /lifts/by-slug/{slug} - Get lift by slug

```bash
curl -X GET "http://localhost:8080/lifts/by-slug/squat" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### POST /lifts - Create a lift (Admin)

```bash
# Create a competition lift (Squat)
curl -X POST "http://localhost:8080/lifts" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Squat",
    "slug": "squat",
    "isCompetitionLift": true
  }'

# Create a competition lift (Bench Press)
curl -X POST "http://localhost:8080/lifts" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bench Press",
    "slug": "bench-press",
    "isCompetitionLift": true
  }'

# Create an accessory lift with parent reference
curl -X POST "http://localhost:8080/lifts" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Pause Squat",
    "slug": "pause-squat",
    "isCompetitionLift": false,
    "parentLiftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  }'
```

### PUT /lifts/{id} - Update a lift (Admin)

```bash
# Update lift name and slug
curl -X PUT "http://localhost:8080/lifts/a1b2c3d4-e5f6-7890-abcd-ef1234567890" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Back Squat",
    "slug": "back-squat"
  }'

# Remove parent lift reference
curl -X PUT "http://localhost:8080/lifts/b2c3d4e5-f6a7-8901-bcde-f12345678901" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "clearParentLift": true
  }'
```

### DELETE /lifts/{id} - Delete a lift (Admin)

```bash
curl -X DELETE "http://localhost:8080/lifts/a1b2c3d4-e5f6-7890-abcd-ef1234567890" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true"
```

---

## Lift Maxes

Manage user's personal records (1RM, training max).

### GET /users/{userId}/lift-maxes - List user's lift maxes

```bash
# List all lift maxes for a user
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/lift-maxes" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# Filter by lift ID
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/lift-maxes?lift_id=a1b2c3d4-e5f6-7890-abcd-ef1234567890" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# Filter by type (TRAINING_MAX)
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/lift-maxes?type=TRAINING_MAX" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# With pagination
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/lift-maxes?page=1&pageSize=5" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### GET /users/{userId}/lift-maxes/current - Get current lift max

```bash
# Get current training max for squat
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/lift-maxes/current?lift=a1b2c3d4-e5f6-7890-abcd-ef1234567890&type=TRAINING_MAX" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# Get current 1RM for bench press
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/lift-maxes/current?lift=b2c3d4e5-f6a7-8901-bcde-f12345678901&type=ONE_RM" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### GET /lift-maxes/{id} - Get lift max by ID

```bash
curl -X GET "http://localhost:8080/lift-maxes/c3d4e5f6-a7b8-9012-cdef-123456789012" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### GET /lift-maxes/{id}/convert - Convert lift max

```bash
# Convert 1RM to Training Max (default 90%)
curl -X GET "http://localhost:8080/lift-maxes/c3d4e5f6-a7b8-9012-cdef-123456789012/convert?to_type=TRAINING_MAX" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# Convert with custom percentage (85%)
curl -X GET "http://localhost:8080/lift-maxes/c3d4e5f6-a7b8-9012-cdef-123456789012/convert?to_type=TRAINING_MAX&percentage=85" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# Convert Training Max to 1RM
curl -X GET "http://localhost:8080/lift-maxes/c3d4e5f6-a7b8-9012-cdef-123456789012/convert?to_type=ONE_RM&percentage=90" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### POST /users/{userId}/lift-maxes - Create lift max

```bash
# Create a training max for squat (315 lbs)
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/lift-maxes" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "type": "TRAINING_MAX",
    "value": 315.0,
    "effectiveDate": "2024-01-15T00:00:00Z"
  }'

# Create a 1RM for bench press (225 lbs)
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/lift-maxes" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
    "type": "ONE_RM",
    "value": 225.0
  }'

# Create a training max for deadlift (405 lbs)
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

### PUT /lift-maxes/{id} - Update lift max

```bash
# Update the value
curl -X PUT "http://localhost:8080/lift-maxes/c3d4e5f6-a7b8-9012-cdef-123456789012" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "value": 320.0
  }'

# Update value and effective date
curl -X PUT "http://localhost:8080/lift-maxes/c3d4e5f6-a7b8-9012-cdef-123456789012" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "value": 325.0,
    "effectiveDate": "2024-02-01T00:00:00Z"
  }'
```

### DELETE /lift-maxes/{id} - Delete lift max

```bash
curl -X DELETE "http://localhost:8080/lift-maxes/c3d4e5f6-a7b8-9012-cdef-123456789012" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

---

## Prescriptions

Manage exercise prescriptions (what to do for a single exercise slot).

### GET /prescriptions - List all prescriptions

```bash
# Basic request
curl -X GET "http://localhost:8080/prescriptions" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# Filter by lift
curl -X GET "http://localhost:8080/prescriptions?lift_id=a1b2c3d4-e5f6-7890-abcd-ef1234567890" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# Sorted by order
curl -X GET "http://localhost:8080/prescriptions?sortBy=order&sortOrder=asc" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### GET /prescriptions/{id} - Get prescription by ID

```bash
curl -X GET "http://localhost:8080/prescriptions/e5f6a7b8-c9d0-1234-efgh-567890123456" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### POST /prescriptions - Create a prescription (Admin)

```bash
# Fixed 5x5 at 85% of training max (classic strength prescription)
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "TRAINING_MAX",
      "percentage": 85.0,
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
    "restSeconds": 180
  }'

# AMRAP set at 95% (Wendler-style top set)
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
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
    "order": 1,
    "notes": "Top set - as many reps as possible",
    "restSeconds": 300
  }'

# Ramping warmup sets
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "TRAINING_MAX",
      "percentage": 100.0
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
    "notes": "Warmup - build up to working weight"
  }'

# Boring But Big accessory (5x10 at 50%)
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
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
    "notes": "BBB volume work - focus on quality reps",
    "restSeconds": 90
  }'
```

### PUT /prescriptions/{id} - Update a prescription (Admin)

```bash
# Update percentage and rest time
curl -X PUT "http://localhost:8080/prescriptions/e5f6a7b8-c9d0-1234-efgh-567890123456" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "loadStrategy": {
      "type": "PERCENT_OF",
      "maxType": "TRAINING_MAX",
      "percentage": 90.0
    },
    "restSeconds": 240
  }'

# Update notes
curl -X PUT "http://localhost:8080/prescriptions/e5f6a7b8-c9d0-1234-efgh-567890123456" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "notes": "Pause at the bottom for 2 seconds"
  }'
```

### DELETE /prescriptions/{id} - Delete a prescription (Admin)

```bash
curl -X DELETE "http://localhost:8080/prescriptions/e5f6a7b8-c9d0-1234-efgh-567890123456" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true"
```

### POST /prescriptions/{id}/resolve - Resolve prescription to concrete sets

```bash
curl -X POST "http://localhost:8080/prescriptions/e5f6a7b8-c9d0-1234-efgh-567890123456/resolve" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "550e8400-e29b-41d4-a716-446655440000"
  }'
```

### POST /prescriptions/resolve-batch - Resolve multiple prescriptions

```bash
curl -X POST "http://localhost:8080/prescriptions/resolve-batch" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionIds": [
      "e5f6a7b8-c9d0-1234-efgh-567890123456",
      "f6a7b8c9-d0e1-2345-fghi-678901234567"
    ],
    "userId": "550e8400-e29b-41d4-a716-446655440000"
  }'
```

---

## Days

Manage training days.

### GET /days - List all days

```bash
# Basic request
curl -X GET "http://localhost:8080/days" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# Filter by program
curl -X GET "http://localhost:8080/days?program_id=g7h8i9j0-k1l2-3456-mnop-789012345678" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# With pagination and sorting
curl -X GET "http://localhost:8080/days?page=1&pageSize=10&sortBy=name&sortOrder=asc" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### GET /days/{id} - Get day by ID (includes prescriptions)

```bash
curl -X GET "http://localhost:8080/days/h8i9j0k1-l2m3-4567-nopq-890123456789" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### GET /days/by-slug/{slug} - Get day by slug

```bash
# Get by slug
curl -X GET "http://localhost:8080/days/by-slug/squat-day" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# Scoped to a specific program
curl -X GET "http://localhost:8080/days/by-slug/squat-day?program_id=g7h8i9j0-k1l2-3456-mnop-789012345678" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### POST /days - Create a day (Admin)

```bash
# Squat-focused day
curl -X POST "http://localhost:8080/days" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Squat Day",
    "slug": "squat-day",
    "metadata": {
      "focus": "lower body",
      "primaryLift": "squat"
    },
    "programId": "g7h8i9j0-k1l2-3456-mnop-789012345678"
  }'

# Bench-focused day
curl -X POST "http://localhost:8080/days" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bench Day",
    "slug": "bench-day",
    "metadata": {
      "focus": "upper body",
      "primaryLift": "bench"
    },
    "programId": "g7h8i9j0-k1l2-3456-mnop-789012345678"
  }'

# Deadlift-focused day
curl -X POST "http://localhost:8080/days" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Deadlift Day",
    "slug": "deadlift-day",
    "metadata": {
      "focus": "posterior chain",
      "primaryLift": "deadlift"
    },
    "programId": "g7h8i9j0-k1l2-3456-mnop-789012345678"
  }'
```

### PUT /days/{id} - Update a day (Admin)

```bash
# Update name and metadata
curl -X PUT "http://localhost:8080/days/h8i9j0k1-l2m3-4567-nopq-890123456789" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Heavy Squat Day",
    "metadata": {
      "focus": "max effort",
      "primaryLift": "squat",
      "intensity": "high"
    }
  }'

# Clear metadata
curl -X PUT "http://localhost:8080/days/h8i9j0k1-l2m3-4567-nopq-890123456789" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "clearMetadata": true
  }'
```

### DELETE /days/{id} - Delete a day (Admin)

```bash
curl -X DELETE "http://localhost:8080/days/h8i9j0k1-l2m3-4567-nopq-890123456789" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true"
```

### POST /days/{id}/prescriptions - Add prescription to day (Admin)

```bash
# Add first prescription (order 1)
curl -X POST "http://localhost:8080/days/h8i9j0k1-l2m3-4567-nopq-890123456789/prescriptions" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionId": "e5f6a7b8-c9d0-1234-efgh-567890123456",
    "order": 1
  }'

# Add second prescription (order 2)
curl -X POST "http://localhost:8080/days/h8i9j0k1-l2m3-4567-nopq-890123456789/prescriptions" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionId": "f6a7b8c9-d0e1-2345-fghi-678901234567",
    "order": 2
  }'
```

### DELETE /days/{id}/prescriptions/{prescriptionId} - Remove prescription from day (Admin)

```bash
curl -X DELETE "http://localhost:8080/days/h8i9j0k1-l2m3-4567-nopq-890123456789/prescriptions/e5f6a7b8-c9d0-1234-efgh-567890123456" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true"
```

### PUT /days/{id}/prescriptions/reorder - Reorder prescriptions in day (Admin)

```bash
curl -X PUT "http://localhost:8080/days/h8i9j0k1-l2m3-4567-nopq-890123456789/prescriptions/reorder" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "prescriptionIds": [
      "f6a7b8c9-d0e1-2345-fghi-678901234567",
      "e5f6a7b8-c9d0-1234-efgh-567890123456",
      "a7b8c9d0-e1f2-3456-ghij-789012345678"
    ]
  }'
```

---

## Weeks

Manage training weeks.

### GET /weeks - List all weeks

```bash
# Basic request
curl -X GET "http://localhost:8080/weeks" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# Sorted by week number
curl -X GET "http://localhost:8080/weeks?sortBy=week_number&sortOrder=asc" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### GET /weeks/{id} - Get week by ID (includes days)

```bash
curl -X GET "http://localhost:8080/weeks/i9j0k1l2-m3n4-5678-opqr-901234567890" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### POST /weeks - Create a week (Admin)

```bash
# Week 1 of a 4-week cycle
curl -X POST "http://localhost:8080/weeks" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "cycleId": "j0k1l2m3-n4o5-6789-pqrs-012345678901",
    "weekNumber": 1,
    "name": "Week 1 - 5s"
  }'

# Week 2
curl -X POST "http://localhost:8080/weeks" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "cycleId": "j0k1l2m3-n4o5-6789-pqrs-012345678901",
    "weekNumber": 2,
    "name": "Week 2 - 3s"
  }'

# Week 3
curl -X POST "http://localhost:8080/weeks" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "cycleId": "j0k1l2m3-n4o5-6789-pqrs-012345678901",
    "weekNumber": 3,
    "name": "Week 3 - 5/3/1"
  }'

# Week 4 (Deload)
curl -X POST "http://localhost:8080/weeks" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "cycleId": "j0k1l2m3-n4o5-6789-pqrs-012345678901",
    "weekNumber": 4,
    "name": "Week 4 - Deload"
  }'
```

### PUT /weeks/{id} - Update a week (Admin)

```bash
curl -X PUT "http://localhost:8080/weeks/i9j0k1l2-m3n4-5678-opqr-901234567890" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Week 1 - Volume"
  }'
```

### DELETE /weeks/{id} - Delete a week (Admin)

```bash
curl -X DELETE "http://localhost:8080/weeks/i9j0k1l2-m3n4-5678-opqr-901234567890" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true"
```

### POST /weeks/{id}/days - Add day to week (Admin)

```bash
# Add Monday (position 0)
curl -X POST "http://localhost:8080/weeks/i9j0k1l2-m3n4-5678-opqr-901234567890/days" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "dayId": "h8i9j0k1-l2m3-4567-nopq-890123456789",
    "position": 0
  }'

# Add Wednesday (position 1)
curl -X POST "http://localhost:8080/weeks/i9j0k1l2-m3n4-5678-opqr-901234567890/days" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "dayId": "k1l2m3n4-o5p6-7890-qrst-123456789012",
    "position": 1
  }'

# Add Friday (position 2)
curl -X POST "http://localhost:8080/weeks/i9j0k1l2-m3n4-5678-opqr-901234567890/days" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "dayId": "l2m3n4o5-p6q7-8901-rstu-234567890123",
    "position": 2
  }'
```

### DELETE /weeks/{id}/days/{dayId} - Remove day from week (Admin)

```bash
curl -X DELETE "http://localhost:8080/weeks/i9j0k1l2-m3n4-5678-opqr-901234567890/days/h8i9j0k1-l2m3-4567-nopq-890123456789" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true"
```

---

## Cycles

Manage training cycles.

### GET /cycles - List all cycles

```bash
curl -X GET "http://localhost:8080/cycles" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### GET /cycles/{id} - Get cycle by ID

```bash
curl -X GET "http://localhost:8080/cycles/j0k1l2m3-n4o5-6789-pqrs-012345678901" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### POST /cycles - Create a cycle (Admin)

```bash
# 4-week cycle (typical for 5/3/1)
curl -X POST "http://localhost:8080/cycles" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "4 Week 5/3/1 Cycle",
    "lengthWeeks": 4
  }'

# 1-week cycle (typical for linear progression)
curl -X POST "http://localhost:8080/cycles" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Weekly Linear Cycle",
    "lengthWeeks": 1
  }'

# 3-week cycle (for Greg Nuckols programs)
curl -X POST "http://localhost:8080/cycles" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "3 Week Undulating Cycle",
    "lengthWeeks": 3
  }'
```

### PUT /cycles/{id} - Update a cycle (Admin)

```bash
curl -X PUT "http://localhost:8080/cycles/j0k1l2m3-n4o5-6789-pqrs-012345678901" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "4 Week BBB Cycle",
    "lengthWeeks": 4
  }'
```

### DELETE /cycles/{id} - Delete a cycle (Admin)

```bash
curl -X DELETE "http://localhost:8080/cycles/j0k1l2m3-n4o5-6789-pqrs-012345678901" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true"
```

---

## Weekly Lookups

Manage weekly lookup tables (varying parameters by week).

### GET /weekly-lookups - List all weekly lookups

```bash
curl -X GET "http://localhost:8080/weekly-lookups" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### GET /weekly-lookups/{id} - Get weekly lookup by ID

```bash
curl -X GET "http://localhost:8080/weekly-lookups/m3n4o5p6-q7r8-9012-stuv-345678901234" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### POST /weekly-lookups - Create a weekly lookup (Admin)

```bash
# 5/3/1 top set percentages
curl -X POST "http://localhost:8080/weekly-lookups" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "5/3/1 Top Set Percentages",
    "entries": {
      "1": 85,
      "2": 90,
      "3": 95,
      "4": 60
    }
  }'

# 5/3/1 working set percentages (set 1 of each week)
curl -X POST "http://localhost:8080/weekly-lookups" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "5/3/1 Week 1 Percentages",
    "entries": {
      "1": 65,
      "2": 75,
      "3": 85
    }
  }'

# Greg Nuckols 3-week cycle percentages
curl -X POST "http://localhost:8080/weekly-lookups" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Nuckols 3-Week Cycle",
    "entries": {
      "1": 80,
      "2": 85,
      "3": 90
    }
  }'
```

### PUT /weekly-lookups/{id} - Update a weekly lookup (Admin)

```bash
curl -X PUT "http://localhost:8080/weekly-lookups/m3n4o5p6-q7r8-9012-stuv-345678901234" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "5/3/1 Top Set Percentages (Modified)",
    "entries": {
      "1": 80,
      "2": 85,
      "3": 90,
      "4": 50
    }
  }'
```

### DELETE /weekly-lookups/{id} - Delete a weekly lookup (Admin)

```bash
curl -X DELETE "http://localhost:8080/weekly-lookups/m3n4o5p6-q7r8-9012-stuv-345678901234" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true"
```

---

## Daily Lookups

Manage daily lookup tables (varying parameters by day).

### GET /daily-lookups - List all daily lookups

```bash
curl -X GET "http://localhost:8080/daily-lookups" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### GET /daily-lookups/{id} - Get daily lookup by ID

```bash
curl -X GET "http://localhost:8080/daily-lookups/n4o5p6q7-r8s9-0123-uvwx-456789012345" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### POST /daily-lookups - Create a daily lookup (Admin)

```bash
# Heavy/Light/Medium daily variation (Bill Starr style)
curl -X POST "http://localhost:8080/daily-lookups" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Heavy/Light/Medium",
    "entries": {
      "monday": 100,
      "wednesday": 80,
      "friday": 90
    }
  }'

# Daily undulation for high frequency programs
curl -X POST "http://localhost:8080/daily-lookups" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Daily Undulating Periodization",
    "entries": {
      "day1": 85,
      "day2": 70,
      "day3": 90,
      "day4": 75,
      "day5": 95
    }
  }'
```

### PUT /daily-lookups/{id} - Update a daily lookup (Admin)

```bash
curl -X PUT "http://localhost:8080/daily-lookups/n4o5p6q7-r8s9-0123-uvwx-456789012345" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Modified H/L/M",
    "entries": {
      "monday": 100,
      "wednesday": 70,
      "friday": 85
    }
  }'
```

### DELETE /daily-lookups/{id} - Delete a daily lookup (Admin)

```bash
curl -X DELETE "http://localhost:8080/daily-lookups/n4o5p6q7-r8s9-0123-uvwx-456789012345" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true"
```

---

## Programs

Manage training programs.

### GET /programs - List all programs

```bash
# Basic request
curl -X GET "http://localhost:8080/programs" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# With pagination and sorting
curl -X GET "http://localhost:8080/programs?page=1&pageSize=10&sortBy=name&sortOrder=asc" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### GET /programs/{id} - Get program by ID (with embedded cycle details)

```bash
curl -X GET "http://localhost:8080/programs/g7h8i9j0-k1l2-3456-mnop-789012345678" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### POST /programs - Create a program (Admin)

```bash
# Wendler 5/3/1 BBB
curl -X POST "http://localhost:8080/programs" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Wendler 5/3/1 BBB",
    "slug": "wendler-531-bbb",
    "description": "The classic Boring But Big variant of 5/3/1 featuring high-volume accessory work at 50% of training max.",
    "cycleId": "j0k1l2m3-n4o5-6789-pqrs-012345678901",
    "weeklyLookupId": "m3n4o5p6-q7r8-9012-stuv-345678901234",
    "defaultRounding": 5.0
  }'

# Starting Strength (linear progression)
curl -X POST "http://localhost:8080/programs" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Starting Strength",
    "slug": "starting-strength",
    "description": "Classic linear progression program for novice lifters featuring 3x5 working sets with weight increases every session.",
    "cycleId": "o5p6q7r8-s9t0-1234-vwxy-567890123456",
    "defaultRounding": 5.0
  }'

# Bill Starr 5x5 (with daily lookup for H/L/M)
curl -X POST "http://localhost:8080/programs" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bill Starr 5x5",
    "slug": "bill-starr-5x5",
    "description": "Heavy/Light/Medium 5x5 program with ramping sets and daily intensity variation.",
    "cycleId": "o5p6q7r8-s9t0-1234-vwxy-567890123456",
    "dailyLookupId": "n4o5p6q7-r8s9-0123-uvwx-456789012345",
    "defaultRounding": 5.0
  }'

# Greg Nuckols High Frequency
curl -X POST "http://localhost:8080/programs" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Greg Nuckols High Frequency",
    "slug": "nuckols-high-frequency",
    "description": "High frequency squat/bench/deadlift program with daily undulation and 3-week cycles.",
    "cycleId": "p6q7r8s9-t0u1-2345-wxyz-678901234567",
    "weeklyLookupId": "q7r8s9t0-u1v2-3456-xyza-789012345678",
    "dailyLookupId": "r8s9t0u1-v2w3-4567-yzab-890123456789",
    "defaultRounding": 2.5
  }'
```

### PUT /programs/{id} - Update a program (Admin)

```bash
# Update description and rounding
curl -X PUT "http://localhost:8080/programs/g7h8i9j0-k1l2-3456-mnop-789012345678" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Updated description with more detail about the program structure.",
    "defaultRounding": 2.5
  }'
```

### DELETE /programs/{id} - Delete a program (Admin)

```bash
curl -X DELETE "http://localhost:8080/programs/g7h8i9j0-k1l2-3456-mnop-789012345678" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true"
```

---

## Progressions

Manage progression rules (how to increase weights over time).

### GET /progressions - List all progressions

```bash
# Basic request
curl -X GET "http://localhost:8080/progressions" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# Filter by type
curl -X GET "http://localhost:8080/progressions?type=LINEAR" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### GET /progressions/{id} - Get progression by ID

```bash
curl -X GET "http://localhost:8080/progressions/s9t0u1v2-w3x4-5678-bcde-901234567890" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### POST /progressions - Create a progression (Admin)

```bash
# Linear progression +5lb after each session (Starting Strength style)
curl -X POST "http://localhost:8080/progressions" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Linear +5lb per session",
    "type": "LINEAR",
    "parameters": {
      "increment": 5.0,
      "maxType": "TRAINING_MAX",
      "triggerType": "AFTER_SESSION"
    }
  }'

# Linear progression +2.5lb (for upper body or slower progression)
curl -X POST "http://localhost:8080/progressions" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Linear +2.5lb per session",
    "type": "LINEAR",
    "parameters": {
      "increment": 2.5,
      "maxType": "TRAINING_MAX",
      "triggerType": "AFTER_SESSION"
    }
  }'

# Weekly linear progression
curl -X POST "http://localhost:8080/progressions" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Weekly +5lb",
    "type": "LINEAR",
    "parameters": {
      "increment": 5.0,
      "maxType": "TRAINING_MAX",
      "triggerType": "AFTER_WEEK"
    }
  }'

# Cycle-based progression (5/3/1 style - +5lb at end of cycle)
curl -X POST "http://localhost:8080/progressions" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Cycle +5lb (lower body)",
    "type": "CYCLE",
    "parameters": {
      "increment": 5.0,
      "maxType": "TRAINING_MAX"
    }
  }'

# Cycle-based progression for upper body (+2.5lb)
curl -X POST "http://localhost:8080/progressions" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Cycle +2.5lb (upper body)",
    "type": "CYCLE",
    "parameters": {
      "increment": 2.5,
      "maxType": "TRAINING_MAX"
    }
  }'
```

### PUT /progressions/{id} - Update a progression (Admin)

```bash
curl -X PUT "http://localhost:8080/progressions/s9t0u1v2-w3x4-5678-bcde-901234567890" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Modified Linear +10lb",
    "parameters": {
      "increment": 10.0,
      "maxType": "TRAINING_MAX",
      "triggerType": "AFTER_SESSION"
    }
  }'
```

### DELETE /progressions/{id} - Delete a progression (Admin)

```bash
curl -X DELETE "http://localhost:8080/progressions/s9t0u1v2-w3x4-5678-bcde-901234567890" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true"
```

---

## Program Progressions

Configure which progressions apply to which programs/lifts.

### GET /programs/{programId}/progressions - List program progressions

```bash
curl -X GET "http://localhost:8080/programs/g7h8i9j0-k1l2-3456-mnop-789012345678/progressions" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### GET /programs/{programId}/progressions/{configId} - Get specific configuration

```bash
curl -X GET "http://localhost:8080/programs/g7h8i9j0-k1l2-3456-mnop-789012345678/progressions/t0u1v2w3-x4y5-6789-cdef-012345678901" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### POST /programs/{programId}/progressions - Create program progression (Admin)

```bash
# Apply cycle progression to all lifts (5/3/1 default)
curl -X POST "http://localhost:8080/programs/g7h8i9j0-k1l2-3456-mnop-789012345678/progressions" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "s9t0u1v2-w3x4-5678-bcde-901234567890",
    "priority": 1
  }'

# Apply specific progression to squat only (+5lb per session)
curl -X POST "http://localhost:8080/programs/g7h8i9j0-k1l2-3456-mnop-789012345678/progressions" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "u1v2w3x4-y5z6-7890-defg-123456789012",
    "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "priority": 2
  }'

# Apply +2.5lb progression to bench and overhead press
curl -X POST "http://localhost:8080/programs/g7h8i9j0-k1l2-3456-mnop-789012345678/progressions" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "v2w3x4y5-z6a7-8901-efgh-234567890123",
    "liftId": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
    "priority": 2
  }'
```

### PUT /programs/{programId}/progressions/{configId} - Update configuration (Admin)

```bash
curl -X PUT "http://localhost:8080/programs/g7h8i9j0-k1l2-3456-mnop-789012345678/progressions/t0u1v2w3-x4y5-6789-cdef-012345678901" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "priority": 3
  }'
```

### DELETE /programs/{programId}/progressions/{configId} - Delete configuration (Admin)

```bash
curl -X DELETE "http://localhost:8080/programs/g7h8i9j0-k1l2-3456-mnop-789012345678/progressions/t0u1v2w3-x4y5-6789-cdef-012345678901" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true"
```

---

## User Program Enrollment

Manage user enrollment in programs.

### GET /users/{userId}/program - Get current enrollment

```bash
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/program" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### POST /users/{userId}/program - Enroll user in program

```bash
# Enroll in Wendler 5/3/1 BBB
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/program" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "programId": "g7h8i9j0-k1l2-3456-mnop-789012345678"
  }'

# Enroll in Starting Strength
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/program" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "programId": "w3x4y5z6-a7b8-9012-fghi-345678901234"
  }'
```

### DELETE /users/{userId}/program - Unenroll user

```bash
curl -X DELETE "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/program" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

---

## State Advancement

Advance a user's program state (move to next day/week).

### POST /users/{userId}/program-state/advance - Advance state

```bash
# Advance to next day
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/program-state/advance" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "advanceType": "day"
  }'

# Advance to next week
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/program-state/advance" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "advanceType": "week"
  }'
```

---

## Workout Generation

Generate workouts based on user's program and state.

### GET /users/{userId}/workout - Generate current workout

```bash
# Generate workout for current state
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/workout" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# Generate workout for specific date
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/workout?date=2024-01-15" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# Generate workout with specific week override
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/workout?weekNumber=3" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# Generate workout with specific day slug override
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/workout?daySlug=squat-day" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### GET /users/{userId}/workout/preview - Preview a workout

```bash
# Preview workout for week 3, squat day
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/workout/preview?week=3&day=squat-day" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# Preview deload week, bench day
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/workout/preview?week=4&day=bench-day" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

---

## Progression History

Query a user's progression history.

### GET /users/{userId}/progression-history - List progression history

```bash
# Basic request
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/progression-history" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# Filter by lift
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/progression-history?lift_id=a1b2c3d4-e5f6-7890-abcd-ef1234567890" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# Filter by progression rule
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/progression-history?progression_id=s9t0u1v2-w3x4-5678-bcde-901234567890" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# With pagination
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/progression-history?page=1&pageSize=20" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

---

## Manual Progression Trigger

Manually trigger a progression for a user.

### POST /users/{userId}/progressions/trigger - Trigger progression

```bash
# Trigger progression for all configured lifts
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/progressions/trigger" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "s9t0u1v2-w3x4-5678-bcde-901234567890"
  }'

# Trigger progression for specific lift only
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/progressions/trigger" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "s9t0u1v2-w3x4-5678-bcde-901234567890",
    "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  }'

# Force trigger even if already applied this period
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/progressions/trigger" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "s9t0u1v2-w3x4-5678-bcde-901234567890",
    "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "force": true
  }'
```

---

## Workout Sessions

Manage workout session lifecycle.

### POST /workouts/start - Start a workout session

```bash
# Start a new workout session for the current user
curl -X POST "http://localhost:8080/workouts/start" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### GET /workouts/{id} - Get workout session by ID

```bash
curl -X GET "http://localhost:8080/workouts/a1b2c3d4-e5f6-7890-abcd-ef1234567890" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### POST /workouts/{id}/finish - Complete a workout session

```bash
curl -X POST "http://localhost:8080/workouts/a1b2c3d4-e5f6-7890-abcd-ef1234567890/finish" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### POST /workouts/{id}/abandon - Abandon a workout session

```bash
curl -X POST "http://localhost:8080/workouts/a1b2c3d4-e5f6-7890-abcd-ef1234567890/abandon" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### GET /users/{userId}/workouts - List user's workout history

```bash
# List all workout sessions
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/workouts" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# Filter by status
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/workouts?status=COMPLETED" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# With pagination
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/workouts?limit=10&offset=0" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### GET /users/{userId}/workouts/current - Get current in-progress workout

```bash
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/workouts/current" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

---

## Enrollment State Management

Manage enrollment state transitions.

### POST /users/{userId}/enrollment/next-cycle - Start next cycle

```bash
# Start the next cycle (only valid when enrollmentStatus is BETWEEN_CYCLES)
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/enrollment/next-cycle" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### POST /users/{userId}/enrollment/advance-week - Advance to next week

```bash
# Advance to the next week in the cycle
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/enrollment/advance-week" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

---

## Complete Workflow Examples

### New User Setup Workflow

```bash
# 1. Create lift maxes for all competition lifts
# Squat training max: 315 lbs
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/lift-maxes" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "type": "TRAINING_MAX",
    "value": 315.0
  }'

# Bench training max: 225 lbs
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/lift-maxes" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
    "type": "TRAINING_MAX",
    "value": 225.0
  }'

# Deadlift training max: 405 lbs
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/lift-maxes" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "d4e5f6a7-b8c9-0123-defg-456789012345",
    "type": "TRAINING_MAX",
    "value": 405.0
  }'

# Overhead Press training max: 135 lbs
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/lift-maxes" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "c3d4e5f6-a7b8-9012-cdef-345678901234",
    "type": "TRAINING_MAX",
    "value": 135.0
  }'

# 2. Enroll in a program
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/program" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "programId": "g7h8i9j0-k1l2-3456-mnop-789012345678"
  }'

# 3. Generate today's workout
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/workout" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### Daily Training Workflow

```bash
# 1. Get today's workout
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/workout" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# 2. After completing workout, advance to next day
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/program-state/advance" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "advanceType": "day"
  }'
```

### End of Cycle Workflow (5/3/1)

```bash
# 1. Preview next week's workout (should be week 1 of new cycle)
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/workout/preview?week=1&day=squat-day" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# 2. Manually trigger cycle-end progressions (if not automatic)
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/progressions/trigger" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "progressionId": "s9t0u1v2-w3x4-5678-bcde-901234567890"
  }'

# 3. Check progression history to verify
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/progression-history?page=1&pageSize=5" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### State Machine Workout Workflow

Complete workout session workflow using the state machine.

```bash
# 1. Check enrollment status
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/program" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# 2. Start a workout session
curl -X POST "http://localhost:8080/workouts/start" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
# Returns session ID, e.g., a1b2c3d4-e5f6-7890-abcd-ef1234567890

# 3. Generate the workout
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/workout" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# 4. Log sets during workout
curl -X POST "http://localhost:8080/sessions/a1b2c3d4-e5f6-7890-abcd-ef1234567890/sets" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{
    "sets": [
      {"prescriptionId": "prescription-uuid", "setNumber": 1, "weight": 265.0, "reps": 5},
      {"prescriptionId": "prescription-uuid", "setNumber": 2, "weight": 265.0, "reps": 5},
      {"prescriptionId": "prescription-uuid", "setNumber": 3, "weight": 265.0, "reps": 5}
    ]
  }'

# 5. Finish the workout
curl -X POST "http://localhost:8080/workouts/a1b2c3d4-e5f6-7890-abcd-ef1234567890/finish" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# 6. After completing all days in the week, advance to next week
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/enrollment/advance-week" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# 7. At end of cycle (when enrollmentStatus becomes BETWEEN_CYCLES), trigger progressions
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/progressions/trigger" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{"progressionId": "cycle-progression-uuid"}'

# 8. Start the next cycle
curl -X POST "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/enrollment/next-cycle" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### Handling Abandoned Workouts

```bash
# Check for current in-progress workout
curl -X GET "http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000/workouts/current" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# If there's an abandoned session blocking new workouts, abandon it
curl -X POST "http://localhost:8080/workouts/a1b2c3d4-e5f6-7890-abcd-ef1234567890/abandon" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"

# Now you can start a new workout
curl -X POST "http://localhost:8080/workouts/start" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000"
```

### Program Administrator Setup Workflow

```bash
# 1. Create base lifts (admin)
curl -X POST "http://localhost:8080/lifts" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"name": "Squat", "slug": "squat", "isCompetitionLift": true}'

# 2. Create a cycle
curl -X POST "http://localhost:8080/cycles" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"name": "4 Week Cycle", "lengthWeeks": 4}'

# 3. Create weekly lookups
curl -X POST "http://localhost:8080/weekly-lookups" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"name": "5/3/1 Percentages", "entries": {"1": 85, "2": 90, "3": 95, "4": 60}}'

# 4. Create prescriptions
curl -X POST "http://localhost:8080/prescriptions" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "liftId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "loadStrategy": {"type": "PERCENT_OF", "maxType": "TRAINING_MAX", "percentage": 85.0, "lookupKey": "week"},
    "setScheme": {"type": "FIXED", "sets": 1, "reps": 5, "isAmrap": true},
    "order": 1
  }'

# 5. Create days
curl -X POST "http://localhost:8080/days" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"name": "Squat Day", "slug": "squat-day"}'

# 6. Add prescriptions to days
curl -X POST "http://localhost:8080/days/h8i9j0k1-l2m3-4567-nopq-890123456789/prescriptions" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"prescriptionId": "e5f6a7b8-c9d0-1234-efgh-567890123456", "order": 1}'

# 7. Create weeks and add days
curl -X POST "http://localhost:8080/weeks" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"cycleId": "j0k1l2m3-n4o5-6789-pqrs-012345678901", "weekNumber": 1, "name": "Week 1"}'

curl -X POST "http://localhost:8080/weeks/i9j0k1l2-m3n4-5678-opqr-901234567890/days" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"dayId": "h8i9j0k1-l2m3-4567-nopq-890123456789", "position": 0}'

# 8. Create program
curl -X POST "http://localhost:8080/programs" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My 5/3/1 Program",
    "slug": "my-531-program",
    "description": "Custom 5/3/1 variant",
    "cycleId": "j0k1l2m3-n4o5-6789-pqrs-012345678901",
    "weeklyLookupId": "m3n4o5p6-q7r8-9012-stuv-345678901234",
    "defaultRounding": 5.0
  }'

# 9. Create progressions and attach to program
curl -X POST "http://localhost:8080/progressions" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"name": "Cycle +5lb", "type": "CYCLE", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX"}}'

curl -X POST "http://localhost:8080/programs/g7h8i9j0-k1l2-3456-mnop-789012345678/progressions" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Admin: true" \
  -H "Content-Type: application/json" \
  -d '{"progressionId": "s9t0u1v2-w3x4-5678-bcde-901234567890", "priority": 1}'
```

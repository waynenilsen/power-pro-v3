# PowerPro API Error Documentation

This document provides comprehensive documentation of all errors returned by the PowerPro API. Understanding these error formats and scenarios will help you implement robust error handling in your API clients.

## Table of Contents

- [Error Response Format](#error-response-format)
- [HTTP Status Codes](#http-status-codes)
- [Error Categories](#error-categories)
- [Common Error Scenarios](#common-error-scenarios)
  - [Authentication Errors](#authentication-errors)
  - [Authorization Errors](#authorization-errors)
  - [Validation Errors](#validation-errors)
  - [Resource Errors](#resource-errors)
  - [Business Rule Violations](#business-rule-violations)
  - [Server Errors](#server-errors)
- [Error Handling Best Practices](#error-handling-best-practices)

---

## Error Response Format

All API errors are returned as JSON with a consistent, structured envelope:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {}
  }
}
```

### Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `error.code` | string | Yes | Machine-readable error code (e.g., `NOT_FOUND`, `VALIDATION_ERROR`) |
| `error.message` | string | Yes | Human-readable error message describing what went wrong |
| `error.details` | object | No | Optional structured details, typically containing validation errors |

### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `NOT_FOUND` | 404 | Resource does not exist |
| `VALIDATION_ERROR` | 400 | Input validation failed |
| `BAD_REQUEST` | 400 | Malformed request |
| `CONFLICT` | 409 | Resource conflict (e.g., duplicate slug) |
| `FORBIDDEN` | 403 | Permission denied |
| `UNAUTHORIZED` | 401 | Authentication required |
| `UNPROCESSABLE_ENTITY` | 422 | Valid request but cannot be processed |
| `INTERNAL_ERROR` | 500 | Server error |

### Example Responses

**Resource not found:**
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "lift not found: abc123"
  }
}
```

**Validation failure with details:**
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "validation failed",
    "details": {
      "validationErrors": [
        "name is required",
        "value must be positive",
        "slug must be unique"
      ]
    }
  }
}
```

**Conflict error:**
```json
{
  "error": {
    "code": "CONFLICT",
    "message": "a lift with this slug already exists"
  }
}
```

---

## HTTP Status Codes

The API uses standard HTTP status codes to indicate the result of each request:

### Success Codes

| Code | Name | Description | Used For |
|------|------|-------------|----------|
| `200` | OK | Request succeeded | GET requests, PUT updates |
| `201` | Created | Resource created successfully | POST requests creating new resources |
| `204` | No Content | Request succeeded with no response body | DELETE requests |

### Client Error Codes

| Code | Name | Description | Used For |
|------|------|-------------|----------|
| `400` | Bad Request | Invalid request format or validation failure | Malformed JSON, missing required fields, invalid field values |
| `401` | Unauthorized | Authentication required | Missing or invalid authentication headers |
| `403` | Forbidden | Authenticated but not authorized | Accessing another user's resources, admin-only endpoints |
| `404` | Not Found | Resource does not exist | Invalid ID or slug in URL |
| `409` | Conflict | Request conflicts with current state | Duplicate slug, foreign key constraints, state conflicts |
| `422` | Unprocessable Entity | Valid request but cannot be processed | Business rule violations (missing lift max, etc.) |

### Server Error Codes

| Code | Name | Description | Used For |
|------|------|-------------|----------|
| `500` | Internal Server Error | Unexpected server error | Database errors, unexpected conditions |

---

## Error Categories

The API uses seven internal error categories that map to HTTP status codes:

| Category | HTTP Status | Description |
|----------|-------------|-------------|
| `not found` | 404 | Resource does not exist |
| `validation failed` | 400 | Input validation failed |
| `bad request` | 400 | Malformed request |
| `unauthorized` | 401 | Authentication required |
| `forbidden` | 403 | Permission denied |
| `conflict` | 409 | State/data conflict |
| `internal error` | 500 | Server-side error |

---

## Common Error Scenarios

### Authentication Errors

#### Missing Authentication Headers

**When**: No `Authorization` or `X-User-ID` header is provided for a protected endpoint.

**HTTP Status**: `401 Unauthorized`

```http
POST /lifts HTTP/1.1
Content-Type: application/json

{"name": "Squat", "slug": "squat"}
```

```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Authentication required"
  }
}
```

**Resolution**: Include either:
- `Authorization: Bearer {session-token}` header (obtained from login), or
- `X-User-ID: {userId}` header (for development/testing only)

#### Invalid Session Token

**When**: The provided session token is expired, revoked, or malformed.

**HTTP Status**: `401 Unauthorized`

```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "invalid or expired session"
  }
}
```

**Resolution**: Obtain a new session token by logging in again via `POST /auth/login`.

#### Invalid Login Credentials

**When**: Email or password is incorrect during login.

**HTTP Status**: `401 Unauthorized`

```http
POST /auth/login HTTP/1.1
Content-Type: application/json

{"email": "user@example.com", "password": "wrongpassword"}
```

```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "invalid email or password"
  }
}
```

**Resolution**: Verify the email and password are correct.

#### Email Already Registered

**When**: Attempting to register with an email that already exists.

**HTTP Status**: `409 Conflict`

```http
POST /auth/register HTTP/1.1
Content-Type: application/json

{"email": "existing@example.com", "password": "password123"}
```

```json
{
  "error": {
    "code": "CONFLICT",
    "message": "email already registered"
  }
}
```

**Resolution**: Use a different email address or login with the existing account.

#### Missing Required Registration Fields

**When**: Email or password is missing during registration.

**HTTP Status**: `400 Bad Request`

```http
POST /auth/register HTTP/1.1
Content-Type: application/json

{"email": "user@example.com"}
```

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "validation failed",
    "details": {
      "validationErrors": ["password is required"]
    }
  }
}
```

**Resolution**: Provide both email and password fields.

---

### Authorization Errors

#### Admin Privileges Required

**When**: A non-admin user attempts to access an admin-only endpoint.

**HTTP Status**: `403 Forbidden`

```http
POST /lifts HTTP/1.1
Authorization: Bearer session-token
Content-Type: application/json

{"name": "Squat", "slug": "squat"}
```

```json
{
  "error": {
    "code": "FORBIDDEN",
    "message": "Admin privileges required"
  }
}
```

**Resolution**: Include the `X-Admin: true` header for admin operations:
```http
POST /lifts HTTP/1.1
Authorization: Bearer session-token
X-Admin: true
Content-Type: application/json
```

#### Accessing Another User's Resources

**When**: A user attempts to access or modify resources belonging to another user.

**HTTP Status**: `403 Forbidden`

```json
{
  "error": {
    "code": "FORBIDDEN",
    "message": "you can only view your own enrollment"
  }
}
```

**Resolution**: Ensure you're accessing resources that belong to the authenticated user, or use admin privileges.

#### Accessing Another User's Profile (Non-Admin)

**When**: A user attempts to view another user's profile without admin privileges.

**HTTP Status**: `403 Forbidden`

```http
GET /users/another-user-uuid/profile HTTP/1.1
Authorization: Bearer your-session-token
```

```json
{
  "error": {
    "code": "FORBIDDEN",
    "message": "access denied: you do not have permission to view this profile"
  }
}
```

**Resolution**: Only access your own profile, or use admin privileges to view others.

#### Updating Another User's Profile

**When**: Any user (including admins) attempts to update another user's profile.

**HTTP Status**: `403 Forbidden`

```http
PUT /users/another-user-uuid/profile HTTP/1.1
Authorization: Bearer admin-session-token
X-Admin: true
Content-Type: application/json

{"name": "New Name"}
```

```json
{
  "error": {
    "code": "FORBIDDEN",
    "message": "access denied: only the profile owner can update their profile"
  }
}
```

**Resolution**: Profile updates are strictly owner-only; even admins cannot modify another user's profile.

#### Accessing Another User's Dashboard

**When**: Any user (including admins) attempts to access another user's dashboard.

**HTTP Status**: `403 Forbidden`

```http
GET /users/another-user-uuid/dashboard HTTP/1.1
Authorization: Bearer admin-session-token
X-Admin: true
```

```json
{
  "error": {
    "code": "FORBIDDEN",
    "message": "access denied: you can only view your own dashboard"
  }
}
```

**Resolution**: Dashboard access is strictly owner-only; only the authenticated user can view their own dashboard.

---

### Validation Errors

#### Missing Required Fields

**When**: Required fields are not included in the request body.

**HTTP Status**: `400 Bad Request`

```http
POST /lifts HTTP/1.1
Authorization: Bearer 123e4567-e89b-12d3-a456-426614174000
X-Admin: true
Content-Type: application/json

{"name": ""}
```

```json
{
  "error": "validation failed",
  "details": [
    "name is required",
    "slug is required"
  ]
}
```

**Resolution**: Provide all required fields with valid values.

#### Invalid Field Format

**When**: A field value doesn't match the expected format.

**HTTP Status**: `400 Bad Request`

```json
{
  "error": "validation failed",
  "details": [
    "value must be positive",
    "slug must contain only lowercase letters, numbers, and hyphens"
  ]
}
```

**Resolution**: Review field constraints in the API documentation and ensure values match the expected format.

#### Invalid Request Body

**When**: The request body is not valid JSON.

**HTTP Status**: `400 Bad Request`

```http
POST /lifts HTTP/1.1
Content-Type: application/json

{invalid json here}
```

```json
{
  "error": "invalid request body"
}
```

**Resolution**: Ensure the request body is valid JSON with correct syntax.

#### Missing Required Query Parameter

**When**: A required query parameter is not provided.

**HTTP Status**: `400 Bad Request`

```http
GET /lookups/daily?cycle=123 HTTP/1.1
Authorization: Bearer 123e4567-e89b-12d3-a456-426614174000
```

```json
{
  "error": "week: missing required parameter"
}
```

**Resolution**: Include all required query parameters.

#### Invalid Path Parameter

**When**: A path parameter is missing or invalid (e.g., invalid UUID).

**HTTP Status**: `400 Bad Request`

```http
GET /lifts/not-a-uuid HTTP/1.1
Authorization: Bearer 123e4567-e89b-12d3-a456-426614174000
```

```json
{
  "error": "missing lift ID"
}
```

**Resolution**: Ensure path parameters are valid UUIDs where required.

---

### Resource Errors

#### Resource Not Found

**When**: The requested resource does not exist.

**HTTP Status**: `404 Not Found`

```http
GET /lifts/00000000-0000-0000-0000-000000000000 HTTP/1.1
Authorization: Bearer 123e4567-e89b-12d3-a456-426614174000
```

```json
{
  "error": "lift not found: 00000000-0000-0000-0000-000000000000"
}
```

**Resolution**: Verify the resource ID or slug is correct. List resources first if unsure.

#### Resource Not Found by Slug

**When**: Looking up a resource by slug that doesn't exist.

**HTTP Status**: `404 Not Found`

```http
GET /lifts/by-slug/nonexistent-lift HTTP/1.1
Authorization: Bearer 123e4567-e89b-12d3-a456-426614174000
```

```json
{
  "error": "lift not found: nonexistent-lift"
}
```

#### Parent Resource Not Found

**When**: Creating/updating a resource with a reference to a non-existent parent.

**HTTP Status**: `400 Bad Request`

```json
{
  "error": "parentLiftId: parent lift not found"
}
```

**Resolution**: Verify the referenced parent resource exists before creating child resources.

---

### Conflict Errors

#### Duplicate Slug

**When**: Creating or updating a resource with a slug that already exists.

**HTTP Status**: `409 Conflict`

```http
POST /lifts HTTP/1.1
Authorization: Bearer 123e4567-e89b-12d3-a456-426614174000
X-Admin: true
Content-Type: application/json

{"name": "Squat", "slug": "squat"}
```

```json
{
  "error": "slug already exists"
}
```

**Resolution**: Choose a unique slug or update the existing resource instead.

#### Foreign Key Constraint Violation

**When**: Attempting to delete a resource that is referenced by other records.

**HTTP Status**: `409 Conflict`

```http
DELETE /lifts/123e4567-e89b-12d3-a456-426614174000 HTTP/1.1
Authorization: Bearer 123e4567-e89b-12d3-a456-426614174000
X-Admin: true
```

```json
{
  "error": "cannot delete lift: it is referenced by other lifts as a parent"
}
```

**Resolution**: Delete or update dependent resources first, or modify them to remove the reference.

#### Duplicate Resource

**When**: Attempting to create a resource that would duplicate an existing one (beyond just slug).

**HTTP Status**: `409 Conflict`

```json
{
  "error": "user is already enrolled in this program"
}
```

**Resolution**: Check for existing resources before creating, or update the existing resource instead.

---

### Business Rule Violations

#### Missing Lift Max for Workout Generation

**When**: Attempting to generate a workout without the required lift maxes set up.

**HTTP Status**: `422 Unprocessable Entity`

```http
GET /workout?week=1&day=1 HTTP/1.1
Authorization: Bearer 123e4567-e89b-12d3-a456-426614174000
```

```json
{
  "error": "missing lift max: set up your training maxes to generate workouts",
  "details": [
    "missing lift max for squat (TRAINING_MAX)",
    "missing lift max for bench-press (TRAINING_MAX)"
  ]
}
```

**Resolution**: Create the required lift max records before generating workouts:
```http
POST /lift-maxes HTTP/1.1
Authorization: Bearer 123e4567-e89b-12d3-a456-426614174000
Content-Type: application/json

{
  "liftId": "...",
  "maxType": "TRAINING_MAX",
  "value": 315,
  "unit": "lb"
}
```

#### Invalid State Transition

**When**: Attempting to advance state in an invalid way.

**HTTP Status**: `400 Bad Request` or `422 Unprocessable Entity`

```json
{
  "error": "cannot advance: already at end of cycle"
}
```

**Resolution**: Check current state before attempting transitions.

---

### Server Errors

#### Internal Server Error

**When**: An unexpected error occurs on the server.

**HTTP Status**: `500 Internal Server Error`

```json
{
  "error": "internal server error"
}
```

**Note**: Internal error details are logged server-side but not exposed to clients for security reasons.

**Resolution**: If this error persists:
1. Check that your request is correctly formatted
2. Try the request again (may be a transient issue)
3. Contact API support if the problem continues

---

## Error Handling Best Practices

### 1. Check HTTP Status Code First

Use the HTTP status code to determine the general category of error:

```go
switch resp.StatusCode {
case 400:
    // Validation error - check details
case 401:
    // Authentication needed - redirect to login
case 403:
    // Permission denied - show access error
case 404:
    // Resource not found
case 409:
    // Conflict - handle duplicate or constraint
case 422:
    // Business rule - guide user to fix
case 500:
    // Server error - retry or show generic error
}
```

### 2. Parse the Error Response

Always attempt to parse the error body for detailed information:

```go
type ErrorResponse struct {
    Error   string   `json:"error"`
    Details []string `json:"details,omitempty"`
}

var errResp ErrorResponse
if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil {
    // Use errResp.Error and errResp.Details
}
```

### 3. Handle Validation Details

When `details` is present, display all issues to the user:

```go
if len(errResp.Details) > 0 {
    for _, detail := range errResp.Details {
        displayFieldError(detail)
    }
} else {
    displayError(errResp.Error)
}
```

### 4. Implement Retry Logic for Server Errors

For `500` errors, implement exponential backoff:

```go
const maxRetries = 3
for attempt := 0; attempt < maxRetries; attempt++ {
    resp, err := makeRequest()
    if err != nil || resp.StatusCode >= 500 {
        time.Sleep(time.Duration(1<<attempt) * time.Second)
        continue
    }
    break
}
```

### 5. Never Trust Client-Side Validation Alone

Always handle server validation errors, even if you validate client-side:

- Server may have additional constraints
- Multiple clients may have different validation logic
- Business rules may change without client updates

### 6. Log Errors for Debugging

Log full error responses to help diagnose issues:

```go
if resp.StatusCode >= 400 {
    log.Printf("API error: status=%d error=%s details=%v",
        resp.StatusCode, errResp.Error, errResp.Details)
}
```

---

## See Also

- [API Reference](./api-reference.md) - Complete endpoint documentation
- [Example Requests](./example-requests.md) - Copy-paste ready examples
- [Example Responses](./example-responses.md) - Response format examples

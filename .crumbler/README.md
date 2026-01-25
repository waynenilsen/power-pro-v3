# MVP Frontend-Readiness

## Overview

The core workout engine is complete: program configuration, workout generation, set logging, progression tracking, and state machine management all work end-to-end (validated by Texas Method E2E test). However, several user-facing features are missing before a frontend team can build a complete app.

This plan addresses the gaps needed to "throw it over the fence" to frontend development.

---

## What's Done (Not In Scope)

- 78 API endpoints for program config, workouts, progressions
- State machine for enrollment/cycle/week/workout lifecycle
- Workout generation with daily/weekly lookups
- Progression system (linear, AMRAP-based) with history
- Consistent envelope response format with error codes
- E2E tests documenting the full workout flow

---

## Scope: 5 Work Items

### 1. Authentication System

**Current state:** Fake `X-User-ID` and `X-Admin` headers for testing.

**Target state:** Real session-based authentication.

**Technical requirements:**
- **Session auth, NOT JWT** - Server-side session storage in SQLite
- **Headless API** - Auth token continues to come in header (e.g., `Authorization: Bearer <session_token>`)
- **Email/password signup** - No OAuth, no magic links
- **Name is optional** - Only email + password required for registration
- **No email verification for now** - We don't have email infrastructure, defer this
- **Password hashing** - bcrypt or argon2

**Endpoints needed:**
```
POST /auth/register    - Create account (email, password, optional name)
POST /auth/login       - Get session token
POST /auth/logout      - Invalidate session
GET  /auth/me          - Get current user from session
```

**Schema changes:**
- Add `email`, `password_hash`, `name` to `users` table
- Create `sessions` table (id, user_id, token, expires_at, created_at)

---

### 2. User Profile

**Current state:** Users exist but have no editable profile data.

**Target state:** Users can view/update their profile and preferences.

**Endpoints needed:**
```
GET  /users/{id}/profile     - Get profile (name, email, preferences)
PUT  /users/{id}/profile     - Update profile
```

**Profile fields:**
- `name` (optional, display name)
- `weight_unit` (lb/kg, default: lb)
- `created_at`, `updated_at`

**Note:** Email changes should require re-authentication (defer for now, make email immutable initially).

---

### 3. Dashboard Endpoint

**Current state:** Frontend must make 4+ calls to understand "what's happening now."

**Target state:** Single aggregated endpoint for the home screen.

**Endpoint:**
```
GET /users/{id}/dashboard
```

**Response shape:**
```json
{
  "enrollment": {
    "status": "ACTIVE",
    "program_name": "Texas Method",
    "cycle_iteration": 2,
    "cycle_status": "IN_PROGRESS",
    "week_number": 1,
    "week_status": "IN_PROGRESS"
  },
  "next_workout": {
    "day_name": "Volume Day",
    "day_slug": "volume",
    "exercise_count": 2,
    "estimated_sets": 10
  },
  "current_session": null,  // or {id, started_at} if IN_PROGRESS
  "recent_workouts": [
    {"date": "2024-01-15", "day_name": "Intensity Day", "sets_completed": 4}
  ],
  "current_maxes": [
    {"lift": "Squat", "value": 320, "type": "TRAINING_MAX"},
    {"lift": "Bench", "value": 227.5, "type": "TRAINING_MAX"}
  ]
}
```

This is an aggregation endpoint - no new data, just combines existing queries.

---

### 4. Seed Canonical Programs

**Current state:** No programs exist. E2E tests create ad-hoc programs with unique slugs, then discard them.

**Target state:** 3-5 popular programs seeded and ready for users to enroll.

**Programs to seed:**
- Starting Strength (beginner, 3 days/week, linear progression)
- Texas Method (intermediate, 3 days/week, weekly progression)
- 5/3/1 (intermediate, 4 days/week, monthly progression)
- GZCLP (beginner/intermediate, 3-4 days/week, tiered progression)

**Implementation:**
- Add migration that creates full program structure (lookups, prescriptions, days, weeks, cycles, programs, progressions)
- Use canonical slugs: `starting-strength`, `texas-method`, `531`, `gzclp`
- Tests continue creating their own programs with unique slugs (`texas-method-{testID}`) - no conflict

**Note:** This is a large migration but mostly mechanical - translate E2E test setup code into SQL INSERTs.

---

### 5. Program Discovery

**Current state:** `GET /programs` returns flat list, no filtering.

**Target state:** Frontend can help users find appropriate programs.

**Enhancements to existing endpoint:**
```
GET /programs?difficulty=beginner&days_per_week=3&focus=strength
```

**Filter parameters:**
- `difficulty` (beginner, intermediate, advanced)
- `days_per_week` (1-7)
- `focus` (strength, hypertrophy, peaking)
- `has_amrap` (true/false)
- `search` (name substring)

**Schema changes:**
- Add metadata columns to `programs` table: `difficulty`, `days_per_week`, `focus`, `has_amrap`
- Backfill existing programs with appropriate values

**Program detail enhancement:**
```
GET /programs/{id}
```
Should include:
- Sample week preview (day names, exercise counts)
- Lift requirements (what lifts the program uses)
- Estimated session duration

---

## Out of Scope (Deferred)

- **Analytics/Insights** - Progress charts, volume tracking, PR history visualization
- **Email verification** - No email infrastructure yet
- **Password reset** - Requires email
- **OAuth/social login** - Keep it simple
- **Real-time features** - No WebSockets needed for MVP

---

## Implementation Order

1. **Auth** - Everything else depends on having real users
2. **User Profile** - Simple, needed for preferences
3. **Dashboard** - High frontend value, aggregation only
4. **Seed Programs** - Users need programs to enroll in
5. **Program Discovery** - Enhances existing endpoint (useful once programs exist)

---

## Success Criteria

- [ ] User can register with email/password
- [ ] User can login and receive session token
- [ ] Session token works in Authorization header for all protected endpoints
- [ ] User can view/update their profile
- [ ] Dashboard endpoint returns aggregated current state
- [ ] At least 3 canonical programs seeded and enrollable
- [ ] Programs can be filtered by difficulty, days per week, focus
- [ ] E2E tests cover auth flow and new endpoints
- [ ] Existing E2E tests continue to pass (backwards compatible)

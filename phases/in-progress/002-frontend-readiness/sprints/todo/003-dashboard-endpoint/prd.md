# PRD 003: Dashboard Endpoint

## Product Vision

PowerPro users need a single aggregated endpoint that provides all information necessary to render a home screen. Rather than making multiple API calls, the dashboard endpoint combines enrollment status, next workout, current session, recent workouts, and current maxes into one efficient response.

## Strategic Objectives

1. **Reduce API Round-trips**: Aggregate multiple queries into a single endpoint
2. **Enable Home Screen**: Provide all data needed for a mobile/web home screen
3. **Surface Workout Context**: Show users where they are in their program
4. **Prepare for Offline**: Single response can be cached for offline viewing

## Themes & Initiatives

### Theme 1: Enrollment Status Aggregation
- **Strategic Objective**: Surface Workout Context
- **Rationale**: Users need to know their current program, cycle, and week status at a glance
- **Initiatives**:
  - Initiative A: Return active enrollment status (ACTIVE, PAUSED, NONE)
  - Initiative B: Include program name and cycle/week position
  - Initiative C: Show cycle and week progress status (IN_PROGRESS, COMPLETED)

### Theme 2: Next Workout Preview
- **Strategic Objective**: Enable Home Screen
- **Rationale**: Users want to see what workout comes next without navigating
- **Initiatives**:
  - Initiative A: Calculate the next scheduled workout day
  - Initiative B: Return day name, slug, and exercise count
  - Initiative C: Estimate total sets for workout planning

### Theme 3: Current Session Tracking
- **Strategic Objective**: Surface Workout Context
- **Rationale**: If a user is mid-workout, they should see that session immediately
- **Initiatives**:
  - Initiative A: Detect if user has an active workout session
  - Initiative B: Return current session details if present
  - Initiative C: Return null if no active session

### Theme 4: Recent Workout History
- **Strategic Objective**: Enable Home Screen, Surface Workout Context
- **Rationale**: Users want to see their recent activity and progress
- **Initiatives**:
  - Initiative A: Query last N completed workouts
  - Initiative B: Return date, day name, and sets completed
  - Initiative C: Limit to manageable count (e.g., last 5)

### Theme 5: Current Training Maxes
- **Strategic Objective**: Enable Home Screen
- **Rationale**: Powerlifters constantly reference their current maxes
- **Initiatives**:
  - Initiative A: Retrieve current training maxes for all lifts
  - Initiative B: Return lift name, value, and max type
  - Initiative C: Handle case where maxes are not yet set

## Success Metrics

| Metric | Target |
|--------|--------|
| Dashboard returns all sections | Complete |
| Single API call replaces 5+ separate calls | Complete |
| Response time < 200ms | Complete |
| Handles no enrollment gracefully | Complete |
| Handles mid-workout state correctly | Complete |
| Returns empty arrays for missing data (not errors) | Complete |

## Timeline

| Phase | Scope |
|-------|-------|
| Now | Dashboard service aggregation logic |
| Now | Enrollment status aggregation |
| Now | Next workout calculation |
| Now | Recent workouts query |
| Now | Current maxes query |
| Now | Dashboard API endpoint |
| Now | E2E tests |

## Dependencies

- Sprint 001: Authentication System (session-based auth)
- Sprint 002: User Profile (user preferences)
- Existing domain models: Enrollment, Cycle, Week, WorkoutSession, TrainingMax

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Performance if queries not optimized | Medium | Medium | Use efficient joins, consider denormalization later |
| Complex null handling across sections | Low | Low | Define clear contract for empty/null states |
| Response size too large | Low | Low | Limit recent workouts, omit unnecessary fields |

## Out of Scope

- Real-time updates (polling or websockets)
- Caching layer (deferred to tech debt sprint)
- Analytics/statistics aggregation (separate endpoint)
- Program recommendations
- Social features (leaderboards, friends)

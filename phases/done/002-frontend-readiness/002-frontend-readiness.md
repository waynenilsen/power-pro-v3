# Phase 002: Frontend Readiness

## Vision & Strategic Objectives

### Product Vision
Make the PowerPro API frontend-ready by implementing real authentication, user profiles, and aggregated endpoints. The core workout engine is complete—this phase bridges the gap between a functional API and one that frontend developers can build a complete app against.

### Strategic Objectives
1. **Enable Frontend Development**: Provide all endpoints and data structures a frontend team needs to build a complete powerlifting app
2. **Real User Management**: Replace fake authentication headers with real session-based auth and user profiles
3. **Improved Developer Experience**: Aggregated endpoints reduce frontend complexity and API calls
4. **Production Readiness**: Seed canonical programs so users have real content to enroll in

### Context
Phase 001 delivered a complete workout engine with 78 API endpoints covering program configuration, workout generation, set logging, progression tracking, and state machine management. However, the current implementation uses fake `X-User-ID` headers for authentication and requires frontends to make multiple API calls to understand the current state. This phase addresses these gaps to enable frontend development.

## Themes & Initiatives

### Theme 1: Authentication System
- **Strategic Objective**: Real User Management, Enable Frontend Development
- **Rationale**: Every protected endpoint requires authentication. Without real auth, no production frontend can be built.
- **Initiatives**:
  - Initiative A: Implement session-based authentication with SQLite storage
  - Initiative B: Email/password registration (no OAuth, no magic links)
  - Initiative C: Session token in Authorization header for API access
  - Initiative D: Password hashing with bcrypt or argon2
- **Dependencies**: None—this is the foundation of Phase 002
- **Risks**: Session management complexity; ensure proper security practices

### Theme 2: User Management
- **Strategic Objective**: Real User Management, Improved Developer Experience
- **Rationale**: Users need profiles with preferences (e.g., weight units) for a personalized experience.
- **Initiatives**:
  - Initiative A: User profile endpoints (view/update name, preferences)
  - Initiative B: Weight unit preference (lb/kg)
  - Initiative C: Extend users table with profile fields
- **Dependencies**: Theme 1 (Authentication)
- **Risks**: Scope creep—keep profile minimal for MVP

### Theme 3: API Aggregation
- **Strategic Objective**: Improved Developer Experience, Enable Frontend Development
- **Rationale**: Frontends shouldn't need 4+ API calls to render a home screen.
- **Initiatives**:
  - Initiative A: Dashboard endpoint aggregating enrollment status, next workout, recent workouts, current maxes
  - Initiative B: Enhanced program detail endpoint with sample week preview
- **Dependencies**: Themes 1-2, existing workout/enrollment endpoints
- **Risks**: Aggregation logic may be complex; ensure performance is acceptable

### Theme 4: Program Content
- **Strategic Objective**: Enable Frontend Development, Production Readiness
- **Rationale**: Users need real programs to enroll in. Without seeded programs, the app has no content.
- **Initiatives**:
  - Initiative A: Seed 3-5 canonical programs (Starting Strength, Texas Method, 5/3/1, GZCLP)
  - Initiative B: Add program metadata (difficulty, days_per_week, focus, has_amrap)
  - Initiative C: Program discovery with filtering
- **Dependencies**: Existing program infrastructure from Phase 001
- **Risks**: Large migration for seeding programs; ensure data integrity

## Timeline

| Phase | Timeline | Themes/Initiatives |
|-------|----------|-------------------|
| Now | Current | Theme 1: Authentication System (sessions, registration, login) |
| Now | Current | Theme 2: User Management (profiles, preferences) |
| Next | Following | Theme 3: API Aggregation (dashboard, enhanced program detail) |
| Next | Following | Theme 4: Program Content (seed programs, discovery filters) |

## Success Metrics

| Metric | Target |
|--------|--------|
| User can register with email/password | Complete |
| User can login and receive session token | Complete |
| Session token works in Authorization header for all protected endpoints | Complete |
| User can view/update their profile | Complete |
| Dashboard endpoint returns aggregated current state | Complete |
| At least 3 canonical programs seeded and enrollable | Complete |
| Programs can be filtered by difficulty, days per week, focus | Complete |
| E2E tests cover auth flow and new endpoints | Complete |
| Existing E2E tests continue to pass (backwards compatible) | Complete |

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Session security vulnerabilities | Medium | High | Use established patterns; secure token generation; proper expiration |
| Breaking existing E2E tests | Low | High | Maintain backwards compatibility for test auth headers |
| Program seeding data errors | Medium | Medium | Validate against existing E2E test program configs |
| Dashboard performance issues | Low | Medium | Profile queries; add indices if needed |

## Review & Update Process

- **Review cadence**: After each sprint completion
- **Owner**: Engineering Lead
- **Approval**: Product Owner reviews sprint completion, Engineering Lead approves technical implementation
- **Update triggers**:
  - Discovery of additional frontend requirements
  - Security considerations requiring architecture changes
  - Completion of sprints enabling next phase readiness assessment

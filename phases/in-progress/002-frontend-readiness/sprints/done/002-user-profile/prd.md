# PRD 002: User Profile

## Product Vision

PowerPro users need to view and manage their profile data and preferences. This PRD establishes the user profile functionality that enables personalization and prepares for future features like weight-based program calculations.

## Strategic Objectives

1. **Enable Profile Management**: Users can view and update their profile information
2. **Support User Preferences**: Store user preferences like weight unit (lb/kg)
3. **Prepare for Program Features**: Weight unit preference will affect how weights are displayed in future program features
4. **Foundation for Personalization**: Establish the pattern for user-specific settings

## Themes & Initiatives

### Theme 1: Profile Viewing
- **Strategic Objective**: Enable Profile Management
- **Rationale**: Users need to see their account information and verify their data
- **Initiatives**:
  - Initiative A: GET endpoint to retrieve user profile
  - Initiative B: Return profile fields (name, email, preferences)
  - Initiative C: Include timestamps (created_at, updated_at)

### Theme 2: Profile Updates
- **Strategic Objective**: Enable Profile Management, Support User Preferences
- **Rationale**: Users should be able to update their display name and preferences
- **Initiatives**:
  - Initiative A: PUT endpoint to update profile
  - Initiative B: Support partial updates (update only provided fields)
  - Initiative C: Validate input constraints (name length, valid weight unit)

### Theme 3: Weight Unit Preference
- **Strategic Objective**: Support User Preferences, Prepare for Program Features
- **Rationale**: Powerlifters use different weight systems; users should choose their preferred unit
- **Initiatives**:
  - Initiative A: Store weight_unit preference (lb/kg)
  - Initiative B: Default to lb for new users
  - Initiative C: Validate weight_unit values

### Theme 4: Authorization
- **Strategic Objective**: Enable Profile Management
- **Rationale**: Users should only access their own profile; admins can access any profile
- **Initiatives**:
  - Initiative A: Users can only GET/PUT their own profile
  - Initiative B: Admin users can GET any user's profile
  - Initiative C: Return 403 for unauthorized access attempts

## Success Metrics

| Metric | Target |
|--------|--------|
| User can view their profile | Complete |
| User can update their name | Complete |
| User can update weight_unit preference | Complete |
| Profile endpoint returns correct fields | Complete |
| Non-owner cannot access another user's profile | Complete |
| Admin can access any profile | Complete |
| Invalid weight_unit rejected with 400 | Complete |

## Timeline

| Phase | Scope |
|-------|-------|
| Now | Schema migration (weight_unit column) |
| Now | Profile domain logic |
| Now | Profile API endpoints |
| Now | Authorization rules |
| Now | E2E tests |

## Dependencies

- Sprint 001: Authentication System (must be complete for session-based auth)
- Users table with email, name fields (added in Sprint 001)

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Breaking existing tests | Low | Medium | Profile endpoints are new; existing endpoints unchanged |
| Weight unit validation edge cases | Low | Low | Use strict enum validation (only "lb" or "kg") |

## Out of Scope

- Email changes (require re-authentication - deferred)
- Password changes (separate change-password flow - deferred)
- Avatar/profile picture (no file storage infrastructure)
- Profile visibility settings (no social features planned)
- User deletion/account deactivation (deferred)

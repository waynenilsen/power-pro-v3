# Sprint 002 - User Profile

## Overview

Allow users to view and update their profile data and preferences.

## Endpoints Needed

```
GET  /users/{id}/profile     - Get profile (name, email, preferences)
PUT  /users/{id}/profile     - Update profile
```

## Profile Fields

- `name` (optional, display name)
- `weight_unit` (lb/kg, default: lb)
- `created_at`, `updated_at`

## Design Decisions

- Email is immutable initially (changes require re-authentication - deferred)
- Weight unit preference affects how weights are displayed/returned

## Tasks

1. Create sprint directory: `phases/in-progress/002-frontend-readiness/sprints/todo/002-user-profile/`
2. Write `prd.md` - Product requirements for user profile
3. Write `erd.md` - Engineering requirements with detailed specs
4. Create ticket directory structure
5. Create tickets (at least 5):
   - Schema migration for profile fields (weight_unit column)
   - Profile domain logic (validation, preferences)
   - Profile API endpoints (GET, PUT)
   - Authorization rules (users can only access own profile, admins can access all)
   - E2E tests for profile endpoints

## Acceptance Criteria

- [ ] Sprint directory structure exists
- [ ] PRD document completed
- [ ] ERD document completed with REQ-XXX requirements
- [ ] At least 5 tickets created in `tickets/todo/`
- [ ] Tickets reference ERD requirements

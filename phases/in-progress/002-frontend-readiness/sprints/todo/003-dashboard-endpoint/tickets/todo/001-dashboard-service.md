# 001: Dashboard Service

## ERD Reference
Implements: REQ-DASH-001, REQ-DASH-002, NFR-006, NFR-007, NFR-008

## Description
Create the dashboard service that orchestrates aggregation of all dashboard sections. This service acts as the coordination layer, calling individual aggregation functions and combining results into a single response.

## Context / Background
The dashboard endpoint requires data from multiple sources: enrollment, workouts, sessions, and maxes. Rather than coupling these into a single monolithic function, we create a service that delegates to specialized aggregation functions. This follows the existing service pattern in PowerPro.

## Acceptance Criteria
- [ ] Create `internal/dashboard/service.go` with DashboardService struct
- [ ] DashboardService constructor accepts dependencies (enrollment service, workout service, etc.)
- [ ] Implement `GetDashboard(ctx context.Context, userID uuid.UUID) (*Dashboard, error)` method
- [ ] Dashboard struct contains all five sections: Enrollment, NextWorkout, CurrentSession, RecentWorkouts, CurrentMaxes
- [ ] Each section can be null (pointer) or empty slice as appropriate
- [ ] Service calls individual aggregation functions in parallel where possible
- [ ] Error in one section does not fail entire dashboard (log and return null for that section)
- [ ] Unit tests for service orchestration logic

## Technical Notes
- Follow existing service patterns in `internal/` directory
- Use dependency injection for all sub-services
- Consider using errgroup for parallel section fetching
- Dashboard struct should use pointers for nullable sections
- Use existing domain types where they exist

## Dependencies
- Blocks: 007 (API endpoint needs service)
- Blocked by: 002, 003, 004, 005, 006 (aggregation functions)

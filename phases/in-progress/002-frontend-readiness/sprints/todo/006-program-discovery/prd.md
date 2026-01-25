# PRD 006: Program Discovery

## Product Vision

PowerPro needs robust program discovery capabilities to help users find the right training program for their needs. With canonical programs seeded in Sprint 004, users need ways to filter, search, and evaluate programs based on their experience level, schedule, and training goals. Program discovery transforms a simple program list into an intelligent recommendation system.

## Strategic Objectives

1. **Discoverability**: Users can find programs matching their specific needs within seconds
2. **Informed Decisions**: Users have sufficient information to choose the right program
3. **Scalability**: Discovery system works for 4 programs today and 400 programs in the future
4. **API Completeness**: Frontend teams have all data needed for program selection UIs

## Themes & Initiatives

### Theme 1: Program Metadata
- **Strategic Objective**: Discoverability
- **Rationale**: Programs need standardized metadata to enable filtering. Without metadata columns, filtering would require parsing program structure at query time.
- **Initiatives**:
  - Initiative A: Add difficulty column (beginner, intermediate, advanced)
  - Initiative B: Add days_per_week column (1-7)
  - Initiative C: Add focus column (strength, hypertrophy, peaking)
  - Initiative D: Add has_amrap boolean column

### Theme 2: Program Filtering
- **Strategic Objective**: Discoverability, Scalability
- **Rationale**: Users need to narrow down programs based on their constraints (schedule, experience level) and goals (strength vs hypertrophy).
- **Initiatives**:
  - Initiative A: Filter by difficulty level
  - Initiative B: Filter by days per week
  - Initiative C: Filter by training focus
  - Initiative D: Filter by AMRAP presence (some lifters prefer or avoid AMRAP sets)
  - Initiative E: Combine multiple filters with AND logic

### Theme 3: Program Search
- **Strategic Objective**: Discoverability
- **Rationale**: Users who know a program name (or partial name) should be able to find it quickly without browsing.
- **Initiatives**:
  - Initiative A: Search by program name substring
  - Initiative B: Case-insensitive search
  - Initiative C: Search combined with filters

### Theme 4: Program Detail Enhancement
- **Strategic Objective**: Informed Decisions, API Completeness
- **Rationale**: Users need to understand what a program involves before enrolling. Sample week preview and lift requirements help users evaluate fit.
- **Initiatives**:
  - Initiative A: Include sample week preview (day names, exercise counts per day)
  - Initiative B: Include lift requirements (what lifts the program uses)
  - Initiative C: Include estimated session duration (based on exercise count)

## Success Metrics

| Metric | Target |
|--------|--------|
| Programs filterable by all 4 metadata fields | Complete |
| Programs searchable by name | Complete |
| Multiple filters combinable | Complete |
| Program detail includes sample week preview | Complete |
| Program detail includes lift requirements | Complete |
| Existing canonical programs backfilled with metadata | Complete |
| All endpoints respond in < 100ms for 100 programs | Complete |

## Timeline

| Phase | Scope |
|-------|-------|
| Now | Schema migration adding metadata columns |
| Now | Backfill migration for existing programs |
| Now | Program filtering implementation |
| Now | Program search implementation |
| Now | Program detail enhancement |
| Now | E2E tests for discovery features |

## Dependencies

- Sprint 004 (Seed Canonical Programs) - Programs must exist before discovery features
- Phase 001 core domain entities - Program schema and repository patterns

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Metadata values inconsistent across programs | Medium | Medium | Define clear metadata guidelines; review during backfill |
| Filter combinations slow on large datasets | Low | Medium | Add composite indices; test with realistic data volumes |
| Search substring matching slow | Low | Low | Use SQLite LIKE with index prefix matching |
| Estimated duration accuracy | Medium | Low | Use conservative estimates; document as approximate |

## Out of Scope

- AI-powered program recommendations - future enhancement
- User preference profiles - future enhancement
- Program comparison feature - future enhancement
- Program ratings and reviews - future enhancement
- Equipment filtering (barbell only, home gym, etc.) - future enhancement
- Program duration filtering (4 weeks vs 12 weeks) - future enhancement

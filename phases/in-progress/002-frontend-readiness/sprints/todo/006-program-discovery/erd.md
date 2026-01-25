# ERD 006: Program Discovery

## 1. Introduction / Overview

### Purpose
This document specifies the engineering requirements for implementing program discovery features in PowerPro. These features enable users to find programs through filtering, searching, and enhanced detail views.

### Scope
- Database schema changes to add program metadata columns
- Backfill migration for existing canonical programs
- Program filtering by metadata fields
- Program search by name
- Enhanced program detail response with sample week and lift requirements
- Does NOT include: AI recommendations, user preferences, program comparison, ratings

### Definitions/Glossary
| Term | Definition |
|------|------------|
| Difficulty | Program experience level: beginner, intermediate, or advanced |
| Focus | Primary training goal: strength (1-5 reps), hypertrophy (8-12 reps), or peaking (competition prep) |
| AMRAP | As Many Reps As Possible - a set type where lifter performs maximum reps |
| Sample Week | A preview of one week's structure showing days and exercise counts |
| Lift Requirements | The list of unique lifts a program uses |

### Stakeholders
- Engineering Team: Implementation
- Product Owner: Feature validation
- API Consumers (Frontend Teams): Program discovery UI implementation

## 2. Business / Stakeholder Needs

### Why
Users need to find the right program efficiently. A flat list of programs doesn't scale. Filtering by experience level, schedule constraints, and training goals enables users to narrow down options quickly. Enhanced detail views help users make informed enrollment decisions.

### Constraints
- SQLite only - queries must use SQLite-compatible SQL
- Go only - all implementation in Go
- Must extend existing schema - add columns to programs table
- Must maintain backward compatibility - existing API responses should not break
- Must backfill existing programs - Sprint 004 programs need metadata

### Success Criteria
- All filtering options work correctly with single and combined filters
- Search returns relevant results for partial name matches
- Program detail response includes sample week and lift requirements
- Canonical programs have accurate metadata values
- All queries perform within acceptable latency (< 100ms)

## 3. Functional Requirements

### Schema Changes

#### REQ-SCHEMA-001: Difficulty Column
- **Description**: The system shall add a `difficulty` column to the programs table
- **Rationale**: Enables filtering by experience level
- **Priority**: Must
- **Acceptance Criteria**:
  - Column type: TEXT with CHECK constraint
  - Valid values: 'beginner', 'intermediate', 'advanced'
  - Nullable: No (with default 'beginner' for migration safety)
  - Column is indexed for filter performance
- **Dependencies**: None

#### REQ-SCHEMA-002: Days Per Week Column
- **Description**: The system shall add a `days_per_week` column to the programs table
- **Rationale**: Enables filtering by schedule requirements
- **Priority**: Must
- **Acceptance Criteria**:
  - Column type: INTEGER with CHECK constraint
  - Valid values: 1-7
  - Nullable: No (with default 3 for migration safety)
  - Column is indexed for filter performance
- **Dependencies**: None

#### REQ-SCHEMA-003: Focus Column
- **Description**: The system shall add a `focus` column to the programs table
- **Rationale**: Enables filtering by training goal
- **Priority**: Must
- **Acceptance Criteria**:
  - Column type: TEXT with CHECK constraint
  - Valid values: 'strength', 'hypertrophy', 'peaking'
  - Nullable: No (with default 'strength' for migration safety)
  - Column is indexed for filter performance
- **Dependencies**: None

#### REQ-SCHEMA-004: Has AMRAP Column
- **Description**: The system shall add a `has_amrap` column to the programs table
- **Rationale**: Enables filtering by AMRAP preference
- **Priority**: Must
- **Acceptance Criteria**:
  - Column type: INTEGER (SQLite boolean: 0 or 1)
  - Default: 0 (false)
  - Nullable: No
  - Column is indexed for filter performance
- **Dependencies**: None

### Backfill Migration

#### REQ-BACKFILL-001: Canonical Program Metadata
- **Description**: The system shall backfill metadata for existing canonical programs
- **Rationale**: Sprint 004 programs need discovery metadata
- **Priority**: Must
- **Acceptance Criteria**:
  - Starting Strength: difficulty='beginner', days_per_week=3, focus='strength', has_amrap=0
  - Texas Method: difficulty='intermediate', days_per_week=3, focus='strength', has_amrap=0
  - Wendler 5/3/1: difficulty='intermediate', days_per_week=4, focus='strength', has_amrap=1
  - GZCLP: difficulty='beginner', days_per_week=4, focus='strength', has_amrap=1
- **Dependencies**: REQ-SCHEMA-001 through REQ-SCHEMA-004

### Program Filtering

#### REQ-FILTER-001: Filter by Difficulty
- **Description**: The system shall support filtering programs by difficulty level
- **Rationale**: Users need to find programs matching their experience
- **Priority**: Must
- **Acceptance Criteria**:
  - GET /programs?difficulty=beginner returns only beginner programs
  - GET /programs?difficulty=intermediate returns only intermediate programs
  - GET /programs?difficulty=advanced returns only advanced programs
  - Invalid difficulty value returns 400 Bad Request
- **Dependencies**: REQ-SCHEMA-001

#### REQ-FILTER-002: Filter by Days Per Week
- **Description**: The system shall support filtering programs by days per week
- **Rationale**: Users have schedule constraints
- **Priority**: Must
- **Acceptance Criteria**:
  - GET /programs?days_per_week=3 returns only 3-day programs
  - GET /programs?days_per_week=4 returns only 4-day programs
  - Value must be 1-7; invalid values return 400 Bad Request
- **Dependencies**: REQ-SCHEMA-002

#### REQ-FILTER-003: Filter by Focus
- **Description**: The system shall support filtering programs by training focus
- **Rationale**: Users have different training goals
- **Priority**: Must
- **Acceptance Criteria**:
  - GET /programs?focus=strength returns only strength-focused programs
  - GET /programs?focus=hypertrophy returns only hypertrophy programs
  - GET /programs?focus=peaking returns only peaking programs
  - Invalid focus value returns 400 Bad Request
- **Dependencies**: REQ-SCHEMA-003

#### REQ-FILTER-004: Filter by AMRAP
- **Description**: The system shall support filtering programs by AMRAP presence
- **Rationale**: Some lifters prefer or avoid AMRAP sets
- **Priority**: Must
- **Acceptance Criteria**:
  - GET /programs?has_amrap=true returns only programs with AMRAP
  - GET /programs?has_amrap=false returns only programs without AMRAP
  - Invalid boolean value returns 400 Bad Request
- **Dependencies**: REQ-SCHEMA-004

#### REQ-FILTER-005: Combined Filters
- **Description**: The system shall support combining multiple filters
- **Rationale**: Users often have multiple constraints
- **Priority**: Must
- **Acceptance Criteria**:
  - GET /programs?difficulty=beginner&days_per_week=3 returns programs matching BOTH criteria
  - All filter combinations work correctly with AND logic
  - Empty result set returns 200 with empty array, not 404
- **Dependencies**: REQ-FILTER-001 through REQ-FILTER-004

### Program Search

#### REQ-SEARCH-001: Search by Name
- **Description**: The system shall support searching programs by name substring
- **Rationale**: Users who know a program name should find it quickly
- **Priority**: Must
- **Acceptance Criteria**:
  - GET /programs?search=strength returns Starting Strength
  - GET /programs?search=531 returns Wendler 5/3/1
  - Search is case-insensitive
  - Search matches anywhere in name (not just prefix)
  - Empty search parameter is ignored (returns all programs)
- **Dependencies**: None

#### REQ-SEARCH-002: Search with Filters
- **Description**: The system shall support combining search with filters
- **Rationale**: Users may want to search within a filtered set
- **Priority**: Must
- **Acceptance Criteria**:
  - GET /programs?search=strength&difficulty=beginner returns Starting Strength
  - Search and all filters combine with AND logic
- **Dependencies**: REQ-SEARCH-001, REQ-FILTER-005

### Program Detail Enhancement

#### REQ-DETAIL-001: Sample Week Preview
- **Description**: The program detail endpoint shall include a sample week preview
- **Rationale**: Users need to understand weekly structure before enrolling
- **Priority**: Must
- **Acceptance Criteria**:
  - GET /programs/{id} includes `sampleWeek` array in response
  - Each day includes: day name/number, exercise count
  - Shows first week of program (or first cycle for cyclical programs)
  - For A/B rotation programs, shows both A and B days
- **Dependencies**: None

#### REQ-DETAIL-002: Lift Requirements
- **Description**: The program detail endpoint shall include lift requirements
- **Rationale**: Users need to know what equipment/lifts are required
- **Priority**: Must
- **Acceptance Criteria**:
  - GET /programs/{id} includes `liftRequirements` array in response
  - Array contains unique lift names used in the program
  - Lifts are sorted alphabetically
  - Example: ["Bench Press", "Deadlift", "Overhead Press", "Squat"]
- **Dependencies**: None

#### REQ-DETAIL-003: Estimated Session Duration
- **Description**: The program detail endpoint shall include estimated session duration
- **Rationale**: Helps users understand time commitment
- **Priority**: Should
- **Acceptance Criteria**:
  - GET /programs/{id} includes `estimatedSessionMinutes` in response
  - Calculation: (total sets per day * 3 minutes) + (exercises * 2 minutes warmup)
  - Value is an integer representing minutes
  - Documented as approximate estimate
- **Dependencies**: REQ-DETAIL-001

## 4. Non-Functional Requirements

### Performance
- **NFR-001**: Filter queries shall complete in < 50ms for up to 100 programs
- **NFR-002**: Search queries shall complete in < 50ms for up to 100 programs
- **NFR-003**: Combined filter+search queries shall complete in < 100ms for up to 100 programs
- **NFR-004**: Program detail enhancement shall add < 20ms to existing detail query

### Reliability
- **NFR-005**: Schema migration shall be reversible (include down migration)
- **NFR-006**: Backfill migration shall be idempotent

### Maintainability
- **NFR-007**: Filter validation shall use shared validation functions
- **NFR-008**: Query building shall use parameterized queries (no SQL injection risk)
- **NFR-009**: New columns shall be added to existing domain model and DTOs

## 5. External Interfaces

### System Interfaces
- Database: SQLite with goose migrations
- Existing tables: programs (to be extended)
- Related tables: program_days, prescriptions, lifts (for detail enhancement)

### Data Formats

#### Program List Response (enhanced)
```json
{
  "programs": [
    {
      "id": "uuid",
      "slug": "starting-strength",
      "name": "Starting Strength",
      "description": "A novice strength program...",
      "difficulty": "beginner",
      "daysPerWeek": 3,
      "focus": "strength",
      "hasAmrap": false,
      "createdAt": "2024-01-15T10:00:00Z",
      "updatedAt": "2024-01-15T10:00:00Z"
    }
  ]
}
```

#### Program Detail Response (enhanced)
```json
{
  "id": "uuid",
  "slug": "starting-strength",
  "name": "Starting Strength",
  "description": "A novice strength program...",
  "difficulty": "beginner",
  "daysPerWeek": 3,
  "focus": "strength",
  "hasAmrap": false,
  "sampleWeek": [
    {"day": 1, "name": "Workout A", "exerciseCount": 3},
    {"day": 2, "name": "Workout B", "exerciseCount": 3}
  ],
  "liftRequirements": ["Bench Press", "Deadlift", "Overhead Press", "Power Clean", "Squat"],
  "estimatedSessionMinutes": 45,
  "createdAt": "2024-01-15T10:00:00Z",
  "updatedAt": "2024-01-15T10:00:00Z"
}
```

#### Filter Parameters
| Parameter | Type | Values | Example |
|-----------|------|--------|---------|
| difficulty | string | beginner, intermediate, advanced | ?difficulty=beginner |
| days_per_week | integer | 1-7 | ?days_per_week=3 |
| focus | string | strength, hypertrophy, peaking | ?focus=strength |
| has_amrap | boolean | true, false | ?has_amrap=true |
| search | string | any | ?search=strength |

## 6. Constraints & Assumptions

### Technical Constraints
- Go backend
- SQLite database
- goose migrations
- Must extend existing programs table (not create new table)
- Must maintain backward compatibility with existing endpoints

### Assumptions
- Programs table exists from Phase 001
- Canonical programs exist from Sprint 004
- Program days and prescriptions tables exist for detail queries
- Lifts table exists for lift requirements query

## 7. Acceptance Criteria & Verification Methods

| Requirement | Verification Method | Success Criteria |
|-------------|---------------------|------------------|
| REQ-SCHEMA-001 to REQ-SCHEMA-004 | Migration tests | Columns exist with correct constraints |
| REQ-BACKFILL-001 | Query verification | Canonical programs have correct metadata |
| REQ-FILTER-001 to REQ-FILTER-005 | Integration tests | All filter combinations return correct results |
| REQ-SEARCH-001, REQ-SEARCH-002 | Integration tests | Search returns expected programs |
| REQ-DETAIL-001 to REQ-DETAIL-003 | Integration tests | Detail response includes new fields |
| NFR-001 to NFR-004 | Performance tests | Queries complete within time limits |

## 8. Prioritization / Roadmap

### Must-Have (This ERD)
- Schema migration with all 4 metadata columns
- Backfill migration for canonical programs
- All filter implementations (difficulty, days_per_week, focus, has_amrap)
- Combined filter support
- Name search implementation
- Sample week preview in detail response
- Lift requirements in detail response

### Should-Have (This ERD)
- Estimated session duration
- Composite indices for common filter combinations

### Won't-Have (Future ERDs)
- AI-powered recommendations
- User preference profiles
- Program comparison
- Ratings and reviews
- Equipment filtering
- Program duration filtering

## 9. Traceability

### Links to PRD
- PRD-006: Program Discovery

### Links to Phase Document
- Phase 002: Frontend Readiness - Theme 4 (Program Content)

### Forward Links (to Tickets)
- 001-schema-migration.md (REQ-SCHEMA-001 through REQ-SCHEMA-004)
- 002-backfill-migration.md (REQ-BACKFILL-001)
- 003-filtering-implementation.md (REQ-FILTER-001 through REQ-FILTER-005)
- 004-search-implementation.md (REQ-SEARCH-001, REQ-SEARCH-002)
- 005-detail-enhancement.md (REQ-DETAIL-001 through REQ-DETAIL-003)
- 006-e2e-tests.md (All requirements)

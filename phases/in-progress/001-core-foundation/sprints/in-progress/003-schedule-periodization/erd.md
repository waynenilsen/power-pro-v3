# ERD 003: Schedule & Periodization

## 1. Introduction / Overview

### Purpose
This document specifies the engineering requirements for the Schedule & Periodization system - the temporal organization that structures training across days, weeks, and cycles. This system answers "when do I do what" and enables percentage/prescription variation based on cycle position.

### Scope
- Day entity: named training sessions with prescription slots
- Week entity: collections of training days
- Cycle entity: repeating program units
- WeeklyLookup and DailyLookup tables for prescription variation
- Workout generation API
- Does NOT include: lift rotation, phase blocks, peaking/taper

### Definitions/Glossary
| Term | Definition |
|------|------------|
| Day | A single training session with ordered exercises |
| Week | A collection of training days within a program cycle |
| Cycle | The repeating unit of a program (e.g., 4 weeks for 5/3/1) |
| WeeklyLookup | Table mapping week number to prescription parameters |
| DailyLookup | Table mapping day identifier to prescription parameters |
| A/B Rotation | Alternating between two day templates |
| Training Session | A resolved workout for a specific calendar date |

### Stakeholders
- Engineering Team: Implementation
- Product Owner: Feature validation
- API Consumers: Workout generation and scheduling

## 2. Business / Stakeholder Needs

### Why
Programs organize training temporally. 5/3/1 varies percentages by week within a 4-week cycle. Bill Starr varies intensity by day (Heavy/Light/Medium). Greg Nuckols combines both. Without scheduling primitives, each program requires custom implementation.

### Constraints
- Must integrate with Prescription system from ERD-002
- Must support diverse cycle lengths (1, 3, 4+ weeks)
- Must support lookup table integration with LoadStrategy

### Success Criteria
- All 5 Phase 1 programs can be represented using Day/Week/Cycle + Lookups
- Workout generation returns correct prescriptions for any date
- Lookup tables correctly modify prescription parameters

## 3. Functional Requirements

### Day Entity

#### REQ-DAY-001: Day Identification
- **Description**: The system shall provide unique identification for each day template
- **Rationale**: Enables reference from weeks and workout generation
- **Priority**: Must
- **Acceptance Criteria**: Each day has a unique UUID identifier
- **Dependencies**: None

#### REQ-DAY-002: Day Name
- **Description**: The system shall store a name for each day
- **Rationale**: Human-readable identification (e.g., "Day A", "Heavy Day", "Monday")
- **Priority**: Must
- **Acceptance Criteria**: Name required, max 50 characters
- **Dependencies**: REQ-DAY-001

#### REQ-DAY-003: Day Prescriptions
- **Description**: The system shall associate ordered prescriptions with each day
- **Rationale**: A day is defined by its exercises in order
- **Priority**: Must
- **Acceptance Criteria**:
  - Many-to-many relationship with Prescription (via join table)
  - Order field on relationship
  - Same prescription can appear in multiple days
- **Dependencies**: REQ-DAY-001, ERD-002 (Prescription)

#### REQ-DAY-004: Day Metadata
- **Description**: The system shall support optional metadata on days
- **Rationale**: Programs may have day-specific parameters (intensity level, focus)
- **Priority**: Should
- **Acceptance Criteria**:
  - JSONB metadata field
  - Common keys: intensityLevel (HEAVY/LIGHT/MEDIUM), focus (string)
- **Dependencies**: REQ-DAY-001

#### REQ-DAY-005: Day Slug
- **Description**: The system shall store a URL-safe identifier for each day
- **Rationale**: Enables clean API URLs and day lookup in lookups
- **Priority**: Must
- **Acceptance Criteria**: Lowercase alphanumeric with hyphens, unique within program
- **Dependencies**: REQ-DAY-001

### Week Entity

#### REQ-WEEK-001: Week Identification
- **Description**: The system shall provide unique identification for each week template
- **Rationale**: Enables reference from cycles and lookup tables
- **Priority**: Must
- **Acceptance Criteria**: Each week has a unique UUID identifier
- **Dependencies**: None

#### REQ-WEEK-002: Week Number
- **Description**: The system shall track week position within a cycle
- **Rationale**: Enables week-based lookup tables (Week 1 = 65%, Week 2 = 70%)
- **Priority**: Must
- **Acceptance Criteria**: Integer >= 1; unique within cycle
- **Dependencies**: REQ-WEEK-001

#### REQ-WEEK-003: Week Days
- **Description**: The system shall associate ordered days with each week
- **Rationale**: A week contains specific training days
- **Priority**: Must
- **Acceptance Criteria**:
  - Many-to-many relationship with Day (via join table)
  - Day-of-week field (1-7 or named enum: MONDAY-SUNDAY)
  - Multiple days can map to same day-of-week (for multiple sessions)
- **Dependencies**: REQ-WEEK-001, REQ-DAY-001

#### REQ-WEEK-004: A/B Week Support
- **Description**: The system shall support week alternation
- **Rationale**: Some programs alternate week patterns (Week A / Week B)
- **Priority**: Should
- **Acceptance Criteria**:
  - Optional variant field (A, B, or null)
  - Cycle can specify which variant to use per week number
- **Dependencies**: REQ-WEEK-001

### Cycle Entity

#### REQ-CYCLE-001: Cycle Identification
- **Description**: The system shall provide unique identification for each cycle
- **Rationale**: Enables program definition and user progress tracking
- **Priority**: Must
- **Acceptance Criteria**: Each cycle has a unique UUID identifier
- **Dependencies**: None

#### REQ-CYCLE-002: Cycle Length
- **Description**: The system shall track the number of weeks in a cycle
- **Rationale**: Different programs have different cycle lengths
- **Priority**: Must
- **Acceptance Criteria**: Integer >= 1 (1-week, 3-week, 4-week cycles supported)
- **Dependencies**: REQ-CYCLE-001

#### REQ-CYCLE-003: Cycle Weeks
- **Description**: The system shall associate weeks with each cycle
- **Rationale**: A cycle is composed of week templates
- **Priority**: Must
- **Acceptance Criteria**:
  - One-to-many relationship with Week
  - Week number determines order
  - Count of weeks must equal cycle length
- **Dependencies**: REQ-CYCLE-001, REQ-WEEK-001

#### REQ-CYCLE-004: Cycle Name
- **Description**: The system shall store a name for each cycle
- **Rationale**: Human identification of program structure
- **Priority**: Must
- **Acceptance Criteria**: Name required, max 100 characters
- **Dependencies**: REQ-CYCLE-001

### Lookup Tables

#### REQ-LOOKUP-001: WeeklyLookup Entity
- **Description**: The system shall implement weekly lookup tables
- **Rationale**: Programs vary percentages by week (5/3/1: Week 1 = 65/75/85%)
- **Priority**: Must
- **Acceptance Criteria**:
  - Maps week number to parameters (percentage adjustments, rep targets)
  - Can be associated with Prescription or LoadStrategy
  - Returns appropriate values given week number
- **Dependencies**: REQ-WEEK-002

#### REQ-LOOKUP-002: WeeklyLookup Structure
- **Description**: The system shall define weekly lookup data structure
- **Rationale**: Standardized structure enables consistent processing
- **Priority**: Must
- **Acceptance Criteria**:
  - Array of (weekNumber, parameters) entries
  - Parameters: percentageModifier (additive or multiplicative), repTarget
  - Example: Week 1 = {percentages: [65, 75, 85], reps: 5}
- **Dependencies**: REQ-LOOKUP-001

#### REQ-LOOKUP-003: DailyLookup Entity
- **Description**: The system shall implement daily lookup tables
- **Rationale**: Programs vary percentages by day (Bill Starr: Heavy/Light/Medium)
- **Priority**: Must
- **Acceptance Criteria**:
  - Maps day slug or day-of-week to parameters
  - Returns appropriate values given day identifier
- **Dependencies**: REQ-DAY-005

#### REQ-LOOKUP-004: DailyLookup Structure
- **Description**: The system shall define daily lookup data structure
- **Rationale**: Standardized structure enables consistent processing
- **Priority**: Must
- **Acceptance Criteria**:
  - Map of dayIdentifier to parameters
  - Parameters: percentageModifier, intensityLevel
  - Example: {heavy: 100%, light: 70%, medium: 80%}
- **Dependencies**: REQ-LOOKUP-003

#### REQ-LOOKUP-005: Lookup Integration
- **Description**: The system shall integrate lookups with LoadStrategy resolution
- **Rationale**: Lookups must affect final weight calculation
- **Priority**: Must
- **Acceptance Criteria**:
  - PercentOf strategy accepts optional lookup reference
  - During resolution, lookup values modify base percentage
  - Lookup context (week number, day) passed through resolution chain
- **Dependencies**: REQ-LOOKUP-001, REQ-LOOKUP-003, ERD-002 (LoadStrategy)

### User Program State

#### REQ-STATE-001: User Cycle State
- **Description**: The system shall track user's current position in program
- **Rationale**: Enables workout generation for "today" and progression triggers
- **Priority**: Must
- **Acceptance Criteria**:
  - UserProgramState entity: userId, cycleId, currentWeek, currentCycleIteration
  - currentWeek: 1 to cycle.length
  - currentCycleIteration: which time through the cycle (1, 2, 3...)
- **Dependencies**: REQ-CYCLE-001

#### REQ-STATE-002: State Advancement
- **Description**: The system shall advance user state after completing workouts
- **Rationale**: Tracks progress through program
- **Priority**: Must
- **Acceptance Criteria**:
  - Endpoint: POST /users/{userId}/program-state/advance
  - Advances to next day, wrapping to next week, wrapping to next cycle
  - Returns new state
- **Dependencies**: REQ-STATE-001

#### REQ-STATE-003: Cycle Completion Detection
- **Description**: The system shall detect when a user completes a cycle
- **Rationale**: Triggers cycle-based progressions
- **Priority**: Must
- **Acceptance Criteria**:
  - Event/flag when advancing from week N to week 1
  - Can be used as progression trigger
- **Dependencies**: REQ-STATE-002

### Workout Generation

#### REQ-GEN-001: Generate Workout
- **Description**: The system shall generate workouts for a user
- **Rationale**: Core API function - returns today's training
- **Priority**: Must
- **Acceptance Criteria**:
  - Endpoint: GET /users/{userId}/workout
  - Optional query params: date, weekNumber, daySlug
  - Returns resolved prescriptions for the specified day
  - Uses current program state if no overrides provided
- **Dependencies**: REQ-STATE-001, ERD-002 (Prescription resolution)

#### REQ-GEN-002: Workout with Lookup Context
- **Description**: The system shall apply lookups during workout generation
- **Rationale**: Generated workouts must reflect week/day specific parameters
- **Priority**: Must
- **Acceptance Criteria**:
  - Week number from state passed to WeeklyLookup
  - Day identifier passed to DailyLookup
  - Resolved weights reflect lookup modifications
- **Dependencies**: REQ-GEN-001, REQ-LOOKUP-005

#### REQ-GEN-003: Workout Preview
- **Description**: The system shall support previewing future workouts
- **Rationale**: Users may want to see upcoming training
- **Priority**: Should
- **Acceptance Criteria**:
  - Endpoint: GET /users/{userId}/workout/preview?week={n}&day={slug}
  - Does not require state advancement
  - Returns resolved prescriptions for specified position
- **Dependencies**: REQ-GEN-001

### Program Definition

#### REQ-PROG-001: Program Entity
- **Description**: The system shall define a Program as a named configuration
- **Rationale**: Programs bundle cycles, lookups, and defaults
- **Priority**: Must
- **Acceptance Criteria**:
  - Program entity: id, name, slug, description, cycleId, defaultRounding
  - One program references one cycle
  - Program can reference lookup tables
- **Dependencies**: REQ-CYCLE-001

#### REQ-PROG-002: Program CRUD
- **Description**: The system shall provide API endpoints for program management
- **Rationale**: Required for program configuration
- **Priority**: Must
- **Acceptance Criteria**:
  - GET /programs - list programs
  - GET /programs/{id} - get program with full structure
  - POST /programs - create program
  - PUT /programs/{id} - update program
  - DELETE /programs/{id} - delete program (fails if users enrolled)
- **Dependencies**: REQ-PROG-001

#### REQ-PROG-003: User Program Enrollment
- **Description**: The system shall support enrolling users in programs
- **Rationale**: Users follow specific programs
- **Priority**: Must
- **Acceptance Criteria**:
  - Endpoint: POST /users/{userId}/program
  - Creates UserProgramState with initial position
  - User can only be enrolled in one program at a time
- **Dependencies**: REQ-PROG-001, REQ-STATE-001

## 4. Non-Functional Requirements

### Performance
- **NFR-001**: Workout generation shall complete in < 500ms (p95)
- **NFR-002**: State advancement shall complete in < 100ms (p95)
- **NFR-003**: Program CRUD shall complete in < 100ms (p95)

### Reliability
- **NFR-004**: State advancement shall be atomic (no partial updates)
- **NFR-005**: Cycle completion detection shall be reliable (no missed triggers)

### Maintainability
- **NFR-006**: Lookup tables shall be addable without code changes
- **NFR-007**: New cycle lengths shall work without code changes

## 5. External Interfaces

### Data Formats

#### Day API Response
```json
{
  "id": "uuid",
  "name": "Day A",
  "slug": "day-a",
  "metadata": {
    "intensityLevel": "HEAVY"
  },
  "prescriptions": [
    {
      "id": "uuid",
      "order": 1,
      "lift": { "id": "uuid", "name": "Squat" }
    }
  ]
}
```

#### Week API Response
```json
{
  "id": "uuid",
  "weekNumber": 1,
  "variant": null,
  "days": [
    {
      "dayOfWeek": "MONDAY",
      "day": { "id": "uuid", "name": "Day A", "slug": "day-a" }
    },
    {
      "dayOfWeek": "WEDNESDAY",
      "day": { "id": "uuid", "name": "Day B", "slug": "day-b" }
    },
    {
      "dayOfWeek": "FRIDAY",
      "day": { "id": "uuid", "name": "Day A", "slug": "day-a" }
    }
  ]
}
```

#### Cycle API Response
```json
{
  "id": "uuid",
  "name": "5/3/1 Cycle",
  "lengthWeeks": 4,
  "weeks": [
    { "weekNumber": 1, "id": "uuid" },
    { "weekNumber": 2, "id": "uuid" },
    { "weekNumber": 3, "id": "uuid" },
    { "weekNumber": 4, "id": "uuid" }
  ]
}
```

#### WeeklyLookup Example
```json
{
  "id": "uuid",
  "name": "5/3/1 Percentages",
  "entries": [
    { "weekNumber": 1, "percentages": [65, 75, 85], "reps": [5, 5, 5] },
    { "weekNumber": 2, "percentages": [70, 80, 90], "reps": [3, 3, 3] },
    { "weekNumber": 3, "percentages": [75, 85, 95], "reps": [5, 3, 1] },
    { "weekNumber": 4, "percentages": [40, 50, 60], "reps": [5, 5, 5] }
  ]
}
```

#### Generated Workout Response
```json
{
  "userId": "uuid",
  "programId": "uuid",
  "cycleIteration": 1,
  "weekNumber": 2,
  "daySlug": "day-a",
  "date": "2024-01-15",
  "exercises": [
    {
      "prescriptionId": "uuid",
      "lift": { "id": "uuid", "name": "Squat", "slug": "squat" },
      "sets": [
        { "setNumber": 1, "weight": 185, "targetReps": 3, "isWorkSet": false },
        { "setNumber": 2, "weight": 215, "targetReps": 3, "isWorkSet": false },
        { "setNumber": 3, "weight": 245, "targetReps": 3, "isWorkSet": true }
      ],
      "notes": null,
      "restSeconds": 180
    }
  ]
}
```

## 6. Constraints & Assumptions

### Technical Constraints
- Day/Week/Cycle stored as relational entities (not embedded JSONB)
- Lookups stored as JSONB for flexibility
- State tracking requires reliable timestamp handling

### Assumptions
- User entity exists with authentication
- Prescription and LiftMax entities exist (ERD-001, ERD-002)
- All users in same timezone (or timezone handled at API layer)

## 7. Acceptance Criteria & Verification Methods

| Requirement | Verification Method | Success Criteria |
|-------------|---------------------|------------------|
| REQ-DAY-*, REQ-WEEK-*, REQ-CYCLE-* | Unit tests, Integration tests | All entity operations work correctly |
| REQ-LOOKUP-* | Unit tests | Lookups return correct values |
| REQ-STATE-* | Integration tests | State advances correctly through program |
| REQ-GEN-* | Integration tests | Generated workouts match expected output |
| NFR-001 | Load testing | Workout generation < 500ms |

## 8. Prioritization / Roadmap

### Must-Have (This ERD)
- Day, Week, Cycle entities with relationships
- WeeklyLookup and DailyLookup
- User program state tracking
- Workout generation with lookup integration

### Should-Have (This ERD)
- A/B week variants
- Day metadata
- Workout preview endpoint

### Won't-Have (Future ERDs)
- Lift rotation (Phase 5)
- Phase blocks (Phase 5)
- DaysOut countdown (Phase 10)
- Taper protocols (Phase 10)

## 9. Traceability

### Links to PRD
- PRD-003: Schedule & Periodization

### Links to Phase Document
- Phase 001: Core Foundation - Theme 4 (Schedule & Periodization)

### Dependencies
- ERD-001: Core Domain Entities (Lift, LiftMax)
- ERD-002: Prescription System (Prescription, LoadStrategy, SetScheme)

### Forward Links (to Tickets)
- To be created after ERD approval

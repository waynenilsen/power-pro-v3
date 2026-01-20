# ERD 001: Core Domain Entities

## 1. Introduction / Overview

### Purpose
This document specifies the engineering requirements for implementing the foundational domain entities of PowerPro: Lift and LiftMax. These entities serve as the atomic building blocks upon which all other features depend.

### Scope
- Lift entity: representation of exercises
- LiftMax entity: user-specific reference values for load calculation
- API endpoints for CRUD operations on both entities
- Does NOT include: E1RM calculation, RPE integration, historical tracking

### Definitions/Glossary
| Term | Definition |
|------|------------|
| Lift | A specific exercise (e.g., Squat, Bench Press, Deadlift) |
| LiftMax | A user's reference number for a lift used in load calculation |
| 1RM | One-rep max - the maximum weight a lifter can lift for one repetition |
| TM | Training Max - typically 85-90% of 1RM, used as the basis for program percentages |
| Competition Lift | The three main powerlifting lifts: Squat, Bench Press, Deadlift |
| Variation | A modified version of a lift (e.g., Pause Squat, Close-Grip Bench) |

### Stakeholders
- Engineering Team: Implementation
- Product Owner: Feature validation
- API Consumers: Integration and usage

## 2. Business / Stakeholder Needs

### Why
The Lift and LiftMax entities are required by every downstream feature: prescriptions reference lifts, load strategies reference maxes, progressions mutate maxes. Without these foundational entities, no program can be represented.

### Constraints
- Must support extensibility for future lift types and max types
- Must maintain referential integrity between lifts and maxes
- API must be RESTful and follow established conventions

### Success Criteria
- All 5 target programs (Starting Strength, Bill Starr 5x5, 5/3/1 BBB, Sheiko Beginner, Greg Nuckols HF) can reference lifts and maxes without special cases
- Test coverage > 90% for domain logic
- API response time < 100ms for entity operations

## 3. Functional Requirements

### Lift Entity

#### REQ-LIFT-001: Lift Identification
- **Description**: The system shall provide unique identification for each lift
- **Rationale**: Enables unambiguous reference from prescriptions and maxes
- **Priority**: Must
- **Acceptance Criteria**: Each lift has a unique identifier; duplicate identifiers are rejected
- **Dependencies**: None

#### REQ-LIFT-002: Lift Name
- **Description**: The system shall store a human-readable name for each lift
- **Rationale**: Required for display in API responses and user interfaces
- **Priority**: Must
- **Acceptance Criteria**: Name is required, non-empty, max 100 characters
- **Dependencies**: REQ-LIFT-001

#### REQ-LIFT-003: Lift Slug
- **Description**: The system shall store a URL-safe slug for each lift
- **Rationale**: Enables RESTful resource identification and clean URLs
- **Priority**: Must
- **Acceptance Criteria**: Slug is unique, lowercase, alphanumeric with hyphens only
- **Dependencies**: REQ-LIFT-001

#### REQ-LIFT-004: Competition Lift Flag
- **Description**: The system shall track whether a lift is a competition lift
- **Rationale**: Competition lifts (SBD) have special significance in programming
- **Priority**: Must
- **Acceptance Criteria**: Boolean flag, default false; Squat, Bench, Deadlift are true
- **Dependencies**: REQ-LIFT-001

#### REQ-LIFT-005: Parent Lift Reference
- **Description**: The system shall support linking variation lifts to their parent lift
- **Rationale**: Enables grouping variations (Pause Squat → Squat) for max inheritance
- **Priority**: Should
- **Acceptance Criteria**: Optional reference to another lift; circular references rejected
- **Dependencies**: REQ-LIFT-001

#### REQ-LIFT-006: Lift CRUD Operations
- **Description**: The system shall provide API endpoints for lift management
- **Rationale**: Required for program configuration and administration
- **Priority**: Must
- **Acceptance Criteria**:
  - GET /lifts - list all lifts
  - GET /lifts/{id} - get single lift
  - POST /lifts - create lift
  - PUT /lifts/{id} - update lift
  - DELETE /lifts/{id} - delete lift (fails if referenced by maxes)
- **Dependencies**: REQ-LIFT-001 through REQ-LIFT-005

### LiftMax Entity

#### REQ-MAX-001: LiftMax Identification
- **Description**: The system shall provide unique identification for each LiftMax
- **Rationale**: Enables unambiguous reference from load strategies
- **Priority**: Must
- **Acceptance Criteria**: Each LiftMax has a unique identifier
- **Dependencies**: None

#### REQ-MAX-002: User Association
- **Description**: The system shall associate each LiftMax with a specific user
- **Rationale**: Maxes are user-specific; different users have different maxes
- **Priority**: Must
- **Acceptance Criteria**: User reference required; foreign key to user entity
- **Dependencies**: REQ-MAX-001

#### REQ-MAX-003: Lift Association
- **Description**: The system shall associate each LiftMax with a specific lift
- **Rationale**: A max is always for a specific lift
- **Priority**: Must
- **Acceptance Criteria**: Lift reference required; foreign key to lift entity
- **Dependencies**: REQ-MAX-001, REQ-LIFT-001

#### REQ-MAX-004: Max Type
- **Description**: The system shall support multiple max types (1RM, TM)
- **Rationale**: Different programs use different reference types for percentages
- **Priority**: Must
- **Acceptance Criteria**:
  - Enum type with values: ONE_RM, TRAINING_MAX
  - Required field, no default
- **Dependencies**: REQ-MAX-001

#### REQ-MAX-005: Max Value
- **Description**: The system shall store the weight value for each max
- **Rationale**: The actual number used in load calculations
- **Priority**: Must
- **Acceptance Criteria**:
  - Positive decimal value
  - Unit-agnostic (stored as raw number, unit handling at API layer)
  - Precision to 0.25 (quarter pounds/kg)
- **Dependencies**: REQ-MAX-001

#### REQ-MAX-006: Effective Date
- **Description**: The system shall track when each max value became effective
- **Rationale**: Enables temporal queries and progression tracking
- **Priority**: Must
- **Acceptance Criteria**: Required timestamp; defaults to creation time
- **Dependencies**: REQ-MAX-001

#### REQ-MAX-007: Max Uniqueness
- **Description**: The system shall enforce unique (user, lift, type, effective_date) combinations
- **Rationale**: Prevents duplicate max entries for the same context
- **Priority**: Must
- **Acceptance Criteria**: Database constraint enforces uniqueness; API returns 409 on conflict
- **Dependencies**: REQ-MAX-002, REQ-MAX-003, REQ-MAX-004, REQ-MAX-006

#### REQ-MAX-008: TM Validation
- **Description**: The system shall validate Training Max is within expected range of 1RM
- **Rationale**: TM outside 80-95% of 1RM may indicate data entry error
- **Priority**: Should
- **Acceptance Criteria**:
  - Warning (not error) if TM < 80% or > 95% of existing 1RM for same lift
  - Warning logged but creation/update proceeds
- **Dependencies**: REQ-MAX-004, REQ-MAX-005

#### REQ-MAX-009: Max Conversion
- **Description**: The system shall provide conversion between 1RM and TM
- **Rationale**: Users may have one type and need to calculate the other
- **Priority**: Should
- **Acceptance Criteria**:
  - Endpoint: GET /lift-maxes/{id}/convert?to_type={type}&percentage={pct}
  - Default TM percentage: 90%
  - Returns calculated value without persisting
- **Dependencies**: REQ-MAX-004, REQ-MAX-005

#### REQ-MAX-010: LiftMax CRUD Operations
- **Description**: The system shall provide API endpoints for LiftMax management
- **Rationale**: Required for user onboarding and progression tracking
- **Priority**: Must
- **Acceptance Criteria**:
  - GET /users/{userId}/lift-maxes - list user's maxes
  - GET /lift-maxes/{id} - get single max
  - POST /users/{userId}/lift-maxes - create max
  - PUT /lift-maxes/{id} - update max
  - DELETE /lift-maxes/{id} - delete max
- **Dependencies**: REQ-MAX-001 through REQ-MAX-007

#### REQ-MAX-011: Current Max Query
- **Description**: The system shall provide efficient lookup of current max value
- **Rationale**: Load calculation requires finding the most recent max for a given (user, lift, type)
- **Priority**: Must
- **Acceptance Criteria**:
  - Endpoint: GET /users/{userId}/lift-maxes/current?lift={liftId}&type={type}
  - Returns most recent max by effective_date
  - Returns 404 if no max exists
- **Dependencies**: REQ-MAX-002, REQ-MAX-003, REQ-MAX-004, REQ-MAX-006

## 4. Non-Functional Requirements

### Performance
- **NFR-001**: Entity CRUD operations shall respond in < 100ms (p95)
- **NFR-002**: Current max lookup shall respond in < 50ms (p95)
- **NFR-003**: List operations shall support pagination with default page size of 20

### Reliability
- **NFR-004**: Database transactions shall ensure consistency (no partial updates)
- **NFR-005**: Referential integrity shall prevent orphaned maxes

### Security
- **NFR-006**: LiftMax data shall be accessible only to the owning user or admin
- **NFR-007**: Lift data shall be readable by all authenticated users

### Maintainability
- **NFR-008**: Domain logic shall be isolated from API layer
- **NFR-009**: Entity changes shall be validated at domain layer, not just API layer

## 5. External Interfaces

### System Interfaces
- Database: PostgreSQL with TypeORM entities
- API: REST over HTTP/HTTPS

### Data Formats

#### Lift API Response
```json
{
  "id": "uuid",
  "name": "Squat",
  "slug": "squat",
  "isCompetitionLift": true,
  "parentLiftId": null,
  "createdAt": "2024-01-15T10:00:00Z",
  "updatedAt": "2024-01-15T10:00:00Z"
}
```

#### LiftMax API Response
```json
{
  "id": "uuid",
  "userId": "uuid",
  "liftId": "uuid",
  "type": "TRAINING_MAX",
  "value": 315.0,
  "effectiveDate": "2024-01-15T10:00:00Z",
  "createdAt": "2024-01-15T10:00:00Z",
  "updatedAt": "2024-01-15T10:00:00Z"
}
```

## 6. Constraints & Assumptions

### Technical Constraints
- TypeScript/Node.js backend
- PostgreSQL database
- RESTful API design

### Assumptions
- User entity exists and is referenced by UUID
- Authentication/authorization system is in place
- Database migrations are managed via standard tooling

## 7. Acceptance Criteria & Verification Methods

| Requirement | Verification Method | Success Criteria |
|-------------|---------------------|------------------|
| REQ-LIFT-001 to REQ-LIFT-006 | Unit tests, Integration tests | All lift CRUD operations work correctly |
| REQ-MAX-001 to REQ-MAX-011 | Unit tests, Integration tests | All max CRUD operations work correctly |
| NFR-001, NFR-002 | Load testing | Response times within specified limits |
| NFR-006, NFR-007 | Security tests | Authorization rules enforced correctly |

## 8. Prioritization / Roadmap

### Must-Have (This ERD)
- Lift entity with all fields
- LiftMax entity with 1RM and TM types
- All CRUD operations
- Current max lookup

### Should-Have (This ERD)
- Parent lift reference
- TM validation warning
- Max conversion endpoint

### Won't-Have (Future ERDs)
- E1RM calculation
- xRM support
- RPE-based estimation
- Historical max graphing

## 9. Traceability

### Links to PRD
- PRD-001: Core Domain Entities

### Links to Phase Document
- Phase 001: Core Foundation - Theme 1 (Core Domain Entities)

### Forward Links (to Tickets)
- To be created after ERD approval

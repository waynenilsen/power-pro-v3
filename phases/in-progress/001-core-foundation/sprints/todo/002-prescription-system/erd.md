# ERD 002: Prescription System

## 1. Introduction / Overview

### Purpose
This document specifies the engineering requirements for the Prescription system - the core abstraction that defines what a lifter should do for a single exercise slot. A prescription combines a lift, load calculation strategy, and set/rep scheme.

### Scope
- Prescription entity linking lifts to load and set specifications
- PercentOf LoadStrategy for percentage-based load calculation
- Fixed and Ramp SetScheme variants
- Weight rounding system
- Does NOT include: AMRAP, RPE-based loading, dynamic set schemes

### Definitions/Glossary
| Term | Definition |
|------|------------|
| Prescription | A complete instruction for one exercise: lift + load + sets/reps |
| LoadStrategy | Algorithm for calculating the target weight |
| SetScheme | Algorithm for determining sets, reps, and their structure |
| PercentOf | Load strategy that calculates weight as percentage of a max |
| Fixed | Set scheme with constant sets and reps (e.g., 5x5) |
| Ramp | Set scheme with progressive percentages across sets |
| Rounding | Adjusting calculated weight to nearest available plate increment |

### Stakeholders
- Engineering Team: Implementation
- Product Owner: Feature validation
- API Consumers: Integration and workout generation

## 2. Business / Stakeholder Needs

### Why
Programs are composed of prescriptions. Without a flexible prescription system, each program would require custom code. The prescription abstraction enables DRY representation of diverse programming styles through composition.

### Constraints
- Must compose with existing Lift and LiftMax entities
- LoadStrategy and SetScheme must be extensible for future phases
- Must support all 5 Phase 1 target programs

### Success Criteria
- Starting Strength, Bill Starr 5x5, 5/3/1 BBB, Sheiko Beginner, and Greg Nuckols HF can all be represented using Prescription + PercentOf + Fixed/Ramp
- Test coverage > 90%
- Workout generation < 200ms

## 3. Functional Requirements

### Prescription Entity

#### REQ-PRSC-001: Prescription Identification
- **Description**: The system shall provide unique identification for each prescription
- **Rationale**: Enables unambiguous reference from schedules and logs
- **Priority**: Must
- **Acceptance Criteria**: Each prescription has a unique UUID identifier
- **Dependencies**: None

#### REQ-PRSC-002: Lift Association
- **Description**: The system shall associate each prescription with a specific lift
- **Rationale**: A prescription always specifies which exercise to perform
- **Priority**: Must
- **Acceptance Criteria**: Lift reference required; foreign key to Lift entity
- **Dependencies**: ERD-001 (Lift entity)

#### REQ-PRSC-003: LoadStrategy Association
- **Description**: The system shall associate each prescription with a LoadStrategy
- **Rationale**: Defines how to calculate the target weight
- **Priority**: Must
- **Acceptance Criteria**: LoadStrategy reference required; polymorphic association
- **Dependencies**: REQ-LOAD-001

#### REQ-PRSC-004: SetScheme Association
- **Description**: The system shall associate each prescription with a SetScheme
- **Rationale**: Defines the sets/reps structure
- **Priority**: Must
- **Acceptance Criteria**: SetScheme reference required; polymorphic association
- **Dependencies**: REQ-SET-001

#### REQ-PRSC-005: Prescription Ordering
- **Description**: The system shall support ordering prescriptions within a context
- **Rationale**: Exercise order matters for programming
- **Priority**: Must
- **Acceptance Criteria**: Integer order field; default 0; unique within parent context
- **Dependencies**: REQ-PRSC-001

#### REQ-PRSC-006: Prescription Notes
- **Description**: The system shall support optional notes/cues on prescriptions
- **Rationale**: Coaches may include technique cues or special instructions
- **Priority**: Should
- **Acceptance Criteria**: Optional text field; max 500 characters
- **Dependencies**: REQ-PRSC-001

#### REQ-PRSC-007: Rest Period
- **Description**: The system shall support optional rest period specification
- **Rationale**: Some programs prescribe specific rest periods
- **Priority**: Should
- **Acceptance Criteria**: Optional integer field (seconds); null means "use default"
- **Dependencies**: REQ-PRSC-001

### LoadStrategy - PercentOf

#### REQ-LOAD-001: LoadStrategy Interface
- **Description**: The system shall define a common interface for all load strategies
- **Rationale**: Enables polymorphic load calculation
- **Priority**: Must
- **Acceptance Criteria**:
  - Interface: `calculateLoad(userId, liftId, context) -> weight`
  - All strategies implement this interface
- **Dependencies**: None

#### REQ-LOAD-002: PercentOf Strategy
- **Description**: The system shall implement PercentOf load strategy
- **Rationale**: Most programs use percentage-based loading
- **Priority**: Must
- **Acceptance Criteria**:
  - Parameters: reference type (1RM or TM), percentage (0-100+)
  - Fetches current max for user/lift/type
  - Returns max * (percentage / 100)
- **Dependencies**: REQ-LOAD-001, ERD-001 (LiftMax)

#### REQ-LOAD-003: PercentOf Reference Type
- **Description**: The system shall support configurable reference max type
- **Rationale**: Different programs reference 1RM vs TM
- **Priority**: Must
- **Acceptance Criteria**: Enum parameter: ONE_RM or TRAINING_MAX
- **Dependencies**: REQ-LOAD-002

#### REQ-LOAD-004: Weight Rounding
- **Description**: The system shall round calculated weights to specified increment
- **Rationale**: Gyms have specific plate increments; unrounded weights are impractical
- **Priority**: Must
- **Acceptance Criteria**:
  - Configurable rounding increment (e.g., 2.5, 5.0)
  - Rounds to nearest increment (standard rounding)
  - Example: 142.5 with 5.0 increment → 145
- **Dependencies**: REQ-LOAD-002

#### REQ-LOAD-005: Rounding Direction
- **Description**: The system shall support configurable rounding direction
- **Rationale**: Some prefer conservative (down) vs standard (nearest) rounding
- **Priority**: Should
- **Acceptance Criteria**:
  - Enum: NEAREST, DOWN, UP
  - Default: NEAREST
- **Dependencies**: REQ-LOAD-004

#### REQ-LOAD-006: PercentOf Validation
- **Description**: The system shall validate PercentOf parameters
- **Rationale**: Invalid percentages should be rejected early
- **Priority**: Must
- **Acceptance Criteria**:
  - Percentage must be > 0
  - Percentage > 100 allowed (for overload work)
  - Reference type must be valid enum value
- **Dependencies**: REQ-LOAD-002

### SetScheme - Fixed and Ramp

#### REQ-SET-001: SetScheme Interface
- **Description**: The system shall define a common interface for all set schemes
- **Rationale**: Enables polymorphic set generation
- **Priority**: Must
- **Acceptance Criteria**:
  - Interface: `generateSets(baseWeight, context) -> SetList`
  - SetList contains: weight, targetReps, setNumber, isWorkSet flag
  - All schemes implement this interface
- **Dependencies**: None

#### REQ-SET-002: Fixed SetScheme
- **Description**: The system shall implement Fixed set scheme
- **Rationale**: Most common scheme: same weight and reps for all sets (5x5)
- **Priority**: Must
- **Acceptance Criteria**:
  - Parameters: sets (integer), reps (integer)
  - Generates `sets` number of identical sets at baseWeight for `reps`
  - All sets are work sets
- **Dependencies**: REQ-SET-001

#### REQ-SET-003: Ramp SetScheme
- **Description**: The system shall implement Ramp set scheme
- **Rationale**: Warmup progressions and Bill Starr style ramping
- **Priority**: Must
- **Acceptance Criteria**:
  - Parameters: array of (percentage, reps) pairs
  - Generates one set per pair: baseWeight * percentage, specified reps
  - Final set (100%) is work set; others are warmup
- **Dependencies**: REQ-SET-001

#### REQ-SET-004: Ramp Set Classification
- **Description**: The system shall classify ramp sets as warmup or work
- **Rationale**: Enables filtering work sets for volume tracking
- **Priority**: Must
- **Acceptance Criteria**:
  - Sets below configurable threshold (default 80%) are warmup
  - Work set threshold configurable per scheme instance
- **Dependencies**: REQ-SET-003

#### REQ-SET-005: SetScheme Validation
- **Description**: The system shall validate set scheme parameters
- **Rationale**: Invalid parameters should be rejected early
- **Priority**: Must
- **Acceptance Criteria**:
  - Fixed: sets >= 1, reps >= 1
  - Ramp: at least one percentage/rep pair, percentages > 0
- **Dependencies**: REQ-SET-002, REQ-SET-003

### Prescription Resolution

#### REQ-PRSC-008: Prescription Resolution
- **Description**: The system shall resolve a prescription to concrete sets
- **Rationale**: API consumers need the actual workout, not abstract prescription
- **Priority**: Must
- **Acceptance Criteria**:
  - Endpoint: POST /prescriptions/{id}/resolve
  - Input: userId (for max lookup)
  - Output: array of sets with weight, reps, setNumber, isWorkSet
  - Applies LoadStrategy then SetScheme
- **Dependencies**: REQ-PRSC-001 through REQ-PRSC-004, REQ-LOAD-002, REQ-SET-002, REQ-SET-003

#### REQ-PRSC-009: Batch Resolution
- **Description**: The system shall support resolving multiple prescriptions at once
- **Rationale**: Workout generation requires resolving all prescriptions for a session
- **Priority**: Must
- **Acceptance Criteria**:
  - Endpoint: POST /prescriptions/resolve-batch
  - Input: array of prescription IDs + userId
  - Output: array of resolved prescription results
  - Returns partial results if some fail
- **Dependencies**: REQ-PRSC-008

#### REQ-PRSC-010: Resolution Caching
- **Description**: The system shall cache max lookups during batch resolution
- **Rationale**: Multiple prescriptions may reference the same max
- **Priority**: Should
- **Acceptance Criteria**: Single max lookup per (user, lift, type) combination per batch
- **Dependencies**: REQ-PRSC-009

### Prescription CRUD

#### REQ-PRSC-011: Prescription CRUD Operations
- **Description**: The system shall provide API endpoints for prescription management
- **Rationale**: Required for program configuration
- **Priority**: Must
- **Acceptance Criteria**:
  - GET /prescriptions - list prescriptions with filtering
  - GET /prescriptions/{id} - get single prescription
  - POST /prescriptions - create prescription
  - PUT /prescriptions/{id} - update prescription
  - DELETE /prescriptions/{id} - delete prescription
- **Dependencies**: REQ-PRSC-001 through REQ-PRSC-007

## 4. Non-Functional Requirements

### Performance
- **NFR-001**: Single prescription resolution shall complete in < 100ms (p95)
- **NFR-002**: Batch resolution of 20 prescriptions shall complete in < 500ms (p95)
- **NFR-003**: CRUD operations shall complete in < 100ms (p95)

### Reliability
- **NFR-004**: Resolution shall fail gracefully if max not found (clear error message)
- **NFR-005**: Batch resolution shall return partial results on partial failure

### Maintainability
- **NFR-006**: LoadStrategy and SetScheme shall use strategy pattern for extensibility
- **NFR-007**: New strategies/schemes shall be addable without modifying existing code

## 5. External Interfaces

### Data Formats

#### Prescription API Response
```json
{
  "id": "uuid",
  "liftId": "uuid",
  "loadStrategy": {
    "type": "PERCENT_OF",
    "referenceType": "TRAINING_MAX",
    "percentage": 85,
    "roundingIncrement": 5.0,
    "roundingDirection": "NEAREST"
  },
  "setScheme": {
    "type": "FIXED",
    "sets": 5,
    "reps": 5
  },
  "order": 1,
  "notes": "Focus on depth",
  "restSeconds": 180,
  "createdAt": "2024-01-15T10:00:00Z",
  "updatedAt": "2024-01-15T10:00:00Z"
}
```

#### Resolved Prescription Response
```json
{
  "prescriptionId": "uuid",
  "lift": {
    "id": "uuid",
    "name": "Squat",
    "slug": "squat"
  },
  "sets": [
    {
      "setNumber": 1,
      "weight": 265.0,
      "targetReps": 5,
      "isWorkSet": true
    },
    {
      "setNumber": 2,
      "weight": 265.0,
      "targetReps": 5,
      "isWorkSet": true
    }
  ],
  "notes": "Focus on depth",
  "restSeconds": 180
}
```

#### Ramp SetScheme Example
```json
{
  "type": "RAMP",
  "steps": [
    { "percentage": 50, "reps": 5 },
    { "percentage": 63, "reps": 5 },
    { "percentage": 75, "reps": 5 },
    { "percentage": 88, "reps": 5 },
    { "percentage": 100, "reps": 5 }
  ],
  "workSetThreshold": 80
}
```

## 6. Constraints & Assumptions

### Technical Constraints
- LoadStrategy and SetScheme stored as JSONB for flexibility
- Strategy pattern implemented via discriminated unions in TypeScript
- Weight stored as decimal with 0.25 precision

### Assumptions
- Lift and LiftMax entities exist (ERD-001)
- User entity exists with authentication
- Default rounding increment configurable at system level

## 7. Acceptance Criteria & Verification Methods

| Requirement | Verification Method | Success Criteria |
|-------------|---------------------|------------------|
| REQ-PRSC-001 to REQ-PRSC-011 | Unit tests, Integration tests | All prescription operations work correctly |
| REQ-LOAD-001 to REQ-LOAD-006 | Unit tests | PercentOf calculates correctly for all cases |
| REQ-SET-001 to REQ-SET-005 | Unit tests | Fixed and Ramp generate correct sets |
| REQ-PRSC-008, REQ-PRSC-009 | Integration tests | Resolution produces expected output |
| NFR-001, NFR-002 | Load testing | Response times within limits |

## 8. Prioritization / Roadmap

### Must-Have (This ERD)
- Prescription entity with all associations
- PercentOf LoadStrategy with rounding
- Fixed SetScheme
- Ramp SetScheme
- Prescription resolution

### Should-Have (This ERD)
- Batch resolution with caching
- Notes and rest period fields
- Rounding direction option

### Won't-Have (Future ERDs)
- AMRAP SetScheme (Phase 2)
- RPETarget LoadStrategy (Phase 6)
- TopBackoff SetScheme (future)
- RelativeTo LoadStrategy (Phase 7)

## 9. Traceability

### Links to PRD
- PRD-002: Prescription System

### Links to Phase Document
- Phase 001: Core Foundation - Theme 2 (Prescription System)

### Dependencies
- ERD-001: Core Domain Entities (Lift, LiftMax)

### Forward Links (to Tickets)
- To be created after ERD approval

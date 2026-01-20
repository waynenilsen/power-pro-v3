# ERD 004: Progression Rules

## 1. Introduction / Overview

### Purpose
This document specifies the engineering requirements for the Progression Rules system - the logic that mutates LiftMax values over time. This system encodes how lifters systematically increase training loads.

### Scope
- Progression interface and strategy pattern
- Trigger system for controlling when progressions fire
- LinearProgression for session/week-based advancement
- CycleProgression for cycle-end advancement
- Integration with LiftMax mutation and Schedule triggers
- Does NOT include: AMRAP, failure-based, or RPE-based progressions

### Definitions/Glossary
| Term | Definition |
|------|------------|
| Progression | A rule that mutates LiftMax values based on triggers |
| Trigger | An event that causes a progression to evaluate/apply |
| LinearProgression | Add fixed increment at regular intervals |
| CycleProgression | Add fixed increment at cycle completion |
| Increment | The weight added during progression (e.g., 5lb, 2.5kg) |
| Frequency | How often a progression fires (per-session, per-week, per-cycle) |

### Stakeholders
- Engineering Team: Implementation
- Product Owner: Feature validation
- API Consumers: Progression configuration and tracking

## 2. Business / Stakeholder Needs

### Why
Progression rules are the "engine" of powerlifting programs. Without systematic load increases, training stagnates. Different programs use different progression strategies - linear for beginners, cycle-based for intermediate/advanced. A flexible progression system enables representing all programs through configuration.

### Constraints
- Must integrate with LiftMax entity for mutations
- Must integrate with Schedule system for trigger events
- Progressions must be auditable (track what changed and when)

### Success Criteria
- Starting Strength: +5lb squat/deadlift per session, +2.5lb bench/press per session
- Bill Starr 5x5: +5lb per week on main lifts
- 5/3/1: +5lb upper, +10lb lower at end of each 4-week cycle
- Sheiko: No auto-progression (manual)
- Greg Nuckols HF: +5lb at end of 3-week cycle

## 3. Functional Requirements

### Progression Interface

#### REQ-PROG-001: Progression Entity
- **Description**: The system shall define a Progression entity with common interface
- **Rationale**: Enables polymorphic progression handling
- **Priority**: Must
- **Acceptance Criteria**:
  - Base entity with id, name, type discriminator
  - Type-specific parameters stored as JSONB
  - All progressions implement common interface
- **Dependencies**: None

#### REQ-PROG-002: Progression Apply Method
- **Description**: The system shall define a standard apply method for progressions
- **Rationale**: Uniform way to execute progressions regardless of type
- **Priority**: Must
- **Acceptance Criteria**:
  - Interface: `apply(context) -> ProgressionResult`
  - Context includes: userId, liftId, maxType, currentValue, triggerEvent
  - Result includes: newValue, delta, applied (boolean)
- **Dependencies**: REQ-PROG-001

#### REQ-PROG-003: Progression Configuration
- **Description**: The system shall support associating progressions with lifts/programs
- **Rationale**: Different lifts may have different progression rules
- **Priority**: Must
- **Acceptance Criteria**:
  - ProgramProgression join entity: programId, liftId (optional), progressionId, priority
  - Lift-specific progressions override program-level defaults
  - Priority determines order of evaluation
- **Dependencies**: REQ-PROG-001, ERD-003 (Program entity)

### Trigger System

#### REQ-TRIG-001: Trigger Types
- **Description**: The system shall support multiple trigger types
- **Rationale**: Progressions fire at different times for different programs
- **Priority**: Must
- **Acceptance Criteria**:
  - Enum: AFTER_SESSION, AFTER_WEEK, AFTER_CYCLE
  - Each progression specifies its trigger type
- **Dependencies**: None

#### REQ-TRIG-002: Session Trigger
- **Description**: The system shall fire AFTER_SESSION triggers after workout completion
- **Rationale**: Linear progression per-session (Starting Strength)
- **Priority**: Must
- **Acceptance Criteria**:
  - Fired when user completes a training day
  - Passes session context (lifts performed, day, week)
  - Only fires for lifts in the completed session
- **Dependencies**: REQ-TRIG-001

#### REQ-TRIG-003: Week Trigger
- **Description**: The system shall fire AFTER_WEEK triggers after week completion
- **Rationale**: Linear progression per-week (Bill Starr)
- **Priority**: Must
- **Acceptance Criteria**:
  - Fired when user advances from week N to week N+1
  - Passes week context (week number, lifts performed)
  - Fires for all lifts in the progression configuration
- **Dependencies**: REQ-TRIG-001, ERD-003 (state advancement)

#### REQ-TRIG-004: Cycle Trigger
- **Description**: The system shall fire AFTER_CYCLE triggers after cycle completion
- **Rationale**: Cycle progression (5/3/1, Greg Nuckols)
- **Priority**: Must
- **Acceptance Criteria**:
  - Fired when user advances from week N to week 1 (cycle wrap)
  - Passes cycle context (cycle iteration, all lifts)
  - Fires for all lifts in the progression configuration
- **Dependencies**: REQ-TRIG-001, ERD-003 (cycle completion detection)

#### REQ-TRIG-005: Trigger Idempotency
- **Description**: The system shall ensure progressions are not applied multiple times
- **Rationale**: Prevents double-incrementing on retry/error
- **Priority**: Must
- **Acceptance Criteria**:
  - ProgressionLog entity tracks all applications
  - Before applying, check if already applied for this trigger instance
  - Unique constraint on (userId, progressionId, triggerType, triggerTimestamp)
- **Dependencies**: REQ-TRIG-001

### LinearProgression

#### REQ-LINEAR-001: LinearProgression Entity
- **Description**: The system shall implement LinearProgression strategy
- **Rationale**: Most common progression for beginners/intermediates
- **Priority**: Must
- **Acceptance Criteria**:
  - Extends Progression entity
  - Parameters: increment (decimal), triggerType (enum), maxType (1RM or TM)
- **Dependencies**: REQ-PROG-001

#### REQ-LINEAR-002: LinearProgression Apply
- **Description**: The system shall add fixed increment when LinearProgression fires
- **Rationale**: Core progression logic
- **Priority**: Must
- **Acceptance Criteria**:
  - Fetches current LiftMax for (userId, liftId, maxType)
  - Creates new LiftMax with value = current + increment
  - Returns ProgressionResult with delta = increment
- **Dependencies**: REQ-LINEAR-001, ERD-001 (LiftMax)

#### REQ-LINEAR-003: Per-Session Linear
- **Description**: The system shall support per-session linear progression
- **Rationale**: Starting Strength adds weight each session
- **Priority**: Must
- **Acceptance Criteria**:
  - LinearProgression with triggerType = AFTER_SESSION
  - Fires after each session containing the configured lift
- **Dependencies**: REQ-LINEAR-001, REQ-TRIG-002

#### REQ-LINEAR-004: Per-Week Linear
- **Description**: The system shall support per-week linear progression
- **Rationale**: Bill Starr adds weight each week
- **Priority**: Must
- **Acceptance Criteria**:
  - LinearProgression with triggerType = AFTER_WEEK
  - Fires once per week for configured lifts
- **Dependencies**: REQ-LINEAR-001, REQ-TRIG-003

### CycleProgression

#### REQ-CYCLE-001: CycleProgression Entity
- **Description**: The system shall implement CycleProgression strategy
- **Rationale**: Periodized programs progress at cycle boundaries
- **Priority**: Must
- **Acceptance Criteria**:
  - Extends Progression entity
  - Parameters: increment (decimal), maxType (1RM or TM)
  - Implicit triggerType = AFTER_CYCLE
- **Dependencies**: REQ-PROG-001

#### REQ-CYCLE-002: CycleProgression Apply
- **Description**: The system shall add increment when cycle completes
- **Rationale**: 5/3/1 adds 5/10lb after each 4-week cycle
- **Priority**: Must
- **Acceptance Criteria**:
  - Triggered by AFTER_CYCLE event
  - Creates new LiftMax with value = current + increment
  - Works with any cycle length (3-week, 4-week, etc.)
- **Dependencies**: REQ-CYCLE-001, REQ-TRIG-004

#### REQ-CYCLE-003: Lift-Specific Increments
- **Description**: The system shall support different increments per lift
- **Rationale**: 5/3/1: +5lb upper body, +10lb lower body
- **Priority**: Must
- **Acceptance Criteria**:
  - ProgramProgression can have lift-specific increment override
  - If not specified, uses progression default increment
- **Dependencies**: REQ-CYCLE-001, REQ-PROG-003

### Progression History

#### REQ-HIST-001: Progression Log
- **Description**: The system shall log all progression applications
- **Rationale**: Audit trail and debugging
- **Priority**: Must
- **Acceptance Criteria**:
  - ProgressionLog entity: id, userId, progressionId, liftId, previousValue, newValue, delta, appliedAt, triggerType, triggerContext
  - Created on every successful progression application
- **Dependencies**: REQ-PROG-002

#### REQ-HIST-002: Progression History Query
- **Description**: The system shall provide API to query progression history
- **Rationale**: Users want to see their progress over time
- **Priority**: Should
- **Acceptance Criteria**:
  - Endpoint: GET /users/{userId}/progression-history
  - Filter by lift, date range, progression type
  - Returns chronological list of progression applications
- **Dependencies**: REQ-HIST-001

### Manual Progression

#### REQ-MANUAL-001: Manual Max Update
- **Description**: The system shall support manual LiftMax updates
- **Rationale**: Users may test new maxes or need corrections
- **Priority**: Must
- **Acceptance Criteria**:
  - Existing LiftMax CRUD from ERD-001
  - Manual updates logged separately (not as progression)
  - Can override auto-progression values
- **Dependencies**: ERD-001 (LiftMax CRUD)

#### REQ-MANUAL-002: Disable Progression
- **Description**: The system shall support disabling progression for specific lifts
- **Rationale**: Some programs (Sheiko) don't auto-progress
- **Priority**: Should
- **Acceptance Criteria**:
  - ProgramProgression can be marked as disabled
  - Disabled progressions don't fire even when trigger occurs
- **Dependencies**: REQ-PROG-003

### API Endpoints

#### REQ-API-001: Progression CRUD
- **Description**: The system shall provide API endpoints for progression management
- **Rationale**: Required for program configuration
- **Priority**: Must
- **Acceptance Criteria**:
  - GET /progressions - list progressions
  - GET /progressions/{id} - get progression
  - POST /progressions - create progression
  - PUT /progressions/{id} - update progression
  - DELETE /progressions/{id} - delete progression
- **Dependencies**: REQ-PROG-001

#### REQ-API-002: Program Progression Configuration
- **Description**: The system shall provide API for configuring program progressions
- **Rationale**: Link progressions to programs and lifts
- **Priority**: Must
- **Acceptance Criteria**:
  - GET /programs/{id}/progressions - list program's progressions
  - POST /programs/{id}/progressions - add progression to program
  - PUT /programs/{id}/progressions/{progId} - update configuration
  - DELETE /programs/{id}/progressions/{progId} - remove progression
- **Dependencies**: REQ-PROG-003

#### REQ-API-003: Trigger Progression Manually
- **Description**: The system shall support manually triggering progressions
- **Rationale**: Testing and administrative override
- **Priority**: Should
- **Acceptance Criteria**:
  - Endpoint: POST /users/{userId}/progressions/trigger
  - Body: { progressionId, liftId (optional), force: boolean }
  - force=true bypasses idempotency check
- **Dependencies**: REQ-PROG-002, REQ-TRIG-005

## 4. Non-Functional Requirements

### Performance
- **NFR-001**: Progression evaluation shall complete in < 100ms (p95)
- **NFR-002**: Batch progression (multiple lifts) shall complete in < 500ms (p95)

### Reliability
- **NFR-003**: Progression application shall be atomic (no partial updates)
- **NFR-004**: Idempotency shall prevent double-application (100% reliability)
- **NFR-005**: Trigger events shall not be lost (reliable event delivery)

### Maintainability
- **NFR-006**: New progression types shall be addable via strategy pattern
- **NFR-007**: Trigger types shall be extensible without modifying existing code

## 5. External Interfaces

### Data Formats

#### Progression API Response
```json
{
  "id": "uuid",
  "name": "5/3/1 Lower Body Progression",
  "type": "CYCLE_PROGRESSION",
  "parameters": {
    "increment": 10.0,
    "maxType": "TRAINING_MAX"
  },
  "createdAt": "2024-01-15T10:00:00Z",
  "updatedAt": "2024-01-15T10:00:00Z"
}
```

#### LinearProgression Parameters
```json
{
  "increment": 5.0,
  "maxType": "TRAINING_MAX",
  "triggerType": "AFTER_SESSION"
}
```

#### CycleProgression Parameters
```json
{
  "increment": 5.0,
  "maxType": "TRAINING_MAX"
}
```

#### ProgramProgression Configuration
```json
{
  "programId": "uuid",
  "progressionId": "uuid",
  "liftId": "uuid",
  "priority": 1,
  "enabled": true,
  "overrideIncrement": 10.0
}
```

#### ProgressionResult
```json
{
  "applied": true,
  "previousValue": 300.0,
  "newValue": 305.0,
  "delta": 5.0,
  "liftId": "uuid",
  "maxType": "TRAINING_MAX",
  "appliedAt": "2024-01-15T10:00:00Z"
}
```

#### ProgressionLog Entry
```json
{
  "id": "uuid",
  "userId": "uuid",
  "progressionId": "uuid",
  "liftId": "uuid",
  "previousValue": 300.0,
  "newValue": 305.0,
  "delta": 5.0,
  "triggerType": "AFTER_SESSION",
  "triggerContext": {
    "sessionId": "uuid",
    "daySlug": "day-a",
    "weekNumber": 2
  },
  "appliedAt": "2024-01-15T10:00:00Z"
}
```

## 6. Constraints & Assumptions

### Technical Constraints
- Progressions stored as discriminated unions (type + JSONB parameters)
- Trigger events delivered via internal event bus (not external queue in Phase 1)
- Transaction required for LiftMax update + ProgressionLog creation

### Assumptions
- LiftMax entity exists (ERD-001)
- Schedule system provides trigger events (ERD-003)
- Single-user operations (no concurrent progression for same user/lift)

## 7. Acceptance Criteria & Verification Methods

| Requirement | Verification Method | Success Criteria |
|-------------|---------------------|------------------|
| REQ-PROG-* | Unit tests | Progression interface and strategies work correctly |
| REQ-TRIG-* | Integration tests | Triggers fire at correct times |
| REQ-LINEAR-* | Unit tests, Integration tests | Linear progression applies correctly |
| REQ-CYCLE-* | Integration tests | Cycle progression applies at cycle end |
| REQ-HIST-* | Integration tests | Progression history logged and queryable |
| NFR-004 | Retry tests | Double-apply prevented |

## 8. Prioritization / Roadmap

### Must-Have (This ERD)
- Progression interface and entity
- LinearProgression (session, week triggers)
- CycleProgression (cycle trigger)
- Progression logging
- Trigger system integration

### Should-Have (This ERD)
- Manual trigger endpoint
- Progression history query
- Disable progression option

### Won't-Have (Future ERDs)
- AMRAPProgression (Phase 2)
- DeloadOnFailure (Phase 3)
- StageProgression (Phase 3)
- DoubleProgression (Phase 4)

## 9. Traceability

### Links to PRD
- PRD-004: Progression Rules

### Links to Phase Document
- Phase 001: Core Foundation - Theme 3 (Progression Rules)

### Dependencies
- ERD-001: Core Domain Entities (LiftMax to mutate)
- ERD-003: Schedule & Periodization (trigger events)

### Forward Links (to Tickets)
- To be created after ERD approval

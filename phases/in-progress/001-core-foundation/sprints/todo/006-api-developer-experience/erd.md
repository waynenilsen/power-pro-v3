# ERD 006: API & Developer Experience

## 1. Introduction / Overview

### Purpose
This document specifies the engineering requirements for the API & Developer Experience sprint - ensuring PowerPro's API is well-documented, consistent, and demonstrates all 5 target programs.

### Scope
- API documentation for all endpoints
- API consistency audit and improvements
- Program configuration examples for 5 target programs
- Does NOT include: new endpoints, authentication changes, client SDKs, frontend

### Definitions/Glossary
| Term | Definition |
|------|------------|
| OpenAPI | API specification standard (formerly Swagger) |
| RESTful | Architectural style for APIs using HTTP methods semantically |
| Program Configuration | The complete setup of entities to represent a training program |

### Stakeholders
- Engineering Team: Implementation
- Product Owner: Feature validation
- External Developers: API consumers

## 2. Business / Stakeholder Needs

### Why
PowerPro is a headless API. The API design and documentation directly impact adoption and usability. External developers need clear, consistent, well-documented endpoints to integrate successfully.

### Constraints
- Must document existing API without breaking changes
- Must follow REST best practices
- Must demonstrate all 5 Phase 1 target programs

### Success Criteria
- All API endpoints have documentation
- API follows consistent patterns
- All 5 programs can be configured and demonstrated

## 3. Functional Requirements

### API Documentation

#### REQ-DOC-001: Endpoint Documentation
- **Description**: The system shall have documentation for all API endpoints
- **Rationale**: Developers need to know what endpoints exist and how to use them
- **Priority**: Must
- **Acceptance Criteria**:
  - Every endpoint documented with HTTP method, path, description
  - Request body schema documented where applicable
  - Response schema documented for success and error cases
- **Dependencies**: None

#### REQ-DOC-002: Example Requests
- **Description**: The system shall provide example requests for all endpoints
- **Rationale**: Examples help developers understand expected formats
- **Priority**: Must
- **Acceptance Criteria**:
  - Each endpoint has at least one example request
  - Examples use realistic data
  - Examples are copy-paste ready
- **Dependencies**: REQ-DOC-001

#### REQ-DOC-003: Example Responses
- **Description**: The system shall provide example responses for all endpoints
- **Rationale**: Developers need to know what to expect from the API
- **Priority**: Must
- **Acceptance Criteria**:
  - Each endpoint has example success response
  - Common error responses documented
  - Response formats match actual API behavior
- **Dependencies**: REQ-DOC-001

#### REQ-DOC-004: Error Documentation
- **Description**: The system shall document all error codes and formats
- **Rationale**: Proper error handling requires understanding error responses
- **Priority**: Must
- **Acceptance Criteria**:
  - All HTTP status codes used are documented
  - Error response format is documented
  - Common error scenarios are explained
- **Dependencies**: REQ-DOC-001

#### REQ-DOC-005: Workflow Documentation
- **Description**: The system shall document common API workflows
- **Rationale**: Multi-step workflows need guidance beyond individual endpoints
- **Priority**: Should
- **Acceptance Criteria**:
  - User onboarding workflow documented
  - Program enrollment workflow documented
  - Workout generation workflow documented
  - Progression trigger workflow documented
- **Dependencies**: REQ-DOC-001

### API Consistency

#### REQ-API-001: RESTful Conventions Audit
- **Description**: The system shall follow RESTful conventions for all endpoints
- **Rationale**: Consistent conventions make the API predictable
- **Priority**: Must
- **Acceptance Criteria**:
  - GET for reads, POST for creates, PUT for updates, DELETE for deletes
  - Resource-oriented URLs (nouns, not verbs)
  - Proper HTTP status codes (200, 201, 400, 401, 404, etc.)
- **Dependencies**: None

#### REQ-API-002: Response Envelope
- **Description**: The system shall use consistent response envelope format
- **Rationale**: Predictable response structure simplifies client code
- **Priority**: Should
- **Acceptance Criteria**:
  - Success responses follow consistent structure
  - Error responses follow consistent structure
  - Pagination follows consistent pattern
- **Dependencies**: None

#### REQ-API-003: Pagination Pattern
- **Description**: The system shall use consistent pagination for list endpoints
- **Rationale**: Different pagination patterns confuse developers
- **Priority**: Should
- **Acceptance Criteria**:
  - All list endpoints support pagination
  - Pagination parameters are consistent (limit, offset or cursor)
  - Response includes total count and pagination metadata
- **Dependencies**: None

#### REQ-API-004: Filtering Pattern
- **Description**: The system shall use consistent filtering for list endpoints
- **Rationale**: Consistent filtering makes queries predictable
- **Priority**: Should
- **Acceptance Criteria**:
  - Filter parameters follow consistent naming
  - Filter behavior is documented
  - Complex filters use consistent syntax
- **Dependencies**: None

### Program Configuration Examples

#### REQ-PROG-001: Starting Strength Configuration
- **Description**: The system shall document Starting Strength program configuration
- **Rationale**: Validates linear per-session progression
- **Priority**: Must
- **Acceptance Criteria**:
  - Complete configuration documented (lifts, prescriptions, days, cycle, progression)
  - Example workout output shown
  - Progression behavior demonstrated
- **Dependencies**: ERD-001 through ERD-004

#### REQ-PROG-002: Bill Starr 5x5 Configuration
- **Description**: The system shall document Bill Starr 5x5 program configuration
- **Rationale**: Validates ramping sets and Heavy/Light/Medium days
- **Priority**: Must
- **Acceptance Criteria**:
  - Complete configuration documented
  - Daily intensity variation demonstrated
  - Ramping set scheme shown
- **Dependencies**: ERD-001 through ERD-004

#### REQ-PROG-003: Wendler 5/3/1 BBB Configuration
- **Description**: The system shall document Wendler 5/3/1 BBB program configuration
- **Rationale**: Validates weekly percentage variation and cycle progression
- **Priority**: Must
- **Acceptance Criteria**:
  - Complete configuration documented
  - 4-week cycle structure shown
  - Weekly percentage lookup demonstrated
  - Cycle-end progression shown
- **Dependencies**: ERD-001 through ERD-004

#### REQ-PROG-004: Sheiko Beginner Configuration
- **Description**: The system shall document Sheiko Beginner program configuration
- **Rationale**: Validates high-volume programming with no auto-progression
- **Priority**: Must
- **Acceptance Criteria**:
  - Complete configuration documented
  - Multiple daily sessions shown if applicable
  - Manual progression approach documented
- **Dependencies**: ERD-001 through ERD-004

#### REQ-PROG-005: Greg Nuckols HF Configuration
- **Description**: The system shall document Greg Nuckols High Frequency program configuration
- **Rationale**: Validates daily undulation and 3-week cycles
- **Priority**: Must
- **Acceptance Criteria**:
  - Complete configuration documented
  - Daily intensity variation demonstrated
  - 3-week cycle structure shown
- **Dependencies**: ERD-001 through ERD-004

## 4. Non-Functional Requirements

### Usability
- **NFR-001**: Documentation shall be clear and accessible to developers unfamiliar with PowerPro
- **NFR-002**: Examples shall be copy-paste ready

### Maintainability
- **NFR-003**: Documentation shall be maintainable (preferably generated from code/schemas)
- **NFR-004**: Program examples shall use real API calls that can be verified

### Accuracy
- **NFR-005**: Documentation shall accurately reflect current API behavior
- **NFR-006**: Examples shall produce expected outputs when executed

## 5. External Interfaces

### Documentation Format
- OpenAPI/Swagger specification for machine-readable docs
- Markdown documentation for human-readable guides
- Example JSON for request/response bodies

### Program Configuration Format
```json
{
  "program": {
    "name": "Starting Strength",
    "slug": "starting-strength",
    "cycleId": "uuid",
    "progressions": [...]
  },
  "cycle": {
    "lengthWeeks": 1,
    "weeks": [...]
  },
  "days": [...],
  "prescriptions": [...],
  "lifts": [...],
  "lookups": [...]
}
```

## 6. Constraints & Assumptions

### Technical Constraints
- No breaking API changes
- Documentation must match implementation
- Examples must be executable

### Assumptions
- All Phase 1 entities are implemented and working
- Authentication system is in place
- API endpoints are stable

## 7. Acceptance Criteria & Verification Methods

| Requirement | Verification Method | Success Criteria |
|-------------|---------------------|------------------|
| REQ-DOC-001 to REQ-DOC-005 | Documentation review | All endpoints documented with examples |
| REQ-API-001 to REQ-API-004 | API audit | All endpoints follow conventions |
| REQ-PROG-001 to REQ-PROG-005 | Integration test | All programs can be configured and generate workouts |

## 8. Prioritization / Roadmap

### Must-Have (This ERD)
- All endpoints documented
- Example requests and responses
- Error documentation
- RESTful conventions audit
- All 5 program configurations

### Should-Have (This ERD)
- Workflow documentation
- Consistent response envelope
- Pagination and filtering patterns

### Won't-Have (Future)
- Client SDKs
- Interactive API explorer
- Video tutorials

## 9. Traceability

### Links to PRD
- PRD-006: API & Developer Experience

### Links to Phase Document
- Phase 001: Core Foundation - Theme 5 (API & Developer Experience)

### Dependencies
- ERD-001: Core Domain Entities
- ERD-002: Prescription System
- ERD-003: Schedule & Periodization
- ERD-004: Progression Rules
- ERD-005: Technical Debt (recommended but not required)

### Forward Links (to Tickets)
- To be created after ERD approval

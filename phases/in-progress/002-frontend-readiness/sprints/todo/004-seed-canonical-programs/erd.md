# ERD 004: Seed Canonical Programs

## 1. Introduction / Overview

### Purpose
This document specifies the engineering requirements for seeding canonical powerlifting programs into the PowerPro database. These programs provide users with immediate access to proven training methodologies without requiring manual configuration.

### Scope
- Database migrations to seed four canonical programs
- Verification tests to validate seeded program structures
- Documentation of canonical program slugs
- Does NOT include: warmup prescriptions, accessory work, program variants, user customization

### Definitions/Glossary
| Term | Definition |
|------|------------|
| Canonical Program | A pre-seeded, system-owned program that users can enroll in but not modify |
| Seed Migration | A goose migration that inserts pre-defined data rather than schema changes |
| Lookup | Reference data for lift types, rep schemes, intensity modifiers |
| Prescription | A specific set/rep/intensity configuration for an exercise |
| Progression Model | The algorithm for increasing weight session-to-session or week-to-week |

### Stakeholders
- Engineering Team: Implementation
- Product Owner: Feature validation
- API Consumers (Frontend Teams): Program discovery and enrollment

## 2. Business / Stakeholder Needs

### Why
New users need immediate access to training programs. Requiring users to manually configure a program before they can start training creates friction and requires expertise most beginners don't have. Canonical programs solve this by providing ready-to-use, professionally-designed training methodologies.

### Constraints
- SQLite only - migrations must use SQLite-compatible SQL
- Go only - any supporting code must be in Go
- Must use existing schema - no schema changes allowed in this sprint
- Must preserve idempotency - migrations should be safe to re-run
- System user required - programs need an author; use system/admin user

### Success Criteria
- Four canonical programs exist in database after migration
- Each program has correct structure (days, weeks, cycles) per documentation
- Each program has accurate prescriptions (sets, reps, percentages)
- Verification tests confirm program accuracy
- Programs are discoverable via existing program list endpoint

## 3. Functional Requirements

### Program Seeding

#### REQ-SEED-001: Starting Strength Program
- **Description**: The system shall seed the Starting Strength novice program
- **Rationale**: Starting Strength is the most popular beginner program; 3 days/week, linear progression
- **Priority**: Must
- **Acceptance Criteria**:
  - Program created with slug `starting-strength`
  - Program has name "Starting Strength"
  - Program has description matching program methodology
  - Program has 2 workout days (A/B rotation)
  - Workout A: Squat 3x5, Bench Press 3x5, Deadlift 1x5
  - Workout B: Squat 3x5, Overhead Press 3x5, Power Clean 5x3
  - Linear progression model: +5 lbs per session (squat, bench, press), +10 lbs (deadlift)
  - Program marked as canonical (is_canonical flag or similar)
- **Dependencies**: None

#### REQ-SEED-002: Texas Method Program
- **Description**: The system shall seed the Texas Method intermediate program
- **Rationale**: Texas Method is a popular weekly periodization program for intermediates
- **Priority**: Must
- **Acceptance Criteria**:
  - Program created with slug `texas-method`
  - Program has name "Texas Method"
  - Program has 3 workout days per week (Volume/Recovery/Intensity)
  - Volume Day (Monday): Squat 5x5 @ 90%, Bench/Press 5x5 @ 90%, Deadlift 1x5
  - Recovery Day (Wednesday): Squat 2x5 @ 80%, Bench/Press 3x5 @ 90%
  - Intensity Day (Friday): Squat 1x5 @ 100%, Bench/Press 1x5 @ 100%
  - Weekly progression model: +5 lbs per week (upper), +5 lbs per week (lower)
  - Bench and Press alternate each week
- **Dependencies**: REQ-SEED-001

#### REQ-SEED-003: Wendler 5/3/1 Program
- **Description**: The system shall seed the Wendler 5/3/1 program
- **Rationale**: 5/3/1 is the most popular percentage-based monthly periodization program
- **Priority**: Must
- **Acceptance Criteria**:
  - Program created with slug `531`
  - Program has name "Wendler 5/3/1"
  - Program has 4 workout days per week (one main lift per day)
  - Day 1: Overhead Press, Day 2: Deadlift, Day 3: Bench Press, Day 4: Squat
  - 4-week cycle structure (5s week, 3s week, 5/3/1 week, deload week)
  - Week 1: 65%x5, 75%x5, 85%x5+
  - Week 2: 70%x3, 80%x3, 90%x3+
  - Week 3: 75%x5, 85%x3, 95%x1+
  - Week 4 (deload): 40%x5, 50%x5, 60%x5
  - Cycle progression: +5 lbs upper body, +10 lbs lower body per cycle
  - Percentages based on Training Max (90% of 1RM)
- **Dependencies**: REQ-SEED-001

#### REQ-SEED-004: GZCLP Program
- **Description**: The system shall seed the GZCLP linear progression program
- **Rationale**: GZCLP is a popular tiered progression system for beginners/intermediates
- **Priority**: Must
- **Acceptance Criteria**:
  - Program created with slug `gzclp`
  - Program has name "GZCLP"
  - Program has 4 workout days per week
  - Day 1: T1 Squat (5x3+), T2 Bench (3x10)
  - Day 2: T1 OHP (5x3+), T2 Deadlift (3x10)
  - Day 3: T1 Bench (5x3+), T2 Squat (3x10)
  - Day 4: T1 Deadlift (5x3+), T2 OHP (3x10)
  - T1 progression: 5x3+ -> 6x2+ -> 10x1+ (on failure)
  - T2 progression: 3x10 -> 3x8 -> 3x6 (on failure)
  - Linear progression: +5 lbs lower, +2.5 lbs upper per successful session
- **Dependencies**: REQ-SEED-001

### Lookup Data

#### REQ-LOOKUP-001: Lift Type Lookups
- **Description**: The system shall have lookup entries for all lifts used in canonical programs
- **Rationale**: Prescriptions reference lift types; all referenced lifts must exist
- **Priority**: Must
- **Acceptance Criteria**:
  - Squat lift type exists
  - Bench Press lift type exists
  - Deadlift lift type exists
  - Overhead Press lift type exists
  - Power Clean lift type exists
  - Front Squat lift type exists (for program variants)
- **Dependencies**: None

#### REQ-LOOKUP-002: Rep Scheme Lookups
- **Description**: The system shall have lookup entries for all rep schemes used
- **Rationale**: Prescriptions reference rep schemes; all referenced schemes must exist
- **Priority**: Must
- **Acceptance Criteria**:
  - 3x5, 1x5, 5x3, 5x5, 2x5, 1x3, 3x10, 3x8, 3x6, 5x3+, 6x2+, 10x1+ schemes exist
  - AMRAP variants properly marked
- **Dependencies**: None

### Verification

#### REQ-VERIFY-001: Program Structure Tests
- **Description**: The system shall have tests verifying seeded program structures
- **Rationale**: Ensures migrations created correct program topology
- **Priority**: Must
- **Acceptance Criteria**:
  - Test verifies Starting Strength has 2 days, correct exercises per day
  - Test verifies Texas Method has 3 days, correct exercises per day
  - Test verifies 5/3/1 has 4 days, 4 weeks per cycle
  - Test verifies GZCLP has 4 days, correct T1/T2 pairings
- **Dependencies**: REQ-SEED-001 through REQ-SEED-004

#### REQ-VERIFY-002: Prescription Accuracy Tests
- **Description**: The system shall have tests verifying prescription accuracy
- **Rationale**: Ensures sets, reps, and percentages match program documentation
- **Priority**: Must
- **Acceptance Criteria**:
  - Tests verify correct set/rep counts for each exercise
  - Tests verify correct intensity percentages for percentage-based programs
  - Tests verify correct progression increments
- **Dependencies**: REQ-VERIFY-001

### Documentation

#### REQ-DOC-001: Canonical Slug Documentation
- **Description**: The system shall document canonical program slugs
- **Rationale**: Frontend teams need to know which slugs to use for program discovery
- **Priority**: Should
- **Acceptance Criteria**:
  - README or API doc lists canonical slugs: `starting-strength`, `texas-method`, `531`, `gzclp`
  - Documentation explains that canonical programs cannot be modified by users
  - Documentation explains enrollment process for canonical programs
- **Dependencies**: REQ-SEED-001 through REQ-SEED-004

## 4. Non-Functional Requirements

### Performance
- **NFR-001**: Seed migration shall complete in < 5 seconds
- **NFR-002**: Seeded programs shall be queryable with same performance as user programs

### Reliability
- **NFR-003**: Seed migrations shall be idempotent (safe to run multiple times)
- **NFR-004**: Seed migrations shall use transactions to ensure atomicity

### Maintainability
- **NFR-005**: Each program seed shall be in a separate migration file for clarity
- **NFR-006**: Migration SQL shall be commented to explain program structure
- **NFR-007**: Prescription data shall match programs/*.md documentation

## 5. External Interfaces

### System Interfaces
- Database: SQLite with goose migrations
- Existing tables: programs, program_days, program_weeks, program_cycles, prescriptions, lookups

### Data Formats

#### Program Query Response (existing)
```json
{
  "id": "uuid",
  "slug": "starting-strength",
  "name": "Starting Strength",
  "description": "A novice strength program...",
  "isCanonical": true,
  "createdAt": "2024-01-15T10:00:00Z",
  "updatedAt": "2024-01-15T10:00:00Z"
}
```

#### Canonical Program Slugs
```
starting-strength - Starting Strength novice program
texas-method - Texas Method intermediate program
531 - Wendler 5/3/1 program
gzclp - GZCLP linear progression program
```

## 6. Constraints & Assumptions

### Technical Constraints
- Go backend
- SQLite database
- goose migrations
- No schema changes in this sprint
- Existing program entity model must accommodate canonical programs

### Assumptions
- System/admin user exists for program authorship
- Lookup tables have standard lift types and rep schemes
- Existing program query endpoints work for canonical programs
- is_canonical or similar flag exists or can be inferred (e.g., author is system user)

## 7. Acceptance Criteria & Verification Methods

| Requirement | Verification Method | Success Criteria |
|-------------|---------------------|------------------|
| REQ-SEED-001 to REQ-SEED-004 | Integration tests | Programs exist, structure matches spec |
| REQ-LOOKUP-001, REQ-LOOKUP-002 | Migration tests | Lookups exist before program seeds |
| REQ-VERIFY-001, REQ-VERIFY-002 | Automated tests | All verification tests pass |
| REQ-DOC-001 | Manual review | Documentation exists and is accurate |
| NFR-001, NFR-002 | Performance tests | Migrations fast, queries performant |
| NFR-003, NFR-004 | Migration tests | Idempotent, atomic |

## 8. Prioritization / Roadmap

### Must-Have (This ERD)
- Starting Strength program seed
- Texas Method program seed
- Wendler 5/3/1 program seed
- GZCLP program seed
- Required lookup data
- Structure verification tests
- Prescription accuracy tests

### Should-Have (This ERD)
- Canonical slug documentation
- Detailed program descriptions

### Won't-Have (Future ERDs)
- Warmup prescriptions
- Accessory work prescriptions
- Program variants (5/3/1 BBB, 5/3/1 FSL, etc.)
- User program templates
- Program customization

## 9. Traceability

### Links to PRD
- PRD-004: Seed Canonical Programs

### Links to Phase Document
- Phase 002: Frontend Readiness - Theme 2 (Program Discovery)

### Forward Links (to Tickets)
- 001-starting-strength-seed.md (REQ-SEED-001, REQ-LOOKUP-001, REQ-LOOKUP-002)
- 002-texas-method-seed.md (REQ-SEED-002)
- 003-531-seed.md (REQ-SEED-003)
- 004-gzclp-seed.md (REQ-SEED-004)
- 005-program-verification-tests.md (REQ-VERIFY-001, REQ-VERIFY-002)
- 006-canonical-programs-documentation.md (REQ-DOC-001)

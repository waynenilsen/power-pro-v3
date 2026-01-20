# Phase 001: Core Foundation

## Vision & Strategic Objectives

### Product Vision
PowerPro is a headless API for powerlifting programming that enables lifters to follow structured training programs with intelligent load calculation, set schemes, and progression rules. The platform abstracts the complexity of diverse powerlifting methodologies into a unified, DRY domain model that can represent any program through composable primitives.

### Strategic Objectives
1. **Establish Core Domain Model**: Build the foundational entities (Lift, LiftMax, Prescription, LoadStrategy, SetScheme, Progression, Schedule) that all future features will build upon
2. **Enable 5 Programs**: Unlock Starting Strength, Bill Starr 5x5, Wendler 5/3/1 BBB, Sheiko Beginner, and Greg Nuckols High Frequency programs
3. **Validate Architecture**: Prove the DRY architecture can represent fundamentally different programming philosophies without code duplication
4. **Developer Experience**: Create a clean, well-documented API that external developers can integrate with

### Context
This is the foundational phase of PowerPro. The core domain objects defined here will be extended and composed in all future phases. Getting the abstractions right is critical—premature optimization or over-engineering will create technical debt, but insufficient abstraction will lead to spaghetti code as more programs are added. The 5 programs unlocked in this phase represent diverse programming philosophies (linear progression, daily undulation, weekly periodization, high volume) which validates the domain model's flexibility.

## Themes & Initiatives

### Theme 1: Core Domain Entities
- **Strategic Objective**: Establish Core Domain Model
- **Rationale**: The Lift and LiftMax entities are the atomic building blocks. Every prescription, progression, and schedule references these. Getting these abstractions right enables all downstream features.
- **Initiatives**:
  - Initiative A: Implement Lift entity with support for main lifts (Squat, Bench, Deadlift) and extensibility for variations
  - Initiative B: Implement LiftMax entity supporting 1RM and TM (Training Max) types with validation and conversion logic
- **Dependencies**: None—this is the foundation
- **Risks**: Over-engineering the Lift entity with too many fields/relationships before understanding real usage patterns

### Theme 2: Prescription System
- **Strategic Objective**: Establish Core Domain Model, Enable 5 Programs
- **Rationale**: Prescriptions combine what lift to do, how to calculate the weight, and how sets/reps are structured. This is the core "what to do today" abstraction.
- **Initiatives**:
  - Initiative A: Implement Prescription entity linking Lift, LoadStrategy, and SetScheme
  - Initiative B: Implement `PercentOf` LoadStrategy (e.g., 85% of TM)
  - Initiative C: Implement `Fixed` SetScheme (e.g., 5x5) and `Ramp` SetScheme (e.g., warmup percentages)
  - Initiative D: Implement configurable weight rounding (to nearest 2.5lb, 5lb, etc.)
- **Dependencies**: Theme 1 (Lift, LiftMax)
- **Risks**: SetScheme abstraction may need revision when AMRAP and other dynamic schemes are added in Phase 2

### Theme 3: Progression Rules
- **Strategic Objective**: Enable 5 Programs, Validate Architecture
- **Rationale**: Progression rules mutate LiftMax values over time. Linear progression (add weight each session/week) and cycle progression (add weight at end of cycle) cover the 5 target programs.
- **Initiatives**:
  - Initiative A: Implement `LinearProgression` (configurable increment, frequency: per-session/per-week)
  - Initiative B: Implement `CycleProgression` (increment applied at cycle completion)
- **Dependencies**: Theme 1 (LiftMax to mutate), Theme 4 (Schedule for timing)
- **Risks**: Progression trigger timing logic may be complex; ensure clean separation between "when" and "what"

### Theme 4: Schedule & Periodization
- **Strategic Objective**: Enable 5 Programs, Validate Architecture
- **Rationale**: Programs structure training across days, weeks, and cycles. This theme implements the temporal organization that ties prescriptions to calendar time.
- **Initiatives**:
  - Initiative A: Implement Day entity (named training day with exercise slots)
  - Initiative B: Implement Week entity (collection of days with A/B rotation support)
  - Initiative C: Implement Cycle entity (repeating unit: 1-week, 3-week, 4-week, etc.)
  - Initiative D: Implement `WeeklyLookup` (week 1 = X%, week 2 = Y%) and `DailyLookup` (day-specific percentages)
- **Dependencies**: Theme 2 (Prescriptions to schedule)
- **Risks**: Rotation logic (A/B days, lift focus changes) may require more abstraction than initially apparent

### Theme 5: API & Developer Experience
- **Strategic Objective**: Developer Experience
- **Rationale**: PowerPro is a headless API. The API design directly impacts adoption and usability.
- **Initiatives**:
  - Initiative A: Define RESTful API endpoints for program configuration and workout generation
  - Initiative B: Implement authentication system for API access
  - Initiative C: Create API documentation with examples for each of the 5 unlocked programs
- **Dependencies**: Themes 1-4 (domain model must exist to expose via API)
- **Risks**: API design may need revision as more features are added; design for extensibility

## Timeline

| Phase | Timeline | Themes/Initiatives |
|-------|----------|-------------------|
| Now | Current | Theme 1: Core Domain Entities (Lift, LiftMax) |
| Now | Current | Theme 2: Prescription System (Prescription, PercentOf, Fixed, Ramp, Rounding) |
| Next | Following | Theme 3: Progression Rules (LinearProgression, CycleProgression) |
| Next | Following | Theme 4: Schedule & Periodization (Day, Week, Cycle, Lookups) |
| Later | After Next | Theme 5: API & Developer Experience (Endpoints, Auth, Documentation) |

## Success Metrics

| Theme | Metric | Target |
|-------|--------|--------|
| Core Domain Entities | Lift and LiftMax entities implemented | Complete with tests |
| Prescription System | Prescription generation accuracy | 100% match to expected outputs for test cases |
| Progression Rules | LinearProgression and CycleProgression working | Correctly mutate LiftMax values per program specs |
| Schedule & Periodization | Full program representation | Can represent all 5 target programs without special cases |
| API & Developer Experience | API endpoint coverage | All core CRUD operations exposed |
| Overall | Programs unlocked | 5 (Starting Strength, Bill Starr 5x5, 5/3/1 BBB, Sheiko Beginner, Greg Nuckols HF) |

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Over-engineering domain model | Medium | High | Start minimal, add abstraction only when needed by real programs |
| Under-abstraction leading to code duplication | Medium | High | Validate against all 5 target programs before finalizing |
| SetScheme abstraction insufficient for Phase 2 | Low | Medium | Design with extension points; don't close off AMRAP addition |
| API design requires major revision | Low | Medium | Version API from start; design for backward compatibility |

## Review & Update Process

- **Review cadence**: After each theme completion
- **Owner**: Engineering Lead
- **Approval**: Product Owner reviews theme completion, Engineering Lead approves technical implementation
- **Update triggers**:
  - Discovery of new requirements from program analysis
  - Technical constraints requiring architecture changes
  - Completion of themes enabling next phase readiness assessment

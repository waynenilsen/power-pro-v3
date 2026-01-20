# PRD 001: Core Domain Entities

## Product Vision

PowerPro is a headless API for powerlifting programming that enables lifters to follow structured training programs with intelligent load calculation, set schemes, and progression rules. This PRD establishes the foundational domain entities that all other features will build upon.

## Strategic Objectives

1. **Establish Core Domain Model**: Build the foundational Lift and LiftMax entities that all future features will reference
2. **Enable Extensibility**: Design entities to support diverse programming philosophies without code duplication
3. **Validate Architecture**: Prove the DRY architecture can represent fundamentally different lifts and max types

## Themes & Initiatives

### Theme 1: Lift Entity
- **Strategic Objective**: Establish Core Domain Model
- **Rationale**: The Lift entity is the atomic building block. Every prescription, progression, and schedule references a lift. Getting this abstraction right enables all downstream features.
- **Initiatives**:
  - Initiative A: Implement Lift entity with support for main lifts (Squat, Bench, Deadlift)
  - Initiative B: Support lift variations (Pause Squat, Close-Grip Bench, etc.)
  - Initiative C: Define lift categories/tags for grouping (competition lifts, accessories, etc.)

### Theme 2: LiftMax Entity
- **Strategic Objective**: Establish Core Domain Model, Enable Extensibility
- **Rationale**: LiftMax stores the reference numbers used for load calculation. Different programs use different reference types (1RM, Training Max, etc.). A flexible LiftMax system enables diverse programming without special cases.
- **Initiatives**:
  - Initiative A: Implement LiftMax entity supporting 1RM and TM (Training Max) types
  - Initiative B: Implement validation logic (TM typically 85-90% of 1RM)
  - Initiative C: Implement conversion logic between max types
  - Initiative D: Support user-specific LiftMax values with timestamp tracking

## Success Metrics

| Metric | Target |
|--------|--------|
| Lift entity implemented | Complete with full test coverage |
| LiftMax entity implemented | Complete with full test coverage |
| Max type support | 1RM and TM types working correctly |
| Conversion accuracy | 100% match to expected conversion formulas |
| API response time | < 100ms for entity CRUD operations |

## Timeline

| Phase | Scope |
|-------|-------|
| Now | Lift entity with main lifts and variations |
| Now | LiftMax entity with 1RM and TM support |
| Next | Lift categories and tagging |
| Next | Advanced max types (xRM, E1RM) prepared for Phase 2+ |

## Dependencies

- None - this is the foundational layer

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Over-engineering Lift entity | Medium | High | Start minimal, add fields only when needed by real programs |
| Insufficient LiftMax types | Low | Medium | Design with extension points for future max types |
| Performance issues with max lookups | Low | Low | Index properly, benchmark early |

## Out of Scope

- E1RM (Estimated 1RM) calculations - deferred to Phase 2+
- xRM (Rep Max) support - deferred to Phase 2+
- RPE-based max estimation - deferred to Phase 6
- Historical max tracking over time - future enhancement

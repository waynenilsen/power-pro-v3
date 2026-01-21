# PRD 004: Progression Rules

## Product Vision

PowerPro is a headless API for powerlifting programming. This PRD establishes the Progression Rules system - the logic that mutates LiftMax values over time. Progression rules are what make programs "work" - they encode how lifters get stronger by systematically increasing their training loads.

## Strategic Objectives

1. **Enable 5 Programs**: Implement LinearProgression and CycleProgression sufficient for all Phase 1 target programs
2. **Validate Architecture**: Prove that diverse progression strategies can be expressed through a common interface
3. **Clean Separation**: Maintain clear separation between "when" (triggers) and "what" (progression logic)

## Themes & Initiatives

### Theme 1: Progression Interface
- **Strategic Objective**: Validate Architecture
- **Rationale**: A common interface enables diverse progression strategies to be composed and applied uniformly.
- **Initiatives**:
  - Initiative A: Define Progression interface with apply() method
  - Initiative B: Define Trigger system for controlling when progressions fire
  - Initiative C: Support progression chaining and prioritization

### Theme 2: LinearProgression
- **Strategic Objective**: Enable 5 Programs
- **Rationale**: Linear progression (add weight each session/week) is the most common beginner/intermediate strategy. Starting Strength, Bill Starr 5x5, and Greg Nuckols HF all use linear progression.
- **Initiatives**:
  - Initiative A: Implement LinearProgression with configurable increment
  - Initiative B: Support per-session, per-week, and per-cycle frequencies
  - Initiative C: Support lift-specific increment configuration

### Theme 3: CycleProgression
- **Strategic Objective**: Enable 5 Programs
- **Rationale**: Cycle progression (add weight at end of cycle) is used by 5/3/1 and other periodized programs. This enables longer adaptation windows.
- **Initiatives**:
  - Initiative A: Implement CycleProgression triggered on cycle completion
  - Initiative B: Support configurable cycle-end increment
  - Initiative C: Integrate with cycle completion detection from Schedule system

## Success Metrics

| Metric | Target |
|--------|--------|
| Progression interface implemented | Clean, extensible design |
| LinearProgression working | Correctly increments per session/week |
| CycleProgression working | Correctly increments at cycle end |
| All 5 programs representable | Progression rules match program specs |
| No missed progressions | 100% reliability in progression triggers |

## Timeline

| Phase | Scope |
|-------|-------|
| Now | Progression interface and Trigger system |
| Now | LinearProgression with frequency options |
| Now | CycleProgression with cycle integration |
| Later | AMRAPProgression (Phase 2) |
| Later | DeloadOnFailure, StageProgression (Phase 3) |

## Dependencies

- PRD-001/ERD-001: Core Domain Entities (LiftMax to mutate)
- PRD-003/ERD-003: Schedule & Periodization (triggers for cycle completion)

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Trigger timing complexity | Medium | Medium | Clean separation between scheduling and progression |
| Progression conflicts | Low | Low | Clear priority/order for multiple progressions |
| Missed progression events | Low | High | Reliable event system with idempotency |

## Out of Scope

- AMRAPProgression - deferred to Phase 2
- DeloadOnFailure - deferred to Phase 3
- StageProgression - deferred to Phase 3
- DoubleProgression - deferred to Phase 4
- RPE-based progression - deferred to Phase 6

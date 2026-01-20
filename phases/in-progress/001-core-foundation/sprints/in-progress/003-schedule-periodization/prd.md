# PRD 003: Schedule & Periodization

## Product Vision

PowerPro is a headless API for powerlifting programming. This PRD establishes the Schedule & Periodization system - the temporal organization that ties prescriptions to calendar time. Programs structure training across days, weeks, and cycles, and this system enables that organization.

## Strategic Objectives

1. **Enable 5 Programs**: Implement scheduling primitives sufficient for Starting Strength, Bill Starr 5x5, 5/3/1 BBB, Sheiko Beginner, and Greg Nuckols HF
2. **Validate Architecture**: Prove that diverse periodization schemes can be expressed through composable primitives
3. **Support Lookup Tables**: Implement weekly and daily lookup tables that vary prescriptions based on cycle position

## Themes & Initiatives

### Theme 1: Day Entity
- **Strategic Objective**: Enable 5 Programs
- **Rationale**: A Day is a named training session with ordered exercise slots. All programs organize work by training days.
- **Initiatives**:
  - Initiative A: Implement Day entity with name and prescription slots
  - Initiative B: Support A/B day naming for alternating programs (Starting Strength)
  - Initiative C: Support day-specific metadata (Heavy/Light/Medium for Bill Starr)

### Theme 2: Week Entity
- **Strategic Objective**: Enable 5 Programs
- **Rationale**: A Week is a collection of training days. Most programs operate on weekly cycles.
- **Initiatives**:
  - Initiative A: Implement Week entity as a collection of Days
  - Initiative B: Support week numbering within cycles (Week 1, 2, 3, etc.)
  - Initiative C: Support A/B week rotation for alternating programs

### Theme 3: Cycle Entity
- **Strategic Objective**: Enable 5 Programs, Validate Architecture
- **Rationale**: A Cycle is the repeating unit of a program (1-week, 3-week, 4-week, etc.). 5/3/1 uses 4-week cycles, Greg Nuckols uses 3-week cycles.
- **Initiatives**:
  - Initiative A: Implement Cycle entity with configurable length
  - Initiative B: Support cycle position tracking for progression triggers
  - Initiative C: Support cycle completion detection

### Theme 4: Lookup Tables
- **Strategic Objective**: Enable 5 Programs
- **Rationale**: Programs like 5/3/1 and Greg Nuckols vary percentages by week or day. Lookup tables enable this variation without duplicating prescriptions.
- **Initiatives**:
  - Initiative A: Implement WeeklyLookup (Week 1 = X%, Week 2 = Y%)
  - Initiative B: Implement DailyLookup (Monday = X%, Tuesday = Y%)
  - Initiative C: Integrate lookups with LoadStrategy resolution

## Success Metrics

| Metric | Target |
|--------|--------|
| Day entity implemented | Complete with test coverage |
| Week entity implemented | Complete with test coverage |
| Cycle entity implemented | Complete with test coverage |
| WeeklyLookup working | Correct percentage variation by week |
| DailyLookup working | Correct percentage variation by day |
| All 5 programs representable | No special cases required |

## Timeline

| Phase | Scope |
|-------|-------|
| Now | Day entity with prescription slots |
| Now | Week entity with day collection |
| Now | Cycle entity with length configuration |
| Now | WeeklyLookup and DailyLookup tables |
| Later | Rotation logic for lift focus changes (Phase 5) |

## Dependencies

- PRD-001/ERD-001: Core Domain Entities (Lift, LiftMax)
- PRD-002/ERD-002: Prescription System (prescriptions to schedule)

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Rotation logic more complex than expected | Medium | Medium | Defer rotation to Phase 5; validate simple A/B works |
| Week/Cycle position tracking complexity | Low | Medium | Start with simple counter; add state machine if needed |
| Lookup table integration with LoadStrategy | Low | Low | Design clean interface; test extensively |

## Out of Scope

- Lift rotation (which lift is primary this week) - deferred to Phase 5
- Phase blocks with different parameters - deferred to Phase 5
- DaysOut countdown scheduling - deferred to Phase 10
- Taper protocols - deferred to Phase 10

# PRD 002: Prescription System

## Product Vision

PowerPro is a headless API for powerlifting programming. This PRD establishes the Prescription system - the core abstraction that tells a lifter "what to do" for a single exercise slot. Prescriptions combine a lift, a load calculation strategy, and a set/rep scheme into a unified instruction.

## Strategic Objectives

1. **Establish Core Domain Model**: Build the Prescription entity that links Lift, LoadStrategy, and SetScheme
2. **Enable 5 Programs**: Implement LoadStrategy and SetScheme variants sufficient for all Phase 1 target programs
3. **Validate Architecture**: Prove that diverse programming styles can be expressed through composable primitives

## Themes & Initiatives

### Theme 1: Prescription Entity
- **Strategic Objective**: Establish Core Domain Model
- **Rationale**: The Prescription is the fundamental unit of programming - it answers "what lift, at what weight, for how many sets and reps". All programs are composed of prescriptions.
- **Initiatives**:
  - Initiative A: Implement Prescription entity linking Lift, LoadStrategy, and SetScheme
  - Initiative B: Support prescription ordering within a training day
  - Initiative C: Support optional prescription metadata (notes, cues, rest periods)

### Theme 2: LoadStrategy - PercentOf
- **Strategic Objective**: Enable 5 Programs
- **Rationale**: The vast majority of programs prescribe load as a percentage of a reference max (TM or 1RM). PercentOf is the foundational load strategy.
- **Initiatives**:
  - Initiative A: Implement PercentOf LoadStrategy (e.g., 85% of TM)
  - Initiative B: Support configurable reference max type (1RM or TM)
  - Initiative C: Implement weight rounding to configurable increments (2.5lb, 5lb, 2.5kg, etc.)

### Theme 3: SetScheme - Fixed and Ramp
- **Strategic Objective**: Enable 5 Programs, Validate Architecture
- **Rationale**: Fixed sets (5x5) and ramping sets (warmup progression) cover the majority of Phase 1 programs. These two schemes demonstrate the abstraction works.
- **Initiatives**:
  - Initiative A: Implement Fixed SetScheme (e.g., 5 sets of 5 reps)
  - Initiative B: Implement Ramp SetScheme (e.g., 50%x5, 63%x5, 75%x5, 88%x5, 100%x5)
  - Initiative C: Support set-specific metadata (target RPE, notes)

## Success Metrics

| Metric | Target |
|--------|--------|
| Prescription entity implemented | Complete with full test coverage |
| PercentOf LoadStrategy working | Correct weight calculation for all test cases |
| Fixed SetScheme working | Correctly generates set list |
| Ramp SetScheme working | Correctly generates progressive warmup sets |
| Weight rounding | Accurate to specified increment in all cases |

## Timeline

| Phase | Scope |
|-------|-------|
| Now | Prescription entity with Lift reference |
| Now | PercentOf LoadStrategy with rounding |
| Now | Fixed SetScheme |
| Now | Ramp SetScheme |
| Later | Additional LoadStrategies (RPETarget, etc.) in future phases |

## Dependencies

- PRD-001/ERD-001: Core Domain Entities (Lift, LiftMax) must be complete
- User entity must exist for LiftMax lookup

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| SetScheme abstraction insufficient for AMRAP | Low | Medium | Design with extension points; verify AMRAP can be added later |
| LoadStrategy interface too rigid | Low | Medium | Use strategy pattern with common interface |
| Rounding edge cases | Low | Low | Comprehensive test coverage for edge cases |

## Out of Scope

- AMRAP SetScheme - deferred to Phase 2
- RPETarget LoadStrategy - deferred to Phase 6
- TopBackoff SetScheme - future enhancement
- MRS/FatigueDrop SetSchemes - deferred to Phase 8
- RelativeTo LoadStrategy - deferred to Phase 7

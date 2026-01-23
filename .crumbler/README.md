# PowerPro Development Roadmap

## Product Vision

PowerPro is a headless API for powerlifting programming that enables lifters to follow structured training programs with intelligent load calculation, set schemes, and progression rules. The platform abstracts the complexity of diverse powerlifting methodologies into a unified, DRY domain model that can represent any program through composable primitives.

## Completed Work

**Phase 1: Core Foundation** - Core domain model established
- ✅ Core Domain Entities (Lift, LiftMax)
- ✅ Prescription System (Prescription, PercentOf LoadStrategy, Fixed/Ramp SetSchemes)
- ✅ Schedule & Periodization (Day, Week, Cycle, WeeklyLookup, DailyLookup)
- ✅ Progression Rules (LinearProgression, CycleProgression)

**Programs Unlocked:** 5 (Starting Strength, Bill Starr 5x5, Wendler 5/3/1 BBB, Sheiko Beginner, Greg Nuckols High Frequency)

## Remaining Phase 1 Work

### Technical Debt Cleanup
Address code quality, test coverage, and documentation alignment from Phase 1 rapid development.

### API & Developer Experience
Complete API documentation, ensure RESTful consistency, and provide program configuration examples for all 5 unlocked programs.

### E2E Tests for Phase 1 Programs
**Requirement:** Each phase must include end-to-end (e2e) tests for all programs unlocked in that phase. These tests validate that the program can be fully configured and executed through the API.

**Phase 1 programs requiring e2e tests:**
- [Starting Strength](./programs/004-starting-strength.md)
- [Bill Starr 5x5](./programs/015-bill-starr-5x5.md)
- [Wendler 5/3/1 BBB](./programs/003-wendler-531-bbb.md)
- [Sheiko Beginner](./programs/010-sheiko-beginner.md)
- [Greg Nuckols High Frequency](./programs/020-greg-nuckols-frequency.md)

## Future Phases

### Phase 2: AMRAP
Add AMRAP SetScheme, LoggedSet tracking, AMRAPProgression, and AfterSet triggers. Unlocks 4 programs (Greyskull LP, Reddit PPL, Greg Nuckols Beginner, nSuns 5/3/1 LP).

**E2E tests required for:**
- [Greyskull LP](./programs/008-greyskull-lp.md)
- [Reddit PPL 6-Day](./programs/005-reddit-ppl-6-day.md)
- [Greg Nuckols Beginner](./programs/012-greg-nuckols-beginner.md)
- [nSuns 5/3/1 LP 5-Day](./programs/001-nsuns-531-lp-5-day.md)

### Phase 3: Failure Handling
Add failure tracking, DeloadOnFailure and StageProgression rules, OnFailure triggers. Unlocks 2 programs (GZCLP, Texas Method).

**E2E tests required for:**
- [GZCLP](./programs/002-gzclp-linear-progression.md)
- [Texas Method](./programs/006-texas-method.md)

### Phase 4: Double Progression
Add DoubleProgression rule and RepRange SetScheme. Formalizes existing mechanics.

**E2E tests:** No new programs unlocked (formalization only).

### Phase 5: Rotation
Add Rotation schedule support for exercise slot rotation and cycle-based lift focus changes. Unlocks 2 programs (nSuns CAP3, Inverted Juggernaut 5/3/1).

**E2E tests required for:**
- [nSuns CAP3](./programs/013-nsuns-cap3.md)
- [Inverted Juggernaut 5/3/1](./programs/014-inverted-juggernaut-531.md)

### Phase 6: RPE
Add RPETarget LoadStrategy, RPEChart lookup, and RPE field in LoggedSet. Unlocks 1 partial program (RTS Intermediate).

**E2E tests required for:**
- [RTS Intermediate](./programs/019-rts-intermediate.md) (partial)

### Phase 7: E1RM
Add E1RM LiftMax type, FindRM and RelativeTo LoadStrategies, E1RM calculation from LoggedSet, and prescription chaining. Unlocks 3 programs (GZCL Jacked & Tan 2.0, Calgary Barbell 8-Week, Calgary Barbell 16-Week).

**E2E tests required for:**
- [GZCL Jacked & Tan 2.0](./programs/009-gzcl-jacked-and-tan-2.md)
- [Calgary Barbell 8-Week](./programs/018-calgary-barbell-8-week.md)
- [Calgary Barbell 16-Week](./programs/011-calgary-barbell-16-week.md)

### Phase 8: Fatigue Protocols
Add FatigueDrop and MRS SetSchemes, variable set count, and AfterSet triggers with RPE conditions. Unlocks 2 programs (GZCL Compendium, RTS Intermediate complete).

**E2E tests required for:**
- [GZCL Compendium](./programs/016-gzcl-compendium.md)
- [RTS Intermediate](./programs/019-rts-intermediate.md) (complete)

### Phase 9: Volume Targets
Add TotalReps SetScheme, session-level rep tracking, and superset notation. Unlocks 1 program (5/3/1 Building the Monolith).

**E2E tests required for:**
- [5/3/1 Building the Monolith](./programs/007-531-building-the-monolith.md)

### Phase 10: Peaking
Add DaysOut schedule, phase transitions based on calendar, opener practice sessions, and taper protocols. Unlocks 1 program (Sheiko Intermediate).

**E2E tests required for:**
- [Sheiko Intermediate](./programs/017-sheiko-intermediate.md)

## Strategic Objectives

1. **Complete Phase 1**: Finish technical debt cleanup and API developer experience work
2. **Enable 20 Programs**: Progressively unlock programs through phases 2-10
3. **Maintain DRY Architecture**: Ensure all features build on composable primitives without code duplication
4. **Developer Experience**: Maintain clean, well-documented API throughout all phases
5. **Test Coverage**: Ensure all unlocked programs have comprehensive e2e tests validating full program execution

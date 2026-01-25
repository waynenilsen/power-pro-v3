# PRD 004: Seed Canonical Programs

## Product Vision

PowerPro needs pre-configured, real-world powerlifting programs available to users immediately upon registration. Without canonical programs, users would have to manually configure complex program structures before they can start training. Seeding popular programs eliminates this friction and provides users with proven, ready-to-use training methodologies.

## Strategic Objectives

1. **Immediate Value**: New users can enroll in a program within seconds of creating an account
2. **Real-World Programs**: Provide actual powerlifting programs that match established methodologies
3. **Reference Implementations**: Canonical programs serve as templates for user-created programs
4. **Test Data Quality**: Move from synthetic test data to production-quality program structures

## Themes & Initiatives

### Theme 1: Beginner Programs
- **Strategic Objective**: Immediate Value
- **Rationale**: New lifters need simple, proven programs with linear progression that they can follow without deep understanding of periodization.
- **Initiatives**:
  - Initiative A: Seed Starting Strength program (3 days/week, linear progression)
  - Initiative B: Seed GZCLP program (3-4 days/week, tiered progression)

### Theme 2: Intermediate Programs
- **Strategic Objective**: Real-World Programs
- **Rationale**: Lifters who've exhausted linear progression need weekly or monthly periodization to continue making progress.
- **Initiatives**:
  - Initiative A: Seed Texas Method program (3 days/week, weekly periodization)
  - Initiative B: Seed Wendler 5/3/1 program (4 days/week, monthly periodization)

### Theme 3: Program Completeness
- **Strategic Objective**: Reference Implementations
- **Rationale**: Each seeded program must be complete and accurate, demonstrating proper use of the program entity relationships.
- **Initiatives**:
  - Initiative A: Include all program days, weeks, and cycles per program specification
  - Initiative B: Include accurate warmup protocols per program documentation
  - Initiative C: Include correct progression models per program methodology

### Theme 4: Quality Assurance
- **Strategic Objective**: Test Data Quality
- **Rationale**: Seeded programs must be verified for accuracy against program documentation.
- **Initiatives**:
  - Initiative A: Create verification tests that validate program structure
  - Initiative B: Create tests that validate prescription percentages and rep schemes
  - Initiative C: Document canonical program slugs for API consumers

## Success Metrics

| Metric | Target |
|--------|--------|
| Starting Strength program seeded and queryable | Complete |
| Texas Method program seeded and queryable | Complete |
| Wendler 5/3/1 program seeded and queryable | Complete |
| GZCLP program seeded and queryable | Complete |
| All seeded programs have correct day/week/cycle structure | Complete |
| All seeded programs have accurate prescription percentages | Complete |
| Verification tests pass for all seeded programs | Complete |
| Canonical slugs documented in API documentation | Complete |

## Timeline

| Phase | Scope |
|-------|-------|
| Now | Migration with Starting Strength program seed |
| Now | Migration with Texas Method program seed |
| Now | Migration with 5/3/1 program seed |
| Now | Migration with GZCLP program seed |
| Now | Verification tests for all seeded programs |

## Dependencies

- Sprint 001 (Authentication System) - Users must exist to be program authors
- Core domain entities from Phase 001 - Programs, prescriptions, lookups tables must exist

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Program specification inaccuracy | Medium | Medium | Reference programs/*.md documentation; create detailed verification tests |
| Migration complexity | Medium | Low | Use existing E2E test patterns as reference for INSERT structure |
| Slug conflicts with user programs | Low | Low | Reserve canonical slugs; use namespace prefixes if needed |
| Large migration file | Low | Low | Consider one migration per program if review becomes difficult |

## Out of Scope

- User-created program templates - future enhancement
- Program variant selection (e.g., 5/3/1 BBB vs 5/3/1 FSL) - seed one variant per program
- Warmup prescription details - focus on working sets; warmups can be added later
- Accessory work prescriptions - focus on main lifts only
- Program customization options - seed fixed configurations

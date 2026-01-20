# Software Roadmap

Guidelines for creating effective software roadmaps based on best practices from product leaders, engineering teams, and industry experts.

## Directory Structure

Roadmap phase documents are organized in a nested structure with sprints and tickets:

```
phases/
  {state}/                    # todo, in-progress, done, not-doing
    NNN-phase-name/
      NNN-phase-name.md       # Phase document
      sprints/
        {state}/              # todo, in-progress, done, not-doing
          NNN-sprint-name/
            prd.md            # Product Requirements Document
            erd.md            # Engineering Requirements Document
            tickets/
              {state}/        # todo, in-progress, done, not-doing
                NNN-ticket-name.md
```

## Roadmap Document: README.md

### Requirement: `README.md` is the Roadmap Document

The root-level `README.md` file serves as the roadmap document for the entire project.

### Purpose
- Provides a high-level overview of all phase documents
- Links to individual phase documents
- Shares the overarching product vision that guides all phases
- Serves as the entry point for understanding the project roadmap

### Content Requirements

The README.md should include:

1. **Product Vision**
   - The long-term vision for the product
   - The "why" behind the product
   - Strategic direction that guides all phases

2. **Phase Document Overview**
   - Brief description of what phase documents are
   - How they relate to the overall roadmap
   - Link to this roadmap guidelines document (`prompts/roadmap.md`)

3. **Phase Document Links**
   - Links to all phase documents organized by status (todo, in-progress, done)
   - Brief one-sentence description of each phase document
   - Phase documents should be listed in numerical order

4. **Roadmap Structure**
   - Explanation of how phases relate to ERDs and tickets
   - Reference to the traceability chain (Phase → ERD → Tickets)

### Format Example

```markdown
# [Product Name]

## Product Vision

[Overarching vision statement that guides all phases]

## Roadmap

### In Progress
- [Phase 001: Foundation](./phases/in-progress/001-foundation.md) - Establishing core infrastructure and authentication
- [Phase 002: Core Features](./phases/in-progress/002-core-features.md) - Building essential user-facing functionality

### Planned
- [Phase 003: Platform Scaling](./phases/todo/003-platform-scaling.md) - Scaling infrastructure for growth

### Completed
- [Phase 000: Initial Setup](./phases/done/000-initial-setup.md) - Project initialization and tooling

## Roadmap Structure

Phase documents break down into ERDs, which in turn break down into tickets. See [roadmap guidelines](./prompts/roadmap.md) for details.
```

### Maintenance
- Update README.md whenever phase documents are created, moved, or completed
- Keep phase descriptions brief (one sentence)
- Ensure all phase documents are linked
- Vision statement should be stable but can evolve as the product matures

## Directory Descriptions

### `phases/`
The primary directory containing all phase directories. Each phase is a directory containing the phase document and nested sprints.

### `phases/{state}/NNN-phase-name/`
Each phase is a directory containing:
- `NNN-phase-name.md` - The phase document itself
- `sprints/` - Directory containing sprints organized by state

### Phase States
- `phases/todo/` - Phase directories that are planned but not yet started
- `phases/in-progress/` - Phase directories currently being actively worked on
- `phases/done/` - Completed phase directories (only if all sprints are done)
- `phases/not-doing/` - Phase directories that are cancelled or deferred

### Sprint States
Within each phase, sprints have their own states:
- `sprints/todo/` - Sprints planned but not started
- `sprints/in-progress/` - Sprints currently being worked on
- `sprints/done/` - Completed sprints (only if all tickets are done)
- `sprints/not-doing/` - Sprints that are cancelled or deferred

## Workflow

1. **Create**: New phase directories start in `phases/todo/`
2. **Start**: Move phase directories to `phases/in-progress/` when work begins
3. **Complete**: Move phase directories to `phases/done/` when all sprints are completed
4. **Cancel**: Move phase directories to `phases/not-doing/` if they won't be completed

**Important**: Phases and sprints have independent states. A phase can be in-progress while its sprints are in various states.

## Phase Document File Format

### File Naming Convention

Phase directories and documents must follow a specific naming format:

**Phase Directory Format**: `NNN-description/`

Where:
- `NNN` is a zero-padded phase number (e.g., `001`, `002`, `010`, `100`)
- `description` is a short, descriptive name using hyphens or underscores
- The directory contains `NNN-description.md` (the phase document)

**Examples**:
- `001-q1-foundation/` containing `001-q1-foundation.md`
- `002-q2-user-authentication/` containing `002-q2-user-authentication.md`
- `010-q3-platform-scaling/` containing `010-q3-platform-scaling.md`

**Why zero-padding?**
Zero-padding ensures that when directories are sorted lexicographically (alphabetically), they are also sorted numerically. Without zero-padding, `10-phase/` would sort before `2-phase/`, which is incorrect.

### File Content

Each phase document should contain the roadmap structure outlined below.

**Important**: Phase documents must NOT contain status information. Status is always implied by the directory location (`todo/`, `in-progress/`, `done/`, or `not-doing/`). Including status in the file body would create redundancy and potential inconsistencies.

## Core Principles of Excellent Roadmaps

| Principle | Description | Why It Matters |
|-----------|-------------|----------------|
| **Vision-Aligned** | Every roadmap connects to a clear product vision and strategic objectives | Provides north star for decision-making and prioritization |
| **Outcome-Oriented** | Focus on themes and outcomes (e.g., "improve onboarding success") rather than just features | Ensures work ties back to business/customer impact |
| **Data-Driven** | Prioritization based on customer feedback, usage data, support tickets, research | Keeps roadmap grounded in reality, not just loudest voices |
| **Flexible & Living** | Regular reviews and updates (monthly/quarterly) with adaptable time horizons | Responds to changing market conditions, tech constraints, priorities |
| **Audience-Aware** | Tailored detail levels for executives, engineers, customers | Ensures each stakeholder gets relevant information |
| **Dependency-Aware** | Explicitly identifies dependencies, resource constraints, risks | Prevents unrealistic commitments and blockers |
| **Measurable** | Clear success metrics and KPIs for each initiative | Enables validation and demonstrates value |

## Roadmap Structure

### 1. Vision & Strategic Objectives
- **Product Vision**: Long-term direction and "why" behind the product
- **Strategic Objectives**: 2-4 measurable goals (e.g., OKRs) tied to business outcomes
- **Context**: Market position, competitive landscape, key trends

### 2. Themes & Key Initiatives
- **Outcome-Oriented Themes**: 2-5 high-level themes (not feature lists)
  - Examples: "Reduce onboarding friction", "Increase user retention", "Improve platform reliability"
- **Theme Rationale**: Why each theme matters and what strategic objective it supports
- **Initiatives/Epics**: 2-4 major initiatives per theme
  - Each initiative should be outcome-focused, not just a feature description

### 3. Timeline & Phasing
- **Time Horizons**: Use flexible buckets:
  - **Now** (current quarter/sprint): More detail, specific features
  - **Next** (next quarter): Moderate detail, epics/initiatives
  - **Later** (beyond next quarter): High-level themes, exploration/discovery
- **Avoid**: Rigid dates for far-future items; use loose time windows instead
- **Format**: Can be quarterly, by release, or "Now/Next/Later" structure

### 4. Prioritization Framework
- **Use frameworks** like RICE, ICE, or Value vs Effort to evaluate items
- **Data sources**:
  - Customer feedback and support tickets
  - Usage analytics and user behavior data
  - Market research and competitive analysis
  - Business metrics (revenue, retention, growth)
- **Document rationale**: Why each item is prioritized as it is

### 5. Ownership & Dependencies
- **Owners**: Assign owners to initiatives, themes, or milestones
- **Dependencies**: 
  - Cross-team dependencies
  - Technical dependencies
  - External dependencies (vendors, partners)
- **Resource Constraints**: Team bandwidth, budget, skills required

### 6. Success Metrics & KPIs
- **For each theme/initiative**: Define 1-2 measurable success criteria
- **Types of metrics**:
  - Leading indicators (early signals of progress)
  - Lagging indicators (final outcomes)
  - Examples: Engagement, performance, revenue, churn, time-to-value
- **Targets**: Specific, measurable targets where applicable

### 7. Risk Assessment & Mitigation
- **Known Risks**: Technical challenges, market shifts, resource constraints
- **Assumptions**: What you're assuming to be true
- **Unknowns**: Areas requiring discovery or validation
- **Mitigation Strategies**: Contingency plans or decision points
- **Risk Flags**: Items that need special attention or early validation

### 8. Review Cadence & Process
- **Update Frequency**: Monthly or quarterly reviews (define schedule)
- **Roles & Responsibilities**:
  - Who owns the roadmap
  - Who reviews and provides input
  - Who approves changes
- **Change Process**: When and how roadmap changes are made
- **Communication Plan**: How roadmap is shared with different audiences

### 9. Technical Debt Sprints

**Rule**: Every 5th sprint (e.g., Sprint-005, Sprint-010, Sprint-015, Sprint-020, etc.) must be a **technical debt paydown sprint**.

#### Purpose
This ensures regular, systematic attention to technical debt at the roadmap level. Technical debt sprints can be very short and focus on high-level technical improvements rather than feature development.

#### What to Include
- **Theme**: Technical debt reduction (e.g., "Improve Code Quality", "Modernize Infrastructure", "Reduce Security Debt")
- **Strategic Objective**: How reducing this debt supports business goals (e.g., faster feature delivery, reduced risk, improved reliability)
- **Key Initiatives**: High-level technical debt items to address
- **Success Metrics**: How to measure debt reduction (e.g., reduced bug rate, faster build times, improved test coverage)

#### Examples
- Sprint-005: Reduce code debt and improve maintainability
- Sprint-010: Modernize infrastructure and update dependencies
- Sprint-015: Improve test coverage and reduce flaky tests
- Sprint-020: Enhance security posture and patch vulnerabilities

#### Sprint Structure
- Each technical debt sprint contains both a PRD (`prd.md`) and ERD (`erd.md`)
- The PRD provides high-level context and business justification
- The ERD provides detailed requirements for addressing the technical debt
- See `erd.md` for technical debt ERD guidelines

## Best Practices

### ✅ Do
- **Focus on outcomes over features**: Themes should describe business impact, not just features
- **Use flexible time horizons**: Avoid rigid dates for items far in the future
- **Tailor to audience**: Create different views for executives, engineers, customers
- **Incorporate data**: Use customer feedback, analytics, and research to inform priorities
- **Show dependencies**: Make blocking relationships explicit
- **Define success metrics**: Every initiative should have measurable outcomes
- **Regular reviews**: Update roadmap monthly or quarterly to keep it current
- **Be transparent**: Acknowledge uncertainty, assumptions, and risks
- **Keep it visual**: Use clear visuals (swimlanes, timelines, milestones) with consistent formatting
- **Maintain clarity**: Avoid clutter; focus on what matters most

### ❌ Don't
- **Don't be a feature factory**: Avoid roadmaps that are just feature lists without narrative or goals
- **Don't over-detail**: Don't include user stories, specs, or resource-level details that belong in ERDs or tickets
- **Don't ignore dependencies**: Account for resource constraints and technical dependencies
- **Don't let it go stale**: Regular updates prevent mismatches between roadmap and reality
- **Don't overcommit**: Avoid rigid deadlines that can't adapt to changing conditions
- **Don't ignore risks**: Surface and address risks proactively
- **Don't skip prioritization**: Use frameworks, not just intuition or loudest voices
- **Don't forget the "why"**: Always connect items back to vision and strategic objectives

## Audience-Specific Views

### Executive View
- **Focus**: Strategic themes, business impact, high-level timing
- **Include**: Vision, strategic objectives, themes, key metrics, major milestones
- **Exclude**: Technical details, implementation specifics, resource allocation

### Engineering View
- **Focus**: Technical initiatives, dependencies, rough durations
- **Include**: Themes, initiatives, technical dependencies, resource needs, constraints
- **Exclude**: Business strategy details, customer-facing messaging

### Customer/Sales View
- **Focus**: Benefits, progress, what's coming
- **Include**: Themes (customer-benefit focused), high-level timeline, value propositions
- **Exclude**: Internal dependencies, technical details, resource constraints

## Integration with Sprint and Ticket System

Phase documents, sprints, and tickets are organized in a nested hierarchy:
- Phase directories are managed in `phases/` with subdirectories: `todo/`, `in-progress/`, `done/`, `not-doing/`
- Sprints are nested under phases in `sprints/` with the same subdirectories
- Tickets are nested under sprints in `tickets/` with the same subdirectories
- All use zero-padded numbering for lexicographical sorting
- All avoid status in file content (status is implied by directory location)

### From Roadmap to Sprint
- Phase document themes and initiatives inform sprint creation
- Sprints should align with phase document priorities
- Each sprint contains both a PRD (`prd.md`) and ERD (`erd.md`)
- **Every 5th sprint must be a technical debt paydown sprint** (see ERD guidelines)

### From Sprint ERD to Tickets
- Sprint ERDs decompose phase document initiatives into detailed requirements
- Tickets implement sprint ERD requirements
- Tickets are created within the sprint's `tickets/todo/` directory
- Maintain traceability: Phase Document → Sprint (PRD + ERD) → Tickets

### Traceability Chain
```
Phase Document Theme/Initiative (phases/{state}/NNN-phase-name/NNN-phase-name.md)
  ↓
Sprint (phases/{state}/NNN-phase-name/sprints/{state}/NNN-sprint-name/)
  ├── prd.md (Product Requirements)
  └── erd.md (Engineering Requirements with REQ-001, REQ-002, ...)
  ↓
Tickets (phases/{state}/NNN-phase-name/sprints/{state}/NNN-sprint-name/tickets/{state}/NNN-ticket-name.md)
```

### Sprint Completion Rules
- A sprint cannot be moved to `done` if it has any tickets in `todo` or `in-progress` states
- All tickets must be completed or moved to `not-doing` before closing a sprint
- Use `./sdlc.sh move sprint <number> done` to move sprints (validation is automatic)

## Phase Document Format

Each phase document follows this structure. Phase documents are stored in the `phases/` directory structure.

### Document Structure
```markdown
# Phase [Number]: [Phase Name] — [Date]

## Vision & Strategic Objectives
[Product vision and 2-4 strategic goals]

## Themes & Initiatives

### Theme 1: [Outcome-Oriented Theme Name]
- **Strategic Objective**: [Which objective this supports]
- **Rationale**: [Why this theme matters]
- **Initiatives**:
  - Initiative A: [Description, owner, success metrics]
  - Initiative B: [Description, owner, success metrics]
- **Dependencies**: [What blocks or enables this]
- **Risks**: [Known risks and mitigation]

[Repeat for each theme]

## Timeline

| Phase | Timeline | Themes/Initiatives |
|-------|----------|-------------------|
| Now | [Current quarter] | [List items] |
| Next | [Next quarter] | [List items] |
| Later | [Future] | [List items] |

## Success Metrics
[KPIs and targets for each theme/initiative]

## Review & Update Process
- Review cadence: [Monthly/Quarterly]
- Owner: [Role/Name]
- Approval: [Who approves changes]
```

## Prioritization Frameworks

### RICE Framework
- **Reach**: How many users/events affected (per time period)
- **Impact**: How much it matters (0.25 = minimal, 0.5 = low, 1 = medium, 2 = high, 3 = massive)
- **Confidence**: How certain you are (50% = low, 80% = medium, 100% = high)
- **Effort**: Person-months of work
- **Score**: (Reach × Impact × Confidence) / Effort

### ICE Framework
- **Impact**: How much it matters (1-10)
- **Confidence**: How certain you are (1-10)
- **Ease**: How easy to implement (1-10)
- **Score**: (Impact + Confidence + Ease) / 3

### Value vs Effort
- Plot items on 2x2 matrix: High/Low Value × High/Low Effort
- Prioritize: High Value / Low Effort first, then High Value / High Effort

## Common Roadmap Types

### Product Roadmap
- Focus: Customer-facing features and improvements
- Audience: Executives, customers, sales, marketing
- Timeframe: Quarterly to annual

### Engineering/Technical Roadmap
- Focus: Technical initiatives, infrastructure, platform improvements
- Audience: Engineering teams, technical leadership
- Timeframe: Quarterly to multi-year

### Release Roadmap
- Focus: Specific releases and what's included
- Audience: Product, engineering, customers
- Timeframe: Sprint to quarterly

## Validation Checklist

Before finalizing a phase document:

- [ ] Vision and strategic objectives are clearly defined
- [ ] Themes are outcome-oriented, not just feature lists
- [ ] Time horizons are flexible (especially for far-future items)
- [ ] Prioritization uses data and frameworks, not just intuition
- [ ] Dependencies and constraints are identified
- [ ] Success metrics are defined for each theme/initiative
- [ ] Risks are acknowledged with mitigation strategies
- [ ] Review cadence and process are established
- [ ] Phase document is tailored for different audiences
- [ ] Phase document aligns with sprint priorities
- [ ] Phase document can be traced to sprints and tickets
- [ ] Phase directory follows naming convention (zero-padded number)
- [ ] Phase document does not include status (status is implied by directory location)
- [ ] Sprints are properly nested under the phase directory

### README.md Roadmap Checklist

Before considering the roadmap system complete:

- [ ] `README.md` exists and serves as the roadmap document
- [ ] Product vision is clearly stated in README.md
- [ ] All phase documents are linked in README.md
- [ ] Phase documents are organized by status (todo, in-progress, done)
- [ ] Each phase document has a brief one-sentence description
- [ ] README.md explains the roadmap structure and traceability chain

# Engineering Requirements Document (ERD)

Guidelines for creating high-quality Engineering Requirements Documents based on industry best practices and standards (IEEE 29148, ISO/IEC).

## Directory Structure

ERD documents follow the same directory structure as tickets:

```
erds/
├── todo/          # ERDs that need to be created or refined
├── in-progress/   # ERDs currently being worked on
├── done/          # Completed ERDs
└── not-doing/     # ERDs that are cancelled or will not be done
```

## Directory Descriptions

### `erds/`
The primary directory containing all ERD subdirectories.

### `erds/todo/`
ERDs that are planned but not yet started. These represent requirements documents that need to be created or refined.

### `erds/in-progress/`
ERDs that are currently being actively worked on. Move ERDs here when work begins.

### `erds/done/`
Completed ERDs. Move ERDs here when they are finished and approved.

### `erds/not-doing/`
ERDs that have been cancelled, deferred indefinitely, or decided against. Use this for ERDs that will not be completed.

## Workflow

1. **Create**: New ERDs start in `erds/todo/`
2. **Start**: Move ERDs to `erds/in-progress/` when work begins
3. **Complete**: Move ERDs to `erds/done/` when finished and approved
4. **Cancel**: Move ERDs to `erds/not-doing/` if they won't be completed

## ERD File Format

### File Naming Convention

ERDs must follow a specific naming format to enable lexicographical sorting by ERD number:

**Format**: `NNN-description.md`

Where:
- `NNN` is a zero-padded ERD number (e.g., `001`, `002`, `010`, `100`)
- `description` is a short, descriptive name using hyphens or underscores
- File extension is `.md`

**Examples**:
- `001-user-authentication-system.md`
- `002-payment-processing.md`
- `010-api-integration.md`
- `100-deployment-infrastructure.md`

**Why zero-padding?**
Zero-padding ensures that when files are sorted lexicographically (alphabetically), they are also sorted numerically. Without zero-padding, `10-erd.md` would sort before `2-erd.md`, which is incorrect.

### File Content

Each ERD file should contain the document structure outlined below.

**Important**: ERD files must NOT contain status information. Status is always implied by the directory location (`todo/`, `in-progress/`, `done/`, or `not-doing/`). Including status in the file body would create redundancy and potential inconsistencies.

## Core Qualities of Excellent Requirements

Every requirement and the document as a whole should possess these qualities:

| Quality | Description | Why It Matters |
|---------|-------------|----------------|
| **Correct** | Accurately reflects stakeholder needs and domain realities | Prevents building the wrong thing or misaligned expectations |
| **Unambiguous** | Each requirement has only one interpretation. Avoid vague terms like "fast", "easy", "user-friendly" | Ensures all readers (developers, testers, stakeholders) understand consistently |
| **Verifiable/Testable** | Can concretely test or demonstrate each requirement is met | Enables QA, minimizes risk of hidden assumptions |
| **Complete** | Includes ALL functional & non-functional requirements, external interfaces, constraints | Gaps cause scope creep; missing pieces lead to costly rework |
| **Consistent** | Uses consistent terms, formatting, style; no contradictory requirements | Makes document easier to read, automate, and verify |
| **Traceable** | Each requirement has unique identifier; linked to source and to test cases/design artifacts | Enables understanding of "why", managing changes, closing gaps |
| **Prioritized** | Know what's essential vs optional; what's liable to change vs stable | Helps planning and trade-offs when encountering constraints |
| **Modifiable** | Document can evolve with version control, change history, flexible structuring | Requirements often change; document must support that |
| **Feasible** | No requirement demands something outside technical, budgetary, or resource constraints | Prevents unrealistic expectations and project failures |
| **Necessary/Justified** | Each requirement corresponds to real business or user need | Avoids scope bloat and unnecessary complexity |
| **Understandable** | Balance technical detail with clarity for non-technical readers | Ensures alignment, avoids miscommunication |

## Document Structure

Follow this structure to ensure completeness:

### 1. Introduction / Overview
- **Purpose**: What this document is and is not
- **Scope**: System boundaries and relationship to other systems
- **Definitions/Glossary**: Key terms and domain-specific vocabulary
- **Stakeholders**: Who is involved and their roles

### 2. Business / Stakeholder Needs / Drivers
- **Why**: Business objectives, problems being solved, goals
- **Constraints**: Legal, regulatory, environmental, operational
- **Success Criteria**: How success will be measured

### 3. Functional Requirements
- **What the system must do**
- Organized by features or modules
- Use-case descriptions or user stories
- Each requirement should include:
  - Unique ID (e.g., `REQ-001`, `SYS-UI-005`)
  - Description (what the system must do)
  - Acceptance criteria / verification method
  - Rationale (why this requirement exists)

### 4. Non-Functional Requirements (Quality Attributes)
- **Performance**: Response times, throughput, load capacity (with measurable metrics)
- **Reliability**: Uptime, error rates, fault tolerance
- **Scalability**: Growth expectations, capacity planning
- **Security**: Authentication, authorization, data protection
- **Usability**: User experience, accessibility, learnability
- **Maintainability**: Code quality, documentation, testability
- **Compatibility**: Platform, browser, device support

**Critical**: Use measurable metrics. Instead of "fast", specify "respond within 2 seconds for 95% of queries under normal load".

### 5. External Interfaces
- **System Interfaces**: APIs, protocols, communication methods
- **User Interfaces**: Screens, workflows, interaction patterns
- **Hardware/Software**: Dependencies, integrations, data sources
- **Data Formats**: Input/output formats, schemas

### 6. Constraints & Assumptions
- **Technical Constraints**: Technology stack, platform limitations
- **Resource Constraints**: Budget, personnel, timeline
- **Environmental Constraints**: Deployment environment, infrastructure
- **Assumptions**: What we assume to be true (dependencies, availability, etc.)

### 7. Acceptance Criteria & Verification Methods
- **How each requirement will be validated**
- Test cases or test types (unit, integration, acceptance, performance)
- Success metrics and thresholds
- Validation approach for each requirement category

### 8. Prioritization / Roadmap
- **Must-have vs Nice-to-have**: Use frameworks like MoSCoW (Must/Should/Could/Won't)
- **Phased Approach**: What can be deferred or implemented in phases
- **Dependencies**: Requirements that depend on others

### 9. Change Management & Traceability
- **Version Control**: Document version, revision history
- **Traceability Matrix**: Link requirements to:
  - Business goals (backward traceability)
  - Design components (forward traceability)
  - Test cases (verification traceability)
- **Change Process**: How requirements changes are managed

### 10. Appendices / Supporting Material
- **Diagrams**: Architecture, flowcharts, data models, user journeys
- **Mockups/Wireframes**: Visual representations of interfaces
- **References**: External documents, standards, related specifications
- **Glossary**: Expanded definitions

## Writing Guidelines

### Language & Style
- **Use imperative language**: "The system shall..." or "The system must..."
- **Avoid ambiguous modifiers**: Replace "should", "nice to have", "fast", "secure" with specific, measurable criteria
- **Use active voice**: Clear, direct statements
- **Define domain terms**: Include glossary for technical or domain-specific terms

### Requirement Format
Each requirement should follow this structure:

```
REQ-001: [Unique Identifier]
Description: [What the system must do]
Rationale: [Why this requirement exists]
Priority: [Must/Should/Could/Won't]
Acceptance Criteria: [How to verify this requirement is met]
Dependencies: [Other requirements this depends on]
```

### Visual Aids
- Use diagrams for complex flows, architectures, or interactions
- Include mockups for user interfaces
- Use flowcharts for business processes
- **Note**: Keep text as the "source of truth"; visuals support understanding

### Avoiding Common Pitfalls
- ❌ **Don't include design/implementation details** - Specify WHAT, not HOW
- ❌ **Don't use vague qualifiers** - Always provide measurable metrics
- ❌ **Don't mix status in requirements** - Status is managed separately (e.g., via ticket system)
- ❌ **Don't neglect traceability** - Always use unique identifiers and maintain links
- ❌ **Don't skip acceptance criteria** - Every requirement must be verifiable

## Review & Validation Process

1. **Stakeholder Review**: Involve users, domain experts, designers, QA early and often
2. **Technical Feasibility Review**: Ensure developers validate feasibility
3. **Completeness Check**: Verify all functional and non-functional requirements are captured
4. **Consistency Check**: Ensure no contradictions, uniform terminology
5. **Traceability Check**: Verify all requirements link to goals and can link to tests
6. **Iterative Refinement**: Requirements evolve; update document as understanding improves

## Version Control

- Maintain revision history with:
  - Version number
  - Date of change
  - Author of change
  - Description of change
  - Reason for change
- Use version control system (git) for the document itself
- Track changes to individual requirements, not just document versions

## Integration with Ticket System

ERDs and tickets share the same directory structure pattern:
- ERDs are managed in `erds/` with subdirectories: `todo/`, `in-progress/`, `done/`, `not-doing/`
- Tickets are managed in `tickets/` with the same subdirectories
- Both use zero-padded numbering for lexicographical sorting
- Both avoid status in file content (status is implied by directory location)

When creating tickets from ERD requirements:
- Link ticket to requirement ID (e.g., "Implements REQ-001")
- Ensure ticket acceptance criteria match ERD acceptance criteria
- Maintain traceability between ERD and implementation tickets

## Technical Debt ERDs

### Every 5th ERD Must Be Technical Debt

**Rule**: Every 5th ERD (e.g., ERD-005, ERD-010, ERD-015, ERD-020, etc.) must be a **technical debt paydown ERD**.

### Purpose
This ensures regular, systematic attention to technical debt and prevents it from accumulating unchecked. Technical debt ERDs focus on improving code quality, architecture, infrastructure, or processes rather than adding new features.

### Scope & Format
- **Technical debt ERDs can be very short** - they don't need the full structure of feature ERDs
- Focus on the specific technical debt items being addressed
- Reference the technical debt tracking system (see `tech-debt.md`)
- Include clear acceptance criteria for debt resolution
- Prioritize high-impact, high-risk technical debt

### What to Include
- **Debt Type**: Code, Architecture, Test, Documentation, Security, Process, Defect, Data, Design, or Knowledge debt
- **Impact**: How the debt affects development/productivity/security/reliability
- **Classification**: Deliberate/Inadvertent × Prudent/Reckless (see `tech-debt.md`)
- **Remediation Requirements**: What needs to be done to resolve the debt
- **Success Criteria**: How to verify the debt is resolved

### Examples
- ERD-005: Refactor large files for AI compatibility
- ERD-010: Improve test infrastructure and reduce flakiness
- ERD-015: Update security dependencies and patch vulnerabilities
- ERD-020: Break down monolithic components into smaller modules

### Counting
- Count ERDs sequentially by their number (005, 010, 015, 020, etc.)
- If an ERD is moved to `not-doing/`, it still counts toward the sequence
- Technical debt ERDs themselves count toward the sequence

# Breaking Down Sprint ERD into Tickets

Guidelines for decomposing a sprint's Engineering Requirements Document (ERD) into actionable tickets that can be managed in the ticket system.

## Overview

A sprint contains both a PRD (`prd.md`) and an ERD (`erd.md`) that together define **what** needs to be built at a high level. Tickets break this down into **actionable work items** that can be implemented, tested, and completed. This process bridges requirements and implementation.

**Structure**: Tickets are created within the sprint's `tickets/todo/` directory, nested under the sprint.

## Decomposition Principles

### 1. One Requirement May Become Multiple Tickets
- Large requirements should be broken into smaller, independent tickets
- Each ticket should be **completable** within a reasonable timeframe (typically 1-5 days)
- Tickets should represent **cohesive units of work** that deliver value

### 2. Ticket Scope Guidelines
- **Too Large**: If a ticket takes more than a week, break it down further
- **Too Small**: If a ticket takes less than a few hours, consider combining with related work
- **Just Right**: A ticket should be meaningful enough to stand alone but small enough to complete in one iteration

### 3. Maintain Traceability
- Every ticket must reference the ERD requirement(s) it implements
- Use ERD requirement IDs (e.g., "Implements REQ-001" or "Addresses REQ-005, REQ-007")
- This enables impact analysis when requirements change

## Mapping ERD Sections to Tickets

### Functional Requirements → Feature Tickets
**Process**:
1. Review each functional requirement in the ERD
2. Identify if it can be implemented as a single ticket or needs decomposition
3. Break down by:
   - **User-facing features** (one ticket per feature)
   - **Technical components** (backend, frontend, API endpoints)
   - **User flows** (authentication flow, checkout flow)
   - **Data models** (database schema, data access layer)

**Important**: Database schema changes must be in **separate tickets** from the code that uses them. See "Schema Changes" section below.

**Example**:
```
ERD REQ-001: "System shall allow users to authenticate via email and password"
↓
Ticket 001: Implement email/password authentication backend
Ticket 002: Create login UI component
Ticket 003: Add session management
Ticket 004: Implement logout functionality
```

### Non-Functional Requirements → Infrastructure/Quality Tickets
**Process**:
1. Non-functional requirements often span multiple features
2. Create separate tickets for:
   - **Performance requirements** (e.g., "Optimize API response time to <2s")
   - **Security requirements** (e.g., "Implement rate limiting")
   - **Scalability requirements** (e.g., "Add database connection pooling")
   - **Usability requirements** (e.g., "Add loading states to all async operations")

**Example**:
```
ERD NFR-001: "System shall respond within 2 seconds for 95% of API requests"
↓
Ticket 010: Add database query optimization
Ticket 011: Implement API response caching
Ticket 012: Add performance monitoring and alerting
```

### External Interfaces → Integration Tickets
**Process**:
1. Each external interface becomes one or more tickets
2. Break down by:
   - **API integration** (authentication, data exchange)
   - **Third-party services** (payment, email, analytics)
   - **Data import/export** (CSV, JSON, etc.)

### Constraints & Assumptions → Technical Debt/Setup Tickets
**Process**:
1. Constraints may require infrastructure or tooling tickets
2. Assumptions may need validation tickets
3. Create tickets for:
   - **Infrastructure setup** (deployment, CI/CD)
   - **Tooling** (monitoring, logging)
   - **Validation** (proof of concept, spike)

### Schema Changes → Separate Migration Tickets
**Critical Rule**: Database schema changes must be in **separate tickets** from the code that uses them.

**Process**:
1. Identify schema changes needed for ERD requirements
2. Create a dedicated ticket for each schema change:
   - Schema modification (table creation, column additions, etc.)
   - Goose migration file
   - Migration must consider current production state
   - Data migration strategies for updating existing data
   - Performance optimization not required (yet)
3. Create separate tickets for code that uses the schema:
   - These tickets depend on the schema change ticket
   - Implement features using the new schema

**Example**:
```
ERD REQ-001: "System shall store user accounts"
↓
Ticket 001: Create users table schema and migration
  - Schema: users table with email, password_hash, name
  - Migration: goose migration file
  - Data migration: handles any existing data if needed
  
Ticket 002: Implement user registration API
  - Blocked by: 001
  - Uses: users table from Ticket 001
```

**Why Separate?**
- Schema changes can be reviewed independently
- Migrations tested before dependent code
- Clear separation of concerns
- Easier rollback if needed
- Better traceability

See `tech-stack.md` and `ticket-system.md` for more details.

## Using sdlc.sh for Ticket Management

All ticket operations should use the `sdlc.sh` script rather than manually moving files. This ensures consistency and proper tracking.

### Creating Tickets
- Create ticket files directly in the sprint's `tickets/todo/` directory with the naming format `NNN-description.md`
- Use zero-padded numbers (001, 002, etc.) for proper lexicographical sorting
- The ticket will be automatically discovered by `sdlc.sh` which searches across all sprints
- Example path: `phases/{phase-state}/NNN-phase-name/sprints/{sprint-state}/NNN-sprint-name/tickets/todo/NNN-description.md`

### Moving Tickets Between States
Use `sdlc.sh move` command to move tickets:

```bash
# Move ticket to in-progress
./sdlc.sh move ticket 001 in-progress

# Move ticket to done
./sdlc.sh move ticket 001 done

# Move ticket to not-doing (if requirements removed)
./sdlc.sh move ticket 001 not-doing
```

**Note**: `sdlc.sh` automatically finds tickets across all sprints by number or name.

### Finding Tickets
Use `sdlc.sh get-next` to find tickets:

```bash
# Get next todo ticket (searches across all sprints)
./sdlc.sh get-next ticket todo

# Get next in-progress ticket (searches across all sprints)
./sdlc.sh get-next ticket in-progress
```

### Listing Tickets
Use `sdlc.sh list` to see all tickets in a state:

```bash
# List all todo tickets (across all sprints)
./sdlc.sh list ticket todo

# List all in-progress tickets (across all sprints)
./sdlc.sh list ticket in-progress
```

**Important**: Always use `sdlc.sh` commands rather than manually moving files between directories. This ensures the system maintains proper state tracking and respects sprint completion rules.

## Ticket Creation Process

### Step 1: Review Sprint Requirements
- Read through the entire sprint's ERD (`erd.md`) and PRD (`prd.md`)
- Identify all functional and non-functional requirements from the ERD
- Note dependencies between requirements
- Understand priorities (Must/Should/Could/Won't)
- Consider the sprint's context from the PRD

### Step 2: Group Related Requirements
- Identify requirements that work together
- Find natural boundaries for decomposition
- Consider user journeys and workflows

### Step 3: Create Ticket List
- Generate initial list of tickets
- Assign ticket numbers (zero-padded: 001, 002, etc.)
- Write descriptive ticket names

### Step 4: Define Dependencies
- Identify which tickets must be completed before others
- Document blocking relationships
- Consider parallel work opportunities

### Step 5: Map Acceptance Criteria
- For each ticket, extract relevant acceptance criteria from ERD
- Ensure ticket acceptance criteria align with ERD requirements
- Add ticket-specific acceptance criteria if needed

### Step 6: Prioritize Tickets
- Respect ERD priorities (Must/Should/Could/Won't)
- Consider technical dependencies
- Balance high-value work with foundational work

### Step 7: Commit and Push Changes
After creating all tickets, commit and push the changes using conventional commits:

```bash
# Stage all changes (ticket files are file-based, so use git add . to catch everything)
git add .

# Commit using conventional commit format
git commit -m "feat(tickets): break down sprint-XXX ERD into implementation tickets

- Created N tickets for [brief description of what was broken down]
- All tickets reference ERD requirements
- Dependencies documented between tickets
- Tickets created in sprint's tickets/todo/ directory"

# Push to remote repository
git push
```

**Conventional Commit Format**:
- **Type**: `feat` (for new tickets), `docs` (if only documentation), `refactor` (if reorganizing)
- **Scope**: `tickets` (always use this for ticket creation)
- **Description**: Brief summary of what was done
- **Body**: List key details about the tickets created

**Important**: 
- Always use `git add .` to ensure all ticket files are staged (file-based changes can be easy to miss)
- Use descriptive commit messages that reference the sprint number
- Include a count of tickets created in the commit message
- Push immediately after committing to ensure changes are saved

## Ticket Content Template

When creating tickets from ERD requirements, use this structure:

```markdown
# [Ticket Number]: [Descriptive Title]

## ERD Reference
Implements: REQ-001, REQ-003
Related to: REQ-002

## Description
[Clear description of what needs to be done]

## Context / Background
[Why this ticket exists, business value, user need]

## Acceptance Criteria
- [ ] [Specific, testable criterion from ERD]
- [ ] [Additional criterion]
- [ ] [Edge cases or non-happy paths]

## Technical Notes
[Implementation approach, technical constraints, design decisions]

## Dependencies
- Blocks: [Ticket numbers that depend on this]
- Blocked by: [Ticket numbers this depends on]
- Related: [Related ticket numbers]

## Resources / Links
- ERD: [link to ERD document]
- Designs: [link to mockups/wireframes]
- Docs: [link to relevant documentation]
```

## Best Practices

### ✅ Do
- **Maintain traceability**: Always link tickets to ERD requirements
- **Keep tickets focused**: One ticket = one cohesive piece of work
- **Use ERD acceptance criteria**: Copy relevant criteria from ERD to tickets
- **Document dependencies**: Make blocking relationships explicit
- **Respect priorities**: Honor ERD prioritization (Must/Should/Could/Won't)
- **Break down large requirements**: Don't create tickets that are too large
- **Consider user value**: Each ticket should deliver some value, even if small
- **Follow TDD approach**: Create endpoint structure first, hardcode simple responses, write tests that fail for legitimate reasons, then implement (see `tech-stack.md`)

### ❌ Don't
- **Don't skip requirements**: Every ERD requirement should map to at least one ticket
- **Don't mix unrelated work**: Keep tickets focused on single concerns
- **Don't lose context**: Include enough background from ERD in tickets
- **Don't ignore dependencies**: Document what must happen before/after
- **Don't change acceptance criteria**: Ticket criteria should match ERD unless explicitly refined
- **Don't create tickets without ERD reference**: Always link back to source requirement

## Handling Sprint ERD Changes

When sprint ERD requirements change:

1. **Identify affected tickets**: Use traceability to find tickets implementing changed requirements (tickets are in the sprint's `tickets/` directory)
2. **Update or create tickets**: 
   - Update existing tickets if requirements refined
   - Create new tickets in the sprint's `tickets/todo/` if requirements added
   - Move tickets to `not-doing/` using `./sdlc.sh move ticket <number> not-doing` if requirements removed
3. **Review dependencies**: Check if changes affect ticket dependencies
4. **Update acceptance criteria**: Ensure ticket criteria match updated ERD requirements
5. **Sprint completion**: Remember that a sprint cannot be moved to `done` if it has any tickets in `todo` or `in-progress`

## Example: Complete Breakdown

**ERD Requirement**:
```
REQ-001: User Authentication
Description: System shall allow users to authenticate using email and password
Acceptance Criteria:
- User can register with email and password
- User can login with valid credentials
- User receives error for invalid credentials
- Password must meet security requirements (8+ chars, special chars)
Priority: Must
```

**Resulting Tickets**:
```
001-implement-user-registration-api.md
  - Implements: REQ-001
  - Creates user registration endpoint
  - Validates password requirements
  - Stores hashed passwords

002-create-registration-ui.md
  - Implements: REQ-001
  - Blocked by: 001
  - Registration form component
  - Password validation UI
  - Success/error handling

003-implement-login-api.md
  - Implements: REQ-001
  - Blocked by: 001
  - Login endpoint
  - Credential validation
  - Session/token generation

004-create-login-ui.md
  - Implements: REQ-001
  - Blocked by: 003
  - Login form component
  - Error message display
  - Redirect after successful login
```

## Validation Checklist

Before considering sprint ERD breakdown complete:

- [ ] Every ERD requirement maps to at least one ticket
- [ ] All tickets reference their ERD requirement(s)
- [ ] Ticket acceptance criteria align with ERD acceptance criteria
- [ ] Dependencies are documented
- [ ] Priorities are respected
- [ ] Tickets are appropriately sized (not too large, not too small)
- [ ] Related tickets are grouped logically
- [ ] Non-functional requirements have corresponding tickets
- [ ] All tickets are created in the sprint's `tickets/todo/` directory
- [ ] All ticket files have been committed using conventional commits
- [ ] Changes have been pushed to the remote repository
- [ ] Sprint structure is correct (tickets nested under sprint)
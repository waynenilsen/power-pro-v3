# Breaking Down ERD into Tickets

Guidelines for decomposing an Engineering Requirements Document (ERD) into actionable tickets that can be managed in the ticket system.

## Overview

An ERD defines **what** needs to be built at a high level. Tickets break this down into **actionable work items** that can be implemented, tested, and completed. This process bridges requirements and implementation.

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

## Ticket Creation Process

### Step 1: Review ERD Requirements
- Read through the entire ERD
- Identify all functional and non-functional requirements
- Note dependencies between requirements
- Understand priorities (Must/Should/Could/Won't)

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

### ❌ Don't
- **Don't skip requirements**: Every ERD requirement should map to at least one ticket
- **Don't mix unrelated work**: Keep tickets focused on single concerns
- **Don't lose context**: Include enough background from ERD in tickets
- **Don't ignore dependencies**: Document what must happen before/after
- **Don't change acceptance criteria**: Ticket criteria should match ERD unless explicitly refined
- **Don't create tickets without ERD reference**: Always link back to source requirement

## Handling ERD Changes

When ERD requirements change:

1. **Identify affected tickets**: Use traceability to find tickets implementing changed requirements
2. **Update or create tickets**: 
   - Update existing tickets if requirements refined
   - Create new tickets if requirements added
   - Move tickets to `not-doing/` if requirements removed
3. **Review dependencies**: Check if changes affect ticket dependencies
4. **Update acceptance criteria**: Ensure ticket criteria match updated ERD requirements

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

Before considering ERD breakdown complete:

- [ ] Every ERD requirement maps to at least one ticket
- [ ] All tickets reference their ERD requirement(s)
- [ ] Ticket acceptance criteria align with ERD acceptance criteria
- [ ] Dependencies are documented
- [ ] Priorities are respected
- [ ] Tickets are appropriately sized (not too large, not too small)
- [ ] Related tickets are grouped logically
- [ ] Non-functional requirements have corresponding tickets

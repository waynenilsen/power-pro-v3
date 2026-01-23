# Breaking Down ERD Requirements into Crumbs

Guidelines for decomposing an Engineering Requirements Document (ERD) into actionable crumbs that can be managed with crumbler.

## Overview

An ERD defines **what** needs to be built at a high level. Crumbs break this down into **actionable work items** that can be implemented, tested, and completed. This process bridges requirements and implementation.

**Structure**: Crumbs are created using `./crumbler create "Task Name"` which creates a directory with a README.md file containing the task instructions.

## Decomposition Principles

### 1. One Requirement May Become Multiple Crumbs
- Large requirements should be broken into smaller, independent crumbs
- Each crumb should be **completable** within a reasonable timeframe (typically 1-5 days)
- Crumbs should represent **cohesive units of work** that deliver value

### 2. Crumb Scope Guidelines
- **Too Large**: If a crumb takes more than a week, break it down further using `./crumbler create`
- **Too Small**: If a crumb takes less than a few hours, consider combining with related work
- **Just Right**: A crumb should be meaningful enough to stand alone but small enough to complete in one iteration

### 3. Maintain Traceability
- Every crumb must reference the ERD requirement(s) it implements
- Use ERD requirement IDs (e.g., "Implements REQ-001" or "Addresses REQ-005, REQ-007")
- This enables impact analysis when requirements change

## Mapping ERD Sections to Crumbs

### Functional Requirements → Feature Crumbs
**Process**:
1. Review each functional requirement in the ERD
2. Identify if it can be implemented as a single crumb or needs decomposition
3. Break down by:
   - **User-facing features** (one crumb per feature)
   - **Technical components** (backend, API endpoints)
   - **User flows** (authentication flow, checkout flow)
   - **Data models** (database schema, data access layer)

**Important**: Database schema changes must be in **separate crumbs** from the code that uses them. See "Schema Changes" section below.

**Example**:
```
ERD REQ-001: "System shall allow users to authenticate via email and password"
↓
Crumb: Implement email/password authentication backend
Crumb: Add session management
Crumb: Implement logout functionality
```

### Non-Functional Requirements → Infrastructure/Quality Crumbs
**Process**:
1. Non-functional requirements often span multiple features
2. Create separate crumbs for:
   - **Performance requirements** (e.g., "Optimize API response time to <2s")
   - **Security requirements** (e.g., "Implement rate limiting")
   - **Scalability requirements** (e.g., "Add database connection pooling")
   - **Usability requirements** (e.g., "Add loading states to all async operations")

**Example**:
```
ERD NFR-001: "System shall respond within 2 seconds for 95% of API requests"
↓
Crumb: Add database query optimization
Crumb: Implement API response caching
Crumb: Add performance monitoring and alerting
```

### External Interfaces → Integration Crumbs
**Process**:
1. Each external interface becomes one or more crumbs
2. Break down by:
   - **API integration** (authentication, data exchange)
   - **Third-party services** (payment, email, analytics)
   - **Data import/export** (CSV, JSON, etc.)

### Constraints & Assumptions → Technical Debt/Setup Crumbs
**Process**:
1. Constraints may require infrastructure or tooling crumbs
2. Assumptions may need validation crumbs
3. Create crumbs for:
   - **Infrastructure setup** (deployment, CI/CD)
   - **Tooling** (monitoring, logging)
   - **Validation** (proof of concept, spike)

### Schema Changes → Separate Migration Crumbs
**Critical Rule**: Database schema changes must be in **separate crumbs** from the code that uses them.

**Process**:
1. Identify schema changes needed for ERD requirements
2. Create a dedicated crumb for each schema change:
   - Schema modification (table creation, column additions, etc.)
   - Goose migration file
   - Migration must consider current production state
   - Data migration strategies for updating existing data
   - Performance optimization not required (yet)
3. Create separate crumbs for code that uses the schema:
   - These crumbs depend on the schema change crumb
   - Implement features using the new schema

**Example**:
```
ERD REQ-001: "System shall store user accounts"
↓
Crumb: Create users table schema and migration
  - Schema: users table with email, password_hash, name
  - Migration: goose migration file
  - Data migration: handles any existing data if needed
  
Crumb: Implement user registration API
  - Depends on: schema crumb
  - Uses: users table from schema crumb
```

**Why Separate?**
- Schema changes can be reviewed independently
- Migrations tested before dependent code
- Clear separation of concerns
- Easier rollback if needed
- Better traceability

See `tech-stack.md` and `ticket-system.md` for more details.

## Using Crumbler for Work Management

Work is managed using crumbler - a simple task decomposition system. The filesystem IS the state - existence means work to do, deletion means done. Work depth-first: complete children before parents.

**Key Commands**:
```bash
# Get AI instructions for current work
./crumbler prompt

# Create a new crumb (sub-task)
./crumbler create "Task Name"

# Mark current crumb as done
./crumbler delete

# View crumb tree and status
./crumbler status

# Get help
./crumbler help
```

**Workflow**:
1. `crumbler prompt` - Get AI instructions for what to do next
2. Do the work - Follow the instructions
3. `crumbler delete` - Mark crumb as done when complete
4. Repeat

## Crumb Creation Process

### Step 1: Review ERD Requirements
- Read through the entire ERD
- Identify all functional and non-functional requirements
- Note dependencies between requirements
- Understand priorities (Must/Should/Could/Won't)

### Step 2: Group Related Requirements
- Identify requirements that work together
- Find natural boundaries for decomposition
- Consider user journeys and workflows

### Step 3: Create Crumbs
- Use `./crumbler create "Task Name"` to create each crumb
- Crumbs are automatically numbered and organized
- Write descriptive crumb names

### Step 4: Define Dependencies
- Identify which crumbs must be completed before others
- Document blocking relationships in crumb README files
- Consider parallel work opportunities

### Step 5: Map Acceptance Criteria
- For each crumb, extract relevant acceptance criteria from ERD
- Ensure crumb acceptance criteria align with ERD requirements
- Add crumb-specific acceptance criteria if needed

### Step 6: Prioritize Crumbs
- Respect ERD priorities (Must/Should/Could/Won't)
- Consider technical dependencies
- Balance high-value work with foundational work

## Crumb Content Template

When creating crumbs from ERD requirements, use this structure in the README.md:

```markdown
# [Descriptive Title]

## ERD Reference
Implements: REQ-001, REQ-003
Related to: REQ-002

## Description
[Clear description of what needs to be done]

## Context / Background
[Why this crumb exists, business value, user need]

## Acceptance Criteria
- [ ] [Specific, testable criterion from ERD]
- [ ] [Additional criterion]
- [ ] [Edge cases or non-happy paths]

## Technical Notes
[Implementation approach, technical constraints, design decisions]

## Dependencies
- Blocks: [Crumbs that depend on this]
- Blocked by: [Crumbs this depends on]
- Related: [Related crumbs]
```

## Best Practices

### ✅ Do
- **Maintain traceability**: Always link crumbs to ERD requirements
- **Keep crumbs focused**: One crumb = one cohesive piece of work
- **Use ERD acceptance criteria**: Copy relevant criteria from ERD to crumbs
- **Document dependencies**: Make blocking relationships explicit in crumb READMEs
- **Respect priorities**: Honor ERD prioritization (Must/Should/Could/Won't)
- **Break down large requirements**: Don't create crumbs that are too large
- **Consider user value**: Each crumb should deliver some value, even if small
- **Follow TDD approach**: Create endpoint structure first, hardcode simple responses, write tests that fail for legitimate reasons, then implement (see `tech-stack.md`)

### ❌ Don't
- **Don't skip requirements**: Every ERD requirement should map to at least one crumb
- **Don't mix unrelated work**: Keep crumbs focused on single concerns
- **Don't lose context**: Include enough background from ERD in crumbs
- **Don't ignore dependencies**: Document what must happen before/after
- **Don't change acceptance criteria**: Crumb criteria should match ERD unless explicitly refined
- **Don't create crumbs without ERD reference**: Always link back to source requirement

## Handling ERD Changes

When ERD requirements change:

1. **Identify affected crumbs**: Use traceability to find crumbs implementing changed requirements
2. **Update or create crumbs**: 
   - Update existing crumb READMEs if requirements refined
   - Create new crumbs using `./crumbler create` if requirements added
   - Delete crumbs using `./crumbler delete` if requirements removed
3. **Review dependencies**: Check if changes affect crumb dependencies
4. **Update acceptance criteria**: Ensure crumb criteria match updated ERD requirements

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

**Resulting Crumbs**:
```
Crumb: Implement user registration API
  - Implements: REQ-001
  - Creates user registration endpoint
  - Validates password requirements
  - Stores hashed passwords

Crumb: Implement login API
  - Implements: REQ-001
  - Depends on: registration crumb
  - Login endpoint
  - Credential validation
  - Session/token generation
```

## Validation Checklist

Before considering ERD breakdown complete:

- [ ] Every ERD requirement maps to at least one crumb
- [ ] All crumbs reference their ERD requirement(s)
- [ ] Crumb acceptance criteria align with ERD acceptance criteria
- [ ] Dependencies are documented in crumb READMEs
- [ ] Priorities are respected
- [ ] Crumbs are appropriately sized (not too large, not too small)
- [ ] Related crumbs are grouped logically
- [ ] Non-functional requirements have corresponding crumbs
- [ ] All crumbs are created using `./crumbler create`

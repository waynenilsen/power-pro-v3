# Technical Debt

Guidelines for identifying, categorizing, and managing technical debt based on industry best practices and expert recommendations.

## Common Types of Technical Debt

Technical debt can be categorized into the following common types:

### 1. Code Debt
**Definition**: Issues in the source code that make it fragile, hard to understand, modify, or extend.

**Examples**:
- Duplicated logic across multiple files
- Hard-coded values instead of configuration
- Large, complex functions that do too much
- **Very long files** (exceeding AI context limits, typically 1000+ lines)
- Files that are too large to fit in AI context windows
- Tight coupling between components
- Poor naming conventions
- Lack of modularization or proper layering
- Deeply nested conditionals
- Missing error handling

**Risks**: Increased bugs, slower feature delivery, high maintenance cost, reduced developer productivity

**Typical Causes**: Tight deadlines, insufficient code review, lack of coding standards, pressure to ship quickly

**Remediation Strategies**: 
- Enforce style and architecture guidelines
- Schedule regular refactoring sprints
- Incorporate static analysis tools
- Promote pair programming and code reviews
- **Break long files into smaller, focused modules** (critical for AI compatibility)

### 2. Architecture & Infrastructure Debt
**Definition**: Structural decisions at the system level that don't scale or adapt to evolving requirements.

**Examples**:
- Monolithic architecture when microservices would be better
- Tightly coupled components
- Outdated frameworks or libraries
- Single points of failure
- Performance bottlenecks
- Inefficient data storage or access patterns
- Legacy systems that are hard to maintain
- Infrastructure drift (servers/OS out of date)

**Risks**: High cost to change, limits future adaptability, increasing risk over time, scalability issues

**Typical Causes**: Rapid growth, changing requirements, lack of architectural planning, budget constraints

**Remediation Strategies**:
- Architectural reviews and refactoring
- Gradual migration strategies
- Performance testing and optimization
- Infrastructure as code (IaC)
- Regular dependency updates

### 3. Test Debt
**Definition**: Insufficient test coverage or weak/fragile testing infrastructure.

**Examples**:
- Missing unit tests for critical paths
- Low test coverage overall
- Flaky tests that fail intermittently
- Slow test suites that delay feedback
- Manual testing where automation is needed
- Missing integration or end-to-end tests
- Tests that don't catch real bugs
- No regression safety nets

**Risks**: Leads to regressions, high maintenance effort, reduced confidence in shipping changes, production bugs

**Typical Causes**: Time pressure, lack of testing culture, insufficient tooling, legacy codebases

**Remediation Strategies**:
- Increase test coverage incrementally
- Fix flaky tests immediately
- Automate manual test processes
- Implement CI/CD with test gates
- Test-driven development (TDD) practices

### 4. Documentation Debt
**Definition**: Missing, outdated, or inaccurate documentation.

**Examples**:
- Outdated architecture diagrams
- Missing or stale API documentation
- Incomplete onboarding guides
- Absent design decision records (ADRs)
- Outdated setup/installation instructions
- Missing code comments for complex logic
- Stale README files
- No runbooks for operations

**Risks**: Slows new developers, causes misalignment about system behavior, increases maintenance cost, knowledge loss

**Typical Causes**: Documentation seen as low priority, rapid changes, lack of documentation standards, time pressure

**Remediation Strategies**:
- Documentation as part of definition of done
- Automated documentation generation where possible
- Regular documentation reviews
- Living documentation (docs close to code)
- Knowledge sharing sessions

### 5. Security Debt
**Definition**: Postponed or inadequate security practices and controls.

**Examples**:
- Unpatched dependencies with known vulnerabilities
- Hard-coded credentials or secrets
- Weak authentication or authorization mechanisms
- Missing security audits or scans
- Lack of proper access controls
- Outdated encryption standards
- No security monitoring or logging
- Missing security headers or configurations

**Risks**: Data breaches, compliance failures, legal exposure, reputation damage, financial losses

**Typical Causes**: Deadline pressure, lack of security expertise, deprioritized security work, complexity

**Remediation Strategies**:
- Automated dependency scanning
- Regular security audits
- Security training and secure coding standards
- Automated vulnerability checks in CI/CD
- Secrets management systems
- Security reviews in code review process

### 6. Process & People Debt
**Definition**: Inefficiencies or gaps in team workflows and organizational practices.

**Examples**:
- Lack of consistent code review process
- Missing or unclear coding standards
- Knowledge silos (only one person knows X)
- Poor cross-team communication
- Unclear ownership of components
- Inadequate requirements gathering
- Skill gaps in the team
- No retrospectives or improvement process

**Risks**: Undermines code quality and velocity, escalates other types of debt indirectly, reduces team effectiveness

**Typical Causes**: Rapid growth, lack of process definition, insufficient training, organizational issues

**Remediation Strategies**:
- Establish clear processes and standards
- Regular team retrospectives
- Knowledge sharing and documentation
- Cross-training and pair programming
- Clear ownership and responsibilities

### 7. Defect/Bug Debt
**Definition**: Known bugs or technical issues that have been accepted but deferred.

**Examples**:
- Known bugs left unresolved
- Edge cases not handled
- Accumulating defects that degrade user experience
- Workarounds instead of proper fixes
- Low-priority bugs that never get fixed
- Technical issues shipped with features

**Risks**: Bugs accumulate, leading to crashes or regressions, inflating future fix costs, user frustration

**Typical Causes**: Feature pressure, unclear bug prioritization, lack of bug triage process, time constraints

**Remediation Strategies**:
- Regular bug triage sessions
- Bug debt reduction sprints
- Clear bug prioritization framework
- Track bug debt metrics
- Balance new features with bug fixes

### 8. Data Debt
**Definition**: Issues with data models, schemas, or data quality.

**Examples**:
- Poor schema design
- Inconsistent data formats
- Missing data validation
- Outdated or messy data pipelines
- Lack of data lineage or monitoring
- Inefficient database queries
- Missing indexes
- Duplicate or incorrect data

**Risks**: Impacts reliability, performance, integration, and reporting; data quality issues compound over time

**Typical Causes**: Rapid development, changing requirements, lack of data governance, insufficient planning

**Remediation Strategies**:
- Data modeling reviews
- Data quality monitoring
- Schema migrations and cleanup
- Query optimization
- Data governance practices

### 9. Design/UI Debt
**Definition**: Hasty or inconsistent design decisions, especially in user interfaces.

**Examples**:
- Inconsistent design language or components
- Confusing navigation patterns
- Quick UI implementations without design review
- Accessibility issues
- Responsive design gaps
- Outdated UI patterns
- Inconsistent user experience

**Risks**: Poor user experience, increased support requests, redesign costs, accessibility compliance issues

**Typical Causes**: Rapid prototyping, lack of design system, insufficient design review, time pressure

**Remediation Strategies**:
- Design system implementation
- Regular design reviews
- Accessibility audits
- User testing and feedback
- Consistent design patterns

### 10. Knowledge/People Debt
**Definition**: Loss of tacit or explicit knowledge, lack of training, or over-reliance on certain individuals.

**Examples**:
- Key engineer leaves with undocumented knowledge
- Onboarding takes too long
- Only one person understands critical systems
- Missing knowledge transfer
- Lack of training materials
- Tribal knowledge not documented

**Risks**: Knowledge loss, slower onboarding, single points of failure, reduced team resilience

**Typical Causes**: Rapid turnover, insufficient documentation, lack of knowledge sharing culture, time pressure

**Remediation Strategies**:
- Knowledge sharing sessions
- Documentation and runbooks
- Pair programming and mentoring
- Regular knowledge transfer
- Cross-training programs

## File Chunking for AI Compatibility

### Why File Size Matters for AI

AI code assistants (like LLMs) have context window limitations. Very long files cannot be fully loaded into the AI's context, making it difficult or impossible for AI to:
- Understand the full file structure
- Make accurate changes across the entire file
- Provide comprehensive code analysis
- Generate code that fits the existing patterns

**Critical for AI-assisted development**: Files that exceed AI context limits become effectively unmaintainable with AI assistance.

### File Size Guidelines

**Recommended Limits**:
- **Target**: Files under 500 lines of code
- **Warning**: Files between 500-1000 lines should be reviewed for splitting opportunities
- **Critical**: Files over 1000 lines should be broken down into smaller modules

**Context Window Considerations**:
- Most AI assistants have context windows of 8K-128K tokens
- Code typically uses 1-2 tokens per line
- Account for other context (conversation history, other files, system prompts)
- **Practical limit**: ~1000-2000 lines per file for reliable AI assistance

### When to Chunk Files

Break files into smaller chunks when:
- File exceeds 1000 lines
- File contains multiple distinct responsibilities
- File mixes concerns (e.g., data models + business logic + API handlers)
- AI assistant struggles to understand or modify the file
- File has grown organically without refactoring
- File contains multiple classes/modules that could be separate

### Chunking Strategies

1. **By Responsibility**: Split files by single responsibility principle
   - One class/module per file
   - Separate data models from business logic
   - Separate API handlers from core logic

2. **By Feature/Module**: Group related functionality
   - Feature-based organization
   - Domain-driven design boundaries
   - Module boundaries

3. **By Layer**: Separate by architectural layers
   - Presentation layer (UI/API)
   - Business logic layer
   - Data access layer
   - Infrastructure layer

4. **Extract Common Patterns**: Identify and extract reusable code
   - Shared utilities
   - Common interfaces
   - Shared types/interfaces

### Chunking Best Practices

**✅ Do**:
- Create focused, single-purpose files
- Use clear, descriptive file names
- Maintain clear module boundaries
- Keep related code together
- Use proper imports/exports
- Document module interfaces
- Test each chunk independently

**❌ Don't**:
- Create artificial splits that break logical cohesion
- Split files just to meet line count (split by responsibility)
- Create circular dependencies between chunks
- Lose the big picture (maintain clear module relationships)

### Example: Before and After

**Before** (2000+ lines):
```
utils.js - Contains: validation, formatting, API helpers, data transformation, error handling, logging, configuration
```

**After** (chunked):
```
validation.js - Validation logic (~200 lines)
formatting.js - Data formatting utilities (~150 lines)
api-helpers.js - API request/response helpers (~300 lines)
data-transform.js - Data transformation functions (~250 lines)
error-handling.js - Error handling utilities (~200 lines)
logging.js - Logging utilities (~150 lines)
config.js - Configuration management (~100 lines)
```

### Remediation Steps

1. **Identify large files**: Use tools to find files exceeding thresholds
2. **Analyze responsibilities**: Map what each file does
3. **Plan chunking**: Design how to split while maintaining functionality
4. **Extract incrementally**: Move code to new files one chunk at a time
5. **Update imports**: Fix all references to moved code
6. **Test thoroughly**: Ensure functionality remains intact
7. **Update documentation**: Reflect new file structure

### Tools for Identifying Large Files

- `find . -name "*.js" -exec wc -l {} \; | sort -rn | head -20`
- IDE plugins that show file size
- Static analysis tools
- Git history analysis (track file growth)

### Priority

**High Priority** if:
- File is actively being modified
- File blocks AI-assisted development
- File is in critical path
- Multiple developers struggle with the file

**Medium Priority** if:
- File is stable but large
- File is rarely modified
- Chunking would improve maintainability

## Technical Debt Quadrants

Classify technical debt by intent and recklessness (Martin Fowler's Technical Debt Quadrant):

| | **Prudent** | **Reckless** |
|---|---|---|
| **Deliberate** | **Deliberate & Prudent**: You know you're incurring debt and have a plan to repay it. Example: MVP trade-offs, early launch decisions. | **Deliberate & Reckless**: Conscious shortcuts without a repayment plan. Example: Skipping critical security steps, ignoring best practices. |
| **Inadvertent** | **Inadvertent & Prudent**: Debt that arises from growth or learning. What was acceptable becomes suboptimal as the system evolves. Example: Architecture that worked for 100 users but not 1M. | **Inadvertent & Reckless**: Unnoticed bad practices from lack of training or oversight. Example: Poor coding practices due to insufficient code review. |

**Priority Guidelines**:
- **Deliberate & Reckless**: Highest priority - address immediately
- **Inadvertent & Reckless**: High priority - address soon
- **Deliberate & Prudent**: Medium priority - follow repayment plan
- **Inadvertent & Prudent**: Lower priority - address when refactoring

## Identifying Technical Debt

### Signs of Technical Debt
- **Velocity Decline**: Team velocity decreasing over time
- **Increasing Bug Rate**: More bugs appearing in production
- **Longer Development Cycles**: Features taking longer to implement
- **Fear of Change**: Developers hesitant to modify certain areas
- **Knowledge Silos**: Only specific people can work on certain areas
- **Build/Deploy Issues**: Frequent build failures or deployment problems
- **Technical Spikes**: Regular need for investigation spikes
- **Code Smells**: Duplication, complexity, tight coupling

### Assessment Questions
- How long does it take to add a new feature?
- How often do we need to fix bugs in this area?
- How confident are we making changes here?
- What would happen if [key person] left?
- How much time is spent on maintenance vs new features?
- Are there areas developers avoid?

## Managing Technical Debt

### Prioritization Framework

When prioritizing technical debt, consider:

1. **Impact**: How does this debt affect:
   - Developer productivity
   - System reliability
   - Security
   - User experience
   - Business goals

2. **Cost of Delay**: What happens if we don't address this?
   - Increasing maintenance cost
   - Risk of incidents
   - Slowing feature development
   - Technical constraints

3. **Cost to Fix**: How much effort to resolve?
   - Time required
   - Risk of breaking changes
   - Dependencies
   - Resource availability

4. **Quadrant Classification**: Deliberate/Inadvertent × Prudent/Reckless

### Integration with Sprint and Ticket System

Technical debt should be tracked as sprints and tickets:

- **Sprint Structure**: Technical debt sprints follow the same structure as feature sprints (contain `prd.md` and `erd.md`)
- **Every 5th Sprint**: Must be a technical debt sprint (see `roadmap.md` and `erd.md`)
- **Ticket Naming**: Use zero-padded numbers (e.g., `001-fix-duplicated-logic.md`)
- **Category Tag**: Include debt type in ticket (e.g., `[Code Debt]`, `[Security Debt]`)
- **Priority**: Based on quadrant and impact assessment
- **Acceptance Criteria**: Define what "debt resolved" means
- **Link to Sprint ERD**: Tickets reference sprint ERD requirements

### Example Ticket Format

```markdown
# 001-[Category] Short description of debt

## Debt Type
[Code/Architecture/Test/Documentation/Security/Process/Defect/Data/Design/Knowledge Debt]

## Description
[What is the debt and where is it located]

## Impact
- [How this affects development/productivity/security/etc.]

## Classification
- Intent: [Deliberate/Inadvertent]
- Recklessness: [Prudent/Reckless]
- Priority: [High/Medium/Low]

## Remediation Steps
- [ ] Step 1
- [ ] Step 2
- [ ] Step 3

## Acceptance Criteria
- [ ] [Criterion 1]
- [ ] [Criterion 2]

## Related
- Blocks: [Ticket numbers]
- Related to: [ERD requirement IDs]
```

## Best Practices

### ✅ Do
- **Track debt explicitly**: Make technical debt visible in your planning
- **Balance debt and features**: Allocate time for debt reduction
- **Prioritize by impact**: Focus on high-impact, high-risk debt first
- **Create repayment plans**: For deliberate debt, have a plan
- **Regular assessment**: Review and update debt inventory regularly
- **Prevent new debt**: Establish practices to avoid creating new debt
- **Measure debt**: Track metrics like bug rate, velocity, build times

### ❌ Don't
- **Don't ignore debt**: It compounds over time
- **Don't accumulate reckless debt**: Address deliberate & reckless immediately
- **Don't let debt block features**: Balance is key
- **Don't skip code reviews**: They catch inadvertent debt
- **Don't defer security debt**: Security debt is high risk
- **Don't create knowledge silos**: Share knowledge proactively

## Debt Reduction Strategies

1. **Dedicated Debt Sprints**: Allocate entire sprints to debt reduction
2. **20% Rule**: Spend 20% of time on technical improvement
3. **Boy Scout Rule**: Leave code better than you found it
4. **Refactoring as You Go**: Fix debt when touching related code
5. **Debt Triage**: Regular sessions to assess and prioritize debt
6. **Prevention**: Establish practices to prevent new debt

## Validation Checklist

When creating a technical debt ticket:

- [ ] Debt type is clearly identified
- [ ] Impact is described (productivity, security, reliability, etc.)
- [ ] Classification (quadrant) is assigned
- [ ] Priority is set based on impact and classification
- [ ] Remediation steps are defined
- [ ] Acceptance criteria are clear
- [ ] Related tickets/ERDs are linked
- [ ] Cost to fix is estimated
- [ ] Cost of delay is understood

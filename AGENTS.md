# Agents

## Ticket System
[Ticket System](./prompts/ticket-system.md)

## Engineering Requirements Document
[ERD](./prompts/erd.md)

## Breaking Down ERD into Tickets
[ERD to Tickets](./prompts/erd-to-tickets.md)

## Software Roadmap
[Roadmap](./prompts/roadmap.md)

## Technical Debt
[Technical Debt](./prompts/tech-debt.md)

## Technology Stack
[Tech Stack](./prompts/tech-stack.md)

## Authentication System
[Authentication](./prompts/auth.md)

## Document Hierarchy & Relationships

The system follows a hierarchical breakdown structure:

- **1 Phase Document** → **at least 3 PRDs** (Product Requirements Documents)
- **1 PRD** → **precisely 1 ERD** (Engineering Requirements Document)
- **1 ERD** → **at least 5 Tickets**

This ensures proper decomposition and traceability throughout the development process.

### Technical Debt Rule

**Every 5th PRD/ERD must be a technical debt paydown PRD/ERD** (e.g., PRD-005, PRD-010, ERD-005, ERD-010, etc.). Technical debt PRDs/ERDs can be very short and focus on improving code quality, architecture, infrastructure, or processes rather than adding new features. See the ERD and Roadmap guidelines for details.

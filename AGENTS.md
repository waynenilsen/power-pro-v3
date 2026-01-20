# Agents

Read the README to learn more about what we are doing. This project is PowerPro. It is a headless api there is no frontend and there should not be a frontend. The core challenge of powerpro is keeping the code DRY. All of these powerlifters are doing things slightly differently but at the end of the day they're really doing very similar things. It is our job as programmers to manage this complexity in a re-usable way. The difficulty will come from managing the complexity related to that. When designing the entities and their relationship to eachother and the classes even in early phases it will be critical to design it such that we don't end up with spaghetti and you must keep your eye on the horizon. Don't preemptively add fields but you may preemptively add needed layers of abstraction that may well start as trivial / passthrough for initial implementations.

If you are working on a ticket, ERD, PRD, or roadmap Phase you must use sdlc.sh to manipulate those objects.

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

## Powerlifting Programs In Depth
[Programs](./programs/*.md)

## Document Hierarchy & Relationships

The system follows a hierarchical breakdown structure:

- **1 Phase Document** → **at least 3 PRDs** (Product Requirements Documents)
- **1 PRD** → **precisely 1 ERD** (Engineering Requirements Document)
- **1 ERD** → **at least 5 Tickets**

This ensures proper decomposition and traceability throughout the development process.

### Technical Debt Rule

**Every 5th PRD/ERD must be a technical debt paydown PRD/ERD** (e.g., PRD-005, PRD-010, ERD-005, ERD-010, etc.). Technical debt PRDs/ERDs can be very short and focus on improving code quality, architecture, infrastructure, or processes rather than adding new features. See the ERD and Roadmap guidelines for details.

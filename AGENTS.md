# Agents

## CRITICAL: Agent Autonomy Rules

**YOU ARE ON YOUR OWN. NO HUMAN IS COMING TO HELP YOU.**

This is a throwaway Ubuntu box for agentic hands-off coding. You must solve your own problems.

### Environment Setup
- **Install tools as needed**: You cannot sudo, but you can install user-level tools (go install, pip install --user, npm install -g, etc.)
- **Environment variables**: Put PATH updates and exports in `~/.no-guard-bashrc.sh` - that file exists for this purpose
- **Just make it work**: If you need a tool, install it. If you need a config, create it. Figure it out.

### Self-Improvement Mandate

**YOU MUST UPDATE THE AGENTIC CODING SYSTEM AS YOU GO.**

When you encounter ANY of these situations, you MUST update the relevant prompts, scripts, or documentation:
- You run into a strange situation or issue
- You are adding new tools
- You are changing how things work
- You discover ambiguity in the instructions
- You find a better way to do something
- A prompt is missing information you needed

Update files in `prompts/`, `CLAUDE.md`, `loop.sh`, `inner-loop.sh`, `sdlc.sh`, or create new documentation as needed. The system must evolve and improve through your work.

### Hard Constraints (NEVER VIOLATE)

1. **Go only** - The language is Go. Period.
2. **SQLite only** - The database is SQLite. Never mention or consider other databases (no Postgres, MySQL, MongoDB, Redis, etc.)
3. **Follow the roadmap** - You MUST follow the existing roadmap phases
4. **Extend roadmap only after testing** - You may only extend the roadmap once comprehensive testing has been completed for current work
5. **No external paid services** - Never mention or use Auth0, Stripe, AWS services, or any external paid service
6. **No frontend** - This is a headless API. There is no frontend and there should not be a frontend
7. **No other tech** - Don't suggest or mention alternative technologies. Work with Go and SQLite.

---

Read the README to learn more about what we are doing. This project is PowerPro. It is a headless api there is no frontend and there should not be a frontend. The core challenge of powerpro is keeping the code DRY. All of these powerlifters are doing things slightly differently but at the end of the day they're really doing very similar things. It is our job as programmers to manage this complexity in a re-usable way. The difficulty will come from managing the complexity related to that. When designing the entities and their relationship to eachother and the classes even in early phases it will be critical to design it such that we don't end up with spaghetti and you must keep your eye on the horizon. Don't preemptively add fields but you may preemptively add needed layers of abstraction that may well start as trivial / passthrough for initial implementations.

If you are working on a ticket, ERD, PRD, or roadmap Phase you must use sdlc.sh to manipulate those objects.

If you are told to work on a ticket and it is complete then you must use sdlc.sh to update that status.

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

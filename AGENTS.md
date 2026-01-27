# Agents

frontend uses bun 

use crumbler as ./crumbler

whenever you break down a crumb include as a last crumb a tech debt audit and cleanup crumb from the resulting work

### Hard Constraints (NEVER VIOLATE)

1. **Go only** - The language is Go. Period.
2. **SQLite only** - The database is SQLite. Never mention or consider other databases (no Postgres, MySQL, MongoDB, Redis, etc.)
3. **Follow the roadmap** - You MUST follow the existing roadmap phases
4. **Extend roadmap only after testing** - You may only extend the roadmap once comprehensive testing has been completed for current work
5. **No external paid services** - Never mention or use Auth0, Stripe, AWS services, or any external paid service
6. **frontend** - in the frontend folder
7. **No other tech** - Don't suggest or mention alternative technologies. Work with Go and SQLite.

---

Read the README to learn more about what we are doing. This project is PowerPro. It is a headless api there is no frontend and there should not be a frontend. The core challenge of powerpro is keeping the code DRY. All of these powerlifters are doing things slightly differently but at the end of the day they're really doing very similar things. It is our job as programmers to manage this complexity in a re-usable way. The difficulty will come from managing the complexity related to that. When designing the entities and their relationship to eachother and the classes even in early phases it will be critical to design it such that we don't end up with spaghetti and you must keep your eye on the horizon. Don't preemptively add fields but you may preemptively add needed layers of abstraction that may well start as trivial / passthrough for initial implementations.

Work is organized using crumbler - a simple task decomposition system. Use `./crumbler help` to learn more about the crumbler system.

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

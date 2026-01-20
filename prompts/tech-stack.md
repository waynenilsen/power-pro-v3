# Technology Stack

This document defines the technology stack and development practices used in this project. All agents working on this codebase must be aware of and follow these technology choices.

## Core Technologies

### Backend Language & Runtime
- **Go (Golang)**: The primary backend language
- **Vanilla Go**: Use standard library Go. Avoid heavy frameworks or ORMs
- **RPC-based System**: The system uses RPC (Remote Procedure Call) architecture for communication between frontend and backend

### Database & Data Access
- **SQLite**: The database system
- **sqlc**: Code generation tool for type-safe SQL queries. Use sqlc to generate Go code from SQL queries
- **goose**: Database migration tool. All database schema changes must be managed through goose migrations

### Migration Management
- All database schema changes must be done via goose migrations
- Migrations are versioned and tracked
- Never modify the database schema directly

## Testing Approach

### E2E Testing Only
This project uses an **unusual but deliberate testing strategy**: **End-to-end testing only**. There are no unit tests or integration tests in the traditional sense.

### Test Execution Model
1. **Dev Server Per Test**: Each test stands up its own development server
2. **Isolated Database**: Each dev server uses its own temporary SQLite database
3. **Isolated Port**: Each dev server runs on its own port
4. **Port Allocation**: Ports are randomly generated between **30000 and 60000**
5. **Collision Handling**: In the unlikely case of port collision, the system handles it appropriately

### Test Cleanup
After each test completes:
- The dev server process is shut down
- The temporary SQLite database file is deleted
- All resources are cleaned up

### Why This Approach?
- Tests the entire system as it will run in production
- Catches integration issues that unit tests might miss
- Ensures the dev server startup/shutdown process works correctly
- Validates the full request/response cycle

### Test-Driven Development (TDD) Process

**Critical**: Tests must use TDD, but in a specific way:

1. **Create the Code Under Test First**: Create the actual endpoint/handler structure
2. **Mock/Hardcode the Body**: Return a simple hardcoded value (e.g., empty list `[]` for a list endpoint)
3. **Write Tests**: Write tests that exercise the endpoint
4. **Tests Should Fail for Legitimate Reasons**: Tests should fail because the implementation doesn't match expectations (e.g., test expects non-empty list but gets `[]`), NOT because:
   - The endpoint doesn't exist
   - The endpoint isn't set up
   - Trivial setup issues

5. **Implement to Pass**: Once tests fail for the right reasons, implement the actual functionality to make tests pass

**Example**:
```go
// Step 1: Create endpoint with hardcoded response
func ListUsers() []User {
    return []User{}  // Hardcoded empty list
}

// Step 2: Write test
func TestListUsers(t *testing.T) {
    // Setup dev server, make request to /users endpoint
    users := listUsers()
    assert.NotEmpty(users)  // This will fail, but for the RIGHT reason
}

// Step 3: Implement actual functionality
func ListUsers() []User {
    // Now implement actual database query
    return db.QueryUsers()
}
```

**Key Principle**: Tests should fail for **legitimate business logic reasons**, not trivial setup/configuration issues. The endpoint structure must exist and be callable before writing tests.

## Development Practices

### Code Generation
- Use **sqlc** to generate type-safe database access code from SQL
- Never write raw database queries without sqlc generation
- Keep SQL queries in `.sql` files for sqlc to process

### Database Migrations
- Use **goose** for all schema changes
- Create migration files for every database change
- Migrations are run automatically during test setup

### Schema Change Process

**Critical**: When making tickets that involve schema changes, follow this process:

1. **Separate Ticket/Commit**: Schema changes must be in a **separate ticket and commit** from the code that uses them
2. **Migration Required**: Every schema change must have a corresponding goose migration
3. **Production Considerations**: Migrations must consider:
   - What is currently checked in to production
   - Migration strategies for updating existing data
   - Backward compatibility where necessary
4. **Data Migration**: Include data transformation/migration logic in the migration file
5. **Performance**: Don't worry about migration performance yet (optimize later if needed)

**Example Workflow**:
- Ticket 001: Add user table (schema change + migration)
- Ticket 002: Implement user registration API (uses the schema from Ticket 001)

This separation ensures:
- Schema changes can be reviewed independently
- Migrations are tested before code that depends on them
- Clear separation of concerns
- Easier rollback if needed

### RPC Communication
- The system uses RPC for frontend/backend communication
- This affects authentication and session management (see `auth.md`)
- Bearer tokens are used for authentication (session ID sent as bearer token)

## External Services & Dependencies

### External Services Policy
- **External services are explicitly denied until otherwise noted**
- The entire system must run on a laptop after checkout
- No external API dependencies (payment processors, email services, etc.) unless explicitly approved
- No cloud services required for local development
- All functionality must be self-contained

### Local Development Requirements
- **Everything must run on a laptop on checkout**
- No external infrastructure dependencies
- All services must be runnable locally
- SQLite database (file-based, no external database server)

### Setup Script
- **A setup script must exist** that gets everything set up including:
  - Installing dependencies
  - Setting up the development environment
  - Running database migrations
  - Any other initialization steps
- The setup script should enable a developer to run `./setup.sh` (or equivalent) and have a working environment
- Migrations must be run automatically as part of setup

## What We Don't Use

### Avoid These Technologies
- **Heavy frameworks**: No Gin, Echo, or similar web frameworks (vanilla Go only)
- **ORMs**: No GORM, Ent, or other ORMs (use sqlc instead)
- **Unit testing frameworks**: No traditional unit tests (e2e only)
- **JWT tokens**: Not used due to RPC nature and revocation requirements (see `auth.md`)
- **Cookie-based sessions**: Sessions are managed via bearer tokens (see `auth.md`)
- **External services**: No external APIs or cloud services unless explicitly approved

## File Organization

### Database Code
- SQL queries: `.sql` files (processed by sqlc)
- Generated code: Generated by sqlc (typically in a `db/` or `gen/` directory)
- Migrations: Goose migration files (typically in `migrations/` directory)

### RPC Definitions
- RPC service definitions and handlers
- Request/response types

### Test Files
- E2E test files that spin up dev servers
- Test utilities for port allocation and cleanup

## Important Considerations

### Port Management
- Tests must handle random port allocation (30000-60000)
- Port collision is rare but possible - handle gracefully
- Each test must clean up its allocated port

### Database Isolation
- Each test gets a fresh database
- No test should affect another test's database
- Cleanup must be thorough (delete temp DB files)

### RPC Architecture
- Frontend and backend are separate
- Communication happens via RPC calls
- This affects how authentication tokens are passed (bearer tokens, not cookies)

## Integration with Other Systems

This tech stack document should be referenced when:
- Creating ERDs (Engineering Requirements Documents)
- Breaking down ERDs into tickets
- Implementing features
- Writing tests
- Making architectural decisions

Agents should always consider these technology constraints when proposing solutions or implementing features.

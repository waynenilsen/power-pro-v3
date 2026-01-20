# Ticket System

An in-filesystem ticket system for managing tasks and work items.

## Directory Structure

```
tickets/
├── todo/          # Tickets that need to be done
├── in-progress/   # Tickets currently being worked on
├── done/          # Completed tickets
└── not-doing/     # Tickets that are cancelled or will not be done
```

## Directory Descriptions

### `tickets/`
The primary directory containing all ticket subdirectories.

### `tickets/todo/`
Tickets that are planned but not yet started. These represent work items that need to be completed.

### `tickets/in-progress/`
Tickets that are currently being actively worked on. Move tickets here when work begins.

### `tickets/done/`
Completed tickets. Move tickets here when they are finished.

### `tickets/not-doing/`
Tickets that have been cancelled, deferred indefinitely, or decided against. Use this for tickets that will not be completed.

## Workflow

1. **Create**: New tickets start in `tickets/todo/`
2. **Start**: Move tickets to `tickets/in-progress/` when work begins
3. **Complete**: Move tickets to `tickets/done/` when finished
4. **Cancel**: Move tickets to `tickets/not-doing/` if they won't be completed

## Ticket File Format

### File Naming Convention

Tickets must follow a specific naming format to enable lexicographical sorting by ticket number:

**Format**: `NNN-description.md`

Where:
- `NNN` is a zero-padded ticket number (e.g., `001`, `002`, `010`, `100`)
- `description` is a short, descriptive name using hyphens or underscores
- File extension is `.md` (or other appropriate format)

**Examples**:
- `001-initial-setup.md`
- `002-add-authentication.md`
- `010-refactor-api.md`
- `100-deploy-production.md`

**Why zero-padding?**
Zero-padding ensures that when files are sorted lexicographically (alphabetically), they are also sorted numerically. Without zero-padding, `10-ticket.md` would sort before `2-ticket.md`, which is incorrect.

### File Content

Each ticket file should contain:
- Title/Description
- Creation date
- Any relevant details, notes, or context

**Important**: Ticket files must NOT contain status information. Status is always implied by the directory location (`todo/`, `in-progress/`, `done/`, or `not-doing/`). Including status in the file body would create redundancy and potential inconsistencies.

## Schema Changes

### Separate Tickets for Schema Changes

**Critical Rule**: When a ticket involves database schema changes, the schema change must be in a **separate ticket and commit** from the code that uses the schema.

### Process
1. **Schema Change Ticket**: Create a ticket for the schema change itself
   - Includes the schema modification
   - Includes the goose migration file
   - Migration must consider what's currently in production
   - Migration must include data migration strategies
   - Performance optimization is not required (yet)

2. **Implementation Ticket**: Create a separate ticket for code that uses the schema
   - Depends on the schema change ticket
   - Uses the new schema structure
   - Implements the feature/functionality

### Example
- **Ticket 001**: Add users table (schema change + migration)
  - Creates `users` table with email, password_hash, name columns
  - Includes goose migration
  - Handles data migration if needed
  
- **Ticket 002**: Implement user registration API (blocked by 001)
  - Uses the `users` table created in Ticket 001
  - Implements registration logic

### Why Separate?
- Schema changes can be reviewed independently
- Migrations are tested before dependent code
- Clear separation of concerns
- Easier rollback if needed
- Better traceability

See `tech-stack.md` for more details on migration requirements.

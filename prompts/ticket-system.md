# Work Management with Crumbler

Work is organized using crumbler - a simple task decomposition system. The filesystem IS the state - existence means work to do, deletion means done. Work depth-first: complete children before parents.

## Directory Structure

Crumbs are directories with README.md files organized in a tree structure:

```
.crumbler/                           # Project root (auto-created)
├── README.md                        # Root crumb
├── 01-setup/                        # First child crumb
│   ├── README.md                    # Task instructions
│   └── 01-database/                 # Nested crumb
│       └── README.md
└── 02-features/
    └── README.md
```

## How Crumbler Works

### Core Concepts

- **Crumb**: A directory containing a `README.md` file with task instructions
- **State**: The filesystem IS the state - existence means work to do, deletion means done
- **Workflow**: Depth-first - complete children before parents
- **Naming**: Use zero-padded numbers (01, 02, etc.) for proper lexicographical sorting

### Basic Commands

```bash
# Get AI instructions for current work
./crumbler prompt

# Create a new crumb (sub-task)
./crumbler create "Task Name"

# Mark current crumb as done (deletes the crumb directory)
./crumbler delete

# View crumb tree and status
./crumbler status

# Get help
./crumbler help
```

## Workflow

1. **Get Instructions**: Run `./crumbler prompt` to get AI instructions for what to do next
2. **Do the Work**: Follow the instructions in the current crumb's README.md
3. **Create Sub-tasks**: If work needs decomposition, use `./crumbler create "Sub-task Name"` to create child crumbs
4. **Complete**: When a crumb is done, run `./crumbler delete` to mark it complete
5. **Repeat**: Run `./crumbler prompt` again to get next instructions

## Crumb README Format

Each crumb directory contains a `README.md` file with:

- **Title**: Clear description of the work
- **Description**: What needs to be done
- **Context**: Why this work exists, dependencies, background
- **Acceptance Criteria**: How to know when it's done
- **Technical Notes**: Implementation approach, constraints, design decisions

**Important**: Crumb README files should NOT contain status information. Status is implied by existence (work to do) or deletion (done).

## Schema Changes

### Separate Crumbs for Schema Changes

**Critical Rule**: When work involves database schema changes, the schema change must be in a **separate crumb and commit** from the code that uses the schema.

### Process

1. **Schema Change Crumb**: Create a crumb for the schema change itself
   - Includes the schema modification
   - Includes the goose migration file
   - Migration must consider what's currently in production
   - Migration must include data migration strategies
   - Performance optimization is not required (yet)

2. **Implementation Crumb**: Create a separate crumb for code that uses the schema
   - Depends on the schema change crumb
   - Uses the new schema structure
   - Implements the feature/functionality

### Example

- **Crumb: Create users table schema**
  - Creates `users` table with email, password_hash, name columns
  - Includes goose migration
  - Handles data migration if needed
  
- **Crumb: Implement user registration API** (depends on schema crumb)
  - Uses the `users` table created in schema crumb
  - Implements registration logic

### Why Separate?

- Schema changes can be reviewed independently
- Migrations are tested before dependent code
- Clear separation of concerns
- Easier rollback if needed
- Better traceability

See `tech-stack.md` for more details on migration requirements.

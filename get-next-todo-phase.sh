#!/usr/bin/env bash
#
# get-next-todo-phase.sh - Get the next phase from phases/todo/
#
# Usage:
#   ./get-next-todo-phase.sh
#
# Returns:
#   - Prints the path to the next phase (lexicographically first in todo/)
#   - Exits with code 0 if phase found, 1 if no phases available

# Get absolute path to this script's directory
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# Call the generic get-next-phase script with 'todo' state
exec "$SCRIPT_DIR/get-next-phase.sh" todo

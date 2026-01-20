#!/usr/bin/env bash
#
# get-next-todo-prd.sh - Get the next PRD from prds/todo/
#
# Usage:
#   ./get-next-todo-prd.sh
#
# Returns:
#   - Prints the path to the next PRD (lexicographically first in todo/)
#   - Exits with code 0 if PRD found, 1 if no PRDs available

# Get absolute path to this script's directory
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# Call the generic get-next-prd script with 'todo' state
exec "$SCRIPT_DIR/get-next-prd.sh" todo

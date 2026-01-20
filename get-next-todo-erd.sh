#!/usr/bin/env bash
#
# get-next-todo-erd.sh - Get the next ERD from erds/todo/
#
# Usage:
#   ./get-next-todo-erd.sh
#
# Returns:
#   - Prints the path to the next ERD (lexicographically first in todo/)
#   - Exits with code 0 if ERD found, 1 if no ERDs available

# Get absolute path to this script's directory
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# Call the generic get-next-erd script with 'todo' state
exec "$SCRIPT_DIR/get-next-erd.sh" todo

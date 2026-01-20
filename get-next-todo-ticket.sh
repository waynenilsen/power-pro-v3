#!/usr/bin/env bash
#
# get-next-todo-ticket.sh - Get the next ticket from tickets/todo/
#
# Usage:
#   ./get-next-todo-ticket.sh
#
# Returns:
#   - Prints the path to the next ticket (lexicographically first in todo/)
#   - Exits with code 0 if ticket found, 1 if no tickets available

# Get absolute path to this script's directory
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# Call the generic get-next-ticket script with 'todo' state
exec "$SCRIPT_DIR/get-next-ticket.sh" todo

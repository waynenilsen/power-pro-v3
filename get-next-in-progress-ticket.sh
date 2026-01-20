#!/usr/bin/env bash
#
# get-next-in-progress-ticket.sh - Get the next ticket from tickets/in-progress/
#
# Usage:
#   ./get-next-in-progress-ticket.sh
#
# Returns:
#   - Prints the path to the next ticket (lexicographically first in in-progress/)
#   - Exits with code 0 if ticket found, 1 if no tickets available

# Get absolute path to this script's directory
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# Call the generic get-next-ticket script with 'in-progress' state
exec "$SCRIPT_DIR/get-next-ticket.sh" in-progress

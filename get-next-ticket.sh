#!/usr/bin/env bash
#
# get-next-ticket - Get the next ticket from tickets/[state]/
#
# Usage:
#   ./get-next-ticket [state]
#
# Arguments:
#   state - The ticket state directory (todo, in-progress, done, not-doing)
#
# Returns:
#   - Prints the path to the next ticket (lexicographically first in [state]/)
#   - Exits with code 0 if ticket found, 1 if no tickets available or invalid state

# Get absolute path to this script's directory
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# Validate state argument
if [ -z "$1" ]; then
  echo "Error: state argument required" >&2
  echo "Usage: $0 [state]" >&2
  echo "Valid states: todo, in-progress, done, not-doing" >&2
  exit 1
fi

STATE="$1"
TICKETS_STATE_DIR="$SCRIPT_DIR/tickets/$STATE"

# Check if state directory exists
if [ ! -d "$TICKETS_STATE_DIR" ]; then
  echo "Error: tickets/$STATE/ directory not found" >&2
  exit 1
fi

# Find the first ticket file (lexicographically sorted)
# Tickets follow format: NNN-description.md
NEXT_TICKET=$(find "$TICKETS_STATE_DIR" -maxdepth 1 -type f -name "*.md" | sort | head -n 1)

# Check if a ticket was found
if [ -z "$NEXT_TICKET" ]; then
  echo "No tickets found in tickets/$STATE/" >&2
  exit 1
fi

# Output the path
echo "$NEXT_TICKET"

exit 0

#!/usr/bin/env bash
#
# get-next-phase - Get the next phase from phases/[state]/
#
# Usage:
#   ./get-next-phase.sh [state]
#
# Arguments:
#   state - The phase state directory (todo, in-progress, done, not-doing)
#
# Returns:
#   - Prints the path to the next phase (lexicographically first in [state]/)
#   - Exits with code 0 if phase found, 1 if no phases available or invalid state

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
PHASES_STATE_DIR="$SCRIPT_DIR/phases/$STATE"

# Check if state directory exists
if [ ! -d "$PHASES_STATE_DIR" ]; then
  echo "Error: phases/$STATE/ directory not found" >&2
  exit 1
fi

# Find the first phase file (lexicographically sorted)
# Phases follow format: NNN-description.md
NEXT_PHASE=$(find "$PHASES_STATE_DIR" -maxdepth 1 -type f -name "*.md" | sort | head -n 1)

# Check if a phase was found
if [ -z "$NEXT_PHASE" ]; then
  echo "No phases found in phases/$STATE/" >&2
  exit 1
fi

# Output the path
echo "$NEXT_PHASE"

exit 0

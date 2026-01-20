#!/usr/bin/env bash
#
# get-next-prd - Get the next PRD from prds/[state]/
#
# Usage:
#   ./get-next-prd.sh [state]
#
# Arguments:
#   state - The PRD state directory (todo, in-progress, done, not-doing)
#
# Returns:
#   - Prints the path to the next PRD (lexicographically first in [state]/)
#   - Exits with code 0 if PRD found, 1 if no PRDs available or invalid state

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
PRDS_STATE_DIR="$SCRIPT_DIR/prds/$STATE"

# Check if state directory exists
if [ ! -d "$PRDS_STATE_DIR" ]; then
  echo "Error: prds/$STATE/ directory not found" >&2
  exit 1
fi

# Find the first PRD file (lexicographically sorted)
# PRDs follow format: NNN-description.md
NEXT_PRD=$(find "$PRDS_STATE_DIR" -maxdepth 1 -type f -name "*.md" | sort | head -n 1)

# Check if a PRD was found
if [ -z "$NEXT_PRD" ]; then
  echo "No PRDs found in prds/$STATE/" >&2
  exit 1
fi

# Output the path
echo "$NEXT_PRD"

exit 0

#!/usr/bin/env bash
#
# get-next-erd - Get the next ERD from erds/[state]/
#
# Usage:
#   ./get-next-erd.sh [state]
#
# Arguments:
#   state - The ERD state directory (todo, in-progress, done, not-doing)
#
# Returns:
#   - Prints the path to the next ERD (lexicographically first in [state]/)
#   - Exits with code 0 if ERD found, 1 if no ERDs available or invalid state

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
ERDS_STATE_DIR="$SCRIPT_DIR/erds/$STATE"

# Check if state directory exists
if [ ! -d "$ERDS_STATE_DIR" ]; then
  echo "Error: erds/$STATE/ directory not found" >&2
  exit 1
fi

# Find the first ERD file (lexicographically sorted)
# ERDs follow format: NNN-description.md
NEXT_ERD=$(find "$ERDS_STATE_DIR" -maxdepth 1 -type f -name "*.md" | sort | head -n 1)

# Check if an ERD was found
if [ -z "$NEXT_ERD" ]; then
  echo "No ERDs found in erds/$STATE/" >&2
  exit 1
fi

# Output the path
echo "$NEXT_ERD"

exit 0

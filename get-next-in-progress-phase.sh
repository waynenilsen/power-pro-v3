#!/usr/bin/env bash
#
# get-next-in-progress-phase.sh - Get the next phase from phases/in-progress/
#
# Usage:
#   ./get-next-in-progress-phase.sh
#
# Returns:
#   - Prints the path to the next phase (lexicographically first in in-progress/)
#   - Exits with code 0 if phase found, 1 if no phases available

# Get absolute path to this script's directory
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# Call the generic get-next-phase script with 'in-progress' state
exec "$SCRIPT_DIR/get-next-phase.sh" in-progress

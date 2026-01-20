#!/usr/bin/env bash
#
# get-next-in-progress-prd.sh - Get the next PRD from prds/in-progress/
#
# Usage:
#   ./get-next-in-progress-prd.sh
#
# Returns:
#   - Prints the path to the next PRD (lexicographically first in in-progress/)
#   - Exits with code 0 if PRD found, 1 if no PRDs available

# Get absolute path to this script's directory
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# Call the generic get-next-prd script with 'in-progress' state
exec "$SCRIPT_DIR/get-next-prd.sh" in-progress

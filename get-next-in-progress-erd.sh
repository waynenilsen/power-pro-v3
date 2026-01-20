#!/usr/bin/env bash
#
# get-next-in-progress-erd.sh - Get the next ERD from erds/in-progress/
#
# Usage:
#   ./get-next-in-progress-erd.sh
#
# Returns:
#   - Prints the path to the next ERD (lexicographically first in in-progress/)
#   - Exits with code 0 if ERD found, 1 if no ERDs available

# Get absolute path to this script's directory
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# Call the generic get-next-erd script with 'in-progress' state
exec "$SCRIPT_DIR/get-next-erd.sh" in-progress

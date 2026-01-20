#!/usr/bin/env bash
#
# move-prd.sh - Move a PRD between directories
#
# Usage:
#   ./move-prd.sh <prd> <destination>
#
# PRD can be specified as:
#   - Full path: /path/to/prds/todo/001-prd.md
#   - Relative path: prds/todo/001-prd.md
#   - Filename: 001-prd.md (searches all PRD directories)
#   - PRD number: 001 (searches all PRD directories for matching PRD)
#
# Destination can be:
#   - todo
#   - in-progress
#   - done
#   - not-doing
#
# Examples:
#   ./move-prd.sh 001-prd.md in-progress
#   ./move-prd.sh 001 done
#   ./move-prd.sh prds/todo/001-prd.md in-progress

# Get absolute path to this script's directory
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PRDS_DIR="$SCRIPT_DIR/prds"

# Check arguments
if [ $# -lt 2 ]; then
  echo "Usage: $0 <prd> <destination>" >&2
  echo "" >&2
  echo "PRD can be: full path, relative path, filename, or PRD number" >&2
  echo "Destination can be: todo, in-progress, done, not-doing" >&2
  exit 1
fi

PRD_INPUT="$1"
DESTINATION="$2"

# Validate destination
case "$DESTINATION" in
  todo|in-progress|done|not-doing)
    DEST_DIR="$PRDS_DIR/$DESTINATION"
    ;;
  *)
    echo "Error: Invalid destination '$DESTINATION'" >&2
    echo "Destination must be: todo, in-progress, done, or not-doing" >&2
    exit 1
    ;;
esac

# Function to find PRD file
find_prd() {
  local search="$1"
  local found=""
  
  # If it's already a full path and exists, use it
  if [ -f "$search" ]; then
    # Check if it's within prds directory
    if [[ "$(realpath "$search" 2>/dev/null)" == "$(realpath "$PRDS_DIR")"* ]]; then
      echo "$(realpath "$search")"
      return 0
    fi
  fi
  
  # If it's a relative path, try resolving it
  if [[ "$search" == *"/"* ]]; then
    local resolved=""
    if [[ "$search" == /* ]]; then
      resolved="$search"
    else
      resolved="$SCRIPT_DIR/$search"
    fi
    if [ -f "$resolved" ] && [[ "$(realpath "$resolved" 2>/dev/null)" == "$(realpath "$PRDS_DIR")"* ]]; then
      echo "$(realpath "$resolved")"
      return 0
    fi
  fi
  
  # Search all PRD directories
  for dir in todo in-progress done not-doing; do
    local dir_path="$PRDS_DIR/$dir"
    if [ ! -d "$dir_path" ]; then
      continue
    fi
    
    # Try exact filename match
    if [ -f "$dir_path/$search" ]; then
      echo "$(realpath "$dir_path/$search")"
      return 0
    fi
    
    # Try PRD number match (NNN-*.md)
    if [[ "$search" =~ ^[0-9]+$ ]]; then
      local padded_search=$(printf "%03d" "$search" 2>/dev/null || echo "$search")
      local matches=$(find "$dir_path" -maxdepth 1 -type f -name "${padded_search}-*.md" 2>/dev/null | head -n 1)
      if [ -n "$matches" ]; then
        echo "$(realpath "$matches")"
        return 0
      fi
    fi
    
    # Try partial match (contains the search string)
    local matches=$(find "$dir_path" -maxdepth 1 -type f -name "*${search}*" 2>/dev/null | head -n 1)
    if [ -n "$matches" ]; then
      echo "$(realpath "$matches")"
      return 0
    fi
  done
  
  return 1
}

# Find the PRD
PRD_PATH=$(find_prd "$PRD_INPUT")

if [ -z "$PRD_PATH" ] || [ ! -f "$PRD_PATH" ]; then
  echo "Error: PRD not found: $PRD_INPUT" >&2
  exit 1
fi

# Get the filename
PRD_FILENAME=$(basename "$PRD_PATH")

# Check if already in destination
CURRENT_DIR=$(dirname "$PRD_PATH")
if [ "$(realpath "$CURRENT_DIR")" = "$(realpath "$DEST_DIR")" ]; then
  echo "PRD is already in $DESTINATION/: $PRD_FILENAME" >&2
  exit 0
fi

# Ensure destination directory exists
mkdir -p "$DEST_DIR"

# Move the PRD
mv "$PRD_PATH" "$DEST_DIR/$PRD_FILENAME"

if [ $? -eq 0 ]; then
  echo "Moved $PRD_FILENAME to $DESTINATION/"
  exit 0
else
  echo "Error: Failed to move PRD" >&2
  exit 1
fi

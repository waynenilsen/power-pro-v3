#!/usr/bin/env bash
#
# move-erd.sh - Move an ERD between directories
#
# Usage:
#   ./move-erd.sh <erd> <destination>
#
# ERD can be specified as:
#   - Full path: /path/to/erds/todo/001-erd.md
#   - Relative path: erds/todo/001-erd.md
#   - Filename: 001-erd.md (searches all ERD directories)
#   - ERD number: 001 (searches all ERD directories for matching ERD)
#
# Destination can be:
#   - todo
#   - in-progress
#   - done
#   - not-doing
#
# Examples:
#   ./move-erd.sh 001-erd.md in-progress
#   ./move-erd.sh 001 done
#   ./move-erd.sh erds/todo/001-erd.md in-progress

# Get absolute path to this script's directory
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ERDS_DIR="$SCRIPT_DIR/erds"

# Check arguments
if [ $# -lt 2 ]; then
  echo "Usage: $0 <erd> <destination>" >&2
  echo "" >&2
  echo "ERD can be: full path, relative path, filename, or ERD number" >&2
  echo "Destination can be: todo, in-progress, done, not-doing" >&2
  exit 1
fi

ERD_INPUT="$1"
DESTINATION="$2"

# Validate destination
case "$DESTINATION" in
  todo|in-progress|done|not-doing)
    DEST_DIR="$ERDS_DIR/$DESTINATION"
    ;;
  *)
    echo "Error: Invalid destination '$DESTINATION'" >&2
    echo "Destination must be: todo, in-progress, done, or not-doing" >&2
    exit 1
    ;;
esac

# Function to find ERD file
find_erd() {
  local search="$1"
  local found=""
  
  # If it's already a full path and exists, use it
  if [ -f "$search" ]; then
    # Check if it's within erds directory
    if [[ "$(realpath "$search" 2>/dev/null)" == "$(realpath "$ERDS_DIR")"* ]]; then
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
    if [ -f "$resolved" ] && [[ "$(realpath "$resolved" 2>/dev/null)" == "$(realpath "$ERDS_DIR")"* ]]; then
      echo "$(realpath "$resolved")"
      return 0
    fi
  fi
  
  # Search all ERD directories
  for dir in todo in-progress done not-doing; do
    local dir_path="$ERDS_DIR/$dir"
    if [ ! -d "$dir_path" ]; then
      continue
    fi
    
    # Try exact filename match
    if [ -f "$dir_path/$search" ]; then
      echo "$(realpath "$dir_path/$search")"
      return 0
    fi
    
    # Try ERD number match (NNN-*.md)
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

# Find the ERD
ERD_PATH=$(find_erd "$ERD_INPUT")

if [ -z "$ERD_PATH" ] || [ ! -f "$ERD_PATH" ]; then
  echo "Error: ERD not found: $ERD_INPUT" >&2
  exit 1
fi

# Get the filename
ERD_FILENAME=$(basename "$ERD_PATH")

# Check if already in destination
CURRENT_DIR=$(dirname "$ERD_PATH")
if [ "$(realpath "$CURRENT_DIR")" = "$(realpath "$DEST_DIR")" ]; then
  echo "ERD is already in $DESTINATION/: $ERD_FILENAME" >&2
  exit 0
fi

# Ensure destination directory exists
mkdir -p "$DEST_DIR"

# Move the ERD
mv "$ERD_PATH" "$DEST_DIR/$ERD_FILENAME"

if [ $? -eq 0 ]; then
  echo "Moved $ERD_FILENAME to $DESTINATION/"
  exit 0
else
  echo "Error: Failed to move ERD" >&2
  exit 1
fi

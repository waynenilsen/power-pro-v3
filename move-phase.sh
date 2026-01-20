#!/usr/bin/env bash
#
# move-phase.sh - Move a phase between directories
#
# Usage:
#   ./move-phase.sh <phase> <destination>
#
# Phase can be specified as:
#   - Full path: /path/to/phases/todo/001-phase.md
#   - Relative path: phases/todo/001-phase.md
#   - Filename: 001-phase.md (searches all phase directories)
#   - Phase number: 001 (searches all phase directories for matching phase)
#
# Destination can be:
#   - todo
#   - in-progress
#   - done
#   - not-doing
#
# Examples:
#   ./move-phase.sh 001-phase.md in-progress
#   ./move-phase.sh 001 done
#   ./move-phase.sh phases/todo/001-phase.md in-progress

# Get absolute path to this script's directory
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PHASES_DIR="$SCRIPT_DIR/phases"

# Check arguments
if [ $# -lt 2 ]; then
  echo "Usage: $0 <phase> <destination>" >&2
  echo "" >&2
  echo "Phase can be: full path, relative path, filename, or phase number" >&2
  echo "Destination can be: todo, in-progress, done, not-doing" >&2
  exit 1
fi

PHASE_INPUT="$1"
DESTINATION="$2"

# Validate destination
case "$DESTINATION" in
  todo|in-progress|done|not-doing)
    DEST_DIR="$PHASES_DIR/$DESTINATION"
    ;;
  *)
    echo "Error: Invalid destination '$DESTINATION'" >&2
    echo "Destination must be: todo, in-progress, done, or not-doing" >&2
    exit 1
    ;;
esac

# Function to find phase file
find_phase() {
  local search="$1"
  local found=""
  
  # If it's already a full path and exists, use it
  if [ -f "$search" ]; then
    # Check if it's within phases directory
    if [[ "$(realpath "$search" 2>/dev/null)" == "$(realpath "$PHASES_DIR")"* ]]; then
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
    if [ -f "$resolved" ] && [[ "$(realpath "$resolved" 2>/dev/null)" == "$(realpath "$PHASES_DIR")"* ]]; then
      echo "$(realpath "$resolved")"
      return 0
    fi
  fi
  
  # Search all phase directories
  for dir in todo in-progress done not-doing; do
    local dir_path="$PHASES_DIR/$dir"
    if [ ! -d "$dir_path" ]; then
      continue
    fi
    
    # Try exact filename match
    if [ -f "$dir_path/$search" ]; then
      echo "$(realpath "$dir_path/$search")"
      return 0
    fi
    
    # Try phase number match (NNN-*.md)
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

# Find the phase
PHASE_PATH=$(find_phase "$PHASE_INPUT")

if [ -z "$PHASE_PATH" ] || [ ! -f "$PHASE_PATH" ]; then
  echo "Error: Phase not found: $PHASE_INPUT" >&2
  exit 1
fi

# Get the filename
PHASE_FILENAME=$(basename "$PHASE_PATH")

# Check if already in destination
CURRENT_DIR=$(dirname "$PHASE_PATH")
if [ "$(realpath "$CURRENT_DIR")" = "$(realpath "$DEST_DIR")" ]; then
  echo "Phase is already in $DESTINATION/: $PHASE_FILENAME" >&2
  exit 0
fi

# Ensure destination directory exists
mkdir -p "$DEST_DIR"

# Move the phase
mv "$PHASE_PATH" "$DEST_DIR/$PHASE_FILENAME"

if [ $? -eq 0 ]; then
  echo "Moved $PHASE_FILENAME to $DESTINATION/"
  exit 0
else
  echo "Error: Failed to move phase" >&2
  exit 1
fi

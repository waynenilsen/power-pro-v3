#!/usr/bin/env bash
#
# move-ticket.sh - Move a ticket between directories
#
# Usage:
#   ./move-ticket.sh <ticket> <destination>
#
# Ticket can be specified as:
#   - Full path: /path/to/tickets/todo/001-ticket.md
#   - Relative path: tickets/todo/001-ticket.md
#   - Filename: 001-ticket.md (searches all ticket directories)
#   - Ticket number: 001 (searches all ticket directories for matching ticket)
#
# Destination can be:
#   - todo
#   - in-progress
#   - done
#   - not-doing
#
# Examples:
#   ./move-ticket.sh 001-ticket.md in-progress
#   ./move-ticket.sh 001 done
#   ./move-ticket.sh tickets/todo/001-ticket.md in-progress

# Get absolute path to this script's directory
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
TICKETS_DIR="$SCRIPT_DIR/tickets"

# Check arguments
if [ $# -lt 2 ]; then
  echo "Usage: $0 <ticket> <destination>" >&2
  echo "" >&2
  echo "Ticket can be: full path, relative path, filename, or ticket number" >&2
  echo "Destination can be: todo, in-progress, done, not-doing" >&2
  exit 1
fi

TICKET_INPUT="$1"
DESTINATION="$2"

# Validate destination
case "$DESTINATION" in
  todo|in-progress|done|not-doing)
    DEST_DIR="$TICKETS_DIR/$DESTINATION"
    ;;
  *)
    echo "Error: Invalid destination '$DESTINATION'" >&2
    echo "Destination must be: todo, in-progress, done, or not-doing" >&2
    exit 1
    ;;
esac

# Function to find ticket file
find_ticket() {
  local search="$1"
  local found=""
  
  # If it's already a full path and exists, use it
  if [ -f "$search" ]; then
    # Check if it's within tickets directory
    if [[ "$(realpath "$search" 2>/dev/null)" == "$(realpath "$TICKETS_DIR")"* ]]; then
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
    if [ -f "$resolved" ] && [[ "$(realpath "$resolved" 2>/dev/null)" == "$(realpath "$TICKETS_DIR")"* ]]; then
      echo "$(realpath "$resolved")"
      return 0
    fi
  fi
  
  # Search all ticket directories
  for dir in todo in-progress done not-doing; do
    local dir_path="$TICKETS_DIR/$dir"
    if [ ! -d "$dir_path" ]; then
      continue
    fi
    
    # Try exact filename match
    if [ -f "$dir_path/$search" ]; then
      echo "$(realpath "$dir_path/$search")"
      return 0
    fi
    
    # Try ticket number match (NNN-*.md)
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

# Find the ticket
TICKET_PATH=$(find_ticket "$TICKET_INPUT")

if [ -z "$TICKET_PATH" ] || [ ! -f "$TICKET_PATH" ]; then
  echo "Error: Ticket not found: $TICKET_INPUT" >&2
  exit 1
fi

# Get the filename
TICKET_FILENAME=$(basename "$TICKET_PATH")

# Check if already in destination
CURRENT_DIR=$(dirname "$TICKET_PATH")
if [ "$(realpath "$CURRENT_DIR")" = "$(realpath "$DEST_DIR")" ]; then
  echo "Ticket is already in $DESTINATION/: $TICKET_FILENAME" >&2
  exit 0
fi

# Ensure destination directory exists
mkdir -p "$DEST_DIR"

# Move the ticket
mv "$TICKET_PATH" "$DEST_DIR/$TICKET_FILENAME"

if [ $? -eq 0 ]; then
  echo "Moved $TICKET_FILENAME to $DESTINATION/"
  exit 0
else
  echo "Error: Failed to move ticket" >&2
  exit 1
fi

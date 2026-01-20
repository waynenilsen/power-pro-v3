#!/usr/bin/env bash
#
# sdlc.sh - SDLC Management Tool
#
# Unified entrypoint for managing tickets, ERDs, PRDs, and phases
#
# Usage:
#   ./sdlc.sh <command> [options]
#
# Commands:
#   get-next    Get the next item from a state directory
#   move        Move an item between state directories
#   list        List items in a state directory
#   help        Show help information
#
# Examples:
#   ./sdlc.sh get-next ticket todo
#   ./sdlc.sh move ticket 001 in-progress
#   ./sdlc.sh list erd todo
#   ./sdlc.sh help

# Get absolute path to this script's directory
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Valid entity types
VALID_ENTITIES=("ticket" "erd" "prd" "phase")
VALID_STATES=("todo" "in-progress" "done" "not-doing")

# Show usage information
show_usage() {
  cat << EOF
${BLUE}SDLC Management Tool${NC}

Usage: $0 <command> [options]

Commands:
  get-next <entity> <state>    Get the next item from a state directory
  move <entity> <item> <dest>   Move an item between state directories
  list <entity> <state>         List items in a state directory
  help [command]                Show help information

Entity Types:
  ticket                        Work tickets
  erd                           Engineering Requirements Documents
  prd                           Product Requirements Documents
  phase                         Phase documents

States:
  todo                          Items to be done
  in-progress                   Items currently being worked on
  done                          Completed items
  not-doing                     Cancelled or deferred items

Examples:
  $0 get-next ticket todo
  $0 get-next erd in-progress
  $0 move ticket 001 in-progress
  $0 move erd 001-test-erd.md done
  $0 list prd todo
  $0 list phase in-progress
  $0 help get-next
  $0 help move

EOF
}

# Show help for a specific command
show_command_help() {
  local cmd="$1"
  case "$cmd" in
    get-next)
      cat << EOF
${BLUE}get-next${NC} - Get the next item from a state directory

Usage: $0 get-next <entity> <state>

Arguments:
  entity    Entity type (ticket, erd, prd, phase)
  state     State directory (todo, in-progress, done, not-doing)

Returns:
  Prints the path to the next item (lexicographically first)
  Exits with code 0 if item found, 1 if no items available

Examples:
  $0 get-next ticket todo
  $0 get-next erd in-progress
  $0 get-next prd done

EOF
      ;;
    move)
      cat << EOF
${BLUE}move${NC} - Move an item between state directories

Usage: $0 move <entity> <item> <destination>

Arguments:
  entity       Entity type (ticket, erd, prd, phase)
  item         Item identifier (filename, number, or path)
  destination  Destination state (todo, in-progress, done, not-doing)

Item can be specified as:
  - Filename:  001-item.md (searches all directories)
  - Number:    001 (searches for NNN-*.md)
  - Full path: /path/to/item.md
  - Rel path:  tickets/todo/001-item.md

Examples:
  $0 move ticket 001 in-progress
  $0 move erd 001-test-erd.md done
  $0 move prd 001 done
  $0 move phase tickets/todo/001-phase.md in-progress

EOF
      ;;
    list)
      cat << EOF
${BLUE}list${NC} - List items in a state directory

Usage: $0 list <entity> <state>

Arguments:
  entity    Entity type (ticket, erd, prd, phase)
  state     State directory (todo, in-progress, done, not-doing)

Returns:
  Lists all items in the specified state directory, sorted lexicographically

Examples:
  $0 list ticket todo
  $0 list erd in-progress
  $0 list prd done

EOF
      ;;
    *)
      echo "Unknown command: $cmd" >&2
      show_usage
      exit 1
      ;;
  esac
}

# Validate entity type
validate_entity() {
  local entity="$1"
  for valid in "${VALID_ENTITIES[@]}"; do
    if [ "$entity" = "$valid" ]; then
      return 0
    fi
  done
  echo "Error: Invalid entity type '$entity'" >&2
  echo "Valid entities: ${VALID_ENTITIES[*]}" >&2
  return 1
}

# Validate state
validate_state() {
  local state="$1"
  for valid in "${VALID_STATES[@]}"; do
    if [ "$state" = "$valid" ]; then
      return 0
    fi
  done
  echo "Error: Invalid state '$state'" >&2
  echo "Valid states: ${VALID_STATES[*]}" >&2
  return 1
}

# Get entity directory name (plural)
get_entity_dir() {
  local entity="$1"
  case "$entity" in
    ticket) echo "tickets" ;;
    erd) echo "erds" ;;
    prd) echo "prds" ;;
    phase) echo "phases" ;;
    *) echo "" ;;
  esac
}

# Get script name for entity
get_script_name() {
  local prefix="$1"
  local entity="$2"
  echo "${prefix}-${entity}.sh"
}

# Command: get-next
cmd_get_next() {
  if [ $# -lt 2 ]; then
    echo "Error: get-next requires entity and state arguments" >&2
    show_command_help get-next
    exit 1
  fi

  local entity="$1"
  local state="$2"

  if ! validate_entity "$entity"; then
    exit 1
  fi

  if ! validate_state "$state"; then
    exit 1
  fi

  local script_name=$(get_script_name "get-next" "$entity")
  local script_path="$SCRIPT_DIR/$script_name"

  if [ ! -f "$script_path" ]; then
    echo "Error: Script not found: $script_name" >&2
    exit 1
  fi

  if [ ! -x "$script_path" ]; then
    echo "Error: Script not executable: $script_name" >&2
    exit 1
  fi

  # Execute the script
  exec "$script_path" "$state"
}

# Command: move
cmd_move() {
  if [ $# -lt 3 ]; then
    echo "Error: move requires entity, item, and destination arguments" >&2
    show_command_help move
    exit 1
  fi

  local entity="$1"
  local item="$2"
  local destination="$3"

  if ! validate_entity "$entity"; then
    exit 1
  fi

  if ! validate_state "$destination"; then
    exit 1
  fi

  local script_name=$(get_script_name "move" "$entity")
  local script_path="$SCRIPT_DIR/$script_name"

  if [ ! -f "$script_path" ]; then
    echo "Error: Script not found: $script_name" >&2
    exit 1
  fi

  if [ ! -x "$script_path" ]; then
    echo "Error: Script not executable: $script_name" >&2
    exit 1
  fi

  # Execute the script
  exec "$script_path" "$item" "$destination"
}

# Command: list
cmd_list() {
  if [ $# -lt 2 ]; then
    echo "Error: list requires entity and state arguments" >&2
    show_command_help list
    exit 1
  fi

  local entity="$1"
  local state="$2"

  if ! validate_entity "$entity"; then
    exit 1
  fi

  if ! validate_state "$state"; then
    exit 1
  fi

  local entity_dir=$(get_entity_dir "$entity")
  local state_dir="$SCRIPT_DIR/$entity_dir/$state"

  if [ ! -d "$state_dir" ]; then
    echo "Error: Directory not found: $entity_dir/$state/" >&2
    exit 1
  fi

  # List all .md files, sorted lexicographically
  local files=$(find "$state_dir" -maxdepth 1 -type f -name "*.md" | sort)

  if [ -z "$files" ]; then
    echo "No ${entity}s found in $entity_dir/$state/" >&2
    exit 1
  fi

  # Display files
  echo "${BLUE}${entity}s in $state:${NC}"
  echo "$files" | while read -r file; do
    if [ -n "$file" ]; then
      echo "  $(basename "$file")"
    fi
  done
}

# Command: help
cmd_help() {
  if [ $# -eq 0 ]; then
    show_usage
  else
    show_command_help "$1"
  fi
}

# Main command dispatcher
main() {
  if [ $# -eq 0 ]; then
    show_usage
    exit 0
  fi

  local command="$1"
  shift

  case "$command" in
    get-next)
      cmd_get_next "$@"
      ;;
    move)
      cmd_move "$@"
      ;;
    list)
      cmd_list "$@"
      ;;
    help|--help|-h)
      cmd_help "$@"
      ;;
    *)
      echo "Error: Unknown command '$command'" >&2
      echo "" >&2
      show_usage
      exit 1
      ;;
  esac
}

# Run main function with all arguments
main "$@"

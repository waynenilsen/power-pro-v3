#!/usr/bin/env bash
#
# sdlc.sh - SDLC Management Tool
#
# Unified tool for managing tickets, sprints, and phases in nested structure:
#   phases/{state}/NNN-phase-name/
#     NNN-phase-name.md
#     sprints/{state}/NNN-sprint-name/
#       prd.md
#       erd.md
#       tickets/{state}/NNN-ticket-name.md


# Enable alias expansion and source environment
shopt -s expand_aliases
[ -f ~/.no-guard-bashrc.sh ] && source ~/.no-guard-bashrc.sh

# Get absolute path to this script's directory
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PHASES_DIR="$SCRIPT_DIR/phases"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Valid entity types and states
VALID_ENTITIES=("ticket" "sprint" "phase")
VALID_STATES=("todo" "in-progress" "done" "not-doing")

# ============================================================================
# Validation Functions
# ============================================================================

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

# ============================================================================
# Path Resolution Functions
# ============================================================================

# Find phase directory by number or name
find_phase() {
  local search="$1"
  local padded_search
  
  # If it's a number, pad it
  if [[ "$search" =~ ^[0-9]+$ ]]; then
    padded_search=$(printf "%03d" "$search")
  else
    padded_search="$search"
  fi
  
  # Search all state directories
  for state_dir in "$PHASES_DIR"/*; do
    if [ ! -d "$state_dir" ]; then
      continue
    fi
    
    # Try exact match
    local phase_dir="$state_dir/$padded_search"
    if [ -d "$phase_dir" ] && [ -f "$phase_dir"/*.md ]; then
      echo "$phase_dir"
      return 0
    fi
    
    # Try prefix match (NNN-*)
    local matches=$(find "$state_dir" -maxdepth 1 -type d -name "${padded_search}-*" | head -n 1)
    if [ -n "$matches" ]; then
      echo "$matches"
      return 0
    fi
  done
  
  return 1
}

# Find sprint directory by number or name
find_sprint() {
  local search="$1"
  local padded_search
  
  # If it's a number, pad it
  if [[ "$search" =~ ^[0-9]+$ ]]; then
    padded_search=$(printf "%03d" "$search")
  else
    padded_search="$search"
  fi
  
  # Search all phases and their sprints
  for phase_state_dir in "$PHASES_DIR"/*; do
    if [ ! -d "$phase_state_dir" ]; then
      continue
    fi
    
    for phase_dir in "$phase_state_dir"/*; do
      if [ ! -d "$phase_dir" ] || [ ! -f "$phase_dir"/*.md ]; then
        continue
      fi
      
      # Search all sprint state directories
      local sprints_dir="$phase_dir/sprints"
      if [ ! -d "$sprints_dir" ]; then
        continue
      fi
      
      for sprint_state_dir in "$sprints_dir"/*; do
        if [ ! -d "$sprint_state_dir" ]; then
          continue
        fi
        
        # Try exact match
        local sprint_dir="$sprint_state_dir/$padded_search"
        if [ -d "$sprint_dir" ] && [ -f "$sprint_dir/prd.md" ] && [ -f "$sprint_dir/erd.md" ]; then
          echo "$sprint_dir"
          return 0
        fi
        
        # Try prefix match
        local matches=$(find "$sprint_state_dir" -maxdepth 1 -type d -name "${padded_search}-*" | head -n 1)
        if [ -n "$matches" ] && [ -f "$matches/prd.md" ] && [ -f "$matches/erd.md" ]; then
          echo "$matches"
          return 0
        fi
      done
    done
  done
  
  return 1
}

# Find ticket file by number or name
find_ticket() {
  local search="$1"
  local padded_search
  
  # If it's a number, pad it
  if [[ "$search" =~ ^[0-9]+$ ]]; then
    padded_search=$(printf "%03d" "$search")
  else
    padded_search="$search"
  fi
  
  # Search all phases -> sprints -> tickets
  for phase_state_dir in "$PHASES_DIR"/*; do
    if [ ! -d "$phase_state_dir" ]; then
      continue
    fi
    
    for phase_dir in "$phase_state_dir"/*; do
      if [ ! -d "$phase_dir" ]; then
        continue
      fi
      
      local sprints_dir="$phase_dir/sprints"
      if [ ! -d "$sprints_dir" ]; then
        continue
      fi
      
      for sprint_state_dir in "$sprints_dir"/*; do
        if [ ! -d "$sprint_state_dir" ]; then
          continue
        fi
        
        for sprint_dir in "$sprint_state_dir"/*; do
          if [ ! -d "$sprint_dir" ]; then
            continue
          fi
          
          local tickets_dir="$sprint_dir/tickets"
          if [ ! -d "$tickets_dir" ]; then
            continue
          fi
          
          for ticket_state_dir in "$tickets_dir"/*; do
            if [ ! -d "$ticket_state_dir" ]; then
              continue
            fi
            
            # Try exact match
            local ticket_file="$ticket_state_dir/${padded_search}.md"
            if [ -f "$ticket_file" ]; then
              echo "$ticket_file"
              return 0
            fi
            
            # Try prefix match
            local matches=$(find "$ticket_state_dir" -maxdepth 1 -type f -name "${padded_search}-*.md" | head -n 1)
            if [ -n "$matches" ]; then
              echo "$matches"
              return 0
            fi
          done
        done
      done
    done
  done
  
  return 1
}

# ============================================================================
# Validation: Check if sprint can be moved to done
# ============================================================================

can_close_sprint() {
  local sprint_dir="$1"
  local tickets_dir="$sprint_dir/tickets"
  
  # Check if tickets directory exists
  if [ ! -d "$tickets_dir" ]; then
    return 0  # No tickets, can close
  fi
  
  # Check for todo or in-progress tickets
  for state in todo in-progress; do
    local state_dir="$tickets_dir/$state"
    if [ -d "$state_dir" ]; then
      local count=$(find "$state_dir" -maxdepth 1 -type f -name "*.md" | wc -l | tr -d ' ')
      if [ "$count" -gt 0 ]; then
        echo "Error: Cannot close sprint: $count ticket(s) in '$state' state" >&2
        return 1
      fi
    fi
  done
  
  return 0
}

# ============================================================================
# Get-Next Commands
# ============================================================================

get_next_phase() {
  local state="$1"
  local state_dir="$PHASES_DIR/$state"
  
  if [ ! -d "$state_dir" ]; then
    echo "Error: Directory not found: phases/$state/" >&2
    exit 1
  fi
  
  local phase_dir=$(find "$state_dir" -maxdepth 1 -type d ! -path "$state_dir" | sort | head -n 1)
  
  if [ -z "$phase_dir" ]; then
    echo "No phases found in phases/$state/" >&2
    exit 1
  fi
  
  echo "$phase_dir"
}

get_next_sprint() {
  local state="$1"
  local found=""
  
  # Search all phases for sprints in the specified state
  for phase_state_dir in "$PHASES_DIR"/*; do
    if [ ! -d "$phase_state_dir" ]; then
      continue
    fi
    
    for phase_dir in "$phase_state_dir"/*; do
      if [ ! -d "$phase_dir" ] || [ ! -f "$phase_dir"/*.md ]; then
        continue
      fi
      
      local sprint_state_dir="$phase_dir/sprints/$state"
      if [ ! -d "$sprint_state_dir" ]; then
        continue
      fi
      
      local sprint_dir=$(find "$sprint_state_dir" -maxdepth 1 -type d ! -path "$sprint_state_dir" | sort | head -n 1)
      if [ -n "$sprint_dir" ]; then
        echo "$sprint_dir"
        return 0
      fi
    done
  done
  
  echo "No sprints found in $state state" >&2
  exit 1
}

get_next_ticket() {
  local state="$1"
  local found=""
  
  # Search all phases -> sprints -> tickets in the specified state
  for phase_state_dir in "$PHASES_DIR"/*; do
    if [ ! -d "$phase_state_dir" ]; then
      continue
    fi
    
    for phase_dir in "$phase_state_dir"/*; do
      if [ ! -d "$phase_dir" ]; then
        continue
      fi
      
      local sprints_dir="$phase_dir/sprints"
      if [ ! -d "$sprints_dir" ]; then
        continue
      fi
      
      for sprint_state_dir in "$sprints_dir"/*; do
        if [ ! -d "$sprint_state_dir" ]; then
          continue
        fi
        
        for sprint_dir in "$sprint_state_dir"/*; do
          if [ ! -d "$sprint_dir" ]; then
            continue
          fi
          
          local ticket_state_dir="$sprint_dir/tickets/$state"
          if [ ! -d "$ticket_state_dir" ]; then
            continue
          fi
          
          local ticket=$(find "$ticket_state_dir" -maxdepth 1 -type f -name "*.md" | sort | head -n 1)
          if [ -n "$ticket" ]; then
            echo "$ticket"
            return 0
          fi
        done
      done
    done
  done
  
  echo "No tickets found in $state state" >&2
  exit 1
}

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
  
  case "$entity" in
    phase)
      get_next_phase "$state"
      ;;
    sprint)
      get_next_sprint "$state"
      ;;
    ticket)
      get_next_ticket "$state"
      ;;
    *)
      echo "Error: Unsupported entity type: $entity" >&2
      exit 1
      ;;
  esac
}

# ============================================================================
# Move Commands
# ============================================================================

move_phase() {
  local phase_input="$1"
  local destination="$2"
  local dest_dir="$PHASES_DIR/$destination"
  
  # Find the phase
  local phase_dir=$(find_phase "$phase_input")
  if [ -z "$phase_dir" ] || [ ! -d "$phase_dir" ]; then
    echo "Error: Phase not found: $phase_input" >&2
    exit 1
  fi
  
  local phase_name=$(basename "$phase_dir")
  local current_dir=$(dirname "$phase_dir")
  
  # Check if already in destination
  if [ "$(realpath "$current_dir")" = "$(realpath "$dest_dir")" ]; then
    echo "Phase is already in $destination/: $phase_name" >&2
    exit 0
  fi
  
  # Ensure destination exists
  mkdir -p "$dest_dir"
  
  # Move the phase directory
  mv "$phase_dir" "$dest_dir/$phase_name"
  echo "Moved phase $phase_name to $destination/"
}

move_sprint() {
  local sprint_input="$1"
  local destination="$2"
  
  # Find the sprint
  local sprint_dir=$(find_sprint "$sprint_input")
  if [ -z "$sprint_dir" ] || [ ! -d "$sprint_dir" ]; then
    echo "Error: Sprint not found: $sprint_input" >&2
    exit 1
  fi
  
  # Validate closing sprint
  if [ "$destination" = "done" ]; then
    if ! can_close_sprint "$sprint_dir"; then
      exit 1
    fi
  fi
  
  local sprint_name=$(basename "$sprint_dir")
  local current_sprint_state_dir=$(dirname "$sprint_dir")
  local phase_dir=$(dirname "$current_sprint_state_dir")
  phase_dir=$(dirname "$phase_dir")  # Go up from sprints/ to phase dir
  
  local dest_sprint_state_dir="$phase_dir/sprints/$destination"
  
  # Check if already in destination
  if [ "$(realpath "$current_sprint_state_dir")" = "$(realpath "$dest_sprint_state_dir")" ]; then
    echo "Sprint is already in $destination/: $sprint_name" >&2
    exit 0
  fi
  
  # Ensure destination exists
  mkdir -p "$dest_sprint_state_dir"
  
  # Move the sprint directory
  mv "$sprint_dir" "$dest_sprint_state_dir/$sprint_name"
  echo "Moved sprint $sprint_name to $destination/"
}

move_ticket() {
  local ticket_input="$1"
  local destination="$2"
  
  # Find the ticket
  local ticket_file=$(find_ticket "$ticket_input")
  if [ -z "$ticket_file" ] || [ ! -f "$ticket_file" ]; then
    echo "Error: Ticket not found: $ticket_input" >&2
    exit 1
  fi
  
  local ticket_name=$(basename "$ticket_file")
  local current_ticket_state_dir=$(dirname "$ticket_file")
  local sprint_dir=$(dirname "$current_ticket_state_dir")
  sprint_dir=$(dirname "$sprint_dir")  # Go up from tickets/ to sprint dir
  
  local dest_ticket_state_dir="$sprint_dir/tickets/$destination"
  
  # Check if already in destination
  if [ "$(realpath "$current_ticket_state_dir")" = "$(realpath "$dest_ticket_state_dir")" ]; then
    echo "Ticket is already in $destination/: $ticket_name" >&2
    exit 0
  fi
  
  # Ensure destination exists
  mkdir -p "$dest_ticket_state_dir"
  
  # Move the ticket file
  mv "$ticket_file" "$dest_ticket_state_dir/$ticket_name"
  echo "Moved ticket $ticket_name to $destination/"
}

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
  
  case "$entity" in
    phase)
      move_phase "$item" "$destination"
      ;;
    sprint)
      move_sprint "$item" "$destination"
      ;;
    ticket)
      move_ticket "$item" "$destination"
      ;;
    *)
      echo "Error: Unsupported entity type: $entity" >&2
      exit 1
      ;;
  esac
}

# ============================================================================
# List Commands
# ============================================================================

list_phases() {
  local state="$1"
  local state_dir="$PHASES_DIR/$state"
  
  if [ ! -d "$state_dir" ]; then
    echo "Error: Directory not found: phases/$state/" >&2
    exit 1
  fi
  
  local dirs=$(find "$state_dir" -maxdepth 1 -type d ! -path "$state_dir" | sort)
  
  if [ -z "$dirs" ]; then
    echo "No phases found in phases/$state/" >&2
    exit 1
  fi
  
  echo "${BLUE}Phases in $state:${NC}"
  echo "$dirs" | while read -r dir; do
    if [ -n "$dir" ]; then
      echo "  $(basename "$dir")"
    fi
  done
}

list_sprints() {
  local state="$1"
  local found=0
  
  echo "${BLUE}Sprints in $state:${NC}"
  for phase_state_dir in "$PHASES_DIR"/*; do
    if [ ! -d "$phase_state_dir" ]; then
      continue
    fi
    
    for phase_dir in "$phase_state_dir"/*; do
      if [ ! -d "$phase_dir" ] || [ ! -f "$phase_dir"/*.md ]; then
        continue
      fi
      
      local sprint_state_dir="$phase_dir/sprints/$state"
      if [ -d "$sprint_state_dir" ]; then
        local sprint_dirs=$(find "$sprint_state_dir" -maxdepth 1 -type d ! -path "$sprint_state_dir" | sort)
        echo "$sprint_dirs" | while read -r sprint_dir; do
          if [ -n "$sprint_dir" ]; then
            local phase_name=$(basename "$phase_dir")
            local sprint_name=$(basename "$sprint_dir")
            echo "  $phase_name/sprints/$state/$sprint_name"
            found=1
          fi
        done
      fi
    done
  done
  
  if [ $found -eq 0 ]; then
    echo "No sprints found in $state" >&2
    exit 1
  fi
}

list_tickets() {
  local state="$1"
  local found=0
  
  echo "${BLUE}Tickets in $state:${NC}"
  for phase_state_dir in "$PHASES_DIR"/*; do
    if [ ! -d "$phase_state_dir" ]; then
      continue
    fi
    
    for phase_dir in "$phase_state_dir"/*; do
      if [ ! -d "$phase_dir" ]; then
        continue
      fi
      
      local sprints_dir="$phase_dir/sprints"
      if [ ! -d "$sprints_dir" ]; then
        continue
      fi
      
      for sprint_state_dir in "$sprints_dir"/*; do
        if [ ! -d "$sprint_state_dir" ]; then
          continue
        fi
        
        for sprint_dir in "$sprint_state_dir"/*; do
          if [ ! -d "$sprint_dir" ]; then
            continue
          fi
          
          local ticket_state_dir="$sprint_dir/tickets/$state"
          if [ -d "$ticket_state_dir" ]; then
            local tickets=$(find "$ticket_state_dir" -maxdepth 1 -type f -name "*.md" | sort)
            echo "$tickets" | while read -r ticket; do
              if [ -n "$ticket" ]; then
                echo "  $(basename "$ticket")"
                found=1
              fi
            done
          fi
        done
      done
    done
  done
  
  if [ $found -eq 0 ]; then
    echo "No tickets found in $state" >&2
    exit 1
  fi
}

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
  
  case "$entity" in
    phase)
      list_phases "$state"
      ;;
    sprint)
      list_sprints "$state"
      ;;
    ticket)
      list_tickets "$state"
      ;;
    *)
      echo "Error: Unsupported entity type: $entity" >&2
      exit 1
      ;;
  esac
}

# ============================================================================
# Help Functions
# ============================================================================

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
  ticket                        Work tickets (nested under sprints)
  sprint                        Sprint documents (contain PRD and ERD, nested under phases)
  phase                         Phase documents

States:
  todo                          Items to be done
  in-progress                   Items currently being worked on
  done                          Completed items
  not-doing                     Cancelled or deferred items

Examples:
  $0 get-next ticket todo
  $0 get-next sprint in-progress
  $0 move ticket 001 in-progress
  $0 move sprint 001 done
  $0 list sprint todo
  $0 list phase in-progress
  $0 help get-next
  $0 help move

EOF
}

show_command_help() {
  local cmd="$1"
  case "$cmd" in
    get-next)
      cat << EOF
${BLUE}get-next${NC} - Get the next item from a state directory

Usage: $0 get-next <entity> <state>

Arguments:
  entity    Entity type (ticket, sprint, phase)
  state     State directory (todo, in-progress, done, not-doing)

Returns:
  Prints the path to the next item (lexicographically first)
  Exits with code 0 if item found, 1 if no items available

Examples:
  $0 get-next ticket todo
  $0 get-next sprint in-progress
  $0 get-next phase done

EOF
      ;;
    move)
      cat << EOF
${BLUE}move${NC} - Move an item between state directories

Usage: $0 move <entity> <item> <destination>

Arguments:
  entity       Entity type (ticket, sprint, phase)
  item         Item identifier (number, name, or path)
  destination  Destination state (todo, in-progress, done, not-doing)

Item can be specified as:
  - Number:    001 (searches for matching entity)
  - Name:      NNN-entity-name (searches for matching entity)
  - Full path: /path/to/entity

Note: Cannot move sprint to 'done' if it has todo/in-progress tickets.

Examples:
  $0 move ticket 001 in-progress
  $0 move sprint 001 done
  $0 move phase 001 in-progress

EOF
      ;;
    list)
      cat << EOF
${BLUE}list${NC} - List items in a state directory

Usage: $0 list <entity> <state>

Arguments:
  entity    Entity type (ticket, sprint, phase)
  state     State directory (todo, in-progress, done, not-doing)

Returns:
  Lists all items in the specified state directory, sorted lexicographically

Examples:
  $0 list ticket todo
  $0 list sprint in-progress
  $0 list phase done

EOF
      ;;
    *)
      echo "Unknown command: $cmd" >&2
      show_usage
      exit 1
      ;;
  esac
}

cmd_help() {
  if [ $# -eq 0 ]; then
    show_usage
  else
    show_command_help "$1"
  fi
}

# ============================================================================
# Main Command Dispatcher
# ============================================================================

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

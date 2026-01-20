#!/usr/bin/env bash
#
# inner-loop.sh - SDLC workflow automation
#
# Hands-off workflow automation that manages the full SDLC:
# - Implements tickets
# - Moves tickets through states
# - Closes sprints when all tickets are done
# - Breaks down sprint ERDs into tickets
# - Creates sprints from phases
# - Creates phases from README


# Enable alias expansion
shopt -s expand_aliases

# Source environment
[ -f ~/.no-guard-bashrc.sh ] && source ~/.no-guard-bashrc.sh

# Get absolute path to script directory
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# Colors for output
CYAN='\033[0;36m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# ============================================================================
# Helper Functions
# ============================================================================

log() {
  echo -e "${CYAN}[inner-loop]${NC} $*"
}

log_success() {
  echo -e "${GREEN}[inner-loop]${NC} $*"
}

log_warning() {
  echo -e "${YELLOW}[inner-loop]${NC} $*"
}

# Commit, add, and push changes with conventional commit format
commit_changes() {
  local type="$1"      # feat, chore, etc.
  local scope="$2"     # ticket, sprint, phase
  local description="$3"
  local body="${4:-}"
  
  # Check if there are any changes to commit
  if ! git diff --quiet || ! git diff --cached --quiet; then
    log "Staging, committing, and pushing changes..."
    git add -A
    
    local commit_msg="$type($scope): $description"
    if [ -n "$body" ]; then
      commit_msg="$commit_msg

$body"
    fi
    
    git commit -m "$commit_msg"
    git push
    log_success "Committed and pushed: $description"
  else
    log "No changes to commit"
  fi
}

# Find next available sprint number across all phases
find_next_sprint_number() {
  local max_num=0
  
  # Search all phases for sprint directories
  for phase_state_dir in "$SCRIPT_DIR/phases"/*; do
    if [ ! -d "$phase_state_dir" ]; then
      continue
    fi
    
    for phase_dir in "$phase_state_dir"/*; do
      if [ ! -d "$phase_dir" ] || [ ! -f "$phase_dir"/*.md ]; then
        continue
      fi
      
      local sprints_dir="$phase_dir/sprints"
      if [ ! -d "$sprints_dir" ]; then
        continue
      fi
      
      # Check all sprint state directories
      for sprint_state_dir in "$sprints_dir"/*; do
        if [ ! -d "$sprint_state_dir" ]; then
          continue
        fi
        
        # Check all sprint directories
        for sprint_dir in "$sprint_state_dir"/*; do
          if [ ! -d "$sprint_dir" ]; then
            continue
          fi
          
          local sprint_name=$(basename "$sprint_dir")
          local sprint_num=$(echo "$sprint_name" | sed -E 's/^([0-9]+)-.*/\1/' | sed 's/^0*//')
          if [ -n "$sprint_num" ] && [ "$sprint_num" -gt "$max_num" ] 2>/dev/null; then
            max_num=$sprint_num
          fi
        done
      done
    done
  done
  
  echo $((max_num + 1))
}

# Find next available phase number
find_next_phase_number() {
  local max_num=0
  
  for phase_state_dir in "$SCRIPT_DIR/phases"/*; do
    if [ ! -d "$phase_state_dir" ]; then
      continue
    fi
    
    for phase_dir in "$phase_state_dir"/*; do
      if [ ! -d "$phase_dir" ]; then
        continue
      fi
      
      local phase_name=$(basename "$phase_dir")
      local phase_num=$(echo "$phase_name" | sed -E 's/^([0-9]+)-.*/\1/' | sed 's/^0*//')
      if [ -n "$phase_num" ] && [ "$phase_num" -gt "$max_num" ] 2>/dev/null; then
        max_num=$phase_num
      fi
    done
  done
  
  echo $((max_num + 1))
}

# ============================================================================
# Main Workflow
# ============================================================================

main() {
  log "Starting SDLC workflow automation"
  
  # ========================================================================
  # Step 1: Implement in-progress tickets
  # ========================================================================
  log "Checking for in-progress tickets..."
  if NEXT_TICKET=$("./sdlc.sh" get-next ticket in-progress 2>&1); then
    log_success "Found in-progress ticket: $NEXT_TICKET"
    log "Implementing ticket..."
    "$SCRIPT_DIR/claude-wrapper.sh" "Implement the ticket at $NEXT_TICKET. Read the ticket file and implement all requirements specified in it."
    exit 0
  fi
  
  # ========================================================================
  # Step 2: Move todo tickets to in-progress
  # ========================================================================
  log "Checking for todo tickets..."
  if NEXT_TICKET=$("./sdlc.sh" get-next ticket todo 2>&1); then
    log_success "Found todo ticket: $NEXT_TICKET"
    log "Moving ticket to in-progress..."
    "./sdlc.sh" move ticket "$(basename "$NEXT_TICKET" .md)" in-progress
    commit_changes "chore" "ticket" "move $(basename "$NEXT_TICKET" .md) to in-progress"
    exit 0
  fi
  
  # ========================================================================
  # Step 3: Break down in-progress sprint ERD into tickets
  # ========================================================================
  log "Checking for in-progress sprints to break down..."
  if NEXT_SPRINT=$("./sdlc.sh" get-next sprint in-progress 2>&1); then
    log_success "Found in-progress sprint: $NEXT_SPRINT"
    
    # Check if sprint already has tickets
    local tickets_dir="$NEXT_SPRINT/tickets"
    local has_tickets=false
    
    if [ -d "$tickets_dir" ]; then
      for state_dir in "$tickets_dir"/*; do
        if [ -d "$state_dir" ] && [ -n "$(find "$state_dir" -maxdepth 1 -type f -name "*.md" 2>/dev/null)" ]; then
          has_tickets=true
          break
        fi
      done
    fi
    
    if [ "$has_tickets" = "false" ]; then
      log "Breaking down sprint ERD into tickets..."
      
      local erd_file="$NEXT_SPRINT/erd.md"
      local prd_file="$NEXT_SPRINT/prd.md"
      
      local prompt="Read the sprint ERD at $erd_file"
      if [ -f "$prd_file" ]; then
        prompt="$prompt and the associated PRD at $prd_file"
      fi
      prompt="$prompt. Also read the guide at $SCRIPT_DIR/prompts/erd-to-tickets.md. Break down the sprint ERD into tickets following the guidelines in erd-to-tickets.md. Use ./sdlc.sh commands to manage tickets. Create tickets in the sprint's tickets/todo/ directory (at $tickets_dir/todo/) with appropriate naming (NNN-description.md format). Ensure all tickets reference the ERD requirements they implement."
      
      "$SCRIPT_DIR/claude-wrapper.sh" "$prompt"
      exit 0
    else
      log "Sprint already has tickets, skipping breakdown"
    fi
  fi
  
  # ========================================================================
  # Step 4: Move todo sprints to in-progress
  # ========================================================================
  log "Checking for todo sprints..."
  if NEXT_SPRINT=$("./sdlc.sh" get-next sprint todo 2>&1); then
    log_success "Found todo sprint: $NEXT_SPRINT"
    log "Moving sprint to in-progress..."
    "./sdlc.sh" move sprint "$(basename "$NEXT_SPRINT")" in-progress
    commit_changes "chore" "sprint" "move $(basename "$NEXT_SPRINT") to in-progress"
    exit 0
  fi
  
  # ========================================================================
  # Step 5: Close ONE sprint that is ready (all tickets done)
  # Only close one sprint per iteration to ensure orderly progression
  # ========================================================================
  log "Checking for sprints ready to close..."
  if NEXT_SPRINT=$("./sdlc.sh" get-next sprint in-progress 2>&1); then
    # Try to move to done - sdlc.sh will validate that all tickets are done
    if "./sdlc.sh" move sprint "$(basename "$NEXT_SPRINT")" done 2>&1; then
      log_success "Closed sprint: $(basename "$NEXT_SPRINT")"
      commit_changes "chore" "sprint" "close sprint $(basename "$NEXT_SPRINT")"
      exit 0
    fi
    # Sprint has active tickets or can't be closed, continue to next steps
  fi
  
  # ========================================================================
  # Step 6: Create sprints from in-progress phase
  # ========================================================================
  log "Checking for in-progress phases to create sprints from..."
  if NEXT_PHASE=$("./sdlc.sh" get-next phase in-progress 2>&1); then
    log_success "Found in-progress phase: $NEXT_PHASE"
    
    # Check if phase already has sprints
    local sprints_dir="$NEXT_PHASE/sprints"
    local has_sprints=false
    
    if [ -d "$sprints_dir" ]; then
      for state_dir in "$sprints_dir"/*; do
        if [ -d "$state_dir" ] && [ -n "$(find "$state_dir" -maxdepth 1 -type d ! -path "$state_dir" 2>/dev/null)" ]; then
          has_sprints=true
          break
        fi
      done
    fi
    
    if [ "$has_sprints" = "false" ]; then
      log "Creating sprints from phase..."
      
      local phase_file="$NEXT_PHASE/$(basename "$NEXT_PHASE").md"
      if [ ! -f "$phase_file" ]; then
        # Try to find the phase document
        phase_file=$(find "$NEXT_PHASE" -maxdepth 1 -type f -name "*.md" | head -n 1)
      fi
      
      # Find next sprint numbers (need at least 3 sprints)
      local first_num=$(find_next_sprint_number)
      local second_num=$((first_num + 1))
      local third_num=$((first_num + 2))
      local first_padded=$(printf "%03d" "$first_num")
      local second_padded=$(printf "%03d" "$second_num")
      local third_padded=$(printf "%03d" "$third_num")
      
      local prompt="Read the phase document at $phase_file. Also read the ERD guidelines at $SCRIPT_DIR/prompts/erd.md and the roadmap guidelines at $SCRIPT_DIR/prompts/roadmap.md. One phase must correspond to at least 3 sprints. Create at least 3 sprints from this phase with the following numbers: $first_padded, $second_padded, $third_padded (create more if the phase warrants it). Each sprint is a directory containing both prd.md and erd.md files. Create sprint directories in $sprints_dir/todo/ with the naming format NNN-description/. After creating all sprints, use ./sdlc.sh move sprint commands to move ONLY the FIRST sprint ($first_padded) to in-progress. Leave all other sprints in todo. Remember: every 5th sprint must be a technical debt paydown sprint."
      
      "$SCRIPT_DIR/claude-wrapper.sh" "$prompt"
      exit 0
    else
      log "Phase already has sprints, skipping sprint creation"
    fi
  fi
  
  # ========================================================================
  # Step 8: Move todo phase to in-progress
  # ========================================================================
  log "Checking for todo phases..."
  if NEXT_PHASE=$("./sdlc.sh" get-next phase todo 2>&1); then
    log_success "Found todo phase: $NEXT_PHASE"
    log "Moving phase to in-progress..."
    "./sdlc.sh" move phase "$(basename "$NEXT_PHASE")" in-progress
    commit_changes "chore" "phase" "move $(basename "$NEXT_PHASE") to in-progress"
    exit 0
  fi
  
  # ========================================================================
  # Step 9: Create new phase from README
  # ========================================================================
  log "No phases available, creating new phase from README..."
  
  local next_phase_num=$(find_next_phase_number)
  local next_phase_padded=$(printf "%03d" "$next_phase_num")
  local phases_todo_dir="$SCRIPT_DIR/phases/todo"
  
  local prompt="Read the README.md file at $SCRIPT_DIR/README.md which serves as the roadmap document. Also read the roadmap guidelines at $SCRIPT_DIR/prompts/roadmap.md. Create a new phase directory and document based on the product vision and roadmap in README.md. The phase should be numbered $next_phase_padded and follow the naming format NNN-description. Create the phase directory in $phases_todo_dir/ with the naming format NNN-description/. Inside the phase directory, create the phase document file NNN-description.md. Follow the phase document format specified in the roadmap guidelines, including Vision & Strategic Objectives, Themes & Initiatives, Timeline, Success Metrics, and Review & Update Process. The phase should align with the product vision in README.md. Do not create any sprints, PRDs, ERDs, or tickets, just the phase directory and document."
  
  "$SCRIPT_DIR/claude-wrapper.sh" "$prompt"
  exit 0
}

# Run main function
main

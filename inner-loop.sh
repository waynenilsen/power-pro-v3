#!/usr/bin/env bash
#
# new-inner-loop.sh - Deterministically compose and run prompts based on ralph logic
#
# This script deterministically builds prompts using bash logic instead of
# relying on relative path resolution in Claude Max.

# Enable alias expansion in non-interactive shell
shopt -s expand_aliases

# Source no-guard-bashrc.sh to give node bun bla bla all tools to claude as well as to get the claude alias
[ -f ~/.no-guard-bashrc.sh ] && source ~/.no-guard-bashrc.sh

# Debug: Script start
echo "[DEBUG] new-inner-loop.sh: Starting"

# Get absolute path to this script's directory
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROMPTS_DIR="$SCRIPT_DIR"
echo "[DEBUG] new-inner-loop.sh: SCRIPT_DIR=$SCRIPT_DIR"
echo "[DEBUG] new-inner-loop.sh: PROMPTS_DIR=$PROMPTS_DIR"

# Main execution
main() {
  echo "[DEBUG] new-inner-loop.sh: Entering main()"
  
  # Get the next in-progress ticket using sdlc.sh
  echo "[DEBUG] Getting next in-progress ticket..."
  NEXT_TICKET=$(./sdlc.sh get-next ticket in-progress 2>&1)
  
  if [ $? -eq 0 ]; then
    echo "[DEBUG] Found ticket: $NEXT_TICKET"
    echo "Next ticket: $NEXT_TICKET"
    
    # Tell Claude to implement the ticket
    echo "[DEBUG] Calling Claude to implement ticket..."
    ./claude-wrapper.sh "Implement the ticket at $NEXT_TICKET. Read the ticket file and implement all requirements specified in it."
    exit 0
  fi
  
  echo "[DEBUG] No in-progress tickets found: $NEXT_TICKET"
  echo "No in-progress tickets available"
  
  # Get the next todo ticket using sdlc.sh
  echo "[DEBUG] Getting next todo ticket..."
  NEXT_TICKET=$(./sdlc.sh get-next ticket todo 2>&1)
  
  if [ $? -eq 0 ]; then
    echo "[DEBUG] Found ticket: $NEXT_TICKET"
    echo "Next ticket: $NEXT_TICKET"
    
    # Move ticket to in-progress
    echo "[DEBUG] Moving ticket to in-progress..."
    ./sdlc.sh move ticket "$NEXT_TICKET" in-progress
    exit 0
  fi
  
  echo "[DEBUG] No todo tickets found: $NEXT_TICKET"
  echo "No todo tickets available"
  
  # Move all in-progress ERDs to done
  echo "[DEBUG] Moving all in-progress ERDs to done..."
  while true; do
    NEXT_ERD=$(./sdlc.sh get-next erd in-progress 2>&1)
    if [ $? -ne 0 ]; then
      break
    fi
    echo "[DEBUG] Moving ERD to done: $NEXT_ERD"
    ./sdlc.sh move erd "$NEXT_ERD" done
  done
  
  # Move all in-progress PRDs to done
  echo "[DEBUG] Moving all in-progress PRDs to done..."
  while true; do
    NEXT_PRD=$(./sdlc.sh get-next prd in-progress 2>&1)
    if [ $? -ne 0 ]; then
      break
    fi
    echo "[DEBUG] Moving PRD to done: $NEXT_PRD"
    ./sdlc.sh move prd "$NEXT_PRD" done
  done
  
  # Move next todo ERD to in-progress
  echo "[DEBUG] Getting next todo ERD..."
  NEXT_ERD=$(./sdlc.sh get-next erd todo 2>&1)
  
  if [ $? -eq 0 ]; then
    echo "[DEBUG] Found ERD: $NEXT_ERD"
    echo "Next ERD: $NEXT_ERD"
    
    # Move ERD to in-progress
    echo "[DEBUG] Moving ERD to in-progress..."
    ./sdlc.sh move erd "$NEXT_ERD" in-progress
    exit 0
  fi
  
  # Move next todo PRD to in-progress
  echo "[DEBUG] Getting next todo PRD..."
  NEXT_PRD=$(./sdlc.sh get-next prd todo 2>&1)
  
  if [ $? -eq 0 ]; then
    echo "[DEBUG] Found PRD: $NEXT_PRD"
    echo "Next PRD: $NEXT_PRD"
    
    # Move PRD to in-progress
    echo "[DEBUG] Moving PRD to in-progress..."
    ./sdlc.sh move prd "$NEXT_PRD" in-progress
    exit 0
  fi
  
  # Check for in-progress ERD and break it down into tickets
  echo "[DEBUG] Checking for in-progress ERD to break down into tickets..."
  NEXT_ERD=$(./sdlc.sh get-next erd in-progress 2>&1)
  
  if [ $? -eq 0 ]; then
    echo "[DEBUG] Found in-progress ERD: $NEXT_ERD"
    echo "Breaking down ERD into tickets: $NEXT_ERD"
    
    # Extract ERD number from filename (e.g., "001-erd-name.md" -> "001")
    ERD_FILENAME=$(basename "$NEXT_ERD")
    ERD_NUMBER=$(echo "$ERD_FILENAME" | sed -E 's/^([0-9]+)-.*/\1/')
    
    # Find associated PRD (same number)
    PRD_PATH=""
    for dir in prds/todo prds/in-progress prds/done prds/not-doing; do
      if [ -d "$SCRIPT_DIR/$dir" ]; then
        PRD_CANDIDATE=$(find "$SCRIPT_DIR/$dir" -maxdepth 1 -type f -name "${ERD_NUMBER}-*.md" 2>/dev/null | head -n 1)
        if [ -n "$PRD_CANDIDATE" ]; then
          PRD_PATH="$PRD_CANDIDATE"
          break
        fi
      fi
    done
    
    # Build prompt for Claude
    PROMPT="Read the ERD at $NEXT_ERD"
    if [ -n "$PRD_PATH" ]; then
      PROMPT="$PROMPT and the associated PRD at $PRD_PATH"
    fi
    PROMPT="$PROMPT. Also read the guide at $SCRIPT_DIR/prompts/erd-to-tickets.md. Break down the ERD into tickets following the guidelines in erd-to-tickets.md. Use ./sdlc.sh commands to manage tickets. Create tickets in the tickets/todo/ directory with appropriate naming (NNN-description.md format). Ensure all tickets reference the ERD requirements they implement."
    
    echo "[DEBUG] Calling Claude to break down ERD into tickets..."
    ./claude-wrapper.sh "$PROMPT"
    exit 0
  fi
  
  echo "[DEBUG] No in-progress ERDs found to break down"
  echo "[DEBUG] No todo ERDs or PRDs available"
  
  # Check if there are any todo phases that need to be broken down into PRDs
  echo "[DEBUG] Checking for todo phases..."
  NEXT_PHASE=$(./sdlc.sh get-next phase todo 2>&1)
  
  if [ $? -eq 0 ]; then
    echo "[DEBUG] Found todo phase: $NEXT_PHASE"
    echo "Creating PRD and ERD from phase: $NEXT_PHASE"
    
    # Move phase to in-progress
    echo "[DEBUG] Moving phase to in-progress..."
    ./sdlc.sh move phase "$NEXT_PHASE" in-progress
    
    # Extract phase number from filename (e.g., "001-phase-name.md" -> "001")
    PHASE_FILENAME=$(basename "$NEXT_PHASE")
    PHASE_NUMBER=$(echo "$PHASE_FILENAME" | sed -E 's/^([0-9]+)-.*/\1/')
    
    # Find the next available PRD/ERD number
    # Check existing PRDs and ERDs to determine next number
    MAX_NUM=0
    for dir in prds/todo prds/in-progress prds/done prds/not-doing erds/todo erds/in-progress erds/done erds/not-doing; do
      if [ -d "$SCRIPT_DIR/$dir" ]; then
        for file in "$SCRIPT_DIR/$dir"/*.md; do
          if [ -f "$file" ]; then
            FILENAME=$(basename "$file")
            FILE_NUM=$(echo "$FILENAME" | sed -E 's/^([0-9]+)-.*/\1/' | sed 's/^0*//')
            if [ -n "$FILE_NUM" ] && [ "$FILE_NUM" -gt "$MAX_NUM" ] 2>/dev/null; then
              MAX_NUM=$FILE_NUM
            fi
          fi
        done
      fi
    done
    
    # Calculate numbers for at least 3 PRD/ERD pairs
    FIRST_NUM=$((MAX_NUM + 1))
    SECOND_NUM=$((MAX_NUM + 2))
    THIRD_NUM=$((MAX_NUM + 3))
    FIRST_NUM_PADDED=$(printf "%03d" "$FIRST_NUM")
    SECOND_NUM_PADDED=$(printf "%03d" "$SECOND_NUM")
    THIRD_NUM_PADDED=$(printf "%03d" "$THIRD_NUM")
    
    # Build prompt for Claude
    PROMPT="Read the phase document at $NEXT_PHASE. Also read the ERD guidelines at $SCRIPT_DIR/prompts/erd.md and the roadmap guidelines at $SCRIPT_DIR/prompts/roadmap.md. One phase must correspond to at least 3 PRDs and ERDs. Create at least 3 PRD (Product Requirements Document) and ERD (Engineering Requirements Document) pairs from this phase with the following numbers: $FIRST_NUM_PADDED, $SECOND_NUM_PADDED, $THIRD_NUM_PADDED (create more if the phase warrants it). Each PRD and its associated ERD must share the same number (e.g., PRD $FIRST_NUM_PADDED pairs with ERD $FIRST_NUM_PADDED). Create all files in their respective todo directories (prds/todo/ and erds/todo/) with the naming format NNN-description.md. After creating all files, use ./sdlc.sh move prd and ./sdlc.sh move erd commands to move ONLY the FIRST pair ($FIRST_NUM_PADDED) to in-progress. Leave all other PRD/ERD pairs in todo. Remember: 1 PRD maps to precisely 1 ERD, and every 5th PRD/ERD must be a technical debt paydown PRD/ERD."
    
    echo "[DEBUG] Calling Claude to create PRD and ERD from phase..."
    ./claude-wrapper.sh "$PROMPT"
    exit 0
  fi
  
  # If we get here, there are no phases - create one from README
  echo "[DEBUG] No phases available, creating new phase from README..."
  
  # Find the next available phase number
  MAX_NUM=0
  for dir in phases/todo phases/in-progress phases/done phases/not-doing; do
    if [ -d "$SCRIPT_DIR/$dir" ]; then
      for file in "$SCRIPT_DIR/$dir"/*.md; do
        if [ -f "$file" ]; then
          FILENAME=$(basename "$file")
          FILE_NUM=$(echo "$FILENAME" | sed -E 's/^([0-9]+)-.*/\1/' | sed 's/^0*//')
          if [ -n "$FILE_NUM" ] && [ "$FILE_NUM" -gt "$MAX_NUM" ] 2>/dev/null; then
            MAX_NUM=$FILE_NUM
          fi
        fi
      done
    fi
  done
  NEXT_NUM=$((MAX_NUM + 1))
  NEXT_NUM_PADDED=$(printf "%03d" "$NEXT_NUM")
  
  # Build prompt for Claude
  PROMPT="Read the README.md file at $SCRIPT_DIR/README.md which serves as the roadmap document. Also read the roadmap guidelines at $SCRIPT_DIR/prompts/roadmap.md. Create a new phase document based on the product vision and roadmap in README.md. The phase should be numbered $NEXT_NUM_PADDED and follow the naming format NNN-description.md. Create the phase document in the phases/todo/ directory. Follow the phase document format specified in the roadmap guidelines, including Vision & Strategic Objectives, Themes & Initiatives, Timeline, Success Metrics, and Review & Update Process. The phase should align with the product vision in README.md do not create any PRDs or ERDs or tickets, just the phase document."
  
  echo "[DEBUG] Calling Claude to create phase from README..."
  ./claude-wrapper.sh "$PROMPT"
  exit 0
}

main

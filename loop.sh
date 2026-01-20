#!/usr/bin/env bash
#
# new-outer-loop.sh - Run new-inner-loop.sh in a continuous loop
#
# Usage:
#   ./new-outer-loop.sh          # Run indefinitely
#   ./new-outer-loop.sh 5         # Run 5 iterations
#
# Each iteration:
#   1. Runs new-inner-loop.sh
#   2. Repeats
#


# Enable alias expansion in non-interactive shell
shopt -s expand_aliases

# Source no-guard-bashrc.sh to give node bun bla bla all tools to claude as well as to get the claude alias
[ -f ~/.no-guard-bashrc.sh ] && source ~/.no-guard-bashrc.sh

MAX_ITERATIONS="${1:-0}"  # 0 = unlimited
ITERATION=0

# Get absolute path to this script's directory
pushd "$(dirname "$0")" >/dev/null
SCRIPT_DIR="$(pwd)"
popd >/dev/null

# Colors
CYAN='\033[0;36m'
DIM='\033[2m'
RESET='\033[0m'

main() {
  echo -e "${CYAN}outer loop starting${RESET}"
  [ "$MAX_ITERATIONS" -gt 0 ] && echo -e "${DIM}max iterations: ${MAX_ITERATIONS}${RESET}"
  echo ""

  # Source bashrc to ensure agent has access to environment
  [ -f ~/.bashrc ] && source ~/.bashrc || [ -f "$HOME/.bashrc" ] && source "$HOME/.bashrc"

  while true; do
    ITERATION=$((ITERATION + 1))

    echo -e "${CYAN}━━━ iteration ${ITERATION} ━━━${RESET}"
    echo ""

    "$SCRIPT_DIR/inner-loop.sh"

    # Check iteration limit
    if [ "$MAX_ITERATIONS" -gt 0 ] && [ "$ITERATION" -ge "$MAX_ITERATIONS" ]; then
      echo ""
      echo -e "${CYAN}completed ${ITERATION} iterations${RESET}"
      break
    fi

    echo ""
  done
}

main

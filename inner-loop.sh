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
  SLEEP 10
}

main

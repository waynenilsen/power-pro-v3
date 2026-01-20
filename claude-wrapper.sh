#!/usr/bin/env bash
#
# claude-wrapper.sh - Wrapper for claude command with standard flags
#
# Usage:
#   ./claude-wrapper.sh "<any prompt string>"


# Enable alias expansion in non-interactive shell
shopt -s expand_aliases

# Source no-guard-bashrc.sh to give node bun bla bla all tools to claude as well as to get the claude alias
[ -f ~/.no-guard-bashrc.sh ] && source ~/.no-guard-bashrc.sh

if [ -z "${1:-}" ]; then
  echo "Usage: claude-wrapper.sh <prompt>"
  echo "  e.g., claude-wrapper.sh \"run the promptgram @promptgrams/ralph.md\""
  exit 1
fi

PROMPT="$1"

claude -p "$PROMPT if you have created any files, you must commit and push them using conventional commits, update gitignore if needed, this is a unit of work for you." \
  --dangerously-skip-permissions \
  --output-format stream-json \
  --verbose | cclean

#!/usr/bin/env bash
#
# test-get-next-scripts.sh - Test all get-next-* scripts and clean up
#
# This script creates test files, runs all the get-next-* scripts,
# verifies they work correctly, and then cleans up all test files.

# Don't use set -e - we want to count failures manually

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Track test files for cleanup
TEST_FILES=()

# Function to create a test file
create_test_file() {
  local dir="$1"
  local filename="$2"
  local full_path="$dir/$filename"
  mkdir -p "$dir"
  echo "# Test $filename" > "$full_path"
  TEST_FILES+=("$full_path")
}

# Function to test a script
test_script() {
  local script_cmd="$1"
  local expected_file="$2"
  local description="$3"
  
  echo -n "Testing $description... "
  
  # Extract script path (first word)
  local script="${script_cmd%% *}"
  
  if [ ! -f "$script" ]; then
    echo -e "${RED}FAIL${NC} - Script not found: $script"
    return 1
  fi
  
  if [ ! -x "$script" ]; then
    echo -e "${RED}FAIL${NC} - Script not executable: $script"
    return 1
  fi
  
  local result
  if result=$(eval "$script_cmd" 2>&1); then
    if [ "$result" = "$expected_file" ]; then
      echo -e "${GREEN}PASS${NC}"
      return 0
    else
      echo -e "${RED}FAIL${NC} - Expected: $expected_file, Got: $result"
      return 1
    fi
  else
    echo -e "${RED}FAIL${NC} - Script failed: $result"
    return 1
  fi
}

# Function to test error case (should fail)
test_error_case() {
  local script_cmd="$1"
  local description="$2"
  
  echo -n "Testing $description (should fail)... "
  
  if eval "$script_cmd" 2>&1 > /dev/null; then
    echo -e "${RED}FAIL${NC} - Should have failed but succeeded"
    return 1
  else
    echo -e "${GREEN}PASS${NC}"
    return 0
  fi
}

# Function to test move script
test_move() {
  local script_cmd="$1"
  local expected_location="$2"
  local description="$3"
  
  echo -n "Testing $description... "
  
  # Extract script path (first word)
  local script="${script_cmd%% *}"
  
  if [ ! -f "$script" ]; then
    echo -e "${RED}FAIL${NC} - Script not found: $script"
    return 1
  fi
  
  if [ ! -x "$script" ]; then
    echo -e "${RED}FAIL${NC} - Script not executable: $script"
    return 1
  fi
  
  local result
  if result=$(eval "$script_cmd" 2>&1); then
    # Check if file exists in expected location
    if [ -f "$expected_location" ]; then
      echo -e "${GREEN}PASS${NC}"
      return 0
    else
      echo -e "${RED}FAIL${NC} - File not found at expected location: $expected_location"
      return 1
    fi
  else
    echo -e "${RED}FAIL${NC} - Script failed: $result"
    return 1
  fi
}

echo -e "${YELLOW}Creating test files...${NC}"

# Create test files for tickets
create_test_file "tickets/todo" "001-test-ticket-a.md"
create_test_file "tickets/todo" "002-test-ticket-b.md"
create_test_file "tickets/in-progress" "010-test-ticket-c.md"
create_test_file "tickets/in-progress" "020-test-ticket-d.md"

# Create test files for ERDs
create_test_file "erds/todo" "001-test-erd-a.md"
create_test_file "erds/todo" "002-test-erd-b.md"
create_test_file "erds/in-progress" "010-test-erd-c.md"
create_test_file "erds/in-progress" "020-test-erd-d.md"

# Create test files for PRDs
create_test_file "prds/todo" "001-test-prd-a.md"
create_test_file "prds/todo" "002-test-prd-b.md"
create_test_file "prds/in-progress" "010-test-prd-c.md"
create_test_file "prds/in-progress" "020-test-prd-d.md"

# Create test files for phases
create_test_file "phases/todo" "001-test-phase-a.md"
create_test_file "phases/todo" "002-test-phase-b.md"
create_test_file "phases/in-progress" "010-test-phase-c.md"
create_test_file "phases/in-progress" "020-test-phase-d.md"

echo -e "${YELLOW}Running tests...${NC}"
echo ""

PASSED=0
FAILED=0

# Test ticket scripts
echo "=== Ticket Scripts ==="
if test_script "./get-next-ticket.sh todo" "$SCRIPT_DIR/tickets/todo/001-test-ticket-a.md" "get-next-ticket.sh todo"; then ((PASSED++)); else ((FAILED++)); fi
if test_script "./get-next-ticket.sh in-progress" "$SCRIPT_DIR/tickets/in-progress/010-test-ticket-c.md" "get-next-ticket.sh in-progress"; then ((PASSED++)); else ((FAILED++)); fi
if test_script "./get-next-todo-ticket.sh" "$SCRIPT_DIR/tickets/todo/001-test-ticket-a.md" "get-next-todo-ticket.sh"; then ((PASSED++)); else ((FAILED++)); fi
if test_script "./get-next-in-progress-ticket.sh" "$SCRIPT_DIR/tickets/in-progress/010-test-ticket-c.md" "get-next-in-progress-ticket.sh"; then ((PASSED++)); else ((FAILED++)); fi
if test_error_case "./get-next-ticket.sh" "get-next-ticket.sh (no args)"; then ((PASSED++)); else ((FAILED++)); fi
if test_error_case "./get-next-ticket.sh invalid-state" "get-next-ticket.sh invalid-state"; then ((PASSED++)); else ((FAILED++)); fi
echo ""

# Test ERD scripts
echo "=== ERD Scripts ==="
if test_script "./get-next-erd.sh todo" "$SCRIPT_DIR/erds/todo/001-test-erd-a.md" "get-next-erd.sh todo"; then ((PASSED++)); else ((FAILED++)); fi
if test_script "./get-next-erd.sh in-progress" "$SCRIPT_DIR/erds/in-progress/010-test-erd-c.md" "get-next-erd.sh in-progress"; then ((PASSED++)); else ((FAILED++)); fi
if test_script "./get-next-todo-erd.sh" "$SCRIPT_DIR/erds/todo/001-test-erd-a.md" "get-next-todo-erd.sh"; then ((PASSED++)); else ((FAILED++)); fi
if test_script "./get-next-in-progress-erd.sh" "$SCRIPT_DIR/erds/in-progress/010-test-erd-c.md" "get-next-in-progress-erd.sh"; then ((PASSED++)); else ((FAILED++)); fi
if test_error_case "./get-next-erd.sh" "get-next-erd.sh (no args)"; then ((PASSED++)); else ((FAILED++)); fi
if test_error_case "./get-next-erd.sh invalid-state" "get-next-erd.sh invalid-state"; then ((PASSED++)); else ((FAILED++)); fi
echo ""

# Test PRD scripts
echo "=== PRD Scripts ==="
if test_script "./get-next-prd.sh todo" "$SCRIPT_DIR/prds/todo/001-test-prd-a.md" "get-next-prd.sh todo"; then ((PASSED++)); else ((FAILED++)); fi
if test_script "./get-next-prd.sh in-progress" "$SCRIPT_DIR/prds/in-progress/010-test-prd-c.md" "get-next-prd.sh in-progress"; then ((PASSED++)); else ((FAILED++)); fi
if test_script "./get-next-todo-prd.sh" "$SCRIPT_DIR/prds/todo/001-test-prd-a.md" "get-next-todo-prd.sh"; then ((PASSED++)); else ((FAILED++)); fi
if test_script "./get-next-in-progress-prd.sh" "$SCRIPT_DIR/prds/in-progress/010-test-prd-c.md" "get-next-in-progress-prd.sh"; then ((PASSED++)); else ((FAILED++)); fi
if test_error_case "./get-next-prd.sh" "get-next-prd.sh (no args)"; then ((PASSED++)); else ((FAILED++)); fi
if test_error_case "./get-next-prd.sh invalid-state" "get-next-prd.sh invalid-state"; then ((PASSED++)); else ((FAILED++)); fi
echo ""

# Test phase scripts
echo "=== Phase Scripts ==="
if test_script "./get-next-phase.sh todo" "$SCRIPT_DIR/phases/todo/001-test-phase-a.md" "get-next-phase.sh todo"; then ((PASSED++)); else ((FAILED++)); fi
if test_script "./get-next-phase.sh in-progress" "$SCRIPT_DIR/phases/in-progress/010-test-phase-c.md" "get-next-phase.sh in-progress"; then ((PASSED++)); else ((FAILED++)); fi
if test_script "./get-next-todo-phase.sh" "$SCRIPT_DIR/phases/todo/001-test-phase-a.md" "get-next-todo-phase.sh"; then ((PASSED++)); else ((FAILED++)); fi
if test_script "./get-next-in-progress-phase.sh" "$SCRIPT_DIR/phases/in-progress/010-test-phase-c.md" "get-next-in-progress-phase.sh"; then ((PASSED++)); else ((FAILED++)); fi
if test_error_case "./get-next-phase.sh" "get-next-phase.sh (no args)"; then ((PASSED++)); else ((FAILED++)); fi
if test_error_case "./get-next-phase.sh invalid-state" "get-next-phase.sh invalid-state"; then ((PASSED++)); else ((FAILED++)); fi
echo ""

# Test empty directory case
echo "=== Empty Directory Tests ==="
# Create empty directories
mkdir -p "tickets/done"
mkdir -p "erds/done"
mkdir -p "prds/done"
mkdir -p "phases/done"
if test_error_case "./get-next-ticket.sh done" "get-next-ticket.sh done (empty)"; then ((PASSED++)); else ((FAILED++)); fi
if test_error_case "./get-next-erd.sh done" "get-next-erd.sh done (empty)"; then ((PASSED++)); else ((FAILED++)); fi
if test_error_case "./get-next-prd.sh done" "get-next-prd.sh done (empty)"; then ((PASSED++)); else ((FAILED++)); fi
if test_error_case "./get-next-phase.sh done" "get-next-phase.sh done (empty)"; then ((PASSED++)); else ((FAILED++)); fi
echo ""

# Test move scripts
echo "=== Move Script Tests ==="
# Test move-ticket.sh
if test_move "./move-ticket.sh 001-test-ticket-a.md in-progress" "$SCRIPT_DIR/tickets/in-progress/001-test-ticket-a.md" "move-ticket.sh (filename)"; then ((PASSED++)); else ((FAILED++)); fi
if test_move "./move-ticket.sh 001 done" "$SCRIPT_DIR/tickets/done/001-test-ticket-a.md" "move-ticket.sh (number)"; then ((PASSED++)); else ((FAILED++)); fi
if test_error_case "./move-ticket.sh" "move-ticket.sh (no args)"; then ((PASSED++)); else ((FAILED++)); fi
if test_error_case "./move-ticket.sh nonexistent.md todo" "move-ticket.sh (file not found)"; then ((PASSED++)); else ((FAILED++)); fi
if test_error_case "./move-ticket.sh 001-test-ticket-a.md invalid-dest" "move-ticket.sh (invalid destination)"; then ((PASSED++)); else ((FAILED++)); fi
# Move it back for other tests
./move-ticket.sh 001-test-ticket-a.md todo > /dev/null 2>&1 || true
echo ""

# Test move-erd.sh
if test_move "./move-erd.sh 001-test-erd-a.md in-progress" "$SCRIPT_DIR/erds/in-progress/001-test-erd-a.md" "move-erd.sh (filename)"; then ((PASSED++)); else ((FAILED++)); fi
if test_move "./move-erd.sh 001 done" "$SCRIPT_DIR/erds/done/001-test-erd-a.md" "move-erd.sh (number)"; then ((PASSED++)); else ((FAILED++)); fi
if test_error_case "./move-erd.sh" "move-erd.sh (no args)"; then ((PASSED++)); else ((FAILED++)); fi
if test_error_case "./move-erd.sh nonexistent.md todo" "move-erd.sh (file not found)"; then ((PASSED++)); else ((FAILED++)); fi
if test_error_case "./move-erd.sh 001-test-erd-a.md invalid-dest" "move-erd.sh (invalid destination)"; then ((PASSED++)); else ((FAILED++)); fi
# Move it back for other tests
./move-erd.sh 001-test-erd-a.md todo > /dev/null 2>&1 || true
echo ""

# Test move-prd.sh
if test_move "./move-prd.sh 001-test-prd-a.md in-progress" "$SCRIPT_DIR/prds/in-progress/001-test-prd-a.md" "move-prd.sh (filename)"; then ((PASSED++)); else ((FAILED++)); fi
if test_move "./move-prd.sh 001 done" "$SCRIPT_DIR/prds/done/001-test-prd-a.md" "move-prd.sh (number)"; then ((PASSED++)); else ((FAILED++)); fi
if test_error_case "./move-prd.sh" "move-prd.sh (no args)"; then ((PASSED++)); else ((FAILED++)); fi
if test_error_case "./move-prd.sh nonexistent.md todo" "move-prd.sh (file not found)"; then ((PASSED++)); else ((FAILED++)); fi
if test_error_case "./move-prd.sh 001-test-prd-a.md invalid-dest" "move-prd.sh (invalid destination)"; then ((PASSED++)); else ((FAILED++)); fi
# Move it back for other tests
./move-prd.sh 001-test-prd-a.md todo > /dev/null 2>&1 || true
echo ""

# Test move-phase.sh
if test_move "./move-phase.sh 001-test-phase-a.md in-progress" "$SCRIPT_DIR/phases/in-progress/001-test-phase-a.md" "move-phase.sh (filename)"; then ((PASSED++)); else ((FAILED++)); fi
if test_move "./move-phase.sh 001 done" "$SCRIPT_DIR/phases/done/001-test-phase-a.md" "move-phase.sh (number)"; then ((PASSED++)); else ((FAILED++)); fi
if test_error_case "./move-phase.sh" "move-phase.sh (no args)"; then ((PASSED++)); else ((FAILED++)); fi
if test_error_case "./move-phase.sh nonexistent.md todo" "move-phase.sh (file not found)"; then ((PASSED++)); else ((FAILED++)); fi
if test_error_case "./move-phase.sh 001-test-phase-a.md invalid-dest" "move-phase.sh (invalid destination)"; then ((PASSED++)); else ((FAILED++)); fi
# Move it back for other tests
./move-phase.sh 001-test-phase-a.md todo > /dev/null 2>&1 || true
echo ""

# Summary
echo "=== Test Summary ==="
echo -e "Passed: ${GREEN}$PASSED${NC}"
echo -e "Failed: ${RED}$FAILED${NC}"
echo ""

# Cleanup
echo -e "${YELLOW}Cleaning up test files...${NC}"
for file in "${TEST_FILES[@]}"; do
  if [ -f "$file" ]; then
    rm -f "$file"
  fi
done

# Remove empty test directories (but keep .gitkeep files)
find tickets/todo tickets/in-progress tickets/done erds/todo erds/in-progress erds/done prds/todo prds/in-progress prds/done phases/todo phases/in-progress phases/done -type d -empty -delete 2>/dev/null || true

echo -e "${GREEN}Cleanup complete!${NC}"

# Exit with appropriate code
if [ $FAILED -eq 0 ]; then
  exit 0
else
  exit 1
fi

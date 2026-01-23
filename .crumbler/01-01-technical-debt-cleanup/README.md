# Technical Debt Cleanup

## Objective
Address code quality, test coverage, and documentation alignment issues from Phase 1 rapid development.

## Scope

Based on tech-debt.md guidelines, identify and remediate:

### 1. Code Debt
- Files exceeding 1000 lines (particularly test files)
- Duplicated logic across files
- Large, complex functions
- Missing error handling
- Tight coupling

### 2. Test Debt
- Ensure all domain entities have comprehensive tests
- Fix any flaky tests
- Improve test coverage for edge cases

### 3. Documentation Debt
- Ensure domain models are documented
- API handlers have clear godoc comments
- Any complex logic has explanatory comments

## Known Large Files (>1000 lines)
Based on analysis, the following test files exceed 1000 lines:
- `internal/api/prescription_handler_test.go` (1752 lines)
- `internal/domain/prescription/prescription_test.go` (1513 lines)
- `internal/domain/setscheme/ramp_test.go` (1416 lines)
- `internal/domain/liftmax/liftmax_test.go` (1384 lines)
- `internal/domain/progression/trigger_test.go` (1259 lines)
- `internal/domain/progression/progression_test.go` (1254 lines)
- `internal/service/progression_service_test.go` (1198 lines)
- `internal/api/liftmax_handler_test.go` (1168 lines)
- `internal/domain/loadstrategy/percentof_test.go` (1117 lines)
- `internal/domain/userprogramstate/userprogramstate_test.go` (1004 lines)

Note: Large test files are acceptable if well-organized. Focus on production code quality.

## Acceptance Criteria
- [ ] All production code files under 500 lines (or justified)
- [ ] No duplicated logic across files
- [ ] All tests passing
- [ ] Code follows Go best practices
- [ ] Domain model documentation complete

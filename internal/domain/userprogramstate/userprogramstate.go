// Package userprogramstate provides domain logic for the UserProgramState entity.
// This package contains pure business logic with no database dependencies,
// making it testable in isolation.
package userprogramstate

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Validation errors
var (
	ErrUserIDRequired            = errors.New("user_id is required")
	ErrProgramIDRequired         = errors.New("program_id is required")
	ErrCurrentWeekInvalid        = errors.New("current_week must be at least 1")
	ErrCurrentCycleIterationInvalid = errors.New("current_cycle_iteration must be at least 1")
	ErrCurrentDayIndexInvalid    = errors.New("current_day_index must be at least 0 if provided")
)

// UserProgramState represents a user's enrollment in a program with their current position.
type UserProgramState struct {
	ID                    string
	UserID                string
	ProgramID             string
	CurrentWeek           int
	CurrentCycleIteration int
	CurrentDayIndex       *int
	EnrolledAt            time.Time
	UpdatedAt             time.Time
}

// EnrollmentWithProgram represents a user's enrollment with program details for responses.
type EnrollmentWithProgram struct {
	State            *UserProgramState
	ProgramName      string
	ProgramSlug      string
	ProgramDescription *string
	CycleLengthWeeks int
}

// ValidationResult contains the result of validating a user program state.
type ValidationResult struct {
	Valid  bool
	Errors []error
}

// NewValidationResult creates a valid result.
func NewValidationResult() *ValidationResult {
	return &ValidationResult{Valid: true, Errors: []error{}}
}

// AddError adds an error to the validation result and marks it invalid.
func (v *ValidationResult) AddError(err error) {
	v.Valid = false
	v.Errors = append(v.Errors, err)
}

// Error returns a combined error message if there are validation errors.
func (v *ValidationResult) Error() error {
	if v.Valid {
		return nil
	}
	var msgs []string
	for _, err := range v.Errors {
		msgs = append(msgs, err.Error())
	}
	return fmt.Errorf("validation failed: %s", strings.Join(msgs, "; "))
}

// ValidateUserID validates the user ID field.
func ValidateUserID(userID string) error {
	if strings.TrimSpace(userID) == "" {
		return ErrUserIDRequired
	}
	return nil
}

// ValidateProgramID validates the program ID field.
func ValidateProgramID(programID string) error {
	if strings.TrimSpace(programID) == "" {
		return ErrProgramIDRequired
	}
	return nil
}

// ValidateCurrentWeek validates the current week field.
func ValidateCurrentWeek(week int) error {
	if week < 1 {
		return ErrCurrentWeekInvalid
	}
	return nil
}

// ValidateCurrentCycleIteration validates the current cycle iteration field.
func ValidateCurrentCycleIteration(iteration int) error {
	if iteration < 1 {
		return ErrCurrentCycleIterationInvalid
	}
	return nil
}

// ValidateCurrentDayIndex validates the current day index field.
func ValidateCurrentDayIndex(dayIndex *int) error {
	if dayIndex != nil && *dayIndex < 0 {
		return ErrCurrentDayIndexInvalid
	}
	return nil
}

// EnrollUserInput contains the input data for enrolling a user in a program.
type EnrollUserInput struct {
	UserID    string
	ProgramID string
}

// EnrollUser validates input and creates a new UserProgramState for enrollment.
// Initial state: currentWeek = 1, currentCycleIteration = 1
func EnrollUser(input EnrollUserInput, id string) (*UserProgramState, *ValidationResult) {
	result := NewValidationResult()

	// Validate user_id
	if err := ValidateUserID(input.UserID); err != nil {
		result.AddError(err)
	}

	// Validate program_id
	if err := ValidateProgramID(input.ProgramID); err != nil {
		result.AddError(err)
	}

	if !result.Valid {
		return nil, result
	}

	now := time.Now()
	return &UserProgramState{
		ID:                    id,
		UserID:                strings.TrimSpace(input.UserID),
		ProgramID:             strings.TrimSpace(input.ProgramID),
		CurrentWeek:           1, // Initial state
		CurrentCycleIteration: 1, // Initial state
		CurrentDayIndex:       nil,
		EnrolledAt:            now,
		UpdatedAt:             now,
	}, result
}

// UpdatePositionInput contains the input data for updating user position in the program.
type UpdatePositionInput struct {
	CurrentWeek           *int
	CurrentCycleIteration *int
	CurrentDayIndex       **int // Double pointer: nil = no change, *nil = clear, *value = set
}

// UpdatePosition updates the user's position within the program.
func UpdatePosition(state *UserProgramState, input UpdatePositionInput) *ValidationResult {
	result := NewValidationResult()

	if input.CurrentWeek != nil {
		if err := ValidateCurrentWeek(*input.CurrentWeek); err != nil {
			result.AddError(err)
		} else {
			state.CurrentWeek = *input.CurrentWeek
		}
	}

	if input.CurrentCycleIteration != nil {
		if err := ValidateCurrentCycleIteration(*input.CurrentCycleIteration); err != nil {
			result.AddError(err)
		} else {
			state.CurrentCycleIteration = *input.CurrentCycleIteration
		}
	}

	if input.CurrentDayIndex != nil {
		newDayIndex := *input.CurrentDayIndex
		if err := ValidateCurrentDayIndex(newDayIndex); err != nil {
			result.AddError(err)
		} else {
			state.CurrentDayIndex = newDayIndex
		}
	}

	if result.Valid {
		state.UpdatedAt = time.Now()
	}

	return result
}

// Validate performs full validation on an existing user program state.
func (s *UserProgramState) Validate() *ValidationResult {
	result := NewValidationResult()

	if err := ValidateUserID(s.UserID); err != nil {
		result.AddError(err)
	}

	if err := ValidateProgramID(s.ProgramID); err != nil {
		result.AddError(err)
	}

	if err := ValidateCurrentWeek(s.CurrentWeek); err != nil {
		result.AddError(err)
	}

	if err := ValidateCurrentCycleIteration(s.CurrentCycleIteration); err != nil {
		result.AddError(err)
	}

	if err := ValidateCurrentDayIndex(s.CurrentDayIndex); err != nil {
		result.AddError(err)
	}

	return result
}

// AdvancementContext contains the context needed for state advancement.
type AdvancementContext struct {
	DaysInCurrentWeek int // Number of training days in the current week
	CycleLengthWeeks  int // Total number of weeks in the cycle
}

// AdvancementResult contains the result of a state advancement operation.
type AdvancementResult struct {
	NewState       *UserProgramState
	CycleCompleted bool
}

// AdvanceState advances the user's position in the program.
// Logic:
// 1. Increment day index
// 2. If day index >= days in week, reset day index to 0, increment week
// 3. If week > cycle length, reset week to 1, increment cycle iteration (cycle completed)
func AdvanceState(state *UserProgramState, ctx AdvancementContext) (*AdvancementResult, *ValidationResult) {
	result := NewValidationResult()

	// Validate context
	if ctx.DaysInCurrentWeek < 1 {
		result.AddError(errors.New("days_in_current_week must be at least 1"))
		return nil, result
	}
	if ctx.CycleLengthWeeks < 1 {
		result.AddError(errors.New("cycle_length_weeks must be at least 1"))
		return nil, result
	}

	// Copy state to avoid mutating original
	newState := &UserProgramState{
		ID:                    state.ID,
		UserID:                state.UserID,
		ProgramID:             state.ProgramID,
		CurrentWeek:           state.CurrentWeek,
		CurrentCycleIteration: state.CurrentCycleIteration,
		CurrentDayIndex:       copyIntPtr(state.CurrentDayIndex),
		EnrolledAt:            state.EnrolledAt,
		UpdatedAt:             time.Now(),
	}

	cycleCompleted := false

	// Initialize day index if nil (first advancement)
	currentDayIndex := 0
	if newState.CurrentDayIndex != nil {
		currentDayIndex = *newState.CurrentDayIndex
	}

	// Increment day index
	currentDayIndex++

	// Check if we've exceeded days in the week
	if currentDayIndex >= ctx.DaysInCurrentWeek {
		// Reset day index to 0
		currentDayIndex = 0

		// Increment week
		newState.CurrentWeek++

		// Check if we've exceeded the cycle length
		if newState.CurrentWeek > ctx.CycleLengthWeeks {
			// Reset to week 1
			newState.CurrentWeek = 1

			// Increment cycle iteration
			newState.CurrentCycleIteration++

			// Mark cycle as completed
			cycleCompleted = true
		}
	}

	// Set the new day index
	newState.CurrentDayIndex = &currentDayIndex

	return &AdvancementResult{
		NewState:       newState,
		CycleCompleted: cycleCompleted,
	}, result
}

// copyIntPtr creates a copy of an int pointer.
func copyIntPtr(p *int) *int {
	if p == nil {
		return nil
	}
	v := *p
	return &v
}

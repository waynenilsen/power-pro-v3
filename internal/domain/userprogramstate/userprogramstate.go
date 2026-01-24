// Package userprogramstate provides domain logic for the UserProgramState entity.
// This package contains pure business logic with no database dependencies,
// making it testable in isolation.
package userprogramstate

import (
	"errors"
	"strings"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/validation"
)

// ScheduleType represents how a program's schedule is determined.
type ScheduleType string

const (
	// ScheduleTypeRotation means the program follows a rotating schedule (default).
	ScheduleTypeRotation ScheduleType = "rotation"
	// ScheduleTypeDaysOut means the program schedule is determined by days until meet date.
	ScheduleTypeDaysOut ScheduleType = "days_out"
)

// Validation errors
var (
	ErrUserIDRequired               = errors.New("user_id is required")
	ErrProgramIDRequired            = errors.New("program_id is required")
	ErrCurrentWeekInvalid           = errors.New("current_week must be at least 1")
	ErrCurrentCycleIterationInvalid = errors.New("current_cycle_iteration must be at least 1")
	ErrCurrentDayIndexInvalid       = errors.New("current_day_index must be at least 0 if provided")
	ErrMeetDateInPast               = errors.New("meet_date must be in the future")
	ErrInvalidScheduleType          = errors.New("schedule_type must be 'rotation' or 'days_out'")
)

// UserProgramState represents a user's enrollment in a program with their current position.
type UserProgramState struct {
	ID                    string
	UserID                string
	ProgramID             string
	CurrentWeek           int
	CurrentCycleIteration int
	CurrentDayIndex       *int
	RotationPosition      int          // 0-based position in rotation for programs like Conjugate/Westside
	CyclesSinceStart      int          // Number of complete cycles since enrollment
	MeetDate              *time.Time   // Optional meet date for peaking programs
	ScheduleType          ScheduleType // "rotation" (default) or "days_out"
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

// ValidationResult is an alias for the shared validation.Result type.
type ValidationResult = validation.Result

// NewValidationResult creates a valid result.
func NewValidationResult() *ValidationResult {
	return validation.NewResult()
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

// ValidateMeetDate validates the meet date field.
// If provided, the meet date must be in the future.
func ValidateMeetDate(meetDate *time.Time) error {
	if meetDate != nil && !meetDate.After(time.Now()) {
		return ErrMeetDateInPast
	}
	return nil
}

// ValidateScheduleType validates the schedule type field.
func ValidateScheduleType(scheduleType ScheduleType) error {
	if scheduleType != ScheduleTypeRotation && scheduleType != ScheduleTypeDaysOut {
		return ErrInvalidScheduleType
	}
	return nil
}

// EnrollUserInput contains the input data for enrolling a user in a program.
type EnrollUserInput struct {
	UserID       string
	ProgramID    string
	MeetDate     *time.Time   // Optional meet date for peaking programs
	ScheduleType ScheduleType // Optional; defaults to "rotation" if empty
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

	// Validate meet_date if provided
	if err := ValidateMeetDate(input.MeetDate); err != nil {
		result.AddError(err)
	}

	// Default schedule type to rotation if not specified
	scheduleType := input.ScheduleType
	if scheduleType == "" {
		scheduleType = ScheduleTypeRotation
	}

	// Validate schedule type
	if err := ValidateScheduleType(scheduleType); err != nil {
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
		RotationPosition:      0,            // Start at first position in rotation
		CyclesSinceStart:      0,            // No cycles completed yet
		MeetDate:              input.MeetDate,
		ScheduleType:          scheduleType,
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

	if err := ValidateMeetDate(s.MeetDate); err != nil {
		result.AddError(err)
	}

	// Validate schedule type only if set (empty is treated as rotation)
	if s.ScheduleType != "" {
		if err := ValidateScheduleType(s.ScheduleType); err != nil {
			result.AddError(err)
		}
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
		RotationPosition:      state.RotationPosition,
		CyclesSinceStart:      state.CyclesSinceStart,
		MeetDate:              copyTimePtr(state.MeetDate),
		ScheduleType:          state.ScheduleType,
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

			// Track total cycles completed since enrollment
			newState.CyclesSinceStart++

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

// copyTimePtr creates a copy of a time.Time pointer.
func copyTimePtr(p *time.Time) *time.Time {
	if p == nil {
		return nil
	}
	v := *p
	return &v
}

// AdvanceRotation advances the rotation position by 1 and wraps around when it
// reaches rotationLength. This is used by programs like Conjugate/Westside that
// cycle through different lift focuses.
//
// Example: With rotationLength=3 (deadlift, squat, bench):
//   - Position 0 → 1
//   - Position 1 → 2
//   - Position 2 → 0 (wraps around)
//
// If rotationLength is <= 0, this function does nothing.
func (s *UserProgramState) AdvanceRotation(rotationLength int) {
	if rotationLength <= 0 {
		return
	}
	s.RotationPosition = (s.RotationPosition + 1) % rotationLength
	s.UpdatedAt = time.Now()
}

// IncrementCyclesSinceStart increments the count of completed cycles.
// This is called when a full program cycle completes.
func (s *UserProgramState) IncrementCyclesSinceStart() {
	s.CyclesSinceStart++
	s.UpdatedAt = time.Now()
}

// UpdateMeetDateInput contains the input data for updating a user's meet date.
type UpdateMeetDateInput struct {
	MeetDate *time.Time // nil means clear the meet date
}

// UpdateMeetDate updates the user's meet date.
// If the meet date is nil, it clears the meet date and sets schedule type to rotation.
// If the meet date is provided, it validates the date is in the future.
func (s *UserProgramState) UpdateMeetDate(input UpdateMeetDateInput) *ValidationResult {
	result := NewValidationResult()

	// Validate meet date if provided
	if err := ValidateMeetDate(input.MeetDate); err != nil {
		result.AddError(err)
		return result
	}

	// Update meet date
	s.MeetDate = input.MeetDate

	// If clearing meet date, also reset schedule type to rotation
	if input.MeetDate == nil {
		s.ScheduleType = ScheduleTypeRotation
	} else {
		// When setting a meet date, default to days_out schedule
		s.ScheduleType = ScheduleTypeDaysOut
	}

	s.UpdatedAt = time.Now()
	return result
}

// DaysOut calculates the number of days until the meet date.
// Returns 0 if meet date is not set.
func (s *UserProgramState) DaysOut() int {
	if s.MeetDate == nil {
		return 0
	}
	now := time.Now()
	duration := s.MeetDate.Sub(now)
	days := int(duration.Hours() / 24)
	if days < 0 {
		return 0
	}
	return days
}

// WeeksToMeet calculates the number of weeks until the meet date.
// Returns 0 if meet date is not set.
func (s *UserProgramState) WeeksToMeet() int {
	days := s.DaysOut()
	return days / 7
}

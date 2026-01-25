// Package workoutsession provides domain logic for the WorkoutSession entity.
// This package contains pure business logic with no database dependencies,
// making it testable in isolation.
package workoutsession

import (
	"errors"
	"strings"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/validation"
)

// Status represents the status of a workout session.
type Status string

const (
	StatusInProgress Status = "IN_PROGRESS"
	StatusCompleted  Status = "COMPLETED"
	StatusAbandoned  Status = "ABANDONED"
)

// Validation errors
var (
	ErrIDRequired                 = errors.New("id is required")
	ErrUserProgramStateIDRequired = errors.New("user_program_state_id is required")
	ErrWeekNumberInvalid          = errors.New("week_number must be at least 1")
	ErrDayIndexInvalid            = errors.New("day_index must be at least 0")
	ErrInvalidStatus              = errors.New("status must be 'IN_PROGRESS', 'COMPLETED', or 'ABANDONED'")
	ErrAlreadyCompleted           = errors.New("session is already completed")
	ErrAlreadyAbandoned           = errors.New("session is already abandoned")
	ErrNotInProgress              = errors.New("session is not in progress")
)

// WorkoutSession represents a single workout session within a user's program.
type WorkoutSession struct {
	ID                 string
	UserProgramStateID string
	WeekNumber         int
	DayIndex           int
	Status             Status
	StartedAt          time.Time
	FinishedAt         *time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// ValidationResult is an alias for the shared validation.Result type.
type ValidationResult = validation.Result

// NewValidationResult creates a valid result.
func NewValidationResult() *ValidationResult {
	return validation.NewResult()
}

// ValidateID validates the session ID field.
func ValidateID(id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrIDRequired
	}
	return nil
}

// ValidateUserProgramStateID validates the user program state ID field.
func ValidateUserProgramStateID(userProgramStateID string) error {
	if strings.TrimSpace(userProgramStateID) == "" {
		return ErrUserProgramStateIDRequired
	}
	return nil
}

// ValidateWeekNumber validates the week number field.
func ValidateWeekNumber(weekNumber int) error {
	if weekNumber < 1 {
		return ErrWeekNumberInvalid
	}
	return nil
}

// ValidateDayIndex validates the day index field.
func ValidateDayIndex(dayIndex int) error {
	if dayIndex < 0 {
		return ErrDayIndexInvalid
	}
	return nil
}

// ValidateStatus validates the status field.
func ValidateStatus(status Status) error {
	switch status {
	case StatusInProgress, StatusCompleted, StatusAbandoned:
		return nil
	default:
		return ErrInvalidStatus
	}
}

// NewWorkoutSessionInput contains the input data for creating a new workout session.
type NewWorkoutSessionInput struct {
	UserProgramStateID string
	WeekNumber         int
	DayIndex           int
}

// NewWorkoutSession validates input and creates a new WorkoutSession.
// The session starts with IN_PROGRESS status.
func NewWorkoutSession(input NewWorkoutSessionInput, id string) (*WorkoutSession, *ValidationResult) {
	result := NewValidationResult()

	// Validate ID
	if err := ValidateID(id); err != nil {
		result.AddError(err)
	}

	// Validate user_program_state_id
	if err := ValidateUserProgramStateID(input.UserProgramStateID); err != nil {
		result.AddError(err)
	}

	// Validate week_number
	if err := ValidateWeekNumber(input.WeekNumber); err != nil {
		result.AddError(err)
	}

	// Validate day_index
	if err := ValidateDayIndex(input.DayIndex); err != nil {
		result.AddError(err)
	}

	if !result.Valid {
		return nil, result
	}

	now := time.Now()
	return &WorkoutSession{
		ID:                 id,
		UserProgramStateID: strings.TrimSpace(input.UserProgramStateID),
		WeekNumber:         input.WeekNumber,
		DayIndex:           input.DayIndex,
		Status:             StatusInProgress,
		StartedAt:          now,
		FinishedAt:         nil,
		CreatedAt:          now,
		UpdatedAt:          now,
	}, result
}

// Complete marks the workout session as completed.
// Returns an error if the session is not in progress.
func (ws *WorkoutSession) Complete() error {
	if ws.Status == StatusCompleted {
		return ErrAlreadyCompleted
	}
	if ws.Status == StatusAbandoned {
		return ErrNotInProgress
	}
	if ws.Status != StatusInProgress {
		return ErrNotInProgress
	}

	now := time.Now()
	ws.Status = StatusCompleted
	ws.FinishedAt = &now
	ws.UpdatedAt = now
	return nil
}

// Abandon marks the workout session as abandoned.
// Returns an error if the session is not in progress.
func (ws *WorkoutSession) Abandon() error {
	if ws.Status == StatusAbandoned {
		return ErrAlreadyAbandoned
	}
	if ws.Status == StatusCompleted {
		return ErrNotInProgress
	}
	if ws.Status != StatusInProgress {
		return ErrNotInProgress
	}

	now := time.Now()
	ws.Status = StatusAbandoned
	ws.FinishedAt = &now
	ws.UpdatedAt = now
	return nil
}

// IsActive returns true if the session is currently in progress.
func (ws *WorkoutSession) IsActive() bool {
	return ws.Status == StatusInProgress
}

// IsFinished returns true if the session has been completed or abandoned.
func (ws *WorkoutSession) IsFinished() bool {
	return ws.Status == StatusCompleted || ws.Status == StatusAbandoned
}

// Validate performs full validation on an existing workout session.
func (ws *WorkoutSession) Validate() *ValidationResult {
	result := NewValidationResult()

	if err := ValidateID(ws.ID); err != nil {
		result.AddError(err)
	}

	if err := ValidateUserProgramStateID(ws.UserProgramStateID); err != nil {
		result.AddError(err)
	}

	if err := ValidateWeekNumber(ws.WeekNumber); err != nil {
		result.AddError(err)
	}

	if err := ValidateDayIndex(ws.DayIndex); err != nil {
		result.AddError(err)
	}

	if err := ValidateStatus(ws.Status); err != nil {
		result.AddError(err)
	}

	return result
}

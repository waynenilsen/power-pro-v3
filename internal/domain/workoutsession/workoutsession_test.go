package workoutsession

import (
	"errors"
	"testing"
)

// ==================== ID Validation Tests ====================

func TestValidateID_Valid(t *testing.T) {
	tests := []struct {
		name string
		id   string
	}{
		{"uuid", "550e8400-e29b-41d4-a716-446655440000"},
		{"simple id", "session-1"},
		{"short id", "s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateID(tt.id)
			if err != nil {
				t.Errorf("ValidateID(%q) = %v, want nil", tt.id, err)
			}
		})
	}
}

func TestValidateID_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		expectedErr error
	}{
		{"empty string", "", ErrIDRequired},
		{"only spaces", "   ", ErrIDRequired},
		{"only tabs", "\t\t", ErrIDRequired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateID(tt.id)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateID(%q) = %v, want %v", tt.id, err, tt.expectedErr)
			}
		})
	}
}

// ==================== UserProgramStateID Validation Tests ====================

func TestValidateUserProgramStateID_Valid(t *testing.T) {
	tests := []struct {
		name               string
		userProgramStateID string
	}{
		{"uuid", "550e8400-e29b-41d4-a716-446655440000"},
		{"simple id", "state-1"},
		{"short id", "s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserProgramStateID(tt.userProgramStateID)
			if err != nil {
				t.Errorf("ValidateUserProgramStateID(%q) = %v, want nil", tt.userProgramStateID, err)
			}
		})
	}
}

func TestValidateUserProgramStateID_Invalid(t *testing.T) {
	tests := []struct {
		name               string
		userProgramStateID string
		expectedErr        error
	}{
		{"empty string", "", ErrUserProgramStateIDRequired},
		{"only spaces", "   ", ErrUserProgramStateIDRequired},
		{"only tabs", "\t\t", ErrUserProgramStateIDRequired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserProgramStateID(tt.userProgramStateID)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateUserProgramStateID(%q) = %v, want %v", tt.userProgramStateID, err, tt.expectedErr)
			}
		})
	}
}

// ==================== WeekNumber Validation Tests ====================

func TestValidateWeekNumber_Valid(t *testing.T) {
	tests := []struct {
		name       string
		weekNumber int
	}{
		{"week 1", 1},
		{"week 10", 10},
		{"large week", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateWeekNumber(tt.weekNumber)
			if err != nil {
				t.Errorf("ValidateWeekNumber(%d) = %v, want nil", tt.weekNumber, err)
			}
		})
	}
}

func TestValidateWeekNumber_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		weekNumber  int
		expectedErr error
	}{
		{"zero", 0, ErrWeekNumberInvalid},
		{"negative", -1, ErrWeekNumberInvalid},
		{"very negative", -100, ErrWeekNumberInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateWeekNumber(tt.weekNumber)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateWeekNumber(%d) = %v, want %v", tt.weekNumber, err, tt.expectedErr)
			}
		})
	}
}

// ==================== DayIndex Validation Tests ====================

func TestValidateDayIndex_Valid(t *testing.T) {
	tests := []struct {
		name     string
		dayIndex int
	}{
		{"zero", 0},
		{"positive", 5},
		{"large", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDayIndex(tt.dayIndex)
			if err != nil {
				t.Errorf("ValidateDayIndex(%d) = %v, want nil", tt.dayIndex, err)
			}
		})
	}
}

func TestValidateDayIndex_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		dayIndex    int
		expectedErr error
	}{
		{"negative", -1, ErrDayIndexInvalid},
		{"very negative", -100, ErrDayIndexInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDayIndex(tt.dayIndex)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateDayIndex(%d) = %v, want %v", tt.dayIndex, err, tt.expectedErr)
			}
		})
	}
}

// ==================== Status Validation Tests ====================

func TestValidateStatus_Valid(t *testing.T) {
	tests := []struct {
		name   string
		status Status
	}{
		{"in progress", StatusInProgress},
		{"completed", StatusCompleted},
		{"abandoned", StatusAbandoned},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStatus(tt.status)
			if err != nil {
				t.Errorf("ValidateStatus(%q) = %v, want nil", tt.status, err)
			}
		})
	}
}

func TestValidateStatus_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		status      Status
		expectedErr error
	}{
		{"empty string", Status(""), ErrInvalidStatus},
		{"invalid value", Status("INVALID"), ErrInvalidStatus},
		{"lowercase in_progress", Status("in_progress"), ErrInvalidStatus},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStatus(tt.status)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateStatus(%q) = %v, want %v", tt.status, err, tt.expectedErr)
			}
		})
	}
}

// ==================== NewWorkoutSession Tests ====================

func TestNewWorkoutSession_ValidInput(t *testing.T) {
	input := NewWorkoutSessionInput{
		UserProgramStateID: "state-1",
		WeekNumber:         1,
		DayIndex:           0,
	}

	session, result := NewWorkoutSession(input, "session-id")

	if !result.Valid {
		t.Errorf("NewWorkoutSession returned invalid result: %v", result.Errors)
	}
	if session == nil {
		t.Fatal("NewWorkoutSession returned nil session")
	}
	if session.ID != "session-id" {
		t.Errorf("session.ID = %q, want %q", session.ID, "session-id")
	}
	if session.UserProgramStateID != "state-1" {
		t.Errorf("session.UserProgramStateID = %q, want %q", session.UserProgramStateID, "state-1")
	}
	if session.WeekNumber != 1 {
		t.Errorf("session.WeekNumber = %d, want %d", session.WeekNumber, 1)
	}
	if session.DayIndex != 0 {
		t.Errorf("session.DayIndex = %d, want %d", session.DayIndex, 0)
	}
	if session.Status != StatusInProgress {
		t.Errorf("session.Status = %q, want %q", session.Status, StatusInProgress)
	}
	if session.FinishedAt != nil {
		t.Errorf("session.FinishedAt = %v, want nil", session.FinishedAt)
	}
}

func TestNewWorkoutSession_TrimsWhitespace(t *testing.T) {
	input := NewWorkoutSessionInput{
		UserProgramStateID: "  state-1  ",
		WeekNumber:         1,
		DayIndex:           0,
	}

	session, result := NewWorkoutSession(input, "session-id")

	if !result.Valid {
		t.Errorf("NewWorkoutSession returned invalid result: %v", result.Errors)
	}
	if session.UserProgramStateID != "state-1" {
		t.Errorf("session.UserProgramStateID = %q, want %q (trimmed)", session.UserProgramStateID, "state-1")
	}
}

func TestNewWorkoutSession_EmptyID(t *testing.T) {
	input := NewWorkoutSessionInput{
		UserProgramStateID: "state-1",
		WeekNumber:         1,
		DayIndex:           0,
	}

	session, result := NewWorkoutSession(input, "")

	if result.Valid {
		t.Error("NewWorkoutSession with empty id returned valid result")
	}
	if session != nil {
		t.Error("NewWorkoutSession with invalid input returned non-nil session")
	}
}

func TestNewWorkoutSession_EmptyUserProgramStateID(t *testing.T) {
	input := NewWorkoutSessionInput{
		UserProgramStateID: "",
		WeekNumber:         1,
		DayIndex:           0,
	}

	session, result := NewWorkoutSession(input, "session-id")

	if result.Valid {
		t.Error("NewWorkoutSession with empty user_program_state_id returned valid result")
	}
	if session != nil {
		t.Error("NewWorkoutSession with invalid input returned non-nil session")
	}
}

func TestNewWorkoutSession_InvalidWeekNumber(t *testing.T) {
	input := NewWorkoutSessionInput{
		UserProgramStateID: "state-1",
		WeekNumber:         0,
		DayIndex:           0,
	}

	session, result := NewWorkoutSession(input, "session-id")

	if result.Valid {
		t.Error("NewWorkoutSession with invalid week_number returned valid result")
	}
	if session != nil {
		t.Error("NewWorkoutSession with invalid input returned non-nil session")
	}
}

func TestNewWorkoutSession_InvalidDayIndex(t *testing.T) {
	input := NewWorkoutSessionInput{
		UserProgramStateID: "state-1",
		WeekNumber:         1,
		DayIndex:           -1,
	}

	session, result := NewWorkoutSession(input, "session-id")

	if result.Valid {
		t.Error("NewWorkoutSession with invalid day_index returned valid result")
	}
	if session != nil {
		t.Error("NewWorkoutSession with invalid input returned non-nil session")
	}
}

func TestNewWorkoutSession_MultipleErrors(t *testing.T) {
	input := NewWorkoutSessionInput{
		UserProgramStateID: "",
		WeekNumber:         0,
		DayIndex:           -1,
	}

	session, result := NewWorkoutSession(input, "")

	if result.Valid {
		t.Error("NewWorkoutSession with multiple errors returned valid result")
	}
	if session != nil {
		t.Error("NewWorkoutSession with invalid input returned non-nil session")
	}
	if len(result.Errors) != 4 {
		t.Errorf("Expected 4 errors, got %d", len(result.Errors))
	}
}

// ==================== Complete Tests ====================

func TestComplete_Success(t *testing.T) {
	session := &WorkoutSession{
		ID:                 "session-id",
		UserProgramStateID: "state-1",
		WeekNumber:         1,
		DayIndex:           0,
		Status:             StatusInProgress,
	}

	err := session.Complete()

	if err != nil {
		t.Errorf("Complete() = %v, want nil", err)
	}
	if session.Status != StatusCompleted {
		t.Errorf("session.Status = %q, want %q", session.Status, StatusCompleted)
	}
	if session.FinishedAt == nil {
		t.Error("session.FinishedAt should be set after completion")
	}
}

func TestComplete_AlreadyCompleted(t *testing.T) {
	session := &WorkoutSession{
		ID:                 "session-id",
		UserProgramStateID: "state-1",
		WeekNumber:         1,
		DayIndex:           0,
		Status:             StatusCompleted,
	}

	err := session.Complete()

	if !errors.Is(err, ErrAlreadyCompleted) {
		t.Errorf("Complete() = %v, want %v", err, ErrAlreadyCompleted)
	}
}

func TestComplete_AlreadyAbandoned(t *testing.T) {
	session := &WorkoutSession{
		ID:                 "session-id",
		UserProgramStateID: "state-1",
		WeekNumber:         1,
		DayIndex:           0,
		Status:             StatusAbandoned,
	}

	err := session.Complete()

	if !errors.Is(err, ErrNotInProgress) {
		t.Errorf("Complete() = %v, want %v", err, ErrNotInProgress)
	}
}

// ==================== Abandon Tests ====================

func TestAbandon_Success(t *testing.T) {
	session := &WorkoutSession{
		ID:                 "session-id",
		UserProgramStateID: "state-1",
		WeekNumber:         1,
		DayIndex:           0,
		Status:             StatusInProgress,
	}

	err := session.Abandon()

	if err != nil {
		t.Errorf("Abandon() = %v, want nil", err)
	}
	if session.Status != StatusAbandoned {
		t.Errorf("session.Status = %q, want %q", session.Status, StatusAbandoned)
	}
	if session.FinishedAt == nil {
		t.Error("session.FinishedAt should be set after abandonment")
	}
}

func TestAbandon_AlreadyAbandoned(t *testing.T) {
	session := &WorkoutSession{
		ID:                 "session-id",
		UserProgramStateID: "state-1",
		WeekNumber:         1,
		DayIndex:           0,
		Status:             StatusAbandoned,
	}

	err := session.Abandon()

	if !errors.Is(err, ErrAlreadyAbandoned) {
		t.Errorf("Abandon() = %v, want %v", err, ErrAlreadyAbandoned)
	}
}

func TestAbandon_AlreadyCompleted(t *testing.T) {
	session := &WorkoutSession{
		ID:                 "session-id",
		UserProgramStateID: "state-1",
		WeekNumber:         1,
		DayIndex:           0,
		Status:             StatusCompleted,
	}

	err := session.Abandon()

	if !errors.Is(err, ErrNotInProgress) {
		t.Errorf("Abandon() = %v, want %v", err, ErrNotInProgress)
	}
}

// ==================== IsActive Tests ====================

func TestIsActive_InProgress(t *testing.T) {
	session := &WorkoutSession{Status: StatusInProgress}
	if !session.IsActive() {
		t.Error("IsActive() = false, want true for IN_PROGRESS")
	}
}

func TestIsActive_Completed(t *testing.T) {
	session := &WorkoutSession{Status: StatusCompleted}
	if session.IsActive() {
		t.Error("IsActive() = true, want false for COMPLETED")
	}
}

func TestIsActive_Abandoned(t *testing.T) {
	session := &WorkoutSession{Status: StatusAbandoned}
	if session.IsActive() {
		t.Error("IsActive() = true, want false for ABANDONED")
	}
}

// ==================== IsFinished Tests ====================

func TestIsFinished_InProgress(t *testing.T) {
	session := &WorkoutSession{Status: StatusInProgress}
	if session.IsFinished() {
		t.Error("IsFinished() = true, want false for IN_PROGRESS")
	}
}

func TestIsFinished_Completed(t *testing.T) {
	session := &WorkoutSession{Status: StatusCompleted}
	if !session.IsFinished() {
		t.Error("IsFinished() = false, want true for COMPLETED")
	}
}

func TestIsFinished_Abandoned(t *testing.T) {
	session := &WorkoutSession{Status: StatusAbandoned}
	if !session.IsFinished() {
		t.Error("IsFinished() = false, want true for ABANDONED")
	}
}

// ==================== WorkoutSession.Validate Tests ====================

func TestWorkoutSession_Validate_Valid(t *testing.T) {
	session := &WorkoutSession{
		ID:                 "session-id",
		UserProgramStateID: "state-1",
		WeekNumber:         1,
		DayIndex:           0,
		Status:             StatusInProgress,
	}

	result := session.Validate()

	if !result.Valid {
		t.Errorf("Validate returned invalid result for valid session: %v", result.Errors)
	}
}

func TestWorkoutSession_Validate_InvalidID(t *testing.T) {
	session := &WorkoutSession{
		ID:                 "",
		UserProgramStateID: "state-1",
		WeekNumber:         1,
		DayIndex:           0,
		Status:             StatusInProgress,
	}

	result := session.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for session with invalid id")
	}
}

func TestWorkoutSession_Validate_InvalidUserProgramStateID(t *testing.T) {
	session := &WorkoutSession{
		ID:                 "session-id",
		UserProgramStateID: "",
		WeekNumber:         1,
		DayIndex:           0,
		Status:             StatusInProgress,
	}

	result := session.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for session with invalid user_program_state_id")
	}
}

func TestWorkoutSession_Validate_InvalidWeekNumber(t *testing.T) {
	session := &WorkoutSession{
		ID:                 "session-id",
		UserProgramStateID: "state-1",
		WeekNumber:         0,
		DayIndex:           0,
		Status:             StatusInProgress,
	}

	result := session.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for session with invalid week_number")
	}
}

func TestWorkoutSession_Validate_InvalidDayIndex(t *testing.T) {
	session := &WorkoutSession{
		ID:                 "session-id",
		UserProgramStateID: "state-1",
		WeekNumber:         1,
		DayIndex:           -1,
		Status:             StatusInProgress,
	}

	result := session.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for session with invalid day_index")
	}
}

func TestWorkoutSession_Validate_InvalidStatus(t *testing.T) {
	session := &WorkoutSession{
		ID:                 "session-id",
		UserProgramStateID: "state-1",
		WeekNumber:         1,
		DayIndex:           0,
		Status:             Status("INVALID"),
	}

	result := session.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for session with invalid status")
	}
}

func TestWorkoutSession_Validate_MultipleErrors(t *testing.T) {
	session := &WorkoutSession{
		ID:                 "",
		UserProgramStateID: "",
		WeekNumber:         0,
		DayIndex:           -1,
		Status:             Status("INVALID"),
	}

	result := session.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for session with multiple errors")
	}
	if len(result.Errors) != 5 {
		t.Errorf("Expected 5 errors, got %d", len(result.Errors))
	}
}

// ==================== ValidationResult Tests ====================

func TestValidationResult_Error_Valid(t *testing.T) {
	result := NewValidationResult()

	err := result.Error()
	if err != nil {
		t.Errorf("ValidationResult.Error() = %v, want nil", err)
	}
}

func TestValidationResult_Error_Invalid(t *testing.T) {
	result := NewValidationResult()
	result.AddError(ErrIDRequired)
	result.AddError(ErrUserProgramStateIDRequired)

	err := result.Error()
	if err == nil {
		t.Error("ValidationResult.Error() = nil, want error")
	}

	errMsg := err.Error()
	if errMsg == "" {
		t.Error("Error message should not be empty")
	}
}

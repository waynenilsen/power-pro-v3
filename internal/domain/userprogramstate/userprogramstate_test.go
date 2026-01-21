package userprogramstate

import (
	"errors"
	"testing"
)

// ==================== UserID Validation Tests ====================

func TestValidateUserID_Valid(t *testing.T) {
	tests := []struct {
		name   string
		userID string
	}{
		{"uuid", "550e8400-e29b-41d4-a716-446655440000"},
		{"simple id", "user-1"},
		{"short id", "u"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserID(tt.userID)
			if err != nil {
				t.Errorf("ValidateUserID(%q) = %v, want nil", tt.userID, err)
			}
		})
	}
}

func TestValidateUserID_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		expectedErr error
	}{
		{"empty string", "", ErrUserIDRequired},
		{"only spaces", "   ", ErrUserIDRequired},
		{"only tabs", "\t\t", ErrUserIDRequired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserID(tt.userID)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateUserID(%q) = %v, want %v", tt.userID, err, tt.expectedErr)
			}
		})
	}
}

// ==================== ProgramID Validation Tests ====================

func TestValidateProgramID_Valid(t *testing.T) {
	tests := []struct {
		name      string
		programID string
	}{
		{"uuid", "550e8400-e29b-41d4-a716-446655440000"},
		{"simple id", "program-1"},
		{"short id", "p"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProgramID(tt.programID)
			if err != nil {
				t.Errorf("ValidateProgramID(%q) = %v, want nil", tt.programID, err)
			}
		})
	}
}

func TestValidateProgramID_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		programID   string
		expectedErr error
	}{
		{"empty string", "", ErrProgramIDRequired},
		{"only spaces", "   ", ErrProgramIDRequired},
		{"only tabs", "\t\t", ErrProgramIDRequired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProgramID(tt.programID)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateProgramID(%q) = %v, want %v", tt.programID, err, tt.expectedErr)
			}
		})
	}
}

// ==================== CurrentWeek Validation Tests ====================

func TestValidateCurrentWeek_Valid(t *testing.T) {
	tests := []struct {
		name string
		week int
	}{
		{"week 1", 1},
		{"week 10", 10},
		{"large week", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCurrentWeek(tt.week)
			if err != nil {
				t.Errorf("ValidateCurrentWeek(%d) = %v, want nil", tt.week, err)
			}
		})
	}
}

func TestValidateCurrentWeek_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		week        int
		expectedErr error
	}{
		{"zero", 0, ErrCurrentWeekInvalid},
		{"negative", -1, ErrCurrentWeekInvalid},
		{"very negative", -100, ErrCurrentWeekInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCurrentWeek(tt.week)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateCurrentWeek(%d) = %v, want %v", tt.week, err, tt.expectedErr)
			}
		})
	}
}

// ==================== CurrentCycleIteration Validation Tests ====================

func TestValidateCurrentCycleIteration_Valid(t *testing.T) {
	tests := []struct {
		name      string
		iteration int
	}{
		{"iteration 1", 1},
		{"iteration 5", 5},
		{"large iteration", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCurrentCycleIteration(tt.iteration)
			if err != nil {
				t.Errorf("ValidateCurrentCycleIteration(%d) = %v, want nil", tt.iteration, err)
			}
		})
	}
}

func TestValidateCurrentCycleIteration_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		iteration   int
		expectedErr error
	}{
		{"zero", 0, ErrCurrentCycleIterationInvalid},
		{"negative", -1, ErrCurrentCycleIterationInvalid},
		{"very negative", -100, ErrCurrentCycleIterationInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCurrentCycleIteration(tt.iteration)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateCurrentCycleIteration(%d) = %v, want %v", tt.iteration, err, tt.expectedErr)
			}
		})
	}
}

// ==================== CurrentDayIndex Validation Tests ====================

func TestValidateCurrentDayIndex_Valid(t *testing.T) {
	tests := []struct {
		name     string
		dayIndex *int
	}{
		{"nil", nil},
		{"zero", ptrInt(0)},
		{"positive", ptrInt(5)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCurrentDayIndex(tt.dayIndex)
			if err != nil {
				t.Errorf("ValidateCurrentDayIndex(%v) = %v, want nil", tt.dayIndex, err)
			}
		})
	}
}

func TestValidateCurrentDayIndex_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		dayIndex    *int
		expectedErr error
	}{
		{"negative", ptrInt(-1), ErrCurrentDayIndexInvalid},
		{"very negative", ptrInt(-100), ErrCurrentDayIndexInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCurrentDayIndex(tt.dayIndex)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateCurrentDayIndex(%v) = %v, want %v", tt.dayIndex, err, tt.expectedErr)
			}
		})
	}
}

// ==================== EnrollUser Tests ====================

func TestEnrollUser_ValidInput(t *testing.T) {
	input := EnrollUserInput{
		UserID:    "user-1",
		ProgramID: "program-1",
	}

	state, result := EnrollUser(input, "state-id")

	if !result.Valid {
		t.Errorf("EnrollUser returned invalid result: %v", result.Errors)
	}
	if state == nil {
		t.Fatal("EnrollUser returned nil state")
	}
	if state.ID != "state-id" {
		t.Errorf("state.ID = %q, want %q", state.ID, "state-id")
	}
	if state.UserID != "user-1" {
		t.Errorf("state.UserID = %q, want %q", state.UserID, "user-1")
	}
	if state.ProgramID != "program-1" {
		t.Errorf("state.ProgramID = %q, want %q", state.ProgramID, "program-1")
	}
	if state.CurrentWeek != 1 {
		t.Errorf("state.CurrentWeek = %d, want %d", state.CurrentWeek, 1)
	}
	if state.CurrentCycleIteration != 1 {
		t.Errorf("state.CurrentCycleIteration = %d, want %d", state.CurrentCycleIteration, 1)
	}
	if state.CurrentDayIndex != nil {
		t.Errorf("state.CurrentDayIndex = %v, want nil", state.CurrentDayIndex)
	}
}

func TestEnrollUser_TrimsWhitespace(t *testing.T) {
	input := EnrollUserInput{
		UserID:    "  user-1  ",
		ProgramID: "  program-1  ",
	}

	state, result := EnrollUser(input, "state-id")

	if !result.Valid {
		t.Errorf("EnrollUser returned invalid result: %v", result.Errors)
	}
	if state.UserID != "user-1" {
		t.Errorf("state.UserID = %q, want %q (trimmed)", state.UserID, "user-1")
	}
	if state.ProgramID != "program-1" {
		t.Errorf("state.ProgramID = %q, want %q (trimmed)", state.ProgramID, "program-1")
	}
}

func TestEnrollUser_EmptyUserID(t *testing.T) {
	input := EnrollUserInput{
		UserID:    "",
		ProgramID: "program-1",
	}

	state, result := EnrollUser(input, "state-id")

	if result.Valid {
		t.Error("EnrollUser with empty user_id returned valid result")
	}
	if state != nil {
		t.Error("EnrollUser with invalid input returned non-nil state")
	}
}

func TestEnrollUser_EmptyProgramID(t *testing.T) {
	input := EnrollUserInput{
		UserID:    "user-1",
		ProgramID: "",
	}

	state, result := EnrollUser(input, "state-id")

	if result.Valid {
		t.Error("EnrollUser with empty program_id returned valid result")
	}
	if state != nil {
		t.Error("EnrollUser with invalid input returned non-nil state")
	}
}

func TestEnrollUser_MultipleErrors(t *testing.T) {
	input := EnrollUserInput{
		UserID:    "",
		ProgramID: "",
	}

	state, result := EnrollUser(input, "state-id")

	if result.Valid {
		t.Error("EnrollUser with multiple errors returned valid result")
	}
	if state != nil {
		t.Error("EnrollUser with invalid input returned non-nil state")
	}
	if len(result.Errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(result.Errors))
	}
}

// ==================== UpdatePosition Tests ====================

func TestUpdatePosition_UpdateWeek(t *testing.T) {
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
	}

	newWeek := 3
	input := UpdatePositionInput{CurrentWeek: &newWeek}

	result := UpdatePosition(state, input)

	if !result.Valid {
		t.Errorf("UpdatePosition returned invalid result: %v", result.Errors)
	}
	if state.CurrentWeek != 3 {
		t.Errorf("state.CurrentWeek = %d, want %d", state.CurrentWeek, 3)
	}
}

func TestUpdatePosition_UpdateCycleIteration(t *testing.T) {
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
	}

	newIteration := 2
	input := UpdatePositionInput{CurrentCycleIteration: &newIteration}

	result := UpdatePosition(state, input)

	if !result.Valid {
		t.Errorf("UpdatePosition returned invalid result: %v", result.Errors)
	}
	if state.CurrentCycleIteration != 2 {
		t.Errorf("state.CurrentCycleIteration = %d, want %d", state.CurrentCycleIteration, 2)
	}
}

func TestUpdatePosition_UpdateDayIndex(t *testing.T) {
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
	}

	newDayIndex := 2
	dayIndexPtr := &newDayIndex
	input := UpdatePositionInput{CurrentDayIndex: &dayIndexPtr}

	result := UpdatePosition(state, input)

	if !result.Valid {
		t.Errorf("UpdatePosition returned invalid result: %v", result.Errors)
	}
	if state.CurrentDayIndex == nil || *state.CurrentDayIndex != 2 {
		t.Errorf("state.CurrentDayIndex = %v, want %d", state.CurrentDayIndex, 2)
	}
}

func TestUpdatePosition_ClearDayIndex(t *testing.T) {
	dayIndex := 2
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
		CurrentDayIndex:       &dayIndex,
	}

	var nilDayIndex *int = nil
	input := UpdatePositionInput{CurrentDayIndex: &nilDayIndex}

	result := UpdatePosition(state, input)

	if !result.Valid {
		t.Errorf("UpdatePosition returned invalid result: %v", result.Errors)
	}
	if state.CurrentDayIndex != nil {
		t.Errorf("state.CurrentDayIndex = %v, want nil", state.CurrentDayIndex)
	}
}

func TestUpdatePosition_InvalidWeek(t *testing.T) {
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
	}
	originalWeek := state.CurrentWeek

	invalidWeek := 0
	input := UpdatePositionInput{CurrentWeek: &invalidWeek}

	result := UpdatePosition(state, input)

	if result.Valid {
		t.Error("UpdatePosition with invalid week returned valid result")
	}
	if state.CurrentWeek != originalWeek {
		t.Errorf("state.CurrentWeek was changed despite validation failure")
	}
}

func TestUpdatePosition_InvalidCycleIteration(t *testing.T) {
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
	}
	originalIteration := state.CurrentCycleIteration

	invalidIteration := -1
	input := UpdatePositionInput{CurrentCycleIteration: &invalidIteration}

	result := UpdatePosition(state, input)

	if result.Valid {
		t.Error("UpdatePosition with invalid cycle iteration returned valid result")
	}
	if state.CurrentCycleIteration != originalIteration {
		t.Errorf("state.CurrentCycleIteration was changed despite validation failure")
	}
}

func TestUpdatePosition_InvalidDayIndex(t *testing.T) {
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
	}

	invalidDayIndex := -1
	invalidDayIndexPtr := &invalidDayIndex
	input := UpdatePositionInput{CurrentDayIndex: &invalidDayIndexPtr}

	result := UpdatePosition(state, input)

	if result.Valid {
		t.Error("UpdatePosition with invalid day index returned valid result")
	}
}

func TestUpdatePosition_NoChanges(t *testing.T) {
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
	}

	input := UpdatePositionInput{} // No changes

	result := UpdatePosition(state, input)

	if !result.Valid {
		t.Errorf("UpdatePosition with no changes returned invalid result: %v", result.Errors)
	}
}

func TestUpdatePosition_MultipleUpdates(t *testing.T) {
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
	}

	newWeek := 3
	newIteration := 2
	input := UpdatePositionInput{
		CurrentWeek:           &newWeek,
		CurrentCycleIteration: &newIteration,
	}

	result := UpdatePosition(state, input)

	if !result.Valid {
		t.Errorf("UpdatePosition returned invalid result: %v", result.Errors)
	}
	if state.CurrentWeek != 3 {
		t.Errorf("state.CurrentWeek = %d, want %d", state.CurrentWeek, 3)
	}
	if state.CurrentCycleIteration != 2 {
		t.Errorf("state.CurrentCycleIteration = %d, want %d", state.CurrentCycleIteration, 2)
	}
}

// ==================== UserProgramState.Validate Tests ====================

func TestUserProgramState_Validate_Valid(t *testing.T) {
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
	}

	result := state.Validate()

	if !result.Valid {
		t.Errorf("Validate returned invalid result for valid state: %v", result.Errors)
	}
}

func TestUserProgramState_Validate_InvalidUserID(t *testing.T) {
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
	}

	result := state.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for state with invalid user_id")
	}
}

func TestUserProgramState_Validate_InvalidProgramID(t *testing.T) {
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
	}

	result := state.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for state with invalid program_id")
	}
}

func TestUserProgramState_Validate_InvalidCurrentWeek(t *testing.T) {
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           0,
		CurrentCycleIteration: 1,
	}

	result := state.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for state with invalid current_week")
	}
}

func TestUserProgramState_Validate_InvalidCycleIteration(t *testing.T) {
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: -1,
	}

	result := state.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for state with invalid current_cycle_iteration")
	}
}

func TestUserProgramState_Validate_InvalidDayIndex(t *testing.T) {
	dayIndex := -1
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
		CurrentDayIndex:       &dayIndex,
	}

	result := state.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for state with invalid current_day_index")
	}
}

func TestUserProgramState_Validate_MultipleErrors(t *testing.T) {
	dayIndex := -1
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "",
		ProgramID:             "",
		CurrentWeek:           0,
		CurrentCycleIteration: 0,
		CurrentDayIndex:       &dayIndex,
	}

	result := state.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for state with multiple errors")
	}
	if len(result.Errors) < 5 {
		t.Errorf("Expected at least 5 errors, got %d", len(result.Errors))
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
	result.AddError(ErrUserIDRequired)
	result.AddError(ErrProgramIDRequired)

	err := result.Error()
	if err == nil {
		t.Error("ValidationResult.Error() = nil, want error")
	}

	errMsg := err.Error()
	if errMsg == "" {
		t.Error("Error message should not be empty")
	}
}

// ==================== Helper Functions for Tests ====================

func ptrInt(i int) *int {
	return &i
}

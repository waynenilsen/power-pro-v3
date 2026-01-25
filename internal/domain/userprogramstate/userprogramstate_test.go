package userprogramstate

import (
	"errors"
	"testing"
	"time"
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
	// Verify initial status values
	if state.EnrollmentStatus != EnrollmentStatusActive {
		t.Errorf("state.EnrollmentStatus = %q, want %q", state.EnrollmentStatus, EnrollmentStatusActive)
	}
	if state.CycleStatus != CycleStatusPending {
		t.Errorf("state.CycleStatus = %q, want %q", state.CycleStatus, CycleStatusPending)
	}
	if state.WeekStatus != WeekStatusPending {
		t.Errorf("state.WeekStatus = %q, want %q", state.WeekStatus, WeekStatusPending)
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
		EnrollmentStatus:      EnrollmentStatusActive,
		CycleStatus:           CycleStatusPending,
		WeekStatus:            WeekStatusPending,
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
		EnrollmentStatus:      EnrollmentStatusActive,
		CycleStatus:           CycleStatusPending,
		WeekStatus:            WeekStatusPending,
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
		EnrollmentStatus:      EnrollmentStatusActive,
		CycleStatus:           CycleStatusPending,
		WeekStatus:            WeekStatusPending,
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
		EnrollmentStatus:      EnrollmentStatusActive,
		CycleStatus:           CycleStatusPending,
		WeekStatus:            WeekStatusPending,
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
		EnrollmentStatus:      EnrollmentStatusActive,
		CycleStatus:           CycleStatusPending,
		WeekStatus:            WeekStatusPending,
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
		EnrollmentStatus:      EnrollmentStatusActive,
		CycleStatus:           CycleStatusPending,
		WeekStatus:            WeekStatusPending,
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
		EnrollmentStatus:      EnrollmentStatus("INVALID"),
		CycleStatus:           CycleStatus("INVALID"),
		WeekStatus:            WeekStatus("INVALID"),
	}

	result := state.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for state with multiple errors")
	}
	if len(result.Errors) < 8 {
		t.Errorf("Expected at least 8 errors, got %d", len(result.Errors))
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

// ==================== AdvanceState Tests ====================

func TestAdvanceState_FirstAdvancement(t *testing.T) {
	// First advancement from nil day index
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
		CurrentDayIndex:       nil, // nil = not yet started
	}

	ctx := AdvancementContext{
		DaysInCurrentWeek: 3,
		CycleLengthWeeks:  4,
	}

	result, validation := AdvanceState(state, ctx)

	if !validation.Valid {
		t.Fatalf("AdvanceState returned invalid result: %v", validation.Errors)
	}
	if result == nil {
		t.Fatal("AdvanceState returned nil result")
	}
	if result.NewState.CurrentDayIndex == nil || *result.NewState.CurrentDayIndex != 1 {
		t.Errorf("CurrentDayIndex = %v, want 1", result.NewState.CurrentDayIndex)
	}
	if result.NewState.CurrentWeek != 1 {
		t.Errorf("CurrentWeek = %d, want 1", result.NewState.CurrentWeek)
	}
	if result.CycleCompleted {
		t.Error("CycleCompleted should be false for first advancement")
	}
}

func TestAdvanceState_AdvanceWithinWeek(t *testing.T) {
	dayIdx := 0
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
		CurrentDayIndex:       &dayIdx,
	}

	ctx := AdvancementContext{
		DaysInCurrentWeek: 3,
		CycleLengthWeeks:  4,
	}

	result, validation := AdvanceState(state, ctx)

	if !validation.Valid {
		t.Fatalf("AdvanceState returned invalid result: %v", validation.Errors)
	}
	if result.NewState.CurrentDayIndex == nil || *result.NewState.CurrentDayIndex != 1 {
		t.Errorf("CurrentDayIndex = %v, want 1", result.NewState.CurrentDayIndex)
	}
	if result.NewState.CurrentWeek != 1 {
		t.Errorf("CurrentWeek = %d, want 1 (should stay same)", result.NewState.CurrentWeek)
	}
	if result.CycleCompleted {
		t.Error("CycleCompleted should be false when advancing within week")
	}
}

func TestAdvanceState_AdvanceToNextWeek(t *testing.T) {
	// Day index is at 2 (last day of 3-day week), should advance to next week
	dayIdx := 2
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
		CurrentDayIndex:       &dayIdx, // Last day of 3-day week
	}

	ctx := AdvancementContext{
		DaysInCurrentWeek: 3,
		CycleLengthWeeks:  4,
	}

	result, validation := AdvanceState(state, ctx)

	if !validation.Valid {
		t.Fatalf("AdvanceState returned invalid result: %v", validation.Errors)
	}
	// After incrementing day index becomes 3, which >= 3 days, so it wraps
	if result.NewState.CurrentDayIndex == nil || *result.NewState.CurrentDayIndex != 0 {
		t.Errorf("CurrentDayIndex = %v, want 0 (reset for new week)", result.NewState.CurrentDayIndex)
	}
	if result.NewState.CurrentWeek != 2 {
		t.Errorf("CurrentWeek = %d, want 2 (advanced to next week)", result.NewState.CurrentWeek)
	}
	if result.CycleCompleted {
		t.Error("CycleCompleted should be false when advancing to next week")
	}
}

func TestAdvanceState_CycleCompletion(t *testing.T) {
	// Last day of last week of 4-week cycle
	dayIdx := 2
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           4, // Last week
		CurrentCycleIteration: 1,
		CurrentDayIndex:       &dayIdx, // Last day
	}

	ctx := AdvancementContext{
		DaysInCurrentWeek: 3, // 3-day week
		CycleLengthWeeks:  4,
	}

	result, validation := AdvanceState(state, ctx)

	if !validation.Valid {
		t.Fatalf("AdvanceState returned invalid result: %v", validation.Errors)
	}
	// After incrementing: day index becomes 3 >= 3, wraps to 0
	// Week increments to 5 > 4, wraps to 1
	// Cycle iteration increments to 2
	if result.NewState.CurrentDayIndex == nil || *result.NewState.CurrentDayIndex != 0 {
		t.Errorf("CurrentDayIndex = %v, want 0 (reset for new cycle)", result.NewState.CurrentDayIndex)
	}
	if result.NewState.CurrentWeek != 1 {
		t.Errorf("CurrentWeek = %d, want 1 (reset for new cycle)", result.NewState.CurrentWeek)
	}
	if result.NewState.CurrentCycleIteration != 2 {
		t.Errorf("CurrentCycleIteration = %d, want 2 (incremented)", result.NewState.CurrentCycleIteration)
	}
	if !result.CycleCompleted {
		t.Error("CycleCompleted should be true when completing cycle")
	}
}

func TestAdvanceState_MultipleAdvancementsThroughCycle(t *testing.T) {
	// Simulate advancing through an entire 2-week cycle with 2 days per week
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
		CurrentDayIndex:       nil, // Starting fresh
	}

	ctx := AdvancementContext{
		DaysInCurrentWeek: 2,
		CycleLengthWeeks:  2,
	}

	// Advancement 1: nil -> day 1
	result, _ := AdvanceState(state, ctx)
	if *result.NewState.CurrentDayIndex != 1 || result.NewState.CurrentWeek != 1 {
		t.Errorf("Advancement 1: day=%v week=%d, want day=1 week=1",
			*result.NewState.CurrentDayIndex, result.NewState.CurrentWeek)
	}
	if result.CycleCompleted {
		t.Error("Advancement 1: unexpected cycle completion")
	}

	// Advancement 2: day 1 -> day 0, week 2
	result, _ = AdvanceState(result.NewState, ctx)
	if *result.NewState.CurrentDayIndex != 0 || result.NewState.CurrentWeek != 2 {
		t.Errorf("Advancement 2: day=%v week=%d, want day=0 week=2",
			*result.NewState.CurrentDayIndex, result.NewState.CurrentWeek)
	}
	if result.CycleCompleted {
		t.Error("Advancement 2: unexpected cycle completion")
	}

	// Advancement 3: day 0, week 2 -> day 1, week 2
	result, _ = AdvanceState(result.NewState, ctx)
	if *result.NewState.CurrentDayIndex != 1 || result.NewState.CurrentWeek != 2 {
		t.Errorf("Advancement 3: day=%v week=%d, want day=1 week=2",
			*result.NewState.CurrentDayIndex, result.NewState.CurrentWeek)
	}
	if result.CycleCompleted {
		t.Error("Advancement 3: unexpected cycle completion")
	}

	// Advancement 4: day 1, week 2 -> day 0, week 1, cycle 2 (CYCLE COMPLETE!)
	result, _ = AdvanceState(result.NewState, ctx)
	if *result.NewState.CurrentDayIndex != 0 || result.NewState.CurrentWeek != 1 {
		t.Errorf("Advancement 4: day=%v week=%d, want day=0 week=1",
			*result.NewState.CurrentDayIndex, result.NewState.CurrentWeek)
	}
	if result.NewState.CurrentCycleIteration != 2 {
		t.Errorf("Advancement 4: cycle=%d, want 2", result.NewState.CurrentCycleIteration)
	}
	if !result.CycleCompleted {
		t.Error("Advancement 4: expected cycle completion")
	}
}

func TestAdvanceState_SingleDaySingleWeekCycle(t *testing.T) {
	// Edge case: 1 day per week, 1 week cycle
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
		CurrentDayIndex:       ptrInt(0),
	}

	ctx := AdvancementContext{
		DaysInCurrentWeek: 1,
		CycleLengthWeeks:  1,
	}

	result, validation := AdvanceState(state, ctx)

	if !validation.Valid {
		t.Fatalf("AdvanceState returned invalid result: %v", validation.Errors)
	}
	// Should complete cycle every advancement
	if !result.CycleCompleted {
		t.Error("CycleCompleted should be true for 1-day, 1-week cycle")
	}
	if result.NewState.CurrentCycleIteration != 2 {
		t.Errorf("CurrentCycleIteration = %d, want 2", result.NewState.CurrentCycleIteration)
	}
	if result.NewState.CurrentWeek != 1 {
		t.Errorf("CurrentWeek = %d, want 1", result.NewState.CurrentWeek)
	}
	if *result.NewState.CurrentDayIndex != 0 {
		t.Errorf("CurrentDayIndex = %d, want 0", *result.NewState.CurrentDayIndex)
	}
}

func TestAdvanceState_DoesNotMutateOriginal(t *testing.T) {
	dayIdx := 0
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
		CurrentDayIndex:       &dayIdx,
	}

	ctx := AdvancementContext{
		DaysInCurrentWeek: 3,
		CycleLengthWeeks:  4,
	}

	_, _ = AdvanceState(state, ctx)

	// Original state should be unchanged
	if state.CurrentWeek != 1 {
		t.Errorf("Original state.CurrentWeek was mutated: %d", state.CurrentWeek)
	}
	if state.CurrentCycleIteration != 1 {
		t.Errorf("Original state.CurrentCycleIteration was mutated: %d", state.CurrentCycleIteration)
	}
	if *state.CurrentDayIndex != 0 {
		t.Errorf("Original state.CurrentDayIndex was mutated: %d", *state.CurrentDayIndex)
	}
}

func TestAdvanceState_InvalidContext_ZeroDays(t *testing.T) {
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
		CurrentDayIndex:       ptrInt(0),
	}

	ctx := AdvancementContext{
		DaysInCurrentWeek: 0, // Invalid
		CycleLengthWeeks:  4,
	}

	result, validation := AdvanceState(state, ctx)

	if validation.Valid {
		t.Error("AdvanceState should fail with 0 days in week")
	}
	if result != nil {
		t.Error("AdvanceState should return nil result on validation failure")
	}
}

func TestAdvanceState_InvalidContext_ZeroWeeks(t *testing.T) {
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
		CurrentDayIndex:       ptrInt(0),
	}

	ctx := AdvancementContext{
		DaysInCurrentWeek: 3,
		CycleLengthWeeks:  0, // Invalid
	}

	result, validation := AdvanceState(state, ctx)

	if validation.Valid {
		t.Error("AdvanceState should fail with 0 weeks in cycle")
	}
	if result != nil {
		t.Error("AdvanceState should return nil result on validation failure")
	}
}

// ==================== MeetDate Validation Tests ====================

func TestValidateMeetDate_Valid(t *testing.T) {
	tests := []struct {
		name     string
		meetDate *time.Time
	}{
		{"nil", nil},
		{"future date", ptrTime(time.Now().Add(24 * time.Hour))},
		{"far future date", ptrTime(time.Now().Add(365 * 24 * time.Hour))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMeetDate(tt.meetDate)
			if err != nil {
				t.Errorf("ValidateMeetDate(%v) = %v, want nil", tt.meetDate, err)
			}
		})
	}
}

func TestValidateMeetDate_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		meetDate    *time.Time
		expectedErr error
	}{
		{"past date", ptrTime(time.Now().Add(-24 * time.Hour)), ErrMeetDateInPast},
		{"now (edge case)", ptrTime(time.Now().Add(-1 * time.Second)), ErrMeetDateInPast},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMeetDate(tt.meetDate)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateMeetDate(%v) = %v, want %v", tt.meetDate, err, tt.expectedErr)
			}
		})
	}
}

// ==================== EnrollmentStatus Validation Tests ====================

func TestValidateEnrollmentStatus_Valid(t *testing.T) {
	tests := []struct {
		name   string
		status EnrollmentStatus
	}{
		{"active", EnrollmentStatusActive},
		{"between cycles", EnrollmentStatusBetweenCycles},
		{"quit", EnrollmentStatusQuit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEnrollmentStatus(tt.status)
			if err != nil {
				t.Errorf("ValidateEnrollmentStatus(%q) = %v, want nil", tt.status, err)
			}
		})
	}
}

func TestValidateEnrollmentStatus_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		status      EnrollmentStatus
		expectedErr error
	}{
		{"empty string", EnrollmentStatus(""), ErrInvalidEnrollmentStatus},
		{"invalid value", EnrollmentStatus("INVALID"), ErrInvalidEnrollmentStatus},
		{"lowercase active", EnrollmentStatus("active"), ErrInvalidEnrollmentStatus},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEnrollmentStatus(tt.status)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateEnrollmentStatus(%q) = %v, want %v", tt.status, err, tt.expectedErr)
			}
		})
	}
}

// ==================== CycleStatus Validation Tests ====================

func TestValidateCycleStatus_Valid(t *testing.T) {
	tests := []struct {
		name   string
		status CycleStatus
	}{
		{"pending", CycleStatusPending},
		{"in progress", CycleStatusInProgress},
		{"completed", CycleStatusCompleted},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCycleStatus(tt.status)
			if err != nil {
				t.Errorf("ValidateCycleStatus(%q) = %v, want nil", tt.status, err)
			}
		})
	}
}

func TestValidateCycleStatus_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		status      CycleStatus
		expectedErr error
	}{
		{"empty string", CycleStatus(""), ErrInvalidCycleStatus},
		{"invalid value", CycleStatus("INVALID"), ErrInvalidCycleStatus},
		{"lowercase pending", CycleStatus("pending"), ErrInvalidCycleStatus},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCycleStatus(tt.status)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateCycleStatus(%q) = %v, want %v", tt.status, err, tt.expectedErr)
			}
		})
	}
}

// ==================== WeekStatus Validation Tests ====================

func TestValidateWeekStatus_Valid(t *testing.T) {
	tests := []struct {
		name   string
		status WeekStatus
	}{
		{"pending", WeekStatusPending},
		{"in progress", WeekStatusInProgress},
		{"completed", WeekStatusCompleted},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateWeekStatus(tt.status)
			if err != nil {
				t.Errorf("ValidateWeekStatus(%q) = %v, want nil", tt.status, err)
			}
		})
	}
}

func TestValidateWeekStatus_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		status      WeekStatus
		expectedErr error
	}{
		{"empty string", WeekStatus(""), ErrInvalidWeekStatus},
		{"invalid value", WeekStatus("INVALID"), ErrInvalidWeekStatus},
		{"lowercase pending", WeekStatus("pending"), ErrInvalidWeekStatus},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateWeekStatus(tt.status)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateWeekStatus(%q) = %v, want %v", tt.status, err, tt.expectedErr)
			}
		})
	}
}

// ==================== ScheduleType Validation Tests ====================

func TestValidateScheduleType_Valid(t *testing.T) {
	tests := []struct {
		name         string
		scheduleType ScheduleType
	}{
		{"rotation", ScheduleTypeRotation},
		{"days_out", ScheduleTypeDaysOut},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateScheduleType(tt.scheduleType)
			if err != nil {
				t.Errorf("ValidateScheduleType(%q) = %v, want nil", tt.scheduleType, err)
			}
		})
	}
}

func TestValidateScheduleType_Invalid(t *testing.T) {
	tests := []struct {
		name         string
		scheduleType ScheduleType
		expectedErr  error
	}{
		{"empty string", ScheduleType(""), ErrInvalidScheduleType},
		{"invalid value", ScheduleType("invalid"), ErrInvalidScheduleType},
		{"typo", ScheduleType("rotaiton"), ErrInvalidScheduleType},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateScheduleType(tt.scheduleType)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateScheduleType(%q) = %v, want %v", tt.scheduleType, err, tt.expectedErr)
			}
		})
	}
}

// ==================== EnrollUser with MeetDate Tests ====================

func TestEnrollUser_WithMeetDate(t *testing.T) {
	futureDate := time.Now().Add(30 * 24 * time.Hour)
	input := EnrollUserInput{
		UserID:       "user-1",
		ProgramID:    "program-1",
		MeetDate:     &futureDate,
		ScheduleType: ScheduleTypeDaysOut,
	}

	state, result := EnrollUser(input, "state-id")

	if !result.Valid {
		t.Errorf("EnrollUser returned invalid result: %v", result.Errors)
	}
	if state == nil {
		t.Fatal("EnrollUser returned nil state")
	}
	if state.MeetDate == nil {
		t.Fatal("state.MeetDate is nil, want non-nil")
	}
	if !state.MeetDate.Equal(futureDate) {
		t.Errorf("state.MeetDate = %v, want %v", state.MeetDate, futureDate)
	}
	if state.ScheduleType != ScheduleTypeDaysOut {
		t.Errorf("state.ScheduleType = %q, want %q", state.ScheduleType, ScheduleTypeDaysOut)
	}
}

func TestEnrollUser_WithPastMeetDate(t *testing.T) {
	pastDate := time.Now().Add(-24 * time.Hour)
	input := EnrollUserInput{
		UserID:    "user-1",
		ProgramID: "program-1",
		MeetDate:  &pastDate,
	}

	state, result := EnrollUser(input, "state-id")

	if result.Valid {
		t.Error("EnrollUser with past meet date returned valid result")
	}
	if state != nil {
		t.Error("EnrollUser with invalid input returned non-nil state")
	}
}

func TestEnrollUser_DefaultsToRotationScheduleType(t *testing.T) {
	input := EnrollUserInput{
		UserID:    "user-1",
		ProgramID: "program-1",
		// ScheduleType not set
	}

	state, result := EnrollUser(input, "state-id")

	if !result.Valid {
		t.Errorf("EnrollUser returned invalid result: %v", result.Errors)
	}
	if state == nil {
		t.Fatal("EnrollUser returned nil state")
	}
	if state.ScheduleType != ScheduleTypeRotation {
		t.Errorf("state.ScheduleType = %q, want %q (default)", state.ScheduleType, ScheduleTypeRotation)
	}
}

func TestEnrollUser_InvalidScheduleType(t *testing.T) {
	input := EnrollUserInput{
		UserID:       "user-1",
		ProgramID:    "program-1",
		ScheduleType: ScheduleType("invalid"),
	}

	state, result := EnrollUser(input, "state-id")

	if result.Valid {
		t.Error("EnrollUser with invalid schedule type returned valid result")
	}
	if state != nil {
		t.Error("EnrollUser with invalid input returned non-nil state")
	}
}

// ==================== UserProgramState.Validate with MeetDate Tests ====================

func TestUserProgramState_Validate_WithValidMeetDate(t *testing.T) {
	futureDate := time.Now().Add(30 * 24 * time.Hour)
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
		MeetDate:              &futureDate,
		ScheduleType:          ScheduleTypeDaysOut,
		EnrollmentStatus:      EnrollmentStatusActive,
		CycleStatus:           CycleStatusPending,
		WeekStatus:            WeekStatusPending,
	}

	result := state.Validate()

	if !result.Valid {
		t.Errorf("Validate returned invalid result for valid state: %v", result.Errors)
	}
}

func TestUserProgramState_Validate_WithPastMeetDate(t *testing.T) {
	pastDate := time.Now().Add(-24 * time.Hour)
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
		MeetDate:              &pastDate,
		ScheduleType:          ScheduleTypeDaysOut,
		EnrollmentStatus:      EnrollmentStatusActive,
		CycleStatus:           CycleStatusPending,
		WeekStatus:            WeekStatusPending,
	}

	result := state.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for state with past meet date")
	}
}

func TestUserProgramState_Validate_WithInvalidScheduleType(t *testing.T) {
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
		ScheduleType:          ScheduleType("invalid"),
		EnrollmentStatus:      EnrollmentStatusActive,
		CycleStatus:           CycleStatusPending,
		WeekStatus:            WeekStatusPending,
	}

	result := state.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for state with invalid schedule type")
	}
}

func TestUserProgramState_Validate_WithInvalidEnrollmentStatus(t *testing.T) {
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
		EnrollmentStatus:      EnrollmentStatus("INVALID"),
		CycleStatus:           CycleStatusPending,
		WeekStatus:            WeekStatusPending,
	}

	result := state.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for state with invalid enrollment status")
	}
}

func TestUserProgramState_Validate_WithInvalidCycleStatus(t *testing.T) {
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
		EnrollmentStatus:      EnrollmentStatusActive,
		CycleStatus:           CycleStatus("INVALID"),
		WeekStatus:            WeekStatusPending,
	}

	result := state.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for state with invalid cycle status")
	}
}

func TestUserProgramState_Validate_WithInvalidWeekStatus(t *testing.T) {
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
		EnrollmentStatus:      EnrollmentStatusActive,
		CycleStatus:           CycleStatusPending,
		WeekStatus:            WeekStatus("INVALID"),
	}

	result := state.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for state with invalid week status")
	}
}

// ==================== AdvanceState preserves MeetDate/ScheduleType Tests ====================

func TestAdvanceState_PreservesMeetDateAndScheduleType(t *testing.T) {
	futureDate := time.Now().Add(30 * 24 * time.Hour)
	dayIdx := 0
	state := &UserProgramState{
		ID:                    "state-id",
		UserID:                "user-1",
		ProgramID:             "program-1",
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
		CurrentDayIndex:       &dayIdx,
		MeetDate:              &futureDate,
		ScheduleType:          ScheduleTypeDaysOut,
	}

	ctx := AdvancementContext{
		DaysInCurrentWeek: 3,
		CycleLengthWeeks:  4,
	}

	result, validation := AdvanceState(state, ctx)

	if !validation.Valid {
		t.Fatalf("AdvanceState returned invalid result: %v", validation.Errors)
	}
	if result.NewState.MeetDate == nil {
		t.Fatal("MeetDate was not preserved after advancement")
	}
	if !result.NewState.MeetDate.Equal(futureDate) {
		t.Errorf("MeetDate = %v, want %v", result.NewState.MeetDate, futureDate)
	}
	if result.NewState.ScheduleType != ScheduleTypeDaysOut {
		t.Errorf("ScheduleType = %q, want %q", result.NewState.ScheduleType, ScheduleTypeDaysOut)
	}
	// Ensure original wasn't mutated
	if state.MeetDate == nil || !state.MeetDate.Equal(futureDate) {
		t.Error("Original state.MeetDate was mutated")
	}
}

// ==================== Helper Functions for Tests ====================

func ptrInt(i int) *int {
	return &i
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

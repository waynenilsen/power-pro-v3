package loggedset

import (
	"errors"
	"testing"
)

// ==================== Validation Tests ====================

func TestValidateUserID_Valid(t *testing.T) {
	err := ValidateUserID("user-123")
	if err != nil {
		t.Errorf("ValidateUserID(\"user-123\") = %v, want nil", err)
	}
}

func TestValidateUserID_Empty(t *testing.T) {
	err := ValidateUserID("")
	if !errors.Is(err, ErrUserIDRequired) {
		t.Errorf("ValidateUserID(\"\") = %v, want %v", err, ErrUserIDRequired)
	}
}

func TestValidateSessionID_Valid(t *testing.T) {
	err := ValidateSessionID("session-123")
	if err != nil {
		t.Errorf("ValidateSessionID(\"session-123\") = %v, want nil", err)
	}
}

func TestValidateSessionID_Empty(t *testing.T) {
	err := ValidateSessionID("")
	if !errors.Is(err, ErrSessionIDRequired) {
		t.Errorf("ValidateSessionID(\"\") = %v, want %v", err, ErrSessionIDRequired)
	}
}

func TestValidatePrescriptionID_Valid(t *testing.T) {
	err := ValidatePrescriptionID("prescription-123")
	if err != nil {
		t.Errorf("ValidatePrescriptionID(\"prescription-123\") = %v, want nil", err)
	}
}

func TestValidatePrescriptionID_Empty(t *testing.T) {
	err := ValidatePrescriptionID("")
	if !errors.Is(err, ErrPrescriptionIDRequired) {
		t.Errorf("ValidatePrescriptionID(\"\") = %v, want %v", err, ErrPrescriptionIDRequired)
	}
}

func TestValidateLiftID_Valid(t *testing.T) {
	err := ValidateLiftID("lift-123")
	if err != nil {
		t.Errorf("ValidateLiftID(\"lift-123\") = %v, want nil", err)
	}
}

func TestValidateLiftID_Empty(t *testing.T) {
	err := ValidateLiftID("")
	if !errors.Is(err, ErrLiftIDRequired) {
		t.Errorf("ValidateLiftID(\"\") = %v, want %v", err, ErrLiftIDRequired)
	}
}

func TestValidateSetNumber_Valid(t *testing.T) {
	tests := []int{1, 5, 10, 100}
	for _, n := range tests {
		err := ValidateSetNumber(n)
		if err != nil {
			t.Errorf("ValidateSetNumber(%d) = %v, want nil", n, err)
		}
	}
}

func TestValidateSetNumber_Invalid(t *testing.T) {
	tests := []int{0, -1, -100}
	for _, n := range tests {
		err := ValidateSetNumber(n)
		if !errors.Is(err, ErrSetNumberInvalid) {
			t.Errorf("ValidateSetNumber(%d) = %v, want %v", n, err, ErrSetNumberInvalid)
		}
	}
}

func TestValidateWeight_Valid(t *testing.T) {
	tests := []float64{0, 0.5, 100, 225.5, 500}
	for _, w := range tests {
		err := ValidateWeight(w)
		if err != nil {
			t.Errorf("ValidateWeight(%v) = %v, want nil", w, err)
		}
	}
}

func TestValidateWeight_Invalid(t *testing.T) {
	tests := []float64{-1, -0.5, -100}
	for _, w := range tests {
		err := ValidateWeight(w)
		if !errors.Is(err, ErrWeightInvalid) {
			t.Errorf("ValidateWeight(%v) = %v, want %v", w, err, ErrWeightInvalid)
		}
	}
}

func TestValidateTargetReps_Valid(t *testing.T) {
	tests := []int{1, 5, 10, 100}
	for _, r := range tests {
		err := ValidateTargetReps(r)
		if err != nil {
			t.Errorf("ValidateTargetReps(%d) = %v, want nil", r, err)
		}
	}
}

func TestValidateTargetReps_Invalid(t *testing.T) {
	tests := []int{0, -1, -100}
	for _, r := range tests {
		err := ValidateTargetReps(r)
		if !errors.Is(err, ErrTargetRepsInvalid) {
			t.Errorf("ValidateTargetReps(%d) = %v, want %v", r, err, ErrTargetRepsInvalid)
		}
	}
}

func TestValidateRepsPerformed_Valid(t *testing.T) {
	tests := []int{0, 1, 5, 10, 100}
	for _, r := range tests {
		err := ValidateRepsPerformed(r)
		if err != nil {
			t.Errorf("ValidateRepsPerformed(%d) = %v, want nil", r, err)
		}
	}
}

func TestValidateRepsPerformed_Invalid(t *testing.T) {
	tests := []int{-1, -100}
	for _, r := range tests {
		err := ValidateRepsPerformed(r)
		if !errors.Is(err, ErrRepsPerformedInvalid) {
			t.Errorf("ValidateRepsPerformed(%d) = %v, want %v", r, err, ErrRepsPerformedInvalid)
		}
	}
}

// ==================== NewLoggedSet Tests ====================

func TestNewLoggedSet_ValidInput(t *testing.T) {
	input := CreateLoggedSetInput{
		UserID:         "user-123",
		SessionID:      "session-456",
		PrescriptionID: "prescription-789",
		LiftID:         "lift-abc",
		SetNumber:      1,
		Weight:         225.0,
		TargetReps:     5,
		RepsPerformed:  7,
		IsAMRAP:        true,
	}

	ls, result := NewLoggedSet(input, "test-id")

	if !result.Valid {
		t.Errorf("NewLoggedSet returned invalid result: %v", result.Errors)
	}
	if ls == nil {
		t.Fatal("NewLoggedSet returned nil logged set")
	}
	if ls.ID != "test-id" {
		t.Errorf("ls.ID = %q, want %q", ls.ID, "test-id")
	}
	if ls.UserID != "user-123" {
		t.Errorf("ls.UserID = %q, want %q", ls.UserID, "user-123")
	}
	if ls.SessionID != "session-456" {
		t.Errorf("ls.SessionID = %q, want %q", ls.SessionID, "session-456")
	}
	if ls.PrescriptionID != "prescription-789" {
		t.Errorf("ls.PrescriptionID = %q, want %q", ls.PrescriptionID, "prescription-789")
	}
	if ls.LiftID != "lift-abc" {
		t.Errorf("ls.LiftID = %q, want %q", ls.LiftID, "lift-abc")
	}
	if ls.SetNumber != 1 {
		t.Errorf("ls.SetNumber = %d, want %d", ls.SetNumber, 1)
	}
	if ls.Weight != 225.0 {
		t.Errorf("ls.Weight = %v, want %v", ls.Weight, 225.0)
	}
	if ls.TargetReps != 5 {
		t.Errorf("ls.TargetReps = %d, want %d", ls.TargetReps, 5)
	}
	if ls.RepsPerformed != 7 {
		t.Errorf("ls.RepsPerformed = %d, want %d", ls.RepsPerformed, 7)
	}
	if !ls.IsAMRAP {
		t.Error("ls.IsAMRAP = false, want true")
	}
	if ls.CreatedAt.IsZero() {
		t.Error("ls.CreatedAt is zero time")
	}
}

func TestNewLoggedSet_ZeroWeight(t *testing.T) {
	input := CreateLoggedSetInput{
		UserID:         "user-123",
		SessionID:      "session-456",
		PrescriptionID: "prescription-789",
		LiftID:         "lift-abc",
		SetNumber:      1,
		Weight:         0, // Bodyweight exercise
		TargetReps:     10,
		RepsPerformed:  15,
		IsAMRAP:        true,
	}

	ls, result := NewLoggedSet(input, "test-id")

	if !result.Valid {
		t.Errorf("NewLoggedSet with zero weight returned invalid result: %v", result.Errors)
	}
	if ls == nil {
		t.Fatal("NewLoggedSet returned nil logged set")
	}
}

func TestNewLoggedSet_ZeroRepsPerformed(t *testing.T) {
	input := CreateLoggedSetInput{
		UserID:         "user-123",
		SessionID:      "session-456",
		PrescriptionID: "prescription-789",
		LiftID:         "lift-abc",
		SetNumber:      1,
		Weight:         225.0,
		TargetReps:     5,
		RepsPerformed:  0, // Failed set
		IsAMRAP:        false,
	}

	ls, result := NewLoggedSet(input, "test-id")

	if !result.Valid {
		t.Errorf("NewLoggedSet with zero reps performed returned invalid result: %v", result.Errors)
	}
	if ls == nil {
		t.Fatal("NewLoggedSet returned nil logged set")
	}
}

func TestNewLoggedSet_EmptyUserID(t *testing.T) {
	input := CreateLoggedSetInput{
		UserID:         "",
		SessionID:      "session-456",
		PrescriptionID: "prescription-789",
		LiftID:         "lift-abc",
		SetNumber:      1,
		Weight:         225.0,
		TargetReps:     5,
		RepsPerformed:  5,
	}

	ls, result := NewLoggedSet(input, "test-id")

	if result.Valid {
		t.Error("NewLoggedSet with empty user ID returned valid result")
	}
	if ls != nil {
		t.Error("NewLoggedSet with invalid input returned non-nil logged set")
	}
}

func TestNewLoggedSet_EmptySessionID(t *testing.T) {
	input := CreateLoggedSetInput{
		UserID:         "user-123",
		SessionID:      "",
		PrescriptionID: "prescription-789",
		LiftID:         "lift-abc",
		SetNumber:      1,
		Weight:         225.0,
		TargetReps:     5,
		RepsPerformed:  5,
	}

	ls, result := NewLoggedSet(input, "test-id")

	if result.Valid {
		t.Error("NewLoggedSet with empty session ID returned valid result")
	}
	if ls != nil {
		t.Error("NewLoggedSet with invalid input returned non-nil logged set")
	}
}

func TestNewLoggedSet_MultipleErrors(t *testing.T) {
	input := CreateLoggedSetInput{
		UserID:         "",  // Invalid
		SessionID:      "",  // Invalid
		PrescriptionID: "",  // Invalid
		LiftID:         "",  // Invalid
		SetNumber:      0,   // Invalid
		Weight:         -1,  // Invalid
		TargetReps:     0,   // Invalid
		RepsPerformed:  -1,  // Invalid
	}

	ls, result := NewLoggedSet(input, "test-id")

	if result.Valid {
		t.Error("NewLoggedSet with multiple errors returned valid result")
	}
	if ls != nil {
		t.Error("NewLoggedSet with invalid input returned non-nil logged set")
	}
	if len(result.Errors) != 8 {
		t.Errorf("Expected 8 errors, got %d", len(result.Errors))
	}
}

// ==================== LoggedSet.Validate Tests ====================

func TestLoggedSet_Validate_Valid(t *testing.T) {
	ls := &LoggedSet{
		ID:             "test-id",
		UserID:         "user-123",
		SessionID:      "session-456",
		PrescriptionID: "prescription-789",
		LiftID:         "lift-abc",
		SetNumber:      1,
		Weight:         225.0,
		TargetReps:     5,
		RepsPerformed:  5,
	}

	result := ls.Validate()

	if !result.Valid {
		t.Errorf("Validate returned invalid result for valid logged set: %v", result.Errors)
	}
}

func TestLoggedSet_Validate_Invalid(t *testing.T) {
	ls := &LoggedSet{
		ID:             "test-id",
		UserID:         "",
		SessionID:      "session-456",
		PrescriptionID: "prescription-789",
		LiftID:         "lift-abc",
		SetNumber:      1,
		Weight:         225.0,
		TargetReps:     5,
		RepsPerformed:  5,
	}

	result := ls.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for logged set with empty user ID")
	}
}

// ==================== Helper Method Tests ====================

func TestLoggedSet_ExceededTarget(t *testing.T) {
	tests := []struct {
		name          string
		targetReps    int
		repsPerformed int
		expected      bool
	}{
		{"exceeded", 5, 7, true},
		{"met exactly", 5, 5, false},
		{"fell short", 5, 3, false},
		{"zero performed", 5, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := &LoggedSet{
				TargetReps:    tt.targetReps,
				RepsPerformed: tt.repsPerformed,
			}
			if got := ls.ExceededTarget(); got != tt.expected {
				t.Errorf("ExceededTarget() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLoggedSet_RepsDifference(t *testing.T) {
	tests := []struct {
		name          string
		targetReps    int
		repsPerformed int
		expected      int
	}{
		{"exceeded by 2", 5, 7, 2},
		{"met exactly", 5, 5, 0},
		{"fell short by 2", 5, 3, -2},
		{"zero performed", 5, 0, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := &LoggedSet{
				TargetReps:    tt.targetReps,
				RepsPerformed: tt.repsPerformed,
			}
			if got := ls.RepsDifference(); got != tt.expected {
				t.Errorf("RepsDifference() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// ==================== ValidationResult Tests ====================

func TestValidationResult_AddError(t *testing.T) {
	result := NewValidationResult()

	if !result.Valid {
		t.Error("NewValidationResult() should be valid")
	}

	result.AddError(ErrUserIDRequired)

	if result.Valid {
		t.Error("ValidationResult should be invalid after AddError")
	}
	if len(result.Errors) != 1 {
		t.Errorf("len(result.Errors) = %d, want 1", len(result.Errors))
	}
}

package liftmax

import (
	"errors"
	"math"
	"strings"
	"testing"
	"time"
)

// mockRepository implements LiftMaxRepository for testing
type mockRepository struct {
	oneRMs map[string]*LiftMax // key: userID:liftID
	err    error               // simulated error
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		oneRMs: make(map[string]*LiftMax),
	}
}

func (m *mockRepository) GetCurrentOneRM(userID, liftID string) (*LiftMax, error) {
	if m.err != nil {
		return nil, m.err
	}
	key := userID + ":" + liftID
	return m.oneRMs[key], nil
}

func (m *mockRepository) SetOneRM(userID, liftID string, oneRM *LiftMax) {
	key := userID + ":" + liftID
	m.oneRMs[key] = oneRM
}

func (m *mockRepository) SetError(err error) {
	m.err = err
}

// ==================== Value Validation Tests ====================

func TestValidateValue_Valid(t *testing.T) {
	tests := []struct {
		name  string
		value float64
	}{
		{"whole number", 315},
		{"quarter", 315.25},
		{"half", 315.5},
		{"three quarters", 315.75},
		{"minimum positive", 0.25},
		{"large value", 1000.5},
		{"small value", 2.5},
		{"zero quarters", 100.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateValue(tt.value)
			if err != nil {
				t.Errorf("ValidateValue(%v) = %v, want nil", tt.value, err)
			}
		})
	}
}

func TestValidateValue_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		value       float64
		expectedErr error
	}{
		{"zero", 0, ErrValueNotPositive},
		{"negative", -100, ErrValueNotPositive},
		{"invalid precision 0.1", 315.1, ErrValueInvalidPrecision},
		{"invalid precision 0.3", 315.3, ErrValueInvalidPrecision},
		{"invalid precision 0.33", 315.33, ErrValueInvalidPrecision},
		{"invalid precision 0.125", 315.125, ErrValueInvalidPrecision},
		{"invalid precision 0.01", 315.01, ErrValueInvalidPrecision},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateValue(tt.value)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateValue(%v) = %v, want %v", tt.value, err, tt.expectedErr)
			}
		})
	}
}

// ==================== Type Validation Tests ====================

func TestValidateType_Valid(t *testing.T) {
	tests := []struct {
		name    string
		maxType MaxType
	}{
		{"ONE_RM", OneRM},
		{"TRAINING_MAX", TrainingMax},
		{"E1RM", E1RM},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateType(tt.maxType)
			if err != nil {
				t.Errorf("ValidateType(%q) = %v, want nil", tt.maxType, err)
			}
		})
	}
}

func TestValidateType_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		maxType     MaxType
		expectedErr error
	}{
		{"empty", "", ErrTypeRequired},
		{"invalid type", MaxType("INVALID"), ErrTypeInvalid},
		{"lowercase one_rm", MaxType("one_rm"), ErrTypeInvalid},
		{"partial match", MaxType("ONE"), ErrTypeInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateType(tt.maxType)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateType(%q) = %v, want %v", tt.maxType, err, tt.expectedErr)
			}
		})
	}
}

// ==================== UserID Validation Tests ====================

func TestValidateUserID_Valid(t *testing.T) {
	tests := []struct {
		name   string
		userID string
	}{
		{"normal uuid", "user-123-456"},
		{"simple", "user1"},
		{"with dashes", "550e8400-e29b-41d4-a716-446655440000"},
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
		name   string
		userID string
	}{
		{"empty", ""},
		{"only spaces", "   "},
		{"only tabs", "\t\t"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserID(tt.userID)
			if !errors.Is(err, ErrUserIDRequired) {
				t.Errorf("ValidateUserID(%q) = %v, want %v", tt.userID, err, ErrUserIDRequired)
			}
		})
	}
}

// ==================== LiftID Validation Tests ====================

func TestValidateLiftID_Valid(t *testing.T) {
	tests := []struct {
		name   string
		liftID string
	}{
		{"normal uuid", "lift-123-456"},
		{"simple", "squat"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLiftID(tt.liftID)
			if err != nil {
				t.Errorf("ValidateLiftID(%q) = %v, want nil", tt.liftID, err)
			}
		})
	}
}

func TestValidateLiftID_Invalid(t *testing.T) {
	tests := []struct {
		name   string
		liftID string
	}{
		{"empty", ""},
		{"only spaces", "   "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLiftID(tt.liftID)
			if !errors.Is(err, ErrLiftIDRequired) {
				t.Errorf("ValidateLiftID(%q) = %v, want %v", tt.liftID, err, ErrLiftIDRequired)
			}
		})
	}
}

// ==================== TM Validation Warning Tests ====================

func TestValidateTMAgainstOneRM_NoWarning(t *testing.T) {
	oneRM := &LiftMax{Value: 400}

	tests := []struct {
		name    string
		tmValue float64
	}{
		{"exactly 80%", 320},
		{"exactly 90%", 360},
		{"exactly 95%", 380},
		{"85%", 340},
		{"92%", 368},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warning := ValidateTMAgainstOneRM(tt.tmValue, oneRM)
			if warning != "" {
				t.Errorf("ValidateTMAgainstOneRM(%v, oneRM=%v) = %q, want empty", tt.tmValue, oneRM.Value, warning)
			}
		})
	}
}

func TestValidateTMAgainstOneRM_WarningBelowRange(t *testing.T) {
	oneRM := &LiftMax{Value: 400}

	tests := []struct {
		name    string
		tmValue float64
	}{
		{"79%", 316},
		{"75%", 300},
		{"50%", 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warning := ValidateTMAgainstOneRM(tt.tmValue, oneRM)
			if warning == "" {
				t.Errorf("ValidateTMAgainstOneRM(%v, oneRM=%v) returned no warning, expected below range warning", tt.tmValue, oneRM.Value)
			}
			if !strings.Contains(warning, "below") {
				t.Errorf("Warning should mention 'below': %s", warning)
			}
		})
	}
}

func TestValidateTMAgainstOneRM_WarningAboveRange(t *testing.T) {
	oneRM := &LiftMax{Value: 400}

	tests := []struct {
		name    string
		tmValue float64
	}{
		{"96%", 384},
		{"100%", 400},
		{"110%", 440},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warning := ValidateTMAgainstOneRM(tt.tmValue, oneRM)
			if warning == "" {
				t.Errorf("ValidateTMAgainstOneRM(%v, oneRM=%v) returned no warning, expected above range warning", tt.tmValue, oneRM.Value)
			}
			if !strings.Contains(warning, "above") {
				t.Errorf("Warning should mention 'above': %s", warning)
			}
		})
	}
}

func TestValidateTMAgainstOneRM_NilOneRM(t *testing.T) {
	warning := ValidateTMAgainstOneRM(350, nil)
	if warning != "" {
		t.Errorf("ValidateTMAgainstOneRM with nil oneRM = %q, want empty", warning)
	}
}

// ==================== CreateLiftMax Tests ====================

func TestCreateLiftMax_ValidInput(t *testing.T) {
	repo := newMockRepository()

	input := CreateLiftMaxInput{
		UserID: "user-123",
		LiftID: "lift-456",
		Type:   OneRM,
		Value:  315.5,
	}

	liftMax, result := CreateLiftMax(input, "max-id", repo)

	if !result.Valid {
		t.Errorf("CreateLiftMax returned invalid result: %v", result.Errors)
	}
	if liftMax == nil {
		t.Fatal("CreateLiftMax returned nil liftMax")
	}
	if liftMax.ID != "max-id" {
		t.Errorf("liftMax.ID = %q, want %q", liftMax.ID, "max-id")
	}
	if liftMax.UserID != "user-123" {
		t.Errorf("liftMax.UserID = %q, want %q", liftMax.UserID, "user-123")
	}
	if liftMax.LiftID != "lift-456" {
		t.Errorf("liftMax.LiftID = %q, want %q", liftMax.LiftID, "lift-456")
	}
	if liftMax.Type != OneRM {
		t.Errorf("liftMax.Type = %q, want %q", liftMax.Type, OneRM)
	}
	if liftMax.Value != 315.5 {
		t.Errorf("liftMax.Value = %v, want %v", liftMax.Value, 315.5)
	}
}

func TestCreateLiftMax_WithEffectiveDate(t *testing.T) {
	repo := newMockRepository()
	effectiveDate := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	input := CreateLiftMaxInput{
		UserID:        "user-123",
		LiftID:        "lift-456",
		Type:          OneRM,
		Value:         315,
		EffectiveDate: &effectiveDate,
	}

	liftMax, result := CreateLiftMax(input, "max-id", repo)

	if !result.Valid {
		t.Errorf("CreateLiftMax returned invalid result: %v", result.Errors)
	}
	if !liftMax.EffectiveDate.Equal(effectiveDate) {
		t.Errorf("liftMax.EffectiveDate = %v, want %v", liftMax.EffectiveDate, effectiveDate)
	}
}

func TestCreateLiftMax_DefaultEffectiveDate(t *testing.T) {
	repo := newMockRepository()
	before := time.Now()

	input := CreateLiftMaxInput{
		UserID: "user-123",
		LiftID: "lift-456",
		Type:   OneRM,
		Value:  315,
		// No EffectiveDate provided
	}

	liftMax, result := CreateLiftMax(input, "max-id", repo)
	after := time.Now()

	if !result.Valid {
		t.Errorf("CreateLiftMax returned invalid result: %v", result.Errors)
	}
	if liftMax.EffectiveDate.Before(before) || liftMax.EffectiveDate.After(after) {
		t.Errorf("liftMax.EffectiveDate should be between %v and %v, got %v", before, after, liftMax.EffectiveDate)
	}
}

func TestCreateLiftMax_E1RM(t *testing.T) {
	repo := newMockRepository()

	input := CreateLiftMaxInput{
		UserID: "user-123",
		LiftID: "lift-456",
		Type:   E1RM,
		Value:  320.5,
	}

	liftMax, result := CreateLiftMax(input, "max-id", repo)

	if !result.Valid {
		t.Errorf("CreateLiftMax returned invalid result: %v", result.Errors)
	}
	if liftMax == nil {
		t.Fatal("CreateLiftMax returned nil liftMax")
	}
	if liftMax.Type != E1RM {
		t.Errorf("liftMax.Type = %q, want %q", liftMax.Type, E1RM)
	}
	if liftMax.Value != 320.5 {
		t.Errorf("liftMax.Value = %v, want %v", liftMax.Value, 320.5)
	}
	// E1RM should not trigger TM-related warnings
	if result.HasWarnings() {
		t.Errorf("CreateLiftMax with E1RM should not have warnings: %v", result.Warnings)
	}
}

func TestCreateLiftMax_TrainingMaxWithWarning(t *testing.T) {
	repo := newMockRepository()
	// Set up existing 1RM of 400
	repo.SetOneRM("user-123", "lift-456", &LiftMax{
		Value: 400,
		Type:  OneRM,
	})

	input := CreateLiftMaxInput{
		UserID: "user-123",
		LiftID: "lift-456",
		Type:   TrainingMax,
		Value:  300, // 75% of 400 - below 80% threshold
	}

	liftMax, result := CreateLiftMax(input, "max-id", repo)

	if !result.Valid {
		t.Errorf("CreateLiftMax should be valid even with warning: %v", result.Errors)
	}
	if liftMax == nil {
		t.Fatal("CreateLiftMax should return liftMax even with warning")
	}
	if !result.HasWarnings() {
		t.Error("CreateLiftMax should have warning for TM below 80%")
	}
}

func TestCreateLiftMax_TrainingMaxNoWarning(t *testing.T) {
	repo := newMockRepository()
	repo.SetOneRM("user-123", "lift-456", &LiftMax{
		Value: 400,
		Type:  OneRM,
	})

	input := CreateLiftMaxInput{
		UserID: "user-123",
		LiftID: "lift-456",
		Type:   TrainingMax,
		Value:  360, // 90% of 400 - within range
	}

	liftMax, result := CreateLiftMax(input, "max-id", repo)

	if !result.Valid {
		t.Errorf("CreateLiftMax returned invalid result: %v", result.Errors)
	}
	if liftMax == nil {
		t.Fatal("CreateLiftMax returned nil liftMax")
	}
	if result.HasWarnings() {
		t.Errorf("CreateLiftMax should not have warnings: %v", result.Warnings)
	}
}

func TestCreateLiftMax_TrainingMaxNoExistingOneRM(t *testing.T) {
	repo := newMockRepository()
	// No 1RM exists

	input := CreateLiftMaxInput{
		UserID: "user-123",
		LiftID: "lift-456",
		Type:   TrainingMax,
		Value:  300,
	}

	liftMax, result := CreateLiftMax(input, "max-id", repo)

	if !result.Valid {
		t.Errorf("CreateLiftMax returned invalid result: %v", result.Errors)
	}
	if liftMax == nil {
		t.Fatal("CreateLiftMax returned nil liftMax")
	}
	// No warning since no 1RM to compare against
	if result.HasWarnings() {
		t.Errorf("CreateLiftMax should not have warnings when no 1RM exists: %v", result.Warnings)
	}
}

func TestCreateLiftMax_TrainingMaxRepoError(t *testing.T) {
	repo := newMockRepository()
	repo.SetError(errors.New("database error"))

	input := CreateLiftMaxInput{
		UserID: "user-123",
		LiftID: "lift-456",
		Type:   TrainingMax,
		Value:  300,
	}

	liftMax, result := CreateLiftMax(input, "max-id", repo)

	if !result.Valid {
		t.Errorf("CreateLiftMax should be valid even with repo error: %v", result.Errors)
	}
	if liftMax == nil {
		t.Fatal("CreateLiftMax should return liftMax even with repo error")
	}
	// Should have warning about unable to validate
	if !result.HasWarnings() {
		t.Error("CreateLiftMax should have warning when repo returns error")
	}
}

func TestCreateLiftMax_InvalidUserID(t *testing.T) {
	repo := newMockRepository()

	input := CreateLiftMaxInput{
		UserID: "",
		LiftID: "lift-456",
		Type:   OneRM,
		Value:  315,
	}

	liftMax, result := CreateLiftMax(input, "max-id", repo)

	if result.Valid {
		t.Error("CreateLiftMax with empty userID should be invalid")
	}
	if liftMax != nil {
		t.Error("CreateLiftMax with invalid input should return nil")
	}
}

func TestCreateLiftMax_InvalidLiftID(t *testing.T) {
	repo := newMockRepository()

	input := CreateLiftMaxInput{
		UserID: "user-123",
		LiftID: "",
		Type:   OneRM,
		Value:  315,
	}

	liftMax, result := CreateLiftMax(input, "max-id", repo)

	if result.Valid {
		t.Error("CreateLiftMax with empty liftID should be invalid")
	}
	if liftMax != nil {
		t.Error("CreateLiftMax with invalid input should return nil")
	}
}

func TestCreateLiftMax_InvalidType(t *testing.T) {
	repo := newMockRepository()

	input := CreateLiftMaxInput{
		UserID: "user-123",
		LiftID: "lift-456",
		Type:   "",
		Value:  315,
	}

	liftMax, result := CreateLiftMax(input, "max-id", repo)

	if result.Valid {
		t.Error("CreateLiftMax with empty type should be invalid")
	}
	if liftMax != nil {
		t.Error("CreateLiftMax with invalid input should return nil")
	}
}

func TestCreateLiftMax_InvalidValue(t *testing.T) {
	repo := newMockRepository()

	input := CreateLiftMaxInput{
		UserID: "user-123",
		LiftID: "lift-456",
		Type:   OneRM,
		Value:  0,
	}

	liftMax, result := CreateLiftMax(input, "max-id", repo)

	if result.Valid {
		t.Error("CreateLiftMax with zero value should be invalid")
	}
	if liftMax != nil {
		t.Error("CreateLiftMax with invalid input should return nil")
	}
}

func TestCreateLiftMax_MultipleErrors(t *testing.T) {
	repo := newMockRepository()

	input := CreateLiftMaxInput{
		UserID: "",      // Invalid
		LiftID: "",      // Invalid
		Type:   "",      // Invalid
		Value:  -100,    // Invalid
	}

	liftMax, result := CreateLiftMax(input, "max-id", repo)

	if result.Valid {
		t.Error("CreateLiftMax with multiple invalid fields should be invalid")
	}
	if liftMax != nil {
		t.Error("CreateLiftMax with invalid input should return nil")
	}
	if len(result.Errors) < 4 {
		t.Errorf("Expected at least 4 errors, got %d", len(result.Errors))
	}
}

func TestCreateLiftMax_NilRepo(t *testing.T) {
	input := CreateLiftMaxInput{
		UserID: "user-123",
		LiftID: "lift-456",
		Type:   TrainingMax,
		Value:  315,
	}

	liftMax, result := CreateLiftMax(input, "max-id", nil)

	if !result.Valid {
		t.Errorf("CreateLiftMax with nil repo should be valid: %v", result.Errors)
	}
	if liftMax == nil {
		t.Fatal("CreateLiftMax with nil repo should return liftMax")
	}
}

// ==================== UpdateLiftMax Tests ====================

func TestUpdateLiftMax_UpdateValue(t *testing.T) {
	repo := newMockRepository()
	liftMax := &LiftMax{
		ID:     "max-id",
		UserID: "user-123",
		LiftID: "lift-456",
		Type:   OneRM,
		Value:  315,
	}

	newValue := 325.5
	input := UpdateLiftMaxInput{Value: &newValue}

	result := UpdateLiftMax(liftMax, input, repo)

	if !result.Valid {
		t.Errorf("UpdateLiftMax returned invalid result: %v", result.Errors)
	}
	if liftMax.Value != 325.5 {
		t.Errorf("liftMax.Value = %v, want %v", liftMax.Value, 325.5)
	}
}

func TestUpdateLiftMax_UpdateType(t *testing.T) {
	repo := newMockRepository()
	liftMax := &LiftMax{
		ID:     "max-id",
		UserID: "user-123",
		LiftID: "lift-456",
		Type:   OneRM,
		Value:  315,
	}

	newType := TrainingMax
	input := UpdateLiftMaxInput{Type: &newType}

	result := UpdateLiftMax(liftMax, input, repo)

	if !result.Valid {
		t.Errorf("UpdateLiftMax returned invalid result: %v", result.Errors)
	}
	if liftMax.Type != TrainingMax {
		t.Errorf("liftMax.Type = %v, want %v", liftMax.Type, TrainingMax)
	}
}

func TestUpdateLiftMax_UpdateEffectiveDate(t *testing.T) {
	repo := newMockRepository()
	originalDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	liftMax := &LiftMax{
		ID:            "max-id",
		UserID:        "user-123",
		LiftID:        "lift-456",
		Type:          OneRM,
		Value:         315,
		EffectiveDate: originalDate,
	}

	newDate := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	input := UpdateLiftMaxInput{EffectiveDate: &newDate}

	result := UpdateLiftMax(liftMax, input, repo)

	if !result.Valid {
		t.Errorf("UpdateLiftMax returned invalid result: %v", result.Errors)
	}
	if !liftMax.EffectiveDate.Equal(newDate) {
		t.Errorf("liftMax.EffectiveDate = %v, want %v", liftMax.EffectiveDate, newDate)
	}
}

func TestUpdateLiftMax_UpdateMultipleFields(t *testing.T) {
	repo := newMockRepository()
	liftMax := &LiftMax{
		ID:     "max-id",
		UserID: "user-123",
		LiftID: "lift-456",
		Type:   OneRM,
		Value:  315,
	}

	newType := TrainingMax
	newValue := 285.0
	input := UpdateLiftMaxInput{
		Type:  &newType,
		Value: &newValue,
	}

	result := UpdateLiftMax(liftMax, input, repo)

	if !result.Valid {
		t.Errorf("UpdateLiftMax returned invalid result: %v", result.Errors)
	}
	if liftMax.Type != TrainingMax {
		t.Errorf("liftMax.Type = %v, want %v", liftMax.Type, TrainingMax)
	}
	if liftMax.Value != 285.0 {
		t.Errorf("liftMax.Value = %v, want %v", liftMax.Value, 285.0)
	}
}

func TestUpdateLiftMax_InvalidValue(t *testing.T) {
	repo := newMockRepository()
	liftMax := &LiftMax{
		ID:     "max-id",
		UserID: "user-123",
		LiftID: "lift-456",
		Type:   OneRM,
		Value:  315,
	}
	originalValue := liftMax.Value

	invalidValue := -100.0
	input := UpdateLiftMaxInput{Value: &invalidValue}

	result := UpdateLiftMax(liftMax, input, repo)

	if result.Valid {
		t.Error("UpdateLiftMax with invalid value should be invalid")
	}
	if liftMax.Value != originalValue {
		t.Errorf("liftMax.Value should not change on validation failure")
	}
}

func TestUpdateLiftMax_InvalidType(t *testing.T) {
	repo := newMockRepository()
	liftMax := &LiftMax{
		ID:     "max-id",
		UserID: "user-123",
		LiftID: "lift-456",
		Type:   OneRM,
		Value:  315,
	}
	originalType := liftMax.Type

	invalidType := MaxType("INVALID")
	input := UpdateLiftMaxInput{Type: &invalidType}

	result := UpdateLiftMax(liftMax, input, repo)

	if result.Valid {
		t.Error("UpdateLiftMax with invalid type should be invalid")
	}
	if liftMax.Type != originalType {
		t.Errorf("liftMax.Type should not change on validation failure")
	}
}

func TestUpdateLiftMax_NoChanges(t *testing.T) {
	repo := newMockRepository()
	liftMax := &LiftMax{
		ID:     "max-id",
		UserID: "user-123",
		LiftID: "lift-456",
		Type:   OneRM,
		Value:  315,
	}

	input := UpdateLiftMaxInput{} // No changes

	result := UpdateLiftMax(liftMax, input, repo)

	if !result.Valid {
		t.Errorf("UpdateLiftMax with no changes should be valid: %v", result.Errors)
	}
}

func TestUpdateLiftMax_TrainingMaxWithWarning(t *testing.T) {
	repo := newMockRepository()
	repo.SetOneRM("user-123", "lift-456", &LiftMax{
		Value: 400,
		Type:  OneRM,
	})

	liftMax := &LiftMax{
		ID:     "max-id",
		UserID: "user-123",
		LiftID: "lift-456",
		Type:   TrainingMax,
		Value:  360,
	}

	newValue := 300.0 // 75% of 400 - below threshold
	input := UpdateLiftMaxInput{Value: &newValue}

	result := UpdateLiftMax(liftMax, input, repo)

	if !result.Valid {
		t.Errorf("UpdateLiftMax should be valid even with warning: %v", result.Errors)
	}
	if !result.HasWarnings() {
		t.Error("UpdateLiftMax should have warning for TM below 80%")
	}
	if liftMax.Value != 300.0 {
		t.Errorf("liftMax.Value should be updated: got %v, want %v", liftMax.Value, 300.0)
	}
}

func TestUpdateLiftMax_ChangeToTrainingMaxWithWarning(t *testing.T) {
	repo := newMockRepository()
	repo.SetOneRM("user-123", "lift-456", &LiftMax{
		Value: 400,
		Type:  OneRM,
	})

	liftMax := &LiftMax{
		ID:     "max-id",
		UserID: "user-123",
		LiftID: "lift-456",
		Type:   OneRM,
		Value:  300, // When changed to TM, this is 75% of 400
	}

	newType := TrainingMax
	input := UpdateLiftMaxInput{Type: &newType}

	result := UpdateLiftMax(liftMax, input, repo)

	if !result.Valid {
		t.Errorf("UpdateLiftMax should be valid even with warning: %v", result.Errors)
	}
	if !result.HasWarnings() {
		t.Error("UpdateLiftMax should have warning when changing to TM below 80%")
	}
}

func TestUpdateLiftMax_UpdatesTimestamp(t *testing.T) {
	repo := newMockRepository()
	originalUpdatedAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	liftMax := &LiftMax{
		ID:        "max-id",
		UserID:    "user-123",
		LiftID:    "lift-456",
		Type:      OneRM,
		Value:     315,
		UpdatedAt: originalUpdatedAt,
	}

	newValue := 325.0
	input := UpdateLiftMaxInput{Value: &newValue}

	before := time.Now()
	result := UpdateLiftMax(liftMax, input, repo)
	after := time.Now()

	if !result.Valid {
		t.Errorf("UpdateLiftMax returned invalid result: %v", result.Errors)
	}
	if liftMax.UpdatedAt.Before(before) || liftMax.UpdatedAt.After(after) {
		t.Errorf("liftMax.UpdatedAt should be between %v and %v, got %v", before, after, liftMax.UpdatedAt)
	}
}

// ==================== LiftMax.Validate Tests ====================

func TestLiftMax_Validate_Valid(t *testing.T) {
	liftMax := &LiftMax{
		ID:     "max-id",
		UserID: "user-123",
		LiftID: "lift-456",
		Type:   OneRM,
		Value:  315,
	}

	result := liftMax.Validate()

	if !result.Valid {
		t.Errorf("Validate returned invalid for valid liftMax: %v", result.Errors)
	}
}

func TestLiftMax_Validate_InvalidUserID(t *testing.T) {
	liftMax := &LiftMax{
		ID:     "max-id",
		UserID: "",
		LiftID: "lift-456",
		Type:   OneRM,
		Value:  315,
	}

	result := liftMax.Validate()

	if result.Valid {
		t.Error("Validate should be invalid for empty userID")
	}
}

func TestLiftMax_Validate_InvalidLiftID(t *testing.T) {
	liftMax := &LiftMax{
		ID:     "max-id",
		UserID: "user-123",
		LiftID: "",
		Type:   OneRM,
		Value:  315,
	}

	result := liftMax.Validate()

	if result.Valid {
		t.Error("Validate should be invalid for empty liftID")
	}
}

func TestLiftMax_Validate_InvalidType(t *testing.T) {
	liftMax := &LiftMax{
		ID:     "max-id",
		UserID: "user-123",
		LiftID: "lift-456",
		Type:   "",
		Value:  315,
	}

	result := liftMax.Validate()

	if result.Valid {
		t.Error("Validate should be invalid for empty type")
	}
}

func TestLiftMax_Validate_InvalidValue(t *testing.T) {
	liftMax := &LiftMax{
		ID:     "max-id",
		UserID: "user-123",
		LiftID: "lift-456",
		Type:   OneRM,
		Value:  0,
	}

	result := liftMax.Validate()

	if result.Valid {
		t.Error("Validate should be invalid for zero value")
	}
}

// ==================== ValidationResult Tests ====================

func TestValidationResult_NewValidationResult(t *testing.T) {
	result := NewValidationResult()

	if !result.Valid {
		t.Error("NewValidationResult should be valid")
	}
	if len(result.Errors) != 0 {
		t.Error("NewValidationResult should have no errors")
	}
	if len(result.Warnings) != 0 {
		t.Error("NewValidationResult should have no warnings")
	}
}

func TestValidationResult_AddError(t *testing.T) {
	result := NewValidationResult()
	result.AddError(ErrValueNotPositive)

	if result.Valid {
		t.Error("AddError should mark result as invalid")
	}
	if len(result.Errors) != 1 {
		t.Errorf("AddError should add one error, got %d", len(result.Errors))
	}
}

func TestValidationResult_AddWarning(t *testing.T) {
	result := NewValidationResult()
	result.AddWarning("test warning")

	if !result.Valid {
		t.Error("AddWarning should not mark result as invalid")
	}
	if !result.HasWarnings() {
		t.Error("HasWarnings should return true after AddWarning")
	}
	if len(result.Warnings) != 1 {
		t.Errorf("AddWarning should add one warning, got %d", len(result.Warnings))
	}
}

func TestValidationResult_Error_Valid(t *testing.T) {
	result := NewValidationResult()

	err := result.Error()
	if err != nil {
		t.Errorf("Error() on valid result should return nil, got %v", err)
	}
}

func TestValidationResult_Error_Invalid(t *testing.T) {
	result := NewValidationResult()
	result.AddError(ErrValueNotPositive)
	result.AddError(ErrTypeRequired)

	err := result.Error()
	if err == nil {
		t.Error("Error() on invalid result should return error")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "positive") {
		t.Errorf("Error message should contain value error: %s", errMsg)
	}
	if !strings.Contains(errMsg, "required") {
		t.Errorf("Error message should contain type error: %s", errMsg)
	}
}

// ==================== MaxCalculator Tests ====================

func TestMaxCalculator_ConvertToTM_DefaultPercentage(t *testing.T) {
	calc := NewMaxCalculator()

	tm, err := calc.ConvertToTM(400, nil)

	if err != nil {
		t.Errorf("ConvertToTM returned error: %v", err)
	}
	if tm != 360 {
		t.Errorf("ConvertToTM(400, nil) = %v, want %v (90%% of 400)", tm, 360)
	}
}

func TestMaxCalculator_ConvertToTM_CustomPercentage(t *testing.T) {
	calc := NewMaxCalculator()

	tests := []struct {
		name       string
		oneRM      float64
		percentage float64
		expected   float64
	}{
		{"85%", 400, 85, 340},
		{"90%", 400, 90, 360},
		{"95%", 400, 95, 380},
		{"80%", 315, 80, 252},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pct := tt.percentage
			tm, err := calc.ConvertToTM(tt.oneRM, &pct)
			if err != nil {
				t.Errorf("ConvertToTM returned error: %v", err)
			}
			if tm != tt.expected {
				t.Errorf("ConvertToTM(%v, %v) = %v, want %v", tt.oneRM, tt.percentage, tm, tt.expected)
			}
		})
	}
}

func TestMaxCalculator_ConvertToTM_RoundsToQuarter(t *testing.T) {
	calc := NewMaxCalculator()

	// 85% of 315 = 267.75
	pct := 85.0
	tm, err := calc.ConvertToTM(315, &pct)

	if err != nil {
		t.Errorf("ConvertToTM returned error: %v", err)
	}
	if tm != 267.75 {
		t.Errorf("ConvertToTM(315, 85) = %v, want %v", tm, 267.75)
	}
}

func TestMaxCalculator_ConvertToTM_InvalidPercentage(t *testing.T) {
	calc := NewMaxCalculator()

	tests := []struct {
		name       string
		percentage float64
	}{
		{"zero", 0},
		{"negative", -10},
		{"over 100", 101},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pct := tt.percentage
			_, err := calc.ConvertToTM(400, &pct)
			if !errors.Is(err, ErrConversionPercentageInvalid) {
				t.Errorf("ConvertToTM with %v%% should return ErrConversionPercentageInvalid, got %v", tt.percentage, err)
			}
		})
	}
}

func TestMaxCalculator_ConvertToOneRM_DefaultPercentage(t *testing.T) {
	calc := NewMaxCalculator()

	oneRM, err := calc.ConvertToOneRM(360, nil)

	if err != nil {
		t.Errorf("ConvertToOneRM returned error: %v", err)
	}
	if oneRM != 400 {
		t.Errorf("ConvertToOneRM(360, nil) = %v, want %v", oneRM, 400)
	}
}

func TestMaxCalculator_ConvertToOneRM_CustomPercentage(t *testing.T) {
	calc := NewMaxCalculator()

	tests := []struct {
		name       string
		tm         float64
		percentage float64
		expected   float64
	}{
		{"from 85%", 340, 85, 400},
		{"from 90%", 360, 90, 400},
		{"from 80%", 252, 80, 315},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pct := tt.percentage
			oneRM, err := calc.ConvertToOneRM(tt.tm, &pct)
			if err != nil {
				t.Errorf("ConvertToOneRM returned error: %v", err)
			}
			if oneRM != tt.expected {
				t.Errorf("ConvertToOneRM(%v, %v) = %v, want %v", tt.tm, tt.percentage, oneRM, tt.expected)
			}
		})
	}
}

func TestMaxCalculator_ConvertToOneRM_RoundsToQuarter(t *testing.T) {
	calc := NewMaxCalculator()

	// 267 / 0.85 = 314.117... should round to 314.25
	pct := 85.0
	oneRM, err := calc.ConvertToOneRM(267, &pct)

	if err != nil {
		t.Errorf("ConvertToOneRM returned error: %v", err)
	}
	// 267 / 0.85 = 314.117647... rounds to 314.0
	expected := RoundToQuarter(267 / 0.85)
	if oneRM != expected {
		t.Errorf("ConvertToOneRM(267, 85) = %v, want %v", oneRM, expected)
	}
}

func TestMaxCalculator_ConvertToOneRM_InvalidPercentage(t *testing.T) {
	calc := NewMaxCalculator()

	tests := []struct {
		name       string
		percentage float64
	}{
		{"zero", 0},
		{"negative", -10},
		{"over 100", 101},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pct := tt.percentage
			_, err := calc.ConvertToOneRM(360, &pct)
			if !errors.Is(err, ErrConversionPercentageInvalid) {
				t.Errorf("ConvertToOneRM with %v%% should return ErrConversionPercentageInvalid, got %v", tt.percentage, err)
			}
		})
	}
}

func TestMaxCalculator_Convert_OneRMToTM(t *testing.T) {
	calc := NewMaxCalculator()

	result, err := calc.Convert(400, OneRM, nil)

	if err != nil {
		t.Errorf("Convert returned error: %v", err)
	}
	if result.OriginalValue != 400 {
		t.Errorf("result.OriginalValue = %v, want %v", result.OriginalValue, 400)
	}
	if result.OriginalType != OneRM {
		t.Errorf("result.OriginalType = %v, want %v", result.OriginalType, OneRM)
	}
	if result.ConvertedValue != 360 {
		t.Errorf("result.ConvertedValue = %v, want %v", result.ConvertedValue, 360)
	}
	if result.ConvertedType != TrainingMax {
		t.Errorf("result.ConvertedType = %v, want %v", result.ConvertedType, TrainingMax)
	}
	if result.Percentage != 90 {
		t.Errorf("result.Percentage = %v, want %v", result.Percentage, 90)
	}
}

func TestMaxCalculator_Convert_TMToOneRM(t *testing.T) {
	calc := NewMaxCalculator()

	result, err := calc.Convert(360, TrainingMax, nil)

	if err != nil {
		t.Errorf("Convert returned error: %v", err)
	}
	if result.OriginalValue != 360 {
		t.Errorf("result.OriginalValue = %v, want %v", result.OriginalValue, 360)
	}
	if result.OriginalType != TrainingMax {
		t.Errorf("result.OriginalType = %v, want %v", result.OriginalType, TrainingMax)
	}
	if result.ConvertedValue != 400 {
		t.Errorf("result.ConvertedValue = %v, want %v", result.ConvertedValue, 400)
	}
	if result.ConvertedType != OneRM {
		t.Errorf("result.ConvertedType = %v, want %v", result.ConvertedType, OneRM)
	}
}

func TestMaxCalculator_Convert_CustomPercentage(t *testing.T) {
	calc := NewMaxCalculator()

	pct := 85.0
	result, err := calc.Convert(400, OneRM, &pct)

	if err != nil {
		t.Errorf("Convert returned error: %v", err)
	}
	if result.ConvertedValue != 340 {
		t.Errorf("result.ConvertedValue = %v, want %v (85%% of 400)", result.ConvertedValue, 340)
	}
	if result.Percentage != 85 {
		t.Errorf("result.Percentage = %v, want %v", result.Percentage, 85)
	}
}

func TestMaxCalculator_Convert_InvalidType(t *testing.T) {
	calc := NewMaxCalculator()

	_, err := calc.Convert(400, MaxType("INVALID"), nil)

	if !errors.Is(err, ErrTypeInvalid) {
		t.Errorf("Convert with invalid type should return ErrTypeInvalid, got %v", err)
	}
}

func TestMaxCalculator_Convert_InvalidPercentage(t *testing.T) {
	calc := NewMaxCalculator()

	pct := 0.0
	_, err := calc.Convert(400, OneRM, &pct)

	if !errors.Is(err, ErrConversionPercentageInvalid) {
		t.Errorf("Convert with invalid percentage should return ErrConversionPercentageInvalid, got %v", err)
	}
}

// ==================== RoundToQuarter Tests ====================

func TestRoundToQuarter(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected float64
	}{
		{"exact quarter", 315.25, 315.25},
		{"exact half", 315.5, 315.5},
		{"exact three quarters", 315.75, 315.75},
		{"whole number", 315.0, 315.0},
		{"rounds down to quarter", 315.1, 315.0},
		{"rounds up to quarter", 315.2, 315.25},
		{"rounds down to half", 315.4, 315.5},
		{"rounds up to half", 315.45, 315.5},
		{"rounds down to three quarters", 315.65, 315.75},
		{"rounds up to whole", 315.9, 316.0},
		{"rounds 0.125 to 0.25", 315.125, 315.25},
		{"rounds 0.375 to 0.5", 315.375, 315.5},
		{"rounds 0.625 to 0.75", 315.625, 315.75},
		{"rounds 0.875 to 1.0", 315.875, 316.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RoundToQuarter(tt.value)
			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("RoundToQuarter(%v) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}

// ==================== Edge Cases ====================

func TestValidateValue_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		valid bool
	}{
		{"minimum valid", 0.25, true},
		{"large value", 1000000.0, true},
		{"very small positive", 0.1, false}, // Not divisible by 0.25
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateValue(tt.value)
			if tt.valid && err != nil {
				t.Errorf("ValidateValue(%v) = %v, want nil", tt.value, err)
			}
			if !tt.valid && err == nil {
				t.Errorf("ValidateValue(%v) = nil, want error", tt.value)
			}
		})
	}
}

func TestConversionRoundTrip(t *testing.T) {
	calc := NewMaxCalculator()

	// Start with 1RM of 400, convert to TM, then back to 1RM
	tm, err := calc.ConvertToTM(400, nil)
	if err != nil {
		t.Fatalf("ConvertToTM failed: %v", err)
	}

	oneRM, err := calc.ConvertToOneRM(tm, nil)
	if err != nil {
		t.Fatalf("ConvertToOneRM failed: %v", err)
	}

	if oneRM != 400 {
		t.Errorf("Round trip conversion: 400 -> %v -> %v, expected 400", tm, oneRM)
	}
}

func TestConversionRoundTrip_CustomPercentage(t *testing.T) {
	calc := NewMaxCalculator()
	pct := 85.0

	// 400 * 0.85 = 340, 340 / 0.85 = 400
	tm, err := calc.ConvertToTM(400, &pct)
	if err != nil {
		t.Fatalf("ConvertToTM failed: %v", err)
	}

	oneRM, err := calc.ConvertToOneRM(tm, &pct)
	if err != nil {
		t.Fatalf("ConvertToOneRM failed: %v", err)
	}

	if oneRM != 400 {
		t.Errorf("Round trip conversion with 85%%: 400 -> %v -> %v, expected 400", tm, oneRM)
	}
}

func TestMaxTypeConstants(t *testing.T) {
	// Ensure constants match expected database values
	if OneRM != "ONE_RM" {
		t.Errorf("OneRM = %q, want %q", OneRM, "ONE_RM")
	}
	if TrainingMax != "TRAINING_MAX" {
		t.Errorf("TrainingMax = %q, want %q", TrainingMax, "TRAINING_MAX")
	}
	if E1RM != "E1RM" {
		t.Errorf("E1RM = %q, want %q", E1RM, "E1RM")
	}
}

func TestDefaultConstants(t *testing.T) {
	if DefaultTMPercentage != 90.0 {
		t.Errorf("DefaultTMPercentage = %v, want %v", DefaultTMPercentage, 90.0)
	}
	if TMWarningLowerBound != 80.0 {
		t.Errorf("TMWarningLowerBound = %v, want %v", TMWarningLowerBound, 80.0)
	}
	if TMWarningUpperBound != 95.0 {
		t.Errorf("TMWarningUpperBound = %v, want %v", TMWarningUpperBound, 95.0)
	}
}

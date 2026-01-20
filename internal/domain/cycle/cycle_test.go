package cycle

import (
	"errors"
	"strings"
	"testing"
)

// ==================== Name Validation Tests ====================

func TestValidateName_Valid(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"5/3/1 Cycle"},
		{"Greg Nuckols 3-Week Cycle"},
		{"Starting Strength"},
		{"A"},
		{strings.Repeat("a", MaxNameLength)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.name)
			if err != nil {
				t.Errorf("ValidateName(%q) = %v, want nil", tt.name, err)
			}
		})
	}
}

func TestValidateName_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedErr error
	}{
		{"empty string", "", ErrNameRequired},
		{"only spaces", "   ", ErrNameRequired},
		{"only tabs", "\t\t", ErrNameRequired},
		{"too long", strings.Repeat("a", MaxNameLength+1), ErrNameTooLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateName(%q) = %v, want %v", tt.input, err, tt.expectedErr)
			}
		})
	}
}

// ==================== LengthWeeks Validation Tests ====================

func TestValidateLengthWeeks_Valid(t *testing.T) {
	tests := []struct {
		name        string
		lengthWeeks int
	}{
		{"1 week", 1},
		{"3 weeks", 3},
		{"4 weeks", 4},
		{"12 weeks", 12},
		{"52 weeks", 52},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLengthWeeks(tt.lengthWeeks)
			if err != nil {
				t.Errorf("ValidateLengthWeeks(%d) = %v, want nil", tt.lengthWeeks, err)
			}
		})
	}
}

func TestValidateLengthWeeks_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		lengthWeeks int
		expectedErr error
	}{
		{"zero", 0, ErrLengthWeeksInvalid},
		{"negative", -1, ErrLengthWeeksInvalid},
		{"very negative", -100, ErrLengthWeeksInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLengthWeeks(tt.lengthWeeks)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateLengthWeeks(%d) = %v, want %v", tt.lengthWeeks, err, tt.expectedErr)
			}
		})
	}
}

// ==================== CreateCycle Tests ====================

func TestCreateCycle_ValidInput(t *testing.T) {
	input := CreateCycleInput{
		Name:        "5/3/1 Cycle",
		LengthWeeks: 4,
	}

	cycle, result := CreateCycle(input, "test-id")

	if !result.Valid {
		t.Errorf("CreateCycle returned invalid result: %v", result.Errors)
	}
	if cycle == nil {
		t.Fatal("CreateCycle returned nil cycle")
	}
	if cycle.ID != "test-id" {
		t.Errorf("cycle.ID = %q, want %q", cycle.ID, "test-id")
	}
	if cycle.Name != "5/3/1 Cycle" {
		t.Errorf("cycle.Name = %q, want %q", cycle.Name, "5/3/1 Cycle")
	}
	if cycle.LengthWeeks != 4 {
		t.Errorf("cycle.LengthWeeks = %d, want %d", cycle.LengthWeeks, 4)
	}
}

func TestCreateCycle_TrimsWhitespace(t *testing.T) {
	input := CreateCycleInput{
		Name:        "  5/3/1 Cycle  ",
		LengthWeeks: 4,
	}

	cycle, result := CreateCycle(input, "test-id")

	if !result.Valid {
		t.Errorf("CreateCycle returned invalid result: %v", result.Errors)
	}
	if cycle.Name != "5/3/1 Cycle" {
		t.Errorf("cycle.Name = %q, want %q (trimmed)", cycle.Name, "5/3/1 Cycle")
	}
}

func TestCreateCycle_OneWeekCycle(t *testing.T) {
	input := CreateCycleInput{
		Name:        "Starting Strength Cycle",
		LengthWeeks: 1,
	}

	cycle, result := CreateCycle(input, "test-id")

	if !result.Valid {
		t.Errorf("CreateCycle returned invalid result: %v", result.Errors)
	}
	if cycle.LengthWeeks != 1 {
		t.Errorf("cycle.LengthWeeks = %d, want %d", cycle.LengthWeeks, 1)
	}
}

func TestCreateCycle_EmptyName(t *testing.T) {
	input := CreateCycleInput{
		Name:        "",
		LengthWeeks: 4,
	}

	cycle, result := CreateCycle(input, "test-id")

	if result.Valid {
		t.Error("CreateCycle with empty name returned valid result")
	}
	if cycle != nil {
		t.Error("CreateCycle with invalid input returned non-nil cycle")
	}
}

func TestCreateCycle_InvalidLengthWeeks(t *testing.T) {
	input := CreateCycleInput{
		Name:        "Test Cycle",
		LengthWeeks: 0,
	}

	cycle, result := CreateCycle(input, "test-id")

	if result.Valid {
		t.Error("CreateCycle with zero length_weeks returned valid result")
	}
	if cycle != nil {
		t.Error("CreateCycle with invalid input returned non-nil cycle")
	}
}

func TestCreateCycle_MultipleErrors(t *testing.T) {
	input := CreateCycleInput{
		Name:        "",  // Invalid
		LengthWeeks: -1, // Invalid
	}

	cycle, result := CreateCycle(input, "test-id")

	if result.Valid {
		t.Error("CreateCycle with multiple errors returned valid result")
	}
	if cycle != nil {
		t.Error("CreateCycle with invalid input returned non-nil cycle")
	}
	if len(result.Errors) < 2 {
		t.Errorf("Expected at least 2 errors, got %d", len(result.Errors))
	}
}

// ==================== UpdateCycle Tests ====================

func TestUpdateCycle_UpdateName(t *testing.T) {
	cycle := &Cycle{ID: "test-id", Name: "Old Name", LengthWeeks: 4}

	newName := "New Name"
	input := UpdateCycleInput{Name: &newName}

	result := UpdateCycle(cycle, input)

	if !result.Valid {
		t.Errorf("UpdateCycle returned invalid result: %v", result.Errors)
	}
	if cycle.Name != "New Name" {
		t.Errorf("cycle.Name = %q, want %q", cycle.Name, "New Name")
	}
}

func TestUpdateCycle_UpdateLengthWeeks(t *testing.T) {
	cycle := &Cycle{ID: "test-id", Name: "Test Cycle", LengthWeeks: 4}

	newLengthWeeks := 3
	input := UpdateCycleInput{LengthWeeks: &newLengthWeeks}

	result := UpdateCycle(cycle, input)

	if !result.Valid {
		t.Errorf("UpdateCycle returned invalid result: %v", result.Errors)
	}
	if cycle.LengthWeeks != 3 {
		t.Errorf("cycle.LengthWeeks = %d, want %d", cycle.LengthWeeks, 3)
	}
}

func TestUpdateCycle_UpdateBothFields(t *testing.T) {
	cycle := &Cycle{ID: "test-id", Name: "Old Name", LengthWeeks: 4}

	newName := "New Name"
	newLengthWeeks := 3
	input := UpdateCycleInput{
		Name:        &newName,
		LengthWeeks: &newLengthWeeks,
	}

	result := UpdateCycle(cycle, input)

	if !result.Valid {
		t.Errorf("UpdateCycle returned invalid result: %v", result.Errors)
	}
	if cycle.Name != "New Name" {
		t.Errorf("cycle.Name = %q, want %q", cycle.Name, "New Name")
	}
	if cycle.LengthWeeks != 3 {
		t.Errorf("cycle.LengthWeeks = %d, want %d", cycle.LengthWeeks, 3)
	}
}

func TestUpdateCycle_TrimsWhitespace(t *testing.T) {
	cycle := &Cycle{ID: "test-id", Name: "Old Name", LengthWeeks: 4}

	newName := "  New Name  "
	input := UpdateCycleInput{Name: &newName}

	result := UpdateCycle(cycle, input)

	if !result.Valid {
		t.Errorf("UpdateCycle returned invalid result: %v", result.Errors)
	}
	if cycle.Name != "New Name" {
		t.Errorf("cycle.Name = %q, want %q (trimmed)", cycle.Name, "New Name")
	}
}

func TestUpdateCycle_InvalidName(t *testing.T) {
	cycle := &Cycle{ID: "test-id", Name: "Old Name", LengthWeeks: 4}
	originalName := cycle.Name

	emptyName := ""
	input := UpdateCycleInput{Name: &emptyName}

	result := UpdateCycle(cycle, input)

	if result.Valid {
		t.Error("UpdateCycle with invalid name returned valid result")
	}
	if cycle.Name != originalName {
		t.Errorf("cycle.Name was changed despite validation failure")
	}
}

func TestUpdateCycle_InvalidLengthWeeks(t *testing.T) {
	cycle := &Cycle{ID: "test-id", Name: "Test Cycle", LengthWeeks: 4}
	originalLengthWeeks := cycle.LengthWeeks

	zeroLengthWeeks := 0
	input := UpdateCycleInput{LengthWeeks: &zeroLengthWeeks}

	result := UpdateCycle(cycle, input)

	if result.Valid {
		t.Error("UpdateCycle with invalid length_weeks returned valid result")
	}
	if cycle.LengthWeeks != originalLengthWeeks {
		t.Errorf("cycle.LengthWeeks was changed despite validation failure")
	}
}

func TestUpdateCycle_NoChanges(t *testing.T) {
	cycle := &Cycle{ID: "test-id", Name: "Test Cycle", LengthWeeks: 4}
	originalUpdatedAt := cycle.UpdatedAt

	input := UpdateCycleInput{} // No changes

	result := UpdateCycle(cycle, input)

	if !result.Valid {
		t.Errorf("UpdateCycle with no changes returned invalid result: %v", result.Errors)
	}
	if !cycle.UpdatedAt.After(originalUpdatedAt) && cycle.UpdatedAt.Equal(originalUpdatedAt) {
		// This is expected - UpdatedAt is updated when result is valid
	}
}

// ==================== Cycle.Validate Tests ====================

func TestCycle_Validate_Valid(t *testing.T) {
	cycle := &Cycle{
		ID:          "test-id",
		Name:        "Test Cycle",
		LengthWeeks: 4,
	}

	result := cycle.Validate()

	if !result.Valid {
		t.Errorf("Validate returned invalid result for valid cycle: %v", result.Errors)
	}
}

func TestCycle_Validate_InvalidName(t *testing.T) {
	cycle := &Cycle{
		ID:          "test-id",
		Name:        "",
		LengthWeeks: 4,
	}

	result := cycle.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for cycle with invalid name")
	}
}

func TestCycle_Validate_InvalidLengthWeeks(t *testing.T) {
	cycle := &Cycle{
		ID:          "test-id",
		Name:        "Test Cycle",
		LengthWeeks: 0,
	}

	result := cycle.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for cycle with invalid length_weeks")
	}
}

func TestCycle_Validate_MultipleErrors(t *testing.T) {
	cycle := &Cycle{
		ID:          "test-id",
		Name:        "",
		LengthWeeks: 0,
	}

	result := cycle.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for cycle with multiple errors")
	}
	if len(result.Errors) < 2 {
		t.Errorf("Expected at least 2 errors, got %d", len(result.Errors))
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
	result.AddError(ErrNameRequired)
	result.AddError(ErrLengthWeeksInvalid)

	err := result.Error()
	if err == nil {
		t.Error("ValidationResult.Error() = nil, want error")
	}

	errMsg := err.Error()
	if errMsg == "" {
		t.Error("Error message should not be empty")
	}
}

// ==================== Edge Case Tests ====================

func TestCreateCycle_MaxLengthName(t *testing.T) {
	input := CreateCycleInput{
		Name:        strings.Repeat("a", MaxNameLength),
		LengthWeeks: 4,
	}

	cycle, result := CreateCycle(input, "test-id")

	if !result.Valid {
		t.Errorf("CreateCycle with max length name returned invalid result: %v", result.Errors)
	}
	if cycle == nil {
		t.Fatal("CreateCycle returned nil cycle")
	}
	if len(cycle.Name) != MaxNameLength {
		t.Errorf("cycle.Name length = %d, want %d", len(cycle.Name), MaxNameLength)
	}
}

func TestCreateCycle_NameTooLong(t *testing.T) {
	input := CreateCycleInput{
		Name:        strings.Repeat("a", MaxNameLength+1),
		LengthWeeks: 4,
	}

	cycle, result := CreateCycle(input, "test-id")

	if result.Valid {
		t.Error("CreateCycle with name too long returned valid result")
	}
	if cycle != nil {
		t.Error("CreateCycle with invalid input returned non-nil cycle")
	}
}

func TestCreateCycle_NegativeLengthWeeks(t *testing.T) {
	input := CreateCycleInput{
		Name:        "Test Cycle",
		LengthWeeks: -5,
	}

	cycle, result := CreateCycle(input, "test-id")

	if result.Valid {
		t.Error("CreateCycle with negative length_weeks returned valid result")
	}
	if cycle != nil {
		t.Error("CreateCycle with invalid input returned non-nil cycle")
	}
}

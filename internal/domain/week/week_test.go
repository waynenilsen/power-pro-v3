package week

import (
	"errors"
	"testing"
)

// ==================== Week Number Validation Tests ====================

func TestValidateWeekNumber_Valid(t *testing.T) {
	tests := []struct {
		name       string
		weekNumber int
	}{
		{"week 1", 1},
		{"week 2", 2},
		{"week 4", 4},
		{"week 12", 12},
		{"week 52", 52},
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

// ==================== Cycle ID Validation Tests ====================

func TestValidateCycleID_Valid(t *testing.T) {
	tests := []string{
		"cycle-123",
		"abc",
		"a",
		"00000000-0000-0000-0000-000000000001",
	}

	for _, cycleID := range tests {
		t.Run(cycleID, func(t *testing.T) {
			err := ValidateCycleID(cycleID)
			if err != nil {
				t.Errorf("ValidateCycleID(%q) = %v, want nil", cycleID, err)
			}
		})
	}
}

func TestValidateCycleID_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		cycleID     string
		expectedErr error
	}{
		{"empty string", "", ErrCycleIDRequired},
		{"only spaces", "   ", ErrCycleIDRequired},
		{"only tabs", "\t\t", ErrCycleIDRequired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCycleID(tt.cycleID)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateCycleID(%q) = %v, want %v", tt.cycleID, err, tt.expectedErr)
			}
		})
	}
}

// ==================== Variant Validation Tests ====================

func TestValidateVariant_Valid(t *testing.T) {
	tests := []struct {
		name    string
		variant *string
	}{
		{"nil", nil},
		{"variant A", strPtr("A")},
		{"variant B", strPtr("B")},
		{"empty string", strPtr("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVariant(tt.variant)
			if err != nil {
				t.Errorf("ValidateVariant(%v) = %v, want nil", ptrStr(tt.variant), err)
			}
		})
	}
}

func TestValidateVariant_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		variant     *string
		expectedErr error
	}{
		{"lowercase a", strPtr("a"), ErrVariantInvalid},
		{"lowercase b", strPtr("b"), ErrVariantInvalid},
		{"variant C", strPtr("C"), ErrVariantInvalid},
		{"some text", strPtr("some text"), ErrVariantInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVariant(tt.variant)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateVariant(%v) = %v, want %v", ptrStr(tt.variant), err, tt.expectedErr)
			}
		})
	}
}

// ==================== DayOfWeek Validation Tests ====================

func TestValidateDayOfWeek_Valid(t *testing.T) {
	tests := []string{
		"MONDAY",
		"TUESDAY",
		"WEDNESDAY",
		"THURSDAY",
		"FRIDAY",
		"SATURDAY",
		"SUNDAY",
	}

	for _, dow := range tests {
		t.Run(dow, func(t *testing.T) {
			err := ValidateDayOfWeek(dow)
			if err != nil {
				t.Errorf("ValidateDayOfWeek(%q) = %v, want nil", dow, err)
			}
		})
	}
}

func TestValidateDayOfWeek_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		dayOfWeek   string
		expectedErr error
	}{
		{"empty string", "", ErrDayOfWeekRequired},
		{"lowercase monday", "monday", ErrDayOfWeekInvalid},
		{"abbreviated", "MON", ErrDayOfWeekInvalid},
		{"number", "1", ErrDayOfWeekInvalid},
		{"random text", "someday", ErrDayOfWeekInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDayOfWeek(tt.dayOfWeek)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateDayOfWeek(%q) = %v, want %v", tt.dayOfWeek, err, tt.expectedErr)
			}
		})
	}
}

// ==================== CreateWeek Tests ====================

func TestCreateWeek_ValidInput(t *testing.T) {
	input := CreateWeekInput{
		WeekNumber: 1,
		CycleID:    "cycle-123",
	}

	week, result := CreateWeek(input, "test-id")

	if !result.Valid {
		t.Errorf("CreateWeek returned invalid result: %v", result.Errors)
	}
	if week == nil {
		t.Fatal("CreateWeek returned nil week")
	}
	if week.ID != "test-id" {
		t.Errorf("week.ID = %q, want %q", week.ID, "test-id")
	}
	if week.WeekNumber != 1 {
		t.Errorf("week.WeekNumber = %d, want %d", week.WeekNumber, 1)
	}
	if week.CycleID != "cycle-123" {
		t.Errorf("week.CycleID = %q, want %q", week.CycleID, "cycle-123")
	}
	if week.Variant != nil {
		t.Errorf("week.Variant = %v, want nil", week.Variant)
	}
}

func TestCreateWeek_WithVariantA(t *testing.T) {
	variant := "A"
	input := CreateWeekInput{
		WeekNumber: 1,
		Variant:    &variant,
		CycleID:    "cycle-123",
	}

	week, result := CreateWeek(input, "test-id")

	if !result.Valid {
		t.Errorf("CreateWeek returned invalid result: %v", result.Errors)
	}
	if week.Variant == nil || *week.Variant != "A" {
		t.Errorf("week.Variant = %v, want %q", ptrStr(week.Variant), "A")
	}
}

func TestCreateWeek_WithVariantB(t *testing.T) {
	variant := "B"
	input := CreateWeekInput{
		WeekNumber: 2,
		Variant:    &variant,
		CycleID:    "cycle-456",
	}

	week, result := CreateWeek(input, "test-id")

	if !result.Valid {
		t.Errorf("CreateWeek returned invalid result: %v", result.Errors)
	}
	if week.Variant == nil || *week.Variant != "B" {
		t.Errorf("week.Variant = %v, want %q", ptrStr(week.Variant), "B")
	}
}

func TestCreateWeek_EmptyVariantBecomesNil(t *testing.T) {
	variant := ""
	input := CreateWeekInput{
		WeekNumber: 1,
		Variant:    &variant,
		CycleID:    "cycle-123",
	}

	week, result := CreateWeek(input, "test-id")

	if !result.Valid {
		t.Errorf("CreateWeek returned invalid result: %v", result.Errors)
	}
	if week.Variant != nil {
		t.Errorf("week.Variant = %v, want nil (empty string normalized to nil)", week.Variant)
	}
}

func TestCreateWeek_InvalidWeekNumber(t *testing.T) {
	input := CreateWeekInput{
		WeekNumber: 0,
		CycleID:    "cycle-123",
	}

	week, result := CreateWeek(input, "test-id")

	if result.Valid {
		t.Error("CreateWeek with zero week number returned valid result")
	}
	if week != nil {
		t.Error("CreateWeek with invalid input returned non-nil week")
	}
}

func TestCreateWeek_EmptyCycleID(t *testing.T) {
	input := CreateWeekInput{
		WeekNumber: 1,
		CycleID:    "",
	}

	week, result := CreateWeek(input, "test-id")

	if result.Valid {
		t.Error("CreateWeek with empty cycle ID returned valid result")
	}
	if week != nil {
		t.Error("CreateWeek with invalid input returned non-nil week")
	}
}

func TestCreateWeek_InvalidVariant(t *testing.T) {
	variant := "C"
	input := CreateWeekInput{
		WeekNumber: 1,
		Variant:    &variant,
		CycleID:    "cycle-123",
	}

	week, result := CreateWeek(input, "test-id")

	if result.Valid {
		t.Error("CreateWeek with invalid variant returned valid result")
	}
	if week != nil {
		t.Error("CreateWeek with invalid input returned non-nil week")
	}
}

func TestCreateWeek_MultipleErrors(t *testing.T) {
	variant := "invalid"
	input := CreateWeekInput{
		WeekNumber: -1,    // Invalid
		Variant:    &variant, // Invalid
		CycleID:    "",    // Invalid
	}

	week, result := CreateWeek(input, "test-id")

	if result.Valid {
		t.Error("CreateWeek with multiple errors returned valid result")
	}
	if week != nil {
		t.Error("CreateWeek with invalid input returned non-nil week")
	}
	if len(result.Errors) < 3 {
		t.Errorf("Expected at least 3 errors, got %d", len(result.Errors))
	}
}

// ==================== UpdateWeek Tests ====================

func TestUpdateWeek_UpdateWeekNumber(t *testing.T) {
	week := &Week{ID: "test-id", WeekNumber: 1, CycleID: "cycle-123"}

	newWeekNumber := 2
	input := UpdateWeekInput{WeekNumber: &newWeekNumber}

	result := UpdateWeek(week, input)

	if !result.Valid {
		t.Errorf("UpdateWeek returned invalid result: %v", result.Errors)
	}
	if week.WeekNumber != 2 {
		t.Errorf("week.WeekNumber = %d, want %d", week.WeekNumber, 2)
	}
}

func TestUpdateWeek_UpdateVariant(t *testing.T) {
	week := &Week{ID: "test-id", WeekNumber: 1, CycleID: "cycle-123"}

	variant := "A"
	input := UpdateWeekInput{Variant: &variant}

	result := UpdateWeek(week, input)

	if !result.Valid {
		t.Errorf("UpdateWeek returned invalid result: %v", result.Errors)
	}
	if week.Variant == nil || *week.Variant != "A" {
		t.Errorf("week.Variant = %v, want %q", ptrStr(week.Variant), "A")
	}
}

func TestUpdateWeek_ClearVariant(t *testing.T) {
	variant := "A"
	week := &Week{ID: "test-id", WeekNumber: 1, Variant: &variant, CycleID: "cycle-123"}

	input := UpdateWeekInput{ClearVariant: true}

	result := UpdateWeek(week, input)

	if !result.Valid {
		t.Errorf("UpdateWeek returned invalid result: %v", result.Errors)
	}
	if week.Variant != nil {
		t.Errorf("week.Variant = %v, want nil", week.Variant)
	}
}

func TestUpdateWeek_UpdateCycleID(t *testing.T) {
	week := &Week{ID: "test-id", WeekNumber: 1, CycleID: "cycle-123"}

	newCycleID := "cycle-456"
	input := UpdateWeekInput{CycleID: &newCycleID}

	result := UpdateWeek(week, input)

	if !result.Valid {
		t.Errorf("UpdateWeek returned invalid result: %v", result.Errors)
	}
	if week.CycleID != "cycle-456" {
		t.Errorf("week.CycleID = %q, want %q", week.CycleID, "cycle-456")
	}
}

func TestUpdateWeek_InvalidWeekNumber(t *testing.T) {
	week := &Week{ID: "test-id", WeekNumber: 1, CycleID: "cycle-123"}
	originalWeekNumber := week.WeekNumber

	zeroWeekNumber := 0
	input := UpdateWeekInput{WeekNumber: &zeroWeekNumber}

	result := UpdateWeek(week, input)

	if result.Valid {
		t.Error("UpdateWeek with invalid week number returned valid result")
	}
	if week.WeekNumber != originalWeekNumber {
		t.Errorf("week.WeekNumber was changed despite validation failure")
	}
}

func TestUpdateWeek_NoChanges(t *testing.T) {
	week := &Week{ID: "test-id", WeekNumber: 1, CycleID: "cycle-123"}
	originalUpdatedAt := week.UpdatedAt

	input := UpdateWeekInput{} // No changes

	result := UpdateWeek(week, input)

	if !result.Valid {
		t.Errorf("UpdateWeek with no changes returned invalid result: %v", result.Errors)
	}
	if !week.UpdatedAt.After(originalUpdatedAt) && week.UpdatedAt.Equal(originalUpdatedAt) {
		// This is expected - UpdatedAt is updated when result is valid
	}
}

func TestUpdateWeek_InvalidVariant(t *testing.T) {
	variant := "A"
	week := &Week{ID: "test-id", WeekNumber: 1, Variant: &variant, CycleID: "cycle-123"}
	originalVariant := week.Variant

	invalidVariant := "C" // Invalid variant
	input := UpdateWeekInput{Variant: &invalidVariant}

	result := UpdateWeek(week, input)

	if result.Valid {
		t.Error("UpdateWeek with invalid variant returned valid result")
	}
	if week.Variant != originalVariant {
		t.Errorf("week.Variant was changed despite validation failure")
	}
}

func TestUpdateWeek_InvalidCycleID(t *testing.T) {
	week := &Week{ID: "test-id", WeekNumber: 1, CycleID: "cycle-123"}
	originalCycleID := week.CycleID

	emptyCycleID := ""
	input := UpdateWeekInput{CycleID: &emptyCycleID}

	result := UpdateWeek(week, input)

	if result.Valid {
		t.Error("UpdateWeek with invalid cycle ID returned valid result")
	}
	if week.CycleID != originalCycleID {
		t.Errorf("week.CycleID was changed despite validation failure")
	}
}

func TestUpdateWeek_EmptyVariantBecomesNil(t *testing.T) {
	variant := "A"
	week := &Week{ID: "test-id", WeekNumber: 1, Variant: &variant, CycleID: "cycle-123"}

	emptyVariant := ""
	input := UpdateWeekInput{Variant: &emptyVariant}

	result := UpdateWeek(week, input)

	if !result.Valid {
		t.Errorf("UpdateWeek returned invalid result: %v", result.Errors)
	}
	if week.Variant != nil {
		t.Errorf("week.Variant = %v, want nil (empty string normalized to nil)", week.Variant)
	}
}

// ==================== Week.Validate Tests ====================

func TestWeek_Validate_Valid(t *testing.T) {
	week := &Week{
		ID:         "test-id",
		WeekNumber: 1,
		CycleID:    "cycle-123",
	}

	result := week.Validate()

	if !result.Valid {
		t.Errorf("Validate returned invalid result for valid week: %v", result.Errors)
	}
}

func TestWeek_Validate_InvalidWeekNumber(t *testing.T) {
	week := &Week{
		ID:         "test-id",
		WeekNumber: 0,
		CycleID:    "cycle-123",
	}

	result := week.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for week with invalid week number")
	}
}

func TestWeek_Validate_InvalidCycleID(t *testing.T) {
	week := &Week{
		ID:         "test-id",
		WeekNumber: 1,
		CycleID:    "",
	}

	result := week.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for week with empty cycle ID")
	}
}

func TestWeek_Validate_InvalidVariant(t *testing.T) {
	invalidVariant := "C"
	week := &Week{
		ID:         "test-id",
		WeekNumber: 1,
		CycleID:    "cycle-123",
		Variant:    &invalidVariant,
	}

	result := week.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for week with invalid variant")
	}
}

func TestWeek_Validate_MultipleErrors(t *testing.T) {
	invalidVariant := "invalid"
	week := &Week{
		ID:         "test-id",
		WeekNumber: 0,
		CycleID:    "",
		Variant:    &invalidVariant,
	}

	result := week.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for week with multiple errors")
	}
	if len(result.Errors) < 2 {
		t.Errorf("Expected at least 2 errors, got %d", len(result.Errors))
	}
}

// ==================== CreateWeekDay Tests ====================

func TestCreateWeekDay_Valid(t *testing.T) {
	input := CreateWeekDayInput{
		WeekID:    "week-123",
		DayID:     "day-456",
		DayOfWeek: "MONDAY",
	}

	wd, result := CreateWeekDay(input, "wd-id")

	if !result.Valid {
		t.Errorf("CreateWeekDay returned invalid result: %v", result.Errors)
	}
	if wd == nil {
		t.Fatal("CreateWeekDay returned nil")
	}
	if wd.ID != "wd-id" {
		t.Errorf("wd.ID = %q, want %q", wd.ID, "wd-id")
	}
	if wd.WeekID != "week-123" {
		t.Errorf("wd.WeekID = %q, want %q", wd.WeekID, "week-123")
	}
	if wd.DayID != "day-456" {
		t.Errorf("wd.DayID = %q, want %q", wd.DayID, "day-456")
	}
	if wd.DayOfWeek != Monday {
		t.Errorf("wd.DayOfWeek = %q, want %q", wd.DayOfWeek, Monday)
	}
}

func TestCreateWeekDay_AllDaysOfWeek(t *testing.T) {
	daysOfWeek := []string{"MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY", "SATURDAY", "SUNDAY"}

	for _, dow := range daysOfWeek {
		t.Run(dow, func(t *testing.T) {
			input := CreateWeekDayInput{
				WeekID:    "week-123",
				DayID:     "day-456",
				DayOfWeek: dow,
			}

			wd, result := CreateWeekDay(input, "wd-id")

			if !result.Valid {
				t.Errorf("CreateWeekDay returned invalid result: %v", result.Errors)
			}
			if wd.DayOfWeek != DayOfWeek(dow) {
				t.Errorf("wd.DayOfWeek = %q, want %q", wd.DayOfWeek, dow)
			}
		})
	}
}

func TestCreateWeekDay_MissingWeekID(t *testing.T) {
	input := CreateWeekDayInput{
		WeekID:    "",
		DayID:     "day-456",
		DayOfWeek: "MONDAY",
	}

	wd, result := CreateWeekDay(input, "wd-id")

	if result.Valid {
		t.Error("CreateWeekDay with missing week ID returned valid result")
	}
	if wd != nil {
		t.Error("CreateWeekDay with invalid input returned non-nil")
	}
}

func TestCreateWeekDay_MissingDayID(t *testing.T) {
	input := CreateWeekDayInput{
		WeekID:    "week-123",
		DayID:     "",
		DayOfWeek: "MONDAY",
	}

	wd, result := CreateWeekDay(input, "wd-id")

	if result.Valid {
		t.Error("CreateWeekDay with missing day ID returned valid result")
	}
	if wd != nil {
		t.Error("CreateWeekDay with invalid input returned non-nil")
	}
}

func TestCreateWeekDay_InvalidDayOfWeek(t *testing.T) {
	input := CreateWeekDayInput{
		WeekID:    "week-123",
		DayID:     "day-456",
		DayOfWeek: "invalid",
	}

	wd, result := CreateWeekDay(input, "wd-id")

	if result.Valid {
		t.Error("CreateWeekDay with invalid day_of_week returned valid result")
	}
	if wd != nil {
		t.Error("CreateWeekDay with invalid input returned non-nil")
	}
}

// ==================== DayOfWeekOrder Tests ====================

func TestDayOfWeekOrder(t *testing.T) {
	tests := []struct {
		dow           DayOfWeek
		expectedOrder int
	}{
		{Monday, 1},
		{Tuesday, 2},
		{Wednesday, 3},
		{Thursday, 4},
		{Friday, 5},
		{Saturday, 6},
		{Sunday, 7},
		{DayOfWeek("INVALID"), 8},
	}

	for _, tt := range tests {
		t.Run(string(tt.dow), func(t *testing.T) {
			order := DayOfWeekOrder(tt.dow)
			if order != tt.expectedOrder {
				t.Errorf("DayOfWeekOrder(%q) = %d, want %d", tt.dow, order, tt.expectedOrder)
			}
		})
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
	result.AddError(ErrWeekNumberInvalid)
	result.AddError(ErrCycleIDRequired)

	err := result.Error()
	if err == nil {
		t.Error("ValidationResult.Error() = nil, want error")
	}

	errMsg := err.Error()
	if errMsg == "" {
		t.Error("Error message should not be empty")
	}
}

// ==================== Helper Functions ====================

func strPtr(s string) *string {
	return &s
}

func ptrStr(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}

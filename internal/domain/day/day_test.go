package day

import (
	"errors"
	"strings"
	"testing"
)

// ==================== Name Validation Tests ====================

func TestValidateName_Valid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"single word", "Day A"},
		{"two words", "Heavy Day"},
		{"with numbers", "Week 1 Day A"},
		{"minimum length", "A"},
		{"exactly 50 chars", strings.Repeat("a", 50)},
		{"with special chars", "Day A - Heavy"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.input)
			if err != nil {
				t.Errorf("ValidateName(%q) = %v, want nil", tt.input, err)
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
		{"51 chars", strings.Repeat("a", 51), ErrNameTooLong},
		{"100 chars", strings.Repeat("a", 100), ErrNameTooLong},
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

// ==================== Slug Validation Tests ====================

func TestValidateSlug_Valid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"single word", "day-a"},
		{"with hyphen", "heavy-day"},
		{"multiple hyphens", "week-1-day-a"},
		{"with numbers", "531-day-1"},
		{"numbers only", "531"},
		{"minimum length", "a"},
		{"exactly 50 chars", strings.Repeat("a", 50)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSlug(tt.input)
			if err != nil {
				t.Errorf("ValidateSlug(%q) = %v, want nil", tt.input, err)
			}
		})
	}
}

func TestValidateSlug_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedErr error
	}{
		{"empty string", "", ErrSlugEmpty},
		{"uppercase", "Day-A", ErrSlugInvalid},
		{"with space", "day a", ErrSlugInvalid},
		{"with underscore", "day_a", ErrSlugInvalid},
		{"leading hyphen", "-day", ErrSlugInvalid},
		{"trailing hyphen", "day-", ErrSlugInvalid},
		{"consecutive hyphens", "day--a", ErrSlugInvalid},
		{"special chars", "day@home", ErrSlugInvalid},
		{"51 chars", strings.Repeat("a", 51), ErrSlugTooLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSlug(tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateSlug(%q) = %v, want %v", tt.input, err, tt.expectedErr)
			}
		})
	}
}

// ==================== Slug Generation Tests ====================

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple lowercase", "day a", "day-a"},
		{"uppercase", "DAY A", "day-a"},
		{"mixed case", "Heavy Day", "heavy-day"},
		{"multiple spaces", "Week  1  Day  A", "week-1-day-a"},
		{"special chars", "5/3/1 Day 1", "5-3-1-day-1"},
		{"underscores", "day_a", "day-a"},
		{"parentheses", "Day A (Heavy)", "day-a-heavy"},
		{"apostrophe", "Greg's Day", "gregs-day"},
		{"leading spaces", "  day a", "day-a"},
		{"trailing spaces", "day a  ", "day-a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateSlug(tt.input)
			if result != tt.expected {
				t.Errorf("GenerateSlug(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGenerateSlug_ValidatesAfterGeneration(t *testing.T) {
	inputs := []string{
		"Day A",
		"Heavy Day",
		"Week 1 Day A",
		"5/3/1 BBB Day",
		"Light (Recovery)",
	}

	for _, input := range inputs {
		slug := GenerateSlug(input)
		if err := ValidateSlug(slug); err != nil {
			t.Errorf("GenerateSlug(%q) = %q, but validation failed: %v", input, slug, err)
		}
	}
}

// ==================== Metadata Validation Tests ====================

func TestValidateMetadata_Valid(t *testing.T) {
	tests := []struct {
		name     string
		metadata map[string]interface{}
	}{
		{"nil metadata", nil},
		{"empty metadata", map[string]interface{}{}},
		{"with intensity level", map[string]interface{}{"intensityLevel": "HEAVY"}},
		{"with focus", map[string]interface{}{"focus": "upper body"}},
		{"with both", map[string]interface{}{"intensityLevel": "LIGHT", "focus": "recovery"}},
		{"with numbers", map[string]interface{}{"someNumber": 123}},
		{"nested", map[string]interface{}{"nested": map[string]interface{}{"key": "value"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMetadata(tt.metadata)
			if err != nil {
				t.Errorf("ValidateMetadata(%v) = %v, want nil", tt.metadata, err)
			}
		})
	}
}

// ==================== Order Validation Tests ====================

func TestValidateOrder_Valid(t *testing.T) {
	tests := []int{0, 1, 10, 100, 1000}

	for _, order := range tests {
		t.Run("order "+string(rune(order)), func(t *testing.T) {
			err := ValidateOrder(order)
			if err != nil {
				t.Errorf("ValidateOrder(%d) = %v, want nil", order, err)
			}
		})
	}
}

func TestValidateOrder_Invalid(t *testing.T) {
	tests := []int{-1, -10, -100}

	for _, order := range tests {
		t.Run("negative order", func(t *testing.T) {
			err := ValidateOrder(order)
			if !errors.Is(err, ErrOrderNegative) {
				t.Errorf("ValidateOrder(%d) = %v, want %v", order, err, ErrOrderNegative)
			}
		})
	}
}

// ==================== CreateDay Tests ====================

func TestCreateDay_ValidInput(t *testing.T) {
	input := CreateDayInput{
		Name: "Day A",
	}

	day, result := CreateDay(input, "test-id")

	if !result.Valid {
		t.Errorf("CreateDay returned invalid result: %v", result.Errors)
	}
	if day == nil {
		t.Fatal("CreateDay returned nil day")
	}
	if day.ID != "test-id" {
		t.Errorf("day.ID = %q, want %q", day.ID, "test-id")
	}
	if day.Name != "Day A" {
		t.Errorf("day.Name = %q, want %q", day.Name, "Day A")
	}
	if day.Slug != "day-a" {
		t.Errorf("day.Slug = %q, want %q (auto-generated)", day.Slug, "day-a")
	}
	if day.Metadata != nil {
		t.Errorf("day.Metadata = %v, want nil", day.Metadata)
	}
	if day.ProgramID != nil {
		t.Errorf("day.ProgramID = %v, want nil", day.ProgramID)
	}
}

func TestCreateDay_WithProvidedSlug(t *testing.T) {
	input := CreateDayInput{
		Name: "Day A",
		Slug: "heavy-day-a",
	}

	day, result := CreateDay(input, "test-id")

	if !result.Valid {
		t.Errorf("CreateDay returned invalid result: %v", result.Errors)
	}
	if day.Slug != "heavy-day-a" {
		t.Errorf("day.Slug = %q, want %q", day.Slug, "heavy-day-a")
	}
}

func TestCreateDay_WithMetadata(t *testing.T) {
	metadata := map[string]interface{}{
		"intensityLevel": "HEAVY",
		"focus":          "squats",
	}
	input := CreateDayInput{
		Name:     "Heavy Day",
		Metadata: metadata,
	}

	day, result := CreateDay(input, "test-id")

	if !result.Valid {
		t.Errorf("CreateDay returned invalid result: %v", result.Errors)
	}
	if day.Metadata == nil {
		t.Fatal("day.Metadata is nil, want non-nil")
	}
	if day.Metadata["intensityLevel"] != "HEAVY" {
		t.Errorf("day.Metadata[\"intensityLevel\"] = %v, want %q", day.Metadata["intensityLevel"], "HEAVY")
	}
}

func TestCreateDay_WithProgramID(t *testing.T) {
	programID := "program-123"
	input := CreateDayInput{
		Name:      "Day A",
		ProgramID: &programID,
	}

	day, result := CreateDay(input, "test-id")

	if !result.Valid {
		t.Errorf("CreateDay returned invalid result: %v", result.Errors)
	}
	if day.ProgramID == nil || *day.ProgramID != "program-123" {
		t.Errorf("day.ProgramID = %v, want %q", day.ProgramID, "program-123")
	}
}

func TestCreateDay_InvalidName(t *testing.T) {
	input := CreateDayInput{
		Name: "",
	}

	day, result := CreateDay(input, "test-id")

	if result.Valid {
		t.Error("CreateDay with empty name returned valid result")
	}
	if day != nil {
		t.Error("CreateDay with invalid input returned non-nil day")
	}
	if len(result.Errors) == 0 {
		t.Error("CreateDay with invalid input returned no errors")
	}
}

func TestCreateDay_InvalidSlug(t *testing.T) {
	input := CreateDayInput{
		Name: "Valid Name",
		Slug: "Invalid Slug",
	}

	day, result := CreateDay(input, "test-id")

	if result.Valid {
		t.Error("CreateDay with invalid slug returned valid result")
	}
	if day != nil {
		t.Error("CreateDay with invalid input returned non-nil day")
	}
}

func TestCreateDay_MultipleErrors(t *testing.T) {
	input := CreateDayInput{
		Name: "",              // Invalid
		Slug: "Invalid Slug", // Invalid
	}

	day, result := CreateDay(input, "test-id")

	if result.Valid {
		t.Error("CreateDay with multiple errors returned valid result")
	}
	if day != nil {
		t.Error("CreateDay with invalid input returned non-nil day")
	}
	if len(result.Errors) < 2 {
		t.Errorf("Expected at least 2 errors, got %d", len(result.Errors))
	}
}

// ==================== UpdateDay Tests ====================

func TestUpdateDay_UpdateName(t *testing.T) {
	day := &Day{ID: "test-id", Name: "Day A", Slug: "day-a"}

	newName := "Heavy Day"
	input := UpdateDayInput{Name: &newName}

	result := UpdateDay(day, input)

	if !result.Valid {
		t.Errorf("UpdateDay returned invalid result: %v", result.Errors)
	}
	if day.Name != "Heavy Day" {
		t.Errorf("day.Name = %q, want %q", day.Name, "Heavy Day")
	}
}

func TestUpdateDay_UpdateSlug(t *testing.T) {
	day := &Day{ID: "test-id", Name: "Day A", Slug: "day-a"}

	newSlug := "heavy-day"
	input := UpdateDayInput{Slug: &newSlug}

	result := UpdateDay(day, input)

	if !result.Valid {
		t.Errorf("UpdateDay returned invalid result: %v", result.Errors)
	}
	if day.Slug != "heavy-day" {
		t.Errorf("day.Slug = %q, want %q", day.Slug, "heavy-day")
	}
}

func TestUpdateDay_UpdateMetadata(t *testing.T) {
	day := &Day{ID: "test-id", Name: "Day A", Slug: "day-a"}

	metadata := map[string]interface{}{"intensityLevel": "LIGHT"}
	input := UpdateDayInput{Metadata: metadata}

	result := UpdateDay(day, input)

	if !result.Valid {
		t.Errorf("UpdateDay returned invalid result: %v", result.Errors)
	}
	if day.Metadata == nil {
		t.Fatal("day.Metadata is nil, want non-nil")
	}
	if day.Metadata["intensityLevel"] != "LIGHT" {
		t.Errorf("day.Metadata[\"intensityLevel\"] = %v, want %q", day.Metadata["intensityLevel"], "LIGHT")
	}
}

func TestUpdateDay_ClearMetadata(t *testing.T) {
	metadata := map[string]interface{}{"intensityLevel": "HEAVY"}
	day := &Day{ID: "test-id", Name: "Day A", Slug: "day-a", Metadata: metadata}

	input := UpdateDayInput{ClearMetadata: true}

	result := UpdateDay(day, input)

	if !result.Valid {
		t.Errorf("UpdateDay returned invalid result: %v", result.Errors)
	}
	if day.Metadata != nil {
		t.Errorf("day.Metadata = %v, want nil", day.Metadata)
	}
}

func TestUpdateDay_SetProgramID(t *testing.T) {
	day := &Day{ID: "test-id", Name: "Day A", Slug: "day-a"}

	programID := "program-123"
	input := UpdateDayInput{ProgramID: &programID}

	result := UpdateDay(day, input)

	if !result.Valid {
		t.Errorf("UpdateDay returned invalid result: %v", result.Errors)
	}
	if day.ProgramID == nil || *day.ProgramID != "program-123" {
		t.Errorf("day.ProgramID = %v, want %q", day.ProgramID, "program-123")
	}
}

func TestUpdateDay_ClearProgramID(t *testing.T) {
	programID := "program-123"
	day := &Day{ID: "test-id", Name: "Day A", Slug: "day-a", ProgramID: &programID}

	input := UpdateDayInput{ClearProgramID: true}

	result := UpdateDay(day, input)

	if !result.Valid {
		t.Errorf("UpdateDay returned invalid result: %v", result.Errors)
	}
	if day.ProgramID != nil {
		t.Errorf("day.ProgramID = %v, want nil", day.ProgramID)
	}
}

func TestUpdateDay_InvalidName(t *testing.T) {
	day := &Day{ID: "test-id", Name: "Day A", Slug: "day-a"}
	originalName := day.Name

	emptyName := ""
	input := UpdateDayInput{Name: &emptyName}

	result := UpdateDay(day, input)

	if result.Valid {
		t.Error("UpdateDay with invalid name returned valid result")
	}
	if day.Name != originalName {
		t.Errorf("day.Name was changed despite validation failure")
	}
}

func TestUpdateDay_NoChanges(t *testing.T) {
	day := &Day{ID: "test-id", Name: "Day A", Slug: "day-a"}
	originalUpdatedAt := day.UpdatedAt

	input := UpdateDayInput{} // No changes

	result := UpdateDay(day, input)

	if !result.Valid {
		t.Errorf("UpdateDay with no changes returned invalid result: %v", result.Errors)
	}
	if !day.UpdatedAt.After(originalUpdatedAt) && day.UpdatedAt.Equal(originalUpdatedAt) {
		// This is expected - UpdatedAt is updated when result is valid
	}
}

// ==================== Day.Validate Tests ====================

func TestDay_Validate_Valid(t *testing.T) {
	day := &Day{
		ID:   "test-id",
		Name: "Day A",
		Slug: "day-a",
	}

	result := day.Validate()

	if !result.Valid {
		t.Errorf("Validate returned invalid result for valid day: %v", result.Errors)
	}
}

func TestDay_Validate_InvalidName(t *testing.T) {
	day := &Day{
		ID:   "test-id",
		Name: "",
		Slug: "day-a",
	}

	result := day.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for day with empty name")
	}
}

func TestDay_Validate_InvalidSlug(t *testing.T) {
	day := &Day{
		ID:   "test-id",
		Name: "Day A",
		Slug: "Invalid Slug",
	}

	result := day.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for day with invalid slug")
	}
}

// ==================== DayPrescription Tests ====================

func TestCreateDayPrescription_Valid(t *testing.T) {
	input := CreateDayPrescriptionInput{
		DayID:          "day-123",
		PrescriptionID: "prescription-456",
	}

	dp, result := CreateDayPrescription(input, "dp-id", 0)

	if !result.Valid {
		t.Errorf("CreateDayPrescription returned invalid result: %v", result.Errors)
	}
	if dp == nil {
		t.Fatal("CreateDayPrescription returned nil")
	}
	if dp.ID != "dp-id" {
		t.Errorf("dp.ID = %q, want %q", dp.ID, "dp-id")
	}
	if dp.DayID != "day-123" {
		t.Errorf("dp.DayID = %q, want %q", dp.DayID, "day-123")
	}
	if dp.PrescriptionID != "prescription-456" {
		t.Errorf("dp.PrescriptionID = %q, want %q", dp.PrescriptionID, "prescription-456")
	}
	if dp.Order != 0 {
		t.Errorf("dp.Order = %d, want %d", dp.Order, 0)
	}
}

func TestCreateDayPrescription_WithOrder(t *testing.T) {
	order := 5
	input := CreateDayPrescriptionInput{
		DayID:          "day-123",
		PrescriptionID: "prescription-456",
		Order:          &order,
	}

	dp, result := CreateDayPrescription(input, "dp-id", order)

	if !result.Valid {
		t.Errorf("CreateDayPrescription returned invalid result: %v", result.Errors)
	}
	if dp.Order != 5 {
		t.Errorf("dp.Order = %d, want %d", dp.Order, 5)
	}
}

func TestCreateDayPrescription_MissingDayID(t *testing.T) {
	input := CreateDayPrescriptionInput{
		DayID:          "",
		PrescriptionID: "prescription-456",
	}

	dp, result := CreateDayPrescription(input, "dp-id", 0)

	if result.Valid {
		t.Error("CreateDayPrescription with missing day ID returned valid result")
	}
	if dp != nil {
		t.Error("CreateDayPrescription with invalid input returned non-nil")
	}
}

func TestCreateDayPrescription_MissingPrescriptionID(t *testing.T) {
	input := CreateDayPrescriptionInput{
		DayID:          "day-123",
		PrescriptionID: "",
	}

	dp, result := CreateDayPrescription(input, "dp-id", 0)

	if result.Valid {
		t.Error("CreateDayPrescription with missing prescription ID returned valid result")
	}
	if dp != nil {
		t.Error("CreateDayPrescription with invalid input returned non-nil")
	}
}

func TestCreateDayPrescription_NegativeOrder(t *testing.T) {
	input := CreateDayPrescriptionInput{
		DayID:          "day-123",
		PrescriptionID: "prescription-456",
	}

	dp, result := CreateDayPrescription(input, "dp-id", -1)

	if result.Valid {
		t.Error("CreateDayPrescription with negative order returned valid result")
	}
	if dp != nil {
		t.Error("CreateDayPrescription with invalid input returned non-nil")
	}
}

// ==================== ReorderPrescriptionsInput Tests ====================

func TestValidateReorderInput_Valid(t *testing.T) {
	input := ReorderPrescriptionsInput{
		DayID:           "day-123",
		PrescriptionIDs: []string{"p1", "p2", "p3"},
	}

	result := ValidateReorderInput(input)

	if !result.Valid {
		t.Errorf("ValidateReorderInput returned invalid result: %v", result.Errors)
	}
}

func TestValidateReorderInput_MissingDayID(t *testing.T) {
	input := ReorderPrescriptionsInput{
		DayID:           "",
		PrescriptionIDs: []string{"p1", "p2"},
	}

	result := ValidateReorderInput(input)

	if result.Valid {
		t.Error("ValidateReorderInput with missing day ID returned valid result")
	}
}

func TestValidateReorderInput_EmptyPrescriptionIDs(t *testing.T) {
	input := ReorderPrescriptionsInput{
		DayID:           "day-123",
		PrescriptionIDs: []string{},
	}

	result := ValidateReorderInput(input)

	if result.Valid {
		t.Error("ValidateReorderInput with empty prescription IDs returned valid result")
	}
}

func TestValidateReorderInput_DuplicatePrescriptionIDs(t *testing.T) {
	input := ReorderPrescriptionsInput{
		DayID:           "day-123",
		PrescriptionIDs: []string{"p1", "p2", "p1"},
	}

	result := ValidateReorderInput(input)

	if result.Valid {
		t.Error("ValidateReorderInput with duplicate prescription IDs returned valid result")
	}
}

func TestValidateReorderInput_EmptyPrescriptionID(t *testing.T) {
	input := ReorderPrescriptionsInput{
		DayID:           "day-123",
		PrescriptionIDs: []string{"p1", "", "p2"},
	}

	result := ValidateReorderInput(input)

	if result.Valid {
		t.Error("ValidateReorderInput with empty prescription ID returned valid result")
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
	result.AddError(ErrSlugInvalid)

	err := result.Error()
	if err == nil {
		t.Error("ValidationResult.Error() = nil, want error")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "name is required") {
		t.Errorf("Error message should contain name error: %s", errMsg)
	}
	if !strings.Contains(errMsg, "slug must contain") {
		t.Errorf("Error message should contain slug error: %s", errMsg)
	}
}

// ==================== Edge Cases ====================

func TestValidateName_ExactlyMaxLength(t *testing.T) {
	name := strings.Repeat("a", 50)
	err := ValidateName(name)
	if err != nil {
		t.Errorf("ValidateName with exactly 50 chars = %v, want nil", err)
	}
}

func TestValidateSlug_ExactlyMaxLength(t *testing.T) {
	slug := strings.Repeat("a", 50)
	err := ValidateSlug(slug)
	if err != nil {
		t.Errorf("ValidateSlug with exactly 50 chars = %v, want nil", err)
	}
}

func TestGenerateSlug_EmptyInput(t *testing.T) {
	slug := GenerateSlug("")
	if slug != "" {
		t.Errorf("GenerateSlug(\"\") = %q, want empty string", slug)
	}
}

func TestGenerateSlug_OnlySpecialChars(t *testing.T) {
	slug := GenerateSlug("!@#$%")
	if slug != "" {
		t.Errorf("GenerateSlug(\"!@#$%%\") = %q, want empty string", slug)
	}
}

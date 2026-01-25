package program

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
		{"5/3/1 BBB"},
		{"Starting Strength"},
		{"Bill Starr 5x5"},
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

// ==================== Slug Validation Tests ====================

func TestValidateSlug_Valid(t *testing.T) {
	tests := []struct {
		name string
		slug string
	}{
		{"simple lowercase", "starting-strength"},
		{"with numbers", "531-bbb"},
		{"single char", "a"},
		{"numbers only", "123"},
		{"max length", strings.Repeat("a", MaxSlugLength)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSlug(tt.slug)
			if err != nil {
				t.Errorf("ValidateSlug(%q) = %v, want nil", tt.slug, err)
			}
		})
	}
}

func TestValidateSlug_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		slug        string
		expectedErr error
	}{
		{"empty string", "", ErrSlugRequired},
		{"only spaces", "   ", ErrSlugRequired},
		{"uppercase letters", "Starting-Strength", ErrSlugInvalid},
		{"spaces", "starting strength", ErrSlugInvalid},
		{"underscores", "starting_strength", ErrSlugInvalid},
		{"special chars", "5/3/1-bbb", ErrSlugInvalid},
		{"too long", strings.Repeat("a", MaxSlugLength+1), ErrSlugTooLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSlug(tt.slug)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateSlug(%q) = %v, want %v", tt.slug, err, tt.expectedErr)
			}
		})
	}
}

// ==================== CycleID Validation Tests ====================

func TestValidateCycleID_Valid(t *testing.T) {
	tests := []struct {
		name    string
		cycleID string
	}{
		{"uuid", "550e8400-e29b-41d4-a716-446655440000"},
		{"simple id", "cycle-1"},
		{"short id", "c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCycleID(tt.cycleID)
			if err != nil {
				t.Errorf("ValidateCycleID(%q) = %v, want nil", tt.cycleID, err)
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

// ==================== DefaultRounding Validation Tests ====================

func TestValidateDefaultRounding_Valid(t *testing.T) {
	tests := []struct {
		name     string
		rounding *float64
	}{
		{"nil", nil},
		{"2.5 lbs", ptrFloat64(2.5)},
		{"5 lbs", ptrFloat64(5.0)},
		{"1 kg", ptrFloat64(1.0)},
		{"small positive", ptrFloat64(0.1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDefaultRounding(tt.rounding)
			if err != nil {
				t.Errorf("ValidateDefaultRounding(%v) = %v, want nil", tt.rounding, err)
			}
		})
	}
}

func TestValidateDefaultRounding_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		rounding    *float64
		expectedErr error
	}{
		{"zero", ptrFloat64(0), ErrDefaultRoundingInvalid},
		{"negative", ptrFloat64(-2.5), ErrDefaultRoundingInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDefaultRounding(tt.rounding)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateDefaultRounding(%v) = %v, want %v", tt.rounding, err, tt.expectedErr)
			}
		})
	}
}

// ==================== CreateProgram Tests ====================

func TestCreateProgram_ValidInput(t *testing.T) {
	input := CreateProgramInput{
		Name:    "5/3/1 BBB",
		Slug:    "531-bbb",
		CycleID: "cycle-id-1",
	}

	prog, result := CreateProgram(input, "test-id")

	if !result.Valid {
		t.Errorf("CreateProgram returned invalid result: %v", result.Errors)
	}
	if prog == nil {
		t.Fatal("CreateProgram returned nil program")
	}
	if prog.ID != "test-id" {
		t.Errorf("program.ID = %q, want %q", prog.ID, "test-id")
	}
	if prog.Name != "5/3/1 BBB" {
		t.Errorf("program.Name = %q, want %q", prog.Name, "5/3/1 BBB")
	}
	if prog.Slug != "531-bbb" {
		t.Errorf("program.Slug = %q, want %q", prog.Slug, "531-bbb")
	}
	if prog.CycleID != "cycle-id-1" {
		t.Errorf("program.CycleID = %q, want %q", prog.CycleID, "cycle-id-1")
	}
}

func TestCreateProgram_FullInput(t *testing.T) {
	input := CreateProgramInput{
		Name:            "5/3/1 BBB",
		Slug:            "531-bbb",
		Description:     ptrString("Wendler's 5/3/1 Boring But Big variant"),
		CycleID:         "cycle-id-1",
		WeeklyLookupID:  ptrString("weekly-lookup-id"),
		DailyLookupID:   ptrString("daily-lookup-id"),
		DefaultRounding: ptrFloat64(2.5),
	}

	prog, result := CreateProgram(input, "test-id")

	if !result.Valid {
		t.Errorf("CreateProgram returned invalid result: %v", result.Errors)
	}
	if prog == nil {
		t.Fatal("CreateProgram returned nil program")
	}
	if prog.Description == nil || *prog.Description != "Wendler's 5/3/1 Boring But Big variant" {
		t.Errorf("program.Description = %v, want %q", prog.Description, "Wendler's 5/3/1 Boring But Big variant")
	}
	if prog.WeeklyLookupID == nil || *prog.WeeklyLookupID != "weekly-lookup-id" {
		t.Errorf("program.WeeklyLookupID = %v, want %q", prog.WeeklyLookupID, "weekly-lookup-id")
	}
	if prog.DailyLookupID == nil || *prog.DailyLookupID != "daily-lookup-id" {
		t.Errorf("program.DailyLookupID = %v, want %q", prog.DailyLookupID, "daily-lookup-id")
	}
	if prog.DefaultRounding == nil || *prog.DefaultRounding != 2.5 {
		t.Errorf("program.DefaultRounding = %v, want %v", prog.DefaultRounding, 2.5)
	}
}

func TestCreateProgram_TrimsWhitespace(t *testing.T) {
	// Note: Slug validation happens before trimming, so slug cannot have spaces
	// We test trimming for name, description, and cycle_id
	input := CreateProgramInput{
		Name:        "  5/3/1 BBB  ",
		Slug:        "531-bbb",
		Description: ptrString("  Description  "),
		CycleID:     "  cycle-id-1  ",
	}

	prog, result := CreateProgram(input, "test-id")

	if !result.Valid {
		t.Errorf("CreateProgram returned invalid result: %v", result.Errors)
	}
	if prog.Name != "5/3/1 BBB" {
		t.Errorf("program.Name = %q, want %q (trimmed)", prog.Name, "5/3/1 BBB")
	}
	if prog.Slug != "531-bbb" {
		t.Errorf("program.Slug = %q, want %q", prog.Slug, "531-bbb")
	}
	if prog.Description == nil || *prog.Description != "Description" {
		t.Errorf("program.Description = %v, want %q (trimmed)", prog.Description, "Description")
	}
	if prog.CycleID != "cycle-id-1" {
		t.Errorf("program.CycleID = %q, want %q (trimmed)", prog.CycleID, "cycle-id-1")
	}
}

func TestCreateProgram_EmptyName(t *testing.T) {
	input := CreateProgramInput{
		Name:    "",
		Slug:    "531-bbb",
		CycleID: "cycle-id-1",
	}

	prog, result := CreateProgram(input, "test-id")

	if result.Valid {
		t.Error("CreateProgram with empty name returned valid result")
	}
	if prog != nil {
		t.Error("CreateProgram with invalid input returned non-nil program")
	}
}

func TestCreateProgram_EmptySlug(t *testing.T) {
	input := CreateProgramInput{
		Name:    "5/3/1 BBB",
		Slug:    "",
		CycleID: "cycle-id-1",
	}

	prog, result := CreateProgram(input, "test-id")

	if result.Valid {
		t.Error("CreateProgram with empty slug returned valid result")
	}
	if prog != nil {
		t.Error("CreateProgram with invalid input returned non-nil program")
	}
}

func TestCreateProgram_InvalidSlug(t *testing.T) {
	input := CreateProgramInput{
		Name:    "5/3/1 BBB",
		Slug:    "Invalid Slug",
		CycleID: "cycle-id-1",
	}

	prog, result := CreateProgram(input, "test-id")

	if result.Valid {
		t.Error("CreateProgram with invalid slug returned valid result")
	}
	if prog != nil {
		t.Error("CreateProgram with invalid input returned non-nil program")
	}
}

func TestCreateProgram_EmptyCycleID(t *testing.T) {
	input := CreateProgramInput{
		Name:    "5/3/1 BBB",
		Slug:    "531-bbb",
		CycleID: "",
	}

	prog, result := CreateProgram(input, "test-id")

	if result.Valid {
		t.Error("CreateProgram with empty cycle_id returned valid result")
	}
	if prog != nil {
		t.Error("CreateProgram with invalid input returned non-nil program")
	}
}

func TestCreateProgram_InvalidDefaultRounding(t *testing.T) {
	input := CreateProgramInput{
		Name:            "5/3/1 BBB",
		Slug:            "531-bbb",
		CycleID:         "cycle-id-1",
		DefaultRounding: ptrFloat64(0),
	}

	prog, result := CreateProgram(input, "test-id")

	if result.Valid {
		t.Error("CreateProgram with invalid default_rounding returned valid result")
	}
	if prog != nil {
		t.Error("CreateProgram with invalid input returned non-nil program")
	}
}

func TestCreateProgram_MultipleErrors(t *testing.T) {
	input := CreateProgramInput{
		Name:            "",          // Invalid
		Slug:            "Invalid!",  // Invalid
		CycleID:         "",          // Invalid
		DefaultRounding: ptrFloat64(-1), // Invalid
	}

	prog, result := CreateProgram(input, "test-id")

	if result.Valid {
		t.Error("CreateProgram with multiple errors returned valid result")
	}
	if prog != nil {
		t.Error("CreateProgram with invalid input returned non-nil program")
	}
	if len(result.Errors) < 4 {
		t.Errorf("Expected at least 4 errors, got %d", len(result.Errors))
	}
}

func TestCreateProgram_EmptyDescriptionBecomesNil(t *testing.T) {
	input := CreateProgramInput{
		Name:        "5/3/1 BBB",
		Slug:        "531-bbb",
		Description: ptrString("   "), // Empty after trim
		CycleID:     "cycle-id-1",
	}

	prog, result := CreateProgram(input, "test-id")

	if !result.Valid {
		t.Errorf("CreateProgram returned invalid result: %v", result.Errors)
	}
	if prog.Description != nil {
		t.Errorf("program.Description = %v, want nil (empty after trim)", prog.Description)
	}
}

// ==================== UpdateProgram Tests ====================

func TestUpdateProgram_UpdateName(t *testing.T) {
	prog := &Program{ID: "test-id", Name: "Old Name", Slug: "old-slug", CycleID: "cycle-1"}

	newName := "New Name"
	input := UpdateProgramInput{Name: &newName}

	result := UpdateProgram(prog, input)

	if !result.Valid {
		t.Errorf("UpdateProgram returned invalid result: %v", result.Errors)
	}
	if prog.Name != "New Name" {
		t.Errorf("program.Name = %q, want %q", prog.Name, "New Name")
	}
}

func TestUpdateProgram_UpdateSlug(t *testing.T) {
	prog := &Program{ID: "test-id", Name: "Test Program", Slug: "old-slug", CycleID: "cycle-1"}

	newSlug := "new-slug"
	input := UpdateProgramInput{Slug: &newSlug}

	result := UpdateProgram(prog, input)

	if !result.Valid {
		t.Errorf("UpdateProgram returned invalid result: %v", result.Errors)
	}
	if prog.Slug != "new-slug" {
		t.Errorf("program.Slug = %q, want %q", prog.Slug, "new-slug")
	}
}

func TestUpdateProgram_UpdateCycleID(t *testing.T) {
	prog := &Program{ID: "test-id", Name: "Test Program", Slug: "test-slug", CycleID: "old-cycle"}

	newCycleID := "new-cycle"
	input := UpdateProgramInput{CycleID: &newCycleID}

	result := UpdateProgram(prog, input)

	if !result.Valid {
		t.Errorf("UpdateProgram returned invalid result: %v", result.Errors)
	}
	if prog.CycleID != "new-cycle" {
		t.Errorf("program.CycleID = %q, want %q", prog.CycleID, "new-cycle")
	}
}

func TestUpdateProgram_UpdateDescription(t *testing.T) {
	prog := &Program{ID: "test-id", Name: "Test Program", Slug: "test-slug", CycleID: "cycle-1"}

	newDesc := "New Description"
	newDescPtr := &newDesc
	input := UpdateProgramInput{Description: &newDescPtr}

	result := UpdateProgram(prog, input)

	if !result.Valid {
		t.Errorf("UpdateProgram returned invalid result: %v", result.Errors)
	}
	if prog.Description == nil || *prog.Description != "New Description" {
		t.Errorf("program.Description = %v, want %q", prog.Description, "New Description")
	}
}

func TestUpdateProgram_ClearDescription(t *testing.T) {
	desc := "Old Description"
	prog := &Program{ID: "test-id", Name: "Test Program", Slug: "test-slug", CycleID: "cycle-1", Description: &desc}

	var nilDesc *string = nil
	input := UpdateProgramInput{Description: &nilDesc}

	result := UpdateProgram(prog, input)

	if !result.Valid {
		t.Errorf("UpdateProgram returned invalid result: %v", result.Errors)
	}
	if prog.Description != nil {
		t.Errorf("program.Description = %v, want nil", prog.Description)
	}
}

func TestUpdateProgram_UpdateDefaultRounding(t *testing.T) {
	prog := &Program{ID: "test-id", Name: "Test Program", Slug: "test-slug", CycleID: "cycle-1"}

	newRounding := 2.5
	newRoundingPtr := &newRounding
	input := UpdateProgramInput{DefaultRounding: &newRoundingPtr}

	result := UpdateProgram(prog, input)

	if !result.Valid {
		t.Errorf("UpdateProgram returned invalid result: %v", result.Errors)
	}
	if prog.DefaultRounding == nil || *prog.DefaultRounding != 2.5 {
		t.Errorf("program.DefaultRounding = %v, want %v", prog.DefaultRounding, 2.5)
	}
}

func TestUpdateProgram_ClearDefaultRounding(t *testing.T) {
	rounding := 2.5
	prog := &Program{ID: "test-id", Name: "Test Program", Slug: "test-slug", CycleID: "cycle-1", DefaultRounding: &rounding}

	var nilRounding *float64 = nil
	input := UpdateProgramInput{DefaultRounding: &nilRounding}

	result := UpdateProgram(prog, input)

	if !result.Valid {
		t.Errorf("UpdateProgram returned invalid result: %v", result.Errors)
	}
	if prog.DefaultRounding != nil {
		t.Errorf("program.DefaultRounding = %v, want nil", prog.DefaultRounding)
	}
}

func TestUpdateProgram_InvalidName(t *testing.T) {
	prog := &Program{ID: "test-id", Name: "Old Name", Slug: "test-slug", CycleID: "cycle-1"}
	originalName := prog.Name

	emptyName := ""
	input := UpdateProgramInput{Name: &emptyName}

	result := UpdateProgram(prog, input)

	if result.Valid {
		t.Error("UpdateProgram with invalid name returned valid result")
	}
	if prog.Name != originalName {
		t.Errorf("program.Name was changed despite validation failure")
	}
}

func TestUpdateProgram_InvalidSlug(t *testing.T) {
	prog := &Program{ID: "test-id", Name: "Test Program", Slug: "old-slug", CycleID: "cycle-1"}
	originalSlug := prog.Slug

	invalidSlug := "Invalid Slug"
	input := UpdateProgramInput{Slug: &invalidSlug}

	result := UpdateProgram(prog, input)

	if result.Valid {
		t.Error("UpdateProgram with invalid slug returned valid result")
	}
	if prog.Slug != originalSlug {
		t.Errorf("program.Slug was changed despite validation failure")
	}
}

func TestUpdateProgram_InvalidCycleID(t *testing.T) {
	prog := &Program{ID: "test-id", Name: "Test Program", Slug: "test-slug", CycleID: "old-cycle"}
	originalCycleID := prog.CycleID

	emptyCycleID := ""
	input := UpdateProgramInput{CycleID: &emptyCycleID}

	result := UpdateProgram(prog, input)

	if result.Valid {
		t.Error("UpdateProgram with invalid cycle_id returned valid result")
	}
	if prog.CycleID != originalCycleID {
		t.Errorf("program.CycleID was changed despite validation failure")
	}
}

func TestUpdateProgram_InvalidDefaultRounding(t *testing.T) {
	prog := &Program{ID: "test-id", Name: "Test Program", Slug: "test-slug", CycleID: "cycle-1"}

	zeroRounding := 0.0
	zeroRoundingPtr := &zeroRounding
	input := UpdateProgramInput{DefaultRounding: &zeroRoundingPtr}

	result := UpdateProgram(prog, input)

	if result.Valid {
		t.Error("UpdateProgram with invalid default_rounding returned valid result")
	}
}

func TestUpdateProgram_NoChanges(t *testing.T) {
	prog := &Program{ID: "test-id", Name: "Test Program", Slug: "test-slug", CycleID: "cycle-1"}

	input := UpdateProgramInput{} // No changes

	result := UpdateProgram(prog, input)

	if !result.Valid {
		t.Errorf("UpdateProgram with no changes returned invalid result: %v", result.Errors)
	}
}

func TestUpdateProgram_TrimsWhitespace(t *testing.T) {
	prog := &Program{ID: "test-id", Name: "Old Name", Slug: "old-slug", CycleID: "old-cycle"}

	newName := "  New Name  "
	input := UpdateProgramInput{Name: &newName}

	result := UpdateProgram(prog, input)

	if !result.Valid {
		t.Errorf("UpdateProgram returned invalid result: %v", result.Errors)
	}
	if prog.Name != "New Name" {
		t.Errorf("program.Name = %q, want %q (trimmed)", prog.Name, "New Name")
	}
}

// ==================== Program.Validate Tests ====================

func TestProgram_Validate_Valid(t *testing.T) {
	prog := &Program{
		ID:      "test-id",
		Name:    "Test Program",
		Slug:    "test-program",
		CycleID: "cycle-1",
	}

	result := prog.Validate()

	if !result.Valid {
		t.Errorf("Validate returned invalid result for valid program: %v", result.Errors)
	}
}

func TestProgram_Validate_InvalidName(t *testing.T) {
	prog := &Program{
		ID:      "test-id",
		Name:    "",
		Slug:    "test-program",
		CycleID: "cycle-1",
	}

	result := prog.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for program with invalid name")
	}
}

func TestProgram_Validate_InvalidSlug(t *testing.T) {
	prog := &Program{
		ID:      "test-id",
		Name:    "Test Program",
		Slug:    "Invalid Slug",
		CycleID: "cycle-1",
	}

	result := prog.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for program with invalid slug")
	}
}

func TestProgram_Validate_InvalidCycleID(t *testing.T) {
	prog := &Program{
		ID:      "test-id",
		Name:    "Test Program",
		Slug:    "test-program",
		CycleID: "",
	}

	result := prog.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for program with invalid cycle_id")
	}
}

func TestProgram_Validate_InvalidDefaultRounding(t *testing.T) {
	rounding := -1.0
	prog := &Program{
		ID:              "test-id",
		Name:            "Test Program",
		Slug:            "test-program",
		CycleID:         "cycle-1",
		DefaultRounding: &rounding,
	}

	result := prog.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for program with invalid default_rounding")
	}
}

func TestProgram_Validate_MultipleErrors(t *testing.T) {
	rounding := -1.0
	prog := &Program{
		ID:              "test-id",
		Name:            "",
		Slug:            "Invalid!",
		CycleID:         "",
		DefaultRounding: &rounding,
	}

	result := prog.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for program with multiple errors")
	}
	if len(result.Errors) < 4 {
		t.Errorf("Expected at least 4 errors, got %d", len(result.Errors))
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
	if errMsg == "" {
		t.Error("Error message should not be empty")
	}
}

// ==================== Filter Validation Tests ====================

func TestValidateDifficulty_Valid(t *testing.T) {
	tests := []struct {
		name       string
		difficulty string
	}{
		{"beginner", "beginner"},
		{"intermediate", "intermediate"},
		{"advanced", "advanced"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDifficulty(tt.difficulty)
			if err != nil {
				t.Errorf("ValidateDifficulty(%q) = %v, want nil", tt.difficulty, err)
			}
		})
	}
}

func TestValidateDifficulty_Invalid(t *testing.T) {
	tests := []struct {
		name       string
		difficulty string
	}{
		{"empty", ""},
		{"capitalized", "Beginner"},
		{"unknown", "expert"},
		{"typo", "begginer"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDifficulty(tt.difficulty)
			if !errors.Is(err, ErrInvalidDifficulty) {
				t.Errorf("ValidateDifficulty(%q) = %v, want %v", tt.difficulty, err, ErrInvalidDifficulty)
			}
		})
	}
}

func TestValidateDaysPerWeek_Valid(t *testing.T) {
	for days := 1; days <= 7; days++ {
		t.Run(string(rune('0'+days)), func(t *testing.T) {
			err := ValidateDaysPerWeek(days)
			if err != nil {
				t.Errorf("ValidateDaysPerWeek(%d) = %v, want nil", days, err)
			}
		})
	}
}

func TestValidateDaysPerWeek_Invalid(t *testing.T) {
	tests := []struct {
		name string
		days int
	}{
		{"zero", 0},
		{"negative", -1},
		{"too high", 8},
		{"way too high", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDaysPerWeek(tt.days)
			if !errors.Is(err, ErrInvalidDaysPerWeek) {
				t.Errorf("ValidateDaysPerWeek(%d) = %v, want %v", tt.days, err, ErrInvalidDaysPerWeek)
			}
		})
	}
}

func TestValidateFocus_Valid(t *testing.T) {
	tests := []struct {
		name  string
		focus string
	}{
		{"strength", "strength"},
		{"hypertrophy", "hypertrophy"},
		{"peaking", "peaking"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFocus(tt.focus)
			if err != nil {
				t.Errorf("ValidateFocus(%q) = %v, want nil", tt.focus, err)
			}
		})
	}
}

func TestValidateFocus_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		focus string
	}{
		{"empty", ""},
		{"capitalized", "Strength"},
		{"unknown", "cardio"},
		{"typo", "strenght"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFocus(tt.focus)
			if !errors.Is(err, ErrInvalidFocus) {
				t.Errorf("ValidateFocus(%q) = %v, want %v", tt.focus, err, ErrInvalidFocus)
			}
		})
	}
}

func TestFilterOptions_Validate_Valid(t *testing.T) {
	tests := []struct {
		name    string
		filters FilterOptions
	}{
		{"empty filters", FilterOptions{}},
		{"difficulty only", FilterOptions{Difficulty: ptrString("beginner")}},
		{"days only", FilterOptions{DaysPerWeek: ptrInt(3)}},
		{"focus only", FilterOptions{Focus: ptrString("strength")}},
		{"has_amrap true", FilterOptions{HasAmrap: ptrBool(true)}},
		{"has_amrap false", FilterOptions{HasAmrap: ptrBool(false)}},
		{"all filters", FilterOptions{
			Difficulty:  ptrString("advanced"),
			DaysPerWeek: ptrInt(4),
			Focus:       ptrString("hypertrophy"),
			HasAmrap:    ptrBool(true),
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filters.Validate()
			if !result.Valid {
				t.Errorf("FilterOptions.Validate() = invalid with errors: %v", result.Errors)
			}
		})
	}
}

func TestFilterOptions_Validate_Invalid(t *testing.T) {
	tests := []struct {
		name          string
		filters       FilterOptions
		expectedError error
	}{
		{"invalid difficulty", FilterOptions{Difficulty: ptrString("expert")}, ErrInvalidDifficulty},
		{"invalid days too low", FilterOptions{DaysPerWeek: ptrInt(0)}, ErrInvalidDaysPerWeek},
		{"invalid days too high", FilterOptions{DaysPerWeek: ptrInt(8)}, ErrInvalidDaysPerWeek},
		{"invalid focus", FilterOptions{Focus: ptrString("cardio")}, ErrInvalidFocus},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filters.Validate()
			if result.Valid {
				t.Errorf("FilterOptions.Validate() = valid, want invalid")
			}
			// Check that the expected error is in the errors list
			found := false
			for _, err := range result.Errors {
				if errors.Is(err, tt.expectedError) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected error %v not found in errors: %v", tt.expectedError, result.Errors)
			}
		})
	}
}

func TestFilterOptions_Validate_MultipleErrors(t *testing.T) {
	filters := FilterOptions{
		Difficulty:  ptrString("expert"),    // invalid
		DaysPerWeek: ptrInt(0),              // invalid
		Focus:       ptrString("cardio"),    // invalid
		HasAmrap:    ptrBool(true),          // valid (boolean always valid)
	}

	result := filters.Validate()
	if result.Valid {
		t.Error("FilterOptions.Validate() = valid, want invalid")
	}
	if len(result.Errors) != 3 {
		t.Errorf("Expected 3 errors, got %d: %v", len(result.Errors), result.Errors)
	}
}

// ==================== Helper Functions for Tests ====================

func ptrString(s string) *string {
	return &s
}

func ptrFloat64(f float64) *float64 {
	return &f
}

func ptrInt(i int) *int {
	return &i
}

func ptrBool(b bool) *bool {
	return &b
}

package dailylookup

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
		{"Bill Starr Day Intensities"},
		{"Heavy/Light/Medium"},
		{"Daily Modifiers"},
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

// ==================== IsValidIntensityLevel Tests ====================

func TestIsValidIntensityLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		expected bool
	}{
		{"HEAVY uppercase", "HEAVY", true},
		{"LIGHT uppercase", "LIGHT", true},
		{"MEDIUM uppercase", "MEDIUM", true},
		{"heavy lowercase", "heavy", true},
		{"light lowercase", "light", true},
		{"medium lowercase", "medium", true},
		{"Heavy mixed case", "Heavy", true},
		{"invalid", "INVALID", false},
		{"empty", "", false},
		{"high", "HIGH", false},
		{"low", "LOW", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidIntensityLevel(tt.level)
			if result != tt.expected {
				t.Errorf("IsValidIntensityLevel(%q) = %v, want %v", tt.level, result, tt.expected)
			}
		})
	}
}

// ==================== Entries Validation Tests ====================

func TestValidateEntries_Valid(t *testing.T) {
	tests := []struct {
		name    string
		entries []DailyLookupEntry
	}{
		{
			name: "single entry",
			entries: []DailyLookupEntry{
				{DayIdentifier: "heavy", PercentageModifier: 100},
			},
		},
		{
			name: "multiple entries",
			entries: []DailyLookupEntry{
				{DayIdentifier: "heavy", PercentageModifier: 100},
				{DayIdentifier: "light", PercentageModifier: 70},
				{DayIdentifier: "medium", PercentageModifier: 80},
			},
		},
		{
			name: "with intensity level",
			entries: []DailyLookupEntry{
				{DayIdentifier: "monday", PercentageModifier: 100, IntensityLevel: stringPtr("HEAVY")},
				{DayIdentifier: "wednesday", PercentageModifier: 70, IntensityLevel: stringPtr("LIGHT")},
			},
		},
		{
			name: "slug-style identifiers",
			entries: []DailyLookupEntry{
				{DayIdentifier: "squat-day", PercentageModifier: 100},
				{DayIdentifier: "bench-day", PercentageModifier: 95},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEntries(tt.entries)
			if err != nil {
				t.Errorf("ValidateEntries() = %v, want nil", err)
			}
		})
	}
}

func TestValidateEntries_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		entries     []DailyLookupEntry
		expectedErr error
	}{
		{
			name:        "empty entries",
			entries:     []DailyLookupEntry{},
			expectedErr: ErrEntriesRequired,
		},
		{
			name:        "nil entries",
			entries:     nil,
			expectedErr: ErrEntriesRequired,
		},
		{
			name: "empty day identifier",
			entries: []DailyLookupEntry{
				{DayIdentifier: "", PercentageModifier: 100},
			},
			expectedErr: ErrDayIdentifierRequired,
		},
		{
			name: "whitespace-only day identifier",
			entries: []DailyLookupEntry{
				{DayIdentifier: "   ", PercentageModifier: 100},
			},
			expectedErr: ErrDayIdentifierRequired,
		},
		{
			name: "duplicate day identifiers",
			entries: []DailyLookupEntry{
				{DayIdentifier: "heavy", PercentageModifier: 100},
				{DayIdentifier: "heavy", PercentageModifier: 90},
			},
			expectedErr: ErrDuplicateDayIdentifier,
		},
		{
			name: "duplicate day identifiers case insensitive",
			entries: []DailyLookupEntry{
				{DayIdentifier: "Heavy", PercentageModifier: 100},
				{DayIdentifier: "HEAVY", PercentageModifier: 90},
			},
			expectedErr: ErrDuplicateDayIdentifier,
		},
		{
			name: "invalid intensity level",
			entries: []DailyLookupEntry{
				{DayIdentifier: "monday", PercentageModifier: 100, IntensityLevel: stringPtr("INVALID")},
			},
			expectedErr: ErrIntensityLevelInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEntries(tt.entries)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateEntries() = %v, want %v", err, tt.expectedErr)
			}
		})
	}
}

// ==================== CreateDailyLookup Tests ====================

func TestCreateDailyLookup_ValidInput(t *testing.T) {
	input := CreateDailyLookupInput{
		Name: "Bill Starr Intensities",
		Entries: []DailyLookupEntry{
			{DayIdentifier: "heavy", PercentageModifier: 100, IntensityLevel: stringPtr("HEAVY")},
			{DayIdentifier: "light", PercentageModifier: 70, IntensityLevel: stringPtr("LIGHT")},
			{DayIdentifier: "medium", PercentageModifier: 80, IntensityLevel: stringPtr("MEDIUM")},
		},
	}

	lookup, result := CreateDailyLookup(input, "test-id")

	if !result.Valid {
		t.Errorf("CreateDailyLookup returned invalid result: %v", result.Errors)
	}
	if lookup == nil {
		t.Fatal("CreateDailyLookup returned nil lookup")
	}
	if lookup.ID != "test-id" {
		t.Errorf("lookup.ID = %q, want %q", lookup.ID, "test-id")
	}
	if lookup.Name != "Bill Starr Intensities" {
		t.Errorf("lookup.Name = %q, want %q", lookup.Name, "Bill Starr Intensities")
	}
	if len(lookup.Entries) != 3 {
		t.Errorf("len(lookup.Entries) = %d, want %d", len(lookup.Entries), 3)
	}
}

func TestCreateDailyLookup_WithProgramID(t *testing.T) {
	programID := "program-123"
	input := CreateDailyLookupInput{
		Name: "Custom Intensities",
		Entries: []DailyLookupEntry{
			{DayIdentifier: "day-a", PercentageModifier: 100},
		},
		ProgramID: &programID,
	}

	lookup, result := CreateDailyLookup(input, "test-id")

	if !result.Valid {
		t.Errorf("CreateDailyLookup returned invalid result: %v", result.Errors)
	}
	if lookup.ProgramID == nil || *lookup.ProgramID != programID {
		t.Errorf("lookup.ProgramID = %v, want %q", lookup.ProgramID, programID)
	}
}

func TestCreateDailyLookup_TrimsWhitespace(t *testing.T) {
	input := CreateDailyLookupInput{
		Name: "  Daily Intensities  ",
		Entries: []DailyLookupEntry{
			{DayIdentifier: "heavy", PercentageModifier: 100},
		},
	}

	lookup, result := CreateDailyLookup(input, "test-id")

	if !result.Valid {
		t.Errorf("CreateDailyLookup returned invalid result: %v", result.Errors)
	}
	if lookup.Name != "Daily Intensities" {
		t.Errorf("lookup.Name = %q, want %q (trimmed)", lookup.Name, "Daily Intensities")
	}
}

func TestCreateDailyLookup_EmptyName(t *testing.T) {
	input := CreateDailyLookupInput{
		Name: "",
		Entries: []DailyLookupEntry{
			{DayIdentifier: "heavy", PercentageModifier: 100},
		},
	}

	lookup, result := CreateDailyLookup(input, "test-id")

	if result.Valid {
		t.Error("CreateDailyLookup with empty name returned valid result")
	}
	if lookup != nil {
		t.Error("CreateDailyLookup with invalid input returned non-nil lookup")
	}
}

func TestCreateDailyLookup_EmptyEntries(t *testing.T) {
	input := CreateDailyLookupInput{
		Name:    "Test Lookup",
		Entries: []DailyLookupEntry{},
	}

	lookup, result := CreateDailyLookup(input, "test-id")

	if result.Valid {
		t.Error("CreateDailyLookup with empty entries returned valid result")
	}
	if lookup != nil {
		t.Error("CreateDailyLookup with invalid input returned non-nil lookup")
	}
}

func TestCreateDailyLookup_MultipleErrors(t *testing.T) {
	input := CreateDailyLookupInput{
		Name:    "",                      // Invalid
		Entries: []DailyLookupEntry{},    // Invalid
	}

	lookup, result := CreateDailyLookup(input, "test-id")

	if result.Valid {
		t.Error("CreateDailyLookup with multiple errors returned valid result")
	}
	if lookup != nil {
		t.Error("CreateDailyLookup with invalid input returned non-nil lookup")
	}
	if len(result.Errors) < 2 {
		t.Errorf("Expected at least 2 errors, got %d", len(result.Errors))
	}
}

// ==================== UpdateDailyLookup Tests ====================

func TestUpdateDailyLookup_UpdateName(t *testing.T) {
	lookup := &DailyLookup{
		ID:   "test-id",
		Name: "Old Name",
		Entries: []DailyLookupEntry{
			{DayIdentifier: "heavy", PercentageModifier: 100},
		},
	}

	newName := "New Name"
	input := UpdateDailyLookupInput{Name: &newName}

	result := UpdateDailyLookup(lookup, input)

	if !result.Valid {
		t.Errorf("UpdateDailyLookup returned invalid result: %v", result.Errors)
	}
	if lookup.Name != "New Name" {
		t.Errorf("lookup.Name = %q, want %q", lookup.Name, "New Name")
	}
}

func TestUpdateDailyLookup_UpdateEntries(t *testing.T) {
	lookup := &DailyLookup{
		ID:   "test-id",
		Name: "Test Lookup",
		Entries: []DailyLookupEntry{
			{DayIdentifier: "heavy", PercentageModifier: 100},
		},
	}

	newEntries := []DailyLookupEntry{
		{DayIdentifier: "heavy", PercentageModifier: 100},
		{DayIdentifier: "light", PercentageModifier: 70},
		{DayIdentifier: "medium", PercentageModifier: 80},
	}
	input := UpdateDailyLookupInput{Entries: &newEntries}

	result := UpdateDailyLookup(lookup, input)

	if !result.Valid {
		t.Errorf("UpdateDailyLookup returned invalid result: %v", result.Errors)
	}
	if len(lookup.Entries) != 3 {
		t.Errorf("len(lookup.Entries) = %d, want %d", len(lookup.Entries), 3)
	}
}

func TestUpdateDailyLookup_UpdateProgramID(t *testing.T) {
	lookup := &DailyLookup{
		ID:   "test-id",
		Name: "Test Lookup",
		Entries: []DailyLookupEntry{
			{DayIdentifier: "heavy", PercentageModifier: 100},
		},
		ProgramID: nil,
	}

	newProgramID := "program-456"
	newProgramIDPtr := &newProgramID
	input := UpdateDailyLookupInput{ProgramID: &newProgramIDPtr}

	result := UpdateDailyLookup(lookup, input)

	if !result.Valid {
		t.Errorf("UpdateDailyLookup returned invalid result: %v", result.Errors)
	}
	if lookup.ProgramID == nil || *lookup.ProgramID != newProgramID {
		t.Errorf("lookup.ProgramID = %v, want %q", lookup.ProgramID, newProgramID)
	}
}

func TestUpdateDailyLookup_ClearProgramID(t *testing.T) {
	programID := "program-123"
	lookup := &DailyLookup{
		ID:   "test-id",
		Name: "Test Lookup",
		Entries: []DailyLookupEntry{
			{DayIdentifier: "heavy", PercentageModifier: 100},
		},
		ProgramID: &programID,
	}

	var nilProgramID *string = nil
	input := UpdateDailyLookupInput{ProgramID: &nilProgramID}

	result := UpdateDailyLookup(lookup, input)

	if !result.Valid {
		t.Errorf("UpdateDailyLookup returned invalid result: %v", result.Errors)
	}
	if lookup.ProgramID != nil {
		t.Errorf("lookup.ProgramID = %v, want nil", lookup.ProgramID)
	}
}

func TestUpdateDailyLookup_InvalidName(t *testing.T) {
	lookup := &DailyLookup{
		ID:   "test-id",
		Name: "Old Name",
		Entries: []DailyLookupEntry{
			{DayIdentifier: "heavy", PercentageModifier: 100},
		},
	}
	originalName := lookup.Name

	emptyName := ""
	input := UpdateDailyLookupInput{Name: &emptyName}

	result := UpdateDailyLookup(lookup, input)

	if result.Valid {
		t.Error("UpdateDailyLookup with invalid name returned valid result")
	}
	if lookup.Name != originalName {
		t.Errorf("lookup.Name was changed despite validation failure")
	}
}

func TestUpdateDailyLookup_InvalidEntries(t *testing.T) {
	lookup := &DailyLookup{
		ID:   "test-id",
		Name: "Test Lookup",
		Entries: []DailyLookupEntry{
			{DayIdentifier: "heavy", PercentageModifier: 100},
		},
	}
	originalEntries := lookup.Entries

	emptyEntries := []DailyLookupEntry{}
	input := UpdateDailyLookupInput{Entries: &emptyEntries}

	result := UpdateDailyLookup(lookup, input)

	if result.Valid {
		t.Error("UpdateDailyLookup with invalid entries returned valid result")
	}
	if len(lookup.Entries) != len(originalEntries) {
		t.Errorf("lookup.Entries was changed despite validation failure")
	}
}

// ==================== GetByDayIdentifier Tests ====================

func TestDailyLookup_GetByDayIdentifier_Found(t *testing.T) {
	lookup := &DailyLookup{
		ID:   "test-id",
		Name: "Test Lookup",
		Entries: []DailyLookupEntry{
			{DayIdentifier: "heavy", PercentageModifier: 100, IntensityLevel: stringPtr("HEAVY")},
			{DayIdentifier: "light", PercentageModifier: 70, IntensityLevel: stringPtr("LIGHT")},
			{DayIdentifier: "medium", PercentageModifier: 80, IntensityLevel: stringPtr("MEDIUM")},
		},
	}

	entry := lookup.GetByDayIdentifier("light")

	if entry == nil {
		t.Fatal("GetByDayIdentifier returned nil for existing day")
	}
	if entry.DayIdentifier != "light" {
		t.Errorf("entry.DayIdentifier = %q, want %q", entry.DayIdentifier, "light")
	}
	if entry.PercentageModifier != 70 {
		t.Errorf("entry.PercentageModifier = %f, want %f", entry.PercentageModifier, 70.0)
	}
}

func TestDailyLookup_GetByDayIdentifier_CaseInsensitive(t *testing.T) {
	lookup := &DailyLookup{
		ID:   "test-id",
		Name: "Test Lookup",
		Entries: []DailyLookupEntry{
			{DayIdentifier: "Heavy", PercentageModifier: 100},
		},
	}

	tests := []string{"Heavy", "heavy", "HEAVY", "hEaVy"}
	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			entry := lookup.GetByDayIdentifier(test)
			if entry == nil {
				t.Fatalf("GetByDayIdentifier(%q) returned nil", test)
			}
		})
	}
}

func TestDailyLookup_GetByDayIdentifier_NotFound(t *testing.T) {
	lookup := &DailyLookup{
		ID:   "test-id",
		Name: "Test Lookup",
		Entries: []DailyLookupEntry{
			{DayIdentifier: "heavy", PercentageModifier: 100},
			{DayIdentifier: "light", PercentageModifier: 70},
		},
	}

	entry := lookup.GetByDayIdentifier("nonexistent")

	if entry != nil {
		t.Errorf("GetByDayIdentifier returned non-nil for non-existing day: %v", entry)
	}
}

// ==================== Validate Tests ====================

func TestDailyLookup_Validate_Valid(t *testing.T) {
	lookup := &DailyLookup{
		ID:   "test-id",
		Name: "Test Lookup",
		Entries: []DailyLookupEntry{
			{DayIdentifier: "heavy", PercentageModifier: 100},
		},
	}

	result := lookup.Validate()

	if !result.Valid {
		t.Errorf("Validate returned invalid result for valid lookup: %v", result.Errors)
	}
}

func TestDailyLookup_Validate_InvalidName(t *testing.T) {
	lookup := &DailyLookup{
		ID:   "test-id",
		Name: "",
		Entries: []DailyLookupEntry{
			{DayIdentifier: "heavy", PercentageModifier: 100},
		},
	}

	result := lookup.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for lookup with invalid name")
	}
}

func TestDailyLookup_Validate_InvalidEntries(t *testing.T) {
	lookup := &DailyLookup{
		ID:      "test-id",
		Name:    "Test Lookup",
		Entries: []DailyLookupEntry{},
	}

	result := lookup.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for lookup with empty entries")
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
	result.AddError(ErrEntriesRequired)

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

func stringPtr(s string) *string {
	return &s
}

package weeklylookup

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
		{"5/3/1 Percentages"},
		{"Greg Nuckols 3-Week Percentages"},
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

// ==================== Entries Validation Tests ====================

func TestValidateEntries_Valid(t *testing.T) {
	tests := []struct {
		name    string
		entries []WeeklyLookupEntry
	}{
		{
			name: "single entry",
			entries: []WeeklyLookupEntry{
				{WeekNumber: 1, Percentages: []float64{65, 75, 85}, Reps: []int{5, 5, 5}},
			},
		},
		{
			name: "multiple entries",
			entries: []WeeklyLookupEntry{
				{WeekNumber: 1, Percentages: []float64{65, 75, 85}, Reps: []int{5, 5, 5}},
				{WeekNumber: 2, Percentages: []float64{70, 80, 90}, Reps: []int{3, 3, 3}},
				{WeekNumber: 3, Percentages: []float64{75, 85, 95}, Reps: []int{5, 3, 1}},
			},
		},
		{
			name: "with percentage modifier",
			entries: []WeeklyLookupEntry{
				{WeekNumber: 1, Percentages: []float64{65}, Reps: []int{5}, PercentageModifier: floatPtr(0.9)},
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
		entries     []WeeklyLookupEntry
		expectedErr error
	}{
		{
			name:        "empty entries",
			entries:     []WeeklyLookupEntry{},
			expectedErr: ErrEntriesRequired,
		},
		{
			name:        "nil entries",
			entries:     nil,
			expectedErr: ErrEntriesRequired,
		},
		{
			name: "week number zero",
			entries: []WeeklyLookupEntry{
				{WeekNumber: 0, Percentages: []float64{65}, Reps: []int{5}},
			},
			expectedErr: ErrWeekNumberInvalid,
		},
		{
			name: "week number negative",
			entries: []WeeklyLookupEntry{
				{WeekNumber: -1, Percentages: []float64{65}, Reps: []int{5}},
			},
			expectedErr: ErrWeekNumberInvalid,
		},
		{
			name: "empty percentages",
			entries: []WeeklyLookupEntry{
				{WeekNumber: 1, Percentages: []float64{}, Reps: []int{5}},
			},
			expectedErr: ErrPercentagesRequired,
		},
		{
			name: "empty reps",
			entries: []WeeklyLookupEntry{
				{WeekNumber: 1, Percentages: []float64{65}, Reps: []int{}},
			},
			expectedErr: ErrRepsRequired,
		},
		{
			name: "mismatched lengths",
			entries: []WeeklyLookupEntry{
				{WeekNumber: 1, Percentages: []float64{65, 75, 85}, Reps: []int{5, 5}},
			},
			expectedErr: ErrPercentageRepsLengthMismatch,
		},
		{
			name: "duplicate week numbers",
			entries: []WeeklyLookupEntry{
				{WeekNumber: 1, Percentages: []float64{65}, Reps: []int{5}},
				{WeekNumber: 1, Percentages: []float64{70}, Reps: []int{3}},
			},
			expectedErr: ErrDuplicateWeekNumber,
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

// ==================== CreateWeeklyLookup Tests ====================

func TestCreateWeeklyLookup_ValidInput(t *testing.T) {
	input := CreateWeeklyLookupInput{
		Name: "5/3/1 Percentages",
		Entries: []WeeklyLookupEntry{
			{WeekNumber: 1, Percentages: []float64{65, 75, 85}, Reps: []int{5, 5, 5}},
			{WeekNumber: 2, Percentages: []float64{70, 80, 90}, Reps: []int{3, 3, 3}},
		},
	}

	lookup, result := CreateWeeklyLookup(input, "test-id")

	if !result.Valid {
		t.Errorf("CreateWeeklyLookup returned invalid result: %v", result.Errors)
	}
	if lookup == nil {
		t.Fatal("CreateWeeklyLookup returned nil lookup")
	}
	if lookup.ID != "test-id" {
		t.Errorf("lookup.ID = %q, want %q", lookup.ID, "test-id")
	}
	if lookup.Name != "5/3/1 Percentages" {
		t.Errorf("lookup.Name = %q, want %q", lookup.Name, "5/3/1 Percentages")
	}
	if len(lookup.Entries) != 2 {
		t.Errorf("len(lookup.Entries) = %d, want %d", len(lookup.Entries), 2)
	}
}

func TestCreateWeeklyLookup_WithProgramID(t *testing.T) {
	programID := "program-123"
	input := CreateWeeklyLookupInput{
		Name: "Custom Percentages",
		Entries: []WeeklyLookupEntry{
			{WeekNumber: 1, Percentages: []float64{65}, Reps: []int{5}},
		},
		ProgramID: &programID,
	}

	lookup, result := CreateWeeklyLookup(input, "test-id")

	if !result.Valid {
		t.Errorf("CreateWeeklyLookup returned invalid result: %v", result.Errors)
	}
	if lookup.ProgramID == nil || *lookup.ProgramID != programID {
		t.Errorf("lookup.ProgramID = %v, want %q", lookup.ProgramID, programID)
	}
}

func TestCreateWeeklyLookup_TrimsWhitespace(t *testing.T) {
	input := CreateWeeklyLookupInput{
		Name: "  5/3/1 Percentages  ",
		Entries: []WeeklyLookupEntry{
			{WeekNumber: 1, Percentages: []float64{65}, Reps: []int{5}},
		},
	}

	lookup, result := CreateWeeklyLookup(input, "test-id")

	if !result.Valid {
		t.Errorf("CreateWeeklyLookup returned invalid result: %v", result.Errors)
	}
	if lookup.Name != "5/3/1 Percentages" {
		t.Errorf("lookup.Name = %q, want %q (trimmed)", lookup.Name, "5/3/1 Percentages")
	}
}

func TestCreateWeeklyLookup_EmptyName(t *testing.T) {
	input := CreateWeeklyLookupInput{
		Name: "",
		Entries: []WeeklyLookupEntry{
			{WeekNumber: 1, Percentages: []float64{65}, Reps: []int{5}},
		},
	}

	lookup, result := CreateWeeklyLookup(input, "test-id")

	if result.Valid {
		t.Error("CreateWeeklyLookup with empty name returned valid result")
	}
	if lookup != nil {
		t.Error("CreateWeeklyLookup with invalid input returned non-nil lookup")
	}
}

func TestCreateWeeklyLookup_EmptyEntries(t *testing.T) {
	input := CreateWeeklyLookupInput{
		Name:    "Test Lookup",
		Entries: []WeeklyLookupEntry{},
	}

	lookup, result := CreateWeeklyLookup(input, "test-id")

	if result.Valid {
		t.Error("CreateWeeklyLookup with empty entries returned valid result")
	}
	if lookup != nil {
		t.Error("CreateWeeklyLookup with invalid input returned non-nil lookup")
	}
}

func TestCreateWeeklyLookup_MultipleErrors(t *testing.T) {
	input := CreateWeeklyLookupInput{
		Name:    "",                       // Invalid
		Entries: []WeeklyLookupEntry{},    // Invalid
	}

	lookup, result := CreateWeeklyLookup(input, "test-id")

	if result.Valid {
		t.Error("CreateWeeklyLookup with multiple errors returned valid result")
	}
	if lookup != nil {
		t.Error("CreateWeeklyLookup with invalid input returned non-nil lookup")
	}
	if len(result.Errors) < 2 {
		t.Errorf("Expected at least 2 errors, got %d", len(result.Errors))
	}
}

// ==================== UpdateWeeklyLookup Tests ====================

func TestUpdateWeeklyLookup_UpdateName(t *testing.T) {
	lookup := &WeeklyLookup{
		ID:   "test-id",
		Name: "Old Name",
		Entries: []WeeklyLookupEntry{
			{WeekNumber: 1, Percentages: []float64{65}, Reps: []int{5}},
		},
	}

	newName := "New Name"
	input := UpdateWeeklyLookupInput{Name: &newName}

	result := UpdateWeeklyLookup(lookup, input)

	if !result.Valid {
		t.Errorf("UpdateWeeklyLookup returned invalid result: %v", result.Errors)
	}
	if lookup.Name != "New Name" {
		t.Errorf("lookup.Name = %q, want %q", lookup.Name, "New Name")
	}
}

func TestUpdateWeeklyLookup_UpdateEntries(t *testing.T) {
	lookup := &WeeklyLookup{
		ID:   "test-id",
		Name: "Test Lookup",
		Entries: []WeeklyLookupEntry{
			{WeekNumber: 1, Percentages: []float64{65}, Reps: []int{5}},
		},
	}

	newEntries := []WeeklyLookupEntry{
		{WeekNumber: 1, Percentages: []float64{70, 80, 90}, Reps: []int{3, 3, 3}},
		{WeekNumber: 2, Percentages: []float64{75, 85, 95}, Reps: []int{5, 3, 1}},
	}
	input := UpdateWeeklyLookupInput{Entries: &newEntries}

	result := UpdateWeeklyLookup(lookup, input)

	if !result.Valid {
		t.Errorf("UpdateWeeklyLookup returned invalid result: %v", result.Errors)
	}
	if len(lookup.Entries) != 2 {
		t.Errorf("len(lookup.Entries) = %d, want %d", len(lookup.Entries), 2)
	}
}

func TestUpdateWeeklyLookup_UpdateProgramID(t *testing.T) {
	lookup := &WeeklyLookup{
		ID:   "test-id",
		Name: "Test Lookup",
		Entries: []WeeklyLookupEntry{
			{WeekNumber: 1, Percentages: []float64{65}, Reps: []int{5}},
		},
		ProgramID: nil,
	}

	newProgramID := "program-456"
	newProgramIDPtr := &newProgramID
	input := UpdateWeeklyLookupInput{ProgramID: &newProgramIDPtr}

	result := UpdateWeeklyLookup(lookup, input)

	if !result.Valid {
		t.Errorf("UpdateWeeklyLookup returned invalid result: %v", result.Errors)
	}
	if lookup.ProgramID == nil || *lookup.ProgramID != newProgramID {
		t.Errorf("lookup.ProgramID = %v, want %q", lookup.ProgramID, newProgramID)
	}
}

func TestUpdateWeeklyLookup_ClearProgramID(t *testing.T) {
	programID := "program-123"
	lookup := &WeeklyLookup{
		ID:   "test-id",
		Name: "Test Lookup",
		Entries: []WeeklyLookupEntry{
			{WeekNumber: 1, Percentages: []float64{65}, Reps: []int{5}},
		},
		ProgramID: &programID,
	}

	var nilProgramID *string = nil
	input := UpdateWeeklyLookupInput{ProgramID: &nilProgramID}

	result := UpdateWeeklyLookup(lookup, input)

	if !result.Valid {
		t.Errorf("UpdateWeeklyLookup returned invalid result: %v", result.Errors)
	}
	if lookup.ProgramID != nil {
		t.Errorf("lookup.ProgramID = %v, want nil", lookup.ProgramID)
	}
}

func TestUpdateWeeklyLookup_InvalidName(t *testing.T) {
	lookup := &WeeklyLookup{
		ID:   "test-id",
		Name: "Old Name",
		Entries: []WeeklyLookupEntry{
			{WeekNumber: 1, Percentages: []float64{65}, Reps: []int{5}},
		},
	}
	originalName := lookup.Name

	emptyName := ""
	input := UpdateWeeklyLookupInput{Name: &emptyName}

	result := UpdateWeeklyLookup(lookup, input)

	if result.Valid {
		t.Error("UpdateWeeklyLookup with invalid name returned valid result")
	}
	if lookup.Name != originalName {
		t.Errorf("lookup.Name was changed despite validation failure")
	}
}

func TestUpdateWeeklyLookup_InvalidEntries(t *testing.T) {
	lookup := &WeeklyLookup{
		ID:   "test-id",
		Name: "Test Lookup",
		Entries: []WeeklyLookupEntry{
			{WeekNumber: 1, Percentages: []float64{65}, Reps: []int{5}},
		},
	}
	originalEntries := lookup.Entries

	emptyEntries := []WeeklyLookupEntry{}
	input := UpdateWeeklyLookupInput{Entries: &emptyEntries}

	result := UpdateWeeklyLookup(lookup, input)

	if result.Valid {
		t.Error("UpdateWeeklyLookup with invalid entries returned valid result")
	}
	if len(lookup.Entries) != len(originalEntries) {
		t.Errorf("lookup.Entries was changed despite validation failure")
	}
}

// ==================== GetByWeekNumber Tests ====================

func TestWeeklyLookup_GetByWeekNumber_Found(t *testing.T) {
	lookup := &WeeklyLookup{
		ID:   "test-id",
		Name: "Test Lookup",
		Entries: []WeeklyLookupEntry{
			{WeekNumber: 1, Percentages: []float64{65, 75, 85}, Reps: []int{5, 5, 5}},
			{WeekNumber: 2, Percentages: []float64{70, 80, 90}, Reps: []int{3, 3, 3}},
			{WeekNumber: 3, Percentages: []float64{75, 85, 95}, Reps: []int{5, 3, 1}},
		},
	}

	entry := lookup.GetByWeekNumber(2)

	if entry == nil {
		t.Fatal("GetByWeekNumber returned nil for existing week")
	}
	if entry.WeekNumber != 2 {
		t.Errorf("entry.WeekNumber = %d, want %d", entry.WeekNumber, 2)
	}
	if entry.Percentages[0] != 70 {
		t.Errorf("entry.Percentages[0] = %f, want %f", entry.Percentages[0], 70.0)
	}
}

func TestWeeklyLookup_GetByWeekNumber_NotFound(t *testing.T) {
	lookup := &WeeklyLookup{
		ID:   "test-id",
		Name: "Test Lookup",
		Entries: []WeeklyLookupEntry{
			{WeekNumber: 1, Percentages: []float64{65}, Reps: []int{5}},
			{WeekNumber: 2, Percentages: []float64{70}, Reps: []int{3}},
		},
	}

	entry := lookup.GetByWeekNumber(5)

	if entry != nil {
		t.Errorf("GetByWeekNumber returned non-nil for non-existing week: %v", entry)
	}
}

// ==================== Validate Tests ====================

func TestWeeklyLookup_Validate_Valid(t *testing.T) {
	lookup := &WeeklyLookup{
		ID:   "test-id",
		Name: "Test Lookup",
		Entries: []WeeklyLookupEntry{
			{WeekNumber: 1, Percentages: []float64{65}, Reps: []int{5}},
		},
	}

	result := lookup.Validate()

	if !result.Valid {
		t.Errorf("Validate returned invalid result for valid lookup: %v", result.Errors)
	}
}

func TestWeeklyLookup_Validate_InvalidName(t *testing.T) {
	lookup := &WeeklyLookup{
		ID:   "test-id",
		Name: "",
		Entries: []WeeklyLookupEntry{
			{WeekNumber: 1, Percentages: []float64{65}, Reps: []int{5}},
		},
	}

	result := lookup.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for lookup with invalid name")
	}
}

func TestWeeklyLookup_Validate_InvalidEntries(t *testing.T) {
	lookup := &WeeklyLookup{
		ID:      "test-id",
		Name:    "Test Lookup",
		Entries: []WeeklyLookupEntry{},
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

func floatPtr(f float64) *float64 {
	return &f
}

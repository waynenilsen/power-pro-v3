package rotationlookup

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
		{"Conjugate Max Effort Rotation"},
		{"Westside Main Lift Rotation"},
		{"Simple Rotation"},
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

// ==================== Entry Validation Tests ====================

func TestValidateEntry_Valid(t *testing.T) {
	tests := []struct {
		name  string
		entry RotationLookupEntry
	}{
		{
			name:  "basic entry position 0",
			entry: RotationLookupEntry{Position: 0, LiftIdentifier: "deadlift"},
		},
		{
			name:  "entry with description",
			entry: RotationLookupEntry{Position: 1, LiftIdentifier: "squat", Description: "High bar squat focus"},
		},
		{
			name:  "entry with higher position",
			entry: RotationLookupEntry{Position: 10, LiftIdentifier: "bench"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEntry(tt.entry)
			if err != nil {
				t.Errorf("ValidateEntry() = %v, want nil", err)
			}
		})
	}
}

func TestValidateEntry_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		entry       RotationLookupEntry
		expectedErr error
	}{
		{
			name:        "negative position",
			entry:       RotationLookupEntry{Position: -1, LiftIdentifier: "deadlift"},
			expectedErr: ErrPositionInvalid,
		},
		{
			name:        "empty lift identifier",
			entry:       RotationLookupEntry{Position: 0, LiftIdentifier: ""},
			expectedErr: ErrLiftIdentifierRequired,
		},
		{
			name:        "whitespace only lift identifier",
			entry:       RotationLookupEntry{Position: 0, LiftIdentifier: "   "},
			expectedErr: ErrLiftIdentifierRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEntry(tt.entry)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateEntry() = %v, want %v", err, tt.expectedErr)
			}
		})
	}
}

// ==================== Entries Validation Tests ====================

func TestValidateEntries_Valid(t *testing.T) {
	tests := []struct {
		name    string
		entries []RotationLookupEntry
	}{
		{
			name: "single entry",
			entries: []RotationLookupEntry{
				{Position: 0, LiftIdentifier: "deadlift"},
			},
		},
		{
			name: "multiple entries",
			entries: []RotationLookupEntry{
				{Position: 0, LiftIdentifier: "deadlift", Description: "Deadlift Focus"},
				{Position: 1, LiftIdentifier: "squat", Description: "Squat Focus"},
				{Position: 2, LiftIdentifier: "bench", Description: "Bench Focus"},
			},
		},
		{
			name: "non-sequential positions",
			entries: []RotationLookupEntry{
				{Position: 0, LiftIdentifier: "deadlift"},
				{Position: 2, LiftIdentifier: "squat"},
				{Position: 5, LiftIdentifier: "bench"},
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
		entries     []RotationLookupEntry
		expectedErr error
	}{
		{
			name:        "empty entries",
			entries:     []RotationLookupEntry{},
			expectedErr: ErrEntriesRequired,
		},
		{
			name:        "nil entries",
			entries:     nil,
			expectedErr: ErrEntriesRequired,
		},
		{
			name: "negative position",
			entries: []RotationLookupEntry{
				{Position: -1, LiftIdentifier: "deadlift"},
			},
			expectedErr: ErrPositionInvalid,
		},
		{
			name: "empty lift identifier",
			entries: []RotationLookupEntry{
				{Position: 0, LiftIdentifier: ""},
			},
			expectedErr: ErrLiftIdentifierRequired,
		},
		{
			name: "duplicate positions",
			entries: []RotationLookupEntry{
				{Position: 0, LiftIdentifier: "deadlift"},
				{Position: 0, LiftIdentifier: "squat"},
			},
			expectedErr: ErrDuplicatePosition,
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

// ==================== CreateRotationLookup Tests ====================

func TestCreateRotationLookup_ValidInput(t *testing.T) {
	input := CreateRotationLookupInput{
		Name: "Conjugate Max Effort Rotation",
		Entries: []RotationLookupEntry{
			{Position: 0, LiftIdentifier: "deadlift", Description: "Deadlift Focus"},
			{Position: 1, LiftIdentifier: "squat", Description: "Squat Focus"},
		},
	}

	lookup, result := CreateRotationLookup(input, "test-id")

	if !result.Valid {
		t.Errorf("CreateRotationLookup returned invalid result: %v", result.Errors)
	}
	if lookup == nil {
		t.Fatal("CreateRotationLookup returned nil lookup")
	}
	if lookup.ID != "test-id" {
		t.Errorf("lookup.ID = %q, want %q", lookup.ID, "test-id")
	}
	if lookup.Name != "Conjugate Max Effort Rotation" {
		t.Errorf("lookup.Name = %q, want %q", lookup.Name, "Conjugate Max Effort Rotation")
	}
	if len(lookup.Entries) != 2 {
		t.Errorf("len(lookup.Entries) = %d, want %d", len(lookup.Entries), 2)
	}
}

func TestCreateRotationLookup_WithProgramID(t *testing.T) {
	programID := "program-123"
	input := CreateRotationLookupInput{
		Name: "Custom Rotation",
		Entries: []RotationLookupEntry{
			{Position: 0, LiftIdentifier: "deadlift"},
		},
		ProgramID: &programID,
	}

	lookup, result := CreateRotationLookup(input, "test-id")

	if !result.Valid {
		t.Errorf("CreateRotationLookup returned invalid result: %v", result.Errors)
	}
	if lookup.ProgramID == nil || *lookup.ProgramID != programID {
		t.Errorf("lookup.ProgramID = %v, want %q", lookup.ProgramID, programID)
	}
}

func TestCreateRotationLookup_TrimsWhitespace(t *testing.T) {
	input := CreateRotationLookupInput{
		Name: "  Conjugate Rotation  ",
		Entries: []RotationLookupEntry{
			{Position: 0, LiftIdentifier: "deadlift"},
		},
	}

	lookup, result := CreateRotationLookup(input, "test-id")

	if !result.Valid {
		t.Errorf("CreateRotationLookup returned invalid result: %v", result.Errors)
	}
	if lookup.Name != "Conjugate Rotation" {
		t.Errorf("lookup.Name = %q, want %q (trimmed)", lookup.Name, "Conjugate Rotation")
	}
}

func TestCreateRotationLookup_EmptyName(t *testing.T) {
	input := CreateRotationLookupInput{
		Name: "",
		Entries: []RotationLookupEntry{
			{Position: 0, LiftIdentifier: "deadlift"},
		},
	}

	lookup, result := CreateRotationLookup(input, "test-id")

	if result.Valid {
		t.Error("CreateRotationLookup with empty name returned valid result")
	}
	if lookup != nil {
		t.Error("CreateRotationLookup with invalid input returned non-nil lookup")
	}
}

func TestCreateRotationLookup_EmptyEntries(t *testing.T) {
	input := CreateRotationLookupInput{
		Name:    "Test Rotation",
		Entries: []RotationLookupEntry{},
	}

	lookup, result := CreateRotationLookup(input, "test-id")

	if result.Valid {
		t.Error("CreateRotationLookup with empty entries returned valid result")
	}
	if lookup != nil {
		t.Error("CreateRotationLookup with invalid input returned non-nil lookup")
	}
}

func TestCreateRotationLookup_MultipleErrors(t *testing.T) {
	input := CreateRotationLookupInput{
		Name:    "",                          // Invalid
		Entries: []RotationLookupEntry{},     // Invalid
	}

	lookup, result := CreateRotationLookup(input, "test-id")

	if result.Valid {
		t.Error("CreateRotationLookup with multiple errors returned valid result")
	}
	if lookup != nil {
		t.Error("CreateRotationLookup with invalid input returned non-nil lookup")
	}
	if len(result.Errors) < 2 {
		t.Errorf("Expected at least 2 errors, got %d", len(result.Errors))
	}
}

// ==================== UpdateRotationLookup Tests ====================

func TestUpdateRotationLookup_UpdateName(t *testing.T) {
	lookup := &RotationLookup{
		ID:   "test-id",
		Name: "Old Name",
		Entries: []RotationLookupEntry{
			{Position: 0, LiftIdentifier: "deadlift"},
		},
	}

	newName := "New Name"
	input := UpdateRotationLookupInput{Name: &newName}

	result := UpdateRotationLookup(lookup, input)

	if !result.Valid {
		t.Errorf("UpdateRotationLookup returned invalid result: %v", result.Errors)
	}
	if lookup.Name != "New Name" {
		t.Errorf("lookup.Name = %q, want %q", lookup.Name, "New Name")
	}
}

func TestUpdateRotationLookup_UpdateEntries(t *testing.T) {
	lookup := &RotationLookup{
		ID:   "test-id",
		Name: "Test Rotation",
		Entries: []RotationLookupEntry{
			{Position: 0, LiftIdentifier: "deadlift"},
		},
	}

	newEntries := []RotationLookupEntry{
		{Position: 0, LiftIdentifier: "squat", Description: "Squat Focus"},
		{Position: 1, LiftIdentifier: "bench", Description: "Bench Focus"},
	}
	input := UpdateRotationLookupInput{Entries: &newEntries}

	result := UpdateRotationLookup(lookup, input)

	if !result.Valid {
		t.Errorf("UpdateRotationLookup returned invalid result: %v", result.Errors)
	}
	if len(lookup.Entries) != 2 {
		t.Errorf("len(lookup.Entries) = %d, want %d", len(lookup.Entries), 2)
	}
}

func TestUpdateRotationLookup_UpdateProgramID(t *testing.T) {
	lookup := &RotationLookup{
		ID:   "test-id",
		Name: "Test Rotation",
		Entries: []RotationLookupEntry{
			{Position: 0, LiftIdentifier: "deadlift"},
		},
		ProgramID: nil,
	}

	newProgramID := "program-456"
	newProgramIDPtr := &newProgramID
	input := UpdateRotationLookupInput{ProgramID: &newProgramIDPtr}

	result := UpdateRotationLookup(lookup, input)

	if !result.Valid {
		t.Errorf("UpdateRotationLookup returned invalid result: %v", result.Errors)
	}
	if lookup.ProgramID == nil || *lookup.ProgramID != newProgramID {
		t.Errorf("lookup.ProgramID = %v, want %q", lookup.ProgramID, newProgramID)
	}
}

func TestUpdateRotationLookup_ClearProgramID(t *testing.T) {
	programID := "program-123"
	lookup := &RotationLookup{
		ID:   "test-id",
		Name: "Test Rotation",
		Entries: []RotationLookupEntry{
			{Position: 0, LiftIdentifier: "deadlift"},
		},
		ProgramID: &programID,
	}

	var nilProgramID *string = nil
	input := UpdateRotationLookupInput{ProgramID: &nilProgramID}

	result := UpdateRotationLookup(lookup, input)

	if !result.Valid {
		t.Errorf("UpdateRotationLookup returned invalid result: %v", result.Errors)
	}
	if lookup.ProgramID != nil {
		t.Errorf("lookup.ProgramID = %v, want nil", lookup.ProgramID)
	}
}

func TestUpdateRotationLookup_InvalidName(t *testing.T) {
	lookup := &RotationLookup{
		ID:   "test-id",
		Name: "Old Name",
		Entries: []RotationLookupEntry{
			{Position: 0, LiftIdentifier: "deadlift"},
		},
	}
	originalName := lookup.Name

	emptyName := ""
	input := UpdateRotationLookupInput{Name: &emptyName}

	result := UpdateRotationLookup(lookup, input)

	if result.Valid {
		t.Error("UpdateRotationLookup with invalid name returned valid result")
	}
	if lookup.Name != originalName {
		t.Errorf("lookup.Name was changed despite validation failure")
	}
}

func TestUpdateRotationLookup_InvalidEntries(t *testing.T) {
	lookup := &RotationLookup{
		ID:   "test-id",
		Name: "Test Rotation",
		Entries: []RotationLookupEntry{
			{Position: 0, LiftIdentifier: "deadlift"},
		},
	}
	originalEntries := lookup.Entries

	emptyEntries := []RotationLookupEntry{}
	input := UpdateRotationLookupInput{Entries: &emptyEntries}

	result := UpdateRotationLookup(lookup, input)

	if result.Valid {
		t.Error("UpdateRotationLookup with invalid entries returned valid result")
	}
	if len(lookup.Entries) != len(originalEntries) {
		t.Errorf("lookup.Entries was changed despite validation failure")
	}
}

// ==================== GetByPosition Tests ====================

func TestRotationLookup_GetByPosition_Found(t *testing.T) {
	lookup := &RotationLookup{
		ID:   "test-id",
		Name: "Test Rotation",
		Entries: []RotationLookupEntry{
			{Position: 0, LiftIdentifier: "deadlift", Description: "Deadlift Focus"},
			{Position: 1, LiftIdentifier: "squat", Description: "Squat Focus"},
			{Position: 2, LiftIdentifier: "bench", Description: "Bench Focus"},
		},
	}

	entry := lookup.GetByPosition(1)

	if entry == nil {
		t.Fatal("GetByPosition returned nil for existing position")
	}
	if entry.Position != 1 {
		t.Errorf("entry.Position = %d, want %d", entry.Position, 1)
	}
	if entry.LiftIdentifier != "squat" {
		t.Errorf("entry.LiftIdentifier = %q, want %q", entry.LiftIdentifier, "squat")
	}
}

func TestRotationLookup_GetByPosition_NotFound(t *testing.T) {
	lookup := &RotationLookup{
		ID:   "test-id",
		Name: "Test Rotation",
		Entries: []RotationLookupEntry{
			{Position: 0, LiftIdentifier: "deadlift"},
			{Position: 1, LiftIdentifier: "squat"},
		},
	}

	entry := lookup.GetByPosition(5)

	if entry != nil {
		t.Errorf("GetByPosition returned non-nil for non-existing position: %v", entry)
	}
}

func TestRotationLookup_GetByPosition_ZeroPosition(t *testing.T) {
	lookup := &RotationLookup{
		ID:   "test-id",
		Name: "Test Rotation",
		Entries: []RotationLookupEntry{
			{Position: 0, LiftIdentifier: "deadlift"},
		},
	}

	entry := lookup.GetByPosition(0)

	if entry == nil {
		t.Fatal("GetByPosition returned nil for position 0")
	}
	if entry.LiftIdentifier != "deadlift" {
		t.Errorf("entry.LiftIdentifier = %q, want %q", entry.LiftIdentifier, "deadlift")
	}
}

// ==================== Length Tests ====================

func TestRotationLookup_Length(t *testing.T) {
	tests := []struct {
		name     string
		entries  []RotationLookupEntry
		expected int
	}{
		{
			name:     "empty",
			entries:  []RotationLookupEntry{},
			expected: 0,
		},
		{
			name: "single entry",
			entries: []RotationLookupEntry{
				{Position: 0, LiftIdentifier: "deadlift"},
			},
			expected: 1,
		},
		{
			name: "multiple entries",
			entries: []RotationLookupEntry{
				{Position: 0, LiftIdentifier: "deadlift"},
				{Position: 1, LiftIdentifier: "squat"},
				{Position: 2, LiftIdentifier: "bench"},
			},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lookup := &RotationLookup{
				ID:      "test-id",
				Name:    "Test Rotation",
				Entries: tt.entries,
			}

			if got := lookup.Length(); got != tt.expected {
				t.Errorf("Length() = %d, want %d", got, tt.expected)
			}
		})
	}
}

// ==================== ContainsLift Tests ====================

func TestRotationLookup_ContainsLift(t *testing.T) {
	lookup := &RotationLookup{
		ID:   "test-id",
		Name: "Test Rotation",
		Entries: []RotationLookupEntry{
			{Position: 0, LiftIdentifier: "deadlift"},
			{Position: 1, LiftIdentifier: "squat"},
			{Position: 2, LiftIdentifier: "bench"},
		},
	}

	tests := []struct {
		liftIdentifier string
		expected       bool
	}{
		{"deadlift", true},
		{"squat", true},
		{"bench", true},
		{"overhead_press", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.liftIdentifier, func(t *testing.T) {
			if got := lookup.ContainsLift(tt.liftIdentifier); got != tt.expected {
				t.Errorf("ContainsLift(%q) = %v, want %v", tt.liftIdentifier, got, tt.expected)
			}
		})
	}
}

// ==================== Validate Tests ====================

func TestRotationLookup_Validate_Valid(t *testing.T) {
	lookup := &RotationLookup{
		ID:   "test-id",
		Name: "Test Rotation",
		Entries: []RotationLookupEntry{
			{Position: 0, LiftIdentifier: "deadlift"},
		},
	}

	result := lookup.Validate()

	if !result.Valid {
		t.Errorf("Validate returned invalid result for valid lookup: %v", result.Errors)
	}
}

func TestRotationLookup_Validate_InvalidName(t *testing.T) {
	lookup := &RotationLookup{
		ID:   "test-id",
		Name: "",
		Entries: []RotationLookupEntry{
			{Position: 0, LiftIdentifier: "deadlift"},
		},
	}

	result := lookup.Validate()

	if result.Valid {
		t.Error("Validate returned valid result for lookup with invalid name")
	}
}

func TestRotationLookup_Validate_InvalidEntries(t *testing.T) {
	lookup := &RotationLookup{
		ID:      "test-id",
		Name:    "Test Rotation",
		Entries: []RotationLookupEntry{},
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

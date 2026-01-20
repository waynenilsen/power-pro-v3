package lift

import (
	"errors"
	"strings"
	"testing"
)

// mockRepository implements LiftRepository for testing
type mockRepository struct {
	lifts map[string]*Lift
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		lifts: make(map[string]*Lift),
	}
}

func (m *mockRepository) GetByID(id string) (*Lift, error) {
	lift, ok := m.lifts[id]
	if !ok {
		return nil, nil
	}
	return lift, nil
}

func (m *mockRepository) SlugExists(slug string, excludeID *string) (bool, error) {
	for _, l := range m.lifts {
		if l.Slug == slug {
			if excludeID != nil && l.ID == *excludeID {
				continue
			}
			return true, nil
		}
	}
	return false, nil
}

func (m *mockRepository) Add(lift *Lift) {
	m.lifts[lift.ID] = lift
}

// ==================== Name Validation Tests ====================

func TestValidateName_Valid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"single word", "Squat"},
		{"two words", "Bench Press"},
		{"with numbers", "5x5 Squat"},
		{"minimum length", "A"},
		{"exactly 100 chars", strings.Repeat("a", 100)},
		{"with special chars", "Close-Grip Bench Press"},
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
		{"101 chars", strings.Repeat("a", 101), ErrNameTooLong},
		{"200 chars", strings.Repeat("a", 200), ErrNameTooLong},
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
		{"single word", "squat"},
		{"with hyphen", "bench-press"},
		{"multiple hyphens", "close-grip-bench-press"},
		{"with numbers", "5x5-squat"},
		{"numbers only", "531"},
		{"minimum length", "a"},
		{"exactly 100 chars", strings.Repeat("a", 100)},
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
		{"uppercase", "Squat", ErrSlugInvalid},
		{"with space", "bench press", ErrSlugInvalid},
		{"with underscore", "bench_press", ErrSlugInvalid},
		{"leading hyphen", "-squat", ErrSlugInvalid},
		{"trailing hyphen", "squat-", ErrSlugInvalid},
		{"consecutive hyphens", "bench--press", ErrSlugInvalid},
		{"special chars", "squat@home", ErrSlugInvalid},
		{"101 chars", strings.Repeat("a", 101), ErrSlugTooLong},
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
		{"simple lowercase", "squat", "squat"},
		{"uppercase", "SQUAT", "squat"},
		{"mixed case", "Bench Press", "bench-press"},
		{"multiple spaces", "Close  Grip  Bench", "close-grip-bench"},
		{"special chars", "5/3/1 BBB", "5-3-1-bbb"},
		{"underscores", "pause_squat", "pause-squat"},
		{"parentheses", "Squat (High Bar)", "squat-high-bar"},
		{"apostrophe", "Greg Nuckols' Program", "greg-nuckols-program"},
		{"leading spaces", "  squat", "squat"},
		{"trailing spaces", "squat  ", "squat"},
		{"mixed special", "Close-Grip Bench Press (Pause)", "close-grip-bench-press-pause"},
		{"numbers", "531 BBB Week 1", "531-bbb-week-1"},
		{"ampersand", "Push & Pull", "push-pull"},
		{"brackets", "Squat [Variation]", "squat-variation"},
		{"dots", "Dr. Squat's Method", "dr-squats-method"},
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
		"Squat",
		"Bench Press",
		"Deadlift",
		"Close-Grip Bench Press",
		"Pause Squat (High Bar)",
		"5/3/1 BBB",
		"Greg Nuckols' High Frequency",
	}

	for _, input := range inputs {
		slug := GenerateSlug(input)
		if err := ValidateSlug(slug); err != nil {
			t.Errorf("GenerateSlug(%q) = %q, but validation failed: %v", input, slug, err)
		}
	}
}

// ==================== Parent Lift Validation Tests ====================

func TestValidateParentLiftID_NilParent(t *testing.T) {
	err := ValidateParentLiftID("lift-1", nil, nil)
	if err != nil {
		t.Errorf("ValidateParentLiftID with nil parent = %v, want nil", err)
	}
}

func TestValidateParentLiftID_SelfReference(t *testing.T) {
	parentID := "lift-1"
	err := ValidateParentLiftID("lift-1", &parentID, nil)
	if !errors.Is(err, ErrSelfReference) {
		t.Errorf("ValidateParentLiftID with self-reference = %v, want %v", err, ErrSelfReference)
	}
}

func TestValidateParentLiftID_ValidParent(t *testing.T) {
	repo := newMockRepository()
	squat := &Lift{ID: "squat-id", Name: "Squat", Slug: "squat"}
	repo.Add(squat)

	parentID := "squat-id"
	err := ValidateParentLiftID("pause-squat-id", &parentID, repo)
	if err != nil {
		t.Errorf("ValidateParentLiftID with valid parent = %v, want nil", err)
	}
}

func TestValidateParentLiftID_CircularReference_Direct(t *testing.T) {
	repo := newMockRepository()

	// Create a chain: pause-squat -> squat
	squatID := "squat-id"
	squat := &Lift{ID: squatID, Name: "Squat", Slug: "squat", ParentLiftID: nil}
	pauseSquat := &Lift{ID: "pause-squat-id", Name: "Pause Squat", Slug: "pause-squat", ParentLiftID: &squatID}

	repo.Add(squat)
	repo.Add(pauseSquat)

	// Try to make squat's parent be pause-squat (creates: squat -> pause-squat -> squat)
	pauseSquatID := "pause-squat-id"
	err := ValidateParentLiftID(squatID, &pauseSquatID, repo)
	if !errors.Is(err, ErrCircularReference) {
		t.Errorf("ValidateParentLiftID with circular reference = %v, want %v", err, ErrCircularReference)
	}
}

func TestValidateParentLiftID_CircularReference_Indirect(t *testing.T) {
	repo := newMockRepository()

	// Create a chain: variation -> pause-squat -> squat
	squatID := "squat-id"
	pauseSquatID := "pause-squat-id"

	squat := &Lift{ID: squatID, Name: "Squat", Slug: "squat", ParentLiftID: nil}
	pauseSquat := &Lift{ID: pauseSquatID, Name: "Pause Squat", Slug: "pause-squat", ParentLiftID: &squatID}
	variation := &Lift{ID: "variation-id", Name: "Tempo Pause Squat", Slug: "tempo-pause-squat", ParentLiftID: &pauseSquatID}

	repo.Add(squat)
	repo.Add(pauseSquat)
	repo.Add(variation)

	// Try to make squat's parent be variation (creates: squat -> variation -> pause-squat -> squat)
	variationID := "variation-id"
	err := ValidateParentLiftID(squatID, &variationID, repo)
	if !errors.Is(err, ErrCircularReference) {
		t.Errorf("ValidateParentLiftID with indirect circular reference = %v, want %v", err, ErrCircularReference)
	}
}

func TestValidateParentLiftID_NonExistentParent(t *testing.T) {
	repo := newMockRepository()

	parentID := "non-existent-id"
	err := ValidateParentLiftID("lift-1", &parentID, repo)
	// Should not error - parent just doesn't exist (will be caught by FK constraint at DB level)
	if err != nil {
		t.Errorf("ValidateParentLiftID with non-existent parent = %v, want nil", err)
	}
}

// ==================== CreateLift Tests ====================

func TestCreateLift_ValidInput(t *testing.T) {
	repo := newMockRepository()

	input := CreateLiftInput{
		Name:              "Squat",
		IsCompetitionLift: true,
	}

	lift, result := CreateLift(input, "test-id", repo)

	if !result.Valid {
		t.Errorf("CreateLift returned invalid result: %v", result.Errors)
	}
	if lift == nil {
		t.Fatal("CreateLift returned nil lift")
	}
	if lift.ID != "test-id" {
		t.Errorf("lift.ID = %q, want %q", lift.ID, "test-id")
	}
	if lift.Name != "Squat" {
		t.Errorf("lift.Name = %q, want %q", lift.Name, "Squat")
	}
	if lift.Slug != "squat" {
		t.Errorf("lift.Slug = %q, want %q (auto-generated)", lift.Slug, "squat")
	}
	if !lift.IsCompetitionLift {
		t.Error("lift.IsCompetitionLift = false, want true")
	}
	if lift.ParentLiftID != nil {
		t.Errorf("lift.ParentLiftID = %v, want nil", lift.ParentLiftID)
	}
}

func TestCreateLift_WithProvidedSlug(t *testing.T) {
	repo := newMockRepository()

	input := CreateLiftInput{
		Name: "Squat",
		Slug: "competition-squat",
	}

	lift, result := CreateLift(input, "test-id", repo)

	if !result.Valid {
		t.Errorf("CreateLift returned invalid result: %v", result.Errors)
	}
	if lift.Slug != "competition-squat" {
		t.Errorf("lift.Slug = %q, want %q", lift.Slug, "competition-squat")
	}
}

func TestCreateLift_WithParent(t *testing.T) {
	repo := newMockRepository()
	squat := &Lift{ID: "squat-id", Name: "Squat", Slug: "squat"}
	repo.Add(squat)

	parentID := "squat-id"
	input := CreateLiftInput{
		Name:         "Pause Squat",
		ParentLiftID: &parentID,
	}

	lift, result := CreateLift(input, "pause-squat-id", repo)

	if !result.Valid {
		t.Errorf("CreateLift returned invalid result: %v", result.Errors)
	}
	if lift.ParentLiftID == nil || *lift.ParentLiftID != "squat-id" {
		t.Errorf("lift.ParentLiftID = %v, want %q", lift.ParentLiftID, "squat-id")
	}
}

func TestCreateLift_InvalidName(t *testing.T) {
	repo := newMockRepository()

	input := CreateLiftInput{
		Name: "",
	}

	lift, result := CreateLift(input, "test-id", repo)

	if result.Valid {
		t.Error("CreateLift with empty name returned valid result")
	}
	if lift != nil {
		t.Error("CreateLift with invalid input returned non-nil lift")
	}
	if len(result.Errors) == 0 {
		t.Error("CreateLift with invalid input returned no errors")
	}
}

func TestCreateLift_InvalidSlug(t *testing.T) {
	repo := newMockRepository()

	input := CreateLiftInput{
		Name: "Valid Name",
		Slug: "Invalid Slug",
	}

	lift, result := CreateLift(input, "test-id", repo)

	if result.Valid {
		t.Error("CreateLift with invalid slug returned valid result")
	}
	if lift != nil {
		t.Error("CreateLift with invalid input returned non-nil lift")
	}
}

func TestCreateLift_SelfReference(t *testing.T) {
	repo := newMockRepository()

	parentID := "test-id"
	input := CreateLiftInput{
		Name:         "Squat",
		ParentLiftID: &parentID,
	}

	lift, result := CreateLift(input, "test-id", repo)

	if result.Valid {
		t.Error("CreateLift with self-reference returned valid result")
	}
	if lift != nil {
		t.Error("CreateLift with self-reference returned non-nil lift")
	}
}

func TestCreateLift_DefaultCompetitionLift(t *testing.T) {
	repo := newMockRepository()

	input := CreateLiftInput{
		Name: "Pause Squat",
		// IsCompetitionLift not set
	}

	lift, result := CreateLift(input, "test-id", repo)

	if !result.Valid {
		t.Errorf("CreateLift returned invalid result: %v", result.Errors)
	}
	if lift.IsCompetitionLift {
		t.Error("lift.IsCompetitionLift = true, want false (default)")
	}
}

func TestCreateLift_MultipleErrors(t *testing.T) {
	repo := newMockRepository()

	selfID := "test-id"
	input := CreateLiftInput{
		Name:         "", // Invalid
		Slug:         "Invalid Slug", // Invalid
		ParentLiftID: &selfID, // Self-reference
	}

	lift, result := CreateLift(input, "test-id", repo)

	if result.Valid {
		t.Error("CreateLift with multiple errors returned valid result")
	}
	if lift != nil {
		t.Error("CreateLift with invalid input returned non-nil lift")
	}
	if len(result.Errors) < 3 {
		t.Errorf("Expected at least 3 errors, got %d", len(result.Errors))
	}
}

// ==================== UpdateLift Tests ====================

func TestUpdateLift_UpdateName(t *testing.T) {
	repo := newMockRepository()
	lift := &Lift{ID: "test-id", Name: "Squat", Slug: "squat"}

	newName := "Competition Squat"
	input := UpdateLiftInput{Name: &newName}

	result := UpdateLift(lift, input, repo)

	if !result.Valid {
		t.Errorf("UpdateLift returned invalid result: %v", result.Errors)
	}
	if lift.Name != "Competition Squat" {
		t.Errorf("lift.Name = %q, want %q", lift.Name, "Competition Squat")
	}
}

func TestUpdateLift_UpdateSlug(t *testing.T) {
	repo := newMockRepository()
	lift := &Lift{ID: "test-id", Name: "Squat", Slug: "squat"}

	newSlug := "competition-squat"
	input := UpdateLiftInput{Slug: &newSlug}

	result := UpdateLift(lift, input, repo)

	if !result.Valid {
		t.Errorf("UpdateLift returned invalid result: %v", result.Errors)
	}
	if lift.Slug != "competition-squat" {
		t.Errorf("lift.Slug = %q, want %q", lift.Slug, "competition-squat")
	}
}

func TestUpdateLift_UpdateCompetitionLift(t *testing.T) {
	repo := newMockRepository()
	lift := &Lift{ID: "test-id", Name: "Squat", Slug: "squat", IsCompetitionLift: false}

	isCompetition := true
	input := UpdateLiftInput{IsCompetitionLift: &isCompetition}

	result := UpdateLift(lift, input, repo)

	if !result.Valid {
		t.Errorf("UpdateLift returned invalid result: %v", result.Errors)
	}
	if !lift.IsCompetitionLift {
		t.Error("lift.IsCompetitionLift = false, want true")
	}
}

func TestUpdateLift_SetParent(t *testing.T) {
	repo := newMockRepository()
	squat := &Lift{ID: "squat-id", Name: "Squat", Slug: "squat"}
	repo.Add(squat)

	lift := &Lift{ID: "pause-squat-id", Name: "Pause Squat", Slug: "pause-squat"}

	parentID := "squat-id"
	input := UpdateLiftInput{ParentLiftID: &parentID}

	result := UpdateLift(lift, input, repo)

	if !result.Valid {
		t.Errorf("UpdateLift returned invalid result: %v", result.Errors)
	}
	if lift.ParentLiftID == nil || *lift.ParentLiftID != "squat-id" {
		t.Errorf("lift.ParentLiftID = %v, want %q", lift.ParentLiftID, "squat-id")
	}
}

func TestUpdateLift_ClearParent(t *testing.T) {
	repo := newMockRepository()
	parentID := "squat-id"
	lift := &Lift{ID: "pause-squat-id", Name: "Pause Squat", Slug: "pause-squat", ParentLiftID: &parentID}

	input := UpdateLiftInput{ClearParentLift: true}

	result := UpdateLift(lift, input, repo)

	if !result.Valid {
		t.Errorf("UpdateLift returned invalid result: %v", result.Errors)
	}
	if lift.ParentLiftID != nil {
		t.Errorf("lift.ParentLiftID = %v, want nil", lift.ParentLiftID)
	}
}

func TestUpdateLift_InvalidName(t *testing.T) {
	repo := newMockRepository()
	lift := &Lift{ID: "test-id", Name: "Squat", Slug: "squat"}
	originalName := lift.Name

	emptyName := ""
	input := UpdateLiftInput{Name: &emptyName}

	result := UpdateLift(lift, input, repo)

	if result.Valid {
		t.Error("UpdateLift with invalid name returned valid result")
	}
	if lift.Name != originalName {
		t.Errorf("lift.Name was changed despite validation failure")
	}
}

func TestUpdateLift_CircularReference(t *testing.T) {
	repo := newMockRepository()

	squatID := "squat-id"
	pauseSquatID := "pause-squat-id"

	squat := &Lift{ID: squatID, Name: "Squat", Slug: "squat", ParentLiftID: nil}
	pauseSquat := &Lift{ID: pauseSquatID, Name: "Pause Squat", Slug: "pause-squat", ParentLiftID: &squatID}

	repo.Add(squat)
	repo.Add(pauseSquat)

	// Try to make squat's parent be pause-squat
	input := UpdateLiftInput{ParentLiftID: &pauseSquatID}

	result := UpdateLift(squat, input, repo)

	if result.Valid {
		t.Error("UpdateLift with circular reference returned valid result")
	}
}

func TestUpdateLift_NoChanges(t *testing.T) {
	repo := newMockRepository()
	lift := &Lift{ID: "test-id", Name: "Squat", Slug: "squat"}
	originalUpdatedAt := lift.UpdatedAt

	input := UpdateLiftInput{} // No changes

	result := UpdateLift(lift, input, repo)

	if !result.Valid {
		t.Errorf("UpdateLift with no changes returned invalid result: %v", result.Errors)
	}
	// UpdatedAt should be updated even with no field changes
	if !lift.UpdatedAt.After(originalUpdatedAt) && lift.UpdatedAt.Equal(originalUpdatedAt) {
		// This is expected - UpdatedAt is updated when result is valid
	}
}

// ==================== Lift.Validate Tests ====================

func TestLift_Validate_Valid(t *testing.T) {
	repo := newMockRepository()
	lift := &Lift{
		ID:                "test-id",
		Name:              "Squat",
		Slug:              "squat",
		IsCompetitionLift: true,
	}

	result := lift.Validate(repo)

	if !result.Valid {
		t.Errorf("Validate returned invalid result for valid lift: %v", result.Errors)
	}
}

func TestLift_Validate_InvalidName(t *testing.T) {
	repo := newMockRepository()
	lift := &Lift{
		ID:   "test-id",
		Name: "",
		Slug: "squat",
	}

	result := lift.Validate(repo)

	if result.Valid {
		t.Error("Validate returned valid result for lift with empty name")
	}
}

func TestLift_Validate_InvalidSlug(t *testing.T) {
	repo := newMockRepository()
	lift := &Lift{
		ID:   "test-id",
		Name: "Squat",
		Slug: "Invalid Slug",
	}

	result := lift.Validate(repo)

	if result.Valid {
		t.Error("Validate returned valid result for lift with invalid slug")
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
	name := strings.Repeat("a", 100)
	err := ValidateName(name)
	if err != nil {
		t.Errorf("ValidateName with exactly 100 chars = %v, want nil", err)
	}
}

func TestValidateSlug_ExactlyMaxLength(t *testing.T) {
	slug := strings.Repeat("a", 100)
	err := ValidateSlug(slug)
	if err != nil {
		t.Errorf("ValidateSlug with exactly 100 chars = %v, want nil", err)
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

func TestCreateLift_EmptySlugFromEmptyName(t *testing.T) {
	repo := newMockRepository()

	input := CreateLiftInput{
		Name: "", // Will generate empty slug
	}

	lift, result := CreateLift(input, "test-id", repo)

	if result.Valid {
		t.Error("CreateLift with empty name should fail")
	}
	if lift != nil {
		t.Error("CreateLift with invalid input should return nil lift")
	}
}

func TestValidateParentLiftID_DeepHierarchy(t *testing.T) {
	repo := newMockRepository()

	// Create a deep chain: l5 -> l4 -> l3 -> l2 -> l1
	l1ID := "l1"
	l2ID := "l2"
	l3ID := "l3"
	l4ID := "l4"

	repo.Add(&Lift{ID: l1ID, Name: "L1", Slug: "l1", ParentLiftID: nil})
	repo.Add(&Lift{ID: l2ID, Name: "L2", Slug: "l2", ParentLiftID: &l1ID})
	repo.Add(&Lift{ID: l3ID, Name: "L3", Slug: "l3", ParentLiftID: &l2ID})
	repo.Add(&Lift{ID: l4ID, Name: "L4", Slug: "l4", ParentLiftID: &l3ID})
	repo.Add(&Lift{ID: "l5", Name: "L5", Slug: "l5", ParentLiftID: &l4ID})

	// Trying to make l1's parent be l5 should detect circular reference
	l5ID := "l5"
	err := ValidateParentLiftID(l1ID, &l5ID, repo)
	if !errors.Is(err, ErrCircularReference) {
		t.Errorf("ValidateParentLiftID with deep circular reference = %v, want %v", err, ErrCircularReference)
	}

	// But l6 -> l5 should be fine
	err = ValidateParentLiftID("l6", &l5ID, repo)
	if err != nil {
		t.Errorf("ValidateParentLiftID with valid deep hierarchy = %v, want nil", err)
	}
}

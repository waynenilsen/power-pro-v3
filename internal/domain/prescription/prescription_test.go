package prescription

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/domain/loadstrategy"
	"github.com/waynenilsen/power-pro-v3/internal/domain/setscheme"
)

// ==================== Mock Types ====================

// mockLoadStrategy implements loadstrategy.LoadStrategy for testing
type mockLoadStrategy struct {
	strategyType    loadstrategy.LoadStrategyType
	calculateResult float64
	calculateErr    error
	validateErr     error
}

func (m *mockLoadStrategy) Type() loadstrategy.LoadStrategyType {
	return m.strategyType
}

func (m *mockLoadStrategy) CalculateLoad(ctx context.Context, params loadstrategy.LoadCalculationParams) (float64, error) {
	if m.calculateErr != nil {
		return 0, m.calculateErr
	}
	return m.calculateResult, nil
}

func (m *mockLoadStrategy) Validate() error {
	return m.validateErr
}

// mockSetScheme implements setscheme.SetScheme for testing
type mockSetScheme struct {
	schemeType     setscheme.SetSchemeType
	generateResult []setscheme.GeneratedSet
	generateErr    error
	validateErr    error
}

func (m *mockSetScheme) Type() setscheme.SetSchemeType {
	return m.schemeType
}

func (m *mockSetScheme) GenerateSets(baseWeight float64, ctx setscheme.SetGenerationContext) ([]setscheme.GeneratedSet, error) {
	if m.generateErr != nil {
		return nil, m.generateErr
	}
	return m.generateResult, nil
}

func (m *mockSetScheme) Validate() error {
	return m.validateErr
}

// mockLiftLookup implements LiftLookup for testing
type mockLiftLookup struct {
	lifts map[string]*LiftInfo
	err   error
}

func newMockLiftLookup() *mockLiftLookup {
	return &mockLiftLookup{
		lifts: make(map[string]*LiftInfo),
	}
}

func (m *mockLiftLookup) GetLiftByID(ctx context.Context, liftID string) (*LiftInfo, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.lifts[liftID], nil
}

func (m *mockLiftLookup) SetLift(liftID string, lift *LiftInfo) {
	m.lifts[liftID] = lift
}

func (m *mockLiftLookup) SetError(err error) {
	m.err = err
}

// ==================== Helper Functions ====================

func validUUID() string {
	return "550e8400-e29b-41d4-a716-446655440000"
}

func anotherValidUUID() string {
	return "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
}

func validLoadStrategy() *mockLoadStrategy {
	return &mockLoadStrategy{
		strategyType:    loadstrategy.TypePercentOf,
		calculateResult: 225.0,
	}
}

func validSetScheme() *mockSetScheme {
	return &mockSetScheme{
		schemeType: setscheme.TypeFixed,
		generateResult: []setscheme.GeneratedSet{
			{SetNumber: 1, Weight: 225.0, TargetReps: 5, IsWorkSet: true},
			{SetNumber: 2, Weight: 225.0, TargetReps: 5, IsWorkSet: true},
			{SetNumber: 3, Weight: 225.0, TargetReps: 5, IsWorkSet: true},
		},
	}
}

// ==================== UUID Validation Tests ====================

func TestIsValidUUID_Valid(t *testing.T) {
	tests := []struct {
		name string
		uuid string
	}{
		{"standard lowercase", "550e8400-e29b-41d4-a716-446655440000"},
		{"uppercase", "550E8400-E29B-41D4-A716-446655440000"},
		{"mixed case", "550e8400-E29B-41d4-a716-446655440000"},
		{"all zeros", "00000000-0000-0000-0000-000000000000"},
		{"all fs", "ffffffff-ffff-ffff-ffff-ffffffffffff"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !isValidUUID(tt.uuid) {
				t.Errorf("isValidUUID(%q) = false, want true", tt.uuid)
			}
		})
	}
}

func TestIsValidUUID_Invalid(t *testing.T) {
	tests := []struct {
		name string
		uuid string
	}{
		{"empty", ""},
		{"too short", "550e8400-e29b-41d4-a716-44665544000"},
		{"too long", "550e8400-e29b-41d4-a716-4466554400000"},
		{"missing hyphen", "550e8400e29b-41d4-a716-446655440000"},
		{"extra hyphen", "550e-8400-e29b-41d4-a716-446655440000"},
		{"invalid character", "550e8400-e29b-41d4-a716-44665544000g"},
		{"hyphen in wrong position", "550e8400-e29-b41d4-a716-446655440000"},
		{"spaces", "550e8400 e29b 41d4 a716 446655440000"},
		{"no hyphens", "550e8400e29b41d4a716446655440000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if isValidUUID(tt.uuid) {
				t.Errorf("isValidUUID(%q) = true, want false", tt.uuid)
			}
		})
	}
}

// ==================== LiftID Validation Tests ====================

func TestValidateLiftID_Valid(t *testing.T) {
	err := ValidateLiftID(validUUID())
	if err != nil {
		t.Errorf("ValidateLiftID(%q) = %v, want nil", validUUID(), err)
	}
}

func TestValidateLiftID_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		liftID      string
		expectedErr error
	}{
		{"empty", "", ErrLiftIDRequired},
		{"only spaces", "   ", ErrLiftIDRequired},
		{"only tabs", "\t\t", ErrLiftIDRequired},
		{"invalid format", "not-a-uuid", ErrLiftIDInvalid},
		{"too short", "550e8400-e29b-41d4", ErrLiftIDInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLiftID(tt.liftID)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ValidateLiftID(%q) = %v, want %v", tt.liftID, err, tt.expectedErr)
			}
		})
	}
}

// ==================== LoadStrategy Validation Tests ====================

func TestValidateLoadStrategy_Valid(t *testing.T) {
	strategy := validLoadStrategy()
	err := ValidateLoadStrategy(strategy)
	if err != nil {
		t.Errorf("ValidateLoadStrategy() = %v, want nil", err)
	}
}

func TestValidateLoadStrategy_Nil(t *testing.T) {
	err := ValidateLoadStrategy(nil)
	if !errors.Is(err, ErrLoadStrategyRequired) {
		t.Errorf("ValidateLoadStrategy(nil) = %v, want %v", err, ErrLoadStrategyRequired)
	}
}

func TestValidateLoadStrategy_Invalid(t *testing.T) {
	strategy := &mockLoadStrategy{
		validateErr: errors.New("strategy is misconfigured"),
	}
	err := ValidateLoadStrategy(strategy)
	if err == nil {
		t.Error("ValidateLoadStrategy with invalid strategy should return error")
	}
}

// ==================== SetScheme Validation Tests ====================

func TestValidateSetScheme_Valid(t *testing.T) {
	scheme := validSetScheme()
	err := ValidateSetScheme(scheme)
	if err != nil {
		t.Errorf("ValidateSetScheme() = %v, want nil", err)
	}
}

func TestValidateSetScheme_Nil(t *testing.T) {
	err := ValidateSetScheme(nil)
	if !errors.Is(err, ErrSetSchemeRequired) {
		t.Errorf("ValidateSetScheme(nil) = %v, want %v", err, ErrSetSchemeRequired)
	}
}

func TestValidateSetScheme_Invalid(t *testing.T) {
	scheme := &mockSetScheme{
		validateErr: errors.New("scheme is misconfigured"),
	}
	err := ValidateSetScheme(scheme)
	if err == nil {
		t.Error("ValidateSetScheme with invalid scheme should return error")
	}
}

// ==================== Order Validation Tests ====================

func TestValidateOrder_Valid(t *testing.T) {
	tests := []struct {
		name  string
		order int
	}{
		{"zero", 0},
		{"positive", 5},
		{"large positive", 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOrder(tt.order)
			if err != nil {
				t.Errorf("ValidateOrder(%d) = %v, want nil", tt.order, err)
			}
		})
	}
}

func TestValidateOrder_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		order int
	}{
		{"negative one", -1},
		{"large negative", -100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOrder(tt.order)
			if !errors.Is(err, ErrOrderNegative) {
				t.Errorf("ValidateOrder(%d) = %v, want %v", tt.order, err, ErrOrderNegative)
			}
		})
	}
}

// ==================== Notes Validation Tests ====================

func TestValidateNotes_Valid(t *testing.T) {
	tests := []struct {
		name  string
		notes string
	}{
		{"empty", ""},
		{"short", "Focus on form"},
		{"exactly 500 chars", strings.Repeat("a", 500)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNotes(tt.notes)
			if err != nil {
				t.Errorf("ValidateNotes(%q) = %v, want nil", tt.notes, err)
			}
		})
	}
}

func TestValidateNotes_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		notes string
	}{
		{"501 chars", strings.Repeat("a", 501)},
		{"1000 chars", strings.Repeat("a", 1000)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNotes(tt.notes)
			if !errors.Is(err, ErrNotesTooLong) {
				t.Errorf("ValidateNotes(len=%d) = %v, want %v", len(tt.notes), err, ErrNotesTooLong)
			}
		})
	}
}

// ==================== RestSeconds Validation Tests ====================

func TestValidateRestSeconds_Valid(t *testing.T) {
	tests := []struct {
		name        string
		restSeconds *int
	}{
		{"nil", nil},
		{"zero", intPtr(0)},
		{"positive", intPtr(90)},
		{"large positive", intPtr(300)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRestSeconds(tt.restSeconds)
			if err != nil {
				t.Errorf("ValidateRestSeconds(%v) = %v, want nil", tt.restSeconds, err)
			}
		})
	}
}

func TestValidateRestSeconds_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		restSeconds *int
	}{
		{"negative one", intPtr(-1)},
		{"large negative", intPtr(-100)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRestSeconds(tt.restSeconds)
			if !errors.Is(err, ErrRestSecondsNegative) {
				t.Errorf("ValidateRestSeconds(%v) = %v, want %v", *tt.restSeconds, err, ErrRestSecondsNegative)
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}

// ==================== CreatePrescription Tests ====================

func TestCreatePrescription_ValidInput(t *testing.T) {
	input := CreatePrescriptionInput{
		LiftID:       validUUID(),
		LoadStrategy: validLoadStrategy(),
		SetScheme:    validSetScheme(),
		Order:        1,
		Notes:        "Focus on form",
		RestSeconds:  intPtr(90),
	}

	prescription, result := CreatePrescription(input, anotherValidUUID())

	if !result.Valid {
		t.Errorf("CreatePrescription returned invalid result: %v", result.Errors)
	}
	if prescription == nil {
		t.Fatal("CreatePrescription returned nil prescription")
	}
	if prescription.ID != anotherValidUUID() {
		t.Errorf("prescription.ID = %q, want %q", prescription.ID, anotherValidUUID())
	}
	if prescription.LiftID != validUUID() {
		t.Errorf("prescription.LiftID = %q, want %q", prescription.LiftID, validUUID())
	}
	if prescription.LoadStrategy == nil {
		t.Error("prescription.LoadStrategy should not be nil")
	}
	if prescription.SetScheme == nil {
		t.Error("prescription.SetScheme should not be nil")
	}
	if prescription.Order != 1 {
		t.Errorf("prescription.Order = %d, want %d", prescription.Order, 1)
	}
	if prescription.Notes != "Focus on form" {
		t.Errorf("prescription.Notes = %q, want %q", prescription.Notes, "Focus on form")
	}
	if prescription.RestSeconds == nil || *prescription.RestSeconds != 90 {
		t.Errorf("prescription.RestSeconds = %v, want 90", prescription.RestSeconds)
	}
}

func TestCreatePrescription_MinimalValidInput(t *testing.T) {
	input := CreatePrescriptionInput{
		LiftID:       validUUID(),
		LoadStrategy: validLoadStrategy(),
		SetScheme:    validSetScheme(),
		// Order defaults to 0, Notes empty, RestSeconds nil
	}

	prescription, result := CreatePrescription(input, anotherValidUUID())

	if !result.Valid {
		t.Errorf("CreatePrescription returned invalid result: %v", result.Errors)
	}
	if prescription == nil {
		t.Fatal("CreatePrescription returned nil prescription")
	}
	if prescription.Order != 0 {
		t.Errorf("prescription.Order = %d, want %d", prescription.Order, 0)
	}
	if prescription.Notes != "" {
		t.Errorf("prescription.Notes = %q, want empty", prescription.Notes)
	}
	if prescription.RestSeconds != nil {
		t.Errorf("prescription.RestSeconds = %v, want nil", prescription.RestSeconds)
	}
}

func TestCreatePrescription_InvalidLiftID(t *testing.T) {
	input := CreatePrescriptionInput{
		LiftID:       "",
		LoadStrategy: validLoadStrategy(),
		SetScheme:    validSetScheme(),
	}

	prescription, result := CreatePrescription(input, anotherValidUUID())

	if result.Valid {
		t.Error("CreatePrescription with empty LiftID should be invalid")
	}
	if prescription != nil {
		t.Error("CreatePrescription with invalid input should return nil")
	}
}

func TestCreatePrescription_InvalidLoadStrategy(t *testing.T) {
	input := CreatePrescriptionInput{
		LiftID:       validUUID(),
		LoadStrategy: nil,
		SetScheme:    validSetScheme(),
	}

	prescription, result := CreatePrescription(input, anotherValidUUID())

	if result.Valid {
		t.Error("CreatePrescription with nil LoadStrategy should be invalid")
	}
	if prescription != nil {
		t.Error("CreatePrescription with invalid input should return nil")
	}
}

func TestCreatePrescription_InvalidSetScheme(t *testing.T) {
	input := CreatePrescriptionInput{
		LiftID:       validUUID(),
		LoadStrategy: validLoadStrategy(),
		SetScheme:    nil,
	}

	prescription, result := CreatePrescription(input, anotherValidUUID())

	if result.Valid {
		t.Error("CreatePrescription with nil SetScheme should be invalid")
	}
	if prescription != nil {
		t.Error("CreatePrescription with invalid input should return nil")
	}
}

func TestCreatePrescription_InvalidOrder(t *testing.T) {
	input := CreatePrescriptionInput{
		LiftID:       validUUID(),
		LoadStrategy: validLoadStrategy(),
		SetScheme:    validSetScheme(),
		Order:        -1,
	}

	prescription, result := CreatePrescription(input, anotherValidUUID())

	if result.Valid {
		t.Error("CreatePrescription with negative Order should be invalid")
	}
	if prescription != nil {
		t.Error("CreatePrescription with invalid input should return nil")
	}
}

func TestCreatePrescription_InvalidNotes(t *testing.T) {
	input := CreatePrescriptionInput{
		LiftID:       validUUID(),
		LoadStrategy: validLoadStrategy(),
		SetScheme:    validSetScheme(),
		Notes:        strings.Repeat("a", 501),
	}

	prescription, result := CreatePrescription(input, anotherValidUUID())

	if result.Valid {
		t.Error("CreatePrescription with notes > 500 chars should be invalid")
	}
	if prescription != nil {
		t.Error("CreatePrescription with invalid input should return nil")
	}
}

func TestCreatePrescription_InvalidRestSeconds(t *testing.T) {
	input := CreatePrescriptionInput{
		LiftID:       validUUID(),
		LoadStrategy: validLoadStrategy(),
		SetScheme:    validSetScheme(),
		RestSeconds:  intPtr(-1),
	}

	prescription, result := CreatePrescription(input, anotherValidUUID())

	if result.Valid {
		t.Error("CreatePrescription with negative RestSeconds should be invalid")
	}
	if prescription != nil {
		t.Error("CreatePrescription with invalid input should return nil")
	}
}

func TestCreatePrescription_MultipleErrors(t *testing.T) {
	input := CreatePrescriptionInput{
		LiftID:       "",
		LoadStrategy: nil,
		SetScheme:    nil,
		Order:        -1,
		Notes:        strings.Repeat("a", 501),
		RestSeconds:  intPtr(-1),
	}

	prescription, result := CreatePrescription(input, anotherValidUUID())

	if result.Valid {
		t.Error("CreatePrescription with multiple invalid fields should be invalid")
	}
	if prescription != nil {
		t.Error("CreatePrescription with invalid input should return nil")
	}
	if len(result.Errors) < 6 {
		t.Errorf("Expected at least 6 errors, got %d", len(result.Errors))
	}
}

func TestCreatePrescription_SetsTimestamps(t *testing.T) {
	input := CreatePrescriptionInput{
		LiftID:       validUUID(),
		LoadStrategy: validLoadStrategy(),
		SetScheme:    validSetScheme(),
	}

	before := time.Now()
	prescription, result := CreatePrescription(input, anotherValidUUID())
	after := time.Now()

	if !result.Valid {
		t.Errorf("CreatePrescription returned invalid result: %v", result.Errors)
	}

	if prescription.CreatedAt.Before(before) || prescription.CreatedAt.After(after) {
		t.Errorf("CreatedAt should be between %v and %v, got %v", before, after, prescription.CreatedAt)
	}
	if prescription.UpdatedAt.Before(before) || prescription.UpdatedAt.After(after) {
		t.Errorf("UpdatedAt should be between %v and %v, got %v", before, after, prescription.UpdatedAt)
	}
}

// ==================== UpdatePrescription Tests ====================

func TestUpdatePrescription_UpdateLiftID(t *testing.T) {
	prescription := createValidPrescription()
	newLiftID := anotherValidUUID()

	input := UpdatePrescriptionInput{LiftID: &newLiftID}
	result := UpdatePrescription(prescription, input)

	if !result.Valid {
		t.Errorf("UpdatePrescription returned invalid result: %v", result.Errors)
	}
	if prescription.LiftID != newLiftID {
		t.Errorf("prescription.LiftID = %q, want %q", prescription.LiftID, newLiftID)
	}
}

func TestUpdatePrescription_UpdateLoadStrategy(t *testing.T) {
	prescription := createValidPrescription()
	newStrategy := &mockLoadStrategy{
		strategyType:    loadstrategy.TypeFixedWeight,
		calculateResult: 135.0,
	}

	input := UpdatePrescriptionInput{LoadStrategy: newStrategy}
	result := UpdatePrescription(prescription, input)

	if !result.Valid {
		t.Errorf("UpdatePrescription returned invalid result: %v", result.Errors)
	}
	if prescription.LoadStrategy.Type() != loadstrategy.TypeFixedWeight {
		t.Errorf("prescription.LoadStrategy.Type() = %q, want %q", prescription.LoadStrategy.Type(), loadstrategy.TypeFixedWeight)
	}
}

func TestUpdatePrescription_UpdateSetScheme(t *testing.T) {
	prescription := createValidPrescription()
	newScheme := &mockSetScheme{
		schemeType: setscheme.TypeRamp,
		generateResult: []setscheme.GeneratedSet{
			{SetNumber: 1, Weight: 135.0, TargetReps: 5, IsWorkSet: false},
			{SetNumber: 2, Weight: 185.0, TargetReps: 3, IsWorkSet: true},
		},
	}

	input := UpdatePrescriptionInput{SetScheme: newScheme}
	result := UpdatePrescription(prescription, input)

	if !result.Valid {
		t.Errorf("UpdatePrescription returned invalid result: %v", result.Errors)
	}
	if prescription.SetScheme.Type() != setscheme.TypeRamp {
		t.Errorf("prescription.SetScheme.Type() = %q, want %q", prescription.SetScheme.Type(), setscheme.TypeRamp)
	}
}

func TestUpdatePrescription_UpdateOrder(t *testing.T) {
	prescription := createValidPrescription()
	newOrder := 5

	input := UpdatePrescriptionInput{Order: &newOrder}
	result := UpdatePrescription(prescription, input)

	if !result.Valid {
		t.Errorf("UpdatePrescription returned invalid result: %v", result.Errors)
	}
	if prescription.Order != 5 {
		t.Errorf("prescription.Order = %d, want %d", prescription.Order, 5)
	}
}

func TestUpdatePrescription_UpdateNotes(t *testing.T) {
	prescription := createValidPrescription()
	newNotes := "Updated notes"

	input := UpdatePrescriptionInput{Notes: &newNotes}
	result := UpdatePrescription(prescription, input)

	if !result.Valid {
		t.Errorf("UpdatePrescription returned invalid result: %v", result.Errors)
	}
	if prescription.Notes != "Updated notes" {
		t.Errorf("prescription.Notes = %q, want %q", prescription.Notes, "Updated notes")
	}
}

func TestUpdatePrescription_UpdateRestSeconds(t *testing.T) {
	prescription := createValidPrescription()
	newRestSeconds := 120

	input := UpdatePrescriptionInput{RestSeconds: &newRestSeconds}
	result := UpdatePrescription(prescription, input)

	if !result.Valid {
		t.Errorf("UpdatePrescription returned invalid result: %v", result.Errors)
	}
	if prescription.RestSeconds == nil || *prescription.RestSeconds != 120 {
		t.Errorf("prescription.RestSeconds = %v, want 120", prescription.RestSeconds)
	}
}

func TestUpdatePrescription_ClearRestSeconds(t *testing.T) {
	prescription := createValidPrescription()
	prescription.RestSeconds = intPtr(90)

	input := UpdatePrescriptionInput{ClearRestSeconds: true}
	result := UpdatePrescription(prescription, input)

	if !result.Valid {
		t.Errorf("UpdatePrescription returned invalid result: %v", result.Errors)
	}
	if prescription.RestSeconds != nil {
		t.Errorf("prescription.RestSeconds = %v, want nil", prescription.RestSeconds)
	}
}

func TestUpdatePrescription_InvalidLiftID(t *testing.T) {
	prescription := createValidPrescription()
	originalLiftID := prescription.LiftID
	invalidLiftID := "not-a-uuid"

	input := UpdatePrescriptionInput{LiftID: &invalidLiftID}
	result := UpdatePrescription(prescription, input)

	if result.Valid {
		t.Error("UpdatePrescription with invalid LiftID should be invalid")
	}
	if prescription.LiftID != originalLiftID {
		t.Error("prescription.LiftID should not change on validation failure")
	}
}

func TestUpdatePrescription_InvalidLoadStrategy(t *testing.T) {
	prescription := createValidPrescription()
	originalStrategy := prescription.LoadStrategy
	invalidStrategy := &mockLoadStrategy{validateErr: errors.New("invalid")}

	input := UpdatePrescriptionInput{LoadStrategy: invalidStrategy}
	result := UpdatePrescription(prescription, input)

	if result.Valid {
		t.Error("UpdatePrescription with invalid LoadStrategy should be invalid")
	}
	if prescription.LoadStrategy != originalStrategy {
		t.Error("prescription.LoadStrategy should not change on validation failure")
	}
}

func TestUpdatePrescription_InvalidSetScheme(t *testing.T) {
	prescription := createValidPrescription()
	originalScheme := prescription.SetScheme
	invalidScheme := &mockSetScheme{validateErr: errors.New("invalid")}

	input := UpdatePrescriptionInput{SetScheme: invalidScheme}
	result := UpdatePrescription(prescription, input)

	if result.Valid {
		t.Error("UpdatePrescription with invalid SetScheme should be invalid")
	}
	if prescription.SetScheme != originalScheme {
		t.Error("prescription.SetScheme should not change on validation failure")
	}
}

func TestUpdatePrescription_InvalidOrder(t *testing.T) {
	prescription := createValidPrescription()
	originalOrder := prescription.Order
	invalidOrder := -5

	input := UpdatePrescriptionInput{Order: &invalidOrder}
	result := UpdatePrescription(prescription, input)

	if result.Valid {
		t.Error("UpdatePrescription with invalid Order should be invalid")
	}
	if prescription.Order != originalOrder {
		t.Error("prescription.Order should not change on validation failure")
	}
}

func TestUpdatePrescription_InvalidNotes(t *testing.T) {
	prescription := createValidPrescription()
	originalNotes := prescription.Notes
	invalidNotes := strings.Repeat("a", 501)

	input := UpdatePrescriptionInput{Notes: &invalidNotes}
	result := UpdatePrescription(prescription, input)

	if result.Valid {
		t.Error("UpdatePrescription with invalid Notes should be invalid")
	}
	if prescription.Notes != originalNotes {
		t.Error("prescription.Notes should not change on validation failure")
	}
}

func TestUpdatePrescription_InvalidRestSeconds(t *testing.T) {
	prescription := createValidPrescription()
	originalRestSeconds := prescription.RestSeconds
	invalidRestSeconds := -10

	input := UpdatePrescriptionInput{RestSeconds: &invalidRestSeconds}
	result := UpdatePrescription(prescription, input)

	if result.Valid {
		t.Error("UpdatePrescription with invalid RestSeconds should be invalid")
	}
	if prescription.RestSeconds != originalRestSeconds {
		t.Error("prescription.RestSeconds should not change on validation failure")
	}
}

func TestUpdatePrescription_NoChanges(t *testing.T) {
	prescription := createValidPrescription()

	input := UpdatePrescriptionInput{}
	result := UpdatePrescription(prescription, input)

	if !result.Valid {
		t.Errorf("UpdatePrescription with no changes should be valid: %v", result.Errors)
	}
}

func TestUpdatePrescription_UpdatesTimestamp(t *testing.T) {
	prescription := createValidPrescription()
	originalUpdatedAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	prescription.UpdatedAt = originalUpdatedAt

	newOrder := 10
	input := UpdatePrescriptionInput{Order: &newOrder}

	before := time.Now()
	result := UpdatePrescription(prescription, input)
	after := time.Now()

	if !result.Valid {
		t.Errorf("UpdatePrescription returned invalid result: %v", result.Errors)
	}
	if prescription.UpdatedAt.Before(before) || prescription.UpdatedAt.After(after) {
		t.Errorf("UpdatedAt should be between %v and %v, got %v", before, after, prescription.UpdatedAt)
	}
}

func TestUpdatePrescription_MultipleFields(t *testing.T) {
	prescription := createValidPrescription()
	newLiftID := anotherValidUUID()
	newOrder := 3
	newNotes := "Multiple updates"
	newRestSeconds := 60

	input := UpdatePrescriptionInput{
		LiftID:      &newLiftID,
		Order:       &newOrder,
		Notes:       &newNotes,
		RestSeconds: &newRestSeconds,
	}
	result := UpdatePrescription(prescription, input)

	if !result.Valid {
		t.Errorf("UpdatePrescription returned invalid result: %v", result.Errors)
	}
	if prescription.LiftID != newLiftID {
		t.Errorf("prescription.LiftID = %q, want %q", prescription.LiftID, newLiftID)
	}
	if prescription.Order != 3 {
		t.Errorf("prescription.Order = %d, want %d", prescription.Order, 3)
	}
	if prescription.Notes != "Multiple updates" {
		t.Errorf("prescription.Notes = %q, want %q", prescription.Notes, "Multiple updates")
	}
	if prescription.RestSeconds == nil || *prescription.RestSeconds != 60 {
		t.Errorf("prescription.RestSeconds = %v, want 60", prescription.RestSeconds)
	}
}

func createValidPrescription() *Prescription {
	prescription, _ := CreatePrescription(CreatePrescriptionInput{
		LiftID:       validUUID(),
		LoadStrategy: validLoadStrategy(),
		SetScheme:    validSetScheme(),
		Order:        1,
		Notes:        "Test notes",
		RestSeconds:  intPtr(90),
	}, anotherValidUUID())
	return prescription
}

// ==================== Prescription.Validate Tests ====================

func TestPrescription_Validate_Valid(t *testing.T) {
	prescription := createValidPrescription()

	result := prescription.Validate()

	if !result.Valid {
		t.Errorf("Validate returned invalid for valid prescription: %v", result.Errors)
	}
}

func TestPrescription_Validate_InvalidLiftID(t *testing.T) {
	prescription := createValidPrescription()
	prescription.LiftID = ""

	result := prescription.Validate()

	if result.Valid {
		t.Error("Validate should be invalid for empty LiftID")
	}
}

func TestPrescription_Validate_InvalidLoadStrategy(t *testing.T) {
	prescription := createValidPrescription()
	prescription.LoadStrategy = nil

	result := prescription.Validate()

	if result.Valid {
		t.Error("Validate should be invalid for nil LoadStrategy")
	}
}

func TestPrescription_Validate_InvalidSetScheme(t *testing.T) {
	prescription := createValidPrescription()
	prescription.SetScheme = nil

	result := prescription.Validate()

	if result.Valid {
		t.Error("Validate should be invalid for nil SetScheme")
	}
}

func TestPrescription_Validate_InvalidOrder(t *testing.T) {
	prescription := createValidPrescription()
	prescription.Order = -1

	result := prescription.Validate()

	if result.Valid {
		t.Error("Validate should be invalid for negative Order")
	}
}

func TestPrescription_Validate_InvalidNotes(t *testing.T) {
	prescription := createValidPrescription()
	prescription.Notes = strings.Repeat("a", 501)

	result := prescription.Validate()

	if result.Valid {
		t.Error("Validate should be invalid for notes > 500 chars")
	}
}

func TestPrescription_Validate_InvalidRestSeconds(t *testing.T) {
	prescription := createValidPrescription()
	negativeRest := -10
	prescription.RestSeconds = &negativeRest

	result := prescription.Validate()

	if result.Valid {
		t.Error("Validate should be invalid for negative RestSeconds")
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
}

func TestValidationResult_AddError(t *testing.T) {
	result := NewValidationResult()
	result.AddError(ErrLiftIDRequired)

	if result.Valid {
		t.Error("AddError should mark result as invalid")
	}
	if len(result.Errors) != 1 {
		t.Errorf("AddError should add one error, got %d", len(result.Errors))
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
	result.AddError(ErrLiftIDRequired)
	result.AddError(ErrLoadStrategyRequired)

	err := result.Error()
	if err == nil {
		t.Error("Error() on invalid result should return error")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "lift ID") {
		t.Errorf("Error message should contain lift ID error: %s", errMsg)
	}
	if !strings.Contains(errMsg, "load strategy") {
		t.Errorf("Error message should contain load strategy error: %s", errMsg)
	}
}

// ==================== Resolve Tests ====================

func TestPrescription_Resolve_Success(t *testing.T) {
	prescription := createValidPrescription()

	liftLookup := newMockLiftLookup()
	liftLookup.SetLift(prescription.LiftID, &LiftInfo{
		ID:   prescription.LiftID,
		Name: "Back Squat",
		Slug: "back-squat",
	})

	ctx := context.Background()
	resCtx := DefaultResolutionContext(liftLookup)

	resolved, err := prescription.Resolve(ctx, "user-123", resCtx)

	if err != nil {
		t.Errorf("Resolve returned error: %v", err)
	}
	if resolved == nil {
		t.Fatal("Resolve returned nil ResolvedPrescription")
	}
	if resolved.PrescriptionID != prescription.ID {
		t.Errorf("resolved.PrescriptionID = %q, want %q", resolved.PrescriptionID, prescription.ID)
	}
	if resolved.Lift.ID != prescription.LiftID {
		t.Errorf("resolved.Lift.ID = %q, want %q", resolved.Lift.ID, prescription.LiftID)
	}
	if resolved.Lift.Name != "Back Squat" {
		t.Errorf("resolved.Lift.Name = %q, want %q", resolved.Lift.Name, "Back Squat")
	}
	if resolved.Lift.Slug != "back-squat" {
		t.Errorf("resolved.Lift.Slug = %q, want %q", resolved.Lift.Slug, "back-squat")
	}
	if len(resolved.Sets) != 3 {
		t.Errorf("resolved.Sets length = %d, want %d", len(resolved.Sets), 3)
	}
	if resolved.Notes != prescription.Notes {
		t.Errorf("resolved.Notes = %q, want %q", resolved.Notes, prescription.Notes)
	}
	if resolved.RestSeconds == nil || *resolved.RestSeconds != *prescription.RestSeconds {
		t.Errorf("resolved.RestSeconds = %v, want %v", resolved.RestSeconds, prescription.RestSeconds)
	}
}

func TestPrescription_Resolve_EmptyUserID(t *testing.T) {
	prescription := createValidPrescription()

	ctx := context.Background()
	resCtx := DefaultResolutionContext(nil)

	_, err := prescription.Resolve(ctx, "", resCtx)

	if err == nil {
		t.Error("Resolve with empty userID should return error")
	}
	if !strings.Contains(err.Error(), "user ID") {
		t.Errorf("Error should mention user ID: %v", err)
	}
}

func TestPrescription_Resolve_LiftNotFound(t *testing.T) {
	prescription := createValidPrescription()

	liftLookup := newMockLiftLookup()
	// Don't add the lift - it won't be found

	ctx := context.Background()
	resCtx := DefaultResolutionContext(liftLookup)

	_, err := prescription.Resolve(ctx, "user-123", resCtx)

	if err == nil {
		t.Error("Resolve with non-existent lift should return error")
	}
	if !errors.Is(err, ErrLiftNotFound) {
		t.Errorf("Error should be ErrLiftNotFound, got: %v", err)
	}
}

func TestPrescription_Resolve_LiftLookupError(t *testing.T) {
	prescription := createValidPrescription()

	liftLookup := newMockLiftLookup()
	liftLookup.SetError(errors.New("database connection failed"))

	ctx := context.Background()
	resCtx := DefaultResolutionContext(liftLookup)

	_, err := prescription.Resolve(ctx, "user-123", resCtx)

	if err == nil {
		t.Error("Resolve with lift lookup error should return error")
	}
	if !strings.Contains(err.Error(), "failed to look up lift") {
		t.Errorf("Error should mention lift lookup failure: %v", err)
	}
}

func TestPrescription_Resolve_LoadCalculationError(t *testing.T) {
	strategy := &mockLoadStrategy{
		strategyType: loadstrategy.TypePercentOf,
		calculateErr: errors.New("calculation failed"),
	}

	prescription := &Prescription{
		ID:           anotherValidUUID(),
		LiftID:       validUUID(),
		LoadStrategy: strategy,
		SetScheme:    validSetScheme(),
	}

	liftLookup := newMockLiftLookup()
	liftLookup.SetLift(prescription.LiftID, &LiftInfo{
		ID:   prescription.LiftID,
		Name: "Back Squat",
		Slug: "back-squat",
	})

	ctx := context.Background()
	resCtx := DefaultResolutionContext(liftLookup)

	_, err := prescription.Resolve(ctx, "user-123", resCtx)

	if err == nil {
		t.Error("Resolve with load calculation error should return error")
	}
	if !strings.Contains(err.Error(), "failed to calculate load") {
		t.Errorf("Error should mention load calculation failure: %v", err)
	}
}

func TestPrescription_Resolve_MaxNotFound(t *testing.T) {
	strategy := &mockLoadStrategy{
		strategyType: loadstrategy.TypePercentOf,
		calculateErr: loadstrategy.ErrMaxNotFound,
	}

	prescription := &Prescription{
		ID:           anotherValidUUID(),
		LiftID:       validUUID(),
		LoadStrategy: strategy,
		SetScheme:    validSetScheme(),
	}

	liftLookup := newMockLiftLookup()
	liftLookup.SetLift(prescription.LiftID, &LiftInfo{
		ID:   prescription.LiftID,
		Name: "Back Squat",
		Slug: "back-squat",
	})

	ctx := context.Background()
	resCtx := DefaultResolutionContext(liftLookup)

	_, err := prescription.Resolve(ctx, "user-123", resCtx)

	if err == nil {
		t.Error("Resolve when max not found should return error")
	}
	if !errors.Is(err, ErrMaxNotFound) {
		t.Errorf("Error should be ErrMaxNotFound, got: %v", err)
	}
}

func TestPrescription_Resolve_SetGenerationError(t *testing.T) {
	scheme := &mockSetScheme{
		schemeType:  setscheme.TypeFixed,
		generateErr: errors.New("set generation failed"),
	}

	prescription := &Prescription{
		ID:           anotherValidUUID(),
		LiftID:       validUUID(),
		LoadStrategy: validLoadStrategy(),
		SetScheme:    scheme,
	}

	liftLookup := newMockLiftLookup()
	liftLookup.SetLift(prescription.LiftID, &LiftInfo{
		ID:   prescription.LiftID,
		Name: "Back Squat",
		Slug: "back-squat",
	})

	ctx := context.Background()
	resCtx := DefaultResolutionContext(liftLookup)

	_, err := prescription.Resolve(ctx, "user-123", resCtx)

	if err == nil {
		t.Error("Resolve with set generation error should return error")
	}
	if !strings.Contains(err.Error(), "failed to generate sets") {
		t.Errorf("Error should mention set generation failure: %v", err)
	}
}

func TestPrescription_Resolve_NoLiftLookup(t *testing.T) {
	prescription := createValidPrescription()

	ctx := context.Background()
	resCtx := ResolutionContext{
		LiftLookup:    nil,
		SetGenContext: setscheme.DefaultSetGenerationContext(),
	}

	resolved, err := prescription.Resolve(ctx, "user-123", resCtx)

	if err != nil {
		t.Errorf("Resolve without lift lookup should not return error: %v", err)
	}
	if resolved == nil {
		t.Fatal("Resolve returned nil ResolvedPrescription")
	}
	// Should use minimal lift info
	if resolved.Lift.ID != prescription.LiftID {
		t.Errorf("resolved.Lift.ID = %q, want %q", resolved.Lift.ID, prescription.LiftID)
	}
	if resolved.Lift.Name != "" {
		t.Errorf("resolved.Lift.Name = %q, want empty", resolved.Lift.Name)
	}
}

func TestPrescription_Resolve_WhitespaceUserID(t *testing.T) {
	prescription := createValidPrescription()

	ctx := context.Background()
	resCtx := DefaultResolutionContext(nil)

	_, err := prescription.Resolve(ctx, "   ", resCtx)

	if err == nil {
		t.Error("Resolve with whitespace-only userID should return error")
	}
}

func TestPrescription_Resolve_SetsContainCorrectValues(t *testing.T) {
	strategy := &mockLoadStrategy{
		strategyType:    loadstrategy.TypePercentOf,
		calculateResult: 200.0,
	}
	scheme := &mockSetScheme{
		schemeType: setscheme.TypeFixed,
		generateResult: []setscheme.GeneratedSet{
			{SetNumber: 1, Weight: 200.0, TargetReps: 5, IsWorkSet: true},
			{SetNumber: 2, Weight: 200.0, TargetReps: 5, IsWorkSet: true},
		},
	}

	prescription := &Prescription{
		ID:           anotherValidUUID(),
		LiftID:       validUUID(),
		LoadStrategy: strategy,
		SetScheme:    scheme,
	}

	liftLookup := newMockLiftLookup()
	liftLookup.SetLift(prescription.LiftID, &LiftInfo{
		ID:   prescription.LiftID,
		Name: "Bench Press",
		Slug: "bench-press",
	})

	ctx := context.Background()
	resCtx := DefaultResolutionContext(liftLookup)

	resolved, err := prescription.Resolve(ctx, "user-123", resCtx)

	if err != nil {
		t.Errorf("Resolve returned error: %v", err)
	}

	if len(resolved.Sets) != 2 {
		t.Fatalf("resolved.Sets length = %d, want 2", len(resolved.Sets))
	}

	// Verify first set
	if resolved.Sets[0].SetNumber != 1 {
		t.Errorf("Set 1 SetNumber = %d, want 1", resolved.Sets[0].SetNumber)
	}
	if resolved.Sets[0].Weight != 200.0 {
		t.Errorf("Set 1 Weight = %v, want 200.0", resolved.Sets[0].Weight)
	}
	if resolved.Sets[0].TargetReps != 5 {
		t.Errorf("Set 1 TargetReps = %d, want 5", resolved.Sets[0].TargetReps)
	}
	if !resolved.Sets[0].IsWorkSet {
		t.Error("Set 1 should be work set")
	}
}

// ==================== DefaultResolutionContext Tests ====================

func TestDefaultResolutionContext(t *testing.T) {
	liftLookup := newMockLiftLookup()
	resCtx := DefaultResolutionContext(liftLookup)

	if resCtx.LiftLookup != liftLookup {
		t.Error("LiftLookup should be set")
	}
	if resCtx.SetGenContext.WorkSetThreshold != 80.0 {
		t.Errorf("WorkSetThreshold = %v, want 80.0", resCtx.SetGenContext.WorkSetThreshold)
	}
}

// ==================== Constants Tests ====================

func TestMaxNotesLength(t *testing.T) {
	if MaxNotesLength != 500 {
		t.Errorf("MaxNotesLength = %d, want %d", MaxNotesLength, 500)
	}
}

// ==================== Error Messages Tests ====================

func TestErrorMessages(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		contains string
	}{
		{"ErrLiftIDRequired", ErrLiftIDRequired, "lift ID is required"},
		{"ErrLiftIDInvalid", ErrLiftIDInvalid, "valid UUID"},
		{"ErrLoadStrategyRequired", ErrLoadStrategyRequired, "load strategy is required"},
		{"ErrSetSchemeRequired", ErrSetSchemeRequired, "set scheme is required"},
		{"ErrOrderNegative", ErrOrderNegative, "order must be >= 0"},
		{"ErrNotesTooLong", ErrNotesTooLong, "500 characters"},
		{"ErrRestSecondsNegative", ErrRestSecondsNegative, "rest seconds must be >= 0"},
		{"ErrLiftNotFound", ErrLiftNotFound, "lift not found"},
		{"ErrMaxNotFound", ErrMaxNotFound, "max not found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(tt.err.Error(), tt.contains) {
				t.Errorf("%s.Error() = %q, should contain %q", tt.name, tt.err.Error(), tt.contains)
			}
		})
	}
}

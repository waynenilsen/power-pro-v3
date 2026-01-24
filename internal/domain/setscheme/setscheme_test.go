package setscheme

import (
	"encoding/json"
	"errors"
	"testing"
)

// mockScheme implements SetScheme for testing purposes.
type mockScheme struct {
	schemeType SetSchemeType
	sets       []GeneratedSet
	err        error
	valid      bool
}

func (m *mockScheme) Type() SetSchemeType {
	return m.schemeType
}

func (m *mockScheme) GenerateSets(baseWeight float64, context SetGenerationContext) ([]GeneratedSet, error) {
	if m.err != nil {
		return nil, m.err
	}
	// Apply baseWeight to sets for testing
	result := make([]GeneratedSet, len(m.sets))
	for i, s := range m.sets {
		result[i] = GeneratedSet{
			SetNumber:  s.SetNumber,
			Weight:     baseWeight * (s.Weight / 100), // Interpret Weight as percentage
			TargetReps: s.TargetReps,
			IsWorkSet:  s.IsWorkSet,
		}
	}
	return result, nil
}

func (m *mockScheme) Validate() error {
	if !m.valid {
		return errors.New("mock validation error")
	}
	return nil
}

func (m *mockScheme) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type SetSchemeType `json:"type"`
		Sets int           `json:"sets"`
		Reps int           `json:"reps"`
	}{
		Type: m.schemeType,
		Sets: len(m.sets),
		Reps: 5,
	})
}

// TestSetSchemeTypeConstants verifies all scheme type constants are defined.
func TestSetSchemeTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant SetSchemeType
		expected string
	}{
		{"TypeFixed", TypeFixed, "FIXED"},
		{"TypeRamp", TypeRamp, "RAMP"},
		{"TypeAMRAP", TypeAMRAP, "AMRAP"},
		{"TypeTopBackoff", TypeTopBackoff, "TOP_BACKOFF"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.constant) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}

// TestValidSchemeTypes verifies all scheme types are in the valid set.
func TestValidSchemeTypes(t *testing.T) {
	expectedTypes := []SetSchemeType{
		TypeFixed,
		TypeRamp,
		TypeAMRAP,
		TypeTopBackoff,
		TypeRepRange,
		TypeGreySkull,
	}

	for _, schemeType := range expectedTypes {
		if !ValidSchemeTypes[schemeType] {
			t.Errorf("expected %s to be in ValidSchemeTypes", schemeType)
		}
	}

	if len(ValidSchemeTypes) != len(expectedTypes) {
		t.Errorf("expected %d types in ValidSchemeTypes, got %d", len(expectedTypes), len(ValidSchemeTypes))
	}
}

// TestValidateSchemeType tests scheme type validation.
func TestValidateSchemeType(t *testing.T) {
	tests := []struct {
		name       string
		schemeType SetSchemeType
		wantErr    bool
	}{
		{"valid FIXED", TypeFixed, false},
		{"valid RAMP", TypeRamp, false},
		{"valid AMRAP", TypeAMRAP, false},
		{"valid TOP_BACKOFF", TypeTopBackoff, false},
		{"valid REP_RANGE", TypeRepRange, false},
		{"empty type", "", true},
		{"unknown type", "UNKNOWN_TYPE", true},
		{"lowercase type", "fixed", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSchemeType(tt.schemeType)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if !errors.Is(err, ErrUnknownSchemeType) {
					t.Errorf("expected ErrUnknownSchemeType, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// TestGeneratedSet tests the GeneratedSet struct.
func TestGeneratedSet(t *testing.T) {
	set := GeneratedSet{
		SetNumber:  1,
		Weight:     265.0,
		TargetReps: 5,
		IsWorkSet:  true,
	}

	if set.SetNumber != 1 {
		t.Errorf("expected SetNumber 1, got %d", set.SetNumber)
	}
	if set.Weight != 265.0 {
		t.Errorf("expected Weight 265.0, got %f", set.Weight)
	}
	if set.TargetReps != 5 {
		t.Errorf("expected TargetReps 5, got %d", set.TargetReps)
	}
	if !set.IsWorkSet {
		t.Error("expected IsWorkSet to be true")
	}
}

// TestGeneratedSet_JSON tests JSON serialization of GeneratedSet.
func TestGeneratedSet_JSON(t *testing.T) {
	set := GeneratedSet{
		SetNumber:  2,
		Weight:     315.0,
		TargetReps: 3,
		IsWorkSet:  true,
	}

	data, err := json.Marshal(set)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var parsed GeneratedSet
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if parsed.SetNumber != set.SetNumber {
		t.Errorf("expected SetNumber %d, got %d", set.SetNumber, parsed.SetNumber)
	}
	if parsed.Weight != set.Weight {
		t.Errorf("expected Weight %f, got %f", set.Weight, parsed.Weight)
	}
	if parsed.TargetReps != set.TargetReps {
		t.Errorf("expected TargetReps %d, got %d", set.TargetReps, parsed.TargetReps)
	}
	if parsed.IsWorkSet != set.IsWorkSet {
		t.Errorf("expected IsWorkSet %v, got %v", set.IsWorkSet, parsed.IsWorkSet)
	}
}

// TestSetGenerationContext tests the SetGenerationContext struct.
func TestSetGenerationContext(t *testing.T) {
	ctx := SetGenerationContext{
		WorkSetThreshold: 75.0,
	}

	if ctx.WorkSetThreshold != 75.0 {
		t.Errorf("expected WorkSetThreshold 75.0, got %f", ctx.WorkSetThreshold)
	}
}

// TestDefaultSetGenerationContext tests the default context factory.
func TestDefaultSetGenerationContext(t *testing.T) {
	ctx := DefaultSetGenerationContext()

	if ctx.WorkSetThreshold != 80.0 {
		t.Errorf("expected default WorkSetThreshold 80.0, got %f", ctx.WorkSetThreshold)
	}
}

// TestSchemeFactory_Register tests scheme registration.
func TestSchemeFactory_Register(t *testing.T) {
	factory := NewSchemeFactory()

	// Verify factory starts empty
	if factory.IsRegistered(TypeFixed) {
		t.Error("factory should not have TypeFixed registered initially")
	}

	// Register a scheme
	factory.Register(TypeFixed, func(data json.RawMessage) (SetScheme, error) {
		return &mockScheme{schemeType: TypeFixed, valid: true}, nil
	})

	// Verify it's registered
	if !factory.IsRegistered(TypeFixed) {
		t.Error("TypeFixed should be registered")
	}

	// Verify other types are not registered
	if factory.IsRegistered(TypeRamp) {
		t.Error("TypeRamp should not be registered")
	}
}

// TestSchemeFactory_Create tests scheme creation.
func TestSchemeFactory_Create(t *testing.T) {
	factory := NewSchemeFactory()
	factory.Register(TypeFixed, func(data json.RawMessage) (SetScheme, error) {
		return &mockScheme{
			schemeType: TypeFixed,
			sets: []GeneratedSet{
				{SetNumber: 1, Weight: 100, TargetReps: 5, IsWorkSet: true},
			},
			valid: true,
		}, nil
	})

	t.Run("create registered scheme", func(t *testing.T) {
		scheme, err := factory.Create(TypeFixed, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if scheme.Type() != TypeFixed {
			t.Errorf("expected type %s, got %s", TypeFixed, scheme.Type())
		}
	})

	t.Run("create unregistered scheme", func(t *testing.T) {
		_, err := factory.Create(TypeRamp, nil)
		if err == nil {
			t.Error("expected error for unregistered type")
		}
		if !errors.Is(err, ErrSchemeNotRegistered) {
			t.Errorf("expected ErrSchemeNotRegistered, got %v", err)
		}
	})
}

// TestSchemeFactory_CreateFromJSON tests scheme creation from JSON.
func TestSchemeFactory_CreateFromJSON(t *testing.T) {
	factory := NewSchemeFactory()
	factory.Register(TypeFixed, func(data json.RawMessage) (SetScheme, error) {
		var parsed struct {
			Type SetSchemeType `json:"type"`
			Sets int           `json:"sets"`
			Reps int           `json:"reps"`
		}
		if err := json.Unmarshal(data, &parsed); err != nil {
			return nil, err
		}
		sets := make([]GeneratedSet, parsed.Sets)
		for i := range sets {
			sets[i] = GeneratedSet{
				SetNumber:  i + 1,
				Weight:     100,
				TargetReps: parsed.Reps,
				IsWorkSet:  true,
			}
		}
		return &mockScheme{
			schemeType: parsed.Type,
			sets:       sets,
			valid:      true,
		}, nil
	})

	t.Run("valid JSON with type", func(t *testing.T) {
		jsonData := []byte(`{"type": "FIXED", "sets": 5, "reps": 5}`)
		scheme, err := factory.CreateFromJSON(jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if scheme.Type() != TypeFixed {
			t.Errorf("expected type %s, got %s", TypeFixed, scheme.Type())
		}
	})

	t.Run("unregistered type in JSON", func(t *testing.T) {
		jsonData := []byte(`{"type": "RAMP", "steps": []}`)
		_, err := factory.CreateFromJSON(jsonData)
		if err == nil {
			t.Error("expected error for unregistered type")
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		jsonData := []byte(`{invalid json}`)
		_, err := factory.CreateFromJSON(jsonData)
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("missing type field", func(t *testing.T) {
		jsonData := []byte(`{"sets": 5, "reps": 5}`)
		_, err := factory.CreateFromJSON(jsonData)
		if err == nil {
			t.Error("expected error for missing type")
		}
	})
}

// TestSchemeFactory_RegisteredTypes tests retrieval of registered types.
func TestSchemeFactory_RegisteredTypes(t *testing.T) {
	factory := NewSchemeFactory()

	// Initially empty
	types := factory.RegisteredTypes()
	if len(types) != 0 {
		t.Errorf("expected 0 registered types, got %d", len(types))
	}

	// Register some types
	factory.Register(TypeFixed, func(data json.RawMessage) (SetScheme, error) {
		return &mockScheme{schemeType: TypeFixed}, nil
	})
	factory.Register(TypeRamp, func(data json.RawMessage) (SetScheme, error) {
		return &mockScheme{schemeType: TypeRamp}, nil
	})

	types = factory.RegisteredTypes()
	if len(types) != 2 {
		t.Errorf("expected 2 registered types, got %d", len(types))
	}

	// Verify both types are in the list
	hasFixed := false
	hasRamp := false
	for _, tt := range types {
		if tt == TypeFixed {
			hasFixed = true
		}
		if tt == TypeRamp {
			hasRamp = true
		}
	}
	if !hasFixed {
		t.Error("expected TypeFixed in registered types")
	}
	if !hasRamp {
		t.Error("expected TypeRamp in registered types")
	}
}

// TestSchemeEnvelope_UnmarshalJSON tests envelope unmarshaling.
func TestSchemeEnvelope_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name         string
		json         string
		expectedType SetSchemeType
		wantErr      bool
	}{
		{
			name:         "valid FIXED",
			json:         `{"type": "FIXED", "sets": 5, "reps": 5}`,
			expectedType: TypeFixed,
			wantErr:      false,
		},
		{
			name:         "valid RAMP",
			json:         `{"type": "RAMP", "steps": []}`,
			expectedType: TypeRamp,
			wantErr:      false,
		},
		{
			name:         "empty type",
			json:         `{"type": "", "sets": 5}`,
			expectedType: "",
			wantErr:      false, // Parsing succeeds, validation catches this
		},
		{
			name:    "invalid JSON",
			json:    `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var envelope SchemeEnvelope
			err := json.Unmarshal([]byte(tt.json), &envelope)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if envelope.Type != tt.expectedType {
					t.Errorf("expected type %s, got %s", tt.expectedType, envelope.Type)
				}
				if len(envelope.Raw) == 0 {
					t.Error("expected Raw to contain the original JSON")
				}
			}
		})
	}
}

// TestMarshalScheme tests scheme serialization.
func TestMarshalScheme(t *testing.T) {
	scheme := &mockScheme{
		schemeType: TypeFixed,
		sets: []GeneratedSet{
			{SetNumber: 1, Weight: 100, TargetReps: 5, IsWorkSet: true},
		},
		valid: true,
	}

	data, err := MarshalScheme(scheme)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the JSON contains the type field
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse marshaled JSON: %v", err)
	}

	typeVal, ok := parsed["type"]
	if !ok {
		t.Error("expected 'type' field in marshaled JSON")
	}
	if typeVal != string(TypeFixed) {
		t.Errorf("expected type %s, got %v", TypeFixed, typeVal)
	}
}

// TestSetScheme_Interface tests that the interface contract is correct.
func TestSetScheme_Interface(t *testing.T) {
	// Test that mockScheme implements SetScheme
	var _ SetScheme = &mockScheme{}

	scheme := &mockScheme{
		schemeType: TypeFixed,
		sets: []GeneratedSet{
			{SetNumber: 1, Weight: 100, TargetReps: 5, IsWorkSet: true},
			{SetNumber: 2, Weight: 100, TargetReps: 5, IsWorkSet: true},
			{SetNumber: 3, Weight: 100, TargetReps: 5, IsWorkSet: true},
		},
		valid: true,
	}

	t.Run("Type returns correct type", func(t *testing.T) {
		if scheme.Type() != TypeFixed {
			t.Errorf("expected %s, got %s", TypeFixed, scheme.Type())
		}
	})

	t.Run("GenerateSets returns sets", func(t *testing.T) {
		baseWeight := 265.0
		ctx := DefaultSetGenerationContext()
		sets, err := scheme.GenerateSets(baseWeight, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sets) != 3 {
			t.Errorf("expected 3 sets, got %d", len(sets))
		}
		for i, set := range sets {
			if set.SetNumber != i+1 {
				t.Errorf("expected SetNumber %d, got %d", i+1, set.SetNumber)
			}
			if set.Weight != baseWeight {
				t.Errorf("expected Weight %f, got %f", baseWeight, set.Weight)
			}
			if set.TargetReps != 5 {
				t.Errorf("expected TargetReps 5, got %d", set.TargetReps)
			}
			if !set.IsWorkSet {
				t.Error("expected IsWorkSet to be true")
			}
		}
	})

	t.Run("GenerateSets returns error when configured", func(t *testing.T) {
		errorScheme := &mockScheme{
			schemeType: TypeFixed,
			err:        ErrInvalidParams,
		}
		ctx := DefaultSetGenerationContext()
		_, err := errorScheme.GenerateSets(100, ctx)
		if !errors.Is(err, ErrInvalidParams) {
			t.Errorf("expected ErrInvalidParams, got %v", err)
		}
	})

	t.Run("Validate returns nil for valid scheme", func(t *testing.T) {
		if err := scheme.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Validate returns error for invalid scheme", func(t *testing.T) {
		invalidScheme := &mockScheme{valid: false}
		if err := invalidScheme.Validate(); err == nil {
			t.Error("expected validation error")
		}
	})
}

// TestErrors tests that error variables are defined correctly.
func TestErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		msg  string
	}{
		{"ErrUnknownSchemeType", ErrUnknownSchemeType, "unknown set scheme type"},
		{"ErrInvalidParams", ErrInvalidParams, "invalid set scheme parameters"},
		{"ErrSchemeNotRegistered", ErrSchemeNotRegistered, "scheme type not registered in factory"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.msg {
				t.Errorf("expected message %q, got %q", tt.msg, tt.err.Error())
			}
		})
	}
}

// TestNewSchemeFactory tests factory creation.
func TestNewSchemeFactory(t *testing.T) {
	factory := NewSchemeFactory()
	if factory == nil {
		t.Fatal("expected non-nil factory")
	}
	if factory.creators == nil {
		t.Error("expected creators map to be initialized")
	}
	if len(factory.creators) != 0 {
		t.Error("expected creators map to be empty initially")
	}
}

// TestSchemeFactory_OverwriteRegistration tests that re-registering overwrites.
func TestSchemeFactory_OverwriteRegistration(t *testing.T) {
	factory := NewSchemeFactory()

	// Register first version - 3 sets
	factory.Register(TypeFixed, func(data json.RawMessage) (SetScheme, error) {
		return &mockScheme{
			schemeType: TypeFixed,
			sets: []GeneratedSet{
				{SetNumber: 1, Weight: 100, TargetReps: 5, IsWorkSet: true},
				{SetNumber: 2, Weight: 100, TargetReps: 5, IsWorkSet: true},
				{SetNumber: 3, Weight: 100, TargetReps: 5, IsWorkSet: true},
			},
		}, nil
	})

	// Register second version - 5 sets (overwrite)
	factory.Register(TypeFixed, func(data json.RawMessage) (SetScheme, error) {
		return &mockScheme{
			schemeType: TypeFixed,
			sets: []GeneratedSet{
				{SetNumber: 1, Weight: 100, TargetReps: 5, IsWorkSet: true},
				{SetNumber: 2, Weight: 100, TargetReps: 5, IsWorkSet: true},
				{SetNumber: 3, Weight: 100, TargetReps: 5, IsWorkSet: true},
				{SetNumber: 4, Weight: 100, TargetReps: 5, IsWorkSet: true},
				{SetNumber: 5, Weight: 100, TargetReps: 5, IsWorkSet: true},
			},
		}, nil
	})

	// Create scheme and verify it uses the second version
	scheme, err := factory.Create(TypeFixed, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := DefaultSetGenerationContext()
	sets, err := scheme.GenerateSets(100, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sets) != 5 {
		t.Errorf("expected 5 sets from second registration, got %d", len(sets))
	}
}

// TestSetNumberIsOneIndexed verifies that SetNumber starts at 1.
func TestSetNumberIsOneIndexed(t *testing.T) {
	scheme := &mockScheme{
		schemeType: TypeFixed,
		sets: []GeneratedSet{
			{SetNumber: 1, Weight: 100, TargetReps: 5, IsWorkSet: true},
			{SetNumber: 2, Weight: 100, TargetReps: 5, IsWorkSet: true},
			{SetNumber: 3, Weight: 100, TargetReps: 5, IsWorkSet: true},
		},
		valid: true,
	}

	ctx := DefaultSetGenerationContext()
	sets, err := scheme.GenerateSets(100, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify sets are 1-indexed
	for i, set := range sets {
		expectedSetNumber := i + 1
		if set.SetNumber != expectedSetNumber {
			t.Errorf("set %d: expected SetNumber %d, got %d", i, expectedSetNumber, set.SetNumber)
		}
	}

	// Verify first set is not 0-indexed
	if len(sets) > 0 && sets[0].SetNumber == 0 {
		t.Error("SetNumber should be 1-indexed, not 0-indexed")
	}
}

// TestIsWorkSetFlag tests that IsWorkSet properly distinguishes warmup and work sets.
func TestIsWorkSetFlag(t *testing.T) {
	// Simulate a ramp scheme with warmup and work sets
	scheme := &mockScheme{
		schemeType: TypeRamp,
		sets: []GeneratedSet{
			{SetNumber: 1, Weight: 50, TargetReps: 5, IsWorkSet: false},  // warmup
			{SetNumber: 2, Weight: 65, TargetReps: 5, IsWorkSet: false},  // warmup
			{SetNumber: 3, Weight: 80, TargetReps: 5, IsWorkSet: false},  // warmup
			{SetNumber: 4, Weight: 90, TargetReps: 5, IsWorkSet: true},   // work
			{SetNumber: 5, Weight: 100, TargetReps: 5, IsWorkSet: true},  // work
		},
		valid: true,
	}

	ctx := DefaultSetGenerationContext()
	sets, err := scheme.GenerateSets(100, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	warmupCount := 0
	workSetCount := 0
	for _, set := range sets {
		if set.IsWorkSet {
			workSetCount++
		} else {
			warmupCount++
		}
	}

	if warmupCount != 3 {
		t.Errorf("expected 3 warmup sets, got %d", warmupCount)
	}
	if workSetCount != 2 {
		t.Errorf("expected 2 work sets, got %d", workSetCount)
	}
}

// TestSchemeEnvelope_MarshalJSON tests envelope marshaling.
func TestSchemeEnvelope_MarshalJSON(t *testing.T) {
	envelope := SchemeEnvelope{
		Type: TypeFixed,
	}

	data, err := json.Marshal(&envelope)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if parsed["type"] != string(TypeFixed) {
		t.Errorf("expected type %s, got %v", TypeFixed, parsed["type"])
	}
}

// TestSchemeEnvelope_Raw tests that Raw contains the original JSON data.
func TestSchemeEnvelope_Raw(t *testing.T) {
	jsonData := `{"type": "FIXED", "sets": 5, "reps": 5}`
	var envelope SchemeEnvelope
	if err := json.Unmarshal([]byte(jsonData), &envelope); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Raw should contain the full original JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal(envelope.Raw, &parsed); err != nil {
		t.Fatalf("failed to parse Raw: %v", err)
	}

	if parsed["type"] != "FIXED" {
		t.Errorf("expected type FIXED in Raw, got %v", parsed["type"])
	}
	if parsed["sets"] != float64(5) {
		t.Errorf("expected sets 5 in Raw, got %v", parsed["sets"])
	}
	if parsed["reps"] != float64(5) {
		t.Errorf("expected reps 5 in Raw, got %v", parsed["reps"])
	}
}

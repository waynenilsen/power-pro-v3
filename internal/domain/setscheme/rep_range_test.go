package setscheme

import (
	"encoding/json"
	"errors"
	"testing"
)

// TestRepRangeSetScheme_Type verifies the type discriminator.
func TestRepRangeSetScheme_Type(t *testing.T) {
	scheme := &RepRangeSetScheme{Sets: 3, MinReps: 8, MaxReps: 12}
	if scheme.Type() != TypeRepRange {
		t.Errorf("expected type %s, got %s", TypeRepRange, scheme.Type())
	}
}

// TestNewRepRangeSetScheme tests the constructor function.
func TestNewRepRangeSetScheme(t *testing.T) {
	tests := []struct {
		name    string
		sets    int
		minReps int
		maxReps int
		wantErr bool
	}{
		{"valid 3x8-12", 3, 8, 12, false},
		{"valid 4x6-8", 4, 6, 8, false},
		{"valid 1x5-10", 1, 5, 10, false},
		{"valid equal min/max", 3, 8, 8, false},
		{"valid 1x1-1 edge case", 1, 1, 1, false},
		{"invalid zero sets", 0, 8, 12, true},
		{"invalid negative sets", -1, 8, 12, true},
		{"invalid zero minReps", 3, 0, 12, true},
		{"invalid negative minReps", 3, -1, 12, true},
		{"invalid maxReps < minReps", 3, 12, 8, true},
		{"invalid zero maxReps", 3, 8, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme, err := NewRepRangeSetScheme(tt.sets, tt.minReps, tt.maxReps)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if scheme != nil {
					t.Error("expected nil scheme on error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if scheme == nil {
					t.Fatal("expected non-nil scheme")
				}
				if scheme.Sets != tt.sets {
					t.Errorf("expected Sets %d, got %d", tt.sets, scheme.Sets)
				}
				if scheme.MinReps != tt.minReps {
					t.Errorf("expected MinReps %d, got %d", tt.minReps, scheme.MinReps)
				}
				if scheme.MaxReps != tt.maxReps {
					t.Errorf("expected MaxReps %d, got %d", tt.maxReps, scheme.MaxReps)
				}
			}
		})
	}
}

// TestRepRangeSetScheme_Validate tests validation logic.
func TestRepRangeSetScheme_Validate(t *testing.T) {
	tests := []struct {
		name    string
		scheme  RepRangeSetScheme
		wantErr bool
	}{
		{"valid 3x8-12", RepRangeSetScheme{Sets: 3, MinReps: 8, MaxReps: 12}, false},
		{"valid 4x6-8", RepRangeSetScheme{Sets: 4, MinReps: 6, MaxReps: 8}, false},
		{"valid equal min/max", RepRangeSetScheme{Sets: 3, MinReps: 10, MaxReps: 10}, false},
		{"valid 1x1-1", RepRangeSetScheme{Sets: 1, MinReps: 1, MaxReps: 1}, false},
		{"valid large values", RepRangeSetScheme{Sets: 100, MinReps: 50, MaxReps: 100}, false},
		{"invalid zero sets", RepRangeSetScheme{Sets: 0, MinReps: 8, MaxReps: 12}, true},
		{"invalid negative sets", RepRangeSetScheme{Sets: -1, MinReps: 8, MaxReps: 12}, true},
		{"invalid zero minReps", RepRangeSetScheme{Sets: 3, MinReps: 0, MaxReps: 12}, true},
		{"invalid negative minReps", RepRangeSetScheme{Sets: 3, MinReps: -1, MaxReps: 12}, true},
		{"invalid maxReps < minReps", RepRangeSetScheme{Sets: 3, MinReps: 12, MaxReps: 8}, true},
		{"invalid zero maxReps with valid minReps", RepRangeSetScheme{Sets: 3, MinReps: 8, MaxReps: 0}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.scheme.Validate()
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if !errors.Is(err, ErrInvalidParams) {
					t.Errorf("expected ErrInvalidParams, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// TestRepRangeSetScheme_GenerateSets tests set generation.
func TestRepRangeSetScheme_GenerateSets(t *testing.T) {
	tests := []struct {
		name       string
		scheme     RepRangeSetScheme
		baseWeight float64
		wantSets   int
		wantReps   int // MinReps is used as target
	}{
		{"3x8-12 at 135", RepRangeSetScheme{Sets: 3, MinReps: 8, MaxReps: 12}, 135.0, 3, 8},
		{"4x6-8 at 185", RepRangeSetScheme{Sets: 4, MinReps: 6, MaxReps: 8}, 185.0, 4, 6},
		{"1x5-10 at 225", RepRangeSetScheme{Sets: 1, MinReps: 5, MaxReps: 10}, 225.0, 1, 5},
		{"3x10-15 at 95", RepRangeSetScheme{Sets: 3, MinReps: 10, MaxReps: 15}, 95.0, 3, 10},
		{"equal min/max", RepRangeSetScheme{Sets: 3, MinReps: 8, MaxReps: 8}, 135.0, 3, 8},
		{"zero weight", RepRangeSetScheme{Sets: 3, MinReps: 8, MaxReps: 12}, 0.0, 3, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := DefaultSetGenerationContext()
			sets, err := tt.scheme.GenerateSets(tt.baseWeight, ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(sets) != tt.wantSets {
				t.Errorf("expected %d sets, got %d", tt.wantSets, len(sets))
			}

			for i, set := range sets {
				// SetNumber should be 1-indexed
				expectedSetNumber := i + 1
				if set.SetNumber != expectedSetNumber {
					t.Errorf("set %d: expected SetNumber %d, got %d", i, expectedSetNumber, set.SetNumber)
				}

				// Weight should match baseWeight
				if set.Weight != tt.baseWeight {
					t.Errorf("set %d: expected Weight %f, got %f", i, tt.baseWeight, set.Weight)
				}

				// TargetReps should be MinReps
				if set.TargetReps != tt.wantReps {
					t.Errorf("set %d: expected TargetReps %d, got %d", i, tt.wantReps, set.TargetReps)
				}

				// All RepRange sets should be work sets
				if !set.IsWorkSet {
					t.Errorf("set %d: expected IsWorkSet to be true", i)
				}
			}
		})
	}
}

// TestRepRangeSetScheme_GenerateSets_InvalidScheme tests generation with invalid params.
func TestRepRangeSetScheme_GenerateSets_InvalidScheme(t *testing.T) {
	tests := []struct {
		name   string
		scheme RepRangeSetScheme
	}{
		{"zero sets", RepRangeSetScheme{Sets: 0, MinReps: 8, MaxReps: 12}},
		{"negative sets", RepRangeSetScheme{Sets: -1, MinReps: 8, MaxReps: 12}},
		{"zero minReps", RepRangeSetScheme{Sets: 3, MinReps: 0, MaxReps: 12}},
		{"negative minReps", RepRangeSetScheme{Sets: 3, MinReps: -1, MaxReps: 12}},
		{"maxReps < minReps", RepRangeSetScheme{Sets: 3, MinReps: 12, MaxReps: 8}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := DefaultSetGenerationContext()
			sets, err := tt.scheme.GenerateSets(135.0, ctx)
			if err == nil {
				t.Error("expected error, got nil")
			}
			if !errors.Is(err, ErrInvalidParams) {
				t.Errorf("expected ErrInvalidParams, got %v", err)
			}
			if sets != nil {
				t.Error("expected nil sets on error")
			}
		})
	}
}

// TestRepRangeSetScheme_GenerateSets_ContextIgnored tests that context is appropriately ignored.
func TestRepRangeSetScheme_GenerateSets_ContextIgnored(t *testing.T) {
	scheme := RepRangeSetScheme{Sets: 3, MinReps: 8, MaxReps: 12}

	// RepRange scheme ignores WorkSetThreshold - all sets are work sets
	contexts := []SetGenerationContext{
		{WorkSetThreshold: 0},
		{WorkSetThreshold: 50},
		{WorkSetThreshold: 80},
		{WorkSetThreshold: 100},
		DefaultSetGenerationContext(),
	}

	for i, ctx := range contexts {
		sets, err := scheme.GenerateSets(135.0, ctx)
		if err != nil {
			t.Fatalf("context %d: unexpected error: %v", i, err)
		}

		for j, set := range sets {
			if !set.IsWorkSet {
				t.Errorf("context %d, set %d: expected IsWorkSet true regardless of context", i, j)
			}
		}
	}
}

// TestRepRangeSetScheme_MarshalJSON tests JSON serialization.
func TestRepRangeSetScheme_MarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		scheme RepRangeSetScheme
	}{
		{"3x8-12", RepRangeSetScheme{Sets: 3, MinReps: 8, MaxReps: 12}},
		{"4x6-8", RepRangeSetScheme{Sets: 4, MinReps: 6, MaxReps: 8}},
		{"1x1-1", RepRangeSetScheme{Sets: 1, MinReps: 1, MaxReps: 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(&tt.scheme)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Parse back and verify structure
			var parsed map[string]interface{}
			if err := json.Unmarshal(data, &parsed); err != nil {
				t.Fatalf("failed to parse JSON: %v", err)
			}

			// Verify type discriminator
			if parsed["type"] != string(TypeRepRange) {
				t.Errorf("expected type %s, got %v", TypeRepRange, parsed["type"])
			}

			// Verify sets
			if int(parsed["sets"].(float64)) != tt.scheme.Sets {
				t.Errorf("expected sets %d, got %v", tt.scheme.Sets, parsed["sets"])
			}

			// Verify minReps
			if int(parsed["minReps"].(float64)) != tt.scheme.MinReps {
				t.Errorf("expected minReps %d, got %v", tt.scheme.MinReps, parsed["minReps"])
			}

			// Verify maxReps
			if int(parsed["maxReps"].(float64)) != tt.scheme.MaxReps {
				t.Errorf("expected maxReps %d, got %v", tt.scheme.MaxReps, parsed["maxReps"])
			}
		})
	}
}

// TestUnmarshalRepRangeSetScheme tests JSON deserialization.
func TestUnmarshalRepRangeSetScheme(t *testing.T) {
	tests := []struct {
		name        string
		json        string
		wantSets    int
		wantMinReps int
		wantMaxReps int
		wantErr     bool
	}{
		{
			name:        "valid 3x8-12",
			json:        `{"type": "REP_RANGE", "sets": 3, "minReps": 8, "maxReps": 12}`,
			wantSets:    3,
			wantMinReps: 8,
			wantMaxReps: 12,
			wantErr:     false,
		},
		{
			name:        "valid 4x6-8",
			json:        `{"type": "REP_RANGE", "sets": 4, "minReps": 6, "maxReps": 8}`,
			wantSets:    4,
			wantMinReps: 6,
			wantMaxReps: 8,
			wantErr:     false,
		},
		{
			name:        "valid equal min/max",
			json:        `{"type": "REP_RANGE", "sets": 3, "minReps": 10, "maxReps": 10}`,
			wantSets:    3,
			wantMinReps: 10,
			wantMaxReps: 10,
			wantErr:     false,
		},
		{
			name:        "without type (still valid)",
			json:        `{"sets": 3, "minReps": 8, "maxReps": 12}`,
			wantSets:    3,
			wantMinReps: 8,
			wantMaxReps: 12,
			wantErr:     false,
		},
		{
			name:    "invalid zero sets",
			json:    `{"type": "REP_RANGE", "sets": 0, "minReps": 8, "maxReps": 12}`,
			wantErr: true,
		},
		{
			name:    "invalid negative sets",
			json:    `{"type": "REP_RANGE", "sets": -1, "minReps": 8, "maxReps": 12}`,
			wantErr: true,
		},
		{
			name:    "invalid zero minReps",
			json:    `{"type": "REP_RANGE", "sets": 3, "minReps": 0, "maxReps": 12}`,
			wantErr: true,
		},
		{
			name:    "invalid negative minReps",
			json:    `{"type": "REP_RANGE", "sets": 3, "minReps": -1, "maxReps": 12}`,
			wantErr: true,
		},
		{
			name:    "invalid maxReps < minReps",
			json:    `{"type": "REP_RANGE", "sets": 3, "minReps": 12, "maxReps": 8}`,
			wantErr: true,
		},
		{
			name:    "missing sets field",
			json:    `{"type": "REP_RANGE", "minReps": 8, "maxReps": 12}`,
			wantErr: true,
		},
		{
			name:    "missing minReps field",
			json:    `{"type": "REP_RANGE", "sets": 3, "maxReps": 12}`,
			wantErr: true,
		},
		{
			name:    "missing maxReps field",
			json:    `{"type": "REP_RANGE", "sets": 3, "minReps": 8}`,
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			json:    `{invalid}`,
			wantErr: true,
		},
		{
			name:    "empty object",
			json:    `{}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme, err := UnmarshalRepRangeSetScheme(json.RawMessage(tt.json))
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if scheme == nil {
					t.Fatal("expected non-nil scheme")
				}

				repRange, ok := scheme.(*RepRangeSetScheme)
				if !ok {
					t.Fatal("expected *RepRangeSetScheme")
				}
				if repRange.Sets != tt.wantSets {
					t.Errorf("expected Sets %d, got %d", tt.wantSets, repRange.Sets)
				}
				if repRange.MinReps != tt.wantMinReps {
					t.Errorf("expected MinReps %d, got %d", tt.wantMinReps, repRange.MinReps)
				}
				if repRange.MaxReps != tt.wantMaxReps {
					t.Errorf("expected MaxReps %d, got %d", tt.wantMaxReps, repRange.MaxReps)
				}
			}
		})
	}
}

// TestRepRangeSetScheme_RoundTrip tests JSON round-trip serialization.
func TestRepRangeSetScheme_RoundTrip(t *testing.T) {
	original := RepRangeSetScheme{Sets: 3, MinReps: 8, MaxReps: 12}

	// Marshal
	data, err := json.Marshal(&original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	// Unmarshal
	scheme, err := UnmarshalRepRangeSetScheme(data)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	// Verify
	repRange, ok := scheme.(*RepRangeSetScheme)
	if !ok {
		t.Fatal("expected *RepRangeSetScheme")
	}
	if repRange.Sets != original.Sets {
		t.Errorf("expected Sets %d, got %d", original.Sets, repRange.Sets)
	}
	if repRange.MinReps != original.MinReps {
		t.Errorf("expected MinReps %d, got %d", original.MinReps, repRange.MinReps)
	}
	if repRange.MaxReps != original.MaxReps {
		t.Errorf("expected MaxReps %d, got %d", original.MaxReps, repRange.MaxReps)
	}
}

// TestRegisterRepRangeScheme tests factory registration.
func TestRegisterRepRangeScheme(t *testing.T) {
	factory := NewSchemeFactory()

	// Should not be registered initially
	if factory.IsRegistered(TypeRepRange) {
		t.Error("TypeRepRange should not be registered initially")
	}

	// Register
	RegisterRepRangeScheme(factory)

	// Should be registered now
	if !factory.IsRegistered(TypeRepRange) {
		t.Error("TypeRepRange should be registered after RegisterRepRangeScheme")
	}
}

// TestRepRangeSetScheme_FactoryIntegration tests full factory workflow.
func TestRepRangeSetScheme_FactoryIntegration(t *testing.T) {
	factory := NewSchemeFactory()
	RegisterRepRangeScheme(factory)

	t.Run("Create from type and data", func(t *testing.T) {
		jsonData := json.RawMessage(`{"sets": 3, "minReps": 8, "maxReps": 12}`)
		scheme, err := factory.Create(TypeRepRange, jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if scheme.Type() != TypeRepRange {
			t.Errorf("expected type %s, got %s", TypeRepRange, scheme.Type())
		}
	})

	t.Run("CreateFromJSON", func(t *testing.T) {
		jsonData := json.RawMessage(`{"type": "REP_RANGE", "sets": 4, "minReps": 6, "maxReps": 8}`)
		scheme, err := factory.CreateFromJSON(jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if scheme.Type() != TypeRepRange {
			t.Errorf("expected type %s, got %s", TypeRepRange, scheme.Type())
		}

		// Generate sets
		ctx := DefaultSetGenerationContext()
		sets, err := scheme.GenerateSets(185.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sets) != 4 {
			t.Errorf("expected 4 sets, got %d", len(sets))
		}
		for _, set := range sets {
			if set.TargetReps != 6 {
				t.Errorf("expected 6 reps (minReps), got %d", set.TargetReps)
			}
		}
	})

	t.Run("Invalid JSON in CreateFromJSON", func(t *testing.T) {
		jsonData := json.RawMessage(`{"type": "REP_RANGE", "sets": 0, "minReps": 8, "maxReps": 12}`)
		_, err := factory.CreateFromJSON(jsonData)
		if err == nil {
			t.Error("expected error for invalid sets")
		}
	})
}

// TestRepRangeSetScheme_Implements_SetScheme verifies interface implementation.
func TestRepRangeSetScheme_Implements_SetScheme(t *testing.T) {
	var _ SetScheme = (*RepRangeSetScheme)(nil)
}

// TestRepRangeSetScheme_EdgeCases tests edge cases.
func TestRepRangeSetScheme_EdgeCases(t *testing.T) {
	t.Run("single set narrow range", func(t *testing.T) {
		scheme := RepRangeSetScheme{Sets: 1, MinReps: 5, MaxReps: 5}
		ctx := DefaultSetGenerationContext()
		sets, err := scheme.GenerateSets(225.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sets) != 1 {
			t.Errorf("expected 1 set, got %d", len(sets))
		}
		if sets[0].SetNumber != 1 {
			t.Errorf("expected SetNumber 1, got %d", sets[0].SetNumber)
		}
		if sets[0].TargetReps != 5 {
			t.Errorf("expected TargetReps 5, got %d", sets[0].TargetReps)
		}
		if sets[0].Weight != 225.0 {
			t.Errorf("expected Weight 225.0, got %f", sets[0].Weight)
		}
		if !sets[0].IsWorkSet {
			t.Error("expected IsWorkSet true")
		}
	})

	t.Run("many sets wide range", func(t *testing.T) {
		scheme := RepRangeSetScheme{Sets: 10, MinReps: 8, MaxReps: 20}
		ctx := DefaultSetGenerationContext()
		sets, err := scheme.GenerateSets(100.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sets) != 10 {
			t.Errorf("expected 10 sets, got %d", len(sets))
		}
		for i, set := range sets {
			if set.SetNumber != i+1 {
				t.Errorf("set %d: expected SetNumber %d, got %d", i, i+1, set.SetNumber)
			}
			// TargetReps should be MinReps
			if set.TargetReps != 8 {
				t.Errorf("set %d: expected TargetReps 8, got %d", i, set.TargetReps)
			}
		}
	})

	t.Run("fractional weight preserved", func(t *testing.T) {
		scheme := RepRangeSetScheme{Sets: 3, MinReps: 8, MaxReps: 12}
		ctx := DefaultSetGenerationContext()
		sets, err := scheme.GenerateSets(142.5, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for i, set := range sets {
			if set.Weight != 142.5 {
				t.Errorf("set %d: expected Weight 142.5, got %f", i, set.Weight)
			}
		}
	})

	t.Run("very large weight", func(t *testing.T) {
		scheme := RepRangeSetScheme{Sets: 1, MinReps: 8, MaxReps: 12}
		ctx := DefaultSetGenerationContext()
		sets, err := scheme.GenerateSets(1000.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sets[0].Weight != 1000.0 {
			t.Errorf("expected Weight 1000.0, got %f", sets[0].Weight)
		}
	})
}

// TestRepRangeSetScheme_ProgramExamples tests common program configurations.
func TestRepRangeSetScheme_ProgramExamples(t *testing.T) {
	tests := []struct {
		name       string
		program    string
		sets       int
		minReps    int
		maxReps    int
		baseWeight float64
	}{
		{"Hypertrophy Accessory", "3x8-12", 3, 8, 12, 135.0},
		{"Double Progression", "3x6-10", 3, 6, 10, 155.0},
		{"Pump Work", "4x10-15", 4, 10, 15, 95.0},
		{"GZCLP Tier 3", "3x15-25", 3, 15, 25, 50.0},
		{"Strength Range", "3x3-5", 3, 3, 5, 275.0},
	}

	for _, tt := range tests {
		t.Run(tt.name+" "+tt.program, func(t *testing.T) {
			scheme, err := NewRepRangeSetScheme(tt.sets, tt.minReps, tt.maxReps)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			ctx := DefaultSetGenerationContext()
			sets, err := scheme.GenerateSets(tt.baseWeight, ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Verify correct number of sets
			if len(sets) != tt.sets {
				t.Errorf("expected %d sets, got %d", tt.sets, len(sets))
			}

			// Verify each set
			for i, set := range sets {
				if set.SetNumber != i+1 {
					t.Errorf("set %d: expected SetNumber %d, got %d", i, i+1, set.SetNumber)
				}
				if set.Weight != tt.baseWeight {
					t.Errorf("set %d: expected Weight %f, got %f", i, tt.baseWeight, set.Weight)
				}
				// TargetReps should be MinReps
				if set.TargetReps != tt.minReps {
					t.Errorf("set %d: expected TargetReps %d, got %d", i, tt.minReps, set.TargetReps)
				}
				if !set.IsWorkSet {
					t.Errorf("set %d: expected IsWorkSet true", i)
				}
			}
		})
	}
}

// TestRepRangeSetScheme_ValidationErrorMessages tests error message clarity.
func TestRepRangeSetScheme_ValidationErrorMessages(t *testing.T) {
	t.Run("zero sets error message", func(t *testing.T) {
		scheme := RepRangeSetScheme{Sets: 0, MinReps: 8, MaxReps: 12}
		err := scheme.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, ErrInvalidParams) {
			t.Errorf("expected ErrInvalidParams, got %v", err)
		}
		if err.Error() == "" {
			t.Error("expected non-empty error message")
		}
	})

	t.Run("zero minReps error message", func(t *testing.T) {
		scheme := RepRangeSetScheme{Sets: 3, MinReps: 0, MaxReps: 12}
		err := scheme.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, ErrInvalidParams) {
			t.Errorf("expected ErrInvalidParams, got %v", err)
		}
	})

	t.Run("maxReps < minReps error message", func(t *testing.T) {
		scheme := RepRangeSetScheme{Sets: 3, MinReps: 12, MaxReps: 8}
		err := scheme.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, ErrInvalidParams) {
			t.Errorf("expected ErrInvalidParams, got %v", err)
		}
	})

	t.Run("sets error takes precedence", func(t *testing.T) {
		// When sets is invalid, that error comes first
		scheme := RepRangeSetScheme{Sets: 0, MinReps: 0, MaxReps: 0}
		err := scheme.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		// Should mention sets since that's checked first
	})
}

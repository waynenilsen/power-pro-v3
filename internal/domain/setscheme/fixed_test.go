package setscheme

import (
	"encoding/json"
	"errors"
	"testing"
)

// TestFixedSetScheme_Type verifies the type discriminator.
func TestFixedSetScheme_Type(t *testing.T) {
	scheme := &FixedSetScheme{Sets: 5, Reps: 5}
	if scheme.Type() != TypeFixed {
		t.Errorf("expected type %s, got %s", TypeFixed, scheme.Type())
	}
}

// TestNewFixedSetScheme tests the constructor function.
func TestNewFixedSetScheme(t *testing.T) {
	tests := []struct {
		name    string
		sets    int
		reps    int
		wantErr bool
	}{
		{"valid 5x5", 5, 5, false},
		{"valid 3x8", 3, 8, false},
		{"valid 1x5", 1, 5, false},
		{"valid 1x1 edge case", 1, 1, false},
		{"valid 10x10", 10, 10, false},
		{"invalid zero sets", 0, 5, true},
		{"invalid negative sets", -1, 5, true},
		{"invalid zero reps", 5, 0, true},
		{"invalid negative reps", 5, -1, true},
		{"invalid both zero", 0, 0, true},
		{"invalid both negative", -1, -2, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme, err := NewFixedSetScheme(tt.sets, tt.reps)
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
				if scheme.Reps != tt.reps {
					t.Errorf("expected Reps %d, got %d", tt.reps, scheme.Reps)
				}
			}
		})
	}
}

// TestFixedSetScheme_Validate tests validation logic.
func TestFixedSetScheme_Validate(t *testing.T) {
	tests := []struct {
		name    string
		scheme  FixedSetScheme
		wantErr bool
	}{
		{"valid 5x5", FixedSetScheme{Sets: 5, Reps: 5}, false},
		{"valid 3x8", FixedSetScheme{Sets: 3, Reps: 8}, false},
		{"valid 1x5", FixedSetScheme{Sets: 1, Reps: 5}, false},
		{"valid 1x1", FixedSetScheme{Sets: 1, Reps: 1}, false},
		{"valid large values", FixedSetScheme{Sets: 100, Reps: 100}, false},
		{"invalid zero sets", FixedSetScheme{Sets: 0, Reps: 5}, true},
		{"invalid negative sets", FixedSetScheme{Sets: -1, Reps: 5}, true},
		{"invalid zero reps", FixedSetScheme{Sets: 5, Reps: 0}, true},
		{"invalid negative reps", FixedSetScheme{Sets: 5, Reps: -1}, true},
		{"invalid both zero", FixedSetScheme{Sets: 0, Reps: 0}, true},
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

// TestFixedSetScheme_GenerateSets tests set generation.
func TestFixedSetScheme_GenerateSets(t *testing.T) {
	tests := []struct {
		name       string
		scheme     FixedSetScheme
		baseWeight float64
		wantSets   int
		wantReps   int
	}{
		{"5x5 at 265", FixedSetScheme{Sets: 5, Reps: 5}, 265.0, 5, 5},
		{"3x8 at 185", FixedSetScheme{Sets: 3, Reps: 8}, 185.0, 3, 8},
		{"1x5 at 315", FixedSetScheme{Sets: 1, Reps: 5}, 315.0, 1, 5},
		{"5x10 BBB at 135", FixedSetScheme{Sets: 5, Reps: 10}, 135.0, 5, 10},
		{"4x6 at 225", FixedSetScheme{Sets: 4, Reps: 6}, 225.0, 4, 6},
		{"1x1 single", FixedSetScheme{Sets: 1, Reps: 1}, 405.0, 1, 1},
		{"zero weight", FixedSetScheme{Sets: 3, Reps: 5}, 0.0, 3, 5},
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

				// TargetReps should match scheme reps
				if set.TargetReps != tt.wantReps {
					t.Errorf("set %d: expected TargetReps %d, got %d", i, tt.wantReps, set.TargetReps)
				}

				// All Fixed sets should be work sets
				if !set.IsWorkSet {
					t.Errorf("set %d: expected IsWorkSet to be true", i)
				}
			}
		})
	}
}

// TestFixedSetScheme_GenerateSets_InvalidScheme tests generation with invalid params.
func TestFixedSetScheme_GenerateSets_InvalidScheme(t *testing.T) {
	tests := []struct {
		name   string
		scheme FixedSetScheme
	}{
		{"zero sets", FixedSetScheme{Sets: 0, Reps: 5}},
		{"negative sets", FixedSetScheme{Sets: -1, Reps: 5}},
		{"zero reps", FixedSetScheme{Sets: 5, Reps: 0}},
		{"negative reps", FixedSetScheme{Sets: 5, Reps: -1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := DefaultSetGenerationContext()
			sets, err := tt.scheme.GenerateSets(265.0, ctx)
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

// TestFixedSetScheme_GenerateSets_ContextIgnored tests that context is appropriately ignored.
func TestFixedSetScheme_GenerateSets_ContextIgnored(t *testing.T) {
	scheme := FixedSetScheme{Sets: 5, Reps: 5}

	// Fixed scheme ignores WorkSetThreshold - all sets are work sets
	contexts := []SetGenerationContext{
		{WorkSetThreshold: 0},
		{WorkSetThreshold: 50},
		{WorkSetThreshold: 80},
		{WorkSetThreshold: 100},
		DefaultSetGenerationContext(),
	}

	for i, ctx := range contexts {
		sets, err := scheme.GenerateSets(265.0, ctx)
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

// TestFixedSetScheme_MarshalJSON tests JSON serialization.
func TestFixedSetScheme_MarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		scheme FixedSetScheme
	}{
		{"5x5", FixedSetScheme{Sets: 5, Reps: 5}},
		{"3x8", FixedSetScheme{Sets: 3, Reps: 8}},
		{"1x1", FixedSetScheme{Sets: 1, Reps: 1}},
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
			if parsed["type"] != string(TypeFixed) {
				t.Errorf("expected type %s, got %v", TypeFixed, parsed["type"])
			}

			// Verify sets
			if int(parsed["sets"].(float64)) != tt.scheme.Sets {
				t.Errorf("expected sets %d, got %v", tt.scheme.Sets, parsed["sets"])
			}

			// Verify reps
			if int(parsed["reps"].(float64)) != tt.scheme.Reps {
				t.Errorf("expected reps %d, got %v", tt.scheme.Reps, parsed["reps"])
			}
		})
	}
}

// TestUnmarshalFixedSetScheme tests JSON deserialization.
func TestUnmarshalFixedSetScheme(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		wantSets int
		wantReps int
		wantErr  bool
	}{
		{
			name:     "valid 5x5",
			json:     `{"type": "FIXED", "sets": 5, "reps": 5}`,
			wantSets: 5,
			wantReps: 5,
			wantErr:  false,
		},
		{
			name:     "valid 3x8",
			json:     `{"type": "FIXED", "sets": 3, "reps": 8}`,
			wantSets: 3,
			wantReps: 8,
			wantErr:  false,
		},
		{
			name:     "valid 1x1",
			json:     `{"type": "FIXED", "sets": 1, "reps": 1}`,
			wantSets: 1,
			wantReps: 1,
			wantErr:  false,
		},
		{
			name:     "without type (still valid)",
			json:     `{"sets": 5, "reps": 5}`,
			wantSets: 5,
			wantReps: 5,
			wantErr:  false,
		},
		{
			name:    "invalid zero sets",
			json:    `{"type": "FIXED", "sets": 0, "reps": 5}`,
			wantErr: true,
		},
		{
			name:    "invalid negative sets",
			json:    `{"type": "FIXED", "sets": -1, "reps": 5}`,
			wantErr: true,
		},
		{
			name:    "invalid zero reps",
			json:    `{"type": "FIXED", "sets": 5, "reps": 0}`,
			wantErr: true,
		},
		{
			name:    "invalid negative reps",
			json:    `{"type": "FIXED", "sets": 5, "reps": -1}`,
			wantErr: true,
		},
		{
			name:    "missing sets field",
			json:    `{"type": "FIXED", "reps": 5}`,
			wantErr: true,
		},
		{
			name:    "missing reps field",
			json:    `{"type": "FIXED", "sets": 5}`,
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
			scheme, err := UnmarshalFixedSetScheme(json.RawMessage(tt.json))
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

				fixed, ok := scheme.(*FixedSetScheme)
				if !ok {
					t.Fatal("expected *FixedSetScheme")
				}
				if fixed.Sets != tt.wantSets {
					t.Errorf("expected Sets %d, got %d", tt.wantSets, fixed.Sets)
				}
				if fixed.Reps != tt.wantReps {
					t.Errorf("expected Reps %d, got %d", tt.wantReps, fixed.Reps)
				}
			}
		})
	}
}

// TestFixedSetScheme_RoundTrip tests JSON round-trip serialization.
func TestFixedSetScheme_RoundTrip(t *testing.T) {
	original := FixedSetScheme{Sets: 5, Reps: 5}

	// Marshal
	data, err := json.Marshal(&original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	// Unmarshal
	scheme, err := UnmarshalFixedSetScheme(data)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	// Verify
	fixed, ok := scheme.(*FixedSetScheme)
	if !ok {
		t.Fatal("expected *FixedSetScheme")
	}
	if fixed.Sets != original.Sets {
		t.Errorf("expected Sets %d, got %d", original.Sets, fixed.Sets)
	}
	if fixed.Reps != original.Reps {
		t.Errorf("expected Reps %d, got %d", original.Reps, fixed.Reps)
	}
}

// TestRegisterFixedScheme tests factory registration.
func TestRegisterFixedScheme(t *testing.T) {
	factory := NewSchemeFactory()

	// Should not be registered initially
	if factory.IsRegistered(TypeFixed) {
		t.Error("TypeFixed should not be registered initially")
	}

	// Register
	RegisterFixedScheme(factory)

	// Should be registered now
	if !factory.IsRegistered(TypeFixed) {
		t.Error("TypeFixed should be registered after RegisterFixedScheme")
	}
}

// TestFixedSetScheme_FactoryIntegration tests full factory workflow.
func TestFixedSetScheme_FactoryIntegration(t *testing.T) {
	factory := NewSchemeFactory()
	RegisterFixedScheme(factory)

	t.Run("Create from type and data", func(t *testing.T) {
		jsonData := json.RawMessage(`{"sets": 5, "reps": 5}`)
		scheme, err := factory.Create(TypeFixed, jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if scheme.Type() != TypeFixed {
			t.Errorf("expected type %s, got %s", TypeFixed, scheme.Type())
		}
	})

	t.Run("CreateFromJSON", func(t *testing.T) {
		jsonData := json.RawMessage(`{"type": "FIXED", "sets": 3, "reps": 8}`)
		scheme, err := factory.CreateFromJSON(jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if scheme.Type() != TypeFixed {
			t.Errorf("expected type %s, got %s", TypeFixed, scheme.Type())
		}

		// Generate sets
		ctx := DefaultSetGenerationContext()
		sets, err := scheme.GenerateSets(185.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sets) != 3 {
			t.Errorf("expected 3 sets, got %d", len(sets))
		}
		for _, set := range sets {
			if set.TargetReps != 8 {
				t.Errorf("expected 8 reps, got %d", set.TargetReps)
			}
		}
	})

	t.Run("Invalid JSON in CreateFromJSON", func(t *testing.T) {
		jsonData := json.RawMessage(`{"type": "FIXED", "sets": 0, "reps": 5}`)
		_, err := factory.CreateFromJSON(jsonData)
		if err == nil {
			t.Error("expected error for invalid sets")
		}
	})
}

// TestFixedSetScheme_Implements_SetScheme verifies interface implementation.
func TestFixedSetScheme_Implements_SetScheme(t *testing.T) {
	var _ SetScheme = (*FixedSetScheme)(nil)
}

// TestFixedSetScheme_EdgeCases tests edge cases.
func TestFixedSetScheme_EdgeCases(t *testing.T) {
	t.Run("single set single rep", func(t *testing.T) {
		scheme := FixedSetScheme{Sets: 1, Reps: 1}
		ctx := DefaultSetGenerationContext()
		sets, err := scheme.GenerateSets(405.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sets) != 1 {
			t.Errorf("expected 1 set, got %d", len(sets))
		}
		if sets[0].SetNumber != 1 {
			t.Errorf("expected SetNumber 1, got %d", sets[0].SetNumber)
		}
		if sets[0].TargetReps != 1 {
			t.Errorf("expected TargetReps 1, got %d", sets[0].TargetReps)
		}
		if sets[0].Weight != 405.0 {
			t.Errorf("expected Weight 405.0, got %f", sets[0].Weight)
		}
		if !sets[0].IsWorkSet {
			t.Error("expected IsWorkSet true")
		}
	})

	t.Run("many sets", func(t *testing.T) {
		scheme := FixedSetScheme{Sets: 20, Reps: 5}
		ctx := DefaultSetGenerationContext()
		sets, err := scheme.GenerateSets(135.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sets) != 20 {
			t.Errorf("expected 20 sets, got %d", len(sets))
		}
		for i, set := range sets {
			if set.SetNumber != i+1 {
				t.Errorf("set %d: expected SetNumber %d, got %d", i, i+1, set.SetNumber)
			}
		}
	})

	t.Run("fractional weight preserved", func(t *testing.T) {
		scheme := FixedSetScheme{Sets: 3, Reps: 5}
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
		scheme := FixedSetScheme{Sets: 1, Reps: 1}
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

// TestFixedSetScheme_ProgramExamples tests common program configurations.
func TestFixedSetScheme_ProgramExamples(t *testing.T) {
	tests := []struct {
		name       string
		program    string
		sets       int
		reps       int
		baseWeight float64
	}{
		{"Starting Strength", "3x5", 3, 5, 265.0},
		{"Bill Starr", "5x5", 5, 5, 225.0},
		{"BBB", "5x10", 5, 10, 135.0},
		{"Generic Hypertrophy", "4x8", 4, 8, 185.0},
		{"Volume Work", "6x6", 6, 6, 175.0},
	}

	for _, tt := range tests {
		t.Run(tt.name+" "+tt.program, func(t *testing.T) {
			scheme, err := NewFixedSetScheme(tt.sets, tt.reps)
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
				if set.TargetReps != tt.reps {
					t.Errorf("set %d: expected TargetReps %d, got %d", i, tt.reps, set.TargetReps)
				}
				if !set.IsWorkSet {
					t.Errorf("set %d: expected IsWorkSet true", i)
				}
			}
		})
	}
}

// TestFixedSetScheme_ValidationErrorMessages tests error message clarity.
func TestFixedSetScheme_ValidationErrorMessages(t *testing.T) {
	t.Run("zero sets error message", func(t *testing.T) {
		scheme := FixedSetScheme{Sets: 0, Reps: 5}
		err := scheme.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, ErrInvalidParams) {
			t.Errorf("expected ErrInvalidParams, got %v", err)
		}
		// Error message should mention "sets"
		if err.Error() == "" {
			t.Error("expected non-empty error message")
		}
	})

	t.Run("zero reps error message", func(t *testing.T) {
		scheme := FixedSetScheme{Sets: 5, Reps: 0}
		err := scheme.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, ErrInvalidParams) {
			t.Errorf("expected ErrInvalidParams, got %v", err)
		}
	})

	t.Run("sets error takes precedence", func(t *testing.T) {
		// When both are invalid, sets error comes first
		scheme := FixedSetScheme{Sets: 0, Reps: 0}
		err := scheme.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		// Should mention sets since that's checked first
	})
}

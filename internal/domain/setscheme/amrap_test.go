package setscheme

import (
	"encoding/json"
	"errors"
	"testing"
)

// TestAMRAPSetScheme_Type verifies the type discriminator.
func TestAMRAPSetScheme_Type(t *testing.T) {
	scheme := &AMRAPSetScheme{Sets: 1, MinReps: 5}
	if scheme.Type() != TypeAMRAP {
		t.Errorf("expected type %s, got %s", TypeAMRAP, scheme.Type())
	}
}

// TestNewAMRAPSetScheme tests the constructor function.
func TestNewAMRAPSetScheme(t *testing.T) {
	tests := []struct {
		name    string
		sets    int
		minReps int
		wantErr bool
	}{
		{"valid 1x5+", 1, 5, false},
		{"valid 1x3+", 1, 3, false},
		{"valid 1x1+", 1, 1, false},
		{"valid 2x5+ multiple AMRAP sets", 2, 5, false},
		{"valid 3x8+", 3, 8, false},
		{"invalid zero sets", 0, 5, true},
		{"invalid negative sets", -1, 5, true},
		{"invalid zero minReps", 1, 0, true},
		{"invalid negative minReps", 1, -1, true},
		{"invalid both zero", 0, 0, true},
		{"invalid both negative", -1, -2, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme, err := NewAMRAPSetScheme(tt.sets, tt.minReps)
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
			}
		})
	}
}

// TestAMRAPSetScheme_Validate tests validation logic.
func TestAMRAPSetScheme_Validate(t *testing.T) {
	tests := []struct {
		name    string
		scheme  AMRAPSetScheme
		wantErr bool
	}{
		{"valid 1x5+", AMRAPSetScheme{Sets: 1, MinReps: 5}, false},
		{"valid 1x3+", AMRAPSetScheme{Sets: 1, MinReps: 3}, false},
		{"valid 1x1+", AMRAPSetScheme{Sets: 1, MinReps: 1}, false},
		{"valid large values", AMRAPSetScheme{Sets: 10, MinReps: 20}, false},
		{"invalid zero sets", AMRAPSetScheme{Sets: 0, MinReps: 5}, true},
		{"invalid negative sets", AMRAPSetScheme{Sets: -1, MinReps: 5}, true},
		{"invalid zero minReps", AMRAPSetScheme{Sets: 1, MinReps: 0}, true},
		{"invalid negative minReps", AMRAPSetScheme{Sets: 1, MinReps: -1}, true},
		{"invalid both zero", AMRAPSetScheme{Sets: 0, MinReps: 0}, true},
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

// TestAMRAPSetScheme_GenerateSets tests set generation.
func TestAMRAPSetScheme_GenerateSets(t *testing.T) {
	tests := []struct {
		name        string
		scheme      AMRAPSetScheme
		baseWeight  float64
		wantSets    int
		wantMinReps int
	}{
		{"1x5+ at 315 (Wendler)", AMRAPSetScheme{Sets: 1, MinReps: 5}, 315.0, 1, 5},
		{"1x3+ at 335", AMRAPSetScheme{Sets: 1, MinReps: 3}, 335.0, 1, 3},
		{"1x1+ at 365", AMRAPSetScheme{Sets: 1, MinReps: 1}, 365.0, 1, 1},
		{"2x5+ multiple AMRAP", AMRAPSetScheme{Sets: 2, MinReps: 5}, 225.0, 2, 5},
		{"1x8+ at 185", AMRAPSetScheme{Sets: 1, MinReps: 8}, 185.0, 1, 8},
		{"zero weight", AMRAPSetScheme{Sets: 1, MinReps: 5}, 0.0, 1, 5},
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

				// TargetReps should match MinReps
				if set.TargetReps != tt.wantMinReps {
					t.Errorf("set %d: expected TargetReps %d, got %d", i, tt.wantMinReps, set.TargetReps)
				}

				// All AMRAP sets should be work sets
				if !set.IsWorkSet {
					t.Errorf("set %d: expected IsWorkSet to be true", i)
				}
			}
		})
	}
}

// TestAMRAPSetScheme_GenerateSets_InvalidScheme tests generation with invalid params.
func TestAMRAPSetScheme_GenerateSets_InvalidScheme(t *testing.T) {
	tests := []struct {
		name   string
		scheme AMRAPSetScheme
	}{
		{"zero sets", AMRAPSetScheme{Sets: 0, MinReps: 5}},
		{"negative sets", AMRAPSetScheme{Sets: -1, MinReps: 5}},
		{"zero minReps", AMRAPSetScheme{Sets: 1, MinReps: 0}},
		{"negative minReps", AMRAPSetScheme{Sets: 1, MinReps: -1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := DefaultSetGenerationContext()
			sets, err := tt.scheme.GenerateSets(315.0, ctx)
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

// TestAMRAPSetScheme_GenerateSets_ContextIgnored tests that context is appropriately ignored.
func TestAMRAPSetScheme_GenerateSets_ContextIgnored(t *testing.T) {
	scheme := AMRAPSetScheme{Sets: 1, MinReps: 5}

	// AMRAP scheme ignores WorkSetThreshold - all sets are work sets
	contexts := []SetGenerationContext{
		{WorkSetThreshold: 0},
		{WorkSetThreshold: 50},
		{WorkSetThreshold: 80},
		{WorkSetThreshold: 100},
		DefaultSetGenerationContext(),
	}

	for i, ctx := range contexts {
		sets, err := scheme.GenerateSets(315.0, ctx)
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

// TestAMRAPSetScheme_MarshalJSON tests JSON serialization.
func TestAMRAPSetScheme_MarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		scheme AMRAPSetScheme
	}{
		{"1x5+", AMRAPSetScheme{Sets: 1, MinReps: 5}},
		{"1x3+", AMRAPSetScheme{Sets: 1, MinReps: 3}},
		{"1x1+", AMRAPSetScheme{Sets: 1, MinReps: 1}},
		{"2x8+", AMRAPSetScheme{Sets: 2, MinReps: 8}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(&tt.scheme)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Parse back and verify structure
			var parsed map[string]any
			if err := json.Unmarshal(data, &parsed); err != nil {
				t.Fatalf("failed to parse JSON: %v", err)
			}

			// Verify type discriminator
			if parsed["type"] != string(TypeAMRAP) {
				t.Errorf("expected type %s, got %v", TypeAMRAP, parsed["type"])
			}

			// Verify sets
			if int(parsed["sets"].(float64)) != tt.scheme.Sets {
				t.Errorf("expected sets %d, got %v", tt.scheme.Sets, parsed["sets"])
			}

			// Verify minReps
			if int(parsed["minReps"].(float64)) != tt.scheme.MinReps {
				t.Errorf("expected minReps %d, got %v", tt.scheme.MinReps, parsed["minReps"])
			}
		})
	}
}

// TestUnmarshalAMRAPSetScheme tests JSON deserialization.
func TestUnmarshalAMRAPSetScheme(t *testing.T) {
	tests := []struct {
		name        string
		json        string
		wantSets    int
		wantMinReps int
		wantErr     bool
	}{
		{
			name:        "valid 1x5+",
			json:        `{"type": "AMRAP", "sets": 1, "minReps": 5}`,
			wantSets:    1,
			wantMinReps: 5,
			wantErr:     false,
		},
		{
			name:        "valid 1x3+",
			json:        `{"type": "AMRAP", "sets": 1, "minReps": 3}`,
			wantSets:    1,
			wantMinReps: 3,
			wantErr:     false,
		},
		{
			name:        "valid 2x8+",
			json:        `{"type": "AMRAP", "sets": 2, "minReps": 8}`,
			wantSets:    2,
			wantMinReps: 8,
			wantErr:     false,
		},
		{
			name:        "without type (still valid)",
			json:        `{"sets": 1, "minReps": 5}`,
			wantSets:    1,
			wantMinReps: 5,
			wantErr:     false,
		},
		{
			name:    "invalid zero sets",
			json:    `{"type": "AMRAP", "sets": 0, "minReps": 5}`,
			wantErr: true,
		},
		{
			name:    "invalid negative sets",
			json:    `{"type": "AMRAP", "sets": -1, "minReps": 5}`,
			wantErr: true,
		},
		{
			name:    "invalid zero minReps",
			json:    `{"type": "AMRAP", "sets": 1, "minReps": 0}`,
			wantErr: true,
		},
		{
			name:    "invalid negative minReps",
			json:    `{"type": "AMRAP", "sets": 1, "minReps": -1}`,
			wantErr: true,
		},
		{
			name:    "missing sets field",
			json:    `{"type": "AMRAP", "minReps": 5}`,
			wantErr: true,
		},
		{
			name:    "missing minReps field",
			json:    `{"type": "AMRAP", "sets": 1}`,
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
			scheme, err := UnmarshalAMRAPSetScheme(json.RawMessage(tt.json))
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

				amrap, ok := scheme.(*AMRAPSetScheme)
				if !ok {
					t.Fatal("expected *AMRAPSetScheme")
				}
				if amrap.Sets != tt.wantSets {
					t.Errorf("expected Sets %d, got %d", tt.wantSets, amrap.Sets)
				}
				if amrap.MinReps != tt.wantMinReps {
					t.Errorf("expected MinReps %d, got %d", tt.wantMinReps, amrap.MinReps)
				}
			}
		})
	}
}

// TestAMRAPSetScheme_RoundTrip tests JSON round-trip serialization.
func TestAMRAPSetScheme_RoundTrip(t *testing.T) {
	original := AMRAPSetScheme{Sets: 1, MinReps: 5}

	// Marshal
	data, err := json.Marshal(&original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	// Unmarshal
	scheme, err := UnmarshalAMRAPSetScheme(data)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	// Verify
	amrap, ok := scheme.(*AMRAPSetScheme)
	if !ok {
		t.Fatal("expected *AMRAPSetScheme")
	}
	if amrap.Sets != original.Sets {
		t.Errorf("expected Sets %d, got %d", original.Sets, amrap.Sets)
	}
	if amrap.MinReps != original.MinReps {
		t.Errorf("expected MinReps %d, got %d", original.MinReps, amrap.MinReps)
	}
}

// TestRegisterAMRAPScheme tests factory registration.
func TestRegisterAMRAPScheme(t *testing.T) {
	factory := NewSchemeFactory()

	// Should not be registered initially
	if factory.IsRegistered(TypeAMRAP) {
		t.Error("TypeAMRAP should not be registered initially")
	}

	// Register
	RegisterAMRAPScheme(factory)

	// Should be registered now
	if !factory.IsRegistered(TypeAMRAP) {
		t.Error("TypeAMRAP should be registered after RegisterAMRAPScheme")
	}
}

// TestAMRAPSetScheme_FactoryIntegration tests full factory workflow.
func TestAMRAPSetScheme_FactoryIntegration(t *testing.T) {
	factory := NewSchemeFactory()
	RegisterAMRAPScheme(factory)

	t.Run("Create from type and data", func(t *testing.T) {
		jsonData := json.RawMessage(`{"sets": 1, "minReps": 5}`)
		scheme, err := factory.Create(TypeAMRAP, jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if scheme.Type() != TypeAMRAP {
			t.Errorf("expected type %s, got %s", TypeAMRAP, scheme.Type())
		}
	})

	t.Run("CreateFromJSON", func(t *testing.T) {
		jsonData := json.RawMessage(`{"type": "AMRAP", "sets": 1, "minReps": 3}`)
		scheme, err := factory.CreateFromJSON(jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if scheme.Type() != TypeAMRAP {
			t.Errorf("expected type %s, got %s", TypeAMRAP, scheme.Type())
		}

		// Generate sets
		ctx := DefaultSetGenerationContext()
		sets, err := scheme.GenerateSets(335.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sets) != 1 {
			t.Errorf("expected 1 set, got %d", len(sets))
		}
		if sets[0].TargetReps != 3 {
			t.Errorf("expected 3 minReps, got %d", sets[0].TargetReps)
		}
		if sets[0].Weight != 335.0 {
			t.Errorf("expected weight 335.0, got %f", sets[0].Weight)
		}
	})

	t.Run("Invalid JSON in CreateFromJSON", func(t *testing.T) {
		jsonData := json.RawMessage(`{"type": "AMRAP", "sets": 0, "minReps": 5}`)
		_, err := factory.CreateFromJSON(jsonData)
		if err == nil {
			t.Error("expected error for invalid sets")
		}
	})
}

// TestAMRAPSetScheme_Implements_SetScheme verifies interface implementation.
func TestAMRAPSetScheme_Implements_SetScheme(t *testing.T) {
	var _ SetScheme = (*AMRAPSetScheme)(nil)
}

// TestAMRAPSetScheme_EdgeCases tests edge cases.
func TestAMRAPSetScheme_EdgeCases(t *testing.T) {
	t.Run("single set single minRep", func(t *testing.T) {
		scheme := AMRAPSetScheme{Sets: 1, MinReps: 1}
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

	t.Run("multiple AMRAP sets", func(t *testing.T) {
		scheme := AMRAPSetScheme{Sets: 3, MinReps: 5}
		ctx := DefaultSetGenerationContext()
		sets, err := scheme.GenerateSets(225.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sets) != 3 {
			t.Errorf("expected 3 sets, got %d", len(sets))
		}
		for i, set := range sets {
			if set.SetNumber != i+1 {
				t.Errorf("set %d: expected SetNumber %d, got %d", i, i+1, set.SetNumber)
			}
		}
	})

	t.Run("fractional weight preserved", func(t *testing.T) {
		scheme := AMRAPSetScheme{Sets: 1, MinReps: 5}
		ctx := DefaultSetGenerationContext()
		sets, err := scheme.GenerateSets(312.5, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sets[0].Weight != 312.5 {
			t.Errorf("expected Weight 312.5, got %f", sets[0].Weight)
		}
	})

	t.Run("very large weight", func(t *testing.T) {
		scheme := AMRAPSetScheme{Sets: 1, MinReps: 1}
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

// TestAMRAPSetScheme_WendlerExamples tests Wendler 5/3/1 program configurations.
func TestAMRAPSetScheme_WendlerExamples(t *testing.T) {
	tests := []struct {
		name       string
		week       string
		minReps    int
		baseWeight float64
	}{
		{"5/3/1 Week 1", "5+", 5, 285.0},
		{"5/3/1 Week 2", "3+", 3, 305.0},
		{"5/3/1 Week 3", "1+", 1, 325.0},
	}

	for _, tt := range tests {
		t.Run(tt.name+" "+tt.week, func(t *testing.T) {
			scheme, err := NewAMRAPSetScheme(1, tt.minReps)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			ctx := DefaultSetGenerationContext()
			sets, err := scheme.GenerateSets(tt.baseWeight, ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Verify single AMRAP set
			if len(sets) != 1 {
				t.Errorf("expected 1 set, got %d", len(sets))
			}

			set := sets[0]
			if set.SetNumber != 1 {
				t.Errorf("expected SetNumber 1, got %d", set.SetNumber)
			}
			if set.Weight != tt.baseWeight {
				t.Errorf("expected Weight %f, got %f", tt.baseWeight, set.Weight)
			}
			if set.TargetReps != tt.minReps {
				t.Errorf("expected TargetReps %d, got %d", tt.minReps, set.TargetReps)
			}
			if !set.IsWorkSet {
				t.Error("expected IsWorkSet true")
			}
		})
	}
}

// TestAMRAPSetScheme_ValidationErrorMessages tests error message clarity.
func TestAMRAPSetScheme_ValidationErrorMessages(t *testing.T) {
	t.Run("zero sets error message", func(t *testing.T) {
		scheme := AMRAPSetScheme{Sets: 0, MinReps: 5}
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
		scheme := AMRAPSetScheme{Sets: 1, MinReps: 0}
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
		scheme := AMRAPSetScheme{Sets: 0, MinReps: 0}
		err := scheme.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		// Should mention sets since that's checked first
	})
}

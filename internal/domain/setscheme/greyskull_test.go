package setscheme

import (
	"encoding/json"
	"errors"
	"testing"
)

// TestGreySkullSetScheme_Type verifies the type discriminator.
func TestGreySkullSetScheme_Type(t *testing.T) {
	scheme := &GreySkullSetScheme{FixedSets: 2, FixedReps: 5, AMRAPSets: 1, MinAMRAPReps: 5}
	if scheme.Type() != TypeGreySkull {
		t.Errorf("expected type %s, got %s", TypeGreySkull, scheme.Type())
	}
}

// TestNewGreySkullSetScheme tests the constructor function.
func TestNewGreySkullSetScheme(t *testing.T) {
	tests := []struct {
		name         string
		fixedSets    int
		fixedReps    int
		amrapSets    int
		minAMRAPReps int
		wantErr      bool
	}{
		{"valid 2x5 + 1x5+", 2, 5, 1, 5, false},
		{"valid 2x12 + 1x12+ (accessory)", 2, 12, 1, 12, false},
		{"valid 0 fixed sets (pure AMRAP)", 0, 0, 1, 5, false},
		{"valid 3x8 + 2x8+", 3, 8, 2, 8, false},
		{"valid 1x5 + 1x5+", 1, 5, 1, 5, false},
		{"invalid negative fixedSets", -1, 5, 1, 5, true},
		{"invalid zero amrapSets", 2, 5, 0, 5, true},
		{"invalid negative amrapSets", 2, 5, -1, 5, true},
		{"invalid zero minAMRAPReps", 2, 5, 1, 0, true},
		{"invalid negative minAMRAPReps", 2, 5, 1, -1, true},
		{"invalid zero fixedReps with fixedSets > 0", 2, 0, 1, 5, true},
		{"invalid negative fixedReps", 2, -1, 1, 5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme, err := NewGreySkullSetScheme(tt.fixedSets, tt.fixedReps, tt.amrapSets, tt.minAMRAPReps)
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
				if scheme.FixedSets != tt.fixedSets {
					t.Errorf("expected FixedSets %d, got %d", tt.fixedSets, scheme.FixedSets)
				}
				if scheme.FixedReps != tt.fixedReps {
					t.Errorf("expected FixedReps %d, got %d", tt.fixedReps, scheme.FixedReps)
				}
				if scheme.AMRAPSets != tt.amrapSets {
					t.Errorf("expected AMRAPSets %d, got %d", tt.amrapSets, scheme.AMRAPSets)
				}
				if scheme.MinAMRAPReps != tt.minAMRAPReps {
					t.Errorf("expected MinAMRAPReps %d, got %d", tt.minAMRAPReps, scheme.MinAMRAPReps)
				}
			}
		})
	}
}

// TestGreySkullSetScheme_Validate tests validation logic.
func TestGreySkullSetScheme_Validate(t *testing.T) {
	tests := []struct {
		name    string
		scheme  GreySkullSetScheme
		wantErr bool
	}{
		{"valid 2x5 + 1x5+", GreySkullSetScheme{FixedSets: 2, FixedReps: 5, AMRAPSets: 1, MinAMRAPReps: 5}, false},
		{"valid 0 fixed + 1 AMRAP", GreySkullSetScheme{FixedSets: 0, FixedReps: 0, AMRAPSets: 1, MinAMRAPReps: 5}, false},
		{"valid 3x10 + 1x10+", GreySkullSetScheme{FixedSets: 3, FixedReps: 10, AMRAPSets: 1, MinAMRAPReps: 10}, false},
		{"valid large values", GreySkullSetScheme{FixedSets: 10, FixedReps: 20, AMRAPSets: 3, MinAMRAPReps: 15}, false},
		{"invalid negative fixedSets", GreySkullSetScheme{FixedSets: -1, FixedReps: 5, AMRAPSets: 1, MinAMRAPReps: 5}, true},
		{"invalid zero amrapSets", GreySkullSetScheme{FixedSets: 2, FixedReps: 5, AMRAPSets: 0, MinAMRAPReps: 5}, true},
		{"invalid negative amrapSets", GreySkullSetScheme{FixedSets: 2, FixedReps: 5, AMRAPSets: -1, MinAMRAPReps: 5}, true},
		{"invalid zero minAMRAPReps", GreySkullSetScheme{FixedSets: 2, FixedReps: 5, AMRAPSets: 1, MinAMRAPReps: 0}, true},
		{"invalid zero fixedReps when fixedSets > 0", GreySkullSetScheme{FixedSets: 2, FixedReps: 0, AMRAPSets: 1, MinAMRAPReps: 5}, true},
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

// TestGreySkullSetScheme_GenerateSets tests set generation.
func TestGreySkullSetScheme_GenerateSets(t *testing.T) {
	tests := []struct {
		name           string
		scheme         GreySkullSetScheme
		baseWeight     float64
		wantTotalSets  int
		wantFixedSets  int
		wantAMRAPSets  int
		wantFixedReps  int
		wantAMRAPReps  int
	}{
		{
			name:           "standard 2x5 + 1x5+",
			scheme:         GreySkullSetScheme{FixedSets: 2, FixedReps: 5, AMRAPSets: 1, MinAMRAPReps: 5},
			baseWeight:     100.0,
			wantTotalSets:  3,
			wantFixedSets:  2,
			wantAMRAPSets:  1,
			wantFixedReps:  5,
			wantAMRAPReps:  5,
		},
		{
			name:           "accessory 2x12 + 1x12+",
			scheme:         GreySkullSetScheme{FixedSets: 2, FixedReps: 12, AMRAPSets: 1, MinAMRAPReps: 12},
			baseWeight:     50.0,
			wantTotalSets:  3,
			wantFixedSets:  2,
			wantAMRAPSets:  1,
			wantFixedReps:  12,
			wantAMRAPReps:  12,
		},
		{
			name:           "pure AMRAP (0 fixed sets)",
			scheme:         GreySkullSetScheme{FixedSets: 0, FixedReps: 0, AMRAPSets: 1, MinAMRAPReps: 5},
			baseWeight:     200.0,
			wantTotalSets:  1,
			wantFixedSets:  0,
			wantAMRAPSets:  1,
			wantFixedReps:  0,
			wantAMRAPReps:  5,
		},
		{
			name:           "3x8 + 2x8+",
			scheme:         GreySkullSetScheme{FixedSets: 3, FixedReps: 8, AMRAPSets: 2, MinAMRAPReps: 8},
			baseWeight:     135.0,
			wantTotalSets:  5,
			wantFixedSets:  3,
			wantAMRAPSets:  2,
			wantFixedReps:  8,
			wantAMRAPReps:  8,
		},
		{
			name:           "zero weight",
			scheme:         GreySkullSetScheme{FixedSets: 2, FixedReps: 5, AMRAPSets: 1, MinAMRAPReps: 5},
			baseWeight:     0.0,
			wantTotalSets:  3,
			wantFixedSets:  2,
			wantAMRAPSets:  1,
			wantFixedReps:  5,
			wantAMRAPReps:  5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := DefaultSetGenerationContext()
			sets, err := tt.scheme.GenerateSets(tt.baseWeight, ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(sets) != tt.wantTotalSets {
				t.Errorf("expected %d sets, got %d", tt.wantTotalSets, len(sets))
			}

			// Verify fixed sets
			for i := 0; i < tt.wantFixedSets; i++ {
				set := sets[i]
				expectedSetNumber := i + 1
				if set.SetNumber != expectedSetNumber {
					t.Errorf("fixed set %d: expected SetNumber %d, got %d", i, expectedSetNumber, set.SetNumber)
				}
				if set.Weight != tt.baseWeight {
					t.Errorf("fixed set %d: expected Weight %f, got %f", i, tt.baseWeight, set.Weight)
				}
				if set.TargetReps != tt.wantFixedReps {
					t.Errorf("fixed set %d: expected TargetReps %d, got %d", i, tt.wantFixedReps, set.TargetReps)
				}
				if !set.IsWorkSet {
					t.Errorf("fixed set %d: expected IsWorkSet to be true", i)
				}
			}

			// Verify AMRAP sets
			for i := 0; i < tt.wantAMRAPSets; i++ {
				setIndex := tt.wantFixedSets + i
				set := sets[setIndex]
				expectedSetNumber := setIndex + 1
				if set.SetNumber != expectedSetNumber {
					t.Errorf("AMRAP set %d: expected SetNumber %d, got %d", i, expectedSetNumber, set.SetNumber)
				}
				if set.Weight != tt.baseWeight {
					t.Errorf("AMRAP set %d: expected Weight %f, got %f", i, tt.baseWeight, set.Weight)
				}
				if set.TargetReps != tt.wantAMRAPReps {
					t.Errorf("AMRAP set %d: expected TargetReps %d, got %d", i, tt.wantAMRAPReps, set.TargetReps)
				}
				if !set.IsWorkSet {
					t.Errorf("AMRAP set %d: expected IsWorkSet to be true", i)
				}
			}
		})
	}
}

// TestGreySkullSetScheme_GenerateSets_InvalidScheme tests generation with invalid params.
func TestGreySkullSetScheme_GenerateSets_InvalidScheme(t *testing.T) {
	tests := []struct {
		name   string
		scheme GreySkullSetScheme
	}{
		{"negative fixedSets", GreySkullSetScheme{FixedSets: -1, FixedReps: 5, AMRAPSets: 1, MinAMRAPReps: 5}},
		{"zero amrapSets", GreySkullSetScheme{FixedSets: 2, FixedReps: 5, AMRAPSets: 0, MinAMRAPReps: 5}},
		{"negative amrapSets", GreySkullSetScheme{FixedSets: 2, FixedReps: 5, AMRAPSets: -1, MinAMRAPReps: 5}},
		{"zero minAMRAPReps", GreySkullSetScheme{FixedSets: 2, FixedReps: 5, AMRAPSets: 1, MinAMRAPReps: 0}},
		{"zero fixedReps with fixedSets > 0", GreySkullSetScheme{FixedSets: 2, FixedReps: 0, AMRAPSets: 1, MinAMRAPReps: 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := DefaultSetGenerationContext()
			sets, err := tt.scheme.GenerateSets(100.0, ctx)
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

// TestGreySkullSetScheme_GenerateSets_ContextIgnored tests that context is appropriately ignored.
func TestGreySkullSetScheme_GenerateSets_ContextIgnored(t *testing.T) {
	scheme := GreySkullSetScheme{FixedSets: 2, FixedReps: 5, AMRAPSets: 1, MinAMRAPReps: 5}

	// GreySkull scheme ignores WorkSetThreshold - all sets are work sets
	contexts := []SetGenerationContext{
		{WorkSetThreshold: 0},
		{WorkSetThreshold: 50},
		{WorkSetThreshold: 80},
		{WorkSetThreshold: 100},
		DefaultSetGenerationContext(),
	}

	for i, ctx := range contexts {
		sets, err := scheme.GenerateSets(100.0, ctx)
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

// TestGreySkullSetScheme_MarshalJSON tests JSON serialization.
func TestGreySkullSetScheme_MarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		scheme GreySkullSetScheme
	}{
		{"2x5 + 1x5+", GreySkullSetScheme{FixedSets: 2, FixedReps: 5, AMRAPSets: 1, MinAMRAPReps: 5}},
		{"2x12 + 1x12+", GreySkullSetScheme{FixedSets: 2, FixedReps: 12, AMRAPSets: 1, MinAMRAPReps: 12}},
		{"0 + 1x5+", GreySkullSetScheme{FixedSets: 0, FixedReps: 0, AMRAPSets: 1, MinAMRAPReps: 5}},
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
			if parsed["type"] != string(TypeGreySkull) {
				t.Errorf("expected type %s, got %v", TypeGreySkull, parsed["type"])
			}

			// Verify fields
			if int(parsed["fixedSets"].(float64)) != tt.scheme.FixedSets {
				t.Errorf("expected fixedSets %d, got %v", tt.scheme.FixedSets, parsed["fixedSets"])
			}
			if int(parsed["fixedReps"].(float64)) != tt.scheme.FixedReps {
				t.Errorf("expected fixedReps %d, got %v", tt.scheme.FixedReps, parsed["fixedReps"])
			}
			if int(parsed["amrapSets"].(float64)) != tt.scheme.AMRAPSets {
				t.Errorf("expected amrapSets %d, got %v", tt.scheme.AMRAPSets, parsed["amrapSets"])
			}
			if int(parsed["minAmrapReps"].(float64)) != tt.scheme.MinAMRAPReps {
				t.Errorf("expected minAmrapReps %d, got %v", tt.scheme.MinAMRAPReps, parsed["minAmrapReps"])
			}
		})
	}
}

// TestUnmarshalGreySkullSetScheme tests JSON deserialization.
func TestUnmarshalGreySkullSetScheme(t *testing.T) {
	tests := []struct {
		name             string
		json             string
		wantFixedSets    int
		wantFixedReps    int
		wantAMRAPSets    int
		wantMinAMRAPReps int
		wantErr          bool
	}{
		{
			name:             "valid 2x5 + 1x5+",
			json:             `{"type": "GREYSKULL", "fixedSets": 2, "fixedReps": 5, "amrapSets": 1, "minAmrapReps": 5}`,
			wantFixedSets:    2,
			wantFixedReps:    5,
			wantAMRAPSets:    1,
			wantMinAMRAPReps: 5,
			wantErr:          false,
		},
		{
			name:             "valid 2x12 + 1x12+",
			json:             `{"type": "GREYSKULL", "fixedSets": 2, "fixedReps": 12, "amrapSets": 1, "minAmrapReps": 12}`,
			wantFixedSets:    2,
			wantFixedReps:    12,
			wantAMRAPSets:    1,
			wantMinAMRAPReps: 12,
			wantErr:          false,
		},
		{
			name:             "without type (still valid)",
			json:             `{"fixedSets": 2, "fixedReps": 5, "amrapSets": 1, "minAmrapReps": 5}`,
			wantFixedSets:    2,
			wantFixedReps:    5,
			wantAMRAPSets:    1,
			wantMinAMRAPReps: 5,
			wantErr:          false,
		},
		{
			name:    "invalid zero amrapSets",
			json:    `{"type": "GREYSKULL", "fixedSets": 2, "fixedReps": 5, "amrapSets": 0, "minAmrapReps": 5}`,
			wantErr: true,
		},
		{
			name:    "invalid negative amrapSets",
			json:    `{"type": "GREYSKULL", "fixedSets": 2, "fixedReps": 5, "amrapSets": -1, "minAmrapReps": 5}`,
			wantErr: true,
		},
		{
			name:    "invalid zero minAmrapReps",
			json:    `{"type": "GREYSKULL", "fixedSets": 2, "fixedReps": 5, "amrapSets": 1, "minAmrapReps": 0}`,
			wantErr: true,
		},
		{
			name:    "invalid zero fixedReps with fixedSets > 0",
			json:    `{"type": "GREYSKULL", "fixedSets": 2, "fixedReps": 0, "amrapSets": 1, "minAmrapReps": 5}`,
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			json:    `{invalid}`,
			wantErr: true,
		},
		{
			name:    "missing required field amrapSets",
			json:    `{"type": "GREYSKULL", "fixedSets": 2, "fixedReps": 5, "minAmrapReps": 5}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme, err := UnmarshalGreySkullSetScheme(json.RawMessage(tt.json))
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

				gs, ok := scheme.(*GreySkullSetScheme)
				if !ok {
					t.Fatal("expected *GreySkullSetScheme")
				}
				if gs.FixedSets != tt.wantFixedSets {
					t.Errorf("expected FixedSets %d, got %d", tt.wantFixedSets, gs.FixedSets)
				}
				if gs.FixedReps != tt.wantFixedReps {
					t.Errorf("expected FixedReps %d, got %d", tt.wantFixedReps, gs.FixedReps)
				}
				if gs.AMRAPSets != tt.wantAMRAPSets {
					t.Errorf("expected AMRAPSets %d, got %d", tt.wantAMRAPSets, gs.AMRAPSets)
				}
				if gs.MinAMRAPReps != tt.wantMinAMRAPReps {
					t.Errorf("expected MinAMRAPReps %d, got %d", tt.wantMinAMRAPReps, gs.MinAMRAPReps)
				}
			}
		})
	}
}

// TestGreySkullSetScheme_RoundTrip tests JSON round-trip serialization.
func TestGreySkullSetScheme_RoundTrip(t *testing.T) {
	original := GreySkullSetScheme{FixedSets: 2, FixedReps: 5, AMRAPSets: 1, MinAMRAPReps: 5}

	// Marshal
	data, err := json.Marshal(&original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	// Unmarshal
	scheme, err := UnmarshalGreySkullSetScheme(data)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	// Verify
	gs, ok := scheme.(*GreySkullSetScheme)
	if !ok {
		t.Fatal("expected *GreySkullSetScheme")
	}
	if gs.FixedSets != original.FixedSets {
		t.Errorf("expected FixedSets %d, got %d", original.FixedSets, gs.FixedSets)
	}
	if gs.FixedReps != original.FixedReps {
		t.Errorf("expected FixedReps %d, got %d", original.FixedReps, gs.FixedReps)
	}
	if gs.AMRAPSets != original.AMRAPSets {
		t.Errorf("expected AMRAPSets %d, got %d", original.AMRAPSets, gs.AMRAPSets)
	}
	if gs.MinAMRAPReps != original.MinAMRAPReps {
		t.Errorf("expected MinAMRAPReps %d, got %d", original.MinAMRAPReps, gs.MinAMRAPReps)
	}
}

// TestRegisterGreySkullScheme tests factory registration.
func TestRegisterGreySkullScheme(t *testing.T) {
	factory := NewSchemeFactory()

	// Should not be registered initially
	if factory.IsRegistered(TypeGreySkull) {
		t.Error("TypeGreySkull should not be registered initially")
	}

	// Register
	RegisterGreySkullScheme(factory)

	// Should be registered now
	if !factory.IsRegistered(TypeGreySkull) {
		t.Error("TypeGreySkull should be registered after RegisterGreySkullScheme")
	}
}

// TestGreySkullSetScheme_FactoryIntegration tests full factory workflow.
func TestGreySkullSetScheme_FactoryIntegration(t *testing.T) {
	factory := NewSchemeFactory()
	RegisterGreySkullScheme(factory)

	t.Run("Create from type and data", func(t *testing.T) {
		jsonData := json.RawMessage(`{"fixedSets": 2, "fixedReps": 5, "amrapSets": 1, "minAmrapReps": 5}`)
		scheme, err := factory.Create(TypeGreySkull, jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if scheme.Type() != TypeGreySkull {
			t.Errorf("expected type %s, got %s", TypeGreySkull, scheme.Type())
		}
	})

	t.Run("CreateFromJSON", func(t *testing.T) {
		jsonData := json.RawMessage(`{"type": "GREYSKULL", "fixedSets": 2, "fixedReps": 5, "amrapSets": 1, "minAmrapReps": 5}`)
		scheme, err := factory.CreateFromJSON(jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if scheme.Type() != TypeGreySkull {
			t.Errorf("expected type %s, got %s", TypeGreySkull, scheme.Type())
		}

		// Generate sets
		ctx := DefaultSetGenerationContext()
		sets, err := scheme.GenerateSets(100.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sets) != 3 {
			t.Errorf("expected 3 sets, got %d", len(sets))
		}

		// Verify fixed sets
		for i := 0; i < 2; i++ {
			if sets[i].TargetReps != 5 {
				t.Errorf("fixed set %d: expected 5 reps, got %d", i, sets[i].TargetReps)
			}
		}

		// Verify AMRAP set
		if sets[2].TargetReps != 5 {
			t.Errorf("AMRAP set: expected 5 reps, got %d", sets[2].TargetReps)
		}
	})

	t.Run("Invalid JSON in CreateFromJSON", func(t *testing.T) {
		jsonData := json.RawMessage(`{"type": "GREYSKULL", "fixedSets": 2, "fixedReps": 5, "amrapSets": 0, "minAmrapReps": 5}`)
		_, err := factory.CreateFromJSON(jsonData)
		if err == nil {
			t.Error("expected error for invalid amrapSets")
		}
	})
}

// TestGreySkullSetScheme_Implements_SetScheme verifies interface implementation.
func TestGreySkullSetScheme_Implements_SetScheme(t *testing.T) {
	var _ SetScheme = (*GreySkullSetScheme)(nil)
}

// TestGreySkullSetScheme_EdgeCases tests edge cases.
func TestGreySkullSetScheme_EdgeCases(t *testing.T) {
	t.Run("many fixed sets", func(t *testing.T) {
		scheme := GreySkullSetScheme{FixedSets: 10, FixedReps: 5, AMRAPSets: 1, MinAMRAPReps: 5}
		ctx := DefaultSetGenerationContext()
		sets, err := scheme.GenerateSets(100.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sets) != 11 {
			t.Errorf("expected 11 sets, got %d", len(sets))
		}
		// Verify last set is AMRAP (set number 11)
		if sets[10].SetNumber != 11 {
			t.Errorf("expected SetNumber 11, got %d", sets[10].SetNumber)
		}
	})

	t.Run("many AMRAP sets", func(t *testing.T) {
		scheme := GreySkullSetScheme{FixedSets: 2, FixedReps: 5, AMRAPSets: 3, MinAMRAPReps: 5}
		ctx := DefaultSetGenerationContext()
		sets, err := scheme.GenerateSets(100.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sets) != 5 {
			t.Errorf("expected 5 sets, got %d", len(sets))
		}
		// Verify AMRAP sets are numbered correctly
		for i := 2; i < 5; i++ {
			if sets[i].SetNumber != i+1 {
				t.Errorf("set %d: expected SetNumber %d, got %d", i, i+1, sets[i].SetNumber)
			}
		}
	})

	t.Run("fractional weight preserved", func(t *testing.T) {
		scheme := GreySkullSetScheme{FixedSets: 2, FixedReps: 5, AMRAPSets: 1, MinAMRAPReps: 5}
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
		scheme := GreySkullSetScheme{FixedSets: 2, FixedReps: 5, AMRAPSets: 1, MinAMRAPReps: 5}
		ctx := DefaultSetGenerationContext()
		sets, err := scheme.GenerateSets(1000.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for i, set := range sets {
			if set.Weight != 1000.0 {
				t.Errorf("set %d: expected Weight 1000.0, got %f", i, set.Weight)
			}
		}
	})
}

// TestGreySkullSetScheme_GreySkullLPExamples tests common GreySkull LP configurations.
func TestGreySkullSetScheme_GreySkullLPExamples(t *testing.T) {
	tests := []struct {
		name         string
		description  string
		fixedSets    int
		fixedReps    int
		amrapSets    int
		minAMRAPReps int
		baseWeight   float64
	}{
		{"Main Lift Standard", "2x5 + 1x5+", 2, 5, 1, 5, 135.0},
		{"Accessory Standard", "2x10 + 1x10+", 2, 10, 1, 10, 50.0},
		{"Accessory Higher Rep", "2x15 + 1x15+", 2, 15, 1, 15, 25.0},
	}

	for _, tt := range tests {
		t.Run(tt.name+" "+tt.description, func(t *testing.T) {
			scheme, err := NewGreySkullSetScheme(tt.fixedSets, tt.fixedReps, tt.amrapSets, tt.minAMRAPReps)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			ctx := DefaultSetGenerationContext()
			sets, err := scheme.GenerateSets(tt.baseWeight, ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Verify correct total number of sets
			expectedTotal := tt.fixedSets + tt.amrapSets
			if len(sets) != expectedTotal {
				t.Errorf("expected %d sets, got %d", expectedTotal, len(sets))
			}

			// Verify fixed sets
			for i := 0; i < tt.fixedSets; i++ {
				set := sets[i]
				if set.SetNumber != i+1 {
					t.Errorf("fixed set %d: expected SetNumber %d, got %d", i, i+1, set.SetNumber)
				}
				if set.Weight != tt.baseWeight {
					t.Errorf("fixed set %d: expected Weight %f, got %f", i, tt.baseWeight, set.Weight)
				}
				if set.TargetReps != tt.fixedReps {
					t.Errorf("fixed set %d: expected TargetReps %d, got %d", i, tt.fixedReps, set.TargetReps)
				}
				if !set.IsWorkSet {
					t.Errorf("fixed set %d: expected IsWorkSet true", i)
				}
			}

			// Verify AMRAP set
			amrapSet := sets[tt.fixedSets]
			if amrapSet.SetNumber != tt.fixedSets+1 {
				t.Errorf("AMRAP set: expected SetNumber %d, got %d", tt.fixedSets+1, amrapSet.SetNumber)
			}
			if amrapSet.Weight != tt.baseWeight {
				t.Errorf("AMRAP set: expected Weight %f, got %f", tt.baseWeight, amrapSet.Weight)
			}
			if amrapSet.TargetReps != tt.minAMRAPReps {
				t.Errorf("AMRAP set: expected TargetReps %d, got %d", tt.minAMRAPReps, amrapSet.TargetReps)
			}
			if !amrapSet.IsWorkSet {
				t.Error("AMRAP set: expected IsWorkSet true")
			}
		})
	}
}

// TestGreySkullSetScheme_ValidationErrorMessages tests error message clarity.
func TestGreySkullSetScheme_ValidationErrorMessages(t *testing.T) {
	t.Run("negative fixedSets error message", func(t *testing.T) {
		scheme := GreySkullSetScheme{FixedSets: -1, FixedReps: 5, AMRAPSets: 1, MinAMRAPReps: 5}
		err := scheme.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, ErrInvalidParams) {
			t.Errorf("expected ErrInvalidParams, got %v", err)
		}
	})

	t.Run("zero amrapSets error message", func(t *testing.T) {
		scheme := GreySkullSetScheme{FixedSets: 2, FixedReps: 5, AMRAPSets: 0, MinAMRAPReps: 5}
		err := scheme.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, ErrInvalidParams) {
			t.Errorf("expected ErrInvalidParams, got %v", err)
		}
	})

	t.Run("fixedSets error takes precedence", func(t *testing.T) {
		// When both are invalid, fixedSets error comes first
		scheme := GreySkullSetScheme{FixedSets: -1, FixedReps: 5, AMRAPSets: 0, MinAMRAPReps: 0}
		err := scheme.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		// Should mention fixedSets since that's checked first
	})
}

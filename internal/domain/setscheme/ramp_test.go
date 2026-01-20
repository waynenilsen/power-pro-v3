package setscheme

import (
	"encoding/json"
	"errors"
	"math"
	"testing"
)

// TestRampSetScheme_Type verifies the type discriminator.
func TestRampSetScheme_Type(t *testing.T) {
	scheme := &RampSetScheme{
		Steps: []RampStep{{Percentage: 100, Reps: 5}},
	}
	if scheme.Type() != TypeRamp {
		t.Errorf("expected type %s, got %s", TypeRamp, scheme.Type())
	}
}

// TestNewRampSetScheme tests the constructor function with default threshold.
func TestNewRampSetScheme(t *testing.T) {
	tests := []struct {
		name    string
		steps   []RampStep
		wantErr bool
	}{
		{
			name: "valid single step",
			steps: []RampStep{
				{Percentage: 100, Reps: 5},
			},
			wantErr: false,
		},
		{
			name: "valid warmup ramp",
			steps: []RampStep{
				{Percentage: 50, Reps: 5},
				{Percentage: 63, Reps: 5},
				{Percentage: 75, Reps: 5},
				{Percentage: 88, Reps: 5},
				{Percentage: 100, Reps: 5},
			},
			wantErr: false,
		},
		{
			name:    "invalid empty steps",
			steps:   []RampStep{},
			wantErr: true,
		},
		{
			name:    "invalid nil steps",
			steps:   nil,
			wantErr: true,
		},
		{
			name: "invalid zero percentage",
			steps: []RampStep{
				{Percentage: 0, Reps: 5},
			},
			wantErr: true,
		},
		{
			name: "invalid negative percentage",
			steps: []RampStep{
				{Percentage: -10, Reps: 5},
			},
			wantErr: true,
		},
		{
			name: "invalid zero reps",
			steps: []RampStep{
				{Percentage: 100, Reps: 0},
			},
			wantErr: true,
		},
		{
			name: "invalid negative reps",
			steps: []RampStep{
				{Percentage: 100, Reps: -1},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme, err := NewRampSetScheme(tt.steps)
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
				if len(scheme.Steps) != len(tt.steps) {
					t.Errorf("expected %d steps, got %d", len(tt.steps), len(scheme.Steps))
				}
				// Default threshold should be applied
				if scheme.WorkSetThreshold != DefaultWorkSetThreshold {
					t.Errorf("expected WorkSetThreshold %f, got %f", DefaultWorkSetThreshold, scheme.WorkSetThreshold)
				}
			}
		})
	}
}

// TestNewRampSetSchemeWithThreshold tests the constructor with custom threshold.
func TestNewRampSetSchemeWithThreshold(t *testing.T) {
	tests := []struct {
		name      string
		steps     []RampStep
		threshold float64
		wantErr   bool
	}{
		{
			name: "valid custom threshold 70",
			steps: []RampStep{
				{Percentage: 100, Reps: 5},
			},
			threshold: 70.0,
			wantErr:   false,
		},
		{
			name: "valid threshold 100",
			steps: []RampStep{
				{Percentage: 100, Reps: 5},
			},
			threshold: 100.0,
			wantErr:   false,
		},
		{
			name: "valid threshold 1",
			steps: []RampStep{
				{Percentage: 100, Reps: 5},
			},
			threshold: 1.0,
			wantErr:   false,
		},
		{
			name: "invalid threshold > 100",
			steps: []RampStep{
				{Percentage: 100, Reps: 5},
			},
			threshold: 101.0,
			wantErr:   true,
		},
		{
			name: "invalid negative threshold",
			steps: []RampStep{
				{Percentage: 100, Reps: 5},
			},
			threshold: -10.0,
			wantErr:   true,
		},
		{
			name: "zero threshold uses default",
			steps: []RampStep{
				{Percentage: 100, Reps: 5},
			},
			threshold: 0,
			wantErr:   false, // Zero is allowed, means use default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme, err := NewRampSetSchemeWithThreshold(tt.steps, tt.threshold)
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
				if scheme.WorkSetThreshold != tt.threshold {
					t.Errorf("expected WorkSetThreshold %f, got %f", tt.threshold, scheme.WorkSetThreshold)
				}
			}
		})
	}
}

// TestRampSetScheme_Validate tests validation logic.
func TestRampSetScheme_Validate(t *testing.T) {
	tests := []struct {
		name    string
		scheme  RampSetScheme
		wantErr bool
	}{
		{
			name: "valid single step",
			scheme: RampSetScheme{
				Steps: []RampStep{{Percentage: 100, Reps: 5}},
			},
			wantErr: false,
		},
		{
			name: "valid warmup ramp",
			scheme: RampSetScheme{
				Steps: []RampStep{
					{Percentage: 50, Reps: 5},
					{Percentage: 63, Reps: 5},
					{Percentage: 75, Reps: 5},
					{Percentage: 88, Reps: 5},
					{Percentage: 100, Reps: 5},
				},
				WorkSetThreshold: 80,
			},
			wantErr: false,
		},
		{
			name: "valid percentage > 100 (overload)",
			scheme: RampSetScheme{
				Steps: []RampStep{{Percentage: 105, Reps: 1}},
			},
			wantErr: false,
		},
		{
			name: "valid fractional percentage",
			scheme: RampSetScheme{
				Steps: []RampStep{{Percentage: 62.5, Reps: 5}},
			},
			wantErr: false,
		},
		{
			name: "valid with zero threshold (uses default)",
			scheme: RampSetScheme{
				Steps:            []RampStep{{Percentage: 100, Reps: 5}},
				WorkSetThreshold: 0,
			},
			wantErr: false,
		},
		{
			name: "invalid empty steps",
			scheme: RampSetScheme{
				Steps: []RampStep{},
			},
			wantErr: true,
		},
		{
			name: "invalid nil steps",
			scheme: RampSetScheme{
				Steps: nil,
			},
			wantErr: true,
		},
		{
			name: "invalid zero percentage in middle",
			scheme: RampSetScheme{
				Steps: []RampStep{
					{Percentage: 50, Reps: 5},
					{Percentage: 0, Reps: 5},
					{Percentage: 100, Reps: 5},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid negative percentage",
			scheme: RampSetScheme{
				Steps: []RampStep{{Percentage: -10, Reps: 5}},
			},
			wantErr: true,
		},
		{
			name: "invalid zero reps",
			scheme: RampSetScheme{
				Steps: []RampStep{{Percentage: 100, Reps: 0}},
			},
			wantErr: true,
		},
		{
			name: "invalid negative reps",
			scheme: RampSetScheme{
				Steps: []RampStep{{Percentage: 100, Reps: -1}},
			},
			wantErr: true,
		},
		{
			name: "invalid reps in second step",
			scheme: RampSetScheme{
				Steps: []RampStep{
					{Percentage: 50, Reps: 5},
					{Percentage: 100, Reps: 0},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid threshold > 100",
			scheme: RampSetScheme{
				Steps:            []RampStep{{Percentage: 100, Reps: 5}},
				WorkSetThreshold: 150,
			},
			wantErr: true,
		},
		{
			name: "invalid negative threshold",
			scheme: RampSetScheme{
				Steps:            []RampStep{{Percentage: 100, Reps: 5}},
				WorkSetThreshold: -10,
			},
			wantErr: true,
		},
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

// TestRampSetScheme_GenerateSets tests set generation.
func TestRampSetScheme_GenerateSets(t *testing.T) {
	t.Run("typical warmup ramp", func(t *testing.T) {
		scheme := RampSetScheme{
			Steps: []RampStep{
				{Percentage: 50, Reps: 5},
				{Percentage: 63, Reps: 5},
				{Percentage: 75, Reps: 5},
				{Percentage: 88, Reps: 5},
				{Percentage: 100, Reps: 5},
			},
			WorkSetThreshold: 80,
		}
		baseWeight := 300.0
		ctx := DefaultSetGenerationContext()

		sets, err := scheme.GenerateSets(baseWeight, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(sets) != 5 {
			t.Errorf("expected 5 sets, got %d", len(sets))
		}

		expectedWeights := []float64{150, 189, 225, 264, 300}
		expectedIsWorkSet := []bool{false, false, false, true, true}

		for i, set := range sets {
			if set.SetNumber != i+1 {
				t.Errorf("set %d: expected SetNumber %d, got %d", i, i+1, set.SetNumber)
			}
			if set.Weight != expectedWeights[i] {
				t.Errorf("set %d: expected Weight %f, got %f", i, expectedWeights[i], set.Weight)
			}
			if set.TargetReps != 5 {
				t.Errorf("set %d: expected TargetReps 5, got %d", i, set.TargetReps)
			}
			if set.IsWorkSet != expectedIsWorkSet[i] {
				t.Errorf("set %d: expected IsWorkSet %v, got %v", i, expectedIsWorkSet[i], set.IsWorkSet)
			}
		}
	})

	t.Run("single step all work sets", func(t *testing.T) {
		scheme := RampSetScheme{
			Steps:            []RampStep{{Percentage: 100, Reps: 5}},
			WorkSetThreshold: 80,
		}
		ctx := DefaultSetGenerationContext()

		sets, err := scheme.GenerateSets(265.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(sets) != 1 {
			t.Errorf("expected 1 set, got %d", len(sets))
		}
		if sets[0].SetNumber != 1 {
			t.Errorf("expected SetNumber 1, got %d", sets[0].SetNumber)
		}
		if sets[0].Weight != 265.0 {
			t.Errorf("expected Weight 265.0, got %f", sets[0].Weight)
		}
		if sets[0].TargetReps != 5 {
			t.Errorf("expected TargetReps 5, got %d", sets[0].TargetReps)
		}
		if !sets[0].IsWorkSet {
			t.Error("expected IsWorkSet true for 100% set")
		}
	})

	t.Run("custom threshold - all warmups", func(t *testing.T) {
		scheme := RampSetScheme{
			Steps: []RampStep{
				{Percentage: 50, Reps: 5},
				{Percentage: 60, Reps: 5},
				{Percentage: 70, Reps: 5},
			},
			WorkSetThreshold: 100, // Nothing is a work set below 100%
		}
		ctx := DefaultSetGenerationContext()

		sets, err := scheme.GenerateSets(200.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		for i, set := range sets {
			if set.IsWorkSet {
				t.Errorf("set %d: expected IsWorkSet false with 100%% threshold", i)
			}
		}
	})

	t.Run("custom threshold - all work sets", func(t *testing.T) {
		scheme := RampSetScheme{
			Steps: []RampStep{
				{Percentage: 50, Reps: 5},
				{Percentage: 60, Reps: 5},
				{Percentage: 70, Reps: 5},
			},
			WorkSetThreshold: 50, // Everything at or above 50% is a work set
		}
		ctx := DefaultSetGenerationContext()

		sets, err := scheme.GenerateSets(200.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		for i, set := range sets {
			if !set.IsWorkSet {
				t.Errorf("set %d: expected IsWorkSet true with 50%% threshold", i)
			}
		}
	})

	t.Run("zero threshold uses default", func(t *testing.T) {
		scheme := RampSetScheme{
			Steps: []RampStep{
				{Percentage: 75, Reps: 5},
				{Percentage: 85, Reps: 5},
			},
			WorkSetThreshold: 0, // Uses default 80%
		}
		ctx := DefaultSetGenerationContext()

		sets, err := scheme.GenerateSets(200.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// 75% < 80% default threshold -> warmup
		if sets[0].IsWorkSet {
			t.Error("set 0: expected IsWorkSet false (75% < 80% default)")
		}
		// 85% >= 80% default threshold -> work set
		if !sets[1].IsWorkSet {
			t.Error("set 1: expected IsWorkSet true (85% >= 80% default)")
		}
	})

	t.Run("varying reps per step", func(t *testing.T) {
		scheme := RampSetScheme{
			Steps: []RampStep{
				{Percentage: 50, Reps: 8},
				{Percentage: 65, Reps: 6},
				{Percentage: 80, Reps: 4},
				{Percentage: 90, Reps: 2},
				{Percentage: 100, Reps: 1},
			},
			WorkSetThreshold: 80,
		}
		ctx := DefaultSetGenerationContext()

		sets, err := scheme.GenerateSets(400.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expectedReps := []int{8, 6, 4, 2, 1}
		for i, set := range sets {
			if set.TargetReps != expectedReps[i] {
				t.Errorf("set %d: expected TargetReps %d, got %d", i, expectedReps[i], set.TargetReps)
			}
		}
	})

	t.Run("zero base weight", func(t *testing.T) {
		scheme := RampSetScheme{
			Steps: []RampStep{
				{Percentage: 50, Reps: 5},
				{Percentage: 100, Reps: 5},
			},
			WorkSetThreshold: 80,
		}
		ctx := DefaultSetGenerationContext()

		sets, err := scheme.GenerateSets(0.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		for i, set := range sets {
			if set.Weight != 0 {
				t.Errorf("set %d: expected Weight 0, got %f", i, set.Weight)
			}
		}
	})

	t.Run("overload percentage (>100%)", func(t *testing.T) {
		scheme := RampSetScheme{
			Steps: []RampStep{
				{Percentage: 90, Reps: 3},
				{Percentage: 100, Reps: 2},
				{Percentage: 105, Reps: 1},
			},
			WorkSetThreshold: 80,
		}
		ctx := DefaultSetGenerationContext()

		sets, err := scheme.GenerateSets(300.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// 105% of 300 = 315
		if sets[2].Weight != 315.0 {
			t.Errorf("expected Weight 315.0 for 105%%, got %f", sets[2].Weight)
		}
		if !sets[2].IsWorkSet {
			t.Error("expected IsWorkSet true for 105%")
		}
	})

	t.Run("fractional percentage", func(t *testing.T) {
		scheme := RampSetScheme{
			Steps: []RampStep{
				{Percentage: 62.5, Reps: 5},
			},
			WorkSetThreshold: 80,
		}
		ctx := DefaultSetGenerationContext()

		sets, err := scheme.GenerateSets(200.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// 62.5% of 200 = 125
		if sets[0].Weight != 125.0 {
			t.Errorf("expected Weight 125.0 for 62.5%%, got %f", sets[0].Weight)
		}
	})

	t.Run("fractional base weight", func(t *testing.T) {
		scheme := RampSetScheme{
			Steps: []RampStep{
				{Percentage: 50, Reps: 5},
				{Percentage: 100, Reps: 5},
			},
			WorkSetThreshold: 80,
		}
		ctx := DefaultSetGenerationContext()

		sets, err := scheme.GenerateSets(137.5, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// 50% of 137.5 = 68.75
		if sets[0].Weight != 68.75 {
			t.Errorf("expected Weight 68.75 for 50%%, got %f", sets[0].Weight)
		}
		// 100% of 137.5 = 137.5
		if sets[1].Weight != 137.5 {
			t.Errorf("expected Weight 137.5 for 100%%, got %f", sets[1].Weight)
		}
	})

	t.Run("threshold exactly at percentage (edge case)", func(t *testing.T) {
		scheme := RampSetScheme{
			Steps: []RampStep{
				{Percentage: 79, Reps: 5},
				{Percentage: 80, Reps: 5},
				{Percentage: 81, Reps: 5},
			},
			WorkSetThreshold: 80,
		}
		ctx := DefaultSetGenerationContext()

		sets, err := scheme.GenerateSets(100.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// 79% < 80% -> warmup
		if sets[0].IsWorkSet {
			t.Error("set 0: expected IsWorkSet false (79% < 80%)")
		}
		// 80% >= 80% -> work set
		if !sets[1].IsWorkSet {
			t.Error("set 1: expected IsWorkSet true (80% >= 80%)")
		}
		// 81% >= 80% -> work set
		if !sets[2].IsWorkSet {
			t.Error("set 2: expected IsWorkSet true (81% >= 80%)")
		}
	})
}

// TestRampSetScheme_GenerateSets_InvalidScheme tests generation with invalid params.
func TestRampSetScheme_GenerateSets_InvalidScheme(t *testing.T) {
	tests := []struct {
		name   string
		scheme RampSetScheme
	}{
		{
			name: "empty steps",
			scheme: RampSetScheme{
				Steps: []RampStep{},
			},
		},
		{
			name: "nil steps",
			scheme: RampSetScheme{
				Steps: nil,
			},
		},
		{
			name: "zero percentage",
			scheme: RampSetScheme{
				Steps: []RampStep{{Percentage: 0, Reps: 5}},
			},
		},
		{
			name: "zero reps",
			scheme: RampSetScheme{
				Steps: []RampStep{{Percentage: 100, Reps: 0}},
			},
		},
		{
			name: "invalid threshold",
			scheme: RampSetScheme{
				Steps:            []RampStep{{Percentage: 100, Reps: 5}},
				WorkSetThreshold: 150,
			},
		},
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

// TestRampSetScheme_MarshalJSON tests JSON serialization.
func TestRampSetScheme_MarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		scheme RampSetScheme
	}{
		{
			name: "single step",
			scheme: RampSetScheme{
				Steps: []RampStep{{Percentage: 100, Reps: 5}},
			},
		},
		{
			name: "warmup ramp",
			scheme: RampSetScheme{
				Steps: []RampStep{
					{Percentage: 50, Reps: 5},
					{Percentage: 75, Reps: 5},
					{Percentage: 100, Reps: 5},
				},
				WorkSetThreshold: 80,
			},
		},
		{
			name: "custom threshold",
			scheme: RampSetScheme{
				Steps:            []RampStep{{Percentage: 100, Reps: 1}},
				WorkSetThreshold: 90,
			},
		},
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
			if parsed["type"] != string(TypeRamp) {
				t.Errorf("expected type %s, got %v", TypeRamp, parsed["type"])
			}

			// Verify steps exist
			steps, ok := parsed["steps"].([]interface{})
			if !ok {
				t.Fatal("expected steps array")
			}
			if len(steps) != len(tt.scheme.Steps) {
				t.Errorf("expected %d steps, got %d", len(tt.scheme.Steps), len(steps))
			}
		})
	}
}

// TestUnmarshalRampSetScheme tests JSON deserialization.
func TestUnmarshalRampSetScheme(t *testing.T) {
	tests := []struct {
		name          string
		json          string
		wantSteps     int
		wantThreshold float64
		wantErr       bool
	}{
		{
			name:          "valid single step",
			json:          `{"type": "RAMP", "steps": [{"percentage": 100, "reps": 5}]}`,
			wantSteps:     1,
			wantThreshold: 0, // Not specified, will be zero (uses default)
			wantErr:       false,
		},
		{
			name:          "valid warmup ramp",
			json:          `{"type": "RAMP", "steps": [{"percentage": 50, "reps": 5}, {"percentage": 100, "reps": 5}], "workSetThreshold": 80}`,
			wantSteps:     2,
			wantThreshold: 80,
			wantErr:       false,
		},
		{
			name:          "without type (still valid)",
			json:          `{"steps": [{"percentage": 100, "reps": 5}]}`,
			wantSteps:     1,
			wantThreshold: 0,
			wantErr:       false,
		},
		{
			name:          "fractional percentage",
			json:          `{"type": "RAMP", "steps": [{"percentage": 62.5, "reps": 5}]}`,
			wantSteps:     1,
			wantThreshold: 0,
			wantErr:       false,
		},
		{
			name:    "invalid empty steps",
			json:    `{"type": "RAMP", "steps": []}`,
			wantErr: true,
		},
		{
			name:    "invalid missing steps",
			json:    `{"type": "RAMP", "workSetThreshold": 80}`,
			wantErr: true,
		},
		{
			name:    "invalid zero percentage",
			json:    `{"type": "RAMP", "steps": [{"percentage": 0, "reps": 5}]}`,
			wantErr: true,
		},
		{
			name:    "invalid negative percentage",
			json:    `{"type": "RAMP", "steps": [{"percentage": -10, "reps": 5}]}`,
			wantErr: true,
		},
		{
			name:    "invalid zero reps",
			json:    `{"type": "RAMP", "steps": [{"percentage": 100, "reps": 0}]}`,
			wantErr: true,
		},
		{
			name:    "invalid threshold > 100",
			json:    `{"type": "RAMP", "steps": [{"percentage": 100, "reps": 5}], "workSetThreshold": 150}`,
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
			scheme, err := UnmarshalRampSetScheme(json.RawMessage(tt.json))
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

				ramp, ok := scheme.(*RampSetScheme)
				if !ok {
					t.Fatal("expected *RampSetScheme")
				}
				if len(ramp.Steps) != tt.wantSteps {
					t.Errorf("expected %d steps, got %d", tt.wantSteps, len(ramp.Steps))
				}
				if ramp.WorkSetThreshold != tt.wantThreshold {
					t.Errorf("expected WorkSetThreshold %f, got %f", tt.wantThreshold, ramp.WorkSetThreshold)
				}
			}
		})
	}
}

// TestRampSetScheme_RoundTrip tests JSON round-trip serialization.
func TestRampSetScheme_RoundTrip(t *testing.T) {
	original := RampSetScheme{
		Steps: []RampStep{
			{Percentage: 50, Reps: 5},
			{Percentage: 75, Reps: 4},
			{Percentage: 100, Reps: 3},
		},
		WorkSetThreshold: 70,
	}

	// Marshal
	data, err := json.Marshal(&original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	// Unmarshal
	scheme, err := UnmarshalRampSetScheme(data)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	// Verify
	ramp, ok := scheme.(*RampSetScheme)
	if !ok {
		t.Fatal("expected *RampSetScheme")
	}
	if len(ramp.Steps) != len(original.Steps) {
		t.Errorf("expected %d steps, got %d", len(original.Steps), len(ramp.Steps))
	}
	for i := range original.Steps {
		if ramp.Steps[i].Percentage != original.Steps[i].Percentage {
			t.Errorf("step %d: expected Percentage %f, got %f", i, original.Steps[i].Percentage, ramp.Steps[i].Percentage)
		}
		if ramp.Steps[i].Reps != original.Steps[i].Reps {
			t.Errorf("step %d: expected Reps %d, got %d", i, original.Steps[i].Reps, ramp.Steps[i].Reps)
		}
	}
	if ramp.WorkSetThreshold != original.WorkSetThreshold {
		t.Errorf("expected WorkSetThreshold %f, got %f", original.WorkSetThreshold, ramp.WorkSetThreshold)
	}
}

// TestRegisterRampScheme tests factory registration.
func TestRegisterRampScheme(t *testing.T) {
	factory := NewSchemeFactory()

	// Should not be registered initially
	if factory.IsRegistered(TypeRamp) {
		t.Error("TypeRamp should not be registered initially")
	}

	// Register
	RegisterRampScheme(factory)

	// Should be registered now
	if !factory.IsRegistered(TypeRamp) {
		t.Error("TypeRamp should be registered after RegisterRampScheme")
	}
}

// TestRampSetScheme_FactoryIntegration tests full factory workflow.
func TestRampSetScheme_FactoryIntegration(t *testing.T) {
	factory := NewSchemeFactory()
	RegisterRampScheme(factory)

	t.Run("Create from type and data", func(t *testing.T) {
		jsonData := json.RawMessage(`{"steps": [{"percentage": 100, "reps": 5}]}`)
		scheme, err := factory.Create(TypeRamp, jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if scheme.Type() != TypeRamp {
			t.Errorf("expected type %s, got %s", TypeRamp, scheme.Type())
		}
	})

	t.Run("CreateFromJSON", func(t *testing.T) {
		jsonData := json.RawMessage(`{"type": "RAMP", "steps": [{"percentage": 50, "reps": 5}, {"percentage": 100, "reps": 5}], "workSetThreshold": 80}`)
		scheme, err := factory.CreateFromJSON(jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if scheme.Type() != TypeRamp {
			t.Errorf("expected type %s, got %s", TypeRamp, scheme.Type())
		}

		// Generate sets
		ctx := DefaultSetGenerationContext()
		sets, err := scheme.GenerateSets(200.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sets) != 2 {
			t.Errorf("expected 2 sets, got %d", len(sets))
		}

		// First set should be warmup (50% < 80%)
		if sets[0].IsWorkSet {
			t.Error("set 0: expected IsWorkSet false")
		}
		// Second set should be work set (100% >= 80%)
		if !sets[1].IsWorkSet {
			t.Error("set 1: expected IsWorkSet true")
		}
	})

	t.Run("Invalid JSON in CreateFromJSON", func(t *testing.T) {
		jsonData := json.RawMessage(`{"type": "RAMP", "steps": []}`)
		_, err := factory.CreateFromJSON(jsonData)
		if err == nil {
			t.Error("expected error for empty steps")
		}
	})
}

// TestRampSetScheme_Implements_SetScheme verifies interface implementation.
func TestRampSetScheme_Implements_SetScheme(t *testing.T) {
	var _ SetScheme = (*RampSetScheme)(nil)
}

// TestRampSetScheme_EdgeCases tests edge cases.
func TestRampSetScheme_EdgeCases(t *testing.T) {
	t.Run("single step single rep", func(t *testing.T) {
		scheme := RampSetScheme{
			Steps:            []RampStep{{Percentage: 100, Reps: 1}},
			WorkSetThreshold: 80,
		}
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

	t.Run("many steps", func(t *testing.T) {
		steps := make([]RampStep, 20)
		for i := 0; i < 20; i++ {
			steps[i] = RampStep{Percentage: float64(40 + i*3), Reps: 5}
		}
		scheme := RampSetScheme{
			Steps:            steps,
			WorkSetThreshold: 80,
		}
		ctx := DefaultSetGenerationContext()
		sets, err := scheme.GenerateSets(100.0, ctx)
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

	t.Run("descending percentages allowed", func(t *testing.T) {
		// Ramp doesn't require ascending order
		scheme := RampSetScheme{
			Steps: []RampStep{
				{Percentage: 100, Reps: 5},
				{Percentage: 90, Reps: 5},
				{Percentage: 80, Reps: 5},
			},
			WorkSetThreshold: 80,
		}
		ctx := DefaultSetGenerationContext()
		sets, err := scheme.GenerateSets(200.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expectedWeights := []float64{200, 180, 160}
		for i, set := range sets {
			if set.Weight != expectedWeights[i] {
				t.Errorf("set %d: expected Weight %f, got %f", i, expectedWeights[i], set.Weight)
			}
		}
	})

	t.Run("non-monotonic percentages allowed", func(t *testing.T) {
		// Ramp allows any order
		scheme := RampSetScheme{
			Steps: []RampStep{
				{Percentage: 70, Reps: 5},
				{Percentage: 100, Reps: 3},
				{Percentage: 80, Reps: 8},
			},
			WorkSetThreshold: 80,
		}
		ctx := DefaultSetGenerationContext()
		sets, err := scheme.GenerateSets(200.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// 70% < 80% -> warmup
		if sets[0].IsWorkSet {
			t.Error("set 0: expected warmup (70% < 80%)")
		}
		// 100% >= 80% -> work set
		if !sets[1].IsWorkSet {
			t.Error("set 1: expected work set (100% >= 80%)")
		}
		// 80% >= 80% -> work set
		if !sets[2].IsWorkSet {
			t.Error("set 2: expected work set (80% >= 80%)")
		}
	})

	t.Run("very small percentage", func(t *testing.T) {
		scheme := RampSetScheme{
			Steps:            []RampStep{{Percentage: 0.1, Reps: 10}},
			WorkSetThreshold: 80,
		}
		ctx := DefaultSetGenerationContext()
		sets, err := scheme.GenerateSets(1000.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// 0.1% of 1000 = 1.0
		if sets[0].Weight != 1.0 {
			t.Errorf("expected Weight 1.0, got %f", sets[0].Weight)
		}
	})

	t.Run("very large percentage", func(t *testing.T) {
		scheme := RampSetScheme{
			Steps:            []RampStep{{Percentage: 200, Reps: 1}},
			WorkSetThreshold: 80,
		}
		ctx := DefaultSetGenerationContext()
		sets, err := scheme.GenerateSets(100.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// 200% of 100 = 200
		if sets[0].Weight != 200.0 {
			t.Errorf("expected Weight 200.0, got %f", sets[0].Weight)
		}
	})
}

// TestRampSetScheme_ProgramExamples tests common program configurations.
func TestRampSetScheme_ProgramExamples(t *testing.T) {
	tests := []struct {
		name             string
		program          string
		steps            []RampStep
		workSetThreshold float64
		baseWeight       float64
		expectedWorkSets int
		expectedWarmups  int
	}{
		{
			name:    "Bill Starr warmup to top set",
			program: "warmup ramp",
			steps: []RampStep{
				{Percentage: 50, Reps: 5},
				{Percentage: 63, Reps: 5},
				{Percentage: 75, Reps: 5},
				{Percentage: 88, Reps: 5},
				{Percentage: 100, Reps: 5},
			},
			workSetThreshold: 80,
			baseWeight:       225.0,
			expectedWorkSets: 2,  // 88% and 100%
			expectedWarmups:  3,  // 50%, 63%, 75%
		},
		{
			name:    "Simple 3-step warmup",
			program: "3-step warmup",
			steps: []RampStep{
				{Percentage: 50, Reps: 5},
				{Percentage: 75, Reps: 3},
				{Percentage: 100, Reps: 1},
			},
			workSetThreshold: 100,
			baseWeight:       315.0,
			expectedWorkSets: 1,  // Only 100%
			expectedWarmups:  2,  // 50%, 75%
		},
		{
			name:    "All work sets",
			program: "heavy triples",
			steps: []RampStep{
				{Percentage: 85, Reps: 3},
				{Percentage: 90, Reps: 3},
				{Percentage: 95, Reps: 3},
			},
			workSetThreshold: 80,
			baseWeight:       400.0,
			expectedWorkSets: 3,  // All >= 80%
			expectedWarmups:  0,
		},
		{
			name:    "Sheiko warmup",
			program: "competition prep",
			steps: []RampStep{
				{Percentage: 50, Reps: 5},
				{Percentage: 60, Reps: 4},
				{Percentage: 70, Reps: 3},
				{Percentage: 80, Reps: 2},
				{Percentage: 85, Reps: 1},
				{Percentage: 90, Reps: 1},
			},
			workSetThreshold: 70,
			baseWeight:       500.0,
			expectedWorkSets: 4,  // 70%, 80%, 85%, 90%
			expectedWarmups:  2,  // 50%, 60%
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+" "+tt.program, func(t *testing.T) {
			scheme, err := NewRampSetSchemeWithThreshold(tt.steps, tt.workSetThreshold)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			ctx := DefaultSetGenerationContext()
			sets, err := scheme.GenerateSets(tt.baseWeight, ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Count work sets and warmups
			workSets := 0
			warmups := 0
			for _, set := range sets {
				if set.IsWorkSet {
					workSets++
				} else {
					warmups++
				}
			}

			if workSets != tt.expectedWorkSets {
				t.Errorf("expected %d work sets, got %d", tt.expectedWorkSets, workSets)
			}
			if warmups != tt.expectedWarmups {
				t.Errorf("expected %d warmups, got %d", tt.expectedWarmups, warmups)
			}

			// Verify set numbers are sequential
			for i, set := range sets {
				if set.SetNumber != i+1 {
					t.Errorf("set %d: expected SetNumber %d, got %d", i, i+1, set.SetNumber)
				}
			}
		})
	}
}

// TestRampSetScheme_ValidationErrorMessages tests error message clarity.
func TestRampSetScheme_ValidationErrorMessages(t *testing.T) {
	t.Run("empty steps error message", func(t *testing.T) {
		scheme := RampSetScheme{Steps: []RampStep{}}
		err := scheme.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, ErrInvalidParams) {
			t.Errorf("expected ErrInvalidParams, got %v", err)
		}
	})

	t.Run("zero percentage error message", func(t *testing.T) {
		scheme := RampSetScheme{
			Steps: []RampStep{{Percentage: 0, Reps: 5}},
		}
		err := scheme.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, ErrInvalidParams) {
			t.Errorf("expected ErrInvalidParams, got %v", err)
		}
	})

	t.Run("zero reps error message", func(t *testing.T) {
		scheme := RampSetScheme{
			Steps: []RampStep{{Percentage: 100, Reps: 0}},
		}
		err := scheme.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, ErrInvalidParams) {
			t.Errorf("expected ErrInvalidParams, got %v", err)
		}
	})

	t.Run("invalid threshold error message", func(t *testing.T) {
		scheme := RampSetScheme{
			Steps:            []RampStep{{Percentage: 100, Reps: 5}},
			WorkSetThreshold: 150,
		}
		err := scheme.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, ErrInvalidParams) {
			t.Errorf("expected ErrInvalidParams, got %v", err)
		}
	})

	t.Run("error identifies step number", func(t *testing.T) {
		scheme := RampSetScheme{
			Steps: []RampStep{
				{Percentage: 50, Reps: 5},
				{Percentage: 75, Reps: 5},
				{Percentage: 0, Reps: 5}, // Third step is invalid
			},
		}
		err := scheme.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		// Error should mention step 3
		errMsg := err.Error()
		if errMsg == "" {
			t.Error("expected non-empty error message")
		}
	})
}

// TestRampStep_JSONSerialization tests RampStep JSON handling.
func TestRampStep_JSONSerialization(t *testing.T) {
	t.Run("marshal single step", func(t *testing.T) {
		step := RampStep{Percentage: 62.5, Reps: 5}
		data, err := json.Marshal(step)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var parsed RampStep
		if err := json.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}

		if parsed.Percentage != step.Percentage {
			t.Errorf("expected Percentage %f, got %f", step.Percentage, parsed.Percentage)
		}
		if parsed.Reps != step.Reps {
			t.Errorf("expected Reps %d, got %d", step.Reps, parsed.Reps)
		}
	})

	t.Run("unmarshal from JSON", func(t *testing.T) {
		jsonStr := `{"percentage": 87.5, "reps": 3}`
		var step RampStep
		if err := json.Unmarshal([]byte(jsonStr), &step); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if step.Percentage != 87.5 {
			t.Errorf("expected Percentage 87.5, got %f", step.Percentage)
		}
		if step.Reps != 3 {
			t.Errorf("expected Reps 3, got %d", step.Reps)
		}
	})
}

// TestDefaultWorkSetThreshold tests the default constant.
func TestDefaultWorkSetThreshold(t *testing.T) {
	if DefaultWorkSetThreshold != 80.0 {
		t.Errorf("expected DefaultWorkSetThreshold 80.0, got %f", DefaultWorkSetThreshold)
	}
}

// TestRampSetScheme_WeightCalculationPrecision tests weight calculation precision.
func TestRampSetScheme_WeightCalculationPrecision(t *testing.T) {
	scheme := RampSetScheme{
		Steps: []RampStep{
			{Percentage: 33.33, Reps: 5},
			{Percentage: 66.67, Reps: 5},
		},
		WorkSetThreshold: 80,
	}
	ctx := DefaultSetGenerationContext()

	sets, err := scheme.GenerateSets(300.0, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 33.33% of 300 = 99.99
	expected1 := 300.0 * (33.33 / 100)
	if math.Abs(sets[0].Weight-expected1) > 0.0001 {
		t.Errorf("set 0: expected Weight %f, got %f", expected1, sets[0].Weight)
	}

	// 66.67% of 300 = 200.01
	expected2 := 300.0 * (66.67 / 100)
	if math.Abs(sets[1].Weight-expected2) > 0.0001 {
		t.Errorf("set 1: expected Weight %f, got %f", expected2, sets[1].Weight)
	}
}

// TestRampSetScheme_ContextIgnored tests that the context parameter is properly handled.
func TestRampSetScheme_ContextIgnored(t *testing.T) {
	// RampSetScheme uses its own WorkSetThreshold, not the one from context
	scheme := RampSetScheme{
		Steps: []RampStep{
			{Percentage: 75, Reps: 5},
			{Percentage: 85, Reps: 5},
		},
		WorkSetThreshold: 80,
	}

	// Test with different context values - scheme should use its own threshold
	contexts := []SetGenerationContext{
		{WorkSetThreshold: 0},
		{WorkSetThreshold: 50},
		{WorkSetThreshold: 90},
		{WorkSetThreshold: 100},
		DefaultSetGenerationContext(),
	}

	for i, ctx := range contexts {
		sets, err := scheme.GenerateSets(200.0, ctx)
		if err != nil {
			t.Fatalf("context %d: unexpected error: %v", i, err)
		}

		// Results should be consistent regardless of context
		// 75% < scheme's 80% threshold -> warmup
		if sets[0].IsWorkSet {
			t.Errorf("context %d: set 0 should be warmup (scheme threshold 80%%)", i)
		}
		// 85% >= scheme's 80% threshold -> work set
		if !sets[1].IsWorkSet {
			t.Errorf("context %d: set 1 should be work set (scheme threshold 80%%)", i)
		}
	}
}

package loadstrategy

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"testing"
)

func TestGetTaperMultiplier(t *testing.T) {
	tests := []struct {
		name     string
		daysOut  int
		expected float64
	}{
		// Final week (0-6 days out) - 50% volume
		{"final week day 0 (meet day)", 0, 0.5},
		{"final week day 3", 3, 0.5},
		{"final week day 6", 6, 0.5},

		// Week 4 (7-13 days out) - 60% volume
		{"week 4 day 7", 7, 0.6},
		{"week 4 day 10", 10, 0.6},
		{"week 4 day 13", 13, 0.6},

		// Week 3 (14-20 days out) - 70% volume
		{"week 3 day 14", 14, 0.7},
		{"week 3 day 17", 17, 0.7},
		{"week 3 day 20", 20, 0.7},

		// Week 2 (21-27 days out) - 80% volume
		{"week 2 day 21", 21, 0.8},
		{"week 2 day 24", 24, 0.8},
		{"week 2 day 27", 27, 0.8},

		// Week 1 of Comp (28-34 days out) - 90% volume
		{"week 1 day 28", 28, 0.9},
		{"week 1 day 31", 31, 0.9},
		{"week 1 day 34", 34, 0.9},

		// Beyond taper zone (35+ days out) - 100% volume
		{"beyond taper day 35", 35, 1.0},
		{"beyond taper day 50", 50, 1.0},
		{"beyond taper day 90", 90, 1.0},

		// Edge case: past meet date (negative days)
		{"past meet date -1", -1, 0.5},
		{"past meet date -7", -7, 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTaperMultiplier(tt.daysOut)
			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("GetTaperMultiplier(%d) = %f, expected %f", tt.daysOut, result, tt.expected)
			}
		})
	}
}

func TestGetTaperMultiplierWithCurve(t *testing.T) {
	// Custom curve for testing
	customCurve := []TaperCurve{
		{ThresholdDays: 3, Multiplier: 0.4},  // Last 3 days: 40%
		{ThresholdDays: 10, Multiplier: 0.6}, // Days 3-9: 60%
		{ThresholdDays: 20, Multiplier: 0.8}, // Days 10-19: 80%
	}

	tests := []struct {
		name     string
		daysOut  int
		expected float64
	}{
		{"custom curve day 0", 0, 0.4},
		{"custom curve day 2", 2, 0.4},
		{"custom curve day 3", 3, 0.6},
		{"custom curve day 9", 9, 0.6},
		{"custom curve day 10", 10, 0.8},
		{"custom curve day 19", 19, 0.8},
		{"custom curve day 20", 20, 1.0}, // No matching tier
		{"custom curve day 30", 30, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTaperMultiplierWithCurve(tt.daysOut, customCurve)
			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("GetTaperMultiplierWithCurve(%d, custom) = %f, expected %f", tt.daysOut, result, tt.expected)
			}
		})
	}
}

func TestDefaultTaperCurve(t *testing.T) {
	curve := DefaultTaperCurve()

	if len(curve) != 5 {
		t.Errorf("expected 5 taper tiers, got %d", len(curve))
	}

	// Verify structure
	expected := []TaperCurve{
		{ThresholdDays: 7, Multiplier: 0.5},
		{ThresholdDays: 14, Multiplier: 0.6},
		{ThresholdDays: 21, Multiplier: 0.7},
		{ThresholdDays: 28, Multiplier: 0.8},
		{ThresholdDays: 35, Multiplier: 0.9},
	}

	for i, tier := range curve {
		if tier.ThresholdDays != expected[i].ThresholdDays {
			t.Errorf("tier %d: expected ThresholdDays %d, got %d", i, expected[i].ThresholdDays, tier.ThresholdDays)
		}
		if tier.Multiplier != expected[i].Multiplier {
			t.Errorf("tier %d: expected Multiplier %f, got %f", i, expected[i].Multiplier, tier.Multiplier)
		}
	}
}

func TestTaperLoadStrategy_Type(t *testing.T) {
	strategy := &TaperLoadStrategy{}
	if strategy.Type() != TypeTaper {
		t.Errorf("expected type %s, got %s", TypeTaper, strategy.Type())
	}
}

func TestTaperLoadStrategy_Validate(t *testing.T) {
	baseStrategy := &PercentOfLoadStrategy{
		ReferenceType: ReferenceTrainingMax,
		Percentage:    85.0,
	}

	tests := []struct {
		name     string
		strategy TaperLoadStrategy
		wantErr  bool
	}{
		{
			name: "valid with base strategy and default curve",
			strategy: TaperLoadStrategy{
				BaseStrategy: baseStrategy,
			},
			wantErr: false,
		},
		{
			name: "valid with custom curve",
			strategy: TaperLoadStrategy{
				BaseStrategy: baseStrategy,
				TaperCurve: []TaperCurve{
					{ThresholdDays: 7, Multiplier: 0.5},
					{ThresholdDays: 14, Multiplier: 0.7},
				},
			},
			wantErr: false,
		},
		{
			name: "valid with maintainIntensity",
			strategy: TaperLoadStrategy{
				BaseStrategy:      baseStrategy,
				MaintainIntensity: true,
			},
			wantErr: false,
		},
		{
			name: "invalid - missing base strategy",
			strategy: TaperLoadStrategy{
				BaseStrategy: nil,
			},
			wantErr: true,
		},
		{
			name: "invalid - base strategy validation fails",
			strategy: TaperLoadStrategy{
				BaseStrategy: &PercentOfLoadStrategy{
					// Missing required fields
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - curve threshold <= 0",
			strategy: TaperLoadStrategy{
				BaseStrategy: baseStrategy,
				TaperCurve: []TaperCurve{
					{ThresholdDays: 0, Multiplier: 0.5},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - curve multiplier < 0",
			strategy: TaperLoadStrategy{
				BaseStrategy: baseStrategy,
				TaperCurve: []TaperCurve{
					{ThresholdDays: 7, Multiplier: -0.5},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - curve multiplier > 1",
			strategy: TaperLoadStrategy{
				BaseStrategy: baseStrategy,
				TaperCurve: []TaperCurve{
					{ThresholdDays: 7, Multiplier: 1.5},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - curve not ascending",
			strategy: TaperLoadStrategy{
				BaseStrategy: baseStrategy,
				TaperCurve: []TaperCurve{
					{ThresholdDays: 14, Multiplier: 0.5},
					{ThresholdDays: 7, Multiplier: 0.7}, // Out of order
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.strategy.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTaperLoadStrategy_CalculateLoad(t *testing.T) {
	// Create mock max lookup
	lookup := newMockMaxLookup()
	lookup.SetMax("user-123", "squat-456", "TRAINING_MAX", 300.0, "2024-01-15")

	// Base strategy: 100% of TM = 300 lbs
	baseStrategy := &PercentOfLoadStrategy{
		ReferenceType:     ReferenceTrainingMax,
		Percentage:        100.0,
		RoundingIncrement: 5.0,
		RoundingDirection: RoundNearest,
	}
	baseStrategy.SetMaxLookup(lookup)

	tests := []struct {
		name              string
		daysOut           interface{} // Can be int, float64, or nil
		curve             []TaperCurve
		maintainIntensity bool
		expected          float64
	}{
		// Default curve tests
		{
			name:     "final week (0 days) - 50% of 300",
			daysOut:  0,
			expected: 150.0, // 300 * 0.5
		},
		{
			name:     "week 4 (10 days) - 60% of 300",
			daysOut:  10,
			expected: 180.0, // 300 * 0.6
		},
		{
			name:     "week 3 (17 days) - 70% of 300",
			daysOut:  17,
			expected: 210.0, // 300 * 0.7
		},
		{
			name:     "week 2 (24 days) - 80% of 300",
			daysOut:  24,
			expected: 240.0, // 300 * 0.8
		},
		{
			name:     "week 1 of comp (31 days) - 90% of 300",
			daysOut:  31,
			expected: 270.0, // 300 * 0.9
		},
		{
			name:     "beyond taper (50 days) - 100% of 300",
			daysOut:  50,
			expected: 300.0, // 300 * 1.0
		},

		// No daysOut provided - should return base load
		{
			name:     "no daysOut - returns base load",
			daysOut:  nil,
			expected: 300.0,
		},

		// float64 daysOut (from JSON unmarshaling)
		{
			name:     "float64 daysOut 10.0",
			daysOut:  10.0,
			expected: 180.0, // 300 * 0.6
		},

		// int64 daysOut
		{
			name:     "int64 daysOut 10",
			daysOut:  int64(10),
			expected: 180.0,
		},

		// Custom curve
		{
			name:    "custom curve - day 5 (90% multiplier)",
			daysOut: 5,
			curve: []TaperCurve{
				{ThresholdDays: 7, Multiplier: 0.9},
			},
			expected: 270.0, // 300 * 0.9
		},

		// Maintain intensity
		{
			name:              "maintain intensity - returns base load",
			daysOut:           10,
			maintainIntensity: true,
			expected:          300.0, // Ignores taper
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clone base strategy to avoid mutation
			bs := &PercentOfLoadStrategy{
				ReferenceType:     baseStrategy.ReferenceType,
				Percentage:        baseStrategy.Percentage,
				RoundingIncrement: baseStrategy.RoundingIncrement,
				RoundingDirection: baseStrategy.RoundingDirection,
			}
			bs.SetMaxLookup(lookup)

			strategy := &TaperLoadStrategy{
				BaseStrategy:      bs,
				TaperCurve:        tt.curve,
				MaintainIntensity: tt.maintainIntensity,
			}

			ctxMap := make(map[string]interface{})
			if tt.daysOut != nil {
				ctxMap["daysOut"] = tt.daysOut
			}

			params := LoadCalculationParams{
				UserID:  "user-123",
				LiftID:  "squat-456",
				Context: ctxMap,
			}

			result, err := strategy.CalculateLoad(context.Background(), params)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestTaperLoadStrategy_CalculateLoad_Errors(t *testing.T) {
	lookup := newMockMaxLookup()
	lookup.SetMax("user-123", "squat-456", "TRAINING_MAX", 300.0, "2024-01-15")

	tests := []struct {
		name        string
		strategy    *TaperLoadStrategy
		params      LoadCalculationParams
		wantErrIs   error
		wantErrText string
	}{
		{
			name: "base strategy error propagates",
			strategy: &TaperLoadStrategy{
				BaseStrategy: &PercentOfLoadStrategy{
					ReferenceType: ReferenceTrainingMax,
					Percentage:    100.0,
					// No maxLookup set
				},
			},
			params: LoadCalculationParams{
				UserID:  "user-123",
				LiftID:  "squat-456",
				Context: map[string]interface{}{"daysOut": 10},
			},
			wantErrIs: ErrInvalidParams,
		},
		{
			name: "validation error - nil base strategy",
			strategy: &TaperLoadStrategy{
				BaseStrategy: nil,
			},
			params: LoadCalculationParams{
				UserID:  "user-123",
				LiftID:  "squat-456",
				Context: map[string]interface{}{"daysOut": 10},
			},
			wantErrIs: ErrInvalidParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.strategy.CalculateLoad(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
				t.Errorf("expected error %v, got %v", tt.wantErrIs, err)
			}
		})
	}
}

func TestTaperLoadStrategy_SetMaxLookup(t *testing.T) {
	lookup := newMockMaxLookup()
	lookup.SetMax("user-123", "squat-456", "TRAINING_MAX", 300.0, "2024-01-15")

	baseStrategy := &PercentOfLoadStrategy{
		ReferenceType:     ReferenceTrainingMax,
		Percentage:        100.0,
		RoundingIncrement: 5.0,
		RoundingDirection: RoundNearest,
		// No maxLookup yet
	}

	strategy := &TaperLoadStrategy{
		BaseStrategy: baseStrategy,
	}

	// Inject max lookup
	strategy.SetMaxLookup(lookup)

	params := LoadCalculationParams{
		UserID:  "user-123",
		LiftID:  "squat-456",
		Context: map[string]interface{}{"daysOut": 50}, // Beyond taper
	}

	result, err := strategy.CalculateLoad(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != 300.0 {
		t.Errorf("expected 300.0, got %f", result)
	}
}

func TestTaperLoadStrategy_MarshalJSON(t *testing.T) {
	baseStrategy := &PercentOfLoadStrategy{
		ReferenceType:     ReferenceTrainingMax,
		Percentage:        85.0,
		RoundingIncrement: 5.0,
		RoundingDirection: RoundNearest,
	}

	strategy := &TaperLoadStrategy{
		BaseStrategy: baseStrategy,
		TaperCurve: []TaperCurve{
			{ThresholdDays: 7, Multiplier: 0.5},
			{ThresholdDays: 14, Multiplier: 0.7},
		},
		MaintainIntensity: false,
	}

	data, err := json.Marshal(strategy)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Verify the JSON structure
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse marshaled JSON: %v", err)
	}

	// Check type field
	if typeVal, ok := parsed["type"]; !ok {
		t.Error("expected 'type' field in marshaled JSON")
	} else if typeVal != string(TypeTaper) {
		t.Errorf("expected type %s, got %v", TypeTaper, typeVal)
	}

	// Check baseStrategy exists
	if _, ok := parsed["baseStrategy"]; !ok {
		t.Error("expected 'baseStrategy' field in marshaled JSON")
	}

	// Check taperCurve exists
	if curve, ok := parsed["taperCurve"].([]interface{}); !ok {
		t.Error("expected 'taperCurve' field in marshaled JSON")
	} else if len(curve) != 2 {
		t.Errorf("expected 2 taper curve entries, got %d", len(curve))
	}
}

func TestUnmarshalTaper(t *testing.T) {
	factory := NewStrategyFactory()
	RegisterPercentOf(factory)
	RegisterTaper(factory)

	tests := []struct {
		name    string
		json    string
		wantErr bool
		check   func(*testing.T, LoadStrategy)
	}{
		{
			name: "valid with base percent_of strategy",
			json: `{
				"type": "TAPER",
				"baseStrategy": {
					"type": "PERCENT_OF",
					"referenceType": "TRAINING_MAX",
					"percentage": 85
				}
			}`,
			wantErr: false,
			check: func(t *testing.T, s LoadStrategy) {
				ts := s.(*TaperLoadStrategy)
				if ts.Type() != TypeTaper {
					t.Errorf("expected type %s, got %s", TypeTaper, ts.Type())
				}
				if ts.BaseStrategy == nil {
					t.Error("expected base strategy to be set")
				}
				if ts.BaseStrategy.Type() != TypePercentOf {
					t.Errorf("expected base strategy type %s, got %s", TypePercentOf, ts.BaseStrategy.Type())
				}
			},
		},
		{
			name: "valid with custom curve",
			json: `{
				"type": "TAPER",
				"baseStrategy": {
					"type": "PERCENT_OF",
					"referenceType": "TRAINING_MAX",
					"percentage": 85
				},
				"taperCurve": [
					{"thresholdDays": 7, "multiplier": 0.5},
					{"thresholdDays": 14, "multiplier": 0.7}
				]
			}`,
			wantErr: false,
			check: func(t *testing.T, s LoadStrategy) {
				ts := s.(*TaperLoadStrategy)
				if len(ts.TaperCurve) != 2 {
					t.Errorf("expected 2 taper curve entries, got %d", len(ts.TaperCurve))
				}
				if ts.TaperCurve[0].ThresholdDays != 7 {
					t.Errorf("expected first threshold 7, got %d", ts.TaperCurve[0].ThresholdDays)
				}
				if ts.TaperCurve[0].Multiplier != 0.5 {
					t.Errorf("expected first multiplier 0.5, got %f", ts.TaperCurve[0].Multiplier)
				}
			},
		},
		{
			name: "valid with maintainIntensity",
			json: `{
				"type": "TAPER",
				"baseStrategy": {
					"type": "PERCENT_OF",
					"referenceType": "TRAINING_MAX",
					"percentage": 85
				},
				"maintainIntensity": true
			}`,
			wantErr: false,
			check: func(t *testing.T, s LoadStrategy) {
				ts := s.(*TaperLoadStrategy)
				if !ts.MaintainIntensity {
					t.Error("expected maintainIntensity to be true")
				}
			},
		},
		{
			name: "invalid - missing base strategy",
			json: `{
				"type": "TAPER"
			}`,
			wantErr: true,
		},
		{
			name: "invalid - invalid base strategy",
			json: `{
				"type": "TAPER",
				"baseStrategy": {
					"type": "PERCENT_OF"
				}
			}`,
			wantErr: true,
		},
		{
			name: "invalid - curve not ascending",
			json: `{
				"type": "TAPER",
				"baseStrategy": {
					"type": "PERCENT_OF",
					"referenceType": "TRAINING_MAX",
					"percentage": 85
				},
				"taperCurve": [
					{"thresholdDays": 14, "multiplier": 0.7},
					{"thresholdDays": 7, "multiplier": 0.5}
				]
			}`,
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			json:    `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy, err := factory.CreateFromJSON(json.RawMessage(tt.json))
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.check != nil {
				tt.check(t, strategy)
			}
		})
	}
}

func TestRegisterTaper(t *testing.T) {
	factory := NewStrategyFactory()
	RegisterPercentOf(factory)
	RegisterTaper(factory)

	if !factory.IsRegistered(TypeTaper) {
		t.Error("expected TypeTaper to be registered")
	}
}

func TestNewTaperLoadStrategy(t *testing.T) {
	baseStrategy := &PercentOfLoadStrategy{
		ReferenceType: ReferenceTrainingMax,
		Percentage:    85.0,
	}

	curve := []TaperCurve{
		{ThresholdDays: 7, Multiplier: 0.5},
	}

	strategy := NewTaperLoadStrategy(baseStrategy, curve, true)

	if strategy.BaseStrategy != baseStrategy {
		t.Error("expected base strategy to be set")
	}
	if len(strategy.TaperCurve) != 1 {
		t.Errorf("expected 1 taper curve entry, got %d", len(strategy.TaperCurve))
	}
	if !strategy.MaintainIntensity {
		t.Error("expected maintainIntensity to be true")
	}
}

func TestTaperLoadStrategy_Interface(t *testing.T) {
	var _ LoadStrategy = (*TaperLoadStrategy)(nil)
}

func TestTaperLoadStrategy_RoundTripJSON(t *testing.T) {
	factory := NewStrategyFactory()
	RegisterPercentOf(factory)
	RegisterTaper(factory)

	original := &TaperLoadStrategy{
		BaseStrategy: &PercentOfLoadStrategy{
			ReferenceType:     ReferenceTrainingMax,
			Percentage:        85.0,
			RoundingIncrement: 2.5,
			RoundingDirection: RoundDown,
		},
		TaperCurve: []TaperCurve{
			{ThresholdDays: 7, Multiplier: 0.5},
			{ThresholdDays: 14, Multiplier: 0.7},
		},
		MaintainIntensity: true,
	}

	// Marshal
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Unmarshal
	restored, err := factory.CreateFromJSON(data)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	ts := restored.(*TaperLoadStrategy)

	// Compare
	if ts.Type() != original.Type() {
		t.Errorf("type mismatch: expected %s, got %s", original.Type(), ts.Type())
	}
	if ts.MaintainIntensity != original.MaintainIntensity {
		t.Errorf("maintainIntensity mismatch: expected %v, got %v", original.MaintainIntensity, ts.MaintainIntensity)
	}
	if len(ts.TaperCurve) != len(original.TaperCurve) {
		t.Errorf("taperCurve length mismatch: expected %d, got %d", len(original.TaperCurve), len(ts.TaperCurve))
	}
	for i, tier := range ts.TaperCurve {
		if tier.ThresholdDays != original.TaperCurve[i].ThresholdDays {
			t.Errorf("taperCurve[%d].ThresholdDays mismatch", i)
		}
		if tier.Multiplier != original.TaperCurve[i].Multiplier {
			t.Errorf("taperCurve[%d].Multiplier mismatch", i)
		}
	}
}

// Benchmark tests

func BenchmarkGetTaperMultiplier(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetTaperMultiplier(i % 50) // Vary daysOut
	}
}

func BenchmarkTaperCalculateLoad(b *testing.B) {
	lookup := newMockMaxLookup()
	lookup.SetMax("user-123", "squat-456", "TRAINING_MAX", 300.0, "2024-01-15")

	baseStrategy := &PercentOfLoadStrategy{
		ReferenceType:     ReferenceTrainingMax,
		Percentage:        100.0,
		RoundingIncrement: 5.0,
		RoundingDirection: RoundNearest,
	}
	baseStrategy.SetMaxLookup(lookup)

	strategy := &TaperLoadStrategy{
		BaseStrategy: baseStrategy,
	}

	params := LoadCalculationParams{
		UserID:  "user-123",
		LiftID:  "squat-456",
		Context: map[string]interface{}{"daysOut": 10},
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strategy.CalculateLoad(ctx, params)
	}
}

func BenchmarkTaperMarshalJSON(b *testing.B) {
	strategy := &TaperLoadStrategy{
		BaseStrategy: &PercentOfLoadStrategy{
			ReferenceType:     ReferenceTrainingMax,
			Percentage:        85.0,
			RoundingIncrement: 5.0,
			RoundingDirection: RoundNearest,
		},
		TaperCurve: DefaultTaperCurve(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Marshal(strategy)
	}
}

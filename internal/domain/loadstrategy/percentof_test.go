package loadstrategy

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"testing"
)

// mockMaxLookup implements MaxLookup for testing.
type mockMaxLookup struct {
	maxes map[string]*MaxValue // key: "userID:liftID:maxType"
	err   error
}

func newMockMaxLookup() *mockMaxLookup {
	return &mockMaxLookup{
		maxes: make(map[string]*MaxValue),
	}
}

func (m *mockMaxLookup) SetMax(userID, liftID, maxType string, value float64, date string) {
	key := userID + ":" + liftID + ":" + maxType
	m.maxes[key] = &MaxValue{
		Value:         value,
		EffectiveDate: date,
	}
}

func (m *mockMaxLookup) SetError(err error) {
	m.err = err
}

func (m *mockMaxLookup) GetCurrentMax(ctx context.Context, userID, liftID, maxType string) (*MaxValue, error) {
	if m.err != nil {
		return nil, m.err
	}
	key := userID + ":" + liftID + ":" + maxType
	return m.maxes[key], nil
}

func TestNewPercentOfLoadStrategy(t *testing.T) {
	lookup := newMockMaxLookup()
	strategy := NewPercentOfLoadStrategy(
		ReferenceTrainingMax,
		85.0,
		5.0,
		RoundNearest,
		lookup,
	)

	if strategy.ReferenceType != ReferenceTrainingMax {
		t.Errorf("expected reference type %s, got %s", ReferenceTrainingMax, strategy.ReferenceType)
	}
	if strategy.Percentage != 85.0 {
		t.Errorf("expected percentage 85.0, got %f", strategy.Percentage)
	}
	if strategy.RoundingIncrement != 5.0 {
		t.Errorf("expected rounding increment 5.0, got %f", strategy.RoundingIncrement)
	}
	if strategy.RoundingDirection != RoundNearest {
		t.Errorf("expected rounding direction %s, got %s", RoundNearest, strategy.RoundingDirection)
	}
}

func TestPercentOfLoadStrategy_Type(t *testing.T) {
	strategy := &PercentOfLoadStrategy{}
	if strategy.Type() != TypePercentOf {
		t.Errorf("expected type %s, got %s", TypePercentOf, strategy.Type())
	}
}

func TestPercentOfLoadStrategy_Validate(t *testing.T) {
	tests := []struct {
		name    string
		strategy PercentOfLoadStrategy
		wantErr error
	}{
		{
			name: "valid with all fields",
			strategy: PercentOfLoadStrategy{
				ReferenceType:     ReferenceTrainingMax,
				Percentage:        85.0,
				RoundingIncrement: 5.0,
				RoundingDirection: RoundNearest,
			},
			wantErr: nil,
		},
		{
			name: "valid with ONE_RM",
			strategy: PercentOfLoadStrategy{
				ReferenceType: ReferenceOneRM,
				Percentage:    90.0,
			},
			wantErr: nil,
		},
		{
			name: "valid with percentage over 100 (overload)",
			strategy: PercentOfLoadStrategy{
				ReferenceType: ReferenceTrainingMax,
				Percentage:    105.0,
			},
			wantErr: nil,
		},
		{
			name: "valid with default rounding",
			strategy: PercentOfLoadStrategy{
				ReferenceType: ReferenceTrainingMax,
				Percentage:    85.0,
				// No rounding fields - will use defaults
			},
			wantErr: nil,
		},
		{
			name: "missing reference type",
			strategy: PercentOfLoadStrategy{
				Percentage: 85.0,
			},
			wantErr: ErrReferenceTypeRequired,
		},
		{
			name: "invalid reference type",
			strategy: PercentOfLoadStrategy{
				ReferenceType: "INVALID",
				Percentage:    85.0,
			},
			wantErr: ErrReferenceTypeInvalid,
		},
		{
			name: "zero percentage",
			strategy: PercentOfLoadStrategy{
				ReferenceType: ReferenceTrainingMax,
				Percentage:    0,
			},
			wantErr: ErrPercentageNotPositive,
		},
		{
			name: "negative percentage",
			strategy: PercentOfLoadStrategy{
				ReferenceType: ReferenceTrainingMax,
				Percentage:    -85.0,
			},
			wantErr: ErrPercentageNotPositive,
		},
		{
			name: "invalid rounding direction",
			strategy: PercentOfLoadStrategy{
				ReferenceType:     ReferenceTrainingMax,
				Percentage:        85.0,
				RoundingDirection: "INVALID",
			},
			wantErr: ErrInvalidRoundingDirection,
		},
		{
			name: "negative rounding increment",
			strategy: PercentOfLoadStrategy{
				ReferenceType:     ReferenceTrainingMax,
				Percentage:        85.0,
				RoundingIncrement: -5.0,
			},
			wantErr: ErrInvalidParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.strategy.Validate()
			if tt.wantErr != nil {
				if err == nil {
					t.Error("expected error, got nil")
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestPercentOfLoadStrategy_CalculateLoad(t *testing.T) {
	tests := []struct {
		name         string
		strategy     PercentOfLoadStrategy
		params       LoadCalculationParams
		setupMaxes   func(*mockMaxLookup)
		expected     float64
		wantErr      error
		wantErrMsg   string
	}{
		{
			name: "85% of TM (315) = 270 (rounded from 267.75)",
			strategy: PercentOfLoadStrategy{
				ReferenceType:     ReferenceTrainingMax,
				Percentage:        85.0,
				RoundingIncrement: 5.0,
				RoundingDirection: RoundNearest,
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat-456",
			},
			setupMaxes: func(m *mockMaxLookup) {
				m.SetMax("user-123", "squat-456", "TRAINING_MAX", 315.0, "2024-01-15")
			},
			expected: 270.0,
		},
		{
			name: "90% of 1RM (400) = 360",
			strategy: PercentOfLoadStrategy{
				ReferenceType:     ReferenceOneRM,
				Percentage:        90.0,
				RoundingIncrement: 5.0,
				RoundingDirection: RoundNearest,
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "deadlift-789",
			},
			setupMaxes: func(m *mockMaxLookup) {
				m.SetMax("user-123", "deadlift-789", "ONE_RM", 400.0, "2024-01-15")
			},
			expected: 360.0,
		},
		{
			name: "105% of TM (overload) = 330.75 -> 330",
			strategy: PercentOfLoadStrategy{
				ReferenceType:     ReferenceTrainingMax,
				Percentage:        105.0,
				RoundingIncrement: 5.0,
				RoundingDirection: RoundNearest,
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat-456",
			},
			setupMaxes: func(m *mockMaxLookup) {
				m.SetMax("user-123", "squat-456", "TRAINING_MAX", 315.0, "2024-01-15")
			},
			expected: 330.0,
		},
		{
			name: "rounding down (conservative)",
			strategy: PercentOfLoadStrategy{
				ReferenceType:     ReferenceTrainingMax,
				Percentage:        85.0,
				RoundingIncrement: 5.0,
				RoundingDirection: RoundDown,
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat-456",
			},
			setupMaxes: func(m *mockMaxLookup) {
				m.SetMax("user-123", "squat-456", "TRAINING_MAX", 315.0, "2024-01-15")
			},
			expected: 265.0, // 267.75 rounded down
		},
		{
			name: "rounding up",
			strategy: PercentOfLoadStrategy{
				ReferenceType:     ReferenceTrainingMax,
				Percentage:        85.0,
				RoundingIncrement: 5.0,
				RoundingDirection: RoundUp,
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat-456",
			},
			setupMaxes: func(m *mockMaxLookup) {
				m.SetMax("user-123", "squat-456", "TRAINING_MAX", 315.0, "2024-01-15")
			},
			expected: 270.0, // 267.75 rounded up
		},
		{
			name: "2.5 lb increment",
			strategy: PercentOfLoadStrategy{
				ReferenceType:     ReferenceTrainingMax,
				Percentage:        75.0,
				RoundingIncrement: 2.5,
				RoundingDirection: RoundNearest,
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "bench-123",
			},
			setupMaxes: func(m *mockMaxLookup) {
				m.SetMax("user-123", "bench-123", "TRAINING_MAX", 200.0, "2024-01-15")
			},
			expected: 150.0, // 200 * 0.75 = 150, exactly on increment
		},
		{
			name: "default rounding when not specified",
			strategy: PercentOfLoadStrategy{
				ReferenceType: ReferenceTrainingMax,
				Percentage:    85.0,
				// No rounding fields - should use defaults
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat-456",
			},
			setupMaxes: func(m *mockMaxLookup) {
				m.SetMax("user-123", "squat-456", "TRAINING_MAX", 315.0, "2024-01-15")
			},
			expected: 270.0, // Uses default 5.0 increment and NEAREST direction
		},
		{
			name: "max not found",
			strategy: PercentOfLoadStrategy{
				ReferenceType: ReferenceTrainingMax,
				Percentage:    85.0,
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat-456",
			},
			setupMaxes: func(m *mockMaxLookup) {
				// No max set
			},
			wantErr: ErrMaxNotFound,
		},
		{
			name: "missing user ID in params",
			strategy: PercentOfLoadStrategy{
				ReferenceType: ReferenceTrainingMax,
				Percentage:    85.0,
			},
			params: LoadCalculationParams{
				UserID: "",
				LiftID: "squat-456",
			},
			wantErr: ErrInvalidParams,
		},
		{
			name: "missing lift ID in params",
			strategy: PercentOfLoadStrategy{
				ReferenceType: ReferenceTrainingMax,
				Percentage:    85.0,
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "",
			},
			wantErr: ErrInvalidParams,
		},
		{
			name: "repository error",
			strategy: PercentOfLoadStrategy{
				ReferenceType: ReferenceTrainingMax,
				Percentage:    85.0,
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat-456",
			},
			setupMaxes: func(m *mockMaxLookup) {
				m.SetError(errors.New("database error"))
			},
			wantErrMsg: "failed to lookup max",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			lookup := newMockMaxLookup()
			if tt.setupMaxes != nil {
				tt.setupMaxes(lookup)
			}
			tt.strategy.SetMaxLookup(lookup)

			result, err := tt.strategy.CalculateLoad(ctx, tt.params)

			if tt.wantErr != nil {
				if err == nil {
					t.Error("expected error, got nil")
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if tt.wantErrMsg != "" {
				if err == nil {
					t.Error("expected error, got nil")
					return
				}
				if err.Error()[:len(tt.wantErrMsg)] != tt.wantErrMsg {
					t.Errorf("expected error message to start with %q, got %q", tt.wantErrMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestPercentOfLoadStrategy_CalculateLoad_NoMaxLookup(t *testing.T) {
	strategy := &PercentOfLoadStrategy{
		ReferenceType: ReferenceTrainingMax,
		Percentage:    85.0,
		// No maxLookup set
	}

	params := LoadCalculationParams{
		UserID: "user-123",
		LiftID: "squat-456",
	}

	_, err := strategy.CalculateLoad(context.Background(), params)
	if err == nil {
		t.Error("expected error when maxLookup is not set")
	}
	if !errors.Is(err, ErrInvalidParams) {
		t.Errorf("expected ErrInvalidParams, got %v", err)
	}
}

func TestPercentOfLoadStrategy_SetMaxLookup(t *testing.T) {
	strategy := &PercentOfLoadStrategy{
		ReferenceType: ReferenceTrainingMax,
		Percentage:    85.0,
	}

	lookup := newMockMaxLookup()
	lookup.SetMax("user-123", "squat-456", "TRAINING_MAX", 300.0, "2024-01-15")

	strategy.SetMaxLookup(lookup)

	params := LoadCalculationParams{
		UserID: "user-123",
		LiftID: "squat-456",
	}

	result, err := strategy.CalculateLoad(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 300 * 0.85 = 255, exactly on 5lb increment
	if result != 255.0 {
		t.Errorf("expected 255.0, got %f", result)
	}
}

func TestPercentOfLoadStrategy_MarshalJSON(t *testing.T) {
	strategy := &PercentOfLoadStrategy{
		ReferenceType:     ReferenceTrainingMax,
		Percentage:        85.0,
		RoundingIncrement: 5.0,
		RoundingDirection: RoundNearest,
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
	} else if typeVal != string(TypePercentOf) {
		t.Errorf("expected type %s, got %v", TypePercentOf, typeVal)
	}

	// Check referenceType field
	if refType, ok := parsed["referenceType"]; !ok {
		t.Error("expected 'referenceType' field in marshaled JSON")
	} else if refType != string(ReferenceTrainingMax) {
		t.Errorf("expected referenceType %s, got %v", ReferenceTrainingMax, refType)
	}

	// Check percentage field
	if pct, ok := parsed["percentage"]; !ok {
		t.Error("expected 'percentage' field in marshaled JSON")
	} else if pct != 85.0 {
		t.Errorf("expected percentage 85.0, got %v", pct)
	}

	// Check rounding fields
	if inc, ok := parsed["roundingIncrement"]; !ok {
		t.Error("expected 'roundingIncrement' field in marshaled JSON")
	} else if inc != 5.0 {
		t.Errorf("expected roundingIncrement 5.0, got %v", inc)
	}

	if dir, ok := parsed["roundingDirection"]; !ok {
		t.Error("expected 'roundingDirection' field in marshaled JSON")
	} else if dir != string(RoundNearest) {
		t.Errorf("expected roundingDirection %s, got %v", RoundNearest, dir)
	}
}

func TestPercentOfLoadStrategy_MarshalJSON_OmitEmpty(t *testing.T) {
	strategy := &PercentOfLoadStrategy{
		ReferenceType: ReferenceTrainingMax,
		Percentage:    85.0,
		// RoundingIncrement and RoundingDirection not set
	}

	data, err := json.Marshal(strategy)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse marshaled JSON: %v", err)
	}

	// roundingIncrement should be omitted when zero
	if _, ok := parsed["roundingIncrement"]; ok {
		t.Error("expected 'roundingIncrement' to be omitted when zero")
	}

	// roundingDirection should be omitted when empty
	if _, ok := parsed["roundingDirection"]; ok {
		t.Error("expected 'roundingDirection' to be omitted when empty")
	}
}

func TestUnmarshalPercentOf(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
		check   func(*testing.T, LoadStrategy)
	}{
		{
			name: "valid with all fields",
			json: `{
				"type": "PERCENT_OF",
				"referenceType": "TRAINING_MAX",
				"percentage": 85,
				"roundingIncrement": 5.0,
				"roundingDirection": "NEAREST"
			}`,
			wantErr: false,
			check: func(t *testing.T, s LoadStrategy) {
				ps := s.(*PercentOfLoadStrategy)
				if ps.ReferenceType != ReferenceTrainingMax {
					t.Errorf("expected referenceType %s, got %s", ReferenceTrainingMax, ps.ReferenceType)
				}
				if ps.Percentage != 85.0 {
					t.Errorf("expected percentage 85.0, got %f", ps.Percentage)
				}
				if ps.RoundingIncrement != 5.0 {
					t.Errorf("expected roundingIncrement 5.0, got %f", ps.RoundingIncrement)
				}
				if ps.RoundingDirection != RoundNearest {
					t.Errorf("expected roundingDirection %s, got %s", RoundNearest, ps.RoundingDirection)
				}
			},
		},
		{
			name: "valid with ONE_RM",
			json: `{
				"type": "PERCENT_OF",
				"referenceType": "ONE_RM",
				"percentage": 90
			}`,
			wantErr: false,
			check: func(t *testing.T, s LoadStrategy) {
				ps := s.(*PercentOfLoadStrategy)
				if ps.ReferenceType != ReferenceOneRM {
					t.Errorf("expected referenceType %s, got %s", ReferenceOneRM, ps.ReferenceType)
				}
			},
		},
		{
			name: "valid with rounding DOWN",
			json: `{
				"type": "PERCENT_OF",
				"referenceType": "TRAINING_MAX",
				"percentage": 85,
				"roundingDirection": "DOWN"
			}`,
			wantErr: false,
			check: func(t *testing.T, s LoadStrategy) {
				ps := s.(*PercentOfLoadStrategy)
				if ps.RoundingDirection != RoundDown {
					t.Errorf("expected roundingDirection %s, got %s", RoundDown, ps.RoundingDirection)
				}
			},
		},
		{
			name: "missing reference type",
			json: `{
				"type": "PERCENT_OF",
				"percentage": 85
			}`,
			wantErr: true,
		},
		{
			name: "invalid reference type",
			json: `{
				"type": "PERCENT_OF",
				"referenceType": "INVALID",
				"percentage": 85
			}`,
			wantErr: true,
		},
		{
			name: "zero percentage",
			json: `{
				"type": "PERCENT_OF",
				"referenceType": "TRAINING_MAX",
				"percentage": 0
			}`,
			wantErr: true,
		},
		{
			name: "negative percentage",
			json: `{
				"type": "PERCENT_OF",
				"referenceType": "TRAINING_MAX",
				"percentage": -85
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
			strategy, err := UnmarshalPercentOf(json.RawMessage(tt.json))
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

func TestRegisterPercentOf(t *testing.T) {
	factory := NewStrategyFactory()
	RegisterPercentOf(factory)

	if !factory.IsRegistered(TypePercentOf) {
		t.Error("expected TypePercentOf to be registered")
	}

	// Test that we can create a strategy from JSON
	jsonData := []byte(`{
		"type": "PERCENT_OF",
		"referenceType": "TRAINING_MAX",
		"percentage": 85
	}`)

	strategy, err := factory.CreateFromJSON(jsonData)
	if err != nil {
		t.Fatalf("failed to create strategy: %v", err)
	}

	if strategy.Type() != TypePercentOf {
		t.Errorf("expected type %s, got %s", TypePercentOf, strategy.Type())
	}
}

func TestReferenceTypeConstants(t *testing.T) {
	// Verify the constant values match what's expected
	if string(ReferenceOneRM) != "ONE_RM" {
		t.Errorf("expected ReferenceOneRM to be 'ONE_RM', got %s", ReferenceOneRM)
	}
	if string(ReferenceTrainingMax) != "TRAINING_MAX" {
		t.Errorf("expected ReferenceTrainingMax to be 'TRAINING_MAX', got %s", ReferenceTrainingMax)
	}

	// Verify aliases
	if OneRM != ReferenceOneRM {
		t.Error("OneRM alias should equal ReferenceOneRM")
	}
	if TrainingMax != ReferenceTrainingMax {
		t.Error("TrainingMax alias should equal ReferenceTrainingMax")
	}
}

func TestValidReferenceTypes(t *testing.T) {
	expectedTypes := []ReferenceType{
		ReferenceOneRM,
		ReferenceTrainingMax,
	}

	for _, refType := range expectedTypes {
		if !ValidReferenceTypes[refType] {
			t.Errorf("expected %s to be in ValidReferenceTypes", refType)
		}
	}

	if len(ValidReferenceTypes) != len(expectedTypes) {
		t.Errorf("expected %d types in ValidReferenceTypes, got %d",
			len(expectedTypes), len(ValidReferenceTypes))
	}
}

func TestPercentOfErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		msg  string
	}{
		{"ErrPercentageRequired", ErrPercentageRequired, "percentage is required"},
		{"ErrPercentageNotPositive", ErrPercentageNotPositive, "percentage must be greater than 0"},
		{"ErrReferenceTypeRequired", ErrReferenceTypeRequired, "reference type is required"},
		{"ErrReferenceTypeInvalid", ErrReferenceTypeInvalid, "reference type must be ONE_RM or TRAINING_MAX"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.msg {
				t.Errorf("expected message %q, got %q", tt.msg, tt.err.Error())
			}
		})
	}
}

// TestPercentOfLoadStrategy_Interface ensures the struct implements LoadStrategy
func TestPercentOfLoadStrategy_Interface(t *testing.T) {
	var _ LoadStrategy = (*PercentOfLoadStrategy)(nil)
}

// TestRoundTripJSON tests that marshaling and unmarshaling produces equivalent strategies
func TestRoundTripJSON(t *testing.T) {
	original := &PercentOfLoadStrategy{
		ReferenceType:     ReferenceTrainingMax,
		Percentage:        85.0,
		RoundingIncrement: 2.5,
		RoundingDirection: RoundDown,
	}

	// Marshal
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Unmarshal
	restored, err := UnmarshalPercentOf(data)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	ps := restored.(*PercentOfLoadStrategy)

	// Compare
	if ps.ReferenceType != original.ReferenceType {
		t.Errorf("referenceType mismatch: expected %s, got %s", original.ReferenceType, ps.ReferenceType)
	}
	if ps.Percentage != original.Percentage {
		t.Errorf("percentage mismatch: expected %f, got %f", original.Percentage, ps.Percentage)
	}
	if ps.RoundingIncrement != original.RoundingIncrement {
		t.Errorf("roundingIncrement mismatch: expected %f, got %f", original.RoundingIncrement, ps.RoundingIncrement)
	}
	if ps.RoundingDirection != original.RoundingDirection {
		t.Errorf("roundingDirection mismatch: expected %s, got %s", original.RoundingDirection, ps.RoundingDirection)
	}
}

// TestCalculateLoadWithVariousPercentages tests a range of realistic percentage values
func TestCalculateLoadWithVariousPercentages(t *testing.T) {
	lookup := newMockMaxLookup()
	lookup.SetMax("user-123", "squat-456", "TRAINING_MAX", 315.0, "2024-01-15")

	tests := []struct {
		percentage float64
		expected   float64
	}{
		{50.0, 160.0},  // 315 * 0.50 = 157.5 -> round(31.5)*5 = 32*5 = 160
		{60.0, 190.0},  // 315 * 0.60 = 189 -> round(37.8)*5 = 38*5 = 190
		{65.0, 205.0},  // 315 * 0.65 = 204.75 -> round(40.95)*5 = 41*5 = 205
		{70.0, 220.0},  // 315 * 0.70 = 220.5 -> round(44.1)*5 = 44*5 = 220
		{75.0, 235.0},  // 315 * 0.75 = 236.25 -> round(47.25)*5 = 47*5 = 235
		{80.0, 250.0},  // 315 * 0.80 = 252 -> round(50.4)*5 = 50*5 = 250
		{85.0, 270.0},  // 315 * 0.85 = 267.75 -> round(53.55)*5 = 54*5 = 270
		{90.0, 285.0},  // 315 * 0.90 = 283.5 -> round(56.7)*5 = 57*5 = 285
		{95.0, 300.0},  // 315 * 0.95 = 299.25 -> round(59.85)*5 = 60*5 = 300
		{100.0, 315.0}, // 315 * 1.00 = 315 -> round(63)*5 = 63*5 = 315
	}

	for _, tt := range tests {
		t.Run(
			// Format percentage as string for test name
			"pct_" + func(f float64) string {
				return string(rune('0'+int(f/10))) + string(rune('0'+int(f)%10))
			}(tt.percentage),
			func(t *testing.T) {
				strategy := &PercentOfLoadStrategy{
					ReferenceType:     ReferenceTrainingMax,
					Percentage:        tt.percentage,
					RoundingIncrement: 5.0,
					RoundingDirection: RoundNearest,
				}
				strategy.SetMaxLookup(lookup)

				params := LoadCalculationParams{
					UserID: "user-123",
					LiftID: "squat-456",
				}

				result, err := strategy.CalculateLoad(context.Background(), params)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if math.Abs(result-tt.expected) > 0.0001 {
					t.Errorf("at %.0f%%: expected %f, got %f", tt.percentage, tt.expected, result)
				}
			},
		)
	}
}

// Benchmark tests
func BenchmarkPercentOfCalculateLoad(b *testing.B) {
	lookup := newMockMaxLookup()
	lookup.SetMax("user-123", "squat-456", "TRAINING_MAX", 315.0, "2024-01-15")

	strategy := &PercentOfLoadStrategy{
		ReferenceType:     ReferenceTrainingMax,
		Percentage:        85.0,
		RoundingIncrement: 5.0,
		RoundingDirection: RoundNearest,
	}
	strategy.SetMaxLookup(lookup)

	params := LoadCalculationParams{
		UserID: "user-123",
		LiftID: "squat-456",
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strategy.CalculateLoad(ctx, params)
	}
}

func BenchmarkPercentOfMarshalJSON(b *testing.B) {
	strategy := &PercentOfLoadStrategy{
		ReferenceType:     ReferenceTrainingMax,
		Percentage:        85.0,
		RoundingIncrement: 5.0,
		RoundingDirection: RoundNearest,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Marshal(strategy)
	}
}

func BenchmarkUnmarshalPercentOf(b *testing.B) {
	data := []byte(`{
		"type": "PERCENT_OF",
		"referenceType": "TRAINING_MAX",
		"percentage": 85,
		"roundingIncrement": 5.0,
		"roundingDirection": "NEAREST"
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		UnmarshalPercentOf(data)
	}
}

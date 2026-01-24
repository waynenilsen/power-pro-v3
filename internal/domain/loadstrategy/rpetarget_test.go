package loadstrategy

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"testing"

	"github.com/waynenilsen/power-pro-v3/internal/domain/rpechart"
)

func TestNewRPETargetLoadStrategy(t *testing.T) {
	lookup := newMockMaxLookup()
	chart := rpechart.NewDefaultRPEChart()
	strategy := NewRPETargetLoadStrategy(
		5,
		8.0,
		5.0,
		RoundNearest,
		lookup,
		chart,
	)

	if strategy.TargetReps != 5 {
		t.Errorf("expected target reps 5, got %d", strategy.TargetReps)
	}
	if strategy.TargetRPE != 8.0 {
		t.Errorf("expected target RPE 8.0, got %f", strategy.TargetRPE)
	}
	if strategy.RoundingIncrement != 5.0 {
		t.Errorf("expected rounding increment 5.0, got %f", strategy.RoundingIncrement)
	}
	if strategy.RoundingDirection != RoundNearest {
		t.Errorf("expected rounding direction %s, got %s", RoundNearest, strategy.RoundingDirection)
	}
}

func TestRPETargetLoadStrategy_Type(t *testing.T) {
	strategy := &RPETargetLoadStrategy{}
	if strategy.Type() != TypeRPETarget {
		t.Errorf("expected type %s, got %s", TypeRPETarget, strategy.Type())
	}
}

func TestRPETargetLoadStrategy_Validate(t *testing.T) {
	tests := []struct {
		name     string
		strategy RPETargetLoadStrategy
		wantErr  error
	}{
		{
			name: "valid with all fields",
			strategy: RPETargetLoadStrategy{
				TargetReps:        5,
				TargetRPE:         8.0,
				RoundingIncrement: 5.0,
				RoundingDirection: RoundNearest,
			},
			wantErr: nil,
		},
		{
			name: "valid with minimal fields",
			strategy: RPETargetLoadStrategy{
				TargetReps: 3,
				TargetRPE:  9.0,
			},
			wantErr: nil,
		},
		{
			name: "valid RPE 7.0",
			strategy: RPETargetLoadStrategy{
				TargetReps: 10,
				TargetRPE:  7.0,
			},
			wantErr: nil,
		},
		{
			name: "valid RPE 10.0",
			strategy: RPETargetLoadStrategy{
				TargetReps: 1,
				TargetRPE:  10.0,
			},
			wantErr: nil,
		},
		{
			name: "valid RPE 8.5",
			strategy: RPETargetLoadStrategy{
				TargetReps: 5,
				TargetRPE:  8.5,
			},
			wantErr: nil,
		},
		{
			name: "valid reps 1",
			strategy: RPETargetLoadStrategy{
				TargetReps: 1,
				TargetRPE:  9.0,
			},
			wantErr: nil,
		},
		{
			name: "valid reps 12",
			strategy: RPETargetLoadStrategy{
				TargetReps: 12,
				TargetRPE:  8.0,
			},
			wantErr: nil,
		},
		{
			name: "invalid reps 0",
			strategy: RPETargetLoadStrategy{
				TargetReps: 0,
				TargetRPE:  8.0,
			},
			wantErr: ErrTargetRepsInvalid,
		},
		{
			name: "invalid reps 13",
			strategy: RPETargetLoadStrategy{
				TargetReps: 13,
				TargetRPE:  8.0,
			},
			wantErr: ErrTargetRepsInvalid,
		},
		{
			name: "invalid reps negative",
			strategy: RPETargetLoadStrategy{
				TargetReps: -1,
				TargetRPE:  8.0,
			},
			wantErr: ErrTargetRepsInvalid,
		},
		{
			name: "invalid RPE 6.5",
			strategy: RPETargetLoadStrategy{
				TargetReps: 5,
				TargetRPE:  6.5,
			},
			wantErr: ErrTargetRPEInvalid,
		},
		{
			name: "invalid RPE 10.5",
			strategy: RPETargetLoadStrategy{
				TargetReps: 5,
				TargetRPE:  10.5,
			},
			wantErr: ErrTargetRPEInvalid,
		},
		{
			name: "invalid RPE 8.3 (not 0.5 increment)",
			strategy: RPETargetLoadStrategy{
				TargetReps: 5,
				TargetRPE:  8.3,
			},
			wantErr: ErrTargetRPEInvalid,
		},
		{
			name: "invalid rounding direction",
			strategy: RPETargetLoadStrategy{
				TargetReps:        5,
				TargetRPE:         8.0,
				RoundingDirection: "INVALID",
			},
			wantErr: ErrInvalidRoundingDirection,
		},
		{
			name: "negative rounding increment",
			strategy: RPETargetLoadStrategy{
				TargetReps:        5,
				TargetRPE:         8.0,
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

func TestRPETargetLoadStrategy_CalculateLoad(t *testing.T) {
	tests := []struct {
		name       string
		strategy   RPETargetLoadStrategy
		params     LoadCalculationParams
		setupMaxes func(*mockMaxLookup)
		expected   float64
		wantErr    error
		wantErrMsg string
	}{
		{
			name: "5 reps @ RPE 8 (77%) = 310 lbs (from 400 1RM)",
			strategy: RPETargetLoadStrategy{
				TargetReps:        5,
				TargetRPE:         8.0,
				RoundingIncrement: 5.0,
				RoundingDirection: RoundNearest,
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat-456",
			},
			setupMaxes: func(m *mockMaxLookup) {
				m.SetMax("user-123", "squat-456", "ONE_RM", 400.0, "2024-01-15")
			},
			expected: 310.0, // 400 * 0.77 = 308 -> 310 rounded
		},
		{
			name: "1 rep @ RPE 10 (100%) = 400 lbs",
			strategy: RPETargetLoadStrategy{
				TargetReps:        1,
				TargetRPE:         10.0,
				RoundingIncrement: 5.0,
				RoundingDirection: RoundNearest,
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat-456",
			},
			setupMaxes: func(m *mockMaxLookup) {
				m.SetMax("user-123", "squat-456", "ONE_RM", 400.0, "2024-01-15")
			},
			expected: 400.0, // 400 * 1.00 = 400
		},
		{
			name: "3 reps @ RPE 9 (89%) = 355 lbs (from 400 1RM)",
			strategy: RPETargetLoadStrategy{
				TargetReps:        3,
				TargetRPE:         9.0,
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
			expected: 355.0, // 400 * 0.89 = 356 -> 355 rounded
		},
		{
			name: "rounding down (conservative)",
			strategy: RPETargetLoadStrategy{
				TargetReps:        5,
				TargetRPE:         8.0,
				RoundingIncrement: 5.0,
				RoundingDirection: RoundDown,
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat-456",
			},
			setupMaxes: func(m *mockMaxLookup) {
				m.SetMax("user-123", "squat-456", "ONE_RM", 400.0, "2024-01-15")
			},
			expected: 305.0, // 400 * 0.77 = 308 -> 305 rounded down
		},
		{
			name: "rounding up",
			strategy: RPETargetLoadStrategy{
				TargetReps:        5,
				TargetRPE:         8.0,
				RoundingIncrement: 5.0,
				RoundingDirection: RoundUp,
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat-456",
			},
			setupMaxes: func(m *mockMaxLookup) {
				m.SetMax("user-123", "squat-456", "ONE_RM", 400.0, "2024-01-15")
			},
			expected: 310.0, // 400 * 0.77 = 308 -> 310 rounded up
		},
		{
			name: "2.5 lb increment",
			strategy: RPETargetLoadStrategy{
				TargetReps:        5,
				TargetRPE:         8.0,
				RoundingIncrement: 2.5,
				RoundingDirection: RoundNearest,
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "bench-123",
			},
			setupMaxes: func(m *mockMaxLookup) {
				m.SetMax("user-123", "bench-123", "ONE_RM", 250.0, "2024-01-15")
			},
			expected: 192.5, // 250 * 0.77 = 192.5 exactly
		},
		{
			name: "default rounding when not specified",
			strategy: RPETargetLoadStrategy{
				TargetReps: 5,
				TargetRPE:  8.0,
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat-456",
			},
			setupMaxes: func(m *mockMaxLookup) {
				m.SetMax("user-123", "squat-456", "ONE_RM", 400.0, "2024-01-15")
			},
			expected: 310.0, // Uses default 5.0 increment and NEAREST direction
		},
		{
			name: "max not found",
			strategy: RPETargetLoadStrategy{
				TargetReps: 5,
				TargetRPE:  8.0,
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
			strategy: RPETargetLoadStrategy{
				TargetReps: 5,
				TargetRPE:  8.0,
			},
			params: LoadCalculationParams{
				UserID: "",
				LiftID: "squat-456",
			},
			wantErr: ErrInvalidParams,
		},
		{
			name: "missing lift ID in params",
			strategy: RPETargetLoadStrategy{
				TargetReps: 5,
				TargetRPE:  8.0,
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "",
			},
			wantErr: ErrInvalidParams,
		},
		{
			name: "repository error",
			strategy: RPETargetLoadStrategy{
				TargetReps: 5,
				TargetRPE:  8.0,
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
			tt.strategy.SetRPEChart(rpechart.NewDefaultRPEChart())

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
				if len(err.Error()) < len(tt.wantErrMsg) || err.Error()[:len(tt.wantErrMsg)] != tt.wantErrMsg {
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

func TestRPETargetLoadStrategy_CalculateLoad_NoMaxLookup(t *testing.T) {
	strategy := &RPETargetLoadStrategy{
		TargetReps: 5,
		TargetRPE:  8.0,
	}
	strategy.SetRPEChart(rpechart.NewDefaultRPEChart())
	// No maxLookup set

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

func TestRPETargetLoadStrategy_CalculateLoad_NoRPEChart(t *testing.T) {
	lookup := newMockMaxLookup()
	lookup.SetMax("user-123", "squat-456", "ONE_RM", 400.0, "2024-01-15")

	strategy := &RPETargetLoadStrategy{
		TargetReps: 5,
		TargetRPE:  8.0,
	}
	strategy.SetMaxLookup(lookup)
	// No rpeChart set

	params := LoadCalculationParams{
		UserID: "user-123",
		LiftID: "squat-456",
	}

	_, err := strategy.CalculateLoad(context.Background(), params)
	if err == nil {
		t.Error("expected error when RPE chart is not set")
	}
	if !errors.Is(err, ErrRPEChartRequired) {
		t.Errorf("expected ErrRPEChartRequired, got %v", err)
	}
}

func TestRPETargetLoadStrategy_CalculateLoad_UsesLookupContextChart(t *testing.T) {
	lookup := newMockMaxLookup()
	lookup.SetMax("user-123", "squat-456", "ONE_RM", 400.0, "2024-01-15")

	strategy := &RPETargetLoadStrategy{
		TargetReps:        5,
		TargetRPE:         8.0,
		RoundingIncrement: 5.0,
		RoundingDirection: RoundNearest,
	}
	strategy.SetMaxLookup(lookup)
	// No injected chart, will use LookupContext chart

	params := LoadCalculationParams{
		UserID: "user-123",
		LiftID: "squat-456",
		LookupContext: &LookupContext{
			RPEChart: rpechart.NewDefaultRPEChart(),
		},
	}

	result, err := strategy.CalculateLoad(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 400 * 0.77 = 308 -> 310 rounded
	expected := 310.0
	if math.Abs(result-expected) > 0.0001 {
		t.Errorf("expected %f, got %f", expected, result)
	}
}

func TestRPETargetLoadStrategy_SetDependencies(t *testing.T) {
	strategy := &RPETargetLoadStrategy{
		TargetReps: 5,
		TargetRPE:  8.0,
	}

	lookup := newMockMaxLookup()
	lookup.SetMax("user-123", "squat-456", "ONE_RM", 300.0, "2024-01-15")
	chart := rpechart.NewDefaultRPEChart()

	strategy.SetMaxLookup(lookup)
	strategy.SetRPEChart(chart)

	params := LoadCalculationParams{
		UserID: "user-123",
		LiftID: "squat-456",
	}

	result, err := strategy.CalculateLoad(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 300 * 0.77 = 231 -> 230 rounded
	if result != 230.0 {
		t.Errorf("expected 230.0, got %f", result)
	}
}

func TestRPETargetLoadStrategy_MarshalJSON(t *testing.T) {
	strategy := &RPETargetLoadStrategy{
		TargetReps:        5,
		TargetRPE:         8.0,
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
	} else if typeVal != string(TypeRPETarget) {
		t.Errorf("expected type %s, got %v", TypeRPETarget, typeVal)
	}

	// Check targetReps field
	if reps, ok := parsed["targetReps"]; !ok {
		t.Error("expected 'targetReps' field in marshaled JSON")
	} else if reps != 5.0 { // JSON numbers are float64
		t.Errorf("expected targetReps 5, got %v", reps)
	}

	// Check targetRpe field
	if rpe, ok := parsed["targetRpe"]; !ok {
		t.Error("expected 'targetRpe' field in marshaled JSON")
	} else if rpe != 8.0 {
		t.Errorf("expected targetRpe 8.0, got %v", rpe)
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

func TestRPETargetLoadStrategy_MarshalJSON_OmitEmpty(t *testing.T) {
	strategy := &RPETargetLoadStrategy{
		TargetReps: 5,
		TargetRPE:  8.0,
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

func TestUnmarshalRPETarget(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
		check   func(*testing.T, LoadStrategy)
	}{
		{
			name: "valid with all fields",
			json: `{
				"type": "RPE_TARGET",
				"targetReps": 5,
				"targetRpe": 8.0,
				"roundingIncrement": 5.0,
				"roundingDirection": "NEAREST"
			}`,
			wantErr: false,
			check: func(t *testing.T, s LoadStrategy) {
				rs := s.(*RPETargetLoadStrategy)
				if rs.TargetReps != 5 {
					t.Errorf("expected targetReps 5, got %d", rs.TargetReps)
				}
				if rs.TargetRPE != 8.0 {
					t.Errorf("expected targetRpe 8.0, got %f", rs.TargetRPE)
				}
				if rs.RoundingIncrement != 5.0 {
					t.Errorf("expected roundingIncrement 5.0, got %f", rs.RoundingIncrement)
				}
				if rs.RoundingDirection != RoundNearest {
					t.Errorf("expected roundingDirection %s, got %s", RoundNearest, rs.RoundingDirection)
				}
			},
		},
		{
			name: "valid with minimal fields",
			json: `{
				"type": "RPE_TARGET",
				"targetReps": 3,
				"targetRpe": 9.0
			}`,
			wantErr: false,
			check: func(t *testing.T, s LoadStrategy) {
				rs := s.(*RPETargetLoadStrategy)
				if rs.TargetReps != 3 {
					t.Errorf("expected targetReps 3, got %d", rs.TargetReps)
				}
				if rs.TargetRPE != 9.0 {
					t.Errorf("expected targetRpe 9.0, got %f", rs.TargetRPE)
				}
			},
		},
		{
			name: "valid with rounding DOWN",
			json: `{
				"type": "RPE_TARGET",
				"targetReps": 5,
				"targetRpe": 8.0,
				"roundingDirection": "DOWN"
			}`,
			wantErr: false,
			check: func(t *testing.T, s LoadStrategy) {
				rs := s.(*RPETargetLoadStrategy)
				if rs.RoundingDirection != RoundDown {
					t.Errorf("expected roundingDirection %s, got %s", RoundDown, rs.RoundingDirection)
				}
			},
		},
		{
			name: "invalid reps 0",
			json: `{
				"type": "RPE_TARGET",
				"targetReps": 0,
				"targetRpe": 8.0
			}`,
			wantErr: true,
		},
		{
			name: "invalid reps 13",
			json: `{
				"type": "RPE_TARGET",
				"targetReps": 13,
				"targetRpe": 8.0
			}`,
			wantErr: true,
		},
		{
			name: "invalid RPE 6.5",
			json: `{
				"type": "RPE_TARGET",
				"targetReps": 5,
				"targetRpe": 6.5
			}`,
			wantErr: true,
		},
		{
			name: "invalid RPE 10.5",
			json: `{
				"type": "RPE_TARGET",
				"targetReps": 5,
				"targetRpe": 10.5
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
			strategy, err := UnmarshalRPETarget(json.RawMessage(tt.json))
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

func TestRegisterRPETarget(t *testing.T) {
	factory := NewStrategyFactory()
	RegisterRPETarget(factory)

	if !factory.IsRegistered(TypeRPETarget) {
		t.Error("expected TypeRPETarget to be registered")
	}

	// Test that we can create a strategy from JSON
	jsonData := []byte(`{
		"type": "RPE_TARGET",
		"targetReps": 5,
		"targetRpe": 8.0
	}`)

	strategy, err := factory.CreateFromJSON(jsonData)
	if err != nil {
		t.Fatalf("failed to create strategy: %v", err)
	}

	if strategy.Type() != TypeRPETarget {
		t.Errorf("expected type %s, got %s", TypeRPETarget, strategy.Type())
	}
}

func TestRPETargetErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		msg  string
	}{
		{"ErrTargetRepsInvalid", ErrTargetRepsInvalid, "target reps must be between 1 and 12"},
		{"ErrTargetRPEInvalid", ErrTargetRPEInvalid, "target RPE must be between 7.0 and 10.0"},
		{"ErrRPEChartRequired", ErrRPEChartRequired, "RPE chart is required for RPE_TARGET strategy"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.msg {
				t.Errorf("expected message %q, got %q", tt.msg, tt.err.Error())
			}
		})
	}
}

// TestRPETargetLoadStrategy_Interface ensures the struct implements LoadStrategy
func TestRPETargetLoadStrategy_Interface(t *testing.T) {
	var _ LoadStrategy = (*RPETargetLoadStrategy)(nil)
}

// TestRPETargetRoundTripJSON tests that marshaling and unmarshaling produces equivalent strategies
func TestRPETargetRoundTripJSON(t *testing.T) {
	original := &RPETargetLoadStrategy{
		TargetReps:        5,
		TargetRPE:         8.5,
		RoundingIncrement: 2.5,
		RoundingDirection: RoundDown,
	}

	// Marshal
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Unmarshal
	restored, err := UnmarshalRPETarget(data)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	rs := restored.(*RPETargetLoadStrategy)

	// Compare
	if rs.TargetReps != original.TargetReps {
		t.Errorf("targetReps mismatch: expected %d, got %d", original.TargetReps, rs.TargetReps)
	}
	if rs.TargetRPE != original.TargetRPE {
		t.Errorf("targetRpe mismatch: expected %f, got %f", original.TargetRPE, rs.TargetRPE)
	}
	if rs.RoundingIncrement != original.RoundingIncrement {
		t.Errorf("roundingIncrement mismatch: expected %f, got %f", original.RoundingIncrement, rs.RoundingIncrement)
	}
	if rs.RoundingDirection != original.RoundingDirection {
		t.Errorf("roundingDirection mismatch: expected %s, got %s", original.RoundingDirection, rs.RoundingDirection)
	}
}

// TestCalculateLoadWithVariousRPEValues tests the strategy across the RPE chart
func TestCalculateLoadWithVariousRPEValues(t *testing.T) {
	lookup := newMockMaxLookup()
	lookup.SetMax("user-123", "squat-456", "ONE_RM", 400.0, "2024-01-15")
	chart := rpechart.NewDefaultRPEChart()

	// Test some key RPE chart values
	tests := []struct {
		reps     int
		rpe      float64
		expected float64
	}{
		{1, 10.0, 400.0},  // 400 * 1.00 = 400
		{1, 9.0, 380.0},   // 400 * 0.95 = 380
		{3, 8.0, 330.0},   // 400 * 0.82 = 328 -> 330
		{5, 8.0, 310.0},   // 400 * 0.77 = 308 -> 310
		{5, 9.0, 320.0},   // 400 * 0.80 = 320
		{8, 8.0, 265.0},   // 400 * 0.66 = 264 -> 265
		{10, 10.0, 265.0}, // 400 * 0.66 = 264 -> 265
	}

	for _, tt := range tests {
		t.Run(
			// Format test name
			"reps_" + string(rune('0'+tt.reps)) + "_rpe_" + string(rune('0'+int(tt.rpe))),
			func(t *testing.T) {
				strategy := &RPETargetLoadStrategy{
					TargetReps:        tt.reps,
					TargetRPE:         tt.rpe,
					RoundingIncrement: 5.0,
					RoundingDirection: RoundNearest,
				}
				strategy.SetMaxLookup(lookup)
				strategy.SetRPEChart(chart)

				params := LoadCalculationParams{
					UserID: "user-123",
					LiftID: "squat-456",
				}

				result, err := strategy.CalculateLoad(context.Background(), params)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if math.Abs(result-tt.expected) > 0.0001 {
					t.Errorf("at %d reps @ RPE %.1f: expected %f, got %f", tt.reps, tt.rpe, tt.expected, result)
				}
			},
		)
	}
}

// Benchmark tests
func BenchmarkRPETargetCalculateLoad(b *testing.B) {
	lookup := newMockMaxLookup()
	lookup.SetMax("user-123", "squat-456", "ONE_RM", 400.0, "2024-01-15")
	chart := rpechart.NewDefaultRPEChart()

	strategy := &RPETargetLoadStrategy{
		TargetReps:        5,
		TargetRPE:         8.0,
		RoundingIncrement: 5.0,
		RoundingDirection: RoundNearest,
	}
	strategy.SetMaxLookup(lookup)
	strategy.SetRPEChart(chart)

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

func BenchmarkRPETargetMarshalJSON(b *testing.B) {
	strategy := &RPETargetLoadStrategy{
		TargetReps:        5,
		TargetRPE:         8.0,
		RoundingIncrement: 5.0,
		RoundingDirection: RoundNearest,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Marshal(strategy)
	}
}

func BenchmarkUnmarshalRPETarget(b *testing.B) {
	data := []byte(`{
		"type": "RPE_TARGET",
		"targetReps": 5,
		"targetRpe": 8.0,
		"roundingIncrement": 5.0,
		"roundingDirection": "NEAREST"
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		UnmarshalRPETarget(data)
	}
}

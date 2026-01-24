package loadstrategy

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

func TestNewFindRMLoadStrategy(t *testing.T) {
	strategy := NewFindRMLoadStrategy(10)

	if strategy.TargetReps != 10 {
		t.Errorf("expected target reps 10, got %d", strategy.TargetReps)
	}
}

func TestFindRMLoadStrategy_Type(t *testing.T) {
	strategy := &FindRMLoadStrategy{}
	if strategy.Type() != TypeFindRM {
		t.Errorf("expected type %s, got %s", TypeFindRM, strategy.Type())
	}
}

func TestFindRMLoadStrategy_Validate(t *testing.T) {
	tests := []struct {
		name     string
		strategy FindRMLoadStrategy
		wantErr  error
	}{
		{
			name:     "valid 1 rep (1RM)",
			strategy: FindRMLoadStrategy{TargetReps: 1},
			wantErr:  nil,
		},
		{
			name:     "valid 5 reps (5RM)",
			strategy: FindRMLoadStrategy{TargetReps: 5},
			wantErr:  nil,
		},
		{
			name:     "valid 10 reps (10RM)",
			strategy: FindRMLoadStrategy{TargetReps: 10},
			wantErr:  nil,
		},
		{
			name:     "valid 12 reps (12RM - max)",
			strategy: FindRMLoadStrategy{TargetReps: 12},
			wantErr:  nil,
		},
		{
			name:     "invalid 0 reps",
			strategy: FindRMLoadStrategy{TargetReps: 0},
			wantErr:  ErrFindRMTargetRepsInvalid,
		},
		{
			name:     "invalid 13 reps (exceeds RPE chart range)",
			strategy: FindRMLoadStrategy{TargetReps: 13},
			wantErr:  ErrFindRMTargetRepsInvalid,
		},
		{
			name:     "invalid negative reps",
			strategy: FindRMLoadStrategy{TargetReps: -1},
			wantErr:  ErrFindRMTargetRepsInvalid,
		},
		{
			name:     "invalid 20 reps",
			strategy: FindRMLoadStrategy{TargetReps: 20},
			wantErr:  ErrFindRMTargetRepsInvalid,
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

func TestFindRMLoadStrategy_CalculateLoad(t *testing.T) {
	tests := []struct {
		name       string
		strategy   FindRMLoadStrategy
		params     LoadCalculationParams
		expected   float64
		wantErr    error
		wantErrMsg string
	}{
		{
			name:     "returns 0 for 10RM",
			strategy: FindRMLoadStrategy{TargetReps: 10},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat-456",
			},
			expected: 0,
		},
		{
			name:     "returns 0 for 5RM",
			strategy: FindRMLoadStrategy{TargetReps: 5},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat-456",
			},
			expected: 0,
		},
		{
			name:     "returns 0 for 1RM",
			strategy: FindRMLoadStrategy{TargetReps: 1},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "deadlift-789",
			},
			expected: 0,
		},
		{
			name:     "missing user ID in params",
			strategy: FindRMLoadStrategy{TargetReps: 10},
			params: LoadCalculationParams{
				UserID: "",
				LiftID: "squat-456",
			},
			wantErr: ErrInvalidParams,
		},
		{
			name:     "missing lift ID in params",
			strategy: FindRMLoadStrategy{TargetReps: 10},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "",
			},
			wantErr: ErrInvalidParams,
		},
		{
			name:     "invalid strategy (0 reps)",
			strategy: FindRMLoadStrategy{TargetReps: 0},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat-456",
			},
			wantErr: ErrFindRMTargetRepsInvalid,
		},
		{
			name:     "invalid strategy (13 reps)",
			strategy: FindRMLoadStrategy{TargetReps: 13},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat-456",
			},
			wantErr: ErrFindRMTargetRepsInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

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

			if result != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestFindRMLoadStrategy_MarshalJSON(t *testing.T) {
	strategy := &FindRMLoadStrategy{
		TargetReps: 10,
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
	} else if typeVal != string(TypeFindRM) {
		t.Errorf("expected type %s, got %v", TypeFindRM, typeVal)
	}

	// Check targetReps field
	if reps, ok := parsed["targetReps"]; !ok {
		t.Error("expected 'targetReps' field in marshaled JSON")
	} else if reps != 10.0 { // JSON numbers are float64
		t.Errorf("expected targetReps 10, got %v", reps)
	}
}

func TestUnmarshalFindRM(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
		check   func(*testing.T, LoadStrategy)
	}{
		{
			name: "valid 10RM",
			json: `{
				"type": "FIND_RM",
				"targetReps": 10
			}`,
			wantErr: false,
			check: func(t *testing.T, s LoadStrategy) {
				fs := s.(*FindRMLoadStrategy)
				if fs.TargetReps != 10 {
					t.Errorf("expected targetReps 10, got %d", fs.TargetReps)
				}
			},
		},
		{
			name: "valid 1RM",
			json: `{
				"type": "FIND_RM",
				"targetReps": 1
			}`,
			wantErr: false,
			check: func(t *testing.T, s LoadStrategy) {
				fs := s.(*FindRMLoadStrategy)
				if fs.TargetReps != 1 {
					t.Errorf("expected targetReps 1, got %d", fs.TargetReps)
				}
			},
		},
		{
			name: "valid 8RM",
			json: `{
				"type": "FIND_RM",
				"targetReps": 8
			}`,
			wantErr: false,
			check: func(t *testing.T, s LoadStrategy) {
				fs := s.(*FindRMLoadStrategy)
				if fs.TargetReps != 8 {
					t.Errorf("expected targetReps 8, got %d", fs.TargetReps)
				}
			},
		},
		{
			name: "valid 12RM (boundary)",
			json: `{
				"type": "FIND_RM",
				"targetReps": 12
			}`,
			wantErr: false,
			check: func(t *testing.T, s LoadStrategy) {
				fs := s.(*FindRMLoadStrategy)
				if fs.TargetReps != 12 {
					t.Errorf("expected targetReps 12, got %d", fs.TargetReps)
				}
			},
		},
		{
			name: "invalid 0 reps",
			json: `{
				"type": "FIND_RM",
				"targetReps": 0
			}`,
			wantErr: true,
		},
		{
			name: "invalid 13 reps",
			json: `{
				"type": "FIND_RM",
				"targetReps": 13
			}`,
			wantErr: true,
		},
		{
			name: "invalid negative reps",
			json: `{
				"type": "FIND_RM",
				"targetReps": -5
			}`,
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			json:    `{invalid}`,
			wantErr: true,
		},
		{
			name: "missing targetReps defaults to 0 (invalid)",
			json: `{
				"type": "FIND_RM"
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy, err := UnmarshalFindRM(json.RawMessage(tt.json))
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

func TestRegisterFindRM(t *testing.T) {
	factory := NewStrategyFactory()
	RegisterFindRM(factory)

	if !factory.IsRegistered(TypeFindRM) {
		t.Error("expected TypeFindRM to be registered")
	}

	// Test that we can create a strategy from JSON
	jsonData := []byte(`{
		"type": "FIND_RM",
		"targetReps": 10
	}`)

	strategy, err := factory.CreateFromJSON(jsonData)
	if err != nil {
		t.Fatalf("failed to create strategy: %v", err)
	}

	if strategy.Type() != TypeFindRM {
		t.Errorf("expected type %s, got %s", TypeFindRM, strategy.Type())
	}
}

func TestFindRMErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		msg  string
	}{
		{"ErrFindRMTargetRepsInvalid", ErrFindRMTargetRepsInvalid, "target reps must be between 1 and 12"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.msg {
				t.Errorf("expected message %q, got %q", tt.msg, tt.err.Error())
			}
		})
	}
}

// TestFindRMLoadStrategy_Interface ensures the struct implements LoadStrategy
func TestFindRMLoadStrategy_Interface(t *testing.T) {
	var _ LoadStrategy = (*FindRMLoadStrategy)(nil)
}

// TestFindRMRoundTripJSON tests that marshaling and unmarshaling produces equivalent strategies
func TestFindRMRoundTripJSON(t *testing.T) {
	original := &FindRMLoadStrategy{
		TargetReps: 8,
	}

	// Marshal
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Unmarshal
	restored, err := UnmarshalFindRM(data)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	fs := restored.(*FindRMLoadStrategy)

	// Compare
	if fs.TargetReps != original.TargetReps {
		t.Errorf("targetReps mismatch: expected %d, got %d", original.TargetReps, fs.TargetReps)
	}
}

// TestFindRMCalculateLoad_AllValidRepCounts tests all valid rep counts (1-12)
func TestFindRMCalculateLoad_AllValidRepCounts(t *testing.T) {
	for reps := 1; reps <= 12; reps++ {
		t.Run(
			"reps_"+string(rune('0'+reps/10))+string(rune('0'+reps%10)),
			func(t *testing.T) {
				strategy := &FindRMLoadStrategy{
					TargetReps: reps,
				}

				params := LoadCalculationParams{
					UserID: "user-123",
					LiftID: "squat-456",
				}

				result, err := strategy.CalculateLoad(context.Background(), params)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				// FindRM should always return 0 (no prescribed weight)
				if result != 0 {
					t.Errorf("at %d reps: expected 0, got %f", reps, result)
				}
			},
		)
	}
}

// TestTypeFindRMConstant verifies the constant value
func TestTypeFindRMConstant(t *testing.T) {
	if string(TypeFindRM) != "FIND_RM" {
		t.Errorf("expected TypeFindRM to be 'FIND_RM', got %s", TypeFindRM)
	}
}

// TestFindRMInValidStrategyTypes verifies TypeFindRM is in ValidStrategyTypes
func TestFindRMInValidStrategyTypes(t *testing.T) {
	if !ValidStrategyTypes[TypeFindRM] {
		t.Error("expected TypeFindRM to be in ValidStrategyTypes")
	}
}

// TestFindRMValidateStrategyType tests ValidateStrategyType with TypeFindRM
func TestFindRMValidateStrategyType(t *testing.T) {
	err := ValidateStrategyType(TypeFindRM)
	if err != nil {
		t.Errorf("ValidateStrategyType should accept TypeFindRM: %v", err)
	}
}

// Benchmark tests
func BenchmarkFindRMCalculateLoad(b *testing.B) {
	strategy := &FindRMLoadStrategy{
		TargetReps: 10,
	}

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

func BenchmarkFindRMMarshalJSON(b *testing.B) {
	strategy := &FindRMLoadStrategy{
		TargetReps: 10,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Marshal(strategy)
	}
}

func BenchmarkUnmarshalFindRM(b *testing.B) {
	data := []byte(`{
		"type": "FIND_RM",
		"targetReps": 10
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		UnmarshalFindRM(data)
	}
}

package loadstrategy

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

// mockSessionLookup implements SessionLookup for testing.
type mockSessionLookup struct {
	sets map[string]map[string]map[int]*LoggedSetResult // sessionID -> liftID -> setIndex -> result
	err  error
}

func newMockSessionLookup() *mockSessionLookup {
	return &mockSessionLookup{
		sets: make(map[string]map[string]map[int]*LoggedSetResult),
	}
}

func (m *mockSessionLookup) addSet(sessionID, liftID string, setIndex int, result *LoggedSetResult) {
	if m.sets[sessionID] == nil {
		m.sets[sessionID] = make(map[string]map[int]*LoggedSetResult)
	}
	if m.sets[sessionID][liftID] == nil {
		m.sets[sessionID][liftID] = make(map[int]*LoggedSetResult)
	}
	m.sets[sessionID][liftID][setIndex] = result
}

func (m *mockSessionLookup) GetLoggedSetByIndex(ctx context.Context, sessionID, liftID string, setIndex int) (*LoggedSetResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.sets[sessionID] == nil {
		return nil, nil
	}
	if m.sets[sessionID][liftID] == nil {
		return nil, nil
	}
	return m.sets[sessionID][liftID][setIndex], nil
}

func TestNewRelativeToLoadStrategy(t *testing.T) {
	mock := newMockSessionLookup()
	strategy := NewRelativeToLoadStrategy(0, 85.0, 5.0, RoundNearest, mock)

	if strategy.ReferenceSetIndex != 0 {
		t.Errorf("expected reference set index 0, got %d", strategy.ReferenceSetIndex)
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

func TestRelativeToLoadStrategy_Type(t *testing.T) {
	strategy := &RelativeToLoadStrategy{}
	if strategy.Type() != TypeRelativeTo {
		t.Errorf("expected type %s, got %s", TypeRelativeTo, strategy.Type())
	}
}

func TestRelativeToLoadStrategy_Validate(t *testing.T) {
	tests := []struct {
		name     string
		strategy RelativeToLoadStrategy
		wantErr  error
	}{
		{
			name: "valid basic strategy",
			strategy: RelativeToLoadStrategy{
				ReferenceSetIndex: 0,
				Percentage:        85.0,
			},
			wantErr: nil,
		},
		{
			name: "valid with all options",
			strategy: RelativeToLoadStrategy{
				ReferenceSetIndex: 0,
				Percentage:        85.0,
				RoundingIncrement: 2.5,
				RoundingDirection: RoundDown,
			},
			wantErr: nil,
		},
		{
			name: "valid high set index",
			strategy: RelativeToLoadStrategy{
				ReferenceSetIndex: 5,
				Percentage:        90.0,
			},
			wantErr: nil,
		},
		{
			name: "valid percentage over 100 (overload)",
			strategy: RelativeToLoadStrategy{
				ReferenceSetIndex: 0,
				Percentage:        105.0,
			},
			wantErr: nil,
		},
		{
			name: "invalid negative set index",
			strategy: RelativeToLoadStrategy{
				ReferenceSetIndex: -1,
				Percentage:        85.0,
			},
			wantErr: ErrReferenceSetIndexInvalid,
		},
		{
			name: "invalid zero percentage",
			strategy: RelativeToLoadStrategy{
				ReferenceSetIndex: 0,
				Percentage:        0,
			},
			wantErr: ErrRelativeToPercentageInvalid,
		},
		{
			name: "invalid negative percentage",
			strategy: RelativeToLoadStrategy{
				ReferenceSetIndex: 0,
				Percentage:        -10.0,
			},
			wantErr: ErrRelativeToPercentageInvalid,
		},
		{
			name: "invalid rounding direction",
			strategy: RelativeToLoadStrategy{
				ReferenceSetIndex: 0,
				Percentage:        85.0,
				RoundingDirection: "INVALID",
			},
			wantErr: ErrInvalidRoundingDirection,
		},
		{
			name: "invalid negative rounding increment",
			strategy: RelativeToLoadStrategy{
				ReferenceSetIndex: 0,
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

func TestRelativeToLoadStrategy_CalculateLoad(t *testing.T) {
	tests := []struct {
		name     string
		strategy func() *RelativeToLoadStrategy
		params   LoadCalculationParams
		expected float64
		wantErr  error
	}{
		{
			name: "basic calculation - 400 lbs top set × 85% = 340 lbs",
			strategy: func() *RelativeToLoadStrategy {
				mock := newMockSessionLookup()
				mock.addSet("session-1", "squat", 0, &LoggedSetResult{Weight: 400, Reps: 5})
				return NewRelativeToLoadStrategy(0, 85.0, 5.0, RoundNearest, mock)
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat",
				Context: map[string]interface{}{
					"sessionID": "session-1",
				},
			},
			expected: 340.0,
		},
		{
			name: "rounding - 405 lbs × 87% = 352.35 → 350 (round nearest 5)",
			strategy: func() *RelativeToLoadStrategy {
				mock := newMockSessionLookup()
				mock.addSet("session-1", "squat", 0, &LoggedSetResult{Weight: 405, Reps: 5})
				return NewRelativeToLoadStrategy(0, 87.0, 5.0, RoundNearest, mock)
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat",
				Context: map[string]interface{}{
					"sessionID": "session-1",
				},
			},
			expected: 350.0, // 405 * 0.87 = 352.35, rounds to 350 (nearest 5)
		},
		{
			name: "round down - 405 lbs × 87% = 352.35 → 350 (round down 5)",
			strategy: func() *RelativeToLoadStrategy {
				mock := newMockSessionLookup()
				mock.addSet("session-1", "squat", 0, &LoggedSetResult{Weight: 405, Reps: 5})
				return NewRelativeToLoadStrategy(0, 87.0, 5.0, RoundDown, mock)
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat",
				Context: map[string]interface{}{
					"sessionID": "session-1",
				},
			},
			expected: 350.0,
		},
		{
			name: "round up - 405 lbs × 87% = 352.35 → 355 (round up 5)",
			strategy: func() *RelativeToLoadStrategy {
				mock := newMockSessionLookup()
				mock.addSet("session-1", "squat", 0, &LoggedSetResult{Weight: 405, Reps: 5})
				return NewRelativeToLoadStrategy(0, 87.0, 5.0, RoundUp, mock)
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat",
				Context: map[string]interface{}{
					"sessionID": "session-1",
				},
			},
			expected: 355.0,
		},
		{
			name: "2.5 lb rounding increment",
			strategy: func() *RelativeToLoadStrategy {
				mock := newMockSessionLookup()
				mock.addSet("session-1", "bench", 0, &LoggedSetResult{Weight: 200, Reps: 3})
				return NewRelativeToLoadStrategy(0, 88.0, 2.5, RoundNearest, mock)
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "bench",
				Context: map[string]interface{}{
					"sessionID": "session-1",
				},
			},
			expected: 175.0, // 200 * 0.88 = 176, rounds to 175 with 2.5 increment
		},
		{
			name: "reference second set (index 1)",
			strategy: func() *RelativeToLoadStrategy {
				mock := newMockSessionLookup()
				mock.addSet("session-1", "squat", 0, &LoggedSetResult{Weight: 400, Reps: 5})
				mock.addSet("session-1", "squat", 1, &LoggedSetResult{Weight: 350, Reps: 8})
				return NewRelativeToLoadStrategy(1, 90.0, 5.0, RoundNearest, mock)
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat",
				Context: map[string]interface{}{
					"sessionID": "session-1",
				},
			},
			expected: 315.0, // 350 * 0.90 = 315
		},
		{
			name: "default rounding increment (5.0)",
			strategy: func() *RelativeToLoadStrategy {
				mock := newMockSessionLookup()
				mock.addSet("session-1", "squat", 0, &LoggedSetResult{Weight: 400, Reps: 5})
				s := NewRelativeToLoadStrategy(0, 85.0, 0, "", mock) // 0 increment = use default
				return s
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat",
				Context: map[string]interface{}{
					"sessionID": "session-1",
				},
			},
			expected: 340.0,
		},
		{
			name: "error - reference set not found",
			strategy: func() *RelativeToLoadStrategy {
				mock := newMockSessionLookup()
				// No sets added
				return NewRelativeToLoadStrategy(0, 85.0, 5.0, RoundNearest, mock)
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat",
				Context: map[string]interface{}{
					"sessionID": "session-1",
				},
			},
			wantErr: ErrReferenceSetNotFound,
		},
		{
			name: "error - missing session context",
			strategy: func() *RelativeToLoadStrategy {
				mock := newMockSessionLookup()
				mock.addSet("session-1", "squat", 0, &LoggedSetResult{Weight: 400, Reps: 5})
				return NewRelativeToLoadStrategy(0, 85.0, 5.0, RoundNearest, mock)
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat",
				// No Context
			},
			wantErr: ErrSessionIDRequired,
		},
		{
			name: "error - nil context map",
			strategy: func() *RelativeToLoadStrategy {
				mock := newMockSessionLookup()
				mock.addSet("session-1", "squat", 0, &LoggedSetResult{Weight: 400, Reps: 5})
				return NewRelativeToLoadStrategy(0, 85.0, 5.0, RoundNearest, mock)
			},
			params: LoadCalculationParams{
				UserID:  "user-123",
				LiftID:  "squat",
				Context: nil,
			},
			wantErr: ErrSessionIDRequired,
		},
		{
			name: "error - empty session ID in context",
			strategy: func() *RelativeToLoadStrategy {
				mock := newMockSessionLookup()
				mock.addSet("session-1", "squat", 0, &LoggedSetResult{Weight: 400, Reps: 5})
				return NewRelativeToLoadStrategy(0, 85.0, 5.0, RoundNearest, mock)
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat",
				Context: map[string]interface{}{
					"sessionID": "",
				},
			},
			wantErr: ErrSessionIDRequired,
		},
		{
			name: "error - missing session lookup",
			strategy: func() *RelativeToLoadStrategy {
				return &RelativeToLoadStrategy{
					ReferenceSetIndex: 0,
					Percentage:        85.0,
					// No session lookup set
				}
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat",
				Context: map[string]interface{}{
					"sessionID": "session-1",
				},
			},
			wantErr: ErrSessionLookupRequired,
		},
		{
			name: "error - missing user ID in params",
			strategy: func() *RelativeToLoadStrategy {
				mock := newMockSessionLookup()
				return NewRelativeToLoadStrategy(0, 85.0, 5.0, RoundNearest, mock)
			},
			params: LoadCalculationParams{
				UserID: "",
				LiftID: "squat",
				Context: map[string]interface{}{
					"sessionID": "session-1",
				},
			},
			wantErr: ErrInvalidParams,
		},
		{
			name: "error - missing lift ID in params",
			strategy: func() *RelativeToLoadStrategy {
				mock := newMockSessionLookup()
				return NewRelativeToLoadStrategy(0, 85.0, 5.0, RoundNearest, mock)
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "",
				Context: map[string]interface{}{
					"sessionID": "session-1",
				},
			},
			wantErr: ErrInvalidParams,
		},
		{
			name: "error - invalid strategy configuration",
			strategy: func() *RelativeToLoadStrategy {
				mock := newMockSessionLookup()
				mock.addSet("session-1", "squat", 0, &LoggedSetResult{Weight: 400, Reps: 5})
				return NewRelativeToLoadStrategy(-1, 85.0, 5.0, RoundNearest, mock) // Invalid index
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat",
				Context: map[string]interface{}{
					"sessionID": "session-1",
				},
			},
			wantErr: ErrReferenceSetIndexInvalid,
		},
		{
			name: "error - session lookup returns error",
			strategy: func() *RelativeToLoadStrategy {
				mock := newMockSessionLookup()
				mock.err = errors.New("database error")
				return NewRelativeToLoadStrategy(0, 85.0, 5.0, RoundNearest, mock)
			},
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "squat",
				Context: map[string]interface{}{
					"sessionID": "session-1",
				},
			},
			wantErr: nil, // Will contain wrapped error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			strategy := tt.strategy()

			result, err := strategy.CalculateLoad(ctx, tt.params)

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

			if tt.name == "error - session lookup returns error" {
				if err == nil {
					t.Error("expected error from session lookup, got nil")
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

func TestRelativeToLoadStrategy_SetSessionLookup(t *testing.T) {
	strategy := &RelativeToLoadStrategy{
		ReferenceSetIndex: 0,
		Percentage:        85.0,
	}

	mock := newMockSessionLookup()
	mock.addSet("session-1", "squat", 0, &LoggedSetResult{Weight: 400, Reps: 5})

	strategy.SetSessionLookup(mock)

	// Should now be able to calculate
	params := LoadCalculationParams{
		UserID: "user-123",
		LiftID: "squat",
		Context: map[string]interface{}{
			"sessionID": "session-1",
		},
	}

	result, err := strategy.CalculateLoad(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := 340.0 // 400 * 0.85
	if result != expected {
		t.Errorf("expected %f, got %f", expected, result)
	}
}

func TestRelativeToLoadStrategy_MarshalJSON(t *testing.T) {
	strategy := &RelativeToLoadStrategy{
		ReferenceSetIndex: 0,
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
	} else if typeVal != string(TypeRelativeTo) {
		t.Errorf("expected type %s, got %v", TypeRelativeTo, typeVal)
	}

	// Check referenceSetIndex field
	if idx, ok := parsed["referenceSetIndex"]; !ok {
		t.Error("expected 'referenceSetIndex' field in marshaled JSON")
	} else if idx != 0.0 { // JSON numbers are float64
		t.Errorf("expected referenceSetIndex 0, got %v", idx)
	}

	// Check percentage field
	if pct, ok := parsed["percentage"]; !ok {
		t.Error("expected 'percentage' field in marshaled JSON")
	} else if pct != 85.0 {
		t.Errorf("expected percentage 85, got %v", pct)
	}

	// Check roundingIncrement field
	if inc, ok := parsed["roundingIncrement"]; !ok {
		t.Error("expected 'roundingIncrement' field in marshaled JSON")
	} else if inc != 5.0 {
		t.Errorf("expected roundingIncrement 5.0, got %v", inc)
	}

	// Check roundingDirection field
	if dir, ok := parsed["roundingDirection"]; !ok {
		t.Error("expected 'roundingDirection' field in marshaled JSON")
	} else if dir != string(RoundNearest) {
		t.Errorf("expected roundingDirection %s, got %v", RoundNearest, dir)
	}
}

func TestUnmarshalRelativeTo(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
		check   func(*testing.T, LoadStrategy)
	}{
		{
			name: "valid basic",
			json: `{
				"type": "RELATIVE_TO",
				"referenceSetIndex": 0,
				"percentage": 85.0
			}`,
			wantErr: false,
			check: func(t *testing.T, s LoadStrategy) {
				rs := s.(*RelativeToLoadStrategy)
				if rs.ReferenceSetIndex != 0 {
					t.Errorf("expected referenceSetIndex 0, got %d", rs.ReferenceSetIndex)
				}
				if rs.Percentage != 85.0 {
					t.Errorf("expected percentage 85.0, got %f", rs.Percentage)
				}
			},
		},
		{
			name: "valid with all options",
			json: `{
				"type": "RELATIVE_TO",
				"referenceSetIndex": 0,
				"percentage": 85.0,
				"roundingIncrement": 2.5,
				"roundingDirection": "DOWN"
			}`,
			wantErr: false,
			check: func(t *testing.T, s LoadStrategy) {
				rs := s.(*RelativeToLoadStrategy)
				if rs.RoundingIncrement != 2.5 {
					t.Errorf("expected roundingIncrement 2.5, got %f", rs.RoundingIncrement)
				}
				if rs.RoundingDirection != RoundDown {
					t.Errorf("expected roundingDirection DOWN, got %s", rs.RoundingDirection)
				}
			},
		},
		{
			name: "valid high set index",
			json: `{
				"type": "RELATIVE_TO",
				"referenceSetIndex": 5,
				"percentage": 90.0
			}`,
			wantErr: false,
			check: func(t *testing.T, s LoadStrategy) {
				rs := s.(*RelativeToLoadStrategy)
				if rs.ReferenceSetIndex != 5 {
					t.Errorf("expected referenceSetIndex 5, got %d", rs.ReferenceSetIndex)
				}
			},
		},
		{
			name: "invalid negative set index",
			json: `{
				"type": "RELATIVE_TO",
				"referenceSetIndex": -1,
				"percentage": 85.0
			}`,
			wantErr: true,
		},
		{
			name: "invalid zero percentage",
			json: `{
				"type": "RELATIVE_TO",
				"referenceSetIndex": 0,
				"percentage": 0
			}`,
			wantErr: true,
		},
		{
			name: "invalid negative percentage",
			json: `{
				"type": "RELATIVE_TO",
				"referenceSetIndex": 0,
				"percentage": -10.0
			}`,
			wantErr: true,
		},
		{
			name: "invalid rounding direction",
			json: `{
				"type": "RELATIVE_TO",
				"referenceSetIndex": 0,
				"percentage": 85.0,
				"roundingDirection": "INVALID"
			}`,
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			json:    `{invalid}`,
			wantErr: true,
		},
		{
			name: "missing percentage defaults to 0 (invalid)",
			json: `{
				"type": "RELATIVE_TO",
				"referenceSetIndex": 0
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy, err := UnmarshalRelativeTo(json.RawMessage(tt.json))
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

func TestRegisterRelativeTo(t *testing.T) {
	factory := NewStrategyFactory()
	RegisterRelativeTo(factory)

	if !factory.IsRegistered(TypeRelativeTo) {
		t.Error("expected TypeRelativeTo to be registered")
	}

	// Test that we can create a strategy from JSON
	jsonData := []byte(`{
		"type": "RELATIVE_TO",
		"referenceSetIndex": 0,
		"percentage": 85.0
	}`)

	strategy, err := factory.CreateFromJSON(jsonData)
	if err != nil {
		t.Fatalf("failed to create strategy: %v", err)
	}

	if strategy.Type() != TypeRelativeTo {
		t.Errorf("expected type %s, got %s", TypeRelativeTo, strategy.Type())
	}
}

func TestRelativeToErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		msg  string
	}{
		{"ErrReferenceSetIndexInvalid", ErrReferenceSetIndexInvalid, "reference set index must be non-negative"},
		{"ErrRelativeToPercentageInvalid", ErrRelativeToPercentageInvalid, "percentage must be greater than 0"},
		{"ErrSessionLookupRequired", ErrSessionLookupRequired, "session lookup is required for RELATIVE_TO strategy"},
		{"ErrSessionIDRequired", ErrSessionIDRequired, "session ID is required in context for RELATIVE_TO strategy"},
		{"ErrReferenceSetNotFound", ErrReferenceSetNotFound, "reference set not found (may not be logged yet)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.msg {
				t.Errorf("expected message %q, got %q", tt.msg, tt.err.Error())
			}
		})
	}
}

// TestRelativeToLoadStrategy_Interface ensures the struct implements LoadStrategy
func TestRelativeToLoadStrategy_Interface(t *testing.T) {
	var _ LoadStrategy = (*RelativeToLoadStrategy)(nil)
}

// TestRelativeToRoundTripJSON tests that marshaling and unmarshaling produces equivalent strategies
func TestRelativeToRoundTripJSON(t *testing.T) {
	original := &RelativeToLoadStrategy{
		ReferenceSetIndex: 0,
		Percentage:        85.0,
		RoundingIncrement: 2.5,
		RoundingDirection: RoundUp,
	}

	// Marshal
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Unmarshal
	restored, err := UnmarshalRelativeTo(data)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	rs := restored.(*RelativeToLoadStrategy)

	// Compare
	if rs.ReferenceSetIndex != original.ReferenceSetIndex {
		t.Errorf("referenceSetIndex mismatch: expected %d, got %d", original.ReferenceSetIndex, rs.ReferenceSetIndex)
	}
	if rs.Percentage != original.Percentage {
		t.Errorf("percentage mismatch: expected %f, got %f", original.Percentage, rs.Percentage)
	}
	if rs.RoundingIncrement != original.RoundingIncrement {
		t.Errorf("roundingIncrement mismatch: expected %f, got %f", original.RoundingIncrement, rs.RoundingIncrement)
	}
	if rs.RoundingDirection != original.RoundingDirection {
		t.Errorf("roundingDirection mismatch: expected %s, got %s", original.RoundingDirection, rs.RoundingDirection)
	}
}

// TestTypeRelativeToConstant verifies the constant value
func TestTypeRelativeToConstant(t *testing.T) {
	if string(TypeRelativeTo) != "RELATIVE_TO" {
		t.Errorf("expected TypeRelativeTo to be 'RELATIVE_TO', got %s", TypeRelativeTo)
	}
}

// TestRelativeToInValidStrategyTypes verifies TypeRelativeTo is in ValidStrategyTypes
func TestRelativeToInValidStrategyTypes(t *testing.T) {
	if !ValidStrategyTypes[TypeRelativeTo] {
		t.Error("expected TypeRelativeTo to be in ValidStrategyTypes")
	}
}

// TestRelativeToValidateStrategyType tests ValidateStrategyType with TypeRelativeTo
func TestRelativeToValidateStrategyType(t *testing.T) {
	err := ValidateStrategyType(TypeRelativeTo)
	if err != nil {
		t.Errorf("ValidateStrategyType should accept TypeRelativeTo: %v", err)
	}
}

// TestLoggedSetResult tests the LoggedSetResult struct
func TestLoggedSetResult(t *testing.T) {
	rpe := 8.5
	result := LoggedSetResult{
		Weight: 400.0,
		Reps:   5,
		RPE:    &rpe,
	}

	if result.Weight != 400.0 {
		t.Errorf("expected weight 400.0, got %f", result.Weight)
	}
	if result.Reps != 5 {
		t.Errorf("expected reps 5, got %d", result.Reps)
	}
	if result.RPE == nil {
		t.Error("expected RPE to be set")
	} else if *result.RPE != 8.5 {
		t.Errorf("expected RPE 8.5, got %f", *result.RPE)
	}
}

// TestLoggedSetResultNilRPE tests LoggedSetResult with nil RPE
func TestLoggedSetResultNilRPE(t *testing.T) {
	result := LoggedSetResult{
		Weight: 400.0,
		Reps:   5,
		RPE:    nil,
	}

	if result.RPE != nil {
		t.Error("expected RPE to be nil")
	}
}

// TestExtractSessionID tests the private extractSessionID method through CalculateLoad
func TestExtractSessionID_WrongType(t *testing.T) {
	mock := newMockSessionLookup()
	mock.addSet("session-1", "squat", 0, &LoggedSetResult{Weight: 400, Reps: 5})
	strategy := NewRelativeToLoadStrategy(0, 85.0, 5.0, RoundNearest, mock)

	// Session ID is wrong type (int instead of string)
	params := LoadCalculationParams{
		UserID: "user-123",
		LiftID: "squat",
		Context: map[string]interface{}{
			"sessionID": 12345, // Wrong type!
		},
	}

	_, err := strategy.CalculateLoad(context.Background(), params)
	if err == nil {
		t.Error("expected error for wrong sessionID type")
	}
	if !errors.Is(err, ErrSessionIDRequired) {
		t.Errorf("expected ErrSessionIDRequired, got %v", err)
	}
}

// Benchmark tests
func BenchmarkRelativeToCalculateLoad(b *testing.B) {
	mock := newMockSessionLookup()
	mock.addSet("session-1", "squat", 0, &LoggedSetResult{Weight: 400, Reps: 5})
	strategy := NewRelativeToLoadStrategy(0, 85.0, 5.0, RoundNearest, mock)

	params := LoadCalculationParams{
		UserID: "user-123",
		LiftID: "squat",
		Context: map[string]interface{}{
			"sessionID": "session-1",
		},
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strategy.CalculateLoad(ctx, params)
	}
}

func BenchmarkRelativeToMarshalJSON(b *testing.B) {
	strategy := &RelativeToLoadStrategy{
		ReferenceSetIndex: 0,
		Percentage:        85.0,
		RoundingIncrement: 5.0,
		RoundingDirection: RoundNearest,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Marshal(strategy)
	}
}

func BenchmarkUnmarshalRelativeTo(b *testing.B) {
	data := []byte(`{
		"type": "RELATIVE_TO",
		"referenceSetIndex": 0,
		"percentage": 85.0,
		"roundingIncrement": 5.0,
		"roundingDirection": "NEAREST"
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		UnmarshalRelativeTo(data)
	}
}

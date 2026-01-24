package setscheme

import (
	"encoding/json"
	"errors"
	"testing"
)

// === RPEThreshold Tests ===

func TestRPEThreshold_Type(t *testing.T) {
	cond := &RPEThreshold{Threshold: 10}
	if cond.Type() != TerminationTypeRPEThreshold {
		t.Errorf("expected type %s, got %s", TerminationTypeRPEThreshold, cond.Type())
	}
}

func TestNewRPEThreshold(t *testing.T) {
	tests := []struct {
		name      string
		threshold float64
		wantErr   bool
	}{
		{"valid RPE 10", 10, false},
		{"valid RPE 9", 9, false},
		{"valid RPE 8.5", 8.5, false},
		{"valid RPE 1", 1, false},
		{"valid RPE 5", 5, false},
		{"invalid RPE 0", 0, true},
		{"invalid RPE -1", -1, true},
		{"invalid RPE 11", 11, true},
		{"invalid RPE 0.5", 0.5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond, err := NewRPEThreshold(tt.threshold)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if cond != nil {
					t.Error("expected nil condition on error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if cond == nil {
					t.Fatal("expected non-nil condition")
				}
				if cond.Threshold != tt.threshold {
					t.Errorf("expected threshold %v, got %v", tt.threshold, cond.Threshold)
				}
			}
		})
	}
}

func TestRPEThreshold_ShouldTerminate(t *testing.T) {
	cond := &RPEThreshold{Threshold: 9}

	tests := []struct {
		name     string
		ctx      TerminationContext
		expected bool
	}{
		{
			name:     "RPE above threshold",
			ctx:      TerminationContext{LastRPE: floatPtr(9.5)},
			expected: true,
		},
		{
			name:     "RPE at threshold",
			ctx:      TerminationContext{LastRPE: floatPtr(9)},
			expected: true,
		},
		{
			name:     "RPE below threshold",
			ctx:      TerminationContext{LastRPE: floatPtr(8)},
			expected: false,
		},
		{
			name:     "RPE well below threshold",
			ctx:      TerminationContext{LastRPE: floatPtr(5)},
			expected: false,
		},
		{
			name:     "no RPE provided",
			ctx:      TerminationContext{LastRPE: nil},
			expected: false,
		},
		{
			name:     "RPE exactly 10 at threshold 9",
			ctx:      TerminationContext{LastRPE: floatPtr(10)},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cond.ShouldTerminate(tt.ctx)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestRPEThreshold_ShouldTerminate_Threshold10(t *testing.T) {
	cond := &RPEThreshold{Threshold: 10}

	tests := []struct {
		name     string
		ctx      TerminationContext
		expected bool
	}{
		{
			name:     "RPE exactly 10",
			ctx:      TerminationContext{LastRPE: floatPtr(10)},
			expected: true,
		},
		{
			name:     "RPE below 10",
			ctx:      TerminationContext{LastRPE: floatPtr(9.5)},
			expected: false,
		},
		{
			name:     "RPE just under 10",
			ctx:      TerminationContext{LastRPE: floatPtr(9.9)},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cond.ShouldTerminate(tt.ctx)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestRPEThreshold_Validate(t *testing.T) {
	tests := []struct {
		name      string
		threshold float64
		wantErr   bool
	}{
		{"valid 10", 10, false},
		{"valid 1", 1, false},
		{"valid 5.5", 5.5, false},
		{"valid 7", 7, false},
		{"invalid 0", 0, true},
		{"invalid 11", 11, true},
		{"invalid -1", -1, true},
		{"invalid 0.99", 0.99, true},
		{"invalid 10.01", 10.01, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond := &RPEThreshold{Threshold: tt.threshold}
			err := cond.Validate()
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if !errors.Is(err, ErrInvalidTermination) {
					t.Errorf("expected ErrInvalidTermination, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestRPEThreshold_MarshalJSON(t *testing.T) {
	cond := &RPEThreshold{Threshold: 9}
	data, err := json.Marshal(cond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if parsed["type"] != string(TerminationTypeRPEThreshold) {
		t.Errorf("expected type %s, got %v", TerminationTypeRPEThreshold, parsed["type"])
	}
	if parsed["threshold"].(float64) != 9 {
		t.Errorf("expected threshold 9, got %v", parsed["threshold"])
	}
}

// === RepFailure Tests ===

func TestRepFailure_Type(t *testing.T) {
	cond := NewRepFailure()
	if cond.Type() != TerminationTypeRepFailure {
		t.Errorf("expected type %s, got %s", TerminationTypeRepFailure, cond.Type())
	}
}

func TestRepFailure_ShouldTerminate(t *testing.T) {
	cond := NewRepFailure()

	tests := []struct {
		name     string
		ctx      TerminationContext
		expected bool
	}{
		{
			name:     "reps below target",
			ctx:      TerminationContext{LastReps: 4, TargetReps: 5},
			expected: true,
		},
		{
			name:     "reps at target",
			ctx:      TerminationContext{LastReps: 5, TargetReps: 5},
			expected: false,
		},
		{
			name:     "reps above target",
			ctx:      TerminationContext{LastReps: 7, TargetReps: 5},
			expected: false,
		},
		{
			name:     "zero reps performed",
			ctx:      TerminationContext{LastReps: 0, TargetReps: 5},
			expected: true,
		},
		{
			name:     "one rep short",
			ctx:      TerminationContext{LastReps: 4, TargetReps: 5},
			expected: true,
		},
		{
			name:     "target is 1, performed 1",
			ctx:      TerminationContext{LastReps: 1, TargetReps: 1},
			expected: false,
		},
		{
			name:     "target is 1, performed 0",
			ctx:      TerminationContext{LastReps: 0, TargetReps: 1},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cond.ShouldTerminate(tt.ctx)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestRepFailure_Validate(t *testing.T) {
	cond := NewRepFailure()
	err := cond.Validate()
	if err != nil {
		t.Errorf("RepFailure.Validate() should always return nil, got %v", err)
	}
}

func TestRepFailure_MarshalJSON(t *testing.T) {
	cond := NewRepFailure()
	data, err := json.Marshal(cond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if parsed["type"] != string(TerminationTypeRepFailure) {
		t.Errorf("expected type %s, got %v", TerminationTypeRepFailure, parsed["type"])
	}
}

// === MaxSets Tests ===

func TestMaxSets_Type(t *testing.T) {
	cond := &MaxSets{Max: 10}
	if cond.Type() != TerminationTypeMaxSets {
		t.Errorf("expected type %s, got %s", TerminationTypeMaxSets, cond.Type())
	}
}

func TestNewMaxSets(t *testing.T) {
	tests := []struct {
		name    string
		max     int
		wantErr bool
	}{
		{"valid 10", 10, false},
		{"valid 1", 1, false},
		{"valid 5", 5, false},
		{"valid 100", 100, false},
		{"invalid 0", 0, true},
		{"invalid -1", -1, true},
		{"invalid -10", -10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond, err := NewMaxSets(tt.max)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if cond != nil {
					t.Error("expected nil condition on error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if cond == nil {
					t.Fatal("expected non-nil condition")
				}
				if cond.Max != tt.max {
					t.Errorf("expected max %d, got %d", tt.max, cond.Max)
				}
			}
		})
	}
}

func TestMaxSets_ShouldTerminate(t *testing.T) {
	cond := &MaxSets{Max: 5}

	tests := []struct {
		name     string
		ctx      TerminationContext
		expected bool
	}{
		{
			name:     "below max",
			ctx:      TerminationContext{TotalSets: 3},
			expected: false,
		},
		{
			name:     "one below max",
			ctx:      TerminationContext{TotalSets: 4},
			expected: false,
		},
		{
			name:     "at max",
			ctx:      TerminationContext{TotalSets: 5},
			expected: true,
		},
		{
			name:     "above max",
			ctx:      TerminationContext{TotalSets: 6},
			expected: true,
		},
		{
			name:     "zero sets",
			ctx:      TerminationContext{TotalSets: 0},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cond.ShouldTerminate(tt.ctx)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMaxSets_Validate(t *testing.T) {
	tests := []struct {
		name    string
		max     int
		wantErr bool
	}{
		{"valid 1", 1, false},
		{"valid 10", 10, false},
		{"valid 100", 100, false},
		{"invalid 0", 0, true},
		{"invalid -1", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond := &MaxSets{Max: tt.max}
			err := cond.Validate()
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if !errors.Is(err, ErrInvalidTermination) {
					t.Errorf("expected ErrInvalidTermination, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestMaxSets_MarshalJSON(t *testing.T) {
	cond := &MaxSets{Max: 10}
	data, err := json.Marshal(cond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if parsed["type"] != string(TerminationTypeMaxSets) {
		t.Errorf("expected type %s, got %v", TerminationTypeMaxSets, parsed["type"])
	}
	if int(parsed["max"].(float64)) != 10 {
		t.Errorf("expected max 10, got %v", parsed["max"])
	}
}

// === UnmarshalTerminationCondition Tests ===

func TestUnmarshalTerminationCondition(t *testing.T) {
	tests := []struct {
		name        string
		json        string
		expectType  TerminationConditionType
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid RPEThreshold",
			json:       `{"type": "RPE_THRESHOLD", "threshold": 9}`,
			expectType: TerminationTypeRPEThreshold,
			wantErr:    false,
		},
		{
			name:       "valid RPEThreshold with 10",
			json:       `{"type": "RPE_THRESHOLD", "threshold": 10}`,
			expectType: TerminationTypeRPEThreshold,
			wantErr:    false,
		},
		{
			name:       "valid RepFailure",
			json:       `{"type": "REP_FAILURE"}`,
			expectType: TerminationTypeRepFailure,
			wantErr:    false,
		},
		{
			name:       "valid MaxSets",
			json:       `{"type": "MAX_SETS", "max": 10}`,
			expectType: TerminationTypeMaxSets,
			wantErr:    false,
		},
		{
			name:        "invalid RPEThreshold",
			json:        `{"type": "RPE_THRESHOLD", "threshold": 11}`,
			wantErr:     true,
			errContains: "threshold",
		},
		{
			name:        "invalid MaxSets",
			json:        `{"type": "MAX_SETS", "max": 0}`,
			wantErr:     true,
			errContains: "max",
		},
		{
			name:        "unknown type",
			json:        `{"type": "UNKNOWN"}`,
			wantErr:     true,
			errContains: "unknown type",
		},
		{
			name:    "invalid JSON",
			json:    `{invalid}`,
			wantErr: true,
		},
		{
			name:    "empty JSON",
			json:    `{}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond, err := UnmarshalTerminationCondition(json.RawMessage(tt.json))
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if cond == nil {
					t.Fatal("expected non-nil condition")
				}
				if cond.Type() != tt.expectType {
					t.Errorf("expected type %s, got %s", tt.expectType, cond.Type())
				}
			}
		})
	}
}

// === Round-trip Tests ===

func TestRPEThreshold_RoundTrip(t *testing.T) {
	original := &RPEThreshold{Threshold: 8.5}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	cond, err := UnmarshalTerminationCondition(data)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	rpe, ok := cond.(*RPEThreshold)
	if !ok {
		t.Fatal("expected *RPEThreshold")
	}
	if rpe.Threshold != original.Threshold {
		t.Errorf("expected threshold %v, got %v", original.Threshold, rpe.Threshold)
	}
}

func TestMaxSets_RoundTrip(t *testing.T) {
	original := &MaxSets{Max: 15}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	cond, err := UnmarshalTerminationCondition(data)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	maxSets, ok := cond.(*MaxSets)
	if !ok {
		t.Fatal("expected *MaxSets")
	}
	if maxSets.Max != original.Max {
		t.Errorf("expected max %d, got %d", original.Max, maxSets.Max)
	}
}

func TestRepFailure_RoundTrip(t *testing.T) {
	original := NewRepFailure()

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	cond, err := UnmarshalTerminationCondition(data)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	_, ok := cond.(*RepFailure)
	if !ok {
		t.Fatal("expected *RepFailure")
	}
}

// === Interface Tests ===

func TestTerminationCondition_InterfaceCompliance(t *testing.T) {
	// Verify all types implement TerminationCondition
	var _ TerminationCondition = (*RPEThreshold)(nil)
	var _ TerminationCondition = (*RepFailure)(nil)
	var _ TerminationCondition = (*MaxSets)(nil)
}

// === TerminationContext Tests ===

func TestTerminationContext_Fields(t *testing.T) {
	rpe := 8.5
	ctx := TerminationContext{
		SetNumber:  3,
		LastRPE:    &rpe,
		LastReps:   5,
		TotalReps:  15,
		TotalSets:  3,
		TargetReps: 5,
	}

	if ctx.SetNumber != 3 {
		t.Errorf("expected SetNumber 3, got %d", ctx.SetNumber)
	}
	if *ctx.LastRPE != 8.5 {
		t.Errorf("expected LastRPE 8.5, got %v", *ctx.LastRPE)
	}
	if ctx.LastReps != 5 {
		t.Errorf("expected LastReps 5, got %d", ctx.LastReps)
	}
	if ctx.TotalReps != 15 {
		t.Errorf("expected TotalReps 15, got %d", ctx.TotalReps)
	}
	if ctx.TotalSets != 3 {
		t.Errorf("expected TotalSets 3, got %d", ctx.TotalSets)
	}
	if ctx.TargetReps != 5 {
		t.Errorf("expected TargetReps 5, got %d", ctx.TargetReps)
	}
}

// === Real-World Scenario Tests ===

func TestTermination_FatigueDropSetScenario(t *testing.T) {
	// Scenario: Lifter doing sets of 5 until they can't hit 5 reps
	repFailure := NewRepFailure()
	maxSets, _ := NewMaxSets(10) // Safety limit

	sets := []struct {
		reps        int
		shouldStop  bool
		description string
	}{
		{5, false, "First set - hit target"},
		{5, false, "Second set - hit target"},
		{5, false, "Third set - hit target"},
		{4, true, "Fourth set - missed target"},
	}

	for i, s := range sets {
		ctx := TerminationContext{
			SetNumber:  i + 1,
			LastReps:   s.reps,
			TargetReps: 5,
			TotalSets:  i + 1,
		}

		shouldStopRep := repFailure.ShouldTerminate(ctx)
		shouldStopMax := maxSets.ShouldTerminate(ctx)
		shouldStop := shouldStopRep || shouldStopMax

		if shouldStop != s.shouldStop {
			t.Errorf("Set %d (%s): expected shouldStop=%v, got %v", i+1, s.description, s.shouldStop, shouldStop)
		}
	}
}

func TestTermination_RPE10StopScenario(t *testing.T) {
	// Scenario: Lifter doing sets until RPE hits 10
	rpeThreshold, _ := NewRPEThreshold(10)
	maxSets, _ := NewMaxSets(8) // Safety limit

	sets := []struct {
		rpe         float64
		shouldStop  bool
		description string
	}{
		{7, false, "First set - easy"},
		{8, false, "Second set - moderate"},
		{9, false, "Third set - hard"},
		{9.5, false, "Fourth set - very hard"},
		{10, true, "Fifth set - maximal"},
	}

	for i, s := range sets {
		ctx := TerminationContext{
			SetNumber: i + 1,
			LastRPE:   &s.rpe,
			TotalSets: i + 1,
		}

		shouldStopRPE := rpeThreshold.ShouldTerminate(ctx)
		shouldStopMax := maxSets.ShouldTerminate(ctx)
		shouldStop := shouldStopRPE || shouldStopMax

		if shouldStop != s.shouldStop {
			t.Errorf("Set %d (%s): expected shouldStop=%v, got %v", i+1, s.description, s.shouldStop, shouldStop)
		}
	}
}

func TestTermination_MaxSetsSafetyLimit(t *testing.T) {
	// Scenario: Lifter keeps hitting reps but we have a safety limit
	maxSets, _ := NewMaxSets(5)

	for i := 1; i <= 6; i++ {
		ctx := TerminationContext{
			TotalSets:  i,
			LastReps:   5,
			TargetReps: 5,
		}

		shouldStop := maxSets.ShouldTerminate(ctx)
		expectedStop := i >= 5

		if shouldStop != expectedStop {
			t.Errorf("Set %d: expected shouldStop=%v, got %v", i, expectedStop, shouldStop)
		}
	}
}

// Helper function to create a float64 pointer
func floatPtr(f float64) *float64 {
	return &f
}

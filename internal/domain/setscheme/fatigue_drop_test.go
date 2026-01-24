package setscheme

import (
	"encoding/json"
	"errors"
	"testing"
)

// Helper function to create a float64 pointer (if not already defined)
func floatPtrFD(f float64) *float64 {
	return &f
}

// === Type Tests ===

func TestFatigueDrop_Type(t *testing.T) {
	scheme := &FatigueDrop{
		TargetReps:  3,
		StartRPE:    8,
		StopRPE:     10,
		DropPercent: 0.05,
		MaxSets:     10,
	}
	if scheme.Type() != TypeFatigueDrop {
		t.Errorf("expected type %s, got %s", TypeFatigueDrop, scheme.Type())
	}
}

// === Constructor Tests ===

func TestNewFatigueDrop(t *testing.T) {
	tests := []struct {
		name        string
		targetReps  int
		startRPE    float64
		stopRPE     float64
		dropPercent float64
		maxSets     int
		wantErr     bool
	}{
		{"valid basic", 3, 8, 10, 0.05, 10, false},
		{"valid with different reps", 5, 7, 9, 0.03, 8, false},
		{"valid min reps", 1, 6, 8, 0.10, 5, false},
		{"valid zero max sets (uses default)", 3, 8, 10, 0.05, 0, false},
		{"valid zero drop percent", 3, 8, 10, 0, 10, false},
		{"valid min RPE range", 3, 1, 2, 0.05, 10, false},
		{"valid max RPE range", 3, 9, 10, 0.05, 10, false},
		{"invalid zero reps", 0, 8, 10, 0.05, 10, true},
		{"invalid negative reps", -1, 8, 10, 0.05, 10, true},
		{"invalid start RPE too low", 3, 0, 10, 0.05, 10, true},
		{"invalid start RPE too high", 3, 11, 12, 0.05, 10, true},
		{"invalid stop RPE too low", 3, 8, 0, 0.05, 10, true},
		{"invalid stop RPE too high", 3, 8, 11, 0.05, 10, true},
		{"invalid stop <= start", 3, 8, 8, 0.05, 10, true},
		{"invalid stop < start", 3, 9, 8, 0.05, 10, true},
		{"invalid drop percent negative", 3, 8, 10, -0.05, 10, true},
		{"invalid drop percent > 1", 3, 8, 10, 1.5, 10, true},
		{"invalid negative max sets", 3, 8, 10, 0.05, -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme, err := NewFatigueDrop(tt.targetReps, tt.startRPE, tt.stopRPE, tt.dropPercent, tt.maxSets)
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
				if scheme.TargetReps != tt.targetReps {
					t.Errorf("expected TargetReps %d, got %d", tt.targetReps, scheme.TargetReps)
				}
				if scheme.StartRPE != tt.startRPE {
					t.Errorf("expected StartRPE %v, got %v", tt.startRPE, scheme.StartRPE)
				}
				if scheme.StopRPE != tt.stopRPE {
					t.Errorf("expected StopRPE %v, got %v", tt.stopRPE, scheme.StopRPE)
				}
				if scheme.DropPercent != tt.dropPercent {
					t.Errorf("expected DropPercent %v, got %v", tt.dropPercent, scheme.DropPercent)
				}
				if scheme.MaxSets != tt.maxSets {
					t.Errorf("expected MaxSets %d, got %d", tt.maxSets, scheme.MaxSets)
				}
			}
		})
	}
}

// === Validate Tests ===

func TestFatigueDrop_Validate(t *testing.T) {
	tests := []struct {
		name    string
		scheme  FatigueDrop
		wantErr bool
	}{
		{
			name: "valid",
			scheme: FatigueDrop{
				TargetReps: 3, StartRPE: 8, StopRPE: 10, DropPercent: 0.05, MaxSets: 10,
			},
			wantErr: false,
		},
		{
			name: "valid zero max sets",
			scheme: FatigueDrop{
				TargetReps: 3, StartRPE: 8, StopRPE: 10, DropPercent: 0.05, MaxSets: 0,
			},
			wantErr: false,
		},
		{
			name: "valid boundary RPE",
			scheme: FatigueDrop{
				TargetReps: 3, StartRPE: 9.5, StopRPE: 10, DropPercent: 0.05, MaxSets: 10,
			},
			wantErr: false,
		},
		{
			name: "invalid zero reps",
			scheme: FatigueDrop{
				TargetReps: 0, StartRPE: 8, StopRPE: 10, DropPercent: 0.05, MaxSets: 10,
			},
			wantErr: true,
		},
		{
			name: "invalid stop equals start",
			scheme: FatigueDrop{
				TargetReps: 3, StartRPE: 10, StopRPE: 10, DropPercent: 0.05, MaxSets: 10,
			},
			wantErr: true,
		},
		{
			name: "invalid drop percent above 1",
			scheme: FatigueDrop{
				TargetReps: 3, StartRPE: 8, StopRPE: 10, DropPercent: 1.0001, MaxSets: 10,
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

// === VariableSetScheme Interface Tests ===

func TestFatigueDrop_IsVariableCount(t *testing.T) {
	scheme := &FatigueDrop{
		TargetReps: 3, StartRPE: 8, StopRPE: 10, DropPercent: 0.05, MaxSets: 10,
	}
	if !scheme.IsVariableCount() {
		t.Error("FatigueDrop.IsVariableCount() should return true")
	}
}

func TestFatigueDrop_GetTerminationCondition(t *testing.T) {
	scheme := &FatigueDrop{
		TargetReps: 3, StartRPE: 8, StopRPE: 10, DropPercent: 0.05, MaxSets: 10,
	}
	cond := scheme.GetTerminationCondition()
	if cond == nil {
		t.Fatal("expected non-nil termination condition")
	}
	if cond.Type() != TerminationTypeRPEThreshold {
		t.Errorf("expected type %s, got %s", TerminationTypeRPEThreshold, cond.Type())
	}
	rpe, ok := cond.(*RPEThreshold)
	if !ok {
		t.Fatal("expected *RPEThreshold")
	}
	if rpe.Threshold != 10 {
		t.Errorf("expected threshold 10, got %v", rpe.Threshold)
	}
}

// === GenerateSets Tests ===

func TestFatigueDrop_GenerateSets(t *testing.T) {
	tests := []struct {
		name       string
		scheme     FatigueDrop
		baseWeight float64
		wantSets   int
		wantReps   int
	}{
		{
			name: "basic first set",
			scheme: FatigueDrop{
				TargetReps: 3, StartRPE: 8, StopRPE: 10, DropPercent: 0.05, MaxSets: 10,
			},
			baseWeight: 315.0,
			wantSets:   1, // Only first provisional set
			wantReps:   3,
		},
		{
			name: "different reps",
			scheme: FatigueDrop{
				TargetReps: 5, StartRPE: 7, StopRPE: 9, DropPercent: 0.03, MaxSets: 8,
			},
			baseWeight: 225.0,
			wantSets:   1,
			wantReps:   5,
		},
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

			set := sets[0]
			if set.SetNumber != 1 {
				t.Errorf("expected SetNumber 1, got %d", set.SetNumber)
			}
			if set.Weight != tt.baseWeight {
				t.Errorf("expected Weight %f, got %f", tt.baseWeight, set.Weight)
			}
			if set.TargetReps != tt.wantReps {
				t.Errorf("expected TargetReps %d, got %d", tt.wantReps, set.TargetReps)
			}
			if !set.IsWorkSet {
				t.Error("expected IsWorkSet to be true")
			}
			if !set.IsProvisional {
				t.Error("expected IsProvisional to be true")
			}
		})
	}
}

func TestFatigueDrop_GenerateSets_InvalidScheme(t *testing.T) {
	scheme := FatigueDrop{
		TargetReps: 0, StartRPE: 8, StopRPE: 10, DropPercent: 0.05, MaxSets: 10,
	}
	ctx := DefaultSetGenerationContext()
	sets, err := scheme.GenerateSets(315.0, ctx)
	if err == nil {
		t.Error("expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidParams) {
		t.Errorf("expected ErrInvalidParams, got %v", err)
	}
	if sets != nil {
		t.Error("expected nil sets on error")
	}
}

// === GenerateNextSet Tests ===

func TestFatigueDrop_GenerateNextSet(t *testing.T) {
	scheme := &FatigueDrop{
		TargetReps:  3,
		StartRPE:    8,
		StopRPE:     10,
		DropPercent: 0.05,
		MaxSets:     10,
	}
	ctx := DefaultSetGenerationContext()

	// Initial set at 315 lbs
	history := []GeneratedSet{
		{SetNumber: 1, Weight: 315.0, TargetReps: 3, IsWorkSet: true, IsProvisional: true},
	}

	// First continuation - RPE 8.5, should continue
	termCtx := TerminationContext{
		SetNumber:  1,
		TotalSets:  1,
		LastRPE:    floatPtrFD(8.5),
		LastReps:   3,
		TargetReps: 3,
	}

	nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
	if !shouldContinue {
		t.Error("expected shouldContinue=true")
	}
	if nextSet == nil {
		t.Fatal("expected non-nil next set")
	}
	if nextSet.SetNumber != 2 {
		t.Errorf("expected SetNumber 2, got %d", nextSet.SetNumber)
	}
	// 315 * 0.95 = 299.25, rounded down to 295
	if nextSet.Weight != 295.0 {
		t.Errorf("expected Weight 295.0, got %f", nextSet.Weight)
	}
	if nextSet.TargetReps != 3 {
		t.Errorf("expected TargetReps 3, got %d", nextSet.TargetReps)
	}
	if !nextSet.IsProvisional {
		t.Error("expected IsProvisional to be true")
	}
}

func TestFatigueDrop_GenerateNextSet_TerminatesAtRPE(t *testing.T) {
	scheme := &FatigueDrop{
		TargetReps:  3,
		StartRPE:    8,
		StopRPE:     10,
		DropPercent: 0.05,
		MaxSets:     10,
	}
	ctx := DefaultSetGenerationContext()

	history := []GeneratedSet{
		{SetNumber: 1, Weight: 315.0, TargetReps: 3, IsWorkSet: true},
	}

	// RPE hits 10 - should terminate
	termCtx := TerminationContext{
		SetNumber:  1,
		TotalSets:  1,
		LastRPE:    floatPtrFD(10),
		LastReps:   3,
		TargetReps: 3,
	}

	nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
	if shouldContinue {
		t.Error("expected shouldContinue=false when RPE hits threshold")
	}
	if nextSet != nil {
		t.Error("expected nil set when terminating")
	}
}

func TestFatigueDrop_GenerateNextSet_TerminatesAtMaxSets(t *testing.T) {
	scheme := &FatigueDrop{
		TargetReps:  3,
		StartRPE:    8,
		StopRPE:     10,
		DropPercent: 0.05,
		MaxSets:     5, // Low max sets for testing
	}
	ctx := DefaultSetGenerationContext()

	history := []GeneratedSet{
		{SetNumber: 5, Weight: 250.0, TargetReps: 3, IsWorkSet: true},
	}

	// RPE still below threshold but max sets reached
	termCtx := TerminationContext{
		SetNumber:  5,
		TotalSets:  5,
		LastRPE:    floatPtrFD(9),
		LastReps:   3,
		TargetReps: 3,
	}

	nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
	if shouldContinue {
		t.Error("expected shouldContinue=false when max sets reached")
	}
	if nextSet != nil {
		t.Error("expected nil set when terminating")
	}
}

func TestFatigueDrop_GenerateNextSet_DefaultMaxSets(t *testing.T) {
	scheme := &FatigueDrop{
		TargetReps:  3,
		StartRPE:    8,
		StopRPE:     10,
		DropPercent: 0.05,
		MaxSets:     0, // Uses default of 10
	}
	ctx := DefaultSetGenerationContext()

	history := []GeneratedSet{
		{SetNumber: 10, Weight: 200.0, TargetReps: 3, IsWorkSet: true},
	}

	termCtx := TerminationContext{
		SetNumber:  10,
		TotalSets:  10, // Default max reached
		LastRPE:    floatPtrFD(9),
		LastReps:   3,
		TargetReps: 3,
	}

	nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
	if shouldContinue {
		t.Error("expected shouldContinue=false when default max sets reached")
	}
	if nextSet != nil {
		t.Error("expected nil set when terminating")
	}
}

func TestFatigueDrop_GenerateNextSet_EmptyHistory(t *testing.T) {
	scheme := &FatigueDrop{
		TargetReps:  3,
		StartRPE:    8,
		StopRPE:     10,
		DropPercent: 0.05,
		MaxSets:     10,
	}
	ctx := DefaultSetGenerationContext()

	// Empty history - should return nil
	termCtx := TerminationContext{
		TotalSets: 0,
		LastRPE:   floatPtrFD(8),
	}

	nextSet, shouldContinue := scheme.GenerateNextSet(ctx, []GeneratedSet{}, termCtx)
	if shouldContinue {
		t.Error("expected shouldContinue=false with empty history")
	}
	if nextSet != nil {
		t.Error("expected nil set with empty history")
	}
}

func TestFatigueDrop_GenerateNextSet_ZeroDropPercent(t *testing.T) {
	scheme := &FatigueDrop{
		TargetReps:  3,
		StartRPE:    8,
		StopRPE:     10,
		DropPercent: 0, // No weight drop
		MaxSets:     10,
	}
	ctx := DefaultSetGenerationContext()

	history := []GeneratedSet{
		{SetNumber: 1, Weight: 315.0, TargetReps: 3, IsWorkSet: true},
	}

	termCtx := TerminationContext{
		SetNumber:  1,
		TotalSets:  1,
		LastRPE:    floatPtrFD(8.5),
		LastReps:   3,
		TargetReps: 3,
	}

	nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
	if !shouldContinue {
		t.Error("expected shouldContinue=true")
	}
	if nextSet == nil {
		t.Fatal("expected non-nil next set")
	}
	// With 0% drop, weight stays the same
	if nextSet.Weight != 315.0 {
		t.Errorf("expected Weight 315.0 (no drop), got %f", nextSet.Weight)
	}
}

// === Full Progression Scenario Tests ===

func TestFatigueDrop_FullProgression(t *testing.T) {
	// Simulate a full fatigue drop session:
	// Squat @ 3 reps, start at RPE 8, drop 5%, stop at RPE 10
	scheme := &FatigueDrop{
		TargetReps:  3,
		StartRPE:    8,
		StopRPE:     10,
		DropPercent: 0.05,
		MaxSets:     10,
	}
	ctx := DefaultSetGenerationContext()

	// Generate first set
	initialSets, err := scheme.GenerateSets(315.0, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(initialSets) != 1 {
		t.Fatalf("expected 1 initial set, got %d", len(initialSets))
	}

	history := initialSets

	// Simulate session with RPE progression
	rpeProgression := []float64{8.0, 8.5, 9.0, 9.5, 10.0}
	expectedWeights := []float64{315.0, 295.0, 280.0, 265.0, 250.0} // Rounded down

	for i, rpe := range rpeProgression {
		if i >= len(history) {
			t.Fatalf("ran out of history at iteration %d", i)
		}

		// Verify current set weight
		if i < len(expectedWeights) && history[i].Weight != expectedWeights[i] {
			t.Errorf("set %d: expected weight %f, got %f", i+1, expectedWeights[i], history[i].Weight)
		}

		termCtx := TerminationContext{
			SetNumber:  i + 1,
			TotalSets:  i + 1,
			LastRPE:    floatPtrFD(rpe),
			LastReps:   3,
			TargetReps: 3,
		}

		nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)

		if rpe >= 10 {
			// Should terminate
			if shouldContinue {
				t.Errorf("set %d: expected termination at RPE %v", i+1, rpe)
			}
			if nextSet != nil {
				t.Errorf("set %d: expected nil set at termination", i+1)
			}
			break
		} else {
			// Should continue
			if !shouldContinue {
				t.Errorf("set %d: expected continuation at RPE %v", i+1, rpe)
			}
			if nextSet == nil {
				t.Fatalf("set %d: expected non-nil next set", i+1)
			}
			history = append(history, *nextSet)
		}
	}

	// Verify we ended up with 5 sets
	if len(history) != 5 {
		t.Errorf("expected 5 sets total, got %d", len(history))
	}
}

func TestFatigueDrop_ImmediateTermination(t *testing.T) {
	// If first set hits stop RPE, should terminate immediately
	scheme := &FatigueDrop{
		TargetReps:  3,
		StartRPE:    8,
		StopRPE:     10,
		DropPercent: 0.05,
		MaxSets:     10,
	}
	ctx := DefaultSetGenerationContext()

	history := []GeneratedSet{
		{SetNumber: 1, Weight: 315.0, TargetReps: 3, IsWorkSet: true},
	}

	// First set is RPE 10
	termCtx := TerminationContext{
		SetNumber:  1,
		TotalSets:  1,
		LastRPE:    floatPtrFD(10),
		LastReps:   3,
		TargetReps: 3,
	}

	nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
	if shouldContinue {
		t.Error("expected immediate termination when first set is at stop RPE")
	}
	if nextSet != nil {
		t.Error("expected nil set on immediate termination")
	}
}

// === JSON Serialization Tests ===

func TestFatigueDrop_MarshalJSON(t *testing.T) {
	scheme := &FatigueDrop{
		TargetReps:  3,
		StartRPE:    8,
		StopRPE:     10,
		DropPercent: 0.05,
		MaxSets:     10,
	}

	data, err := json.Marshal(scheme)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if parsed["type"] != string(TypeFatigueDrop) {
		t.Errorf("expected type %s, got %v", TypeFatigueDrop, parsed["type"])
	}
	if int(parsed["target_reps"].(float64)) != 3 {
		t.Errorf("expected target_reps 3, got %v", parsed["target_reps"])
	}
	if parsed["start_rpe"].(float64) != 8 {
		t.Errorf("expected start_rpe 8, got %v", parsed["start_rpe"])
	}
	if parsed["stop_rpe"].(float64) != 10 {
		t.Errorf("expected stop_rpe 10, got %v", parsed["stop_rpe"])
	}
	if parsed["drop_percent"].(float64) != 0.05 {
		t.Errorf("expected drop_percent 0.05, got %v", parsed["drop_percent"])
	}
	if int(parsed["max_sets"].(float64)) != 10 {
		t.Errorf("expected max_sets 10, got %v", parsed["max_sets"])
	}
}

func TestUnmarshalFatigueDrop(t *testing.T) {
	tests := []struct {
		name        string
		json        string
		wantReps    int
		wantStartRPE float64
		wantStopRPE  float64
		wantErr     bool
	}{
		{
			name:        "valid",
			json:        `{"type": "FATIGUE_DROP", "target_reps": 3, "start_rpe": 8, "stop_rpe": 10, "drop_percent": 0.05, "max_sets": 10}`,
			wantReps:    3,
			wantStartRPE: 8,
			wantStopRPE:  10,
			wantErr:     false,
		},
		{
			name:        "valid without max_sets",
			json:        `{"type": "FATIGUE_DROP", "target_reps": 5, "start_rpe": 7, "stop_rpe": 9, "drop_percent": 0.03}`,
			wantReps:    5,
			wantStartRPE: 7,
			wantStopRPE:  9,
			wantErr:     false,
		},
		{
			name:    "invalid zero reps",
			json:    `{"type": "FATIGUE_DROP", "target_reps": 0, "start_rpe": 8, "stop_rpe": 10, "drop_percent": 0.05}`,
			wantErr: true,
		},
		{
			name:    "invalid stop <= start",
			json:    `{"type": "FATIGUE_DROP", "target_reps": 3, "start_rpe": 10, "stop_rpe": 8, "drop_percent": 0.05}`,
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
			scheme, err := UnmarshalFatigueDrop(json.RawMessage(tt.json))
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

				fd, ok := scheme.(*FatigueDrop)
				if !ok {
					t.Fatal("expected *FatigueDrop")
				}
				if fd.TargetReps != tt.wantReps {
					t.Errorf("expected TargetReps %d, got %d", tt.wantReps, fd.TargetReps)
				}
				if fd.StartRPE != tt.wantStartRPE {
					t.Errorf("expected StartRPE %v, got %v", tt.wantStartRPE, fd.StartRPE)
				}
				if fd.StopRPE != tt.wantStopRPE {
					t.Errorf("expected StopRPE %v, got %v", tt.wantStopRPE, fd.StopRPE)
				}
			}
		})
	}
}

func TestFatigueDrop_RoundTrip(t *testing.T) {
	original := &FatigueDrop{
		TargetReps:  3,
		StartRPE:    8,
		StopRPE:     10,
		DropPercent: 0.05,
		MaxSets:     10,
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	scheme, err := UnmarshalFatigueDrop(data)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	fd, ok := scheme.(*FatigueDrop)
	if !ok {
		t.Fatal("expected *FatigueDrop")
	}
	if fd.TargetReps != original.TargetReps {
		t.Errorf("expected TargetReps %d, got %d", original.TargetReps, fd.TargetReps)
	}
	if fd.StartRPE != original.StartRPE {
		t.Errorf("expected StartRPE %v, got %v", original.StartRPE, fd.StartRPE)
	}
	if fd.StopRPE != original.StopRPE {
		t.Errorf("expected StopRPE %v, got %v", original.StopRPE, fd.StopRPE)
	}
	if fd.DropPercent != original.DropPercent {
		t.Errorf("expected DropPercent %v, got %v", original.DropPercent, fd.DropPercent)
	}
	if fd.MaxSets != original.MaxSets {
		t.Errorf("expected MaxSets %d, got %d", original.MaxSets, fd.MaxSets)
	}
}

// === Factory Registration Tests ===

func TestRegisterFatigueDrop(t *testing.T) {
	factory := NewSchemeFactory()

	if factory.IsRegistered(TypeFatigueDrop) {
		t.Error("TypeFatigueDrop should not be registered initially")
	}

	RegisterFatigueDrop(factory)

	if !factory.IsRegistered(TypeFatigueDrop) {
		t.Error("TypeFatigueDrop should be registered after RegisterFatigueDrop")
	}
}

func TestFatigueDrop_FactoryIntegration(t *testing.T) {
	factory := NewSchemeFactory()
	RegisterFatigueDrop(factory)

	t.Run("Create from type and data", func(t *testing.T) {
		jsonData := json.RawMessage(`{"target_reps": 3, "start_rpe": 8, "stop_rpe": 10, "drop_percent": 0.05, "max_sets": 10}`)
		scheme, err := factory.Create(TypeFatigueDrop, jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if scheme.Type() != TypeFatigueDrop {
			t.Errorf("expected type %s, got %s", TypeFatigueDrop, scheme.Type())
		}
	})

	t.Run("CreateFromJSON", func(t *testing.T) {
		jsonData := json.RawMessage(`{"type": "FATIGUE_DROP", "target_reps": 3, "start_rpe": 8, "stop_rpe": 10, "drop_percent": 0.05, "max_sets": 10}`)
		scheme, err := factory.CreateFromJSON(jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if scheme.Type() != TypeFatigueDrop {
			t.Errorf("expected type %s, got %s", TypeFatigueDrop, scheme.Type())
		}

		// Generate sets
		ctx := DefaultSetGenerationContext()
		sets, err := scheme.GenerateSets(315.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sets) != 1 {
			t.Errorf("expected 1 set, got %d", len(sets))
		}
	})
}

// === Interface Compliance Tests ===

func TestFatigueDrop_Implements_SetScheme(t *testing.T) {
	var _ SetScheme = (*FatigueDrop)(nil)
}

func TestFatigueDrop_Implements_VariableSetScheme(t *testing.T) {
	var _ VariableSetScheme = (*FatigueDrop)(nil)
}

// === Edge Cases ===

func TestFatigueDrop_EdgeCases(t *testing.T) {
	t.Run("very small drop percent", func(t *testing.T) {
		scheme := &FatigueDrop{
			TargetReps:  3,
			StartRPE:    8,
			StopRPE:     10,
			DropPercent: 0.01, // 1% drop
			MaxSets:     10,
		}
		ctx := DefaultSetGenerationContext()

		history := []GeneratedSet{
			{SetNumber: 1, Weight: 315.0, TargetReps: 3, IsWorkSet: true},
		}

		termCtx := TerminationContext{
			TotalSets: 1,
			LastRPE:   floatPtrFD(8.5),
		}

		nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
		if !shouldContinue {
			t.Error("expected to continue")
		}
		// 315 * 0.99 = 311.85, rounded down to 310
		if nextSet.Weight != 310.0 {
			t.Errorf("expected Weight 310.0, got %f", nextSet.Weight)
		}
	})

	t.Run("large weight with standard drop", func(t *testing.T) {
		scheme := &FatigueDrop{
			TargetReps:  3,
			StartRPE:    8,
			StopRPE:     10,
			DropPercent: 0.05,
			MaxSets:     10,
		}
		ctx := DefaultSetGenerationContext()

		history := []GeneratedSet{
			{SetNumber: 1, Weight: 500.0, TargetReps: 3, IsWorkSet: true},
		}

		termCtx := TerminationContext{
			TotalSets: 1,
			LastRPE:   floatPtrFD(8.5),
		}

		nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
		if !shouldContinue {
			t.Error("expected to continue")
		}
		// 500 * 0.95 = 475.0
		if nextSet.Weight != 475.0 {
			t.Errorf("expected Weight 475.0, got %f", nextSet.Weight)
		}
	})

	t.Run("drop percent boundary at 1.0", func(t *testing.T) {
		// 100% drop - valid but results in 0 weight, which should terminate
		scheme := &FatigueDrop{
			TargetReps:  3,
			StartRPE:    8,
			StopRPE:     10,
			DropPercent: 1.0, // 100% drop
			MaxSets:     10,
		}
		ctx := DefaultSetGenerationContext()

		history := []GeneratedSet{
			{SetNumber: 1, Weight: 315.0, TargetReps: 3, IsWorkSet: true},
		}

		termCtx := TerminationContext{
			TotalSets: 1,
			LastRPE:   floatPtrFD(8.5),
		}

		nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
		// Weight would be 0, should terminate
		if shouldContinue {
			t.Error("expected termination when weight drops to 0")
		}
		if nextSet != nil {
			t.Error("expected nil set when weight drops to 0")
		}
	})

	t.Run("no RPE provided continues", func(t *testing.T) {
		scheme := &FatigueDrop{
			TargetReps:  3,
			StartRPE:    8,
			StopRPE:     10,
			DropPercent: 0.05,
			MaxSets:     10,
		}
		ctx := DefaultSetGenerationContext()

		history := []GeneratedSet{
			{SetNumber: 1, Weight: 315.0, TargetReps: 3, IsWorkSet: true},
		}

		// No RPE reported (nil)
		termCtx := TerminationContext{
			TotalSets:  1,
			LastRPE:    nil, // No RPE
			LastReps:   3,
			TargetReps: 3,
		}

		nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
		// When no RPE is provided, termination check returns false (continue)
		if !shouldContinue {
			t.Error("expected to continue when no RPE provided")
		}
		if nextSet == nil {
			t.Error("expected non-nil set when no RPE provided")
		}
	})
}

// === RTS Program Example Test ===

func TestFatigueDrop_RTSStyleExample(t *testing.T) {
	// Real-world RTS example:
	// Competition Squat @ 3 reps
	// Start at RPE 8, drop 5%, stop at RPE 10
	// Lifter works up to 315 at RPE 8, then does fatigue drops

	scheme, err := NewFatigueDrop(3, 8, 10, 0.05, 10)
	if err != nil {
		t.Fatalf("unexpected error creating scheme: %v", err)
	}

	ctx := DefaultSetGenerationContext()
	sets, err := scheme.GenerateSets(315.0, ctx)
	if err != nil {
		t.Fatalf("unexpected error generating sets: %v", err)
	}

	if len(sets) != 1 {
		t.Fatalf("expected 1 initial set, got %d", len(sets))
	}

	// Verify first set
	firstSet := sets[0]
	if firstSet.Weight != 315.0 {
		t.Errorf("first set: expected Weight 315.0, got %f", firstSet.Weight)
	}
	if firstSet.TargetReps != 3 {
		t.Errorf("first set: expected TargetReps 3, got %d", firstSet.TargetReps)
	}
	if !firstSet.IsWorkSet {
		t.Error("first set: expected IsWorkSet true")
	}
	if !firstSet.IsProvisional {
		t.Error("first set: expected IsProvisional true")
	}

	// Verify termination condition
	termCond := scheme.GetTerminationCondition()
	if termCond.Type() != TerminationTypeRPEThreshold {
		t.Errorf("expected termination type %s, got %s", TerminationTypeRPEThreshold, termCond.Type())
	}

	// Verify IsVariableCount
	if !scheme.IsVariableCount() {
		t.Error("expected IsVariableCount true")
	}
}

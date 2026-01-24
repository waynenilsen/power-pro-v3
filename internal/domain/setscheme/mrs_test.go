package setscheme

import (
	"encoding/json"
	"errors"
	"testing"
)

// === Type Tests ===

func TestMRS_Type(t *testing.T) {
	scheme := &MRS{
		TargetTotalReps: 25,
		MinRepsPerSet:   3,
		MaxSets:         10,
		NumberOfMRS:     3,
	}
	if scheme.Type() != TypeMRS {
		t.Errorf("expected type %s, got %s", TypeMRS, scheme.Type())
	}
}

// === Constructor Tests ===

func TestNewMRS(t *testing.T) {
	tests := []struct {
		name            string
		targetTotalReps int
		minRepsPerSet   int
		maxSets         int
		numberOfMRS     int
		wantErr         bool
	}{
		{"valid basic", 25, 3, 10, 3, false},
		{"valid T1 style", 15, 3, 10, 3, false},
		{"valid T3 style", 40, 5, 10, 4, false},
		{"valid min values", 1, 1, 1, 0, false},
		{"valid zero max sets (uses default)", 25, 3, 0, 3, false},
		{"valid zero number of MRS", 25, 3, 10, 0, false},
		{"valid target equals min", 5, 5, 10, 1, false},
		{"invalid zero target reps", 0, 3, 10, 3, true},
		{"invalid negative target reps", -1, 3, 10, 3, true},
		{"invalid zero min reps", 25, 0, 10, 3, true},
		{"invalid negative min reps", 25, -1, 10, 3, true},
		{"invalid target less than min", 3, 5, 10, 3, true},
		{"invalid negative max sets", 25, 3, -1, 3, true},
		{"invalid negative number of MRS", 25, 3, 10, -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme, err := NewMRS(tt.targetTotalReps, tt.minRepsPerSet, tt.maxSets, tt.numberOfMRS)
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
				if scheme.TargetTotalReps != tt.targetTotalReps {
					t.Errorf("expected TargetTotalReps %d, got %d", tt.targetTotalReps, scheme.TargetTotalReps)
				}
				if scheme.MinRepsPerSet != tt.minRepsPerSet {
					t.Errorf("expected MinRepsPerSet %d, got %d", tt.minRepsPerSet, scheme.MinRepsPerSet)
				}
				if scheme.MaxSets != tt.maxSets {
					t.Errorf("expected MaxSets %d, got %d", tt.maxSets, scheme.MaxSets)
				}
				if scheme.NumberOfMRS != tt.numberOfMRS {
					t.Errorf("expected NumberOfMRS %d, got %d", tt.numberOfMRS, scheme.NumberOfMRS)
				}
			}
		})
	}
}

// === Validate Tests ===

func TestMRS_Validate(t *testing.T) {
	tests := []struct {
		name    string
		scheme  MRS
		wantErr bool
	}{
		{
			name: "valid",
			scheme: MRS{
				TargetTotalReps: 25, MinRepsPerSet: 3, MaxSets: 10, NumberOfMRS: 3,
			},
			wantErr: false,
		},
		{
			name: "valid zero max sets",
			scheme: MRS{
				TargetTotalReps: 25, MinRepsPerSet: 3, MaxSets: 0, NumberOfMRS: 3,
			},
			wantErr: false,
		},
		{
			name: "valid target equals min",
			scheme: MRS{
				TargetTotalReps: 5, MinRepsPerSet: 5, MaxSets: 10, NumberOfMRS: 1,
			},
			wantErr: false,
		},
		{
			name: "invalid zero target reps",
			scheme: MRS{
				TargetTotalReps: 0, MinRepsPerSet: 3, MaxSets: 10, NumberOfMRS: 3,
			},
			wantErr: true,
		},
		{
			name: "invalid zero min reps",
			scheme: MRS{
				TargetTotalReps: 25, MinRepsPerSet: 0, MaxSets: 10, NumberOfMRS: 3,
			},
			wantErr: true,
		},
		{
			name: "invalid target less than min",
			scheme: MRS{
				TargetTotalReps: 3, MinRepsPerSet: 5, MaxSets: 10, NumberOfMRS: 3,
			},
			wantErr: true,
		},
		{
			name: "invalid negative max sets",
			scheme: MRS{
				TargetTotalReps: 25, MinRepsPerSet: 3, MaxSets: -1, NumberOfMRS: 3,
			},
			wantErr: true,
		},
		{
			name: "invalid negative number of MRS",
			scheme: MRS{
				TargetTotalReps: 25, MinRepsPerSet: 3, MaxSets: 10, NumberOfMRS: -1,
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

func TestMRS_IsVariableCount(t *testing.T) {
	scheme := &MRS{
		TargetTotalReps: 25, MinRepsPerSet: 3, MaxSets: 10, NumberOfMRS: 3,
	}
	if !scheme.IsVariableCount() {
		t.Error("MRS.IsVariableCount() should return true")
	}
}

func TestMRS_GetTerminationCondition(t *testing.T) {
	scheme := &MRS{
		TargetTotalReps: 25, MinRepsPerSet: 3, MaxSets: 10, NumberOfMRS: 3,
	}
	cond := scheme.GetTerminationCondition()
	if cond == nil {
		t.Fatal("expected non-nil termination condition")
	}
	if cond.Type() != TerminationTypeTotalReps {
		t.Errorf("expected type %s, got %s", TerminationTypeTotalReps, cond.Type())
	}
	totalReps, ok := cond.(*TotalReps)
	if !ok {
		t.Fatal("expected *TotalReps")
	}
	if totalReps.Target != 25 {
		t.Errorf("expected target 25, got %d", totalReps.Target)
	}
}

// === GenerateSets Tests ===

func TestMRS_GenerateSets(t *testing.T) {
	tests := []struct {
		name       string
		scheme     MRS
		baseWeight float64
		wantSets   int
		wantReps   int
	}{
		{
			name: "basic first set",
			scheme: MRS{
				TargetTotalReps: 25, MinRepsPerSet: 3, MaxSets: 10, NumberOfMRS: 3,
			},
			baseWeight: 225.0,
			wantSets:   1,
			wantReps:   3, // MinRepsPerSet
		},
		{
			name: "different min reps",
			scheme: MRS{
				TargetTotalReps: 40, MinRepsPerSet: 5, MaxSets: 8, NumberOfMRS: 4,
			},
			baseWeight: 135.0,
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

func TestMRS_GenerateSets_InvalidScheme(t *testing.T) {
	scheme := MRS{
		TargetTotalReps: 0, MinRepsPerSet: 3, MaxSets: 10, NumberOfMRS: 3,
	}
	ctx := DefaultSetGenerationContext()
	sets, err := scheme.GenerateSets(225.0, ctx)
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

func TestMRS_GenerateNextSet(t *testing.T) {
	scheme := &MRS{
		TargetTotalReps: 25,
		MinRepsPerSet:   3,
		MaxSets:         10,
		NumberOfMRS:     3,
	}
	ctx := DefaultSetGenerationContext()

	// Initial set at 225 lbs
	history := []GeneratedSet{
		{SetNumber: 1, Weight: 225.0, TargetReps: 3, IsWorkSet: true, IsProvisional: true},
	}

	// First continuation - did 10 reps, total 10, should continue
	termCtx := TerminationContext{
		SetNumber:  1,
		TotalSets:  1,
		LastReps:   10,
		TotalReps:  10,
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
	// MRS uses same weight for all sets
	if nextSet.Weight != 225.0 {
		t.Errorf("expected Weight 225.0 (same as first set), got %f", nextSet.Weight)
	}
	if nextSet.TargetReps != 3 {
		t.Errorf("expected TargetReps 3, got %d", nextSet.TargetReps)
	}
	if !nextSet.IsProvisional {
		t.Error("expected IsProvisional to be true")
	}
}

func TestMRS_GenerateNextSet_TerminatesAtTotalReps(t *testing.T) {
	scheme := &MRS{
		TargetTotalReps: 25,
		MinRepsPerSet:   3,
		MaxSets:         10,
		NumberOfMRS:     3,
	}
	ctx := DefaultSetGenerationContext()

	history := []GeneratedSet{
		{SetNumber: 1, Weight: 225.0, TargetReps: 3, IsWorkSet: true},
		{SetNumber: 2, Weight: 225.0, TargetReps: 3, IsWorkSet: true},
		{SetNumber: 3, Weight: 225.0, TargetReps: 3, IsWorkSet: true},
		{SetNumber: 4, Weight: 225.0, TargetReps: 3, IsWorkSet: true},
	}

	// Total reps hits 27 (>= 25) - should terminate
	termCtx := TerminationContext{
		SetNumber:  4,
		TotalSets:  4,
		LastReps:   4,
		TotalReps:  27, // 10 + 8 + 5 + 4 = 27
		TargetReps: 3,
	}

	nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
	if shouldContinue {
		t.Error("expected shouldContinue=false when total reps hits target")
	}
	if nextSet != nil {
		t.Error("expected nil set when terminating")
	}
}

func TestMRS_GenerateNextSet_TerminatesAtRepFailure(t *testing.T) {
	scheme := &MRS{
		TargetTotalReps: 25,
		MinRepsPerSet:   3,
		MaxSets:         10,
		NumberOfMRS:     3,
	}
	ctx := DefaultSetGenerationContext()

	history := []GeneratedSet{
		{SetNumber: 1, Weight: 225.0, TargetReps: 3, IsWorkSet: true},
		{SetNumber: 2, Weight: 225.0, TargetReps: 3, IsWorkSet: true},
		{SetNumber: 3, Weight: 225.0, TargetReps: 3, IsWorkSet: true},
	}

	// Last set only got 2 reps (< 3 minimum) - should terminate due to failure
	termCtx := TerminationContext{
		SetNumber:  3,
		TotalSets:  3,
		LastReps:   2, // Failed to hit minimum
		TotalReps:  18,
		TargetReps: 3,
	}

	nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
	if shouldContinue {
		t.Error("expected shouldContinue=false when reps < minimum")
	}
	if nextSet != nil {
		t.Error("expected nil set when terminating due to failure")
	}
}

func TestMRS_GenerateNextSet_TerminatesAtMaxSets(t *testing.T) {
	scheme := &MRS{
		TargetTotalReps: 100, // Very high target
		MinRepsPerSet:   3,
		MaxSets:         5, // Low max sets for testing
		NumberOfMRS:     3,
	}
	ctx := DefaultSetGenerationContext()

	history := []GeneratedSet{
		{SetNumber: 5, Weight: 225.0, TargetReps: 3, IsWorkSet: true},
	}

	// Total reps still below target but max sets reached
	termCtx := TerminationContext{
		SetNumber:  5,
		TotalSets:  5,
		LastReps:   5,
		TotalReps:  30, // Still below 100
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

func TestMRS_GenerateNextSet_DefaultMaxSets(t *testing.T) {
	scheme := &MRS{
		TargetTotalReps: 200, // Very high target
		MinRepsPerSet:   3,
		MaxSets:         0, // Uses default of 10
		NumberOfMRS:     3,
	}
	ctx := DefaultSetGenerationContext()

	history := []GeneratedSet{
		{SetNumber: 10, Weight: 225.0, TargetReps: 3, IsWorkSet: true},
	}

	termCtx := TerminationContext{
		SetNumber:  10,
		TotalSets:  10, // Default max reached
		LastReps:   5,
		TotalReps:  50,
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

func TestMRS_GenerateNextSet_EmptyHistory(t *testing.T) {
	scheme := &MRS{
		TargetTotalReps: 25,
		MinRepsPerSet:   3,
		MaxSets:         10,
		NumberOfMRS:     3,
	}
	ctx := DefaultSetGenerationContext()

	// Empty history - should return nil
	termCtx := TerminationContext{
		TotalSets: 0,
		LastReps:  10,
	}

	nextSet, shouldContinue := scheme.GenerateNextSet(ctx, []GeneratedSet{}, termCtx)
	if shouldContinue {
		t.Error("expected shouldContinue=false with empty history")
	}
	if nextSet != nil {
		t.Error("expected nil set with empty history")
	}
}

// === Full Progression Scenario Tests ===

func TestMRS_FullProgression_Success(t *testing.T) {
	// Simulate a full MRS session reaching target:
	// Bench Press MRS x 3, target 25 total reps
	scheme := &MRS{
		TargetTotalReps: 25,
		MinRepsPerSet:   3,
		MaxSets:         10,
		NumberOfMRS:     3,
	}
	ctx := DefaultSetGenerationContext()

	// Generate first set
	initialSets, err := scheme.GenerateSets(225.0, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(initialSets) != 1 {
		t.Fatalf("expected 1 initial set, got %d", len(initialSets))
	}

	history := initialSets

	// Simulate session: 10, 8, 5, 4 reps (total 27)
	repsPerSet := []int{10, 8, 5, 4}
	cumulativeReps := 0

	for i, reps := range repsPerSet {
		if i >= len(history) {
			t.Fatalf("ran out of history at iteration %d", i)
		}

		cumulativeReps += reps

		// Verify all sets have same weight
		if history[i].Weight != 225.0 {
			t.Errorf("set %d: expected weight 225.0, got %f", i+1, history[i].Weight)
		}

		termCtx := TerminationContext{
			SetNumber:  i + 1,
			TotalSets:  i + 1,
			LastReps:   reps,
			TotalReps:  cumulativeReps,
			TargetReps: 3,
		}

		nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)

		if cumulativeReps >= 25 {
			// Should terminate
			if shouldContinue {
				t.Errorf("set %d: expected termination at total reps %d", i+1, cumulativeReps)
			}
			if nextSet != nil {
				t.Errorf("set %d: expected nil set at termination", i+1)
			}
			break
		} else {
			// Should continue
			if !shouldContinue {
				t.Errorf("set %d: expected continuation at total reps %d", i+1, cumulativeReps)
			}
			if nextSet == nil {
				t.Fatalf("set %d: expected non-nil next set", i+1)
			}
			history = append(history, *nextSet)
		}
	}

	// Verify we ended up with 4 sets
	if len(history) != 4 {
		t.Errorf("expected 4 sets total, got %d", len(history))
	}
}

func TestMRS_FullProgression_Failure(t *testing.T) {
	// Simulate a failure scenario:
	// Last set fails to hit minimum reps
	scheme := &MRS{
		TargetTotalReps: 25,
		MinRepsPerSet:   3,
		MaxSets:         10,
		NumberOfMRS:     3,
	}
	ctx := DefaultSetGenerationContext()

	initialSets, err := scheme.GenerateSets(225.0, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	history := initialSets

	// Simulate session: 10, 6, 2 (failure on last set)
	repsPerSet := []int{10, 6, 2}
	cumulativeReps := 0

	for i, reps := range repsPerSet {
		if i >= len(history) {
			t.Fatalf("ran out of history at iteration %d", i)
		}

		cumulativeReps += reps

		termCtx := TerminationContext{
			SetNumber:  i + 1,
			TotalSets:  i + 1,
			LastReps:   reps,
			TotalReps:  cumulativeReps,
			TargetReps: 3,
		}

		nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)

		if reps < 3 {
			// Should terminate due to failure
			if shouldContinue {
				t.Errorf("set %d: expected termination at rep failure (got %d reps)", i+1, reps)
			}
			if nextSet != nil {
				t.Errorf("set %d: expected nil set at termination", i+1)
			}
			break
		} else {
			if !shouldContinue {
				t.Errorf("set %d: expected continuation at %d reps", i+1, reps)
			}
			if nextSet == nil {
				t.Fatalf("set %d: expected non-nil next set", i+1)
			}
			history = append(history, *nextSet)
		}
	}

	// Verify we ended up with 3 sets (stopped at failure)
	if len(history) != 3 {
		t.Errorf("expected 3 sets total, got %d", len(history))
	}
}

func TestMRS_ImmediateTermination(t *testing.T) {
	// If first set hits total reps target, should terminate immediately
	scheme := &MRS{
		TargetTotalReps: 10,
		MinRepsPerSet:   3,
		MaxSets:         10,
		NumberOfMRS:     3,
	}
	ctx := DefaultSetGenerationContext()

	history := []GeneratedSet{
		{SetNumber: 1, Weight: 225.0, TargetReps: 3, IsWorkSet: true},
	}

	// First set hit 12 reps (>= 10 target)
	termCtx := TerminationContext{
		SetNumber:  1,
		TotalSets:  1,
		LastReps:   12,
		TotalReps:  12,
		TargetReps: 3,
	}

	nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
	if shouldContinue {
		t.Error("expected immediate termination when first set exceeds target")
	}
	if nextSet != nil {
		t.Error("expected nil set on immediate termination")
	}
}

// === JSON Serialization Tests ===

func TestMRS_MarshalJSON(t *testing.T) {
	scheme := &MRS{
		TargetTotalReps: 25,
		MinRepsPerSet:   3,
		MaxSets:         10,
		NumberOfMRS:     3,
	}

	data, err := json.Marshal(scheme)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if parsed["type"] != string(TypeMRS) {
		t.Errorf("expected type %s, got %v", TypeMRS, parsed["type"])
	}
	if int(parsed["target_total_reps"].(float64)) != 25 {
		t.Errorf("expected target_total_reps 25, got %v", parsed["target_total_reps"])
	}
	if int(parsed["min_reps_per_set"].(float64)) != 3 {
		t.Errorf("expected min_reps_per_set 3, got %v", parsed["min_reps_per_set"])
	}
	if int(parsed["max_sets"].(float64)) != 10 {
		t.Errorf("expected max_sets 10, got %v", parsed["max_sets"])
	}
	if int(parsed["number_of_mrs"].(float64)) != 3 {
		t.Errorf("expected number_of_mrs 3, got %v", parsed["number_of_mrs"])
	}
}

func TestUnmarshalMRS(t *testing.T) {
	tests := []struct {
		name            string
		json            string
		wantTargetReps  int
		wantMinReps     int
		wantErr         bool
	}{
		{
			name:           "valid",
			json:           `{"type": "MRS", "target_total_reps": 25, "min_reps_per_set": 3, "max_sets": 10, "number_of_mrs": 3}`,
			wantTargetReps: 25,
			wantMinReps:    3,
			wantErr:        false,
		},
		{
			name:           "valid without optional fields",
			json:           `{"type": "MRS", "target_total_reps": 15, "min_reps_per_set": 5}`,
			wantTargetReps: 15,
			wantMinReps:    5,
			wantErr:        false,
		},
		{
			name:    "invalid zero target reps",
			json:    `{"type": "MRS", "target_total_reps": 0, "min_reps_per_set": 3}`,
			wantErr: true,
		},
		{
			name:    "invalid target less than min",
			json:    `{"type": "MRS", "target_total_reps": 3, "min_reps_per_set": 5}`,
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
			scheme, err := UnmarshalMRS(json.RawMessage(tt.json))
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

				mrs, ok := scheme.(*MRS)
				if !ok {
					t.Fatal("expected *MRS")
				}
				if mrs.TargetTotalReps != tt.wantTargetReps {
					t.Errorf("expected TargetTotalReps %d, got %d", tt.wantTargetReps, mrs.TargetTotalReps)
				}
				if mrs.MinRepsPerSet != tt.wantMinReps {
					t.Errorf("expected MinRepsPerSet %d, got %d", tt.wantMinReps, mrs.MinRepsPerSet)
				}
			}
		})
	}
}

func TestMRS_RoundTrip(t *testing.T) {
	original := &MRS{
		TargetTotalReps: 25,
		MinRepsPerSet:   3,
		MaxSets:         10,
		NumberOfMRS:     3,
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	scheme, err := UnmarshalMRS(data)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	mrs, ok := scheme.(*MRS)
	if !ok {
		t.Fatal("expected *MRS")
	}
	if mrs.TargetTotalReps != original.TargetTotalReps {
		t.Errorf("expected TargetTotalReps %d, got %d", original.TargetTotalReps, mrs.TargetTotalReps)
	}
	if mrs.MinRepsPerSet != original.MinRepsPerSet {
		t.Errorf("expected MinRepsPerSet %d, got %d", original.MinRepsPerSet, mrs.MinRepsPerSet)
	}
	if mrs.MaxSets != original.MaxSets {
		t.Errorf("expected MaxSets %d, got %d", original.MaxSets, mrs.MaxSets)
	}
	if mrs.NumberOfMRS != original.NumberOfMRS {
		t.Errorf("expected NumberOfMRS %d, got %d", original.NumberOfMRS, mrs.NumberOfMRS)
	}
}

// === Factory Registration Tests ===

func TestRegisterMRS(t *testing.T) {
	factory := NewSchemeFactory()

	if factory.IsRegistered(TypeMRS) {
		t.Error("TypeMRS should not be registered initially")
	}

	RegisterMRS(factory)

	if !factory.IsRegistered(TypeMRS) {
		t.Error("TypeMRS should be registered after RegisterMRS")
	}
}

func TestMRS_FactoryIntegration(t *testing.T) {
	factory := NewSchemeFactory()
	RegisterMRS(factory)

	t.Run("Create from type and data", func(t *testing.T) {
		jsonData := json.RawMessage(`{"target_total_reps": 25, "min_reps_per_set": 3, "max_sets": 10, "number_of_mrs": 3}`)
		scheme, err := factory.Create(TypeMRS, jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if scheme.Type() != TypeMRS {
			t.Errorf("expected type %s, got %s", TypeMRS, scheme.Type())
		}
	})

	t.Run("CreateFromJSON", func(t *testing.T) {
		jsonData := json.RawMessage(`{"type": "MRS", "target_total_reps": 25, "min_reps_per_set": 3, "max_sets": 10, "number_of_mrs": 3}`)
		scheme, err := factory.CreateFromJSON(jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if scheme.Type() != TypeMRS {
			t.Errorf("expected type %s, got %s", TypeMRS, scheme.Type())
		}

		// Generate sets
		ctx := DefaultSetGenerationContext()
		sets, err := scheme.GenerateSets(225.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sets) != 1 {
			t.Errorf("expected 1 set, got %d", len(sets))
		}
	})
}

// === Interface Compliance Tests ===

func TestMRS_Implements_SetScheme(t *testing.T) {
	var _ SetScheme = (*MRS)(nil)
}

func TestMRS_Implements_VariableSetScheme(t *testing.T) {
	var _ VariableSetScheme = (*MRS)(nil)
}

// === Edge Cases ===

func TestMRS_EdgeCases(t *testing.T) {
	t.Run("single rep min", func(t *testing.T) {
		scheme := &MRS{
			TargetTotalReps: 10,
			MinRepsPerSet:   1,
			MaxSets:         10,
			NumberOfMRS:     1,
		}
		ctx := DefaultSetGenerationContext()

		history := []GeneratedSet{
			{SetNumber: 1, Weight: 315.0, TargetReps: 1, IsWorkSet: true},
		}

		// Did 1 rep (exactly min), should continue
		termCtx := TerminationContext{
			TotalSets:  1,
			LastReps:   1,
			TotalReps:  1,
			TargetReps: 1,
		}

		nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
		if !shouldContinue {
			t.Error("expected to continue when hitting exactly min reps")
		}
		if nextSet == nil {
			t.Error("expected non-nil next set")
		}
	})

	t.Run("target equals min reps", func(t *testing.T) {
		scheme := &MRS{
			TargetTotalReps: 5,
			MinRepsPerSet:   5,
			MaxSets:         10,
			NumberOfMRS:     1,
		}
		ctx := DefaultSetGenerationContext()

		history := []GeneratedSet{
			{SetNumber: 1, Weight: 225.0, TargetReps: 5, IsWorkSet: true},
		}

		// First set did exactly 5 reps (meets target)
		termCtx := TerminationContext{
			TotalSets:  1,
			LastReps:   5,
			TotalReps:  5,
			TargetReps: 5,
		}

		nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
		if shouldContinue {
			t.Error("expected termination when target met")
		}
		if nextSet != nil {
			t.Error("expected nil set when terminating")
		}
	})

	t.Run("large rep counts", func(t *testing.T) {
		scheme := &MRS{
			TargetTotalReps: 100,
			MinRepsPerSet:   10,
			MaxSets:         20,
			NumberOfMRS:     5,
		}
		ctx := DefaultSetGenerationContext()

		history := []GeneratedSet{
			{SetNumber: 1, Weight: 135.0, TargetReps: 10, IsWorkSet: true},
		}

		// Did 25 reps first set
		termCtx := TerminationContext{
			TotalSets:  1,
			LastReps:   25,
			TotalReps:  25,
			TargetReps: 10,
		}

		nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
		if !shouldContinue {
			t.Error("expected to continue")
		}
		if nextSet.Weight != 135.0 {
			t.Errorf("expected Weight 135.0, got %f", nextSet.Weight)
		}
	})

	t.Run("boundary at max sets - 1", func(t *testing.T) {
		scheme := &MRS{
			TargetTotalReps: 100,
			MinRepsPerSet:   3,
			MaxSets:         5,
			NumberOfMRS:     3,
		}
		ctx := DefaultSetGenerationContext()

		history := []GeneratedSet{
			{SetNumber: 4, Weight: 225.0, TargetReps: 3, IsWorkSet: true},
		}

		// At 4 sets, still below max of 5
		termCtx := TerminationContext{
			TotalSets:  4,
			LastReps:   5,
			TotalReps:  20,
			TargetReps: 3,
		}

		nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
		if !shouldContinue {
			t.Error("expected to continue at max-1 sets")
		}
		if nextSet == nil {
			t.Error("expected non-nil next set")
		}
		if nextSet.SetNumber != 5 {
			t.Errorf("expected SetNumber 5, got %d", nextSet.SetNumber)
		}
	})
}

// === GZCL Program Example Tests ===

func TestMRS_GZCLStyle_T1(t *testing.T) {
	// GZCL T1: 3 MRS, target ~15 total reps, minimum 3 reps per set
	scheme, err := NewMRS(15, 3, 10, 3)
	if err != nil {
		t.Fatalf("unexpected error creating scheme: %v", err)
	}

	ctx := DefaultSetGenerationContext()
	sets, err := scheme.GenerateSets(225.0, ctx)
	if err != nil {
		t.Fatalf("unexpected error generating sets: %v", err)
	}

	if len(sets) != 1 {
		t.Fatalf("expected 1 initial set, got %d", len(sets))
	}

	firstSet := sets[0]
	if firstSet.Weight != 225.0 {
		t.Errorf("first set: expected Weight 225.0, got %f", firstSet.Weight)
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
	if termCond.Type() != TerminationTypeTotalReps {
		t.Errorf("expected termination type %s, got %s", TerminationTypeTotalReps, termCond.Type())
	}

	// Verify NumberOfMRS field
	if scheme.NumberOfMRS != 3 {
		t.Errorf("expected NumberOfMRS 3, got %d", scheme.NumberOfMRS)
	}
}

func TestMRS_GZCLStyle_T3(t *testing.T) {
	// GZCL T3: 4 MRS, target ~40 total reps, minimum 5 reps per set
	scheme, err := NewMRS(40, 5, 10, 4)
	if err != nil {
		t.Fatalf("unexpected error creating scheme: %v", err)
	}

	if scheme.NumberOfMRS != 4 {
		t.Errorf("expected NumberOfMRS 4 for T3, got %d", scheme.NumberOfMRS)
	}
	if scheme.MinRepsPerSet != 5 {
		t.Errorf("expected MinRepsPerSet 5 for T3, got %d", scheme.MinRepsPerSet)
	}
}

// === TotalReps Termination Condition Tests ===

func TestTotalReps_Type(t *testing.T) {
	cond := &TotalReps{Target: 25}
	if cond.Type() != TerminationTypeTotalReps {
		t.Errorf("expected type %s, got %s", TerminationTypeTotalReps, cond.Type())
	}
}

func TestNewTotalReps(t *testing.T) {
	tests := []struct {
		name    string
		target  int
		wantErr bool
	}{
		{"valid 25", 25, false},
		{"valid 1", 1, false},
		{"valid 100", 100, false},
		{"invalid 0", 0, true},
		{"invalid -1", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond, err := NewTotalReps(tt.target)
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
				if cond.Target != tt.target {
					t.Errorf("expected target %d, got %d", tt.target, cond.Target)
				}
			}
		})
	}
}

func TestTotalReps_ShouldTerminate(t *testing.T) {
	cond := &TotalReps{Target: 25}

	tests := []struct {
		name     string
		ctx      TerminationContext
		expected bool
	}{
		{
			name:     "below target",
			ctx:      TerminationContext{TotalReps: 20},
			expected: false,
		},
		{
			name:     "at target",
			ctx:      TerminationContext{TotalReps: 25},
			expected: true,
		},
		{
			name:     "above target",
			ctx:      TerminationContext{TotalReps: 30},
			expected: true,
		},
		{
			name:     "zero reps",
			ctx:      TerminationContext{TotalReps: 0},
			expected: false,
		},
		{
			name:     "one below target",
			ctx:      TerminationContext{TotalReps: 24},
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

func TestTotalReps_Validate(t *testing.T) {
	tests := []struct {
		name    string
		target  int
		wantErr bool
	}{
		{"valid 1", 1, false},
		{"valid 25", 25, false},
		{"valid 100", 100, false},
		{"invalid 0", 0, true},
		{"invalid -1", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond := &TotalReps{Target: tt.target}
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

func TestTotalReps_MarshalJSON(t *testing.T) {
	cond := &TotalReps{Target: 25}
	data, err := json.Marshal(cond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if parsed["type"] != string(TerminationTypeTotalReps) {
		t.Errorf("expected type %s, got %v", TerminationTypeTotalReps, parsed["type"])
	}
	if int(parsed["target"].(float64)) != 25 {
		t.Errorf("expected target 25, got %v", parsed["target"])
	}
}

func TestTotalReps_RoundTrip(t *testing.T) {
	original := &TotalReps{Target: 30}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	cond, err := UnmarshalTerminationCondition(data)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	totalReps, ok := cond.(*TotalReps)
	if !ok {
		t.Fatal("expected *TotalReps")
	}
	if totalReps.Target != original.Target {
		t.Errorf("expected target %d, got %d", original.Target, totalReps.Target)
	}
}

func TestTotalReps_InterfaceCompliance(t *testing.T) {
	var _ TerminationCondition = (*TotalReps)(nil)
}

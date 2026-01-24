package setscheme

import (
	"encoding/json"
	"errors"
	"testing"
)

// === Type Tests ===

func TestTotalRepsScheme_Type(t *testing.T) {
	scheme := &TotalRepsScheme{
		TargetTotalReps:     100,
		SuggestedRepsPerSet: 10,
		MaxSets:             20,
	}
	if scheme.Type() != TypeTotalReps {
		t.Errorf("expected type %s, got %s", TypeTotalReps, scheme.Type())
	}
}

// === Constructor Tests ===

func TestNewTotalRepsScheme(t *testing.T) {
	tests := []struct {
		name                string
		targetTotalReps     int
		suggestedRepsPerSet int
		maxSets             int
		wantErr             bool
	}{
		{"valid basic", 100, 10, 20, false},
		{"valid chin-ups BTM style", 100, 15, 20, false},
		{"valid dips range lower bound", 100, 20, 15, false},
		{"valid min values", 1, 0, 0, false},
		{"valid zero max sets (uses default)", 100, 10, 0, false},
		{"valid zero suggested reps (uses default)", 100, 0, 20, false},
		{"invalid zero target reps", 0, 10, 20, true},
		{"invalid negative target reps", -1, 10, 20, true},
		{"invalid negative suggested reps", 100, -1, 20, true},
		{"invalid negative max sets", 100, 10, -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme, err := NewTotalRepsScheme(tt.targetTotalReps, tt.suggestedRepsPerSet, tt.maxSets)
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
				if scheme.SuggestedRepsPerSet != tt.suggestedRepsPerSet {
					t.Errorf("expected SuggestedRepsPerSet %d, got %d", tt.suggestedRepsPerSet, scheme.SuggestedRepsPerSet)
				}
				if scheme.MaxSets != tt.maxSets {
					t.Errorf("expected MaxSets %d, got %d", tt.maxSets, scheme.MaxSets)
				}
			}
		})
	}
}

// === Validate Tests ===

func TestTotalRepsScheme_Validate(t *testing.T) {
	tests := []struct {
		name    string
		scheme  TotalRepsScheme
		wantErr bool
	}{
		{
			name: "valid",
			scheme: TotalRepsScheme{
				TargetTotalReps: 100, SuggestedRepsPerSet: 10, MaxSets: 20,
			},
			wantErr: false,
		},
		{
			name: "valid zero max sets",
			scheme: TotalRepsScheme{
				TargetTotalReps: 100, SuggestedRepsPerSet: 10, MaxSets: 0,
			},
			wantErr: false,
		},
		{
			name: "valid zero suggested reps",
			scheme: TotalRepsScheme{
				TargetTotalReps: 100, SuggestedRepsPerSet: 0, MaxSets: 20,
			},
			wantErr: false,
		},
		{
			name: "valid min target",
			scheme: TotalRepsScheme{
				TargetTotalReps: 1, SuggestedRepsPerSet: 1, MaxSets: 1,
			},
			wantErr: false,
		},
		{
			name: "invalid zero target reps",
			scheme: TotalRepsScheme{
				TargetTotalReps: 0, SuggestedRepsPerSet: 10, MaxSets: 20,
			},
			wantErr: true,
		},
		{
			name: "invalid negative suggested reps",
			scheme: TotalRepsScheme{
				TargetTotalReps: 100, SuggestedRepsPerSet: -1, MaxSets: 20,
			},
			wantErr: true,
		},
		{
			name: "invalid negative max sets",
			scheme: TotalRepsScheme{
				TargetTotalReps: 100, SuggestedRepsPerSet: 10, MaxSets: -1,
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

func TestTotalRepsScheme_IsVariableCount(t *testing.T) {
	scheme := &TotalRepsScheme{
		TargetTotalReps: 100, SuggestedRepsPerSet: 10, MaxSets: 20,
	}
	if !scheme.IsVariableCount() {
		t.Error("TotalRepsScheme.IsVariableCount() should return true")
	}
}

func TestTotalRepsScheme_GetTerminationCondition(t *testing.T) {
	scheme := &TotalRepsScheme{
		TargetTotalReps: 100, SuggestedRepsPerSet: 10, MaxSets: 20,
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
	if totalReps.Target != 100 {
		t.Errorf("expected target 100, got %d", totalReps.Target)
	}
}

// === GenerateSets Tests ===

func TestTotalRepsScheme_GenerateSets(t *testing.T) {
	tests := []struct {
		name       string
		scheme     TotalRepsScheme
		baseWeight float64
		wantSets   int
		wantReps   int
	}{
		{
			name: "basic first set with suggested reps",
			scheme: TotalRepsScheme{
				TargetTotalReps: 100, SuggestedRepsPerSet: 15, MaxSets: 20,
			},
			baseWeight: 0.0, // Bodyweight
			wantSets:   1,
			wantReps:   15,
		},
		{
			name: "first set without suggested reps (uses default 10)",
			scheme: TotalRepsScheme{
				TargetTotalReps: 100, SuggestedRepsPerSet: 0, MaxSets: 20,
			},
			baseWeight: 0.0,
			wantSets:   1,
			wantReps:   10, // Default
		},
		{
			name: "with weight",
			scheme: TotalRepsScheme{
				TargetTotalReps: 50, SuggestedRepsPerSet: 8, MaxSets: 15,
			},
			baseWeight: 25.0, // Weighted dips
			wantSets:   1,
			wantReps:   8,
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

func TestTotalRepsScheme_GenerateSets_InvalidScheme(t *testing.T) {
	scheme := TotalRepsScheme{
		TargetTotalReps: 0, SuggestedRepsPerSet: 10, MaxSets: 20,
	}
	ctx := DefaultSetGenerationContext()
	sets, err := scheme.GenerateSets(0.0, ctx)
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

func TestTotalRepsScheme_GenerateNextSet(t *testing.T) {
	scheme := &TotalRepsScheme{
		TargetTotalReps:     100,
		SuggestedRepsPerSet: 15,
		MaxSets:             20,
	}
	ctx := DefaultSetGenerationContext()

	// Initial set (bodyweight chin-ups)
	history := []GeneratedSet{
		{SetNumber: 1, Weight: 0.0, TargetReps: 15, IsWorkSet: true, IsProvisional: true},
	}

	// First continuation - did 15 reps, total 15, should continue
	termCtx := TerminationContext{
		SetNumber:  1,
		TotalSets:  1,
		LastReps:   15,
		TotalReps:  15,
		TargetReps: 15,
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
	// TotalReps uses same weight for all sets
	if nextSet.Weight != 0.0 {
		t.Errorf("expected Weight 0.0 (same as first set), got %f", nextSet.Weight)
	}
	if nextSet.TargetReps != 15 {
		t.Errorf("expected TargetReps 15, got %d", nextSet.TargetReps)
	}
	if !nextSet.IsProvisional {
		t.Error("expected IsProvisional to be true")
	}
}

func TestTotalRepsScheme_GenerateNextSet_TerminatesAtTotalReps(t *testing.T) {
	scheme := &TotalRepsScheme{
		TargetTotalReps:     100,
		SuggestedRepsPerSet: 15,
		MaxSets:             20,
	}
	ctx := DefaultSetGenerationContext()

	history := []GeneratedSet{
		{SetNumber: 1, Weight: 0.0, TargetReps: 15, IsWorkSet: true},
		{SetNumber: 2, Weight: 0.0, TargetReps: 15, IsWorkSet: true},
		{SetNumber: 3, Weight: 0.0, TargetReps: 15, IsWorkSet: true},
		{SetNumber: 4, Weight: 0.0, TargetReps: 15, IsWorkSet: true},
		{SetNumber: 5, Weight: 0.0, TargetReps: 15, IsWorkSet: true},
		{SetNumber: 6, Weight: 0.0, TargetReps: 15, IsWorkSet: true},
		{SetNumber: 7, Weight: 0.0, TargetReps: 15, IsWorkSet: true},
	}

	// Total reps hits 105 (>= 100) - should terminate
	termCtx := TerminationContext{
		SetNumber:  7,
		TotalSets:  7,
		LastReps:   15,
		TotalReps:  105, // 15 * 7 = 105
		TargetReps: 15,
	}

	nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
	if shouldContinue {
		t.Error("expected shouldContinue=false when total reps hits target")
	}
	if nextSet != nil {
		t.Error("expected nil set when terminating")
	}
}

func TestTotalRepsScheme_GenerateNextSet_NoRepFailureCondition(t *testing.T) {
	// Unlike MRS, TotalReps has NO rep failure condition
	// Users can do any number of reps per set
	scheme := &TotalRepsScheme{
		TargetTotalReps:     100,
		SuggestedRepsPerSet: 15,
		MaxSets:             20,
	}
	ctx := DefaultSetGenerationContext()

	history := []GeneratedSet{
		{SetNumber: 1, Weight: 0.0, TargetReps: 15, IsWorkSet: true},
		{SetNumber: 2, Weight: 0.0, TargetReps: 15, IsWorkSet: true},
		{SetNumber: 3, Weight: 0.0, TargetReps: 15, IsWorkSet: true},
	}

	// Last set only got 3 reps - unlike MRS, this should NOT terminate
	termCtx := TerminationContext{
		SetNumber:  3,
		TotalSets:  3,
		LastReps:   3, // Very low reps
		TotalReps:  33,
		TargetReps: 15,
	}

	nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
	if !shouldContinue {
		t.Error("expected shouldContinue=true - TotalReps has no rep failure condition")
	}
	if nextSet == nil {
		t.Error("expected non-nil next set")
	}
}

func TestTotalRepsScheme_GenerateNextSet_TerminatesAtMaxSets(t *testing.T) {
	scheme := &TotalRepsScheme{
		TargetTotalReps:     1000, // Very high target
		SuggestedRepsPerSet: 10,
		MaxSets:             5, // Low max sets for testing
	}
	ctx := DefaultSetGenerationContext()

	history := []GeneratedSet{
		{SetNumber: 5, Weight: 0.0, TargetReps: 10, IsWorkSet: true},
	}

	// Total reps still below target but max sets reached
	termCtx := TerminationContext{
		SetNumber:  5,
		TotalSets:  5,
		LastReps:   10,
		TotalReps:  50, // Still below 1000
		TargetReps: 10,
	}

	nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
	if shouldContinue {
		t.Error("expected shouldContinue=false when max sets reached")
	}
	if nextSet != nil {
		t.Error("expected nil set when terminating")
	}
}

func TestTotalRepsScheme_GenerateNextSet_DefaultMaxSets(t *testing.T) {
	scheme := &TotalRepsScheme{
		TargetTotalReps:     2000, // Very high target
		SuggestedRepsPerSet: 10,
		MaxSets:             0, // Uses default of 20
	}
	ctx := DefaultSetGenerationContext()

	history := []GeneratedSet{
		{SetNumber: 20, Weight: 0.0, TargetReps: 10, IsWorkSet: true},
	}

	termCtx := TerminationContext{
		SetNumber:  20,
		TotalSets:  20, // Default max reached
		LastReps:   10,
		TotalReps:  200,
		TargetReps: 10,
	}

	nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
	if shouldContinue {
		t.Error("expected shouldContinue=false when default max sets reached")
	}
	if nextSet != nil {
		t.Error("expected nil set when terminating")
	}
}

func TestTotalRepsScheme_GenerateNextSet_EmptyHistory(t *testing.T) {
	scheme := &TotalRepsScheme{
		TargetTotalReps:     100,
		SuggestedRepsPerSet: 15,
		MaxSets:             20,
	}
	ctx := DefaultSetGenerationContext()

	// Empty history - should return nil
	termCtx := TerminationContext{
		TotalSets: 0,
		LastReps:  15,
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

func TestTotalRepsScheme_FullProgression_ChinUps(t *testing.T) {
	// Simulate 100 chin-ups from Building the Monolith
	scheme := &TotalRepsScheme{
		TargetTotalReps:     100,
		SuggestedRepsPerSet: 15,
		MaxSets:             20,
	}
	ctx := DefaultSetGenerationContext()

	// Generate first set
	initialSets, err := scheme.GenerateSets(0.0, ctx) // Bodyweight
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(initialSets) != 1 {
		t.Fatalf("expected 1 initial set, got %d", len(initialSets))
	}

	history := initialSets

	// Simulate session: 15, 12, 12, 10, 10, 10, 10, 8, 8, 5 (total 100)
	repsPerSet := []int{15, 12, 12, 10, 10, 10, 10, 8, 8, 5}
	cumulativeReps := 0

	for i, reps := range repsPerSet {
		if i >= len(history) {
			t.Fatalf("ran out of history at iteration %d", i)
		}

		// Verify all sets have same weight (bodyweight = 0)
		if history[i].Weight != 0.0 {
			t.Errorf("set %d: expected weight 0.0, got %f", i+1, history[i].Weight)
		}

		cumulativeReps += reps

		termCtx := TerminationContext{
			SetNumber:  i + 1,
			TotalSets:  i + 1,
			LastReps:   reps,
			TotalReps:  cumulativeReps,
			TargetReps: 15,
		}

		nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)

		if cumulativeReps >= 100 {
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

	// Verify we ended up with 10 sets
	if len(history) != 10 {
		t.Errorf("expected 10 sets total, got %d", len(history))
	}
}

func TestTotalRepsScheme_FullProgression_SmallSets(t *testing.T) {
	// Simulate a scenario where user does many small sets
	scheme := &TotalRepsScheme{
		TargetTotalReps:     50,
		SuggestedRepsPerSet: 10,
		MaxSets:             20,
	}
	ctx := DefaultSetGenerationContext()

	initialSets, err := scheme.GenerateSets(0.0, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	history := initialSets

	// Do sets of 5 reps each - should take 10 sets to reach 50
	cumulativeReps := 0
	setCount := 0

	for cumulativeReps < 50 && setCount < 15 {
		reps := 5 // Small consistent sets
		cumulativeReps += reps
		setCount++

		termCtx := TerminationContext{
			SetNumber:  setCount,
			TotalSets:  setCount,
			LastReps:   reps,
			TotalReps:  cumulativeReps,
			TargetReps: 10,
		}

		nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)

		if cumulativeReps >= 50 {
			if shouldContinue {
				t.Error("expected termination when target reached")
			}
			break
		}

		if !shouldContinue {
			t.Errorf("unexpected termination at %d total reps", cumulativeReps)
		}
		if nextSet == nil {
			t.Fatal("expected non-nil next set")
		}
		history = append(history, *nextSet)
	}

	if cumulativeReps < 50 {
		t.Errorf("expected to reach 50 reps, got %d", cumulativeReps)
	}
}

func TestTotalRepsScheme_ImmediateTermination(t *testing.T) {
	// If first set hits total reps target, should terminate immediately
	scheme := &TotalRepsScheme{
		TargetTotalReps:     10,
		SuggestedRepsPerSet: 15,
		MaxSets:             20,
	}
	ctx := DefaultSetGenerationContext()

	history := []GeneratedSet{
		{SetNumber: 1, Weight: 0.0, TargetReps: 15, IsWorkSet: true},
	}

	// First set hit 12 reps (>= 10 target)
	termCtx := TerminationContext{
		SetNumber:  1,
		TotalSets:  1,
		LastReps:   12,
		TotalReps:  12,
		TargetReps: 15,
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

func TestTotalRepsScheme_MarshalJSON(t *testing.T) {
	scheme := &TotalRepsScheme{
		TargetTotalReps:     100,
		SuggestedRepsPerSet: 10,
		MaxSets:             20,
	}

	data, err := json.Marshal(scheme)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if parsed["type"] != string(TypeTotalReps) {
		t.Errorf("expected type %s, got %v", TypeTotalReps, parsed["type"])
	}
	if int(parsed["target_total_reps"].(float64)) != 100 {
		t.Errorf("expected target_total_reps 100, got %v", parsed["target_total_reps"])
	}
	if int(parsed["suggested_reps_per_set"].(float64)) != 10 {
		t.Errorf("expected suggested_reps_per_set 10, got %v", parsed["suggested_reps_per_set"])
	}
	if int(parsed["max_sets"].(float64)) != 20 {
		t.Errorf("expected max_sets 20, got %v", parsed["max_sets"])
	}
}

func TestTotalRepsScheme_MarshalJSON_OmitsZeroOptionalFields(t *testing.T) {
	scheme := &TotalRepsScheme{
		TargetTotalReps: 100,
		// SuggestedRepsPerSet and MaxSets are 0 (optional)
	}

	data, err := json.Marshal(scheme)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	// Zero optional fields should be omitted
	if _, exists := parsed["suggested_reps_per_set"]; exists {
		t.Error("expected suggested_reps_per_set to be omitted when 0")
	}
	if _, exists := parsed["max_sets"]; exists {
		t.Error("expected max_sets to be omitted when 0")
	}
}

func TestUnmarshalTotalRepsScheme(t *testing.T) {
	tests := []struct {
		name            string
		json            string
		wantTargetReps  int
		wantSuggested   int
		wantErr         bool
	}{
		{
			name:           "valid",
			json:           `{"type": "TOTAL_REPS", "target_total_reps": 100, "suggested_reps_per_set": 10, "max_sets": 20}`,
			wantTargetReps: 100,
			wantSuggested:  10,
			wantErr:        false,
		},
		{
			name:           "valid without optional fields",
			json:           `{"type": "TOTAL_REPS", "target_total_reps": 100}`,
			wantTargetReps: 100,
			wantSuggested:  0,
			wantErr:        false,
		},
		{
			name:    "invalid zero target reps",
			json:    `{"type": "TOTAL_REPS", "target_total_reps": 0}`,
			wantErr: true,
		},
		{
			name:    "invalid negative suggested reps",
			json:    `{"type": "TOTAL_REPS", "target_total_reps": 100, "suggested_reps_per_set": -1}`,
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
			scheme, err := UnmarshalTotalRepsScheme(json.RawMessage(tt.json))
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

				totalReps, ok := scheme.(*TotalRepsScheme)
				if !ok {
					t.Fatal("expected *TotalRepsScheme")
				}
				if totalReps.TargetTotalReps != tt.wantTargetReps {
					t.Errorf("expected TargetTotalReps %d, got %d", tt.wantTargetReps, totalReps.TargetTotalReps)
				}
				if totalReps.SuggestedRepsPerSet != tt.wantSuggested {
					t.Errorf("expected SuggestedRepsPerSet %d, got %d", tt.wantSuggested, totalReps.SuggestedRepsPerSet)
				}
			}
		})
	}
}

func TestTotalRepsScheme_RoundTrip(t *testing.T) {
	original := &TotalRepsScheme{
		TargetTotalReps:     100,
		SuggestedRepsPerSet: 15,
		MaxSets:             20,
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	scheme, err := UnmarshalTotalRepsScheme(data)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	totalReps, ok := scheme.(*TotalRepsScheme)
	if !ok {
		t.Fatal("expected *TotalRepsScheme")
	}
	if totalReps.TargetTotalReps != original.TargetTotalReps {
		t.Errorf("expected TargetTotalReps %d, got %d", original.TargetTotalReps, totalReps.TargetTotalReps)
	}
	if totalReps.SuggestedRepsPerSet != original.SuggestedRepsPerSet {
		t.Errorf("expected SuggestedRepsPerSet %d, got %d", original.SuggestedRepsPerSet, totalReps.SuggestedRepsPerSet)
	}
	if totalReps.MaxSets != original.MaxSets {
		t.Errorf("expected MaxSets %d, got %d", original.MaxSets, totalReps.MaxSets)
	}
}

// === Factory Registration Tests ===

func TestRegisterTotalRepsScheme(t *testing.T) {
	factory := NewSchemeFactory()

	if factory.IsRegistered(TypeTotalReps) {
		t.Error("TypeTotalReps should not be registered initially")
	}

	RegisterTotalRepsScheme(factory)

	if !factory.IsRegistered(TypeTotalReps) {
		t.Error("TypeTotalReps should be registered after RegisterTotalRepsScheme")
	}
}

func TestTotalRepsScheme_FactoryIntegration(t *testing.T) {
	factory := NewSchemeFactory()
	RegisterTotalRepsScheme(factory)

	t.Run("Create from type and data", func(t *testing.T) {
		jsonData := json.RawMessage(`{"target_total_reps": 100, "suggested_reps_per_set": 10, "max_sets": 20}`)
		scheme, err := factory.Create(TypeTotalReps, jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if scheme.Type() != TypeTotalReps {
			t.Errorf("expected type %s, got %s", TypeTotalReps, scheme.Type())
		}
	})

	t.Run("CreateFromJSON", func(t *testing.T) {
		jsonData := json.RawMessage(`{"type": "TOTAL_REPS", "target_total_reps": 100, "suggested_reps_per_set": 15, "max_sets": 20}`)
		scheme, err := factory.CreateFromJSON(jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if scheme.Type() != TypeTotalReps {
			t.Errorf("expected type %s, got %s", TypeTotalReps, scheme.Type())
		}

		// Generate sets
		ctx := DefaultSetGenerationContext()
		sets, err := scheme.GenerateSets(0.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sets) != 1 {
			t.Errorf("expected 1 set, got %d", len(sets))
		}
	})
}

// === Interface Compliance Tests ===

func TestTotalRepsScheme_Implements_SetScheme(t *testing.T) {
	var _ SetScheme = (*TotalRepsScheme)(nil)
}

func TestTotalRepsScheme_Implements_VariableSetScheme(t *testing.T) {
	var _ VariableSetScheme = (*TotalRepsScheme)(nil)
}

// === Edge Cases ===

func TestTotalRepsScheme_EdgeCases(t *testing.T) {
	t.Run("single rep target", func(t *testing.T) {
		scheme := &TotalRepsScheme{
			TargetTotalReps:     1,
			SuggestedRepsPerSet: 1,
			MaxSets:             5,
		}
		ctx := DefaultSetGenerationContext()

		history := []GeneratedSet{
			{SetNumber: 1, Weight: 0.0, TargetReps: 1, IsWorkSet: true},
		}

		// Did 1 rep - should terminate immediately
		termCtx := TerminationContext{
			TotalSets:  1,
			LastReps:   1,
			TotalReps:  1,
			TargetReps: 1,
		}

		nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
		if shouldContinue {
			t.Error("expected termination at single rep target")
		}
		if nextSet != nil {
			t.Error("expected nil set")
		}
	})

	t.Run("zero reps in set - should continue", func(t *testing.T) {
		scheme := &TotalRepsScheme{
			TargetTotalReps:     50,
			SuggestedRepsPerSet: 10,
			MaxSets:             20,
		}
		ctx := DefaultSetGenerationContext()

		history := []GeneratedSet{
			{SetNumber: 1, Weight: 0.0, TargetReps: 10, IsWorkSet: true},
		}

		// Even 0 reps should not terminate (unlike MRS)
		termCtx := TerminationContext{
			TotalSets:  1,
			LastReps:   0,
			TotalReps:  0,
			TargetReps: 10,
		}

		nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
		if !shouldContinue {
			t.Error("expected to continue even with 0 reps")
		}
		if nextSet == nil {
			t.Error("expected non-nil next set")
		}
	})

	t.Run("large rep counts", func(t *testing.T) {
		scheme := &TotalRepsScheme{
			TargetTotalReps:     500,
			SuggestedRepsPerSet: 25,
			MaxSets:             30,
		}
		ctx := DefaultSetGenerationContext()

		history := []GeneratedSet{
			{SetNumber: 1, Weight: 0.0, TargetReps: 25, IsWorkSet: true},
		}

		// Did 50 reps first set (big set!)
		termCtx := TerminationContext{
			TotalSets:  1,
			LastReps:   50,
			TotalReps:  50,
			TargetReps: 25,
		}

		nextSet, shouldContinue := scheme.GenerateNextSet(ctx, history, termCtx)
		if !shouldContinue {
			t.Error("expected to continue")
		}
		if nextSet.Weight != 0.0 {
			t.Errorf("expected Weight 0.0, got %f", nextSet.Weight)
		}
	})

	t.Run("boundary at max sets - 1", func(t *testing.T) {
		scheme := &TotalRepsScheme{
			TargetTotalReps:     1000,
			SuggestedRepsPerSet: 10,
			MaxSets:             5,
		}
		ctx := DefaultSetGenerationContext()

		history := []GeneratedSet{
			{SetNumber: 4, Weight: 0.0, TargetReps: 10, IsWorkSet: true},
		}

		// At 4 sets, still below max of 5
		termCtx := TerminationContext{
			TotalSets:  4,
			LastReps:   10,
			TotalReps:  40,
			TargetReps: 10,
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

	t.Run("weighted exercise", func(t *testing.T) {
		scheme := &TotalRepsScheme{
			TargetTotalReps:     50,
			SuggestedRepsPerSet: 8,
			MaxSets:             15,
		}
		ctx := DefaultSetGenerationContext()

		// Weighted dips at 45 lbs
		sets, err := scheme.GenerateSets(45.0, ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if sets[0].Weight != 45.0 {
			t.Errorf("expected Weight 45.0, got %f", sets[0].Weight)
		}

		history := sets
		termCtx := TerminationContext{
			TotalSets:  1,
			LastReps:   8,
			TotalReps:  8,
			TargetReps: 8,
		}

		nextSet, _ := scheme.GenerateNextSet(ctx, history, termCtx)
		if nextSet.Weight != 45.0 {
			t.Errorf("expected Weight 45.0 for all sets, got %f", nextSet.Weight)
		}
	})
}

// === Program Example Tests ===

func TestTotalRepsScheme_BTM_ChinUps(t *testing.T) {
	// 5/3/1 Building the Monolith: 100 chin-ups
	scheme, err := NewTotalRepsScheme(100, 15, 20)
	if err != nil {
		t.Fatalf("unexpected error creating scheme: %v", err)
	}

	ctx := DefaultSetGenerationContext()
	sets, err := scheme.GenerateSets(0.0, ctx) // Bodyweight
	if err != nil {
		t.Fatalf("unexpected error generating sets: %v", err)
	}

	if len(sets) != 1 {
		t.Fatalf("expected 1 initial set, got %d", len(sets))
	}

	firstSet := sets[0]
	if firstSet.Weight != 0.0 {
		t.Errorf("first set: expected Weight 0.0 (bodyweight), got %f", firstSet.Weight)
	}
	if firstSet.TargetReps != 15 {
		t.Errorf("first set: expected TargetReps 15, got %d", firstSet.TargetReps)
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
}

func TestTotalRepsScheme_BTM_Dips(t *testing.T) {
	// 5/3/1 Building the Monolith: 100-200 dips
	// Using lower bound (100)
	scheme, err := NewTotalRepsScheme(100, 20, 20)
	if err != nil {
		t.Fatalf("unexpected error creating scheme: %v", err)
	}

	if scheme.TargetTotalReps != 100 {
		t.Errorf("expected TargetTotalReps 100 for dips, got %d", scheme.TargetTotalReps)
	}
	if scheme.SuggestedRepsPerSet != 20 {
		t.Errorf("expected SuggestedRepsPerSet 20 for dips, got %d", scheme.SuggestedRepsPerSet)
	}
}

func TestTotalRepsScheme_HighVolume(t *testing.T) {
	// Extreme volume work: 200 face pulls
	scheme, err := NewTotalRepsScheme(200, 25, 20)
	if err != nil {
		t.Fatalf("unexpected error creating scheme: %v", err)
	}

	ctx := DefaultSetGenerationContext()
	sets, err := scheme.GenerateSets(0.0, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// First set should have suggested reps
	if sets[0].TargetReps != 25 {
		t.Errorf("expected TargetReps 25, got %d", sets[0].TargetReps)
	}
}

// === Comparison with MRS ===

func TestTotalRepsScheme_DiffersFromMRS_NoRepFailure(t *testing.T) {
	// This test documents the key difference: TotalReps has no rep failure

	// Create both schemes with same target
	totalRepsScheme := &TotalRepsScheme{
		TargetTotalReps:     50,
		SuggestedRepsPerSet: 10,
		MaxSets:             20,
	}

	mrsScheme := &MRS{
		TargetTotalReps: 50,
		MinRepsPerSet:   5,
		MaxSets:         10,
	}

	ctx := DefaultSetGenerationContext()

	// Both started at 20 reps total
	history := []GeneratedSet{
		{SetNumber: 1, Weight: 100.0, TargetReps: 10, IsWorkSet: true},
		{SetNumber: 2, Weight: 100.0, TargetReps: 10, IsWorkSet: true},
	}

	// User got only 2 reps on last set (well below min)
	termCtx := TerminationContext{
		SetNumber:  2,
		TotalSets:  2,
		LastReps:   2, // Very low reps
		TotalReps:  22,
		TargetReps: 10,
	}

	// TotalReps should CONTINUE
	trNextSet, trShouldContinue := totalRepsScheme.GenerateNextSet(ctx, history, termCtx)
	if !trShouldContinue {
		t.Error("TotalReps should continue even with 2 reps")
	}
	if trNextSet == nil {
		t.Error("TotalReps should return next set")
	}

	// MRS should TERMINATE (due to rep failure)
	mrsNextSet, mrsShouldContinue := mrsScheme.GenerateNextSet(ctx, history, termCtx)
	if mrsShouldContinue {
		t.Error("MRS should terminate with 2 reps < 5 min")
	}
	if mrsNextSet != nil {
		t.Error("MRS should return nil set on termination")
	}
}

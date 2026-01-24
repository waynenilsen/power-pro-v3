// Package integration provides integration tests for cross-component behavior.
// This file tests the StageProgression domain logic.
package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/progression"
	"github.com/waynenilsen/power-pro-v3/internal/domain/setscheme"
)

// TestStageProgressionCreate tests creating stage progressions.
func TestStageProgressionCreate(t *testing.T) {
	t.Run("GZCLP T1 default creation", func(t *testing.T) {
		prog, err := progression.NewGZCLPT1DefaultProgression(uuid.New().String(), "GZCLP T1")
		if err != nil {
			t.Fatalf("Failed to create GZCLP T1: %v", err)
		}

		if prog.Type() != progression.TypeStage {
			t.Errorf("Expected type %s, got %s", progression.TypeStage, prog.Type())
		}
		if prog.TriggerType() != progression.TriggerOnFailure {
			t.Errorf("Expected trigger %s, got %s", progression.TriggerOnFailure, prog.TriggerType())
		}
		if prog.StageCount() != 3 {
			t.Errorf("Expected 3 stages, got %d", prog.StageCount())
		}
		if prog.CurrentStage != 0 {
			t.Errorf("Expected initial stage 0, got %d", prog.CurrentStage)
		}
		if !prog.ResetOnExhaustion {
			t.Error("Expected ResetOnExhaustion to be true")
		}
		if !prog.DeloadOnReset {
			t.Error("Expected DeloadOnReset to be true")
		}
		if prog.DeloadPercent != 0.15 {
			t.Errorf("Expected 15%% deload, got %.2f%%", prog.DeloadPercent*100)
		}
	})

	t.Run("GZCLP T2 default creation", func(t *testing.T) {
		prog, err := progression.NewGZCLPT2DefaultProgression(uuid.New().String(), "GZCLP T2")
		if err != nil {
			t.Fatalf("Failed to create GZCLP T2: %v", err)
		}

		if prog.StageCount() != 3 {
			t.Errorf("Expected 3 stages, got %d", prog.StageCount())
		}
		if prog.DeloadOnReset {
			t.Error("Expected DeloadOnReset to be false for T2")
		}

		// Verify stages have no AMRAP
		for i, stage := range prog.Stages {
			if stage.IsAMRAP {
				t.Errorf("Stage %d should not be AMRAP for T2", i)
			}
		}
	})

	t.Run("requires at least 2 stages", func(t *testing.T) {
		_, err := progression.NewStageProgression(
			uuid.New().String(),
			"Single Stage",
			[]progression.Stage{
				{Name: "Only", Sets: 3, Reps: 5, IsAMRAP: false, MinVolume: 15},
			},
			true, false, 0, progression.TrainingMax,
		)
		if err == nil {
			t.Error("Expected error for single stage")
		}
	})
}

// TestStageProgressionAdvance tests stage advancement on failure.
func TestStageProgressionAdvance(t *testing.T) {
	ctx := context.Background()

	t.Run("failure advances to next stage", func(t *testing.T) {
		prog, _ := progression.NewGZCLPT1DefaultProgression(uuid.New().String(), "Test")
		currentWeight := 200.0

		// At stage 0, trigger failure
		consecutiveFailures := 1
		params := progression.ProgressionContext{
			UserID:       uuid.New().String(),
			LiftID:       uuid.New().String(),
			CurrentValue: currentWeight,
			MaxType:      progression.TrainingMax,
			TriggerEvent: progression.TriggerEvent{
				Type:                progression.TriggerOnFailure,
				Timestamp:           time.Now(),
				ConsecutiveFailures: &consecutiveFailures,
			},
		}

		result, err := prog.Apply(ctx, params)
		if err != nil {
			t.Fatalf("Apply failed: %v", err)
		}

		if !result.Applied {
			t.Errorf("Expected progression to apply, got Applied=false, Reason: %s", result.Reason)
		}
		if result.Delta != 0 {
			t.Errorf("Expected delta=0 (no weight change on stage advance), got %.1f", result.Delta)
		}
		if prog.CurrentStage != 1 {
			t.Errorf("Expected stage to advance to 1, got %d", prog.CurrentStage)
		}
	})

	t.Run("multiple failures advance through all stages", func(t *testing.T) {
		prog, _ := progression.NewGZCLPT1DefaultProgression(uuid.New().String(), "Test")
		currentWeight := 200.0
		consecutiveFailures := 1

		params := progression.ProgressionContext{
			UserID:       uuid.New().String(),
			LiftID:       uuid.New().String(),
			CurrentValue: currentWeight,
			MaxType:      progression.TrainingMax,
			TriggerEvent: progression.TriggerEvent{
				Type:                progression.TriggerOnFailure,
				Timestamp:           time.Now(),
				ConsecutiveFailures: &consecutiveFailures,
			},
		}

		// First failure: 0 -> 1
		prog.Apply(ctx, params)
		if prog.CurrentStage != 1 {
			t.Errorf("After first failure, expected stage 1, got %d", prog.CurrentStage)
		}

		// Second failure: 1 -> 2
		prog.Apply(ctx, params)
		if prog.CurrentStage != 2 {
			t.Errorf("After second failure, expected stage 2, got %d", prog.CurrentStage)
		}

		// Third failure: 2 -> reset to 0 with deload
		result, _ := prog.Apply(ctx, params)
		if prog.CurrentStage != 0 {
			t.Errorf("After exhaustion, expected stage 0, got %d", prog.CurrentStage)
		}
		if result.Delta >= 0 {
			t.Errorf("Expected negative delta (deload), got %.1f", result.Delta)
		}
	})
}

// TestStageProgressionReset tests reset behavior at stage exhaustion.
func TestStageProgressionReset(t *testing.T) {
	ctx := context.Background()

	t.Run("reset with deload reduces weight", func(t *testing.T) {
		prog, _ := progression.NewStageProgression(
			uuid.New().String(),
			"Test",
			[]progression.Stage{
				{Name: "A", Sets: 3, Reps: 5, IsAMRAP: false, MinVolume: 15},
				{Name: "B", Sets: 3, Reps: 3, IsAMRAP: false, MinVolume: 9},
			},
			true,  // ResetOnExhaustion
			true,  // DeloadOnReset
			0.10,  // 10% deload
			progression.TrainingMax,
		)

		// Move to last stage
		prog.SetCurrentStage(1)

		currentWeight := 200.0
		consecutiveFailures := 1
		params := progression.ProgressionContext{
			UserID:       uuid.New().String(),
			LiftID:       uuid.New().String(),
			CurrentValue: currentWeight,
			MaxType:      progression.TrainingMax,
			TriggerEvent: progression.TriggerEvent{
				Type:                progression.TriggerOnFailure,
				Timestamp:           time.Now(),
				ConsecutiveFailures: &consecutiveFailures,
			},
		}

		result, _ := prog.Apply(ctx, params)

		expectedDeload := currentWeight * 0.10
		expectedNewWeight := currentWeight - expectedDeload

		if result.Delta != -expectedDeload {
			t.Errorf("Expected delta=%.1f, got %.1f", -expectedDeload, result.Delta)
		}
		if result.NewValue != expectedNewWeight {
			t.Errorf("Expected new weight=%.1f, got %.1f", expectedNewWeight, result.NewValue)
		}
		if prog.CurrentStage != 0 {
			t.Errorf("Expected reset to stage 0, got %d", prog.CurrentStage)
		}
	})

	t.Run("reset without deload keeps weight", func(t *testing.T) {
		prog, _ := progression.NewStageProgression(
			uuid.New().String(),
			"Test",
			[]progression.Stage{
				{Name: "A", Sets: 3, Reps: 10, IsAMRAP: false, MinVolume: 30},
				{Name: "B", Sets: 3, Reps: 8, IsAMRAP: false, MinVolume: 24},
			},
			true,  // ResetOnExhaustion
			false, // DeloadOnReset = false
			0,
			progression.TrainingMax,
		)

		prog.SetCurrentStage(1)

		currentWeight := 100.0
		consecutiveFailures := 1
		params := progression.ProgressionContext{
			UserID:       uuid.New().String(),
			LiftID:       uuid.New().String(),
			CurrentValue: currentWeight,
			MaxType:      progression.TrainingMax,
			TriggerEvent: progression.TriggerEvent{
				Type:                progression.TriggerOnFailure,
				Timestamp:           time.Now(),
				ConsecutiveFailures: &consecutiveFailures,
			},
		}

		result, _ := prog.Apply(ctx, params)

		if result.Delta != 0 {
			t.Errorf("Expected delta=0 (no deload), got %.1f", result.Delta)
		}
		if result.NewValue != currentWeight {
			t.Errorf("Expected weight unchanged=%.1f, got %.1f", currentWeight, result.NewValue)
		}
		if prog.CurrentStage != 0 {
			t.Errorf("Expected reset to stage 0, got %d", prog.CurrentStage)
		}
	})

	t.Run("no reset on exhaustion returns not applied", func(t *testing.T) {
		prog, _ := progression.NewStageProgression(
			uuid.New().String(),
			"Test",
			[]progression.Stage{
				{Name: "A", Sets: 3, Reps: 10, IsAMRAP: false, MinVolume: 30},
				{Name: "B", Sets: 3, Reps: 8, IsAMRAP: false, MinVolume: 24},
			},
			false, // ResetOnExhaustion = false
			false,
			0,
			progression.TrainingMax,
		)

		prog.SetCurrentStage(1)

		currentWeight := 100.0
		consecutiveFailures := 1
		params := progression.ProgressionContext{
			UserID:       uuid.New().String(),
			LiftID:       uuid.New().String(),
			CurrentValue: currentWeight,
			MaxType:      progression.TrainingMax,
			TriggerEvent: progression.TriggerEvent{
				Type:                progression.TriggerOnFailure,
				Timestamp:           time.Now(),
				ConsecutiveFailures: &consecutiveFailures,
			},
		}

		result, _ := prog.Apply(ctx, params)

		if result.Applied {
			t.Error("Expected Applied=false for manual intervention case")
		}
		if result.Reason == "" {
			t.Error("Expected reason to be set for non-applied result")
		}
	})
}

// TestStageProgressionTriggerMismatch tests trigger type validation.
func TestStageProgressionTriggerMismatch(t *testing.T) {
	ctx := context.Background()

	prog, _ := progression.NewGZCLPT1DefaultProgression(uuid.New().String(), "Test")

	tests := []struct {
		name        string
		triggerType progression.TriggerType
	}{
		{"AFTER_SESSION", progression.TriggerAfterSession},
		{"AFTER_WEEK", progression.TriggerAfterWeek},
		{"AFTER_CYCLE", progression.TriggerAfterCycle},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := progression.ProgressionContext{
				UserID:       uuid.New().String(),
				LiftID:       uuid.New().String(),
				CurrentValue: 200.0,
				MaxType:      progression.TrainingMax,
				TriggerEvent: progression.TriggerEvent{
					Type:      tt.triggerType,
					Timestamp: time.Now(),
				},
			}

			result, err := prog.Apply(ctx, params)
			if err != nil {
				t.Fatalf("Apply failed: %v", err)
			}

			if result.Applied {
				t.Errorf("Expected Applied=false for trigger type %s", tt.triggerType)
			}
		})
	}
}

// TestStageProgressionMaxTypeMismatch tests max type validation.
func TestStageProgressionMaxTypeMismatch(t *testing.T) {
	ctx := context.Background()

	// Create progression with TRAINING_MAX
	prog, _ := progression.NewGZCLPT1DefaultProgression(uuid.New().String(), "Test")

	consecutiveFailures := 1
	params := progression.ProgressionContext{
		UserID:       uuid.New().String(),
		LiftID:       uuid.New().String(),
		CurrentValue: 200.0,
		MaxType:      progression.OneRM, // Mismatch!
		TriggerEvent: progression.TriggerEvent{
			Type:                progression.TriggerOnFailure,
			Timestamp:           time.Now(),
			ConsecutiveFailures: &consecutiveFailures,
		},
	}

	result, _ := prog.Apply(ctx, params)

	if result.Applied {
		t.Error("Expected Applied=false for max type mismatch")
	}
}

// TestStageProgressionShouldResetFailureCounter tests the reset flag.
func TestStageProgressionShouldResetFailureCounter(t *testing.T) {
	prog, _ := progression.NewGZCLPT1DefaultProgression(uuid.New().String(), "Test")

	// StageProgression should always reset failure counter on stage change
	if !prog.ShouldResetFailureCounter() {
		t.Error("StageProgression should reset failure counter")
	}
}

// TestStageGetCurrentSetScheme tests the set scheme retrieval.
func TestStageGetCurrentSetScheme(t *testing.T) {
	prog, _ := progression.NewGZCLPT1DefaultProgression(uuid.New().String(), "Test")

	// Stage 0: 5x3+
	scheme := prog.GetCurrentSetScheme()
	// The scheme should be AMRAP type
	sets, err := scheme.GenerateSets(200.0, setscheme.DefaultSetGenerationContext())
	if err != nil {
		t.Fatalf("Failed to generate sets: %v", err)
	}
	if len(sets) != 5 {
		t.Errorf("Expected 5 sets for stage 0, got %d", len(sets))
	}

	// Advance to stage 1: 6x2+
	prog.SetCurrentStage(1)
	scheme = prog.GetCurrentSetScheme()
	sets, err = scheme.GenerateSets(200.0, setscheme.DefaultSetGenerationContext())
	if err != nil {
		t.Fatalf("Failed to generate sets: %v", err)
	}
	if len(sets) != 6 {
		t.Errorf("Expected 6 sets for stage 1, got %d", len(sets))
	}
}

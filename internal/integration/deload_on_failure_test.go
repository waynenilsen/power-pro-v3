// Package integration provides integration tests for cross-component behavior.
// This file tests the DeloadOnFailure progression domain logic.
package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/progression"
)

// TestDeloadOnFailureCreate tests creating DeloadOnFailure progressions.
func TestDeloadOnFailureCreate(t *testing.T) {
	t.Run("percent deload creation", func(t *testing.T) {
		prog, err := progression.NewDeloadOnFailure(
			uuid.New().String(),
			"Test Percent Deload",
			2,           // 2 failures threshold
			"percent",   // percent deload
			0.10,        // 10% deload
			0,           // no fixed amount
			true,        // reset on deload
			progression.TrainingMax,
		)
		if err != nil {
			t.Fatalf("Failed to create: %v", err)
		}

		if prog.Type() != progression.TypeDeloadOnFailure {
			t.Errorf("Expected type %s, got %s", progression.TypeDeloadOnFailure, prog.Type())
		}
		if prog.TriggerType() != progression.TriggerOnFailure {
			t.Errorf("Expected trigger %s, got %s", progression.TriggerOnFailure, prog.TriggerType())
		}
		if prog.FailureThreshold != 2 {
			t.Errorf("Expected threshold 2, got %d", prog.FailureThreshold)
		}
	})

	t.Run("fixed deload creation", func(t *testing.T) {
		prog, err := progression.NewDeloadOnFailure(
			uuid.New().String(),
			"Test Fixed Deload",
			3,         // 3 failures threshold
			"fixed",   // fixed deload
			0,         // no percent
			10.0,      // 10lb deload
			true,
			progression.TrainingMax,
		)
		if err != nil {
			t.Fatalf("Failed to create: %v", err)
		}

		if prog.DeloadType != "fixed" {
			t.Errorf("Expected deload type 'fixed', got '%s'", prog.DeloadType)
		}
		if prog.DeloadAmount != 10.0 {
			t.Errorf("Expected deload amount 10, got %.1f", prog.DeloadAmount)
		}
	})

	t.Run("validation errors", func(t *testing.T) {
		// Zero threshold
		_, err := progression.NewDeloadOnFailure(
			uuid.New().String(), "Test", 0, "percent", 0.10, 0, true, progression.TrainingMax,
		)
		if err == nil {
			t.Error("Expected error for zero threshold")
		}

		// Invalid deload type
		_, err = progression.NewDeloadOnFailure(
			uuid.New().String(), "Test", 1, "invalid", 0.10, 0, true, progression.TrainingMax,
		)
		if err == nil {
			t.Error("Expected error for invalid deload type")
		}

		// Zero percent for percent type
		_, err = progression.NewDeloadOnFailure(
			uuid.New().String(), "Test", 1, "percent", 0, 0, true, progression.TrainingMax,
		)
		if err == nil {
			t.Error("Expected error for zero percent")
		}

		// Percent > 1
		_, err = progression.NewDeloadOnFailure(
			uuid.New().String(), "Test", 1, "percent", 1.5, 0, true, progression.TrainingMax,
		)
		if err == nil {
			t.Error("Expected error for percent > 1")
		}

		// Zero amount for fixed type
		_, err = progression.NewDeloadOnFailure(
			uuid.New().String(), "Test", 1, "fixed", 0, 0, true, progression.TrainingMax,
		)
		if err == nil {
			t.Error("Expected error for zero fixed amount")
		}
	})
}

// TestDeloadOnFailureThreshold tests that deload only triggers at threshold.
func TestDeloadOnFailureThreshold(t *testing.T) {
	ctx := context.Background()

	prog, _ := progression.NewDeloadOnFailure(
		uuid.New().String(),
		"Test",
		3,         // 3 failures threshold
		"percent",
		0.10,
		0,
		true,
		progression.TrainingMax,
	)

	currentWeight := 300.0
	userID := uuid.New().String()
	liftID := uuid.New().String()

	tests := []struct {
		name        string
		failures    int
		shouldApply bool
	}{
		{"1 failure - below threshold", 1, false},
		{"2 failures - below threshold", 2, false},
		{"3 failures - at threshold", 3, true},
		{"4 failures - above threshold", 4, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := progression.ProgressionContext{
				UserID:       userID,
				LiftID:       liftID,
				CurrentValue: currentWeight,
				MaxType:      progression.TrainingMax,
				TriggerEvent: progression.TriggerEvent{
					Type:                progression.TriggerOnFailure,
					Timestamp:           time.Now(),
					ConsecutiveFailures: &tt.failures,
				},
			}

			result, err := prog.Apply(ctx, params)
			if err != nil {
				t.Fatalf("Apply failed: %v", err)
			}

			if result.Applied != tt.shouldApply {
				t.Errorf("Expected Applied=%v, got %v. Reason: %s",
					tt.shouldApply, result.Applied, result.Reason)
			}
		})
	}
}

// TestDeloadOnFailurePercentCalculation tests percent-based deload calculation.
func TestDeloadOnFailurePercentCalculation(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		currentWeight  float64
		deloadPercent  float64
		expectedDelta  float64
		expectedNew    float64
	}{
		{"10% of 200", 200.0, 0.10, -20.0, 180.0},
		{"15% of 200", 200.0, 0.15, -30.0, 170.0},
		{"10% of 315", 315.0, 0.10, -31.5, 283.5},
		{"5% of 100", 100.0, 0.05, -5.0, 95.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prog, _ := progression.NewDeloadOnFailure(
				uuid.New().String(),
				"Test",
				1, // threshold 1
				"percent",
				tt.deloadPercent,
				0,
				true,
				progression.TrainingMax,
			)

			consecutiveFailures := 1
			params := progression.ProgressionContext{
				UserID:       uuid.New().String(),
				LiftID:       uuid.New().String(),
				CurrentValue: tt.currentWeight,
				MaxType:      progression.TrainingMax,
				TriggerEvent: progression.TriggerEvent{
					Type:                progression.TriggerOnFailure,
					Timestamp:           time.Now(),
					ConsecutiveFailures: &consecutiveFailures,
				},
			}

			result, _ := prog.Apply(ctx, params)

			if !result.Applied {
				t.Fatal("Expected progression to apply")
			}
			if result.Delta != tt.expectedDelta {
				t.Errorf("Expected delta=%.1f, got %.1f", tt.expectedDelta, result.Delta)
			}
			if result.NewValue != tt.expectedNew {
				t.Errorf("Expected new value=%.1f, got %.1f", tt.expectedNew, result.NewValue)
			}
		})
	}
}

// TestDeloadOnFailureFixedCalculation tests fixed-amount deload calculation.
func TestDeloadOnFailureFixedCalculation(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		currentWeight float64
		deloadAmount  float64
		expectedDelta float64
		expectedNew   float64
	}{
		{"5lb from 200", 200.0, 5.0, -5.0, 195.0},
		{"10lb from 315", 315.0, 10.0, -10.0, 305.0},
		{"2.5lb from 100", 100.0, 2.5, -2.5, 97.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prog, _ := progression.NewDeloadOnFailure(
				uuid.New().String(),
				"Test",
				1,
				"fixed",
				0,
				tt.deloadAmount,
				true,
				progression.TrainingMax,
			)

			consecutiveFailures := 1
			params := progression.ProgressionContext{
				UserID:       uuid.New().String(),
				LiftID:       uuid.New().String(),
				CurrentValue: tt.currentWeight,
				MaxType:      progression.TrainingMax,
				TriggerEvent: progression.TriggerEvent{
					Type:                progression.TriggerOnFailure,
					Timestamp:           time.Now(),
					ConsecutiveFailures: &consecutiveFailures,
				},
			}

			result, _ := prog.Apply(ctx, params)

			if !result.Applied {
				t.Fatal("Expected progression to apply")
			}
			if result.Delta != tt.expectedDelta {
				t.Errorf("Expected delta=%.1f, got %.1f", tt.expectedDelta, result.Delta)
			}
			if result.NewValue != tt.expectedNew {
				t.Errorf("Expected new value=%.1f, got %.1f", tt.expectedNew, result.NewValue)
			}
		})
	}
}

// TestDeloadOnFailureNoNegativeWeight tests that weight doesn't go below zero.
func TestDeloadOnFailureNoNegativeWeight(t *testing.T) {
	ctx := context.Background()

	// Fixed deload larger than current weight
	prog, _ := progression.NewDeloadOnFailure(
		uuid.New().String(),
		"Test",
		1,
		"fixed",
		0,
		100.0, // 100lb deload
		true,
		progression.TrainingMax,
	)

	consecutiveFailures := 1
	params := progression.ProgressionContext{
		UserID:       uuid.New().String(),
		LiftID:       uuid.New().String(),
		CurrentValue: 50.0, // Only 50lb
		MaxType:      progression.TrainingMax,
		TriggerEvent: progression.TriggerEvent{
			Type:                progression.TriggerOnFailure,
			ConsecutiveFailures: &consecutiveFailures,
		},
	}

	result, _ := prog.Apply(ctx, params)

	if result.NewValue < 0 {
		t.Errorf("Weight should not go below 0, got %.1f", result.NewValue)
	}
	if result.NewValue != 0 {
		t.Errorf("Expected new value=0 when deload exceeds current, got %.1f", result.NewValue)
	}
}

// TestDeloadOnFailureShouldResetFailureCounter tests the reset flag.
func TestDeloadOnFailureShouldResetFailureCounter(t *testing.T) {
	// With ResetOnDeload=true
	prog1, _ := progression.NewDeloadOnFailure(
		uuid.New().String(), "Test", 1, "percent", 0.10, 0, true, progression.TrainingMax,
	)
	if !prog1.ShouldResetFailureCounter() {
		t.Error("Should reset counter when ResetOnDeload=true")
	}

	// With ResetOnDeload=false
	prog2, _ := progression.NewDeloadOnFailure(
		uuid.New().String(), "Test", 1, "percent", 0.10, 0, false, progression.TrainingMax,
	)
	if prog2.ShouldResetFailureCounter() {
		t.Error("Should not reset counter when ResetOnDeload=false")
	}
}

// TestDeloadOnFailureTriggerMismatch tests trigger type validation.
func TestDeloadOnFailureTriggerMismatch(t *testing.T) {
	ctx := context.Background()

	prog, _ := progression.NewDeloadOnFailure(
		uuid.New().String(), "Test", 1, "percent", 0.10, 0, true, progression.TrainingMax,
	)

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
					Type: tt.triggerType,
				},
			}

			result, _ := prog.Apply(ctx, params)

			if result.Applied {
				t.Errorf("Expected Applied=false for trigger type %s", tt.triggerType)
			}
		})
	}
}

// TestDeloadOnFailureMissingFailureCount tests handling of missing failure count.
func TestDeloadOnFailureMissingFailureCount(t *testing.T) {
	ctx := context.Background()

	prog, _ := progression.NewDeloadOnFailure(
		uuid.New().String(), "Test", 1, "percent", 0.10, 0, true, progression.TrainingMax,
	)

	params := progression.ProgressionContext{
		UserID:       uuid.New().String(),
		LiftID:       uuid.New().String(),
		CurrentValue: 200.0,
		MaxType:      progression.TrainingMax,
		TriggerEvent: progression.TriggerEvent{
			Type:                progression.TriggerOnFailure,
			ConsecutiveFailures: nil, // Missing!
		},
	}

	result, _ := prog.Apply(ctx, params)

	if result.Applied {
		t.Error("Expected Applied=false when ConsecutiveFailures is nil")
	}
}

// TestDeloadOnFailureMaxTypeMismatch tests max type validation.
func TestDeloadOnFailureMaxTypeMismatch(t *testing.T) {
	ctx := context.Background()

	// Create with TRAINING_MAX
	prog, _ := progression.NewDeloadOnFailure(
		uuid.New().String(), "Test", 1, "percent", 0.10, 0, true, progression.TrainingMax,
	)

	consecutiveFailures := 1
	params := progression.ProgressionContext{
		UserID:       uuid.New().String(),
		LiftID:       uuid.New().String(),
		CurrentValue: 200.0,
		MaxType:      progression.OneRM, // Mismatch!
		TriggerEvent: progression.TriggerEvent{
			Type:                progression.TriggerOnFailure,
			ConsecutiveFailures: &consecutiveFailures,
		},
	}

	result, _ := prog.Apply(ctx, params)

	if result.Applied {
		t.Error("Expected Applied=false for max type mismatch")
	}
}

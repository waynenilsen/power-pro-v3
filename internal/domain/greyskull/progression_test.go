package greyskull

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/domain/progression"
)

// intPtr is a helper function for creating int pointers.
func intPtr(i int) *int {
	return &i
}

// TestGreySkullProgression_Type tests that GreySkullProgression returns correct type.
func TestGreySkullProgression_Type(t *testing.T) {
	gsp := &GreySkullProgression{
		ID:              "prog-1",
		Name:            "Test Progression",
		WeightIncrement: 2.5,
		MinReps:         5,
		DoubleThreshold: 10,
		DeloadPercent:   0.10,
		MaxTypeValue:    progression.TrainingMax,
	}
	if gsp.Type() != progression.TypeGreySkull {
		t.Errorf("expected %s, got %s", progression.TypeGreySkull, gsp.Type())
	}
}

// TestGreySkullProgression_TriggerType tests that GreySkullProgression returns correct trigger type.
func TestGreySkullProgression_TriggerType(t *testing.T) {
	gsp := &GreySkullProgression{
		ID:              "prog-1",
		Name:            "Test",
		WeightIncrement: 2.5,
		MinReps:         5,
		DoubleThreshold: 10,
		DeloadPercent:   0.10,
		MaxTypeValue:    progression.TrainingMax,
	}
	if gsp.TriggerType() != progression.TriggerAfterSet {
		t.Errorf("expected %s, got %s", progression.TriggerAfterSet, gsp.TriggerType())
	}
}

// TestGreySkullProgression_Validate tests GreySkullProgression validation.
func TestGreySkullProgression_Validate(t *testing.T) {
	tests := []struct {
		name    string
		gsp     GreySkullProgression
		wantErr bool
		errType error
	}{
		{
			name: "valid main lift progression",
			gsp: GreySkullProgression{
				ID:              "prog-1",
				Name:            "Bench Press Progression",
				WeightIncrement: 2.5,
				MinReps:         5,
				DoubleThreshold: 10,
				DeloadPercent:   0.10,
				MaxTypeValue:    progression.TrainingMax,
			},
			wantErr: false,
		},
		{
			name: "valid accessory progression",
			gsp: GreySkullProgression{
				ID:              "prog-2",
				Name:            "Curl Progression",
				WeightIncrement: 2.5,
				MinReps:         10,
				DoubleThreshold: 15,
				DeloadPercent:   0.10,
				MaxTypeValue:    progression.TrainingMax,
			},
			wantErr: false,
		},
		{
			name: "valid lower body progression",
			gsp: GreySkullProgression{
				ID:              "prog-3",
				Name:            "Squat Progression",
				WeightIncrement: 5.0,
				MinReps:         5,
				DoubleThreshold: 10,
				DeloadPercent:   0.10,
				MaxTypeValue:    progression.TrainingMax,
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			gsp: GreySkullProgression{
				ID:              "",
				Name:            "Test",
				WeightIncrement: 2.5,
				MinReps:         5,
				DoubleThreshold: 10,
				DeloadPercent:   0.10,
				MaxTypeValue:    progression.TrainingMax,
			},
			wantErr: true,
			errType: progression.ErrInvalidParams,
		},
		{
			name: "missing name",
			gsp: GreySkullProgression{
				ID:              "prog-1",
				Name:            "",
				WeightIncrement: 2.5,
				MinReps:         5,
				DoubleThreshold: 10,
				DeloadPercent:   0.10,
				MaxTypeValue:    progression.TrainingMax,
			},
			wantErr: true,
			errType: progression.ErrInvalidParams,
		},
		{
			name: "zero weight increment",
			gsp: GreySkullProgression{
				ID:              "prog-1",
				Name:            "Test",
				WeightIncrement: 0,
				MinReps:         5,
				DoubleThreshold: 10,
				DeloadPercent:   0.10,
				MaxTypeValue:    progression.TrainingMax,
			},
			wantErr: true,
			errType: progression.ErrInvalidParams,
		},
		{
			name: "negative weight increment",
			gsp: GreySkullProgression{
				ID:              "prog-1",
				Name:            "Test",
				WeightIncrement: -2.5,
				MinReps:         5,
				DoubleThreshold: 10,
				DeloadPercent:   0.10,
				MaxTypeValue:    progression.TrainingMax,
			},
			wantErr: true,
			errType: progression.ErrInvalidParams,
		},
		{
			name: "minReps zero",
			gsp: GreySkullProgression{
				ID:              "prog-1",
				Name:            "Test",
				WeightIncrement: 2.5,
				MinReps:         0,
				DoubleThreshold: 10,
				DeloadPercent:   0.10,
				MaxTypeValue:    progression.TrainingMax,
			},
			wantErr: true,
			errType: progression.ErrInvalidParams,
		},
		{
			name: "doubleThreshold equals minReps",
			gsp: GreySkullProgression{
				ID:              "prog-1",
				Name:            "Test",
				WeightIncrement: 2.5,
				MinReps:         5,
				DoubleThreshold: 5,
				DeloadPercent:   0.10,
				MaxTypeValue:    progression.TrainingMax,
			},
			wantErr: true,
			errType: progression.ErrInvalidParams,
		},
		{
			name: "doubleThreshold less than minReps",
			gsp: GreySkullProgression{
				ID:              "prog-1",
				Name:            "Test",
				WeightIncrement: 2.5,
				MinReps:         10,
				DoubleThreshold: 5,
				DeloadPercent:   0.10,
				MaxTypeValue:    progression.TrainingMax,
			},
			wantErr: true,
			errType: progression.ErrInvalidParams,
		},
		{
			name: "zero deload percent",
			gsp: GreySkullProgression{
				ID:              "prog-1",
				Name:            "Test",
				WeightIncrement: 2.5,
				MinReps:         5,
				DoubleThreshold: 10,
				DeloadPercent:   0,
				MaxTypeValue:    progression.TrainingMax,
			},
			wantErr: true,
			errType: progression.ErrInvalidParams,
		},
		{
			name: "deload percent greater than 1",
			gsp: GreySkullProgression{
				ID:              "prog-1",
				Name:            "Test",
				WeightIncrement: 2.5,
				MinReps:         5,
				DoubleThreshold: 10,
				DeloadPercent:   1.5,
				MaxTypeValue:    progression.TrainingMax,
			},
			wantErr: true,
			errType: progression.ErrInvalidParams,
		},
		{
			name: "negative deload percent",
			gsp: GreySkullProgression{
				ID:              "prog-1",
				Name:            "Test",
				WeightIncrement: 2.5,
				MinReps:         5,
				DoubleThreshold: 10,
				DeloadPercent:   -0.10,
				MaxTypeValue:    progression.TrainingMax,
			},
			wantErr: true,
			errType: progression.ErrInvalidParams,
		},
		{
			name: "missing max type",
			gsp: GreySkullProgression{
				ID:              "prog-1",
				Name:            "Test",
				WeightIncrement: 2.5,
				MinReps:         5,
				DoubleThreshold: 10,
				DeloadPercent:   0.10,
				MaxTypeValue:    "",
			},
			wantErr: true,
			errType: progression.ErrUnknownMaxType,
		},
		{
			name: "invalid max type",
			gsp: GreySkullProgression{
				ID:              "prog-1",
				Name:            "Test",
				WeightIncrement: 2.5,
				MinReps:         5,
				DoubleThreshold: 10,
				DeloadPercent:   0.10,
				MaxTypeValue:    "INVALID",
			},
			wantErr: true,
			errType: progression.ErrUnknownMaxType,
		},
		{
			name: "deload percent exactly 1 (100%)",
			gsp: GreySkullProgression{
				ID:              "prog-1",
				Name:            "Test",
				WeightIncrement: 2.5,
				MinReps:         5,
				DoubleThreshold: 10,
				DeloadPercent:   1.0,
				MaxTypeValue:    progression.TrainingMax,
			},
			wantErr: false, // 100% deload is valid (though unusual)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.gsp.Validate()
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if tt.errType != nil && !errors.Is(err, tt.errType) {
					t.Errorf("expected %v, got %v", tt.errType, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// TestNewGreySkullProgression tests the factory function.
func TestNewGreySkullProgression(t *testing.T) {
	t.Run("valid parameters", func(t *testing.T) {
		gsp, err := NewGreySkullProgression("prog-1", "Bench Progress", 2.5, 5, 10, 0.10, progression.TrainingMax)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if gsp.ID != "prog-1" {
			t.Errorf("expected ID 'prog-1', got '%s'", gsp.ID)
		}
		if gsp.Name != "Bench Progress" {
			t.Errorf("expected Name 'Bench Progress', got '%s'", gsp.Name)
		}
		if gsp.WeightIncrement != 2.5 {
			t.Errorf("expected WeightIncrement 2.5, got %f", gsp.WeightIncrement)
		}
		if gsp.MinReps != 5 {
			t.Errorf("expected MinReps 5, got %d", gsp.MinReps)
		}
		if gsp.DoubleThreshold != 10 {
			t.Errorf("expected DoubleThreshold 10, got %d", gsp.DoubleThreshold)
		}
		if gsp.DeloadPercent != 0.10 {
			t.Errorf("expected DeloadPercent 0.10, got %f", gsp.DeloadPercent)
		}
	})

	t.Run("invalid parameters", func(t *testing.T) {
		_, err := NewGreySkullProgression("", "Test", 2.5, 5, 10, 0.10, progression.TrainingMax)
		if err == nil {
			t.Error("expected error for empty ID")
		}
	})
}

// TestNewGreySkullMainLiftProgression tests the main lift factory function.
func TestNewGreySkullMainLiftProgression(t *testing.T) {
	gsp, err := NewGreySkullMainLiftProgression("prog-1", "Squat Progress", 5.0, progression.TrainingMax)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gsp.MinReps != 5 {
		t.Errorf("expected MinReps 5, got %d", gsp.MinReps)
	}
	if gsp.DoubleThreshold != 10 {
		t.Errorf("expected DoubleThreshold 10, got %d", gsp.DoubleThreshold)
	}
	if gsp.DeloadPercent != 0.10 {
		t.Errorf("expected DeloadPercent 0.10, got %f", gsp.DeloadPercent)
	}
}

// TestNewGreySkullAccessoryProgression tests the accessory factory function.
func TestNewGreySkullAccessoryProgression(t *testing.T) {
	gsp, err := NewGreySkullAccessoryProgression("prog-1", "Curl Progress", 2.5, progression.TrainingMax)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gsp.MinReps != 10 {
		t.Errorf("expected MinReps 10, got %d", gsp.MinReps)
	}
	if gsp.DoubleThreshold != 15 {
		t.Errorf("expected DoubleThreshold 15, got %d", gsp.DoubleThreshold)
	}
	if gsp.DeloadPercent != 0.10 {
		t.Errorf("expected DeloadPercent 0.10, got %f", gsp.DeloadPercent)
	}
}

// TestGreySkullProgression_Apply tests the Apply method.
func TestGreySkullProgression_Apply(t *testing.T) {
	gsp := &GreySkullProgression{
		ID:              "prog-1",
		Name:            "GreySkull Main Lift",
		WeightIncrement: 2.5,
		MinReps:         5,
		DoubleThreshold: 10,
		DeloadPercent:   0.10,
		MaxTypeValue:    progression.TrainingMax,
	}

	ctx := context.Background()

	t.Run("deload - reps less than minReps (3 reps)", func(t *testing.T) {
		params := progression.ProgressionContext{
			UserID:       "user-123",
			LiftID:       "bench-uuid",
			MaxType:      progression.TrainingMax,
			CurrentValue: 200,
			TriggerEvent: progression.TriggerEvent{
				Type:          progression.TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(3),
				IsAMRAP:       true,
			},
		}

		result, err := gsp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true, reason: %s", result.Reason)
		}
		// 10% deload: 200 * 0.10 = 20, new value = 180
		expectedDelta := -20.0
		expectedNewValue := 180.0
		if result.Delta != expectedDelta {
			t.Errorf("expected Delta %f, got %f", expectedDelta, result.Delta)
		}
		if result.NewValue != expectedNewValue {
			t.Errorf("expected NewValue %f, got %f", expectedNewValue, result.NewValue)
		}
	})

	t.Run("deload - exactly at failure threshold (4 reps)", func(t *testing.T) {
		params := progression.ProgressionContext{
			UserID:       "user-123",
			LiftID:       "bench-uuid",
			MaxType:      progression.TrainingMax,
			CurrentValue: 200,
			TriggerEvent: progression.TriggerEvent{
				Type:          progression.TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(4),
				IsAMRAP:       true,
			},
		}

		result, err := gsp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true, reason: %s", result.Reason)
		}
		// 4 < 5, so deload
		expectedDelta := -20.0
		if result.Delta != expectedDelta {
			t.Errorf("expected Delta %f (deload), got %f", expectedDelta, result.Delta)
		}
	})

	t.Run("standard increment - exactly at minReps (5 reps)", func(t *testing.T) {
		params := progression.ProgressionContext{
			UserID:       "user-123",
			LiftID:       "bench-uuid",
			MaxType:      progression.TrainingMax,
			CurrentValue: 200,
			TriggerEvent: progression.TriggerEvent{
				Type:          progression.TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(5),
				IsAMRAP:       true,
			},
		}

		result, err := gsp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true, reason: %s", result.Reason)
		}
		if result.Delta != 2.5 {
			t.Errorf("expected Delta 2.5, got %f", result.Delta)
		}
		if result.NewValue != 202.5 {
			t.Errorf("expected NewValue 202.5, got %f", result.NewValue)
		}
	})

	t.Run("standard increment - 7 reps (in standard range)", func(t *testing.T) {
		params := progression.ProgressionContext{
			UserID:       "user-123",
			LiftID:       "bench-uuid",
			MaxType:      progression.TrainingMax,
			CurrentValue: 200,
			TriggerEvent: progression.TriggerEvent{
				Type:          progression.TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(7),
				IsAMRAP:       true,
			},
		}

		result, err := gsp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true, reason: %s", result.Reason)
		}
		// 7 reps is in [5, 10) range, standard increment
		if result.Delta != 2.5 {
			t.Errorf("expected Delta 2.5, got %f", result.Delta)
		}
	})

	t.Run("standard increment - 9 reps (edge of standard range)", func(t *testing.T) {
		params := progression.ProgressionContext{
			UserID:       "user-123",
			LiftID:       "bench-uuid",
			MaxType:      progression.TrainingMax,
			CurrentValue: 200,
			TriggerEvent: progression.TriggerEvent{
				Type:          progression.TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(9),
				IsAMRAP:       true,
			},
		}

		result, err := gsp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// 9 reps is in [5, 10) range, standard increment
		if result.Delta != 2.5 {
			t.Errorf("expected Delta 2.5, got %f", result.Delta)
		}
	})

	t.Run("double increment - exactly at doubleThreshold (10 reps)", func(t *testing.T) {
		params := progression.ProgressionContext{
			UserID:       "user-123",
			LiftID:       "bench-uuid",
			MaxType:      progression.TrainingMax,
			CurrentValue: 200,
			TriggerEvent: progression.TriggerEvent{
				Type:          progression.TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(10),
				IsAMRAP:       true,
			},
		}

		result, err := gsp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true, reason: %s", result.Reason)
		}
		// 10+ reps, double increment: 2.5 * 2 = 5
		if result.Delta != 5.0 {
			t.Errorf("expected Delta 5.0, got %f", result.Delta)
		}
		if result.NewValue != 205.0 {
			t.Errorf("expected NewValue 205.0, got %f", result.NewValue)
		}
	})

	t.Run("double increment - 15 reps (well above threshold)", func(t *testing.T) {
		params := progression.ProgressionContext{
			UserID:       "user-123",
			LiftID:       "bench-uuid",
			MaxType:      progression.TrainingMax,
			CurrentValue: 200,
			TriggerEvent: progression.TriggerEvent{
				Type:          progression.TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(15),
				IsAMRAP:       true,
			},
		}

		result, err := gsp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true, reason: %s", result.Reason)
		}
		// 15 >= 10, double increment
		if result.Delta != 5.0 {
			t.Errorf("expected Delta 5.0, got %f", result.Delta)
		}
	})

	t.Run("trigger type mismatch", func(t *testing.T) {
		params := progression.ProgressionContext{
			UserID:       "user-123",
			LiftID:       "bench-uuid",
			MaxType:      progression.TrainingMax,
			CurrentValue: 200,
			TriggerEvent: progression.TriggerEvent{
				Type:      progression.TriggerAfterSession, // Wrong trigger type
				Timestamp: time.Now(),
			},
		}

		result, err := gsp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied to be false for trigger type mismatch")
		}
		if result.Reason == "" {
			t.Error("expected Reason to be set")
		}
	})

	t.Run("max type mismatch", func(t *testing.T) {
		params := progression.ProgressionContext{
			UserID:       "user-123",
			LiftID:       "bench-uuid",
			MaxType:      progression.OneRM, // Wrong max type
			CurrentValue: 225,
			TriggerEvent: progression.TriggerEvent{
				Type:          progression.TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(7),
				IsAMRAP:       true,
			},
		}

		result, err := gsp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied to be false for max type mismatch")
		}
		if result.Reason == "" {
			t.Error("expected Reason to be set")
		}
	})

	t.Run("not an AMRAP set", func(t *testing.T) {
		params := progression.ProgressionContext{
			UserID:       "user-123",
			LiftID:       "bench-uuid",
			MaxType:      progression.TrainingMax,
			CurrentValue: 200,
			TriggerEvent: progression.TriggerEvent{
				Type:          progression.TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(8),
				IsAMRAP:       false, // Not an AMRAP set
			},
		}

		result, err := gsp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied to be false for non-AMRAP set")
		}
		if result.Reason == "" {
			t.Error("expected Reason to be set")
		}
	})

	t.Run("reps not provided", func(t *testing.T) {
		params := progression.ProgressionContext{
			UserID:       "user-123",
			LiftID:       "bench-uuid",
			MaxType:      progression.TrainingMax,
			CurrentValue: 200,
			TriggerEvent: progression.TriggerEvent{
				Type:      progression.TriggerAfterSet,
				Timestamp: time.Now(),
				IsAMRAP:   true,
				// RepsPerformed is nil
			},
		}

		result, err := gsp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied to be false when reps not provided")
		}
		if result.Reason == "" {
			t.Error("expected Reason to be set")
		}
	})

	t.Run("invalid context (missing userID)", func(t *testing.T) {
		params := progression.ProgressionContext{
			UserID:       "", // Missing userID
			LiftID:       "bench-uuid",
			MaxType:      progression.TrainingMax,
			CurrentValue: 200,
			TriggerEvent: progression.TriggerEvent{
				Type:          progression.TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(7),
				IsAMRAP:       true,
			},
		}

		_, err := gsp.Apply(ctx, params)
		if err == nil {
			t.Error("expected error for invalid context")
		}
	})
}

// TestGreySkullProgression_Apply_DeloadDoesNotGoBelowZero tests deload edge case.
func TestGreySkullProgression_Apply_DeloadDoesNotGoBelowZero(t *testing.T) {
	gsp := &GreySkullProgression{
		ID:              "prog-1",
		Name:            "Test",
		WeightIncrement: 2.5,
		MinReps:         5,
		DoubleThreshold: 10,
		DeloadPercent:   0.10,
		MaxTypeValue:    progression.TrainingMax,
	}

	ctx := context.Background()

	t.Run("deload on small weight does not go negative", func(t *testing.T) {
		params := progression.ProgressionContext{
			UserID:       "user-123",
			LiftID:       "bench-uuid",
			MaxType:      progression.TrainingMax,
			CurrentValue: 5, // Very low weight
			TriggerEvent: progression.TriggerEvent{
				Type:          progression.TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(2), // Triggers deload
				IsAMRAP:       true,
			},
		}

		result, err := gsp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.NewValue < 0 {
			t.Errorf("NewValue should not be negative, got %f", result.NewValue)
		}
	})

	t.Run("100% deload results in zero", func(t *testing.T) {
		gspFull := &GreySkullProgression{
			ID:              "prog-1",
			Name:            "Test",
			WeightIncrement: 2.5,
			MinReps:         5,
			DoubleThreshold: 10,
			DeloadPercent:   1.0, // 100% deload
			MaxTypeValue:    progression.TrainingMax,
		}

		params := progression.ProgressionContext{
			UserID:       "user-123",
			LiftID:       "bench-uuid",
			MaxType:      progression.TrainingMax,
			CurrentValue: 100,
			TriggerEvent: progression.TriggerEvent{
				Type:          progression.TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(2),
				IsAMRAP:       true,
			},
		}

		result, err := gspFull.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.NewValue != 0 {
			t.Errorf("expected NewValue 0 with 100%% deload, got %f", result.NewValue)
		}
		if result.Delta != -100 {
			t.Errorf("expected Delta -100, got %f", result.Delta)
		}
	})
}

// TestGreySkullProgression_Apply_AccessoryVariant tests accessory progression behavior.
func TestGreySkullProgression_Apply_AccessoryVariant(t *testing.T) {
	gsp := &GreySkullProgression{
		ID:              "prog-1",
		Name:            "Curl Progression",
		WeightIncrement: 2.5,
		MinReps:         10, // Accessory thresholds
		DoubleThreshold: 15,
		DeloadPercent:   0.10,
		MaxTypeValue:    progression.TrainingMax,
	}

	ctx := context.Background()

	testCases := []struct {
		reps          int
		expectedDelta float64
		description   string
	}{
		{reps: 8, expectedDelta: -10.0, description: "8 reps (< 10): deload 10% of 100"},
		{reps: 9, expectedDelta: -10.0, description: "9 reps (< 10): deload"},
		{reps: 10, expectedDelta: 2.5, description: "10 reps (exactly minReps): standard"},
		{reps: 12, expectedDelta: 2.5, description: "12 reps (in standard range): standard"},
		{reps: 14, expectedDelta: 2.5, description: "14 reps (edge of standard): standard"},
		{reps: 15, expectedDelta: 5.0, description: "15 reps (exactly doubleThreshold): double"},
		{reps: 20, expectedDelta: 5.0, description: "20 reps (well above): double"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			params := progression.ProgressionContext{
				UserID:       "user-123",
				LiftID:       "curl-uuid",
				MaxType:      progression.TrainingMax,
				CurrentValue: 100,
				TriggerEvent: progression.TriggerEvent{
					Type:          progression.TriggerAfterSet,
					Timestamp:     time.Now(),
					RepsPerformed: intPtr(tc.reps),
					IsAMRAP:       true,
				},
			}

			result, err := gsp.Apply(ctx, params)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !result.Applied {
				t.Errorf("expected Applied=true, reason: %s", result.Reason)
			}
			if result.Delta != tc.expectedDelta {
				t.Errorf("expected Delta %f, got %f", tc.expectedDelta, result.Delta)
			}
		})
	}
}

// TestGreySkullProgression_Apply_LowerBody tests lower body progression with 5lb increment.
func TestGreySkullProgression_Apply_LowerBody(t *testing.T) {
	gsp := &GreySkullProgression{
		ID:              "prog-1",
		Name:            "Squat Progression",
		WeightIncrement: 5.0, // Lower body uses 5lb
		MinReps:         5,
		DoubleThreshold: 10,
		DeloadPercent:   0.10,
		MaxTypeValue:    progression.TrainingMax,
	}

	ctx := context.Background()

	t.Run("standard increment lower body", func(t *testing.T) {
		params := progression.ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      progression.TrainingMax,
			CurrentValue: 300,
			TriggerEvent: progression.TriggerEvent{
				Type:          progression.TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(7),
				IsAMRAP:       true,
			},
		}

		result, err := gsp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Delta != 5.0 {
			t.Errorf("expected Delta 5.0, got %f", result.Delta)
		}
		if result.NewValue != 305.0 {
			t.Errorf("expected NewValue 305.0, got %f", result.NewValue)
		}
	})

	t.Run("double increment lower body", func(t *testing.T) {
		params := progression.ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      progression.TrainingMax,
			CurrentValue: 300,
			TriggerEvent: progression.TriggerEvent{
				Type:          progression.TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(12),
				IsAMRAP:       true,
			},
		}

		result, err := gsp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Double: 5.0 * 2 = 10.0
		if result.Delta != 10.0 {
			t.Errorf("expected Delta 10.0, got %f", result.Delta)
		}
		if result.NewValue != 310.0 {
			t.Errorf("expected NewValue 310.0, got %f", result.NewValue)
		}
	})
}

// TestGreySkullProgression_JSON tests JSON serialization roundtrip.
func TestGreySkullProgression_JSON(t *testing.T) {
	gsp := &GreySkullProgression{
		ID:              "prog-123",
		Name:            "Test GreySkull Progression",
		WeightIncrement: 2.5,
		MinReps:         5,
		DoubleThreshold: 10,
		DeloadPercent:   0.10,
		MaxTypeValue:    progression.TrainingMax,
	}

	// Marshal
	data, err := json.Marshal(gsp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Verify JSON structure
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if parsed["type"] != string(progression.TypeGreySkull) {
		t.Errorf("expected type %s, got %v", progression.TypeGreySkull, parsed["type"])
	}
	if parsed["id"] != "prog-123" {
		t.Errorf("expected id 'prog-123', got %v", parsed["id"])
	}
	if parsed["name"] != "Test GreySkull Progression" {
		t.Errorf("expected name 'Test GreySkull Progression', got %v", parsed["name"])
	}
	if parsed["maxType"] != string(progression.TrainingMax) {
		t.Errorf("expected maxType %s, got %v", progression.TrainingMax, parsed["maxType"])
	}
	if parsed["weightIncrement"].(float64) != 2.5 {
		t.Errorf("expected weightIncrement 2.5, got %v", parsed["weightIncrement"])
	}
	if int(parsed["minReps"].(float64)) != 5 {
		t.Errorf("expected minReps 5, got %v", parsed["minReps"])
	}
	if int(parsed["doubleThreshold"].(float64)) != 10 {
		t.Errorf("expected doubleThreshold 10, got %v", parsed["doubleThreshold"])
	}
	if parsed["deloadPercent"].(float64) != 0.10 {
		t.Errorf("expected deloadPercent 0.10, got %v", parsed["deloadPercent"])
	}

	// Unmarshal back
	var restored GreySkullProgression
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if restored.ID != gsp.ID {
		t.Errorf("ID mismatch: expected %s, got %s", gsp.ID, restored.ID)
	}
	if restored.Name != gsp.Name {
		t.Errorf("Name mismatch: expected %s, got %s", gsp.Name, restored.Name)
	}
	if restored.WeightIncrement != gsp.WeightIncrement {
		t.Errorf("WeightIncrement mismatch: expected %f, got %f", gsp.WeightIncrement, restored.WeightIncrement)
	}
	if restored.MinReps != gsp.MinReps {
		t.Errorf("MinReps mismatch: expected %d, got %d", gsp.MinReps, restored.MinReps)
	}
	if restored.DoubleThreshold != gsp.DoubleThreshold {
		t.Errorf("DoubleThreshold mismatch: expected %d, got %d", gsp.DoubleThreshold, restored.DoubleThreshold)
	}
	if restored.DeloadPercent != gsp.DeloadPercent {
		t.Errorf("DeloadPercent mismatch: expected %f, got %f", gsp.DeloadPercent, restored.DeloadPercent)
	}
	if restored.MaxTypeValue != gsp.MaxTypeValue {
		t.Errorf("MaxTypeValue mismatch: expected %s, got %s", gsp.MaxTypeValue, restored.MaxTypeValue)
	}
}

// TestUnmarshalGreySkullProgression tests deserialization function.
func TestUnmarshalGreySkullProgression(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		jsonData := []byte(`{
			"id": "prog-1",
			"name": "Test",
			"weightIncrement": 2.5,
			"minReps": 5,
			"doubleThreshold": 10,
			"deloadPercent": 0.10,
			"maxType": "TRAINING_MAX"
		}`)

		prog, err := UnmarshalGreySkullProgression(jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		gsp, ok := prog.(*GreySkullProgression)
		if !ok {
			t.Fatalf("expected *GreySkullProgression, got %T", prog)
		}
		if gsp.ID != "prog-1" {
			t.Errorf("expected ID 'prog-1', got '%s'", gsp.ID)
		}
		if gsp.Type() != progression.TypeGreySkull {
			t.Errorf("expected type %s, got %s", progression.TypeGreySkull, gsp.Type())
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		jsonData := []byte(`{invalid}`)
		_, err := UnmarshalGreySkullProgression(jsonData)
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("invalid progression data (missing id)", func(t *testing.T) {
		jsonData := []byte(`{
			"id": "",
			"name": "Test",
			"weightIncrement": 2.5,
			"minReps": 5,
			"doubleThreshold": 10,
			"deloadPercent": 0.10,
			"maxType": "TRAINING_MAX"
		}`)

		_, err := UnmarshalGreySkullProgression(jsonData)
		if err == nil {
			t.Error("expected error for invalid progression data")
		}
	})

	t.Run("invalid progression data (doubleThreshold <= minReps)", func(t *testing.T) {
		jsonData := []byte(`{
			"id": "prog-1",
			"name": "Test",
			"weightIncrement": 2.5,
			"minReps": 10,
			"doubleThreshold": 5,
			"deloadPercent": 0.10,
			"maxType": "TRAINING_MAX"
		}`)

		_, err := UnmarshalGreySkullProgression(jsonData)
		if err == nil {
			t.Error("expected error for invalid thresholds")
		}
	})
}

// TestRegisterGreySkullProgression tests factory registration.
func TestRegisterGreySkullProgression(t *testing.T) {
	factory := progression.NewProgressionFactory()

	// Verify not registered initially
	if factory.IsRegistered(progression.TypeGreySkull) {
		t.Error("TypeGreySkull should not be registered initially")
	}

	// Register
	RegisterGreySkullProgression(factory)

	// Verify registered
	if !factory.IsRegistered(progression.TypeGreySkull) {
		t.Error("TypeGreySkull should be registered after calling RegisterGreySkullProgression")
	}

	// Create from factory
	jsonData := []byte(`{
		"type": "GREYSKULL_PROGRESSION",
		"id": "prog-1",
		"name": "Factory Test",
		"weightIncrement": 2.5,
		"minReps": 5,
		"doubleThreshold": 10,
		"deloadPercent": 0.10,
		"maxType": "TRAINING_MAX"
	}`)

	prog, err := factory.CreateFromJSON(jsonData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prog.Type() != progression.TypeGreySkull {
		t.Errorf("expected type %s, got %s", progression.TypeGreySkull, prog.Type())
	}
}

// TestGreySkullProgression_Interface verifies that GreySkullProgression implements Progression.
func TestGreySkullProgression_Interface(t *testing.T) {
	var _ progression.Progression = &GreySkullProgression{}
}

// TestGreySkullProgression_AppliedAt tests that AppliedAt is set correctly.
func TestGreySkullProgression_AppliedAt(t *testing.T) {
	gsp := &GreySkullProgression{
		ID:              "prog-1",
		Name:            "Test",
		WeightIncrement: 2.5,
		MinReps:         5,
		DoubleThreshold: 10,
		DeloadPercent:   0.10,
		MaxTypeValue:    progression.TrainingMax,
	}

	before := time.Now()

	params := progression.ProgressionContext{
		UserID:       "user-1",
		LiftID:       "lift-1",
		MaxType:      progression.TrainingMax,
		CurrentValue: 100,
		TriggerEvent: progression.TriggerEvent{
			Type:          progression.TriggerAfterSet,
			Timestamp:     time.Now(),
			RepsPerformed: intPtr(7),
			IsAMRAP:       true,
		},
	}

	result, err := gsp.Apply(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	after := time.Now()

	if result.AppliedAt.Before(before) {
		t.Error("AppliedAt should not be before the test started")
	}
	if result.AppliedAt.After(after) {
		t.Error("AppliedAt should not be after the test ended")
	}
}

// TestGreySkullProgression_RealisticScenarios tests realistic GreySkull LP scenarios.
func TestGreySkullProgression_RealisticScenarios(t *testing.T) {
	// Upper body progression (bench press, OHP)
	upperGsp, err := NewGreySkullMainLiftProgression("upper-prog", "Upper Body Main", 2.5, progression.TrainingMax)
	if err != nil {
		t.Fatalf("failed to create upper body progression: %v", err)
	}

	// Lower body progression (squat, deadlift)
	lowerGsp, err := NewGreySkullProgression("lower-prog", "Lower Body Main", 5.0, 5, 10, 0.10, progression.TrainingMax)
	if err != nil {
		t.Fatalf("failed to create lower body progression: %v", err)
	}

	ctx := context.Background()

	t.Run("beginner bench press scenario - standard progress", func(t *testing.T) {
		// Beginner hits 5 reps on 135lb AMRAP set
		params := progression.ProgressionContext{
			UserID:       "beginner-1",
			LiftID:       "bench-uuid",
			MaxType:      progression.TrainingMax,
			CurrentValue: 135,
			TriggerEvent: progression.TriggerEvent{
				Type:          progression.TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(5),
				IsAMRAP:       true,
			},
		}

		result, err := upperGsp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.NewValue != 137.5 {
			t.Errorf("expected next session at 137.5lb, got %f", result.NewValue)
		}
	})

	t.Run("beginner squat scenario - great AMRAP performance", func(t *testing.T) {
		// Beginner crushes it with 12 reps on 185lb squat AMRAP
		params := progression.ProgressionContext{
			UserID:       "beginner-1",
			LiftID:       "squat-uuid",
			MaxType:      progression.TrainingMax,
			CurrentValue: 185,
			TriggerEvent: progression.TriggerEvent{
				Type:          progression.TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(12),
				IsAMRAP:       true,
			},
		}

		result, err := lowerGsp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Double increment: 5.0 * 2 = 10.0
		if result.NewValue != 195 {
			t.Errorf("expected next session at 195lb (double increment), got %f", result.NewValue)
		}
	})

	t.Run("intermediate stall scenario - deload triggered", func(t *testing.T) {
		// Intermediate lifter fails on OHP, only getting 4 reps at 135lb
		params := progression.ProgressionContext{
			UserID:       "intermediate-1",
			LiftID:       "ohp-uuid",
			MaxType:      progression.TrainingMax,
			CurrentValue: 135,
			TriggerEvent: progression.TriggerEvent{
				Type:          progression.TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(4),
				IsAMRAP:       true,
			},
		}

		result, err := upperGsp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// 10% deload: 135 * 0.10 = 13.5, new = 121.5
		if result.NewValue != 121.5 {
			t.Errorf("expected deload to 121.5lb, got %f", result.NewValue)
		}
	})
}

// TestGreySkullProgression_OneRM tests progression with ONE_RM max type.
func TestGreySkullProgression_OneRM(t *testing.T) {
	gsp := &GreySkullProgression{
		ID:              "prog-1",
		Name:            "1RM GreySkull Progression",
		WeightIncrement: 2.5,
		MinReps:         5,
		DoubleThreshold: 10,
		DeloadPercent:   0.10,
		MaxTypeValue:    progression.OneRM,
	}

	ctx := context.Background()

	t.Run("applies to ONE_RM", func(t *testing.T) {
		params := progression.ProgressionContext{
			UserID:       "user-1",
			LiftID:       "bench-uuid",
			MaxType:      progression.OneRM,
			CurrentValue: 225,
			TriggerEvent: progression.TriggerEvent{
				Type:          progression.TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(7),
				IsAMRAP:       true,
			},
		}

		result, err := gsp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied=true, reason: %s", result.Reason)
		}
		if result.NewValue != 227.5 {
			t.Errorf("expected NewValue 227.5, got %f", result.NewValue)
		}
	})

	t.Run("does not apply to TRAINING_MAX", func(t *testing.T) {
		params := progression.ProgressionContext{
			UserID:       "user-1",
			LiftID:       "bench-uuid",
			MaxType:      progression.TrainingMax, // Mismatch
			CurrentValue: 200,
			TriggerEvent: progression.TriggerEvent{
				Type:          progression.TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(8),
				IsAMRAP:       true,
			},
		}

		result, err := gsp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied=false for max type mismatch")
		}
	})
}

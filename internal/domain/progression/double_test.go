package progression

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

// Helper function for creating int pointers (local to this test file)
func doubleIntPtr(i int) *int {
	return &i
}

// TestDoubleProgression_Type tests that DoubleProgression returns correct type.
func TestDoubleProgression_Type(t *testing.T) {
	dp := &DoubleProgression{
		ID:               "prog-1",
		Name:             "Test Progression",
		WeightIncrement:  5.0,
		MaxTypeValue:     TrainingMax,
		TriggerTypeValue: TriggerAfterSet,
	}
	if dp.Type() != TypeDouble {
		t.Errorf("expected %s, got %s", TypeDouble, dp.Type())
	}
}

// TestDoubleProgression_TriggerType tests that DoubleProgression returns correct trigger type.
func TestDoubleProgression_TriggerType(t *testing.T) {
	dp := &DoubleProgression{
		ID:               "prog-1",
		Name:             "Test",
		WeightIncrement:  5.0,
		MaxTypeValue:     TrainingMax,
		TriggerTypeValue: TriggerAfterSet,
	}
	if dp.TriggerType() != TriggerAfterSet {
		t.Errorf("expected %s, got %s", TriggerAfterSet, dp.TriggerType())
	}
}

// TestDoubleProgression_Validate tests DoubleProgression validation.
func TestDoubleProgression_Validate(t *testing.T) {
	tests := []struct {
		name    string
		dp      DoubleProgression
		wantErr bool
		errType error
	}{
		{
			name: "valid progression",
			dp: DoubleProgression{
				ID:               "prog-1",
				Name:             "Double Progression",
				WeightIncrement:  5.0,
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: TriggerAfterSet,
			},
			wantErr: false,
		},
		{
			name: "valid progression with ONE_RM",
			dp: DoubleProgression{
				ID:               "prog-2",
				Name:             "1RM Double Progression",
				WeightIncrement:  2.5,
				MaxTypeValue:     OneRM,
				TriggerTypeValue: TriggerAfterSet,
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			dp: DoubleProgression{
				ID:               "",
				Name:             "Test",
				WeightIncrement:  5.0,
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: TriggerAfterSet,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "missing name",
			dp: DoubleProgression{
				ID:               "prog-1",
				Name:             "",
				WeightIncrement:  5.0,
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: TriggerAfterSet,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "zero weight increment",
			dp: DoubleProgression{
				ID:               "prog-1",
				Name:             "Test",
				WeightIncrement:  0,
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: TriggerAfterSet,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "negative weight increment",
			dp: DoubleProgression{
				ID:               "prog-1",
				Name:             "Test",
				WeightIncrement:  -5.0,
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: TriggerAfterSet,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "missing max type",
			dp: DoubleProgression{
				ID:               "prog-1",
				Name:             "Test",
				WeightIncrement:  5.0,
				MaxTypeValue:     "",
				TriggerTypeValue: TriggerAfterSet,
			},
			wantErr: true,
			errType: ErrUnknownMaxType,
		},
		{
			name: "invalid max type",
			dp: DoubleProgression{
				ID:               "prog-1",
				Name:             "Test",
				WeightIncrement:  5.0,
				MaxTypeValue:     "INVALID",
				TriggerTypeValue: TriggerAfterSet,
			},
			wantErr: true,
			errType: ErrUnknownMaxType,
		},
		{
			name: "wrong trigger type (AFTER_SESSION)",
			dp: DoubleProgression{
				ID:               "prog-1",
				Name:             "Test",
				WeightIncrement:  5.0,
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: TriggerAfterSession,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "wrong trigger type (ON_FAILURE)",
			dp: DoubleProgression{
				ID:               "prog-1",
				Name:             "Test",
				WeightIncrement:  5.0,
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: TriggerOnFailure,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "missing trigger type",
			dp: DoubleProgression{
				ID:               "prog-1",
				Name:             "Test",
				WeightIncrement:  5.0,
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: "",
			},
			wantErr: true,
			errType: ErrUnknownTriggerType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dp.Validate()
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

// TestNewDoubleProgression tests the factory function.
func TestNewDoubleProgression(t *testing.T) {
	t.Run("valid parameters", func(t *testing.T) {
		dp, err := NewDoubleProgression("prog-1", "Test Double", 5.0, TrainingMax, TriggerAfterSet)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if dp.ID != "prog-1" {
			t.Errorf("expected ID 'prog-1', got '%s'", dp.ID)
		}
		if dp.Name != "Test Double" {
			t.Errorf("expected Name 'Test Double', got '%s'", dp.Name)
		}
		if dp.WeightIncrement != 5.0 {
			t.Errorf("expected WeightIncrement 5.0, got %f", dp.WeightIncrement)
		}
	})

	t.Run("invalid parameters - empty ID", func(t *testing.T) {
		_, err := NewDoubleProgression("", "Test", 5.0, TrainingMax, TriggerAfterSet)
		if err == nil {
			t.Error("expected error for empty ID")
		}
	})

	t.Run("invalid parameters - zero increment", func(t *testing.T) {
		_, err := NewDoubleProgression("prog-1", "Test", 0, TrainingMax, TriggerAfterSet)
		if err == nil {
			t.Error("expected error for zero increment")
		}
	})
}

// TestDoubleProgression_Apply tests the Apply method.
func TestDoubleProgression_Apply(t *testing.T) {
	dp := &DoubleProgression{
		ID:               "prog-1",
		Name:             "3x8-12 Progression",
		WeightIncrement:  5.0,
		MaxTypeValue:     TrainingMax,
		TriggerTypeValue: TriggerAfterSet,
	}

	ctx := context.Background()

	t.Run("applies when reps equal max reps (ceiling)", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "bench-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 100,
			TriggerEvent: TriggerEvent{
				Type:          TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: doubleIntPtr(12),
				MaxReps:       doubleIntPtr(12), // Ceiling is 12
			},
		}

		result, err := dp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true, reason: %s", result.Reason)
		}
		if result.Delta != 5.0 {
			t.Errorf("expected Delta 5.0, got %f", result.Delta)
		}
		if result.NewValue != 105 {
			t.Errorf("expected NewValue 105, got %f", result.NewValue)
		}
	})

	t.Run("applies when reps exceed max reps", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "bench-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 100,
			TriggerEvent: TriggerEvent{
				Type:          TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: doubleIntPtr(14), // Exceeded ceiling
				MaxReps:       doubleIntPtr(12),
			},
		}

		result, err := dp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true, reason: %s", result.Reason)
		}
		if result.Delta != 5.0 {
			t.Errorf("expected Delta 5.0, got %f", result.Delta)
		}
		if result.NewValue != 105 {
			t.Errorf("expected NewValue 105, got %f", result.NewValue)
		}
	})

	t.Run("does not apply when below ceiling", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "bench-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 100,
			TriggerEvent: TriggerEvent{
				Type:          TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: doubleIntPtr(10), // Below ceiling
				MaxReps:       doubleIntPtr(12),
			},
		}

		result, err := dp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied to be false when below ceiling")
		}
		if result.Delta != 0 {
			t.Errorf("expected Delta 0, got %f", result.Delta)
		}
		if result.NewValue != 100 {
			t.Errorf("expected NewValue unchanged at 100, got %f", result.NewValue)
		}
		if result.Reason == "" {
			t.Error("expected Reason to be set")
		}
	})

	t.Run("does not apply when at minimum reps", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "bench-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 100,
			TriggerEvent: TriggerEvent{
				Type:          TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: doubleIntPtr(8), // At minimum
				MaxReps:       doubleIntPtr(12),
			},
		}

		result, err := dp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied to be false when at minimum")
		}
		if result.Delta != 0 {
			t.Errorf("expected Delta 0, got %f", result.Delta)
		}
	})

	t.Run("trigger type mismatch", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "bench-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 100,
			TriggerEvent: TriggerEvent{
				Type:      TriggerAfterSession, // Wrong trigger type
				Timestamp: time.Now(),
			},
		}

		result, err := dp.Apply(ctx, params)
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
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "bench-uuid",
			MaxType:      OneRM, // Wrong max type
			CurrentValue: 120,
			TriggerEvent: TriggerEvent{
				Type:          TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: doubleIntPtr(12),
				MaxReps:       doubleIntPtr(12),
			},
		}

		result, err := dp.Apply(ctx, params)
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

	t.Run("reps not provided", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "bench-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 100,
			TriggerEvent: TriggerEvent{
				Type:      TriggerAfterSet,
				Timestamp: time.Now(),
				MaxReps:   doubleIntPtr(12),
				// RepsPerformed is nil
			},
		}

		result, err := dp.Apply(ctx, params)
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

	t.Run("max reps (ceiling) not provided", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "bench-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 100,
			TriggerEvent: TriggerEvent{
				Type:          TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: doubleIntPtr(12),
				// MaxReps is nil
			},
		}

		result, err := dp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied to be false when max reps not provided")
		}
		if result.Reason == "" {
			t.Error("expected Reason to be set")
		}
	})

	t.Run("invalid context - missing userID", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "", // Missing userID
			LiftID:       "bench-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 100,
			TriggerEvent: TriggerEvent{
				Type:          TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: doubleIntPtr(12),
				MaxReps:       doubleIntPtr(12),
			},
		}

		_, err := dp.Apply(ctx, params)
		if err == nil {
			t.Error("expected error for invalid context")
		}
	})
}

// TestDoubleProgression_Apply_OneRM tests progression with ONE_RM max type.
func TestDoubleProgression_Apply_OneRM(t *testing.T) {
	dp := &DoubleProgression{
		ID:               "prog-1",
		Name:             "1RM Double Progression",
		WeightIncrement:  2.5,
		MaxTypeValue:     OneRM,
		TriggerTypeValue: TriggerAfterSet,
	}

	ctx := context.Background()

	t.Run("applies to ONE_RM", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-1",
			LiftID:       "squat-uuid",
			MaxType:      OneRM,
			CurrentValue: 315,
			TriggerEvent: TriggerEvent{
				Type:          TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: doubleIntPtr(12),
				MaxReps:       doubleIntPtr(12),
			},
		}

		result, err := dp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied=true, reason: %s", result.Reason)
		}
		if result.NewValue != 317.5 {
			t.Errorf("expected NewValue 317.5, got %f", result.NewValue)
		}
	})

	t.Run("does not apply to TRAINING_MAX", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-1",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax, // Mismatch
			CurrentValue: 285,
			TriggerEvent: TriggerEvent{
				Type:          TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: doubleIntPtr(12),
				MaxReps:       doubleIntPtr(12),
			},
		}

		result, err := dp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied=false for max type mismatch")
		}
	})
}

// TestDoubleProgression_Apply_EdgeCases tests edge cases in Apply.
func TestDoubleProgression_Apply_EdgeCases(t *testing.T) {
	dp := &DoubleProgression{
		ID:               "prog-1",
		Name:             "Test",
		WeightIncrement:  5.0,
		MaxTypeValue:     TrainingMax,
		TriggerTypeValue: TriggerAfterSet,
	}

	ctx := context.Background()

	t.Run("zero reps performed", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-1",
			LiftID:       "lift-1",
			MaxType:      TrainingMax,
			CurrentValue: 100,
			TriggerEvent: TriggerEvent{
				Type:          TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: doubleIntPtr(0),
				MaxReps:       doubleIntPtr(12),
			},
		}

		result, err := dp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied=false for zero reps")
		}
	})

	t.Run("ceiling of 1 rep", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-1",
			LiftID:       "lift-1",
			MaxType:      TrainingMax,
			CurrentValue: 100,
			TriggerEvent: TriggerEvent{
				Type:          TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: doubleIntPtr(1),
				MaxReps:       doubleIntPtr(1), // Ceiling of 1
			},
		}

		result, err := dp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied=true when at ceiling of 1, reason: %s", result.Reason)
		}
	})

	t.Run("very large rep count", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-1",
			LiftID:       "lift-1",
			MaxType:      TrainingMax,
			CurrentValue: 50,
			TriggerEvent: TriggerEvent{
				Type:          TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: doubleIntPtr(50),
				MaxReps:       doubleIntPtr(20),
			},
		}

		result, err := dp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied=true when exceeding ceiling, reason: %s", result.Reason)
		}
	})

	t.Run("fractional weight increment", func(t *testing.T) {
		dpFractional := &DoubleProgression{
			ID:               "prog-frac",
			Name:             "Fractional",
			WeightIncrement:  1.25, // Fractional increment
			MaxTypeValue:     TrainingMax,
			TriggerTypeValue: TriggerAfterSet,
		}

		params := ProgressionContext{
			UserID:       "user-1",
			LiftID:       "lift-1",
			MaxType:      TrainingMax,
			CurrentValue: 100,
			TriggerEvent: TriggerEvent{
				Type:          TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: doubleIntPtr(12),
				MaxReps:       doubleIntPtr(12),
			},
		}

		result, err := dpFractional.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.NewValue != 101.25 {
			t.Errorf("expected NewValue 101.25, got %f", result.NewValue)
		}
	})
}

// TestDoubleProgression_Apply_AppliedAt tests that AppliedAt is set correctly.
func TestDoubleProgression_Apply_AppliedAt(t *testing.T) {
	dp := &DoubleProgression{
		ID:               "prog-1",
		Name:             "Test",
		WeightIncrement:  5.0,
		MaxTypeValue:     TrainingMax,
		TriggerTypeValue: TriggerAfterSet,
	}

	before := time.Now()

	params := ProgressionContext{
		UserID:       "user-1",
		LiftID:       "lift-1",
		MaxType:      TrainingMax,
		CurrentValue: 100,
		TriggerEvent: TriggerEvent{
			Type:          TriggerAfterSet,
			Timestamp:     time.Now(),
			RepsPerformed: doubleIntPtr(12),
			MaxReps:       doubleIntPtr(12),
		},
	}

	result, err := dp.Apply(context.Background(), params)
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

// TestDoubleProgression_JSON tests JSON serialization roundtrip.
func TestDoubleProgression_JSON(t *testing.T) {
	dp := &DoubleProgression{
		ID:               "prog-123",
		Name:             "Test Double Progression",
		WeightIncrement:  5.0,
		MaxTypeValue:     TrainingMax,
		TriggerTypeValue: TriggerAfterSet,
	}

	// Marshal
	data, err := json.Marshal(dp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Verify JSON structure
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if parsed["type"] != string(TypeDouble) {
		t.Errorf("expected type %s, got %v", TypeDouble, parsed["type"])
	}
	if parsed["id"] != "prog-123" {
		t.Errorf("expected id 'prog-123', got %v", parsed["id"])
	}
	if parsed["name"] != "Test Double Progression" {
		t.Errorf("expected name 'Test Double Progression', got %v", parsed["name"])
	}
	if parsed["weightIncrement"] != 5.0 {
		t.Errorf("expected weightIncrement 5.0, got %v", parsed["weightIncrement"])
	}
	if parsed["maxType"] != string(TrainingMax) {
		t.Errorf("expected maxType %s, got %v", TrainingMax, parsed["maxType"])
	}
	if parsed["triggerType"] != string(TriggerAfterSet) {
		t.Errorf("expected triggerType %s, got %v", TriggerAfterSet, parsed["triggerType"])
	}

	// Unmarshal back
	var restored DoubleProgression
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if restored.ID != dp.ID {
		t.Errorf("ID mismatch: expected %s, got %s", dp.ID, restored.ID)
	}
	if restored.Name != dp.Name {
		t.Errorf("Name mismatch: expected %s, got %s", dp.Name, restored.Name)
	}
	if restored.WeightIncrement != dp.WeightIncrement {
		t.Errorf("WeightIncrement mismatch: expected %f, got %f", dp.WeightIncrement, restored.WeightIncrement)
	}
	if restored.MaxTypeValue != dp.MaxTypeValue {
		t.Errorf("MaxTypeValue mismatch: expected %s, got %s", dp.MaxTypeValue, restored.MaxTypeValue)
	}
	if restored.TriggerTypeValue != dp.TriggerTypeValue {
		t.Errorf("TriggerTypeValue mismatch: expected %s, got %s", dp.TriggerTypeValue, restored.TriggerTypeValue)
	}
}

// TestUnmarshalDoubleProgression tests deserialization function.
func TestUnmarshalDoubleProgression(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		jsonData := []byte(`{
			"id": "prog-1",
			"name": "Test",
			"weightIncrement": 5.0,
			"maxType": "TRAINING_MAX",
			"triggerType": "AFTER_SET"
		}`)

		progression, err := UnmarshalDoubleProgression(jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		dp, ok := progression.(*DoubleProgression)
		if !ok {
			t.Fatalf("expected *DoubleProgression, got %T", progression)
		}
		if dp.ID != "prog-1" {
			t.Errorf("expected ID 'prog-1', got '%s'", dp.ID)
		}
		if dp.Type() != TypeDouble {
			t.Errorf("expected type %s, got %s", TypeDouble, dp.Type())
		}
		if dp.WeightIncrement != 5.0 {
			t.Errorf("expected WeightIncrement 5.0, got %f", dp.WeightIncrement)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		jsonData := []byte(`{invalid}`)
		_, err := UnmarshalDoubleProgression(jsonData)
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("invalid progression data - empty ID", func(t *testing.T) {
		jsonData := []byte(`{
			"id": "",
			"name": "Test",
			"weightIncrement": 5.0,
			"maxType": "TRAINING_MAX",
			"triggerType": "AFTER_SET"
		}`)

		_, err := UnmarshalDoubleProgression(jsonData)
		if err == nil {
			t.Error("expected error for invalid progression data")
		}
	})

	t.Run("invalid progression data - zero increment", func(t *testing.T) {
		jsonData := []byte(`{
			"id": "prog-1",
			"name": "Test",
			"weightIncrement": 0,
			"maxType": "TRAINING_MAX",
			"triggerType": "AFTER_SET"
		}`)

		_, err := UnmarshalDoubleProgression(jsonData)
		if err == nil {
			t.Error("expected error for zero increment")
		}
	})
}

// TestRegisterDoubleProgression tests factory registration.
func TestRegisterDoubleProgression(t *testing.T) {
	factory := NewProgressionFactory()

	// Verify not registered initially
	if factory.IsRegistered(TypeDouble) {
		t.Error("TypeDouble should not be registered initially")
	}

	// Register
	RegisterDoubleProgression(factory)

	// Verify registered
	if !factory.IsRegistered(TypeDouble) {
		t.Error("TypeDouble should be registered after calling RegisterDoubleProgression")
	}

	// Create from factory
	jsonData := []byte(`{
		"type": "DOUBLE_PROGRESSION",
		"id": "prog-1",
		"name": "Factory Test",
		"weightIncrement": 5.0,
		"maxType": "TRAINING_MAX",
		"triggerType": "AFTER_SET"
	}`)

	progression, err := factory.CreateFromJSON(jsonData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if progression.Type() != TypeDouble {
		t.Errorf("expected type %s, got %s", TypeDouble, progression.Type())
	}
}

// TestDoubleProgression_Interface verifies that DoubleProgression implements Progression.
func TestDoubleProgression_Interface(t *testing.T) {
	var _ Progression = &DoubleProgression{}
}

// TestDoubleProgression_RealisticExample tests a realistic bodybuilding-style configuration.
func TestDoubleProgression_RealisticExample(t *testing.T) {
	// Typical 3x8-12 scheme with 5lb increments
	dp := &DoubleProgression{
		ID:               "bb-prog",
		Name:             "Bodybuilding 3x8-12",
		WeightIncrement:  5.0,
		MaxTypeValue:     TrainingMax,
		TriggerTypeValue: TriggerAfterSet,
	}

	ctx := context.Background()

	testCases := []struct {
		reps            int
		ceiling         int
		currentValue    float64
		expectedApplied bool
		expectedDelta   float64
		description     string
	}{
		{reps: 8, ceiling: 12, currentValue: 100, expectedApplied: false, expectedDelta: 0, description: "just started, at 8 reps"},
		{reps: 9, ceiling: 12, currentValue: 100, expectedApplied: false, expectedDelta: 0, description: "progressing, at 9 reps"},
		{reps: 10, ceiling: 12, currentValue: 100, expectedApplied: false, expectedDelta: 0, description: "progressing, at 10 reps"},
		{reps: 11, ceiling: 12, currentValue: 100, expectedApplied: false, expectedDelta: 0, description: "almost there, at 11 reps"},
		{reps: 12, ceiling: 12, currentValue: 100, expectedApplied: true, expectedDelta: 5.0, description: "hit ceiling, at 12 reps"},
		{reps: 13, ceiling: 12, currentValue: 105, expectedApplied: true, expectedDelta: 5.0, description: "exceeded ceiling, at 13 reps"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			params := ProgressionContext{
				UserID:       "user-1",
				LiftID:       "curls-uuid",
				MaxType:      TrainingMax,
				CurrentValue: tc.currentValue,
				TriggerEvent: TriggerEvent{
					Type:          TriggerAfterSet,
					Timestamp:     time.Now(),
					RepsPerformed: doubleIntPtr(tc.reps),
					MaxReps:       doubleIntPtr(tc.ceiling),
				},
			}

			result, err := dp.Apply(ctx, params)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Applied != tc.expectedApplied {
				if tc.expectedApplied {
					t.Errorf("expected Applied=true for %s, reason: %s", tc.description, result.Reason)
				} else {
					t.Errorf("expected Applied=false for %s", tc.description)
				}
			}

			if result.Delta != tc.expectedDelta {
				t.Errorf("expected Delta %f, got %f", tc.expectedDelta, result.Delta)
			}
		})
	}
}

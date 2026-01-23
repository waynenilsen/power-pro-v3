package progression

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

// TestAMRAPProgression_Type tests that AMRAPProgression returns correct type.
func TestAMRAPProgression_Type(t *testing.T) {
	ap := &AMRAPProgression{
		ID:               "prog-1",
		Name:             "Test Progression",
		MaxTypeValue:     TrainingMax,
		TriggerTypeValue: TriggerAfterSet,
		Thresholds: []RepsThreshold{
			{MinReps: 2, Increment: 5.0},
		},
	}
	if ap.Type() != TypeAMRAP {
		t.Errorf("expected %s, got %s", TypeAMRAP, ap.Type())
	}
}

// TestAMRAPProgression_TriggerType tests that AMRAPProgression returns correct trigger type.
func TestAMRAPProgression_TriggerType(t *testing.T) {
	ap := &AMRAPProgression{
		ID:               "prog-1",
		Name:             "Test",
		MaxTypeValue:     TrainingMax,
		TriggerTypeValue: TriggerAfterSet,
		Thresholds: []RepsThreshold{
			{MinReps: 2, Increment: 5.0},
		},
	}
	if ap.TriggerType() != TriggerAfterSet {
		t.Errorf("expected %s, got %s", TriggerAfterSet, ap.TriggerType())
	}
}

// TestAMRAPProgression_Validate tests AMRAPProgression validation.
func TestAMRAPProgression_Validate(t *testing.T) {
	tests := []struct {
		name    string
		ap      AMRAPProgression
		wantErr bool
		errType error
	}{
		{
			name: "valid progression with single threshold",
			ap: AMRAPProgression{
				ID:               "prog-1",
				Name:             "nSuns Progression",
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: TriggerAfterSet,
				Thresholds: []RepsThreshold{
					{MinReps: 2, Increment: 5.0},
				},
			},
			wantErr: false,
		},
		{
			name: "valid progression with multiple thresholds",
			ap: AMRAPProgression{
				ID:               "prog-2",
				Name:             "nSuns Multi-tier",
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: TriggerAfterSet,
				Thresholds: []RepsThreshold{
					{MinReps: 2, Increment: 5.0},
					{MinReps: 4, Increment: 10.0},
					{MinReps: 6, Increment: 15.0},
				},
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			ap: AMRAPProgression{
				ID:               "",
				Name:             "Test",
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: TriggerAfterSet,
				Thresholds: []RepsThreshold{
					{MinReps: 2, Increment: 5.0},
				},
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "missing name",
			ap: AMRAPProgression{
				ID:               "prog-1",
				Name:             "",
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: TriggerAfterSet,
				Thresholds: []RepsThreshold{
					{MinReps: 2, Increment: 5.0},
				},
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "missing max type",
			ap: AMRAPProgression{
				ID:               "prog-1",
				Name:             "Test",
				MaxTypeValue:     "",
				TriggerTypeValue: TriggerAfterSet,
				Thresholds: []RepsThreshold{
					{MinReps: 2, Increment: 5.0},
				},
			},
			wantErr: true,
			errType: ErrUnknownMaxType,
		},
		{
			name: "invalid max type",
			ap: AMRAPProgression{
				ID:               "prog-1",
				Name:             "Test",
				MaxTypeValue:     "INVALID",
				TriggerTypeValue: TriggerAfterSet,
				Thresholds: []RepsThreshold{
					{MinReps: 2, Increment: 5.0},
				},
			},
			wantErr: true,
			errType: ErrUnknownMaxType,
		},
		{
			name: "wrong trigger type (AFTER_SESSION)",
			ap: AMRAPProgression{
				ID:               "prog-1",
				Name:             "Test",
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: TriggerAfterSession,
				Thresholds: []RepsThreshold{
					{MinReps: 2, Increment: 5.0},
				},
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "empty thresholds",
			ap: AMRAPProgression{
				ID:               "prog-1",
				Name:             "Test",
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: TriggerAfterSet,
				Thresholds:       []RepsThreshold{},
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "negative minReps",
			ap: AMRAPProgression{
				ID:               "prog-1",
				Name:             "Test",
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: TriggerAfterSet,
				Thresholds: []RepsThreshold{
					{MinReps: -1, Increment: 5.0},
				},
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "zero increment",
			ap: AMRAPProgression{
				ID:               "prog-1",
				Name:             "Test",
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: TriggerAfterSet,
				Thresholds: []RepsThreshold{
					{MinReps: 2, Increment: 0},
				},
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "negative increment",
			ap: AMRAPProgression{
				ID:               "prog-1",
				Name:             "Test",
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: TriggerAfterSet,
				Thresholds: []RepsThreshold{
					{MinReps: 2, Increment: -5.0},
				},
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "thresholds not sorted ascending",
			ap: AMRAPProgression{
				ID:               "prog-1",
				Name:             "Test",
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: TriggerAfterSet,
				Thresholds: []RepsThreshold{
					{MinReps: 4, Increment: 10.0},
					{MinReps: 2, Increment: 5.0}, // Out of order
				},
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "thresholds with duplicate minReps",
			ap: AMRAPProgression{
				ID:               "prog-1",
				Name:             "Test",
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: TriggerAfterSet,
				Thresholds: []RepsThreshold{
					{MinReps: 2, Increment: 5.0},
					{MinReps: 2, Increment: 10.0}, // Duplicate minReps
				},
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ap.Validate()
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

// TestNewAMRAPProgression tests the factory function.
func TestNewAMRAPProgression(t *testing.T) {
	t.Run("valid parameters", func(t *testing.T) {
		thresholds := []RepsThreshold{
			{MinReps: 2, Increment: 5.0},
			{MinReps: 4, Increment: 10.0},
		}
		ap, err := NewAMRAPProgression("prog-1", "nSuns", TrainingMax, TriggerAfterSet, thresholds)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ap.ID != "prog-1" {
			t.Errorf("expected ID 'prog-1', got '%s'", ap.ID)
		}
		if ap.Name != "nSuns" {
			t.Errorf("expected Name 'nSuns', got '%s'", ap.Name)
		}
		if len(ap.Thresholds) != 2 {
			t.Errorf("expected 2 thresholds, got %d", len(ap.Thresholds))
		}
	})

	t.Run("invalid parameters", func(t *testing.T) {
		_, err := NewAMRAPProgression("", "Test", TrainingMax, TriggerAfterSet, []RepsThreshold{
			{MinReps: 2, Increment: 5.0},
		})
		if err == nil {
			t.Error("expected error for empty ID")
		}
	})
}

// TestAMRAPProgression_Apply tests the Apply method.
func TestAMRAPProgression_Apply(t *testing.T) {
	ap := &AMRAPProgression{
		ID:               "prog-1",
		Name:             "nSuns Progression",
		MaxTypeValue:     TrainingMax,
		TriggerTypeValue: TriggerAfterSet,
		Thresholds: []RepsThreshold{
			{MinReps: 2, Increment: 5.0},
			{MinReps: 4, Increment: 10.0},
			{MinReps: 6, Increment: 15.0},
		},
	}

	ctx := context.Background()

	t.Run("hits lowest threshold (2 reps)", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:          TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(2),
				IsAMRAP:       true,
			},
		}

		result, err := ap.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true, reason: %s", result.Reason)
		}
		if result.Delta != 5.0 {
			t.Errorf("expected Delta 5.0, got %f", result.Delta)
		}
		if result.NewValue != 305 {
			t.Errorf("expected NewValue 305, got %f", result.NewValue)
		}
	})

	t.Run("hits lowest threshold (3 reps)", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:          TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(3),
				IsAMRAP:       true,
			},
		}

		result, err := ap.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true, reason: %s", result.Reason)
		}
		if result.Delta != 5.0 {
			t.Errorf("expected Delta 5.0 (still in 2-3 range), got %f", result.Delta)
		}
	})

	t.Run("hits middle threshold (4 reps)", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:          TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(4),
				IsAMRAP:       true,
			},
		}

		result, err := ap.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true, reason: %s", result.Reason)
		}
		if result.Delta != 10.0 {
			t.Errorf("expected Delta 10.0, got %f", result.Delta)
		}
		if result.NewValue != 310 {
			t.Errorf("expected NewValue 310, got %f", result.NewValue)
		}
	})

	t.Run("hits highest threshold (6 reps)", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:          TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(6),
				IsAMRAP:       true,
			},
		}

		result, err := ap.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true, reason: %s", result.Reason)
		}
		if result.Delta != 15.0 {
			t.Errorf("expected Delta 15.0, got %f", result.Delta)
		}
		if result.NewValue != 315 {
			t.Errorf("expected NewValue 315, got %f", result.NewValue)
		}
	})

	t.Run("exceeds highest threshold (10 reps)", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:          TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(10),
				IsAMRAP:       true,
			},
		}

		result, err := ap.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true, reason: %s", result.Reason)
		}
		if result.Delta != 15.0 {
			t.Errorf("expected Delta 15.0 (highest threshold), got %f", result.Delta)
		}
	})

	t.Run("below minimum threshold (1 rep)", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:          TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(1),
				IsAMRAP:       true,
			},
		}

		result, err := ap.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied to be false when below minimum threshold")
		}
		if result.Delta != 0 {
			t.Errorf("expected Delta 0, got %f", result.Delta)
		}
		if result.NewValue != 300 {
			t.Errorf("expected NewValue unchanged at 300, got %f", result.NewValue)
		}
		if result.Reason == "" {
			t.Error("expected Reason to be set")
		}
	})

	t.Run("zero reps", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:          TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(0),
				IsAMRAP:       true,
			},
		}

		result, err := ap.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied to be false for zero reps")
		}
	})

	t.Run("trigger type mismatch", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:      TriggerAfterSession, // Wrong trigger type
				Timestamp: time.Now(),
			},
		}

		result, err := ap.Apply(ctx, params)
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
			LiftID:       "squat-uuid",
			MaxType:      OneRM, // Wrong max type
			CurrentValue: 335,
			TriggerEvent: TriggerEvent{
				Type:          TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(5),
				IsAMRAP:       true,
			},
		}

		result, err := ap.Apply(ctx, params)
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
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:          TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(5),
				IsAMRAP:       false, // Not an AMRAP set
			},
		}

		result, err := ap.Apply(ctx, params)
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
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:      TriggerAfterSet,
				Timestamp: time.Now(),
				IsAMRAP:   true,
				// RepsPerformed is nil
			},
		}

		result, err := ap.Apply(ctx, params)
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

	t.Run("invalid context", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "", // Missing userID
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:          TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(5),
				IsAMRAP:       true,
			},
		}

		_, err := ap.Apply(ctx, params)
		if err == nil {
			t.Error("expected error for invalid context")
		}
	})
}

// TestAMRAPProgression_Apply_SingleThreshold tests Apply with a single threshold.
func TestAMRAPProgression_Apply_SingleThreshold(t *testing.T) {
	ap := &AMRAPProgression{
		ID:               "prog-1",
		Name:             "Simple AMRAP",
		MaxTypeValue:     TrainingMax,
		TriggerTypeValue: TriggerAfterSet,
		Thresholds: []RepsThreshold{
			{MinReps: 3, Increment: 5.0},
		},
	}

	ctx := context.Background()

	t.Run("meets threshold", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-1",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 200,
			TriggerEvent: TriggerEvent{
				Type:          TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(3),
				IsAMRAP:       true,
			},
		}

		result, err := ap.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied=true, reason: %s", result.Reason)
		}
		if result.NewValue != 205 {
			t.Errorf("expected NewValue 205, got %f", result.NewValue)
		}
	})

	t.Run("below threshold", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-1",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 200,
			TriggerEvent: TriggerEvent{
				Type:          TriggerAfterSet,
				Timestamp:     time.Now(),
				RepsPerformed: intPtr(2),
				IsAMRAP:       true,
			},
		}

		result, err := ap.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied=false when below threshold")
		}
		if result.NewValue != 200 {
			t.Errorf("expected NewValue unchanged at 200, got %f", result.NewValue)
		}
	})
}

// TestAMRAPProgression_Apply_ZeroMinReps tests Apply with minReps=0 threshold.
func TestAMRAPProgression_Apply_ZeroMinReps(t *testing.T) {
	ap := &AMRAPProgression{
		ID:               "prog-1",
		Name:             "Always Progress",
		MaxTypeValue:     TrainingMax,
		TriggerTypeValue: TriggerAfterSet,
		Thresholds: []RepsThreshold{
			{MinReps: 0, Increment: 2.5}, // Always applies if AMRAP
		},
	}

	ctx := context.Background()

	params := ProgressionContext{
		UserID:       "user-1",
		LiftID:       "squat-uuid",
		MaxType:      TrainingMax,
		CurrentValue: 100,
		TriggerEvent: TriggerEvent{
			Type:          TriggerAfterSet,
			Timestamp:     time.Now(),
			RepsPerformed: intPtr(0),
			IsAMRAP:       true,
		},
	}

	result, err := ap.Apply(ctx, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Applied {
		t.Errorf("expected Applied=true with minReps=0, reason: %s", result.Reason)
	}
	if result.NewValue != 102.5 {
		t.Errorf("expected NewValue 102.5, got %f", result.NewValue)
	}
}

// TestAMRAPProgression_JSON tests JSON serialization roundtrip.
func TestAMRAPProgression_JSON(t *testing.T) {
	ap := &AMRAPProgression{
		ID:               "prog-123",
		Name:             "Test AMRAP Progression",
		MaxTypeValue:     TrainingMax,
		TriggerTypeValue: TriggerAfterSet,
		Thresholds: []RepsThreshold{
			{MinReps: 2, Increment: 5.0},
			{MinReps: 4, Increment: 10.0},
			{MinReps: 6, Increment: 15.0},
		},
	}

	// Marshal
	data, err := json.Marshal(ap)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Verify JSON structure
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if parsed["type"] != string(TypeAMRAP) {
		t.Errorf("expected type %s, got %v", TypeAMRAP, parsed["type"])
	}
	if parsed["id"] != "prog-123" {
		t.Errorf("expected id 'prog-123', got %v", parsed["id"])
	}
	if parsed["name"] != "Test AMRAP Progression" {
		t.Errorf("expected name 'Test AMRAP Progression', got %v", parsed["name"])
	}
	if parsed["maxType"] != string(TrainingMax) {
		t.Errorf("expected maxType %s, got %v", TrainingMax, parsed["maxType"])
	}
	if parsed["triggerType"] != string(TriggerAfterSet) {
		t.Errorf("expected triggerType %s, got %v", TriggerAfterSet, parsed["triggerType"])
	}

	thresholds, ok := parsed["thresholds"].([]interface{})
	if !ok {
		t.Fatalf("expected thresholds to be an array")
	}
	if len(thresholds) != 3 {
		t.Errorf("expected 3 thresholds, got %d", len(thresholds))
	}

	// Unmarshal back
	var restored AMRAPProgression
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if restored.ID != ap.ID {
		t.Errorf("ID mismatch: expected %s, got %s", ap.ID, restored.ID)
	}
	if restored.Name != ap.Name {
		t.Errorf("Name mismatch: expected %s, got %s", ap.Name, restored.Name)
	}
	if restored.MaxTypeValue != ap.MaxTypeValue {
		t.Errorf("MaxTypeValue mismatch: expected %s, got %s", ap.MaxTypeValue, restored.MaxTypeValue)
	}
	if restored.TriggerTypeValue != ap.TriggerTypeValue {
		t.Errorf("TriggerTypeValue mismatch: expected %s, got %s", ap.TriggerTypeValue, restored.TriggerTypeValue)
	}
	if len(restored.Thresholds) != len(ap.Thresholds) {
		t.Errorf("Thresholds length mismatch: expected %d, got %d", len(ap.Thresholds), len(restored.Thresholds))
	}
	for i, th := range restored.Thresholds {
		if th.MinReps != ap.Thresholds[i].MinReps {
			t.Errorf("Threshold[%d].MinReps mismatch: expected %d, got %d", i, ap.Thresholds[i].MinReps, th.MinReps)
		}
		if th.Increment != ap.Thresholds[i].Increment {
			t.Errorf("Threshold[%d].Increment mismatch: expected %f, got %f", i, ap.Thresholds[i].Increment, th.Increment)
		}
	}
}

// TestUnmarshalAMRAPProgression tests deserialization function.
func TestUnmarshalAMRAPProgression(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		jsonData := []byte(`{
			"id": "prog-1",
			"name": "Test",
			"maxType": "TRAINING_MAX",
			"triggerType": "AFTER_SET",
			"thresholds": [
				{"minReps": 2, "increment": 5.0},
				{"minReps": 4, "increment": 10.0}
			]
		}`)

		progression, err := UnmarshalAMRAPProgression(jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		ap, ok := progression.(*AMRAPProgression)
		if !ok {
			t.Fatalf("expected *AMRAPProgression, got %T", progression)
		}
		if ap.ID != "prog-1" {
			t.Errorf("expected ID 'prog-1', got '%s'", ap.ID)
		}
		if ap.Type() != TypeAMRAP {
			t.Errorf("expected type %s, got %s", TypeAMRAP, ap.Type())
		}
		if len(ap.Thresholds) != 2 {
			t.Errorf("expected 2 thresholds, got %d", len(ap.Thresholds))
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		jsonData := []byte(`{invalid}`)
		_, err := UnmarshalAMRAPProgression(jsonData)
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("invalid progression data", func(t *testing.T) {
		jsonData := []byte(`{
			"id": "",
			"name": "Test",
			"maxType": "TRAINING_MAX",
			"triggerType": "AFTER_SET",
			"thresholds": [{"minReps": 2, "increment": 5.0}]
		}`)

		_, err := UnmarshalAMRAPProgression(jsonData)
		if err == nil {
			t.Error("expected error for invalid progression data")
		}
	})

	t.Run("sorts thresholds on unmarshal", func(t *testing.T) {
		// JSON with thresholds out of order (this tests the defensive sorting)
		jsonData := []byte(`{
			"id": "prog-1",
			"name": "Test",
			"maxType": "TRAINING_MAX",
			"triggerType": "AFTER_SET",
			"thresholds": [
				{"minReps": 6, "increment": 15.0},
				{"minReps": 2, "increment": 5.0},
				{"minReps": 4, "increment": 10.0}
			]
		}`)

		progression, err := UnmarshalAMRAPProgression(jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		ap := progression.(*AMRAPProgression)
		// After sorting: [2, 4, 6]
		if ap.Thresholds[0].MinReps != 2 {
			t.Errorf("expected first threshold minReps=2, got %d", ap.Thresholds[0].MinReps)
		}
		if ap.Thresholds[1].MinReps != 4 {
			t.Errorf("expected second threshold minReps=4, got %d", ap.Thresholds[1].MinReps)
		}
		if ap.Thresholds[2].MinReps != 6 {
			t.Errorf("expected third threshold minReps=6, got %d", ap.Thresholds[2].MinReps)
		}
	})
}

// TestRegisterAMRAPProgression tests factory registration.
func TestRegisterAMRAPProgression(t *testing.T) {
	factory := NewProgressionFactory()

	// Verify not registered initially
	if factory.IsRegistered(TypeAMRAP) {
		t.Error("TypeAMRAP should not be registered initially")
	}

	// Register
	RegisterAMRAPProgression(factory)

	// Verify registered
	if !factory.IsRegistered(TypeAMRAP) {
		t.Error("TypeAMRAP should be registered after calling RegisterAMRAPProgression")
	}

	// Create from factory
	jsonData := []byte(`{
		"type": "AMRAP_PROGRESSION",
		"id": "prog-1",
		"name": "Factory Test",
		"maxType": "TRAINING_MAX",
		"triggerType": "AFTER_SET",
		"thresholds": [
			{"minReps": 2, "increment": 5.0}
		]
	}`)

	progression, err := factory.CreateFromJSON(jsonData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if progression.Type() != TypeAMRAP {
		t.Errorf("expected type %s, got %s", TypeAMRAP, progression.Type())
	}
}

// TestAMRAPProgression_Interface verifies that AMRAPProgression implements Progression.
func TestAMRAPProgression_Interface(t *testing.T) {
	var _ Progression = &AMRAPProgression{}
}

// TestAMRAPProgression_AppliedAt tests that AppliedAt is set correctly.
func TestAMRAPProgression_AppliedAt(t *testing.T) {
	ap := &AMRAPProgression{
		ID:               "prog-1",
		Name:             "Test",
		MaxTypeValue:     TrainingMax,
		TriggerTypeValue: TriggerAfterSet,
		Thresholds: []RepsThreshold{
			{MinReps: 1, Increment: 5.0},
		},
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
			RepsPerformed: intPtr(3),
			IsAMRAP:       true,
		},
	}

	result, err := ap.Apply(context.Background(), params)
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

// TestAMRAPProgression_NSunsExample tests a realistic nSuns-style configuration.
func TestAMRAPProgression_NSunsExample(t *testing.T) {
	// nSuns 5/3/1 LP progression scheme
	ap := &AMRAPProgression{
		ID:               "nsuns-prog",
		Name:             "nSuns 5/3/1 LP",
		MaxTypeValue:     TrainingMax,
		TriggerTypeValue: TriggerAfterSet,
		Thresholds: []RepsThreshold{
			{MinReps: 2, Increment: 5.0},   // 2-3 reps: +5lb
			{MinReps: 4, Increment: 10.0},  // 4-5 reps: +10lb
			{MinReps: 6, Increment: 15.0},  // 6+ reps: +15lb
		},
	}

	ctx := context.Background()

	testCases := []struct {
		reps              int
		currentValue      float64
		expectedIncrement float64
		expectedNewValue  float64
		description       string
	}{
		{reps: 1, currentValue: 300, expectedIncrement: 0, expectedNewValue: 300, description: "failed AMRAP (1 rep)"},
		{reps: 2, currentValue: 300, expectedIncrement: 5.0, expectedNewValue: 305, description: "minimum (2 reps)"},
		{reps: 3, currentValue: 300, expectedIncrement: 5.0, expectedNewValue: 305, description: "low (3 reps)"},
		{reps: 4, currentValue: 300, expectedIncrement: 10.0, expectedNewValue: 310, description: "medium (4 reps)"},
		{reps: 5, currentValue: 300, expectedIncrement: 10.0, expectedNewValue: 310, description: "medium (5 reps)"},
		{reps: 6, currentValue: 300, expectedIncrement: 15.0, expectedNewValue: 315, description: "high (6 reps)"},
		{reps: 10, currentValue: 300, expectedIncrement: 15.0, expectedNewValue: 315, description: "excellent (10 reps)"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			params := ProgressionContext{
				UserID:       "user-1",
				LiftID:       "bench-uuid",
				MaxType:      TrainingMax,
				CurrentValue: tc.currentValue,
				TriggerEvent: TriggerEvent{
					Type:          TriggerAfterSet,
					Timestamp:     time.Now(),
					RepsPerformed: intPtr(tc.reps),
					IsAMRAP:       true,
				},
			}

			result, err := ap.Apply(ctx, params)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.expectedIncrement == 0 {
				if result.Applied {
					t.Errorf("expected Applied=false for %s", tc.description)
				}
			} else {
				if !result.Applied {
					t.Errorf("expected Applied=true for %s, reason: %s", tc.description, result.Reason)
				}
			}

			if result.Delta != tc.expectedIncrement {
				t.Errorf("expected Delta %f, got %f", tc.expectedIncrement, result.Delta)
			}
			if result.NewValue != tc.expectedNewValue {
				t.Errorf("expected NewValue %f, got %f", tc.expectedNewValue, result.NewValue)
			}
		})
	}
}

// TestAMRAPProgression_OneRM tests progression with ONE_RM max type.
func TestAMRAPProgression_OneRM(t *testing.T) {
	ap := &AMRAPProgression{
		ID:               "prog-1",
		Name:             "1RM AMRAP Progression",
		MaxTypeValue:     OneRM,
		TriggerTypeValue: TriggerAfterSet,
		Thresholds: []RepsThreshold{
			{MinReps: 2, Increment: 5.0},
		},
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
				RepsPerformed: intPtr(3),
				IsAMRAP:       true,
			},
		}

		result, err := ap.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied=true, reason: %s", result.Reason)
		}
		if result.NewValue != 320 {
			t.Errorf("expected NewValue 320, got %f", result.NewValue)
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
				RepsPerformed: intPtr(5),
				IsAMRAP:       true,
			},
		}

		result, err := ap.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied=false for max type mismatch")
		}
	})
}

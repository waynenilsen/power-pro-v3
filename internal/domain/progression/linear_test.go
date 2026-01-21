package progression

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

// TestLinearProgression_Type tests that LinearProgression returns correct type.
func TestLinearProgression_Type(t *testing.T) {
	lp := &LinearProgression{
		ID:               "prog-1",
		Name:             "Test Progression",
		Increment:        5.0,
		MaxTypeValue:     TrainingMax,
		TriggerTypeValue: TriggerAfterSession,
	}
	if lp.Type() != TypeLinear {
		t.Errorf("expected %s, got %s", TypeLinear, lp.Type())
	}
}

// TestLinearProgression_TriggerType tests that LinearProgression returns correct trigger type.
func TestLinearProgression_TriggerType(t *testing.T) {
	tests := []struct {
		name        string
		triggerType TriggerType
	}{
		{"AFTER_SESSION", TriggerAfterSession},
		{"AFTER_WEEK", TriggerAfterWeek},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lp := &LinearProgression{
				ID:               "prog-1",
				Name:             "Test",
				Increment:        5.0,
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: tt.triggerType,
			}
			if lp.TriggerType() != tt.triggerType {
				t.Errorf("expected %s, got %s", tt.triggerType, lp.TriggerType())
			}
		})
	}
}

// TestLinearProgression_Validate tests LinearProgression validation.
func TestLinearProgression_Validate(t *testing.T) {
	tests := []struct {
		name    string
		lp      LinearProgression
		wantErr bool
		errType error
	}{
		{
			name: "valid AFTER_SESSION progression",
			lp: LinearProgression{
				ID:               "prog-1",
				Name:             "Squat Progression",
				Increment:        5.0,
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: TriggerAfterSession,
			},
			wantErr: false,
		},
		{
			name: "valid AFTER_WEEK progression",
			lp: LinearProgression{
				ID:               "prog-2",
				Name:             "Weekly Progression",
				Increment:        2.5,
				MaxTypeValue:     OneRM,
				TriggerTypeValue: TriggerAfterWeek,
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			lp: LinearProgression{
				ID:               "",
				Name:             "Test",
				Increment:        5.0,
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: TriggerAfterSession,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "missing name",
			lp: LinearProgression{
				ID:               "prog-1",
				Name:             "",
				Increment:        5.0,
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: TriggerAfterSession,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "zero increment",
			lp: LinearProgression{
				ID:               "prog-1",
				Name:             "Test",
				Increment:        0,
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: TriggerAfterSession,
			},
			wantErr: true,
			errType: ErrIncrementNotPositive,
		},
		{
			name: "negative increment",
			lp: LinearProgression{
				ID:               "prog-1",
				Name:             "Test",
				Increment:        -5.0,
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: TriggerAfterSession,
			},
			wantErr: true,
			errType: ErrIncrementNotPositive,
		},
		{
			name: "missing max type",
			lp: LinearProgression{
				ID:               "prog-1",
				Name:             "Test",
				Increment:        5.0,
				MaxTypeValue:     "",
				TriggerTypeValue: TriggerAfterSession,
			},
			wantErr: true,
			errType: ErrUnknownMaxType,
		},
		{
			name: "invalid max type",
			lp: LinearProgression{
				ID:               "prog-1",
				Name:             "Test",
				Increment:        5.0,
				MaxTypeValue:     "INVALID",
				TriggerTypeValue: TriggerAfterSession,
			},
			wantErr: true,
			errType: ErrUnknownMaxType,
		},
		{
			name: "missing trigger type",
			lp: LinearProgression{
				ID:               "prog-1",
				Name:             "Test",
				Increment:        5.0,
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: "",
			},
			wantErr: true,
			errType: ErrUnknownTriggerType,
		},
		{
			name: "invalid trigger type",
			lp: LinearProgression{
				ID:               "prog-1",
				Name:             "Test",
				Increment:        5.0,
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: "INVALID",
			},
			wantErr: true,
			errType: ErrUnknownTriggerType,
		},
		{
			name: "unsupported trigger type AFTER_CYCLE",
			lp: LinearProgression{
				ID:               "prog-1",
				Name:             "Test",
				Increment:        5.0,
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: TriggerAfterCycle,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.lp.Validate()
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

// TestNewLinearProgression tests the factory function.
func TestNewLinearProgression(t *testing.T) {
	t.Run("valid parameters", func(t *testing.T) {
		lp, err := NewLinearProgression("prog-1", "Squat +5lb", 5.0, TrainingMax, TriggerAfterSession)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if lp.ID != "prog-1" {
			t.Errorf("expected ID 'prog-1', got '%s'", lp.ID)
		}
		if lp.Name != "Squat +5lb" {
			t.Errorf("expected Name 'Squat +5lb', got '%s'", lp.Name)
		}
		if lp.Increment != 5.0 {
			t.Errorf("expected Increment 5.0, got %f", lp.Increment)
		}
		if lp.MaxTypeValue != TrainingMax {
			t.Errorf("expected MaxTypeValue TrainingMax, got %s", lp.MaxTypeValue)
		}
		if lp.TriggerTypeValue != TriggerAfterSession {
			t.Errorf("expected TriggerTypeValue TriggerAfterSession, got %s", lp.TriggerTypeValue)
		}
	})

	t.Run("invalid parameters", func(t *testing.T) {
		_, err := NewLinearProgression("", "Test", 5.0, TrainingMax, TriggerAfterSession)
		if err == nil {
			t.Error("expected error for empty ID")
		}
	})
}

// TestLinearProgression_Apply tests the Apply method.
func TestLinearProgression_Apply(t *testing.T) {
	lp := &LinearProgression{
		ID:               "prog-1",
		Name:             "Squat Progression",
		Increment:        5.0,
		MaxTypeValue:     TrainingMax,
		TriggerTypeValue: TriggerAfterSession,
	}

	ctx := context.Background()

	t.Run("successful application", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:           TriggerAfterSession,
				Timestamp:      time.Now(),
				LiftsPerformed: []string{"squat-uuid", "bench-uuid"},
			},
		}

		result, err := lp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true, reason: %s", result.Reason)
		}
		if result.PreviousValue != 300 {
			t.Errorf("expected PreviousValue 300, got %f", result.PreviousValue)
		}
		if result.NewValue != 305 {
			t.Errorf("expected NewValue 305, got %f", result.NewValue)
		}
		if result.Delta != 5 {
			t.Errorf("expected Delta 5, got %f", result.Delta)
		}
		if result.LiftID != "squat-uuid" {
			t.Errorf("expected LiftID 'squat-uuid', got '%s'", result.LiftID)
		}
		if result.MaxType != TrainingMax {
			t.Errorf("expected MaxType TrainingMax, got %s", result.MaxType)
		}
	})

	t.Run("trigger type mismatch", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:      TriggerAfterWeek, // Wrong trigger type
				Timestamp: time.Now(),
			},
		}

		result, err := lp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied to be false for trigger type mismatch")
		}
		if result.NewValue != result.PreviousValue {
			t.Error("expected NewValue to equal PreviousValue when not applied")
		}
		if result.Delta != 0 {
			t.Error("expected Delta to be 0 when not applied")
		}
		if result.Reason == "" {
			t.Error("expected Reason to be set when not applied")
		}
	})

	t.Run("max type mismatch", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      OneRM, // Wrong max type
			CurrentValue: 335,
			TriggerEvent: TriggerEvent{
				Type:           TriggerAfterSession,
				Timestamp:      time.Now(),
				LiftsPerformed: []string{"squat-uuid"},
			},
		}

		result, err := lp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied to be false for max type mismatch")
		}
		if result.Reason == "" {
			t.Error("expected Reason to be set when not applied")
		}
	})

	t.Run("lift not in session", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "deadlift-uuid", // Not in liftsPerformed
			MaxType:      TrainingMax,
			CurrentValue: 400,
			TriggerEvent: TriggerEvent{
				Type:           TriggerAfterSession,
				Timestamp:      time.Now(),
				LiftsPerformed: []string{"squat-uuid", "bench-uuid"}, // No deadlift
			},
		}

		result, err := lp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied to be false when lift not in session")
		}
		if result.Reason == "" {
			t.Error("expected Reason to be set when lift not in session")
		}
	})

	t.Run("empty lifts performed (applies to all)", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:           TriggerAfterSession,
				Timestamp:      time.Now(),
				LiftsPerformed: []string{}, // Empty list - applies to all
			},
		}

		result, err := lp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true when liftsPerformed is empty, reason: %s", result.Reason)
		}
	})

	t.Run("nil lifts performed (applies to all)", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:      TriggerAfterSession,
				Timestamp: time.Now(),
				// LiftsPerformed is nil
			},
		}

		result, err := lp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true when liftsPerformed is nil, reason: %s", result.Reason)
		}
	})

	t.Run("invalid context", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "", // Missing userID
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:      TriggerAfterSession,
				Timestamp: time.Now(),
			},
		}

		_, err := lp.Apply(ctx, params)
		if err == nil {
			t.Error("expected error for invalid context")
		}
	})
}

// TestLinearProgression_Apply_PerWeek tests Apply for per-week progression.
func TestLinearProgression_Apply_PerWeek(t *testing.T) {
	lp := &LinearProgression{
		ID:               "prog-1",
		Name:             "Bill Starr Weekly",
		Increment:        5.0,
		MaxTypeValue:     TrainingMax,
		TriggerTypeValue: TriggerAfterWeek,
	}

	ctx := context.Background()

	t.Run("successful weekly application", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 200,
			TriggerEvent: TriggerEvent{
				Type:       TriggerAfterWeek,
				Timestamp:  time.Now(),
				WeekNumber: intPtr(2),
			},
		}

		result, err := lp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true, reason: %s", result.Reason)
		}
		if result.NewValue != 205 {
			t.Errorf("expected NewValue 205, got %f", result.NewValue)
		}
	})

	t.Run("weekly progression ignores liftsPerformed", func(t *testing.T) {
		// For AFTER_WEEK, we don't check liftsPerformed
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "deadlift-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:           TriggerAfterWeek,
				Timestamp:      time.Now(),
				WeekNumber:     intPtr(2),
				LiftsPerformed: []string{"squat-uuid"}, // Doesn't include deadlift
			},
		}

		result, err := lp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Weekly progression should apply regardless of liftsPerformed
		if !result.Applied {
			t.Errorf("expected Applied to be true for weekly progression, reason: %s", result.Reason)
		}
	})
}

// TestLinearProgression_Apply_StartingStrength tests Starting Strength pattern.
func TestLinearProgression_Apply_StartingStrength(t *testing.T) {
	// Starting Strength: +5lb per session on squat
	lpSquat, _ := NewLinearProgression("ss-squat", "SS Squat", 5.0, TrainingMax, TriggerAfterSession)
	// Starting Strength: +2.5lb per session on bench/press
	lpBench, _ := NewLinearProgression("ss-bench", "SS Bench", 2.5, TrainingMax, TriggerAfterSession)

	ctx := context.Background()

	t.Run("squat increments by 5", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-1",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 135,
			TriggerEvent: TriggerEvent{
				Type:           TriggerAfterSession,
				Timestamp:      time.Now(),
				LiftsPerformed: []string{"squat-uuid", "bench-uuid", "deadlift-uuid"},
			},
		}

		result, err := lpSquat.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied=true, reason: %s", result.Reason)
		}
		if result.Delta != 5.0 {
			t.Errorf("expected Delta 5.0, got %f", result.Delta)
		}
		if result.NewValue != 140 {
			t.Errorf("expected NewValue 140, got %f", result.NewValue)
		}
	})

	t.Run("bench increments by 2.5", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-1",
			LiftID:       "bench-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 100,
			TriggerEvent: TriggerEvent{
				Type:           TriggerAfterSession,
				Timestamp:      time.Now(),
				LiftsPerformed: []string{"squat-uuid", "bench-uuid", "row-uuid"},
			},
		}

		result, err := lpBench.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied=true, reason: %s", result.Reason)
		}
		if result.Delta != 2.5 {
			t.Errorf("expected Delta 2.5, got %f", result.Delta)
		}
		if result.NewValue != 102.5 {
			t.Errorf("expected NewValue 102.5, got %f", result.NewValue)
		}
	})
}

// TestLinearProgression_Apply_BillStarr tests Bill Starr pattern.
func TestLinearProgression_Apply_BillStarr(t *testing.T) {
	// Bill Starr 5x5: +5lb per week on all lifts
	lp, _ := NewLinearProgression("bs-prog", "Bill Starr Weekly", 5.0, TrainingMax, TriggerAfterWeek)

	ctx := context.Background()

	t.Run("applies to all lifts at week end", func(t *testing.T) {
		lifts := []struct {
			liftID   string
			current  float64
			expected float64
		}{
			{"squat-uuid", 200, 205},
			{"bench-uuid", 150, 155},
			{"row-uuid", 135, 140},
		}

		for _, lift := range lifts {
			params := ProgressionContext{
				UserID:       "user-1",
				LiftID:       lift.liftID,
				MaxType:      TrainingMax,
				CurrentValue: lift.current,
				TriggerEvent: TriggerEvent{
					Type:       TriggerAfterWeek,
					Timestamp:  time.Now(),
					WeekNumber: intPtr(1),
				},
			}

			result, err := lp.Apply(ctx, params)
			if err != nil {
				t.Fatalf("unexpected error for %s: %v", lift.liftID, err)
			}
			if !result.Applied {
				t.Errorf("expected Applied=true for %s, reason: %s", lift.liftID, result.Reason)
			}
			if result.NewValue != lift.expected {
				t.Errorf("for %s: expected NewValue %f, got %f", lift.liftID, lift.expected, result.NewValue)
			}
		}
	})
}

// TestLinearProgression_JSON tests JSON serialization roundtrip.
func TestLinearProgression_JSON(t *testing.T) {
	lp := &LinearProgression{
		ID:               "prog-123",
		Name:             "Test Progression",
		Increment:        5.0,
		MaxTypeValue:     TrainingMax,
		TriggerTypeValue: TriggerAfterSession,
	}

	// Marshal
	data, err := json.Marshal(lp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Verify JSON structure
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if parsed["type"] != string(TypeLinear) {
		t.Errorf("expected type %s, got %v", TypeLinear, parsed["type"])
	}
	if parsed["id"] != "prog-123" {
		t.Errorf("expected id 'prog-123', got %v", parsed["id"])
	}
	if parsed["name"] != "Test Progression" {
		t.Errorf("expected name 'Test Progression', got %v", parsed["name"])
	}
	if parsed["increment"] != 5.0 {
		t.Errorf("expected increment 5.0, got %v", parsed["increment"])
	}
	if parsed["maxType"] != string(TrainingMax) {
		t.Errorf("expected maxType %s, got %v", TrainingMax, parsed["maxType"])
	}
	if parsed["triggerType"] != string(TriggerAfterSession) {
		t.Errorf("expected triggerType %s, got %v", TriggerAfterSession, parsed["triggerType"])
	}

	// Unmarshal back
	var restored LinearProgression
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if restored.ID != lp.ID {
		t.Errorf("ID mismatch: expected %s, got %s", lp.ID, restored.ID)
	}
	if restored.Name != lp.Name {
		t.Errorf("Name mismatch: expected %s, got %s", lp.Name, restored.Name)
	}
	if restored.Increment != lp.Increment {
		t.Errorf("Increment mismatch: expected %f, got %f", lp.Increment, restored.Increment)
	}
	if restored.MaxTypeValue != lp.MaxTypeValue {
		t.Errorf("MaxTypeValue mismatch: expected %s, got %s", lp.MaxTypeValue, restored.MaxTypeValue)
	}
	if restored.TriggerTypeValue != lp.TriggerTypeValue {
		t.Errorf("TriggerTypeValue mismatch: expected %s, got %s", lp.TriggerTypeValue, restored.TriggerTypeValue)
	}
}

// TestUnmarshalLinearProgression tests deserialization function.
func TestUnmarshalLinearProgression(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		jsonData := []byte(`{
			"id": "prog-1",
			"name": "Test",
			"increment": 5.0,
			"maxType": "TRAINING_MAX",
			"triggerType": "AFTER_SESSION"
		}`)

		progression, err := UnmarshalLinearProgression(jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		lp, ok := progression.(*LinearProgression)
		if !ok {
			t.Fatalf("expected *LinearProgression, got %T", progression)
		}
		if lp.ID != "prog-1" {
			t.Errorf("expected ID 'prog-1', got '%s'", lp.ID)
		}
		if lp.Type() != TypeLinear {
			t.Errorf("expected type %s, got %s", TypeLinear, lp.Type())
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		jsonData := []byte(`{invalid}`)
		_, err := UnmarshalLinearProgression(jsonData)
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("invalid progression data", func(t *testing.T) {
		jsonData := []byte(`{
			"id": "",
			"name": "Test",
			"increment": 5.0,
			"maxType": "TRAINING_MAX",
			"triggerType": "AFTER_SESSION"
		}`)

		_, err := UnmarshalLinearProgression(jsonData)
		if err == nil {
			t.Error("expected error for invalid progression data")
		}
	})
}

// TestRegisterLinearProgression tests factory registration.
func TestRegisterLinearProgression(t *testing.T) {
	factory := NewProgressionFactory()

	// Verify not registered initially
	if factory.IsRegistered(TypeLinear) {
		t.Error("TypeLinear should not be registered initially")
	}

	// Register
	RegisterLinearProgression(factory)

	// Verify registered
	if !factory.IsRegistered(TypeLinear) {
		t.Error("TypeLinear should be registered after calling RegisterLinearProgression")
	}

	// Create from factory
	jsonData := []byte(`{
		"type": "LINEAR_PROGRESSION",
		"id": "prog-1",
		"name": "Factory Test",
		"increment": 5.0,
		"maxType": "TRAINING_MAX",
		"triggerType": "AFTER_SESSION"
	}`)

	progression, err := factory.CreateFromJSON(jsonData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if progression.Type() != TypeLinear {
		t.Errorf("expected type %s, got %s", TypeLinear, progression.Type())
	}
}

// TestLinearProgression_Interface verifies that LinearProgression implements Progression.
func TestLinearProgression_Interface(t *testing.T) {
	var _ Progression = &LinearProgression{}
}

// TestLinearProgression_AppliedAt tests that AppliedAt is set correctly.
func TestLinearProgression_AppliedAt(t *testing.T) {
	lp := &LinearProgression{
		ID:               "prog-1",
		Name:             "Test",
		Increment:        5.0,
		MaxTypeValue:     TrainingMax,
		TriggerTypeValue: TriggerAfterSession,
	}

	before := time.Now()

	params := ProgressionContext{
		UserID:       "user-1",
		LiftID:       "lift-1",
		MaxType:      TrainingMax,
		CurrentValue: 100,
		TriggerEvent: TriggerEvent{
			Type:      TriggerAfterSession,
			Timestamp: time.Now(),
		},
	}

	result, err := lp.Apply(context.Background(), params)
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

// TestLinearProgression_VariousIncrements tests various increment values.
func TestLinearProgression_VariousIncrements(t *testing.T) {
	increments := []struct {
		increment float64
		current   float64
		expected  float64
	}{
		{2.5, 100, 102.5},
		{5.0, 200, 205},
		{10.0, 300, 310},
		{1.25, 45, 46.25},
		{0.5, 50, 50.5},
	}

	for _, test := range increments {
		t.Run(formatFloat(test.increment)+"lb", func(t *testing.T) {
			lp := &LinearProgression{
				ID:               "prog-1",
				Name:             "Test",
				Increment:        test.increment,
				MaxTypeValue:     TrainingMax,
				TriggerTypeValue: TriggerAfterSession,
			}

			params := ProgressionContext{
				UserID:       "user-1",
				LiftID:       "lift-1",
				MaxType:      TrainingMax,
				CurrentValue: test.current,
				TriggerEvent: TriggerEvent{
					Type:      TriggerAfterSession,
					Timestamp: time.Now(),
				},
			}

			result, err := lp.Apply(context.Background(), params)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !result.Applied {
				t.Error("expected Applied=true")
			}
			if result.NewValue != test.expected {
				t.Errorf("expected NewValue %f, got %f", test.expected, result.NewValue)
			}
			if result.Delta != test.increment {
				t.Errorf("expected Delta %f, got %f", test.increment, result.Delta)
			}
		})
	}
}

// TestLinearProgression_OneRM tests progression with ONE_RM max type.
func TestLinearProgression_OneRM(t *testing.T) {
	lp := &LinearProgression{
		ID:               "prog-1",
		Name:             "1RM Progression",
		Increment:        5.0,
		MaxTypeValue:     OneRM,
		TriggerTypeValue: TriggerAfterSession,
	}

	ctx := context.Background()

	t.Run("applies to ONE_RM", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-1",
			LiftID:       "squat-uuid",
			MaxType:      OneRM,
			CurrentValue: 315,
			TriggerEvent: TriggerEvent{
				Type:           TriggerAfterSession,
				Timestamp:      time.Now(),
				LiftsPerformed: []string{"squat-uuid"},
			},
		}

		result, err := lp.Apply(ctx, params)
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
				Type:           TriggerAfterSession,
				Timestamp:      time.Now(),
				LiftsPerformed: []string{"squat-uuid"},
			},
		}

		result, err := lp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied=false for max type mismatch")
		}
	})
}

// formatFloat formats a float for test names.
func formatFloat(f float64) string {
	if f == float64(int(f)) {
		return json.Number(string(rune(int(f) + '0'))).String()
	}
	return json.Number(string(rune(int(f*10)/10+'0')) + "." + string(rune(int(f*10)%10+'0'))).String()
}

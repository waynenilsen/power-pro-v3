package progression

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

// TestCycleProgression_Type tests that CycleProgression returns correct type.
func TestCycleProgression_Type(t *testing.T) {
	cp := &CycleProgression{
		ID:           "prog-1",
		Name:         "Test Progression",
		Increment:    10.0,
		MaxTypeValue: TrainingMax,
	}
	if cp.Type() != TypeCycle {
		t.Errorf("expected %s, got %s", TypeCycle, cp.Type())
	}
}

// TestCycleProgression_TriggerType tests that CycleProgression always returns AFTER_CYCLE.
func TestCycleProgression_TriggerType(t *testing.T) {
	cp := &CycleProgression{
		ID:           "prog-1",
		Name:         "Test",
		Increment:    10.0,
		MaxTypeValue: TrainingMax,
	}
	if cp.TriggerType() != TriggerAfterCycle {
		t.Errorf("expected %s, got %s", TriggerAfterCycle, cp.TriggerType())
	}
}

// TestCycleProgression_Validate tests CycleProgression validation.
func TestCycleProgression_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cp      CycleProgression
		wantErr bool
		errType error
	}{
		{
			name: "valid progression",
			cp: CycleProgression{
				ID:           "prog-1",
				Name:         "5/3/1 Lower Body",
				Increment:    10.0,
				MaxTypeValue: TrainingMax,
			},
			wantErr: false,
		},
		{
			name: "valid progression with ONE_RM",
			cp: CycleProgression{
				ID:           "prog-2",
				Name:         "1RM Cycle Progression",
				Increment:    5.0,
				MaxTypeValue: OneRM,
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			cp: CycleProgression{
				ID:           "",
				Name:         "Test",
				Increment:    10.0,
				MaxTypeValue: TrainingMax,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "missing name",
			cp: CycleProgression{
				ID:           "prog-1",
				Name:         "",
				Increment:    10.0,
				MaxTypeValue: TrainingMax,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "zero increment",
			cp: CycleProgression{
				ID:           "prog-1",
				Name:         "Test",
				Increment:    0,
				MaxTypeValue: TrainingMax,
			},
			wantErr: true,
			errType: ErrIncrementNotPositive,
		},
		{
			name: "negative increment",
			cp: CycleProgression{
				ID:           "prog-1",
				Name:         "Test",
				Increment:    -10.0,
				MaxTypeValue: TrainingMax,
			},
			wantErr: true,
			errType: ErrIncrementNotPositive,
		},
		{
			name: "missing max type",
			cp: CycleProgression{
				ID:           "prog-1",
				Name:         "Test",
				Increment:    10.0,
				MaxTypeValue: "",
			},
			wantErr: true,
			errType: ErrUnknownMaxType,
		},
		{
			name: "invalid max type",
			cp: CycleProgression{
				ID:           "prog-1",
				Name:         "Test",
				Increment:    10.0,
				MaxTypeValue: "INVALID",
			},
			wantErr: true,
			errType: ErrUnknownMaxType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cp.Validate()
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

// TestNewCycleProgression tests the factory function.
func TestNewCycleProgression(t *testing.T) {
	t.Run("valid parameters", func(t *testing.T) {
		cp, err := NewCycleProgression("prog-1", "5/3/1 Squat +10lb", 10.0, TrainingMax)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cp.ID != "prog-1" {
			t.Errorf("expected ID 'prog-1', got '%s'", cp.ID)
		}
		if cp.Name != "5/3/1 Squat +10lb" {
			t.Errorf("expected Name '5/3/1 Squat +10lb', got '%s'", cp.Name)
		}
		if cp.Increment != 10.0 {
			t.Errorf("expected Increment 10.0, got %f", cp.Increment)
		}
		if cp.MaxTypeValue != TrainingMax {
			t.Errorf("expected MaxTypeValue TrainingMax, got %s", cp.MaxTypeValue)
		}
	})

	t.Run("invalid parameters", func(t *testing.T) {
		_, err := NewCycleProgression("", "Test", 10.0, TrainingMax)
		if err == nil {
			t.Error("expected error for empty ID")
		}
	})
}

// TestCycleProgression_Apply tests the Apply method.
func TestCycleProgression_Apply(t *testing.T) {
	cp := &CycleProgression{
		ID:           "prog-1",
		Name:         "5/3/1 Progression",
		Increment:    10.0,
		MaxTypeValue: TrainingMax,
	}

	ctx := context.Background()

	t.Run("successful application at cycle end", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:           TriggerAfterCycle,
				Timestamp:      time.Now(),
				CycleIteration: intPtr(1),
				WeekNumber:     intPtr(4), // End of 4-week cycle
			},
		}

		result, err := cp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true, reason: %s", result.Reason)
		}
		if result.PreviousValue != 300 {
			t.Errorf("expected PreviousValue 300, got %f", result.PreviousValue)
		}
		if result.NewValue != 310 {
			t.Errorf("expected NewValue 310, got %f", result.NewValue)
		}
		if result.Delta != 10 {
			t.Errorf("expected Delta 10, got %f", result.Delta)
		}
		if result.LiftID != "squat-uuid" {
			t.Errorf("expected LiftID 'squat-uuid', got '%s'", result.LiftID)
		}
		if result.MaxType != TrainingMax {
			t.Errorf("expected MaxType TrainingMax, got %s", result.MaxType)
		}
	})

	t.Run("trigger type mismatch - AFTER_SESSION", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:           TriggerAfterSession,
				Timestamp:      time.Now(),
				LiftsPerformed: []string{"squat-uuid"},
			},
		}

		result, err := cp.Apply(ctx, params)
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

	t.Run("trigger type mismatch - AFTER_WEEK", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:       TriggerAfterWeek,
				Timestamp:  time.Now(),
				WeekNumber: intPtr(2),
			},
		}

		result, err := cp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied to be false for AFTER_WEEK trigger")
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
				Type:           TriggerAfterCycle,
				Timestamp:      time.Now(),
				CycleIteration: intPtr(1),
			},
		}

		result, err := cp.Apply(ctx, params)
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

	t.Run("invalid context", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "", // Missing userID
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:      TriggerAfterCycle,
				Timestamp: time.Now(),
			},
		}

		_, err := cp.Apply(ctx, params)
		if err == nil {
			t.Error("expected error for invalid context")
		}
	})
}

// TestCycleProgression_ApplyWithOverride tests lift-specific increment overrides.
func TestCycleProgression_ApplyWithOverride(t *testing.T) {
	// Default progression with 5lb increment (typical upper body)
	cp := &CycleProgression{
		ID:           "prog-1",
		Name:         "5/3/1 Progression",
		Increment:    5.0,
		MaxTypeValue: TrainingMax,
	}

	ctx := context.Background()

	t.Run("uses default increment when override is nil", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "bench-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 200,
			TriggerEvent: TriggerEvent{
				Type:           TriggerAfterCycle,
				Timestamp:      time.Now(),
				CycleIteration: intPtr(1),
			},
		}

		result, err := cp.ApplyWithOverride(ctx, params, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied=true, reason: %s", result.Reason)
		}
		if result.Delta != 5.0 {
			t.Errorf("expected Delta 5.0 (default), got %f", result.Delta)
		}
		if result.NewValue != 205 {
			t.Errorf("expected NewValue 205, got %f", result.NewValue)
		}
	})

	t.Run("uses override increment for lower body (10lb)", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:           TriggerAfterCycle,
				Timestamp:      time.Now(),
				CycleIteration: intPtr(1),
			},
		}

		override := 10.0
		result, err := cp.ApplyWithOverride(ctx, params, &override)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied=true, reason: %s", result.Reason)
		}
		if result.Delta != 10.0 {
			t.Errorf("expected Delta 10.0 (override), got %f", result.Delta)
		}
		if result.NewValue != 310 {
			t.Errorf("expected NewValue 310, got %f", result.NewValue)
		}
	})

	t.Run("rejects zero override increment", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:           TriggerAfterCycle,
				Timestamp:      time.Now(),
				CycleIteration: intPtr(1),
			},
		}

		override := 0.0
		_, err := cp.ApplyWithOverride(ctx, params, &override)
		if err == nil {
			t.Error("expected error for zero override increment")
		}
		if !errors.Is(err, ErrIncrementNotPositive) {
			t.Errorf("expected ErrIncrementNotPositive, got %v", err)
		}
	})

	t.Run("rejects negative override increment", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:           TriggerAfterCycle,
				Timestamp:      time.Now(),
				CycleIteration: intPtr(1),
			},
		}

		override := -5.0
		_, err := cp.ApplyWithOverride(ctx, params, &override)
		if err == nil {
			t.Error("expected error for negative override increment")
		}
		if !errors.Is(err, ErrIncrementNotPositive) {
			t.Errorf("expected ErrIncrementNotPositive, got %v", err)
		}
	})
}

// TestCycleProgression_Apply_531Pattern tests 5/3/1 pattern: +5lb upper, +10lb lower.
func TestCycleProgression_Apply_531Pattern(t *testing.T) {
	// 5/3/1 has a single progression per max type, with lift-specific overrides
	cpUpper, _ := NewCycleProgression("531-upper", "5/3/1 Upper Body", 5.0, TrainingMax)
	cpLower, _ := NewCycleProgression("531-lower", "5/3/1 Lower Body", 10.0, TrainingMax)

	ctx := context.Background()

	t.Run("bench press +5lb", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-1",
			LiftID:       "bench-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 185,
			TriggerEvent: TriggerEvent{
				Type:           TriggerAfterCycle,
				Timestamp:      time.Now(),
				CycleIteration: intPtr(1),
				WeekNumber:     intPtr(4),
			},
		}

		result, err := cpUpper.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied=true, reason: %s", result.Reason)
		}
		if result.Delta != 5.0 {
			t.Errorf("expected Delta 5.0, got %f", result.Delta)
		}
		if result.NewValue != 190 {
			t.Errorf("expected NewValue 190, got %f", result.NewValue)
		}
	})

	t.Run("overhead press +5lb", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-1",
			LiftID:       "press-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 135,
			TriggerEvent: TriggerEvent{
				Type:           TriggerAfterCycle,
				Timestamp:      time.Now(),
				CycleIteration: intPtr(1),
			},
		}

		result, err := cpUpper.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Delta != 5.0 {
			t.Errorf("expected Delta 5.0, got %f", result.Delta)
		}
		if result.NewValue != 140 {
			t.Errorf("expected NewValue 140, got %f", result.NewValue)
		}
	})

	t.Run("squat +10lb", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-1",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 315,
			TriggerEvent: TriggerEvent{
				Type:           TriggerAfterCycle,
				Timestamp:      time.Now(),
				CycleIteration: intPtr(1),
			},
		}

		result, err := cpLower.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Delta != 10.0 {
			t.Errorf("expected Delta 10.0, got %f", result.Delta)
		}
		if result.NewValue != 325 {
			t.Errorf("expected NewValue 325, got %f", result.NewValue)
		}
	})

	t.Run("deadlift +10lb", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-1",
			LiftID:       "deadlift-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 405,
			TriggerEvent: TriggerEvent{
				Type:           TriggerAfterCycle,
				Timestamp:      time.Now(),
				CycleIteration: intPtr(1),
			},
		}

		result, err := cpLower.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Delta != 10.0 {
			t.Errorf("expected Delta 10.0, got %f", result.Delta)
		}
		if result.NewValue != 415 {
			t.Errorf("expected NewValue 415, got %f", result.NewValue)
		}
	})
}

// TestCycleProgression_Apply_GregNuckolsHF tests Greg Nuckols High Frequency 3-week cycle pattern.
func TestCycleProgression_Apply_GregNuckolsHF(t *testing.T) {
	// Greg Nuckols HF: +5lb at end of 3-week cycle for all lifts
	cp, _ := NewCycleProgression("gn-hf", "Greg Nuckols HF +5lb", 5.0, TrainingMax)

	ctx := context.Background()

	lifts := []struct {
		liftID   string
		current  float64
		expected float64
	}{
		{"squat-uuid", 275, 280},
		{"bench-uuid", 200, 205},
		{"deadlift-uuid", 365, 370},
	}

	for _, lift := range lifts {
		t.Run(lift.liftID, func(t *testing.T) {
			params := ProgressionContext{
				UserID:       "user-1",
				LiftID:       lift.liftID,
				MaxType:      TrainingMax,
				CurrentValue: lift.current,
				TriggerEvent: TriggerEvent{
					Type:           TriggerAfterCycle,
					Timestamp:      time.Now(),
					CycleIteration: intPtr(1),
					WeekNumber:     intPtr(3), // 3-week cycle
				},
			}

			result, err := cp.Apply(ctx, params)
			if err != nil {
				t.Fatalf("unexpected error for %s: %v", lift.liftID, err)
			}
			if !result.Applied {
				t.Errorf("expected Applied=true for %s, reason: %s", lift.liftID, result.Reason)
			}
			if result.NewValue != lift.expected {
				t.Errorf("for %s: expected NewValue %f, got %f", lift.liftID, lift.expected, result.NewValue)
			}
		})
	}
}

// TestCycleProgression_Apply_MultipleCycles tests progression across multiple cycle iterations.
func TestCycleProgression_Apply_MultipleCycles(t *testing.T) {
	cp, _ := NewCycleProgression("531-squat", "5/3/1 Squat", 10.0, TrainingMax)

	ctx := context.Background()

	// Simulate progressing through 4 cycles
	currentValue := 300.0
	for cycle := 1; cycle <= 4; cycle++ {
		params := ProgressionContext{
			UserID:       "user-1",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: currentValue,
			TriggerEvent: TriggerEvent{
				Type:           TriggerAfterCycle,
				Timestamp:      time.Now(),
				CycleIteration: intPtr(cycle),
			},
		}

		result, err := cp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("cycle %d: unexpected error: %v", cycle, err)
		}
		if !result.Applied {
			t.Errorf("cycle %d: expected Applied=true", cycle)
		}

		expectedNew := currentValue + 10.0
		if result.NewValue != expectedNew {
			t.Errorf("cycle %d: expected NewValue %f, got %f", cycle, expectedNew, result.NewValue)
		}

		// Update for next cycle
		currentValue = result.NewValue
	}

	// After 4 cycles: 300 + (4 * 10) = 340
	if currentValue != 340 {
		t.Errorf("expected final value 340, got %f", currentValue)
	}
}

// TestCycleProgression_JSON tests JSON serialization roundtrip.
func TestCycleProgression_JSON(t *testing.T) {
	cp := &CycleProgression{
		ID:           "prog-123",
		Name:         "Test Cycle Progression",
		Increment:    10.0,
		MaxTypeValue: TrainingMax,
	}

	// Marshal
	data, err := json.Marshal(cp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Verify JSON structure
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if parsed["type"] != string(TypeCycle) {
		t.Errorf("expected type %s, got %v", TypeCycle, parsed["type"])
	}
	if parsed["id"] != "prog-123" {
		t.Errorf("expected id 'prog-123', got %v", parsed["id"])
	}
	if parsed["name"] != "Test Cycle Progression" {
		t.Errorf("expected name 'Test Cycle Progression', got %v", parsed["name"])
	}
	if parsed["increment"] != 10.0 {
		t.Errorf("expected increment 10.0, got %v", parsed["increment"])
	}
	if parsed["maxType"] != string(TrainingMax) {
		t.Errorf("expected maxType %s, got %v", TrainingMax, parsed["maxType"])
	}

	// Unmarshal back
	var restored CycleProgression
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if restored.ID != cp.ID {
		t.Errorf("ID mismatch: expected %s, got %s", cp.ID, restored.ID)
	}
	if restored.Name != cp.Name {
		t.Errorf("Name mismatch: expected %s, got %s", cp.Name, restored.Name)
	}
	if restored.Increment != cp.Increment {
		t.Errorf("Increment mismatch: expected %f, got %f", cp.Increment, restored.Increment)
	}
	if restored.MaxTypeValue != cp.MaxTypeValue {
		t.Errorf("MaxTypeValue mismatch: expected %s, got %s", cp.MaxTypeValue, restored.MaxTypeValue)
	}
}

// TestUnmarshalCycleProgression tests deserialization function.
func TestUnmarshalCycleProgression(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		jsonData := []byte(`{
			"id": "prog-1",
			"name": "Test",
			"increment": 10.0,
			"maxType": "TRAINING_MAX"
		}`)

		progression, err := UnmarshalCycleProgression(jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		cp, ok := progression.(*CycleProgression)
		if !ok {
			t.Fatalf("expected *CycleProgression, got %T", progression)
		}
		if cp.ID != "prog-1" {
			t.Errorf("expected ID 'prog-1', got '%s'", cp.ID)
		}
		if cp.Type() != TypeCycle {
			t.Errorf("expected type %s, got %s", TypeCycle, cp.Type())
		}
		if cp.TriggerType() != TriggerAfterCycle {
			t.Errorf("expected trigger type %s, got %s", TriggerAfterCycle, cp.TriggerType())
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		jsonData := []byte(`{invalid}`)
		_, err := UnmarshalCycleProgression(jsonData)
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("invalid progression data", func(t *testing.T) {
		jsonData := []byte(`{
			"id": "",
			"name": "Test",
			"increment": 10.0,
			"maxType": "TRAINING_MAX"
		}`)

		_, err := UnmarshalCycleProgression(jsonData)
		if err == nil {
			t.Error("expected error for invalid progression data")
		}
	})
}

// TestRegisterCycleProgression tests factory registration.
func TestRegisterCycleProgression(t *testing.T) {
	factory := NewProgressionFactory()

	// Verify not registered initially
	if factory.IsRegistered(TypeCycle) {
		t.Error("TypeCycle should not be registered initially")
	}

	// Register
	RegisterCycleProgression(factory)

	// Verify registered
	if !factory.IsRegistered(TypeCycle) {
		t.Error("TypeCycle should be registered after calling RegisterCycleProgression")
	}

	// Create from factory
	jsonData := []byte(`{
		"type": "CYCLE_PROGRESSION",
		"id": "prog-1",
		"name": "Factory Test",
		"increment": 10.0,
		"maxType": "TRAINING_MAX"
	}`)

	progression, err := factory.CreateFromJSON(jsonData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if progression.Type() != TypeCycle {
		t.Errorf("expected type %s, got %s", TypeCycle, progression.Type())
	}
	if progression.TriggerType() != TriggerAfterCycle {
		t.Errorf("expected trigger type %s, got %s", TriggerAfterCycle, progression.TriggerType())
	}
}

// TestCycleProgression_Interface verifies that CycleProgression implements Progression.
func TestCycleProgression_Interface(t *testing.T) {
	var _ Progression = &CycleProgression{}
}

// TestCycleProgression_AppliedAt tests that AppliedAt is set correctly.
func TestCycleProgression_AppliedAt(t *testing.T) {
	cp := &CycleProgression{
		ID:           "prog-1",
		Name:         "Test",
		Increment:    10.0,
		MaxTypeValue: TrainingMax,
	}

	before := time.Now()

	params := ProgressionContext{
		UserID:       "user-1",
		LiftID:       "lift-1",
		MaxType:      TrainingMax,
		CurrentValue: 100,
		TriggerEvent: TriggerEvent{
			Type:      TriggerAfterCycle,
			Timestamp: time.Now(),
		},
	}

	result, err := cp.Apply(context.Background(), params)
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

// TestCycleProgression_VariousIncrements tests various increment values.
func TestCycleProgression_VariousIncrements(t *testing.T) {
	increments := []struct {
		increment float64
		current   float64
		expected  float64
	}{
		{5.0, 200, 205},
		{10.0, 300, 310},
		{2.5, 100, 102.5},
		{15.0, 400, 415},
		{7.5, 250, 257.5},
	}

	for _, test := range increments {
		t.Run(formatCycleFloat(test.increment)+"lb", func(t *testing.T) {
			cp := &CycleProgression{
				ID:           "prog-1",
				Name:         "Test",
				Increment:    test.increment,
				MaxTypeValue: TrainingMax,
			}

			params := ProgressionContext{
				UserID:       "user-1",
				LiftID:       "lift-1",
				MaxType:      TrainingMax,
				CurrentValue: test.current,
				TriggerEvent: TriggerEvent{
					Type:      TriggerAfterCycle,
					Timestamp: time.Now(),
				},
			}

			result, err := cp.Apply(context.Background(), params)
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

// TestCycleProgression_OneRM tests progression with ONE_RM max type.
func TestCycleProgression_OneRM(t *testing.T) {
	cp := &CycleProgression{
		ID:           "prog-1",
		Name:         "1RM Cycle Progression",
		Increment:    10.0,
		MaxTypeValue: OneRM,
	}

	ctx := context.Background()

	t.Run("applies to ONE_RM", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-1",
			LiftID:       "squat-uuid",
			MaxType:      OneRM,
			CurrentValue: 350,
			TriggerEvent: TriggerEvent{
				Type:           TriggerAfterCycle,
				Timestamp:      time.Now(),
				CycleIteration: intPtr(1),
			},
		}

		result, err := cp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied=true, reason: %s", result.Reason)
		}
		if result.NewValue != 360 {
			t.Errorf("expected NewValue 360, got %f", result.NewValue)
		}
	})

	t.Run("does not apply to TRAINING_MAX", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-1",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax, // Mismatch
			CurrentValue: 285,
			TriggerEvent: TriggerEvent{
				Type:           TriggerAfterCycle,
				Timestamp:      time.Now(),
				CycleIteration: intPtr(1),
			},
		}

		result, err := cp.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied=false for max type mismatch")
		}
	})
}

// TestCycleProgression_WorksWithAnyCycleLength tests that cycle length doesn't matter.
func TestCycleProgression_WorksWithAnyCycleLength(t *testing.T) {
	cp, _ := NewCycleProgression("prog-1", "Any Cycle Length", 5.0, TrainingMax)

	ctx := context.Background()

	cycleLengths := []int{3, 4, 5, 6, 8, 12} // Various cycle lengths

	for _, length := range cycleLengths {
		t.Run(formatCycleWeeks(length), func(t *testing.T) {
			params := ProgressionContext{
				UserID:       "user-1",
				LiftID:       "squat-uuid",
				MaxType:      TrainingMax,
				CurrentValue: 300,
				TriggerEvent: TriggerEvent{
					Type:           TriggerAfterCycle,
					Timestamp:      time.Now(),
					CycleIteration: intPtr(1),
					WeekNumber:     intPtr(length), // Week number at cycle end
				},
			}

			result, err := cp.Apply(ctx, params)
			if err != nil {
				t.Fatalf("unexpected error for %d-week cycle: %v", length, err)
			}
			if !result.Applied {
				t.Errorf("expected Applied=true for %d-week cycle, reason: %s", length, result.Reason)
			}
			if result.Delta != 5.0 {
				t.Errorf("expected Delta 5.0, got %f", result.Delta)
			}
		})
	}
}

// Helper to format float for test names
func formatCycleFloat(f float64) string {
	if f == float64(int(f)) {
		return string(rune(int(f)%10 + '0'))
	}
	whole := int(f)
	frac := int((f - float64(whole)) * 10)
	return string([]byte{byte(whole%10 + '0'), '.', byte(frac + '0')})
}

// Helper to format cycle weeks for test names
func formatCycleWeeks(weeks int) string {
	return string([]byte{byte(weeks/10 + '0'), byte(weeks%10 + '0'), '-', 'w', 'e', 'e', 'k'})
}

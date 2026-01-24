package progression

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

// Helper functions for pointer creation
func strPtr(s string) *string    { return &s }
func floatPtr(f float64) *float64 { return &f }

// TestDeloadOnFailure_Type tests that DeloadOnFailure returns correct type.
func TestDeloadOnFailure_Type(t *testing.T) {
	d := &DeloadOnFailure{
		ID:               "prog-1",
		Name:             "Test Deload",
		FailureThreshold: 3,
		DeloadType:       DeloadTypePercent,
		DeloadPercent:    0.10,
		MaxTypeValue:     TrainingMax,
	}
	if d.Type() != TypeDeloadOnFailure {
		t.Errorf("expected %s, got %s", TypeDeloadOnFailure, d.Type())
	}
}

// TestDeloadOnFailure_TriggerType tests that DeloadOnFailure returns ON_FAILURE trigger.
func TestDeloadOnFailure_TriggerType(t *testing.T) {
	d := &DeloadOnFailure{
		ID:               "prog-1",
		Name:             "Test Deload",
		FailureThreshold: 1,
		DeloadType:       DeloadTypePercent,
		DeloadPercent:    0.10,
		MaxTypeValue:     TrainingMax,
	}
	if d.TriggerType() != TriggerOnFailure {
		t.Errorf("expected %s, got %s", TriggerOnFailure, d.TriggerType())
	}
}

// TestDeloadOnFailure_Validate tests DeloadOnFailure validation.
func TestDeloadOnFailure_Validate(t *testing.T) {
	tests := []struct {
		name    string
		d       DeloadOnFailure
		wantErr bool
		errType error
	}{
		{
			name: "valid percent deload",
			d: DeloadOnFailure{
				ID:               "prog-1",
				Name:             "GZCLP Deload",
				FailureThreshold: 1,
				DeloadType:       DeloadTypePercent,
				DeloadPercent:    0.15,
				MaxTypeValue:     TrainingMax,
			},
			wantErr: false,
		},
		{
			name: "valid fixed deload",
			d: DeloadOnFailure{
				ID:               "prog-2",
				Name:             "Texas Method Deload",
				FailureThreshold: 2,
				DeloadType:       DeloadTypeFixed,
				DeloadAmount:     5.0,
				MaxTypeValue:     TrainingMax,
			},
			wantErr: false,
		},
		{
			name: "valid with reset on deload",
			d: DeloadOnFailure{
				ID:               "prog-3",
				Name:             "Reset Deload",
				FailureThreshold: 3,
				DeloadType:       DeloadTypePercent,
				DeloadPercent:    0.10,
				ResetOnDeload:    true,
				MaxTypeValue:     OneRM,
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			d: DeloadOnFailure{
				ID:               "",
				Name:             "Test",
				FailureThreshold: 1,
				DeloadType:       DeloadTypePercent,
				DeloadPercent:    0.10,
				MaxTypeValue:     TrainingMax,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "missing name",
			d: DeloadOnFailure{
				ID:               "prog-1",
				Name:             "",
				FailureThreshold: 1,
				DeloadType:       DeloadTypePercent,
				DeloadPercent:    0.10,
				MaxTypeValue:     TrainingMax,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "zero failure threshold",
			d: DeloadOnFailure{
				ID:               "prog-1",
				Name:             "Test",
				FailureThreshold: 0,
				DeloadType:       DeloadTypePercent,
				DeloadPercent:    0.10,
				MaxTypeValue:     TrainingMax,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "negative failure threshold",
			d: DeloadOnFailure{
				ID:               "prog-1",
				Name:             "Test",
				FailureThreshold: -1,
				DeloadType:       DeloadTypePercent,
				DeloadPercent:    0.10,
				MaxTypeValue:     TrainingMax,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "invalid deload type",
			d: DeloadOnFailure{
				ID:               "prog-1",
				Name:             "Test",
				FailureThreshold: 1,
				DeloadType:       "invalid",
				DeloadPercent:    0.10,
				MaxTypeValue:     TrainingMax,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "empty deload type",
			d: DeloadOnFailure{
				ID:               "prog-1",
				Name:             "Test",
				FailureThreshold: 1,
				DeloadType:       "",
				DeloadPercent:    0.10,
				MaxTypeValue:     TrainingMax,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "percent deload with zero percent",
			d: DeloadOnFailure{
				ID:               "prog-1",
				Name:             "Test",
				FailureThreshold: 1,
				DeloadType:       DeloadTypePercent,
				DeloadPercent:    0,
				MaxTypeValue:     TrainingMax,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "percent deload with negative percent",
			d: DeloadOnFailure{
				ID:               "prog-1",
				Name:             "Test",
				FailureThreshold: 1,
				DeloadType:       DeloadTypePercent,
				DeloadPercent:    -0.10,
				MaxTypeValue:     TrainingMax,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "percent deload exceeds 100%",
			d: DeloadOnFailure{
				ID:               "prog-1",
				Name:             "Test",
				FailureThreshold: 1,
				DeloadType:       DeloadTypePercent,
				DeloadPercent:    1.5,
				MaxTypeValue:     TrainingMax,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "fixed deload with zero amount",
			d: DeloadOnFailure{
				ID:               "prog-1",
				Name:             "Test",
				FailureThreshold: 1,
				DeloadType:       DeloadTypeFixed,
				DeloadAmount:     0,
				MaxTypeValue:     TrainingMax,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "fixed deload with negative amount",
			d: DeloadOnFailure{
				ID:               "prog-1",
				Name:             "Test",
				FailureThreshold: 1,
				DeloadType:       DeloadTypeFixed,
				DeloadAmount:     -5.0,
				MaxTypeValue:     TrainingMax,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "missing max type",
			d: DeloadOnFailure{
				ID:               "prog-1",
				Name:             "Test",
				FailureThreshold: 1,
				DeloadType:       DeloadTypePercent,
				DeloadPercent:    0.10,
				MaxTypeValue:     "",
			},
			wantErr: true,
			errType: ErrUnknownMaxType,
		},
		{
			name: "invalid max type",
			d: DeloadOnFailure{
				ID:               "prog-1",
				Name:             "Test",
				FailureThreshold: 1,
				DeloadType:       DeloadTypePercent,
				DeloadPercent:    0.10,
				MaxTypeValue:     "INVALID",
			},
			wantErr: true,
			errType: ErrUnknownMaxType,
		},
		{
			name: "100% percent deload is valid",
			d: DeloadOnFailure{
				ID:               "prog-1",
				Name:             "Full Reset",
				FailureThreshold: 1,
				DeloadType:       DeloadTypePercent,
				DeloadPercent:    1.0,
				MaxTypeValue:     TrainingMax,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.d.Validate()
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

// TestNewDeloadOnFailure tests the factory function.
func TestNewDeloadOnFailure(t *testing.T) {
	t.Run("valid percent parameters", func(t *testing.T) {
		d, err := NewDeloadOnFailure("prog-1", "GZCLP Deload", 1, DeloadTypePercent, 0.15, 0, true, TrainingMax)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if d.ID != "prog-1" {
			t.Errorf("expected ID 'prog-1', got '%s'", d.ID)
		}
		if d.Name != "GZCLP Deload" {
			t.Errorf("expected Name 'GZCLP Deload', got '%s'", d.Name)
		}
		if d.FailureThreshold != 1 {
			t.Errorf("expected FailureThreshold 1, got %d", d.FailureThreshold)
		}
		if d.DeloadType != DeloadTypePercent {
			t.Errorf("expected DeloadType 'percent', got '%s'", d.DeloadType)
		}
		if d.DeloadPercent != 0.15 {
			t.Errorf("expected DeloadPercent 0.15, got %f", d.DeloadPercent)
		}
		if !d.ResetOnDeload {
			t.Error("expected ResetOnDeload to be true")
		}
	})

	t.Run("valid fixed parameters", func(t *testing.T) {
		d, err := NewDeloadOnFailure("prog-2", "Texas Method Deload", 2, DeloadTypeFixed, 0, 5.0, false, TrainingMax)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if d.DeloadType != DeloadTypeFixed {
			t.Errorf("expected DeloadType 'fixed', got '%s'", d.DeloadType)
		}
		if d.DeloadAmount != 5.0 {
			t.Errorf("expected DeloadAmount 5.0, got %f", d.DeloadAmount)
		}
	})

	t.Run("invalid parameters", func(t *testing.T) {
		_, err := NewDeloadOnFailure("", "Test", 1, DeloadTypePercent, 0.10, 0, true, TrainingMax)
		if err == nil {
			t.Error("expected error for empty ID")
		}
	})
}

// makeFailureTriggerEvent creates a TriggerEvent for failure tests.
func makeFailureTriggerEvent(consecutiveFailures int) TriggerEvent {
	return TriggerEvent{
		Type:                TriggerOnFailure,
		Timestamp:           time.Now(),
		FailedSetID:         strPtr("set-1"),
		TargetReps:          intPtr(5),
		RepsPerformed:       intPtr(3),
		ConsecutiveFailures: intPtr(consecutiveFailures),
		SetWeight:           floatPtr(200),
		ProgressionID:       strPtr("prog-1"),
	}
}

// TestDeloadOnFailure_Apply tests the Apply method.
func TestDeloadOnFailure_Apply(t *testing.T) {
	ctx := context.Background()

	// Create a percent-based deload progression
	percentDeload := &DeloadOnFailure{
		ID:               "prog-1",
		Name:             "10% Deload",
		FailureThreshold: 3,
		DeloadType:       DeloadTypePercent,
		DeloadPercent:    0.10,
		MaxTypeValue:     TrainingMax,
	}

	// Create a fixed-amount deload progression
	fixedDeload := &DeloadOnFailure{
		ID:               "prog-2",
		Name:             "5lb Deload",
		FailureThreshold: 2,
		DeloadType:       DeloadTypeFixed,
		DeloadAmount:     5.0,
		MaxTypeValue:     TrainingMax,
	}

	t.Run("percent deload when threshold met", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 200,
			TriggerEvent: makeFailureTriggerEvent(3),
		}

		result, err := percentDeload.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true, reason: %s", result.Reason)
		}
		if result.PreviousValue != 200 {
			t.Errorf("expected PreviousValue 200, got %f", result.PreviousValue)
		}
		if result.NewValue != 180 {
			t.Errorf("expected NewValue 180 (200 - 10%%), got %f", result.NewValue)
		}
		if result.Delta != -20 {
			t.Errorf("expected Delta -20, got %f", result.Delta)
		}
	})

	t.Run("fixed deload when threshold met", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 200,
			TriggerEvent: makeFailureTriggerEvent(2),
		}

		result, err := fixedDeload.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true, reason: %s", result.Reason)
		}
		if result.NewValue != 195 {
			t.Errorf("expected NewValue 195 (200 - 5lb), got %f", result.NewValue)
		}
		if result.Delta != -5 {
			t.Errorf("expected Delta -5, got %f", result.Delta)
		}
	})

	t.Run("no deload when threshold not met", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 200,
			TriggerEvent: makeFailureTriggerEvent(2), // Threshold is 3
		}

		result, err := percentDeload.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied to be false when threshold not met")
		}
		if result.NewValue != 200 {
			t.Errorf("expected NewValue unchanged at 200, got %f", result.NewValue)
		}
		if result.Delta != 0 {
			t.Errorf("expected Delta 0, got %f", result.Delta)
		}
		if result.Reason == "" {
			t.Error("expected Reason to be set")
		}
	})

	t.Run("trigger type mismatch", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 200,
			TriggerEvent: TriggerEvent{
				Type:      TriggerAfterSession, // Wrong trigger type
				Timestamp: time.Now(),
			},
		}

		result, err := percentDeload.Apply(ctx, params)
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
			CurrentValue: 225,
			TriggerEvent: makeFailureTriggerEvent(3),
		}

		result, err := percentDeload.Apply(ctx, params)
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

	t.Run("invalid context", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "", // Missing userID
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 200,
			TriggerEvent: makeFailureTriggerEvent(3),
		}

		_, err := percentDeload.Apply(ctx, params)
		if err == nil {
			t.Error("expected error for invalid context")
		}
	})

	t.Run("missing consecutive failures", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 200,
			TriggerEvent: TriggerEvent{
				Type:      TriggerOnFailure,
				Timestamp: time.Now(),
				// ConsecutiveFailures is nil
			},
		}

		result, err := percentDeload.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied to be false when ConsecutiveFailures is nil")
		}
		if result.Reason == "" {
			t.Error("expected Reason to be set")
		}
	})
}

// TestDeloadOnFailure_Apply_EdgeCases tests edge cases for the Apply method.
func TestDeloadOnFailure_Apply_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("deload does not go below zero", func(t *testing.T) {
		// Fixed deload larger than current value
		largeDeload := &DeloadOnFailure{
			ID:               "prog-1",
			Name:             "Large Deload",
			FailureThreshold: 1,
			DeloadType:       DeloadTypeFixed,
			DeloadAmount:     100.0,
			MaxTypeValue:     TrainingMax,
		}

		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 50, // Less than deload amount
			TriggerEvent: makeFailureTriggerEvent(1),
		}

		result, err := largeDeload.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true, reason: %s", result.Reason)
		}
		if result.NewValue != 0 {
			t.Errorf("expected NewValue to be 0 (clamped), got %f", result.NewValue)
		}
		if result.Delta != -50 {
			t.Errorf("expected Delta -50 (full reduction to 0), got %f", result.Delta)
		}
	})

	t.Run("threshold of 1 triggers on first failure", func(t *testing.T) {
		immediateDeload := &DeloadOnFailure{
			ID:               "prog-1",
			Name:             "Immediate Deload",
			FailureThreshold: 1,
			DeloadType:       DeloadTypePercent,
			DeloadPercent:    0.10,
			MaxTypeValue:     TrainingMax,
		}

		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 200,
			TriggerEvent: makeFailureTriggerEvent(1),
		}

		result, err := immediateDeload.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true with threshold=1, reason: %s", result.Reason)
		}
	})

	t.Run("exceeds threshold still triggers", func(t *testing.T) {
		d := &DeloadOnFailure{
			ID:               "prog-1",
			Name:             "Deload",
			FailureThreshold: 2,
			DeloadType:       DeloadTypePercent,
			DeloadPercent:    0.10,
			MaxTypeValue:     TrainingMax,
		}

		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax,
			CurrentValue: 200,
			TriggerEvent: makeFailureTriggerEvent(5), // Far exceeds threshold
		}

		result, err := d.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied to be true when exceeding threshold, reason: %s", result.Reason)
		}
	})
}

// TestDeloadOnFailure_ShouldResetFailureCounter tests the ResetOnDeload helper.
func TestDeloadOnFailure_ShouldResetFailureCounter(t *testing.T) {
	t.Run("returns true when ResetOnDeload is true", func(t *testing.T) {
		d := &DeloadOnFailure{ResetOnDeload: true}
		if !d.ShouldResetFailureCounter() {
			t.Error("expected ShouldResetFailureCounter to return true")
		}
	})

	t.Run("returns false when ResetOnDeload is false", func(t *testing.T) {
		d := &DeloadOnFailure{ResetOnDeload: false}
		if d.ShouldResetFailureCounter() {
			t.Error("expected ShouldResetFailureCounter to return false")
		}
	})
}

// TestDeloadOnFailure_JSON tests JSON serialization roundtrip.
func TestDeloadOnFailure_JSON(t *testing.T) {
	d := &DeloadOnFailure{
		ID:               "prog-123",
		Name:             "Test Deload Progression",
		FailureThreshold: 3,
		DeloadType:       DeloadTypePercent,
		DeloadPercent:    0.15,
		ResetOnDeload:    true,
		MaxTypeValue:     TrainingMax,
	}

	// Marshal
	data, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Verify JSON structure
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if parsed["type"] != string(TypeDeloadOnFailure) {
		t.Errorf("expected type %s, got %v", TypeDeloadOnFailure, parsed["type"])
	}
	if parsed["id"] != "prog-123" {
		t.Errorf("expected id 'prog-123', got %v", parsed["id"])
	}
	if parsed["name"] != "Test Deload Progression" {
		t.Errorf("expected name 'Test Deload Progression', got %v", parsed["name"])
	}
	if parsed["failureThreshold"] != float64(3) {
		t.Errorf("expected failureThreshold 3, got %v", parsed["failureThreshold"])
	}
	if parsed["deloadType"] != DeloadTypePercent {
		t.Errorf("expected deloadType 'percent', got %v", parsed["deloadType"])
	}
	if parsed["deloadPercent"] != 0.15 {
		t.Errorf("expected deloadPercent 0.15, got %v", parsed["deloadPercent"])
	}
	if parsed["resetOnDeload"] != true {
		t.Errorf("expected resetOnDeload true, got %v", parsed["resetOnDeload"])
	}
	if parsed["maxType"] != string(TrainingMax) {
		t.Errorf("expected maxType %s, got %v", TrainingMax, parsed["maxType"])
	}

	// Unmarshal back
	var restored DeloadOnFailure
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if restored.ID != d.ID {
		t.Errorf("ID mismatch: expected %s, got %s", d.ID, restored.ID)
	}
	if restored.Name != d.Name {
		t.Errorf("Name mismatch: expected %s, got %s", d.Name, restored.Name)
	}
	if restored.FailureThreshold != d.FailureThreshold {
		t.Errorf("FailureThreshold mismatch: expected %d, got %d", d.FailureThreshold, restored.FailureThreshold)
	}
	if restored.DeloadType != d.DeloadType {
		t.Errorf("DeloadType mismatch: expected %s, got %s", d.DeloadType, restored.DeloadType)
	}
	if restored.DeloadPercent != d.DeloadPercent {
		t.Errorf("DeloadPercent mismatch: expected %f, got %f", d.DeloadPercent, restored.DeloadPercent)
	}
	if restored.ResetOnDeload != d.ResetOnDeload {
		t.Errorf("ResetOnDeload mismatch: expected %v, got %v", d.ResetOnDeload, restored.ResetOnDeload)
	}
	if restored.MaxTypeValue != d.MaxTypeValue {
		t.Errorf("MaxTypeValue mismatch: expected %s, got %s", d.MaxTypeValue, restored.MaxTypeValue)
	}
}

// TestDeloadOnFailure_JSON_FixedDeload tests JSON with fixed deload type.
func TestDeloadOnFailure_JSON_FixedDeload(t *testing.T) {
	d := &DeloadOnFailure{
		ID:               "prog-456",
		Name:             "Fixed Deload",
		FailureThreshold: 2,
		DeloadType:       DeloadTypeFixed,
		DeloadAmount:     10.0,
		ResetOnDeload:    false,
		MaxTypeValue:     OneRM,
	}

	data, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if parsed["deloadType"] != DeloadTypeFixed {
		t.Errorf("expected deloadType 'fixed', got %v", parsed["deloadType"])
	}
	if parsed["deloadAmount"] != 10.0 {
		t.Errorf("expected deloadAmount 10.0, got %v", parsed["deloadAmount"])
	}
}

// TestUnmarshalDeloadOnFailure tests deserialization function.
func TestUnmarshalDeloadOnFailure(t *testing.T) {
	t.Run("valid JSON percent type", func(t *testing.T) {
		jsonData := []byte(`{
			"id": "prog-1",
			"name": "Test",
			"failureThreshold": 3,
			"deloadType": "percent",
			"deloadPercent": 0.10,
			"resetOnDeload": true,
			"maxType": "TRAINING_MAX"
		}`)

		progression, err := UnmarshalDeloadOnFailure(jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		d, ok := progression.(*DeloadOnFailure)
		if !ok {
			t.Fatalf("expected *DeloadOnFailure, got %T", progression)
		}
		if d.ID != "prog-1" {
			t.Errorf("expected ID 'prog-1', got '%s'", d.ID)
		}
		if d.Type() != TypeDeloadOnFailure {
			t.Errorf("expected type %s, got %s", TypeDeloadOnFailure, d.Type())
		}
		if d.FailureThreshold != 3 {
			t.Errorf("expected FailureThreshold 3, got %d", d.FailureThreshold)
		}
	})

	t.Run("valid JSON fixed type", func(t *testing.T) {
		jsonData := []byte(`{
			"id": "prog-2",
			"name": "Fixed Deload",
			"failureThreshold": 2,
			"deloadType": "fixed",
			"deloadAmount": 5.0,
			"maxType": "ONE_RM"
		}`)

		progression, err := UnmarshalDeloadOnFailure(jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		d := progression.(*DeloadOnFailure)
		if d.DeloadType != DeloadTypeFixed {
			t.Errorf("expected DeloadType 'fixed', got '%s'", d.DeloadType)
		}
		if d.DeloadAmount != 5.0 {
			t.Errorf("expected DeloadAmount 5.0, got %f", d.DeloadAmount)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		jsonData := []byte(`{invalid}`)
		_, err := UnmarshalDeloadOnFailure(jsonData)
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("invalid progression data", func(t *testing.T) {
		jsonData := []byte(`{
			"id": "",
			"name": "Test",
			"failureThreshold": 3,
			"deloadType": "percent",
			"deloadPercent": 0.10,
			"maxType": "TRAINING_MAX"
		}`)

		_, err := UnmarshalDeloadOnFailure(jsonData)
		if err == nil {
			t.Error("expected error for invalid progression data")
		}
	})
}

// TestRegisterDeloadOnFailure tests factory registration.
func TestRegisterDeloadOnFailure(t *testing.T) {
	factory := NewProgressionFactory()

	// Verify not registered initially
	if factory.IsRegistered(TypeDeloadOnFailure) {
		t.Error("TypeDeloadOnFailure should not be registered initially")
	}

	// Register
	RegisterDeloadOnFailure(factory)

	// Verify registered
	if !factory.IsRegistered(TypeDeloadOnFailure) {
		t.Error("TypeDeloadOnFailure should be registered after calling RegisterDeloadOnFailure")
	}

	// Create from factory
	jsonData := []byte(`{
		"type": "DELOAD_ON_FAILURE",
		"id": "prog-1",
		"name": "Factory Test",
		"failureThreshold": 2,
		"deloadType": "percent",
		"deloadPercent": 0.15,
		"resetOnDeload": true,
		"maxType": "TRAINING_MAX"
	}`)

	progression, err := factory.CreateFromJSON(jsonData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if progression.Type() != TypeDeloadOnFailure {
		t.Errorf("expected type %s, got %s", TypeDeloadOnFailure, progression.Type())
	}
}

// TestDeloadOnFailure_Interface verifies that DeloadOnFailure implements Progression.
func TestDeloadOnFailure_Interface(t *testing.T) {
	var _ Progression = &DeloadOnFailure{}
}

// TestDeloadOnFailure_AppliedAt tests that AppliedAt is set correctly.
func TestDeloadOnFailure_AppliedAt(t *testing.T) {
	d := &DeloadOnFailure{
		ID:               "prog-1",
		Name:             "Test",
		FailureThreshold: 1,
		DeloadType:       DeloadTypePercent,
		DeloadPercent:    0.10,
		MaxTypeValue:     TrainingMax,
	}

	before := time.Now()

	params := ProgressionContext{
		UserID:       "user-1",
		LiftID:       "lift-1",
		MaxType:      TrainingMax,
		CurrentValue: 100,
		TriggerEvent: makeFailureTriggerEvent(1),
	}

	result, err := d.Apply(context.Background(), params)
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

// TestDeloadOnFailure_GZCLPExample tests a realistic GZCLP T2 style configuration.
func TestDeloadOnFailure_GZCLPExample(t *testing.T) {
	// GZCLP T2: 1 failure -> 15% deload (then move to next rep stage)
	d := &DeloadOnFailure{
		ID:               "gzclp-t2-deload",
		Name:             "GZCLP T2 Deload",
		FailureThreshold: 1,
		DeloadType:       DeloadTypePercent,
		DeloadPercent:    0.15,
		ResetOnDeload:    true,
		MaxTypeValue:     TrainingMax,
	}

	ctx := context.Background()

	testCases := []struct {
		currentValue float64
		failures     int
		shouldApply  bool
		expectedNew  float64
		description  string
	}{
		{100, 1, true, 85, "first failure on 100lb"},
		{200, 1, true, 170, "first failure on 200lb"},
		{150, 2, true, 127.5, "second failure (still triggers)"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			params := ProgressionContext{
				UserID:       "user-1",
				LiftID:       "squat-uuid",
				MaxType:      TrainingMax,
				CurrentValue: tc.currentValue,
				TriggerEvent: makeFailureTriggerEvent(tc.failures),
			}

			result, err := d.Apply(ctx, params)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.shouldApply {
				if !result.Applied {
					t.Errorf("expected Applied=true, reason: %s", result.Reason)
				}
				if result.NewValue != tc.expectedNew {
					t.Errorf("expected NewValue %f, got %f", tc.expectedNew, result.NewValue)
				}
			} else {
				if result.Applied {
					t.Error("expected Applied=false")
				}
			}
		})
	}
}

// TestDeloadOnFailure_TexasMethodExample tests a realistic Texas Method configuration.
func TestDeloadOnFailure_TexasMethodExample(t *testing.T) {
	// Texas Method: 2 consecutive stalls -> reduce 5lb
	d := &DeloadOnFailure{
		ID:               "tm-deload",
		Name:             "Texas Method Deload",
		FailureThreshold: 2,
		DeloadType:       DeloadTypeFixed,
		DeloadAmount:     5.0,
		ResetOnDeload:    true,
		MaxTypeValue:     TrainingMax,
	}

	ctx := context.Background()

	testCases := []struct {
		currentValue float64
		failures     int
		shouldApply  bool
		expectedNew  float64
		description  string
	}{
		{315, 1, false, 315, "first failure (no deload yet)"},
		{315, 2, true, 310, "second consecutive failure triggers 5lb deload"},
		{310, 3, true, 305, "third failure still triggers"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			params := ProgressionContext{
				UserID:       "user-1",
				LiftID:       "squat-uuid",
				MaxType:      TrainingMax,
				CurrentValue: tc.currentValue,
				TriggerEvent: makeFailureTriggerEvent(tc.failures),
			}

			result, err := d.Apply(ctx, params)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.shouldApply {
				if !result.Applied {
					t.Errorf("expected Applied=true, reason: %s", result.Reason)
				}
				if result.NewValue != tc.expectedNew {
					t.Errorf("expected NewValue %f, got %f", tc.expectedNew, result.NewValue)
				}
			} else {
				if result.Applied {
					t.Error("expected Applied=false")
				}
				if result.NewValue != tc.currentValue {
					t.Errorf("expected NewValue unchanged at %f, got %f", tc.currentValue, result.NewValue)
				}
			}
		})
	}
}

// TestDeloadOnFailure_OneRM tests progression with ONE_RM max type.
func TestDeloadOnFailure_OneRM(t *testing.T) {
	d := &DeloadOnFailure{
		ID:               "prog-1",
		Name:             "1RM Deload",
		FailureThreshold: 1,
		DeloadType:       DeloadTypePercent,
		DeloadPercent:    0.10,
		MaxTypeValue:     OneRM,
	}

	ctx := context.Background()

	t.Run("applies to ONE_RM", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-1",
			LiftID:       "squat-uuid",
			MaxType:      OneRM,
			CurrentValue: 400,
			TriggerEvent: makeFailureTriggerEvent(1),
		}

		result, err := d.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Errorf("expected Applied=true, reason: %s", result.Reason)
		}
		if result.NewValue != 360 {
			t.Errorf("expected NewValue 360 (400 - 10%%), got %f", result.NewValue)
		}
	})

	t.Run("does not apply to TRAINING_MAX", func(t *testing.T) {
		params := ProgressionContext{
			UserID:       "user-1",
			LiftID:       "squat-uuid",
			MaxType:      TrainingMax, // Mismatch
			CurrentValue: 340,
			TriggerEvent: makeFailureTriggerEvent(1),
		}

		result, err := d.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied=false for max type mismatch")
		}
	})
}

// TestDeloadOnFailure_VariousPercentages tests various deload percentages.
func TestDeloadOnFailure_VariousPercentages(t *testing.T) {
	ctx := context.Background()

	percentages := []struct {
		percent     float64
		current     float64
		expectedNew float64
		description string
	}{
		{0.05, 200, 190, "5% deload"},
		{0.10, 200, 180, "10% deload"},
		{0.15, 200, 170, "15% deload"},
		{0.20, 200, 160, "20% deload"},
		{0.50, 200, 100, "50% deload"},
		{1.00, 200, 0, "100% deload (full reset)"},
	}

	for _, tc := range percentages {
		t.Run(tc.description, func(t *testing.T) {
			d := &DeloadOnFailure{
				ID:               "prog-1",
				Name:             "Test",
				FailureThreshold: 1,
				DeloadType:       DeloadTypePercent,
				DeloadPercent:    tc.percent,
				MaxTypeValue:     TrainingMax,
			}

			params := ProgressionContext{
				UserID:       "user-1",
				LiftID:       "lift-1",
				MaxType:      TrainingMax,
				CurrentValue: tc.current,
				TriggerEvent: makeFailureTriggerEvent(1),
			}

			result, err := d.Apply(ctx, params)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !result.Applied {
				t.Errorf("expected Applied=true, reason: %s", result.Reason)
			}
			if result.NewValue != tc.expectedNew {
				t.Errorf("expected NewValue %f, got %f", tc.expectedNew, result.NewValue)
			}
		})
	}
}

// TestDeloadOnFailure_VariousFixedAmounts tests various fixed deload amounts.
func TestDeloadOnFailure_VariousFixedAmounts(t *testing.T) {
	ctx := context.Background()

	amounts := []struct {
		amount      float64
		current     float64
		expectedNew float64
		description string
	}{
		{2.5, 100, 97.5, "2.5lb deload"},
		{5.0, 200, 195, "5lb deload"},
		{10.0, 300, 290, "10lb deload"},
		{20.0, 200, 180, "20lb deload"},
	}

	for _, tc := range amounts {
		t.Run(tc.description, func(t *testing.T) {
			d := &DeloadOnFailure{
				ID:               "prog-1",
				Name:             "Test",
				FailureThreshold: 1,
				DeloadType:       DeloadTypeFixed,
				DeloadAmount:     tc.amount,
				MaxTypeValue:     TrainingMax,
			}

			params := ProgressionContext{
				UserID:       "user-1",
				LiftID:       "lift-1",
				MaxType:      TrainingMax,
				CurrentValue: tc.current,
				TriggerEvent: makeFailureTriggerEvent(1),
			}

			result, err := d.Apply(ctx, params)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !result.Applied {
				t.Errorf("expected Applied=true, reason: %s", result.Reason)
			}
			if result.NewValue != tc.expectedNew {
				t.Errorf("expected NewValue %f, got %f", tc.expectedNew, result.NewValue)
			}
		})
	}
}

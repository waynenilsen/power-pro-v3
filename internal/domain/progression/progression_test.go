package progression

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

// mockProgression implements Progression for testing purposes.
type mockProgression struct {
	progressionType ProgressionType
	triggerType     TriggerType
	increment       float64
	shouldApply     bool
	valid           bool
	err             error
}

func (m *mockProgression) Type() ProgressionType {
	return m.progressionType
}

func (m *mockProgression) TriggerType() TriggerType {
	return m.triggerType
}

func (m *mockProgression) Apply(ctx context.Context, params ProgressionContext) (ProgressionResult, error) {
	if m.err != nil {
		return ProgressionResult{}, m.err
	}

	if !m.shouldApply {
		return ProgressionResult{
			Applied:       false,
			PreviousValue: params.CurrentValue,
			NewValue:      params.CurrentValue,
			Delta:         0,
			LiftID:        params.LiftID,
			MaxType:       params.MaxType,
			AppliedAt:     time.Now(),
			Reason:        "conditions not met",
		}, nil
	}

	newValue := params.CurrentValue + m.increment
	return ProgressionResult{
		Applied:       true,
		PreviousValue: params.CurrentValue,
		NewValue:      newValue,
		Delta:         m.increment,
		LiftID:        params.LiftID,
		MaxType:       params.MaxType,
		AppliedAt:     time.Now(),
	}, nil
}

func (m *mockProgression) Validate() error {
	if !m.valid {
		return errors.New("mock validation error")
	}
	return nil
}

func (m *mockProgression) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type        ProgressionType `json:"type"`
		TriggerType TriggerType     `json:"triggerType"`
		Increment   float64         `json:"increment"`
	}{
		Type:        m.progressionType,
		TriggerType: m.triggerType,
		Increment:   m.increment,
	})
}

// TestProgressionTypeConstants verifies all progression type constants are defined.
func TestProgressionTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant ProgressionType
		expected string
	}{
		{"TypeLinear", TypeLinear, "LINEAR_PROGRESSION"},
		{"TypeCycle", TypeCycle, "CYCLE_PROGRESSION"},
		{"TypeAMRAP", TypeAMRAP, "AMRAP_PROGRESSION"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.constant) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}

// TestTriggerTypeConstants verifies all trigger type constants are defined.
func TestTriggerTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant TriggerType
		expected string
	}{
		{"TriggerAfterSession", TriggerAfterSession, "AFTER_SESSION"},
		{"TriggerAfterWeek", TriggerAfterWeek, "AFTER_WEEK"},
		{"TriggerAfterCycle", TriggerAfterCycle, "AFTER_CYCLE"},
		{"TriggerAfterSet", TriggerAfterSet, "AFTER_SET"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.constant) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}

// TestMaxTypeConstants verifies all max type constants are defined.
func TestMaxTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant MaxType
		expected string
	}{
		{"OneRM", OneRM, "ONE_RM"},
		{"TrainingMax", TrainingMax, "TRAINING_MAX"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.constant) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}

// TestValidProgressionTypes verifies all progression types are in the valid set.
func TestValidProgressionTypes(t *testing.T) {
	expectedTypes := []ProgressionType{
		TypeLinear,
		TypeCycle,
		TypeAMRAP,
		TypeDeloadOnFailure,
		TypeStage,
		TypeDouble,
	}

	for _, progressionType := range expectedTypes {
		if !ValidProgressionTypes[progressionType] {
			t.Errorf("expected %s to be in ValidProgressionTypes", progressionType)
		}
	}

	if len(ValidProgressionTypes) != len(expectedTypes) {
		t.Errorf("expected %d types in ValidProgressionTypes, got %d", len(expectedTypes), len(ValidProgressionTypes))
	}
}

// TestValidTriggerTypes verifies all trigger types are in the valid set.
func TestValidTriggerTypes(t *testing.T) {
	expectedTypes := []TriggerType{
		TriggerAfterSession,
		TriggerAfterWeek,
		TriggerAfterCycle,
		TriggerAfterSet,
		TriggerOnFailure,
	}

	for _, triggerType := range expectedTypes {
		if !ValidTriggerTypes[triggerType] {
			t.Errorf("expected %s to be in ValidTriggerTypes", triggerType)
		}
	}

	if len(ValidTriggerTypes) != len(expectedTypes) {
		t.Errorf("expected %d types in ValidTriggerTypes, got %d", len(expectedTypes), len(ValidTriggerTypes))
	}
}

// TestValidMaxTypes verifies all max types are in the valid set.
func TestValidMaxTypes(t *testing.T) {
	expectedTypes := []MaxType{
		OneRM,
		TrainingMax,
	}

	for _, maxType := range expectedTypes {
		if !ValidMaxTypes[maxType] {
			t.Errorf("expected %s to be in ValidMaxTypes", maxType)
		}
	}

	if len(ValidMaxTypes) != len(expectedTypes) {
		t.Errorf("expected %d types in ValidMaxTypes, got %d", len(expectedTypes), len(ValidMaxTypes))
	}
}

// TestTriggerEvent_Validate tests TriggerEvent validation.
func TestTriggerEvent_Validate(t *testing.T) {
	tests := []struct {
		name    string
		event   TriggerEvent
		wantErr bool
		errType error
	}{
		{
			name: "valid AFTER_SESSION trigger",
			event: TriggerEvent{
				Type:      TriggerAfterSession,
				Timestamp: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "valid AFTER_WEEK trigger with context",
			event: TriggerEvent{
				Type:       TriggerAfterWeek,
				Timestamp:  time.Now(),
				WeekNumber: intPtr(2),
			},
			wantErr: false,
		},
		{
			name: "valid AFTER_CYCLE trigger with full context",
			event: TriggerEvent{
				Type:           TriggerAfterCycle,
				Timestamp:      time.Now(),
				WeekNumber:     intPtr(4),
				CycleIteration: intPtr(1),
			},
			wantErr: false,
		},
		{
			name: "missing trigger type",
			event: TriggerEvent{
				Type:      "",
				Timestamp: time.Now(),
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
		{
			name: "invalid trigger type",
			event: TriggerEvent{
				Type:      "INVALID_TYPE",
				Timestamp: time.Now(),
			},
			wantErr: true,
			errType: ErrUnknownTriggerType,
		},
		{
			name: "missing timestamp",
			event: TriggerEvent{
				Type: TriggerAfterSession,
			},
			wantErr: true,
			errType: ErrInvalidParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// TestProgressionContext_Validate tests ProgressionContext validation.
func TestProgressionContext_Validate(t *testing.T) {
	validTrigger := TriggerEvent{
		Type:      TriggerAfterSession,
		Timestamp: time.Now(),
	}

	tests := []struct {
		name    string
		ctx     ProgressionContext
		wantErr bool
		errType error
	}{
		{
			name: "valid context",
			ctx: ProgressionContext{
				UserID:       "user-123",
				LiftID:       "lift-456",
				MaxType:      TrainingMax,
				CurrentValue: 300,
				TriggerEvent: validTrigger,
			},
			wantErr: false,
		},
		{
			name: "valid context with ONE_RM",
			ctx: ProgressionContext{
				UserID:       "user-123",
				LiftID:       "lift-456",
				MaxType:      OneRM,
				CurrentValue: 315,
				TriggerEvent: validTrigger,
			},
			wantErr: false,
		},
		{
			name: "missing user ID",
			ctx: ProgressionContext{
				UserID:       "",
				LiftID:       "lift-456",
				MaxType:      TrainingMax,
				CurrentValue: 300,
				TriggerEvent: validTrigger,
			},
			wantErr: true,
			errType: ErrUserIDRequired,
		},
		{
			name: "missing lift ID",
			ctx: ProgressionContext{
				UserID:       "user-123",
				LiftID:       "",
				MaxType:      TrainingMax,
				CurrentValue: 300,
				TriggerEvent: validTrigger,
			},
			wantErr: true,
			errType: ErrLiftIDRequired,
		},
		{
			name: "missing max type",
			ctx: ProgressionContext{
				UserID:       "user-123",
				LiftID:       "lift-456",
				MaxType:      "",
				CurrentValue: 300,
				TriggerEvent: validTrigger,
			},
			wantErr: true,
			errType: ErrMaxTypeRequired,
		},
		{
			name: "invalid max type",
			ctx: ProgressionContext{
				UserID:       "user-123",
				LiftID:       "lift-456",
				MaxType:      "INVALID",
				CurrentValue: 300,
				TriggerEvent: validTrigger,
			},
			wantErr: true,
			errType: ErrUnknownMaxType,
		},
		{
			name: "zero current value",
			ctx: ProgressionContext{
				UserID:       "user-123",
				LiftID:       "lift-456",
				MaxType:      TrainingMax,
				CurrentValue: 0,
				TriggerEvent: validTrigger,
			},
			wantErr: true,
			errType: ErrCurrentValueNotPositive,
		},
		{
			name: "negative current value",
			ctx: ProgressionContext{
				UserID:       "user-123",
				LiftID:       "lift-456",
				MaxType:      TrainingMax,
				CurrentValue: -100,
				TriggerEvent: validTrigger,
			},
			wantErr: true,
			errType: ErrCurrentValueNotPositive,
		},
		{
			name: "invalid trigger event",
			ctx: ProgressionContext{
				UserID:       "user-123",
				LiftID:       "lift-456",
				MaxType:      TrainingMax,
				CurrentValue: 300,
				TriggerEvent: TriggerEvent{
					Type: "INVALID",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ctx.Validate()
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

// TestProgressionResult tests ProgressionResult struct.
func TestProgressionResult(t *testing.T) {
	now := time.Now()

	t.Run("applied progression", func(t *testing.T) {
		result := ProgressionResult{
			Applied:       true,
			PreviousValue: 300,
			NewValue:      305,
			Delta:         5,
			LiftID:        "lift-123",
			MaxType:       TrainingMax,
			AppliedAt:     now,
		}

		if !result.Applied {
			t.Error("expected Applied to be true")
		}
		if result.Delta != 5 {
			t.Errorf("expected Delta 5, got %f", result.Delta)
		}
		if result.NewValue != result.PreviousValue+result.Delta {
			t.Error("NewValue should equal PreviousValue + Delta")
		}
	})

	t.Run("not applied progression", func(t *testing.T) {
		result := ProgressionResult{
			Applied:       false,
			PreviousValue: 300,
			NewValue:      300,
			Delta:         0,
			LiftID:        "lift-123",
			MaxType:       TrainingMax,
			AppliedAt:     now,
			Reason:        "trigger type mismatch",
		}

		if result.Applied {
			t.Error("expected Applied to be false")
		}
		if result.Reason == "" {
			t.Error("expected Reason to be set for non-applied progression")
		}
	})
}

// TestValidateProgressionType tests progression type validation.
func TestValidateProgressionType(t *testing.T) {
	tests := []struct {
		name            string
		progressionType ProgressionType
		wantErr         bool
	}{
		{"valid LINEAR_PROGRESSION", TypeLinear, false},
		{"valid CYCLE_PROGRESSION", TypeCycle, false},
		{"valid AMRAP_PROGRESSION", TypeAMRAP, false},
		{"empty type", "", true},
		{"unknown type", "UNKNOWN_TYPE", true},
		{"lowercase type", "linear_progression", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProgressionType(tt.progressionType)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if !errors.Is(err, ErrUnknownProgressionType) {
					t.Errorf("expected ErrUnknownProgressionType, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// TestValidateTriggerType tests trigger type validation.
func TestValidateTriggerType(t *testing.T) {
	tests := []struct {
		name        string
		triggerType TriggerType
		wantErr     bool
	}{
		{"valid AFTER_SESSION", TriggerAfterSession, false},
		{"valid AFTER_WEEK", TriggerAfterWeek, false},
		{"valid AFTER_CYCLE", TriggerAfterCycle, false},
		{"valid AFTER_SET", TriggerAfterSet, false},
		{"empty type", "", true},
		{"unknown type", "UNKNOWN_TYPE", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTriggerType(tt.triggerType)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if !errors.Is(err, ErrUnknownTriggerType) {
					t.Errorf("expected ErrUnknownTriggerType, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// TestValidateMaxType tests max type validation.
func TestValidateMaxType(t *testing.T) {
	tests := []struct {
		name    string
		maxType MaxType
		wantErr bool
	}{
		{"valid ONE_RM", OneRM, false},
		{"valid TRAINING_MAX", TrainingMax, false},
		{"empty type", "", true},
		{"unknown type", "UNKNOWN_TYPE", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMaxType(tt.maxType)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if !errors.Is(err, ErrUnknownMaxType) {
					t.Errorf("expected ErrUnknownMaxType, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// TestProgressionFactory_Register tests progression registration.
func TestProgressionFactory_Register(t *testing.T) {
	factory := NewProgressionFactory()

	// Verify factory starts empty
	if factory.IsRegistered(TypeLinear) {
		t.Error("factory should not have TypeLinear registered initially")
	}

	// Register a progression
	factory.Register(TypeLinear, func(data json.RawMessage) (Progression, error) {
		return &mockProgression{progressionType: TypeLinear, valid: true}, nil
	})

	// Verify it's registered
	if !factory.IsRegistered(TypeLinear) {
		t.Error("TypeLinear should be registered")
	}

	// Verify other types are not registered
	if factory.IsRegistered(TypeCycle) {
		t.Error("TypeCycle should not be registered")
	}
}

// TestProgressionFactory_Create tests progression creation.
func TestProgressionFactory_Create(t *testing.T) {
	factory := NewProgressionFactory()
	factory.Register(TypeLinear, func(data json.RawMessage) (Progression, error) {
		return &mockProgression{
			progressionType: TypeLinear,
			triggerType:     TriggerAfterSession,
			increment:       5,
			shouldApply:     true,
			valid:           true,
		}, nil
	})

	t.Run("create registered progression", func(t *testing.T) {
		progression, err := factory.Create(TypeLinear, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if progression.Type() != TypeLinear {
			t.Errorf("expected type %s, got %s", TypeLinear, progression.Type())
		}
		if progression.TriggerType() != TriggerAfterSession {
			t.Errorf("expected trigger type %s, got %s", TriggerAfterSession, progression.TriggerType())
		}
	})

	t.Run("create unregistered progression", func(t *testing.T) {
		_, err := factory.Create(TypeCycle, nil)
		if err == nil {
			t.Error("expected error for unregistered type")
		}
		if !errors.Is(err, ErrProgressionNotRegistered) {
			t.Errorf("expected ErrProgressionNotRegistered, got %v", err)
		}
	})
}

// TestProgressionFactory_CreateFromJSON tests progression creation from JSON.
func TestProgressionFactory_CreateFromJSON(t *testing.T) {
	factory := NewProgressionFactory()
	factory.Register(TypeLinear, func(data json.RawMessage) (Progression, error) {
		var parsed struct {
			Type        ProgressionType `json:"type"`
			TriggerType TriggerType     `json:"triggerType"`
			Increment   float64         `json:"increment"`
		}
		if err := json.Unmarshal(data, &parsed); err != nil {
			return nil, err
		}
		return &mockProgression{
			progressionType: parsed.Type,
			triggerType:     parsed.TriggerType,
			increment:       parsed.Increment,
			shouldApply:     true,
			valid:           true,
		}, nil
	})

	t.Run("valid JSON with type", func(t *testing.T) {
		jsonData := []byte(`{"type": "LINEAR_PROGRESSION", "triggerType": "AFTER_SESSION", "increment": 5}`)
		progression, err := factory.CreateFromJSON(jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if progression.Type() != TypeLinear {
			t.Errorf("expected type %s, got %s", TypeLinear, progression.Type())
		}
	})

	t.Run("unregistered type in JSON", func(t *testing.T) {
		jsonData := []byte(`{"type": "CYCLE_PROGRESSION", "increment": 10}`)
		_, err := factory.CreateFromJSON(jsonData)
		if err == nil {
			t.Error("expected error for unregistered type")
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		jsonData := []byte(`{invalid json}`)
		_, err := factory.CreateFromJSON(jsonData)
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("missing type field", func(t *testing.T) {
		jsonData := []byte(`{"increment": 5}`)
		_, err := factory.CreateFromJSON(jsonData)
		if err == nil {
			t.Error("expected error for missing type")
		}
	})
}

// TestProgressionFactory_RegisteredTypes tests retrieval of registered types.
func TestProgressionFactory_RegisteredTypes(t *testing.T) {
	factory := NewProgressionFactory()

	// Initially empty
	types := factory.RegisteredTypes()
	if len(types) != 0 {
		t.Errorf("expected 0 registered types, got %d", len(types))
	}

	// Register some types
	factory.Register(TypeLinear, func(data json.RawMessage) (Progression, error) {
		return &mockProgression{progressionType: TypeLinear}, nil
	})
	factory.Register(TypeCycle, func(data json.RawMessage) (Progression, error) {
		return &mockProgression{progressionType: TypeCycle}, nil
	})

	types = factory.RegisteredTypes()
	if len(types) != 2 {
		t.Errorf("expected 2 registered types, got %d", len(types))
	}

	// Verify both types are in the list
	hasLinear := false
	hasCycle := false
	for _, tt := range types {
		if tt == TypeLinear {
			hasLinear = true
		}
		if tt == TypeCycle {
			hasCycle = true
		}
	}
	if !hasLinear {
		t.Error("expected TypeLinear in registered types")
	}
	if !hasCycle {
		t.Error("expected TypeCycle in registered types")
	}
}

// TestProgressionEnvelope_UnmarshalJSON tests envelope unmarshaling.
func TestProgressionEnvelope_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name         string
		json         string
		expectedType ProgressionType
		wantErr      bool
	}{
		{
			name:         "valid LINEAR_PROGRESSION",
			json:         `{"type": "LINEAR_PROGRESSION", "increment": 5, "triggerType": "AFTER_SESSION"}`,
			expectedType: TypeLinear,
			wantErr:      false,
		},
		{
			name:         "valid CYCLE_PROGRESSION",
			json:         `{"type": "CYCLE_PROGRESSION", "increment": 10}`,
			expectedType: TypeCycle,
			wantErr:      false,
		},
		{
			name:         "empty type",
			json:         `{"type": "", "increment": 5}`,
			expectedType: "",
			wantErr:      false, // Parsing succeeds, validation catches this
		},
		{
			name:    "invalid JSON",
			json:    `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var envelope ProgressionEnvelope
			err := json.Unmarshal([]byte(tt.json), &envelope)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if envelope.Type != tt.expectedType {
					t.Errorf("expected type %s, got %s", tt.expectedType, envelope.Type)
				}
				if len(envelope.Raw) == 0 {
					t.Error("expected Raw to contain the original JSON")
				}
			}
		})
	}
}

// TestMarshalProgression tests progression serialization.
func TestMarshalProgression(t *testing.T) {
	progression := &mockProgression{
		progressionType: TypeLinear,
		triggerType:     TriggerAfterSession,
		increment:       5,
		valid:           true,
	}

	data, err := MarshalProgression(progression)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the JSON contains the type field
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse marshaled JSON: %v", err)
	}

	typeVal, ok := parsed["type"]
	if !ok {
		t.Error("expected 'type' field in marshaled JSON")
	}
	if typeVal != string(TypeLinear) {
		t.Errorf("expected type %s, got %v", TypeLinear, typeVal)
	}

	triggerVal, ok := parsed["triggerType"]
	if !ok {
		t.Error("expected 'triggerType' field in marshaled JSON")
	}
	if triggerVal != string(TriggerAfterSession) {
		t.Errorf("expected triggerType %s, got %v", TriggerAfterSession, triggerVal)
	}
}

// TestProgression_Interface tests that the interface contract is correct.
func TestProgression_Interface(t *testing.T) {
	// Test that mockProgression implements Progression
	var _ Progression = &mockProgression{}

	progression := &mockProgression{
		progressionType: TypeLinear,
		triggerType:     TriggerAfterSession,
		increment:       5,
		shouldApply:     true,
		valid:           true,
	}

	t.Run("Type returns correct type", func(t *testing.T) {
		if progression.Type() != TypeLinear {
			t.Errorf("expected %s, got %s", TypeLinear, progression.Type())
		}
	})

	t.Run("TriggerType returns correct trigger type", func(t *testing.T) {
		if progression.TriggerType() != TriggerAfterSession {
			t.Errorf("expected %s, got %s", TriggerAfterSession, progression.TriggerType())
		}
	})

	t.Run("Apply returns result with applied=true", func(t *testing.T) {
		ctx := context.Background()
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "lift-456",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:      TriggerAfterSession,
				Timestamp: time.Now(),
			},
		}
		result, err := progression.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Applied {
			t.Error("expected Applied to be true")
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
	})

	t.Run("Apply returns result with applied=false when conditions not met", func(t *testing.T) {
		notApplyProgression := &mockProgression{
			progressionType: TypeLinear,
			triggerType:     TriggerAfterSession,
			increment:       5,
			shouldApply:     false,
			valid:           true,
		}

		ctx := context.Background()
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "lift-456",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:      TriggerAfterWeek, // Different trigger type
				Timestamp: time.Now(),
			},
		}
		result, err := notApplyProgression.Apply(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Applied {
			t.Error("expected Applied to be false")
		}
		if result.Reason == "" {
			t.Error("expected Reason to be set")
		}
	})

	t.Run("Apply returns error when configured", func(t *testing.T) {
		errorProgression := &mockProgression{
			progressionType: TypeLinear,
			err:             errors.New("test error"),
		}
		ctx := context.Background()
		params := ProgressionContext{
			UserID:       "user-123",
			LiftID:       "lift-456",
			MaxType:      TrainingMax,
			CurrentValue: 300,
			TriggerEvent: TriggerEvent{
				Type:      TriggerAfterSession,
				Timestamp: time.Now(),
			},
		}
		_, err := errorProgression.Apply(ctx, params)
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("Validate returns nil for valid progression", func(t *testing.T) {
		if err := progression.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Validate returns error for invalid progression", func(t *testing.T) {
		invalidProgression := &mockProgression{valid: false}
		if err := invalidProgression.Validate(); err == nil {
			t.Error("expected validation error")
		}
	})
}

// TestErrors tests that error variables are defined correctly.
func TestErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		msg  string
	}{
		{"ErrUnknownProgressionType", ErrUnknownProgressionType, "unknown progression type"},
		{"ErrUnknownTriggerType", ErrUnknownTriggerType, "unknown trigger type"},
		{"ErrUnknownMaxType", ErrUnknownMaxType, "unknown max type"},
		{"ErrInvalidParams", ErrInvalidParams, "invalid progression parameters"},
		{"ErrProgressionNotRegistered", ErrProgressionNotRegistered, "progression type not registered in factory"},
		{"ErrUserIDRequired", ErrUserIDRequired, "user ID is required"},
		{"ErrLiftIDRequired", ErrLiftIDRequired, "lift ID is required"},
		{"ErrMaxTypeRequired", ErrMaxTypeRequired, "max type is required"},
		{"ErrCurrentValueNotPositive", ErrCurrentValueNotPositive, "current value must be positive"},
		{"ErrTriggerEventRequired", ErrTriggerEventRequired, "trigger event is required"},
		{"ErrIncrementNotPositive", ErrIncrementNotPositive, "increment must be positive"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.msg {
				t.Errorf("expected message %q, got %q", tt.msg, tt.err.Error())
			}
		})
	}
}

// TestNewProgressionFactory tests factory creation.
func TestNewProgressionFactory(t *testing.T) {
	factory := NewProgressionFactory()
	if factory == nil {
		t.Fatal("expected non-nil factory")
	}
	if factory.creators == nil {
		t.Error("expected creators map to be initialized")
	}
	if len(factory.creators) != 0 {
		t.Error("expected creators map to be empty initially")
	}
}

// TestProgressionFactory_OverwriteRegistration tests that re-registering overwrites.
func TestProgressionFactory_OverwriteRegistration(t *testing.T) {
	factory := NewProgressionFactory()

	// Register first version
	factory.Register(TypeLinear, func(data json.RawMessage) (Progression, error) {
		return &mockProgression{progressionType: TypeLinear, increment: 5}, nil
	})

	// Register second version (overwrite)
	factory.Register(TypeLinear, func(data json.RawMessage) (Progression, error) {
		return &mockProgression{progressionType: TypeLinear, increment: 10, shouldApply: true}, nil
	})

	// Create progression and verify it uses the second version
	progression, err := factory.Create(TypeLinear, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()
	params := ProgressionContext{
		UserID:       "u",
		LiftID:       "l",
		MaxType:      TrainingMax,
		CurrentValue: 300,
		TriggerEvent: TriggerEvent{
			Type:      TriggerAfterSession,
			Timestamp: time.Now(),
		},
	}
	result, _ := progression.Apply(ctx, params)
	if result.Delta != 10 {
		t.Errorf("expected increment 10 from second registration, got %f", result.Delta)
	}
}

// TestNewProgression tests the convenience NewProgression function.
func TestNewProgression(t *testing.T) {
	factory := NewProgressionFactory()
	factory.Register(TypeLinear, func(data json.RawMessage) (Progression, error) {
		return &mockProgression{progressionType: TypeLinear, valid: true}, nil
	})

	t.Run("valid progression type", func(t *testing.T) {
		jsonData := json.RawMessage(`{"increment": 5}`)
		progression, err := NewProgression(factory, "LINEAR_PROGRESSION", jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if progression.Type() != TypeLinear {
			t.Errorf("expected type %s, got %s", TypeLinear, progression.Type())
		}
	})

	t.Run("invalid progression type", func(t *testing.T) {
		jsonData := json.RawMessage(`{"increment": 5}`)
		_, err := NewProgression(factory, "INVALID_TYPE", jsonData)
		if err == nil {
			t.Error("expected error for invalid type")
		}
		if !errors.Is(err, ErrUnknownProgressionType) {
			t.Errorf("expected ErrUnknownProgressionType, got %v", err)
		}
	})

	t.Run("unregistered progression type", func(t *testing.T) {
		jsonData := json.RawMessage(`{"increment": 10}`)
		_, err := NewProgression(factory, "CYCLE_PROGRESSION", jsonData)
		if err == nil {
			t.Error("expected error for unregistered type")
		}
	})
}

// TestTriggerEvent_JSON tests TriggerEvent JSON serialization.
func TestTriggerEvent_JSON(t *testing.T) {
	sessionID := "session-123"
	weekNum := 2
	cycleIter := 1
	daySlug := "day-a"

	event := TriggerEvent{
		Type:           TriggerAfterCycle,
		Timestamp:      time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		SessionID:      &sessionID,
		WeekNumber:     &weekNum,
		CycleIteration: &cycleIter,
		DaySlug:        &daySlug,
		LiftsPerformed: []string{"lift-1", "lift-2"},
	}

	// Serialize
	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Deserialize
	var parsed TriggerEvent
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Verify
	if parsed.Type != event.Type {
		t.Errorf("expected type %s, got %s", event.Type, parsed.Type)
	}
	if parsed.SessionID == nil || *parsed.SessionID != sessionID {
		t.Error("SessionID mismatch")
	}
	if parsed.WeekNumber == nil || *parsed.WeekNumber != weekNum {
		t.Error("WeekNumber mismatch")
	}
	if len(parsed.LiftsPerformed) != 2 {
		t.Errorf("expected 2 lifts, got %d", len(parsed.LiftsPerformed))
	}
}

// TestProgressionContext_JSON tests ProgressionContext JSON serialization.
func TestProgressionContext_JSON(t *testing.T) {
	ctx := ProgressionContext{
		UserID:       "user-123",
		LiftID:       "lift-456",
		MaxType:      TrainingMax,
		CurrentValue: 300,
		TriggerEvent: TriggerEvent{
			Type:      TriggerAfterSession,
			Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
	}

	// Serialize
	data, err := json.Marshal(ctx)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Deserialize
	var parsed ProgressionContext
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Verify
	if parsed.UserID != ctx.UserID {
		t.Errorf("expected UserID %s, got %s", ctx.UserID, parsed.UserID)
	}
	if parsed.LiftID != ctx.LiftID {
		t.Errorf("expected LiftID %s, got %s", ctx.LiftID, parsed.LiftID)
	}
	if parsed.MaxType != ctx.MaxType {
		t.Errorf("expected MaxType %s, got %s", ctx.MaxType, parsed.MaxType)
	}
	if parsed.CurrentValue != ctx.CurrentValue {
		t.Errorf("expected CurrentValue %f, got %f", ctx.CurrentValue, parsed.CurrentValue)
	}
}

// TestProgressionResult_JSON tests ProgressionResult JSON serialization.
func TestProgressionResult_JSON(t *testing.T) {
	result := ProgressionResult{
		Applied:       true,
		PreviousValue: 300,
		NewValue:      305,
		Delta:         5,
		LiftID:        "lift-123",
		MaxType:       TrainingMax,
		AppliedAt:     time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}

	// Serialize
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Deserialize
	var parsed ProgressionResult
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Verify
	if parsed.Applied != result.Applied {
		t.Error("Applied mismatch")
	}
	if parsed.PreviousValue != result.PreviousValue {
		t.Errorf("expected PreviousValue %f, got %f", result.PreviousValue, parsed.PreviousValue)
	}
	if parsed.NewValue != result.NewValue {
		t.Errorf("expected NewValue %f, got %f", result.NewValue, parsed.NewValue)
	}
	if parsed.Delta != result.Delta {
		t.Errorf("expected Delta %f, got %f", result.Delta, parsed.Delta)
	}
}

// Helper function for creating int pointers
func intPtr(i int) *int {
	return &i
}

// TestProgressionEnvelope_MarshalJSON tests ProgressionEnvelope.MarshalJSON.
func TestProgressionEnvelope_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		envelope ProgressionEnvelope
		wantType ProgressionType
	}{
		{
			name: "LINEAR_PROGRESSION type",
			envelope: ProgressionEnvelope{
				Type: TypeLinear,
			},
			wantType: TypeLinear,
		},
		{
			name: "CYCLE_PROGRESSION type",
			envelope: ProgressionEnvelope{
				Type: TypeCycle,
			},
			wantType: TypeCycle,
		},
		{
			name: "empty type",
			envelope: ProgressionEnvelope{
				Type: "",
			},
			wantType: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.envelope.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON failed: %v", err)
			}

			// Verify JSON structure
			var parsed map[string]interface{}
			if err := json.Unmarshal(data, &parsed); err != nil {
				t.Fatalf("failed to parse JSON: %v", err)
			}

			typeVal, ok := parsed["type"]
			if !ok {
				t.Error("expected 'type' field in JSON")
			}
			if typeVal != string(tt.wantType) {
				t.Errorf("expected type %q, got %v", tt.wantType, typeVal)
			}
		})
	}
}

// TestProgressionEnvelope_MarshalJSON_Roundtrip tests marshal/unmarshal roundtrip.
func TestProgressionEnvelope_MarshalJSON_Roundtrip(t *testing.T) {
	original := ProgressionEnvelope{
		Type: TypeLinear,
	}

	// Marshal
	data, err := original.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	// Unmarshal
	var restored ProgressionEnvelope
	if err := restored.UnmarshalJSON(data); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	// Verify
	if restored.Type != original.Type {
		t.Errorf("Type mismatch: expected %s, got %s", original.Type, restored.Type)
	}
}

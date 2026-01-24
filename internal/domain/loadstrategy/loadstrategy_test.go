package loadstrategy

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

// mockStrategy implements LoadStrategy for testing purposes.
type mockStrategy struct {
	strategyType LoadStrategyType
	value        float64
	err          error
	valid        bool
}

func (m *mockStrategy) Type() LoadStrategyType {
	return m.strategyType
}

func (m *mockStrategy) CalculateLoad(ctx context.Context, params LoadCalculationParams) (float64, error) {
	if m.err != nil {
		return 0, m.err
	}
	return m.value, nil
}

func (m *mockStrategy) Validate() error {
	if !m.valid {
		return errors.New("mock validation error")
	}
	return nil
}

func (m *mockStrategy) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type  LoadStrategyType `json:"type"`
		Value float64          `json:"value"`
	}{
		Type:  m.strategyType,
		Value: m.value,
	})
}

// TestLoadStrategyTypeConstants verifies all strategy type constants are defined.
func TestLoadStrategyTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant LoadStrategyType
		expected string
	}{
		{"TypePercentOf", TypePercentOf, "PERCENT_OF"},
		{"TypeRPETarget", TypeRPETarget, "RPE_TARGET"},
		{"TypeFixedWeight", TypeFixedWeight, "FIXED_WEIGHT"},
		{"TypeRelativeTo", TypeRelativeTo, "RELATIVE_TO"},
		{"TypeFindRM", TypeFindRM, "FIND_RM"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.constant) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}

// TestValidStrategyTypes verifies all strategy types are in the valid set.
func TestValidStrategyTypes(t *testing.T) {
	expectedTypes := []LoadStrategyType{
		TypePercentOf,
		TypeRPETarget,
		TypeFixedWeight,
		TypeRelativeTo,
		TypeFindRM,
	}

	for _, strategyType := range expectedTypes {
		if !ValidStrategyTypes[strategyType] {
			t.Errorf("expected %s to be in ValidStrategyTypes", strategyType)
		}
	}

	if len(ValidStrategyTypes) != len(expectedTypes) {
		t.Errorf("expected %d types in ValidStrategyTypes, got %d", len(expectedTypes), len(ValidStrategyTypes))
	}
}

// TestLoadCalculationParams_Validate tests parameter validation.
func TestLoadCalculationParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  LoadCalculationParams
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid params",
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "lift-456",
			},
			wantErr: false,
		},
		{
			name: "valid params with context",
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "lift-456",
				Context: map[string]interface{}{
					"extra": "data",
				},
			},
			wantErr: false,
		},
		{
			name: "missing user ID",
			params: LoadCalculationParams{
				UserID: "",
				LiftID: "lift-456",
			},
			wantErr: true,
			errMsg:  "user ID is required",
		},
		{
			name: "missing lift ID",
			params: LoadCalculationParams{
				UserID: "user-123",
				LiftID: "",
			},
			wantErr: true,
			errMsg:  "lift ID is required",
		},
		{
			name: "both missing",
			params: LoadCalculationParams{
				UserID: "",
				LiftID: "",
			},
			wantErr: true,
			errMsg:  "user ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if !errors.Is(err, ErrInvalidParams) {
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

// TestValidateStrategyType tests strategy type validation.
func TestValidateStrategyType(t *testing.T) {
	tests := []struct {
		name         string
		strategyType LoadStrategyType
		wantErr      bool
	}{
		{"valid PERCENT_OF", TypePercentOf, false},
		{"valid RPE_TARGET", TypeRPETarget, false},
		{"valid FIXED_WEIGHT", TypeFixedWeight, false},
		{"valid RELATIVE_TO", TypeRelativeTo, false},
		{"empty type", "", true},
		{"unknown type", "UNKNOWN_TYPE", true},
		{"lowercase type", "percent_of", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStrategyType(tt.strategyType)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if !errors.Is(err, ErrUnknownStrategyType) {
					t.Errorf("expected ErrUnknownStrategyType, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// TestStrategyFactory_Register tests strategy registration.
func TestStrategyFactory_Register(t *testing.T) {
	factory := NewStrategyFactory()

	// Verify factory starts empty
	if factory.IsRegistered(TypePercentOf) {
		t.Error("factory should not have TypePercentOf registered initially")
	}

	// Register a strategy
	factory.Register(TypePercentOf, func(data json.RawMessage) (LoadStrategy, error) {
		return &mockStrategy{strategyType: TypePercentOf, value: 100, valid: true}, nil
	})

	// Verify it's registered
	if !factory.IsRegistered(TypePercentOf) {
		t.Error("TypePercentOf should be registered")
	}

	// Verify other types are not registered
	if factory.IsRegistered(TypeRPETarget) {
		t.Error("TypeRPETarget should not be registered")
	}
}

// TestStrategyFactory_Create tests strategy creation.
func TestStrategyFactory_Create(t *testing.T) {
	factory := NewStrategyFactory()
	factory.Register(TypePercentOf, func(data json.RawMessage) (LoadStrategy, error) {
		return &mockStrategy{strategyType: TypePercentOf, value: 200, valid: true}, nil
	})

	t.Run("create registered strategy", func(t *testing.T) {
		strategy, err := factory.Create(TypePercentOf, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if strategy.Type() != TypePercentOf {
			t.Errorf("expected type %s, got %s", TypePercentOf, strategy.Type())
		}
	})

	t.Run("create unregistered strategy", func(t *testing.T) {
		_, err := factory.Create(TypeRPETarget, nil)
		if err == nil {
			t.Error("expected error for unregistered type")
		}
		if !errors.Is(err, ErrStrategyNotRegistered) {
			t.Errorf("expected ErrStrategyNotRegistered, got %v", err)
		}
	})
}

// TestStrategyFactory_CreateFromJSON tests strategy creation from JSON.
func TestStrategyFactory_CreateFromJSON(t *testing.T) {
	factory := NewStrategyFactory()
	factory.Register(TypePercentOf, func(data json.RawMessage) (LoadStrategy, error) {
		// Parse the JSON to extract any fields
		var parsed struct {
			Type       LoadStrategyType `json:"type"`
			Percentage float64          `json:"percentage"`
		}
		if err := json.Unmarshal(data, &parsed); err != nil {
			return nil, err
		}
		return &mockStrategy{
			strategyType: parsed.Type,
			value:        parsed.Percentage,
			valid:        true,
		}, nil
	})

	t.Run("valid JSON with type", func(t *testing.T) {
		jsonData := []byte(`{"type": "PERCENT_OF", "percentage": 85}`)
		strategy, err := factory.CreateFromJSON(jsonData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if strategy.Type() != TypePercentOf {
			t.Errorf("expected type %s, got %s", TypePercentOf, strategy.Type())
		}
	})

	t.Run("unregistered type in JSON", func(t *testing.T) {
		jsonData := []byte(`{"type": "RPE_TARGET", "rpe": 8}`)
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
		jsonData := []byte(`{"percentage": 85}`)
		_, err := factory.CreateFromJSON(jsonData)
		if err == nil {
			t.Error("expected error for missing type")
		}
	})
}

// TestStrategyFactory_RegisteredTypes tests retrieval of registered types.
func TestStrategyFactory_RegisteredTypes(t *testing.T) {
	factory := NewStrategyFactory()

	// Initially empty
	types := factory.RegisteredTypes()
	if len(types) != 0 {
		t.Errorf("expected 0 registered types, got %d", len(types))
	}

	// Register some types
	factory.Register(TypePercentOf, func(data json.RawMessage) (LoadStrategy, error) {
		return &mockStrategy{strategyType: TypePercentOf}, nil
	})
	factory.Register(TypeFixedWeight, func(data json.RawMessage) (LoadStrategy, error) {
		return &mockStrategy{strategyType: TypeFixedWeight}, nil
	})

	types = factory.RegisteredTypes()
	if len(types) != 2 {
		t.Errorf("expected 2 registered types, got %d", len(types))
	}

	// Verify both types are in the list
	hasPercentOf := false
	hasFixedWeight := false
	for _, tt := range types {
		if tt == TypePercentOf {
			hasPercentOf = true
		}
		if tt == TypeFixedWeight {
			hasFixedWeight = true
		}
	}
	if !hasPercentOf {
		t.Error("expected TypePercentOf in registered types")
	}
	if !hasFixedWeight {
		t.Error("expected TypeFixedWeight in registered types")
	}
}

// TestStrategyEnvelope_UnmarshalJSON tests envelope unmarshaling.
func TestStrategyEnvelope_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name         string
		json         string
		expectedType LoadStrategyType
		wantErr      bool
	}{
		{
			name:         "valid PERCENT_OF",
			json:         `{"type": "PERCENT_OF", "percentage": 85}`,
			expectedType: TypePercentOf,
			wantErr:      false,
		},
		{
			name:         "valid RPE_TARGET",
			json:         `{"type": "RPE_TARGET", "rpe": 8}`,
			expectedType: TypeRPETarget,
			wantErr:      false,
		},
		{
			name:         "empty type",
			json:         `{"type": "", "percentage": 85}`,
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
			var envelope StrategyEnvelope
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

// TestMarshalStrategy tests strategy serialization.
func TestMarshalStrategy(t *testing.T) {
	strategy := &mockStrategy{
		strategyType: TypePercentOf,
		value:        85,
		valid:        true,
	}

	data, err := MarshalStrategy(strategy)
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
	if typeVal != string(TypePercentOf) {
		t.Errorf("expected type %s, got %v", TypePercentOf, typeVal)
	}
}

// TestLoadStrategy_Interface tests that the interface contract is correct.
func TestLoadStrategy_Interface(t *testing.T) {
	// Test that mockStrategy implements LoadStrategy
	var _ LoadStrategy = &mockStrategy{}

	strategy := &mockStrategy{
		strategyType: TypePercentOf,
		value:        315,
		valid:        true,
	}

	t.Run("Type returns correct type", func(t *testing.T) {
		if strategy.Type() != TypePercentOf {
			t.Errorf("expected %s, got %s", TypePercentOf, strategy.Type())
		}
	})

	t.Run("CalculateLoad returns value", func(t *testing.T) {
		ctx := context.Background()
		params := LoadCalculationParams{
			UserID: "user-123",
			LiftID: "lift-456",
		}
		value, err := strategy.CalculateLoad(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if value != 315 {
			t.Errorf("expected 315, got %f", value)
		}
	})

	t.Run("CalculateLoad returns error when configured", func(t *testing.T) {
		errorStrategy := &mockStrategy{
			strategyType: TypePercentOf,
			err:          ErrMaxNotFound,
		}
		ctx := context.Background()
		params := LoadCalculationParams{
			UserID: "user-123",
			LiftID: "lift-456",
		}
		_, err := errorStrategy.CalculateLoad(ctx, params)
		if !errors.Is(err, ErrMaxNotFound) {
			t.Errorf("expected ErrMaxNotFound, got %v", err)
		}
	})

	t.Run("Validate returns nil for valid strategy", func(t *testing.T) {
		if err := strategy.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Validate returns error for invalid strategy", func(t *testing.T) {
		invalidStrategy := &mockStrategy{valid: false}
		if err := invalidStrategy.Validate(); err == nil {
			t.Error("expected validation error")
		}
	})
}

// TestMaxValue tests the MaxValue struct.
func TestMaxValue(t *testing.T) {
	maxVal := MaxValue{
		Value:         315.0,
		EffectiveDate: "2024-01-15",
	}

	if maxVal.Value != 315.0 {
		t.Errorf("expected value 315.0, got %f", maxVal.Value)
	}
	if maxVal.EffectiveDate != "2024-01-15" {
		t.Errorf("expected date 2024-01-15, got %s", maxVal.EffectiveDate)
	}
}

// TestErrors tests that error variables are defined correctly.
func TestErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		msg  string
	}{
		{"ErrUnknownStrategyType", ErrUnknownStrategyType, "unknown load strategy type"},
		{"ErrInvalidParams", ErrInvalidParams, "invalid load calculation parameters"},
		{"ErrMaxNotFound", ErrMaxNotFound, "max not found for user/lift combination"},
		{"ErrStrategyNotRegistered", ErrStrategyNotRegistered, "strategy type not registered in factory"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.msg {
				t.Errorf("expected message %q, got %q", tt.msg, tt.err.Error())
			}
		})
	}
}

// TestNewStrategyFactory tests factory creation.
func TestNewStrategyFactory(t *testing.T) {
	factory := NewStrategyFactory()
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

// TestStrategyFactory_OverwriteRegistration tests that re-registering overwrites.
func TestStrategyFactory_OverwriteRegistration(t *testing.T) {
	factory := NewStrategyFactory()

	// Register first version
	factory.Register(TypePercentOf, func(data json.RawMessage) (LoadStrategy, error) {
		return &mockStrategy{strategyType: TypePercentOf, value: 100}, nil
	})

	// Register second version (overwrite)
	factory.Register(TypePercentOf, func(data json.RawMessage) (LoadStrategy, error) {
		return &mockStrategy{strategyType: TypePercentOf, value: 200}, nil
	})

	// Create strategy and verify it uses the second version
	strategy, err := factory.Create(TypePercentOf, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()
	params := LoadCalculationParams{UserID: "u", LiftID: "l"}
	value, _ := strategy.CalculateLoad(ctx, params)
	if value != 200 {
		t.Errorf("expected value 200 from second registration, got %f", value)
	}
}

// TestStrategyEnvelope_MarshalJSON tests StrategyEnvelope.MarshalJSON.
func TestStrategyEnvelope_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		envelope StrategyEnvelope
		wantType LoadStrategyType
	}{
		{
			name: "PERCENT_OF type",
			envelope: StrategyEnvelope{
				Type: TypePercentOf,
			},
			wantType: TypePercentOf,
		},
		{
			name: "RPE_TARGET type",
			envelope: StrategyEnvelope{
				Type: TypeRPETarget,
			},
			wantType: TypeRPETarget,
		},
		{
			name: "empty type",
			envelope: StrategyEnvelope{
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

// TestStrategyEnvelope_MarshalJSON_Roundtrip tests marshal/unmarshal roundtrip.
func TestStrategyEnvelope_MarshalJSON_Roundtrip(t *testing.T) {
	original := StrategyEnvelope{
		Type: TypePercentOf,
	}

	// Marshal
	data, err := original.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	// Unmarshal
	var restored StrategyEnvelope
	if err := restored.UnmarshalJSON(data); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	// Verify
	if restored.Type != original.Type {
		t.Errorf("Type mismatch: expected %s, got %s", original.Type, restored.Type)
	}
}

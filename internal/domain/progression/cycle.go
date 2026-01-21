// Package progression provides domain logic for progression strategies.
// This file implements the CycleProgression strategy that adds a fixed increment
// at cycle completion. This supports periodized programs like 5/3/1 and Greg Nuckols HF.
package progression

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// CycleProgression implements the Progression interface for cycle-based weight increases.
// This is used by periodized programs that need longer adaptation windows.
// Examples:
//   - 5/3/1: +5lb upper body, +10lb lower body at end of each 4-week cycle
//   - Greg Nuckols HF: +5lb at end of each 3-week cycle
//
// CycleProgression has an implicit AFTER_CYCLE trigger type - it only fires when a cycle completes.
// The cycle length is determined by the program schedule, not the progression configuration.
type CycleProgression struct {
	// ID is the unique identifier for this progression.
	ID string `json:"id"`
	// Name is the human-readable name for this progression.
	Name string `json:"name"`
	// Increment is the default weight to add on cycle completion.
	// This can be overridden per-lift via ProgramProgression.OverrideIncrement.
	Increment float64 `json:"increment"`
	// MaxTypeValue specifies which max to update (ONE_RM or TRAINING_MAX).
	MaxTypeValue MaxType `json:"maxType"`
}

// NewCycleProgression creates a new CycleProgression with the given parameters.
func NewCycleProgression(id, name string, increment float64, maxType MaxType) (*CycleProgression, error) {
	cp := &CycleProgression{
		ID:           id,
		Name:         name,
		Increment:    increment,
		MaxTypeValue: maxType,
	}
	if err := cp.Validate(); err != nil {
		return nil, err
	}
	return cp, nil
}

// Type returns the discriminator string for this progression.
// Implements Progression interface.
func (c *CycleProgression) Type() ProgressionType {
	return TypeCycle
}

// TriggerType returns the trigger type this progression responds to.
// Implements Progression interface.
// CycleProgression always uses AFTER_CYCLE trigger (implicit).
func (c *CycleProgression) TriggerType() TriggerType {
	return TriggerAfterCycle
}

// Validate validates the progression's configuration parameters.
// Implements Progression interface.
func (c *CycleProgression) Validate() error {
	if c.ID == "" {
		return fmt.Errorf("%w: id is required", ErrInvalidParams)
	}
	if c.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidParams)
	}
	if c.Increment <= 0 {
		return ErrIncrementNotPositive
	}
	if err := ValidateMaxType(c.MaxTypeValue); err != nil {
		return err
	}
	return nil
}

// Apply evaluates and applies the progression given the context.
// Implements Progression interface.
//
// For AFTER_CYCLE triggers:
//   - Verifies the trigger event type matches AFTER_CYCLE
//   - Verifies the max type matches the progression's configured max type
//   - Returns applied=false if trigger type is not AFTER_CYCLE
//   - Returns applied=false if max type does not match
//
// Returns ProgressionResult with applied=true and delta=increment on success.
// The increment can be overridden per-lift via ApplyWithOverride.
func (c *CycleProgression) Apply(ctx context.Context, params ProgressionContext) (ProgressionResult, error) {
	return c.ApplyWithOverride(ctx, params, nil)
}

// ApplyWithOverride applies the progression with an optional increment override.
// This enables lift-specific increments (e.g., 5/3/1: +5lb upper, +10lb lower)
// via the ProgramProgression.OverrideIncrement field.
//
// If overrideIncrement is nil, uses the default CycleProgression.Increment.
// If overrideIncrement is not nil, uses the override value.
func (c *CycleProgression) ApplyWithOverride(ctx context.Context, params ProgressionContext, overrideIncrement *float64) (ProgressionResult, error) {
	// Validate context
	if err := params.Validate(); err != nil {
		return ProgressionResult{}, fmt.Errorf("invalid progression context: %w", err)
	}

	now := time.Now()

	// Check if trigger type matches - CycleProgression only responds to AFTER_CYCLE
	if params.TriggerEvent.Type != TriggerAfterCycle {
		return ProgressionResult{
			Applied:       false,
			PreviousValue: params.CurrentValue,
			NewValue:      params.CurrentValue,
			Delta:         0,
			LiftID:        params.LiftID,
			MaxType:       params.MaxType,
			AppliedAt:     now,
			Reason:        fmt.Sprintf("trigger type mismatch: expected %s, got %s", TriggerAfterCycle, params.TriggerEvent.Type),
		}, nil
	}

	// Check if max type matches
	if params.MaxType != c.MaxTypeValue {
		return ProgressionResult{
			Applied:       false,
			PreviousValue: params.CurrentValue,
			NewValue:      params.CurrentValue,
			Delta:         0,
			LiftID:        params.LiftID,
			MaxType:       params.MaxType,
			AppliedAt:     now,
			Reason:        fmt.Sprintf("max type mismatch: expected %s, got %s", c.MaxTypeValue, params.MaxType),
		}, nil
	}

	// Determine increment to use (override or default)
	increment := c.Increment
	if overrideIncrement != nil {
		if *overrideIncrement <= 0 {
			return ProgressionResult{}, ErrIncrementNotPositive
		}
		increment = *overrideIncrement
	}

	// Apply the progression
	newValue := params.CurrentValue + increment

	return ProgressionResult{
		Applied:       true,
		PreviousValue: params.CurrentValue,
		NewValue:      newValue,
		Delta:         increment,
		LiftID:        params.LiftID,
		MaxType:       params.MaxType,
		AppliedAt:     now,
	}, nil
}

// MarshalJSON implements json.Marshaler for CycleProgression.
// Ensures the type discriminator is always included in serialized output.
func (c *CycleProgression) MarshalJSON() ([]byte, error) {
	type Alias CycleProgression
	return json.Marshal(&struct {
		Type ProgressionType `json:"type"`
		*Alias
	}{
		Type:  TypeCycle,
		Alias: (*Alias)(c),
	})
}

// UnmarshalCycleProgression deserializes a CycleProgression from JSON.
// This is used by the ProgressionFactory for type-safe deserialization.
func UnmarshalCycleProgression(data json.RawMessage) (Progression, error) {
	var cp CycleProgression
	if err := json.Unmarshal(data, &cp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cycle progression: %w", err)
	}
	if err := cp.Validate(); err != nil {
		return nil, fmt.Errorf("invalid cycle progression: %w", err)
	}
	return &cp, nil
}

// RegisterCycleProgression registers the CycleProgression type with a factory.
// This should be called during application initialization.
func RegisterCycleProgression(factory *ProgressionFactory) {
	factory.Register(TypeCycle, UnmarshalCycleProgression)
}

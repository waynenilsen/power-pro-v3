// Package progression provides domain logic for progression strategies.
// This file implements the DoubleProgression strategy that increases reps until
// hitting a ceiling, then adds weight and resets reps.
package progression

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// DoubleProgression implements the Progression interface for double progression.
// Double progression works by first increasing reps until a ceiling is reached,
// then adding weight and resetting to the minimum rep target.
//
// Example: 3x8-12 scheme
//   - Start at 3x8 with 100lb
//   - Add reps each session until reaching 3x12
//   - When all sets hit 12 reps, add 5lb and reset to 3x8 at 105lb
//
// This progression type uses AFTER_SET trigger and checks if reps performed
// have reached or exceeded the MaxReps ceiling provided in the TriggerEvent.
type DoubleProgression struct {
	// ID is the unique identifier for this progression.
	ID string `json:"id"`
	// Name is the human-readable name for this progression.
	Name string `json:"name"`
	// WeightIncrement is the weight to add when rep ceiling is reached.
	WeightIncrement float64 `json:"weightIncrement"`
	// MaxTypeValue specifies which max to update (ONE_RM or TRAINING_MAX).
	MaxTypeValue MaxType `json:"maxType"`
	// TriggerTypeValue should be TriggerAfterSet for double progressions.
	TriggerTypeValue TriggerType `json:"triggerType"`
}

// NewDoubleProgression creates a new DoubleProgression with the given parameters.
func NewDoubleProgression(id, name string, weightIncrement float64, maxType MaxType, triggerType TriggerType) (*DoubleProgression, error) {
	dp := &DoubleProgression{
		ID:               id,
		Name:             name,
		WeightIncrement:  weightIncrement,
		MaxTypeValue:     maxType,
		TriggerTypeValue: triggerType,
	}
	if err := dp.Validate(); err != nil {
		return nil, err
	}
	return dp, nil
}

// Type returns the discriminator string for this progression.
// Implements Progression interface.
func (d *DoubleProgression) Type() ProgressionType {
	return TypeDouble
}

// TriggerType returns the trigger type this progression responds to.
// Implements Progression interface.
func (d *DoubleProgression) TriggerType() TriggerType {
	return d.TriggerTypeValue
}

// Validate validates the progression's configuration parameters.
// Implements Progression interface.
func (d *DoubleProgression) Validate() error {
	if d.ID == "" {
		return fmt.Errorf("%w: id is required", ErrInvalidParams)
	}
	if d.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidParams)
	}
	if d.WeightIncrement <= 0 {
		return fmt.Errorf("%w: weight increment must be positive", ErrInvalidParams)
	}
	if err := ValidateMaxType(d.MaxTypeValue); err != nil {
		return err
	}
	if err := ValidateTriggerType(d.TriggerTypeValue); err != nil {
		return err
	}
	// DoubleProgression requires AFTER_SET trigger
	if d.TriggerTypeValue != TriggerAfterSet {
		return fmt.Errorf("%w: double progression requires AFTER_SET trigger type", ErrInvalidParams)
	}
	return nil
}

// Apply evaluates and applies the progression given the context.
// Implements Progression interface.
//
// DoubleProgression implements the double progression pattern commonly used in
// bodybuilding-style rep range schemes (e.g., 3x8-12). The key characteristics are:
//
//  1. Trigger type must be AFTER_SET - fires after logging a set
//  2. The TriggerEvent must include RepsPerformed and MaxReps
//  3. When reps performed >= MaxReps (the ceiling), weight is added:
//     - Applied=true, Delta=WeightIncrement
//  4. When reps performed < MaxReps, no weight change:
//     - Applied=false with reason explaining not at ceiling yet
//
// Note: The rep reset to minimum is handled by the set scheme, not this progression.
// This progression only handles the weight increase decision.
func (d *DoubleProgression) Apply(ctx context.Context, params ProgressionContext) (ProgressionResult, error) {
	if err := params.Validate(); err != nil {
		return ProgressionResult{}, fmt.Errorf("invalid progression context: %w", err)
	}

	now := time.Now()

	// Trigger type must match
	if params.TriggerEvent.Type != d.TriggerTypeValue {
		return ProgressionResult{
			Applied:       false,
			PreviousValue: params.CurrentValue,
			NewValue:      params.CurrentValue,
			Delta:         0,
			LiftID:        params.LiftID,
			MaxType:       params.MaxType,
			AppliedAt:     now,
			Reason:        fmt.Sprintf("trigger type mismatch: expected %s, got %s", d.TriggerTypeValue, params.TriggerEvent.Type),
		}, nil
	}

	// Max type must match
	if params.MaxType != d.MaxTypeValue {
		return ProgressionResult{
			Applied:       false,
			PreviousValue: params.CurrentValue,
			NewValue:      params.CurrentValue,
			Delta:         0,
			LiftID:        params.LiftID,
			MaxType:       params.MaxType,
			AppliedAt:     now,
			Reason:        fmt.Sprintf("max type mismatch: expected %s, got %s", d.MaxTypeValue, params.MaxType),
		}, nil
	}

	// RepsPerformed must be provided
	if params.TriggerEvent.RepsPerformed == nil {
		return ProgressionResult{
			Applied:       false,
			PreviousValue: params.CurrentValue,
			NewValue:      params.CurrentValue,
			Delta:         0,
			LiftID:        params.LiftID,
			MaxType:       params.MaxType,
			AppliedAt:     now,
			Reason:        "reps performed not provided",
		}, nil
	}

	// MaxReps (the ceiling) must be provided
	if params.TriggerEvent.MaxReps == nil {
		return ProgressionResult{
			Applied:       false,
			PreviousValue: params.CurrentValue,
			NewValue:      params.CurrentValue,
			Delta:         0,
			LiftID:        params.LiftID,
			MaxType:       params.MaxType,
			AppliedAt:     now,
			Reason:        "max reps (ceiling) not provided",
		}, nil
	}

	reps := *params.TriggerEvent.RepsPerformed
	maxReps := *params.TriggerEvent.MaxReps

	// Check if rep ceiling has been reached
	if reps >= maxReps {
		newValue := params.CurrentValue + d.WeightIncrement
		return ProgressionResult{
			Applied:       true,
			PreviousValue: params.CurrentValue,
			NewValue:      newValue,
			Delta:         d.WeightIncrement,
			LiftID:        params.LiftID,
			MaxType:       params.MaxType,
			AppliedAt:     now,
		}, nil
	}

	// Rep ceiling not yet reached
	return ProgressionResult{
		Applied:       false,
		PreviousValue: params.CurrentValue,
		NewValue:      params.CurrentValue,
		Delta:         0,
		LiftID:        params.LiftID,
		MaxType:       params.MaxType,
		AppliedAt:     now,
		Reason:        fmt.Sprintf("rep ceiling not reached: performed %d, need %d", reps, maxReps),
	}, nil
}

// MarshalJSON implements json.Marshaler for DoubleProgression.
// Ensures the type discriminator is always included in serialized output.
func (d *DoubleProgression) MarshalJSON() ([]byte, error) {
	type Alias DoubleProgression
	return json.Marshal(&struct {
		Type ProgressionType `json:"type"`
		*Alias
	}{
		Type:  TypeDouble,
		Alias: (*Alias)(d),
	})
}

// UnmarshalDoubleProgression deserializes a DoubleProgression from JSON.
// This is used by the ProgressionFactory for type-safe deserialization.
func UnmarshalDoubleProgression(data json.RawMessage) (Progression, error) {
	var dp DoubleProgression
	if err := json.Unmarshal(data, &dp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal double progression: %w", err)
	}
	if err := dp.Validate(); err != nil {
		return nil, fmt.Errorf("invalid double progression: %w", err)
	}
	return &dp, nil
}

// RegisterDoubleProgression registers the DoubleProgression type with a factory.
// This should be called during application initialization.
func RegisterDoubleProgression(factory *ProgressionFactory) {
	factory.Register(TypeDouble, UnmarshalDoubleProgression)
}

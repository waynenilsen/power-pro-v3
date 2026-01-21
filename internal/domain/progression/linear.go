// Package progression provides domain logic for progression strategies.
// This file implements the LinearProgression strategy that adds a fixed increment
// at regular intervals (per-session or per-week).
package progression

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// LinearProgression implements the Progression interface for linear weight increases.
// This is the most common progression type for beginner and intermediate programs.
// Examples:
//   - Starting Strength: +5lb per session for squat/deadlift, +2.5lb for press/bench
//   - Bill Starr 5x5: +5lb per week on all main lifts
type LinearProgression struct {
	// ID is the unique identifier for this progression.
	ID string `json:"id"`
	// Name is the human-readable name for this progression.
	Name string `json:"name"`
	// Increment is the weight to add on each application.
	Increment float64 `json:"increment"`
	// MaxTypeValue specifies which max to update (ONE_RM or TRAINING_MAX).
	MaxTypeValue MaxType `json:"maxType"`
	// TriggerTypeValue specifies when to fire (AFTER_SESSION or AFTER_WEEK).
	TriggerTypeValue TriggerType `json:"triggerType"`
}

// NewLinearProgression creates a new LinearProgression with the given parameters.
func NewLinearProgression(id, name string, increment float64, maxType MaxType, triggerType TriggerType) (*LinearProgression, error) {
	lp := &LinearProgression{
		ID:               id,
		Name:             name,
		Increment:        increment,
		MaxTypeValue:     maxType,
		TriggerTypeValue: triggerType,
	}
	if err := lp.Validate(); err != nil {
		return nil, err
	}
	return lp, nil
}

// Type returns the discriminator string for this progression.
// Implements Progression interface.
func (l *LinearProgression) Type() ProgressionType {
	return TypeLinear
}

// TriggerType returns the trigger type this progression responds to.
// Implements Progression interface.
func (l *LinearProgression) TriggerType() TriggerType {
	return l.TriggerTypeValue
}

// Validate validates the progression's configuration parameters.
// Implements Progression interface.
func (l *LinearProgression) Validate() error {
	if l.ID == "" {
		return fmt.Errorf("%w: id is required", ErrInvalidParams)
	}
	if l.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidParams)
	}
	if l.Increment <= 0 {
		return ErrIncrementNotPositive
	}
	if err := ValidateMaxType(l.MaxTypeValue); err != nil {
		return err
	}
	if err := ValidateTriggerType(l.TriggerTypeValue); err != nil {
		return err
	}
	// LinearProgression only supports AFTER_SESSION and AFTER_WEEK triggers
	if l.TriggerTypeValue != TriggerAfterSession && l.TriggerTypeValue != TriggerAfterWeek {
		return fmt.Errorf("%w: linear progression only supports AFTER_SESSION and AFTER_WEEK triggers", ErrInvalidParams)
	}
	return nil
}

// Apply evaluates and applies the progression given the context.
// Implements Progression interface.
//
// For AFTER_SESSION triggers:
//   - Verifies the trigger event type matches AFTER_SESSION
//   - Checks if the lift was performed in the session (if lifts are specified)
//   - Returns applied=false if trigger type mismatches or lift not in session
//
// For AFTER_WEEK triggers:
//   - Verifies the trigger event type matches AFTER_WEEK
//   - Applies to all configured lifts regardless of session contents
//   - Returns applied=false if trigger type mismatches
//
// Returns ProgressionResult with applied=true and delta=increment on success.
func (l *LinearProgression) Apply(ctx context.Context, params ProgressionContext) (ProgressionResult, error) {
	// Validate context
	if err := params.Validate(); err != nil {
		return ProgressionResult{}, fmt.Errorf("invalid progression context: %w", err)
	}

	now := time.Now()

	// Check if trigger type matches
	if params.TriggerEvent.Type != l.TriggerTypeValue {
		return ProgressionResult{
			Applied:       false,
			PreviousValue: params.CurrentValue,
			NewValue:      params.CurrentValue,
			Delta:         0,
			LiftID:        params.LiftID,
			MaxType:       params.MaxType,
			AppliedAt:     now,
			Reason:        fmt.Sprintf("trigger type mismatch: expected %s, got %s", l.TriggerTypeValue, params.TriggerEvent.Type),
		}, nil
	}

	// Check if max type matches
	if params.MaxType != l.MaxTypeValue {
		return ProgressionResult{
			Applied:       false,
			PreviousValue: params.CurrentValue,
			NewValue:      params.CurrentValue,
			Delta:         0,
			LiftID:        params.LiftID,
			MaxType:       params.MaxType,
			AppliedAt:     now,
			Reason:        fmt.Sprintf("max type mismatch: expected %s, got %s", l.MaxTypeValue, params.MaxType),
		}, nil
	}

	// For AFTER_SESSION triggers, verify the lift was performed in the session
	if l.TriggerTypeValue == TriggerAfterSession {
		liftsPerformed := params.TriggerEvent.LiftsPerformed
		if len(liftsPerformed) > 0 {
			liftFound := false
			for _, liftID := range liftsPerformed {
				if liftID == params.LiftID {
					liftFound = true
					break
				}
			}
			if !liftFound {
				return ProgressionResult{
					Applied:       false,
					PreviousValue: params.CurrentValue,
					NewValue:      params.CurrentValue,
					Delta:         0,
					LiftID:        params.LiftID,
					MaxType:       params.MaxType,
					AppliedAt:     now,
					Reason:        fmt.Sprintf("lift %s was not performed in this session", params.LiftID),
				}, nil
			}
		}
	}

	// Apply the progression
	newValue := params.CurrentValue + l.Increment

	return ProgressionResult{
		Applied:       true,
		PreviousValue: params.CurrentValue,
		NewValue:      newValue,
		Delta:         l.Increment,
		LiftID:        params.LiftID,
		MaxType:       params.MaxType,
		AppliedAt:     now,
	}, nil
}

// MarshalJSON implements json.Marshaler for LinearProgression.
// Ensures the type discriminator is always included in serialized output.
func (l *LinearProgression) MarshalJSON() ([]byte, error) {
	type Alias LinearProgression
	return json.Marshal(&struct {
		Type ProgressionType `json:"type"`
		*Alias
	}{
		Type:  TypeLinear,
		Alias: (*Alias)(l),
	})
}

// UnmarshalLinearProgression deserializes a LinearProgression from JSON.
// This is used by the ProgressionFactory for type-safe deserialization.
func UnmarshalLinearProgression(data json.RawMessage) (Progression, error) {
	var lp LinearProgression
	if err := json.Unmarshal(data, &lp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal linear progression: %w", err)
	}
	if err := lp.Validate(); err != nil {
		return nil, fmt.Errorf("invalid linear progression: %w", err)
	}
	return &lp, nil
}

// RegisterLinearProgression registers the LinearProgression type with a factory.
// This should be called during application initialization.
func RegisterLinearProgression(factory *ProgressionFactory) {
	factory.Register(TypeLinear, UnmarshalLinearProgression)
}

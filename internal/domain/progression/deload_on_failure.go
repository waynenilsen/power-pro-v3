// Package progression provides domain logic for progression strategies.
// This file implements the DeloadOnFailure strategy that reduces weight after
// consecutive failures, commonly used in programs like GZCLP and Texas Method.
package progression

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// DeloadOnFailure implements the Progression interface for deload-based weight reduction.
// This progression type reduces weight after a configurable number of consecutive failures.
// It is triggered by ON_FAILURE events and only applies when the failure threshold is met.
//
// Examples:
//   - GZCLP T2: 1 failure -> 15% deload (then move to next rep scheme stage)
//   - Texas Method: 2 consecutive stalls -> 5-10lb deload, reset counter
//
// The deload can be configured as either a percentage of current weight or a fixed amount.
type DeloadOnFailure struct {
	// ID is the unique identifier for this progression.
	ID string `json:"id"`
	// Name is the human-readable name for this progression.
	Name string `json:"name"`
	// FailureThreshold is the number of consecutive failures before deload triggers.
	// Must be >= 1. A value of 1 means deload on first failure.
	FailureThreshold int `json:"failureThreshold"`
	// DeloadType specifies how the deload is calculated: "percent" or "fixed".
	DeloadType string `json:"deloadType"`
	// DeloadPercent is the percentage to reduce (e.g., 0.10 for 10% deload).
	// Only used when DeloadType is "percent". Must be between 0 and 1.
	DeloadPercent float64 `json:"deloadPercent,omitempty"`
	// DeloadAmount is the fixed amount to reduce (e.g., 5.0 for 5lb deload).
	// Only used when DeloadType is "fixed". Must be > 0.
	DeloadAmount float64 `json:"deloadAmount,omitempty"`
	// ResetOnDeload indicates whether to reset the failure counter after deload.
	// When true, the failure counter returns to 0 after deload is applied.
	ResetOnDeload bool `json:"resetOnDeload"`
	// MaxTypeValue specifies which max to update (ONE_RM or TRAINING_MAX).
	MaxTypeValue MaxType `json:"maxType"`
}

// DeloadType constants for configuration validation.
const (
	DeloadTypePercent = "percent"
	DeloadTypeFixed   = "fixed"
)

// NewDeloadOnFailure creates a new DeloadOnFailure with the given parameters.
func NewDeloadOnFailure(id, name string, failureThreshold int, deloadType string, deloadPercent, deloadAmount float64, resetOnDeload bool, maxType MaxType) (*DeloadOnFailure, error) {
	d := &DeloadOnFailure{
		ID:               id,
		Name:             name,
		FailureThreshold: failureThreshold,
		DeloadType:       deloadType,
		DeloadPercent:    deloadPercent,
		DeloadAmount:     deloadAmount,
		ResetOnDeload:    resetOnDeload,
		MaxTypeValue:     maxType,
	}
	if err := d.Validate(); err != nil {
		return nil, err
	}
	return d, nil
}

// Type returns the discriminator string for this progression.
// Implements Progression interface.
func (d *DeloadOnFailure) Type() ProgressionType {
	return TypeDeloadOnFailure
}

// TriggerType returns the trigger type this progression responds to.
// DeloadOnFailure only responds to ON_FAILURE triggers.
// Implements Progression interface.
func (d *DeloadOnFailure) TriggerType() TriggerType {
	return TriggerOnFailure
}

// Validate validates the progression's configuration parameters.
// Implements Progression interface.
func (d *DeloadOnFailure) Validate() error {
	if d.ID == "" {
		return fmt.Errorf("%w: id is required", ErrInvalidParams)
	}
	if d.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidParams)
	}
	if d.FailureThreshold < 1 {
		return fmt.Errorf("%w: failureThreshold must be at least 1", ErrInvalidParams)
	}
	if err := ValidateMaxType(d.MaxTypeValue); err != nil {
		return err
	}

	// Validate deload type and associated parameters
	switch d.DeloadType {
	case DeloadTypePercent:
		if d.DeloadPercent <= 0 || d.DeloadPercent > 1 {
			return fmt.Errorf("%w: deloadPercent must be between 0 (exclusive) and 1 (inclusive)", ErrInvalidParams)
		}
	case DeloadTypeFixed:
		if d.DeloadAmount <= 0 {
			return fmt.Errorf("%w: deloadAmount must be positive for fixed deload type", ErrInvalidParams)
		}
	default:
		return fmt.Errorf("%w: deloadType must be 'percent' or 'fixed', got '%s'", ErrInvalidParams, d.DeloadType)
	}

	return nil
}

// Apply evaluates and applies the progression given the context.
// Implements Progression interface.
//
// DeloadOnFailure implements the failure-based deload pattern used in many intermediate
// and advanced programs. The key characteristics are:
//
//  1. Trigger type must be ON_FAILURE - this only fires when a set fails
//  2. The TriggerEvent must contain a FailureTriggerContext with ConsecutiveFailures
//  3. The deload only applies when ConsecutiveFailures >= FailureThreshold
//  4. When threshold is met, weight is reduced by either a percentage or fixed amount
//  5. The delta returned is negative (weight goes down, not up)
//
// This approach allows programs to implement automatic deload protocols that
// help lifters break through plateaus without manual intervention.
func (d *DeloadOnFailure) Apply(ctx context.Context, params ProgressionContext) (ProgressionResult, error) {
	if err := params.Validate(); err != nil {
		return ProgressionResult{}, fmt.Errorf("invalid progression context: %w", err)
	}

	now := time.Now()

	// Trigger type must be ON_FAILURE
	if params.TriggerEvent.Type != TriggerOnFailure {
		return ProgressionResult{
			Applied:       false,
			PreviousValue: params.CurrentValue,
			NewValue:      params.CurrentValue,
			Delta:         0,
			LiftID:        params.LiftID,
			MaxType:       params.MaxType,
			AppliedAt:     now,
			Reason:        fmt.Sprintf("trigger type mismatch: expected %s, got %s", TriggerOnFailure, params.TriggerEvent.Type),
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

	// Extract failure context from the trigger event
	// The trigger event should have ConsecutiveFailures set for ON_FAILURE triggers
	if params.TriggerEvent.ConsecutiveFailures == nil {
		return ProgressionResult{
			Applied:       false,
			PreviousValue: params.CurrentValue,
			NewValue:      params.CurrentValue,
			Delta:         0,
			LiftID:        params.LiftID,
			MaxType:       params.MaxType,
			AppliedAt:     now,
			Reason:        "consecutiveFailures not provided in trigger event",
		}, nil
	}

	consecutiveFailures := *params.TriggerEvent.ConsecutiveFailures

	// Check if failure threshold is met
	if consecutiveFailures < d.FailureThreshold {
		return ProgressionResult{
			Applied:       false,
			PreviousValue: params.CurrentValue,
			NewValue:      params.CurrentValue,
			Delta:         0,
			LiftID:        params.LiftID,
			MaxType:       params.MaxType,
			AppliedAt:     now,
			Reason:        fmt.Sprintf("failure threshold not met: %d consecutive failures, threshold is %d", consecutiveFailures, d.FailureThreshold),
		}, nil
	}

	// Calculate deload amount based on deload type
	var deloadAmount float64
	switch d.DeloadType {
	case DeloadTypePercent:
		// Percent deload: reduce by percentage of current value
		// e.g., 10% deload on 200lb = 200 * 0.10 = 20lb reduction
		deloadAmount = params.CurrentValue * d.DeloadPercent
	case DeloadTypeFixed:
		// Fixed deload: reduce by fixed amount
		deloadAmount = d.DeloadAmount
	}

	// Calculate new value (delta is negative for deload)
	newValue := params.CurrentValue - deloadAmount
	delta := -deloadAmount

	// Ensure we don't go below zero
	if newValue < 0 {
		newValue = 0
		delta = -params.CurrentValue
	}

	return ProgressionResult{
		Applied:       true,
		PreviousValue: params.CurrentValue,
		NewValue:      newValue,
		Delta:         delta,
		LiftID:        params.LiftID,
		MaxType:       params.MaxType,
		AppliedAt:     now,
	}, nil
}

// ShouldResetFailureCounter returns whether the failure counter should be reset
// after this progression is applied. This is used by the service layer to
// coordinate between the progression and the failure tracking system.
func (d *DeloadOnFailure) ShouldResetFailureCounter() bool {
	return d.ResetOnDeload
}

// MarshalJSON implements json.Marshaler for DeloadOnFailure.
// Ensures the type discriminator is always included in serialized output.
func (d *DeloadOnFailure) MarshalJSON() ([]byte, error) {
	type Alias DeloadOnFailure
	return json.Marshal(&struct {
		Type ProgressionType `json:"type"`
		*Alias
	}{
		Type:  TypeDeloadOnFailure,
		Alias: (*Alias)(d),
	})
}

// UnmarshalDeloadOnFailure deserializes a DeloadOnFailure from JSON.
// This is used by the ProgressionFactory for type-safe deserialization.
func UnmarshalDeloadOnFailure(data json.RawMessage) (Progression, error) {
	var d DeloadOnFailure
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, fmt.Errorf("failed to unmarshal deload on failure progression: %w", err)
	}
	if err := d.Validate(); err != nil {
		return nil, fmt.Errorf("invalid deload on failure progression: %w", err)
	}
	return &d, nil
}

// RegisterDeloadOnFailure registers the DeloadOnFailure type with a factory.
// This should be called during application initialization.
func RegisterDeloadOnFailure(factory *ProgressionFactory) {
	factory.Register(TypeDeloadOnFailure, UnmarshalDeloadOnFailure)
}

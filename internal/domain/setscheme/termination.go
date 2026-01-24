// Package setscheme provides domain logic for set/rep scheme strategies.
package setscheme

import (
	"encoding/json"
	"errors"
	"fmt"
)

// TerminationConditionType identifies the type of termination condition.
type TerminationConditionType string

const (
	// TerminationTypeRPEThreshold stops when RPE >= target threshold.
	TerminationTypeRPEThreshold TerminationConditionType = "RPE_THRESHOLD"
	// TerminationTypeRepFailure stops when reps fall below target.
	TerminationTypeRepFailure TerminationConditionType = "REP_FAILURE"
	// TerminationTypeMaxSets stops after a maximum number of sets (safety limit).
	TerminationTypeMaxSets TerminationConditionType = "MAX_SETS"
)

// ErrInvalidTermination indicates invalid termination condition configuration.
var ErrInvalidTermination = errors.New("invalid termination condition")

// TerminationContext provides information about the current set for termination decisions.
type TerminationContext struct {
	// SetNumber is the current set number (1-indexed).
	SetNumber int
	// LastRPE is the RPE of the last completed set. Nil if not tracked.
	LastRPE *float64
	// LastReps is the number of reps performed in the last set.
	LastReps int
	// TotalReps is the cumulative reps performed so far.
	TotalReps int
	// TotalSets is the number of sets completed so far.
	TotalSets int
	// TargetReps is the target reps for the set (what we wanted).
	TargetReps int
}

// TerminationCondition determines when to stop generating more sets.
type TerminationCondition interface {
	// Type returns the discriminator string for this condition.
	Type() TerminationConditionType
	// ShouldTerminate returns true if we should stop generating sets.
	ShouldTerminate(ctx TerminationContext) bool
	// Validate validates the condition's configuration.
	Validate() error
}

// RPEThreshold stops when RPE reaches or exceeds a target threshold.
// Common use case: "Do sets until you hit RPE 10" or "Stop at RPE 9".
type RPEThreshold struct {
	// Threshold is the RPE at or above which we terminate (required, 1-10).
	Threshold float64 `json:"threshold"`
}

// NewRPEThreshold creates a new RPEThreshold condition.
func NewRPEThreshold(threshold float64) (*RPEThreshold, error) {
	c := &RPEThreshold{Threshold: threshold}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return c, nil
}

// Type returns the discriminator string for RPEThreshold.
func (r *RPEThreshold) Type() TerminationConditionType {
	return TerminationTypeRPEThreshold
}

// ShouldTerminate returns true if RPE >= threshold.
// If no RPE is provided (nil), returns false (continue).
func (r *RPEThreshold) ShouldTerminate(ctx TerminationContext) bool {
	if ctx.LastRPE == nil {
		return false
	}
	return *ctx.LastRPE >= r.Threshold
}

// Validate validates the RPEThreshold configuration.
func (r *RPEThreshold) Validate() error {
	if r.Threshold < 1 || r.Threshold > 10 {
		return fmt.Errorf("%w: threshold must be between 1 and 10, got %v", ErrInvalidTermination, r.Threshold)
	}
	return nil
}

// MarshalJSON implements json.Marshaler for RPEThreshold.
func (r *RPEThreshold) MarshalJSON() ([]byte, error) {
	type Alias RPEThreshold
	return json.Marshal(&struct {
		Type TerminationConditionType `json:"type"`
		*Alias
	}{
		Type:  TerminationTypeRPEThreshold,
		Alias: (*Alias)(r),
	})
}

// RepFailure stops when reps performed fall below target.
// Common use case: "Do sets until you can't hit 5 reps".
type RepFailure struct{}

// NewRepFailure creates a new RepFailure condition.
func NewRepFailure() *RepFailure {
	return &RepFailure{}
}

// Type returns the discriminator string for RepFailure.
func (r *RepFailure) Type() TerminationConditionType {
	return TerminationTypeRepFailure
}

// ShouldTerminate returns true if LastReps < TargetReps.
func (r *RepFailure) ShouldTerminate(ctx TerminationContext) bool {
	return ctx.LastReps < ctx.TargetReps
}

// Validate validates the RepFailure configuration (always valid).
func (r *RepFailure) Validate() error {
	return nil
}

// MarshalJSON implements json.Marshaler for RepFailure.
func (r *RepFailure) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type TerminationConditionType `json:"type"`
	}{
		Type: TerminationTypeRepFailure,
	})
}

// MaxSets stops after a maximum number of sets (safety limit).
// This is typically used as a secondary condition to prevent infinite sets.
type MaxSets struct {
	// Max is the maximum number of sets allowed (required, >= 1).
	Max int `json:"max"`
}

// NewMaxSets creates a new MaxSets condition.
func NewMaxSets(max int) (*MaxSets, error) {
	c := &MaxSets{Max: max}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return c, nil
}

// Type returns the discriminator string for MaxSets.
func (m *MaxSets) Type() TerminationConditionType {
	return TerminationTypeMaxSets
}

// ShouldTerminate returns true if TotalSets >= Max.
func (m *MaxSets) ShouldTerminate(ctx TerminationContext) bool {
	return ctx.TotalSets >= m.Max
}

// Validate validates the MaxSets configuration.
func (m *MaxSets) Validate() error {
	if m.Max < 1 {
		return fmt.Errorf("%w: max must be >= 1, got %d", ErrInvalidTermination, m.Max)
	}
	return nil
}

// MarshalJSON implements json.Marshaler for MaxSets.
func (m *MaxSets) MarshalJSON() ([]byte, error) {
	type Alias MaxSets
	return json.Marshal(&struct {
		Type TerminationConditionType `json:"type"`
		*Alias
	}{
		Type:  TerminationTypeMaxSets,
		Alias: (*Alias)(m),
	})
}

// TerminationConditionEnvelope is the JSON wrapper for polymorphic TerminationCondition serialization.
type TerminationConditionEnvelope struct {
	Type TerminationConditionType `json:"type"`
	Raw  json.RawMessage          `json:"-"`
}

// UnmarshalJSON implements json.Unmarshaler for TerminationConditionEnvelope.
func (e *TerminationConditionEnvelope) UnmarshalJSON(data []byte) error {
	var typeOnly struct {
		Type TerminationConditionType `json:"type"`
	}
	if err := json.Unmarshal(data, &typeOnly); err != nil {
		return fmt.Errorf("failed to parse termination condition type: %w", err)
	}
	e.Type = typeOnly.Type
	e.Raw = data
	return nil
}

// UnmarshalTerminationCondition deserializes a TerminationCondition from JSON.
func UnmarshalTerminationCondition(data json.RawMessage) (TerminationCondition, error) {
	var envelope TerminationConditionEnvelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, fmt.Errorf("failed to parse termination envelope: %w", err)
	}

	switch envelope.Type {
	case TerminationTypeRPEThreshold:
		var cond RPEThreshold
		if err := json.Unmarshal(data, &cond); err != nil {
			return nil, fmt.Errorf("failed to unmarshal RPEThreshold: %w", err)
		}
		if err := cond.Validate(); err != nil {
			return nil, err
		}
		return &cond, nil
	case TerminationTypeRepFailure:
		return NewRepFailure(), nil
	case TerminationTypeMaxSets:
		var cond MaxSets
		if err := json.Unmarshal(data, &cond); err != nil {
			return nil, fmt.Errorf("failed to unmarshal MaxSets: %w", err)
		}
		if err := cond.Validate(); err != nil {
			return nil, err
		}
		return &cond, nil
	default:
		return nil, fmt.Errorf("%w: unknown type %s", ErrInvalidTermination, envelope.Type)
	}
}

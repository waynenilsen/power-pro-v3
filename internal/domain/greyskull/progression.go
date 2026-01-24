// Package greyskull provides domain logic for the GreySkull LP program.
// This file implements the GreySkullProgression strategy that adjusts training weights
// based on AMRAP (As Many Reps As Possible) set performance using a three-tier system:
// deload, standard increment, and double increment.
package greyskull

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/domain/progression"
)

// GreySkullProgression implements the Progression interface for GreySkull LP's
// AMRAP-based weight progression. This progression type uniquely combines both
// positive progression and deload logic in a single progression type.
//
// The three-tier progression rules:
//   - If reps < MinReps (failure): Deload by DeloadPercent (e.g., 10%)
//   - If MinReps <= reps < DoubleThreshold: Add standard increment (WeightIncrement)
//   - If reps >= DoubleThreshold: Add double increment (WeightIncrement * 2)
//
// Example for main lifts (upper body):
//   - MinReps: 5, DoubleThreshold: 10, WeightIncrement: 2.5, DeloadPercent: 0.10
//   - 3 reps (< 5): Deload 10%
//   - 7 reps (5-9): Add 2.5 lbs
//   - 12 reps (>= 10): Add 5 lbs (double)
//
// Example for accessory lifts:
//   - MinReps: 10, DoubleThreshold: 15, WeightIncrement: 2.5, DeloadPercent: 0.10
//   - 8 reps (< 10): Deload 10%
//   - 12 reps (10-14): Add 2.5 lbs
//   - 18 reps (>= 15): Add 5 lbs (double)
type GreySkullProgression struct {
	// ID is the unique identifier for this progression.
	ID string `json:"id"`
	// Name is the human-readable name for this progression.
	Name string `json:"name"`
	// WeightIncrement is the standard weight increment (e.g., 2.5 for upper body, 5 for lower body).
	WeightIncrement float64 `json:"weightIncrement"`
	// MinReps is the minimum reps to avoid deload (5 for main lifts, 10 for accessories).
	MinReps int `json:"minReps"`
	// DoubleThreshold is the reps threshold to trigger double increment (10 for main, 15 for accessory).
	DoubleThreshold int `json:"doubleThreshold"`
	// DeloadPercent is the percentage to reduce weight on failure (e.g., 0.10 for 10%).
	DeloadPercent float64 `json:"deloadPercent"`
	// MaxTypeValue specifies which max to update (ONE_RM or TRAINING_MAX).
	MaxTypeValue progression.MaxType `json:"maxType"`
}

// NewGreySkullProgression creates a new GreySkullProgression with the given parameters.
func NewGreySkullProgression(id, name string, weightIncrement float64, minReps, doubleThreshold int, deloadPercent float64, maxType progression.MaxType) (*GreySkullProgression, error) {
	gsp := &GreySkullProgression{
		ID:              id,
		Name:            name,
		WeightIncrement: weightIncrement,
		MinReps:         minReps,
		DoubleThreshold: doubleThreshold,
		DeloadPercent:   deloadPercent,
		MaxTypeValue:    maxType,
	}
	if err := gsp.Validate(); err != nil {
		return nil, err
	}
	return gsp, nil
}

// NewGreySkullMainLiftProgression creates a GreySkullProgression configured for main lifts.
// Main lift defaults: MinReps=5, DoubleThreshold=10, DeloadPercent=0.10
func NewGreySkullMainLiftProgression(id, name string, weightIncrement float64, maxType progression.MaxType) (*GreySkullProgression, error) {
	return NewGreySkullProgression(id, name, weightIncrement, 5, 10, 0.10, maxType)
}

// NewGreySkullAccessoryProgression creates a GreySkullProgression configured for accessory lifts.
// Accessory defaults: MinReps=10, DoubleThreshold=15, DeloadPercent=0.10
func NewGreySkullAccessoryProgression(id, name string, weightIncrement float64, maxType progression.MaxType) (*GreySkullProgression, error) {
	return NewGreySkullProgression(id, name, weightIncrement, 10, 15, 0.10, maxType)
}

// Type returns the discriminator string for this progression.
// Implements Progression interface.
func (g *GreySkullProgression) Type() progression.ProgressionType {
	return progression.TypeGreySkull
}

// TriggerType returns the trigger type this progression responds to.
// GreySkullProgression fires AFTER_SET since it evaluates AMRAP performance.
// Implements Progression interface.
func (g *GreySkullProgression) TriggerType() progression.TriggerType {
	return progression.TriggerAfterSet
}

// Validate validates the progression's configuration parameters.
// Implements Progression interface.
func (g *GreySkullProgression) Validate() error {
	if g.ID == "" {
		return fmt.Errorf("%w: id is required", progression.ErrInvalidParams)
	}
	if g.Name == "" {
		return fmt.Errorf("%w: name is required", progression.ErrInvalidParams)
	}
	if g.WeightIncrement <= 0 {
		return fmt.Errorf("%w: weightIncrement must be positive", progression.ErrInvalidParams)
	}
	if g.MinReps < 1 {
		return fmt.Errorf("%w: minReps must be at least 1", progression.ErrInvalidParams)
	}
	if g.DoubleThreshold <= g.MinReps {
		return fmt.Errorf("%w: doubleThreshold must be greater than minReps", progression.ErrInvalidParams)
	}
	if g.DeloadPercent <= 0 || g.DeloadPercent > 1 {
		return fmt.Errorf("%w: deloadPercent must be between 0 (exclusive) and 1 (inclusive)", progression.ErrInvalidParams)
	}
	if err := progression.ValidateMaxType(g.MaxTypeValue); err != nil {
		return err
	}
	return nil
}

// Apply evaluates and applies the progression given the context.
// Implements Progression interface.
//
// GreySkullProgression implements the GreySkull LP AMRAP-based progression pattern.
// The key characteristics are:
//
//  1. Trigger type must be AFTER_SET - fires immediately after logging an AMRAP set
//  2. The TriggerEvent must include RepsPerformed and IsAMRAP=true
//  3. Three-tier progression based on reps performed:
//     a. reps < MinReps: Deload (negative delta)
//     b. MinReps <= reps < DoubleThreshold: Standard increment
//     c. reps >= DoubleThreshold: Double increment
//
// This approach allows the program to automatically handle both progression
// and deload logic based on AMRAP performance in a single progression type.
func (g *GreySkullProgression) Apply(ctx context.Context, params progression.ProgressionContext) (progression.ProgressionResult, error) {
	if err := params.Validate(); err != nil {
		return progression.ProgressionResult{}, fmt.Errorf("invalid progression context: %w", err)
	}

	now := time.Now()

	// Trigger type must be AFTER_SET
	if params.TriggerEvent.Type != progression.TriggerAfterSet {
		return progression.ProgressionResult{
			Applied:       false,
			PreviousValue: params.CurrentValue,
			NewValue:      params.CurrentValue,
			Delta:         0,
			LiftID:        params.LiftID,
			MaxType:       params.MaxType,
			AppliedAt:     now,
			Reason:        fmt.Sprintf("trigger type mismatch: expected %s, got %s", progression.TriggerAfterSet, params.TriggerEvent.Type),
		}, nil
	}

	// Max type must match
	if params.MaxType != g.MaxTypeValue {
		return progression.ProgressionResult{
			Applied:       false,
			PreviousValue: params.CurrentValue,
			NewValue:      params.CurrentValue,
			Delta:         0,
			LiftID:        params.LiftID,
			MaxType:       params.MaxType,
			AppliedAt:     now,
			Reason:        fmt.Sprintf("max type mismatch: expected %s, got %s", g.MaxTypeValue, params.MaxType),
		}, nil
	}

	// Must be an AMRAP set
	if !params.TriggerEvent.IsAMRAP {
		return progression.ProgressionResult{
			Applied:       false,
			PreviousValue: params.CurrentValue,
			NewValue:      params.CurrentValue,
			Delta:         0,
			LiftID:        params.LiftID,
			MaxType:       params.MaxType,
			AppliedAt:     now,
			Reason:        "set is not marked as AMRAP",
		}, nil
	}

	// RepsPerformed must be provided
	if params.TriggerEvent.RepsPerformed == nil {
		return progression.ProgressionResult{
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

	reps := *params.TriggerEvent.RepsPerformed

	var newValue float64
	var delta float64

	if reps < g.MinReps {
		// Deload: reduce by DeloadPercent
		deloadAmount := params.CurrentValue * g.DeloadPercent
		newValue = params.CurrentValue - deloadAmount
		delta = -deloadAmount

		// Ensure we don't go below zero
		if newValue < 0 {
			newValue = 0
			delta = -params.CurrentValue
		}
	} else if reps >= g.DoubleThreshold {
		// Double increment
		delta = g.WeightIncrement * 2
		newValue = params.CurrentValue + delta
	} else {
		// Standard increment (MinReps <= reps < DoubleThreshold)
		delta = g.WeightIncrement
		newValue = params.CurrentValue + delta
	}

	return progression.ProgressionResult{
		Applied:       true,
		PreviousValue: params.CurrentValue,
		NewValue:      newValue,
		Delta:         delta,
		LiftID:        params.LiftID,
		MaxType:       params.MaxType,
		AppliedAt:     now,
	}, nil
}

// MarshalJSON implements json.Marshaler for GreySkullProgression.
// Ensures the type discriminator is always included in serialized output.
func (g *GreySkullProgression) MarshalJSON() ([]byte, error) {
	type Alias GreySkullProgression
	return json.Marshal(&struct {
		Type progression.ProgressionType `json:"type"`
		*Alias
	}{
		Type:  progression.TypeGreySkull,
		Alias: (*Alias)(g),
	})
}

// UnmarshalGreySkullProgression deserializes a GreySkullProgression from JSON.
// This is used by the ProgressionFactory for type-safe deserialization.
func UnmarshalGreySkullProgression(data json.RawMessage) (progression.Progression, error) {
	var gsp GreySkullProgression
	if err := json.Unmarshal(data, &gsp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal greyskull progression: %w", err)
	}
	if err := gsp.Validate(); err != nil {
		return nil, fmt.Errorf("invalid greyskull progression: %w", err)
	}
	return &gsp, nil
}

// RegisterGreySkullProgression registers the GreySkullProgression type with a factory.
// This should be called during application initialization.
func RegisterGreySkullProgression(factory *progression.ProgressionFactory) {
	factory.Register(progression.TypeGreySkull, UnmarshalGreySkullProgression)
}

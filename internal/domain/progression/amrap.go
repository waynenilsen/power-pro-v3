// Package progression provides domain logic for progression strategies.
// This file implements the AMRAPProgression strategy that adjusts training maxes
// based on AMRAP (As Many Reps As Possible) set performance.
package progression

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

// RepsThreshold defines a threshold for AMRAP-based progression.
// When reps performed >= MinReps, the Increment is applied.
type RepsThreshold struct {
	// MinReps is the minimum number of reps to trigger this threshold.
	MinReps int `json:"minReps"`
	// Increment is the weight to add when this threshold is met.
	Increment float64 `json:"increment"`
}

// AMRAPProgression implements the Progression interface for AMRAP-based weight increases.
// This is used by programs like nSuns where the weight increase depends on how many
// reps the lifter achieves on their AMRAP set.
//
// Example nSuns-style thresholds:
//
//	{minReps: 2, increment: 5.0}   - 2-3 reps: +5lb
//	{minReps: 4, increment: 10.0}  - 4-5 reps: +10lb
//	{minReps: 6, increment: 15.0}  - 6+ reps: +15lb
//
// The Apply logic finds the highest threshold where repsPerformed >= minReps,
// then applies that threshold's increment. If no threshold is met, delta is 0.
type AMRAPProgression struct {
	// ID is the unique identifier for this progression.
	ID string `json:"id"`
	// Name is the human-readable name for this progression.
	Name string `json:"name"`
	// MaxTypeValue specifies which max to update (ONE_RM or TRAINING_MAX).
	MaxTypeValue MaxType `json:"maxType"`
	// TriggerTypeValue should be TriggerAfterSet for AMRAP progressions.
	TriggerTypeValue TriggerType `json:"triggerType"`
	// Thresholds defines the reps-to-increment mapping.
	// Must be sorted by MinReps ascending.
	Thresholds []RepsThreshold `json:"thresholds"`
}

// NewAMRAPProgression creates a new AMRAPProgression with the given parameters.
func NewAMRAPProgression(id, name string, maxType MaxType, triggerType TriggerType, thresholds []RepsThreshold) (*AMRAPProgression, error) {
	ap := &AMRAPProgression{
		ID:               id,
		Name:             name,
		MaxTypeValue:     maxType,
		TriggerTypeValue: triggerType,
		Thresholds:       thresholds,
	}
	if err := ap.Validate(); err != nil {
		return nil, err
	}
	return ap, nil
}

// Type returns the discriminator string for this progression.
// Implements Progression interface.
func (a *AMRAPProgression) Type() ProgressionType {
	return TypeAMRAP
}

// TriggerType returns the trigger type this progression responds to.
// Implements Progression interface.
func (a *AMRAPProgression) TriggerType() TriggerType {
	return a.TriggerTypeValue
}

// Validate validates the progression's configuration parameters.
// Implements Progression interface.
func (a *AMRAPProgression) Validate() error {
	if a.ID == "" {
		return fmt.Errorf("%w: id is required", ErrInvalidParams)
	}
	if a.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidParams)
	}
	if err := ValidateMaxType(a.MaxTypeValue); err != nil {
		return err
	}
	if err := ValidateTriggerType(a.TriggerTypeValue); err != nil {
		return err
	}
	// AMRAPProgression typically uses AFTER_SET trigger
	if a.TriggerTypeValue != TriggerAfterSet {
		return fmt.Errorf("%w: AMRAP progression requires AFTER_SET trigger type", ErrInvalidParams)
	}
	if len(a.Thresholds) == 0 {
		return fmt.Errorf("%w: at least one threshold is required", ErrInvalidParams)
	}

	// Validate thresholds: sorted by minReps ascending, positive increments
	for i, t := range a.Thresholds {
		if t.MinReps < 0 {
			return fmt.Errorf("%w: threshold[%d].minReps must be non-negative", ErrInvalidParams, i)
		}
		if t.Increment <= 0 {
			return fmt.Errorf("%w: threshold[%d].increment must be positive", ErrInvalidParams, i)
		}
		if i > 0 && t.MinReps <= a.Thresholds[i-1].MinReps {
			return fmt.Errorf("%w: thresholds must be sorted by minReps ascending (threshold[%d].minReps=%d <= threshold[%d].minReps=%d)",
				ErrInvalidParams, i, t.MinReps, i-1, a.Thresholds[i-1].MinReps)
		}
	}

	return nil
}

// Apply evaluates and applies the progression given the context.
// Implements Progression interface.
//
// AMRAPProgression implements the AMRAP-based progression pattern used in programs
// like nSuns. The key characteristics are:
//
//  1. Trigger type must be AFTER_SET - fires immediately after logging an AMRAP set
//  2. The TriggerEvent must include RepsPerformed and IsAMRAP=true
//  3. The increment depends on how many reps were performed:
//     - Find the highest threshold where repsPerformed >= minReps
//     - Apply that threshold's increment
//     - If no threshold is met, apply 0 (no progression)
//
// This approach allows programs to reward better performance with larger jumps,
// encouraging lifters to push for more reps while ensuring sustainable progression.
func (a *AMRAPProgression) Apply(ctx context.Context, params ProgressionContext) (ProgressionResult, error) {
	if err := params.Validate(); err != nil {
		return ProgressionResult{}, fmt.Errorf("invalid progression context: %w", err)
	}

	now := time.Now()

	// Trigger type must match
	if params.TriggerEvent.Type != a.TriggerTypeValue {
		return ProgressionResult{
			Applied:       false,
			PreviousValue: params.CurrentValue,
			NewValue:      params.CurrentValue,
			Delta:         0,
			LiftID:        params.LiftID,
			MaxType:       params.MaxType,
			AppliedAt:     now,
			Reason:        fmt.Sprintf("trigger type mismatch: expected %s, got %s", a.TriggerTypeValue, params.TriggerEvent.Type),
		}, nil
	}

	// Max type must match
	if params.MaxType != a.MaxTypeValue {
		return ProgressionResult{
			Applied:       false,
			PreviousValue: params.CurrentValue,
			NewValue:      params.CurrentValue,
			Delta:         0,
			LiftID:        params.LiftID,
			MaxType:       params.MaxType,
			AppliedAt:     now,
			Reason:        fmt.Sprintf("max type mismatch: expected %s, got %s", a.MaxTypeValue, params.MaxType),
		}, nil
	}

	// Must be an AMRAP set
	if !params.TriggerEvent.IsAMRAP {
		return ProgressionResult{
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

	reps := *params.TriggerEvent.RepsPerformed

	// Find the highest threshold where reps >= minReps
	// Thresholds are sorted ascending, so we iterate in reverse to find the highest match
	var increment float64 = 0
	var thresholdMet bool = false
	for i := len(a.Thresholds) - 1; i >= 0; i-- {
		if reps >= a.Thresholds[i].MinReps {
			increment = a.Thresholds[i].Increment
			thresholdMet = true
			break
		}
	}

	if !thresholdMet {
		return ProgressionResult{
			Applied:       false,
			PreviousValue: params.CurrentValue,
			NewValue:      params.CurrentValue,
			Delta:         0,
			LiftID:        params.LiftID,
			MaxType:       params.MaxType,
			AppliedAt:     now,
			Reason:        fmt.Sprintf("no threshold met: reps=%d, minimum required=%d", reps, a.Thresholds[0].MinReps),
		}, nil
	}

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

// MarshalJSON implements json.Marshaler for AMRAPProgression.
// Ensures the type discriminator is always included in serialized output.
func (a *AMRAPProgression) MarshalJSON() ([]byte, error) {
	type Alias AMRAPProgression
	return json.Marshal(&struct {
		Type ProgressionType `json:"type"`
		*Alias
	}{
		Type:  TypeAMRAP,
		Alias: (*Alias)(a),
	})
}

// UnmarshalAMRAPProgression deserializes an AMRAPProgression from JSON.
// This is used by the ProgressionFactory for type-safe deserialization.
func UnmarshalAMRAPProgression(data json.RawMessage) (Progression, error) {
	var ap AMRAPProgression
	if err := json.Unmarshal(data, &ap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal AMRAP progression: %w", err)
	}

	// Sort thresholds by minReps ascending if not already sorted
	// (defensive, validation will catch unsorted input but this makes it more robust)
	sort.Slice(ap.Thresholds, func(i, j int) bool {
		return ap.Thresholds[i].MinReps < ap.Thresholds[j].MinReps
	})

	if err := ap.Validate(); err != nil {
		return nil, fmt.Errorf("invalid AMRAP progression: %w", err)
	}
	return &ap, nil
}

// RegisterAMRAPProgression registers the AMRAPProgression type with a factory.
// This should be called during application initialization.
func RegisterAMRAPProgression(factory *ProgressionFactory) {
	factory.Register(TypeAMRAP, UnmarshalAMRAPProgression)
}

// Package setscheme provides domain logic for set/rep scheme strategies.
package setscheme

import (
	"encoding/json"
	"fmt"

	"github.com/waynenilsen/power-pro-v3/internal/domain/loadstrategy"
)

// FatigueDrop implements RTS-style load drop training where sets continue
// at progressively lower weights until a target RPE is reached.
//
// Example: Squat @ 3 reps, start at RPE 8, drop 5%, stop at RPE 10
//  1. Set 1: 315 lbs x 3 @ RPE 8.0 (target achieved)
//  2. Set 2: 299 lbs x 3 @ RPE 8.5 (dropped 5%, continue)
//  3. Set 3: 284 lbs x 3 @ RPE 9.0 (continue)
//  4. Set 4: 270 lbs x 3 @ RPE 9.5 (continue)
//  5. Set 5: 256 lbs x 3 @ RPE 10 (STOP - target RPE reached)
type FatigueDrop struct {
	// TargetReps is the number of repetitions per set (required, must be >= 1).
	TargetReps int `json:"target_reps"`
	// StartRPE is the initial target RPE for the first set (required, 1-10).
	StartRPE float64 `json:"start_rpe"`
	// StopRPE is the RPE threshold at which to stop (required, must be > StartRPE).
	StopRPE float64 `json:"stop_rpe"`
	// DropPercent is the percentage to drop weight each set (e.g., 0.05 for 5%).
	DropPercent float64 `json:"drop_percent"`
	// MaxSets is the safety limit for maximum number of sets (default 10 if 0).
	MaxSets int `json:"max_sets,omitempty"`
}

// DefaultMaxSets is the default safety limit for FatigueDrop.
const DefaultMaxSets = 10

// NewFatigueDrop creates a new FatigueDrop set scheme.
// Returns an error if validation fails.
func NewFatigueDrop(targetReps int, startRPE, stopRPE, dropPercent float64, maxSets int) (*FatigueDrop, error) {
	scheme := &FatigueDrop{
		TargetReps:  targetReps,
		StartRPE:    startRPE,
		StopRPE:     stopRPE,
		DropPercent: dropPercent,
		MaxSets:     maxSets,
	}
	if err := scheme.Validate(); err != nil {
		return nil, err
	}
	return scheme, nil
}

// Type returns the discriminator string for this scheme.
func (f *FatigueDrop) Type() SetSchemeType {
	return TypeFatigueDrop
}

// GenerateSets generates the first set at the starting weight.
// For variable schemes, this returns only the first (provisional) set.
// Subsequent sets are generated via GenerateNextSet based on session performance.
func (f *FatigueDrop) GenerateSets(baseWeight float64, _ SetGenerationContext) ([]GeneratedSet, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}

	// Return only the first set, marked as provisional
	return []GeneratedSet{
		{
			SetNumber:     1,
			Weight:        baseWeight,
			TargetReps:    f.TargetReps,
			IsWorkSet:     true,
			IsProvisional: true,
		},
	}, nil
}

// Validate validates the scheme's configuration parameters.
func (f *FatigueDrop) Validate() error {
	if f.TargetReps < 1 {
		return fmt.Errorf("%w: target_reps must be >= 1, got %d", ErrInvalidParams, f.TargetReps)
	}
	if f.StartRPE < 1 || f.StartRPE > 10 {
		return fmt.Errorf("%w: start_rpe must be between 1 and 10, got %v", ErrInvalidParams, f.StartRPE)
	}
	if f.StopRPE < 1 || f.StopRPE > 10 {
		return fmt.Errorf("%w: stop_rpe must be between 1 and 10, got %v", ErrInvalidParams, f.StopRPE)
	}
	if f.StopRPE <= f.StartRPE {
		return fmt.Errorf("%w: stop_rpe (%v) must be greater than start_rpe (%v)", ErrInvalidParams, f.StopRPE, f.StartRPE)
	}
	if f.DropPercent < 0 || f.DropPercent > 1 {
		return fmt.Errorf("%w: drop_percent must be between 0 and 1, got %v", ErrInvalidParams, f.DropPercent)
	}
	if f.MaxSets < 0 {
		return fmt.Errorf("%w: max_sets must be >= 0, got %d", ErrInvalidParams, f.MaxSets)
	}
	return nil
}

// IsVariableCount returns true, indicating this scheme has variable set counts.
func (f *FatigueDrop) IsVariableCount() bool {
	return true
}

// GetTerminationCondition returns the RPE threshold condition for termination.
func (f *FatigueDrop) GetTerminationCondition() TerminationCondition {
	return &RPEThreshold{Threshold: f.StopRPE}
}

// getEffectiveMaxSets returns the max sets limit, applying default if not set.
func (f *FatigueDrop) getEffectiveMaxSets() int {
	if f.MaxSets == 0 {
		return DefaultMaxSets
	}
	return f.MaxSets
}

// GenerateNextSet generates the next set based on history and termination context.
// Returns the next set and true if generation should continue,
// or nil and false if the termination condition is met.
func (f *FatigueDrop) GenerateNextSet(ctx SetGenerationContext, history []GeneratedSet, termCtx TerminationContext) (*GeneratedSet, bool) {
	// Check termination conditions first
	// 1. Check if RPE threshold is met (primary termination)
	if f.GetTerminationCondition().ShouldTerminate(termCtx) {
		return nil, false
	}

	// 2. Check max sets safety limit
	if termCtx.TotalSets >= f.getEffectiveMaxSets() {
		return nil, false
	}

	// Calculate next set number
	nextSetNumber := termCtx.TotalSets + 1

	// Calculate next weight
	var nextWeight float64
	if len(history) == 0 {
		// This shouldn't happen in normal flow (GenerateSets gives first set),
		// but handle gracefully by returning nil
		return nil, false
	}

	// Get the last set's weight and apply the drop
	lastWeight := history[len(history)-1].Weight
	nextWeight = lastWeight * (1 - f.DropPercent)

	// Apply rounding (round down to be conservative with fatigue)
	roundedWeight, err := loadstrategy.RoundWeightDown(nextWeight, loadstrategy.DefaultRoundingIncrement)
	if err == nil {
		nextWeight = roundedWeight
	}
	// If rounding fails, use the unrounded weight

	// Ensure weight doesn't go negative
	if nextWeight <= 0 {
		return nil, false
	}

	return &GeneratedSet{
		SetNumber:     nextSetNumber,
		Weight:        nextWeight,
		TargetReps:    f.TargetReps,
		IsWorkSet:     true,
		IsProvisional: true,
	}, true
}

// MarshalJSON implements json.Marshaler for FatigueDrop.
// Includes the type discriminator for polymorphic deserialization.
func (f *FatigueDrop) MarshalJSON() ([]byte, error) {
	type Alias FatigueDrop
	return json.Marshal(&struct {
		Type SetSchemeType `json:"type"`
		*Alias
	}{
		Type:  TypeFatigueDrop,
		Alias: (*Alias)(f),
	})
}

// UnmarshalFatigueDrop deserializes a FatigueDrop from JSON.
// This is used by the SchemeFactory.
func UnmarshalFatigueDrop(data json.RawMessage) (SetScheme, error) {
	var scheme FatigueDrop
	if err := json.Unmarshal(data, &scheme); err != nil {
		return nil, fmt.Errorf("failed to unmarshal FatigueDrop: %w", err)
	}
	if err := scheme.Validate(); err != nil {
		return nil, err
	}
	return &scheme, nil
}

// RegisterFatigueDrop registers the FatigueDrop scheme with the given factory.
func RegisterFatigueDrop(factory *SchemeFactory) {
	factory.Register(TypeFatigueDrop, UnmarshalFatigueDrop)
}

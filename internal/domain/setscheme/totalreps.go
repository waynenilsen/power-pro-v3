// Package setscheme provides domain logic for set/rep scheme strategies.
package setscheme

import (
	"encoding/json"
	"fmt"
)

// TotalRepsScheme implements a variable set scheme where users accumulate reps
// until reaching a target total. Unlike MRS, there's no minimum reps requirement
// per set - users can distribute reps however they want.
//
// Example: "100 chin-ups" from 5/3/1 Building the Monolith
//  1. Set 1: 15 reps (total: 15, continue)
//  2. Set 2: 12 reps (total: 27, continue)
//  3. Set 3: 10 reps (total: 37, continue)
//  ...continue until total >= 100 or max sets reached...
type TotalRepsScheme struct {
	// TargetTotalReps is the cumulative rep target to reach (required, >= 1).
	// Generation stops when total reps >= this value.
	TargetTotalReps int `json:"target_total_reps"`
	// SuggestedRepsPerSet is a hint for initial set size (optional, >= 1 if set).
	// This is purely advisory and helps users plan their first set.
	SuggestedRepsPerSet int `json:"suggested_reps_per_set,omitempty"`
	// MaxSets is the safety limit for maximum number of sets (default 20 if 0).
	// Higher than MRS since accessory volume work may require many sets.
	MaxSets int `json:"max_sets,omitempty"`
}

// DefaultTotalRepsMaxSets is the default safety limit for TotalReps.
// Higher than MRS (10) since accessory work often requires more sets.
const DefaultTotalRepsMaxSets = 20

// NewTotalRepsScheme creates a new TotalReps set scheme.
// Returns an error if validation fails.
func NewTotalRepsScheme(targetTotalReps, suggestedRepsPerSet, maxSets int) (*TotalRepsScheme, error) {
	scheme := &TotalRepsScheme{
		TargetTotalReps:     targetTotalReps,
		SuggestedRepsPerSet: suggestedRepsPerSet,
		MaxSets:             maxSets,
	}
	if err := scheme.Validate(); err != nil {
		return nil, err
	}
	return scheme, nil
}

// Type returns the discriminator string for this scheme.
func (t *TotalRepsScheme) Type() SetSchemeType {
	return TypeTotalReps
}

// GenerateSets generates the first set at the given weight.
// For variable schemes, this returns only the first (provisional) set.
// Subsequent sets are generated via GenerateNextSet based on session performance.
func (t *TotalRepsScheme) GenerateSets(baseWeight float64, _ SetGenerationContext) ([]GeneratedSet, error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}

	// Use suggested reps if provided, otherwise use a reasonable default
	targetReps := t.SuggestedRepsPerSet
	if targetReps == 0 {
		targetReps = 10 // Default suggestion
	}

	// Return only the first set, marked as provisional.
	return []GeneratedSet{
		{
			SetNumber:     1,
			Weight:        baseWeight,
			TargetReps:    targetReps,
			IsWorkSet:     true,
			IsProvisional: true,
		},
	}, nil
}

// Validate validates the scheme's configuration parameters.
func (t *TotalRepsScheme) Validate() error {
	if t.TargetTotalReps < 1 {
		return fmt.Errorf("%w: target_total_reps must be >= 1, got %d", ErrInvalidParams, t.TargetTotalReps)
	}
	if t.SuggestedRepsPerSet < 0 {
		return fmt.Errorf("%w: suggested_reps_per_set must be >= 0, got %d", ErrInvalidParams, t.SuggestedRepsPerSet)
	}
	if t.MaxSets < 0 {
		return fmt.Errorf("%w: max_sets must be >= 0, got %d", ErrInvalidParams, t.MaxSets)
	}
	return nil
}

// IsVariableCount returns true, indicating this scheme has variable set counts.
func (t *TotalRepsScheme) IsVariableCount() bool {
	return true
}

// GetTerminationCondition returns the total reps threshold condition for termination.
func (t *TotalRepsScheme) GetTerminationCondition() TerminationCondition {
	return &TotalReps{Target: t.TargetTotalReps}
}

// getEffectiveMaxSets returns the max sets limit, applying default if not set.
func (t *TotalRepsScheme) getEffectiveMaxSets() int {
	if t.MaxSets == 0 {
		return DefaultTotalRepsMaxSets
	}
	return t.MaxSets
}

// getEffectiveSuggestedReps returns the suggested reps, applying default if not set.
func (t *TotalRepsScheme) getEffectiveSuggestedReps() int {
	if t.SuggestedRepsPerSet == 0 {
		return 10 // Default suggestion
	}
	return t.SuggestedRepsPerSet
}

// GenerateNextSet generates the next set based on history and termination context.
// Returns the next set and true if generation should continue,
// or nil and false if the termination condition is met.
//
// Termination occurs when EITHER of:
// 1. TotalReps >= TargetTotalReps (success - hit target)
// 2. TotalSets >= MaxSets (safety limit)
//
// Unlike MRS, there is NO rep failure condition - users can do sets of any size.
func (t *TotalRepsScheme) GenerateNextSet(ctx SetGenerationContext, history []GeneratedSet, termCtx TerminationContext) (*GeneratedSet, bool) {
	// Check termination conditions

	// 1. Check if total reps target is met (primary termination - success)
	if t.GetTerminationCondition().ShouldTerminate(termCtx) {
		return nil, false
	}

	// 2. Check max sets safety limit
	if termCtx.TotalSets >= t.getEffectiveMaxSets() {
		return nil, false
	}

	// Calculate next set number
	nextSetNumber := termCtx.TotalSets + 1

	// Get weight from history - use same weight for all sets
	var weight float64
	if len(history) == 0 {
		// This shouldn't happen in normal flow (GenerateSets gives first set),
		// but handle gracefully by returning nil
		return nil, false
	}
	weight = history[0].Weight // Use first set's weight (all sets same weight)

	return &GeneratedSet{
		SetNumber:     nextSetNumber,
		Weight:        weight,
		TargetReps:    t.getEffectiveSuggestedReps(),
		IsWorkSet:     true,
		IsProvisional: true,
	}, true
}

// MarshalJSON implements json.Marshaler for TotalRepsScheme.
// Includes the type discriminator for polymorphic deserialization.
func (t *TotalRepsScheme) MarshalJSON() ([]byte, error) {
	type Alias TotalRepsScheme
	return json.Marshal(&struct {
		Type SetSchemeType `json:"type"`
		*Alias
	}{
		Type:  TypeTotalReps,
		Alias: (*Alias)(t),
	})
}

// UnmarshalTotalRepsScheme deserializes a TotalRepsScheme from JSON.
// This is used by the SchemeFactory.
func UnmarshalTotalRepsScheme(data json.RawMessage) (SetScheme, error) {
	var scheme TotalRepsScheme
	if err := json.Unmarshal(data, &scheme); err != nil {
		return nil, fmt.Errorf("failed to unmarshal TotalRepsScheme: %w", err)
	}
	if err := scheme.Validate(); err != nil {
		return nil, err
	}
	return &scheme, nil
}

// RegisterTotalRepsScheme registers the TotalReps scheme with the given factory.
func RegisterTotalRepsScheme(factory *SchemeFactory) {
	factory.Register(TypeTotalReps, UnmarshalTotalRepsScheme)
}

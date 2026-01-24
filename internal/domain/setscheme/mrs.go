// Package setscheme provides domain logic for set/rep scheme strategies.
package setscheme

import (
	"encoding/json"
	"fmt"
)

// MRS implements GZCL-style Max Rep Sets training where sets continue
// at a fixed weight until a target total rep count is reached or
// the lifter fails to hit minimum reps.
//
// Example: Bench Press MRS x 3, target 25 total reps
//  1. Set 1: 225 lbs x 10 (total: 10, continue)
//  2. Set 2: 225 lbs x 8 (total: 18, continue)
//  3. Set 3: 225 lbs x 5 (total: 23, continue)
//  4. Set 4: 225 lbs x 4 (total: 27, STOP - exceeded 25)
//
// Alternative termination - failure:
//  1. Set 1: 225 lbs x 10 (total: 10)
//  2. Set 2: 225 lbs x 6 (total: 16)
//  3. Set 3: 225 lbs x 2 (STOP - failed to hit minimum of 3)
type MRS struct {
	// TargetTotalReps is the cumulative rep target to reach (required, >= 1).
	// Generation stops when total reps >= this value.
	TargetTotalReps int `json:"target_total_reps"`
	// MinRepsPerSet is the minimum reps required per set (required, >= 1).
	// If a set fails to reach this, it's considered technical failure.
	MinRepsPerSet int `json:"min_reps_per_set"`
	// MaxSets is the safety limit for maximum number of sets (default 10 if 0).
	MaxSets int `json:"max_sets,omitempty"`
	// NumberOfMRS is how many MRS blocks in the session (for GZCL: T1=3, T3=4).
	// This field is for metadata/planning purposes.
	NumberOfMRS int `json:"number_of_mrs,omitempty"`
}

// DefaultMRSMaxSets is the default safety limit for MRS.
const DefaultMRSMaxSets = 10

// NewMRS creates a new MRS set scheme.
// Returns an error if validation fails.
func NewMRS(targetTotalReps, minRepsPerSet, maxSets, numberOfMRS int) (*MRS, error) {
	scheme := &MRS{
		TargetTotalReps: targetTotalReps,
		MinRepsPerSet:   minRepsPerSet,
		MaxSets:         maxSets,
		NumberOfMRS:     numberOfMRS,
	}
	if err := scheme.Validate(); err != nil {
		return nil, err
	}
	return scheme, nil
}

// Type returns the discriminator string for this scheme.
func (m *MRS) Type() SetSchemeType {
	return TypeMRS
}

// GenerateSets generates the first set at the given weight.
// For variable schemes, this returns only the first (provisional) set.
// Subsequent sets are generated via GenerateNextSet based on session performance.
func (m *MRS) GenerateSets(baseWeight float64, _ SetGenerationContext) ([]GeneratedSet, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}

	// Return only the first set, marked as provisional.
	// TargetReps is set to MinRepsPerSet as the minimum expectation.
	return []GeneratedSet{
		{
			SetNumber:     1,
			Weight:        baseWeight,
			TargetReps:    m.MinRepsPerSet,
			IsWorkSet:     true,
			IsProvisional: true,
		},
	}, nil
}

// Validate validates the scheme's configuration parameters.
func (m *MRS) Validate() error {
	if m.TargetTotalReps < 1 {
		return fmt.Errorf("%w: target_total_reps must be >= 1, got %d", ErrInvalidParams, m.TargetTotalReps)
	}
	if m.MinRepsPerSet < 1 {
		return fmt.Errorf("%w: min_reps_per_set must be >= 1, got %d", ErrInvalidParams, m.MinRepsPerSet)
	}
	if m.TargetTotalReps < m.MinRepsPerSet {
		return fmt.Errorf("%w: target_total_reps (%d) must be >= min_reps_per_set (%d)",
			ErrInvalidParams, m.TargetTotalReps, m.MinRepsPerSet)
	}
	if m.MaxSets < 0 {
		return fmt.Errorf("%w: max_sets must be >= 0, got %d", ErrInvalidParams, m.MaxSets)
	}
	if m.NumberOfMRS < 0 {
		return fmt.Errorf("%w: number_of_mrs must be >= 0, got %d", ErrInvalidParams, m.NumberOfMRS)
	}
	return nil
}

// IsVariableCount returns true, indicating this scheme has variable set counts.
func (m *MRS) IsVariableCount() bool {
	return true
}

// GetTerminationCondition returns the total reps threshold condition for termination.
// Note: This returns the primary condition. Additional conditions (MinReps, MaxSets)
// are checked in GenerateNextSet.
func (m *MRS) GetTerminationCondition() TerminationCondition {
	return &TotalReps{Target: m.TargetTotalReps}
}

// getEffectiveMaxSets returns the max sets limit, applying default if not set.
func (m *MRS) getEffectiveMaxSets() int {
	if m.MaxSets == 0 {
		return DefaultMRSMaxSets
	}
	return m.MaxSets
}

// GenerateNextSet generates the next set based on history and termination context.
// Returns the next set and true if generation should continue,
// or nil and false if the termination condition is met.
//
// Termination occurs when ANY of:
// 1. TotalReps >= TargetTotalReps (success - hit target)
// 2. LastReps < MinRepsPerSet (failure - couldn't hit minimum)
// 3. TotalSets >= MaxSets (safety limit)
func (m *MRS) GenerateNextSet(ctx SetGenerationContext, history []GeneratedSet, termCtx TerminationContext) (*GeneratedSet, bool) {
	// Check termination conditions

	// 1. Check if total reps target is met (primary termination - success)
	if m.GetTerminationCondition().ShouldTerminate(termCtx) {
		return nil, false
	}

	// 2. Check if last set failed to hit minimum reps (failure termination)
	if termCtx.LastReps < m.MinRepsPerSet && termCtx.TotalSets > 0 {
		return nil, false
	}

	// 3. Check max sets safety limit
	if termCtx.TotalSets >= m.getEffectiveMaxSets() {
		return nil, false
	}

	// Calculate next set number
	nextSetNumber := termCtx.TotalSets + 1

	// MRS uses same weight for all sets - get it from history
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
		TargetReps:    m.MinRepsPerSet,
		IsWorkSet:     true,
		IsProvisional: true,
	}, true
}

// MarshalJSON implements json.Marshaler for MRS.
// Includes the type discriminator for polymorphic deserialization.
func (m *MRS) MarshalJSON() ([]byte, error) {
	type Alias MRS
	return json.Marshal(&struct {
		Type SetSchemeType `json:"type"`
		*Alias
	}{
		Type:  TypeMRS,
		Alias: (*Alias)(m),
	})
}

// UnmarshalMRS deserializes a MRS from JSON.
// This is used by the SchemeFactory.
func UnmarshalMRS(data json.RawMessage) (SetScheme, error) {
	var scheme MRS
	if err := json.Unmarshal(data, &scheme); err != nil {
		return nil, fmt.Errorf("failed to unmarshal MRS: %w", err)
	}
	if err := scheme.Validate(); err != nil {
		return nil, err
	}
	return &scheme, nil
}

// RegisterMRS registers the MRS scheme with the given factory.
func RegisterMRS(factory *SchemeFactory) {
	factory.Register(TypeMRS, UnmarshalMRS)
}

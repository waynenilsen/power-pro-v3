// Package setscheme provides domain logic for set/rep scheme strategies.
package setscheme

import (
	"encoding/json"
	"fmt"
)

// GreySkullSetScheme generates sets for GreySkull LP style workouts.
// GreySkull LP uses a unique set scheme of fixed sets followed by an AMRAP set.
// Example: 2x5 + 1x5+ means 2 sets of 5 reps, then 1 AMRAP set with minimum 5 reps.
// This is a composite scheme that combines fixed sets with a final AMRAP set.
type GreySkullSetScheme struct {
	// FixedSets is the number of fixed sets to generate (required, must be >= 0).
	// Typically 2 for GreySkull LP main lifts.
	FixedSets int `json:"fixedSets"`
	// FixedReps is the number of repetitions per fixed set (required, must be >= 1).
	// Typically 5 for GreySkull LP main lifts.
	FixedReps int `json:"fixedReps"`
	// AMRAPSets is the number of AMRAP sets to generate (required, must be >= 1).
	// Typically 1 for GreySkull LP.
	AMRAPSets int `json:"amrapSets"`
	// MinAMRAPReps is the minimum expected reps for the AMRAP set (required, must be >= 1).
	// Typically 5 for GreySkull LP main lifts.
	MinAMRAPReps int `json:"minAmrapReps"`
}

// NewGreySkullSetScheme creates a new GreySkullSetScheme with the given parameters.
// Returns an error if validation fails.
func NewGreySkullSetScheme(fixedSets, fixedReps, amrapSets, minAMRAPReps int) (*GreySkullSetScheme, error) {
	scheme := &GreySkullSetScheme{
		FixedSets:    fixedSets,
		FixedReps:    fixedReps,
		AMRAPSets:    amrapSets,
		MinAMRAPReps: minAMRAPReps,
	}
	if err := scheme.Validate(); err != nil {
		return nil, err
	}
	return scheme, nil
}

// Type returns the discriminator string for this scheme.
func (g *GreySkullSetScheme) Type() SetSchemeType {
	return TypeGreySkull
}

// GenerateSets generates concrete sets from a base weight.
// Each generated set has:
//   - SetNumber: 1 through (FixedSets + AMRAPSets) (1-indexed)
//   - Weight: baseWeight (unchanged)
//   - TargetReps: FixedReps for fixed sets, MinAMRAPReps for AMRAP sets
//   - IsWorkSet: true (all GreySkull sets are work sets)
//
// The first FixedSets sets are fixed sets, followed by AMRAPSets AMRAP sets.
// AMRAP sets use MinAMRAPReps as the target, representing the minimum to "succeed".
func (g *GreySkullSetScheme) GenerateSets(baseWeight float64, _ SetGenerationContext) ([]GeneratedSet, error) {
	if err := g.Validate(); err != nil {
		return nil, err
	}

	totalSets := g.FixedSets + g.AMRAPSets
	sets := make([]GeneratedSet, 0, totalSets)

	// Generate fixed sets
	for i := 0; i < g.FixedSets; i++ {
		sets = append(sets, GeneratedSet{
			SetNumber:  i + 1,
			Weight:     baseWeight,
			TargetReps: g.FixedReps,
			IsWorkSet:  true,
		})
	}

	// Generate AMRAP sets
	for i := 0; i < g.AMRAPSets; i++ {
		sets = append(sets, GeneratedSet{
			SetNumber:  g.FixedSets + i + 1,
			Weight:     baseWeight,
			TargetReps: g.MinAMRAPReps,
			IsWorkSet:  true,
		})
	}

	return sets, nil
}

// Validate validates the scheme's configuration parameters.
// Returns an error if:
//   - FixedSets is less than 0
//   - FixedReps is less than 1 (when FixedSets > 0)
//   - AMRAPSets is less than 1
//   - MinAMRAPReps is less than 1
func (g *GreySkullSetScheme) Validate() error {
	if g.FixedSets < 0 {
		return fmt.Errorf("%w: fixedSets must be >= 0, got %d", ErrInvalidParams, g.FixedSets)
	}
	if g.FixedSets > 0 && g.FixedReps < 1 {
		return fmt.Errorf("%w: fixedReps must be >= 1 when fixedSets > 0, got %d", ErrInvalidParams, g.FixedReps)
	}
	if g.AMRAPSets < 1 {
		return fmt.Errorf("%w: amrapSets must be >= 1, got %d", ErrInvalidParams, g.AMRAPSets)
	}
	if g.MinAMRAPReps < 1 {
		return fmt.Errorf("%w: minAmrapReps must be >= 1, got %d", ErrInvalidParams, g.MinAMRAPReps)
	}
	return nil
}

// MarshalJSON implements json.Marshaler for GreySkullSetScheme.
// Includes the type discriminator for polymorphic deserialization.
func (g *GreySkullSetScheme) MarshalJSON() ([]byte, error) {
	type Alias GreySkullSetScheme
	return json.Marshal(&struct {
		Type SetSchemeType `json:"type"`
		*Alias
	}{
		Type:  TypeGreySkull,
		Alias: (*Alias)(g),
	})
}

// UnmarshalGreySkullSetScheme deserializes a GreySkullSetScheme from JSON.
// This is used by the SchemeFactory.
func UnmarshalGreySkullSetScheme(data json.RawMessage) (SetScheme, error) {
	var scheme GreySkullSetScheme
	if err := json.Unmarshal(data, &scheme); err != nil {
		return nil, fmt.Errorf("failed to unmarshal GreySkullSetScheme: %w", err)
	}
	if err := scheme.Validate(); err != nil {
		return nil, err
	}
	return &scheme, nil
}

// RegisterGreySkullScheme registers the GreySkullSetScheme with the given factory.
func RegisterGreySkullScheme(factory *SchemeFactory) {
	factory.Register(TypeGreySkull, UnmarshalGreySkullSetScheme)
}

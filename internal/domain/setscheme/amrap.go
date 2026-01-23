// Package setscheme provides domain logic for set/rep scheme strategies.
package setscheme

import (
	"encoding/json"
	"fmt"
)

// AMRAPSetScheme generates sets where the lifter performs as many reps as possible.
// This is common in programs like 5/3/1 Wendler (e.g., "5+" means do at least 5, then AMRAP).
// The MinReps field represents the minimum reps to "succeed" - the actual reps performed
// are logged when the workout is completed.
type AMRAPSetScheme struct {
	// Sets is the number of AMRAP sets to generate (required, must be >= 1).
	// Usually 1 since AMRAP is typically a single top set.
	Sets int `json:"sets"`
	// MinReps is the minimum expected reps (required, must be >= 1).
	// This is used for display/logging and represents the "floor" for the AMRAP.
	// For example, Wendler's "5+" has MinReps=5.
	MinReps int `json:"minReps"`
}

// NewAMRAPSetScheme creates a new AMRAPSetScheme with the given sets and minReps.
// Returns an error if validation fails.
func NewAMRAPSetScheme(sets, minReps int) (*AMRAPSetScheme, error) {
	scheme := &AMRAPSetScheme{
		Sets:    sets,
		MinReps: minReps,
	}
	if err := scheme.Validate(); err != nil {
		return nil, err
	}
	return scheme, nil
}

// Type returns the discriminator string for this scheme.
func (a *AMRAPSetScheme) Type() SetSchemeType {
	return TypeAMRAP
}

// GenerateSets generates concrete sets from a base weight.
// Each generated set has:
//   - SetNumber: 1 through Sets (1-indexed)
//   - Weight: baseWeight (unchanged)
//   - TargetReps: MinReps value (the minimum to "succeed")
//   - IsWorkSet: true (all AMRAP sets are work sets)
//
// Note: The actual "AMRAP" behavior (logging more reps than target) happens
// at workout logging time. For generation purposes, AMRAP works like Fixed.
func (a *AMRAPSetScheme) GenerateSets(baseWeight float64, _ SetGenerationContext) ([]GeneratedSet, error) {
	if err := a.Validate(); err != nil {
		return nil, err
	}

	sets := make([]GeneratedSet, a.Sets)
	for i := 0; i < a.Sets; i++ {
		sets[i] = GeneratedSet{
			SetNumber:  i + 1,
			Weight:     baseWeight,
			TargetReps: a.MinReps,
			IsWorkSet:  true,
		}
	}
	return sets, nil
}

// Validate validates the scheme's configuration parameters.
// Returns an error if:
//   - Sets is less than 1
//   - MinReps is less than 1
func (a *AMRAPSetScheme) Validate() error {
	if a.Sets < 1 {
		return fmt.Errorf("%w: sets must be >= 1, got %d", ErrInvalidParams, a.Sets)
	}
	if a.MinReps < 1 {
		return fmt.Errorf("%w: minReps must be >= 1, got %d", ErrInvalidParams, a.MinReps)
	}
	return nil
}

// MarshalJSON implements json.Marshaler for AMRAPSetScheme.
// Includes the type discriminator for polymorphic deserialization.
func (a *AMRAPSetScheme) MarshalJSON() ([]byte, error) {
	type Alias AMRAPSetScheme
	return json.Marshal(&struct {
		Type SetSchemeType `json:"type"`
		*Alias
	}{
		Type:  TypeAMRAP,
		Alias: (*Alias)(a),
	})
}

// UnmarshalAMRAPSetScheme deserializes an AMRAPSetScheme from JSON.
// This is used by the SchemeFactory.
func UnmarshalAMRAPSetScheme(data json.RawMessage) (SetScheme, error) {
	var scheme AMRAPSetScheme
	if err := json.Unmarshal(data, &scheme); err != nil {
		return nil, fmt.Errorf("failed to unmarshal AMRAPSetScheme: %w", err)
	}
	if err := scheme.Validate(); err != nil {
		return nil, err
	}
	return &scheme, nil
}

// RegisterAMRAPScheme registers the AMRAPSetScheme with the given factory.
func RegisterAMRAPScheme(factory *SchemeFactory) {
	factory.Register(TypeAMRAP, UnmarshalAMRAPSetScheme)
}

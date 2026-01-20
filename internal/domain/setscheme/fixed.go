// Package setscheme provides domain logic for set/rep scheme strategies.
package setscheme

import (
	"encoding/json"
	"fmt"
)

// FixedSetScheme generates a specified number of sets with the same weight and reps.
// This is the most common scheme in powerlifting (e.g., 5x5, 3x8, 4x6).
// All generated sets are work sets.
type FixedSetScheme struct {
	// Sets is the number of sets to generate (required, must be >= 1).
	Sets int `json:"sets"`
	// Reps is the number of repetitions per set (required, must be >= 1).
	Reps int `json:"reps"`
}

// NewFixedSetScheme creates a new FixedSetScheme with the given sets and reps.
// Returns an error if validation fails.
func NewFixedSetScheme(sets, reps int) (*FixedSetScheme, error) {
	scheme := &FixedSetScheme{
		Sets: sets,
		Reps: reps,
	}
	if err := scheme.Validate(); err != nil {
		return nil, err
	}
	return scheme, nil
}

// Type returns the discriminator string for this scheme.
func (f *FixedSetScheme) Type() SetSchemeType {
	return TypeFixed
}

// GenerateSets generates concrete sets from a base weight.
// Each generated set has:
//   - SetNumber: 1 through Sets (1-indexed)
//   - Weight: baseWeight (unchanged)
//   - TargetReps: Reps value
//   - IsWorkSet: true (all Fixed sets are work sets)
func (f *FixedSetScheme) GenerateSets(baseWeight float64, _ SetGenerationContext) ([]GeneratedSet, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}

	sets := make([]GeneratedSet, f.Sets)
	for i := 0; i < f.Sets; i++ {
		sets[i] = GeneratedSet{
			SetNumber:  i + 1,
			Weight:     baseWeight,
			TargetReps: f.Reps,
			IsWorkSet:  true,
		}
	}
	return sets, nil
}

// Validate validates the scheme's configuration parameters.
// Returns an error if:
//   - Sets is less than 1
//   - Reps is less than 1
func (f *FixedSetScheme) Validate() error {
	if f.Sets < 1 {
		return fmt.Errorf("%w: sets must be >= 1, got %d", ErrInvalidParams, f.Sets)
	}
	if f.Reps < 1 {
		return fmt.Errorf("%w: reps must be >= 1, got %d", ErrInvalidParams, f.Reps)
	}
	return nil
}

// MarshalJSON implements json.Marshaler for FixedSetScheme.
// Includes the type discriminator for polymorphic deserialization.
func (f *FixedSetScheme) MarshalJSON() ([]byte, error) {
	type Alias FixedSetScheme
	return json.Marshal(&struct {
		Type SetSchemeType `json:"type"`
		*Alias
	}{
		Type:  TypeFixed,
		Alias: (*Alias)(f),
	})
}

// UnmarshalFixedSetScheme deserializes a FixedSetScheme from JSON.
// This is used by the SchemeFactory.
func UnmarshalFixedSetScheme(data json.RawMessage) (SetScheme, error) {
	var scheme FixedSetScheme
	if err := json.Unmarshal(data, &scheme); err != nil {
		return nil, fmt.Errorf("failed to unmarshal FixedSetScheme: %w", err)
	}
	if err := scheme.Validate(); err != nil {
		return nil, err
	}
	return &scheme, nil
}

// RegisterFixedScheme registers the FixedSetScheme with the given factory.
func RegisterFixedScheme(factory *SchemeFactory) {
	factory.Register(TypeFixed, UnmarshalFixedSetScheme)
}

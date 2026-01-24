// Package setscheme provides domain logic for set/rep scheme strategies.
package setscheme

import (
	"encoding/json"
	"fmt"
)

// RepRangeSetScheme generates sets with a target rep range (e.g., 3x8-12).
// This scheme is common in hypertrophy and double progression programs where
// lifters work within a rep range and progress weight once they hit the max.
// GenerateSets uses MinReps as the target for planning purposes.
type RepRangeSetScheme struct {
	// Sets is the number of sets to generate (required, must be >= 1).
	Sets int `json:"sets"`
	// MinReps is the minimum reps for the range (required, must be >= 1).
	MinReps int `json:"minReps"`
	// MaxReps is the maximum reps for the range (required, must be >= MinReps).
	MaxReps int `json:"maxReps"`
}

// NewRepRangeSetScheme creates a new RepRangeSetScheme with the given parameters.
// Returns an error if validation fails.
func NewRepRangeSetScheme(sets, minReps, maxReps int) (*RepRangeSetScheme, error) {
	scheme := &RepRangeSetScheme{
		Sets:    sets,
		MinReps: minReps,
		MaxReps: maxReps,
	}
	if err := scheme.Validate(); err != nil {
		return nil, err
	}
	return scheme, nil
}

// Type returns the discriminator string for this scheme.
func (r *RepRangeSetScheme) Type() SetSchemeType {
	return TypeRepRange
}

// GenerateSets generates concrete sets from a base weight.
// Each generated set has:
//   - SetNumber: 1 through Sets (1-indexed)
//   - Weight: baseWeight (unchanged)
//   - TargetReps: MinReps value (the minimum of the range)
//   - IsWorkSet: true (all RepRange sets are work sets)
//
// Note: The actual rep range behavior (tracking progress toward MaxReps) happens
// at workout logging time. For generation purposes, MinReps is used as the target.
func (r *RepRangeSetScheme) GenerateSets(baseWeight float64, _ SetGenerationContext) ([]GeneratedSet, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}

	sets := make([]GeneratedSet, r.Sets)
	for i := 0; i < r.Sets; i++ {
		sets[i] = GeneratedSet{
			SetNumber:  i + 1,
			Weight:     baseWeight,
			TargetReps: r.MinReps,
			IsWorkSet:  true,
		}
	}
	return sets, nil
}

// Validate validates the scheme's configuration parameters.
// Returns an error if:
//   - Sets is less than 1
//   - MinReps is less than 1
//   - MaxReps is less than MinReps
func (r *RepRangeSetScheme) Validate() error {
	if r.Sets < 1 {
		return fmt.Errorf("%w: sets must be >= 1, got %d", ErrInvalidParams, r.Sets)
	}
	if r.MinReps < 1 {
		return fmt.Errorf("%w: minReps must be >= 1, got %d", ErrInvalidParams, r.MinReps)
	}
	if r.MaxReps < r.MinReps {
		return fmt.Errorf("%w: maxReps must be >= minReps, got maxReps=%d minReps=%d", ErrInvalidParams, r.MaxReps, r.MinReps)
	}
	return nil
}

// MarshalJSON implements json.Marshaler for RepRangeSetScheme.
// Includes the type discriminator for polymorphic deserialization.
func (r *RepRangeSetScheme) MarshalJSON() ([]byte, error) {
	type Alias RepRangeSetScheme
	return json.Marshal(&struct {
		Type SetSchemeType `json:"type"`
		*Alias
	}{
		Type:  TypeRepRange,
		Alias: (*Alias)(r),
	})
}

// UnmarshalRepRangeSetScheme deserializes a RepRangeSetScheme from JSON.
// This is used by the SchemeFactory.
func UnmarshalRepRangeSetScheme(data json.RawMessage) (SetScheme, error) {
	var scheme RepRangeSetScheme
	if err := json.Unmarshal(data, &scheme); err != nil {
		return nil, fmt.Errorf("failed to unmarshal RepRangeSetScheme: %w", err)
	}
	if err := scheme.Validate(); err != nil {
		return nil, err
	}
	return &scheme, nil
}

// RegisterRepRangeScheme registers the RepRangeSetScheme with the given factory.
func RegisterRepRangeScheme(factory *SchemeFactory) {
	factory.Register(TypeRepRange, UnmarshalRepRangeSetScheme)
}

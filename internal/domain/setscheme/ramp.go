// Package setscheme provides domain logic for set/rep scheme strategies.
package setscheme

import (
	"encoding/json"
	"fmt"
)

// DefaultWorkSetThreshold is the default percentage threshold for work set classification.
// Sets at or above this percentage are considered work sets.
const DefaultWorkSetThreshold = 80.0

// RampStep represents a single step in a ramp progression.
// Each step defines the percentage of baseWeight and reps for one set.
type RampStep struct {
	// Percentage is the percentage of baseWeight for this set (required, must be > 0).
	Percentage float64 `json:"percentage"`
	// Reps is the number of repetitions for this set (required, must be >= 1).
	Reps int `json:"reps"`
}

// RampSetScheme generates sets with progressive percentages across a series of sets.
// Used for warmup progressions and Bill Starr style ramping sets.
// Each step generates one set, with weight calculated as baseWeight * (step.Percentage / 100).
type RampSetScheme struct {
	// Steps is the array of percentage/rep pairs defining the ramp progression (required, at least one step).
	Steps []RampStep `json:"steps"`
	// WorkSetThreshold is the percentage above which sets are classified as work sets (default 80).
	// Sets with percentage >= WorkSetThreshold are work sets; below are warmup sets.
	WorkSetThreshold float64 `json:"workSetThreshold,omitempty"`
}

// NewRampSetScheme creates a new RampSetScheme with the given steps.
// Uses the default WorkSetThreshold of 80%.
// Returns an error if validation fails.
func NewRampSetScheme(steps []RampStep) (*RampSetScheme, error) {
	return NewRampSetSchemeWithThreshold(steps, DefaultWorkSetThreshold)
}

// NewRampSetSchemeWithThreshold creates a new RampSetScheme with the given steps and work set threshold.
// Returns an error if validation fails.
func NewRampSetSchemeWithThreshold(steps []RampStep, workSetThreshold float64) (*RampSetScheme, error) {
	scheme := &RampSetScheme{
		Steps:            steps,
		WorkSetThreshold: workSetThreshold,
	}
	if err := scheme.Validate(); err != nil {
		return nil, err
	}
	return scheme, nil
}

// Type returns the discriminator string for this scheme.
func (r *RampSetScheme) Type() SetSchemeType {
	return TypeRamp
}

// GenerateSets generates concrete sets from a base weight.
// Each generated set has:
//   - SetNumber: 1 through len(Steps) (1-indexed)
//   - Weight: baseWeight * (step.Percentage / 100)
//   - TargetReps: step.Reps value
//   - IsWorkSet: true if step.Percentage >= WorkSetThreshold
func (r *RampSetScheme) GenerateSets(baseWeight float64, ctx SetGenerationContext) ([]GeneratedSet, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}

	// Use the scheme's WorkSetThreshold, falling back to default if not set
	threshold := r.WorkSetThreshold
	if threshold <= 0 {
		threshold = DefaultWorkSetThreshold
	}

	sets := make([]GeneratedSet, len(r.Steps))
	for i, step := range r.Steps {
		weight := baseWeight * (step.Percentage / 100)
		sets[i] = GeneratedSet{
			SetNumber:  i + 1,
			Weight:     weight,
			TargetReps: step.Reps,
			IsWorkSet:  step.Percentage >= threshold,
		}
	}
	return sets, nil
}

// Validate validates the scheme's configuration parameters.
// Returns an error if:
//   - Steps is empty (at least one step required)
//   - Any step percentage is <= 0
//   - Any step reps is < 1
//   - WorkSetThreshold is <= 0 or > 100
func (r *RampSetScheme) Validate() error {
	// At least one step required
	if len(r.Steps) < 1 {
		return fmt.Errorf("%w: at least one step required", ErrInvalidParams)
	}

	// Validate WorkSetThreshold if explicitly set (non-zero)
	// A zero value means "use default", which is valid
	if r.WorkSetThreshold < 0 {
		return fmt.Errorf("%w: workSetThreshold must be > 0, got %.2f", ErrInvalidParams, r.WorkSetThreshold)
	}
	if r.WorkSetThreshold > 100 {
		return fmt.Errorf("%w: workSetThreshold must be <= 100, got %.2f", ErrInvalidParams, r.WorkSetThreshold)
	}

	// Validate each step
	for i, step := range r.Steps {
		if step.Percentage <= 0 {
			return fmt.Errorf("%w: step %d percentage must be > 0, got %.2f", ErrInvalidParams, i+1, step.Percentage)
		}
		if step.Reps < 1 {
			return fmt.Errorf("%w: step %d reps must be >= 1, got %d", ErrInvalidParams, i+1, step.Reps)
		}
	}

	return nil
}

// MarshalJSON implements json.Marshaler for RampSetScheme.
// Includes the type discriminator for polymorphic deserialization.
func (r *RampSetScheme) MarshalJSON() ([]byte, error) {
	type Alias RampSetScheme
	return json.Marshal(&struct {
		Type SetSchemeType `json:"type"`
		*Alias
	}{
		Type:  TypeRamp,
		Alias: (*Alias)(r),
	})
}

// UnmarshalRampSetScheme deserializes a RampSetScheme from JSON.
// This is used by the SchemeFactory.
func UnmarshalRampSetScheme(data json.RawMessage) (SetScheme, error) {
	var scheme RampSetScheme
	if err := json.Unmarshal(data, &scheme); err != nil {
		return nil, fmt.Errorf("failed to unmarshal RampSetScheme: %w", err)
	}
	// Apply default WorkSetThreshold if not specified in JSON (will be zero)
	// Validation will pass because zero is allowed (means use default)
	if err := scheme.Validate(); err != nil {
		return nil, err
	}
	return &scheme, nil
}

// RegisterRampScheme registers the RampSetScheme with the given factory.
func RegisterRampScheme(factory *SchemeFactory) {
	factory.Register(TypeRamp, UnmarshalRampSetScheme)
}

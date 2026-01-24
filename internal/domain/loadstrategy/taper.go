// Package loadstrategy provides domain logic for load calculation strategies.
package loadstrategy

import (
	"context"
	"encoding/json"
	"fmt"
)

// TypeTaper is the strategy type for taper-based load modification.
const TypeTaper LoadStrategyType = "TAPER"

// TaperCurve defines the volume reduction curve as meet approaches.
// Each entry maps a maximum days-out threshold to a volume multiplier.
// The curve is evaluated in order; the first threshold that daysOut is less than
// determines the multiplier.
type TaperCurve struct {
	// ThresholdDays is the upper bound (exclusive) of days for this tier.
	// E.g., ThresholdDays=7 means "if daysOut < 7, use this multiplier".
	ThresholdDays int `json:"thresholdDays"`

	// Multiplier is the volume multiplier (0.0-1.0) applied when in this tier.
	// E.g., 0.5 means 50% of normal volume (50% reduction).
	Multiplier float64 `json:"multiplier"`
}

// DefaultTaperCurve returns the standard competition taper curve.
// Volume is progressively reduced as the meet approaches:
//   - Week 5 of Comp (days 0-6): 50% volume (final week)
//   - Week 4 (days 7-13): 60% volume
//   - Week 3 (days 14-20): 70% volume
//   - Week 2 (days 21-27): 80% volume
//   - Week 1 (days 28-34): 90% volume
//   - Beyond 35 days: 100% volume (no taper)
func DefaultTaperCurve() []TaperCurve {
	return []TaperCurve{
		{ThresholdDays: 7, Multiplier: 0.5},   // Final week
		{ThresholdDays: 14, Multiplier: 0.6},  // Week 4
		{ThresholdDays: 21, Multiplier: 0.7},  // Week 3
		{ThresholdDays: 28, Multiplier: 0.8},  // Week 2
		{ThresholdDays: 35, Multiplier: 0.9},  // Week 1 of Comp
	}
}

// GetTaperMultiplier calculates the volume multiplier based on days out from meet.
// Uses the default taper curve. Returns 1.0 (no taper) if daysOut >= 35.
func GetTaperMultiplier(daysOut int) float64 {
	return GetTaperMultiplierWithCurve(daysOut, DefaultTaperCurve())
}

// GetTaperMultiplierWithCurve calculates the volume multiplier using a custom curve.
// The curve is evaluated in order of threshold; returns 1.0 if no threshold matches.
func GetTaperMultiplierWithCurve(daysOut int, curve []TaperCurve) float64 {
	// If daysOut is negative (past meet date), still apply final week taper
	if daysOut < 0 {
		daysOut = 0
	}

	for _, tier := range curve {
		if daysOut < tier.ThresholdDays {
			return tier.Multiplier
		}
	}

	// No matching tier means no taper (100% volume)
	return 1.0
}

// TaperLoadStrategy applies a taper multiplier to a base load strategy.
// This is a decorator pattern: it wraps another LoadStrategy and modifies
// the calculated load based on days remaining until meet.
//
// The taper affects VOLUME (number of sets/reps or total tonnage), not intensity.
// In practice, this means reducing the calculated load proportionally to
// reduce training stress as the competition approaches.
//
// Note: Some programs may prefer to reduce sets rather than load. This strategy
// provides a load-based taper that can be used in conjunction with set reduction.
type TaperLoadStrategy struct {
	// BaseStrategy is the underlying load calculation strategy.
	// The taper multiplier is applied to the result of this strategy's CalculateLoad.
	BaseStrategy LoadStrategy `json:"baseStrategy"`

	// TaperCurve defines the volume reduction schedule.
	// Optional: if nil or empty, DefaultTaperCurve() is used.
	TaperCurve []TaperCurve `json:"taperCurve,omitempty"`

	// MaintainIntensity, when true, does not reduce the calculated load.
	// Instead, the taper multiplier is included in the result context for
	// external handling (e.g., reducing sets/reps instead of load).
	// Default is false (load is reduced by the taper multiplier).
	MaintainIntensity bool `json:"maintainIntensity,omitempty"`

	// rawBaseStrategy stores the raw JSON for delayed deserialization.
	rawBaseStrategy json.RawMessage `json:"-"`
}

// Type returns the strategy type discriminator.
func (s *TaperLoadStrategy) Type() LoadStrategyType {
	return TypeTaper
}

// CalculateLoad calculates the tapered load based on the base strategy and days out.
// The daysOut value is extracted from params.Context["daysOut"].
// If daysOut is not provided, no taper is applied (returns base load).
func (s *TaperLoadStrategy) CalculateLoad(ctx context.Context, params LoadCalculationParams) (float64, error) {
	if err := s.Validate(); err != nil {
		return 0, err
	}

	// Calculate base load from the wrapped strategy
	baseLoad, err := s.BaseStrategy.CalculateLoad(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("taper: base strategy calculation failed: %w", err)
	}

	// Extract daysOut from context
	daysOut, ok := s.extractDaysOut(params.Context)
	if !ok {
		// No daysOut provided; return base load without taper
		return baseLoad, nil
	}

	// Get the taper multiplier
	curve := s.TaperCurve
	if len(curve) == 0 {
		curve = DefaultTaperCurve()
	}
	multiplier := GetTaperMultiplierWithCurve(daysOut, curve)

	// If maintaining intensity, don't modify the load
	if s.MaintainIntensity {
		// The multiplier could be stored back in context for external use,
		// but for now we just return the base load unchanged.
		return baseLoad, nil
	}

	// Apply taper multiplier to the load
	return baseLoad * multiplier, nil
}

// extractDaysOut extracts the days-out value from the context map.
// Returns the value and true if found, or 0 and false if not present or invalid.
func (s *TaperLoadStrategy) extractDaysOut(ctxMap map[string]interface{}) (int, bool) {
	if ctxMap == nil {
		return 0, false
	}

	daysOutRaw, exists := ctxMap["daysOut"]
	if !exists {
		return 0, false
	}

	switch v := daysOutRaw.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	default:
		return 0, false
	}
}

// Validate validates the strategy's configuration parameters.
func (s *TaperLoadStrategy) Validate() error {
	if s.BaseStrategy == nil {
		return fmt.Errorf("%w: base strategy is required", ErrInvalidParams)
	}

	// Validate the base strategy
	if err := s.BaseStrategy.Validate(); err != nil {
		return fmt.Errorf("taper: invalid base strategy: %w", err)
	}

	// Validate taper curve if provided
	for i, tier := range s.TaperCurve {
		if tier.ThresholdDays <= 0 {
			return fmt.Errorf("%w: taper curve entry %d has invalid threshold days", ErrInvalidParams, i)
		}
		if tier.Multiplier < 0 || tier.Multiplier > 1 {
			return fmt.Errorf("%w: taper curve entry %d has invalid multiplier (must be 0-1)", ErrInvalidParams, i)
		}
	}

	// Validate curve is in ascending order of thresholds
	for i := 1; i < len(s.TaperCurve); i++ {
		if s.TaperCurve[i].ThresholdDays <= s.TaperCurve[i-1].ThresholdDays {
			return fmt.Errorf("%w: taper curve must have ascending threshold days", ErrInvalidParams)
		}
	}

	return nil
}

// SetMaxLookup sets the max lookup on the base strategy if it supports it.
func (s *TaperLoadStrategy) SetMaxLookup(maxLookup MaxLookup) {
	if setter, ok := s.BaseStrategy.(interface{ SetMaxLookup(MaxLookup) }); ok {
		setter.SetMaxLookup(maxLookup)
	}
}

// MarshalJSON implements json.Marshaler.
// Includes the type discriminator in the JSON output.
func (s *TaperLoadStrategy) MarshalJSON() ([]byte, error) {
	type Alias TaperLoadStrategy
	return json.Marshal(&struct {
		Type LoadStrategyType `json:"type"`
		*Alias
	}{
		Type:  TypeTaper,
		Alias: (*Alias)(s),
	})
}

// UnmarshalTaper deserializes a TaperLoadStrategy from JSON.
// This is a factory function that can be registered with StrategyFactory.
// Note: The base strategy is deserialized using the provided factory.
func UnmarshalTaper(factory *StrategyFactory) func(json.RawMessage) (LoadStrategy, error) {
	return func(data json.RawMessage) (LoadStrategy, error) {
		// First, unmarshal the taper-specific fields
		var envelope struct {
			BaseStrategy json.RawMessage `json:"baseStrategy"`
			TaperCurve   []TaperCurve    `json:"taperCurve,omitempty"`
			MaintainIntensity bool       `json:"maintainIntensity,omitempty"`
		}
		if err := json.Unmarshal(data, &envelope); err != nil {
			return nil, fmt.Errorf("failed to unmarshal Taper strategy: %w", err)
		}

		// Deserialize the base strategy using the factory
		baseStrategy, err := factory.CreateFromJSON(envelope.BaseStrategy)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal Taper base strategy: %w", err)
		}

		s := &TaperLoadStrategy{
			BaseStrategy:      baseStrategy,
			TaperCurve:        envelope.TaperCurve,
			MaintainIntensity: envelope.MaintainIntensity,
		}

		// Validate the deserialized strategy
		if err := s.Validate(); err != nil {
			return nil, fmt.Errorf("invalid Taper strategy: %w", err)
		}

		return s, nil
	}
}

// RegisterTaper registers the Taper strategy with a factory.
// This is a convenience function for setting up the factory.
func RegisterTaper(factory *StrategyFactory) {
	factory.Register(TypeTaper, UnmarshalTaper(factory))
}

// NewTaperLoadStrategy creates a new TaperLoadStrategy wrapping the given base strategy.
func NewTaperLoadStrategy(baseStrategy LoadStrategy, curve []TaperCurve, maintainIntensity bool) *TaperLoadStrategy {
	return &TaperLoadStrategy{
		BaseStrategy:      baseStrategy,
		TaperCurve:        curve,
		MaintainIntensity: maintainIntensity,
	}
}

// Ensure TaperLoadStrategy implements LoadStrategy.
var _ LoadStrategy = (*TaperLoadStrategy)(nil)

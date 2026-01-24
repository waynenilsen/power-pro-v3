// Package progression provides domain logic for progression strategies.
// This file implements the StageProgression strategy that changes set/rep schemes on failure,
// commonly used in programs like GZCLP where lifters cycle through different rep schemes
// (e.g., 5x3+ -> 6x2+ -> 10x1+) at the same weight before resetting.
package progression

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/domain/setscheme"
)

// TypeStage is the progression type constant for stage progressions.
const TypeStage ProgressionType = "STAGE_PROGRESSION"

// Stage represents a single stage in a stage-based progression scheme.
// Each stage defines a set/rep configuration that the lifter moves to on failure.
type Stage struct {
	// Name is a human-readable identifier for this stage (e.g., "5x3+", "6x2+", "10x1+").
	Name string `json:"name"`
	// Sets is the number of sets for this stage.
	Sets int `json:"sets"`
	// Reps is the target reps per set.
	Reps int `json:"reps"`
	// IsAMRAP indicates if the last set is an AMRAP set.
	IsAMRAP bool `json:"isAmrap"`
	// MinVolume is the minimum total reps required to pass this stage.
	// For example, 5x3 would have MinVolume=15 (must get at least 15 total reps).
	MinVolume int `json:"minVolume"`
}

// Validate validates a Stage configuration.
func (s Stage) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("%w: stage name is required", ErrInvalidParams)
	}
	if s.Sets < 1 {
		return fmt.Errorf("%w: stage sets must be at least 1", ErrInvalidParams)
	}
	if s.Reps < 1 {
		return fmt.Errorf("%w: stage reps must be at least 1", ErrInvalidParams)
	}
	if s.MinVolume < 1 {
		return fmt.Errorf("%w: stage minVolume must be at least 1", ErrInvalidParams)
	}
	return nil
}

// ToSetScheme converts this Stage to a SetScheme for prescription generation.
// Returns either an AMRAPSetScheme or FixedSetScheme depending on IsAMRAP.
func (s Stage) ToSetScheme() setscheme.SetScheme {
	if s.IsAMRAP {
		scheme, _ := setscheme.NewAMRAPSetScheme(s.Sets, s.Reps)
		return scheme
	}
	scheme, _ := setscheme.NewFixedSetScheme(s.Sets, s.Reps)
	return scheme
}

// StageProgression implements the Progression interface for stage-based rep scheme changes.
// This progression type cycles through different set/rep configurations on failure,
// keeping the weight constant until all stages are exhausted.
//
// Examples:
//   - GZCLP T1 Default: 5x3+ -> 6x2+ -> 10x1+ (on failure, advance stage; after 10x1+ fails, reset)
//   - GZCLP T2 Default: 3x10 -> 3x8 -> 3x6 (on failure, advance stage; after 3x6 fails, reset)
//
// When the lifter fails the current stage (total reps < MinVolume), they advance to the next stage.
// The same weight is maintained across all stages. After exhausting all stages, the progression
// can optionally reset to stage 0 with an optional deload.
type StageProgression struct {
	// ID is the unique identifier for this progression.
	ID string `json:"id"`
	// Name is the human-readable name for this progression.
	Name string `json:"name"`
	// Stages is the ordered list of stages to cycle through.
	// The lifter starts at index 0 and advances on failure.
	Stages []Stage `json:"stages"`
	// CurrentStage is the current stage index (0-based).
	// This value is updated as the lifter progresses through stages.
	CurrentStage int `json:"currentStage"`
	// ResetOnExhaustion indicates whether to reset to stage 0 after the last stage fails.
	// If false, the progression returns a special result requiring manual intervention.
	ResetOnExhaustion bool `json:"resetOnExhaustion"`
	// DeloadOnReset indicates whether to apply a deload when resetting to stage 0.
	// Only applicable when ResetOnExhaustion is true.
	DeloadOnReset bool `json:"deloadOnReset"`
	// DeloadPercent is the percentage to deload when resetting (e.g., 0.15 for 15%).
	// Only used when both ResetOnExhaustion and DeloadOnReset are true.
	DeloadPercent float64 `json:"deloadPercent,omitempty"`
	// MaxTypeValue specifies which max to update (ONE_RM or TRAINING_MAX).
	MaxTypeValue MaxType `json:"maxType"`
}

// NewStageProgression creates a new StageProgression with the given parameters.
func NewStageProgression(id, name string, stages []Stage, resetOnExhaustion, deloadOnReset bool, deloadPercent float64, maxType MaxType) (*StageProgression, error) {
	s := &StageProgression{
		ID:                id,
		Name:              name,
		Stages:            stages,
		CurrentStage:      0,
		ResetOnExhaustion: resetOnExhaustion,
		DeloadOnReset:     deloadOnReset,
		DeloadPercent:     deloadPercent,
		MaxTypeValue:      maxType,
	}
	if err := s.Validate(); err != nil {
		return nil, err
	}
	return s, nil
}

// Type returns the discriminator string for this progression.
// Implements Progression interface.
func (s *StageProgression) Type() ProgressionType {
	return TypeStage
}

// TriggerType returns the trigger type this progression responds to.
// StageProgression only responds to ON_FAILURE triggers.
// Implements Progression interface.
func (s *StageProgression) TriggerType() TriggerType {
	return TriggerOnFailure
}

// Validate validates the progression's configuration parameters.
// Implements Progression interface.
func (s *StageProgression) Validate() error {
	if s.ID == "" {
		return fmt.Errorf("%w: id is required", ErrInvalidParams)
	}
	if s.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidParams)
	}
	if len(s.Stages) < 2 {
		return fmt.Errorf("%w: at least 2 stages are required", ErrInvalidParams)
	}
	if s.CurrentStage < 0 || s.CurrentStage >= len(s.Stages) {
		return fmt.Errorf("%w: currentStage must be between 0 and %d", ErrInvalidParams, len(s.Stages)-1)
	}
	if err := ValidateMaxType(s.MaxTypeValue); err != nil {
		return err
	}

	// Validate each stage
	for i, stage := range s.Stages {
		if err := stage.Validate(); err != nil {
			return fmt.Errorf("stage %d: %w", i, err)
		}
	}

	// Validate deload configuration
	if s.DeloadOnReset {
		if !s.ResetOnExhaustion {
			return fmt.Errorf("%w: deloadOnReset requires resetOnExhaustion to be true", ErrInvalidParams)
		}
		if s.DeloadPercent <= 0 || s.DeloadPercent > 1 {
			return fmt.Errorf("%w: deloadPercent must be between 0 (exclusive) and 1 (inclusive)", ErrInvalidParams)
		}
	}

	return nil
}

// Apply evaluates and applies the progression given the context.
// Implements Progression interface.
//
// StageProgression implements the stage-based progression pattern used in programs like GZCLP.
// The key characteristics are:
//
//  1. Trigger type must be ON_FAILURE - this only fires when a set fails
//  2. On failure, advance to the next stage (change set/rep scheme, keep weight)
//  3. If at the last stage and failure occurs:
//     - If ResetOnExhaustion: reset to stage 0, optionally deload
//     - Otherwise: return Applied=false with reason for manual intervention
//  4. The result indicates the new stage and any weight change
//
// Note: The actual SetScheme change is communicated via NewStage and NewSetScheme fields
// in the StageProgressionResult. The caller is responsible for updating the prescription.
func (s *StageProgression) Apply(ctx context.Context, params ProgressionContext) (ProgressionResult, error) {
	if err := params.Validate(); err != nil {
		return ProgressionResult{}, fmt.Errorf("invalid progression context: %w", err)
	}

	now := time.Now()

	// Trigger type must be ON_FAILURE
	if params.TriggerEvent.Type != TriggerOnFailure {
		return ProgressionResult{
			Applied:       false,
			PreviousValue: params.CurrentValue,
			NewValue:      params.CurrentValue,
			Delta:         0,
			LiftID:        params.LiftID,
			MaxType:       params.MaxType,
			AppliedAt:     now,
			Reason:        fmt.Sprintf("trigger type mismatch: expected %s, got %s", TriggerOnFailure, params.TriggerEvent.Type),
		}, nil
	}

	// Max type must match
	if params.MaxType != s.MaxTypeValue {
		return ProgressionResult{
			Applied:       false,
			PreviousValue: params.CurrentValue,
			NewValue:      params.CurrentValue,
			Delta:         0,
			LiftID:        params.LiftID,
			MaxType:       params.MaxType,
			AppliedAt:     now,
			Reason:        fmt.Sprintf("max type mismatch: expected %s, got %s", s.MaxTypeValue, params.MaxType),
		}, nil
	}

	// Check if we're at the last stage
	isLastStage := s.CurrentStage >= len(s.Stages)-1

	if isLastStage {
		// At the last stage - handle exhaustion
		if !s.ResetOnExhaustion {
			// Manual intervention required
			return ProgressionResult{
				Applied:       false,
				PreviousValue: params.CurrentValue,
				NewValue:      params.CurrentValue,
				Delta:         0,
				LiftID:        params.LiftID,
				MaxType:       params.MaxType,
				AppliedAt:     now,
				Reason:        "all stages exhausted; manual intervention required",
			}, nil
		}

		// Reset to stage 0
		newStage := 0
		var delta float64
		newValue := params.CurrentValue

		if s.DeloadOnReset {
			// Apply deload
			deloadAmount := params.CurrentValue * s.DeloadPercent
			newValue = params.CurrentValue - deloadAmount
			delta = -deloadAmount

			// Ensure we don't go below zero
			if newValue < 0 {
				newValue = 0
				delta = -params.CurrentValue
			}
		}

		// Update current stage for next evaluation
		s.CurrentStage = newStage

		return ProgressionResult{
			Applied:       true,
			PreviousValue: params.CurrentValue,
			NewValue:      newValue,
			Delta:         delta,
			LiftID:        params.LiftID,
			MaxType:       params.MaxType,
			AppliedAt:     now,
		}, nil
	}

	// Not at last stage - advance to next stage
	s.CurrentStage++

	return ProgressionResult{
		Applied:       true,
		PreviousValue: params.CurrentValue,
		NewValue:      params.CurrentValue, // Weight stays the same
		Delta:         0,
		LiftID:        params.LiftID,
		MaxType:       params.MaxType,
		AppliedAt:     now,
	}, nil
}

// GetCurrentStage returns the current stage configuration.
func (s *StageProgression) GetCurrentStage() Stage {
	if s.CurrentStage < 0 || s.CurrentStage >= len(s.Stages) {
		// Return first stage as fallback
		return s.Stages[0]
	}
	return s.Stages[s.CurrentStage]
}

// GetCurrentSetScheme returns the SetScheme for the current stage.
// This is used when generating prescriptions.
func (s *StageProgression) GetCurrentSetScheme() setscheme.SetScheme {
	return s.GetCurrentStage().ToSetScheme()
}

// SetCurrentStage sets the current stage index.
// Returns an error if the index is out of bounds.
func (s *StageProgression) SetCurrentStage(stage int) error {
	if stage < 0 || stage >= len(s.Stages) {
		return fmt.Errorf("%w: stage index %d out of bounds [0, %d)", ErrInvalidParams, stage, len(s.Stages))
	}
	s.CurrentStage = stage
	return nil
}

// StageCount returns the total number of stages.
func (s *StageProgression) StageCount() int {
	return len(s.Stages)
}

// IsAtLastStage returns true if the progression is at the last stage.
func (s *StageProgression) IsAtLastStage() bool {
	return s.CurrentStage >= len(s.Stages)-1
}

// ShouldResetFailureCounter returns whether the failure counter should be reset
// after this progression is applied. For StageProgression, we reset on stage change
// since the set/rep scheme is changing.
func (s *StageProgression) ShouldResetFailureCounter() bool {
	return true
}

// MarshalJSON implements json.Marshaler for StageProgression.
// Ensures the type discriminator is always included in serialized output.
func (s *StageProgression) MarshalJSON() ([]byte, error) {
	type Alias StageProgression
	return json.Marshal(&struct {
		Type ProgressionType `json:"type"`
		*Alias
	}{
		Type:  TypeStage,
		Alias: (*Alias)(s),
	})
}

// UnmarshalStageProgression deserializes a StageProgression from JSON.
// This is used by the ProgressionFactory for type-safe deserialization.
func UnmarshalStageProgression(data json.RawMessage) (Progression, error) {
	var s StageProgression
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stage progression: %w", err)
	}
	if err := s.Validate(); err != nil {
		return nil, fmt.Errorf("invalid stage progression: %w", err)
	}
	return &s, nil
}

// RegisterStageProgression registers the StageProgression type with a factory.
// This should be called during application initialization.
func RegisterStageProgression(factory *ProgressionFactory) {
	factory.Register(TypeStage, UnmarshalStageProgression)
}

// GZCLP Stage Preset Functions

// NewGZCLPT1DefaultProgression creates a StageProgression for GZCLP T1 default scheme.
// Stages: 5x3+ -> 6x2+ -> 10x1+ with 15% deload on reset.
func NewGZCLPT1DefaultProgression(id, name string) (*StageProgression, error) {
	return NewStageProgression(id, name, []Stage{
		{Name: "5x3+", Sets: 5, Reps: 3, IsAMRAP: true, MinVolume: 15},
		{Name: "6x2+", Sets: 6, Reps: 2, IsAMRAP: true, MinVolume: 12},
		{Name: "10x1+", Sets: 10, Reps: 1, IsAMRAP: true, MinVolume: 10},
	}, true, true, 0.15, TrainingMax)
}

// NewGZCLPT1ModifiedProgression creates a StageProgression for GZCLP T1 modified scheme.
// Stages: 3x5+ -> 4x3+ -> 5x2+ with 15% deload on reset.
func NewGZCLPT1ModifiedProgression(id, name string) (*StageProgression, error) {
	return NewStageProgression(id, name, []Stage{
		{Name: "3x5+", Sets: 3, Reps: 5, IsAMRAP: true, MinVolume: 15},
		{Name: "4x3+", Sets: 4, Reps: 3, IsAMRAP: true, MinVolume: 12},
		{Name: "5x2+", Sets: 5, Reps: 2, IsAMRAP: true, MinVolume: 10},
	}, true, true, 0.15, TrainingMax)
}

// NewGZCLPT2DefaultProgression creates a StageProgression for GZCLP T2 default scheme.
// Stages: 3x10 -> 3x8 -> 3x6 with no deload on reset.
func NewGZCLPT2DefaultProgression(id, name string) (*StageProgression, error) {
	return NewStageProgression(id, name, []Stage{
		{Name: "3x10", Sets: 3, Reps: 10, IsAMRAP: false, MinVolume: 30},
		{Name: "3x8", Sets: 3, Reps: 8, IsAMRAP: false, MinVolume: 24},
		{Name: "3x6", Sets: 3, Reps: 6, IsAMRAP: false, MinVolume: 18},
	}, true, false, 0, TrainingMax)
}

// Package workout provides domain logic for workout generation.
// This package contains pure business logic with no database dependencies,
// making it testable in isolation.
package workout

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/domain/loadstrategy"
	"github.com/waynenilsen/power-pro-v3/internal/domain/prescription"
	"github.com/waynenilsen/power-pro-v3/internal/domain/setscheme"
)

// Validation errors
var (
	ErrUserNotEnrolled     = errors.New("user is not enrolled in any program")
	ErrDayNotFound         = errors.New("day not found for the specified position")
	ErrInvalidWeekNumber   = errors.New("week number must be >= 1")
	ErrInvalidDaySlug      = errors.New("day slug is required")
	ErrNoPrescriptions     = errors.New("day has no prescriptions")
	ErrProgramNotFound     = errors.New("program not found")
	ErrCycleNotFound       = errors.New("cycle not found")
	ErrWeekNotFound        = errors.New("week not found in cycle")
)

// LiftInfo contains minimal lift information for workout response.
type LiftInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// SetInfo represents a resolved set in a workout.
type SetInfo struct {
	SetNumber  int     `json:"setNumber"`
	Weight     float64 `json:"weight"`
	TargetReps int     `json:"targetReps"`
	IsWorkSet  bool    `json:"isWorkSet"`
}

// ExerciseInfo represents a resolved exercise in a workout.
type ExerciseInfo struct {
	PrescriptionID string    `json:"prescriptionId"`
	Lift           LiftInfo  `json:"lift"`
	Sets           []SetInfo `json:"sets"`
	Notes          string    `json:"notes,omitempty"`
	RestSeconds    *int      `json:"restSeconds,omitempty"`
}

// Workout represents a fully resolved workout for a user.
type Workout struct {
	UserID         string         `json:"userId"`
	ProgramID      string         `json:"programId"`
	CycleIteration int            `json:"cycleIteration"`
	WeekNumber     int            `json:"weekNumber"`
	DaySlug        string         `json:"daySlug"`
	Date           string         `json:"date"`
	Exercises      []ExerciseInfo `json:"exercises"`
}

// GenerationParams contains parameters for generating a workout.
type GenerationParams struct {
	UserID     string
	WeekNumber *int    // Optional: use state if not provided
	DaySlug    *string // Optional: use state if not provided
	Date       *string // Optional: defaults to today
}

// PreviewParams contains parameters for previewing a workout.
type PreviewParams struct {
	UserID     string
	WeekNumber int    // Required for preview
	DaySlug    string // Required for preview
}

// ProgramContext contains the program structure needed for workout generation.
type ProgramContext struct {
	ProgramID      string
	ProgramName    string
	CycleID        string
	CycleLengthWeeks int
	WeeklyLookupID *string
	DailyLookupID  *string
}

// UserState contains the user's current position in the program.
type UserState struct {
	CurrentWeek           int
	CurrentCycleIteration int
	CurrentDayIndex       *int
}

// DayContext contains the day information needed for workout generation.
type DayContext struct {
	DayID   string
	DaySlug string
	DayName string
}

// GenerationContext provides all dependencies needed for workout generation.
type GenerationContext struct {
	// LiftLookup is used to look up lift information during resolution.
	LiftLookup prescription.LiftLookup

	// SetGenContext provides context for set generation.
	SetGenContext setscheme.SetGenerationContext

	// LookupContext provides week/day context for lookup-based load modifications.
	LookupContext *loadstrategy.LookupContext
}

// DefaultGenerationContext returns a GenerationContext with default values.
func DefaultGenerationContext(liftLookup prescription.LiftLookup) GenerationContext {
	return GenerationContext{
		LiftLookup:    liftLookup,
		SetGenContext: setscheme.DefaultSetGenerationContext(),
	}
}

// GenerateWorkout transforms abstract program prescriptions into a concrete, user-specific workout.
//
// This is the culmination of the workout generation pipeline, combining:
//   - Program structure (which exercises in what order)
//   - User state (current week, cycle iteration for periodization context)
//   - User maxes (for load calculation)
//   - Lookup tables (for week/day-specific intensity modifications)
//
// The generation process resolves each prescription sequentially, maintaining exercise order
// as defined in the program. Exercise ordering is critical in powerlifting - compound movements
// (squat, bench, deadlift) typically come first when the lifter is freshest, followed by
// accessory work.
//
// The GenContext carries dependencies through the resolution chain:
//   - LiftLookup: Provides lift names/slugs for display
//   - SetGenContext: Configuration like work set threshold (default 80%)
//   - LookupContext: Week/day-specific modifiers for periodization
//
// If any prescription fails to resolve (e.g., missing max value), the entire workout
// generation fails. This is intentional - partial workouts could lead to imbalanced training.
func GenerateWorkout(
	ctx context.Context,
	userID string,
	programCtx ProgramContext,
	userState UserState,
	dayCtx DayContext,
	prescriptions []*prescription.Prescription,
	genCtx GenerationContext,
	date string,
) (*Workout, error) {
	if len(prescriptions) == 0 {
		return nil, ErrNoPrescriptions
	}

	exercises := make([]ExerciseInfo, 0, len(prescriptions))

	// Resolve each prescription in order. The order is preserved from the day's prescription
	// list, which was defined by the program author. This ensures primary lifts come first.
	for _, p := range prescriptions {
		resCtx := prescription.ResolutionContext{
			LiftLookup:    genCtx.LiftLookup,
			SetGenContext: genCtx.SetGenContext,
			LookupContext: genCtx.LookupContext,
		}

		resolved, err := p.Resolve(ctx, userID, resCtx)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve prescription %s: %w", p.ID, err)
		}

		exercise := ExerciseInfo{
			PrescriptionID: resolved.PrescriptionID,
			Lift: LiftInfo{
				ID:   resolved.Lift.ID,
				Name: resolved.Lift.Name,
				Slug: resolved.Lift.Slug,
			},
			Sets:        convertSets(resolved.Sets),
			Notes:       resolved.Notes,
			RestSeconds: resolved.RestSeconds,
		}

		exercises = append(exercises, exercise)
	}

	return &Workout{
		UserID:         userID,
		ProgramID:      programCtx.ProgramID,
		CycleIteration: userState.CurrentCycleIteration,
		WeekNumber:     userState.CurrentWeek,
		DaySlug:        dayCtx.DaySlug,
		Date:           date,
		Exercises:      exercises,
	}, nil
}

// convertSets converts generated sets to set info.
func convertSets(sets []setscheme.GeneratedSet) []SetInfo {
	result := make([]SetInfo, len(sets))
	for i, s := range sets {
		result[i] = SetInfo{
			SetNumber:  s.SetNumber,
			Weight:     s.Weight,
			TargetReps: s.TargetReps,
			IsWorkSet:  s.IsWorkSet,
		}
	}
	return result
}

// GetDateString returns today's date in YYYY-MM-DD format.
func GetDateString() string {
	return time.Now().Format("2006-01-02")
}

// ValidateGenerationParams validates the generation parameters.
func ValidateGenerationParams(params GenerationParams) error {
	if params.WeekNumber != nil && *params.WeekNumber < 1 {
		return ErrInvalidWeekNumber
	}
	return nil
}

// ValidatePreviewParams validates the preview parameters.
func ValidatePreviewParams(params PreviewParams) error {
	if params.WeekNumber < 1 {
		return ErrInvalidWeekNumber
	}
	if params.DaySlug == "" {
		return ErrInvalidDaySlug
	}
	return nil
}

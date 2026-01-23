// Package service provides application service layer implementations.
// This file implements the ProgressionService which handles trigger integration,
// idempotency enforcement, and atomic transactions for progression applications.
package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/waynenilsen/power-pro-v3/internal/db"
	"github.com/waynenilsen/power-pro-v3/internal/domain/progression"
)

// Errors for progression service operations.
var (
	ErrUserNotEnrolled           = errors.New("user is not enrolled in any program")
	ErrNoApplicableProgressions  = errors.New("no applicable progressions found")
	ErrProgressionAlreadyApplied = errors.New("progression already applied (idempotent skip)")
	ErrNoCurrentMax              = errors.New("no current max found for lift")
	ErrProgressionNotFound       = errors.New("progression not found")
	ErrInvalidTriggerContext     = errors.New("invalid trigger context")
	ErrTransactionFailed         = errors.New("transaction failed")
	ErrLiftNotFound              = errors.New("lift not found")
)

// ProgressionService handles trigger events and applies progressions atomically.
// It orchestrates the connection between state advancement triggers and progression application.
type ProgressionService struct {
	sqlDB   *sql.DB
	queries *db.Queries
	factory *progression.ProgressionFactory
}

// NewProgressionService creates a new ProgressionService.
func NewProgressionService(sqlDB *sql.DB, factory *progression.ProgressionFactory) *ProgressionService {
	return &ProgressionService{
		sqlDB:   sqlDB,
		queries: db.New(sqlDB),
		factory: factory,
	}
}

// TriggerResult represents the result of processing a single progression during a trigger.
type TriggerResult struct {
	ProgressionID string                         `json:"progressionId"`
	LiftID        string                         `json:"liftId"`
	Applied       bool                           `json:"applied"`
	Skipped       bool                           `json:"skipped"`
	SkipReason    string                         `json:"skipReason,omitempty"`
	Result        *progression.ProgressionResult `json:"result,omitempty"`
	Error         string                         `json:"error,omitempty"`
}

// AggregateResult represents the result of processing all progressions for a trigger event.
type AggregateResult struct {
	TriggerType  progression.TriggerType `json:"triggerType"`
	UserID       string                  `json:"userId"`
	ProgramID    string                  `json:"programId"`
	Timestamp    interface{}             `json:"timestamp"`
	Results      []TriggerResult         `json:"results"`
	TotalApplied int                     `json:"totalApplied"`
	TotalSkipped int                     `json:"totalSkipped"`
	TotalErrors  int                     `json:"totalErrors"`
}

// HandleSessionComplete handles AFTER_SESSION triggers.
// This is called when a user completes a training day.
// Only applies progressions to lifts that were performed in the session.
func (s *ProgressionService) HandleSessionComplete(ctx context.Context, event *progression.TriggerEventV2) (*AggregateResult, error) {
	if event.Type != progression.TriggerAfterSession {
		return nil, fmt.Errorf("%w: expected AFTER_SESSION, got %s", ErrInvalidTriggerContext, event.Type)
	}

	sessionCtx, ok := event.Context.(progression.SessionTriggerContext)
	if !ok {
		return nil, fmt.Errorf("%w: context is not SessionTriggerContext", ErrInvalidTriggerContext)
	}

	if err := event.Validate(); err != nil {
		return nil, fmt.Errorf("invalid trigger event: %w", err)
	}

	return s.processProgressions(ctx, event, sessionCtx.LiftsPerformed)
}

// HandleWeekAdvance handles AFTER_WEEK triggers.
// This is called when a user advances from week N to week N+1.
// Applies progressions to all configured lifts.
func (s *ProgressionService) HandleWeekAdvance(ctx context.Context, event *progression.TriggerEventV2) (*AggregateResult, error) {
	if event.Type != progression.TriggerAfterWeek {
		return nil, fmt.Errorf("%w: expected AFTER_WEEK, got %s", ErrInvalidTriggerContext, event.Type)
	}

	if err := event.Validate(); err != nil {
		return nil, fmt.Errorf("invalid trigger event: %w", err)
	}

	// For week triggers, apply to all configured lifts (nil = all lifts)
	return s.processProgressions(ctx, event, nil)
}

// HandleCycleComplete handles AFTER_CYCLE triggers.
// This is called when a user completes a cycle (week wraps to 1).
// Applies progressions to all configured lifts.
func (s *ProgressionService) HandleCycleComplete(ctx context.Context, event *progression.TriggerEventV2) (*AggregateResult, error) {
	if event.Type != progression.TriggerAfterCycle {
		return nil, fmt.Errorf("%w: expected AFTER_CYCLE, got %s", ErrInvalidTriggerContext, event.Type)
	}

	if err := event.Validate(); err != nil {
		return nil, fmt.Errorf("invalid trigger event: %w", err)
	}

	// For cycle triggers, apply to all configured lifts (nil = all lifts)
	return s.processProgressions(ctx, event, nil)
}

// processProgressions is the core method that processes all applicable progressions for a trigger event.
// It orchestrates the complete progression evaluation pipeline.
//
// The progression system follows this evaluation flow:
//
//  1. Enrollment Lookup: Verify the user is enrolled in a program (progression is meaningless without one)
//  2. Progression Discovery: Fetch all enabled progressions for the user's program
//  3. Lift Filtering: For session triggers, only consider lifts that were actually performed
//  4. Per-Progression Evaluation: For each progression-lift combination:
//     a. Check trigger type compatibility (e.g., don't fire cycle progressions on session complete)
//     b. Check idempotency (prevent duplicate applications for the same trigger event)
//     c. Apply progression atomically (fetch current max, calculate new value, persist)
//  5. Aggregate Results: Return a summary of what was applied, skipped, or errored
//
// The liftsFilter parameter serves different purposes depending on trigger type:
//   - AFTER_SESSION: Only progress lifts that were actually trained (prevents phantom progression)
//   - AFTER_WEEK/AFTER_CYCLE: Typically nil, meaning progress all configured lifts
//
// This design ensures that:
//   - A squat progression doesn't fire when the user only did bench press
//   - Progressions are never applied twice for the same trigger event
//   - Failures in one progression don't prevent others from being processed
func (s *ProgressionService) processProgressions(ctx context.Context, event *progression.TriggerEventV2, liftsFilter []string) (*AggregateResult, error) {
	// User must be enrolled in a program - progressions are always program-specific
	enrollment, err := s.queries.GetUserProgramStateByUserID(ctx, event.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotEnrolled
		}
		return nil, wrapError("failed to get user enrollment", err)
	}

	// Fetch all enabled progressions for this program, ordered by priority.
	// Priority ordering ensures consistent application order when multiple
	// progressions might affect the same lift (though this is typically avoided).
	programProgressions, err := s.queries.ListEnabledProgramProgressionsByProgram(ctx, enrollment.ProgramID)
	if err != nil {
		return nil, wrapError("failed to get program progressions", err)
	}

	// No progressions configured is a valid state - return empty results, not an error
	if len(programProgressions) == 0 {
		return &AggregateResult{
			TriggerType:  event.Type,
			UserID:       event.UserID,
			ProgramID:    enrollment.ProgramID,
			Timestamp:    event.Timestamp,
			Results:      []TriggerResult{},
			TotalApplied: 0,
			TotalSkipped: 0,
			TotalErrors:  0,
		}, nil
	}

	// Build lift filter set for O(1) lookup during per-progression processing
	liftsFilterSet := make(map[string]bool)
	for _, liftID := range liftsFilter {
		liftsFilterSet[liftID] = true
	}

	result := &AggregateResult{
		TriggerType:  event.Type,
		UserID:       event.UserID,
		ProgramID:    enrollment.ProgramID,
		Timestamp:    event.Timestamp,
		Results:      []TriggerResult{},
		TotalApplied: 0,
		TotalSkipped: 0,
		TotalErrors:  0,
	}

	// Process each progression independently - failures in one don't affect others
	for _, pp := range programProgressions {
		triggerResult := s.processSingleProgression(ctx, event, pp, liftsFilter, liftsFilterSet)
		if triggerResult != nil {
			result.Results = append(result.Results, *triggerResult)
			if triggerResult.Applied {
				result.TotalApplied++
			} else if triggerResult.Skipped {
				result.TotalSkipped++
			} else if triggerResult.Error != "" {
				result.TotalErrors++
			}
		}
	}

	return result, nil
}

// processSingleProgression processes a single program progression entry.
// Returns nil if the progression should be skipped entirely (not counted in results).
func (s *ProgressionService) processSingleProgression(
	ctx context.Context,
	event *progression.TriggerEventV2,
	pp db.ProgramProgression,
	liftsFilter []string,
	liftsFilterSet map[string]bool,
) *TriggerResult {
	// Skip if no lift ID and we have a lift filter (can't apply program-wide progression to specific lifts)
	if !pp.LiftID.Valid && liftsFilter != nil {
		return nil
	}

	// Skip if lift-specific and not in filter
	if pp.LiftID.Valid && liftsFilter != nil && !liftsFilterSet[pp.LiftID.String] {
		return nil
	}

	// Get the progression definition
	progressionDef, err := s.queries.GetProgression(ctx, pp.ProgressionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &TriggerResult{
				ProgressionID: pp.ProgressionID,
				LiftID:        pp.LiftID.String,
				Applied:       false,
				Error:         ErrProgressionNotFound.Error(),
			}
		}
		return &TriggerResult{
			ProgressionID: pp.ProgressionID,
			LiftID:        pp.LiftID.String,
			Applied:       false,
			Error:         fmt.Sprintf("failed to get progression: %v", err),
		}
	}

	// Parse the progression using the factory
	prog, err := s.factory.Create(progression.ProgressionType(progressionDef.Type), json.RawMessage(progressionDef.Parameters))
	if err != nil {
		return &TriggerResult{
			ProgressionID: pp.ProgressionID,
			LiftID:        pp.LiftID.String,
			Applied:       false,
			Error:         fmt.Sprintf("failed to parse progression: %v", err),
		}
	}

	// Check if trigger type matches
	if prog.TriggerType() != event.Type {
		return &TriggerResult{
			ProgressionID: pp.ProgressionID,
			LiftID:        pp.LiftID.String,
			Applied:       false,
			Skipped:       true,
			SkipReason:    fmt.Sprintf("trigger type mismatch: progression expects %s", prog.TriggerType()),
		}
	}

	// If no lift ID on program progression, we need to get all configured lifts
	// For now, skip program-wide progressions without specific lift
	if !pp.LiftID.Valid {
		return &TriggerResult{
			ProgressionID: pp.ProgressionID,
			Applied:       false,
			Skipped:       true,
			SkipReason:    "program-wide progressions without specific lift not yet supported",
		}
	}

	// Apply progression for this lift
	triggerResult := s.applyProgressionWithTransaction(ctx, event, pp, prog)
	return &triggerResult
}

// GetDefaultFactory returns a factory with all built-in progression types registered.
func GetDefaultFactory() *progression.ProgressionFactory {
	factory := progression.NewProgressionFactory()
	progression.RegisterLinearProgression(factory)
	progression.RegisterCycleProgression(factory)
	return factory
}


// wrapError creates a formatted error with context.
func wrapError(context string, err error) error {
	return fmt.Errorf("%s: %w", context, err)
}

// wrapErrorString creates an error with additional string context.
func wrapErrorString(baseErr error, context string) error {
	return fmt.Errorf("%w: %s", baseErr, context)
}

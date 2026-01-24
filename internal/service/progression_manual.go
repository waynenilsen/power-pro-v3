// Package service provides application service layer implementations.
// This file implements manual progression trigger functionality for testing and admin overrides.
package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/db"
	"github.com/waynenilsen/power-pro-v3/internal/domain/progression"
)

// ManualTriggerResult represents the result of a manual progression trigger.
// It may contain multiple TriggerResults if multiple lifts were affected.
type ManualTriggerResult struct {
	Results      []TriggerResult `json:"results"`
	TotalApplied int             `json:"totalApplied"`
	TotalSkipped int             `json:"totalSkipped"`
	TotalErrors  int             `json:"totalErrors"`
}

// ApplyProgressionManually applies a progression manually (for testing/admin override).
// If liftID is empty, applies to all lifts configured for this progression in the user's program.
// If force is true, bypasses idempotency check.
func (s *ProgressionService) ApplyProgressionManually(
	ctx context.Context,
	userID, progressionID, liftID string,
	force bool,
) (*ManualTriggerResult, error) {
	// Get the progression definition
	progressionDef, err := s.queries.GetProgression(ctx, progressionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProgressionNotFound
		}
		return nil, wrapError("failed to get progression", err)
	}

	// Parse the progression using the factory
	prog, err := s.factory.Create(progression.ProgressionType(progressionDef.Type), json.RawMessage(progressionDef.Parameters))
	if err != nil {
		return nil, wrapError("failed to parse progression", err)
	}

	// Get user enrollment to find program
	enrollment, err := s.queries.GetUserProgramStateByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotEnrolled
		}
		return nil, wrapError("failed to get user enrollment", err)
	}

	// Determine which lifts to apply the progression to
	var liftIDs []string
	if liftID != "" {
		// Verify the lift exists
		_, err := s.queries.GetLift(ctx, liftID)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrLiftNotFound
			}
			return nil, wrapError("failed to get lift", err)
		}
		liftIDs = []string{liftID}
	} else {
		// Get all lifts configured for this progression in the user's program
		programProgressions, err := s.queries.ListEnabledProgramProgressionsByProgramAndProgression(ctx, db.ListEnabledProgramProgressionsByProgramAndProgressionParams{
			ProgramID:     enrollment.ProgramID,
			ProgressionID: progressionID,
		})
		if err != nil {
			return nil, wrapError("failed to get program progressions", err)
		}

		if len(programProgressions) == 0 {
			return nil, ErrNoApplicableProgressions
		}

		// Collect unique lift IDs (only those that are lift-specific)
		seen := make(map[string]bool)
		for _, pp := range programProgressions {
			if pp.LiftID.Valid && !seen[pp.LiftID.String] {
				liftIDs = append(liftIDs, pp.LiftID.String)
				seen[pp.LiftID.String] = true
			}
		}

		if len(liftIDs) == 0 {
			return nil, wrapErrorString(ErrNoApplicableProgressions, "progression has no lift-specific configurations")
		}
	}

	// Build the results
	result := &ManualTriggerResult{
		Results: make([]TriggerResult, 0, len(liftIDs)),
	}

	// Apply progression to each lift
	for _, lid := range liftIDs {
		triggerResult := s.applyManualProgressionToLift(ctx, userID, progressionID, lid, enrollment.ProgramID, prog, force)
		result.Results = append(result.Results, triggerResult)

		if triggerResult.Applied {
			result.TotalApplied++
		} else if triggerResult.Skipped {
			result.TotalSkipped++
		} else if triggerResult.Error != "" {
			result.TotalErrors++
		}
	}

	return result, nil
}

// applyManualProgressionToLift applies a manual progression to a specific lift.
func (s *ProgressionService) applyManualProgressionToLift(
	ctx context.Context,
	userID, progressionID, liftID, programID string,
	prog progression.Progression,
	force bool,
) TriggerResult {
	// Build a synthetic trigger event with ManualTriggerContext
	now := time.Now()

	// Create underlying context based on trigger type
	var underlyingContext progression.TriggerContext
	switch prog.TriggerType() {
	case progression.TriggerAfterSession:
		underlyingContext = progression.SessionTriggerContext{
			SessionID:      "manual-trigger",
			DaySlug:        "manual",
			WeekNumber:     1,
			LiftsPerformed: []string{liftID},
		}
	case progression.TriggerAfterWeek:
		underlyingContext = progression.WeekTriggerContext{
			PreviousWeek:   1,
			NewWeek:        2,
			CycleIteration: 1,
		}
	case progression.TriggerAfterCycle:
		underlyingContext = progression.CycleTriggerContext{
			CompletedCycle: 1,
			NewCycle:       2,
			TotalWeeks:     4,
		}
	case progression.TriggerOnFailure:
		// For ON_FAILURE triggers, we need to get the current failure count
		// from the failure counter repository and create a synthetic failure context.
		// Since this is a manual trigger, we use synthetic values for logged set info.
		consecutiveFailures := 1 // Default if no counter exists
		counter, err := s.queries.GetFailureCounterByKey(ctx, db.GetFailureCounterByKeyParams{
			UserID:        userID,
			LiftID:        liftID,
			ProgressionID: progressionID,
		})
		if err == nil {
			consecutiveFailures = int(counter.ConsecutiveFailures)
		}
		// If force mode is enabled, ensure at least 1 failure is reported
		if force && consecutiveFailures < 1 {
			consecutiveFailures = 1
		}
		underlyingContext = progression.FailureTriggerContext{
			LoggedSetID:         "manual-trigger",
			LiftID:              liftID,
			TargetReps:          5,              // Synthetic - doesn't affect progression logic
			RepsPerformed:       0,              // Synthetic - indicates failure
			RepsDifference:      -5,             // Synthetic
			ConsecutiveFailures: consecutiveFailures,
			Weight:              0,              // Will be looked up during progression application
			ProgressionID:       progressionID,
		}
	}

	// Wrap with ManualTriggerContext for audit purposes
	manualContext := progression.NewManualTriggerContext(underlyingContext, liftID, force)

	event := &progression.TriggerEventV2{
		Type:      prog.TriggerType(),
		UserID:    userID,
		Timestamp: now,
		Context:   manualContext,
	}

	// Build a synthetic program progression entry
	pp := db.ProgramProgression{
		ID:            "manual",
		ProgramID:     programID,
		ProgressionID: progressionID,
		LiftID:        sql.NullString{String: liftID, Valid: true},
		Priority:      0,
		Enabled:       1,
	}

	// If force is true, use applyProgressionWithTransactionForce which bypasses idempotency
	if force {
		return s.applyProgressionWithTransactionForce(ctx, event, pp, prog)
	}

	return s.applyProgressionWithTransaction(ctx, event, pp, prog)
}

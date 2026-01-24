// Package service provides application service layer implementations.
// This file implements the FailureService which handles failure detection and tracking.
package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/waynenilsen/power-pro-v3/internal/db"
	"github.com/waynenilsen/power-pro-v3/internal/domain/loggedset"
	"github.com/waynenilsen/power-pro-v3/internal/domain/progression"
	"github.com/waynenilsen/power-pro-v3/internal/repository"
)

// Errors for failure service operations.
var (
	ErrLoggedSetRequired = errors.New("logged set is required")
	ErrNotAFailure       = errors.New("set did not fail (reps_performed >= target_reps)")
)

// FailureService handles failure detection and counter management.
// It integrates with LoggedSet creation to track failures and fire OnFailure triggers.
type FailureService struct {
	sqlDB       *sql.DB
	queries     *db.Queries
	counterRepo *repository.FailureCounterRepository
	factory     *progression.ProgressionFactory
}

// NewFailureService creates a new FailureService.
func NewFailureService(sqlDB *sql.DB, factory *progression.ProgressionFactory) *FailureService {
	return &FailureService{
		sqlDB:       sqlDB,
		queries:     db.New(sqlDB),
		counterRepo: repository.NewFailureCounterRepository(sqlDB),
		factory:     factory,
	}
}

// FailureCheckResult contains the result of checking a logged set for failure.
type FailureCheckResult struct {
	// IsFailure indicates if the set was a failure (reps_performed < target_reps).
	IsFailure bool
	// ConsecutiveFailures is the current count after processing (0 if success).
	ConsecutiveFailures int
	// ProgressionID is the progression this failure was tracked against.
	ProgressionID string
	// TriggerFired indicates if an OnFailure trigger was fired.
	TriggerFired bool
}

// ProcessSetResult contains all failure check results for a logged set.
type ProcessSetResult struct {
	// LoggedSetID is the ID of the logged set that was processed.
	LoggedSetID string
	// Results contains failure check results for each applicable progression.
	Results []FailureCheckResult
}

// CheckForFailure determines if a logged set is a failure.
// A failure is defined as: reps_performed < target_reps.
func (s *FailureService) CheckForFailure(ls *loggedset.LoggedSet) bool {
	if ls == nil {
		return false
	}
	return ls.RepsPerformed < ls.TargetReps
}

// IsSuccess determines if a logged set is a success.
// A success is defined as: reps_performed >= target_reps.
func (s *FailureService) IsSuccess(ls *loggedset.LoggedSet) bool {
	if ls == nil {
		return false
	}
	return ls.RepsPerformed >= ls.TargetReps
}

// ProcessLoggedSet processes a logged set to update failure counters and fire triggers.
// This should be called after a LoggedSet is created.
//
// For each applicable progression:
//   - If failure: increment failure counter, potentially fire OnFailure trigger
//   - If success: reset failure counter
//
// Returns results for each progression that was processed.
func (s *FailureService) ProcessLoggedSet(ctx context.Context, ls *loggedset.LoggedSet) (*ProcessSetResult, error) {
	if ls == nil {
		return nil, ErrLoggedSetRequired
	}

	// Get the user's enrolled program
	enrollment, err := s.queries.GetUserProgramStateByUserID(ctx, ls.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			// User not enrolled in a program, nothing to do
			return &ProcessSetResult{
				LoggedSetID: ls.ID,
				Results:     []FailureCheckResult{},
			}, nil
		}
		return nil, fmt.Errorf("failed to get user enrollment: %w", err)
	}

	// Get all enabled progressions for the program that use OnFailure trigger
	programProgressions, err := s.queries.ListEnabledProgramProgressionsByProgram(ctx, enrollment.ProgramID)
	if err != nil {
		return nil, fmt.Errorf("failed to get program progressions: %w", err)
	}

	isFailure := s.CheckForFailure(ls)
	result := &ProcessSetResult{
		LoggedSetID: ls.ID,
		Results:     []FailureCheckResult{},
	}

	for _, pp := range programProgressions {
		// Skip if not for this lift
		if pp.LiftID.Valid && pp.LiftID.String != ls.LiftID {
			continue
		}

		// Get the progression definition
		progressionDef, err := s.queries.GetProgression(ctx, pp.ProgressionID)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			return nil, fmt.Errorf("failed to get progression %s: %w", pp.ProgressionID, err)
		}

		// Parse the progression to check its trigger type
		prog, err := s.factory.Create(progression.ProgressionType(progressionDef.Type), []byte(progressionDef.Parameters))
		if err != nil {
			continue // Skip progressions that can't be parsed
		}

		// Process based on whether this is a failure or success
		var checkResult FailureCheckResult
		checkResult.ProgressionID = pp.ProgressionID

		if isFailure {
			// Increment failure counter
			count, err := s.counterRepo.IncrementOnFailure(ls.UserID, ls.LiftID, pp.ProgressionID)
			if err != nil {
				return nil, fmt.Errorf("failed to increment failure counter: %w", err)
			}
			checkResult.IsFailure = true
			checkResult.ConsecutiveFailures = count

			// Check if this progression responds to OnFailure triggers
			if prog.TriggerType() == progression.TriggerOnFailure {
				checkResult.TriggerFired = true
				// Note: Actual trigger processing would be done by ProgressionService
				// This service just tracks failures and prepares the trigger context
			}
		} else {
			// Reset failure counter on success
			err := s.counterRepo.ResetOnSuccess(ls.UserID, ls.LiftID, pp.ProgressionID)
			if err != nil {
				return nil, fmt.Errorf("failed to reset failure counter: %w", err)
			}
			checkResult.IsFailure = false
			checkResult.ConsecutiveFailures = 0
		}

		result.Results = append(result.Results, checkResult)
	}

	return result, nil
}

// GetFailureCount retrieves the current consecutive failure count for a user/lift/progression.
// Returns 0 if no counter exists.
func (s *FailureService) GetFailureCount(userID, liftID, progressionID string) (int, error) {
	return s.counterRepo.GetFailureCount(userID, liftID, progressionID)
}

// IncrementFailureCounter increments the failure counter for a user/lift/progression.
// Returns the new failure count.
func (s *FailureService) IncrementFailureCounter(userID, liftID, progressionID string) (int, error) {
	return s.counterRepo.IncrementOnFailure(userID, liftID, progressionID)
}

// ResetFailureCounter resets the failure counter for a user/lift/progression.
func (s *FailureService) ResetFailureCounter(userID, liftID, progressionID string) error {
	return s.counterRepo.ResetOnSuccess(userID, liftID, progressionID)
}

// BuildFailureTriggerContext creates a FailureTriggerContext from a logged set and failure count.
func (s *FailureService) BuildFailureTriggerContext(ls *loggedset.LoggedSet, consecutiveFailures int, progressionID string) progression.FailureTriggerContext {
	return progression.FailureTriggerContext{
		LoggedSetID:         ls.ID,
		LiftID:              ls.LiftID,
		TargetReps:          ls.TargetReps,
		RepsPerformed:       ls.RepsPerformed,
		RepsDifference:      ls.RepsDifference(),
		ConsecutiveFailures: consecutiveFailures,
		Weight:              ls.Weight,
		ProgressionID:       progressionID,
	}
}

// CreateFailureTriggerEvent creates a complete TriggerEventV2 for an OnFailure trigger.
func (s *FailureService) CreateFailureTriggerEvent(ls *loggedset.LoggedSet, consecutiveFailures int, progressionID string) *progression.TriggerEventV2 {
	ctx := s.BuildFailureTriggerContext(ls, consecutiveFailures, progressionID)
	return progression.NewFailureTriggerEvent(ls.UserID, ctx)
}

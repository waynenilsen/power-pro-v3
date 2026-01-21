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
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/db"
	"github.com/waynenilsen/power-pro-v3/internal/domain/liftmax"
	"github.com/waynenilsen/power-pro-v3/internal/domain/progression"
)

// Errors for progression service operations.
var (
	ErrUserNotEnrolled             = errors.New("user is not enrolled in any program")
	ErrNoApplicableProgressions    = errors.New("no applicable progressions found")
	ErrProgressionAlreadyApplied   = errors.New("progression already applied (idempotent skip)")
	ErrNoCurrentMax                = errors.New("no current max found for lift")
	ErrProgressionNotFound         = errors.New("progression not found")
	ErrInvalidTriggerContext       = errors.New("invalid trigger context")
	ErrTransactionFailed           = errors.New("transaction failed")
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
	ProgressionID   string               `json:"progressionId"`
	LiftID          string               `json:"liftId"`
	Applied         bool                 `json:"applied"`
	Skipped         bool                 `json:"skipped"`
	SkipReason      string               `json:"skipReason,omitempty"`
	Result          *progression.ProgressionResult `json:"result,omitempty"`
	Error           string               `json:"error,omitempty"`
}

// AggregateResult represents the result of processing all progressions for a trigger event.
type AggregateResult struct {
	TriggerType     progression.TriggerType `json:"triggerType"`
	UserID          string                  `json:"userId"`
	ProgramID       string                  `json:"programId"`
	Timestamp       time.Time               `json:"timestamp"`
	Results         []TriggerResult         `json:"results"`
	TotalApplied    int                     `json:"totalApplied"`
	TotalSkipped    int                     `json:"totalSkipped"`
	TotalErrors     int                     `json:"totalErrors"`
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
// It handles lookup, filtering, idempotency checks, and atomic application.
func (s *ProgressionService) processProgressions(ctx context.Context, event *progression.TriggerEventV2, liftsFilter []string) (*AggregateResult, error) {
	// Get user's enrolled program
	enrollment, err := s.queries.GetUserProgramStateByUserID(ctx, event.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotEnrolled
		}
		return nil, fmt.Errorf("failed to get user enrollment: %w", err)
	}

	// Fetch enabled program progressions ordered by priority
	programProgressions, err := s.queries.ListEnabledProgramProgressionsByProgram(ctx, enrollment.ProgramID)
	if err != nil {
		return nil, fmt.Errorf("failed to get program progressions: %w", err)
	}

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

	// Build a set of lifts to filter by (if any)
	liftsFilterSet := make(map[string]bool)
	if liftsFilter != nil {
		for _, liftID := range liftsFilter {
			liftsFilterSet[liftID] = true
		}
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

	// Process each program progression
	for _, pp := range programProgressions {
		// Skip if no lift ID and we have a lift filter (can't apply program-wide progression to specific lifts)
		if !pp.LiftID.Valid && liftsFilter != nil {
			continue
		}

		// Skip if lift-specific and not in filter
		if pp.LiftID.Valid && liftsFilter != nil && !liftsFilterSet[pp.LiftID.String] {
			continue
		}

		// Get the progression definition
		progressionDef, err := s.queries.GetProgression(ctx, pp.ProgressionID)
		if err != nil {
			if err == sql.ErrNoRows {
				result.Results = append(result.Results, TriggerResult{
					ProgressionID: pp.ProgressionID,
					LiftID:        pp.LiftID.String,
					Applied:       false,
					Error:         ErrProgressionNotFound.Error(),
				})
				result.TotalErrors++
				continue
			}
			return nil, fmt.Errorf("failed to get progression: %w", err)
		}

		// Parse the progression using the factory
		prog, err := s.factory.Create(progression.ProgressionType(progressionDef.Type), json.RawMessage(progressionDef.Parameters))
		if err != nil {
			result.Results = append(result.Results, TriggerResult{
				ProgressionID: pp.ProgressionID,
				LiftID:        pp.LiftID.String,
				Applied:       false,
				Error:         fmt.Sprintf("failed to parse progression: %v", err),
			})
			result.TotalErrors++
			continue
		}

		// Check if trigger type matches
		if prog.TriggerType() != event.Type {
			result.Results = append(result.Results, TriggerResult{
				ProgressionID: pp.ProgressionID,
				LiftID:        pp.LiftID.String,
				Applied:       false,
				Skipped:       true,
				SkipReason:    fmt.Sprintf("trigger type mismatch: progression expects %s", prog.TriggerType()),
			})
			result.TotalSkipped++
			continue
		}

		// If no lift ID on program progression, we need to get all configured lifts
		// For now, skip program-wide progressions without specific lift
		if !pp.LiftID.Valid {
			result.Results = append(result.Results, TriggerResult{
				ProgressionID: pp.ProgressionID,
				Applied:       false,
				Skipped:       true,
				SkipReason:    "program-wide progressions without specific lift not yet supported",
			})
			result.TotalSkipped++
			continue
		}

		// Apply progression for this lift
		triggerResult := s.applyProgressionWithTransaction(ctx, event, pp, prog)
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

// applyProgressionWithTransaction applies a single progression in an atomic transaction.
// It performs idempotency check, updates LiftMax, and logs the application.
func (s *ProgressionService) applyProgressionWithTransaction(
	ctx context.Context,
	event *progression.TriggerEventV2,
	pp db.ProgramProgression,
	prog progression.Progression,
) TriggerResult {
	liftID := pp.LiftID.String

	// Begin transaction
	tx, err := s.sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return TriggerResult{
			ProgressionID: pp.ProgressionID,
			LiftID:        liftID,
			Applied:       false,
			Error:         fmt.Sprintf("failed to begin transaction: %v", err),
		}
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	txQueries := db.New(tx)

	// Check idempotency (per user, progression, lift, trigger type, and timestamp)
	appliedAtStr := event.Timestamp.Format(time.RFC3339)
	alreadyApplied, err := txQueries.CheckIdempotency(ctx, db.CheckIdempotencyParams{
		UserID:        event.UserID,
		ProgressionID: pp.ProgressionID,
		LiftID:        liftID,
		TriggerType:   string(event.Type),
		AppliedAt:     appliedAtStr,
	})
	if err != nil {
		return TriggerResult{
			ProgressionID: pp.ProgressionID,
			LiftID:        liftID,
			Applied:       false,
			Error:         fmt.Sprintf("failed to check idempotency: %v", err),
		}
	}
	if alreadyApplied == 1 {
		tx.Rollback()
		return TriggerResult{
			ProgressionID: pp.ProgressionID,
			LiftID:        liftID,
			Applied:       false,
			Skipped:       true,
			SkipReason:    "already applied (idempotent skip)",
		}
	}

	// Get current max for the lift
	// Determine max type based on progression
	var maxType progression.MaxType
	switch p := prog.(type) {
	case *progression.LinearProgression:
		maxType = p.MaxTypeValue
	case *progression.CycleProgression:
		maxType = p.MaxTypeValue
	default:
		return TriggerResult{
			ProgressionID: pp.ProgressionID,
			LiftID:        liftID,
			Applied:       false,
			Error:         "unknown progression type",
		}
	}

	currentMax, err := txQueries.GetCurrentMax(ctx, db.GetCurrentMaxParams{
		UserID: event.UserID,
		LiftID: liftID,
		Type:   string(maxType),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			tx.Rollback()
			return TriggerResult{
				ProgressionID: pp.ProgressionID,
				LiftID:        liftID,
				Applied:       false,
				Skipped:       true,
				SkipReason:    fmt.Sprintf("no current %s found for lift", maxType),
			}
		}
		return TriggerResult{
			ProgressionID: pp.ProgressionID,
			LiftID:        liftID,
			Applied:       false,
			Error:         fmt.Sprintf("failed to get current max: %v", err),
		}
	}

	// Build progression context
	// Convert trigger context to the flat TriggerEvent structure used by progression.Apply
	triggerEvent := buildTriggerEvent(event)

	progressionCtx := progression.ProgressionContext{
		UserID:       event.UserID,
		LiftID:       liftID,
		MaxType:      maxType,
		CurrentValue: currentMax.Value,
		TriggerEvent: triggerEvent,
	}

	// Apply the progression
	var progressionResult progression.ProgressionResult

	// Handle override increment for CycleProgression
	if cp, ok := prog.(*progression.CycleProgression); ok && pp.OverrideIncrement.Valid {
		override := pp.OverrideIncrement.Float64
		progressionResult, err = cp.ApplyWithOverride(ctx, progressionCtx, &override)
	} else {
		progressionResult, err = prog.Apply(ctx, progressionCtx)
	}

	if err != nil {
		return TriggerResult{
			ProgressionID: pp.ProgressionID,
			LiftID:        liftID,
			Applied:       false,
			Error:         fmt.Sprintf("failed to apply progression: %v", err),
		}
	}

	if !progressionResult.Applied {
		tx.Rollback()
		return TriggerResult{
			ProgressionID: pp.ProgressionID,
			LiftID:        liftID,
			Applied:       false,
			Skipped:       true,
			SkipReason:    progressionResult.Reason,
			Result:        &progressionResult,
		}
	}

	// Create new LiftMax entry
	// Use the event timestamp as the effective date (this ensures uniqueness per trigger)
	newMaxID := uuid.New().String()
	now := time.Now()
	nowStr := now.Format(time.RFC3339)
	effectiveDateStr := event.Timestamp.Format(time.RFC3339)

	err = txQueries.CreateLiftMax(ctx, db.CreateLiftMaxParams{
		ID:            newMaxID,
		UserID:        event.UserID,
		LiftID:        liftID,
		Type:          string(maxType),
		Value:         progressionResult.NewValue,
		EffectiveDate: effectiveDateStr,
		CreatedAt:     nowStr,
		UpdatedAt:     nowStr,
	})
	if err != nil {
		return TriggerResult{
			ProgressionID: pp.ProgressionID,
			LiftID:        liftID,
			Applied:       false,
			Error:         fmt.Sprintf("failed to create new lift max: %v", err),
		}
	}

	// Create progression log entry
	triggerContextJSON, err := json.Marshal(event.Context)
	if err != nil {
		return TriggerResult{
			ProgressionID: pp.ProgressionID,
			LiftID:        liftID,
			Applied:       false,
			Error:         fmt.Sprintf("failed to serialize trigger context: %v", err),
		}
	}

	logID := uuid.New().String()
	err = txQueries.CreateProgressionLog(ctx, db.CreateProgressionLogParams{
		ID:             logID,
		UserID:         event.UserID,
		ProgressionID:  pp.ProgressionID,
		LiftID:         liftID,
		PreviousValue:  progressionResult.PreviousValue,
		NewValue:       progressionResult.NewValue,
		Delta:          progressionResult.Delta,
		TriggerType:    string(event.Type),
		TriggerContext: string(triggerContextJSON),
		AppliedAt:      appliedAtStr,
	})
	if err != nil {
		return TriggerResult{
			ProgressionID: pp.ProgressionID,
			LiftID:        liftID,
			Applied:       false,
			Error:         fmt.Sprintf("failed to create progression log: %v", err),
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return TriggerResult{
			ProgressionID: pp.ProgressionID,
			LiftID:        liftID,
			Applied:       false,
			Error:         fmt.Sprintf("failed to commit transaction: %v", err),
		}
	}

	return TriggerResult{
		ProgressionID: pp.ProgressionID,
		LiftID:        liftID,
		Applied:       true,
		Result:        &progressionResult,
	}
}

// buildTriggerEvent converts a TriggerEventV2 to the flat TriggerEvent structure.
// This bridges the new strongly-typed trigger context with the existing progression interface.
func buildTriggerEvent(event *progression.TriggerEventV2) progression.TriggerEvent {
	triggerEvent := progression.TriggerEvent{
		Type:      event.Type,
		Timestamp: event.Timestamp,
	}

	switch ctx := event.Context.(type) {
	case progression.SessionTriggerContext:
		triggerEvent.SessionID = &ctx.SessionID
		triggerEvent.DaySlug = &ctx.DaySlug
		triggerEvent.WeekNumber = &ctx.WeekNumber
		triggerEvent.LiftsPerformed = ctx.LiftsPerformed
	case progression.WeekTriggerContext:
		triggerEvent.WeekNumber = &ctx.NewWeek
		triggerEvent.CycleIteration = &ctx.CycleIteration
	case progression.CycleTriggerContext:
		triggerEvent.CycleIteration = &ctx.CompletedCycle
	}

	return triggerEvent
}

// ManualTriggerResult represents the result of a manual progression trigger.
// It may contain multiple TriggerResults if multiple lifts were affected.
type ManualTriggerResult struct {
	Results      []TriggerResult `json:"results"`
	TotalApplied int             `json:"totalApplied"`
	TotalSkipped int             `json:"totalSkipped"`
	TotalErrors  int             `json:"totalErrors"`
}

// ErrLiftNotFound is returned when a specified lift does not exist.
var ErrLiftNotFound = errors.New("lift not found")

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
		return nil, fmt.Errorf("failed to get progression: %w", err)
	}

	// Parse the progression using the factory
	prog, err := s.factory.Create(progression.ProgressionType(progressionDef.Type), json.RawMessage(progressionDef.Parameters))
	if err != nil {
		return nil, fmt.Errorf("failed to parse progression: %w", err)
	}

	// Get user enrollment to find program
	enrollment, err := s.queries.GetUserProgramStateByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotEnrolled
		}
		return nil, fmt.Errorf("failed to get user enrollment: %w", err)
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
			return nil, fmt.Errorf("failed to get lift: %w", err)
		}
		liftIDs = []string{liftID}
	} else {
		// Get all lifts configured for this progression in the user's program
		programProgressions, err := s.queries.ListEnabledProgramProgressionsByProgramAndProgression(ctx, db.ListEnabledProgramProgressionsByProgramAndProgressionParams{
			ProgramID:     enrollment.ProgramID,
			ProgressionID: progressionID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get program progressions: %w", err)
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
			return nil, fmt.Errorf("%w: progression has no lift-specific configurations", ErrNoApplicableProgressions)
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

// applyProgressionWithTransactionForce applies a progression bypassing idempotency checks.
// It is identical to applyProgressionWithTransaction except it skips the idempotency check.
// This is used for manual force=true triggers.
func (s *ProgressionService) applyProgressionWithTransactionForce(
	ctx context.Context,
	event *progression.TriggerEventV2,
	pp db.ProgramProgression,
	prog progression.Progression,
) TriggerResult {
	liftID := pp.LiftID.String

	// Begin transaction
	tx, err := s.sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return TriggerResult{
			ProgressionID: pp.ProgressionID,
			LiftID:        liftID,
			Applied:       false,
			Error:         fmt.Sprintf("failed to begin transaction: %v", err),
		}
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	txQueries := db.New(tx)

	// SKIP idempotency check for force mode

	// Determine max type based on progression
	var maxType progression.MaxType
	switch p := prog.(type) {
	case *progression.LinearProgression:
		maxType = p.MaxTypeValue
	case *progression.CycleProgression:
		maxType = p.MaxTypeValue
	default:
		return TriggerResult{
			ProgressionID: pp.ProgressionID,
			LiftID:        liftID,
			Applied:       false,
			Error:         "unknown progression type",
		}
	}

	currentMax, err := txQueries.GetCurrentMax(ctx, db.GetCurrentMaxParams{
		UserID: event.UserID,
		LiftID: liftID,
		Type:   string(maxType),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			tx.Rollback()
			return TriggerResult{
				ProgressionID: pp.ProgressionID,
				LiftID:        liftID,
				Applied:       false,
				Skipped:       true,
				SkipReason:    fmt.Sprintf("no current %s found for lift", maxType),
			}
		}
		return TriggerResult{
			ProgressionID: pp.ProgressionID,
			LiftID:        liftID,
			Applied:       false,
			Error:         fmt.Sprintf("failed to get current max: %v", err),
		}
	}

	// Build progression context
	triggerEvent := buildTriggerEvent(event)

	progressionCtx := progression.ProgressionContext{
		UserID:       event.UserID,
		LiftID:       liftID,
		MaxType:      maxType,
		CurrentValue: currentMax.Value,
		TriggerEvent: triggerEvent,
	}

	// Apply the progression
	var progressionResult progression.ProgressionResult

	// Handle override increment for CycleProgression
	if cp, ok := prog.(*progression.CycleProgression); ok && pp.OverrideIncrement.Valid {
		override := pp.OverrideIncrement.Float64
		progressionResult, err = cp.ApplyWithOverride(ctx, progressionCtx, &override)
	} else {
		progressionResult, err = prog.Apply(ctx, progressionCtx)
	}

	if err != nil {
		return TriggerResult{
			ProgressionID: pp.ProgressionID,
			LiftID:        liftID,
			Applied:       false,
			Error:         fmt.Sprintf("failed to apply progression: %v", err),
		}
	}

	if !progressionResult.Applied {
		tx.Rollback()
		return TriggerResult{
			ProgressionID: pp.ProgressionID,
			LiftID:        liftID,
			Applied:       false,
			Skipped:       true,
			SkipReason:    progressionResult.Reason,
			Result:        &progressionResult,
		}
	}

	// Create new LiftMax entry with a unique timestamp to avoid conflicts
	// Use nanoseconds to ensure uniqueness for force mode
	newMaxID := uuid.New().String()
	now := time.Now()
	nowStr := now.Format(time.RFC3339Nano) // Use RFC3339Nano for uniqueness
	effectiveDateStr := now.Format(time.RFC3339Nano)

	err = txQueries.CreateLiftMax(ctx, db.CreateLiftMaxParams{
		ID:            newMaxID,
		UserID:        event.UserID,
		LiftID:        liftID,
		Type:          string(maxType),
		Value:         progressionResult.NewValue,
		EffectiveDate: effectiveDateStr,
		CreatedAt:     nowStr,
		UpdatedAt:     nowStr,
	})
	if err != nil {
		return TriggerResult{
			ProgressionID: pp.ProgressionID,
			LiftID:        liftID,
			Applied:       false,
			Error:         fmt.Sprintf("failed to create new lift max: %v", err),
		}
	}

	// Create progression log entry with ManualTriggerContext
	triggerContextJSON, err := json.Marshal(event.Context)
	if err != nil {
		return TriggerResult{
			ProgressionID: pp.ProgressionID,
			LiftID:        liftID,
			Applied:       false,
			Error:         fmt.Sprintf("failed to serialize trigger context: %v", err),
		}
	}

	// Use unique timestamp for idempotency (for force mode logs)
	appliedAtStr := now.Format(time.RFC3339Nano)

	logID := uuid.New().String()
	err = txQueries.CreateProgressionLog(ctx, db.CreateProgressionLogParams{
		ID:             logID,
		UserID:         event.UserID,
		ProgressionID:  pp.ProgressionID,
		LiftID:         liftID,
		PreviousValue:  progressionResult.PreviousValue,
		NewValue:       progressionResult.NewValue,
		Delta:          progressionResult.Delta,
		TriggerType:    string(event.Type),
		TriggerContext: string(triggerContextJSON),
		AppliedAt:      appliedAtStr,
	})
	if err != nil {
		return TriggerResult{
			ProgressionID: pp.ProgressionID,
			LiftID:        liftID,
			Applied:       false,
			Error:         fmt.Sprintf("failed to create progression log: %v", err),
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return TriggerResult{
			ProgressionID: pp.ProgressionID,
			LiftID:        liftID,
			Applied:       false,
			Error:         fmt.Sprintf("failed to commit transaction: %v", err),
		}
	}

	return TriggerResult{
		ProgressionID: pp.ProgressionID,
		LiftID:        liftID,
		Applied:       true,
		Result:        &progressionResult,
	}
}

// GetDefaultFactory returns a factory with all built-in progression types registered.
func GetDefaultFactory() *progression.ProgressionFactory {
	factory := progression.NewProgressionFactory()
	progression.RegisterLinearProgression(factory)
	progression.RegisterCycleProgression(factory)
	return factory
}

// MaxType conversion helper for the service layer
func maxTypeToLiftMaxType(mt progression.MaxType) liftmax.MaxType {
	switch mt {
	case progression.OneRM:
		return liftmax.OneRM
	case progression.TrainingMax:
		return liftmax.TrainingMax
	default:
		return liftmax.OneRM
	}
}

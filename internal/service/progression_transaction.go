// Package service provides application service layer implementations.
// This file implements transaction handling for progression application.
package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/db"
	"github.com/waynenilsen/power-pro-v3/internal/domain/progression"
)

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
			_ = tx.Rollback()
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
		_ = tx.Rollback()
		return TriggerResult{
			ProgressionID: pp.ProgressionID,
			LiftID:        liftID,
			Applied:       false,
			Skipped:       true,
			SkipReason:    "already applied (idempotent skip)",
		}
	}

	// Get current max and apply progression
	return s.applyProgressionCore(ctx, tx, txQueries, event, pp, prog, liftID, appliedAtStr)
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
			_ = tx.Rollback()
		}
	}()

	txQueries := db.New(tx)

	// SKIP idempotency check for force mode
	// Use unique timestamp for force mode
	// NOTE: Add 1 second to ensure RFC3339Nano timestamps sort AFTER RFC3339 timestamps
	// SQLite sorts strings lexicographically. When comparing same-second timestamps:
	// - "2024-01-01T00:00:00Z" (RFC3339)
	// - "2024-01-01T00:00:00.123456789Z" (RFC3339Nano at same second)
	// The 'Z' > '.' so RFC3339 sorts higher! Adding 1 second fixes this.
	now := time.Now().Add(time.Second)
	appliedAtStr := now.Format(time.RFC3339Nano)

	return s.applyProgressionCore(ctx, tx, txQueries, event, pp, prog, liftID, appliedAtStr)
}

// applyProgressionCore contains the shared logic for applying a progression.
// It handles max lookup, progression application, LiftMax creation, and logging.
func (s *ProgressionService) applyProgressionCore(
	ctx context.Context,
	tx *sql.Tx,
	txQueries *db.Queries,
	event *progression.TriggerEventV2,
	pp db.ProgramProgression,
	prog progression.Progression,
	liftID string,
	appliedAtStr string,
) TriggerResult {
	// Determine max type based on progression
	maxType, err := getMaxTypeFromProgression(prog)
	if err != nil {
		return TriggerResult{
			ProgressionID: pp.ProgressionID,
			LiftID:        liftID,
			Applied:       false,
			Error:         err.Error(),
		}
	}

	currentMax, err := txQueries.GetCurrentMax(ctx, db.GetCurrentMaxParams{
		UserID: event.UserID,
		LiftID: liftID,
		Type:   string(maxType),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			_ = tx.Rollback()
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

	// For StageProgression, load the current stage from persistent storage
	// This is necessary because the progression is re-created from stored params each time
	var stageProg *progression.StageProgression
	if sp, ok := prog.(*progression.StageProgression); ok {
		stageProg = sp
		// Look up current stage from user_progression_states
		state, err := txQueries.GetUserProgressionState(ctx, db.GetUserProgressionStateParams{
			UserID:        event.UserID,
			LiftID:        liftID,
			ProgressionID: pp.ProgressionID,
		})
		if err == nil {
			// Found existing state, set the current stage
			if setErr := stageProg.SetCurrentStage(int(state.CurrentStage)); setErr != nil {
				// Stage out of bounds, reset to 0
				_ = stageProg.SetCurrentStage(0)
			}
		}
		// If not found (sql.ErrNoRows), use default stage 0
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
		_ = tx.Rollback()
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
	// NOTE on timestamp formats and SQLite sorting:
	// SQLite sorts strings lexicographically. When comparing timestamps:
	// - "2024-01-01T00:00:00Z" (RFC3339)
	// - "2024-01-01T00:00:00.123456789Z" (RFC3339Nano at same second)
	// The 'Z' (ASCII 90) > '.' (ASCII 46), so RFC3339 sorts HIGHER than RFC3339Nano at same second!
	// This would cause GetCurrentMax (ORDER BY effective_date DESC) to return stale values.
	//
	// Solution: For force mode (which uses RFC3339Nano), we add 1 second to the timestamp.
	// This ensures "2024-01-01T00:00:01.xxx" sorts AFTER "2024-01-01T00:00:00Z".
	now := time.Now().Add(time.Second)
	nowStr := now.Format(time.RFC3339Nano)
	effectiveDateStr := event.Timestamp.Format(time.RFC3339)

	// For force mode, use the unique RFC3339Nano timestamp (which has +1 second offset applied)
	if appliedAtStr != event.Timestamp.Format(time.RFC3339) {
		effectiveDateStr = appliedAtStr
	}

	newMaxID := uuid.New().String()
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
		TriggerContext: sql.NullString{String: string(triggerContextJSON), Valid: true},
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

	// For StageProgression, persist the new stage after successful application
	// The stageProg.CurrentStage was updated in-memory by Apply()
	if stageProg != nil {
		stateID := uuid.New().String()
		err = txQueries.UpsertUserProgressionState(ctx, db.UpsertUserProgressionStateParams{
			ID:            stateID,
			UserID:        event.UserID,
			LiftID:        liftID,
			ProgressionID: pp.ProgressionID,
			CurrentStage:  int64(stageProg.CurrentStage),
			StateData:     sql.NullString{}, // No additional state data needed
		})
		if err != nil {
			return TriggerResult{
				ProgressionID: pp.ProgressionID,
				LiftID:        liftID,
				Applied:       false,
				Error:         fmt.Sprintf("failed to persist stage state: %v", err),
			}
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

// getMaxTypeFromProgression extracts the MaxType from a progression.
func getMaxTypeFromProgression(prog progression.Progression) (progression.MaxType, error) {
	switch p := prog.(type) {
	case *progression.LinearProgression:
		return p.MaxTypeValue, nil
	case *progression.CycleProgression:
		return p.MaxTypeValue, nil
	case *progression.AMRAPProgression:
		return p.MaxTypeValue, nil
	case *progression.DeloadOnFailure:
		return p.MaxTypeValue, nil
	case *progression.StageProgression:
		return p.MaxTypeValue, nil
	default:
		return "", fmt.Errorf("unknown progression type: %T", prog)
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
	case progression.FailureTriggerContext:
		triggerEvent.ConsecutiveFailures = &ctx.ConsecutiveFailures
		triggerEvent.RepsPerformed = &ctx.RepsPerformed
		triggerEvent.TargetReps = &ctx.TargetReps
	case *progression.ManualTriggerContext:
		// For manual triggers, extract from the underlying context
		underlyingEvent := &progression.TriggerEventV2{
			Type:      ctx.InnerTriggerType,
			UserID:    event.UserID,
			Timestamp: event.Timestamp,
			Context:   ctx.UnderlyingContext,
		}
		return buildTriggerEvent(underlyingEvent)
	}

	return triggerEvent
}

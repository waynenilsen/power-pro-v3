// Package repository provides database repository implementations.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/db"
	"github.com/waynenilsen/power-pro-v3/internal/domain/workoutsession"
)

// WorkoutSessionRepository implements workout session persistence using sqlc-generated queries.
type WorkoutSessionRepository struct {
	queries *db.Queries
}

// NewWorkoutSessionRepository creates a new WorkoutSessionRepository.
func NewWorkoutSessionRepository(sqlDB *sql.DB) *WorkoutSessionRepository {
	return &WorkoutSessionRepository{
		queries: db.New(sqlDB),
	}
}

// GetByID retrieves a workout session by its ID.
func (r *WorkoutSessionRepository) GetByID(id string) (*workoutsession.WorkoutSession, error) {
	ctx := context.Background()
	dbSession, err := r.queries.GetWorkoutSessionByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get workout session by ID: %w", err)
	}
	return dbWorkoutSessionToDomain(dbSession), nil
}

// GetActiveByUserProgramStateID retrieves the active (IN_PROGRESS) session for a user program state.
func (r *WorkoutSessionRepository) GetActiveByUserProgramStateID(userProgramStateID string) (*workoutsession.WorkoutSession, error) {
	ctx := context.Background()
	dbSession, err := r.queries.GetActiveWorkoutSession(ctx, userProgramStateID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get active workout session: %w", err)
	}
	return dbWorkoutSessionToDomain(dbSession), nil
}

// GetByUserProgramStateID retrieves all workout sessions for a user program state, ordered by created_at DESC.
func (r *WorkoutSessionRepository) GetByUserProgramStateID(userProgramStateID string) ([]*workoutsession.WorkoutSession, error) {
	ctx := context.Background()
	dbSessions, err := r.queries.GetWorkoutSessionsByState(ctx, userProgramStateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workout sessions by state: %w", err)
	}

	sessions := make([]*workoutsession.WorkoutSession, len(dbSessions))
	for i, dbSession := range dbSessions {
		sessions[i] = dbWorkoutSessionToDomain(dbSession)
	}
	return sessions, nil
}

// Create persists a new workout session to the database.
func (r *WorkoutSessionRepository) Create(session *workoutsession.WorkoutSession) error {
	ctx := context.Background()

	err := r.queries.CreateWorkoutSession(ctx, db.CreateWorkoutSessionParams{
		ID:                 session.ID,
		UserProgramStateID: session.UserProgramStateID,
		WeekNumber:         int64(session.WeekNumber),
		DayIndex:           int64(session.DayIndex),
		Status:             string(session.Status),
		StartedAt:          session.StartedAt.Format(time.RFC3339),
		CreatedAt:          session.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          session.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to create workout session: %w", err)
	}
	return nil
}

// Complete marks a workout session as completed.
func (r *WorkoutSessionRepository) Complete(session *workoutsession.WorkoutSession) error {
	ctx := context.Background()

	err := r.queries.CompleteWorkoutSession(ctx, db.CompleteWorkoutSessionParams{
		FinishedAt: timePtrToNullStringWS(session.FinishedAt),
		UpdatedAt:  session.UpdatedAt.Format(time.RFC3339),
		ID:         session.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to complete workout session: %w", err)
	}
	return nil
}

// Abandon marks a workout session as abandoned.
func (r *WorkoutSessionRepository) Abandon(session *workoutsession.WorkoutSession) error {
	ctx := context.Background()

	err := r.queries.AbandonWorkoutSession(ctx, db.AbandonWorkoutSessionParams{
		FinishedAt: timePtrToNullStringWS(session.FinishedAt),
		UpdatedAt:  session.UpdatedAt.Format(time.RFC3339),
		ID:         session.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to abandon workout session: %w", err)
	}
	return nil
}

// UpdateStatus updates the status of a workout session.
func (r *WorkoutSessionRepository) UpdateStatus(session *workoutsession.WorkoutSession) error {
	ctx := context.Background()

	err := r.queries.UpdateWorkoutSessionStatus(ctx, db.UpdateWorkoutSessionStatusParams{
		Status:    string(session.Status),
		UpdatedAt: session.UpdatedAt.Format(time.RFC3339),
		ID:        session.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to update workout session status: %w", err)
	}
	return nil
}

// Delete removes a workout session from the database.
func (r *WorkoutSessionRepository) Delete(id string) error {
	ctx := context.Background()

	err := r.queries.DeleteWorkoutSession(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete workout session: %w", err)
	}
	return nil
}

// Helper functions

func dbWorkoutSessionToDomain(dbSession db.WorkoutSession) *workoutsession.WorkoutSession {
	startedAt, _ := time.Parse(time.RFC3339, dbSession.StartedAt)
	createdAt, _ := time.Parse(time.RFC3339, dbSession.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, dbSession.UpdatedAt)

	return &workoutsession.WorkoutSession{
		ID:                 dbSession.ID,
		UserProgramStateID: dbSession.UserProgramStateID,
		WeekNumber:         int(dbSession.WeekNumber),
		DayIndex:           int(dbSession.DayIndex),
		Status:             workoutsession.Status(dbSession.Status),
		StartedAt:          startedAt,
		FinishedAt:         nullStringToTimePtrWS(dbSession.FinishedAt),
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
	}
}

// nullStringToTimePtrWS converts a sql.NullString containing an RFC3339 time to a *time.Time.
func nullStringToTimePtrWS(ns sql.NullString) *time.Time {
	if !ns.Valid || ns.String == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, ns.String)
	if err != nil {
		return nil
	}
	return &t
}

// timePtrToNullStringWS converts a *time.Time to a sql.NullString in RFC3339 format.
func timePtrToNullStringWS(t *time.Time) sql.NullString {
	if t == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: t.Format(time.RFC3339), Valid: true}
}

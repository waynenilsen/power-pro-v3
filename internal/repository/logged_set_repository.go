// Package repository provides database repository implementations.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/db"
	"github.com/waynenilsen/power-pro-v3/internal/domain/loggedset"
)

// LoggedSetRepository implements persistence for LoggedSet entities using sqlc-generated queries.
type LoggedSetRepository struct {
	queries *db.Queries
}

// NewLoggedSetRepository creates a new LoggedSetRepository.
func NewLoggedSetRepository(sqlDB *sql.DB) *LoggedSetRepository {
	return &LoggedSetRepository{
		queries: db.New(sqlDB),
	}
}

// GetByID retrieves a logged set by its ID.
func (r *LoggedSetRepository) GetByID(id string) (*loggedset.LoggedSet, error) {
	ctx := context.Background()
	dbSet, err := r.queries.GetLoggedSet(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get logged set: %w", err)
	}
	return dbGetLoggedSetRowToDomain(dbSet), nil
}

// ListBySession retrieves all logged sets for a session.
func (r *LoggedSetRepository) ListBySession(sessionID string) ([]loggedset.LoggedSet, error) {
	ctx := context.Background()
	dbSets, err := r.queries.ListLoggedSetsBySession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list logged sets by session: %w", err)
	}

	sets := make([]loggedset.LoggedSet, len(dbSets))
	for i, dbSet := range dbSets {
		sets[i] = *dbListLoggedSetsBySessionRowToDomain(dbSet)
	}
	return sets, nil
}

// LoggedSetListParams contains parameters for listing logged sets by user.
type LoggedSetListParams struct {
	UserID string
	Limit  int64
	Offset int64
}

// ListByUser retrieves logged sets for a user with pagination.
func (r *LoggedSetRepository) ListByUser(params LoggedSetListParams) ([]loggedset.LoggedSet, int64, error) {
	ctx := context.Background()

	if params.Limit <= 0 {
		params.Limit = 20
	}

	total, err := r.queries.CountLoggedSetsByUser(ctx, params.UserID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count logged sets: %w", err)
	}

	dbSets, err := r.queries.ListLoggedSetsByUser(ctx, db.ListLoggedSetsByUserParams{
		UserID: params.UserID,
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list logged sets by user: %w", err)
	}

	sets := make([]loggedset.LoggedSet, len(dbSets))
	for i, dbSet := range dbSets {
		sets[i] = *dbListLoggedSetsByUserRowToDomain(dbSet)
	}
	return sets, total, nil
}

// GetLatestAMRAPForLift retrieves the most recent AMRAP set for a user's lift.
func (r *LoggedSetRepository) GetLatestAMRAPForLift(userID, liftID string) (*loggedset.LoggedSet, error) {
	ctx := context.Background()
	dbSet, err := r.queries.GetLatestAMRAPForLift(ctx, db.GetLatestAMRAPForLiftParams{
		UserID: userID,
		LiftID: liftID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get latest AMRAP for lift: %w", err)
	}
	return dbGetLatestAMRAPForLiftRowToDomain(dbSet), nil
}

// Create persists a new logged set to the database.
func (r *LoggedSetRepository) Create(ls *loggedset.LoggedSet) error {
	ctx := context.Background()

	var rpe sql.NullFloat64
	if ls.RPE != nil {
		rpe = sql.NullFloat64{Float64: *ls.RPE, Valid: true}
	}

	err := r.queries.CreateLoggedSet(ctx, db.CreateLoggedSetParams{
		ID:             ls.ID,
		UserID:         ls.UserID,
		SessionID:      ls.SessionID,
		PrescriptionID: ls.PrescriptionID,
		LiftID:         ls.LiftID,
		SetNumber:      int64(ls.SetNumber),
		Weight:         ls.Weight,
		TargetReps:     int64(ls.TargetReps),
		RepsPerformed:  int64(ls.RepsPerformed),
		IsAmrap:        ls.IsAMRAP,
		Rpe:            rpe,
		CreatedAt:      ls.CreatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to create logged set: %w", err)
	}
	return nil
}

// Delete removes a logged set from the database.
func (r *LoggedSetRepository) Delete(id string) error {
	ctx := context.Background()

	err := r.queries.DeleteLoggedSet(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete logged set: %w", err)
	}
	return nil
}

// DeleteBySession removes all logged sets for a session.
func (r *LoggedSetRepository) DeleteBySession(sessionID string) error {
	ctx := context.Background()

	err := r.queries.DeleteLoggedSetsBySession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete logged sets by session: %w", err)
	}
	return nil
}

// ListBySessionAndPrescription retrieves all logged sets for a session and prescription.
func (r *LoggedSetRepository) ListBySessionAndPrescription(sessionID, prescriptionID string) ([]loggedset.LoggedSet, error) {
	ctx := context.Background()
	dbSets, err := r.queries.ListLoggedSetsBySessionAndPrescription(ctx, db.ListLoggedSetsBySessionAndPrescriptionParams{
		SessionID:      sessionID,
		PrescriptionID: prescriptionID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list logged sets by session and prescription: %w", err)
	}

	sets := make([]loggedset.LoggedSet, len(dbSets))
	for i, dbSet := range dbSets {
		sets[i] = *dbListLoggedSetsBySessionAndPrescriptionRowToDomain(dbSet)
	}
	return sets, nil
}

// Helper functions

// nullFloat64ToPtr converts a sql.NullFloat64 to a *float64.
func nullFloat64ToPtr(nf sql.NullFloat64) *float64 {
	if nf.Valid {
		return &nf.Float64
	}
	return nil
}

func dbGetLoggedSetRowToDomain(dbSet db.GetLoggedSetRow) *loggedset.LoggedSet {
	createdAt, _ := time.Parse(time.RFC3339, dbSet.CreatedAt)

	return &loggedset.LoggedSet{
		ID:             dbSet.ID,
		UserID:         dbSet.UserID,
		SessionID:      dbSet.SessionID,
		PrescriptionID: dbSet.PrescriptionID,
		LiftID:         dbSet.LiftID,
		SetNumber:      int(dbSet.SetNumber),
		Weight:         dbSet.Weight,
		TargetReps:     int(dbSet.TargetReps),
		RepsPerformed:  int(dbSet.RepsPerformed),
		IsAMRAP:        dbSet.IsAmrap,
		RPE:            nullFloat64ToPtr(dbSet.Rpe),
		CreatedAt:      createdAt,
	}
}

func dbListLoggedSetsBySessionRowToDomain(dbSet db.ListLoggedSetsBySessionRow) *loggedset.LoggedSet {
	createdAt, _ := time.Parse(time.RFC3339, dbSet.CreatedAt)

	return &loggedset.LoggedSet{
		ID:             dbSet.ID,
		UserID:         dbSet.UserID,
		SessionID:      dbSet.SessionID,
		PrescriptionID: dbSet.PrescriptionID,
		LiftID:         dbSet.LiftID,
		SetNumber:      int(dbSet.SetNumber),
		Weight:         dbSet.Weight,
		TargetReps:     int(dbSet.TargetReps),
		RepsPerformed:  int(dbSet.RepsPerformed),
		IsAMRAP:        dbSet.IsAmrap,
		RPE:            nullFloat64ToPtr(dbSet.Rpe),
		CreatedAt:      createdAt,
	}
}

func dbListLoggedSetsByUserRowToDomain(dbSet db.ListLoggedSetsByUserRow) *loggedset.LoggedSet {
	createdAt, _ := time.Parse(time.RFC3339, dbSet.CreatedAt)

	return &loggedset.LoggedSet{
		ID:             dbSet.ID,
		UserID:         dbSet.UserID,
		SessionID:      dbSet.SessionID,
		PrescriptionID: dbSet.PrescriptionID,
		LiftID:         dbSet.LiftID,
		SetNumber:      int(dbSet.SetNumber),
		Weight:         dbSet.Weight,
		TargetReps:     int(dbSet.TargetReps),
		RepsPerformed:  int(dbSet.RepsPerformed),
		IsAMRAP:        dbSet.IsAmrap,
		RPE:            nullFloat64ToPtr(dbSet.Rpe),
		CreatedAt:      createdAt,
	}
}

func dbGetLatestAMRAPForLiftRowToDomain(dbSet db.GetLatestAMRAPForLiftRow) *loggedset.LoggedSet {
	createdAt, _ := time.Parse(time.RFC3339, dbSet.CreatedAt)

	return &loggedset.LoggedSet{
		ID:             dbSet.ID,
		UserID:         dbSet.UserID,
		SessionID:      dbSet.SessionID,
		PrescriptionID: dbSet.PrescriptionID,
		LiftID:         dbSet.LiftID,
		SetNumber:      int(dbSet.SetNumber),
		Weight:         dbSet.Weight,
		TargetReps:     int(dbSet.TargetReps),
		RepsPerformed:  int(dbSet.RepsPerformed),
		IsAMRAP:        dbSet.IsAmrap,
		RPE:            nullFloat64ToPtr(dbSet.Rpe),
		CreatedAt:      createdAt,
	}
}

func dbListLoggedSetsBySessionAndPrescriptionRowToDomain(dbSet db.ListLoggedSetsBySessionAndPrescriptionRow) *loggedset.LoggedSet {
	createdAt, _ := time.Parse(time.RFC3339, dbSet.CreatedAt)

	return &loggedset.LoggedSet{
		ID:             dbSet.ID,
		UserID:         dbSet.UserID,
		SessionID:      dbSet.SessionID,
		PrescriptionID: dbSet.PrescriptionID,
		LiftID:         dbSet.LiftID,
		SetNumber:      int(dbSet.SetNumber),
		Weight:         dbSet.Weight,
		TargetReps:     int(dbSet.TargetReps),
		RepsPerformed:  int(dbSet.RepsPerformed),
		IsAMRAP:        dbSet.IsAmrap,
		RPE:            nullFloat64ToPtr(dbSet.Rpe),
		CreatedAt:      createdAt,
	}
}

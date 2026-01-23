// Package repository provides database repository implementations.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/db"
	"github.com/waynenilsen/power-pro-v3/internal/domain/progression"
)

// FailureCounterRepository implements persistence for FailureCounter entities using sqlc-generated queries.
type FailureCounterRepository struct {
	queries *db.Queries
}

// NewFailureCounterRepository creates a new FailureCounterRepository.
func NewFailureCounterRepository(sqlDB *sql.DB) *FailureCounterRepository {
	return &FailureCounterRepository{
		queries: db.New(sqlDB),
	}
}

// GetByID retrieves a failure counter by its ID.
func (r *FailureCounterRepository) GetByID(id string) (*progression.FailureCounter, error) {
	ctx := context.Background()
	dbCounter, err := r.queries.GetFailureCounter(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get failure counter: %w", err)
	}
	return dbFailureCounterToDomain(dbCounter), nil
}

// GetByKey retrieves a failure counter by its composite key (user_id, lift_id, progression_id).
func (r *FailureCounterRepository) GetByKey(userID, liftID, progressionID string) (*progression.FailureCounter, error) {
	ctx := context.Background()
	dbCounter, err := r.queries.GetFailureCounterByKey(ctx, db.GetFailureCounterByKeyParams{
		UserID:        userID,
		LiftID:        liftID,
		ProgressionID: progressionID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get failure counter by key: %w", err)
	}
	return dbFailureCounterToDomain(dbCounter), nil
}

// ListByUser retrieves all failure counters for a user.
func (r *FailureCounterRepository) ListByUser(userID string) ([]progression.FailureCounter, error) {
	ctx := context.Background()
	dbCounters, err := r.queries.ListFailureCountersByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list failure counters by user: %w", err)
	}

	counters := make([]progression.FailureCounter, len(dbCounters))
	for i, dbCounter := range dbCounters {
		counters[i] = *dbFailureCounterToDomain(dbCounter)
	}
	return counters, nil
}

// ListByUserAndLift retrieves all failure counters for a user's lift.
func (r *FailureCounterRepository) ListByUserAndLift(userID, liftID string) ([]progression.FailureCounter, error) {
	ctx := context.Background()
	dbCounters, err := r.queries.ListFailureCountersByUserAndLift(ctx, db.ListFailureCountersByUserAndLiftParams{
		UserID: userID,
		LiftID: liftID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list failure counters by user and lift: %w", err)
	}

	counters := make([]progression.FailureCounter, len(dbCounters))
	for i, dbCounter := range dbCounters {
		counters[i] = *dbFailureCounterToDomain(dbCounter)
	}
	return counters, nil
}

// ListByProgression retrieves all failure counters for a progression.
func (r *FailureCounterRepository) ListByProgression(progressionID string) ([]progression.FailureCounter, error) {
	ctx := context.Background()
	dbCounters, err := r.queries.ListFailureCountersByProgression(ctx, progressionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list failure counters by progression: %w", err)
	}

	counters := make([]progression.FailureCounter, len(dbCounters))
	for i, dbCounter := range dbCounters {
		counters[i] = *dbFailureCounterToDomain(dbCounter)
	}
	return counters, nil
}

// Create persists a new failure counter to the database.
func (r *FailureCounterRepository) Create(fc *progression.FailureCounter) error {
	ctx := context.Background()

	var lastFailureAt, lastSuccessAt sql.NullString
	if fc.LastFailureAt != nil {
		lastFailureAt = sql.NullString{String: fc.LastFailureAt.Format(time.RFC3339), Valid: true}
	}
	if fc.LastSuccessAt != nil {
		lastSuccessAt = sql.NullString{String: fc.LastSuccessAt.Format(time.RFC3339), Valid: true}
	}

	err := r.queries.CreateFailureCounter(ctx, db.CreateFailureCounterParams{
		ID:                  fc.ID,
		UserID:              fc.UserID,
		LiftID:              fc.LiftID,
		ProgressionID:       fc.ProgressionID,
		ConsecutiveFailures: int64(fc.ConsecutiveFailures),
		LastFailureAt:       lastFailureAt,
		LastSuccessAt:       lastSuccessAt,
		CreatedAt:           fc.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           fc.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to create failure counter: %w", err)
	}
	return nil
}

// Update updates an existing failure counter in the database.
func (r *FailureCounterRepository) Update(fc *progression.FailureCounter) error {
	ctx := context.Background()

	var lastFailureAt, lastSuccessAt sql.NullString
	if fc.LastFailureAt != nil {
		lastFailureAt = sql.NullString{String: fc.LastFailureAt.Format(time.RFC3339), Valid: true}
	}
	if fc.LastSuccessAt != nil {
		lastSuccessAt = sql.NullString{String: fc.LastSuccessAt.Format(time.RFC3339), Valid: true}
	}

	err := r.queries.UpdateFailureCounter(ctx, db.UpdateFailureCounterParams{
		ID:                  fc.ID,
		ConsecutiveFailures: int64(fc.ConsecutiveFailures),
		LastFailureAt:       lastFailureAt,
		LastSuccessAt:       lastSuccessAt,
		UpdatedAt:           fc.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to update failure counter: %w", err)
	}
	return nil
}

// IncrementOnFailure atomically increments the failure counter using upsert.
// Creates the counter if it doesn't exist.
// Returns the current failure count after increment.
func (r *FailureCounterRepository) IncrementOnFailure(userID, liftID, progressionID string) (int, error) {
	ctx := context.Background()
	now := time.Now().Format(time.RFC3339)

	// Use upsert to atomically increment or create
	err := r.queries.UpsertFailureCounterOnFailure(ctx, db.UpsertFailureCounterOnFailureParams{
		ID:            uuid.New().String(),
		UserID:        userID,
		LiftID:        liftID,
		ProgressionID: progressionID,
		LastFailureAt: sql.NullString{String: now, Valid: true},
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to increment failure counter: %w", err)
	}

	// Fetch the updated counter to get the current count
	counter, err := r.GetByKey(userID, liftID, progressionID)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch updated failure counter: %w", err)
	}
	if counter == nil {
		return 0, fmt.Errorf("failure counter not found after upsert")
	}

	return counter.ConsecutiveFailures, nil
}

// ResetOnSuccess atomically resets the failure counter using upsert.
// Creates the counter if it doesn't exist (with 0 failures).
func (r *FailureCounterRepository) ResetOnSuccess(userID, liftID, progressionID string) error {
	ctx := context.Background()
	now := time.Now().Format(time.RFC3339)

	err := r.queries.UpsertFailureCounterOnSuccess(ctx, db.UpsertFailureCounterOnSuccessParams{
		ID:            uuid.New().String(),
		UserID:        userID,
		LiftID:        liftID,
		ProgressionID: progressionID,
		LastSuccessAt: sql.NullString{String: now, Valid: true},
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	if err != nil {
		return fmt.Errorf("failed to reset failure counter: %w", err)
	}
	return nil
}

// GetFailureCount retrieves the current consecutive failure count.
// Returns 0 if no counter exists.
func (r *FailureCounterRepository) GetFailureCount(userID, liftID, progressionID string) (int, error) {
	counter, err := r.GetByKey(userID, liftID, progressionID)
	if err != nil {
		return 0, err
	}
	if counter == nil {
		return 0, nil
	}
	return counter.ConsecutiveFailures, nil
}

// Delete removes a failure counter from the database by ID.
func (r *FailureCounterRepository) Delete(id string) error {
	ctx := context.Background()

	err := r.queries.DeleteFailureCounter(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete failure counter: %w", err)
	}
	return nil
}

// DeleteByKey removes a failure counter from the database by composite key.
func (r *FailureCounterRepository) DeleteByKey(userID, liftID, progressionID string) error {
	ctx := context.Background()

	err := r.queries.DeleteFailureCounterByKey(ctx, db.DeleteFailureCounterByKeyParams{
		UserID:        userID,
		LiftID:        liftID,
		ProgressionID: progressionID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete failure counter by key: %w", err)
	}
	return nil
}

// Helper functions

func dbFailureCounterToDomain(dbCounter db.FailureCounter) *progression.FailureCounter {
	createdAt, _ := time.Parse(time.RFC3339, dbCounter.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, dbCounter.UpdatedAt)

	var lastFailureAt, lastSuccessAt *time.Time
	if dbCounter.LastFailureAt.Valid {
		t, _ := time.Parse(time.RFC3339, dbCounter.LastFailureAt.String)
		lastFailureAt = &t
	}
	if dbCounter.LastSuccessAt.Valid {
		t, _ := time.Parse(time.RFC3339, dbCounter.LastSuccessAt.String)
		lastSuccessAt = &t
	}

	return &progression.FailureCounter{
		ID:                  dbCounter.ID,
		UserID:              dbCounter.UserID,
		LiftID:              dbCounter.LiftID,
		ProgressionID:       dbCounter.ProgressionID,
		ConsecutiveFailures: int(dbCounter.ConsecutiveFailures),
		LastFailureAt:       lastFailureAt,
		LastSuccessAt:       lastSuccessAt,
		CreatedAt:           createdAt,
		UpdatedAt:           updatedAt,
	}
}

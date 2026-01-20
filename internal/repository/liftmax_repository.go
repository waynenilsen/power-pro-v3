// Package repository provides database repository implementations.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/db"
	"github.com/waynenilsen/power-pro-v3/internal/domain/liftmax"
)

// LiftMaxRepository implements persistence operations for LiftMax entities.
type LiftMaxRepository struct {
	queries *db.Queries
}

// NewLiftMaxRepository creates a new LiftMaxRepository.
func NewLiftMaxRepository(sqlDB *sql.DB) *LiftMaxRepository {
	return &LiftMaxRepository{
		queries: db.New(sqlDB),
	}
}

// LiftMaxListParams contains parameters for listing lift maxes.
type LiftMaxListParams struct {
	UserID    string
	LiftID    *string
	Type      *string
	SortOrder SortOrder
	Limit     int64
	Offset    int64
}

// GetByID retrieves a lift max by its ID.
func (r *LiftMaxRepository) GetByID(id string) (*liftmax.LiftMax, error) {
	ctx := context.Background()
	dbMax, err := r.queries.GetLiftMax(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get lift max: %w", err)
	}
	return dbLiftMaxToDomain(dbMax), nil
}

// GetCurrentOneRM retrieves the most recent 1RM for a user and lift.
// Implements liftmax.LiftMaxRepository interface.
func (r *LiftMaxRepository) GetCurrentOneRM(userID, liftID string) (*liftmax.LiftMax, error) {
	ctx := context.Background()
	dbMax, err := r.queries.GetCurrentOneRM(ctx, db.GetCurrentOneRMParams{
		UserID: userID,
		LiftID: liftID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get current 1RM: %w", err)
	}
	return dbLiftMaxToDomain(dbMax), nil
}

// GetCurrentMax retrieves the most recent max for a user, lift, and type.
func (r *LiftMaxRepository) GetCurrentMax(userID, liftID, maxType string) (*liftmax.LiftMax, error) {
	ctx := context.Background()
	dbMax, err := r.queries.GetCurrentMax(ctx, db.GetCurrentMaxParams{
		UserID: userID,
		LiftID: liftID,
		Type:   maxType,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get current max: %w", err)
	}
	return dbLiftMaxToDomain(dbMax), nil
}

// List retrieves lift maxes with pagination and optional filtering.
func (r *LiftMaxRepository) List(params LiftMaxListParams) ([]liftmax.LiftMax, int64, error) {
	ctx := context.Background()

	// Set defaults
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.SortOrder == "" {
		params.SortOrder = SortDesc // Default to descending by effective_date
	}

	var dbMaxes []db.LiftMax
	var total int64
	var err error

	// Determine which query variant to use based on filters
	hasLiftFilter := params.LiftID != nil
	hasTypeFilter := params.Type != nil

	// Get count
	switch {
	case hasLiftFilter && hasTypeFilter:
		total, err = r.queries.CountLiftMaxesByUserFilterLiftAndType(ctx, db.CountLiftMaxesByUserFilterLiftAndTypeParams{
			UserID: params.UserID,
			LiftID: *params.LiftID,
			Type:   *params.Type,
		})
	case hasLiftFilter:
		total, err = r.queries.CountLiftMaxesByUserFilterLift(ctx, db.CountLiftMaxesByUserFilterLiftParams{
			UserID: params.UserID,
			LiftID: *params.LiftID,
		})
	case hasTypeFilter:
		total, err = r.queries.CountLiftMaxesByUserFilterType(ctx, db.CountLiftMaxesByUserFilterTypeParams{
			UserID: params.UserID,
			Type:   *params.Type,
		})
	default:
		total, err = r.queries.CountLiftMaxesByUser(ctx, params.UserID)
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to count lift maxes: %w", err)
	}

	// Get data based on filters and sort order
	switch {
	case hasLiftFilter && hasTypeFilter && params.SortOrder == SortDesc:
		dbMaxes, err = r.queries.ListLiftMaxesByUserFilterLiftAndTypeByEffectiveDateDesc(ctx, db.ListLiftMaxesByUserFilterLiftAndTypeByEffectiveDateDescParams{
			UserID: params.UserID,
			LiftID: *params.LiftID,
			Type:   *params.Type,
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	case hasLiftFilter && hasTypeFilter && params.SortOrder == SortAsc:
		dbMaxes, err = r.queries.ListLiftMaxesByUserFilterLiftAndTypeByEffectiveDateAsc(ctx, db.ListLiftMaxesByUserFilterLiftAndTypeByEffectiveDateAscParams{
			UserID: params.UserID,
			LiftID: *params.LiftID,
			Type:   *params.Type,
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	case hasLiftFilter && params.SortOrder == SortDesc:
		dbMaxes, err = r.queries.ListLiftMaxesByUserFilterLiftByEffectiveDateDesc(ctx, db.ListLiftMaxesByUserFilterLiftByEffectiveDateDescParams{
			UserID: params.UserID,
			LiftID: *params.LiftID,
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	case hasLiftFilter && params.SortOrder == SortAsc:
		dbMaxes, err = r.queries.ListLiftMaxesByUserFilterLiftByEffectiveDateAsc(ctx, db.ListLiftMaxesByUserFilterLiftByEffectiveDateAscParams{
			UserID: params.UserID,
			LiftID: *params.LiftID,
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	case hasTypeFilter && params.SortOrder == SortDesc:
		dbMaxes, err = r.queries.ListLiftMaxesByUserFilterTypeByEffectiveDateDesc(ctx, db.ListLiftMaxesByUserFilterTypeByEffectiveDateDescParams{
			UserID: params.UserID,
			Type:   *params.Type,
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	case hasTypeFilter && params.SortOrder == SortAsc:
		dbMaxes, err = r.queries.ListLiftMaxesByUserFilterTypeByEffectiveDateAsc(ctx, db.ListLiftMaxesByUserFilterTypeByEffectiveDateAscParams{
			UserID: params.UserID,
			Type:   *params.Type,
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	case params.SortOrder == SortAsc:
		dbMaxes, err = r.queries.ListLiftMaxesByUserByEffectiveDateAsc(ctx, db.ListLiftMaxesByUserByEffectiveDateAscParams{
			UserID: params.UserID,
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	default:
		dbMaxes, err = r.queries.ListLiftMaxesByUserByEffectiveDateDesc(ctx, db.ListLiftMaxesByUserByEffectiveDateDescParams{
			UserID: params.UserID,
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list lift maxes: %w", err)
	}

	maxes := make([]liftmax.LiftMax, len(dbMaxes))
	for i, dbMax := range dbMaxes {
		maxes[i] = *dbLiftMaxToDomain(dbMax)
	}

	return maxes, total, nil
}

// Create persists a new lift max to the database.
func (r *LiftMaxRepository) Create(m *liftmax.LiftMax) error {
	ctx := context.Background()

	err := r.queries.CreateLiftMax(ctx, db.CreateLiftMaxParams{
		ID:            m.ID,
		UserID:        m.UserID,
		LiftID:        m.LiftID,
		Type:          string(m.Type),
		Value:         m.Value,
		EffectiveDate: m.EffectiveDate.Format(time.RFC3339),
		CreatedAt:     m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     m.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to create lift max: %w", err)
	}
	return nil
}

// Update persists changes to an existing lift max.
func (r *LiftMaxRepository) Update(m *liftmax.LiftMax) error {
	ctx := context.Background()

	err := r.queries.UpdateLiftMax(ctx, db.UpdateLiftMaxParams{
		ID:            m.ID,
		Value:         m.Value,
		EffectiveDate: m.EffectiveDate.Format(time.RFC3339),
		UpdatedAt:     m.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to update lift max: %w", err)
	}
	return nil
}

// Delete removes a lift max from the database.
func (r *LiftMaxRepository) Delete(id string) error {
	ctx := context.Background()

	err := r.queries.DeleteLiftMax(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete lift max: %w", err)
	}
	return nil
}

// UniqueConstraintExists checks if a lift max with the same (user, lift, type, effective_date) exists.
func (r *LiftMaxRepository) UniqueConstraintExists(userID, liftID, maxType string, effectiveDate time.Time, excludeID *string) (bool, error) {
	ctx := context.Background()
	effectiveDateStr := effectiveDate.Format(time.RFC3339)

	if excludeID == nil {
		exists, err := r.queries.UniqueConstraintExists(ctx, db.UniqueConstraintExistsParams{
			UserID:        userID,
			LiftID:        liftID,
			Type:          maxType,
			EffectiveDate: effectiveDateStr,
		})
		if err != nil {
			return false, fmt.Errorf("failed to check unique constraint: %w", err)
		}
		return exists == 1, nil
	}

	exists, err := r.queries.UniqueConstraintExistsExcluding(ctx, db.UniqueConstraintExistsExcludingParams{
		UserID:        userID,
		LiftID:        liftID,
		Type:          maxType,
		EffectiveDate: effectiveDateStr,
		ID:            *excludeID,
	})
	if err != nil {
		return false, fmt.Errorf("failed to check unique constraint: %w", err)
	}
	return exists == 1, nil
}

// Helper function to convert database model to domain model.
func dbLiftMaxToDomain(dbMax db.LiftMax) *liftmax.LiftMax {
	effectiveDate, _ := time.Parse(time.RFC3339, dbMax.EffectiveDate)
	createdAt, _ := time.Parse(time.RFC3339, dbMax.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, dbMax.UpdatedAt)

	return &liftmax.LiftMax{
		ID:            dbMax.ID,
		UserID:        dbMax.UserID,
		LiftID:        dbMax.LiftID,
		Type:          liftmax.MaxType(dbMax.Type),
		Value:         dbMax.Value,
		EffectiveDate: effectiveDate,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

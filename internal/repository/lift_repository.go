// Package repository provides database repository implementations.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/db"
	"github.com/waynenilsen/power-pro-v3/internal/domain/lift"
)

// LiftRepository implements lift.LiftRepository using sqlc-generated queries.
type LiftRepository struct {
	queries *db.Queries
}

// NewLiftRepository creates a new LiftRepository.
func NewLiftRepository(sqlDB *sql.DB) *LiftRepository {
	return &LiftRepository{
		queries: db.New(sqlDB),
	}
}

// GetByID retrieves a lift by its ID.
func (r *LiftRepository) GetByID(id string) (*lift.Lift, error) {
	ctx := context.Background()
	dbLift, err := r.queries.GetLift(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get lift: %w", err)
	}
	return dbLiftToDomain(dbLift), nil
}

// GetBySlug retrieves a lift by its slug.
func (r *LiftRepository) GetBySlug(slug string) (*lift.Lift, error) {
	ctx := context.Background()
	dbLift, err := r.queries.GetLiftBySlug(ctx, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get lift by slug: %w", err)
	}
	return dbLiftToDomain(dbLift), nil
}

// SortField represents a field to sort by.
type SortField string

const (
	SortByName      SortField = "name"
	SortByCreatedAt SortField = "created_at"
)

// SortOrder represents the sort direction.
type SortOrder string

const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

// ListParams contains parameters for listing lifts.
type ListParams struct {
	Limit             int64
	Offset            int64
	SortBy            SortField
	SortOrder         SortOrder
	FilterCompetition *bool
}

// List retrieves lifts with pagination, sorting, and optional filtering.
func (r *LiftRepository) List(params ListParams) ([]lift.Lift, int64, error) {
	ctx := context.Background()

	// Set defaults
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = SortByName
	}
	if params.SortOrder == "" {
		params.SortOrder = SortAsc
	}

	var dbLifts []db.Lift
	var total int64
	var err error

	if params.FilterCompetition != nil {
		// Filter by competition lift
		isCompetition := int64(0)
		if *params.FilterCompetition {
			isCompetition = 1
		}

		total, err = r.queries.CountLiftsFilteredByCompetition(ctx, isCompetition)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count lifts: %w", err)
		}

		// Select appropriate query based on sort
		switch {
		case params.SortBy == SortByName && params.SortOrder == SortAsc:
			dbLifts, err = r.queries.ListLiftsFilteredByCompetitionByNameAsc(ctx, db.ListLiftsFilteredByCompetitionByNameAscParams{
				IsCompetitionLift: isCompetition,
				Limit:             params.Limit,
				Offset:            params.Offset,
			})
		case params.SortBy == SortByName && params.SortOrder == SortDesc:
			dbLifts, err = r.queries.ListLiftsFilteredByCompetitionByNameDesc(ctx, db.ListLiftsFilteredByCompetitionByNameDescParams{
				IsCompetitionLift: isCompetition,
				Limit:             params.Limit,
				Offset:            params.Offset,
			})
		case params.SortBy == SortByCreatedAt && params.SortOrder == SortAsc:
			dbLifts, err = r.queries.ListLiftsFilteredByCompetitionByCreatedAtAsc(ctx, db.ListLiftsFilteredByCompetitionByCreatedAtAscParams{
				IsCompetitionLift: isCompetition,
				Limit:             params.Limit,
				Offset:            params.Offset,
			})
		case params.SortBy == SortByCreatedAt && params.SortOrder == SortDesc:
			dbLifts, err = r.queries.ListLiftsFilteredByCompetitionByCreatedAtDesc(ctx, db.ListLiftsFilteredByCompetitionByCreatedAtDescParams{
				IsCompetitionLift: isCompetition,
				Limit:             params.Limit,
				Offset:            params.Offset,
			})
		default:
			dbLifts, err = r.queries.ListLiftsFilteredByCompetitionByNameAsc(ctx, db.ListLiftsFilteredByCompetitionByNameAscParams{
				IsCompetitionLift: isCompetition,
				Limit:             params.Limit,
				Offset:            params.Offset,
			})
		}
	} else {
		// No filter
		total, err = r.queries.CountLifts(ctx)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count lifts: %w", err)
		}

		// Select appropriate query based on sort
		switch {
		case params.SortBy == SortByName && params.SortOrder == SortAsc:
			dbLifts, err = r.queries.ListLiftsByNameAsc(ctx, db.ListLiftsByNameAscParams{
				Limit:  params.Limit,
				Offset: params.Offset,
			})
		case params.SortBy == SortByName && params.SortOrder == SortDesc:
			dbLifts, err = r.queries.ListLiftsByNameDesc(ctx, db.ListLiftsByNameDescParams{
				Limit:  params.Limit,
				Offset: params.Offset,
			})
		case params.SortBy == SortByCreatedAt && params.SortOrder == SortAsc:
			dbLifts, err = r.queries.ListLiftsByCreatedAtAsc(ctx, db.ListLiftsByCreatedAtAscParams{
				Limit:  params.Limit,
				Offset: params.Offset,
			})
		case params.SortBy == SortByCreatedAt && params.SortOrder == SortDesc:
			dbLifts, err = r.queries.ListLiftsByCreatedAtDesc(ctx, db.ListLiftsByCreatedAtDescParams{
				Limit:  params.Limit,
				Offset: params.Offset,
			})
		default:
			dbLifts, err = r.queries.ListLiftsByNameAsc(ctx, db.ListLiftsByNameAscParams{
				Limit:  params.Limit,
				Offset: params.Offset,
			})
		}
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list lifts: %w", err)
	}

	lifts := make([]lift.Lift, len(dbLifts))
	for i, dbLift := range dbLifts {
		lifts[i] = *dbLiftToDomain(dbLift)
	}

	return lifts, total, nil
}

// SlugExists checks if a slug already exists, excluding a specific ID.
func (r *LiftRepository) SlugExists(slug string, excludeID *string) (bool, error) {
	ctx := context.Background()

	if excludeID == nil {
		exists, err := r.queries.SlugExistsForNew(ctx, slug)
		if err != nil {
			return false, fmt.Errorf("failed to check slug exists: %w", err)
		}
		return exists == 1, nil
	}

	exists, err := r.queries.SlugExists(ctx, db.SlugExistsParams{
		Slug: slug,
		ID:   *excludeID,
	})
	if err != nil {
		return false, fmt.Errorf("failed to check slug exists: %w", err)
	}
	return exists == 1, nil
}

// Create persists a new lift to the database.
func (r *LiftRepository) Create(l *lift.Lift) error {
	ctx := context.Background()

	err := r.queries.CreateLift(ctx, db.CreateLiftParams{
		ID:                l.ID,
		Name:              l.Name,
		Slug:              l.Slug,
		IsCompetitionLift: boolToInt64(l.IsCompetitionLift),
		ParentLiftID:      stringPtrToNullString(l.ParentLiftID),
		CreatedAt:         l.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         l.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to create lift: %w", err)
	}
	return nil
}

// Update persists changes to an existing lift.
func (r *LiftRepository) Update(l *lift.Lift) error {
	ctx := context.Background()

	err := r.queries.UpdateLift(ctx, db.UpdateLiftParams{
		ID:                l.ID,
		Name:              l.Name,
		Slug:              l.Slug,
		IsCompetitionLift: boolToInt64(l.IsCompetitionLift),
		ParentLiftID:      stringPtrToNullString(l.ParentLiftID),
		UpdatedAt:         l.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to update lift: %w", err)
	}
	return nil
}

// Delete removes a lift from the database.
func (r *LiftRepository) Delete(id string) error {
	ctx := context.Background()

	err := r.queries.DeleteLift(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete lift: %w", err)
	}
	return nil
}

// HasChildReferences checks if a lift is referenced as a parent by other lifts.
func (r *LiftRepository) HasChildReferences(id string) (bool, error) {
	ctx := context.Background()

	hasRefs, err := r.queries.LiftHasChildReferences(ctx, sql.NullString{String: id, Valid: true})
	if err != nil {
		return false, fmt.Errorf("failed to check references: %w", err)
	}
	return hasRefs == 1, nil
}

// Helper functions

func dbLiftToDomain(dbLift db.Lift) *lift.Lift {
	createdAt, _ := time.Parse(time.RFC3339, dbLift.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, dbLift.UpdatedAt)

	return &lift.Lift{
		ID:                dbLift.ID,
		Name:              dbLift.Name,
		Slug:              dbLift.Slug,
		IsCompetitionLift: dbLift.IsCompetitionLift == 1,
		ParentLiftID:      nullStringToStringPtr(dbLift.ParentLiftID),
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}
}

func boolToInt64(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

func stringPtrToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func nullStringToStringPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

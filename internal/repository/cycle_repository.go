// Package repository provides database repository implementations.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/db"
	"github.com/waynenilsen/power-pro-v3/internal/domain/cycle"
)

// CycleRepository implements cycle persistence using sqlc-generated queries.
type CycleRepository struct {
	queries *db.Queries
}

// NewCycleRepository creates a new CycleRepository.
func NewCycleRepository(sqlDB *sql.DB) *CycleRepository {
	return &CycleRepository{
		queries: db.New(sqlDB),
	}
}

// CycleSortField represents a field to sort by.
type CycleSortField string

const (
	CycleSortByName        CycleSortField = "name"
	CycleSortByCreatedAt   CycleSortField = "created_at"
	CycleSortByLengthWeeks CycleSortField = "length_weeks"
)

// CycleListParams contains parameters for listing cycles.
type CycleListParams struct {
	Limit     int64
	Offset    int64
	SortBy    CycleSortField
	SortOrder SortOrder
}

// GetByID retrieves a cycle by its ID.
func (r *CycleRepository) GetByID(id string) (*cycle.Cycle, error) {
	ctx := context.Background()
	dbCycle, err := r.queries.GetCycle(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get cycle: %w", err)
	}
	return dbCycleToDomain(dbCycle), nil
}

// List retrieves cycles with pagination and sorting.
func (r *CycleRepository) List(params CycleListParams) ([]cycle.Cycle, int64, error) {
	ctx := context.Background()

	// Set defaults
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = CycleSortByName
	}
	if params.SortOrder == "" {
		params.SortOrder = SortAsc
	}

	total, err := r.queries.CountCycles(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count cycles: %w", err)
	}

	var dbCycles []db.Cycle

	// Select appropriate query based on sort
	switch {
	case params.SortBy == CycleSortByName && params.SortOrder == SortAsc:
		dbCycles, err = r.queries.ListCyclesByNameAsc(ctx, db.ListCyclesByNameAscParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	case params.SortBy == CycleSortByName && params.SortOrder == SortDesc:
		dbCycles, err = r.queries.ListCyclesByNameDesc(ctx, db.ListCyclesByNameDescParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	case params.SortBy == CycleSortByCreatedAt && params.SortOrder == SortAsc:
		dbCycles, err = r.queries.ListCyclesByCreatedAtAsc(ctx, db.ListCyclesByCreatedAtAscParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	case params.SortBy == CycleSortByCreatedAt && params.SortOrder == SortDesc:
		dbCycles, err = r.queries.ListCyclesByCreatedAtDesc(ctx, db.ListCyclesByCreatedAtDescParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	case params.SortBy == CycleSortByLengthWeeks && params.SortOrder == SortAsc:
		dbCycles, err = r.queries.ListCyclesByLengthWeeksAsc(ctx, db.ListCyclesByLengthWeeksAscParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	case params.SortBy == CycleSortByLengthWeeks && params.SortOrder == SortDesc:
		dbCycles, err = r.queries.ListCyclesByLengthWeeksDesc(ctx, db.ListCyclesByLengthWeeksDescParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	default:
		dbCycles, err = r.queries.ListCyclesByNameAsc(ctx, db.ListCyclesByNameAscParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list cycles: %w", err)
	}

	cycles := make([]cycle.Cycle, len(dbCycles))
	for i, dbCycle := range dbCycles {
		cycles[i] = *dbCycleToDomain(dbCycle)
	}

	return cycles, total, nil
}

// Create persists a new cycle to the database.
func (r *CycleRepository) Create(c *cycle.Cycle) error {
	ctx := context.Background()

	err := r.queries.CreateCycle(ctx, db.CreateCycleParams{
		ID:          c.ID,
		Name:        c.Name,
		LengthWeeks: int64(c.LengthWeeks),
		CreatedAt:   c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   c.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to create cycle: %w", err)
	}
	return nil
}

// Update persists changes to an existing cycle.
func (r *CycleRepository) Update(c *cycle.Cycle) error {
	ctx := context.Background()

	err := r.queries.UpdateCycle(ctx, db.UpdateCycleParams{
		ID:          c.ID,
		Name:        c.Name,
		LengthWeeks: int64(c.LengthWeeks),
		UpdatedAt:   c.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to update cycle: %w", err)
	}
	return nil
}

// Delete removes a cycle from the database.
func (r *CycleRepository) Delete(id string) error {
	ctx := context.Background()

	err := r.queries.DeleteCycle(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete cycle: %w", err)
	}
	return nil
}

// IsUsedByPrograms checks if a cycle is used by any programs.
func (r *CycleRepository) IsUsedByPrograms(id string) (bool, error) {
	ctx := context.Background()

	isUsed, err := r.queries.CycleIsUsedByPrograms(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to check if cycle is used by programs: %w", err)
	}
	return isUsed == 1, nil
}

// CountWeeks counts the number of weeks associated with a cycle.
func (r *CycleRepository) CountWeeks(cycleID string) (int64, error) {
	ctx := context.Background()

	count, err := r.queries.CountWeeksByCycleID(ctx, cycleID)
	if err != nil {
		return 0, fmt.Errorf("failed to count weeks: %w", err)
	}
	return count, nil
}

// ListWeeks retrieves all weeks for a cycle ordered by week number.
func (r *CycleRepository) ListWeeks(cycleID string) ([]cycle.CycleWeek, error) {
	ctx := context.Background()

	dbWeeks, err := r.queries.ListWeeksByCycleID(ctx, cycleID)
	if err != nil {
		return nil, fmt.Errorf("failed to list weeks: %w", err)
	}

	weeks := make([]cycle.CycleWeek, len(dbWeeks))
	for i, dbWeek := range dbWeeks {
		weeks[i] = cycle.CycleWeek{
			ID:         dbWeek.ID,
			WeekNumber: int(dbWeek.WeekNumber),
		}
	}

	return weeks, nil
}

// Helper functions

func dbCycleToDomain(dbCycle db.Cycle) *cycle.Cycle {
	createdAt, _ := time.Parse(time.RFC3339, dbCycle.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, dbCycle.UpdatedAt)

	return &cycle.Cycle{
		ID:          dbCycle.ID,
		Name:        dbCycle.Name,
		LengthWeeks: int(dbCycle.LengthWeeks),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

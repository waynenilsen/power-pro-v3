// Package repository provides database repository implementations.
package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/db"
	"github.com/waynenilsen/power-pro-v3/internal/domain/dailylookup"
)

// DailyLookupRepository implements daily lookup persistence using sqlc-generated queries.
type DailyLookupRepository struct {
	queries *db.Queries
}

// NewDailyLookupRepository creates a new DailyLookupRepository.
func NewDailyLookupRepository(sqlDB *sql.DB) *DailyLookupRepository {
	return &DailyLookupRepository{
		queries: db.New(sqlDB),
	}
}

// DailyLookupSortField represents a field to sort by.
type DailyLookupSortField string

const (
	DailyLookupSortByName      DailyLookupSortField = "name"
	DailyLookupSortByCreatedAt DailyLookupSortField = "created_at"
)

// DailyLookupListParams contains parameters for listing daily lookups.
type DailyLookupListParams struct {
	Limit     int64
	Offset    int64
	SortBy    DailyLookupSortField
	SortOrder SortOrder
}

// GetByID retrieves a daily lookup by its ID.
func (r *DailyLookupRepository) GetByID(id string) (*dailylookup.DailyLookup, error) {
	ctx := context.Background()
	dbLookup, err := r.queries.GetDailyLookup(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get daily lookup: %w", err)
	}
	return dbDailyLookupToDomain(dbLookup)
}

// List retrieves daily lookups with pagination and sorting.
func (r *DailyLookupRepository) List(params DailyLookupListParams) ([]dailylookup.DailyLookup, int64, error) {
	ctx := context.Background()

	// Set defaults
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = DailyLookupSortByName
	}
	if params.SortOrder == "" {
		params.SortOrder = SortAsc
	}

	total, err := r.queries.CountDailyLookups(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count daily lookups: %w", err)
	}

	var dbLookups []db.DailyLookup

	// Select appropriate query based on sort
	switch {
	case params.SortBy == DailyLookupSortByName && params.SortOrder == SortAsc:
		dbLookups, err = r.queries.ListDailyLookupsByNameAsc(ctx, db.ListDailyLookupsByNameAscParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	case params.SortBy == DailyLookupSortByName && params.SortOrder == SortDesc:
		dbLookups, err = r.queries.ListDailyLookupsByNameDesc(ctx, db.ListDailyLookupsByNameDescParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	case params.SortBy == DailyLookupSortByCreatedAt && params.SortOrder == SortAsc:
		dbLookups, err = r.queries.ListDailyLookupsByCreatedAtAsc(ctx, db.ListDailyLookupsByCreatedAtAscParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	case params.SortBy == DailyLookupSortByCreatedAt && params.SortOrder == SortDesc:
		dbLookups, err = r.queries.ListDailyLookupsByCreatedAtDesc(ctx, db.ListDailyLookupsByCreatedAtDescParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	default:
		dbLookups, err = r.queries.ListDailyLookupsByNameAsc(ctx, db.ListDailyLookupsByNameAscParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list daily lookups: %w", err)
	}

	lookups := make([]dailylookup.DailyLookup, 0, len(dbLookups))
	for _, dbLookup := range dbLookups {
		lookup, err := dbDailyLookupToDomain(dbLookup)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to convert daily lookup: %w", err)
		}
		lookups = append(lookups, *lookup)
	}

	return lookups, total, nil
}

// Create persists a new daily lookup to the database.
func (r *DailyLookupRepository) Create(d *dailylookup.DailyLookup) error {
	ctx := context.Background()

	entriesJSON, err := json.Marshal(d.Entries)
	if err != nil {
		return fmt.Errorf("failed to marshal entries: %w", err)
	}

	err = r.queries.CreateDailyLookup(ctx, db.CreateDailyLookupParams{
		ID:        d.ID,
		Name:      d.Name,
		Entries:   string(entriesJSON),
		ProgramID: nullStringPtr(d.ProgramID),
		CreatedAt: d.CreatedAt.Format(time.RFC3339),
		UpdatedAt: d.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to create daily lookup: %w", err)
	}
	return nil
}

// Update persists changes to an existing daily lookup.
func (r *DailyLookupRepository) Update(d *dailylookup.DailyLookup) error {
	ctx := context.Background()

	entriesJSON, err := json.Marshal(d.Entries)
	if err != nil {
		return fmt.Errorf("failed to marshal entries: %w", err)
	}

	err = r.queries.UpdateDailyLookup(ctx, db.UpdateDailyLookupParams{
		ID:        d.ID,
		Name:      d.Name,
		Entries:   string(entriesJSON),
		ProgramID: nullStringPtr(d.ProgramID),
		UpdatedAt: d.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to update daily lookup: %w", err)
	}
	return nil
}

// Delete removes a daily lookup from the database.
func (r *DailyLookupRepository) Delete(id string) error {
	ctx := context.Background()

	err := r.queries.DeleteDailyLookup(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete daily lookup: %w", err)
	}
	return nil
}

// IsUsedByPrograms checks if a daily lookup is used by any programs.
func (r *DailyLookupRepository) IsUsedByPrograms(id string) (bool, error) {
	ctx := context.Background()

	isUsed, err := r.queries.DailyLookupIsUsedByPrograms(ctx, sql.NullString{String: id, Valid: true})
	if err != nil {
		return false, fmt.Errorf("failed to check if daily lookup is used by programs: %w", err)
	}
	return isUsed == 1, nil
}

// Helper functions

func dbDailyLookupToDomain(dbLookup db.DailyLookup) (*dailylookup.DailyLookup, error) {
	createdAt, _ := time.Parse(time.RFC3339, dbLookup.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, dbLookup.UpdatedAt)

	var entries []dailylookup.DailyLookupEntry
	if err := json.Unmarshal([]byte(dbLookup.Entries), &entries); err != nil {
		return nil, fmt.Errorf("failed to unmarshal entries: %w", err)
	}

	var programID *string
	if dbLookup.ProgramID.Valid {
		programID = &dbLookup.ProgramID.String
	}

	return &dailylookup.DailyLookup{
		ID:        dbLookup.ID,
		Name:      dbLookup.Name,
		Entries:   entries,
		ProgramID: programID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

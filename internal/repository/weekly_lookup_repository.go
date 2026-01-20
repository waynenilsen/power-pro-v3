// Package repository provides database repository implementations.
package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/db"
	"github.com/waynenilsen/power-pro-v3/internal/domain/weeklylookup"
)

// WeeklyLookupRepository implements weekly lookup persistence using sqlc-generated queries.
type WeeklyLookupRepository struct {
	queries *db.Queries
}

// NewWeeklyLookupRepository creates a new WeeklyLookupRepository.
func NewWeeklyLookupRepository(sqlDB *sql.DB) *WeeklyLookupRepository {
	return &WeeklyLookupRepository{
		queries: db.New(sqlDB),
	}
}

// WeeklyLookupSortField represents a field to sort by.
type WeeklyLookupSortField string

const (
	WeeklyLookupSortByName      WeeklyLookupSortField = "name"
	WeeklyLookupSortByCreatedAt WeeklyLookupSortField = "created_at"
)

// WeeklyLookupListParams contains parameters for listing weekly lookups.
type WeeklyLookupListParams struct {
	Limit     int64
	Offset    int64
	SortBy    WeeklyLookupSortField
	SortOrder SortOrder
}

// GetByID retrieves a weekly lookup by its ID.
func (r *WeeklyLookupRepository) GetByID(id string) (*weeklylookup.WeeklyLookup, error) {
	ctx := context.Background()
	dbLookup, err := r.queries.GetWeeklyLookup(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get weekly lookup: %w", err)
	}
	return dbWeeklyLookupToDomain(dbLookup)
}

// List retrieves weekly lookups with pagination and sorting.
func (r *WeeklyLookupRepository) List(params WeeklyLookupListParams) ([]weeklylookup.WeeklyLookup, int64, error) {
	ctx := context.Background()

	// Set defaults
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = WeeklyLookupSortByName
	}
	if params.SortOrder == "" {
		params.SortOrder = SortAsc
	}

	total, err := r.queries.CountWeeklyLookups(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count weekly lookups: %w", err)
	}

	var dbLookups []db.WeeklyLookup

	// Select appropriate query based on sort
	switch {
	case params.SortBy == WeeklyLookupSortByName && params.SortOrder == SortAsc:
		dbLookups, err = r.queries.ListWeeklyLookupsByNameAsc(ctx, db.ListWeeklyLookupsByNameAscParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	case params.SortBy == WeeklyLookupSortByName && params.SortOrder == SortDesc:
		dbLookups, err = r.queries.ListWeeklyLookupsByNameDesc(ctx, db.ListWeeklyLookupsByNameDescParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	case params.SortBy == WeeklyLookupSortByCreatedAt && params.SortOrder == SortAsc:
		dbLookups, err = r.queries.ListWeeklyLookupsByCreatedAtAsc(ctx, db.ListWeeklyLookupsByCreatedAtAscParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	case params.SortBy == WeeklyLookupSortByCreatedAt && params.SortOrder == SortDesc:
		dbLookups, err = r.queries.ListWeeklyLookupsByCreatedAtDesc(ctx, db.ListWeeklyLookupsByCreatedAtDescParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	default:
		dbLookups, err = r.queries.ListWeeklyLookupsByNameAsc(ctx, db.ListWeeklyLookupsByNameAscParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list weekly lookups: %w", err)
	}

	lookups := make([]weeklylookup.WeeklyLookup, 0, len(dbLookups))
	for _, dbLookup := range dbLookups {
		lookup, err := dbWeeklyLookupToDomain(dbLookup)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to convert weekly lookup: %w", err)
		}
		lookups = append(lookups, *lookup)
	}

	return lookups, total, nil
}

// Create persists a new weekly lookup to the database.
func (r *WeeklyLookupRepository) Create(w *weeklylookup.WeeklyLookup) error {
	ctx := context.Background()

	entriesJSON, err := json.Marshal(w.Entries)
	if err != nil {
		return fmt.Errorf("failed to marshal entries: %w", err)
	}

	err = r.queries.CreateWeeklyLookup(ctx, db.CreateWeeklyLookupParams{
		ID:        w.ID,
		Name:      w.Name,
		Entries:   string(entriesJSON),
		ProgramID: nullStringPtr(w.ProgramID),
		CreatedAt: w.CreatedAt.Format(time.RFC3339),
		UpdatedAt: w.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to create weekly lookup: %w", err)
	}
	return nil
}

// Update persists changes to an existing weekly lookup.
func (r *WeeklyLookupRepository) Update(w *weeklylookup.WeeklyLookup) error {
	ctx := context.Background()

	entriesJSON, err := json.Marshal(w.Entries)
	if err != nil {
		return fmt.Errorf("failed to marshal entries: %w", err)
	}

	err = r.queries.UpdateWeeklyLookup(ctx, db.UpdateWeeklyLookupParams{
		ID:        w.ID,
		Name:      w.Name,
		Entries:   string(entriesJSON),
		ProgramID: nullStringPtr(w.ProgramID),
		UpdatedAt: w.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to update weekly lookup: %w", err)
	}
	return nil
}

// Delete removes a weekly lookup from the database.
func (r *WeeklyLookupRepository) Delete(id string) error {
	ctx := context.Background()

	err := r.queries.DeleteWeeklyLookup(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete weekly lookup: %w", err)
	}
	return nil
}

// IsUsedByPrograms checks if a weekly lookup is used by any programs.
func (r *WeeklyLookupRepository) IsUsedByPrograms(id string) (bool, error) {
	ctx := context.Background()

	isUsed, err := r.queries.WeeklyLookupIsUsedByPrograms(ctx, sql.NullString{String: id, Valid: true})
	if err != nil {
		return false, fmt.Errorf("failed to check if weekly lookup is used by programs: %w", err)
	}
	return isUsed == 1, nil
}

// Helper functions

func dbWeeklyLookupToDomain(dbLookup db.WeeklyLookup) (*weeklylookup.WeeklyLookup, error) {
	createdAt, _ := time.Parse(time.RFC3339, dbLookup.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, dbLookup.UpdatedAt)

	var entries []weeklylookup.WeeklyLookupEntry
	if err := json.Unmarshal([]byte(dbLookup.Entries), &entries); err != nil {
		return nil, fmt.Errorf("failed to unmarshal entries: %w", err)
	}

	var programID *string
	if dbLookup.ProgramID.Valid {
		programID = &dbLookup.ProgramID.String
	}

	return &weeklylookup.WeeklyLookup{
		ID:        dbLookup.ID,
		Name:      dbLookup.Name,
		Entries:   entries,
		ProgramID: programID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func nullStringPtr(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

// Package repository provides database repository implementations.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/db"
	"github.com/waynenilsen/power-pro-v3/internal/domain/program"
)

// ProgramRepository implements program persistence using sqlc-generated queries.
type ProgramRepository struct {
	queries *db.Queries
}

// NewProgramRepository creates a new ProgramRepository.
func NewProgramRepository(sqlDB *sql.DB) *ProgramRepository {
	return &ProgramRepository{
		queries: db.New(sqlDB),
	}
}

// ProgramSortField represents a field to sort by.
type ProgramSortField string

const (
	ProgramSortByName      ProgramSortField = "name"
	ProgramSortByCreatedAt ProgramSortField = "created_at"
)

// ProgramListParams contains parameters for listing programs.
type ProgramListParams struct {
	Limit     int64
	Offset    int64
	SortBy    ProgramSortField
	SortOrder SortOrder
}

// GetByID retrieves a program by its ID.
func (r *ProgramRepository) GetByID(id string) (*program.Program, error) {
	ctx := context.Background()
	dbProgram, err := r.queries.GetProgram(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get program: %w", err)
	}
	return dbProgramToDomain(dbProgram), nil
}

// GetBySlug retrieves a program by its slug.
func (r *ProgramRepository) GetBySlug(slug string) (*program.Program, error) {
	ctx := context.Background()
	dbProgram, err := r.queries.GetProgramBySlug(ctx, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get program by slug: %w", err)
	}
	return dbProgramToDomain(dbProgram), nil
}

// List retrieves programs with pagination and sorting.
func (r *ProgramRepository) List(params ProgramListParams) ([]program.Program, int64, error) {
	ctx := context.Background()

	// Set defaults
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = ProgramSortByName
	}
	if params.SortOrder == "" {
		params.SortOrder = SortAsc
	}

	total, err := r.queries.CountPrograms(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count programs: %w", err)
	}

	var dbPrograms []db.Program

	// Select appropriate query based on sort
	switch {
	case params.SortBy == ProgramSortByName && params.SortOrder == SortAsc:
		dbPrograms, err = r.queries.ListProgramsByNameAsc(ctx, db.ListProgramsByNameAscParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	case params.SortBy == ProgramSortByName && params.SortOrder == SortDesc:
		dbPrograms, err = r.queries.ListProgramsByNameDesc(ctx, db.ListProgramsByNameDescParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	case params.SortBy == ProgramSortByCreatedAt && params.SortOrder == SortAsc:
		dbPrograms, err = r.queries.ListProgramsByCreatedAtAsc(ctx, db.ListProgramsByCreatedAtAscParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	case params.SortBy == ProgramSortByCreatedAt && params.SortOrder == SortDesc:
		dbPrograms, err = r.queries.ListProgramsByCreatedAtDesc(ctx, db.ListProgramsByCreatedAtDescParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	default:
		dbPrograms, err = r.queries.ListProgramsByNameAsc(ctx, db.ListProgramsByNameAscParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list programs: %w", err)
	}

	programs := make([]program.Program, len(dbPrograms))
	for i, dbProg := range dbPrograms {
		programs[i] = *dbProgramToDomain(dbProg)
	}

	return programs, total, nil
}

// Create persists a new program to the database.
func (r *ProgramRepository) Create(p *program.Program) error {
	ctx := context.Background()

	err := r.queries.CreateProgram(ctx, db.CreateProgramParams{
		ID:              p.ID,
		Name:            p.Name,
		Slug:            p.Slug,
		Description:     stringPtrToNullString(p.Description),
		CycleID:         p.CycleID,
		WeeklyLookupID:  stringPtrToNullString(p.WeeklyLookupID),
		DailyLookupID:   stringPtrToNullString(p.DailyLookupID),
		DefaultRounding: programFloat64PtrToNullFloat64(p.DefaultRounding),
		CreatedAt:       p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       p.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to create program: %w", err)
	}
	return nil
}

// Update persists changes to an existing program.
func (r *ProgramRepository) Update(p *program.Program) error {
	ctx := context.Background()

	err := r.queries.UpdateProgram(ctx, db.UpdateProgramParams{
		ID:              p.ID,
		Name:            p.Name,
		Slug:            p.Slug,
		Description:     stringPtrToNullString(p.Description),
		CycleID:         p.CycleID,
		WeeklyLookupID:  stringPtrToNullString(p.WeeklyLookupID),
		DailyLookupID:   stringPtrToNullString(p.DailyLookupID),
		DefaultRounding: programFloat64PtrToNullFloat64(p.DefaultRounding),
		UpdatedAt:       p.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to update program: %w", err)
	}
	return nil
}

// Delete removes a program from the database.
func (r *ProgramRepository) Delete(id string) error {
	ctx := context.Background()

	err := r.queries.DeleteProgram(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete program: %w", err)
	}
	return nil
}

// SlugExists checks if a program with the given slug exists.
func (r *ProgramRepository) SlugExists(slug string) (bool, error) {
	ctx := context.Background()

	exists, err := r.queries.ProgramSlugExists(ctx, slug)
	if err != nil {
		return false, fmt.Errorf("failed to check if slug exists: %w", err)
	}
	return exists == 1, nil
}

// SlugExistsExcluding checks if a program with the given slug exists, excluding a specific ID.
func (r *ProgramRepository) SlugExistsExcluding(slug string, excludeID string) (bool, error) {
	ctx := context.Background()

	exists, err := r.queries.ProgramSlugExistsExcluding(ctx, db.ProgramSlugExistsExcludingParams{
		Slug: slug,
		ID:   excludeID,
	})
	if err != nil {
		return false, fmt.Errorf("failed to check if slug exists: %w", err)
	}
	return exists == 1, nil
}

// HasEnrolledUsers checks if any users are enrolled in the program.
func (r *ProgramRepository) HasEnrolledUsers(id string) (bool, error) {
	ctx := context.Background()

	hasEnrolled, err := r.queries.ProgramHasEnrolledUsers(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to check if program has enrolled users: %w", err)
	}
	return hasEnrolled == 1, nil
}

// CountEnrolledUsers returns the count of users enrolled in the program.
func (r *ProgramRepository) CountEnrolledUsers(id string) (int64, error) {
	ctx := context.Background()

	count, err := r.queries.CountEnrolledUsers(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("failed to count enrolled users: %w", err)
	}
	return count, nil
}

// GetCycleForProgram retrieves the cycle associated with a program.
func (r *ProgramRepository) GetCycleForProgram(cycleID string) (*program.ProgramCycle, error) {
	ctx := context.Background()

	dbCycle, err := r.queries.GetCycleForProgram(ctx, cycleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get cycle for program: %w", err)
	}

	// Get weeks for the cycle
	dbWeeks, err := r.queries.ListWeeksByCycleID(ctx, cycleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get weeks for cycle: %w", err)
	}

	weeks := make([]program.ProgramCycleWeek, len(dbWeeks))
	for i, w := range dbWeeks {
		weeks[i] = program.ProgramCycleWeek{
			ID:         w.ID,
			WeekNumber: int(w.WeekNumber),
		}
	}

	return &program.ProgramCycle{
		ID:          dbCycle.ID,
		Name:        dbCycle.Name,
		LengthWeeks: int(dbCycle.LengthWeeks),
		Weeks:       weeks,
	}, nil
}

// GetWeeklyLookupReference retrieves a weekly lookup reference (id and name only).
func (r *ProgramRepository) GetWeeklyLookupReference(id string) (*program.LookupReference, error) {
	ctx := context.Background()

	dbLookup, err := r.queries.GetWeeklyLookupForProgram(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get weekly lookup: %w", err)
	}

	return &program.LookupReference{
		ID:   dbLookup.ID,
		Name: dbLookup.Name,
	}, nil
}

// GetDailyLookupReference retrieves a daily lookup reference (id and name only).
func (r *ProgramRepository) GetDailyLookupReference(id string) (*program.LookupReference, error) {
	ctx := context.Background()

	dbLookup, err := r.queries.GetDailyLookupForProgram(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get daily lookup: %w", err)
	}

	return &program.LookupReference{
		ID:   dbLookup.ID,
		Name: dbLookup.Name,
	}, nil
}

// Helper functions

func dbProgramToDomain(dbProg db.Program) *program.Program {
	createdAt, _ := time.Parse(time.RFC3339, dbProg.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, dbProg.UpdatedAt)

	return &program.Program{
		ID:              dbProg.ID,
		Name:            dbProg.Name,
		Slug:            dbProg.Slug,
		Description:     nullStringToStringPtr(dbProg.Description),
		CycleID:         dbProg.CycleID,
		WeeklyLookupID:  nullStringToStringPtr(dbProg.WeeklyLookupID),
		DailyLookupID:   nullStringToStringPtr(dbProg.DailyLookupID),
		DefaultRounding: programNullFloat64ToFloat64Ptr(dbProg.DefaultRounding),
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
}

// Note: stringPtrToNullString and nullStringToStringPtr are defined in lift_repository.go
// They are shared across repositories in the same package.

func programFloat64PtrToNullFloat64(f *float64) sql.NullFloat64 {
	if f == nil {
		return sql.NullFloat64{Valid: false}
	}
	return sql.NullFloat64{Float64: *f, Valid: true}
}

func programNullFloat64ToFloat64Ptr(nf sql.NullFloat64) *float64 {
	if !nf.Valid {
		return nil
	}
	return &nf.Float64
}

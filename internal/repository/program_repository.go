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
	Filters   *program.FilterOptions
}

// GetByID retrieves a program by its ID.
func (r *ProgramRepository) GetByID(id string) (*program.Program, error) {
	ctx := context.Background()
	row, err := r.queries.GetProgram(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get program: %w", err)
	}
	return dbGetProgramRowToDomain(row), nil
}

// GetBySlug retrieves a program by its slug.
func (r *ProgramRepository) GetBySlug(slug string) (*program.Program, error) {
	ctx := context.Background()
	row, err := r.queries.GetProgramBySlug(ctx, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get program by slug: %w", err)
	}
	return dbGetProgramBySlugRowToDomain(row), nil
}

// List retrieves programs with pagination, sorting, and optional filtering.
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

	// Check if we have any filters
	hasFilters := params.Filters != nil && (params.Filters.Difficulty != nil ||
		params.Filters.DaysPerWeek != nil ||
		params.Filters.Focus != nil ||
		params.Filters.HasAmrap != nil ||
		params.Filters.Search != nil)

	var total int64
	var err error

	if hasFilters {
		// Use filtered count
		filterParams := buildFilterParams(params.Filters)
		total, err = r.queries.CountProgramsFiltered(ctx, db.CountProgramsFilteredParams{
			Difficulty:  filterParams.difficulty,
			DaysPerWeek: filterParams.daysPerWeek,
			Focus:       filterParams.focus,
			HasAmrap:    filterParams.hasAmrap,
			Search:      filterParams.search,
		})
	} else {
		total, err = r.queries.CountPrograms(ctx)
	}
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count programs: %w", err)
	}

	var programs []program.Program

	if hasFilters {
		programs, err = r.listProgramsFiltered(ctx, params)
	} else {
		programs, err = r.listProgramsUnfiltered(ctx, params)
	}

	if err != nil {
		return nil, 0, err
	}

	return programs, total, nil
}

// filterParamsHelper holds the converted filter parameters for sqlc queries.
type filterParamsHelper struct {
	difficulty  interface{}
	daysPerWeek interface{}
	focus       interface{}
	hasAmrap    interface{}
	search      interface{}
}

// buildFilterParams converts FilterOptions to the interface{} types sqlc expects.
func buildFilterParams(filters *program.FilterOptions) filterParamsHelper {
	var params filterParamsHelper
	if filters == nil {
		return params
	}
	if filters.Difficulty != nil {
		params.difficulty = *filters.Difficulty
	}
	if filters.DaysPerWeek != nil {
		params.daysPerWeek = int64(*filters.DaysPerWeek)
	}
	if filters.Focus != nil {
		params.focus = *filters.Focus
	}
	if filters.HasAmrap != nil {
		if *filters.HasAmrap {
			params.hasAmrap = int64(1)
		} else {
			params.hasAmrap = int64(0)
		}
	}
	if filters.Search != nil {
		params.search = *filters.Search
	}
	return params
}

// listProgramsFiltered handles filtered program listing with sorting.
func (r *ProgramRepository) listProgramsFiltered(ctx context.Context, params ProgramListParams) ([]program.Program, error) {
	filterParams := buildFilterParams(params.Filters)

	var rows []db.ListProgramsFilteredByNameAscRow
	var err error

	switch {
	case params.SortBy == ProgramSortByName && params.SortOrder == SortAsc:
		rows, err = r.queries.ListProgramsFilteredByNameAsc(ctx, db.ListProgramsFilteredByNameAscParams{
			Difficulty:  filterParams.difficulty,
			DaysPerWeek: filterParams.daysPerWeek,
			Focus:       filterParams.focus,
			HasAmrap:    filterParams.hasAmrap,
			Search:      filterParams.search,
			Limit:       params.Limit,
			Offset:      params.Offset,
		})
	case params.SortBy == ProgramSortByName && params.SortOrder == SortDesc:
		descRows, e := r.queries.ListProgramsFilteredByNameDesc(ctx, db.ListProgramsFilteredByNameDescParams{
			Difficulty:  filterParams.difficulty,
			DaysPerWeek: filterParams.daysPerWeek,
			Focus:       filterParams.focus,
			HasAmrap:    filterParams.hasAmrap,
			Search:      filterParams.search,
			Limit:       params.Limit,
			Offset:      params.Offset,
		})
		if e != nil {
			return nil, fmt.Errorf("failed to list programs: %w", e)
		}
		// Convert to common row type
		rows = make([]db.ListProgramsFilteredByNameAscRow, len(descRows))
		for i, r := range descRows {
			rows[i] = db.ListProgramsFilteredByNameAscRow(r)
		}
	case params.SortBy == ProgramSortByCreatedAt && params.SortOrder == SortAsc:
		ascRows, e := r.queries.ListProgramsFilteredByCreatedAtAsc(ctx, db.ListProgramsFilteredByCreatedAtAscParams{
			Difficulty:  filterParams.difficulty,
			DaysPerWeek: filterParams.daysPerWeek,
			Focus:       filterParams.focus,
			HasAmrap:    filterParams.hasAmrap,
			Search:      filterParams.search,
			Limit:       params.Limit,
			Offset:      params.Offset,
		})
		if e != nil {
			return nil, fmt.Errorf("failed to list programs: %w", e)
		}
		rows = make([]db.ListProgramsFilteredByNameAscRow, len(ascRows))
		for i, r := range ascRows {
			rows[i] = db.ListProgramsFilteredByNameAscRow(r)
		}
	case params.SortBy == ProgramSortByCreatedAt && params.SortOrder == SortDesc:
		descRows, e := r.queries.ListProgramsFilteredByCreatedAtDesc(ctx, db.ListProgramsFilteredByCreatedAtDescParams{
			Difficulty:  filterParams.difficulty,
			DaysPerWeek: filterParams.daysPerWeek,
			Focus:       filterParams.focus,
			HasAmrap:    filterParams.hasAmrap,
			Search:      filterParams.search,
			Limit:       params.Limit,
			Offset:      params.Offset,
		})
		if e != nil {
			return nil, fmt.Errorf("failed to list programs: %w", e)
		}
		rows = make([]db.ListProgramsFilteredByNameAscRow, len(descRows))
		for i, r := range descRows {
			rows[i] = db.ListProgramsFilteredByNameAscRow(r)
		}
	default:
		rows, err = r.queries.ListProgramsFilteredByNameAsc(ctx, db.ListProgramsFilteredByNameAscParams{
			Difficulty:  filterParams.difficulty,
			DaysPerWeek: filterParams.daysPerWeek,
			Focus:       filterParams.focus,
			HasAmrap:    filterParams.hasAmrap,
			Search:      filterParams.search,
			Limit:       params.Limit,
			Offset:      params.Offset,
		})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list programs: %w", err)
	}

	programs := make([]program.Program, len(rows))
	for i, row := range rows {
		programs[i] = *dbFilteredRowToDomain(row)
	}

	return programs, nil
}

// listProgramsUnfiltered handles unfiltered program listing with sorting.
func (r *ProgramRepository) listProgramsUnfiltered(ctx context.Context, params ProgramListParams) ([]program.Program, error) {
	var err error
	var rows []db.ListProgramsByNameAscRow

	switch {
	case params.SortBy == ProgramSortByName && params.SortOrder == SortAsc:
		rows, err = r.queries.ListProgramsByNameAsc(ctx, db.ListProgramsByNameAscParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	case params.SortBy == ProgramSortByName && params.SortOrder == SortDesc:
		descRows, e := r.queries.ListProgramsByNameDesc(ctx, db.ListProgramsByNameDescParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
		if e != nil {
			return nil, fmt.Errorf("failed to list programs: %w", e)
		}
		rows = make([]db.ListProgramsByNameAscRow, len(descRows))
		for i, r := range descRows {
			rows[i] = db.ListProgramsByNameAscRow(r)
		}
	case params.SortBy == ProgramSortByCreatedAt && params.SortOrder == SortAsc:
		ascRows, e := r.queries.ListProgramsByCreatedAtAsc(ctx, db.ListProgramsByCreatedAtAscParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
		if e != nil {
			return nil, fmt.Errorf("failed to list programs: %w", e)
		}
		rows = make([]db.ListProgramsByNameAscRow, len(ascRows))
		for i, r := range ascRows {
			rows[i] = db.ListProgramsByNameAscRow(r)
		}
	case params.SortBy == ProgramSortByCreatedAt && params.SortOrder == SortDesc:
		descRows, e := r.queries.ListProgramsByCreatedAtDesc(ctx, db.ListProgramsByCreatedAtDescParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
		if e != nil {
			return nil, fmt.Errorf("failed to list programs: %w", e)
		}
		rows = make([]db.ListProgramsByNameAscRow, len(descRows))
		for i, r := range descRows {
			rows[i] = db.ListProgramsByNameAscRow(r)
		}
	default:
		rows, err = r.queries.ListProgramsByNameAsc(ctx, db.ListProgramsByNameAscParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list programs: %w", err)
	}

	programs := make([]program.Program, len(rows))
	for i, row := range rows {
		programs[i] = *dbUnfilteredRowToDomain(row)
	}

	return programs, nil
}

// Create persists a new program to the database.
func (r *ProgramRepository) Create(p *program.Program) error {
	ctx := context.Background()

	var hasAmrap int64
	if p.HasAmrap {
		hasAmrap = 1
	}

	err := r.queries.CreateProgram(ctx, db.CreateProgramParams{
		ID:              p.ID,
		Name:            p.Name,
		Slug:            p.Slug,
		Description:     stringPtrToNullString(p.Description),
		CycleID:         p.CycleID,
		WeeklyLookupID:  stringPtrToNullString(p.WeeklyLookupID),
		DailyLookupID:   stringPtrToNullString(p.DailyLookupID),
		DefaultRounding: programFloat64PtrToNullFloat64(p.DefaultRounding),
		Difficulty:      p.Difficulty,
		DaysPerWeek:     int64(p.DaysPerWeek),
		Focus:           p.Focus,
		HasAmrap:        hasAmrap,
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

	var hasAmrap int64
	if p.HasAmrap {
		hasAmrap = 1
	}

	err := r.queries.UpdateProgram(ctx, db.UpdateProgramParams{
		ID:              p.ID,
		Name:            p.Name,
		Slug:            p.Slug,
		Description:     stringPtrToNullString(p.Description),
		CycleID:         p.CycleID,
		WeeklyLookupID:  stringPtrToNullString(p.WeeklyLookupID),
		DailyLookupID:   stringPtrToNullString(p.DailyLookupID),
		DefaultRounding: programFloat64PtrToNullFloat64(p.DefaultRounding),
		Difficulty:      p.Difficulty,
		DaysPerWeek:     int64(p.DaysPerWeek),
		Focus:           p.Focus,
		HasAmrap:        hasAmrap,
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

// SampleWeekDay represents a day in the program's sample week.
type SampleWeekDay struct {
	ID            string
	Name          string
	DayOfWeek     string
	ExerciseCount int
}

// GetSampleWeek retrieves the first week's structure for a program.
func (r *ProgramRepository) GetSampleWeek(programID string) ([]SampleWeekDay, error) {
	ctx := context.Background()

	rows, err := r.queries.GetProgramSampleWeek(ctx, db.GetProgramSampleWeekParams{
		ID:        programID,
		ProgramID: sql.NullString{String: programID, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get sample week: %w", err)
	}

	days := make([]SampleWeekDay, len(rows))
	for i, row := range rows {
		days[i] = SampleWeekDay{
			ID:            row.ID,
			Name:          row.Name,
			DayOfWeek:     row.DayOfWeek,
			ExerciseCount: int(row.ExerciseCount),
		}
	}

	return days, nil
}

// GetLiftRequirements retrieves the unique lifts used in a program.
func (r *ProgramRepository) GetLiftRequirements(programID string) ([]string, error) {
	ctx := context.Background()

	lifts, err := r.queries.GetProgramLiftRequirements(ctx, sql.NullString{String: programID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get lift requirements: %w", err)
	}

	return lifts, nil
}

// SessionStats holds aggregate stats for estimating session duration.
type SessionStats struct {
	TotalSets      int
	TotalDays      int
	TotalExercises int
}

// GetSessionStats retrieves aggregate stats for session duration estimation.
func (r *ProgramRepository) GetSessionStats(programID string) (*SessionStats, error) {
	ctx := context.Background()

	row, err := r.queries.GetProgramSessionStats(ctx, sql.NullString{String: programID, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			return &SessionStats{}, nil
		}
		return nil, fmt.Errorf("failed to get session stats: %w", err)
	}

	// The interface{} types from sqlc need type assertion
	totalSets := toInt(row.TotalSets)
	totalDays := toInt(row.TotalDays)
	totalExercises := toInt(row.TotalExercises)

	return &SessionStats{
		TotalSets:      totalSets,
		TotalDays:      totalDays,
		TotalExercises: totalExercises,
	}, nil
}

// toInt converts an interface{} to int, handling various numeric types.
func toInt(v interface{}) int {
	switch val := v.(type) {
	case int64:
		return int(val)
	case int:
		return val
	case float64:
		return int(val)
	default:
		return 0
	}
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
		Difficulty:      dbProg.Difficulty,
		DaysPerWeek:     int(dbProg.DaysPerWeek),
		Focus:           dbProg.Focus,
		HasAmrap:        dbProg.HasAmrap == 1,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
}

// dbFilteredRowToDomain converts a filtered list row to a domain Program.
func dbFilteredRowToDomain(row db.ListProgramsFilteredByNameAscRow) *program.Program {
	createdAt, _ := time.Parse(time.RFC3339, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, row.UpdatedAt)

	return &program.Program{
		ID:              row.ID,
		Name:            row.Name,
		Slug:            row.Slug,
		Description:     nullStringToStringPtr(row.Description),
		CycleID:         row.CycleID,
		WeeklyLookupID:  nullStringToStringPtr(row.WeeklyLookupID),
		DailyLookupID:   nullStringToStringPtr(row.DailyLookupID),
		DefaultRounding: programNullFloat64ToFloat64Ptr(row.DefaultRounding),
		Difficulty:      row.Difficulty,
		DaysPerWeek:     int(row.DaysPerWeek),
		Focus:           row.Focus,
		HasAmrap:        row.HasAmrap == 1,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
}

// dbUnfilteredRowToDomain converts an unfiltered list row to a domain Program.
func dbUnfilteredRowToDomain(row db.ListProgramsByNameAscRow) *program.Program {
	createdAt, _ := time.Parse(time.RFC3339, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, row.UpdatedAt)

	return &program.Program{
		ID:              row.ID,
		Name:            row.Name,
		Slug:            row.Slug,
		Description:     nullStringToStringPtr(row.Description),
		CycleID:         row.CycleID,
		WeeklyLookupID:  nullStringToStringPtr(row.WeeklyLookupID),
		DailyLookupID:   nullStringToStringPtr(row.DailyLookupID),
		DefaultRounding: programNullFloat64ToFloat64Ptr(row.DefaultRounding),
		Difficulty:      row.Difficulty,
		DaysPerWeek:     int(row.DaysPerWeek),
		Focus:           row.Focus,
		HasAmrap:        row.HasAmrap == 1,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
}

// dbGetProgramRowToDomain converts a GetProgram row to a domain Program.
func dbGetProgramRowToDomain(row db.GetProgramRow) *program.Program {
	createdAt, _ := time.Parse(time.RFC3339, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, row.UpdatedAt)

	return &program.Program{
		ID:              row.ID,
		Name:            row.Name,
		Slug:            row.Slug,
		Description:     nullStringToStringPtr(row.Description),
		CycleID:         row.CycleID,
		WeeklyLookupID:  nullStringToStringPtr(row.WeeklyLookupID),
		DailyLookupID:   nullStringToStringPtr(row.DailyLookupID),
		DefaultRounding: programNullFloat64ToFloat64Ptr(row.DefaultRounding),
		Difficulty:      row.Difficulty,
		DaysPerWeek:     int(row.DaysPerWeek),
		Focus:           row.Focus,
		HasAmrap:        row.HasAmrap == 1,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
}

// dbGetProgramBySlugRowToDomain converts a GetProgramBySlug row to a domain Program.
func dbGetProgramBySlugRowToDomain(row db.GetProgramBySlugRow) *program.Program {
	createdAt, _ := time.Parse(time.RFC3339, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, row.UpdatedAt)

	return &program.Program{
		ID:              row.ID,
		Name:            row.Name,
		Slug:            row.Slug,
		Description:     nullStringToStringPtr(row.Description),
		CycleID:         row.CycleID,
		WeeklyLookupID:  nullStringToStringPtr(row.WeeklyLookupID),
		DailyLookupID:   nullStringToStringPtr(row.DailyLookupID),
		DefaultRounding: programNullFloat64ToFloat64Ptr(row.DefaultRounding),
		Difficulty:      row.Difficulty,
		DaysPerWeek:     int(row.DaysPerWeek),
		Focus:           row.Focus,
		HasAmrap:        row.HasAmrap == 1,
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

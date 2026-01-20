// Package repository provides database repository implementations.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/db"
	"github.com/waynenilsen/power-pro-v3/internal/domain/week"
)

// WeekRepository implements week persistence using sqlc-generated queries.
type WeekRepository struct {
	queries *db.Queries
}

// NewWeekRepository creates a new WeekRepository.
func NewWeekRepository(sqlDB *sql.DB) *WeekRepository {
	return &WeekRepository{
		queries: db.New(sqlDB),
	}
}

// WeekSortField represents a field to sort by.
type WeekSortField string

const (
	WeekSortByWeekNumber WeekSortField = "week_number"
	WeekSortByCreatedAt  WeekSortField = "created_at"
)

// WeekListParams contains parameters for listing weeks.
type WeekListParams struct {
	Limit         int64
	Offset        int64
	SortBy        WeekSortField
	SortOrder     SortOrder
	FilterCycleID *string
}

// GetByID retrieves a week by its ID.
func (r *WeekRepository) GetByID(id string) (*week.Week, error) {
	ctx := context.Background()
	dbWeek, err := r.queries.GetWeek(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get week: %w", err)
	}
	return dbWeekToDomain(dbWeek), nil
}

// List retrieves weeks with pagination, sorting, and optional filtering.
func (r *WeekRepository) List(params WeekListParams) ([]week.Week, int64, error) {
	ctx := context.Background()

	// Set defaults
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = WeekSortByWeekNumber
	}
	if params.SortOrder == "" {
		params.SortOrder = SortAsc
	}

	var dbWeeks []db.Week
	var total int64
	var err error

	if params.FilterCycleID != nil {
		// Filter by cycle
		cycleID := *params.FilterCycleID

		total, err = r.queries.CountWeeksFilteredByCycle(ctx, cycleID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count weeks: %w", err)
		}

		// Select appropriate query based on sort
		switch {
		case params.SortBy == WeekSortByWeekNumber && params.SortOrder == SortAsc:
			dbWeeks, err = r.queries.ListWeeksFilteredByCycleByWeekNumberAsc(ctx, db.ListWeeksFilteredByCycleByWeekNumberAscParams{
				CycleID: cycleID,
				Limit:   params.Limit,
				Offset:  params.Offset,
			})
		case params.SortBy == WeekSortByWeekNumber && params.SortOrder == SortDesc:
			dbWeeks, err = r.queries.ListWeeksFilteredByCycleByWeekNumberDesc(ctx, db.ListWeeksFilteredByCycleByWeekNumberDescParams{
				CycleID: cycleID,
				Limit:   params.Limit,
				Offset:  params.Offset,
			})
		case params.SortBy == WeekSortByCreatedAt && params.SortOrder == SortAsc:
			dbWeeks, err = r.queries.ListWeeksFilteredByCycleByCreatedAtAsc(ctx, db.ListWeeksFilteredByCycleByCreatedAtAscParams{
				CycleID: cycleID,
				Limit:   params.Limit,
				Offset:  params.Offset,
			})
		case params.SortBy == WeekSortByCreatedAt && params.SortOrder == SortDesc:
			dbWeeks, err = r.queries.ListWeeksFilteredByCycleByCreatedAtDesc(ctx, db.ListWeeksFilteredByCycleByCreatedAtDescParams{
				CycleID: cycleID,
				Limit:   params.Limit,
				Offset:  params.Offset,
			})
		default:
			dbWeeks, err = r.queries.ListWeeksFilteredByCycleByWeekNumberAsc(ctx, db.ListWeeksFilteredByCycleByWeekNumberAscParams{
				CycleID: cycleID,
				Limit:   params.Limit,
				Offset:  params.Offset,
			})
		}
	} else {
		// No filter
		total, err = r.queries.CountWeeks(ctx)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count weeks: %w", err)
		}

		// Select appropriate query based on sort
		switch {
		case params.SortBy == WeekSortByWeekNumber && params.SortOrder == SortAsc:
			dbWeeks, err = r.queries.ListWeeksByWeekNumberAsc(ctx, db.ListWeeksByWeekNumberAscParams{
				Limit:  params.Limit,
				Offset: params.Offset,
			})
		case params.SortBy == WeekSortByWeekNumber && params.SortOrder == SortDesc:
			dbWeeks, err = r.queries.ListWeeksByWeekNumberDesc(ctx, db.ListWeeksByWeekNumberDescParams{
				Limit:  params.Limit,
				Offset: params.Offset,
			})
		case params.SortBy == WeekSortByCreatedAt && params.SortOrder == SortAsc:
			dbWeeks, err = r.queries.ListWeeksByCreatedAtAsc(ctx, db.ListWeeksByCreatedAtAscParams{
				Limit:  params.Limit,
				Offset: params.Offset,
			})
		case params.SortBy == WeekSortByCreatedAt && params.SortOrder == SortDesc:
			dbWeeks, err = r.queries.ListWeeksByCreatedAtDesc(ctx, db.ListWeeksByCreatedAtDescParams{
				Limit:  params.Limit,
				Offset: params.Offset,
			})
		default:
			dbWeeks, err = r.queries.ListWeeksByWeekNumberAsc(ctx, db.ListWeeksByWeekNumberAscParams{
				Limit:  params.Limit,
				Offset: params.Offset,
			})
		}
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list weeks: %w", err)
	}

	weeks := make([]week.Week, len(dbWeeks))
	for i, dbWeek := range dbWeeks {
		weeks[i] = *dbWeekToDomain(dbWeek)
	}

	return weeks, total, nil
}

// WeekNumberExistsInCycle checks if a week number already exists in a cycle, excluding a specific ID.
func (r *WeekRepository) WeekNumberExistsInCycle(cycleID string, weekNumber int, excludeID *string) (bool, error) {
	ctx := context.Background()

	if excludeID == nil {
		exists, err := r.queries.WeekNumberExistsInCycleForNew(ctx, db.WeekNumberExistsInCycleForNewParams{
			CycleID:    cycleID,
			WeekNumber: int64(weekNumber),
		})
		if err != nil {
			return false, fmt.Errorf("failed to check week number exists: %w", err)
		}
		return exists == 1, nil
	}

	exists, err := r.queries.WeekNumberExistsInCycle(ctx, db.WeekNumberExistsInCycleParams{
		CycleID:    cycleID,
		WeekNumber: int64(weekNumber),
		ID:         *excludeID,
	})
	if err != nil {
		return false, fmt.Errorf("failed to check week number exists: %w", err)
	}
	return exists == 1, nil
}

// Create persists a new week to the database.
func (r *WeekRepository) Create(w *week.Week) error {
	ctx := context.Background()

	err := r.queries.CreateWeek(ctx, db.CreateWeekParams{
		ID:         w.ID,
		WeekNumber: int64(w.WeekNumber),
		Variant:    stringPtrToNullString(w.Variant),
		CycleID:    w.CycleID,
		CreatedAt:  w.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  w.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to create week: %w", err)
	}
	return nil
}

// Update persists changes to an existing week.
func (r *WeekRepository) Update(w *week.Week) error {
	ctx := context.Background()

	err := r.queries.UpdateWeek(ctx, db.UpdateWeekParams{
		ID:         w.ID,
		WeekNumber: int64(w.WeekNumber),
		Variant:    stringPtrToNullString(w.Variant),
		CycleID:    w.CycleID,
		UpdatedAt:  w.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to update week: %w", err)
	}
	return nil
}

// Delete removes a week from the database.
func (r *WeekRepository) Delete(id string) error {
	ctx := context.Background()

	err := r.queries.DeleteWeek(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete week: %w", err)
	}
	return nil
}

// IsUsedInActiveCycle checks if a week is used in an active cycle (where users are enrolled).
func (r *WeekRepository) IsUsedInActiveCycle(id string) (bool, error) {
	ctx := context.Background()

	isUsed, err := r.queries.WeekIsUsedInActiveCycle(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to check if week is used in active cycle: %w", err)
	}
	return isUsed == 1, nil
}

// CycleExists checks if a cycle exists.
func (r *WeekRepository) CycleExists(id string) (bool, error) {
	ctx := context.Background()

	_, err := r.queries.GetCycleByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if cycle exists: %w", err)
	}
	return true, nil
}

// Week Day methods

// GetWeekDay retrieves a week day by its ID.
func (r *WeekRepository) GetWeekDay(id string) (*week.WeekDay, error) {
	ctx := context.Background()
	dbWeekDay, err := r.queries.GetWeekDay(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get week day: %w", err)
	}
	return dbWeekDayToDomain(dbWeekDay), nil
}

// GetWeekDayByWeekAndDayAndDayOfWeek retrieves a week day by week ID, day ID, and day of week.
func (r *WeekRepository) GetWeekDayByWeekAndDayAndDayOfWeek(weekID, dayID, dayOfWeek string) (*week.WeekDay, error) {
	ctx := context.Background()
	dbWeekDay, err := r.queries.GetWeekDayByWeekAndDayAndDayOfWeek(ctx, db.GetWeekDayByWeekAndDayAndDayOfWeekParams{
		WeekID:    weekID,
		DayID:     dayID,
		DayOfWeek: dayOfWeek,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get week day: %w", err)
	}
	return dbWeekDayToDomain(dbWeekDay), nil
}

// ListWeekDays retrieves all days for a week ordered by day of week.
func (r *WeekRepository) ListWeekDays(weekID string) ([]week.WeekDay, error) {
	ctx := context.Background()
	dbWeekDays, err := r.queries.ListWeekDays(ctx, weekID)
	if err != nil {
		return nil, fmt.Errorf("failed to list week days: %w", err)
	}

	weekDays := make([]week.WeekDay, len(dbWeekDays))
	for i, dbWeekDay := range dbWeekDays {
		weekDays[i] = *dbWeekDayToDomain(dbWeekDay)
	}

	return weekDays, nil
}

// CreateWeekDay adds a day to a week.
func (r *WeekRepository) CreateWeekDay(wd *week.WeekDay) error {
	ctx := context.Background()

	err := r.queries.CreateWeekDay(ctx, db.CreateWeekDayParams{
		ID:        wd.ID,
		WeekID:    wd.WeekID,
		DayID:     wd.DayID,
		DayOfWeek: string(wd.DayOfWeek),
		CreatedAt: wd.CreatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to create week day: %w", err)
	}
	return nil
}

// DeleteWeekDay removes a day from a week by its ID.
func (r *WeekRepository) DeleteWeekDay(id string) error {
	ctx := context.Background()

	err := r.queries.DeleteWeekDay(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete week day: %w", err)
	}
	return nil
}

// DeleteWeekDayByWeekAndDayAndDayOfWeek removes a day from a week by week ID, day ID, and day of week.
func (r *WeekRepository) DeleteWeekDayByWeekAndDayAndDayOfWeek(weekID, dayID, dayOfWeek string) error {
	ctx := context.Background()

	err := r.queries.DeleteWeekDayByWeekAndDay(ctx, db.DeleteWeekDayByWeekAndDayParams{
		WeekID:    weekID,
		DayID:     dayID,
		DayOfWeek: dayOfWeek,
	})
	if err != nil {
		return fmt.Errorf("failed to delete week day: %w", err)
	}
	return nil
}

// CountWeekDays counts the number of days in a week.
func (r *WeekRepository) CountWeekDays(weekID string) (int64, error) {
	ctx := context.Background()
	count, err := r.queries.CountWeekDays(ctx, weekID)
	if err != nil {
		return 0, fmt.Errorf("failed to count week days: %w", err)
	}
	return count, nil
}

// DayExists checks if a day exists.
func (r *WeekRepository) DayExists(id string) (bool, error) {
	ctx := context.Background()

	_, err := r.queries.GetDay(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if day exists: %w", err)
	}
	return true, nil
}

// Helper functions

func dbWeekToDomain(dbWeek db.Week) *week.Week {
	createdAt, _ := time.Parse(time.RFC3339, dbWeek.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, dbWeek.UpdatedAt)

	return &week.Week{
		ID:         dbWeek.ID,
		WeekNumber: int(dbWeek.WeekNumber),
		Variant:    nullStringToStringPtr(dbWeek.Variant),
		CycleID:    dbWeek.CycleID,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}

func dbWeekDayToDomain(dbWeekDay db.WeekDay) *week.WeekDay {
	createdAt, _ := time.Parse(time.RFC3339, dbWeekDay.CreatedAt)

	return &week.WeekDay{
		ID:        dbWeekDay.ID,
		WeekID:    dbWeekDay.WeekID,
		DayID:     dbWeekDay.DayID,
		DayOfWeek: week.DayOfWeek(dbWeekDay.DayOfWeek),
		CreatedAt: createdAt,
	}
}

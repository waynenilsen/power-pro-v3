// Package repository provides database repository implementations.
package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/db"
	"github.com/waynenilsen/power-pro-v3/internal/domain/day"
)

// DayRepository implements day persistence using sqlc-generated queries.
type DayRepository struct {
	queries *db.Queries
}

// NewDayRepository creates a new DayRepository.
func NewDayRepository(sqlDB *sql.DB) *DayRepository {
	return &DayRepository{
		queries: db.New(sqlDB),
	}
}

// DaySortField represents a field to sort by.
type DaySortField string

const (
	DaySortByName      DaySortField = "name"
	DaySortByCreatedAt DaySortField = "created_at"
)

// DayListParams contains parameters for listing days.
type DayListParams struct {
	Limit           int64
	Offset          int64
	SortBy          DaySortField
	SortOrder       SortOrder
	FilterProgramID *string
}

// GetByID retrieves a day by its ID.
func (r *DayRepository) GetByID(id string) (*day.Day, error) {
	ctx := context.Background()
	dbDay, err := r.queries.GetDay(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get day: %w", err)
	}
	return dbDayToDomain(dbDay), nil
}

// GetBySlug retrieves a day by its slug within a program.
func (r *DayRepository) GetBySlug(slug string, programID *string) (*day.Day, error) {
	ctx := context.Background()
	dbDay, err := r.queries.GetDayBySlug(ctx, db.GetDayBySlugParams{
		Slug:      slug,
		ProgramID: stringPtrToNullString(programID),
		Column3:   stringPtrToNullString(programID),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get day by slug: %w", err)
	}
	return dbDayToDomain(dbDay), nil
}

// List retrieves days with pagination, sorting, and optional filtering.
func (r *DayRepository) List(params DayListParams) ([]day.Day, int64, error) {
	ctx := context.Background()

	// Set defaults
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = DaySortByName
	}
	if params.SortOrder == "" {
		params.SortOrder = SortAsc
	}

	var dbDays []db.Day
	var total int64
	var err error

	if params.FilterProgramID != nil {
		// Filter by program
		programID := sql.NullString{String: *params.FilterProgramID, Valid: true}

		total, err = r.queries.CountDaysFilteredByProgram(ctx, programID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count days: %w", err)
		}

		// Select appropriate query based on sort
		switch {
		case params.SortBy == DaySortByName && params.SortOrder == SortAsc:
			dbDays, err = r.queries.ListDaysFilteredByProgramByNameAsc(ctx, db.ListDaysFilteredByProgramByNameAscParams{
				ProgramID: programID,
				Limit:     params.Limit,
				Offset:    params.Offset,
			})
		case params.SortBy == DaySortByName && params.SortOrder == SortDesc:
			dbDays, err = r.queries.ListDaysFilteredByProgramByNameDesc(ctx, db.ListDaysFilteredByProgramByNameDescParams{
				ProgramID: programID,
				Limit:     params.Limit,
				Offset:    params.Offset,
			})
		case params.SortBy == DaySortByCreatedAt && params.SortOrder == SortAsc:
			dbDays, err = r.queries.ListDaysFilteredByProgramByCreatedAtAsc(ctx, db.ListDaysFilteredByProgramByCreatedAtAscParams{
				ProgramID: programID,
				Limit:     params.Limit,
				Offset:    params.Offset,
			})
		case params.SortBy == DaySortByCreatedAt && params.SortOrder == SortDesc:
			dbDays, err = r.queries.ListDaysFilteredByProgramByCreatedAtDesc(ctx, db.ListDaysFilteredByProgramByCreatedAtDescParams{
				ProgramID: programID,
				Limit:     params.Limit,
				Offset:    params.Offset,
			})
		default:
			dbDays, err = r.queries.ListDaysFilteredByProgramByNameAsc(ctx, db.ListDaysFilteredByProgramByNameAscParams{
				ProgramID: programID,
				Limit:     params.Limit,
				Offset:    params.Offset,
			})
		}
	} else {
		// No filter
		total, err = r.queries.CountDays(ctx)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count days: %w", err)
		}

		// Select appropriate query based on sort
		switch {
		case params.SortBy == DaySortByName && params.SortOrder == SortAsc:
			dbDays, err = r.queries.ListDaysByNameAsc(ctx, db.ListDaysByNameAscParams{
				Limit:  params.Limit,
				Offset: params.Offset,
			})
		case params.SortBy == DaySortByName && params.SortOrder == SortDesc:
			dbDays, err = r.queries.ListDaysByNameDesc(ctx, db.ListDaysByNameDescParams{
				Limit:  params.Limit,
				Offset: params.Offset,
			})
		case params.SortBy == DaySortByCreatedAt && params.SortOrder == SortAsc:
			dbDays, err = r.queries.ListDaysByCreatedAtAsc(ctx, db.ListDaysByCreatedAtAscParams{
				Limit:  params.Limit,
				Offset: params.Offset,
			})
		case params.SortBy == DaySortByCreatedAt && params.SortOrder == SortDesc:
			dbDays, err = r.queries.ListDaysByCreatedAtDesc(ctx, db.ListDaysByCreatedAtDescParams{
				Limit:  params.Limit,
				Offset: params.Offset,
			})
		default:
			dbDays, err = r.queries.ListDaysByNameAsc(ctx, db.ListDaysByNameAscParams{
				Limit:  params.Limit,
				Offset: params.Offset,
			})
		}
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list days: %w", err)
	}

	days := make([]day.Day, len(dbDays))
	for i, dbDay := range dbDays {
		days[i] = *dbDayToDomain(dbDay)
	}

	return days, total, nil
}

// SlugExists checks if a slug already exists within a program, excluding a specific ID.
func (r *DayRepository) SlugExists(slug string, programID *string, excludeID *string) (bool, error) {
	ctx := context.Background()

	if programID == nil {
		// Check for null program ID
		if excludeID == nil {
			exists, err := r.queries.DaySlugExistsNullProgramForNew(ctx, slug)
			if err != nil {
				return false, fmt.Errorf("failed to check slug exists: %w", err)
			}
			return exists == 1, nil
		}
		exists, err := r.queries.DaySlugExistsNullProgram(ctx, db.DaySlugExistsNullProgramParams{
			Slug: slug,
			ID:   *excludeID,
		})
		if err != nil {
			return false, fmt.Errorf("failed to check slug exists: %w", err)
		}
		return exists == 1, nil
	}

	// Check for non-null program ID
	programIDStr := sql.NullString{String: *programID, Valid: true}
	if excludeID == nil {
		exists, err := r.queries.DaySlugExistsForNew(ctx, db.DaySlugExistsForNewParams{
			Slug:      slug,
			ProgramID: programIDStr,
		})
		if err != nil {
			return false, fmt.Errorf("failed to check slug exists: %w", err)
		}
		return exists == 1, nil
	}

	exists, err := r.queries.DaySlugExists(ctx, db.DaySlugExistsParams{
		Slug:      slug,
		ProgramID: programIDStr,
		ID:        *excludeID,
	})
	if err != nil {
		return false, fmt.Errorf("failed to check slug exists: %w", err)
	}
	return exists == 1, nil
}

// Create persists a new day to the database.
func (r *DayRepository) Create(d *day.Day) error {
	ctx := context.Background()

	metadata := sql.NullString{}
	if d.Metadata != nil {
		metadataJSON, err := json.Marshal(d.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
		metadata = sql.NullString{String: string(metadataJSON), Valid: true}
	}

	err := r.queries.CreateDay(ctx, db.CreateDayParams{
		ID:        d.ID,
		Name:      d.Name,
		Slug:      d.Slug,
		Metadata:  metadata,
		ProgramID: stringPtrToNullString(d.ProgramID),
		CreatedAt: d.CreatedAt.Format(time.RFC3339),
		UpdatedAt: d.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to create day: %w", err)
	}
	return nil
}

// Update persists changes to an existing day.
func (r *DayRepository) Update(d *day.Day) error {
	ctx := context.Background()

	metadata := sql.NullString{}
	if d.Metadata != nil {
		metadataJSON, err := json.Marshal(d.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
		metadata = sql.NullString{String: string(metadataJSON), Valid: true}
	}

	err := r.queries.UpdateDay(ctx, db.UpdateDayParams{
		ID:        d.ID,
		Name:      d.Name,
		Slug:      d.Slug,
		Metadata:  metadata,
		ProgramID: stringPtrToNullString(d.ProgramID),
		UpdatedAt: d.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to update day: %w", err)
	}
	return nil
}

// Delete removes a day from the database.
func (r *DayRepository) Delete(id string) error {
	ctx := context.Background()

	err := r.queries.DeleteDay(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete day: %w", err)
	}
	return nil
}

// IsUsedInWeeks checks if a day is used in any weeks.
func (r *DayRepository) IsUsedInWeeks(id string) (bool, error) {
	ctx := context.Background()

	isUsed, err := r.queries.DayIsUsedInWeeks(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to check if day is used in weeks: %w", err)
	}
	return isUsed == 1, nil
}

// Day Prescription methods

// GetDayPrescription retrieves a day prescription by its ID.
func (r *DayRepository) GetDayPrescription(id string) (*day.DayPrescription, error) {
	ctx := context.Background()
	dbDayPrescription, err := r.queries.GetDayPrescription(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get day prescription: %w", err)
	}
	return dbDayPrescriptionToDomain(dbDayPrescription), nil
}

// GetDayPrescriptionByDayAndPrescription retrieves a day prescription by day ID and prescription ID.
func (r *DayRepository) GetDayPrescriptionByDayAndPrescription(dayID, prescriptionID string) (*day.DayPrescription, error) {
	ctx := context.Background()
	dbDayPrescription, err := r.queries.GetDayPrescriptionByDayAndPrescription(ctx, db.GetDayPrescriptionByDayAndPrescriptionParams{
		DayID:          dayID,
		PrescriptionID: prescriptionID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get day prescription: %w", err)
	}
	return dbDayPrescriptionToDomain(dbDayPrescription), nil
}

// ListDayPrescriptions retrieves all prescriptions for a day ordered by their order field.
func (r *DayRepository) ListDayPrescriptions(dayID string) ([]day.DayPrescription, error) {
	ctx := context.Background()
	dbDayPrescriptions, err := r.queries.ListDayPrescriptions(ctx, dayID)
	if err != nil {
		return nil, fmt.Errorf("failed to list day prescriptions: %w", err)
	}

	prescriptions := make([]day.DayPrescription, len(dbDayPrescriptions))
	for i, dbDayPrescription := range dbDayPrescriptions {
		prescriptions[i] = *dbDayPrescriptionToDomain(dbDayPrescription)
	}

	return prescriptions, nil
}

// GetMaxDayPrescriptionOrder retrieves the maximum order value for prescriptions in a day.
func (r *DayRepository) GetMaxDayPrescriptionOrder(dayID string) (int, error) {
	ctx := context.Background()
	maxOrder, err := r.queries.GetMaxDayPrescriptionOrder(ctx, dayID)
	if err != nil {
		return -1, fmt.Errorf("failed to get max day prescription order: %w", err)
	}
	// The result is interface{} due to COALESCE, need to convert
	switch v := maxOrder.(type) {
	case int64:
		return int(v), nil
	case int:
		return v, nil
	default:
		return -1, nil
	}
}

// CreateDayPrescription adds a prescription to a day.
func (r *DayRepository) CreateDayPrescription(dp *day.DayPrescription) error {
	ctx := context.Background()

	err := r.queries.CreateDayPrescription(ctx, db.CreateDayPrescriptionParams{
		ID:             dp.ID,
		DayID:          dp.DayID,
		PrescriptionID: dp.PrescriptionID,
		Order:          int64(dp.Order),
		CreatedAt:      dp.CreatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to create day prescription: %w", err)
	}
	return nil
}

// DeleteDayPrescription removes a prescription from a day by its ID.
func (r *DayRepository) DeleteDayPrescription(id string) error {
	ctx := context.Background()

	err := r.queries.DeleteDayPrescription(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete day prescription: %w", err)
	}
	return nil
}

// DeleteDayPrescriptionByDayAndPrescription removes a prescription from a day by day ID and prescription ID.
func (r *DayRepository) DeleteDayPrescriptionByDayAndPrescription(dayID, prescriptionID string) error {
	ctx := context.Background()

	err := r.queries.DeleteDayPrescriptionByDayAndPrescription(ctx, db.DeleteDayPrescriptionByDayAndPrescriptionParams{
		DayID:          dayID,
		PrescriptionID: prescriptionID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete day prescription: %w", err)
	}
	return nil
}

// UpdateDayPrescriptionOrder updates the order of a day prescription.
func (r *DayRepository) UpdateDayPrescriptionOrder(id string, order int) error {
	ctx := context.Background()

	err := r.queries.UpdateDayPrescriptionOrder(ctx, db.UpdateDayPrescriptionOrderParams{
		ID:    id,
		Order: int64(order),
	})
	if err != nil {
		return fmt.Errorf("failed to update day prescription order: %w", err)
	}
	return nil
}

// CountDayPrescriptions counts the number of prescriptions in a day.
func (r *DayRepository) CountDayPrescriptions(dayID string) (int64, error) {
	ctx := context.Background()
	count, err := r.queries.CountDayPrescriptions(ctx, dayID)
	if err != nil {
		return 0, fmt.Errorf("failed to count day prescriptions: %w", err)
	}
	return count, nil
}

// Helper functions

func dbDayToDomain(dbDay db.Day) *day.Day {
	createdAt, _ := time.Parse(time.RFC3339, dbDay.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, dbDay.UpdatedAt)

	var metadata map[string]interface{}
	if dbDay.Metadata.Valid && dbDay.Metadata.String != "" {
		_ = json.Unmarshal([]byte(dbDay.Metadata.String), &metadata)
	}

	return &day.Day{
		ID:        dbDay.ID,
		Name:      dbDay.Name,
		Slug:      dbDay.Slug,
		Metadata:  metadata,
		ProgramID: nullStringToStringPtr(dbDay.ProgramID),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func dbDayPrescriptionToDomain(dbDayPrescription db.DayPrescription) *day.DayPrescription {
	createdAt, _ := time.Parse(time.RFC3339, dbDayPrescription.CreatedAt)

	return &day.DayPrescription{
		ID:             dbDayPrescription.ID,
		DayID:          dbDayPrescription.DayID,
		PrescriptionID: dbDayPrescription.PrescriptionID,
		Order:          int(dbDayPrescription.Order),
		CreatedAt:      createdAt,
	}
}

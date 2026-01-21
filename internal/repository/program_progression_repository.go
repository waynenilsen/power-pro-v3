// Package repository provides database repository implementations.
package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/db"
	"github.com/waynenilsen/power-pro-v3/internal/domain/progression"
)

// ProgramProgressionEntity represents a program-progression configuration as stored in the database.
type ProgramProgressionEntity struct {
	ID                string
	ProgramID         string
	ProgressionID     string
	LiftID            *string
	Priority          int64
	Enabled           bool
	OverrideIncrement *float64
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// ProgramProgressionWithDetails includes progression details joined from the progressions table.
type ProgramProgressionWithDetails struct {
	ProgramProgressionEntity
	ProgressionName       string
	ProgressionType       progression.ProgressionType
	ProgressionParameters json.RawMessage
}

// ProgramProgressionRepository implements CRUD operations for program progressions.
type ProgramProgressionRepository struct {
	queries *db.Queries
}

// NewProgramProgressionRepository creates a new ProgramProgressionRepository.
func NewProgramProgressionRepository(sqlDB *sql.DB) *ProgramProgressionRepository {
	return &ProgramProgressionRepository{
		queries: db.New(sqlDB),
	}
}

// GetByID retrieves a program progression by its ID.
func (r *ProgramProgressionRepository) GetByID(id string) (*ProgramProgressionEntity, error) {
	ctx := context.Background()
	dbPP, err := r.queries.GetProgramProgression(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get program progression: %w", err)
	}
	return dbProgramProgressionToEntity(dbPP), nil
}

// ListByProgram retrieves all program progressions for a program with joined progression details.
func (r *ProgramProgressionRepository) ListByProgram(programID string) ([]ProgramProgressionWithDetails, error) {
	ctx := context.Background()
	rows, err := r.queries.ListProgramProgressionsWithDetailsByProgram(ctx, programID)
	if err != nil {
		return nil, fmt.Errorf("failed to list program progressions: %w", err)
	}

	entities := make([]ProgramProgressionWithDetails, len(rows))
	for i, row := range rows {
		entities[i] = *dbProgramProgressionWithDetailsToEntity(row)
	}
	return entities, nil
}

// ProgramProgressionListParams defines pagination parameters for listing program progressions.
type ProgramProgressionListParams struct {
	ProgramID string
	Limit     int64
	Offset    int64
}

// ListByProgramPaginated retrieves program progressions with pagination.
func (r *ProgramProgressionRepository) ListByProgramPaginated(params ProgramProgressionListParams) ([]ProgramProgressionWithDetails, int64, error) {
	ctx := context.Background()

	// Get total count
	count, err := r.queries.CountProgramProgressionsWithDetailsByProgram(ctx, params.ProgramID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count program progressions: %w", err)
	}

	// Get paginated results
	rows, err := r.queries.ListProgramProgressionsWithDetailsByProgramPaginated(ctx, db.ListProgramProgressionsWithDetailsByProgramPaginatedParams{
		ProgramID: params.ProgramID,
		Limit:     params.Limit,
		Offset:    params.Offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list program progressions: %w", err)
	}

	entities := make([]ProgramProgressionWithDetails, len(rows))
	for i, row := range rows {
		entities[i] = *dbProgramProgressionWithDetailsPaginatedToEntity(row)
	}
	return entities, count, nil
}

// CheckDuplicate checks if a program progression with the same program, progression, and lift already exists.
// Returns the existing entity if found, nil otherwise.
func (r *ProgramProgressionRepository) CheckDuplicate(programID, progressionID string, liftID *string) (*ProgramProgressionEntity, error) {
	ctx := context.Background()

	liftIDParam := sql.NullString{}
	if liftID != nil {
		liftIDParam = sql.NullString{String: *liftID, Valid: true}
	}

	dbPP, err := r.queries.GetProgramProgressionByProgramProgressionLift(ctx, db.GetProgramProgressionByProgramProgressionLiftParams{
		ProgramID:     programID,
		ProgressionID: progressionID,
		LiftID:        liftIDParam,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to check duplicate: %w", err)
	}
	return dbProgramProgressionToEntity(dbPP), nil
}

// Create persists a new program progression to the database.
func (r *ProgramProgressionRepository) Create(entity *ProgramProgressionEntity) error {
	ctx := context.Background()

	enabled := int64(0)
	if entity.Enabled {
		enabled = 1
	}

	liftID := sql.NullString{}
	if entity.LiftID != nil {
		liftID = sql.NullString{String: *entity.LiftID, Valid: true}
	}

	overrideIncrement := sql.NullFloat64{}
	if entity.OverrideIncrement != nil {
		overrideIncrement = sql.NullFloat64{Float64: *entity.OverrideIncrement, Valid: true}
	}

	err := r.queries.CreateProgramProgression(ctx, db.CreateProgramProgressionParams{
		ID:                entity.ID,
		ProgramID:         entity.ProgramID,
		ProgressionID:     entity.ProgressionID,
		LiftID:            liftID,
		Priority:          entity.Priority,
		Enabled:           enabled,
		OverrideIncrement: overrideIncrement,
		CreatedAt:         entity.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         entity.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to create program progression: %w", err)
	}
	return nil
}

// Update persists changes to an existing program progression.
func (r *ProgramProgressionRepository) Update(entity *ProgramProgressionEntity) error {
	ctx := context.Background()

	enabled := int64(0)
	if entity.Enabled {
		enabled = 1
	}

	overrideIncrement := sql.NullFloat64{}
	if entity.OverrideIncrement != nil {
		overrideIncrement = sql.NullFloat64{Float64: *entity.OverrideIncrement, Valid: true}
	}

	err := r.queries.UpdateProgramProgression(ctx, db.UpdateProgramProgressionParams{
		ID:                entity.ID,
		Priority:          entity.Priority,
		Enabled:           enabled,
		OverrideIncrement: overrideIncrement,
		UpdatedAt:         entity.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to update program progression: %w", err)
	}
	return nil
}

// Delete removes a program progression from the database.
func (r *ProgramProgressionRepository) Delete(id string) error {
	ctx := context.Background()

	err := r.queries.DeleteProgramProgression(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete program progression: %w", err)
	}
	return nil
}

// Helper functions

func dbProgramProgressionToEntity(dbPP db.ProgramProgression) *ProgramProgressionEntity {
	createdAt, _ := time.Parse(time.RFC3339, dbPP.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, dbPP.UpdatedAt)

	var liftID *string
	if dbPP.LiftID.Valid {
		liftID = &dbPP.LiftID.String
	}

	var overrideIncrement *float64
	if dbPP.OverrideIncrement.Valid {
		overrideIncrement = &dbPP.OverrideIncrement.Float64
	}

	return &ProgramProgressionEntity{
		ID:                dbPP.ID,
		ProgramID:         dbPP.ProgramID,
		ProgressionID:     dbPP.ProgressionID,
		LiftID:            liftID,
		Priority:          dbPP.Priority,
		Enabled:           dbPP.Enabled == 1,
		OverrideIncrement: overrideIncrement,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}
}

func dbProgramProgressionWithDetailsToEntity(row db.ListProgramProgressionsWithDetailsByProgramRow) *ProgramProgressionWithDetails {
	createdAt, _ := time.Parse(time.RFC3339, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, row.UpdatedAt)

	var liftID *string
	if row.LiftID.Valid {
		liftID = &row.LiftID.String
	}

	var overrideIncrement *float64
	if row.OverrideIncrement.Valid {
		overrideIncrement = &row.OverrideIncrement.Float64
	}

	return &ProgramProgressionWithDetails{
		ProgramProgressionEntity: ProgramProgressionEntity{
			ID:                row.ID,
			ProgramID:         row.ProgramID,
			ProgressionID:     row.ProgressionID,
			LiftID:            liftID,
			Priority:          row.Priority,
			Enabled:           row.Enabled == 1,
			OverrideIncrement: overrideIncrement,
			CreatedAt:         createdAt,
			UpdatedAt:         updatedAt,
		},
		ProgressionName:       row.ProgressionName,
		ProgressionType:       progression.ProgressionType(row.ProgressionType),
		ProgressionParameters: json.RawMessage(row.ProgressionParameters),
	}
}

func dbProgramProgressionWithDetailsPaginatedToEntity(row db.ListProgramProgressionsWithDetailsByProgramPaginatedRow) *ProgramProgressionWithDetails {
	createdAt, _ := time.Parse(time.RFC3339, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, row.UpdatedAt)

	var liftID *string
	if row.LiftID.Valid {
		liftID = &row.LiftID.String
	}

	var overrideIncrement *float64
	if row.OverrideIncrement.Valid {
		overrideIncrement = &row.OverrideIncrement.Float64
	}

	return &ProgramProgressionWithDetails{
		ProgramProgressionEntity: ProgramProgressionEntity{
			ID:                row.ID,
			ProgramID:         row.ProgramID,
			ProgressionID:     row.ProgressionID,
			LiftID:            liftID,
			Priority:          row.Priority,
			Enabled:           row.Enabled == 1,
			OverrideIncrement: overrideIncrement,
			CreatedAt:         createdAt,
			UpdatedAt:         updatedAt,
		},
		ProgressionName:       row.ProgressionName,
		ProgressionType:       progression.ProgressionType(row.ProgressionType),
		ProgressionParameters: json.RawMessage(row.ProgressionParameters),
	}
}

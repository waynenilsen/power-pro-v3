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

// ProgressionEntity represents a progression as stored in the database.
// This is a data transfer object that mirrors the db.Progression model
// but uses domain types for easier handling.
type ProgressionEntity struct {
	ID         string
	Name       string
	Type       progression.ProgressionType
	Parameters json.RawMessage
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ProgressionListParams contains parameters for listing progressions.
type ProgressionListParams struct {
	Limit      int64
	Offset     int64
	FilterType *progression.ProgressionType
}

// ProgressionRepository implements CRUD operations for progressions.
type ProgressionRepository struct {
	queries *db.Queries
}

// NewProgressionRepository creates a new ProgressionRepository.
func NewProgressionRepository(sqlDB *sql.DB) *ProgressionRepository {
	return &ProgressionRepository{
		queries: db.New(sqlDB),
	}
}

// GetByID retrieves a progression by its ID.
func (r *ProgressionRepository) GetByID(id string) (*ProgressionEntity, error) {
	ctx := context.Background()
	dbProg, err := r.queries.GetProgression(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get progression: %w", err)
	}
	return dbProgressionToEntity(dbProg), nil
}

// List retrieves progressions with pagination and optional filtering.
func (r *ProgressionRepository) List(params ProgressionListParams) ([]ProgressionEntity, int64, error) {
	ctx := context.Background()

	// Set defaults
	if params.Limit <= 0 {
		params.Limit = 20
	}

	var dbProgs []db.Progression
	var total int64
	var err error

	if params.FilterType != nil {
		// Filtered by type
		total, err = r.queries.CountProgressionsByType(ctx, string(*params.FilterType))
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count progressions: %w", err)
		}

		dbProgs, err = r.queries.ListProgressionsByType(ctx, db.ListProgressionsByTypeParams{
			Type:   string(*params.FilterType),
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	} else {
		// No filter
		total, err = r.queries.CountProgressions(ctx)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count progressions: %w", err)
		}

		dbProgs, err = r.queries.ListProgressions(ctx, db.ListProgressionsParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list progressions: %w", err)
	}

	entities := make([]ProgressionEntity, len(dbProgs))
	for i, dbProg := range dbProgs {
		entities[i] = *dbProgressionToEntity(dbProg)
	}

	return entities, total, nil
}

// Create persists a new progression to the database.
func (r *ProgressionRepository) Create(entity *ProgressionEntity) error {
	ctx := context.Background()

	err := r.queries.CreateProgression(ctx, db.CreateProgressionParams{
		ID:         entity.ID,
		Name:       entity.Name,
		Type:       string(entity.Type),
		Parameters: string(entity.Parameters),
		CreatedAt:  entity.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  entity.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to create progression: %w", err)
	}
	return nil
}

// Update persists changes to an existing progression.
func (r *ProgressionRepository) Update(entity *ProgressionEntity) error {
	ctx := context.Background()

	err := r.queries.UpdateProgression(ctx, db.UpdateProgressionParams{
		ID:         entity.ID,
		Name:       entity.Name,
		Type:       string(entity.Type),
		Parameters: string(entity.Parameters),
		UpdatedAt:  entity.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to update progression: %w", err)
	}
	return nil
}

// Delete removes a progression from the database.
func (r *ProgressionRepository) Delete(id string) error {
	ctx := context.Background()

	err := r.queries.DeleteProgression(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete progression: %w", err)
	}
	return nil
}

// HasProgramReferences checks if a progression is referenced by any program_progressions.
func (r *ProgressionRepository) HasProgramReferences(id string) (bool, error) {
	ctx := context.Background()

	count, err := r.queries.CountProgramProgressionsByProgression(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to check references: %w", err)
	}
	return count > 0, nil
}

// Helper functions

func dbProgressionToEntity(dbProg db.Progression) *ProgressionEntity {
	createdAt, _ := time.Parse(time.RFC3339, dbProg.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, dbProg.UpdatedAt)

	return &ProgressionEntity{
		ID:         dbProg.ID,
		Name:       dbProg.Name,
		Type:       progression.ProgressionType(dbProg.Type),
		Parameters: json.RawMessage(dbProg.Parameters),
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}

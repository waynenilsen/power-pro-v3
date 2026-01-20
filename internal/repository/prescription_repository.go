// Package repository provides database repository implementations.
package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/db"
	"github.com/waynenilsen/power-pro-v3/internal/domain/loadstrategy"
	"github.com/waynenilsen/power-pro-v3/internal/domain/prescription"
	"github.com/waynenilsen/power-pro-v3/internal/domain/setscheme"
)

// PrescriptionRepository implements prescription persistence using sqlc-generated queries.
type PrescriptionRepository struct {
	queries         *db.Queries
	strategyFactory *loadstrategy.StrategyFactory
	schemeFactory   *setscheme.SchemeFactory
}

// NewPrescriptionRepository creates a new PrescriptionRepository.
func NewPrescriptionRepository(sqlDB *sql.DB, strategyFactory *loadstrategy.StrategyFactory, schemeFactory *setscheme.SchemeFactory) *PrescriptionRepository {
	return &PrescriptionRepository{
		queries:         db.New(sqlDB),
		strategyFactory: strategyFactory,
		schemeFactory:   schemeFactory,
	}
}

// PrescriptionSortField represents a field to sort by.
type PrescriptionSortField string

const (
	PrescriptionSortByOrder     PrescriptionSortField = "order"
	PrescriptionSortByCreatedAt PrescriptionSortField = "created_at"
)

// PrescriptionListParams contains parameters for listing prescriptions.
type PrescriptionListParams struct {
	Limit        int64
	Offset       int64
	SortBy       PrescriptionSortField
	SortOrder    SortOrder
	FilterLiftID *string
}

// GetByID retrieves a prescription by its ID.
func (r *PrescriptionRepository) GetByID(id string) (*prescription.Prescription, error) {
	ctx := context.Background()
	dbPrescription, err := r.queries.GetPrescription(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get prescription: %w", err)
	}
	return r.dbPrescriptionToDomain(dbPrescription)
}

// List retrieves prescriptions with pagination, sorting, and optional filtering.
func (r *PrescriptionRepository) List(params PrescriptionListParams) ([]prescription.Prescription, int64, error) {
	ctx := context.Background()

	// Set defaults
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = PrescriptionSortByOrder
	}
	if params.SortOrder == "" {
		params.SortOrder = SortAsc
	}

	var dbPrescriptions []db.Prescription
	var total int64
	var err error

	if params.FilterLiftID != nil {
		// Filter by lift ID
		total, err = r.queries.CountPrescriptionsFilterLift(ctx, *params.FilterLiftID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count prescriptions: %w", err)
		}

		// Select appropriate query based on sort
		switch {
		case params.SortBy == PrescriptionSortByOrder && params.SortOrder == SortAsc:
			dbPrescriptions, err = r.queries.ListPrescriptionsFilterLiftByOrderAsc(ctx, db.ListPrescriptionsFilterLiftByOrderAscParams{
				LiftID: *params.FilterLiftID,
				Limit:  params.Limit,
				Offset: params.Offset,
			})
		case params.SortBy == PrescriptionSortByOrder && params.SortOrder == SortDesc:
			dbPrescriptions, err = r.queries.ListPrescriptionsFilterLiftByOrderDesc(ctx, db.ListPrescriptionsFilterLiftByOrderDescParams{
				LiftID: *params.FilterLiftID,
				Limit:  params.Limit,
				Offset: params.Offset,
			})
		case params.SortBy == PrescriptionSortByCreatedAt && params.SortOrder == SortAsc:
			dbPrescriptions, err = r.queries.ListPrescriptionsFilterLiftByCreatedAtAsc(ctx, db.ListPrescriptionsFilterLiftByCreatedAtAscParams{
				LiftID: *params.FilterLiftID,
				Limit:  params.Limit,
				Offset: params.Offset,
			})
		case params.SortBy == PrescriptionSortByCreatedAt && params.SortOrder == SortDesc:
			dbPrescriptions, err = r.queries.ListPrescriptionsFilterLiftByCreatedAtDesc(ctx, db.ListPrescriptionsFilterLiftByCreatedAtDescParams{
				LiftID: *params.FilterLiftID,
				Limit:  params.Limit,
				Offset: params.Offset,
			})
		default:
			dbPrescriptions, err = r.queries.ListPrescriptionsFilterLiftByOrderAsc(ctx, db.ListPrescriptionsFilterLiftByOrderAscParams{
				LiftID: *params.FilterLiftID,
				Limit:  params.Limit,
				Offset: params.Offset,
			})
		}
	} else {
		// No filter
		total, err = r.queries.CountPrescriptions(ctx)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count prescriptions: %w", err)
		}

		// Select appropriate query based on sort
		switch {
		case params.SortBy == PrescriptionSortByOrder && params.SortOrder == SortAsc:
			dbPrescriptions, err = r.queries.ListPrescriptionsByOrderAsc(ctx, db.ListPrescriptionsByOrderAscParams{
				Limit:  params.Limit,
				Offset: params.Offset,
			})
		case params.SortBy == PrescriptionSortByOrder && params.SortOrder == SortDesc:
			dbPrescriptions, err = r.queries.ListPrescriptionsByOrderDesc(ctx, db.ListPrescriptionsByOrderDescParams{
				Limit:  params.Limit,
				Offset: params.Offset,
			})
		case params.SortBy == PrescriptionSortByCreatedAt && params.SortOrder == SortAsc:
			dbPrescriptions, err = r.queries.ListPrescriptionsByCreatedAtAsc(ctx, db.ListPrescriptionsByCreatedAtAscParams{
				Limit:  params.Limit,
				Offset: params.Offset,
			})
		case params.SortBy == PrescriptionSortByCreatedAt && params.SortOrder == SortDesc:
			dbPrescriptions, err = r.queries.ListPrescriptionsByCreatedAtDesc(ctx, db.ListPrescriptionsByCreatedAtDescParams{
				Limit:  params.Limit,
				Offset: params.Offset,
			})
		default:
			dbPrescriptions, err = r.queries.ListPrescriptionsByOrderAsc(ctx, db.ListPrescriptionsByOrderAscParams{
				Limit:  params.Limit,
				Offset: params.Offset,
			})
		}
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list prescriptions: %w", err)
	}

	prescriptions := make([]prescription.Prescription, 0, len(dbPrescriptions))
	for _, dbPrescription := range dbPrescriptions {
		p, err := r.dbPrescriptionToDomain(dbPrescription)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to convert prescription: %w", err)
		}
		prescriptions = append(prescriptions, *p)
	}

	return prescriptions, total, nil
}

// Create persists a new prescription to the database.
func (r *PrescriptionRepository) Create(p *prescription.Prescription) error {
	ctx := context.Background()

	loadStrategyJSON, err := json.Marshal(p.LoadStrategy)
	if err != nil {
		return fmt.Errorf("failed to marshal load strategy: %w", err)
	}

	setSchemeJSON, err := json.Marshal(p.SetScheme)
	if err != nil {
		return fmt.Errorf("failed to marshal set scheme: %w", err)
	}

	err = r.queries.CreatePrescription(ctx, db.CreatePrescriptionParams{
		ID:           p.ID,
		LiftID:       p.LiftID,
		LoadStrategy: string(loadStrategyJSON),
		SetScheme:    string(setSchemeJSON),
		Order:        int64(p.Order),
		Notes:        stringToNullString(p.Notes),
		RestSeconds:  intPtrToNullInt64(p.RestSeconds),
		CreatedAt:    p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    p.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to create prescription: %w", err)
	}
	return nil
}

// Update persists changes to an existing prescription.
func (r *PrescriptionRepository) Update(p *prescription.Prescription) error {
	ctx := context.Background()

	loadStrategyJSON, err := json.Marshal(p.LoadStrategy)
	if err != nil {
		return fmt.Errorf("failed to marshal load strategy: %w", err)
	}

	setSchemeJSON, err := json.Marshal(p.SetScheme)
	if err != nil {
		return fmt.Errorf("failed to marshal set scheme: %w", err)
	}

	err = r.queries.UpdatePrescription(ctx, db.UpdatePrescriptionParams{
		ID:           p.ID,
		LiftID:       p.LiftID,
		LoadStrategy: string(loadStrategyJSON),
		SetScheme:    string(setSchemeJSON),
		Order:        int64(p.Order),
		Notes:        stringToNullString(p.Notes),
		RestSeconds:  intPtrToNullInt64(p.RestSeconds),
		UpdatedAt:    p.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to update prescription: %w", err)
	}
	return nil
}

// Delete removes a prescription from the database.
func (r *PrescriptionRepository) Delete(id string) error {
	ctx := context.Background()

	err := r.queries.DeletePrescription(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete prescription: %w", err)
	}
	return nil
}

// LiftHasReferences checks if a lift has prescription references.
func (r *PrescriptionRepository) LiftHasReferences(liftID string) (bool, error) {
	ctx := context.Background()

	hasRefs, err := r.queries.LiftHasPrescriptionReferences(ctx, liftID)
	if err != nil {
		return false, fmt.Errorf("failed to check prescription references: %w", err)
	}
	return hasRefs == 1, nil
}

// Helper functions

func (r *PrescriptionRepository) dbPrescriptionToDomain(dbPrescription db.Prescription) (*prescription.Prescription, error) {
	createdAt, _ := time.Parse(time.RFC3339, dbPrescription.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, dbPrescription.UpdatedAt)

	// Unmarshal load strategy
	loadStrategy, err := r.strategyFactory.CreateFromJSON(json.RawMessage(dbPrescription.LoadStrategy))
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal load strategy: %w", err)
	}

	// Unmarshal set scheme
	setScheme, err := r.schemeFactory.CreateFromJSON(json.RawMessage(dbPrescription.SetScheme))
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal set scheme: %w", err)
	}

	return &prescription.Prescription{
		ID:           dbPrescription.ID,
		LiftID:       dbPrescription.LiftID,
		LoadStrategy: loadStrategy,
		SetScheme:    setScheme,
		Order:        int(dbPrescription.Order),
		Notes:        nullStringToString(dbPrescription.Notes),
		RestSeconds:  nullInt64ToIntPtr(dbPrescription.RestSeconds),
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
}

func stringToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func nullStringToString(ns sql.NullString) string {
	if !ns.Valid {
		return ""
	}
	return ns.String
}

func intPtrToNullInt64(i *int) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: int64(*i), Valid: true}
}

func nullInt64ToIntPtr(ni sql.NullInt64) *int {
	if !ni.Valid {
		return nil
	}
	val := int(ni.Int64)
	return &val
}

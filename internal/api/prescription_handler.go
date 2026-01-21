package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/loadstrategy"
	"github.com/waynenilsen/power-pro-v3/internal/domain/prescription"
	"github.com/waynenilsen/power-pro-v3/internal/domain/setscheme"
	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
	"github.com/waynenilsen/power-pro-v3/internal/repository"
)

// PrescriptionHandler handles HTTP requests for prescription operations.
type PrescriptionHandler struct {
	repo            *repository.PrescriptionRepository
	liftRepo        *repository.LiftRepository
	liftMaxRepo     *repository.LiftMaxRepository
	strategyFactory *loadstrategy.StrategyFactory
	schemeFactory   *setscheme.SchemeFactory
}

// NewPrescriptionHandler creates a new PrescriptionHandler.
func NewPrescriptionHandler(
	repo *repository.PrescriptionRepository,
	liftRepo *repository.LiftRepository,
	liftMaxRepo *repository.LiftMaxRepository,
	strategyFactory *loadstrategy.StrategyFactory,
	schemeFactory *setscheme.SchemeFactory,
) *PrescriptionHandler {
	return &PrescriptionHandler{
		repo:            repo,
		liftRepo:        liftRepo,
		liftMaxRepo:     liftMaxRepo,
		strategyFactory: strategyFactory,
		schemeFactory:   schemeFactory,
	}
}

// PrescriptionResponse represents the API response format for a prescription.
type PrescriptionResponse struct {
	ID           string          `json:"id"`
	LiftID       string          `json:"liftId"`
	LoadStrategy json.RawMessage `json:"loadStrategy"`
	SetScheme    json.RawMessage `json:"setScheme"`
	Order        int             `json:"order"`
	Notes        string          `json:"notes,omitempty"`
	RestSeconds  *int            `json:"restSeconds,omitempty"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
}

// CreatePrescriptionRequest represents the request body for creating a prescription.
type CreatePrescriptionRequest struct {
	LiftID       string          `json:"liftId"`
	LoadStrategy json.RawMessage `json:"loadStrategy"`
	SetScheme    json.RawMessage `json:"setScheme"`
	Order        *int            `json:"order,omitempty"`
	Notes        string          `json:"notes,omitempty"`
	RestSeconds  *int            `json:"restSeconds,omitempty"`
}

// UpdatePrescriptionRequest represents the request body for updating a prescription.
type UpdatePrescriptionRequest struct {
	LiftID           *string         `json:"liftId,omitempty"`
	LoadStrategy     json.RawMessage `json:"loadStrategy,omitempty"`
	SetScheme        json.RawMessage `json:"setScheme,omitempty"`
	Order            *int            `json:"order,omitempty"`
	Notes            *string         `json:"notes,omitempty"`
	RestSeconds      *int            `json:"restSeconds,omitempty"`
	ClearRestSeconds bool            `json:"clearRestSeconds,omitempty"`
}

func (h *PrescriptionHandler) prescriptionToResponse(p *prescription.Prescription) (PrescriptionResponse, error) {
	loadStrategyJSON, err := json.Marshal(p.LoadStrategy)
	if err != nil {
		return PrescriptionResponse{}, err
	}

	setSchemeJSON, err := json.Marshal(p.SetScheme)
	if err != nil {
		return PrescriptionResponse{}, err
	}

	return PrescriptionResponse{
		ID:           p.ID,
		LiftID:       p.LiftID,
		LoadStrategy: loadStrategyJSON,
		SetScheme:    setSchemeJSON,
		Order:        p.Order,
		Notes:        p.Notes,
		RestSeconds:  p.RestSeconds,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}, nil
}

// List handles GET /prescriptions
func (h *PrescriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()

	// Pagination
	page := 1
	pageSize := 20
	if p := query.Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if ps := query.Get("pageSize"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	// Sorting
	sortBy := repository.PrescriptionSortByOrder
	sortOrder := repository.SortAsc
	if s := query.Get("sortBy"); s != "" {
		switch strings.ToLower(s) {
		case "order":
			sortBy = repository.PrescriptionSortByOrder
		case "created_at", "createdat":
			sortBy = repository.PrescriptionSortByCreatedAt
		}
	}
	if o := query.Get("sortOrder"); o != "" {
		switch strings.ToLower(o) {
		case "asc":
			sortOrder = repository.SortAsc
		case "desc":
			sortOrder = repository.SortDesc
		}
	}

	// Filter by lift_id
	var filterLiftID *string
	if liftID := query.Get("lift_id"); liftID != "" {
		filterLiftID = &liftID
	}

	params := repository.PrescriptionListParams{
		Limit:        int64(pageSize),
		Offset:       int64((page - 1) * pageSize),
		SortBy:       sortBy,
		SortOrder:    sortOrder,
		FilterLiftID: filterLiftID,
	}

	prescriptions, total, err := h.repo.List(params)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to list prescriptions", err))
		return
	}

	// Convert to response format
	data := make([]PrescriptionResponse, 0, len(prescriptions))
	for _, p := range prescriptions {
		resp, err := h.prescriptionToResponse(&p)
		if err != nil {
			writeDomainError(w, apperrors.NewInternal("failed to format prescription", err))
			return
		}
		data = append(data, resp)
	}

	// Use standard envelope with offset-based pagination
	offset := (page - 1) * pageSize
	writePaginatedData(w, http.StatusOK, data, total, pageSize, offset)
}

// Get handles GET /prescriptions/{id}
func (h *PrescriptionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing prescription ID"))
		return
	}

	p, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get prescription", err))
		return
	}
	if p == nil {
		writeDomainError(w, apperrors.NewNotFound("prescription", id))
		return
	}

	resp, err := h.prescriptionToResponse(p)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to format prescription", err))
		return
	}

	writeData(w, http.StatusOK, resp)
}

// Create handles POST /prescriptions
func (h *PrescriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreatePrescriptionRequest
	if err := readJSON(r, &req); err != nil {
		writeDomainError(w, apperrors.NewBadRequest("invalid request body"))
		return
	}

	// Parse load strategy
	loadStrategy, err := h.strategyFactory.CreateFromJSON(req.LoadStrategy)
	if err != nil {
		writeDomainError(w, apperrors.NewValidation("loadStrategy", err.Error()))
		return
	}

	// Parse set scheme
	setScheme, err := h.schemeFactory.CreateFromJSON(req.SetScheme)
	if err != nil {
		writeDomainError(w, apperrors.NewValidation("setScheme", err.Error()))
		return
	}

	// Default order to 0 if not provided
	order := 0
	if req.Order != nil {
		order = *req.Order
	}

	// Generate UUID
	id := uuid.New().String()

	// Use domain logic to create and validate
	input := prescription.CreatePrescriptionInput{
		LiftID:       req.LiftID,
		LoadStrategy: loadStrategy,
		SetScheme:    setScheme,
		Order:        order,
		Notes:        req.Notes,
		RestSeconds:  req.RestSeconds,
	}

	newPrescription, result := prescription.CreatePrescription(input, id)
	if !result.Valid {
		details := make([]string, len(result.Errors))
		for i, err := range result.Errors {
			details[i] = err.Error()
		}
		writeDomainError(w, apperrors.NewValidationMsg("validation failed"), details...)
		return
	}

	// Check lift exists
	lift, err := h.liftRepo.GetByID(req.LiftID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to verify lift", err))
		return
	}
	if lift == nil {
		writeDomainError(w, apperrors.NewValidation("liftId", "lift not found"))
		return
	}

	// Persist
	if err := h.repo.Create(newPrescription); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to create prescription", err))
		return
	}

	resp, err := h.prescriptionToResponse(newPrescription)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to format prescription", err))
		return
	}

	writeData(w, http.StatusCreated, resp)
}

// Update handles PUT /prescriptions/{id}
func (h *PrescriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing prescription ID"))
		return
	}

	// Get existing prescription
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get prescription", err))
		return
	}
	if existing == nil {
		writeDomainError(w, apperrors.NewNotFound("prescription", id))
		return
	}

	var req UpdatePrescriptionRequest
	if err := readJSON(r, &req); err != nil {
		writeDomainError(w, apperrors.NewBadRequest("invalid request body"))
		return
	}

	// Parse load strategy if provided
	var newLoadStrategy loadstrategy.LoadStrategy
	if len(req.LoadStrategy) > 0 {
		newLoadStrategy, err = h.strategyFactory.CreateFromJSON(req.LoadStrategy)
		if err != nil {
			writeDomainError(w, apperrors.NewValidation("loadStrategy", err.Error()))
			return
		}
	}

	// Parse set scheme if provided
	var newSetScheme setscheme.SetScheme
	if len(req.SetScheme) > 0 {
		newSetScheme, err = h.schemeFactory.CreateFromJSON(req.SetScheme)
		if err != nil {
			writeDomainError(w, apperrors.NewValidation("setScheme", err.Error()))
			return
		}
	}

	// Check lift exists if being updated
	if req.LiftID != nil {
		lift, err := h.liftRepo.GetByID(*req.LiftID)
		if err != nil {
			writeDomainError(w, apperrors.NewInternal("failed to verify lift", err))
			return
		}
		if lift == nil {
			writeDomainError(w, apperrors.NewValidation("liftId", "lift not found"))
			return
		}
	}

	// Use domain logic to update and validate
	input := prescription.UpdatePrescriptionInput{
		LiftID:           req.LiftID,
		LoadStrategy:     newLoadStrategy,
		SetScheme:        newSetScheme,
		Order:            req.Order,
		Notes:            req.Notes,
		RestSeconds:      req.RestSeconds,
		ClearRestSeconds: req.ClearRestSeconds,
	}

	result := prescription.UpdatePrescription(existing, input)
	if !result.Valid {
		details := make([]string, len(result.Errors))
		for i, err := range result.Errors {
			details[i] = err.Error()
		}
		writeDomainError(w, apperrors.NewValidationMsg("validation failed"), details...)
		return
	}

	// Persist
	if err := h.repo.Update(existing); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to update prescription", err))
		return
	}

	resp, err := h.prescriptionToResponse(existing)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to format prescription", err))
		return
	}

	writeData(w, http.StatusOK, resp)
}

// Delete handles DELETE /prescriptions/{id}
func (h *PrescriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing prescription ID"))
		return
	}

	// Check prescription exists
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get prescription", err))
		return
	}
	if existing == nil {
		writeDomainError(w, apperrors.NewNotFound("prescription", id))
		return
	}

	// Delete
	if err := h.repo.Delete(id); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to delete prescription", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

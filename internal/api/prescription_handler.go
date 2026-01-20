package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/loadstrategy"
	"github.com/waynenilsen/power-pro-v3/internal/domain/prescription"
	"github.com/waynenilsen/power-pro-v3/internal/domain/setscheme"
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
		writeError(w, http.StatusInternalServerError, "Failed to list prescriptions")
		return
	}

	// Convert to response format
	data := make([]PrescriptionResponse, 0, len(prescriptions))
	for _, p := range prescriptions {
		resp, err := h.prescriptionToResponse(&p)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to format prescription")
			return
		}
		data = append(data, resp)
	}

	// Calculate total pages
	totalPages := total / int64(pageSize)
	if total%int64(pageSize) > 0 {
		totalPages++
	}

	resp := PaginatedResponse{
		Data:       data,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}

	writeJSON(w, http.StatusOK, resp)
}

// Get handles GET /prescriptions/{id}
func (h *PrescriptionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing prescription ID")
		return
	}

	p, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get prescription")
		return
	}
	if p == nil {
		writeError(w, http.StatusNotFound, "Prescription not found")
		return
	}

	resp, err := h.prescriptionToResponse(p)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to format prescription")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// Create handles POST /prescriptions
func (h *PrescriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreatePrescriptionRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Parse load strategy
	loadStrategy, err := h.strategyFactory.CreateFromJSON(req.LoadStrategy)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid load strategy", err.Error())
		return
	}

	// Parse set scheme
	setScheme, err := h.schemeFactory.CreateFromJSON(req.SetScheme)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid set scheme", err.Error())
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
		writeError(w, http.StatusBadRequest, "Validation failed", details...)
		return
	}

	// Check lift exists
	lift, err := h.liftRepo.GetByID(req.LiftID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to verify lift")
		return
	}
	if lift == nil {
		writeError(w, http.StatusBadRequest, "Lift not found")
		return
	}

	// Persist
	if err := h.repo.Create(newPrescription); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create prescription")
		return
	}

	resp, err := h.prescriptionToResponse(newPrescription)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to format prescription")
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

// Update handles PUT /prescriptions/{id}
func (h *PrescriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing prescription ID")
		return
	}

	// Get existing prescription
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get prescription")
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "Prescription not found")
		return
	}

	var req UpdatePrescriptionRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Parse load strategy if provided
	var newLoadStrategy loadstrategy.LoadStrategy
	if len(req.LoadStrategy) > 0 {
		newLoadStrategy, err = h.strategyFactory.CreateFromJSON(req.LoadStrategy)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid load strategy", err.Error())
			return
		}
	}

	// Parse set scheme if provided
	var newSetScheme setscheme.SetScheme
	if len(req.SetScheme) > 0 {
		newSetScheme, err = h.schemeFactory.CreateFromJSON(req.SetScheme)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid set scheme", err.Error())
			return
		}
	}

	// Check lift exists if being updated
	if req.LiftID != nil {
		lift, err := h.liftRepo.GetByID(*req.LiftID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to verify lift")
			return
		}
		if lift == nil {
			writeError(w, http.StatusBadRequest, "Lift not found")
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
		writeError(w, http.StatusBadRequest, "Validation failed", details...)
		return
	}

	// Persist
	if err := h.repo.Update(existing); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update prescription")
		return
	}

	resp, err := h.prescriptionToResponse(existing)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to format prescription")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// Delete handles DELETE /prescriptions/{id}
func (h *PrescriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing prescription ID")
		return
	}

	// Check prescription exists
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get prescription")
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "Prescription not found")
		return
	}

	// Delete
	if err := h.repo.Delete(id); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete prescription")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ResolveRequest represents the request body for resolving a prescription.
type ResolveRequest struct {
	UserID string `json:"userId"`
}

// ResolvedPrescriptionResponse represents the API response for a resolved prescription.
type ResolvedPrescriptionResponse struct {
	PrescriptionID string                   `json:"prescriptionId"`
	Lift           prescription.LiftInfo    `json:"lift"`
	Sets           []setscheme.GeneratedSet `json:"sets"`
	Notes          string                   `json:"notes,omitempty"`
	RestSeconds    *int                     `json:"restSeconds,omitempty"`
}

// BatchResolveRequest represents the request body for batch resolving prescriptions.
type BatchResolveRequest struct {
	PrescriptionIDs []string `json:"prescriptionIds"`
	UserID          string   `json:"userId"`
}

// BatchResolveResultItem represents a single item in the batch resolution response.
type BatchResolveResultItem struct {
	PrescriptionID string                        `json:"prescriptionId"`
	Status         string                        `json:"status"`
	Resolved       *ResolvedPrescriptionResponse `json:"resolved,omitempty"`
	Error          string                        `json:"error,omitempty"`
}

// BatchResolveResponse represents the response for batch resolution.
type BatchResolveResponse struct {
	Results []BatchResolveResultItem `json:"results"`
}

// liftLookupAdapter adapts LiftRepository to prescription.LiftLookup interface.
type liftLookupAdapter struct {
	repo *repository.LiftRepository
}

// GetLiftByID implements prescription.LiftLookup.
func (a *liftLookupAdapter) GetLiftByID(ctx context.Context, liftID string) (*prescription.LiftInfo, error) {
	lift, err := a.repo.GetByID(liftID)
	if err != nil {
		return nil, err
	}
	if lift == nil {
		return nil, nil
	}
	return &prescription.LiftInfo{
		ID:   lift.ID,
		Name: lift.Name,
		Slug: lift.Slug,
	}, nil
}

// maxLookupAdapter adapts LiftMaxRepository to loadstrategy.MaxLookup interface.
type maxLookupAdapter struct {
	repo *repository.LiftMaxRepository
}

// GetCurrentMax implements loadstrategy.MaxLookup.
func (a *maxLookupAdapter) GetCurrentMax(ctx context.Context, userID, liftID, maxType string) (*loadstrategy.MaxValue, error) {
	max, err := a.repo.GetCurrentMax(userID, liftID, maxType)
	if err != nil {
		return nil, err
	}
	if max == nil {
		return nil, nil
	}
	return &loadstrategy.MaxValue{
		Value:         max.Value,
		EffectiveDate: max.EffectiveDate.Format(time.RFC3339),
	}, nil
}

// cachedMaxLookup wraps a maxLookupAdapter with per-request caching.
type cachedMaxLookup struct {
	underlying loadstrategy.MaxLookup
	cache      map[string]*loadstrategy.MaxValue
	mu         sync.RWMutex
}

// newCachedMaxLookup creates a new cachedMaxLookup.
func newCachedMaxLookup(underlying loadstrategy.MaxLookup) *cachedMaxLookup {
	return &cachedMaxLookup{
		underlying: underlying,
		cache:      make(map[string]*loadstrategy.MaxValue),
	}
}

// GetCurrentMax implements loadstrategy.MaxLookup with caching.
func (c *cachedMaxLookup) GetCurrentMax(ctx context.Context, userID, liftID, maxType string) (*loadstrategy.MaxValue, error) {
	key := userID + "|" + liftID + "|" + maxType

	// Check cache first
	c.mu.RLock()
	if val, ok := c.cache[key]; ok {
		c.mu.RUnlock()
		return val, nil
	}
	c.mu.RUnlock()

	// Fetch from underlying
	val, err := c.underlying.GetCurrentMax(ctx, userID, liftID, maxType)
	if err != nil {
		return nil, err
	}

	// Store in cache (including nil results)
	c.mu.Lock()
	c.cache[key] = val
	c.mu.Unlock()

	return val, nil
}

// Resolve handles POST /prescriptions/{id}/resolve
func (h *PrescriptionHandler) Resolve(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing prescription ID")
		return
	}

	var req ResolveRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.UserID) == "" {
		writeError(w, http.StatusBadRequest, "userId is required")
		return
	}

	// Fetch prescription
	p, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get prescription")
		return
	}
	if p == nil {
		writeError(w, http.StatusNotFound, "Prescription not found")
		return
	}

	// Set up resolution context
	liftLookup := &liftLookupAdapter{repo: h.liftRepo}
	maxLookup := &maxLookupAdapter{repo: h.liftMaxRepo}

	// Inject MaxLookup into load strategy
	h.injectMaxLookup(p.LoadStrategy, maxLookup)

	resCtx := prescription.DefaultResolutionContext(liftLookup)

	// Resolve
	ctx := r.Context()
	resolved, err := p.Resolve(ctx, req.UserID, resCtx)
	if err != nil {
		// Check for specific error types
		if errors.Is(err, prescription.ErrLiftNotFound) {
			writeError(w, http.StatusNotFound, "Lift not found")
			return
		}
		if errors.Is(err, prescription.ErrMaxNotFound) || errors.Is(err, loadstrategy.ErrMaxNotFound) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to resolve prescription")
		return
	}

	resp := ResolvedPrescriptionResponse{
		PrescriptionID: resolved.PrescriptionID,
		Lift:           resolved.Lift,
		Sets:           resolved.Sets,
		Notes:          resolved.Notes,
		RestSeconds:    resolved.RestSeconds,
	}

	writeJSON(w, http.StatusOK, resp)
}

// ResolveBatch handles POST /prescriptions/resolve-batch
func (h *PrescriptionHandler) ResolveBatch(w http.ResponseWriter, r *http.Request) {
	var req BatchResolveRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.UserID) == "" {
		writeError(w, http.StatusBadRequest, "userId is required")
		return
	}

	if len(req.PrescriptionIDs) == 0 {
		writeError(w, http.StatusBadRequest, "prescriptionIds is required")
		return
	}

	// Set up resolution context with cached max lookup
	liftLookup := &liftLookupAdapter{repo: h.liftRepo}
	underlyingMaxLookup := &maxLookupAdapter{repo: h.liftMaxRepo}
	cachedMaxLookup := newCachedMaxLookup(underlyingMaxLookup)

	resCtx := prescription.DefaultResolutionContext(liftLookup)
	ctx := r.Context()

	results := make([]BatchResolveResultItem, len(req.PrescriptionIDs))

	for i, prescriptionID := range req.PrescriptionIDs {
		result := BatchResolveResultItem{
			PrescriptionID: prescriptionID,
		}

		// Fetch prescription
		p, err := h.repo.GetByID(prescriptionID)
		if err != nil {
			result.Status = "error"
			result.Error = "Failed to get prescription"
			results[i] = result
			continue
		}
		if p == nil {
			result.Status = "error"
			result.Error = "Prescription not found"
			results[i] = result
			continue
		}

		// Inject cached MaxLookup into load strategy
		h.injectMaxLookup(p.LoadStrategy, cachedMaxLookup)

		// Resolve
		resolved, err := p.Resolve(ctx, req.UserID, resCtx)
		if err != nil {
			result.Status = "error"
			if errors.Is(err, prescription.ErrMaxNotFound) || errors.Is(err, loadstrategy.ErrMaxNotFound) {
				result.Error = err.Error()
			} else if errors.Is(err, prescription.ErrLiftNotFound) {
				result.Error = "Lift not found"
			} else {
				result.Error = "Failed to resolve prescription"
			}
			results[i] = result
			continue
		}

		result.Status = "success"
		result.Resolved = &ResolvedPrescriptionResponse{
			PrescriptionID: resolved.PrescriptionID,
			Lift:           resolved.Lift,
			Sets:           resolved.Sets,
			Notes:          resolved.Notes,
			RestSeconds:    resolved.RestSeconds,
		}
		results[i] = result
	}

	resp := BatchResolveResponse{
		Results: results,
	}

	writeJSON(w, http.StatusOK, resp)
}

// injectMaxLookup injects a MaxLookup into a LoadStrategy if it supports it.
func (h *PrescriptionHandler) injectMaxLookup(strategy loadstrategy.LoadStrategy, maxLookup loadstrategy.MaxLookup) {
	// Check if strategy has a SetMaxLookup method (like PercentOfLoadStrategy)
	if setter, ok := strategy.(interface{ SetMaxLookup(loadstrategy.MaxLookup) }); ok {
		setter.SetMaxLookup(maxLookup)
	}
}

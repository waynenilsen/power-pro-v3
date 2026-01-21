package api

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/domain/loadstrategy"
	"github.com/waynenilsen/power-pro-v3/internal/domain/prescription"
	"github.com/waynenilsen/power-pro-v3/internal/domain/setscheme"
	"github.com/waynenilsen/power-pro-v3/internal/repository"
)

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

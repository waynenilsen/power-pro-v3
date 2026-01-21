package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/progression"
	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
	"github.com/waynenilsen/power-pro-v3/internal/repository"
)

// ProgressionHandler handles HTTP requests for progression operations.
type ProgressionHandler struct {
	repo *repository.ProgressionRepository
}

// NewProgressionHandler creates a new ProgressionHandler.
func NewProgressionHandler(repo *repository.ProgressionRepository) *ProgressionHandler {
	return &ProgressionHandler{repo: repo}
}

// ProgressionResponse represents the API response format for a progression.
type ProgressionResponse struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Type       string          `json:"type"`
	Parameters json.RawMessage `json:"parameters"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
}

// LinearProgressionParams represents the parameters for a linear progression.
type LinearProgressionParams struct {
	Increment   float64 `json:"increment"`
	MaxType     string  `json:"maxType"`
	TriggerType string  `json:"triggerType"`
}

// CycleProgressionParams represents the parameters for a cycle progression.
type CycleProgressionParams struct {
	Increment float64 `json:"increment"`
	MaxType   string  `json:"maxType"`
}

// CreateProgressionRequest represents the request body for creating a progression.
type CreateProgressionRequest struct {
	Name       string          `json:"name"`
	Type       string          `json:"type"`
	Parameters json.RawMessage `json:"parameters"`
}

// UpdateProgressionRequest represents the request body for updating a progression.
type UpdateProgressionRequest struct {
	Name       *string          `json:"name,omitempty"`
	Type       *string          `json:"type,omitempty"`
	Parameters *json.RawMessage `json:"parameters,omitempty"`
}

func progressionEntityToResponse(entity *repository.ProgressionEntity) ProgressionResponse {
	return ProgressionResponse{
		ID:         entity.ID,
		Name:       entity.Name,
		Type:       string(entity.Type),
		Parameters: entity.Parameters,
		CreatedAt:  entity.CreatedAt,
		UpdatedAt:  entity.UpdatedAt,
	}
}

// List handles GET /progressions
func (h *ProgressionHandler) List(w http.ResponseWriter, r *http.Request) {
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

	// Filter by type
	var filterType *progression.ProgressionType
	if t := query.Get("type"); t != "" {
		pt := progression.ProgressionType(strings.ToUpper(t))
		if progression.ValidProgressionTypes[pt] {
			filterType = &pt
		}
	}

	params := repository.ProgressionListParams{
		Limit:      int64(pageSize),
		Offset:     int64((page - 1) * pageSize),
		FilterType: filterType,
	}

	entities, total, err := h.repo.List(params)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to list progressions", err))
		return
	}

	// Convert to response format
	data := make([]ProgressionResponse, len(entities))
	for i, entity := range entities {
		data[i] = progressionEntityToResponse(&entity)
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

// Get handles GET /progressions/{id}
func (h *ProgressionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing progression ID"))
		return
	}

	entity, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get progression", err))
		return
	}
	if entity == nil {
		writeDomainError(w, apperrors.NewNotFound("progression", id))
		return
	}

	writeJSON(w, http.StatusOK, progressionEntityToResponse(entity))
}

// Create handles POST /progressions
func (h *ProgressionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateProgressionRequest
	if err := readJSON(r, &req); err != nil {
		writeDomainError(w, apperrors.NewBadRequest("invalid request body"))
		return
	}

	// Validate request
	validationErrors := h.validateProgressionRequest(req.Name, req.Type, req.Parameters)
	if len(validationErrors) > 0 {
		writeDomainError(w, apperrors.NewValidationMsg("validation failed"), validationErrors...)
		return
	}

	// Generate UUID and timestamps
	id := uuid.New().String()
	now := time.Now()

	// Inject ID and Name into the parameters for domain validation
	enrichedParams, err := enrichProgressionParams(id, req.Name, req.Parameters)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to process parameters", err))
		return
	}

	entity := &repository.ProgressionEntity{
		ID:         id,
		Name:       req.Name,
		Type:       progression.ProgressionType(req.Type),
		Parameters: enrichedParams,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// Persist
	if err := h.repo.Create(entity); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to create progression", err))
		return
	}

	writeJSON(w, http.StatusCreated, progressionEntityToResponse(entity))
}

// enrichProgressionParams injects id and name into progression parameters.
// This is needed because the domain Progression types require id and name
// in the parameters JSON for validation and instantiation.
func enrichProgressionParams(id, name string, params json.RawMessage) (json.RawMessage, error) {
	var paramsMap map[string]interface{}
	if err := json.Unmarshal(params, &paramsMap); err != nil {
		return nil, err
	}

	paramsMap["id"] = id
	paramsMap["name"] = name

	return json.Marshal(paramsMap)
}

// Update handles PUT /progressions/{id}
func (h *ProgressionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing progression ID"))
		return
	}

	// Get existing progression
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get progression", err))
		return
	}
	if existing == nil {
		writeDomainError(w, apperrors.NewNotFound("progression", id))
		return
	}

	var req UpdateProgressionRequest
	if err := readJSON(r, &req); err != nil {
		writeDomainError(w, apperrors.NewBadRequest("invalid request body"))
		return
	}

	// Apply updates
	name := existing.Name
	if req.Name != nil {
		name = *req.Name
	}

	progType := string(existing.Type)
	if req.Type != nil {
		progType = *req.Type
	}

	params := existing.Parameters
	if req.Parameters != nil {
		params = *req.Parameters
	}

	// Validate the updated values
	validationErrors := h.validateProgressionRequest(name, progType, params)
	if len(validationErrors) > 0 {
		writeDomainError(w, apperrors.NewValidationMsg("validation failed"), validationErrors...)
		return
	}

	// Inject ID and Name into the parameters
	enrichedParams, err := enrichProgressionParams(id, name, params)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to process parameters", err))
		return
	}

	// Update entity
	existing.Name = name
	existing.Type = progression.ProgressionType(progType)
	existing.Parameters = enrichedParams
	existing.UpdatedAt = time.Now()

	// Persist
	if err := h.repo.Update(existing); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to update progression", err))
		return
	}

	writeJSON(w, http.StatusOK, progressionEntityToResponse(existing))
}

// Delete handles DELETE /progressions/{id}
func (h *ProgressionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing progression ID"))
		return
	}

	// Check progression exists
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get progression", err))
		return
	}
	if existing == nil {
		writeDomainError(w, apperrors.NewNotFound("progression", id))
		return
	}

	// Check for references (program_progressions)
	hasRefs, err := h.repo.HasProgramReferences(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to check references", err))
		return
	}
	if hasRefs {
		writeDomainError(w, apperrors.NewConflict("cannot delete progression: it is referenced by program progressions"))
		return
	}

	// Delete
	if err := h.repo.Delete(id); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to delete progression", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// validateProgressionRequest validates progression request data.
func (h *ProgressionHandler) validateProgressionRequest(name, progType string, params json.RawMessage) []string {
	var errors []string

	// Validate name
	if name == "" {
		errors = append(errors, "name is required")
	}

	// Validate type
	pt := progression.ProgressionType(progType)
	if err := progression.ValidateProgressionType(pt); err != nil {
		errors = append(errors, err.Error())
		return errors // Can't validate params without valid type
	}

	// Validate parameters based on type
	switch pt {
	case progression.TypeLinear:
		errors = append(errors, h.validateLinearParams(params)...)
	case progression.TypeCycle:
		errors = append(errors, h.validateCycleParams(params)...)
	}

	return errors
}

// validateLinearParams validates LinearProgression parameters.
func (h *ProgressionHandler) validateLinearParams(params json.RawMessage) []string {
	var errors []string
	var p LinearProgressionParams

	if err := json.Unmarshal(params, &p); err != nil {
		errors = append(errors, "invalid parameters: failed to parse as LinearProgression params")
		return errors
	}

	// Validate increment (required, positive)
	if p.Increment <= 0 {
		errors = append(errors, "increment must be positive")
	}

	// Validate maxType (required)
	maxType := progression.MaxType(p.MaxType)
	if err := progression.ValidateMaxType(maxType); err != nil {
		errors = append(errors, err.Error())
	}

	// Validate triggerType (required)
	triggerType := progression.TriggerType(p.TriggerType)
	if err := progression.ValidateTriggerType(triggerType); err != nil {
		errors = append(errors, err.Error())
	} else if triggerType != progression.TriggerAfterSession && triggerType != progression.TriggerAfterWeek {
		errors = append(errors, "linear progression only supports AFTER_SESSION and AFTER_WEEK triggers")
	}

	return errors
}

// validateCycleParams validates CycleProgression parameters.
func (h *ProgressionHandler) validateCycleParams(params json.RawMessage) []string {
	var errors []string
	var p CycleProgressionParams

	if err := json.Unmarshal(params, &p); err != nil {
		errors = append(errors, "invalid parameters: failed to parse as CycleProgression params")
		return errors
	}

	// Validate increment (required, positive)
	if p.Increment <= 0 {
		errors = append(errors, "increment must be positive")
	}

	// Validate maxType (required)
	maxType := progression.MaxType(p.MaxType)
	if err := progression.ValidateMaxType(maxType); err != nil {
		errors = append(errors, err.Error())
	}

	return errors
}

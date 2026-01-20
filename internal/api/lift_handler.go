package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/lift"
	"github.com/waynenilsen/power-pro-v3/internal/repository"
)

// LiftHandler handles HTTP requests for lift operations.
type LiftHandler struct {
	repo *repository.LiftRepository
}

// NewLiftHandler creates a new LiftHandler.
func NewLiftHandler(repo *repository.LiftRepository) *LiftHandler {
	return &LiftHandler{repo: repo}
}

// LiftResponse represents the API response format for a lift.
type LiftResponse struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Slug              string    `json:"slug"`
	IsCompetitionLift bool      `json:"isCompetitionLift"`
	ParentLiftID      *string   `json:"parentLiftId"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// CreateLiftRequest represents the request body for creating a lift.
type CreateLiftRequest struct {
	Name              string  `json:"name"`
	Slug              string  `json:"slug,omitempty"`
	IsCompetitionLift bool    `json:"isCompetitionLift"`
	ParentLiftID      *string `json:"parentLiftId,omitempty"`
}

// UpdateLiftRequest represents the request body for updating a lift.
type UpdateLiftRequest struct {
	Name              *string `json:"name,omitempty"`
	Slug              *string `json:"slug,omitempty"`
	IsCompetitionLift *bool   `json:"isCompetitionLift,omitempty"`
	ParentLiftID      *string `json:"parentLiftId,omitempty"`
	ClearParentLift   bool    `json:"clearParentLift,omitempty"`
}

func liftToResponse(l *lift.Lift) LiftResponse {
	return LiftResponse{
		ID:                l.ID,
		Name:              l.Name,
		Slug:              l.Slug,
		IsCompetitionLift: l.IsCompetitionLift,
		ParentLiftID:      l.ParentLiftID,
		CreatedAt:         l.CreatedAt,
		UpdatedAt:         l.UpdatedAt,
	}
}

// List handles GET /lifts
func (h *LiftHandler) List(w http.ResponseWriter, r *http.Request) {
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
	sortBy := repository.SortByName
	sortOrder := repository.SortAsc
	if s := query.Get("sortBy"); s != "" {
		switch strings.ToLower(s) {
		case "name":
			sortBy = repository.SortByName
		case "created_at", "createdat":
			sortBy = repository.SortByCreatedAt
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

	// Filter
	var filterCompetition *bool
	if f := query.Get("is_competition_lift"); f != "" {
		val := strings.ToLower(f) == "true" || f == "1"
		filterCompetition = &val
	}

	params := repository.ListParams{
		Limit:             int64(pageSize),
		Offset:            int64((page - 1) * pageSize),
		SortBy:            sortBy,
		SortOrder:         sortOrder,
		FilterCompetition: filterCompetition,
	}

	lifts, total, err := h.repo.List(params)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list lifts")
		return
	}

	// Convert to response format
	data := make([]LiftResponse, len(lifts))
	for i, l := range lifts {
		data[i] = liftToResponse(&l)
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

// Get handles GET /lifts/{id}
func (h *LiftHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing lift ID")
		return
	}

	l, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get lift")
		return
	}
	if l == nil {
		writeError(w, http.StatusNotFound, "Lift not found")
		return
	}

	writeJSON(w, http.StatusOK, liftToResponse(l))
}

// GetBySlug handles GET /lifts/by-slug/{slug}
func (h *LiftHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		writeError(w, http.StatusBadRequest, "Missing slug")
		return
	}

	l, err := h.repo.GetBySlug(slug)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get lift")
		return
	}
	if l == nil {
		writeError(w, http.StatusNotFound, "Lift not found")
		return
	}

	writeJSON(w, http.StatusOK, liftToResponse(l))
}

// Create handles POST /lifts
func (h *LiftHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateLiftRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Generate UUID
	id := uuid.New().String()

	// Use domain logic to create and validate
	input := lift.CreateLiftInput{
		Name:              req.Name,
		Slug:              req.Slug,
		IsCompetitionLift: req.IsCompetitionLift,
		ParentLiftID:      req.ParentLiftID,
	}

	newLift, result := lift.CreateLift(input, id, h.repo)
	if !result.Valid {
		details := make([]string, len(result.Errors))
		for i, err := range result.Errors {
			details[i] = err.Error()
		}
		writeError(w, http.StatusBadRequest, "Validation failed", details...)
		return
	}

	// Check for slug conflict
	exists, err := h.repo.SlugExists(newLift.Slug, nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to check slug uniqueness")
		return
	}
	if exists {
		writeError(w, http.StatusConflict, "Slug already exists")
		return
	}

	// Check parent lift exists if provided
	if req.ParentLiftID != nil {
		parent, err := h.repo.GetByID(*req.ParentLiftID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to verify parent lift")
			return
		}
		if parent == nil {
			writeError(w, http.StatusBadRequest, "Parent lift not found")
			return
		}
	}

	// Persist
	if err := h.repo.Create(newLift); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create lift")
		return
	}

	writeJSON(w, http.StatusCreated, liftToResponse(newLift))
}

// Update handles PUT /lifts/{id}
func (h *LiftHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing lift ID")
		return
	}

	// Get existing lift
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get lift")
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "Lift not found")
		return
	}

	var req UpdateLiftRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check for slug conflict if slug is being changed
	if req.Slug != nil && *req.Slug != existing.Slug {
		exists, err := h.repo.SlugExists(*req.Slug, &id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to check slug uniqueness")
			return
		}
		if exists {
			writeError(w, http.StatusConflict, "Slug already exists")
			return
		}
	}

	// Check parent lift exists if being set
	if req.ParentLiftID != nil && !req.ClearParentLift {
		parent, err := h.repo.GetByID(*req.ParentLiftID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to verify parent lift")
			return
		}
		if parent == nil {
			writeError(w, http.StatusBadRequest, "Parent lift not found")
			return
		}
	}

	// Use domain logic to update and validate
	input := lift.UpdateLiftInput{
		Name:              req.Name,
		Slug:              req.Slug,
		IsCompetitionLift: req.IsCompetitionLift,
		ParentLiftID:      req.ParentLiftID,
		ClearParentLift:   req.ClearParentLift,
	}

	result := lift.UpdateLift(existing, input, h.repo)
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
		writeError(w, http.StatusInternalServerError, "Failed to update lift")
		return
	}

	writeJSON(w, http.StatusOK, liftToResponse(existing))
}

// Delete handles DELETE /lifts/{id}
func (h *LiftHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing lift ID")
		return
	}

	// Check lift exists
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get lift")
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "Lift not found")
		return
	}

	// Check for references (child lifts)
	hasRefs, err := h.repo.HasChildReferences(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to check references")
		return
	}
	if hasRefs {
		writeError(w, http.StatusConflict, "Cannot delete lift: it is referenced by other lifts as a parent")
		return
	}

	// Note: In the future, also check for LiftMax references per NFR-005
	// For now, we only have the lifts table

	// Delete
	if err := h.repo.Delete(id); err != nil {
		if errors.Is(err, errors.New("FOREIGN KEY constraint failed")) {
			writeError(w, http.StatusConflict, "Cannot delete lift: it is referenced by other records")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to delete lift")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/liftmax"
	"github.com/waynenilsen/power-pro-v3/internal/repository"
)

// LiftMaxHandler handles HTTP requests for lift max operations.
type LiftMaxHandler struct {
	repo     *repository.LiftMaxRepository
	liftRepo *repository.LiftRepository
}

// NewLiftMaxHandler creates a new LiftMaxHandler.
func NewLiftMaxHandler(repo *repository.LiftMaxRepository, liftRepo *repository.LiftRepository) *LiftMaxHandler {
	return &LiftMaxHandler{repo: repo, liftRepo: liftRepo}
}

// LiftMaxResponse represents the API response format for a lift max.
type LiftMaxResponse struct {
	ID            string    `json:"id"`
	UserID        string    `json:"userId"`
	LiftID        string    `json:"liftId"`
	Type          string    `json:"type"`
	Value         float64   `json:"value"`
	EffectiveDate time.Time `json:"effectiveDate"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// CreateLiftMaxRequest represents the request body for creating a lift max.
type CreateLiftMaxRequest struct {
	LiftID        string     `json:"liftId"`
	Type          string     `json:"type"`
	Value         float64    `json:"value"`
	EffectiveDate *time.Time `json:"effectiveDate,omitempty"`
}

// UpdateLiftMaxRequest represents the request body for updating a lift max.
type UpdateLiftMaxRequest struct {
	Value         *float64   `json:"value,omitempty"`
	EffectiveDate *time.Time `json:"effectiveDate,omitempty"`
}

func liftMaxToResponse(m *liftmax.LiftMax) LiftMaxResponse {
	return LiftMaxResponse{
		ID:            m.ID,
		UserID:        m.UserID,
		LiftID:        m.LiftID,
		Type:          string(m.Type),
		Value:         m.Value,
		EffectiveDate: m.EffectiveDate,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

// List handles GET /users/{userId}/lift-maxes
func (h *LiftMaxHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "Missing user ID")
		return
	}

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

	// Sorting (default: descending by effective_date)
	sortOrder := repository.SortDesc
	if o := query.Get("sortOrder"); o != "" {
		switch strings.ToLower(o) {
		case "asc":
			sortOrder = repository.SortAsc
		case "desc":
			sortOrder = repository.SortDesc
		}
	}

	// Filters
	var filterLiftID *string
	var filterType *string

	if liftID := query.Get("lift_id"); liftID != "" {
		filterLiftID = &liftID
	}
	if maxType := query.Get("type"); maxType != "" {
		// Validate type value
		upperType := strings.ToUpper(maxType)
		if upperType != string(liftmax.OneRM) && upperType != string(liftmax.TrainingMax) {
			writeError(w, http.StatusBadRequest, "Invalid type filter: must be ONE_RM or TRAINING_MAX")
			return
		}
		filterType = &upperType
	}

	params := repository.LiftMaxListParams{
		UserID:    userID,
		LiftID:    filterLiftID,
		Type:      filterType,
		SortOrder: sortOrder,
		Limit:     int64(pageSize),
		Offset:    int64((page - 1) * pageSize),
	}

	maxes, total, err := h.repo.List(params)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list lift maxes")
		return
	}

	// Convert to response format
	data := make([]LiftMaxResponse, len(maxes))
	for i, m := range maxes {
		data[i] = liftMaxToResponse(&m)
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

// Get handles GET /lift-maxes/{id}
func (h *LiftMaxHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing lift max ID")
		return
	}

	m, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get lift max")
		return
	}
	if m == nil {
		writeError(w, http.StatusNotFound, "Lift max not found")
		return
	}

	writeJSON(w, http.StatusOK, liftMaxToResponse(m))
}

// Create handles POST /users/{userId}/lift-maxes
func (h *LiftMaxHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "Missing user ID")
		return
	}

	var req CreateLiftMaxRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Verify lift exists
	lift, err := h.liftRepo.GetByID(req.LiftID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to verify lift")
		return
	}
	if lift == nil {
		writeError(w, http.StatusBadRequest, "Lift not found")
		return
	}

	// Generate UUID
	id := uuid.New().String()

	// Use domain logic to create and validate
	input := liftmax.CreateLiftMaxInput{
		UserID:        userID,
		LiftID:        req.LiftID,
		Type:          liftmax.MaxType(req.Type),
		Value:         req.Value,
		EffectiveDate: req.EffectiveDate,
	}

	newMax, result := liftmax.CreateLiftMax(input, id, h.repo)
	if !result.Valid {
		details := make([]string, len(result.Errors))
		for i, err := range result.Errors {
			details[i] = err.Error()
		}
		writeError(w, http.StatusBadRequest, "Validation failed", details...)
		return
	}

	// Check unique constraint
	exists, err := h.repo.UniqueConstraintExists(newMax.UserID, newMax.LiftID, string(newMax.Type), newMax.EffectiveDate, nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to check uniqueness")
		return
	}
	if exists {
		writeError(w, http.StatusConflict, "A lift max with this user, lift, type, and effective date already exists")
		return
	}

	// Persist
	if err := h.repo.Create(newMax); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create lift max")
		return
	}

	// Return with warnings if any
	response := liftMaxToResponse(newMax)
	if result.HasWarnings() {
		writeJSON(w, http.StatusCreated, ResponseWithWarnings{
			Data:     response,
			Warnings: result.Warnings,
		})
		return
	}

	writeJSON(w, http.StatusCreated, response)
}

// Update handles PUT /lift-maxes/{id}
func (h *LiftMaxHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing lift max ID")
		return
	}

	// Get existing lift max
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get lift max")
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "Lift max not found")
		return
	}

	var req UpdateLiftMaxRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Build update input (type and liftId cannot be changed)
	input := liftmax.UpdateLiftMaxInput{
		Value:         req.Value,
		EffectiveDate: req.EffectiveDate,
	}

	// Use domain logic to update and validate
	result := liftmax.UpdateLiftMax(existing, input, h.repo)
	if !result.Valid {
		details := make([]string, len(result.Errors))
		for i, err := range result.Errors {
			details[i] = err.Error()
		}
		writeError(w, http.StatusBadRequest, "Validation failed", details...)
		return
	}

	// Check unique constraint if effective date changed
	if req.EffectiveDate != nil {
		exists, err := h.repo.UniqueConstraintExists(existing.UserID, existing.LiftID, string(existing.Type), existing.EffectiveDate, &id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to check uniqueness")
			return
		}
		if exists {
			writeError(w, http.StatusConflict, "A lift max with this user, lift, type, and effective date already exists")
			return
		}
	}

	// Persist
	if err := h.repo.Update(existing); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update lift max")
		return
	}

	// Return with warnings if any
	response := liftMaxToResponse(existing)
	if result.HasWarnings() {
		writeJSON(w, http.StatusOK, ResponseWithWarnings{
			Data:     response,
			Warnings: result.Warnings,
		})
		return
	}

	writeJSON(w, http.StatusOK, response)
}

// Delete handles DELETE /lift-maxes/{id}
func (h *LiftMaxHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing lift max ID")
		return
	}

	// Check lift max exists
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get lift max")
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "Lift max not found")
		return
	}

	// Delete
	if err := h.repo.Delete(id); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete lift max")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

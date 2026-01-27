package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/liftmax"
	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
	"github.com/waynenilsen/power-pro-v3/internal/middleware"
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
// Note: Type field is ignored - all lift maxes are created as ONE_RM.
// Training Max is automatically calculated at 90% of the 1RM.
type CreateLiftMaxRequest struct {
	LiftID        string     `json:"liftId"`
	Type          string     `json:"type"` // Ignored: always creates ONE_RM, TM is auto-calculated
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
		writeDomainError(w, apperrors.NewBadRequest("missing user ID"))
		return
	}

	// Parse query parameters
	query := r.URL.Query()

	// Pagination (limit/offset)
	pg := ParsePagination(query)

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

	// Filter by lift_id
	filterLiftID := ParseFilterString(query, "lift_id")

	// Filter by type (enum validation)
	filterType, err := ParseFilterEnum(query, "type", []string{string(liftmax.OneRM), string(liftmax.TrainingMax)})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	params := repository.LiftMaxListParams{
		UserID:    userID,
		LiftID:    filterLiftID,
		Type:      filterType,
		SortOrder: sortOrder,
		Limit:     int64(pg.Limit),
		Offset:    int64(pg.Offset),
	}

	maxes, total, err := h.repo.List(params)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to list lift maxes", err))
		return
	}

	// Convert to response format
	data := make([]LiftMaxResponse, len(maxes))
	for i, m := range maxes {
		data[i] = liftMaxToResponse(&m)
	}

	writePaginatedData(w, http.StatusOK, data, total, pg.Limit, pg.Offset)
}

// Get handles GET /lift-maxes/{id}
func (h *LiftMaxHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing lift max ID"))
		return
	}

	m, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get lift max", err))
		return
	}
	if m == nil {
		writeDomainError(w, apperrors.NewNotFound("lift max", id))
		return
	}

	// Check ownership: user must be owner or admin
	requestingUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)
	if requestingUserID != m.UserID && !isAdmin {
		writeDomainError(w, apperrors.NewForbidden("you do not have permission to access this resource"))
		return
	}

	writeData(w, http.StatusOK, liftMaxToResponse(m))
}

// Create handles POST /users/{userId}/lift-maxes
func (h *LiftMaxHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing user ID"))
		return
	}

	var req CreateLiftMaxRequest
	if err := readJSON(r, &req); err != nil {
		writeDomainError(w, apperrors.NewBadRequest("invalid request body"))
		return
	}

	// Verify lift exists
	lift, err := h.liftRepo.GetByID(req.LiftID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to verify lift", err))
		return
	}
	if lift == nil {
		writeDomainError(w, apperrors.NewValidation("liftId", "lift not found"))
		return
	}

	// Force type to ONE_RM - Training Max is auto-calculated
	reqType := liftmax.OneRM

	// Generate UUID
	id := uuid.New().String()

	// Use domain logic to create and validate
	input := liftmax.CreateLiftMaxInput{
		UserID:        userID,
		LiftID:        req.LiftID,
		Type:          reqType,
		Value:         req.Value,
		EffectiveDate: req.EffectiveDate,
	}

	newMax, result := liftmax.CreateLiftMax(input, id, h.repo)
	if !result.Valid {
		details := make([]string, len(result.Errors))
		for i, err := range result.Errors {
			details[i] = err.Error()
		}
		writeDomainError(w, apperrors.NewValidationMsg("validation failed"), details...)
		return
	}

	// Check unique constraint
	exists, err := h.repo.UniqueConstraintExists(newMax.UserID, newMax.LiftID, string(newMax.Type), newMax.EffectiveDate, nil)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to check uniqueness", err))
		return
	}
	if exists {
		writeDomainError(w, apperrors.NewConflict("a lift max with this user, lift, type, and effective date already exists"))
		return
	}

	// Persist the 1RM
	if err := h.repo.Create(newMax); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to create lift max", err))
		return
	}

	// Auto-create/update Training Max at 90% of 1RM
	if err := h.syncTrainingMax(newMax); err != nil {
		// Log but don't fail - the 1RM was created successfully
		result.AddWarning("Failed to auto-calculate Training Max: " + err.Error())
	}

	// Return with warnings if any
	response := liftMaxToResponse(newMax)
	if result.HasWarnings() {
		writeDataWithWarnings(w, http.StatusCreated, response, result.Warnings)
		return
	}

	writeData(w, http.StatusCreated, response)
}

// syncTrainingMax creates or updates a Training Max based on a 1RM value.
// The TM is set to 90% of the 1RM, rounded to the nearest 0.25.
func (h *LiftMaxHandler) syncTrainingMax(oneRM *liftmax.LiftMax) error {
	calculator := liftmax.NewMaxCalculator()
	tmValue, err := calculator.ConvertToTM(oneRM.Value, nil) // Uses default 90%
	if err != nil {
		return err
	}

	// Check if a TM already exists for this user/lift with the same effective date
	existingTM, err := h.repo.GetCurrentMax(oneRM.UserID, oneRM.LiftID, string(liftmax.TrainingMax))
	if err != nil {
		return err
	}

	now := time.Now()

	if existingTM != nil {
		// Update existing TM
		existingTM.Value = tmValue
		existingTM.EffectiveDate = oneRM.EffectiveDate
		existingTM.UpdatedAt = now
		return h.repo.Update(existingTM)
	}

	// Create new TM
	newTM := &liftmax.LiftMax{
		ID:            uuid.New().String(),
		UserID:        oneRM.UserID,
		LiftID:        oneRM.LiftID,
		Type:          liftmax.TrainingMax,
		Value:         tmValue,
		EffectiveDate: oneRM.EffectiveDate,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	return h.repo.Create(newTM)
}

// Update handles PUT /lift-maxes/{id}
func (h *LiftMaxHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing lift max ID"))
		return
	}

	// Get existing lift max
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get lift max", err))
		return
	}
	if existing == nil {
		writeDomainError(w, apperrors.NewNotFound("lift max", id))
		return
	}

	// Prevent direct modification of Training Maxes - they are auto-calculated from 1RM
	if existing.Type == liftmax.TrainingMax {
		writeDomainError(w, apperrors.NewBadRequest("Training Max cannot be modified directly - update your 1RM instead"))
		return
	}

	// Check ownership: user must be owner or admin
	requestingUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)
	if requestingUserID != existing.UserID && !isAdmin {
		writeDomainError(w, apperrors.NewForbidden("you do not have permission to modify this resource"))
		return
	}

	var req UpdateLiftMaxRequest
	if err := readJSON(r, &req); err != nil {
		writeDomainError(w, apperrors.NewBadRequest("invalid request body"))
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
		writeDomainError(w, apperrors.NewValidationMsg("validation failed"), details...)
		return
	}

	// Check unique constraint if effective date changed
	if req.EffectiveDate != nil {
		exists, err := h.repo.UniqueConstraintExists(existing.UserID, existing.LiftID, string(existing.Type), existing.EffectiveDate, &id)
		if err != nil {
			writeDomainError(w, apperrors.NewInternal("failed to check uniqueness", err))
			return
		}
		if exists {
			writeDomainError(w, apperrors.NewConflict("a lift max with this user, lift, type, and effective date already exists"))
			return
		}
	}

	// Persist
	if err := h.repo.Update(existing); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to update lift max", err))
		return
	}

	// Auto-sync Training Max when 1RM is updated
	if existing.Type == liftmax.OneRM {
		if err := h.syncTrainingMax(existing); err != nil {
			result.AddWarning("Failed to auto-calculate Training Max: " + err.Error())
		}
	}

	// Return with warnings if any
	response := liftMaxToResponse(existing)
	if result.HasWarnings() {
		writeDataWithWarnings(w, http.StatusOK, response, result.Warnings)
		return
	}

	writeData(w, http.StatusOK, response)
}

// Delete handles DELETE /lift-maxes/{id}
func (h *LiftMaxHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing lift max ID"))
		return
	}

	// Check lift max exists
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get lift max", err))
		return
	}
	if existing == nil {
		writeDomainError(w, apperrors.NewNotFound("lift max", id))
		return
	}

	// Check ownership: user must be owner or admin
	requestingUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)
	if requestingUserID != existing.UserID && !isAdmin {
		writeDomainError(w, apperrors.NewForbidden("you do not have permission to delete this resource"))
		return
	}

	// Delete
	if err := h.repo.Delete(id); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to delete lift max", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

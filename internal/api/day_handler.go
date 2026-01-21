package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/day"
	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
	"github.com/waynenilsen/power-pro-v3/internal/repository"
)

// DayHandler handles HTTP requests for day operations.
type DayHandler struct {
	repo             *repository.DayRepository
	prescriptionRepo *repository.PrescriptionRepository
}

// NewDayHandler creates a new DayHandler.
func NewDayHandler(repo *repository.DayRepository, prescriptionRepo *repository.PrescriptionRepository) *DayHandler {
	return &DayHandler{
		repo:             repo,
		prescriptionRepo: prescriptionRepo,
	}
}

// DayResponse represents the API response format for a day.
type DayResponse struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Slug      string                 `json:"slug"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	ProgramID *string                `json:"programId,omitempty"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
}

// DayWithPrescriptionsResponse represents the API response format for a day with its prescriptions.
type DayWithPrescriptionsResponse struct {
	ID            string                    `json:"id"`
	Name          string                    `json:"name"`
	Slug          string                    `json:"slug"`
	Metadata      map[string]interface{}    `json:"metadata,omitempty"`
	ProgramID     *string                   `json:"programId,omitempty"`
	Prescriptions []DayPrescriptionResponse `json:"prescriptions"`
	CreatedAt     time.Time                 `json:"createdAt"`
	UpdatedAt     time.Time                 `json:"updatedAt"`
}

// DayPrescriptionResponse represents a prescription within a day.
type DayPrescriptionResponse struct {
	ID             string    `json:"id"`
	PrescriptionID string    `json:"prescriptionId"`
	Order          int       `json:"order"`
	CreatedAt      time.Time `json:"createdAt"`
}

// CreateDayRequest represents the request body for creating a day.
type CreateDayRequest struct {
	Name      string                 `json:"name"`
	Slug      string                 `json:"slug,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	ProgramID *string                `json:"programId,omitempty"`
}

// UpdateDayRequest represents the request body for updating a day.
type UpdateDayRequest struct {
	Name           *string                `json:"name,omitempty"`
	Slug           *string                `json:"slug,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	ClearMetadata  bool                   `json:"clearMetadata,omitempty"`
	ProgramID      *string                `json:"programId,omitempty"`
	ClearProgramID bool                   `json:"clearProgramId,omitempty"`
}

func dayToResponse(d *day.Day) DayResponse {
	return DayResponse{
		ID:        d.ID,
		Name:      d.Name,
		Slug:      d.Slug,
		Metadata:  d.Metadata,
		ProgramID: d.ProgramID,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func dayToResponseWithPrescriptions(d *day.Day, prescriptions []day.DayPrescription) DayWithPrescriptionsResponse {
	prescriptionResponses := make([]DayPrescriptionResponse, len(prescriptions))
	for i, p := range prescriptions {
		prescriptionResponses[i] = DayPrescriptionResponse{
			ID:             p.ID,
			PrescriptionID: p.PrescriptionID,
			Order:          p.Order,
			CreatedAt:      p.CreatedAt,
		}
	}

	return DayWithPrescriptionsResponse{
		ID:            d.ID,
		Name:          d.Name,
		Slug:          d.Slug,
		Metadata:      d.Metadata,
		ProgramID:     d.ProgramID,
		Prescriptions: prescriptionResponses,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
}

// List handles GET /days
func (h *DayHandler) List(w http.ResponseWriter, r *http.Request) {
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
	sortBy := repository.DaySortByName
	sortOrder := repository.SortAsc
	if s := query.Get("sortBy"); s != "" {
		switch strings.ToLower(s) {
		case "name":
			sortBy = repository.DaySortByName
		case "created_at", "createdat":
			sortBy = repository.DaySortByCreatedAt
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

	// Filter by program_id
	var filterProgramID *string
	if programID := query.Get("program_id"); programID != "" {
		filterProgramID = &programID
	}

	params := repository.DayListParams{
		Limit:           int64(pageSize),
		Offset:          int64((page - 1) * pageSize),
		SortBy:          sortBy,
		SortOrder:       sortOrder,
		FilterProgramID: filterProgramID,
	}

	days, total, err := h.repo.List(params)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to list days", err))
		return
	}

	// Convert to response format
	data := make([]DayResponse, len(days))
	for i, d := range days {
		data[i] = dayToResponse(&d)
	}

	// Use standard envelope with offset-based pagination
	offset := (page - 1) * pageSize
	writePaginatedData(w, http.StatusOK, data, total, pageSize, offset)
}

// Get handles GET /days/{id}
func (h *DayHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing day ID"))
		return
	}

	d, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get day", err))
		return
	}
	if d == nil {
		writeDomainError(w, apperrors.NewNotFound("day", id))
		return
	}

	// Get prescriptions for this day
	prescriptions, err := h.repo.ListDayPrescriptions(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get day prescriptions", err))
		return
	}

	writeData(w, http.StatusOK, dayToResponseWithPrescriptions(d, prescriptions))
}

// GetBySlug handles GET /days/by-slug/{slug}
func (h *DayHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing slug"))
		return
	}

	// Optional program_id query parameter
	var programID *string
	if pid := r.URL.Query().Get("program_id"); pid != "" {
		programID = &pid
	}

	d, err := h.repo.GetBySlug(slug, programID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get day", err))
		return
	}
	if d == nil {
		writeDomainError(w, apperrors.NewNotFound("day", slug))
		return
	}

	// Get prescriptions for this day
	prescriptions, err := h.repo.ListDayPrescriptions(d.ID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get day prescriptions", err))
		return
	}

	writeData(w, http.StatusOK, dayToResponseWithPrescriptions(d, prescriptions))
}

// Create handles POST /days
func (h *DayHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateDayRequest
	if err := readJSON(r, &req); err != nil {
		writeDomainError(w, apperrors.NewBadRequest("invalid request body"))
		return
	}

	// Validate metadata JSON if provided
	if req.Metadata != nil {
		if _, err := json.Marshal(req.Metadata); err != nil {
			writeDomainError(w, apperrors.NewValidation("metadata", "invalid JSON"))
			return
		}
	}

	// Generate UUID
	id := uuid.New().String()

	// Use domain logic to create and validate
	input := day.CreateDayInput{
		Name:      req.Name,
		Slug:      req.Slug,
		Metadata:  req.Metadata,
		ProgramID: req.ProgramID,
	}

	newDay, result := day.CreateDay(input, id)
	if !result.Valid {
		details := make([]string, len(result.Errors))
		for i, err := range result.Errors {
			details[i] = err.Error()
		}
		writeDomainError(w, apperrors.NewValidationMsg("validation failed"), details...)
		return
	}

	// Check for slug conflict within the program
	exists, err := h.repo.SlugExists(newDay.Slug, newDay.ProgramID, nil)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to check slug uniqueness", err))
		return
	}
	if exists {
		writeDomainError(w, apperrors.NewConflict("slug already exists within this program"))
		return
	}

	// Persist
	if err := h.repo.Create(newDay); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to create day", err))
		return
	}

	writeData(w, http.StatusCreated, dayToResponse(newDay))
}

// Update handles PUT /days/{id}
func (h *DayHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing day ID"))
		return
	}

	// Get existing day
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get day", err))
		return
	}
	if existing == nil {
		writeDomainError(w, apperrors.NewNotFound("day", id))
		return
	}

	var req UpdateDayRequest
	if err := readJSON(r, &req); err != nil {
		writeDomainError(w, apperrors.NewBadRequest("invalid request body"))
		return
	}

	// Validate metadata JSON if provided
	if req.Metadata != nil {
		if _, err := json.Marshal(req.Metadata); err != nil {
			writeDomainError(w, apperrors.NewValidation("metadata", "invalid JSON"))
			return
		}
	}

	// Determine the new program ID for slug uniqueness check
	newProgramID := existing.ProgramID
	if req.ClearProgramID {
		newProgramID = nil
	} else if req.ProgramID != nil {
		newProgramID = req.ProgramID
	}

	// Check for slug conflict if slug or program is being changed
	newSlug := existing.Slug
	if req.Slug != nil {
		newSlug = *req.Slug
	}
	if newSlug != existing.Slug || !ptrStringEqual(newProgramID, existing.ProgramID) {
		exists, err := h.repo.SlugExists(newSlug, newProgramID, &id)
		if err != nil {
			writeDomainError(w, apperrors.NewInternal("failed to check slug uniqueness", err))
			return
		}
		if exists {
			writeDomainError(w, apperrors.NewConflict("slug already exists within this program"))
			return
		}
	}

	// Use domain logic to update and validate
	input := day.UpdateDayInput{
		Name:           req.Name,
		Slug:           req.Slug,
		Metadata:       req.Metadata,
		ClearMetadata:  req.ClearMetadata,
		ProgramID:      req.ProgramID,
		ClearProgramID: req.ClearProgramID,
	}

	result := day.UpdateDay(existing, input)
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
		writeDomainError(w, apperrors.NewInternal("failed to update day", err))
		return
	}

	writeData(w, http.StatusOK, dayToResponse(existing))
}

// Delete handles DELETE /days/{id}
func (h *DayHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing day ID"))
		return
	}

	// Check day exists
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get day", err))
		return
	}
	if existing == nil {
		writeDomainError(w, apperrors.NewNotFound("day", id))
		return
	}

	// Check if day is used in any weeks
	isUsed, err := h.repo.IsUsedInWeeks(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to check if day is used", err))
		return
	}
	if isUsed {
		writeDomainError(w, apperrors.NewConflict("cannot delete day: it is used in one or more weeks"))
		return
	}

	// Delete
	if err := h.repo.Delete(id); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to delete day", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper function to compare string pointers
func ptrStringEqual(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

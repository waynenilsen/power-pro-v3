package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/week"
	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
	"github.com/waynenilsen/power-pro-v3/internal/repository"
)

// WeekHandler handles HTTP requests for week operations.
type WeekHandler struct {
	repo *repository.WeekRepository
}

// NewWeekHandler creates a new WeekHandler.
func NewWeekHandler(repo *repository.WeekRepository) *WeekHandler {
	return &WeekHandler{
		repo: repo,
	}
}

// WeekResponse represents the API response format for a week.
type WeekResponse struct {
	ID         string    `json:"id"`
	WeekNumber int       `json:"weekNumber"`
	Variant    *string   `json:"variant,omitempty"`
	CycleID    string    `json:"cycleId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// WeekWithDaysResponse represents the API response format for a week with its day mappings.
type WeekWithDaysResponse struct {
	ID         string            `json:"id"`
	WeekNumber int               `json:"weekNumber"`
	Variant    *string           `json:"variant,omitempty"`
	CycleID    string            `json:"cycleId"`
	Days       []WeekDayResponse `json:"days"`
	CreatedAt  time.Time         `json:"createdAt"`
	UpdatedAt  time.Time         `json:"updatedAt"`
}

// WeekDayResponse represents a day within a week.
type WeekDayResponse struct {
	ID        string    `json:"id"`
	DayID     string    `json:"dayId"`
	DayOfWeek string    `json:"dayOfWeek"`
	CreatedAt time.Time `json:"createdAt"`
}

// CreateWeekRequest represents the request body for creating a week.
type CreateWeekRequest struct {
	WeekNumber int     `json:"weekNumber"`
	Variant    *string `json:"variant,omitempty"`
	CycleID    string  `json:"cycleId"`
}

// UpdateWeekRequest represents the request body for updating a week.
type UpdateWeekRequest struct {
	WeekNumber   *int    `json:"weekNumber,omitempty"`
	Variant      *string `json:"variant,omitempty"`
	ClearVariant bool    `json:"clearVariant,omitempty"`
	CycleID      *string `json:"cycleId,omitempty"`
}

func weekToResponse(w *week.Week) WeekResponse {
	return WeekResponse{
		ID:         w.ID,
		WeekNumber: w.WeekNumber,
		Variant:    w.Variant,
		CycleID:    w.CycleID,
		CreatedAt:  w.CreatedAt,
		UpdatedAt:  w.UpdatedAt,
	}
}

func weekToResponseWithDays(w *week.Week, days []week.WeekDay) WeekWithDaysResponse {
	dayResponses := make([]WeekDayResponse, len(days))
	for i, d := range days {
		dayResponses[i] = WeekDayResponse{
			ID:        d.ID,
			DayID:     d.DayID,
			DayOfWeek: string(d.DayOfWeek),
			CreatedAt: d.CreatedAt,
		}
	}

	return WeekWithDaysResponse{
		ID:         w.ID,
		WeekNumber: w.WeekNumber,
		Variant:    w.Variant,
		CycleID:    w.CycleID,
		Days:       dayResponses,
		CreatedAt:  w.CreatedAt,
		UpdatedAt:  w.UpdatedAt,
	}
}

// List handles GET /weeks
func (h *WeekHandler) List(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()

	// Pagination (limit/offset)
	pg := ParsePagination(query)

	// Sorting
	sortBy := repository.WeekSortByWeekNumber
	sortOrder := repository.SortAsc
	if s := query.Get("sortBy"); s != "" {
		switch strings.ToLower(s) {
		case "week_number", "weeknumber":
			sortBy = repository.WeekSortByWeekNumber
		case "created_at", "createdat":
			sortBy = repository.WeekSortByCreatedAt
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

	// Filter by cycle_id
	var filterCycleID *string
	if cycleID := query.Get("cycle_id"); cycleID != "" {
		filterCycleID = &cycleID
	}

	params := repository.WeekListParams{
		Limit:         int64(pg.Limit),
		Offset:        int64(pg.Offset),
		SortBy:        sortBy,
		SortOrder:     sortOrder,
		FilterCycleID: filterCycleID,
	}

	weeks, total, err := h.repo.List(params)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to list weeks", err))
		return
	}

	// Convert to response format
	data := make([]WeekResponse, len(weeks))
	for i, wk := range weeks {
		data[i] = weekToResponse(&wk)
	}

	writePaginatedData(w, http.StatusOK, data, total, pg.Limit, pg.Offset)
}

// Get handles GET /weeks/{id}
func (h *WeekHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing week ID"))
		return
	}

	wk, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get week", err))
		return
	}
	if wk == nil {
		writeDomainError(w, apperrors.NewNotFound("week", id))
		return
	}

	// Get days for this week
	days, err := h.repo.ListWeekDays(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get week days", err))
		return
	}

	writeData(w, http.StatusOK, weekToResponseWithDays(wk, days))
}

// Create handles POST /weeks
func (h *WeekHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateWeekRequest
	if err := readJSON(r, &req); err != nil {
		writeDomainError(w, apperrors.NewBadRequest("invalid request body"))
		return
	}

	// Check if cycle exists
	cycleExists, err := h.repo.CycleExists(req.CycleID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to verify cycle", err))
		return
	}
	if !cycleExists {
		writeDomainError(w, apperrors.NewValidation("cycleId", "cycle not found"))
		return
	}

	// Generate UUID
	id := uuid.New().String()

	// Use domain logic to create and validate
	input := week.CreateWeekInput{
		WeekNumber: req.WeekNumber,
		Variant:    req.Variant,
		CycleID:    req.CycleID,
	}

	newWeek, result := week.CreateWeek(input, id)
	if !result.Valid {
		details := make([]string, len(result.Errors))
		for i, err := range result.Errors {
			details[i] = err.Error()
		}
		writeDomainError(w, apperrors.NewValidationMsg("validation failed"), details...)
		return
	}

	// Check for week number conflict within the cycle
	exists, err := h.repo.WeekNumberExistsInCycle(newWeek.CycleID, newWeek.WeekNumber, nil)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to check week number uniqueness", err))
		return
	}
	if exists {
		writeDomainError(w, apperrors.NewConflict("week number already exists within this cycle"))
		return
	}

	// Persist
	if err := h.repo.Create(newWeek); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to create week", err))
		return
	}

	writeData(w, http.StatusCreated, weekToResponse(newWeek))
}

// Update handles PUT /weeks/{id}
func (h *WeekHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing week ID"))
		return
	}

	// Get existing week
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get week", err))
		return
	}
	if existing == nil {
		writeDomainError(w, apperrors.NewNotFound("week", id))
		return
	}

	var req UpdateWeekRequest
	if err := readJSON(r, &req); err != nil {
		writeDomainError(w, apperrors.NewBadRequest("invalid request body"))
		return
	}

	// If cycle is being changed, verify the new cycle exists
	newCycleID := existing.CycleID
	if req.CycleID != nil {
		cycleExists, err := h.repo.CycleExists(*req.CycleID)
		if err != nil {
			writeDomainError(w, apperrors.NewInternal("failed to verify cycle", err))
			return
		}
		if !cycleExists {
			writeDomainError(w, apperrors.NewValidation("cycleId", "cycle not found"))
			return
		}
		newCycleID = *req.CycleID
	}

	// Check for week number conflict if week number or cycle is being changed
	newWeekNumber := existing.WeekNumber
	if req.WeekNumber != nil {
		newWeekNumber = *req.WeekNumber
	}
	if newWeekNumber != existing.WeekNumber || newCycleID != existing.CycleID {
		exists, err := h.repo.WeekNumberExistsInCycle(newCycleID, newWeekNumber, &id)
		if err != nil {
			writeDomainError(w, apperrors.NewInternal("failed to check week number uniqueness", err))
			return
		}
		if exists {
			writeDomainError(w, apperrors.NewConflict("week number already exists within this cycle"))
			return
		}
	}

	// Use domain logic to update and validate
	input := week.UpdateWeekInput{
		WeekNumber:   req.WeekNumber,
		Variant:      req.Variant,
		ClearVariant: req.ClearVariant,
		CycleID:      req.CycleID,
	}

	result := week.UpdateWeek(existing, input)
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
		writeDomainError(w, apperrors.NewInternal("failed to update week", err))
		return
	}

	writeData(w, http.StatusOK, weekToResponse(existing))
}

// Delete handles DELETE /weeks/{id}
func (h *WeekHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing week ID"))
		return
	}

	// Check week exists
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get week", err))
		return
	}
	if existing == nil {
		writeDomainError(w, apperrors.NewNotFound("week", id))
		return
	}

	// Check if week is used in an active cycle
	isUsed, err := h.repo.IsUsedInActiveCycle(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to check if week is used", err))
		return
	}
	if isUsed {
		writeDomainError(w, apperrors.NewConflict("cannot delete week: it is part of an active cycle with enrolled users"))
		return
	}

	// Delete
	if err := h.repo.Delete(id); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to delete week", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

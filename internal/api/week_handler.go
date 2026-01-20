package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/week"
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

// AddDayRequest represents the request body for adding a day to a week.
type AddDayRequest struct {
	DayID     string `json:"dayId"`
	DayOfWeek string `json:"dayOfWeek"`
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
		Limit:         int64(pageSize),
		Offset:        int64((page - 1) * pageSize),
		SortBy:        sortBy,
		SortOrder:     sortOrder,
		FilterCycleID: filterCycleID,
	}

	weeks, total, err := h.repo.List(params)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list weeks")
		return
	}

	// Convert to response format
	data := make([]WeekResponse, len(weeks))
	for i, wk := range weeks {
		data[i] = weekToResponse(&wk)
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

// Get handles GET /weeks/{id}
func (h *WeekHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing week ID")
		return
	}

	wk, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get week")
		return
	}
	if wk == nil {
		writeError(w, http.StatusNotFound, "Week not found")
		return
	}

	// Get days for this week
	days, err := h.repo.ListWeekDays(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get week days")
		return
	}

	writeJSON(w, http.StatusOK, weekToResponseWithDays(wk, days))
}

// Create handles POST /weeks
func (h *WeekHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateWeekRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check if cycle exists
	cycleExists, err := h.repo.CycleExists(req.CycleID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to verify cycle")
		return
	}
	if !cycleExists {
		writeError(w, http.StatusBadRequest, "Cycle not found")
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
		writeError(w, http.StatusBadRequest, "Validation failed", details...)
		return
	}

	// Check for week number conflict within the cycle
	exists, err := h.repo.WeekNumberExistsInCycle(newWeek.CycleID, newWeek.WeekNumber, nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to check week number uniqueness")
		return
	}
	if exists {
		writeError(w, http.StatusConflict, "Week number already exists within this cycle")
		return
	}

	// Persist
	if err := h.repo.Create(newWeek); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create week")
		return
	}

	writeJSON(w, http.StatusCreated, weekToResponse(newWeek))
}

// Update handles PUT /weeks/{id}
func (h *WeekHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing week ID")
		return
	}

	// Get existing week
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get week")
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "Week not found")
		return
	}

	var req UpdateWeekRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// If cycle is being changed, verify the new cycle exists
	newCycleID := existing.CycleID
	if req.CycleID != nil {
		cycleExists, err := h.repo.CycleExists(*req.CycleID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to verify cycle")
			return
		}
		if !cycleExists {
			writeError(w, http.StatusBadRequest, "Cycle not found")
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
			writeError(w, http.StatusInternalServerError, "Failed to check week number uniqueness")
			return
		}
		if exists {
			writeError(w, http.StatusConflict, "Week number already exists within this cycle")
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
		writeError(w, http.StatusBadRequest, "Validation failed", details...)
		return
	}

	// Persist
	if err := h.repo.Update(existing); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update week")
		return
	}

	writeJSON(w, http.StatusOK, weekToResponse(existing))
}

// Delete handles DELETE /weeks/{id}
func (h *WeekHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing week ID")
		return
	}

	// Check week exists
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get week")
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "Week not found")
		return
	}

	// Check if week is used in an active cycle
	isUsed, err := h.repo.IsUsedInActiveCycle(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to check if week is used")
		return
	}
	if isUsed {
		writeError(w, http.StatusConflict, "Cannot delete week: it is part of an active cycle with enrolled users")
		return
	}

	// Delete
	if err := h.repo.Delete(id); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete week")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddDay handles POST /weeks/{id}/days
func (h *WeekHandler) AddDay(w http.ResponseWriter, r *http.Request) {
	weekID := r.PathValue("id")
	if weekID == "" {
		writeError(w, http.StatusBadRequest, "Missing week ID")
		return
	}

	// Check week exists
	wk, err := h.repo.GetByID(weekID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get week")
		return
	}
	if wk == nil {
		writeError(w, http.StatusNotFound, "Week not found")
		return
	}

	var req AddDayRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.DayID) == "" {
		writeError(w, http.StatusBadRequest, "dayId is required")
		return
	}

	if strings.TrimSpace(req.DayOfWeek) == "" {
		writeError(w, http.StatusBadRequest, "dayOfWeek is required")
		return
	}

	// Check day exists
	dayExists, err := h.repo.DayExists(req.DayID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to verify day")
		return
	}
	if !dayExists {
		writeError(w, http.StatusBadRequest, "Day not found")
		return
	}

	// Generate UUID
	id := uuid.New().String()

	// Create domain entity
	input := week.CreateWeekDayInput{
		WeekID:    weekID,
		DayID:     req.DayID,
		DayOfWeek: req.DayOfWeek,
	}

	newWeekDay, result := week.CreateWeekDay(input, id)
	if !result.Valid {
		details := make([]string, len(result.Errors))
		for i, err := range result.Errors {
			details[i] = err.Error()
		}
		writeError(w, http.StatusBadRequest, "Validation failed", details...)
		return
	}

	// Persist
	if err := h.repo.CreateWeekDay(newWeekDay); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to add day to week")
		return
	}

	resp := WeekDayResponse{
		ID:        newWeekDay.ID,
		DayID:     newWeekDay.DayID,
		DayOfWeek: string(newWeekDay.DayOfWeek),
		CreatedAt: newWeekDay.CreatedAt,
	}

	writeJSON(w, http.StatusCreated, resp)
}

// RemoveDay handles DELETE /weeks/{id}/days/{dayId}
func (h *WeekHandler) RemoveDay(w http.ResponseWriter, r *http.Request) {
	weekID := r.PathValue("id")
	dayID := r.PathValue("dayId")

	if weekID == "" {
		writeError(w, http.StatusBadRequest, "Missing week ID")
		return
	}
	if dayID == "" {
		writeError(w, http.StatusBadRequest, "Missing day ID")
		return
	}

	// day_of_week is required to identify which mapping to delete
	// (since the same day can be mapped to multiple days of the week)
	dayOfWeek := r.URL.Query().Get("day_of_week")
	if dayOfWeek == "" {
		writeError(w, http.StatusBadRequest, "day_of_week query parameter is required")
		return
	}

	// Check week exists
	wk, err := h.repo.GetByID(weekID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get week")
		return
	}
	if wk == nil {
		writeError(w, http.StatusNotFound, "Week not found")
		return
	}

	// Check if day mapping exists in this week
	existing, err := h.repo.GetWeekDayByWeekAndDayAndDayOfWeek(weekID, dayID, dayOfWeek)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to check day mapping")
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "Day mapping not found in this week")
		return
	}

	// Delete
	if err := h.repo.DeleteWeekDay(existing.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to remove day from week")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

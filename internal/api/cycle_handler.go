package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/cycle"
	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
	"github.com/waynenilsen/power-pro-v3/internal/repository"
)

// CycleHandler handles HTTP requests for cycle operations.
type CycleHandler struct {
	repo *repository.CycleRepository
}

// NewCycleHandler creates a new CycleHandler.
func NewCycleHandler(repo *repository.CycleRepository) *CycleHandler {
	return &CycleHandler{
		repo: repo,
	}
}

// CycleResponse represents the API response format for a cycle.
type CycleResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	LengthWeeks int       `json:"lengthWeeks"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CycleWeekResponse represents a week within a cycle response.
type CycleWeekResponse struct {
	ID         string `json:"id"`
	WeekNumber int    `json:"weekNumber"`
}

// CycleWithWeeksResponse represents the API response format for a cycle with its weeks.
type CycleWithWeeksResponse struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	LengthWeeks int                 `json:"lengthWeeks"`
	Weeks       []CycleWeekResponse `json:"weeks"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
}

// CreateCycleRequest represents the request body for creating a cycle.
type CreateCycleRequest struct {
	Name        string `json:"name"`
	LengthWeeks int    `json:"lengthWeeks"`
}

// UpdateCycleRequest represents the request body for updating a cycle.
type UpdateCycleRequest struct {
	Name        *string `json:"name,omitempty"`
	LengthWeeks *int    `json:"lengthWeeks,omitempty"`
}

func cycleToResponse(c *cycle.Cycle) CycleResponse {
	return CycleResponse{
		ID:          c.ID,
		Name:        c.Name,
		LengthWeeks: c.LengthWeeks,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func cycleToResponseWithWeeks(c *cycle.Cycle, weeks []cycle.CycleWeek) CycleWithWeeksResponse {
	weekResponses := make([]CycleWeekResponse, len(weeks))
	for i, w := range weeks {
		weekResponses[i] = CycleWeekResponse{
			ID:         w.ID,
			WeekNumber: w.WeekNumber,
		}
	}

	return CycleWithWeeksResponse{
		ID:          c.ID,
		Name:        c.Name,
		LengthWeeks: c.LengthWeeks,
		Weeks:       weekResponses,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

// List handles GET /cycles
func (h *CycleHandler) List(w http.ResponseWriter, r *http.Request) {
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
	sortBy := repository.CycleSortByName
	sortOrder := repository.SortAsc
	if s := query.Get("sortBy"); s != "" {
		switch strings.ToLower(s) {
		case "name":
			sortBy = repository.CycleSortByName
		case "created_at", "createdat":
			sortBy = repository.CycleSortByCreatedAt
		case "length_weeks", "lengthweeks":
			sortBy = repository.CycleSortByLengthWeeks
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

	params := repository.CycleListParams{
		Limit:     int64(pageSize),
		Offset:    int64((page - 1) * pageSize),
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	cycles, total, err := h.repo.List(params)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to list cycles", err))
		return
	}

	// Convert to response format
	data := make([]CycleResponse, len(cycles))
	for i, c := range cycles {
		data[i] = cycleToResponse(&c)
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

// Get handles GET /cycles/{id}
func (h *CycleHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing cycle ID"))
		return
	}

	c, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get cycle", err))
		return
	}
	if c == nil {
		writeDomainError(w, apperrors.NewNotFound("cycle", id))
		return
	}

	// Get weeks for this cycle
	weeks, err := h.repo.ListWeeks(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get cycle weeks", err))
		return
	}

	writeJSON(w, http.StatusOK, cycleToResponseWithWeeks(c, weeks))
}

// Create handles POST /cycles
func (h *CycleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateCycleRequest
	if err := readJSON(r, &req); err != nil {
		writeDomainError(w, apperrors.NewBadRequest("invalid request body"))
		return
	}

	// Generate UUID
	id := uuid.New().String()

	// Use domain logic to create and validate
	input := cycle.CreateCycleInput{
		Name:        req.Name,
		LengthWeeks: req.LengthWeeks,
	}

	newCycle, result := cycle.CreateCycle(input, id)
	if !result.Valid {
		details := make([]string, len(result.Errors))
		for i, err := range result.Errors {
			details[i] = err.Error()
		}
		writeDomainError(w, apperrors.NewValidationMsg("validation failed"), details...)
		return
	}

	// Persist
	if err := h.repo.Create(newCycle); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to create cycle", err))
		return
	}

	writeJSON(w, http.StatusCreated, cycleToResponse(newCycle))
}

// Update handles PUT /cycles/{id}
func (h *CycleHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing cycle ID"))
		return
	}

	// Get existing cycle
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get cycle", err))
		return
	}
	if existing == nil {
		writeDomainError(w, apperrors.NewNotFound("cycle", id))
		return
	}

	var req UpdateCycleRequest
	if err := readJSON(r, &req); err != nil {
		writeDomainError(w, apperrors.NewBadRequest("invalid request body"))
		return
	}

	// Use domain logic to update and validate
	input := cycle.UpdateCycleInput{
		Name:        req.Name,
		LengthWeeks: req.LengthWeeks,
	}

	result := cycle.UpdateCycle(existing, input)
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
		writeDomainError(w, apperrors.NewInternal("failed to update cycle", err))
		return
	}

	writeJSON(w, http.StatusOK, cycleToResponse(existing))
}

// Delete handles DELETE /cycles/{id}
func (h *CycleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing cycle ID"))
		return
	}

	// Check cycle exists
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get cycle", err))
		return
	}
	if existing == nil {
		writeDomainError(w, apperrors.NewNotFound("cycle", id))
		return
	}

	// Check if cycle is used by programs
	isUsed, err := h.repo.IsUsedByPrograms(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to check if cycle is used", err))
		return
	}
	if isUsed {
		writeDomainError(w, apperrors.NewConflict("cannot delete cycle: it is used by one or more programs"))
		return
	}

	// Delete
	if err := h.repo.Delete(id); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to delete cycle", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

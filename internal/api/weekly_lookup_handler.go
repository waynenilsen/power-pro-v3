package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
	"github.com/waynenilsen/power-pro-v3/internal/domain/weeklylookup"
	"github.com/waynenilsen/power-pro-v3/internal/repository"
)

// WeeklyLookupHandler handles HTTP requests for weekly lookup operations.
type WeeklyLookupHandler struct {
	repo *repository.WeeklyLookupRepository
}

// NewWeeklyLookupHandler creates a new WeeklyLookupHandler.
func NewWeeklyLookupHandler(repo *repository.WeeklyLookupRepository) *WeeklyLookupHandler {
	return &WeeklyLookupHandler{
		repo: repo,
	}
}

// WeeklyLookupEntryResponse represents an entry in the API response.
type WeeklyLookupEntryResponse struct {
	WeekNumber         int       `json:"weekNumber"`
	Percentages        []float64 `json:"percentages"`
	Reps               []int     `json:"reps"`
	PercentageModifier *float64  `json:"percentageModifier,omitempty"`
}

// WeeklyLookupResponse represents the API response format for a weekly lookup.
type WeeklyLookupResponse struct {
	ID        string                     `json:"id"`
	Name      string                     `json:"name"`
	Entries   []WeeklyLookupEntryResponse `json:"entries"`
	ProgramID *string                    `json:"programId,omitempty"`
	CreatedAt time.Time                  `json:"createdAt"`
	UpdatedAt time.Time                  `json:"updatedAt"`
}

// CreateWeeklyLookupRequest represents the request body for creating a weekly lookup.
type CreateWeeklyLookupRequest struct {
	Name      string                       `json:"name"`
	Entries   []WeeklyLookupEntryRequest   `json:"entries"`
	ProgramID *string                      `json:"programId,omitempty"`
}

// WeeklyLookupEntryRequest represents an entry in the create/update request.
type WeeklyLookupEntryRequest struct {
	WeekNumber         int       `json:"weekNumber"`
	Percentages        []float64 `json:"percentages"`
	Reps               []int     `json:"reps"`
	PercentageModifier *float64  `json:"percentageModifier,omitempty"`
}

// UpdateWeeklyLookupRequest represents the request body for updating a weekly lookup.
type UpdateWeeklyLookupRequest struct {
	Name      *string                      `json:"name,omitempty"`
	Entries   *[]WeeklyLookupEntryRequest  `json:"entries,omitempty"`
	ProgramID **string                     `json:"programId,omitempty"`
}

func weeklyLookupToResponse(w *weeklylookup.WeeklyLookup) WeeklyLookupResponse {
	entries := make([]WeeklyLookupEntryResponse, len(w.Entries))
	for i, e := range w.Entries {
		entries[i] = WeeklyLookupEntryResponse{
			WeekNumber:         e.WeekNumber,
			Percentages:        e.Percentages,
			Reps:               e.Reps,
			PercentageModifier: e.PercentageModifier,
		}
	}

	return WeeklyLookupResponse{
		ID:        w.ID,
		Name:      w.Name,
		Entries:   entries,
		ProgramID: w.ProgramID,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
	}
}

func requestEntriesToDomain(entries []WeeklyLookupEntryRequest) []weeklylookup.WeeklyLookupEntry {
	result := make([]weeklylookup.WeeklyLookupEntry, len(entries))
	for i, e := range entries {
		result[i] = weeklylookup.WeeklyLookupEntry{
			WeekNumber:         e.WeekNumber,
			Percentages:        e.Percentages,
			Reps:               e.Reps,
			PercentageModifier: e.PercentageModifier,
		}
	}
	return result
}

// List handles GET /weekly-lookups
func (h *WeeklyLookupHandler) List(w http.ResponseWriter, r *http.Request) {
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
	sortBy := repository.WeeklyLookupSortByName
	sortOrder := repository.SortAsc
	if s := query.Get("sortBy"); s != "" {
		switch strings.ToLower(s) {
		case "name":
			sortBy = repository.WeeklyLookupSortByName
		case "created_at", "createdat":
			sortBy = repository.WeeklyLookupSortByCreatedAt
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

	params := repository.WeeklyLookupListParams{
		Limit:     int64(pageSize),
		Offset:    int64((page - 1) * pageSize),
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	lookups, total, err := h.repo.List(params)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to list weekly lookups", err))
		return
	}

	// Convert to response format
	data := make([]WeeklyLookupResponse, len(lookups))
	for i, l := range lookups {
		data[i] = weeklyLookupToResponse(&l)
	}

	// Use standard envelope with offset-based pagination
	offset := (page - 1) * pageSize
	writePaginatedData(w, http.StatusOK, data, total, pageSize, offset)
}

// Get handles GET /weekly-lookups/{id}
func (h *WeeklyLookupHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing weekly lookup ID"))
		return
	}

	lookup, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get weekly lookup", err))
		return
	}
	if lookup == nil {
		writeDomainError(w, apperrors.NewNotFound("weekly lookup", id))
		return
	}

	writeData(w, http.StatusOK, weeklyLookupToResponse(lookup))
}

// Create handles POST /weekly-lookups
func (h *WeeklyLookupHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateWeeklyLookupRequest
	if err := readJSON(r, &req); err != nil {
		writeDomainError(w, apperrors.NewBadRequest("invalid request body"))
		return
	}

	// Generate UUID
	id := uuid.New().String()

	// Use domain logic to create and validate
	input := weeklylookup.CreateWeeklyLookupInput{
		Name:      req.Name,
		Entries:   requestEntriesToDomain(req.Entries),
		ProgramID: req.ProgramID,
	}

	newLookup, result := weeklylookup.CreateWeeklyLookup(input, id)
	if !result.Valid {
		details := make([]string, len(result.Errors))
		for i, err := range result.Errors {
			details[i] = err.Error()
		}
		writeDomainError(w, apperrors.NewValidationMsg("validation failed"), details...)
		return
	}

	// Persist
	if err := h.repo.Create(newLookup); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to create weekly lookup", err))
		return
	}

	writeData(w, http.StatusCreated, weeklyLookupToResponse(newLookup))
}

// Update handles PUT /weekly-lookups/{id}
func (h *WeeklyLookupHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing weekly lookup ID"))
		return
	}

	// Get existing lookup
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get weekly lookup", err))
		return
	}
	if existing == nil {
		writeDomainError(w, apperrors.NewNotFound("weekly lookup", id))
		return
	}

	var req UpdateWeeklyLookupRequest
	if err := readJSON(r, &req); err != nil {
		writeDomainError(w, apperrors.NewBadRequest("invalid request body"))
		return
	}

	// Use domain logic to update and validate
	input := weeklylookup.UpdateWeeklyLookupInput{
		Name:      req.Name,
		ProgramID: req.ProgramID,
	}
	if req.Entries != nil {
		entries := requestEntriesToDomain(*req.Entries)
		input.Entries = &entries
	}

	result := weeklylookup.UpdateWeeklyLookup(existing, input)
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
		writeDomainError(w, apperrors.NewInternal("failed to update weekly lookup", err))
		return
	}

	writeData(w, http.StatusOK, weeklyLookupToResponse(existing))
}

// Delete handles DELETE /weekly-lookups/{id}
func (h *WeeklyLookupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing weekly lookup ID"))
		return
	}

	// Check lookup exists
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get weekly lookup", err))
		return
	}
	if existing == nil {
		writeDomainError(w, apperrors.NewNotFound("weekly lookup", id))
		return
	}

	// Check if lookup is used by programs
	isUsed, err := h.repo.IsUsedByPrograms(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to check if weekly lookup is used", err))
		return
	}
	if isUsed {
		writeDomainError(w, apperrors.NewConflict("cannot delete weekly lookup: it is used by one or more programs"))
		return
	}

	// Delete
	if err := h.repo.Delete(id); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to delete weekly lookup", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

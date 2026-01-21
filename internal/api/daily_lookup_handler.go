package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/dailylookup"
	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
	"github.com/waynenilsen/power-pro-v3/internal/repository"
)

// DailyLookupHandler handles HTTP requests for daily lookup operations.
type DailyLookupHandler struct {
	repo *repository.DailyLookupRepository
}

// NewDailyLookupHandler creates a new DailyLookupHandler.
func NewDailyLookupHandler(repo *repository.DailyLookupRepository) *DailyLookupHandler {
	return &DailyLookupHandler{
		repo: repo,
	}
}

// DailyLookupEntryResponse represents an entry in the API response.
type DailyLookupEntryResponse struct {
	DayIdentifier      string  `json:"dayIdentifier"`
	PercentageModifier float64 `json:"percentageModifier"`
	IntensityLevel     *string `json:"intensityLevel,omitempty"`
}

// DailyLookupResponse represents the API response format for a daily lookup.
type DailyLookupResponse struct {
	ID        string                    `json:"id"`
	Name      string                    `json:"name"`
	Entries   []DailyLookupEntryResponse `json:"entries"`
	ProgramID *string                   `json:"programId,omitempty"`
	CreatedAt time.Time                 `json:"createdAt"`
	UpdatedAt time.Time                 `json:"updatedAt"`
}

// CreateDailyLookupRequest represents the request body for creating a daily lookup.
type CreateDailyLookupRequest struct {
	Name      string                     `json:"name"`
	Entries   []DailyLookupEntryRequest  `json:"entries"`
	ProgramID *string                    `json:"programId,omitempty"`
}

// DailyLookupEntryRequest represents an entry in the create/update request.
type DailyLookupEntryRequest struct {
	DayIdentifier      string  `json:"dayIdentifier"`
	PercentageModifier float64 `json:"percentageModifier"`
	IntensityLevel     *string `json:"intensityLevel,omitempty"`
}

// UpdateDailyLookupRequest represents the request body for updating a daily lookup.
type UpdateDailyLookupRequest struct {
	Name      *string                    `json:"name,omitempty"`
	Entries   *[]DailyLookupEntryRequest `json:"entries,omitempty"`
	ProgramID **string                   `json:"programId,omitempty"`
}

func dailyLookupToResponse(d *dailylookup.DailyLookup) DailyLookupResponse {
	entries := make([]DailyLookupEntryResponse, len(d.Entries))
	for i, e := range d.Entries {
		entries[i] = DailyLookupEntryResponse{
			DayIdentifier:      e.DayIdentifier,
			PercentageModifier: e.PercentageModifier,
			IntensityLevel:     e.IntensityLevel,
		}
	}

	return DailyLookupResponse{
		ID:        d.ID,
		Name:      d.Name,
		Entries:   entries,
		ProgramID: d.ProgramID,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func dailyRequestEntriesToDomain(entries []DailyLookupEntryRequest) []dailylookup.DailyLookupEntry {
	result := make([]dailylookup.DailyLookupEntry, len(entries))
	for i, e := range entries {
		result[i] = dailylookup.DailyLookupEntry{
			DayIdentifier:      e.DayIdentifier,
			PercentageModifier: e.PercentageModifier,
			IntensityLevel:     e.IntensityLevel,
		}
	}
	return result
}

// List handles GET /daily-lookups
func (h *DailyLookupHandler) List(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()

	// Pagination (limit/offset)
	pg := ParsePagination(query)

	// Sorting
	sortBy := repository.DailyLookupSortByName
	sortOrder := repository.SortAsc
	if s := query.Get("sortBy"); s != "" {
		switch strings.ToLower(s) {
		case "name":
			sortBy = repository.DailyLookupSortByName
		case "created_at", "createdat":
			sortBy = repository.DailyLookupSortByCreatedAt
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

	params := repository.DailyLookupListParams{
		Limit:     int64(pg.Limit),
		Offset:    int64(pg.Offset),
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	lookups, total, err := h.repo.List(params)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to list daily lookups", err))
		return
	}

	// Convert to response format
	data := make([]DailyLookupResponse, len(lookups))
	for i, l := range lookups {
		data[i] = dailyLookupToResponse(&l)
	}

	writePaginatedData(w, http.StatusOK, data, total, pg.Limit, pg.Offset)
}

// Get handles GET /daily-lookups/{id}
func (h *DailyLookupHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing daily lookup ID"))
		return
	}

	lookup, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get daily lookup", err))
		return
	}
	if lookup == nil {
		writeDomainError(w, apperrors.NewNotFound("daily lookup", id))
		return
	}

	writeData(w, http.StatusOK, dailyLookupToResponse(lookup))
}

// Create handles POST /daily-lookups
func (h *DailyLookupHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateDailyLookupRequest
	if err := readJSON(r, &req); err != nil {
		writeDomainError(w, apperrors.NewBadRequest("invalid request body"))
		return
	}

	// Generate UUID
	id := uuid.New().String()

	// Use domain logic to create and validate
	input := dailylookup.CreateDailyLookupInput{
		Name:      req.Name,
		Entries:   dailyRequestEntriesToDomain(req.Entries),
		ProgramID: req.ProgramID,
	}

	newLookup, result := dailylookup.CreateDailyLookup(input, id)
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
		writeDomainError(w, apperrors.NewInternal("failed to create daily lookup", err))
		return
	}

	writeData(w, http.StatusCreated, dailyLookupToResponse(newLookup))
}

// Update handles PUT /daily-lookups/{id}
func (h *DailyLookupHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing daily lookup ID"))
		return
	}

	// Get existing lookup
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get daily lookup", err))
		return
	}
	if existing == nil {
		writeDomainError(w, apperrors.NewNotFound("daily lookup", id))
		return
	}

	var req UpdateDailyLookupRequest
	if err := readJSON(r, &req); err != nil {
		writeDomainError(w, apperrors.NewBadRequest("invalid request body"))
		return
	}

	// Use domain logic to update and validate
	input := dailylookup.UpdateDailyLookupInput{
		Name:      req.Name,
		ProgramID: req.ProgramID,
	}
	if req.Entries != nil {
		entries := dailyRequestEntriesToDomain(*req.Entries)
		input.Entries = &entries
	}

	result := dailylookup.UpdateDailyLookup(existing, input)
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
		writeDomainError(w, apperrors.NewInternal("failed to update daily lookup", err))
		return
	}

	writeData(w, http.StatusOK, dailyLookupToResponse(existing))
}

// Delete handles DELETE /daily-lookups/{id}
func (h *DailyLookupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing daily lookup ID"))
		return
	}

	// Check lookup exists
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get daily lookup", err))
		return
	}
	if existing == nil {
		writeDomainError(w, apperrors.NewNotFound("daily lookup", id))
		return
	}

	// Check if lookup is used by programs
	isUsed, err := h.repo.IsUsedByPrograms(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to check if daily lookup is used", err))
		return
	}
	if isUsed {
		writeDomainError(w, apperrors.NewConflict("cannot delete daily lookup: it is used by one or more programs"))
		return
	}

	// Delete
	if err := h.repo.Delete(id); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to delete daily lookup", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

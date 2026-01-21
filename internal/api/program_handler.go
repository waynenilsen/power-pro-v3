package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/program"
	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
	"github.com/waynenilsen/power-pro-v3/internal/repository"
)

// ProgramHandler handles HTTP requests for program operations.
type ProgramHandler struct {
	repo      *repository.ProgramRepository
	cycleRepo *repository.CycleRepository
}

// NewProgramHandler creates a new ProgramHandler.
func NewProgramHandler(repo *repository.ProgramRepository, cycleRepo *repository.CycleRepository) *ProgramHandler {
	return &ProgramHandler{
		repo:      repo,
		cycleRepo: cycleRepo,
	}
}

// ProgramResponse represents the API response format for a program (list view).
type ProgramResponse struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Slug            string    `json:"slug"`
	Description     *string   `json:"description,omitempty"`
	CycleID         string    `json:"cycleId"`
	WeeklyLookupID  *string   `json:"weeklyLookupId,omitempty"`
	DailyLookupID   *string   `json:"dailyLookupId,omitempty"`
	DefaultRounding *float64  `json:"defaultRounding,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// ProgramCycleWeekResponse represents a week within a cycle.
type ProgramCycleWeekResponse struct {
	ID         string `json:"id"`
	WeekNumber int    `json:"weekNumber"`
}

// ProgramCycleResponse represents embedded cycle info in a program response.
type ProgramCycleResponse struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
	LengthWeeks int                        `json:"lengthWeeks"`
	Weeks       []ProgramCycleWeekResponse `json:"weeks"`
}

// LookupReferenceResponse represents a lookup table reference.
type LookupReferenceResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ProgramDetailResponse represents the API response format for a program (detail view with embedded cycle).
type ProgramDetailResponse struct {
	ID              string                   `json:"id"`
	Name            string                   `json:"name"`
	Slug            string                   `json:"slug"`
	Description     *string                  `json:"description,omitempty"`
	Cycle           *ProgramCycleResponse    `json:"cycle"`
	WeeklyLookup    *LookupReferenceResponse `json:"weeklyLookup,omitempty"`
	DailyLookup     *LookupReferenceResponse `json:"dailyLookup,omitempty"`
	DefaultRounding *float64                 `json:"defaultRounding,omitempty"`
	CreatedAt       time.Time                `json:"createdAt"`
	UpdatedAt       time.Time                `json:"updatedAt"`
}

// CreateProgramRequest represents the request body for creating a program.
type CreateProgramRequest struct {
	Name            string   `json:"name"`
	Slug            string   `json:"slug"`
	Description     *string  `json:"description,omitempty"`
	CycleID         string   `json:"cycleId"`
	WeeklyLookupID  *string  `json:"weeklyLookupId,omitempty"`
	DailyLookupID   *string  `json:"dailyLookupId,omitempty"`
	DefaultRounding *float64 `json:"defaultRounding,omitempty"`
}

// UpdateProgramRequest represents the request body for updating a program.
type UpdateProgramRequest struct {
	Name            *string   `json:"name,omitempty"`
	Slug            *string   `json:"slug,omitempty"`
	Description     **string  `json:"description,omitempty"`
	CycleID         *string   `json:"cycleId,omitempty"`
	WeeklyLookupID  **string  `json:"weeklyLookupId,omitempty"`
	DailyLookupID   **string  `json:"dailyLookupId,omitempty"`
	DefaultRounding **float64 `json:"defaultRounding,omitempty"`
}

func programToResponse(p *program.Program) ProgramResponse {
	return ProgramResponse{
		ID:              p.ID,
		Name:            p.Name,
		Slug:            p.Slug,
		Description:     p.Description,
		CycleID:         p.CycleID,
		WeeklyLookupID:  p.WeeklyLookupID,
		DailyLookupID:   p.DailyLookupID,
		DefaultRounding: p.DefaultRounding,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

func programToDetailResponse(p *program.Program, cycle *program.ProgramCycle, weeklyLookup *program.LookupReference, dailyLookup *program.LookupReference) ProgramDetailResponse {
	var cycleResp *ProgramCycleResponse
	if cycle != nil {
		weeks := make([]ProgramCycleWeekResponse, len(cycle.Weeks))
		for i, w := range cycle.Weeks {
			weeks[i] = ProgramCycleWeekResponse{
				ID:         w.ID,
				WeekNumber: w.WeekNumber,
			}
		}
		cycleResp = &ProgramCycleResponse{
			ID:          cycle.ID,
			Name:        cycle.Name,
			LengthWeeks: cycle.LengthWeeks,
			Weeks:       weeks,
		}
	}

	var weeklyLookupResp *LookupReferenceResponse
	if weeklyLookup != nil {
		weeklyLookupResp = &LookupReferenceResponse{
			ID:   weeklyLookup.ID,
			Name: weeklyLookup.Name,
		}
	}

	var dailyLookupResp *LookupReferenceResponse
	if dailyLookup != nil {
		dailyLookupResp = &LookupReferenceResponse{
			ID:   dailyLookup.ID,
			Name: dailyLookup.Name,
		}
	}

	return ProgramDetailResponse{
		ID:              p.ID,
		Name:            p.Name,
		Slug:            p.Slug,
		Description:     p.Description,
		Cycle:           cycleResp,
		WeeklyLookup:    weeklyLookupResp,
		DailyLookup:     dailyLookupResp,
		DefaultRounding: p.DefaultRounding,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

// List handles GET /programs
func (h *ProgramHandler) List(w http.ResponseWriter, r *http.Request) {
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
	sortBy := repository.ProgramSortByName
	sortOrder := repository.SortAsc
	if s := query.Get("sortBy"); s != "" {
		switch strings.ToLower(s) {
		case "name":
			sortBy = repository.ProgramSortByName
		case "created_at", "createdat":
			sortBy = repository.ProgramSortByCreatedAt
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

	params := repository.ProgramListParams{
		Limit:     int64(pageSize),
		Offset:    int64((page - 1) * pageSize),
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	programs, total, err := h.repo.List(params)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to list programs", err))
		return
	}

	// Convert to response format
	data := make([]ProgramResponse, len(programs))
	for i, p := range programs {
		data[i] = programToResponse(&p)
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

// Get handles GET /programs/{id}
func (h *ProgramHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing program ID"))
		return
	}

	p, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get program", err))
		return
	}
	if p == nil {
		writeDomainError(w, apperrors.NewNotFound("program", id))
		return
	}

	// Get the associated cycle with its weeks
	cycle, err := h.repo.GetCycleForProgram(p.CycleID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get program cycle", err))
		return
	}

	// Get lookup references if present
	var weeklyLookup *program.LookupReference
	if p.WeeklyLookupID != nil {
		weeklyLookup, err = h.repo.GetWeeklyLookupReference(*p.WeeklyLookupID)
		if err != nil {
			writeDomainError(w, apperrors.NewInternal("failed to get weekly lookup", err))
			return
		}
	}

	var dailyLookup *program.LookupReference
	if p.DailyLookupID != nil {
		dailyLookup, err = h.repo.GetDailyLookupReference(*p.DailyLookupID)
		if err != nil {
			writeDomainError(w, apperrors.NewInternal("failed to get daily lookup", err))
			return
		}
	}

	writeJSON(w, http.StatusOK, programToDetailResponse(p, cycle, weeklyLookup, dailyLookup))
}

// Create handles POST /programs
func (h *ProgramHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateProgramRequest
	if err := readJSON(r, &req); err != nil {
		writeDomainError(w, apperrors.NewBadRequest("invalid request body"))
		return
	}

	// Check if cycle exists
	cycle, err := h.cycleRepo.GetByID(req.CycleID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to verify cycle", err))
		return
	}
	if cycle == nil {
		writeDomainError(w, apperrors.NewValidation("cycleId", "cycle not found"))
		return
	}

	// Check slug uniqueness
	slugExists, err := h.repo.SlugExists(req.Slug)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to check slug uniqueness", err))
		return
	}
	if slugExists {
		writeDomainError(w, apperrors.NewConflict("a program with this slug already exists"))
		return
	}

	// Generate UUID
	id := uuid.New().String()

	// Use domain logic to create and validate
	input := program.CreateProgramInput{
		Name:            req.Name,
		Slug:            req.Slug,
		Description:     req.Description,
		CycleID:         req.CycleID,
		WeeklyLookupID:  req.WeeklyLookupID,
		DailyLookupID:   req.DailyLookupID,
		DefaultRounding: req.DefaultRounding,
	}

	newProgram, result := program.CreateProgram(input, id)
	if !result.Valid {
		details := make([]string, len(result.Errors))
		for i, err := range result.Errors {
			details[i] = err.Error()
		}
		writeDomainError(w, apperrors.NewValidationMsg("validation failed"), details...)
		return
	}

	// Persist
	if err := h.repo.Create(newProgram); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to create program", err))
		return
	}

	writeJSON(w, http.StatusCreated, programToResponse(newProgram))
}

// Update handles PUT /programs/{id}
func (h *ProgramHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing program ID"))
		return
	}

	// Get existing program
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get program", err))
		return
	}
	if existing == nil {
		writeDomainError(w, apperrors.NewNotFound("program", id))
		return
	}

	var req UpdateProgramRequest
	if err := readJSON(r, &req); err != nil {
		writeDomainError(w, apperrors.NewBadRequest("invalid request body"))
		return
	}

	// Check cycle_id if provided
	if req.CycleID != nil {
		cycle, err := h.cycleRepo.GetByID(*req.CycleID)
		if err != nil {
			writeDomainError(w, apperrors.NewInternal("failed to verify cycle", err))
			return
		}
		if cycle == nil {
			writeDomainError(w, apperrors.NewValidation("cycleId", "cycle not found"))
			return
		}
	}

	// Check slug uniqueness if changing
	if req.Slug != nil && *req.Slug != existing.Slug {
		slugExists, err := h.repo.SlugExistsExcluding(*req.Slug, id)
		if err != nil {
			writeDomainError(w, apperrors.NewInternal("failed to check slug uniqueness", err))
			return
		}
		if slugExists {
			writeDomainError(w, apperrors.NewConflict("a program with this slug already exists"))
			return
		}
	}

	// Use domain logic to update and validate
	input := program.UpdateProgramInput{
		Name:            req.Name,
		Slug:            req.Slug,
		Description:     req.Description,
		CycleID:         req.CycleID,
		WeeklyLookupID:  req.WeeklyLookupID,
		DailyLookupID:   req.DailyLookupID,
		DefaultRounding: req.DefaultRounding,
	}

	result := program.UpdateProgram(existing, input)
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
		writeDomainError(w, apperrors.NewInternal("failed to update program", err))
		return
	}

	writeJSON(w, http.StatusOK, programToResponse(existing))
}

// Delete handles DELETE /programs/{id}
func (h *ProgramHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing program ID"))
		return
	}

	// Check program exists
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get program", err))
		return
	}
	if existing == nil {
		writeDomainError(w, apperrors.NewNotFound("program", id))
		return
	}

	// Check if any users are enrolled
	hasEnrolled, err := h.repo.HasEnrolledUsers(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to check if users are enrolled", err))
		return
	}
	if hasEnrolled {
		writeDomainError(w, apperrors.NewConflict("cannot delete program: users are enrolled"))
		return
	}

	// Delete
	if err := h.repo.Delete(id); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to delete program", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

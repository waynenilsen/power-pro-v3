package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/repository"
)

// ProgramProgressionHandler handles HTTP requests for program progression configuration.
type ProgramProgressionHandler struct {
	ppRepo         *repository.ProgramProgressionRepository
	programRepo    *repository.ProgramRepository
	progressionRepo *repository.ProgressionRepository
	liftRepo       *repository.LiftRepository
}

// NewProgramProgressionHandler creates a new ProgramProgressionHandler.
func NewProgramProgressionHandler(
	ppRepo *repository.ProgramProgressionRepository,
	programRepo *repository.ProgramRepository,
	progressionRepo *repository.ProgressionRepository,
	liftRepo *repository.LiftRepository,
) *ProgramProgressionHandler {
	return &ProgramProgressionHandler{
		ppRepo:         ppRepo,
		programRepo:    programRepo,
		progressionRepo: progressionRepo,
		liftRepo:       liftRepo,
	}
}

// ProgramProgressionResponse represents the API response format for a program progression configuration.
type ProgramProgressionResponse struct {
	ID                string          `json:"id"`
	ProgramID         string          `json:"programId"`
	ProgressionID     string          `json:"progressionId"`
	LiftID            *string         `json:"liftId"`
	Priority          int64           `json:"priority"`
	Enabled           bool            `json:"enabled"`
	OverrideIncrement *float64        `json:"overrideIncrement,omitempty"`
	CreatedAt         time.Time       `json:"createdAt"`
	UpdatedAt         time.Time       `json:"updatedAt"`
	Progression       *ProgressionRef `json:"progression,omitempty"`
}

// ProgressionRef is a reference to progression details included in the response.
type ProgressionRef struct {
	Name       string          `json:"name"`
	Type       string          `json:"type"`
	Parameters json.RawMessage `json:"parameters"`
}

// CreateProgramProgressionRequest represents the request body for creating a program progression.
type CreateProgramProgressionRequest struct {
	ProgressionID     string   `json:"progressionId"`
	LiftID            *string  `json:"liftId,omitempty"`
	Priority          *int64   `json:"priority,omitempty"`
	Enabled           *bool    `json:"enabled,omitempty"`
	OverrideIncrement *float64 `json:"overrideIncrement,omitempty"`
}

// UpdateProgramProgressionRequest represents the request body for updating a program progression.
type UpdateProgramProgressionRequest struct {
	Priority          *int64   `json:"priority,omitempty"`
	Enabled           *bool    `json:"enabled,omitempty"`
	OverrideIncrement *float64 `json:"overrideIncrement,omitempty"`
}

func programProgressionEntityToResponse(entity *repository.ProgramProgressionEntity) ProgramProgressionResponse {
	return ProgramProgressionResponse{
		ID:                entity.ID,
		ProgramID:         entity.ProgramID,
		ProgressionID:     entity.ProgressionID,
		LiftID:            entity.LiftID,
		Priority:          entity.Priority,
		Enabled:           entity.Enabled,
		OverrideIncrement: entity.OverrideIncrement,
		CreatedAt:         entity.CreatedAt,
		UpdatedAt:         entity.UpdatedAt,
	}
}

func programProgressionWithDetailsToResponse(entity *repository.ProgramProgressionWithDetails) ProgramProgressionResponse {
	return ProgramProgressionResponse{
		ID:                entity.ID,
		ProgramID:         entity.ProgramID,
		ProgressionID:     entity.ProgressionID,
		LiftID:            entity.LiftID,
		Priority:          entity.Priority,
		Enabled:           entity.Enabled,
		OverrideIncrement: entity.OverrideIncrement,
		CreatedAt:         entity.CreatedAt,
		UpdatedAt:         entity.UpdatedAt,
		Progression: &ProgressionRef{
			Name:       entity.ProgressionName,
			Type:       string(entity.ProgressionType),
			Parameters: entity.ProgressionParameters,
		},
	}
}

// List handles GET /programs/{programId}/progressions
func (h *ProgramProgressionHandler) List(w http.ResponseWriter, r *http.Request) {
	programID := r.PathValue("programId")
	if programID == "" {
		writeError(w, http.StatusBadRequest, "Missing program ID")
		return
	}

	// Verify program exists
	program, err := h.programRepo.GetByID(programID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get program")
		return
	}
	if program == nil {
		writeError(w, http.StatusNotFound, "Program not found")
		return
	}

	// Get program progressions with details
	entities, err := h.ppRepo.ListByProgram(programID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list program progressions")
		return
	}

	// Convert to response format
	data := make([]ProgramProgressionResponse, len(entities))
	for i, entity := range entities {
		data[i] = programProgressionWithDetailsToResponse(&entity)
	}

	writeJSON(w, http.StatusOK, data)
}

// Get handles GET /programs/{programId}/progressions/{configId}
func (h *ProgramProgressionHandler) Get(w http.ResponseWriter, r *http.Request) {
	programID := r.PathValue("programId")
	configID := r.PathValue("configId")

	if programID == "" {
		writeError(w, http.StatusBadRequest, "Missing program ID")
		return
	}
	if configID == "" {
		writeError(w, http.StatusBadRequest, "Missing config ID")
		return
	}

	// Get the program progression
	entity, err := h.ppRepo.GetByID(configID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get program progression")
		return
	}
	if entity == nil {
		writeError(w, http.StatusNotFound, "Program progression configuration not found")
		return
	}

	// Verify it belongs to the specified program
	if entity.ProgramID != programID {
		writeError(w, http.StatusNotFound, "Program progression configuration not found")
		return
	}

	writeJSON(w, http.StatusOK, programProgressionEntityToResponse(entity))
}

// Create handles POST /programs/{programId}/progressions
func (h *ProgramProgressionHandler) Create(w http.ResponseWriter, r *http.Request) {
	programID := r.PathValue("programId")
	if programID == "" {
		writeError(w, http.StatusBadRequest, "Missing program ID")
		return
	}

	// Verify program exists
	program, err := h.programRepo.GetByID(programID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get program")
		return
	}
	if program == nil {
		writeError(w, http.StatusNotFound, "Program not found")
		return
	}

	var req CreateProgramProgressionRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	validationErrors := h.validateCreateRequest(&req)
	if len(validationErrors) > 0 {
		writeError(w, http.StatusBadRequest, "Validation failed", validationErrors...)
		return
	}

	// Verify progression exists
	progression, err := h.progressionRepo.GetByID(req.ProgressionID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get progression")
		return
	}
	if progression == nil {
		writeError(w, http.StatusBadRequest, "Progression not found", "progressionId references a non-existent progression")
		return
	}

	// Verify lift exists if provided
	if req.LiftID != nil {
		lift, err := h.liftRepo.GetByID(*req.LiftID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to get lift")
			return
		}
		if lift == nil {
			writeError(w, http.StatusBadRequest, "Lift not found", "liftId references a non-existent lift")
			return
		}
	}

	// Check for duplicate
	existing, err := h.ppRepo.CheckDuplicate(programID, req.ProgressionID, req.LiftID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to check for duplicate")
		return
	}
	if existing != nil {
		writeError(w, http.StatusConflict, "Duplicate configuration", "A program progression with the same program, progression, and lift already exists")
		return
	}

	// Set defaults
	priority := int64(0)
	if req.Priority != nil {
		priority = *req.Priority
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	// Generate UUID and timestamps
	id := uuid.New().String()
	now := time.Now()

	entity := &repository.ProgramProgressionEntity{
		ID:                id,
		ProgramID:         programID,
		ProgressionID:     req.ProgressionID,
		LiftID:            req.LiftID,
		Priority:          priority,
		Enabled:           enabled,
		OverrideIncrement: req.OverrideIncrement,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	// Persist
	if err := h.ppRepo.Create(entity); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create program progression")
		return
	}

	writeJSON(w, http.StatusCreated, programProgressionEntityToResponse(entity))
}

// Update handles PUT /programs/{programId}/progressions/{configId}
func (h *ProgramProgressionHandler) Update(w http.ResponseWriter, r *http.Request) {
	programID := r.PathValue("programId")
	configID := r.PathValue("configId")

	if programID == "" {
		writeError(w, http.StatusBadRequest, "Missing program ID")
		return
	}
	if configID == "" {
		writeError(w, http.StatusBadRequest, "Missing config ID")
		return
	}

	// Get existing configuration
	existing, err := h.ppRepo.GetByID(configID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get program progression")
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "Program progression configuration not found")
		return
	}

	// Verify it belongs to the specified program
	if existing.ProgramID != programID {
		writeError(w, http.StatusNotFound, "Program progression configuration not found")
		return
	}

	var req UpdateProgramProgressionRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Apply updates
	if req.Priority != nil {
		existing.Priority = *req.Priority
	}
	if req.Enabled != nil {
		existing.Enabled = *req.Enabled
	}
	if req.OverrideIncrement != nil {
		existing.OverrideIncrement = req.OverrideIncrement
	}
	// Note: To clear overrideIncrement, the request should explicitly set it to null
	// This is handled by checking if the field was present in the JSON at all

	existing.UpdatedAt = time.Now()

	// Persist
	if err := h.ppRepo.Update(existing); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update program progression")
		return
	}

	writeJSON(w, http.StatusOK, programProgressionEntityToResponse(existing))
}

// Delete handles DELETE /programs/{programId}/progressions/{configId}
func (h *ProgramProgressionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	programID := r.PathValue("programId")
	configID := r.PathValue("configId")

	if programID == "" {
		writeError(w, http.StatusBadRequest, "Missing program ID")
		return
	}
	if configID == "" {
		writeError(w, http.StatusBadRequest, "Missing config ID")
		return
	}

	// Get existing configuration
	existing, err := h.ppRepo.GetByID(configID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get program progression")
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "Program progression configuration not found")
		return
	}

	// Verify it belongs to the specified program
	if existing.ProgramID != programID {
		writeError(w, http.StatusNotFound, "Program progression configuration not found")
		return
	}

	// Delete
	if err := h.ppRepo.Delete(configID); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete program progression")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// validateCreateRequest validates the create request.
func (h *ProgramProgressionHandler) validateCreateRequest(req *CreateProgramProgressionRequest) []string {
	var errors []string

	if req.ProgressionID == "" {
		errors = append(errors, "progressionId is required")
	}

	if req.Priority != nil && *req.Priority < 0 {
		errors = append(errors, "priority must be non-negative")
	}

	if req.OverrideIncrement != nil && *req.OverrideIncrement <= 0 {
		errors = append(errors, "overrideIncrement must be positive")
	}

	return errors
}

package api

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/userprogramstate"
	"github.com/waynenilsen/power-pro-v3/internal/middleware"
	"github.com/waynenilsen/power-pro-v3/internal/repository"
)

// EnrollmentHandler handles HTTP requests for user program enrollment operations.
type EnrollmentHandler struct {
	stateRepo   *repository.UserProgramStateRepository
	programRepo *repository.ProgramRepository
}

// NewEnrollmentHandler creates a new EnrollmentHandler.
func NewEnrollmentHandler(stateRepo *repository.UserProgramStateRepository, programRepo *repository.ProgramRepository) *EnrollmentHandler {
	return &EnrollmentHandler{
		stateRepo:   stateRepo,
		programRepo: programRepo,
	}
}

// EnrollmentProgramResponse represents program info in an enrollment response.
type EnrollmentProgramResponse struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Slug             string  `json:"slug"`
	Description      *string `json:"description,omitempty"`
	CycleLengthWeeks int     `json:"cycleLengthWeeks"`
}

// EnrollmentStateResponse represents the state portion of an enrollment response.
type EnrollmentStateResponse struct {
	CurrentWeek           int  `json:"currentWeek"`
	CurrentCycleIteration int  `json:"currentCycleIteration"`
	CurrentDayIndex       *int `json:"currentDayIndex,omitempty"`
}

// EnrollmentResponse represents the API response format for a user's program enrollment.
type EnrollmentResponse struct {
	ID         string                    `json:"id"`
	UserID     string                    `json:"userId"`
	Program    EnrollmentProgramResponse `json:"program"`
	State      EnrollmentStateResponse   `json:"state"`
	EnrolledAt time.Time                 `json:"enrolledAt"`
	UpdatedAt  time.Time                 `json:"updatedAt"`
}

// EnrollRequest represents the request body for enrolling a user in a program.
type EnrollRequest struct {
	ProgramID string `json:"programId"`
}

func enrollmentToResponse(e *userprogramstate.EnrollmentWithProgram) EnrollmentResponse {
	return EnrollmentResponse{
		ID:     e.State.ID,
		UserID: e.State.UserID,
		Program: EnrollmentProgramResponse{
			ID:               e.State.ProgramID,
			Name:             e.ProgramName,
			Slug:             e.ProgramSlug,
			Description:      e.ProgramDescription,
			CycleLengthWeeks: e.CycleLengthWeeks,
		},
		State: EnrollmentStateResponse{
			CurrentWeek:           e.State.CurrentWeek,
			CurrentCycleIteration: e.State.CurrentCycleIteration,
			CurrentDayIndex:       e.State.CurrentDayIndex,
		},
		EnrolledAt: e.State.EnrolledAt,
		UpdatedAt:  e.State.UpdatedAt,
	}
}

// Enroll handles POST /users/{userId}/program
func (h *EnrollmentHandler) Enroll(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "Missing user ID")
		return
	}

	// Authorization check: only the user themselves or an admin can enroll
	authUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)
	if authUserID != userID && !isAdmin {
		writeError(w, http.StatusForbidden, "Access denied: you can only manage your own enrollment")
		return
	}

	var req EnrollRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate program exists
	program, err := h.programRepo.GetByID(req.ProgramID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to verify program")
		return
	}
	if program == nil {
		writeError(w, http.StatusBadRequest, "Program not found", "program_id does not reference a valid program")
		return
	}

	// Check if user is already enrolled
	isEnrolled, err := h.stateRepo.UserIsEnrolled(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to check enrollment status")
		return
	}

	// If already enrolled, replace existing enrollment (re-enrollment replaces existing)
	if isEnrolled {
		// Delete existing enrollment
		if err := h.stateRepo.DeleteByUserID(userID); err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to remove existing enrollment")
			return
		}
	}

	// Generate UUID for new enrollment
	id := uuid.New().String()

	// Use domain logic to create enrollment
	input := userprogramstate.EnrollUserInput{
		UserID:    userID,
		ProgramID: req.ProgramID,
	}

	newState, result := userprogramstate.EnrollUser(input, id)
	if !result.Valid {
		details := make([]string, len(result.Errors))
		for i, err := range result.Errors {
			details[i] = err.Error()
		}
		writeError(w, http.StatusBadRequest, "Validation failed", details...)
		return
	}

	// Persist
	if err := h.stateRepo.Create(newState); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create enrollment")
		return
	}

	// Fetch the full enrollment with program details for response
	enrollment, err := h.stateRepo.GetEnrollmentWithProgram(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to retrieve enrollment details")
		return
	}
	if enrollment == nil {
		writeError(w, http.StatusInternalServerError, "Enrollment created but could not be retrieved")
		return
	}

	writeJSON(w, http.StatusCreated, enrollmentToResponse(enrollment))
}

// Get handles GET /users/{userId}/program
func (h *EnrollmentHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "Missing user ID")
		return
	}

	// Authorization check: only the user themselves or an admin can view enrollment
	authUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)
	if authUserID != userID && !isAdmin {
		writeError(w, http.StatusForbidden, "Access denied: you can only view your own enrollment")
		return
	}

	enrollment, err := h.stateRepo.GetEnrollmentWithProgram(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get enrollment")
		return
	}
	if enrollment == nil {
		writeError(w, http.StatusNotFound, "Not enrolled in any program")
		return
	}

	writeJSON(w, http.StatusOK, enrollmentToResponse(enrollment))
}

// Unenroll handles DELETE /users/{userId}/program
func (h *EnrollmentHandler) Unenroll(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "Missing user ID")
		return
	}

	// Authorization check: only the user themselves or an admin can unenroll
	authUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)
	if authUserID != userID && !isAdmin {
		writeError(w, http.StatusForbidden, "Access denied: you can only manage your own enrollment")
		return
	}

	// Check if enrolled
	isEnrolled, err := h.stateRepo.UserIsEnrolled(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to check enrollment status")
		return
	}
	if !isEnrolled {
		writeError(w, http.StatusNotFound, "Not enrolled in any program")
		return
	}

	// Delete enrollment
	if err := h.stateRepo.DeleteByUserID(userID); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to unenroll")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

package api

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/userprogramstate"
	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
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
		writeDomainError(w, apperrors.NewBadRequest("missing user ID"))
		return
	}

	// Authorization check: only the user themselves or an admin can enroll
	authUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)
	if authUserID != userID && !isAdmin {
		writeDomainError(w, apperrors.NewForbidden("you can only manage your own enrollment"))
		return
	}

	var req EnrollRequest
	if err := readJSON(r, &req); err != nil {
		writeDomainError(w, apperrors.NewBadRequest("invalid request body"))
		return
	}

	// Validate program exists
	program, err := h.programRepo.GetByID(req.ProgramID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to verify program", err))
		return
	}
	if program == nil {
		writeDomainError(w, apperrors.NewValidation("programId", "program not found"))
		return
	}

	// Check if user is already enrolled
	isEnrolled, err := h.stateRepo.UserIsEnrolled(userID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to check enrollment status", err))
		return
	}

	// If already enrolled, replace existing enrollment (re-enrollment replaces existing)
	if isEnrolled {
		// Delete existing enrollment
		if err := h.stateRepo.DeleteByUserID(userID); err != nil {
			writeDomainError(w, apperrors.NewInternal("failed to remove existing enrollment", err))
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
		writeDomainError(w, apperrors.NewValidationMsg("validation failed"), details...)
		return
	}

	// Persist
	if err := h.stateRepo.Create(newState); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to create enrollment", err))
		return
	}

	// Fetch the full enrollment with program details for response
	enrollment, err := h.stateRepo.GetEnrollmentWithProgram(userID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to retrieve enrollment details", err))
		return
	}
	if enrollment == nil {
		writeDomainError(w, apperrors.NewInternal("enrollment created but could not be retrieved", nil))
		return
	}

	writeData(w, http.StatusCreated, enrollmentToResponse(enrollment))
}

// Get handles GET /users/{userId}/program
func (h *EnrollmentHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing user ID"))
		return
	}

	// Authorization check: only the user themselves or an admin can view enrollment
	authUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)
	if authUserID != userID && !isAdmin {
		writeDomainError(w, apperrors.NewForbidden("you can only view your own enrollment"))
		return
	}

	enrollment, err := h.stateRepo.GetEnrollmentWithProgram(userID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get enrollment", err))
		return
	}
	if enrollment == nil {
		writeDomainError(w, apperrors.NewNotFound("enrollment", userID))
		return
	}

	writeData(w, http.StatusOK, enrollmentToResponse(enrollment))
}

// Unenroll handles DELETE /users/{userId}/program
func (h *EnrollmentHandler) Unenroll(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing user ID"))
		return
	}

	// Authorization check: only the user themselves or an admin can unenroll
	authUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)
	if authUserID != userID && !isAdmin {
		writeDomainError(w, apperrors.NewForbidden("you can only manage your own enrollment"))
		return
	}

	// Check if enrolled
	isEnrolled, err := h.stateRepo.UserIsEnrolled(userID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to check enrollment status", err))
		return
	}
	if !isEnrolled {
		writeDomainError(w, apperrors.NewNotFound("enrollment", userID))
		return
	}

	// Delete enrollment
	if err := h.stateRepo.DeleteByUserID(userID); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to unenroll", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

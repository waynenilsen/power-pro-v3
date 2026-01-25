package api

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/event"
	"github.com/waynenilsen/power-pro-v3/internal/domain/userprogramstate"
	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
	"github.com/waynenilsen/power-pro-v3/internal/middleware"
	"github.com/waynenilsen/power-pro-v3/internal/repository"
)

// EnrollmentHandler handles HTTP requests for user program enrollment operations.
type EnrollmentHandler struct {
	stateRepo   *repository.UserProgramStateRepository
	programRepo *repository.ProgramRepository
	sessionRepo *repository.WorkoutSessionRepository
	eventBus    *event.Bus
}

// NewEnrollmentHandler creates a new EnrollmentHandler.
func NewEnrollmentHandler(
	stateRepo *repository.UserProgramStateRepository,
	programRepo *repository.ProgramRepository,
	sessionRepo *repository.WorkoutSessionRepository,
	eventBus *event.Bus,
) *EnrollmentHandler {
	return &EnrollmentHandler{
		stateRepo:   stateRepo,
		programRepo: programRepo,
		sessionRepo: sessionRepo,
		eventBus:    eventBus,
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

// CurrentWorkoutSessionResponse represents the current workout session in an enrollment response.
type CurrentWorkoutSessionResponse struct {
	ID         string     `json:"id"`
	WeekNumber int        `json:"weekNumber"`
	DayIndex   int        `json:"dayIndex"`
	Status     string     `json:"status"`
	StartedAt  time.Time  `json:"startedAt"`
	FinishedAt *time.Time `json:"finishedAt,omitempty"`
}

// EnrollmentResponse represents the API response format for a user's program enrollment.
type EnrollmentResponse struct {
	ID                    string                         `json:"id"`
	UserID                string                         `json:"userId"`
	Program               EnrollmentProgramResponse      `json:"program"`
	State                 EnrollmentStateResponse        `json:"state"`
	EnrollmentStatus      string                         `json:"enrollmentStatus"`
	CycleStatus           string                         `json:"cycleStatus"`
	WeekStatus            string                         `json:"weekStatus"`
	CurrentWorkoutSession *CurrentWorkoutSessionResponse `json:"currentWorkoutSession"`
	EnrolledAt            time.Time                      `json:"enrolledAt"`
	UpdatedAt             time.Time                      `json:"updatedAt"`
}

// EnrollRequest represents the request body for enrolling a user in a program.
type EnrollRequest struct {
	ProgramID string `json:"programId"`
}

func enrollmentToResponse(e *userprogramstate.EnrollmentWithProgram, currentSession *CurrentWorkoutSessionResponse) EnrollmentResponse {
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
		EnrollmentStatus:      string(e.State.EnrollmentStatus),
		CycleStatus:           string(e.State.CycleStatus),
		WeekStatus:            string(e.State.WeekStatus),
		CurrentWorkoutSession: currentSession,
		EnrolledAt:            e.State.EnrolledAt,
		UpdatedAt:             e.State.UpdatedAt,
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

	// Emit ENROLLED event
	if h.eventBus != nil {
		evt := event.NewStateEvent(event.EventEnrolled, userID, req.ProgramID).
			WithPayload(event.PayloadEnrolledAt, newState.EnrolledAt)
		h.eventBus.PublishAsync(context.Background(), evt)
	}

	// New enrollment has no active workout session
	writeData(w, http.StatusCreated, enrollmentToResponse(enrollment, nil))
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

	// Fetch current workout session if any
	var currentSessionResponse *CurrentWorkoutSessionResponse
	if h.sessionRepo != nil {
		activeSession, err := h.sessionRepo.GetActiveByUserProgramStateID(enrollment.State.ID)
		if err != nil {
			writeDomainError(w, apperrors.NewInternal("failed to get active workout session", err))
			return
		}
		if activeSession != nil {
			currentSessionResponse = &CurrentWorkoutSessionResponse{
				ID:         activeSession.ID,
				WeekNumber: activeSession.WeekNumber,
				DayIndex:   activeSession.DayIndex,
				Status:     string(activeSession.Status),
				StartedAt:  activeSession.StartedAt,
				FinishedAt: activeSession.FinishedAt,
			}
		}
	}

	writeData(w, http.StatusOK, enrollmentToResponse(enrollment, currentSessionResponse))
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

	// Get current enrollment to capture cycles/weeks completed before deletion
	enrollment, err := h.stateRepo.GetEnrollmentWithProgram(userID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get enrollment", err))
		return
	}
	if enrollment == nil {
		writeDomainError(w, apperrors.NewNotFound("enrollment", userID))
		return
	}

	// Calculate cycles and weeks completed
	// Completed cycles = currentCycleIteration - 1 (if currently in a cycle, it's not complete)
	// But if enrollment status is BETWEEN_CYCLES or cycle status is COMPLETED, count current cycle
	cyclesCompleted := enrollment.State.CurrentCycleIteration - 1
	if enrollment.State.EnrollmentStatus == userprogramstate.EnrollmentStatusBetweenCycles ||
		enrollment.State.CycleStatus == userprogramstate.CycleStatusCompleted {
		cyclesCompleted = enrollment.State.CurrentCycleIteration
	}

	// Weeks completed = (completed cycles * weeks per cycle) + weeks in current cycle
	// If week status is COMPLETED, count the current week, otherwise current - 1
	weeksInCurrentCycle := enrollment.State.CurrentWeek - 1
	if enrollment.State.WeekStatus == userprogramstate.WeekStatusCompleted {
		weeksInCurrentCycle = enrollment.State.CurrentWeek
	}
	weeksCompleted := (cyclesCompleted * enrollment.CycleLengthWeeks) + weeksInCurrentCycle

	// Emit QUIT event before deletion
	if h.eventBus != nil {
		evt := event.NewStateEvent(event.EventQuit, userID, enrollment.State.ProgramID).
			WithPayload(event.PayloadCyclesCompleted, cyclesCompleted).
			WithPayload(event.PayloadWeeksCompleted, weeksCompleted)
		h.eventBus.PublishAsync(context.Background(), evt)
	}

	// Delete enrollment
	if err := h.stateRepo.DeleteByUserID(userID); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to unenroll", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// NextCycle handles POST /users/{userId}/enrollment/next-cycle
// Starts a new cycle when the enrollment is in BETWEEN_CYCLES state.
func (h *EnrollmentHandler) NextCycle(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing user ID"))
		return
	}

	// Authorization check
	authUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)
	if authUserID != userID && !isAdmin {
		writeDomainError(w, apperrors.NewForbidden("you can only manage your own enrollment"))
		return
	}

	// Get current enrollment
	enrollment, err := h.stateRepo.GetEnrollmentWithProgram(userID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get enrollment", err))
		return
	}
	if enrollment == nil {
		writeDomainError(w, apperrors.NewNotFound("enrollment", userID))
		return
	}

	// Validate enrollment is BETWEEN_CYCLES
	if enrollment.State.EnrollmentStatus != userprogramstate.EnrollmentStatusBetweenCycles {
		writeDomainError(w, apperrors.NewInvalidEnrollmentState("start new cycle", string(enrollment.State.EnrollmentStatus)))
		return
	}

	// Update state: set ACTIVE, increment cycle iteration, reset week to 1
	enrollment.State.EnrollmentStatus = userprogramstate.EnrollmentStatusActive
	enrollment.State.CurrentCycleIteration++
	enrollment.State.CurrentWeek = 1
	enrollment.State.CurrentDayIndex = nil
	enrollment.State.CycleStatus = userprogramstate.CycleStatusPending
	enrollment.State.WeekStatus = userprogramstate.WeekStatusPending
	enrollment.State.UpdatedAt = time.Now()

	// Persist changes
	if err := h.stateRepo.Update(enrollment.State); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to update enrollment", err))
		return
	}

	// Emit CYCLE_STARTED event
	if h.eventBus != nil {
		evt := event.NewStateEvent(event.EventCycleStarted, userID, enrollment.State.ProgramID).
			WithPayload(event.PayloadCycleIteration, enrollment.State.CurrentCycleIteration).
			WithPayload(event.PayloadWeekNumber, enrollment.State.CurrentWeek)
		h.eventBus.PublishAsync(context.Background(), evt)
	}

	// Fetch updated enrollment for response
	updatedEnrollment, err := h.stateRepo.GetEnrollmentWithProgram(userID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to retrieve updated enrollment", err))
		return
	}

	writeData(w, http.StatusOK, enrollmentToResponse(updatedEnrollment, nil))
}

// AdvanceWeek handles POST /users/{userId}/enrollment/advance-week
// Advances to the next week in the cycle. If at the final week, transitions to BETWEEN_CYCLES.
func (h *EnrollmentHandler) AdvanceWeek(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing user ID"))
		return
	}

	// Authorization check
	authUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)
	if authUserID != userID && !isAdmin {
		writeDomainError(w, apperrors.NewForbidden("you can only manage your own enrollment"))
		return
	}

	// Get current enrollment
	enrollment, err := h.stateRepo.GetEnrollmentWithProgram(userID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get enrollment", err))
		return
	}
	if enrollment == nil {
		writeDomainError(w, apperrors.NewNotFound("enrollment", userID))
		return
	}

	// Validate enrollment is ACTIVE
	if enrollment.State.EnrollmentStatus != userprogramstate.EnrollmentStatusActive {
		writeDomainError(w, apperrors.NewInvalidEnrollmentState("advance week", string(enrollment.State.EnrollmentStatus)))
		return
	}

	previousWeek := enrollment.State.CurrentWeek
	cycleBoundaryReached := false

	// Check if this is the final week of the cycle
	if enrollment.State.CurrentWeek >= enrollment.CycleLengthWeeks {
		// At final week - transition to BETWEEN_CYCLES
		enrollment.State.EnrollmentStatus = userprogramstate.EnrollmentStatusBetweenCycles
		enrollment.State.CycleStatus = userprogramstate.CycleStatusCompleted
		enrollment.State.WeekStatus = userprogramstate.WeekStatusCompleted
		cycleBoundaryReached = true
	} else {
		// Advance to next week
		enrollment.State.CurrentWeek++
		enrollment.State.CurrentDayIndex = nil
		enrollment.State.WeekStatus = userprogramstate.WeekStatusPending
	}

	enrollment.State.UpdatedAt = time.Now()

	// Persist changes
	if err := h.stateRepo.Update(enrollment.State); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to update enrollment", err))
		return
	}

	// Emit WEEK_COMPLETED event
	if h.eventBus != nil {
		evt := event.NewStateEvent(event.EventWeekCompleted, userID, enrollment.State.ProgramID).
			WithPayload(event.PayloadPreviousWeek, previousWeek).
			WithPayload(event.PayloadNewWeek, enrollment.State.CurrentWeek).
			WithPayload(event.PayloadCycleIteration, enrollment.State.CurrentCycleIteration)
		h.eventBus.PublishAsync(context.Background(), evt)

		// Emit CYCLE_BOUNDARY_REACHED if applicable
		if cycleBoundaryReached {
			boundaryEvt := event.NewStateEvent(event.EventCycleBoundaryReached, userID, enrollment.State.ProgramID).
				WithPayload(event.PayloadCompletedCycle, enrollment.State.CurrentCycleIteration).
				WithPayload(event.PayloadTotalWeeks, enrollment.CycleLengthWeeks)
			h.eventBus.PublishAsync(context.Background(), boundaryEvt)
		}
	}

	// Fetch updated enrollment for response
	updatedEnrollment, err := h.stateRepo.GetEnrollmentWithProgram(userID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to retrieve updated enrollment", err))
		return
	}

	writeData(w, http.StatusOK, enrollmentToResponse(updatedEnrollment, nil))
}

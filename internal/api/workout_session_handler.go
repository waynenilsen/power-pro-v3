package api

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/event"
	"github.com/waynenilsen/power-pro-v3/internal/domain/workoutsession"
	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
	"github.com/waynenilsen/power-pro-v3/internal/middleware"
	"github.com/waynenilsen/power-pro-v3/internal/repository"
)

// WorkoutSessionHandler handles HTTP requests for workout session operations.
type WorkoutSessionHandler struct {
	sessionRepo *repository.WorkoutSessionRepository
	stateRepo   *repository.UserProgramStateRepository
	eventBus    *event.Bus
}

// NewWorkoutSessionHandler creates a new WorkoutSessionHandler.
func NewWorkoutSessionHandler(
	sessionRepo *repository.WorkoutSessionRepository,
	stateRepo *repository.UserProgramStateRepository,
	eventBus *event.Bus,
) *WorkoutSessionHandler {
	return &WorkoutSessionHandler{
		sessionRepo: sessionRepo,
		stateRepo:   stateRepo,
		eventBus:    eventBus,
	}
}

// WorkoutSessionResponse represents the API response format for a workout session.
type WorkoutSessionResponse struct {
	ID                 string     `json:"id"`
	UserProgramStateID string     `json:"userProgramStateId"`
	WeekNumber         int        `json:"weekNumber"`
	DayIndex           int        `json:"dayIndex"`
	Status             string     `json:"status"`
	StartedAt          time.Time  `json:"startedAt"`
	FinishedAt         *time.Time `json:"finishedAt,omitempty"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
}

// StartWorkoutRequest represents the request body for starting a workout.
type StartWorkoutRequest struct {
	// No fields required - uses current state from enrollment
}

func workoutSessionToResponse(ws *workoutsession.WorkoutSession) WorkoutSessionResponse {
	return WorkoutSessionResponse{
		ID:                 ws.ID,
		UserProgramStateID: ws.UserProgramStateID,
		WeekNumber:         ws.WeekNumber,
		DayIndex:           ws.DayIndex,
		Status:             string(ws.Status),
		StartedAt:          ws.StartedAt,
		FinishedAt:         ws.FinishedAt,
		CreatedAt:          ws.CreatedAt,
		UpdatedAt:          ws.UpdatedAt,
	}
}

// Start handles POST /workouts/start
// Starts a new workout session for the authenticated user.
func (h *WorkoutSessionHandler) Start(w http.ResponseWriter, r *http.Request) {
	authUserID := middleware.GetUserID(r)
	if authUserID == "" {
		writeDomainError(w, apperrors.NewUnauthorized("authentication required"))
		return
	}

	// Get the user's current enrollment/state
	enrollment, err := h.stateRepo.GetEnrollmentWithProgram(authUserID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get enrollment", err))
		return
	}
	if enrollment == nil {
		writeDomainError(w, apperrors.NewBadRequest("user is not enrolled in a program"))
		return
	}

	// Check if user already has an IN_PROGRESS session
	activeSession, err := h.sessionRepo.GetActiveByUserProgramStateID(enrollment.State.ID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to check for active session", err))
		return
	}
	if activeSession != nil {
		writeDomainError(w, apperrors.NewConflict("workout session already in progress"))
		return
	}

	// Determine day index - if nil, start at 0
	dayIndex := 0
	if enrollment.State.CurrentDayIndex != nil {
		dayIndex = *enrollment.State.CurrentDayIndex
	}

	// Create the session
	id := uuid.New().String()
	session, result := workoutsession.NewWorkoutSession(
		workoutsession.NewWorkoutSessionInput{
			UserProgramStateID: enrollment.State.ID,
			WeekNumber:         enrollment.State.CurrentWeek,
			DayIndex:           dayIndex,
		},
		id,
	)
	if !result.Valid {
		details := make([]string, len(result.Errors))
		for i, err := range result.Errors {
			details[i] = err.Error()
		}
		writeDomainError(w, apperrors.NewValidationMsg("validation failed"), details...)
		return
	}

	// Persist
	if err := h.sessionRepo.Create(session); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to create workout session", err))
		return
	}

	// Emit WORKOUT_STARTED event
	if h.eventBus != nil {
		evt := event.NewStateEvent(event.EventWorkoutStarted, authUserID, enrollment.State.ProgramID).
			WithPayload(event.PayloadSessionID, session.ID).
			WithPayload(event.PayloadWeekNumber, session.WeekNumber).
			WithPayload(event.PayloadDaySlug, session.DayIndex)
		h.eventBus.PublishAsync(context.Background(), evt)
	}

	writeData(w, http.StatusCreated, workoutSessionToResponse(session))
}

// Get handles GET /workouts/{id}
// Returns workout session details.
func (h *WorkoutSessionHandler) Get(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("id")
	if sessionID == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing session ID"))
		return
	}

	session, err := h.sessionRepo.GetByID(sessionID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get session", err))
		return
	}
	if session == nil {
		writeDomainError(w, apperrors.NewNotFound("workout session", sessionID))
		return
	}

	// Authorization: check if user owns this session or is admin
	authUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)

	// Get the user program state to check ownership
	state, err := h.stateRepo.GetByID(session.UserProgramStateID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get program state", err))
		return
	}
	if state == nil {
		writeDomainError(w, apperrors.NewInternal("session references invalid program state", nil))
		return
	}

	if state.UserID != authUserID && !isAdmin {
		writeDomainError(w, apperrors.NewForbidden("you can only view your own workout sessions"))
		return
	}

	writeData(w, http.StatusOK, workoutSessionToResponse(session))
}

// Finish handles POST /workouts/{id}/finish
// Marks a workout session as completed.
func (h *WorkoutSessionHandler) Finish(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("id")
	if sessionID == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing session ID"))
		return
	}

	session, err := h.sessionRepo.GetByID(sessionID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get session", err))
		return
	}
	if session == nil {
		writeDomainError(w, apperrors.NewNotFound("workout session", sessionID))
		return
	}

	// Authorization check
	authUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)

	state, err := h.stateRepo.GetByID(session.UserProgramStateID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get program state", err))
		return
	}
	if state == nil {
		writeDomainError(w, apperrors.NewInternal("session references invalid program state", nil))
		return
	}

	if state.UserID != authUserID && !isAdmin {
		writeDomainError(w, apperrors.NewForbidden("you can only manage your own workout sessions"))
		return
	}

	// Complete the session
	if err := session.Complete(); err != nil {
		switch err {
		case workoutsession.ErrAlreadyCompleted:
			writeDomainError(w, apperrors.NewConflict("session already completed"))
		case workoutsession.ErrNotInProgress:
			writeDomainError(w, apperrors.NewBadRequest("session is not in progress"))
		default:
			writeDomainError(w, apperrors.NewInternal("failed to complete session", err))
		}
		return
	}

	// Persist
	if err := h.sessionRepo.Complete(session); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to save completed session", err))
		return
	}

	// Emit WORKOUT_COMPLETED event
	if h.eventBus != nil {
		evt := event.NewStateEvent(event.EventWorkoutCompleted, state.UserID, state.ProgramID).
			WithPayload(event.PayloadSessionID, session.ID).
			WithPayload(event.PayloadWeekNumber, session.WeekNumber).
			WithPayload(event.PayloadDaySlug, session.DayIndex)
		h.eventBus.PublishAsync(context.Background(), evt)
	}

	writeData(w, http.StatusOK, workoutSessionToResponse(session))
}

// Abandon handles POST /workouts/{id}/abandon
// Marks a workout session as abandoned.
func (h *WorkoutSessionHandler) Abandon(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("id")
	if sessionID == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing session ID"))
		return
	}

	session, err := h.sessionRepo.GetByID(sessionID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get session", err))
		return
	}
	if session == nil {
		writeDomainError(w, apperrors.NewNotFound("workout session", sessionID))
		return
	}

	// Authorization check
	authUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)

	state, err := h.stateRepo.GetByID(session.UserProgramStateID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get program state", err))
		return
	}
	if state == nil {
		writeDomainError(w, apperrors.NewInternal("session references invalid program state", nil))
		return
	}

	if state.UserID != authUserID && !isAdmin {
		writeDomainError(w, apperrors.NewForbidden("you can only manage your own workout sessions"))
		return
	}

	// Abandon the session
	if err := session.Abandon(); err != nil {
		switch err {
		case workoutsession.ErrAlreadyAbandoned:
			writeDomainError(w, apperrors.NewConflict("session already abandoned"))
		case workoutsession.ErrNotInProgress:
			writeDomainError(w, apperrors.NewBadRequest("session is not in progress"))
		default:
			writeDomainError(w, apperrors.NewInternal("failed to abandon session", err))
		}
		return
	}

	// Persist
	if err := h.sessionRepo.Abandon(session); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to save abandoned session", err))
		return
	}

	// Emit WORKOUT_ABANDONED event
	if h.eventBus != nil {
		evt := event.NewStateEvent(event.EventWorkoutAbandoned, state.UserID, state.ProgramID).
			WithPayload(event.PayloadSessionID, session.ID).
			WithPayload(event.PayloadWeekNumber, session.WeekNumber).
			WithPayload(event.PayloadDaySlug, session.DayIndex)
		h.eventBus.PublishAsync(context.Background(), evt)
	}

	writeData(w, http.StatusOK, workoutSessionToResponse(session))
}

// ListByUser handles GET /users/{id}/workouts
// Lists a user's workout history with pagination.
func (h *WorkoutSessionHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	if userID == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing user ID"))
		return
	}

	// Authorization check
	authUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)
	if authUserID != userID && !isAdmin {
		writeDomainError(w, apperrors.NewForbidden("you can only view your own workouts"))
		return
	}

	// Parse pagination
	pagination := ParsePagination(r.URL.Query())

	// Parse optional status filter
	status, err := ParseFilterEnum(r.URL.Query(), "status", []string{"IN_PROGRESS", "COMPLETED", "ABANDONED"})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	// Get total count
	total, err := h.sessionRepo.CountByUserID(userID, status)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to count sessions", err))
		return
	}

	// Get sessions
	sessions, err := h.sessionRepo.GetByUserID(userID, status, pagination.Limit, pagination.Offset)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get sessions", err))
		return
	}

	// Convert to response format
	responses := make([]WorkoutSessionResponse, len(sessions))
	for i, session := range sessions {
		responses[i] = workoutSessionToResponse(session)
	}

	writePaginatedData(w, http.StatusOK, responses, total, pagination.Limit, pagination.Offset)
}

// GetCurrentByUser handles GET /users/{id}/workouts/current
// Gets the user's current in-progress workout if any.
func (h *WorkoutSessionHandler) GetCurrentByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	if userID == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing user ID"))
		return
	}

	// Authorization check
	authUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)
	if authUserID != userID && !isAdmin {
		writeDomainError(w, apperrors.NewForbidden("you can only view your own workouts"))
		return
	}

	// Get active session
	session, err := h.sessionRepo.GetActiveByUserID(userID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get active session", err))
		return
	}
	if session == nil {
		writeDomainError(w, apperrors.NewNotFound("active workout session", userID))
		return
	}

	writeData(w, http.StatusOK, workoutSessionToResponse(session))
}

package api

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/event"
	"github.com/waynenilsen/power-pro-v3/internal/domain/loggedset"
	"github.com/waynenilsen/power-pro-v3/internal/domain/workoutsession"
	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
	"github.com/waynenilsen/power-pro-v3/internal/middleware"
	"github.com/waynenilsen/power-pro-v3/internal/repository"
	"github.com/waynenilsen/power-pro-v3/internal/service"
)

// LoggedSetHandler handles HTTP requests for logged set operations.
type LoggedSetHandler struct {
	repo               *repository.LoggedSetRepository
	workoutSessionRepo *repository.WorkoutSessionRepository
	stateRepo          *repository.UserProgramStateRepository
	failureService     *service.FailureService
	eventBus           *event.Bus
}

// NewLoggedSetHandler creates a new LoggedSetHandler.
func NewLoggedSetHandler(
	repo *repository.LoggedSetRepository,
	workoutSessionRepo *repository.WorkoutSessionRepository,
	stateRepo *repository.UserProgramStateRepository,
	failureService *service.FailureService,
	eventBus *event.Bus,
) *LoggedSetHandler {
	return &LoggedSetHandler{
		repo:               repo,
		workoutSessionRepo: workoutSessionRepo,
		stateRepo:          stateRepo,
		failureService:     failureService,
		eventBus:           eventBus,
	}
}

// LoggedSetResponse represents the API response format for a logged set.
type LoggedSetResponse struct {
	ID             string    `json:"id"`
	UserID         string    `json:"userId"`
	SessionID      string    `json:"sessionId"`
	PrescriptionID string    `json:"prescriptionId"`
	LiftID         string    `json:"liftId"`
	SetNumber      int       `json:"setNumber"`
	Weight         float64   `json:"weight"`
	TargetReps     int       `json:"targetReps"`
	RepsPerformed  int       `json:"repsPerformed"`
	IsAMRAP        bool      `json:"isAmrap"`
	RPE            *float64  `json:"rpe,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
}

// CreateLoggedSetRequest represents a single logged set in the batch request.
type CreateLoggedSetRequest struct {
	PrescriptionID string   `json:"prescriptionId"`
	LiftID         string   `json:"liftId"`
	SetNumber      int      `json:"setNumber"`
	Weight         float64  `json:"weight"`
	TargetReps     int      `json:"targetReps"`
	RepsPerformed  int      `json:"repsPerformed"`
	IsAMRAP        bool     `json:"isAmrap"`
	RPE            *float64 `json:"rpe,omitempty"`
}

// CreateLoggedSetsBatchRequest represents the request body for creating logged sets.
type CreateLoggedSetsBatchRequest struct {
	Sets []CreateLoggedSetRequest `json:"sets"`
}

func loggedSetToResponse(ls *loggedset.LoggedSet) LoggedSetResponse {
	return LoggedSetResponse{
		ID:             ls.ID,
		UserID:         ls.UserID,
		SessionID:      ls.SessionID,
		PrescriptionID: ls.PrescriptionID,
		LiftID:         ls.LiftID,
		SetNumber:      ls.SetNumber,
		Weight:         ls.Weight,
		TargetReps:     ls.TargetReps,
		RepsPerformed:  ls.RepsPerformed,
		IsAMRAP:        ls.IsAMRAP,
		RPE:            ls.RPE,
		CreatedAt:      ls.CreatedAt,
	}
}

// CreateBatch handles POST /sessions/{sessionId}/sets
func (h *LoggedSetHandler) CreateBatch(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionId")
	if sessionID == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing session ID"))
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(r)
	if userID == "" {
		writeDomainError(w, apperrors.NewUnauthorized("authentication required"))
		return
	}

	// Validate that the workout session exists and is IN_PROGRESS
	var programID string
	if h.workoutSessionRepo != nil {
		session, err := h.workoutSessionRepo.GetByID(sessionID)
		if err != nil {
			writeDomainError(w, apperrors.NewInternal("failed to get workout session", err))
			return
		}
		if session == nil {
			writeDomainError(w, apperrors.NewNotFound("workout session", sessionID))
			return
		}
		if session.Status != workoutsession.StatusInProgress {
			writeDomainError(w, apperrors.NewSessionNotActive(string(session.Status)))
			return
		}

		// Get program ID for event emission
		if h.stateRepo != nil {
			state, err := h.stateRepo.GetByID(session.UserProgramStateID)
			if err == nil && state != nil {
				programID = state.ProgramID
			}
		}
	}

	var req CreateLoggedSetsBatchRequest
	if err := readJSON(r, &req); err != nil {
		writeDomainError(w, apperrors.NewBadRequest("invalid request body"))
		return
	}

	if len(req.Sets) == 0 {
		writeDomainError(w, apperrors.NewBadRequest("at least one set is required"))
		return
	}

	responses := make([]LoggedSetResponse, 0, len(req.Sets))

	for i, setReq := range req.Sets {
		id := uuid.New().String()

		input := loggedset.CreateLoggedSetInput{
			UserID:         userID,
			SessionID:      sessionID,
			PrescriptionID: setReq.PrescriptionID,
			LiftID:         setReq.LiftID,
			SetNumber:      setReq.SetNumber,
			Weight:         setReq.Weight,
			TargetReps:     setReq.TargetReps,
			RepsPerformed:  setReq.RepsPerformed,
			IsAMRAP:        setReq.IsAMRAP,
			RPE:            setReq.RPE,
		}

		newSet, result := loggedset.NewLoggedSet(input, id)
		if !result.Valid {
			details := make([]string, len(result.Errors))
			for j, err := range result.Errors {
				details[j] = err.Error()
			}
			writeDomainError(w, apperrors.NewValidationMsg("validation failed for set "+string(rune('0'+i+1))), details...)
			return
		}

		if err := h.repo.Create(newSet); err != nil {
			writeDomainError(w, apperrors.NewInternal("failed to create logged set", err))
			return
		}

		// Process the logged set for failure tracking (if FailureService is configured)
		if h.failureService != nil {
			_, _ = h.failureService.ProcessLoggedSet(r.Context(), newSet)
			// Note: We don't fail the request if failure processing fails,
			// as the set has been successfully logged. Failure processing is
			// best-effort and logged separately.
		}

		// Emit SET_LOGGED event
		if h.eventBus != nil {
			isFailure := newSet.RepsPerformed < newSet.TargetReps
			evt := event.NewStateEvent(event.EventSetLogged, userID, programID).
				WithPayload(event.PayloadLoggedSetID, newSet.ID).
				WithPayload(event.PayloadSessionID, sessionID).
				WithPayload(event.PayloadLiftID, newSet.LiftID).
				WithPayload(event.PayloadRepsPerformed, newSet.RepsPerformed).
				WithPayload(event.PayloadTargetReps, newSet.TargetReps).
				WithPayload(event.PayloadWeight, newSet.Weight).
				WithPayload(event.PayloadIsAMRAP, newSet.IsAMRAP).
				WithPayload(event.PayloadIsFailure, isFailure)
			h.eventBus.PublishAsync(context.Background(), evt)
		}

		responses = append(responses, loggedSetToResponse(newSet))
	}

	writeData(w, http.StatusCreated, responses)
}

// ListBySession handles GET /sessions/{sessionId}/sets
func (h *LoggedSetHandler) ListBySession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionId")
	if sessionID == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing session ID"))
		return
	}

	sets, err := h.repo.ListBySession(sessionID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to list logged sets", err))
		return
	}

	data := make([]LoggedSetResponse, len(sets))
	for i, s := range sets {
		data[i] = loggedSetToResponse(&s)
	}

	writeData(w, http.StatusOK, data)
}

// ListByUser handles GET /users/{userId}/logged-sets
func (h *LoggedSetHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing user ID"))
		return
	}

	// Get current user ID from context
	currentUserID := middleware.GetUserID(r)
	if currentUserID == "" {
		writeDomainError(w, apperrors.NewUnauthorized("authentication required"))
		return
	}

	// Check if the requesting user can access this data
	// Users can only access their own logged sets (or admins can access anyone's)
	isAdmin := middleware.IsAdmin(r)
	if currentUserID != userID && !isAdmin {
		writeDomainError(w, apperrors.NewForbidden("cannot access other user's logged sets"))
		return
	}

	// Parse pagination
	query := r.URL.Query()
	pg := ParsePagination(query)

	params := repository.LoggedSetListParams{
		UserID: userID,
		Limit:  int64(pg.Limit),
		Offset: int64(pg.Offset),
	}

	sets, total, err := h.repo.ListByUser(params)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to list logged sets", err))
		return
	}

	data := make([]LoggedSetResponse, len(sets))
	for i, s := range sets {
		data[i] = loggedSetToResponse(&s)
	}

	writePaginatedData(w, http.StatusOK, data, total, pg.Limit, pg.Offset)
}

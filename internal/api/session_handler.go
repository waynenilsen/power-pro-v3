package api

import (
	"errors"
	"net/http"

	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
	"github.com/waynenilsen/power-pro-v3/internal/middleware"
	"github.com/waynenilsen/power-pro-v3/internal/service"
)

// SessionHandler handles HTTP requests for workout session operations.
type SessionHandler struct {
	sessionService *service.SessionService
}

// NewSessionHandler creates a new SessionHandler.
func NewSessionHandler(sessionService *service.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: sessionService}
}

// NextSetResponse represents the API response for a next set request.
type NextSetResponse struct {
	NextSet            *NextSetInfo `json:"nextSet,omitempty"`
	IsComplete         bool         `json:"isComplete"`
	TotalSetsCompleted int          `json:"totalSetsCompleted"`
	TotalRepsCompleted int          `json:"totalRepsCompleted"`
	TerminationReason  string       `json:"terminationReason,omitempty"`
}

// NextSetInfo represents a generated set in the API response.
type NextSetInfo struct {
	SetNumber  int     `json:"setNumber"`
	Weight     float64 `json:"weight"`
	TargetReps int     `json:"targetReps"`
	IsWorkSet  bool    `json:"isWorkSet"`
}

// GetNextSet handles GET /sessions/{sessionId}/prescriptions/{prescriptionId}/next-set
// Returns the next set to perform for a variable scheme prescription.
func (h *SessionHandler) GetNextSet(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionId")
	if sessionID == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing session ID"))
		return
	}

	prescriptionID := r.PathValue("prescriptionId")
	if prescriptionID == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing prescription ID"))
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(r)
	if userID == "" {
		writeDomainError(w, apperrors.NewUnauthorized("authentication required"))
		return
	}

	req := service.NextSetRequest{
		SessionID:      sessionID,
		PrescriptionID: prescriptionID,
		UserID:         userID,
	}

	result, err := h.sessionService.GetNextSet(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrPrescriptionNotFound) {
			writeDomainError(w, apperrors.NewNotFound("prescription", prescriptionID))
			return
		}
		if errors.Is(err, service.ErrNotVariableScheme) {
			writeDomainError(w, apperrors.NewBadRequest("prescription does not use a variable set scheme"))
			return
		}
		if errors.Is(err, service.ErrNoSetsLogged) {
			writeDomainError(w, apperrors.NewBadRequest("no sets logged yet - log the first set before requesting next set"))
			return
		}
		writeDomainError(w, apperrors.NewInternal("failed to get next set", err))
		return
	}

	response := NextSetResponse{
		IsComplete:         result.IsComplete,
		TotalSetsCompleted: result.TotalSetsCompleted,
		TotalRepsCompleted: result.TotalRepsCompleted,
		TerminationReason:  result.TerminationReason,
	}

	if result.NextSet != nil {
		response.NextSet = &NextSetInfo{
			SetNumber:  result.NextSet.SetNumber,
			Weight:     result.NextSet.Weight,
			TargetReps: result.NextSet.TargetReps,
			IsWorkSet:  result.NextSet.IsWorkSet,
		}
	}

	writeData(w, http.StatusOK, response)
}

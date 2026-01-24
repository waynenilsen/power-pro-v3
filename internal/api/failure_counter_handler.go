// Package api provides HTTP handlers for the API.
// This file implements the FailureCounterHandler for querying failure counter state.
package api

import (
	"net/http"

	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
	"github.com/waynenilsen/power-pro-v3/internal/middleware"
	"github.com/waynenilsen/power-pro-v3/internal/service"
)

// FailureCounterHandler handles HTTP requests for failure counter operations.
type FailureCounterHandler struct {
	failureService *service.FailureService
}

// NewFailureCounterHandler creates a new FailureCounterHandler.
func NewFailureCounterHandler(failureService *service.FailureService) *FailureCounterHandler {
	return &FailureCounterHandler{failureService: failureService}
}

// FailureCounterResponse represents the API response for failure counter state.
type FailureCounterResponse struct {
	UserID              string `json:"userId"`
	LiftID              string `json:"liftId"`
	ProgressionID       string `json:"progressionId"`
	ConsecutiveFailures int    `json:"consecutiveFailures"`
}

// Get handles GET /users/{userId}/failure-counters
// Query parameters:
//   - liftId: required - the lift ID to query
//   - progressionId: required - the progression ID to query
func (h *FailureCounterHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing user ID"))
		return
	}

	// Authorization: user can query their own failure counters, admin can query any
	authUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)
	if !isAdmin && authUserID != userID {
		writeDomainError(w, apperrors.NewForbidden("you do not have permission to view failure counters for this user"))
		return
	}

	// Get query parameters
	liftID := r.URL.Query().Get("liftId")
	progressionID := r.URL.Query().Get("progressionId")

	if liftID == "" {
		writeDomainError(w, apperrors.NewBadRequest("liftId query parameter is required"))
		return
	}
	if progressionID == "" {
		writeDomainError(w, apperrors.NewBadRequest("progressionId query parameter is required"))
		return
	}

	// Get the failure count
	count, err := h.failureService.GetFailureCount(userID, liftID, progressionID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get failure count", err))
		return
	}

	response := FailureCounterResponse{
		UserID:              userID,
		LiftID:              liftID,
		ProgressionID:       progressionID,
		ConsecutiveFailures: count,
	}

	writeData(w, http.StatusOK, response)
}

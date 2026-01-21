// Package api provides HTTP handlers for the API.
package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/domain/progression"
	"github.com/waynenilsen/power-pro-v3/internal/middleware"
	"github.com/waynenilsen/power-pro-v3/internal/service"
)

// ManualTriggerHandler handles HTTP requests for manual progression triggering.
type ManualTriggerHandler struct {
	progressionService *service.ProgressionService
}

// NewManualTriggerHandler creates a new ManualTriggerHandler.
func NewManualTriggerHandler(progressionService *service.ProgressionService) *ManualTriggerHandler {
	return &ManualTriggerHandler{progressionService: progressionService}
}

// TriggerRequest represents the request body for manual progression trigger.
type TriggerRequest struct {
	ProgressionID string `json:"progressionId"`
	LiftID        string `json:"liftId,omitempty"`
	Force         bool   `json:"force"`
}

// TriggerResultResponse represents a single progression result in the API response.
type TriggerResultResponse struct {
	ProgressionID string                   `json:"progressionId"`
	LiftID        string                   `json:"liftId"`
	Applied       bool                     `json:"applied"`
	Skipped       bool                     `json:"skipped,omitempty"`
	SkipReason    string                   `json:"skipReason,omitempty"`
	Result        *ProgressionResultDetail `json:"result,omitempty"`
	Error         string                   `json:"error,omitempty"`
}

// ProgressionResultDetail contains the details of an applied progression.
type ProgressionResultDetail struct {
	PreviousValue float64   `json:"previousValue"`
	NewValue      float64   `json:"newValue"`
	Delta         float64   `json:"delta"`
	MaxType       string    `json:"maxType"`
	AppliedAt     time.Time `json:"appliedAt"`
}

// TriggerResponse represents the response for manual progression trigger.
type TriggerResponse struct {
	Results      []TriggerResultResponse `json:"results"`
	TotalApplied int                     `json:"totalApplied"`
	TotalSkipped int                     `json:"totalSkipped"`
	TotalErrors  int                     `json:"totalErrors"`
}

// Trigger handles POST /users/{userId}/progressions/trigger
func (h *ManualTriggerHandler) Trigger(w http.ResponseWriter, r *http.Request) {
	// Get userId from path
	userID := r.PathValue("userId")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "Missing user ID")
		return
	}

	// Authorization: user can trigger their own progressions, admin can trigger for any user
	authUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)
	if !isAdmin && authUserID != userID {
		writeError(w, http.StatusForbidden, "Access denied: you do not have permission to trigger progressions for this user")
		return
	}

	// Parse request body
	var req TriggerRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Validate required fields
	if req.ProgressionID == "" {
		writeError(w, http.StatusBadRequest, "progressionId is required")
		return
	}

	// Apply progression manually
	result, err := h.progressionService.ApplyProgressionManually(r.Context(), userID, req.ProgressionID, req.LiftID, req.Force)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrProgressionNotFound):
			writeError(w, http.StatusNotFound, "Progression not found")
		case errors.Is(err, service.ErrLiftNotFound):
			writeError(w, http.StatusNotFound, "Lift not found")
		case errors.Is(err, service.ErrUserNotEnrolled):
			writeError(w, http.StatusBadRequest, "User is not enrolled in any program")
		case errors.Is(err, service.ErrNoApplicableProgressions):
			writeError(w, http.StatusBadRequest, "No applicable progressions found", err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "Failed to trigger progression", err.Error())
		}
		return
	}

	// Convert to API response format
	response := TriggerResponse{
		Results:      make([]TriggerResultResponse, len(result.Results)),
		TotalApplied: result.TotalApplied,
		TotalSkipped: result.TotalSkipped,
		TotalErrors:  result.TotalErrors,
	}

	for i, tr := range result.Results {
		resp := TriggerResultResponse{
			ProgressionID: tr.ProgressionID,
			LiftID:        tr.LiftID,
			Applied:       tr.Applied,
			Skipped:       tr.Skipped,
			SkipReason:    tr.SkipReason,
			Error:         tr.Error,
		}

		if tr.Result != nil {
			resp.Result = &ProgressionResultDetail{
				PreviousValue: tr.Result.PreviousValue,
				NewValue:      tr.Result.NewValue,
				Delta:         tr.Result.Delta,
				MaxType:       string(progression.TrainingMax), // Default, actual maxType is in result
				AppliedAt:     tr.Result.AppliedAt,
			}
		}

		response.Results[i] = resp
	}

	writeJSON(w, http.StatusOK, response)
}

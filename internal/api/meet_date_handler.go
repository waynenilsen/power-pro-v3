package api

import (
	"net/http"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/domain/userprogramstate"
	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
	"github.com/waynenilsen/power-pro-v3/internal/middleware"
	"github.com/waynenilsen/power-pro-v3/internal/repository"
)

// MeetDateHandler handles HTTP requests for meet date management.
type MeetDateHandler struct {
	stateRepo *repository.UserProgramStateRepository
}

// NewMeetDateHandler creates a new MeetDateHandler.
func NewMeetDateHandler(stateRepo *repository.UserProgramStateRepository) *MeetDateHandler {
	return &MeetDateHandler{
		stateRepo: stateRepo,
	}
}

// SetMeetDateRequest represents the request body for setting a meet date.
type SetMeetDateRequest struct {
	MeetDate *string `json:"meet_date"` // ISO 8601 date string, or null to clear
}

// MeetDateResponse represents the response for meet date operations.
type MeetDateResponse struct {
	MeetDate     *string `json:"meet_date,omitempty"`
	DaysOut      int     `json:"days_out"`
	CurrentPhase string  `json:"current_phase"`
	WeeksToMeet  int     `json:"weeks_to_meet"`
}

// CountdownResponse represents the response for the countdown endpoint.
type CountdownResponse struct {
	MeetDate         *string `json:"meet_date,omitempty"`
	DaysOut          int     `json:"days_out"`
	CurrentPhase     string  `json:"current_phase"`
	PhaseWeek        int     `json:"phase_week"`
	TaperMultiplier  float64 `json:"taper_multiplier"`
}

// SetMeetDate handles PUT /users/{userId}/programs/{programId}/state/meet-date
func (h *MeetDateHandler) SetMeetDate(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing user ID"))
		return
	}

	// programId is in the path but we don't use it for now since users can only have one enrollment
	// This matches the API design pattern for future multi-program support

	// Authorization check: only the user themselves or an admin can set meet date
	authUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)
	if authUserID != userID && !isAdmin {
		writeDomainError(w, apperrors.NewForbidden("you can only manage your own program state"))
		return
	}

	var req SetMeetDateRequest
	if err := readJSON(r, &req); err != nil {
		writeDomainError(w, apperrors.NewBadRequest("invalid request body"))
		return
	}

	// Get user's current state
	state, err := h.stateRepo.GetByUserID(userID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get user state", err))
		return
	}
	if state == nil {
		writeDomainError(w, apperrors.NewNotFound("enrollment", userID))
		return
	}

	// Parse meet date if provided
	var meetDate *time.Time
	if req.MeetDate != nil && *req.MeetDate != "" {
		// Try RFC3339 first (with time), then date-only format
		parsedTime, err := time.Parse(time.RFC3339, *req.MeetDate)
		if err != nil {
			// Try date-only format
			parsedTime, err = time.Parse("2006-01-02", *req.MeetDate)
			if err != nil {
				writeDomainError(w, apperrors.NewValidation("meet_date", "invalid date format; use ISO 8601 (e.g., 2024-06-15 or 2024-06-15T00:00:00Z)"))
				return
			}
		}
		meetDate = &parsedTime
	}

	// Update meet date using domain logic
	input := userprogramstate.UpdateMeetDateInput{
		MeetDate: meetDate,
	}
	result := state.UpdateMeetDate(input)
	if !result.Valid {
		details := make([]string, len(result.Errors))
		for i, err := range result.Errors {
			details[i] = err.Error()
		}
		writeDomainError(w, apperrors.NewValidationMsg("validation failed"), details...)
		return
	}

	// Persist the updated state
	if err := h.stateRepo.Update(state); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to update state", err))
		return
	}

	// Build response
	response := h.buildMeetDateResponse(state)
	writeData(w, http.StatusOK, response)
}

// GetCountdown handles GET /users/{userId}/programs/{programId}/state/countdown
func (h *MeetDateHandler) GetCountdown(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing user ID"))
		return
	}

	// Authorization check: only the user themselves or an admin can view countdown
	authUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)
	if authUserID != userID && !isAdmin {
		writeDomainError(w, apperrors.NewForbidden("you can only view your own program state"))
		return
	}

	// Get user's current state
	state, err := h.stateRepo.GetByUserID(userID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get user state", err))
		return
	}
	if state == nil {
		writeDomainError(w, apperrors.NewNotFound("enrollment", userID))
		return
	}

	// Build response
	response := h.buildCountdownResponse(state)
	writeData(w, http.StatusOK, response)
}

// buildMeetDateResponse builds the MeetDateResponse from state.
func (h *MeetDateHandler) buildMeetDateResponse(state *userprogramstate.UserProgramState) MeetDateResponse {
	var meetDateStr *string
	if state.MeetDate != nil {
		str := state.MeetDate.Format("2006-01-02")
		meetDateStr = &str
	}

	return MeetDateResponse{
		MeetDate:     meetDateStr,
		DaysOut:      state.DaysOut(),
		CurrentPhase: h.determinePhase(state),
		WeeksToMeet:  state.WeeksToMeet(),
	}
}

// buildCountdownResponse builds the CountdownResponse from state.
func (h *MeetDateHandler) buildCountdownResponse(state *userprogramstate.UserProgramState) CountdownResponse {
	var meetDateStr *string
	if state.MeetDate != nil {
		str := state.MeetDate.Format("2006-01-02")
		meetDateStr = &str
	}

	return CountdownResponse{
		MeetDate:        meetDateStr,
		DaysOut:         state.DaysOut(),
		CurrentPhase:    h.determinePhase(state),
		PhaseWeek:       state.CurrentWeek,
		TaperMultiplier: h.determineTaperMultiplier(state),
	}
}

// determinePhase determines the current phase based on days out from meet.
// This is a simplified phase calculation - programs may override this.
func (h *MeetDateHandler) determinePhase(state *userprogramstate.UserProgramState) string {
	if state.MeetDate == nil {
		return "off_season"
	}

	daysOut := state.DaysOut()
	switch {
	case daysOut <= 7:
		return "meet_week"
	case daysOut <= 14:
		return "peak"
	case daysOut <= 28:
		return "taper"
	case daysOut <= 56:
		return "prep_2"
	case daysOut <= 84:
		return "prep_1"
	default:
		return "base"
	}
}

// determineTaperMultiplier determines the taper multiplier based on phase.
// This is a simplified calculation - the actual taper logic is in loadstrategy.
func (h *MeetDateHandler) determineTaperMultiplier(state *userprogramstate.UserProgramState) float64 {
	if state.MeetDate == nil {
		return 1.0
	}

	daysOut := state.DaysOut()
	switch {
	case daysOut <= 7:
		return 0.4 // Meet week - very light
	case daysOut <= 14:
		return 0.6 // Peak week
	case daysOut <= 21:
		return 0.75 // Taper week 2
	case daysOut <= 28:
		return 0.85 // Taper week 1
	default:
		return 1.0 // Normal training
	}
}

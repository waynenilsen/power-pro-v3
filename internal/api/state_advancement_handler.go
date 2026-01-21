package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/domain/userprogramstate"
	"github.com/waynenilsen/power-pro-v3/internal/middleware"
	"github.com/waynenilsen/power-pro-v3/internal/repository"
)

// StateAdvancementHandler handles HTTP requests for state advancement operations.
type StateAdvancementHandler struct {
	stateRepo *repository.UserProgramStateRepository
	db        *sql.DB
}

// NewStateAdvancementHandler creates a new StateAdvancementHandler.
func NewStateAdvancementHandler(stateRepo *repository.UserProgramStateRepository, db *sql.DB) *StateAdvancementHandler {
	return &StateAdvancementHandler{
		stateRepo: stateRepo,
		db:        db,
	}
}

// StateAdvancementResponse represents the API response for state advancement.
type StateAdvancementResponse struct {
	CurrentWeek           int       `json:"currentWeek"`
	CurrentCycleIteration int       `json:"currentCycleIteration"`
	CurrentDayIndex       *int      `json:"currentDayIndex,omitempty"`
	CycleCompleted        bool      `json:"cycleCompleted"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

// Advance handles POST /users/{userId}/program-state/advance
func (h *StateAdvancementHandler) Advance(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "Missing user ID")
		return
	}

	// Authorization check: only the user themselves or an admin can advance state
	authUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)
	if authUserID != userID && !isAdmin {
		writeError(w, http.StatusForbidden, "Access denied: you can only advance your own state")
		return
	}

	// Get state advancement context
	advCtx, err := h.stateRepo.GetStateAdvancementContext(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get state context")
		return
	}
	if advCtx == nil {
		writeError(w, http.StatusNotFound, "Not enrolled in any program")
		return
	}

	// Handle case where there are no days configured for current week
	daysInWeek := advCtx.DaysInCurrentWeek
	if daysInWeek == 0 {
		// Default to 1 day if no week days configured (allows advancement through weeks)
		daysInWeek = 1
	}

	// Perform state advancement using domain logic
	domainCtx := userprogramstate.AdvancementContext{
		DaysInCurrentWeek: daysInWeek,
		CycleLengthWeeks:  advCtx.CycleLengthWeeks,
	}

	advResult, validation := userprogramstate.AdvanceState(advCtx.State, domainCtx)
	if !validation.Valid {
		details := make([]string, len(validation.Errors))
		for i, err := range validation.Errors {
			details[i] = err.Error()
		}
		writeError(w, http.StatusBadRequest, "State advancement failed", details...)
		return
	}

	// Update state in database (atomic via single UPDATE statement)
	if err := h.stateRepo.Update(advResult.NewState); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update state")
		return
	}

	// Return response
	resp := StateAdvancementResponse{
		CurrentWeek:           advResult.NewState.CurrentWeek,
		CurrentCycleIteration: advResult.NewState.CurrentCycleIteration,
		CurrentDayIndex:       advResult.NewState.CurrentDayIndex,
		CycleCompleted:        advResult.CycleCompleted,
		UpdatedAt:             advResult.NewState.UpdatedAt,
	}

	writeJSON(w, http.StatusOK, resp)
}

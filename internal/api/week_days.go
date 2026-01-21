package api

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/week"
	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
)

// AddDayRequest represents the request body for adding a day to a week.
type AddDayRequest struct {
	DayID     string `json:"dayId"`
	DayOfWeek string `json:"dayOfWeek"`
}

// AddDay handles POST /weeks/{id}/days
func (h *WeekHandler) AddDay(w http.ResponseWriter, r *http.Request) {
	weekID := r.PathValue("id")
	if weekID == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing week ID"))
		return
	}

	// Check week exists
	wk, err := h.repo.GetByID(weekID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get week", err))
		return
	}
	if wk == nil {
		writeDomainError(w, apperrors.NewNotFound("week", weekID))
		return
	}

	var req AddDayRequest
	if err := readJSON(r, &req); err != nil {
		writeDomainError(w, apperrors.NewBadRequest("invalid request body"))
		return
	}

	if strings.TrimSpace(req.DayID) == "" {
		writeDomainError(w, apperrors.NewValidation("dayId", "is required"))
		return
	}

	if strings.TrimSpace(req.DayOfWeek) == "" {
		writeDomainError(w, apperrors.NewValidation("dayOfWeek", "is required"))
		return
	}

	// Check day exists
	dayExists, err := h.repo.DayExists(req.DayID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to verify day", err))
		return
	}
	if !dayExists {
		writeDomainError(w, apperrors.NewValidation("dayId", "day not found"))
		return
	}

	// Generate UUID
	id := uuid.New().String()

	// Create domain entity
	input := week.CreateWeekDayInput{
		WeekID:    weekID,
		DayID:     req.DayID,
		DayOfWeek: req.DayOfWeek,
	}

	newWeekDay, result := week.CreateWeekDay(input, id)
	if !result.Valid {
		details := make([]string, len(result.Errors))
		for i, err := range result.Errors {
			details[i] = err.Error()
		}
		writeDomainError(w, apperrors.NewValidationMsg("validation failed"), details...)
		return
	}

	// Persist
	if err := h.repo.CreateWeekDay(newWeekDay); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to add day to week", err))
		return
	}

	resp := WeekDayResponse{
		ID:        newWeekDay.ID,
		DayID:     newWeekDay.DayID,
		DayOfWeek: string(newWeekDay.DayOfWeek),
		CreatedAt: newWeekDay.CreatedAt,
	}

	writeData(w, http.StatusCreated, resp)
}

// RemoveDay handles DELETE /weeks/{id}/days/{dayId}
func (h *WeekHandler) RemoveDay(w http.ResponseWriter, r *http.Request) {
	weekID := r.PathValue("id")
	dayID := r.PathValue("dayId")

	if weekID == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing week ID"))
		return
	}
	if dayID == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing day ID"))
		return
	}

	// day_of_week is required to identify which mapping to delete
	// (since the same day can be mapped to multiple days of the week)
	dayOfWeek := r.URL.Query().Get("day_of_week")
	if dayOfWeek == "" {
		writeDomainError(w, apperrors.NewValidation("day_of_week", "query parameter is required"))
		return
	}

	// Check week exists
	wk, err := h.repo.GetByID(weekID)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get week", err))
		return
	}
	if wk == nil {
		writeDomainError(w, apperrors.NewNotFound("week", weekID))
		return
	}

	// Check if day mapping exists in this week
	existing, err := h.repo.GetWeekDayByWeekAndDayAndDayOfWeek(weekID, dayID, dayOfWeek)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to check day mapping", err))
		return
	}
	if existing == nil {
		writeDomainError(w, apperrors.NewNotFound("day mapping in week", dayID))
		return
	}

	// Delete
	if err := h.repo.DeleteWeekDay(existing.ID); err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to remove day from week", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

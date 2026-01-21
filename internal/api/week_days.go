package api

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/week"
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
		writeError(w, http.StatusBadRequest, "Missing week ID")
		return
	}

	// Check week exists
	wk, err := h.repo.GetByID(weekID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get week")
		return
	}
	if wk == nil {
		writeError(w, http.StatusNotFound, "Week not found")
		return
	}

	var req AddDayRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.DayID) == "" {
		writeError(w, http.StatusBadRequest, "dayId is required")
		return
	}

	if strings.TrimSpace(req.DayOfWeek) == "" {
		writeError(w, http.StatusBadRequest, "dayOfWeek is required")
		return
	}

	// Check day exists
	dayExists, err := h.repo.DayExists(req.DayID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to verify day")
		return
	}
	if !dayExists {
		writeError(w, http.StatusBadRequest, "Day not found")
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
		writeError(w, http.StatusBadRequest, "Validation failed", details...)
		return
	}

	// Persist
	if err := h.repo.CreateWeekDay(newWeekDay); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to add day to week")
		return
	}

	resp := WeekDayResponse{
		ID:        newWeekDay.ID,
		DayID:     newWeekDay.DayID,
		DayOfWeek: string(newWeekDay.DayOfWeek),
		CreatedAt: newWeekDay.CreatedAt,
	}

	writeJSON(w, http.StatusCreated, resp)
}

// RemoveDay handles DELETE /weeks/{id}/days/{dayId}
func (h *WeekHandler) RemoveDay(w http.ResponseWriter, r *http.Request) {
	weekID := r.PathValue("id")
	dayID := r.PathValue("dayId")

	if weekID == "" {
		writeError(w, http.StatusBadRequest, "Missing week ID")
		return
	}
	if dayID == "" {
		writeError(w, http.StatusBadRequest, "Missing day ID")
		return
	}

	// day_of_week is required to identify which mapping to delete
	// (since the same day can be mapped to multiple days of the week)
	dayOfWeek := r.URL.Query().Get("day_of_week")
	if dayOfWeek == "" {
		writeError(w, http.StatusBadRequest, "day_of_week query parameter is required")
		return
	}

	// Check week exists
	wk, err := h.repo.GetByID(weekID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get week")
		return
	}
	if wk == nil {
		writeError(w, http.StatusNotFound, "Week not found")
		return
	}

	// Check if day mapping exists in this week
	existing, err := h.repo.GetWeekDayByWeekAndDayAndDayOfWeek(weekID, dayID, dayOfWeek)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to check day mapping")
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "Day mapping not found in this week")
		return
	}

	// Delete
	if err := h.repo.DeleteWeekDay(existing.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to remove day from week")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

package api

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/day"
)

// AddPrescriptionRequest represents the request body for adding a prescription to a day.
type AddPrescriptionRequest struct {
	PrescriptionID string `json:"prescriptionId"`
	Order          *int   `json:"order,omitempty"`
}

// ReorderPrescriptionsRequest represents the request body for reordering prescriptions in a day.
type ReorderPrescriptionsRequest struct {
	PrescriptionIDs []string `json:"prescriptionIds"`
}

// AddPrescription handles POST /days/{id}/prescriptions
func (h *DayHandler) AddPrescription(w http.ResponseWriter, r *http.Request) {
	dayID := r.PathValue("id")
	if dayID == "" {
		writeError(w, http.StatusBadRequest, "Missing day ID")
		return
	}

	// Check day exists
	d, err := h.repo.GetByID(dayID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get day")
		return
	}
	if d == nil {
		writeError(w, http.StatusNotFound, "Day not found")
		return
	}

	var req AddPrescriptionRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.PrescriptionID) == "" {
		writeError(w, http.StatusBadRequest, "prescriptionId is required")
		return
	}

	// Check prescription exists
	prescription, err := h.prescriptionRepo.GetByID(req.PrescriptionID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to verify prescription")
		return
	}
	if prescription == nil {
		writeError(w, http.StatusBadRequest, "Prescription not found")
		return
	}

	// Check if this prescription is already in this day
	existing, err := h.repo.GetDayPrescriptionByDayAndPrescription(dayID, req.PrescriptionID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to check existing prescription")
		return
	}
	if existing != nil {
		writeError(w, http.StatusConflict, "Prescription is already in this day")
		return
	}

	// Determine order
	order := 0
	if req.Order != nil {
		if *req.Order < 0 {
			writeError(w, http.StatusBadRequest, "Order must be >= 0")
			return
		}
		order = *req.Order
	} else {
		// Auto-assign next order
		maxOrder, err := h.repo.GetMaxDayPrescriptionOrder(dayID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to determine order")
			return
		}
		order = maxOrder + 1
	}

	// Generate UUID
	id := uuid.New().String()

	// Create domain entity
	input := day.CreateDayPrescriptionInput{
		DayID:          dayID,
		PrescriptionID: req.PrescriptionID,
		Order:          &order,
	}

	newDayPrescription, result := day.CreateDayPrescription(input, id, order)
	if !result.Valid {
		details := make([]string, len(result.Errors))
		for i, err := range result.Errors {
			details[i] = err.Error()
		}
		writeError(w, http.StatusBadRequest, "Validation failed", details...)
		return
	}

	// Persist
	if err := h.repo.CreateDayPrescription(newDayPrescription); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to add prescription to day")
		return
	}

	resp := DayPrescriptionResponse{
		ID:             newDayPrescription.ID,
		PrescriptionID: newDayPrescription.PrescriptionID,
		Order:          newDayPrescription.Order,
		CreatedAt:      newDayPrescription.CreatedAt,
	}

	writeJSON(w, http.StatusCreated, resp)
}

// RemovePrescription handles DELETE /days/{id}/prescriptions/{prescriptionId}
func (h *DayHandler) RemovePrescription(w http.ResponseWriter, r *http.Request) {
	dayID := r.PathValue("id")
	prescriptionID := r.PathValue("prescriptionId")

	if dayID == "" {
		writeError(w, http.StatusBadRequest, "Missing day ID")
		return
	}
	if prescriptionID == "" {
		writeError(w, http.StatusBadRequest, "Missing prescription ID")
		return
	}

	// Check day exists
	d, err := h.repo.GetByID(dayID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get day")
		return
	}
	if d == nil {
		writeError(w, http.StatusNotFound, "Day not found")
		return
	}

	// Check if prescription is in this day
	existing, err := h.repo.GetDayPrescriptionByDayAndPrescription(dayID, prescriptionID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to check prescription")
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "Prescription not found in this day")
		return
	}

	// Delete
	if err := h.repo.DeleteDayPrescription(existing.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to remove prescription from day")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ReorderPrescriptions handles PUT /days/{id}/prescriptions/reorder
func (h *DayHandler) ReorderPrescriptions(w http.ResponseWriter, r *http.Request) {
	dayID := r.PathValue("id")
	if dayID == "" {
		writeError(w, http.StatusBadRequest, "Missing day ID")
		return
	}

	// Check day exists
	d, err := h.repo.GetByID(dayID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get day")
		return
	}
	if d == nil {
		writeError(w, http.StatusNotFound, "Day not found")
		return
	}

	var req ReorderPrescriptionsRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	input := day.ReorderPrescriptionsInput{
		DayID:           dayID,
		PrescriptionIDs: req.PrescriptionIDs,
	}
	result := day.ValidateReorderInput(input)
	if !result.Valid {
		details := make([]string, len(result.Errors))
		for i, err := range result.Errors {
			details[i] = err.Error()
		}
		writeError(w, http.StatusBadRequest, "Validation failed", details...)
		return
	}

	// Get current prescriptions for this day
	currentPrescriptions, err := h.repo.ListDayPrescriptions(dayID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get current prescriptions")
		return
	}

	// Create a map of prescription ID to day prescription
	prescriptionMap := make(map[string]*day.DayPrescription)
	for i := range currentPrescriptions {
		prescriptionMap[currentPrescriptions[i].PrescriptionID] = &currentPrescriptions[i]
	}

	// Verify all prescription IDs exist in this day
	for _, prescriptionID := range req.PrescriptionIDs {
		if _, ok := prescriptionMap[prescriptionID]; !ok {
			writeError(w, http.StatusBadRequest, "Prescription not found in this day: "+prescriptionID)
			return
		}
	}

	// Verify count matches
	if len(req.PrescriptionIDs) != len(currentPrescriptions) {
		writeError(w, http.StatusBadRequest, "Prescription IDs count does not match current prescriptions count")
		return
	}

	// Update orders
	for newOrder, prescriptionID := range req.PrescriptionIDs {
		dayPrescription := prescriptionMap[prescriptionID]
		if dayPrescription.Order != newOrder {
			if err := h.repo.UpdateDayPrescriptionOrder(dayPrescription.ID, newOrder); err != nil {
				writeError(w, http.StatusInternalServerError, "Failed to update prescription order")
				return
			}
		}
	}

	// Return updated day with prescriptions
	updatedPrescriptions, err := h.repo.ListDayPrescriptions(dayID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get updated prescriptions")
		return
	}

	writeJSON(w, http.StatusOK, dayToResponseWithPrescriptions(d, updatedPrescriptions))
}

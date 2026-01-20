package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/day"
	"github.com/waynenilsen/power-pro-v3/internal/repository"
)

// DayHandler handles HTTP requests for day operations.
type DayHandler struct {
	repo             *repository.DayRepository
	prescriptionRepo *repository.PrescriptionRepository
}

// NewDayHandler creates a new DayHandler.
func NewDayHandler(repo *repository.DayRepository, prescriptionRepo *repository.PrescriptionRepository) *DayHandler {
	return &DayHandler{
		repo:             repo,
		prescriptionRepo: prescriptionRepo,
	}
}

// DayResponse represents the API response format for a day.
type DayResponse struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Slug      string                 `json:"slug"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	ProgramID *string                `json:"programId,omitempty"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
}

// DayWithPrescriptionsResponse represents the API response format for a day with its prescriptions.
type DayWithPrescriptionsResponse struct {
	ID            string                           `json:"id"`
	Name          string                           `json:"name"`
	Slug          string                           `json:"slug"`
	Metadata      map[string]interface{}           `json:"metadata,omitempty"`
	ProgramID     *string                          `json:"programId,omitempty"`
	Prescriptions []DayPrescriptionResponse        `json:"prescriptions"`
	CreatedAt     time.Time                        `json:"createdAt"`
	UpdatedAt     time.Time                        `json:"updatedAt"`
}

// DayPrescriptionResponse represents a prescription within a day.
type DayPrescriptionResponse struct {
	ID             string    `json:"id"`
	PrescriptionID string    `json:"prescriptionId"`
	Order          int       `json:"order"`
	CreatedAt      time.Time `json:"createdAt"`
}

// CreateDayRequest represents the request body for creating a day.
type CreateDayRequest struct {
	Name      string                 `json:"name"`
	Slug      string                 `json:"slug,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	ProgramID *string                `json:"programId,omitempty"`
}

// UpdateDayRequest represents the request body for updating a day.
type UpdateDayRequest struct {
	Name           *string                `json:"name,omitempty"`
	Slug           *string                `json:"slug,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	ClearMetadata  bool                   `json:"clearMetadata,omitempty"`
	ProgramID      *string                `json:"programId,omitempty"`
	ClearProgramID bool                   `json:"clearProgramId,omitempty"`
}

// AddPrescriptionRequest represents the request body for adding a prescription to a day.
type AddPrescriptionRequest struct {
	PrescriptionID string `json:"prescriptionId"`
	Order          *int   `json:"order,omitempty"`
}

// ReorderPrescriptionsRequest represents the request body for reordering prescriptions in a day.
type ReorderPrescriptionsRequest struct {
	PrescriptionIDs []string `json:"prescriptionIds"`
}

func dayToResponse(d *day.Day) DayResponse {
	return DayResponse{
		ID:        d.ID,
		Name:      d.Name,
		Slug:      d.Slug,
		Metadata:  d.Metadata,
		ProgramID: d.ProgramID,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func dayToResponseWithPrescriptions(d *day.Day, prescriptions []day.DayPrescription) DayWithPrescriptionsResponse {
	prescriptionResponses := make([]DayPrescriptionResponse, len(prescriptions))
	for i, p := range prescriptions {
		prescriptionResponses[i] = DayPrescriptionResponse{
			ID:             p.ID,
			PrescriptionID: p.PrescriptionID,
			Order:          p.Order,
			CreatedAt:      p.CreatedAt,
		}
	}

	return DayWithPrescriptionsResponse{
		ID:            d.ID,
		Name:          d.Name,
		Slug:          d.Slug,
		Metadata:      d.Metadata,
		ProgramID:     d.ProgramID,
		Prescriptions: prescriptionResponses,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
}

// List handles GET /days
func (h *DayHandler) List(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()

	// Pagination
	page := 1
	pageSize := 20
	if p := query.Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if ps := query.Get("pageSize"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	// Sorting
	sortBy := repository.DaySortByName
	sortOrder := repository.SortAsc
	if s := query.Get("sortBy"); s != "" {
		switch strings.ToLower(s) {
		case "name":
			sortBy = repository.DaySortByName
		case "created_at", "createdat":
			sortBy = repository.DaySortByCreatedAt
		}
	}
	if o := query.Get("sortOrder"); o != "" {
		switch strings.ToLower(o) {
		case "asc":
			sortOrder = repository.SortAsc
		case "desc":
			sortOrder = repository.SortDesc
		}
	}

	// Filter by program_id
	var filterProgramID *string
	if programID := query.Get("program_id"); programID != "" {
		filterProgramID = &programID
	}

	params := repository.DayListParams{
		Limit:           int64(pageSize),
		Offset:          int64((page - 1) * pageSize),
		SortBy:          sortBy,
		SortOrder:       sortOrder,
		FilterProgramID: filterProgramID,
	}

	days, total, err := h.repo.List(params)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list days")
		return
	}

	// Convert to response format
	data := make([]DayResponse, len(days))
	for i, d := range days {
		data[i] = dayToResponse(&d)
	}

	// Calculate total pages
	totalPages := total / int64(pageSize)
	if total%int64(pageSize) > 0 {
		totalPages++
	}

	resp := PaginatedResponse{
		Data:       data,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}

	writeJSON(w, http.StatusOK, resp)
}

// Get handles GET /days/{id}
func (h *DayHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing day ID")
		return
	}

	d, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get day")
		return
	}
	if d == nil {
		writeError(w, http.StatusNotFound, "Day not found")
		return
	}

	// Get prescriptions for this day
	prescriptions, err := h.repo.ListDayPrescriptions(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get day prescriptions")
		return
	}

	writeJSON(w, http.StatusOK, dayToResponseWithPrescriptions(d, prescriptions))
}

// GetBySlug handles GET /days/by-slug/{slug}
func (h *DayHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		writeError(w, http.StatusBadRequest, "Missing slug")
		return
	}

	// Optional program_id query parameter
	var programID *string
	if pid := r.URL.Query().Get("program_id"); pid != "" {
		programID = &pid
	}

	d, err := h.repo.GetBySlug(slug, programID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get day")
		return
	}
	if d == nil {
		writeError(w, http.StatusNotFound, "Day not found")
		return
	}

	// Get prescriptions for this day
	prescriptions, err := h.repo.ListDayPrescriptions(d.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get day prescriptions")
		return
	}

	writeJSON(w, http.StatusOK, dayToResponseWithPrescriptions(d, prescriptions))
}

// Create handles POST /days
func (h *DayHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateDayRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate metadata JSON if provided
	if req.Metadata != nil {
		if _, err := json.Marshal(req.Metadata); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid metadata JSON")
			return
		}
	}

	// Generate UUID
	id := uuid.New().String()

	// Use domain logic to create and validate
	input := day.CreateDayInput{
		Name:      req.Name,
		Slug:      req.Slug,
		Metadata:  req.Metadata,
		ProgramID: req.ProgramID,
	}

	newDay, result := day.CreateDay(input, id)
	if !result.Valid {
		details := make([]string, len(result.Errors))
		for i, err := range result.Errors {
			details[i] = err.Error()
		}
		writeError(w, http.StatusBadRequest, "Validation failed", details...)
		return
	}

	// Check for slug conflict within the program
	exists, err := h.repo.SlugExists(newDay.Slug, newDay.ProgramID, nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to check slug uniqueness")
		return
	}
	if exists {
		writeError(w, http.StatusConflict, "Slug already exists within this program")
		return
	}

	// Persist
	if err := h.repo.Create(newDay); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create day")
		return
	}

	writeJSON(w, http.StatusCreated, dayToResponse(newDay))
}

// Update handles PUT /days/{id}
func (h *DayHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing day ID")
		return
	}

	// Get existing day
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get day")
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "Day not found")
		return
	}

	var req UpdateDayRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate metadata JSON if provided
	if req.Metadata != nil {
		if _, err := json.Marshal(req.Metadata); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid metadata JSON")
			return
		}
	}

	// Determine the new program ID for slug uniqueness check
	newProgramID := existing.ProgramID
	if req.ClearProgramID {
		newProgramID = nil
	} else if req.ProgramID != nil {
		newProgramID = req.ProgramID
	}

	// Check for slug conflict if slug or program is being changed
	newSlug := existing.Slug
	if req.Slug != nil {
		newSlug = *req.Slug
	}
	if newSlug != existing.Slug || !ptrStringEqual(newProgramID, existing.ProgramID) {
		exists, err := h.repo.SlugExists(newSlug, newProgramID, &id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to check slug uniqueness")
			return
		}
		if exists {
			writeError(w, http.StatusConflict, "Slug already exists within this program")
			return
		}
	}

	// Use domain logic to update and validate
	input := day.UpdateDayInput{
		Name:           req.Name,
		Slug:           req.Slug,
		Metadata:       req.Metadata,
		ClearMetadata:  req.ClearMetadata,
		ProgramID:      req.ProgramID,
		ClearProgramID: req.ClearProgramID,
	}

	result := day.UpdateDay(existing, input)
	if !result.Valid {
		details := make([]string, len(result.Errors))
		for i, err := range result.Errors {
			details[i] = err.Error()
		}
		writeError(w, http.StatusBadRequest, "Validation failed", details...)
		return
	}

	// Persist
	if err := h.repo.Update(existing); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update day")
		return
	}

	writeJSON(w, http.StatusOK, dayToResponse(existing))
}

// Delete handles DELETE /days/{id}
func (h *DayHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing day ID")
		return
	}

	// Check day exists
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get day")
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "Day not found")
		return
	}

	// Check if day is used in any weeks
	isUsed, err := h.repo.IsUsedInWeeks(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to check if day is used")
		return
	}
	if isUsed {
		writeError(w, http.StatusConflict, "Cannot delete day: it is used in one or more weeks")
		return
	}

	// Delete
	if err := h.repo.Delete(id); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete day")
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

// Helper function to compare string pointers
func ptrStringEqual(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

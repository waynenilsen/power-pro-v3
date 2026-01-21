package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/liftmax"
	"github.com/waynenilsen/power-pro-v3/internal/middleware"
)

// ConversionResponse represents the API response format for a max conversion.
type ConversionResponse struct {
	OriginalValue  float64 `json:"originalValue"`
	OriginalType   string  `json:"originalType"`
	ConvertedValue float64 `json:"convertedValue"`
	ConvertedType  string  `json:"convertedType"`
	Percentage     float64 `json:"percentage"`
}

// GetCurrent handles GET /users/{userId}/lift-maxes/current
// Returns the most recent lift max for a user, lift, and type combination.
func (h *LiftMaxHandler) GetCurrent(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "Missing user ID")
		return
	}

	// Authorization is handled by middleware for the userId in path
	// The middleware already verified the requesting user is owner or admin

	// Parse and validate query parameters
	query := r.URL.Query()

	// lift is required
	liftID := query.Get("lift")
	if liftID == "" {
		writeError(w, http.StatusBadRequest, "Missing required query parameter: lift")
		return
	}

	// Validate lift is a valid UUID
	if _, err := uuid.Parse(liftID); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid lift parameter: must be a valid UUID")
		return
	}

	// type is required
	maxType := strings.ToUpper(query.Get("type"))
	if maxType == "" {
		writeError(w, http.StatusBadRequest, "Missing required query parameter: type")
		return
	}

	// Validate type
	if maxType != string(liftmax.OneRM) && maxType != string(liftmax.TrainingMax) {
		writeError(w, http.StatusBadRequest, "Invalid type parameter: must be ONE_RM or TRAINING_MAX")
		return
	}

	// Query for the current max
	m, err := h.repo.GetCurrentMax(userID, liftID, maxType)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get current lift max")
		return
	}
	if m == nil {
		writeError(w, http.StatusNotFound, "No lift max found for the specified user, lift, and type")
		return
	}

	writeJSON(w, http.StatusOK, liftMaxToResponse(m))
}

// Convert handles GET /lift-maxes/{id}/convert
// Converts a lift max between 1RM and Training Max without persisting.
func (h *LiftMaxHandler) Convert(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing lift max ID")
		return
	}

	// Get the lift max
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get lift max")
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "Lift max not found")
		return
	}

	// Check ownership: user must be owner or admin
	requestingUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)
	if requestingUserID != existing.UserID && !isAdmin {
		writeError(w, http.StatusForbidden, "Access denied: you do not have permission to access this resource")
		return
	}

	// Parse query parameters
	query := r.URL.Query()

	// to_type is required
	toType := strings.ToUpper(query.Get("to_type"))
	if toType == "" {
		writeError(w, http.StatusBadRequest, "Missing required query parameter: to_type")
		return
	}

	// Validate to_type
	if toType != string(liftmax.OneRM) && toType != string(liftmax.TrainingMax) {
		writeError(w, http.StatusBadRequest, "Invalid to_type: must be ONE_RM or TRAINING_MAX")
		return
	}

	// Check if to_type is same as current type
	if toType == string(existing.Type) {
		writeError(w, http.StatusBadRequest, "Cannot convert to same type: lift max is already "+toType)
		return
	}

	// Parse percentage (optional, default 90)
	percentage := liftmax.DefaultTMPercentage
	if pctStr := query.Get("percentage"); pctStr != "" {
		pct, err := strconv.ParseFloat(pctStr, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid percentage: must be a number")
			return
		}
		if pct < 1 || pct > 100 {
			writeError(w, http.StatusBadRequest, "Invalid percentage: must be between 1 and 100")
			return
		}
		percentage = pct
	}

	// Perform conversion using domain logic
	calculator := liftmax.NewMaxCalculator()
	var convertedValue float64

	if existing.Type == liftmax.OneRM {
		// Converting from 1RM to TM
		convertedValue, err = calculator.ConvertToTM(existing.Value, &percentage)
	} else {
		// Converting from TM to 1RM
		convertedValue, err = calculator.ConvertToOneRM(existing.Value, &percentage)
	}

	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := ConversionResponse{
		OriginalValue:  existing.Value,
		OriginalType:   string(existing.Type),
		ConvertedValue: convertedValue,
		ConvertedType:  toType,
		Percentage:     percentage,
	}

	writeJSON(w, http.StatusOK, response)
}

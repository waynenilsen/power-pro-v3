package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/liftmax"
	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
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
		writeDomainError(w, apperrors.NewBadRequest("missing user ID"))
		return
	}

	// Authorization is handled by middleware for the userId in path
	// The middleware already verified the requesting user is owner or admin

	// Parse and validate query parameters
	query := r.URL.Query()

	// lift is required
	liftID := query.Get("lift")
	if liftID == "" {
		writeDomainError(w, apperrors.NewValidation("lift", "missing required query parameter"))
		return
	}

	// Validate lift is a valid UUID
	if _, err := uuid.Parse(liftID); err != nil {
		writeDomainError(w, apperrors.NewValidation("lift", "must be a valid UUID"))
		return
	}

	// type is required
	maxType := strings.ToUpper(query.Get("type"))
	if maxType == "" {
		writeDomainError(w, apperrors.NewValidation("type", "missing required query parameter"))
		return
	}

	// Validate type
	if maxType != string(liftmax.OneRM) && maxType != string(liftmax.TrainingMax) {
		writeDomainError(w, apperrors.NewValidation("type", "must be ONE_RM or TRAINING_MAX"))
		return
	}

	// Query for the current max
	m, err := h.repo.GetCurrentMax(userID, liftID, maxType)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get current lift max", err))
		return
	}
	if m == nil {
		writeDomainError(w, apperrors.NewNotFound("lift max", "user/lift/type combination"))
		return
	}

	writeData(w, http.StatusOK, liftMaxToResponse(m))
}

// Convert handles GET /lift-maxes/{id}/convert
// Converts a lift max between 1RM and Training Max without persisting.
func (h *LiftMaxHandler) Convert(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing lift max ID"))
		return
	}

	// Get the lift max
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to get lift max", err))
		return
	}
	if existing == nil {
		writeDomainError(w, apperrors.NewNotFound("lift max", id))
		return
	}

	// Check ownership: user must be owner or admin
	requestingUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)
	if requestingUserID != existing.UserID && !isAdmin {
		writeDomainError(w, apperrors.NewForbidden("you do not have permission to access this resource"))
		return
	}

	// Parse query parameters
	query := r.URL.Query()

	// to_type is required
	toType := strings.ToUpper(query.Get("to_type"))
	if toType == "" {
		writeDomainError(w, apperrors.NewValidation("to_type", "missing required query parameter"))
		return
	}

	// Validate to_type
	if toType != string(liftmax.OneRM) && toType != string(liftmax.TrainingMax) {
		writeDomainError(w, apperrors.NewValidation("to_type", "must be ONE_RM or TRAINING_MAX"))
		return
	}

	// Check if to_type is same as current type
	if toType == string(existing.Type) {
		writeDomainError(w, apperrors.NewBadRequest("cannot convert to same type: lift max is already "+toType))
		return
	}

	// Parse percentage (optional, default 90)
	percentage := liftmax.DefaultTMPercentage
	if pctStr := query.Get("percentage"); pctStr != "" {
		pct, err := strconv.ParseFloat(pctStr, 64)
		if err != nil {
			writeDomainError(w, apperrors.NewValidation("percentage", "must be a number"))
			return
		}
		if pct < 1 || pct > 100 {
			writeDomainError(w, apperrors.NewValidation("percentage", "must be between 1 and 100"))
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
		writeDomainError(w, apperrors.NewBadRequest(err.Error()))
		return
	}

	response := ConversionResponse{
		OriginalValue:  existing.Value,
		OriginalType:   string(existing.Type),
		ConvertedValue: convertedValue,
		ConvertedType:  toType,
		Percentage:     percentage,
	}

	writeData(w, http.StatusOK, response)
}

package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/domain/progression"
	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
	"github.com/waynenilsen/power-pro-v3/internal/middleware"
	"github.com/waynenilsen/power-pro-v3/internal/repository"
)

// ProgressionHistoryHandler handles HTTP requests for progression history queries.
type ProgressionHistoryHandler struct {
	repo *repository.ProgressionHistoryRepository
}

// NewProgressionHistoryHandler creates a new ProgressionHistoryHandler.
func NewProgressionHistoryHandler(repo *repository.ProgressionHistoryRepository) *ProgressionHistoryHandler {
	return &ProgressionHistoryHandler{repo: repo}
}

// ProgressionHistoryResponse represents the API response format for a progression history entry.
type ProgressionHistoryResponse struct {
	ID              string          `json:"id"`
	ProgressionID   string          `json:"progressionId"`
	ProgressionName string          `json:"progressionName"`
	ProgressionType string          `json:"progressionType"`
	LiftID          string          `json:"liftId"`
	LiftName        string          `json:"liftName"`
	PreviousValue   float64         `json:"previousValue"`
	NewValue        float64         `json:"newValue"`
	Delta           float64         `json:"delta"`
	TriggerType     string          `json:"triggerType"`
	TriggerContext  json.RawMessage `json:"triggerContext"`
	AppliedAt       time.Time       `json:"appliedAt"`
}

// ProgressionHistoryListResponse wraps the paginated progression history response.
type ProgressionHistoryListResponse struct {
	Data       []ProgressionHistoryResponse `json:"data"`
	Limit      int                          `json:"limit"`
	Offset     int                          `json:"offset"`
	TotalItems int64                        `json:"totalItems"`
}

// List handles GET /users/{userId}/progression-history
func (h *ProgressionHistoryHandler) List(w http.ResponseWriter, r *http.Request) {
	// Get userId from path
	userID := r.PathValue("userId")
	if userID == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing user ID"))
		return
	}

	// Authorization check: user can only query their own history, admin can query any
	authUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)
	if !isAdmin && authUserID != userID {
		writeDomainError(w, apperrors.NewForbidden("you do not have permission to access this resource"))
		return
	}

	// Parse query parameters
	query := r.URL.Query()

	// Pagination
	limit := 20
	offset := 0
	if l := query.Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
			if limit > 100 {
				limit = 100
			}
		}
	}
	if o := query.Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	// Build filter
	filter := repository.ProgressionHistoryFilter{
		UserID: userID,
		Limit:  int64(limit),
		Offset: int64(offset),
	}

	// Filter by liftId
	if liftID := query.Get("liftId"); liftID != "" {
		filter.LiftID = &liftID
	}

	// Filter by progressionType
	if pt := query.Get("progressionType"); pt != "" {
		// Normalize to uppercase
		normalizedPT := strings.ToUpper(pt)
		// Validate progression type
		if progression.ValidProgressionTypes[progression.ProgressionType(normalizedPT)] {
			filter.ProgressionType = &normalizedPT
		} else {
			writeDomainError(w, apperrors.NewValidation("progressionType", "invalid value; valid values: LINEAR_PROGRESSION, CYCLE_PROGRESSION"))
			return
		}
	}

	// Filter by triggerType
	if tt := query.Get("triggerType"); tt != "" {
		// Normalize to uppercase
		normalizedTT := strings.ToUpper(tt)
		// Validate trigger type
		if progression.ValidTriggerTypes[progression.TriggerType(normalizedTT)] {
			filter.TriggerType = &normalizedTT
		} else {
			writeDomainError(w, apperrors.NewValidation("triggerType", "invalid value; valid values: AFTER_SESSION, AFTER_WEEK, AFTER_CYCLE"))
			return
		}
	}

	// Filter by startDate (ISO 8601)
	if sd := query.Get("startDate"); sd != "" {
		// Try parsing as RFC3339 first, then as date-only
		if t, err := time.Parse(time.RFC3339, sd); err == nil {
			filter.StartDate = &t
		} else if t, err := time.Parse("2006-01-02", sd); err == nil {
			filter.StartDate = &t
		} else {
			writeDomainError(w, apperrors.NewValidation("startDate", "invalid format; use ISO 8601 format (e.g., 2024-01-15 or 2024-01-15T10:00:00Z)"))
			return
		}
	}

	// Filter by endDate (ISO 8601)
	if ed := query.Get("endDate"); ed != "" {
		// Try parsing as RFC3339 first, then as date-only
		if t, err := time.Parse(time.RFC3339, ed); err == nil {
			filter.EndDate = &t
		} else if t, err := time.Parse("2006-01-02", ed); err == nil {
			// Set to end of day for date-only format
			t = t.Add(24*time.Hour - time.Second)
			filter.EndDate = &t
		} else {
			writeDomainError(w, apperrors.NewValidation("endDate", "invalid format; use ISO 8601 format (e.g., 2024-01-15 or 2024-01-15T10:00:00Z)"))
			return
		}
	}

	// Fetch data
	entries, total, err := h.repo.List(r.Context(), filter)
	if err != nil {
		writeDomainError(w, apperrors.NewInternal("failed to fetch progression history", err))
		return
	}

	// Convert to response format
	data := make([]ProgressionHistoryResponse, len(entries))
	for i, entry := range entries {
		data[i] = ProgressionHistoryResponse{
			ID:              entry.ID,
			ProgressionID:   entry.ProgressionID,
			ProgressionName: entry.ProgressionName,
			ProgressionType: entry.ProgressionType,
			LiftID:          entry.LiftID,
			LiftName:        entry.LiftName,
			PreviousValue:   entry.PreviousValue,
			NewValue:        entry.NewValue,
			Delta:           entry.Delta,
			TriggerType:     entry.TriggerType,
			TriggerContext:  entry.TriggerContext,
			AppliedAt:       entry.AppliedAt,
		}
	}

	// Use standard envelope with pagination metadata
	writePaginatedData(w, http.StatusOK, data, total, limit, offset)
}

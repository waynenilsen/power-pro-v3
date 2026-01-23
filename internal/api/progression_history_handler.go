package api

import (
	"encoding/json"
	"net/http"
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

	// Pagination (limit/offset)
	pg := ParsePagination(query)

	// Build filter
	filter := repository.ProgressionHistoryFilter{
		UserID: userID,
		Limit:  int64(pg.Limit),
		Offset: int64(pg.Offset),
	}

	// Filter by liftId
	filter.LiftID = ParseFilterString(query, "liftId")

	// Filter by progressionType (enum validation)
	progressionTypes := make([]string, 0, len(progression.ValidProgressionTypes))
	for pt := range progression.ValidProgressionTypes {
		progressionTypes = append(progressionTypes, string(pt))
	}
	progressionType, err := ParseFilterEnum(query, "progressionType", progressionTypes)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	filter.ProgressionType = progressionType

	// Filter by triggerType (enum validation)
	triggerTypes := make([]string, 0, len(progression.ValidTriggerTypes))
	for tt := range progression.ValidTriggerTypes {
		triggerTypes = append(triggerTypes, string(tt))
	}
	triggerType, err := ParseFilterEnum(query, "triggerType", triggerTypes)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	filter.TriggerType = triggerType

	// Filter by startDate (ISO 8601)
	startDate, err := ParseFilterDate(query, "startDate")
	if err != nil {
		writeDomainError(w, err)
		return
	}
	filter.StartDate = startDate

	// Filter by endDate (ISO 8601, end of day for date-only format)
	endDate, err := ParseFilterDateEndOfDay(query, "endDate")
	if err != nil {
		writeDomainError(w, err)
		return
	}
	filter.EndDate = endDate

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
	writePaginatedData(w, http.StatusOK, data, total, pg.Limit, pg.Offset)
}

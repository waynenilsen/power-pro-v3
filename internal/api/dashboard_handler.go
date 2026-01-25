package api

import (
	"net/http"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/dashboard"
	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
	"github.com/waynenilsen/power-pro-v3/internal/middleware"
)

// DashboardHandler handles HTTP requests for dashboard operations.
type DashboardHandler struct {
	dashboardService *dashboard.Service
}

// NewDashboardHandler creates a new DashboardHandler.
func NewDashboardHandler(dashboardService *dashboard.Service) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// DashboardResponse represents the response for GET /users/{id}/dashboard.
type DashboardResponse struct {
	Enrollment     *EnrollmentSummaryResponse   `json:"enrollment"`
	NextWorkout    *NextWorkoutPreviewResponse  `json:"nextWorkout"`
	CurrentSession *SessionSummaryResponse      `json:"currentSession"`
	RecentWorkouts []WorkoutSummaryResponse     `json:"recentWorkouts"`
	CurrentMaxes   []MaxSummaryResponse         `json:"currentMaxes"`
}

// EnrollmentSummaryResponse represents the enrollment section.
type EnrollmentSummaryResponse struct {
	Status         string `json:"status"`
	ProgramName    string `json:"programName"`
	CycleIteration int    `json:"cycleIteration"`
	CycleStatus    string `json:"cycleStatus"`
	WeekNumber     int    `json:"weekNumber"`
	WeekStatus     string `json:"weekStatus"`
}

// NextWorkoutPreviewResponse represents the next workout section.
type NextWorkoutPreviewResponse struct {
	DayName       string `json:"dayName"`
	DaySlug       string `json:"daySlug"`
	ExerciseCount int    `json:"exerciseCount"`
	EstimatedSets int    `json:"estimatedSets"`
}

// SessionSummaryResponse represents the current session section.
type SessionSummaryResponse struct {
	SessionID     string    `json:"sessionId"`
	DayName       string    `json:"dayName"`
	StartedAt     time.Time `json:"startedAt"`
	SetsCompleted int       `json:"setsCompleted"`
	TotalSets     int       `json:"totalSets"`
}

// WorkoutSummaryResponse represents a recent workout.
type WorkoutSummaryResponse struct {
	Date          string `json:"date"`
	DayName       string `json:"dayName"`
	SetsCompleted int    `json:"setsCompleted"`
}

// MaxSummaryResponse represents a current max.
type MaxSummaryResponse struct {
	Lift  string  `json:"lift"`
	Value float64 `json:"value"`
	Type  string  `json:"type"`
}

// Get handles GET /users/{id}/dashboard
func (h *DashboardHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	if userID == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing user ID"))
		return
	}

	// Authorization check: owner-only (not even admins can view other users' dashboards)
	authUserID := middleware.GetUserID(r)
	if authUserID != userID {
		writeDomainError(w, apperrors.NewForbidden("dashboard access is owner-only"))
		return
	}

	// Get dashboard from service
	dash, err := h.dashboardService.GetDashboard(r.Context(), userID)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	// Build and return response
	response := h.buildDashboardResponse(dash)
	writeData(w, http.StatusOK, response)
}

// buildDashboardResponse converts the domain dashboard to a response.
func (h *DashboardHandler) buildDashboardResponse(dash *dashboard.Dashboard) DashboardResponse {
	response := DashboardResponse{
		RecentWorkouts: make([]WorkoutSummaryResponse, 0),
		CurrentMaxes:   make([]MaxSummaryResponse, 0),
	}

	// Enrollment section
	if dash.Enrollment != nil {
		response.Enrollment = &EnrollmentSummaryResponse{
			Status:         dash.Enrollment.Status,
			ProgramName:    dash.Enrollment.ProgramName,
			CycleIteration: dash.Enrollment.CycleIteration,
			CycleStatus:    dash.Enrollment.CycleStatus,
			WeekNumber:     dash.Enrollment.WeekNumber,
			WeekStatus:     dash.Enrollment.WeekStatus,
		}
	}

	// Next workout section
	if dash.NextWorkout != nil {
		response.NextWorkout = &NextWorkoutPreviewResponse{
			DayName:       dash.NextWorkout.DayName,
			DaySlug:       dash.NextWorkout.DaySlug,
			ExerciseCount: dash.NextWorkout.ExerciseCount,
			EstimatedSets: dash.NextWorkout.EstimatedSets,
		}
	}

	// Current session section
	if dash.CurrentSession != nil {
		response.CurrentSession = &SessionSummaryResponse{
			SessionID:     dash.CurrentSession.SessionID,
			DayName:       dash.CurrentSession.DayName,
			StartedAt:     dash.CurrentSession.StartedAt,
			SetsCompleted: dash.CurrentSession.SetsCompleted,
			TotalSets:     dash.CurrentSession.TotalSets,
		}
	}

	// Recent workouts
	for _, ws := range dash.RecentWorkouts {
		response.RecentWorkouts = append(response.RecentWorkouts, WorkoutSummaryResponse{
			Date:          ws.Date,
			DayName:       ws.DayName,
			SetsCompleted: ws.SetsCompleted,
		})
	}

	// Current maxes
	for _, ms := range dash.CurrentMaxes {
		response.CurrentMaxes = append(response.CurrentMaxes, MaxSummaryResponse{
			Lift:  ms.Lift,
			Value: ms.Value,
			Type:  ms.Type,
		})
	}

	return response
}

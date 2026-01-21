package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/waynenilsen/power-pro-v3/internal/domain/loadstrategy"
	"github.com/waynenilsen/power-pro-v3/internal/domain/prescription"
	"github.com/waynenilsen/power-pro-v3/internal/domain/setscheme"
	"github.com/waynenilsen/power-pro-v3/internal/domain/workout"
	"github.com/waynenilsen/power-pro-v3/internal/middleware"
	"github.com/waynenilsen/power-pro-v3/internal/repository"
)

// WorkoutHandler handles HTTP requests for workout generation operations.
type WorkoutHandler struct {
	workoutRepo *repository.WorkoutRepository
	liftLookup  *repository.LiftLookupAdapter
	maxLookup   *repository.MaxLookupAdapter
}

// NewWorkoutHandler creates a new WorkoutHandler.
func NewWorkoutHandler(workoutRepo *repository.WorkoutRepository, sqlDB *sql.DB) *WorkoutHandler {
	return &WorkoutHandler{
		workoutRepo: workoutRepo,
		liftLookup:  repository.NewLiftLookupAdapter(sqlDB),
		maxLookup:   repository.NewMaxLookupAdapter(sqlDB),
	}
}

// WorkoutLiftResponse represents lift info in a workout response.
type WorkoutLiftResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// WorkoutSetResponse represents a set in a workout response.
type WorkoutSetResponse struct {
	SetNumber  int     `json:"setNumber"`
	Weight     float64 `json:"weight"`
	TargetReps int     `json:"targetReps"`
	IsWorkSet  bool    `json:"isWorkSet"`
}

// WorkoutExerciseResponse represents an exercise in a workout response.
type WorkoutExerciseResponse struct {
	PrescriptionID string               `json:"prescriptionId"`
	Lift           WorkoutLiftResponse  `json:"lift"`
	Sets           []WorkoutSetResponse `json:"sets"`
	Notes          string               `json:"notes,omitempty"`
	RestSeconds    *int                 `json:"restSeconds,omitempty"`
}

// WorkoutResponse represents the API response for a generated workout.
type WorkoutResponse struct {
	UserID         string                    `json:"userId"`
	ProgramID      string                    `json:"programId"`
	CycleIteration int                       `json:"cycleIteration"`
	WeekNumber     int                       `json:"weekNumber"`
	DaySlug        string                    `json:"daySlug"`
	Date           string                    `json:"date"`
	Exercises      []WorkoutExerciseResponse `json:"exercises"`
}

func workoutToResponse(w *workout.Workout) WorkoutResponse {
	exercises := make([]WorkoutExerciseResponse, len(w.Exercises))
	for i, e := range w.Exercises {
		sets := make([]WorkoutSetResponse, len(e.Sets))
		for j, s := range e.Sets {
			sets[j] = WorkoutSetResponse{
				SetNumber:  s.SetNumber,
				Weight:     s.Weight,
				TargetReps: s.TargetReps,
				IsWorkSet:  s.IsWorkSet,
			}
		}
		exercises[i] = WorkoutExerciseResponse{
			PrescriptionID: e.PrescriptionID,
			Lift: WorkoutLiftResponse{
				ID:   e.Lift.ID,
				Name: e.Lift.Name,
				Slug: e.Lift.Slug,
			},
			Sets:        sets,
			Notes:       e.Notes,
			RestSeconds: e.RestSeconds,
		}
	}

	return WorkoutResponse{
		UserID:         w.UserID,
		ProgramID:      w.ProgramID,
		CycleIteration: w.CycleIteration,
		WeekNumber:     w.WeekNumber,
		DaySlug:        w.DaySlug,
		Date:           w.Date,
		Exercises:      exercises,
	}
}

// Generate handles GET /users/{userId}/workout
// Generates the current workout for the user based on their program state.
// Optional query params: date, weekNumber, daySlug
func (h *WorkoutHandler) Generate(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "Missing user ID")
		return
	}

	// Authorization check: only the user themselves or an admin can generate workout
	authUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)
	if authUserID != userID && !isAdmin {
		writeError(w, http.StatusForbidden, "Access denied: you can only view your own workouts")
		return
	}

	// Parse optional query parameters
	var weekNumber *int
	var daySlug *string
	var date *string

	if weekStr := r.URL.Query().Get("weekNumber"); weekStr != "" {
		week, err := strconv.Atoi(weekStr)
		if err != nil || week < 1 {
			writeError(w, http.StatusBadRequest, "Invalid weekNumber: must be a positive integer")
			return
		}
		weekNumber = &week
	}

	if ds := r.URL.Query().Get("daySlug"); ds != "" {
		daySlug = &ds
	}

	if d := r.URL.Query().Get("date"); d != "" {
		date = &d
	}

	// Get workout generation data
	data, err := h.workoutRepo.GetWorkoutGenerationData(userID, weekNumber, daySlug)
	if err != nil {
		if errors.Is(err, workout.ErrUserNotEnrolled) {
			writeError(w, http.StatusNotFound, "User is not enrolled in any program")
			return
		}
		if errors.Is(err, workout.ErrWeekNotFound) {
			writeError(w, http.StatusBadRequest, "Week not found in cycle")
			return
		}
		if errors.Is(err, workout.ErrDayNotFound) {
			writeError(w, http.StatusBadRequest, "Day not found for the specified position")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to retrieve workout data")
		return
	}

	if len(data.Prescriptions) == 0 {
		writeError(w, http.StatusNotFound, "Day has no prescriptions")
		return
	}

	// Inject MaxLookup into prescriptions for load strategy resolution
	repository.InjectMaxLookup(data.Prescriptions, h.maxLookup)

	// Determine date
	workoutDate := workout.GetDateString()
	if date != nil {
		workoutDate = *date
	}

	// Build generation context with lookups
	genCtx := workout.GenerationContext{
		LiftLookup:    h.liftLookup,
		SetGenContext: setscheme.DefaultSetGenerationContext(),
	}

	// Build lookup context if lookups are configured
	if data.WeeklyLookup != nil || data.DailyLookup != nil {
		genCtx.LookupContext = &loadstrategy.LookupContext{
			WeekNumber:   data.Enrollment.CurrentWeek,
			DaySlug:      data.Day.Slug,
			WeeklyLookup: data.WeeklyLookup,
			DailyLookup:  data.DailyLookup,
		}
	}

	// Build program context
	programCtx := workout.ProgramContext{
		ProgramID:        data.Enrollment.ProgramID,
		ProgramName:      data.Enrollment.ProgramName,
		CycleID:          data.Enrollment.CycleID,
		CycleLengthWeeks: data.Enrollment.CycleLengthWeeks,
	}

	// Build user state
	userState := workout.UserState{
		CurrentWeek:           data.Enrollment.CurrentWeek,
		CurrentCycleIteration: data.Enrollment.CurrentCycleIteration,
		CurrentDayIndex:       data.Enrollment.CurrentDayIndex,
	}

	// Build day context
	dayCtx := workout.DayContext{
		DayID:   data.Day.ID,
		DaySlug: data.Day.Slug,
		DayName: data.Day.Name,
	}

	// Generate the workout
	generatedWorkout, err := workout.GenerateWorkout(
		context.Background(),
		userID,
		programCtx,
		userState,
		dayCtx,
		data.Prescriptions,
		genCtx,
		workoutDate,
	)
	if err != nil {
		// Check for specific errors
		if errors.Is(err, prescription.ErrMaxNotFound) {
			writeError(w, http.StatusBadRequest, "Missing lift max: set up your training maxes to generate workouts", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to generate workout", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, workoutToResponse(generatedWorkout))
}

// Preview handles GET /users/{userId}/workout/preview
// Previews a workout for a specific week and day without requiring state advancement.
// Required query params: week, day
func (h *WorkoutHandler) Preview(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "Missing user ID")
		return
	}

	// Authorization check: only the user themselves or an admin can preview workout
	authUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)
	if authUserID != userID && !isAdmin {
		writeError(w, http.StatusForbidden, "Access denied: you can only view your own workouts")
		return
	}

	// Parse required query parameters
	weekStr := r.URL.Query().Get("week")
	if weekStr == "" {
		writeError(w, http.StatusBadRequest, "Missing required parameter: week")
		return
	}
	week, err := strconv.Atoi(weekStr)
	if err != nil || week < 1 {
		writeError(w, http.StatusBadRequest, "Invalid week: must be a positive integer")
		return
	}

	daySlug := r.URL.Query().Get("day")
	if daySlug == "" {
		writeError(w, http.StatusBadRequest, "Missing required parameter: day")
		return
	}

	// Get workout generation data
	data, err := h.workoutRepo.GetWorkoutGenerationData(userID, &week, &daySlug)
	if err != nil {
		if errors.Is(err, workout.ErrUserNotEnrolled) {
			writeError(w, http.StatusNotFound, "User is not enrolled in any program")
			return
		}
		if errors.Is(err, workout.ErrWeekNotFound) {
			writeError(w, http.StatusBadRequest, "Week not found in cycle")
			return
		}
		if errors.Is(err, workout.ErrDayNotFound) {
			writeError(w, http.StatusBadRequest, "Day not found for the specified week")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to retrieve workout data")
		return
	}

	if len(data.Prescriptions) == 0 {
		writeError(w, http.StatusNotFound, "Day has no prescriptions")
		return
	}

	// Inject MaxLookup into prescriptions for load strategy resolution
	repository.InjectMaxLookup(data.Prescriptions, h.maxLookup)

	// Build generation context with lookups
	genCtx := workout.GenerationContext{
		LiftLookup:    h.liftLookup,
		SetGenContext: setscheme.DefaultSetGenerationContext(),
	}

	// Build lookup context if lookups are configured
	// For preview, use the specified week/day for lookups
	if data.WeeklyLookup != nil || data.DailyLookup != nil {
		genCtx.LookupContext = &loadstrategy.LookupContext{
			WeekNumber:   week,
			DaySlug:      daySlug,
			WeeklyLookup: data.WeeklyLookup,
			DailyLookup:  data.DailyLookup,
		}
	}

	// Build program context
	programCtx := workout.ProgramContext{
		ProgramID:        data.Enrollment.ProgramID,
		ProgramName:      data.Enrollment.ProgramName,
		CycleID:          data.Enrollment.CycleID,
		CycleLengthWeeks: data.Enrollment.CycleLengthWeeks,
	}

	// Build user state for preview (use specified week, keep iteration from state)
	userState := workout.UserState{
		CurrentWeek:           week,
		CurrentCycleIteration: data.Enrollment.CurrentCycleIteration,
		CurrentDayIndex:       nil, // Not relevant for preview
	}

	// Build day context
	dayCtx := workout.DayContext{
		DayID:   data.Day.ID,
		DaySlug: data.Day.Slug,
		DayName: data.Day.Name,
	}

	// Generate the workout preview (no date needed for preview, use placeholder)
	workoutDate := workout.GetDateString()

	generatedWorkout, err := workout.GenerateWorkout(
		context.Background(),
		userID,
		programCtx,
		userState,
		dayCtx,
		data.Prescriptions,
		genCtx,
		workoutDate,
	)
	if err != nil {
		// Check for specific errors
		if errors.Is(err, prescription.ErrMaxNotFound) {
			writeError(w, http.StatusBadRequest, "Missing lift max: set up your training maxes to generate workouts", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to generate workout preview", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, workoutToResponse(generatedWorkout))
}

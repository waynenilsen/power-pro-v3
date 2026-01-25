// Package e2e provides end-to-end tests for complete program workflows.
// These tests validate entire program configurations from setup through execution.
package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// =============================================================================
// RESPONSE TYPES FOR API DECODING
// =============================================================================

// LiftData matches the lift data format within the response envelope.
type LiftData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// LiftResponse is the standard envelope for single lift responses.
type LiftResponse struct {
	Data LiftData `json:"data"`
}

// LiftMaxData matches the lift max data format.
type LiftMaxData struct {
	ID     string  `json:"id"`
	LiftID string  `json:"liftId"`
	Type   string  `json:"type"`
	Value  float64 `json:"value"`
}

// LiftMaxResponse is the standard envelope for lift max responses.
type LiftMaxResponse struct {
	Data LiftMaxData `json:"data"`
}

// PrescriptionData matches the prescription data format.
type PrescriptionData struct {
	ID string `json:"id"`
}

// PrescriptionResponse is the standard envelope for prescription responses.
type PrescriptionResponse struct {
	Data PrescriptionData `json:"data"`
}

// DayData matches the day data format.
type DayData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// DayResponse is the standard envelope for day responses.
type DayResponse struct {
	Data DayData `json:"data"`
}

// CycleData matches the cycle data format.
type CycleData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	LengthWeeks int    `json:"lengthWeeks"`
}

// CycleResponse is the standard envelope for cycle responses.
type CycleResponse struct {
	Data CycleData `json:"data"`
}

// WeekData matches the week data format.
type WeekData struct {
	ID         string `json:"id"`
	WeekNumber int    `json:"weekNumber"`
}

// WeekResponse is the standard envelope for week responses.
type WeekResponse struct {
	Data WeekData `json:"data"`
}

// ProgramData matches the program data format.
type ProgramData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// ProgramResponse is the standard envelope for program responses.
type ProgramResponse struct {
	Data ProgramData `json:"data"`
}

// ProgressionData matches the progression data format.
type ProgressionData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// ProgressionResponse is the standard envelope for progression responses.
type ProgressionResponse struct {
	Data ProgressionData `json:"data"`
}

// WorkoutSetData represents a set in a workout response.
type WorkoutSetData struct {
	SetNumber  int     `json:"setNumber"`
	Weight     float64 `json:"weight"`
	TargetReps int     `json:"targetReps"`
	IsWorkSet  bool    `json:"isWorkSet"`
}

// WorkoutLiftData represents lift info in a workout response.
type WorkoutLiftData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// WorkoutExerciseData represents an exercise in a workout response.
type WorkoutExerciseData struct {
	PrescriptionID string           `json:"prescriptionId"`
	Lift           WorkoutLiftData  `json:"lift"`
	Sets           []WorkoutSetData `json:"sets"`
	Notes          string           `json:"notes,omitempty"`
	RestSeconds    *int             `json:"restSeconds,omitempty"`
}

// WorkoutData represents the API response for a generated workout.
type WorkoutData struct {
	UserID         string                `json:"userId"`
	ProgramID      string                `json:"programId"`
	CycleIteration int                   `json:"cycleIteration"`
	WeekNumber     int                   `json:"weekNumber"`
	DaySlug        string                `json:"daySlug"`
	Date           string                `json:"date"`
	Exercises      []WorkoutExerciseData `json:"exercises"`
}

// WorkoutResponse is the standard envelope for workout responses.
type WorkoutResponse struct {
	Data WorkoutData `json:"data"`
}

// ManualTriggerRequest matches the API request body for progression triggers.
type ManualTriggerRequest struct {
	ProgressionID string `json:"progressionId"`
	LiftID        string `json:"liftId,omitempty"`
	Force         bool   `json:"force"`
}

// TriggerResultDetail contains the details of an applied progression.
type TriggerResultDetail struct {
	PreviousValue float64   `json:"previousValue"`
	NewValue      float64   `json:"newValue"`
	Delta         float64   `json:"delta"`
	MaxType       string    `json:"maxType"`
	AppliedAt     time.Time `json:"appliedAt"`
}

// TriggerResult represents a single progression result in the API response.
type TriggerResult struct {
	ProgressionID string               `json:"progressionId"`
	LiftID        string               `json:"liftId"`
	Applied       bool                 `json:"applied"`
	Skipped       bool                 `json:"skipped,omitempty"`
	SkipReason    string               `json:"skipReason,omitempty"`
	Result        *TriggerResultDetail `json:"result,omitempty"`
	Error         string               `json:"error,omitempty"`
}

// TriggerResponseData represents the response for manual progression trigger.
type TriggerResponseData struct {
	Results      []TriggerResult `json:"results"`
	TotalApplied int             `json:"totalApplied"`
	TotalSkipped int             `json:"totalSkipped"`
	TotalErrors  int             `json:"totalErrors"`
}

// TriggerResponse is the standard envelope for trigger responses.
type TriggerResponse struct {
	Data TriggerResponseData `json:"data"`
}

// EnrollmentProgramData represents program info in an enrollment response.
type EnrollmentProgramData struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Slug             string  `json:"slug"`
	Description      *string `json:"description,omitempty"`
	CycleLengthWeeks int     `json:"cycleLengthWeeks"`
}

// EnrollmentStateData represents the state portion of an enrollment response.
type EnrollmentStateData struct {
	CurrentWeek           int  `json:"currentWeek"`
	CurrentCycleIteration int  `json:"currentCycleIteration"`
	CurrentDayIndex       *int `json:"currentDayIndex,omitempty"`
}

// CurrentWorkoutSessionData represents the current workout session in an enrollment response.
type CurrentWorkoutSessionData struct {
	ID         string `json:"id"`
	WeekNumber int    `json:"weekNumber"`
	DayIndex   int    `json:"dayIndex"`
	Status     string `json:"status"`
}

// EnrollmentData matches the enrollment data format with all state fields.
type EnrollmentData struct {
	ID                    string                     `json:"id"`
	UserID                string                     `json:"userId"`
	Program               EnrollmentProgramData      `json:"program"`
	State                 EnrollmentStateData        `json:"state"`
	EnrollmentStatus      string                     `json:"enrollmentStatus"`
	CycleStatus           string                     `json:"cycleStatus"`
	WeekStatus            string                     `json:"weekStatus"`
	CurrentWorkoutSession *CurrentWorkoutSessionData `json:"currentWorkoutSession"`
}

// EnrollmentResponse is the standard envelope for enrollment responses.
type EnrollmentResponse struct {
	Data EnrollmentData `json:"data"`
}

// WorkoutSessionData matches the workout session data format.
type WorkoutSessionData struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// WorkoutSessionResponse is the standard envelope for workout session responses.
type WorkoutSessionResponse struct {
	Data WorkoutSessionData `json:"data"`
}

// =============================================================================
// HTTP HELPER FUNCTIONS
// =============================================================================

func adminPost(url string, body string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", testutil.TestAdminID)
	req.Header.Set("X-Admin", "true")
	return http.DefaultClient.Do(req)
}

func userPost(url string, body string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

func userGet(url string, userID string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

func authPostTrigger(url string, body any, userID string) (*http.Response, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)
	return http.DefaultClient.Do(req)
}

// startWorkoutSession starts a new workout session and returns its ID.
func startWorkoutSession(t *testing.T, ts *testutil.TestServer, userID string) string {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, ts.URL("/workouts/start"), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-User-ID", userID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to start workout: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to start workout session, status %d: %s", resp.StatusCode, body)
	}

	var envelope WorkoutSessionResponse
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

// finishWorkoutSession completes a workout session.
func finishWorkoutSession(t *testing.T, ts *testutil.TestServer, sessionID, userID string) {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, ts.URL("/workouts/"+sessionID+"/finish"), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-User-ID", userID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to finish workout: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to finish workout session, status %d: %s", resp.StatusCode, body)
	}
}

// =============================================================================
// STARTING STRENGTH E2E TEST
// =============================================================================

// TestStartingStrengthProgram validates the complete Starting Strength program
// configuration and execution through the API.
//
// Starting Strength characteristics:
// - A/B Rotation: Alternating workouts (A: Squat/Bench/Deadlift, B: Squat/Press/Power Clean)
// - Fixed 3x5: All main lifts use FIXED set scheme with 3 sets of 5 reps
// - LinearProgression: AFTER_SESSION trigger with +5lb for upper body, +10lb for lower body
func TestStartingStrengthProgram(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Test-unique identifiers
	testID := uuid.New().String()[:8]
	// Use a seeded test user (required for foreign key constraints)
	userID := "workout-test-user"

	// Seeded lift IDs
	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"
	deadliftID := "00000000-0000-0000-0000-000000000003"

	// Starting strength training maxes
	squatMax := 225.0   // Starting squat training max
	benchMax := 135.0   // Starting bench training max
	deadliftMax := 275.0 // Starting deadlift training max
	pressMax := 95.0    // Starting press training max
	cleanMax := 135.0   // Starting power clean training max

	// Create additional lifts (Press and Power Clean are not seeded)
	pressSlug := "press-" + testID
	pressBody := fmt.Sprintf(`{"name": "Overhead Press", "slug": "%s", "isCompetitionLift": false}`, pressSlug)
	pressResp, err := adminPost(ts.URL("/lifts"), pressBody)
	if err != nil {
		t.Fatalf("Failed to create press lift: %v", err)
	}
	if pressResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(pressResp.Body)
		pressResp.Body.Close()
		t.Fatalf("Failed to create press lift, status %d: %s", pressResp.StatusCode, body)
	}
	var pressEnvelope LiftResponse
	json.NewDecoder(pressResp.Body).Decode(&pressEnvelope)
	pressResp.Body.Close()
	pressID := pressEnvelope.Data.ID

	cleanSlug := "power-clean-" + testID
	cleanBody := fmt.Sprintf(`{"name": "Power Clean", "slug": "%s", "isCompetitionLift": false}`, cleanSlug)
	cleanResp, err := adminPost(ts.URL("/lifts"), cleanBody)
	if err != nil {
		t.Fatalf("Failed to create power clean lift: %v", err)
	}
	if cleanResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(cleanResp.Body)
		cleanResp.Body.Close()
		t.Fatalf("Failed to create power clean lift, status %d: %s", cleanResp.StatusCode, body)
	}
	var cleanEnvelope LiftResponse
	json.NewDecoder(cleanResp.Body).Decode(&cleanEnvelope)
	cleanResp.Body.Close()
	cleanID := cleanEnvelope.Data.ID

	// Create training maxes for the user
	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", squatMax)
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", benchMax)
	createLiftMax(t, ts, userID, deadliftID, "TRAINING_MAX", deadliftMax)
	createLiftMax(t, ts, userID, pressID, "TRAINING_MAX", pressMax)
	createLiftMax(t, ts, userID, cleanID, "TRAINING_MAX", cleanMax)

	// Create prescriptions for each exercise (FIXED 3x5 at 100% training max)
	squatPrescID := createPrescription(t, ts, squatID, 3, 5, 100.0, 0)
	benchPrescID := createPrescription(t, ts, benchID, 3, 5, 100.0, 1)
	deadliftPrescID := createPrescription(t, ts, deadliftID, 1, 5, 100.0, 2) // 1x5 for deadlift
	pressPrescID := createPrescription(t, ts, pressID, 3, 5, 100.0, 1)
	cleanPrescID := createPrescription(t, ts, cleanID, 5, 3, 100.0, 2) // 5x3 for power clean

	// Create Day A: Squat, Bench, Deadlift
	dayASlug := "day-a-" + testID
	dayABody := fmt.Sprintf(`{"name": "Day A", "slug": "%s"}`, dayASlug)
	dayAResp, _ := adminPost(ts.URL("/days"), dayABody)
	var dayAEnvelope DayResponse
	json.NewDecoder(dayAResp.Body).Decode(&dayAEnvelope)
	dayAResp.Body.Close()
	dayAID := dayAEnvelope.Data.ID

	// Add prescriptions to Day A
	addPrescToDay(t, ts, dayAID, squatPrescID)
	addPrescToDay(t, ts, dayAID, benchPrescID)
	addPrescToDay(t, ts, dayAID, deadliftPrescID)

	// Create Day B: Squat, Press, Power Clean
	dayBSlug := "day-b-" + testID
	dayBBody := fmt.Sprintf(`{"name": "Day B", "slug": "%s"}`, dayBSlug)
	dayBResp, _ := adminPost(ts.URL("/days"), dayBBody)
	var dayBEnvelope DayResponse
	json.NewDecoder(dayBResp.Body).Decode(&dayBEnvelope)
	dayBResp.Body.Close()
	dayBID := dayBEnvelope.Data.ID

	// Add prescriptions to Day B (Squat is shared, Press and Clean are B-only)
	// Need separate squat prescription for Day B to allow independent tracking
	squatPrescBID := createPrescription(t, ts, squatID, 3, 5, 100.0, 0)
	addPrescToDay(t, ts, dayBID, squatPrescBID)
	addPrescToDay(t, ts, dayBID, pressPrescID)
	addPrescToDay(t, ts, dayBID, cleanPrescID)

	// Create 1-week cycle with A/B/A pattern (Mon/Wed/Fri)
	cycleName := "SS Cycle " + testID
	cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 1}`, cycleName)
	cycleResp, _ := adminPost(ts.URL("/cycles"), cycleBody)
	var cycleEnvelope CycleResponse
	json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
	cycleResp.Body.Close()
	cycleID := cycleEnvelope.Data.ID

	// Create week 1 in the cycle
	weekBody := fmt.Sprintf(`{"weekNumber": 1, "cycleId": "%s"}`, cycleID)
	weekResp, _ := adminPost(ts.URL("/weeks"), weekBody)
	var weekEnvelope WeekResponse
	json.NewDecoder(weekResp.Body).Decode(&weekEnvelope)
	weekResp.Body.Close()
	weekID := weekEnvelope.Data.ID

	// Add days to week: A/B/A pattern
	addDayToWeek(t, ts, weekID, dayAID, "MONDAY")
	addDayToWeek(t, ts, weekID, dayBID, "WEDNESDAY")
	addDayToWeek(t, ts, weekID, dayAID, "FRIDAY")

	// Create program
	programSlug := "starting-strength-" + testID
	programBody := fmt.Sprintf(`{"name": "Starting Strength", "slug": "%s", "cycleId": "%s"}`, programSlug, cycleID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID := programEnvelope.Data.ID

	// Create Linear Progressions
	// Upper body progression (+5lb)
	upperProgBody := `{"name": "SS Upper Linear", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"}}`
	upperProgResp, _ := adminPost(ts.URL("/progressions"), upperProgBody)
	var upperProgEnvelope ProgressionResponse
	json.NewDecoder(upperProgResp.Body).Decode(&upperProgEnvelope)
	upperProgResp.Body.Close()
	upperProgID := upperProgEnvelope.Data.ID

	// Lower body progression (+10lb)
	lowerProgBody := `{"name": "SS Lower Linear", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 10.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"}}`
	lowerProgResp, _ := adminPost(ts.URL("/progressions"), lowerProgBody)
	var lowerProgEnvelope ProgressionResponse
	json.NewDecoder(lowerProgResp.Body).Decode(&lowerProgEnvelope)
	lowerProgResp.Body.Close()
	lowerProgID := lowerProgEnvelope.Data.ID

	// Link progressions to program
	// Lower body lifts get +10lb
	linkProgressionToProgram(t, ts, programID, lowerProgID, squatID, 1)
	linkProgressionToProgram(t, ts, programID, lowerProgID, deadliftID, 2)
	// Upper body lifts get +5lb
	linkProgressionToProgram(t, ts, programID, upperProgID, benchID, 3)
	linkProgressionToProgram(t, ts, programID, upperProgID, pressID, 4)
	linkProgressionToProgram(t, ts, programID, upperProgID, cleanID, 5)

	// Enroll user in program
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, err := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	if err != nil {
		t.Fatalf("Failed to enroll user: %v", err)
	}
	if enrollResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(enrollResp.Body)
		enrollResp.Body.Close()
		t.Fatalf("Failed to enroll user, status %d: %s", enrollResp.StatusCode, body)
	}
	enrollResp.Body.Close()

	// =============================================================================
	// EXECUTION PHASE: Day A (Workout 1)
	// Using explicit state machine flow: start workout -> log sets -> finish workout
	// =============================================================================
	var dayASessionID string
	t.Run("Day A workout generates correct exercises and weights", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		if workoutResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(workoutResp.Body)
			t.Fatalf("Expected status 200, got %d: %s", workoutResp.StatusCode, body)
		}

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		// Verify Day A structure
		if workout.Data.DaySlug != dayASlug {
			t.Errorf("Expected Day A slug '%s', got '%s'", dayASlug, workout.Data.DaySlug)
		}

		if len(workout.Data.Exercises) != 3 {
			t.Fatalf("Expected 3 exercises on Day A, got %d", len(workout.Data.Exercises))
		}

		// Verify exercise order and weights
		exercisesByLift := make(map[string]WorkoutExerciseData)
		for _, ex := range workout.Data.Exercises {
			exercisesByLift[ex.Lift.ID] = ex
		}

		// Squat: 3x5 @ 225 (100% of training max)
		if squat, ok := exercisesByLift[squatID]; ok {
			if len(squat.Sets) != 3 {
				t.Errorf("Squat: expected 3 sets, got %d", len(squat.Sets))
			}
			for i, set := range squat.Sets {
				if set.Weight != squatMax {
					t.Errorf("Squat set %d: expected weight %f, got %f", i+1, squatMax, set.Weight)
				}
				if set.TargetReps != 5 {
					t.Errorf("Squat set %d: expected 5 reps, got %d", i+1, set.TargetReps)
				}
			}
		} else {
			t.Error("Day A missing Squat exercise")
		}

		// Bench: 3x5 @ 135
		if bench, ok := exercisesByLift[benchID]; ok {
			if len(bench.Sets) != 3 {
				t.Errorf("Bench: expected 3 sets, got %d", len(bench.Sets))
			}
			for i, set := range bench.Sets {
				if set.Weight != benchMax {
					t.Errorf("Bench set %d: expected weight %f, got %f", i+1, benchMax, set.Weight)
				}
				if set.TargetReps != 5 {
					t.Errorf("Bench set %d: expected 5 reps, got %d", i+1, set.TargetReps)
				}
			}
		} else {
			t.Error("Day A missing Bench exercise")
		}

		// Deadlift: 1x5 @ 275
		if deadlift, ok := exercisesByLift[deadliftID]; ok {
			if len(deadlift.Sets) != 1 {
				t.Errorf("Deadlift: expected 1 set, got %d", len(deadlift.Sets))
			}
			if len(deadlift.Sets) > 0 {
				if deadlift.Sets[0].Weight != deadliftMax {
					t.Errorf("Deadlift: expected weight %f, got %f", deadliftMax, deadlift.Sets[0].Weight)
				}
				if deadlift.Sets[0].TargetReps != 5 {
					t.Errorf("Deadlift: expected 5 reps, got %d", deadlift.Sets[0].TargetReps)
				}
			}
		} else {
			t.Error("Day A missing Deadlift exercise")
		}
	})

	// Start workout session, log sets, and finish - progression applies automatically (AFTER_SESSION trigger)
	t.Run("Day A complete workout with auto-progression", func(t *testing.T) {
		// Start workout session
		dayASessionID = startWorkoutSession(t, ts, userID)

		// Get workout to find prescription IDs
		workoutResp, _ := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)
		workoutResp.Body.Close()

		// Log sets for each exercise (successful completion)
		for _, ex := range workout.Data.Exercises {
			for _, set := range ex.Sets {
				logSSSet(t, ts, userID, dayASessionID, ex.PrescriptionID, ex.Lift.ID, set.SetNumber, set.Weight, set.TargetReps, set.TargetReps)
			}
		}

		// Finish workout
		finishWorkoutSession(t, ts, dayASessionID, userID)

		// Trigger progressions for Day A lifts (Squat, Bench, Deadlift)
		triggerProgressionForLift(t, ts, userID, lowerProgID, squatID)    // +10lb
		triggerProgressionForLift(t, ts, userID, upperProgID, benchID)    // +5lb
		triggerProgressionForLift(t, ts, userID, lowerProgID, deadliftID) // +10lb

		// Advance to next day
		advanceUserState(t, ts, userID)
	})

	// =============================================================================
	// EXECUTION PHASE: Day B (Workout 2)
	// Using explicit state machine flow: start workout -> log sets -> finish workout
	// =============================================================================
	var dayBSessionID string
	t.Run("Day B workout shows updated squat weight and correct B exercises", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		if workoutResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(workoutResp.Body)
			t.Fatalf("Expected status 200, got %d: %s", workoutResp.StatusCode, body)
		}

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		// Verify Day B structure
		if workout.Data.DaySlug != dayBSlug {
			t.Errorf("Expected Day B slug '%s', got '%s'", dayBSlug, workout.Data.DaySlug)
		}

		if len(workout.Data.Exercises) != 3 {
			t.Fatalf("Expected 3 exercises on Day B, got %d", len(workout.Data.Exercises))
		}

		// Verify exercise order and weights
		exercisesByLift := make(map[string]WorkoutExerciseData)
		for _, ex := range workout.Data.Exercises {
			exercisesByLift[ex.Lift.ID] = ex
		}

		// Squat should now be at 235 (+10 from Day A progression)
		expectedSquat := squatMax + 10.0
		if squat, ok := exercisesByLift[squatID]; ok {
			if len(squat.Sets) != 3 {
				t.Errorf("Squat: expected 3 sets, got %d", len(squat.Sets))
			}
			for i, set := range squat.Sets {
				if set.Weight != expectedSquat {
					t.Errorf("Squat set %d: expected weight %f (increased from Day A), got %f", i+1, expectedSquat, set.Weight)
				}
			}
		} else {
			t.Error("Day B missing Squat exercise")
		}

		// Press: 3x5 @ 95 (unchanged, not performed on Day A)
		if press, ok := exercisesByLift[pressID]; ok {
			if len(press.Sets) != 3 {
				t.Errorf("Press: expected 3 sets, got %d", len(press.Sets))
			}
			for i, set := range press.Sets {
				if set.Weight != pressMax {
					t.Errorf("Press set %d: expected weight %f, got %f", i+1, pressMax, set.Weight)
				}
				if set.TargetReps != 5 {
					t.Errorf("Press set %d: expected 5 reps, got %d", i+1, set.TargetReps)
				}
			}
		} else {
			t.Error("Day B missing Press exercise")
		}

		// Power Clean: 5x3 @ 135
		if clean, ok := exercisesByLift[cleanID]; ok {
			if len(clean.Sets) != 5 {
				t.Errorf("Power Clean: expected 5 sets, got %d", len(clean.Sets))
			}
			for i, set := range clean.Sets {
				if set.Weight != cleanMax {
					t.Errorf("Power Clean set %d: expected weight %f, got %f", i+1, cleanMax, set.Weight)
				}
				if set.TargetReps != 3 {
					t.Errorf("Power Clean set %d: expected 3 reps, got %d", i+1, set.TargetReps)
				}
			}
		} else {
			t.Error("Day B missing Power Clean exercise")
		}
	})

	// Start workout session, log sets, and finish - progression applies automatically (AFTER_SESSION trigger)
	t.Run("Day B complete workout with auto-progression", func(t *testing.T) {
		// Start workout session
		dayBSessionID = startWorkoutSession(t, ts, userID)

		// Get workout to find prescription IDs
		workoutResp, _ := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)
		workoutResp.Body.Close()

		// Log sets for each exercise (successful completion)
		for _, ex := range workout.Data.Exercises {
			for _, set := range ex.Sets {
				logSSSet(t, ts, userID, dayBSessionID, ex.PrescriptionID, ex.Lift.ID, set.SetNumber, set.Weight, set.TargetReps, set.TargetReps)
			}
		}

		// Finish workout
		finishWorkoutSession(t, ts, dayBSessionID, userID)

		// Trigger progressions for Day B lifts (Squat, Press, Power Clean)
		triggerProgressionForLift(t, ts, userID, lowerProgID, squatID) // +10lb (another squat session)
		triggerProgressionForLift(t, ts, userID, upperProgID, pressID) // +5lb
		triggerProgressionForLift(t, ts, userID, upperProgID, cleanID) // +5lb

		// Advance to next day
		advanceUserState(t, ts, userID)
	})

	// =============================================================================
	// VALIDATION PHASE: Day A again (Workout 3) with all accumulated progressions
	// =============================================================================
	t.Run("Day A second time shows all accumulated progressions", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		if workoutResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(workoutResp.Body)
			t.Fatalf("Expected status 200, got %d: %s", workoutResp.StatusCode, body)
		}

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		// Should be back to Day A
		if workout.Data.DaySlug != dayASlug {
			t.Errorf("Expected Day A slug '%s', got '%s'", dayASlug, workout.Data.DaySlug)
		}

		exercisesByLift := make(map[string]WorkoutExerciseData)
		for _, ex := range workout.Data.Exercises {
			exercisesByLift[ex.Lift.ID] = ex
		}

		// Squat: now at 245 (225 + 10 + 10)
		expectedSquat := squatMax + 20.0
		if squat, ok := exercisesByLift[squatID]; ok {
			for i, set := range squat.Sets {
				if set.Weight != expectedSquat {
					t.Errorf("Squat set %d: expected weight %f, got %f", i+1, expectedSquat, set.Weight)
				}
			}
		}

		// Bench: now at 140 (135 + 5)
		expectedBench := benchMax + 5.0
		if bench, ok := exercisesByLift[benchID]; ok {
			for i, set := range bench.Sets {
				if set.Weight != expectedBench {
					t.Errorf("Bench set %d: expected weight %f, got %f", i+1, expectedBench, set.Weight)
				}
			}
		}

		// Deadlift: now at 285 (275 + 10)
		expectedDeadlift := deadliftMax + 10.0
		if deadlift, ok := exercisesByLift[deadliftID]; ok {
			for _, set := range deadlift.Sets {
				if set.Weight != expectedDeadlift {
					t.Errorf("Deadlift: expected weight %f, got %f", expectedDeadlift, set.Weight)
				}
			}
		}
	})
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

func createLiftMax(t *testing.T, ts *testutil.TestServer, userID, liftID, maxType string, value float64) {
	t.Helper()
	body := fmt.Sprintf(`{"liftId": "%s", "type": "%s", "value": %f}`, liftID, maxType, value)
	resp, err := userPost(ts.URL("/users/"+userID+"/lift-maxes"), body, userID)
	if err != nil {
		t.Fatalf("Failed to create lift max: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create lift max, status %d: %s", resp.StatusCode, bodyBytes)
	}
}

func createPrescription(t *testing.T, ts *testutil.TestServer, liftID string, sets, reps int, percentage float64, order int) string {
	t.Helper()
	body := fmt.Sprintf(`{
		"liftId": "%s",
		"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": %f},
		"setScheme": {"type": "FIXED", "sets": %d, "reps": %d},
		"order": %d
	}`, liftID, percentage, sets, reps, order)

	resp, err := adminPost(ts.URL("/prescriptions"), body)
	if err != nil {
		t.Fatalf("Failed to create prescription: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create prescription, status %d: %s", resp.StatusCode, bodyBytes)
	}

	var envelope PrescriptionResponse
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

func addPrescToDay(t *testing.T, ts *testutil.TestServer, dayID, prescID string) {
	t.Helper()
	body := fmt.Sprintf(`{"prescriptionId": "%s"}`, prescID)
	resp, err := adminPost(ts.URL("/days/"+dayID+"/prescriptions"), body)
	if err != nil {
		t.Fatalf("Failed to add prescription to day: %v", err)
	}
	resp.Body.Close()
}

func addDayToWeek(t *testing.T, ts *testutil.TestServer, weekID, dayID, dayOfWeek string) {
	t.Helper()
	body := fmt.Sprintf(`{"dayId": "%s", "dayOfWeek": "%s"}`, dayID, dayOfWeek)
	resp, err := adminPost(ts.URL("/weeks/"+weekID+"/days"), body)
	if err != nil {
		t.Fatalf("Failed to add day to week: %v", err)
	}
	resp.Body.Close()
}

func linkProgressionToProgram(t *testing.T, ts *testutil.TestServer, programID, progressionID, liftID string, priority int) {
	t.Helper()
	body := fmt.Sprintf(`{"progressionId": "%s", "liftId": "%s", "priority": %d, "enabled": true}`, progressionID, liftID, priority)
	resp, err := adminPost(ts.URL("/programs/"+programID+"/progressions"), body)
	if err != nil {
		t.Fatalf("Failed to link progression to program: %v", err)
	}
	resp.Body.Close()
}

// advanceUserState advances the user's program state to the next workout.
// Note: In the new state machine flow, this is typically handled automatically by finishWorkoutSession.
// This helper is kept for backwards compatibility with other tests.
func advanceUserState(t *testing.T, ts *testutil.TestServer, userID string) {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, ts.URL("/users/"+userID+"/program-state/advance"), nil)
	if err != nil {
		t.Fatalf("Failed to create advance request: %v", err)
	}
	req.Header.Set("X-User-ID", userID)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to advance user state: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to advance state, status %d: %s", resp.StatusCode, bodyBytes)
	}
}

// logSSSet logs a single set for Starting Strength workout.
func logSSSet(t *testing.T, ts *testutil.TestServer, userID, sessionID, prescriptionID, liftID string, setNumber int, weight float64, targetReps, repsPerformed int) {
	t.Helper()

	type setRequest struct {
		PrescriptionID string  `json:"prescriptionId"`
		LiftID         string  `json:"liftId"`
		SetNumber      int     `json:"setNumber"`
		Weight         float64 `json:"weight"`
		TargetReps     int     `json:"targetReps"`
		RepsPerformed  int     `json:"repsPerformed"`
	}

	setsReq := []setRequest{{
		PrescriptionID: prescriptionID,
		LiftID:         liftID,
		SetNumber:      setNumber,
		Weight:         weight,
		TargetReps:     targetReps,
		RepsPerformed:  repsPerformed,
	}}

	body, _ := json.Marshal(map[string]interface{}{"sets": setsReq})
	req, _ := http.NewRequest(http.MethodPost, ts.URL("/sessions/"+sessionID+"/sets"), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to log set: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to log set, status %d: %s", resp.StatusCode, respBody)
	}
}

// triggerProgressionForLift triggers a progression for a specific lift and returns the response.
// NOTE: Uses Force: true because auto-progression via events is not yet fully implemented.
// This should be changed to Force: false once event-driven progressions are working.
func triggerProgressionForLift(t *testing.T, ts *testutil.TestServer, userID, progressionID, liftID string) TriggerResponse {
	t.Helper()
	reqBody := ManualTriggerRequest{
		ProgressionID: progressionID,
		LiftID:        liftID,
		Force:         true,
	}
	resp, err := authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), reqBody, userID)
	if err != nil {
		t.Fatalf("Failed to trigger progression: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to trigger progression, status %d: %s", resp.StatusCode, bodyBytes)
	}

	var triggerResp TriggerResponse
	json.NewDecoder(resp.Body).Decode(&triggerResp)
	return triggerResp
}

// =============================================================================
// STATE ASSERTION HELPERS
// =============================================================================

// ExpectedEnrollmentState defines expected state values for assertions.
type ExpectedEnrollmentState struct {
	EnrollmentStatus string
	CycleStatus      string
	WeekStatus       string
	CurrentWeek      int
	CycleIteration   int
	HasActiveSession bool
	SessionStatus    string // Optional, only checked if HasActiveSession is true
}

// assertEnrollmentState verifies all state fields in an enrollment response.
func assertEnrollmentState(t *testing.T, enrollment EnrollmentData, expected ExpectedEnrollmentState) {
	t.Helper()

	if enrollment.EnrollmentStatus != expected.EnrollmentStatus {
		t.Errorf("EnrollmentStatus: expected %q, got %q", expected.EnrollmentStatus, enrollment.EnrollmentStatus)
	}
	if enrollment.CycleStatus != expected.CycleStatus {
		t.Errorf("CycleStatus: expected %q, got %q", expected.CycleStatus, enrollment.CycleStatus)
	}
	if enrollment.WeekStatus != expected.WeekStatus {
		t.Errorf("WeekStatus: expected %q, got %q", expected.WeekStatus, enrollment.WeekStatus)
	}
	if enrollment.State.CurrentWeek != expected.CurrentWeek {
		t.Errorf("CurrentWeek: expected %d, got %d", expected.CurrentWeek, enrollment.State.CurrentWeek)
	}
	if enrollment.State.CurrentCycleIteration != expected.CycleIteration {
		t.Errorf("CycleIteration: expected %d, got %d", expected.CycleIteration, enrollment.State.CurrentCycleIteration)
	}

	hasSession := enrollment.CurrentWorkoutSession != nil
	if hasSession != expected.HasActiveSession {
		t.Errorf("HasActiveSession: expected %v, got %v", expected.HasActiveSession, hasSession)
	}

	if expected.HasActiveSession && hasSession && expected.SessionStatus != "" {
		if enrollment.CurrentWorkoutSession.Status != expected.SessionStatus {
			t.Errorf("SessionStatus: expected %q, got %q", expected.SessionStatus, enrollment.CurrentWorkoutSession.Status)
		}
	}
}

// getEnrollment fetches the current enrollment state for a user.
func getEnrollment(t *testing.T, ts *testutil.TestServer, userID string) EnrollmentData {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, ts.URL("/users/"+userID+"/program"), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-User-ID", userID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to get enrollment: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to get enrollment, status %d: %s", resp.StatusCode, body)
	}

	var envelope EnrollmentResponse
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		t.Fatalf("Failed to decode enrollment response: %v", err)
	}
	return envelope.Data
}

// enrollUser enrolls a user in a program and returns the enrollment data.
func enrollUser(t *testing.T, ts *testutil.TestServer, userID, programID string) EnrollmentData {
	t.Helper()
	body := fmt.Sprintf(`{"programId": "%s"}`, programID)
	resp, err := userPost(ts.URL("/users/"+userID+"/program"), body, userID)
	if err != nil {
		t.Fatalf("Failed to enroll user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to enroll user, status %d: %s", resp.StatusCode, bodyBytes)
	}

	var envelope EnrollmentResponse
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		t.Fatalf("Failed to decode enrollment response: %v", err)
	}
	return envelope.Data
}

// unenrollUser removes enrollment for a user.
func unenrollUser(t *testing.T, ts *testutil.TestServer, userID string) {
	t.Helper()
	req, err := http.NewRequest(http.MethodDelete, ts.URL("/users/"+userID+"/program"), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-User-ID", userID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to unenroll user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to unenroll user, status %d: %s", resp.StatusCode, body)
	}
}

// startWorkoutAndVerify starts a workout and verifies state transitions.
// Returns the session ID of the started workout.
func startWorkoutAndVerify(t *testing.T, ts *testutil.TestServer, userID string,
	expectedCycleStatus, expectedWeekStatus string) string {
	t.Helper()

	// Start workout
	sessionID := startWorkoutSession(t, ts, userID)

	// Verify enrollment state changed appropriately
	enrollment := getEnrollment(t, ts, userID)

	if enrollment.CycleStatus != expectedCycleStatus {
		t.Errorf("After starting workout, CycleStatus: expected %q, got %q", expectedCycleStatus, enrollment.CycleStatus)
	}
	if enrollment.WeekStatus != expectedWeekStatus {
		t.Errorf("After starting workout, WeekStatus: expected %q, got %q", expectedWeekStatus, enrollment.WeekStatus)
	}
	if enrollment.CurrentWorkoutSession == nil {
		t.Error("After starting workout, expected an active session but got none")
	} else if enrollment.CurrentWorkoutSession.Status != "IN_PROGRESS" {
		t.Errorf("After starting workout, SessionStatus: expected %q, got %q", "IN_PROGRESS", enrollment.CurrentWorkoutSession.Status)
	}

	return sessionID
}

// finishWorkoutAndVerify finishes a workout and verifies state transitions.
func finishWorkoutAndVerify(t *testing.T, ts *testutil.TestServer, sessionID, userID string,
	expectedEnrollmentStatus, expectedCycleStatus, expectedWeekStatus string) {
	t.Helper()

	// Finish workout
	finishWorkoutSession(t, ts, sessionID, userID)

	// Verify enrollment state changed appropriately
	enrollment := getEnrollment(t, ts, userID)

	if enrollment.EnrollmentStatus != expectedEnrollmentStatus {
		t.Errorf("After finishing workout, EnrollmentStatus: expected %q, got %q", expectedEnrollmentStatus, enrollment.EnrollmentStatus)
	}
	if enrollment.CycleStatus != expectedCycleStatus {
		t.Errorf("After finishing workout, CycleStatus: expected %q, got %q", expectedCycleStatus, enrollment.CycleStatus)
	}
	if enrollment.WeekStatus != expectedWeekStatus {
		t.Errorf("After finishing workout, WeekStatus: expected %q, got %q", expectedWeekStatus, enrollment.WeekStatus)
	}
}

// advanceWeek advances to the next week and returns updated enrollment.
func advanceWeek(t *testing.T, ts *testutil.TestServer, userID string) EnrollmentData {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, ts.URL("/users/"+userID+"/enrollment/advance-week"), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-User-ID", userID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to advance week: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to advance week, status %d: %s", resp.StatusCode, body)
	}

	var envelope EnrollmentResponse
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		t.Fatalf("Failed to decode enrollment response: %v", err)
	}
	return envelope.Data
}

// startNextCycle starts a new cycle when in BETWEEN_CYCLES state.
func startNextCycle(t *testing.T, ts *testutil.TestServer, userID string) EnrollmentData {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, ts.URL("/users/"+userID+"/enrollment/next-cycle"), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-User-ID", userID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to start next cycle: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to start next cycle, status %d: %s", resp.StatusCode, body)
	}

	var envelope EnrollmentResponse
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		t.Fatalf("Failed to decode enrollment response: %v", err)
	}
	return envelope.Data
}

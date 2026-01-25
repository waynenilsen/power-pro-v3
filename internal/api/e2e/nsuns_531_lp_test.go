// Package e2e provides end-to-end tests for complete program workflows.
// These tests validate entire program configurations from setup through execution.
package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// =============================================================================
// NSUNS 5/3/1 LP 5-DAY E2E TEST
// =============================================================================

// TestNSuns531LP5DayProgram validates the complete nSuns 5/3/1 LP 5-Day program
// configuration and execution through the API.
//
// nSuns 5/3/1 LP characteristics:
// - 1-week cycle with weekly linear progression
// - Multiple AMRAP sets per day (1+ sets on primary lifts)
// - AMRAPProgression with threshold-based increments (stored for future automatic triggering)
// - AMRAP logged sets are recorded with isAmrap flag
//
// This test demonstrates:
// - AMRAP set scheme generation
// - Logged set persistence with AMRAP flag
// - Program structure with 5-day weekly cycle
// - AMRAP progression storage (full automatic triggering is a future enhancement)
func TestNSuns531LP5DayProgram(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Test-unique identifiers
	testID := uuid.New().String()[:8]
	userID := "nsuns-lp-test-user"

	// Seeded lift IDs
	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"
	deadliftID := "00000000-0000-0000-0000-000000000003"

	// nSuns training maxes
	benchTM := 225.0    // Bench training max
	squatTM := 315.0    // Squat training max
	deadliftTM := 365.0 // Deadlift training max

	// Create training maxes for the user
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", benchTM)
	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", squatTM)
	createLiftMax(t, ts, userID, deadliftID, "TRAINING_MAX", deadliftTM)

	// =============================================================================
	// Create AMRAP prescriptions for each main lift
	// Each day has a 1+ AMRAP set at varying percentages
	// =============================================================================

	// Day 1: Bench 1+ (AMRAP at 95%)
	benchDay1PrescID := createAMRAPPrescription(t, ts, benchID, 1, 1, 95.0, 0)

	// Day 2: Squat 1+ (AMRAP at 95%)
	squatDay2PrescID := createAMRAPPrescription(t, ts, squatID, 1, 1, 95.0, 0)

	// Day 3: OHP (not testing, use bench for variety)
	// Day 4: Deadlift 1+ (AMRAP at 95%)
	deadliftDay4PrescID := createAMRAPPrescription(t, ts, deadliftID, 1, 1, 95.0, 0)

	// Day 5: Bench volume (use same AMRAP pattern)
	benchDay5PrescID := createAMRAPPrescription(t, ts, benchID, 1, 1, 90.0, 0)

	// =============================================================================
	// Create Days
	// =============================================================================

	// Day 1: Bench Day
	day1Slug := "day1-bench-" + testID
	day1Body := fmt.Sprintf(`{"name": "Day 1 - Bench", "slug": "%s"}`, day1Slug)
	day1Resp, _ := adminPost(ts.URL("/days"), day1Body)
	var day1Envelope DayResponse
	json.NewDecoder(day1Resp.Body).Decode(&day1Envelope)
	day1Resp.Body.Close()
	day1ID := day1Envelope.Data.ID
	addPrescToDay(t, ts, day1ID, benchDay1PrescID)

	// Day 2: Squat Day
	day2Slug := "day2-squat-" + testID
	day2Body := fmt.Sprintf(`{"name": "Day 2 - Squat", "slug": "%s"}`, day2Slug)
	day2Resp, _ := adminPost(ts.URL("/days"), day2Body)
	var day2Envelope DayResponse
	json.NewDecoder(day2Resp.Body).Decode(&day2Envelope)
	day2Resp.Body.Close()
	day2ID := day2Envelope.Data.ID
	addPrescToDay(t, ts, day2ID, squatDay2PrescID)

	// Day 3: OHP Day (use bench as placeholder)
	day3Slug := "day3-ohp-" + testID
	day3Body := fmt.Sprintf(`{"name": "Day 3 - OHP", "slug": "%s"}`, day3Slug)
	day3Resp, _ := adminPost(ts.URL("/days"), day3Body)
	var day3Envelope DayResponse
	json.NewDecoder(day3Resp.Body).Decode(&day3Envelope)
	day3Resp.Body.Close()
	day3ID := day3Envelope.Data.ID
	// Add a simple fixed prescription for day 3
	ohpPrescID := createPrescription(t, ts, benchID, 3, 5, 70.0, 0)
	addPrescToDay(t, ts, day3ID, ohpPrescID)

	// Day 4: Deadlift Day
	day4Slug := "day4-deadlift-" + testID
	day4Body := fmt.Sprintf(`{"name": "Day 4 - Deadlift", "slug": "%s"}`, day4Slug)
	day4Resp, _ := adminPost(ts.URL("/days"), day4Body)
	var day4Envelope DayResponse
	json.NewDecoder(day4Resp.Body).Decode(&day4Envelope)
	day4Resp.Body.Close()
	day4ID := day4Envelope.Data.ID
	addPrescToDay(t, ts, day4ID, deadliftDay4PrescID)

	// Day 5: Bench Volume Day
	day5Slug := "day5-bench-vol-" + testID
	day5Body := fmt.Sprintf(`{"name": "Day 5 - Bench Volume", "slug": "%s"}`, day5Slug)
	day5Resp, _ := adminPost(ts.URL("/days"), day5Body)
	var day5Envelope DayResponse
	json.NewDecoder(day5Resp.Body).Decode(&day5Envelope)
	day5Resp.Body.Close()
	day5ID := day5Envelope.Data.ID
	addPrescToDay(t, ts, day5ID, benchDay5PrescID)

	// =============================================================================
	// Create 1-week cycle
	// =============================================================================
	cycleName := "nSuns LP Cycle " + testID
	cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 1}`, cycleName)
	cycleResp, _ := adminPost(ts.URL("/cycles"), cycleBody)
	var cycleEnvelope CycleResponse
	json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
	cycleResp.Body.Close()
	cycleID := cycleEnvelope.Data.ID

	// Create week 1
	weekBody := fmt.Sprintf(`{"weekNumber": 1, "cycleId": "%s"}`, cycleID)
	weekResp, _ := adminPost(ts.URL("/weeks"), weekBody)
	var weekEnvelope WeekResponse
	json.NewDecoder(weekResp.Body).Decode(&weekEnvelope)
	weekResp.Body.Close()
	weekID := weekEnvelope.Data.ID

	// Add all 5 days to the week
	addDayToWeek(t, ts, weekID, day1ID, "MONDAY")
	addDayToWeek(t, ts, weekID, day2ID, "TUESDAY")
	addDayToWeek(t, ts, weekID, day3ID, "WEDNESDAY")
	addDayToWeek(t, ts, weekID, day4ID, "THURSDAY")
	addDayToWeek(t, ts, weekID, day5ID, "FRIDAY")

	// =============================================================================
	// Create Program
	// =============================================================================
	programSlug := "nsuns-531-lp-" + testID
	programBody := fmt.Sprintf(`{"name": "nSuns 5/3/1 LP 5-Day", "slug": "%s", "cycleId": "%s"}`,
		programSlug, cycleID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID := programEnvelope.Data.ID

	// =============================================================================
	// Create AMRAPProgression with threshold-based increments
	// This is the key feature of nSuns - performance-based progression
	// The progression is stored and will be used for automatic triggering when
	// AMRAP sets are logged (future enhancement)
	// =============================================================================
	amrapProgBody := `{
		"name": "nSuns AMRAP Progression",
		"type": "AMRAP_PROGRESSION",
		"parameters": {
			"maxType": "TRAINING_MAX",
			"triggerType": "AFTER_SET",
			"thresholds": [
				{"minReps": 2, "increment": 5.0},
				{"minReps": 4, "increment": 10.0},
				{"minReps": 6, "increment": 15.0}
			]
		}
	}`
	amrapProgResp, err := adminPost(ts.URL("/progressions"), amrapProgBody)
	if err != nil {
		t.Fatalf("Failed to create AMRAP progression: %v", err)
	}
	if amrapProgResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(amrapProgResp.Body)
		amrapProgResp.Body.Close()
		t.Fatalf("Failed to create AMRAP progression, status %d: %s", amrapProgResp.StatusCode, body)
	}
	var amrapProgEnvelope ProgressionResponse
	json.NewDecoder(amrapProgResp.Body).Decode(&amrapProgEnvelope)
	amrapProgResp.Body.Close()
	amrapProgID := amrapProgEnvelope.Data.ID

	// Link AMRAP progression to program for each lift
	linkProgressionToProgram(t, ts, programID, amrapProgID, benchID, 1)
	linkProgressionToProgram(t, ts, programID, amrapProgID, squatID, 2)
	linkProgressionToProgram(t, ts, programID, amrapProgID, deadliftID, 3)

	// =============================================================================
	// Enroll user in program
	// =============================================================================
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
	// DAY 1: Generate Bench workout with AMRAP set
	// =============================================================================
	t.Run("Day 1 generates Bench AMRAP at 95%", func(t *testing.T) {
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

		if workout.Data.DaySlug != day1Slug {
			t.Errorf("Expected day slug '%s', got '%s'", day1Slug, workout.Data.DaySlug)
		}

		if len(workout.Data.Exercises) < 1 {
			t.Fatalf("Expected at least 1 exercise, got %d", len(workout.Data.Exercises))
		}

		bench := workout.Data.Exercises[0]
		if bench.Lift.ID != benchID {
			t.Errorf("Expected bench lift, got %s", bench.Lift.ID)
		}

		// AMRAP prescription should have 1 set
		if len(bench.Sets) != 1 {
			t.Errorf("Expected 1 AMRAP set, got %d", len(bench.Sets))
		}

		// Verify weight is approximately 95% of TM
		expectedWeight := benchTM * 0.95 // 225 * 0.95 = 213.75
		if !withinTolerance(bench.Sets[0].Weight, expectedWeight, 5.0) {
			t.Errorf("Expected weight ~%.1f (95%% of TM), got %.1f", expectedWeight, bench.Sets[0].Weight)
		}

		// Target reps should be 1 (1+ AMRAP)
		if bench.Sets[0].TargetReps != 1 {
			t.Errorf("Expected target reps 1, got %d", bench.Sets[0].TargetReps)
		}
	})

	// =============================================================================
	// Log Bench AMRAP with 5 reps
	// =============================================================================
	sessionID := startWorkoutSession(t, ts, userID)
	t.Run("Log Day 1 Bench AMRAP with 5 reps", func(t *testing.T) {
		workoutResp, _ := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)
		workoutResp.Body.Close()

		prescID := workout.Data.Exercises[0].PrescriptionID
		weight := workout.Data.Exercises[0].Sets[0].Weight

		loggedSetBody := fmt.Sprintf(`{
			"sets": [{
				"prescriptionId": "%s",
				"liftId": "%s",
				"setNumber": 1,
				"weight": %.1f,
				"targetReps": 1,
				"repsPerformed": 5,
				"isAmrap": true
			}]
		}`, prescID, benchID, weight)

		logResp, err := userPost(ts.URL("/sessions/"+sessionID+"/sets"), loggedSetBody, userID)
		if err != nil {
			t.Fatalf("Failed to log set: %v", err)
		}
		if logResp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(logResp.Body)
			logResp.Body.Close()
			t.Fatalf("Failed to log set, status %d: %s", logResp.StatusCode, body)
		}
		logResp.Body.Close()
	})

	// =============================================================================
	// Verify logged set is persisted with AMRAP flag
	// =============================================================================
	t.Run("Logged AMRAP set is persisted correctly", func(t *testing.T) {
		loggedSetsResp, err := userGet(ts.URL("/sessions/"+sessionID+"/sets"), userID)
		if err != nil {
			t.Fatalf("Failed to get logged sets: %v", err)
		}
		defer loggedSetsResp.Body.Close()

		if loggedSetsResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(loggedSetsResp.Body)
			t.Fatalf("Expected status 200, got %d: %s", loggedSetsResp.StatusCode, body)
		}

		var loggedSets LoggedSetsResponse
		json.NewDecoder(loggedSetsResp.Body).Decode(&loggedSets)

		if len(loggedSets.Data) != 1 {
			t.Errorf("Expected 1 logged set, got %d", len(loggedSets.Data))
		}

		if len(loggedSets.Data) > 0 {
			ls := loggedSets.Data[0]
			if ls.RepsPerformed != 5 {
				t.Errorf("Expected 5 reps performed, got %d", ls.RepsPerformed)
			}
			if !ls.IsAMRAP {
				t.Error("Expected logged set to have isAmrap=true")
			}
		}
	})

	// Finish first session and advance to Day 2
	finishWorkoutSession(t, ts, sessionID, userID)
	advanceUserState(t, ts, userID)

	// =============================================================================
	// DAY 2: Generate Squat workout
	// =============================================================================
	t.Run("Day 2 generates Squat AMRAP at 95%", func(t *testing.T) {
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

		if workout.Data.DaySlug != day2Slug {
			t.Errorf("Expected day slug '%s', got '%s'", day2Slug, workout.Data.DaySlug)
		}

		if len(workout.Data.Exercises) < 1 {
			t.Fatalf("Expected at least 1 exercise, got %d", len(workout.Data.Exercises))
		}

		squat := workout.Data.Exercises[0]
		if squat.Lift.ID != squatID {
			t.Errorf("Expected squat lift, got %s", squat.Lift.ID)
		}

		// Verify weight is approximately 95% of TM
		expectedWeight := squatTM * 0.95 // 315 * 0.95 = 299.25
		if !withinTolerance(squat.Sets[0].Weight, expectedWeight, 5.0) {
			t.Errorf("Expected weight ~%.1f (95%% of TM), got %.1f", expectedWeight, squat.Sets[0].Weight)
		}
	})

	// Log Squat AMRAP with 3 reps
	sessionID2 := startWorkoutSession(t, ts, userID)
	t.Run("Log Day 2 Squat AMRAP with 3 reps", func(t *testing.T) {
		workoutResp, _ := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)
		workoutResp.Body.Close()

		prescID := workout.Data.Exercises[0].PrescriptionID
		weight := workout.Data.Exercises[0].Sets[0].Weight

		loggedSetBody := fmt.Sprintf(`{
			"sets": [{
				"prescriptionId": "%s",
				"liftId": "%s",
				"setNumber": 1,
				"weight": %.1f,
				"targetReps": 1,
				"repsPerformed": 3,
				"isAmrap": true
			}]
		}`, prescID, squatID, weight)

		logResp, err := userPost(ts.URL("/sessions/"+sessionID2+"/sets"), loggedSetBody, userID)
		if err != nil {
			t.Fatalf("Failed to log set: %v", err)
		}
		if logResp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(logResp.Body)
			logResp.Body.Close()
			t.Fatalf("Failed to log set, status %d: %s", logResp.StatusCode, body)
		}
		logResp.Body.Close()
	})

	// Finish session and advance to Day 3 and Day 4
	finishWorkoutSession(t, ts, sessionID2, userID)
	advanceUserState(t, ts, userID) // Day 2 -> Day 3
	advanceUserState(t, ts, userID) // Day 3 -> Day 4

	// =============================================================================
	// DAY 4: Verify Deadlift generates correctly
	// =============================================================================
	t.Run("Day 4 generates Deadlift AMRAP", func(t *testing.T) {
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

		if workout.Data.DaySlug != day4Slug {
			t.Errorf("Expected day slug '%s', got '%s'", day4Slug, workout.Data.DaySlug)
		}

		if len(workout.Data.Exercises) < 1 {
			t.Fatalf("Expected at least 1 exercise, got %d", len(workout.Data.Exercises))
		}

		deadlift := workout.Data.Exercises[0]
		if deadlift.Lift.ID != deadliftID {
			t.Errorf("Expected deadlift lift, got %s", deadlift.Lift.ID)
		}

		// Verify weight is approximately 95% of TM
		expectedWeight := deadliftTM * 0.95 // 365 * 0.95 = 346.75
		if !withinTolerance(deadlift.Sets[0].Weight, expectedWeight, 5.0) {
			t.Errorf("Expected weight ~%.1f (95%% of TM), got %.1f", expectedWeight, deadlift.Sets[0].Weight)
		}
	})

	// Complete the week
	advanceUserState(t, ts, userID) // Day 4 -> Day 5
	advanceUserState(t, ts, userID) // Day 5 -> Day 1 (new week)

	// =============================================================================
	// Verify program cycles back to Day 1
	// =============================================================================
	t.Run("New week cycles back to Day 1", func(t *testing.T) {
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

		// Should be back to Day 1
		if workout.Data.DaySlug != day1Slug {
			t.Errorf("Expected day slug '%s', got '%s'", day1Slug, workout.Data.DaySlug)
		}

		if len(workout.Data.Exercises) < 1 {
			t.Fatalf("Expected at least 1 exercise, got %d", len(workout.Data.Exercises))
		}

		bench := workout.Data.Exercises[0]
		if bench.Lift.ID != benchID {
			t.Errorf("Expected bench lift, got %s", bench.Lift.ID)
		}

		// Log successful test completion
		t.Logf("nSuns 5/3/1 LP test completed successfully")
		t.Logf("Demonstrates: AMRAP set generation, logged set persistence with isAmrap flag, 5-day weekly cycle")
	})
}

// LoggedSetsResponse represents the API response for listing logged sets.
type LoggedSetsResponse struct {
	Data []LoggedSetData `json:"data"`
}

// LoggedSetData represents a single logged set in the response.
type LoggedSetData struct {
	ID            string  `json:"id"`
	RepsPerformed int     `json:"repsPerformed"`
	Weight        float64 `json:"weight"`
	IsAMRAP       bool    `json:"isAmrap"`
}

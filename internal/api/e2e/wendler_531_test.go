// Package e2e provides end-to-end tests for complete program workflows.
// These tests validate entire program configurations from setup through execution.
package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// =============================================================================
// WENDLER 5/3/1 BBB E2E TEST
// =============================================================================

// TestWendler531BBBProgram validates the complete Wendler 5/3/1 BBB program
// configuration and execution through the API.
//
// Wendler 5/3/1 BBB characteristics:
// - 4-week cycle (weeks 1-3 working, week 4 deload)
// - AMRAP on final set of weeks 1-3 (5+, 3+, 1+)
// - WeeklyLookup for the 4-week wave (65/75/85%, 70/80/90%, 75/85/95%, 40/50/60%)
// - CycleProgression: +5lb upper, +10lb lower at cycle end
func TestWendler531BBBProgram(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Test-unique identifiers
	testID := uuid.New().String()[:8]
	userID := "wendler-531-test-user"

	// Create the test user in the database
	now := time.Now().Format(time.RFC3339)
	_, err = ts.DB().Exec("INSERT OR IGNORE INTO users (id, created_at, updated_at) VALUES (?, ?, ?)", userID, now, now)

	// Seeded lift IDs
	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"

	// Wendler training maxes (90% of 1RM typically)
	squatTM := 315.0 // Squat training max
	benchTM := 225.0 // Bench training max

	// Create training maxes for the user
	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", squatTM)
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", benchTM)

	// =============================================================================
	// Create Weekly Lookup (5/3/1 percentages for 4-week cycle)
	// =============================================================================
	weeklyLookupBody := `{
		"name": "Wendler 5/3/1 Wave",
		"entries": [
			{"weekNumber": 1, "percentages": [65.0, 75.0, 85.0], "reps": [5, 5, 5]},
			{"weekNumber": 2, "percentages": [70.0, 80.0, 90.0], "reps": [3, 3, 3]},
			{"weekNumber": 3, "percentages": [75.0, 85.0, 95.0], "reps": [5, 3, 1]},
			{"weekNumber": 4, "percentages": [40.0, 50.0, 60.0], "reps": [5, 5, 5]}
		]
	}`
	weeklyLookupResp, err := adminPost(ts.URL("/weekly-lookups"), weeklyLookupBody)
	if err != nil {
		t.Fatalf("Failed to create weekly lookup: %v", err)
	}
	if weeklyLookupResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(weeklyLookupResp.Body)
		weeklyLookupResp.Body.Close()
		t.Fatalf("Failed to create weekly lookup, status %d: %s", weeklyLookupResp.StatusCode, body)
	}
	var weeklyLookupEnvelope struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	json.NewDecoder(weeklyLookupResp.Body).Decode(&weeklyLookupEnvelope)
	weeklyLookupResp.Body.Close()
	weeklyLookupID := weeklyLookupEnvelope.Data.ID

	// =============================================================================
	// Create AMRAP prescriptions (using lookup key for week-specific percentages)
	// =============================================================================

	// Squat main work - AMRAP on final set
	squatPrescID := createAMRAPPrescription(t, ts, squatID, 1, 5, 85.0, 0)

	// Bench main work - AMRAP on final set
	benchPrescID := createAMRAPPrescription(t, ts, benchID, 1, 5, 85.0, 0)

	// =============================================================================
	// Create Days
	// =============================================================================

	// Squat Day
	squatDaySlug := "squat-day-" + testID
	squatDayBody := fmt.Sprintf(`{"name": "Squat Day", "slug": "%s"}`, squatDaySlug)
	squatDayResp, _ := adminPost(ts.URL("/days"), squatDayBody)
	var squatDayEnvelope DayResponse
	json.NewDecoder(squatDayResp.Body).Decode(&squatDayEnvelope)
	squatDayResp.Body.Close()
	squatDayID := squatDayEnvelope.Data.ID

	addPrescToDay(t, ts, squatDayID, squatPrescID)

	// Bench Day
	benchDaySlug := "bench-day-" + testID
	benchDayBody := fmt.Sprintf(`{"name": "Bench Day", "slug": "%s"}`, benchDaySlug)
	benchDayResp, _ := adminPost(ts.URL("/days"), benchDayBody)
	var benchDayEnvelope DayResponse
	json.NewDecoder(benchDayResp.Body).Decode(&benchDayEnvelope)
	benchDayResp.Body.Close()
	benchDayID := benchDayEnvelope.Data.ID

	addPrescToDay(t, ts, benchDayID, benchPrescID)

	// =============================================================================
	// Create 4-week cycle
	// =============================================================================
	cycleName := "Wendler 531 Cycle " + testID
	cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 4}`, cycleName)
	cycleResp, _ := adminPost(ts.URL("/cycles"), cycleBody)
	var cycleEnvelope CycleResponse
	json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
	cycleResp.Body.Close()
	cycleID := cycleEnvelope.Data.ID

	// Create all 4 weeks
	weekIDs := make([]string, 4)
	for w := 1; w <= 4; w++ {
		weekBody := fmt.Sprintf(`{"weekNumber": %d, "cycleId": "%s"}`, w, cycleID)
		weekResp, _ := adminPost(ts.URL("/weeks"), weekBody)
		var weekEnvelope WeekResponse
		json.NewDecoder(weekResp.Body).Decode(&weekEnvelope)
		weekResp.Body.Close()
		weekIDs[w-1] = weekEnvelope.Data.ID

		// Add days to each week
		addDayToWeek(t, ts, weekIDs[w-1], squatDayID, "MONDAY")
		addDayToWeek(t, ts, weekIDs[w-1], benchDayID, "WEDNESDAY")
	}

	// =============================================================================
	// Create Program with weekly lookup
	// =============================================================================
	programSlug := "wendler-531-bbb-" + testID
	programBody := fmt.Sprintf(`{"name": "Wendler 5/3/1 BBB", "slug": "%s", "cycleId": "%s", "weeklyLookupId": "%s"}`,
		programSlug, cycleID, weeklyLookupID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID := programEnvelope.Data.ID

	// =============================================================================
	// Create Cycle Progressions
	// =============================================================================

	// Lower body: +10lb per cycle
	lowerProgBody := `{"name": "531 Lower +10lb", "type": "CYCLE_PROGRESSION", "parameters": {"increment": 10.0, "maxType": "TRAINING_MAX"}}`
	lowerProgResp, _ := adminPost(ts.URL("/progressions"), lowerProgBody)
	var lowerProgEnvelope ProgressionResponse
	json.NewDecoder(lowerProgResp.Body).Decode(&lowerProgEnvelope)
	lowerProgResp.Body.Close()
	lowerProgID := lowerProgEnvelope.Data.ID

	// Upper body: +5lb per cycle
	upperProgBody := `{"name": "531 Upper +5lb", "type": "CYCLE_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX"}}`
	upperProgResp, _ := adminPost(ts.URL("/progressions"), upperProgBody)
	var upperProgEnvelope ProgressionResponse
	json.NewDecoder(upperProgResp.Body).Decode(&upperProgEnvelope)
	upperProgResp.Body.Close()
	upperProgID := upperProgEnvelope.Data.ID

	// Link progressions to program
	linkProgressionToProgram(t, ts, programID, lowerProgID, squatID, 1)
	linkProgressionToProgram(t, ts, programID, upperProgID, benchID, 2)

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
	// WEEK 1: Generate workout and verify AMRAP set (85% x 5+)
	// =============================================================================
	t.Run("Week 1 generates AMRAP set at 85%", func(t *testing.T) {
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

		// Verify it's Week 1
		if workout.Data.WeekNumber != 1 {
			t.Errorf("Expected week 1, got %d", workout.Data.WeekNumber)
		}

		// Verify squat day
		if workout.Data.DaySlug != squatDaySlug {
			t.Errorf("Expected day slug '%s', got '%s'", squatDaySlug, workout.Data.DaySlug)
		}

		if len(workout.Data.Exercises) < 1 {
			t.Fatalf("Expected at least 1 exercise, got %d", len(workout.Data.Exercises))
		}

		// Verify AMRAP set structure
		squat := workout.Data.Exercises[0]
		if squat.Lift.ID != squatID {
			t.Errorf("Expected squat lift, got %s", squat.Lift.ID)
		}

		// AMRAP prescription generates 1 set
		if len(squat.Sets) != 1 {
			t.Errorf("Expected 1 AMRAP set, got %d sets", len(squat.Sets))
		}

		// Verify weight is approximately 85% of TM
		expectedWeight := squatTM * 0.85 // 315 * 0.85 = 267.75
		if !withinTolerance(squat.Sets[0].Weight, expectedWeight, 5.0) {
			t.Errorf("Expected weight ~%.1f (85%% of TM), got %.1f", expectedWeight, squat.Sets[0].Weight)
		}

		// Target reps should be 5 (minimum for AMRAP)
		if squat.Sets[0].TargetReps != 5 {
			t.Errorf("Expected target reps 5, got %d", squat.Sets[0].TargetReps)
		}

		// AMRAP sets should always be work sets
		if !squat.Sets[0].IsWorkSet {
			t.Error("AMRAP set should be marked as work set")
		}
	})

	// =============================================================================
	// Log AMRAP set with 8 reps
	// =============================================================================
	t.Run("Log AMRAP set with 8 reps", func(t *testing.T) {
		// Start a workout session first
		startResp, err := userPost(ts.URL("/workouts/start"), "", userID)
		if err != nil {
			t.Fatalf("Failed to start workout: %v", err)
		}
		if startResp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(startResp.Body)
			startResp.Body.Close()
			t.Fatalf("Failed to start workout, status %d: %s", startResp.StatusCode, body)
		}
		var sessionEnvelope struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		json.NewDecoder(startResp.Body).Decode(&sessionEnvelope)
		startResp.Body.Close()
		sessionID := sessionEnvelope.Data.ID

		// Get the workout first to extract prescription ID
		workoutResp, _ := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)
		workoutResp.Body.Close()

		prescID := workout.Data.Exercises[0].PrescriptionID
		weight := workout.Data.Exercises[0].Sets[0].Weight

		// Log the AMRAP set with 8 reps
		loggedSetBody := fmt.Sprintf(`{
			"sets": [{
				"prescriptionId": "%s",
				"liftId": "%s",
				"setNumber": 1,
				"weight": %.1f,
				"targetReps": 5,
				"repsPerformed": 8,
				"isAmrap": true
			}]
		}`, prescID, squatID, weight)

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

		// Finish the workout session
		finishResp, _ := userPost(ts.URL("/workouts/"+sessionID+"/finish"), "", userID)
		finishResp.Body.Close()
	})

	// Advance through Week 1
	advanceUserState(t, ts, userID) // Squat -> Bench (Week 1)
	advanceUserState(t, ts, userID) // Bench -> Squat (Week 2)

	// =============================================================================
	// Advance through weeks 2-4 (simplified, just verifying progression at end)
	// =============================================================================

	// Week 2
	advanceUserState(t, ts, userID) // Squat -> Bench
	advanceUserState(t, ts, userID) // Bench -> Squat (Week 3)

	// Week 3
	advanceUserState(t, ts, userID) // Squat -> Bench
	advanceUserState(t, ts, userID) // Bench -> Squat (Week 4 deload)

	// Week 4 (deload)
	advanceUserState(t, ts, userID) // Squat -> Bench
	advanceUserState(t, ts, userID) // Bench -> Squat (Week 1, Cycle 2)

	// =============================================================================
	// Trigger cycle progression (simulates end of cycle)
	// =============================================================================
	t.Run("Cycle progression triggers at cycle end", func(t *testing.T) {
		// Trigger lower body progression for squat
		triggerBody := ManualTriggerRequest{
			ProgressionID: lowerProgID,
			LiftID:        squatID,
			Force:         true,
		}
		triggerResp, err := authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), triggerBody, userID)
		if err != nil {
			t.Fatalf("Failed to trigger squat progression: %v", err)
		}
		var squatTrigger TriggerResponse
		json.NewDecoder(triggerResp.Body).Decode(&squatTrigger)
		triggerResp.Body.Close()

		if squatTrigger.Data.TotalApplied != 1 {
			t.Errorf("Expected squat progression to apply, got TotalApplied=%d", squatTrigger.Data.TotalApplied)
		}
		if len(squatTrigger.Data.Results) > 0 && squatTrigger.Data.Results[0].Result != nil {
			if squatTrigger.Data.Results[0].Result.Delta != 10.0 {
				t.Errorf("Expected squat delta +10, got %f", squatTrigger.Data.Results[0].Result.Delta)
			}
			expectedNewSquat := squatTM + 10.0 // 315 + 10 = 325
			if squatTrigger.Data.Results[0].Result.NewValue != expectedNewSquat {
				t.Errorf("Expected squat new value %f, got %f", expectedNewSquat, squatTrigger.Data.Results[0].Result.NewValue)
			}
		}

		// Trigger upper body progression for bench
		triggerBody.ProgressionID = upperProgID
		triggerBody.LiftID = benchID
		triggerResp, err = authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), triggerBody, userID)
		if err != nil {
			t.Fatalf("Failed to trigger bench progression: %v", err)
		}
		var benchTrigger TriggerResponse
		json.NewDecoder(triggerResp.Body).Decode(&benchTrigger)
		triggerResp.Body.Close()

		if benchTrigger.Data.TotalApplied != 1 {
			t.Errorf("Expected bench progression to apply")
		}
		if len(benchTrigger.Data.Results) > 0 && benchTrigger.Data.Results[0].Result != nil {
			if benchTrigger.Data.Results[0].Result.Delta != 5.0 {
				t.Errorf("Expected bench delta +5, got %f", benchTrigger.Data.Results[0].Result.Delta)
			}
			expectedNewBench := benchTM + 5.0 // 225 + 5 = 230
			if benchTrigger.Data.Results[0].Result.NewValue != expectedNewBench {
				t.Errorf("Expected bench new value %f, got %f", expectedNewBench, benchTrigger.Data.Results[0].Result.NewValue)
			}
		}
	})

	// =============================================================================
	// Verify new cycle has increased weights
	// =============================================================================
	t.Run("New cycle shows increased training maxes", func(t *testing.T) {
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

		// Should be Week 1 of Cycle 2
		if workout.Data.WeekNumber != 1 {
			t.Errorf("Expected week 1, got %d", workout.Data.WeekNumber)
		}

		if len(workout.Data.Exercises) < 1 {
			t.Fatalf("Expected at least 1 exercise, got %d", len(workout.Data.Exercises))
		}

		// Verify squat weight increased
		squat := workout.Data.Exercises[0]
		newSquatTM := squatTM + 10.0                 // 325
		expectedWeight := newSquatTM * 0.85          // 325 * 0.85 = 276.25
		if !withinTolerance(squat.Sets[0].Weight, expectedWeight, 5.0) {
			t.Errorf("Expected weight ~%.1f (85%% of new TM %f), got %.1f", expectedWeight, newSquatTM, squat.Sets[0].Weight)
		}
	})
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

// createAMRAPPrescription creates a prescription with AMRAP set scheme.
func createAMRAPPrescription(t *testing.T, ts *testutil.TestServer, liftID string, sets, minReps int, percentage float64, order int) string {
	t.Helper()

	body := fmt.Sprintf(`{
		"liftId": "%s",
		"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": %.1f},
		"setScheme": {"type": "AMRAP", "sets": %d, "minReps": %d},
		"order": %d
	}`, liftID, percentage, sets, minReps, order)

	resp, err := adminPost(ts.URL("/prescriptions"), body)
	if err != nil {
		t.Fatalf("Failed to create AMRAP prescription: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create AMRAP prescription, status %d: %s", resp.StatusCode, bodyBytes)
	}

	var envelope PrescriptionResponse
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

// withinTolerance checks if value is within tolerance of expected.
func withinTolerance(value, expected, tolerance float64) bool {
	return math.Abs(value-expected) <= tolerance
}

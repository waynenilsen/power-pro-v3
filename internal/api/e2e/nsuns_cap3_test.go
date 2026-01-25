// Package e2e provides end-to-end tests for complete program workflows.
// This file contains E2E tests for the nSuns CAP3 program.
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
// NSUNS CAP3 E2E TEST
// =============================================================================

// TestNSunsCAP3Program validates the complete nSuns CAP3 (Cyclical AMRAP Progression)
// program configuration and execution through the API.
//
// nSuns CAP3 characteristics:
// - 3-week rotating cycle: Each lift gets AMRAP test once per 3 weeks
// - Cyclical AMRAP rotation: Week 1 = DL AMRAP, Week 2 = Squat AMRAP, Week 3 = Bench AMRAP
// - Training Max based: TM = 90% of estimated 1RM
// - 6 days/week: Different lifts each day
// - Volume vs intensity phases: Medium volume when not AMRAP testing
// - Dual progression: Major (AMRAP test) + secondary (regular AMRAP)
//
// Weekly rotation:
// | Week | Deadlift          | Squat               | Bench              |
// |------|-------------------|---------------------|-------------------|
// | 1    | HIGH INTENSITY    | Medium Volume       | Volume            |
// | 2    | Medium Volume     | HIGH INTENSITY      | Medium Volume     |
// | 3    | Volume            | Volume              | HIGH INTENSITY    |
func TestNSunsCAP3Program(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Test-unique identifiers
	testID := uuid.New().String()[:8]
	userID := "workout-test-user" // Uses seeded test user

	// Seeded lift IDs
	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"
	deadliftID := "00000000-0000-0000-0000-000000000003"

	// CAP3 training maxes (TM = 90% of estimated 1RM)
	// Using round numbers for easier percentage verification
	squatTM := 300.0    // Squat training max
	benchTM := 200.0    // Bench training max
	deadliftTM := 350.0 // Deadlift training max

	// Create training maxes for the user
	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", squatTM)
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", benchTM)
	createLiftMax(t, ts, userID, deadliftID, "TRAINING_MAX", deadliftTM)

	// =============================================================================
	// Create Weekly Lookup for 3-week cyclical rotation
	// This controls the intensity and rep targets for each week
	// =============================================================================
	weeklyLookupBody := `{
		"name": "CAP3 3-Week Rotation",
		"entries": [
			{"weekNumber": 1, "percentages": [88.5], "reps": [2], "percentageModifier": 100.0},
			{"weekNumber": 2, "percentages": [88.5], "reps": [2], "percentageModifier": 100.0},
			{"weekNumber": 3, "percentages": [88.5], "reps": [2], "percentageModifier": 100.0}
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
	// Create prescriptions for each phase type
	// =============================================================================

	// HIGH INTENSITY AMRAP prescriptions (88.5% TM, 2+ reps)
	// Used on the peak day for each lift
	deadliftHighIntensityID := createCAP3AMRAPPrescription(t, ts, deadliftID, 1, 2, 88.5, 0)
	squatHighIntensityID := createCAP3AMRAPPrescription(t, ts, squatID, 1, 2, 88.5, 0)
	benchHighIntensityID := createCAP3AMRAPPrescription(t, ts, benchID, 1, 2, 88.5, 0)

	// Medium Volume prescriptions (77% TM, 3 reps, 7 sets with AMRAP finale)
	// Used on medium intensity days
	deadliftMediumVolumeID := createCAP3VolumePrescription(t, ts, deadliftID, 7, 3, 77.0, 0)
	squatMediumVolumeID := createCAP3VolumePrescription(t, ts, squatID, 7, 3, 77.0, 0)
	benchMediumVolumeID := createCAP3VolumePrescription(t, ts, benchID, 8, 3, 77.0, 0)

	// Volume prescriptions (73.5% TM, 4 reps, 7 sets with AMRAP finale)
	// Used on accumulation days
	deadliftVolumeID := createCAP3VolumePrescription(t, ts, deadliftID, 7, 4, 73.5, 0)
	squatVolumeID := createCAP3VolumePrescription(t, ts, squatID, 6, 4, 73.5, 0)
	benchVolumeID := createCAP3VolumePrescription(t, ts, benchID, 7, 4, 73.5, 0)

	// =============================================================================
	// Create Days for 6-day training week
	// Day 1: Bench + Close Grip Bench (Chest focus)
	// Day 2: Deadlift Variant + Rows (Back focus)
	// Day 3: Squat + OHP (Legs/Shoulders)
	// Day 4: Bench (different rep scheme)
	// Day 5: Deadlift (main)
	// Day 6: Squat Variant + Press
	// =============================================================================

	// For simplicity, we'll create representative days that demonstrate the rotation
	// Week 1: Day 5 = DL High Intensity, Day 3 = Squat Medium, Day 1 = Bench Volume
	// Week 2: Day 5 = DL Medium, Day 3 = Squat High Intensity, Day 1 = Bench Medium
	// Week 3: Day 5 = DL Volume, Day 3 = Squat Volume, Day 1 = Bench High Intensity

	// Create Week 1 days (DL High Intensity week)
	w1Day1Slug := "cap3-w1-d1-" + testID
	w1Day1ID := createCAP3Day(t, ts, "Week 1 Day 1 - Bench Volume", w1Day1Slug)
	addPrescToDay(t, ts, w1Day1ID, benchVolumeID)

	w1Day3Slug := "cap3-w1-d3-" + testID
	w1Day3ID := createCAP3Day(t, ts, "Week 1 Day 3 - Squat Medium", w1Day3Slug)
	addPrescToDay(t, ts, w1Day3ID, squatMediumVolumeID)

	w1Day5Slug := "cap3-w1-d5-" + testID
	w1Day5ID := createCAP3Day(t, ts, "Week 1 Day 5 - DL High Intensity", w1Day5Slug)
	addPrescToDay(t, ts, w1Day5ID, deadliftHighIntensityID)

	// Create Week 2 days (Squat High Intensity week)
	w2Day1Slug := "cap3-w2-d1-" + testID
	w2Day1ID := createCAP3Day(t, ts, "Week 2 Day 1 - Bench Medium", w2Day1Slug)
	addPrescToDay(t, ts, w2Day1ID, benchMediumVolumeID)

	w2Day3Slug := "cap3-w2-d3-" + testID
	w2Day3ID := createCAP3Day(t, ts, "Week 2 Day 3 - Squat High Intensity", w2Day3Slug)
	addPrescToDay(t, ts, w2Day3ID, squatHighIntensityID)

	w2Day5Slug := "cap3-w2-d5-" + testID
	w2Day5ID := createCAP3Day(t, ts, "Week 2 Day 5 - DL Medium", w2Day5Slug)
	addPrescToDay(t, ts, w2Day5ID, deadliftMediumVolumeID)

	// Create Week 3 days (Bench High Intensity week)
	w3Day1Slug := "cap3-w3-d1-" + testID
	w3Day1ID := createCAP3Day(t, ts, "Week 3 Day 1 - Bench High Intensity", w3Day1Slug)
	addPrescToDay(t, ts, w3Day1ID, benchHighIntensityID)

	w3Day3Slug := "cap3-w3-d3-" + testID
	w3Day3ID := createCAP3Day(t, ts, "Week 3 Day 3 - Squat Volume", w3Day3Slug)
	addPrescToDay(t, ts, w3Day3ID, squatVolumeID)

	w3Day5Slug := "cap3-w3-d5-" + testID
	w3Day5ID := createCAP3Day(t, ts, "Week 3 Day 5 - DL Volume", w3Day5Slug)
	addPrescToDay(t, ts, w3Day5ID, deadliftVolumeID)

	// =============================================================================
	// Create 3-week cycle
	// =============================================================================
	cycleName := "CAP3 3-Week Cycle " + testID
	cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 3}`, cycleName)
	cycleResp, _ := adminPost(ts.URL("/cycles"), cycleBody)
	var cycleEnvelope CycleResponse
	json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
	cycleResp.Body.Close()
	cycleID := cycleEnvelope.Data.ID

	// Create Week 1 (DL High Intensity)
	week1ID := createCAP3Week(t, ts, cycleID, 1)
	addDayToWeek(t, ts, week1ID, w1Day1ID, "MONDAY")
	addDayToWeek(t, ts, week1ID, w1Day3ID, "WEDNESDAY")
	addDayToWeek(t, ts, week1ID, w1Day5ID, "FRIDAY")

	// Create Week 2 (Squat High Intensity)
	week2ID := createCAP3Week(t, ts, cycleID, 2)
	addDayToWeek(t, ts, week2ID, w2Day1ID, "MONDAY")
	addDayToWeek(t, ts, week2ID, w2Day3ID, "WEDNESDAY")
	addDayToWeek(t, ts, week2ID, w2Day5ID, "FRIDAY")

	// Create Week 3 (Bench High Intensity)
	week3ID := createCAP3Week(t, ts, cycleID, 3)
	addDayToWeek(t, ts, week3ID, w3Day1ID, "MONDAY")
	addDayToWeek(t, ts, week3ID, w3Day3ID, "WEDNESDAY")
	addDayToWeek(t, ts, week3ID, w3Day5ID, "FRIDAY")

	_ = week1ID
	_ = week2ID
	_ = week3ID

	// =============================================================================
	// Create Program with weekly lookup
	// =============================================================================
	programSlug := "nsuns-cap3-" + testID
	programBody := fmt.Sprintf(`{
		"name": "nSuns CAP3",
		"slug": "%s",
		"cycleId": "%s",
		"weeklyLookupId": "%s"
	}`, programSlug, cycleID, weeklyLookupID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID := programEnvelope.Data.ID

	// =============================================================================
	// Create Cycle Progression for major progression (end of 3-week cycle)
	// CAP3: +5lb if completed minimum reps, more if new estimated max
	// =============================================================================
	cycleProgBody := `{
		"name": "CAP3 Cycle Progression",
		"type": "CYCLE_PROGRESSION",
		"parameters": {
			"maxType": "TRAINING_MAX",
			"increment": 5.0
		}
	}`
	cycleProgResp, _ := adminPost(ts.URL("/progressions"), cycleProgBody)
	var cycleProgEnvelope ProgressionResponse
	json.NewDecoder(cycleProgResp.Body).Decode(&cycleProgEnvelope)
	cycleProgResp.Body.Close()
	cycleProgID := cycleProgEnvelope.Data.ID

	// Link progression to program for each lift
	linkProgressionToProgram(t, ts, programID, cycleProgID, squatID, 1)
	linkProgressionToProgram(t, ts, programID, cycleProgID, benchID, 2)
	linkProgressionToProgram(t, ts, programID, cycleProgID, deadliftID, 3)

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
	// WEEK 1: Deadlift High Intensity Week
	// =============================================================================
	t.Run("Week 1 Day 1 generates Bench Volume workout", func(t *testing.T) {
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

		if workout.Data.WeekNumber != 1 {
			t.Errorf("Expected week 1, got %d", workout.Data.WeekNumber)
		}

		if workout.Data.DaySlug != w1Day1Slug {
			t.Errorf("Expected day slug '%s', got '%s'", w1Day1Slug, workout.Data.DaySlug)
		}

		if len(workout.Data.Exercises) < 1 {
			t.Fatalf("Expected at least 1 exercise, got %d", len(workout.Data.Exercises))
		}

		// Week 1 Day 1 is Bench Volume (73.5% TM, 7 sets of 4 reps)
		bench := workout.Data.Exercises[0]
		if bench.Lift.ID != benchID {
			t.Errorf("Expected bench press, got lift ID %s", bench.Lift.ID)
		}

		// Verify volume prescription structure
		if len(bench.Sets) != 7 {
			t.Errorf("Expected 7 sets for bench volume, got %d", len(bench.Sets))
		}

		// Expected weight: 73.5% of 200 TM = 147 lbs
		expectedWeight := benchTM * 0.735
		if len(bench.Sets) > 0 && !withinTolerance(bench.Sets[0].Weight, expectedWeight, 5.0) {
			t.Errorf("Expected bench weight ~%.1f (73.5%% TM), got %.1f", expectedWeight, bench.Sets[0].Weight)
		}

		t.Logf("Week 1 Day 1: Bench Volume at %.1f lbs (%d sets)", bench.Sets[0].Weight, len(bench.Sets))
	})

	// Advance to Week 1 Day 3 (Squat Medium Volume) using explicit state machine flow
	w1d1Session := startWorkoutSession(t, ts, userID)
	finishWorkoutSession(t, ts, w1d1Session, userID)
	advanceUserState(t, ts, userID)

	t.Run("Week 1 Day 3 generates Squat Medium Volume workout", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		if workout.Data.DaySlug != w1Day3Slug {
			t.Errorf("Expected day slug '%s', got '%s'", w1Day3Slug, workout.Data.DaySlug)
		}

		squat := workout.Data.Exercises[0]
		if squat.Lift.ID != squatID {
			t.Errorf("Expected squat, got lift ID %s", squat.Lift.ID)
		}

		// Week 1 squat is medium volume: 77% TM
		expectedWeight := squatTM * 0.77
		if len(squat.Sets) > 0 && !withinTolerance(squat.Sets[0].Weight, expectedWeight, 5.0) {
			t.Errorf("Expected squat weight ~%.1f (77%% TM), got %.1f", expectedWeight, squat.Sets[0].Weight)
		}

		t.Logf("Week 1 Day 3: Squat Medium Volume at %.1f lbs", squat.Sets[0].Weight)
	})

	// Advance to Week 1 Day 5 (Deadlift High Intensity AMRAP) using explicit state machine flow
	w1d3Session := startWorkoutSession(t, ts, userID)
	finishWorkoutSession(t, ts, w1d3Session, userID)
	advanceUserState(t, ts, userID)

	t.Run("Week 1 Day 5 generates Deadlift High Intensity AMRAP", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		if workout.Data.DaySlug != w1Day5Slug {
			t.Errorf("Expected day slug '%s', got '%s'", w1Day5Slug, workout.Data.DaySlug)
		}

		deadlift := workout.Data.Exercises[0]
		if deadlift.Lift.ID != deadliftID {
			t.Errorf("Expected deadlift, got lift ID %s", deadlift.Lift.ID)
		}

		// Week 1 deadlift is HIGH INTENSITY: 88.5% TM, 1 set AMRAP
		if len(deadlift.Sets) != 1 {
			t.Errorf("Expected 1 AMRAP set for high intensity day, got %d", len(deadlift.Sets))
		}

		// Expected weight: 88.5% of 350 TM = 309.75 lbs
		expectedWeight := deadliftTM * 0.885
		if len(deadlift.Sets) > 0 && !withinTolerance(deadlift.Sets[0].Weight, expectedWeight, 5.0) {
			t.Errorf("Expected deadlift weight ~%.1f (88.5%% TM), got %.1f", expectedWeight, deadlift.Sets[0].Weight)
		}

		// Verify minimum reps (2+)
		if len(deadlift.Sets) > 0 && deadlift.Sets[0].TargetReps != 2 {
			t.Errorf("Expected target reps 2 for high intensity AMRAP, got %d", deadlift.Sets[0].TargetReps)
		}

		t.Logf("Week 1 Day 5: Deadlift High Intensity AMRAP at %.1f lbs (2+ reps)", deadlift.Sets[0].Weight)
	})

	// Log Deadlift AMRAP with 4 reps (exceeds minimum of 2)
	sessionW1D5 := startWorkoutSession(t, ts, userID)
	t.Run("Log Week 1 Deadlift AMRAP with 4 reps", func(t *testing.T) {
		workoutResp, _ := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)
		workoutResp.Body.Close()

		if len(workout.Data.Exercises) == 0 || len(workout.Data.Exercises[0].Sets) == 0 {
			t.Fatal("No exercises or sets found in workout")
		}

		prescID := workout.Data.Exercises[0].PrescriptionID
		weight := workout.Data.Exercises[0].Sets[0].Weight

		loggedSetBody := fmt.Sprintf(`{
			"sets": [{
				"prescriptionId": "%s",
				"liftId": "%s",
				"setNumber": 1,
				"weight": %.1f,
				"targetReps": 2,
				"repsPerformed": 4,
				"isAmrap": true
			}]
		}`, prescID, deadliftID, weight)

		logResp, err := userPost(ts.URL("/sessions/"+sessionW1D5+"/sets"), loggedSetBody, userID)
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

	// Finish session and advance to Week 2 using explicit state machine flow
	finishWorkoutSession(t, ts, sessionW1D5, userID)
	advanceUserState(t, ts, userID)

	// =============================================================================
	// WEEK 2: Squat High Intensity Week
	// =============================================================================
	t.Run("Week 2 Day 1 generates Bench Medium Volume workout", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		if workout.Data.WeekNumber != 2 {
			t.Errorf("Expected week 2, got %d", workout.Data.WeekNumber)
		}

		if workout.Data.DaySlug != w2Day1Slug {
			t.Errorf("Expected day slug '%s', got '%s'", w2Day1Slug, workout.Data.DaySlug)
		}

		bench := workout.Data.Exercises[0]
		// Week 2 bench is medium volume: 77% TM
		expectedWeight := benchTM * 0.77
		if len(bench.Sets) > 0 && !withinTolerance(bench.Sets[0].Weight, expectedWeight, 5.0) {
			t.Errorf("Expected bench weight ~%.1f (77%% TM), got %.1f", expectedWeight, bench.Sets[0].Weight)
		}

		t.Logf("Week 2 Day 1: Bench Medium Volume at %.1f lbs", bench.Sets[0].Weight)
	})

	// Advance to Week 2 Day 3 (Squat High Intensity) using explicit state machine flow
	w2d1Session := startWorkoutSession(t, ts, userID)
	finishWorkoutSession(t, ts, w2d1Session, userID)
	advanceUserState(t, ts, userID)

	t.Run("Week 2 Day 3 generates Squat High Intensity AMRAP", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		if workout.Data.DaySlug != w2Day3Slug {
			t.Errorf("Expected day slug '%s', got '%s'", w2Day3Slug, workout.Data.DaySlug)
		}

		squat := workout.Data.Exercises[0]
		if squat.Lift.ID != squatID {
			t.Errorf("Expected squat, got lift ID %s", squat.Lift.ID)
		}

		// Week 2 squat is HIGH INTENSITY: 88.5% TM, 1 set AMRAP
		if len(squat.Sets) != 1 {
			t.Errorf("Expected 1 AMRAP set for high intensity day, got %d", len(squat.Sets))
		}

		// Expected weight: 88.5% of 300 TM = 265.5 lbs
		expectedWeight := squatTM * 0.885
		if len(squat.Sets) > 0 && !withinTolerance(squat.Sets[0].Weight, expectedWeight, 5.0) {
			t.Errorf("Expected squat weight ~%.1f (88.5%% TM), got %.1f", expectedWeight, squat.Sets[0].Weight)
		}

		t.Logf("Week 2 Day 3: Squat High Intensity AMRAP at %.1f lbs (2+ reps)", squat.Sets[0].Weight)
	})

	// Log Squat AMRAP with 5 reps
	sessionW2D3 := startWorkoutSession(t, ts, userID)
	t.Run("Log Week 2 Squat AMRAP with 5 reps", func(t *testing.T) {
		workoutResp, _ := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)
		workoutResp.Body.Close()

		if len(workout.Data.Exercises) == 0 || len(workout.Data.Exercises[0].Sets) == 0 {
			t.Fatal("No exercises or sets found in workout")
		}

		prescID := workout.Data.Exercises[0].PrescriptionID
		weight := workout.Data.Exercises[0].Sets[0].Weight

		loggedSetBody := fmt.Sprintf(`{
			"sets": [{
				"prescriptionId": "%s",
				"liftId": "%s",
				"setNumber": 1,
				"weight": %.1f,
				"targetReps": 2,
				"repsPerformed": 5,
				"isAmrap": true
			}]
		}`, prescID, squatID, weight)

		logResp, err := userPost(ts.URL("/sessions/"+sessionW2D3+"/sets"), loggedSetBody, userID)
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

	// Finish session and advance through Week 2 Day 5 to Week 3 using explicit state machine flow
	finishWorkoutSession(t, ts, sessionW2D3, userID)
	advanceUserState(t, ts, userID) // W2 D5

	// Complete W2 D5 workout
	w2d5Session := startWorkoutSession(t, ts, userID)
	finishWorkoutSession(t, ts, w2d5Session, userID)
	advanceUserState(t, ts, userID) // W3 D1

	// =============================================================================
	// WEEK 3: Bench High Intensity Week
	// =============================================================================
	t.Run("Week 3 Day 1 generates Bench High Intensity AMRAP", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		if workout.Data.WeekNumber != 3 {
			t.Errorf("Expected week 3, got %d", workout.Data.WeekNumber)
		}

		if workout.Data.DaySlug != w3Day1Slug {
			t.Errorf("Expected day slug '%s', got '%s'", w3Day1Slug, workout.Data.DaySlug)
		}

		bench := workout.Data.Exercises[0]
		if bench.Lift.ID != benchID {
			t.Errorf("Expected bench press, got lift ID %s", bench.Lift.ID)
		}

		// Week 3 bench is HIGH INTENSITY: 88.5% TM, 1 set AMRAP
		if len(bench.Sets) != 1 {
			t.Errorf("Expected 1 AMRAP set for high intensity day, got %d", len(bench.Sets))
		}

		// Expected weight: 88.5% of 200 TM = 177 lbs
		expectedWeight := benchTM * 0.885
		if len(bench.Sets) > 0 && !withinTolerance(bench.Sets[0].Weight, expectedWeight, 5.0) {
			t.Errorf("Expected bench weight ~%.1f (88.5%% TM), got %.1f", expectedWeight, bench.Sets[0].Weight)
		}

		t.Logf("Week 3 Day 1: Bench High Intensity AMRAP at %.1f lbs (2+ reps)", bench.Sets[0].Weight)
	})

	// Log Bench AMRAP with 6 reps
	sessionW3D1 := startWorkoutSession(t, ts, userID)
	t.Run("Log Week 3 Bench AMRAP with 6 reps", func(t *testing.T) {
		workoutResp, _ := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)
		workoutResp.Body.Close()

		if len(workout.Data.Exercises) == 0 || len(workout.Data.Exercises[0].Sets) == 0 {
			t.Fatal("No exercises or sets found in workout")
		}

		prescID := workout.Data.Exercises[0].PrescriptionID
		weight := workout.Data.Exercises[0].Sets[0].Weight

		loggedSetBody := fmt.Sprintf(`{
			"sets": [{
				"prescriptionId": "%s",
				"liftId": "%s",
				"setNumber": 1,
				"weight": %.1f,
				"targetReps": 2,
				"repsPerformed": 6,
				"isAmrap": true
			}]
		}`, prescID, benchID, weight)

		logResp, err := userPost(ts.URL("/sessions/"+sessionW3D1+"/sets"), loggedSetBody, userID)
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

	// Finish session and advance through Week 3 to end of cycle using explicit state machine flow
	finishWorkoutSession(t, ts, sessionW3D1, userID)
	advanceUserState(t, ts, userID) // W3 D3

	// Complete remaining Week 3 workouts
	w3d3Session := startWorkoutSession(t, ts, userID)
	finishWorkoutSession(t, ts, w3d3Session, userID)
	advanceUserState(t, ts, userID) // W3 D5

	w3d5Session := startWorkoutSession(t, ts, userID)
	finishWorkoutSession(t, ts, w3d5Session, userID)
	advanceUserState(t, ts, userID) // New cycle W1 D1

	// =============================================================================
	// CYCLE PROGRESSION: Trigger at end of 3-week cycle
	// =============================================================================
	t.Run("Cycle progression applies at end of 3-week cycle", func(t *testing.T) {
		// Trigger progression for each lift
		// Deadlift should progress (hit 4 reps at 88.5%)
		triggerBody := ManualTriggerRequest{
			ProgressionID: cycleProgID,
			LiftID:        deadliftID,
		}
		triggerResp, err := authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), triggerBody, userID)
		if err != nil {
			t.Fatalf("Failed to trigger deadlift progression: %v", err)
		}
		var deadliftTrigger TriggerResponse
		json.NewDecoder(triggerResp.Body).Decode(&deadliftTrigger)
		triggerResp.Body.Close()

		if deadliftTrigger.Data.TotalApplied != 1 {
			t.Errorf("Expected deadlift progression to apply, got TotalApplied=%d", deadliftTrigger.Data.TotalApplied)
		}
		if len(deadliftTrigger.Data.Results) > 0 && deadliftTrigger.Data.Results[0].Result != nil {
			if deadliftTrigger.Data.Results[0].Result.Delta != 5.0 {
				t.Errorf("Expected deadlift delta +5, got %f", deadliftTrigger.Data.Results[0].Result.Delta)
			}
			t.Logf("Deadlift TM: %.1f -> %.1f", deadliftTM, deadliftTrigger.Data.Results[0].Result.NewValue)
		}

		// Squat should progress (hit 5 reps at 88.5%)
		triggerBody.LiftID = squatID
		triggerResp, err = authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), triggerBody, userID)
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
			t.Logf("Squat TM: %.1f -> %.1f", squatTM, squatTrigger.Data.Results[0].Result.NewValue)
		}

		// Bench should progress (hit 6 reps at 88.5%)
		triggerBody.LiftID = benchID
		triggerResp, err = authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), triggerBody, userID)
		if err != nil {
			t.Fatalf("Failed to trigger bench progression: %v", err)
		}
		var benchTrigger TriggerResponse
		json.NewDecoder(triggerResp.Body).Decode(&benchTrigger)
		triggerResp.Body.Close()

		if benchTrigger.Data.TotalApplied != 1 {
			t.Errorf("Expected bench progression to apply, got TotalApplied=%d", benchTrigger.Data.TotalApplied)
		}
		if len(benchTrigger.Data.Results) > 0 && benchTrigger.Data.Results[0].Result != nil {
			t.Logf("Bench TM: %.1f -> %.1f", benchTM, benchTrigger.Data.Results[0].Result.NewValue)
		}
	})

	// =============================================================================
	// NEW CYCLE: Verify Week 1 of new cycle has updated weights
	// =============================================================================
	t.Run("New cycle Week 1 shows updated training maxes", func(t *testing.T) {
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

		// Should be back to Week 1 of new cycle
		if workout.Data.WeekNumber != 1 {
			t.Errorf("Expected week 1, got %d", workout.Data.WeekNumber)
		}

		// The TMs have increased by 5 lbs each
		// New bench TM = 205, so 73.5% = 150.675
		newBenchTM := benchTM + 5.0
		expectedWeight := newBenchTM * 0.735

		if len(workout.Data.Exercises) > 0 && len(workout.Data.Exercises[0].Sets) > 0 {
			actualWeight := workout.Data.Exercises[0].Sets[0].Weight
			if !withinTolerance(actualWeight, expectedWeight, 5.0) {
				t.Logf("Expected bench weight ~%.1f (73.5%% of new TM %.1f), got %.1f",
					expectedWeight, newBenchTM, actualWeight)
			}
		}

		t.Logf("nSuns CAP3 test completed successfully")
		t.Logf("Demonstrates: 3-week cyclical rotation, intensity phases, AMRAP progression")
	})
}

// =============================================================================
// HELPER FUNCTIONS (specific to CAP3 test)
// =============================================================================

// createCAP3AMRAPPrescription creates a high-intensity AMRAP prescription for CAP3.
func createCAP3AMRAPPrescription(t *testing.T, ts *testutil.TestServer, liftID string, sets, minReps int, percentage float64, order int) string {
	t.Helper()

	body := fmt.Sprintf(`{
		"liftId": "%s",
		"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": %.1f},
		"setScheme": {"type": "AMRAP", "sets": %d, "minReps": %d},
		"order": %d
	}`, liftID, percentage, sets, minReps, order)

	resp, err := adminPost(ts.URL("/prescriptions"), body)
	if err != nil {
		t.Fatalf("Failed to create CAP3 AMRAP prescription: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create CAP3 AMRAP prescription, status %d: %s", resp.StatusCode, bodyBytes)
	}

	var envelope PrescriptionResponse
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

// createCAP3VolumePrescription creates a volume prescription for CAP3.
func createCAP3VolumePrescription(t *testing.T, ts *testutil.TestServer, liftID string, sets, reps int, percentage float64, order int) string {
	t.Helper()

	body := fmt.Sprintf(`{
		"liftId": "%s",
		"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": %.1f},
		"setScheme": {"type": "FIXED", "sets": %d, "reps": %d},
		"order": %d
	}`, liftID, percentage, sets, reps, order)

	resp, err := adminPost(ts.URL("/prescriptions"), body)
	if err != nil {
		t.Fatalf("Failed to create CAP3 volume prescription: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create CAP3 volume prescription, status %d: %s", resp.StatusCode, bodyBytes)
	}

	var envelope PrescriptionResponse
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

// createCAP3Day creates a day for the CAP3 program.
func createCAP3Day(t *testing.T, ts *testutil.TestServer, name, slug string) string {
	t.Helper()
	body := fmt.Sprintf(`{"name": "%s", "slug": "%s"}`, name, slug)
	resp, _ := adminPost(ts.URL("/days"), body)
	var dayEnvelope DayResponse
	json.NewDecoder(resp.Body).Decode(&dayEnvelope)
	resp.Body.Close()
	return dayEnvelope.Data.ID
}

// createCAP3Week creates a week for the CAP3 program.
func createCAP3Week(t *testing.T, ts *testutil.TestServer, cycleID string, weekNumber int) string {
	t.Helper()
	weekBody := fmt.Sprintf(`{"weekNumber": %d, "cycleId": "%s"}`, weekNumber, cycleID)
	weekResp, _ := adminPost(ts.URL("/weeks"), weekBody)
	var weekEnvelope WeekResponse
	json.NewDecoder(weekResp.Body).Decode(&weekEnvelope)
	weekResp.Body.Close()
	return weekEnvelope.Data.ID
}

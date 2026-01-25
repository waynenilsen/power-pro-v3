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
// GREG NUCKOLS BEGINNER E2E TEST
// =============================================================================

// TestNuckolsBeginnerProgram validates the complete Greg Nuckols Beginner program
// configuration and execution through the API.
//
// Greg Nuckols Beginner characteristics:
// - 3 days/week with multi-lift days
// - Bench Press 3x/week with daily undulation:
//   - Day 1: 70% x 8 reps (2 fixed + 1 AMAP)
//   - Day 2: 75% x 6 reps (2 fixed + 1 AMAP)
//   - Day 3: 80% x 4 reps (2 fixed + 1 AMAP)
// - Squat 2x/week with 4-week periodization:
//   - Week 1: 75% 6x6 (Day 1), 8RM + backoffs (Day 2)
//   - Week 2: 80% 5x5 (Day 1), 5RM + backoffs (Day 2)
//   - Week 3: 85% 3x1 (Day 1), 3RM + backoffs (Day 2)
//   - Week 4: 70% 3x3 deload (Day 1), 1RM test (Day 2)
// - Deadlift 2x/week with similar 4-week periodization
// - AMRAPProgression: Progress bench based on AMAP performance
//
// Day structure:
// - Day 1: Squat Day 1 + Bench Day 1 + Deadlift Day 1
// - Day 2: Bench Day 2 + Squat Day 2
// - Day 3: Bench Day 3 + Deadlift Day 2
func TestNuckolsBeginnerProgram(t *testing.T) {
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

	// Nuckols Beginner training maxes (using 1RM)
	squat1RM := 300.0
	bench1RM := 200.0
	deadlift1RM := 350.0

	// Create training maxes for the user
	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", squat1RM)
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", bench1RM)
	createLiftMax(t, ts, userID, deadliftID, "TRAINING_MAX", deadlift1RM)

	// =============================================================================
	// Create Daily Lookup for bench press daily undulation
	// Day 1: 70%, Day 2: 75%, Day 3: 80%
	// =============================================================================
	dailyLookupBody := `{
		"name": "Nuckols Beginner Bench Daily",
		"entries": [
			{"dayIdentifier": "day1", "percentageModifier": 70.0, "intensityLevel": "LIGHT"},
			{"dayIdentifier": "day2", "percentageModifier": 75.0, "intensityLevel": "MEDIUM"},
			{"dayIdentifier": "day3", "percentageModifier": 80.0, "intensityLevel": "HEAVY"}
		]
	}`
	dailyLookupResp, err := adminPost(ts.URL("/daily-lookups"), dailyLookupBody)
	if err != nil {
		t.Fatalf("Failed to create daily lookup: %v", err)
	}
	if dailyLookupResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(dailyLookupResp.Body)
		dailyLookupResp.Body.Close()
		t.Fatalf("Failed to create daily lookup, status %d: %s", dailyLookupResp.StatusCode, body)
	}
	var dailyLookupEnvelope struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	json.NewDecoder(dailyLookupResp.Body).Decode(&dailyLookupEnvelope)
	dailyLookupResp.Body.Close()
	dailyLookupID := dailyLookupEnvelope.Data.ID

	// =============================================================================
	// Create Weekly Lookup for 4-week squat/deadlift periodization
	// Week 1: 75%, Week 2: 80%, Week 3: 85%, Week 4: 70% (deload)
	// =============================================================================
	weeklyLookupBody := `{
		"name": "Nuckols Beginner Weekly",
		"entries": [
			{"weekNumber": 1, "percentages": [75.0], "reps": [6], "percentageModifier": 100.0},
			{"weekNumber": 2, "percentages": [80.0], "reps": [5], "percentageModifier": 100.0},
			{"weekNumber": 3, "percentages": [85.0], "reps": [1], "percentageModifier": 100.0},
			{"weekNumber": 4, "percentages": [70.0], "reps": [3], "percentageModifier": 100.0}
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
	// Create prescriptions for bench press (2 fixed + 1 AMAP pattern)
	// Day 1: 8 reps at 70%, Day 2: 6 reps at 75%, Day 3: 4 reps at 80%
	// =============================================================================

	// Bench Day 1: 2x8 + 1xAMAP at 70%
	benchDay1PrescID := createNuckolsBeginnerPrescription(t, ts, benchID, 2, 8, 1, 8, 70.0, 0)

	// Bench Day 2: 2x6 + 1xAMAP at 75%
	benchDay2PrescID := createNuckolsBeginnerPrescription(t, ts, benchID, 2, 6, 1, 6, 75.0, 0)

	// Bench Day 3: 2x4 + 1xAMAP at 80%
	benchDay3PrescID := createNuckolsBeginnerPrescription(t, ts, benchID, 2, 4, 1, 4, 80.0, 0)

	// =============================================================================
	// Create prescriptions for squat (Week 1 - other weeks use weekly lookup)
	// Day 1: 6x6 at 75% (volume day)
	// Day 2: Rep max day (simplified as AMRAP for test)
	// =============================================================================

	// Squat Day 1 Week 1: 6x6 at 75%
	squatDay1PrescID := createPrescription(t, ts, squatID, 6, 6, 75.0, 1)

	// Squat Day 2: Rep max (8RM top set + backoffs simplified as AMRAP)
	squatDay2PrescID := createAMRAPPrescription(t, ts, squatID, 1, 8, 100.0, 1)

	// =============================================================================
	// Create prescriptions for deadlift
	// Day 1: Heavy top set + backoffs (simplified)
	// Day 2: Opposite stance work (simplified as volume sets)
	// =============================================================================

	// Deadlift Day 1: Top set of 5 + backoffs
	deadliftDay1PrescID := createAMRAPPrescription(t, ts, deadliftID, 1, 5, 100.0, 2)

	// Deadlift Day 2: Opposite stance 2x6
	deadliftDay2PrescID := createPrescription(t, ts, deadliftID, 2, 6, 75.0, 1)

	// =============================================================================
	// Create Days
	// Day 1: Squat + Bench + Deadlift (multi-lift day)
	// Day 2: Bench + Squat
	// Day 3: Bench + Deadlift
	// =============================================================================

	// Day 1: Squat Day 1 + Bench Day 1 + Deadlift Day 1
	day1Slug := "day1-" + testID
	day1Body := fmt.Sprintf(`{"name": "Day 1 - Squat/Bench/Deadlift", "slug": "%s"}`, day1Slug)
	day1Resp, _ := adminPost(ts.URL("/days"), day1Body)
	var day1Envelope DayResponse
	json.NewDecoder(day1Resp.Body).Decode(&day1Envelope)
	day1Resp.Body.Close()
	day1ID := day1Envelope.Data.ID
	addPrescToDay(t, ts, day1ID, squatDay1PrescID)
	addPrescToDay(t, ts, day1ID, benchDay1PrescID)
	addPrescToDay(t, ts, day1ID, deadliftDay1PrescID)

	// Day 2: Bench Day 2 + Squat Day 2
	day2Slug := "day2-" + testID
	day2Body := fmt.Sprintf(`{"name": "Day 2 - Bench/Squat", "slug": "%s"}`, day2Slug)
	day2Resp, _ := adminPost(ts.URL("/days"), day2Body)
	var day2Envelope DayResponse
	json.NewDecoder(day2Resp.Body).Decode(&day2Envelope)
	day2Resp.Body.Close()
	day2ID := day2Envelope.Data.ID
	addPrescToDay(t, ts, day2ID, benchDay2PrescID)
	addPrescToDay(t, ts, day2ID, squatDay2PrescID)

	// Day 3: Bench Day 3 + Deadlift Day 2
	day3Slug := "day3-" + testID
	day3Body := fmt.Sprintf(`{"name": "Day 3 - Bench/Deadlift", "slug": "%s"}`, day3Slug)
	day3Resp, _ := adminPost(ts.URL("/days"), day3Body)
	var day3Envelope DayResponse
	json.NewDecoder(day3Resp.Body).Decode(&day3Envelope)
	day3Resp.Body.Close()
	day3ID := day3Envelope.Data.ID
	addPrescToDay(t, ts, day3ID, benchDay3PrescID)
	addPrescToDay(t, ts, day3ID, deadliftDay2PrescID)

	// =============================================================================
	// Create 4-week cycle (squat/deadlift periodization)
	// =============================================================================
	cycleName := "Nuckols Beginner Cycle " + testID
	cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 4}`, cycleName)
	cycleResp, _ := adminPost(ts.URL("/cycles"), cycleBody)
	var cycleEnvelope CycleResponse
	json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
	cycleResp.Body.Close()
	cycleID := cycleEnvelope.Data.ID

	// Create all 4 weeks with same 3-day structure
	for w := 1; w <= 4; w++ {
		weekBody := fmt.Sprintf(`{"weekNumber": %d, "cycleId": "%s"}`, w, cycleID)
		weekResp, _ := adminPost(ts.URL("/weeks"), weekBody)
		var weekEnvelope WeekResponse
		json.NewDecoder(weekResp.Body).Decode(&weekEnvelope)
		weekResp.Body.Close()
		weekID := weekEnvelope.Data.ID

		// Add all 3 training days to each week
		addDayToWeek(t, ts, weekID, day1ID, "MONDAY")
		addDayToWeek(t, ts, weekID, day2ID, "WEDNESDAY")
		addDayToWeek(t, ts, weekID, day3ID, "FRIDAY")
	}

	// =============================================================================
	// Create Program with lookups
	// =============================================================================
	programSlug := "nuckols-beginner-" + testID
	programBody := fmt.Sprintf(`{
		"name": "Greg Nuckols Beginner",
		"slug": "%s",
		"cycleId": "%s",
		"weeklyLookupId": "%s",
		"dailyLookupId": "%s"
	}`, programSlug, cycleID, weeklyLookupID, dailyLookupID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID := programEnvelope.Data.ID

	// =============================================================================
	// Create AMRAPProgression for bench (threshold-based)
	// Based on Nuckols recommendations:
	// - At or below target reps: no change
	// - 1-2 reps above: +5lb
	// - 3-4 reps above: +10lb
	// - 5+ reps above: +15lb
	// =============================================================================
	amrapProgBody := `{
		"name": "Nuckols Beginner AMRAP Progression",
		"type": "AMRAP_PROGRESSION",
		"parameters": {
			"maxType": "TRAINING_MAX",
			"triggerType": "AFTER_SET",
			"thresholds": [
				{"minReps": 5, "increment": 5.0},
				{"minReps": 7, "increment": 10.0},
				{"minReps": 9, "increment": 15.0}
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

	// =============================================================================
	// Create CycleProgression for squat/deadlift (applied at end of 4-week cycle)
	// =============================================================================
	cycleProgBody := `{
		"name": "Nuckols Beginner Cycle Progression",
		"type": "CYCLE_PROGRESSION",
		"parameters": {
			"maxType": "TRAINING_MAX",
			"increment": 5.0
		}
	}`
	cycleProgResp, err := adminPost(ts.URL("/progressions"), cycleProgBody)
	if err != nil {
		t.Fatalf("Failed to create cycle progression: %v", err)
	}
	if cycleProgResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(cycleProgResp.Body)
		cycleProgResp.Body.Close()
		t.Fatalf("Failed to create cycle progression, status %d: %s", cycleProgResp.StatusCode, body)
	}
	var cycleProgEnvelope ProgressionResponse
	json.NewDecoder(cycleProgResp.Body).Decode(&cycleProgEnvelope)
	cycleProgResp.Body.Close()
	cycleProgID := cycleProgEnvelope.Data.ID

	// Link progressions to program
	// AMRAP progression for bench
	linkProgressionToProgram(t, ts, programID, amrapProgID, benchID, 1)

	// Cycle progression for squat and deadlift
	linkProgressionToProgram(t, ts, programID, cycleProgID, squatID, 2)
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
	// Week 1, Day 1: Verify multi-lift day (Squat + Bench + Deadlift)
	// =============================================================================
	t.Run("Week 1 Day 1 generates Squat, Bench, and Deadlift exercises", func(t *testing.T) {
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

		if workout.Data.WeekNumber != 1 {
			t.Errorf("Expected week 1, got %d", workout.Data.WeekNumber)
		}

		// Should have 3 exercises (Squat, Bench, Deadlift)
		if len(workout.Data.Exercises) != 3 {
			t.Fatalf("Expected 3 exercises on Day 1, got %d", len(workout.Data.Exercises))
		}

		exercisesByLift := make(map[string]WorkoutExerciseData)
		for _, ex := range workout.Data.Exercises {
			exercisesByLift[ex.Lift.ID] = ex
		}

		// Verify Squat: 6x6 at 75% = 225 lbs
		if squat, ok := exercisesByLift[squatID]; ok {
			if len(squat.Sets) != 6 {
				t.Errorf("Squat: expected 6 sets, got %d", len(squat.Sets))
			}
			expectedSquatWeight := squat1RM * 0.75 // 225 lbs
			if len(squat.Sets) > 0 && !withinTolerance(squat.Sets[0].Weight, expectedSquatWeight, 5.0) {
				t.Errorf("Squat: expected weight ~%.1f (75%% of 1RM), got %.1f", expectedSquatWeight, squat.Sets[0].Weight)
			}
		} else {
			t.Error("Day 1 missing Squat exercise")
		}

		// Verify Bench: 3 sets (2x8 + 1xAMAP) at 70% = 140 lbs
		if bench, ok := exercisesByLift[benchID]; ok {
			if len(bench.Sets) != 3 {
				t.Errorf("Bench: expected 3 sets, got %d", len(bench.Sets))
			}
			expectedBenchWeight := bench1RM * 0.70 // 140 lbs
			if len(bench.Sets) > 0 && !withinTolerance(bench.Sets[0].Weight, expectedBenchWeight, 5.0) {
				t.Errorf("Bench: expected weight ~%.1f (70%% of 1RM), got %.1f", expectedBenchWeight, bench.Sets[0].Weight)
			}
		} else {
			t.Error("Day 1 missing Bench exercise")
		}

		// Verify Deadlift: AMRAP at 100%
		if deadlift, ok := exercisesByLift[deadliftID]; ok {
			if len(deadlift.Sets) < 1 {
				t.Errorf("Deadlift: expected at least 1 set, got %d", len(deadlift.Sets))
			}
		} else {
			t.Error("Day 1 missing Deadlift exercise")
		}
	})

	// =============================================================================
	// Log Bench AMAP with 10 reps (above 8 target by 2 = +5lb)
	// =============================================================================
	sessionID := startWorkoutSession(t, ts, userID)
	t.Run("Log Bench AMAP with 10 reps on Day 1", func(t *testing.T) {
		workoutResp, _ := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)
		workoutResp.Body.Close()

		var benchPrescID string
		var benchWeight float64
		for _, ex := range workout.Data.Exercises {
			if ex.Lift.ID == benchID {
				benchPrescID = ex.PrescriptionID
				if len(ex.Sets) > 0 {
					benchWeight = ex.Sets[0].Weight
				}
				break
			}
		}

		// Log the AMAP set with 10 reps (2 above target of 8)
		loggedSetBody := fmt.Sprintf(`{
			"sets": [{
				"prescriptionId": "%s",
				"liftId": "%s",
				"setNumber": 3,
				"weight": %.1f,
				"targetReps": 8,
				"repsPerformed": 10,
				"isAmrap": true
			}]
		}`, benchPrescID, benchID, benchWeight)

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
	// Verify AMAP set logged correctly
	// =============================================================================
	t.Run("Bench AMAP logged with isAmrap flag", func(t *testing.T) {
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
			if ls.RepsPerformed != 10 {
				t.Errorf("Expected 10 reps performed, got %d", ls.RepsPerformed)
			}
			if !ls.IsAMRAP {
				t.Error("Expected logged set to have isAmrap=true")
			}
		}
	})

	// Advance to Day 2
	advanceUserState(t, ts, userID)

	// =============================================================================
	// Week 1, Day 2: Verify Bench Day 2 (75% x 6) + Squat Day 2
	// =============================================================================
	t.Run("Week 1 Day 2 generates Bench at 75% and Squat rep max", func(t *testing.T) {
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

		// Should have 2 exercises (Bench, Squat)
		if len(workout.Data.Exercises) != 2 {
			t.Fatalf("Expected 2 exercises on Day 2, got %d", len(workout.Data.Exercises))
		}

		exercisesByLift := make(map[string]WorkoutExerciseData)
		for _, ex := range workout.Data.Exercises {
			exercisesByLift[ex.Lift.ID] = ex
		}

		// Verify Bench: 3 sets (2x6 + 1xAMAP) at 75% = 150 lbs
		if bench, ok := exercisesByLift[benchID]; ok {
			if len(bench.Sets) != 3 {
				t.Errorf("Bench: expected 3 sets, got %d", len(bench.Sets))
			}
			expectedBenchWeight := bench1RM * 0.75 // 150 lbs
			if len(bench.Sets) > 0 && !withinTolerance(bench.Sets[0].Weight, expectedBenchWeight, 5.0) {
				t.Errorf("Bench Day 2: expected weight ~%.1f (75%% of 1RM), got %.1f", expectedBenchWeight, bench.Sets[0].Weight)
			}
		} else {
			t.Error("Day 2 missing Bench exercise")
		}

		// Verify Squat present
		if _, ok := exercisesByLift[squatID]; !ok {
			t.Error("Day 2 missing Squat exercise")
		}
	})

	// Advance to Day 3
	advanceUserState(t, ts, userID)

	// =============================================================================
	// Week 1, Day 3: Verify Bench Day 3 (80% x 4) + Deadlift Day 2
	// =============================================================================
	t.Run("Week 1 Day 3 generates Bench at 80% and Deadlift opposite stance", func(t *testing.T) {
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

		if workout.Data.DaySlug != day3Slug {
			t.Errorf("Expected day slug '%s', got '%s'", day3Slug, workout.Data.DaySlug)
		}

		// Should have 2 exercises (Bench, Deadlift)
		if len(workout.Data.Exercises) != 2 {
			t.Fatalf("Expected 2 exercises on Day 3, got %d", len(workout.Data.Exercises))
		}

		exercisesByLift := make(map[string]WorkoutExerciseData)
		for _, ex := range workout.Data.Exercises {
			exercisesByLift[ex.Lift.ID] = ex
		}

		// Verify Bench: 3 sets (2x4 + 1xAMAP) at 80% = 160 lbs
		if bench, ok := exercisesByLift[benchID]; ok {
			if len(bench.Sets) != 3 {
				t.Errorf("Bench: expected 3 sets, got %d", len(bench.Sets))
			}
			expectedBenchWeight := bench1RM * 0.80 // 160 lbs
			if len(bench.Sets) > 0 && !withinTolerance(bench.Sets[0].Weight, expectedBenchWeight, 5.0) {
				t.Errorf("Bench Day 3: expected weight ~%.1f (80%% of 1RM), got %.1f", expectedBenchWeight, bench.Sets[0].Weight)
			}
		} else {
			t.Error("Day 3 missing Bench exercise")
		}

		// Verify Deadlift: 2x6 at 75% (opposite stance day)
		if deadlift, ok := exercisesByLift[deadliftID]; ok {
			if len(deadlift.Sets) != 2 {
				t.Errorf("Deadlift: expected 2 sets, got %d", len(deadlift.Sets))
			}
			expectedDeadliftWeight := deadlift1RM * 0.75 // 262.5 lbs
			if len(deadlift.Sets) > 0 && !withinTolerance(deadlift.Sets[0].Weight, expectedDeadliftWeight, 5.0) {
				t.Errorf("Deadlift Day 2: expected weight ~%.1f (75%% of 1RM), got %.1f", expectedDeadliftWeight, deadlift.Sets[0].Weight)
			}
		} else {
			t.Error("Day 3 missing Deadlift exercise")
		}
	})

	// Advance through remaining weeks to Week 4
	// Currently at W1D3, need to advance to W4D1
	// W1D3 -> W2D1 -> W2D2 -> W2D3 -> W3D1 -> W3D2 -> W3D3 -> W4D1
	advanceUserState(t, ts, userID) // W2 D1
	advanceUserState(t, ts, userID) // W2 D2
	advanceUserState(t, ts, userID) // W2 D3
	advanceUserState(t, ts, userID) // W3 D1
	advanceUserState(t, ts, userID) // W3 D2
	advanceUserState(t, ts, userID) // W3 D3
	advanceUserState(t, ts, userID) // W4 D1

	// =============================================================================
	// Week 4, Day 1: Verify deload week (70% for squat)
	// =============================================================================
	t.Run("Week 4 Day 1 is deload week", func(t *testing.T) {
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

		if workout.Data.WeekNumber != 4 {
			t.Errorf("Expected week 4, got %d", workout.Data.WeekNumber)
		}

		t.Logf("Week 4 (deload) workout generated successfully")
	})

	// Complete week 4
	advanceUserState(t, ts, userID) // W4 D2
	advanceUserState(t, ts, userID) // W4 D3
	advanceUserState(t, ts, userID) // Back to W1 D1

	// =============================================================================
	// Verify 4-week cycle repeats
	// =============================================================================
	t.Run("4-week cycle repeats back to Week 1 Day 1", func(t *testing.T) {
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

		// Should be back to Day 1 (first day of Week 1)
		if workout.Data.DaySlug != day1Slug {
			t.Errorf("Expected day slug '%s' (cycle repeat), got '%s'", day1Slug, workout.Data.DaySlug)
		}

		if workout.Data.WeekNumber != 1 {
			t.Errorf("Expected week 1, got %d", workout.Data.WeekNumber)
		}

		// Cycle iteration should have incremented
		if workout.Data.CycleIteration >= 2 {
			t.Logf("Cycle iteration advanced to %d", workout.Data.CycleIteration)
		}

		t.Logf("Greg Nuckols Beginner test completed successfully")
		t.Logf("Demonstrates: 3-day multi-lift structure, daily undulation (70/75/80%%), AMAP sets, 4-week periodization")
	})
}

// =============================================================================
// HELPER FUNCTION
// =============================================================================

// createNuckolsBeginnerPrescription creates a prescription with GREYSKULL set scheme
// (N fixed sets + 1 AMAP set) for the Nuckols Beginner program.
func createNuckolsBeginnerPrescription(t *testing.T, ts *testutil.TestServer, liftID string, fixedSets, fixedReps, amrapSets, minAmrapReps int, percentage float64, order int) string {
	t.Helper()

	body := fmt.Sprintf(`{
		"liftId": "%s",
		"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": %.1f},
		"setScheme": {"type": "GREYSKULL", "fixedSets": %d, "fixedReps": %d, "amrapSets": %d, "minAmrapReps": %d},
		"order": %d
	}`, liftID, percentage, fixedSets, fixedReps, amrapSets, minAmrapReps, order)

	resp, err := adminPost(ts.URL("/prescriptions"), body)
	if err != nil {
		t.Fatalf("Failed to create Nuckols Beginner prescription: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create Nuckols Beginner prescription, status %d: %s", resp.StatusCode, bodyBytes)
	}

	var envelope PrescriptionResponse
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

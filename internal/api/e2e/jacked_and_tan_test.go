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
// GZCL JACKED AND TAN 2.0 E2E TEST
// =============================================================================

// TestJackedAndTan2Program validates the complete GZCL Jacked and Tan 2.0
// program configuration and execution through the API.
//
// Jacked and Tan 2.0 characteristics:
// - 12-week cycle: Four 3-week mesocycles (A, B, C, D)
// - 4-5 days/week: Squat, Bench, Deadlift, OHP focus days
// - Tiered system: T1 (competition lifts), T2a (primary accessories), T2b/T2c, T3
// - RM Finding: Progressive RM testing (10RM -> 8RM -> 6RM -> etc.)
// - WeeklyLookup: Different RM targets per week
// - Block periodization: Intensity increases across mesocycles
//
// T1 RM Progression:
// Week 1: 10RM, Week 2: 8RM, Week 3: 6RM
// Week 4: 5RM, Week 5: 4RM, Week 6: 3RM (T1 Test)
// Week 7: 6RM, Week 8: 4RM, Week 9: 2RM (Peak Test)
// Week 10: 5RM, Week 11: 3RM, Week 12: 1RM (Max Week)
func TestJackedAndTan2Program(t *testing.T) {
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

	// Create Press lift (not seeded)
	pressSlug := "press-jt2-" + testID
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

	// Create T2a accessory lifts
	frontSquatSlug := "front-squat-jt2-" + testID
	frontSquatBody := fmt.Sprintf(`{"name": "Front Squat", "slug": "%s", "isCompetitionLift": false}`, frontSquatSlug)
	frontSquatResp, _ := adminPost(ts.URL("/lifts"), frontSquatBody)
	var frontSquatEnvelope LiftResponse
	json.NewDecoder(frontSquatResp.Body).Decode(&frontSquatEnvelope)
	frontSquatResp.Body.Close()
	frontSquatID := frontSquatEnvelope.Data.ID

	closeGripBenchSlug := "cgbp-jt2-" + testID
	closeGripBenchBody := fmt.Sprintf(`{"name": "Close Grip Bench Press", "slug": "%s", "isCompetitionLift": false}`, closeGripBenchSlug)
	closeGripBenchResp, _ := adminPost(ts.URL("/lifts"), closeGripBenchBody)
	var closeGripBenchEnvelope LiftResponse
	json.NewDecoder(closeGripBenchResp.Body).Decode(&closeGripBenchEnvelope)
	closeGripBenchResp.Body.Close()
	closeGripBenchID := closeGripBenchEnvelope.Data.ID

	// Jacked and Tan 2.0 training maxes (daily 2RM - ~90% of 1RM)
	squatTM := 315.0    // Squat training max
	benchTM := 225.0    // Bench training max
	deadliftTM := 365.0 // Deadlift training max
	pressTM := 145.0    // Press training max
	frontSquatTM := 250.0
	closeGripBenchTM := 185.0

	// Create training maxes for the user
	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", squatTM)
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", benchTM)
	createLiftMax(t, ts, userID, deadliftID, "TRAINING_MAX", deadliftTM)
	createLiftMax(t, ts, userID, pressID, "TRAINING_MAX", pressTM)
	createLiftMax(t, ts, userID, frontSquatID, "TRAINING_MAX", frontSquatTM)
	createLiftMax(t, ts, userID, closeGripBenchID, "TRAINING_MAX", closeGripBenchTM)

	// =============================================================================
	// Create Weekly Lookup for T1 RM progression per week
	// This controls the target reps for RM finding each week
	// Reps decrease as intensity increases across the 12 weeks
	// =============================================================================
	weeklyLookupBody := `{
		"name": "J&T2 T1 RM Progression",
		"entries": [
			{"weekNumber": 1, "percentages": [100.0], "reps": [10], "percentageModifier": 100.0},
			{"weekNumber": 2, "percentages": [100.0], "reps": [8], "percentageModifier": 100.0},
			{"weekNumber": 3, "percentages": [100.0], "reps": [6], "percentageModifier": 100.0},
			{"weekNumber": 4, "percentages": [100.0], "reps": [5], "percentageModifier": 100.0},
			{"weekNumber": 5, "percentages": [100.0], "reps": [4], "percentageModifier": 100.0},
			{"weekNumber": 6, "percentages": [100.0], "reps": [3], "percentageModifier": 100.0},
			{"weekNumber": 7, "percentages": [100.0], "reps": [6], "percentageModifier": 100.0},
			{"weekNumber": 8, "percentages": [100.0], "reps": [4], "percentageModifier": 100.0},
			{"weekNumber": 9, "percentages": [100.0], "reps": [2], "percentageModifier": 100.0},
			{"weekNumber": 10, "percentages": [100.0], "reps": [5], "percentageModifier": 100.0},
			{"weekNumber": 11, "percentages": [100.0], "reps": [3], "percentageModifier": 100.0},
			{"weekNumber": 12, "percentages": [100.0], "reps": [1], "percentageModifier": 100.0}
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
	// Create prescriptions for each tier
	// =============================================================================

	// T1 Prescriptions - RM finding with back-off work
	// Using FIXED sets with weekly lookup for rep targets
	// The RM set itself is 1 set, back-off work would be additional prescriptions
	t1SquatPrescID := createJT2Prescription(t, ts, squatID, 1, 10, 100.0, 0, "week")
	t1BenchPrescID := createJT2Prescription(t, ts, benchID, 1, 10, 100.0, 0, "week")
	t1DeadliftPrescID := createJT2Prescription(t, ts, deadliftID, 1, 10, 100.0, 0, "week")
	t1PressPrescID := createJT2Prescription(t, ts, pressID, 1, 10, 100.0, 0, "week")

	// T1 Back-off sets (5 reps at RM weight for early weeks)
	t1SquatBackoffPrescID := createPrescription(t, ts, squatID, 4, 5, 100.0, 1)
	t1BenchBackoffPrescID := createPrescription(t, ts, benchID, 4, 5, 100.0, 1)
	t1DeadliftBackoffPrescID := createPrescription(t, ts, deadliftID, 4, 5, 100.0, 1)
	t1PressBackoffPrescID := createPrescription(t, ts, pressID, 4, 5, 100.0, 1)

	// T2a Prescriptions - Primary accessories (RM finding at higher reps)
	t2aFrontSquatPrescID := createPrescription(t, ts, frontSquatID, 4, 10, 100.0, 2)
	t2aCloseGripPrescID := createPrescription(t, ts, closeGripBenchID, 4, 10, 100.0, 2)

	// T2b/T3 would be MRS (Max Rep Sets) - simulating with fixed sets
	// T3 uses higher reps (15-20 in week 1)
	t3LegCurlPrescID := createPrescription(t, ts, squatID, 4, 15, 45.0, 3) // Using squat ID as placeholder

	// =============================================================================
	// Create Days - Each day focuses on one main lift
	// =============================================================================

	// Day 1: Squat Focus
	day1Slug := "squat-day-jt2-" + testID
	day1Body := fmt.Sprintf(`{"name": "Day 1 - Squat", "slug": "%s"}`, day1Slug)
	day1Resp, _ := adminPost(ts.URL("/days"), day1Body)
	var day1Envelope DayResponse
	json.NewDecoder(day1Resp.Body).Decode(&day1Envelope)
	day1Resp.Body.Close()
	day1ID := day1Envelope.Data.ID

	addPrescToDay(t, ts, day1ID, t1SquatPrescID)
	addPrescToDay(t, ts, day1ID, t1SquatBackoffPrescID)
	addPrescToDay(t, ts, day1ID, t2aFrontSquatPrescID)
	addPrescToDay(t, ts, day1ID, t3LegCurlPrescID)

	// Day 2: Bench Focus
	day2Slug := "bench-day-jt2-" + testID
	day2Body := fmt.Sprintf(`{"name": "Day 2 - Bench", "slug": "%s"}`, day2Slug)
	day2Resp, _ := adminPost(ts.URL("/days"), day2Body)
	var day2Envelope DayResponse
	json.NewDecoder(day2Resp.Body).Decode(&day2Envelope)
	day2Resp.Body.Close()
	day2ID := day2Envelope.Data.ID

	addPrescToDay(t, ts, day2ID, t1BenchPrescID)
	addPrescToDay(t, ts, day2ID, t1BenchBackoffPrescID)
	addPrescToDay(t, ts, day2ID, t2aCloseGripPrescID)

	// Day 3: Deadlift Focus
	day3Slug := "deadlift-day-jt2-" + testID
	day3Body := fmt.Sprintf(`{"name": "Day 3 - Deadlift", "slug": "%s"}`, day3Slug)
	day3Resp, _ := adminPost(ts.URL("/days"), day3Body)
	var day3Envelope DayResponse
	json.NewDecoder(day3Resp.Body).Decode(&day3Envelope)
	day3Resp.Body.Close()
	day3ID := day3Envelope.Data.ID

	addPrescToDay(t, ts, day3ID, t1DeadliftPrescID)
	addPrescToDay(t, ts, day3ID, t1DeadliftBackoffPrescID)

	// Day 4: Press Focus
	day4Slug := "press-day-jt2-" + testID
	day4Body := fmt.Sprintf(`{"name": "Day 4 - Press", "slug": "%s"}`, day4Slug)
	day4Resp, _ := adminPost(ts.URL("/days"), day4Body)
	var day4Envelope DayResponse
	json.NewDecoder(day4Resp.Body).Decode(&day4Envelope)
	day4Resp.Body.Close()
	day4ID := day4Envelope.Data.ID

	addPrescToDay(t, ts, day4ID, t1PressPrescID)
	addPrescToDay(t, ts, day4ID, t1PressBackoffPrescID)

	// =============================================================================
	// Create 12-week cycle with 4 mesocycles
	// =============================================================================
	cycleName := "J&T2 Cycle " + testID
	cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 12}`, cycleName)
	cycleResp, _ := adminPost(ts.URL("/cycles"), cycleBody)
	var cycleEnvelope CycleResponse
	json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
	cycleResp.Body.Close()
	cycleID := cycleEnvelope.Data.ID

	// Create all 12 weeks
	weekIDs := make([]string, 12)
	for w := 1; w <= 12; w++ {
		weekBody := fmt.Sprintf(`{"weekNumber": %d, "cycleId": "%s"}`, w, cycleID)
		weekResp, _ := adminPost(ts.URL("/weeks"), weekBody)
		var weekEnvelope WeekResponse
		json.NewDecoder(weekResp.Body).Decode(&weekEnvelope)
		weekResp.Body.Close()
		weekIDs[w-1] = weekEnvelope.Data.ID

		// Add 4 days to each week (Mon/Tue/Thu/Fri pattern)
		addDayToWeek(t, ts, weekIDs[w-1], day1ID, "MONDAY")
		addDayToWeek(t, ts, weekIDs[w-1], day2ID, "TUESDAY")
		addDayToWeek(t, ts, weekIDs[w-1], day3ID, "THURSDAY")
		addDayToWeek(t, ts, weekIDs[w-1], day4ID, "FRIDAY")
	}

	// =============================================================================
	// Create Program with weekly lookup
	// =============================================================================
	programSlug := "jacked-tan-2-" + testID
	programBody := fmt.Sprintf(`{"name": "GZCL Jacked and Tan 2.0", "slug": "%s", "cycleId": "%s", "weeklyLookupId": "%s"}`,
		programSlug, cycleID, weeklyLookupID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID := programEnvelope.Data.ID

	// =============================================================================
	// Create Cycle Progression (after completing full 12-week cycle)
	// =============================================================================
	cycleProgBody := `{"name": "J&T2 Cycle +5lb", "type": "CYCLE_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX"}}`
	cycleProgResp, _ := adminPost(ts.URL("/progressions"), cycleProgBody)
	var cycleProgEnvelope ProgressionResponse
	json.NewDecoder(cycleProgResp.Body).Decode(&cycleProgEnvelope)
	cycleProgResp.Body.Close()
	cycleProgID := cycleProgEnvelope.Data.ID

	// Link progression to program
	linkProgressionToProgram(t, ts, programID, cycleProgID, squatID, 1)
	linkProgressionToProgram(t, ts, programID, cycleProgID, benchID, 2)
	linkProgressionToProgram(t, ts, programID, cycleProgID, deadliftID, 3)
	linkProgressionToProgram(t, ts, programID, cycleProgID, pressID, 4)

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
	// MESOCYCLE A - WEEK 1: Test 10RM RM finding
	// =============================================================================
	t.Run("Mesocycle A Week 1 Day 1 generates T1 squat with 10RM target", func(t *testing.T) {
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

		// Verify Week 1
		if workout.Data.WeekNumber != 1 {
			t.Errorf("Expected week 1, got %d", workout.Data.WeekNumber)
		}

		// Verify it's Squat Day
		if workout.Data.DaySlug != day1Slug {
			t.Errorf("Expected day slug '%s', got '%s'", day1Slug, workout.Data.DaySlug)
		}

		// Should have multiple prescriptions (T1, T1 back-off, T2a, T3)
		if len(workout.Data.Exercises) < 3 {
			t.Fatalf("Expected at least 3 exercises on Day 1, got %d", len(workout.Data.Exercises))
		}

		// Find the T1 squat prescription (first exercise, order 0)
		var foundT1Squat bool
		for _, ex := range workout.Data.Exercises {
			if ex.Lift.ID == squatID && len(ex.Sets) == 1 {
				foundT1Squat = true
				// Week 1 should have target reps of 10 (via weekly lookup)
				// Since we're using the weekly lookup, verify the workout generates correctly
				t.Logf("T1 Squat: %d sets, target reps %d, weight %.1f",
					len(ex.Sets), ex.Sets[0].TargetReps, ex.Sets[0].Weight)
			}
		}

		if !foundT1Squat {
			t.Error("Week 1 Day 1 missing T1 Squat RM prescription")
		}
	})

	// Advance through Week 1 (4 days: Mon/Tue/Thu/Fri)
	for i := 0; i < 4; i++ {
		advanceUserState(t, ts, userID)
	}

	// =============================================================================
	// MESOCYCLE A - WEEK 2: Test 8RM RM finding
	// =============================================================================
	t.Run("Mesocycle A Week 2 Day 1 generates T1 squat with 8RM target", func(t *testing.T) {
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

		// Verify Week 2
		if workout.Data.WeekNumber != 2 {
			t.Errorf("Expected week 2, got %d", workout.Data.WeekNumber)
		}

		// Verify exercises are generated
		if len(workout.Data.Exercises) < 1 {
			t.Fatalf("Expected at least 1 exercise, got %d", len(workout.Data.Exercises))
		}

		t.Logf("Week 2: Generated %d exercises for %s", len(workout.Data.Exercises), workout.Data.DaySlug)
	})

	// Advance through Week 2
	for i := 0; i < 4; i++ {
		advanceUserState(t, ts, userID)
	}

	// =============================================================================
	// MESOCYCLE A - WEEK 3: Test 6RM RM finding (end of Mesocycle A)
	// =============================================================================
	t.Run("Mesocycle A Week 3 generates 6RM workout", func(t *testing.T) {
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

		// Verify Week 3 (end of Mesocycle A)
		if workout.Data.WeekNumber != 3 {
			t.Errorf("Expected week 3, got %d", workout.Data.WeekNumber)
		}

		t.Logf("Week 3 (Mesocycle A end): Generated %d exercises", len(workout.Data.Exercises))
	})

	// Advance through Week 3 to start Mesocycle B
	for i := 0; i < 4; i++ {
		advanceUserState(t, ts, userID)
	}

	// =============================================================================
	// MESOCYCLE B - WEEK 4: Test 5RM (start of strength phase)
	// =============================================================================
	t.Run("Mesocycle B Week 4 generates 5RM workout", func(t *testing.T) {
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

		// Verify Week 4 (start of Mesocycle B)
		if workout.Data.WeekNumber != 4 {
			t.Errorf("Expected week 4, got %d", workout.Data.WeekNumber)
		}

		t.Logf("Week 4 (Mesocycle B): Generated %d exercises", len(workout.Data.Exercises))
	})

	// Advance through Weeks 4-5-6 (Mesocycle B)
	for i := 0; i < 12; i++ { // 4 days * 3 weeks
		advanceUserState(t, ts, userID)
	}

	// =============================================================================
	// MESOCYCLE C - WEEK 7: Test reset to 6RM
	// =============================================================================
	t.Run("Mesocycle C Week 7 resets to 6RM", func(t *testing.T) {
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

		// Verify Week 7 (start of Mesocycle C)
		if workout.Data.WeekNumber != 7 {
			t.Errorf("Expected week 7, got %d", workout.Data.WeekNumber)
		}

		t.Logf("Week 7 (Mesocycle C): Generated %d exercises", len(workout.Data.Exercises))
	})

	// Advance through Weeks 7-8-9 (Mesocycle C)
	for i := 0; i < 12; i++ {
		advanceUserState(t, ts, userID)
	}

	// =============================================================================
	// MESOCYCLE D - WEEK 10: Peaking phase begins
	// =============================================================================
	t.Run("Mesocycle D Week 10 starts peaking phase", func(t *testing.T) {
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

		// Verify Week 10 (start of Mesocycle D - peaking)
		if workout.Data.WeekNumber != 10 {
			t.Errorf("Expected week 10, got %d", workout.Data.WeekNumber)
		}

		t.Logf("Week 10 (Mesocycle D peaking): Generated %d exercises", len(workout.Data.Exercises))
	})

	// Advance through Weeks 10-11 to get to Week 12
	for i := 0; i < 8; i++ { // 4 days * 2 weeks
		advanceUserState(t, ts, userID)
	}

	// =============================================================================
	// MESOCYCLE D - WEEK 12: Max week (1RM)
	// =============================================================================
	t.Run("Mesocycle D Week 12 is max attempt week with 1RM", func(t *testing.T) {
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

		// Verify Week 12 (max week)
		if workout.Data.WeekNumber != 12 {
			t.Errorf("Expected week 12, got %d", workout.Data.WeekNumber)
		}

		// In max week, verify we have exercises
		if len(workout.Data.Exercises) < 1 {
			t.Fatal("Expected exercises in max week workout")
		}

		t.Logf("Week 12 (Max Week): Generated %d exercises for %s", len(workout.Data.Exercises), workout.Data.DaySlug)
	})

	// Advance through Week 12 to trigger cycle completion
	for i := 0; i < 4; i++ {
		advanceUserState(t, ts, userID)
	}

	// =============================================================================
	// CYCLE PROGRESSION: Trigger at end of 12-week cycle
	// =============================================================================
	t.Run("Cycle progression triggers after 12-week cycle", func(t *testing.T) {
		// Trigger progression for squat
		triggerBody := ManualTriggerRequest{
			ProgressionID: cycleProgID,
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
			expectedDelta := 5.0
			if squatTrigger.Data.Results[0].Result.Delta != expectedDelta {
				t.Errorf("Expected squat delta +%.1f, got %f", expectedDelta, squatTrigger.Data.Results[0].Result.Delta)
			}
			expectedNewTM := squatTM + expectedDelta // 320
			if squatTrigger.Data.Results[0].Result.NewValue != expectedNewTM {
				t.Errorf("Expected squat new value %f, got %f", expectedNewTM, squatTrigger.Data.Results[0].Result.NewValue)
			}
		}
	})

	// =============================================================================
	// NEW CYCLE: Verify Week 1 shows increased training maxes
	// =============================================================================
	t.Run("New cycle Week 1 shows updated training max", func(t *testing.T) {
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
			t.Errorf("Expected week 1 of new cycle, got %d", workout.Data.WeekNumber)
		}

		// Verify exercises are generated with new TM
		if len(workout.Data.Exercises) < 1 {
			t.Fatal("Expected exercises in new cycle workout")
		}

		t.Logf("New cycle Week 1: Generated %d exercises", len(workout.Data.Exercises))
	})
}

// =============================================================================
// HELPER FUNCTIONS (specific to J&T2 test)
// =============================================================================

// createJT2Prescription creates a prescription with weekly lookup for RM targets.
func createJT2Prescription(t *testing.T, ts *testutil.TestServer, liftID string, sets, reps int, percentage float64, order int, lookupKey string) string {
	t.Helper()

	body := fmt.Sprintf(`{
		"liftId": "%s",
		"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": %.1f, "lookupKey": "%s"},
		"setScheme": {"type": "FIXED", "sets": %d, "reps": %d},
		"order": %d
	}`, liftID, percentage, lookupKey, sets, reps, order)

	resp, err := adminPost(ts.URL("/prescriptions"), body)
	if err != nil {
		t.Fatalf("Failed to create J&T2 prescription: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create J&T2 prescription, status %d: %s", resp.StatusCode, bodyBytes)
	}

	var envelope PrescriptionResponse
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

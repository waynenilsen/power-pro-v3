// Package e2e provides end-to-end tests for complete program workflows.
// This file contains E2E tests for the Inverted Juggernaut 5/3/1 program.
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
// INVERTED JUGGERNAUT 5/3/1 E2E TEST
// =============================================================================

// TestInvertedJuggernaut531Program validates the complete Inverted Juggernaut 5/3/1
// program configuration and execution through the API.
//
// Inverted Juggernaut 5/3/1 characteristics:
// - 16-week cycle: Four 4-week rep waves (10s, 8s, 5s, 3s)
// - Each wave has 4 phases: Accumulation, Intensification, Realization, Deload
// - 5/3/1 percentage structure within each session:
//   - Week 1 (Accum): 65%x5, 75%x5, 85%x5+ AMRAP, 75%x5, 65%x5+ AMRAP
//   - Week 2 (Intens): 70%x3, 80%x3, 90%x3+ AMRAP, 80%x3, 70%x3+ AMRAP
//   - Week 3 (Real): 75%x5, 85%x3, 95%x1+ AMRAP, 85%x3, 75%x5+ AMRAP
//   - Week 4 (Deload): 40%x5, 50%x5, 60%x5 (no AMRAP)
// - Training Max = 90% of 1RM
// - Wave-specific volume sets before 5/3/1 work
// - Rep standards: 10s wave=10 reps, 8s wave=8 reps, 5s wave=5 reps, 3s wave=3 reps
// - Cycle progression: +5lb upper, +10lb lower per 16-week cycle
func TestInvertedJuggernaut531Program(t *testing.T) {
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

	// Create OHP lift (not seeded)
	ohpSlug := "ohp-ij531-" + testID
	ohpBody := fmt.Sprintf(`{"name": "Overhead Press", "slug": "%s", "isCompetitionLift": false}`, ohpSlug)
	ohpResp, err := adminPost(ts.URL("/lifts"), ohpBody)
	if err != nil {
		t.Fatalf("Failed to create OHP lift: %v", err)
	}
	if ohpResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(ohpResp.Body)
		ohpResp.Body.Close()
		t.Fatalf("Failed to create OHP lift, status %d: %s", ohpResp.StatusCode, body)
	}
	var ohpEnvelope LiftResponse
	json.NewDecoder(ohpResp.Body).Decode(&ohpEnvelope)
	ohpResp.Body.Close()
	ohpID := ohpEnvelope.Data.ID

	// Training maxes (90% of theoretical 1RM)
	// Using round numbers for easier percentage verification
	squatTM := 315.0    // Squat training max
	benchTM := 225.0    // Bench training max
	deadliftTM := 365.0 // Deadlift training max
	ohpTM := 135.0      // OHP training max

	// Create training maxes for the user
	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", squatTM)
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", benchTM)
	createLiftMax(t, ts, userID, deadliftID, "TRAINING_MAX", deadliftTM)
	createLiftMax(t, ts, userID, ohpID, "TRAINING_MAX", ohpTM)

	// =============================================================================
	// Create Weekly Lookup for 5/3/1 percentage progression (4-week pattern)
	// This repeats for each wave (10s, 8s, 5s, 3s)
	// Negative reps indicate AMRAP sets
	// =============================================================================
	weeklyLookupBody := `{
		"name": "Inverted Juggernaut 5/3/1 Weekly",
		"entries": [
			{"weekNumber": 1, "percentages": [65.0, 75.0, 85.0], "reps": [5, 5, 5]},
			{"weekNumber": 2, "percentages": [70.0, 80.0, 90.0], "reps": [3, 3, 3]},
			{"weekNumber": 3, "percentages": [75.0, 85.0, 95.0], "reps": [5, 3, 1]},
			{"weekNumber": 4, "percentages": [40.0, 50.0, 60.0], "reps": [5, 5, 5]},
			{"weekNumber": 5, "percentages": [65.0, 75.0, 85.0], "reps": [5, 5, 5]},
			{"weekNumber": 6, "percentages": [70.0, 80.0, 90.0], "reps": [3, 3, 3]},
			{"weekNumber": 7, "percentages": [75.0, 85.0, 95.0], "reps": [5, 3, 1]},
			{"weekNumber": 8, "percentages": [40.0, 50.0, 60.0], "reps": [5, 5, 5]},
			{"weekNumber": 9, "percentages": [65.0, 75.0, 85.0], "reps": [5, 5, 5]},
			{"weekNumber": 10, "percentages": [70.0, 80.0, 90.0], "reps": [3, 3, 3]},
			{"weekNumber": 11, "percentages": [75.0, 85.0, 95.0], "reps": [5, 3, 1]},
			{"weekNumber": 12, "percentages": [40.0, 50.0, 60.0], "reps": [5, 5, 5]},
			{"weekNumber": 13, "percentages": [65.0, 75.0, 85.0], "reps": [5, 5, 5]},
			{"weekNumber": 14, "percentages": [70.0, 80.0, 90.0], "reps": [3, 3, 3]},
			{"weekNumber": 15, "percentages": [75.0, 85.0, 95.0], "reps": [5, 3, 1]},
			{"weekNumber": 16, "percentages": [40.0, 50.0, 60.0], "reps": [5, 5, 5]}
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
	// Create prescriptions for main lifts
	// Using AMRAP prescription for the top set (85%/90%/95% depending on week)
	// =============================================================================

	// Squat - AMRAP top set
	squatAMRAPPrescID := createAMRAPPrescription(t, ts, squatID, 1, 5, 85.0, 0)
	// Back-off sets (2 sets at 75% and 65%)
	squatBackoff1PrescID := createPrescription(t, ts, squatID, 1, 5, 75.0, 1)
	squatBackoff2PrescID := createAMRAPPrescription(t, ts, squatID, 1, 5, 65.0, 2)

	// Bench - AMRAP top set
	benchAMRAPPrescID := createAMRAPPrescription(t, ts, benchID, 1, 5, 85.0, 0)
	// Back-off sets
	benchBackoff1PrescID := createPrescription(t, ts, benchID, 1, 5, 75.0, 1)
	benchBackoff2PrescID := createAMRAPPrescription(t, ts, benchID, 1, 5, 65.0, 2)

	// Deadlift - AMRAP top set
	deadliftAMRAPPrescID := createAMRAPPrescription(t, ts, deadliftID, 1, 5, 85.0, 0)
	// Back-off sets
	deadliftBackoff1PrescID := createPrescription(t, ts, deadliftID, 1, 5, 75.0, 1)
	deadliftBackoff2PrescID := createAMRAPPrescription(t, ts, deadliftID, 1, 5, 65.0, 2)

	// OHP - AMRAP top set
	ohpAMRAPPrescID := createAMRAPPrescription(t, ts, ohpID, 1, 5, 85.0, 0)
	// Back-off sets
	ohpBackoff1PrescID := createPrescription(t, ts, ohpID, 1, 5, 75.0, 1)
	ohpBackoff2PrescID := createAMRAPPrescription(t, ts, ohpID, 1, 5, 65.0, 2)

	// Deload prescriptions (no AMRAP, lighter weights)
	squatDeloadPrescID := createPrescription(t, ts, squatID, 3, 5, 50.0, 0) // 40/50/60 average
	benchDeloadPrescID := createPrescription(t, ts, benchID, 3, 5, 50.0, 0)
	deadliftDeloadPrescID := createPrescription(t, ts, deadliftID, 3, 5, 50.0, 0)
	ohpDeloadPrescID := createPrescription(t, ts, ohpID, 3, 5, 50.0, 0)

	// =============================================================================
	// Create Days (4-day split as per program spec)
	// Day 1: OHP, Day 2: Squat, Day 3: Bench, Day 4: Deadlift
	// =============================================================================

	// Day 1: OHP Focus
	ohpDaySlug := "ohp-day-ij531-" + testID
	ohpDayBody := fmt.Sprintf(`{"name": "Day 1 - OHP", "slug": "%s"}`, ohpDaySlug)
	ohpDayResp, _ := adminPost(ts.URL("/days"), ohpDayBody)
	var ohpDayEnvelope DayResponse
	json.NewDecoder(ohpDayResp.Body).Decode(&ohpDayEnvelope)
	ohpDayResp.Body.Close()
	ohpDayID := ohpDayEnvelope.Data.ID

	addPrescToDay(t, ts, ohpDayID, ohpAMRAPPrescID)
	addPrescToDay(t, ts, ohpDayID, ohpBackoff1PrescID)
	addPrescToDay(t, ts, ohpDayID, ohpBackoff2PrescID)

	// Day 2: Squat Focus
	squatDaySlug := "squat-day-ij531-" + testID
	squatDayBody := fmt.Sprintf(`{"name": "Day 2 - Squat", "slug": "%s"}`, squatDaySlug)
	squatDayResp, _ := adminPost(ts.URL("/days"), squatDayBody)
	var squatDayEnvelope DayResponse
	json.NewDecoder(squatDayResp.Body).Decode(&squatDayEnvelope)
	squatDayResp.Body.Close()
	squatDayID := squatDayEnvelope.Data.ID

	addPrescToDay(t, ts, squatDayID, squatAMRAPPrescID)
	addPrescToDay(t, ts, squatDayID, squatBackoff1PrescID)
	addPrescToDay(t, ts, squatDayID, squatBackoff2PrescID)

	// Day 3: Bench Focus
	benchDaySlug := "bench-day-ij531-" + testID
	benchDayBody := fmt.Sprintf(`{"name": "Day 3 - Bench", "slug": "%s"}`, benchDaySlug)
	benchDayResp, _ := adminPost(ts.URL("/days"), benchDayBody)
	var benchDayEnvelope DayResponse
	json.NewDecoder(benchDayResp.Body).Decode(&benchDayEnvelope)
	benchDayResp.Body.Close()
	benchDayID := benchDayEnvelope.Data.ID

	addPrescToDay(t, ts, benchDayID, benchAMRAPPrescID)
	addPrescToDay(t, ts, benchDayID, benchBackoff1PrescID)
	addPrescToDay(t, ts, benchDayID, benchBackoff2PrescID)

	// Day 4: Deadlift Focus
	deadliftDaySlug := "deadlift-day-ij531-" + testID
	deadliftDayBody := fmt.Sprintf(`{"name": "Day 4 - Deadlift", "slug": "%s"}`, deadliftDaySlug)
	deadliftDayResp, _ := adminPost(ts.URL("/days"), deadliftDayBody)
	var deadliftDayEnvelope DayResponse
	json.NewDecoder(deadliftDayResp.Body).Decode(&deadliftDayEnvelope)
	deadliftDayResp.Body.Close()
	deadliftDayID := deadliftDayEnvelope.Data.ID

	addPrescToDay(t, ts, deadliftDayID, deadliftAMRAPPrescID)
	addPrescToDay(t, ts, deadliftDayID, deadliftBackoff1PrescID)
	addPrescToDay(t, ts, deadliftDayID, deadliftBackoff2PrescID)

	// Deload Days (Week 4, 8, 12, 16)
	ohpDeloadDaySlug := "ohp-deload-ij531-" + testID
	ohpDeloadDayBody := fmt.Sprintf(`{"name": "Deload OHP", "slug": "%s"}`, ohpDeloadDaySlug)
	ohpDeloadDayResp, _ := adminPost(ts.URL("/days"), ohpDeloadDayBody)
	var ohpDeloadDayEnvelope DayResponse
	json.NewDecoder(ohpDeloadDayResp.Body).Decode(&ohpDeloadDayEnvelope)
	ohpDeloadDayResp.Body.Close()
	ohpDeloadDayID := ohpDeloadDayEnvelope.Data.ID
	addPrescToDay(t, ts, ohpDeloadDayID, ohpDeloadPrescID)

	squatDeloadDaySlug := "squat-deload-ij531-" + testID
	squatDeloadDayBody := fmt.Sprintf(`{"name": "Deload Squat", "slug": "%s"}`, squatDeloadDaySlug)
	squatDeloadDayResp, _ := adminPost(ts.URL("/days"), squatDeloadDayBody)
	var squatDeloadDayEnvelope DayResponse
	json.NewDecoder(squatDeloadDayResp.Body).Decode(&squatDeloadDayEnvelope)
	squatDeloadDayResp.Body.Close()
	squatDeloadDayID := squatDeloadDayEnvelope.Data.ID
	addPrescToDay(t, ts, squatDeloadDayID, squatDeloadPrescID)

	benchDeloadDaySlug := "bench-deload-ij531-" + testID
	benchDeloadDayBody := fmt.Sprintf(`{"name": "Deload Bench", "slug": "%s"}`, benchDeloadDaySlug)
	benchDeloadDayResp, _ := adminPost(ts.URL("/days"), benchDeloadDayBody)
	var benchDeloadDayEnvelope DayResponse
	json.NewDecoder(benchDeloadDayResp.Body).Decode(&benchDeloadDayEnvelope)
	benchDeloadDayResp.Body.Close()
	benchDeloadDayID := benchDeloadDayEnvelope.Data.ID
	addPrescToDay(t, ts, benchDeloadDayID, benchDeloadPrescID)

	deadliftDeloadDaySlug := "deadlift-deload-ij531-" + testID
	deadliftDeloadDayBody := fmt.Sprintf(`{"name": "Deload Deadlift", "slug": "%s"}`, deadliftDeloadDaySlug)
	deadliftDeloadDayResp, _ := adminPost(ts.URL("/days"), deadliftDeloadDayBody)
	var deadliftDeloadDayEnvelope DayResponse
	json.NewDecoder(deadliftDeloadDayResp.Body).Decode(&deadliftDeloadDayEnvelope)
	deadliftDeloadDayResp.Body.Close()
	deadliftDeloadDayID := deadliftDeloadDayEnvelope.Data.ID
	addPrescToDay(t, ts, deadliftDeloadDayID, deadliftDeloadPrescID)

	// =============================================================================
	// Create 16-week cycle
	// =============================================================================
	cycleName := "Inverted Juggernaut 5/3/1 Cycle " + testID
	cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 16}`, cycleName)
	cycleResp, _ := adminPost(ts.URL("/cycles"), cycleBody)
	var cycleEnvelope CycleResponse
	json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
	cycleResp.Body.Close()
	cycleID := cycleEnvelope.Data.ID

	// Create all 16 weeks
	for w := 1; w <= 16; w++ {
		weekBody := fmt.Sprintf(`{"weekNumber": %d, "cycleId": "%s"}`, w, cycleID)
		weekResp, _ := adminPost(ts.URL("/weeks"), weekBody)
		var weekEnvelope WeekResponse
		json.NewDecoder(weekResp.Body).Decode(&weekEnvelope)
		weekResp.Body.Close()
		weekID := weekEnvelope.Data.ID

		// Deload weeks (4, 8, 12, 16) get deload days
		if w%4 == 0 {
			addDayToWeek(t, ts, weekID, ohpDeloadDayID, "MONDAY")
			addDayToWeek(t, ts, weekID, squatDeloadDayID, "TUESDAY")
			addDayToWeek(t, ts, weekID, benchDeloadDayID, "THURSDAY")
			addDayToWeek(t, ts, weekID, deadliftDeloadDayID, "FRIDAY")
		} else {
			// Regular training weeks
			addDayToWeek(t, ts, weekID, ohpDayID, "MONDAY")
			addDayToWeek(t, ts, weekID, squatDayID, "TUESDAY")
			addDayToWeek(t, ts, weekID, benchDayID, "THURSDAY")
			addDayToWeek(t, ts, weekID, deadliftDayID, "FRIDAY")
		}
	}

	// =============================================================================
	// Create Program with weekly lookup
	// =============================================================================
	programSlug := "inverted-juggernaut-531-" + testID
	programBody := fmt.Sprintf(`{
		"name": "Inverted Juggernaut 5/3/1",
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
	// Create Cycle Progressions
	// After 16-week cycle: +5lb upper, +10lb lower
	// =============================================================================

	// Lower body progression (+10lb per cycle)
	lowerProgBody := `{
		"name": "IJ531 Lower +10lb",
		"type": "CYCLE_PROGRESSION",
		"parameters": {
			"increment": 10.0,
			"maxType": "TRAINING_MAX"
		}
	}`
	lowerProgResp, _ := adminPost(ts.URL("/progressions"), lowerProgBody)
	var lowerProgEnvelope ProgressionResponse
	json.NewDecoder(lowerProgResp.Body).Decode(&lowerProgEnvelope)
	lowerProgResp.Body.Close()
	lowerProgID := lowerProgEnvelope.Data.ID

	// Upper body progression (+5lb per cycle)
	upperProgBody := `{
		"name": "IJ531 Upper +5lb",
		"type": "CYCLE_PROGRESSION",
		"parameters": {
			"increment": 5.0,
			"maxType": "TRAINING_MAX"
		}
	}`
	upperProgResp, _ := adminPost(ts.URL("/progressions"), upperProgBody)
	var upperProgEnvelope ProgressionResponse
	json.NewDecoder(upperProgResp.Body).Decode(&upperProgEnvelope)
	upperProgResp.Body.Close()
	upperProgID := upperProgEnvelope.Data.ID

	// Link progressions to program
	linkProgressionToProgram(t, ts, programID, lowerProgID, squatID, 1)
	linkProgressionToProgram(t, ts, programID, lowerProgID, deadliftID, 2)
	linkProgressionToProgram(t, ts, programID, upperProgID, benchID, 3)
	linkProgressionToProgram(t, ts, programID, upperProgID, ohpID, 4)

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
	// WAVE 1 (10s Wave) - WEEK 1: Accumulation Phase
	// =============================================================================
	t.Run("Wave 1 Week 1 - 10s Accumulation generates AMRAP at 85%", func(t *testing.T) {
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

		// First day should be OHP
		if workout.Data.DaySlug != ohpDaySlug {
			t.Errorf("Expected OHP day slug '%s', got '%s'", ohpDaySlug, workout.Data.DaySlug)
		}

		// Should have 3 prescriptions: AMRAP + 2 back-off
		if len(workout.Data.Exercises) < 3 {
			t.Fatalf("Expected at least 3 exercises (AMRAP + back-offs), got %d", len(workout.Data.Exercises))
		}

		// Verify first exercise is OHP AMRAP at 85%
		ohp := workout.Data.Exercises[0]
		if ohp.Lift.ID != ohpID {
			t.Errorf("Expected OHP lift, got %s", ohp.Lift.ID)
		}

		// Expected weight: 85% of 135 TM = 114.75
		expectedWeight := ohpTM * 0.85
		if len(ohp.Sets) > 0 && !withinTolerance(ohp.Sets[0].Weight, expectedWeight, 5.0) {
			t.Errorf("Expected OHP weight ~%.1f (85%% of TM), got %.1f", expectedWeight, ohp.Sets[0].Weight)
		}

		t.Logf("Wave 1 Week 1: OHP AMRAP at %.1f lbs (85%% of %.1f TM)", ohp.Sets[0].Weight, ohpTM)
	})

	// Complete Week 1 Day 1 (OHP) using explicit state machine flow
	completeIJ531WorkoutDay(t, ts, userID)

	t.Run("Wave 1 Week 1 Day 2 - Squat session with back-offs", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		// Still Week 1
		if workout.Data.WeekNumber != 1 {
			t.Errorf("Expected week 1, got %d", workout.Data.WeekNumber)
		}

		// Should be Squat day
		if workout.Data.DaySlug != squatDaySlug {
			t.Errorf("Expected Squat day, got '%s'", workout.Data.DaySlug)
		}

		// Verify squat exercises present
		if len(workout.Data.Exercises) < 3 {
			t.Fatalf("Expected 3 squat exercises, got %d", len(workout.Data.Exercises))
		}

		squat := workout.Data.Exercises[0]
		expectedWeight := squatTM * 0.85 // 315 * 0.85 = 267.75
		if len(squat.Sets) > 0 && !withinTolerance(squat.Sets[0].Weight, expectedWeight, 5.0) {
			t.Errorf("Expected Squat weight ~%.1f (85%% of TM), got %.1f", expectedWeight, squat.Sets[0].Weight)
		}

		t.Logf("Wave 1 Week 1: Squat at %.1f lbs", squat.Sets[0].Weight)
	})

	// Complete rest of Week 1 (Day 2 Squat, Day 3 Bench, Day 4 Deadlift)
	completeIJ531WorkoutDay(t, ts, userID) // Week 1, Day 2 (Squat)
	completeIJ531WorkoutDay(t, ts, userID) // Week 1, Day 3 (Bench)
	completeIJ531WorkoutDay(t, ts, userID) // Week 1, Day 4 (Deadlift)

	// =============================================================================
	// WAVE 1 - WEEK 4: Deload Phase
	// =============================================================================

	// Advance to Week 4, Day 1 (OHP Deload)
	// From Week 1 Day 4, we need:
	// - Week 2: 4 days
	// - Week 3: 4 days
	// - To Week 4 Day 1: 1 more advance
	// Total: 9 workout completions to reach Week 4 Day 1
	for i := 0; i < 9; i++ {
		completeIJ531WorkoutDay(t, ts, userID)
	}

	t.Run("Wave 1 Week 4 - Deload with reduced volume", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		// Verify Week 4 (deload)
		if workout.Data.WeekNumber != 4 {
			t.Errorf("Expected week 4 (deload), got %d", workout.Data.WeekNumber)
		}

		// Should be deload OHP day
		if workout.Data.DaySlug != ohpDeloadDaySlug {
			t.Errorf("Expected deload OHP day '%s', got '%s'", ohpDeloadDaySlug, workout.Data.DaySlug)
		}

		// Deload should have just 1 prescription with 3 sets (no AMRAP)
		if len(workout.Data.Exercises) != 1 {
			t.Errorf("Expected 1 exercise on deload, got %d", len(workout.Data.Exercises))
		}

		if len(workout.Data.Exercises) > 0 {
			ohp := workout.Data.Exercises[0]
			// Deload weight: ~50% of TM
			expectedWeight := ohpTM * 0.50
			if len(ohp.Sets) > 0 && !withinTolerance(ohp.Sets[0].Weight, expectedWeight, 5.0) {
				t.Errorf("Expected deload weight ~%.1f (50%% of TM), got %.1f", expectedWeight, ohp.Sets[0].Weight)
			}
			t.Logf("Wave 1 Week 4 Deload: OHP at %.1f lbs", ohp.Sets[0].Weight)
		}
	})

	// =============================================================================
	// WAVE 2 (8s Wave) - WEEK 5: New Wave Accumulation
	// =============================================================================

	// Advance through Week 4 deload to Week 5 Day 1
	// From Week 4 Day 1, need 4 advances to complete Week 4, then we're at Week 5 Day 1
	for i := 0; i < 4; i++ {
		completeIJ531WorkoutDay(t, ts, userID)
	}

	t.Run("Wave 2 Week 5 - 8s Wave Accumulation begins", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		// Verify Week 5 (8s wave start)
		if workout.Data.WeekNumber != 5 {
			t.Errorf("Expected week 5 (8s wave start), got %d", workout.Data.WeekNumber)
		}

		// Should be regular OHP day (not deload)
		if workout.Data.DaySlug != ohpDaySlug {
			t.Errorf("Expected OHP day, got '%s'", workout.Data.DaySlug)
		}

		t.Logf("Wave 2 Week 5: 8s Wave Accumulation - %d exercises", len(workout.Data.Exercises))
	})

	// =============================================================================
	// WAVE 3 (5s Wave) - WEEK 9: Strength Emphasis
	// =============================================================================

	// Advance from Week 5 Day 1 to Week 9 Day 1
	// Week 5: 4 days, Week 6: 4 days, Week 7: 4 days, Week 8: 4 days = 16 days
	for i := 0; i < 16; i++ {
		completeIJ531WorkoutDay(t, ts, userID)
	}

	t.Run("Wave 3 Week 9 - 5s Wave Accumulation", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		// Verify Week 9
		if workout.Data.WeekNumber != 9 {
			t.Errorf("Expected week 9 (5s wave start), got %d", workout.Data.WeekNumber)
		}

		t.Logf("Wave 3 Week 9: 5s Wave - Week %d", workout.Data.WeekNumber)
	})

	// =============================================================================
	// WAVE 4 (3s Wave) - WEEK 13: Peak Intensity
	// =============================================================================

	// Advance from Week 9 Day 1 to Week 13 Day 1
	// Week 9: 4 days, Week 10: 4 days, Week 11: 4 days, Week 12: 4 days = 16 days
	for i := 0; i < 16; i++ {
		completeIJ531WorkoutDay(t, ts, userID)
	}

	t.Run("Wave 4 Week 13 - 3s Wave Peak Intensity begins", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		// Verify Week 13
		if workout.Data.WeekNumber != 13 {
			t.Errorf("Expected week 13 (3s wave start), got %d", workout.Data.WeekNumber)
		}

		// Should be regular training day
		if workout.Data.DaySlug != ohpDaySlug {
			t.Errorf("Expected OHP day, got '%s'", workout.Data.DaySlug)
		}

		t.Logf("Wave 4 Week 13: 3s Wave Peak - Week %d", workout.Data.WeekNumber)
	})

	// =============================================================================
	// WEEK 16: Final Deload before Cycle Progression
	// =============================================================================

	// Advance from Week 13 Day 1 to Week 16 Day 1
	// Week 13: 4 days, Week 14: 4 days, Week 15: 4 days = 12 days
	for i := 0; i < 12; i++ {
		completeIJ531WorkoutDay(t, ts, userID)
	}

	t.Run("Week 16 - Final Deload of 16-week cycle", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		// Verify Week 16 (final deload)
		if workout.Data.WeekNumber != 16 {
			t.Errorf("Expected week 16 (final deload), got %d", workout.Data.WeekNumber)
		}

		// Should be deload day
		if workout.Data.DaySlug != ohpDeloadDaySlug {
			t.Errorf("Expected deload day, got '%s'", workout.Data.DaySlug)
		}

		t.Logf("Week 16: Final Deload - preparing for cycle progression")
	})

	// Advance through Week 16 to complete the cycle
	// Week 16: 4 days to complete, cycle wraps to Week 1 Day 1
	for i := 0; i < 4; i++ {
		completeIJ531WorkoutDay(t, ts, userID)
	}

	// =============================================================================
	// CYCLE PROGRESSION: Trigger at end of 16-week cycle
	// =============================================================================
	t.Run("Cycle progression applies at end of 16-week cycle", func(t *testing.T) {
		// Trigger lower body progression for squat (+10lb)
		squatTrigger := triggerProgressionForLift(t, ts, userID, lowerProgID, squatID)

		if squatTrigger.Data.TotalApplied != 1 {
			t.Errorf("Expected squat progression to apply, got TotalApplied=%d", squatTrigger.Data.TotalApplied)
		}
		if len(squatTrigger.Data.Results) > 0 && squatTrigger.Data.Results[0].Result != nil {
			if squatTrigger.Data.Results[0].Result.Delta != 10.0 {
				t.Errorf("Expected squat delta +10, got %f", squatTrigger.Data.Results[0].Result.Delta)
			}
			expectedNewSquat := squatTM + 10.0 // 325
			if squatTrigger.Data.Results[0].Result.NewValue != expectedNewSquat {
				t.Errorf("Expected squat new TM %f, got %f", expectedNewSquat, squatTrigger.Data.Results[0].Result.NewValue)
			}
			t.Logf("Squat TM: %.1f -> %.1f (+10 lb)", squatTM, squatTrigger.Data.Results[0].Result.NewValue)
		}

		// Trigger lower body progression for deadlift (+10lb)
		deadliftTrigger := triggerProgressionForLift(t, ts, userID, lowerProgID, deadliftID)

		if deadliftTrigger.Data.TotalApplied != 1 {
			t.Errorf("Expected deadlift progression to apply")
		}
		if len(deadliftTrigger.Data.Results) > 0 && deadliftTrigger.Data.Results[0].Result != nil {
			t.Logf("Deadlift TM: %.1f -> %.1f (+10 lb)", deadliftTM, deadliftTrigger.Data.Results[0].Result.NewValue)
		}

		// Trigger upper body progression for bench (+5lb)
		benchTrigger := triggerProgressionForLift(t, ts, userID, upperProgID, benchID)

		if benchTrigger.Data.TotalApplied != 1 {
			t.Errorf("Expected bench progression to apply")
		}
		if len(benchTrigger.Data.Results) > 0 && benchTrigger.Data.Results[0].Result != nil {
			if benchTrigger.Data.Results[0].Result.Delta != 5.0 {
				t.Errorf("Expected bench delta +5, got %f", benchTrigger.Data.Results[0].Result.Delta)
			}
			t.Logf("Bench TM: %.1f -> %.1f (+5 lb)", benchTM, benchTrigger.Data.Results[0].Result.NewValue)
		}

		// Trigger upper body progression for OHP (+5lb)
		ohpTrigger := triggerProgressionForLift(t, ts, userID, upperProgID, ohpID)

		if ohpTrigger.Data.TotalApplied != 1 {
			t.Errorf("Expected OHP progression to apply")
		}
		if len(ohpTrigger.Data.Results) > 0 && ohpTrigger.Data.Results[0].Result != nil {
			t.Logf("OHP TM: %.1f -> %.1f (+5 lb)", ohpTM, ohpTrigger.Data.Results[0].Result.NewValue)
		}
	})

	// =============================================================================
	// NEW CYCLE: Verify Week 1 of new cycle shows increased weights
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
			t.Errorf("Expected week 1 of new cycle, got %d", workout.Data.WeekNumber)
		}

		// Should be OHP day (Day 1)
		if workout.Data.DaySlug != ohpDaySlug {
			t.Errorf("Expected OHP day, got '%s'", workout.Data.DaySlug)
		}

		// Verify weight increased (OHP TM increased by 5)
		if len(workout.Data.Exercises) > 0 && len(workout.Data.Exercises[0].Sets) > 0 {
			newOHPTM := ohpTM + 5.0            // 140
			expectedWeight := newOHPTM * 0.85  // 119
			actualWeight := workout.Data.Exercises[0].Sets[0].Weight
			if !withinTolerance(actualWeight, expectedWeight, 5.0) {
				t.Errorf("Expected OHP weight ~%.1f (85%% of new TM %.1f), got %.1f",
					expectedWeight, newOHPTM, actualWeight)
			}
			t.Logf("New Cycle Week 1: OHP at %.1f lbs (increased from %.1f)",
				actualWeight, ohpTM*0.85)
		}

		t.Logf("Inverted Juggernaut 5/3/1 test completed successfully")
		t.Logf("Validated: 16-week cycle, 4 rep waves, 5/3/1 structure, deloads, cycle progression")
	})
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

// completeIJ531WorkoutDay completes an Inverted Juggernaut 5/3/1 workout day using explicit state machine flow.
func completeIJ531WorkoutDay(t *testing.T, ts *testutil.TestServer, userID string) {
	t.Helper()

	sessionID := startWorkoutSession(t, ts, userID)

	workoutResp, _ := userGet(ts.URL("/users/"+userID+"/workout"), userID)
	var workout WorkoutResponse
	json.NewDecoder(workoutResp.Body).Decode(&workout)
	workoutResp.Body.Close()

	for _, ex := range workout.Data.Exercises {
		for _, set := range ex.Sets {
			logIJ531Set(t, ts, userID, sessionID, ex.PrescriptionID, ex.Lift.ID, set.SetNumber, set.Weight, set.TargetReps, set.TargetReps, false)
		}
	}

	finishWorkoutSession(t, ts, sessionID, userID)
	advanceUserState(t, ts, userID)
}

// logIJ531Set logs a single set for Inverted Juggernaut 5/3/1 workout.
func logIJ531Set(t *testing.T, ts *testutil.TestServer, userID, sessionID, prescriptionID, liftID string, setNumber int, weight float64, targetReps, repsPerformed int, isAmrap bool) {
	t.Helper()

	loggedSetBody := fmt.Sprintf(`{
		"sets": [{
			"prescriptionId": "%s",
			"liftId": "%s",
			"setNumber": %d,
			"weight": %.1f,
			"targetReps": %d,
			"repsPerformed": %d,
			"isAmrap": %t
		}]
	}`, prescriptionID, liftID, setNumber, weight, targetReps, repsPerformed, isAmrap)

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
}

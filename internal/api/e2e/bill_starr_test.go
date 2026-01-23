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

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// =============================================================================
// BILL STARR 5x5 E2E TEST
// =============================================================================

// TestBillStarr5x5Program validates the complete Bill Starr 5x5 program
// configuration and execution through the API.
//
// Bill Starr 5x5 characteristics:
// - Heavy/Light/Medium Days: Different intensities per day (DailyLookup pattern)
//   - Heavy Day: 100% of working weight
//   - Light Day: ~80% of working weight
//   - Medium Day: ~90% of working weight
// - Ramp Sets: 5x5 where weight increases each set (50%, 63%, 75%, 88%, 100%)
// - LinearProgression: AFTER_WEEK trigger with +5lb increment
func TestBillStarr5x5Program(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Test-unique identifiers
	testID := uuid.New().String()[:8]
	// Use a seeded test user (required for foreign key constraints)
	userID := "bill-starr-test-user"

	// Seeded lift IDs
	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"

	// Create Row lift (not seeded)
	rowSlug := "row-" + testID
	rowBody := fmt.Sprintf(`{"name": "Bent Over Row", "slug": "%s", "isCompetitionLift": false}`, rowSlug)
	rowResp, err := adminPost(ts.URL("/lifts"), rowBody)
	if err != nil {
		t.Fatalf("Failed to create row lift: %v", err)
	}
	if rowResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(rowResp.Body)
		rowResp.Body.Close()
		t.Fatalf("Failed to create row lift, status %d: %s", rowResp.StatusCode, body)
	}
	var rowEnvelope LiftResponse
	json.NewDecoder(rowResp.Body).Decode(&rowEnvelope)
	rowResp.Body.Close()
	rowID := rowEnvelope.Data.ID

	// Bill Starr training maxes
	squatMax := 300.0 // Squat training max
	benchMax := 200.0 // Bench training max
	rowMax := 185.0   // Row training max

	// Create training maxes for the user
	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", squatMax)
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", benchMax)
	createLiftMax(t, ts, userID, rowID, "TRAINING_MAX", rowMax)

	// =============================================================================
	// Create Daily Lookup for Heavy/Light/Medium intensities
	// =============================================================================
	dailyLookupBody := `{
		"name": "Bill Starr H/L/M Intensity",
		"entries": [
			{"dayIdentifier": "heavy", "percentageModifier": 100.0, "intensityLevel": "HEAVY"},
			{"dayIdentifier": "light", "percentageModifier": 80.0, "intensityLevel": "LIGHT"},
			{"dayIdentifier": "medium", "percentageModifier": 90.0, "intensityLevel": "MEDIUM"}
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
	// Create prescriptions with RAMP set scheme
	// Bill Starr uses 5x5 ramping: 50%, 63%, 75%, 88%, 100% of top set
	// =============================================================================

	// Standard 5x5 Ramp prescription (used for Heavy Day and Medium Day)
	squat5x5RampID := createRampPrescription(t, ts, squatID, []RampStep{
		{Percentage: 50, Reps: 5},
		{Percentage: 63, Reps: 5},
		{Percentage: 75, Reps: 5},
		{Percentage: 88, Reps: 5},
		{Percentage: 100, Reps: 5},
	}, 80.0, 0, "day")

	bench5x5RampID := createRampPrescription(t, ts, benchID, []RampStep{
		{Percentage: 50, Reps: 5},
		{Percentage: 63, Reps: 5},
		{Percentage: 75, Reps: 5},
		{Percentage: 88, Reps: 5},
		{Percentage: 100, Reps: 5},
	}, 80.0, 1, "day")

	row5x5RampID := createRampPrescription(t, ts, rowID, []RampStep{
		{Percentage: 50, Reps: 5},
		{Percentage: 63, Reps: 5},
		{Percentage: 75, Reps: 5},
		{Percentage: 88, Reps: 5},
		{Percentage: 100, Reps: 5},
	}, 80.0, 2, "day")

	// Light Day uses a capped 4x5 ramp (stops at 75%)
	squatLightRampID := createRampPrescription(t, ts, squatID, []RampStep{
		{Percentage: 50, Reps: 5},
		{Percentage: 63, Reps: 5},
		{Percentage: 75, Reps: 5},
		{Percentage: 75, Reps: 5}, // Repeated - capped at 75%
	}, 70.0, 0, "day")

	// =============================================================================
	// Create Days: Heavy, Light, Medium
	// The day slug must match the dailyLookup dayIdentifier for intensity lookup
	// =============================================================================

	// Heavy Day
	heavyDayBody := `{"name": "Heavy Day", "slug": "heavy"}`
	heavyDayResp, _ := adminPost(ts.URL("/days"), heavyDayBody)
	var heavyDayEnvelope DayResponse
	json.NewDecoder(heavyDayResp.Body).Decode(&heavyDayEnvelope)
	heavyDayResp.Body.Close()
	heavyDayID := heavyDayEnvelope.Data.ID

	// Add prescriptions to Heavy Day
	addPrescToDay(t, ts, heavyDayID, squat5x5RampID)
	addPrescToDay(t, ts, heavyDayID, bench5x5RampID)
	addPrescToDay(t, ts, heavyDayID, row5x5RampID)

	// Light Day
	lightDayBody := `{"name": "Light Day", "slug": "light"}`
	lightDayResp, _ := adminPost(ts.URL("/days"), lightDayBody)
	var lightDayEnvelope DayResponse
	json.NewDecoder(lightDayResp.Body).Decode(&lightDayEnvelope)
	lightDayResp.Body.Close()
	lightDayID := lightDayEnvelope.Data.ID

	// Add light prescriptions to Light Day
	addPrescToDay(t, ts, lightDayID, squatLightRampID)

	// Medium Day (uses same 5x5 ramp prescriptions as Heavy, but daily lookup scales to 90%)
	mediumDayBody := `{"name": "Medium Day", "slug": "medium"}`
	mediumDayResp, _ := adminPost(ts.URL("/days"), mediumDayBody)
	var mediumDayEnvelope DayResponse
	json.NewDecoder(mediumDayResp.Body).Decode(&mediumDayEnvelope)
	mediumDayResp.Body.Close()
	mediumDayID := mediumDayEnvelope.Data.ID

	// Need separate prescriptions for medium day to avoid sharing
	squat5x5RampMediumID := createRampPrescription(t, ts, squatID, []RampStep{
		{Percentage: 50, Reps: 5},
		{Percentage: 63, Reps: 5},
		{Percentage: 75, Reps: 5},
		{Percentage: 88, Reps: 5},
		{Percentage: 100, Reps: 5},
	}, 80.0, 0, "day")

	bench5x5RampMediumID := createRampPrescription(t, ts, benchID, []RampStep{
		{Percentage: 50, Reps: 5},
		{Percentage: 63, Reps: 5},
		{Percentage: 75, Reps: 5},
		{Percentage: 88, Reps: 5},
		{Percentage: 100, Reps: 5},
	}, 80.0, 1, "day")

	row5x5RampMediumID := createRampPrescription(t, ts, rowID, []RampStep{
		{Percentage: 50, Reps: 5},
		{Percentage: 63, Reps: 5},
		{Percentage: 75, Reps: 5},
		{Percentage: 88, Reps: 5},
		{Percentage: 100, Reps: 5},
	}, 80.0, 2, "day")

	addPrescToDay(t, ts, mediumDayID, squat5x5RampMediumID)
	addPrescToDay(t, ts, mediumDayID, bench5x5RampMediumID)
	addPrescToDay(t, ts, mediumDayID, row5x5RampMediumID)

	// =============================================================================
	// Create 1-week cycle with H/L/M pattern (Mon/Wed/Fri)
	// =============================================================================
	cycleName := "BS5x5 Cycle " + testID
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

	// Add days to week: Heavy/Light/Medium pattern
	addDayToWeek(t, ts, weekID, heavyDayID, "MONDAY")
	addDayToWeek(t, ts, weekID, lightDayID, "WEDNESDAY")
	addDayToWeek(t, ts, weekID, mediumDayID, "FRIDAY")

	// =============================================================================
	// Create program and link daily lookup
	// =============================================================================
	programSlug := "bill-starr-5x5-" + testID
	programBody := fmt.Sprintf(`{"name": "Bill Starr 5x5", "slug": "%s", "cycleId": "%s", "dailyLookupId": "%s"}`,
		programSlug, cycleID, dailyLookupID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID := programEnvelope.Data.ID

	// =============================================================================
	// Create Linear Progression (AFTER_WEEK trigger with +5lb)
	// =============================================================================
	weeklyProgBody := `{"name": "BS5x5 Weekly Linear", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_WEEK"}}`
	weeklyProgResp, _ := adminPost(ts.URL("/progressions"), weeklyProgBody)
	var weeklyProgEnvelope ProgressionResponse
	json.NewDecoder(weeklyProgResp.Body).Decode(&weeklyProgEnvelope)
	weeklyProgResp.Body.Close()
	weeklyProgID := weeklyProgEnvelope.Data.ID

	// Link progression to program for each lift
	linkProgressionToProgram(t, ts, programID, weeklyProgID, squatID, 1)
	linkProgressionToProgram(t, ts, programID, weeklyProgID, benchID, 2)
	linkProgressionToProgram(t, ts, programID, weeklyProgID, rowID, 3)

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
	// EXECUTION PHASE: Heavy Day (Workout 1)
	// =============================================================================
	t.Run("Heavy Day generates RAMP sets at 100% intensity", func(t *testing.T) {
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

		// Verify Heavy Day
		if workout.Data.DaySlug != "heavy" {
			t.Errorf("Expected day slug 'heavy', got '%s'", workout.Data.DaySlug)
		}

		if len(workout.Data.Exercises) != 3 {
			t.Fatalf("Expected 3 exercises on Heavy Day, got %d", len(workout.Data.Exercises))
		}

		exercisesByLift := make(map[string]WorkoutExerciseData)
		for _, ex := range workout.Data.Exercises {
			exercisesByLift[ex.Lift.ID] = ex
		}

		// Squat on Heavy Day: 5 ramping sets at 100% daily intensity
		// Expected weights: 50%, 63%, 75%, 88%, 100% of 300 = 150, 189, 225, 264, 300
		if squat, ok := exercisesByLift[squatID]; ok {
			if len(squat.Sets) != 5 {
				t.Errorf("Squat: expected 5 sets, got %d", len(squat.Sets))
			}
			expectedPercentages := []float64{0.50, 0.63, 0.75, 0.88, 1.00}
			for i, set := range squat.Sets {
				expectedWeight := squatMax * expectedPercentages[i]
				if !closeEnough(set.Weight, expectedWeight) {
					t.Errorf("Squat set %d: expected weight ~%.1f, got %.1f", i+1, expectedWeight, set.Weight)
				}
				if set.TargetReps != 5 {
					t.Errorf("Squat set %d: expected 5 reps, got %d", i+1, set.TargetReps)
				}
			}
			// Verify work set classification (80% threshold: sets 4 and 5 should be work sets)
			if squat.Sets[3].IsWorkSet != true || squat.Sets[4].IsWorkSet != true {
				t.Error("Squat sets 4 and 5 should be work sets (>=80%)")
			}
			if squat.Sets[0].IsWorkSet != false || squat.Sets[1].IsWorkSet != false || squat.Sets[2].IsWorkSet != false {
				t.Error("Squat sets 1-3 should be warmup sets (<80%)")
			}
		} else {
			t.Error("Heavy Day missing Squat exercise")
		}

		// Bench on Heavy Day
		if bench, ok := exercisesByLift[benchID]; ok {
			if len(bench.Sets) != 5 {
				t.Errorf("Bench: expected 5 sets, got %d", len(bench.Sets))
			}
			// Top set at 100% of 200 = 200
			topSet := bench.Sets[len(bench.Sets)-1]
			if !closeEnough(topSet.Weight, benchMax) {
				t.Errorf("Bench top set: expected %.1f, got %.1f", benchMax, topSet.Weight)
			}
		} else {
			t.Error("Heavy Day missing Bench exercise")
		}

		// Row on Heavy Day
		if row, ok := exercisesByLift[rowID]; ok {
			if len(row.Sets) != 5 {
				t.Errorf("Row: expected 5 sets, got %d", len(row.Sets))
			}
			// Top set at 100% of 185 = 185
			topSet := row.Sets[len(row.Sets)-1]
			if !closeEnough(topSet.Weight, rowMax) {
				t.Errorf("Row top set: expected %.1f, got %.1f", rowMax, topSet.Weight)
			}
		} else {
			t.Error("Heavy Day missing Row exercise")
		}
	})

	// Advance to Light Day
	advanceUserState(t, ts, userID)

	// =============================================================================
	// EXECUTION PHASE: Light Day (Workout 2)
	// =============================================================================
	t.Run("Light Day generates RAMP sets at 80% intensity", func(t *testing.T) {
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

		// Verify Light Day
		if workout.Data.DaySlug != "light" {
			t.Errorf("Expected day slug 'light', got '%s'", workout.Data.DaySlug)
		}

		exercisesByLift := make(map[string]WorkoutExerciseData)
		for _, ex := range workout.Data.Exercises {
			exercisesByLift[ex.Lift.ID] = ex
		}

		// Light Day squat uses capped ramp (4 sets, max 75%) at 80% daily intensity
		// Training max: 300, Daily intensity: 80% = 240 effective top
		// Ramp percentages: 50%, 63%, 75%, 75% (capped)
		// Expected: 50% of 240 = 120, 63% of 240 = 151.2, 75% of 240 = 180, 75% of 240 = 180
		if squat, ok := exercisesByLift[squatID]; ok {
			if len(squat.Sets) != 4 {
				t.Errorf("Light Day Squat: expected 4 sets (capped ramp), got %d", len(squat.Sets))
			}
			// The daily lookup applies 80% modifier, then ramp percentages apply
			// Base is 300 * 80% = 240, then ramp steps apply to that
			dailyModifier := 0.80
			expectedPercentages := []float64{0.50, 0.63, 0.75, 0.75}
			for i, set := range squat.Sets {
				expectedWeight := squatMax * dailyModifier * expectedPercentages[i]
				if !closeEnough(set.Weight, expectedWeight) {
					t.Errorf("Light Squat set %d: expected weight ~%.1f, got %.1f", i+1, expectedWeight, set.Weight)
				}
			}
		} else {
			t.Error("Light Day missing Squat exercise")
		}
	})

	// Advance to Medium Day
	advanceUserState(t, ts, userID)

	// =============================================================================
	// EXECUTION PHASE: Medium Day (Workout 3)
	// =============================================================================
	t.Run("Medium Day generates RAMP sets at 90% intensity", func(t *testing.T) {
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

		// Verify Medium Day
		if workout.Data.DaySlug != "medium" {
			t.Errorf("Expected day slug 'medium', got '%s'", workout.Data.DaySlug)
		}

		if len(workout.Data.Exercises) != 3 {
			t.Fatalf("Expected 3 exercises on Medium Day, got %d", len(workout.Data.Exercises))
		}

		exercisesByLift := make(map[string]WorkoutExerciseData)
		for _, ex := range workout.Data.Exercises {
			exercisesByLift[ex.Lift.ID] = ex
		}

		// Medium Day uses 90% daily intensity
		// Squat: 300 * 90% = 270 effective top, then ramp: 50%, 63%, 75%, 88%, 100%
		if squat, ok := exercisesByLift[squatID]; ok {
			if len(squat.Sets) != 5 {
				t.Errorf("Medium Day Squat: expected 5 sets, got %d", len(squat.Sets))
			}
			dailyModifier := 0.90
			expectedPercentages := []float64{0.50, 0.63, 0.75, 0.88, 1.00}
			for i, set := range squat.Sets {
				expectedWeight := squatMax * dailyModifier * expectedPercentages[i]
				if !closeEnough(set.Weight, expectedWeight) {
					t.Errorf("Medium Squat set %d: expected weight ~%.1f, got %.1f", i+1, expectedWeight, set.Weight)
				}
			}
			// Top set should be 270 (300 * 90%)
			topSet := squat.Sets[len(squat.Sets)-1]
			expectedTop := squatMax * 0.90
			if !closeEnough(topSet.Weight, expectedTop) {
				t.Errorf("Medium Squat top set: expected %.1f, got %.1f", expectedTop, topSet.Weight)
			}
		} else {
			t.Error("Medium Day missing Squat exercise")
		}

		// Verify Bench at 90% intensity
		if bench, ok := exercisesByLift[benchID]; ok {
			topSet := bench.Sets[len(bench.Sets)-1]
			expectedTop := benchMax * 0.90 // 200 * 90% = 180
			if !closeEnough(topSet.Weight, expectedTop) {
				t.Errorf("Medium Bench top set: expected %.1f, got %.1f", expectedTop, topSet.Weight)
			}
		}

		// Verify Row at 90% intensity
		// Note: 185 * 90% = 166.5, may round to 165 depending on rounding config
		if row, ok := exercisesByLift[rowID]; ok {
			topSet := row.Sets[len(row.Sets)-1]
			expectedTop := rowMax * 0.90 // 185 * 90% = 166.5
			// Allow for rounding differences (165 or 166.5)
			if !closeEnough(topSet.Weight, expectedTop) && !closeEnough(topSet.Weight, 165.0) {
				t.Errorf("Medium Row top set: expected ~%.1f (with rounding), got %.1f", expectedTop, topSet.Weight)
			}
		}
	})

	// =============================================================================
	// PROGRESSION PHASE: Trigger weekly progression after completing the week
	// =============================================================================
	t.Run("Weekly progression increases all lifts by +5lb", func(t *testing.T) {
		// Trigger progression for squat
		triggerBody := ManualTriggerRequest{
			ProgressionID: weeklyProgID,
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
			if squatTrigger.Data.Results[0].Result.Delta != 5.0 {
				t.Errorf("Expected squat delta +5, got %f", squatTrigger.Data.Results[0].Result.Delta)
			}
			expectedNewSquat := squatMax + 5.0
			if squatTrigger.Data.Results[0].Result.NewValue != expectedNewSquat {
				t.Errorf("Expected squat new value %f, got %f", expectedNewSquat, squatTrigger.Data.Results[0].Result.NewValue)
			}
		}

		// Trigger progression for bench
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
		}

		// Trigger progression for row
		triggerBody.LiftID = rowID
		triggerResp, err = authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), triggerBody, userID)
		if err != nil {
			t.Fatalf("Failed to trigger row progression: %v", err)
		}
		var rowTrigger TriggerResponse
		json.NewDecoder(triggerResp.Body).Decode(&rowTrigger)
		triggerResp.Body.Close()

		if rowTrigger.Data.TotalApplied != 1 {
			t.Errorf("Expected row progression to apply")
		}
		if len(rowTrigger.Data.Results) > 0 && rowTrigger.Data.Results[0].Result != nil {
			if rowTrigger.Data.Results[0].Result.Delta != 5.0 {
				t.Errorf("Expected row delta +5, got %f", rowTrigger.Data.Results[0].Result.Delta)
			}
		}
	})

	// Advance to next week's Heavy Day
	advanceUserState(t, ts, userID)

	// =============================================================================
	// VALIDATION PHASE: Next week's Heavy Day should show increased weights
	// =============================================================================
	t.Run("Next week Heavy Day shows +5lb progression on all lifts", func(t *testing.T) {
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

		// Should be Heavy Day again
		if workout.Data.DaySlug != "heavy" {
			t.Errorf("Expected day slug 'heavy', got '%s'", workout.Data.DaySlug)
		}

		exercisesByLift := make(map[string]WorkoutExerciseData)
		for _, ex := range workout.Data.Exercises {
			exercisesByLift[ex.Lift.ID] = ex
		}

		// Squat should now be at 305 (300 + 5)
		newSquatMax := squatMax + 5.0
		if squat, ok := exercisesByLift[squatID]; ok {
			topSet := squat.Sets[len(squat.Sets)-1]
			if !closeEnough(topSet.Weight, newSquatMax) {
				t.Errorf("Week 2 Squat top set: expected %.1f, got %.1f", newSquatMax, topSet.Weight)
			}
			// Verify all ramping sets increased proportionally
			expectedPercentages := []float64{0.50, 0.63, 0.75, 0.88, 1.00}
			for i, set := range squat.Sets {
				expectedWeight := newSquatMax * expectedPercentages[i]
				if !closeEnough(set.Weight, expectedWeight) {
					t.Errorf("Week 2 Squat set %d: expected weight ~%.1f, got %.1f", i+1, expectedWeight, set.Weight)
				}
			}
		}

		// Bench should now be at 205 (200 + 5)
		newBenchMax := benchMax + 5.0
		if bench, ok := exercisesByLift[benchID]; ok {
			topSet := bench.Sets[len(bench.Sets)-1]
			if !closeEnough(topSet.Weight, newBenchMax) {
				t.Errorf("Week 2 Bench top set: expected %.1f, got %.1f", newBenchMax, topSet.Weight)
			}
		}

		// Row should now be at 190 (185 + 5)
		newRowMax := rowMax + 5.0
		if row, ok := exercisesByLift[rowID]; ok {
			topSet := row.Sets[len(row.Sets)-1]
			if !closeEnough(topSet.Weight, newRowMax) {
				t.Errorf("Week 2 Row top set: expected %.1f, got %.1f", newRowMax, topSet.Weight)
			}
		}
	})
}

// =============================================================================
// HELPER TYPES AND FUNCTIONS
// =============================================================================

// RampStep represents a single step in a ramp progression for test creation.
type RampStep struct {
	Percentage float64 `json:"percentage"`
	Reps       int     `json:"reps"`
}

// createRampPrescription creates a prescription with RAMP set scheme.
func createRampPrescription(t *testing.T, ts *testutil.TestServer, liftID string, steps []RampStep, workSetThreshold float64, order int, lookupKey string) string {
	t.Helper()

	stepsJSON, err := json.Marshal(steps)
	if err != nil {
		t.Fatalf("Failed to marshal ramp steps: %v", err)
	}

	body := fmt.Sprintf(`{
		"liftId": "%s",
		"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 100.0, "lookupKey": "%s"},
		"setScheme": {"type": "RAMP", "steps": %s, "workSetThreshold": %.1f},
		"order": %d
	}`, liftID, lookupKey, string(stepsJSON), workSetThreshold, order)

	resp, err := adminPost(ts.URL("/prescriptions"), body)
	if err != nil {
		t.Fatalf("Failed to create ramp prescription: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create ramp prescription, status %d: %s", resp.StatusCode, bodyBytes)
	}

	var envelope PrescriptionResponse
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

// closeEnough checks if two floats are close enough (within 1.0 tolerance for rounding).
func closeEnough(a, b float64) bool {
	return math.Abs(a-b) < 1.0
}

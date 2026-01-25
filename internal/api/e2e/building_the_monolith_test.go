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
// 5/3/1 BUILDING THE MONOLITH E2E TEST
// =============================================================================

// TestBuildingTheMonolithProgram validates the complete 5/3/1 Building the Monolith
// program configuration and execution through the API.
//
// Building the Monolith characteristics:
// - 6-week cycle: Two 3-week blocks (same percentages, TM increases in block 2)
// - 3 days/week: Monday (Squat+Press), Wednesday (Deadlift+Bench), Friday (Squat+Press volume)
// - High volume main work: Multiple sets at working percentages
// - Friday Widowmaker: 1x20 squat at lower percentage
// - Press AMRAP on Monday: Final set is 5+
// - Friday Press volume: 10 sets of 5 at moderate percentage
// - CycleProgression: After 3 weeks, +5lb upper, +10lb lower
func TestBuildingTheMonolithProgram(t *testing.T) {
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
	pressSlug := "press-btm-" + testID
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

	// Building the Monolith training maxes (85% of 1RM typically)
	squatTM := 315.0    // Squat training max
	benchTM := 225.0    // Bench training max
	deadliftTM := 365.0 // Deadlift training max
	pressTM := 145.0    // Press training max

	// Create training maxes for the user
	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", squatTM)
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", benchTM)
	createLiftMax(t, ts, userID, deadliftID, "TRAINING_MAX", deadliftTM)
	createLiftMax(t, ts, userID, pressID, "TRAINING_MAX", pressTM)

	// =============================================================================
	// Create prescriptions for each day
	// BTM uses fixed percentages for work sets (90% for week 1)
	// =============================================================================

	// Monday Squat: 5x5 at 90% of TM (work sets only, warmups handled separately in real app)
	mondaySquatPrescID := createPrescription(t, ts, squatID, 5, 5, 90.0, 0)

	// Monday Press: AMRAP at 70% of TM
	mondayPressAMRAPPrescID := createAMRAPPrescription(t, ts, pressID, 1, 5, 70.0, 1)

	// Wednesday Deadlift: 5x5 at 90% of TM
	wednesdayDeadliftPrescID := createPrescription(t, ts, deadliftID, 5, 5, 90.0, 0)

	// Wednesday Bench: 5x5 at 90% of TM
	wednesdayBenchPrescID := createPrescription(t, ts, benchID, 5, 5, 90.0, 1)

	// Friday Squat: 3x5 at 90% (warmup/ramp-up sets)
	fridaySquatPrescID := createPrescription(t, ts, squatID, 3, 5, 90.0, 0)

	// Friday Widowmaker: 1x20 at 45% of TM
	fridayWidowmakerPrescID := createPrescription(t, ts, squatID, 1, 20, 45.0, 1)

	// Friday Press volume: 10x5 at 72% of TM
	fridayPressVolumePrescID := createPrescription(t, ts, pressID, 10, 5, 72.0, 2)

	// =============================================================================
	// Create Days
	// =============================================================================

	// Monday: Squat + Press (AMRAP)
	mondaySlug := "monday-btm-" + testID
	mondayBody := fmt.Sprintf(`{"name": "Monday - Squat/Press", "slug": "%s"}`, mondaySlug)
	mondayResp, _ := adminPost(ts.URL("/days"), mondayBody)
	var mondayEnvelope DayResponse
	json.NewDecoder(mondayResp.Body).Decode(&mondayEnvelope)
	mondayResp.Body.Close()
	mondayID := mondayEnvelope.Data.ID

	addPrescToDay(t, ts, mondayID, mondaySquatPrescID)
	addPrescToDay(t, ts, mondayID, mondayPressAMRAPPrescID)

	// Wednesday: Deadlift + Bench
	wednesdaySlug := "wednesday-btm-" + testID
	wednesdayBody := fmt.Sprintf(`{"name": "Wednesday - Deadlift/Bench", "slug": "%s"}`, wednesdaySlug)
	wednesdayResp, _ := adminPost(ts.URL("/days"), wednesdayBody)
	var wednesdayEnvelope DayResponse
	json.NewDecoder(wednesdayResp.Body).Decode(&wednesdayEnvelope)
	wednesdayResp.Body.Close()
	wednesdayID := wednesdayEnvelope.Data.ID

	addPrescToDay(t, ts, wednesdayID, wednesdayDeadliftPrescID)
	addPrescToDay(t, ts, wednesdayID, wednesdayBenchPrescID)

	// Friday: Squat (with Widowmaker) + Press volume
	fridaySlug := "friday-btm-" + testID
	fridayBody := fmt.Sprintf(`{"name": "Friday - Squat/Press Volume", "slug": "%s"}`, fridaySlug)
	fridayResp, _ := adminPost(ts.URL("/days"), fridayBody)
	var fridayEnvelope DayResponse
	json.NewDecoder(fridayResp.Body).Decode(&fridayEnvelope)
	fridayResp.Body.Close()
	fridayID := fridayEnvelope.Data.ID

	addPrescToDay(t, ts, fridayID, fridaySquatPrescID)
	addPrescToDay(t, ts, fridayID, fridayWidowmakerPrescID)
	addPrescToDay(t, ts, fridayID, fridayPressVolumePrescID)

	// =============================================================================
	// Create 3-week cycle (one block of BTM)
	// =============================================================================
	cycleName := "BTM Cycle " + testID
	cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 3}`, cycleName)
	cycleResp, _ := adminPost(ts.URL("/cycles"), cycleBody)
	var cycleEnvelope CycleResponse
	json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
	cycleResp.Body.Close()
	cycleID := cycleEnvelope.Data.ID

	// Create all 3 weeks
	weekIDs := make([]string, 3)
	for w := 1; w <= 3; w++ {
		weekBody := fmt.Sprintf(`{"weekNumber": %d, "cycleId": "%s"}`, w, cycleID)
		weekResp, _ := adminPost(ts.URL("/weeks"), weekBody)
		var weekEnvelope WeekResponse
		json.NewDecoder(weekResp.Body).Decode(&weekEnvelope)
		weekResp.Body.Close()
		weekIDs[w-1] = weekEnvelope.Data.ID

		// Add days to each week (Mon/Wed/Fri)
		addDayToWeek(t, ts, weekIDs[w-1], mondayID, "MONDAY")
		addDayToWeek(t, ts, weekIDs[w-1], wednesdayID, "WEDNESDAY")
		addDayToWeek(t, ts, weekIDs[w-1], fridayID, "FRIDAY")
	}

	// =============================================================================
	// Create Program
	// =============================================================================
	programSlug := "btm-531-" + testID
	programBody := fmt.Sprintf(`{"name": "5/3/1 Building the Monolith", "slug": "%s", "cycleId": "%s"}`,
		programSlug, cycleID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID := programEnvelope.Data.ID

	// =============================================================================
	// Create Cycle Progressions
	// =============================================================================

	// Lower body: +10lb per cycle
	lowerProgBody := `{"name": "BTM Lower +10lb", "type": "CYCLE_PROGRESSION", "parameters": {"increment": 10.0, "maxType": "TRAINING_MAX"}}`
	lowerProgResp, _ := adminPost(ts.URL("/progressions"), lowerProgBody)
	var lowerProgEnvelope ProgressionResponse
	json.NewDecoder(lowerProgResp.Body).Decode(&lowerProgEnvelope)
	lowerProgResp.Body.Close()
	lowerProgID := lowerProgEnvelope.Data.ID

	// Upper body: +5lb per cycle
	upperProgBody := `{"name": "BTM Upper +5lb", "type": "CYCLE_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX"}}`
	upperProgResp, _ := adminPost(ts.URL("/progressions"), upperProgBody)
	var upperProgEnvelope ProgressionResponse
	json.NewDecoder(upperProgResp.Body).Decode(&upperProgEnvelope)
	upperProgResp.Body.Close()
	upperProgID := upperProgEnvelope.Data.ID

	// Link progressions to program
	linkProgressionToProgram(t, ts, programID, lowerProgID, squatID, 1)
	linkProgressionToProgram(t, ts, programID, lowerProgID, deadliftID, 2)
	linkProgressionToProgram(t, ts, programID, upperProgID, benchID, 3)
	linkProgressionToProgram(t, ts, programID, upperProgID, pressID, 4)

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
	// WEEK 1: Monday - Test Squat 5x5 at 90% and Press AMRAP
	// =============================================================================
	t.Run("Week 1 Monday generates squat 5x5 at 90% and press AMRAP", func(t *testing.T) {
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

		// Verify it's Week 1, Monday
		if workout.Data.WeekNumber != 1 {
			t.Errorf("Expected week 1, got %d", workout.Data.WeekNumber)
		}

		if workout.Data.DaySlug != mondaySlug {
			t.Errorf("Expected day slug '%s', got '%s'", mondaySlug, workout.Data.DaySlug)
		}

		if len(workout.Data.Exercises) < 2 {
			t.Fatalf("Expected at least 2 exercises, got %d", len(workout.Data.Exercises))
		}

		exercisesByLift := make(map[string]WorkoutExerciseData)
		for _, ex := range workout.Data.Exercises {
			exercisesByLift[ex.Lift.ID] = ex
		}

		// Verify squat 5x5 at 90% of TM
		// Expected weight: 315 * 0.90 = 283.5 → rounded to 285
		if squat, ok := exercisesByLift[squatID]; ok {
			if len(squat.Sets) != 5 {
				t.Errorf("Squat: expected 5 work sets, got %d sets", len(squat.Sets))
			}
			expectedWeight := 285.0 // 315 * 0.90 = 283.5 → 285
			for i, set := range squat.Sets {
				if !withinTolerance(set.Weight, expectedWeight, 5.0) {
					t.Errorf("Squat set %d: expected weight ~%.1f (90%% of TM), got %.1f", i+1, expectedWeight, set.Weight)
				}
				if set.TargetReps != 5 {
					t.Errorf("Squat set %d: expected 5 reps, got %d", i+1, set.TargetReps)
				}
			}
		} else {
			t.Error("Monday missing Squat exercise")
		}

		// Verify press AMRAP set at 70%
		if press, ok := exercisesByLift[pressID]; ok {
			if len(press.Sets) != 1 {
				t.Errorf("Press: expected 1 AMRAP set, got %d sets", len(press.Sets))
			}
			// AMRAP at 70% of 145 = 101.5 → rounded to 100
			expectedWeight := 100.0
			if len(press.Sets) > 0 {
				if !withinTolerance(press.Sets[0].Weight, expectedWeight, 5.0) {
					t.Errorf("Press AMRAP: expected weight ~%.1f (70%% of TM), got %.1f", expectedWeight, press.Sets[0].Weight)
				}
				if press.Sets[0].TargetReps != 5 {
					t.Errorf("Press AMRAP: expected target reps 5, got %d", press.Sets[0].TargetReps)
				}
			}
		} else {
			t.Error("Monday missing Press exercise")
		}
	})

	// Complete Monday workout using explicit state machine flow
	completeBTMWorkoutDay(t, ts, userID)

	// =============================================================================
	// WEEK 1: Wednesday - Test Deadlift and Bench 5x5 at 90%
	// =============================================================================
	t.Run("Week 1 Wednesday generates deadlift and bench 5x5 at 90%", func(t *testing.T) {
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

		// Verify it's Wednesday
		if workout.Data.DaySlug != wednesdaySlug {
			t.Errorf("Expected day slug '%s', got '%s'", wednesdaySlug, workout.Data.DaySlug)
		}

		if len(workout.Data.Exercises) < 2 {
			t.Fatalf("Expected at least 2 exercises, got %d", len(workout.Data.Exercises))
		}

		exercisesByLift := make(map[string]WorkoutExerciseData)
		for _, ex := range workout.Data.Exercises {
			exercisesByLift[ex.Lift.ID] = ex
		}

		// Verify deadlift 5x5 at 90%
		// Expected weight: 365 * 0.90 = 328.5 → rounded to 330
		if deadlift, ok := exercisesByLift[deadliftID]; ok {
			if len(deadlift.Sets) != 5 {
				t.Errorf("Deadlift: expected 5 work sets, got %d sets", len(deadlift.Sets))
			}
			expectedWeight := 330.0 // 365 * 0.90 = 328.5 → 330
			for i, set := range deadlift.Sets {
				if !withinTolerance(set.Weight, expectedWeight, 5.0) {
					t.Errorf("Deadlift set %d: expected weight ~%.1f (90%% of TM), got %.1f", i+1, expectedWeight, set.Weight)
				}
				if set.TargetReps != 5 {
					t.Errorf("Deadlift set %d: expected 5 reps, got %d", i+1, set.TargetReps)
				}
			}
		} else {
			t.Error("Wednesday missing Deadlift exercise")
		}

		// Verify bench 5x5 at 90%
		// Expected weight: 225 * 0.90 = 202.5 → rounded to 205
		if bench, ok := exercisesByLift[benchID]; ok {
			if len(bench.Sets) != 5 {
				t.Errorf("Bench: expected 5 work sets, got %d sets", len(bench.Sets))
			}
			expectedWeight := 205.0 // 225 * 0.90 = 202.5 → 205
			for i, set := range bench.Sets {
				if !withinTolerance(set.Weight, expectedWeight, 5.0) {
					t.Errorf("Bench set %d: expected weight ~%.1f (90%% of TM), got %.1f", i+1, expectedWeight, set.Weight)
				}
				if set.TargetReps != 5 {
					t.Errorf("Bench set %d: expected 5 reps, got %d", i+1, set.TargetReps)
				}
			}
		} else {
			t.Error("Wednesday missing Bench exercise")
		}
	})

	// Complete Wednesday workout using explicit state machine flow
	completeBTMWorkoutDay(t, ts, userID)

	// =============================================================================
	// WEEK 1: Friday - Test Widowmaker (1x20 at 45%) and Press volume (10x5)
	// =============================================================================
	t.Run("Week 1 Friday generates widowmaker and press volume", func(t *testing.T) {
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

		// Verify it's Friday
		if workout.Data.DaySlug != fridaySlug {
			t.Errorf("Expected day slug '%s', got '%s'", fridaySlug, workout.Data.DaySlug)
		}

		// Collect exercises by lift (may have multiple prescriptions per lift)
		var squatExercises []WorkoutExerciseData
		var pressExercises []WorkoutExerciseData
		for _, ex := range workout.Data.Exercises {
			if ex.Lift.ID == squatID {
				squatExercises = append(squatExercises, ex)
			}
			if ex.Lift.ID == pressID {
				pressExercises = append(pressExercises, ex)
			}
		}

		// Should have 2 squat prescriptions (warmup + widowmaker)
		if len(squatExercises) != 2 {
			t.Errorf("Friday: expected 2 squat prescriptions, got %d", len(squatExercises))
		}

		// Find the widowmaker set (1x20)
		foundWidowmaker := false
		for _, squat := range squatExercises {
			for _, set := range squat.Sets {
				if set.TargetReps == 20 {
					foundWidowmaker = true
					// Expected weight: 315 * 0.45 = 141.75 → rounded to 140
					expectedWeight := 140.0
					if !withinTolerance(set.Weight, expectedWeight, 5.0) {
						t.Errorf("Widowmaker: expected weight ~%.1f (45%% of TM), got %.1f", expectedWeight, set.Weight)
					}
				}
			}
		}
		if !foundWidowmaker {
			t.Error("Friday missing Widowmaker 1x20 set")
		}

		// Verify press volume (10x5 at 72%)
		if len(pressExercises) < 1 {
			t.Error("Friday missing Press volume exercise")
		} else {
			press := pressExercises[0]
			if len(press.Sets) != 10 {
				t.Errorf("Press volume: expected 10 sets, got %d sets", len(press.Sets))
			}
			// Expected weight: 145 * 0.72 = 104.4 → rounded to 105
			expectedWeight := 105.0
			for i, set := range press.Sets {
				if !withinTolerance(set.Weight, expectedWeight, 5.0) {
					t.Errorf("Press volume set %d: expected weight ~%.1f (72%% of TM), got %.1f", i+1, expectedWeight, set.Weight)
				}
				if set.TargetReps != 5 {
					t.Errorf("Press volume set %d: expected 5 reps, got %d", i+1, set.TargetReps)
				}
			}
		}
	})

	// Complete rest of the cycle using explicit state machine flow
	// Week 1: Friday
	completeBTMWorkoutDay(t, ts, userID) // Week 1 Fri -> Week 2 Mon

	// Week 2: Mon/Wed/Fri
	completeBTMWorkoutDay(t, ts, userID) // Week 2 Mon -> Wed
	completeBTMWorkoutDay(t, ts, userID) // Week 2 Wed -> Fri
	completeBTMWorkoutDay(t, ts, userID) // Week 2 Fri -> Week 3 Mon

	// Week 3: Mon/Wed/Fri
	completeBTMWorkoutDay(t, ts, userID) // Week 3 Mon -> Wed
	completeBTMWorkoutDay(t, ts, userID) // Week 3 Wed -> Fri
	completeBTMWorkoutDay(t, ts, userID) // Week 3 Fri -> Cycle 2, Week 1 Mon

	// =============================================================================
	// CYCLE PROGRESSION: Trigger TM increases at end of cycle (no Force flag)
	// =============================================================================
	t.Run("Cycle progression triggers at cycle end", func(t *testing.T) {
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
				t.Errorf("Expected squat new value %f, got %f", expectedNewSquat, squatTrigger.Data.Results[0].Result.NewValue)
			}
		}

		// Trigger lower body progression for deadlift (+10lb)
		deadliftTrigger := triggerProgressionForLift(t, ts, userID, lowerProgID, deadliftID)

		if deadliftTrigger.Data.TotalApplied != 1 {
			t.Errorf("Expected deadlift progression to apply")
		}
		if len(deadliftTrigger.Data.Results) > 0 && deadliftTrigger.Data.Results[0].Result != nil {
			if deadliftTrigger.Data.Results[0].Result.Delta != 10.0 {
				t.Errorf("Expected deadlift delta +10, got %f", deadliftTrigger.Data.Results[0].Result.Delta)
			}
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
		}

		// Trigger upper body progression for press (+5lb)
		pressTrigger := triggerProgressionForLift(t, ts, userID, upperProgID, pressID)

		if pressTrigger.Data.TotalApplied != 1 {
			t.Errorf("Expected press progression to apply")
		}
		if len(pressTrigger.Data.Results) > 0 && pressTrigger.Data.Results[0].Result != nil {
			if pressTrigger.Data.Results[0].Result.Delta != 5.0 {
				t.Errorf("Expected press delta +5, got %f", pressTrigger.Data.Results[0].Result.Delta)
			}
		}
	})

	// =============================================================================
	// CYCLE 2: Verify new weights reflect TM increases
	// =============================================================================
	t.Run("Cycle 2 Week 1 shows increased training maxes", func(t *testing.T) {
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

		// Should be back to Week 1 (Monday)
		if workout.Data.WeekNumber != 1 {
			t.Errorf("Expected week 1 of cycle 2, got %d", workout.Data.WeekNumber)
		}

		exercisesByLift := make(map[string]WorkoutExerciseData)
		for _, ex := range workout.Data.Exercises {
			exercisesByLift[ex.Lift.ID] = ex
		}

		// Squat should now be at (315 + 10) * 0.90 = 325 * 0.90 = 292.5 → 295
		if squat, ok := exercisesByLift[squatID]; ok {
			expectedWeight := 295.0 // 325 * 0.90 = 292.5 → 295
			if len(squat.Sets) > 0 {
				if !withinTolerance(squat.Sets[0].Weight, expectedWeight, 5.0) {
					t.Errorf("Cycle 2 Squat: expected weight ~%.1f (90%% of new TM 325), got %.1f",
						expectedWeight, squat.Sets[0].Weight)
				}
			}
		}

		// Press AMRAP should now be at (145 + 5) * 0.70 = 150 * 0.70 = 105
		if press, ok := exercisesByLift[pressID]; ok {
			expectedWeight := 105.0 // 150 * 0.70 = 105
			if len(press.Sets) > 0 {
				if !withinTolerance(press.Sets[0].Weight, expectedWeight, 5.0) {
					t.Errorf("Cycle 2 Press: expected weight ~%.1f (70%% of new TM 150), got %.1f",
						expectedWeight, press.Sets[0].Weight)
				}
			}
		}
	})
}

// completeBTMWorkoutDay starts a workout session, logs all sets, finishes, and advances state.
// This is used for simple workout completion in Building the Monolith tests.
func completeBTMWorkoutDay(t *testing.T, ts *testutil.TestServer, userID string) {
	t.Helper()

	sessionID := startWorkoutSession(t, ts, userID)

	// Get the workout to find prescription IDs
	workoutResp, _ := userGet(ts.URL("/users/"+userID+"/workout"), userID)
	var workout WorkoutResponse
	json.NewDecoder(workoutResp.Body).Decode(&workout)
	workoutResp.Body.Close()

	// Log sets for each exercise
	for _, ex := range workout.Data.Exercises {
		for _, set := range ex.Sets {
			// Log successful completion (reps performed = target reps)
			loggedSetBody := fmt.Sprintf(`{
				"sets": [{
					"prescriptionId": "%s",
					"liftId": "%s",
					"setNumber": %d,
					"weight": %.1f,
					"targetReps": %d,
					"repsPerformed": %d,
					"isAmrap": false
				}]
			}`, ex.PrescriptionID, ex.Lift.ID, set.SetNumber, set.Weight, set.TargetReps, set.TargetReps)

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
	}

	finishWorkoutSession(t, ts, sessionID, userID)

	// Advance to next day (required until automatic triggering is implemented)
	advanceUserState(t, ts, userID)
}

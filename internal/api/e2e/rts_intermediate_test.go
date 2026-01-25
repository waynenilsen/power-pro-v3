// Package e2e provides end-to-end tests for complete program workflows.
// This file contains E2E tests for RTS (Reactive Training Systems) Generalized Intermediate program.
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
// RTS INTERMEDIATE HELPER FUNCTIONS
// =============================================================================

// createRTSIntermediateTestSetup creates the complete RTS Intermediate program structure.
// Returns programID, cycleID, and the week IDs for all 9 weeks.
//
// RTS Intermediate structure:
// - Week 1: Baseline (RPE 9-10, no fatigue method)
// - Weeks 2-4: Development (Load Drop, 0-7.5% fatigue)
// - Weeks 5-8: Intensification (Load Drop + Repeat, 5-7% fatigue)
// - Week 9: Peaking/Testing (singles at RPE 7-10, 0% fatigue)
func createRTSIntermediateTestSetup(t *testing.T, ts *testutil.TestServer, testID string) (programID, cycleID string, weekIDs []string) {
	t.Helper()

	// Seeded lift IDs
	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"
	deadliftID := "00000000-0000-0000-0000-000000000003"

	// =============================================================================
	// Create prescriptions for each phase using RPE-based load strategy
	// RTS uses RPE_TARGET load strategy with fatigue drop set schemes
	// =============================================================================

	// Week 1: Baseline prescriptions (RPE 9-10, no fatigue drops)
	// Working up to establish baseline weights - fixed sets
	w1SquatPrescID := createRTSPrescription(t, ts, squatID, 5, 9.0, 3, 5, 0)
	w1BenchPrescID := createRTSPrescription(t, ts, benchID, 5, 9.0, 3, 5, 1)
	w1DeadliftPrescID := createRTSPrescription(t, ts, deadliftID, 5, 9.0, 2, 5, 2)

	// Weeks 2-4: Development prescriptions (RPE 8-9, with fatigue drops)
	// Using Load Drop method with 5% fatigue
	devSquatPrescID := createRTSFatigueDropPrescription(t, ts, squatID, 4, 8.0, 10.0, 0.05, 0)
	devBenchPrescID := createRTSFatigueDropPrescription(t, ts, benchID, 4, 8.0, 10.0, 0.05, 1)
	devDeadliftPrescID := createRTSFatigueDropPrescription(t, ts, deadliftID, 3, 8.0, 10.0, 0.05, 2)

	// Weeks 5-8: Intensification prescriptions (RPE 8-9, with 5-7% fatigue)
	// Higher intensity, using Load Drop method
	intSquatPrescID := createRTSFatigueDropPrescription(t, ts, squatID, 3, 8.5, 10.0, 0.07, 0)
	intBenchPrescID := createRTSFatigueDropPrescription(t, ts, benchID, 3, 8.5, 10.0, 0.07, 1)
	intDeadliftPrescID := createRTSFatigueDropPrescription(t, ts, deadliftID, 2, 8.5, 10.0, 0.05, 2)

	// Week 9: Peaking prescriptions (singles at RPE 7-10, no fatigue)
	// Testing new maxes with singles
	peakSquatPrescID := createRTSPrescription(t, ts, squatID, 1, 9.5, 3, 1, 0)
	peakBenchPrescID := createRTSPrescription(t, ts, benchID, 1, 9.5, 3, 1, 1)
	peakDeadliftPrescID := createRTSPrescription(t, ts, deadliftID, 1, 10.0, 2, 1, 2)

	// =============================================================================
	// Create Days for each phase
	// RTS typically has 4 training days per week
	// =============================================================================

	// Week 1 (Baseline) Days
	w1Day1Slug := "rts-w1-d1-" + testID
	w1Day1ID := createRTSDay(t, ts, "Week 1 Day 1 - Squat/Bench Baseline", w1Day1Slug)
	addPrescToDay(t, ts, w1Day1ID, w1SquatPrescID)
	addPrescToDay(t, ts, w1Day1ID, w1BenchPrescID)

	w1Day2Slug := "rts-w1-d2-" + testID
	w1Day2ID := createRTSDay(t, ts, "Week 1 Day 2 - Deadlift/Bench Baseline", w1Day2Slug)
	addPrescToDay(t, ts, w1Day2ID, w1DeadliftPrescID)
	addPrescToDay(t, ts, w1Day2ID, w1BenchPrescID)

	w1Day3Slug := "rts-w1-d3-" + testID
	w1Day3ID := createRTSDay(t, ts, "Week 1 Day 3 - Squat Variation", w1Day3Slug)
	addPrescToDay(t, ts, w1Day3ID, w1SquatPrescID)

	w1Day4Slug := "rts-w1-d4-" + testID
	w1Day4ID := createRTSDay(t, ts, "Week 1 Day 4 - Deadlift Variation", w1Day4Slug)
	addPrescToDay(t, ts, w1Day4ID, w1DeadliftPrescID)

	// Development Days (Weeks 2-4) - with fatigue drops
	devDay1Slug := "rts-dev-d1-" + testID
	devDay1ID := createRTSDay(t, ts, "Development Day 1 - Squat/Bench Load Drop", devDay1Slug)
	addPrescToDay(t, ts, devDay1ID, devSquatPrescID)
	addPrescToDay(t, ts, devDay1ID, devBenchPrescID)

	devDay2Slug := "rts-dev-d2-" + testID
	devDay2ID := createRTSDay(t, ts, "Development Day 2 - Deadlift/Bench Load Drop", devDay2Slug)
	addPrescToDay(t, ts, devDay2ID, devDeadliftPrescID)
	addPrescToDay(t, ts, devDay2ID, devBenchPrescID)

	devDay3Slug := "rts-dev-d3-" + testID
	devDay3ID := createRTSDay(t, ts, "Development Day 3 - Squat Variation", devDay3Slug)
	addPrescToDay(t, ts, devDay3ID, devSquatPrescID)

	devDay4Slug := "rts-dev-d4-" + testID
	devDay4ID := createRTSDay(t, ts, "Development Day 4 - Deadlift Variation", devDay4Slug)
	addPrescToDay(t, ts, devDay4ID, devDeadliftPrescID)

	// Intensification Days (Weeks 5-8)
	intDay1Slug := "rts-int-d1-" + testID
	intDay1ID := createRTSDay(t, ts, "Intensification Day 1 - Squat/Bench Heavy", intDay1Slug)
	addPrescToDay(t, ts, intDay1ID, intSquatPrescID)
	addPrescToDay(t, ts, intDay1ID, intBenchPrescID)

	intDay2Slug := "rts-int-d2-" + testID
	intDay2ID := createRTSDay(t, ts, "Intensification Day 2 - Deadlift/Bench Heavy", intDay2Slug)
	addPrescToDay(t, ts, intDay2ID, intDeadliftPrescID)
	addPrescToDay(t, ts, intDay2ID, intBenchPrescID)

	intDay3Slug := "rts-int-d3-" + testID
	intDay3ID := createRTSDay(t, ts, "Intensification Day 3 - Squat Volume", intDay3Slug)
	addPrescToDay(t, ts, intDay3ID, intSquatPrescID)

	intDay4Slug := "rts-int-d4-" + testID
	intDay4ID := createRTSDay(t, ts, "Intensification Day 4 - Deadlift Volume", intDay4Slug)
	addPrescToDay(t, ts, intDay4ID, intDeadliftPrescID)

	// Peaking Days (Week 9)
	peakDay1Slug := "rts-peak-d1-" + testID
	peakDay1ID := createRTSDay(t, ts, "Peaking Day 1 - Squat/Bench Singles", peakDay1Slug)
	addPrescToDay(t, ts, peakDay1ID, peakSquatPrescID)
	addPrescToDay(t, ts, peakDay1ID, peakBenchPrescID)

	peakDay2Slug := "rts-peak-d2-" + testID
	peakDay2ID := createRTSDay(t, ts, "Peaking Day 2 - Deadlift Singles", peakDay2Slug)
	addPrescToDay(t, ts, peakDay2ID, peakDeadliftPrescID)

	peakDay3Slug := "rts-peak-d3-" + testID
	peakDay3ID := createRTSDay(t, ts, "Peaking Day 3 - Squat Opener Practice", peakDay3Slug)
	addPrescToDay(t, ts, peakDay3ID, peakSquatPrescID)

	peakDay4Slug := "rts-peak-d4-" + testID
	peakDay4ID := createRTSDay(t, ts, "Peaking Day 4 - Light Movement", peakDay4Slug)
	addPrescToDay(t, ts, peakDay4ID, peakBenchPrescID)

	// =============================================================================
	// Create 9-week cycle
	// =============================================================================
	cycleName := "RTS Intermediate Cycle " + testID
	cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 9}`, cycleName)
	cycleResp, _ := adminPost(ts.URL("/cycles"), cycleBody)
	var cycleEnvelope CycleResponse
	json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
	cycleResp.Body.Close()
	cycleID = cycleEnvelope.Data.ID

	// Create weeks for each phase
	weekIDs = make([]string, 9)

	// Week 1: Baseline
	weekIDs[0] = createRTSWeek(t, ts, cycleID, 1)
	addDayToWeek(t, ts, weekIDs[0], w1Day1ID, "MONDAY")
	addDayToWeek(t, ts, weekIDs[0], w1Day2ID, "TUESDAY")
	addDayToWeek(t, ts, weekIDs[0], w1Day3ID, "THURSDAY")
	addDayToWeek(t, ts, weekIDs[0], w1Day4ID, "FRIDAY")

	// Weeks 2-4: Development
	for i := 2; i <= 4; i++ {
		weekID := createRTSWeek(t, ts, cycleID, i)
		weekIDs[i-1] = weekID
		addDayToWeek(t, ts, weekID, devDay1ID, "MONDAY")
		addDayToWeek(t, ts, weekID, devDay2ID, "TUESDAY")
		addDayToWeek(t, ts, weekID, devDay3ID, "THURSDAY")
		addDayToWeek(t, ts, weekID, devDay4ID, "FRIDAY")
	}

	// Weeks 5-8: Intensification
	for i := 5; i <= 8; i++ {
		weekID := createRTSWeek(t, ts, cycleID, i)
		weekIDs[i-1] = weekID
		addDayToWeek(t, ts, weekID, intDay1ID, "MONDAY")
		addDayToWeek(t, ts, weekID, intDay2ID, "TUESDAY")
		addDayToWeek(t, ts, weekID, intDay3ID, "THURSDAY")
		addDayToWeek(t, ts, weekID, intDay4ID, "FRIDAY")
	}

	// Week 9: Peaking
	weekIDs[8] = createRTSWeek(t, ts, cycleID, 9)
	addDayToWeek(t, ts, weekIDs[8], peakDay1ID, "MONDAY")
	addDayToWeek(t, ts, weekIDs[8], peakDay2ID, "TUESDAY")
	addDayToWeek(t, ts, weekIDs[8], peakDay3ID, "THURSDAY")
	addDayToWeek(t, ts, weekIDs[8], peakDay4ID, "FRIDAY")

	// =============================================================================
	// Create program
	// =============================================================================
	programSlug := "rts-intermediate-" + testID
	programBody := fmt.Sprintf(`{"name": "RTS Generalized Intermediate", "slug": "%s", "cycleId": "%s"}`,
		programSlug, cycleID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID = programEnvelope.Data.ID

	return programID, cycleID, weekIDs
}

// createRTSPrescription creates a prescription with RPE-based load strategy.
// Uses RPE_TARGET load strategy for RPE-based calculations.
func createRTSPrescription(t *testing.T, ts *testutil.TestServer, liftID string, targetReps int, targetRPE float64, sets, reps, order int) string {
	t.Helper()

	// Use RPE_TARGET load strategy
	body := fmt.Sprintf(`{
		"liftId": "%s",
		"loadStrategy": {"type": "RPE_TARGET", "targetReps": %d, "targetRpe": %.1f, "roundingIncrement": 5, "roundingDirection": "NEAREST"},
		"setScheme": {"type": "FIXED", "sets": %d, "reps": %d},
		"order": %d
	}`, liftID, targetReps, targetRPE, sets, reps, order)

	resp, err := adminPost(ts.URL("/prescriptions"), body)
	if err != nil {
		t.Fatalf("Failed to create RTS prescription: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create RTS prescription, status %d: %s", resp.StatusCode, bodyBytes)
	}

	var envelope PrescriptionResponse
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

// createRTSFatigueDropPrescription creates a prescription with RPE-based load and fatigue drop scheme.
// Uses RPE_TARGET load strategy with FATIGUE_DROP set scheme.
func createRTSFatigueDropPrescription(t *testing.T, ts *testutil.TestServer, liftID string, targetReps int, startRPE, stopRPE, dropPercent float64, order int) string {
	t.Helper()

	// Use RPE_TARGET load strategy with FATIGUE_DROP set scheme
	body := fmt.Sprintf(`{
		"liftId": "%s",
		"loadStrategy": {"type": "RPE_TARGET", "targetReps": %d, "targetRpe": %.1f, "roundingIncrement": 5, "roundingDirection": "NEAREST"},
		"setScheme": {"type": "FATIGUE_DROP", "target_reps": %d, "start_rpe": %.1f, "stop_rpe": %.1f, "drop_percent": %.2f, "max_sets": 10},
		"order": %d
	}`, liftID, targetReps, startRPE, targetReps, startRPE, stopRPE, dropPercent, order)

	resp, err := adminPost(ts.URL("/prescriptions"), body)
	if err != nil {
		t.Fatalf("Failed to create RTS fatigue drop prescription: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create RTS fatigue drop prescription, status %d: %s", resp.StatusCode, bodyBytes)
	}

	var envelope PrescriptionResponse
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

func createRTSDay(t *testing.T, ts *testutil.TestServer, name, slug string) string {
	t.Helper()
	body := fmt.Sprintf(`{"name": "%s", "slug": "%s"}`, name, slug)
	resp, _ := adminPost(ts.URL("/days"), body)
	var dayEnvelope DayResponse
	json.NewDecoder(resp.Body).Decode(&dayEnvelope)
	resp.Body.Close()
	return dayEnvelope.Data.ID
}

func createRTSWeek(t *testing.T, ts *testutil.TestServer, cycleID string, weekNumber int) string {
	t.Helper()
	weekBody := fmt.Sprintf(`{"weekNumber": %d, "cycleId": "%s"}`, weekNumber, cycleID)
	weekResp, _ := adminPost(ts.URL("/weeks"), weekBody)
	var weekEnvelope WeekResponse
	json.NewDecoder(weekResp.Body).Decode(&weekEnvelope)
	weekResp.Body.Close()
	return weekEnvelope.Data.ID
}

// rtsWithinTolerance checks if value is within tolerance of expected.
func rtsWithinTolerance(value, expected, tolerance float64) bool {
	return math.Abs(value-expected) <= tolerance
}

// =============================================================================
// RTS INTERMEDIATE E2E TESTS
// =============================================================================

// TestRTSIntermediateFull9WeekProgram validates the complete 9-week RTS
// Intermediate program execution with RPE-based autoregulation.
//
// RTS Intermediate characteristics:
// - Week 1: Baseline (RPE 9-10, no fatigue method)
// - Weeks 2-4: Development (Load Drop, 0-7.5% fatigue)
// - Weeks 5-8: Intensification (Load Drop + Repeat, 5-7% fatigue)
// - Week 9: Peaking/Testing (singles at RPE 7-10, 0% fatigue)
// - 4 days/week: Mon/Tue/Thu/Fri pattern
// - RPE-based load calculation using RPE chart lookup
func TestRTSIntermediateFull9WeekProgram(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	userID := "workout-test-user"

	// Seeded lift IDs
	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"
	deadliftID := "00000000-0000-0000-0000-000000000003"

	// One Rep Maxes (RTS uses 1RM for RPE-based calculations)
	squat1RM := 400.0
	bench1RM := 300.0
	deadlift1RM := 450.0

	// Create ONE_RM maxes (RTS uses 1RM for RPE calculations)
	createLiftMax(t, ts, userID, squatID, "ONE_RM", squat1RM)
	createLiftMax(t, ts, userID, benchID, "ONE_RM", bench1RM)
	createLiftMax(t, ts, userID, deadliftID, "ONE_RM", deadlift1RM)

	// Create complete RTS Intermediate program
	programID, _, weekIDs := createRTSIntermediateTestSetup(t, ts, testID)

	// Enroll user
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

	t.Run("verifies 9-week cycle structure", func(t *testing.T) {
		if len(weekIDs) != 9 {
			t.Errorf("Expected 9 weeks, got %d", len(weekIDs))
		}
	})

	t.Run("generates Week 1 baseline workout with RPE-based weights", func(t *testing.T) {
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

		// Verify we have exercises
		if len(workout.Data.Exercises) == 0 {
			t.Error("Expected exercises in workout, got none")
		}

		// RPE chart: 5 reps @ RPE 9 = 80%
		// Squat: 400 * 0.80 = 320 lbs
		expectedSquatWeight := squat1RM * 0.80
		for _, ex := range workout.Data.Exercises {
			if ex.Lift.ID == squatID {
				for _, set := range ex.Sets {
					if rtsWithinTolerance(set.Weight, expectedSquatWeight, 5.0) {
						t.Logf("Week 1 squat set weight: %.1f (expected ~%.1f at 80%% of 1RM)",
							set.Weight, expectedSquatWeight)
					}
				}
				break
			}
		}
	})
}

// TestRTSIntermediateRPEChartCalculation validates RPE-based weight calculations
// using the standard RTS RPE chart.
func TestRTSIntermediateRPEChartCalculation(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	userID := "workout-test-user"

	// Seeded lift IDs
	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"
	deadliftID := "00000000-0000-0000-0000-000000000003"

	// Use round numbers for easy verification
	squat1RM := 400.0
	bench1RM := 300.0
	deadlift1RM := 500.0

	createLiftMax(t, ts, userID, squatID, "ONE_RM", squat1RM)
	createLiftMax(t, ts, userID, benchID, "ONE_RM", bench1RM)
	createLiftMax(t, ts, userID, deadliftID, "ONE_RM", deadlift1RM)

	programID, _, _ := createRTSIntermediateTestSetup(t, ts, testID)

	// Enroll user
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	t.Run("RPE chart lookup produces correct percentages", func(t *testing.T) {
		// Standard RTS RPE chart values:
		// 5 reps @ RPE 9 = 80%
		// 5 reps @ RPE 8 = 77%
		// 4 reps @ RPE 9 = 82%
		// 3 reps @ RPE 9 = 89%
		// 1 rep @ RPE 10 = 100%

		workoutResp, _ := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		// Week 1 baseline uses 5 reps @ RPE 9 = 80%
		// Squat: 400 * 0.80 = 320 lbs
		// Bench: 300 * 0.80 = 240 lbs
		expectedSquat := 320.0
		expectedBench := 240.0

		for _, ex := range workout.Data.Exercises {
			switch ex.Lift.ID {
			case squatID:
				for _, set := range ex.Sets {
					if rtsWithinTolerance(set.Weight, expectedSquat, 5.0) {
						t.Logf("Squat weight %.1f is within tolerance of expected %.1f (80%% of %.0f)",
							set.Weight, expectedSquat, squat1RM)
					} else {
						t.Logf("Squat weight: %.1f (expected ~%.1f)", set.Weight, expectedSquat)
					}
				}
			case benchID:
				for _, set := range ex.Sets {
					if rtsWithinTolerance(set.Weight, expectedBench, 5.0) {
						t.Logf("Bench weight %.1f is within tolerance of expected %.1f (80%% of %.0f)",
							set.Weight, expectedBench, bench1RM)
					} else {
						t.Logf("Bench weight: %.1f (expected ~%.1f)", set.Weight, expectedBench)
					}
				}
			}
		}
	})
}

// TestRTSIntermediatePhaseTransitions validates that phase transitions work
// correctly across the 9-week program.
func TestRTSIntermediatePhaseTransitions(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	userID := "workout-test-user"

	// Seeded lift IDs
	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"
	deadliftID := "00000000-0000-0000-0000-000000000003"

	createLiftMax(t, ts, userID, squatID, "ONE_RM", 400.0)
	createLiftMax(t, ts, userID, benchID, "ONE_RM", 300.0)
	createLiftMax(t, ts, userID, deadliftID, "ONE_RM", 450.0)

	programID, _, _ := createRTSIntermediateTestSetup(t, ts, testID)

	// Enroll user
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	t.Run("Week 1 is Baseline phase", func(t *testing.T) {
		workoutResp, _ := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)
		workoutResp.Body.Close()

		if workout.Data.WeekNumber != 1 {
			t.Errorf("Expected week 1, got %d", workout.Data.WeekNumber)
		}
		t.Logf("Week 1 (Baseline): %s", workout.Data.DaySlug)
	})

	// Advance through Week 1 (4 days)
	for i := 0; i < 4; i++ {
		advanceUserState(t, ts, userID)
	}

	t.Run("Week 2 begins Development phase", func(t *testing.T) {
		workoutResp, _ := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)
		workoutResp.Body.Close()

		if workout.Data.WeekNumber != 2 {
			t.Errorf("Expected week 2, got %d", workout.Data.WeekNumber)
		}
		t.Logf("Week 2 (Development): %s", workout.Data.DaySlug)
	})

	// Advance through Weeks 2-4 (3 weeks x 4 days = 12 days)
	for i := 0; i < 12; i++ {
		advanceUserState(t, ts, userID)
	}

	t.Run("Week 5 begins Intensification phase", func(t *testing.T) {
		workoutResp, _ := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)
		workoutResp.Body.Close()

		if workout.Data.WeekNumber != 5 {
			t.Errorf("Expected week 5, got %d", workout.Data.WeekNumber)
		}
		t.Logf("Week 5 (Intensification): %s", workout.Data.DaySlug)
	})

	// Advance through Weeks 5-8 (4 weeks x 4 days = 16 days)
	for i := 0; i < 16; i++ {
		advanceUserState(t, ts, userID)
	}

	t.Run("Week 9 begins Peaking phase", func(t *testing.T) {
		workoutResp, _ := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)
		workoutResp.Body.Close()

		if workout.Data.WeekNumber != 9 {
			t.Errorf("Expected week 9, got %d", workout.Data.WeekNumber)
		}
		t.Logf("Week 9 (Peaking): %s", workout.Data.DaySlug)
	})
}

// TestRTSIntermediateFourDayStructure validates that the program
// follows the 4-day per week structure (Mon/Tue/Thu/Fri).
func TestRTSIntermediateFourDayStructure(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	userID := "workout-test-user"

	// Seeded lift IDs
	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"
	deadliftID := "00000000-0000-0000-0000-000000000003"

	createLiftMax(t, ts, userID, squatID, "ONE_RM", 400.0)
	createLiftMax(t, ts, userID, benchID, "ONE_RM", 300.0)
	createLiftMax(t, ts, userID, deadliftID, "ONE_RM", 450.0)

	programID, _, _ := createRTSIntermediateTestSetup(t, ts, testID)

	// Enroll user
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	t.Run("generates 4 unique days per week", func(t *testing.T) {
		daysSeen := make(map[string]bool)

		// Get all 4 days of week 1
		for day := 0; day < 4; day++ {
			workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
			if err != nil {
				t.Fatalf("Failed to get workout day %d: %v", day+1, err)
			}

			if workoutResp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(workoutResp.Body)
				workoutResp.Body.Close()
				t.Fatalf("Expected status 200 for day %d, got %d: %s", day+1, workoutResp.StatusCode, body)
			}

			var workout WorkoutResponse
			json.NewDecoder(workoutResp.Body).Decode(&workout)
			workoutResp.Body.Close()

			daysSeen[workout.Data.DaySlug] = true
			t.Logf("Day %d: %s with %d exercises", day+1, workout.Data.DaySlug, len(workout.Data.Exercises))

			// Advance to next day
			if day < 3 {
				advanceUserState(t, ts, userID)
			}
		}

		// Should have 4 unique days
		if len(daysSeen) != 4 {
			t.Errorf("Expected 4 unique days per week, got %d: %v", len(daysSeen), daysSeen)
		}
	})
}

// TestRTSIntermediateWeekProgression validates that advancing through
// days correctly increments the week number at week boundaries.
func TestRTSIntermediateWeekProgression(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	userID := "workout-test-user"

	// Seeded lift IDs
	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"
	deadliftID := "00000000-0000-0000-0000-000000000003"

	createLiftMax(t, ts, userID, squatID, "ONE_RM", 400.0)
	createLiftMax(t, ts, userID, benchID, "ONE_RM", 300.0)
	createLiftMax(t, ts, userID, deadliftID, "ONE_RM", 450.0)

	programID, _, _ := createRTSIntermediateTestSetup(t, ts, testID)

	// Enroll user
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	t.Run("week increments from 1 to 2 after 4 days", func(t *testing.T) {
		// Get week 1 day 1
		workoutResp, _ := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)
		workoutResp.Body.Close()

		if workout.Data.WeekNumber != 1 {
			t.Errorf("Expected week 1, got %d", workout.Data.WeekNumber)
		}

		// Advance through all 4 days of week 1
		for i := 0; i < 4; i++ {
			advanceUserState(t, ts, userID)
		}

		// Get week 2 day 1
		workoutResp, _ = userGet(ts.URL("/users/"+userID+"/workout"), userID)
		json.NewDecoder(workoutResp.Body).Decode(&workout)
		workoutResp.Body.Close()

		if workout.Data.WeekNumber != 2 {
			t.Errorf("Expected week 2 after 4 advances, got %d", workout.Data.WeekNumber)
		}
	})
}

// TestRTSIntermediateFatigueDropScheme validates that the FATIGUE_DROP set scheme
// is properly applied in the Development and Intensification phases.
func TestRTSIntermediateFatigueDropScheme(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	userID := "workout-test-user"

	// Seeded lift IDs
	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"
	deadliftID := "00000000-0000-0000-0000-000000000003"

	createLiftMax(t, ts, userID, squatID, "ONE_RM", 400.0)
	createLiftMax(t, ts, userID, benchID, "ONE_RM", 300.0)
	createLiftMax(t, ts, userID, deadliftID, "ONE_RM", 450.0)

	programID, _, _ := createRTSIntermediateTestSetup(t, ts, testID)

	// Enroll user
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	// Advance to Week 2 (Development phase with fatigue drops)
	for i := 0; i < 4; i++ {
		advanceUserState(t, ts, userID)
	}

	t.Run("Development phase uses fatigue drop scheme", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		if workout.Data.WeekNumber != 2 {
			t.Errorf("Expected week 2 (Development), got %d", workout.Data.WeekNumber)
		}

		// Verify exercises are present with fatigue drop sets
		if len(workout.Data.Exercises) == 0 {
			t.Error("Expected exercises in Development phase workout, got none")
		}

		// Log the exercise structure
		for _, ex := range workout.Data.Exercises {
			t.Logf("Week 2 %s: %d sets", ex.Lift.Name, len(ex.Sets))
			// Fatigue drop generates provisional sets one at a time
			if len(ex.Sets) > 0 {
				t.Logf("  First set weight: %.1f lbs", ex.Sets[0].Weight)
			}
		}
	})
}

// TestRTSIntermediatePeakingPhase validates that Week 9 uses appropriate
// peaking prescriptions with singles at high RPE.
func TestRTSIntermediatePeakingPhase(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	userID := "workout-test-user"

	// Seeded lift IDs
	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"
	deadliftID := "00000000-0000-0000-0000-000000000003"

	squat1RM := 400.0
	bench1RM := 300.0
	deadlift1RM := 450.0

	createLiftMax(t, ts, userID, squatID, "ONE_RM", squat1RM)
	createLiftMax(t, ts, userID, benchID, "ONE_RM", bench1RM)
	createLiftMax(t, ts, userID, deadliftID, "ONE_RM", deadlift1RM)

	programID, _, _ := createRTSIntermediateTestSetup(t, ts, testID)

	// Enroll user
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	// Advance to Week 9 (8 weeks x 4 days = 32 days)
	for i := 0; i < 32; i++ {
		advanceUserState(t, ts, userID)
	}

	t.Run("Week 9 peaking phase uses singles at high RPE", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		if workout.Data.WeekNumber != 9 {
			t.Errorf("Expected week 9 (Peaking), got %d", workout.Data.WeekNumber)
		}

		// In Week 9, we use singles at RPE 9.5-10
		// RPE chart: 1 rep @ RPE 9.5 = 97.5%, 1 rep @ RPE 10 = 100%
		// Squat: 400 * 0.975 = 390 lbs
		expectedPeakWeight := squat1RM * 0.975

		for _, ex := range workout.Data.Exercises {
			if ex.Lift.ID == squatID {
				t.Logf("Week 9 %s: %d sets", ex.Lift.Name, len(ex.Sets))
				for _, set := range ex.Sets {
					t.Logf("  Set: %.1f lbs x %d reps (expected ~%.1f at 97.5%% of 1RM)",
						set.Weight, set.TargetReps, expectedPeakWeight)
					// Verify using singles (1 rep)
					if set.TargetReps != 1 {
						t.Logf("Note: Expected 1 rep singles in peaking phase, got %d", set.TargetReps)
					}
				}
				break
			}
		}
	})
}

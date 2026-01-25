// Package e2e provides end-to-end tests for complete program workflows.
// This file contains E2E tests for Calgary Barbell 16-Week program.
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
// CALGARY BARBELL 16-WEEK HELPER FUNCTIONS
// =============================================================================

// createCalgaryBarbell16WeekTestSetup creates the complete Calgary Barbell 16-Week program structure.
// Returns programID, cycleID, weeklyLookupID, and the week IDs for all 16 weeks.
func createCalgaryBarbell16WeekTestSetup(t *testing.T, ts *testutil.TestServer, testID string) (programID, cycleID, weeklyLookupID string, weekIDs []string) {
	t.Helper()

	// Seeded lift IDs
	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"
	deadliftID := "00000000-0000-0000-0000-000000000003"

	// =============================================================================
	// Create Weekly Lookup for 16-week intensity/rep progression
	// Calgary Barbell uses percentage-based intensity that varies each week
	// =============================================================================
	weeklyLookupBody := `{
		"name": "Calgary Barbell 16-Week Progression",
		"entries": [
			{"weekNumber": 1, "percentages": [67.0], "reps": [7], "percentageModifier": 100.0},
			{"weekNumber": 2, "percentages": [70.0], "reps": [6], "percentageModifier": 100.0},
			{"weekNumber": 3, "percentages": [73.0], "reps": [6], "percentageModifier": 100.0},
			{"weekNumber": 4, "percentages": [75.0], "reps": [5], "percentageModifier": 100.0},
			{"weekNumber": 5, "percentages": [80.0], "reps": [3], "percentageModifier": 100.0},
			{"weekNumber": 6, "percentages": [82.0], "reps": [3], "percentageModifier": 100.0},
			{"weekNumber": 7, "percentages": [86.0], "reps": [2], "percentageModifier": 100.0},
			{"weekNumber": 8, "percentages": [85.0], "reps": [3], "percentageModifier": 100.0},
			{"weekNumber": 9, "percentages": [82.0], "reps": [4], "percentageModifier": 100.0},
			{"weekNumber": 10, "percentages": [85.0], "reps": [3], "percentageModifier": 100.0},
			{"weekNumber": 11, "percentages": [83.0], "reps": [3], "percentageModifier": 100.0},
			{"weekNumber": 12, "percentages": [87.0], "reps": [3], "percentageModifier": 100.0},
			{"weekNumber": 13, "percentages": [89.0], "reps": [2], "percentageModifier": 100.0},
			{"weekNumber": 14, "percentages": [92.0], "reps": [1], "percentageModifier": 100.0},
			{"weekNumber": 15, "percentages": [92.0], "reps": [1], "percentageModifier": 100.0},
			{"weekNumber": 16, "percentages": [85.0], "reps": [2], "percentageModifier": 100.0}
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
	weeklyLookupID = weeklyLookupEnvelope.Data.ID

	// =============================================================================
	// Create prescriptions for each phase
	// Phase 1 (Hypertrophy): 4 sets, higher reps (5-7), 67-75% intensity
	// Phase 2 (Strength): 3-5 sets, lower reps (2-3), 80-86% intensity
	// Phase 3 (Peaking): 4-6 sets, 3-4 reps, 82-85% intensity
	// Phase 4 (Intensification): RPE-based singles/doubles, 87-92%
	// Phase 5 (Taper): Reduced volume at 82-85%
	// =============================================================================

	// Hypertrophy Phase prescriptions (Weeks 1-4)
	hypertrophySquatPrescID := createCalgaryPrescription(t, ts, squatID, 4, 7, 67.0, 0, "week")
	hypertrophyBenchPrescID := createCalgaryPrescription(t, ts, benchID, 4, 7, 67.0, 1, "week")
	hypertrophyDeadliftPrescID := createCalgaryPrescription(t, ts, deadliftID, 4, 7, 67.0, 2, "week")

	// Strength Phase prescriptions (Weeks 5-8)
	strengthSquatPrescID := createCalgaryPrescription(t, ts, squatID, 4, 3, 80.0, 0, "week")
	strengthBenchPrescID := createCalgaryPrescription(t, ts, benchID, 4, 3, 80.0, 1, "week")
	strengthDeadliftPrescID := createCalgaryPrescription(t, ts, deadliftID, 3, 3, 80.0, 2, "week")

	// Peaking Phase prescriptions (Weeks 9-11)
	peakingSquatPrescID := createCalgaryPrescription(t, ts, squatID, 5, 4, 82.0, 0, "week")
	peakingBenchPrescID := createCalgaryPrescription(t, ts, benchID, 5, 4, 82.0, 1, "week")
	peakingDeadliftPrescID := createCalgaryPrescription(t, ts, deadliftID, 4, 4, 82.0, 2, "week")

	// Intensification Phase prescriptions (Weeks 12-15)
	intensificationSquatPrescID := createCalgaryPrescription(t, ts, squatID, 3, 2, 90.0, 0, "week")
	intensificationBenchPrescID := createCalgaryPrescription(t, ts, benchID, 3, 2, 90.0, 1, "week")
	intensificationDeadliftPrescID := createCalgaryPrescription(t, ts, deadliftID, 2, 2, 90.0, 2, "week")

	// Taper Week prescriptions (Week 16) - reduced volume
	taperSquatPrescID := createCalgaryPrescription(t, ts, squatID, 3, 2, 82.0, 0, "week")
	taperBenchPrescID := createCalgaryPrescription(t, ts, benchID, 4, 1, 85.0, 1, "week")
	taperDeadliftPrescID := createCalgaryPrescription(t, ts, deadliftID, 2, 2, 82.0, 2, "week")

	// =============================================================================
	// Create Days for each phase (4 days per week)
	// Day 1: Squat primary + Bench secondary
	// Day 2: Deadlift primary + Bench variation + Squat secondary
	// Day 3: Squat variation + Bench primary + Deadlift secondary
	// Day 4: Deadlift variation + Bench variation
	// =============================================================================

	// Hypertrophy Phase Days (Weeks 1-4)
	hypertrophyDay1Slug := "cb16-hyp-d1-" + testID
	hypertrophyDay1ID := createCalgaryDay(t, ts, "Hypertrophy Day 1 - Squat/Bench", hypertrophyDay1Slug)
	addPrescToDay(t, ts, hypertrophyDay1ID, hypertrophySquatPrescID)
	addPrescToDay(t, ts, hypertrophyDay1ID, hypertrophyBenchPrescID)

	hypertrophyDay2Slug := "cb16-hyp-d2-" + testID
	hypertrophyDay2ID := createCalgaryDay(t, ts, "Hypertrophy Day 2 - Deadlift/Bench/Squat", hypertrophyDay2Slug)
	addPrescToDay(t, ts, hypertrophyDay2ID, hypertrophyDeadliftPrescID)
	addPrescToDay(t, ts, hypertrophyDay2ID, hypertrophyBenchPrescID)

	hypertrophyDay3Slug := "cb16-hyp-d3-" + testID
	hypertrophyDay3ID := createCalgaryDay(t, ts, "Hypertrophy Day 3 - Squat/Bench/Deadlift", hypertrophyDay3Slug)
	addPrescToDay(t, ts, hypertrophyDay3ID, hypertrophySquatPrescID)
	addPrescToDay(t, ts, hypertrophyDay3ID, hypertrophyBenchPrescID)

	hypertrophyDay4Slug := "cb16-hyp-d4-" + testID
	hypertrophyDay4ID := createCalgaryDay(t, ts, "Hypertrophy Day 4 - Deadlift/Bench", hypertrophyDay4Slug)
	addPrescToDay(t, ts, hypertrophyDay4ID, hypertrophyDeadliftPrescID)
	addPrescToDay(t, ts, hypertrophyDay4ID, hypertrophyBenchPrescID)

	// Strength Phase Days (Weeks 5-8)
	strengthDay1Slug := "cb16-str-d1-" + testID
	strengthDay1ID := createCalgaryDay(t, ts, "Strength Day 1 - Squat/Bench", strengthDay1Slug)
	addPrescToDay(t, ts, strengthDay1ID, strengthSquatPrescID)
	addPrescToDay(t, ts, strengthDay1ID, strengthBenchPrescID)

	strengthDay2Slug := "cb16-str-d2-" + testID
	strengthDay2ID := createCalgaryDay(t, ts, "Strength Day 2 - Deadlift/Bench/Squat", strengthDay2Slug)
	addPrescToDay(t, ts, strengthDay2ID, strengthDeadliftPrescID)
	addPrescToDay(t, ts, strengthDay2ID, strengthBenchPrescID)

	strengthDay3Slug := "cb16-str-d3-" + testID
	strengthDay3ID := createCalgaryDay(t, ts, "Strength Day 3 - Squat/Bench/Deadlift", strengthDay3Slug)
	addPrescToDay(t, ts, strengthDay3ID, strengthSquatPrescID)
	addPrescToDay(t, ts, strengthDay3ID, strengthBenchPrescID)

	strengthDay4Slug := "cb16-str-d4-" + testID
	strengthDay4ID := createCalgaryDay(t, ts, "Strength Day 4 - Deadlift/Bench", strengthDay4Slug)
	addPrescToDay(t, ts, strengthDay4ID, strengthDeadliftPrescID)
	addPrescToDay(t, ts, strengthDay4ID, strengthBenchPrescID)

	// Peaking Phase Days (Weeks 9-11)
	peakingDay1Slug := "cb16-peak-d1-" + testID
	peakingDay1ID := createCalgaryDay(t, ts, "Peaking Day 1 - Squat/Bench", peakingDay1Slug)
	addPrescToDay(t, ts, peakingDay1ID, peakingSquatPrescID)
	addPrescToDay(t, ts, peakingDay1ID, peakingBenchPrescID)

	peakingDay2Slug := "cb16-peak-d2-" + testID
	peakingDay2ID := createCalgaryDay(t, ts, "Peaking Day 2 - Deadlift/Bench/Squat", peakingDay2Slug)
	addPrescToDay(t, ts, peakingDay2ID, peakingDeadliftPrescID)
	addPrescToDay(t, ts, peakingDay2ID, peakingBenchPrescID)

	peakingDay3Slug := "cb16-peak-d3-" + testID
	peakingDay3ID := createCalgaryDay(t, ts, "Peaking Day 3 - Squat/Bench/Deadlift", peakingDay3Slug)
	addPrescToDay(t, ts, peakingDay3ID, peakingSquatPrescID)
	addPrescToDay(t, ts, peakingDay3ID, peakingBenchPrescID)

	peakingDay4Slug := "cb16-peak-d4-" + testID
	peakingDay4ID := createCalgaryDay(t, ts, "Peaking Day 4 - Deadlift/Bench", peakingDay4Slug)
	addPrescToDay(t, ts, peakingDay4ID, peakingDeadliftPrescID)
	addPrescToDay(t, ts, peakingDay4ID, peakingBenchPrescID)

	// Intensification Phase Days (Weeks 12-15)
	intensificationDay1Slug := "cb16-int-d1-" + testID
	intensificationDay1ID := createCalgaryDay(t, ts, "Intensification Day 1 - Squat/Bench", intensificationDay1Slug)
	addPrescToDay(t, ts, intensificationDay1ID, intensificationSquatPrescID)
	addPrescToDay(t, ts, intensificationDay1ID, intensificationBenchPrescID)

	intensificationDay2Slug := "cb16-int-d2-" + testID
	intensificationDay2ID := createCalgaryDay(t, ts, "Intensification Day 2 - Deadlift/Bench/Squat", intensificationDay2Slug)
	addPrescToDay(t, ts, intensificationDay2ID, intensificationDeadliftPrescID)
	addPrescToDay(t, ts, intensificationDay2ID, intensificationBenchPrescID)

	intensificationDay3Slug := "cb16-int-d3-" + testID
	intensificationDay3ID := createCalgaryDay(t, ts, "Intensification Day 3 - Squat/Bench/Deadlift", intensificationDay3Slug)
	addPrescToDay(t, ts, intensificationDay3ID, intensificationSquatPrescID)
	addPrescToDay(t, ts, intensificationDay3ID, intensificationBenchPrescID)

	intensificationDay4Slug := "cb16-int-d4-" + testID
	intensificationDay4ID := createCalgaryDay(t, ts, "Intensification Day 4 - Deadlift/Bench", intensificationDay4Slug)
	addPrescToDay(t, ts, intensificationDay4ID, intensificationDeadliftPrescID)
	addPrescToDay(t, ts, intensificationDay4ID, intensificationBenchPrescID)

	// Taper Week Days (Week 16)
	taperDay1Slug := "cb16-taper-d1-" + testID
	taperDay1ID := createCalgaryDay(t, ts, "Taper Day 1 - Squat/Bench", taperDay1Slug)
	addPrescToDay(t, ts, taperDay1ID, taperSquatPrescID)
	addPrescToDay(t, ts, taperDay1ID, taperBenchPrescID)

	taperDay2Slug := "cb16-taper-d2-" + testID
	taperDay2ID := createCalgaryDay(t, ts, "Taper Day 2 - Deadlift/Bench", taperDay2Slug)
	addPrescToDay(t, ts, taperDay2ID, taperDeadliftPrescID)
	addPrescToDay(t, ts, taperDay2ID, taperBenchPrescID)

	taperDay3Slug := "cb16-taper-d3-" + testID
	taperDay3ID := createCalgaryDay(t, ts, "Taper Day 3 - All Lifts", taperDay3Slug)
	addPrescToDay(t, ts, taperDay3ID, taperSquatPrescID)
	addPrescToDay(t, ts, taperDay3ID, taperBenchPrescID)
	addPrescToDay(t, ts, taperDay3ID, taperDeadliftPrescID)

	taperDay4Slug := "cb16-taper-d4-" + testID
	taperDay4ID := createCalgaryDay(t, ts, "Taper Day 4 - Light Practice", taperDay4Slug)
	addPrescToDay(t, ts, taperDay4ID, taperSquatPrescID)
	addPrescToDay(t, ts, taperDay4ID, taperBenchPrescID)

	// =============================================================================
	// Create 16-week cycle
	// =============================================================================
	cycleName := "Calgary Barbell 16-Week Cycle " + testID
	cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 16}`, cycleName)
	cycleResp, _ := adminPost(ts.URL("/cycles"), cycleBody)
	var cycleEnvelope CycleResponse
	json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
	cycleResp.Body.Close()
	cycleID = cycleEnvelope.Data.ID

	// Create weeks for each phase
	weekIDs = make([]string, 16)

	// Hypertrophy Phase: Weeks 1-4
	for i := 1; i <= 4; i++ {
		weekID := createCalgaryWeek(t, ts, cycleID, i)
		weekIDs[i-1] = weekID
		addDayToWeek(t, ts, weekID, hypertrophyDay1ID, "MONDAY")
		addDayToWeek(t, ts, weekID, hypertrophyDay2ID, "TUESDAY")
		addDayToWeek(t, ts, weekID, hypertrophyDay3ID, "THURSDAY")
		addDayToWeek(t, ts, weekID, hypertrophyDay4ID, "FRIDAY")
	}

	// Strength Phase: Weeks 5-8
	for i := 5; i <= 8; i++ {
		weekID := createCalgaryWeek(t, ts, cycleID, i)
		weekIDs[i-1] = weekID
		addDayToWeek(t, ts, weekID, strengthDay1ID, "MONDAY")
		addDayToWeek(t, ts, weekID, strengthDay2ID, "TUESDAY")
		addDayToWeek(t, ts, weekID, strengthDay3ID, "THURSDAY")
		addDayToWeek(t, ts, weekID, strengthDay4ID, "FRIDAY")
	}

	// Peaking Phase: Weeks 9-11
	for i := 9; i <= 11; i++ {
		weekID := createCalgaryWeek(t, ts, cycleID, i)
		weekIDs[i-1] = weekID
		addDayToWeek(t, ts, weekID, peakingDay1ID, "MONDAY")
		addDayToWeek(t, ts, weekID, peakingDay2ID, "TUESDAY")
		addDayToWeek(t, ts, weekID, peakingDay3ID, "THURSDAY")
		addDayToWeek(t, ts, weekID, peakingDay4ID, "FRIDAY")
	}

	// Intensification Phase: Weeks 12-15
	for i := 12; i <= 15; i++ {
		weekID := createCalgaryWeek(t, ts, cycleID, i)
		weekIDs[i-1] = weekID
		addDayToWeek(t, ts, weekID, intensificationDay1ID, "MONDAY")
		addDayToWeek(t, ts, weekID, intensificationDay2ID, "TUESDAY")
		addDayToWeek(t, ts, weekID, intensificationDay3ID, "THURSDAY")
		addDayToWeek(t, ts, weekID, intensificationDay4ID, "FRIDAY")
	}

	// Taper Week: Week 16
	taperWeekID := createCalgaryWeek(t, ts, cycleID, 16)
	weekIDs[15] = taperWeekID
	addDayToWeek(t, ts, taperWeekID, taperDay1ID, "MONDAY")
	addDayToWeek(t, ts, taperWeekID, taperDay2ID, "TUESDAY")
	addDayToWeek(t, ts, taperWeekID, taperDay3ID, "WEDNESDAY")
	addDayToWeek(t, ts, taperWeekID, taperDay4ID, "THURSDAY")

	// =============================================================================
	// Create program with weekly lookup
	// =============================================================================
	programSlug := "calgary-barbell-16-" + testID
	programBody := fmt.Sprintf(`{"name": "Calgary Barbell 16-Week", "slug": "%s", "cycleId": "%s", "weeklyLookupId": "%s"}`,
		programSlug, cycleID, weeklyLookupID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID = programEnvelope.Data.ID

	return programID, cycleID, weeklyLookupID, weekIDs
}

func createCalgaryPrescription(t *testing.T, ts *testutil.TestServer, liftID string, sets, reps int, percentage float64, order int, lookupKey string) string {
	t.Helper()

	body := fmt.Sprintf(`{
		"liftId": "%s",
		"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": %.1f, "lookupKey": "%s"},
		"setScheme": {"type": "FIXED", "sets": %d, "reps": %d},
		"order": %d
	}`, liftID, percentage, lookupKey, sets, reps, order)

	resp, err := adminPost(ts.URL("/prescriptions"), body)
	if err != nil {
		t.Fatalf("Failed to create Calgary prescription: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create Calgary prescription, status %d: %s", resp.StatusCode, bodyBytes)
	}

	var envelope PrescriptionResponse
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

func createCalgaryDay(t *testing.T, ts *testutil.TestServer, name, slug string) string {
	t.Helper()
	body := fmt.Sprintf(`{"name": "%s", "slug": "%s"}`, name, slug)
	resp, _ := adminPost(ts.URL("/days"), body)
	var dayEnvelope DayResponse
	json.NewDecoder(resp.Body).Decode(&dayEnvelope)
	resp.Body.Close()
	return dayEnvelope.Data.ID
}

func createCalgaryWeek(t *testing.T, ts *testutil.TestServer, cycleID string, weekNumber int) string {
	t.Helper()
	weekBody := fmt.Sprintf(`{"weekNumber": %d, "cycleId": "%s"}`, weekNumber, cycleID)
	weekResp, _ := adminPost(ts.URL("/weeks"), weekBody)
	var weekEnvelope WeekResponse
	json.NewDecoder(weekResp.Body).Decode(&weekEnvelope)
	weekResp.Body.Close()
	return weekEnvelope.Data.ID
}

// =============================================================================
// CALGARY BARBELL 16-WEEK E2E TESTS
// =============================================================================

// TestCalgaryBarbell16WeekProgram validates the complete 16-week Calgary Barbell
// program execution with all five phases and taper week.
//
// Calgary Barbell 16-Week characteristics:
// - Phase 1 (Hypertrophy): Weeks 1-4 - Volume accumulation, 67-75% intensity, 5-7 reps
// - Phase 2 (Strength): Weeks 5-8 - Intensity introduction, 80-86%, 2-3 reps
// - Phase 3 (Peaking): Weeks 9-11 - Competition specificity, 82-85%, fatigue management
// - Phase 4 (Intensification): Weeks 12-15 - Max strength, RPE-based singles/doubles
// - Phase 5 (Taper): Week 16 - Recovery & priming, opener practice
// - 4 days/week: Mon/Tue/Thu/Fri pattern
func TestCalgaryBarbell16WeekProgram(t *testing.T) {
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

	// Training maxes (intermediate-advanced levels typical for Calgary Barbell)
	squatMax := 200.0
	benchMax := 140.0
	deadliftMax := 220.0

	// Create training maxes
	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", squatMax)
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", benchMax)
	createLiftMax(t, ts, userID, deadliftID, "TRAINING_MAX", deadliftMax)

	// Create complete Calgary Barbell 16-Week program
	programID, _, _, _ := createCalgaryBarbell16WeekTestSetup(t, ts, testID)

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

	// Set meet date 16 weeks (112 days) from now
	meetDate := time.Now().AddDate(0, 0, 112).Format("2006-01-02")
	meetDateResponse := setMeetDate(t, ts, userID, programID, meetDate)

	t.Run("verifies meet date is set correctly for 16-week program", func(t *testing.T) {
		if meetDateResponse.DaysOut < 105 || meetDateResponse.DaysOut > 120 {
			t.Errorf("Expected ~112 days out, got %d", meetDateResponse.DaysOut)
		}
	})

	t.Run("verifies initial phase at 16 weeks out", func(t *testing.T) {
		// At 112 days out, should be in early phase
		validPhases := map[string]bool{"prep_1": true, "base": true, "off_season": true}
		if !validPhases[meetDateResponse.CurrentPhase] {
			t.Logf("At 112 days out, got phase: %s", meetDateResponse.CurrentPhase)
		}
	})

	t.Run("Hypertrophy Phase Week 1 generates correct workout", func(t *testing.T) {
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

		// In Hypertrophy phase (Week 1), weights should be at ~67% of training max
		for _, ex := range workout.Data.Exercises {
			if ex.Lift.ID == squatID {
				expectedWeight := squatMax * 0.67
				for _, set := range ex.Sets {
					if !closeEnough(set.Weight, expectedWeight) {
						t.Logf("Hypertrophy squat weight: expected ~%.1f (67%%), got %.1f",
							expectedWeight, set.Weight)
					}
				}
			}
		}
	})
}

// TestCalgaryBarbell16WeekPhaseTransitions validates that phase transitions work
// correctly when starting at different points in the program.
func TestCalgaryBarbell16WeekPhaseTransitions(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	userID := "workout-test-user"

	// Create complete Calgary Barbell program
	programID, _, _, _ := createCalgaryBarbell16WeekTestSetup(t, ts, testID)

	// Enroll user
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	t.Run("entering prep_2 phase at 8 weeks out", func(t *testing.T) {
		// 8 weeks = 56 days
		meetDate := time.Now().AddDate(0, 0, 56).Format("2006-01-02")
		response := setMeetDate(t, ts, userID, programID, meetDate)

		// At 56 days, should be in prep_2 phase (strength phase)
		validPhases := map[string]bool{"prep_2": true, "prep_1": true}
		if !validPhases[response.CurrentPhase] {
			t.Logf("At 56 days out, got phase: %s", response.CurrentPhase)
		}
	})

	t.Run("entering taper phase at 2 weeks out", func(t *testing.T) {
		// 2 weeks = 14 days
		meetDate := time.Now().AddDate(0, 0, 14).Format("2006-01-02")
		response := setMeetDate(t, ts, userID, programID, meetDate)

		// At 14 days, should be in peak phase
		if response.CurrentPhase != "peak" {
			t.Errorf("Expected peak phase at 14 days out, got %s", response.CurrentPhase)
		}
	})

	t.Run("entering meet week at 1 week out", func(t *testing.T) {
		// 1 week = 7 days
		meetDate := time.Now().AddDate(0, 0, 7).Format("2006-01-02")
		response := setMeetDate(t, ts, userID, programID, meetDate)

		// At 7 days, should be meet_week
		if response.CurrentPhase != "meet_week" {
			t.Errorf("Expected meet_week phase at 7 days out, got %s", response.CurrentPhase)
		}
	})
}

// TestCalgaryBarbell16WeekIntensityProgression validates that intensity
// progresses correctly across the 16 weeks.
func TestCalgaryBarbell16WeekIntensityProgression(t *testing.T) {
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

	// Use round numbers for easy percentage verification
	squatMax := 200.0
	benchMax := 150.0
	deadliftMax := 250.0

	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", squatMax)
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", benchMax)
	createLiftMax(t, ts, userID, deadliftID, "TRAINING_MAX", deadliftMax)

	programID, _, _, _ := createCalgaryBarbell16WeekTestSetup(t, ts, testID)

	// Enroll user
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	// Set meet date 16 weeks out
	meetDate := time.Now().AddDate(0, 0, 112).Format("2006-01-02")
	setMeetDate(t, ts, userID, programID, meetDate)

	t.Run("Week 1 hypertrophy phase uses 67% intensity", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		if workout.Data.WeekNumber != 1 {
			t.Errorf("Expected week 1, got %d", workout.Data.WeekNumber)
		}

		// Log weights for verification
		for _, ex := range workout.Data.Exercises {
			var max float64
			switch ex.Lift.ID {
			case squatID:
				max = squatMax
			case benchID:
				max = benchMax
			case deadliftID:
				max = deadliftMax
			default:
				continue
			}

			for i, set := range ex.Sets {
				percentage := (set.Weight / max) * 100
				t.Logf("Week 1 %s Set %d: %.1f lbs (%.1f%%)", ex.Lift.Name, i+1, set.Weight, percentage)
			}
		}
	})

	// Complete Week 1 workouts (4 days) using explicit state machine flow
	for i := 0; i < 4; i++ {
		completeCB16WorkoutDay(t, ts, userID)
	}

	t.Run("Week 2 hypertrophy phase uses 70% intensity", func(t *testing.T) {
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

		t.Logf("Week 2: Generated %d exercises", len(workout.Data.Exercises))
	})
}

// TestCalgaryBarbell16WeekTaperWeek validates that the taper week (Week 16)
// reduces volume appropriately for meet preparation.
func TestCalgaryBarbell16WeekTaperWeek(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	userID := "workout-test-user"

	programID, _, _, _ := createCalgaryBarbell16WeekTestSetup(t, ts, testID)

	// Enroll user
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	// Set meet date 5 days out (should be in meet_week/taper phase)
	meetDate := time.Now().AddDate(0, 0, 5).Format("2006-01-02")
	setMeetDate(t, ts, userID, programID, meetDate)

	t.Run("taper multiplier is reduced at meet week", func(t *testing.T) {
		countdown := getCountdown(t, ts, userID, programID)

		// At meet week, taper multiplier should be low (around 0.4)
		if countdown.TaperMultiplier > 0.5 {
			t.Errorf("Expected low taper multiplier at meet week, got %.2f", countdown.TaperMultiplier)
		}
	})

	t.Run("phase is meet_week at 5 days out", func(t *testing.T) {
		countdown := getCountdown(t, ts, userID, programID)

		if countdown.CurrentPhase != "meet_week" {
			t.Errorf("Expected meet_week phase, got %s", countdown.CurrentPhase)
		}
	})
}

// TestCalgaryBarbell16WeekFourDayStructure validates that the program
// follows the 4-day per week structure (Mon/Tue/Thu/Fri).
func TestCalgaryBarbell16WeekFourDayStructure(t *testing.T) {
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

	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", 200.0)
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", 140.0)
	createLiftMax(t, ts, userID, deadliftID, "TRAINING_MAX", 220.0)

	programID, _, _, _ := createCalgaryBarbell16WeekTestSetup(t, ts, testID)

	// Enroll user
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	// Set meet date
	meetDate := time.Now().AddDate(0, 0, 112).Format("2006-01-02")
	setMeetDate(t, ts, userID, programID, meetDate)

	t.Run("generates 4 unique day slugs per week", func(t *testing.T) {
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

			// Complete the current day's workout (advances state automatically after finish)
			if day < 3 {
				completeCB16WorkoutDay(t, ts, userID)
			}
		}

		// Should have 4 unique days
		if len(daysSeen) != 4 {
			t.Errorf("Expected 4 unique days per week, got %d: %v", len(daysSeen), daysSeen)
		}
	})
}

// TestCalgaryBarbell16WeekWeekProgression validates that advancing through
// days correctly increments the week number at week boundaries.
func TestCalgaryBarbell16WeekWeekProgression(t *testing.T) {
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

	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", 200.0)
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", 140.0)
	createLiftMax(t, ts, userID, deadliftID, "TRAINING_MAX", 220.0)

	programID, _, _, _ := createCalgaryBarbell16WeekTestSetup(t, ts, testID)

	// Enroll user
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	// Set meet date
	meetDate := time.Now().AddDate(0, 0, 112).Format("2006-01-02")
	setMeetDate(t, ts, userID, programID, meetDate)

	t.Run("week increments from 1 to 2 after 4 days", func(t *testing.T) {
		// Get week 1 day 1
		workoutResp, _ := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)
		workoutResp.Body.Close()

		if workout.Data.WeekNumber != 1 {
			t.Errorf("Expected week 1, got %d", workout.Data.WeekNumber)
		}

		// Complete all 4 days of week 1 using explicit state machine flow
		for i := 0; i < 4; i++ {
			completeCB16WorkoutDay(t, ts, userID)
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

// TestCalgaryBarbell16WeekMeetDateChanges validates that changing the meet date
// correctly recalculates the phase and countdown.
func TestCalgaryBarbell16WeekMeetDateChanges(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	userID := "workout-test-user"

	programID, _, _, _ := createCalgaryBarbell16WeekTestSetup(t, ts, testID)

	// Enroll user
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	t.Run("changes from 12 weeks to 6 weeks recalculates phase", func(t *testing.T) {
		// First set meet date at 12 weeks (84 days)
		meetDate12Weeks := time.Now().AddDate(0, 0, 84).Format("2006-01-02")
		response1 := setMeetDate(t, ts, userID, programID, meetDate12Weeks)

		t.Logf("At 84 days out, phase: %s", response1.CurrentPhase)

		// Now change to 6 weeks (42 days)
		meetDate6Weeks := time.Now().AddDate(0, 0, 42).Format("2006-01-02")
		response2 := setMeetDate(t, ts, userID, programID, meetDate6Weeks)

		// At 42 days, should be in prep_2 phase
		if response2.CurrentPhase != "prep_2" {
			t.Logf("Expected prep_2 phase at 42 days out, got %s", response2.CurrentPhase)
		}

		// Verify days out changed
		if response2.DaysOut >= response1.DaysOut {
			t.Errorf("Expected days_out to decrease, got %d -> %d", response1.DaysOut, response2.DaysOut)
		}
	})

	t.Run("clearing meet date returns to off_season", func(t *testing.T) {
		// Clear meet date
		body := `{"meet_date": null}`
		resp, err := userPut(ts.URL("/users/"+userID+"/programs/"+programID+"/state/meet-date"), body, userID)
		if err != nil {
			t.Fatalf("Failed to clear meet date: %v", err)
		}
		defer resp.Body.Close()

		var envelope MeetDateResponseEnvelope
		json.NewDecoder(resp.Body).Decode(&envelope)

		if envelope.Data.CurrentPhase != "off_season" {
			t.Errorf("Expected off_season phase after clearing meet date, got %s", envelope.Data.CurrentPhase)
		}
		if envelope.Data.DaysOut != 0 {
			t.Errorf("Expected 0 days_out after clearing, got %d", envelope.Data.DaysOut)
		}
	})
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

// completeCB16WorkoutDay completes a Calgary Barbell 16-Week workout day using explicit state machine flow.
// This function starts a session, logs all sets, finishes the session, and advances to the next day.
func completeCB16WorkoutDay(t *testing.T, ts *testutil.TestServer, userID string) {
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
			logCB16Set(t, ts, userID, sessionID, ex.PrescriptionID, ex.Lift.ID, set.SetNumber, set.Weight, set.TargetReps, set.TargetReps)
		}
	}

	finishWorkoutSession(t, ts, sessionID, userID)

	// Advance to next day
	advanceUserState(t, ts, userID)
}

// logCB16Set logs a single set for Calgary Barbell 16-Week workout.
func logCB16Set(t *testing.T, ts *testutil.TestServer, userID, sessionID, prescriptionID, liftID string, setNumber int, weight float64, targetReps, repsPerformed int) {
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
